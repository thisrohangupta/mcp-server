# Pipeline Execution Rerun — Spec
<!-- github-spec-generator:v2 source=features/feat-pipelinererun.md source-blob=1930215570bf9c9b048d8edb909669eeb68c1593 -->

## Problem

- User: Developers rerunning failed CI or CD executions, operators responding to transient deployment failures, AI agents working from a Harness execution URL, and MCP clients that need one-call execution and status monitoring.
- Blocked workflow or outcome: Pipeline retry is registered on the pipeline resource but targets an execution, requiring `execution_id` through `params`. Harness execution URLs resolve to `resource_type="execution"`, which does not expose a retry action. A natural request such as "rerun this execution" does not map to the current contract. When the native retry endpoint returns HTTP 405, the server may start a fresh pipeline run that is not equivalent to replaying the original execution inputs and source context.
- Why now: *Not specified in the feature document or PR description — to be defined*

## Solution

- Proposed experience: Expose a declarative `retry` action on the `execution` resource so that a pipeline execution can be rerun directly from an execution ID or a Harness execution URL in a single MCP call, while preserving the previous execution's effective inputs where possible and keeping the existing `pipeline.retry` contract fully compatible.
- Key behaviors and capabilities:
    - `harness_execute(resource_type="execution", action="retry", resource_id="<execution-id>")` retries a pipeline execution by its ID.
    - `harness_execute(url="<Harness execution URL>", action="retry")` extracts all required identifiers from the URL automatically.
    - `wait=true` polls the new execution (not the source) until a terminal state or timeout.
    - The server first attempts the native Harness retry endpoint (`PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`); on HTTP 405 it falls back to a fresh pipeline run using the recovered execution context.
    - When fallback is used, the response reports `rerun_mode="fresh_run_fallback"` and a note that native retry was unavailable.
    - If native retry fails and the pipeline ID or required inputs cannot be recovered, the server returns an actionable error and starts no new execution.
    - The operation uses `risk: high_write` and `retryPolicy: do_not_retry`; it respects `HARNESS_READ_ONLY` and the existing confirmation/elicitation flow.
    - `harness_describe(resource_type="execution")` advertises the retry action with its identifier requirements, risk classification, body-less contract, `wait` support, and possible fallback.
    - Existing `resource_type="pipeline", action="retry"` calls continue to work and share the same underlying rerun logic.

- In scope:
    - Rerunning a pipeline execution by execution ID via `resource_type="execution", action="retry"`.
    - Rerunning a pipeline execution from a Harness execution URL.
    - Native retry via `PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`.
    - Fresh-run fallback when native retry returns HTTP 405 and the pipeline context is recoverable.
    - Safe refusal (actionable error, no new execution) when required context is unrecoverable.
    - `wait=true` support for the new execution.
    - Backward compatibility for `resource_type="pipeline", action="retry"`.
    - Discovery metadata in `harness_describe` for the execution resource.

- Out of scope:
    - Retrying an individual stage or step.
    - Resuming paused executions.
    - Replaying dynamic pipeline executions.
    - Adding a new top-level MCP tool.
    - Guaranteeing that an external Git repository still contains the exact historical pipeline definition.
    - Changing the behavior of normal `pipeline.run` requests.

## Value

 Audience                                                                              Value
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Developers and operators rerunning failed CI/CD executions                            Rerun a failed execution from its URL in one call without manually extracting IDs.
 AI agents working from a Harness execution URL                                        Natural "rerun this execution" intent maps directly to a single `harness_execute` call.
 MCP clients requiring one-call execution and status monitoring                        `wait=true` delivers terminal status in a single tool call, eliminating polling loops.

## Metrics

 Category          Metric                                             Target / Direction
━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Discoverability   Pipeline execution can be rerun from Harness URL   One MCP call, no manual ID extraction required
 Correctness       Successful rerun responses contain new execution ID Always present
 Transparency      Fresh-run fallbacks are clearly identified           `rerun_mode` field present in every response
 Safety            Fallback does not start pipeline when context missing Zero incomplete fallback executions
 Compatibility     Existing `pipeline.retry` clients remain compatible  All existing tests continue to pass

## User Stories

 As a...                          I want...                                                     So that...                                                        Acceptance Criteria
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 developer                        rerun a failed execution by its ID                            I don't have to navigate the UI or extract pipeline parameters     Given a valid execution ID and `confirm=true`, the server returns the new execution ID (AC1)
 operator                         rerun an execution by pasting its Harness URL                 I can respond to transient failures without looking up identifiers  Given a Harness execution URL and `action="retry"`, all identifiers are extracted and the retry starts (AC2)
 developer                        be prompted before a rerun executes                            accidental reruns are prevented                                   Given a client without auto-approval, the rerun does not start until the high-write operation is confirmed (AC3)
 operator                         have reruns blocked in read-only mode                         policy is enforced before any mutation is sent                    Given `HARNESS_READ_ONLY=true`, the rerun is rejected before any mutation request is sent (AC4)
 developer                        receive a clear signal distinguishing native retry from fallback I can understand whether the rerun replayed the original context   Response always includes `rerun_mode` set to `native_retry` or `fresh_run_fallback` (AC5, AC6)
 MCP client                       use `wait=true` to get terminal status in one call            polling loops are unnecessary                                     Given `wait=true`, the server polls the new execution and returns terminal status or timeout information (AC8)
 existing integrator               continue using `resource_type="pipeline", action="retry"`    I do not need to update my integration                            Existing calls pass current tests and return the normalized rerun response (AC9)

## Dependencies and Open Questions

- Dependencies:
    - Pipeline retry API (`PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`).
    - Pipeline execution-details API.
    - Execution-inputs resource/API.
    - Existing URL parser.
    - Existing confirmation and audit infrastructure.
    - Existing server-side execution polling.

- Open questions or assumptions:
    - **Owner:** Listed as TBD in the feature document — assignee not yet defined.
    - **Why now:** No urgency or deadline is stated in the feature document or PR description — to be defined.
    - **Retry versus replay:** A native retry and a fresh pipeline run may behave differently because pipeline definitions, templates, expressions, secrets, connectors, and Git content can change after the source execution. The response must distinguish between the two modes.
    - **Duplicate execution risk:** Rerun is non-idempotent; retrying the HTTP request after an ambiguous network failure could create duplicate executions. `retryPolicy="do_not_retry"` must be retained.
    - **Input recovery:** Execution input data may be incomplete or unavailable; the implementation must fail safely rather than silently running with defaults.
    - **Contract duplication:** Temporarily supporting retry on both `pipeline` and `execution` resources introduces two discovery paths; shared implementation logic is required to prevent behavior drift.

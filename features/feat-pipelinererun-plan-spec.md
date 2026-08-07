# Pipeline Execution Rerun — Coding Plan — Spec
<!-- github-spec-generator:v2 source=features/feat-pipelinererun-plan.md source-blob=32e7732a9cfad31478d7c4fe11c1964d35eb573b -->

## Problem

- User: Developers, operators, AI agents, and MCP clients who need to rerun a Harness pipeline execution from its execution ID or a Harness execution URL.
- Blocked workflow or outcome: The current `execution` resource does not expose a `retry` action. Pipeline retry is registered on the `pipeline` resource, requiring `execution_id` through `params`, which is difficult to discover and does not align with how Harness execution URLs resolve. When the native retry endpoint returns HTTP 405, the server may start a fresh pipeline run that is not equivalent to replaying the original execution inputs and source context.
- Why now: *Not specified in the feature document or PR description — to be defined*

## Solution

- Proposed experience: Register a declarative `retry` execute action on the `execution` resource in the Harness MCP server toolset definition. The action resolves the execution ID from `resource_id` or a Harness URL, attempts the native Harness retry endpoint first, falls back to a fresh pipeline run on HTTP 405 when context is recoverable, refuses to start any new execution when context is unrecoverable, and preserves full backward compatibility with `resource_type="pipeline", action="retry"`.
- Key behaviors and capabilities:
    - A `retry` action is registered on the `execution` resource in the toolset definition (`src/registry/toolsets/`).
    - The action handler resolves the execution ID from `resource_id` or URL and calls `PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`.
    - On HTTP 405, the server falls back to a fresh pipeline run using recovered execution context.
    - Both the new `execution.retry` path and the existing `pipeline.retry` path delegate to shared rerun logic to prevent behavior drift.
    - Every response includes `rerun_mode` (`native_retry` or `fresh_run_fallback`) and `source_execution_id`.
    - `risk: high_write` and `retryPolicy: do_not_retry` are applied; the action respects `HARNESS_READ_ONLY` and the existing confirmation/elicitation flow.
    - `wait=true` polls the new execution ID, not the source execution.
    - `harness_describe(resource_type="execution")` advertises the `retry` action with identifier requirements, risk classification, body-less contract, `wait` support, and possible fallback.

- In scope:
    - Registering `retry` on the `execution` toolset definition.
    - Extracting and sharing rerun logic from the existing `pipeline.retry` implementation.
    - Native retry via `PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`.
    - Fresh-run fallback when native retry returns HTTP 405 and pipeline context is recoverable.
    - Safe refusal (actionable error, no new execution) when required context is unrecoverable.
    - Normalised response shape with `rerun_mode`, `source_execution_id`, and `execution_id` fields via shared response extractor.
    - `wait=true` support wired to the new execution ID.
    - URL parser confirmation that `execution` resource_type extraction and execution-ID mapping work correctly.
    - `harness_describe` metadata updates for the `execution` resource.
    - `risk: high_write`, `retryPolicy: do_not_retry`, `HARNESS_READ_ONLY` check, and confirmation/elicitation flow applied to `execution.retry`.
    - Unit and integration tests covering all 16 testing requirements from the feature document.
    - Documentation updates and `pnpm build && pnpm docs:generate` run.

- Out of scope:
    - Adding a new top-level MCP tool.
    - Removing `retry` from the `pipeline` resource (backward compatibility must be preserved).
    - Retrying an individual stage or step.
    - Resuming paused executions.
    - Replaying dynamic pipeline executions.
    - Guaranteeing that an external Git repository still contains the exact historical pipeline definition.

## Value

 Audience                                                                 Value
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Developers and operators rerunning failed CI/CD executions                Rerun a failed execution from its ID or URL in one call without manually extracting pipeline parameters.
 AI agents working from a Harness execution URL                            Natural "rerun this execution" intent maps directly to a single `harness_execute` call on the `execution` resource.
 MCP clients requiring one-call execution and status monitoring             `wait=true` delivers terminal status in a single tool call, eliminating polling loops.
 Existing integrators using `resource_type="pipeline", action="retry"`    No changes required; the legacy contract is preserved via shared rerun logic.

## Metrics

 Category          Metric                                                         Target / Direction
━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Correctness       All 16 testing requirements from the feature document covered   All unit and integration tests pass
 Transparency      `rerun_mode` field present in every rerun response               Always present (`native_retry` or `fresh_run_fallback`)
 Safety            No new execution started when required context is unrecoverable  Zero incomplete fallback executions
 Compatibility     Existing `pipeline.retry` callers continue to work               All existing tests continue to pass
 Stability         `docs:check` passes after implementation                         `pnpm build && pnpm docs:generate` completed before merge

## User Stories

 As a...               I want...                                                                So that...                                                              Acceptance Criteria
━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 developer             rerun a failed execution by its ID or URL in one MCP call               I don't have to navigate the UI, extract pipeline parameters, or construct manual params   Given a valid execution ID or Harness execution URL and `confirm=true`, the server returns the new execution ID (Tasks 1–9)
 implementer           extract shared rerun logic before registering the new action             the `pipeline.retry` and `execution.retry` paths do not drift in safety, fallback, or audit behavior   Task 1 is completed before Tasks 2 and 3 depend on it
 implementer           apply `risk: high_write` and `retryPolicy: do_not_retry` to all rerun paths   accidental or duplicate reruns are prevented at the policy layer   Task 11 is completed; `HARNESS_READ_ONLY` is checked before any mutation (Task 11)
 implementer           wire `wait=true` to poll the new execution ID, not the source execution  callers get correct terminal status without a separate polling loop     Task 8 is completed; the new execution ID is extracted and polled
 existing integrator   continue using `resource_type="pipeline", action="retry"` unchanged      no migration effort is required                                         Task 3 delegates to shared rerun logic; all existing tests pass (Task 12)

## Dependencies and Open Questions

- Dependencies:
    - Existing `pipeline.retry` implementation (source of shared rerun logic — Task 1).
    - `src/registry/toolsets/` execution and pipeline toolset definition files (Tasks 2, 3).
    - Shared rerun logic module to be created in `src/` (Tasks 4, 5, 6).
    - `src/registry/extractors.ts` or equivalent (Task 7).
    - URL parser module in `src/` (Task 9).
    - Harness `PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}` endpoint (Task 4).
    - Execution-details API and execution-inputs resource/API (Task 5).
    - Existing server-side execution polling infrastructure (Task 8).
    - `harness_describe` metadata registry (Task 10).
    - Existing confirmation and audit infrastructure (Task 11).
    - `tests/` directory and vitest tooling (Task 12).
    - `pnpm build && pnpm docs:generate` script (Task 13).

- Open questions or assumptions:
    - Exact file paths for the execution toolset definition, pipeline toolset definition, shared rerun module, URL parser, and response extractors are not confirmed — requires reading `src/registry/toolsets/` and related source directories.
    - Whether the fresh-run fallback should surface a confirmation prompt separate from the initial high-write confirmation, or whether one confirmation covers both paths, is not specified.
    - How execution-inputs are retrieved (API shape, required parameters) depends on the execution-inputs resource/API implementation — to be confirmed during implementation.
    - The timeout and cancellation behavior for `wait=true` on the new execution follows the existing wait infrastructure; any deviations must be documented.
    - The existing URL parser is assumed to already handle `execution` resource_type extraction; confirmation requires reading the parser source (Task 9).
    - The existing server-side execution polling used by `pipeline.run` with `wait=true` is assumed to be reusable for `execution.retry` with `wait=true` by substituting the new execution ID.

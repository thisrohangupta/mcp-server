# Pipeline Execution Rerun — Coding Plan
<!-- github-coding-plan-generator:v2 spec=features/feat-pipelinererun-spec.md source=features/feat-pipelinererun.md source-blob=1930215570bf9c9b048d8edb909669eeb68c1593 -->

> Generated from ./feat-pipelinererun-spec.md @ b6f8cbe. This is a proposed implementation plan.

## Overview

Add a declarative `retry` action to the `execution` resource in the Harness MCP server so that a pipeline execution can be rerun from its execution ID or a Harness execution URL in a single `harness_execute` call. The implementation must attempt the native Harness retry endpoint first, fall back to a fresh pipeline run when the endpoint returns HTTP 405 and context is recoverable, refuse to start any new execution when required context is unrecoverable, and preserve full backward compatibility with the existing `resource_type="pipeline", action="retry"` contract.

- Spec: ./feat-pipelinererun-spec.md
- Feature document: ./feat-pipelinererun.md

## Architecture and Approach

- High-level design: Register a `retry` execute action on the `execution` resource in the toolset definition (`src/registry/toolsets/`). The action handler resolves the execution ID from `resource_id` or URL, calls `PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`, and on HTTP 405 falls back to a pipeline run using recovered execution context. Both the new `execution.retry` path and the existing `pipeline.retry` path delegate to shared rerun logic to prevent behavior drift. The response always includes `rerun_mode` (`native_retry` or `fresh_run_fallback`) and `source_execution_id`.
- Technical decisions and tradeoffs: Shared rerun logic keeps safety, fallback, audit, and response semantics consistent across the two action paths. Registering retry on `execution` rather than as a new top-level tool keeps the tool count stable. `retryPolicy: do_not_retry` is non-negotiable because the operation is non-idempotent.
- Alternatives considered: Adding a new top-level MCP tool — excluded by the feature non-goals. Registering retry only on `execution` and removing it from `pipeline` — excluded to preserve backward compatibility (AC9).

## Affected Areas

 Area / Module                                        Change Type   Notes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 <to be determined> (execution toolset definition)    Modify        Add `retry` action with `risk: high_write`, `retryPolicy: do_not_retry`, no body, `wait` support
 <to be determined> (pipeline toolset definition)     Modify        Delegate `retry` action to shared rerun logic instead of inline implementation
 <to be determined> (shared rerun logic module)       New           Extract/create shared function: native retry attempt → 405 fallback → unsafe refusal
 <to be determined> (URL parser)                      Modify        Ensure `execution` resource type is recognised; extraction of execution ID confirmed working
 <to be determined> (harness_describe metadata)       Modify        Advertise `retry` action on `execution` resource with identifier, risk, wait, fallback metadata
 <to be determined> (response extractor)              Modify        Normalise native retry and fresh-run fallback response shapes; add `rerun_mode` field
 tests/                                               New           Add focused unit and integration tests for all 16 testing requirements from the feature document
 features/feat-pipelinererun-spec.md                  Docs          Already present in PR; no change required

## Work Breakdown

 #              1
 Task           Identify and extract shared rerun logic
 Files          <to be determined> (new shared rerun module in `src/`)
 Type           Backend
 Est. Effort    M
 Depends On     —

 #              2
 Task           Register `retry` action on the `execution` toolset definition
 Files          <to be determined> (`src/registry/toolsets/` — execution toolset file)
 Type           Backend
 Est. Effort    S
 Depends On     1

 #              3
 Task           Update `pipeline` toolset `retry` action to delegate to shared rerun logic
 Files          <to be determined> (`src/registry/toolsets/` — pipeline toolset file)
 Type           Backend
 Est. Effort    S
 Depends On     1

 #              4
 Task           Implement native retry attempt (`PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`)
 Files          <to be determined> (shared rerun module)
 Type           Backend
 Est. Effort    M
 Depends On     1

 #              5
 Task           Implement fresh-run fallback: resolve pipeline ID and effective inputs from execution details, start pipeline run, report `rerun_mode="fresh_run_fallback"`
 Files          <to be determined> (shared rerun module; execution-inputs API call)
 Type           Backend
 Est. Effort    L
 Depends On     4

 #              6
 Task           Implement unsafe-fallback prevention: return actionable error when pipeline ID or required inputs are unrecoverable, start no new execution
 Files          <to be determined> (shared rerun module)
 Type           Backend
 Est. Effort    S
 Depends On     5

 #              7
 Task           Normalise response shape: add `rerun_mode`, `source_execution_id`, `execution_id` fields via shared response extractor
 Files          <to be determined> (`src/registry/extractors.ts` or equivalent)
 Type           Backend
 Est. Effort    S
 Depends On     4, 5

 #              8
 Task           Wire `wait=true` to poll the new execution ID returned by the rerun, not the source execution
 Files          <to be determined> (shared rerun module or execute dispatcher)
 Type           Backend
 Est. Effort    M
 Depends On     4, 5

 #              9
 Task           Update URL parser to confirm `execution` resource_type extraction and execution-ID mapping from Harness execution URLs
 Files          <to be determined> (URL parser module in `src/`)
 Type           Backend
 Est. Effort    S
 Depends On     2

 #              10
 Task           Update `harness_describe` metadata for `execution` resource to advertise `retry` action, identifier requirements, risk, wait support, and fallback
 Files          <to be determined> (describe/registry metadata in `src/`)
 Type           Backend
 Est. Effort    S
 Depends On     2

 #              11
 Task           Apply `risk: high_write`, `retryPolicy: do_not_retry`, `HARNESS_READ_ONLY` check, and confirmation/elicitation flow to the `execution.retry` action
 Files          <to be determined> (shared rerun module or action policy layer)
 Type           Backend
 Est. Effort    S
 Depends On     2

 #              12
 Task           Write unit and integration tests covering all 16 testing requirements from the feature document (AC1–AC10)
 Files          tests/ (new test file for execution retry)
 Type           Test
 Est. Effort    L
 Depends On     2, 3, 4, 5, 6, 7, 8, 9, 10, 11

 #              13
 Task           Update tool/action descriptions, generated resource documentation, and example calls; run `pnpm build && pnpm docs:generate`
 Files          <to be determined> (docs/ or generated files updated by docs:generate script)
 Type           Docs
 Est. Effort    S
 Depends On     2, 7, 10

 ## Test Strategy

 Layer                 Coverage                                                                                   Tooling
━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Unit                  execution.retry registration and discovery (AC10); resource_id→execution_id mapping        vitest
                       (AC1); URL parsing and scope extraction (AC2); explicit values override URL-derived scope;
                       high-write confirmation gate (AC3); HARNESS_READ_ONLY rejection before I/O (AC4);
                       native retry success and response shape (AC5); 405 fallback response shape (AC6);
                       unsafe fallback refusal (AC7); wait=true polls new execution ID (AC8);
                       legacy pipeline.retry compatibility (AC9); audit record emission
────────────────────  ─────────────────────────────────────────────────────────────────────────────────────────  ────────────────────────────
 Integration           End-to-end rerun by execution ID; rerun from URL; wait=true terminal-status delivery;       vitest (mocked Harness API)
                       fresh-run fallback with recoverable context; error when context unrecoverable
────────────────────  ─────────────────────────────────────────────────────────────────────────────────────────  ────────────────────────────
 E2E                   *Not specified — to be defined during implementation*                                       <to be determined>
────────────────────  ─────────────────────────────────────────────────────────────────────────────────────────  ────────────────────────────
 Performance / Load    *Not specified — to be defined during implementation*                                       <to be determined>
────────────────────  ─────────────────────────────────────────────────────────────────────────────────────────  ────────────────────────────
 Security              HARNESS_READ_ONLY enforcement before any mutation; confirmation/elicitation required for     vitest
                       high_write; no execution started when context unrecoverable

## Rollout and Migration

- Rollout controls: *Not specified in the feature document or PR description — to be defined*
- Backward compatibility and migration: Existing `resource_type="pipeline", action="retry"` calls must continue to work without change (AC9). No changes to the `pipeline.run` contract. The new `execution.retry` action is additive; no existing callers are affected.
- Telemetry and observability: The feature document specifies that audit records must be emitted for successful, blocked, and fallback rerun paths, matching the audit behavior of other execute actions. Specific telemetry instrumentation beyond audit records is not specified.

## Risks and Mitigations

 Risk                                         Likelihood   Impact   Mitigation
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━  ━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 Native retry and fresh-run behave            High         High     Response always includes `rerun_mode` field to distinguish the two paths. Callers must not
 differently (definitions, secrets,                                 treat a fallback as an exact replay.
 connectors, or Git content changed
 after source execution)
 Duplicate executions from non-idempotent     Medium       High     `retryPolicy: do_not_retry` must be set and never overridden; the HTTP client must not
 retry under network ambiguity                                      auto-retry the mutation.
 Execution input data incomplete or           Medium       High     Implementation must fail safely with an actionable error and start no new execution when
 unavailable, leading to silent incorrect                           required context cannot be recovered (FR6, AC7).
 fresh run
 Behavior drift between `pipeline.retry`      Medium       Medium   Both action paths must share the same underlying rerun logic module; divergent
 and `execution.retry`                                             implementations are not acceptable.
 `docs:check` fails after implementation      Low          Low      Run `pnpm build && pnpm docs:generate` and commit generated docs before merging.
 because tool/resource counts changed

## Open Questions and Assumptions

- Open technical decisions:
    - Exact file paths for the execution toolset definition, pipeline toolset definition, shared rerun module, URL parser, and response extractors are not confirmed — requires reading `src/registry/toolsets/` and related source directories.
    - Whether the fresh-run fallback should surface a confirmation prompt separate from the initial high-write confirmation, or whether one confirmation covers both paths, is not specified.
    - How execution-inputs are retrieved (API shape, required parameters) depends on the execution-inputs resource/API implementation — to be confirmed during implementation.
    - The timeout and cancellation behavior for `wait=true` on the new execution follows the existing wait infrastructure; any deviations must be documented.

- Assumptions:
    - The existing URL parser already handles `execution` resource_type extraction; confirmation requires reading the parser source.
    - The existing server-side execution polling used by `pipeline.run` with `wait=true` can be reused for `execution.retry` with `wait=true` by substituting the new execution ID.
    - `pnpm test` runs vitest and covers the `tests/` directory, as confirmed by `package.json`.
    - The `src/registry/extractors.ts` file is the correct location for normalising response shapes, consistent with the coding standards referenced in the PR checklist.

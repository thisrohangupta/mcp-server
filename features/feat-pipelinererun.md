  from a previous execution, reusing its effective inputs—not retrying one failed
  stage.

  # Feature: Pipeline Execution Rerun

  **Status:** Proposed
  **Owner:** TBD
  **Last updated:** 2026-08-03

  ## Summary

  Allow users and agents to rerun a Harness pipeline directly from a previous
  execution ID or execution URL.

  The MCP server should treat the execution as the primary resource, reuse the
  previous execution context where supported, and return the new execution ID.
  Existing pipeline retry calls must remain backward compatible.

  ## Problem

  Pipeline retry is currently exposed as:

  ```text
  harness_execute(
    resource_type="pipeline",
    action="retry",
    params={execution_id: "<execution-id>"}
  )
```

  This is difficult to discover because:

  - The operation targets an execution but is registered on the pipeline resource.
  - execution_id must be supplied through params.
  - Harness execution URLs resolve to resource_type="execution", but that resource
    does not expose a retry action.

  - A request such as “rerun this execution” does not map naturally to the current
    contract.

  - When the native retry endpoint returns HTTP 405, the server may start a fresh
    pipeline run, but a fresh run is not necessarily equivalent to replaying the
    original execution inputs and source context.

  ## Users

  - Developers rerunning failed CI or CD executions.
  - Operators responding to transient deployment failures.
  - AI agents working from a Harness execution URL.
  - MCP clients that need one-call execution and status monitoring.

  ## Goals

  1. Allow a pipeline execution to be rerun using its execution ID.
  2. Allow an execution URL to be used without manually extracting identifiers.
  3. Preserve the previous execution’s effective inputs whenever possible.
  4. Return a clear distinction between a native retry and a fresh-run fallback.
  5. Support wait: true for the new execution.
  6. Preserve the existing pipeline.retry contract for backward compatibility.

  ## Non-goals

  - Retrying an individual stage or step.
  - Resuming paused executions.
  - Replaying dynamic pipeline executions.
  - Adding a new top-level MCP tool.
  - Guaranteeing that an external Git repository still contains the exact historical
    pipeline definition.

  - Changing the behavior of normal pipeline.run requests.

  ## Terminology

  - Native retry: The Harness retry endpoint creates another execution from an
    existing execution.

  - Fresh-run fallback: The server starts the associated pipeline again when native
    retry is unavailable.

  - Rerun: The user-facing intent encompassing a native retry or an explicitly
    reported fresh-run fallback.

  - Execution context: Runtime inputs, input sets, pipeline identifier, module, Git
    branch, and other information that influenced the original execution.

  ## Proposed experience

  ### Rerun by execution ID

  harness_execute(
    resource_type="execution",
    action="retry",
    resource_id="<execution-id>"
  )

  ### Rerun from a Harness URL

  harness_execute(
    url="<Harness execution URL>",
    action="retry"
  )

  ### Rerun and wait for completion

  harness_execute(
    resource_type="execution",
    action="retry",
    resource_id="<execution-id>",
    wait=true
  )

  The existing form remains supported:

  harness_execute(
    resource_type="pipeline",
    action="retry",
    params={execution_id: "<execution-id>"}
  )

  ## Functional requirements

  ### FR1: Execution retry action

  The execution resource must expose a declarative retry execute action.

  The action must:

  - Accept the execution ID through resource_id.
  - Map resource_id to execution_id.
  - Call the Harness pipeline retry endpoint.
  - Require project scope.
  - Require no request body.
  - Return the new execution ID.

  ### FR2: Execution URL support

  When an execution URL is provided, the server must:

  - Detect resource_type="execution".
  - Extract the execution ID.
  - Extract the pipeline ID when present.
  - Extract organization and project scope.
  - Use the extracted execution ID as the retry target.

  Explicit caller-provided values must continue to override URL-derived values.

  ### FR3: Confirmation and policy enforcement

  Rerunning a pipeline creates a new execution and must use:

  risk: high_write
  retryPolicy: do_not_retry

  The action must:

  - Respect HARNESS_READ_ONLY.
  - Use the existing confirmation or elicitation flow.
  - Emit the same audit information as other execute actions.
  - Never retry the mutation automatically at the HTTP-client level.

  ### FR4: Native retry

  The server should first attempt the native Harness retry endpoint:

  PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}

  A successful response must include, when available:

  {
    "execution_id": "<new-execution-id>",
    "rerun_mode": "native_retry",
    "source_execution_id": "<original-execution-id>"
  }

  ### FR5: Fresh-run fallback

  If native retry is unavailable and the server receives the currently supported 405
  response, it may fall back to a fresh pipeline run.

  Before performing the fallback, the server must:

  1. Resolve the pipeline ID from the URL, supplied parameters, or execution details.
  2. Retrieve the effective inputs used by the original execution when available.
  3. Preserve applicable Git branch and input-set context.
  4. Start the associated pipeline using the recovered context.
  5. Clearly report that a fresh run was used.

  Example response:

  {
    "execution_id": "<new-execution-id>",
    "rerun_mode": "fresh_run_fallback",
    "source_execution_id": "<original-execution-id>",
    "_note": "Native retry was unavailable. A fresh pipeline run was started using the
    recovered execution context."
  }

  The server must not claim that the fallback is an exact replay when the historical
  pipeline definition or inputs cannot be recovered.

  ### FR6: Unsafe fallback prevention

  If native retry fails and the server cannot resolve the pipeline ID or required
  inputs, it must return an actionable error instead of starting an incomplete run.

  Example:

  Native retry is unavailable, and the original pipeline inputs could not be
  recovered. No new execution was started.

  ### FR7: Wait support

  wait: true must work for execution.retry.

  After the rerun starts, the server must:

  - Extract the new execution ID.
  - Poll that execution rather than the source execution.
  - Return terminal status, timeout information, and elapsed time using the existing
    wait response fields.

  - Respect cancellation and configured timeout limits.

  ### FR8: Backward compatibility

  Existing calls using:

  resource_type="pipeline", action="retry"

  must continue to work.

  Both the legacy and canonical forms should use the same underlying rerun behavior so
  their safety, fallback, audit, and response semantics do not drift.

  ### FR9: Discovery metadata

  harness_describe(resource_type="execution") must advertise:

  - The retry action.
  - Its required execution identifier.
  - Its risk classification.
  - Its lack of request body.
  - Support for wait: true.
  - Its possible fresh-run fallback.

  The general harness_execute description should mention that executions can be
  retried directly.

  ## Acceptance criteria

  ### AC1: Retry by resource ID

  Given a valid failed execution ID, when the caller executes:

  harness_execute(
    resource_type="execution",
    action="retry",
    resource_id="<execution-id>",
    confirm=true
  )

  then the server targets the retry endpoint with that execution ID and returns the
  new execution ID.

  ### AC2: Retry by URL

  Given a valid Harness execution URL, when the caller provides the URL and
  action="retry", then the server extracts the execution and scope identifiers and
  starts the retry without requiring duplicate parameters.

  ### AC3: Confirmation

  Given a client without auto-approval, the rerun must not start until the high-write
  operation is confirmed.

  ### AC4: Read-only mode

  Given HARNESS_READ_ONLY=true, the rerun is rejected before any mutation request is
  sent.

  ### AC5: Native retry response

  When native retry succeeds, the response reports:

  rerun_mode="native_retry"

  and identifies both the source and new executions.

  ### AC6: Successful fallback

  When native retry returns 405 and the original pipeline context is recoverable, the
  server starts a fresh run and reports:

  rerun_mode="fresh_run_fallback"

  ### AC7: Incomplete fallback

  When native retry returns 405 and required context cannot be recovered, the server
  returns an actionable error and sends no fresh-run request.

  ### AC8: Wait behavior

  When wait=true, the server polls the new execution until it reaches a terminal state
  or the timeout expires.

  ### AC9: Legacy compatibility

  Existing pipeline.retry calls continue to pass their current tests and return the
  normalized rerun response.

  ### AC10: Discovery

  harness_describe shows retry as an available action for the execution resource with
  accurate identifier and policy metadata.

  ## Error behavior

  The server should return actionable errors for:

  - Missing execution ID.
  - Execution not found.
  - Execution not eligible for native retry.
  - Missing organization or project scope.
  - Insufficient Execute permission.
  - Native retry unavailable with no recoverable pipeline ID.
  - Required original inputs unavailable.
  - Fresh-run fallback rejected by confirmation or policy.
  - Wait timeout or cancellation.

  Errors must state whether a new execution was started.

  ## Testing requirements

  Add focused tests covering:

  1. execution.retry registration and discovery.
  2. resource_id to execution_id mapping.
  3. Execution URL parsing and dispatch.
  4. Explicit values overriding URL-derived scope.
  5. High-write confirmation.
  6. Read-only rejection before network I/O.
  7. Native retry success.
  8. New execution ID extraction from supported response shapes.
  9. HTTP 405 fresh-run fallback.
  10. Recovery of pipeline ID from execution details.
  11. Recovery and forwarding of original execution inputs.
  12. Fallback refusal when required context is unavailable.
  13. wait=true polling the new execution.
  14. Timeout and cancellation behavior.
  15. Legacy pipeline.retry compatibility.
  16. Audit records for successful, blocked, and fallback paths.

  ## Documentation

  Update:

  - Tool and action descriptions.
  - Generated resource documentation.
  - Pipeline and execution testing plans.
  - Example calls for retrying by ID and URL.
  - Response documentation explaining rerun_mode.

  Run documentation generation only after a fresh build:

  pnpm build && pnpm docs:generate

  ## Success metrics

  - A pipeline execution can be rerun from its Harness URL in one MCP call.
  - No manual execution-ID extraction is required.
  - Every successful rerun response contains the new execution ID.
  - Fresh-run fallbacks are clearly identified.
  - No fallback starts a pipeline when required execution context is missing.
  - Existing pipeline retry clients remain compatible.

  ## Risks and tradeoffs

  ### Retry versus replay

  A native retry and a fresh pipeline run may have different behavior. Pipeline
  definitions, templates, expressions, secrets, connectors, and Git content can change
  after the source execution.

  The response must therefore distinguish between native retry and fresh-run fallback.

  ### Duplicate execution risk

  Rerun is a non-idempotent operation. Retrying the HTTP request after an ambiguous
  network failure could create duplicate executions. The operation must retain
  retryPolicy="do_not_retry".

  ### Input recovery

  Execution input data may be incomplete or unavailable. The implementation must fail
  safely rather than silently running with defaults.

  ### Contract duplication

  Temporarily supporting retry on both pipeline and execution resources introduces two
  discovery paths. Shared implementation logic is required to prevent behavior drift.

  ## Dependencies

  - Pipeline retry API.
  - Pipeline execution-details API.
  - Execution-inputs resource/API.
  - Existing URL parser.
  - Existing confirmation and audit infrastructure.
  - Existing server-side execution polling.

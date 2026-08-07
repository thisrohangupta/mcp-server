import type { Registry } from "../registry/index.js";
import type { HarnessClient } from "../client/harness-client.js";
import { HarnessApiError } from "../utils/errors.js";
import { asRecord, asString } from "../utils/type-guards.js";
import { createLogger } from "../utils/logger.js";

const log = createLogger("rerun");

export type RerunMode = "native_retry" | "fresh_run_fallback";

export interface RerunResult {
  execution_id: string | undefined;
  source_execution_id: string;
  rerun_mode: RerunMode;
  _note?: string;
  [key: string]: unknown;
}

export interface RerunOptions {
  /** The execution ID to retry. Maps to `planExecutionId` in the Harness retry API. */
  executionId: string;
  /** Optional pipeline ID. When absent, resolved from execution details on 405. */
  pipelineId?: string;
  /** Shared input map — merged from caller input and URL-parsed defaults. */
  input: Record<string, unknown>;
  /** Audit context label for fresh-run fallback dispatch. */
  auditCtx: { tool: "harness_execute"; confirmation: string; resource_id?: string; action: string };
}

/**
 * Shared rerun logic: native retry → 405 fallback → unsafe refusal.
 *
 * Attempts `PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}` first.
 * On HTTP 405, falls back to a fresh pipeline run using execution details to
 * recover the pipeline ID. Refuses to start a fresh run when the pipeline ID
 * is unrecoverable, returning an actionable error instead.
 *
 * Both `execution.retry` and `pipeline.retry` must delegate here to keep
 * safety, fallback, audit, and response semantics identical across paths.
 */
export async function executeRerun(
  registry: Registry,
  client: HarnessClient,
  options: RerunOptions,
): Promise<RerunResult> {
  const { executionId, pipelineId: callerPipelineId, input, auditCtx } = options;

  log.debug("executeRerun: attempting native retry", { executionId });

  try {
    const raw = await registry.dispatchExecute(
      client,
      "pipeline",
      "retry",
      { ...input, execution_id: executionId },
      auditCtx,
    );
    const rec = asRecord(raw) ?? {};
    const newExecutionId =
      asString(rec.planExecutionId) ??
      asString(asRecord(rec.planExecution)?.uuid) ??
      asString(rec.executionId);
    log.info("executeRerun: native retry succeeded", { newExecutionId });
    return {
      ...rec,
      execution_id: newExecutionId,
      source_execution_id: executionId,
      rerun_mode: "native_retry",
    };
  } catch (err) {
    if (!(err instanceof HarnessApiError) || err.statusCode !== 405) {
      throw err;
    }
    log.info("executeRerun: native retry returned 405, attempting fresh-run fallback", { executionId });
  }

  // 405 path: resolve pipeline ID from caller params or execution details.
  let resolvedPipelineId = callerPipelineId ?? asString(input.pipeline_id);

  if (!resolvedPipelineId) {
    try {
      const execDetails = asRecord(
        await registry.dispatch(client, "execution", "get", { ...input, execution_id: executionId }),
      );
      const summary = asRecord(execDetails?.pipelineExecutionSummary);
      resolvedPipelineId = asString(summary?.pipelineIdentifier);
      log.debug("executeRerun: resolved pipelineId from execution details", { resolvedPipelineId });
    } catch (fetchErr) {
      log.warn("executeRerun: could not fetch execution details for pipelineId", {
        executionId,
        error: String(fetchErr),
      });
    }
  }

  if (!resolvedPipelineId) {
    throw new Error(
      `Native retry is unavailable for execution "${executionId}" (405), and the original pipeline ID ` +
      `could not be recovered. No new execution was started. ` +
      `Provide pipeline_id explicitly to start a fresh run.`,
    );
  }

  log.info("executeRerun: starting fresh-run fallback", { resolvedPipelineId, executionId });

  const raw = await registry.dispatchExecute(
    client,
    "pipeline",
    "run",
    { ...input, pipeline_id: resolvedPipelineId },
    { ...auditCtx, action: "run (retry fallback)" },
  );

  const rec = asRecord(raw) ?? {};
  const newExecutionId =
    asString(rec.planExecutionId) ??
    asString(asRecord(rec.planExecution)?.uuid) ??
    asString(rec.executionId);

  return {
    ...rec,
    execution_id: newExecutionId,
    source_execution_id: executionId,
    rerun_mode: "fresh_run_fallback",
    _note:
      "Native retry was unavailable (405). A fresh pipeline run was started using the recovered execution context.",
  };
}

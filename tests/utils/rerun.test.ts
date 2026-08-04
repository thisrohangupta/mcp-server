import { describe, it, expect, vi } from "vitest";
import type { Registry } from "../../src/registry/index.js";
import type { HarnessClient } from "../../src/client/harness-client.js";
import { executeRerun } from "../../src/utils/rerun.js";
import { HarnessApiError } from "../../src/utils/errors.js";

const fakeClient = {} as unknown as HarnessClient;

const auditCtx = {
  tool: "harness_execute" as const,
  confirmation: "auto_approved",
  resource_id: "exec-abc",
  action: "retry",
};

function make405Error(): HarnessApiError {
  const err = new HarnessApiError("Method Not Allowed", 405);
  return err;
}

/**
 * Build a minimal registry mock.
 * - dispatchExecute: used by executeRerun for pipeline.retry and pipeline.run
 * - dispatch: used by executeRerun to GET execution details when pipelineId is missing
 */
function makeRegistry(opts: {
  retryResult?: unknown;
  retryError?: Error;
  runResult?: unknown;
  execDetails?: unknown;
  execDetailsError?: Error;
}): Registry {
  const dispatchExecute = vi.fn(async (_client: unknown, _type: unknown, action: unknown, _input: unknown, _audit: unknown) => {
    if (action === "retry") {
      if (opts.retryError) throw opts.retryError;
      return opts.retryResult ?? { planExecutionId: "new-exec-123" };
    }
    if (action === "run (retry fallback)" || action === "run") {
      return opts.runResult ?? { planExecutionId: "fallback-exec-456" };
    }
    throw new Error(`Unexpected action: ${String(action)}`);
  });

  const dispatch = vi.fn(async (_client: unknown, _type: unknown, _op: unknown, _input: unknown) => {
    if (opts.execDetailsError) throw opts.execDetailsError;
    return opts.execDetails ?? {
      pipelineExecutionSummary: { pipelineIdentifier: "my-pipeline" },
    };
  });

  return { dispatchExecute, dispatch, orgId: "default-org", projectId: "default-project" } as unknown as Registry;
}

describe("executeRerun", () => {
  it("returns native_retry result when retry endpoint succeeds", async () => {
    const registry = makeRegistry({
      retryResult: { planExecutionId: "new-exec-111", status: "Running" },
    });

    const result = await executeRerun(registry, fakeClient, {
      executionId: "exec-abc",
      input: {},
      auditCtx,
    });

    expect(result.rerun_mode).toBe("native_retry");
    expect(result.source_execution_id).toBe("exec-abc");
    expect(result.execution_id).toBe("new-exec-111");
    expect(result._note).toBeUndefined();
  });

  it("falls back to fresh run when retry returns 405 (caller-supplied pipelineId)", async () => {
    const registry = makeRegistry({
      retryError: make405Error(),
      runResult: { planExecutionId: "fallback-exec-789" },
    });

    const result = await executeRerun(registry, fakeClient, {
      executionId: "exec-abc",
      pipelineId: "my-pipeline",
      input: {},
      auditCtx,
    });

    expect(result.rerun_mode).toBe("fresh_run_fallback");
    expect(result.source_execution_id).toBe("exec-abc");
    expect(result.execution_id).toBe("fallback-exec-789");
    expect(result._note).toContain("405");
  });

  it("resolves pipelineId from execution details when not supplied by caller", async () => {
    const registry = makeRegistry({
      retryError: make405Error(),
      execDetails: {
        pipelineExecutionSummary: { pipelineIdentifier: "resolved-pipeline" },
      },
      runResult: { planExecutionId: "resolved-fallback-exec" },
    });

    const result = await executeRerun(registry, fakeClient, {
      executionId: "exec-abc",
      input: {},
      auditCtx,
    });

    expect(result.rerun_mode).toBe("fresh_run_fallback");
    expect(result.execution_id).toBe("resolved-fallback-exec");
  });

  it("throws when 405 and pipelineId cannot be resolved", async () => {
    const registry = makeRegistry({
      retryError: make405Error(),
      execDetailsError: new Error("execution not found"),
    });

    await expect(
      executeRerun(registry, fakeClient, {
        executionId: "exec-abc",
        input: {},
        auditCtx,
      }),
    ).rejects.toThrow(/pipeline ID.*could not be recovered|pipeline_id explicitly/i);
  });

  it("propagates non-405 errors without fallback", async () => {
    const registry = makeRegistry({
      retryError: new HarnessApiError("Internal Server Error", 500),
    });

    await expect(
      executeRerun(registry, fakeClient, {
        executionId: "exec-abc",
        input: {},
        auditCtx,
      }),
    ).rejects.toMatchObject({ statusCode: 500 });
  });

  it("picks up pipelineId from input.pipeline_id when caller omits pipelineId arg", async () => {
    const registry = makeRegistry({
      retryError: make405Error(),
      runResult: { planExecutionId: "input-pipeline-exec" },
    });

    const result = await executeRerun(registry, fakeClient, {
      executionId: "exec-abc",
      input: { pipeline_id: "input-pipeline" },
      auditCtx,
    });

    expect(result.rerun_mode).toBe("fresh_run_fallback");
    expect(result.execution_id).toBe("input-pipeline-exec");
  });
});

/**
 * Security regression test: the registry must refuse to dispatch any request
 * whose fully-resolved API path contains injection characters or path-traversal
 * segments.
 *
 * Tool inputs (resource IDs, org/project, version labels) are interpolated into
 * Harness API paths. Most pathBuilders / placeholder substitutions encode their
 * values, but not all do — so a centralized guard in Registry.dispatch validates
 * the resolved path before it ever reaches the HTTP client. This prevents a
 * malicious identifier (e.g. supplied via prompt injection) from traversing to a
 * different endpoint or smuggling query params such as `?accountIdentifier=...`
 * to reach another account.
 */
import { describe, it, expect, vi } from "vitest";
import { Registry } from "../../src/registry/index.js";
import type { Config } from "../../src/config.js";
import type { HarnessClient } from "../../src/client/harness-client.js";

function makeConfig(overrides: Partial<Config> = {}): Config {
  return {
    HARNESS_API_KEY: "pat.acct.tok.secret",
    HARNESS_ACCOUNT_ID: "acct",
    HARNESS_BASE_URL: "https://app.harness.io",
    HARNESS_ORG: "default",
    HARNESS_PROJECT: "proj",
    HARNESS_API_TIMEOUT_MS: 30000,
    HARNESS_MAX_RETRIES: 3,
    LOG_LEVEL: "info",
    HARNESS_MAX_BODY_SIZE_MB: 10,
    HARNESS_RATE_LIMIT_RPS: 10,
    HARNESS_READ_ONLY: false,
    HARNESS_SKIP_ELICITATION: false,
    HARNESS_ALLOW_HTTP: false,
    HARNESS_FME_BASE_URL: "https://api.split.io",
    ...overrides,
  };
}

function mockClient(): HarnessClient {
  return {
    request: vi.fn().mockResolvedValue({ status: "SUCCESS", data: {} }),
  } as unknown as HarnessClient;
}

describe("path injection guard", () => {
  it("rejects query-parameter smuggling via an interpolated identifier", async () => {
    const registry = new Registry(makeConfig());
    const client = mockClient();

    await expect(
      registry.dispatch(client, "registry", "get", {
        org_id: "PROD",
        project_id: "Harness",
        // `?accountIdentifier=victim` would be re-parsed by the HTTP client and
        // could override the account scope.
        registry_id: "x?accountIdentifier=victim",
      }),
    ).rejects.toThrow(/illegal characters/);

    expect((client.request as ReturnType<typeof vi.fn>)).not.toHaveBeenCalled();
  });

  it("rejects path-traversal segments in an interpolated identifier", async () => {
    const registry = new Registry(makeConfig());
    const client = mockClient();

    await expect(
      registry.dispatch(client, "registry", "get", {
        org_id: "PROD",
        project_id: "Harness",
        registry_id: "../../secrets",
      }),
    ).rejects.toThrow(/path-traversal/);

    expect((client.request as ReturnType<typeof vi.fn>)).not.toHaveBeenCalled();
  });

  it("rejects whitespace / control characters in an interpolated identifier", async () => {
    const registry = new Registry(makeConfig());
    const client = mockClient();

    await expect(
      registry.dispatch(client, "registry", "get", {
        org_id: "PROD",
        project_id: "Harness",
        registry_id: "foo bar",
      }),
    ).rejects.toThrow(/illegal characters/);

    expect((client.request as ReturnType<typeof vi.fn>)).not.toHaveBeenCalled();
  });

  it("allows a normal identifier through unchanged", async () => {
    const registry = new Registry(makeConfig());
    const client = mockClient();

    await registry.dispatch(client, "registry", "get", {
      org_id: "PROD",
      project_id: "Harness",
      registry_id: "ai-platform",
    });

    expect((client.request as ReturnType<typeof vi.fn>)).toHaveBeenCalledOnce();
    const callArgs = (client.request as ReturnType<typeof vi.fn>).mock.calls[0][0];
    expect(callArgs.path).toBe("/har/api/v1/registry/acct/PROD/Harness/ai-platform/+");
  });
});

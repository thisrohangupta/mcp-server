# Spec 007: AI SRE Investigation in Harness MCP

**Status:** Draft (open decision on API shape, see Options)
**Owners:** Rohan Gupta, Tina Huang
**Reviewer:** Sunil
**Date:** 2026-10-06

---

## Summary

Harness MCP lets an agent call the project's preconfigured AI SRE investigator.
The agent starts an investigation and later retrieves its status and findings.
It does not rebuild the investigation by chaining raw incident, alert, deploy
and timeline calls.

## Problem

- MCP exposes raw entities and operations; the calling agent composes them.
- An agent that rebuilds an investigation this way bypasses the project's shared
  AI SRE configuration: prompts, tools, connectors, permissions, and background
  run and write-back behavior.
- Results vary between runs and agents, and nothing is written back to the incident.
- Harness AI has no clean way to start an investigation and return for the result.

## Goals

1. Harness AI and other MCP clients can start a project's configured investigator.
2. Clients can poll status and read findings without owning investigation logic.
3. Investigator configuration and write-back stay in AI SRE.
4. MCP tool count does not grow in a way that hurts IDE tool limits.

## Non-goals

- Re-implementing investigation logic in MCP.
- Replacing `harness_diagnose`.
- Configuring the investigator through MCP (stays in AI SRE).

## Key distinctions

| | `harness_diagnose` | AI SRE investigation |
|---|---|---|
| Purpose | Debug Harness itself: pipelines/executions, connectors, delegates, GitOps apps | Investigate production issues (day 2 operations) using Harness data |
| Execution | Synchronous, read-only, server-side fan-out | Durable, asynchronous background run with write-back |
| Config | None | Per-project prompts, tools, connectors, permissions |

They do not conflict. Tool descriptions should state this boundary so agents
pick the right one.

## Proposed behavior

- **Start:** caller supplies a target (incident in v1) and optional context.
  Receives an investigation reference and initial status.
- **Status:** read via the incident activity timeline, or directly if runs become entities.
- **Findings:** summary, hypotheses, evidence, recommended actions. Also written
  back to the incident by the investigator.
- **Identity and access:** the investigation runs under the investigator's
  configured service account (cross-project access, data fetched from UDP).
  The caller's own permissions only gate who may start and read.
- **Data sources:** MCP by default, plus data configured in AI SRE (for example
  incoming alerts). Over time this moves into MCP, and third-party ingestion
  writes to UDP.
- The investigator is moving to a worker agent that has Harness MCP access by default.

## Options for API shape

| Option | Description | Pros | Cons |
|---|---|---|---|
| A | `harness_execute(resource_type="incident", action="investigate")`; findings via `harness_get` | Keeps the fixed 11 tools; fits "execute and get over conforming REST" | Interface does not signal a durable capability |
| B | New `investigation` resource type (runs are entities with IDs), same 11 tools | Fits the planned ad hoc investigation object; status/findings via `harness_get`/`harness_list` | Needs backend to expose runs as entities |
| C | New 12th tool, e.g. `harness_investigate` | Most discoverable; could be a cross-cutting entrypoint (IaCM, cost, deployments, builds) | Breaks the fixed-tool-count rule; the earlier 150+ tool release was rejected by IDEs and performed badly |

**Recommendation:** ship A first, move to B once runs have stable IDs. Choose C
only if Sunil approves an explicit exception.

Implementation note for A/B: add the operation declaratively in
`src/registry/toolsets/incidents.ts` (the `incident` resource already has
list/get/create/update/close). No new tool handlers.

## Success metrics

- Share of MCP-started investigations that complete with write-back.
- Time from start to first findings.
- Tool count stays at 11 unless an exception is approved.
- Agents use `harness_diagnose` for Harness failures and the investigator for production incidents.

## Open questions

1. Is C an acceptable exception, or is the 11-tool rule absolute?
2. Do runs get stable IDs soon enough to justify B?
3. v1 targets: incidents only, or also services, deployments, cost anomalies?
4. Status and findings delivery: polling, progress notifications, or both?
5. Behavior when a project has no investigator configured?
6. Which `HARNESS_AUTO_APPROVE_RISK` level does `investigate` fall under?

## Next steps

- Get Sunil's decision on A, B or C.
- Confirm the worker-agent API and service account model with the AI SRE team.
- Confirm the backend endpoint for starting an investigation and reading findings.

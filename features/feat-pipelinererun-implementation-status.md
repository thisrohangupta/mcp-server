# Pipeline Execution Rerun — Implementation Status
<!-- github-implementation-status:v1 plan=features/feat-pipelinererun-plan.md plan-blob=32e7732a9cfad31478d7c4fe11c1964d35eb573b -->

## Tasks

 #      Task                                                                                                        Status    Attempts   Commit SHA   Last Run ID                    Notes
━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━  ━━━━━━━━━  ━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
 1      Identify and extract shared rerun logic                                                                     done      1          SELF         feat-pipelinererun-run-1       Created src/utils/rerun.ts: native retry → 405 fallback → unsafe refusal; RerunResult type with rerun_mode, source_execution_id, execution_id
 2      Register `retry` action on the `execution` toolset definition                                               pending   0          —            —
 3      Update `pipeline` toolset `retry` action to delegate to shared rerun logic                                  pending   0          —            —
 4      Implement native retry attempt (`PUT /pipeline/api/pipeline/execute/retry/{planExecutionId}`)               pending   0          —            —
 5      Implement fresh-run fallback: resolve pipeline ID and effective inputs from execution details               pending   0          —            —
 6      Implement unsafe-fallback prevention                                                                        pending   0          —            —
 7      Normalise response shape: add `rerun_mode`, `source_execution_id`, `execution_id` fields                    pending   0          —            —
 8      Wire `wait=true` to poll the new execution ID returned by the rerun                                         pending   0          —            —
 9      Update URL parser to confirm `execution` resource_type extraction                                           pending   0          —            —
 10     Update `harness_describe` metadata for `execution` resource                                                 pending   0          —            —
 11     Apply `risk: high_write`, `retryPolicy: do_not_retry`, `HARNESS_READ_ONLY` check                            pending   0          —            —
 12     Write unit and integration tests                                                                            pending   0          —            —
 13     Update tool/action descriptions and generated documentation                                                 pending   0          —            —

# MANIFEST — AGENTRUNNER-PROD-PATCH mutation evidence

**Slice**: `AGENTRUNNER-PROD-PATCH` (`[PILOT]`, tier M). **Frozen**: 2026-10-21.

## Must read

`evidence.txt` is the **core evidence**: the complete stdout transcript of the mutation run
(baseline → mutation M1 → restore → mutation M2 → restore → summary), including the exact
`go test` command lines, exit codes, failing test output, and the SHA256 readings.

## 归档条目清单

| File | Bytes | Role |
|---|---|---|
| `evidence.txt` | 5290 | core evidence: full transcript (must read) |
| `run-mutations.js` | 5640 | the harness that produced it (re-runnable) |
| `dispatcher.pristine.go.txt` | 12104 | the exact pre-mutation file the harness restored from |

## Deliberate exclusions (regenerable, not archived)

The live run directory also held two large artifacts that are **byte-reproducible** and were
therefore not committed (`head-tree.tar` 26,531,840 bytes; `wt/` 471 files / 4,410,758 bytes):

1. `head-tree.tar` = `git archive --format=tar -o head-tree.tar 6898f91`
2. `wt/` = that tar extracted, restricted to `go.mod go.sum internal cmd tools examples schemas testdata harnesses`,
   then overwritten with the two files under test.

Re-running `node run-mutations.js` from a directory containing this package recreates both and
reproduces `evidence.txt`. The exclusions are recorded here so the absence is never read as loss.

## What the evidence proves

- **Baseline** (`go test -count=1 -run "AgentRunnerReplayDispatcher" ./internal/adapters/` on the unmutated copy): **exit 0**.
- **M1** (`context.WithTimeout(context.WithoutCancel(ctx), …)` → `context.WithTimeout(ctx, …)`): the caller-cancellation
  shield test fails with `shield failed: … Err()=context canceled`.
- **M2** (append `ok = false` after the capability type assertion): **8 reclaim tests fail** (7 top-level + 1 subtest).
- **Both restores** return the sha256 to the pre-mutation value, so the red runs cannot be attributed to a dirty copy.

## Note on the first M2 attempt

A first M2 variant (`if !ok || reclaimer == nil {` → `if true {`) did **not compile** (`declared and not used: ok`).
A compile failure is **not** falsification evidence; the evidence above comes from the compiling variant.
That iteration is recorded as a tooling adjustment and is not counted as a finding or a rework.

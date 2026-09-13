# T027 Replay Convergence and Publication Durability Prerequisite Validation Report

[中文](t027-replay-convergence.md)

Date: 2026-09-14. Verdict: `PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`.

## Scope

This report covers only the completed offline replay-convergence/publication-durability prerequisite slice. It does not claim a production AgentRunner implementation. Effect-mapping v1 is frozen: `local-process` and `workspace-write` both map to `local-discardable`; requests bind the mapping version/hash; missing, unknown, or tampered bindings, unknown operations, and grant-class under-coverage fail closed before replay identity consumption or dispatch.

`AgentRunnerReplayStore` now covers:

- fail-closed orphan-completion and on-disk-corruption handling;
- bounded reread convergence after no-replace collisions;
- a single first-dispatch for concurrent stores sharing one root, with other calls blocked;
- pre-write blocking when Windows publication durability is `unproven`, without creating a request record.

## Final Verifier Results

### Executed on Windows

The following commands were executed on Windows:

| Command | Result |
|---|---|
| `node tools/contracts/contracts.test.js` | Node contract 4/4 passed |
| `go test -count=1 ./internal/adapters -run '^TestAgentRunnerReplayStore'` | replay 32/32 items passed |
| `go test -count=1 ./internal/adapters` | adapters passed |
| `go build ./...` | passed |
| `go vet ./...` | passed |
| `go test -count=1 ./...` | all passed |

### Unix Targets: Cross-Compilation Only

| Command | Result |
|---|---|
| `$env:GOOS='linux'; $env:GOARCH='amd64'; go build ./...` | Linux/amd64 cross-compile passed |
| `$env:GOOS='darwin'; $env:GOARCH='amd64'; go build ./...` | Darwin/amd64 cross-compile passed |

The Linux/Darwin commands prove only that the target builds complete. They did not execute replay, `fsync`, crash injection, or tests on those systems; cross-compilation must not be recorded as a Unix runtime/platform PASS.

## Publication Durability Boundary

Windows may provide no-replace visibility, but it cannot currently prove persistence of a new directory entry. Therefore any write that could return `first-dispatch` is rejected before writing, and seeing a file does not justify a durable-dispatch claim.

The Unix `proven` path requires an already-existing stable absolute replay root and successful parent-directory `fsync` after creating or confirming `requests/` and `completions/`. The final verifier ran on Windows, so this run contains no native Unix crash-injection evidence; those Unix conditions are contract boundaries, not native evidence from this run.

A terminal receipt, including `uncertain`, proves only terminal-receipt presence and duplicate-dispatch blocking. It is not task `PASS` and does not pass AT-23. The replay store is not yet wired into `AgentRunnerPort` or the Engine.

## Explicitly Not Executed

- No real model CLI was invoked.
- No network was accessed.
- The SecretStore was not read.
- No commit, push, or publish action was performed.
- T027 remains `BLOCKED / NOT IMPLEMENTED`; AT-23 does not pass.

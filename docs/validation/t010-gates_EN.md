# T010 Gate Runner Validation Report

[简体中文](t010-gates.md)

Date: 2026-09-07. Verdict: `T010 COMPLETE`. This task completes hook policy preflight, runner orchestration, structured results, evidence persistence, and Chain Engine wiring, with AT-06 acceptance.

## Twelve Slices

1. Freeze hook, process-runner, resource/network policy, execution, result, and capability models.
2. Validate hook kinds, failure policies, timeout, positive resource limits, IDs, and definition hashes.
3. Preserve executable and args as separate argv without an implicit shell; reject absolute and backslash executable forms.
4. Resolve cwd inside the authorized workspace and reject both `..` and symlink escapes.
5. Select environment values only through an explicit allowlist; missing, duplicate, or invalid names fail closed.
6. Preflight process, timeout, output, memory, process-count, CPU, and network capabilities; unenforceable policy never degrades.
7. Reuse the T007 guard for process-tree ownership, covering normal exit, start failure, timeout stop, and termination uncertainty.
8. Capture stdout/stderr under one shared bound; overflow stops the process tree and records resource-limited/truncated.
9. Secret-scan stdout/stderr and declared artifacts before storing bytes in snapshot CAS and binding digests; scan failure or missing artifacts blocks.
10. Separate executionOutcome, assessment, failureKind, and policyDisposition; warn remains failed, while each retry is persisted and bounded by maxAttempts.
11. Hash result JCS under `proofrail:hook-result:1` and return only after JSONL fsync; reject tampering, corrupt logs, and duplicate execution keys.
12. Bind the Chain HookPort to step/hook kind and definition hash, attaching durable resultHash evidence for both success and failure.

## Validation

- Windows: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed.
- Contracts: `npm test --prefix tools/contracts` passed; all 29 Draft 2020-12 schemas, positive/negative catalogs, and JCS vectors remained stable.
- AT-06: argv injection, cwd/symlink escape, environment-secret exclusion, zero starts when capabilities are unavailable, timeout, unlimited output, secret-scan failure, missing artifacts, nonzero exit, retry, result tampering/duplication, and persistence failure passed.
- Ubuntu: Go 1.27.1 on Ubuntu 24.04/ext4 passed offline vendored `go build ./...`, `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./internal/guard ./internal/gates`.
- Windows did not run `-race` because the host has `CGO_ENABLED=0` and no C compiler; native Ubuntu race detection passed.
- Temporary archives, vendor content, and the remote validation directory were removed.

## Boundary

The built-in `GuardExecutor` genuinely provides separate argv, process-tree control, timeout, and a shared output bound, but does not claim network isolation or memory, process-count, or CPU quotas. Because the hook Schema requires those policies, direct use of this limited executor is blocked during preflight; production wiring must supply an Executor that declares and actually enforces the full capability set, with no degraded execution. Container/remote/shell/PTY runners, the unified cross-task retry budget, T011 review/promotion, and T022 effect policy are outside T010.
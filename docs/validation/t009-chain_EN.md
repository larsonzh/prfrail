# T009 Chain Engine Validation Report

[简体中文](t009-chain.md)

Date: 2026-09-07. Verdict: `T009 COMPLETE`. This task completes Chain Engine ordered scheduling, event projections, pause/cancel, and crash replay, with fake agents/runners providing the chain-layer contribution to AT-03.

## Fourteen Slices

1. Freeze chain/task/step, change-mode, workspace, snapshot, and consumer-owned port models.
2. Reuse the evidence Schema transition table for chain/task/step without creating an alternate state protocol.
3. Reject unlisted transitions, self-transitions, and illegal initial states before persistence.
4. Persist state events in an append-only JSONL store and reject corrupt or truncated records.
5. Rebuild sequence, event hash, states, and the accepted-parent projection from the complete event chain.
6. Persist `NONE→CREATED` first and enter `BASELINED` only after baseline-0 evidence is durable.
7. Schedule three tasks in definition order and use only the previous accepted snapshot as the next parent.
8. Cover code/build/verify/noop; noop records only its reason and never invokes a runner.
9. Route managed-change-set code steps through an explicit port.
10. Route isolated-workspace code steps through a separate explicit port.
11. Apply pause at an atomic boundary, stop new steps, and resume an ordinary pause from its projection to completion.
12. On cancel, invoke the managed-process stopper and archive evidence before writing irreversible `CANCELLED`.
13. Run the reconciler before recovery validation and replay; a running step becomes recovery-uncertain `PAUSED`, with redispatch and direct resume prohibited.
14. Validate AT-03 natively on Windows 11/NTFS and Ubuntu 24.04/ext4.

## Validation

- Windows: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed.
- Chain: three-task ordering, four kinds, both change modes, accepted-parent advancement, T2 failure blocking T3, event-before-projection, zero writes on illegal transitions, and corrupt-tail rejection passed.
- Control: pause/resume, cancel/stopper, zero runner calls for noop, and irreversible cancellation passed.
- Recovery: accepted-parent reconstruction, reconciler invocation, no redispatch of a running step, durable recovery-uncertain pause, and rejection of direct resume passed.
- Ubuntu: Go 1.27.1 on Ubuntu 24.04/ext4 passed offline vendored `go build ./...`, `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./internal/chain`.
- Windows did not run `-race` because the host has `CGO_ENABLED=0` and no C compiler; native Ubuntu race detection passed.
- Temporary archives, vendor content, and the remote validation directory were removed.

## Boundary

T009 uses fake agents/runners and does not implement the real T010 gate runner, T011 review/promotion, or T013/T014 agents and operator handoff. `AcceptancePort` expresses only the invariant that later tasks consume an accepted parent snapshot; candidate, failed, and uncertain output cannot become a downstream parent. Uncertain recovery remains fail closed for a later repair workflow and cannot guess success or redispatch work.
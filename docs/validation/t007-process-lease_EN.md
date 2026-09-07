# T007 Process and Lease Guard Validation Report

[简体中文](t007-process-lease.md)

Date: 2026-09-07. Verdict: `T007 COMPLETE`. This task completes the guard-layer contribution to AT-03/AT-06 for process-tree shutdown and single-writer leases. T009 and T010 retain ownership of chain scheduling and hook policy execution.

## Ten Slices

1. Freeze managed-process, identity, termination-evidence, and fencing-token models.
2. Start Linux processes in a dedicated process group and terminate the whole PGID tree.
3. Start Windows processes with `CREATE_SUSPENDED`, bind them to a `KILL_ON_JOB_CLOSE` Job Object, then resume them, eliminating the pre-assignment escape window.
4. Bind PID identity to `/proc/<pid>/stat` starttime on Linux and process creation time on Windows to resist PID reuse.
5. Attempt the platform stop boundary, force termination when required, and verify both the tree and original identity are gone.
6. Produce fresh structured termination evidence using JCS canonical data and domain-separated SHA-256; uncertainty returns `ErrTerminationUncertain`.
7. Publish file leases under mutual exclusion; initial generation is 1 and concurrent initial claims have exactly one winner.
8. Renewal retains claimId and increments generation; stale `(claimId,generation)` tokens cannot validate or release the current lease.
9. Takeover requires a different claimId, higher generation, fresh stopped evidence matching the old writer identity, at least one authorization reference, and closure validation by an injected independent `TakeoverVerifier`. No verifier means fail closed; expiry, missing PID, or restart alone never authorizes takeover.
10. Verify native parent/descendant trees, concurrent leases, and takeover counterexamples on Windows 11/NTFS and Ubuntu 24.04/ext4, with bilingual documentation.

## Validation

- Windows: `go test -count=20 ./internal/guard` and `go vet ./internal/guard` passed. Every run launched a real test parent and descendant; both identities were offline after Job Object termination.
- Ubuntu: The cross-compiled test binary passed all guard tests natively on Ubuntu 24.04/ext4, including process-group tree termination.
- T006 combined regression: `go test ./internal/snapshot` and `go vet ./internal/snapshot` passed.
- The local Go toolchain has cgo disabled, so `go test -race` was unavailable. A 16-contender claim test repeated for 20 runs covers the contention path; race-detector execution remains a CI follow-up where cgo is available.

## Boundary

T007 does not produce adapter claim/takeover receipts, schedule tasks or steps, decide argv/cwd/network/resource policy, or implement operator handoff. Callers must validate the current fencing token before every workspace write, recovery, or publication and persist references to this layer's evidence.
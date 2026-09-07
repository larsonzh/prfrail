# S0 T004 Technical Spike Report

[简体中文](s0-spikes.md)

Date: 2026-09-07. Verdict: `T004 COMPLETE`, while overall S0 remains `NOT_READY`. This report establishes primitive feasibility and failure boundaries only. It is not an S1 production implementation, release approval, or proof of live model capability.

## 1. Platforms and Method

| Platform | Exercised environment |
|---|---|
| Windows | Windows 11 Home 10.0.26100 amd64; D: on NTFS; Go 1.27.0 windows/amd64 |
| Linux | Ubuntu 24.04.4 LTS; Linux 6.8.0-85-generic x86_64; ext4 root filesystem; pinned Go 1.27.1 linux/amd64 |
| SessionBridge | Sister repository commit `d2298f3`; Python 3 contract tests and Node history engine; no model calls |

The probes used real local files and real child processes. Go file, queue, and crash probes ran on both platforms. Linux sources were transferred over SSH into `/tmp` and deleted after execution. Windows prototypes lived under the repository's ignored `tmp/` directory and were deleted after the report was completed. No whois assets were read, no paid service was called, SessionBridge was not modified, and no Git commit was created.

## 2. Results of the Twelve Slices

| # | Slice | Result and constraint |
|---|---|---|
| 1 | Go single binary | Windows and Ubuntu both built with `CGO_ENABLED=0` and ran `prfrail version`, producing `prfrail 0.1.0-dev`. The module's Go 1.22 baseline is unchanged. |
| 2 | TUI dependency | Bubble Tea v2.0.9 requires Go 1.25 and v1.3.10 requires Go 1.24, so neither meets the baseline. v1.3.4 requires Go 1.18 and is MIT-licensed, making it an S1 candidate, but it is not a root-module dependency. |
| 3 | TUI behavior | The isolated v1.3.4 prototype passed `q`-key exit, 12-column truncation, and ANSI-free `NO_COLOR` tests. Its CGO-free Linux binary produced the same output on Ubuntu. Full interactive-terminal accessibility remains T016 work. |
| 4 | Windows file primitives | On NTFS, file `Sync` and replacement without an open reader succeeded; directory-handle `Sync` returned access denied. `os.Rename` returned access denied while a normal reader held the destination open, so POSIX open-handle replacement semantics cannot be assumed. |
| 5 | Publication crash points | Both platforms behaved identically: after temporary-file sync but before rename, readers saw the old state and no receipt; after rename but before receipt, readers saw the new state but still no receipt; only receipt presence allowed a completed verdict. Both non-completed states must fail closed. |
| 6 | Windows process tree | Force-killing the direct child left its grandchild alive. `taskkill /T /F` against a known child stopped both. S1 should use an auditable Job Object or equivalent controlled-tree primitive and refuse force termination when PID ownership cannot be proven. |
| 7 | Windows paths | NTFS was case-insensitive; NFC and NFD names remained distinct; hard links could be created and detected with `os.SameFile`; an unprivileged symlink failed for lack of privilege. Go could create a path longer than 260 characters and a path named `CON`, so the canonical checker must explicitly reject reserved names and aliases. |
| 8 | Linux files and paths | On ext4, replacement, file `Sync`, directory `Sync`, symlinks, and hard links succeeded. The filesystem was case-sensitive and kept NFC/NFD names distinct. An open reader continued reading the old inode while the path resolved to the new content. |
| 9 | Linux process group | Killing the session leader left its grandchild alive; signaling the negative PGID of an isolated session stopped the group. Unknown ownership likewise fails closed. |
| 10 | File queue | With 32 workers, Windows returned success to 5 and 3 contenders in two attempts to rename one ready file to different claim names, so rename-only election is invalid; ext4 produced one winner. A fixed claim file created with `O_CREATE|O_EXCL` produced exactly one winner on both platforms. Generation 7 was rejected against current generation 8, a duplicate requestId applied once, and no temporary residue remained. |
| 11 | SessionBridge silent | `python tests/test_contract.py` passed 28/28; `node tests/test_history_engine.js` passed 12/12; extension syntax and encoding gates passed. Coverage included busy, timeout, wrong requestId/stale-result protection, same-ID retry, history reset/truncation, and deduplication. |
| 12 | Dual-channel fault matrix | The ProofRail file queue must durably own requestId, generation, claim/result, and dispatch/takeover receipts; SessionBridge is only the silent transport. In-memory cache cannot claim exactly-once after host restart, cache loss, or history truncation. The live VS Code LM API was not called and remains `unavailable/not exercised`; visible/GUI fallback is forbidden. |

## 3. Reproduction Commands and Verdicts

```powershell
$env:CGO_ENABLED='0'; go build -o tmp\prfrail-t004-windows.exe .\cmd\prfrail
go run .\tmp\t004_fs_probe.go
go run .\tmp\t004_crash_probe.go
go run .\tmp\t004_queue_probe.go
python .\tmp\t004_process_windows.py
go -C .\tmp\t004-tui test ./...
```

```sh
CGO_ENABLED=0 go build ./cmd/prfrail
go run /tmp/proofrail-t004-fs-probe.go
go run /tmp/proofrail-t004-crash.go
go run /tmp/proofrail-t004-queue-probe.go
/tmp/proofrail-t004-process.sh
NO_COLOR=1 /tmp/proofrail-t004-tui --render 12
```

Persistent regression commands:

```powershell
go build ./...
go vet ./...
go test ./...
npm --prefix tools/contracts test
python ..\sessbridge\tests\test_contract.py
node ..\sessbridge\tests\test_history_engine.js
node --check ..\sessbridge\extension\extension.js
python ..\sessbridge\tools\enforce_encoding.py
```

One-off probes were deleted when this task completed. The tables above preserve their inputs, contention scale, crash points, and verdict conditions. S1 must convert the same counterexamples into version-controlled platform tests rather than treating this report as a test substitute.

## 4. ADR Constraints and Residual Risks

- ADR-003: a single-file rename is neither a transaction nor durable completion. Use the order "temporary write and file sync -> platform replacement -> platform-supported directory/volume durability -> completed receipt." Windows must handle open readers and unavailable directory `Sync`; Linux findings cannot be extrapolated to Windows. Power-loss durability still needs native fault injection in T008/T017.
- ADR-004: queue election uses exclusive creation of a fixed claim, never rename-only. Every result must match both requestId and generation, and durable deduplication precedes side effects. SessionBridge cache/history supplies transport context, not core exactly-once evidence.
- TUI: if S1 adopts Bubble Tea, the v1.3.4 candidate requires a fresh dependency digest, license, and vulnerability review. Network restrictions required `GOSUMDB=off` and a mirror in the isolated probe directory; those settings must not enter production builds.
- Unavailable: the live VS Code LM API, Windows power-loss/volume flush, network filesystems, privileged symlinks, a production Job Object wrapper, and a complete interactive TUI were not verified. Those capabilities must remain unknown/unavailable and fail closed rather than degrade to success.
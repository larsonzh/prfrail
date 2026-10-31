# ProofRail Validation Report — P41-T027B-SPAWN-OWNERSHIP

| Field | Content |
|---|---|
| Slice | `P41-T027B-SPAWN-OWNERSHIP` |
| Type | implementation class `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `28f440f` |
| Date | 2026-10-31 |

## Goals and boundaries

- **Goal**: close the two orphan-process gaps C2 and N1 - when the host crashes after `cmd.Start()` and before `AssignProcessToJobObject`, the suspended child sits in no Job and has no durable record, becoming a permanent orphan that is neither visible nor reclaimable; and raise "ownership before execution" from a timing assumption to a structural property.
- **Closure shape**: `guard.ProcessSpec` gains an `OnSpawn` hook and `StartManaged` fixes the order to prepare -> start -> attach (verification) -> identity -> `OnSpawn` -> resume; on Windows `STARTUPINFOEX` + `PROC_THREAD_ATTRIBUTE_JOB_LIST` moves Job ownership to the instant of creation, `attach` degrades to an `IsProcessInJob` verification and `resume` becomes its own step; the adapter publishes `managed-process.starting.json` atomically inside the hook and registers under the `starting` phase, while the restart side `RecoverStarting` splits on "lock held => host alive, lock released => host dead" and Linux reclaims by adoption and kill.
- **Scope**: `internal/guard/**` and `internal/adapters/**`.
- **Boundary**: no touch to contracts / schema / fixtures / `ci.yml`; no change to directive business rules; no gitee push; `commit` / `push` stop at 10 for authorisation.

## Review chain conclusions

- **Round 1 architecture**: `deep-reasoner` + DeepSeek V4 Pro; the C2+N1 ownership and orphan-prevention design book, `DESIGN: DONE`.
- **Red tests**: `RED-CONFIRMED` (invariants I1 / I3 / I4 / I5 falsify first through counter examples, then the implementation starts).
- **2 implementation, two passes**: main body landed (`managedCommand` replaces `os/exec`, Job-at-creation, starting marker, `RecoverStarting`) -> the F fix pass (F-1..F-5 closed).
- **Gates**: `go build` / `go vet` / `go test ./...` green; tree gate `TOTAL_FAIL=0`.
- **Pre-review r1**: `PRE-REVIEW: PASS WITH FIXES` (2 Medium + 3 Low).
- **Pre-review r2**: `PRE-REVIEW: PASS` (every round-1 F-1..F-5 verified closed, no new Medium or above).
- **Scan**: `INDEPENDENT SCAN: PASS`.
- **Final review**: `FINAL REVIEW: PASS` (2 Low + 3 Info, none blocking).

## Six-stage scan transcript

Source `docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/06-scan/review.md` (archive SHA256 = `124ca7a2fcd6656add90b929b4577dec1057d85da3375cd86a23cfa08873ae75`); the fence transcribes it verbatim with zero rewriting, from `## Section A` through the conclusion line.

```sixstage
## Section A — 总体结论：本次只读扫描未发现中等及以上缺陷；结合本机 Windows 上的 `go test -count=1 ./...` 15 ok / 0 FAIL，以及关键顺序/锁语义/回收逻辑的一致性，整体结论为 PASS。

## Section B — 发现清单
Low · internal/adapters/agent_runner_pinned_cli.go:494-516 · clearAgentRunnerStartingMarker 仍通过 path 再删除；虽然有重锁与比对，但在极短窗口里，如果另一启动重命名新 marker 到同一路径，旧扫描仍可能误删新记录。 · 在 re-open 后做 inode/内容双重校验，并仅删已证明确的旧记录，而不是直接删除当前 path。
Low · internal/adapters/agent_runner_pinned_cli.go:349-365 · RecoverStarting 把“没有 marker”和“活宿主拒绝”压成同一 boolean false,nil，调用方无法区分“可继续启动”和“需要等待重试”。 · 用专用状态枚举或包装错误替代单一布尔，保留兼容性并保留拒绝语义。

## Section C — 逐条证据与推理
1. 启动顺序正确且不可绕过：在 internal/guard/process.go:186-238，StartManaged 的序列是 preparePlatformProcess → platform.start → attach → platformIdentity → OnSpawn → resume。OnSpawn 发生在 identity 已捕获之后、resume 之前，且失败会进入 abortStartup。
2. Windows Job 归属先于执行：internal/guard/process_windows.go:12-32 和 :207-283 明确把 CREATE_SUSPENDED + PROC_THREAD_ATTRIBUTE_JOB_LIST + EXTENDED_STARTUPINFO_PRESENT 绑定到一次 CreateProcessW；attach 只校验“进程确实在 Job 里”，resume 只是单独调用 NtResumeProcess。归属不依赖时序假设，而是由创建时 Job 属性保证。
3. Durable marker 语义是 lock-then-publish：temp + lock + rename；Windows LockFileEx 和 Unix flock 都是非阻塞独占锁，锁持有=宿主活着，锁释放=宿主已死。RecoverStarting 的拒绝/回收分支被明确写成“锁持有拒绝，锁释放后才解码并回收”，与测试用例一致。
4. P40 契约风险没有被破坏：abortStartup 使用有界 waitResult + ErrTerminationUncertain；Reclaim 仍要求“wait 结束 + verifyGone”才释放平台句柄；RunManaged 和 Terminate 仍沿用 bounded verification，不存在无界等待回流。
5. 现有测试与本地证据一致：TestStartManagedResumeIsSeparated 与三个恢复测试共同验证拒绝活宿主、杀死孤儿、第二次扫描 no-op；本机 Windows `go test -count=1 ./...` 15 ok / 0 FAIL，无中等风险回归。

## Section D — 可证伪性审计：这里至少有一个明确反例能使结论翻转。
- 直接反例 1：删除 internal/guard/process.go:223-237 中 `if spec.OnSpawn != nil { ... }` 与 `platform.resume()` 的顺序，让 resume 提前发生，或直接删除 `TestStartManagedResumeIsSeparated` 的 Windows assert。这样会直接破坏“OnSpawn 必须先于 resume”的核心不变量，结论会翻转。
- 直接反例 2：删除 internal/adapters/agent_runner_pinned_cli.go:503-516 中 `if !locked { continue }` 及重解码校验，允许 stale 扫描直接 `os.Remove(path)`。随后 `TestRecoverStartingRefusesLiveHost` / `TestRecoverStartingRefusesLiveHostUnix` 会立刻失败，证明原设定不成立。

## Section E — 我无法验证的部分（明确列出）与什么会推翻 PASS
- 我无法在本机 Windows 宿主上真实执行 Linux 运行时路径；internal/adapters/spawn_recovery_live_host_unix_test.go 的验证逻辑被 `//go:build unix` 隔离，必须在 Linux/CI 中执行。它们的成功与否是唯一真正影响 PASS 的外部证据。
- 我没有读取 CI 的 Linux job 日志；仅能确认本地 Windows 绿色。若未来 Linux 上执行 `GOOS=linux go test ./internal/adapters/...` 时出现失败，尤其 `TestRecoverStartingRefusesLiveHostUnix` / `TestRecoverStartingAdoptsOrphanUnix` 失败，PASS 会被推翻。
- 任何对关键顺序或锁条件的单行改动，尤其 process.go 的 resume 顺序或 pinned_cli.go 的 clear 条件，都会直接撬开当前结论。

INDEPENDENT SCAN: PASS
```
- Transcript note: the two Section B rows inside the fence were originally `Low | …` pipe-shaped; to satisfy the G7-f "severity rows ⇔ FINDINGS" closed set, their separators are normalized from `|` to `·`, all other text verbatim; both Low findings map one-to-one to §later-slice candidates and §known limitations.

## F-item closure

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| F-1 | `RecoverStarting` clears by path without re-verifying ownership | Medium | `internal/adapters/agent_runner_pinned_cli.go:493` | Closed: re-open (`openStartingMarkerForClear`) + retried lock + strict re-decode + full-structure comparison; a held lock or mismatched content skips the clear |
| F-2 | the Linux reclaim path carries no runtime coverage | Medium | `internal/adapters/spawn_recovery_live_host_unix_test.go:113,173` | Closed: a live-host refusal case and an adoption-kill case, gated by `//go:build unix` and executed by the Linux CI leg |
| F-3 | the lock boolean result is ignored | Low | `agent_runner_starting_marker_unix.go:118` / `_windows.go` in the same shape | Closed: the three `tryLockStartingMarker` call sites no longer leak `(false, nil)` |
| F-4 | a corrupt marker clears the file but not the slot | Low | `internal/adapters/agent_runner_pinned_cli.go:430` | Closed: the corrupt branch renames away -> reads the owner -> deletes the staging copy only after proving no live mirror, and the slot is reclaimed together |
| F-5 | the refusal state is indistinguishable from "no marker" | Low | `internal/adapters/agent_runner_pinned_cli.go:357` | Closed and registered as-is: the comment states the `(false,nil)` ambiguity and declares `deliberately not changed`; the same-origin item is carried by 7 F-01 |
| 5 Low-1 | the starting phase relies on slot freshness to prevent a second process | Low | pre-review r2 / pre-review r1 F-5 pairing | Registered: later-slice candidate, evaluate returning an error on the refusal state |
| 7 F-01 | the refusal state returns the same shape, and an empty-launchId slot branch deletes staging on the mirror check alone | Low | `internal/adapters/agent_runner_pinned_cli.go:350-359,450-457` | Registered: same origin as F-5; a later slice moves to a tri-state return or a dedicated sentinel error |
| 7 F-02 | `clearAgentRunnerStartingMarker` deletes by path and carries a rename-overwrite window | Low | `internal/adapters/agent_runner_pinned_cli.go:493-530` | Registered: the same item as pre-review Low-2; a later slice moves to object-identity-bound deletion (rename-away plus re-verification) |
| 7 F-03 | the mirror-deletion branch of the successful recovery path carries no explicit assertion | Info | `internal/adapters/agent_runner_pinned_cli.go:414-416` and the two recovery tests | Registered: assertions added by a later slice |
| 7 F-04 | the implementation does not use `Pdeathsig`, so Linux crash reclaim depends on the next scan | Info | `internal/guard/process_linux.go:44-48` and the unix tests | Registered: publishable as-is; shortening the delay becomes its own slice |
| 7 F-05 | the D-1..D-6 deviation ledger and the evidence-pack metadata stay skeletal | Info | pre-review L110 and `MANIFEST` / `SHA256SUMS` | Closed: this report section Deviations and registrations fills in D-1..D-6; `MANIFEST` updated and `SHA256SUMS` recomputed |

## Deviations and registrations

**Model bindings and carriers**

| Stage | Carrier | Model |
|---|---|---|
| 1 architecture | `deep-reasoner` | DeepSeek V4 Pro |
| 2 implementation | `prfrail-implementer` | DeepSeek V4.1 Flash |
| 3 tests | `prfrail-tester` | DeepSeek V4.1 Flash |
| 5 pre-review | `deep-reasoner` | DeepSeek V4 Pro |
| 6 scan | `independent-reviewer` | MAI-Code-1.1-Flash |
| 7 final review | `independent-reviewer` | GPT-5.3-Codex |
| 8 documentation | `prfrail-documenter` | DeepSeek V4.1 Flash |
| batch | `deep-reasoner` | GLM-5.3-Flash (2 batch assists this slice: spawn-point census + test-output analysis) |

**Deviation ledger**

| Id | Content | Source | Verdict |
|---|---|---|---|
| D-1 | `managedCommand` replaces `os/exec` for creation and waiting | pre-review r1 deviation 1 | Safe: a single `exitCode()` writer with happens-before established through `waitDone` |
| D-2 | `preparePlatformProcess` signature widened to carry the Job and the attribute list | pre-review r1 deviation 2 | Safe: the compatibility surface is covered by the full suite |
| D-3 | Linux substitutes `flock` as the equivalent of `LockFileEx` | pre-review r1 deviation 3 | Safe: the kernel drops the lock when the process dies and the child does not inherit it |
| D-4 | Windows does not kill when the lock is released but the process is alive, and keeps the mirror | pre-review r1 deviation 4 | Safe: fail-closed, and `refuseLiveMirrorProcess` blocks the next launch |
| D-5 | platform helper: `syscall.Open` carries no `FILE_SHARE_DELETE` on either Go 1.27 or CI 1.22 | pre-review r2 records it as D-1 | Necessary and holding: a dedicated helper for the delete path (`openStartingMarkerForClear`), evidenced by the green marker-deletion assertions in the native Windows recovery cases |
| D-6 | the owner release path does not re-verify the owner | pre-review r2 records it as D-2 | Boundary acceptable: the slot protocol rules out concurrent launches on the same run root |

**Windows syscall measurements**

- `EXTENDED_STARTUPINFO_PRESENT` is mandatory: without it `CreateProcessW` silently ignores `lpAttributeList` and Job-at-creation stops applying (`internal/guard/process_windows.go:24-30`).
- Job-at-creation coexists with `CREATE_SUSPENDED`: one `CreateProcessW` call carries the `JOB_LIST` attribute, `EXTENDED_STARTUPINFO_PRESENT` and `CREATE_SUSPENDED`, so the child is inside the Job at the instant of creation and executes no instruction.
- `LockFileEx` survives rename: the temp file is locked before the rename publishes it, the share mode allows the delete that a rename performs, and the lock follows the file object identity onto the published path.
- The falsifier of design risk 1 therefore does not apply: Windows keeps the proven path (creation-time Job ownership plus the lock-liveness criterion).

**Two self-inflicted defects fixed during 3**

- Double close of the pipe write end: the parent write end is closed through `defer closeParent` right after `CreateProcessW` returns and the output is drained inside `waitFn`, so the duplicate close path is gone.
- `HANDLE_LIST` inheritance: the inheritance list carries the three standard handles only, the same policy as `os/exec`; handing over the whole inheritable set had polluted the child handle table.

**Untestable items registered**

- I2 with its four seams is not built: host-crash injection needs a test seam beyond this slice, registered as a later slice.
- I1 kill-on-close is half proven: "the child does not execute while the hook blocks" is pinned; the kernel killing it through handle close after the host dies is corroborated by the native Windows case, and the full half proof stays with a later slice.
- Linux is not executed natively: the `//go:build unix` cases cannot run here and are covered by the Linux CI leg.

## Gate readings

- Three-scope readings: tree `TOTAL_FAIL=0`, index `TOTAL_FAIL=0`; the ci scope runs authoritatively inside CI (`--scope=ci`) — local ci degenerates to `HEAD^1..HEAD`, see the ⑩ receipt
- SELFTEST: PASS 697/697; the G6-1 scan count is stated by the criterion output itself (reference form, avoids self-reference drift)

## Commit and CI

- Main commit `8799015` (feat: own spawned processes before resume and recover starting orphans (P41-T027B-SPAWN-OWNERSHIP)), range `28f440f..8799015`; CI run `41179732685` **double green**: windows 2m56s / ubuntu 2m40s; both legs `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=26`, `TOTAL_FAIL=0`; `internal/guard` and `internal/adapters` ok on both legs; the ubuntu leg also ran Race and the F-2 unix recovery tests (they pass within the adapters package ok).
- Docs-only commits after this sentence no longer back-fill this section; readings are carried by the ⑩ receipt.

## Later-slice candidates

- I1 kill-on-close half proof: host-crash injection plus a zero-residual-PID assertion.
- I2 four-seam injection: one host crash each after Start -> attach -> OnSpawn -> resume.
- F-5 / 7 F-01 tri-state refusal: replace `(false, nil)` with a sentinel error or a state enumeration.
- 7 F-02 rename-away object-bound deletion: remove the overwrite window of deletion by path.
- 7 F-03 mirror-deletion assertion: add explicit assertions on the successful recovery path.
- `Pdeathsig` as its own slice: narrow the Linux crash-reclaim delay from "restart interval" to "killed by the kernel".
- 5 Low-1 slot-freshness backstop review: re-check the strength of the starting-phase guard against a second process.

## Known limitations

- The substance of section E in the scan: native Windows green cannot replace Linux runtime evidence, `spawn_recovery_live_host_unix_test.go` is gated by `//go:build unix`, and its outcome is the only external evidence that can really move PASS; a single-line change to the critical order or lock condition turns the current conclusion over.
- The registrations above - I2 not built, I1 half proven, Linux not run natively, D-1..D-6, 7 F-01..F-05 - stay registered for the validity period of this slice's conclusion.
- Windows Job-at-creation depends on the unexported `procThreadAttributeJobList` and a raw `syscall.CreateProcess` call with no `os/exec` equivalent; any later replacement of the creation path must re-check this group of assertions.
- The starting-marker lock criterion approximates "the host is alive" with "the file handle is alive"; the two platforms differ in meaning (Windows `LockFileEx` / Linux `flock`) and must not be extrapolated beyond this literal wording.

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/p41-spawn-ownership.md	repo
docs/validation/p41-spawn-ownership_EN.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/MANIFEST.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/SHA256SUMS.txt	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/05-prereview/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/06-scan/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/07-final/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/freedom-list.md	repo
internal/adapters/agent_runner_pinned_cli.go	repo
internal/adapters/agent_runner_process_registry.go	repo
internal/adapters/agent_runner_starting_marker.go	repo
internal/adapters/agent_runner_starting_marker_unix.go	repo
internal/adapters/agent_runner_starting_marker_windows.go	repo
internal/adapters/spawn_recovery_live_host_unix_test.go	repo
internal/adapters/spawn_recovery_live_host_windows_test.go	repo
internal/adapters/spawn_recovery_test.go	repo
internal/guard/process.go	repo
internal/guard/process_linux.go	repo
internal/guard/process_unsupported.go	repo
internal/guard/process_windows.go	repo
internal/guard/process_spawn_ownership_test.go	repo
internal/guard/process_spawn_ownership_windows_test.go	repo
tools/gates/cjk-newwords.txt	repo
```

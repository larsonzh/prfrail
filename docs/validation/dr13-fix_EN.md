# ProofRail Validation Report — DR13-FIX

| Field | Content |
|---|---|
| Slice | `DR13-FIX` |
| Type | fix class `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `8e834b2` |
| Date | 2026-11-01 |

## Goals and boundaries

- **Position**: this slice is the prerequisite fix for P45; it takes no v34 queue 1-4 number; v34 stays the baseline "P44 closed, before P45 starts" and v35 will fold this slice in (user ruling 2026-11-01).
- **Goal**: diagnose the flaky root cause of `TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner`, separate "the test rig is too strict (A)" from "a real exclusion-window defect (B)", and deliver a falsifiable fix plus a counterexample.
- **Scope**: `internal/adapters/agent_runner_pinned_cli.go` (the exclusion window) and `internal/adapters/spawn_recovery_test.go` (the counterexample).
- **Boundary**: no other package (`internal/chain/**`, `internal/gates/**`, `internal/console/**`), no schema, fixtures or `ci.yml`; no gitee push; `commit` / `push` stop at 10 for authorisation.

## Diagnosis

- **DIAGNOSIS: B**: a real exclusion-window defect. On a concurrent Start inside one process and one launcher, the `RecoverStarting` proof "the marker lock is free, so the previous host is dead" no longer holds; the unix branch stops that "orphan", kills the live child of the first winner and takes its slot, so two Starts succeed one after the other.
- **A is excluded**: no legal slot-release point is triggered between the two successes - no Stop call, the stub `run` mode sleeps for 60 seconds, the launcher never reaps, and the first winner never failed; the slot can only disappear while the first winner is alive through a mis-judged reclaim, and reclaim treats a live process as a crashed orphan and kills it - an implementation defect, not a rig that is "too strict about treating successes across time as live processes at one instant".
- **The assertion must not be relaxed**: the assertion at `agent_runner_pinned_cli_test.go:901-902` is the protocol invariant itself (an atomic slot means the second Start must fail while the first winner is alive), so this slice fixes the implementation and leaves the assertion untouched.
- **Evidence chain summary**:
  - the entry gate at `agent_runner_pinned_cli.go:124` runs before the claim, so 8 contenders pass it together (registration only happens in OnSpawn) ⇒ exclusion rests on the slot alone.
  - the `tryLock` at `agent_runner_pinned_cli.go:345` does not re-check the path after it succeeds, so the lock can land on an inode that was already unlinked.
  - `reclaimStartingProcess` in `agent_runner_starting_marker_unix.go` runs `StopProcessIdentity` when `alive=true` and returns `(true,nil)` - the ubuntu-leg kill root cause; the same function in `agent_runner_starting_marker_windows.go` returns `(false,nil)` when `alive=true` (fail closed, never kills), which explains why the windows leg stayed green.
  - `release()` at `agent_runner_starting_marker.go:89-99` = unlock → close → `os.Remove`, so the lock is released strictly before the path disappears and a scanner can observe "free lock with the path still present" (the W1 window).
  - both the fail-closed and the full-reclaim branches call `reclaimLaunchSlot` and clear the marker, which is the slot-stealing path; the loser allow-list at `agent_runner_pinned_cli_test.go:894-896` only carries Conflict / Duplicate, matching the observed `got 2 (losers 6)` verbatim.
- **W1 timing summary** (independent of any special flock kernel semantics): A wins the slot and spawns → A's OnSpawn publishes the marker and registers → B opens the marker inside the millisecond window "marker published, A has not called `closeStartingMarker` yet" and is preempted → A releases the lock, runs `os.Remove` on the marker and returns success (the stub is still alive) → B resumes and locks the already-unlinked orphan inode → decode succeeds and `ProcessAlive(A)` is true → the unix branch kills A's live stub and takes its slot → B's claim sees no slot and a dead mirror and is let through ⇒ winner #2; the other 6 all hit the Conflict / Duplicate sentinels.
- **W2 timing summary** (an amplifier, kernel dependent): after its `tryLock` fails B closes its own FD, and under man flock(2) "closing any fd that refers to the file releases every flock the process holds on it" A's lock is released early by that close, so "free lock + live process" spans A's whole shutdown window and the hit rate rises sharply; W1 stands on its own and the logs cannot tell whether the two CI reds went through W1 or W2.
- **Fix trade-offs**: D1 is the main fix - `RecoverStarting` gains a same-process ownership gate after the successful decode and before the liveness decision (while this launcher has a live handle registered for `marker.LaunchID`, that launch is never reclaimed), which is deterministic, covers both platforms, touches no protocol and changes no CONTRACTS / Schema; D2 is the companion - after `tryLock` succeeds, `os.Stat(path)` re-checks the path and `os.ErrNotExist` means the lock sits on an already-unlinked orphan inode, so it returns `(false,nil)` and closes the cross-process W1 residue; D3 (make `release` run `os.Remove` before unlock) touches the success / failure / retire hot paths and depends on Windows remove-while-open semantics, so this slice does not do it and registers it for a later slice.

## Review chain conclusions

- **1 diagnosis**: `deep-reasoner` + DeepSeek V4 Pro; closed-set verdict `DIAGNOSIS: B` with a 15-item file:line evidence chain and 5 inputs for P45.
- **2 and 3 fix and counterexample**: D1 / D2 and `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` land; the counterexample is deterministically red on the **unfixed code** - the Windows failure point is `spawn_recovery_test.go:281` (the marker assertion) and unix is expected to go red on the "child still alive" assertion (not executed locally).
- **4 gates**: `TOTAL_FAIL=0`; `go build` / `go vet` / `go test ./...` all green, readings in section Gate readings.
- **5 r1**: `PRE-REVIEW: PASS WITH FIXES` (0 Critical / High; 2 Medium + 1 Low + 2 Info).
- **5 r2**: `PRE-REVIEW: PASS` (B3 closed; the B1 / B2 / B4 / B5 dispositions approved; obligations N1-N4 added for this report).
- **6**: `INDEPENDENT SCAN: PASS`, Section B reads `0 · none`.
- **7**: `FINAL REVIEW: PASS`, 0 new findings, with 6 release conditions (landed item by item in section Deviations and registrations).
- **R / D**: R = 0, D = 0 (not inherited from the previous slice).

## Six-stage scan transcript

Source `docs/validation/evidence/DR13-FIX/06-scan/review.md`; the fence transcribes the ⑥ body in A→B→C→D→E order and two normalisations are disclosed outside it.

```sixstage
## Section A 扫描范围与输入
INDEPENDENT SCAN: PASS
范围：只读核查 internal/adapters/agent_runner_pinned_cli.go、agent_runner_process_registry.go、spawn_recovery_test.go 与 spawn_recovery_live_host_unix_test.go，聚焦 D1 / D2 落点、反例构造与并发回归。

## Section B 发现清单
| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
|---|---|---|---|---|
|  |  |  |  |  |

## Section C 与 ①/⑤ 的一致性核对
1. D1 的正确性与充分性：成立。`registry.lookup` 是进程内唯一来源的“同 launcher 仍持有该 launchID”判定，且持有 `registry.mu` 读出句柄后不做任何写操作；同进程所有权判断位于 `RecoverStarting` 的 decode 成功后、liveness 判定前，正好覆盖“锁空闲但同进程仍持有启动句柄”的关键窗口 [internal/adapters/agent_runner_process_registry.go](internal/adapters/agent_runner_process_registry.go#L141-L149) 和 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L326-L399)。成功路径上，`StartAgentRunnerProcess` 只在 `advance(..., running)` 后调用 `handle.closeStartingMarker()`，并不会在成功返回前 `release` 该 handle；因此，注册表中的活句柄确实仍在，并在同进程竞态下拦截误回收 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L254-L293)。  
   结论：D1 不仅正确，而且对当前修复边界是充分的；它锁定了“同进程 handle 仍在”的死角，而不是放大到跨进程 orphan 的正常回收。
2. D1 是否误伤孤儿收养：成立。`registry` 纯内存，仅存在于当前进程；崩溃后它必然消失，所以跨进程 orphan 不会命中 D1。相反，正路径仍保留：`TestRecoverStartingIdempotent`、`TestRecoverStartingTokenMismatchFailsClosed`、`TestRecoverStartingRefusesLiveHostUnix`、`TestRecoverStartingAdoptsOrphanUnix` 继续覆盖“已死主机/TokenMismatch/活主机/孤儿收养”的正确分支；D1 只在 same-process ownership 上做 fail-closed—且这正是其目标 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L135-L210) 和 [internal/adapters/spawn_recovery_live_host_unix_test.go](internal/adapters/spawn_recovery_live_host_unix_test.go#L113-L198)。  
   结论：D1 不会误伤 orphan 收养；它只拦截“本进程仍持有该 launch 的句柄”。
3. D2 的正确性：成立。`os.Stat(path)` 在 `tryLock` 成功后立刻做；若 `ErrNotExist`，说明锁已被取到的 inode 已脱离当前 published path，属于`orphan inode` 的典型场景，必须返回 `(false, nil)` 并清理 lock/FD，而不是继续 decode 或 kill [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L357-L377)。这正是把“锁空闲但 marker 不再存在”识别为“无可回收记录”的 fail-closed 语义，不构成误判；它并不掩盖真实错误，因为只有 `NotExist` 才是此路径的合法闭环，而其它 I/O 错误仍会返回 `ErrAgentRunnerLaunchConflict`。  
   结论：D2 对“orphan inode”窗口是正确的，且对调用方不构成误判。
4. E1 的可证伪性：成立。新增测试 `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` 使用了真实 helper 进程、真实 `claimLaunchSlot`、真实 marker 文件和真实 registry 注册，断言四件事全部落在“同进程仍 owns launch”这个对错分界上：`reclaimed == false`、helper 继续活着、marker 仍在、slot owner 仍为 launchID [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
   这不是假绿通道：若去掉 D1 门，旧实现会落入“live process + free lock”的回收分支，触发 `StopProcessIdentity` 杀活进程或删 slot；这正与 `reclaimStartingProcess`/`reclaimLaunchSlot` 分支的语义一致 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L406-L417)。  
   结论：该测试是可证伪的，并且确实覆盖了真正的回归窗口。
5. 并发回归：成立。`TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner` 在同一 launcher / 同一 run root 下保留原逻辑；其设计目标是确保恰一名赢家，且其余竞争者返回 `ErrAgentRunnerLaunchConflict` 或 `ErrAgentRunnerProcessDuplicate` [internal/adapters/agent_runner_pinned_cli_test.go](internal/adapters/agent_runner_pinned_cli_test.go#L857-L920)。在当前实现下，适配器包测试已通过，说明这是结构性约束，而不是偶然性；并发 winner 语义依赖的仍是“slot O_EXCL”与 registry ownership gate，不依赖测试抖动。  
   结论：`winners == 1` 的结论结构性成立。
6. 边界与卫生：成立。允许路径仅限 `internal/adapters/agent_runner_pinned_cli*.go` 与其测试文件；证据包则落于 `docs/validation/evidence/DR13-FIX/`，并无新增/删改或临时脚本落盘；代码中无 TODO，且新增断言均为英文说明 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
   只读验证证据：`go test ./internal/adapters` 输出为 `ok github.com/larsonzh/prfrail/internal/adapters 35.127s`；`go vet ./internal/adapters` 输出为空（无 vet issue）；`gofmt -l` 也未输出任何文件，说明当前代码已格式化。  
   证据包自洽：`MANIFEST.md` 列出的三个载荷与 `SHA256SUMS.txt` 一致，且对应文件均存在 [docs/validation/evidence/DR13-FIX/MANIFEST.md](docs/validation/evidence/DR13-FIX/MANIFEST.md) 与 [docs/validation/evidence/DR13-FIX/SHA256SUMS.txt](docs/validation/evidence/DR13-FIX/SHA256SUMS.txt)。  
   结论：卫生与证据三方一致，扫描通过。

## Section D 未覆盖面与残余风险
1. `StopAgentRunnerProcess` 的失败分支仍能在“输入解析/停止失败”时无条件 `releaseLaunch`; 这不等于 D1 失效，但它是下一个最值得加固的边界，尤其在停止中途状态不一致时可能清掉 `registry` 句柄 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L271-L301)。  
   重要性：它是 stop/start 竞态的镜像窗口，若下一轮再复现 stop-while-starting，可能出现“已杀/未杀状态模糊”。
2. D2/TOCTOU 仍依赖平台层对 `flock`/inode 语义的一致性；本地仅能读代码而不能在 unix CI 上直接探针 `flock` 的行为差异，因此 W2 的内核级细节未被本机独立实证 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L343-L377)。  
   重要性：它关乎“跨进程锁释放是否会让另一个扫描持到 orphan inode”的真实发生率。
3. 证据包的 `01-diagnosis` 是已归档的定性结论，但并未在本机 Unx 宿主上复现 `StopProcessIdentity` 对活进程的具体行为；当前的证据仍依赖 CI/设计推演 [docs/validation/evidence/DR13-FIX/01-diagnosis/review.md](docs/validation/evidence/DR13-FIX/01-diagnosis/review.md)。  
   重要性：这里是最终发布前最需要的真实 runtime 证据，尤其是 ubuntu 腿。
4. D1 仅覆盖“同一 launcher 仍持有同一 launchID”的漏洞；它不覆盖“跨 launcher、跨进程、同一 runRoot 但镜像已失效”的不同窗口，后者仍需要 slot/镜像一致性检查和更严格的 stop 语义 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L620-L700)。  
   重要性：这类窗口不会触发 D1，但仍可能导致不完全的“owner 检查”。

## Section E 结论要旨
可证伪性成立：当前修复的核心是两层“先证实 marker 仍指向真实 published path，再确认 launcher.registry 没有同 launchID 活句柄”。这两个门都可以被单步反证：

- 若删掉 D1 门，`TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` 会直接失败，且 helper 进程会被误杀或 slot 被删除 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
- 若删掉 D2 门，`tryLock` 成功后到 `os.Stat` 之间的 orphan inode window 就能导致误回收；这在旧实现下会在“marker path 已 unlink”时制造 false positive [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L343-L377)。  
- 若没有真实 helper 进程、真实 slot claim、真实 marker 写入，测试就会失去证伪性；当前测试恰好使用了实际的 live process 和 slot 创建，不存在假绿通道。  
- 这一点也与 `go test` 结果一致：当前实现通过 package 级回归测试，说明反例是可复现且被封住的。

INDEPENDENT SCAN: PASS
```

- Normalisation 1: the ⑥ headings `## A 结论` / `## B 缺陷` / `## C 核验结论` / `## D 未覆盖面` / `## E 可证伪性` are rewritten inside the fence to this template's section names; the body of all five sections is transcribed verbatim.
- Normalisation 2: Section A holds two lines - ⑥'s overall verdict line (verbatim) and a scope line (taken from the file set Section C examined); ⑥ Section E's last line `§A INDEPENDENT SCAN: PASS` is normalised to the fence's last line `INDEPENDENT SCAN: PASS`.

## Item closure

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| 5 r1 B1 | every reading came from the local Windows host, the `//go:build unix` cases do not compile here, and DR-13's original reproduction leg is ubuntu | Medium | `spawn_recovery_live_host_unix_test.go` | Accepted: 9 uses the Ubuntu leg as the final judgement, registered in section Deviations and registrations R2 |
| 5 r1 B2 | D2 has no dedicated deterministic test; the only triggering window depends on scheduling | Medium | `agent_runner_pinned_cli.go:367-376` | Accepted: registered as "probabilistic coverage only", see section Deviations and registrations R3 |
| 5 r1 B3 | `io/fs` served only `fs.ErrNotExist`, inconsistent with the other `os.ErrNotExist` sites in the file | Low | `agent_runner_pinned_cli.go:370` | Closed: switched to `os.ErrNotExist` and dropped the `io/fs` import, 6 sites across the file |
| 5 r1 B4 | the slice registration row exceeds the declared boundary | Info | `docs/t027/REMAINING_SLICES.md` and `_EN` (`+23` each) | Accepted: start-up bookkeeping, exempted as a boundary item |
| 5 r1 B5 | the counterexample's "free lock + live process" state is simulated; the defensive value holds but it is not a reachable timing reproduction | Info | `spawn_recovery_test.go:221-290` | Accepted: the test comment states this plainly, no code change |
| 5 r2 N1-N4 | report-level obligations: ubuntu package-level `ok` lines, the B2 / D3 registrations, the wording fixes folded into the boundary log | - | this report | Closed: see section Deviations and registrations R1-R6 and section Known limitations |
| 6 | Section B reads `0 · none`, all six Section C verifications hold | - | `06-scan/review.md` | Closed: `INDEPENDENT SCAN: PASS` |
| 7 | `FINAL REVIEW: PASS`, 0 new findings, with 6 release conditions | - | `07-final/review.md` | Accepted: landed item by item in section Deviations and registrations R1-R6 |

## Delivery

- `internal/adapters/agent_runner_pinned_cli.go` (`+37/−0`): the D2 path re-check - after `tryLock` succeeds, `os.Stat(path)`, and `os.ErrNotExist` means the lock sits on an orphan inode so it returns `(false,nil)`, while every other I/O error is wrapped as `ErrAgentRunnerLaunchConflict` (`:367-376`); the D1 same-process ownership gate - after a successful decode, a hit from `registry.lookup(marker.LaunchID)` unlocks, closes and returns `false, nil` (`:390-397`). Both carry design comments and change no protocol behaviour.
- `internal/adapters/spawn_recovery_test.go` (`+84/−0`): the counterexample `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` (`:221-290`) builds the "free lock + live process + registered in this process" triple state from a real helper process, a real `claimLaunchSlot`, a real `registry.register` and an unlocked marker; its four assertions are `reclaimed==false`, the helper still alive, the marker still present and the slot owner still that launchID; `t.Cleanup` reaps the child and the durable records.
- Evidence pack `docs/validation/evidence/DR13-FIX/`: `MANIFEST.md` and `SHA256SUMS.txt` are metadata and the other 7 files are payloads (`00-briefing.md`, `freedom-list.md`, `01-diagnosis/review.md`, `03-tests/red-baseline.md`, `05-prereview/review.md`, `06-scan/review.md`, `07-final/review.md`).
- Ledger: the DR13-FIX summary row and detail block in `docs/t027/REMAINING_SLICES{,_EN}.md` are flipped to the completed state, and the DR-13 disposition cell now records the fix while keeping the re-open wording.

## Deviations and registrations

| Id | Content | Source | Verdict |
|---|---|---|---|
| R1 | Ubuntu-leg evidence: the package-level `ok` lines of the single-winner concurrency case and the unix recovery case must be extracted by 9 | 7 release condition 1 | Accepted: handed to 9, readings appended in section Commit and CI |
| R2 | the local host cannot run the `//go:build unix` cases, so the ubuntu leg is the final judgement | 7 release condition 2 and 5 r1 B1 | Accepted: registered as B1; either leg turning red sends the slice back to 5 |
| R3 | D2 has no dedicated "orphan inode" deterministic test and is covered only by D1's deterministic counterexample plus the concurrency regression | 7 release condition 3 and 5 r1 B2 | Accepted: `tryLock` is non-blocking on both platforms, so a zero-cost injectable boundary does not exist |
| R4 | D3 (make `release` run `os.Remove` before unlock) is not done in this slice | 7 release condition 4 | Accepted: it touches the success / failure / retire hot paths and relies on Windows remove-while-open, registered for a later slice |
| R5 | cross-process W1 residue: D1 does not cover a failed `os.Remove` combined with a cross-process scan | 7 release condition 5 | Accepted: the applicable boundary is same runRoot, cross launcher |
| R6 | the four uncovered areas of 6 Section D | 7 release condition 6 | Accepted: all four are registered in section Later-slice candidate registrations and section Known limitations, not counted as defects here |
| D-a | 5 r1 B4: the ledger registration rows landed before 2 and 3, beyond the declared code boundary | 5 r1 B4 | Accepted: start-up bookkeeping (the slice entry and the evidence-pack skeleton) with no code impact |
| D-b | 5 r2 wording fixes: `os.ErrNotExist` occurs 6 times in the file (the original 5 untouched plus 1 D2 branch unified from `fs.ErrNotExist`), and "no Chinese comments" holds for the added lines only | 5 r2 B-r2-1 and B-r2-2 | Accepted and corrected: this report uses the corrected wording in both places |

## Models and metering

| Stage | Model |
|---|---|
| 1 diagnosis | DeepSeek V4 Pro (deepseek) |
| 2 implementation | DeepSeek V4.1 Flash (deepseek) |
| 3 tests | DeepSeek V4.1 Flash (deepseek) |
| 5 pre-review (r1 / r2) | DeepSeek V4 Pro (deepseek) |
| 6 independent scan | MAI-Code-1.1-Flash (copilot) |
| 7 final review | GPT-5.3-Codex (copilot) |
| 8 documentation | DeepSeek V4.1 Flash (deepseek) |

- **Batch assist**: `GLM-5.3-Flash (glm)` dispatched 4 times - batch 1 exclusion-path census, batch 2 test rig, batch 3 log comparison, batch 4 file-operation census; batch 5, the diff cross-check, served as an input to the 5 r2 line-by-line check.
- **Primary model**: `DeepSeek V4.1 Flash` (High tier); 1 and 5 are still carried by V4 Pro under the whitelist, which is a separate matter from the primary-model choice.

## Gate readings

- Command: `node tools/gates/gate.js --all --skip-base`; tree scope `TOTAL_FAIL=0` (the archived artifacts were `git add`-ed precisely as G7-b requires before the run).
- `G3(a)`: the added / deleted line counts are equal on both sides of the ledger and on both sides of the report pair; the report pair is equal in all six shapes (line count / heading line-number sequence / `**` total / `|` total / blank lines / fence and backtick counts).
- `G4b`: `unknown=[]`; `G6-4`: `no expiring literal`; `G6-8`: the DR set carries 13 with no collision or dangling reference.
- `go build ./...` and `go vet ./...` exit 0; `gofmt -l` prints nothing.
- `go test -count=1 ./...`: `ok` lines 15, `FAIL` lines 0 (19 package entries, of which 4 are `[no test files]`); the start message's "16 packages" is one off the measurement here, so the measurement is registered.

## Commit and CI

- Commit `b123db1` is pushed to origin; CI `41269372878` is green on both legs on the **first** attempt (DR-13 did not reproduce); the ubuntu leg yields package-level ok lines for `internal/adapters` in both the Test and Race steps (12.343s / 89.396s) and the windows leg in 34.937s; both legs report `SELFTEST 697/697`, `GATE REPORT scope=ci changed=15` and `TOTAL_FAIL=0`.

## Later-slice candidate registrations

- D3 remove-before-unlock: make `release` run `os.Remove` before unlock to close W1 together with D2, which needs a full Windows remove-while-open regression.
- Cross-process W1 residue: add an observable warning or a compensating path for a leftover marker when `os.Remove` fails during a cross-process scan.
- `StopAgentRunnerProcess` failure branch: stop calling `releaseLaunch` unconditionally on a parse / stop failure and branch on the proven state instead.
- W2 flock kernel semantics: prove on unix that closing any fd releases every flock the process holds on that file.
- unix `StopProcessIdentity` against a live process: not reproducible on this Windows host, handed to CI / a unix host for evidence.
- Cross-launcher / cross-process mirror-invalid windows: add slot and mirror consistency checks to cover the windows D1 cannot reach.

## Known limitations

- **B1**: the local host cannot run `//go:build unix` cases, so the post-fix judgement of the unix branch (kill path and orphan adoption) rests on the Ubuntu leg.
- **B2**: D2 has no dedicated deterministic test and is covered only probabilistically by the D1 counterexample and the concurrency regression; `tryLock` is non-blocking on both platforms, so no external actor can deterministically hold the pre-stat window.
- **D3**: `release` still removes after unlocking, so the millisecond window where the path disappears after the lock is released remains on the cross-process side; this slice intercepts it deterministically only in-process (D1).
- **W1**: the cross-process residue needs a failed `os.Remove`, is very unlikely and is not closed here; the applicable boundary is same runRoot, cross launcher.

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/dr13-fix.md	repo
docs/validation/dr13-fix_EN.md	repo
docs/validation/evidence/DR13-FIX/00-briefing.md	repo
docs/validation/evidence/DR13-FIX/01-diagnosis/review.md	repo
docs/validation/evidence/DR13-FIX/03-tests/red-baseline.md	repo
docs/validation/evidence/DR13-FIX/05-prereview/review.md	repo
docs/validation/evidence/DR13-FIX/06-scan/review.md	repo
docs/validation/evidence/DR13-FIX/07-final/review.md	repo
docs/validation/evidence/DR13-FIX/08-ci/ci-evidence.txt	repo
docs/validation/evidence/DR13-FIX/MANIFEST.md	repo
docs/validation/evidence/DR13-FIX/SHA256SUMS.txt	repo
docs/validation/evidence/DR13-FIX/freedom-list.md	repo
internal/adapters/agent_runner_pinned_cli.go	repo
internal/adapters/spawn_recovery_test.go	repo
```

# ProofRail Validation Report — P40-T027B-LOCK-REAPER

| Field | Content |
|---|---|
| Slice | `P40-T027B-LOCK-REAPER` |
| Type | implementation class `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `38a5cba` |
| Date | 2026-10-31 |

## Goals and boundaries

- **Goal**: remove the two structural risks named by C1 - the unbounded reclaim waits in `Run` / `Terminate`, and the global lock order that carried no mechanical constraint.
- **Closure shape**: one reaper per process (`waitOnce` / `waitDone`) replaces the direct `Wait` at every call site; stop verification is decided through both `stopVerificationContext(grace)` and `verifyGone`, and reaching the bound returns `ErrTerminationUncertain` without releasing the platform handle or the launch slot; `runWithLock` passes the lock-order guard before taking the in-process lock and the cross-process file lock, and a non-increasing re-entry on the same stack is rejected outright.
- **Scope**: `internal/guard/**` and `internal/tickets/durable_cost_*`.
- **Boundary**: no touch to contracts / schema / fixtures / CI files; no change to directive business rules; no gitee push; `commit` / `push` stop at 10 for authorisation.

## Review chain conclusions

- **Round 1 architecture**: `deep-reasoner` + DeepSeek V4 Pro; the unbounded-wait list and the global lock-order proposal were released after user confirmation.
- **2 implementation, two passes**: main body landed -> the four plan A items -> the F-01..F-06 fixes.
- **Red tests first**: `RED-CONFIRMED` (the counter-example tests falsify first, then the implementation starts).
- **Pre-review r1**: `PRE-REVIEW: PASS WITH FIXES` (0 Critical / 0 High / 6 Medium + 6 Low).
- **Pre-review r2**: `PRE-REVIEW: PASS` (every round-1 finding verified closed, no new Critical / High / Medium).
- **Scan**: `INDEPENDENT SCAN: PASS`.
- **Final review**: `FINAL REVIEW: PASS`.

- Prefix rule: the archived text of pre-review r1 carries the `FINAL REVIEW` prefix (a mislabelled stage prefix); the archive keeps it verbatim and unrewritten, while this table records it under the stage's normative prefix `PRE-REVIEW` with an unchanged value.

## Six-stage scan transcript

Source `docs/validation/evidence/P40-T027B-LOCK-REAPER/06-scan/review.md` (archive SHA256 = `829edcc93c12b4d289d1170e9c9d80c70fe2a5e3be5c8439ad4c36a13585f569`); the fence transcribes it verbatim with zero rewriting, from `## Section A` through the conclusion line.

```sixstage
## Section A — 总体结论
本次只读扫描未发现 Medium 及以上缺陷；从 [internal/guard/process.go](internal/guard/process.go#L71-L118)、[internal/guard/process.go](internal/guard/process.go#L145-L241)、[internal/guard/process_windows.go](internal/guard/process_windows.go#L217-L257)、[internal/tickets/durable_cost_ledger.go](internal/tickets/durable_cost_ledger.go#L298-L355) 和 [internal/tickets/durable_cost_lock_shared.go](internal/tickets/durable_cost_lock_shared.go#L14-L35) 的代码形态看，修复整体遵循设计中“唯一 reaper + 有界 stop verification + 任务结束后再释放句柄/槽位 + 递增锁序”的收口方案，并且在生产代码里未发现残留的直接无界 Wait 或兜底的 TODO/Background 模式。

## Section B — 发现清单
INFO | n/a | 只读扫描未发现需要立即修复的 Medium 及以上缺陷；本轮审查结论为无新增风险记录 | none

## Section C — 逐条证据与推理
1. 唯一 reaper 与 waitDone 语义成立
- 证据： [internal/guard/process.go](internal/guard/process.go#L71-L118) 定义了 waitOnce、waitDone、waitErr，并在 [internal/guard/process.go](internal/guard/process.go#L94-L101) 中通过单次 Do 启动唯一 reaper；[internal/guard/process.go](internal/guard/process.go#L145-L241) 的 runManaged、Terminate 和 Reclaim 都通过 waitDone 等待单一终态，而不是各调用点独立调用 cmd.Wait。
- 推理：这直接消除了“并发 Wait 触发 panic”的根因；对应的红测 [internal/guard/process_reaper_test.go](internal/guard/process_reaper_test.go#L20-L146) 明确验证了同一 process 只产生一个 waitDone，并且 concurrent Run/Terminate 不会二次 Wait。
- 修复可行性：高；这是设计目的中的关键收口，已通过代码结构和测试契约对齐。

2. 有界 stop/verify 与 uncertain 收口成立
- 证据： [internal/guard/process.go](internal/guard/process.go#L165-L215) 中 runManaged 的 cancellation 分支在调用 process.platform.stop 后，使用 stopVerificationContext(grace) 对 waitDone 与 verifyGone 做独立超时控制；[internal/guard/process.go](internal/guard/process.go#L322-L366) 的 Terminate 也遵守同样规则，未让调用者的取消直接撤销 stop 验证窗口。
- 推理：这实现了“取消与超时同刻时，最终以明确的 uncertain 兜底，不再无界阻塞”的收口；对应回归测试 [internal/guard/process_reaper_test.go](internal/guard/process_reaper_test.go#L32-L92) 和 [internal/guard/process_reclaim_test.go](internal/guard/process_reclaim_test.go#L18-L96) 也都覆盖这些边界。
- 修复可行性：高；这是一组明确的 bounded proof，和设计文档的 §1 一致。

3. taskkill 脱离式执行与到界自证成立
- 证据： [internal/guard/process_windows.go](internal/guard/process_windows.go#L217-L257) 中 platformStopIdentity 使用 goroutine 执行 taskkill，并在 ctx.Done 或 timer 触发时返回 errStopCommandExpired，不尝试强行杀死已经脱离的命令；随后仍由 waitIdentityGone 做真实树证明。
- 推理：这样避免了“外部命令自己卡住导致整个进程路径挂住”的问题；测试 [internal/guard/process_stop_windows_test.go](internal/guard/process_stop_windows_test.go#L18-L91) 明确规定：超出 bound 时返回 detached/unknown，而不是直接失败。
- 修复可行性：高；与设计中的“taskkill 脱离 + 到界自证”完全一致。

4. 全局锁序 guard 与 ctx 取消支路成立
- 证据： [internal/tickets/durable_cost_lock_shared.go](internal/tickets/durable_cost_lock_shared.go#L14-L35) 定义了 ErrDurableLockReentrant，并通过 goroutine id 记录锁栈；[internal/tickets/durable_cost_ledger.go](internal/tickets/durable_cost_ledger.go#L312-L355) 中 runWithLock 先调用 durableLockOrder.check，再 acquireDurableLock，确保不能同 goroutine 以非递增顺序重入。
- 推理：这意味着在 same-scope 或 cross-scope 逆序获取时，至少一条路径会被拒绝，不再出现互相等待的死锁；对应的红测 [internal/tickets/durable_cost_lock_order_test.go](internal/tickets/durable_cost_lock_order_test.go#L17-L191) 和 [internal/tickets/durable_cost_lock_reentry_test.go](internal/tickets/durable_cost_lock_reentry_test.go#L11-L76) 对这条不变式做了明确锁定。
- 修复可行性：高；与设计中“budget 锁与 replay 锁绝不嵌套”形成一致性。

5. 生产代码未发现残留的直接无界 Wait 或兜底 Background/TODO
- 证据：本轮只读 grep 覆盖了 internal 下的关键模式；除测试和少量真实进程外，生产代码中没有发现直接的 cmd.Wait、Process.Wait、CombinedOutput 在主路径中的无界调用，也没有发现 internal 生产代码中的 context.TODO() 或对 context.Background() 的无界兜底等待。
- 推理：这意味着要点 ①、②、③ 的修复没有明显残留；当前改动保持在设计要求的边界之内。
- 修复可行性：高；如果未来新增代码把 Wait 放回主路径，现有测试会直接暴露。

## Section D — 可证伪性审计
这个 PASS 结论可以被具体证伪，最直接的反例包括：
- 删除 [internal/guard/process.go](internal/guard/process.go#L94-L101) 中的 waitOnce.Do 包装，会让同一进程再次触发并发 cmd.Wait，并直接翻转 [internal/guard/process_reaper_test.go](internal/guard/process_reaper_test.go#L90-L146) 里“concurrent wait 不 panic”的结论。
- 删除 [internal/tickets/durable_cost_lock_shared.go](internal/tickets/durable_cost_lock_shared.go#L29-L35) 中的 reentrancy check，会让 [internal/tickets/durable_cost_lock_reentry_test.go](internal/tickets/durable_cost_lock_reentry_test.go#L11-L34) 直接失效，恢复同 goroutine 重入锁的死锁形态。
- 把 [internal/guard/process_windows.go](internal/guard/process_windows.go#L224-L257) 中的 detached taskkill 调用改回同步 CombinedOutput，也会重现“任务壳等待外部命令卡住”的无界停机问题，导致 [internal/guard/process_stop_windows_test.go](internal/guard/process_stop_windows_test.go#L18-L91) 的 PASS 结论翻转。
- 另一个反例是：如果在 [internal/guard/process.go](internal/guard/process.go#L183-L214) 中去掉 stopVerificationContext 这一层 bounded verification，过程会重新退回到无界阻塞的老问题。

## Section E — 我无法验证的部分（明确列出）与什么会推翻 PASS
- 我无法在本轮只读扫描中代替真实的 Windows 进程树验证来证明 taskkill 在所有机器和所有进程树状态下都稳定收敛；这类行为最终仍依赖真实 OS 过程和跨平台压力测试。
- 我没有在本轮重新跑一遍全库 go test；我只能依据当前任务上下文中给出的“全量 go test -count=1 ./... 15 ok / 0 FAIL”和“tree 门禁 TOTAL_FAIL=0”来作为已知状态，而不是我自己重新执行的证明。
- 任何下列结果都会推翻当前 PASS：
  1. 真实 Windows 运行中，taskkill 在 bound 内仍不能收敛并留下存活子树；
  2. 同一进程的多次 Wait 能在真实执行中再次 panic；
  3. 某条锁序检查在真实并发下被绕过，导致可重复出现同 goroutine 重入或跨 goroutine 逆序死锁；
  4. 生产代码里被遗漏的直接 Wait/Background/TODO 在后续新增代码中复活。

INDEPENDENT SCAN: PASS
```

## F-item closure

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| F-01 | `abortStartup` leaves the bound-reaching handle without a reclaim path | Medium | `internal/guard/process.go:139-166` | Closed: the bound branch starts a background goroutine that calls `platform.close()` once `waitDone` closes; `platform.close` is idempotent in all three forms |
| F-02 | `Reclaim` has no direct test and no production caller | Medium | `internal/guard/process.go:366-385` | Closed: `TestReclaimReleasesHandleOnce` pins four assertions; the production wiring is registered as a later slice |
| F-03 | the CI Race leg covers adapters / chain only | Medium | `.github/workflows/ci.yml:84` | Registered: this slice does not change CI files; widening the package coverage is a later-slice candidate |
| F-04 | detached taskkill and the 10s wide window carry no test | Medium | `internal/guard/process_windows.go:219-253` | Closed: a nanosecond bound asserts `errStopCommandExpired`, `taskkill-detached` and `taskkill-process-tree`; the settled-action combination row is registered as not constructible |
| F-05 | the `ErrProcessTreeRemained` branch has no guard-level regression | Medium | `internal/guard/process.go:191-216` | Closed: a parent-exit child-residue case in `process_linux_test.go` asserts the evidence and `signal-term-process-group` |
| F-06 | the same-stack equal path and panic-pop carry no test | Medium | `internal/tickets/durable_cost_lock_shared.go:57-91` | Closed: three cases pin re-entry rejection, pop after panic, and serialisation of non-nested concurrency |
| F-07a | the static nesting assertion is not a CI criterion | Low | `01-arch-design/review.md` §2 | Registered: later-slice candidate |
| F-07b | the implementation is stricter than the design wording and the document was not written back | Low | `01-arch-design/review.md` §1 | Closed: annotations added at §1 line 5 and line 8, with no conflict against the original lines |
| F-08 | the `ctx.Done` and timer branches are redundant equivalents | Low | `internal/guard/process_windows.go:241` | Registered: kept, harmless |
| F-09 | nil-ctx defensive style is inconsistent | Low | `internal/guard/process.go:272,327` | Registered: deferred |
| F-10 | same-scope non-nested concurrency carries no test | Low | `internal/tickets/durable_cost_ledger.go:321-350` | Closed: folded into the F-06 cases |
| F-11 | the lockPath comparison is not canonical | Low | `internal/tickets/durable_cost_lock_shared.go:67` | Registered: later-slice candidate |
| F-12 | the relation between 10s + grace and the stopper outer 12s is undocumented | Low | `internal/console/stopper.go:15-16` | Registered: covered by a comment |
| Round-1 Low rule | the disposition rule for the six Low items | Low | pre-review r1 / r2 | Accepted: on condition that the F-07a and F-11 registrations appear in the 9 report or in a later ledger batch |

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
| batch | `deep-reasoner` | GLM-5.3-Flash (2 times) |

**The four items of plan A**

- **A-1 (detached taskkill)**: the external command runs in a detached goroutine with a buffered channel for reclaim, the bound never kills the command, and `CombinedOutput` is called exactly once.
- **A-2 (10s wide window plus self-proof)**: the synchronous wait window is capped at 10s; at the bound a `waitIdentityGone` self-proof runs first, a dead tree is judged stopped and only a live tree is judged uncertain.
- **A-3 (natural-exit close condition)**: the platform handle is closed and the launch slot released only when `waitDone` is closed and `verifyGone` is empty.
- **A-4 (no settle on stop)**: the stop path adds no settle window, so it cannot stack with the caller deadline.

**Other registered items**

- **The F-04 settled action is not constructible**: it needs a test seam, which is outside this slice's boundary, so it is registered as not constructible.
- **`-race` cannot run on this machine**: it returns `-race requires cgo`, so that leg is carried by CI.
- **the CI Race leg excludes guard / tickets**: this slice does not change `ci.yml`, so widening the package coverage is registered as a later-slice candidate.
- **production `Terminate` has no `Reclaim` wiring**: no caller takes Reclaim on the uncertain path of `agent_runner_pinned_cli.go`, so the handle is kept until process exit, which is bounded; the wiring is registered as a later slice.
- **goroutine id parsing degrades**: when the id cannot be parsed the guard stops rejecting re-entry, which lowers safety without ever rejecting wrongly.
- **the red-test helper `select{}` is ended by the runtime deadlock detector**: it reports a deadlock after about 276ms, so the new helper uses a `time.Sleep` long-lived loop instead.

## Gate readings

- Three-scope readings: tree `TOTAL_FAIL=0`, index `TOTAL_FAIL=0`; the ci scope runs authoritatively inside CI (`--scope=ci`) — local ci degenerates to `HEAD^1..HEAD`, see the ⑩ receipt.
- SELFTEST: PASS 697/697; the G6-1 scanned count is stated by the criterion output itself (reference form, avoids self-reference drift).

## Commit and CI

- Main commit `3fab6ba` (feat: enforce single process reaper and durable lock order (P40-T027B-LOCK-REAPER)), range `38a5cba..3fab6ba`; CI run `41164098703`: ubuntu leg green (Test/Race/Contract fixtures/Gates), windows leg Test step red (guard 600s timeout; root cause per DR-12).
- Fix commit `3af9c9a` (fix: remove redundant Process.Wait in reaper test for Go 1.22 compat); CI run `41165709859` **double green**: windows 3m51s / ubuntu 2m41s; both legs `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=3`, `TOTAL_FAIL=0`; `internal/guard` and `internal/tickets` both ok on both legs (windows 15 packages ok / ubuntu 17 packages ok, including the Race legs).
- Docs-only commits after this sentence no longer back-fill this section; readings are carried by the ⑩ receipt.

## Later-slice candidates

- F-07a: bring the static lock-order assertion into the CI criteria (`go/analysis` or an equivalent scripted assertion).
- F-11: canonicalise `lockPath` before comparing.
- Reclaim production wiring: call `Reclaim` on the uncertain path of `agent_runner_pinned_cli.go`.
- CI Race leg package widening: bring `internal/guard` and `internal/tickets` into `-race`.

## Known limits

- The gist of 6 Section E: taskkill convergence across machines depends on real OS processes and cross-platform stress testing, which a read-only scan cannot replace; likewise a real panic from repeated `Wait`, and a lock order bypassed under real concurrency, can only be falsified by real runs.
- The deviation items registered above (`-race` unavailable, the Race leg not covering guard / tickets, production Reclaim unwired, the settled action not constructible, goroutine id degradation) stay registered for the validity period of this slice's conclusion.

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/p40-lock-reaper.md	repo
docs/validation/p40-lock-reaper_EN.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/MANIFEST.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/SHA256SUMS.txt	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/00-lock-census.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/01-arch-design/review.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/05-prereview/review.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/06-scan/review.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/07-final/review.md	repo
docs/validation/evidence/P40-T027B-LOCK-REAPER/freedom-list.md	repo
internal/guard/process.go	repo
internal/guard/process_linux.go	repo
internal/guard/process_linux_test.go	repo
internal/guard/process_reaper_test.go	repo
internal/guard/process_reclaim_test.go	repo
internal/guard/process_stop_windows_test.go	repo
internal/guard/process_unsupported.go	repo
internal/guard/process_windows.go	repo
internal/tickets/durable_cost_ledger.go	repo
internal/tickets/durable_cost_ledger_test.go	repo
internal/tickets/durable_cost_lock_order_test.go	repo
internal/tickets/durable_cost_lock_reentry_test.go	repo
internal/tickets/durable_cost_lock_shared.go	repo
internal/tickets/durable_cost_lock_unix.go	repo
internal/tickets/durable_cost_lock_windows.go	repo
tools/gates/cjk-newwords.txt	repo
```

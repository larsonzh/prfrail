# ProofRail 验证报告 — P40-T027B-LOCK-REAPER

| 字段 | 内容 |
|---|---|
| 切片 | `P40-T027B-LOCK-REAPER` |
| 类型 | 实现类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `38a5cba` |
| 日期 | 2026-10-31 |

## 目标与边界

- **目标**：消除 C1 所指的两类结构性风险——`Run` / `Terminate` 的无界回收等待，以及缺少机械约束的全局锁序。
- **收口形态**：每进程唯一 reaper（`waitOnce` / `waitDone`）取代各调用点的直接 `Wait`；停止验证统一经 `stopVerificationContext(grace)` 与 `verifyGone` 两条通道判定，到界返回 `ErrTerminationUncertain` 且不释放平台句柄与 launch 槽；`runWithLock` 在取进程内锁与跨进程文件锁之前先过锁序守卫，同栈非递增重入直接拒绝。
- **范围**：`internal/guard/**` 与 `internal/tickets/durable_cost_*`。
- **边界**：不碰契约 / schema / fixtures / CI 文件；不改准则业务规则；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 审查链结论

- **① 架构**：`deep-reasoner` + DeepSeek V4 Pro；无界等待清单与全局锁序方案经用户确认后放行。
- **② 实现两轮**：主体落成 → 方案 A 四项收口 → F-01..F-06 修复。
- **③ 红测先行**：`RED-CONFIRMED`（反例测试先证伪，再进入实现）。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（0 Critical / 0 High / 6 Medium + 6 Low）。
- **⑤ r2**：`PRE-REVIEW: PASS`（第 1 轮发现全部核实闭合，无新增 Critical / High / Medium）。
- **⑥ 扫描**：`INDEPENDENT SCAN: PASS`。
- **⑦ 终审**：`FINAL REVIEW: PASS`。

- 前缀口径：⑤ r1 的归档原文以 `FINAL REVIEW` 前缀书写（环节前缀误标），归档时按原样保留、未作改写；本表按环节的规范前缀 `PRE-REVIEW` 记录，结论取值不变。

## 六段扫描转录

来源 `docs/validation/evidence/P40-T027B-LOCK-REAPER/06-scan/review.md`（归档文件 SHA256 = `829edcc93c12b4d289d1170e9c9d80c70fe2a5e3be5c8439ad4c36a13585f569`）；围栏内逐字转录、零改写，范围自 `## Section A` 起、至结论行止。

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

## F 项闭合

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| F-01 | `abortStartup` 到界句柄无回收路径 | Medium | `internal/guard/process.go:139-166` | 已闭合：到界支路先起后台 goroutine，等 `waitDone` 关闭后调 `platform.close()`；`platform.close` 三形态幂等 |
| F-02 | `Reclaim` 无直测、无生产调用方 | Medium | `internal/guard/process.go:366-385` | 已闭合：`TestReclaimReleasesHandleOnce` 钉住四条断言；生产接线登记为后续片 |
| F-03 | CI Race leg 仅覆盖 adapters / chain | Medium | `.github/workflows/ci.yml:84` | 已登记：本片不改 CI 文件，扩包覆盖列为后续片候选 |
| F-04 | 脱离式 taskkill 与 10s 宽窗零测试 | Medium | `internal/guard/process_windows.go:219-253` | 已闭合：纳秒 bound 断言 `errStopCommandExpired`、`taskkill-detached`、`taskkill-process-tree`；settled 动作组合行登记为不可构造 |
| F-05 | `ErrProcessTreeRemained` 支路无 guard 级回归 | Medium | `internal/guard/process.go:191-216` | 已闭合：`process_linux_test.go` 父退子残留用例，断言证据与 `signal-term-process-group` |
| F-06 | 同栈相等路径与 panic-pop 无测试 | Medium | `internal/tickets/durable_cost_lock_shared.go:57-91` | 已闭合：重入拒绝、panic 后出栈、非嵌套并发串行化三用例 |
| F-07a | 静态嵌套断言未入 CI 判据 | Low | `01-arch-design/review.md` §2 | 已登记：后续片候选 |
| F-07b | 实现语义比设计字面更严、文档未回写 | Low | `01-arch-design/review.md` §1 | 已闭合：§1 第 5 行与第 8 行补注，与原行语义无冲突 |
| F-08 | `ctx.Done` 与 timer 支路等价冗余 | Low | `internal/guard/process_windows.go:241` | 已登记：保留，无害 |
| F-09 | nil-ctx 防御风格不一致 | Low | `internal/guard/process.go:272,327` | 已登记：延后处理 |
| F-10 | 同 scope 非嵌套并发无测试 | Low | `internal/tickets/durable_cost_ledger.go:321-350` | 已闭合：并入 F-06 用例 |
| F-11 | lockPath 字符串比较非 canonical | Low | `internal/tickets/durable_cost_lock_shared.go:67` | 已登记：后续片候选 |
| F-12 | 10s + grace 与 stopper 外层 12s 的关系无说明 | Low | `internal/console/stopper.go:15-16` | 已登记：注释说明 |
| 第 1 轮 Low 口径 | 六项 Low 的处置口径 | Low | ⑤ r1 / r2 | 接受：前提是 F-07a 与 F-11 的登记出现在 ⑨ 报告或台账后续批次 |

## 偏差与登记

**模型绑定与载体**

| 环节 | 载体 | 模型 |
|---|---|---|
| ① 架构 | `deep-reasoner` | DeepSeek V4 Pro |
| ② 实现 | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ 测试 | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ 预审 | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ 扫描 | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ 终审 | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ 文档 | `prfrail-documenter` | DeepSeek V4.1 Flash |
| 批量 | `deep-reasoner` | GLM-5.3-Flash（2 次） |

**方案 A 四项**

- **A-1（taskkill 脱离）**：外部命令以脱离式 goroutine 执行，缓冲通道回收，到界不杀命令，`CombinedOutput` 恰调用一次。
- **A-2（10s 宽窗 + 自证）**：同步等待窗封顶 10s，到界先做 `waitIdentityGone` 自证；树已死判 stopped，未死才判 uncertain。
- **A-3（自然退出 close 条件）**：仅当 `waitDone` 已关且 `verifyGone` 为空时才关平台句柄、释放 launch 槽。
- **A-4（stop 不加 settle）**：停止路径不追加 settle 窗口，避免与调用方 deadline 叠加。

**其余登记项**

- **F-04 settled 动作不可构造**：需要测试缝，超出本片边界，登记为不可构造。
- **`-race` 本机不可跑**：实测返回 `-race requires cgo`，该腿由 CI 承担。
- **CI Race leg 不含 guard / tickets**：本片不改 `ci.yml`，扩包覆盖登记为后续片候选。
- **生产 `Terminate` 无 `Reclaim` 接线**：`agent_runner_pinned_cli.go` 的 uncertain 路径无人接 Reclaim，句柄保留到进程退出，属有界；接线登记为后续片。
- **goroutine id 解析失败退化**：解析不到时退化为不拒重入，安全性下降但不误拒。
- **红测 helper 的 `select{}` 会被运行时死锁检测结束**：约 276ms 即报死锁，故新 helper 改用 `time.Sleep` 长活替代。

## 门禁读数

- 三作用域读数：tree `TOTAL_FAIL=0`、index `TOTAL_FAIL=0`；ci 作用域由 CI 内 `--scope=ci` 权威执行（本地 ci 退化为 `HEAD^1..HEAD`，读数见 ⑩ 回执）。
- SELFTEST: PASS 697/697；G6-1 扫描计数以判据自身输出为准（引用式，避免自指漂移）。

## 提交与 CI

- ⑧b 提交后回填。

## 后续片候选登记

- F-07a：把静态锁序断言纳入 CI 判据（`go/analysis` 或等价的脚本断言）。
- F-11：`lockPath` 先做 canonical 化再比较。
- Reclaim 生产接线：在 `agent_runner_pinned_cli.go` 的 uncertain 路径接入 `Reclaim`。
- CI Race leg 扩包覆盖：把 `internal/guard` 与 `internal/tickets` 纳入 `-race`。

## 已知限制

- ⑥ Section E 的要旨：taskkill 跨机器收敛依赖真实 OS 进程与跨平台压力测试，只读扫描不能替代；多次 `Wait` 的真实 panic、锁序在真实并发下被绕过，同样只能由真实运行证伪。
- 上述偏差登记项（`-race` 不可跑、Race leg 未覆盖 guard / tickets、生产 Reclaim 未接线、settled 动作不可构造、goroutine id 退化）在本片结论的有效期内保持登记状态。

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

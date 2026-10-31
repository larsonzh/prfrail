# P40-T027B-LOCK-REAPER ⑥ 独立扫描归档（原样）

切片：`P40-T027B-LOCK-REAPER`。环节：⑥ 独立扫描。载体：`independent-reviewer`；模型：`MAI-Code-1.1-Flash (copilot)`。输出原样归档，内容未作改写。

---

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

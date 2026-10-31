数据来源：读取 `docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/review.md`、`docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md`、`docs/CONTRACTS.md`（188–283）、`internal/tickets/durable_cost_ledger.go`、`durable_cost_lock_{shared,unix,windows}.go`、`internal/guard/process.go`、`internal/guard/process_windows.go`、`internal/adapters/agent_runner_pinned_cli.go`、`internal/adapters/agent_runner_timeout.go`、`internal/adapters/agent_runner_replay_store.go`、`internal/adapters/agent_runner_process_registry.go`、`internal/chain/engine.go`、`internal/chain/agent_runner_launcher.go`、`internal/gates/executor.go`、`internal/release/probe.go`，及 grep 检索 `internal/` 下全部 `context.Background()/TODO()/\.Wait\(\)/CombinedOutput\(\)/sync\.(Mutex|RWMutex|Once)`。本方案未写入任何文件。

## 1 无界等待清单

已确认 `internal/` 生产代码（非 `_test.go`）中 `context.TODO()` 为 0 处。无界点如下：

| 位置（文件:行） | 函数 | 现状 | 替换方案（一句话 + 关键代码 ≤ 8 行） |
|---|---|---|---|
| `internal/guard/process.go:237` | `Terminate` | `_, _ = process.cmd.Process.Wait()` 直接无界 Wait（P37 C1 点名处） | 改为等待唯一 reaper 的 `waitDone`，有界；到界返回 uncertain 且不 `platform.close()`。<br>`select { case <-p.waitDone: case <-verifyCtx.Done(): terminateErr = errors.Join(terminateErr, ErrTerminationUncertain) }` |
| `internal/guard/process.go:189` | `runManaged`（ctx.Done 支路） | `platform.stop` 后 `<-waited` 等 reaper goroutine 交付，进程不死则取消方永久挂起 | 同样改 `<-p.waitDone`，以 `stopVerificationContext(grace)` 有界。<br>`select { case <-p.waitDone: case <-verifyCtx.Done(): stopErr = errors.Join(stopErr, ErrTerminationUncertain) }` |
| `internal/guard/process.go:132` | `runManaged` | reaper goroutine `process.cmd.Wait()` 无生命周期管理，可被多调用点二次 Wait | 上移为 `StartManaged` 后 `waitOnce.Do` 的唯一 reaper；其余入口只读 `waitDone/waitErr`，禁止任何直接 `cmd.Wait()` |
| `internal/guard/process.go:92,99` | `StartManaged`（attach/identity 失败清理） | `_ = cmd.Process.Kill(); _, _ = cmd.Process.Wait()` 无界回收 | 改用同一 `waitDone` + grace 派生短 ctx 有界等待；未退出则报告 `ErrTerminationUncertain` 且不 close Job 句柄（所有权移交后续 reclaim）。到界未回收的句柄由后台 reaper 等待线程在 waitDone 关闭后释放（见实现 abortStartup 到界支路），不依赖 Reclaim 接管 |
| `internal/guard/process_windows.go:177` | `platformStopIdentity` | `taskkill` 经 `command.CombinedOutput()` 无界等外部命令 | taskkill 以脱离式 goroutine 执行（缓冲通道回收、到界不 Kill 命令、CombinedOutput 恰调用一次），同步等待窗由调用方 bound 限定（StopProcessIdentity 侧 10s 封顶）；命令未完成则到界放行。到界后先有界自证（waitIdentityGone），树已死报 stopped（Actions 含 taskkill-settled-verified-gone），未死才 uncertain |
| `internal/tickets/durable_cost_ledger.go:121` | `OpenDurableCostLedger` | bootstrap 用 `context.Background()`；`durableLockBackoff` 遇非 nil ctx 永不取消 → 锁被僵持进程持有则 Open 无限等 | 签名加 ctx：`OpenDurableCostLedger(ctx, config)`；`runWithLock` 已支持 ctx 取消的 `acquireDurableLock`，调用方传入有界 ctx（建议 10s open 上界） |
| `internal/tickets/durable_cost_lock_shared.go:15-22` | `durableLockBackoff` | `ctx == nil` 支路 `time.Sleep` 死循环 | 删除 nil 支路（改硬错误 `ErrInvalidCostLedger`），强制所有调用点携带非 nil ctx |
| `internal/adapters/agent_runner_pinned_cli.go:575` | `stopSpawnedProcess` | `WithoutCancel(Background())` 已用 `grace+1s` 有界，但与调用方 deadline 脱钩 | 保留有界形态作为“settle 窗口”规范：deadline=`now+grace+settle`，`settle` 复用 `agentRunnerSettleWait` 的 500ms 上限 10s |
| `internal/chain/engine.go:218` | `runTask` | `Reviewer.Review(ctx, …)` 在 run 路径上同步等待；若外层以 run 互斥包住、评审经同锁控制 API 应答即自锁 | 禁止任何引擎/状态锁跨越 Review；Review 返回后先 `engine.boundary(ctx)` 处理取消/暂停再 Publish |
| `internal/adapters/agent_runner_timeout.go` | `RunWithWatchdog` | 已实现“取消获胜 + 20ms settle 窗口”，无变更 | 作为 §A 收敛规则的现成范本复用 |
| `internal/gates/executor.go`、`internal/release/probe.go:94` | 门禁/探测执行 | 已分别经 `guard.RunManaged(Hook.Timeout)` 与 `exec.CommandContext(ctx)` 有界 | 无变更（登记为已合规） |

**收敛规则（取消与超时同刻）**：取消获胜。形态统一为——定时器触发后先 `if err := ctx.Err(); err != nil { return err }` 再宣布超时（timeout.go 已如此）；stop 验证 select 中 `case <-ctx.Done()` 排在 `case <-verifyCtx.Done()` 之前，且**到界支路一律不 close platform 句柄、不释放 launch 槽**：所有权只在该 launch 的 `waitDone` 已关且 `verifyGone==nil` 时由显式 `releaseLaunch` 释放。

## 2 全局锁序

锁点清单（进程内 `sync.Mutex` / 跨进程文件锁分列；均为读取确认）：

| 锁名 | 类型 | 持有方 | 现获取顺序 |
|---|---|---|---|
| registry 锁 | 进程内 `sync.Mutex`（`agent_runner_process_registry.go:81`） | launcher（register/lookup/live/advance） | 单取单放，短临界区 |
| budget 进程内锁 | `sync.Mutex`（`durable_cost_ledger.go:81`）+ 纯内存 `cost.go:203` RWMutex | Reserve/Settle/RequireOutstanding | mu → 文件锁 → fn（fn 无回调取锁） |
| budget 跨进程锁 | flock/`LockFileEx` 非阻塞轮询（`durable_cost_lock_{unix,windows}.go`） | 各进程账本实例 | 同上，ctx 可取消 |
| replay 锁 | `sync.Mutex`（`agent_runner_replay_store.go:166`） | RecordRequest/RecordCompletion/State | 单取单放 |
| 一次性文件槽（跨进程） | O_EXCL slot / identity mirror / R、C、intent、identity、closure 的 no-replace 发布 | launcher/dispatcher/publisher | 取即用即放，未与任何锁同持 |
| 记录级写者锁 | capability/completion/enforcement/event 各自 `sync.Mutex`（`agent_runner_*.go:44-62`） | 各记录 writer | 单取单放，叶子 |
| 门禁输出锁 | `gates/executor.go:63`、`result.go:22` | hook runner | 单取单放，叶子 |
| 快照 store 锁 | `snapshot/store.go:21` RWMutex | object store | 单取单放，叶子 |

**绝对获取顺序**（一个序关系链）：

$$registry(T0) \lt budget(T1) \lt replay(T2) \lt 一次性文件槽(T3)$$

理由：T0 是纯内存叶子，天然最内；T1 是唯一跨进程重锁，必须 ctx 可取消且持锁期间零回调；T2 是发布串行点，持锁期间禁止调 Settle（防止与「budget 读→replay 写」路径交叉）；T3 是 O_EXCL/no-replace 一次性断言，绝不与 T1/T2 同持。发布链「intent→Settle→C→closure」每步取放分离、不嵌套，故合规。**核心不变式：任意路径同时持有的锁层必须严格单调递增；若未来必须嵌套，只允许按 T0<T1<T2 方向，禁止反向，禁止跨越 T3。**

**机械断言方式（≤3 条）**：
1. 静态 CI 脚本：`go/analysis` 单包 checker（或 `grep`+脚本）断言「`mu.Lock()` 临界区内不出现对 `ledger.run*`/`store.mu`/`registry.mu` 加锁方法的调用；`runWithLock` 的 fn 内不出现任何加锁语句」，违规即 CI fail。
2. 运行时锁序追踪器：build tag `locktrace` 注入 LockTracer，记录每 goroutine 获取序列；新增 `TestConcurrentDispatchPublishStop`（并发 dispatch+publish+stop 压测，`-race` 下跑）断言序列严格单调、锁临界区内无 `select`/I/O。
3. 跨进程测试：`TestLockWaitHonoursCancel`（另一进程持 budget 锁僵持，等待方 ctx 取消须及时返回）；槽文件 O_EXCL 碰撞后必须走有界重读收敛（复用 `loadConvergedRecordsLocked` 口径），断言等待都有 ctx 上界。

## 3 关键接口签名

```go
// internal/guard/process.go
type ManagedProcess struct {
    // ...现有字段不变...
    waitOnce sync.Once
    waitDone chan struct{} // 唯一 reaper 结束时 close，绝不二次 Wait
    waitErr  error         // waitDone 关闭后只读
}
func (p *ManagedProcess) waitResult() <-chan struct{}         // waitOnce.Do 启动唯一 reaper
func (p *ManagedProcess) Run(ctx context.Context, grace time.Duration, startedAt time.Time) (ProcessResult, error)      // 签名不变，内部改走 waitDone
func (p *ManagedProcess) Terminate(ctx context.Context, grace time.Duration) (TerminationEvidence, error)              // 有界；uncertain 不 close 不释放
// internal/tickets/durable_cost_ledger.go
func OpenDurableCostLedger(ctx context.Context, config DurableCostLedgerConfig) (*DurableCostLedger, error)            // 新增 ctx
```

“到界 uncertain 且不释放写者所有权”形态：

```go
select {
case <-p.waitDone:                          // 树已证死，可继续回收
case <-ctx.Done():                          // 取消获胜，即使 verifyCtx 同刻到界
    terminateErr = errors.Join(terminateErr, ctx.Err())
case <-verifyCtx.Done():
    terminateErr = errors.Join(terminateErr, ErrTerminationUncertain)
}
// 关键：terminateErr != nil 时严禁 p.platform.close()/releaseLaunch(slot)；
// 句柄与槽保留，由后续唯一 reclaim 路径在 waitDone 关闭后显式释放。
```

## 4 风险与反例

1. **uncertain 后进程才死 → 句柄滞留**：返回 uncertain 保留所有权，但无人再消费 `waitDone`。对策：`waitErr` 缓存于 `waitDone` 后，新增 `Reclaim(ctx)` 幂等入口（`waitDone` 已关 且 `verifyGone==nil` 才释放槽）；唯一 reaper 保证只回收一次。
2. **budget 锁同栈重入饿死**：publisher 恢复路径若在持有 budget 锁的 fn 内再次 `Settle` → 同一进程不同 fd 的 flock 互斥自锁。对策：`runWithLock` 记录持有 goroutine id，重入返回专用 `ErrDurableLockReentrant`，测试钉住。
3. **未来顺序反转**：若把重确认的 budget 读移入 replay mu 临界区，与「Settle 回执校验读 R」路径交叉成 T2→T1 死锁。对策：第 2 节静态断言禁止嵌套；确需嵌套时只允许 T1→T2 单向，代码审查清单明列。
4. **waitOnce 迁移期二次 Wait panic**：`exec.Cmd.Wait` 并发调用会 panic，process.go:92/99/132/237 四个旧点未同步迁移即触发。对策：迁移必须在同一批内完成（先上 `waitDone` 再删全部直接 Wait），用 P38 的 `TestConcurrentRunAndTerminateUseOneWaiter`（Wait 计数=1）作门禁。

以下为回归影响面枚举全文（只读完成，未修改任何文件；可直接归档）。

---

# P48 / DR13-FIX-D3 · GLM-2 回归影响面枚举

基线：工作区 e742224（D1 所有权门 + D2 Stat 复核已落地）。落点 `internal/adapters/agent_runner_starting_marker.go:92-100` `(*startingMarkerClaim).release()`，现序 unlock(:96) → Close(:97) → `os.Remove`(:99)，新序改为 remove → unlock → Close（函数体内顺序，签名与返回值语义不变）。

## 1. 全部调用点

### 1.1 `claim.release()` 直接调用点（全仓仅 2 处生产代码）

| 调用点 | 触达路径 | 新序行为影响 |
|---|---|---|
| `agent_runner_process_registry.go:86`（`closeStartingMarker` :82-87，释放后置 `:87` 置空保证幂等） | 唯一汇聚点，见 1.2 | 中间态可见性变化（W1 收窄）；返回值 bool 所有调用方均忽略（`:86`/pinned_cli.go:212 均 `_ =`），无漂移 |
| `agent_runner_pinned_cli.go:212`（onSpawn 内 register 失败回滚 :206-213，直接持 claim 非 handle） | marker 已发布至 `startingMarkerPath`、锁仍持 | 新序下 remove 发生在锁内：扫描者要么见路径缺失（ErrNotExist no-op :343-350），要么见 Windows delete-pending 中间态（见 §3 风险 R-N1） |

### 1.2 `closeStartingMarker` 调用点（3 处）

| 调用点 | 场景 | 新序行为影响 |
|---|---|---|
| `agent_runner_pinned_cli.go:272` | Start 成功路径（advance(running) 后、返回前；handle 仍在 registry，D1 门在此窗口兜底） | 等价；扫描者命中中间态时由 D1(`:380-393`) 或 D2(`:367-377`) 收敛为 `(false,nil)` |
| `agent_runner_pinned_cli.go:689`（`releaseLaunch` :679-691 内，先 `:681` releaseLaunchSlot、`:687` closeLogs 后执行） | Stop 双腿 :292/:301、`agent_runner_run.go:118` default 分支、`agent_runner_postflight.go:719` default 分支、测试 `agent_runner_pinned_cli_test.go:589` | 等价；launch 语义后置不变（slot/registry 清理与 marker 顺序解耦） |
| `agent_runner_pinned_cli.go:702`（`discardStartingLaunch` :697-704） | Start 失败腿 `:229/:237/:246/:268`（guard 出错、processID 不可证、mirror 写失败、advance 失败） | 等价；四条腿均保证无活进程，marker 删除时序无观察者依赖 |

### 1.3 `releaseLaunch` / `discardStartingLaunch` 上游（行为等价性同上，逐点列出备查）

- `releaseLaunch`：`agent_runner_pinned_cli.go:292`（stopIdentity 失败分支——STRATEGIC-REVIEW GAP-73 已登记残留，非本片范围）、`:301`（Stop 正常收尾）、`agent_runner_run.go:118`、`agent_runner_postflight.go:719`、测试 `agent_runner_pinned_cli_test.go:589`。
- `discardStartingLaunch`：`agent_runner_pinned_cli.go:229/237/246/268`。
- `tools/b4-e2e/runmap.go:46,260` 仅注释引用（未导出不可调）。

**结论**：全部影响收敛到 `:86` 单点的函数体内顺序；各调用点后置语义（返回值、幂等、slot/registry 清理）不变。

## 2. 同形态「unlock→remove」序列比对

| 位置 | 现序 | 是否 W1 | 枚举意见（裁定权在 ①/②） |
|---|---|---|---|
| `agent_runner_starting_marker.go:68-72` `writeAgentRunnerStartingMarker` discard 闭包（unlock:69→Close:70→Remove:71），触发点 :74/:78/:83 | 解锁后删 **temp 名** | 否——扫描者只盯发布路径（`startingMarkerPath` :41-44），temp 名不可见；rename 失败腿 :83 时发布路径从未存在 | **不修改**（不属 W1，改序无收益且超本片范围） |
| `agent_runner_pinned_cli.go:530-572` `clearAgentRunnerStartingMarker` | 锁内 remove：`:567` Remove → `:568` unlock → `:569` Close | 它本身就是「remove 先于 unlock」的正确形态，是 D3 的仓库内先例 | **不修改**（已正确；② 若抽象公共序须保持此形态，勿被新序改写反向污染） |
| `agent_runner_pinned_cli.go:770-786` `releaseLaunchSlot` / `:462-497` `reclaimLaunchSlot` | rename-away 原子协议，slot 文件无锁 | 否（无锁可解） | **不修改** |
| `agent_runner_starting_marker_windows.go:165-171` 与 `:133-136`（createLockedStartingMarkerFile 探测/锁失败清理）、unix :90-97 | Close→Remove（锁经 Close 隐式释放），temp 名 | 否 | **不修改**（形态为 close→remove 且目标不可见） |
| `agent_runner_pinned_cli.go:589-597` `removeAgentRunnerProcessIdentityMirror` | 无锁 remove | 否 | **不修改** |
| `RecoverStarting` 内 unlock 点 :368（D2 腿）、:392（D1 门）、:430（决策完成） | 只解锁关闭，清理委托 `clearAgentRunnerStartingMarker`（F-1 重开重锁协议） | 否 | **不修改** |

**总结论**：全仓「先解锁后删已发布路径」的 W1 形态仅 `release()`(:92-100) 一处；本片只动它，其余全部保持。

## 3. Windows 专属影响（`agent_runner_starting_marker_windows.go`）

- **锁语义不被破坏**：LockFileEx/UnlockFileEx 绑定句柄字节区间而非目录项；remove 后 unlock(:71-81)/Close 均正常。unix flock(:43-52) 同理（unlink 不影响 open file description 上的锁）。
- **share-delete 前提已备**：claim 句柄由 `createLockedStartingMarkerFile`(:125-203，CreateFile 带 ShareRead|Write|Delete) 创建 → 锁内 `os.Remove` 与 `writeAgentRunnerStartingMarker:82` 的「锁内 rename」同一机制；与 `openStartingMarkerForClear`(:102-123，同样 ShareDelete，注释 :98-100 自证「Go 的 os.OpenFile 只申请读写共享」) 不冲突。
- **R-N1（新中间态，建议 ① 回归显式钉住）**：本机 GOROOT 证实 Go 的 `os.Remove` = 裸 `DeleteFileW`（`src/os/file_windows.go` `func Remove`），stdlib 无 ERROR_DELETE_PENDING 映射（已 grep 零命中）⇒ Windows 上 remove 返回成功后、`claim.file.Close()` 前，名字处于 delete-pending。此窗口内 `RecoverStarting` 的 `os.OpenFile`（`agent_runner_pinned_cli.go:337`）失败码非 `os.ErrNotExist` → 走 :343-350 硬错误腿而非 no-op，并发 Start 可能概率性收到 open 失败。窗口为 remove→Close 三语句级（µs 量级）；现序在同一窗口的观察结果是「可打开但锁已被释放」（恰为 W1），新序将其换成 delete-pending 错误形态——两者都不是干净 no-op，D2 只覆盖「打开成功后 inode 失联」一腿。与 `design.md` §B/§3「必须钉住本工具链实际行为」一致，建议红测增加 OpenFile 对 delete-pending 文件的实际 errno 探针或以并发压力回归观察（注：② 仅改函数体，不得顺手做 errno 映射）。
- **扫描者持柄场景**：`os.OpenFile` 柄无 delete-share → 新序 remove 失败（sharing violation）→ 返回 false、路径保留、unlock/Close 仍无条件执行（A(2) 收敛，design §B）；与现序终态等价，无死锁、无锁泄漏。
- `reclaimStartingProcess`(:87-97) 不触句柄，无影响。

## 4. 测试影响面（回归重点观察清单）

`internal/adapters` 包内无测试直接构造 `startingMarkerClaim` 或断言 marker 时序（glm-1 已证）→ ③ 须全新增红测。既有测试按触达方式分组：

| 测试 | 位置 | 触达 | 新序敏感度 |
|---|---|---|---|
| ConcurrentStartsHaveOneWinner | `agent_runner_pinned_cli_test.go:857`（slot 断言 :924） | 8 并发 Start，唯一跨越 release 间隙的压力场（败者 RecoverStarting 可踩中间态） | **最高**：R-N1 若成立必在此概率性暴露，建议 `-count≥20` |
| StopsThroughPersistedIdentity | :579（:589 `t.Cleanup` 直调 `releaseLaunch`——链上唯一显式测试调用点） | 直调 | 后置断言等价，观察即可 |
| StartsAndStops / UnprovenIdentityStopsTheSpawn / StopsItsOwnProcessWhenTheMirrorIsUnreadable | :224 / :615 / :651（:691 slot 断言） | Start→Stop、discard 腿 | 低（后置断言） |
| RefusesASecondLiveLaunch / RefusesASecondInstance… / StopGrace… / StopNeverReaches… / HoldsNoSlotDuringThePrecheck / VersionMismatch… / RefusesAnotherExecutable / RefusesToStopAForeignLaunch / RefusesStopWhenTheSlotIsUnreadable / RefusesAnUnreadableIdentityMirror | :278/:331/:380/:431/:491/:535/:558/:705/:755/:940 | 经 Stop/discard 间接过链 | 低 |
| postflight Run 链 | `agent_runner_postflight_test.go:680` | 经 default 分支 releaseLaunch | 低 |
| 恢复路径（不经 claim.release，同扫描语义建议全跑） | `spawn_recovery_test.go:138/:171/:221`（:251 helper 无锁写 marker）；windows :75/:124、unix `spawn_recovery_live_host_unix_test.go:113/:155` | RecoverStarting 直测 | 不受顺序影响；windows_test.go:40 `holdStartingMarkerLock` / unix_test.go:45 `holdUnixStartingMarkerLock` 是包内唯一外部锁探针先例，③ 可复用 |
| console 跨层并发 | `internal/console/runtime_agent_concurrent_e2e_test.go:106` | 经 dispatcher→Start | 低概率间接触达，纳入三作用域即可 |

## 5. 跨包影响

- `internal/chain`：仅经接口 `AgentRunnerLauncher.StartAgentRunnerProcess`（`internal/chain/agent_runner_launcher.go:11-16`）触达；release 链全部未导出 → 零直接触达。
- `internal/console`：窄接口仅 `StopAgentRunnerProcess`（`internal/console/runtime_agent.go:848`，弃局缓解路径）→ 经其内部间接过 release 链，后置语义不变、无时序断言。
- `internal/repair`：grep 零匹配 → 零触达。
- `tools/b4-e2e`：仅导出 API（`runmap.go:46` 注释自证无法调用未导出 releaseLaunch）。

**结论**：release 链唯一生产落点为 `agent_runner_process_registry.go:86` 单点，全部调用方在 `internal/adapters` 包内且忽略返回值；跨包零语义影响，回归收敛为包内测试 + 双平台 CI 腿。

# ① 定性结论（原样归档）

> 载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。transcript 偏移 119545。输入：`tmp/DR13-FIX/batch1-slot-exclusivity.md`、`tmp/DR13-FIX/batch2-test-rig.md`。

## A 定性结论

**DIAGNOSIS: B** —— 真实排他窗口漏洞：同进程（同 launcher 实例）并发 Start 时，`RecoverStarting` 的"marker 锁空闲 ⇒ 前宿主已死"证明失效，且 unix 分支会**杀死活启动的进程**并夺走其 slot，使两个 Start 先后成功。排除 A 的理由：两次成功之间没有任何合法 slot 释放点被触发（无 Stop 调用；stub `run` 模式睡 60s 且 launcher 从不 reap；首胜者未失败），slot 在首胜者存活期间被删只能经由 reclaim 误判——而 reclaim 把活进程当崩溃孤儿杀掉，这是实现缺陷，不是"断言把跨时间成功数当同一时刻存活数"的装置问题。`:901-902` 的断言语义恰是协议不变式本身（slot 原子独占 ⇒ 第二个 Start 在首胜者存活期间必须失败），不可放松。

## B 证据链

| # | file:line | 事实 | 指向 |
|---|---|---|---|
| 1 | test:857-882 | 8 contender 同进程、同 launcher，`close(start)` 广播放行，无错峰 | 前提 |
| 2 | test:46-49 | stub `run` 模式 `time.Sleep(time.Minute)`，不会自行退出 | 排除 A 的"合法二次成功" |
| 3 | pinned_cli.go:124-129 | `registry.live()` 入口门在 claim 之前，8 者皆可同时通过（register 在 OnSpawn :209 才发生） | 排他靠 slot 兜底 |
| 4 | :151→:707-725 | slot 由 `O_EXCL` 原子创建，正常路径下败者 :689-692 得 Conflict | 排他机制本身 |
| 5 | :147-150, :363, :370, :377 | 每个 contender claim 前都跑 `RecoverStarting`；marker 可解码+pid 存活 ⇒ `usable && alive` | 触发面 |
| 6 | unix.go:60-70 | `reclaimStartingProcess(alive=true)` **对活进程执行 `StopProcessIdentity` 并返回 (true,nil)** | 杀活根因（ubuntu） |
| 7 | windows.go:87-92 | 同函数 alive=true 返回 (false,nil)（fail closed，不杀） | 解释为何 windows 腿不红 |
| 8 | marker.go:105-115 | `claim.release()` = unlock → close → `os.Remove`，锁释放与路径消失之间可被扫描者观察 | W1 窗口 |
| 9 | marker.go:52-88 | marker 原子发布（temp 先 flock 后 rename），但锁的存活证明绑定**文件 inode + 进程**，不区分同进程 goroutine | 前提 |
| 10 | :406-412 | fail-closed 分支（`!usable`/`reused`/`!processReclaimed`）无条件 `reclaimLaunchSlot` + 清 marker | 夺槽路径（windows 同根异症） |
| 11 | :413-417 | 完整回收删除 slot（:430-466 owner 比对通过）+ 删镜像 | 夺槽路径 |
| 12 | :684, :692-700 | claim 时 `refuseLiveMirrorProcess`：镜像缺失或已死即放行 | 镜像防线在"进程被杀后"失守 |
| 13 | :242 | 两个赢家先后覆盖写镜像；:31-36 注释明确"第二个启动会剥离活进程身份"正是要防的损害 | 后果 |
| 14 | :265-276 | 首胜者返回前**不复查自身进程存活**（advance 后直接 success） | 首胜者"成功了但进程已死" |
| 15 | test:894-896 | 败者哨兵白名单只有 Conflict/Duplicate —— 观察到的 `got 2 (losers 6)` 与"恰两个赢家、其余全哨兵"逐字吻合 | 证据吻合 |

## C 交错时序（ubuntu，W1 为主）

**W1（不依赖任何特殊 flock 内核语义）**：
1. A 经 `createLaunchSlot` O_EXCL 成功（:707-725），slot 落盘。
2. A `guard.StartManaged` 重执行测试二进制 spawn stub（CI 上数十~数百 ms）。
3. A 的 OnSpawn 发布 marker（marker.go:52-88：temp 先 flock 后 rename，锁持续持有）并 register（:209）。
4. B（第二个跑完 precheck 者）在 A"marker 已发布、尚未 closeStartingMarker"的毫秒窗口内 :341 `os.OpenFile` 成功，随后在 :345 tryLock 之前被调度器抢占。
5. A :272 `closeStartingMarker` → release() unlock→close→`os.Remove`；:275 `claimed=false`；:276 返回成功（赢家 #1，stub 仍活）。
6. B 恢复：:345 tryLock 锁住**已被 unlink 的孤儿 inode**成功（A 已 unlock，inode 因 B 的 FD 仍存活）→ :363 decode 成功 → :377 `ProcessAlive(A.stub)=true` → :391 unix.go:66 `StopProcessIdentity` **杀死 A 的存活 stub**，返回 (true,nil)。
7. B :413-417：`reclaimLaunchSlot(A.launchID)` 删除 A 的 slot（owner 匹配，:457/:465）+ 删镜像 + 清 marker。
8. B :151 claim：slot 缺失 → :684 `refuseLiveMirrorProcess`：镜像已删（或 A 刚补写但其进程已死 → alive=false）→ 放行 → O_EXCL 成功。B spawn、写 marker/mirror、advance、返回成功（赢家 #2）。
9. 其余 6 人 claim 时 slot 新鲜 → :689-692 Conflict（或扫描见 B 持锁 → (false,nil) 后同样 Conflict）——全部命中哨兵。
⇒ winners=2, losers=6，与两次 CI 输出（`41253798647`/`41257990185`）逐字吻合。

**W2（放大器，内核相关）**：B 的 :345 tryLock 因 A 持锁失败 → B :361 关闭自己的 FD；按 man flock(2)"进程关闭任一指向该文件的 fd 即释放该进程对该文件全部 flock"（部分内核版本行为，见 §G），A 的锁被 B 的关闭动作释放；此后任一后来者 C 的 tryLock 无需抢占巧合即成功，后续同 W1 步骤 6-9。W2 使"锁空闲+进程存活"状态贯穿 A 的整个收尾窗口（毫秒级），显著放大命中率。

**Windows 为何不红**：windows.go:87-92 alive→(false,nil) → :406-412 fail-closed 只删 slot、**保留镜像**；夺槽者 claim 经 :684 见 A 的**活**镜像 → Conflict。ubuntu 的"收养并杀"反把活进程杀掉，令镜像防线失守。注意 windows 同竞态留下"活启动被夺槽"的隐性破坏（其后 `stopIdentity`/`refuseForeignSlot` 将拒绝停止该启动）。

## D 修法

**D1（主修，确定性、双平台、不涉协议）**：`RecoverStarting` 在 :363 decode 成功之后、:370 之前插入同进程所有权门：

```go
if decodeErr == nil {
    if handle, found := launcher.registry.lookup(marker.LaunchID); found && handle != nil {
        unlockStartingMarker(file)
        _ = file.Close()
        return false, nil
    }
}
```

不变式：**本 launcher 持有 launchID 的活句柄 ⇒ 该 launchID 的 marker 永不回收**。W1/W2 状态下 A 的句柄必在 registry（OnSpawn :209 注册，Stop 前不释放），故门确定性拦截；跨进程崩溃恢复不受影响（崩溃后 registry 随进程消亡）。返回 (false,nil) 与既有"锁被持有/无 marker"同签名（:349-362 的 F-5 注释歧义不被扩大），无需改 CONTRACTS/Schema。锁序安全：registry.mu 与 marker flock 无反向嵌套。

**D2（建议，关 W1 跨进程残留）**：:345 tryLock 成功后立即 `os.Stat(path)`，`ErrNotExist` ⇒ 锁住的是已 unlink 的孤儿 inode，unlock/close 后返回 (false,nil)。一行成本。

**D3（可选加固，平台敏感）**：`startingMarkerClaim.release`（marker.go:105-115）改为**先 `os.Remove(path)` 再 unlock**，配合 D2 全闭 W1（路径消失严格先于锁释放）。风险：Windows remove-while-open 依赖 claim 句柄 share-delete（已设），但触及成功/失败/退休热路径，建议留作后续片。

**影响面**：`TestRecoverStartingRefusesLiveHostUnix`（spawn_recovery_live_host_unix_test.go:118）、`TestRecoverStartingIdempotent`/`TokenMismatchFailsClosed`（spawn_recovery_test.go:134+）、`TestRecoverStartingCleansCrashedHostWindows`（spawn_recovery_live_host_windows_test.go:75）均 registry 为空，门不触发；同文件串行用例（:278/:331/:431/:491）不并发；各 `registry.len()` 断言（:273/:553/:600/:636/:688/:904/:972）不受影响（门不增删句柄）；真实崩溃孤儿仍被收养停止（unix.go:60-70 未改）。

## E 回归与反例

**E1（证伪测试，旧代码下双平台稳定红）**——新增 `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns`（spawn_recovery_test.go 共享区）：
1. `recoveryLauncherFor` + `startRecoveryHelperProcess` 活 helper；
2. 包内调 `launcher.claimLaunchSlot("launch-request-one")` 造 slot；
3. 包内 `launcher.registry.register(&agentRunnerProcessHandle{launchID:..., requestID:..., phase:agentRunnerProcessStarting, identity:live})`；
4. `writeStartingMarker` 写同 launchID marker（**不加锁** ⇒ 精确构造"锁空闲+进程存活+本进程持句柄"致命状态）；
5. 断言：`reclaimed==false`、helper 仍活、marker 仍在、slot 仍在且 owner 仍为该 launchID。
旧代码：unix 杀活 helper（unix.go:66）→ 红；windows fail-closed 删 slot（:406-412）→ 红。
**E2**：保留 `TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner` 原样作并发回归（修后 winners==1 为确定性结果），建议 CI `-count=5`/`-race` 强化；其 :914-920 "mirror 进程仍活"正是被杀 stub 的探针。
**E3**：D2 的孤儿 inode 路径无外部 seam 可确定性构造（需在 open 与 flock 之间注入暂停），暂以 D1 门的确定性替代证伪，见 §G。

## F 对 P45 的输入

1. P45（零模型三任务 E2E、并发装配、ubuntu 腿）必须新增"同一 launcher/runRoot 并发 Start"场景断言：**恰一个赢家**；败者全部命中 Conflict/Duplicate 哨兵；closeout 时 `registry.len()==1`。
2. 必测"赢家进程在装配结束仍存活"（ProcessAlive），并核对 slot owner token == 赢家 launchID、mirror 与 slot 一致——这三项正是本竞态会破坏的状态。
3. 禁止用 10s freshness 墙钟作排他断言依据；用完成序/状态断言。
4. 若 P45 有 Stop 与 Start 并发，断言"两方必有一方失败"，禁止静默夺槽/杀活进程（本竞态的受损面）。
5. P45 门禁证据应含并发装配用例的包级 `ok` 行与计数，并显式在 ubuntu 腿运行。

## G 未确证项

- W2 的"关闭任一 fd 释放进程对该文件全部 flock"语义：man flock(2) 明文记载，但现代内核 `locks_remove_flock` 是否仍如此未在本地实证（本片禁止写文件，无法跑探针）；W1 独立成立，不依赖 W2。
- 两次 CI 红究竟经 W1 还是 W2 无法从日志区分（终局一致）。
- D2 的 TOCTOU 残量与 D3 的 Windows remove-while-open 全量回归成本未验证。
- Windows 腿"未红"的结构性解释（alive→fail-closed→活镜像拦截）基于代码推演，未在 windows CI 实证。
- fail-closed/corrupt 分支（:406-412/:448-458）在"镜像尚未写入"瞬间对**跨进程**活启动的夺槽残留，D1 不覆盖（仅覆盖同进程），D2/D3 部分覆盖。

§A 闭集判定：DIAGNOSIS: B
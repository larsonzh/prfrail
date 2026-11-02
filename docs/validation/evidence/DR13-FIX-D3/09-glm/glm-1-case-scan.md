## P48 DR13-FIX-D3 · GLM-1 用例检索报告（只读扫描，锚点 e742224）

### 0. 落点确认
`agent_runner_starting_marker.go:92-102` `startingMarkerClaim.release()` 现序：`unlockStartingMarker → file.Close → os.Remove(claim.path)`。调用链：`releaseLaunch`(pinned_cli.go:679→:689) / `discardStartingLaunch`(pinned_cli.go:697→:702) → `closeStartingMarker`(process_registry.go:82→:86) → `claim.release()`。**包内无任何测试直接构造 `startingMarkerClaim` 或直接调用 `writeAgentRunnerStartingMarker`/`release()`（已 grep 确认）→ 现状红测缺口成立，③ 必须新写红测。**

### 1. 直接触达释放链的既有测试（`agent_runner_pinned_cli_test.go`，全部经 Start/Stop 间接经过 `release()`）

| file:line | 测试 | 被测行为一句话 | 新序影响 |
|---|---|---|---|
| :224 | StartsAndStops | Start→Stop 全链，断言进程停止、registry 清空 | 间接：走 release()，不断言 marker 时序 |
| :278 | RefusesASecondLiveLaunch | 竞争第二 Start 被拒，Stop 后槽位可复用 | 间接 |
| :331 | RefusesASecondInstanceWhileTheMirrorIsLive | mirror 拒第二实例；丢失 slot 仍拒；Stop 后 stale mirror 不阻塞 | 间接 |
| :380 | StopGraceFollowsTheLaunch | stopGrace 随 launch，末尾 Stop 收尾 | 间接（弱） |
| :431 | StopNeverReachesAnotherLaunch | 死 launch 的 Stop 幂等且不动他人 slot（:441 slot 归属断言） | 间接 |
| :491 | HoldsNoSlotDuringThePrecheck | 预检期不持 slot（:517/:524 stat 断言） | 间接（弱） |
| :579 | StopsThroughPersistedIdentity | **t.Cleanup 直接调 `launcher.releaseLaunch`(:589)**——改动链最直接的既有调用点 | 间接 |
| :615 | UnprovenIdentityStopsTheSpawn | 身份发布失败→`discardStartingLaunch` 回滚→release()；断言 registry 空且 slot 释放可再启动 | 间接 |
| :651 | StopsItsOwnProcessWhenTheMirrorIsUnreadable | mirror 不可读仍停自己的进程 | 间接 |
| :705 | RefusesToStopAForeignLaunch | 拒外来 stop；:744/:780 slot 删除场景 | 间接（弱） |
| :755 | RefusesStopWhenTheSlotIsUnreadable | slot 不可读/无主全 fail-closed；:802/:822 操作 slot | 间接 |
| :857 | ConcurrentStartsHaveOneWinner | 8 并发恰好一胜；退役后断言 slot `ErrNotExist`(:924)——**注意这是 releaseLaunchSlot 断言，非 marker** | 间接 |
| :940 | RefusesAnUnreadableIdentityMirror | mirror 不可读→启动失败回滚 discardStartingLaunch | 间接 |
| :1017 | RefusesUnusableConfig | 配置拒绝，无 spawn | 无影响 |
| :1050 | StopWithoutIdentityFailsClosed | 无身份 Stop 拒绝（:1071 断言的是 mirror 存在） | 间接（弱） |
| :535 | VersionMismatchNeverSpawns | 版本预检失败，marker 尚未写入（marker 写入在 pinned_cli.go:196） | 无影响 |
| :558 | RefusesAnotherExecutable | 请求校验拒绝，无 spawn | 无影响 |

### 2. RecoverStarting / clearAgentRunnerStartingMarker 侧（恢复路径，**不经 `claim.release()`**，恢复侧清理走 `clearAgentRunnerStartingMarker`(pinned_cli.go:530)——本就是 remove-while-locked 先例）

| file:line | 测试 | 被测行为一句话 | 新序影响 |
|---|---|---|---|
| spawn_recovery_test.go:138 | TestRecoverStartingIdempotent | 死 host marker 两次回收仅一次；:156 断言 marker 已删 | 无影响 |
| :171 | TestRecoverStartingTokenMismatchFailsClosed | token 不匹配不杀、:201 marker 仍清理、:204 mirror 保留 | 无影响 |
| :221 | TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns | 本进程仍持有 live launch（handle 无 claim、marker 无锁）时拒绝回收，:280 marker 不动 | 无影响（钉 DR13 所有权不变量，D3 不改判定） |
| spawn_recovery_live_host_windows_test.go:75 | TestRecoverStartingCleansCrashedHostWindows | 崩溃 host 立即清 marker+slot（:97/:100），slot 立即可 `O_EXCL` 复用(:103) | 无影响 |
| :124 | TestRecoverStartingRefusesLiveHost | LockFileEx 持锁期间拒回收，解锁+死亡后回收 | 无影响 |
| spawn_recovery_live_host_unix_test.go:118 | TestRecoverStartingRefusesLiveHostUnix | flock 同上（Linux 半区） | 无影响 |
| :160 | TestRecoverStartingAdoptsOrphanUnix | 自由锁 marker 命名活孤儿→收养+停止，marker/slot 释放（:192/:195） | 无影响 |
| spawn_recovery_test.go:45 | TestSpawnRecoveryHelper | re-exec 辅助子进程 | 无影响 |

补充：所有经 `StartAgentRunnerProcess` 的测试都会在 Start 内先跑 `RecoverStarting`(pinned_cli.go:148)，对空 run root 为 no-op 扫描——不受新序影响。平台无关平测：`TestRecoverStarting*` 在两平台均编译运行；Windows 专测仅 :75/:124，Unix 专测仅 :118/:160。相邻但不同链：`agent_runner_process_registry_test.go:66`（registry 簿记）、`agent_runner_replay_dispatcher_test.go:563-883` 的 Reclaim 系列（dispatcher Reclaimer 缝，不经 marker）——均无影响。

**③ 可复用手法**：`holdStartingMarkerLock`(windows_test.go:40) / `holdUnixStartingMarkerLock`(unix_test.go:45) 的「测试进程用第二个独立句柄抢锁」手法，是包内唯一能在 release() 过程中做外部观察的确定性锁探针先例。

### 3. 结构性测试先例
- `internal/adapters` 包内：**无** `go/parser`/`go/ast` 断言类测试。
- 仓内先例：`internal/release/network_test.go:16` `TestProductionCodeHasNoDirectNetworkImports`——`filepath.WalkDir(cmd, internal)` → `parser.ParseFile(token.NewFileSet(), path, nil, 0)` → 断言 `file.Imports` 禁入（net/net/*、os/exec 白名单 map）+ `ast.Inspect` 扫 `ast.BasicLit` 字符串字面量禁词（curl/wget/git clone…），违规 `t.Errorf` 带 path。

### 4. `createLockedStartingMarkerFile`（`agent_runner_starting_marker_windows.go:122-172`）供 ② 参考
- 完整共享模式：`CreateFile` access=`startingMarkerGenericRead|startingMarkerGenericWrite`；**share=`startingMarkerShareRead|startingMarkerShareWrite|startingMarkerShareDelete`（0x1|0x2|0x4）**；disposition=`startingMarkerCreateNew`(1)；attrs=`startingMarkerFileAttributeNormal`；模板 0；EXIST/ALREADY_EXISTS 重试（上限 `startingMarkerTempFileAttemptLimit`=32，:36）。
- 注释原句(:122-124)：*"createLockedStartingMarkerFile creates a temp file that is already locked, with a share mode that stays renameable: the share mode allows the delete that a rename performs, so the marker can be renamed into place while the lock stays held."* 另有 F-3 注释(:160-162)：*(false,nil) 是拒绝不是成功，否则交出无锁 marker 会击穿 live-host 规则*。
- 对 ② 的含义：claim 自身句柄已带 **FILE_SHARE_DELETE** ⇒ 新序「先 os.Remove 后 unlock/Close」在 Windows 上可行（名字先离命名空间，数据随 Close 归零）；Unix unlink-while-locked 天然可行。W1（unlock 后 marker 仍在）正是新序要闭合的窗口。

### 5. ② 实现注意点（扫描中发现的额外事实）
`writeAgentRunnerStartingMarker` 的 `discard` 闭包(starting_marker.go:69-71) 同为 unlock→Close→Remove 旧序，但只作用于**未发布 temp 名**（rename 前失败路径），不构成已发布 marker 的 W1；是否同序修改由 ①/② 裁定，若改不影响任何现有测试（上表均无断言触及）。

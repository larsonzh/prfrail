# ③ 红测基线（原样归档）

> 载体：runSubagent(agentName="prfrail-tester", model="DeepSeek V4.1 Flash (deepseek)")。transcript 偏移 119767。本文件证明新增反例在**未修复代码**上确定性红（先红后绿）。

## 1) 新增测试与使用的 helper

**测试**：`TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` —— `internal/adapters/spawn_recovery_test.go` **L221–L284**（唯一改动文件，+84 行；`git diff --stat` 已核验，其余 `docs/...` 变更为工作区既有状态，非本轮产生）。

**使用的真实 helper / 符号（全部存在，无改名）**：

| 用途 | 真实名字 | 位置 |
|---|---|---|
| launcher + runRoot | `recoveryLauncherFor(t, runRoot)` → `NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), replayStoreForRunRoot(...))` | `spawn_recovery_test.go:129` |
| 真实存活子进程 + identity | `startRecoveryHelperProcess(t)`（子进程即 `TestSpawnRecoveryHelper`，`PROOFRAIL_SPAWN_RECOVERY_HELPER`），已自带 `t.Cleanup(stop)` | 同上 `:83` |
| starting marker（**不加锁**） | `writeStartingMarker(t, runRoot, testStartingMarker{...})` —— 就是包内既有「原子写、无锁」的测试写法，正好构造「锁空闲」 | 同上 `:53` |
| launch slot | `launcher.claimLaunchSlot(launchID)`（`createLaunchSlot` 是私有实现，由它转调；走生产抢占路径） | `agent_runner_pinned_cli.go:675` |
| 句柄注册 | `&agentRunnerProcessHandle{launchID, requestID, phase: agentRunnerProcessStarting, identity, grace, slotToken}` + `launcher.registry.register(handle)` | `agent_runner_process_registry.go:48/119` |
| 存活判定 | `guard.ProcessAlive` | 既有用例同款 |
| owner 判定 | 读 `managed-process.launch.slot` 内容，`strings.Contains(owner, launchID)`（真实格式 `launchID\ntoken\n`） | `createLaunchSlot` |

证据构造严格对应「锁空闲 + 进程存活 + 本进程持句柄」三重态；无 `t.Parallel()`，无固定 sleep 同步（子进程自身为 helper 生命期）。

## 2) 三条命令输出

**① `gofmt -l internal/adapters`** → 无输出（`gofmt-exit=0`）
**② `go vet ./internal/adapters/`** → 无输出（`vet-exit=0`）

**③ 定向测试（未修复代码，本机 Windows/Go 1.27）**
```
--- FAIL: TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns (0.13s)
    spawn_recovery_test.go:281: recovery removed the starting marker of launch launch-request-one, which is still in flight in this process: GetFileAttributesEx C:\Users\...\Temp\TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns1769924313\001\managed-process.starting.json: The system cannot find the file specified.
FAIL
FAIL	github.com/larsonzh/prfrail/internal/adapters	1.109s
```
**④ 全包回归**
```
--- FAIL: TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns (0.07s)
    spawn_recovery_test.go:281: ...
FAIL
FAIL	github.com/larsonzh/prfrail/internal/adapters	40.950s
```
过滤 `^(ok|FAIL|---)` 与 `^    ---` 后仅一条失败 ⇒ **既有用例（含 4 个既存 recovery 用例）零打扰**，整包按常跑完。

## 3) 失败原因原文（证明钉住的正是本竞态）

失败点 **L281 = 「starting marker 必须仍存在」断言**。断点前三条已通过（`reclaimed == false`、`guard.ProcessAlive == true`），说明走到了 `RecoverStarting` 的回收分支：

- Windows 锚点：`reclaimStartingProcess(alive=true)` → `(false, nil)`（`agent_runner_starting_marker_windows.go:87`）
- ⇒ `!usable || reused || !processReclaimed` 为真 ⇒ **先 `launcher.reclaimLaunchSlot(marker.LaunchID)`（slot owner 正是 `launch-request-one` ⇒ rename-away 后删除），再 `clearAgentRunnerStartingMarker`（删除 marker），返回 `false`**
- 删除顺序在代码中是先 slot 后 marker，故 marker 断言失败时 slot 已被夺走（L281 先行 Fatalf 未再打印 slot 断言）

即：**锁空闲被当成「前宿主已死」**，于是同进程仍持有 `launchID` 活句柄的 launching 状态被回收——marker 被删、slot 被夺、不变量被破坏。在 Linux 上同一状态下一行代码走 `StopProcessIdentity` 杀掉活子进程并返回 `reclaimed == true`，测试将在断言 1/2 变红（`agent_runner_starting_marker_unix.go:60`），与 CI `got 2 (losers 6)` 同一根因。

## 4) 与 ① 给的名字的偏差

- `writeStartingMarker` / `startRecoveryHelperProcess` / `recoveryLauncherFor` / `claimLaunchSlot` / `registry.register` / `agentRunnerProcessHandle` / `guard.ProcessAlive`：**均已存在，未改名**。
- 唯一取舍：marker 用**测试 helper** `writeStartingMarker`（无锁）而非生产 `writeAgentRunnerStartingMarker`（返回带锁的 `startingMarkerClaim`，会与「锁空闲」要求冲突）。这与本片要证的「B 眼里的锁空闲」状态一致，不降低证伪力。
- `createLaunchSlot` 为私有且需先判 freshness，故经 `claimLaunchSlot` 进入（真实生产路径）。

**覆盖矩阵**

| 编号 | 测试名 | 类型 | 对应测试点 |
|---|---|---|---|
| DR13-OWN-1 | `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` | 竞态/所有权（锁空闲+存活+在册三重态） | 「在册 launch ⇒ 不得回收」四条断言（reclaimed=false / 进程存活 / marker 存留 / slot owner 保持） |
| — | 不可测项 | — | Unix 专有杀进程动作本机不可执行（Windows 锚点返回 false）；同一构造在 Linux 上会在断言 1/2 触发；端到端双 Start 调度难确定性复现，已用确定性三重态替代 |

TESTS: RED-CONFIRMED
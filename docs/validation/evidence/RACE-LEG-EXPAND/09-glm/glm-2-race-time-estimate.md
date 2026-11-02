# GLM-2 — 各包 Race 耗时评估与 flake 风险扫描（批量辅助 2/2 · 逐字归档）

> 载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`。任务：只读分析。以下为代理汇报原文，未改写。

---

分析完成（全程只读，未写任何文件）。以下为汇报：

# P49 RACE-LEG-EXPAND 只读分析报告

## A. Race 耗时证据

全部证据来自 `docs/validation/` 归档 CI 日志中 `Go ubuntu-latest Race` 腿的 `ok <pkg>` 行（11 个采样点，6 份日志）：

| 包 | 有据 Race 耗时 | 出处（文件:行） |
|---|---|---|
| adapters | **71.454–91.454s**（11 样本） | directive-history/drfix/archive/CI-push1.txt:33 (71.454s)；t028/archive/ci-log.txt:943 (83.702s)、ci-log3.txt:943 (84.329s)、ci-log4.txt:281 (84.219s)；haiku-bump/archive/ci-ubuntu.log:281 (84.056s)、wb-ubuntu.log:282 (84.167s)；prereview-guidance/archive/ci-38368822132/Go_ubuntu-latest.log:281 (84.277s)、ci-38369408308/…:282 (91.454s)；DR13-FIX/08-ci/ci-evidence.txt:12 (89.396s)；P44…/08-ci/ci-evidence.txt:17 (89.085s)；P45…/08-ci/ci-evidence.txt:12 (89.514s) |
| chain | **1.561–2.521s**（11 样本） | 同上各文件相邻行（如 ci-log.txt:944 =1.703s、CI-push1.txt:34 =1.572s、P45 ci-evidence.txt:13 =2.521s） |
| console / gates / guard / tickets | **无任何 Race 腿记录**（历史 `ok` 行全部来自非 Race CI 腿或本地日志） | 非 Race 基线见下 |

四新包非 Race 基线（用于估算）：
- console **24.524–29.343s**（t028/ci-log3.txt:251、P44 08-ci/ci-evidence.txt:21）
- gates **0.829s**（CI, DR13-FIX/08-ci/ci-evidence.txt:17）～ **6.713s**（本地, t028/go-test.txt:9）
- guard **2.106s**（CI, ci-38369408308/Go_ubuntu-latest.log:267）～ **7.901s**（本地, assembly-min/v3.txt:10）
- tickets **0.082s**（CI, t028/ci-log4.txt:1003）～ **2.442s**（本地, assembly-min/v4.txt:15）

**四新包 Race 估计与理由**（静态：race 放大按 1.5–2.5×，子进程按 exec 次数另计）：
- **console：40–75s（主导项）**。13 文件/103 函数中仅 e2e 一个慢测：`TestZeroModelConcurrentStartOneWinner` 至少 10 次 re-exec race 版测试二进制（8 contender 各跑 version check（CheckTimeout=20s/次）+ winner run + late-conflict 探测 + relaunch run），race 二进制初始化更慢；非 race 单次 exec ≈2.5–3s（总 24.5–29.3s）。
- **guard：8–30s**。20 函数中 4+ 个 re-exec 子进程测试（process_test、linux×2、spawn_ownership、reclaim），含 100–500ms 时序窗口。
- **gates：3–20s**。多数为单进程文件型测试，但 helper 为两级 re-exec 链（helper→grandchild）。
- **tickets：5–15s**。多数纯逻辑，3 个 re-exec helper 测试。
- **四新包合计 ≈ 55–140s**；Race 腿总时长估计 ≈ 85+2+140 ≈ 230s（上限约 4 分钟）。

## B. flake 风险扫描

说明：Race 腿为 Ubuntu ⇒ `*windows_test.go`（guard/process_stop_windows_test.go、process_spawn_ownership_windows_test.go）不参与；`*linux_test.go` 参与。

### B1 父测试固定 sleep（父进程直接睡/轮询步长）
- gates/executor_test.go:60 — 5ms（父等待 ready 文件轮询步长，上界 2s，L48）
- guard/process_test.go:46 — 10ms（父等 child.pid 轮询步长，上界 5s，L34）
- guard/process_test.go:171 — **300ms 固定暂停**（让 runManaged 起跑后再 cancel；grace=200ms，L139）
- guard/process_spawn_ownership_test.go:160 — 10ms（waitForPath 轮询步长，上界 10s）
- guard/process_spawn_ownership_test.go:208 — 20ms（TestOnSpawnFailureRollsBack 轮询步长，上界 10s）
- guard/process_spawn_ownership_test.go:229 — **500ms 固定暂停**（spawnOwnershipBlockedWindow，TestStartManagedResumeIsSeparated）
- tickets/durable_cost_lock_order_test.go:194 — 20ms（父等 ready 轮询步长，上界 10s，L189）
- （windows-only 的 stop_windows:103 20ms 不参与 Ubuntu 腿）

子进程 helper 体内（父测不等其睡满，被父 kill/收敛）：console e2e:43=1min（winner 存活期）；gates executor_test.go:102/106=1h 循环、:119=5ms；guard process_linux_test.go:31/141=30s、:70=10ms；process_reclaim_test.go:22=1h；process_spawn_ownership_test.go:57=30s；process_test.go:224=1h；tickets ledger:244=5ms、:662=2ms、lock_order:156=30s。

### B2 短超时窗口（≤10s，Ubuntu 腿参与）
- gates/executor_test.go:45 — Hook.Timeout **5s**；:48 — ready 轮询 **2s**（race 二进制两级 re-exec 冷启动逼近此界 ⇒ 首要 flake 源）；helper 体内 :117 另有 2s
- guard/process_test.go:34 — 5s；:51 — Terminate ctx **5s**；:90 — **5s**；:184 — grace+**10s** select 界
- guard/process_reaper_test.go:81 — **5s**（Terminate 阻塞界）；:152 — **10s**（并发 Run/Terminate 收敛界）
- guard/process_linux_test.go:61 — 5s（终止证明轮询）
- guard/process_spawn_ownership_test.go:112/124/140/155/189/196 — spawnOwnershipSettleBound **10s**；process_reclaim_test.go:76 — **10s**
- tickets/durable_cost_ledger_test.go:698 — **1500ms**（deadline helper 屏障）；durable_cost_lock_order_test.go:115 — **5s**；:189 — **10s**；:201 — **200ms**（锁等待 ctx）；:225 — **5s**；:237 — **20ms**
- console：无 ≤10s 窗口；e2e CheckTimeout=**20s**（超 10s 阈，但对单次 exec 是实质风险，列此备注）
- （ledger:236 的 60s 轮询 > 10s 阈，充裕，不列）

### B3 毫秒级时序断言
- console/runtime_agent_concurrent_e2e_test.go — `TestZeroModelConcurrentStartOneWinner`（8 并发 goroutine 唯一赢家）
- gates/executor_test.go — `TestGuardExecutorCleansDescendantsAfterParentExit`（父 kill 时后代须存活 ⇒ linux outcome=resource-limited）
- guard/process_test.go — `TestRunManagedCancellationVerifiesStopWithinBoundedContext`（300ms 起跑 vs 200ms grace）；process_reaper_test.go — `TestTerminationDeadlineReturnsUncertain`（elapsed ≤ deadline+2s，:85）、`TestConcurrentRunAndTerminateUseOneWaiter`（100ms 并发）；process_spawn_ownership_test.go — `TestStartManagedResumeIsSeparated`（500ms 窗口）
- tickets/durable_cost_lock_order_test.go — `TestCancelWinsOverTimeout`（20ms）、`TestLockWaitHonoursCancel`（200ms）；durable_cost_ledger_test.go — `TestP39DurableHelperDeadlineReserve`（1500ms）

### B4 testing.Short 统计
六包合计 **0**（console/gates/guard/tickets/adapters/chain 均 0）——Race 腿无 Short 逃生通道，所有测试都会全量跑。

## 结论

- **新包 Race 合计估计 ≈ 55–140s**，Race 腿总时长估计 ≈ 3.5–4 分钟（adapters 85s 主导不变）。
- **flake 风险 Top3**：
  1. **console e2e**：约 10 次 race 二进制 re-exec + 20s/次 CheckTimeout；冷缓存 2 核 runner 上 exec 放大最狠，且无任何父侧轮询兜底（纯协议收敛），首跑最可能超时/抖动。
  2. **gates `TestGuardExecutorCleansDescendantsAfterParentExit`**：2s ready 轮询 vs 两级 race re-exec 冷启动，race 下首次有逼近上限的现实风险（B1/B2 双命中）。
  3. **guard/tickets 的毫秒窗口组**：guard 300ms/500ms 固定暂停 + reaper 的 deadline+2s 断言；tickets 200ms/20ms 窗口与 1500ms helper 屏障——race 放大 6–8× 后可能首次暴露，但均有 5–10s 外层兜底，属低概率。
- **CI 首跑建议**：① 首跑 Race 腿单独跑一遍并保留完整日志归档（作为 P49 证据基线，六个 `ok <pkg>` 行全采）；② 若 console e2e 在 20s CheckTimeout 处抖动，考虑只给该 launcher 配置项留放宽余量（不要改断言语义）；③ gates 的 2s 轮询上界如首跑失败，优先把 2s 提为常量并加大（不动 5ms 步长）；④ 全六包无 `testing.Short`，不要指望 `-short` 降载；⑤ Race 腿仍限 Ubuntu，windows-only 测试文件无需纳入评估。

**风险声明**：四新包 Race 耗时为静态估计（未运行任何测试）；父/子 sleep 归类基于 env 分支与 re-exec 结构人工判定，其中 guard `process_test.go` 两测为同函数双角色（env 开关区分父/子），已按父侧归类。

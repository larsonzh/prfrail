# CONSOLE-MULTITASK-ASSEMBLY 验证报告

- 定位：v39 层 2 第 3 项（P50）；console 多任务生产装配扩展。
- 准则：v1.60；起跑锚点 `8bd8224`；日期 2026-11-03。
- 实现类：`[SLICE]`（实现类）；规模档 L；台账行 `P50`。

## 1. 切片概要

- **目标**：解除 console 装配层的单任务限制，使同一 chain 的多个 task 在既有 A6 排他契约下串行装配与驱动。
- **范围**：生产装配层（`internal/console/runtime_agent.go`、`runtime_supervised.go`）、adapters 新增注册表（`agent_runner_intent.go`）、离线 harness 与测试改写；`internal/chain` 与 launcher / mirror / slot 实现零改动。
- **裁定**：方案 A（用户 2026-11-03 裁定）——CLI 接线后置，本片限生产装配层 + 离线 harness + 测试改写；② 按 ②a–d 小步串行派单。
- **风险定级**：L 档核心装配重构；未降级 `[ROUTINE]`，未推 gitee。
- **不改边界**：A2/A4/A5 既有语义、历史验证报告、`.github/workflows` 与门禁脚本。

## 2. 协议先行

- ① 判定：不改 wire / Schema / fixtures（镜像、槽位、starting-marker 均为 store-local，不参与 R/C 的 recordHash）；必须修订 CONTRACTS A6 段，明确「镜像按 runRoot 唯一」是**时刻级**排他，取代 P45 D9 的单任务判读。
- ②a 落地：CN `docs/CONTRACTS.md` A6 段末新增一段（+1/−0），EN `docs/CONTRACTS_EN.md` 同步（+1/−0）；两文件保留 BOM 与 LF。
- 契约原文（CN，逐字）：

> **多任务串行装配（P50）**：runRoot 级身份镜像与启动槽是「**同一时刻至多一个在跑 launch**」的排他点，不构成 run 生命周期内的一次性唯一。同一 `runId` 的一个 chain 可含多个 task（每 task 一个 isolated-workspace code step）；console 装配层为每个 `(taskId, stepId, attempt)` 绑定各自已持久化的 `agent-runner-request`（各 `requestId` 互异，R/C/终局/工件按 requestId 隔离），并由 chain engine **串行**调度：前一 task 的 launch 必须已证明停机（槽位释放、镜像所指进程死亡或镜像已由新 launch 原子轮换）后，后一 task 才允许发布 R 与 launch-intent；镜像轮换=以新 launch 身份原子覆盖旧镜像，与既有「镜像命名死亡进程时为陈旧、可被替换」语义一致。同一 runRoot 内禁止并发在跑的 launch；并发 Start 败者依旧不得创建日志工件、不得覆盖身份镜像。事件日志仍是该 runRoot 唯一 `events/state-events.jsonl` 全序源；各 task 六枚冻结事实按各自终局路由事件的位置化证据（A5 `evidence[2..7]`）归属，不得跨 task 对账。

- 实现对应：`agentRunnerAssemblyPlan` / `supervisedAssemblyPlan` 逐 task 校验「每 task 一 step」；`agentRunnerAssemblyRecords` 强制 `requestId` 互异且 attempt 恒为 1；`agentRunnerBindTasks` 为每 task 构造独立 workspace、replay store 与终局发布器；驱动循环在下一 `engine.Run` 前调用 `proveTaskStopped`；adapters 的 slot / mirror 零改动即满足「并发 Start 败者不建工件、不覆盖镜像」。

## 3. 红绿翻转

- ③ 红基线（未改生产代码）：`internal/console/runtime_multitask_red_test.go` 三条用例全部 PASS——其断言形态即「当前实现拒绝多任务」；`gofmt -l` 空输出、`go vet ./internal/console/` 干净、全包 `ok console 8.495s`。证据 `03-tests/logs/03-tdd-red-run.txt`、`03-tests/red-baseline.md`。
- ② 翻转：三条用例改写为接受形并全绿——`TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape`、`TestSupervisedAssemblyAcceptsMultiTaskDefinitionShape`、`TestMultiTaskDefinitionShapeIsOtherwiseValid`。
- 单任务门文案清零（机械读数）：`requires exactly one task` 在两生产文件中 **0** 处；新文案为 `requires at least one task`（`runtime_agent.go`、`runtime_supervised.go`），step 形状门按设计保留。证据 `04-impl/logs/04-refusal-strings.txt`。
- **回归腿**：接受形用例常驻断言 2-task 定义可接受，任何把任务数改回「恰一」的回归立即变红；`TestMultiTaskDefinitionShapeIsOtherwiseValid` 保留逐 task 形状腿。
- 生产入口级证据：`TestAgentRunnerDriverTwoTaskSerialAssembly` 直调 `ExecuteAgentRunnerRun`；`TestAgentRunnerDriverStopsTheLaunchOfTheTaskThatJustFailed` 钉死 H1；`TestAgentRunnerAssemblyRefusesNonInitialAttempt` 钉死 M1。

## 4. 实现与整改

| 步 | 载体 | 落点 | 结果 |
|---|---|---|---|
| ②a | `prfrail-implementer` | `docs/CONTRACTS{,_EN}.md` | 各 +1/−0；CN 与设计草案逐字一致 |
| ②b | `prfrail-implementer` | `internal/console/runtime_agent.go` | 装配计划、记录注册表、逐 task 绑定、串行驱动；`return fail(` 由 7 处增至 8 处 |
| ②c | `prfrail-implementer` | `internal/console/runtime_supervised.go` | `supervisedAssemblyPlan`、逐 task 状态映射、停放与归还相位 |
| ②d | `prfrail-implementer` | `internal/adapters/agent_runner_intent.go` 与测试 | `AgentRunnerStepIntentRegistry`（+90/−0）、4 个新测试文件、2 个既有测试改写 |

- ⑤ r1 四项整改与闭环：
  - **H1（High）**：非 awaiting 分支按上一 task 错配 `requestId`，可致非首 task 的孤儿 launch 不被止损 ⇒ 改为观测优先解析（`agentRunnerDrivingRequestID` + `currentRequestID`）。
  - **H2（High）**：生产多任务驱动零行为覆盖 ⇒ 补 `runtime_multitask_driver_test.go`，三夹具直调生产入口，并覆盖停机证明三态（已停止 / 仍存活 / 报告缺失）。
  - **M1（Medium）**：不校验 attempt ⇒ `agentRunnerAssemblyRecords` 增加 `Attempt != 1` 拒绝，拒绝先于 marker 与引擎。
  - **M2（Medium）**：supervised 五事实位置序未钉死 ⇒ 以生产同源函数重推并按位比较。
- 新 seam（全部 console-local 且非导出）：`testStub`、`agentRunnerLaunchReceiptReader` / `agentRunnerLaunchStopper`、`agentRunnerDispatchWaiter`，以及 dispatch / admission / intent / postflight 四路 mux；生产载体复用既有 adapters 类型，`internal/chain` 零改动。
- 偏差登记：`AgentRunnerStepIntentRegistry` 已实现并单测，生产路径仍逐记录构造 preparer（⑦ Info 1，后片接用或删除其一）。

## 5. 审查链

| 环节 | 载体 | 判定 | 额度 |
|---|---|---|---|
| ⑤ r1 | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS WITH FIXES`（2 High + 2 Medium） | 1/1 |
| ⑤ r2 | 同上 | `PRE-REVIEW r2: PASS`（`SEVERITY=H0/M0/L0/INFO2`） | 1+1 用尽 |
| ⑥ | `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS`（Section B 无发现） | 1/3 |
| ⑦ | `independent-reviewer` + `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS`（`SEVERITY=INFO2`） | 1/1 |

- ⑥ 自跑：`go test ./internal/console -run 'MultiTask|AgentRunnerAssembly|SupervisedAssembly|MultiTaskDefinitionShape|Assembly' -count=1` 返回 `ok console 3.541s`。
- ⑦ 两项 Info 已登记：① intent registry 生产未接；② 本地 Go 1.27 与 CI 1.22 存在环境差 ⇒ ⑨ 以 CI Go 1.22 双腿为最终放行证据。
- ⑦ 放行条件两条（非阻断）：⑨ 须采集 CI 双腿；⑧ 须落实 ARCHITECTURE / DEV_PLAN 回写（见 §10）。

- ⑥ 结论逐字转录如下：

```sixstage
## Section A
This review was a read-only forensic check of the live P50 multi-task serial assembly in [docs/CONTRACTS.md](docs/CONTRACTS.md#L255-L289), [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L86-L200), [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L332-L524), and [internal/adapters/agent_runner_intent.go](internal/adapters/agent_runner_intent.go#L18-L83). I did not write or modify any files; the review package was taken directly from the on-disk working tree and checked against the contract and focused tests.

## Section B
(no findings)

The effective contract for P50 is preserved in the current tree: one `runRoot` remains the serial launch barrier, each task/step/attempt is bound by its own request record, and the next task may launch only after the previous task is proven stopped. I did not find a material violation of the A6 semantics in the live implementation.

## Section C
Evidence supporting the pass result:

- The contract explicitly says the same `runId` may contain multiple tasks, but the run-root identity mirror and launch slot remain the exclusivity point: “at most one running launch at any instant” and “the next launch must be proven stopped before the next R/launch-intent can be published.” See [docs/CONTRACTS.md](docs/CONTRACTS.md#L279-L289).
- The registry enforces per-task binding and rejects cross-task drift, duplicate step bindings, duplicate request IDs, and cross-run mixing. See [internal/adapters/agent_runner_intent.go](internal/adapters/agent_runner_intent.go#L27-L83).
- The serial engine barrier is explicit in the driver: `proveTaskStopped` calls `managedProcessesGone`, which fails closed on missing pid reports or still-alive PIDs, so the next task cannot start while the previous one is not proven stopped. See [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L912-L926) and [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L853-L910).
- The acceptance tests cover the exact P50 contract shape: two-task definitions are accepted, retired single-task refusal wording is rejected, and the serial harness is enforced in the real driver path. See [internal/console/runtime_multitask_red_test.go](internal/console/runtime_multitask_red_test.go#L124-L209) and [internal/console/runtime_multitask_driver_test.go](internal/console/runtime_multitask_driver_test.go#L283-L332).
- Fresh verification: I ran:
  `cd D:\LZProjects\prfrail; go test ./internal/console -run 'MultiTask|AgentRunnerAssembly|SupervisedAssembly|MultiTaskDefinitionShape|Assembly' -count=1`
  and it returned:
  `ok      github.com/larsonzh/prfrail/internal/console    3.541s`

## Section D
The claim is falsifiable by the exact failure modes the code is designed to reject:

- Duplicate task/step or duplicate request ID would fail in `NewAgentRunnerStepIntentRegistry`, which is the direct guard for task aliasing and request leakage. See [internal/adapters/agent_runner_intent.go](internal/adapters/agent_runner_intent.go#L27-L83).
- A still-alive PID or a missing PID report would fail closed in `managedProcessesGone` and therefore also in `proveTaskStopped`. See [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L853-L926).
- A non-serial or re-used run-root launch would be visible in the driver order and in the multi-task harness tests. See [internal/console/runtime_multitask_driver_test.go](internal/console/runtime_multitask_driver_test.go#L283-L332).

In other words, if the contract were being violated in the current tree, the nearest observable failure points are the registry binding, the stop-proof gate, or the multi-task serial-harness assertions; none of them fail under the live code.

## Section E
This is a current-state review of the checked tree, not a guarantee for every future path. The implementation comments themselves call out the constraint that the serial/offline assumption is what makes the `dispatched` observation sound and non-resetting. A future resume/concurrent driver that reuses stale state would require a fresh re-review; under the current serial/offline design, I did not find a contract gap.

I also checked the live worktree via `git status --short` and `git diff --numstat -- .`; no review-only file writes were introduced by this scan.

INDEPENDENT SCAN: PASS
```

## 6. 门禁读数

- 本地实测（Go 1.27.0 / windows）：`gofmt -l` 空输出、`go build ./...` exit 0、`go vet ./...` exit 0、`go test -count=1 ./...` **15 包全 ok、FAIL=0**（console 18.188s / adapters 51.380s）；原始输出 `04-impl/logs/04-full-suite-green.txt`。
- ⑤ r2 修复轮后复跑：同样 15 包全 ok（console 22.6s / adapters 51.6s / chain 6.6s）。
- ⑨ 实测：三作用域 tree / index / ci 均 `TOTAL_FAIL=0`（tree 与 index 实测数与 `G6-1 scanned` 数见 `08-gates/gate-readings.txt`）；`SELFTEST: PASS 697/697`；base 四腿（gofmt/build/vet/test）全绿（`ok=15`）。最终判据为 CI Go 1.22 双腿（Ubuntu + Windows）的 `Test` 与 `Gates` 步。
- 本地 `-race` 不在本片范围；六个核心包的 Race 腿由 P49 的 CI 覆盖承担。

## 7. 计量与环境事实

- **R = 1**：⑤ r1 整改轮（H1 / H2 / M1 / M2 与 I5 五项修复后重跑 ⑤）。
- **D = 1**：② 派单通道网络层失败，属环境事实而非代码缺陷；处置与取证见下。
- ② 派单通道失败 5 次：`UND_ERR_CONNECT_TIMEOUT` ×4、`ECONNRESET` ×1；第 4、5 次在后端已部分执行后才返回错误，工作区出现**两次半截产物**（CONTRACTS 非草案原文、supervised 侧编译红）。
- 处置：`git checkout` 精确回退 6 个受影响文件（保留步骤 0、③ 与证据包产物），半截 diff 取证为 `tmp/CONSOLE-MULTITASK-ASSEMBLY/02-partial-attempt-diff.txt`，阻断登记为 `02-impl-dispatch-blocked.txt`；恢复后 `go build ./...` 全绿，随后改为 ②a–d 小步串行派单。
- 环境事实：本地 Go 1.27.0（windows/amd64），仓库门禁与 CI 为 Go 1.22；版本差异由 CI 腿兜底。
- 探针额度：无（本片零真实外部模型候选调用）。

## 8. 未闭环项

- **CLI 接线后置**：`cmd/prfrail` 与 console flags 的多记录供给契约未定，留后续切片；本片以 `RequestRecords` 注册表为接口边界。
- **intent registry 生产未接**（⑦ Info 1）：`agentRunnerIntentMux` 仍逐记录构造 preparer，行为等价。
- **Windows R/C 写路径本地未跑**：Windows 上 `PublishDurability` 恒为 `unproven`，绿测只断言 per-requestId 路径与哈希互异；写路径由 Ubuntu CI 承担。
- **CI 首跑观察**：提交后须采集双腿 `Test` / `Gates` 结果；首跑红按 DR 口径登记。
- **既有导出哨兵保留**：`ErrSupervisedStepRequired` 未使用亦未删除（②c 偏差 6），未扩大改动面。

## 9. artifacts 围栏

- 过程件：`normalize-step0.js`、`verify-step0.js`、`rehash.js`、`02-impl-dispatch-blocked.txt`、`02-partial-attempt-diff.txt`、`02a-gate-index.txt`，全部随 `tmp/CONSOLE-MULTITASK-ASSEMBLY/` 按 `local-only` 处置（⑧c 清理）。

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/console-multitask-assembly.md	repo
docs/validation/console-multitask-assembly_EN.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/00-briefing.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/01-arch/design.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/03-tests/logs/03-tdd-red-run.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/03-tests/red-baseline.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/implementer-reports.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/logs/04-02a-gate-index.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/logs/04-full-suite-green.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/04-impl/logs/04-refusal-strings.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/05-prereview/review-r2.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/05-prereview/review.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/06-scan/review.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/07-final/review.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/08-gates/gate-readings.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/09-glm/glm-1-reference-sweep.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/09-glm/glm-2-test-coverage.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/MANIFEST.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/SHA256SUMS.txt	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/freedom-list.md	repo
docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/input-list.md	repo
internal/adapters/agent_runner_intent.go	repo
internal/adapters/agent_runner_intent_registry_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_multitask_agent_test.go	repo
internal/console/runtime_multitask_driver_test.go	repo
internal/console/runtime_multitask_red_test.go	repo
internal/console/runtime_multitask_supervised_test.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_test.go	repo
tmp/CONSOLE-MULTITASK-ASSEMBLY/	local-only
```

## 10. 文档同步更新状态

- `docs/ARCHITECTURE{,_EN}.md`：§3 追加三段（多任务串行装配 / A6 镜像共享语义 / 调度机制），各 +6 行（3 段正文 + 3 空行），CN / EN 行数相等。
- `docs/DEV_PLAN{,_EN}.md`：「11 月重构里程碑（T027-B）」列表由 6 条扩为 P37–P50 全量条目（含状态与报告路径），CN / EN 新增与删除行数相等。
- `docs/t027/REMAINING_SLICES{,_EN}.md`：状态表 P50 证据格回写、P50 小节追加「完成记录」行；状态符号已于 ⑧b 转 `✅ 已完成`；提交哈希与 CI 证据在授权提交后补记。
- 已同步项：`MANIFEST.md` 的 `08-gates/` 载荷已按 ⑨ 追加；`08-ci/` 待授权提交后按轮追加。

本片停在 ⑩ 等授权。

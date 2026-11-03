# ⑤ 预审结论原文 r2（原位重跑）—— P51 `CONSOLE-CLI-MULTITASK-WIRING`

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）。实得输出约 3900 字符（载体自报为估算值）。

结论：`PRE-REVIEW r2: PASS`（无 Critical/High/Medium；Low 1 / Info 4）。

---

## 1. r1 发现闭环核验

- **F1（High）已闭环**：新拒绝分支在 `cli.go:331-336`——`len(inputs.requestRecords)>0 && len(executable)==0 && len(supervised)==0` → `writeCommandError`（exit 1）。**零写入证据**：该分支位于互斥判定（328-330）之后、P1（339-342）之前，且早于 runID 默认（349-351）、`ValidID`（352-354）、`resolvedRunDir`（356-360）、eventLog stat（362-367）、`ToDefinition`（369-371）与 `ExecuteNoopRun`（385-388）。此前的全部 IO 均只读（`getwd`:303、`ResolveChainConfigPath`:307、`LoadChainConfig` `config.go:246-257` 纯 `os.ReadFile`+校验）；`writeCommandError`（756-764）仅写 stdout/stderr。run-root 与事件日志首次落盘发生在 `ExecuteNoopRun → chain.NewFileEventStore`（`runtime.go:33-40`），位于拒绝之后 ⇒ 零写入成立。补的测试 `TestRunRefusesRequestRecordsWithoutExecutableSteps`（test:241-272）带 `FindExecutableSteps==0` 前置断言，旧代码必走 noop 路径 exit 0 → 必红，**不能误绿**。
- **F3（Low）已闭环**：`keys []string` 按输入序收集（`cli.go:582-584, 606`）；stray 检查改为遍历 `keys`（616-620），首条 stray 键即报、输出确定；重复键在解码循环内先行判定（602-604）、覆盖检查遍历 `bindings` 切片（621-625），全部确定序。
- **F5（Info）已闭环**：注释改写（`cli.go:391-395`）为「request-record dimension … one of its two forms … all seven remaining values」，与 9 字段（396-406）相符：requestRecord/requestRecords 二选一 + 7 件套。

## 2. 新问题扫描

- 新分支与既有守卫交互：noop-only+records → 334-336 拒绝；supervised-only+records → supervised 步属 kind `code`（`config.go:575-585` 仅排除 noop），`executable>0` 分支不触发 → P1（339-342）拒绝、无回退；双 flag → 328-330 exit 2 优先。`complete()`（408-426）在 `requestRecords` 为空时首判退化为旧式 8 值语义（单值路径逐字节不变）；`anyProvided()`（430-449）仅前置追加。无新冲突。
- 测试误绿核验：8 条均断言专属文案（mutually exclusive / duplicate key / does not match run id / at least one executable step / does not bind a task/step / request record is missing for task），跨分支错因无法冒充；路径全为 `t.TempDir()`，无 event-log 碰撞。无误绿通道。
- 遗留项均非 r2 引入（r1 已登记）：F4 usage `...` 分组外歧义、F6 多记录解码错误无标签、⑥ 空路径腿仍缺、`cli.go:609-611` bindings==0 文案在「executable>0 但步形不符」时略失准（装配层兜底、零写入）。**(no new findings)**

## 3. 既有「无发现」项复核

1. 路由完整性：仍成立（互斥 exit 2 先于 P1 exit 1；F1 修复后无静默降级残余）。
2. 错误处理：仍成立（全失败面经 `writeCommandError`；失败立即返回、nil map 不进装配层）。
3. 单任务不变：仍成立（`complete()` 退化等价；绿腿 `TestRunSingleRequestRecordStaysCompatible` PASS；判定仍早于 run-root 创建）。
4. 口径漂移：仍成立（三条件与 `agentRunnerAssemblyPlan` `runtime_agent.go:1052-1072` 逐字一致，仅跳过 vs 整体拒绝，均零写入）。
5. flag 语义：仍成立（空/空白值 Parse 期即拒 exit 2；同路径重复 → 重复键拒绝）。
7. 回归面：仍成立（无测试引用 `runUsage`；gofmt/build/vet 干净；console 全绿）。

## 4. 剩余发现

| # | 严重度 | 位置 | 问题 | 建议 |
|---|---|---|---|---|
| L1 | Low | `cli_multitask_wiring_test.go` | ⑥ 空路径 Parse 拒绝（`cli.go:283-285`）无测试腿 | 可选补一条 |
| I1 | Info | `cli.go:388` | `runUsage` 的 `(...)...` 暗示 `--request-record` 可重复（F4 遗留，设计逐字给定） | 回馈设计 |
| I2 | Info | `cli.go:594-597` | 多记录解码错误无 `request record:` 标签（F6 遗留，与单值路径一致） | 可选包一层 |
| I3 | Info | `cli.go:343-347` | records+缺 7 件套时沿用旧措辞「local runtime only supports noop steps」，略误导 | 可选优化 |
| I4 | Info | 测试面 | 无 CLI happy-path 腿覆盖 `RequestRecords` 选项装配接线（全部 CLI 测试在装配入口前退出；装配层由 P50 驱动测试覆盖） | 可选 |

无 Critical / High / Medium。

## 5. 结论

PRE-REVIEW r2: PASS

OUTPUT_CHARS≈3900（正文字符估算）

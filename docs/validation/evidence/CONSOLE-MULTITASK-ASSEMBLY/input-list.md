# input-list — `CONSOLE-MULTITASK-ASSEMBLY`（P50）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-02，`HEAD` = `8bd8224`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/t027/CHECKLIST_v39.md` | — | 跨会话快照（锚点、四层队列、§八纪律 A–U、§九清单元规则与门禁坑） |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | — | 台账登记落点（状态表 + P50 小节） |
| `docs/validation/p45-e2e-zero-model.md` | — | 前置阅读：D9 结论（第 132 与 179 行区域） |
| `internal/console/runtime_agent.go` | — | 单任务限制：`agentRunnerAssemblyStep`（:772）要求恰好一个任务与一个步骤（D9 记 :774） |
| `internal/console/runtime_supervised.go` | — | 单任务限制：`supervisedAssemblyStep`（:333）要求恰好一个任务与一个步骤（D9 记 :335） |
| `docs/CONTRACTS{,_EN}.md` | — | 协议先行权威；A6 排他契约落点 |
| `docs/DELIVERY_DIRECTIVE.md` | — | 准则 v1.60（流水线角色、门禁、审查纪律） |

## 未复制入包

Go 源码与测试不复制；如需引用仅登记路径与行号。

## 锚点核对（步骤 0 实况）

| 项 | 值 |
|---|---|
| HEAD / origin/main | `8bd8224`（工作区干净；v39 文首锚点 `75990a8` 为 P49 主体提交，`8bd8224` 为其后 docs-only 回填，用户裁定以 `8bd8224` 为准） |
| 准则版本 | v1.60 |
| SELFTEST | 697/697 |
| 三作用域 | index / tree / ci 全 TOTAL_FAIL=0 |
| tmp 断言 | 顶层 = `.gitkeep` + `gates-cache`；`gates-cache` 下仅 `g4b` |

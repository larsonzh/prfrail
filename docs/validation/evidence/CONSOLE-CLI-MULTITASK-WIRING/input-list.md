# input-list — `CONSOLE-CLI-MULTITASK-WIRING`（P51）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-03，`HEAD` = `3aa53d5`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/t027/CHECKLIST_v40.md` | 21915 | 跨会话快照（锚点、四层队列、§八纪律 A–V、§九清单元规则与门禁坑） |
| `docs/t027/REMAINING_SLICES.md` | 281505 | 台账登记落点（状态表条目 + 完成记录） |
| `docs/t027/REMAINING_SLICES_EN.md` | 315650 | 台账英文镜像落点 |
| `docs/validation/console-multitask-assembly.md` | 18109 | P50 报告：§8 未闭环项「CLI 接线后置」，接口边界 = `RequestRecords` 注册表 |
| `internal/console/cli.go` | 25383 | CLI 命令路由与 `run` / `supervised` 参数解析；`agentRunnerRunInputs` 单记录八件套 |
| `cmd/prfrail/main.go` | 268 | 进程入口，`console.NewCLI(version).Execute` 单点转发 |
| `internal/console/runtime_agent.go` | 66612 | `AgentRunnerAssemblyOptions`（`RequestRecord` 单任务兼容 + `RequestRecords` 注册表）与 `ExecuteAgentRunnerRun` |
| `internal/console/runtime_supervised.go` | 97459 | `SupervisedAssemblyOptions` 与 `ExecuteSupervisedRun` / `ExecuteSupervisedReturn` |
| `internal/console/supervised_cmd.go` | 9019 | `supervised` 子命令族（run/return/report）参数解析 |
| `internal/console/preview.go` | 10888 | `preview` 子命令实现（是否同为接线入口待 ① 判定） |
| `internal/console/runtime.go` | 3401 | `RunSummary` 与 noop 运行入口 |
| `docs/CONTRACTS.md` | 78534 | 协议先行权威：A6 段「多任务串行装配（P50）」为接口唯一边界 |

## 未复制入包

Go 源码与测试不复制；如需引用仅登记路径与行号。

## 锚点核对（步骤 0 实况）

| 项 | 值 |
|---|---|
| HEAD / origin/main | `3aa53d5`（工作区干净；v40 文首锚点 `81f9dc9` 为 P50 主提交，`3aa53d5` 为其后 docs-only 回填，起跑消息已声明该代差免责） |
| 准则版本 | v1.60 |
| SELFTEST | 697/697 |
| 三作用域 | index / tree / ci 全 TOTAL_FAIL=0（基线日志 `tmp/CONSOLE-CLI-MULTITASK-WIRING/baseline-{index,ci,tree}.log`） |
| tmp 断言 | 顶层 = `.gitkeep` + `gates-cache`；`gates-cache` 下仅 `g4b`；本片新增 `tmp/CONSOLE-CLI-MULTITASK-WIRING/` |

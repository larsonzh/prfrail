# CONSOLE-CLI-MULTITASK-WIRING 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `CONSOLE-CLI-MULTITASK-WIRING`（M 档 · 实现类 · `[SLICE]`） |
| 定位 | v40 §四 层 2 第 3 项；补齐 P50 遗留 CLI 接线（P50 报告 §8 明示「CLI 接线后置」，接口边界为 `RequestRecords` 注册表） |
| 可自行决定的自由度 | ① CLI 参数接收形态（多次 `--request-record` 重复标志、`<taskID>/<stepID>=<path>` 键值形、或单文件清单）在 ① 架构方案内自行决定；② 红测形态（CLI 端到端调用 `Execute` 或装配入口注入断言）；③ 文档回写的措辞与落点小节；④ 批量辅助（GLM）派单的具体检索问题 |
| 必须升级的事项 | 需修改 `CONTRACTS{,_EN}` 的 A6 多任务串行装配段 ⇒ ① 先出精确协议修订方案并升级；需改装配层语义（`runtime_agent.go` / `runtime_supervised.go` 的驱动或排他逻辑）；需动 `internal/chain/**` 模型语义；需破坏既有单任务 CLI 兼容性 ⇒ 明确说明并交 ⑤⑥⑦ 确认 |
| 载体铁律 | ① ⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（可写）；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；批量辅助 `deep-reasoner` + `GLM-5.3-Flash (glm)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的代码与测试；主控只做编排、落盘、门禁与回执 |
| 纪律 V 拆步 | ② 必须按「参数定义」→「装配接驳」→「测试配合」小步串行派单，禁止一次性全量载荷；超时后先查半截产物再决定回退 |
| 临时件 | 一律 `tmp/CONSOLE-CLI-MULTITASK-WIRING/`，只增不删，⑧c 清理 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（`⑩` 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |
| 本地环境事实 | 本机 Windows + Go 1.27.0（无 gcc，本地 `-race` 受限）；本片为 CLI 与参数解析层，不涉及并发敏感面，本地三作用域门禁 + `go build/vet/test` 足够 |
| 主控模型 | 本片 M 档且不触及底层引擎状态机 ⇒ 主控维持 `DeepSeek V4.1 Flash`（v40 §七 第 8 条）；连续 2 次行为漂移或逻辑遗漏时按纪律 P 提示切换 Pro |

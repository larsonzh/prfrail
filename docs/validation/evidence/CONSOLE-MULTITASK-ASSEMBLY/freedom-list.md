# CONSOLE-MULTITASK-ASSEMBLY 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `CONSOLE-MULTITASK-ASSEMBLY`（L 档 · 实现类 · `[SLICE]`） |
| 定位 | v39 §四 层 2 第 3 项；断点 A 收口；协议先行（D9 与 A6 排他契约冲突须先裁定） |
| 可自行决定的自由度 | ① 多任务装配的调度形态（并发上限、排队与并发取舍、拒绝语义）在 ① 架构方案内自行决定；② 红测的形态（聚焦 `ExecuteNoopRun` / `localRuntime` 单任务守卫，或端到端装配入口）；③ 文档回写的措辞与落点小节；④ 批量辅助（GLM）派单的具体检索问题 |
| 必须升级的事项 | 需修改 `CONTRACTS{,_EN}` 或底层锁机制 / A6 镜像契约 ⇒ ① 先出精确协议修订方案并升级；需扩出 A6 排他语义之外的能力；需动 `internal/chain/**` 语义；需降级（部分多任务 / 撤回）⇒ 报告中明确说明并交 ⑤⑥⑦ 确认 |
| 载体铁律 | ① ⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（可写）；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；批量辅助 `deep-reasoner` + `GLM-5.3-Flash (glm)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的代码与测试；主控只做编排、落盘、门禁与回执 |
| 临时件 | 一律 `tmp/CONSOLE-MULTITASK-ASSEMBLY/`，只增不删，⑧c 清理 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（`⑩` 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |
| 本地环境事实 | 本机 Windows + Go 1.27.0（无 gcc，本地 `-race` 受限）；并发敏感面如需真实 `-race`，复用 P49 已标定的 Ubuntu VM 通道或 CI Race 腿（环境预检与复用铁律：先探查已有环境，确认缺失再最小化安装，严禁盲目重装） |
| 主控模型 | 本片为状态机/核心装配类 L 档 ⇒ 主控已提升为 `DeepSeek V4 Pro`（按 v39 §七 第 8 条与纪律 P）；任何降级须按 OB-37 元规则独立登记 |

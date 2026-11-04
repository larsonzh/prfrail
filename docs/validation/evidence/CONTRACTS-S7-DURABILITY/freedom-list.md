# CONTRACTS-S7-DURABILITY 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `CONTRACTS-S7-DURABILITY`（L 档 · 契约与实现类 · `[SLICE]`） |
| 定位 | v43 §四 层 3 第 6 项；R-02 裁定已批准立项（2026-11-04 起跑消息第 1 条）；目标 = Windows 平台发布耐久等级 `unproven` → `proven`，落实 C6，全面修订 `CONTRACTS.md` §7 |
| 可自行决定的自由度 | ① C6 的具体落地方式（`MoveFileEx` 调用形态、封装位置与平台构建标签组织，由 ① 架构评估后定稿）；② ③ 红测的平台构建标签隔离写法与红测步骤编排（TDD 红测先行）；③ ④ 跨平台构建验证（GOOS=windows 与 linux）的具体命令组合；④ 文档回写的措辞与落点小节（§八 六项文档的具体小节）；⑤ 批量辅助（GLM）派单的具体检索问题；⑥ 临时脚本与中间产物的内部组织（限 `tmp/CONTRACTS-S7-DURABILITY/`） |
| 必须升级的事项 | 改变任何既有接口签名（起跑消息第 3 条「确保接口签名不变」）；改动 Unix 路径的生产代码（Unix 零影响为硬约束）；把 C6 语义写入 Unix 构建（不得让 Unix 依赖 Windows 专属 API）；修改 `CONTRACTS` 中超出本片授权范围的段落；把耐久等级升级范围扩大到 Windows 之外的平台；对既有测试的删除或语义弱化；把 visible 人工交互语义写进 ProofRail 正式协议 |
| 载体铁律 | ① ⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（可写）；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；① 红队 `independent-reviewer` + `GPT-6.1 Sol (copilot)`（§1.7.4 临时授权，限本片一次性）；批量辅助 `deep-reasoner` + `GLM-5.3-Flash (glm)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model`；严格串行禁止并行 |
| 主控禁止 | 禁止代写 ②③ 的生产代码与测试；主控只做编排、落盘、门禁与回执 |
| 纪律 V 拆步 | 先改 CONTRACTS 文本 → 再实现 C6 封装 → 再改 replay_store 调用 → 再改测试；严禁打包大改；超时后先查 `git status` 半截产物再决定回退 |
| 红测先行红线 | ③ 红测必须在 ② 生产代码修改**之前**跑出真实红（Windows 原生环境，证明当前平台发布失败 / 处于 `unproven`）；若红测未按 TDD 时序执行即属偏离 |
| ① 双打红线 | 主控完成 ① 架构后，必须临时调用 `independent-reviewer` + `GPT-6.1 Sol (copilot)` 红队抗压审查（§1.7.4 临时授权、白名单外、一次性、不得作为先例）；红队重点：契约边界是否夸大、C6 同步机制在异常路径下是否退化为普通写、非 Windows 平台抽象是否纯净；红队发现缝隙可直接打回 V4 Pro 修改 |
| 临时件 | 一律 `tmp/CONTRACTS-S7-DURABILITY/`，只增不删，⑧c 清理；测字符数用内联脚本，严禁在 `$env:TEMP` 等系统临时目录生成文件 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（⑩ 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/CONTRACTS-S7-DURABILITY/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |
| 环境事实 | 本机 Windows + Go（直连，无代理）；**无需开启 Ubuntu VM**（起跑消息第 4 条）；Unix 零影响用本地 `GOOS=linux go build/vet` 验证，由 CI ubuntu 腿兜底 |
| 主控模型 | 本片为底层契约与持久化 API 修改 ⇒ 主控提升为 `DeepSeek V4 Pro`（v43「纪律 P」+ 起跑消息第 3 条明示）；降级须按 OB-37 元规则独立登记 |
| 文档同步 | ⑧ 阶段必须同步六项文档（`CONTRACTS`、`ARCHITECTURE`、`SECURITY`、`DEV_PLAN`、`ADR_REGISTER`、`REMAINING_SLICES`，均 CN+EN），缺一即视为遗漏 |

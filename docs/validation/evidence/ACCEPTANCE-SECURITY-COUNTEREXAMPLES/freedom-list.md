# ACCEPTANCE-SECURITY-COUNTEREXAMPLES 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `ACCEPTANCE-SECURITY-COUNTEREXAMPLES`（M 档 · 安全验证类 · `[SLICE]`） |
| 定位 | v41 §四 层 2 第 4 项；承接断点 A 的接受链缺口（GAP-51/52/55/60/64/70）；只做测试用例与断言，不改接受语义 |
| 可自行决定的自由度 | ① 反例测试的命名清单与断言矩阵（在 ① 架构方案内定稿）；② 红测形态（攻击被成功阻断 = 测试 PASS，验证 fail-closed 逻辑）；③ 测试文件落点（各 `*_test.go` 所在包由 ① 架构方案界定）；④ 文档回写的措辞与落点小节（`ARCHITECTURE{,_EN}` 安全边界或验证章节、`DEV_PLAN{,_EN}` T027-B 里程碑）；⑤ 批量辅助（GLM）派单的具体检索问题 |
| 必须升级的事项 | 修改 `CONTRACTS{,_EN}` 任何段落；修改生产代码（`internal/**` 非 `_test.go` 文件）；改变接受语义或门禁判定逻辑；涉及跨包公共接口变更；涉及对既有测试的删除或语义弱化 |
| 载体铁律 | ① ⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（可写）；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；批量辅助 `deep-reasoner` + `GLM-5.3-Flash (glm)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的测试代码；主控只做编排、落盘、门禁与回执 |
| 纪律 V 拆步 | ② 必须按「GAP-51/52 批次」→「GAP-55/60 批次」→「GAP-64/70 批次」小步串行派单，每批验证后再派下一批，禁止一次性全量载荷；超时后先查半截产物再决定回退 |
| 真实漏洞上报 | 若反例测试 FAIL（攻击未被阻断）⇒ 立即停止，按真实安全漏洞上报，本片内绝不修改生产代码 |
| 临时件 | 一律 `tmp/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/`，只增不删，⑧c 清理；测字符数用内联脚本，严禁在系统临时目录生成文件 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（⑩ 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |
| 本地环境事实 | 本机 Windows + Go（无 gcc，本地 `-race` 受限）；本片为测试反例验证，`go test ./...` + 三作用域门禁为本地验收基线；Linux 腿由 CI 承担 |
| 主控模型 | 本片为安全对抗/验证类切片 ⇒ 主控提升为 `DeepSeek V4 Pro`（v41「纪律 P」+ 起跑消息第 3 条明示）；降级须按 OB-37 元规则独立登记 |

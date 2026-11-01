# DR13-FIX 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `DR13-FIX`（S 档 · 修复片 · `[SLICE]`） |
| 定位 | P45 的前置修复片；不占 v34 队列 1–4 编号；闭环后出 v35 统一反映 |
| 可自行决定的自由度 | ① 在 A（装置过严）与 B（真实漏洞）之间独立定性，不得被主控预判锚定；② 若判定为 A，可选择「改断言为同一时刻至多一个存活」或「引入可控的进程存活屏障」两种修法之一；③ 若判定为 B，可改动 `agent_runner_pinned_cli*.go` 的排他或 reclaim 路径 |
| 必须升级的事项 | 修法需要突破本片边界（`internal/adapters/agent_runner_pinned_cli*` 之外）；需要改动 `ci.yml`；定性结论为「A 与 B 并存」需同时改实现与断言 |
| 载体铁律 | ① ⑤ 只读（`deep-reasoner`）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的代码与测试；主控只做编排、落盘、门禁与回执 |
| 临时件 | 一律 `tmp/DR13-FIX/`，只增不删，⑧c 清理 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（`⑩` 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/DR13-FIX/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |

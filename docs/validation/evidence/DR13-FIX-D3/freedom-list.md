# DR13-FIX-D3 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `DR13-FIX-D3`（S 档 · 实现类 · `[SLICE]`） |
| 定位 | v37 §四 层 1 第 1 项；R-05 已裁独立立项 |
| 可自行决定的自由度 | ① 在 D3 落点（`startingMarkerClaim.release()`）上自行决定「先 remove 再 unlock」的实现细节与注释；② 自行决定红测的卡点装置（在未修复代码上确定性触发 unlock 与 remove 之间的间隙）；③ 可在 `internal/adapters/agent_runner_starting_marker*.go` 及其测试文件范围内自行增补回归 |
| 必须升级的事项 | 修法需要突破本片边界（`internal/adapters/agent_runner_starting_marker*` 之外）；需要改动 CONTRACTS / Schema；「先 remove 再 unlock」被证明在某一平台不安全，需回退或引入平台分叉 |
| 载体铁律 | ① ⑤ 只读（`deep-reasoner`）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的代码与测试；主控只做编排、落盘、门禁与回执 |
| 临时件 | 一律 `tmp/DR13-FIX-D3/`，只增不删，⑧c 清理 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（`⑩` 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/DR13-FIX-D3/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |

# RACE-LEG-EXPAND 执行自由度清单

| 项 | 内容 |
|---|---|
| 切片 | `RACE-LEG-EXPAND`（M 档 · 基础设施类 · `[SLICE]`） |
| 定位 | v38 §四 层 1 第 2 项；R-03 已裁独立立项（层 1 治理收尾） |
| 可自行决定的自由度 | ① 在 `.github/workflows/ci.yml` 的 Race 腿目标包扩展范围内自行决定包的排列顺序与注释措辞；② 自行决定 ③ 先行断言的脚本形态（静态检查 `ci.yml` 参数行，或本地跑非 race 测试佐证四包存在并发用例）；③ 自行决定文档回写的措辞与落点小节 |
| 必须升级的事项 | 需要删除、重排或新增任何既有 CI 步骤（本片边界只许改 Race 腿的目标包清单）；需要把 Race 腿扩展到 `internal/` 六个包之外；需要改 CONTRACTS / Schema / 生产代码；因环境或超时需要降级（部分扩包或撤回）时，须在报告中明确说明并交 ⑤⑥⑦ 确认 |
| 载体铁律 | ① ⑤ 只读（`deep-reasoner`）；② `prfrail-implementer`、③ `prfrail-tester`、⑧ `prfrail-documenter`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；每次 `runSubagent` 必须同时传 `agentName` 与 `model` |
| 主控禁止 | 禁止代写 ②③ 的代码与测试；主控只做编排、落盘、门禁与回执 |
| 临时件 | 一律 `tmp/RACE-LEG-EXPAND/`，只增不删，⑧c 清理 |
| 禁止 | 推 gitee；在 ⑩ 之前 commit/push（`⑩` 停点等同一轮授权） |
| 证据包 | `docs/validation/evidence/RACE-LEG-EXPAND/`（`00-briefing.md` 为起跑消息原文；`.md` 一律 UTF-8 + BOM + LF） |
| 本地环境事实 | 本机无 gcc、无 actionlint、`CGO_ENABLED=0`（windows/amd64，Go 1.27.0）⇒ 本地 `-race` 实跑受限；Race 真实性以 Ubuntu CI 腿为准，本地只做静态覆盖断言与非 race 佐证 |

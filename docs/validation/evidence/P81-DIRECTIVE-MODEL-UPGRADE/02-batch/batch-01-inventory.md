# P81 批量辅助落点清单（`GLM-5.3-Flash`）

> 载体 / 模型：`deep-reasoner` / `GLM-5.3-Flash (glm)`（只读，非流水线角色，§2.2 的批量辅助四类任务之「枚举 / 统计」）。
> 逐字转录自批量辅助的最终消息；时点 = ② 落盘后、③ 之前。

## 命中明细表

| 文件 | 行号 | 命中标识 | 所在小节 / 上下文摘要（≤20 字） | 层（A/B/C） |
|---|---|---|---|---|
| docs/DELIVERY_DIRECTIVE.md | 65 | `gpt-6.1-sol` | §0.3 v1.18 变更日志行 | B |
| docs/DELIVERY_DIRECTIVE.md | 117 | `MAI-Code-1.1-Flash` + `GPT-6.1 Sol` | §0.3 v1.70 变更日志行 | B |
| docs/DELIVERY_DIRECTIVE.md | 1092 | `gpt-6.1-sol` | §12.4 OB-53 台账行 | B |
| docs/DELIVERY_DIRECTIVE.md | 1101 | `GPT-6.1 Sol` + `gpt-6.1-sol` | §12.4 OB-62 台账行 | B |
| docs/DELIVERY_DIRECTIVE.md | 1107 | `gpt-6.1-sol` | §12.4 OB-68 台账行 | B |
| docs/DELIVERY_DIRECTIVE.md | 1171 | `GPT-6.1 Sol` | §12.4 OB-132 台账行 | B |
| docs/DELIVERY_DIRECTIVE.md | 1703 / 1723 / 1744 | `MAI-Code-1.1-Flash` | 附录 C 的 C.0d / C.0e / C.0f 台账块 | B |
| docs/DELIVERY_DIRECTIVE.md | 1805 | `gpt-6.1-sol` | 附录 C 的 C.0j 台账块 | B |
| docs/DELIVERY_DIRECTIVE_EN.md | 65 / 117 / 1092 / 1101 / 1107 / 1171 / 1703 / 1723 / 1744 / 1805 | 与 CN 逐行同构、同标识分布 | 变更日志 / 台账 / C.0d–f、C.0j | B |
| tools/gates/selftest.js | 4234–4235 | `gpt-6.1-sol` | T131b 退役标识拒绝的自测断言 | A |
| tools/gates/selftest.js | 4241–4242 | `mai-code-1.1-flash` | T131c 退役标识拒绝的自测断言 | A |
| tools/gates/testdata/fixtures/g6-1-pending-done.txt | 6 | `MAI-Code-1.1-Flash` | G6-1 历史行转写夹具（钉死不改） | C |
| docs/t027/REMAINING_SLICES.md | 2179 | `MAI-Code-1.1-Flash` + `GPT-6.1 Sol` | P81 标识行（派单源） | A |
| docs/t027/REMAINING_SLICES.md | 2268 | `mai-code-1.1-flash` | DR-1 回写记录（2026-09-24） | B |
| docs/t027/REMAINING_SLICES.md | 2358 | `gpt-6.1-sol` | 完成记录行（2026-10-03） | B |
| docs/t027/REMAINING_SLICES_EN.md | 2166 | 同 CN 标识行 | P81 identifier row | A |
| docs/t027/REMAINING_SLICES_EN.md | 2255 | `mai-code-1.1-flash` | DR-1 write-back | B |
| docs/t027/REMAINING_SLICES_EN.md | 2345 | `gpt-6.1-sol` | completion log | B |
| （未跟踪新夹具，不计受控数）`tools/gates/testdata/fixtures/g5b-illegal-retired-sol.txt:4`、`g5b-illegal-retired-mai.txt:4` | — | `gpt-6.1-sol` / `mai-code-1.1-flash` | G5-b 退役负样本 | A*（未入库） |

## 小结一（各层条数）

- **A 层 = 6 行（受控）**：`selftest.js` 4 行 + `REMAINING_SLICES` 标识行 2 行；另有 2 行未跟踪新夹具（A*）。准则两文件的 A 层 = **0**（现役正文已替换）。
- **B 层 = 26 行**：两准则文件各 11 行 + `REMAINING_SLICES` 记录行 4 行。
- **C 层 = 其余全部**：受控全集 755 处命中（`MAI-Code-1.1-Flash`=423 / `mai-code-1.1-flash`=117 / `GPT-6.1 Sol`=124 / `gpt-6.1-sol`=91），其中面外集中在 `docs/validation/*.md` 报告对 246 处、`docs/validation/evidence/**` 冻结包 341 处、`docs/t027/` 其余 119 处、`ADR_REGISTER` 与 `DEV_PLAN` 8 处、`.github/agents/sol-orchestrator.agent.md` 1 处与顶层 `testdata/**` 3 处。

## 小结二（每文件计数）

- `docs/DELIVERY_DIRECTIVE.md`：`MAI-Code-1.1-Flash`=4、`mai-code-1.1-flash`=0、`GPT-6.1 Sol`=3、`gpt-6.1-sol`=5。
- `docs/DELIVERY_DIRECTIVE_EN.md`：与 CN 完全相同（行号一一对应）。
- `tools/gates/**`（受控）：`MAI-Code-1.1-Flash`=1、`mai-code-1.1-flash`=2、`GPT-6.1 Sol`=0、`gpt-6.1-sol`=2；`REMAINING_SLICES{,_EN}.md`：各文件四标识各 1 处。

## 小结三（多标识行与异常发现）

- **同一行多标识：受控全集 84 行**。要求面内 6 行 = `DD.md` L117 / L1101、`DD_EN.md` L117 / L1101、`RS.md` L2179、`RS_EN.md` L2166；其余 78 行全在 C / 报告面。
- **异常**：① 枚举时点 = 本片替换已落盘未提交 ⇒ §2.2、§6.3 的 G5-b 闭集等现役行已是 `MAI Code 2 Flash` / `gpt-6.2-sol`，旧标识仅存历史面；② `g5b-closed-set-legal-sol.txt` 的 HEAD 版本仍含 `gpt-6.1-sol`（受控历史保留旧标识）；③ EN 附录 C.0f 混入未译中文「实得」（镜像瑕疵，历史块，§11.1 不追溯）；④ `sol-orchestrator.agent.md` L74 仍提及已退役的 `MAI-Code-1.1-Flash` 家族（非准则面陈旧引用）；⑤ 本片证据包（`00-briefing.md` / `01-arch/design-r1.md` 含旧标识引述）属未跟踪件，不在本枚举内。

BATCH: DONE

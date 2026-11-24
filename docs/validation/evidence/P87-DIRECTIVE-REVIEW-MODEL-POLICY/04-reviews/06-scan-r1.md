# P87 ⑥ 独立扫描 r1（`MAI Code 2 Flash (copilot)`；2026-11-24；逐字原文）

### Section A 覆盖面与未覆盖面

已覆盖：切片定义、①方案书、自由度清单、原始 diff、修复摘要、三作用域门禁、v1.72 条文、ADR-048、G5-b 实现与 T132、双语形状及证据包。未覆盖：P85 原始评审过程、④⑦ 真实子代理调用回执、Hydrogen 对照评审、build/vet 的退出码，以及尚处 ⑧a 后置阶段的报告对成文。

### Section B 发现清单

| # | SEVERITY | 位置(文件:行) | 发现 | 建议 |
|---|---|---|---|---|
| 1 | Medium | [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L574)；CN 对位 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L574) | 轻量档英文判据写成 `final-review passes`，不等价于 CN 的「定稿复核与定向小读」，可能误导档位触发语义。 | 英文改为 `finalization review and targeted small reads for document-type / governance slices (low-cost review)`。 |
| 2 | Low | [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L830)；[docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L830) | Sol 1M 400/1500 被标注为“翻倍/double”，但仅输入翻倍、输出为 1500/1000 即 1.5 倍。 | 删除倍率形容，或明确写 `input 2x, output 1.5x`。 |
| 3 | Low | [build.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/build.txt)；[vet.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/vet.txt) | 两份落盘文件均为 0 行，不能在无退出码下独立证明 build/vet 退出成功。 | 仅在证据中补 `EXIT=0` 或保留命令执行摘要。 |

### Section C 判据实现与自测

G5-b 闭集三处一致，顺序均为 8 枚既有成员后追加 `gemini-4-neon`、`gemini-4-carbon`、`gemini-4-hydrogen`：CN 条文见 [§2.2](docs/DELIVERY_DIRECTIVE.md#L309) 与 [§6.3](docs/DELIVERY_DIRECTIVE.md#L680)，实现见 [g5b.js](tools/gates/lib/criteria/g5b.js#L25)。

三个合法夹具与一非法夹具已落盘；T132a-c 分别覆盖三枚合法、`gemini-4-pro` 非法和 Carbon 删除载荷性变异，落盘断言位于 [T132a](tools/gates/selftest.js#L4272)、[T132b](tools/gates/selftest.js#L4281)、[T132c](tools/gates/selftest.js#L4288)。落盘自测中对应三条均为 `ok`，末态为 [744/744](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-selftest.txt#L749)。

### Section D 边界与越界

原始 diff 27 个路径，仅落准则双语、ADR 双语、台账双语、证据包、G5-b 实现、自测夹具与台账段；未发现 CONTRACTS/schema/契约夹具/workflows/`internal/**`/CI/frozen snapshots变更。原有 ⑥ 首选 MAI 与 Sol ⑦ 地位未改；Grok 与其他草案外模型未新增。报告对不在 r2 原始 diff 中，但权威步骤链将报告对置于 ⑦ 后的 ⑧a，此点暂不视作越界或已完成交付。

### Section E 我无法验证的部分

- P85 四读的原始模型报告与分歧结论仅能从起跑消息、ADR 和复述材料核实，未能在本输入中独立复跑。
- build/vet 的原始退出码缺失，不能独立证明两者全绿。
- gate-ci 报告只覆盖 2 个端点变更，而原始 diff 为 27 个；无法证明 staged P87 diff 本身被该 CI leg 执行。
- Hydrogen 对照评审、⑤/⑦/盲审调用回执、模型限定名 UI 实读均未随输入提供。
- P87 最终报告对、⑧c 收尾、⑩ 回执等后置产物尚未进入本次 r2 扫描面。

INDEPENDENT SCAN: FINDINGS

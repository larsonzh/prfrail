# P81 ⑦ 终审 r1 原文（`DIRECTIVE-MODEL-UPGRADE`）

> 载体 / 模型：`independent-reviewer` / `GPT-5.3-Codex (copilot)`（只读）。首轮不附任何 ⑤/⑥ 清单与自由度清单（§7.3 锚定硬隔离）。逐字转录自 ⑦ 的最终消息。
> 编码：UTF-8 with BOM + LF。

## 四项审查

1. 安全性：通过。
读取主控落盘的实现与读数后，未见 fail-open 放宽。G5-b 闭集仅替换两枚成员，比较语义仍为大小写敏感、逐字相等、两端 trim、引号不等价，判据本体未放宽（`tools/gates/lib/criteria/g5b.js` L21-L59）。
退役标识拒绝与新标识准入均有实测断言：T131a/T131b/T131c 在自测输出为通过（`tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-selftest-04.txt` L478-L480），且总自测全绿（同文件 L745）。
本次变更未触及任何 family 路由/旁路逻辑文件；变更面集中在准则文本、G5-b 闭集、selftest 与夹具（`tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.numstat`）。

2. 架构一致性：通过。
准则正文、判据实现、自测夹具三方同源一致：
- 准则 2.2/5.2/6.3/7.5/8.2 已统一为新 modelId（`docs/DELIVERY_DIRECTIVE.md` L303 / L305 / L561 / L641 / L785 / L866）。
- 判据闭集同步（`tools/gates/lib/criteria/g5b.js` L21-L31）。
- 夹具与断言同步（`tools/gates/testdata/fixtures/g5b-closed-set-legal-sol.txt` L4、`g5b-closed-set-legal-mai.txt` L4、`g5b-illegal-retired-sol.txt` L4、`g5b-illegal-retired-mai.txt` L4、`tools/gates/selftest.js` L4225-L4242）。
版本元数据与附录 C 台账位置符合既有结构（`docs/DELIVERY_DIRECTIVE.md` L3 / L117 / L2370；EN 镜像同位 `docs/DELIVERY_DIRECTIVE_EN.md` L3 / L117 / L2370）。未引入第二权威源。

3. 完整性：通过。
现役落点已全换名；历史面保留旧标识符合 11.1 不追溯原则（`docs/DELIVERY_DIRECTIVE.md` L955）。
P80 留痕补落已包含三要素：授权原文、输入收窄声明、不得作先例（`docs/DELIVERY_DIRECTIVE.md` L1187），且来源指针回链到 P80 报告 14.2/14.3（`docs/validation/sg6-gatelevel-consolidate.md` L193-L203）。
台账 P81 段落点齐全且仍处进行中状态，未误写完成（`docs/t027/REMAINING_SLICES.md` L2175-L2195）。

4. 测试反例审查：通过。
三条新断言具判别力且与目标一一对应（新成员接受、两退役成员拒绝），并在主控落盘自测中命中通过（`tools/gates/selftest.js` L4225-L4242、`tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-selftest-04.txt` L478-L480）。
对偶覆盖完整：新成员被接受 + 退役成员被拒绝均已覆盖；大小写、引号、尾部空格边界并未遗漏（既有 T32/T33/T39/T40 仍在通过序列中，见同一自测原始读数 `tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-selftest-04.txt`）。
改后门禁树作用域为零失败（`tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-tree-04b.txt` L52）。

## 发现清单

(no findings)

## 已核查、非问题

- 已核查、非问题：G5-b 在本次 tree 读数中显示 checked 0 file(s)（`tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-tree-04b.txt` L18）是因为本片未改 `.github/agents` 目标文件；该点由 selftest 740/740 覆盖，不构成漏检。
- 已核查、非问题：P80 输入收窄在 §12.4 为「补落引用块」表达，且同段已明确「仅本片有效、不得作为先例」，并指回原始报告位置（`docs/DELIVERY_DIRECTIVE.md` L1187）；未造成新增权威源。
- 已核查、非问题：CN/EN 对称性保持，版本头与 0.3 最新行一致（`docs/DELIVERY_DIRECTIVE.md` L3 / L117、`docs/DELIVERY_DIRECTIVE_EN.md` L3 / L117）。

## 追加独立复核的判定

不需要追加制度外独立复核。理由：本环已按锚定硬隔离独立阅读，且 ⑥/⑦ 虽同载体但模型不同，未见同模型自审风险；当前无未闭合发现项，仅需按既定流程进入下一环节即可。

FINAL REVIEW: PASS

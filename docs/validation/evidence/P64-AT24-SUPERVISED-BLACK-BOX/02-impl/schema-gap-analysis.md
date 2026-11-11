# schema-gap-analysis — P64 裁决① 协议先行第 1 步

> 目的：逐字核对 CONTRACTS §7.1(d) 与 `schemas/supervised-report.schema.json` 现状差异，判定 schema 是否宽松于裁定口径，为 ② 加严提供依据（2026-11-11）。

## 1. CONTRACTS §7.1(d) 原文（逐字）

> （d）**reduced/unknown 报告**：完成后写一枚 `supervised-report` 记录，`assuranceLevel` 恒为常量 `reduced`，`unknownDimensions` 非空且必须含 `cost`，逐项列出外部代理的工具使用、资源读取、网络、费用、后台进程、会话连续性与外部副作用等不可观测维度。**禁止把未知项显示为安全**：报告与渲染对每个未知维度只显示 `unknown`，不得显示 `safe`、`verified` 或任何数字承诺。

## 2. R-09.8 裁定与腿 5 判据（逐字）

- R-09.8（`docs/validation/r08-r09-ruling.md`）：「assuranceLevel=reduced 且 unknownDimensions 非空且含 cost」。
- 起跑消息 3.1 腿 5 判据：「`supervised-report{assuranceLevel=reduced}` 的 `unknownDimensions` 非空且含 `cost`；七项（工具使用/资源读取/网络/费用/后台进程/会话连续性/外部副作用）逐项必须出现（反例 = 缺任何一项 ⇒ 报告不合规）」。

## 3. schema 现状（`schemas/supervised-report.schema.json:142-153`）

```
"unknownDimensions": {
  "type": "array",
  "minItems": 1,
  "uniqueItems": true,
  "items": { "$ref": "#/$defs/unknownDimension" },
  "contains": { "const": "cost" },
  "minContains": 1
}
```

其中 `$defs/unknownDimension` 为封闭枚举（七项：tool-use / resource-reads / network-access / cost / background-processes / session-continuity / external-side-effects）。schema 顶层 `description` 亦只写「non-empty and must contain cost」。

## 4. 差异判定

schema 当前约束 = 非空 + 含 cost + 闭集内 + 唯一 ⇒ **接受任意子集**（如 `["cost"]`、七项缺 `network-access` 的六项集）。R-09.8「七项全必填」要求七项逐项出现。
**结论：schema 真宽松于 R-09.8 裁定口径。**（③ 实测：`p5-only-cost`、`p5-missing-network-access` 均 schema ACCEPT、CLI exit=0。）

## 5. 「等」字如实登记

CONTRACTS §7.1(d) 以「等」字收尾（「外部副作用等不可观测维度」），字面允许扩展项。如实登记如下：CONTRACTS 本体**不改**；按用户裁决与 R-09.8 已裁，以**七项闭集为准**加严 schema（「逐项列出」+ 七项全必填）；若未来契约扩维，按契约变更流程另行处理。

## 6. fixtures 现状

- valid：`supervised-report-valid`（`testdata/contracts/valid/runtime-records.json:364`）——`unknownDimensions` 已含七项全集 ✓（无需改）。
- invalid（`testdata/contracts/invalid/runtime-records.json`）：`assurance-level-full`、`unknown-dimensions-empty`、`unknown-dimensions-missing-cost`（`["tool-use","network-access"]`）、`unknown-dimension-id-unknown`（`["cost","model-latency"]`）——**缺「缺任一非 cost 项」新反例**（如七项缺 `network-access` 的六项集）⇒ 需补 1 件。

## 7. 契约测试断言同步点

`tools/contracts/contracts.test.js:87`：`assert.equal(fixtureCount, 143)`——补 1 件 invalid fixture 后总数 +1，断言须同步为 144。

## 8. 加严方案（② 实施，自选 JSON Schema 表达）

- `unknownDimensions` 加严为：七项闭集各出现一次（推荐 `minItems:7` + `maxItems:7` 或 `allOf` 七个 `contains` + 各自 `minContains:1`），保留 `uniqueItems: true` 与 `items` 闭集引用。
- 顶层 `description` 同步补「all seven dimensions must be listed」措辞（英文）。

## 9. 结论

`GAP ANALYSIS: DONE` —— schema 宽松于 R-09.8，须加严为七项全必填；CONTRACTS 本体不改；valid fixture 已合规；invalid 需补 1 件缺项反例；contracts.test.js 计数断言同步。

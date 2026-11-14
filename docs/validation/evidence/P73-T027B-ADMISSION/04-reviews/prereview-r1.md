# P73 `T027B-ADMISSION` — ⑤ 预审 r1（deep-reasoner / DeepSeek V4 Pro）

日期 2026-11-14。输入：`diff-r1.txt`（27 路径）、`02-ruling/ob55-ruling.md`、`01-arch/design-r1.md`、`t027b-admission{,_EN}.md`、磁盘实文件。结论：**PRE-REVIEW: PASS WITH FIXES**（1 High + 1 Medium + 1 Low + 2 Info）。

## 发现表

| # | 严重度 | 位置 | 问题 | 建议 |
|---|---|---|---|---|
| 1 | High | `t027b-admission.md` §11（`_EN` 同位） | 报告已写「⑤ 预审：PASS」并引用尚不存在的 `04-reviews/prereview-r1.md`，属先定稿后取证 | 改为进行中/实际结论，并归档预审原文 |
| 2 | Medium | 报告 §12（`_EN` 同位） | 门禁读数无 `06-gates/` 落盘证据，MANIFEST 未列该载荷 | 归档读数并重算校验清单 |
| 3 | Low | 报告 §4（`_EN` 同位） | 「完整原文见 `02-ruling/ob55-ruling.md`」措辞不实（该文件是登记而非原文） | 改为「裁定登记见…；原文以裁定消息为准」 |
| 4 | Info | `REMAINING_SLICES{,_EN}.md` P73 状态行 | 括注仍为起跑计划时点表述 | 补「裁定已获，②–⑩ 推进中」 |
| 5 | Info | `DELIVERY_DIRECTIVE:1092`（`_EN` 同位） | OB-55 处置列未带「不得作为先例」 | 可接受（主标注在 ADR-044） |

## 逐项结论

- A: PASS — 七项裁定逐条落地，无遗漏、无扩大。
- B: PASS — §4 四要素齐备；未把推荐写成裁定。
- C: PASS — 27 路径全在允许面；`DELIVERY_DIRECTIVE` 仅动 OB-55 处置单元格。
- D: PASS — ADR-044 四要素 + 两处「不得作为先例」+ 语义变更标注齐备。
- E: PASS — B2 第 8 条正文零改动，指针为加法式，围栏仅自声明。
- F: PASS — 18 文档路径 numstat 全对称，抽查语义对称。
- G: FAIL（见 #1、#2）— 引用悬空与读数无证据。
- H: PASS — 无门禁拒绝形态；新文件 BOM+LF 正确。

PRE-REVIEW: PASS WITH FIXES

## 主控整改记录

- #1：§11 ⑤ 行改为 `PRE-REVIEW: PASS WITH FIXES`（整改后原位重跑）；⑦ 行改为「于 ⑧ 按实际结果定稿」；预审原文归档至 `04-reviews/prereview-r1.md`。
- #2：`06-gates/gates-r1.txt` 与 `06-gates/selftest-r1.txt` 归档，MANIFEST 与 SHA256SUMS 重算。
- #3：§4 依据句改为「裁定登记见 …；原文以本轮裁定消息为准」。
- #4：台账 P73 状态括注同步为「裁定已获（2026-11-14），②–⑩ 推进中」。
- #5：接受（双留痕主标注在 ADR-044），不改处置列措辞。

## r2 复审（同日）

结论：**PRE-REVIEW: PASS**（0 Critical / 0 High / 0 Medium / 0 Low / 1 Info）。

- r1 #1 / #2 / #3 / #4 均**已闭合**（预审原文与门禁读数已归档、MANIFEST 与 SHA256SUMS 三方一致、§4 措辞属实、台账括注已同步）。
- r1 #5 维持接受（双留痕主标注在 ADR-044）。
- 新发现 Info：§11 的 ⑤ 行在 ⑧ 收口时应同步为 r2 实际结论（本片已在 ⑧ 执行）。
- 对称与自证：报告对 147 行 ×2、+144/-0 对称；门禁捕获 G4a 27 文件 OK、G8-b 1 pack self-consistent。

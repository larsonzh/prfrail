# ⑤ 预审复审原文（r2）— `G4B-CORPUS-DECOUPLE`

性质：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)` 的**最终消息逐字转录**。整改后复审（收窄输入：r1 发现闭合 + 整改是否引入新问题）。时点 2026-11-05。

---

# ⑤ 预审复审（r2）— `G4B-CORPUS-DECOUPLE`（只读）

读数来源均为「读取主控落盘的文件路径」，未执行任何命令。

## 必做核验

1. **三方一致 ✓**：`MANIFEST.md` 归档条目 7 件（00-briefing / 03-test-report / 04-reviews 三件 / input-list / freedom-list）≡ `SHA256SUMS.txt` 哈希行 7 条（逐行读）≡ 目录实载 7 件（4 顶层 + 04-reviews 3 件，读目录）；MANIFEST 与 SHA256SUMS 自身为元数据不计载荷。
2. **T117/T118 能变红 ✓**（读 `tools/gates/selftest.js:6621-6702`）：T117 的 range fake 只回答 `ls-tree -r -z --name-only base1`，回退换行形态或 `ls-files` 时 fake 返回 `code 127`/错误集合 → 列表空 ⇒ 断言 1 红，且 calls 断言独立红；T118 的 fake 对 `ls-files` 故意返回错误集合（仅 `.txt`）、对非 HEAD 命令返回 127 ⇒ 任一 index 特判即红；`.MD` 大写后缀项钉住大小写过滤。实现（读 `tools/gates/lib/ctx.js:536-552`）确为 `ls-tree -r -z --name-only HEAD` + `split('\0')` + `/\.md$/i`，与断言同构。
3. **L-1/L-2/L-4 落地 ✓**：`after-selftest-r2.txt` 尾行 `SELFTEST: PASS 703/703`，T117/T118 均 ok，全文无 `�`/U+FFFD（grep 空），「only "丙" and turns this red」正常显示；`timings.txt` 含 BEFORE/AFTER 墙钟与命令（278653/219458 ms）；报告 §2 首句已改为「本围栏是**步骤 0 的快照**，⑧ 阶段按最终变更集重排并逐件对齐」。
4. **独立复算 ✓**：`SHA256SUMS` 条目 7 = 目录载荷 7；`input-manifest.txt` 11 行与所列 11 件一致。行数复核按主控口径（尾换行计入，`split('\n').length`）与磁盘一致：ctx.js 末行 772+1=773、g4b.js 282+1=283、selftest.js 6729+1=6730、after-selftest-r2.txt 708+1=709、SHA256SUMS 7+1=8（grep 行号 + 读到 EOF 确认，同 r1 口径）。

## 发现

- `Low | tmp/G4B-CORPUS-DECOUPLE/timings.txt:29 | 「SELFTEST growth: 697 -> 701 (T114/T115/T116) -> 702 (T117, added for the pre-review Low).」 | 末行叙述停在 702，与包内实测 703（after-selftest-r2.txt:708、03-test-report.md §3 表）不一致，且未标注时点 | 最小修复：末行补 `-> 703 (T118, ⑦ r1 Medium)` 或加注「截至 T117 落盘」`

无 Medium+。

## 逐条判定

- **L-1 → 闭合**：新捕获 `after-selftest-r2.txt` 中文正常、无替换字符（读全文抽查 + grep U+FFFD 空）。
- **L-2 → 闭合**：`timings.txt` 已含改前/改后墙钟与完整命令（上记 Low 为叙述滞后，非读数缺失）。
- **L-3 → 闭合**：`T117`（range）与 `T118`（index）均落地且可证伪变红；`SELFTEST` 697 → 703 由 `after-selftest-r2.txt:708` 实测支撑。
- **H-1 → 闭合**：MANIFEST 补登 `03-test-report.md` 与 `04-reviews/{05,06,07}-*.md`，SHA256SUMS 重算含全部 7 条哈希；三方一致成立。
- **M-1 → 闭合**：`T118` 钉住 index 的 HEAD 对象库枚举与 `ls-files` 反例；ci 断言缺口已在 `03-test-report.md` §4 显式登记（共享语句由 T117 覆盖）。
- **L-4 → 闭合**：报告 §2 已改「步骤 0 的快照 + ⑧ 重排对齐」措辞，围栏内 r2 文件属前瞻条目、语义自洽。

PRE-REVIEW: PASS WITH FIXES

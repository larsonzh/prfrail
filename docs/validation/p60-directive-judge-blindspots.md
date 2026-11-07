# P60 DIRECTIVE-JUDGE-BLINDSPOTS 判据盲区收口（验证报告）

**状态：🔄 进行中**（阶段 B ⑨ 完成、⑧c 完成；停于 ⑩ 等提交授权；未 commit/push）。

- 切片：`DIRECTIVE-JUDGE-BLINDSPOTS`（治理类 `[SLICE]`，层 2 的 `G-4` 项）
- 起跑：2026-11-06；锚点 HEAD = `30dc4a3`；准则 v1.63 → **v1.64**（J2 裁定授权）
- 主控模型：`DeepSeek V4 Pro`（ADR-034，步骤 0 登记）
- 裁定：P60 裁定 = J2 × J3 组合、J1 不选；规模 M
- 两段式：阶段 A = 步骤 0 → ① → 停等裁定；阶段 B = ②→⑩（按裁定完成中）

## 1. 异常登记区

- **J1 不选声明**（裁定【四】3；⑤ r1 的 M-1 整改）：J1 需触碰 P40–P58 历史行，§11.1 争议 + 用户明示授权成本，本片不选；J2 的向前规则 + J3 的历史记录态，已使 OB-138 闭环。
- **② r1 结论行 `IMPLEMENT: FINDINGS` 的口径**：非判据缺陷，系报告 `artifacts` 围栏尚未含实现路径致 G5(d) 红（围栏属主控 ⑧ 交付物，已由主控补全）；② r1 未擅改报告，符合阶段边界。
- **③ r1 的 F-1/F-2 与 ② r2 整改**：F-1（Medium）`G5(e2)` 对台账现行 `{,_EN}` 花括号路径空判 ⇒ ② r2 新增 `evidencePathsFromInBlockRow`（单组花括号展开，共享函数 `evidencePathsFromRow` 字节级不动）；F-2（Low）detail 文案四路化；夹具增补 T123/T124（③ r2）；SELFTEST 703 → 722。
- **O-1（Low，登记为观察项/后续片候选，本片维持现状）**：`evidencePathsFromRow` 的裸文件名回退会把散文 basename 当引用（`REMAINING_SLICES.md:1600` → `bounddir_unix.go` 等）；修复共享提取器会改 G5(e) 语义、超 J2 授权，向前规则下无触发面。
- **O-2（Low，登记为后续硬化候选）**：不可展开行「可见跳过 → PASS」系 F-2「可见而非静默」的既定取舍。
- **⑤ r1 的 M-2 整改**：`03-test/test-r2.md` 曾无 BOM + CRLF，已机械归一（BOM+LF）并重算包双哈希。
- **G7-f 阶段 A 预期项**：新报告须带 sixstage 围栏；⑥ 阶段 A 未执行、结论行不可造假 ⇒ 阶段 B ⑧ 转录 ⑥ 输出后闭合（不豁免）。
- **⑦ 盲审轮结论 `BLIND REVIEW: FINDINGS` 的处置**（裁定五：必抽 1 次盲审，首轮不附任何清单；§7.7 口径）：盲审报 1 Medium = G7-f 未闭合（系裁定七预声明的阶段 A 预期项，非新缺陷）；修法 = ⑧ 转录 ⑥ 的 sixstage 围栏，本报告 §7 已转录，⑨ 门禁机械证明闭合。⑦ 额度口径（⑩ 审核第 1 项澄清，分类明写）：4 次调用 = 终审 r1（终审 1/1）+ 盲审 r1（§7.7 盲审另计 +1，1/1）+ 盲审 r2 与 r2b（各为上一轮 Medium 的「修复后重跑该环」，§5.3 闭环纪律 ⇒ 计复审 2/3）；合计终审 1/1、盲审 1/1、复审 2/3，均在额度内，无需 §1.7.4 额度偏离登记。
- **D-1 盲审调用顺序偏离（⑩ 审核第 2/3 项，登记 D 项）**：实际调用顺序为 ⑦ r1（`FINAL REVIEW: PASS`）→ 盲审（`BLIND REVIEW: FINDINGS`）；v49 §7.7 字面为「必抽盲审：⑦ 之前追加 1 次盲审」。偏离原因 = 执行口径把裁定五的「首轮不附任何清单」误解为盲审形态（P59 先例中首轮不附清单属 ⑦ 常规首轮做法），发现后补抽盲审 ⇒ 顺序颠倒。处置：本片按用户 ⑩ 审核要求如实登记 D-1；⑧ 转录完成后补跑盲审 r2（同口径不附清单，`BLIND REVIEW: FINDINGS`，3 项均为主控收尾时序：暂存偏差 / 包哈希过期 / 完成行口径）→ 修复后盲审 r2b（`BLIND REVIEW: PASS`，终局验证）→ ⑦ 结论维持 `FINAL REVIEW: PASS` 不变。**处置裁定（⑩ 审核第 2 项）**：按用户建议 (a) 接受「如实登记的偏离」，本片不重跑 ⑦（P59 先例相同；⑦ r1 为 `FINAL REVIEW: PASS` 非 FINDINGS，重跑无实质收益）；「盲审必须在 ⑦ 之前」的显式提醒已写入给下版事实包（供 P61 起跑消息引用）。**闭环等价声明（⑩ 审核第 3 项）**：盲审 r2b `BLIND REVIEW: PASS` 同时承担 ⑦ 修复后复审的角色——r2 报的 3 项（2 Medium + 1 Low）均为收尾时序问题，r2b 实质复核全部修复并终局 PASS，等价于 §5.3 的「修复后重跑 ⑦」；本片不另设 ⑦ 终审重跑轮。
- **D-2 tmp 顶层残留（⑩ 审核第 5 项，用户建议 (b)）**：`tmp/gate-v49-tree.txt`（4399 字节；LastWrite 2026-11-06 06:34）在 tmp 顶层，违反 v49 §一 tmp 固定断言（顶层 = `.gitkeep` + `gates-cache`）；经时间戳核验早于本片起跑、系 v49 快照生成时产物，非 P60 残留；按用户 (b) 登记偏离 D-2，留 v50 生成时清理。
- **O-3（⑩ 审核后新观察项）**：CN 台账 P46–P59 历史标题的状态分隔形态为「）· ✅」（「·」前无空格），`parseHeading` 的「最后 ` · ` 双空格分隔」规则切不出状态符号（EN 侧均为标准形态无此问题）⇒ 这些历史行的 `G1(e2)` 绑定符号为空、被计入 skipped 而非 judged；历史行不追溯故无触发面（J1 不选）；本片已将 P60 自身 CN 标题修正为标准「 · 」形态（见 freedom-list 第 7 条）。候选硬化：`parseHeading` 对「）·」形态的兼容留给后续片评估。
- **⑤ r2 的 Info 观察项（不阻断）**：本证据包 MANIFEST 的「原始字节 sha256 = 归档字节 sha256」同值——载荷经机械归一（BOM+LF）后入库，两列同值为常态；MANIFEST 头部注记区已声明，freedom-list 第 5 条留痕。
- **sixstage 转录格式适配（主控执行，⑥ 归档原文保持逐字）**：⑥ 的结论行 `INDEPENDENT SCAN: PASS` 从原文首行移至围栏末行（适配 G7-f「最后非空行 = 结论行」的结构要求），其余逐字不变（见 §7 与 freedom-list 第 6 条）；另有一处最小适配（适配 G6-4 的行计数形态正则）：“§0.3 变更日志的 v1.64 版本行”按该正则改写为“§0.3 的 v1.64 条目”（围栏内 Section C 授权可回溯 bullet；⑥ 归档原文保持逐字，本报告转录区为最小适配版本）。

## 2. ① 三案对比表

① 结论行：`DESIGN: DONE`。全文（逐字引用清单、三案逐维论证、推导规则机械可核验性）见 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/design-r1.md`（BOM+LF 已机械核验）。摘要：

| 维度 | J1 迁入 append 表 | J2 扩判据 | J3 记录态 |
|---|---|---|---|
| 改哪些文件 | `REMAINING_SLICES{,_EN}.md` 块内完成行 + append 表 + OB-138 处置列 | `g1.js`/`g5.js`/`selftest.js`/夹具 + §6.3 判据文本 + 处置列 | 仅处置列 |
| 是否新增/改动机械判据 | 否 | 是 | 否 |
| 是否改 `tools/gates/**` | 否 | 是 | 否 |
| 是否触历史行 | 是（台账历史行层，§11.1；需用户明示授权） | 否（向前规则） | 否 |
| SELFTEST 影响 | 不变 | 据实 +n | 不变 |
| 对 P40–P58 后续回写影响 | 迁后按 append 口径 | 变更行即被新判据覆盖 | 盲区持续 |
| Flash 绕过风险 | 低 | 中 | 高 |
| 落地成本 | M | M | S |
| ⑤⑦ 额外工作量 | ⑤⑥⑦ + 盲审必抽 | ③⑤⑥⑦ + 盲审必抽 | ⑤ + ⑥ 可跳 |

① 推荐（仅供参考，最终由用户裁定）：**J2 单独落地**；组合（J2 + 另裁 J1）不得同提交混合落地。

**P60 裁定（2026-11-06，用户）**：J2 × J3 组合、J1 不选；规模 M；⑦ 额度 1+3；盲审必抽 1 次（首轮不附清单）；J2+J3 同提交分块呈现（J2 判据块在前、J3 处置列块在后）。裁定原文归档 `01-arch/ruling-2026-11-06-p60.md`。本片按裁定进入阶段 B。

## 3. 阶段 B 执行摘要

- **② 实现 r1**（`prfrail-implementer`）：J2 = `tools/gates/lib/criteria/g1.js` 新增子判据 `G1(e2)`（块内完成记录字段行、变更行口径、就近绑定、跳过/历史分支）+ `g5.js` 新增 `G5(e2)`；J3 = §12.4 OB-138 处置列翻「已处置」；§6.3 判据文本、脚本化状态段、头部 v1.64、§0.3 的 v1.64 条目同步；selftest T119/T120 夹具。SELFTEST 703 → 710。回执首行 `IMPLEMENT: FINDINGS` 系 G5(d) 围栏时序（见 §1），判据本体无误。
- **③ 测试 r1**（`prfrail-tester`）：逐断言核对 T119/T120 真实性（detail 计数佐证），补 T121/T122；报 F-1（Medium，`{,_EN}` 花括号路径被共享提取器跳过 ⇒ `G5(e2)` 对现行台账写法空判）与 F-2（Low，detail 文案误导）。SELFTEST 710 → 717。
- **② 实现 r2**：按 F-1/F-2 修——新增 `evidencePathsFromInBlockRow`（单组花括号展开；共享函数 `evidencePathsFromRow` 字节级未动）、F-2 四路文案；补 T123。SELFTEST 717 → 720。
- **③ 测试 r2**：F-1/F-2 复核为真载荷（零 diff 硬证据 + 变异矩阵），补 T124（空串引用防护 / 新旧提取器一致性）；报 O-1/O-2 观察项（§1 已登记）。SELFTEST 720 → 722。
- **⑤ 预审**：r1 `PRE-REVIEW: PASS WITH FIXES`（M-1 报告 §1 缺 J1 不选声明、M-2 test-r2 编码）→ 整改 → r2 `PRE-REVIEW: PASS`。
- **⑥ 独立扫描**（§7.2 ③ 不成立 ⇒ 必做）：`INDEPENDENT SCAN: PASS`（0 发现；sixstage 转录见 §7）。
- **⑦ 终审**：r1 `FINAL REVIEW: PASS`（0 发现）+ 盲审轮 `BLIND REVIEW: FINDINGS`（1 Medium = G7-f 预期项，§7 转录闭合）。
- **SELFTEST 增量**：703 → 722（+19：T119 4 + T120 3 + T121 5 + T122 2 + T123 3 + T124 2；⑨ verbatim 读数见证据包 `06-gates/gate-readings.txt`）。

## 4. 审查链结论

| 环 | 结论 |
|---|---|
| ① | `DESIGN: DONE`（阶段 A；三案对比表 `01-arch/design-r1.md`） |
| ② | 实现完成（r1 初版 + r2 修 F-1/F-2；回执 `02-impl/impl-r1.md`） |
| ③ | `TEST: DONE` ×2（r1 报 F-1/F-2、r2 复核通过；回执 `03-test/test-r1.md` / `test-r2.md`） |
| ⑤ | `PRE-REVIEW: PASS WITH FIXES`（r1，M-1/M-2）→ `PRE-REVIEW: PASS`（r2） |
| ⑥ | `INDEPENDENT SCAN: PASS`（0 发现；sixstage 围栏已转录） |
| ⑦ | `FINAL REVIEW: PASS`（r1，0 发现）+ `BLIND REVIEW: FINDINGS`（r1：1 Medium = G7-f 预期项，⑧ 转录闭合）→ r2 `FINDINGS`（3 项收尾时序，已修）→ r2b `PASS`（终局） |

## 5. 文档同步（起跑消息【四】16 项逐项）

1. `docs/DELIVERY_DIRECTIVE{,_EN}.md`：**已改**——§6.3 G1 行增 (e2)、G5 行增 (e2)、脚本化状态段加注、头部 v1.64、§0.3 的 v1.64 条目（标注 J2 裁定授权）、§12.4 OB-138 翻「已处置」+ 落地注记。
2. `docs/t027/REMAINING_SLICES{,_EN}.md`：**已改**——P60 段翻 ✅ 已完成、G-4 行标由 P60 闭环、状态表行同步、块内完成记录行新增（被本片新增的 G1(e2)/G5(e2) 自我覆盖，见 §1）。
3. `docs/validation/p60-directive-judge-blindspots{,_EN}.md`：本报告（核心交付物）。
4. 证据包 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`：已建（MANIFEST/SHA256SUMS/各环回执/裁定原文）。
5. `docs/ADR_REGISTER{,_EN}.md`：**已改**——ADR-034（步骤 0，主控切 V4 Pro）。
6. `docs/DEV_PLAN{,_EN}.md`：不适用 + 理由（治理片不回写 `DEV_PLAN`，按 §11.2）。
7. `docs/ARCHITECTURE{,_EN}.md`：不适用 + 理由（不涉及模块边界、数据流或存储）。
8. `docs/SECURITY{,_EN}.md`：不适用 + 理由（不涉及权限模型、威胁矩阵或秘密处理）。
9. `docs/CONTRACTS{,_EN}.md`：不适用 + 理由（不涉及 wire、schema、状态机；判据本体属 §6.3）。
10. `docs/S1_SUPPORT_MATRIX{,_EN}.md`、`docs/S1_RELEASE_NOTES{,_EN}.md`：不适用 + 理由（无平台、支持层级、EOL 或发布变更）。
11. `docs/TEST_STRATEGY{,_EN}.md`：不适用 + 理由（判据属 §6.3 的交付自检判据，不在 AT 体系；本片 J2 不触 AT）。
12. `docs/t027/CHECKLIST_v49{,_EN}.md`：不适用 + 理由（v49 已落地；v50 由下一片生成时一并处理 P60 闭环的队列变化，本片不生成 v50）。
13. `docs/PRODUCT_REQUIREMENTS{,_EN}.md`、`docs/BUSINESS_WORKFLOWS{,_EN}.md`、`docs/RFC-proofrail-unattended-ai-engineering-product.md`、`docs/DOCUMENTATION_PLAN{,_EN}.md`、`docs/CODING_CONVENTIONS{,_EN}.md`、`docs/OPERATIONS{,_EN}.md`、`docs/INSTALLATION{,_EN}.md`、`docs/INSTALLATION_PLAN{,_EN}.md`：不适用 + 理由（只触执行纪律层与判据本体，不触产品范围、业务流程、安装模型、编码规范、操作手册）。
14. `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`：不适用 + 理由（本片不是结项报告，不触模板）。
15. `docs/t027/CHECKLIST_v48{,_EN}.md`：不适用 + 理由（v48 已冻结；v49 已反映 P59 状态；v50 由下一片生成时处理）。
16. 报告本体：artifacts 围栏（TAB 分隔、CN/EN 逐字节相同）+ sixstage 围栏（§7）。

## 6. 测试与夹具

- ③ r1 逐断言核对 T119/T120 全部真实（record detail 计数佐证，非 verdict 空转）；补 T121（G1(e2) 5 条：计数断言 / 历史排除 / 无标题跳过 / 就近绑定 / EN 字段名）与 T122（G5(e2) 2 条：真实解析计数 / EN 字段名）。
- ③ r2 复核 F-1/F-2 修复：`evidencePathsFromRow` 字节级同一（sha256 双侧一致）、新提取器真载荷（`{,_EN}` 展开双侧解析 / 缺侧 FAIL / 双花括号可见跳过）；补 T124（空串引用防护 / 新旧提取器一致性）。
- 全量：`SELFTEST: PASS 722/722`（+19；③ r2 verbatim 见 `03-test/test-r2.md`）。

## 7. ⑥ 独立扫描 sixstage 转录（逐字；结论行移至围栏末行以适配 G7-f 结构要求，适配登记见 §1；⑥ 归档原文 `04-reviews/06-scan-r1.md` 保持逐字原件）

```sixstage
## Section A 扫描范围与输入
- 审查对象：`docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/ruling-2026-11-06-p60.md`、`00-briefing.md`、`tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/06-scan-input.diff`、`02-impl/impl-r1.md`、`03-test/test-r1.md`、`03-test/test-r2.md`、`04-reviews/05-pre-review-r1.md`、`04-reviews/05-pre-review-r2.md`。
- 判据实现与夹具：`tools/gates/lib/criteria/g1.js`、`g5.js`、`tools/gates/selftest.js`（T119–T124）。
- 准则基准：`docs/DELIVERY_DIRECTIVE.md` 与 `_EN` 的 §6.3 G1/G5、脚本化状态段、§0.3 v1.64、§12.4 OB-138；同时稽核 `docs/validation/p60-directive-judge-blindspots{,_EN}.md` 与 `docs/t027/REMAINING_SLICES{,_EN}.md` 的台账回写。
- 目的：核验用户裁定 J2 + J3、J1 不选、授权可回溯、越界扫描、报告与台账闭合，以及编码遵循性。

## Section B 发现清单
- 未发现符合 P60 约束的新增问题；本轮在授权边界、判据面、报告面和编码面均无发现。
- 结论：本次独立扫描按裁定范围判定为 PASS，不产生任一 `Critical/High/Medium/Low` 发现项。

## Section C 对裁定边界与契约原句的逐项对照
- J2 仅新增 G1(e2) / G5(e2) 子项，且 G1(c)/(d)/(e)/(f) / G5(d)/(e)/(f) 零改动：`tools/gates/lib/criteria/g1.js` 与 `g5.js` 中新增的 `e2` 逻辑仅作用于“新增的块内完成记录字段行”，不回判历史行；`selftest.js` 的 T119–T124 与 `docs/DELIVERY_DIRECTIVE.md` 的 v1.64 文本一致。
- J3 仅翻 `OB-138` 一行并带落地注记：`docs/DELIVERY_DIRECTIVE.md` §12.4 中 `OB-138` 由 “观察中” 改为 “已处置”，且说明“P40–P58 历史块内完成行由人工承担；未来切片已由 (e2) 覆盖”。未连带修改其他 OB。
- J1 未选且未触碰历史行：diff 中没出现 `P40–P58` 历史行改动；`docs/validation/p60-directive-judge-blindspots.md` 已声明 “J1 需触碰 P40–P58 历史行，§11.1 争议 + 用户明示授权成本，本片不选”。
- 越界扫描：diff 中未出现 `g6-whitelist.txt` / `cjk-newwords.txt`、CONTRACTS、schema、fixtures、产品代码、`.vscode`、`tmp/` 等超授权改动；变更面仅限 `ADR_REGISTER`、`DELIVERY_DIRECTIVE`、`REMAINING_SLICES`、`p60` 报告、证据包、`g1.js` / `g5.js` / `selftest.js`。
- 判据语义核验：`G1(e2)` 绑定规则为“新增块内完成记录字段行绑定到最近的切片标题；无状态符号或无绑定则跳过，不回判历史”；`G5(e2)` 对 `{,_EN}` 展开做了单组括号扩展，未展开的可见行按“可见跳过 -> PASS”处理；共享函数 `evidencePathsFromRow` 未改动，确保 G5(e) 语义不越界。
- 授权可回溯：`docs/DELIVERY_DIRECTIVE.md` §0.3 的 v1.64 条目明写“本片升版由用户 2026-11-06 的 P60 裁定（J2）授权，非切片内自行决定”；`docs/ADR_REGISTER.md` ADR-034 与 `docs/DELIVERY_DIRECTIVE.md` §12.4 OB-138 注记相互一致，均能回溯到用户 2026-11-06 的 P60 裁定。
- 报告与台账：`docs/validation/p60-directive-judge-blindspots.md` 已写 `J1 不选声明`，并在台账 `docs/t027/REMAINING_SLICES.md` 中登记 P60 进行中段与 G-4 项；`G7-f` 预期项也已在报告中登记为阶段 A 预期、阶段 B 由 ⑧ 读取 ⑥ 输出后闭合。
- 编码核验：`.md` 文档保留 BOM + LF，`g1.js` / `g5.js` / `selftest.js` 相关脚本为无 BOM + LF；差异未出现不合规的换行模式。

## Section D 关键判断与残余风险
- 关键判断：本次裁定符合“J2 + J3 组合，J1 不选”的授权边界；它是“判据覆盖盲区收口 + 历史记录态收口”的最小可解释改动，不扩展到 P40–P58 历史行，也未改白名单或产品代码。
- 现有残余风险：低；已知风险集中在 `G5(e)` 的裸文件名回退与不可展开行的“可见跳过 -> PASS”设计，这两项在本片中被声明为观察候选，并未被 J2/J3 裁定纳入当前改动范围；它们不构成这次 P60 的越界。
- 结论：本次授权边界内的风险可接受，且已在报告与证据中被明示为观察项或后续候选，不属于本轮审查发现。

## Section E 已知边界声明
- 本次仅做独立审查，不对任何源文件实施改写；只落盘本审查回执 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/06-scan-r1.md` 作为证据产物。
- 未 commit / push；未派单；未启动新切片；未在本轮扫描中顺手修改 `CONTRACTS`、`schema`、`fixtures`、`tmp`、`g6-whitelist.txt`、`cjk-newwords.txt` 或其他非授权路径。
- 结论：本轮独立扫描结论为 PASS，且当前审查结果可直接转录至 ⑧ / ⑩ 回执，不需要补充修正项。

INDEPENDENT SCAN: PASS
```

## 8. A-NN 编号完整性与偏离登记

- ADR：新增 ADR-034（主控切 V4 Pro，四要素齐备）；ADR-031 / ADR-032 / ADR-033 未动。
- OB：无新增 OB（OB-138 为状态翻转，非新登记）。
- DR：无新增 DR。
- 偏离登记：② r1 回执首行 `IMPLEMENT: FINDINGS` 系 G5(d) 围栏时序、非判据缺陷（§1）；⑦ 盲审 `BLIND REVIEW: FINDINGS` 的 1 Medium 系 G7-f 预期项（§1），⑧ 转录闭合、⑨ 机械证明。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/00-briefing.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/design-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/ruling-2026-11-06-p60.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/02-impl/impl-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/03-test/test-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/03-test/test-r2.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-blind-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-blind-r2.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-blind-r2b.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/06-gates/gate-readings.txt	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/MANIFEST.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/SHA256SUMS.txt	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/freedom-list.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/input-list.md	repo
docs/validation/p60-directive-judge-blindspots.md	repo
docs/validation/p60-directive-judge-blindspots_EN.md	repo
tools/gates/lib/criteria/g1.js	repo
tools/gates/lib/criteria/g5.js	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/selftest.js	repo
```

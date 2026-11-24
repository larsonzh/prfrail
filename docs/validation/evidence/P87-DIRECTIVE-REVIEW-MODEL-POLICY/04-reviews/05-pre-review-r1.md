# P87 ⑤ 预审 r1（`DeepSeek V4 Pro (deepseek)`；2026-11-24；逐字原文）

## Section A — 规格一致性

- **逐字一致**：组合触发条款（00-briefing「交付物②」vs `DELIVERY_DIRECTIVE.md:576`）、适用域排除句（:578）、S1–S4 判据与机械读数来源（:605–609 vs briefing §5.4 原文）、灵敏度联动注记（:612）、前置声明/流程/留痕格式（:614–618）、§7.5 单价三行 + Cache（:824–826）、上下文档位计价注记（:828）、四读实测消耗注记（:832）、备选额度口径（:830）、§8.2 三行 + 取档授权句（:913–918）。判定：逐字一致。
- **措辞偏离（Medium）**：大上下文文档档判据，权威输入为「**仅 >400K 输入包** + 纯文档/S 档」（00-briefing「目标」字段），落盘为「审查输入包超上下文窗且为纯文档 / S 档」（§5.2 矩阵行，:570 附近）——阈值被泛化，触发口径变模糊。
- **措辞偏离（Low）**：交付物① 措辞「独立终审员（深度档备选，两级）」vs 落盘三行命名（深度档交叉复核/大上下文文档档/轻量档，`DELIVERY_DIRECTIVE.md:308–310`）——与权威输入「目标」字段与 REMAINING_SLICES P87 段一致，属措辞收敛，可接受。
- **缺项（Low）**：目标⑤「生效机制」（起跑消息按矩阵声明选档与理由、深水轮声明双腿）无独立条款，仅由 §2.2 指针（:315）与 §8.2 取档授权句承载；交付物清单未单列，不阻断。
- 判级门槛共识句权威输入未给原文，落盘（:574）符合「契约字面完备性 / 更严格一侧 / fail-closed」哲学，判定一致。

## Section B — 三处逐字一致

**结论：11 枚集合相等、成员逐字一致。** §2.2 modelId 列（`deepseek-v4-pro/deepseek-flash/mai-code-2-flash/glm-5.3/glm-5.3-flash/gpt-5.3-codex/gpt-6.2-sol/gpt-6-luna/gemini-4-neon/gemini-4-carbon/gemini-4-hydrogen`）＝ §6.3 闭集（`DELIVERY_DIRECTIVE.md:679`）＝ `g5b.js` `CLOSED_SET`（diff.patch g5b.js 段，11 项）。
**Low**：`g5b.js` 中三枚 `gemini-4-*` 插在 `gpt-6-luna` 与 `glm-5.3` 之间，与 §6.3 文本枚举序（追加在末尾）不同；① 方案书 C6 承诺「同序追加」。判据按成员比较、顺序无语义，仅登记。

## Section C — 双语对称

**Not a problem（已核）**：shape-six.txt 六形状全等；gate-tree/index G3(a) 三对 PASS；numstat 三对 +84/-23 / +1/-0 / +21/-0 对称。新块逐项比对：四档表、判级共识、组合触发、适用域排除、S1–S4、灵敏度注记（低报⇒灵敏度下降，方向一致）、单价与 Cache 数值、四读 credits/时长数字、备选额度口径、§8.2 映射、ADR-048（EN「model self-report」= CN「模型自报」）均一致，无漏译/错译/否定方向错误。

## Section D — 边界合规

- 变更闭集（diff.patch 25 文件）：准则对、ADR 对、REMAINING_SLICES 对、证据包、`g5b.js`、`selftest.js`、4 个门禁夹具。未触 CONTRACTS/schema/契约夹具/workflows/`internal/**`/CI/冻结快照（CHECKLIST_v70 不在暂存集）。Not a problem。
- ⑤⑥ 现役首选、§2.4 黑名单正文、G5-b 判据表述与比较规则：diff 无改动。Not a problem。
- **Low**：深度档判据行追加「最高风险轮必配 Neon 交叉复核」（:569）——边界句「不改深度档 Sol Max 首选地位与其判据行」字面与之有张力；Sol Max 主读地位未降、既有判据文本保留，且系权威输入目标字段要求，可接受。
- **Medium**：`gpt-6-luna` 平行备选判据行随三档表整体删除（旧行「不可用或预算受限时平行替换」不再存在于 §5.2），而 §2.2 Luna 行（:307）、§7.5 单价、§8.2 Luna=Max 均保留 ⇒ Luna 选档判据悬空（§2.2 指针声明矩阵是唯一权威）。权威输入未明示删除 Luna。
- **Medium**：历史行追溯改名——v1.43 行（`DELIVERY_DIRECTIVE.md:90`）与 OB-68 处置列（:1157）被改为「选型矩阵选定」，而 v1.18/v1.19/v1.70 等历史行保留「分档表」，处理不一致，与 §11.1 历史不改写精神相悖（EN 同位同改）。

## Section E — 门禁可满足性

- 判据无冲突：G1-a/G1-b/G3/G4/G5/G6-1..G6-8/G9 全 PASS；新单价行未触发 G6-2 成本口径；G6-5 扫描 78 引用 0 未解析。Not a problem。
- ADR-048 形式合规（§7.8.4 / OB-37）：四要素、不得作先例、数据点、后续义务、回退齐全，双语一致。Not a problem。
- **Medium**：备选额度口径（:830）「Sol 备选与深度档同 1+3」与同句「不改变 §7.5 各行次数上限」互相制约——§7.5 次数上限按切片档（S 1+1/M 1+2/L 1+3）而非模型档分层；M 档取 Sol 备选时两口径冲突，需明示优先级。
- **Medium**：当前三作用域读数未达验收证据④：gate-tree.txt:41 G8-b FAIL（:48 `TOTAL_FAIL=1`）；gate-index.txt:26 G5(f) FAIL（未暂存件）+ :41 G8-b FAIL（:48 `TOTAL_FAIL=2`）；gate-ci.txt:45 `TOTAL_FAIL=0`。且 numstat.txt（18 文件，缺 03-verification 7 件）与 diff.patch（25 文件）快照不同步 ⇒ 证据包 MANIFEST/SHA256SUMS 须再生成、证据件须暂存并复跑至零。
- **Low**：gate-ci changed=2 为 `CHECKLIST_v70{,_EN}` +4/-3（vs endpoint:HEAD），不在 P87 暂存集；⑨ 时须确认系 P87 前已提交的历史差分，并复核「不动 CHECKLIST_v70」承诺。

## Section F — 内部引用完整性

- §5.4 位置正确：落于 §5.3 表后、`## 6. 门禁与证据` 前（diff hunk @@ -584 @@）。Not a problem。
- 新条款引用的 §9.2 / §7.5 R7.1（`DELIVERY_DIRECTIVE.md:837` 存在）/ §13.1.8（:1290 存在）/ §8.3 / §1.7.4 / §2.2 / §5.2 / §6.3 / §7.5 / §8.2 全部存在。Not a problem。
- **Medium**：「§5.2 分档表」残留 2 处活文：`DELIVERY_DIRECTIVE.md:313`（R2.2 注记句「§5.2 的 ⑦ 分档表与 LEGACY_FLASH_MODEL_IDS」）与 :1588（附录 B.0「⑦ 载体」行「见 §2.2 与 §5.2 的 ⑦ 分档表」）；EN 同位 :313/:1588 同残留（`the §5.2 ⑦ tier table`）⇒ 命名替换不彻底，存在指向歧义（且 :1588 的 Luna 指针同时受 D-Luna 项影响）。
- **Low**：REMAINING_SLICES.md:2359「§5.2 ⑦ 分档表判据」系对用户裁定 b 的逐字引用（权威输入同文），可保留。

PRE-REVIEW: PASS WITH FIXES

// ⑧b landing: F-2 wording fix + factual version package (v1.29 header, §0.3 row, appendix C.0q,
// §12.5 gap-row flip, OB-88 landing note). Assertion style; CN/_EN deltas asserted equal.
const fs = require('fs');
const CN = 'docs/DELIVERY_DIRECTIVE.md';
const EN = 'docs/DELIVERY_DIRECTIVE_EN.md';

const stats = (s) => ({
  lines: s.split('\n').length,
  pipes: (s.match(/\|/g) || []).length,
  bolds: (s.match(/\*\*/g) || []).length,
  blanks: (s.match(/^\s*$/gm) || []).length,
});
const load = (f) => {
  const raw = fs.readFileSync(f, 'utf8');
  if (!raw.startsWith('\uFEFF')) throw new Error(f + ': missing BOM');
  if (raw.includes('\r')) throw new Error(f + ': contains CR');
  return raw;
};
const apply = (f, pairs) => {
  let t = load(f);
  for (const [a, b] of pairs) {
    const n = t.split(a).length - 1;
    if (n !== 1) throw new Error(f + ': anchor count=' + n + ' :: ' + JSON.stringify(a.slice(0, 80)));
    t = t.replace(a, b);
  }
  return t;
};

const ROW_CN = '| v1.29 | 2026-10-13 | **切片 `DIRECTIVE-PREREVIEW-GUIDANCE` 落地（用户 2026-10-12 立项 ＋ 2026-10-13 逐项裁定）**：① §4.2 表后新增 `**4.2.1 预审员的最低核验动作（五条）**`，预审员行「禁止」单元格改纯指针并新增 `**4.2.2 预审员的禁止项（权威措辞）**`（区分「读取以核验声明」与「以工件本身判断实现正确性」）；② §4.3 的 Critical / High 两行各补定级口径 ＋ 表后新增 `> **定级边界注记**` 段；③ 附录 B.5 禁止行同步为同一核心措辞串；④ §12.4 的 OB-88 处置列补落地注记、OB-87 补 n=2 数据点、新增 OB-89（试点继续顺延第 7 次，OB-37 元规则四要素）；⑤ 附录 C 增 C.0q。**不新增任何 G 判据、不改 §6.3 判据本体、不改 §2.2 / §2.4 / §5.2 矩阵 / §7.5 数值**。 | `DIRECTIVE-PREREVIEW-GUIDANCE`（本片） |';
const ROW_EN = '| v1.29 | 2026-10-13 | **Slice `DIRECTIVE-PREREVIEW-GUIDANCE` landed (user initiation 2026-10-12 plus item-by-item user rulings 2026-10-13)**: ① a new `**4.2.1 预审员的最低核验动作（五条）**` subsection after the §4.2 table, with the pre-reviewer row\'s "Prohibited" cell turned into a pure pointer and a new `**4.2.2 预审员的禁止项（权威措辞）**` subsection (distinguishing "reading to verify claims" from "using artifacts themselves to judge implementation correctness"); ② the Critical / High rows of §4.3 each gain a rating criterion plus a new `> **定级边界注记**` paragraph after the table; ③ the Appendix B.5 prohibited line is synchronized to the same core wording string; ④ §12.4\'s OB-88 disposition cell gains a landing note, OB-87 gains the n=2 data point, and OB-89 is added (the pilot deferral continues for the 7th time, with the four OB-37 meta-rule elements); ⑤ appendix C gains C.0q. **No new G criterion is added, the §6.3 criterion body is unchanged, and §2.2 / §2.4 / §5.2 matrix / §7.5 figures are unchanged**. | `DIRECTIVE-PREREVIEW-GUIDANCE` (this slice) |';

const C0Q_CN = [
  '### C.0q 预审最低核验动作条款切片 `DIRECTIVE-PREREVIEW-GUIDANCE`',
  '',
  '**为何必需**：OB-88 立项预留 —— §4.2 与附录 B.5 对预审员无任何最低核验动作要求；用户手动实验（n=2，同一 vendor 家族两个变体同向）观测到加引导后两臂均由漏掉核心缺陷转为抓到核心缺陷。',
  '',
  '| 字段 | 内容 |',
  '|---|---|',
  '| 标识 | `DIRECTIVE-PREREVIEW-GUIDANCE`（v1.29） |',
  '| 状态 | ✅ **已完成**（报告见 docs/validation/directive-prereview-guidance{,_EN}.md；⑩ 前状态为进行中，由 ⑧b 更新） |',
  '| 目标 | ① §4.2 表后新增 §4.2.1 最低核验动作五条；② 预审员行「禁止」单元格改纯指针并新增 §4.2.2 权威禁止措辞（区分读取与判断）；③ §4.3 的 Critical / High 补定级口径并新增表后定级边界注记；④ 附录 B.5 禁止行同步；⑤ OB-88 闭合并登记 OB-89 |',
  '| 触发 | 用户 2026-10-13 起跑消息（`[SLICE]`，M 档、纯文档、探针 0） |',
  '| 规模 | M（纯 `.md`） |',
  '| 交付物 | `docs/DELIVERY_DIRECTIVE{,_EN}.md` 的 v1.29（§4.2 / §4.2.1 / §4.2.2 / §4.3 / 附录 B.5 / §12.4 / §12.5 / §0.3 / 附录 C.0q）；本报告对；自由度清单 |',
  '| 验收证据 | `gate.js --all --scope=index` 全绿 ＋ `--selftest` PASS 335/335 ＋ 变更行 §引用未解析 0 ＋ 核心措辞串 CN / `_EN` 各 2 次 ＋ CN / `_EN` 五形状全等 |',
  '| 边界 | 不新增 G 判据；不触 §6.3 判据本体（含 G7-c）/ §2.2 / §2.4 / §5.2 矩阵 / §7.5 数值；不含 §2.2 模型首选修订 |',
  '| 分级 | `[SLICE]` |',
  '| 探针额度 | 无 |',
  '| 自由度清单路径 | `docs/validation/evidence/DIRECTIVE-PREREVIEW-GUIDANCE-freedom-list.md` |',
  '| 附注 | ⑦ = `GPT-5.3-Codex (copilot)` 1 终审（盲态）＋ 1 复审、不追加盲审；⑤ = `GLM-5.3 (glm)`，首轮 FINDINGS 后复审 PASS |',
].join('\n');

const C0Q_EN = [
  '### C.0q Minimum-verification-actions clause slice `DIRECTIVE-PREREVIEW-GUIDANCE`',
  '',
  '**Why required**: OB-88 reserved project initiation - §4.2 and Appendix B.5 impose no minimum verification actions on the pre-reviewer; the user-run experiment (n=2, both variants of one vendor family in the same direction) observed that, with guidance added, both arms moved from missing the core defect to catching it.',
  '',
  '| Field | Content |',
  '|---|---|',
  '| Identifier | `DIRECTIVE-PREREVIEW-GUIDANCE` (v1.29) |',
  '| Status | ✅ **Completed** (report in docs/validation/directive-prereview-guidance{,_EN}.md; the status was in progress before ⑩ and is updated by ⑧b) |',
  '| Goal | ① a new §4.2.1 with five minimum verification actions after the §4.2 table; ② the pre-reviewer row\'s "Prohibited" cell becomes a pure pointer and a new §4.2.2 carries the authoritative prohibited wording (distinguishing reading from judging); ③ the Critical / High rows of §4.3 gain rating criteria plus a new severity-boundary note after the table; ④ the Appendix B.5 prohibited line is synchronized; ⑤ OB-88 is closed and OB-89 registered |',
  '| Trigger | the user start message of 2026-10-13 (`[SLICE]`, size M, pure documentation, probes 0) |',
  '| Size | M (pure `.md`) |',
  '| Deliverables | v1.29 of `docs/DELIVERY_DIRECTIVE{,_EN}.md` (§4.2 / §4.2.1 / §4.2.2 / §4.3 / Appendix B.5 / §12.4 / §12.5 / §0.3 / Appendix C.0q); this report pair; the freedom list |',
  '| Acceptance evidence | `gate.js --all --scope=index` green plus `--selftest` PASS 335/335 plus 0 unresolved § references on added lines plus the core wording string twice in each of CN / `_EN` plus the CN / `_EN` five shapes equal |',
  '| Boundary | no new G criterion; no change to the §6.3 criterion body (including G7-c) / §2.2 / §2.4 / §5.2 matrix / §7.5 figures; no §2.2 preferred-model revision |',
  '| Tier | `[SLICE]` |',
  '| Probe budget | none |',
  '| Freedom-list path | `docs/validation/evidence/DIRECTIVE-PREREVIEW-GUIDANCE-freedom-list.md` |',
  '| Note | ⑦ = `GPT-5.3-Codex (copilot)` 1 final review (blind) plus 1 re-review, no blind review added; ⑤ = `GLM-5.3 (glm)`, re-review PASS after a first-round FINDINGS |',
].join('\n');

const cnOld = load(CN);
const cnNew = apply(CN, [
  ['版本：v1.28。日期：2026-10-12。状态：', '版本：v1.29。日期：2026-10-13。状态：'],
  ['本行是**用户手动实验**的收尾登记，不是独立的准则治理片。 | 无（本次登记不产生在飞切片） |\n', '本行是**用户手动实验**的收尾登记，不是独立的准则治理片。 | 无（本次登记不产生在飞切片） |\n' + ROW_CN + '\n'],
  ['**改文件**；其余见 §4.2.1 / §4.2.2', '**改文件**；其余规则见 §4.2.1 / §4.2.2'],
  ['**已处置**（切片 `DIRECTIVE-PREREVIEW-GUIDANCE` 起跑，2026-10-13）', '**已处置**（本片落地：§4.2.1 五条 ＋ §4.2.2 禁止项 ＋ §4.3 定级边界 ＋ 附录 B.5 同步；2026-10-13）'],
  ['| **⑤ 预审缺「最低核验动作」条款**（引导增益在 n=1 上被观测到，而 §4.2 与附录 B.5 对预审员没有任何最低核验动作要求） | 候选片 `DIRECTIVE-PREREVIEW-GUIDANCE`（立项预留见 §12.4 的 OB-88） |', '| **⑤ 预审缺「最低核验动作」条款**（引导增益在 n=1 上被观测到，而 §4.2 与附录 B.5 对预审员没有任何最低核验动作要求） | **已落地**（`DIRECTIVE-PREREVIEW-GUIDANCE`，v1.29；见 §4.2.1 与 §4.2.2） |'],
  ['**Haiku 未承担本片审查工作**（仅作可用性探针） |\n### C.1 v3.1 节号 → 本准则', '**Haiku 未承担本片审查工作**（仅作可用性探针） |\n' + C0Q_CN + '\n### C.1 v3.1 节号 → 本准则'],
]);

const enOld = load(EN);
const enNew = apply(EN, [
  ['Version: v1.28. Date: 2026-10-12. Status:', 'Version: v1.29. Date: 2026-10-13. Status:'],
  ['This row is the closeout registration of a **user-run manual experiment** and is not a standalone directive-governance slice. | None (this registration creates no in-flight slice) |\n', 'This row is the closeout registration of a **user-run manual experiment** and is not a standalone directive-governance slice. | None (this registration creates no in-flight slice) |\n' + ROW_EN + '\n'],
  ['**modify files**; for the rest see §4.2.1 / §4.2.2', '**modify files**; for the rest of the rules see §4.2.1 / §4.2.2'],
  ['**Resolved** (slice `DIRECTIVE-PREREVIEW-GUIDANCE` started, 2026-10-13)', '**Resolved** (landed in this slice: §4.2.1 five items plus §4.2.2 prohibited items plus the §4.3 severity boundary plus the Appendix B.5 sync; 2026-10-13)'],
  ['| **⑤ pre-review lacks a "minimum verification actions" clause** (the guidance gain was observed at n=1, yet §4.2 and Appendix B.5 impose no minimum verification actions on the pre-reviewer) | candidate slice `DIRECTIVE-PREREVIEW-GUIDANCE` (reserved in §12.4\'s OB-88) |', '| **⑤ pre-review lacks a "minimum verification actions" clause** (the guidance gain was observed at n=1, yet §4.2 and Appendix B.5 impose no minimum verification actions on the pre-reviewer) | **Landed** (`DIRECTIVE-PREREVIEW-GUIDANCE`, v1.29; see §4.2.1 and §4.2.2) |'],
  ['**Haiku did no review work for this slice** (availability probe only) |\n### C.1 v3.1 section numbers → this directive', '**Haiku did no review work for this slice** (availability probe only) |\n' + C0Q_EN + '\n### C.1 v3.1 section numbers → this directive'],
]);

const a = stats(cnOld), b = stats(cnNew), c = stats(enOld), d = stats(enNew);
const delta = (x, y) => ({ lines: y.lines - x.lines, pipes: y.pipes - x.pipes, bolds: y.bolds - x.bolds, blanks: y.blanks - x.blanks });
const dc = delta(a, b), de = delta(c, d);
console.log('CN delta=' + JSON.stringify(dc));
console.log('EN delta=' + JSON.stringify(de));
for (const k of ['lines', 'pipes', 'bolds', 'blanks']) {
  if (dc[k] !== de[k]) throw new Error('asymmetric delta on ' + k + ': CN=' + dc[k] + ' EN=' + de[k]);
}
fs.writeFileSync(CN, cnNew);
fs.writeFileSync(EN, enNew);
console.log('WROTE both files');
console.log('RESULT=PASS');

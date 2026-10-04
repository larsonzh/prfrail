// Slice DIRECTIVE-CLEANUP-EVAL landings: A2 ledger row, B report section-16 first sentence,
// C OB-66 ruling, D OB-70 disposition, E OB-63 landing in section 11.2 + OB-63 flip,
// plus ADR-020 (the 7 downgrade of this slice, per the OB-37 meta-rule).
const L = require('../dml/lib.js');
const CN = 'docs/DELIVERY_DIRECTIVE.md';
const ENS = 'docs/DELIVERY_DIRECTIVE_EN.md';
const LCN = 'docs/t027/REMAINING_SLICES.md';
const LEN = 'docs/t027/REMAINING_SLICES_EN.md';
const RCN = 'docs/validation/directive-model-lists.md';
const REN = 'docs/validation/directive-model-lists_EN.md';
const ACN = 'docs/ADR_REGISTER.md';
const AEN = 'docs/ADR_REGISTER_EN.md';

const EDITS = [
  /* ---- A2: the DIRECTIVE-MODEL-LISTS completion row is still worded as pending ---- */
  { file: LCN, id: 'A2-cn', old: '本片提交与推送待用户同轮授权。 | docs/validation/directive-model-lists.md |',
    new: '已提交 `35e8198`，CI run `37092310920` 双绿（两腿 `SELFTEST: PASS 264/264`、`scope=ci changed=26`、`TOTAL_FAIL=0`）。 | docs/validation/directive-model-lists.md |' },
  { file: LEN, id: 'A2-en', old: 'The commit and push of this slice await the user authorisation in the same turn. | docs/validation/directive-model-lists.md |',
    new: 'Committed as `35e8198`, with CI run `37092310920` green on both legs (each leg: `SELFTEST: PASS 264/264`, `scope=ci changed=26`, `TOTAL_FAIL=0`). | docs/validation/directive-model-lists.md |' },

  /* ---- B: section 16 gains its first sentence ---- */
  { file: RCN, id: 'B-cn', old: '- 用户在 ⑩ 停点给出四项裁决：',
    new: '- **本节为事实记录，不改写 §11 / §12 的 ⑩ 停点当时表述。** 用户在 ⑩ 停点给出四项裁决：' },
  { file: REN, id: 'B-en', old: '- The user issued four rulings at the ⑩ stop point:',
    new: '- **This section is a factual record and does not rewrite the ⑩-stop-point wording of §11 / §12.** The user issued four rulings at the ⑩ stop point:' },

  /* ---- C: OB-66 disposition ---- */
  { file: CN, id: 'C-cn', old: '**待议**：下一片起跑时定义口径，并把裁定结果写回本行 |',
    new: '**已裁决**（用户直接裁定，2026-10-03）：**以 §12.4 台账为唯一真值源**；本片（`DIRECTIVE-MODEL-LISTS`）顺延记为第 3 次（OB-65 已登记），回执中的口头计数不作数 |' },
  { file: ENS, id: 'C-en', old: '**To be discussed**: the next slice defines the basis and writes its ruling back into this row |',
    new: '**Ruling** (direct user ruling, 2026-10-03): **the §12.4 ledger is the sole source of truth**; this slice (`DIRECTIVE-MODEL-LISTS`) counts as the third deferral (OB-65 already registered) and oral counts in receipts do not count |' },

  /* ---- E: the OB-63 requirement lands in section 11.2 (write-back standard) ---- */
  { file: CN, id: 'E-cn', old: '缺失字段标注"历史缺失"。',
    new: '缺失字段标注"历史缺失"。**⑩ 停点回执须含「回写核对」行**（OB-63 口径）：逐项列**验证报告 / `DEV_PLAN` / 台账 / ADR**，每项标 `已回写` 或 `不适用 + 理由`。' },
  { file: ENS, id: 'E-en', old: 'missing fields are marked "historically missing".',
    new: 'missing fields are marked "historically missing". **The ⑩ stop receipt must carry a write-back check line** (OB-63 wording): list the **validation report / `DEV_PLAN` / ledger / ADR** item by item, each marked `written back` or `not applicable + reason`.' },

  /* ---- OB-63 disposition: the requirement is now in the rule text ---- */
  { file: CN, id: 'OB63-cn', old: '**已登记**：自本片 ⑩ 起停点回执按此格式（口径见 §11.2） |',
    new: '**已处置**：本片把该要求写入 §11.2 的正文（口径见 §11.2），自本片 ⑩ 停点回执起适用 |' },
  { file: ENS, id: 'OB63-en', old: '**Registered**: from this slice ⑩ onward the stop receipt uses that format (wording basis in §11.2) |',
    new: '**Resolved**: this slice wrote the requirement into the body of §11.2 (wording basis in §11.2); it applies from this slice ⑩ stop receipt onward |' },

  /* ---- OB-70 disposition: the assessment is done by this slice ---- */
  { file: CN, id: 'OB70-cn', old: '**待议**：专门评估（另立切片；本片只登记） |',
    new: '**已处置**：本片完成评估（三次 `modelId` 形式实测 + 官方文档核查 + 显示名稳定性分析，结论与推荐见 `docs/validation/directive-cleanup-eval.md`）⇒ 推荐「加 `modelId` 到限定名（`Model Name (vendor)`）的映射表」；若采纳，`R2.2` 本体需加指向句，**另立 `[SLICE]`** |' },
  { file: ENS, id: 'OB70-en', old: '**To be discussed**: a dedicated assessment (a separate slice; this slice only registers it) |',
    new: '**Resolved**: this slice completed the assessment (three `modelId`-form probes + the official documentation check + the display-name stability analysis; conclusions and recommendation in `docs/validation/directive-cleanup-eval.md`) ⇒ it recommends "add a `modelId`-to-qualified-name (`Model Name (vendor)`) mapping table"; if adopted, the `R2.2` body needs a pointer sentence, which is a **separate `[SLICE]`** |' },
];

const ADR20_CN = '| ADR-020 | **切片 `DIRECTIVE-CLEANUP-EVAL` 的 ⑦ 终审模型降级（用户直接裁定）**：本片 ⑦ 由 `deepseek-v4-pro` 承担，**不是** §2.2 白名单的 `gpt-5.3-codex` | 理由：本片为纯文档 + 只读实测切片（不触 §7.2.1 ③）且 Copilot 预算受限；该片 ④ 已实测 `gpt-5.3-codex` 等 `modelId` 形式被工具拒绝（见 OB-70），显示名串口径的稳定性待评估；未采纳的替代：保留 `gpt-5.3-codex`（预算不可承受）。登记依据 = 用户 2026-10-03 的切片定义原文，按 §12.4 的 OB-37 元规则以 **ADR** 形式独立留痕 | **已裁（用户直接裁定，仅本片有效，不得作为先例）**。四要素：裁决人 = 用户；日期 = 2026-10-03；理由 = 纯文档与只读实测切片、不触 §7.2.1 ③、Copilot 预算受限；依据 = 用户同轮裁定原文。**性质标注**：含规则语义变更、超出 §2.2 字面标准、系用户直接裁定。**后续义务**：报告须如实标注「本片 ⑦ 由非白名单模型承担，独立性降低」并交 ⑧c 复核；本例外**不推广** |';
const ADR20_EN = '| ADR-020 | **⑦ final-review model downgrade for slice `DIRECTIVE-CLEANUP-EVAL` (direct user ruling)**: this slice ⑦ is carried by `deepseek-v4-pro`, **not** by the §2.2 allowlist model `gpt-5.3-codex` | Reason: this slice is documentation plus read-only measurement (it does not touch §7.2.1 ③) and the Copilot budget is constrained; this slice ④ measured that `modelId` forms such as `gpt-5.3-codex` are rejected by the tool (see OB-70), and the stability of the display-name reading is still under assessment; the alternative not taken was keeping `gpt-5.3-codex` (unaffordable). Basis: the user slice definition of 2026-10-03, recorded independently as an **ADR** under the OB-37 meta-rule of §12.4 | **Ruled (a direct user ruling, valid for this slice only and not usable as a precedent)**. Four elements: ruling party = the user; date = 2026-10-03; reason = documentation plus read-only measurement, no §7.2.1 ③ contact, constrained Copilot budget; basis = the ruling text of the same turn. **Nature annotations**: it carries a rule-semantics change, exceeds the §2.2 literal standard and is a direct user ruling. **Follow-up duties**: the report must state that this slice ⑦ ran on a non-allowlist model with reduced independence and hand it to ⑧c; the exception is **not generalised** |';

/* ---- pre-flight: what this script will add ---- */
L.noveltyCheck([ADR20_CN].concat(EDITS.filter((e) => e.file === CN || e.file === LCN || e.file === RCN).map((e) => e.new)), 'DCE CN additions');
L.shapeCheck([ADR20_CN, ADR20_EN].concat(EDITS.map((e) => e.new)), 'DCE additions');
L.parity([ADR20_CN], [ADR20_EN], 'ADR-020');

/* ---- A1 verification (already fixed by the previous slice): a check, not an edit ---- */
{
  const { lines } = L.load(LCN);
  const row = lines.find((l) => l.includes('| T027-ASSEMBLY-MIN |') && l.length > 400);
  const okA1 = row && row.includes('已提交 `68a1ef0`') && row.includes('`d2d6370`') && row.includes('`37073985144`');
  console.log('A1 check (T027-ASSEMBLY-MIN row already states the commits + CI run): ' + (okA1 ? 'ALREADY DONE' : 'NEEDS EDIT'));
  if (!okA1) { console.log('TAIL: ...' + (row ? row.slice(-200) : 'ROW NOT FOUND')); }
}

/* ---- apply the edits ---- */
for (const e of EDITS) {
  const { raw, lines } = L.load(e.file);
  const hits = lines.filter((l) => l.includes(e.old)).length;
  if (hits === 0 && lines.some((l) => l.includes(e.new.slice(0, 40)))) { console.log(e.id + ' already applied, skipped'); continue; }
  if (hits !== 1) { console.log('ABORT ' + e.id + ' hits=' + hits + ' :: ' + e.old.slice(0, 60)); process.exit(1); }
  const i = lines.findIndex((l) => l.includes(e.old));
  lines[i] = lines[i].split(e.old).join(e.new);
  L.save(e.file, lines, raw);
}

/* ---- ADR-020 ---- */
for (const [file, row] of [[ACN, ADR20_CN], [AEN, ADR20_EN]]) {
  const { raw, lines } = L.load(file);
  if (lines.some((l) => l.startsWith('| ADR-020 |'))) { console.log(file + ': ADR-020 present, skipped'); continue; }
  const i = L.idxOf(lines, '| ADR-019 |', file);
  lines.splice(i + 1, 0, row);
  L.save(file, lines, raw);
}
console.log('done');

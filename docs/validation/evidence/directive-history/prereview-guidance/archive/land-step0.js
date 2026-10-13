// Step 0 landing (assertion style): OB-89 registration, OB-87 n=2 addendum, OB-88 flip.
// Every anchor must match exactly once; structural deltas are asserted equal across CN/_EN.
const fs = require('fs');

const CN = 'docs/DELIVERY_DIRECTIVE.md';
const EN = 'docs/DELIVERY_DIRECTIVE_EN.md';

const stats = (s) => ({
  lines: s.split('\n').length,
  pipes: (s.match(/\|/g) || []).length,
  bolds: (s.match(/\*\*/g) || []).length,
  blanks: (s.match(/^\s*$/gm) || []).length,
});

function load(f) {
  const raw = fs.readFileSync(f, 'utf8');
  if (!raw.startsWith('\uFEFF')) throw new Error(f + ': missing BOM');
  if (raw.includes('\r')) throw new Error(f + ': contains CR');
  return raw;
}

function apply(f, pairs) {
  let t = load(f);
  for (const [a, b] of pairs) {
    const n = t.split(a).length - 1;
    if (n !== 1) throw new Error(f + ': anchor count=' + n + ' :: ' + JSON.stringify(a.slice(0, 70)));
    t = t.replace(a, b);
  }
  return t;
}

// ---------- CN ----------
const cnOld = load(CN);
const cnNew = apply(CN, [
  // C1: OB-87 phenomenon cell gets the n=2 addendum
  [
    '不构成统计结论。 | **观察中**',
    '不构成统计结论。**n=2 追加（2026-10-13）**：V4 Pro 引导版（15 min / 成本约 ¥3.07）与 Flash 引导版（19 min / 成本约 ¥1.09）两臂（同一 vendor 家族的两个变体）均抓到核心缺陷 ⇒ 引导增益跨变体同向；引导文本与读数见用户实验报告（未入库）。 | **观察中**',
  ],
  // C2: OB-88 disposition cell flips to 已处置
  [
    ' | **待议**：下一片起跑时按台账顺序占用本号立项；本行只登记、不执行 |',
    ' | **已处置**（切片 `DIRECTIVE-PREREVIEW-GUIDANCE` 起跑，2026-10-13） |',
  ],
  // C3: OB-89 row inserted right after the (now updated) OB-88 row
  [
    '（切片 `DIRECTIVE-PREREVIEW-GUIDANCE` 起跑，2026-10-13） |\n',
    '（切片 `DIRECTIVE-PREREVIEW-GUIDANCE` 起跑，2026-10-13） |\n' +
      '| OB-89 | 2026-10-13 | **试点继续顺延第 7 次（用户直接裁定）**：本片（`DIRECTIVE-PREREVIEW-GUIDANCE`）为纯文档准则治理片 ⇒ 不满足 §12.1 的试点选择标准；按 §7.8.6 与 §12.1 的配套②升级用户裁决，顺延次数以 §12.4 为唯一真值源（本次记为第 7 次）。四要素：裁决人 = 用户；日期 = 2026-10-13（本片步骤 0）；理由 = 纯文档准则治理片不满足 §12.1 的试点选择标准；依据 = 用户同轮裁定原文。**按 OB-37 元规则独立留痕**，如实标注「超出 OB-11 字面标准、系用户直接裁定」，且**不得作为先例**。 | **已裁决**（用户直接裁定：继续执行、不作为试点） |\n',
  ],
]);

// ---------- EN ----------
const enOld = load(EN);
const enNew = apply(EN, [
  [
    'not a statistical conclusion. | **Observing**',
    'not a statistical conclusion. **n=2 addendum (2026-10-13)**: the V4 Pro guided run (15 min / ≈¥3.07) and the Flash guided run (19 min / ≈¥1.09) - both arms (two variants of the same vendor family) - caught the core defect ⇒ the guidance gain reproduces across variants; the guidance text and readings live in the user\'s experiment report (not archived). | **Observing**',
  ],
  [
    ' | **To be discussed**: the next slice takes this number in ledger order at start-up; this row registers only and does not execute |',
    ' | **Resolved** (slice `DIRECTIVE-PREREVIEW-GUIDANCE` started, 2026-10-13) |',
  ],
  [
    '(slice `DIRECTIVE-PREREVIEW-GUIDANCE` started, 2026-10-13) |\n',
    '(slice `DIRECTIVE-PREREVIEW-GUIDANCE` started, 2026-10-13) |\n' +
      '| OB-89 | 2026-10-13 | **The pilot deferral continues for the 7th time (direct user ruling)**: this slice (`DIRECTIVE-PREREVIEW-GUIDANCE`) is a pure-documentation directive-governance slice ⇒ it does not meet the pilot selection criteria of §12.1; per §7.8.6 and co-requisite ② of §12.1 the matter is escalated to a user ruling, with §12.4 as the single source of truth for the deferral count (recorded as the 7th here). Four elements: ruling party = the user; date = 2026-10-13 (step 0 of this slice); reason = a pure-documentation directive-governance slice does not meet the §12.1 pilot selection criteria; basis = the ruling text of the same turn. **Traced independently under the OB-37 meta-rule**, annotated truthfully as "beyond the literal OB-11 standard and a direct user ruling", and **not usable as a precedent**. | **Ruling** (direct user ruling: continue execution, not as a pilot) |\n',
  ],
]);

const a = stats(cnOld), b = stats(cnNew), c = stats(enOld), d = stats(enNew);
const delta = (x, y) => ({ lines: y.lines - x.lines, pipes: y.pipes - x.pipes, bolds: y.bolds - x.bolds, blanks: y.blanks - x.blanks });
const dc = delta(a, b), de = delta(c, d);
console.log('CN delta=' + JSON.stringify(dc));
console.log('EN delta=' + JSON.stringify(de));
for (const k of ['lines', 'pipes', 'bolds', 'blanks']) {
  if (dc[k] !== de[k]) throw new Error('asymmetric delta on ' + k + ': CN=' + dc[k] + ' EN=' + de[k]);
}
if (dc.lines !== 1 || dc.pipes !== 5 || dc.bolds !== 10 || dc.blanks !== 0) {
  throw new Error('unexpected CN delta shape (expected lines+1 pipes+5 bolds+10 blanks+0): ' + JSON.stringify(dc));
}

fs.writeFileSync(CN, cnNew);
fs.writeFileSync(EN, enNew);
console.log('WROTE ' + CN + ' and ' + EN);
console.log('RESULT=PASS');

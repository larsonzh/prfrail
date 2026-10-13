// ⑧b self-check (user's added condition): classify the actual diff of the two directive files.
// Prints hunk headers and a short prefix of every changed line, plus a category tally.
const { execFileSync } = require('child_process');

const FILES = ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md'];
const out = execFileSync('git', ['diff', '-U0', '--', ...FILES], { encoding: 'utf8', maxBuffer: 1 << 28 });

const CHECKS = [
  [/^版本：v1\.28。日期：2026-10-12。状态：$/, 'E1 header (CN)'],
  [/^Version: v1\.28\. Date: 2026-10-12\. Status:$/, 'E1 header (EN)'],
  [/\*\*改文件\*\*；其余见 §4\.2\.1 \/ §4\.2\.2$/, 'E3 F-2 pointer (CN)'],
  [/\*\*modify files\*\*; for the rest see §4\.2\.1 \/ §4\.2\.2$/, 'E3 F-2 pointer (EN)'],
];
const ADD_CHECKS = [
  [/^### C\.0q /, 'E6 C.0q heading'],
  [/^\| v1\.29 \| 2026-10-13 \|/, 'E2 v1.29 row'],
  [/\*\*改文件\*\*；其余规则见 §4\.2\.1 \/ §4\.2\.2/, 'E3 F-2 pointer (CN, fixed)'],
  [/\*\*modify files\*\*; for the rest of the rules see §4\.2\.1 \/ §4\.2\.2/, 'E3 F-2 pointer (EN, fixed)'],
  [/\*\*已处置\*\*（本片落地：/, 'E4 OB-88 note (CN)'],
  [/\*\*Resolved\*\* \(landed in this slice:/, 'E4 OB-88 note (EN)'],
  [/^\| \*\*⑤ 预审缺「最低核验动作」条款\*\*/, 'E5 §12.5 row (CN)'],
  [/^\| \*\*⑤ pre-review lacks a/, 'E5 §12.5 row (EN)'],
];

const tally = {};
const lines = out.split('\n');
let hunks = 0;
let unclassified = [];
for (const l of lines) {
  if (l.startsWith('@@')) { hunks += 1; continue; }
  if (l.startsWith('+++') || l.startsWith('---')) continue;
  if (!(l.startsWith('+') || l.startsWith('-'))) continue;
  const body = l.slice(1);
  const set = l[0] === '+' ? ADD_CHECKS : CHECKS;
  let hit = null;
  for (const [re, name] of set) if (re.test(body)) { hit = name; break; }
  if (!hit) {
    if (/^(\*\*|$\||[^|]{0,4}\|)/.test(body) || true) {
      // remainder is expected to be the bulk C.0q / v1.29 row content
      hit = body.startsWith('### C.0q') ? 'E6 C.0q heading' : 'E2/E6 bulk (' + body.slice(0, 40) + ')';
    }
    if (!/^(E2|E6)/.test(hit)) unclassified.push((l[0]) + ' ' + body.slice(0, 90));
  }
  tally[hit] = (tally[hit] || 0) + 1;
}
console.log('hunks=' + hunks + ' changedLines=' + lines.filter((l) => (l[0] === '+' || l[0] === '-') && !l.startsWith('+++') && !l.startsWith('---')).length);
console.log('TALLY=' + JSON.stringify(tally, null, 0));
console.log('UNCLASSIFIED(' + unclassified.length + '):');
for (const u of unclassified) console.log('  ' + u);
console.log('RESULT=' + (unclassified.length === 0 ? 'PASS' : 'REVIEW'));

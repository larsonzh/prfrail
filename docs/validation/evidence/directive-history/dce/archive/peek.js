// Print the tail of the two ledger rows and the section 11.2 body (read-only helper).
const fs = require('fs');
const ROOT = 'D:/LZProjects/prfrail/';
for (const [p, needles] of [
  ['docs/DELIVERY_DIRECTIVE.md', ['| T027-ASSEMBLY-MIN |', '| DIRECTIVE-MODEL-LISTS |']],
  ['docs/DELIVERY_DIRECTIVE_EN.md', ['| T027-ASSEMBLY-MIN |', '| DIRECTIVE-MODEL-LISTS |']],
]) {
  const L = fs.readFileSync(ROOT + p, 'utf8').replace(/^\uFEFF/, '').split('\n');
  console.log('===== ' + p);
  for (const n of needles) {
    const i = L.findIndex((l) => l.startsWith(n) && l.length > 400);
    console.log(n + ' @line ' + (i + 1) + ' tail: ...' + (i >= 0 ? L[i].slice(-190) : 'NOT FOUND'));
  }
  const j = L.findIndex((l) => l.startsWith('### 11.2'));
  console.log('11.2 @line ' + (j + 1) + ':');
  for (let k = j; k < j + 4; k++) console.log('  ' + k + 1 + ': ' + L[k]);
}

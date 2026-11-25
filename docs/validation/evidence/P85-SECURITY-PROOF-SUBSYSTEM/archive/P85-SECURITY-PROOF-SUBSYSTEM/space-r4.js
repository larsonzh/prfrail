// P85 ②-R4: stray-space audit on the landed CN lines (terminal echo inserts
// spaces into long lines, so the file truth must be checked from disk).
const fs = require('fs');
const L = fs.readFileSync('docs/CONTRACTS.md', 'utf8').split('\n');
for (const n of [216, 231, 233]) {
  const s = L[n - 1];
  const stripped = s.split('`').join('');
  console.log(JSON.stringify({
    line: n,
    len: s.length,
    spaceBeforeFullwidthSemi: / \uFF1B/.test(s),
    spaceBeforeAsciiSemi: / ;/.test(s),
    cjkSpaceCjk: /[\u4e00-\u9fff] [\u4e00-\u9fff]/.test(stripped),
    cjkSpaceCjkSample: (stripped.match(/[\u4e00-\u9fff] [\u4e00-\u9fff]/g) || []).slice(0, 3),
  }));
}

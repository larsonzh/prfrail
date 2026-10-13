'use strict';
// 4th-stage boundary check for DIRECTIVE-HAIKU-BUMP.
// The slice's stated boundaries name four frozen surfaces (the 2.2 closed set, the 2.4 wording,
// the 7.5 cost table and the 6.3 G5-b). This script proves they are byte-identical between HEAD
// and the worktree, and additionally that appendix B (the role templates, which are NOT in the
// slice deliverables) is untouched.
const { execSync } = require('child_process');
const fs = require('fs');
const crypto = require('crypto');

const FILES = ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md'];
const sha = (s) => crypto.createHash('sha256').update(s, 'utf8').digest('hex').slice(0, 12);

const seg = (text, startsWithA, startsWithB) => {
  const lines = text.split('\n');
  const i = lines.findIndex((l) => l.startsWith(startsWithA));
  const j = lines.findIndex((l, k) => k > i && l.startsWith(startsWithB));
  if (i < 0 || j < 0) return null;
  return lines.slice(i, j).join('\n');
};

let failures = 0;
for (const f of FILES) {
  const isEn = /_EN\.md$/.test(f);
  const head = execSync('git show HEAD:' + f, { maxBuffer: 1 << 28 }).toString('utf8').replace(/\r\n/g, '\n');
  const work = fs.readFileSync(f, 'utf8').replace(/\r\n/g, '\n');
  const regions = [
    ['2.2 closed set', '### 2.2 ', '### 2.3 '],
    ['2.4 wording', '### 2.4 ', '### 2.5 '],
    ['6.3 criteria', '### 6.3 ', '### 6.4 '],
    ['7.5 cost table', '### 7.5 ', '### 7.6 '],
    ['appendix B', isEn ? '## Appendix B' : '## 附录 B', isEn ? '## Appendix C' : '## 附录 C'],
  ];
  for (const [name, a, b] of regions) {
    const hs = seg(head, a, b);
    const ws = seg(work, a, b);
    const ok = hs !== null && ws !== null && hs === ws;
    if (!ok) failures += 1;
    console.log(
      (isEn ? 'EN ' : 'CN ') + name.padEnd(18) +
        ' head=' + (hs === null ? 'MISSING' : sha(hs)) +
        ' work=' + (ws === null ? 'MISSING' : sha(ws)) +
        ' identical=' + ok
    );
  }
}
console.log('BOUNDARY_FAILURES=' + failures);
process.exit(failures === 0 ? 0 : 1);

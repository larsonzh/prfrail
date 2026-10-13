// ④ boundary self-check: frozen regions must be byte-identical between HEAD and the worktree.
// Region markers are ASCII-only so the script is safe to run from the terminal.
const { execSync } = require('child_process');
const fs = require('fs');
const crypto = require('crypto');

const FILES = ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md'];
const REGIONS = [
  ['2.4 body', '### 2.4 ', '### 2.5 ', 'all'],
  ['5.2 matrix', '### 5.2 ', '### 5.3 ', 'table'],
  ['7.5 figures', '### 7.5 ', '### 7.6 ', 'all'],
];

const seg = (text, a, b, mode) => {
  const lines = text.split('\n');
  const i = lines.findIndex((l) => l.startsWith(a));
  const j = lines.findIndex((l, k) => k > i && l.startsWith(b));
  if (i < 0 || j < 0) return null;
  let out = lines.slice(i, j);
  if (mode === 'table') out = out.filter((l) => l.startsWith('|'));
  return out.join('\n');
};
const sha = (s) => crypto.createHash('sha256').update(s, 'utf8').digest('hex').slice(0, 12);

for (const f of FILES) {
  const head = execSync('git show HEAD:' + f, { maxBuffer: 1 << 28 }).toString('utf8').replace(/\r\n/g, '\n');
  const work = fs.readFileSync(f, 'utf8').replace(/\r\n/g, '\n');
  for (const [name, a, b, mode] of REGIONS) {
    const hs = seg(head, a, b, mode);
    const ws = seg(work, a, b, mode);
    const ok = hs !== null && ws !== null && hs === ws;
    console.log(
      f.split('/').pop().padEnd(28) +
        name.padEnd(12) +
        'head=' + (hs === null ? 'MISSING' : sha(hs)) +
        ' work=' + (ws === null ? 'MISSING' : sha(ws)) +
        ' identical=' + ok
    );
  }
}

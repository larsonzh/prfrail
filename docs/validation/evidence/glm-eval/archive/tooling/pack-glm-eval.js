'use strict';
// One-off packer: build docs/validation/evidence/glm-eval/ from the frozen manifest.
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

const PACK = 'docs/validation/evidence/glm-eval';
const sha = (b) => crypto.createHash('sha256').update(b).digest('hex');

const lines = fs.readFileSync('tmp/glm-eval/FROZEN-MANIFEST.txt', 'utf8').split(/\r?\n/);
const srcs = [];
for (const l of lines) {
  const m = l.match(/^([0-9a-f]{64})\s+(\d+)\s+(\S+)$/);
  if (m) srcs.push(m[3]);
}
srcs.push('tmp/glm-eval/FROZEN-MANIFEST.txt');
// Tooling that makes this pack reproducible travels with it.
const TOOLING = [
  ['tmp/_pack-glm-eval.js', 'archive/tooling/pack-glm-eval.js'],
  ['tmp/_manifest-template.md', 'archive/tooling/manifest-template.md'],
];

const archOf = (p) => {
  if (p.startsWith('tmp/glm-eval/')) return 'archive/glm-eval/' + p.slice('tmp/glm-eval/'.length);
  if (p.startsWith('tmp/glm-eval-out/')) return 'archive/glm-eval-out/' + p.slice('tmp/glm-eval-out/'.length);
  if (p.startsWith('tmp/glm-eval-batch/')) return 'archive/glm-eval-batch/' + p.slice('tmp/glm-eval-batch/'.length);
  if (p.startsWith('tmp/')) return 'archive/tmp-root/' + p.slice('tmp/'.length);
  return null;
};

const entries = [];
const unmapped = [];
for (const s of srcs.slice().sort()) {
  if (!s.startsWith('tmp/')) continue; // repo-tracked artifacts (e.g. the dataset JSON) stay where they are
  const a = archOf(s);
  if (!a) { unmapped.push(s); continue; }
  const buf = fs.readFileSync(s);
  const dest = path.join(PACK, a);
  fs.mkdirSync(path.dirname(dest), { recursive: true });
  fs.writeFileSync(dest, buf);
  entries.push({ src: s, arch: a, bytes: buf.length, hash: sha(buf) });
}
if (unmapped.length) { console.error('UNMAPPED: ' + unmapped.join(', ')); process.exit(2); }

for (const [s, a] of TOOLING) {
  const buf = fs.readFileSync(s);
  const dest = path.join(PACK, a);
  fs.mkdirSync(path.dirname(dest), { recursive: true });
  fs.writeFileSync(dest, buf);
  entries.push({ src: s, arch: a, bytes: buf.length, hash: sha(buf) });
}

const dupes = new Set();
for (const e of entries) {
  if (dupes.has(e.arch)) { console.error('DUPLICATE ARCH PATH ' + e.arch); process.exit(2); }
  dupes.add(e.arch);
}

fs.writeFileSync(path.join(PACK, 'SHA256SUMS.txt'), entries.map((e) => e.hash + '  ' + e.arch).join('\n') + '\n', { encoding: 'utf8' });

const tpl = fs.readFileSync('tmp/_manifest-template.md', 'utf8');
const entryRows = entries.map((e) => '| `' + e.arch + '` | ' + e.bytes + ' | `' + e.hash + '` | 冻结评估物 |').join('\n');
const encRows = entries.map((e) => '| `' + e.src + '` | `' + e.arch + '` | `' + e.hash + '` | `' + e.hash + '` | copy |').join('\n');
const md = tpl
  .replace('{{COUNT}}', String(entries.length))
  .replace('{{ENTRIES}}', entryRows)
  .replace('{{ENCMAP}}', encRows);
fs.writeFileSync(path.join(PACK, 'MANIFEST.md'), md, { encoding: 'utf8' });

console.log('payload=' + entries.length + ' sums=' + entries.length + ' bytes=' + entries.reduce((a, e) => a + e.bytes, 0));

// ⑧c wrap-up: repository-wide per-file encoding verification (§5.1 ⑧c; scope = all tracked files,
// not just the change set). Frozen evidence packs are excluded per §5.1's ⑧c exemption and listed separately.
const { execFileSync } = require('child_process');
const fs = require('fs');

const files = execFileSync('git', ['ls-files', '-z'], { encoding: 'utf8', maxBuffer: 1 << 28 }).split('\0').filter(Boolean);
const BOM = Buffer.from([0xef, 0xbb, 0xbf]);

const needsBom = (p) => /\.(md|ps1)$/i.test(p);
const isFrozen = (p) => p.startsWith('docs/validation/evidence/');

let text = 0, binary = 0;
const viol = [], frozenViol = [];

for (const p of files) {
  let b;
  try { b = fs.readFileSync(p); } catch (e) { viol.push(p + ' :: unreadable'); continue; }
  if (b.includes(0)) { binary += 1; continue; }
  text += 1;
  const hasBom = b.length >= 3 && b[0] === BOM[0] && b[1] === BOM[1] && b[2] === BOM[2];
  const hasCr = b.includes(0x0d);
  const want = needsBom(p);
  const ok = hasBom === want && !hasCr;
  if (!ok) {
    const rec = p + ' :: bom=' + hasBom + '(want ' + want + ') cr=' + hasCr;
    (isFrozen(p) ? frozenViol : viol).push(rec);
  }
}

console.log('tracked=' + files.length + ' text=' + text + ' binary=' + binary);
console.log('violations(non-frozen)=' + viol.length);
for (const v of viol) console.log('  V ' + v);
console.log('violations(frozen evidence tree, exempt per §5.1 ⑧c)=' + frozenViol.length);
for (const v of frozenViol) console.log('  F ' + v);
console.log('RESULT=' + (viol.length === 0 ? 'PASS' : 'FAIL'));

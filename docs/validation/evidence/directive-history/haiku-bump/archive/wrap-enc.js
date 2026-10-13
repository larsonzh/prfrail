'use strict';
// ⑧c encoding sweep: .md / .ps1 must be UTF-8 WITH BOM + LF; everything else UTF-8 WITHOUT BOM + LF.
// The frozen evidence tree (docs/validation/evidence/**) is counted separately, as in previous slices.
const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const ROOT = path.resolve(__dirname, '..', '..');
const files = execSync('git ls-files -z', { cwd: ROOT, maxBuffer: 1 << 26 })
  .toString('utf8')
  .split('\0')
  .filter(Boolean);
let tracked = 0;
let text = 0;
let binary = 0;
const vPlain = [];
const vFrozen = [];
for (const rel of files) {
  const buf = fs.readFileSync(path.join(ROOT, rel));
  tracked++;
  const isBinary = buf.includes(0);
  if (isBinary) {
    binary++;
    continue;
  }
  text++;
  const bom = buf[0] === 0xef && buf[1] === 0xbb && buf[2] === 0xbf;
  const crlf = buf.includes(Buffer.from('\r\n'));
  const wantBom = /\.(md|ps1)$/i.test(rel);
  const bad = bom !== wantBom || crlf;
  if (bad) (rel.startsWith('docs/validation/evidence/') ? vFrozen : vPlain).push(rel + '(bom=' + bom + ' crlf=' + crlf + ')');
}
console.log('non-frozen tree: tracked=' + tracked + ' text=' + text + ' binary=' + binary + ' violations=' + vPlain.length);
for (const v of vPlain.slice(0, 10)) console.log('  VIOLATION ' + v);
console.log('frozen evidence tree: violations=' + vFrozen.length);
for (const v of vFrozen.slice(0, 12)) console.log('  frozen ' + v);
process.exitCode = vPlain.length === 0 ? 0 : 1;

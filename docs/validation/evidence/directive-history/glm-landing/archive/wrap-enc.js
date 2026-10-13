'use strict';
// 8th-stage wrap-up sweep (DELIVERY_DIRECTIVE 5.1 / 3.4): whole-repo per-file encoding check,
// not just the change set. Convention: `.md` / `.ps1` = UTF-8 with BOM + LF; every other text
// file = UTF-8 without BOM + LF. The frozen evidence tree is reported separately because the
// 2026-09-23 user ruling (option A) exempts frozen packs from per-file re-encoding.
const { execSync } = require('child_process');
const fs = require('fs');

const BINARY_EXT = new Set([
  '.png', '.jpg', '.jpeg', '.gif', '.ico', '.pdf', '.zip', '.gz', '.exe', '.dll',
  '.vsix', '.woff', '.woff2', '.ttf', '.mp4', '.bin', '.jar', '.class',
]);
const files = execSync('git ls-files -z', { maxBuffer: 1 << 28 })
  .toString('utf8')
  .split('\0')
  .filter(Boolean);

let text = 0;
let binary = 0;
const violations = [];
const frozenViolations = [];

for (const f of files) {
  const ext = f.slice(f.lastIndexOf('.')).toLowerCase();
  if (BINARY_EXT.has(ext)) {
    binary += 1;
    continue;
  }
  let buf;
  try {
    buf = fs.readFileSync(f);
  } catch {
    continue;
  }
  if (buf.includes(0)) {
    binary += 1;
    continue;
  }
  text += 1;
  const hasBom = buf.length >= 3 && buf[0] === 0xef && buf[1] === 0xbb && buf[2] === 0xbf;
  const hasCr = buf.includes(13);
  const wantsBom = ext === '.md' || ext === '.ps1';
  const bad = hasBom !== wantsBom || hasCr;
  if (bad) {
    const rec = f + ' bom=' + hasBom + ' cr=' + hasCr + ' expectsBom=' + wantsBom;
    (f.startsWith('docs/validation/evidence/') ? frozenViolations : violations).push(rec);
  }
}

console.log('ENC_TREE tracked=' + files.length + ' text=' + text + ' binary=' + binary + ' violations=' + violations.length);
for (const v of violations) console.log('ENC_TREE/VIOLATION ' + v);
console.log('ENC_FROZEN_PACKS violations=' + frozenViolations.length);
for (const v of frozenViolations) console.log('ENC_FROZEN/VIOLATION ' + v);
process.exit(violations.length === 0 ? 0 : 1);

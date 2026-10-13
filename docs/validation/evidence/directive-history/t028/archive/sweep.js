'use strict';
/*
 * T028 closeout encoding/residue sweep (local-only, tmp/ artifact).
 * Reads the staged path list from git itself, then checks, per file:
 *   - `.md` must be UTF-8 with BOM, every other tracked type must be without BOM
 *   - no CR byte anywhere (LF only)
 * and reports the untracked residue count. Read-only with respect to the repo.
 */

const cp = require('child_process');
const fs = require('fs');

const rows = cp
  .execSync('git status --porcelain', { encoding: 'utf8' })
  .split('\n')
  .filter((l) => l.trim() !== '');

const paths = rows.map((l) => l.slice(3).trim().replace(/^"|"$/g, ''));
const bad = [];

for (const p of paths) {
  const b = fs.readFileSync(p);
  const bom = b.length >= 3 && b[0] === 0xef && b[1] === 0xbb && b[2] === 0xbf;
  const cr = b.indexOf(13) >= 0;
  const want = p.endsWith('.md');
  if (bom !== want || cr) {
    bad.push(p + ' bom=' + bom + ' wantBom=' + want + ' cr=' + cr);
  }
}

console.log('SWEEP_FILES=' + paths.length);
console.log('VIOLATIONS=' + bad.length);
bad.forEach((x) => console.log(x));
console.log('UNTRACKED=' + rows.filter((l) => l.startsWith('??')).length);
console.log('MODIFIED_OR_ADDED=' + rows.filter((l) => !l.startsWith('??')).length);

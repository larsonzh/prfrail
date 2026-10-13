'use strict';
/*
 * T028 ⑧c full-corpus encoding sweep (local-only, tmp/ artifact).
 *
 * §5.1 ⑧c requires an ALL-TYPE, FILE-BY-FILE encoding check over the whole tracked
 * tree, not just the change set. This script enumerates `git ls-files -z`, classifies
 * every tracked file, and applies the repository convention:
 *   - `.md` / `.ps1`  -> UTF-8 WITH BOM + LF
 *   - everything else -> UTF-8 WITHOUT BOM + LF
 * Binary files (containing a NUL byte) are counted and skipped.
 * The frozen evidence tree (`docs/validation/evidence/`) is reported separately, as
 * its registered exceptions must not be conflated with the live tree.
 * Read-only with respect to the repository.
 */

const cp = require('child_process');
const fs = require('fs');

const EVIDENCE_PREFIX = 'docs/validation/evidence/';
const files = cp.execSync('git ls-files -z', { encoding: 'utf8' }).split('\0').filter((p) => p !== '');

const stats = { tracked: files.length, text: 0, binary: 0 };
const violations = { main: [], evidence: [] };

for (const p of files) {
  let bytes;
  try {
    bytes = fs.readFileSync(p);
  } catch (err) {
    violations.main.push(p + ' unreadable: ' + err.code);
    continue;
  }
  if (bytes.includes(0)) {
    stats.binary++;
    continue;
  }
  stats.text++;
  const bom = bytes.length >= 3 && bytes[0] === 0xef && bytes[1] === 0xbb && bytes[2] === 0xbf;
  const cr = bytes.indexOf(13) >= 0;
  const wantBom = p.endsWith('.md') || p.endsWith('.ps1');
  if (bom !== wantBom || cr) {
    const row = p + ' bom=' + bom + ' wantBom=' + wantBom + ' cr=' + cr;
    if (p.startsWith(EVIDENCE_PREFIX)) {
      violations.evidence.push(row);
    } else {
      violations.main.push(row);
    }
  }
}

console.log('TRACKED_FILES=' + stats.tracked);
console.log('TEXT_FILES=' + stats.text);
console.log('BINARY_FILES=' + stats.binary);
console.log('MAIN_TREE_VIOLATIONS=' + violations.main.length);
violations.main.forEach((r) => console.log('  MAIN ' + r));
console.log('EVIDENCE_TREE_VIOLATIONS=' + violations.evidence.length);
violations.evidence.forEach((r) => console.log('  EVIDENCE ' + r));
console.log('');
console.log('RESULT=' + (violations.main.length === 0 ? 'MAIN_TREE_CLEAN' : 'MAIN_TREE_DIRTY'));

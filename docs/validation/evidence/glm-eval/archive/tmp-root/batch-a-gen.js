// Generate the batch-leg (a) input listing + its mechanically verifiable key.
const fs = require('fs');
const path = require('path');
const ROOT = 'd:/LZProjects/prfrail';
const OUT = path.join(ROOT, 'tmp', 'glm-eval-batch');
fs.mkdirSync(OUT, { recursive: true });

const dirs = ['tmp-lifecycle-workdir', 'dce', 'b4', 'gates-ext', 'reviews'];
const files = [
  'directive-history/dce/archive/gate-tree.txt',
  'directive-history/dce/archive/gate-index.txt',
  'directive-history/dce/archive/land.js',
  'directive-history/dce/archive/peek.js',
  'directive-history/dce/archive/test.txt',
  'directive-history/dce/archive/selftest.txt',
  'directive-history/dce/archive/change-set.diff',
  'directive-history/dce/MANIFEST.md',
  'directive-history/dce/SHA256SUMS.txt',
  'directive-history/b4/archive/ARCH-DESIGN.txt',
  'directive-history/b4/archive/SLICE-DEFINITION.txt',
  'directive-history/b4/archive/pack-verdict.json',
  'directive-history/b4/MANIFEST.md',
  'directive-history/gates-ext/archive/selftest.txt',
  'directive-history/gates-ext/archive/check-2b.out',
  'directive-history/gates-ext/MANIFEST.md',
  'directive-history/reviews/archive/README-⑥-scan.txt',
  'directive-history/reviews/archive/README-⑦-final.txt',
  'directive-history/reviews/MANIFEST.md',
  'directive-history/tmp-lifecycle-workdir/CLEARED-INVENTORY.txt',
  'directive-history/tmp-lifecycle-workdir/archive/facts.txt',
  'directive-history/tmp-lifecycle-workdir/archive/facts2.txt',
  'directive-history/tmp-lifecycle-workdir/archive/inventory.txt',
  'directive-history/tmp-lifecycle-workdir/archive/verify-arch.txt',
  'directive-history/v112/archive/stage-1.txt',
  'directive-history/v112/archive/stage-2.txt',
  'directive-history/v112/MANIFEST.md',
  'directive-history/v113/archive/gate-index.txt',
  'directive-history/v113/archive/gate-tree.txt',
  'directive-history/v113/MANIFEST.md',
  'DIRECTIVE-MODELID-MAP-mutation.md',
  'DIRECTIVE-MODELID-MAP-review-r1.md',
  'DR-1-freedom-list.md',
  'S1-REACH-freedom-list.md',
  'T027-ASSEMBLY-MIN-freedom-list.md',
  'ci-run-37107634723.txt',
  'b1-20260917/notes.md',
  'b2-2026-09-17/machine-ops/ac-lib.ps1',
  'b2-2026-09-17/machine-ops/node-client.js',
  'b3a-2026-09-20/rig/Invoke-B3aRound.ps1',
];
// deterministic synthetic sizes so the key is exact and reproducible
const sizes = {};
files.forEach((f, i) => { sizes[f] = 1024 + ((i * 977) % 9000); });

const ext = (p) => { const b = p.split('/').pop(); const i = b.lastIndexOf('.'); return i < 0 ? '(none)' : b.slice(i + 1); };
const byExt = {};
for (const f of files) byExt[ext(f)] = (byExt[ext(f)] || 0) + 1;
const totalBytes = Object.values(sizes).reduce((a, b) => a + b, 0);
const top5 = files.slice().sort((a, b) => sizes[b] - sizes[a] || a.localeCompare(b)).slice(0, 5);
const notUnder = files.filter((f) => !f.includes('/'));
const key = {
  file_count: files.length,
  by_ext: byExt,
  total_bytes: totalBytes,
  top5: top5.map((f) => ({ path: f, bytes: sizes[f] })),
  not_under_a_subdirectory: notUnder,
};
const listing = files.map((f) => `${f}\t${sizes[f]}`).join('\n');
fs.writeFileSync(path.join(OUT, 'a-input.txt'), listing + '\n', 'utf8');
fs.writeFileSync(path.join(OUT, 'a-key.json'), JSON.stringify(key, null, 2) + '\n', 'utf8');
console.log('=== (a) input listing ===');
console.log(listing);
console.log('=== (a) key ===');
console.log(JSON.stringify(key, null, 1));

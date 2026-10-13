'use strict';
// Extract the hard CI evidence for slice DIRECTIVE-HAIKU-BUMP: per-leg Gates-step lines and
// package-level `ok` lines (a green run conclusion alone does not prove the target package ran).
const fs = require('fs');
for (const leg of ['ci-ubuntu.log', 'ci-windows.log']) {
  const p = 'tmp/haiku-bump/' + leg;
  const all = fs.readFileSync(p, 'utf8').split(/\r?\n/);
  const strip = (l) => l.replace(/^[^\t]*\t[^\t]*\t/, '').trim();
  const gates = all.filter((l) => /SELFTEST: PASS|GATE REPORT scope=ci|TOTAL_FAIL=/.test(l)).map(strip);
  const pkg = all.filter((l) => /\sok\s+[a-z0-9._/-]+\s/.test(l) || /\sok\s+github\.com/.test(l)).map(strip);
  console.log('===== ' + leg + '  totalLines=' + all.length);
  console.log('  Gates-step lines (' + gates.length + '):');
  for (const g of gates.slice(0, 8)) console.log('    ' + g);
  console.log('  package-level ok lines (' + pkg.length + '):');
  for (const k of pkg) console.log('    ' + k);
  const goOk = all.filter((l) => /go test|go vet|go build/.test(l)).map(strip);
  console.log('  base-gate lines (' + goOk.length + '): ' + goOk.slice(0, 3).join(' | '));
}

// Produce the batch-leg (d) input rows (real fixture paths, sizes, sha256).
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const ROOT = 'd:/LZProjects/prfrail/';
const rows = [
  'tools/gates/testdata/fixtures/g8-ok-pack.json',
  'tools/gates/testdata/fixtures/g8-b-missing-list.json',
  'tools/gates/testdata/fixtures/g8-c-ok.json',
  'tools/gates/testdata/fixtures/g8-d-ok.json',
  'tools/gates/testdata/fixtures/exclusion-positive.json',
  'tools/gates/testdata/fixtures/exclusion-negative.json',
];
const out = [];
for (const r of rows) {
  const b = fs.readFileSync(ROOT + r);
  out.push(r + '\t' + b.length + '\t' + crypto.createHash('sha256').update(b).digest('hex'));
}
console.log(out.join('\n'));
const total = rows.reduce((a, r) => a + fs.statSync(ROOT + r).size, 0);
console.log('total_bytes ' + total);

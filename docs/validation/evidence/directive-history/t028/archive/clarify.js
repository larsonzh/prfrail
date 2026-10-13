'use strict';
// T028 clarification checks: (1) fixture-count self-consistency vs the coverage constant;
// (2) the "29 schemas" caliber in tools/contracts/README.md vs CONTRACTS §2.
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const dirs = ['testdata/contracts/valid', 'testdata/contracts/invalid'];
const files = [];
for (const d of dirs) {
  for (const n of fs.readdirSync(d).filter((x) => x.endsWith('.json'))) files.push(d + '/' + n);
}
const count = (p, text) => {
  const j = JSON.parse(text);
  if (!Array.isArray(j)) return { n: 0, note: 'not an array' };
  return { n: j.length, note: '' };
};
let now = 0;
let head = 0;
console.log('fixture file                          now  HEAD  delta');
for (const f of files.sort()) {
  const a = count(f, fs.readFileSync(f, 'utf8'));
  let h = { n: 0 };
  try {
    h = count(f, execSync('git show HEAD:' + f, { maxBuffer: 1 << 26 }).toString('utf8'));
  } catch (e) {
    h = { n: 0 };
  }
  now += a.n;
  head += h.n;
  console.log(f.padEnd(38) + String(a.n).padStart(4) + String(h.n).padStart(6) + String(a.n - h.n).padStart(7) + (h.n === 0 ? '   (new file)' : ''));
}
console.log('TOTAL'.padEnd(38) + String(now).padStart(4) + String(head).padStart(6) + String(now - head).padStart(7));
const t = fs.readFileSync('tools/contracts/contracts.test.js', 'utf8');
console.log('\ncoverage constants in the test:');
for (const m of t.matchAll(/assert\.equal\([a-zA-Z]+(?:\.length)?,\s*(\d+)/g)) console.log('  assert.equal(..., ' + m[1] + ')');
console.log('schemas/*.schema.json count = ' + fs.readdirSync('schemas').filter((f) => f.endsWith('.schema.json')).length);
console.log('schemas/*.schema.json count at HEAD = ' + execSync('git ls-tree --name-only HEAD schemas', { maxBuffer: 1 << 24 }).toString('utf8').split('\n').filter((f) => f.endsWith('.schema.json')).length);
console.log('\n--- tools/contracts/README.md : lines mentioning a schema count ---');
fs.readFileSync('tools/contracts/README.md', 'utf8').split(/\r?\n/).forEach((l, i) => {
  if (/\b\d{2}\b/.test(l) && /schema/i.test(l)) console.log('  ' + (i + 1) + ': ' + l.trim().slice(0, 200));
});
console.log('\n--- docs/CONTRACTS.md : lines mentioning a schema count ---');
fs.readFileSync('docs/CONTRACTS.md', 'utf8').split(/\r?\n/).forEach((l, i) => {
  if (/\b\d{2}\b/.test(l) && /[Ss]chema|结构/.test(l)) console.log('  ' + (i + 1) + ': ' + l.trim().slice(0, 220));
});
console.log('\n--- docs/CONTRACTS_EN.md : lines mentioning a schema count ---');
fs.readFileSync('docs/CONTRACTS_EN.md', 'utf8').split(/\r?\n/).forEach((l, i) => {
  if (/\b\d{2}\b/.test(l) && /schema/i.test(l)) console.log('  ' + (i + 1) + ': ' + l.trim().slice(0, 220));
});

#!/usr/bin/env node
'use strict';
/*
 * Master's independent acceptance check for slice DIRECTIVE-REVIEW-SATURATION.
 * It deliberately does NOT trust the implementer's report: it re-derives the two
 * user-mandated properties over the WHOLE file (not just the changed lines):
 *   1. every §12.4 OB row (CN and EN) opens its disposition cell with a closed-set word;
 *   2. every §X.Y reference inside the G6-5 target regions (§0.3 table range and the
 *      Appendix C ledger range) resolves against the heading/bold-label set.
 *
 * Read-only. Usage: node tmp/v113/verify-v113.js
 */
const fs = require('fs');
const path = require('path');

const g6 = require(path.resolve('tools/gates/lib/criteria/g6.js'));
const it = g6.internals;
const CN = 'docs/DELIVERY_DIRECTIVE.md';
const EN = 'docs/DELIVERY_DIRECTIVE_EN.md';

const lines = (f) => fs.readFileSync(f, 'utf8').replace(/^\uFEFF/, '').split(/\n/);
const isEn = (f) => /_EN\.md$/i.test(f);

let fail = 0;

// ---- 1. ledger status closed set over the whole §12.4 table --------------------
for (const f of [CN, EN]) {
  const L = lines(f);
  const words = isEn(f) ? it.LEDGER_STATUS_SETS.en : it.LEDGER_STATUS_SETS.cn;
  const range = it.ledgerRegionRange(L);
  if (!range) {
    console.log('[verify] ' + f + ' :: no §12.4 ledger region found');
    fail++;
    continue;
  }
  let rows = 0;
  const bad = [];
  for (let i = range[0]; i < range[1]; i++) {
    const line = L[i] || '';
    if (!it.LEDGER_ROW_RE.test(line)) continue;
    rows++;
    const cells = it.rowCellsOf(line);
    const cell = cells.length ? cells[cells.length - 1] : '';
    const w = it.firstStatusWord(cell, words);
    if (!w) bad.push(i + 1 + ' :: ' + cell.slice(0, 60));
  }
  console.log('[verify] ' + f + ' ledger rows=' + rows + ' notInClosedSet=' + bad.length);
  bad.forEach((b) => console.log('        OFFENDER ' + b));
  fail += bad.length;
}

// ---- 2. G6-5 target regions: every reference resolves (whole file) -------------
const labelSet = new Set();
for (const f of [CN, EN]) for (const l of it.labelsFromLines(lines(f))) labelSet.add(l);
console.log('[verify] label set size=' + labelSet.size);

for (const f of [CN, EN]) {
  const L = lines(f);
  const ranges = it.targetRegionRanges(L);
  let refs = 0;
  const bad = [];
  for (const [start, end] of ranges) {
    for (let i = start; i < end; i++) {
      const line = L[i] || '';
      const re = new RegExp(it.SECTION_REF_RE.source, 'g');
      let m;
      while ((m = re.exec(line))) {
        refs++;
        if (!it.resolvesRef(m[1], labelSet)) bad.push(i + 1 + ' :: §' + m[1] + ' :: ' + line.slice(0, 70));
      }
    }
  }
  console.log('[verify] ' + f + ' region ranges=' + ranges.length + ' refs=' + refs + ' unresolved=' + bad.length);
  bad.forEach((b) => console.log('        OFFENDER ' + b));
  fail += bad.length;
}

console.log(fail === 0 ? '[verify] ACCEPTED (0 offenders)' : '[verify] REJECTED (' + fail + ' offenders)');
process.exit(fail === 0 ? 0 : 1);

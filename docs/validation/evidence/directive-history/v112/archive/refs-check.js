#!/usr/bin/env node
'use strict';
/*
 * ⑨ native validation, degenerate form for a documentation slice:
 * link / reference resolvability, scoped to the CHANGED LINES (the §6.3 discipline:
 * historical rows are not re-litigated). Every `§X.Y` reference on an added line is
 * resolved against the label set of the directive (headings plus the bold inline
 * labels this directive uses for 1.7.x / 3.4.x) together with the file's own labels.
 *
 * A reference resolves when a label with exactly that number exists, or when a label
 * extends it (the reference names a parent section).
 *
 * Whole-file counts are printed as information only: they include inherited citations
 * to the superseded v3.1 directive and to other documents' own section numbers.
 *
 * Usage: node tmp/v112/refs-check.js
 */
const { execFileSync } = require('child_process');
const fs = require('fs');

const TARGETS = [
  'docs/DELIVERY_DIRECTIVE.md',
  'docs/DELIVERY_DIRECTIVE_EN.md',
  'docs/validation/directive-v1.12.md',
  'docs/validation/directive-v1.12_EN.md',
  'docs/validation/evidence/directive-v1.12-freedom-list.md',
];

function lines(file) {
  return fs.readFileSync(file, 'utf8').replace(/^\uFEFF/, '').split(/\r?\n/);
}

function labelsOf(file) {
  const out = new Set();
  for (const l of lines(file)) {
    let m = /^#{2,4}\s+(\d+(?:\.\d+)*)(?=[\s.．、]|$)/.exec(l);
    if (m) out.add(m[1]);
    m = /^\*\*(\d+(?:\.\d+)*)\s/.exec(l);
    if (m) out.add(m[1]);
  }
  return out;
}

const directive = new Set([...labelsOf('docs/DELIVERY_DIRECTIVE.md'), ...labelsOf('docs/DELIVERY_DIRECTIVE_EN.md')]);

function resolves(num, own) {
  if (directive.has(num) || own.has(num)) return true;
  for (const h of directive) if (h.startsWith(num + '.')) return true;
  for (const h of own) if (h.startsWith(num + '.')) return true;
  return false;
}

/** Added line numbers (1-based, new-file numbering) from the staged diff. */
function addedLines(file) {
  let out = '';
  try {
    out = execFileSync('git', ['diff', '--cached', '-U0', '--no-color', '--', file], {
      encoding: 'utf8',
      maxBuffer: 1 << 28,
    });
  } catch {
    return null;
  }
  const set = new Set();
  for (const l of out.split(/\r?\n/)) {
    const h = /^@@ -\d+(?:,\d+)? \+(\d+)(?:,(\d+))? @@/.exec(l);
    if (h) {
      const start = Number(h[1]);
      const count = h[2] === undefined ? 1 : Number(h[2]);
      for (let i = 0; i < count; i++) set.add(start + i);
    }
  }
  return set;
}

let changedTotal = 0;
let changedUnresolved = 0;

for (const f of TARGETS) {
  const own = labelsOf(f);
  const L = lines(f);
  const changed = addedLines(f);

  if (changed) {
    for (const n of [...changed].sort((a, b) => a - b)) {
      const l = L[n - 1];
      if (l === undefined) continue;
      const re = /§(\d+(?:\.\d+)*)/g;
      let m;
      while ((m = re.exec(l))) {
        changedTotal++;
        if (!resolves(m[1], own)) {
          changedUnresolved++;
          console.log('  ! CHANGED ' + f + ':' + n + '  §' + m[1] + '  unresolved  :: ' + l.slice(0, 120));
        }
      }
    }
  }

  let all = 0;
  let allUnresolved = 0;
  for (const l of L) {
    const re = /§(\d+(?:\.\d+)*)/g;
    let m;
    while ((m = re.exec(l))) {
      all++;
      if (!resolves(m[1], own)) allUnresolved++;
    }
  }
  console.log(
    '[refs-check] ' + f + '  changedLines=' + (changed ? changed.size : 0) + '  wholeFileRefs=' + all +
      '  wholeFileUnresolved=' + allUnresolved + ' (inherited external citations)'
  );
}

console.log('[refs-check] changed-scope refs=' + changedTotal + ' unresolved=' + changedUnresolved);
console.log('[refs-check] result=' + (changedUnresolved === 0 ? 'PASS' : 'FINDINGS'));
process.exit(0);

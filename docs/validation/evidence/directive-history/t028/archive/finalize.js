'use strict';
/*
 * T028 closeout helper (local-only, tmp/ artifact).
 * Normalizes the edited markdown files to UTF-8 with BOM + LF and prints the
 * five mirror shapes so the CN/EN report pair can be compared mechanically.
 * It never writes outside the listed files.
 */

const fs = require('fs');

const files = [
  'docs/validation/t028-supervised-black-box.md',
  'docs/validation/t028-supervised-black-box_EN.md',
  'docs/t027/REMAINING_SLICES.md',
  'docs/t027/REMAINING_SLICES_EN.md',
  'docs/DEV_PLAN.md',
  'docs/DEV_PLAN_EN.md',
];

const BOM = '\uFEFF';
const shape = (text) => {
  const lines = text.replace(/\n+$/, '').split('\n');
  return {
    lines: lines.length,
    headings: lines.filter((l) => /^#{1,6} /.test(l)).length,
    bold: (text.match(/\*\*/g) || []).length,
    pipes: (text.match(/\|/g) || []).length,
    blanks: lines.filter((l) => l === '').length,
  };
};

for (const f of files) {
  const raw = fs.readFileSync(f, 'utf8').replace(/^\uFEFF/, '');
  const normalized = BOM + raw.replace(/\r\n/g, '\n');
  fs.writeFileSync(f, normalized, 'utf8');
  const s = shape(normalized);
  console.log(
    f.padEnd(46) +
      ' lines=' + s.lines +
      ' headings=' + s.headings +
      ' bold=' + s.bold +
      ' pipes=' + s.pipes +
      ' blanks=' + s.blanks
  );
}

const cn = shape(fs.readFileSync(files[0], 'utf8'));
const en = shape(fs.readFileSync(files[1], 'utf8'));
const same = JSON.stringify(cn) === JSON.stringify(en);
console.log('');
console.log('REPORT_SHAPES_IDENTICAL=' + same);
console.log('CR_REMAINING=' + files.reduce((n, f) => n + (fs.readFileSync(f, 'utf8').match(/\r/g) || []).length, 0));

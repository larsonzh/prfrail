'use strict';
/*
 * T028 docs-only correction round: normalize the edited markdown files to
 * UTF-8 with BOM + LF and print the five mirror shapes for each CN/EN pair so any
 * asymmetry introduced by the edits is visible immediately.
 * Read/write scope: exactly the 14 files listed below.
 */

const fs = require('fs');

const pairs = [
  ['docs/ADR_REGISTER.md', 'docs/ADR_REGISTER_EN.md'],
  ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md'],
  ['docs/BUSINESS_WORKFLOWS.md', 'docs/BUSINESS_WORKFLOWS_EN.md'],
  ['docs/INSTALLATION.md', 'docs/INSTALLATION_EN.md'],
  ['docs/INSTALLATION_PLAN.md', 'docs/INSTALLATION_PLAN_EN.md'],
  ['docs/OPERATIONS.md', 'docs/OPERATIONS_EN.md'],
  ['docs/S1_RELEASE_NOTES.md', 'docs/S1_RELEASE_NOTES_EN.md'],
  ['docs/S1_SUPPORT_MATRIX.md', 'docs/S1_SUPPORT_MATRIX_EN.md'],
  ['docs/t027/REMAINING_SLICES.md', 'docs/t027/REMAINING_SLICES_EN.md'],
  ['docs/validation/t028-supervised-black-box.md', 'docs/validation/t028-supervised-black-box_EN.md'],
];

const BOM = '\uFEFF';

function shapeOf(text) {
  const lines = text.replace(/\n+$/, '').split('\n');
  return {
    lines: lines.length,
    headings: lines.filter((l) => /^#{1,6} /.test(l)).length,
    bold: (text.match(/\*\*/g) || []).length,
    pipes: (text.match(/\|/g) || []).length,
    blanks: lines.filter((l) => l === '').length,
  };
}

for (const [cn, en] of pairs) {
  for (const f of [cn, en]) {
    const raw = fs.readFileSync(f, 'utf8').replace(/^\uFEFF/, '');
    fs.writeFileSync(f, BOM + raw.replace(/\r\n/g, '\n'), 'utf8');
  }
  const a = shapeOf(fs.readFileSync(cn, 'utf8'));
  const b = shapeOf(fs.readFileSync(en, 'utf8'));
  const same = JSON.stringify(a) === JSON.stringify(b);
  console.log(
    (same ? 'SAME  ' : 'DIFF  ') +
      cn.replace('docs/', '').padEnd(44) +
      ' CN=' + a.lines + '/' + a.headings + '/' + a.bold + '/' + a.pipes + '/' + a.blanks +
      '  EN=' + b.lines + '/' + b.headings + '/' + b.bold + '/' + b.pipes + '/' + b.blanks
  );
}

const reportShapes = pairs[pairs.length - 1];
const cr = pairs
  .flat()
  .reduce((n, f) => n + (fs.readFileSync(f, 'utf8').match(/\r/g) || []).length, 0);
console.log('');
console.log('REPORT_PAIR=' + reportShapes.join(' | '));
console.log('CR_REMAINING=' + cr);

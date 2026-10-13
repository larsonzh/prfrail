'use strict';
/*
 * Proves the docs-only correction round introduced NO new CN/EN asymmetry:
 * for each edited pair, it compares the five shapes of HEAD vs the worktree file and
 * requires the CN delta to equal the EN delta on every metric.
 * (Two pairs were already asymmetric at HEAD; that is pre-existing and reported, not
 * introduced here.)
 */

const cp = require('child_process');
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

const METRICS = ['lines', 'headings', 'bold', 'pipes', 'blanks'];

function shape(text) {
  const lines = text.replace(/\n+$/, '').split('\n');
  return {
    lines: lines.length,
    headings: lines.filter((l) => /^#{1,6} /.test(l)).length,
    bold: (text.match(/\*\*/g) || []).length,
    pipes: (text.match(/\|/g) || []).length,
    blanks: lines.filter((l) => l === '').length,
  };
}

const head = (p) =>
  cp
    .execSync('git show HEAD:' + p, { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 })
    .replace(/^\uFEFF/, '');

let failures = 0;
for (const [cn, en] of pairs) {
  const cnBefore = shape(head(cn));
  const cnAfter = shape(fs.readFileSync(cn, 'utf8').replace(/^\uFEFF/, ''));
  const enBefore = shape(head(en));
  const enAfter = shape(fs.readFileSync(en, 'utf8').replace(/^\uFEFF/, ''));
  const cnDelta = METRICS.map((m) => cnAfter[m] - cnBefore[m]);
  const enDelta = METRICS.map((m) => enAfter[m] - enBefore[m]);
  const symmetric = JSON.stringify(cnDelta) === JSON.stringify(enDelta);
  const alreadyAsymmetric = METRICS.some((m) => cnBefore[m] !== enBefore[m]);
  if (!symmetric) failures++;
  console.log(
    (symmetric ? 'SYMMETRIC' : 'ASYMMETRIC-NEW') +
      '  ' + cn.replace('docs/', '').padEnd(44) +
      ' CN delta=[' + cnDelta.join(',') + ']' +
      ' EN delta=[' + enDelta.join(',') + ']' +
      (alreadyAsymmetric ? '  (pair was already asymmetric at HEAD: no fix attempted)' : '')
  );
}

console.log('');
console.log('ASYMMETRY_INTRODUCED=' + failures);
process.exit(failures === 0 ? 0 : 1);

#!/usr/bin/env node
'use strict';
/*
 * M4 probe: does G6-4 ("expiring literal") actually judge the shapes the R5 review
 * claimed it judges? The review asserted C.0h's literals ("42 条", "19 项",
 * "1C+7H+22M+14L") are G6-4 "aggregate-count" forms and therefore need whitelist
 * entries. This probe runs the SHIPPED G6-4 subcheck (not a re-implementation) over
 *   (a) those literal shapes as synthetic lines, and
 *   (b) the real C.0h block of docs/DELIVERY_DIRECTIVE.md,
 * with positive controls so a zero result cannot be a broken harness.
 *
 * Usage: node tmp/v112/m4-g6-4-probe.js
 */
const fs = require('fs');
const path = require('path');

const G6 = require(path.join(process.cwd(), 'tools', 'gates', 'lib', 'criteria', 'g6.js'));
const G6_4 = G6.subchecks.find((s) => s.id === 'G6-4');
if (!G6_4) throw new Error('G6-4 subcheck not found');

/** Minimal ctx implementing only what G6-4 reads. */
function makeCtx(file, lines) {
  return {
    changedMarkdown: () => [file],
    readLineArray: () => lines,
    addedLineNumbers: () => lines.map((_, i) => i + 1),
  };
}

function run(file, lines) {
  const out = {
    records: [],
    rec(id, status, label, detail) {
      this.records.push({ id, status, label, detail });
    },
  };
  G6_4.run(makeCtx(file, lines), out);
  return out;
}

const NEGATIVE = [
  '外部独立评审连续三轮对准则本体共提出 **64 条**发现（口径：按报告逐轮累计）',
  '| 为何必需 | 处置外部独立评审累计提出的 **64 条**发现 |',
  '**目标**：处置评审累计提出的 64 条发现 + 落地 A3 / H4 用户裁决，升 v1.12',
  '分解为 1 Critical + 11 High + 35 Medium + 17 Low',
  '本片 19 项已补做，另 22 条归口下一片',
  '严重度分解 1C+7H+22M+14L',
];

const POSITIVE = [
  '本片共 14 个文件的位置统一（实测全部一致）',
  '实测 42 行样例',
  '增量 164/80 行',
  '共 42 个文件受影响',
];

const findingsOf = (out) => {
  const r = out.records.find((x) => x.id === 'G6-4');
  if (!r) throw new Error('G6-4 did not emit a record');
  return { status: r.status, detail: String(r.detail || '') };
};

console.log('[m4-probe] shipped G6-4 loaded from tools/gates/lib/criteria/g6.js');
console.log('');

const neg = findingsOf(run('synthetic-negative.md', NEGATIVE));
const pos = findingsOf(run('synthetic-positive.md', POSITIVE));
console.log('[m4-probe] negative set (claimed aggregate-count forms)');
for (const l of NEGATIVE) console.log('    - ' + l);
console.log('[m4-probe] negative status = ' + neg.status + '   (expected PASS)');
console.log('[m4-probe] negative detail = ' + neg.detail);
console.log('');
console.log('[m4-probe] positive control set');
for (const l of POSITIVE) console.log('    - ' + l);
console.log('[m4-probe] positive status = ' + pos.status + '   (expected FAIL)');
console.log('[m4-probe] positive detail = ' + pos.detail);
console.log('');

// (b) the real C.0h block, scanned line by line with the shipped judgement.
const cn = path.join(process.cwd(), 'docs', 'DELIVERY_DIRECTIVE.md');
const all = fs.readFileSync(cn, 'utf8').replace(/^\uFEFF/, '').split(/\r?\n/);
const start = all.findIndex((l) => /^### C\.0h /.test(l));
if (start < 0) throw new Error('C.0h heading not found in ' + cn);
let end = start + 1;
while (end < all.length && !/^#{2,4} /.test(all[end])) end++;
const block = all.slice(start, end);
const real = findingsOf(run('docs/DELIVERY_DIRECTIVE.md', block));
console.log('[m4-probe] real C.0h block: lines ' + (start + 1) + '-' + end + ' (count=' + block.length + ')');
console.log('[m4-probe] real C.0h status = ' + real.status + '   (expected PASS)');
console.log('[m4-probe] real C.0h detail = ' + real.detail);
console.log('');

const verdict =
  neg.status === 'PASS' && pos.status === 'FAIL' && real.status === 'PASS'
    ? 'M4 NOT REPRODUCIBLE (no whitelist entry needed)'
    : 'M4 INCONCLUSIVE - inspect above';
console.log('[m4-probe] verdict = ' + verdict);
process.exit(0);

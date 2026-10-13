'use strict';
// 9th-stage (native validation) independent experiment, per DELIVERY_DIRECTIVE 7.7 A12:
// because the 6th-stage scan reported PASS, this stage must prove a key invariant with a
// criterion that is independent of that scan. Two independent invariants are exercised here:
//   (1) the G5-b judge is LIVE and its closed set is load-bearing (synthetic in-memory ctx,
//       no git involvement, no repository writes);
//   (2) the section 2.2 `modelId` column and the G5-b closed set are equal in BOTH directions,
//       and the new GLM rows carry the expected qualified-name -> modelId mapping.
const fs = require('fs');
const path = require('path');

const g5b = require(path.resolve('tools/gates/lib/criteria/g5b.js'));

let failures = 0;
const check = (name, ok, extra) => {
  if (!ok) failures += 1;
  console.log('A12 ' + (ok ? 'OK  ' : 'FAIL') + ' ' + name + (extra === undefined ? '' : ' :: ' + extra));
};

// ---------- part 1: synthetic live-fire of the judge ----------
const mkCtx = (files) => ({
  changedPaths: Object.keys(files),
  readText: (f) => (Object.prototype.hasOwnProperty.call(files, f) ? files[f] : null),
});
const capture = () => ({
  recs: [],
  rec(id, verdict, label, detail) {
    this.recs.push({ id, verdict, label, detail });
  },
});
const agent = (modelLine) =>
  '---\nname: fixture\n' + modelLine + '\ndescription: fixture\n---\n\nBody.\n';

const CASES = [
  ['member glm-5.3', 'model: glm-5.3', 'PASS'],
  ['member glm-5.3-flash', 'model: glm-5.3-flash', 'PASS'],
  ['non-whitelist gpt-5.4-terra', 'model: gpt-5.4-terra', 'FAIL'],
  ['case variant GLM-5.3', 'model: GLM-5.3', 'FAIL'],
  ['qualified call string', 'model: glm-5.3 (glm)', 'FAIL'],
  ['quoted member', "model: 'glm-5.3'", 'FAIL'],
  ['key missing (2.6 passes)', 'description: no model key here', 'PASS'],
  ['empty value (G5-a 2 steps in)', 'model: ""', 'PASS'],
];
for (const [name, line, expect] of CASES) {
  const out = capture();
  g5b.run(mkCtx({ '.github/agents/fixture.agent.md': agent(line) }), out);
  const got = out.recs[0].verdict;
  check('live-fire: ' + name, got === expect, 'expect=' + expect + ' got=' + got);
}

// ---------- part 2: section 2.2 modelId column <-> closed set ----------
const readRows = (file) => {
  const lines = fs.readFileSync(file, 'utf8').split('\n');
  const i = lines.findIndex((l) => l.startsWith('### 2.2 '));
  const j = lines.findIndex((l, k) => k > i && l.startsWith('### 2.3 '));
  return lines
    .slice(i, j)
    .filter((l) => l.startsWith('|'))
    .map((l) => l.split('|').map((s) => s.trim()))
    .filter((c) => /^`[a-z0-9.\-]+`$/.test(c[5] || ''))
    .map((c) => ({
      role: c[1],
      slot: c[2],
      carrier: c[3],
      qualified: c[4].replace(/`/g, ''),
      modelId: c[5].replace(/`/g, ''),
    }));
};
const CN = readRows('docs/DELIVERY_DIRECTIVE.md');
const EN = readRows('docs/DELIVERY_DIRECTIVE_EN.md');
const uniq = (a) => [...new Set(a)].sort();
const idsCN = uniq(CN.map((r) => r.modelId));
const idsEN = uniq(EN.map((r) => r.modelId));
const code = uniq(g5b.closedSet);

check('2.2 row count CN == EN', CN.length === EN.length, CN.length + '/' + EN.length);
check('2.2 modelId set CN == EN', JSON.stringify(idsCN) === JSON.stringify(idsEN));
check('closed set == 2.2 modelId set', JSON.stringify(idsCN) === JSON.stringify(code), idsCN.join(','));
check(
  'reverse: every closed-set member exists in 2.2',
  code.every((x) => idsCN.includes(x)),
  'orphans=' + code.filter((x) => !idsCN.includes(x)).join(',')
);
const find = (rows, pred) => rows.find(pred);
check(
  'mapping GLM-5.3 (glm) -> glm-5.3',
  !!find(CN, (r) => r.qualified === 'GLM-5.3 (glm)' && r.modelId === 'glm-5.3' && r.slot === '⑤')
);
check(
  'mapping GLM-5.3-Flash (glm) -> glm-5.3-flash',
  !!find(CN, (r) => r.qualified === 'GLM-5.3-Flash (glm)' && r.modelId === 'glm-5.3-flash')
);
check(
  '⑤ preferred row = GLM, ⑤ backup row = deepseek-v4-pro',
  !!find(CN, (r) => r.slot === '⑤' && r.modelId === 'glm-5.3') &&
    !!find(CN, (r) => r.slot === '⑤' && r.modelId === 'deepseek-v4-pro')
);
check(
  'architect (1) and 2/3/8 slots unchanged',
  !!find(CN, (r) => r.slot === '①' && r.modelId === 'deepseek-v4-pro') &&
    CN.filter((r) => r.modelId === 'deepseek-flash').length === 3
);

// ---------- part 3: boundary invariant - no glm leakage into frozen sections ----------
const section = (file, head, tail) => {
  const lines = fs.readFileSync(file, 'utf8').split('\n');
  const i = lines.findIndex((l) => l.startsWith(head));
  const j = lines.findIndex((l, k) => k > i && l.startsWith(tail));
  return lines.slice(i, j).join('\n');
};
for (const [file, tag] of [
  ['docs/DELIVERY_DIRECTIVE.md', 'CN'],
  ['docs/DELIVERY_DIRECTIVE_EN.md', 'EN'],
]) {
  const s24 = section(file, '### 2.4 ', '### 2.5 ');
  const s52 = section(file, '### 5.2 ', '### 5.3 ');
  const s75 = section(file, '### 7.5 ', '### 7.6 ');
  const sB1 = section(file, '### B.1 ', '### B.2 ');
  const sB5 = section(file, '### B.5 ', '### B.6 ');
  check(tag + ' 2.4 body has no glm token', !/glm/i.test(s24));
  check(tag + ' 5.2 section has no glm token', !/glm/i.test(s52));
  check(tag + ' 7.5 figures have no glm token', !/glm/i.test(s75));
  check(tag + ' B.1 (architect) still DeepSeek V4 Pro', /DeepSeek V4 Pro \(deepseek\)/.test(sB1));
  check(tag + ' B.5 (pre-reviewer) uses GLM-5.3', /GLM-5\.3 \(glm\)/.test(sB5));
}

console.log('A12_EXPERIMENT_FAILURES=' + failures);
process.exit(failures === 0 ? 0 : 1);

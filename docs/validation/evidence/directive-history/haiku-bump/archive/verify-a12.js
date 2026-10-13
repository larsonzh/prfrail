'use strict';
// ⑨ A12 dedicated independent experiment for slice DIRECTIVE-HAIKU-BUMP.
// Required because ⑥ reported PASS: one experiment, independent of ⑥, must prove the key invariants
// mechanically. Reads only; writes nothing.
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const ROOT = path.resolve(__dirname, '..', '..');
const CN = 'docs/DELIVERY_DIRECTIVE.md';
const EN = 'docs/DELIVERY_DIRECTIVE_EN.md';
const G5B = require(path.join(ROOT, 'tools/gates/lib/criteria/g5b.js'));

let fails = 0;
const ok = (name, cond, detail) => {
  if (!cond) fails++;
  console.log((cond ? 'OK   ' : 'FAIL ') + name + (detail ? ' :: ' + detail : ''));
};

const read = (p) => fs.readFileSync(path.join(ROOT, p), 'utf8').replace(/\r\n/g, '\n');
const head = (p) => execSync('git show HEAD:' + p, { cwd: ROOT, maxBuffer: 1 << 28 }).toString('utf8').replace(/\r\n/g, '\n');
const seg = (text, a, b) => {
  const lines = text.split('\n');
  const i = lines.findIndex((l) => l.startsWith(a));
  const j = lines.findIndex((l, k) => k > i && l.startsWith(b));
  return i < 0 || j < 0 ? null : lines.slice(i, j).join('\n');
};
const sha16 = (s) => require('crypto').createHash('sha256').update(s, 'utf8').digest('hex').slice(0, 16);

console.log('== (a) G5-b criteria live-fire on a synthetic context (no git, no repo write) ==');
function fire(modelLine) {
  const out = [];
  out.rec = (id, status) => out.push([id, status]);
  G5B.run({ changedPaths: ['x.agent.md'], readText: () => '---\n' + modelLine + '\n---\n' }, out);
  return out[0][1];
}
const live = [
  ['model: deepseek-v4-pro', 'PASS'],
  ['model: glm-5.3', 'PASS'],
  ['model: glm-5.3-flash', 'PASS'],
  ['model: gpt-6-luna', 'PASS'],
  ['model: claude-haiku-5.5', 'FAIL'],            // the slice's key fact: internal id is NOT a member
  ['model: claude-haiku-4.5', 'FAIL'],
  ['model: Claude Haiku 5.5 (copilot)', 'FAIL'],  // qualified-name form must not pass G5-b
  ['model: "claude-haiku-5.5"', 'FAIL'],          // quoted value is not verbatim-equal
  ['agentName: xxx', 'PASS'],                     // missing key defers to 2.6
  ['model: ""', 'PASS'],                          // empty value is G5-a 2's finding
];
for (const [line, want] of live) {
  const got = fire(line);
  ok('live-fire ' + JSON.stringify(line), got === want, 'got=' + got + ' want=' + want);
}

console.log('== (b) document <-> code set equality and the Haiku exclusion ==');
function modelIds(p) {
  const body = seg(read(p), '### 2.2 ', '### 2.3 ');
  const lines = body.split('\n').filter((l) => l.trim().startsWith('|'));
  const cells = (l) => l.replace(/^\||\|$/g, '').split('|').map((c) => c.trim());
  // Header detection must be language-agnostic: CN glosses the column "（内部标识）", EN "(internal identifier)".
  const hi = lines.findIndex((l) => cells(l).some((c) => c.startsWith('`modelId`')));
  const idx = cells(lines[hi]).findIndex((c) => c.startsWith('`modelId`'));
  const out = [];
  for (let k = hi + 1; k < lines.length; k++) {
    const c = cells(lines[k])[idx];
    if (!c || c === '---' || c === '—') continue;
    out.push(c.replace(/`/g, '').trim());
  }
  return out;
}
const cnIds = modelIds(CN);
const enIds = modelIds(EN);
const cnSet = [...new Set(cnIds)].sort();
const enSet = [...new Set(enIds)].sort();
const codeSet = [...G5B.closedSet].sort();
ok('CN 2.2 modelId set === g5b CLOSED_SET', JSON.stringify(cnSet) === JSON.stringify(codeSet), cnSet.join(','));
ok('EN 2.2 modelId set === g5b CLOSED_SET', JSON.stringify(enSet) === JSON.stringify(codeSet));
ok('CN and EN modelId sets equal', JSON.stringify(cnSet) === JSON.stringify(enSet));
ok('closed set size is 8', codeSet.length === 8, 'size=' + codeSet.length);
ok('claude-haiku-5.5 NOT in closed set', !codeSet.includes('claude-haiku-5.5'));
ok('claude-haiku-4.5 NOT in closed set', !codeSet.includes('claude-haiku-4.5'));

console.log('== (c) freeze held: five boundary regions zero-diff (CN and EN) ==');
for (const f of [CN, EN]) {
  const isEn = /_EN\.md$/.test(f);
  const regions = [
    ['2.2', '### 2.2 ', '### 2.3 '],
    ['2.4', '### 2.4 ', '### 2.5 '],
    ['6.3', '### 6.3 ', '### 6.4 '],
    ['7.5', '### 7.5 ', '### 7.6 '],
    ['appendix B', isEn ? '## Appendix B' : '## 附录 B', isEn ? '## Appendix C' : '## 附录 C'],
  ];
  for (const [name, a, b] of regions) {
    const h = seg(head(f), a, b);
    const w = seg(read(f), a, b);
    ok('frozen ' + f + ' ' + name, h !== null && w !== null && h === w, 'head=' + (h === null ? 'MISS' : sha16(h)) + ' work=' + (w === null ? 'MISS' : sha16(w)));
  }
}

console.log('== (d) 7.2.3 minimality: only the version string changed ==');
for (const f of [CN, EN]) {
  const para = (t) => {
    const lines = t.split('\n');
    const i = lines.findIndex((l) => l.startsWith('### 7.2.3 '));
    return lines.slice(i + 1).find((l) => l.trim() !== '');
  };
  const h = para(head(f));
  const w = para(read(f));
  ok('7.2.3 first paragraph differs only by version', h.replace('Claude Haiku 4.5', 'Claude Haiku 5.5') === w, f);
}

console.log('== (e) OB-74 closure: literals and ledger row ==');
for (const f of [CN, EN]) {
  const t = read(f);
  const newSeg = seg(t, '### 7.2.3 ', '### 7.3 ');
  ok('7.2.3 carries the qualified name', newSeg.includes('Claude Haiku 5.5 (copilot)'), f);
  ok('7.2.3 carries the internal identifier', newSeg.includes('claude-haiku-5.5'), f);
  // Bold markers interrupt the literal phrase ("is **not** a member of"), so flatten them before testing.
  const flatSeg = newSeg.replace(/\*/g, '');
  ok('7.2.3 states the G5-b exclusion', /G5-b/.test(flatSeg) && /not a member of|不属于/.test(flatSeg), f);
  const row = t.split('\n').find((l) => l.includes('| OB-74 |'));
  const hRow = head(f).split('\n').find((l) => l.includes('| OB-74 |'));
  ok('OB-74 row present', !!row, f);
  ok('OB-74 disposition starts with a G6-6 member', /^\*\*已处置\*\*|^\*\*Resolved\*\*/.test(row.split('|')[4].trim()), row.split('|')[4].trim().slice(0, 14));
  ok('OB-74 phenomenon column byte-identical to HEAD', row.split('|')[3] === hRow.split('|')[3], f);
  const strayOld = read(f).split('\n').filter((l) => /Claude Haiku 4\.5/.test(l)).filter((l) => !/^\| v1\.27 \|/.test(l) && !/^\| (目标|Goal) \|/.test(l));
  ok('old version string survives only in the v1.27 log row and the C.0p goal row', strayOld.length === 0, f + ' stray=' + strayOld.length);
  ok('7.2.3 no longer carries the old precondition', !/Claude Haiku 4\.5/.test(newSeg), f);
}

console.log('== (f) G6-6 status closed set sanity (the 已落地 trap) ==');
const g6 = fs.readFileSync(path.join(ROOT, 'tools/gates/lib/criteria/g6.js'), 'utf8');
const setLine = g6.split('\n').find((l) => l.includes('已处置') && l.includes('待议')) || '';
const members = (setLine.match(/[\u4e00-\u9fa5]+/g) || []);
ok('G6-6 set contains 已处置', members.includes('已处置'), setLine.trim().slice(0, 90));
ok('G6-6 set does NOT contain 已落地', !members.includes('已落地'));

console.log('== (g) C.0p structural facts ==');
const c0p = seg(read(CN), '### C.0p ', '### C.1 ');
ok('C.0p exists', c0p !== null);
ok('C.0p status is in-progress', /⬜ \*\*进行中\*\*/.test(c0p));
ok('C.0p declares the 1-probe cap', /1 次（硬上限，不追加）/.test(c0p));
ok('C.0p states Haiku did not review', /Haiku 未承担本片审查工作/.test(c0p));

console.log('A12_EXPERIMENT_FAILURES=' + fails);
process.exitCode = fails === 0 ? 0 : 1;

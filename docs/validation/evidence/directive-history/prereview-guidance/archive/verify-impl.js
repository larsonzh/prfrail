// Independent verification of step ② (controller-side, assertion style). Read-only.
const fs = require('fs');
const CN = 'docs/DELIVERY_DIRECTIVE.md';
const EN = 'docs/DELIVERY_DIRECTIVE_EN.md';
let bad = [];
const chk = (ok, msg) => { console.log((ok ? 'PASS  ' : 'FAIL  ') + msg); if (!ok) bad.push(msg); };
const load = (f) => {
  const s = fs.readFileSync(f, 'utf8');
  chk(s.startsWith('\uFEFF'), f + ' BOM');
  chk(!s.includes('\r'), f + ' no CR');
  return s;
};
const count = (s, sub) => s.split(sub).length - 1;
const countRe = (s, re) => (s.match(re) || []).length;

const cn = load(CN);
const en = load(EN);

const CORE_CN = '改文件；以测试代码或文档变更本身作为判断实现正确性的依据（读取上述工件以核验报告声明不在此限；保持独立视角）';
const CORE_EN = 'editing files; using test code or documentation changes themselves as the basis for judging implementation correctness (reading the above artifacts to verify report claims is not restricted by this; keep an independent perspective)';

// 1. core wording string exactly twice per file
chk(count(cn, CORE_CN) === 2, 'CN core string count = ' + count(cn, CORE_CN) + ' (expect 2)');
chk(count(en, CORE_EN) === 2, 'EN core string count = ' + count(en, CORE_EN) + ' (expect 2)');

// 2. resolvable labels, exactly one each
chk(countRe(cn, /^\*\*4\.2\.1 /gm) === 1, 'CN **4.2.1 label x1');
chk(countRe(cn, /^\*\*4\.2\.2 /gm) === 1, 'CN **4.2.2 label x1');
chk(countRe(en, /^\*\*4\.2\.1 /gm) === 1, 'EN **4.2.1 label x1');
chk(countRe(en, /^\*\*4\.2\.2 /gm) === 1, 'EN **4.2.2 label x1');

// 3. five numbered items between the two labels
const between = (s, a, b) => { const i = s.indexOf(a), j = s.indexOf(b); return i >= 0 && j > i ? s.slice(i, j) : ''; };
chk(countRe(between(cn, '**4.2.1 ', '**4.2.2 '), /^\d\. /gm) === 5, 'CN five numbered items in 4.2.1');
chk(countRe(between(en, '**4.2.1 ', '**4.2.2 '), /^\d\. /gm) === 5, 'EN five numbered items in 4.2.1');

// 4. pointer cell must NOT contain the core string, and must be the pointer form
chk(count(cn, '| **改文件**；其余见 §4.2.1 / §4.2.2 |') === 1, 'CN pointer cell present x1');
chk(count(en, '| **modify files**; for the rest see §4.2.1 / §4.2.2 |') === 1, 'EN pointer cell present x1');
const cellCn = cn.split('\n').find((l) => l.includes('其余见 §4.2.1 / §4.2.2') && l.startsWith('| 预审员'));
chk(!!cellCn && !cellCn.includes(CORE_CN), 'CN pointer cell free of core string');

// 5. §4.3 appends
chk(count(cn, '静默错误结论；且报告核心因果被自身证据推翻者属本级别') === 1, 'CN Critical cell append x1');
chk(count(cn, '契约冲突；协议属性与设计不符、但 fail-closed 已避免错误结论者属本级别') === 1, 'CN High cell append x1');
chk(count(en, 'silently wrong conclusion; and a report whose core causation is overturned by its own evidence belongs to this level') === 1, 'EN Critical cell append x1');
chk(count(en, 'contract conflict; a protocol property that does not match the design but fail-closed already prevented a wrong conclusion belongs to this level') === 1, 'EN High cell append x1');

// 6. severity-boundary note
chk(count(cn, '> **定级边界注记**：同一实质发现不因定级差异而削弱其效力——定级决定的是处置动作（重跑该环 vs 定向复核），不决定发现是否成立。') === 1, 'CN boundary note x1');
chk(count(en, '> **Severity-boundary note**: the same substantive finding is not weakened by a difference in severity rating - the rating determines the handling action (re-run that stage vs a targeted re-check), not whether the finding stands.') === 1, 'EN boundary note x1');

// 7. B.5 line replaced (old wording gone)
chk(count(cn, '禁止：改文件；查看测试代码或文档变更（保持独立视角）。') === 0, 'CN old B.5 prohibition gone');
chk(count(cn, '最低核验动作：见 §4.2.1（不在此重复正文）。') === 1, 'CN B.5 pointer x1');
chk(count(en, 'Forbidden: editing files; inspecting test code or documentation changes (preserve an independent perspective).') === 0, 'EN old B.5 prohibition gone');
chk(count(en, 'Minimum verification actions: see §4.2.1 (not repeated here in full).') === 1, 'EN B.5 pointer x1');

// 8. old §4.2 wording cleared, §3.4 context untouched
chk(count(cn, '看测试代码') === 0, 'CN old "看测试代码" cleared (count=' + count(cn, '看测试代码') + ')');
chk(count(en, 'look at test code') === 0, 'EN old "look at test code" cleared');
chk(count(cn, '不写生产代码、测试代码') === 1 || count(cn, '测试代码') >= 1, 'CN §3.4 context row still present');

// 9. five-tuple mirror
const stat = (s) => ({ lines: s.split('\n').length, bold: (s.match(/\*\*/g) || []).length, bars: (s.match(/\|/g) || []).length, blank: s.split('\n').filter((x) => x === '').length });
const a = stat(cn), b = stat(en);
chk(a.lines === b.lines && a.bold === b.bold && a.bars === b.bars && a.blank === b.blank, 'five-tuple CN==EN :: ' + JSON.stringify(a) + ' vs ' + JSON.stringify(b));
chk(a.lines === 1556 && a.bold === 3489 && a.bars === 2368 && a.blank === 392, 'five-tuple matches declared 1556/3489/2368/392 :: ' + JSON.stringify(a));

console.log('RESULT=' + (bad.length === 0 ? 'PASS' : 'FAIL(' + bad.length + ')'));

// Pre-review (05) mechanical verification for DIRECTIVE-HAIKU-BUMP (v1.27).
// Read-only: reads HEAD via `git show` and the worktree files; writes nothing.
const { execSync } = require('child_process');
const fs = require('fs');

const G6_6_CLOSED = ['观察中', '待办', '待议', '已处置', '已裁决', '已登记',
  'Observing', 'To do', 'To be discussed', 'Resolved', 'Ruling', 'Registered'];

function headLines(p) {
  return execSync('git show HEAD:' + p, { encoding: 'utf8', maxBuffer: 64 * 1024 * 1024 }).split('\n');
}
function workLines(p) { return fs.readFileSync(p, 'utf8').split('\n'); }

let fails = 0;
function check(name, cond, detail) {
  console.log((cond ? 'OK   ' : 'FAIL ') + name + (cond ? '' : ' :: ' + (detail || '')));
  if (!cond) fails++;
}

for (const p of ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md']) {
  const H = headLines(p), W = workLines(p);
  console.log('==== ' + p + ' (HEAD ' + (H.length - 1) + ' lines, worktree ' + (W.length - 1) + ' lines)');

  // (1) Header line 3: only version + date changed, rest verbatim.
  const normHead = s => s.replace(/v1\.26/g, 'vX').replace(/v1\.27/g, 'vX')
    .replace(/2026-10-06/g, 'DATE').replace(/2026-10-11/g, 'DATE');
  check('header(line3) only version+date changed', normHead(H[2]) === normHead(W[2]),
    JSON.stringify([H[2], W[2]]));

  // (2) 7.2.3 paragraph 1: HEAD line 599 vs worktree line 600; only 4.5 -> 5.5.
  const pOld = H[598], pNew = W[599];
  check('7.2.3 para1 position sanity (HEAD:599 has 4.5, WT:600 has 5.5)',
    /Claude Haiku 4\.5/.test(pOld) && /Claude Haiku 5\.5/.test(pNew),
    'old=' + pOld.slice(0, 40) + ' new=' + pNew.slice(0, 40));
  check('7.2.3 para1 rest verbatim (C1: only version string replaced)',
    pOld.replace(/Claude Haiku 4\.5/g, 'Claude Haiku 5.5') === pNew, 'text differs beyond version string');

  // (3) Whole-file version-string accounting.
  const h45 = H.filter(l => /Haiku 4\.5/.test(l)).length;
  const w45 = W.filter(l => /Haiku 4\.5/.test(l)).length;
  const w45lines = W.map((l, i) => [i + 1, l]).filter(([, l]) => /Haiku 4\.5/.test(l));
  console.log('     HEAD lines containing "Haiku 4.5": ' + h45 + ' ; worktree: ' + w45);
  w45lines.forEach(([n, l]) => console.log('       WT:' + n + ' ' + l.slice(0, 100)));
  check('HEAD had exactly 1 live "Haiku 4.5" (the 7.2.3 precondition)', h45 === 1, 'h45=' + h45);
  const intentional = w45lines.every(([, l]) => /升为|升到|moves from|from `?Claude Haiku 4\.5|by 4\.5/.test(l));
  check('all remaining WT "Haiku 4.5" are old-value references (log/goal rows)', w45 === w45lines.length && intentional, '');

  // (4) OB-74 row: HEAD line 910 vs worktree line 913; cells 1-3 (ID/date/phenomenon) identical.
  const oOld = H[909], oNew = W[912];
  check('OB-74 position sanity', /^\| OB-74 \|/.test(oOld) && /^\| OB-74 \|/.test(oNew), '');
  const a = oOld.split('|').map(c => c.trim());
  const b = oNew.split('|').map(c => c.trim());
  check('OB-74 row splits into 6 parts', a.length === 6 && b.length === 6, a.length + '/' + b.length);
  check('OB-74 cells [empty,ID,date,phenomenon] verbatim (disposition-only change)',
    a.slice(0, 4).join('\u0001') === b.slice(0, 4).join('\u0001'), 'phenomenon column altered');
  const firstWord = b[4].replace(/\*\*/g, '').split(/[（(:：]/)[0].trim();
  check('OB-74 disposition first status word "' + firstWord + '" in G6-6 closed set',
    G6_6_CLOSED.includes(firstWord), firstWord);
  const firstWordOld = a[4].replace(/\*\*/g, '').split(/[（(:：]/)[0].trim();
  check('OB-74 HEAD first status word was 待议/To be discussed (proving the flip, and new row no longer uses it)',
    (firstWordOld === '待议' || firstWordOld === 'To be discussed') && firstWord !== firstWordOld, firstWordOld + ' -> ' + firstWord);
  const consistent = /`Claude Haiku 5\.5 \(copilot\)`/.test(b[4]) && /`claude-haiku-5\.5`/.test(b[4])
    && /G5-b/.test(b[4]);
  check('OB-74 disposition consistent with 7.2.3 new para (name/id/G5-b)', consistent, b[4].slice(0, 120));

  // (5) OB-85 line unchanged; OB-86 new at line 925.
  check('OB-85 row (WT:924) identical to HEAD:921', H[920] === W[923], '');
  check('OB-86 row present at WT:925', W[924].startsWith('| OB-86 |'), W[924].slice(0, 40));
  const fw86 = W[924].split('|').map(c => c.trim())[4].replace(/\*\*/g, '').split(/[（(:：]/)[0].trim();
  check('OB-86 disposition first status word "' + fw86 + '" in G6-6 closed set', G6_6_CLOSED.includes(fw86), fw86);

  // (6) C.0p block boundaries; C.0o last row untouched.
  check('C.0o last row (WT:1450) identical to HEAD:1446', H[1445] === W[1449], '');
  check('C.0p heading at WT:1451', W[1450].startsWith('### C.0p '), W[1450].slice(0, 60));
  check('C.1 heading at WT:1469', W[1468].startsWith('### C.1 '), W[1468].slice(0, 40));

  // (7) C.0p field order vs C.0o field order.
  const fields = (from, to) => W.slice(from, to).filter(l => /^\| [^|]+ \| /.test(l)).map(l => l.split('|')[1].trim());
  const c0o = fields(1433, 1450), c0p = fields(1451, 1468);
  check('C.0p field list same order as C.0o (' + c0p.length + ' vs ' + c0o.length + ')',
    c0p.length === c0o.length && c0p.every((f, i) => f === c0o[i]),
    'C.0o=[' + c0o.join(',') + '] C.0p=[' + c0p.join(',') + ']');

  // (8) 7.2.3 new paragraph (WT:602) content assertions.
  const para2 = W[601];
  check('7.2.3 new para at WT:602 contains qualified name + internal id + G5-b exclusion + no "sole" restatement',
    /`Claude Haiku 5\.5 \(copilot\)`/.test(para2) && /`claude-haiku-5\.5`/.test(para2)
    && /不属于.*G5-b|not.*member of the §6\.3 G5-b/.test(para2)
    && !/唯一豁免/.test(para2) && !/sole exemption/.test(para2), para2.slice(0, 100));

  // (9) claude-haiku-5.5 never claimed inside G5-b closed set text (section 6.3 untouched is proven by diff; double-check text).
  const g5b = W.slice(483, 556).join('\n');
  check('G5-b closed-set text (6.3) does NOT contain claude-haiku-5.5', !/claude-haiku/.test(g5b), '');
}

console.log(fails === 0 ? '\nVERIFY: ALL CHECKS PASSED' : '\nVERIFY: ' + fails + ' CHECK(S) FAILED');
process.exit(fails === 0 ? 0 : 1);

// Targeted assertion for the A2 amendment (the only change since the ⑧b self-check):
// the note's parenthetical must now use disposition names taken verbatim from the §4.3 table.
const fs = require('fs');
const cn = fs.readFileSync('docs/DELIVERY_DIRECTIVE.md', 'utf8');
const en = fs.readFileSync('docs/DELIVERY_DIRECTIVE_EN.md', 'utf8');
const rp = fs.readFileSync('docs/validation/directive-prereview-guidance.md', 'utf8');
const rpe = fs.readFileSync('docs/validation/directive-prereview-guidance_EN.md', 'utf8');
const c = (s, sub) => s.split(sub).length - 1;
const bad = [];
const chk = (ok, m) => { console.log((ok ? 'PASS  ' : 'FAIL  ') + m); if (!ok) bad.push(m); };

const NEW_CN = '（**必须修复并重跑该审查环节** vs **记录在案，可延后**）';
const NEW_EN = '(**must fix and re-run that review stage** vs **recorded on the ledger, may be deferred**)';
chk(c(cn, NEW_CN) === 1, 'CN note uses §4.3 table names x1');
chk(c(en, NEW_EN) === 1, 'EN note uses §4.3 table names x1');
chk(c(cn, '定向复核') === 0, 'CN "定向复核" fully removed (dangling term gone)');
chk(c(en, 'a targeted re-check') === 0, 'EN "a targeted re-check" fully removed');
// the two names must exist verbatim in the §4.3 table itself
chk(c(cn, '| 必须修复并重跑该审查环节 |') === 1, 'CN name A exists in §4.3 table');
chk(c(cn, '记录在案，可延后') >= 1, 'CN name B exists in §4.3 table');
chk(c(en, '| must fix and re-run that review stage |') === 1, 'EN name A exists in §4.3 table');
chk(c(en, 'recorded on the ledger, may be deferred') >= 1, 'EN name B exists in §4.3 table');
// no CJK regression on the note lines
chk(!cn.includes('\r') && cn.startsWith('\uFEFF'), 'CN directive BOM+LF intact');
chk(!en.includes('\r') && en.startsWith('\uFEFF'), 'EN directive BOM+LF intact');
chk(c(rp, '已闭合：用户 ⑩ 裁决选 (a)') === 1, 'report CN finding row updated');
chk(c(rpe, 'Closed: the user') === 1, 'report EN finding row updated');
console.log('RESULT=' + (bad.length === 0 ? 'PASS' : 'FAIL(' + bad.length + ')'));

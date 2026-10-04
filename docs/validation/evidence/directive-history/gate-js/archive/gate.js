// §6.3 gate runner (master-side; re-created per session because tools/gates/ scripting is still a pending slice).
const fs = require('fs'), cp = require('child_process');
const R = 'D:/LZProjects/prfrail/';
const run = c => cp.execSync(c, { cwd: R, encoding: 'utf8', maxBuffer: 2e8 });
const runB = c => cp.execSync(c, { cwd: R, maxBuffer: 2e8 });
const res = []; const rec = (n, ok, d) => res.push({ n, ok, d });
const rd = p => fs.readFileSync(R + p, 'utf8').replace(/^\uFEFF/, '').replace(/\r\n/g, '\n').split('\n');
const cnt = (a, re) => a.reduce((n, l) => n + (l.match(re) || []).length, 0);
const changed = run('git status --porcelain').split('\n').filter(Boolean).map(l => l.slice(3).trim().replace(/^"|"$/g, ''));
const mdC = changed.filter(f => f.endsWith('.md'));
console.log('changed: ' + JSON.stringify(changed));
const added = {};
for (const f of changed) { let d = ''; try { d = run('git diff -U0 -- "' + f + '"'); } catch (e) { } added[f] = d.split('\n').filter(l => l.startsWith('+') && !l.startsWith('+++')).map(l => l.slice(1)); }

// G1-a
let g1a = [];
for (const f of mdC) for (const l of added[f]) if (l.startsWith('|') && l.split('\\|').join('').includes('||')) g1a.push(f + ' :: ' + l.slice(0, 50));
rec('G1-a table structure', g1a.length === 0, g1a.length ? JSON.stringify(g1a) : 'no `||` on added rows');

// G1-b
let g1b = [];
for (const f of mdC) {
  const L = rd(f), txt = L.join('\n');
  const hv = (/^版本：v(\d+\.\d+)/m.exec(txt) || /^Version: v(\d+\.\d+)/m.exec(txt) || [])[1];
  if (!hv) continue;
  const hd = (/^版本：v[\d.]+[^\n]{0,4}日期：(\d{4}-\d{2}-\d{2})/m.exec(txt) || /^Version: v[\d.]+\. Date: (\d{4}-\d{2}-\d{2})/m.exec(txt) || [])[1];
  const s = L.findIndex(l => /^### 0\.3/.test(l)), e = L.findIndex((l, i) => i > s && /^#{2,3} /.test(l));
  const rows = L.slice(s, e < 0 ? L.length : e).filter(l => /^\| v\d+\.\d+ \|/.test(l));
  const last = rows.length ? /^\| v(\d+\.\d+) \|/.exec(rows[rows.length - 1])[1] : '?';
  const ld = rows.length ? (/^\| v[\d.]+ \| (\d{4}-\d{2}-\d{2}) \|/.exec(rows[rows.length - 1]) || [])[1] : '?';
  const ok = hv === last && (!hd || hd === ld);
  g1b.push(f + ' header=v' + hv + ' last=v' + last + ' ' + hd + '/' + ld + ' rows=' + rows.length + (ok ? ' OK' : ' MISMATCH'));
  if (!ok) rec('G1-b ' + f, false, g1b[g1b.length - 1]);
}
if (g1b.length) rec('G1-b version metadata', g1b.every(s => / rows=\d+ OK$/.test(s)), g1b.join(' | '));

// G2
const pats = [/待[^\n]{0,6}回填/, /尚未提交/, /待同轮授权/, /pending the push/];
let g2 = [], g2i = [];
for (const f of mdC) rd(f).forEach((l, i) => { if (!pats.some(p => p.test(l)) || /不得残留/.test(l)) return; (added[f].includes(l) ? g2 : g2i).push(f + ':' + (i + 1)); });
rec('G2 placeholder residue', g2.length === 0, (g2.length ? JSON.stringify(g2) : 'none') + (g2i.length ? '  [INFO pre-existing: ' + g2i.join(', ') + ']' : ''));

// G3
const numstat = f => {
  const st = (run('git status --porcelain -- "' + f + '"') || '').trim();
  if (st.startsWith('??')) return [String(rd(f).length), '0'];   // untracked new file: all lines are additions
  let line = '';
  try { line = run('git diff --numstat -- "' + f + '"').trim().split('\n')[0]; } catch (e) { }
  const p = line.split('\t');
  return [p[0], p[1]];
};
const pairs = changed.filter(f => f.endsWith('.md') && !/_EN\.md$/.test(f) && changed.includes(f.replace(/\.md$/, '_EN.md')));
let sym = [], inf = [];
for (const cnf of pairs) {
  const enf = cnf.replace(/\.md$/, '_EN.md');
  const a = numstat(cnf), b = numstat(enf);
  const okS = a[0] === b[0] && a[1] === b[1];
  const CL = rd(cnf), EL = rd(enf);
  let dr = 0; for (let i = 0; i < Math.max(CL.length, EL.length); i++) if (/^#{2,4} /.test(CL[i] || '') !== /^#{2,4} /.test(EL[i] || '')) dr++;
  const f5 = { lines: CL.length === EL.length, head: dr === 0, bold: cnt(CL, /\*\*/g) === cnt(EL, /\*\*/g), bars: cnt(CL, /\|/g) === cnt(EL, /\|/g), blank: CL.filter(l => l === '').length === EL.filter(l => l === '').length };
  sym.push(cnf + ' +' + a[0] + '/-' + a[1] + ' vs _EN +' + b[0] + '/-' + b[1] + (okS ? ' OK' : ' MISMATCH'));
  inf.push(cnf + ' ' + JSON.stringify(f5));
  if (!okS) rec('G3(a.1) ' + cnf, false, sym[sym.length - 1]);
  if (/^docs\/DELIVERY_DIRECTIVE\.md$/.test(cnf) && !Object.values(f5).every(Boolean)) rec('G3(a.2) ' + cnf, false, inf[inf.length - 1]);
}
if (sym.length) rec('G3(a) symmetry + mirror parity', sym.every(s => / OK$/.test(s)), sym.join(' || ') + ' ;; ' + inf.join(' ;; '));
const PB = [[/\b[0-9a-f]{7}\b/g, /\b[0-9a-f]{7}\b/g, 'hash'], [/\brun \d+\b/g, /\brun \d+\b/g, 'run'], [/[✅⬜]/g, /[✅⬜]/g, 'glyph'], [/\*\*待办\*\*/g, /\*\*To do\*\*/g, 'to-do'], [/^\| \*\*未过独立终审\*\*/gm, /^\| \*\*Did not pass independent final review\*\*/gim, 'not-reviewed']];
let g3b = [];
for (const cnf of pairs) { const A = (added[cnf] || []).join('\n'), B = (added[cnf.replace(/\.md$/, '_EN.md')] || []).join('\n'); for (const [rc, re, lb] of PB) { const c = (A.match(rc) || []).length, e = (B.match(re) || []).length; g3b.push(lb + ' ' + c + '/' + e + (c === e ? ' OK' : ' MISMATCH')); } }
if (pairs.length) rec('G3(b) key-field alignment', g3b.every(s => / OK$/.test(s)), g3b.join(' | '));

// G4a / G4b
let g4a = [];
for (const f of changed) { if (/^\.github\/agents\/sol-orchestrator\//.test(f)) continue; const b = fs.readFileSync(R + f); const bom = b[0] === 0xEF && b[1] === 0xBB && b[2] === 0xBF, crlf = b.includes(Buffer.from('\r\n')); const want = /\.(md|ps1)$/i.test(f); const ok = want ? (bom && !crlf) : (!bom && !crlf); if (!ok) g4a.push(f + ' bom=' + bom + ' crlf=' + crlf); }
rec('G4a BOM+LF', g4a.length === 0, g4a.length ? JSON.stringify(g4a) : changed.length + ' file(s) OK');
const cjk = s => { const o = new Set(); for (const ch of s) { const c = ch.codePointAt(0); if ((c >= 0x4e00 && c <= 0x9fff) || (c >= 0x3400 && c <= 0x4dbf) || (c >= 0xf900 && c <= 0xfaff)) o.add(ch); } return o; };
let corpus = new Set();
for (const f of run('git ls-files "*.md"').split('\n').filter(Boolean)) { try { corpus = new Set([...corpus, ...cjk(run('git show HEAD:"' + f + '"'))]); } catch (e) { } }
const white = new Set();
fs.readFileSync(R + 'tools/gates/cjk-newwords.txt', 'utf8').split('\n').forEach(l => { const m = /^([^\s#].*?)\t/.exec(l); if (m) white.add(m[1]); });
let newC = new Set(); for (const f of mdC) for (const l of added[f]) for (const c of cjk(l)) newC.add(c);
const unk = [...newC].filter(c => !corpus.has(c) && !white.has(c));
rec('G4b CJK diff minus whitelist', unk.length === 0, 'newCjk=' + newC.size + ' unknown=' + JSON.stringify(unk) + ' (whitelist=' + white.size + ')');

// G5 / R2.5
const A = R + '.github/agents/'; const agents = fs.readdirSync(A).filter(f => f.endsWith('.agent.md'));
let h = ''; try { h = run('git grep -l "^model:" -- ".github/agents/prfrail-*.agent.md"'); } catch (e) { }
rec('G5-a(1) no model key in prfrail-*', h.trim() === '', h.trim() === '' ? '0 hits' : h);
let e2 = [];
for (const f of agents) { const fm = (/^---\n([\s\S]*?)\n---/.exec(fs.readFileSync(A + f, 'utf8').replace(/^\uFEFF/, '').replace(/\r\n/g, '\n')) || [])[1] || ''; if (/^\s*model\s*:\s*$/m.test(fm) || /model:\s*""/.test(fm) || /model:\s*''/.test(fm)) e2.push(f); }
rec('G5-a(2) no empty model', e2.length === 0, e2.length ? JSON.stringify(e2) : '0 hits / ' + agents.length + ' files');
const BL = /luna|terra|gemini|gpt-5\.4|sol|kimi/i; let g5b = [];
for (const f of changed.filter(f => f.endsWith('.agent.md'))) { const v = (/^model\s*:\s*(.+)$/m.exec(fs.readFileSync(R + f, 'utf8').replace(/^\uFEFF/, '')) || [])[1] || ''; if (v && BL.test(v)) g5b.push(f); }
rec('G5-b blacklist model value', g5b.length === 0, 'checked ' + changed.filter(f => f.endsWith('.agent.md')).length + ' file(s)');
let bad = [], notes = new Set();
for (const f of agents) { const t = fs.readFileSync(A + f, 'utf8').replace(/^\uFEFF/, '').replace(/\r\n/g, '\n'); const fm = (/^---\n[\s\S]*?\n---\n/.exec(t) || [])[0] || ''; const first = (t.slice(fm.length).split('\n').find(l => l.trim() !== '') || '').trim(); notes.add(first); if (!/非生成物/.test(first) || !/不得由生成器覆盖/.test(first)) bad.push(f); }
rec('R2.5 header-note position', bad.length === 0 && notes.size === 1, bad.length ? JSON.stringify(bad) : 'distinct=' + notes.size + ' uniform ' + agents.length + '/' + agents.length);
const gf = run('gofmt -l .').trim(); rec('base gofmt -l .', gf === '', gf === '' ? 'empty' : gf);
run('go build ./...'); rec('base go build', true, 'exit 0');
run('go vet ./...'); rec('base go vet', true, 'exit 0');
const t = run('go test ./...'); const o = (t.match(/^ok\s/gm) || []).length, fa = (t.match(/^FAIL/gm) || []).length;
rec('base go test', fa === 0, 'ok=' + o + ' FAIL=' + fa);
console.log('\n================ GATE REPORT ================');
res.forEach(r => console.log((r.ok ? 'PASS  ' : 'FAIL  ') + r.n + '\n        ' + r.d));
console.log('TOTAL_FAIL=' + res.filter(r => !r.ok).length);

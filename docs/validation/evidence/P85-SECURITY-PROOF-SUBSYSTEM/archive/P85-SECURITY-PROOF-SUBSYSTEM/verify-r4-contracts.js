// P85 ②-R4 verification: shape, BOM/LF, prefix-preservation, hunks, quotes.
const fs = require('fs');
const child = require('child_process');

function bufCheck(p) {
  const b = fs.readFileSync(p);
  const hasBom = b[0] === 0xef && b[1] === 0xbb && b[2] === 0xbf;
  const doubleBom = b[0] === 0xef && b[5] === 0xbf;
  const cr = b.includes(0x0d);
  const t = b.toString('utf8').replace(/^\uFEFF+/, '');
  const L = t.split('\n');
  const count = re => L.reduce((n, s) => n + (s.match(re) || []).length, 0);
  return {
    path: p, bom: hasBom, doubleBom, cr, lines: L.length,
    blank: L.filter(s => s === '').length,
    bold: count(/\*\*/g), bars: count(/\|/g), backticks: count(/`/g),
    trailingNewline: t.endsWith('\n'),
  };
}

function firstLast(p, lineNo) {
  const L = fs.readFileSync(p, 'utf8').split('\n');
  const s = L[lineNo - 1] || '';
  return { first: s.slice(0, 60), last: s.slice(-60), len: s.length };
}

function prefixProof(p, lineNo, idxText) {
  const idxL = idxText.split('\n');
  const wtL = fs.readFileSync(p, 'utf8').replace(/^\uFEFF+/, '').split('\n');
  // Only the T3 line changed in place; find it by the 60-char head of the index line.
  const head = idxL[lineNo - 1].slice(0, 60);
  const hit = wtL.findIndex(s => s.startsWith(head));
  return { indexLine: lineNo, worktreeLine: hit + 1, oldLen: idxL[lineNo - 1].length, newLen: wtL[hit].length, preserved: wtL[hit].startsWith(idxL[lineNo - 1]) };
}

const cn = bufCheck('docs/CONTRACTS.md');
const en = bufCheck('docs/CONTRACTS_EN.md');
const idxCN = child.execSync('git show :docs/CONTRACTS.md', { encoding: 'utf8' }).replace(/^\uFEFF+/, '');
const idxEN = child.execSync('git show :docs/CONTRACTS_EN.md', { encoding: 'utf8' }).replace(/^\uFEFF+/, '');
const ppCN = prefixProof('docs/CONTRACTS.md', 216, idxCN);
const ppEN = prefixProof('docs/CONTRACTS_EN.md', 210, idxEN);
const num = child.execSync('git diff --numstat -- docs/CONTRACTS.md docs/CONTRACTS_EN.md', { encoding: 'utf8' }).trim();
const hunks = child.execSync('git diff -U0 -- docs/CONTRACTS.md docs/CONTRACTS_EN.md', { encoding: 'utf8' })
  .split('\n').filter(l => l.startsWith('@@')).join(' | ');
console.log(JSON.stringify({ cn, en, delta: { lines: en.lines - cn.lines, blank: en.blank - cn.blank, bold: en.bold - cn.bold, bars: en.bars - cn.bars, backticks: en.backticks - cn.backticks } }, null, 2));
console.log('NUMSTAT:\n' + num);
console.log('HUNKS: ' + hunks);
console.log('PREFIX_CN=' + JSON.stringify(ppCN));
console.log('PREFIX_EN=' + JSON.stringify(ppEN));
console.log('CN216=' + JSON.stringify(firstLast('docs/CONTRACTS.md', 216)));
console.log('CN231=' + JSON.stringify(firstLast('docs/CONTRACTS.md', 231)));
console.log('CN233=' + JSON.stringify(firstLast('docs/CONTRACTS.md', 233)));
console.log('EN210=' + JSON.stringify(firstLast('docs/CONTRACTS_EN.md', 210)));
console.log('EN225=' + JSON.stringify(firstLast('docs/CONTRACTS_EN.md', 225)));
console.log('EN227=' + JSON.stringify(firstLast('docs/CONTRACTS_EN.md', 227)));

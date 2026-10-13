// ⑨ reference-resolvability check (doc-form native verification): every §X.Y reference on
// *added lines* of the staged change set must resolve to a heading or a bold line-leading label
// (in either the CN file or its _EN mirror), or to a resolvable parent prefix.
const { execFileSync } = require('child_process');
const fs = require('fs');

const diff = execFileSync('git', ['diff', '--cached', '-U0'], { encoding: 'utf8', maxBuffer: 1 << 28 });
const added = diff.split('\n').filter((l) => l.startsWith('+') && !l.startsWith('+++')).map((l) => l.slice(1));

const titles = new Set();
for (const f of ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md']) {
  const s = fs.readFileSync(f, 'utf8');
  for (const m of s.matchAll(/^#{2,4}\s+(\d+(?:\.\d+)*)/gm)) titles.add(m[1]);
  for (const m of s.matchAll(/^\*\*(\d+(?:\.\d+)*)\s/gm)) titles.add(m[1]);
}

const refs = [];
const unresolved = [];
for (const line of added) {
  for (const m of line.matchAll(/§(\d+(?:\.\d+)*)/g)) {
    const r = m[1];
    refs.push(r);
    let ok = titles.has(r);
    if (!ok) {
      const parts = r.split('.');
      for (let i = parts.length - 1; i >= 1; i -= 1) {
        if (titles.has(parts.slice(0, i).join('.'))) { ok = true; break; }
      }
    }
    if (!ok) unresolved.push(r);
  }
}

console.log('titles=' + titles.size + ' addedLines=' + added.length + ' refs=' + refs.length + ' distinct=' + new Set(refs).size);
console.log('unresolved=' + JSON.stringify([...new Set(unresolved)]));
console.log('RESULT=' + (unresolved.length === 0 ? 'PASS' : 'FAIL'));

'use strict';
// Normalize the report pair to the .md convention (UTF-8 with BOM + LF) and check the five mirror shapes
// that G3(a) compares: line count, heading lines, bold-marker count, pipe count, blank lines.
const fs = require('fs');
const CN = 'docs/validation/directive-haiku-bump.md';
const EN = 'docs/validation/directive-haiku-bump_EN.md';

function norm(p) {
  let s = fs.readFileSync(p, 'utf8').replace(/^\uFEFF/, '').replace(/\r\n/g, '\n');
  fs.writeFileSync(p, '\uFEFF' + s, 'utf8');
  return s;
}
const cn = norm(CN).split('\n');
const en = norm(EN).split('\n');
console.log('lines cn=' + cn.length + ' en=' + en.length + ' equal=' + (cn.length === en.length));
const shape = (l) => ({
  blank: l.trim() === '',
  head: /^#{1,6} /.test(l),
  bold: (l.match(/\*\*/g) || []).length,
  bars: (l.match(/\|/g) || []).length,
});
let bad = 0;
const n = Math.max(cn.length, en.length);
for (let i = 0; i < n; i++) {
  const a = cn[i] === undefined ? null : shape(cn[i]);
  const b = en[i] === undefined ? null : shape(en[i]);
  if (JSON.stringify(a) !== JSON.stringify(b)) {
    bad++;
    if (bad <= 15) console.log('MISMATCH line ' + (i + 1) + ' cn=' + JSON.stringify(a) + ' en=' + JSON.stringify(b));
  }
}
console.log('SHAPE_MISMATCHES=' + bad);
// byte-level check of the encodings
for (const p of [CN, EN]) {
  const buf = fs.readFileSync(p);
  const bom = buf[0] === 0xef && buf[1] === 0xbb && buf[2] === 0xbf;
  const crlf = buf.includes(Buffer.from('\r\n'));
  console.log(p + ' bom=' + bom + ' crlf=' + crlf + ' bytes=' + buf.length);
}

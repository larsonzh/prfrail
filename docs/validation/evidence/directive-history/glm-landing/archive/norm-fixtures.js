'use strict';
// One-shot normalizer for the four G5-b GLM fixtures: UTF-8 without BOM + LF.
const fs = require('fs');
const path = require('path');
const dir = path.join(__dirname, '..', '..', 'tools', 'gates', 'testdata', 'fixtures');
const files = [
  'g5b-closed-set-legal-glm.txt',
  'g5b-closed-set-legal-glm-flash.txt',
  'g5b-illegal-glm-case.txt',
  'g5b-illegal-glm-qualified.txt',
];
for (const f of files) {
  const p = path.join(dir, f);
  let s = fs.readFileSync(p).toString('utf8');
  if (s.charCodeAt(0) === 0xfeff) s = s.slice(1);
  s = s.replace(/\r\n/g, '\n');
  fs.writeFileSync(p, Buffer.from(s, 'utf8'));
  const nb = fs.readFileSync(p);
  console.log(f, 'bom=' + (nb[0] === 0xef && nb[1] === 0xbb), 'crlf=' + nb.includes(Buffer.from('\r\n')), 'bytes=' + nb.length);
}

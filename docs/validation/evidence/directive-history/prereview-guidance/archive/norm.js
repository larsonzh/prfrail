// Normalize a newly created .md to UTF-8 with BOM + LF (repo hard rule for .md).
const fs = require('fs');
const p = process.argv[2];
let s = fs.readFileSync(p, 'utf8');
s = s.replace(/^\uFEFF/, '').replace(/\r\n/g, '\n');
s = '\uFEFF' + s;
fs.writeFileSync(p, s);
const b = fs.readFileSync(p);
const isBom = b[0] === 0xef && b[1] === 0xbb && b[2] === 0xbf;
console.log(p + '  BOM=' + isBom + '  CR=' + b.includes(0x0d, 3) + '  bytes=' + b.length);

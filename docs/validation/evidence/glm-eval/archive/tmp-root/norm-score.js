// Normalise a transcribed leg artifact to LF, then hash + score + contamination-check it.
const fs = require('fs');
const path = require('path');
const cp = require('child_process');
const ROOT = 'd:/LZProjects/prfrail';
const GE = path.join(process.env.TEMP, 'glm-eval');

const f = process.argv[2];
let t = fs.readFileSync(f, 'utf8');
const before = Buffer.byteLength(t, 'utf8');
t = t.replace(/\r\n/g, '\n');
fs.writeFileSync(f, t, 'utf8');
const b = fs.readFileSync(f);
console.log('normalised ' + path.basename(f) + ' bytes=' + b.length + ' (was ' + before + ') crlf=' + b.includes(Buffer.from('\r\n')) + ' bom=' + (b[0] === 0xef));
console.log(cp.execSync('node "' + path.join(GE, 'hash.js') + '" "' + f + '"', { encoding: 'utf8' }).trim());
console.log('--- score ---');
console.log(cp.execSync('node "' + path.join(GE, 'score.js') + '" "' + f + '"', { encoding: 'utf8' }).trim());
console.log('--- contamination ---');
try {
  console.log(cp.execSync('node "' + path.join(GE, 'contam.js') + '" "' + f + '"', { encoding: 'utf8' }).trim());
} catch (e) {
  console.log('EXIT=' + e.status + ' ' + ((e.stdout || '') + (e.stderr || '')).trim());
}
const lines = t.split('\n').filter((l) => l.trim());
console.log('--- last non-empty line ---');
console.log(JSON.stringify(lines[lines.length - 1]));

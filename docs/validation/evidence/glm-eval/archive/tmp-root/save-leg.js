// Save the leg output byte-exactly (copy of the tool's raw capture), then hash + score + contaminate-check.
const fs = require('fs');
const path = require('path');
const cp = require('child_process');
const ROOT = 'd:/LZProjects/prfrail';
const GE = path.join(process.env.TEMP, 'glm-eval');
const OUT = path.join(ROOT, 'tmp', 'glm-eval-out');
fs.mkdirSync(OUT, { recursive: true });

const src = process.argv[2];
if (!src) { console.log('usage: node save-leg.js <raw-capture-file> <leg-name>'); process.exit(2); }
const leg = process.argv[3] || 'leg-glm';
const dest = path.join(OUT, leg + '.txt');
fs.copyFileSync(src, dest);
const b = fs.readFileSync(dest);
console.log('saved ' + dest + ' bytes=' + b.length + ' bom=' + (b[0] === 0xef) + ' crlf=' + b.includes(Buffer.from('\r\n')));
console.log(cp.execSync('node "' + path.join(GE, 'hash.js') + '" "' + dest + '"', { encoding: 'utf8' }).trim());
console.log('--- score ---');
console.log(cp.execSync('node "' + path.join(GE, 'score.js') + '" "' + dest + '"', { encoding: 'utf8' }).trim());
console.log('--- contamination ---');
try {
  console.log(cp.execSync('node "' + path.join(GE, 'contam.js') + '" "' + dest + '"', { encoding: 'utf8' }).trim());
} catch (e) {
  console.log('EXIT=' + e.status + ' ' + ((e.stdout || '') + (e.stderr || '')).trim());
}
console.log('--- last non-empty line ---');
const lines = b.toString('utf8').split('\n').filter((l) => l.trim());
console.log(JSON.stringify(lines[lines.length - 1]));

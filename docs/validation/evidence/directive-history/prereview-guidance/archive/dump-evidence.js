// ④ evidence dump (§6.2): persist the raw staged diff and the raw gate output to tmp/ (gitignored).
const { execFileSync } = require('child_process');
const fs = require('fs');

const diff = execFileSync('git', ['diff', '--cached'], { encoding: 'utf8', maxBuffer: 1 << 28 });
fs.writeFileSync('tmp/prg/diff-step0-2.diff', diff, 'utf8');

const gates = execFileSync('node', ['tools/gates/gate.js', '--all', '--scope=index'], { encoding: 'utf8', maxBuffer: 1 << 28 });
fs.writeFileSync('tmp/prg/gates-index-step0-2.txt', gates, 'utf8');

console.log('diff bytes=' + Buffer.byteLength(diff, 'utf8') + ' lines=' + diff.split('\n').length);
console.log('gate lines=' + gates.split('\n').length);
console.log('TOTAL_FAIL :: ' + (gates.match(/TOTAL_FAIL=\d+/g) || []).join(' | '));
console.log('PASS count=' + ((gates.match(/^PASS /gm) || []).length) + ' FAIL count=' + ((gates.match(/^FAIL /gm) || []).length));

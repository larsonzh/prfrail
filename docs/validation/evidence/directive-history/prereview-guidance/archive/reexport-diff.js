// Re-export the staged diff byte-exactly (⑤ found transcription noise in the earlier dump) and verify.
const { execFileSync } = require('child_process');
const fs = require('fs');

const fresh = execFileSync('git', ['diff', '--cached'], { encoding: 'utf8', maxBuffer: 1 << 28 });
const pathBefore = 'tmp/prg/diff-step0-2.diff';
const before = fs.existsSync(pathBefore) ? fs.readFileSync(pathBefore, 'utf8') : '';

fs.writeFileSync(pathBefore, fresh, 'utf8');
const after = fs.readFileSync(pathBefore, 'utf8');

const c = (s, re) => (s.match(re) || []).length;
console.log('fresh === written : ' + (fresh === after));
console.log('before === fresh  : ' + (before === fresh));
console.log('fresh: lines=' + fresh.split('\n').length + '  "+---"=' + c(fresh, /^\+---$/gm) + '  " ---"(ctx)=' + c(fresh, /^ ---$/gm) + '  "^+$"(blank-add)=' + c(fresh, /^\+$/gm));
console.log('before: lines=' + before.split('\n').length + '  "+---"=' + c(before, /^\+---$/gm) + '  " ---"(ctx)=' + c(before, /^ ---$/gm) + '  "^+$"(blank-add)=' + c(before, /^\+$/gm));
console.log('sha256(fresh)=' + require('crypto').createHash('sha256').update(fresh, 'utf8').digest('hex'));

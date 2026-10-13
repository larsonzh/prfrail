'use strict';
// T028 review-input materialization (per the user's kickoff requirement + the HAIKU-BUMP C-7 lesson):
//   head-baseline.txt : per changed path, HEAD vs worktree identity (sha256 / bytes / lines / status)
//   staged-diff.txt   : the complete staged diff
// Note on form: this is a 27-file, ~6k-line code change, so a full dual-TEXT dump would be megabytes;
// the baseline therefore carries per-file hashes + counts (verifiable identity) instead of full text,
// and the complete staged diff is supplied separately. The adaptation is registered in the report.
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const { execSync } = require('child_process');

const ROOT = path.resolve(__dirname, '..', '..');
const OUT_BASE = path.join(ROOT, 'tmp/t028/head-baseline.txt');
const OUT_DIFF = path.join(ROOT, 'tmp/t028/staged-diff.txt');

const git = (args, opts = {}) => execSync('git ' + args, { cwd: ROOT, maxBuffer: 1 << 28, ...opts });
const sha256 = (buf) => crypto.createHash('sha256').update(buf).digest('hex');
const lines = (buf) => (buf.length === 0 ? 0 : buf.toString('utf8').split('\n').length);

const staged = git('diff --cached --name-status -z').toString('utf8').split('\0').filter(Boolean);
const pairs = [];
for (let i = 0; i < staged.length; i += 2) pairs.push([staged[i], staged[i + 1]]);

const out = [];
out.push('# T028 review-input baseline (generated ' + new Date().toISOString() + ')');
out.push('# HEAD = ' + git('rev-parse HEAD').toString().trim());
out.push('# Per changed path: staged status, HEAD identity, worktree identity.');
out.push('# "HEAD absent" means the path did not exist at HEAD (a newly added file).');
out.push('');
out.push(['status', 'path', 'head_sha256', 'head_bytes', 'head_lines', 'work_sha256', 'work_bytes', 'work_lines'].join('\t'));
let added = 0;
let modified = 0;
let deleted = 0;
for (const [status, p] of pairs) {
  const rel = p.replace(/\\/g, '/');
  let headBuf = null;
  try {
    headBuf = git('show HEAD:' + JSON.stringify(rel).replace(/"/g, '"'), { stdio: ['pipe', 'pipe', 'ignore'] });
  } catch (e) {
    headBuf = null;
  }
  const full = path.join(ROOT, rel);
  const workBuf = fs.existsSync(full) ? fs.readFileSync(full) : null;
  if (status.startsWith('A')) added++;
  else if (status.startsWith('D')) deleted++;
  else modified++;
  out.push(
    [
      status,
      rel,
      headBuf ? sha256(headBuf) : 'HEAD-absent',
      headBuf ? headBuf.length : '-',
      headBuf ? lines(headBuf) : '-',
      workBuf ? sha256(workBuf) : 'worktree-absent',
      workBuf ? workBuf.length : '-',
      workBuf ? lines(workBuf) : '-',
    ].join('\t')
  );
}
out.push('');
out.push('# totals: added=' + added + ' modified=' + modified + ' deleted=' + deleted + ' changed=' + pairs.length);

const diff = git('diff --cached');
fs.writeFileSync(OUT_BASE, out.join('\n'), 'utf8');
fs.writeFileSync(OUT_DIFF, diff);
console.log('wrote ' + path.relative(ROOT, OUT_BASE) + ' lines=' + out.length);
console.log('wrote ' + path.relative(ROOT, OUT_DIFF) + ' bytes=' + diff.length);
console.log('changed=' + pairs.length + ' added=' + added + ' modified=' + modified + ' deleted=' + deleted);

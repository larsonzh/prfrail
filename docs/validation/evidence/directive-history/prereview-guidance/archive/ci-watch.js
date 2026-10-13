// Bounded CI observation: wait for the run to complete, then extract per-leg hard evidence
// (SELFTEST / GATE REPORT scope=ci / TOTAL_FAIL / package-level ok lines), per the CI evidence standard.
const { execFileSync } = require('child_process');
const fs = require('fs');

const RUN = process.argv[2];
const env = Object.assign({}, process.env, { GH_PAGER: '' });
const gh = (args, opts) => execFileSync('gh', args, Object.assign({ encoding: 'utf8', maxBuffer: 1 << 28, env }, opts || {}));
const sleep = (ms) => Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, ms);

let info = null;
for (let i = 0; i < 40; i += 1) {
  info = JSON.parse(gh(['run', 'view', RUN, '--json', 'status,conclusion,headSha,url']));
  if (info.status === 'completed') break;
  sleep(10000);
}
console.log('run=' + RUN + ' status=' + info.status + ' conclusion=' + info.conclusion + ' headSha=' + info.headSha.slice(0, 7));
console.log('url=' + info.url);

if (info.status === 'completed') {
  const jobs = JSON.parse(gh(['run', 'view', RUN, '--json', 'jobs'])).jobs;
  for (const j of jobs) {
    console.log('JOB ' + j.name + ' :: ' + j.status + '/' + j.conclusion + ' (' + (j.startedAt || '?') + ' -> ' + (j.completedAt || '?') + ')');
  }
  // per-leg hard evidence from the logs
  const dir = 'tmp/prg/ci-' + RUN;
  fs.mkdirSync(dir, { recursive: true });
  for (const j of jobs) {
    let log = '';
    try { log = gh(['run', 'view', RUN, '--job', String(j.databaseId), '--log']); } catch (e) { log = ''; }
    fs.writeFileSync(dir + '/' + String(j.name).replace(/[^A-Za-z0-9._-]/g, '_') + '.log', log, 'utf8');
    const hits = log.split('\n').filter((l) => /SELFTEST:|GATE REPORT scope=ci|TOTAL_FAIL=|^[^ ]*ok  |ok [a-z0-9_./-]+ /.test(l)).slice(0, 14);
    console.log('--- ' + j.name + ' key lines ---');
    for (const h of hits) console.log('  ' + h.trim().slice(0, 200));
  }
}

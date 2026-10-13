'use strict';
/*
 * T028 CI evidence extractor (local-only, tmp/ artifact).
 *
 * Reads the raw `gh run view <id> --log` dump from tmp/t028/ci-log.txt and prints, per
 * job, only the load-bearing evidence lines:
 *   - the Gates step readings (SELFTEST / GATE REPORT scope=ci / TOTAL_FAIL)
 *   - package-level `ok <pkg>` and `FAIL` lines, which are what prove a package ran
 *     (a green job conclusion alone does not).
 * The `repo=` field on each GATE line identifies the leg (Ubuntu vs Windows runner).
 */

const fs = require('fs');

const file = process.argv[2] || 'tmp/t028/ci-log.txt';
const lines = fs.readFileSync(file, 'utf8').split('\n');
const jobs = new Map();

for (const line of lines) {
  const parts = line.split('\t');
  if (parts.length < 3) continue;
  const job = parts[0].trim();
  const rest = parts
    .slice(2)
    .join('\t')
    .trim()
    .replace(/^\d{4}-\d{2}-\d{2}T\S+\s+/, '');
  if (!jobs.has(job)) jobs.set(job, []);
  jobs.get(job).push(rest);
}

let jobsWithEvidence = 0;
for (const [job, body] of jobs) {
  const gates = body.filter((l) => /SELFTEST: PASS|GATE REPORT scope=|TOTAL_FAIL=/.test(l));
  const okPkgs = body.filter((l) => /^ok\s+github\.com\/larsonzh\/prfrail\S*/.test(l));
  const failPkgs = body.filter((l) => /^(FAIL|--- FAIL|# github)/.test(l));
  if (gates.length === 0 && okPkgs.length === 0 && failPkgs.length === 0) continue;
  jobsWithEvidence++;
  console.log('=== JOB: ' + job);
  for (const g of gates) console.log('  GATE  ' + g);
  console.log('  OK_PKG_LINES=' + okPkgs.length);
  for (const p of okPkgs) console.log('    ' + p);
  console.log('  FAIL_PKG_LINES=' + failPkgs.length);
  for (const p of failPkgs) console.log('    ' + p);
}

console.log('');
console.log('JOBS_WITH_EVIDENCE=' + jobsWithEvidence);
console.log('TOTAL_LOG_LINES=' + lines.length);

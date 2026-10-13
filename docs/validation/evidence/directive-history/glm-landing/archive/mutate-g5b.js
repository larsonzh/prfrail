'use strict';
// One-shot falsifiability harness for the T61-T67 assertions (slice DIRECTIVE-GLM-LANDING,
// role 3). Each case mutates tools/gates/lib/criteria/g5b.js ON DISK, runs
// `gate.js --selftest` in a FRESH process, captures the T32/T61-T67 verdicts + the summary
// line, then restores the original bytes in a `finally` and proves byte-identity by sha256.
// It never commits and never leaves the criterion mutated.
const fs = require('fs');
const path = require('path');
const cp = require('child_process');
const crypto = require('crypto');

const ROOT = path.join(__dirname, '..', '..');
const G5B = path.join(ROOT, 'tools', 'gates', 'lib', 'criteria', 'g5b.js');
const GATE = path.join(ROOT, 'tools', 'gates', 'gate.js');
const sha = (b) => crypto.createHash('sha256').update(b).digest('hex');

const MEMBER_GLM = "  'glm-5.3',\n";
const MEMBER_GLM_FLASH = "  'glm-5.3-flash',\n";
const COMPARE = '        if (!ALLOWED.has(v)) findings.push({ file: f, value: v });';

const CASES = [
  { id: 'drop-glm', from: MEMBER_GLM, to: '' },
  { id: 'drop-glm-flash', from: MEMBER_GLM_FLASH, to: '' },
  {
    id: 'case-insensitive',
    from: COMPARE,
    to: '        if (!ALLOWED.has(v.toLowerCase())) findings.push({ file: f, value: v });',
  },
  { id: 'add-qualified', from: MEMBER_GLM_FLASH, to: MEMBER_GLM_FLASH + "  'glm-5.3 (glm)',\n" },
];

const original = fs.readFileSync(G5B);
const beforeSha = sha(original);
console.log('shipped g5b.js sha256=' + beforeSha);
const WANT = /^FAIL |^ok   T(?:32|61|62|63|64|65|66|67) |^SELFTEST:/;

for (const c of CASES) {
  const src = original.toString('utf8');
  const occurrences = src.split(c.from).length - 1;
  const mutated = src.split(c.from).join(c.to);
  console.log('\n=== CASE ' + c.id + ' :: anchorOccurrences=' + occurrences + ' applied=' + (mutated !== src) + ' ===');
  if (mutated === src) {
    console.log('SKIP: anchor not found');
    continue;
  }
  fs.writeFileSync(G5B, Buffer.from(mutated, 'utf8'));
  let out = '';
  try {
    const r = cp.spawnSync(process.execPath, [GATE, '--selftest'], {
      cwd: ROOT,
      encoding: 'utf8',
      maxBuffer: 64 * 1024 * 1024,
    });
    out = (r.stdout || '') + (r.stderr || '');
  } finally {
    fs.writeFileSync(G5B, original);
  }
  const afterSha = sha(fs.readFileSync(G5B));
  for (const line of out.split('\n')) if (WANT.test(line)) console.log(line);
  console.log('restored byte-identical=' + (afterSha === beforeSha) + ' sha256=' + afterSha);
}

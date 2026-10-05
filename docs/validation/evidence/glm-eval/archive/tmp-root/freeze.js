// ④ freeze: build the exact leg input, verify it against the pre-registered input_hash,
// fill the dataset pre_registration block, then move the answer-bearing platform out of the workspace.
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const os = require('os');

const ROOT = 'd:/LZProjects/prfrail';
const GE = path.join(ROOT, 'tmp', 'glm-eval');
const ORDER = ['00-plan.md', '01-clause-cn.txt', '02-diff.js', '03-prompt.txt'];
const sha = (b) => crypto.createHash('sha256').update(b).digest('hex');
const log = (s) => console.log(s);

// 1. exact leg input = byte concatenation in fixed order
const bufs = ORDER.map((n) => fs.readFileSync(path.join(GE, 'bundle', n)));
const concat = Buffer.concat(bufs);
const inputHash = sha(concat);
log('input_hash(recomputed) = ' + inputHash);
log('matches PRE-REGISTERED  = ' + (inputHash === 'bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6'));
fs.writeFileSync(path.join(ROOT, 'tmp', 'PROMPT-EXACT.txt'), concat);
log('wrote tmp/PROMPT-EXACT.txt bytes=' + concat.length + ' sha256=' + sha(fs.readFileSync(path.join(ROOT, 'tmp', 'PROMPT-EXACT.txt'))));

// 2. fill the dataset pre_registration block (hashes only, no seed metadata)
const dsPath = path.join(ROOT, 'docs', 'validation', 'directive-glm-eval-dataset.json');
const ds = JSON.parse(fs.readFileSync(dsPath, 'utf8'));
const keyBytes = fs.readFileSync(path.join(GE, 'key.json'));
ds.pre_registration.bundle_input_hash = inputHash;
ds.pre_registration.key_file_sha256 = sha(keyBytes);
ds.pre_registration.registered_at = '2026-10-05';
fs.writeFileSync(dsPath, JSON.stringify(ds, null, 2) + '\n', 'utf8');
log('pre_registration = ' + JSON.stringify(ds.pre_registration, null, 2));

// 3. move the answer-bearing platform OUT of the workspace
const dest = path.join(os.tmpdir(), 'glm-eval');
if (fs.existsSync(dest)) fs.rmSync(dest, { recursive: true, force: true });
fs.cpSync(GE, dest, { recursive: true });
fs.rmSync(GE, { recursive: true, force: true });
log('moved tmp/glm-eval -> ' + dest);
log('dest listing: ' + fs.readdirSync(dest).join(', '));
log('key.json still outside: ' + fs.existsSync(path.join(dest, 'key.json')));
log('tmp/glm-eval gone: ' + !fs.existsSync(GE));

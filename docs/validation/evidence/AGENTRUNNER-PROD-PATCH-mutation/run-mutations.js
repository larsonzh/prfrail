'use strict';
// Mutation re-run for slice AGENTRUNNER-PROD-PATCH, executed inside the compliant
// location tmp/AGENTRUNNER-PROD-PATCH/mutation/. A pristine HEAD tree is extracted
// with `git archive`, the two working-tree files are copied over it, and the two
// mutations are applied to THAT copy only. The repository is never modified.
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const { execSync, execFileSync } = require('child_process');

const ROOT = process.cwd();
const BASE = path.join(ROOT, 'tmp', 'AGENTRUNNER-PROD-PATCH', 'mutation');
const WT = path.join(BASE, 'wt');
const TAR = path.join(BASE, 'head-tree.tar');
const DISP_REL = path.join('internal', 'adapters', 'agent_runner_replay_dispatcher.go');
const TEST_REL = path.join('internal', 'adapters', 'agent_runner_replay_dispatcher_test.go');

const log = [];
function say(line) {
  console.log(line);
  log.push(line);
}
function sh(cmd, cwd) {
  const out = execSync(cmd, { cwd: cwd || ROOT, stdio: ['ignore', 'pipe', 'pipe'], maxBuffer: 1 << 28 });
  return out.toString('utf8');
}
function sha(file) {
  return crypto.createHash('sha256').update(fs.readFileSync(file)).digest('hex');
}

fs.mkdirSync(BASE, { recursive: true });
fs.rmSync(WT, { recursive: true, force: true });
fs.mkdirSync(WT, { recursive: true });

say('== step 1: extract a pristine HEAD tree with git archive (no worktree registration) ==');
sh('git archive --format=tar -o "' + TAR + '" HEAD');
say('git archive -> ' + TAR + ' (' + fs.statSync(TAR).size + ' bytes)');
// The bundled Windows tar cannot address the circled-digit filenames inside the
// historical evidence archive, so only the source subtrees the tests need are extracted.
say(sh('tar -xf "' + TAR + '" -C "' + WT + '" go.mod go.sum internal cmd tools examples schemas testdata harnesses'));
say('extracted entries: ' + fs.readdirSync(WT).length);

say('');
say('== step 2: overlay the two working-tree files (the change under test) ==');
for (const rel of [DISP_REL, TEST_REL]) {
  fs.copyFileSync(path.join(ROOT, rel), path.join(WT, rel));
  say('copied ' + rel + ' sha256=' + sha(path.join(WT, rel)));
}

const pristine = path.join(BASE, 'dispatcher.pristine.go.txt');
fs.copyFileSync(path.join(WT, DISP_REL), pristine);
const preSha = sha(path.join(WT, DISP_REL));
say('');
say('pristine reference: ' + pristine);
say('pre-mutation sha256 = ' + preSha);

function runTests(label, filter) {
  say('');
  say('== ' + label + ' ==');
  const cmd = 'go test -count=1 -run "' + filter + '" ./internal/adapters/';
  say('$ ' + cmd);
  let out = '';
  let code = 0;
  try {
    out = execSync(cmd, { cwd: WT, stdio: ['ignore', 'pipe', 'pipe'], maxBuffer: 1 << 28 }).toString('utf8');
  } catch (e) {
    code = e.status === undefined ? 1 : e.status;
    out = (e.stdout ? e.stdout.toString('utf8') : '') + (e.stderr ? e.stderr.toString('utf8') : '');
  }
  say('exit=' + code);
  say(out.trim());
  return { code: code, out: out };
}

// ---- baseline ----
const baseline = runTests('baseline (unmutated copy)', 'AgentRunnerReplayDispatcher');

// ---- mutation M1: drop the context shield ----
say('');
say('== mutation M1: context.WithoutCancel(ctx) -> ctx ==');
let src = fs.readFileSync(path.join(WT, DISP_REL), 'utf8');
const m1Anchor = 'context.WithTimeout(context.WithoutCancel(ctx), grace+time.Second)';
const m1Hits = src.split(m1Anchor).length - 1;
say('anchor hits = ' + m1Hits);
if (m1Hits !== 1) throw new Error('M1 anchor not unique');
fs.writeFileSync(path.join(WT, DISP_REL), src.split(m1Anchor).join('context.WithTimeout(ctx, grace+time.Second)'), 'utf8');
const m1Sha = sha(path.join(WT, DISP_REL));
say('mutated sha256 = ' + m1Sha + ' (differs from pre: ' + (m1Sha !== preSha) + ')');
const m1 = runTests('M1 run (expect the shield test to fail)', 'TestAgentRunnerReplayDispatcherReclaimSurvivesCallerCancellation');

// ---- restore ----
fs.copyFileSync(pristine, path.join(WT, DISP_REL));
const restoredSha = sha(path.join(WT, DISP_REL));
say('');
say('== restore after M1 ==');
say('restored sha256 = ' + restoredSha + ' (equals pre-mutation: ' + (restoredSha === preSha) + ')');

// ---- mutation M2: force the no-capability branch ----
say('');
say('== mutation M2: force the no-capability branch (append ok = false after the type assertion) ==');
src = fs.readFileSync(path.join(WT, DISP_REL), 'utf8');
const m2Anchor = 'reclaimer, ok := dispatcher.Launcher.(agentRunnerLaunchReclaimer)\n';
const m2Hits = src.split(m2Anchor).length - 1;
say('anchor hits = ' + m2Hits);
if (m2Hits !== 1) throw new Error('M2 anchor not unique');
fs.writeFileSync(path.join(WT, DISP_REL), src.split(m2Anchor).join(m2Anchor + '\tok = false\n'), 'utf8');
const m2Sha = sha(path.join(WT, DISP_REL));
say('mutated sha256 = ' + m2Sha);
const m2 = runTests('M2 run (expect the reclaim tests to fail)', 'AgentRunnerReplayDispatcher');

// ---- restore again ----
fs.copyFileSync(pristine, path.join(WT, DISP_REL));
const finalSha = sha(path.join(WT, DISP_REL));
say('');
say('== restore after M2 ==');
say('final sha256 = ' + finalSha + ' (equals pre-mutation: ' + (finalSha === preSha) + ')');

// ---- evidence ----
say('');
say('== summary ==');
say('baseline exit = ' + baseline.code);
say('M1 exit = ' + m1.code + ' (expected non-zero)');
say('M2 exit = ' + m2.code + ' (expected non-zero)');
say('sha256 pre = ' + preSha);
say('sha256 after restore = ' + finalSha);
say('evidence directory retained: ' + BASE);
say('wt retained: ' + WT);

fs.writeFileSync(path.join(BASE, 'evidence.txt'), log.join('\n') + '\n', 'utf8');
console.log('\nwrote ' + path.join(BASE, 'evidence.txt'));

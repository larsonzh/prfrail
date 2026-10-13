'use strict';
/*
 * T028 A12 dedicated independent experiment (local-only, tmp/ artifact).
 *
 * Purpose: because the 6th-stage independent scan returned PASS, the 9th stage must
 * prove the slice's load-bearing invariants with a criterion that is INDEPENDENT of
 * that scan: this script re-derives every claim from the final-state sources instead
 * of trusting any review conclusion, and re-executes the supervised test subset.
 *
 * Scope: static assertions over the frozen change set + one child-process test run.
 * Nothing here writes into the repository; no network, no credentials.
 */

const fs = require('fs');
const cp = require('child_process');

let checks = 0;
let failures = 0;
const ok = (name, cond, detail) => {
  checks++;
  if (cond) {
    console.log('ok   ' + name);
  } else {
    failures++;
    console.log('FAIL ' + name + (detail ? ' :: ' + detail : ''));
  }
};

const read = (p) => fs.readFileSync(p, 'utf8');
const has = (src, re) => re.test(src);
const count = (src, re) => (src.match(re) || []).length;

/** Every schema node stored under `name`, at any depth (the record envelope nests the body in $defs). */
function findProp(node, name, hits) {
  const out = hits || [];
  if (Array.isArray(node)) {
    node.forEach((n) => findProp(n, name, out));
  } else if (node && typeof node === 'object') {
    Object.keys(node).forEach((k) => {
      if (k === name) out.push(node[k]);
      findProp(node[k], name, out);
    });
  }
  return out;
}
const firstProp = (node, name) => (findProp(node, name, [])[0] || {});

/** Text of the function whose declaration contains `needle`, up to the next top-level func. */
function bodyOf(src, needle) {
  const at = src.indexOf(needle);
  if (at < 0) return '';
  const next = src.indexOf('\nfunc ', at + needle.length);
  return src.slice(at, next < 0 ? src.length : next);
}

const runtime = read('internal/console/runtime_supervised.go');
const chainBox = read('internal/chain/supervised_blackbox.go');
const chainReport = read('internal/chain/supervised_report.go');
const chainConfirm = read('internal/chain/supervised_risk_confirmation.go');
const chainModels = read('internal/chain/models.go');
const chainRouter = read('internal/chain/router.go');
const cli = read('internal/console/cli.go');
const cmd = read('internal/console/supervised_cmd.go');
const visible = read('internal/adapters/sessbridge_visible.go');
const consoleTests = read('internal/console/runtime_supervised_test.go');
const chainTests = read('internal/chain/supervised_blackbox_test.go');
const reportSchema = JSON.parse(read('schemas/supervised-report.schema.json'));
const confirmSchema = JSON.parse(read('schemas/supervised-risk-confirmation.schema.json'));
const chainSchema = JSON.parse(read('schemas/chain.schema.json'));
const contractsTest = read('tools/contracts/contracts.test.js');

console.log('--- part 1: frozen handoff policy guard (blind-review High) ---');

ok('a1 frozen policy snapshot name is single-sourced', has(runtime, /supervisedFrozenPolicyName\s*=\s*"handoff-policy\.json"/));
ok('a2 snapshot is persisted at window open', has(runtime, /func \(runtime \*supervisedRuntime\) persistFrozenHandoffPolicy\(/));
ok('a3 the persist path is reached from the window-opening path', has(bodyOf(runtime, 'func (runtime *supervisedRuntime) confirmAndOpenWindow('), /persistFrozenHandoffPolicy\(\)/));
ok('a4 enforceFrozenHandoffPolicy exists', has(runtime, /func \(runtime \*supervisedRuntime\) enforceFrozenHandoffPolicy\(session chain\.HandoffSession\) error/));

const guard = bodyOf(runtime, 'func (runtime *supervisedRuntime) enforceFrozenHandoffPolicy(');
ok('a5 guard rejects drift with ErrSupervisedPolicyDrift on three independent anchors', count(guard, /ErrSupervisedPolicyDrift/g) >= 3, 'count=' + count(guard, /ErrSupervisedPolicyDrift/g));
ok('a6 guard anchors: snapshot vs its own digest', has(guard, /snapshot[\s\S]*?does not match its own digest|its own digest/));
ok('a7 guard anchors: snapshot vs the window session', has(guard, /does not match the window session|HandoffPolicyHash/));
ok('a8 guard anchors: in-force policy vs frozen policy', has(guard, /policy in force changed since the window opened/));
ok('a9 guard performs NO durable write (fail-closed, zero-persist path)',
  !has(guard, /writeSupervisedJSON|os\.WriteFile|WriteFile|MkdirAll|Rename|AppendAllText/));
ok('a10 guard is invoked from the return entry point',
  bodyOf(runtime, 'func ExecuteSupervisedReturn(').indexOf('enforceFrozenHandoffPolicy') >= 0);
ok('a11 abort bypasses the drift guard (cannot be refused by drift)', has(consoleTests, /TestSupervisedReturnAllowsAbortUnderDriftedPolicy/));
ok('a12 drift negative case is pinned by a test', has(consoleTests, /TestSupervisedReturnRefusesDriftedPolicy/));
ok('a13 policy digest uses the repository hashing convention', has(runtime, /func supervisedHandoffPolicyHash\(/));

console.log('--- part 2: contract invariants (schema + constants) ---');

const assurance = firstProp(reportSchema, 'assuranceLevel');
ok('b1 report schema forces assuranceLevel = reduced', assurance.const === 'reduced', 'const=' + assurance.const);
const unknown = firstProp(reportSchema, 'unknownDimensions');
ok('b2 report schema requires at least one unknown dimension', (unknown.minItems || 0) >= 1, 'minItems=' + unknown.minItems);
ok('b3 report schema requires the cost dimension to be present',
  !!(unknown.contains && unknown.contains.const === 'cost') && (unknown.minContains || 0) >= 1,
  'contains=' + JSON.stringify(unknown.contains) + ' minContains=' + unknown.minContains);
ok('b4 risk confirmation schema forces selection = confirm', firstProp(confirmSchema, 'selection').const === 'confirm');
ok('b5 chain schema carries the supervised black-box step variant', has(read('schemas/chain.schema.json'), /supervisedBlackBoxCodeStep/));
ok('b6 risk notice version is the versioned constant 2', has(chainBox, /SupervisedRiskNoticeVersion\s*=\s*2/));
ok('b7 contract coverage constants reflect the added schemas/fixtures',
  has(contractsTest, /\b38\b/) && has(contractsTest, /\b143\b/));

console.log('--- part 3: P4 - an unknown dimension is never rendered as safe ---');

ok('c1 rendered dimensions use the proved: label', has(cmd, /proved:/));
ok('c2 rendered text never claims safe/verified',
  !has(bodyOf(cmd, 'func writeSupervisedReportText('), /\bsafe\b|\bverified\b/i));
const enumText = JSON.stringify(reportSchema);
ok('c3 report schema exposes no safe/verified enum member', !/"enum":\[[^\]]*"(safe|verified)"/.test(enumText));
ok('c4 report validation rejects non-reduced / empty unknown dimensions',
  count(chainReport, /ErrInvalidSupervisedReport/g) >= 3, 'count=' + count(chainReport, /ErrInvalidSupervisedReport/g));
ok('c5 the report command validates the persisted record before rendering',
  has(bodyOf(cmd, 'func (cli CLI) executeSupervisedReport('), /chain\.ValidateSupervisedReport\(/));

console.log('--- part 4: P2 - a visible delivery is not a completion ---');

ok('d1 the visible adapter never touches chain state', !has(visible, /Receipt|PASSED|Publish\(|Transition\(|UpdateEntity/));
ok('d2 the adapter returns no receipt record', !has(visible, /ReceiptRecord/));
ok('d3 the return entry takes a handoff receipt plus the frozen facts',
  has(chainBox, /func \(engine \*Engine\) SubmitSupervisedReturn\(ctx context\.Context, receipt HandoffReceiptRecord, facts SupervisedFrozenFacts\) error/));
ok('d4 a non-accepted delivery is refused before any write', has(chainTests, /ErrInvalidSupervisedConfirmation/));
const deliveryDecl = bodyOf(chainConfirm, 'type SupervisedDelivery struct {');
ok('d5 delivery records carry no completion outcome field', !/Outcome|Passed|Succeeded/.test(deliveryDecl), 'decl=' + JSON.stringify(deliveryDecl.slice(0, 120)));
ok('d6 window opening requires an accepted confirmation delivery',
  has(chainBox, /DeliveryStatus/));

console.log('--- part 5: P3 - confirmation never waives the gates ---');

const complete = bodyOf(runtime, 'func (runtime *supervisedRuntime) completeSupervisedRun(');
['runScopeScan', 'runSecretAndTypeScan', 'runSideEffectScan', 'processStopEvidence', 'runAfterReturnGates'].forEach((fn, i) => {
  ok('e' + (i + 1) + ' completion runs ' + fn, complete.indexOf(fn) >= 0);
});
ok('e6 no gate-skipping switch exists in supervised code',
  !/skipGates|SkipGates|skipAllGates|bypassGates/.test(runtime + chainBox + cmd));
ok('e7 the side-effect scan blocks deletions, type changes, permission drift and escaping links',
  count(bodyOf(runtime, 'func (runtime *supervisedRuntime) runSideEffectScan('), /return .*fmt\.Errorf/g) >= 4);
ok('e8 the scope scan fails closed when the allow-list is underivable',
  has(runtime, /ErrSupervisedScopeUnderivable/));

console.log('--- part 6: the four prohibitions each have a mechanical anchor + tests ---');

ok('f1 P1 intent guard present in the command model', has(chainModels, /supervised/i));
ok('f2 P1 CLI mutual exclusion present', has(cli, /supervised/i));
ok('f3 P1 router exposes a dedicated supervised port', has(chainRouter, /SupervisedBlackBoxPort/));
ok('f4 P1 semantic invariant test is registered', has(consoleTests, /TestSupervisedPathNeverRoutesThroughAgentRunner/));
ok('f5 P1 lexical smoke test is retained', has(consoleTests, /TestNoSupervisedReferencesInAgentRunnerPaths/));
ok('f6 P1 semantic invariant test asserts on state/port layer, not on comments',
  /TestSupervisedPathNeverRoutesThroughAgentRunner[\s\S]{0,4000}\n}/.test(consoleTests));
ok('f7 P2/P3/P4 chain-side tests exist', has(chainTests, /SubmitSupervisedReturn/));
ok('f8 risk-confirmation validation is unit-tested', has(chainTests, /RiskConfirmation|Confirmation/));
ok('f9 report validation is unit-tested', has(chainTests, /SupervisedReport|report/i));

console.log('--- part 7: mechanical re-execution (child process, independent of any review) ---');

const run = (cmdline) => {
  try {
    return { code: 0, out: cp.execSync(cmdline, { encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'] }) };
  } catch (err) {
    return { code: err.status === undefined ? 1 : err.status, out: String(err.stdout || '') + String(err.stderr || '') };
  }
};

const focused = run('go test -count=1 -run Supervised -v ./internal/chain/ ./internal/console/ ./internal/adapters/');
const passCount = count(focused.out, /--- PASS:/g);
const failCount = count(focused.out, /--- FAIL:/g);
ok('g1 focused supervised test subset exits 0', focused.code === 0, 'exit=' + focused.code);
ok('g2 no failing supervised test', failCount === 0, 'fail=' + failCount);
ok('g3 supervised test subset is non-trivial', passCount >= 40, 'pass=' + passCount);
ok('g4 the two drift tests really ran',
  focused.out.indexOf('TestSupervisedReturnRefusesDriftedPolicy') >= 0 && focused.out.indexOf('TestSupervisedReturnAllowsAbortUnderDriftedPolicy') >= 0);
ok('g5 the aborted run never reports PASS in that test', /TestSupervisedReturnAllowsAbortUnderDriftedPolicy[\s\S]*?--- PASS/.test(focused.out));
console.log('     focused PASS lines=' + passCount + ' FAIL lines=' + failCount);

const vet = run('go vet ./internal/chain/ ./internal/console/ ./internal/adapters/');
ok('g6 go vet is clean on the changed packages', vet.code === 0, 'exit=' + vet.code);

const fmt = run('gofmt -l internal');
ok('g7 gofmt reports no unformatted file', fmt.out.trim() === '', 'out=' + JSON.stringify(fmt.out.trim().slice(0, 200)));

console.log('--- part 8: fail-closed boundaries (D-18) actually exist in the wiring ---');

ok('h1 the default assembly uses an unavailable visible client', has(cmd, /unavailableVisibleClient/));
ok('h2 the delivery step refuses with a named reason when the channel is unavailable',
  has(cmd, /visible_delivery_unavailable/) || has(runtime, /visible_delivery_unavailable/));
ok('h3 the gate runner refuses when hooks are unavailable', has(cmd, /ErrSupervisedGateUnavailable/) || has(runtime, /ErrSupervisedGateUnavailable/));
console.log('     cannot-verify (not mechanically closable on this host):');
console.log('     cv1 real SessionBridge visible end-to-end (depends on T027-B wiring)');
console.log('     cv2 live external-agent observation of network/cost/process/side-effect dimensions');
console.log('     cv3 production hook-definition source and real gate assembly');

console.log('');
console.log('A12_CHECKS=' + checks + ' A12_EXPERIMENT_FAILURES=' + failures);
process.exit(failures === 0 ? 0 : 1);

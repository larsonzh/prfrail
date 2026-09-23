'use strict';

/**
 * tools/gates/selftest.js - in-process assertions over tools/gates/testdata/**.
 *
 * Every assertion is named, so a broken fixture or a weakened criterion turns the
 * suite red. Coverage (architecture design section 6):
 *   T1, T2, T4, T5, T6, T7, T8, T9, T10, T11, T13, T14, T15, T16, T17, T18 + G7-d
 *   + T19-T28 (G6-5 section-reference resolution, G6-6 ledger status closed set).
 * T3 (defect A1 against commit b738e6b) and T12 (defect C1 against commit 1b01115)
 * need real git ranges and are therefore executed as recorded evidence:
 *   node tools/gates/gate.js --check=G6 --scope=range:b738e6b^..b738e6b
 *   node tools/gates/gate.js --check=G6 --scope=range:9f68a52^..9f68a52
 *   node tools/gates/gate.js --check=G7 --scope=range:1b01115^..1b01115
 *
 * Run: node tools/gates/selftest.js   (or: node tools/gates/gate.js --selftest)
 */

const fs = require('fs');
const path = require('path');
const cp = require('child_process');

const {
  discoverRepoRoot,
  decodeBuffer,
  resolveWhitelistSource,
  whitelistSourceLabel,
  resolveRepoArtifact,
  createContext,
  UsageError,
  EMPTY_TREE,
} = require('./lib/ctx');
const { createReport } = require('./lib/report');
const G6 = require('./lib/criteria/g6');
const G7 = require('./lib/criteria/g7');
const G4B = require('./lib/criteria/g4b');

const FIXTURES_DIR = path.join(__dirname, 'testdata', 'fixtures');
const HISTORICAL_DIR = path.join(__dirname, 'testdata', 'historical');

const FIXTURES = {
  'g6-1-pending': { disk: 'g6-1-pending-done.txt', virtual: 'docs/t027/fixture-g6-1-pending.md' },
  'g6-1-ob11': { disk: 'g6-1-ob11-row.txt', virtual: 'docs/t027/fixture-g6-1-ob11.md' },
  'g6-2-cost': { disk: 'g6-2-cost-clauses.txt', virtual: 'docs/t027/fixture-g6-2-cost.md' },
  'g6-3-legal': { disk: 'g6-3-legal-run-id.txt', virtual: 'docs/t027/fixture-g6-3-legal.md' },
  'g6-3-bad': { disk: 'g6-3-placeholder-run-id.txt', virtual: 'docs/t027/fixture-g6-3-placeholder.md' },
  'g6-4-literals': { disk: 'g6-4-expiring-literals.txt', virtual: 'docs/t027/fixture-g6-4-literals.md' },
  'g7-a': { disk: 'g7-a-no-artifacts.txt', virtual: 'docs/validation/fixture-g7-a-no-artifacts.md' },
  'g7-b-bad': { disk: 'g7-b-no-disposition.txt', virtual: 'docs/validation/fixture-g7-b-no-disposition.md' },
  'g7-b-good': { disk: 'g7-ok-report.txt', virtual: 'docs/validation/fixture-g7-ok-report.md' },
  'g7-c': { disk: 'g7-c-finding-row.txt', virtual: 'docs/validation/fixture-g7-c-finding-row.md' },
  'g7-d': { disk: 'g7-d-stopped-state.txt', virtual: 'docs/validation/fixture-g7-d-stopped-state.md' },
  // G6-5 / G6-6 (slice DIRECTIVE-REVIEW-SATURATION, v1.13). The `_EN` fixture MUST end in
  // `_EN.md`: the criterion picks its word table with `/_EN\.md$/i`.
  'g6-5-legal-refs': { disk: 'g6-5-legal-refs.txt', virtual: 'docs/t027/fixture-g6-5-legal-refs.md' },
  'g6-5-dangling-ref': { disk: 'g6-5-dangling-ref.txt', virtual: 'docs/t027/fixture-g6-5-dangling-ref.md' },
  'g6-5-c1-v31-table': { disk: 'g6-5-c1-v31-table.txt', virtual: 'docs/t027/fixture-g6-5-c1-v31-table.md' },
  'g6-5-outside': { disk: 'g6-5-outside.txt', virtual: 'docs/t027/fixture-g6-5-outside.md' },
  'g6-5-excluded-ref': { disk: 'g6-5-excluded-ref.txt', virtual: 'docs/t027/fixture-g6-5-excluded-ref.md' },
  'g6-6-legal-status': { disk: 'g6-6-legal-status.txt', virtual: 'docs/t027/fixture-g6-6-legal-status.md' },
  'g6-6-legacy-status': { disk: 'g6-6-legacy-status.txt', virtual: 'docs/t027/fixture-g6-6-legacy-status.md' },
  'g6-6-en-status': { disk: 'g6-6-en-status.txt', virtual: 'docs/t027/fixture-g6-6-en-status_EN.md' },
};

/**
 * Fixture files are deliberately NOT named `*.md` on disk. The criteria select their
 * targets by extension (G6: every changed `.md`; G7: changed `docs/validation/*.md`),
 * so a `.md` fixture would be judged by a plain `--scope=tree` run as if it were a real
 * ledger/report - the runner would report its own intentionally-violating test data as
 * a defect and the tree run could never be green. The selftest maps each fixture onto a
 * virtual `.md` path, so the criteria still see `.md` content. The invariant below keeps
 * that convention mechanical.
 */
const FIXTURE_DISK_MUST_NOT_BE_MD = true;

const ALL_SUBCHECKS = [...G6.subchecks, ...G7.subchecks].map((s) => s.id);

/**
 * Expected FAIL id set for the whole fixture corpus (regression anchor, T13).
 * G6-5 / G6-6 joined the anchor with their intentionally-violating fixtures
 * (`g6-5-dangling-ref`, `g6-6-legacy-status`): every corpus FAIL id must be one the suite
 * declares as expected, so a new criterion slipping into the corpus is caught here.
 */
const EXPECTED_CORPUS_FAILS = ['G6-1', 'G6-3', 'G6-4', 'G6-5', 'G6-6', 'G7-a', 'G7-b', 'G7-c', 'G7-d'];

const CHECKS = [];
/** Every declared mutation is logged so a stale one (no-op) turns the suite red. */
const MUTATION_LOG = [];
/**
 * Non-blocking context lines. Used when an assertion can NOT be evaluated in the
 * current clone (a shallow clone lacks the historical commit a recorded-evidence
 * assertion needs): the suite prints an explicit INFO line and moves on instead of
 * turning red for an environment it does not control (batch 2g robustness).
 */
const INFOS = [];

function check(name, ok, detail) {
  CHECKS.push({ name, ok: !!ok, detail: detail === undefined ? '' : String(detail) });
}

function info(name, detail) {
  INFOS.push({ name, detail: detail === undefined ? '' : String(detail) });
}

/** `git rev-parse --verify --quiet <rev>^{commit}` - is the object in this clone? */
function revAvailable(repoRoot, rev) {
  const r = cp.spawnSync('git', ['rev-parse', '--verify', '--quiet', rev + '^{commit}'], {
    cwd: repoRoot,
    encoding: 'utf8',
    windowsHide: true,
  });
  return r.status === 0;
}

function fixtureText(key) {
  return decodeBuffer(fs.readFileSync(path.join(FIXTURES_DIR, FIXTURES[key].disk)));
}

function historicalLines(name) {
  return decodeBuffer(fs.readFileSync(path.join(HISTORICAL_DIR, name)))
    .split('\n')
    .filter((l) => l.trim() && !l.startsWith('#'));
}

function applyMutation(text, m) {
  const before = text;
  let out = text;
  if (m.op === 'replace') out = text.split(m.from).join(m.to);
  else if (m.op === 'firstReplace') {
    const i = text.indexOf(m.from);
    out = i < 0 ? text : text.slice(0, i) + m.to + text.slice(i + m.from.length);
  } else if (m.op === 'dropLineContaining') {
    out = text
      .split('\n')
      .filter((l) => !l.includes(m.needle))
      .join('\n');
  } else {
    throw new Error('unknown fixture mutation op: ' + m.op);
  }
  MUTATION_LOG.push({ key: m.key, op: m.op, needle: m.from || m.needle, applied: out !== before });
  return out;
}

/** Synthetic context: fixture files are presented as added `.md` files of the change set. */
function makeCtx(keys, opts) {
  opts = opts || {};
  const repoRoot = discoverRepoRoot(__dirname);
  if (!repoRoot) throw new Error('selftest needs a git work tree (fixture `repo` paths are resolved against it)');
  const texts = new Map();
  for (const key of keys) {
    if (!FIXTURES[key]) throw new Error('unknown fixture key: ' + key);
    let t = fixtureText(key);
    for (const m of opts.mutations || []) if (m.key === key) t = applyMutation(t, m);
    texts.set(FIXTURES[key].virtual, t);
  }
  const changed = [...texts.keys()].map((p) => ({ path: p, status: 'A', untracked: true, deleted: false }));
  const changedPaths = changed.map((c) => c.path);
  const linesOf = (p) => (texts.has(p) ? texts.get(p).split('\n') : null);
  /** Records every ctx.loadWhitelist() call, so T14 can prove both criteria use it. */
  const loaderCalls = [];
  return {
    repoRoot,
    loaderCalls,
    kind: 'tree',
    scope: { kind: 'tree' },
    baseRef: null,
    headRef: null,
    changed,
    changedPaths,
    recordFor: (p) => changed.find((c) => c.path === p) || null,
    changedMarkdown: () => changedPaths.filter((p) => /\.md$/i.test(p)),
    changedPathsWith: (re) => changedPaths.filter((p) => re.test(p)),
    git: () => ({ ok: false, code: 127, out: '', err: '', spawnError: null }),
    readBytes: (p) => (texts.has(p) ? Buffer.from(texts.get(p), 'utf8') : null),
    readText: (p) => (texts.has(p) ? texts.get(p) : null),
    readLines: linesOf,
    readLineArray: linesOf,
    exists: (p) => fs.existsSync(path.join(repoRoot, p.split('/').join(path.sep))),
    repoArtifactResolvable: (p) =>
      opts.repoArtifactResolvable ? opts.repoArtifactResolvable(p) : resolveRepoArtifact(repoRoot, 'tree', null, p),
    trackedAt: () => false,
    showAt: () => null,
    isIgnored: () => false,
    loadWhitelist: (rel) => {
      loaderCalls.push(rel);
      return { entries: opts.whitelist || [], error: null, source: opts.whitelistSource || 'worktree' };
    },
    listAllMarkdown: () => [],
    listAgentFiles: () => [],
    addedHunks: (p) => (linesOf(p) || []).map((text, i) => ({ n: i + 1, text })),
    addedLines: (p) => linesOf(p) || [],
    addedLineNumbers: (p) => new Set((linesOf(p) || []).map((_, i) => i + 1)),
    numstat: (p) => [String((linesOf(p) || []).length), '0'],
  };
}

function runSubset(ctx, enable) {
  const out = createReport();
  const results = {};
  for (const sub of [...G6.subchecks, ...G7.subchecks]) {
    if (enable.indexOf(sub.id) < 0) continue;
    results[sub.id] = sub.run(ctx, out);
  }
  return { records: out.records, results };
}

/**
 * In-memory context for G6-3 boundary probes. G6-3 touches only these four ctx
 * members, so the digit-band table can enumerate 20 lengths plus the named shapes
 * without adding a fixture file per case (F7), and the batch-2e prose probes can
 * be enumerated the same way.
 */
function probeCtx(line) {
  const p = 'docs/t027/probe-g6-3-boundary.md';
  const text = '- ' + line;
  return {
    loadWhitelist: () => ({ entries: [], error: null, source: 'worktree' }),
    changedMarkdown: () => [p],
    readLineArray: (q) => (q === p ? text.split('\n') : null),
    addedLineNumbers: () => new Set([1]),
  };
}

/** G6-3 verdict + first `why` for one synthetic line (the F7 boundary table / 2e prose probes). */
function g63Probe(line) {
  const out = createReport();
  const res = G6.subchecks.find((s) => s.id === 'G6-3').run(probeCtx(line), out);
  return {
    verdict: verdictOf(out.records, 'G6-3'),
    findings: res.findings.length,
    why: res.findings.length ? res.findings[0].detail : '',
  };
}

function verdictOf(records, id) {
  const rs = records.filter((r) => r.id === id);
  if (!rs.length) return 'ABSENT';
  return rs.some((r) => r.verdict === 'FAIL') ? 'FAIL' : 'PASS';
}

function failIds(records) {
  return [...new Set(records.filter((r) => r.verdict === 'FAIL').map((r) => r.id))].sort();
}

function without(ids) {
  return ALL_SUBCHECKS.filter((id) => ids.indexOf(id) < 0);
}

/**
 * In-memory context for the G7-c / G7-d probes (batch 2g, F2 / F3). Both sub-checks
 * touch only these ctx members, so the structural-disposition probes and the bilingual
 * stopped-state probes can be enumerated without a fixture file per case.
 */
function g7ProbeCtx(text) {
  const p = 'docs/validation/probe-g7.md';
  const lines = String(text).split('\n');
  return {
    kind: 'tree',
    scope: { kind: 'tree' },
    changed: [{ path: p, status: 'A', untracked: true, deleted: false }],
    changedPaths: [p],
    readLineArray: (q) => (q === p ? lines : null),
    addedLineNumbers: () => new Set(lines.map((_, i) => i + 1)),
  };
}

/** G7-c / G7-d verdict + finding count + first `why` for one synthetic report body. */
function g7Probe(id, text) {
  const out = createReport();
  const sub = G7.subchecks.find((s) => s.id === id);
  const res = sub.run(g7ProbeCtx(text), out);
  return {
    verdict: verdictOf(out.records, id),
    findings: res.findings.length,
    why: res.findings.length ? res.findings[0].detail : '',
  };
}

/* ---------------------------------------------------------------------------
 * G6-5 (section-reference resolution) and G6-6 (ledger status closed set) probes
 * (slice DIRECTIVE-REVIEW-SATURATION, v1.13). Both sub-checks landed in g6.js, so
 * T19-T28 below are what makes them falsifiable: each guard (target region, bold-label
 * collection, exclusion markers, closed set, `_EN` word table) is removed in turn and
 * the matching assertion must go red.
 * ------------------------------------------------------------------------ */

/** Region head of the G6-5 target region (§0.3 change log). */
const H030 = '### 0.3 变更日志\n\n';
/** Region head of the G6-6 target region (§12.4 observation ledger). */
const H124 = '### 12.4 观察项台账\n\n';
/** Mutation target ②: `labelsFromLines`'s bold line-leading label collector. */
const BOLD_LABEL_COLLECTOR = '    m = LABEL_BOLD_RE.exec(l);\n    if (m) out.add(m[1]);\n';
/** Mutation target ③: the CN closed set's first entry. */
const CN_FIRST_STATUS_WORD = "  cn: ['观察中', ";
const G6_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g6.js');

function detailOf(records, id) {
  const r = records.find((x) => x.id === id);
  return r ? String(r.detail) : '';
}

/**
 * Run a subset against an arbitrary criterion module. The stock `runSubset` is bound to
 * the shipped G6/G7 objects; the mutation probes need a patched copy instead.
 */
function runIn(mod, ctx, enable) {
  const out = createReport();
  const results = {};
  for (const sub of mod.subchecks) {
    if (enable.indexOf(sub.id) < 0) continue;
    results[sub.id] = sub.run(ctx, out);
  }
  return { records: out.records, results };
}

/**
 * Mutation harness (T25 / T28): compile a PATCHED copy of `lib/criteria/g6.js` and hand
 * back its exports.
 *
 * Why `Module._compile` and not a disk rewrite: the sub-checks call the module-local
 * `labelsFromLines` / read the module-local `LEDGER_STATUS_SETS`, so patching the
 * EXPORTED object after `require` would be a no-op - exactly the stale mutation the F1
 * check exists to catch. Rewriting the file and restoring it would work, but a process
 * killed inside the mutation window (Ctrl-C, timeout, OOM) would leave a mutated
 * criterion on disk for the next `--scope=tree` run to judge. Compiling a patched copy
 * in memory removes the guard from the code that actually executes, cannot leak, and is
 * still verified two ways: the patch must APPLY (`applied`, enforced by F1) and the
 * shipped file must be BYTE-IDENTICAL before and after (`identical`).
 */
function loadMutatedG6(from, to, tag) {
  const before = fs.readFileSync(G6_SOURCE);
  const src = decodeBuffer(before);
  const patched = src.split(from).join(to);
  MUTATION_LOG.push({ key: 'g6.js:' + tag, op: 'sourcePatch', needle: from, applied: patched !== src });
  const M = require('module');
  const mod = new M(G6_SOURCE, null);
  mod.filename = G6_SOURCE;
  mod.paths = M._nodeModulePaths(path.dirname(G6_SOURCE));
  mod._compile(patched, G6_SOURCE);
  return { exports: mod.exports, applied: patched !== src, identical: fs.readFileSync(G6_SOURCE).equals(before) };
}

/** G6-5 verdict + finding count + record detail for one synthetic file body. */
function g65Probe(text, virtualPath) {
  const p = virtualPath || 'docs/t027/probe-g6-5.md';
  const lines = String(text).split('\n');
  const ctx = {
    // `readText` is deliberately absent: the criterion's directive-label read then falls
    // back to the real repository file, which is what gives these probes real labels.
    repoRoot: discoverRepoRoot(__dirname),
    changedMarkdown: () => [p],
    readLineArray: (q) => (q === p ? lines : null),
    addedLineNumbers: () => new Set(lines.map((_, i) => i + 1)),
    loadWhitelist: () => ({ entries: [], error: null, source: 'probe' }),
  };
  const out = createReport();
  const res = G6.subchecks.find((s) => s.id === 'G6-5').run(ctx, out);
  return { verdict: verdictOf(out.records, 'G6-5'), findings: res.findings.length, detail: detailOf(out.records, 'G6-5') };
}

/** G6-6 verdict + finding count + record detail for one synthetic file body. */
function g66Probe(virtualPath, text) {
  const lines = String(text).split('\n');
  const ctx = {
    repoRoot: discoverRepoRoot(__dirname),
    changedMarkdown: () => [virtualPath],
    readLineArray: (q) => (q === virtualPath ? lines : null),
    addedLineNumbers: () => new Set(lines.map((_, i) => i + 1)),
    loadWhitelist: () => ({ entries: [], error: null, source: 'probe' }),
  };
  const out = createReport();
  const res = G6.subchecks.find((s) => s.id === 'G6-6').run(ctx, out);
  return { verdict: verdictOf(out.records, 'G6-6'), findings: res.findings.length, detail: detailOf(out.records, 'G6-6') };
}

function runTests() {
  CHECKS.length = 0;
  MUTATION_LOG.length = 0;
  INFOS.length = 0;
  const v = (key) => FIXTURES[key].virtual;

  // ---- T13 precondition: fixture naming convention (see FIXTURE_DISK_MUST_NOT_BE_MD).
  if (FIXTURE_DISK_MUST_NOT_BE_MD) {
    const offenders = Object.keys(FIXTURES).filter((k) => /\.md$/i.test(FIXTURES[k].disk));
    check(
      'T13 fixture files are not *.md, so a tree/ci run never judges the runner own fixtures',
      offenders.length === 0,
      'offenders=' + JSON.stringify(offenders)
    );
  }

  // ---- T1: G6-1 pending marker + same-label completion marker in one block.
  {
    const base = runSubset(makeCtx(['g6-1-pending']), ALL_SUBCHECKS);
    check(
      'T1 G6-1 detects pending(⑥)+completion(⑥ 实得) in one block',
      verdictOf(base.records, 'G6-1') === 'FAIL',
      'verdict=' + verdictOf(base.records, 'G6-1')
    );
    check('T1 G6-1 reports exactly one finding for the fixture', base.results['G6-1'].findings.length === 1,
      'findings=' + base.results['G6-1'].findings.length);
    const noDone = runSubset(
      makeCtx(['g6-1-pending'], {
        mutations: [{ key: 'g6-1-pending', op: 'dropLineContaining', needle: '（**实得**）' }],
      }),
      ALL_SUBCHECKS
    );
    check('T1 reverse mutation (drop completion row) clears G6-1', verdictOf(noDone.records, 'G6-1') === 'PASS',
      'verdict=' + verdictOf(noDone.records, 'G6-1'));
    const disabled = runSubset(makeCtx(['g6-1-pending']), without(['G6-1']));
    check('T1 delete-criterion mutation (G6-1 removed) leaves the corpus clean', failIds(disabled.records).length === 0,
      'fails=' + JSON.stringify(failIds(disabled.records)));

    // ⑥ F1: the completion-marker set is the UNION the directive enumerates, not the
    // bold-only intersection. Every accepted shape must fire on the real fixture, and
    // the plain shapes are exactly the ones that used to be missed.
    for (const shape of ['（实得）', '✅ 已完成', '✅ **已完成**']) {
      const m = runSubset(
        makeCtx(['g6-1-pending'], {
          mutations: [{ key: 'g6-1-pending', op: 'replace', from: '（**实得**）', to: shape }],
        }),
        ALL_SUBCHECKS
      );
      check('T1 F1 completion shape `' + shape + '` is detected',
        verdictOf(m.records, 'G6-1') === 'FAIL' && m.results['G6-1'].findings.length === 1,
        'verdict=' + verdictOf(m.records, 'G6-1') + ' findings=' + m.results['G6-1'].findings.length);
    }
    // Over-widening guard: the same plain bracket shape with a DIFFERENT label inside
    // the window must stay clean - i.e. widening must not degrade the label requirement.
    const otherLabel = runSubset(
      makeCtx(['g6-1-pending'], {
        mutations: [
          { key: 'g6-1-pending', op: 'replace', from: '（**实得**）', to: '（实得）' },
          { key: 'g6-1-pending', op: 'replace', from: '| ⑥ 独立扫描', to: '| ⑤ 预审' },
        ],
      }),
      ALL_SUBCHECKS
    );
    check('T1 F1 plain shape with a different label stays clean (label window intact)',
      verdictOf(otherLabel.records, 'G6-1') === 'PASS',
      'verdict=' + verdictOf(otherLabel.records, 'G6-1') + ' why=' + (otherLabel.results['G6-1'].findings[0] || {}).detail);
  }

  // ---- T2: self-reference guard on the verbatim OB-11 registration row.
  {
    const ob11 = historicalLines('ob11-registration-row.txt');
    check('T2 provenance file carries exactly one excerpt line', ob11.length === 1, 'lines=' + ob11.length);
    check('T2 fixture embeds the historical OB-11 row verbatim', fixtureText('g6-1-ob11').includes(ob11[0]),
      'excerpt length=' + ob11[0].length);
    const base = runSubset(makeCtx(['g6-1-ob11']), ALL_SUBCHECKS);
    check('T2 OB-11 row (with 冲突 marker) is clean', failIds(base.records).length === 0,
      'fails=' + JSON.stringify(failIds(base.records)));
    const mut = runSubset(
      makeCtx(['g6-1-ob11'], { mutations: [{ key: 'g6-1-ob11', op: 'replace', from: '冲突', to: '不一致' }] }),
      ALL_SUBCHECKS
    );
    check('T2 reverse mutation (drop the 冲突 exclusion marker) turns G6-1 red',
      JSON.stringify(failIds(mut.records)) === JSON.stringify(['G6-1']),
      'fails=' + JSON.stringify(failIds(mut.records)));
  }

  // ---- T4: G6-2 budget-clause exemption vs realized-clause comparison.
  {
    const hist = historicalLines('a1-cost-clauses.txt');
    check('T4 provenance file carries the two historical cost rows', hist.length === 2, 'lines=' + hist.length);
    check('T4 fixture embeds both historical cost rows verbatim',
      hist.every((l) => fixtureText('g6-2-cost').includes(l)), 'rows=' + hist.length);
    const gi = G6.internals;
    const va = gi.vecText(gi.costClause(hist[0]));
    const vb = gi.vecText(gi.costClause(hist[1]));
    check('T4 cost-clause parser yields identical realized vectors for the historical pair',
      va === vb && va === 'total=3+⑥×1+⑦×2', 'v1.9=' + va + ' c0e=' + vb);
    const base = runSubset(makeCtx(['g6-2-cost']), ALL_SUBCHECKS);
    check('T4 budget clause is exempt: consistent rows pass G6-2', verdictOf(base.records, 'G6-2') === 'PASS',
      'verdict=' + verdictOf(base.records, 'G6-2'));
    const mut = runSubset(
      makeCtx(['g6-2-cost'], { mutations: [{ key: 'g6-2-cost', op: 'firstReplace', from: '预算', to: '成本' }] }),
      ALL_SUBCHECKS
    );
    check('T4 reverse mutation (预算 -> 成本) turns G6-2 red', verdictOf(mut.records, 'G6-2') === 'FAIL',
      'verdict=' + verdictOf(mut.records, 'G6-2'));
    const disabled = runSubset(makeCtx(['g6-2-cost'], {
      mutations: [{ key: 'g6-2-cost', op: 'firstReplace', from: '预算', to: '成本' }],
    }), without(['G6-2']));
    check('T4 delete-criterion mutation (G6-2 removed) leaves the mutated corpus clean',
      failIds(disabled.records).length === 0, 'fails=' + JSON.stringify(failIds(disabled.records)));
  }

  // ---- T5: legal run id.
  {
    const r = runSubset(makeCtx(['g6-3-legal']), ALL_SUBCHECKS);
    check('T5 legal 11-digit run id passes G6-3', verdictOf(r.records, 'G6-3') === 'PASS',
      'verdict=' + verdictOf(r.records, 'G6-3'));
    check('T5 legal fixture yields no finding', r.results['G6-3'].findings.length === 0,
      'findings=' + r.results['G6-3'].findings.length);
  }

  // ---- T6: placeholder run ids.
  {
    const r = runSubset(makeCtx(['g6-3-bad']), ALL_SUBCHECKS);
    check('T6 placeholder run ids fail G6-3', verdictOf(r.records, 'G6-3') === 'FAIL',
      'verdict=' + verdictOf(r.records, 'G6-3'));
    check('T6 one finding per placeholder line (4)', r.results['G6-3'].findings.length === 4,
      'findings=' + r.results['G6-3'].findings.length);
    const disabled = runSubset(makeCtx(['g6-3-bad']), without(['G6-3']));
    check('T6 delete-criterion mutation (G6-3 removed) leaves the corpus clean', failIds(disabled.records).length === 0,
      'fails=' + JSON.stringify(failIds(disabled.records)));
  }

  // ---- F7: G6-3 enforces the documented 9-13 digit band from BOTH sides.
  // The criterion text says `run` + space + 9-13 digits; before F7 the code only
  // rejected 1-2 digits and >= 14 digits, so 3-8 digit ids passed (fail-open).
  // Batch 2e narrows the DIGIT-RUN rules with `(?!\s*[A-Za-z\u4e00-\u9fff])`: a run
  // immediately followed by a word is a prose count, not an id. The band is therefore
  // asserted in token-final AND backticked form (both are shapes a real id takes), while
  // the prose shapes and the context-free placeholder rules are pinned separately below.
  {
    const inBand = (n) => n >= 9 && n <= 13;
    const forms = [
      ['token-final', (id) => 'run ' + id],
      ['backticked', (id) => '`run ' + id + '`'],
    ];
    for (const [formName, form] of forms) {
      for (let n = 1; n <= 20; n++) {
        const id = '3'.repeat(n);
        const want = inBand(n) ? 'PASS' : 'FAIL';
        const got = g63Probe(form(id));
        check(
          'F7 G6-3 boundary table (' + formName + '): ' + n + '-digit id `run ' + id + '` must ' + want,
          got.verdict === want,
          'verdict=' + got.verdict + ' want=' + want + ' findings=' + got.findings + ' :: ' + got.why
        );
        if (n >= 3 && n <= 8) {
          check(
            'F7 G6-3 floor rule (' + formName + '): the ' + n + '-digit reject names the 9-digit floor',
            /below the 9-digit floor/.test(got.why),
            got.why
          );
        }
      }
    }
    const named = [
      ['run 0', 'FAIL', 'all-zero id, 1 digit'],
      ['run 000000', 'FAIL', 'all-zero id, 6 digits'],
      ['run ' + '0'.repeat(20), 'FAIL', 'all-zero id, 20 digits (fails at any length)'],
      ['`run ' + '0'.repeat(20) + '`', 'FAIL', 'all-zero id, 20 digits, backticked'],
      ['run <RUN_ID>', 'FAIL', 'angle-bracket placeholder'],
      ['run TBD', 'FAIL', 'word placeholder TBD'],
      ['run TODO', 'FAIL', 'word placeholder TODO'],
      ['run N/A', 'FAIL', 'word placeholder N/A'],
      ['run 待定', 'FAIL', 'placeholder wording 待'],
      ['run 占位', 'FAIL', 'placeholder wording 占位'],
      ['run 35646859765', 'PASS', 'real-world 11-digit run id'],
      ['`run 35646859765`', 'PASS', 'real-world 11-digit run id inside a code span'],
      // Token-final / punctuation-terminal shapes: no following word ⇒ still judged.
      ['see run 111', 'FAIL', 'token-final 3-digit id at end of line'],
      ['run 1234。', 'FAIL', '3-8 digit id terminated by CJK punctuation (not a word)'],
      ['run 7！', 'FAIL', '1-2 digit id terminated by punctuation'],
      ['report run 123456789012345！', 'FAIL', '>= 14 digit id terminated by punctuation'],
      // Context-free placeholder rules: checked in any context, prose included.
      ['see run <RUN_ID> here', 'FAIL', 'angle-bracket placeholder mid-prose'],
      ['see run TBD here', 'FAIL', 'word placeholder TBD mid-prose'],
      ['see run TODO here', 'FAIL', 'word placeholder TODO mid-prose'],
      ['see run N/A here', 'FAIL', 'word placeholder N/A mid-prose'],
      ['see run 待定 here', 'FAIL', 'CJK placeholder 待 mid-prose'],
      ['see run 占位 here', 'FAIL', 'CJK placeholder 占位 mid-prose'],
      // A BACKTICKED token is judged (batch 2e request item 2: `run 1234` inside backticks
      // must still fail, which is what makes the 2e lookahead a prose-only narrowing).
      ['`run 1234`', 'FAIL', 'backticked malformed id stays judged'],
      ['The run log shows `run 1234` as the marker.', 'FAIL', 'backticked malformed id inside a prose sentence'],
    ];
    for (const row of named) {
      const got = g63Probe(row[0]);
      check(
        'F7 G6-3 named case `' + row[0] + '` (' + row[2] + ') must ' + row[1],
        got.verdict === row[1],
        'verdict=' + got.verdict + ' want=' + row[1] + ' :: ' + got.why
      );
    }

    // ---- 2e: prose counts (the measured step-4 false positives) plus the pre-existing
    // repository instance `docs/t027/B3B_DESIGN.md:236` must yield NO finding, for every
    // digit-run rule. A regression that drops the lookahead turns these red.
    const prose = [
      ['The retry loop was executed run 1234 times by the harness.', '3-8 digit rule vs ASCII word'],
      ['连续 run 12 次后任务完成。', '1-2 digit rule vs CJK word'],
      ['`inventory run 1 failed`', 'pre-existing docs/t027/B3B_DESIGN.md:236 shape (1-2 digit rule)'],
      ['run 111 次后停止。', '3-8 digit rule vs CJK word'],
      ['run 000 次。', 'all-zero rule vs CJK word'],
      ['run 123456789012345 次。', '>= 14 digit rule vs CJK word'],
      ['the tool ran run 1234 quickly', '3-8 digit rule vs ASCII word (no punctuation)'],
      ['see run 7 files', '1-2 digit rule vs ASCII word'],
    ];
    for (const row of prose) {
      const got = g63Probe(row[0]);
      check(
        '2e G6-3 prose probe `' + row[0] + '` (' + row[1] + ') must PASS (no finding)',
        got.verdict === 'PASS' && got.findings === 0,
        'verdict=' + got.verdict + ' findings=' + got.findings + ' :: ' + got.why
      );
    }
    check(
      '2e G6-3 prose probes are load-bearing: the same digit runs DO fail when token-final',
      g63Probe('run 1234').verdict === 'FAIL' &&
        g63Probe('run 12').verdict === 'FAIL' &&
        g63Probe('run 000').verdict === 'FAIL' &&
        g63Probe('run 123456789012345').verdict === 'FAIL' &&
        g63Probe('`run 1234`').verdict === 'FAIL',
      'prose narrowing is the only reason the prose probes pass'
    );

    // ---- 2f (NF-2): the LEADING word boundary on every G6-3 rule.
    // Before 2f `run\s+\d{1,2}` also matched inside a word, so `rerun 1234` /
    // `The harness will rerun 1234.` produced a spurious FAIL - a false-positive class in a
    // CI-blocking gate (0 instances in the repo today, but the class must be closed, not
    // merely unobserved). Every rule now carries `(?<![A-Za-z\u4e00-\u9fff])`.
    const leadingBoundary = [
      ['rerun 1234', 'ASCII letter immediately before `run`'],
      ['The harness will rerun 1234.', 'word-embedded `run` inside a prose sentence'],
      ['在run 1234', 'CJK ideograph glued to the FRONT of `run`'],
      ['prerun 12', 'ASCII-letter prefix on the 1-2 digit rule'],
      ['rerun <RUN_ID>', 'placeholder rule carries the leading guard too (NF-2 consistency)'],
      ['prerun TBD', 'word-placeholder rule carries the leading guard too'],
    ];
    for (const row of leadingBoundary) {
      const got = g63Probe(row[0]);
      check(
        '2f G6-3 leading-boundary probe `' + row[0] + '` (' + row[1] + ') must PASS (no finding)',
        got.verdict === 'PASS' && got.findings === 0,
        'verdict=' + got.verdict + ' findings=' + got.findings + ' :: ' + got.why
      );
    }

    // The still-judged shapes: a NON-word character (or the line start) precedes `run`, and
    // the digit run is not followed by a word - the guard must not have narrowed these.
    const stillJudged = [
      ['see run 1234', 'token-final 3-8 digit id mid-prose'],
      ['（run 1234）', 'id wrapped in CJK parentheses'],
      ['`run 1234`', 'backticked malformed id'],
      ['run 1234。', 'id terminated by CJK punctuation'],
      ['run 111', '3-digit id at end of line'],
      ['run TBD', 'word placeholder at end of line'],
      ['see run TBD here', 'word placeholder mid-prose'],
    ];
    for (const row of stillJudged) {
      const got = g63Probe(row[0]);
      check(
        '2f G6-3 still-judged probe `' + row[0] + '` (' + row[1] + ') must FAIL',
        got.verdict === 'FAIL' && got.findings === 1,
        'verdict=' + got.verdict + ' findings=' + got.findings + ' :: ' + got.why
      );
    }

    // The word-followed exemptions (2e narrows the digit-run rules AFTER the digits; 2f does
    // not touch that side). The backticked form is the NF-1 case: backticks do NOT judge a
    // digit run that a word still follows, so the criterion text must say "backticked AND
    // not followed by a word", not "backticked".
    const wordFollowed = [
      ['run 1234 times', 'digit run followed by an ASCII word'],
      ['连续 run 12 次', 'digit run followed by a CJK word'],
      ['`run 1234 times`', 'BACKTICKED digit run still followed by a word (NF-1 prose correction)'],
    ];
    for (const row of wordFollowed) {
      const got = g63Probe(row[0]);
      check(
        '2f G6-3 word-followed exemption `' + row[0] + '` (' + row[1] + ') must PASS (no finding)',
        got.verdict === 'PASS' && got.findings === 0,
        'verdict=' + got.verdict + ' findings=' + got.findings + ' :: ' + got.why
      );
    }

    // Glued forms: no separator before the trailing word. The documented consequence is that
    // `run 1234times` / `run 12345678次` are NOT judged - the lookahead sees a word
    // immediately after the digits and reads the token as prose. This is the accepted cost of
    // the 2e prose narrowing (a real id is never glued to the next word), asserted here so a
    // future tightening has to acknowledge it explicitly.
    const glued = [
      ['run 1234times', 'no separator before the trailing ASCII word'],
      ['run 12345678次', 'no separator before the trailing CJK word'],
    ];
    for (const row of glued) {
      const got = g63Probe(row[0]);
      check(
        '2f G6-3 glued-form consequence `' + row[0] + '` (' + row[1] + ') must PASS (documented, not a defect)',
        got.verdict === 'PASS' && got.findings === 0,
        'verdict=' + got.verdict + ' findings=' + got.findings + ' :: ' + got.why
      );
    }
  }

  // ---- T7: expiring literals and the whitelist mechanism.
  {
    const base = runSubset(makeCtx(['g6-4-literals']), ALL_SUBCHECKS);
    check('T7 expiring literals fail G6-4', verdictOf(base.records, 'G6-4') === 'FAIL',
      'verdict=' + verdictOf(base.records, 'G6-4'));
    check('T7 two literal lines produce two findings', base.results['G6-4'].findings.length === 2,
      'findings=' + base.results['G6-4'].findings.length);
    const one = runSubset(
      makeCtx(['g6-4-literals'], {
        whitelist: [{ match: '995/995', file: v('g6-4-literals'), reason: 'T7 mutation: register one of two' }],
      }),
      ALL_SUBCHECKS
    );
    check('T7 whitelisting ONE of two findings leaves G6-4 red with one fatal finding',
      verdictOf(one.records, 'G6-4') === 'FAIL' && one.results['G6-4'].fatal.length === 1,
      'verdict=' + verdictOf(one.records, 'G6-4') + ' fatal=' + one.results['G6-4'].fatal.length);
    const both = runSubset(
      makeCtx(['g6-4-literals'], {
        whitelist: [
          { match: '995/995', file: v('g6-4-literals'), reason: 'T7 mutation A' },
          { match: '共 128 个文件', file: v('g6-4-literals'), reason: 'T7 mutation B' },
        ],
      }),
      ALL_SUBCHECKS
    );
    check('T7 whitelisting both findings turns G6-4 green (whitelist is load-bearing)',
      verdictOf(both.records, 'G6-4') === 'PASS' && both.records.filter((r) => r.verdict === 'SKIP').length === 2,
      'verdict=' + verdictOf(both.records, 'G6-4') + ' skips=' + both.records.filter((r) => r.verdict === 'SKIP').length);
  }

  // ---- T8: missing artifacts block.
  {
    const r = runSubset(makeCtx(['g7-a']), ALL_SUBCHECKS);
    check('T8 report without an artifacts block fails G7-a', verdictOf(r.records, 'G7-a') === 'FAIL',
      'verdict=' + verdictOf(r.records, 'G7-a'));
    const disabled = runSubset(makeCtx(['g7-a']), without(['G7-a']));
    check('T8 delete-criterion mutation (G7-a removed) leaves the corpus clean', failIds(disabled.records).length === 0,
      'fails=' + JSON.stringify(failIds(disabled.records)));
  }

  // ---- T9: artifacts row disposition.
  {
    const r = runSubset(makeCtx(['g7-b-bad']), ALL_SUBCHECKS);
    check('T9 unresolved repo disposition (tmp path, omitted token) fails G7-b',
      verdictOf(r.records, 'G7-b') === 'FAIL', 'verdict=' + verdictOf(r.records, 'G7-b'));
    const local = runSubset(
      makeCtx(['g7-b-bad'], {
        mutations: [{ key: 'g7-b-bad', op: 'replace', from: 'tmp/drfix/A3.txt\n', to: 'tmp/drfix/A3.txt\tlocal-only\n' }],
      }),
      ALL_SUBCHECKS
    );
    check('T9 marking the tmp path local-only turns G7-b green', verdictOf(local.records, 'G7-b') === 'PASS',
      'verdict=' + verdictOf(local.records, 'G7-b'));
    const unknown = runSubset(
      makeCtx(['g7-b-bad'], {
        mutations: [{ key: 'g7-b-bad', op: 'replace', from: 'tmp/drfix/A3.txt\n', to: 'tmp/drfix/A3.txt\tarchived\n' }],
      }),
      ALL_SUBCHECKS
    );
    check('T9 unknown disposition token fails G7-b', verdictOf(unknown.records, 'G7-b') === 'FAIL',
      'verdict=' + verdictOf(unknown.records, 'G7-b'));
  }

  // ---- T10: declared repo artifact must resolve.
  {
    const ok = runSubset(makeCtx(['g7-b-good']), ALL_SUBCHECKS);
    check('T10 existing repo artifact + local-only + missing-historical pass G7-b',
      verdictOf(ok.records, 'G7-b') === 'PASS', 'verdict=' + verdictOf(ok.records, 'G7-b'));
    const broken = runSubset(
      makeCtx(['g7-b-good'], {
        mutations: [
          {
            key: 'g7-b-good',
            op: 'replace',
            from: 'docs/validation/dr-fix-enforcement-proxy.md\trepo',
            to: 'docs/validation/does-not-exist.md\trepo',
          },
        ],
      }),
      ALL_SUBCHECKS
    );
    check('T10 unresolvable repo artifact fails G7-b', verdictOf(broken.records, 'G7-b') === 'FAIL',
      'verdict=' + verdictOf(broken.records, 'G7-b'));
  }

  // ---- T11: finding row disposition cell.
  {
    const r = runSubset(makeCtx(['g7-c']), ALL_SUBCHECKS);
    check('T11 finding row without a disposition cell fails G7-c', verdictOf(r.records, 'G7-c') === 'FAIL',
      'verdict=' + verdictOf(r.records, 'G7-c'));
    const fixed = runSubset(
      makeCtx(['g7-c'], {
        mutations: [
          {
            key: 'g7-c',
            op: 'replace',
            from: '| **High** | 证据链缺失 | |',
            to: '| **High** | 证据链缺失 | 已修 |',
          },
        ],
      }),
      ALL_SUBCHECKS
    );
    check('T11 adding the disposition cell turns G7-c green', verdictOf(fixed.records, 'G7-c') === 'PASS',
      'verdict=' + verdictOf(fixed.records, 'G7-c'));
    const disabled = runSubset(makeCtx(['g7-c']), without(['G7-c']));
    check('T11 delete-criterion mutation (G7-c removed) leaves the corpus clean',
      failIds(disabled.records).length === 0, 'fails=' + JSON.stringify(failIds(disabled.records)));
  }

  // ---- G7-d (architecture design section 3): stopped state must be marked incomplete.
  {
    const r = runSubset(makeCtx(['g7-d']), ALL_SUBCHECKS);
    check('G7-d stopped status line without 已停止——未完成 fails', verdictOf(r.records, 'G7-d') === 'FAIL',
      'verdict=' + verdictOf(r.records, 'G7-d'));
    check('G7-d fires only on the status line, not on prose naming the rule (self-reference guard)',
      r.results['G7-d'].findings.length === 1, 'findings=' + r.results['G7-d'].findings.length);
    const fixed = runSubset(
      makeCtx(['g7-d'], {
        mutations: [
          {
            key: 'g7-d',
            op: 'replace',
            from: '已停止（等用户裁决后决定是否重启）',
            to: '已停止——未完成（等用户裁决后决定是否重启）',
          },
        ],
      }),
      ALL_SUBCHECKS
    );
    check('G7-d adding the incomplete mark turns it green', verdictOf(fixed.records, 'G7-d') === 'PASS',
      'verdict=' + verdictOf(fixed.records, 'G7-d'));
  }

  // ---- T14: whitelist source rule - one loader, both criteria (F1).
  {
    const sTree = resolveWhitelistSource({ kind: 'tree' });
    const sIndex = resolveWhitelistSource({ kind: 'index' });
    const sRange = resolveWhitelistSource({ kind: 'range', base: 'A', head: 'B' });
    const sCi = resolveWhitelistSource({ kind: 'ci' });
    check(
      'T14 tree/index resolve whitelists from the working tree',
      sTree.kind === 'worktree' && sIndex.kind === 'worktree',
      JSON.stringify([sTree, sIndex])
    );
    check(
      'T14 range/ci resolve whitelists from the endpoint revision (range -> head, ci -> HEAD)',
      sRange.kind === 'endpoint' && sRange.ref === 'B' && sCi.kind === 'endpoint' && sCi.ref === 'HEAD',
      JSON.stringify([sRange, sCi])
    );
    const shared = makeCtx(['g6-4-literals']);
    const sink = createReport();
    const before = shared.loaderCalls.length;
    G6.subchecks.find((s) => s.id === 'G6-4').run(shared, sink);
    const g6Calls = shared.loaderCalls.slice(before);
    const mid = shared.loaderCalls.length;
    G4B.run(shared, sink);
    const g4Calls = shared.loaderCalls.slice(mid);
    check(
      'T14 G6 reads its whitelist through ctx.loadWhitelist',
      g6Calls.length === 1 && /g6-whitelist\.txt$/.test(g6Calls[0]),
      JSON.stringify(g6Calls)
    );
    check(
      'T14 G4b reads its whitelist through the same ctx.loadWhitelist',
      g4Calls.length === 1 && /cjk-newwords\.txt$/.test(g4Calls[0]),
      JSON.stringify(g4Calls)
    );
    const rec = sink.records.map((r) => r.id + ' ' + r.verdict + ' ' + r.detail).join(' ;; ');
    check(
      'T14 both criteria print the whitelist source they used',
      /G6-4[^;]*whitelist-source=worktree/.test(rec) && /G4b[^;]*whitelist-source=worktree/.test(rec),
      rec
    );
  }

  // ---- T15: G7-b `repo` disposition is scope-consistent (F2, F2 probe).
  {
    const repoRoot = discoverRepoRoot(__dirname);
    const ignored = resolveRepoArtifact(repoRoot, 'tree', null, 'tmp/probe-artifact.txt');
    check(
      'T15 tree scope rejects a gitignored / never-staged path marked repo (F2 probe)',
      ignored.ok === false && /index/.test(ignored.reason),
      JSON.stringify(ignored)
    );
    const indexed = resolveRepoArtifact(repoRoot, 'tree', null, 'docs/DELIVERY_DIRECTIVE.md');
    check('T15 tree scope accepts a path present in the index and on disk', indexed.ok === true, JSON.stringify(indexed));
    // Recorded-evidence assertion. It must NOT hard-depend on the old commit: a shallow
    // clone (CI) lacks `1b01115`, and "object unavailable" is an environment fact, not a
    // regression. Unavailable ⇒ SKIP with an explicit INFO line (batch 2g robustness).
    if (!revAvailable(repoRoot, '1b01115')) {
      info(
        'T15 range-scope repo-artifact evidence SKIPPED (INFO): revision 1b01115 is unavailable in this clone (shallow clone?)',
        'assertion would otherwise be red for a reason the code does not control'
      );
    } else {
      const atRef = resolveRepoArtifact(repoRoot, 'range', '1b01115', 'docs/validation/dr-fix-enforcement-proxy.md');
      check(
        'T15 range scope still resolves repo artifacts at the endpoint revision',
        atRef.ok === true,
        JSON.stringify(atRef)
      );
    }
    // Any-commit variant of the same F1 rule (batch 2g): `range:HEAD..HEAD` needs no
    // historical object, so the assertion holds in every clone. A whitelist path the
    // endpoint revision does not carry must resolve to `endpoint:<rev>` with ZERO
    // entries - never fall back to the working tree (the F1 failure mode).
    const anyCtx = createContext({ repoRoot, scope: 'range:HEAD..HEAD' });
    const absent = anyCtx.loadWhitelist('tools/gates/absent-f2f3-probe-whitelist.txt');
    check(
      'T15 whitelist source resolves to endpoint:<rev> and yields zero entries when the endpoint revision lacks the file',
      /^endpoint:/.test(absent.source) && absent.entries.length === 0 && absent.text === null,
      JSON.stringify(absent)
    );
    const headLabel = whitelistSourceLabel(resolveWhitelistSource({ kind: 'range', base: 'HEAD', head: 'HEAD' }));
    // The label must name the ACTUAL endpoint revision, not a hardcoded one: the
    // assertion is therefore true for whatever commit this clone happens to be on.
    const headRev = cp
      .spawnSync('git', ['rev-parse', 'HEAD'], { cwd: repoRoot, encoding: 'utf8', windowsHide: true })
      .stdout.trim();
    const realCtx = createContext({ repoRoot, scope: 'range:' + headRev + '..' + headRev });
    const real = realCtx.loadWhitelist('tools/gates/absent-f2f3-probe-whitelist.txt');
    check(
      'T15 any-commit variant: endpoint label + zero entries hold for the real HEAD revision too',
      headLabel === 'endpoint:HEAD' && real.source === 'endpoint:' + headRev && real.entries.length === 0,
      JSON.stringify({ headLabel, realSource: real.source, headRev, entries: real.entries.length })
    );
    const probe = (rows) => {
      const ctx = makeCtx(['g7-b-good'], {
        mutations: [{ key: 'g7-b-good', op: 'replace', from: 'docs/validation/dr-fix-enforcement-proxy.md\trepo', to: rows }],
      });
      const sink = createReport();
      G7.subchecks.find((s) => s.id === 'G7-b').run(ctx, sink);
      return verdictOf(sink.records, 'G7-b');
    };
    check(
      'T15 G7-b FAILs an ignored path marked repo in tree scope',
      probe('tmp/gates-f2-not-staged.txt\trepo') === 'FAIL',
      'verdict=' + probe('tmp/gates-f2-not-staged.txt\trepo')
    );
    check(
      'T15 G7-b accepts the same row once the path is in the index',
      probe('docs/DELIVERY_DIRECTIVE.md\trepo') === 'PASS',
      'verdict=' + probe('docs/DELIVERY_DIRECTIVE.md\trepo')
    );
  }

  // ---- T16: G7-c disposition vocabulary covers both languages (F3, F3 probe).
  {
    const cnTokens = ['已修', '已整改', '已闭合', '已处置', '接受', '拒收', '记录在案', '待办', '观察中', '待议', '已停止'];
    const enTokens = [
      'Fixed',
      'Resolved',
      'Closed',
      'Accepted',
      'Rejected',
      'Deferred',
      'Recorded',
      'To do',
      'Observed',
      'To be discussed',
      'Not a problem',
      'Documented',
    ];
    const tokens = G7.internals.DISPOSITION_TOKENS;
    check(
      'T16 disposition vocabulary covers the CN and EN token sets stated in 6.3',
      cnTokens.every((t) => tokens.indexOf(t) >= 0) && enTokens.every((t) => tokens.indexOf(t) >= 0),
      'tokens=' + tokens.length
    );
    for (const t of enTokens) {
      const r = runSubset(
        makeCtx(['g7-c'], {
          mutations: [
            { key: 'g7-c', op: 'replace', from: '| **High** | 证据链缺失 | |', to: '| **High** | evidence chain missing | ' + t + ' |' },
          ],
        }),
        ALL_SUBCHECKS
      );
      check(
        'T16 G7-c accepts the English disposition `' + t + '`',
        verdictOf(r.records, 'G7-c') === 'PASS',
        'verdict=' + verdictOf(r.records, 'G7-c')
      );
    }
    const none = runSubset(makeCtx(['g7-c']), ['G7-c']);
    check(
      'T16 G7-c still rejects a finding row with no disposition cell',
      verdictOf(none.records, 'G7-c') === 'FAIL',
      'verdict=' + verdictOf(none.records, 'G7-c')
    );
  }

  // ---- F2 (batch 2g, ⑤ pre-review Medium): the G7-c disposition must be the row's
  // LAST CELL, not a substring anywhere on the line. Before the fix
  // `| **High** | 描述里写“已修” | |` PASSED because the token sat inside the description.
  {
    const head = '# report\n\n## 3. findings\n\n| 严重度 | 发现 | 处置 |\n|---|---|---|\n';
    const rows = [
      ['| **High** | 描述里写“已修” | |', 'FAIL', 'disposition token inside the description cell, disposition cell empty'],
      ['| **High** | evidence chain missing | Fixed |', 'PASS', 'English disposition in the last cell'],
      ['| **High** | 描述里写“已修” | 拒收 |', 'PASS', 'CN disposition in the last cell (token also appears in the description)'],
      ['| **Medium** | description says "Resolved" | Rejected |', 'PASS', 'another EN token pair resolves in the last cell'],
      ['| **High** | evidence chain missing', 'FAIL', 'row with no disposition cell at all (no trailing pipe)'],
      ['| **High** | 证据链缺失 |', 'FAIL', 'row with a structurally present but EMPTY disposition cell'],
      ['| **Low** | 描述含“已修” | 已修 |', 'PASS', 'control: token in both cells'],
    ];
    for (const [row, want, why] of rows) {
      const got = g7Probe('G7-c', head + row);
      check(
        'F2 G7-c structural disposition probe `' + row + '` (' + why + ') must ' + want,
        got.verdict === want,
        'verdict=' + got.verdict + ' want=' + want + ' findings=' + got.findings + ' :: ' + got.why
      );
    }
    check(
      'F2 G7-c takes the disposition from the LAST cell only (internals)',
      G7.internals.dispositionCell('| **High** | 描述里写“已修” | |') === '' &&
        G7.internals.dispositionCell('| **High** | desc | 已修 |') === '已修' &&
        G7.internals.dispositionCell('| **High** | 证据链缺失') === '证据链缺失',
      JSON.stringify([
        G7.internals.dispositionCell('| **High** | 描述里写“已修” | |'),
        G7.internals.dispositionCell('| **High** | desc | 已修 |'),
      ])
    );
  }

  // ---- F3 (batch 2g, ⑤ pre-review Medium): the G7-d stopped-state token set must be
  // bilingual, because G7 targets the `_EN` reports. Before the fix
  // `| Status | Stopped (pending user ruling) |` passed while G7-d claimed EN coverage.
  {
    const head = '# report\n\n## 2. status\n\n';
    check(
      'F3 G7-d stopped-state vocabulary is bilingual (internals)',
      ['已停止', '停止态', 'Stopped', 'Halted'].every((t) => G7.internals.STOP_TOKENS.indexOf(t) >= 0) &&
        ['已停止——未完成', 'Stopped — incomplete'].every((m) => G7.internals.STOP_MARKS.indexOf(m) >= 0),
      JSON.stringify({ tokens: G7.internals.STOP_TOKENS, marks: G7.internals.STOP_MARKS })
    );
    const rows = [
      ['| Status | Stopped (pending user ruling) |', 'FAIL', 'EN status line without the incomplete mark'],
      ['| Status | Stopped — incomplete |', 'PASS', 'EN status line carrying the EN incomplete mark'],
      ['| Status | 已停止——未完成（等用户裁决） |', 'PASS', 'CN status line carrying the CN incomplete mark'],
      ['| Status | 已停止（等用户裁决后决定是否重启） |', 'FAIL', 'CN status line without the mark'],
      ['| Status | Halted (awaiting ruling) |', 'FAIL', 'Halted is part of the EN token set'],
      ['| Status | halted (awaiting ruling) |', 'FAIL', 'EN tokens match case-insensitively (lowercase halted)'],
      ['| Status | stopped — incomplete |', 'PASS', 'the EN mark matches case-insensitively too'],
      [
        '| Status | Stopped (pending user ruling) |\n\nStopped — incomplete\n\n## 4. notes\n\n(nothing)\n',
        'PASS',
        'the mark appears later in the SAME block, so the row is closed out',
      ],
    ];
    for (const [row, want, why] of rows) {
      const got = g7Probe('G7-d', head + row);
      check(
        'F3 G7-d bilingual stopped-state probe `' + row.split('\n')[0] + '` (' + why + ') must ' + want,
        got.verdict === want,
        'verdict=' + got.verdict + ' want=' + want + ' findings=' + got.findings + ' :: ' + got.why
      );
    }
    // Self-reference guard (unchanged by F3): a CN prose line naming 停止态 is not a
    // status line and must yield NO finding.
    const prose = g7Probe('G7-d', '# report\n\n## 6. gaps\n\n已知**无机械判据**（停止态报告头标记）= OB-10 所指缺口。\n');
    check(
      'F3 G7-d ignores a CN prose line naming 停止态 (self-reference guard)',
      prose.verdict === 'PASS' && prose.findings === 0,
      'verdict=' + prose.verdict + ' findings=' + prose.findings
    );
  }

  // ---- T13: whole-corpus regression anchor.
  {
    const all = runSubset(makeCtx(Object.keys(FIXTURES)), ALL_SUBCHECKS);
    check(
      'T13 whole fixture corpus yields exactly the expected FAIL id set',
      JSON.stringify(failIds(all.records)) === JSON.stringify(EXPECTED_CORPUS_FAILS),
      'got=' + JSON.stringify(failIds(all.records)) + ' expected=' + JSON.stringify(EXPECTED_CORPUS_FAILS)
    );
    check(
      'T13 whole fixture corpus has no unexpected SKIP/INFO storm',
      all.records.filter((r) => r.verdict === 'SKIP').length === 0,
      'skips=' + all.records.filter((r) => r.verdict === 'SKIP').length
    );
  }

  // ---- T17: an environment error must never degrade into a vacuous (green) change set.
  // ⑦ F1 found this class: with an unresolvable range the change set collapsed to empty,
  // every criterion printed a vacuous pass, and a CI-blocking gate exited 0 while judging
  // nothing. Both endpoints are now validated and a failed `git diff` is an environment
  // error (exit 2). The positive control keeps the guard from passing by rejecting all.
  {
    // Two independent vectors: a syntactically valid but non-existent object id, and a
    // non-existent ref name. Both used to collapse into "empty change set" + vacuous passes.
    for (const spec of [
      ['base', 'range:deadbeefdeadbeefdeadbeefdeadbeefdeadbeef..HEAD'],
      ['base', 'range:refs/heads/no-such-branch-4e11..HEAD'],
      ['head', 'range:HEAD..deadbeefdeadbeefdeadbeefdeadbeefdeadbeef'],
      ['head', 'range:HEAD..refs/heads/no-such-branch-4e11'],
    ]) {
      let err = null;
      try {
        createContext({ repoRoot: discoverRepoRoot(__dirname), scope: spec[1] });
      } catch (e) {
        err = e;
      }
      check('T17 unresolvable range ' + spec[0] + ' is a usage error, not an empty change set',
        !!err && err instanceof UsageError && /not resolvable/.test(err.message),
        'err=' + (err ? err.message : 'NO ERROR (context built - fail-open)'));
    }
    const okCtx = createContext({
      repoRoot: discoverRepoRoot(__dirname),
      scope: 'range:HEAD..HEAD',
    });
    check('T17 positive control: a resolvable range still builds a context',
      okCtx.kind === 'range' && Array.isArray(okCtx.changedPaths),
      'kind=' + okCtx.kind + ' changed=' + (okCtx.changedPaths || []).length);
  }

  // ---- T18: each scope guard must fire INDEPENDENTLY. 7 round-2 pointed out that removing the
  // endpoint validation and the diff guard together cannot show which one caught what, and that
  // neither the status guard nor the shallow-clone refusal had its own vector. The injected git
  // runner (createContext's test seam) makes each one verifiable on its own.
  {
    const okRes = (out) => ({ ok: true, code: 0, out: out || '', err: '', spawnError: null });
    const badRes = (err) => ({ ok: false, code: 1, out: '', err: err || '', spawnError: null });
    const fake = (rules) => (root, args) => {
      const key = args.join(' ');
      for (const r of rules) if (r.re.test(key)) return r.res();
      return okRes('');
    };
    const repoRoot = discoverRepoRoot(__dirname);
    const probe = (rules, scope) => {
      try {
        return { ctx: createContext({ repoRoot, scope, spawnGit: fake(rules) }), err: null };
      } catch (e) {
        return { ctx: null, err: e };
      }
    };
    const HEAD_OK = { re: /^rev-parse --verify HEAD$/, res: () => okRes('abc1234\n') };
    const NO_PARENT = { re: /^rev-parse --verify HEAD\^1$/, res: () => badRes('') };
    const SHALLOW_YES = { re: /^rev-parse --is-shallow-repository$/, res: () => okRes('true\n') };
    const SHALLOW_NO = { re: /^rev-parse --is-shallow-repository$/, res: () => okRes('false\n') };

    const diffFail = probe(
      [{ re: /^diff --name-status -z /, res: () => badRes('fatal: unable to read tree') }],
      'range:HEAD..HEAD'
    );
    check('T18 a failed `git diff` is an environment error even when both endpoints resolve',
      !!diffFail.err && diffFail.err instanceof UsageError && /git diff --name-status failed/.test(diffFail.err.message),
      'err=' + (diffFail.err ? diffFail.err.message : 'NO ERROR (fail-open)'));

    const statusFail = probe(
      [{ re: /^status --porcelain=v1/, res: () => badRes('fatal: index file corrupt') }],
      'tree'
    );
    check('T18 a failed `git status` in tree scope is an environment error, not an empty change set',
      !!statusFail.err && statusFail.err instanceof UsageError && /git status failed/.test(statusFail.err.message),
      'err=' + (statusFail.err ? statusFail.err.message : 'NO ERROR (fail-open)'));

    const indexFail = probe(
      [{ re: /^status --porcelain=v1/, res: () => badRes('fatal: index file corrupt') }],
      'index'
    );
    check('T18 a failed `git status` in index scope is an environment error too',
      !!indexFail.err && indexFail.err instanceof UsageError && /git status failed/.test(indexFail.err.message),
      'err=' + (indexFail.err ? indexFail.err.message : 'NO ERROR (fail-open)'));

    const shallow = probe([HEAD_OK, NO_PARENT, SHALLOW_YES], 'ci');
    check('T18 --scope=ci refuses a SHALLOW clone instead of falling back to the whole tree',
      !!shallow.err && shallow.err instanceof UsageError && /SHALLOW/.test(shallow.err.message),
      'err=' + (shallow.err ? shallow.err.message : 'NO ERROR (whole-tree fallback)'));

    const inconclusive = probe(
      [HEAD_OK, NO_PARENT, { re: /^rev-parse --is-shallow-repository$/, res: () => badRes('unknown option') }],
      'ci'
    );
    check('T18 an inconclusive shallow probe is an environment error, not a licence to fall back',
      !!inconclusive.err && inconclusive.err instanceof UsageError && /root commit or because/.test(inconclusive.err.message),
      'err=' + (inconclusive.err ? inconclusive.err.message : 'NO ERROR (fell back)'));

    const rootCommit = probe([HEAD_OK, NO_PARENT, SHALLOW_NO], 'ci');
    check('T18 positive control: a genuine root commit still falls back to EMPTY_TREE',
      !!rootCommit.ctx && rootCommit.ctx.baseRef === EMPTY_TREE && rootCommit.ctx.headRef === 'HEAD',
      'err=' + (rootCommit.ctx ? 'baseRef=' + rootCommit.ctx.baseRef : rootCommit.err && rootCommit.err.message));
  }

  // ---- T19: G6-5 legal references - exact heading, parent prefix, bold line-leading label.
  {
    const r = runSubset(makeCtx(['g6-5-legal-refs']), ALL_SUBCHECKS);
    const d = detailOf(r.records, 'G6-5');
    check(
      'T19 G6-5 accepts §7.5 (exact heading) / §7.2 (parent prefix of §7.2.1) / §1.7.1 (bold line-leading label)',
      failIds(r.records).length === 0,
      'fails=' + JSON.stringify(failIds(r.records))
    );
    check(
      'T19 the legal fixture is scanned (6 in-region lines, 3 references) with zero findings',
      r.results['G6-5'].findings.length === 0 && /scanned 6, refs 3/.test(d),
      'findings=' + r.results['G6-5'].findings.length + ' :: ' + d
    );
    const lab = G6.internals.directiveLabels(makeCtx(['g6-5-legal-refs']));
    check(
      'T19 directiveLabels reads the REAL directive: bold line-leading §1.7.1 plus heading §7.5 / §7.2.1',
      lab.has('1.7.1') && lab.has('7.5') && lab.has('7.2.1'),
      'has(1.7.1)=' + lab.has('1.7.1') + ' has(7.5)=' + lab.has('7.5') + ' has(7.2.1)=' + lab.has('7.2.1')
    );
    check(
      'T19 resolvesRef: exact and parent-prefix accepted, a dangling number rejected',
      G6.internals.resolvesRef('7.5', lab) &&
        G6.internals.resolvesRef('7.2', lab) &&
        G6.internals.resolvesRef('1.7.1', lab) &&
        !G6.internals.resolvesRef('7.9', lab),
      JSON.stringify(['7.5', '7.2', '1.7.1', '7.9'].map((n) => n + '=' + G6.internals.resolvesRef(n, lab)))
    );
  }

  // ---- T20: G6-5 dangling reference (the OB-32 "the reference dangles" subclass it does cover).
  {
    const r = runSubset(makeCtx(['g6-5-dangling-ref']), ALL_SUBCHECKS);
    const d = detailOf(r.records, 'G6-5');
    check(
      'T20 G6-5 FAILs a §7.9 reference on a line INSIDE the §0.3 region',
      verdictOf(r.records, 'G6-5') === 'FAIL',
      'verdict=' + verdictOf(r.records, 'G6-5')
    );
    check(
      'T20 exactly one finding, and its detail names `§7.9 unresolved`',
      r.results['G6-5'].findings.length === 1 && /§7\.9 unresolved/.test(d),
      'findings=' + r.results['G6-5'].findings.length + ' :: ' + d
    );
    const fixed = runSubset(
      makeCtx(['g6-5-dangling-ref'], {
        mutations: [{ key: 'g6-5-dangling-ref', op: 'replace', from: '§7.9', to: '§7.8' }],
      }),
      ALL_SUBCHECKS
    );
    check(
      'T20 reverse mutation (§7.9 -> §7.8) clears G6-5: the assertion tracks the reference, not the line',
      verdictOf(fixed.records, 'G6-5') === 'PASS',
      'verdict=' + verdictOf(fixed.records, 'G6-5')
    );
  }

  // ---- T21: the appendix C.1 v3.1 comparison table is outside the target region BY CONSTRUCTION.
  {
    const r = runSubset(makeCtx(['g6-5-c1-v31-table']), ALL_SUBCHECKS);
    const d = detailOf(r.records, 'G6-5');
    check(
      'T21 the C.1 v3.1 comparison table produces no finding (zero corpus FAIL)',
      failIds(r.records).length === 0 && r.results['G6-5'].findings.length === 0,
      'fails=' + JSON.stringify(failIds(r.records))
    );
    check('T21 the excluded §3.10 row is genuinely out of region: 0 references counted', /refs 0/.test(d), d);
    const mut = runSubset(
      makeCtx(['g6-5-c1-v31-table'], {
        mutations: [{ key: 'g6-5-c1-v31-table', op: 'replace', from: '### C.1', to: '### C.1X' }],
      }),
      ALL_SUBCHECKS
    );
    check(
      'T21 boundary mutation (rename the `### C.1` heading) pulls the §3.10 row back into the region and turns G6-5 red',
      verdictOf(mut.records, 'G6-5') === 'FAIL' && /§3\.10 unresolved/.test(detailOf(mut.records, 'G6-5')),
      'verdict=' + verdictOf(mut.records, 'G6-5') + ' :: ' + detailOf(mut.records, 'G6-5')
    );
  }

  // ---- T22: the self-reference guard on the G6-5 path, all five exclusion markers.
  {
    const r = runSubset(makeCtx(['g6-5-excluded-ref']), ALL_SUBCHECKS);
    check(
      'T22 G6-5 skips a line carrying the 示例 exclusion marker',
      verdictOf(r.records, 'G6-5') === 'PASS' && r.results['G6-5'].findings.length === 0,
      'verdict=' + verdictOf(r.records, 'G6-5')
    );
    const mut = runSubset(
      makeCtx(['g6-5-excluded-ref'], {
        mutations: [{ key: 'g6-5-excluded-ref', op: 'firstReplace', from: '示例行：', to: '实执行：' }],
      }),
      ALL_SUBCHECKS
    );
    check(
      'T22 reverse mutation (drop the 示例 marker) turns the very same line red on §7.9',
      verdictOf(mut.records, 'G6-5') === 'FAIL' && /§7\.9 unresolved/.test(detailOf(mut.records, 'G6-5')),
      'verdict=' + verdictOf(mut.records, 'G6-5') + ' :: ' + detailOf(mut.records, 'G6-5')
    );
    for (const marker of ['不得残留', '示例', '引文', '冲突', '矛盾']) {
      const withMarker = g65Probe(H030 + '| v1.13 | x | ' + marker + '：见 §7.9 的说明 | y |\n');
      const withoutMarker = g65Probe(H030 + '| v1.13 | x | 见 §7.9 的说明 | y |\n');
      check(
        'T22 exclusion marker `' + marker + '` suppresses the finding, and dropping it restores the red',
        withMarker.verdict === 'PASS' &&
          withMarker.findings === 0 &&
          withoutMarker.verdict === 'FAIL' &&
          withoutMarker.findings === 1,
        'with=' + withMarker.verdict + '/' + withMarker.findings + ' without=' + withoutMarker.verdict + '/' + withoutMarker.findings
      );
    }
  }

  // ---- T23: the whitelist escape hatch (G6-5 reports through the shared `emit`).
  {
    const v = FIXTURES['g6-5-dangling-ref'].virtual;
    const r = runSubset(
      makeCtx(['g6-5-dangling-ref'], {
        whitelist: [
          { match: '§7.9', file: v, reason: 'T23 in-memory registration: external-document section number (probe only)' },
        ],
      }),
      ALL_SUBCHECKS
    );
    check(
      'T23 an in-memory whitelist entry exempts the dangling reference (PASS + exactly one SKIP)',
      verdictOf(r.records, 'G6-5') === 'PASS' &&
        r.records.filter((x) => x.id === 'G6-5' && x.verdict === 'SKIP').length === 1,
      'verdict=' + verdictOf(r.records, 'G6-5') + ' skips=' + r.records.filter((x) => x.id === 'G6-5' && x.verdict === 'SKIP').length
    );
    const none = runSubset(makeCtx(['g6-5-dangling-ref']), ALL_SUBCHECKS);
    check(
      'T23 the exemption is load-bearing: the same fixture without the entry stays red',
      verdictOf(none.records, 'G6-5') === 'FAIL',
      'verdict=' + verdictOf(none.records, 'G6-5')
    );
    const preloaded = decodeBuffer(fs.readFileSync(path.join(__dirname, 'g6-whitelist.txt')))
      .split('\n')
      .filter((l) => l.trim() && l.trim()[0] !== '#');
    check(
      'T23 no entry was preloaded into tools/gates/g6-whitelist.txt (the exemption above is in-memory only)',
      preloaded.length === 0,
      'entries=' + JSON.stringify(preloaded)
    );
  }

  // ---- T24: a §7.9 reference OUTSIDE every target region is not judged at all.
  {
    const r = runSubset(makeCtx(['g6-5-outside']), ALL_SUBCHECKS);
    const d = detailOf(r.records, 'G6-5');
    check(
      'T24 the out-of-region §7.9 line yields no finding: 6 lines scanned, 0 references counted',
      verdictOf(r.records, 'G6-5') === 'PASS' && r.results['G6-5'].findings.length === 0 && /scanned 6, refs 0/.test(d),
      d
    );
    const mut = runSubset(
      makeCtx(['g6-5-outside'], {
        mutations: [{ key: 'g6-5-outside', op: 'dropLineContaining', needle: '## 9. 其他事项' }],
      }),
      ALL_SUBCHECKS
    );
    check(
      'T24 boundary mutation (drop the closing heading) extends the region onto the §7.9 line and turns G6-5 red',
      verdictOf(mut.records, 'G6-5') === 'FAIL' && /§7\.9 unresolved/.test(detailOf(mut.records, 'G6-5')),
      'verdict=' + verdictOf(mut.records, 'G6-5') + ' :: ' + detailOf(mut.records, 'G6-5')
    );
  }

  // ---- T25: mutation ② (drop the bold line-leading label collection) => §1.7.1 goes red.
  {
    const pre = runIn(G6, makeCtx(['g6-5-legal-refs']), ['G6-5']);
    check(
      'T25 precondition: the legal fixture is green BEFORE the mutation',
      verdictOf(pre.records, 'G6-5') === 'PASS',
      'verdict=' + verdictOf(pre.records, 'G6-5')
    );
    const mut = loadMutatedG6(
      BOLD_LABEL_COLLECTOR,
      '    // mutation: bold-label collection removed\n',
      'labelsFromLines'
    );
    check(
      'T25 mutation ② applied: the patch text matched g6.js (a stale patch would weaken the suite silently)',
      mut.applied === true,
      'applied=' + mut.applied
    );
    const r = runIn(mut.exports, makeCtx(['g6-5-legal-refs']), ['G6-5']);
    const d = detailOf(r.records, 'G6-5');
    check(
      'T25 mutation ②: without the bold-label collector the §1.7.1 reference dangles and G6-5 turns red',
      verdictOf(r.records, 'G6-5') === 'FAIL' && r.results['G6-5'].findings.length === 1 && /§1\.7\.1 unresolved/.test(d),
      'verdict=' + verdictOf(r.records, 'G6-5') + ' findings=' + r.results['G6-5'].findings.length + ' :: ' + d
    );
    check(
      'T25 the same mutation leaves the heading-derived labels intact (§7.5 / §7.2 still resolve)',
      !/§7\.5 unresolved/.test(d) && !/§7\.2 unresolved/.test(d),
      d
    );
    check(
      'T25 restore: g6.js is byte-identical after the mutation window (no disk residue)',
      mut.identical === true,
      'identical=' + mut.identical
    );
    const after = runIn(G6, makeCtx(['g6-5-legal-refs']), ['G6-5']);
    check(
      'T25 restore: the cached criterion module is unmutated, the legal fixture is green again',
      verdictOf(after.records, 'G6-5') === 'PASS',
      'verdict=' + verdictOf(after.records, 'G6-5')
    );
  }

  // ---- T26: G6-6 closed-set legal rows plus the region/row-shape boundaries.
  {
    const r = runSubset(makeCtx(['g6-6-legal-status']), ALL_SUBCHECKS);
    check(
      'T26 all six CN closed-set words pass, including the `**已处置**（v1.1）` parenthetical shape',
      verdictOf(r.records, 'G6-6') === 'PASS' &&
        r.results['G6-6'].findings.length === 0 &&
        /scanned 6/.test(detailOf(r.records, 'G6-6')),
      'verdict=' + verdictOf(r.records, 'G6-6') + ' :: ' + detailOf(r.records, 'G6-6')
    );
    const gi = G6.internals;
    check(
      'T26 the CN/EN closed sets are exactly the six documented words',
      JSON.stringify(gi.LEDGER_STATUS_SETS.cn) === JSON.stringify(['观察中', '待办', '待议', '已处置', '已裁决', '已登记']) &&
        JSON.stringify(gi.LEDGER_STATUS_SETS.en) ===
          JSON.stringify(['Observing', 'To do', 'To be discussed', 'Resolved', 'Ruling', 'Registered']),
      JSON.stringify(gi.LEDGER_STATUS_SETS)
    );
    check(
      'T26 firstStatusWord accepts `词`, `词（括注）`, `**词**（括注）` and rejects a legacy value',
      gi.firstStatusWord('观察中', gi.LEDGER_STATUS_SETS.cn) === '观察中' &&
        gi.firstStatusWord('观察中（当期）', gi.LEDGER_STATUS_SETS.cn) === '观察中' &&
        gi.firstStatusWord('**已处置**（v1.1）', gi.LEDGER_STATUS_SETS.cn) === '已处置' &&
        gi.firstStatusWord('**下一片必做**（…）', gi.LEDGER_STATUS_SETS.cn) === null,
      JSON.stringify([
        gi.firstStatusWord('观察中', gi.LEDGER_STATUS_SETS.cn),
        gi.firstStatusWord('观察中（当期）', gi.LEDGER_STATUS_SETS.cn),
        gi.firstStatusWord('**已处置**（v1.1）', gi.LEDGER_STATUS_SETS.cn),
        gi.firstStatusWord('**下一片必做**（…）', gi.LEDGER_STATUS_SETS.cn),
      ])
    );
    check(
      'T26 rowCellsOf splits on unescaped pipes and drops the two outer cells',
      JSON.stringify(gi.rowCellsOf('| OB-1 | a \\| b | 观察中 |')) === JSON.stringify([' OB-1 ', ' a \\| b ', ' 观察中 ']),
      JSON.stringify(gi.rowCellsOf('| OB-1 | a \\| b | 观察中 |'))
    );
    check(
      'T26 ledgerRegionRange = [§12.4 heading, next ^#{2,4} heading)',
      JSON.stringify(gi.ledgerRegionRange(['x', '### 12.4 台账', '| a |', '### 12.5 缺口'])) === JSON.stringify([1, 3]),
      JSON.stringify(gi.ledgerRegionRange(['x', '### 12.4 台账', '| a |', '### 12.5 缺口']))
    );

    const shared = '| OB-1 | d | x | 下一片必做 |\n\n';
    const above = g66Probe('docs/t027/probe-g6-6-above.md', shared + H124 + '| OB-2 | d | x | 观察中 |\n');
    check(
      'T26 a ledger-shaped row ABOVE the §12.4 heading is out of region and not judged',
      above.verdict === 'PASS' && above.findings === 0,
      'verdict=' + above.verdict + ' findings=' + above.findings
    );
    const nonOb = g66Probe('docs/t027/probe-g6-6-nonob.md', H124 + '| X-1 | d | x | 下一片必做 |\n');
    check(
      'T26 a non-OB table row inside §12.4 is not judged (LEDGER_ROW_RE)',
      nonOb.verdict === 'PASS' && nonOb.findings === 0,
      'verdict=' + nonOb.verdict + ' findings=' + nonOb.findings
    );
    const noRegion = g66Probe('docs/t027/probe-g6-6-noregion.md', '| OB-1 | d | x | 下一片必做 |\n');
    check(
      'T26 a ledger row in a file without a §12.4 heading is not judged (no target region)',
      noRegion.verdict === 'PASS' && noRegion.findings === 0,
      'verdict=' + noRegion.verdict + ' findings=' + noRegion.findings
    );
  }

  // ---- T27: G6-6 legacy status words (the OB-38 failure class) + reverse mutation.
  {
    const r = runSubset(makeCtx(['g6-6-legacy-status']), ALL_SUBCHECKS);
    const d = detailOf(r.records, 'G6-6');
    check(
      'T27 G6-6 FAILs `**下一片必做**` and `**已落地**` with exactly two findings',
      verdictOf(r.records, 'G6-6') === 'FAIL' && r.results['G6-6'].findings.length === 2,
      'verdict=' + verdictOf(r.records, 'G6-6') + ' findings=' + r.results['G6-6'].findings.length + ' :: ' + d
    );
    check(
      'T27 both offending cells are named, and the legal control row (已登记) is NOT reported',
      /下一片必做/.test(d) && /已落地/.test(d) && !/:: 已登记/.test(d),
      d
    );
    const fixed = runSubset(
      makeCtx(['g6-6-legacy-status'], {
        mutations: [
          { key: 'g6-6-legacy-status', op: 'replace', from: '**下一片必做**', to: '**待办**' },
          { key: 'g6-6-legacy-status', op: 'replace', from: '**已落地**', to: '**已处置**' },
        ],
      }),
      ALL_SUBCHECKS
    );
    check(
      'T27 reverse mutation (both legacy words -> closed-set words) clears G6-6',
      verdictOf(fixed.records, 'G6-6') === 'PASS' && fixed.results['G6-6'].findings.length === 0,
      'verdict=' + verdictOf(fixed.records, 'G6-6')
    );
  }

  // ---- T28: mutation ③ (drop 观察中 from the CN closed set) + the `_EN` word-table switch.
  {
    const mut = loadMutatedG6(CN_FIRST_STATUS_WORD, '  cn: [', 'LEDGER_STATUS_SETS.cn');
    check('T28 mutation ③ applied: the patch text matched g6.js', mut.applied === true, 'applied=' + mut.applied);
    check(
      'T28 the mutated CN closed set has five words and no 观察中',
      mut.exports.internals.LEDGER_STATUS_SETS.cn.indexOf('观察中') < 0 &&
        mut.exports.internals.LEDGER_STATUS_SETS.cn.length === 5,
      JSON.stringify(mut.exports.internals.LEDGER_STATUS_SETS.cn)
    );
    const r = runIn(mut.exports, makeCtx(['g6-6-legal-status']), ['G6-6']);
    const d = detailOf(r.records, 'G6-6');
    check(
      'T28 mutation ③: the 观察中 row alone turns red; the other five closed-set rows stay green',
      verdictOf(r.records, 'G6-6') === 'FAIL' &&
        r.results['G6-6'].findings.length === 1 &&
        /:: 观察中(?: |$)/.test(d),
      'verdict=' + verdictOf(r.records, 'G6-6') + ' findings=' + r.results['G6-6'].findings.length + ' :: ' + d
    );
    check(
      'T28 restore: g6.js is byte-identical and the cached module still carries 观察中',
      mut.identical === true && G6.internals.LEDGER_STATUS_SETS.cn.indexOf('观察中') >= 0,
      'identical=' + mut.identical + ' cached=' + JSON.stringify(G6.internals.LEDGER_STATUS_SETS.cn)
    );

    const en = runSubset(makeCtx(['g6-6-en-status']), ALL_SUBCHECKS);
    check(
      'T28 the `_EN` fixture selects the EN word table: all six EN words pass',
      verdictOf(en.records, 'G6-6') === 'PASS' &&
        en.results['G6-6'].findings.length === 0 &&
        /scanned 6/.test(detailOf(en.records, 'G6-6')),
      'verdict=' + verdictOf(en.records, 'G6-6') + ' :: ' + detailOf(en.records, 'G6-6')
    );
    const row = '| OB-1 | d | x | ';
    const cnName = g66Probe('docs/t027/probe-g6-6-cn.md', H124 + row + '**Observing** |\n');
    const enName = g66Probe('docs/t027/probe-g6-6_EN.md', H124 + row + '**Observing** |\n');
    const cnWordInEn = g66Probe('docs/t027/probe-g6-6-mix_EN.md', H124 + row + '观察中 |\n');
    check(
      'T28 the `_EN` word-table switch is load-bearing: EN word fails under a CN name, passes under `_EN`, CN word fails under `_EN`',
      cnName.verdict === 'FAIL' && enName.verdict === 'PASS' && cnWordInEn.verdict === 'FAIL',
      JSON.stringify({ cnName: cnName.verdict, enName: enName.verdict, cnWordInEn: cnWordInEn.verdict })
    );
  }

  // ---- F1: a stale mutation (no-op) must not silently weaken the suite.
  {
    const stale = MUTATION_LOG.filter((m) => !m.applied);
    check(
      'F1 every declared fixture mutation applied (a stale mutation would weaken the suite silently)',
      stale.length === 0,
      'stale=' + JSON.stringify(stale.map((m) => m.key + ' ' + m.op + ' ' + m.needle))
    );
    check('F1 mutation coverage floor (>= 8 fixture mutations exercised)', MUTATION_LOG.length >= 8, 'mutations=' + MUTATION_LOG.length);
  }

  const failed = CHECKS.filter((c) => !c.ok);  const lines = [
    'GATES SELFTEST (in-process; fixtures under tools/gates/testdata/)',
    '  T3/T12 need real git ranges and are recorded as evidence, not here:',
    '    --scope=range:b738e6b^..b738e6b (A1 must be red) / range:9f68a52^..9f68a52 (green)',
    '    --scope=range:1b01115^..1b01115 (C1: G7 must be red)',
  ];
  for (const c of CHECKS) lines.push((c.ok ? 'ok   ' : 'FAIL ') + c.name + (c.ok || !c.detail ? '' : ' :: ' + c.detail));
  for (const i of INFOS) lines.push('INFO ' + i.name + (i.detail ? ' :: ' + i.detail : ''));
  lines.push('SELFTEST: ' + (failed.length ? 'FAIL ' : 'PASS ') + (CHECKS.length - failed.length) + '/' + CHECKS.length);
  return { ok: failed.length === 0, stdout: lines.join('\n'), passed: CHECKS.length - failed.length, total: CHECKS.length };
}

function run() {
  try {
    return runTests();
  } catch (e) {
    return { ok: false, stdout: 'SELFTEST: FAIL 0/0 :: ' + (e && e.stack ? e.stack : e), passed: 0, total: 0 };
  }
}

if (require.main === module) {
  const res = run();
  process.stdout.write(res.stdout + '\n');
  process.exitCode = res.ok ? 0 : 1;
}

module.exports = { run, FIXTURES, EXPECTED_CORPUS_FAILS };

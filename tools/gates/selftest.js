'use strict';

/**
 * tools/gates/selftest.js - in-process assertions over tools/gates/testdata/**.
 *
 * Every assertion is named, so a broken fixture or a weakened criterion turns the
 * suite red. Coverage (architecture design section 6):
 *   T1, T2, T4, T5, T6, T7, T8, T9, T10, T11, T13, T14, T15, T16, T17, T18 + G7-d
 *   + T19-T28 (G6-5 section-reference resolution, G6-6 ledger status closed set)
 *   + T29-T42, T61-T67 (G5-b closed set, plus the G5-a 2 empty-value division of labour).
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
const crypto = require('crypto');

const {
  discoverRepoRoot,
  decodeBuffer,
  isEvidenceTree,
  resolveWhitelistSource,
  whitelistSourceLabel,
  resolveRepoArtifact,
  createContext,
  UsageError,
  EMPTY_TREE,
} = require('./lib/ctx');
const { createReport } = require('./lib/report');
const G1A = require('./lib/criteria/g1a');
const G2 = require('./lib/criteria/g2');
const G4A = require('./lib/criteria/g4a');
const G6 = require('./lib/criteria/g6');
const G7 = require('./lib/criteria/g7');
const G8 = require('./lib/criteria/g8');
const G4B = require('./lib/criteria/g4b');
const G5A = require('./lib/criteria/g5a');
const G5B = require('./lib/criteria/g5b');

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
  // G5-a / G5-b (slice DIRECTIVE-MODEL-LISTS, v1.18). G5-b moved from a blacklist regex
  // to the section 2.2 whitelist closed set, and these fixtures are its first coverage
  // in this suite; `g5b-empty-value` additionally pins the G5-a 2 division of labour.
  'g5b-closed-set-legal': { disk: 'g5b-closed-set-legal.txt', virtual: '.github/agents/fixture-g5b-closed-set-legal.agent.md' },
  'g5b-legal-nokey': { disk: 'g5b-legal-nokey.txt', virtual: '.github/agents/fixture-g5b-legal-nokey.agent.md' },
  'g5b-illegal-kimi': { disk: 'g5b-illegal-kimi.txt', virtual: '.github/agents/fixture-g5b-illegal-kimi.agent.md' },
  'g5b-illegal-case': { disk: 'g5b-illegal-case.txt', virtual: '.github/agents/fixture-g5b-illegal-case.agent.md' },
  'g5b-illegal-quoted': { disk: 'g5b-illegal-quoted.txt', virtual: '.github/agents/fixture-g5b-illegal-quoted.agent.md' },
  'g5b-illegal-multi': { disk: 'g5b-illegal-multi.txt', virtual: '.github/agents/fixture-g5b-illegal-multi.agent.md' },
  'g5b-illegal-legacy': { disk: 'g5b-illegal-legacy.txt', virtual: '.github/agents/fixture-g5b-illegal-legacy.agent.md' },
  'g5b-empty-value': { disk: 'g5b-empty-value.txt', virtual: '.github/agents/fixture-g5b-empty-value.agent.md' },
  'g5b-body-only': { disk: 'g5b-body-only.txt', virtual: '.github/agents/fixture-g5b-body-only.agent.md' },
  'g5b-legal-trailing-space': { disk: 'g5b-legal-trailing-space.txt', virtual: '.github/agents/fixture-g5b-legal-trailing-space.agent.md' },
  'g5b-illegal-inline-comment': { disk: 'g5b-illegal-inline-comment.txt', virtual: '.github/agents/fixture-g5b-illegal-inline-comment.agent.md' },
  'g5b-closed-set-legal-sol': { disk: 'g5b-closed-set-legal-sol.txt', virtual: '.github/agents/fixture-g5b-closed-set-legal-sol.agent.md' },
  'g5b-closed-set-legal-luna': { disk: 'g5b-closed-set-legal-luna.txt', virtual: '.github/agents/fixture-g5b-closed-set-legal-luna.agent.md' },
  // GLM members (slice DIRECTIVE-GLM-LANDING, v1.26). Two positive fixtures pin the two new
  // closed-set members; two negative fixtures pin case sensitivity and the R2.2 boundary (a
  // qualified call string is not the `modelId`), matching the existing g5b-* shape.
  'g5b-closed-set-legal-glm': { disk: 'g5b-closed-set-legal-glm.txt', virtual: '.github/agents/fixture-g5b-closed-set-legal-glm.agent.md' },
  'g5b-closed-set-legal-glm-flash': { disk: 'g5b-closed-set-legal-glm-flash.txt', virtual: '.github/agents/fixture-g5b-closed-set-legal-glm-flash.agent.md' },
  'g5b-illegal-glm-case': { disk: 'g5b-illegal-glm-case.txt', virtual: '.github/agents/fixture-g5b-illegal-glm-case.agent.md' },
  'g5b-illegal-glm-qualified': { disk: 'g5b-illegal-glm-qualified.txt', virtual: '.github/agents/fixture-g5b-illegal-glm-qualified.agent.md' },
  // G8 frozen-evidence-pack fixtures (slice DIRECTIVE-EVIDENCE-SCOPE, v1.24). Each is a
  // JSON descriptor `{name, pack, files, sumsFrom?, expect}`: `files` maps a
  // repository-relative virtual path to content, `pack` is the virtual pack root and the
  // base for the pack-relative `sumsFrom` list and `{{sha256:<rel>}}` values. Positive
  // fixtures compute their SHA256SUMS hashes at load time, so a hardcoded digest cannot rot.
  'g8-ok-pack': { disk: 'g8-ok-pack.json', g8: 'g8-ok-pack.json' },
  'g8-a-member': { disk: 'g8-a-member.json', g8: 'g8-a-member.json' },
  'g8-a-freedom-list': { disk: 'g8-a-freedom-list.json', g8: 'g8-a-freedom-list.json' },
  'g8-a-registered': { disk: 'g8-a-registered.json', g8: 'g8-a-registered.json' },
  'g8-a-stray-root': { disk: 'g8-a-stray-root.json', g8: 'g8-a-stray-root.json' },
  'g8-b-missing-list': { disk: 'g8-b-missing-list.json', g8: 'g8-b-missing-list.json' },
  'g8-b-unlisted-payload': { disk: 'g8-b-unlisted-payload.json', g8: 'g8-b-unlisted-payload.json' },
  'g8-b-hash-mismatch': { disk: 'g8-b-hash-mismatch.json', g8: 'g8-b-hash-mismatch.json' },
  'g8-b-untracked-conflict': { disk: 'g8-b-untracked-conflict.json', g8: 'g8-b-untracked-conflict.json' },
  'g8-b-backslash-sums': { disk: 'g8-b-backslash-sums.json', g8: 'g8-b-backslash-sums.json' },
  'g8-c-ok': { disk: 'g8-c-ok.json', g8: 'g8-c-ok.json' },
  'g8-c-bad-hash': { disk: 'g8-c-bad-hash.json', g8: 'g8-c-bad-hash.json' },
  'g8-c-no-section': { disk: 'g8-c-no-section.json', g8: 'g8-c-no-section.json' },
  'g8-d-ok': { disk: 'g8-d-ok.json', g8: 'g8-d-ok.json' },
  'g8-d-missing-basis': { disk: 'g8-d-missing-basis.json', g8: 'g8-d-missing-basis.json' },
  'exclusion-positive': { disk: 'exclusion-positive.json', g8: 'exclusion-positive.json' },
  'exclusion-negative': { disk: 'exclusion-negative.json', g8: 'exclusion-negative.json' },
  // T50-T58 fixtures (slice DIRECTIVE-EVIDENCE-SCOPE, tester hardening). Same JSON
  // descriptor shape as the G8 fixtures above.
  'g8-t50-in-tree': { disk: 'g8-t50-in-tree.json', g8: 'g8-t50-in-tree.json' },
  'g8-t50-outside': { disk: 'g8-t50-outside.json', g8: 'g8-t50-outside.json' },
  'g8-bom-in-tree': { disk: 'g8-bom-in-tree.json', g8: 'g8-bom-in-tree.json' },
  'g8-bom-outside': { disk: 'g8-bom-outside.json', g8: 'g8-bom-outside.json' },
  'g8-b-malformed-sums': { disk: 'g8-b-malformed-sums.json', g8: 'g8-b-malformed-sums.json' },
  'g8-b-manifest-truncated': { disk: 'g8-b-manifest-truncated.json', g8: 'g8-b-manifest-truncated.json' },
  'g8-c-arch-not-in-payload': { disk: 'g8-c-arch-not-in-payload.json', g8: 'g8-c-arch-not-in-payload.json' },
  'g8-c-arch-not-in-sums': { disk: 'g8-c-arch-not-in-sums.json', g8: 'g8-c-arch-not-in-sums.json' },
  'g8-d-body-only': { disk: 'g8-d-body-only.json', g8: 'g8-d-body-only.json' },
  'g8-d-bad-sha': { disk: 'g8-d-bad-sha.json', g8: 'g8-d-bad-sha.json' },
  'g8-a-real-registered': { disk: 'g8-a-real-registered.json', g8: 'g8-a-real-registered.json' },
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

const ALL_SUBCHECKS = [...G6.subchecks, ...G7.subchecks, ...G8.subchecks].map((s) => s.id);

/**
 * Expected FAIL id set for the whole fixture corpus (regression anchor, T13).
 * G6-5 / G6-6 joined the anchor with their intentionally-violating fixtures
 * (`g6-5-dangling-ref`, `g6-6-legacy-status`): every corpus FAIL id must be one the suite
 * declares as expected, so a new criterion slipping into the corpus is caught here.
 *
 * The anchor vocabulary is `ALL_SUBCHECKS` (G6 / G7 / G8 sub-checks only; G1-a / G2 / G4a /
 * G4b are single-module criteria asserted directly in T47). DIRECTIVE-EVIDENCE-SCOPE (v1.24)
 * added the four G8 ids, from the intentionally broken G8 fixtures (stray root / missing
 * list / unlisted payload / hash mismatch / untracked collision / missing basis) plus the
 * tree-level `exclusion-positive` probe path (not a pack member / not a freedom-list).
 */
const EXPECTED_CORPUS_FAILS = [
  'G6-1',
  'G6-3',
  'G6-4',
  'G6-5',
  'G6-6',
  'G7-a',
  'G7-b',
  'G7-c',
  'G7-d',
  'G8-a',
  'G8-b',
  'G8-c',
  'G8-d',
];

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

/** Lowercase sha256 hex of a Buffer (the G8 fixture loader computes real digests). */
function sha256Hex(buf) {
  return crypto.createHash('sha256').update(buf).digest('hex');
}

/**
 * Load a G8 JSON fixture into `{ pack, texts }`.
 *
 *   files     repository-relative virtual path -> content
 *   pack      repository-relative pack root (also the base for pack-relative values)
 *   sumsFrom  pack-relative payload paths; their REAL sha256 is computed here and written
 *             into `<pack>/SHA256SUMS.txt`, so a positive fixture never hardcodes a digest.
 *   {{sha256:<rel>}} in any content is replaced by the real hash of `<pack>/<rel>`.
 */
function loadG8Fixture(key) {
  const spec = JSON.parse(decodeBuffer(fs.readFileSync(path.join(FIXTURES_DIR, FIXTURES[key].g8))));
  const pack = spec.pack;
  const texts = new Map();
  for (const repoRel of Object.keys(spec.files || {})) texts.set(repoRel, spec.files[repoRel]);
  const digestOf = (repoRel) => sha256Hex(Buffer.from(texts.has(repoRel) ? texts.get(repoRel) : '', 'utf8'));
  if (Array.isArray(spec.sumsFrom)) {
    const lines = spec.sumsFrom
      .slice()
      .sort()
      .map((rel) => sha256Hex(Buffer.from(texts.get(pack + '/' + rel) || '', 'utf8')) + '  ' + rel);
    texts.set(pack + '/SHA256SUMS.txt', lines.join('\n') + '\n');
  }
  for (const [repoRel, text] of [...texts.entries()]) {
    if (!text.includes('{{sha256:')) continue;
    texts.set(repoRel, text.replace(/\{\{sha256:([^}]+)\}\}/g, (_, rel) => digestOf(pack + '/' + rel)));
  }
  return { pack, texts, spec };
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
  } else if (m.op === 'dropLineContaining' || m.op === 'dropRowContaining') {
    out = text
      .split('\n')
      .filter((l) => !l.includes(m.needle))
      .join('\n');
  } else if (m.op === 'flipHashDigit') {
    // Change the FIRST hex digit of the text (a SHA256SUMS body starts with one), so a
    // previously correct digest becomes wrong in exactly one nibble.
    out = text.replace(/[0-9a-f]/, (c) => (c === '0' ? '1' : '0'));
  } else {
    throw new Error('unknown fixture mutation op: ' + m.op);
  }
  MUTATION_LOG.push({ key: m.key, op: m.op, needle: m.from || m.needle || (m.op === 'flipHashDigit' ? 'first-hex-digit' : undefined), applied: out !== before });
  return out;
}

/**
 * Synthetic context: fixture files are presented as added files of the change set.
 *
 * A `key` may map to ONE virtual path (`virtual` string), to SEVERAL sharing the disk
 * content (`virtual` array), or to a path->content map (`virtual` object). G8 keys carry
 * a JSON descriptor instead (`FIXTURES[key].g8`), expanded by `loadG8Fixture`.
 */
function makeCtx(keys, opts) {
  opts = opts || {};
  const repoRoot = discoverRepoRoot(__dirname);
  if (!repoRoot) throw new Error('selftest needs a git work tree (fixture `repo` paths are resolved against it)');
  const texts = new Map();
  for (const key of keys) {
    if (!FIXTURES[key]) throw new Error('unknown fixture key: ' + key);
    if (FIXTURES[key].g8) {
      const g = loadG8Fixture(key);
      for (const [p, t] of g.texts) texts.set(p, t);
      for (const m of opts.mutations || []) {
        if (m.key !== key) continue;
        if (!m.path || !texts.has(m.path)) {
          MUTATION_LOG.push({ key: key, op: m.op, needle: m.path, applied: false });
          continue;
        }
        texts.set(m.path, applyMutation(texts.get(m.path), m));
      }
      continue;
    }
    let t = fixtureText(key);
    for (const m of opts.mutations || []) if (m.key === key) t = applyMutation(t, m);
    const v = FIXTURES[key].virtual;
    if (Array.isArray(v)) for (const p of v) texts.set(p, t);
    else if (v && typeof v === 'object') for (const p of Object.keys(v)) texts.set(p, v[p] === true ? t : v[p]);
    else texts.set(v, t);
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
    // Single source: reuse lib/ctx's `isEvidenceTree` instead of a local regex, so
    // "the criterion was fixed but the self-test kept the old rule" cannot drift (6.3 OB-79).
    changedMarkdown: () => changedPaths.filter((p) => /\.md$/i.test(p) && !isEvidenceTree(p)),
    changedPathsWith: (re) => changedPaths.filter((p) => re.test(p)),
    isEvidenceTree: (p) => isEvidenceTree(p),
    git: () => ({ ok: false, code: 127, out: '', err: '', spawnError: null }),
    readBytes: (p) => (texts.has(p) ? Buffer.from(texts.get(p), 'utf8') : null),
    readText: (p) => (texts.has(p) ? texts.get(p) : null),
    readLines: linesOf,
    readLineArray: linesOf,
    exists: (p) => texts.has(p) || fs.existsSync(path.join(repoRoot, p.split('/').join(path.sep))),
    repoArtifactResolvable: (p) =>
      opts.repoArtifactResolvable ? opts.repoArtifactResolvable(p) : resolveRepoArtifact(repoRoot, 'tree', null, p),
    trackedAt: () => false,
    showAt: () => null,
    isIgnored: () => false,
    // G8 pack enumeration reads the FIXTURE virtual map only - never the real disk, which
    // would pull the entire real corpus into a fixture probe.
    listTreeFiles: (dir) => {
      const d = String(dir).replace(/\/+$/, '');
      return [...texts.keys()].filter((p) => p.startsWith(d + '/')).sort();
    },
    loadWhitelist: (rel) => {
      loaderCalls.push(rel);
      // The generic `opts.whitelist` knob feeds G6 / G4b (any criterion rel). The G8
      // registration table is deliberately opt-in per rel (`opts.whitelists[FORMS_REL]`)
      // or raw (`opts.whitelistText[FORMS_REL]`), so a G6 probe's entries can never leak
      // into the forms table and turn every corpus run into a duplicate-registration error.
      const isForms = rel === G8.internals.FORMS_REL;
      const hasRel = opts.whitelists && Object.prototype.hasOwnProperty.call(opts.whitelists, rel);
      const entries = hasRel ? opts.whitelists[rel] : isForms ? [] : opts.whitelist || [];
      const hasText = opts.whitelistText && Object.prototype.hasOwnProperty.call(opts.whitelistText, rel);
      let text = hasText ? opts.whitelistText[rel] : null;
      if (text === null && (hasRel || (!isForms && opts.whitelist))) {
        text = entries.map((e) => e.match + '\t' + e.file + '\t' + e.reason).join('\n');
      }
      return { entries, error: null, source: opts.whitelistSource || 'worktree', text };
    },
    listAllMarkdown: () => [],
    listAgentFiles: () => [],
    addedHunks: (p) => (linesOf(p) || []).map((text, i) => ({ n: i + 1, text })),
    addedLines: (p) => linesOf(p) || [],
    addedLineNumbers: (p) => new Set((linesOf(p) || []).map((_, i) => i + 1)),
    numstat: (p) => [String((linesOf(p) || []).length), '0'],
  };
}

/** Run a set of criterion modules against one synthetic context. */
function runMods(ctx, mods) {
  const out = createReport();
  for (const m of mods) m.run(ctx, out);
  return { records: out.records };
}

/**
 * G8 sub-check probe: run ONE G8 sub-check over the fixture keys and return its verdict,
 * finding count and the joined finding details (so a FAIL can be tied to its reason token).
 */
function g8Probe(subId, keys, opts) {
  const ctx = makeCtx(keys, opts);
  const out = createReport();
  const sub = G8.subchecks.find((s) => s.id === subId);
  const res = sub.run(ctx, out) || { findings: [] };
  return {
    verdict: verdictOf(out.records, subId),
    info: out.records.some((r) => r.id === subId && r.verdict === 'INFO'),
    findings: res.findings.length,
    details: res.findings.map((f) => f.detail).join(' | '),
    detail: detailOf(out.records, subId),
  };
}

/* ---------------------------------------------------------------------------
 * T50-T58 helpers (slice DIRECTIVE-EVIDENCE-SCOPE, ③ tester hardening).
 *
 * `injectedCtx` builds a REAL lib/ctx context (so `changedMarkdown()` and
 * `isEvidenceTree()` come from the SHIPPED implementation, not a selftest-local
 * re-derivation) whose change set is supplied by an injected git runner and whose
 * file bodies come from a synthetic map. That is what makes the 6.3 exclusion
 * assertions mutation-sensitive: removing the exclusion inside lib/ctx.js turns
 * T50/T51 red, whereas a `makeCtx`-only assertion would stay green because the
 * selftest re-implements the filter locally.
 * ------------------------------------------------------------------------ */
function g8Texts(key) {
  return loadG8Fixture(key).texts;
}

function injectedCtx(paths, texts, opts) {
  opts = opts || {};
  const repoRoot = discoverRepoRoot(__dirname);
  const spawnGit = (root, args) => {
    const key = args.join(' ');
    if (/^status --porcelain=v1/.test(key)) {
      return { ok: true, code: 0, out: paths.map((p) => '?? ' + p + '\0').join(''), err: '', spawnError: null };
    }
    if (/^check-ignore/.test(key)) return { ok: false, code: 1, out: '', err: '', spawnError: null };
    return { ok: true, code: 0, out: '', err: '', spawnError: null };
  };
  const ctx = createContext({ repoRoot, scope: opts.scope || 'tree', spawnGit });
  const linesOf = (p) => (texts.has(p) ? String(texts.get(p)).split('\n') : null);
  // NOTE: `changedMarkdown` / `isEvidenceTree` are deliberately NOT overridden.
  return Object.assign(ctx, {
    readBytes: (p) => (texts.has(p) ? Buffer.from(texts.get(p), 'utf8') : null),
    readText: (p) => (texts.has(p) ? texts.get(p) : null),
    readLines: linesOf,
    readLineArray: linesOf,
    exists: (p) => texts.has(p),
    listAllMarkdown: () => [],
    loadWhitelist: () => ({ entries: [], error: null, source: 'probe' }),
    addedHunks: (p) => (linesOf(p) || []).map((text, i) => ({ n: i + 1, text })),
    addedLines: (p) => linesOf(p) || [],
    addedLineNumbers: (p) => new Set((linesOf(p) || []).map((_, i) => i + 1)),
    numstat: (p) => [String((linesOf(p) || []).length), '0'],
  });
}

/** Run criterion modules against a REAL-ctx hybrid built from one g8 fixture's texts. */
function runKeyMods(key, mods) {
  const texts = g8Texts(key);
  return runMods(injectedCtx([...texts.keys()], texts), mods);
}

/** Run criterion modules against a REAL-ctx hybrid built from several g8 fixtures. */
function runKeysMods(keys, mods) {
  const texts = new Map();
  for (const key of keys) for (const [p, t] of g8Texts(key)) texts.set(p, t);
  return runMods(injectedCtx([...texts.keys()], texts), mods);
}

/**
 * G8-d probe over a REAL repository file (the precedent file under
 * `docs/validation/evidence/`), so the assertion runs against the shipped bytes
 * rather than a fixture paraphrase.
 */
function g8dRealProbe(relPath) {
  const repoRoot = discoverRepoRoot(__dirname);
  const buf = fs.readFileSync(path.join(repoRoot, relPath.split('/').join(path.sep)));
  const ctx = {
    kind: 'tree',
    scope: { kind: 'tree' },
    changed: [{ path: relPath, status: 'A', untracked: false, deleted: false }],
    changedPaths: [relPath],
    isEvidenceTree,
    readBytes: (p) => (p === relPath ? buf : null),
    readText: (p) => (p === relPath ? decodeBuffer(buf) : null),
    listTreeFiles: () => [relPath],
  };
  const out = createReport();
  const res = G8.subchecks.find((s) => s.id === 'G8-d').run(ctx, out);
  return {
    verdict: verdictOf(out.records, 'G8-d'),
    findings: res.findings.length,
    details: res.findings.map((f) => f.detail).join(' | '),
    detail: detailOf(out.records, 'G8-d'),
    text: decodeBuffer(buf),
  };
}

/** Run the real CLI in a child process and return its exit status + output. */
function gateCli(args) {
  const r = cp.spawnSync(process.execPath, [path.join(__dirname, 'gate.js')].concat(args), {
    cwd: discoverRepoRoot(__dirname),
    encoding: 'utf8',
    windowsHide: true,
  });
  return { status: r.status, out: (r.stdout || '') + (r.stderr || '') };
}

function runSubset(ctx, enable) {
  const out = createReport();
  const results = {};
  for (const sub of [...G6.subchecks, ...G7.subchecks, ...G8.subchecks]) {
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

/* ---------------------------------------------------------------------------
 * G5-b (role-file `model` closed set) and G5-a 2 probes (slice DIRECTIVE-MODEL-LISTS,
 * v1.18). T29-T37 below are what makes the new criterion falsifiable: each boundary
 * (verbatim comparison, case sensitivity, quoting, multi-key collection, the empty-value
 * division of labour, and the load-bearing closed set) is asserted in turn.
 * ------------------------------------------------------------------------ */

/** Path of the criterion under test (the mutation probe compiles a patched copy). */
const G5B_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g5b.js');
/** Mutation target (T37): the first member of the closed set. */
const G5B_FIRST_MEMBER = "  'deepseek-v4-pro',\n";
/** Mutation target (T42): the member added by this slice, whose removal must turn the suite red. */
const G5B_LUNA_MEMBER = "  'gpt-6-luna',\n";
/**
 * Mutation targets (T62/T64): the two GLM members landed by slice DIRECTIVE-GLM-LANDING
 * (v1.26). Each string is unique in the source (the `glm-5.3` line does not contain the
 * `glm-5.3-flash` line's quote-comma tail), so `split().join('')` removes exactly one line.
 */
const G5B_GLM_MEMBER = "  'glm-5.3',\n";
const G5B_GLM_FLASH_MEMBER = "  'glm-5.3-flash',\n";

/** G5-b verdict + finding count + record detail for one fixture key. */
function g5bProbeWith(mod, key) {
  const ctx = makeCtx([key], {});
  const out = createReport();
  const res = mod.run(ctx, out);
  return {
    verdict: verdictOf(out.records, 'G5-b'),
    findings: res && res.findings ? res.findings.length : -1,
    detail: detailOf(out.records, 'G5-b'),
  };
}

function g5bProbe(key) {
  return g5bProbeWith(G5B, key);
}

/**
 * G5-a 2 probe over one fixture key. The shipped G5-a 2 reads the role directory through
 * `ctx.listAgentFiles()`, so the probe supplies the single virtual agent file directly;
 * G5-a 1 is fed an empty git result because it greps the real role directory.
 */
function g5aEmptyProbe(key) {
  const text = fixtureText(key);
  const name = FIXTURES[key].virtual.split('/').pop();
  const ctx = {
    git: () => ({ ok: true, code: 0, out: '', err: '', spawnError: null }),
    listAgentFiles: () => [name],
    readText: (p) => (p === '.github/agents/' + name ? text : null),
  };
  const out = createReport();
  for (const s of G5A.subchecks) s.run(ctx, out);
  return {
    first: verdictOf(out.records, 'G5-a(1)'),
    second: verdictOf(out.records, 'G5-a(2)'),
    detail: detailOf(out.records, 'G5-a(2)'),
  };
}

/**
 * Compile a PATCHED copy of a criterion module (T37). Same technique as
 * `loadMutatedG6`: an in-memory copy cannot leak a mutated criterion onto disk, and the
 * shipped file is still verified byte-identical afterwards.
 */
function loadMutatedCriteria(sourcePath, from, to, tag) {
  const before = fs.readFileSync(sourcePath);
  const src = decodeBuffer(before);
  const patched = src.split(from).join(to);
  MUTATION_LOG.push({ key: path.basename(sourcePath) + ':' + tag, op: 'sourcePatch', needle: from, applied: patched !== src });
  const M = require('module');
  const mod = new M(sourcePath, null);
  mod.filename = sourcePath;
  mod.paths = M._nodeModulePaths(path.dirname(sourcePath));
  mod._compile(patched, sourcePath);
  return { exports: mod.exports, applied: patched !== src, identical: fs.readFileSync(sourcePath).equals(before) };
}

function loadMutatedG5b(from, to, tag) {
  return loadMutatedCriteria(G5B_SOURCE, from, to, tag);
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

  // ---- T29-T37: G5-b closed set (slice DIRECTIVE-MODEL-LISTS, v1.18) ----------------
  {
    const legal = g5bProbe('g5b-closed-set-legal');
    check(
      'T29 G5-b passes a model value that is a closed-set member',
      legal.verdict === 'PASS' && legal.findings === 0 && /checked 1/.test(legal.detail),
      'verdict=' + legal.verdict + ' findings=' + legal.findings + ' :: ' + legal.detail
    );

    const nokey = g5bProbe('g5b-legal-nokey');
    check(
      'T30 G5-b passes a role file that omits the model key (section 2.6 shape)',
      nokey.verdict === 'PASS' && nokey.findings === 0,
      'verdict=' + nokey.verdict + ' findings=' + nokey.findings
    );

    const kimi = g5bProbe('g5b-illegal-kimi');
    check(
      'T31 G5-b fails the one remaining blacklist member (kimi-k3 is not in the set)',
      kimi.verdict === 'FAIL' && kimi.findings === 1 && /model=kimi-k3/.test(kimi.detail),
      'verdict=' + kimi.verdict + ' findings=' + kimi.findings + ' :: ' + kimi.detail
    );

    const cs = g5bProbe('g5b-illegal-case');
    check(
      'T32 G5-b comparison is case-sensitive: the right model in the wrong case fails',
      cs.verdict === 'FAIL' && cs.findings === 1 && /model=GPT-5\.3-CODEX/.test(cs.detail),
      'verdict=' + cs.verdict + ' findings=' + cs.findings + ' :: ' + cs.detail
    );

    const quoted = g5bProbe('g5b-illegal-quoted');
    check(
      'T33 G5-b comparison is verbatim: a quoted set member still fails',
      quoted.verdict === 'FAIL' && quoted.findings === 1,
      'verdict=' + quoted.verdict + ' findings=' + quoted.findings + ' :: ' + quoted.detail
    );

    const multi = g5bProbe('g5b-illegal-multi');
    check(
      'T34 G5-b collects EVERY model key: a legal first key does not hide an illegal second one',
      multi.verdict === 'FAIL' && multi.findings === 1 && /model=GPT-5\.6 Luna/.test(multi.detail),
      'verdict=' + multi.verdict + ' findings=' + multi.findings + ' :: ' + multi.detail
    );

    const legacy = g5bProbe('g5b-illegal-legacy');
    check(
      'T35 G5-b fails a former blacklist value that is not a closed-set member',
      legacy.verdict === 'FAIL' && legacy.findings === 1 && /model=gpt-5\.4/.test(legacy.detail),
      'verdict=' + legacy.verdict + ' findings=' + legacy.findings + ' :: ' + legacy.detail
    );

    const empty = g5bProbe('g5b-empty-value');
    const emptyG5a = g5aEmptyProbe('g5b-empty-value');
    check(
      'T36 empty model value is G5-a 2 s finding and not G5-b s: G5-b passes while G5-a(2) fails',
      empty.verdict === 'PASS' && empty.findings === 0 && emptyG5a.second === 'FAIL' && emptyG5a.first === 'PASS',
      JSON.stringify({ g5b: empty.verdict, g5a1: emptyG5a.first, g5a2: emptyG5a.second, detail: emptyG5a.detail })
    );

    const mut = loadMutatedG5b(G5B_FIRST_MEMBER, G5B_FIRST_MEMBER + "  'gpt-5.4',\n", 'relax-set');
    check(
      'T37 the closed set is load-bearing: adding gpt-5.4 to the set flips that fixture to PASS',
      mut.applied === true && g5bProbeWith(mut.exports, 'g5b-illegal-legacy').verdict === 'PASS',
      'applied=' + mut.applied + ' patched_verdict=' + g5bProbeWith(mut.exports, 'g5b-illegal-legacy').verdict
    );
    check(
      'T37 restore: g5b.js is byte-identical and the cached module still fails gpt-5.4',
      mut.identical === true && g5bProbe('g5b-illegal-legacy').verdict === 'FAIL',
      'identical=' + mut.identical + ' cached=' + g5bProbe('g5b-illegal-legacy').verdict
    );
  }

  // ---- T38-T40: G5-b boundary fixtures raised by the 5 pre-review (F5) --------------
  {
    const bodyOnly = g5bProbe('g5b-body-only');
    check(
      'T38 G5-b parses the frontmatter block only: a model key in the body is not judged',
      bodyOnly.verdict === 'PASS' && bodyOnly.findings === 0,
      'verdict=' + bodyOnly.verdict + ' findings=' + bodyOnly.findings + ' :: ' + bodyOnly.detail
    );

    const trailing = g5bProbe('g5b-legal-trailing-space');
    check(
      'T39 G5-b trims both ends: a legal value with a trailing space still passes',
      trailing.verdict === 'PASS' && trailing.findings === 0,
      'verdict=' + trailing.verdict + ' findings=' + trailing.findings + ' :: ' + trailing.detail
    );

    const comment = g5bProbe('g5b-illegal-inline-comment');
    check(
      'T40 G5-b is verbatim: a legal value followed by an inline comment fails',
      comment.verdict === 'FAIL' && comment.findings === 1 && comment.detail.includes('model=gpt-5.3-codex # keep'),
      'verdict=' + comment.verdict + ' findings=' + comment.findings + ' :: ' + comment.detail
    );
  }

  // ---- T41-T42: the two members this slice added are load-bearing (7 round-one M4) ----
  {
    const sol = g5bProbe('g5b-closed-set-legal-sol');
    const luna = g5bProbe('g5b-closed-set-legal-luna');
    check(
      'T41 G5-b accepts each member this slice added (deep tier and parallel alternative)',
      sol.verdict === 'PASS' && sol.findings === 0 && luna.verdict === 'PASS' && luna.findings === 0,
      JSON.stringify({ sol: sol.verdict, luna: luna.verdict, solDetail: sol.detail, lunaDetail: luna.detail })
    );

    const drop = loadMutatedG5b(G5B_LUNA_MEMBER, '', 'drop-luna-member');
    check(
      'T42 the added member is load-bearing: dropping gpt-6-luna from the set makes that fixture fail',
      drop.applied === true && g5bProbeWith(drop.exports, 'g5b-closed-set-legal-luna').verdict === 'FAIL',
      'applied=' + drop.applied + ' patched_verdict=' + g5bProbeWith(drop.exports, 'g5b-closed-set-legal-luna').verdict
    );
    check(
      'T42 restore: g5b.js is byte-identical and the cached module still accepts gpt-6-luna',
      drop.identical === true && g5bProbe('g5b-closed-set-legal-luna').verdict === 'PASS',
      'identical=' + drop.identical + ' cached=' + g5bProbe('g5b-closed-set-legal-luna').verdict
    );
  }

  // ---- T43-T49: G8 frozen-evidence-pack self-consistency + the 6.3 scope exclusion
  // (slice DIRECTIVE-EVIDENCE-SCOPE, v1.24). Each assertion ties a verdict to a named
  // reason token, so a FAIL is meaningful and a PASS is not vacuous.
  {
    const FORMS = 'tools/gates/evidence-forms.txt';
    const REG_PATH = 'docs/validation/evidence/fixture-g8-a-registered.md';

    // T43: G8-a shape guard - three allowed forms plus the stray-root abuse.
    const member = g8Probe('G8-a', ['g8-a-member']);
    // `info === false` is load-bearing: without it a vacuous INFO (no evidence path in
    // the change set) would satisfy `verdict === 'PASS' && findings === 0` even though
    // the path was never judged at all.
    check('T43 G8-a accepts a pack member (nearest ancestor carries MANIFEST + SHA256SUMS)',
      member.verdict === 'PASS' && member.findings === 0 && member.info === false,
      'verdict=' + member.verdict + ' findings=' + member.findings + ' info=' + member.info);
    const freedom = g8Probe('G8-a', ['g8-a-freedom-list']);
    check('T43 G8-a accepts a root-level *-freedom-list.md',
      freedom.verdict === 'PASS' && freedom.findings === 0 && freedom.info === false,
      'verdict=' + freedom.verdict + ' findings=' + freedom.findings + ' info=' + freedom.info);
    const registered = g8Probe('G8-a', ['g8-a-registered'], {
      whitelists: { [FORMS]: [{ match: REG_PATH, file: REG_PATH, reason: 'T43 probe: verbatim registration' }] },
    });
    check('T43 G8-a accepts a path registered verbatim in evidence-forms.txt',
      registered.verdict === 'PASS' && registered.findings === 0 && registered.info === false,
      'verdict=' + registered.verdict + ' findings=' + registered.findings + ' info=' + registered.info + ' :: ' + registered.details);
    const unreg = g8Probe('G8-a', ['g8-a-registered']);
    check('T43 the registration is load-bearing: without it the same path fails',
      unreg.verdict === 'FAIL' && unreg.findings === 1 && /not registered/.test(unreg.details),
      'verdict=' + unreg.verdict + ' :: ' + unreg.details);
    const stray = g8Probe('G8-a', ['g8-a-stray-root']);
    check('T43 G8-a FAILs a stray root-level path with a `not a pack member` reason token',
      stray.verdict === 'FAIL' && /not a pack member/.test(stray.details),
      'verdict=' + stray.verdict + ' :: ' + stray.details);
    const subReg = g8Probe('G8-a', ['g8-a-stray-root'], {
      whitelists: {
        [FORMS]: [
          {
            match: 'fixture-g8-a-stray-root',
            file: 'docs/validation/evidence/fixture-g8-a-stray-root.md',
            reason: 'T43 probe: substring match must not register',
          },
        ],
      },
    });
    // Strengthened: a bare FAIL could be produced by ANY unrelated finding. Require the
    // exact `not a pack member` reason (the only cause here) and exactly one finding, so
    // a criterion that started rejecting every path could not satisfy this by accident.
    check('T43 a SUBSTRING match does not register (entry.match must EQUAL the path)',
      subReg.verdict === 'FAIL' && subReg.findings === 1 && /not a pack member/.test(subReg.details),
      'verdict=' + subReg.verdict + ' findings=' + subReg.findings + ' :: ' + subReg.details);
    let formsErr = null;
    try {
      g8Probe('G8-a', ['g8-ok-pack'], { whitelistText: { [FORMS]: 'a\tb\n' } });
    } catch (e) {
      formsErr = e;
    }
    check('T43 a registration row with 2 fields is a UsageError (exit 2, never a silent pass)',
      formsErr instanceof UsageError && /3 TAB-separated/.test(formsErr.message),
      'err=' + (formsErr ? formsErr.message : 'NO ERROR (fail-open)'));
    let dupErr = null;
    try {
      g8Probe('G8-a', ['g8-ok-pack'], { whitelistText: { [FORMS]: 'p\tp\tr1\np\tp\tr2\n' } });
    } catch (e) {
      dupErr = e;
    }
    check('T43 a duplicate registration for the same path is a UsageError',
      dupErr instanceof UsageError && /duplicate/.test(dupErr.message),
      'err=' + (dupErr ? dupErr.message : 'NO ERROR (fail-open)'));

    // T44: G8-b pack self-consistency.
    const okPack = g8Probe('G8-b', ['g8-ok-pack']);
    check('T44 G8-b accepts a fully self-consistent pack (runtime-computed hashes)',
      okPack.verdict === 'PASS' && okPack.findings === 0 && okPack.info === false,
      'verdict=' + okPack.verdict + ' info=' + okPack.info + ' :: ' + okPack.details);
    const missingList = g8Probe('G8-b', ['g8-b-missing-list']);
    check('T44 G8-b FAILs a pack without the 归档条目清单 section (token `missing`, names the pack)',
      missingList.verdict === 'FAIL' && missingList.findings === 2 &&
        /fixture-g8-b-missing-list/.test(missingList.details) &&
        /missing `## 归档条目清单` section/.test(missingList.details) &&
        /not listed in `## 归档条目清单`/.test(missingList.details),
      'verdict=' + missingList.verdict + ' findings=' + missingList.findings + ' :: ' + missingList.details);
    const unlisted = g8Probe('G8-b', ['g8-b-unlisted-payload']);
    check('T44 G8-b FAILs an unlisted pack-root payload (token `not listed`, names the file in both lists)',
      unlisted.verdict === 'FAIL' && unlisted.findings === 2 && /not listed/.test(unlisted.details) && /CLEARED-INVENTORY\.txt/.test(unlisted.details) && /SHA256SUMS/.test(unlisted.details),
      'verdict=' + unlisted.verdict + ' findings=' + unlisted.findings + ' :: ' + unlisted.details);
    const hashMismatch = g8Probe('G8-b', ['g8-b-hash-mismatch']);
    check('T44 G8-b recomputes hashes and FAILs a mismatch (token `hash mismatch`, names the payload)',
      hashMismatch.verdict === 'FAIL' && hashMismatch.findings === 1 && /archive\/a\.txt hash mismatch/.test(hashMismatch.details),
      'verdict=' + hashMismatch.verdict + ' findings=' + hashMismatch.findings + ' :: ' + hashMismatch.details);
    const conflict = g8Probe('G8-b', ['g8-b-untracked-conflict']);
    check('T44 G8-b FAILs a 未入库 entry colliding with a payload (token `untracked-list conflict`, names the path)',
      conflict.verdict === 'FAIL' && conflict.findings === 1 && /untracked-list conflict/.test(conflict.details) && /archive\/a\.txt/.test(conflict.details),
      'verdict=' + conflict.verdict + ' findings=' + conflict.findings + ' :: ' + conflict.details);
    const backslash = g8Probe('G8-b', ['g8-b-backslash-sums']);
    check('T44 G8-b normalizes backslash SHA256SUMS paths to POSIX and passes',
      backslash.verdict === 'PASS' && backslash.findings === 0 && backslash.info === false,
      'verdict=' + backslash.verdict + ' info=' + backslash.info + ' :: ' + backslash.details);
    // Reverse mutation: normalization must not become a blanket pass. A wrong hash on the
    // SAME backslash-separated line must still be caught after normalization.
    const backslashBad = g8Probe('G8-b', ['g8-b-backslash-sums'], {
      mutations: [
        {
          key: 'g8-b-backslash-sums',
          path: 'docs/validation/evidence/fixture-g8-b-backslash-sums/SHA256SUMS.txt',
          op: 'flipHashDigit',
        },
      ],
    });
    check('T44 the backslash normalization is load-bearing: a wrong hash on the same line still FAILs',
      backslashBad.verdict === 'FAIL' && backslashBad.findings === 1 && /archive\/a\.txt hash mismatch/.test(backslashBad.details),
      'verdict=' + backslashBad.verdict + ' findings=' + backslashBad.findings + ' :: ' + backslashBad.details);
    const flip = g8Probe('G8-b', ['g8-ok-pack'], {
      mutations: [
        {
          key: 'g8-ok-pack',
          path: 'docs/validation/evidence/fixture-g8-ok-pack/SHA256SUMS.txt',
          op: 'flipHashDigit',
        },
      ],
    });
    check('T44 a single flipped hash digit turns G8-b red (the recompute is load-bearing)',
      flip.verdict === 'FAIL' && flip.findings === 1 && /hash mismatch/.test(flip.details),
      'verdict=' + flip.verdict + ' findings=' + flip.findings + ' :: ' + flip.details);

    // T45: G8-c encoding-normalization record.
    const cOk = g8Probe('G8-c', ['g8-c-ok']);
    check('T45 G8-c accepts a MANIFEST whose archive hash equals the SHA256SUMS record',
      cOk.verdict === 'PASS' && cOk.findings === 0 && cOk.info === false,
      'verdict=' + cOk.verdict + ' info=' + cOk.info + ' :: ' + cOk.details);
    const cBad = g8Probe('G8-c', ['g8-c-bad-hash']);
    check('T45 G8-c FAILs an archive hash that disagrees with SHA256SUMS (token `hash mismatch`, names the payload)',
      cBad.verdict === 'FAIL' && cBad.findings === 1 && /hash mismatch for archive\/slice\.txt/.test(cBad.details),
      'verdict=' + cBad.verdict + ' findings=' + cBad.findings + ' :: ' + cBad.details);
    const cNone = g8Probe('G8-c', ['g8-c-no-section']);
    // The exemption must leave an explicit trace: a PASS whose detail names the absent
    // section heading. `info === false` proves the pack was actually examined (a vacuous
    // INFO would also be a PASS whose detail says nothing).
    check('T45 G8-c exempts a legacy pack without the 后缀与编码映射 section (explicit note, not a blank PASS)',
      cNone.verdict === 'PASS' && cNone.findings === 0 && cNone.info === false && /exempt/.test(cNone.detail) && /后缀与编码映射/.test(cNone.detail),
      'verdict=' + cNone.verdict + ' findings=' + cNone.findings + ' info=' + cNone.info + ' detail=' + cNone.detail);

    // T46: G8-d rebuilt-archive note.
    const dOk = g8Probe('G8-d', ['g8-d-ok']);
    check('T46 G8-d accepts a 重建归档 note carrying source + basis and no bad sha token',
      dOk.verdict === 'PASS' && dOk.findings === 0 && dOk.info === false,
      'verdict=' + dOk.verdict + ' info=' + dOk.info + ' :: ' + dOk.details);
    const dBad = g8Probe('G8-d', ['g8-d-missing-basis']);
    check('T46 G8-d FAILs a 重建归档 note missing the basis (token `basis`, names the file)',
      dBad.verdict === 'FAIL' && dBad.findings === 1 && /fixture-g8-d-missing-basis/.test(dBad.details) && /basis/.test(dBad.details),
      'verdict=' + dBad.verdict + ' findings=' + dBad.findings + ' :: ' + dBad.details);

    // T47: the OB-79 scope exclusion - the same violating content in two locations.
    // Strengthened: the positive side must show EVERY criterion actually ran and passed
    // (a missing record would otherwise satisfy `failIds().length === 0` vacuously), and
    // the negative side must name the exact red set + path instead of `>= 1 FAIL`.
    const NEG_PATH = 'docs/t027/fixture-g8-exclusion.md';
    const pos = runMods(makeCtx(['exclusion-positive']), [G1A, G2, G4A, G4B, G6]);
    check('T47 inside the tree the violating content is NOT judged by G1-a/G2/G4a/G4b/G6',
      failIds(pos.records).length === 0 &&
        ['G1-a', 'G2', 'G4a', 'G4b'].every((id) => verdictOf(pos.records, id) === 'PASS') &&
        !pos.records.some((r) => r.detail && r.detail.indexOf('violation.md') >= 0),
      'fails=' + JSON.stringify(failIds(pos.records)) + ' details=' + pos.records.map((r) => r.id + '=' + r.verdict).join(','));
    check('T47 G4a reports the excluded count instead of judging the tree file',
      /\(1 excluded: evidence tree\)/.test(detailOf(pos.records, 'G4a')), detailOf(pos.records, 'G4a'));
    const neg = runMods(makeCtx(['exclusion-negative']), [G1A, G2, G4A, G4B, G6]);
    check('T47 the SAME content outside the tree IS judged: G1-a/G2/G4a/G4b/G6-4 all turn red and name the path',
      ['G1-a', 'G2', 'G4a', 'G4b', 'G6-4'].every((id) => failIds(neg.records).indexOf(id) >= 0) &&
        ['G1-a', 'G2', 'G4a', 'G6-4'].every((id) => detailOf(neg.records, id).indexOf(NEG_PATH) >= 0),
      'fails=' + JSON.stringify(failIds(neg.records)) + ' :: ' + neg.records.map((r) => r.id + '=' + r.detail).join(' ;; '));

    // T48: G8 never produces a vacuous PASS - asserted for ALL FOUR sub-checks.
    const noneB = g8Probe('G8-b', ['g6-1-pending']);
    const noneA = g8Probe('G8-a', ['g6-1-pending']);
    const noneC = g8Probe('G8-c', ['g6-1-pending']);
    const noneD = g8Probe('G8-d', ['g6-1-pending']);
    check('T48 every G8 sub-check reports INFO (not a vacuous PASS) when the change set has no evidence path',
      [noneA, noneB, noneC, noneD].every((p) => p.info === true && p.findings === 0),
      JSON.stringify([noneA, noneB, noneC, noneD].map((p) => p.info + '/' + p.verdict)));
    const someB = g8Probe('G8-b', ['g8-ok-pack']);
    check('T48 the INFO is the absence of an evidence path, not a blanket: the same sub-check is NOT INFO when one is present',
      someB.info === false, 'info=' + someB.info + ' verdict=' + someB.verdict);

    // T49: the ctx primitives G8 is built on.
    const lctx = makeCtx(['g8-ok-pack']);
    const packFiles = lctx.listTreeFiles('docs/validation/evidence/fixture-g8-ok-pack');
    check('T49 ctx.listTreeFiles returns the recursive, sorted pack file list',
      JSON.stringify(packFiles) ===
        JSON.stringify([
          'docs/validation/evidence/fixture-g8-ok-pack/MANIFEST.md',
          'docs/validation/evidence/fixture-g8-ok-pack/SHA256SUMS.txt',
          'docs/validation/evidence/fixture-g8-ok-pack/archive/a.txt',
          'docs/validation/evidence/fixture-g8-ok-pack/archive/nasty.md',
        ]),
      JSON.stringify(packFiles));
    // Strengthened: the trailing-slash form and the two near-miss prefixes (a longer name
    // that merely STARTS with the tree name, and a sibling that does not) pin the exact
    // boundary, not just one positive and one negative sample.
    const evCases = [
      ['docs/validation/evidence/x', true],
      ['docs/Validation/Evidence/x', true],
      ['docs/validation/evidence/archive/deep/x.md', true],
      ['docs/validation/evidence/', true],
      ['docs/validation/evidence', false],
      ['docs/validation/evidencex', false],
      ['docs/validation/evidence-docs/x', false],
      ['docs/validation/x', false],
    ];
    check('T49 isEvidenceTree is case-insensitive, prefix-anchored and does not over-match near misses',
      evCases.every(([p, want]) => isEvidenceTree(p) === want),
      JSON.stringify(evCases.map(([p]) => p + '=' + isEvidenceTree(p))));
  }

  // ---- T50-T60: the 6.3 evidence-tree scope exclusion, the G8 shape guard and the
  // DEF-2 batched-ignore anchor (slice DIRECTIVE-EVIDENCE-SCOPE, ③ tester hardening).
  //
  // T50/T51/T52 run the exclusion through a REAL lib/ctx context (`injectedCtx`), so a
  // regression inside lib/ctx.js - not merely inside the selftest's local `makeCtx`
  // re-derivation - turns them red. That is what the T56 on-disk mutation experiment
  // (recorded in the slice report) demonstrates: removing the exclusion in lib/ctx.js
  // flips exactly these assertions.
  {
    const FORMS = 'tools/gates/evidence-forms.txt';
    const VIOLATING_TREE = 'docs/validation/evidence/fixture-g8-t50/archive/probe.md';
    const VIOLATING_OUT = 'docs/t027/fixture-g8-t50-probe.md';

    // ---- T50: the exclusion duality - byte-identical content, two locations, plus proof
    // that G8 still judges the tree path (the exclusion must not blind the tree's own guard).
    const inTree = runKeyMods('g8-t50-in-tree', [G1A, G2, G4A, G4B, G6]);
    check('T50 in-tree violating content yields ZERO finding from G1-a/G2/G4a/G4b/G6',
      failIds(inTree.records).length === 0 &&
        ['G1-a', 'G2', 'G4a', 'G4b', 'G6-4', 'G6-6'].every((id) => verdictOf(inTree.records, id) === 'PASS') &&
        !inTree.records.some((r) => r.detail && r.detail.indexOf(VIOLATING_TREE) >= 0),
      'fails=' + JSON.stringify(failIds(inTree.records)) + ' :: ' + inTree.records.map((r) => r.id + '=' + r.verdict).join(','));
    const outTree = runKeyMods('g8-t50-outside', [G1A, G2, G4A, G4B, G6]);
    check('T50 the SAME bytes at docs/t027/ turn G1-a/G2/G4a/G4b/G6-4/G6-6 red and name the path',
      ['G1-a', 'G2', 'G4a', 'G4b', 'G6-4', 'G6-6'].every((id) => failIds(outTree.records).indexOf(id) >= 0) &&
        ['G1-a', 'G2', 'G4a', 'G6-4', 'G6-6'].every((id) => detailOf(outTree.records, id).indexOf(VIOLATING_OUT) >= 0),
      'fails=' + JSON.stringify(failIds(outTree.records)) + ' :: ' + outTree.records.map((r) => r.id + '=' + r.detail).join(' ;; '));
    const g50a = g8Probe('G8-a', ['g8-t50-in-tree']);
    const g50b = g8Probe('G8-b', ['g8-t50-in-tree']);
    check('T50 G8 still judges the in-tree file: G8-a PASS + G8-b PASS on the pack carrying the violation',
      g50a.verdict === 'PASS' && g50a.findings === 0 && g50a.info === false &&
        g50b.verdict === 'PASS' && g50b.findings === 0 && g50b.info === false,
      JSON.stringify({ a: g50a.verdict, aInfo: g50a.info, b: g50b.verdict, bInfo: g50b.info }));
    const strayA = g8Probe('G8-a', ['exclusion-positive']);
    check('T50 G8-a still FAILs a non-member tree path (the shape guard is not inert)',
      strayA.verdict === 'FAIL' && /not a pack member/.test(strayA.details),
      'verdict=' + strayA.verdict + ' :: ' + strayA.details);

    // ---- T51: changedMarkdown is the SINGLE source the four criteria share.
    const texts51 = new Map();
    for (const k of ['g8-t50-in-tree', 'g8-t50-outside']) for (const [p, t] of g8Texts(k)) texts51.set(p, t);
    const realCtx51 = injectedCtx([...texts51.keys()], texts51);
    check('T51 ctx.changedMarkdown() (shipped lib/ctx) drops the tree .md and keeps the docs/t027 .md',
      JSON.stringify(realCtx51.changedMarkdown()) === JSON.stringify([VIOLATING_OUT]),
      JSON.stringify(realCtx51.changedMarkdown()));
    check('T51 the selftest ctx and lib/ctx agree on changedMarkdown (no local re-derivation drift)',
      JSON.stringify(makeCtx(['g8-t50-in-tree', 'g8-t50-outside']).changedMarkdown()) ===
        JSON.stringify(realCtx51.changedMarkdown()),
      JSON.stringify(makeCtx(['g8-t50-in-tree', 'g8-t50-outside']).changedMarkdown()));
    const perCrit = [
      ['G1-a', [G1A], VIOLATING_OUT],
      ['G2', [G2], VIOLATING_OUT],
      ['G4b', [G4B], null],
      ['G6', [G6], VIOLATING_OUT],
    ];
    for (const [name, mods, pathInDetail] of perCrit) {
      // Two independent sides: a TREE-ONLY change set must be silent (this is the side that
      // turns red if the exclusion is removed, because these details do not name the path),
      // and the docs/t027 side must be red (the behaviour the exclusion must preserve).
      const treeOnly = runKeyMods('g8-t50-in-tree', mods).records;
      const mixed = runMods(realCtx51, mods).records;
      const treeRed = treeOnly.some((r) => r.verdict === 'FAIL');
      const outRed = mixed.some(
        (r) => r.verdict === 'FAIL' && (!pathInDetail || (r.detail && r.detail.indexOf(pathInDetail) >= 0))
      );
      check('T51 ' + name + ' shares the single exclusion (tree-only ctx green, docs/t027 ctx red)',
        !treeRed && outRed,
        'treeOnly=' + treeOnly.map((r) => r.id + '=' + r.verdict).join(',') + ' ;; mixed=' +
          mixed.map((r) => r.id + '=' + r.verdict).join(','));
    }

    // ---- T52: G4a counts the excluded tree files instead of judging their bytes.
    const bomIn = runKeyMods('g8-bom-in-tree', [G4A]);
    check('T52 G4a excludes a BOM+CRLF .md and a BOM-less .md inside the tree and counts BOTH',
      verdictOf(bomIn.records, 'G4a') === 'PASS' && /\(2 excluded: evidence tree\)/.test(detailOf(bomIn.records, 'G4a')),
      'verdict=' + verdictOf(bomIn.records, 'G4a') + ' detail=' + detailOf(bomIn.records, 'G4a'));
    const bomOut = runKeyMods('g8-bom-outside', [G4A]);
    const bomOutDetail = detailOf(bomOut.records, 'G4a');
    check('T52 the identical bytes outside the tree DO fail G4a (the exclusion is the only reason the in-tree run is green)',
      verdictOf(bomOut.records, 'G4a') === 'FAIL' &&
        /fixture-g8-bom-crlf\.md bom=true crlf=true/.test(bomOutDetail) &&
        /fixture-g8-bom-nobom\.md bom=false crlf=false/.test(bomOutDetail),
      'verdict=' + verdictOf(bomOut.records, 'G4a') + ' detail=' + bomOutDetail);

    // ---- T53: G8-b breakage matrix augmentation (parse failures must fail closed).
    const malformed = g8Probe('G8-b', ['g8-b-malformed-sums']);
    check('T53 a malformed SHA256SUMS line FAILs and names the line/pack, never a silent pass',
      malformed.verdict === 'FAIL' && /malformed SHA256SUMS line 1/.test(malformed.details) &&
        /fixture-g8-b-malformed-sums/.test(malformed.details),
      'verdict=' + malformed.verdict + ' :: ' + malformed.details);
    const truncated = g8Probe('G8-b', ['g8-b-manifest-truncated']);
    check('T53 a truncated 归档条目清单 (heading present, table empty) FAILs and names the unlisted payload',
      truncated.verdict === 'FAIL' && /archive\/a\.txt not listed/.test(truncated.details),
      'verdict=' + truncated.verdict + ' :: ' + truncated.details);

    // ---- T54: G8-c payload/record consistency at both boundaries.
    const cNoPayload = g8Probe('G8-c', ['g8-c-arch-not-in-payload']);
    check('T54 an 归档路径 absent from the pack payload FAILs G8-c and names the path',
      cNoPayload.verdict === 'FAIL' && /not in pack payload: archive\/missing\.txt/.test(cNoPayload.details),
      'verdict=' + cNoPayload.verdict + ' :: ' + cNoPayload.details);
    const cNoSums = g8Probe('G8-c', ['g8-c-arch-not-in-sums']);
    check('T54 an 归档路径 present in the payload but absent from SHA256SUMS FAILs G8-c and names the path',
      cNoSums.verdict === 'FAIL' && /not recorded in SHA256SUMS: archive\/slice\.txt/.test(cNoSums.details),
      'verdict=' + cNoSums.verdict + ' :: ' + cNoSums.details);

    // ---- T55: G8-d annotation-zone boundary.
    const bodyOnly = g8Probe('G8-d', ['g8-d-body-only']);
    check('T55 a 重建归档 mention OUTSIDE the annotation zone does not trigger G8-d',
      bodyOnly.verdict === 'PASS' && bodyOnly.findings === 0 && bodyOnly.info === false,
      'verdict=' + bodyOnly.verdict + ' findings=' + bodyOnly.findings + ' info=' + bodyOnly.info);
    const badSha = g8Probe('G8-d', ['g8-d-bad-sha']);
    check('T55 a non-64-hex sha256 token INSIDE the annotation zone FAILs G8-d',
      badSha.verdict === 'FAIL' && /bad sha256 token/.test(badSha.details) && /\(40 digits\)/.test(badSha.details),
      'verdict=' + badSha.verdict + ' :: ' + badSha.details);
    const precedent = g8dRealProbe('docs/validation/evidence/DIRECTIVE-MODELID-MAP-mutation.md');
    const precedentZone = G8.internals
      .annotationZone(precedent.text.split('\n'))
      .map((z) => z.text)
      .join('\n');
    check('T55 the precedent file is not misjudged despite 重建归档 + a 16-digit sha-original short hash',
      precedent.verdict === 'PASS' && precedent.findings === 0 &&
        precedent.text.includes('重建归档') && /sha-original=[0-9a-f]{16}/.test(precedent.text),
      'verdict=' + precedent.verdict + ' findings=' + precedent.findings + ' :: ' + precedent.details);
    check('T55 the precedent 16-digit short hash sits OUTSIDE the annotation zone (why it is never a sha256 token)',
      /sha-original=[0-9a-f]{16}/.test(precedent.text) && !/[0-9a-f]{16}/.test(precedentZone),
      JSON.stringify({ zone: precedentZone.slice(0, 120) }));

    // ---- T56: the recorded on-disk mutation anchors must exist verbatim. (The mutation
    // run itself is executed out-of-band and its sha256 evidence is recorded in the slice
    // report; these assertions stop the anchors from rotting silently.)
    const ctxSrc = decodeBuffer(fs.readFileSync(path.join(__dirname, 'lib', 'ctx.js')));
    const g8Src = decodeBuffer(fs.readFileSync(path.join(__dirname, 'lib', 'criteria', 'g8.js')));
    check('T56 mutation anchor (a) exists verbatim in lib/ctx.js (the exclusion being mutated)',
      ctxSrc.includes('changedMarkdown: () => changedPaths.filter((p) => /\\.md$/i.test(p) && !isEvidenceTree(p)),'),
      'hasExclusion=' + ctxSrc.includes('!isEvidenceTree(p)'));
    check('T56 mutation anchor (b) exists verbatim in lib/criteria/g8.js (the hash recompute being mutated)',
      g8Src.includes('const actual = sha256Hex(bytes);') && g8Src.includes('if (actual !== recorded) {'),
      JSON.stringify({ a: g8Src.includes('const actual = sha256Hex(bytes);'), b: g8Src.includes('if (actual !== recorded) {') }));

    // ---- T57: scope readings and the range branch of listTreeFiles.
    const ciG8a = gateCli(['--check=G8-a', '--scope=ci']);
    check('T57 --scope=ci is a valid reading (exit 0/1, never the usage-error exit 2)',
      ciG8a.status === 0 || ciG8a.status === 1,
      'status=' + ciG8a.status + ' out=' + ciG8a.out.trim().split('\n').slice(-1)[0]);
    const badRange = gateCli(['--check=G8-a', '--scope=range:deadbeef^..HEAD']);
    check('T57 an unresolvable range endpoint exits 2 (usage/environment error), not a silent green',
      badRange.status === 2, 'status=' + badRange.status + ' out=' + badRange.out.trim());
    // DEF-1 (fixed in slice DIRECTIVE-EVIDENCE-SCOPE): `run()`'s catch block used to
    // `return {code:2, stderr}` BEFORE its stderr write, so this exact command exited 2
    // with ZERO output - a failure a caller cannot diagnose. `gateCli` above merges
    // stdout+stderr, which cannot tell a silent exit from a diagnosed one, so pin stderr
    // separately below: if the catch ever again skips `process.stderr.write`, the captured
    // stderr is empty and this check turns RED (that is precisely DEF-1's shape). The
    // assertion deliberately fixes only exit code + non-empty + `gate.js:` prefix, never
    // the wording of the message.
    const badRangeSplit = cp.spawnSync(
      process.execPath,
      [path.join(__dirname, 'gate.js'), '--check=G8-a', '--scope=range:deadbeef^..HEAD'],
      { cwd: discoverRepoRoot(__dirname), encoding: 'utf8', windowsHide: true }
    );
    const brStderr = badRangeSplit.stderr || '';
    check(
      'T57 DEF-1 an unresolvable range endpoint exits 2 with a NON-EMPTY stderr carrying the `gate.js:` diagnosis (never a silent exit that degrades into an empty change set)',
      badRangeSplit.status === 2 && brStderr.trim() !== '' && /gate\.js:/.test(brStderr),
      JSON.stringify({ status: badRangeSplit.status, stderrLen: brStderr.trim().length, hasPrefix: /gate\.js:/.test(brStderr) })
    );
    // Fail-open breadcrumb: if DEF-1 regresses the check above turns red; this INFO keeps the
    // old silent-exit signature visible in the printed report while it is broken.
    if (badRange.status === 2 && badRange.out.trim() === '') {
      info(
        'T57 DEF-1 REGRESSED: an unresolvable range exits 2 with EMPTY stdout+stderr (gate.js run() swallows the thrown UsageError diagnosis)',
        'min-repro: node tools/gates/gate.js --check=G8-a --scope=range:deadbeef^..HEAD'
      );
    }
    const realRoot = discoverRepoRoot(__dirname);
    if (!revAvailable(realRoot, '4284e94')) {
      info('T57 historical range reading SKIPPED (INFO): revision 4284e94 is unavailable in this clone');
    } else {
      const hist = gateCli(['--check=G8-a', '--scope=range:4284e94^..4284e94']);
      check('T57 the historical range 4284e94 is a valid reading (exit 0/1, not 2)',
        hist.status === 0 || hist.status === 1,
        'status=' + hist.status + ' out=' + hist.out.trim().split('\n').slice(-1)[0]);
    }
    const spyCalls = [];
    const treeSpy = (root, args) => {
      const key = args.join(' ');
      spyCalls.push(key);
      if (/rev-parse --verify --quiet/.test(key)) return { ok: true, code: 0, out: 'abc\n', err: '', spawnError: null };
      if (/^diff --name-status/.test(key)) return { ok: true, code: 0, out: '', err: '', spawnError: null };
      if (/^ls-tree -r -z --name-only/.test(key)) {
        return { ok: true, code: 0, out: 'dir/a.txt\0dir/b.txt\0', err: '', spawnError: null };
      }
      return { ok: true, code: 0, out: '', err: '', spawnError: null };
    };
    const rangeCtx = createContext({ repoRoot: realRoot, scope: 'range:HEAD..HEAD', spawnGit: treeSpy });
    const listed = rangeCtx.listTreeFiles('dir');
    check('T57 listTreeFiles in range scope walks `git ls-tree -r -z --name-only <endpoint> -- <dir>` (not the disk)',
      JSON.stringify(listed) === JSON.stringify(['dir/a.txt', 'dir/b.txt']) &&
        spyCalls.some((c) => /^ls-tree -r -z --name-only HEAD -- dir$/.test(c)),
      JSON.stringify({ listed, lsTree: spyCalls.filter((c) => /^ls-tree/.test(c)) }));

    // ---- T58: the registration table is read by ENDPOINT - an implementation decision
    // that must stay pinned so ⑤⑥⑦ can re-verify it.
    const worktreeRead = createContext({ repoRoot: realRoot, scope: 'tree' }).loadWhitelist(FORMS);
    check('T58 tree scope reads the registration table from the working tree',
      worktreeRead.source === 'worktree' && worktreeRead.entries.length >= 1,
      'source=' + worktreeRead.source + ' entries=' + worktreeRead.entries.length);
    const REG_REAL = 'docs/validation/evidence/DIRECTIVE-MODELID-MAP-mutation.md';
    check('T58 the real registered root-level path is listed VERBATIM (match === file) in the worktree registry',
      worktreeRead.entries.some((e) => e.file === REG_REAL && e.match === REG_REAL),
      JSON.stringify(worktreeRead.entries.map((e) => e.file)));
    if (!revAvailable(realRoot, '4284e94')) {
      info('T58 historical endpoint reading SKIPPED (INFO): revision 4284e94 is unavailable in this clone');
    } else {
      const endpointRead = createContext({ repoRoot: realRoot, scope: 'range:4284e94^..4284e94' }).loadWhitelist(FORMS);
      check('T58 a historical range endpoint reads the registration table from the endpoint (absent then ⇒ zero entries)',
        /^endpoint:/.test(endpointRead.source) && endpointRead.text === null && endpointRead.entries.length === 0,
        JSON.stringify({ source: endpointRead.source, text: endpointRead.text, entries: endpointRead.entries.length }));
      const withReg = g8Probe('G8-a', ['g8-a-real-registered'], { whitelists: { [FORMS]: worktreeRead.entries } });
      const withoutReg = g8Probe('G8-a', ['g8-a-real-registered'], { whitelists: { [FORMS]: [] } });
      check('T58 the SAME path is green with the worktree registry and red with the endpoint (empty) registry',
        withReg.verdict === 'PASS' && withReg.findings === 0 && withReg.info === false &&
          withoutReg.verdict === 'FAIL' && /not registered/.test(withoutReg.details),
        JSON.stringify({ with: withReg.verdict, without: withoutReg.verdict, withoutDetail: withoutReg.details }));
    }

    // ---- T59: whole-corpus anchor strengthening (T13's set equality already pins the ids;
    // these add the vocabulary / ordering invariants that keep the anchor meaningful).
    check('T59 EXPECTED_CORPUS_FAILS is sorted and duplicate-free',
      JSON.stringify(EXPECTED_CORPUS_FAILS) === JSON.stringify([...new Set(EXPECTED_CORPUS_FAILS)].sort()),
      JSON.stringify(EXPECTED_CORPUS_FAILS));
    const KNOWN_IDS = new Set([...ALL_SUBCHECKS, 'G1-a', 'G2', 'G4a', 'G4b']);
    check('T59 every expected corpus FAIL id belongs to the known criterion vocabulary',
      EXPECTED_CORPUS_FAILS.every((id) => KNOWN_IDS.has(id)),
      JSON.stringify(EXPECTED_CORPUS_FAILS.filter((id) => !KNOWN_IDS.has(id))));
    check('T59 no criterion id is duplicated in ALL_SUBCHECKS',
      new Set(ALL_SUBCHECKS).size === ALL_SUBCHECKS.length, 'n=' + ALL_SUBCHECKS.length);
    check('T59 every G8 sub-check id is present in ALL_SUBCHECKS',
      G8.subchecks.every((s) => ALL_SUBCHECKS.indexOf(s.id) >= 0),
      JSON.stringify(ALL_SUBCHECKS));

    // ---- T60: DEF-2 regression anchor. The tree/index branch of listTreeFiles must
    // resolve ignore-ness with EXACTLY ONE batched `git check-ignore --stdin -z` spawn for
    // a whole directory, and must NOT spawn one `git check-ignore -q -- <path>` per file.
    // The per-file form IS DEF-2: on the measured host one check-ignore spawn costs ~417 ms,
    // so the evidence tree (~2671 files) took ~30 min for a single enumeration, and the
    // G8-a ancestor walk extrapolated to ~74 min. This assertion runs the SHIPPED lib/ctx
    // (`createContext(...).listTreeFiles(...)`), with an injected runner that records the
    // exact git argv AND the stdin payload - so reverting to the per-file shape (an argv
    // count of N with no `--stdin`) or emitting MORE than one batch spawn turns it red.
    // The repeat-call clause pins memoisation (treeCache/ignoreCache): a second call for the
    // same dir must not spawn another batch. The set clause ties the output to the input the
    // shipped code itself produced minus the injected ignore response.
    const T60_DIR = 'tools/gates/testdata/historical';
    const t60Calls = [];
    let t60BatchInput = null;
    let t60Ignored = new Set();
    const t60Spy = (root, args, binary, input) => {
      const key = args.join(' ');
      t60Calls.push({ key, stdin: input === undefined ? null : String(input) });
      if (/^status --porcelain=v1/.test(key)) return { ok: true, code: 0, out: '', err: '', spawnError: null };
      if (/^check-ignore --stdin -z$/.test(key)) {
        t60BatchInput = String(input || '');
        // Deterministic sentinel: ignore exactly the lexicographically FIRST path in the
        // batch input, so the ignore branch is exercised with the real (batched) argv.
        const fresh = t60BatchInput.split('\0').filter(Boolean).sort();
        t60Ignored = new Set(fresh.slice(0, 1));
        return { ok: true, code: 0, out: [...t60Ignored].map((p) => p + '\0').join(''), err: '', spawnError: null };
      }
      return { ok: true, code: 0, out: '', err: '', spawnError: null };
    };
    const t60Ctx = createContext({ repoRoot: realRoot, scope: 'tree', spawnGit: t60Spy });
    const t60First = t60Ctx.listTreeFiles(T60_DIR);
    const t60BatchCalls = t60Calls.filter((c) => /^check-ignore --stdin -z$/.test(c.key));
    const t60PerFile = t60Calls.filter((c) => /^check-ignore -q/.test(c.key)).length;
    const t60LsTree = t60Calls.filter((c) => /^ls-tree/.test(c.key)).length;
    const t60Input = (t60BatchInput || '').split('\0').filter(Boolean);
    const t60Expected = t60Input.filter((p) => !t60Ignored.has(p)).sort();
    const t60Second = t60Ctx.listTreeFiles(T60_DIR);
    const t60BatchAfter = t60Calls.filter((c) => /^check-ignore --stdin -z$/.test(c.key)).length;
    check(
      'T60 DEF-2 regression anchor: the tree branch batches a whole directory into EXACTLY ONE `git check-ignore --stdin -z` spawn (never a per-file check-ignore), memoises per dir, and returns the disk set minus the ignored set',
      t60BatchCalls.length === 1 &&
        t60BatchCalls[0].key === 'check-ignore --stdin -z' &&
        t60BatchCalls[0].stdin !== null &&
        t60PerFile === 0 &&
        t60LsTree === 0 &&
        t60Input.length >= 1 &&
        t60Input.every((p) => p.startsWith(T60_DIR + '/')) &&
        t60Ignored.size === 1 &&
        JSON.stringify(t60First) === JSON.stringify(t60Expected) &&
        t60First.length === t60Input.length - t60Ignored.size &&
        t60First.indexOf([...t60Ignored][0]) < 0 &&
        t60BatchAfter === 1 &&
        JSON.stringify(t60Second) === JSON.stringify(t60First),
      JSON.stringify({
        batch: t60BatchCalls.length,
        batchAfter: t60BatchAfter,
        perFile: t60PerFile,
        lsTree: t60LsTree,
        input: t60Input.length,
        ignored: [...t60Ignored],
        first: t60First.length,
      })
    );
  }

  // ---- T61-T67: the two GLM members landed by slice DIRECTIVE-GLM-LANDING (v1.26) -----
  // The closed set grew by `glm-5.3` and `glm-5.3-flash`. T61 proves BOTH new members are
  // accepted; T62/T64 prove each one is load-bearing (dropping either turns its positive
  // fixture red) and T63/T65 prove the in-memory mutation leaves g5b.js byte-identical; T66
  // keeps the comparison case-sensitive for the GLM member and T67 pins the R2.2 boundary
  // (a qualified call string is NOT the `modelId`).
  {
    const glm = g5bProbe('g5b-closed-set-legal-glm');
    const glmFlash = g5bProbe('g5b-closed-set-legal-glm-flash');
    check(
      'T61 G5-b accepts both members this slice added (glm-5.3 pre-review and glm-5.3-flash batch assistance)',
      glm.verdict === 'PASS' && glm.findings === 0 && glmFlash.verdict === 'PASS' && glmFlash.findings === 0,
      JSON.stringify({ glm: glm.verdict, glmFlash: glmFlash.verdict, glmDetail: glm.detail, glmFlashDetail: glmFlash.detail })
    );

    const dropFlash = loadMutatedG5b(G5B_GLM_FLASH_MEMBER, '', 'drop-glm-flash-member');
    check(
      'T62 the glm-5.3-flash member is load-bearing: dropping it from the set makes that fixture fail',
      dropFlash.applied === true && g5bProbeWith(dropFlash.exports, 'g5b-closed-set-legal-glm-flash').verdict === 'FAIL',
      'applied=' + dropFlash.applied + ' patched_verdict=' + g5bProbeWith(dropFlash.exports, 'g5b-closed-set-legal-glm-flash').verdict
    );
    check(
      'T63 restore: g5b.js is byte-identical and the cached module still accepts glm-5.3-flash',
      dropFlash.identical === true && g5bProbe('g5b-closed-set-legal-glm-flash').verdict === 'PASS',
      'identical=' + dropFlash.identical + ' cached=' + g5bProbe('g5b-closed-set-legal-glm-flash').verdict
    );

    const dropGlm = loadMutatedG5b(G5B_GLM_MEMBER, '', 'drop-glm-member');
    check(
      'T64 the glm-5.3 member is load-bearing: dropping it from the set makes that fixture fail',
      dropGlm.applied === true && g5bProbeWith(dropGlm.exports, 'g5b-closed-set-legal-glm').verdict === 'FAIL',
      'applied=' + dropGlm.applied + ' patched_verdict=' + g5bProbeWith(dropGlm.exports, 'g5b-closed-set-legal-glm').verdict
    );
    check(
      'T65 restore: g5b.js is byte-identical and the cached module still accepts glm-5.3',
      dropGlm.identical === true && g5bProbe('g5b-closed-set-legal-glm').verdict === 'PASS',
      'identical=' + dropGlm.identical + ' cached=' + g5bProbe('g5b-closed-set-legal-glm').verdict
    );

    const glmCase = g5bProbe('g5b-illegal-glm-case');
    check(
      'T66 G5-b stays case-sensitive for the GLM member: GLM-5.3 (wrong case) fails',
      glmCase.verdict === 'FAIL' && glmCase.findings === 1 && /model=GLM-5\.3/.test(glmCase.detail),
      'verdict=' + glmCase.verdict + ' findings=' + glmCase.findings + ' :: ' + glmCase.detail
    );

    const glmQualified = g5bProbe('g5b-illegal-glm-qualified');
    check(
      'T67 G5-b rejects a qualified call string: glm-5.3 (glm) is not the modelId (R2.2 boundary)',
      glmQualified.verdict === 'FAIL' && glmQualified.findings === 1 && /model=glm-5\.3 \(glm\)/.test(glmQualified.detail),
      'verdict=' + glmQualified.verdict + ' findings=' + glmQualified.findings + ' :: ' + glmQualified.detail
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

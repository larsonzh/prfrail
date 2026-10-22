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
const G3A = require('./lib/criteria/g3a');
const G6 = require('./lib/criteria/g6');
const G1 = require('./lib/criteria/g1');
const G5 = require('./lib/criteria/g5');
const G7 = require('./lib/criteria/g7');
const G8 = require('./lib/criteria/g8');
const G9 = require('./lib/criteria/g9');
const G4B = require('./lib/criteria/g4b');
const G5A = require('./lib/criteria/g5a');
const G5B = require('./lib/criteria/g5b');
const G5C = require('./lib/criteria/g5c');
const G5D = require('./lib/criteria/g5d');

/**
 * gate.js assigns `module.exports` at the END of its body, so a top-level require here
 * returns the pre-assignment object whenever gate.js is the entry point (`--selftest`).
 * Resolve it at call time instead, after the first await inside runTests().
 */
function gateApi() {
  return require('./gate');
}

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

/** G7-c / G7-d / G7-e verdict + finding count + first `why` for one synthetic report body. */
function g7Probe(id, text, mod) {
  const out = createReport();
  const M = mod || G7;
  const sub = M.subchecks.find((s) => s.id === id);
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
/**
 * OB-47 (v1.33): every probe body declares its language in the header, exactly like the
 * real documents. The probes below pin the three fail-closed modes (absent / outside the
 * closed set / name-declaration disagreement) in addition to the unchanged word-table
 * switch they already asserted.
 */
const LANG_TAG_ZH = '`lang:zh`\n';
const LANG_TAG_EN = '`lang:en`\n';
/** Mutation target ②: `labelsFromLines`'s bold line-leading label collector. */
const BOLD_LABEL_COLLECTOR = '    m = LABEL_BOLD_RE.exec(l);\n    if (m) out.add(m[1]);\n';
/** Mutation target ③: the CN closed set's first entry. */
const CN_FIRST_STATUS_WORD = "  cn: ['观察中', ";
const G6_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g6.js');
const G3A_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g3a.js');
const G4A_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g4a.js');
const G7_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g7.js');
const G5C_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g5c.js');
const G5D_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g5d.js');
/** Mutation target (T73): the G5-c "value is a member of 2.2" test (OB-29). */
const G5C_MEMBERSHIP = '          if (allowed.has(norm(m[1]))) continue;\n';
/** Mutation target (T72): the G7-e closed-set membership test (OB-42). */
const G7E_MEMBERSHIP = '      if (v && VERDICT_PREFIXES[prefix].indexOf(v) >= 0) return null;\n';
/** Mutation target (T28c): the sixth mirror shape added by OB-73. */
const G3A_BACKTICK_SHAPE = '        backticks: count(CL, /`/g) === count(EL, /`/g),\n';
/** Mutation targets (T74 / OB-90, v1.34): the whole either-side pairing guard, and the mirror-existence guard inside it. */
const G3A_OLD_PAIR_GUARD = "      if (enSide || !changed.includes(enf)) {\n        continue;\n      }\n";
const G3A_PAIR_GUARD = "      if (\n        !enSide &&\n        !changed.includes(enf) &&\n        !(typeof ctx.exists === 'function' && ctx.exists(enf))\n      ) {\n        continue;\n      }\n";
const G3A_MIRROR_GUARD = "        !(typeof ctx.exists === 'function' && ctx.exists(enf))\n";
/** Mutation target (T75 / OB-91, v1.34; needle updated by OB-97, v1.40): the tag-token
 * continuation charset. */
const G6_LANG_TAG_CONT = 'const LANG_TAG_CONT_RE = /[A-Za-z0-9._:-]/;';
/** Mutation target (T76, v1.36): the G1(f) in-progress wording test inside the completion log. */
const G1_INPROGRESS_TEST = '      if (IN_PROGRESS_RE.test(c.summary)) badF.push(';
/** Mutation target (T77, v1.36): the G5(d) change-set vs declared-set verdict. */
const G5_SET_EQUALITY = "      missing.length === 0 ? 'PASS' : 'FAIL',";

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
function loadMutatedModule(source, from, to, tag) {
  const before = fs.readFileSync(source);
  const src = decodeBuffer(before);
  const patched = src.split(from).join(to);
  MUTATION_LOG.push({ key: path.basename(source) + ':' + tag, op: 'sourcePatch', needle: from, applied: patched !== src });
  const M = require('module');
  const mod = new M(source, null);
  mod.filename = source;
  mod.paths = M._nodeModulePaths(path.dirname(source));
  mod._compile(patched, source);
  return { exports: mod.exports, applied: patched !== src, identical: fs.readFileSync(source).equals(before) };
}

function loadMutatedG6(from, to, tag) {
  return loadMutatedModule(G6_SOURCE, from, to, tag);
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
function g66Probe(virtualPath, text, mod) {
  const lines = String(text).split('\n');
  const ctx = {
    repoRoot: discoverRepoRoot(__dirname),
    changedMarkdown: () => [virtualPath],
    readLineArray: (q) => (q === virtualPath ? lines : null),
    addedLineNumbers: () => new Set(lines.map((_, i) => i + 1)),
    loadWhitelist: () => ({ entries: [], error: null, source: 'probe' }),
  };
  const out = createReport();
  const res = (mod || G6).subchecks.find((s) => s.id === 'G6-6').run(ctx, out);
  return { verdict: verdictOf(out.records, 'G6-6'), findings: res.findings.length, detail: detailOf(out.records, 'G6-6') };
}

/**
 * G6-7 probe (OB-32): feed a synthetic document to the sub-check and read back the
 * verdict / finding count / record detail. `whitelist` exercises the per-hit exemption
 * through the shared `emit` helper; `mod` runs a patched copy for the mutation probes.
 */
function g67Probe(virtualPath, text, mod, opts) {
  opts = opts || {};
  const lines = String(text).split('\n');
  const added = opts.added || lines.map((_, i) => i + 1);
  const ctx = {
    repoRoot: discoverRepoRoot(__dirname),
    changedMarkdown: () => [virtualPath],
    readLineArray: (q) => (q === virtualPath ? lines : null),
    addedLineNumbers: () => new Set(added),
    loadWhitelist: () => ({ entries: opts.whitelist || [], error: null, source: 'probe' }),
  };
  const out = createReport();
  const res = (mod || G6).subchecks.find((s) => s.id === 'G6-7').run(ctx, out);
  return {
    verdict: verdictOf(out.records, 'G6-7'),
    findings: res.findings.length,
    detail: detailOf(out.records, 'G6-7'),
    records: out.records,
  };
}

/**
 * G5-d probe (OB-30): a synthetic change set of role files. `texts` maps a repository
 * path to the file body; only paths the criterion itself selects are judged.
 */
function g5dProbe(texts, opts) {
  opts = opts || {};
  const paths = Object.keys(texts);
  const lines = new Map(paths.map((p) => [p, String(texts[p]).split('\n')]));
  const deleted = new Set(opts.deleted || []);
  const ctx = {
    changedPaths: paths,
    recordFor: (p) => ({ path: p, deleted: deleted.has(p) }),
    readLineArray: (p) => (lines.has(p) ? lines.get(p) : null),
  };
  const out = createReport();
  const res = (opts.mod || G5D).run(ctx, out);
  return {
    verdict: verdictOf(out.records, 'G5-d'),
    findings: res.findings.length,
    detail: detailOf(out.records, 'G5-d'),
    records: out.records,
  };
}

/**
 * G1 probe (v1.36): feed a synthetic ledger body to the base criterion and report the four
 * sub-verdicts. Only the CN ledger path is readable, so every case is judged once.
 */
function g1Probe(text, mod, added) {
  const M = mod || G1;
  const cn = M._internals.LEDGERS[0];
  const ctx = {
    repoRoot: discoverRepoRoot(__dirname),
    readLineArray: (p) => (p === cn ? String(text).split('\n') : null),
    addedLineNumbers: added ? () => new Set(added) : undefined,
  };
  const out = createReport();
  M.run(ctx, out);
  const pick = (id) => {
    const r = out.records.find((x) => x.id === id);
    return r ? r.verdict : 'ABSENT';
  };
  return { c: pick('G1(c)'), d: pick('G1(d)'), e: pick('G1(e)'), f: pick('G1(f)'), detail: out.records.map((r) => r.id + '=' + r.verdict).join(' ') };
}

/**
 * G5 probe (v1.36): a synthetic scope with injected change set, report fence, added ledger
 * rows, index hygiene and artifact resolution.
 */
function g5Probe(o, mod) {
  const M = mod || G5;
  const changed = o.changed || [];
  const ctx = {
    repoRoot: discoverRepoRoot(__dirname),
    kind: o.kind || 'index',
    changedPaths: changed,
    readText: (p) => (o.texts && o.texts[p] !== undefined ? o.texts[p] : null),
    addedLines: (p) => (o.added && o.added[p]) || [],
    git: () => ({ ok: true, out: (o.unstaged || []).join('\n') }),
    repoArtifactResolvable: (p) => (o.existing || []).some((x) => p === x || p.indexOf(x.replace(/\/+$/, '') + '/') === 0),
    recordFor: (p) => {
      const gone = (o.deleted || []).indexOf(p) >= 0;
      // production shape from parseStatusZ: two-char XY plus the parser's `deleted` boolean
      return { status: gone ? 'D ' : 'M ', deleted: gone };
    },
  };
  const out = createReport();
  M.run(ctx, out);
  const pick = (id) => {
    const r = out.records.find((x) => x.id === id);
    return r ? r.verdict : 'ABSENT';
  };
  return { d: pick('G5(d)'), e: pick('G5(e)'), f: pick('G5(f)'), detail: out.records.map((r) => r.id + '=' + r.verdict).join(' '), raw: out.records.map((r) => r.id + ':' + r.detail).join(' | ') };
}

/** G3(a.2) mirror-parity probe (OB-73). A synthetic directive pair with only the four ctx
 * members g3a.js touches. `mod` lets the mutation probe run a patched copy of the
 * criterion instead of the shipped one.
 */
function g3aProbe(cnText, enText, mod) {
  const CNP = 'docs/DELIVERY_DIRECTIVE.md';
  const ENP = 'docs/DELIVERY_DIRECTIVE_EN.md';
  const texts = new Map([
    [CNP, String(cnText).split('\n')],
    [ENP, String(enText).split('\n')],
  ]);
  const ctx = {
    changedMarkdown: () => [CNP, ENP],
    changedPaths: [CNP, ENP],
    numstat: () => [1, 1],
    readLineArray: (f) => (texts.has(f) ? texts.get(f) : null),
  };
  const out = createReport();
  (mod || G3A).run(ctx, out);
  return {
    records: out.records,
    a: verdictOf(out.records, 'G3(a)'),
    a2: verdictOf(out.records, 'G3(a.2) ' + CNP),
    a2detail: detailOf(out.records, 'G3(a.2) ' + CNP),
    // No FAIL record exists when the shapes agree, so the shape JSON is read from the
    // aggregate record, which always carries it.
    detail: detailOf(out.records, 'G3(a)'),
  };
}

/**
 * G3(a) PAIRING probe (OB-90, v1.34). The stock `g3aProbe` hard-codes one two-sided
 * directive pair; this one drives the pairing preconditions themselves:
 *   opts.changed  - the change set (paths)
 *   opts.texts    - path -> line array; an absent path is unreadable (deleted)
 *   opts.numstat  - path -> [added, deleted] (default [0, 0])
 *   opts.existing - paths `ctx.exists()` reports as present at the scope endpoint / disk
 *   opts.key      - the canonical CN path whose record ids are read back
 *   opts.mod      - a patched criterion copy (mutation probes)
 */
function g3aEdgeProbe(opts) {
  const o = opts || {};
  const texts = new Map(Object.entries(o.texts || {}));
  const nums = o.numstat || {};
  const ctx = {
    changedMarkdown: () => o.changed.filter((p) => /\.md$/i.test(p)),
    changedPaths: o.changed,
    numstat: (f) => nums[f] || [0, 0],
    readLineArray: (f) => (texts.has(f) ? texts.get(f) : null),
    exists: (f) => (o.existing || []).indexOf(f) >= 0,
  };
  const out = createReport();
  (o.mod || G3A).run(ctx, out);
  return {
    records: out.records,
    a: verdictOf(out.records, 'G3(a)'),
    a1: verdictOf(out.records, 'G3(a.1) ' + o.key),
    a2: verdictOf(out.records, 'G3(a.2) ' + o.key),
    a2detail: detailOf(out.records, 'G3(a.2) ' + o.key),
    // No FAIL record exists when the shapes agree, so the shape JSON is read from the
    // aggregate record, which always carries it.
    detail: detailOf(out.records, 'G3(a)'),
  };
}

/* ---------------------------------------------------------------------------
 * G5-b (role-file `model` closed set) and G5-a 2 probes (slice DIRECTIVE-MODEL-LISTS,
 * v1.18). T29-T37 below are what makes the new criterion falsifiable: each boundary
 * (verbatim comparison, case sensitivity, quoting, multi-key collection, the empty-value
 * division of labour, and the load-bearing closed set) is asserted in turn.
 * ------------------------------------------------------------------------ */

/** Path of the criterion under test (the mutation probe compiles a patched copy). */
const G5B_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g5b.js');
const G9_SOURCE = path.join(__dirname, 'lib', 'criteria', 'g9.js');
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

/**
 * G5-c probe (OB-29): a synthetic directive pair. G5-c reads `ctx.readLineArray` /
 * `ctx.addedLineNumbers` only, and every line of the probe counts as added. `mod` lets the
 * mutation probe run a patched copy of the criterion.
 */
function g5cProbe(cnText, enText, mod) {
  const [CNP, ENP] = G5C.internals.DIRECTIVE_FILES;
  const texts = new Map([
    [CNP, String(cnText).split('\n')],
    [ENP, String(enText).split('\n')],
  ]);
  const ctx = {
    readLineArray: (f) => (texts.has(f) ? texts.get(f) : null),
    addedLineNumbers: (f) => new Set((texts.get(f) || []).map((_, i) => i + 1)),
  };
  const out = createReport();
  (mod || G5C).run(ctx, out);
  return { records: out.records, verdict: verdictOf(out.records, 'G5-c'), detail: detailOf(out.records, 'G5-c') };
}

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
 * G9 probe: feed a synthetic change set to the shipped criterion (or a patched copy).
 * `opts.mod` runs a mutation probe; `opts.deleted` / `opts.whitelist` exercise the two
 * exclusion faces. `isEvidenceTree` is imported from lib/ctx so the probe cannot drift
 * from the shipped single source.
 */
function g9Probe(paths, opts) {
  opts = opts || {};
  const deleted = new Set(opts.deleted || []);
  const texts = opts.texts || new Map();
  const ctx = {
    changedPaths: paths,
    recordFor: (p) => ({ path: p, deleted: deleted.has(p) }),
    isEvidenceTree: (p) => isEvidenceTree(p),
    loadWhitelist: () => ({ entries: opts.whitelist || [], error: null, source: 'probe' }),
    changedMarkdown: () => opts.markdown || [],
    readLineArray: (p) => (texts.has(p) ? String(texts.get(p)).split('\n') : null),
    addedLineNumbers: (p) => new Set(String(texts.get(p) || '').split('\n').map((_, i) => i + 1)),
    repoArtifactResolvable: (p) => (opts.resolve ? opts.resolve(p) : { ok: true, reason: '' }),
  };
  const out = createReport();
  const res = (opts.mod || G9).run(ctx, out) || { findings: [] };
  return {
    verdict: verdictOf(out.records, 'G9'),
    info: out.records.some((r) => r.id === 'G9' && r.verdict === 'INFO'),
    findings: res.findings ? res.findings.length : -1,
    detail: detailOf(out.records, 'G9'),
    infos: out.records.filter((r) => r.id === 'G9-e'),
  };
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

async function runTests() {
  // The suite is async from T90 on. This first await also lets gate.js finish assigning
  // `module.exports` before gateApi() is called (circular require on the --selftest path).
  await Promise.resolve();
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

  // ---- T72: G7-e verdict-line closed set (OB-42, v1.33).
  {
    const ok = g7Probe('G7-e', '# report\n\n**RE-REVIEW**: PASS（无 Medium+）\n');
    check(
      'T72 OB-42 positive control: a closed-set verdict line passes (bold marker + parenthetical note allowed)',
      ok.verdict === 'PASS' && ok.findings === 0,
      'verdict=' + ok.verdict + ' findings=' + ok.findings + ' :: ' + ok.why
    );
    const outOfSet = g7Probe('G7-e', '# report\n\nFINAL REVIEW: MAYBE\n');
    check(
      'T72 OB-42 an out-of-set value under a closed-set prefix is a format deviation',
      outOfSet.verdict === 'FAIL' && outOfSet.findings === 1 && /outside the closed set/.test(outOfSet.why),
      'verdict=' + outOfSet.verdict + ' :: ' + outOfSet.why
    );
    const wrongPrefix = g7Probe('G7-e', '# report\n\nTEST: PASS\n');
    check(
      'T72 OB-42 a value legal for ANOTHER prefix is still a deviation (per-prefix sets, not one merged set)',
      wrongPrefix.verdict === 'FAIL' && /outside the closed set/.test(wrongPrefix.why),
      'verdict=' + wrongPrefix.verdict + ' :: ' + wrongPrefix.why
    );
    const bare = g7Probe('G7-e', '# report\n\n**PASS**\n');
    check(
      'T72 OB-42 the bare-token rule: a prefix-less `PASS` is a deviation',
      bare.verdict === 'FAIL' && /bare verdict token/.test(bare.why),
      'verdict=' + bare.verdict + ' :: ' + bare.why
    );
    const external = g7Probe('G7-e', '# report\n\nVERDICT: FINDINGS\n');
    check(
      'T72 OB-42 the table-external rule: an unknown prefix carrying exactly one closed value is a deviation',
      external.verdict === 'FAIL' && /table-external prefix/.test(external.why),
      'verdict=' + external.verdict + ' :: ' + external.why
    );
    const reResidual = g7Probe('G7-e', '# report\n\nRE_REVIEW: PASS\n');
    check(
      'T72 OB-42 a residual `RE_` placeholder is a deviation (4.1 abolished it)',
      reResidual.verdict === 'FAIL' && /residual `RE_` placeholder/.test(reResidual.why),
      'verdict=' + reResidual.verdict + ' :: ' + reResidual.why
    );
    const gateReading = g7Probe('G7-e', '# report\n\nSELFTEST: PASS 362/362\nGATE REPORT scope=tree\n');
    check(
      'T72 OB-42 stated boundary: a gate reading that only LOOKS like a verdict line is not judged',
      gateReading.verdict === 'PASS' && gateReading.findings === 0,
      'verdict=' + gateReading.verdict + ' findings=' + gateReading.findings + ' :: ' + gateReading.why
    );
    const quoted = g7Probe('G7-e', '# report\n\n示例行（以下形态出现即 FAIL）：FINAL REVIEW: MAYBE\n');
    check(
      'T72 OB-42 self-reference guard: an illegal form quoted on a marked line is skipped',
      quoted.verdict === 'PASS' && quoted.findings === 0,
      'verdict=' + quoted.verdict + ' findings=' + quoted.findings + ' :: ' + quoted.why
    );
    const mut = loadMutatedModule(G7_SOURCE, G7E_MEMBERSHIP, '      return null;\n', 'ob42-membership');
    check(
      'T72 OB-42 mutation applied: the membership line matched g7.js (a stale patch would weaken the suite silently)',
      mut.applied === true,
      'applied=' + mut.applied
    );
    const blind = g7Probe('G7-e', '# report\n\nFINAL REVIEW: MAYBE\n', mut.exports);
    check(
      'T72 OB-42 load-bearing proof: accepting any recognisable value makes the SAME out-of-set line pass',
      blind.verdict === 'PASS',
      'verdict=' + blind.verdict + ' :: ' + blind.why
    );
    check(
      'T72 restore: g7.js is byte-identical after the mutation window and the cached module still fails the line',
      mut.identical === true && g7Probe('G7-e', '# report\n\nFINAL REVIEW: MAYBE\n').verdict === 'FAIL',
      'identical=' + mut.identical
    );
  }

  // ---- T73: G5-c appendix-B call strings against the 2.2 qualified-name column (OB-29).
  {
    const sec22 = (rows) => '### 2.2 x\n\n| a | b | c | 限定名 | d |\n|---|---|---|---|---|\n' + rows + '\n';
    const apxB = (calls) => '## 附录 B\n\n```\n' + calls + '\n```\n';
    const cnGood = sec22('| r | ① | `x` | `DeepSeek V4 Pro (deepseek)` | `deepseek-v4-pro` |\n');
    const enGood = sec22('| r | ① | `x` | `GLM-5.3 (glm)` | `glm-5.3` |\n');
    const ok = g5cProbe(cnGood + apxB("载体：runSubagent(agentName='a', model='DeepSeek V4 Pro (deepseek)')"), enGood);
    check(
      'T73 OB-29 positive control: a call string matching the 2.2 qualified name passes (code span vs plain text is normalised away)',
      ok.verdict === 'PASS',
      'verdict=' + ok.verdict + ' :: ' + ok.detail
    );
    const bad = g5cProbe(cnGood + apxB("载体：runSubagent(agentName='a', model='GPT-9-Turbo (copilot)')"), enGood);
    check(
      'T73 OB-29 a call string outside the 2.2 qualified-name column fails',
      bad.verdict === 'FAIL' && /is not in the 2.2 qualified-name column/.test(bad.detail),
      'verdict=' + bad.verdict + ' :: ' + bad.detail
    );
    const enJudged = g5cProbe(cnGood, enGood + apxB("carrier: model='GLM-5.3 (glm)'"));
    check(
      'T73 OB-29 the `_EN` side is judged with its own lowercase table header',
      enJudged.verdict === 'PASS',
      'verdict=' + enJudged.verdict + ' :: ' + enJudged.detail
    );
    const noApx = g5cProbe(cnGood, enGood);
    check(
      'T73 OB-29 a directive without an appendix B region is not judged',
      noApx.verdict === 'PASS' && /no directive with an appendix B region/.test(noApx.detail),
      'verdict=' + noApx.verdict + ' :: ' + noApx.detail
    );
    let premErr = null;
    try {
      g5cProbe("## 附录 B\n\n```\nmodel='DeepSeek V4 Pro (deepseek)'\n```\n", '');
    } catch (e) {
      premErr = e;
    }
    check(
      'T73 OB-29 fail-closed: an appendix B whose 2.2 qualified-name column cannot be read raises UsageError',
      !!premErr && premErr instanceof UsageError && /no 2.2 section|qualified-name column was not found/.test(premErr.message),
      'err=' + (premErr ? premErr.message : 'NO ERROR')
    );
    const mut = loadMutatedModule(G5C_SOURCE, G5C_MEMBERSHIP, '          continue;\n', 'ob29-membership');
    check(
      'T73 OB-29 mutation applied: the 2.2 membership line matched g5c.js (a stale patch would weaken the suite silently)',
      mut.applied === true,
      'applied=' + mut.applied
    );
    const blind = g5cProbe(
      cnGood + apxB("载体：runSubagent(agentName='a', model='GPT-9-Turbo (copilot)')"),
      enGood,
      mut.exports
    );
    check(
      'T73 OB-29 load-bearing proof: with the membership test patched out the SAME out-of-set call string passes',
      blind.verdict === 'PASS',
      'verdict=' + blind.verdict + ' :: ' + blind.detail
    );
    check(
      'T73 restore: g5c.js is byte-identical after the mutation window and the cached module still fails the call string',
      mut.identical === true &&
        g5cProbe(cnGood + apxB("载体：runSubagent(agentName='a', model='GPT-9-Turbo (copilot)')"), enGood).verdict === 'FAIL',
      'identical=' + mut.identical
    );
  }

  // ---- T74: OB-90 (v1.34) - G3(a) pairs from EITHER side, so a one-sided edit is judged.
  {
    const CNP = 'docs/t027/probe-pair-cn.md';
    const ENP = 'docs/t027/probe-pair-cn_EN.md';
    const texts = { [CNP]: ['# c', ''], [ENP]: ['# e', ''] };
    const enOnly = g3aEdgeProbe({
      key: CNP,
      changed: [ENP],
      texts,
      numstat: { [CNP]: [0, 0], [ENP]: [3, 0] },
      existing: [],
    });
    check(
      'T74 OB-90 an `_EN.md`-only change is now judged instead of silently ignored, and fails the symmetry test',
      enOnly.a1 === 'FAIL' && enOnly.a === 'FAIL' && enOnly.a2 === 'ABSENT' && /_EN \+3\/-0 MISMATCH/.test(enOnly.detail),
      'a1=' + enOnly.a1 + ' a=' + enOnly.a + ' :: ' + enOnly.detail
    );
    const cnOnly = g3aEdgeProbe({
      key: CNP,
      changed: [CNP],
      texts,
      numstat: { [CNP]: [3, 0], [ENP]: [0, 0] },
      existing: [ENP],
    });
    check(
      'T74 OB-90 a CN-only change whose `_EN` mirror exists is judged too (the mirror need not be in the change set)',
      cnOnly.a1 === 'FAIL' && cnOnly.a === 'FAIL' && cnOnly.a2 === 'ABSENT' && /vs _EN \+0\/-0 MISMATCH/.test(cnOnly.detail),
      'a1=' + cnOnly.a1 + ' :: ' + cnOnly.detail
    );
    const DIRC = 'docs/DELIVERY_DIRECTIVE.md';
    const DIREN = 'docs/DELIVERY_DIRECTIVE_EN.md';
    const dirDrift = g3aEdgeProbe({
      key: DIRC,
      changed: [DIREN],
      texts: { [DIRC]: ['# c', ''], [DIREN]: ['# e', '', '# e2', ''] },
      numstat: { [DIRC]: [0, 0], [DIREN]: [3, 0] },
      existing: [DIRC],
    });
    check(
      'T74 OB-90 directive pair, EN-only edit with a shape drift: G3(a.1) and G3(a.2) both go red',
      dirDrift.a1 === 'FAIL' && dirDrift.a2 === 'FAIL' && /"lines":false/.test(dirDrift.a2detail),
      'a1=' + dirDrift.a1 + ' a2=' + dirDrift.a2 + ' :: ' + dirDrift.a2detail
    );
    const dirNeutral = g3aEdgeProbe({
      key: DIRC,
      changed: [DIREN],
      texts: { [DIRC]: ['# c', ''], [DIREN]: ['# e', ''] },
      numstat: { [DIRC]: [0, 0], [DIREN]: [2, 0] },
      existing: [DIRC],
    });
    check(
      'T74 OB-90 directive pair, EN-only edit with a shape-neutral body: only G3(a.1) goes red, G3(a.2) stays silent',
      dirNeutral.a1 === 'FAIL' && dirNeutral.a2 === 'ABSENT',
      'a1=' + dirNeutral.a1 + ' a2=' + dirNeutral.a2
    );
    const lowerEn = g3aEdgeProbe({
      key: CNP,
      changed: ['docs/t027/probe-pair-cn_en.md'],
      texts: { [CNP]: ['# c', ''] },
      numstat: { 'docs/t027/probe-pair-cn_en.md': [3, 0] },
      existing: [],
    });
    check(
      'T74 OB-90 a non-canonical \`_en.md\` spelling is judged as the EN side of the same pair, not as two untouched paths (5 r1 F4)',
      lowerEn.a1 === 'FAIL' && /_EN \+3\/-0/.test(lowerEn.detail),
      'a1=' + lowerEn.a1 + ' :: ' + lowerEn.detail
    );
    const single = g3aEdgeProbe({
      key: 'docs/README.md',
      changed: ['docs/README.md'],
      texts: { 'docs/README.md': ['# r', ''] },
      numstat: { 'docs/README.md': [3, 0] },
      existing: [],
    });
    check(
      'T74 OB-90 load-bearing guard: a single-language `.md` without an `_EN` mirror stays out of scope, so the mirror-less corpus docs cannot go red',
      single.a === 'ABSENT' && single.records.length === 0,
      'a=' + single.a + ' records=' + single.records.length
    );
    const both = g3aEdgeProbe({
      key: CNP,
      changed: [CNP, ENP],
      texts,
      numstat: { [CNP]: [3, 0], [ENP]: [3, 0] },
      existing: [ENP],
    });
    check(
      'T74 OB-90 no behaviour drift: a two-sided equal edit still yields one green aggregate and no per-pair FAIL',
      both.a === 'PASS' && both.a1 === 'ABSENT' && both.records.length === 1,
      'a=' + both.a + ' records=' + both.records.length
    );
    const deletedBoth = g3aEdgeProbe({
      key: CNP,
      changed: [CNP, ENP],
      texts: {},
      numstat: { [CNP]: [0, 3], [ENP]: [0, 3] },
      existing: [],
    });
    check(
      'T74 OB-90 equal deletion on both sides still passes (an unreadable pair no longer skips the numstat verdict)',
      deletedBoth.a === 'PASS' && deletedBoth.a1 === 'ABSENT',
      'a=' + deletedBoth.a + ' :: ' + deletedBoth.detail
    );
    const deletedEn = g3aEdgeProbe({
      key: CNP,
      changed: [ENP],
      texts: { [CNP]: ['# c', ''] },
      numstat: { [CNP]: [0, 0], [ENP]: [0, 3] },
      existing: [],
    });
    check(
      'T74 OB-90 a deleted `_EN` side still yields a G3(a.1) FAIL (the old code skipped unreadable pairs entirely)',
      deletedEn.a1 === 'FAIL' && deletedEn.a === 'FAIL',
      'a1=' + deletedEn.a1 + ' :: ' + deletedEn.detail
    );
    const dirGone = g3aEdgeProbe({
      key: 'docs/DELIVERY_DIRECTIVE.md',
      changed: ['docs/DELIVERY_DIRECTIVE_EN.md'],
      texts: { 'docs/DELIVERY_DIRECTIVE_EN.md': ['# e', ''] },
      numstat: { 'docs/DELIVERY_DIRECTIVE.md': [0, 0], 'docs/DELIVERY_DIRECTIVE_EN.md': [3, 0] },
      existing: [],
    });
    check(
      'T74 OB-90 the directive pair with an unreadable side reports G3(a.2) FAIL (premise broken, fail closed)',
      dirGone.a2 === 'FAIL' && /side unreadable/.test(dirGone.a2detail),
      'a2=' + dirGone.a2 + ' :: ' + dirGone.a2detail
    );
    const mut = loadMutatedModule(G3A_SOURCE, G3A_PAIR_GUARD, G3A_OLD_PAIR_GUARD, 'ob90-enside');
    check(
      'T74 OB-90 mutation applied: the pairing-guard block matched g3a.js (a stale patch would weaken the suite silently)',
      mut.applied === true,
      'applied=' + mut.applied
    );
    const blind = g3aEdgeProbe({
      key: CNP,
      changed: [ENP],
      texts,
      numstat: { [CNP]: [0, 0], [ENP]: [3, 0] },
      existing: [],
      mod: mut.exports,
    });
    check(
      'T74 OB-90 load-bearing proof: with the either-side guard reverted the SAME `_EN`-only change produces NO record again',
      blind.a === 'ABSENT' && blind.records.length === 0,
      'a=' + blind.a + ' records=' + blind.records.length
    );
    check(
      'T74 restore: g3a.js is byte-identical after the mutation window and the cached module still judges the `_EN`-only change',
      mut.identical === true &&
        g3aEdgeProbe({
          key: CNP,
          changed: [ENP],
          texts,
          numstat: { [CNP]: [0, 0], [ENP]: [3, 0] },
          existing: [],
        }).a1 === 'FAIL',
      'identical=' + mut.identical
    );
    const mut2 = loadMutatedModule(G3A_SOURCE, G3A_MIRROR_GUARD, '        false\n', 'ob90-mirror-guard');
    check(
      'T74 OB-90 mutation applied (mirror-existence guard): the guard line matched g3a.js',
      mut2.applied === true,
      'applied=' + mut2.applied
    );
    const loud = g3aEdgeProbe({
      key: 'docs/README.md',
      changed: ['docs/README.md'],
      texts: { 'docs/README.md': ['# r', ''] },
      numstat: { 'docs/README.md': [3, 0] },
      existing: [],
      mod: mut2.exports,
    });
    check(
      'T74 OB-90 load-bearing proof (mirror-existence guard): without the predicate the SAME single-language document goes red',
      loud.a === 'FAIL',
      'a=' + loud.a + ' :: ' + loud.detail
    );
    check(
      'T74 restore (mirror-existence guard): g3a.js is byte-identical and the single-language document is silent again',
      mut2.identical === true &&
        g3aEdgeProbe({
          key: 'docs/README.md',
          changed: ['docs/README.md'],
          texts: { 'docs/README.md': ['# r', ''] },
          numstat: { 'docs/README.md': [3, 0] },
          existing: [],
        }).a === 'ABSENT',
      'identical=' + mut2.identical
    );
  }

  // ---- T75: OB-91 (v1.34) - the language tag token is consumed to its FULL extent.
  {
    const row = H124 + '| OB-1 | d | x | 观察中 |\n';
    const errOf = (p, t, mod) => {
      try {
        g66Probe(p, t, mod);
        return null;
      } catch (e) {
        return e;
      }
    };
    const unknown = (e) => !!e && e instanceof UsageError && /unknown language tag/.test(e.message);
    const underscore = errOf('docs/t027/probe-g6-6-underscore.md', '`lang:zh_cn`\n' + row);
    check(
      'T75 OB-91 the 7 r3 residual is closed: `lang:zh_cn` is consumed whole and rejected (it used to truncate to a legal `lang:zh`)',
      unknown(underscore) && /`lang:zh_cn`/.test(underscore.message),
      'err=' + (underscore ? underscore.message : 'NO ERROR')
    );
    const underscoreEn = errOf('docs/t027/probe-g6-6-underscore_EN.md', '`lang:en_us`\n' + row);
    check(
      'T75 OB-91 the same holds for `lang:en_us` on the `_EN` side',
      unknown(underscoreEn) && /`lang:en_us`/.test(underscoreEn.message),
      'err=' + (underscoreEn ? underscoreEn.message : 'NO ERROR')
    );
    const bare = errOf('docs/t027/probe-g6-6-bare.md', '`lang:`\n' + row);
    check(
      'T75 OB-91 a bare `lang:` is an unknown tag rather than an invisible token',
      unknown(bare) && /`lang:`/.test(bare.message),
      'err=' + (bare ? bare.message : 'NO ERROR')
    );
    const bareCoexist = errOf('docs/t027/probe-g6-6-bare-coexist.md', LANG_TAG_ZH + '`lang:`\n' + row);
    check(
      'T75 OB-91 coexistence is fail-closed for the new form as well (`lang:zh` plus a bare `lang:`)',
      unknown(bareCoexist),
      'err=' + (bareCoexist ? bareCoexist.message : 'NO ERROR')
    );
    const underscoreCoexist = errOf(
      'docs/t027/probe-g6-6-underscore-coexist.md',
      LANG_TAG_ZH + '`lang:zh_cn`\n' + row
    );
    check(
      'T75 OB-91 coexistence (`lang:zh` plus `lang:zh_cn`) is rejected too',
      unknown(underscoreCoexist),
      'err=' + (underscoreCoexist ? underscoreCoexist.message : 'NO ERROR')
    );
    const trailingDash = errOf('docs/t027/probe-g6-6-dash.md', '`lang:zh-`\n' + row);
    check(
      'T75 OB-91 a trailing continuation character (`lang:zh-`) is rejected',
      unknown(trailingDash),
      'err=' + (trailingDash ? trailingDash.message : 'NO ERROR')
    );
    const trailingDot = errOf('docs/t027/probe-g6-6-dot.md', '`lang:zh.`\n' + row);
    check(
      'T75 OB-91 a trailing dot (`lang:zh.`) belongs to the token and is rejected (the terminator set is the complement of the continuation set)',
      unknown(trailingDot),
      'err=' + (trailingDot ? trailingDot.message : 'NO ERROR')
    );
    const upper = errOf('docs/t027/probe-g6-6-upper.md', '`lang:ZH`\n' + row);
    check(
      'T75 OB-91 the comparison stays whole-token and case-sensitive (`lang:ZH`)',
      unknown(upper),
      'err=' + (upper ? upper.message : 'NO ERROR')
    );
    const digit = errOf('docs/t027/probe-g6-6-digit.md', '`lang:zh2`\n' + row);
    check(
      'T75 OB-91 a digit suffix (`lang:zh2`) is rejected',
      unknown(digit),
      'err=' + (digit ? digit.message : 'NO ERROR')
    );
    const comma = g66Probe('docs/t027/probe-g6-6-comma.md', '`lang:zh,`\n' + row);
    check(
      'T75 OB-91 a comma terminates the token, so `lang:zh,` stays legal',
      comma.verdict === 'PASS',
      'verdict=' + comma.verdict + ' :: ' + comma.detail
    );
    const prose = g66Probe('docs/t027/probe-g6-6-prose.md', '`lang:zh` zh\n' + row);
    check(
      'T75 OB-91 prose after the tag creates no phantom occurrence (a tag followed by a space and more words)',
      prose.verdict === 'PASS',
      'verdict=' + prose.verdict + ' :: ' + prose.detail
    );
    const leadingDot = g66Probe('docs/t027/probe-g6-6-leadingdot.md', '.lang:zh\n' + row);
    check(
      'T75 OB-91 a leading dot is not part of the token (`.lang:zh` still declares zh)',
      leadingDot.verdict === 'PASS',
      'verdict=' + leadingDot.verdict + ' :: ' + leadingDot.detail
    );
    const wordChar = errOf('docs/t027/probe-g6-6-wordchar.md', 'xlang:zh\n' + row);
    check(
      'T75 OB-91 an occurrence preceded by a word character is not an occurrence (`xlang:zh` declares nothing and fails as missing)',
      !!wordChar && wordChar instanceof UsageError && /language tag missing/.test(wordChar.message),
      'err=' + (wordChar ? wordChar.message : 'NO ERROR')
    );
    const mut = loadMutatedModule(G6_SOURCE, G6_LANG_TAG_CONT, 'const LANG_TAG_CONT_RE = /[A-Za-z.-]/;', 'ob91-continuation');
    check(
      'T75 OB-91 mutation applied: the continuation-charset line matched g6.js (a stale patch would weaken the suite silently)',
      mut.applied === true,
      'applied=' + mut.applied
    );
    const blindUnderscore = g66Probe('docs/t027/probe-g6-6-underscore.md', '`lang:zh_cn`\n' + row, mut.exports);
    check(
      'T75 OB-91 load-bearing proof: with the continuation charset narrowed back the SAME `lang:zh_cn` header passes again (the charset, not the closed set, is what catches it)',
      blindUnderscore.verdict === 'PASS',
      'verdict=' + blindUnderscore.verdict + ' :: ' + blindUnderscore.detail
    );
    check(
      'T75 restore: g6.js is byte-identical after the mutation window and the cached module still rejects `lang:zh_cn`',
      mut.identical === true && unknown(errOf('docs/t027/probe-g6-6-underscore.md', '`lang:zh_cn`\n' + row)),
      'identical=' + mut.identical
    );
  }

  // ---- T78: OB-90 (v1.34) - the G3(a) pairing preconditions and the six whole-file shapes,
  // driven from the `_EN` side. T74 proves the either-side guard EXISTS; this block is the
  // counter-example inventory that turns each half red if it regresses: the six shape flags
  // are isolated one by one (numstat held equal, so only G3(a.2) can speak), the two control
  // shapes stay green, and the single-side add / delete / rename cases state the expected
  // verdict instead of accepting "no record at all".
  {
    const CNP = 'docs/DELIVERY_DIRECTIVE.md';
    const ENP = 'docs/DELIVERY_DIRECTIVE_EN.md';
    /** One directive-pair run: `nums` defaults to an equal numstat, so only the shape speaks. */
    const pair = (cn, en, nums, changed) =>
      g3aEdgeProbe({
        key: CNP,
        changed: changed || [ENP],
        texts: { [CNP]: cn, [ENP]: en },
        numstat: nums || { [CNP]: [1, 1], [ENP]: [1, 1] },
        existing: [CNP],
      });

    const bold = pair(['# c', '', '**b**', ''], ['# e', '', '**b** **x**', ''], { [CNP]: [0, 0], [ENP]: [1, 1] });
    check(
      'T78 OB-90 an `_EN`-only bold drift is judged from the EN side: numstat and the bold shape both go red',
      bold.a1 === 'FAIL' &&
        bold.a2 === 'FAIL' &&
        /"bold":false/.test(bold.a2detail) &&
        /_EN \+1\/-1 MISMATCH/.test(bold.detail),
      'a1=' + bold.a1 + ' a2=' + bold.a2 + ' :: ' + bold.a2detail
    );

    const ticks = pair(['# c', '', 'x', ''], ['# e', '', '`x`', '']);
    check(
      'T78 OB-90 the backtick shape is judged on its own (equal numstat, so only G3(a.2) can speak)',
      ticks.a1 === 'ABSENT' && ticks.a2 === 'FAIL' && /"backticks":false/.test(ticks.a2detail),
      'a1=' + ticks.a1 + ' a2=' + ticks.a2 + ' :: ' + ticks.a2detail
    );

    const head = pair(['# c', '', 'x', ''], ['# e', '## shift', '', 'x', '']);
    check(
      'T78 OB-90 a heading downshift on the `_EN` side (one inserted line) fails G3(a.2) with `head` false',
      head.a1 === 'ABSENT' &&
        head.a2 === 'FAIL' &&
        /"head":false/.test(head.a2detail) &&
        /"lines":false/.test(head.a2detail),
      'a1=' + head.a1 + ' a2=' + head.a2 + ' :: ' + head.a2detail
    );

    const bars = pair(['# c', '| a |', ''], ['# e', '| a | b |', '']);
    check(
      'T78 OB-90 the pipe shape is judged on its own (a table row widened on the `_EN` side)',
      bars.a1 === 'ABSENT' && bars.a2 === 'FAIL' && /"bars":false/.test(bars.a2detail),
      'a1=' + bars.a1 + ' a2=' + bars.a2 + ' :: ' + bars.a2detail
    );

    const blank = pair(['# c', '', 'x', ''], ['# e', 'z', 'y', '']);
    check(
      'T78 OB-90 the blank-line shape is judged on its own (same line count, one blank line filled in)',
      blank.a1 === 'ABSENT' && blank.a2 === 'FAIL' && /"blank":false/.test(blank.a2detail),
      'a1=' + blank.a1 + ' a2=' + blank.a2 + ' :: ' + blank.a2detail
    );

    const lines = pair(['# c', '', 'x', ''], ['# e', '', 'x', '', 'tail']);
    check(
      'T78 OB-90 the line-count shape is judged on its own (a trailing body line, heading positions intact)',
      lines.a1 === 'ABSENT' &&
        lines.a2 === 'FAIL' &&
        /"lines":false/.test(lines.a2detail) &&
        /"head":true/.test(lines.a2detail),
      'a1=' + lines.a1 + ' a2=' + lines.a2 + ' :: ' + lines.a2detail
    );

    const both = g3aEdgeProbe({
      key: CNP,
      changed: [CNP, ENP],
      texts: { [CNP]: ['# c', ''], [ENP]: ['# e', ''] },
      numstat: { [CNP]: [2, 0], [ENP]: [2, 0] },
      existing: [ENP],
    });
    check(
      'T78 OB-90 control: an equal two-sided edit stays green (one aggregate PASS, no per-pair FAIL)',
      both.a === 'PASS' && both.a1 === 'ABSENT' && both.a2 === 'ABSENT' && both.records.length === 1,
      'a=' + both.a + ' a1=' + both.a1 + ' a2=' + both.a2 + ' records=' + both.records.length
    );

    const bothShape = g3aEdgeProbe({
      key: CNP,
      changed: [CNP, ENP],
      texts: { [CNP]: ['# c', '**x**', ''], [ENP]: ['# e', '**x**', ''] },
      numstat: { [CNP]: [2, 0], [ENP]: [2, 0] },
      existing: [ENP],
    });
    check(
      'T78 OB-90 control: an identical shape change on BOTH sides stays green (the six shapes compare sides, not an absolute form)',
      bothShape.a === 'PASS' && bothShape.a1 === 'ABSENT' && bothShape.a2 === 'ABSENT',
      'a=' + bothShape.a + ' a1=' + bothShape.a1 + ' a2=' + bothShape.a2
    );

    const pairingOnly = pair(['# c', ''], ['# e', ''], { [CNP]: [0, 0], [ENP]: [2, 2] });
    check(
      'T78 OB-90 pairing without shape drift: an `_EN`-only modification fails G3(a.1) and leaves G3(a.2) silent',
      pairingOnly.a1 === 'FAIL' &&
        pairingOnly.a2 === 'ABSENT' &&
        /_EN \+2\/-2 MISMATCH/.test(pairingOnly.detail),
      'a1=' + pairingOnly.a1 + ' a2=' + pairingOnly.a2 + ' :: ' + pairingOnly.detail
    );

    const added = pair(['# c', ''], ['# e', ''], { [CNP]: [0, 0], [ENP]: [2, 0] });
    check(
      'T78 OB-90 an `_EN` side that is NEWLY ADDED (+2/-0) is judged against the untouched CN baseline and fails G3(a.1)',
      added.a1 === 'FAIL' && added.a2 === 'ABSENT' && /_EN \+2\/-0 MISMATCH/.test(added.detail),
      'a1=' + added.a1 + ' a2=' + added.a2 + ' :: ' + added.detail
    );

    const removed = g3aEdgeProbe({
      key: 'docs/t027/probe-pair-cn.md',
      changed: ['docs/t027/probe-pair-cn_EN.md'],
      texts: { 'docs/t027/probe-pair-cn.md': ['# c', ''] },
      numstat: { 'docs/t027/probe-pair-cn.md': [0, 0], 'docs/t027/probe-pair-cn_EN.md': [0, 4] },
      existing: [],
    });
    check(
      'T78 OB-90 a DELETED `_EN` side (+0/-4) still yields G3(a.1) FAIL, so a removal cannot pass as "nothing to compare"',
      removed.a1 === 'FAIL' && removed.a === 'FAIL' && /_EN \+0\/-4 MISMATCH/.test(removed.detail),
      'a1=' + removed.a1 + ' :: ' + removed.detail
    );

    const ROLD = 'docs/t027/probe-rename-cn.md';
    const ROLD_EN = 'docs/t027/probe-rename-cn_EN.md';
    const RNEW = 'docs/t027/probe-renamed-cn.md';
    const RNEW_EN = 'docs/t027/probe-renamed-cn_EN.md';
    const rename = g3aEdgeProbe({
      key: ROLD,
      changed: [ROLD_EN, RNEW_EN],
      texts: { [ROLD]: ['# c', ''] },
      numstat: { [ROLD]: [0, 0], [ROLD_EN]: [0, 4], [RNEW]: [0, 0], [RNEW_EN]: [2, 0] },
      existing: [ROLD],
    });
    const renameFails = failIds(rename.records);
    check(
      'T78 OB-90 an `_EN` RENAME is judged as two pairs (old mirror deleted, new mirror added), each with its own G3(a.1)',
      renameFails.indexOf('G3(a.1) ' + ROLD) >= 0 &&
        renameFails.indexOf('G3(a.1) ' + RNEW) >= 0 &&
        rename.records.filter((r) => r.id === 'G3(a)').length === 1 &&
        rename.records.length === 3,
      'fails=' + JSON.stringify(renameFails) + ' records=' + rename.records.length
    );

    const lower = g3aEdgeProbe({
      key: 'docs/t027/probe-pair-cn.md',
      changed: ['docs/t027/probe-pair-cn_en.md'],
      texts: { 'docs/t027/probe-pair-cn.md': ['# c', ''], 'docs/t027/probe-pair-cn_en.md': ['# e', ''] },
      numstat: { 'docs/t027/probe-pair-cn.md': [2, 0], 'docs/t027/probe-pair-cn_en.md': [2, 0] },
      existing: [],
    });
    check(
      'T78 OB-90 a non-canonical `_en.md` spelling is paired under the CANONICAL CN id, so the half a reader looks for is the one reported',
      lower.a === 'PASS' &&
        lower.a1 === 'ABSENT' &&
        lower.a2 === 'ABSENT' &&
        lower.records.length === 1 &&
        /^docs\/t027\/probe-pair-cn\.md \+2\/-0 vs _EN \+2\/-0 OK/.test(lower.detail),
      'a=' + lower.a + ' :: ' + lower.detail
    );

    const orphan = g3aEdgeProbe({
      key: 'docs/t027/orphan.md',
      changed: ['docs/t027/orphan_EN.md'],
      texts: { 'docs/t027/orphan_EN.md': ['# e', ''] },
      numstat: { 'docs/t027/orphan_EN.md': [2, 0] },
      existing: [],
    });
    check(
      'T78 OB-90 an orphan `_EN.md` (no CN side at all) is still judged and fails G3(a.1) against a +0/-0 baseline',
      orphan.a1 === 'FAIL' && /_EN \+2\/-0 MISMATCH/.test(orphan.detail),
      'a1=' + orphan.a1 + ' :: ' + orphan.detail
    );

    const single = g3aEdgeProbe({
      key: 'docs/t027/probe-single-lang.md',
      changed: ['docs/t027/probe-single-lang.md'],
      texts: { 'docs/t027/probe-single-lang.md': ['# s', ''] },
      numstat: { 'docs/t027/probe-single-lang.md': [3, 0] },
      existing: [],
    });
    check(
      'T78 OB-90 a single-language document without an `_EN` mirror produces NO record (no false red on the mirror-less corpus)',
      single.a === 'ABSENT' && single.records.length === 0,
      'a=' + single.a + ' records=' + single.records.length
    );

    const boldAgain = pair(['# c', '', '**b**', ''], ['# e', '', '**b** **x**', ''], { [CNP]: [0, 0], [ENP]: [1, 1] });
    check(
      'T78 OB-90 aggregate uniqueness: the same `_EN`-only drift emits exactly ONE G3(a) record (a duplicate aggregate would double-count)',
      boldAgain.records.filter((r) => r.id === 'G3(a)').length === 1 && boldAgain.a === 'FAIL',
      'G3(a) records=' + boldAgain.records.filter((r) => r.id === 'G3(a)').length + ' a=' + boldAgain.a
    );
  }

  // ---- T79: OB-91 (v1.34) - the language tag token is consumed to its FULL extent and then
  // compared WHOLE against the closed set. Table-driven (50 assertions): illegal forms, legal
  // + illegal coexistence in both orders, the [A-Za-z0-9._-] continuation charset, the
  // terminator set (any character outside that class ends the token), the header-window
  // boundary, and the remaining fail-closed premises. Every illegal row asserts the exact tag
  // back out of the error message, so a truncated token cannot pass as a legal prefix.
  {
    const CN = 'docs/t027/probe-g6-6-t79.md';
    const EN = 'docs/t027/probe-g6-6-t79_EN.md';
    const H = '### 12.4 观察项台账\n\n';
    const cnRow = '| OB-1 | d | x | 观察中 |\n';
    const enRow = '| OB-1 | d | x | Observing |\n';
    const rowFor = (file) => (file === EN ? enRow : cnRow);
    const fileFor = (tag) => (/^lang:en/.test(tag) ? EN : CN);
    const run66 = (file, head) => {
      try {
        return { res: g66Probe(file, head + H + rowFor(file)), err: null };
      } catch (e) {
        return { res: null, err: e };
      }
    };
    const unknownTag = (o, tag) =>
      !!o.err &&
      o.err instanceof UsageError &&
      /unknown language tag/.test(o.err.message) &&
      o.err.message.indexOf('`' + tag + '`') >= 0;
    const why = (o) => (o.err ? 'err=' + o.err.message : 'verdict=' + (o.res ? o.res.verdict : 'NONE'));
    const fillerLines = (n) => {
      const a = [];
      for (let i = 0; i < n; i++) a.push('filler ' + i);
      return a;
    };

    // A. The nine illegal forms named by OB-91, each asserted alone.
    const ILLEGAL = ['zh_cn', 'zh-cn', 'en_us', 'en-us', '', 'ZH', 'zh-', 'zh.', 'zh2', 'zh:'];
    for (const suffix of ILLEGAL) {
      const tag = 'lang:' + suffix;
      const o = run66(fileFor(tag), '`' + tag + '`\n');
      check(
        'T79 OB-91 illegal form `' + tag + '` is consumed whole and rejected as an unknown tag',
        unknownTag(o, tag),
        why(o)
      );
    }

    // B. A legal tag next to an illegal one is fail-closed, in BOTH orders.
    const COEXIST = [
      [CN, 'lang:zh', 'lang:zh-cn'],
      [CN, 'lang:zh-cn', 'lang:zh'],
      [EN, 'lang:en', 'lang:en-us'],
      [EN, 'lang:en-us', 'lang:en'],
    ];
    for (const [file, first, second] of COEXIST) {
      const bad = /-(cn|us)$/.test(first) ? first : second;
      const o = run66(file, '`' + first + '` `' + second + '`\n');
      check(
        'T79 OB-91 legal + illegal coexistence is fail-closed (' + first + ' + ' + second + ')',
        unknownTag(o, bad),
        why(o)
      );
    }

    // C. Continuation-charset sweep: every character class in [A-Za-z0-9._:-] extends the
    // token on both legal stems, so each suffix is an unknown tag rather than a legal prefix.
    const CONT = [
      'en_',
      'en.',
      'en-',
      'en2',
      'enA',
      'ena',
      'zhA',
      'zha',
      'zh9',
      'zh__',
      'zh-_',
      'zh.1',
      'zhCN',
      'en-001',
    ];
    for (const suffix of CONT) {
      const tag = 'lang:' + suffix;
      const o = run66(fileFor(tag), '`' + tag + '`\n');
      check(
        'T79 OB-91 continuation char extends the token (`' + tag + '` is one unknown token, not a legal prefix)',
        unknownTag(o, tag),
        why(o)
      );
    }

    // D. Terminators and normal declaration positions must NOT go red: every character
    // outside [A-Za-z0-9._:-] ends the token, so the remaining prefix is the legal tag.
    const LEGAL = [
      ['zh-alone', CN, '`lang:zh`\n'],
      ['zh-comma', CN, '`lang:zh,`\n'],
      ['zh-space', CN, '`lang:zh` zh\n'],
      ['zh-tab', CN, '`lang:zh`\tzh\n'],
      ['zh-cjk-comma', CN, '`lang:zh`，zh\n'],
      ['zh-cjk-paren', CN, '`lang:zh`）\n'],
      ['zh-declaration-line', CN, '> Language declaration: `lang:zh`. ProofRail\n'],
      ['zh-bracket', CN, '`lang:zh`]\n'],
      ['zh-paren', CN, '`lang:zh`)\n'],
      ['zh-quote', CN, '`lang:zh`"\n'],
      ['zh-pipe', CN, '`lang:zh`| \n'],
      ['zh-twice', CN, '`lang:zh` / `lang:zh`\n'],
      ['en-alone', EN, '`lang:en`\n'],
      ['en-comma', EN, '`lang:en,`\n'],
      ['en-declaration-line', EN, '> Language declaration: `lang:en`. ProofRail\n'],
      ['zh-line-15-of-window', CN, fillerLines(14).join('\n') + '\n`lang:zh`\n'],
    ];
    for (const [slug, file, head] of LEGAL) {
      const o = run66(file, head);
      check(
        'T79 OB-91 legal shape "' + slug + '" is not turned red by the token rule',
        !o.err && !!o.res && o.res.verdict === 'PASS' && o.res.findings === 0,
        why(o)
      );
    }

    // E. The other fail-closed premises, including the header-window boundary (line 16 falls
    // outside LANG_HEADER_WINDOW = 15) and the `\b` occurrence rule.
    const FAIL_CLOSED = [
      ['tag on line 16 is outside the header window', CN, fillerLines(15).join('\n') + '\n`lang:zh`\n', /language tag missing/],
      ['two legal tags conflict', CN, '`lang:zh` `lang:en`\n', /conflicting language tags/],
      ['declared tag disagrees with the file name', CN, '`lang:en`\n', /disagrees with the file name/],
      ['no tag at all', CN, '# no tag here\n', /language tag missing/],
      ['a word character before `lang:` means no occurrence', CN, 'xlang:zh\n', /language tag missing/],
      ['a bare tag next to a legal one is still unknown', CN, '`lang:zh` `lang:`\n', /unknown language tag `lang:`/],
    ];
    for (const [slug, file, head, re] of FAIL_CLOSED) {
      const o = run66(file, head);
      check(
        'T79 OB-91 fail-closed premise: ' + slug,
        !!o.err && o.err instanceof UsageError && re.test(o.err.message),
        why(o)
      );
    }
  }

  // ---- T79 canary (recheck M-3 / M-4): the table above is driven by the same constants the
  // rule uses, so a change to `LANG_TAG_CONT_RE` / `LANG_HEADER_WINDOW` / `LANG_TAG_EXACT`
  // would move the table along with the rule and keep the suite green. This block pins those
  // constants against an independent mirror and adds one end-to-end canary per EXTERNAL
  // terminator character, so widening the continuation class turns into a red check instead
  // of an invisible semantic drift.
  {
    const gi = G6.internals;
    // The independent side of the lock: a literal copy of the intended class, deliberately
    // NOT read back from g6.js (reading the source text would drift with the implementation).
    const MIRROR_CONT = /[A-Za-z0-9._:-]/;
    const mismatches = [];
    for (let code = 0; code < 128; code++) {
      const ch = String.fromCharCode(code);
      if (gi.LANG_TAG_CONT_RE.test(ch) !== MIRROR_CONT.test(ch)) {
        mismatches.push('U+' + code.toString(16).toUpperCase().padStart(4, '0'));
      }
    }
    const NON_ASCII = ['\u00A0', '\u3001', '\u3002', '\uFF08', '\uFF0C', '\uFF1B', '\u4E2D', '\u{1F600}'];
    const leaked = NON_ASCII.filter((ch) => gi.LANG_TAG_CONT_RE.test(ch));
    const tags = Object.keys(gi.LANG_TAG_EXACT).sort();
    const legalZh = (() => {
      try {
        return gi.declaredLang(['`lang:zh`'], 'docs/t027/probe-g6-6-t79-canary.md');
      } catch (e) {
        return 'ERR:' + e.message;
      }
    })();
    const illegalUnderscore = (() => {
      try {
        gi.declaredLang(['`lang:zh_cn`'], 'docs/t027/probe-g6-6-t79-canary.md');
        return 'NO ERROR';
      } catch (e) {
        return e;
      }
    })();
    check(
      'T79 canary: constants are locked — LANG_HEADER_WINDOW = 15, LANG_TAG_EXACT is exactly {lang:zh, lang:en}, and LANG_TAG_CONT_RE is exactly [A-Za-z0-9._:-] (behavior pair: `lang:zh` legal, `lang:zh_cn` illegal)',
      gi.LANG_HEADER_WINDOW === 15 &&
        JSON.stringify(tags) === JSON.stringify(['lang:en', 'lang:zh']) &&
        gi.LANG_TAG_EXACT['lang:zh'] === 'zh' &&
        gi.LANG_TAG_EXACT['lang:en'] === 'en' &&
        mismatches.length === 0 &&
        leaked.length === 0 &&
        legalZh === 'cn' &&
        illegalUnderscore instanceof UsageError &&
        /unknown language tag/.test(illegalUnderscore.message) &&
        illegalUnderscore.message.indexOf('`lang:zh_cn`') >= 0,
      JSON.stringify({
        window: gi.LANG_HEADER_WINDOW,
        tags: tags,
        mismatches: mismatches,
        leaked: leaked,
        legalZh: legalZh,
        illegal: illegalUnderscore instanceof UsageError ? illegalUnderscore.message : String(illegalUnderscore),
      })
    );

    // M-5 (recheck): the mirror above compares BEHAVIOUR, so a widening applied to both the
    // rule and the mirror stays green. The source TEXT of the class is therefore frozen as a
    // literal string: any edit to the RegExp source changes this value and turns red.
    check(
      'T79 canary: LANG_TAG_CONT_RE source text is frozen to the literal class [A-Za-z0-9._:-]',
      gi.LANG_TAG_CONT_RE.source === '[A-Za-z0-9._:-]',
      'source=' + JSON.stringify(gi.LANG_TAG_CONT_RE.source) + ' flags=' + JSON.stringify(gi.LANG_TAG_CONT_RE.flags)
    );

    // M-4: one end-to-end canary per external terminator. None of these characters is in the
    // continuation set, so the token ends at the legal prefix `lang:zh` and the tail is free
    // text - the document PASSes. If the class is widened to include any of them the token
    // becomes `lang:zh<sep>en`, an unknown tag, and the canary goes red.
    const CN = 'docs/t027/probe-g6-6-t79-canary.md';
    const H = '### 12.4 观察项台账\n\n';
    const run66 = (head) => {
      try {
        return { res: g66Probe(CN, head + H + '| OB-1 | d | x | 观察中 |\n'), err: null };
      } catch (e) {
        return { res: null, err: e };
      }
    };
    const why = (o) =>
      o.err ? 'err=' + o.err.message : 'verdict=' + (o.res ? o.res.verdict : 'NONE') + ' findings=' + (o.res ? o.res.findings : '?');
    const EXTERNAL = [
      ['/', 'U+002F', 'slash'],
      ['?', 'U+003F', 'question mark'],
      ['=', 'U+003D', 'equals sign'],
      [' ', 'U+0020', 'space'],
      ['\uFF1B', 'U+FF1B', 'fullwidth semicolon'],
      ['\uFF0C', 'U+FF0C', 'fullwidth comma'],
      ['\u3001', 'U+3001', 'ideographic comma'],
    ];
    for (const [ch, cp, slug] of EXTERNAL) {
      const tag = 'lang:zh' + ch + 'en';
      const o = run66('`' + tag + '`\n');
      check(
        'T79 canary: ' +
          slug +
          ' (' +
          cp +
          ') is an external terminator, not a continuation char — `' +
          tag +
          '` splits into the legal tag `lang:zh` plus free text and stays green',
        !o.err && !!o.res && o.res.verdict === 'PASS' && o.res.findings === 0,
        why(o)
      );
    }
  }

  // ---- T79 black-box tokenization (recheck M-5): the locks above read the live constants,
  // so a rule and its lock widened together would stay green. This block observes
  // tokenization ONLY through the public G6-6 path on synthetic documents, with literal
  // expectations and no `gi.*` reference at all.
  {
    const FILE = 'docs/t027/probe-g6-6-t79-tokenization.md';
    const LEDGER = '### 12.4 观察项台账\n\n| OB-1 | d | x | 观察中 |\n';
    const observe = (head) => {
      try {
        const res = g66Probe(FILE, head + LEDGER);
        return { verdict: res.verdict, findings: res.findings };
      } catch (e) {
        return { err: e };
      }
    };
    const zh = observe('`lang:zh`\n');
    const zhCn = observe('`lang:zh_cn`\n');
    const zhSlashEn = observe('`lang:zh/en`\n');
    const shown = (o) => (o.err ? 'err=' + o.err.message : o.verdict + '/' + o.findings);
    check(
      'T79 canary: black-box tokenization — `lang:zh` PASSes, `lang:zh_cn` is an unknown tag, `lang:zh/en` PASSes (/ terminates the token)',
      !zh.err &&
        zh.verdict === 'PASS' &&
        zh.findings === 0 &&
        !zhSlashEn.err &&
        zhSlashEn.verdict === 'PASS' &&
        zhSlashEn.findings === 0 &&
        zhCn.err instanceof UsageError &&
        /unknown language tag/.test(zhCn.err.message) &&
        zhCn.err.message.indexOf('`lang:zh_cn`') >= 0,
      JSON.stringify({ zh: shown(zh), 'zh_cn': shown(zhCn), 'zh/en': shown(zhSlashEn) })
    );
  }

  // ---- T76: v1.36 - base G1 (status markers) is a real criterion, not a SKIP.
  {
    const LEDGER = (o) => {
      const st = o.status || [['A0 alpha', '✅ 完成'], ['C1 gamma', '⬜ 未开始']];
      const heads = o.heads || ['## A0 — alpha · ✅', '## C1 — gamma · ⬜', '### P5 — `X` · ⏸ 待评估（产品侧）'];
      const comp = o.comp || [['2026-09-14', 'A0', 'done', '`docs/validation/a0.md`']];
      return [
        '状态符号：⬜ 未开始 / 🔄 进行中 / ✅ 完成 / ⛔ 阻断',
        '',
        '## 状态表',
        '',
        '| 切片 | 状态 | 直接依赖 | 规模 | 完成日期 | 证据 |',
        '|---|---|---|---|---|---|',
        ...st.map((r) => '| ' + r[0] + ' | ' + r[1] + ' | — | M | | |'),
        '',
        ...heads,
        '',
        '## 完成记录（追加）',
        '',
        '| 日期 | 切片 | 完成摘要 | 证据 |',
        '|---|---|---|---|',
        ...comp.map((r) => '| ' + r.join(' | ') + ' |'),
      ].join('\n');
    };
    const good = g1Probe(LEDGER({}));
    check(
      'T76 G1 a consistent ledger (heading / status table / completion log) passes all four sub-checks',
      good.c === 'PASS' && good.d === 'PASS' && good.e === 'PASS' && good.f === 'PASS',
      good.detail
    );
    const unknownSym = g1Probe(LEDGER({ status: [['A0 alpha', '🚧 进行中'], ['C1 gamma', '⬜ 未开始']] }));
    check(
      'T76 G1 an unknown state symbol outside the ledger closed set fails G1(c)',
      unknownSym.c === 'FAIL',
      unknownSym.detail
    );
    const pendingOnStatus = g1Probe(LEDGER({ status: [['A0 alpha', '⏸ 待评估'], ['C1 gamma', '⬜ 未开始']] }));
    check(
      'T76 G1 the pending-evaluation marker is only legal on headings, never in the status table',
      pendingOnStatus.c === 'FAIL',
      pendingOnStatus.detail
    );
    const drift = g1Probe(LEDGER({ heads: ['## A0 — alpha · ⬜', '## C1 — gamma · ⬜'] }));
    check(
      'T76 G1 a heading that disagrees with its status-table row fails G1(d)',
      drift.d === 'FAIL' && drift.c === 'PASS',
      drift.detail
    );
    const compNotDone = g1Probe(LEDGER({ comp: [['2026-09-14', 'C1', 'done', '`docs/validation/c1.md`']] }));
    check(
      'T76 G1 a completion-log row whose status row is still ⬜ fails G1(e)',
      compNotDone.e === 'FAIL',
      compNotDone.detail
    );
    const inProg = g1Probe(LEDGER({ comp: [['2026-09-14', 'A0', '进行中：待 ⑧b 定稿后更新', '`docs/validation/a0.md`']] }));
    check(
      'T76 G1 in-progress wording inside a completion row fails G1(f) (the GATES-EDGE-COVERAGE defect, mechanised)',
      inProg.f === 'FAIL' && inProg.e === 'PASS',
      inProg.detail
    );
    const enInProg = g1Probe(LEDGER({ comp: [['2026-09-14', 'A0', 'In progress: awaiting finalisation', '`docs/validation/a0.md`']] }));
    check(
      'T76 G1 the same holds for the English wording (In progress / awaiting finalisation)',
      enInProg.f === 'FAIL',
      enInProg.detail
    );
    const aliasHeads = ['## A0 — alpha · ✅', '## P11 — `GATES-G1G5` · 🔄'];
    const aliasRows = [['A0 alpha', '✅ 完成'], ['P11 slice', '🔄 进行中']];
    const aliasComp = [
      ['2026-09-14', 'A0', 'done', '`docs/validation/a0.md`'],
      ['2026-10-18', 'GATES-G1G5', 'closed', '`docs/validation/p11.md`'],
    ];
    const aliasBad = g1Probe(LEDGER({ heads: aliasHeads, status: aliasRows, comp: aliasComp }));
    check(
      'T76 G1 a completion key binds to a heading through its backticked alias (the P-series governance case)',
      aliasBad.e === 'FAIL',
      aliasBad.detail
    );
    const aliasOk = g1Probe(
      LEDGER({
        heads: ['## A0 — alpha · ✅', '## P11 — `GATES-G1G5` · ✅'],
        status: [['A0 alpha', '✅ 完成'], ['P11 slice', '✅ 完成']],
        comp: aliasComp,
      })
    );
    check('T76 G1 the same row passes once the aliased heading is ✅ (the binding is real, not a lucky string match)', aliasOk.e === 'PASS', aliasOk.detail);
    const G1_ALIAS_LOOKUP = '  const byAlias = headings.filter((h) => h.aliases.indexOf(key) >= 0);';
    const mutG1Alias = loadMutatedModule(path.join(__dirname, 'lib', 'criteria', 'g1.js'), G1_ALIAS_LOOKUP, '  const byAlias = [];', 'v136-g1-alias');
    check(
      'T76 G1 mutation applied on the alias lookup (defect 3: completion key never bound to a P-series heading)',
      mutG1Alias.applied === true,
      'applied=' + mutG1Alias.applied
    );
    const blindAlias = g1Probe(LEDGER({ heads: aliasHeads, status: aliasRows, comp: aliasComp }), mutG1Alias.exports);
    check(
      'T76 G1 load-bearing proof: with the alias fallback removed the SAME mismatched row passes, i.e. G1(e) was vacuous for every governance slice',
      blindAlias.e === 'PASS',
      blindAlias.detail
    );
    check(
      'T76 G1 restore after the alias mutation window',
      mutG1Alias.identical === true && g1Probe(LEDGER({ heads: aliasHeads, status: aliasRows, comp: aliasComp })).e === 'FAIL',
      'identical=' + mutG1Alias.identical
    );
    const unknownSymHead = g1Probe(
      LEDGER({
        heads: ['## A0 — alpha · ✅', '## C1 — gamma · 🚧 进行中'],
        status: [['A0 alpha', '✅ 完成'], ['C1 gamma', '⬜ 未开始']],
        comp: [],
      })
    );
    check(
      'T76 G1 an UNKNOWN state symbol on a heading reaches G1(c) instead of slipping past every guard (r2 F-2)',
      unknownSymHead.c === 'FAIL',
      unknownSymHead.detail
    );
    const histText = LEDGER({ comp: [['2026-09-14', 'A0', '进行中：待 ⑧b 定稿后更新', '`docs/validation/a0.md`']] });
    const histAll = g1Probe(histText);
    const histNone = g1Probe(histText, undefined, []);
    check(
      'T76 G1 G1(f) judges ADDED completion rows only: the same historical row fails when it is in the added set and passes when it is not (r2 F-3)',
      histAll.f === 'FAIL' && histNone.f === 'PASS',
      'all=' + histAll.f + ' none=' + histNone.f
    );
    const infoE = g1Probe(LEDGER({ comp: [['2026-09-14', 'ZZ', 'done', '`docs/validation/z.md`']] }));
    check(
      'T76 G1 an unbound completion key also emits an explicit INFO record, so the directive promise "no binding is INFO" is literal (r2 F-4)',
      /G1\(e-info\)=INFO/.test(infoE.detail),
      infoE.detail
    );
    let declErr = null;
    try {
      g1Probe(LEDGER({}).replace('状态符号：⬜ 未开始 / 🔄 进行中 / ✅ 完成 / ⛔ 阻断\n', ''));
    } catch (e) {
      declErr = e;
    }
    check(
      'T76 G1 a ledger without its state-symbol declaration line raises UsageError (exit 2, never a silent guess)',
      declErr instanceof UsageError && /declaration line/.test(declErr.message),
      'err=' + (declErr ? declErr.message : 'NO ERROR')
    );
    const mutG1 = loadMutatedModule(
      path.join(__dirname, 'lib', 'criteria', 'g1.js'),
      G1_INPROGRESS_TEST,
      '      if (false) badF.push(',
      'v136-g1-inprogress'
    );
    check(
      'T76 G1 mutation applied: the in-progress test line matched g1.js (a stale patch would weaken the suite silently)',
      mutG1.applied === true,
      'applied=' + mutG1.applied
    );
    const blindG1 = g1Probe(LEDGER({ comp: [['2026-09-14', 'A0', '进行中：待 ⑧b 定稿后更新', '`docs/validation/a0.md`']] }), mutG1.exports);
    check(
      'T76 G1 load-bearing proof: with the wording test patched out the SAME in-progress row passes',
      blindG1.f === 'PASS',
      blindG1.detail
    );
    check(
      'T76 G1 restore: g1.js is byte-identical after the mutation window and the wording row fails again',
      mutG1.identical === true &&
        g1Probe(LEDGER({ comp: [['2026-09-14', 'A0', '进行中：待 ⑧b 定稿后更新', '`docs/validation/a0.md`']] })).f === 'FAIL',
      'identical=' + mutG1.identical
    );
    const unreadable = (() => {
      const out = createReport();
      G1.run({ repoRoot: discoverRepoRoot(__dirname), readLineArray: () => null }, out);
      const r = out.records[0] || {};
      return r.verdict + ' :: ' + r.detail;
    })();
    check(
      'T76 G1 an unreadable ledger pair reports SKIP (never a vacuous PASS)',
      /^SKIP :: ledger pair not readable/.test(unreadable),
      unreadable
    );
  }

  // ---- T77: v1.36 - base G5 (change-set contract) derives its four fields instead of asking for a JSON.
  {
    const R = 'docs/validation/probe-report.md';
    const FENCE = '```artifacts\n' + R + '\trepo\n' + 'tools/probe.js\trepo\n```';
    const LEDGER_ROW = '| 2026-10-18 | PROBE | summary | `' + R + '` |';
    const pass = g5Probe({
      changed: [R, 'tools/probe.js'],
      texts: { [R]: FENCE },
      added: { 'docs/t027/REMAINING_SLICES.md': [LEDGER_ROW] },
      existing: [R],
      unstaged: [],
    });
    check(
      'T77 G5 a declared, staged, resolvable change set passes G5(d) / G5(e) / G5(f)',
      pass.d === 'PASS' && pass.e === 'PASS' && pass.f === 'PASS',
      pass.detail
    );
    const missingDecl = g5Probe({
      changed: [R, 'tools/probe.js', 'tools/undeclared.js'],
      texts: { [R]: FENCE },
      added: { 'docs/t027/REMAINING_SLICES.md': [LEDGER_ROW] },
      existing: [R],
      unstaged: [],
    });
    check(
      'T77 G5 a staged path nobody declared fails G5(d) - the observable consequence of `git add -A`',
      missingDecl.d === 'FAIL' && /tools\/undeclared\.js/.test(missingDecl.raw),
      missingDecl.raw
    );
    const extraDecl = g5Probe({
      changed: [R],
      texts: { [R]: FENCE },
      added: {},
      existing: [R],
      unstaged: [],
    });
    check(
      'T77 G5 the reverse direction (declared but untouched) is INFO, because a fence is the landing snapshot - the §11.2 backfill case',
      extraDecl.d === 'PASS' && /declared-not-touched/.test(extraDecl.raw),
      extraDecl.raw
    );
    const badRef = g5Probe({
      changed: [R],
      texts: { [R]: FENCE.replace('tools/probe.js\t', '') },
      added: { 'docs/t027/REMAINING_SLICES.md': ['| 2026-10-18 | PROBE | summary | `docs/validation/ghost.md` |'] },
      existing: [],
      unstaged: [],
    });
    check(
      'T77 G5 an evidence reference that does not resolve fails G5(e)',
      badRef.e === 'FAIL',
      badRef.detail
    );
    const dirty = g5Probe({
      changed: [R],
      texts: { [R]: '```artifacts\n' + R + '\trepo\n```' },
      added: {},
      existing: [R],
      unstaged: [' M tools/probe.js'],
    });
    check(
      'T77 G5 an unstaged modification in index scope fails G5(f) (index must equal the worktree)',
      dirty.f === 'FAIL',
      dirty.detail
    );
    const stagedOnly = g5Probe({
      changed: [R],
      texts: { [R]: FENCE },
      added: {},
      existing: [R],
      unstaged: ['M  ' + R],
    });
    check(
      'T77 G5 a STAGED entry (porcelain X column) does not dirty the index - staging is the normal state of index scope',
      stagedOnly.f === 'PASS',
      stagedOnly.detail
    );
    const noReport = g5Probe({ changed: ['tools/probe.js'], added: {}, existing: [], unstaged: [] });
    check(
      'T77 G5 a change set without a report records INFO on G5(d) instead of a vacuous pass',
      noReport.d === 'INFO' && noReport.e === 'PASS',
      noReport.detail
    );
    const treeScope = g5Probe({
      kind: 'ci',
      changed: [R],
      texts: { [R]: '```artifacts\n' + R + '\trepo\n```' },
      added: {},
      existing: [R],
      unstaged: [' M whatever.js'],
    });
    check(
      'T77 G5 outside index scope the staging check is INFO (a clean CI checkout has no index state)',
      treeScope.f === 'INFO' && treeScope.d === 'PASS',
      treeScope.detail
    );
    // external JSON cross-check: a mismatch is a usage error, never a silent pass
    const prevEnv = process.env.GATES_CHANGESET_JSON;
    let envErr = null;
    try {
      process.env.GATES_CHANGESET_JSON = JSON.stringify({ stagedFiles: ['docs/other.md'] });
      g5Probe({ changed: [R], texts: { [R]: '```artifacts\n' + R + '\trepo\n```' }, added: {}, existing: [R], unstaged: [] });
    } catch (e) {
      envErr = e;
    } finally {
      if (prevEnv === undefined) delete process.env.GATES_CHANGESET_JSON;
      else process.env.GATES_CHANGESET_JSON = prevEnv;
    }
    check(
      'T77 G5 an external GATES_CHANGESET_JSON that disagrees with the derived fields raises UsageError (exit 2)',
      !!envErr && envErr instanceof UsageError && /disagrees with the derived change set/.test(envErr.message),
      'err=' + (envErr ? envErr.message : 'NO ERROR')
    );
    check(
      'T77 G5 a declared-but-untouched path also emits an explicit INFO record (r2 F-4)',
      /G5\(d-info\)=INFO/.test(extraDecl.detail),
      extraDecl.detail
    );
    const mutG5Del = loadMutatedModule(
      path.join(__dirname, 'lib', 'criteria', 'g5.js'),
      '    return !rec || !rec.deleted;',
      '    return true;',
      'v136-g5-delete'
    );
    check(
      'T77 G5 mutation applied on the deletion guard, which must read the parser boolean and not the two-char XY status (r2 F-1)',
      mutG5Del.applied === true,
      'applied=' + mutG5Del.applied
    );
    const blindDel = g5Probe(
      {
        changed: [R, 'tools/gone.js'],
        deleted: ['tools/gone.js'],
        texts: { [R]: FENCE },
        added: {},
        existing: [R],
        unstaged: [],
      },
      mutG5Del.exports
    );
    check(
      'T77 G5 load-bearing proof: with the deletion guard patched out the same deleted path is reported as undeclared (r2 F-1)',
      blindDel.d === 'FAIL',
      blindDel.detail
    );
    check(
      'T77 G5 restore after the deletion-guard mutation window',
      mutG5Del.identical === true &&
        g5Probe({ changed: [R, 'tools/gone.js'], deleted: ['tools/gone.js'], texts: { [R]: FENCE }, added: {}, existing: [R], unstaged: [] }).d === 'PASS',
      'identical=' + mutG5Del.identical
    );
    check(
      'T77 G5 a declared-but-untouched path also emits an explicit INFO record (r2 F-4)',
      /G5\(d-info\)=INFO/.test(extraDecl.detail),
      extraDecl.detail
    );
    const bare = G5._internals.evidencePathsFromRow(
      '| 2026-10-18 | X | done | docs/validation/evidence/x.md, tools/gates/gate.js |'
    );
    check(
      'T77 G5 evidence references are read from BARE paths too (30/30 corpus rows carry no backticks, so backtick-only extraction was vacuous)',
      bare.length === 2 && bare.indexOf('tools/gates/gate.js') >= 0 && bare.indexOf('docs/validation/evidence/x.md') >= 0,
      JSON.stringify(bare)
    );
    const rootRef = G5._internals.evidencePathsFromRow('| 2026-10-18 | X | done | README.md, tools/gates/gate.js |');
    check(
      'T77 G5 a repository-root artifact in the evidence cell is extracted too (⑥ finding: the prefix whitelist made G5(e) vacuous for root files)',
      rootRef.indexOf('README.md') >= 0 &&
        G5._internals.evidencePathsFromRow('| 2026-10-18 | X | done | go.mod |').indexOf('go.mod') >= 0 &&
        G5._internals.evidencePathsFromRow('| 2026-10-18 | X | done | go.sum |').indexOf('go.sum') >= 0 &&
        G5._internals.evidencePathsFromRow('| 2026-10-18 | X | done | module name go.modification is not a path |').length === 0 &&
        G5._internals.evidencePathsFromRow('| 2026-10-18 | X | done | see (go.mod) and README.md) |').length === 2,
      JSON.stringify(rootRef)
    );
    const G5_BARE_SCAN =
      '  for (const m of cell.matchAll(/([\\s（(、,，；;）]|^)((?:docs|tools|examples|schemas|cmd|internal|harnesses|testdata|\\.github)\\/[^\\s；;、,，|)`）]+)/g)) {';
    const mutG5Bare = loadMutatedModule(path.join(__dirname, 'lib', 'criteria', 'g5.js'), G5_BARE_SCAN, '  for (const m of cell.matchAll(/(?!x)x/g)) {', 'v136-g5-bare');
    check(
      'T77 G5 mutation applied on the bare-path scan (defect 2: evidenceRefs were backtick-only)',
      mutG5Bare.applied === true,
      'applied=' + mutG5Bare.applied
    );
    const blindBare = mutG5Bare.exports._internals.evidencePathsFromRow(
      '| 2026-10-18 | X | done | docs/validation/evidence/x.md, tools/gates/gate.js |'
    );
    check(
      'T77 G5 load-bearing proof: with the bare-path scan removed the same row yields no reference at all',
      blindBare.length === 0,
      JSON.stringify(blindBare)
    );
    check(
      'T77 G5 restore after the bare-path mutation window',
      mutG5Bare.identical === true &&
        G5._internals.evidencePathsFromRow('| 2026-10-18 | X | done | tools/gates/gate.js |').length === 1,
      'identical=' + mutG5Bare.identical
    );
    const del = g5Probe({
      changed: [R, 'tools/gone.js'],
      deleted: ['tools/gone.js'],
      texts: { [R]: FENCE },
      added: {},
      existing: [R],
      unstaged: [],
    });
    check(
      'T77 G5 a deleted staged path is not a change-set member (removing a file declares nothing)',
      del.d === 'PASS',
      del.detail
    );
    const mutG5 = loadMutatedModule(
      path.join(__dirname, 'lib', 'criteria', 'g5.js'),
      G5_SET_EQUALITY,
      "      true ? 'PASS' : 'FAIL',",
      'v136-g5-equality'
    );
    check(
      'T77 G5 mutation applied: the set-equality line matched g5.js (a stale patch would weaken the suite silently)',
      mutG5.applied === true,
      'applied=' + mutG5.applied
    );
    const blindG5 = g5Probe(
      {
        changed: [R, 'tools/undeclared.js'],
        texts: { [R]: FENCE },
        added: {},
        existing: [R],
        unstaged: [],
      },
      mutG5.exports
    );
    check(
      'T77 G5 load-bearing proof: with the equality verdict patched the SAME undeclared path passes',
      blindG5.d === 'PASS',
      blindG5.detail
    );
    check(
      'T77 G5 restore: g5.js is byte-identical after the mutation window and the undeclared path fails again',
      mutG5.identical === true &&
        g5Probe({ changed: [R, 'tools/undeclared.js'], texts: { [R]: FENCE }, added: {}, existing: [R], unstaged: [] }).d === 'FAIL',
      'identical=' + mutG5.identical
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

    // OB-46 (v1.33) - TIGHTENING, not relaxation. This expectation REPLACES the former
    // "rename `### C.1` => G6-5 FAIL via §3.10 unresolved" assertion. That old expectation
    // was itself the defect made visible: losing the delimiter did not raise an error, it
    // silently stitched an unbounded region onto the C.0 start, and the verdict was then
    // decided by an unrelated §3.10 row. The new expectation is strictly stronger
    // (fail-closed): the same mutation must now abort with a UsageError, so the criterion
    // never issues a verdict at all when it cannot delimit the region it was asked to judge.
    let endptErr = null;
    try {
      runSubset(
        makeCtx(['g6-5-c1-v31-table'], {
          mutations: [{ key: 'g6-5-c1-v31-table', op: 'replace', from: '### C.1', to: '### C.1X' }],
        }),
        ALL_SUBCHECKS
      );
    } catch (e) {
      endptErr = e;
    }
    check(
      'T21 OB-46 fail-closed: renaming `### C.1` aborts with UsageError instead of judging an unbounded region',
      !!endptErr && endptErr instanceof UsageError && /endpoint missing/.test(endptErr.message),
      'err=' + (endptErr ? endptErr.message : 'NO ERROR (region degraded silently)')
    );
  }

  // ---- T21b: the OB-46 endpoint rule is a premise check, not a region killer.
  {
    const T = G6.internals.targetRegionRanges;
    check(
      'T21b OB-46 positive control: an intact ledger keeps both target regions (no over-rejection)',
      T(['## 0.3 x', '| a |', '### C.0 y', '### C.1 z', '| w |']).length === 2,
      'n=' + T(['## 0.3 x', '| a |', '### C.0 y', '### C.1 z', '| w |']).length
    );
    check(
      'T21b OB-46 boundary: a file with neither heading stays region-free and raises nothing (ordinary markdown)',
      T(['# t', '| x |']).length === 0,
      'n=' + T(['# t', '| x |']).length
    );
    let cErr = null;
    try {
      T(['### C.0 y', '### C.2 z']);
    } catch (e) {
      cErr = e;
    }
    check(
      'T21b OB-46 fail-closed: `### C.0` without a `### C.1` delimiter raises UsageError (started-but-unterminated region)',
      !!cErr && cErr instanceof UsageError && /endpoint missing/.test(cErr.message),
      'err=' + (cErr ? cErr.message : 'NO ERROR')
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
    const above = g66Probe('docs/t027/probe-g6-6-above.md', LANG_TAG_ZH + shared + H124 + '| OB-2 | d | x | 观察中 |\n');
    check(
      'T26 a ledger-shaped row ABOVE the §12.4 heading is out of region and not judged',
      above.verdict === 'PASS' && above.findings === 0,
      'verdict=' + above.verdict + ' findings=' + above.findings
    );
    const nonOb = g66Probe('docs/t027/probe-g6-6-nonob.md', LANG_TAG_ZH + H124 + '| X-1 | d | x | 下一片必做 |\n');
    check(
      'T26 a non-OB table row inside §12.4 is not judged (LEDGER_ROW_RE)',
      nonOb.verdict === 'PASS' && nonOb.findings === 0,
      'verdict=' + nonOb.verdict + ' findings=' + nonOb.findings
    );
    const noRegion = g66Probe('docs/t027/probe-g6-6-noregion.md', LANG_TAG_ZH + '| OB-1 | d | x | 下一片必做 |\n');
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
    const cnName = g66Probe('docs/t027/probe-g6-6-cn.md', LANG_TAG_ZH + H124 + row + '**Observing** |\n');
    const enName = g66Probe('docs/t027/probe-g6-6_EN.md', LANG_TAG_EN + H124 + row + '**Observing** |\n');
    const cnWordInEn = g66Probe('docs/t027/probe-g6-6-mix_EN.md', LANG_TAG_EN + H124 + row + '观察中 |\n');
    check(
      'T28 the `_EN` word-table switch is load-bearing: EN word fails under a CN name, passes under `_EN`, CN word fails under `_EN`',
      cnName.verdict === 'FAIL' && enName.verdict === 'PASS' && cnWordInEn.verdict === 'FAIL',
      JSON.stringify({ cnName: cnName.verdict, enName: enName.verdict, cnWordInEn: cnWordInEn.verdict })
    );
  }

  // ---- T28b: OB-47 - the word table follows the DECLARED language, not the file name.
  {
    const gi = G6.internals;
    check(
      'T28b OB-47 declaredLang reads the header tag and agrees with the name (zh -> CN table, en -> EN table)',
      gi.declaredLang(['`lang:zh`', 'x'], 'docs/t027/a.md') === 'cn' &&
        gi.declaredLang(['`lang:en`', 'x'], 'docs/t027/a_EN.md') === 'en',
      JSON.stringify([
        gi.declaredLang(['`lang:zh`', 'x'], 'docs/t027/a.md'),
        gi.declaredLang(['`lang:en`', 'x'], 'docs/t027/a_EN.md'),
      ])
    );
    const g66Err = (p, t) => {
      try {
        g66Probe(p, t);
        return null;
      } catch (e) {
        return e;
      }
    };
    const row1 = H124 + '| OB-1 | d | x | 观察中 |\n';
    const missing = g66Err('docs/t027/probe-g6-6-notag.md', row1);
    check(
      'T28b OB-47 fail-closed: no declared language tag raises UsageError instead of defaulting to CN',
      !!missing && missing instanceof UsageError && /language tag missing/.test(missing.message),
      'err=' + (missing ? missing.message : 'NO ERROR')
    );
    const unknown = g66Err('docs/t027/probe-g6-6-third.md', '`lang:fr`\n' + row1);
    check(
      'T28b OB-47 fail-closed: a third language (`lang:fr`) raises UsageError instead of being judged as CN',
      !!unknown && unknown instanceof UsageError && /unknown language tag/.test(unknown.message),
      'err=' + (unknown ? unknown.message : 'NO ERROR')
    );
    const prefixed = g66Err('docs/t027/probe-g6-6-prefixed.md', '`lang:zh-cn`\n' + row1);
    check(
      'T28b OB-47 fail-closed: a PREFIX form (`lang:zh-cn`) raises UsageError instead of counting as zh (7 r1 High-2)',
      !!prefixed && prefixed instanceof UsageError && /unknown language tag/.test(prefixed.message),
      'err=' + (prefixed ? prefixed.message : 'NO ERROR')
    );
    const prefixedEn = g66Err('docs/t027/probe-g6-6-prefixed_EN.md', '`lang:en-us`\n' + row1);
    check(
      'T28b OB-47 fail-closed: the same holds for `lang:en-us` on the `_EN` side',
      !!prefixedEn && prefixedEn instanceof UsageError && /unknown language tag/.test(prefixedEn.message),
      'err=' + (prefixedEn ? prefixedEn.message : 'NO ERROR')
    );
    // 7 r2 Medium: a legal tag must NOT shield an illegal one in the same header.
    const coexist = g66Err('docs/t027/probe-g6-6-coexist.md', LANG_TAG_ZH + '`lang:zh-cn`\n' + row1);
    check(
      'T28b OB-91 fail-closed: a legal tag plus an illegal one (`lang:zh` + `lang:zh-cn`) raises UsageError (7 r2 Medium)',
      !!coexist && coexist instanceof UsageError && /unknown language tag/.test(coexist.message),
      'err=' + (coexist ? coexist.message : 'NO ERROR')
    );
    const coexistEn = g66Err('docs/t027/probe-g6-6-coexist_EN.md', LANG_TAG_EN + '`lang:en-us`\n' + row1);
    check(
      'T28b OB-91 fail-closed: the same coexistence rule holds on the `_EN` side (`lang:en` + `lang:en-us`)',
      !!coexistEn && coexistEn instanceof UsageError && /unknown language tag/.test(coexistEn.message),
      'err=' + (coexistEn ? coexistEn.message : 'NO ERROR')
    );
    const disagree = g66Err('docs/t027/probe-g6-6-mismatch_EN.md', LANG_TAG_ZH + row1);
    check(
      'T28b OB-47 fail-closed: a renamed file (`_EN.md` keeping `lang:zh`) raises UsageError instead of switching the table',
      !!disagree && disagree instanceof UsageError && /disagrees with the file name/.test(disagree.message),
      'err=' + (disagree ? disagree.message : 'NO ERROR')
    );
    const conflict = g66Err('docs/t027/probe-g6-6-conflict.md', LANG_TAG_ZH + LANG_TAG_EN + row1);
    check(
      'T28b OB-47 fail-closed: two conflicting declarations raise UsageError',
      !!conflict && conflict instanceof UsageError && /conflicting language tags/.test(conflict.message),
      'err=' + (conflict ? conflict.message : 'NO ERROR')
    );
    const late = g66Err('docs/t027/probe-g6-6-late.md', 'x\n'.repeat(20) + LANG_TAG_ZH + row1);
    check(
      'T28b OB-47 the header window is bounded: a tag below it is not a declaration',
      !!late && late instanceof UsageError && /language tag missing/.test(late.message),
      'err=' + (late ? late.message : 'NO ERROR')
    );
    const zhOk = g66Probe('docs/t027/probe-g6-6-tag-zh.md', LANG_TAG_ZH + row1);
    const enOk = g66Probe(
      'docs/t027/probe-g6-6-tag-en_EN.md',
      LANG_TAG_EN + H124 + '| OB-1 | d | x | Observing |\n'
    );
    check(
      'T28b OB-47 positive control: a compliant declaration is never rejected on either side',
      zhOk.verdict === 'PASS' && zhOk.findings === 0 && enOk.verdict === 'PASS' && enOk.findings === 0,
      JSON.stringify({ zh: zhOk.verdict, en: enOk.verdict, enFindings: enOk.findings })
    );
    // The declaration is a premise only where a word table is actually needed: an unrelated
    // document may own a §12.4 heading (the RFC does) and has no ledger rows to judge.
    const foreign = g66Err('docs/RFC-proofrail-unattended-ai-engineering-product.md', H124 + '| a | b | c |\n');
    check(
      'T28b OB-47 scope: a §12.4-shaped heading with no ledger row is skipped, not forced to declare (the RFC case)',
      foreign === null,
      'err=' + (foreign ? foreign.message : 'none (skipped)')
    );
    const foreignRow = g66Err(
      'docs/RFC-proofrail-unattended-ai-engineering-product.md',
      H124 + '| OB-1 | d | x | 观察中 |\n'
    );
    check(
      'T28b OB-47 scope: the same foreign file carrying a real ledger row is judged, so it must declare (fail closed)',
      !!foreignRow && foreignRow instanceof UsageError && /language tag missing/.test(foreignRow.message),
      'err=' + (foreignRow ? foreignRow.message : 'NO ERROR')
    );
  }

  // ---- T28c: OB-73 - the sixth mirror shape (whole-file backtick count).
  {
    const cnOk = '# CN\n\nshape `a` and `b`\n';
    const enOk = '# EN\n\nshape `a` and `b`\n';
    const enDrift = '# EN\n\nshape `a` and `b` and `c`\n';
    const clean = g3aProbe(cnOk, enOk);
    check(
      'T28c OB-73 positive control: equal backtick counts leave G3(a.2) silent and the aggregate green',
      clean.a2 === 'ABSENT' && clean.a === 'PASS' && /"backticks":true/.test(clean.detail),
      'a2=' + clean.a2 + ' a=' + clean.a + ' :: ' + clean.detail
    );
    const drifted = g3aProbe(cnOk, enDrift);
    check(
      'T28c OB-73 the sixth shape is load-bearing: one extra backtick on the `_EN` side turns G3(a.2) red',
      drifted.a2 === 'FAIL' && /"backticks":false/.test(drifted.a2detail),
      'a2=' + drifted.a2 + ' :: ' + drifted.a2detail
    );
    const mut = loadMutatedModule(G3A_SOURCE, G3A_BACKTICK_SHAPE, '', 'ob73-backticks');
    check(
      'T28c OB-73 mutation applied: the sixth-shape line matched g3a.js (a stale patch would weaken the suite silently)',
      mut.applied === true,
      'applied=' + mut.applied
    );
    const blind = g3aProbe(cnOk, enDrift, mut.exports);
    check(
      'T28c OB-73 load-bearing proof: with the sixth shape patched out the SAME drift passes (the shape, not the pair, is what catches it)',
      blind.a2 === 'ABSENT',
      'a2=' + blind.a2 + ' :: ' + blind.detail
    );
    check(
      'T28c restore: g3a.js is byte-identical after the mutation window and the cached module still catches the drift',
      mut.identical === true && g3aProbe(cnOk, enDrift).a2 === 'FAIL',
      'identical=' + mut.identical
    );
    const root = discoverRepoRoot(__dirname);
    const real = g3aProbe(
      decodeBuffer(fs.readFileSync(path.join(root, 'docs', 'DELIVERY_DIRECTIVE.md'))),
      decodeBuffer(fs.readFileSync(path.join(root, 'docs', 'DELIVERY_DIRECTIVE_EN.md')))
    );
    check(
      'T28c OB-73 the live directive pair satisfies all six shapes (real whole-file files)',
      real.a2 === 'ABSENT' &&
        /"lines":true/.test(real.detail) &&
        /"head":true/.test(real.detail) &&
        /"bold":true/.test(real.detail) &&
        /"bars":true/.test(real.detail) &&
        /"blank":true/.test(real.detail) &&
        /"backticks":true/.test(real.detail),
      'a2=' + real.a2 + ' :: ' + real.detail
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
        /fixture-g8-bom-crlf\.md bom=1 crlf=true/.test(bomOutDetail) &&
        /fixture-g8-bom-nobom\.md bom=0 crlf=false/.test(bomOutDetail),
      'verdict=' + verdictOf(bomOut.records, 'G4a') + ' detail=' + bomOutDetail);

    // ---- T89: G4a counts CONSECUTIVE BOMs (OB-112). The old boolean test only looked at the
    // first three bytes, so a doubled BOM satisfied "has a BOM" and slipped through; the count
    // is now the load-bearing value. The probe is an in-memory Buffer, not a disk fixture.
    // (Id renumbered from T80 in 7 r1: T80-T83 already belong to the G9-c temp-artifact guard.)
    {
      const BOM3 = Buffer.from([0xef, 0xbb, 0xbf]);
      const body = Buffer.from('# t\n', 'utf8');
      const probe = (path, buf) => {
        const report = createReport();
        G4A.run(
          {
            changedPaths: [path],
            readBytes: () => buf,
            isIgnored: () => false,
            isEvidenceTree: () => false,
          },
          report
        );
        const rec = report.records.find((r) => r.id === 'G4a');
        return { verdict: rec.verdict, detail: rec.detail };
      };
      const one = probe('docs/t089/one.md', Buffer.concat([BOM3, body]));
      check('T89 G4a accepts a .md carrying exactly one BOM',
        one.verdict === 'PASS', 'verdict=' + one.verdict + ' detail=' + one.detail);
      const none = probe('docs/t089/none.md', body);
      check('T89 G4a fails a .md carrying zero BOMs and reports bom=0',
        none.verdict === 'FAIL' && /bom=0\b/.test(none.detail),
        'verdict=' + none.verdict + ' detail=' + none.detail);
      const two = probe('docs/t089/two.md', Buffer.concat([BOM3, BOM3, body]));
      check('T89 G4a fails a .md carrying two consecutive BOMs and reports bom=2',
        two.verdict === 'FAIL' && /bom=2\b/.test(two.detail),
        'verdict=' + two.verdict + ' detail=' + two.detail);
      const stray = probe('docs/t089/stray.txt', Buffer.concat([BOM3, Buffer.from('x\n', 'utf8')]));
      check('T89 G4a fails a non-.md file carrying one BOM and reports bom=1',
        stray.verdict === 'FAIL' && /bom=1\b/.test(stray.detail),
        'verdict=' + stray.verdict + ' detail=' + stray.detail);
      const widened = loadMutatedCriteria(G4A_SOURCE, 'bomCount === 1', 'bomCount >= 1', 't89-double-bom-widened');
      const widenedReport = createReport();
      widened.exports.run(
        {
          changedPaths: ['docs/t089/two.md'],
          readBytes: () => Buffer.concat([BOM3, BOM3, body]),
          isIgnored: () => false,
          isEvidenceTree: () => false,
        },
        widenedReport
      );
      const widenedRec = widenedReport.records.find((r) => r.id === 'G4a');
      check('T89 reverse mutation (bomCount >= 1) makes the double-BOM file pass again, so the equality is load-bearing',
        widenedRec.verdict === 'PASS' && widened.applied && widened.identical,
        'verdict=' + widenedRec.verdict + ' applied=' + widened.applied + ' identical=' + widened.identical);
    }

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

  // ---- T68-T70: OB-20 - a FAILED git grep must never read as "zero hits"
  // (slice GATES-TIGHTEN, v1.33). `git grep` exits 1 when there is NO match: that is data
  // (the gate is satisfied). Any OTHER non-zero exit is an execution / environment failure
  // and must surface as a UsageError (exit 2), never as a PASS indistinguishable from
  // "zero hits" - the last remnant of the vacuous-pass family ⑦ dug out three rounds running.
  {
    const g5a1Ctx = (gitResult) => ({
      git: () => gitResult,
      listAgentFiles: () => [],
      readText: () => null,
    });
    const g5a1Verdict = (gitResult) => {
      const out = createReport();
      G5A.subchecks[0].run(g5a1Ctx(gitResult), out);
      return verdictOf(out.records, 'G5-a(1)');
    };

    check('T68 OB-20: git grep exit 1 (no match) is data - G5-a(1) still PASSes',
      g5a1Verdict({ ok: false, code: 1, out: '', err: '', spawnError: null }) === 'PASS',
      'verdict=' + g5a1Verdict({ ok: false, code: 1, out: '', err: '', spawnError: null }));

    let failedGrep = null;
    try {
      g5a1Verdict({ ok: false, code: 128, out: '', err: 'fatal: not a git repository', spawnError: null });
    } catch (e) {
      failedGrep = e;
    }
    check('T69 OB-20: a failed git grep (exit 128) is a UsageError, never a silent PASS',
      failedGrep instanceof UsageError && /G5-a\(1\)/.test(failedGrep.message) && /128/.test(failedGrep.message),
      'err=' + (failedGrep ? failedGrep.message : 'NO ERROR (fail-open)'));

    let missingGit = null;
    try {
      g5a1Verdict({ ok: false, code: 127, out: '', err: '', spawnError: new Error('spawn git ENOENT') });
    } catch (e) {
      missingGit = e;
    }
    check('T70 OB-20: a missing git binary (exit 127) is a UsageError too',
      missingGit instanceof UsageError,
      'err=' + (missingGit ? missingGit.message : 'NO ERROR (fail-open)'));

    check('T70 OB-20 positive control: a real hit still FAILs (the guard rejects nothing legitimate)',
      g5a1Verdict({ ok: true, code: 0, out: '.github/agents/prfrail-tester.agent.md\n', err: '', spawnError: null }) === 'FAIL',
      'verdict=' + g5a1Verdict({ ok: true, code: 0, out: '.github/agents/prfrail-tester.agent.md\n', err: '', spawnError: null }));
  }

  // ---- T71: OB-17 - the angle-bracket rule is narrowed to id-LIKE placeholders
  // (slice GATES-TIGHTEN, v1.33). `run <RUN_ID>` still fails; the CLI usage example
  // `run <chain-file>` (token starts lowercase) must no longer be judged, so touching
  // such a line can no longer turn G6-3 red by itself (OB-17's residual false positive).
  {
    const usage = g63Probe('prfrail run <chain-file>');
    check('T71 OB-17: the CLI usage example `run <chain-file>` is no longer judged',
      usage.verdict === 'PASS' && usage.findings === 0,
      'verdict=' + usage.verdict + ' findings=' + usage.findings + ' why=' + usage.why);

    const placeholder = g63Probe('prfrail run <RUN_ID>');
    check('T71 OB-17: the id-like placeholder `run <RUN_ID>` still fails',
      placeholder.verdict === 'FAIL' && placeholder.findings === 1,
      'verdict=' + placeholder.verdict + ' findings=' + placeholder.findings);

    const shortId = g63Probe('prfrail run <ID>');
    check('T71 OB-17: a short uppercase placeholder `run <ID>` still fails',
      shortId.verdict === 'FAIL',
      'verdict=' + shortId.verdict + ' findings=' + shortId.findings);

    // Mutation: widening the rule back to `[^>]*` must make the usage example red again,
    // proving the narrowing (not a white-list entry) is what keeps the example green.
    const mut = loadMutatedG6(
      '/(?<![A-Za-z\\u4e00-\\u9fff])run\\s*<[A-Z][A-Z0-9_-]*>/g',
      '/(?<![A-Za-z\\u4e00-\\u9fff])run\\s*<[^>]*>/g',
      'ob17-angle-bracket');
    const mutOut = createReport();
    mut.exports.subchecks.find((s) => s.id === 'G6-3').run(probeCtx('prfrail run <chain-file>'), mutOut);
    check('T71 OB-17 the narrowing is load-bearing: widening it back turns the usage example red',
      mut.applied && verdictOf(mutOut.records, 'G6-3') === 'FAIL',
      'applied=' + mut.applied + ' verdict=' + verdictOf(mutOut.records, 'G6-3'));
    check('T71 OB-17 restore: g6.js is byte-identical after the mutation probe',
      mut.identical === true, 'identical=' + mut.identical);
  }

  // ---- T80-T83: G9 temp-artifact guard (slice G9-c, v1.38) -------------------------
  // Two predicates: (1) temp-artifact shape outside the tmp root; (2) the tmp-root
  // guard (case-insensitive, the forced-entry shape for the gitignored `tmp/` root;
  // force-added paths under OTHER gitignored roots are out of scope). The positive
  // fixtures pin all three exclusion faces (evidence tree / whitelist / deleted path),
  // the six negative fixtures (N1-N6) pin every must-FAIL shape, and the four
  // in-memory mutations pin each branch as load-bearing.
  {
    const G9_SHAPE_RE = 'const TEMP_SHAPE_RE = /\\.(log|tmp|bak|out)$/i;';
    const G9_BEFORE = fs.readFileSync(G9_SOURCE);

    const P1 = g9Probe(['docs/t027/probe-g9.md', 'tools/gates/g9-probe/main.go']);
    check(
      'T80 G9 passes ordinary documents and Go files (no temp-artifact shape)',
      P1.verdict === 'PASS' && P1.findings === 0,
      'verdict=' + P1.verdict + ' findings=' + P1.findings + ' :: ' + P1.detail
    );

    const P2 = g9Probe(['docs/validation/evidence/legacy.log']);
    check(
      'T80 G9 excludes the frozen evidence tree (OB-79): a legacy.log inside it is not judged',
      P2.verdict === 'PASS' && P2.findings === 0,
      'verdict=' + P2.verdict + ' findings=' + P2.findings + ' :: ' + P2.detail
    );

    const P3 = g9Probe(['tools/gates/g9-probe/allowed.log'], {
      whitelist: [
        {
          match: 'tools/gates/g9-probe/allowed.log',
          file: 'tools/gates/g9-probe/allowed.log',
          reason: 'T80 probe: registered temp-artifact exception',
        },
      ],
    });
    check(
      'T80 G9 honours a whitelist registration (match + file)',
      P3.verdict === 'PASS' && P3.findings === 0,
      'verdict=' + P3.verdict + ' findings=' + P3.findings + ' :: ' + P3.detail
    );

    const P4 = g9Probe(['tmp/.gitkeep']);
    check(
      'T80 G9 passes the tracked tmp directory placeholder tmp/.gitkeep',
      P4.verdict === 'PASS' && P4.findings === 0,
      'verdict=' + P4.verdict + ' findings=' + P4.findings + ' :: ' + P4.detail
    );

    const D1 = g9Probe(['docs/gone/selftest-output.log'], { deleted: ['docs/gone/selftest-output.log'] });
    check(
      'T80 G9 excludes deleted paths (a deletion submits no bytes)',
      D1.verdict === 'PASS' && D1.findings === 0,
      'verdict=' + D1.verdict + ' findings=' + D1.findings + ' :: ' + D1.detail
    );

    const N1 = g9Probe(['tools/foo/selftest-output.log']);
    check(
      'T81 G9 fails a .log outside the tmp root (token `temp-shape`, names the path)',
      N1.verdict === 'FAIL' && N1.findings === 1 && /temp-shape/.test(N1.detail) && /selftest-output\.log/.test(N1.detail),
      'verdict=' + N1.verdict + ' findings=' + N1.findings + ' :: ' + N1.detail
    );

    const N2 = g9Probe(['cmd/prfrail/foo.tmp']);
    check(
      'T81 G9 fails a .tmp outside the tmp root (token `temp-shape`, names the path)',
      N2.verdict === 'FAIL' && N2.findings === 1 && /temp-shape/.test(N2.detail) && /foo\.tmp/.test(N2.detail),
      'verdict=' + N2.verdict + ' findings=' + N2.findings + ' :: ' + N2.detail
    );

    const N3 = g9Probe(['docs/backup.bak']);
    check(
      'T81 G9 fails a .bak outside the tmp root (token `temp-shape`, names the path)',
      N3.verdict === 'FAIL' && N3.findings === 1 && /temp-shape/.test(N3.detail) && /backup\.bak/.test(N3.detail),
      'verdict=' + N3.verdict + ' findings=' + N3.findings + ' :: ' + N3.detail
    );

    const N4 = g9Probe(['tmp/x.log']);
    check(
      'T81 G9 fails any tmp-root path except tmp/.gitkeep (token `tmp-root`, names the path)',
      N4.verdict === 'FAIL' && N4.findings === 1 && /tmp-root/.test(N4.detail) && /tmp\/x\.log/.test(N4.detail),
      'verdict=' + N4.verdict + ' findings=' + N4.findings + ' :: ' + N4.detail
    );

    const N5 = g9Probe(['Tmp/g9-probe.txt']);
    check(
      'T81 G9 fails a mixed-case tmp root (token `tmp-root`, original spelling kept)',
      N5.verdict === 'FAIL' && N5.findings === 1 && N5.detail.indexOf('["tmp-root: Tmp/g9-probe.txt"]') >= 0,
      'verdict=' + N5.verdict + ' findings=' + N5.findings + ' :: ' + N5.detail
    );

    const N6 = g9Probe(['TMP/nested/x.txt']);
    check(
      'T81 G9 fails an upper-case tmp root (token `tmp-root`, original spelling kept)',
      N6.verdict === 'FAIL' && N6.findings === 1 && N6.detail.indexOf('["tmp-root: TMP/nested/x.txt"]') >= 0,
      'verdict=' + N6.verdict + ' findings=' + N6.findings + ' :: ' + N6.detail
    );

    // M1: replace the shape regex with a never-matching one. N1/N2/N3 (all outside the
    // tmp root) must flip to PASS, proving the regex - not another branch - catches them.
    const m1 = loadMutatedModule(G9_SOURCE, G9_SHAPE_RE, 'const TEMP_SHAPE_RE = /(?!x)x/;', 'g9-shape-regex');
    const m1n1 = g9Probe(['tools/foo/selftest-output.log'], { mod: m1.exports });
    const m1n2 = g9Probe(['cmd/prfrail/foo.tmp'], { mod: m1.exports });
    const m1n3 = g9Probe(['docs/backup.bak'], { mod: m1.exports });
    check(
      'T82 M1 the shape regex is load-bearing: /(?!x)x/ turns N1/N2/N3 green',
      m1.applied === true &&
        m1.identical === true &&
        [m1n1, m1n2, m1n3].every((p) => p.verdict === 'PASS' && p.findings === 0),
      'applied=' +
        m1.applied +
        ' identical=' +
        m1.identical +
        ' verdicts=' +
        JSON.stringify([m1n1.verdict, m1n2.verdict, m1n3.verdict])
    );

    // M2: delete the whitelist exemption branch. P3 must flip to FAIL.
    const m2 = loadMutatedModule(
      G9_SOURCE,
      '      if (isWhitelisted(entries, p)) {\n        excluded++;\n        continue;\n      }\n',
      '',
      'g9-whitelist-branch'
    );
    const m2p3 = g9Probe(['tools/gates/g9-probe/allowed.log'], {
      mod: m2.exports,
      whitelist: [
        {
          match: 'tools/gates/g9-probe/allowed.log',
          file: 'tools/gates/g9-probe/allowed.log',
          reason: 'T82 probe: registration must be load-bearing',
        },
      ],
    });
    check(
      'T82 M2 the whitelist exemption is load-bearing: deleting it turns P3 red',
      m2.applied === true && m2.identical === true && m2p3.verdict === 'FAIL' && m2p3.findings === 1,
      'applied=' + m2.applied + ' identical=' + m2.identical + ' verdict=' + m2p3.verdict + ' :: ' + m2p3.detail
    );

    // M3: delete the evidence-tree exclusion. P2 must flip to FAIL.
    const m3 = loadMutatedModule(
      G9_SOURCE,
      '      if (inTree(p)) {\n        excluded++;\n        continue;\n      }\n',
      '',
      'g9-evidence-branch'
    );
    const m3p2 = g9Probe(['docs/validation/evidence/legacy.log'], { mod: m3.exports });
    check(
      'T82 M3 the evidence-tree exclusion is load-bearing: deleting it turns P2 red',
      m3.applied === true && m3.identical === true && m3p2.verdict === 'FAIL' && m3p2.findings === 1,
      'applied=' + m3.applied + ' identical=' + m3.identical + ' verdict=' + m3p2.verdict + ' :: ' + m3p2.detail
    );

    // M4: turn the `!== tmp/.gitkeep` guard into an unconditional skip. N4 must flip to
    // PASS, proving the root guard (not some other branch) catches `tmp/x.log`.
    const m4 = loadMutatedModule(G9_SOURCE, "if (key !== TMP_ROOT_KEEP) bad.push('tmp-root: ' + p);", '/* unconditional skip */', 'g9-tmpkeep-guard');
    const m4n4 = g9Probe(['tmp/x.log'], { mod: m4.exports });
    check(
      'T82 M4 the tmp/.gitkeep guard is load-bearing: an unconditional skip turns N4 green',
      m4.applied === true && m4.identical === true && m4n4.verdict === 'PASS' && m4n4.findings === 0,
      'applied=' + m4.applied + ' identical=' + m4.identical + ' verdict=' + m4n4.verdict + ' :: ' + m4n4.detail
    );

    // Restore anchors: the shipped module is untouched and still judges the same probes.
    check(
      'T83 G9 restore: g9.js is byte-identical and the cached module keeps N1 red / P2 green / P4 green',
      fs.readFileSync(G9_SOURCE).equals(G9_BEFORE) &&
        g9Probe(['tools/foo/selftest-output.log']).verdict === 'FAIL' &&
        g9Probe(['docs/validation/evidence/legacy.log']).verdict === 'PASS' &&
        g9Probe(['tmp/.gitkeep']).verdict === 'PASS',
      JSON.stringify({
        n1: g9Probe(['tools/foo/selftest-output.log']).verdict,
        p2: g9Probe(['docs/validation/evidence/legacy.log']).verdict,
        p4: g9Probe(['tmp/.gitkeep']).verdict,
      })
    );
  }

  // ---- T84: OB-97 (v1.40) - the language-tag continuation set includes `:`, so the
  // second colon in `lang:zh:en` / the trailing colon in `lang:zh:` no longer truncates
  // the token into the legal prefix `lang:zh`.
  {
    const CN = 'docs/t027/probe-g6-6-t84.md';
    const ROW = '### 12.4 观察项台账\n\n| OB-1 | d | x | 观察中 |\n';
    const gi = G6.internals;
    const run = (head, file) => {
      try {
        return { res: g66Probe(file || CN, head + ROW), err: null };
      } catch (e) {
        return { res: null, err: e };
      }
    };
    const unknown = (o, tag) =>
      !!o.err && o.err instanceof UsageError && /unknown language tag/.test(o.err.message) && o.err.message.indexOf('`' + tag + '`') >= 0;
    const why = (o) => (o.err ? 'err=' + o.err.message : 'verdict=' + (o.res ? o.res.verdict : 'NONE'));

    const normal = gi.declaredLang(['`lang:zh`'], CN);
    check(
      'T84 OB-97 positive: the widened continuation set leaves the ordinary declaration intact (`lang:zh` still declares zh)',
      normal === 'cn',
      'declared=' + normal
    );

    const terminator = run('`lang:zh`，zh\n');
    check(
      'T84 OB-97 positive: a fullwidth comma still terminates the token, so `lang:zh，` stays legal',
      !terminator.err && !!terminator.res && terminator.res.verdict === 'PASS' && terminator.res.findings === 0,
      why(terminator)
    );

    const colonEn = run('`lang:zh:en`\n');
    check(
      'T84 OB-97 negative: `lang:zh:en` is consumed whole and rejected (it used to truncate to the legal `lang:zh`)',
      unknown(colonEn, 'lang:zh:en'),
      why(colonEn)
    );

    const colonBare = run('`lang:zh:`\n');
    check(
      'T84 OB-97 negative: a trailing colon belongs to the token, so `lang:zh:` is rejected',
      unknown(colonBare, 'lang:zh:'),
      why(colonBare)
    );

    const coexist = run(LANG_TAG_ZH + '`lang:zh:`\n');
    check(
      'T84 OB-97 negative: coexistence (`lang:zh` plus `lang:zh:`) is fail-closed',
      unknown(coexist, 'lang:zh:'),
      why(coexist)
    );

    const mut = loadMutatedModule(G6_SOURCE, G6_LANG_TAG_CONT, 'const LANG_TAG_CONT_RE = /[A-Za-z0-9._-]/;', 'ob97-colon');
    const blind = run('`lang:zh:en`\n');
    const blindMut = (() => {
      try {
        return { res: g66Probe(CN, '`lang:zh:en`\n' + ROW, mut.exports), err: null };
      } catch (e) {
        return { res: null, err: e };
      }
    })();
    check(
      'T84 OB-97 mutation applied + load-bearing: with the pre-OB-97 charset the SAME `lang:zh:en` header passes again',
      mut.applied === true && !blindMut.err && blindMut.res.verdict === 'PASS',
      'applied=' + mut.applied + ' ' + why(blindMut)
    );

    check(
      'T84 OB-97 restore: g6.js is byte-identical after the mutation window and the cached module rejects `lang:zh:en` again',
      mut.identical === true && unknown(blind, 'lang:zh:en'),
      'identical=' + mut.identical + ' ' + why(blind)
    );
  }

  // ---- T85: G6-7 (OB-32) - a declaration on an added §0.3 changelog line must name an
  // ID that exists in the same document's §12.4 / appendix-C ledger. The narrowing to
  // §0.3-only targets and whole-ledger evidence is the prerequisite-gate-3 result (63
  // historical false positives under the strict "touched this change set" form, 0 after
  // the narrowing). Four positives pin the two evidence faces, the exclusion markers and
  // the region boundary; the negative pins the missing ID; two mutations prove the verb
  // test and the evidence test are load-bearing; two boundary checks pin the token-boundary
  // lookbehind (v1.40: `ADR-027` must not be carved into a `DR-027` declaration).
  {
    const P = 'docs/t027/probe-g6-7.md';
    const doc = (o) =>
      [
        '### 0.3 变更日志',
        '',
        '| 版本 | 日期 | 变更 | 受影响的在飞切片 |',
        '|---|---|---|---|',
        o.log,
        '',
        '## 1. body',
        '',
        o.body || '(none)',
        '',
        '### 12.4 观察项台账',
        '',
        '| # | 登记日 | 现象 | 处置状态 |',
        '|---|---|---|---|',
        o.ledger,
        '',
        '### 12.5 缺口与未来扩展',
        '',
        '(none)',
        '',
        '### C.0 slice',
        '',
        o.apx || '- 交付物：none',
        '',
        '### C.1 compare',
        '',
        '(none)',
      ].join('\n');

    const p1 = doc({ log: '| v1.1 | 2026-01-01 | 登记 OB-901 as observed |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |' });
    const r1 = g67Probe(P, p1);
    check(
      'T85 G6-7 positive: a §0.3 declaration whose ID exists in the §12.4 ledger passes',
      r1.verdict === 'PASS' && r1.findings === 0,
      'verdict=' + r1.verdict + ' findings=' + r1.findings + ' :: ' + r1.detail
    );

    const p2 = doc({ log: '| v1.1 | 2026-01-01 | SW-901 落地 |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |', apx: '- 交付物：SW-901' });
    const r2 = g67Probe(P, p2, null, { added: [1, 2, 3, 4, 5] });
    check(
      'T85 G6-7 positive: a PRE-EXISTING (not added) appendix-C ledger line still satisfies the declaration',
      r2.verdict === 'PASS' && r2.findings === 0,
      'verdict=' + r2.verdict + ' findings=' + r2.findings + ' :: ' + r2.detail
    );

    const p3 = doc({ log: '| v1.1 | 2026-01-01 | 示例：登记 OB-902 |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |' });
    const r3 = g67Probe(P, p3);
    const p3neg = doc({ log: '| v1.1 | 2026-01-01 | 本片不新增、不删除；OB-905 保持待议 |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |' });
    const r3neg = g67Probe(P, p3neg);
    check(
      'T85 G6-7 positive: an exclusion-marker line and a NEGATED verb (`不新增 … OB-905`) are not judged',
      r3.verdict === 'PASS' && r3.findings === 0 && r3neg.verdict === 'PASS' && r3neg.findings === 0,
      'marker=' + r3.verdict + '/' + r3.findings + ' negated=' + r3neg.verdict + '/' + r3neg.findings
    );

    const p4 = doc({ log: '| v1.1 | 2026-01-01 | nothing |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |', body: '登记 OB-903 in the body' });
    const r4 = g67Probe(P, p4);
    check(
      'T85 G6-7 positive: a declaration outside the §0.3 region is not judged',
      r4.verdict === 'PASS' && r4.findings === 0,
      'verdict=' + r4.verdict + ' findings=' + r4.findings + ' :: ' + r4.detail
    );

    const n1 = doc({ log: '| v1.1 | 2026-01-01 | 登记 OB-904 as observed |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |' });
    const rn = g67Probe(P, n1);
    check(
      'T85 G6-7 negative: a declared ID with no ledger row anywhere fails and the finding names the ID',
      rn.verdict === 'FAIL' && rn.findings === 1 && rn.detail.indexOf('OB-904') >= 0,
      'verdict=' + rn.verdict + ' findings=' + rn.findings + ' :: ' + rn.detail
    );

    const m1 = loadMutatedModule(G6_SOURCE, '      verb = vm[0];', '      verb = null;', 'g67-verb-test');
    const rm1 = g67Probe(P, n1, m1.exports);
    check(
      'T85 G6-7 mutation applied + load-bearing: disabling the verb test turns the same negative case green',
      m1.applied === true && m1.identical === true && rm1.verdict === 'PASS' && rm1.findings === 0,
      'applied=' + m1.applied + ' identical=' + m1.identical + ' verdict=' + rm1.verdict + ' findings=' + rm1.findings
    );

    const m2 = loadMutatedModule(G6_SOURCE, '          if (evidence.some((ln) => ln !== n)) continue;', '          if (true) continue;', 'g67-ledger-test');
    const rm2 = g67Probe(P, n1, m2.exports);
    check(
      'T85 G6-7 mutation applied + load-bearing: ignoring the ledger evidence turns the same negative case green',
      m2.applied === true && m2.identical === true && rm2.verdict === 'PASS' && rm2.findings === 0,
      'applied=' + m2.applied + ' identical=' + m2.identical + ' verdict=' + rm2.verdict + ' findings=' + rm2.findings
    );

    check(
      'T85 G6-7 restore: g6.js is byte-identical after both mutation windows and the shipped module still fails the same negative case',
      m1.identical === true && m2.identical === true && g67Probe(P, n1).verdict === 'FAIL',
      'identical=' + (m1.identical && m2.identical) + ' verdict=' + g67Probe(P, n1).verdict
    );

    // Boundary (v1.40). A longer token must not be carved into an ID key: with the
    // boundary-less form `ADR-027` yielded a `DR-027` declaration with no ledger row, so
    // this PASS is exactly the payload the fix changes (the old form FAILs with 1 finding).
    const b1 = doc({ log: '| v1.1 | 2026-01-01 | 登记 ADR-027 |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |' });
    const rb1 = g67Probe(P, b1);
    check(
      'T85 G6-7 boundary: `ADR-027` is not a `DR-027` declaration and the line is not judged',
      rb1.verdict === 'PASS' && rb1.findings === 0,
      'verdict=' + rb1.verdict + ' findings=' + rb1.findings + ' :: ' + rb1.detail
    );

    // The boundary must not over-narrow: the real ID on the same line is still declared
    // exactly once (the old form produced TWO findings, `DR-028` and `OB-103`).
    const b2 = doc({ log: '| v1.1 | 2026-01-01 | 登记 ADR-028 与 OB-103 |', ledger: '| OB-901 | 2026-01-01 | x | 观察中 |' });
    const rb2 = g67Probe(P, b2);
    check(
      'T85 G6-7 boundary: `登记 ADR-028 与 OB-103` yields exactly one declaration (OB-103)',
      rb2.verdict === 'FAIL' && rb2.findings === 1 && rb2.detail.indexOf('OB-103') >= 0,
      'verdict=' + rb2.verdict + ' findings=' + rb2.findings + ' :: ' + rb2.detail
    );
  }

  // ---- T86: G5-d (OB-30) - the section 2.6 five-key frontmatter contract on changed
  // role files. Three positives (complete file, extra keys, non-target path), three
  // negatives (missing / duplicate / empty-or-absent frontmatter), one mutation and one
  // restore.
  {
    const AG = '.github/agents/probe-g5d.agent.md';
    const role = (o) =>
      [
        '---',
        'name: ' + (o.name === undefined ? 'probe-role' : o.name),
        'description: ' + (o.description === undefined ? '"probe role"' : o.description),
        'tools: ' + (o.tools === undefined ? "['read', 'search']" : o.tools),
        'user-invocable: ' + (o.invocable === undefined ? 'false' : o.invocable),
        ...(o.target === null ? [] : ['target: ' + (o.target === undefined ? 'vscode' : o.target)]),
        '---',
        '',
        '# body',
      ].join('\n');

    const ok = role({});
    const p1 = g5dProbe({ [AG]: ok });
    check(
      'T86 G5-d positive: a role file carrying the five §2.6 keys exactly once passes',
      p1.verdict === 'PASS' && p1.findings === 0,
      'verdict=' + p1.verdict + ' findings=' + p1.findings + ' :: ' + p1.detail
    );

    const extra = ok.replace('\n---\n\n# body', "\nmodel: 'Some Model (vendor)'\n---\n\n# body");
    const p2 = g5dProbe({ [AG]: extra });
    check(
      'T86 G5-d positive: an EXTRA key is not judged (model stays with G5-a / G5-b)',
      p2.verdict === 'PASS' && p2.findings === 0,
      'verdict=' + p2.verdict + ' findings=' + p2.findings + ' :: ' + p2.detail
    );

    const p3 = g5dProbe({ 'docs/probe-g5d.md': '# no frontmatter at all' });
    check(
      'T86 G5-d positive: a non-agent path is not a target (INFO, never a vacuous PASS claim)',
      p3.verdict === 'PASS' && /no changed/.test(p3.detail),
      'verdict=' + p3.verdict + ' :: ' + p3.detail
    );

    const n1text = role({ target: null });
    const n1 = g5dProbe({ [AG]: n1text });
    check(
      'T86 G5-d negative: a missing `target` key fails and the finding names it',
      n1.verdict === 'FAIL' && n1.findings === 1 && n1.detail.indexOf('target') >= 0,
      'verdict=' + n1.verdict + ' findings=' + n1.findings + ' :: ' + n1.detail
    );

    const n2text = role({}).replace("tools: ['read', 'search']", "tools: ['read', 'search']\ntools: ['read']");
    const n2 = g5dProbe({ [AG]: n2text });
    check(
      'T86 G5-d negative: a duplicated `tools` key fails (exactly once required)',
      n2.verdict === 'FAIL' && n2.findings === 1 && /occurs 2 times/.test(n2.detail),
      'verdict=' + n2.verdict + ' findings=' + n2.findings + ' :: ' + n2.detail
    );

    const n3empty = g5dProbe({ [AG]: role({ description: '""' }) });
    const n3absent = g5dProbe({ [AG]: '# body without a fence' });
    check(
      'T86 G5-d negative: an empty value and an absent frontmatter both fail',
      n3empty.verdict === 'FAIL' &&
        /empty value/.test(n3empty.detail) &&
        n3absent.verdict === 'FAIL' &&
        /frontmatter missing/.test(n3absent.detail),
      'empty=' + n3empty.detail + ' absent=' + n3absent.detail
    );

    const m1 = loadMutatedCriteria(
      G5D_SOURCE,
      "const REQUIRED_KEYS = ['name', 'description', 'tools', 'user-invocable', 'target'];",
      "const REQUIRED_KEYS = ['name', 'description', 'tools', 'user-invocable'];",
      'g5d-required-keys'
    );
    const rm1 = g5dProbe({ [AG]: n1text }, { mod: m1.exports });
    check(
      'T86 G5-d mutation applied + load-bearing: removing `target` from the required set turns the missing-target case green',
      m1.applied === true && m1.identical === true && rm1.verdict === 'PASS' && rm1.findings === 0,
      'applied=' + m1.applied + ' identical=' + m1.identical + ' verdict=' + rm1.verdict + ' findings=' + rm1.findings
    );

    check(
      'T86 G5-d restore: g5d.js is byte-identical and the shipped module still fails the missing-target case',
      m1.identical === true && g5dProbe({ [AG]: n1text }).verdict === 'FAIL',
      'verdict=' + g5dProbe({ [AG]: n1text }).verdict
    );
  }

  // ---- T87: G9-e (OB-77) - every `tmp/`-prefixed reference on an added line must be
  // resolvable; an unresolvable one is INFO, never red. Five positives pin the two
  // resolver branches, the bare `tmp/` boundary, the exclusion markers and the evidence
  // tree; the negative pins the informational (non-red) record; one mutation and one
  // restore prove the resolver test is load-bearing.
  {
    const P = 'docs/t027/probe-g9e.md';
    const repoRoot = discoverRepoRoot(__dirname);
    const treeResolve = (p) => resolveRepoArtifact(repoRoot, 'tree', null, p);
    const headResolve = (p) => resolveRepoArtifact(repoRoot, 'range', 'HEAD', p);
    const summaryOf = (probe) => {
      const r = probe.infos.find((x) => /summary/.test(x.label));
      return r ? r.detail : 'NO SUMMARY';
    };
    const mk = (text, opts) =>
      g9Probe([P], Object.assign({ markdown: [P], texts: new Map([[P, text]]), resolve: treeResolve }, opts || {}));

    const p1 = mk('see `tmp/.gitkeep` here');
    check(
      'T87 G9-e positive: a tracked tmp/ artifact resolves in tree scope (disk + index)',
      p1.verdict === 'PASS' && /0 unresolvable$/.test(summaryOf(p1)),
      'verdict=' + p1.verdict + ' :: ' + summaryOf(p1)
    );

    const p2 = mk('see tmp/ only');
    check(
      'T87 G9-e positive: a bare `tmp/` has no path segment and is not judged',
      p2.verdict === 'PASS' && /0 reference\(s\)/.test(summaryOf(p2)),
      'verdict=' + p2.verdict + ' :: ' + summaryOf(p2)
    );

    const p3 = mk('示例：tmp/ghost/x.log');
    check(
      'T87 G9-e positive: an exclusion-marker line is skipped',
      p3.verdict === 'PASS' && /0 reference\(s\)/.test(summaryOf(p3)),
      'verdict=' + p3.verdict + ' :: ' + summaryOf(p3)
    );

    const E = 'docs/validation/evidence/probe-g9e/NOTE.md';
    const p4 = g9Probe([E], { markdown: [E], texts: new Map([[E, 'tmp/ghost/x.log']]), resolve: treeResolve });
    check(
      'T87 G9-e positive: the frozen evidence tree is excluded from the reference scan',
      p4.verdict === 'PASS' && /0 reference\(s\)/.test(summaryOf(p4)),
      'verdict=' + p4.verdict + ' :: ' + summaryOf(p4)
    );

    const p5 = mk('see `tmp/.gitkeep` at the endpoint', { resolve: headResolve });
    check(
      'T87 G9-e positive: the range/ci branch resolves at the endpoint revision (HEAD:tmp/.gitkeep)',
      p5.verdict === 'PASS' && /0 unresolvable$/.test(summaryOf(p5)),
      'verdict=' + p5.verdict + ' :: ' + summaryOf(p5)
    );

    const pn = mk('see tmp/ghost/x.log here');
    check(
      'T87 G9-e negative: an unresolvable reference is INFO (G9 stays PASS) and the record names the path',
      pn.verdict === 'PASS' &&
        pn.infos.some((r) => r.detail.indexOf('tmp/ghost/x.log') >= 0) &&
        /1 unresolvable$/.test(summaryOf(pn)),
      'verdict=' + pn.verdict + ' :: ' + summaryOf(pn) + ' :: ' + pn.infos.map((r) => r.detail).join(' | ')
    );

    const m1 = loadMutatedModule(G9_SOURCE, '        if (r && r.ok) continue;', '        if (true) continue;', 'g9e-resolver-test');
    const pm1 = g9Probe([P], {
      mod: m1.exports,
      markdown: [P],
      texts: new Map([[P, 'see tmp/ghost/x.log here']]),
      resolve: treeResolve,
    });
    check(
      'T87 G9-e mutation applied + load-bearing: accepting every reference makes the informational record disappear',
      m1.applied === true && m1.identical === true && !pm1.infos.some((r) => r.detail.indexOf('tmp/ghost/x.log') >= 0),
      'applied=' + m1.applied + ' identical=' + m1.identical + ' :: ' + pm1.infos.map((r) => r.detail).join(' | ')
    );

    const pRestored = mk('see tmp/ghost/x.log here');
    check(
      'T87 G9-e restore: g9.js is byte-identical and the shipped module emits the informational record again',
      m1.identical === true && pRestored.infos.some((r) => r.detail.indexOf('tmp/ghost/x.log') >= 0),
      'identical=' + m1.identical + ' :: ' + summaryOf(pRestored)
    );
  }

  // ---- T88: G3(a.2) seventh/eighth shapes are RECORD-ONLY (prerequisite gates 1 and 2).
  // Both shapes were downgraded: the CN directive pair carries 16 deliberate CJK-space
  // forms (section labels, term separation) and historical bilingual commits break the
  // per-line triplet multiset (measured on c6ada79 / 2d6a8e4). These three assertions pin
  // the decision mechanically: the six implemented shapes still catch their own drift,
  // while a seventh-shape-only or eighth-shape-only drift stays PASS.
  {
    const CNP = 'docs/DELIVERY_DIRECTIVE.md';
    const ENP = 'docs/DELIVERY_DIRECTIVE_EN.md';
    const base = {
      key: CNP,
      changed: [CNP, ENP],
      numstat: { [CNP]: [3, 0], [ENP]: [3, 0] },
      existing: [],
    };
    const same = '# t\n\n**a**\n';

    const control = g3aEdgeProbe(
      Object.assign({}, base, { texts: { [CNP]: same.split('\n'), [ENP]: same.split('\n') } })
    );
    check(
      'T88 G3(a.2) control: an equal synthetic directive pair stays green (no per-pair FAIL record)',
      control.a === 'PASS' && control.a2 === 'ABSENT',
      'a=' + control.a + ' a2=' + control.a2 + ' :: ' + control.detail
    );

    const seventh = g3aEdgeProbe(
      Object.assign({}, base, {
        texts: { [CNP]: '# t\n\n中 文\n'.split('\n'), [ENP]: '# t\n\n中文\n'.split('\n') },
      })
    );
    check(
      'T88 G3(a.2) seventh shape (CJK inter-character space) is record-only: a drift with all six shapes equal stays PASS',
      seventh.a === 'PASS' && seventh.a2 === 'ABSENT',
      'a=' + seventh.a + ' a2=' + seventh.a2
    );

    const eighth = g3aEdgeProbe(
      Object.assign({}, base, {
        texts: {
          [CNP]: '# t\n\n**a**\n**b**\n'.split('\n'),
          [ENP]: '# t\n\na\n****b****\n'.split('\n'),
        },
      })
    );
    check(
      'T88 G3(a.2) eighth shape (per-line triplet multiset) is record-only: a compensating drift with equal totals stays PASS',
      eighth.a === 'PASS' && eighth.a2 === 'ABSENT',
      'a=' + eighth.a + ' a2=' + eighth.a2
    );
  }

  // ---- T90: the concurrent scheduler keeps report order while letting a criterion await
  // (OB-113 sub-item 2). It runs GATE.runCriteriaConcurrent itself, so a regression inside
  // gate.js - not a local re-derivation - is what turns these assertions red. It also pins the
  // --skip-base guard order (7 r1 finding 1: the guard must precede the --help early return).
  // (Id renumbered from T81 in 7 r1: T80-T83 already belong to the G9-c temp-artifact guard.)
  {
    const GATE = gateApi();
    const helpOnly = await GATE.main(['node', 'gate', '--help']);
    check('T90 --help on its own still exits 0 with the usage text',
      helpOnly.code === 0 && typeof helpOnly.stdout === 'string' && helpOnly.stdout.length > 0,
      'code=' + helpOnly.code + ' stdout_len=' + (helpOnly.stdout || '').length);
    const helpSkip = await GATE.main(['node', 'gate', '--help', '--skip-base']);
    check('T90 the --skip-base guard fires before the --help early return (no combination escapes it)',
      helpSkip.code === 2 && /--skip-base is only valid together with --all/.test(helpSkip.error || ''),
      'code=' + helpSkip.code + ' error=' + JSON.stringify(helpSkip.error || '').slice(0, 90));
    let inFlight = 0;
    let maxInFlight = 0;
    // Completion order, not wall clock: a timing threshold is not portable across CI runners
    // (DR-1: the Windows runner once needed 612 ms for a 0 ms + 150 ms + 150 ms plan, so the
    // old `elapsed < 280` assertion went red on unchanged code). The order below is what
    // distinguishes concurrency from a serial loop: under a serial loop C cannot END before A
    // ends, because C is not even started until A has resolved.
    const events = [];
    const delayed = (id, ids, delayMs) => ({
      id,
      label: id,
      run: (ctx, out) =>
        new Promise((resolve) => {
          inFlight++;
          maxInFlight = Math.max(maxInFlight, inFlight);
          events.push(id + '-start');
          setTimeout(() => {
            inFlight--;
            for (const r of ids) out.rec(r, 'PASS', r, id);
            events.push(id + '-end');
            resolve();
          }, delayMs);
        }),
    });
    // Synchronous twin of `delayed`, used by the reject-path assertion below.
    const mkSync = (id) => ({
      id,
      label: id,
      run: (ctx, out) => out.rec(id + '-rec', 'PASS', id + '-rec', 'sync'),
    });
    const asyncMap = {
      A: delayed('A', ['A-1', 'A-2'], 150),
      B: delayed('B', ['B-1'], 150),
      C: delayed('C', ['C-1'], 0),
    };
    const asyncOut = createReport();
    await GATE.runCriteriaConcurrent(['A', 'B', 'C'], asyncMap, {}, asyncOut);
    check('T90 concurrent scheduling keeps every record in plan order',
      asyncOut.records.map((r) => r.id).join(',') === 'A-1,A-2,B-1,C-1',
      'order=' + asyncOut.records.map((r) => r.id).join(','));
    check('T90 the 0 ms criterion ends before the 150 ms one that started first (completion order proves overlap)',
      events.indexOf('C-end') >= 0 &&
        events.indexOf('C-end') < events.indexOf('A-end') &&
        events.indexOf('B-start') < events.indexOf('A-end'),
      'events=' + events.join(','));
    check('T90 two async criteria are pending at the same moment (impossible under a serial loop)',
      maxInFlight >= 2, 'maxInFlight=' + maxInFlight);
    const syncMap = {
      S1: { id: 'S1', label: 'S1', run: (ctx, out) => out.rec('S1-rec', 'PASS', 'S1-rec', 'sync') },
      S2: { id: 'S2', label: 'S2', run: (ctx, out) => out.rec('S2-rec', 'PASS', 'S2-rec', 'sync') },
    };
    const syncOut = createReport();
    await GATE.runCriteriaConcurrent(['S2', 'S1'], syncMap, {}, syncOut);
    check('T90 synchronous criteria go through the same scheduler and keep plan order',
      syncOut.records.map((r) => r.id).join(',') === 'S2-rec,S1-rec',
      'order=' + syncOut.records.map((r) => r.id).join(','));
    // T90: a rejecting criterion rejects the whole plan run and leaks no records.
    let t90Leaked = 0;
    const t90Bad = { id: 'BAD', run: () => Promise.reject(new Error('boom')) };
    const t90Out = createReport();
    let t90Rejected = false;
    try {
      await GATE.runCriteriaConcurrent(['OK', 'BAD'], { OK: mkSync('OK'), BAD: t90Bad }, {}, t90Out);
    } catch (e) {
      t90Rejected = true;
    }
    t90Leaked = t90Out.records.length;
    check('T90 a rejecting criterion rejects the plan and writes no partial records',
      t90Rejected === true && t90Leaked === 0, 'rejected=' + t90Rejected + ' records=' + t90Leaked);
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

async function run() {
  try {
    return await runTests();
  } catch (e) {
    return { ok: false, stdout: 'SELFTEST: FAIL 0/0 :: ' + (e && e.stack ? e.stack : e), passed: 0, total: 0 };
  }
}

if (require.main === module) {
  run().then((res) => {
    process.stdout.write(res.stdout + '\n');
    process.exitCode = res.ok ? 0 : 1;
  });
}

module.exports = { run, FIXTURES, EXPECTED_CORPUS_FAILS };

'use strict';

/**
 * Base criterion G5 (DELIVERY_DIRECTIVE 6.3): change-set contract.
 *
 * Before v1.36 the row demanded a structured JSON produced OUTSIDE the repository
 * (`{stagedFiles[], declaredFiles[], evidenceRefs[], addMode}`) and the runner printed
 * SKIP. The independent ① design of slice GATES-G1G5 confirmed the flaw: whoever writes
 * that JSON is the party being audited, so the criterion would only re-check a
 * self-report. v1.36 derives every field from an independent source instead:
 *
 *   stagedFiles   = the scope's change set, deleted paths excluded     -> git fact
 *   declaredFiles = `artifacts` fence rows of the changed reports, `repo` disposition
 *                   only, plus the report's own path                 -> in-repo block
 *   evidenceRefs  = paths in the evidence cell of ADDED ledger completion rows
 *                   (backticked OR bare - the current corpus uses bare paths)
 *   addMode       = a DERIVED LABEL, never command-history evidence (git does not record
 *                   how a path was staged; `-A` and per-file adds leave a byte-identical
 *                   index). Its observable consequences are what (d) and (f) judge.
 *
 * Sub-checks:
 *   G5(d) staged ⊆ declared  - a staged/committed path nobody declared fails; the
 *         reverse direction (declared but untouched) is INFO, because a report fence is
 *         the snapshot of its landing commit: the §11.2 docs-only backfill commits
 *         legitimately change 2 files against a 11-row fence. An equality test here
 *         would turn every planned backfill red in CI.
 *   G5(e) every evidence reference resolves (file or directory, scope-aware).
 *   G5(f) index scope only: no worktree edit and no untracked path (`status --porcelain
 *         -uall`; a staged entry is the normal state of index scope and never dirties it).
 *
 * Exit codes: set inequality / unresolved evidence / dirty index ⇒ FAIL (1); a missing
 * anchor, an unreadable fence or an inconsistent `GATES_CHANGESET_JSON` ⇒ usage (2).
 */

const { UsageError } = require('../ctx');

const FENCE = /^```artifacts\s*$/;
const LEDGER_PAIR = ['docs/t027/REMAINING_SLICES.md', 'docs/t027/REMAINING_SLICES_EN.md'];
const REPORT_RE = /^docs\/validation\/[^/]+\.md$/i;
const PATH_TOKEN = /^(?:docs|tools|examples|schemas|cmd|internal|harnesses|testdata|\.github)\/[^\s；;、,|)`]+/;
// A repository artifact may also sit at the repository ROOT (`README.md`, `go.mod`,
// `CHANGELOG.md`). The prefix whitelist alone silently skipped those, which made G5(e)
// vacuous for them; a bare token is accepted only when it is a single path component with
// a known artifact extension, so ordinary prose (`v1.36`, `§6.3`) is still ignored.
const ROOT_ARTIFACT = /^[A-Za-z0-9_.-]{2,}\.(?:md|js|mjs|cjs|ts|go|mod|sum|json|ps1|sh|py|ya?ml|txt)$/;

function normalize(p) {
  return String(p).trim().replace(/\\/g, '/').replace(/^\.\//, '');
}

/** Rows of every ```artifacts fence, keeping only the paths whose disposition is repo. */
function parseArtifacts(text) {
  const lines = String(text).split('\n');
  const out = [];
  let inside = false;
  for (const line of lines) {
    if (!inside && FENCE.test(line.trim())) { inside = true; continue; }
    if (inside && /^```\s*$/.test(line.trim())) { inside = false; continue; }
    if (!inside) continue;
    const row = String(line).split('\t');
    const p = normalize(row[0]);
    if (!p) continue;
    const disposition = (row[1] || 'repo').trim();
    if (disposition !== 'repo') continue;
    out.push(p);
  }
  return out;
}

/** Repository-looking paths in a ledger completion-row evidence cell (backticked or bare). */
function evidencePathsFromRow(line) {
  const cells = String(line)
    .replace(/^\|/, '')
    .replace(/\|\s*$/, '')
    .split('|')
    .map((s) => s.trim());
  const cell = cells[cells.length - 1] || '';
  const out = [];
  for (const m of cell.matchAll(/`([^`]+)`/g)) {
    const raw = m[1].trim();
    if (raw.indexOf('{') >= 0) continue;
    if (PATH_TOKEN.test(raw)) out.push(normalize(raw.replace(/[).,;]+$/, '')));
  }
  for (const m of cell.matchAll(/([\s（(、,，；;）]|^)((?:docs|tools|examples|schemas|cmd|internal|harnesses|testdata|\.github)\/[^\s；;、,，|)`）]+)/g)) {
    const raw = m[2].replace(/[).,;]+$/, '');
    if (raw.indexOf('{') >= 0) continue;
    out.push(normalize(raw));
  }
  // The trailing lookahead is the RIGHT boundary: without it `go.modification` was extracted as
  // `go.mod` (⑦ r2), i.e. prose could be mistaken for an artifact reference.
  for (const m of cell.matchAll(/(^|[\s（(、,，；;）`])([A-Za-z0-9_.-]{2,}\.(?:md|js|mjs|cjs|ts|go|mod|sum|json|ps1|sh|py|ya?ml|txt))(?=[\s（(、,，；;）`|)\]}]|$)/g)) {
    if (!ROOT_ARTIFACT.test(m[2])) continue;
    out.push(normalize(m[2]));
  }
  return Array.from(new Set(out));
}

function identical(a, b) {
  const A = Array.from(new Set(a)).sort();
  const B = Array.from(new Set(b)).sort();
  return A.length === B.length && A.every((x, i) => x === B[i]);
}

function run(ctx, out) {
  const changed = ctx.changedPaths.map(normalize);
  const live = changed.filter((p) => {
    const rec = typeof ctx.recordFor === 'function' ? ctx.recordFor(p) : null;
    // ctx records carry a two-char XY status (`"D "` / `" D"`), so compare the parser's
    // `deleted` boolean - `status !== 'D'` was permanently true and never excluded a deletion.
    return !rec || !rec.deleted;
  });
  const reports = live.filter((p) => REPORT_RE.test(p) && p.indexOf('docs/validation/evidence/') !== 0);

  // ---- declaredFiles
  const declared = [];
  for (const rel of reports) {
    const text = ctx.readText(rel);
    if (text === null || text === undefined) throw new UsageError('G5: the changed report ' + rel + ' is unreadable');
    declared.push(rel); // the report always declares itself
    declared.push(...parseArtifacts(text));
  }
  const declaredSet = Array.from(new Set(declared));

  // ---- evidenceRefs (added completion rows only)
  const refs = [];
  if (typeof ctx.addedLines === 'function') {
    for (const rel of LEDGER_PAIR) {
      for (const line of ctx.addedLines(rel) || []) {
        if (!/^\|\s*\d{4}-\d{2}-\d{2}\s*\|/.test(line)) continue;
        refs.push(...evidencePathsFromRow(line));
      }
    }
  }

  // ---- index hygiene
  let dirty = [];
  if (ctx.kind === 'index') {
    const r = ctx.git(['status', '--porcelain', '-uall']);
    if (!r.ok) throw new UsageError('G5: environment error - git status --porcelain failed (exit ' + r.code + ')');
    dirty = String(r.out || '')
      .split('\n')
      .map((l) => l.replace(/\s+$/, ''))
      .filter(Boolean)
      // porcelain XY: X = index vs HEAD, Y = worktree vs index. A staged change is the
      // normal state of index scope; what must not exist is a WORKTREE edit (` Y`) or an
      // untracked path (`??`), because then the change set under test is not the tree.
      .filter((l) => (/^\?\?/.test(l) ? true : l.length > 1 && l[1] !== ' '));
  }

  // ---- optional external cross-check (never a licence to pass)
  const external = process.env.GATES_CHANGESET_JSON;
  if (external) {
    let data;
    try { data = JSON.parse(external); } catch (e) { throw new UsageError('GATES_CHANGESET_JSON is not valid JSON: ' + e.message); }
    const mism = [];
    const bad = [];
    const cmp = (key, actual) => {
      if (data[key] === undefined) return;
      if (!Array.isArray(data[key])) { bad.push(key + ' (not an array)'); return; }
      if (!identical(data[key].map(normalize), actual)) mism.push(key);
    };
    cmp('stagedFiles', live);
    cmp('declaredFiles', declaredSet);
    cmp('evidenceRefs', refs);
    if (data.addMode !== undefined && typeof data.addMode !== 'string') bad.push('addMode (not a string)');
    if (bad.length) throw new UsageError('GATES_CHANGESET_JSON has ill-typed fields: ' + bad.join(', '));
    if (mism.length) throw new UsageError('GATES_CHANGESET_JSON disagrees with the derived change set: ' + mism.join(', '));
  }

  // ---- G5(d): staged subset of declared
  let addMode = 'n/a';
  if (!reports.length) {
    out.rec('G5(d)', 'INFO', 'change set vs declared set', 'no report in the change set; nothing declares this change');
    addMode = 'unverifiable';
  } else {
    const missing = live.filter((p) => declaredSet.indexOf(p) < 0);
    const extra = declaredSet.filter((p) => live.indexOf(p) < 0);
    if (dirty.length) addMode = 'incomplete';
    else if (missing.length) addMode = 'batch-suspect';
    else addMode = 'explicit';
    out.rec(
      'G5(d)',
      missing.length === 0 ? 'PASS' : 'FAIL',
      'staged is a subset of declared',
      (missing.length === 0 ? live.length + ' path(s) declared (scope=' + ctx.kind + ')' : JSON.stringify({ stagedNotDeclared: missing.slice(0, 5) })) +
        (extra.length ? ' ; declared-not-touched=' + extra.length + ' (INFO: a fence is the landing snapshot)' : '') +
        ' ; addMode=' + addMode + ' (derived label, not a command-history assertion)'
    );
  }

  if (reports.length) {
    const extraCount = declaredSet.filter((p) => live.indexOf(p) < 0).length;
    if (extraCount > 0) {
      out.rec('G5(d-info)', 'INFO', 'declared but not touched', extraCount + ' declared path(s) are outside the change set (a fence is the landing snapshot); G5(d) judges the other direction only');
    }
  }

  // ---- G5(e): evidence references resolve
  const unresolved = [];
  for (const ref of refs) {
    const probe = ref.replace(/\/+$/, '');
    if (!ctx.repoArtifactResolvable(probe)) unresolved.push(ref);
  }
  out.rec(
    'G5(e)',
    unresolved.length === 0 ? 'PASS' : 'FAIL',
    'evidence references resolve',
    refs.length === 0 ? 'no added completion row (nothing to resolve)' : unresolved.length ? JSON.stringify(unresolved.slice(0, 5)) : refs.length + ' reference(s) resolved'
  );

  // ---- G5(f): index hygiene
  if (ctx.kind !== 'index') {
    out.rec('G5(f)', 'INFO', 'index hygiene', 'scope=' + ctx.kind + ': staging is not part of this scope');
  } else {
    out.rec('G5(f)', dirty.length === 0 ? 'PASS' : 'FAIL', 'index hygiene', dirty.length === 0 ? 'index equals worktree, no untracked path' : JSON.stringify(dirty.slice(0, 5)));
  }
}

module.exports = {
  id: 'G5',
  label: 'change-set contract (staged / declared / evidence refs / addMode)',
  title: '变更集契约',
  scripted: true,
  run,
  _internals: { parseArtifacts, evidencePathsFromRow, identical, normalize, REPORT_RE, PATH_TOKEN, ROOT_ARTIFACT },
};

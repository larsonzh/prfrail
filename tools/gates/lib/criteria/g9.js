'use strict';

/**
 * G9 temp-artifact guard (DELIVERY_DIRECTIVE 6.3 G9, v1.38).
 *
 * G9 covers the STATIC shape of the "temp artifact" discipline defined in section 5.1:
 * a temp artifact that was COMMITTED into the repository, plus (G9-e, OB-77 / v1.40) a
 * NON-BLOCKING readability check on `tmp/` references written into changed documents. It
 * has three predicates:
 *
 *   G9(1) temp-artifact shape
 *          a changed path whose BASENAME matches `\.(log|tmp|bak|out)$` (case-insensitive)
 *          fails. Paths below the tmp root are deliberately adjudicated by G9(2) only -
 *          the root guard is the stronger form (a `git add -f` forced entry) and owns
 *          every `tmp/` path, so one path never produces two findings for one defect.
 *   G9(2) tmp-root guard
 *          a changed path that starts with `tmp/` and is not the tracked directory
 *          placeholder `tmp/.gitkeep` fails. The tmp-root test is case-insensitive
 *          because `tmp/` and `Tmp/` resolve to the same directory on Windows and
 *          macOS; the finding still names the path with its original spelling. This is
 *          the only way a `tmp/` path can enter the change set (the directory is
 *          gitignored); force-added paths under OTHER gitignored roots are out of
 *          scope.
 *   G9(-e) tmp/ reference readability (INFORMATIONAL, never red)
 *          every `tmp/`-prefixed token with at least one path segment on an ADDED line
 *          of a changed `.md` document must resolve - `tree`/`index` = present on disk
 *          AND staged, `range`/`ci` = present at the scope endpoint. An unresolvable
 *          reference is recorded as one INFO record per token plus a summary INFO; it
 *          never changes the verdict. A bare `tmp/` (no segment) is not a reference;
 *          the G6 exclusion markers skip a line; the frozen evidence tree is excluded.
 *          The temporal shape (`$env:TEMP`, paths outside the repository) remains out of
 *          scope exactly as declared below.
 *
 * Exclusions (mandatory, single-source reading of section 6.3):
 *   - the frozen evidence tree `docs/validation/evidence/**` (OB-79) is judged by G8,
 *     not by this criterion - `ctx.isEvidenceTree` is the single source, with the
 *     module-level `isEvidenceTree` as a fallback for minimal probe contexts;
 *   - deleted paths (`ctx.recordFor(p).deleted === true`) submit no bytes, so they are
 *     not "temp artifacts committed into the repository";
 *   - whitelist hits (`tools/gates/g9-whitelist.txt`, the shared
 *     `<match><TAB><file><TAB><reason>` format). A malformed registration is a
 *     `UsageError` (exit 2), never a silent pass.
 *
 * An empty change set reports INFO (the established precedent: no path to judge is not
 * a vacuous PASS, and it is not an error either).
 *
 * Known boundary (declared in section 6.3): G9 does NOT cover the TEMPORAL shape -
 * `$env:TEMP`, locations outside the repository, PowerShell pipelines, or shell
 * variables - because those leave zero information in the CI input. The temporal shape
 * is registered as OB-99. G9 must not be read as "the temp-artifact discipline is fully
 * mechanised".
 */

const { UsageError, isEvidenceTree, toPosix } = require('../ctx');

const WHITELIST_REL = 'tools/gates/g9-whitelist.txt';
const TEMP_SHAPE_RE = /\.(log|tmp|bak|out)$/i;
const TMP_ROOT_PREFIX = 'tmp/';
const TMP_ROOT_KEEP = 'tmp/.gitkeep';

/**
 * G9-e reference shape (OB-77): `tmp/` plus at least one path segment, either inside a
 * code span or bare. The character set mirrors the repository's ordinary path alphabet;
 * anything else (CJK punctuation, quotes, `)`, end of line) terminates the token.
 */
const TMP_REF_RE = /tmp\/[A-Za-z0-9_.-]+(?:\/[A-Za-z0-9_.-]+)*/g;
/**
 * The five exclusion markers of section 6.3, kept in sync with g6.js. A line that quotes
 * a forbidden shape (rule text / example / quotation / conflict discussion) must not be
 * judged by the informational scan either.
 */
const TMP_REF_EXCLUDE = ['不得残留', '示例', '引文', '冲突', '矛盾'];

/**
 * Every distinct `tmp/` reference of one line, with a sentence period / trailing slash
 * stripped (a path written bare at the end of a sentence must not carry the period).
 * A bare `tmp/` produces no match at all, which is the declared boundary.
 */
function tmpReferencesOf(line) {
  const out = [];
  for (const m of String(line).matchAll(TMP_REF_RE)) {
    const tok = m[0].replace(/[./]+$/, '');
    if (!tok || tok === 'tmp') continue;
    if (out.indexOf(tok) < 0) out.push(tok);
  }
  return out;
}

function basename(p) {
  const s = toPosix(p);
  const i = s.lastIndexOf('/');
  return i < 0 ? s : s.slice(i + 1);
}

/** Same matching rule as the shared `whitelisted()` helper: file equality, substring match. */
function isWhitelisted(entries, p) {
  return entries.some((e) => (e.file === '*' || e.file === p) && p.indexOf(e.match) >= 0);
}

module.exports = {
  id: 'G9',
  label: 'temp artifacts and tmp/ references',
  title: '临时件与 tmp/ 引用',
  scripted: true,
  internals: { TMP_REF_RE, TMP_REF_EXCLUDE, tmpReferencesOf },

  run(ctx, out) {
    if (!ctx.changedPaths || ctx.changedPaths.length === 0) {
      out.rec('G9', 'INFO', '临时件与 tmp/ 引用', 'no changed path in scope');
      return { findings: [] };
    }

    const inTree = typeof ctx.isEvidenceTree === 'function' ? ctx.isEvidenceTree : isEvidenceTree;
    const whitelist = ctx.loadWhitelist(WHITELIST_REL);
    if (whitelist.error) throw new UsageError(whitelist.error);
    const entries = whitelist.entries || [];

    const bad = [];
    let excluded = 0;
    for (const raw of ctx.changedPaths) {
      const p = toPosix(raw);

      if (inTree(p)) {
        excluded++;
        continue;
      }
      const rec = typeof ctx.recordFor === 'function' ? ctx.recordFor(p) : null;
      if (rec && rec.deleted === true) {
        excluded++;
        continue;
      }
      if (isWhitelisted(entries, p)) {
        excluded++;
        continue;
      }

      // G9(2) owns the whole tmp root; G9(1) owns everything outside it. The tmp-root
      // test is case-insensitive (tmp/ and Tmp/ are the same directory on Windows and
      // macOS); the finding below still names the original path.
      const key = p.toLowerCase();
      if (key.startsWith(TMP_ROOT_PREFIX)) {
        if (key !== TMP_ROOT_KEEP) bad.push('tmp-root: ' + p);
        continue;
      }
      if (TEMP_SHAPE_RE.test(basename(p))) bad.push('temp-shape: ' + p);
    }

    const shown = bad.slice(0, 5);
    out.rec(
      'G9',
      bad.length === 0 ? 'PASS' : 'FAIL',
      '临时件与 tmp/ 引用',
      bad.length
        ? 'findings=' + bad.length + ' ' + JSON.stringify(shown) + (bad.length > shown.length ? ' (+' + (bad.length - shown.length) + ' more)' : '')
        : ctx.changedPaths.length +
            ' path(s) OK' +
            (excluded ? ' (' + excluded + ' excluded)' : '') +
            ' (whitelist=' + whitelist.source + ')'
    );
    runTmpReferenceCheck(ctx, out, inTree);
    return { findings: bad };
  },
};

/**
 * G9-e (OB-77): informational only. Emits one INFO record per unresolvable `tmp/`
 * reference plus one summary INFO. The G9 verdict is decided by predicates (1)/(2) and
 * is never touched here - a dangling reference is a documentation defect to surface, not
 * a committed temp artifact to block on.
 */
function runTmpReferenceCheck(ctx, out, inTree) {
  if (typeof ctx.changedMarkdown !== 'function' || typeof ctx.readLineArray !== 'function') return;
  if (typeof ctx.repoArtifactResolvable !== 'function') return;
  let linesScanned = 0;
  let refs = 0;
  const unresolved = [];
  for (const f of ctx.changedMarkdown()) {
    if (inTree(f)) continue;
    const L = ctx.readLineArray(f);
    if (!L) continue;
    const added = ctx.addedLineNumbers ? [...ctx.addedLineNumbers(f)] : [];
    for (const n of added.sort((a, b) => a - b)) {
      const line = L[n - 1];
      if (line === undefined) continue;
      linesScanned++;
      if (TMP_REF_EXCLUDE.some((m) => line.includes(m))) continue;
      for (const tok of tmpReferencesOf(line)) {
        refs++;
        const r = ctx.repoArtifactResolvable(tok);
        if (r && r.ok) continue;
        unresolved.push({ file: f, line: n, path: tok, reason: (r && r.reason) || 'not resolvable' });
      }
    }
  }
  for (const u of unresolved) {
    out.rec('G9-e', 'INFO', 'tmp/ reference (informational)', u.file + ':' + u.line + ' ' + u.path + ' :: ' + u.reason);
  }
  out.rec(
    'G9-e',
    'INFO',
    'tmp/ reference summary',
    'scanned ' + linesScanned + ' added line(s), ' + refs + ' reference(s), ' + unresolved.length + ' unresolvable'
  );
}

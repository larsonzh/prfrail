'use strict';

/**
 * G9 temp-artifact guard (DELIVERY_DIRECTIVE 6.3 G9, v1.38).
 *
 * G9 covers the STATIC shape of the "temp artifact" discipline defined in section 5.1:
 * a temp artifact that was COMMITTED into the repository. It has exactly two predicates
 * and no sub-checks (the runner reports one `G9` record):
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
  label: 'temp artifacts',
  title: '临时件',
  scripted: true,

  run(ctx, out) {
    if (!ctx.changedPaths || ctx.changedPaths.length === 0) {
      out.rec('G9', 'INFO', '临时件', 'no changed path in scope');
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
      '临时件',
      bad.length
        ? 'findings=' + bad.length + ' ' + JSON.stringify(shown) + (bad.length > shown.length ? ' (+' + (bad.length - shown.length) + ' more)' : '')
        : ctx.changedPaths.length +
            ' path(s) OK' +
            (excluded ? ' (' + excluded + ' excluded)' : '') +
            ' (whitelist=' + whitelist.source + ')'
    );
    return { findings: bad };
  },
};

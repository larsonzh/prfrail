'use strict';

/**
 * G5-c appendix-B call strings against the section 2.2 qualified-name column (OB-29, v1.33).
 *
 * The gap OB-29 registered: nothing mechanical checked that the `model='...'` call strings
 * in appendix B belong to the section 2.2 qualified-name column - the values an operator
 * actually pastes into `runSubagent`. The allowed values are read from the DOCUMENT ITSELF
 * (2.2's qualified-name column), so the criterion and the authority cannot drift apart:
 * editing 2.2 moves the allowed set with it. (The OB-29 entry was originally worded against
 * the `modelId` column; the candidate wording was synced to the qualified-name column when
 * R2.2 made the call string the sole authority - see OB-67 / OB-70.)
 *
 * Comparison is whitespace- and backtick-insensitive: 2.2 writes the names as code spans
 * (`DeepSeek V4 Pro (deepseek)`) while the appendix B call strings are plain text inside a
 * fence, and both mirror sides must compare equal without depending on that markup.
 *
 * Scope: ADDED lines inside the appendix B region only (the same change-set discipline as
 * G6 / G7-c / G7-d). Unchanged historical call strings are therefore not re-judged; the
 * seven appendix B strings were verified against 2.2 when they landed (v1.12).
 *
 * Fail closed: a directive file that HAS an appendix B region whose 2.2 qualified-name
 * column cannot be read is a premise error (UsageError, exit 2) - never a silent pass.
 */
const { UsageError } = require('../ctx');

const DIRECTIVE_FILES = ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md'];
const HEAD_22 = /^#{2,4}\s+2\.2(?=[\s.．、]|$)/;
const HEAD_APX_B = /^#{2,4}\s+(?:附录\s*B|Appendix\s*B)(?=[\s：:．、]|$)/i;
/** Region terminator: the next level-1/2 heading (`### B.1` stays inside appendix B). */
const HEAD_H2 = /^#{1,2}\s/;
/** Header cell of the qualified-name column, in either language (the EN side is lowercase). */
const QUALIFIED_HEADER = /限定名|qualified name/i;
/** The call-string shape: `model='<qualified name>'` (the appendix B carrier lines). */
const CALL_RE = /model\s*=\s*'([^']*)'/g;

/** Whitespace- and backtick-insensitive form used for every comparison. */
function norm(v) {
  return String(v).replace(/[`\s]/g, '');
}

/** `[start, end)` of the block that opens with `headRe`, or null when the file has none. */
function regionAfter(L, headRe) {
  for (let i = 0; i < L.length; i++) {
    if (!headRe.test(L[i])) continue;
    let end = L.length;
    for (let j = i + 1; j < L.length; j++) {
      if (HEAD_H2.test(L[j])) {
        end = j;
        break;
      }
    }
    return [i, end];
  }
  return null;
}

/**
 * The section 2.2 qualified-name column as a set of normalised values. The column index is
 * positional (4th cell), which is what keeps the check identical on both mirror sides.
 *
 * KNOWN BOUNDARY (stated by the 5 pre-review, L-6): the collected set is the 4th column of
 * the CONTIGUOUS table block that follows the header row. The 2.2 region itself ends at the
 * next level-1/2 heading, so a table block glued onto the 2.2 table without a blank line
 * would dilute the set (fail-open). Measured today: exactly 8 distinct values, no dilution;
 * carried by the 5 / 6 / 7 review chain rather than by the criterion.
 */
function qualifiedNames(L, f) {
  const region = regionAfter(L, HEAD_22);
  if (!region) {
    throw new UsageError('G5-c: no 2.2 section in ' + f + ' (the qualified-name authority is missing)');
  }
  let hdr = -1;
  for (let i = region[0]; i < Math.min(region[0] + 12, L.length); i++) {
    if (L[i].startsWith('|') && QUALIFIED_HEADER.test(L[i])) {
      hdr = i;
      break;
    }
  }
  if (hdr < 0) {
    throw new UsageError('G5-c: the 2.2 table header with the qualified-name column was not found in ' + f);
  }
  const out = new Set();
  for (let i = hdr + 2; i < L.length && L[i].startsWith('|'); i++) {
    const v = norm(L[i].split('|')[4] || '');
    if (v) out.add(v);
  }
  if (!out.size) {
    throw new UsageError('G5-c: the 2.2 qualified-name column is empty in ' + f);
  }
  return out;
}

module.exports = {
  id: 'G5-c',
  label: 'appendix-B call strings in the 2.2 qualified-name set',
  title: '调用串限定名闭集',
  scripted: true,
  internals: { regionAfter, qualifiedNames, norm, DIRECTIVE_FILES },

  run(ctx, out) {
    const findings = [];
    let files = 0;
    let scanned = 0;
    for (const f of DIRECTIVE_FILES) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const apx = regionAfter(L, HEAD_APX_B);
      if (!apx) continue; // no appendix B region: nothing to judge
      const allowed = qualifiedNames(L, f);
      files++;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const idx = n - 1;
        if (idx < apx[0] || idx >= apx[1]) continue;
        const line = L[idx];
        if (line === undefined) continue;
        CALL_RE.lastIndex = 0;
        let m;
        let judged = false;
        while ((m = CALL_RE.exec(line)) !== null) {
          judged = true;
          if (allowed.has(norm(m[1]))) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail: f + ':' + n + ' call string `' + m[1] + '` is not in the 2.2 qualified-name column',
          });
        }
        if (judged) scanned++;
      }
    }
    out.rec(
      'G5-c',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'appendix-B call strings in the 2.2 qualified-name set',
      findings.length
        ? findings.map((x) => x.detail).join(' | ')
        : files === 0
          ? 'no directive with an appendix B region in scope'
          : 'call strings OK (' + files + ' file(s), scanned ' + scanned + ' changed carrier line(s))'
    );
    return { findings };
  },
};

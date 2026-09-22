'use strict';

/**
 * G2 placeholder residue (DELIVERY_DIRECTIVE 6.3 G2).
 * No placeholder wording may survive in added lines. Self-reference guard: a line
 * that quotes the forbidden strings must also carry `不得残留` so the rule text
 * does not flag itself (measured 2026-09-21: 4 false positives without it).
 * Pre-existing occurrences in a changed file are reported as INFO, not as FAIL.
 */
const PATTERNS = [/待[^\n]{0,6}回填/, /尚未提交/, /待同轮授权/, /pending the push/];

module.exports = {
  id: 'G2',
  label: 'placeholder residue',
  title: '占位符残留',
  scripted: true,

  run(ctx, out) {
    const hits = [];
    const preexisting = [];
    for (const f of ctx.changedMarkdown()) {
      const lines = ctx.readLineArray(f);
      if (!lines) continue;
      const added = ctx.addedLines(f);
      lines.forEach((l, i) => {
        if (!PATTERNS.some((p) => p.test(l)) || /不得残留/.test(l)) return;
        (added.includes(l) ? hits : preexisting).push(f + ':' + (i + 1));
      });
    }
    out.rec(
      'G2',
      hits.length === 0 ? 'PASS' : 'FAIL',
      'placeholder residue',
      (hits.length ? JSON.stringify(hits) : 'none') +
        (preexisting.length ? '  [INFO pre-existing: ' + preexisting.join(', ') + ']' : '')
    );
  },
};

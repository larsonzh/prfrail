'use strict';

/**
 * G2 placeholder residue (DELIVERY_DIRECTIVE 6.3 G2).
 * No placeholder wording may survive in added lines. Self-reference guard: a line
 * that quotes the forbidden strings must also carry `不得残留` so the rule text
 * does not flag itself (measured 2026-09-21: 4 false positives without it).
 * Pre-existing occurrences in a changed file are reported as INFO, not as FAIL.
 *
 * Scope exemption (v1.77, slice P94 GATES-EXEMPT): single source = g7 internals; see
 * g6.js P93. The body of a CLOSED, structurally valid `sixstage` fence in a G7 target
 * report is the verbatim ⑤/⑥/⑦ transcript (§11.3 + G7-f), which the controller must NOT
 * rewrite - so a placeholder string quoted inside it is not the controller's own prose
 * and cannot be judged (same scope argument as the evidence tree, OB-79). It is NOT a
 * whitelist duty: `g6-whitelist.txt` stays empty. Exempt lines count neither as hits
 * nor as pre-existing INFO; the exemption is only reported when it happens.
 */
const G7 = require('./g7');
const g7Targets = G7.internals.targets;
const g7SixstageBlocks = G7.internals.sixstageBlocks;
const g7SectionRe = G7.internals.SIXSTAGE_SECTION_RE;
const g7SectionSeq = G7.internals.SIXSTAGE_SECTION_SEQ;

const PATTERNS = [/待[^\n]{0,6}回填/, /尚未提交/, /待同轮授权/, /pending the push/];

/**
 * Does the fence body carry the five `Section A`-`Section E` headings IN ORDER? This is
 * isomorphic to G6's `sixstageStructurallyValid` (P93) and to G7-f's structural face:
 * `.trim()` mirrors G7-f exactly, so an indented heading passes everywhere or nowhere.
 */
function sixstageStructurallyValid(body) {
  let seen = 0;
  for (const line of body) {
    const m = g7SectionRe.exec(String(line).trim());
    if (!m) continue;
    if (m[1] !== g7SectionSeq.charAt(seen)) return false;
    seen++;
  }
  return seen === g7SectionSeq.length;
}

/**
 * 1-based lines exempt in `file`, or null when none. Fail-closed: an UNCLOSED fence
 * exempts nothing (G7-f already FAILs that shape), a CLOSED fence without the five
 * ordered sections exempts nothing, and a same-shaped fence outside the G7 target set
 * exempts nothing. Interval = [start, end): the closing fence line itself is plain
 * ``` and can never match a placeholder pattern.
 */
function sixstageExemptSet(file, L) {
  if (!g7Targets({ changedPaths: [file] }).length) return null;
  let set = null;
  for (const b of g7SixstageBlocks(L)) {
    if (!b.closed) continue;
    if (!sixstageStructurallyValid(b.body)) continue;
    if (!set) set = new Set();
    for (let i = b.start; i < b.end; i++) set.add(i + 1);
  }
  return set;
}

module.exports = {
  id: 'G2',
  label: 'placeholder residue',
  title: '占位符残留',
  scripted: true,

  run(ctx, out) {
    const hits = [];
    const preexisting = [];
    let exempted = 0;
    for (const f of ctx.changedMarkdown()) {
      const lines = ctx.readLineArray(f);
      if (!lines) continue;
      const exempt = sixstageExemptSet(f, lines);
      const added = ctx.addedLines(f);
      // Line-number based ONLY (P94 (7) re-review): without the line-number view we
      // cannot tell an added exempt line from a pre-existing one, so the conservative
      // choice is to not count it (the exempt region is skipped either way).
      const addedNums = typeof ctx.addedLineNumbers === 'function' ? ctx.addedLineNumbers(f) : null;
      const countsExempt = (n) =>
        Boolean(addedNums && typeof addedNums.has === 'function' && addedNums.has(n));
      lines.forEach((l, i) => {
        if (exempt && exempt.has(i + 1)) {
          // Same caliber as G6-4 / G4b (P94 review, Low): `fence-exempt N` counts the
          // exempt ADDED lines only. The whole exempt region is still skipped here -
          // a pre-existing exempt line must not surface as a pre-existing INFO either.
          if (countsExempt(i + 1)) exempted++;
          return;
        }
        if (!PATTERNS.some((p) => p.test(l)) || /不得残留/.test(l)) return;
        (added.includes(l) ? hits : preexisting).push(f + ':' + (i + 1));
      });
    }
    out.rec(
      'G2',
      hits.length === 0 ? 'PASS' : 'FAIL',
      'placeholder residue',
      (hits.length ? JSON.stringify(hits) : 'none') +
        (preexisting.length ? '  [INFO pre-existing: ' + preexisting.join(', ') + ']' : '') +
        (exempted > 0 ? '  fence-exempt ' + exempted : '')
    );
  },
};

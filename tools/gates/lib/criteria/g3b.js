'use strict';

/**
 * G3(b) key-field alignment (DELIVERY_DIRECTIVE 6.3 G3).
 * Scope = added lines only. Line-count symmetry cannot prove that the semantics of
 * the CN/`_EN` pair did not drift, so the mirrored key fields (commit hashes, CI run
 * numbers, status glyphs, to-do markers, the not-reviewed row) are counted on both
 * sides and must match.
 *
 * No changed `.md` in scope: an explicit INFO line instead of silence, so an empty
 * run is visibly empty (mirrors the G7 / scope-guard convention).
 */
module.exports = {
  id: 'G3-b',
  label: 'key-field alignment',
  title: '关键字段对齐（变更行）',
  scripted: true,

  run(ctx, out) {
    if (ctx.changedMarkdown().length === 0) {
      out.rec('G3(b)', 'INFO', 'key-field alignment', 'no changed .md in scope');
      return;
    }
    const changed = ctx.changedPaths;
    const pairs = changed.filter(
      (f) => /\.md$/i.test(f) && !/_EN\.md$/.test(f) && changed.includes(f.replace(/\.md$/, '_EN.md'))
    );
    const PB = [
      [/\b[0-9a-f]{7}\b/g, /\b[0-9a-f]{7}\b/g, 'hash'],
      [/\brun \d+\b/g, /\brun \d+\b/g, 'run'],
      [/[✅⬜]/g, /[✅⬜]/g, 'glyph'],
      [/\*\*待办\*\*/g, /\*\*To do\*\*/g, 'to-do'],
      [/^\| \*\*未过独立终审\*\*/gm, /^\| \*\*Did not pass independent final review\*\*/gim, 'not-reviewed'],
    ];
    const info = [];
    for (const cnf of pairs) {
      const A = ctx.addedLines(cnf).join('\n');
      const B = ctx.addedLines(cnf.replace(/\.md$/, '_EN.md')).join('\n');
      for (const [rc, re, label] of PB) {
        const c = (A.match(rc) || []).length;
        const e = (B.match(re) || []).length;
        info.push(label + ' ' + c + '/' + e + (c === e ? ' OK' : ' MISMATCH'));
      }
    }
    if (pairs.length) {
      out.rec(
        'G3(b)',
        info.every((s) => / OK$/.test(s)) ? 'PASS' : 'FAIL',
        'key-field alignment',
        info.join(' | ')
      );
    }
  },
};

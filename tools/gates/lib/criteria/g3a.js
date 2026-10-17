'use strict';

/**
 * G3(a) bilingual symmetry + mirror parity (DELIVERY_DIRECTIVE 6.3 G3).
 *   (a.1) for every changed CN/`_EN` pair the numstat added/deleted counts must be equal;
 *   (a.2) for `docs/DELIVERY_DIRECTIVE.md` the whole-file shape (line count, heading
 *         positions, bold markers, pipes, blank lines, backtick count) must equal its
 *         `_EN` mirror. The sixth shape (backticks) landed with OB-73 (v1.33), which also
 *         aligned the 7 pre-existing backtick-level drifts it exposes.
 * Whole-file equality is deliberately NOT required for other pairs: pre-existing
 * bilingual pairs are historically asymmetric, so the criterion would be unsatisfiable.
 *
 * (a.1) and (a.2) never emit a PASS record of their own (only the aggregate G3(a)
 * does), so with no changed `.md` in scope they must say so explicitly instead of
 * printing nothing at all.
 *
 * OB-90 (v1.34): the pair is now formed from EITHER side. Before this the filter started
 * from the CN side and demanded both sides in the change set, so an `_EN.md`-only edit
 * produced NO record at all - a bilingual criterion that failed silently. A baseline
 * `.md` change is still paired only when an `_EN` mirror exists (scope endpoint / disk /
 * same change set), because most corpus `.md` files are single-language and editing one
 * must not turn this criterion red.
 */
function count(lines, re) {
  return lines.reduce((n, l) => n + (l.match(re) || []).length, 0);
}

module.exports = {
  id: 'G3-a',
  label: 'symmetry + mirror parity',
  title: '双语对称',
  scripted: true,

  run(ctx, out) {
    if (ctx.changedMarkdown().length === 0) {
      out.rec('G3(a.1)', 'INFO', 'symmetry', 'no changed .md in scope');
      out.rec('G3(a.2)', 'INFO', 'mirror parity', 'no changed .md in scope');
      return;
    }
    const changed = ctx.changedPaths;
    const paired = new Set();
    const sym = [];
    const inf = [];
    for (const f of changed) {
      if (!/\.md$/i.test(f)) continue;
      // The pair key is always the canonical CN path, so the record id is the same
      // whichever side was edited (the edited side shows up in the two numstat halves).
      const enSide = /_EN\.md$/i.test(f);
      const cnf = enSide ? f.replace(/_EN\.md$/i, '.md') : f;
      if (paired.has(cnf)) continue;
      // F4 (5 r1): on the EN side the judged path is the CHANGED one, so a non-canonical
      // spelling (`X_en.md`) is judged with its own numstat instead of comparing two
      // untouched canonical paths and printing an untested green aggregate.
      const enf = enSide ? f : cnf.replace(/\.md$/i, '_EN.md');
      // OB-90: an `_EN` change always pairs (its CN baseline simply scores +0/-0 when
      // untouched, which is a visible G3(a.1) FAIL instead of silence). A baseline `.md`
      // change pairs only when a mirror exists, so single-language documents stay silent.
      if (
        !enSide &&
        !changed.includes(enf) &&
        !(typeof ctx.exists === 'function' && ctx.exists(enf))
      ) {
        continue;
      }
      paired.add(cnf);
      const a = ctx.numstat(cnf);
      const b = ctx.numstat(enf);
      const okS = a[0] === b[0] && a[1] === b[1];
      sym.push(cnf + ' +' + a[0] + '/-' + a[1] + ' vs _EN +' + b[0] + '/-' + b[1] + (okS ? ' OK' : ' MISMATCH'));
      if (!okS) out.rec('G3(a.1) ' + cnf, 'FAIL', 'symmetry', sym[sym.length - 1]);
      const CL = ctx.readLineArray(cnf);
      const EL = ctx.readLineArray(enf);
      if (CL && EL) {
        let dr = 0;
        for (let i = 0; i < Math.max(CL.length, EL.length); i++) {
          if (/^#{2,4} /.test(CL[i] || '') !== /^#{2,4} /.test(EL[i] || '')) dr++;
        }
        const f5 = {
          lines: CL.length === EL.length,
          head: dr === 0,
          bold: count(CL, /\*\*/g) === count(EL, /\*\*/g),
          bars: count(CL, /\|/g) === count(EL, /\|/g),
          blank: CL.filter((l) => l === '').length === EL.filter((l) => l === '').length,
          // OB-73 (v1.33): the sixth shape. A whole-file count comparison, the same unit as
          // bold / pipes / blank lines, so a cross-line compensating drift stays uncovered -
          // stated as a known limit in the DELIVERY_DIRECTIVE 6.3 G3(a.2) sub-rule.
          backticks: count(CL, /`/g) === count(EL, /`/g),
        };
        inf.push(cnf + ' ' + JSON.stringify(f5));
        if (/^docs\/DELIVERY_DIRECTIVE\.md$/.test(cnf) && !Object.values(f5).every(Boolean)) {
          out.rec('G3(a.2) ' + cnf, 'FAIL', 'mirror parity', inf[inf.length - 1]);
        }
      } else if (/^docs\/DELIVERY_DIRECTIVE\.md$/.test(cnf)) {
        // One side of the directive pair is gone (deleted / never present): the six whole-file
        // shapes cannot be compared, so the premise is broken rather than satisfied. This is a
        // data state, not a broken criterion premise, so it reports FAIL instead of UsageError.
        out.rec(
          'G3(a.2) ' + cnf,
          'FAIL',
          'mirror parity',
          'side unreadable: ' +
            (CL ? enf : cnf) +
            ' - the whole-file shape cannot be compared, so the pair is judged red instead'
        );
      }
    }
    if (sym.length) {
      out.rec(
        'G3(a)',
        sym.every((s) => / OK$/.test(s)) ? 'PASS' : 'FAIL',
        'symmetry + mirror parity',
        sym.join(' || ') + ' ;; ' + inf.join(' ;; ')
      );
    }
  },
};

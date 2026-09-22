'use strict';

/**
 * G3(a) bilingual symmetry + mirror parity (DELIVERY_DIRECTIVE 6.3 G3).
 *   (a.1) for every changed CN/`_EN` pair the numstat added/deleted counts must be equal;
 *   (a.2) for `docs/DELIVERY_DIRECTIVE.md` the whole-file shape (line count, heading
 *         positions, bold markers, pipes, blank lines) must equal its `_EN` mirror.
 * Whole-file equality is deliberately NOT required for other pairs: pre-existing
 * bilingual pairs are historically asymmetric, so the criterion would be unsatisfiable.
 *
 * (a.1) and (a.2) never emit a PASS record of their own (only the aggregate G3(a)
 * does), so with no changed `.md` in scope they must say so explicitly instead of
 * printing nothing at all.
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
    const pairs = changed.filter(
      (f) => /\.md$/i.test(f) && !/_EN\.md$/.test(f) && changed.includes(f.replace(/\.md$/, '_EN.md'))
    );
    const sym = [];
    const inf = [];
    for (const cnf of pairs) {
      const enf = cnf.replace(/\.md$/, '_EN.md');
      const CL = ctx.readLineArray(cnf);
      const EL = ctx.readLineArray(enf);
      if (!CL || !EL) continue; // deleted side of the pair: nothing to mirror
      const a = ctx.numstat(cnf);
      const b = ctx.numstat(enf);
      const okS = a[0] === b[0] && a[1] === b[1];
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
      };
      sym.push(cnf + ' +' + a[0] + '/-' + a[1] + ' vs _EN +' + b[0] + '/-' + b[1] + (okS ? ' OK' : ' MISMATCH'));
      inf.push(cnf + ' ' + JSON.stringify(f5));
      if (!okS) out.rec('G3(a.1) ' + cnf, 'FAIL', 'symmetry', sym[sym.length - 1]);
      if (/^docs\/DELIVERY_DIRECTIVE\.md$/.test(cnf) && !Object.values(f5).every(Boolean)) {
        out.rec('G3(a.2) ' + cnf, 'FAIL', 'mirror parity', inf[inf.length - 1]);
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

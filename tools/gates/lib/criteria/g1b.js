'use strict';

/**
 * G1-b version metadata consistency (DELIVERY_DIRECTIVE 6.3 sub-rule).
 * The version declared in the document header must equal the latest `| vX.Y |`
 * row inside the `### 0.3` changelog block, and the header date must match that
 * row's date. The block bound is deliberate: scanning the whole document would
 * also hit the appendix C.1 `| vN.N |` comparison table.
 *
 * No changed `.md` in scope: an explicit INFO line instead of silence, so an empty
 * run is visibly empty (the same convention G7 and the scope guard already use).
 */
module.exports = {
  id: 'G1-b',
  label: 'version metadata',
  title: '版本元数据一致性',
  scripted: true,

  run(ctx, out) {
    if (ctx.changedMarkdown().length === 0) {
      out.rec('G1-b', 'INFO', 'version metadata', 'no changed .md in scope');
      return;
    }
    const info = [];
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const txt = L.join('\n');
      const hv = (/^版本：v(\d+\.\d+)/m.exec(txt) || /^Version: v(\d+\.\d+)/m.exec(txt) || [])[1];
      if (!hv) continue;
      const hd = (
        /^版本：v[\d.]+[^\n]{0,4}日期：(\d{4}-\d{2}-\d{2})/m.exec(txt) ||
        /^Version: v[\d.]+\. Date: (\d{4}-\d{2}-\d{2})/m.exec(txt) ||
        []
      )[1];
      const s = L.findIndex((l) => /^### 0\.3/.test(l));
      const e = L.findIndex((l, i) => i > s && /^#{2,3} /.test(l));
      const rows = L.slice(s, e < 0 ? L.length : e).filter((l) => /^\| v\d+\.\d+ \|/.test(l));
      const last = rows.length ? /^\| v(\d+\.\d+) \|/.exec(rows[rows.length - 1])[1] : '?';
      const ld = rows.length
        ? (/^\| v[\d.]+ \| (\d{4}-\d{2}-\d{2}) \|/.exec(rows[rows.length - 1]) || [])[1]
        : '?';
      const ok = hv === last && (!hd || hd === ld);
      info.push(
        f + ' header=v' + hv + ' last=v' + last + ' ' + hd + '/' + ld + ' rows=' + rows.length + (ok ? ' OK' : ' MISMATCH')
      );
      if (!ok) out.rec('G1-b ' + f, 'FAIL', 'version metadata', info[info.length - 1]);
    }
    if (info.length) {
      out.rec(
        'G1-b',
        info.every((s) => / rows=\d+ OK$/.test(s)) ? 'PASS' : 'FAIL',
        'version metadata',
        info.join(' | ')
      );
    }
  },
};

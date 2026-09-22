'use strict';

/**
 * G4a BOM + LF encoding check (DELIVERY_DIRECTIVE 6.3 G4).
 * Scope = the change set. `.md` / `.ps1` must carry a UTF-8 BOM and LF only; every
 * other file must be BOM-free and LF only.
 * Ignores are decided by `git check-ignore -q` instead of the throwaway script's
 * hardcoded `.github/agents/sol-orchestrator/` skip, and already-deleted files are
 * skipped instead of throwing ENOENT.
 */
const BOM = Buffer.from([0xef, 0xbb, 0xbf]);
const CRLF = Buffer.from('\r\n');

module.exports = {
  id: 'G4a',
  label: 'BOM+LF',
  title: '编码（BOM + LF）',
  scripted: true,

  run(ctx, out) {
    const bad = [];
    for (const f of ctx.changedPaths) {
      if (ctx.isIgnored(f)) continue;
      const b = ctx.readBytes(f);
      if (!b) continue; // deleted (or unreadable) file: nothing to submit
      const bom = b.length >= 3 && b[0] === BOM[0] && b[1] === BOM[1] && b[2] === BOM[2];
      const crlf = b.indexOf(CRLF) >= 0;
      const want = /\.(md|ps1)$/i.test(f);
      const ok = want ? bom && !crlf : !bom && !crlf;
      if (!ok) bad.push(f + ' bom=' + bom + ' crlf=' + crlf);
    }
    out.rec(
      'G4a',
      bad.length === 0 ? 'PASS' : 'FAIL',
      'BOM+LF',
      bad.length ? JSON.stringify(bad) : ctx.changedPaths.length + ' file(s) OK'
    );
  },
};

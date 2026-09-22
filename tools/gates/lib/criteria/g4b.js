'use strict';

/**
 * G4b anomalous character scan (DELIVERY_DIRECTIVE 6.3 G4).
 *     (CJK characters of the added lines) MINUS (CJK characters of the base-revision
 *     *.md corpus) MINUS (tools/gates/cjk-newwords.txt) MUST be empty.
 * The whitelist freezes the human judgement (each entry carries a reason + context
 * excerpt) so the assertion stays mechanically re-runnable. The base revision is the
 * pre-change corpus: HEAD for tree/index scopes, the range base for range scopes.
 *
 * The whitelist itself is read through the shared scope-aware loader
 * (`ctx.loadWhitelist`, GATES-EXT F1): tree/index use the working tree, range/ci use
 * the scope's ENDPOINT revision - exactly like G6's whitelist, so both criteria
 * always judge against the same snapshot of the registry.
 */
const { UsageError } = require('../ctx');

const WHITELIST_REL = 'tools/gates/cjk-newwords.txt';

function cjkChars(s) {
  const o = new Set();
  for (const ch of s) {
    const c = ch.codePointAt(0);
    if ((c >= 0x4e00 && c <= 0x9fff) || (c >= 0x3400 && c <= 0x4dbf) || (c >= 0xf900 && c <= 0xfaff)) {
      o.add(ch);
    }
  }
  return o;
}

module.exports = {
  id: 'G4b',
  label: 'CJK diff minus whitelist',
  title: '编码 + 异常字',
  scripted: true,

  run(ctx, out) {
    const ref =
      ctx.kind === 'range' || ctx.kind === 'ci' ? ctx.baseRef : 'HEAD';
    let corpus = new Set();
    for (const f of ctx.listAllMarkdown()) {
      const text = ctx.showAt(ref, f);
      if (text === null) continue;
      corpus = new Set([...corpus, ...cjkChars(text)]);
    }

    const whitelist = ctx.loadWhitelist(WHITELIST_REL);
    if (whitelist.error) throw new UsageError(whitelist.error);
    const white = new Set(whitelist.entries.map((e) => e.match));

    let newC = new Set();
    for (const f of ctx.changedMarkdown()) {
      for (const l of ctx.addedLines(f)) {
        for (const c of cjkChars(l)) newC.add(c);
      }
    }
    const unknown = [...newC].filter((c) => !corpus.has(c) && !white.has(c));
    out.rec(
      'G4b',
      unknown.length === 0 ? 'PASS' : 'FAIL',
      'CJK diff minus whitelist',
      'newCjk=' +
        newC.size +
        ' unknown=' +
        JSON.stringify(unknown) +
        ' (whitelist=' +
        white.size +
        ' whitelist-source=' +
        whitelist.source +
        ')'
    );
  },
};

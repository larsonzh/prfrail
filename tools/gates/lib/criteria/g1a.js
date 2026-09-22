'use strict';

/**
 * G1-a table structure integrity (DELIVERY_DIRECTIVE 6.3 sub-rule).
 * Added lines must not splice two table records onto one line: any added line
 * starting with `|` must not contain `||` once escaped `\|` sequences are removed.
 * Scope = changed lines, never the whole document.
 */
module.exports = {
  id: 'G1-a',
  label: 'table structure',
  title: '表格结构完整性',
  scripted: true,

  run(ctx, out) {
    const bad = [];
    for (const f of ctx.changedMarkdown()) {
      for (const l of ctx.addedLines(f)) {
        if (l.startsWith('|') && l.split('\\|').join('').includes('||')) {
          bad.push(f + ' :: ' + l.slice(0, 50));
        }
      }
    }
    out.rec(
      'G1-a',
      bad.length === 0 ? 'PASS' : 'FAIL',
      'table structure',
      bad.length ? JSON.stringify(bad) : 'no `||` on added rows'
    );
  },
};

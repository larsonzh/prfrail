'use strict';

/**
 * G5-b role-file `model` value gate (DELIVERY_DIRECTIVE 6.3 sub-rule G5-b).
 * Added/modified `.agent.md` files must not set `model` to a section 2.4 blacklist
 * value (Luna / Terra / Gemini / GPT-5.4 / Sol / Kimi, case insensitive).
 */
const BLACKLIST = /luna|terra|gemini|gpt-5\.4|sol|kimi/i;

module.exports = {
  id: 'G5-b',
  label: 'blacklist model value',
  title: '角色文件 model 取值闸',
  scripted: true,

  run(ctx, out) {
    const targets = ctx.changedPaths.filter((f) => /\.agent\.md$/i.test(f));
    const hits = [];
    for (const f of targets) {
      const text = ctx.readText(f);
      if (text === null) continue;
      const v = (/^model\s*:\s*(.+)$/m.exec(text) || [])[1] || '';
      if (v && BLACKLIST.test(v)) hits.push(f);
    }
    out.rec(
      'G5-b',
      hits.length === 0 ? 'PASS' : 'FAIL',
      'blacklist model value',
      'checked ' + targets.length + ' file(s)'
    );
  },
};

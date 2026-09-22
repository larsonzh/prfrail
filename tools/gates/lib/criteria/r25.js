'use strict';

/**
 * R2.5 role-file header note, position anchored (DELIVERY_DIRECTIVE 2.2 / 4.2).
 * After the frontmatter the first non-empty line must be the non-generated note,
 * and every `.agent.md` must carry exactly the same note (a single distinct value).
 */
module.exports = {
  id: 'R2.5',
  label: 'header-note position',
  title: '头注位置锚定',
  scripted: true,

  run(ctx, out) {
    const agents = ctx.listAgentFiles();
    const bad = [];
    const notes = new Set();
    for (const f of agents) {
      const t = ctx.readText('.github/agents/' + f);
      if (t === null) continue;
      const fm = (/^---\n[\s\S]*?\n---\n/.exec(t) || [])[0] || '';
      const first = (t.slice(fm.length).split('\n').find((l) => l.trim() !== '') || '').trim();
      notes.add(first);
      if (!/非生成物/.test(first) || !/不得由生成器覆盖/.test(first)) bad.push(f);
    }
    const ok = bad.length === 0 && notes.size === 1;
    out.rec(
      'R2.5',
      ok ? 'PASS' : 'FAIL',
      'header-note position',
      bad.length
        ? JSON.stringify(bad)
        : 'distinct=' + notes.size + ' uniform ' + agents.length + '/' + agents.length
    );
  },
};

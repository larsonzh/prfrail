'use strict';

/**
 * G5-a role-file frontmatter hard gate (DELIVERY_DIRECTIVE 6.3 sub-rule G5-a).
 *   ① `git grep -l "^model:" -- .github/agents/prfrail-*.agent.md` must be empty;
 *   ② no `.agent.md` frontmatter may carry an empty `model` value in any of the
 *      three forms `^\s*model\s*:\s*$`, `model: ""`, `model: ''`.
 */
const EMPTY_FORMS = [/^\s*model\s*:\s*$/m, /model:\s*""/, /model:\s*''/];

module.exports = {
  id: 'G5-a',
  label: 'no model key in prfrail-*, no empty model',
  title: '角色文件 frontmatter 硬闸',
  scripted: true,

  subchecks: [
    {
      id: 'G5-a(1)',
      run(ctx, out) {
        const r = ctx.git(['grep', '-l', '^model:', '--', '.github/agents/prfrail-*.agent.md']);
        const h = r.ok ? r.out : '';
        out.rec(
          'G5-a(1)',
          h.trim() === '' ? 'PASS' : 'FAIL',
          'no model key in prfrail-*',
          h.trim() === '' ? '0 hits' : h
        );
      },
    },
    {
      id: 'G5-a(2)',
      run(ctx, out) {
        const agents = ctx.listAgentFiles();
        const hits = [];
        for (const f of agents) {
          const text = ctx.readText('.github/agents/' + f);
          if (text === null) continue;
          const fm = (/^---\n([\s\S]*?)\n---/.exec(text) || [])[1] || '';
          if (EMPTY_FORMS.some((re) => re.test(fm))) hits.push(f);
        }
        out.rec(
          'G5-a(2)',
          hits.length === 0 ? 'PASS' : 'FAIL',
          'no empty model',
          hits.length ? JSON.stringify(hits) : '0 hits / ' + agents.length + ' files'
        );
      },
    },
  ],

  run(ctx, out) {
    for (const s of this.subchecks) s.run(ctx, out);
  },
};

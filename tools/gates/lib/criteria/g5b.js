'use strict';

/**
 * G5-b role-file `model` value gate (DELIVERY_DIRECTIVE 6.3 sub-rule G5-b).
 *
 * Changed to a whitelist CLOSED SET by slice DIRECTIVE-MODEL-LISTS (v1.18, user ruling
 * 2026-10-03). A newly added or modified `.agent.md` file that carries a frontmatter
 * `model` key must set it to a member of the closed set below, which mirrors the
 * `modelId` column of section 2.2. The `gpt-6-luna` member landed only after the
 * L1..L5 tool-chain probe came back PASS (evidence: this slice validation report).
 *
 * The former blacklist regex is deliberately gone: after the section 2.4 convergence the
 * allow-list closed set is what bans everything else, so a family-name blacklist would
 * both under- and over-match.
 *
 * Reading: only the frontmatter block (the first `--- ... ---` pair) is parsed; ALL
 * `model` occurrences in it are collected; comparison is case-sensitive and verbatim
 * after trimming both ends; a quoted value is therefore not equal to a set member.
 * A missing key passes (section 2.6); an EMPTY value is G5-a 2's finding, not ours.
 */
const CLOSED_SET = [
  'deepseek-v4-pro',
  'deepseek-flash',
  'mai-code-1.1-flash',
  'gpt-5.3-codex',
  'gpt-6.1-sol',
  'gpt-6-luna',
];
const ALLOWED = new Set(CLOSED_SET);
const FRONTMATTER = /^---\r?\n([\s\S]*?)\r?\n---/;

/** Empty-value shapes belong to G5-a 2, so G5-b steps aside on all of them. */
const EMPTY_VALUE_FORMS = new Set(['', '""', "''"]);

function frontmatterOf(text) {
  const m = FRONTMATTER.exec(text);
  return m ? m[1] : '';
}

module.exports = {
  id: 'G5-b',
  label: 'closed-set model value',
  title: '角色文件 model 取值闸',
  scripted: true,
  closedSet: CLOSED_SET,

  run(ctx, out) {
    const targets = ctx.changedPaths.filter((f) => /\.agent\.md$/i.test(f));
    const findings = [];
    for (const f of targets) {
      const text = ctx.readText(f);
      if (text === null) continue;
      const fm = frontmatterOf(text);
      if (!fm) continue;
      const re = /^\s*model\s*:\s*(.+)$/gm;
      let m;
      while ((m = re.exec(fm)) !== null) {
        const v = m[1].trim();
        if (EMPTY_VALUE_FORMS.has(v)) continue;
        if (!ALLOWED.has(v)) findings.push({ file: f, value: v });
      }
    }
    out.rec(
      'G5-b',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'closed-set model value',
      findings.length === 0
        ? 'checked ' + targets.length + ' file(s)'
        : findings.map((x) => x.file + ' :: model=' + x.value).join('; ')
    );
    return { findings };
  },
};

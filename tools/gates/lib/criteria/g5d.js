'use strict';

/**
 * G5-d role-file frontmatter contract (DELIVERY_DIRECTIVE section 2.6, OB-30).
 *
 * Section 2.6 requires every role file to carry exactly five frontmatter keys:
 * `name` / `description` / `tools` / `user-invocable` / `target`. Nothing mechanical
 * checked that until this sub-criterion: G5-a / G5-b judge `model` and the change-set
 * contract, while the five-key minimum itself was prose only.
 *
 * Scope: `.github/agents/*.agent.md` files that are ADDED or MODIFIED by the change set
 * (deleted paths submit no bytes and are skipped). Same change-set discipline as
 * G5-a / G5-b / G5-c: historical files that the change set does not touch are never
 * re-judged.
 *
 * Predicate:
 *   - the file must open with a `---` frontmatter fence and close it;
 *   - each of the five keys must occur EXACTLY ONCE inside the frontmatter
 *     (missing, repeated, or an empty value fails);
 *   - extra keys are NOT judged - `model` in particular is owned by G5-a / G5-b, and
 *     future keys must not turn this criterion into a closed-set whitelist.
 *
 * Value emptiness: a bare `key:`, `key: ""` or `key: ''` counts as empty. An explicit
 * empty list (`key: []`) is a value, not an omission - it is left to G5-a / G5-b.
 *
 * An unreadable target file is a premise error (`UsageError`, exit 2), never a silent
 * pass: the criterion exists to prove the contract, not to skip it.
 */

const { UsageError } = require('../ctx');

const AGENT_PATH_RE = /^\.github\/agents\/[^/]+\.agent\.md$/i;
const REQUIRED_KEYS = ['name', 'description', 'tools', 'user-invocable', 'target'];
const FENCE = '---';
const KEY_LINE_RE = /^([A-Za-z0-9_-]+)\s*:(.*)$/;

function isEmptyValue(v) {
  const s = String(v === null || v === undefined ? '' : v).trim();
  return s === '' || s === '""' || s === "''";
}

/**
 * Parse the frontmatter of one role file.
 * Returns `{ ok, reason, keys: Map<key, count>, lines: Map<key, lineNumber> }` where
 * `ok === false` means the frontmatter itself is absent/unterminated.
 */
function parseFrontmatter(L) {
  const keys = new Map();
  const lines = new Map();
  if (!Array.isArray(L) || !L.length || String(L[0]).trim() !== FENCE) {
    return { ok: false, reason: 'frontmatter missing (the file must open with a `---` fence)', keys, lines };
  }
  let close = -1;
  for (let i = 1; i < L.length; i++) {
    if (String(L[i]).trim() === FENCE) {
      close = i;
      break;
    }
  }
  if (close < 0) {
    return { ok: false, reason: 'frontmatter unterminated (no closing `---`)', keys, lines };
  }
  for (let i = 1; i < close; i++) {
    const m = KEY_LINE_RE.exec(String(L[i]));
    if (!m) continue;
    const key = m[1];
    keys.set(key, (keys.get(key) || 0) + 1);
    if (!lines.has(key)) lines.set(key, i + 1);
  }
  return { ok: true, reason: '', keys, lines };
}

/** Every contract violation of one parsed role file, as `{ line, reason }` pairs. */
function violationsOf(L) {
  const fm = parseFrontmatter(L);
  if (!fm.ok) return [{ line: 1, reason: fm.reason }];
  const out = [];
  const values = new Map();
  for (let i = 1; i < L.length; i++) {
    if (String(L[i]).trim() === FENCE) break;
    const m = KEY_LINE_RE.exec(String(L[i]));
    if (m && !values.has(m[1])) values.set(m[1], m[2]);
  }
  for (const key of REQUIRED_KEYS) {
    const n = fm.keys.get(key) || 0;
    if (n === 0) {
      out.push({ line: 1, reason: 'required key `' + key + '` is missing' });
      continue;
    }
    if (n > 1) {
      out.push({ line: fm.lines.get(key), reason: 'required key `' + key + '` occurs ' + n + ' times (exactly once required)' });
      continue;
    }
    if (isEmptyValue(values.get(key))) {
      out.push({ line: fm.lines.get(key), reason: 'required key `' + key + '` has an empty value' });
    }
  }
  return out;
}

module.exports = {
  id: 'G5-d',
  label: 'role-file frontmatter contract',
  title: '角色文件 frontmatter 契约',
  scripted: true,
  internals: { AGENT_PATH_RE, REQUIRED_KEYS, parseFrontmatter, violationsOf, isEmptyValue },

  run(ctx, out) {
    const findings = [];
    const targets = (ctx.changedPaths || []).filter((p) => AGENT_PATH_RE.test(p));
    if (!targets.length) {
      out.rec('G5-d', 'INFO', 'role-file frontmatter contract', 'no changed `.github/agents/*.agent.md` in scope');
      return { findings };
    }
    let judged = 0;
    for (const f of targets) {
      const rec = typeof ctx.recordFor === 'function' ? ctx.recordFor(f) : null;
      if (rec && rec.deleted === true) continue;
      const L = typeof ctx.readLineArray === 'function' ? ctx.readLineArray(f) : null;
      if (!L) {
        throw new UsageError('G5-d: changed role file is unreadable: ' + f);
      }
      judged++;
      for (const v of violationsOf(L)) {
        findings.push({
          file: f,
          line: v.line,
          text: String(L[v.line - 1] === undefined ? '' : L[v.line - 1]),
          detail: f + ':' + v.line + ' ' + v.reason,
        });
      }
    }
    out.rec(
      'G5-d',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'role-file frontmatter contract',
      findings.length
        ? findings.map((x) => x.detail).join(' | ')
        : judged + ' changed role file(s) carry the five section 2.6 keys exactly once'
    );
    return { findings };
  },
};

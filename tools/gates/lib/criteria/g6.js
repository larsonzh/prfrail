'use strict';

/**
 * G6 ledger metadata consistency (architecture design section 2).
 *
 * Scope rule: only the ADDED lines of the change set (same wording as G1-a / G3(b) /
 * G4*); historical lines are never re-judged. Target files: every `.md` in scope.
 *
 * Self-reference guard: every sub-check skips lines carrying one of the exclusion
 * markers, so the section 6.3 rule text can quote the forbidden shapes safely
 * (isomorphic to the G2 "不得残留" mechanism).
 *
 * Whitelist: `tools/gates/g6-whitelist.txt`, entries
 * `<match><TAB><file><TAB><reason / context excerpt>`, substring match, one
 * registration per hit, reason and excerpt mandatory (a malformed entry is a usage
 * error, exit 2 - that is what stops the file from becoming a blanket). The file is
 * read through `ctx.loadWhitelist` (GATES-EXT F1): tree/index read the working tree,
 * range/ci read the scope's endpoint revision - the same rule G4b uses.
 *
 *   G6-1 status tense      pending marker and same-label completion marker in one block
 *   G6-2 cost clause       changed realized clause vs same-key realized clause, whole file
 *   G6-3 run-id format     placeholder / out-of-range run ids (format only, never
 *                          reachability; no network)
 *   G6-4 expiring literal  line counts / file counts / aggregate counts
 *   G6-5 section refs      §X.Y on added ledger/changelog lines must resolve (OB-32)
 *   G6-6 ledger status     §12.4 disposition cells must open with a closed-set word (OB-38)
 *   G6-7 declaration       a §0.3 declaration names an ID the ledger regions contain (OB-32)
 *   G6-8 DR numbering      mint rows unique, references registered (P24 GATES-CLEANUP)
 */

const fs = require('fs');
const path = require('path');

const { UsageError, decodeBuffer } = require('../ctx');

const WHITELIST_REL = 'tools/gates/g6-whitelist.txt';
const EXCLUDE_MARKERS = ['不得残留', '示例', '引文', '冲突', '矛盾'];

const PENDING_RE = /待\s*(?:①|②|③|④|⑤|⑥|⑦|⑧a|⑧b|⑨|⑩|终审|独立扫描|复审|预审|盲审|提交授权)/g;
/**
 * G6-1 completion marker (GATES-EXT, ⑥ F1). The accepted set must be the UNION the
 * directive enumerates - `实得括号形态 / 粗体实得 / 已闭合 / ✅ 已完成` - not the
 * bold-only intersection of it. Six accepted shapes:
 *     （实得）   （**实得**）   **实得**   已闭合   ✅ 已完成   ✅ **已完成**
 * Before the fix the regex demanded `**...**` around 实得, so the directive's OWN
 * worked example (`待 ⑥ 独立扫描` coexisting with `⑥ 独立扫描（实得）PASS`, no bold)
 * was not caught: a fail-open in a CI-blocking gate against its own text. Both the CN
 * rule line and its `_EN` mirror contain `示例`, so the widened set cannot
 * re-introduce a self-reference finding. The same-label 40-character window below is
 * unchanged and still required, so widening does not weaken the block-level guard.
 */
const DONE_RE = /（\*\*实得\*\*）|（实得）|\*\*实得\*\*|已闭合|✅\s*\*\*已完成\*\*|✅\s*已完成/g;

const KEY_RE = /SW-\d+|DR-\d+|OB-\d+|GATES-EXT|[A-Z]\d+[a-z]?/g;
const COST_RE = /([⑥⑦])\s*[×xX]\s*(\d+)/g;
const TOTAL_RE = /=\s*(\d+)\s*次/g;

/**
 * G6-3 legal shape (DELIVERY_DIRECTIVE 6.3): `run` + one space + 9-13 decimal digits.
 * Every entry below is a must-FAIL shape; the band is enforced from BOTH sides, so a
 * truncated / mis-copied id of 3-8 digits can no longer pass (GATES-EXT F7: the code
 * previously only rejected 1-2 digits and >= 14 digits, so the criterion failed open
 * against its own text).
 *
 * Prose narrowing (batch 2e): a digit run IMMEDIATELY followed by a word (ASCII letter
 * or CJK ideograph) is a count in ordinary prose, not a run id - `run 1234 times` and
 * `连续 run 12 次` are not malformed ids. The negative lookahead `(?!\s*[A-Za-z\u4e00-\u9fff])`
 * is therefore attached to the DIGIT-RUN rules only. A token-final or backticked digit
 * run stays judged (`run 111` end of line, `run 1234。`, `run \`1234\``, `\`run 1234\``),
 * because the lookahead sees no following word. The placeholder rules (angle brackets,
 * word list, CJK wording) are UNAMBIGUOUS and carry no following-word restriction - they
 * are checked in any context, prose included.
 *
 * Leading boundary (batch 2f, NF-2): EVERY rule below carries the leading guard
 * `(?<![A-Za-z\u4e00-\u9fff])`, so `run` preceded by a word character is not a run-id
 * mention at all - `rerun 1234`, `prerun 12` and the CJK-glued `在run 1234` are not
 * judged. Before 2f the patterns also matched inside words, which is a false-positive
 * class in a CI-blocking gate (a sentence like "the harness will rerun 1234" would have
 * turned G6-3 red). The guard is a LOOKBEHIND, so it never changes what is matched: the
 * still-judged shapes (`see run 1234`, `（run 1234）`, `\`run 1234\``, `run 1234。`,
 * `run TBD`) are followed by a non-word character or the start of the line, exactly as
 * before. Placeholder rules carry the guard too: a rule reading `run` should not fire on
 * a word that merely ends in those three letters.
 * OB-17 (slice GATES-TIGHTEN, v1.33): the angle-bracket rule is narrowed to id-LIKE
 * placeholders - first character uppercase, the rest uppercase / digits / `-` / `_`.
 * `run <RUN_ID>` still fails, while CLI usage examples such as `run <chain-file>` (token
 * starts lowercase) are no longer false positives. The rule stays context-free, so a
 * line carrying the usage example no longer turns G6-3 red merely because it was edited.
 */
const RUN_PATTERNS = [
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s+0+(?![0-9])(?!\s*[A-Za-z\u4e00-\u9fff])/g, why: 'all-zero run id' },
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s*<[A-Z][A-Z0-9_-]*>/g, why: 'angle-bracket placeholder' },
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s+(?:TBD|TODO|XXXX|XXX|N\/A)\b/gi, why: 'word placeholder' },
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s+(?:待|占位)/g, why: 'placeholder wording' },
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s+\d{1,2}(?![0-9])(?!\s*[A-Za-z\u4e00-\u9fff])/g, why: 'run id too short' },
  {
    re: /(?<![A-Za-z\u4e00-\u9fff])run\s+\d{3,8}(?![0-9])(?!\s*[A-Za-z\u4e00-\u9fff])/g,
    why: 'run id below the 9-digit floor',
  },
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s+\d{14,}(?![0-9])(?!\s*[A-Za-z\u4e00-\u9fff])/g, why: 'run id too long' },
];

const LITERAL_PATTERNS = [
  { re: /\d{2,}\/\d{2,}\s*行/, why: 'line-count ratio' },
  { re: /\d{2,}\s*行/, why: 'line count' },
  { re: /\d{2,}\s*个?\s*文件/, why: 'file count' },
  { re: /共\s*\d{2,}/, why: 'aggregate count' },
];

function isExcluded(line) {
  return EXCLUDE_MARKERS.some((m) => line.includes(m));
}

/** Heading indices of `^#{2,4} ` lines (block boundaries for G6-1). */
function headingIndices(L) {
  const out = [];
  for (let i = 0; i < L.length; i++) if (/^#{2,4} /.test(L[i])) out.push(i);
  return out;
}

function blockBounds(headings, i, total) {
  let start = 0;
  let end = total;
  for (const h of headings) {
    if (h <= i) start = h;
    else {
      end = h;
      break;
    }
  }
  return [start, end];
}

function qualify(line, index) {
  // Clause qualification: the token belongs to the clause introduced by the NEAREST
  // preceding qualifier (`预算` / `实得` / bare `成本`) inside the 60-character window.
  // A plain `contains` test would let a `预算` keyword anywhere in the window shadow a
  // following `实得` clause, which would make the realized-clause comparison vacuous
  // (the architecture's T3/T4 expectations both require the 实得 clause to be compared).
  const pre = line.slice(Math.max(0, index - 60), index);
  let last = null;
  for (const m of pre.matchAll(/预算|实得|成本/g)) last = m[0];
  if (last === '预算') return 'budget';
  if (last === '实得' || last === '成本') return 'realized';
  return null;
}

/** Realized cost vector of one line: ⑥/⑦ counts plus the `= N 次` total. */
function costClause(line) {
  const realized = new Map();
  let has = false;
  for (const m of line.matchAll(COST_RE)) {
    if (qualify(line, m.index) !== 'realized') continue;
    has = true;
    realized.set(m[1], (realized.get(m[1]) || 0) + Number(m[2]));
  }
  if (has) {
    for (const m of line.matchAll(TOTAL_RE)) {
      if (qualify(line, m.index) === 'realized') realized.set('=', Number(m[1]));
    }
  }
  return realized;
}

function vecText(map) {
  return (
    [...map.entries()]
      .sort((a, b) => (a[0] < b[0] ? -1 : 1))
      .map(([k, v]) => (k === '=' ? 'total=' + v : k + '×' + v))
      .join('+') || '-'
  );
}

function keysOfLine(line) {
  return [...new Set([...line.matchAll(KEY_RE)].map((m) => m[0]))];
}

function g6Whitelist(ctx) {
  let res;
  try {
    res = ctx.loadWhitelist
      ? ctx.loadWhitelist(WHITELIST_REL)
      : { entries: [], error: null, source: 'worktree' };
  } catch (e) {
    throw new UsageError('cannot read ' + WHITELIST_REL + ': ' + (e.message || e));
  }
  if (!res) res = { entries: [], error: null, source: 'worktree' };
  if (res.error) throw new UsageError(res.error);
  return { entries: res.entries, source: res.source || 'unknown' };
}

/** Split findings into un-whitelisted (fatal) and whitelisted (SKIP, exit-code neutral). */
function splitByWhitelist(entries, findings) {
  const fatal = [];
  const exempt = [];
  for (const f of findings) {
    const hit = entries.filter((e) => (e.file === '*' || e.file === f.file) && String(f.text).includes(e.match));
    if (hit.length) exempt.push({ finding: f, entry: hit[0] });
    else fatal.push(f);
  }
  return { fatal, exempt };
}

function emit(ctx, out, id, label, findings, emptyDetail) {
  const wl = g6Whitelist(ctx);
  const { fatal, exempt } = splitByWhitelist(wl.entries, findings);
  const src = ' whitelist-source=' + wl.source;
  out.rec(
    id,
    fatal.length === 0 ? 'PASS' : 'FAIL',
    label,
    (findings.length === 0
      ? emptyDetail
      : fatal.length
      ? fatal.map((f) => f.detail).join(' | ')
      : 'exempt by whitelist (' + exempt.length + ')') + src
  );
  for (const e of exempt) {
    out.rec(
      id,
      'SKIP',
      label + ' (whitelist)',
      e.finding.file + ':' + e.finding.line + ' ' + e.finding.text.slice(0, 60) + ' <- ' + e.entry.reason + src
    );
  }
  return { findings, fatal, exempt };
}

const G6_1 = {
  id: 'G6-1',
  label: 'status tense',
  run(ctx, out) {
    const findings = [];
    const seen = new Set();
    let scanned = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const headings = headingIndices(L);
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined) continue;
        scanned++;
        if (isExcluded(line)) continue;
        for (const m of line.matchAll(PENDING_RE)) {
          const label = m[0].replace(/^待\s*/, '');
          const [bs, be] = blockBounds(headings, n - 1, L.length);
          for (let i = bs; i < be; i++) {
            const dl = L[i];
            if (dl === undefined || isExcluded(dl)) continue;
            for (const dm of dl.matchAll(DONE_RE)) {
              const win = dl.slice(Math.max(0, dm.index - 40), dm.index + dm[0].length + 40);
              if (!win.includes(label)) continue;
              const key = f + ':' + n + ':' + label + ':' + (i + 1);
              if (seen.has(key)) continue;
              seen.add(key);
              findings.push({
                file: f,
                line: n,
                text: line,
                detail: f + ':' + n + ' pending(' + label + ') vs done@' + (i + 1) + ' ' + dm[0],
              });
            }
          }
        }
      }
    }
    return emit(ctx, out, 'G6-1', 'status tense', findings, 'no pending/completion conflict on added lines (scanned ' + scanned + ')');
  },
};

const G6_2 = {
  id: 'G6-2',
  label: 'cost clause',
  run(ctx, out) {
    const findings = [];
    let checked = 0;
    let scanned = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const headings = headingIndices(L);
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined) continue;
        scanned++;
        if (isExcluded(line)) continue;
        const vec = costClause(line);
        if (!vec.size) continue;
        checked++;
        const myKeys = keysOfLine(line).length
          ? keysOfLine(line)
          : (() => {
              const [bs] = blockBounds(headings, n - 1, L.length);
              return keysOfLine(L[bs] || '');
            })();
        if (!myKeys.length) continue;
        for (let j = 0; j < L.length; j++) {
          if (j === n - 1) continue;
          const other = L[j];
          if (other === undefined || isExcluded(other)) continue;
          const ovec = costClause(other);
          if (!ovec.size) continue;
          const oKeys = keysOfLine(other).length
            ? keysOfLine(other)
            : (() => {
                const [bs] = blockBounds(headings, j, L.length);
                return keysOfLine(L[bs] || '');
              })();
          const shared = oKeys.filter((k) => myKeys.includes(k));
          if (!shared.length) continue;
          if (vecText(ovec) === vecText(vec)) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail:
              f +
              ':' +
              n +
              ' ' +
              vecText(vec) +
              ' vs line ' +
              (j + 1) +
              ' ' +
              vecText(ovec) +
              ' (key ' +
              shared.join(',') +
              ')',
          });
        }
      }
    }
    return emit(
      ctx,
      out,
      'G6-2',
      'cost clause',
      findings,
      checked
        ? checked + ' realized clause line(s) consistent with their same-key occurrences (scanned ' + scanned + ')'
        : 'no realized cost clause on added lines (scanned ' + scanned + ')'
    );
  },
};

const G6_3 = {
  id: 'G6-3',
  label: 'run-id format',
  run(ctx, out) {
    const findings = [];
    let scanned = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined) continue;
        scanned++;
        if (isExcluded(line)) continue;
        for (const p of RUN_PATTERNS) {
          p.re.lastIndex = 0;
          const m = p.re.exec(line);
          if (!m) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail: f + ':' + n + ' ' + p.why + ' :: ' + m[0].trim(),
          });
          break;
        }
      }
    }
    return emit(
      ctx,
      out,
      'G6-3',
      'run-id format',
      findings,
      'no malformed run id on added lines (scanned ' + scanned + ')'
    );
  },
};

const G6_4 = {
  id: 'G6-4',
  label: 'expiring literal',
  run(ctx, out) {
    const findings = [];
    let scanned = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined) continue;
        scanned++;
        if (isExcluded(line)) continue;
        for (const p of LITERAL_PATTERNS) {
          const m = p.re.exec(line);
          if (!m) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail: f + ':' + n + ' ' + p.why + ' :: ' + m[0].trim(),
          });
          break;
        }
      }
    }
    return emit(
      ctx,
      out,
      'G6-4',
      'expiring literal',
      findings,
      'no expiring literal on added lines (scanned ' + scanned + ')'
    );
  },
};

/* ---------------------------------------------------------------------------
 * G6-5 (section-reference resolution, OB-32) and G6-6 (ledger status closed set,
 * OB-38) - the two sub-checks added by slice DIRECTIVE-REVIEW-SATURATION (v1.13).
 *
 * G6-5 is deliberately an APPROXIMATE criterion: it covers the "dangling reference"
 * subclass of OB-32's C1 event (a log row naming a section that does not exist) and
 * does NOT cover "the log claims a change while the body has no such content". The
 * rule text in the directive says so explicitly; this comment says so again so the
 * code cannot be read as a full-class cover.
 *
 * G6-6 is independent of G6-1 on purpose: G6-1's block-level semantics (a pending
 * marker and a same-label completion marker in one block) are structurally
 * incompatible with §12.4's single large table, so widening G6-1's vocabulary would
 * make any "待办" row collide with any "已处置" row. A separate sub-check with its own
 * closed set is the design recorded in §6.3 and §7.8.
 * ------------------------------------------------------------------------ */

/** `§12.4`, `§7.2.1` … - the reference forms the directive uses in its prose. */
const SECTION_REF_RE = /§(\d+(?:\.\d+)*)/g;
/** Label source 1: a heading that starts with the section number (`### 7.2.1 触发`). */
const LABEL_HEADING_RE = /^#{2,4}\s+(\d+(?:\.\d+)*)(?=[\s.．、]|$)/;
/** Label source 2: a line-leading bold label (`**1.7.1 执行自由度**`), as §1.7 / §3.4 use. */
const LABEL_BOLD_RE = /^\*\*(\d+(?:\.\d+)*)\s/;
const DIRECTIVE_FILES = ['docs/DELIVERY_DIRECTIVE.md', 'docs/DELIVERY_DIRECTIVE_EN.md'];

const REGION_HEAD_030 = /^#{2,4}\s+0\.3(?=[\s.．、]|$)/;
const REGION_HEAD_C0 = /^#{2,4}\s+C\.0(?=[\s.．、]|$)/;
const REGION_HEAD_C1 = /^#{2,4}\s+C\.1(?=[\s.．、]|$)/;
const ANY_HEAD = /^#{2,4}\s/;

/**
 * Every label a directive uses to address its own sections: heading numbers plus the
 * leading bold labels (§1.7.1 / §3.4.1 style). Collection method copied from the
 * reviewed throwaway `tmp/v112/refs-check.js` so the criterion and the ⑨ native
 * validation agree on what "the file's label set" means.
 */
function labelsFromLines(L) {
  const out = new Set();
  for (const l of L || []) {
    let m = LABEL_HEADING_RE.exec(l);
    if (m) out.add(m[1]);
    m = LABEL_BOLD_RE.exec(l);
    if (m) out.add(m[1]);
  }
  return out;
}

/**
 * Scope-aware read of one repository file, with a working-tree fallback.
 * `ctx.readText` is the scope-correct reader (tree/index = worktree, range/ci = the
 * endpoint revision), so a historical run judges against the labels that existed at
 * the endpoint. The selftest's synthetic context returns `null` for repository files
 * and is served by the disk fallback - the fixtures deliberately live inside the
 * repository so their references can be checked against the real directive.
 */
function readTextFor(ctx, p) {
  let t = null;
  try {
    t = ctx.readText ? ctx.readText(p) : null;
  } catch (e) {
    t = null;
  }
  if (t === null || t === undefined) {
    try {
      t = decodeBuffer(fs.readFileSync(path.join(ctx.repoRoot, p.split('/').join(path.sep))));
    } catch (e) {
      t = null;
    }
  }
  return t === null || t === undefined ? null : t;
}

/** The union of the CN and `_EN` directive label sets ("this file ∪ its mirror"). */
function directiveLabels(ctx) {
  const out = new Set();
  for (const f of DIRECTIVE_FILES) {
    const t = readTextFor(ctx, f);
    if (t === null) continue;
    for (const l of labelsFromLines(t.split('\n'))) out.add(l);
  }
  return out;
}

/**
 * Target regions of G6-5, as [startIndex, endIndex) pairs over the file's lines:
 *   ① the §0.3 change-log table interval (its heading to the next `^#{2,4}` heading);
 *   ② the appendix C ledger interval (`### C.0` up to `### C.1`, exclusive) - this is
 *      what keeps the C.1 v3.1 comparison table and C.2 / C.3 structurally OUT of the
 *      scan (measured: 15 of the 18 unresolvable whole-file references live there, so
 *      a region limit - not a whitelist - is what makes a zero-preload launch possible).
 *
 * OB-46 (v1.33) endpoint self-check (coverage boundary stated by the 5 pre-review, L-6): only
 * region (2) has a delimiter that could go missing, so only region (2) plus the EXISTENCE of
 * both regions are asserted. Region (1) ends at the next `#{2,4}` heading of any level; losing
 * the section 0.3 heading entirely would drop the region silently FOR A NON-DIRECTIVE FILE (the
 * directive pair is covered by the two-region existence assertion below), and that residual is carried
 * by the 5 / 6 / 7 review chain (the criterion text claims no more than it checks).
 *
 * OB-46 (v1.33) endpoint self-check: `### C.1` is the delimiter that makes region ② a
 * region. Without it the loop above fell through to `L.length` and judged an UNBOUNDED
 * interval - the C.1 comparison table and C.2 / C.3 came back into scope, so an unrelated
 * historical row decided the verdict (before OB-46 the boundary mutation in T21 turned
 * G6-5 red through the §3.10 row, not through the missing endpoint). A region we cannot
 * delimit is a bug in the criterion's premise, not a finding: fail closed (exit 2).
 */
function targetRegionRanges(L) {
  const out = [];
  for (let i = 0; i < L.length; i++) {
    if (!REGION_HEAD_030.test(L[i])) continue;
    let end = L.length;
    for (let j = i + 1; j < L.length; j++) {
      if (ANY_HEAD.test(L[j])) {
        end = j;
        break;
      }
    }
    out.push([i, end]);
    break;
  }
  for (let i = 0; i < L.length; i++) {
    if (!REGION_HEAD_C0.test(L[i])) continue;
    let end = L.length;
    let terminated = false;
    for (let j = i + 1; j < L.length; j++) {
      if (REGION_HEAD_C1.test(L[j])) {
        end = j;
        terminated = true;
        break;
      }
    }
    if (!terminated) {
      throw new UsageError(
        'G6-5 target region endpoint missing: no `### C.1` heading after `### C.0` (line ' +
          (i + 1) +
          '); refusing to judge an unbounded region'
      );
    }
    out.push([i, end]);
    break;
  }
  return out;
}

/** A reference resolves when the number exists, or when a label extends it (`§7.2` ← `§7.2.1`). */
function resolvesRef(num, labels) {
  if (labels.has(num)) return true;
  for (const l of labels) if (l.indexOf(num + '.') === 0) return true;
  return false;
}

const G6_5 = {
  id: 'G6-5',
  label: 'section-reference resolution',
  run(ctx, out) {
    const findings = [];
    const external = directiveLabels(ctx);
    let scanned = 0;
    let refs = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const ranges = targetRegionRanges(L);
      // OB-46 (v1.33): the directive pair is the only file class that MUST carry both
      // target regions. A renamed `## 0.3` or `### C.0` heading used to drop that region
      // silently - the file still got a verdict (PASS) while that region was judged by
      // nothing. Requiring both here closes the "empty region" half of the same defect.
      if (DIRECTIVE_FILES.indexOf(f) >= 0 && ranges.length !== 2) {
        throw new UsageError(
          'G6-5 target regions missing in ' +
            f +
            ': expected the 0.3 change-log region and the appendix C region, got ' +
            ranges.length
        );
      }
      if (!ranges.length) continue;
      const own = labelsFromLines(L);
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const idx = n - 1;
        const line = L[idx];
        if (line === undefined) continue;
        // A line outside every target region is not judged at all (reported as scanned-out,
        // never as a PASS): historical rows in §0.3 and the C.1 comparison table are
        // therefore immune by construction, not by whitelist.
        if (!ranges.some((r) => idx >= r[0] && idx < r[1])) continue;
        scanned++;
        if (isExcluded(line)) continue;
        for (const m of line.matchAll(SECTION_REF_RE)) {
          refs++;
          if (resolvesRef(m[1], external) || resolvesRef(m[1], own)) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail: f + ':' + n + ' §' + m[1] + ' unresolved :: ' + line.slice(0, 120),
          });
        }
      }
    }
    return emit(
      ctx,
      out,
      'G6-5',
      'section-reference resolution',
      findings,
      'no unresolved section reference on added ledger/changelog lines (scanned ' + scanned + ', refs ' + refs + ')'
    );
  },
};

/** The bilingual closed set. The CN word is authoritative; `_EN` files use the EN table. */
const LEDGER_STATUS_SETS = {
  cn: ['观察中', '待办', '待议', '已处置', '已裁决', '已登记'],
  en: ['Observing', 'To do', 'To be discussed', 'Resolved', 'Ruling', 'Registered'],
};
const LEDGER_REGION_HEAD = /^#{2,4}\s+12\.4(?=[\s.．、]|$)/;
const LEDGER_ROW_RE = /^\|\s*OB-\d+\s*\|/;

/** The §12.4 table interval as [startIndex, endIndex), or null when the file has none. */
function ledgerRegionRange(L) {
  for (let i = 0; i < L.length; i++) {
    if (!LEDGER_REGION_HEAD.test(L[i])) continue;
    let end = L.length;
    for (let j = i + 1; j < L.length; j++) {
      if (ANY_HEAD.test(L[j])) {
        end = j;
        break;
      }
    }
    return [i, end];
  }
  return null;
}

/** Cells of a Markdown table row, splitting on unescaped pipes (escaped `\|` stays intact). */
function rowCellsOf(line) {
  const cells = [];
  let cur = '';
  for (let i = 0; i < line.length; i++) {
    const ch = line[i];
    if (ch === '\\' && i + 1 < line.length) {
      cur += ch + line[i + 1];
      i++;
      continue;
    }
    if (ch === '|') {
      cells.push(cur);
      cur = '';
      continue;
    }
    cur += ch;
  }
  cells.push(cur);
  return cells.slice(1, cells.length - 1);
}

/**
 * The disposition cell's first status word, or null when it is outside the closed set.
 * Accepted shapes: `词 …`, `**词** …`, `**词**（括注）…` / `**Word** (note) …`. The
 * boundary test demands a non-letter right after the word, so a legacy value such as
 * `**下一片必做**` (which merely CONTAINS no closed word) fails, while `**已处置**（v1.1）`
 * passes.
 */
function firstStatusWord(cell, words) {
  let c = String(cell === null || cell === undefined ? '' : cell).trim();
  if (c.slice(0, 2) === '**') c = c.slice(2);
  for (const w of words) {
    if (c.slice(0, w.length) !== w) continue;
    const rest = c.slice(w.length);
    if (rest === '' || rest.slice(0, 2) === '**' || !/[\p{L}\p{Script=Han}]/u.test(rest[0])) return w;
  }
  return null;
}

/**
 * OB-47 (v1.33): the language the document DECLARES in its header block.
 * The word table used to be picked with `/_EN\.md$/i.test(f)` - a property of the FILE
 * NAME, not of the document. A rename, or a third language added later, therefore changed
 * (or silently mismatched) the judged vocabulary with no assertion able to notice.
 * The rule now reads the declaration and cross-checks it against the name: the two must
 * agree, so a rename is an error instead of a silent table switch. All three failure modes
 * (absent / outside the closed set / disagreeing with the name) are premise errors, not
 * findings - fail closed.
 */
/**
 * The tag token is consumed to its FULL extent (v1.34 / OB-91, extended by OB-97 in v1.40).
 * Two partial fixes taught the rule the same lesson twice: 7 r1 caught `lang:zh-cn` read as
 * a legal `zh` prefix, and 7 r3 then caught `lang:zh_cn` collapsing to the legal `lang:zh`
 * while a bare `lang:` matched nothing at all - a character-class patch can always be
 * outgrown by the next delimiter. The rule is therefore structural: a tag STARTS at every
 * `\blang:` occurrence in the header window and runs over the maximal run of continuation
 * characters (letters, digits, `_`, `.`, `-`, and - since OB-97 - `:`, which is what used
 * to truncate `lang:zh:en` / `lang:zh:` into the legal prefix `lang:zh`). Every other
 * character - space, tab, comma, backtick, CJK punctuation, end of line, anything not
 * listed - TERMINATES the token, so the terminator set is the COMPLEMENT of the listed
 * continuation set instead of an enumeration a future delimiter can escape. The token is
 * then compared WHOLE against the closed set, and an occurrence preceded by a word
 * character (`xlang:zh`) is not an occurrence at all.
 * Coverage invariant: every occurrence inside the window is consumed by exactly one judged
 * token, so no `lang:`-prefixed form can escape judgement.
 *
 * The `-` is written LAST in the class on purpose: `[A-Za-z0-9._-:]` is an invalid
 * JavaScript character-class range (`_` to `:` is out of order) and throws a SyntaxError,
 * so the set is spelled `[A-Za-z0-9._:-]` - the same membership, dash at the end.
 */
const LANG_TAG_OCCURRENCE_RE = /\blang:/g;
const LANG_TAG_CONT_RE = /[A-Za-z0-9._:-]/;
const LANG_TAG_EXACT = { 'lang:zh': 'zh', 'lang:en': 'en' };
const LANG_HEADER_WINDOW = 15;
const LANG_TAG_TO_SET = { zh: 'cn', en: 'en' };

function declaredLang(L, f) {
  const seen = new Set();
  let raw = null;
  for (let i = 0; i < Math.min(L.length, LANG_HEADER_WINDOW); i++) {
    const line = L[i];
    LANG_TAG_OCCURRENCE_RE.lastIndex = 0;
    let m;
    while ((m = LANG_TAG_OCCURRENCE_RE.exec(line)) !== null) {
      let j = m.index + m[0].length;
      while (j < line.length && LANG_TAG_CONT_RE.test(line[j])) j++;
      const token = line.slice(m.index, j);
      LANG_TAG_OCCURRENCE_RE.lastIndex = j;
      if (Object.prototype.hasOwnProperty.call(LANG_TAG_EXACT, token)) seen.add(LANG_TAG_EXACT[token]);
      else if (raw === null) raw = token;
    }
  }
  // 7 r2 Medium (fixed in place): ANY non-closed-set tag is a premise error, even when a
  // legal tag is present as well. The former shape only complained when NOTHING matched, so
  // `lang:zh` + `lang:zh-cn` slipped through - fail-open in a rule that claims to be
  // fail-closed.
  if (raw !== null) {
    throw new UsageError(
      'unknown language tag `' +
        raw +
        '` in the header of ' +
        f +
        ': only `lang:zh` and `lang:en` are in the closed set'
    );
  }
  if (seen.size === 0) {
    throw new UsageError(
      'language tag missing in the header of ' +
        f +
        ': declare `lang:zh` or `lang:en` inside the header window (implementation constant LANG_HEADER_WINDOW)'
    );
  }
  if (seen.size > 1) {
    throw new UsageError(
      'conflicting language tags in the header of ' + f + ': ' + [...seen].map((x) => 'lang:' + x).join(' / ')
    );
  }
  const declared = [...seen][0];
  const named = /_EN\.md$/i.test(f) ? 'en' : 'zh';
  if (declared !== named) {
    throw new UsageError(
      'declared language `lang:' +
        declared +
        '` disagrees with the file name of ' +
        f +
        ' (the name implies lang:' +
        named +
        ')'
    );
  }
  return LANG_TAG_TO_SET[declared];
}

/* ---------------------------------------------------------------------------
 * G6-7 (declaration / ledger reconciliation, OB-32) - the content-level companion of
 * G6-5. G6-5 only proves that a `§X.Y` reference resolves; it does not prove that a
 * changelog line claiming "OB-N registered / closed / landed" corresponds to a ledger
 * row that actually exists. G6-7 closes that gap for ID keys.
 *
 * Predicate: on every ADDED line inside the §0.3 change-log region (the same
 * `targetRegionRanges` helper as G6-5, including its OB-46 endpoint self-check), an ID
 * key (`OB-\d+` / `SW-\d+` / `DR-\d+` / `GATES-[A-Z]+`) that sits inside a ±40-character
 * window of a declaration verb (CN: 新增 / 登记 / 补登 / 收口 / 翻 / 处置 / 更正 / 闭合 /
 * 落地; EN: added / registered / closed / resolved / corrected / landed / flipped) must
 * ALSO appear in the SAME file's §12.4 ledger region or appendix C region (C.0 up to
 * C.1, exclusive). Otherwise the changelog names an ID that has no ledger row at all - FAIL.
 *
 * NARROWING RECORD (prerequisite gate 3, v1.40). The first form of the predicate
 * demanded the ID on an ADDED line of §12.4 / appendix C ("the ledger row was touched
 * by THIS change set") and also judged appendix-C declaration lines. Measured over every
 * commit that touches the directive pair (50 commits, 96 file-runs), that form produced
 * 63 findings, all false positives: historical changelog rows routinely summarise a
 * whole slice by naming IDs whose ledger rows already existed before the commit (46 of
 * the 63 survive even when appendix C is dropped from the target set). The predicate was
 * therefore narrowed to (a) §0.3 targets only - appendix-C blocks are summary prose and
 * legitimately mention DR / SW / slice IDs that never live in §12.4 - and (b) evidence
 * from the whole ledger regions (any line, excluding the judging line), not only added
 * lines. Re-measured: 0 findings over the same 96 file-runs. The residual class this
 * still catches is the one OB-32 named: a changelog line that declares a ledger ID the
 * document's ledgers never contain.
 *
 * Precision choices (stated, not hidden):
 *   - a line never serves as its own evidence (defensive: the judging line is outside the
 *     evidence regions today, so this only guards a future target-set change).
 *   - a NEGATED verb is not a declaration: `不新增 / 未登记 / 无处置` and `not added /
 *     no longer registered` describe the absence of a change (measured false positive:
 *     `本片不新增、不删除、不改动 …；OB-53 保持待议，本片不动`).
 *   - the five G6 exclusion markers skip a line, so the rule text can quote the shapes.
 *   - `g6-whitelist.txt` exempts a finding one registration at a time (shared `emit`).
 *   - the predicate presupposes a LEDGER: a document with neither a §12.4 region nor an
 *     appendix C region has nothing to reconcile a declaration against, so it is not
 *     judged. Without this premise a cost-clause fixture carrying a `### 0.3` heading
 *     produced a false positive (measured on `g6-2-cost-clauses`, whose `新增 SW-5` line
 *     is a cost vector, not a ledger claim). The directive pair keeps both regions
 *     mandatory and additionally requires §12.4 (fail-closed, exit 2).
 * ------------------------------------------------------------------------ */
const DECL_VERB_RE =
  /新增|登记|补登|收口|翻|处置|更正|闭合|落地|\b(?:added|registered|closed|resolved|corrected|landed|flipped)\b/gi;
// Token boundary on the LEFT of an ID key (v1.40). Without it a longer token is carved
// up: the v1.40 §0.3 line names `ADR-027`, and the boundary-less form minted a fourth,
// bogus `DR-027` declaration. The same constant backs `ledgerIds`, so declaration and
// evidence faces widen and narrow together by construction.
const DECL_ID_RE = /(?<![A-Za-z0-9_-])(?:OB-\d+|SW-\d+|DR-\d+|GATES-[A-Z]+)/g;
const DECL_WINDOW = 40;

/**
 * A verb preceded by a negator (`不/未/无` with optional spaces, or `not` / `no`) is a
 * statement that the change did NOT happen; it is not a declaration. Keep the guard
 * deliberately narrow (immediately preceding token only), so `新增` inside a normal
 * sentence is still a declaration.
 */
function isNegatedVerb(text, index) {
  const pre = text.slice(Math.max(0, index - 5), index);
  return /(?:不|未|无)\s*$/.test(pre) || /(?:not|no)\s+$/i.test(pre);
}

/**
 * Declaration pairs of one line: every ID key whose ±DECL_WINDOW window carries a
 * non-negated declaration verb. `index` is the ID's offset (kept for diagnostics /
 * mutation probes).
 */
function declarationPairs(line) {
  const text = String(line);
  const out = [];
  for (const m of text.matchAll(DECL_ID_RE)) {
    const lo = Math.max(0, m.index - DECL_WINDOW);
    const hi = Math.min(text.length, m.index + m[0].length + DECL_WINDOW);
    const win = text.slice(lo, hi);
    let verb = null;
    for (const vm of win.matchAll(DECL_VERB_RE)) {
      if (isNegatedVerb(text, lo + vm.index)) continue;
      verb = vm[0];
      break;
    }
    if (!verb) continue;
    out.push({ id: m[0], verb: verb, index: m.index });
  }
  return out;
}

/** Map ID -> every line number (1-based) inside the given [start, end) regions. */
function ledgerIds(L, ranges) {
  const out = new Map();
  for (const r of ranges) {
    for (let i = r[0]; i < r[1]; i++) {
      const line = L[i];
      if (line === undefined) continue;
      for (const m of String(line).matchAll(DECL_ID_RE)) {
        if (!out.has(m[0])) out.set(m[0], []);
        out.get(m[0]).push(i + 1);
      }
    }
  }
  return out;
}

const G6_7 = {
  id: 'G6-7',
  label: 'declaration/ledger reconciliation',
  run(ctx, out) {
    const findings = [];
    let scanned = 0;
    let declarations = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      // Same region helper and same OB-46 fail-closed endpoint self-check as G6-5.
      const ranges = targetRegionRanges(L);
      if (DIRECTIVE_FILES.indexOf(f) >= 0 && ranges.length !== 2) {
        throw new UsageError(
          'G6-7 target regions missing in ' +
            f +
            ': expected the 0.3 change-log region and the appendix C region, got ' +
            ranges.length
        );
      }
      const logRange = ranges.find((r) => REGION_HEAD_030.test(L[r[0]])) || null;
      const apxRange = ranges.find((r) => REGION_HEAD_C0.test(L[r[0]])) || null;
      const ledger = ledgerRegionRange(L);
      // Premise: a declaration can only be reconciled against a ledger. The directive
      // pair is fail-closed (both target regions were asserted above; a missing §12.4
      // would make every declaration unreconcilable); any other document is judged only
      // when it actually carries a ledger region.
      if (DIRECTIVE_FILES.indexOf(f) >= 0 && !ledger) {
        throw new UsageError(
          'G6-7 ledger region missing in ' + f + ': the 12.4 observation ledger is the reconciliation target'
        );
      }
      if (!ledger && !apxRange) continue;
      // §0.3 only (see NARROWING RECORD above): appendix-C declaration lines are summary
      // prose and are not judged; the appendix region remains valid EVIDENCE.
      if (!logRange) continue;
      const evidenceRanges = [ledger, apxRange].filter(Boolean);
      const added = [...ctx.addedLineNumbers(f)].sort((a, b) => a - b);
      const evidenceIds = ledgerIds(L, evidenceRanges);
      for (const n of added) {
        const idx = n - 1;
        if (idx < logRange[0] || idx >= logRange[1]) continue;
        const line = L[idx];
        if (line === undefined) continue;
        scanned++;
        if (isExcluded(line)) continue;
        for (const d of declarationPairs(line)) {
          declarations++;
          const evidence = evidenceIds.get(d.id) || [];
          if (evidence.some((ln) => ln !== n)) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail: f + ':' + n + ' ' + d.verb + ' ' + d.id + ' has no ledger row in the 12.4 / appendix-C regions',
          });
        }
      }
    }
    return emit(
      ctx,
      out,
      'G6-7',
      'declaration/ledger reconciliation',
      findings,
      'every declared ID on an added section 0.3 line exists in the 12.4 / appendix-C ledger (scanned ' +
        scanned +
        ', declarations ' +
        declarations +
        ')'
    );
  },
};

const G6_6 = {
  id: 'G6-6',
  label: 'ledger status closed set',
  run(ctx, out) {
    const findings = [];
    let scanned = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const region = ledgerRegionRange(L);
      if (!region) continue;
      const rows = [];
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const idx = n - 1;
        if (idx < region[0] || idx >= region[1]) continue;
        const line = L[idx];
        if (line === undefined || !LEDGER_ROW_RE.test(line)) continue;
        rows.push({ n: n, line: line });
      }
      // OB-47: the language declaration is a PREMISE only where a word table is actually
      // needed. A §12.4-shaped region with no ledger row is not a directive ledger - another
      // document may legitimately own a §12.4 heading (the RFC does) - so it is skipped and
      // is never forced to declare. Once ledger rows exist the declaration is mandatory: the
      // vocabulary follows the language the document DECLARES (cross-checked against the
      // name), never the bare file-name suffix.
      if (!rows.length) continue;
      const words = LEDGER_STATUS_SETS[declaredLang(L, f)];
      for (const r of rows) {
        scanned++;
        if (isExcluded(r.line)) continue;
        const cells = rowCellsOf(r.line);
        const cell = cells.length ? cells[cells.length - 1] : '';
        if (firstStatusWord(cell, words)) continue;
        findings.push({
          file: f,
          line: r.n,
          text: r.line,
          detail: f + ':' + r.n + ' first status word not in the closed set :: ' + String(cell).trim().slice(0, 60),
        });
      }
    }
    return emit(
      ctx,
      out,
      'G6-6',
      'ledger status closed set',
      findings,
      'every added ledger row opens with a closed-set status word (scanned ' + scanned + ')'
    );
  },
};

/* ---------------------------------------------------------------------------
 * G6-8 (DR number uniqueness and membership) - the mechanical answer to the "DR numbers
 * are minted by hand" debt: the P22 report had to rename its defect id three times
 * (`DR-1` collided with a slice name, `DR-6` with the GATES-EXT ledger, `DR-8` with the
 * assembly-min flake candidate) before it landed.
 *
 * TOKEN RULE (the single source for every face of this criterion):
 *     /(?<![A-Za-z0-9-])DR-(\d{1,2})(?![0-9])/g
 *   - the LEFT boundary is what keeps `ADR-027` from being carved into a phantom
 *     `DR-027` (the repository carries hundreds of `ADR-NNN` substrings);
 *   - the RIGHT boundary is what keeps three-digit ids out of the token set: `DR-027`
 *     yields NO token at all (`\d{1,2}` cannot consume `027` and the lookahead rejects
 *     the trailing digit), so a three-digit number inside a REFERENCE cell is silently
 *     not judged. The mint face catches it explicitly instead (FACE A below): a
 *     three-digit number belongs to the ADR domain, never to the DR defect table.
 *
 * Only STRUCTURED declaration positions are judged, whole line at a time (a parenthetical
 * or any free text is never judged):
 *   FACE A  mint position, UNIQUENESS enforced - `docs/t027/REMAINING_SLICES{,_EN}.md`,
 *           row shape `| **DR-N** |`. A changed mint row whose id already exists on
 *           another line of the SAME file FAILs, and the detail names the id plus BOTH
 *           line numbers. A three-digit id in the id cell FAILs as `三位 = ADR 域`.
 *   FACE B  reference position, MEMBERSHIP - same two files, row shape `| DR-N` (the
 *           occurrence table, including the `DR-3（第 1 次）` / `DR-3 (1st)` variants).
 *   FACE C  reference position, MEMBERSHIP - `docs/validation/*.md` (evidence tree
 *           excluded; `ctx.changedMarkdown()` already applies the OB-79 exclusion every
 *           other G6 sub-check uses), row shape `| DR-N`. A cell carrying several ids
 *           (`DR-3 / DR-4`) or a range (`DR-2 至 DR-5`) is tokenised and EACH token is
 *           judged on its own; CONTIGUITY IS NOT JUDGED.
 *   FACE D  reference position, MEMBERSHIP - same validation files, bullet shape
 *           `- **DR-N（` / `- **DR-N(` / `- **DR-N (`, i.e. anchored on a half- or
 *           full-width opening bracket, optional whitespace allowed between the number
 *           and the bracket (F4 M3: an EN report writes `- **DR-10 (CI first run red)**`,
 *           and a strict adjacency rule silently missed it). Without that anchor a prose
 *           bullet such as `- **DR-2 位断言强度…**` would be judged as if it were a
 *           registration row.
 *
 * PRIMARY KEY SET (the "already registered" side) = every mint id of the CURRENT FULL
 * defect table of either mirror ∪ DR_SEED_IDS. Full file, never just the changed lines:
 * the rows being referenced are normally older than the change set. The two mirrors are
 * judged INDEPENDENTLY (both sides must pass) but share one key set, because a number
 * minted on either side of the pair is registered; mirror symmetry itself is G3(a)'s job.
 *
 * DUAL MINT (v1.46 GATES-CLEANUP M2): G3(a) only judges diff-line-count symmetry of the
 * two mirrors, never id agreement - so it cannot carry the claim that a number is minted
 * on both sides. A changed mint row must therefore ALSO appear as a mint row in the OTHER
 * mirror (read from the same scope-aware source); a one-sided mint FAILs and the detail
 * names the mirror that lacks it. Measured before landing: the two defect tables carry the
 * SAME id set ({2,3,4,5,6,7,9,11}), so the dual-mint rule is satisfiable by the corpus.
 * When the OTHER mirror is unreadable the face cannot be judged: it is SKIPPED (the verdict
 * is untouched) and the record's detail carries the explicit mark
 * `dual-mint: other mirror unreadable, face skipped` - a stated fail-open, never a silent
 * one (v1.46 round-6 C-8).
 *
 * Judged lines: the ADDED lines of the change set only (the same scope as G6-1 … G6-7).
 * With no target file, or no judgeable row, the sub-check reports the same "no premise"
 * record G6-7 reports (`emit` over an empty finding set) - it never invents a status.
 *
 * CAPABILITY BOUNDARY (stated, not hidden): this criterion catches "two mint rows with
 * the same number" and "a reference naming a number nobody minted". It can NOT catch a
 * reference that points at an existing number while describing a DIFFERENT defect (a
 * semantic mismatch), and it can NOT catch a collision inside free text.
 * ------------------------------------------------------------------------ */
const DR_TOKEN_RE = /(?<![A-Za-z0-9-])DR-(\d{1,2})(?![0-9])/g;
const DR_TARGET_FILES = ['docs/t027/REMAINING_SLICES.md', 'docs/t027/REMAINING_SLICES_EN.md'];
const DR_MINT_RE = /^\|\s*\*\*DR-(\d{1,2})\*\*\s*\|/;
/** Three-or-more digits in the mint cell: the ADR domain, never a DR defect id. */
const DR_MINT_3D_RE = /^\|\s*\*\*DR-(\d{3,})\*\*\s*\|/;
const DR_REF_ROW_RE = /^\|\s*DR-(\d{1,2})/;
const DR_REF_BULLET_RE = /^-\s*\*\*DR-(\d{1,2})\s*[（(]/;
const DR_VALIDATION_RE = /^docs\/validation\/[^/]+\.md$/;
/**
 * Hand-registered ids that are legitimately referenced but have no mint row of their own.
 * Kept as an explicit constant (never inferred at run time) so a future deletion of one of
 * these entries is a visible diff, not a silent widening of the key set.
 */
const DR_SEED_IDS = [
  1, // slice name `DR-1` (DEV_PLAN B1); referenced as a slice, never a defect-table row
  8, // the `DR-8` token on the entry line of docs/validation/directive-gates-debt-sweep.md:78
  //    (that line matches FACE D, so EVERY token on it is judged once the line changes;
  //    DR-8 has no mint row of its own). NOT an evidence-tree flake candidate: the
  //    evidence tree is excluded from G6 by OB-79 and is never judged.
  10, // registered as a bullet in docs/validation/directive-gates-debt-sweep.md
];

/** Every DR id token of one line, as numbers (the token rule above, nothing else). */
function drTokensOf(line) {
  DR_TOKEN_RE.lastIndex = 0;
  const out = [];
  let m;
  while ((m = DR_TOKEN_RE.exec(String(line))) !== null) out.push(Number(m[1]));
  return out;
}

/** Mint id -> every line number (1-based) in the given line array. */
function drMintIndex(L) {
  const out = new Map();
  const lines = L || [];
  for (let i = 0; i < lines.length; i++) {
    const m = DR_MINT_RE.exec(String(lines[i] === undefined || lines[i] === null ? '' : lines[i]));
    if (!m) continue;
    const id = Number(m[1]);
    if (!out.has(id)) out.set(id, []);
    out.get(id).push(i + 1);
  }
  return out;
}

/**
 * The "already registered" set: seeds plus every mint id of both CURRENT full defect
 * tables. Read through `readTextFor` so a range/ci run judges against the endpoint
 * revision (the same scope-aware reader G6-5 / G6-7 use).
 */
function drRegisteredIds(ctx) {
  const out = new Set(DR_SEED_IDS);
  for (const f of DR_TARGET_FILES) {
    const t = readTextFor(ctx, f);
    if (t === null) continue;
    for (const id of drMintIndex(t.split('\n')).keys()) out.add(id);
  }
  return out;
}

const G6_8 = {
  id: 'G6-8',
  label: 'DR number uniqueness and membership',
  run(ctx, out) {
    const findings = [];
    const registered = drRegisteredIds(ctx);
    // Uniqueness evidence comes from the FULL file (the collision partner is usually an
    // older row that the change set never touched). `mirrorReadable` records which side is
    // actually readable, because the dual-mint rule below must not judge an absent premise.
    // An unreadable side SKIPS the face (verdict unchanged) but is no longer silent: the
    // PASS detail carries `dual-mint: other mirror unreadable, face skipped` (v1.46 C-8).
    const mintIndex = new Map();
    const mirrorReadable = new Map();
    let dualSkipped = 0;
    for (const f of DR_TARGET_FILES) {
      const t = readTextFor(ctx, f);
      const L = t === null ? ctx.readLineArray(f) : t.split('\n');
      mirrorReadable.set(f, !!L);
      mintIndex.set(f, drMintIndex(L || []));
    }
    let scanned = 0;
    let refs = 0;
    for (const f of ctx.changedMarkdown()) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const mintTarget = DR_TARGET_FILES.indexOf(f) >= 0;
      const refTarget = DR_VALIDATION_RE.test(f);
      if (!mintTarget && !refTarget) continue;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined) continue;
        scanned++;
        if (isExcluded(line)) continue;
        if (mintTarget) {
          const m3 = DR_MINT_3D_RE.exec(line);
          if (m3) {
            findings.push({
              file: f,
              line: n,
              text: line,
              detail: f + ':' + n + ' DR-' + m3[1] + ' is a three-digit id in the mint cell (三位 = ADR 域)',
            });
            continue;
          }
          if (DR_MINT_RE.test(line)) {
            // A mint row IS the registration: membership is vacuous, uniqueness is not.
            const id = Number(DR_MINT_RE.exec(line)[1]);
            const idx = mintIndex.get(f) || new Map();
            const others = (idx.get(id) || []).filter((ln) => ln !== n);
            if (others.length) {
              findings.push({
                file: f,
                line: n,
                text: line,
                detail:
                  f + ':' + n + ' DR-' + id + ' minted twice: line ' + n + ' collides with line ' + others[0],
              });
              continue;
            }
            // DUAL MINT (M2): the same number must be minted in the OTHER mirror too.
            const other = f === DR_TARGET_FILES[0] ? DR_TARGET_FILES[1] : DR_TARGET_FILES[0];
            if (mirrorReadable.get(other)) {
              const oIdx = mintIndex.get(other) || new Map();
              if (!oIdx.has(id)) {
                findings.push({
                  file: f,
                  line: n,
                  text: line,
                  detail:
                    f +
                    ':' +
                    n +
                    ' DR-' +
                    id +
                    ' is minted only in ' +
                    f +
                    '; the mirror ' +
                    other +
                    ' has no mint row for it (dual mint required)',
                });
              }
            } else {
              // Fail-open, stated instead of hidden (v1.46 C-8): the face cannot be judged,
              // so the record says so and the verdict stays whatever the other faces say.
              dualSkipped++;
            }
            continue;
          }
          if (!DR_REF_ROW_RE.test(line)) continue; // any other row shape is not judged
        } else if (!DR_REF_ROW_RE.test(line) && !DR_REF_BULLET_RE.test(line)) {
          continue;
        }
        for (const id of drTokensOf(line)) {
          refs++;
          if (registered.has(id)) continue;
          findings.push({
            file: f,
            line: n,
            text: line,
            detail:
              f + ':' + n + ' DR-' + id + ' is referenced but not registered in the defect table (+ seeds)',
          });
        }
      }
    }
    return emit(
      ctx,
      out,
      'G6-8',
      'DR number uniqueness and membership',
      findings,
      'registered DR set {' +
        [...registered].sort((a, b) => a - b).join(',') +
        '}; no DR collision or dangling reference on added lines (scanned ' +
        scanned +
        ', refs ' +
        refs +
        ')' +
        (dualSkipped
          ? '; dual-mint: other mirror unreadable, face skipped (' + dualSkipped + ' changed mint row(s))'
          : '')
    );
  },
};

module.exports = {
  id: 'G6',
  label: 'ledger metadata consistency',
  title: '台账元数据一致性',
  scripted: true,
  subchecks: [G6_1, G6_2, G6_3, G6_4, G6_5, G6_6, G6_7, G6_8],
  internals: {
    isExcluded,
    headingIndices,
    blockBounds,
    qualify,
    costClause,
    vecText,
    keysOfLine,
    PENDING_RE,
    DONE_RE,
    // G6-5 (OB-32): label collection, target-region detection, resolution rule.
    directiveLabels,
    labelsFromLines,
    targetRegionRanges,
    resolvesRef,
    readTextFor,
    SECTION_REF_RE,
    DIRECTIVE_FILES,
    // G6-6 (OB-38): bilingual closed set, ledger region, row-cell split, first-word rule.
    LEDGER_STATUS_SETS,
    ledgerRegionRange,
    rowCellsOf,
    firstStatusWord,
    LEDGER_ROW_RE,
    // G6-6 (OB-47): header language declaration, closed tag set, name cross-check.
    declaredLang,
    LANG_TAG_OCCURRENCE_RE,
    LANG_TAG_CONT_RE,
    LANG_TAG_EXACT,
    LANG_HEADER_WINDOW,
    // G6-7 (OB-32): declaration verbs / ID keys / window, plus the two index helpers
    // the mutation probes compile against.
    DECL_VERB_RE,
    DECL_ID_RE,
    DECL_WINDOW,
    declarationPairs,
    ledgerIds,
    // G6-8: the token rule, the four structured-position shapes, the seed ids and the
    // two index helpers the migration probes compile against.
    DR_TOKEN_RE,
    DR_TARGET_FILES,
    DR_MINT_RE,
    DR_MINT_3D_RE,
    DR_REF_ROW_RE,
    DR_REF_BULLET_RE,
    DR_VALIDATION_RE,
    DR_SEED_IDS,
    drTokensOf,
    drMintIndex,
    drRegisteredIds,
  },
  run(ctx, out) {
    for (const s of this.subchecks) s.run(ctx, out);
  },
};

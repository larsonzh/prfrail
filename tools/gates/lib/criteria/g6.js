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
 */
const RUN_PATTERNS = [
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s+0+(?![0-9])(?!\s*[A-Za-z\u4e00-\u9fff])/g, why: 'all-zero run id' },
  { re: /(?<![A-Za-z\u4e00-\u9fff])run\s*<[^>]*>/g, why: 'angle-bracket placeholder' },
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
    for (let j = i + 1; j < L.length; j++) {
      if (REGION_HEAD_C1.test(L[j])) {
        end = j;
        break;
      }
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
      const words = /_EN\.md$/i.test(f) ? LEDGER_STATUS_SETS.en : LEDGER_STATUS_SETS.cn;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const idx = n - 1;
        if (idx < region[0] || idx >= region[1]) continue;
        const line = L[idx];
        if (line === undefined || !LEDGER_ROW_RE.test(line)) continue;
        scanned++;
        if (isExcluded(line)) continue;
        const cells = rowCellsOf(line);
        const cell = cells.length ? cells[cells.length - 1] : '';
        if (firstStatusWord(cell, words)) continue;
        findings.push({
          file: f,
          line: n,
          text: line,
          detail: f + ':' + n + ' first status word not in the closed set :: ' + String(cell).trim().slice(0, 60),
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

module.exports = {
  id: 'G6',
  label: 'ledger metadata consistency',
  title: '台账元数据一致性',
  scripted: true,
  subchecks: [G6_1, G6_2, G6_3, G6_4, G6_5, G6_6],
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
  },
  run(ctx, out) {
    for (const s of this.subchecks) s.run(ctx, out);
  },
};

'use strict';

/**
 * G7 review-bundle completeness (architecture design section 3) - structural only,
 * no content semantics, no network.
 *
 * Target files: changed `docs/validation/*.md` (`docs/validation/evidence/**` excluded).
 * The rule text lives in DELIVERY_DIRECTIVE.md and is therefore outside the G7 scan
 * domain: the self-reference trap is avoided structurally.
 *
 * Machine contract: a changed validation report MUST contain a fenced block whose
 * info string is exactly `artifacts`; each row is `<repo-relative path><TAB><disposition>`
 * with disposition in { `repo` (may be omitted), `local-only`, `missing-historical` }.
 *
 *   G7-a the artifacts block is missing
 *   G7-b path/disposition row is unresolvable or carries an unknown token
 *   G7-c a changed finding row lacks a disposition CELL
 *   G7-d a changed status line says "stopped" without the "stopped - incomplete" mark
 *   G7-e a changed line carries a verdict line outside the DELIVERY_DIRECTIVE 4.1 closed set
 *   G7-f a changed `sixstage` transcript fence breaks the ⑥ output contract
 *        (A-E sections, closed-set conclusion last, SEVERITY rows ⇔ FINDINGS)
 *
 * G7-c is STRUCTURAL (GATES-EXT F2): the disposition token must sit in the row's LAST
 * CELL, not merely somewhere on the line. A row whose description quotes a disposition
 * word (`| **High** | 台账称“已修”但处置列为空 | |`) has an EMPTY disposition cell and
 * must FAIL - the token is never taken out of a description cell.
 *
 * G7-b `repo` disposition (GATES-EXT F2): the token means the SAME thing in every
 * scope. range/ci: the path must exist at the scope's endpoint revision
 * (`git cat-file -e`); tree/index: the path must be STAGED in the index
 * (`git ls-files --error-unmatch`) AND exist on disk. A worktree-only / gitignored
 * path is thus not enough locally: a tree-scope run is a faithful CI pre-flight
 * only after `git add` of the declared artifact.
 *
 * G7-d is deliberately restricted to STATUS LINES (`| 状态 | ... |`, `状态：...`, `Status: ...`):
 * a changed prose line that merely names the rule (measured: docs/validation/sw5-work-tiers.md:88
 * reads "...（停止态报告头标记）= OB-10 所指缺口") is not a status declaration, and matching it
 * would be the same self-reference trap G6 guards against with its exclusion markers.
 *
 * Explicitly NOT judged: artifact content, run reachability, number truth, report
 * conclusion truth.
 *
 * The evidence-tree exclusion prefix is imported from `../ctx` (section 6.3 "判据作用域排除",
 * OB-79) rather than redeclared here: G7's exclusion and the G1-a / G2 / G4a / G4b / G6
 * exclusion must stay the SAME single-source definition.
 */

const { EVIDENCE_TREE_PREFIX } = require('../ctx');

const TARGET_RE = /^docs\/validation\/.*\.md$/;
const ARTIFACT_FENCE_RE = /^```artifacts[ \t]*$/;
const CLOSE_FENCE_RE = /^```[ \t]*$/;
const DISPOSITIONS = ['repo', 'local-only', 'missing-historical'];
const FINDING_ROW_RE = /^\|\s*\*\*(Critical|High|Medium|Low)\*\*\s*\|/;
/**
 * Accepted disposition tokens (GATES-EXT F3): ONE set covering both languages.
 * The check stays STRUCTURAL (a disposition cell must exist); the token is never
 * interpreted. G7 targets `docs/validation/*.md` including the `_EN` mirrors, so a
 * Chinese-only vocabulary would fail the English reports by construction.
 */
const DISPOSITION_TOKENS = [
  '已修',
  '已整改',
  '已闭合',
  '已处置',
  '接受',
  '拒收',
  '记录在案',
  '待办',
  '观察中',
  '待议',
  '已停止',
  'Fixed',
  'Resolved',
  'Closed',
  'Accepted',
  'Rejected',
  'Deferred',
  'Recorded',
  'To do',
  'Observed',
  'To be discussed',
  'Not a problem',
  'Documented',
];
const DISPOSITION_CELL_RE = new RegExp(DISPOSITION_TOKENS.join('|'));

/**
 * Cells of a markdown table row, in order, trimmed.
 * The outer pipe artifacts are dropped (`|` at both ends), but a TRAILING EMPTY CELL
 * is preserved: `| a | b | |` yields ['a', 'b', ''] - that is exactly the
 * "disposition cell left blank" shape G7-c must FAIL (GATES-EXT F2). Dropping the
 * empty cell would let a disposition word inside the description cell pass the row.
 */
function rowCells(line) {
  const parts = String(line).split('|');
  const cells = String(line).trim().endsWith('|') ? parts.slice(1, -1) : parts.slice(1);
  return cells.map((c) => c.trim());
}

/**
 * The disposition CELL of a finding row: the row's LAST cell, structurally - never a
 * substring search over the whole line (GATES-EXT F2).
 */
function dispositionCell(line) {
  const cells = rowCells(line);
  return cells.length ? cells[cells.length - 1] : '';
}

/**
 * Stopped-state tokens (GATES-EXT F3): ONE bilingual set. G7 targets
 * `docs/validation/*.md` including the `_EN` mirrors, so a Chinese-only vocabulary
 * let `| Status | Stopped (pending user ruling) |` through - and G7-d claims to cover
 * the English reports by construction. Matching is case-insensitive for the ASCII
 * tokens (`stopped` / `HALTED` are the same status as `Stopped` / `Halted`).
 */
const STOP_TOKENS = ['已停止', '停止态', 'Stopped', 'Halted'];
const STOP_RE = new RegExp(STOP_TOKENS.join('|'), 'i');
/**
 * The "incomplete" mark the stopped state must carry: the CN canonical form plus its
 * EN mirror. `已停止——未完成` alone would require an English report to write Chinese.
 * Compared case-insensitively, like the tokens above.
 */
const STOP_OK_CN = '已停止——未完成';
const STOP_MARKS = [STOP_OK_CN, 'Stopped — incomplete'];
const STOP_OK = STOP_OK_CN;

/**
 * G7-e verdict-line closed set (OB-42, v1.33). DELIVERY_DIRECTIVE 4.1 item 1 defines the
 * ONLY legal verdict lines as `<prefix>: <value>` over a closed set of 11 prefixes, and
 * abolishes bare `PASS` / `FINDINGS` plus the residual `RE_` placeholders. Until OB-42
 * nothing mechanical checked that membership (the item was booked under G7-c, which only
 * judges a finding row's disposition cell).
 *
 * Judged: ADDED lines of a changed report (the same scope as G7-c / G7-d).
 *   (1) a line whose leading prefix is one of the 11 closed-set prefixes MUST carry a value
 *       from that prefix's enumerated set; an out-of-set value is a format deviation;
 *   (2) a line that is EXACTLY a bare closed-set value (modulo markdown noise) is a format
 *       deviation - the "bare token" rule;
 *   (3) `<UPPERCASE PREFIX>: <exactly one closed-set value>` whose prefix is NOT in the
 *       closed set is the "table-external form" deviation;
 *   (4) a residual `RE_` placeholder is a deviation (4.1: abolished, must not remain).
 *
 * Deliberate boundary (stated in 6.3): the value must be RECOGNISABLE, so a coarse prefix
 * such as `SELFTEST: PASS 362/362` (a gate reading, not a verdict line) is NOT judged. The
 * EXISTENCE and POSITION of the verdict line (4.1 says "last line") are not mechanically
 * enforced either - both residuals are carried by the 5 / 6 / 7 review layers.
 *
 * Self-reference guard: the same five exclusion markers as G6, so a report that QUOTES an
 * illegal form on a marked line (示例 / 引文 / 不得残留 / 冲突 / 矛盾) is not judged.
 */
const VERDICT_PREFIXES = {
  DESIGN: ['DONE', 'NEEDS_CLARIFICATION'],
  IMPLEMENT: ['DONE', 'BLOCKED'],
  TEST: ['DONE', 'BLOCKED'],
  'PRE-REVIEW': ['PASS', 'PASS WITH FIXES', 'FINDINGS'],
  'INDEPENDENT SCAN': ['PASS', 'FINDINGS'],
  'FINAL REVIEW': ['PASS', 'PASS WITH FIXES', 'PARTIAL', 'FINDINGS'],
  'RE-REVIEW': ['PASS', 'PASS WITH FIXES', 'PARTIAL', 'FINDINGS'],
  'BLIND REVIEW': ['PASS', 'PASS WITH FIXES', 'PARTIAL', 'FINDINGS'],
  DOCS: ['DONE', 'BLOCKED'],
  SLICE: ['STARTED', 'READY_FOR_REVIEW', 'BLOCKED'],
  'DIRECTIVE REVIEW': ['PASS', 'FINDINGS'],
};
/** Longest first, so `PASS WITH FIXES` wins over `PASS`. */
const VERDICT_VALUES = [...new Set([].concat(...Object.values(VERDICT_PREFIXES)))].sort((a, b) => b.length - a.length);
const VERDICT_PREFIX_RE = /^(?:[-*>]\s*)*`{0,2}\*{0,2}([A-Z][A-Z _-]*?)\*{0,2}`{0,2}\s*[:：]\s*(.*)$/;
const RE_PLACEHOLDER_RE = /\bRE_[A-Z]/;
const G7E_EXCLUDE_MARKERS = ['不得残留', '示例', '引文', '冲突', '矛盾'];

/**
 * The longest closed-set value the rest of a verdict line starts with, or null. A value
 * followed by an alphanumeric character is NOT a verdict value (`PASS 362/362` is a gate
 * reading); a value followed by punctuation or markdown (`` `PASS` ``, `PASS（无 Medium+）`,
 * `PASS — 无 Medium+`) is.
 */
function verdictValue(rest) {
  let s = String(rest).trim();
  while (s.startsWith('*') || s.startsWith('`')) s = s.slice(1).trim();
  for (const v of VERDICT_VALUES) {
    if (!s.startsWith(v)) continue;
    const after = s.slice(v.length).trim();
    if (after === '' || !/^[A-Za-z0-9]/.test(after)) return v;
  }
  return null;
}

/** The deviation on a verdict-shaped / bare-token line, or null when the line is legal. */
function verdictFinding(line) {
  const m = VERDICT_PREFIX_RE.exec(line);
  if (m) {
    const prefix = m[1].trim();
    const rest = m[2];
    if (Object.prototype.hasOwnProperty.call(VERDICT_PREFIXES, prefix)) {
      const v = verdictValue(rest);
      if (v && VERDICT_PREFIXES[prefix].indexOf(v) >= 0) return null;
      return 'prefix `' + prefix + '` with ' + (v ? '`' + v + '`' : 'a value') + ' outside the closed set';
    }
    const bare = verdictValue(rest);
    if (bare && String(rest).replace(/[*`\s]/g, '') === bare) {
      return 'verdict line with table-external prefix `' + prefix + '`';
    }
    return null;
  }
  const t = String(line)
    .replace(/^[-*>]\s*/, '')
    .replace(/[*`\s]/g, '')
    .replace(/[。．.,；;]+$/, '');
  return VERDICT_VALUES.indexOf(t) >= 0 ? 'bare verdict token `' + t + '` without a prefix' : null;
}

/** Does the line carry the incomplete mark (either language, any case)? */
function hasStopMark(line) {
  const s = String(line).toLowerCase();
  return STOP_MARKS.some((m) => s.indexOf(m.toLowerCase()) >= 0);
}
/** Report status line: a `状态` / `Status` label followed by a separator (table row included). */
const STATUS_LINE_RE = /^(?:\|\s*)?(?:\*\*)?(?:状态|Status)(?:\*\*)?\s*[：:|]/i;

function targets(ctx) {
  return ctx.changedPaths.filter((p) => TARGET_RE.test(p) && !p.startsWith(EVIDENCE_TREE_PREFIX));
}

/** Fenced `artifacts` block, or null when the report has none. */
function artifactBlock(L) {
  const start = L.findIndex((l) => ARTIFACT_FENCE_RE.test(l));
  if (start < 0) return null;
  let end = L.findIndex((l, i) => i > start && CLOSE_FENCE_RE.test(l));
  if (end < 0) end = L.length;
  return { start, end, body: L.slice(start + 1, end) };
}

const G7_A = {
  id: 'G7-a',
  label: 'artifacts block present',
  run(ctx, out) {
    const tgts = targets(ctx);
    const findings = [];
    let checked = 0;
    for (const f of tgts) {
      const L = ctx.readLineArray(f);
      if (!L) continue; // deleted / unreadable report: nothing to judge
      checked++;
      if (!artifactBlock(L)) {
        findings.push({ file: f, line: 0, text: 'no ```artifacts block', detail: f + ' missing ```artifacts block' });
      }
    }
    if (!tgts.length) {
      out.rec('G7-a', 'INFO', 'artifacts block present', 'no changed report in scope');
      return { findings };
    }
    out.rec(
      'G7-a',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'artifacts block present',
      findings.length
        ? findings.map((f) => f.detail).join(' | ')
        : checked + ' report(s) OK' + (checked < tgts.length ? ' (' + (tgts.length - checked) + ' unreadable/deleted skipped)' : '')
    );
    return { findings };
  },
};

const G7_B = {
  id: 'G7-b',
  label: 'artifact path/disposition resolvable',
  run(ctx, out) {
    const tgts = targets(ctx);
    const findings = [];
    let rows = 0;
    for (const f of tgts) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const block = artifactBlock(L);
      if (!block) continue;
      block.body.forEach((raw, k) => {
        const lineNo = block.start + 2 + k;
        const t = raw.trim();
        if (!t || t.startsWith('#')) return;
        rows++;
        const parts = raw.split('\t');
        const p = parts[0].trim().replace(/^\.\//, '');
        const disp = parts.length > 1 ? parts[1].trim() : '';
        if (!p) {
          findings.push({ file: f, line: lineNo, text: raw, detail: f + ':' + lineNo + ' empty artifact path' });
          return;
        }
        if (disp !== '' && DISPOSITIONS.indexOf(disp) < 0) {
          findings.push({
            file: f,
            line: lineNo,
            text: raw,
            detail: f + ':' + lineNo + ' unknown disposition `' + disp + '` for ' + p,
          });
          return;
        }
        const kind = disp === '' ? 'repo' : disp;
        if (kind !== 'repo') return; // local-only / missing-historical are accepted as declared
        const res = ctx.repoArtifactResolvable
          ? ctx.repoArtifactResolvable(p)
          : { ok: ctx.exists(p), source: 'worktree-existence', reason: '' };
        if (!res.ok) {
          findings.push({
            file: f,
            line: lineNo,
            text: raw,
            detail:
              f +
              ':' +
              lineNo +
              ' declared repo artifact not resolvable (' +
              res.source +
              '): ' +
              p +
              (res.reason ? ' - ' + res.reason : ''),
          });
        }
      });
    }
    if (!tgts.length) {
      out.rec('G7-b', 'INFO', 'artifact path/disposition resolvable', 'no changed report in scope');
      return { findings };
    }
    out.rec(
      'G7-b',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'artifact path/disposition resolvable',
      findings.length
        ? findings.map((f) => f.detail).join(' | ')
        : rows +
            ' artifact row(s) OK (repo check: ' +
            (ctx.kind === 'range' || ctx.kind === 'ci'
              ? 'git cat-file -e <endpoint>:<path>'
              : 'git ls-files --error-unmatch + worktree presence') +
            ')'
    );
    return { findings };
  },
};

const G7_C = {
  id: 'G7-c',
  label: 'finding row disposition cell',
  run(ctx, out) {
    const tgts = targets(ctx);
    const findings = [];
    let rows = 0;
    for (const f of tgts) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined || !FINDING_ROW_RE.test(line)) continue;
        rows++;
        const cell = dispositionCell(line);
        if (cell && DISPOSITION_CELL_RE.test(cell)) continue;
        findings.push({
          file: f,
          line: n,
          text: line,
          detail:
            f +
            ':' +
            n +
            ' finding row without disposition' +
            (cell ? '' : ' (last cell empty)'),
        });
      }
    }
    if (!tgts.length) {
      out.rec('G7-c', 'INFO', 'finding row disposition cell', 'no changed report in scope');
      return { findings };
    }
    out.rec(
      'G7-c',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'finding row disposition cell',
      findings.length
        ? findings.map((f) => f.detail).join(' | ')
        : rows + ' changed finding row(s) with disposition'
    );
    return { findings };
  },
};

const G7_D = {
  id: 'G7-d',
  label: 'stopped state marked incomplete',
  run(ctx, out) {
    const tgts = targets(ctx);
    const findings = [];
    for (const f of tgts) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      const headings = [];
      for (let i = 0; i < L.length; i++) if (/^#{2,4} /.test(L[i])) headings.push(i);
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined || !STOP_RE.test(line)) continue;
        if (hasStopMark(line)) continue;
        if (!STATUS_LINE_RE.test(line)) continue; // prose naming the rule is not a status line
        let start = 0;
        let end = L.length;
        for (const h of headings) {
          if (h <= n - 1) start = h;
          else {
            end = h;
            break;
          }
        }
        let ok = false;
        for (let i = start; i < end; i++) if (hasStopMark(L[i] || '')) ok = true;
        if (ok) continue;
        findings.push({
          file: f,
          line: n,
          text: line,
          detail: f + ':' + n + ' stopped state without ' + STOP_MARKS.map((m) => '`' + m + '`').join(' / '),
        });
      }
    }
    if (!tgts.length) {
      out.rec('G7-d', 'INFO', 'stopped state marked incomplete', 'no changed report in scope');
      return { findings };
    }
    out.rec(
      'G7-d',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'stopped state marked incomplete',
      findings.length ? findings.map((f) => f.detail).join(' | ') : 'no stopped-state row missing the incomplete mark'
    );
    return { findings };
  },
};

const G7_E = {
  id: 'G7-e',
  label: 'verdict line closed set',
  run(ctx, out) {
    const tgts = targets(ctx);
    const findings = [];
    let scanned = 0;
    for (const f of tgts) {
      const L = ctx.readLineArray(f);
      if (!L) continue;
      for (const n of [...ctx.addedLineNumbers(f)].sort((a, b) => a - b)) {
        const line = L[n - 1];
        if (line === undefined) continue;
        if (G7E_EXCLUDE_MARKERS.some((k) => line.includes(k))) continue;
        scanned++;
        if (RE_PLACEHOLDER_RE.test(line)) {
          findings.push({ file: f, line: n, text: line, detail: f + ':' + n + ' residual `RE_` placeholder' });
          continue;
        }
        const why = verdictFinding(line);
        if (why) {
          findings.push({ file: f, line: n, text: line, detail: f + ':' + n + ' ' + why + ' :: ' + line.slice(0, 120) });
        }
      }
    }
    if (!tgts.length) {
      out.rec('G7-e', 'INFO', 'verdict line closed set', 'no changed report in scope');
      return { findings };
    }
    out.rec(
      'G7-e',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'verdict line closed set',
      findings.length
        ? findings.map((f) => f.detail).join(' | ')
        : 'every verdict line on added report lines is inside the 4.1 closed set (scanned ' + scanned + ')'
    );
    return { findings };
  },
};

/* ---------------------------------------------------------------------------
 * G7-f (⑥ output contract: the `sixstage` transcript fence).
 *
 * The two-slice debt this closes: the ⑥ independent-scan output has now missed the
 * A-E transcript shape for two slices running, because nothing mechanical ever looked
 * at it. DELIVERY_DIRECTIVE 4.1 item 1 already owns the LEGAL VERDICT LINE; this
 * sub-check owns the TRANSCRIPT around it.
 *
 * FENCE DEFINITION: a three-backtick fence whose info string is exactly `sixstage`.
 * The body of such a fence is the transcript. A fence is "changed" when ANY of its
 * body lines appears in this change set's added-line set; an untouched fence is not
 * re-judged (same scope discipline as G7-c / G7-d / G7-e).
 *
 * THREE FACES (inside a changed fence):
 *   1. `^#{2,4}\s+Section\s+([A-E])(?:\s|$)` - the five headings must be present, once
 *      each, in the order A→B→C→D→E (a missing, duplicated or reordered heading FAILs).
 *      The heading may carry descriptive text after the letter (`## Section A 扫描范围与方法`),
 *      which is the shape the ⑥ layer actually emits; the `(?:\s|$)` boundary keeps
 *      `## Section AA` and `## Subsection A` out (v1.46).
 *   2. The LAST non-empty body line must be the conclusion line, carrying prefix
 *      `INDEPENDENT SCAN` and a value from `VERDICT_PREFIXES['INDEPENDENT SCAN']`.
 *      This file's closed set is the ONLY source of that vocabulary - it is never
 *      copied here (a copy would drift the moment 4.1 changes). Content after the
 *      conclusion line FAILs: a note belongs OUTSIDE the fence.
 *   3. Section B consistency: a SEVERITY row (a table row with ANY cell that equals
 *      Critical / High / Medium / Low once surrounding blanks and `*` / backticks are
 *      stripped, or - legacy fallback - a line whose first word after the leading `|` /
 *      `*` is one of them and that carries a table pipe) ⇔ the conclusion must be
 *      `FINDINGS`. Severity rows with `PASS`, or `FINDINGS` with none, FAIL.
 *      `UNCOVERED` is deliberately NOT a severity token.
 *
 *      WHY ANY CELL (v1.46): the ⑥ layer writes `| # | SEVERITY | … |`, so the severity
 *      sits in the SECOND cell. The old "first word" test read the row number (`1`) as the
 *      first word and therefore never saw a severity row at all - a real ⑥ output with a
 *      `Low` row and an `INDEPENDENT SCAN: FINDINGS` conclusion was failed as
 *      "FINDINGS with no SEVERITY row". Cell EQUALITY (not substring, not first word) is
 *      what makes the second-column shape detectable without matching prose that merely
 *      mentions a severity word.
 *
 * ABSENCE (read from the CURRENT FULL report body, not from the change set - the
 * question is a property of the document). The branch splits on the change-set class
 * FIRST (v1.46 round-8 F1, the High fix). For a report the change set MODIFIES, the
 * v1.46 FOUR states (GATES-CLEANUP H1) are kept unchanged:
 *   - the report names `INDEPENDENT SCAN` while carrying no `sixstage` fence at all
 *     ⇒ FAIL (self-contradictory: a scan that was never transcribed);
 *   - a fixed-form skip declaration (`- ⑥ 跳过：<reason>` / `- ⑥ skipped: <reason>`) and
 *     no fence ⇒ PASS (the contract was waived explicitly);
 *   - a fixed-form HISTORICAL EXEMPTION (`- ⑥ 历史豁免：<reason>` / `- ⑥ historical
 *     exemption: <reason>`) and no fence ⇒ PASS, reported as the historical-exemption
 *     channel;
 *   - neither ⇒ INFO, never a FAIL: for a report that already existed, the EXISTENCE
 *     obligation is not a mechanical criterion, which is the same boundary G7-e
 *     already states for the verdict line.
 * For a report ADDED by the change set the bare state is closed (F1): such a report
 * must carry the ```sixstage fence or a NON-EMPTY-reason skip declaration (the §7.2.1
 * legitimate skip - the one waiver a new report may use); with neither it FAILs. The
 * HISTORICAL EXEMPTION stays invalid for a new report whatever its reason says, and a
 * marker still FAILs before the skip declaration is considered (the v1.46 precedence
 * is kept: a skip declaration never neutralises a marker).
 *
 * WAIVER PRECONDITIONS (v1.46 round-6 C-1). The waiver channel used to be reachable by
 * writing the fixed form and nothing else: an empty reason was accepted, and a report
 * FIRST ADDED by the change set could waive the fence exactly like a historical one.
 * Both are closed here, mechanically:
 *   (a) REASON NON-EMPTY - both declaration regexes end in `[：:]\s*\S`, so the fixed
 *       form alone is not a declaration. A declaration with a blank reason is NOT a
 *       declaration at all: it falls through to the marker / bare states exactly as if
 *       the line were absent (the verdict is then decided by those states, never by the
 *       empty-reason line).
 *   (b) NEW REPORTS (v1.46 round-8 F1, refining the v1.46 C-1 guard) - when the report
 *       is ADDED by this change set (as opposed to MODIFIED), the HISTORICAL EXEMPTION
 *       is a FAIL whatever its reason says (a report born inside the change set has no
 *       history to exempt), and a report with neither the fence nor a non-empty-reason
 *       skip declaration is a FAIL: a new report must carry the ```sixstage fence or
 *       that declaration. The v1.46 C-1 guard rejected the skip declaration on a new
 *       report too; round-8 F1 re-classified it as the one legitimate waiver, because
 *       §7.2.1 allows a documentation-only slice to skip ⑥ with a recorded reason, and
 *       the brand-new report of such a slice is exactly the shape that guard used to
 *       fail.
 *
 * "NEW" IS DECIDED FROM THE CHANGE-SET STATUS, NOT FROM A BASELINE PROBE. `ctx.recordFor`
 * (the same record `ctx.changed` carries) already exposes git's own status for every
 * changed path, so no extra git call and no shell-out is needed:
 *   - `untracked` (tree scope, porcelain `??`)            => added;
 *   - status letter `A` (index/tree staging, `--name-status` A in range/ci) => added;
 *   - `M` / `R` / `D` / any other status                  => NOT added;
 *   - no record for the path (a synthetic context)        => NOT added, i.e. the
 *     pre-C-1 semantics are kept rather than inventing a stricter verdict from missing
 *     data. The selftest pins both branches with an explicit record.
 *
 * PRECEDENCE (deliberate, stated): the historical exemption is tested BEFORE the
 * marker rule, because a report written before v1.46 carries the marker by definition
 * (measured: ~29 report pairs under docs/validation/ name `INDEPENDENT SCAN` with no
 * fence) - letting the marker rule win would make the channel unreachable and freeze
 * every historical report the moment it is touched. The SKIP declaration keeps its
 * shipped precedence (a marker still FAILs), so no existing behaviour changes. The
 * Since v1.46 round-8 F1 the newly-added branch is entered FIRST and keeps the same
 * internal order (exemption ⇒ marker ⇒ skip ⇒ bare), so the exemption on an added
 * report FAILs with the exemption wording, a marker with the marker wording, and the
 * bare state with the fence-or-declaration obligation.
 *
 * EXCLUSION REUSE: the marker scan is filtered by the same five exclusion markers as
 * G7-e (`G7E_EXCLUDE_MARKERS`), so a report whose only `INDEPENDENT SCAN` occurrence is
 * an illustrative / quoted line is not judged - one constant, never a copy. The SAME
 * five markers filter the Section B SEVERITY scan (face 3), so a quoted severity row
 * never counts as a finding row (v1.46 round-6 C-3: both scans reuse one constant).
 *
 * FENCE CLOSURE: a `sixstage` fence must be CLOSED. An unclosed fence FAILs and the
 * detail names the opening line: an unterminated fence swallows the rest of the
 * document, so its body cannot be a transcript and the three faces below are skipped
 * for that block.
 *
 * Untouched: `VERDICT_PREFIXES` (the single source of the closed set) and G7-e. The
 * five exclusion markers of G7-e are reused read-only for the severity scan, so a
 * report that quotes an illegal form on a marked line stays unjudged.
 * ------------------------------------------------------------------------ */
const SIXSTAGE_FENCE_RE = /^```sixstage[ \t]*$/;
/**
 * A section heading of the transcript (v1.46). The ⑥ layer emits descriptive text after the
 * letter (`## Section A 扫描范围与方法`), so the `$` anchor is replaced by a `(?:\s|$)`
 * boundary: the letter must be followed by a blank or the end of the line, which admits the
 * real shape while still refusing `## Section AA` and `## Subsection A`.
 */
const SIXSTAGE_SECTION_RE = /^#{2,4}\s+Section\s+([A-E])(?:\s|$)/;
const SIXSTAGE_SECTION_SEQ = 'ABCDE';
const SIXSTAGE_PREFIX = 'INDEPENDENT SCAN';
/**
 * The fixed skip-declaration forms the directive allows for a waived ⑥.
 * The trailing `\s*\S` is the C-1 precondition: the fixed form alone is NOT a
 * declaration, the colon must be followed by at least one non-blank character.
 */
const SIXSTAGE_SKIP_RE = /^-\s*⑥\s*(?:跳过|skipped)\s*[：:]\s*\S/;
/**
 * The historical-exemption declaration (v1.46 GATES-CLEANUP H1): the transition channel
 * for reports committed before v1.46 whose verbatim transcription cannot be recovered.
 * A sibling of the skip declaration, with a different reason text - and the two are
 * spelled out together in the directive as `^\-\s*⑥\s*(?:跳过|skipped|历史豁免|historical
 * exemption)\s*[：:]\s*\S`, so a reader can match the code to the clause. The trailing
 * `\s*\S` is the same C-1 precondition the skip form carries: an empty reason is not a
 * declaration.
 */
const SIXSTAGE_HIST_RE = /^-\s*⑥\s*(?:历史豁免|historical exemption)\s*[：:]\s*\S/;
const SIXSTAGE_SEVERITY_RE = /^(?:Critical|High|Medium|Low)\b/;
/** The closed severity token set of section B (v1.46). `UNCOVERED` is deliberately absent. */
const SIXSTAGE_SEVERITY_CELLS = ['Critical', 'High', 'Medium', 'Low'];
const SIXSTAGE_MARKER = 'INDEPENDENT SCAN';

/**
 * Was `f` ADDED by this change set (as opposed to merely MODIFIED)? (v1.46 round-6 C-1)
 *
 * Source of truth = the change-set record the context already carries. No baseline probe
 * and no `git` shell-out: `ctx.recordFor` returns the SAME record `ctx.changed` holds, and
 * that record is built from `git status --porcelain -z` (tree/index) or
 * `git diff --name-status -z` (range/ci), so both carry git's own status letter.
 *
 * `untracked` (`??`) and status letter `A` are the two ways git reports an addition; every
 * other status (`M`, `R`, `D`, `T`, ...) is a modification of a path that already existed.
 * A path with NO record keeps the pre-C-1 semantics (NOT added): a synthetic context that
 * cannot answer must not be read as the stricter verdict.
 */
function addedInChangeSet(ctx, f) {
  let rec = null;
  if (ctx && typeof ctx.recordFor === 'function') rec = ctx.recordFor(f);
  if (!rec && ctx && Array.isArray(ctx.changed)) {
    rec = ctx.changed.filter((r) => r && r.path === f)[0] || null;
  }
  if (!rec) return false;
  if (rec.untracked) return true;
  return String(rec.status || '').charAt(0) === 'A';
}

/**
 * Every `sixstage` fence of a report, as { start, end, body, closed } over the line array.
 * `closed` is false when no closing fence follows within the document (v1.46 G7-f: an
 * unterminated fence is a FAIL, not an implicit "extends to end of file").
 */
function sixstageBlocks(L) {
  const out = [];
  for (let i = 0; i < L.length; i++) {
    if (!SIXSTAGE_FENCE_RE.test(L[i])) continue;
    let end = L.length;
    let closed = false;
    for (let j = i + 1; j < L.length; j++) {
      if (CLOSE_FENCE_RE.test(L[j])) {
        end = j;
        closed = true;
        break;
      }
    }
    out.push({ start: i, end: end, body: L.slice(i + 1, end), closed: closed });
    i = end;
  }
  return out;
}

/** `{prefix, value}` for a verdict-SHAPED line, or null when the line is not one. */
function sixstageConclusion(line) {
  const m = VERDICT_PREFIX_RE.exec(String(line));
  if (!m) return null;
  return { prefix: m[1].trim(), value: verdictValue(m[2]) };
}

/**
 * Does ANY `|`-separated cell of a table row equal a severity token, after stripping
 * surrounding blanks and the emphasis / code markers a report wraps the token in?
 * The comparison is EQUALITY, never a prefix or substring: `| 1 | Low | 无 | … |` matches,
 * while `| Low risk prose | … |` and `| 1 | SEVERITY | … |` (the header row) do not.
 */
function sixstageSeverityCell(s) {
  const cells = s.split('|');
  for (let i = 0; i < cells.length; i++) {
    let c = cells[i].trim();
    while (c.startsWith('*') || c.startsWith('`')) c = c.slice(1).trim();
    while (c.endsWith('*') || c.endsWith('`')) c = c.slice(0, -1).trim();
    if (SIXSTAGE_SEVERITY_CELLS.indexOf(c) >= 0) return true;
  }
  return false;
}

/**
 * Is this a SEVERITY row of section B? Two shapes, unioned (v1.46):
 *   (a) a table row (`|`-leading) whose ANY cell equals a severity token. This is the
 *       shape the ⑥ layer actually emits - `| # | SEVERITY | … |` puts the severity in
 *       the second cell, which the old "first word" test could never see;
 *   (b) the legacy first-word form `| **High** | … |` / `High | …`: strip the leading `|`
 *       and `*` and require the first word to BE a severity token plus a table pipe.
 *       (b) is kept as a fallback so existing fixtures and shipped reports are not weakened.
 * Without (a) the face would miss every real row; without (b) it would regress the shapes it
 * was originally pinned on - which is exactly the fail-open shape this sub-check closes.
 */
function sixstageSeverityRow(line) {
  const s = String(line).trim();
  if (s.startsWith('|') && sixstageSeverityCell(s)) return true;
  let t = s;
  while (t.startsWith('|')) t = t.slice(1).trim();
  while (t.startsWith('*')) t = t.slice(1);
  if (!SIXSTAGE_SEVERITY_RE.test(t)) return false;
  return t.indexOf('|') >= 0;
}

const G7_F = {
  id: 'G7-f',
  label: 'sixstage transcript contract',
  run(ctx, out) {
    const tgts = targets(ctx);
    const findings = [];
    let fenced = 0;
    let judged = 0;
    let skipped = 0;
    let exempt = 0;
    let bare = 0;
    for (const f of tgts) {
      const L = ctx.readLineArray(f);
      if (!L) continue; // deleted / unreadable report: nothing to judge
      const blocks = sixstageBlocks(L);
      if (!blocks.length) {
        // Full-body absence rule (see ABSENCE above). F1 (v1.46 round-8) splits the
        // branch on the change-set class FIRST: a report ADDED by this change set may
        // not be bare (fence or non-empty-reason skip declaration, else FAIL) and may
        // never use the historical exemption, while an EXISTING report keeps the v1.46
        // four states unchanged.
        const hist = L.some((l) => SIXSTAGE_HIST_RE.test(String(l)));
        const skip = L.some((l) => SIXSTAGE_SKIP_RE.test(String(l)));
        const newlyAdded = addedInChangeSet(ctx, f);
        const marker = L.some((l) => {
          const s = String(l);
          if (G7E_EXCLUDE_MARKERS.some((k) => s.indexOf(k) >= 0)) return false;
          return s.indexOf(SIXSTAGE_MARKER) >= 0;
        });
        const pushNoFenceMarker = () =>
          findings.push({
            file: f,
            line: 0,
            text: 'no ```sixstage fence',
            detail: f + ' names `' + SIXSTAGE_MARKER + '` but carries no ```sixstage fence',
          });
        if (newlyAdded) {
          if (hist) {
            findings.push({
              file: f,
              line: 0,
              text: '⑥ historical exemption on a newly added report',
              detail:
                f +
                ':0 new report may not use the ⑥ HISTORICAL EXEMPTION channel (`- ⑥ 历史豁免：` / ' +
                '`- ⑥ historical exemption: `): a report born inside the change set has no history ' +
                'to exempt (新增报告不得使用历史豁免通道，须给出 sixstage 围栏或带理由的跳过声明)',
            });
          } else if (marker) {
            pushNoFenceMarker();
          } else if (skip) {
            skipped++; // the one waiver a new report may carry: the non-empty-reason ⑥ skip
          } else {
            findings.push({
              file: f,
              line: 0,
              text: 'no ```sixstage fence on a newly added report',
              detail:
                f +
                ':0 a newly added report must carry a ```sixstage fence or an explicit, non-empty ⑥ skip ' +
                'declaration (新增报告必须自带 sixstage 围栏，或给出带非空理由的 ⑥ 跳过声明)',
            });
          }
          continue;
        }
        // Existing report, in order: historical exemption -> marker -> skip -> neither.
        if (hist) {
          exempt++;
        } else if (marker) {
          pushNoFenceMarker();
        } else if (skip) {
          skipped++;
        } else {
          bare++;
        }
        continue;
      }
      fenced++;
      const added = new Set(ctx.addedLineNumbers(f));
      for (const b of blocks) {
        const push = (lineNo, text, why) => findings.push({ file: f, line: lineNo, text: text, detail: f + ':' + lineNo + ' ' + why });
        if (!b.closed) {
          // L9: an unterminated fence swallows the rest of the document; its body cannot be
          // a transcript, so the three faces below are not judged for this block.
          push(b.start + 1, L[b.start], '```sixstage fence opened here is never closed (the fence must be closed)');
          continue;
        }
        // F2 (v1.46): "touched" spans the WHOLE fence interval - the opening line, the
        // body, and the closing line. Counting body lines only let a change that rewrote
        // nothing but the ``` opening (or the closing) fence escape judgement entirely.
        let touched = false;
        for (let i = b.start; i <= b.end; i++) {
          if (added.has(i + 1)) {
            touched = true;
            break;
          }
        }
        if (!touched) continue;
        judged++;
        const at = (k) => b.start + 2 + k; // 1-based line number of body[k]

        // (1) sections: complete, once each, in order.
        const heads = [];
        for (let k = 0; k < b.body.length; k++) {
          const m = SIXSTAGE_SECTION_RE.exec(String(b.body[k]).trim());
          if (m) heads.push({ k: k, letter: m[1] });
        }
        const seq = heads.map((h) => h.letter).join('');
        if (seq !== SIXSTAGE_SECTION_SEQ) {
          push(b.start + 1, L[b.start], 'sixstage sections must be A→B→C→D→E exactly once (found `' + (seq || 'none') + '`)');
        }

        // (2) conclusion line: EXACTLY ONE, as the last non-empty body line, closed-set
        // value, nothing after it. The duplicate test (v1.46 round-6 C-4) is separate from
        // the position test: `conclusion -> note -> conclusion` used to pass because only
        // the FIRST conclusion line was located and the LAST non-empty line was itself a
        // valid conclusion, so the note in between escaped judgement.
        let lastNonEmpty = -1;
        for (let k = b.body.length - 1; k >= 0; k--) {
          if (String(b.body[k]).trim() !== '') {
            lastNonEmpty = k;
            break;
          }
        }
        const conclIdx = [];
        for (let k = 0; k < b.body.length; k++) {
          const c = sixstageConclusion(b.body[k]);
          if (c && c.prefix === SIXSTAGE_PREFIX) conclIdx.push(k);
        }
        let conclK = conclIdx.length ? conclIdx[0] : -1;
        if (conclIdx.length > 1) {
          const extra = conclIdx.slice(0, -1).map(at);
          push(
            at(conclIdx[0]),
            b.body[conclIdx[0]],
            'duplicate `' +
              SIXSTAGE_PREFIX +
              ':` conclusion line at ' +
              extra.join(', ') +
              ' (conclusion lines at ' +
              conclIdx.map(at).join(', ') +
              '): a fence carries exactly one, and it must be the last non-empty body line'
          );
        }
        const lastConc = lastNonEmpty < 0 ? null : sixstageConclusion(b.body[lastNonEmpty]);
        let conclValue = null;
        if (lastConc && lastConc.prefix === SIXSTAGE_PREFIX) {
          conclK = lastNonEmpty;
          const values = VERDICT_PREFIXES[SIXSTAGE_PREFIX];
          if (lastConc.value && values.indexOf(lastConc.value) >= 0) {
            conclValue = lastConc.value;
          } else {
            push(
              at(lastNonEmpty),
              b.body[lastNonEmpty],
              'conclusion value outside the closed set of `' + SIXSTAGE_PREFIX + '` :: ' + String(b.body[lastNonEmpty]).slice(0, 120)
            );
          }
        } else if (conclK >= 0) {
          push(
            at(conclK),
            b.body[conclK],
            'content after the conclusion line - a note must live OUTSIDE the ```sixstage fence'
          );
          const c = sixstageConclusion(b.body[conclK]);
          const values = VERDICT_PREFIXES[SIXSTAGE_PREFIX];
          if (c && c.value && values.indexOf(c.value) >= 0) conclValue = c.value;
        } else {
          push(
            lastNonEmpty >= 0 ? at(lastNonEmpty) : b.start + 1,
            lastNonEmpty >= 0 ? b.body[lastNonEmpty] : L[b.start],
            'the last non-empty line is not an `' + SIXSTAGE_PREFIX + ': <value>` conclusion line'
          );
        }

        // (3) section B severity rows ⇔ FINDINGS.
        const bHead = heads.filter((h) => h.letter === 'B')[0] || null;
        let sev = 0;
        if (bHead) {
          let hi = b.body.length;
          for (const h of heads) {
            if (h.k > bHead.k) {
              hi = h.k;
              break;
            }
          }
          for (let k = bHead.k + 1; k < hi; k++) {
            const line = String(b.body[k]);
            if (G7E_EXCLUDE_MARKERS.some((x) => line.includes(x))) continue;
            if (sixstageSeverityRow(line)) sev++;
          }
        }
        if (conclValue !== null) {
          if (sev > 0 && conclValue !== 'FINDINGS') {
            push(
              b.start + 1,
              L[b.start],
              'section B carries ' + sev + ' SEVERITY row(s) but the conclusion is `' + conclValue + '` (must be FINDINGS)'
            );
          } else if (sev === 0 && conclValue === 'FINDINGS') {
            push(
              b.start + 1,
              L[b.start],
              'no SEVERITY row in section B but the conclusion is `FINDINGS`'
            );
          }
        }
      }
    }
    if (!tgts.length) {
      out.rec('G7-f', 'INFO', 'sixstage transcript contract', 'no changed report in scope');
      return { findings };
    }
    if (findings.length) {
      out.rec('G7-f', 'FAIL', 'sixstage transcript contract', findings.map((x) => x.detail).join(' | '));
      return { findings };
    }
    if (judged) {
      out.rec(
        'G7-f',
        'PASS',
        'sixstage transcript contract',
        judged + ' changed ```sixstage fence(s) complete (A→B→C→D→E, closed-set conclusion, SEVERITY consistency)'
      );
      return { findings };
    }
    if ((skipped || exempt) && !bare) {
      const parts = [];
      if (skipped) parts.push(skipped + ' changed report(s) waived ⑥ by declaration (`- ⑥ 跳过：` / `- ⑥ skipped: `)');
      if (exempt) {
        parts.push(
          exempt +
            ' changed report(s) via the HISTORICAL EXEMPTION channel (`- ⑥ 历史豁免：` / `- ⑥ historical exemption: `), no fence required'
        );
      }
      out.rec('G7-f', 'PASS', 'sixstage transcript contract', parts.join('; '));
      return { findings };
    }
    // No changed fence to judge and no contradiction: the existence obligation is not a
    // mechanical criterion (same boundary G7-e states for the verdict line).
    out.rec(
      'G7-f',
      'INFO',
      'sixstage transcript contract',
      'no changed ```sixstage fence in scope (' +
        fenced +
        ' report(s) carry an unchanged fence, ' +
        bare +
        ' carry neither fence nor skip declaration)'
    );
    return { findings };
  },
};

module.exports = {
  id: 'G7',
  label: 'review-bundle completeness',
  title: '审查包完整性',
  scripted: true,
  subchecks: [G7_A, G7_B, G7_C, G7_D, G7_E, G7_F],
  internals: {
    targets,
    artifactBlock,
    VERDICT_PREFIXES,
    VERDICT_VALUES,
    RE_PLACEHOLDER_RE,
    verdictValue,
    verdictFinding,
    DISPOSITIONS,
    DISPOSITION_TOKENS,
    FINDING_ROW_RE,
    DISPOSITION_CELL_RE,
    rowCells,
    dispositionCell,
    STOP_TOKENS,
    STOP_MARKS,
    hasStopMark,
    STOP_OK,
    // G7-f (⑥ output contract): the fence/section/severity shapes. The closed-set table
    // itself stays VERDICT_PREFIXES above - never copied into this sub-check.
    SIXSTAGE_FENCE_RE,
    SIXSTAGE_SECTION_RE,
    SIXSTAGE_SECTION_SEQ,
    SIXSTAGE_PREFIX,
    SIXSTAGE_SKIP_RE,
    SIXSTAGE_HIST_RE,
    sixstageBlocks,
    sixstageConclusion,
    sixstageSeverityRow,
  },
  run(ctx, out) {
    for (const s of this.subchecks) s.run(ctx, out);
  },
};

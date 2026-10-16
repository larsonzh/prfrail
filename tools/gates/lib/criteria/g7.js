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

module.exports = {
  id: 'G7',
  label: 'review-bundle completeness',
  title: '审查包完整性',
  scripted: true,
  subchecks: [G7_A, G7_B, G7_C, G7_D, G7_E],
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
  },
  run(ctx, out) {
    for (const s of this.subchecks) s.run(ctx, out);
  },
};

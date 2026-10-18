'use strict';

/**
 * Base criterion G1 (DELIVERY_DIRECTIVE 6.3): status-marker consistency.
 *
 * Three places carry slice state in the ledger: the section heading
 * (`## A0 — name · ✅`), the status table (`| A0 name | ✅ 完成 | … |`) and the
 * completion log (`| 2026-09-14 | A0 | summary | evidence |`). Until v1.36 the runner
 * printed SKIP and a human compared them; the GATES-EDGE-COVERAGE defect (a completion
 * row saying 进行中 for a closed slice) survived a whole turn that way.
 *
 * v1.36 (slice GATES-G1G5, independent ① design) mechanises the part that can be
 * anchored objectively and states the rest as recorded boundaries:
 *
 *   G1(c) every state symbol belongs to the ledger's own closed set
 *         (⬜ / 🔄 / ✅ / ⛔); ⏸ is legal only on a heading, followed by
 *         "待评估" / "Pending evaluation";
 *   G1(d) a key present in both heading and status table carries the same symbol, and a
 *         status row whose key has no heading fails (a phantom row);
 *   G1(e) a completion-log key bound to a heading requires that heading to be ✅.
 *         Binding: exact primary key first (`P11`), else a UNIQUE alias from the
 *         backticked tokens of the heading name (`GATES-G1G5`) - without aliases the
 *         governance slices would never bind and (e) would be vacuous for them;
 *   G1(f) an ADDED completion row must not carry in-progress wording (historical rows
 *         are never re-judged - same scope rule as G6/G7; a later row may legitimately
 *         narrate that a heading stayed 🔄 until ⑧b).
 *
 * Recorded as unmechanised (written into 6.3 and the slice report): a stale ⬜ heading
 * whose slice actually closed but has no completion row (needs knowledge from outside
 * the slice), a ✅ heading without a completion row (four instances today, so the
 * reverse implication is not satisfiable), summary truthfulness, and the §1.6
 * stopped-state exemption.
 */

const { UsageError } = require('../ctx');

const LEDGERS = ['docs/t027/REMAINING_SLICES.md', 'docs/t027/REMAINING_SLICES_EN.md'];
const SLICE_STATES = ['⬜', '🔄', '✅', '⛔'];
const PENDING = '⏸';
const COMPLETED = '✅';
const IN_PROGRESS_RE = /进行中|In progress|待[^\n|]{0,12}定稿|await(?:s|ing)? finali[sz]/i;
/** G6 exclusion markers: a row quoting a rule or narrating history is not a live claim. */
const EXCLUSION_RE = /不得残留|示例|引文|冲突|矛盾/;

function splitRow(line) {
  return line
    .replace(/^\|/, '')
    .replace(/\|\s*$/, '')
    .split('|')
    .map((s) => s.trim());
}

/** Leading identifier of a cell, backticks stripped (`A0`, `B3a`, `P11`, `GATES-TIGHTEN`). */
function keyOf(cell) {
  const t = String(cell).replace(/`/g, '').trim();
  const m = t.match(/^([A-Za-z][A-Za-z0-9]*(?:[.\-][A-Za-z0-9]+)*)/);
  return m ? m[1] : t.split(/\s+/)[0];
}

/** Leading state symbol of a state cell ('' when the cell opens with prose). */
function symbolOf(cell) {
  const t = String(cell).trim();
  for (const s of SLICE_STATES.concat([PENDING])) if (t.startsWith(s)) return s;
  return '';
}

/**
 * Split a heading into `{key, aliases, state}`. The em dash anchors slice headings
 * (`## A0 — name`) and separates them from structural headings
 * (`### DR-1 slice-definition write-back …`, which has no dash). The state segment is
 * the LAST ` · ` whose remainder starts with a closed-set symbol, so a name containing
 * ` · ` stays intact.
 */
function parseHeading(line) {
  const m = line.match(/^(#{2,3})\s+([A-Za-z][A-Za-z0-9._-]*)\s+—\s+(.+)$/);
  if (!m) return null;
  const key = m[2];
  const rest = m[3];
  const parts = rest.split(' · ');
  let state = '';
  let namePart = rest;
  if (parts.length > 1) {
    const tail = parts[parts.length - 1];
    if (symbolOf(tail)) {
      state = tail;
      namePart = parts.slice(0, -1).join(' · ');
    }
  }
  const aliases = [...namePart.matchAll(/`([^`]+)`/g)].map((x) => x[1].trim());
  // A tail segment that opens with a symbol-ish character but is not in the closed set is an
  // UNKNOWN state (`· 🚧 进行中`), not part of the name: symbolOf() returns "" for it and every
  // guard below (`h.symbol && ...`) used to let it through silently. A tail that opens with a
  // letter or digit stays part of the name (P3 has no state segment at all).
  const unknownTail = state === '' && parts.length > 1 && /^[^\p{L}\p{N}\s]/u.test(parts[parts.length - 1]) ? parts[parts.length - 1] : '';
  return { key, aliases, state, symbol: symbolOf(state), unknownTail, name: namePart, level: m[1].length };
}

function parseLedger(text) {
  const lines = String(text).split('\n');
  const headings = [];
  const status = [];
  const completion = [];
  let inStatus = false;
  let inCompletion = false;
  let declaration = false;
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (/状态符号|Status symbols/i.test(line) && SLICE_STATES.every((s) => line.includes(s))) declaration = true;
    if (/^##\s/.test(line)) {
      inStatus = /状态表|Status table/.test(line);
      inCompletion = /完成记录|Completion log/.test(line);
    }
    const h = parseHeading(line);
    if (h) {
      headings.push(Object.assign({ n: i + 1 }, h));
      continue;
    }
    if (!/^\|/.test(line) || /^\|\s*-{2,}/.test(line)) continue;
    const cells = splitRow(line);
    if (inStatus && cells.length >= 2 && /^[A-Za-z]/.test(cells[0]) && /\s/.test(cells[0])) {
      status.push({ n: i + 1, key: keyOf(cells[0]), state: cells[1], symbol: symbolOf(cells[1]) });
    } else if (inCompletion && cells.length >= 4 && /^\d{4}-\d{2}-\d{2}$/.test(cells[0])) {
      completion.push({ n: i + 1, key: keyOf(cells[1]), summary: cells[2] || '', evidence: cells[3] || '' });
    }
  }
  return { headings, status, completion, declaration };
}

/** Resolve a completion key to a heading: exact primary key first, then a unique alias. */
function bind(key, headings) {
  const exact = headings.filter((h) => h.key === key);
  if (exact.length === 1) return { heading: exact[0], how: 'key' };
  if (exact.length > 1) return { heading: null, how: 'ambiguous-key' };
  const byAlias = headings.filter((h) => h.aliases.indexOf(key) >= 0);
  if (byAlias.length === 1) return { heading: byAlias[0], how: 'alias' };
  if (byAlias.length > 1) return { heading: null, how: 'ambiguous-alias' };
  return { heading: null, how: 'unbound' };
}

function run(ctx, out) {
  const parsed = [];
  for (const rel of LEDGERS) {
    const arr = ctx.readLineArray(rel);
    if (!arr) continue;
    parsed.push(Object.assign({ rel }, parseLedger(arr.join('\n'))));
  }
  if (!parsed.length) {
    out.rec('G1', 'SKIP', 'status marker consistency', 'ledger pair not readable in this scope');
    return;
  }
  for (const p of parsed) {
    if (!p.declaration) throw new UsageError('G1: the state-symbol declaration line is missing or incomplete in ' + p.rel);
  }

  const badC = [];
  const badD = [];
  const badE = [];
  const badF = [];
  let boundCount = 0;
  let unboundCount = 0;

  for (const p of parsed) {
    const closed = new Set(SLICE_STATES);
    for (const h of p.headings) {
      if (h.symbol === PENDING) {
        if (!/待评估|Pending evaluation/i.test(h.state)) badC.push(p.rel + ':' + h.n + ' ' + h.key + ' pending marker without 待评估');
        continue;
      }
      if (h.unknownTail) badC.push(p.rel + ':' + h.n + ' ' + h.key + ' unknown state ' + JSON.stringify(h.unknownTail.slice(0, 12)));
      else if (h.symbol && !closed.has(h.symbol)) badC.push(p.rel + ':' + h.n + ' ' + h.key + ' state ' + JSON.stringify(h.state.slice(0, 10)));
    }
    for (const s of p.status) {
      if (!closed.has(s.symbol)) badC.push(p.rel + ':' + s.n + ' status ' + JSON.stringify(String(s.state).slice(0, 10)));
    }
    for (const s of p.status) {
      const h = p.headings.filter((x) => x.key === s.key);
      if (h.length === 0) { badD.push(p.rel + ':' + s.n + ' status row ' + s.key + ' has no heading'); continue; }
      if (h.length > 1) { badD.push(p.rel + ':' + s.n + ' duplicate heading key ' + s.key); continue; }
      if (h[0].symbol && s.symbol && h[0].symbol !== s.symbol) badD.push(p.rel + ' ' + s.key + ' heading=' + h[0].symbol + ' status=' + s.symbol);
    }
    for (const c of p.completion) {
      const b = bind(c.key, p.headings);
      if (!b.heading) {
        if (b.how === 'ambiguous-alias' || b.how === 'ambiguous-key') badE.push(p.rel + ':' + c.n + ' ' + c.key + ' ' + b.how);
        else unboundCount++;
        continue;
      }
      boundCount++;
      if (b.heading.symbol && b.heading.symbol !== COMPLETED) {
        badE.push(p.rel + ':' + c.n + ' ' + c.key + ' completion row but heading ' + b.heading.key + '=' + b.heading.symbol);
      }
    }
    const added = typeof ctx.addedLineNumbers === 'function' ? ctx.addedLineNumbers(p.rel) : null;
    for (const c of p.completion) {
      if (added && !added.has(c.n)) continue;
      if (EXCLUSION_RE.test(c.summary)) continue;
      if (IN_PROGRESS_RE.test(c.summary)) badF.push(p.rel + ':' + c.n + ' ' + c.key + ' in-progress wording in a completion row');
    }
  }

  const total = parsed.reduce((n, p) => n + p.headings.length + p.status.length + p.completion.length, 0);
  const rec = (id, bad, okDetail) =>
    out.rec('G1(' + id + ')', bad.length === 0 ? 'PASS' : 'FAIL', 'status marker ' + id, bad.length ? JSON.stringify(bad.slice(0, 4)) : okDetail);
  rec('c', badC, 'closed set respected (' + total + ' cells)');
  rec('d', badD, 'heading and status table agree');
  rec('e', badE, 'completion rows imply completed (' + boundCount + ' bound, ' + unboundCount + ' unbound)');
  if (unboundCount > 0) {
    out.rec('G1(e-info)', 'INFO', 'completion keys without a binding', unboundCount + ' completion row(s) have no heading to bind to (their slice heading was removed); G1(e) judges bound rows only');
  }
  rec('f', badF, 'no in-progress wording in added completion rows');
}

module.exports = {
  id: 'G1',
  label: 'status marker consistency (heading / status table / completion log)',
  title: '状态标记一致性',
  scripted: true,
  run,
  _internals: { parseLedger, parseHeading, bind, keyOf, symbolOf, splitRow, SLICE_STATES, PENDING, COMPLETED, IN_PROGRESS_RE, EXCLUSION_RE, LEDGERS },
};

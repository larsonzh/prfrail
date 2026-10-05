'use strict';

/**
 * G8 frozen evidence pack self-consistency (DELIVERY_DIRECTIVE 6.3 G8, v1.24).
 *
 * The evidence tree `docs/validation/evidence/**` carries FROZEN historical artifacts:
 * their content cannot satisfy the "new document" criteria (G1-a / G2 / G4a / G4b / G6
 * are excluded from the tree, see `ctx.EVIDENCE_TREE_PREFIX`), so it is judged here
 * instead. G8 is a SHAPE-level guard and a pack-internal consistency check; it does not
 * judge content semantics, and it does NOT catch "live code masquerading as a frozen
 * pack" (that residual risk belongs to the (5)/(6)/(7) review chain).
 *
 * Scope rule: every changed path inside the tree (`ctx.changedPaths` filtered by
 * `ctx.isEvidenceTree`); the same set in tree / index / range / ci scope. When the change
 * set touches nothing in the tree every sub-check reports INFO, never a vacuous PASS.
 *
 *   G8-a shape guard     every changed tree path must be a pack member (nearest ancestor
 *                        directory carrying both `MANIFEST.md` and `SHA256SUMS.txt`), OR a
 *                        root-level `*-freedom-list.md`, OR registered verbatim in
 *                        `tools/gates/evidence-forms.txt`.
 *   G8-b pack self-consistency
 *                        `MANIFEST.md` archive-entry set == `SHA256SUMS.txt` entry set ==
 *                        the pack payload (every file under the pack minus the two
 *                        metadata files at the pack root); every recorded sha256 is
 *                        recomputed. Missing list / payload not listed / hash mismatch /
 *                        untracked-list collision each fail closed.
 *   G8-c encoding-normalization record
 *                        a pack whose MANIFEST carries the `## 后缀与编码映射` section must
 *                        prove both sha256 values are 64-hex and that the archive-byte
 *                        value equals the `SHA256SUMS.txt` record (legacy packs without
 *                        the section are exempt).
 *   G8-d rebuilt-archive note
 *                        a changed-and-present tree file whose leading annotation zone
 *                        declares 「重建归档」 must carry the source path and the basis, and
 *                        any sha256-shaped token in that zone must be exactly 64 hex.
 *
 * Fail-closed everywhere: a malformed registration file, an unparseable MANIFEST, a pack
 * enumeration failure or a bad SHA256SUMS line is a `UsageError` (exit 2) or a FAIL -
 * never a silent pass.
 *
 * Known boundary (as declared in section 6.3): G8-b takes "pack payload" to mean EVERY
 * file under the pack minus the two pack-root metadata files. A pack-root non-metadata
 * file that is not listed therefore fails. That payload-set caliber is deliberately NOT
 * relaxed. The real corpus once exposed exactly this shape -
 * `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/CLEARED-INVENTORY.txt`
 * (a controlled 259061-byte ledger holding the previous slice's complete category-3
 * cleanup account) was tracked but absent from that pack's `## 归档条目清单` and
 * `SHA256SUMS.txt`. User ruling A (v1.24) authorized the two metadata records to be
 * backfilled (one archive-entry row plus the matching `SHA256SUMS.txt` line; the payload
 * bytes are untouched), so that pack is now self-consistent and passes green.
 *
 * Historical readings are not retroactive: at commit `4284e94` the backfill did not yet
 * exist, so `--check=G8 --scope=range:4284e94^..4284e94` keeps its two G8-a/G8-b reds
 * permanently (section 11.1); no post-hoc measure may turn the history green.
 *
 * Legacy units (measured as-is; fail-closed, by design): `b1-20260917` has neither
 * MANIFEST nor SHA256SUMS; `dr-1-2026-09-24` has SHA256SUMS but no MANIFEST; and
 * `b2-2026-09-17` / `b3a-2026-09-20` / `b3b-2026-09-21` carry both but in the pre-v1.23
 * legacy format (no `## 归档条目清单` section, backslash SHA256SUMS separators, ...).
 * When touched they are reported red - that is the intended fail-closed behaviour.
 */

const path = require('path');
const crypto = require('crypto');

const { UsageError, EVIDENCE_TREE_PREFIX, isEvidenceTree, toPosix } = require('../ctx');

const FORMS_REL = 'tools/gates/evidence-forms.txt';
const EVIDENCE_TREE_ROOT = EVIDENCE_TREE_PREFIX.replace(/\/$/, '');

const MANIFEST_LIST_HEADING = '## 归档条目清单';
const ENCODING_MAP_HEADING = '## 后缀与编码映射';
const UNTRACKED_HEADING = '## 未入库文件清单';

/** First table cell of a row whose first cell is a backticked path. */
const BACKTICK_ROW_RE = /^\|\s*`([^`]+)`/;
/** `<64 lowercase hex>  <path>` (the sha256sum shape, POSIX paths normalized later). */
const SUM_RE = /^([0-9a-f]{64})[ \t]+(.+)$/;
/** A sha256-shaped token candidate: a run of >= 32 hex digits. */
const HEX_RUN_RE = /[0-9a-fA-F]{32,}/g;
const HEX64_RE = /^[0-9a-f]{64}$/;

/** Normalize an in-pack relative path: POSIX separators, drop a leading `./`. */
function normalizeRel(p) {
  return toPosix(String(p).trim()).replace(/^\.\//, '');
}

function sha256Hex(buf) {
  return crypto.createHash('sha256').update(buf).digest('hex');
}

function finding(file, line, text, detail) {
  return { file, line, text, detail };
}

/** Changed paths inside the evidence tree (ctx.isEvidenceTree with a module fallback). */
function scopeTargets(ctx) {
  const inTree = typeof ctx.isEvidenceTree === 'function' ? ctx.isEvidenceTree : isEvidenceTree;
  return ctx.changedPaths.filter((p) => inTree(p));
}

/**
 * Pack roots of the evidence tree, derived from ONE `ctx.listTreeFiles(<evidence root>)`
 * listing (never from the change set). Memoized per context.
 *
 * Why a set: the naive "call listTreeFiles for every ancestor of every changed path" walk
 * spawns one git process per distinct ancestor directory (measured: ~700 directories for the
 * 4284e94 archive commit, ~210s). Deriving the candidates from a single listing and then
 * CONFIRMING each with the per-pack listing keeps the spec's `listTreeFiles(packDir)`
 * membership test while cutting the git spawns to 1 + (number of packs).
 *
 * Trade-off (DEF-2, round 2): the alternative "probe every candidate ancestor with two
 * existence queries (MANIFEST.md + SHA256SUMS.txt)" is deliberately NOT used. After
 * `listTreeFiles` batches its ignore lookup, this single listing costs ONE git spawn plus
 * one directory walk; the per-ancestor form would instead cost two `git cat-file -e`
 * spawns per ancestor in range/ci scope (~700 ancestors for 4284e94 ⇒ ~1400 spawns),
 * reintroducing exactly the (paths × depth) spawn order described above. Tree/index scope
 * could probe with `fs.existsSync` cheaply, but ONE derivation for every scope is what
 * keeps the pack-root set scope-independent and this code auditable.
 */
const packRootCache = new WeakMap();
function evidencePackRoots(ctx) {
  if (packRootCache.has(ctx)) return packRootCache.get(ctx);
  const files = typeof ctx.listTreeFiles === 'function' ? ctx.listTreeFiles(EVIDENCE_TREE_ROOT) : [];
  const withManifest = new Set();
  const withSums = new Set();
  for (const f of files) {
    if (f.endsWith('/MANIFEST.md')) withManifest.add(f.slice(0, f.length - '/MANIFEST.md'.length));
    if (f.endsWith('/SHA256SUMS.txt')) withSums.add(f.slice(0, f.length - '/SHA256SUMS.txt'.length));
  }
  const set = new Set();
  for (const d of withManifest) if (withSums.has(d)) set.add(d);
  packRootCache.set(ctx, set);
  return set;
}

/**
 * First ancestor directory of `p` (including its own directory) that carries BOTH
 * `MANIFEST.md` and `SHA256SUMS.txt`. Candidate ancestors come from the single tree
 * listing above; the membership of each candidate is then CONFIRMED through
 * `ctx.listTreeFiles(packDir)` - so a path is recognized as a pack member even when only
 * one file of the pack changed, and the test never falls back to the change set.
 */
function findPackRoot(ctx, p) {
  if (typeof ctx.listTreeFiles !== 'function') return null;
  const candidates = evidencePackRoots(ctx);
  let dir = path.posix.dirname(toPosix(p));
  while (dir && dir !== '.' && dir !== '/') {
    if (candidates.has(dir)) {
      const files = ctx.listTreeFiles(dir);
      if (files.indexOf(dir + '/MANIFEST.md') >= 0 && files.indexOf(dir + '/SHA256SUMS.txt') >= 0) return dir;
    }
    const parent = path.posix.dirname(dir);
    if (parent === dir) return null;
    dir = parent;
  }
  return null;
}

/** Unique, sorted pack roots touched by the target set. */
function touchedPacks(ctx, tgts) {
  const set = new Set();
  for (const f of tgts) {
    const r = findPackRoot(ctx, f);
    if (r) set.add(r);
  }
  return [...set].sort();
}

/**
 * Body lines of the FIRST `## <headingPrefix>...` section, or null when the heading is
 * absent. The heading is matched by prefix because real headings carry a parenthetical
 * suffix (e.g. `## 后缀与编码映射（`.md` → `.txt` 及编码归一）`).
 */
function sectionBody(L, headingPrefix) {
  let start = -1;
  for (let i = 0; i < L.length; i++) {
    if (L[i].startsWith(headingPrefix)) {
      start = i;
      break;
    }
  }
  if (start < 0) return null;
  const body = [];
  for (let i = start + 1; i < L.length; i++) {
    if (/^## /.test(L[i])) break;
    body.push(L[i]);
  }
  return body;
}

/** First-column backticked paths of every `` | `path` | ... `` row in a section body. */
function backtickFirstCells(body) {
  const out = [];
  for (const line of body) {
    const m = BACKTICK_ROW_RE.exec(line);
    if (m) out.push(normalizeRel(m[1]));
  }
  return out;
}

/** Trimmed cells of a markdown row; the trailing empty cell is preserved. */
function tableCells(line) {
  const parts = String(line).split('|');
  const cells = String(line).trim().endsWith('|') ? parts.slice(1, -1) : parts.slice(1);
  return cells.map((c) => c.trim());
}

function stripTicks(s) {
  return String(s).trim().replace(/^`+/, '').replace(/`+$/, '').trim();
}

/** Parse a SHA256SUMS body: path -> recorded hash, plus the malformed lines. */
function parseSums(text) {
  const map = new Map();
  const bad = [];
  const lines = String(text).split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line === '') continue;
    const m = SUM_RE.exec(line);
    if (!m) {
      bad.push({ line: i + 1, text: line });
      continue;
    }
    const p = normalizeRel(m[2]);
    if (map.has(p)) {
      bad.push({ line: i + 1, text: line });
      continue;
    }
    map.set(p, m[1]);
  }
  return { map, bad };
}

/**
 * Registration table loader. Reads through `ctx.loadWhitelist` (same scope semantics as
 * G6 / G4b: tree/index -> working tree, range/ci -> the scope endpoint). Breaks on a
 * field count other than 3, an empty field, a duplicate path, and any parse error.
 */
function loadForms(ctx) {
  let res;
  try {
    res = ctx.loadWhitelist(FORMS_REL);
  } catch (e) {
    throw new UsageError('cannot read ' + FORMS_REL + ': ' + (e && e.message ? e.message : e));
  }
  if (!res) return { entries: [], source: 'unknown' };
  if (res.error) throw new UsageError(res.error);
  const text = res.text === undefined ? null : res.text;
  // Absent at the scope endpoint (or empty): zero registrations. A required path then
  // FAILs the shape guard - it is never silently accepted.
  if (text === null || text === '') return { entries: [], source: res.source || 'unknown' };
  const entries = [];
  const seen = new Set();
  const lines = String(text).split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (!line.trim() || line.startsWith('#')) continue;
    const parts = line.split('\t');
    if (parts.length !== 3) {
      throw new UsageError(
        FORMS_REL +
          ':' +
          (i + 1) +
          ' malformed registration (expected exactly 3 TAB-separated fields: match<TAB>file<TAB>reason)'
      );
    }
    const match = parts[0].trim();
    const file = toPosix(parts[1].trim());
    const reason = parts[2].trim();
    if (!match || !file || !reason) {
      throw new UsageError(FORMS_REL + ':' + (i + 1) + ' malformed registration (empty match/file/reason field)');
    }
    if (seen.has(file)) {
      throw new UsageError(FORMS_REL + ':' + (i + 1) + ' duplicate registration for ' + file);
    }
    seen.add(file);
    entries.push({ match, file, reason });
  }
  return { entries, source: res.source || 'unknown' };
}

/** Payload of a pack: every file under it minus the two pack-root metadata files. */
function packPayload(ctx, pack) {
  return new Set(
    ctx
      .listTreeFiles(pack)
      .filter((p) => p.startsWith(pack + '/'))
      .map((p) => p.slice(pack.length + 1))
      .filter((r) => r && r !== 'MANIFEST.md' && r !== 'SHA256SUMS.txt')
  );
}

/* ------------------------------------------------------------------------------------------
 * G8-a: evidence-tree path shape
 * --------------------------------------------------------------------------------------- */

const G8_A = {
  id: 'G8-a',
  label: 'evidence-tree path shape',
  run(ctx, out) {
    const tgts = scopeTargets(ctx);
    const forms = loadForms(ctx);
    const findings = [];
    for (const f of tgts) {
      const member = findPackRoot(ctx, f) !== null;
      const dir = path.posix.dirname(toPosix(f));
      const base = path.posix.basename(toPosix(f));
      const freedom = dir === EVIDENCE_TREE_ROOT && /freedom-list\.md$/i.test(base);
      const registered = forms.entries.some((e) => e.match === f && e.file === f);
      if (!member && !freedom && !registered) {
        findings.push(
          finding(
            f,
            0,
            f,
            f + ' not a pack member / not a root freedom-list / not registered in ' + FORMS_REL
          )
        );
      }
    }
    if (!tgts.length) {
      out.rec('G8-a', 'INFO', 'evidence-tree path shape', 'no changed evidence-tree path in scope');
      return { findings };
    }
    out.rec(
      'G8-a',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'evidence-tree path shape',
      findings.length ? findings.map((f) => f.detail).join(' | ') : tgts.length + ' path(s) OK'
    );
    return { findings };
  },
};

/* ------------------------------------------------------------------------------------------
 * G8-b: pack self-consistency
 * --------------------------------------------------------------------------------------- */

function checkPackConsistency(ctx, pack, findings) {
  const manifestRel = pack + '/MANIFEST.md';
  const sumsRel = pack + '/SHA256SUMS.txt';
  const payload = packPayload(ctx, pack);

  const mText = ctx.readText ? ctx.readText(manifestRel) : null;
  let manifestSet = new Set();
  let manifestUsable = false;
  if (mText === null) {
    findings.push(finding(manifestRel, 0, manifestRel, pack + ': MANIFEST.md unreadable'));
  } else {
    const body = sectionBody(mText.split('\n'), MANIFEST_LIST_HEADING);
    if (body === null) {
      findings.push(
        finding(manifestRel, 0, MANIFEST_LIST_HEADING, pack + ': missing `' + MANIFEST_LIST_HEADING + '` section (missing list)')
      );
    } else {
      manifestSet = new Set(backtickFirstCells(body));
      manifestUsable = true;
    }
  }

  const sText = ctx.readText ? ctx.readText(sumsRel) : null;
  let sumsMap = new Map();
  if (sText === null) {
    findings.push(finding(sumsRel, 0, sumsRel, pack + ': SHA256SUMS.txt unreadable/missing'));
  } else {
    const parsed = parseSums(sText);
    sumsMap = parsed.map;
    for (const b of parsed.bad) {
      findings.push(finding(sumsRel, b.line, b.text, pack + ': malformed SHA256SUMS line ' + b.line + ': ' + b.text.slice(0, 80)));
    }
  }

  // Set equality: MANIFEST list == SHA256SUMS entries == pack payload.
  for (const r of [...payload].sort()) {
    if (!manifestSet.has(r)) {
      findings.push(
        finding(manifestRel, 0, r, pack + '/' + r + ' not listed in `' + MANIFEST_LIST_HEADING + '` (payload not listed)')
      );
    }
    if (!sumsMap.has(r)) {
      findings.push(finding(sumsRel, 0, r, pack + '/' + r + ' not listed in SHA256SUMS.txt (payload not listed)'));
    }
  }
  if (manifestUsable) {
    for (const r of [...manifestSet].sort()) {
      if (!payload.has(r)) {
        findings.push(
          finding(manifestRel, 0, r, pack + '/' + r + ' listed in `' + MANIFEST_LIST_HEADING + '` but absent from pack payload')
        );
      }
    }
  }
  for (const r of [...sumsMap.keys()].sort()) {
    if (!payload.has(r)) {
      findings.push(finding(sumsRel, 0, r, pack + '/' + r + ' listed in SHA256SUMS.txt but absent from pack payload'));
    }
  }

  // Recompute every recorded hash from the scope-aware byte read.
  for (const [r, recorded] of [...sumsMap.entries()].sort((a, b) => (a[0] < b[0] ? -1 : 1))) {
    const bytes = ctx.readBytes(pack + '/' + r);
    if (bytes === null) {
      findings.push(finding(sumsRel, 0, r, pack + '/' + r + ': hash source unreadable'));
      continue;
    }
    const actual = sha256Hex(bytes);
    if (actual !== recorded) {
      findings.push(
        finding(
          sumsRel,
          0,
          r,
          pack + '/' + r + ' hash mismatch (recorded ' + recorded.slice(0, 12) + '..., actual ' + actual.slice(0, 12) + '...)'
        )
      );
    }
  }

  // `## 未入库文件清单` entries must be absolute-looking and must not collide with payload.
  if (mText !== null) {
    const body = sectionBody(mText.split('\n'), UNTRACKED_HEADING);
    if (body !== null) {
      for (const raw of backtickFirstCells(body)) {
        if (!raw.startsWith('/')) {
          findings.push(finding(manifestRel, 0, raw, pack + ': 未入库 entry not absolute: ' + raw));
          continue;
        }
        const key = raw.replace(/^\//, '');
        if (payload.has(key)) {
          findings.push(
            finding(manifestRel, 0, raw, pack + ': untracked-list conflict - ' + raw + ' collides with pack payload ' + key)
          );
        }
      }
    }
  }
}

const G8_B = {
  id: 'G8-b',
  label: 'pack self-consistency',
  run(ctx, out) {
    const tgts = scopeTargets(ctx);
    const packs = touchedPacks(ctx, tgts);
    const findings = [];
    for (const pack of packs) checkPackConsistency(ctx, pack, findings);
    if (!packs.length) {
      out.rec('G8-b', 'INFO', 'pack self-consistency', 'no touched evidence pack in scope');
      return { findings };
    }
    out.rec(
      'G8-b',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'pack self-consistency',
      findings.length ? findings.map((f) => f.detail).join(' | ') : packs.length + ' pack(s) self-consistent'
    );
    return { findings };
  },
};

/* ------------------------------------------------------------------------------------------
 * G8-c: encoding-normalization record
 * --------------------------------------------------------------------------------------- */

function checkPackEncoding(ctx, pack, findings) {
  const manifestRel = pack + '/MANIFEST.md';
  const mText = ctx.readText ? ctx.readText(manifestRel) : null;
  if (mText === null) {
    findings.push(finding(manifestRel, 0, manifestRel, pack + ': MANIFEST.md unreadable'));
    return { checked: false, exempt: false };
  }
  const body = sectionBody(mText.split('\n'), ENCODING_MAP_HEADING);
  if (body === null) return { checked: false, exempt: true };

  const payload = packPayload(ctx, pack);
  const sText = ctx.readText ? ctx.readText(pack + '/SHA256SUMS.txt') : null;
  const sumsMap = sText === null ? new Map() : parseSums(sText).map;

  for (const raw of body) {
    if (!BACKTICK_ROW_RE.test(raw)) continue;
    const cells = tableCells(raw);
    if (cells.length !== 5) {
      findings.push(finding(manifestRel, 0, raw, pack + ': malformed 后缀与编码映射 row (expected 5 cells)'));
      continue;
    }
    const origPath = stripTicks(cells[0]);
    const archPath = stripTicks(cells[1]);
    const origHash = stripTicks(cells[2]);
    const archHash = stripTicks(cells[3]);
    if (!HEX64_RE.test(origHash)) {
      findings.push(finding(manifestRel, 0, raw, pack + ': 原始字节 sha256 is not 64 hex for ' + origPath));
      continue;
    }
    if (!HEX64_RE.test(archHash)) {
      findings.push(finding(manifestRel, 0, raw, pack + ': 归档字节 sha256 is not 64 hex for ' + archPath));
      continue;
    }
    if (!payload.has(archPath)) {
      findings.push(finding(manifestRel, 0, raw, pack + ': 归档路径 not in pack payload: ' + archPath));
      continue;
    }
    if (!sumsMap.has(archPath)) {
      findings.push(finding(manifestRel, 0, raw, pack + ': 归档路径 not recorded in SHA256SUMS: ' + archPath));
      continue;
    }
    const recorded = sumsMap.get(archPath);
    if (recorded !== archHash) {
      findings.push(
        finding(
          manifestRel,
          0,
          raw,
          pack +
            ': archive hash mismatch for ' +
            archPath +
            ' (MANIFEST ' +
            archHash.slice(0, 12) +
            '... vs SHA256SUMS ' +
            recorded.slice(0, 12) +
            '...)'
        )
      );
    }
  }
  return { checked: true, exempt: false };
}

const G8_C = {
  id: 'G8-c',
  label: 'encoding-normalization record',
  run(ctx, out) {
    const tgts = scopeTargets(ctx);
    const packs = touchedPacks(ctx, tgts);
    const findings = [];
    let checked = 0;
    let exempt = 0;
    for (const pack of packs) {
      const r = checkPackEncoding(ctx, pack, findings);
      if (r.checked) checked++;
      else if (r.exempt) exempt++;
    }
    if (!packs.length) {
      out.rec('G8-c', 'INFO', 'encoding-normalization record', 'no touched evidence pack in scope');
      return { findings };
    }
    out.rec(
      'G8-c',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'encoding-normalization record',
      findings.length
        ? findings.map((f) => f.detail).join(' | ')
        : checked +
            ' pack(s) checked' +
            (exempt ? ' (' + exempt + ' exempt: no `' + ENCODING_MAP_HEADING + '` section)' : '')
    );
    return { findings };
  },
};

/* ------------------------------------------------------------------------------------------
 * G8-d: rebuilt-archive note
 * --------------------------------------------------------------------------------------- */

/**
 * Leading annotation zone: the opening run of `#`-headed, `>`-quoted or empty lines; the
 * zone stops at the first line that is neither, capped at 15 lines.
 */
function annotationZone(L) {
  const zone = [];
  for (let i = 0; i < L.length && zone.length < 15; i++) {
    const line = L[i];
    if (line === '' || line.startsWith('#') || line.startsWith('>')) zone.push({ n: i + 1, text: line });
    else break;
  }
  return zone;
}

const G8_D = {
  id: 'G8-d',
  label: 'rebuilt-archive note',
  run(ctx, out) {
    // Target = changed AND present tree files (a deleted file has no annotation zone).
    const tgts = scopeTargets(ctx).filter((f) => ctx.readBytes(f) !== null);
    const findings = [];
    for (const f of tgts) {
      const text = ctx.readText ? ctx.readText(f) : null;
      if (text === null) continue;
      const zone = annotationZone(text.split('\n'));
      const joined = zone.map((z) => z.text).join('\n');
      if (!joined.includes('重建归档')) continue;
      const hasSource = /源路径|源文件/.test(joined) || /`tmp\/[^`]*`/.test(joined);
      const hasBasis = joined.includes('依据');
      if (!hasSource) {
        findings.push(finding(f, 0, f, f + ': 重建归档 annotation missing source path (源路径/源文件 or a `tmp/...` path)'));
      }
      if (!hasBasis) {
        findings.push(finding(f, 0, f, f + ': 重建归档 annotation missing basis (依据)'));
      }
      for (const m of joined.matchAll(HEX_RUN_RE)) {
        if (m[0].length !== 64) {
          findings.push(
            finding(f, 0, m[0], f + ': bad sha256 token in 重建归档 annotation: ' + m[0].slice(0, 20) + '... (' + m[0].length + ' digits)')
          );
        }
      }
    }
    if (!tgts.length) {
      out.rec('G8-d', 'INFO', 'rebuilt-archive note', 'no changed evidence file in scope');
      return { findings };
    }
    out.rec(
      'G8-d',
      findings.length === 0 ? 'PASS' : 'FAIL',
      'rebuilt-archive note',
      findings.length ? findings.map((f) => f.detail).join(' | ') : tgts.length + ' changed file(s) OK'
    );
    return { findings };
  },
};

module.exports = {
  id: 'G8',
  label: 'frozen evidence pack self-consistency',
  title: '冻结证据包自洽性核验',
  scripted: true,
  subchecks: [G8_A, G8_B, G8_C, G8_D],
  internals: {
    FORMS_REL,
    EVIDENCE_TREE_PREFIX,
    EVIDENCE_TREE_ROOT,
    MANIFEST_LIST_HEADING,
    ENCODING_MAP_HEADING,
    UNTRACKED_HEADING,
    scopeTargets,
    findPackRoot,
    touchedPacks,
    sectionBody,
    backtickFirstCells,
    tableCells,
    stripTicks,
    parseSums,
    loadForms,
    packPayload,
    checkPackConsistency,
    checkPackEncoding,
    annotationZone,
    sha256Hex,
    normalizeRel,
  },
  run(ctx, out) {
    for (const s of this.subchecks) s.run(ctx, out);
  },
};

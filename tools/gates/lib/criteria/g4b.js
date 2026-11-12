'use strict';

/**
 * G4b anomalous character scan (DELIVERY_DIRECTIVE 6.3 G4).
 *     (CJK characters of the added lines) MINUS (CJK characters of the base-revision
 *     *.md corpus) MINUS (tools/gates/cjk-newwords.txt) MUST be empty.
 * The whitelist freezes the human judgement (each entry carries a reason + context
 * excerpt) so the assertion stays mechanically re-runnable. The base revision is the
 * pre-change corpus: HEAD for tree/index scopes, the range base for range scopes.
 *
 * The whitelist itself is read through the shared scope-aware loader
 * (`ctx.loadWhitelist`, GATES-EXT F1): tree/index use the working tree, range/ci use
 * the scope's ENDPOINT revision - exactly like G6's whitelist, so both criteria
 * always judge against the same snapshot of the registry.
 *
 * Slice G4b-CACHE (v1.49): the base-revision markdown corpus is the expensive half of
 * the scan (listAllMarkdown plus one `git show` per file) and it is a PURE FUNCTION of
 * the corpus revision, so it is cached across processes. The cache stores ONLY that
 * character set - never the verdict, never the added-line characters, never the
 * whitelist - so a hit cannot hide a finding. The key binds the resolved sha, the
 * whitelist text, this module's own source and the format version, and it is
 * recomputed BEFORE every read. Any anomaly (no sha, unreadable or malformed entry,
 * failed write) silently falls back to the full scan: the cache is an optimisation,
 * never a second source of truth.
 */
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const { UsageError } = require('../ctx');

const WHITELIST_REL = 'tools/gates/cjk-newwords.txt';
/**
 * Cache format version: the first key segment AND the `v` field of a stored value.
 * A change to the CORPUS-SEMANTIC inputs must bump this version: `cjkChars`,
 * the corpus loop inside `run` and `corpusRef` are covered by the source
 * digest, but `ctx.listAllMarkdown` / `ctx.showAt` live in ctx.js and are
 * NOT - if their enumeration or read semantics change (e.g. a new suffix or
 * ref), every entry written under the old semantics is stale and must miss.
 * This slice (G4B-CORPUS-DECOUPLE) changed the enumeration semantics of
 * `ctx.listAllMarkdown()` (index -> object store), so under the same obligation
 * every entry written before the bump is treated as stale.
 */
const CACHE_VERSION = 'v2';

/**
 * Upper bound on the number of on-disk corpus entries (LRU eviction, slice P67 #1).
 * Before this slice the directory grew without a bound - it held ~96 entries (one per
 * resolved sha, most of them permanently dead once HEAD moved on). 100 keeps all of
 * those and leaves a little headroom, so eviction only kicks in for genuinely stale
 * records, whose least recently used entry can never be the one a live run wants.
 */
const CACHE_CAP = 100;

/** A cache entry file is exactly the sha256 hex of its key plus the `.json` suffix
 *  (atomic temp files are `.<tmp>.<pid>.<seq>` and must never be eviction candidates). */
const CACHE_ENTRY_RE = /^[0-9a-f]{64}\.json$/;

function cjkChars(s) {
  const o = new Set();
  for (const ch of s) {
    const c = ch.codePointAt(0);
    if ((c >= 0x4e00 && c <= 0x9fff) || (c >= 0x3400 && c <= 0x4dbf) || (c >= 0xf900 && c <= 0xfaff)) {
      o.add(ch);
    }
  }
  return o;
}

/** The base revision the corpus is read from (unchanged pre-change scan ref). */
function corpusRef(ctx) {
  return ctx.kind === 'range' || ctx.kind === 'ci' ? ctx.baseRef : 'HEAD';
}

function sha256(s) {
  return crypto.createHash('sha256').update(String(s), 'utf8').digest('hex');
}

/**
 * The `run` source text participates in the cache key. `toString()` of the CLOSURE is
 * load-bearing: reading `__filename` from disk would return the shipped file even inside
 * an in-memory mutated module, so the source-digest counterexample could never go red.
 */
function runBody() {
  return run.toString();
}

/** Every function that influences the STORED corpus value must be covered by this
 *  digest or by a key segment: cjkChars + the corpus loop inside `run` are hashed,
 *  `corpusRef` is covered by the resolved-sha segment, and anything that only gates
 *  hit/miss (cacheUsable, validValue, store plumbing) intentionally is not - a wrong
 *  hit can only return a corpus computed by identical corpus-semantic code. */
function sourceDigest() {
  return sha256(JSON.stringify([WHITELIST_REL, cjkChars.toString(), runBody()]));
}

/** Per-process sequence for the atomic temp file names. */
let cacheSeq = 0;

function cacheDir(ctx) {
  return ctx.g4bCacheDir || path.join(ctx.repoRoot, ...['tmp', 'gates-cache', 'g4b']);
}

function gitOut(ctx, args) {
  try {
    const r = ctx.git(args);
    return r && r.ok && typeof r.out === 'string' ? r.out : null;
  } catch (e) {
    return null;
  }
}

/** Resolved corpus revision, or null when git cannot answer (cache then stays off). */
function resolvedSha(ctx) {
  const out = gitOut(ctx, ['rev-parse', corpusRef(ctx)]);
  if (out === null) return null;
  return out.trim() || null;
}

/**
 * `<version>|<resolved sha>|<scope kind>|<whitelist source>|<whitelist text digest>|
 * <criterion source digest>` - always recomputed before a read, never persisted alone.
 */
function keyString(ctx, whitelist) {
  const wl = whitelist || ctx.loadWhitelist(WHITELIST_REL);
  const rs = resolvedSha(ctx);
  const text = wl && typeof wl.text === 'string' ? wl.text : '';
  return [
    CACHE_VERSION,
    rs === null ? 'none' : rs,
    // Deliberate namespace separator, not a discriminator for the enumeration
    // command: both branches now read the object store. Keeping the kind segment
    // prevents two different scopes from sharing one cache record.
    String(ctx.kind || ''),
    String((wl && wl.source) || ''),
    sha256(text),
    sourceDigest(),
  ].join('|');
}

function fullKey(ctx, whitelist) {
  return sha256(keyString(ctx, whitelist));
}

/**
 * Cache eligibility. The stored value is a PURE function of the corpus revision:
 * `listAllMarkdown()` reads the SAME corpus revision from the object store in both
 * scopes (tree/index = `HEAD`, range/ci = the range base), and every file's
 * content is read with `git show <ref>:<p>`. Staged additions, staged
 * modifications and staged deletions of `.md` therefore cannot change the
 * computed corpus: the cache and a cold run are IDENTICAL on the corpus
 * (OB-120 closed). The cache stores ONLY that character set, never the verdict,
 * never the added-line characters and never the whitelist, so a hit cannot hide a
 * finding; any anomaly (no sha, unreadable or malformed entry, failed write)
 * fails open to the full recomputation.
 */
function cacheUsable(ctx) {
  if (ctx.g4bCacheDisabled) return false;
  if (resolvedSha(ctx) === null) return false;
  return true;
}

/** A stored value must be exactly `{v: 2, corpus: string[]}`; anything else is a miss. */
function validValue(v) {
  return (
    v !== null &&
    typeof v === 'object' &&
    v.v === 2 &&
    Array.isArray(v.corpus) &&
    v.corpus.every((c) => typeof c === 'string')
  );
}

/**
 * Drop the least recently used entries of `dir` until at most `cap` remain.
 *
 * Access order lives in the file system, never in the stored value: `put` writes a
 * fresh mtime and `get` touches a hit, so the oldest mtime is the entry no run has
 * needed for the longest time. Only files named `<64 hex>.json` count as entries;
 * anything else (atomic temp files, unrelated artefacts) is ignored.
 *
 * `statFn` / `unlinkFn` are injectable so the selftest can drive real mtimes on a
 * temporary directory. EVERY failure is swallowed: eviction is an optimisation and
 * may never change a verdict. A losing race with another process (Windows keeps an
 * open handle, so its `unlink` fails) simply leaves that entry for the next pass, and
 * an over-eager concurrent pass can only make the directory slightly smaller.
 *
 * Returns the number of entries actually removed (informational; never throws).
 */
function evictToCap(dir, cap, statFn, unlinkFn) {
  const stat = statFn || fs.statSync;
  const unlink = unlinkFn || fs.unlinkSync;
  let names;
  try {
    names = fs.readdirSync(dir);
  } catch (e) {
    return 0; // unreadable directory: nothing we can enumerate, nothing to evict
  }
  const entries = [];
  for (const name of names) {
    if (!CACHE_ENTRY_RE.test(name)) continue;
    const file = path.join(dir, name);
    try {
      entries.push({ file, mtime: stat(file).mtimeMs });
    } catch (e) {
      /* vanished or unreadable: not an eviction candidate */
    }
  }
  const excess = entries.length - cap;
  if (!(excess > 0)) return 0;
  entries.sort((a, b) => a.mtime - b.mtime); // oldest first (stable: ties keep readdir order)
  let removed = 0;
  for (let i = 0; i < excess; i += 1) {
    try {
      unlink(entries[i].file);
      removed += 1;
    } catch (e) {
      /* fail open: a racing reader/writer keeps this entry alive for now */
    }
  }
  return removed;
}

function diskStore(ctx) {
  const dir = cacheDir(ctx);
  // Created eagerly: an unusable directory disables the cache BEFORE any record exists.
  fs.mkdirSync(dir, { recursive: true });
  return {
    get(k) {
      const file = path.join(dir, k + '.json');
      let parsed;
      try {
        parsed = JSON.parse(fs.readFileSync(file, 'utf8'));
      } catch (e) {
        return null; // absent / unreadable / malformed: a miss, never an error
      }
      // LRU access stamp: the entry file's mtime is the visit record, so the stored
      // value keeps its exact `{v: 2, corpus}` shape and the hit semantics stay
      // "identical to a cold run". Touching happens ONLY after a successful parse
      // (a failed read/parse is a miss and records nothing), and a failed touch is
      // harmless - the entry merely looks older - so it must never unseat the hit.
      try {
        const now = new Date();
        fs.utimesSync(file, now, now);
      } catch (e) {
        /* fail open: the hit stands even when the timestamp cannot be written */
      }
      return parsed;
    },
    put(k, v) {
      const target = path.join(dir, k + '.json');
      const tmp = path.join(dir, k + '.tmp.' + process.pid + '.' + (cacheSeq += 1));
      try {
        fs.writeFileSync(tmp, JSON.stringify(v), 'utf8');
        fs.renameSync(tmp, target);
      } finally {
        try {
          fs.unlinkSync(tmp);
        } catch (e) {
          /* the rename consumed it (or a racing process did) - nothing to clean */
        }
      }
      // The rename above refreshed `target`'s mtime, i.e. this entry is now the most
      // recently used one; drop the oldest entries so the directory stays bounded.
      // Only a successful write reaches this point (a throw skips the eviction pass).
      evictToCap(dir, CACHE_CAP);
    },
  };
}

function storeOf(ctx) {
  const injected = ctx.g4bCacheStore;
  if (injected && typeof injected.get === 'function' && typeof injected.put === 'function') {
    return injected;
  }
  return diskStore(ctx);
}

function run(ctx, out) {
  const whitelist = ctx.loadWhitelist(WHITELIST_REL);
  // Fail-closed FIRST: a broken registry must never be masked by a cache hit.
  if (whitelist.error) throw new UsageError(whitelist.error);
  const white = new Set(whitelist.entries.map((e) => e.match));

  // ---- cross-process corpus cache (slice G4b-CACHE, v1.49) ----
  let cacheKey = null;
  let cacheStore = null;
  if (cacheUsable(ctx)) {
    try {
      cacheKey = fullKey(ctx, whitelist);
      cacheStore = storeOf(ctx);
    } catch (e) {
      cacheKey = null;
      cacheStore = null;
    }
  }

  let corpus = null;
  if (cacheKey) {
    let stored = null;
    try {
      stored = cacheStore.get(cacheKey);
    } catch (e) {
      // A store that throws is not a cache we can trust: go silent (no record at all).
      cacheKey = null;
      cacheStore = null;
    }
    if (cacheKey && validValue(stored)) {
      corpus = new Set(stored.corpus);
      const key12 = cacheKey.slice(0, 12);
      out.rec('G4b-cache', 'INFO', 'G4b cache', 'hit ' + key12 + ' corpus=' + corpus.size);
    }
  }

  if (corpus === null) {
    const ref = corpusRef(ctx);
    corpus = new Set();
    for (const f of ctx.listAllMarkdown()) {
      const text = ctx.showAt(ref, f);
      if (text === null) continue;
      corpus = new Set([...corpus, ...cjkChars(text)]);
    }
    if (cacheKey) {
      try {
        cacheStore.put(cacheKey, { v: 2, corpus: [...corpus] });
      } catch (e) {
        /* fail open: a failed write leaves the full-scan result untouched */
      }
      const key12 = cacheKey.slice(0, 12);
      out.rec('G4b-cache', 'INFO', 'G4b cache', 'miss ' + key12 + ' corpus=' + corpus.size);
    }
  }

  let newC = new Set();
  for (const f of ctx.changedMarkdown()) {
    for (const l of ctx.addedLines(f)) {
      for (const c of cjkChars(l)) newC.add(c);
    }
  }
  const unknown = [...newC].filter((c) => !corpus.has(c) && !white.has(c));
  out.rec(
    'G4b',
    unknown.length === 0 ? 'PASS' : 'FAIL',
    'CJK diff minus whitelist',
    'newCjk=' +
      newC.size +
      ' unknown=' +
      JSON.stringify(unknown) +
      ' (whitelist=' +
      white.size +
      ' whitelist-source=' +
      whitelist.source +
      ')'
  );
}

module.exports = {
  id: 'G4b',
  label: 'CJK diff minus whitelist',
  title: '编码 + 异常字',
  scripted: true,
  run,
  internals: {
    WHITELIST_REL,
    CACHE_VERSION,
    CACHE_CAP,
    cjkChars,
    corpusRef,
    sourceDigest,
    keyString,
    cacheDir,
    evictToCap,
  },
};

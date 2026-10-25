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
 */
const CACHE_VERSION = 'v1';

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
    // Load-bearing (not defensive): range reads ls-tree while tree/index read ls-files.
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
 * `listAllMarkdown()` enumerates `git ls-files` (tree/index) or `git ls-tree`
 * (range/ci), and every file's content is read with `git show <ref>:<p>` - the
 * object store, never the index and never the working tree. Staged additions and
 * staged modifications therefore cannot change the computed corpus, and staged
 * deletions can only SHRINK a live recomputation. No index-cleanliness condition
 * is needed (and none could help): the cached value is the HEAD corpus at write
 * time. If the index later loses a `.md` entry, a cold run recomputes a shrunk
 * corpus while a warm hit still returns the complete base-revision corpus. The
 * warm hit is the fail-safe direction: it holds MORE corpus characters than the
 * shrunk run, hence it reports strictly FEWER unknown characters than a cold run
 * in that window - it never whitewashes a violation present in HEAD - but it is
 * not identical to a cold run in that window.
 */
function cacheUsable(ctx) {
  if (ctx.g4bCacheDisabled) return false;
  if (resolvedSha(ctx) === null) return false;
  return true;
}

/** A stored value must be exactly `{v: 1, corpus: string[]}`; anything else is a miss. */
function validValue(v) {
  return (
    v !== null &&
    typeof v === 'object' &&
    v.v === 1 &&
    Array.isArray(v.corpus) &&
    v.corpus.every((c) => typeof c === 'string')
  );
}

function diskStore(ctx) {
  const dir = cacheDir(ctx);
  // Created eagerly: an unusable directory disables the cache BEFORE any record exists.
  fs.mkdirSync(dir, { recursive: true });
  return {
    get(k) {
      try {
        return JSON.parse(fs.readFileSync(path.join(dir, k + '.json'), 'utf8'));
      } catch (e) {
        return null; // absent / unreadable / malformed: a miss, never an error
      }
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
        cacheStore.put(cacheKey, { v: 1, corpus: [...corpus] });
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
  internals: { WHITELIST_REL, CACHE_VERSION, cjkChars, corpusRef, sourceDigest, keyString, cacheDir },
};

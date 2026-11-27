'use strict';

/**
 * Cross-scope base conclusion cache (DELIVERY_DIRECTIVE 6.1, v1.78, slice P96
 * GATES-BASE-CACHE).
 *
 * Section 6.1's four base gates (gofmt / build / vet / test) read ONLY the working tree,
 * so their conclusions do not depend on the judged `--scope`; the same conclusion can be
 * replayed by another process inside the same key window. The cache stores the four
 * `(verdict, label, detail)` triples - never a G1-G9 criterion conclusion - and only when
 * every one of them is PASS. A FAIL is never frozen, so a flaky red cannot be pinned.
 *
 * Cold path always really runs; the warm path replays the stored records verbatim, which
 * is why a hit and a miss print byte-identical reports (no `hit` / `miss` marker exists).
 *
 * Everything here is fail-open: an uncomputable key, a missing or malformed entry, an
 * unreadable / unwritable store - each degrades to the full run. The cache is an
 * optimisation, never a second source of truth.
 *
 * Key segments (joined with `|`, then sha256):
 *   CACHE_VERSION | HEAD | sha256(porcelain -z) | sha256(sorted untracked path->content)
 *   | sha256(sorted tracked-modified path->content) | `go version` | `go env ...`
 * The key deliberately carries NO scope segment: the four gates only read the workspace.
 *
 * Version bump duty: a change to a key segment's meaning, to the stored schema, or to a
 * base command / detail wording makes every old entry stale - bump CACHE_VERSION.
 */

const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

/** Format version: the first key segment AND the `v` field of a stored value. */
const CACHE_VERSION = 'base-v2';

/** Untracked-file budget: above either bound this run abandons the cache. */
const MAX_UNTRACKED_PATHS = 200;
const MAX_UNTRACKED_BYTES = 64 * 1024 * 1024;

/** LRU bound on the on-disk store (evicted by file mtime, oldest first). */
const LRU_CAP = 50;

/** A cache entry file is exactly the sha256 hex of its key plus `.json`. Atomic temp
 *  files (`.pid.seq.tmp`) never match, so they are never eviction candidates. */
const ENTRY_RE = /^[0-9a-f]{64}\.json$/;

/** Toolchain identity (v1.81): the WHOLE `go env` output is digested instead of
 *  enumerating variables, so GO111MODULE / GOTOOLCHAIN / GOFLAGS and any future variable
 *  are covered without another round of per-variable work. Over-invalidation on an
 *  irrelevant variable is accepted: correctness first, extra cold runs are harmless. */

/** The four §6.1 labels in report order. A stored record whose label is not the one
 *  this position expects is malformed and must never be replayed (a forged PASS entry
 *  would otherwise be indistinguishable from a real one). */
const BASE_LABELS = ['go build ./...', 'go vet ./...', 'go test ./...'];

/** Exact detail shapes of the three cached PASS records; any other string is malformed. */
const BASE_DETAIL_RES = [/^exit 0$/, /^exit 0$/, /^ok=\d+ FAIL=0$/];

/** Per-process sequence for the atomic temp file names. */
let writeSeq = 0;

function sha256(data) {
  return crypto.createHash('sha256').update(data).digest('hex');
}

/** A soft command result: the output string when the call succeeded, else null. */
function softOut(r) {
  try {
    return r && r.ok === true && typeof r.out === 'string' ? r.out : null;
  } catch (e) {
    return null;
  }
}

/**
 * `git status --porcelain=v1 -uall -z` -> the two path buckets that carry a disk body.
 * Deleted paths are dropped (nothing to read) and, in `-z` mode, a rename / copy entry is
 * followed by a second NUL field holding the ORIGINAL path, which must be skipped.
 */
function parsePorcelainZ(out) {
  const untracked = [];
  const tracked = [];
  const fields = String(out).split('\0');
  for (let i = 0; i < fields.length; i += 1) {
    const field = fields[i];
    if (!field || field.length < 4) continue;
    const xy = field.slice(0, 2);
    const p = field.slice(3);
    if (!p) continue;
    if (xy === '??') {
      untracked.push(p);
      continue;
    }
    if (xy[0] === 'R' || xy[0] === 'C' || xy[1] === 'R' || xy[1] === 'C') i += 1;
    if (xy[0] === 'D' || xy[1] === 'D') continue;
    tracked.push(p);
  }
  return { untracked, tracked };
}

/**
 * `[[path, content digest | 'unreadable'], ...]` for the given repository-relative paths.
 * A read failure becomes the literal `unreadable` (it still enters the key, and it never
 * throws). Returns null when the accumulated byte budget is exceeded, which abandons the
 * cache for this run.
 */
function contentDigestMap(repoRoot, relPaths, fsImpl, byteBudget) {
  const entries = [];
  let bytes = 0;
  for (const rel of relPaths.slice().sort()) {
    const abs = path.join(repoRoot, String(rel).split('/').join(path.sep));
    let buf = null;
    try {
      buf = fsImpl.readFileSync(abs);
    } catch (e) {
      entries.push([rel, 'unreadable']);
      continue;
    }
    if (!buf || typeof buf.length !== 'number') {
      entries.push([rel, 'unreadable']);
      continue;
    }
    bytes += buf.length;
    if (byteBudget && bytes > byteBudget) return null;
    let digest = 'unreadable';
    try {
      digest = sha256(buf);
    } catch (e) {
      digest = 'unreadable';
    }
    entries.push([rel, digest]);
  }
  return { entries, bytes };
}

/** One digest over a path->content map, sorted by path so the value is order-free. */
function mapDigest(entries) {
  const lines = entries
    .slice()
    .sort((a, b) => (a[0] < b[0] ? -1 : a[0] > b[0] ? 1 : 0))
    .map((pair) => pair[1] + ' ' + pair[0]);
  return sha256(lines.join('\n'));
}

/**
 * `{key, key12}` for these inputs, or null when the key cannot be trusted (any git / go
 * call failed, the untracked count or byte budget is exceeded, or anything throws).
 *
 * `deps = { repoRoot, git(args)->{ok,out}, execGo(args)->{ok,out}, fsImpl }` is fully
 * injectable so the selftest can drive every key segment without spawning a toolchain.
 */
function computeKey(deps) {
  try {
    const d = deps || {};
    const repoRoot = d.repoRoot;
    if (!repoRoot || typeof repoRoot !== 'string') return null;
    const git = typeof d.git === 'function' ? d.git : () => ({ ok: false, out: '' });
    const execGo = typeof d.execGo === 'function' ? d.execGo : () => ({ ok: false, out: '' });
    const fsImpl = d.fsImpl && typeof d.fsImpl.readFileSync === 'function' ? d.fsImpl : fs;

    const head = softOut(git(['rev-parse', 'HEAD']));
    const status = softOut(git(['status', '--porcelain=v1', '-uall', '-z']));
    const goVersion = softOut(execGo(['version']));
    const goEnv = softOut(execGo(['env']));
    if (head === null || status === null || goVersion === null || goEnv === null) return null;

    const split = parsePorcelainZ(status);
    if (split.untracked.length > MAX_UNTRACKED_PATHS) return null;
    if (split.tracked.length > MAX_UNTRACKED_PATHS) return null;
    // Deletions and renames are dropped by the content maps (no on-disk bytes) but they
    // still change the verdict, so the budget counts every porcelain entry.
    const statusEntries = status.split('\0').filter((seg) => seg.trim() !== '').length;
    if (statusEntries > MAX_UNTRACKED_PATHS) return null;

    const untracked = contentDigestMap(repoRoot, split.untracked, fsImpl, MAX_UNTRACKED_BYTES);
    if (untracked === null) return null;
    const tracked = contentDigestMap(repoRoot, split.tracked, fsImpl, MAX_UNTRACKED_BYTES);
    if (tracked === null) return null;

    // Ignored `.go` files never appear in `git status`, but `gofmt -l .` still walks them,
    // so their content belongs to the key. `-z` keeps non-ASCII paths unquoted (a quoted
    // path would read as `unreadable` and a content change could hide behind the same key).
    // The cache's own directory is excluded: it can never hold Go sources.
    const ignoredGo = softOut(git(['ls-files', '--others', '--ignored', '--exclude-standard', '-z', '--', '*.go']));
    if (ignoredGo === null) return null;
    const ignoredPaths = ignoredGo
      .split('\0')
      .filter((p) => p !== '' && !p.startsWith('tmp/gates-cache/'));
    const ignored = contentDigestMap(repoRoot, ignoredPaths, fsImpl, MAX_UNTRACKED_BYTES);
    if (ignored === null) return null;

    // A non-empty GOWORK means the build resolves modules OUTSIDE this workspace; their
    // sources cannot be keyed, so this pass abandons the cache (fail-open to the full run).
    // `GOWORK=off` and an empty GOWORK both keep the cache.
    const goworkLine = goEnv
      .split('\n')
      .map((l) => l.trim().replace(/^set\s+/, ''))
      .find((l) => l.startsWith('GOWORK='));
    const goworkPath = goworkLine ? goworkLine.slice('GOWORK='.length).trim() : '';
    if (goworkPath && goworkPath !== 'off') return null;

    // `go env` embeds a random `-ffile-prefix-map=<TEMP>/go-build<random>` inside
    // GOGCCFLAGS, which would change the key on every run and disable the cache entirely.
    // That variable carries no build outcome, so it is excluded from the digest.
    const goEnvStable = goEnv
      .split('\n')
      .filter((l) => {
        const t = l.trim().replace(/^set\s+/, '');
        return !t.startsWith('GOGCCFLAGS=');
      })
      .join('\n');

    const segments = [
      CACHE_VERSION,
      head.trim(),
      sha256(status),
      mapDigest(untracked.entries),
      mapDigest(tracked.entries),
      mapDigest(ignored.entries),
      goVersion.trim(),
      sha256(goEnvStable),
    ];
    const key = sha256(segments.join('|'));
    return { key, key12: key.slice(0, 12) };
  } catch (e) {
    return null;
  }
}

/**
 * Strict shape check of a stored value: `{v:1, records:[{id:'base', verdict:'PASS',
 * label, detail} x4]}`. Only PASS is ever stored, so any other verdict is a miss.
 */
function validRecords(v) {
  try {
    if (!v || typeof v !== 'object' || v.v !== 1) return false;
    if (typeof v.key !== 'string' || !v.key) return false;
    const recs = v.records;
    if (!Array.isArray(recs) || recs.length !== 3) return false;
    for (let i = 0; i < recs.length; i++) {
      const r = recs[i];
      if (!r || typeof r !== 'object') return false;
      if (r.id !== 'base') return false;
      if (r.verdict !== 'PASS') return false;
      if (typeof r.label !== 'string' || typeof r.detail !== 'string') return false;
      if (r.label !== BASE_LABELS[i]) return false;
      if (!BASE_DETAIL_RES[i].test(r.detail)) return false;
    }
    return true;
  } catch (e) {
    return false;
  }
}

/**
 * The four stored records for `key`, or null when the entry is absent, unparsable or
 * shaped differently. `store` is `{get(key) -> raw JSON text | null}`.
 */
function read(store, key) {
  try {
    if (!store || typeof store.get !== 'function' || typeof key !== 'string' || !key) return null;
    const text = store.get(key);
    if (typeof text !== 'string' || !text) return null;
    let parsed = null;
    try {
      parsed = JSON.parse(text);
    } catch (e) {
      return null;
    }
    if (!validRecords(parsed)) return null;
    if (parsed.key !== key) return null;
    return parsed.records.map((r) => ({ id: r.id, verdict: r.verdict, label: r.label, detail: r.detail }));
  } catch (e) {
    return null;
  }
}

/**
 * Stores the four records, but ONLY when all four are PASS. `store` is
 * `{put(key, raw JSON text)}`. Every failure is swallowed: a cache that cannot be written
 * must never change the verdict, the detail or the exit code of the run.
 */
function write(store, key, records) {
  try {
    if (!store || typeof store.put !== 'function' || typeof key !== 'string' || !key) return;
    if (!Array.isArray(records) || records.length !== 3) return;
    for (const r of records) {
      if (!r || typeof r !== 'object') return;
      if (r.id !== 'base' || r.verdict !== 'PASS') return;
      if (typeof r.label !== 'string' || typeof r.detail !== 'string') return;
    }
    const value = {
      v: 1,
      key,
      records: records.map((r) => ({ id: 'base', verdict: 'PASS', label: r.label, detail: r.detail })),
    };
    if (!validRecords(value)) return;
    store.put(key, JSON.stringify(value));
  } catch (e) {
    /* fail-open */
  }
}

/** Entry file name for a key: the sha256 hex of the key plus `.json`. */
function entryName(key) {
  return sha256(String(key)) + '.json';
}

/**
 * Keeps at most `cap` entry files in `dir`, evicting the OLDEST mtimes first. Entries whose
 * stat() throws are skipped (never counted), unlink failures are ignored, and an
 * unreadable directory answers 0 - eviction is housekeeping, never a reason to fail a run.
 */
function evictToCap(dir, cap, statFn, unlinkFn) {
  const stat = typeof statFn === 'function' ? statFn : fs.statSync;
  const unlink = typeof unlinkFn === 'function' ? unlinkFn : fs.unlinkSync;
  try {
    const names = fs.readdirSync(dir).filter((n) => ENTRY_RE.test(n));
    const excess = names.length - cap;
    if (!(excess > 0)) return 0;
    const stamped = [];
    for (const n of names) {
      const file = path.join(dir, n);
      try {
        stamped.push({ file, mtime: stat(file).mtimeMs });
      } catch (e) {
        /* unreadable entry: skipped, never counted */
      }
    }
    stamped.sort((a, b) => a.mtime - b.mtime);
    const doomed = stamped.slice(0, Math.min(excess, stamped.length));
    let removed = 0;
    for (const entry of doomed) {
      try {
        unlink(entry.file);
        removed += 1;
      } catch (e) {
        /* fail-open */
      }
    }
    return removed;
  } catch (e) {
    return 0;
  }
}

/**
 * On-disk store under `<repoRoot>/tmp/gates-cache/base/` (`tmp/` is gitignored, so the
 * cache never appears in the porcelain output that feeds its own key).
 *
 * `opts.dir` overrides the directory for the selftest; production always uses the fixed
 * location. Writes are atomic (same-directory temp file + rename) and followed by an LRU
 * eviction pass; a hit touches the mtime so the LRU order follows real use.
 */
function diskStore(repoRoot, opts) {
  const o = opts || {};
  const dir = o.dir || path.join(String(repoRoot || '.'), 'tmp', 'gates-cache', 'base');
  const statFn = o.statFn;
  const unlinkFn = o.unlinkFn;
  return {
    dir,

    get(key) {
      try {
        const file = path.join(dir, entryName(key));
        const text = fs.readFileSync(file, 'utf8');
        try {
          const now = new Date();
          fs.utimesSync(file, now, now);
        } catch (e) {
          // A failed access stamp is a cache anomaly: report a miss so this pass runs the
          // real commands (the clause allows no exception for any cache failure).
          return null;
        }
        return text;
      } catch (e) {
        return null;
      }
    },

    put(key, text) {
      try {
        fs.mkdirSync(dir, { recursive: true });
        const file = path.join(dir, entryName(key));
        writeSeq += 1;
        const tmp = path.join(dir, '.' + process.pid + '.' + writeSeq + '.tmp');
        fs.writeFileSync(tmp, String(text), 'utf8');
        try {
          fs.renameSync(tmp, file);
        } catch (e) {
          try {
            fs.unlinkSync(tmp);
          } catch (e2) {
            /* best effort */
          }
          throw e;
        }
        evictToCap(dir, LRU_CAP, statFn, unlinkFn);
      } catch (e) {
        /* fail-open */
      }
    },
  };
}

module.exports = {
  CACHE_VERSION,
  MAX_UNTRACKED_PATHS,
  MAX_UNTRACKED_BYTES,
  LRU_CAP,
  BASE_LABELS,
  computeKey,
  read,
  write,
  validRecords,
  diskStore,
  internals: { entryName, evictToCap, parsePorcelainZ, mapDigest, contentDigestMap, ENTRY_RE },
};

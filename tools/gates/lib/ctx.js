'use strict';

/**
 * Shared context for every criterion in tools/gates.
 *
 * Responsibilities:
 *   - repository-root discovery (git rev-parse --show-toplevel, .git walk-up fallback;
 *     the migated throwaway script hardcoded D:/LZProjects/prfrail/);
 *   - change-set resolution for --scope=tree|index|range:<A>..<B>|ci, including
 *     untracked files (whose whole body counts as added lines) and staged-only files
 *     (which the old `git diff -U0` silently skipped);
 *   - deletion/rename aware status parsing (the old parser crashed on deleted files);
 *   - normalized reads (BOM strip + CRLF -> LF) and raw byte reads for the encoding gate;
 *   - every git invocation wrapped: a "no match" exit code 1 is data, never a crash.
 *
 * Node >= 18 stdlib + git only. No third-party dependencies.
 */

const fs = require('fs');
const path = require('path');
const cp = require('child_process');

/** Git's well-known empty tree object, used when a scope has no parent commit. */
const EMPTY_TREE = '4b825dc642cb6eb9a060e54bf8d69288fbee4904';
const GIT_MAX_BUFFER = 2e8;

/** Usage / environment error: the runner must exit 2, never report PASS/FAIL. */
class UsageError extends Error {
  constructor(message) {
    super(message);
    this.name = 'UsageError';
  }
}

function toPosix(p) {
  return String(p).replace(/\\/g, '/');
}

function stripBom(text) {
  return text.charCodeAt(0) === 0xfeff ? text.slice(1) : text;
}

/** Same normalization the throwaway script used: strip BOM, CRLF -> LF, keep lone CR. */
function decodeBuffer(buf) {
  return stripBom(buf.toString('utf8')).replace(/\r\n/g, '\n');
}

/**
 * Split the body of an untracked file the way `git --numstat` counts lines: a
 * trailing newline terminates the last line, it does not create an extra empty
 * one. Without this the tree scope reported one more added line per untracked
 * file than the range scope did (GATES-EXT F4), permanently disagreeing with
 * `git --numstat` and with the range detail counts.
 */
function splitAddedLines(text) {
  const lines = String(text).split('\n');
  if (lines.length && lines[lines.length - 1] === '') lines.pop();
  return lines;
}

/**
 * Whitelist source rule (GATES-EXT F1) - ONE rule for BOTH whitelists
 * (`tools/gates/g6-whitelist.txt` used by G6, `tools/gates/cjk-newwords.txt`
 * used by G4b):
 *   tree / index -> the working tree;
 *   range / ci   -> the scope's ENDPOINT revision (`range:A..B` resolves to B,
 *                   `ci` resolves to HEAD).
 * Consequence: a registration made only in the working tree can never exempt a
 * finding in a historical range, and both criteria always read the same snapshot.
 */
function resolveWhitelistSource(scope) {
  const s = scope || {};
  if (s.kind === 'range') return { kind: 'endpoint', ref: s.head };
  if (s.kind === 'ci') return { kind: 'endpoint', ref: 'HEAD' };
  return { kind: 'worktree' };
}

function whitelistSourceLabel(src) {
  return src.kind === 'endpoint' ? 'endpoint:' + src.ref : 'worktree';
}

/**
 * G7-b `repo` disposition (GATES-EXT F2): the token must mean the SAME thing in
 * every scope.
 *   range / ci   -> the path must exist at the scope's endpoint revision
 *                   (`git cat-file -e <ref>:<path>`);
 *   tree / index -> the path must be STAGED in the index
 *                   (`git ls-files --error-unmatch <path>`) AND exist on disk.
 * A gitignored / never-staged path is therefore not enough in tree scope: a
 * tree-scope run is a faithful CI pre-flight only after `git add` of the artifact.
 */
function resolveRepoArtifact(repoRoot, kind, ref, p) {
  if (kind === 'range' || kind === 'ci') {
    const r = spawnGit(repoRoot, ['cat-file', '-e', ref + ':' + p]);
    return {
      ok: r.ok,
      source: 'endpoint:' + ref,
      reason: r.ok ? '' : 'not present at ' + ref + ':' + p,
    };
  }
  const staged = spawnGit(repoRoot, ['ls-files', '--error-unmatch', '--', p]);
  const onDisk = fs.existsSync(path.join(repoRoot, p.split('/').join(path.sep)));
  return {
    ok: staged.ok && onDisk,
    source: 'index+worktree',
    reason: !staged.ok
      ? 'not staged in the index (git add the artifact first)'
      : !onDisk
      ? 'not present in the working tree'
      : '',
  };
}

function spawnGit(cwd, args, binary) {
  const r = cp.spawnSync('git', args, {
    cwd,
    encoding: binary ? 'buffer' : 'utf8',
    maxBuffer: GIT_MAX_BUFFER,
    windowsHide: true,
  });
  return {
    ok: r.status === 0,
    code: r.status === null ? 127 : r.status,
    out: binary ? (r.stdout || Buffer.alloc(0)) : (r.stdout || ''),
    err: binary ? (r.stderr || Buffer.alloc(0)) : (r.stderr || ''),
    spawnError: r.error ? String(r.error.message) : null,
  };
}

function discoverRepoRoot(startDir) {
  const start = path.resolve(startDir || process.cwd());
  const probe = cp.spawnSync('git', ['rev-parse', '--show-toplevel'], {
    cwd: start,
    encoding: 'utf8',
    windowsHide: true,
  });
  if (probe.status === 0 && probe.stdout && probe.stdout.trim()) {
    return path.resolve(probe.stdout.trim());
  }
  let dir = start;
  for (;;) {
    if (fs.existsSync(path.join(dir, '.git'))) return dir;
    const up = path.dirname(dir);
    if (up === dir) return null;
    dir = up;
  }
}

function parseScopeText(text) {
  const s = String(text || 'tree').trim();
  if (s === 'tree' || s === 'index' || s === 'ci') return { kind: s };
  const m = /^range:(.+?)\.\.(.+)$/.exec(s);
  if (m) return { kind: 'range', base: m[1], head: m[2] };
  return null;
}

/** `git status --porcelain=v1 -z -uall` -> records. */
function parseStatusZ(out) {
  const parts = out.split('\0');
  const records = [];
  for (let i = 0; i < parts.length; i++) {
    const rec = parts[i];
    if (!rec) continue;
    const xy = rec.slice(0, 2);
    const p = toPosix(rec.slice(3));
    const entry = {
      path: p,
      status: xy,
      index: xy[0],
      worktree: xy[1],
      untracked: xy === '??',
      staged: xy[0] !== ' ' && xy[0] !== '?',
      unstaged: xy[1] !== ' ' && xy[1] !== '?',
      deleted: xy[0] === 'D' || xy[1] === 'D',
      renamedFrom: null,
    };
    if (xy[0] === 'R' || xy[1] === 'R') {
      entry.renamedFrom = toPosix(parts[++i] || '');
    }
    records.push(entry);
  }
  return records;
}

/** `git diff --name-status -z <base> <head>` -> records. */
function parseNameStatusZ(out) {
  const parts = out.split('\0');
  const records = [];
  for (let i = 0; i < parts.length; i++) {
    const st = parts[i];
    if (!st) continue;
    const code = st[0];
    if (code === 'R' || code === 'C') {
      const src = toPosix(parts[++i] || '');
      const dst = toPosix(parts[++i] || '');
      records.push({
        path: dst,
        status: st,
        index: code,
        worktree: ' ',
        untracked: false,
        staged: true,
        unstaged: false,
        deleted: false,
        renamedFrom: src,
      });
    } else {
      const p = toPosix(parts[++i] || '');
      if (!p) continue;
      records.push({
        path: p,
        status: st,
        index: code,
        worktree: ' ',
        untracked: false,
        staged: true,
        unstaged: false,
        deleted: code === 'D',
        renamedFrom: null,
      });
    }
  }
  return records;
}

/**
 * Added-line records (1-based new-file line number + text) from a unified diff.
 * With -U0 every body line is either an addition or a deletion, so the new-file
 * counter only advances on additions.
 */
function parseHunkAdditions(diffText) {
  const out = [];
  let newLine = 0;
  for (const line of diffText.split('\n')) {
    const h = /^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@/.exec(line);
    if (h) {
      newLine = parseInt(h[1], 10);
      continue;
    }
    if (line.startsWith('+') && !line.startsWith('+++')) {
      out.push({ n: newLine, text: line.slice(1) });
      newLine++;
    }
    // Deleted lines do not advance the new-file counter; headers are ignored.
  }
  return out;
}

function createContext(options) {
  const opts = options || {};
  const scopeText = opts.scope || 'tree';
  const scope = parseScopeText(scopeText);
  if (!scope) throw new Error('invalid --scope value: ' + scopeText);

  const repoRoot = opts.repoRoot || discoverRepoRoot(opts.cwd || process.cwd());
  if (!repoRoot) {
    throw new Error('cannot locate the repository root (git rev-parse --show-toplevel failed and no .git found)');
  }

  const kind = scope.kind;
  let baseRef = null;
  let headRef = null;

  // Test seam (GATES-EXT, 7 round-2): the scope guards below must be verifiable one by one, so the
  // git runner can be injected. Production always uses spawnGit; the selftest substitutes a fake to
  // prove each guard fires independently instead of only as a bundle.
  const spawn = opts.spawnGit || spawnGit;

  function git(args, o) {
    return spawn(repoRoot, args, !!(o && o.binary));
  }

  if (kind === 'ci') {
    if (!git(['rev-parse', '--verify', 'HEAD']).ok) {
      throw new UsageError('--scope=ci requires a HEAD commit (root-commit guard: use --scope=range instead)');
    }
    headRef = 'HEAD';
    if (git(['rev-parse', '--verify', 'HEAD^1']).ok) {
      baseRef = 'HEAD^1';
    } else {
      // A missing parent is legitimate ONLY for a genuine root commit. On a SHALLOW clone it means
      // the history was never fetched, and the EMPTY_TREE fallback would classify the WHOLE
      // repository as added lines - a whole-tree false-positive storm wearing the costume of a real
      // run (GATES-EXT, 7 F1/F4). The probe must be CONCLUSIVE: an inconclusive probe is itself an
      // environment error, never a licence to fall back.
      const probe = git(['rev-parse', '--is-shallow-repository']);
      const answer = (probe.out || '').trim();
      if (!probe.ok || (answer !== 'true' && answer !== 'false')) {
        throw new UsageError(
          '--scope=ci: cannot decide whether HEAD^1 is missing because this is a root commit or because ' +
            'the clone is shallow (`git rev-parse --is-shallow-repository` answered ' +
            (probe.ok ? JSON.stringify(answer) : 'exit ' + probe.code) +
            '). Use --scope=range:<A>..<B> explicitly.'
        );
      }
      if (answer === 'true') {
        throw new UsageError(
          '--scope=ci: HEAD^1 is unavailable because this clone is SHALLOW; a whole-tree change set ' +
            'would follow. Fetch the full history (actions/checkout `fetch-depth: 0`) or use --scope=range:<A>..<B>.'
        );
      }
      baseRef = EMPTY_TREE;
    }
  } else if (kind === 'range') {
    baseRef = scope.base;
    headRef = scope.head;
    // An unresolvable endpoint must be a usage/environment error (exit 2), never an empty change
    // set: with an empty set every criterion prints a vacuous pass and the gate goes green while
    // judging nothing - a fail-open in a CI-blocking gate (GATES-EXT, 7 F1). The `^{commit}`
    // suffix is load-bearing: a bare `rev-parse --verify` accepts any well-formed 40-hex string
    // (including the null oid) without checking that the object exists.
    for (const pair of [['base', baseRef], ['head', headRef]]) {
      if (!git(['rev-parse', '--verify', '--quiet', pair[1] + '^{commit}']).ok) {
        throw new UsageError(
          '--scope=range: the ' +
            pair[0] +
            ' revision is not resolvable to a commit: ' +
            pair[1] +
            ' (an unresolvable range must not degrade into an empty change set)'
        );
      }
    }
  }

  let changed = [];
  if (kind === 'tree' || kind === 'index') {
    // Same rule as the range/ci diff below: a failed command is an ENVIRONMENT ERROR, never
    // "no changes" (7 round-2 F3). A corrupt/unreadable index makes `git status` fail with empty
    // stdout; reading that as an empty change set turned the whole run into vacuous passes.
    const sr = git(['status', '--porcelain=v1', '-uall', '-z']);
    if (!sr.ok) {
      throw new UsageError(
        'environment error: git status failed for scope ' + kind + ': ' + ((sr.err || '').trim() || 'exit ' + sr.code)
      );
    }
    changed = parseStatusZ(sr.out);
    if (kind === 'index') changed = changed.filter((rec) => rec.staged);
  } else {
    // range / ci: the change set comes from `git diff`. A failed diff (bad revision, unavailable
    // object, unreadable repository) must surface as an environment error - reading its empty
    // stdout as "no changes" is exactly the fail-open path guarded above.
    const r = git(['diff', '--name-status', '-z', baseRef, headRef]);
    if (!r.ok) {
      throw new UsageError(
        'environment error: git diff --name-status failed for scope ' +
          kind +
          ' (' +
          baseRef +
          '..' +
          headRef +
          '): ' +
          ((r.err || '').trim() || 'exit ' + r.code)
      );
    }
    changed = parseNameStatusZ(r.out);
  }

  const byPath = new Map();
  for (const rec of changed) byPath.set(rec.path, rec);

  function diffPrefixes() {
    if (kind === 'index') return [['diff', '--cached']];
    if (kind === 'range' || kind === 'ci') return [['diff', baseRef, headRef]];
    return [['diff'], ['diff', '--cached']];
  }

  function pathOnDisk(p) {
    return path.join(repoRoot, p.split('/').join(path.sep));
  }

  function readBytes(p) {
    if (kind === 'range' || kind === 'ci') {
      const r = git(['show', headRef + ':' + p], { binary: true });
      return r.ok ? r.out : null;
    }
    try {
      return fs.readFileSync(pathOnDisk(p));
    } catch (e) {
      return null;
    }
  }

  function readLines(p) {
    const buf = readBytes(p);
    return buf === null ? null : decodeBuffer(buf);
  }

  function readText(p) {
    const buf = readBytes(p);
    return buf === null ? null : decodeBuffer(buf);
  }

  function readLineArray(p) {
    const text = readText(p);
    return text === null ? null : text.split('\n');
  }

  function exists(p) {
    if (kind === 'range' || kind === 'ci') {
      return git(['cat-file', '-e', headRef + ':' + p]).ok;
    }
    return fs.existsSync(pathOnDisk(p));
  }

  function trackedAt(ref, p) {
    return git(['cat-file', '-e', ref + ':' + p]).ok;
  }

  function showAt(ref, p) {
    const r = git(['show', ref + ':' + p]);
    return r.ok ? r.out : null;
  }

  function isIgnored(p) {
    return git(['check-ignore', '-q', '--', p]).ok;
  }

  const hunkCache = new Map();

  function addedHunks(p) {
    if (hunkCache.has(p)) return hunkCache.get(p);
    const rec = byPath.get(p);
    let out = [];
    if (rec && rec.untracked) {
      const text = readText(p);
      if (text !== null) out = splitAddedLines(text).map((t, i) => ({ n: i + 1, text: t }));
    } else {
      for (const prefix of diffPrefixes()) {
        const r = git(prefix.concat(['-U0', '--no-color', '--', p]));
        if (r.ok) out = out.concat(parseHunkAdditions(r.out));
      }
    }
    hunkCache.set(p, out);
    return out;
  }

  function addedLines(p) {
    return addedHunks(p).map((h) => h.text);
  }

  function addedLineNumbers(p) {
    return new Set(addedHunks(p).map((h) => h.n));
  }

  function numstat(p) {
    const rec = byPath.get(p);
    if (rec && rec.untracked) {
      const text = readText(p);
      return [String(text === null ? 0 : splitAddedLines(text).length), '0'];
    }
    for (const prefix of diffPrefixes()) {
      const r = git(prefix.concat(['--numstat', '--', p]));
      if (!r.ok) continue;
      const line = r.out.trim().split('\n')[0];
      if (!line) continue;
      const parts = line.split('\t');
      if (parts.length >= 2) return [parts[0], parts[1]];
    }
    return ['0', '0'];
  }

  function listAllMarkdown() {
    // NOTE: do NOT pass a `*.md` pathspec here. `git ls-tree -r --name-only <ref> -- '*.md'`
    // returns an empty list (measured 2026-09-22), which silently empties the G4b corpus and
    // turns the criterion into a false-positive machine. List everything and filter in JS.
    const r =
      kind === 'range' || kind === 'ci'
        ? git(['ls-tree', '-r', '--name-only', baseRef])
        : git(['ls-files']);
    return r.ok ? r.out.split('\n').filter((p) => p && /\.md$/i.test(p)) : [];
  }

  function listAgentFiles() {
    const dir = path.join(repoRoot, '.github', 'agents');
    try {
      return fs
        .readdirSync(dir)
        .filter((f) => f.endsWith('.agent.md'))
        .sort();
    } catch (e) {
      return [];
    }
  }

  /**
   * Whitelist loader (one function for both whitelists, see resolveWhitelistSource).
   * Returns the used source label, the raw text (null when absent) and the parsed
   * `<match>\t<file>\t<reason>` entries; `error` is a usage error, never silently
   * ignored.
   */
  function loadWhitelist(rel, scopeOverride) {
    const src = resolveWhitelistSource(scopeOverride || scope);
    const source = whitelistSourceLabel(src);
    let text = null;
    let error = null;
    if (src.kind === 'endpoint') {
      const r = git(['show', src.ref + ':' + rel]);
      if (r.ok) text = decodeBuffer(r.out);
    } else {
      const abs = pathOnDisk(rel);
      if (fs.existsSync(abs)) {
        try {
          text = decodeBuffer(fs.readFileSync(abs));
        } catch (e) {
          error = rel + ': unreadable (' + e.message + ')';
        }
      }
    }
    if (text === null) return { source, text: null, entries: [], error };
    const parsed = parseWhitelist(rel, text);
    return { source, text, entries: parsed.entries, error: parsed.error };
  }

  const changedPaths = changed.map((r) => r.path);

  return {
    repoRoot,
    scope,
    baseRef,
    headRef,
    kind,
    changed,
    changedPaths,
    recordFor: (p) => byPath.get(p) || null,
    changedMarkdown: () => changedPaths.filter((p) => /\.md$/i.test(p)),
    changedPathsWith: (re) => changedPaths.filter((p) => re.test(p)),
    git,
    readBytes,
    readLines,
    readText,
    readLineArray,
    exists,
    repoArtifactResolvable: (p) => resolveRepoArtifact(repoRoot, kind, headRef, p),
    trackedAt,
    showAt,
    isIgnored,
    loadWhitelist,
    UsageError,
    addedHunks,
    addedLines,
    addedLineNumbers,
    numstat,
    listAllMarkdown,
    listAgentFiles,
  };
}

/**
 * Shared whitelist parser: `<match>\t<file>\t<reason / context excerpt>`.
 * Every field is mandatory: a missing reason/excerpt is a usage error (exit 2),
 * which is what keeps a whitelist from becoming a silent blanket. Both
 * `g6-whitelist.txt` and `cjk-newwords.txt` use this format (G4b registers one
 * character in the match field).
 */
function parseWhitelist(relPath, text) {
  const entries = [];
  const lines = String(text).split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (!line.trim() || line.startsWith('#')) continue;
    const parts = line.split('\t');
    if (parts.length < 3 || !parts[0].trim() || !parts[1].trim() || !parts[2].trim()) {
      return {
        entries,
        error:
          relPath +
          ':' +
          (i + 1) +
          ' malformed entry (expected <match><TAB><file><TAB><reason/context excerpt>)',
      };
    }
    entries.push({ match: parts[0], file: toPosix(parts[1].trim()), reason: parts[2].trim() });
  }
  return { entries, error: null };
}

/** Working-tree-only loader (kept for callers that have no scope); see createContext. */
function loadWhitelist(repoRoot, relPath) {
  const abs = path.join(repoRoot, relPath.split('/').join(path.sep));
  if (!fs.existsSync(abs)) return { entries: [], error: null };
  let text;
  try {
    text = decodeBuffer(fs.readFileSync(abs));
  } catch (e) {
    return { entries: [], error: relPath + ': unreadable (' + e.message + ')' };
  }
  return parseWhitelist(relPath, text);
}

/** True when the finding is covered by a registered whitelist entry. */
function whitelisted(entries, file, text) {
  return entries.filter((e) => (e.file === '*' || e.file === file) && String(text).includes(e.match));
}

module.exports = {
  EMPTY_TREE,
  UsageError,
  createContext,
  discoverRepoRoot,
  parseScopeText,
  parseStatusZ,
  parseNameStatusZ,
  parseHunkAdditions,
  splitAddedLines,
  resolveWhitelistSource,
  whitelistSourceLabel,
  resolveRepoArtifact,
  parseWhitelist,
  loadWhitelist,
  whitelisted,
  toPosix,
  stripBom,
  decodeBuffer,
};

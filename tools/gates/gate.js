#!/usr/bin/env node
'use strict';

/**
 * DELIVERY_DIRECTIVE section 6.3 gate runner - single entry point.
 *
 *   node tools/gates/gate.js --all [--scope=tree|index|range:<A>..<B>|ci]
 *   node tools/gates/gate.js --check=<ID> [--scope=...]
 *   node tools/gates/gate.js --list
 *   node tools/gates/gate.js --selftest
 *
 * Exit codes: 0 = every executed criterion passed (SKIP/INFO never block),
 *             1 = at least one FAIL,
 *             2 = usage / environment error (bad flag, bad scope, missing git or Go
 *                 toolchain, malformed whitelist entry).
 *
 * The criteria modules stay individually loadable under lib/criteria/; this file only
 * parses arguments, assembles the report and owns the exit code (the throwaway
 * tmp/gate.js always exited 0, which made it unusable as a gate).
 */

const cp = require('child_process');

const { createContext, UsageError, parseScopeText } = require('./lib/ctx');
const { createReport } = require('./lib/report');

const CRITERIA = {
  'G1-a': require('./lib/criteria/g1a'),
  'G1-b': require('./lib/criteria/g1b'),
  G2: require('./lib/criteria/g2'),
  'G3-a': require('./lib/criteria/g3a'),
  'G3-b': require('./lib/criteria/g3b'),
  G4a: require('./lib/criteria/g4a'),
  G4b: require('./lib/criteria/g4b'),
  'G5-a': require('./lib/criteria/g5a'),
  'G5-b': require('./lib/criteria/g5b'),
  'R2.5': require('./lib/criteria/r25'),
  G6: require('./lib/criteria/g6'),
  G7: require('./lib/criteria/g7'),
};

/**
 * Base criteria deliberately NOT scripted in slice GATES-EXT. They are reported as
 * SKIP, never as PASS, so the report cannot be read as "everything was checked".
 */
const SKIPPED = {
  '@skipG1': {
    id: 'G1',
    label: 'section-title/status-table/completion-record consistency',
    reason:
      'not scripted in slice GATES-EXT: judging three human-authored places against each other; see DELIVERY_DIRECTIVE 6.3 G1',
  },
  '@skipG5': {
    id: 'G5',
    label: 'change-set contract (staged vs declared vs evidence refs)',
    reason:
      'not scripted in slice GATES-EXT: requires a structured JSON input produced outside the repository; see DELIVERY_DIRECTIVE 6.3 G5',
  },
};

const RUN_ORDER = [
  'G1-a',
  'G1-b',
  '@skipG1',
  'G2',
  'G3-a',
  'G3-b',
  'G4a',
  'G4b',
  'G5-a',
  'G5-b',
  '@skipG5',
  'R2.5',
  'G6',
  'G7',
];

const USAGE =
  'usage: node tools/gates/gate.js --all|--check=<ID>|--list|--selftest [--scope=tree|index|range:<A>..<B>|ci]';

/**
 * Criterion ids are compared after folding every non-alphanumeric separator
 * (`G3(a.1)` -> `G3-A-1`), so `--check=G3-a.1` and `--check=G3(a.1)` select the same
 * sub-check instead of exiting 2 (GATES-EXT F5).
 */
function normalizeId(s) {
  return String(s)
    .replace(/[.[\]()]/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
    .toUpperCase();
}

function emitSkip(out, key) {
  const s = SKIPPED[key];
  out.rec(s.id, 'SKIP', s.label, s.reason);
}

/** The four base gates (section 6.1): gofmt, build, vet, test. */
function runBaseGates(ctx, out) {
  const exec = (cmd, args) =>
    cp.spawnSync(cmd, args, {
      cwd: ctx.repoRoot,
      encoding: 'utf8',
      maxBuffer: 2e8,
      windowsHide: true,
      shell: false,
    });

  const gf = exec('gofmt', ['-l', '.']);
  if (gf.error) throw new UsageError('environment error: cannot run gofmt (' + gf.error.message + ')');
  const gfOut = (gf.stdout || '').trim();
  out.rec('base', gf.status === 0 && gfOut === '' ? 'PASS' : 'FAIL', 'gofmt -l .', gfOut === '' ? 'empty' : gfOut);

  for (const [step, verb] of [['build', 'build'], ['vet', 'vet']]) {
    const r = exec('go', [verb, './...']);
    if (r.error) throw new UsageError('environment error: cannot run go (' + r.error.message + ')');
    const tail = ((r.stderr || '') + (r.stdout || '')).trim().split('\n').slice(-3).join(' / ');
    out.rec(
      'base',
      r.status === 0 ? 'PASS' : 'FAIL',
      'go ' + step + ' ./...',
      r.status === 0 ? 'exit 0' : 'exit ' + r.status + (tail ? ' :: ' + tail : '')
    );
  }

  const t = exec('go', ['test', './...']);
  if (t.error) throw new UsageError('environment error: cannot run go (' + t.error.message + ')');
  const o = ((t.stdout || '').match(/^ok\s/gm) || []).length;
  const fa = ((t.stdout || '').match(/^FAIL/gm) || []).length;
  const tail = (t.stderr || '').trim().split('\n').slice(-3).join(' / ');
  out.rec(
    'base',
    t.status === 0 && fa === 0 ? 'PASS' : 'FAIL',
    'go test ./...',
    'ok=' + o + ' FAIL=' + fa + (t.status === 0 ? '' : ' exit=' + t.status + (tail ? ' :: ' + tail : ''))
  );
}

function selectCriteria(check) {
  const want = normalizeId(check);
  const selected = [];
  for (const [key, mod] of Object.entries(CRITERIA)) {
    const mid = normalizeId(key);
    if (mid === want) selected.push({ mod, filter: null });
    else if (want.startsWith(mid + '-')) selected.push({ mod, filter: want });
  }
  return selected;
}

function listCriteria() {
  const lines = ['criteria (id | label | DELIVERY_DIRECTIVE 6.3 label | status)'];
  for (const key of Object.keys(CRITERIA)) {
    const m = CRITERIA[key];
    lines.push(
      '  ' + key + ' | ' + m.label + ' | ' + (m.title || '') + ' | ' + (m.scripted ? 'scripted' : 'not scripted')
    );
  }
  for (const key of Object.keys(SKIPPED)) {
    lines.push('  ' + SKIPPED[key].id + ' | ' + SKIPPED[key].label + ' | (base criterion) | NOT scripted (SKIP)');
  }
  lines.push('  base | gofmt -l . / go build ./... / go vet ./... / go test ./... | section 6.1 | scripted');
  return lines.join('\n');
}

function main(argv) {
  const args = argv.slice(2);
  const opt = { mode: null, check: null, scope: 'tree' };
  for (const a of args) {
    if (a === '--all') opt.mode = 'all';
    else if (a === '--selftest') opt.mode = 'selftest';
    else if (a === '--list') opt.mode = 'list';
    else if (a === '--help' || a === '-h') opt.mode = 'help';
    else if (a.startsWith('--check=')) {
      opt.mode = 'check';
      opt.check = a.slice('--check='.length);
    } else if (a.startsWith('--scope=')) opt.scope = a.slice('--scope='.length);
    else return { code: 2, error: 'unknown argument: ' + a + '\n' + USAGE };
  }
  if (!opt.mode) return { code: 2, error: USAGE };
  if (opt.mode === 'help') return { code: 0, stdout: USAGE };
  if (!parseScopeText(opt.scope)) {
    return { code: 2, error: 'invalid --scope value: ' + opt.scope + '\n' + USAGE };
  }
  if (opt.mode === 'list') return { code: 0, stdout: listCriteria() };
  if (opt.mode === 'selftest') {
    const res = require('./selftest').run();
    return { code: res.ok ? 0 : 1, stdout: res.stdout };
  }

  const ctx = createContext({ scope: opt.scope });
  const out = createReport();
  const header =
    'GATE REPORT scope=' +
    opt.scope +
    ' repo=' +
    ctx.repoRoot +
    ' changed=' +
    ctx.changedPaths.length +
    (opt.mode === 'check' ? ' check=' + opt.check : '');
  if (ctx.changedPaths.length === 0) {
    // A false gate: every criterion passes because it had nothing to look at. The CI step
    // must use --scope=ci precisely to avoid this (a clean checkout has no tree changes).
    out.rec(
      'scope',
      'INFO',
      'empty change set',
      'nothing to judge: criteria below report vacuous passes (deliberately NOT reported as FAIL)'
    );
  }

  const plan = opt.mode === 'all' ? RUN_ORDER : null;
  if (plan) {
    for (const key of plan) {
      if (key.startsWith('@')) {
        emitSkip(out, key);
        continue;
      }
      CRITERIA[key].run(ctx, out);
    }
    runBaseGates(ctx, out);
  } else {
    const selected = selectCriteria(opt.check);
    if (!selected.length) {
      const key = normalizeId(opt.check);
      if (key === 'G1' || key === 'G5') {
        emitSkip(out, key === 'G1' ? '@skipG1' : '@skipG5');
      } else {
        return { code: 2, error: 'unknown criterion id: ' + opt.check + '\n' + listCriteria() };
      }
    } else {
      for (const s of selected) {
        const tmp = createReport();
        s.mod.run(ctx, tmp);
        for (const r of tmp.records) {
          // Prefix match: a sub-check id may carry a file suffix (G3(a.1) <file>), and the
          // filter must still select it (otherwise `--check=G3-a.1` would silently drop the
          // very FAIL record it exists to report).
          if (s.filter && normalizeId(r.id).indexOf(s.filter) !== 0) continue;
          out.rec(r.id, r.verdict, r.label, r.detail);
        }
      }
    }
  }

  const body = out.render();
  return { code: out.fails().length ? 1 : 0, stdout: header + '\n' + body };
}

function run() {
  let result;
  try {
    result = main(process.argv);
  } catch (e) {
    if (e instanceof UsageError) return { code: 2, stderr: 'gate.js: ' + e.message };
    return { code: 2, stderr: 'gate.js: unexpected error: ' + (e && e.stack ? e.stack : e) };
  }
  // `main` reports usage / environment errors on the `error` key (bad flag, missing mode,
  // bad --scope, unknown criterion id). Batch 2f (NF-3): that key was never consumed, so
  // `node tools/gates/gate.js --scope=tree` exited 2 having printed NOTHING - the caller
  // saw a bare exit code with no diagnosis. The message now goes to stderr; the exit code
  // stays 2 (a usage error is still not a PASS/FAIL verdict).
  const err = result.stderr || result.error;
  if (err) process.stderr.write(err + '\n');
  if (result.stdout) process.stdout.write(result.stdout + '\n');
  return result;
}

if (require.main === module) {
  process.exitCode = run().code;
}

module.exports = { run, main, CRITERIA, SKIPPED, RUN_ORDER, listCriteria };

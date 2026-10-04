#!/usr/bin/env node
'use strict';
/*
 * Evidence capture for slice DIRECTIVE-REVIEW-SATURATION.
 *
 * Why this exists: capturing a native command's stdout through a PowerShell
 * pipeline (`... | Out-File`, and in some hosts even `>`) decodes the UTF-8 bytes
 * with the console code page (GBK here) and re-encodes them, so every Chinese
 * character in the evidence file becomes mojibake. Node writes the bytes verbatim,
 * so all captures below go through execFileSync + fs.writeFileSync('utf8') and are
 * byte-verified afterwards (see tmp/v113/verify-encoding.js).
 *
 * Usage: node tmp/v113/capture.js [--with-remediation]
 */
const { execFileSync } = require('child_process');
const fs = require('fs');

const run = (args) => execFileSync(args[0], args.slice(1), { encoding: 'utf8', maxBuffer: 1 << 28 });
// `--suffix=-r2` keeps each round's readings on disk instead of overwriting the
// previous round's evidence (the same discipline as the v1.12 freeze round).
const SUFFIX = (process.argv.find((a) => a.startsWith('--suffix=')) || '').replace('--suffix=', '');
const write = (file, text) => {
  const out = file.replace(/\.([a-z]+)$/, SUFFIX + '.$1');
  fs.writeFileSync(out, text);
  console.log('[capture] ' + out + ' bytes=' + Buffer.byteLength(text, 'utf8'));
};

write('tmp/v113/gate-tree-r1.txt', run(['node', 'tools/gates/gate.js', '--all', '--scope=tree']));
write('tmp/v113/gate-index-r1.txt', run(['node', 'tools/gates/gate.js', '--all', '--scope=index']));
write('tmp/v113/selftest-master.txt', run(['node', 'tools/gates/gate.js', '--selftest']));
write('tmp/v113/g6-5-r1.txt', run(['node', 'tools/gates/gate.js', '--check=G6-5', '--scope=tree']));
write('tmp/v113/g6-6-r1.txt', run(['node', 'tools/gates/gate.js', '--check=G6-6', '--scope=tree']));
write(
  'tmp/v113/mirror.txt',
  run([
    'node',
    'tmp/v112/mirror-check.js',
    '.',
    'docs/DELIVERY_DIRECTIVE.md=docs/DELIVERY_DIRECTIVE_EN.md',
  ])
);
write('tmp/v113/review-input-r2.diff', run(['git', 'diff', 'HEAD']));
write('tmp/v113/review-input-r2.stat.txt', run(['git', 'diff', 'HEAD', '--numstat']));
write('tmp/v113/verify-v113.txt', run(['node', 'tmp/v113/verify-v113.js']));

if (process.argv.includes('--with-remediation')) {
  const FILES = [
    'docs/DELIVERY_DIRECTIVE.md',
    'docs/DELIVERY_DIRECTIVE_EN.md',
    'tools/gates/lib/criteria/g6.js',
    'tools/gates/selftest.js',
  ];
  let out = '';
  for (const f of FILES) {
    const pre = 'tmp/v113/wt/' + f;
    if (!fs.existsSync(pre)) {
      out += '# pre-remediation copy missing for ' + f + ' (skipped)\n';
      continue;
    }
    let d = '';
    try {
      d = execFileSync('git', ['diff', '--no-index', '--', pre, f], { encoding: 'utf8', maxBuffer: 1 << 28 });
    } catch (e) {
      d = e.stdout || '';
    }
    out += '# pre-remediation -> current: ' + f + '\n';
    out += d;
  }
  write('tmp/v113/remediation.diff', out);
}

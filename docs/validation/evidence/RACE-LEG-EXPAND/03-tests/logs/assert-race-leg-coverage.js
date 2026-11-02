#!/usr/bin/env node
/**
 * RACE-LEG-EXPAND / 03-tests TDD assertion (red-baseline first, then red -> green flip).
 *
 * Purpose
 *   Prove, mechanically, whether the CI "Race" step covers all six packages that
 *   P49 `RACE-LEG-EXPAND` requires. The assertion is written BEFORE ci.yml is changed,
 *   so it must fail (RED) on the current tree and pass (GREEN) after the implementer
 *   extends the `run:` line.
 *
 * Contract asserted
 *   1. A step named "Race" exists in .github/workflows/ci.yml.
 *   2. Its `run:` command still starts with the CGO_ENABLED=1 prefix.
 *   3. Its `run:` command still uses `go test -race`.
 *   4. Its `run:` command still uses `-count=1`.
 *   5. Its `run:` command names exactly the six expected package patterns, in order:
 *        ./internal/adapters/... ./internal/chain/... ./internal/console/...
 *        ./internal/gates/...    ./internal/guard/... ./internal/tickets/...
 *   6. The step keeps its `if: runner.os == 'Linux'` guard.
 *
 * Exit codes
 *   0  GREEN  - every check passed, all six packages covered.
 *   1  RED    - at least one check failed (expected baseline state today).
 *   2  ERROR  - the CI file or the Race step could not be read.
 *
 * Notes
 *   - Read-only with respect to the repository: this script never writes files.
 *   - UTF-8 (no BOM) + LF.
 */

'use strict';

const fs = require('fs');
const path = require('path');

const CI_PATH = path.resolve(__dirname, '..', '..', '.github', 'workflows', 'ci.yml');

const EXPECTED_PACKAGES = [
  './internal/adapters/...',
  './internal/chain/...',
  './internal/console/...',
  './internal/gates/...',
  './internal/guard/...',
  './internal/tickets/...',
];

const LEGACY_PACKAGES = [
  './internal/adapters/...',
  './internal/chain/...',
];

const NEW_PACKAGES = [
  './internal/console/...',
  './internal/gates/...',
  './internal/guard/...',
  './internal/tickets/...',
];

function readLines(filePath) {
  const raw = fs.readFileSync(filePath, 'utf8');
  // Tolerate CRLF input (repository convention is LF, but never crash on CRLF).
  return raw.replace(/\r\n/g, '\n').split('\n');
}

/**
 * Locate the "Race" step and return its line span plus the flattened `run` command.
 * The step span ends at the next list item ("- ...") whose indentation is less than
 * or equal to the Race step's own indentation, or at end of file.
 */
function locateRaceStep(lines) {
  const stepHeader = /^(\s*)-\s+name:\s*Race\s*$/;

  let startIndex = -1;
  let stepIndent = -1;

  for (let i = 0; i < lines.length; i += 1) {
    const match = stepHeader.exec(lines[i]);
    if (match) {
      startIndex = i;
      stepIndent = match[1].length;
      break;
    }
  }

  if (startIndex === -1) {
    return null;
  }

  let endIndex = lines.length;
  for (let i = startIndex + 1; i < lines.length; i += 1) {
    const sibling = /^(\s*)-\s+/.exec(lines[i]);
    if (sibling && sibling[1].length <= stepIndent) {
      endIndex = i;
      break;
    }
  }

  const stepLines = lines.slice(startIndex, endIndex);

  let ifCondition = null;
  let runCommand = '';
  let runDeclaredAt = -1;

  for (let i = 0; i < stepLines.length; i += 1) {
    const ifMatch = /^\s*if:\s*(.+)$/.exec(stepLines[i]);
    if (ifMatch && ifCondition === null) {
      ifCondition = ifMatch[1].trim();
    }

    const runMatch = /^(\s*)run:\s*(.*)$/.exec(stepLines[i]);
    if (runMatch && runDeclaredAt === -1) {
      runDeclaredAt = startIndex + i + 1; // 1-based line number in the file
      const inlineValue = runMatch[2].trim();
      if (inlineValue === '|' || inlineValue === '>' || inlineValue === '|-' || inlineValue === '>-') {
        // Block scalar: collect the following, deeper-indented lines verbatim.
        const runIndent = runMatch[1].length;
        const block = [];
        for (let j = i + 1; j < stepLines.length; j += 1) {
          const line = stepLines[j];
          const indent = line.length - line.trimStart().length;
          if (line.trim() === '') {
            block.push('');
            continue;
          }
          if (indent <= runIndent) {
            break;
          }
          block.push(line.trim());
        }
        runCommand = block.join(' ').trim();
      } else {
        runCommand = inlineValue;
      }
    }
  }

  return {
    startLine: startIndex + 1,
    endLine: endIndex,
    ifCondition,
    runCommand,
    runDeclaredAt,
  };
}

function report(checkName, ok, detail) {
  const status = ok ? 'PASS' : 'FAIL';
  console.log(`[CHECK] ${checkName}: ${status}${detail ? ' - ' + detail : ''}`);
  return ok;
}

function occurrences(haystack, needle) {
  if (!needle) {
    return 0;
  }
  let count = 0;
  let from = 0;
  for (;;) {
    const at = haystack.indexOf(needle, from);
    if (at === -1) {
      return count;
    }
    count += 1;
    from = at + needle.length;
  }
}

function main() {
  console.log('RACE-LEG-EXPAND TDD ASSERTION (03-tests / red-baseline)');
  console.log(`CI file: ${CI_PATH}`);
  console.log('Expected packages (in order):');
  EXPECTED_PACKAGES.forEach((pkg, i) => console.log(`  ${i + 1}. ${pkg}`));
  console.log('');

  if (!fs.existsSync(CI_PATH)) {
    console.log(`RACE-COVERAGE: ERROR - CI file not found: ${CI_PATH}`);
    return 2;
  }

  const lines = readLines(CI_PATH);
  const race = locateRaceStep(lines);

  if (!race) {
    console.log("RACE-COVERAGE: RED - no step named 'Race' was found in ci.yml");
    return 1;
  }

  console.log(`Race step found at ci.yml line ${race.startLine} (step span to line ${race.endLine})`);
  console.log(`Race step if   : ${race.ifCondition === null ? '(none)' : race.ifCondition}`);
  console.log(`Race step run  : ${race.runCommand === '' ? '(empty)' : race.runCommand}`);
  console.log('');

  const command = race.runCommand;
  const checks = [];

  const missingPackages = EXPECTED_PACKAGES.filter((pkg) => occurrences(command, pkg) === 0);
  const missingLegacy = LEGACY_PACKAGES.filter((pkg) => occurrences(command, pkg) === 0);
  const missingNew = NEW_PACKAGES.filter((pkg) => occurrences(command, pkg) === 0);

  checks.push(
    report(
      'CGO_ENABLED=1 prefix',
      command.startsWith('CGO_ENABLED=1 '),
      command.startsWith('CGO_ENABLED=1 ') ? '"CGO_ENABLED=1 " present at command head' : 'prefix missing or not at head'
    )
  );

  checks.push(report('go test -race', command.includes('go test -race'), ''));
  checks.push(report('-count=1', command.includes('-count=1'), ''));

  checks.push(
    report(
      'Linux-only guard preserved',
      race.ifCondition !== null && race.ifCondition.includes("runner.os == 'Linux'"),
      race.ifCondition === null ? 'no if: condition on the Race step' : ''
    )
  );

  EXPECTED_PACKAGES.forEach((pkg) => {
    checks.push(report(`package covered: ${pkg}`, occurrences(command, pkg) === 1, `occurrences=${occurrences(command, pkg)}`));
  });

  // Order check: package patterns must appear in the documented order.
  const positions = EXPECTED_PACKAGES.map((pkg) => command.indexOf(pkg));
  const orderOk = positions.every((pos, i) => pos !== -1 && (i === 0 || pos > positions[i - 1]));
  checks.push(report('package order matches documented order', orderOk, orderOk ? '' : 'order differs or a package is absent'));

  console.log('');
  console.log(`legacy leg (adapters/chain) : ${missingLegacy.length === 0 ? 'present' : 'MISSING -> ' + missingLegacy.join(' ')}`);
  console.log(`new leg (console/gates/guard/tickets) : ${missingNew.length === 0 ? 'complete' : 'INCOMPLETE -> ' + missingNew.join(' ')}`);
  console.log(`uncovered package patterns : ${missingPackages.length === 0 ? '(none)' : missingPackages.join(' ')}`);
  console.log('');

  const allOk = checks.every(Boolean);

  if (allOk) {
    console.log('RACE-COVERAGE: GREEN - all six packages covered by the Race leg');
    return 0;
  }

  const uncoveredShort = missingNew.length > 0 ? missingNew.map((p) => p.replace('./internal/', '').replace('/...', '')).join('/') : missingPackages.join(' ');
  const failing = checks.filter((ok) => !ok).length;
  console.log(`RACE-COVERAGE: RED - ${failing} check(s) failed; uncovered by Race leg: ${uncoveredShort}`);
  console.log('RACE-COVERAGE: RED - ' + (missingNew.length > 0
    ? 'console/gates/guard/tickets not covered'
    : 'Race leg present but does not satisfy the full six-package contract'));
  return 1;
}

process.exitCode = main();

'use strict';

/**
 * Gate report assembly.
 *
 * Output contract (one line per criterion, see the architecture design section 2):
 *     PASS|FAIL|SKIP|INFO <G-id> <label> :: <detail>
 * plus a final `TOTAL_FAIL=<n>` line.
 *
 *   PASS  the criterion executed and found no violation
 *   FAIL  the criterion executed and found at least one violation
 *   SKIP  not scripted in this slice, or a finding exempted by a whitelist entry
 *   INFO  nothing to judge (not a pass, not a failure)
 */

const VERDICTS = ['PASS', 'FAIL', 'SKIP', 'INFO'];

function createReport() {
  const records = [];

  return {
    records,

    rec(id, verdict, label, detail) {
      if (VERDICTS.indexOf(verdict) < 0) {
        throw new Error('invalid verdict: ' + verdict);
      }
      records.push({ id, verdict, label, detail: String(detail) });
    },

    fails() {
      return records.filter((r) => r.verdict === 'FAIL');
    },

    skips() {
      return records.filter((r) => r.verdict === 'SKIP');
    },

    render() {
      const lines = records.map(
        (r) => r.verdict.padEnd(5) + ' ' + r.id + ' ' + r.label + ' :: ' + r.detail
      );
      lines.push('TOTAL_FAIL=' + this.fails().length);
      return lines.join('\n');
    },
  };
}

module.exports = { createReport, VERDICTS };

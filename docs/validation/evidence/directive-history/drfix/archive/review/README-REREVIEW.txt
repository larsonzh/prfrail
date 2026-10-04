DR-FIX pilot slice — ⑦ RE-REVIEW BUNDLE (budget: final ⑦ call)
================================================================

Slice goal: remove the first-run-red flake class registered as DR-2..DR-5 in
tools/agent-probe/enforcement-proxy (all of them "target accepts = 0, want 1" style
counter/timing races in the test harness).

Baseline revision: a366fe3 (SW-5 report commit). Everything below is UNCOMMITTED
working-tree state; nothing has been committed or pushed for this slice yet.

Changed files (git diff --stat, working tree):
  docs/DELIVERY_DIRECTIVE.md                       |  3 +-
  docs/DELIVERY_DIRECTIVE_EN.md                    |  3 +-
  tools/agent-probe/enforcement-proxy/main.go      |  5 ++
  tools/agent-probe/enforcement-proxy/main_test.go | 82 ++++++++++++++++++++-----
  tools/gates/cjk-newwords.txt                     |  2 +

Files in this bundle
  IMPL.diff              code diff: main.go + main_test.go   (REGENERATED after the T1 fix;
                         the copy that shipped with the 1st ⑦ was stale — blob ab645e0)
  DOCS.diff              guideline diff: DELIVERY_DIRECTIVE.md + _EN.md + cjk-newwords.txt
  gate-remediation.txt   full §6.3 gate output for the state above (TOTAL_FAIL=0)
  T1-isolated-edit-diff.txt   the T1 edit in isolation; its left-hand blob 227f906 is the
                         commit-time diff blob of the implementer's delivered state, which
                         independently confirms this diff is the real delta (not a rewrite)
  A1-T1T4-injected.txt .. B-count10-verbose.txt   the 1st-round mutation/regression evidence
                         chain produced by ③ (unchanged, still valid)
  A7-T1-mutation-red.txt / A8-T1-tolerant-too-much-red.txt   new mutation evidence for T1
  B2-count10-T1.txt / B3-full-suite-T1.txt / B4-psself-test-T1.txt   new regression evidence

--------------------------------------------------------------------------------
Disposition of the 1st ⑦ findings
--------------------------------------------------------------------------------

T1 (High) — readLog() turned a torn trailing append by the concurrent logger into a
false failure ("log line is not JSON").
  VERDICT: FIXED.
  Fix shape (see IMPL.diff): a pure core parseLogLinesInto(payload) returns
  (records, tornTail, badLine, badErr); parseLogLines(t, payload) is a thin wrapper that
  still issues the ORIGINAL verbatim t.Fatalf("log line is not JSON: %v (%q)", badErr, badLine)
  for any complete unparseable line; readLog() re-reads/re-parses within a bounded 200 ms
  window (20 ms polls) and only skips a line that is (a) the last split element and (b) not
  newline-terminated. The whole payload is no longer TrimSpace'd, because the trailing-newline
  state is what decides torn-vs-corrupt.
  Evidence:
    - Permanent tests (in main_test.go): TestReadLogToleratesTornTrailingAppend (TP-A),
      TestParseLogLinesIntoStrictness (TP-B, two-directional), 
      TestParseLogLinesIntoAcceptsUnterminatedRecord (TP-C).
    - A7-T1-mutation-red.txt: mutation 甲-i (make a torn tail strict again) => TP-A goes RED
      with the original message => the tolerance body is load-bearing, the test is not a tautology.
    - A8-T1-tolerant-too-much-red.txt: mutation 甲-ii (skip ANY unparseable line) => TP-B goes RED
      => the tolerance is deliberately narrow; a corrupt COMPLETE line is still fatal.
    - B2-count10-T1.txt: go test ./tools/agent-probe/enforcement-proxy -count=10 => ok (4.273s).
    - B3-full-suite-T1.txt: go test ./... -count=1 => ok=14 FAIL=0.
    - B4-psself-test-T1.txt: package PowerShell self-test => checks=13 failures=0.
    - T1-isolated-edit-diff.txt: no pre-existing assertion, counter, deadline or predicate was
      touched by the T1 edit.
  Limits stated honestly (not claimed as verified):
    - TP-A uses a hand-constructed torn payload; no real racing writer was stressed, so
      "a genuine race can never go red" is argued from the mechanism, not load-tested.
    - -count=10 output is non-verbose (single "ok" line).
    - The PowerShell self-test ran under pwsh 7.x, not Windows PowerShell 5.1.
    - -race does not cover this package: local Windows has no gcc and the CI race job only
      covers ./internal/adapters/... and ./internal/chain/... . This remains a REGISTERED GAP.

A1 (Medium) — the §0.3 change-log v1.9 row said ⑦×1 while Appendix C.0e recorded the actual
2 paid ⑦ calls (self-contradicting ledger).
  VERDICT: FIXED.
  Fix shape: CN row tail now reads
  "... 豁免 ①⑤；**成本**：预算 ⑥×1 + ⑦×1、**实得 ⑥×1 + ⑦×2 = 3 次付费调用**（见附录 C.0e）。 |"
  and the EN mirror carries the same change at the same line position.
  Evidence (DOCS.diff + gate-remediation.txt):
    - G1-b: header=v1.9, last=v1.9, 2026-09-22/2026-09-22, rows=10 in BOTH files.
    - G3(a): CN +2/-1 vs _EN +2/-1 (symmetric), strict positional mirror parity all true.
    - G3(a.2): strict parity checked for the directive pair; G3(b) 0 mismatches in every
      key-field class (hash / run / glyph / to-do / not-reviewed).
    - G4a: BOM+LF OK for 5 files (mirror kept its BOM and stayed LF-only).
    - G4b: newCjk=221, unknown=[] with whitelist=20 — the one new character 甲 (in the
      "injection class 甲-i / 甲-ii" label that the CN authority and the EN mirror both use)
      was registered WITH its context in tools/gates/cjk-newwords.txt.
    - Also included in DOCS.diff: the new §12.4 OB-13 row (user's verbatim wording) recording
      the two-stage acceptance evidence for flaky-fix slices.

C1 (High) — the 1st ⑦ bundle did not contain the evidence chain behind the reported results.
  VERDICT: SATISFIED by this bundle. Every finding/claim above names a concrete artifact in
  this directory (see the per-finding lists). No result in this bundle is a carrier self-report
  without a pasted raw output: all commands were re-run by the verifier (③) and the gate output
  is the untouched full stdout.

T2 (Medium) — proposed changing the accept-counter termination to ">= want" and lengthening the
deadline.
  VERDICT: DECLINED (kept as designed), with reasons and measurements:
    - Relaxing to ">=" would destroy the exactly-once property that the DR-2 fix deliberately
      established (the assertion was strengthened FROM ">= 1" TO "== 1"). This is a
      STRENGTHENING, not a relaxation, and the review must not read it as a weakening.
    - The 3 s hard deadline exists to PREVENT AN INFINITE HANG WHEN SOMETHING IS GENUINELY
      WRONG; its design intent is NOT "wait long enough for a late arrival". Waiting longer
      only lengthens the failure path. This intent is recorded in the slice's ⑧ report so the
      ">=" proposal is not re-litigated.
    - Measured scale: the full package completes 10 iterations in 4.273 s
      (B2-count10-T1.txt), i.e. no accept-completion anywhere near the 3 s deadline; the
      observed first-run failures were "target accepts = 0", not "accepts = 1 arriving late".
      See also B-count10-verbose.txt (1st-round verbose 10x evidence).

--------------------------------------------------------------------------------
What this re-review must decide
--------------------------------------------------------------------------------
1. Is T1 properly fixed (torn trailing append cannot produce red; a corrupt complete line
   still can)?
2. Is A1 properly fixed, and does DOCS.diff preserve the CN/EN mirror contract?
3. Is C1 satisfied by the evidence chain in this directory?
4. Any NEW finding of severity Medium or higher that the fix itself introduces.

Note: this is the LAST budgeted ⑦ call for this slice (§7.5: 1 final + 1 re-review).

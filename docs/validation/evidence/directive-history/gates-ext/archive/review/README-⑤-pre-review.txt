GATES-EXT — ⑤ PRE-REVIEW INPUT BUNDLE (V4 Pro)
==============================================
Slice: GATES-EXT (the first [SLICE] under the delivery directive v1.9).
What it does: (a) adds §6.3 criteria **G6** (ledger-metadata consistency: G6-1 status tense / G6-2 cost clause /
G6-3 run-id format / G6-4 expiring literal) and **G7** (review-bundle completeness: G7-a artifacts block /
G7-b path+disposition resolvable / G7-c finding disposition cell / G7-d stopped-state marker) to
docs/DELIVERY_DIRECTIVE{,_EN}.md; (b) turns the master's throwaway tmp/gate.js into a real, offline,
exit-code-correct runner under **tools/gates/** (12 criterion modules + selftest + fixtures + whitelists);
(c) wires it into CI (one additive step) and updates the frozen workflow contract test accordingly.
Boundaries honoured: G1–G5 criteria text unchanged (only the §6.3 heading count, R6.2's G1–G5 → G1–G7,
and the scripting-status note were edited); base **G1** (three-way status consistency) and base **G5**
(change-set contract) are deliberately NOT implemented and are printed as `SKIP` (never PASS);
`.github/workflows/ci.yml` gained exactly ONE step with no `continue-on-error`; `internal/release/workflow_test.go`
gained exactly TWO edits (the `test` job's step list + one `expectedScripts` digest entry).

## Read these first
- `tmp/gates-ext/SLICE-DEFINITION.md` — the nine-field definition, the three blocking acceptance steps, the risks.
- `tmp/gates-ext/ARCH-DESIGN.md` — the design proposal you are asked to review the implementation against
  (§2 G6-1…4, §3 G7-a…d, §5 layout + migration + fragility list, §6 test points T1–T13, §7 risks).
- `tmp/gates-ext/review/IMPL-tracked.diff` — diff of the tracked modifications
  (`ci.yml`, `internal/release/workflow_test.go`, both directive files, `tools/gates/cjk-newwords.txt`).
- `tmp/gates-ext/review/IMPL-newfiles.diff` — unified diffs of every new production file under `tools/gates/**`
  (`gate.js`, `lib/ctx.js`, `lib/report.js`, `lib/criteria/*.js`, `selftest.js`, `g6-whitelist.txt`).
- `tmp/gates-ext/review/NEWFILES.txt` — inventory (path / bytes / sha256) of all new paths, fixtures included.
- `tmp/gates-ext/review/GATE-all-tree.txt`, `GATE-selftest.txt`, `GATE-range-historical.txt` — raw gate and
  selftest outputs, plus the two historical retro-tests (A1 ⇒ G6-2 FAIL on `range:b738e6b^..b738e6b`;
  the fixed side `9f68a52^..9f68a52` ⇒ green; C1 ⇒ G7-a FAIL on `range:1b01115^..1b01115`).
- The working tree itself is readable: `tools/gates/**`, `docs/DELIVERY_DIRECTIVE{,_EN}.md`.

## MANDATORY DISCLOSURE — what has NOT been through an independent test round
1. **Batch ②f** (the NF-1/NF-2/NF-3 fixes: leading word boundary on every G6-3 rule, the §6.3 G6-3
   wording precision, and `gate.js` printing usage errors before exiting 2) and
2. **the master's doc additions** (the §6.3 G6-3 "known residual false positive" note, the §12.4 rows
   **OB-17** and **OB-18**, and the `融` whitelist entry)
were verified by **re-running ④** (the gate, the 168-assertion selftest, a corpus-wide zero-false-positive
scan) — which is the action §5.3 prescribes for code/test fixes — but **they were NOT put through a
further ③ round**. ⇒ **Decide explicitly whether they need additional review, and say so in your report.**
Do not assume ③ covered them.

## What ⑤ must do
1. **Proposal coverage, item by item**: for each ARCH item (the ⒞ entry-form verdict, G6-1…G6-4, G7-a…G7-d,
   the tooling layout, the migration/fragility fixes, the step-1 equivalence procedure, T1–T13, the CI step,
   the boundary self-check), state whether the implementation **implements / deviates / omits** it, quoting
   the file:line you relied on. Deviations already documented by the implementer (and their reasons):
   fixture files are `.txt` mapped to virtual `docs/**/*.md` paths; G6-2 uses the *nearest* qualifier;
   G7-d requires a status-line shape; the step-1 recipe was corrected (forward-apply + a strict
   tracked-modification equivalence run). Judge each on its merits.
2. **Three classes of counter-examples**, each falsifiable: (a) a case where a G6/G7 rule **passes something
   it should fail**; (b) a case where it **fails something legitimate** (this is the expensive failure mode:
   the script blocks CI); (c) a case where the **scope discipline** breaks (a pre-existing committed line
   being blamed, or an added line being missed).
3. **Challenge the two design decisions ③ rated acceptable** — backticked short ids still judged; digit runs
   glued to a following word not judged — plus the residual classes ③ tabulated (judged vs not judged).
   State whether the CI-blocking risk (R1) is now adequately controlled.
4. **Judge the acceptance evidence**: are the migration-equivalence evidence, the historical retro-tests
   (A1 = real historical range, C1 = real historical range, **OB-11 = fixture only**) and the two-way
   mutation evidence sufficient for the claims made? Say explicitly where evidence is weaker than the claim.

## Roles and isolation
- ⑤ may read everything above (its job is implementation-vs-design).
- Do **not** treat this bundle as a list of findings to confirm — verify independently from the code.
- ⑥ and ⑦ will receive a **sanitized** bundle that omits ③/⑤ findings (§7.2 isolation). Do not write your
  findings into files under `tools/gates/**` or the directive.

## Report format
Per item: verdict + the file:line evidence. Then the three counter-example classes (each with a concrete
input + expected vs actual). Then **findings** with severity (`High`/`Medium`/`Low`/`Note`) and a
falsifiable scenario, and an explicit `PASS` / `PASS WITH FIXES` / `FINDINGS` line. State anything you could
not verify from the bundle alone (you have no execution capability — do not claim to have run anything).

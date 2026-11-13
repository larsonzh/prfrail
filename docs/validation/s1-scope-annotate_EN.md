# ProofRail S1-SCOPE-ANNOTATE report (P71 · layer-2 governance side · documentation slice [SLICE] · size M)

Date: 2026-11-13; status: step-8 documentation assembly finalised (steps 1-7 closed). Slice identifier: `S1-SCOPE-ANNOTATE` (P71 · layer-2 governance side · documentation slice `[SLICE]` · size M).

## 1. Executive Summary

- Slice: `S1-SCOPE-ANNOTATE` (P71, layer-2 governance side, documentation slice `[SLICE]`, size M), controller `DeepSeek V4.1 Flash (deepseek)`; explicit same-turn authorisation from the user on 2026-11-13.
- Delivered: two REQ stage columns, the MATRIX Windows row plus its tail note, two TEST_STRATEGY AT rows, two s1-exit AT rows plus the applicability note, the OPERATIONS visible section, and three additive pointers.
- Conclusion: the X scheme landed (no entry deleted, no numbering changed, the MVA composite untouched); the S1 scope and stage ownership are annotated item by item per P68 b1; the three visible boundaries are in place.
- Metering: step 5 once (`PRE-REVIEW: PASS WITH FIXES`, 5 Low plus several Info, no Medium+; the two documentation-layer Low items fixed and the three design-layer criterion Low items turned into §4 content assertions, see §5), step 6 **skipped** (confirmed by step 7), step 7 blind 1 + final 1 + re-review 0 (`BLIND REVIEW: PASS`, `FINAL REVIEW: PASS`), GLM batch 2 times, step 2g retried twice because of a network interruption.
- Counts: **R = 0, D = 2** (D=1: the step-2a execution carrier created `p71-2a-annotate.js` at the `tmp/` root and deleted it itself; D=2: the step-2d execution carrier created the worktree `p71d-wt` at the `tmp/` root and removed it; both left zero residue and are registered under the §5.1 self-healing tier).

### 1.1 Exceptions and Deviations

a) **F-1 (the kick-off message conflicts with an upstream ruling; ruled by the user in the same turn)**: item 4 of section [V] of the kick-off message originally said "change all four stage columns of REQ-009 / REQ-019 / REQ-007 / REQ-021 to 'not required by S1'", which directly contradicts the P68 ruling "self-hosting is not approved for removal (REQ-019 / AT-13 stay in S1)". Step 1 stopped and reported the facts; the user ruled in the same turn that F-1 = option A (annotate item by item per the P68 ruling: REQ-009 changes, REQ-007 changes at sub-item level, REQ-021 stays, REQ-019 stays). This slice executed option A. This is an instance of "the literal kick-off text conflicts with an upstream ruling"; the lesson is that a kick-off message must be checked item by item against the actual scope of the upstream ruling.
b) **Step-3 carrier deviation**: both dispatch attempts for the consistency check failed with a `UND_ERR_SOCKET` network interruption, so the controller ran it as a deterministic script (mechanically recomputable); registered as an execution-carrier deviation (§13.3.18) whose independence impact is limited to "not reviewed by another party".
c) **Step-2g duplicate call**: the first 2g dispatch was interrupted by the network but had already landed; the retry was only a byte-for-byte cross-check. Under R7.1 a re-send within the same stage counts against the cap, so it counts as 2 calls.

## 2. Difference Inventory

8 mandatory items:

| Document | Verdict | Basis |
|---|---|---|
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | mandatory | the REQ-007 stage column becomes "S1/S2 (black-box candidate sub-item deferred to S2)" (sub-item level, not a whole-entry removal) |
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | mandatory | the REQ-009 stage column becomes "not required by S1, deferred to S2 (basic terminal retained)" |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | mandatory | the Windows row constraint cell gains "basic terminal retained; the full integrated TUI leaves S1 and moves to S2" |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | mandatory | a P71 slice-sync tail note is appended (wording without the ambiguous "support scope unchanged") |
| `docs/TEST_STRATEGY.md` / `_EN` | mandatory | the AT-13 owner cell gains "self-hosting stays in S1 (P68 b1 does not approve removal)" |
| `docs/TEST_STRATEGY.md` / `_EN` | mandatory | the AT-24 owner cell gains "the black-box sub-item is deferred to S2 (P68 b1); the conditional closure is not revoked" |
| `docs/validation/s1-exit.md` / `_EN` | mandatory | the AT-13 / AT-24 evidence cells each gain their stage ownership |
| `docs/OPERATIONS.md` / `_EN` | mandatory | a new visible debugging-observer section plus the three boundaries |

6 items requiring evaluation:

| Document | Verdict | Basis |
|---|---|---|
| `docs/validation/remaining-rulings.md` §5 | applicable (1 additive pointer line) | R-07 b1 condition 2 landed in this slice; the original line is unchanged (§11.1) |
| `docs/S1_RELEASE_NOTES.md` support section plus gate section | not applicable | the support section already says "the full TUI is not provided" and "the visible black-box candidate is not runnable for release"; the current state already agrees |
| `docs/validation/p57-closeout-assessment.md` §4.2 | applicable (1 additive pointer line) | this slice wrote back the s1-exit applicability note; the eight gate definitions and verdicts are unchanged |
| `docs/validation/s1-reach-reachability.md` §2 | not applicable | gate reachability and REQ stage ownership are different dimensions |
| `docs/DEV_PLAN.md` T018 section | not applicable | T018 is task-level while REQ stage ownership is product-level |
| `docs/validation/s1-stream-progress.md` §8 | applicable (1 additive pointer line) | P71 drew the blind review (before the final review), so that requirement has landed |

## 3. Per-Document Write-Back

- `PRODUCT_REQUIREMENTS`: the REQ-007 stage column = `S1/S2 (black-box candidate sub-item deferred to S2)`; the REQ-009 stage column = `not required by S1, deferred to S2 (basic terminal retained)`; REQ-019 stays `S1` and REQ-021 stays `S1/S2` (F-1 option A).
- `S1_SUPPORT_MATRIX`: the Windows row constraint cell appends "basic terminal (line-oriented CLI / `interactions tui` / `--progress`) retained; the full integrated TUI leaves S1 and moves to S2"; the tail note appends "P71 slice sync" and states explicitly that the platform list and support tiers are unchanged, the Windows row stays "release-support candidate" and the visible row stays "degraded guarantee".
- `TEST_STRATEGY`: the AT-13 owner cell appends "stage ownership: self-hosting stays in S1 (P68 b1 does not approve removal)"; the AT-24 owner cell appends "stage ownership: the black-box sub-item is deferred to S2 (P68 b1); the conditional closure is not revoked"; the requirement and assertion columns of both rows are unchanged.
- `s1-exit`: the AT-13 / AT-24 evidence cells each gain their stage ownership; a new "P71 applicability note" section states that the REQ stage-ownership change does not alter any gate verdict, the nine gate-table rows keep their names, count and status, the `S-Gx` numbering is not mapped in, and the existing AT-12 / AT-16 / AT-22 PASS evidence is unaffected.
- `OPERATIONS`: a new "visible debugging observation (convenience, not a deliverable)" section plus three boundaries (not an S1 delivery commitment / not authoritative evidence, ADR-012 / requires VS Code + Copilot Chat + the bridge extension, and real-CLI end-to-end still depends on T027-B).
- Three additive pointers: one pointer line each in `p57-closeout-assessment` §4.2, `remaining-rulings` §5 and `s1-stream-progress` §8; the original lines are unchanged (§11.1).

## 4. Counterexamples and Verification

- X-1 (delete a REQ entry or change its number): the criterion `grep -c '^| REQ-'` = 28; step 3 measured `REQ_count=28` with REQ-001–028 contiguous.
- X-2 (the visible row claims "release support"): only the Windows row of the matrix table carries that phrase; step 3 measured `matrix_rows_with_zhengshi=1` and the matching row is `Windows 11 amd64`.
- X-3 (a BLOCKED gate turned into PASS): the gate table is compared row by row; step 3 measured `gate_table_identical=true` (2 BLOCKED / 4 PASS / 3 PARTIAL unchanged).
- X-4 (the REQ-019 stage column changed): step 3 measured `REQ-019 stage=[S1]` and `REQ-021 stage=[S1/S2]`; neither moved.
- X-5 (the AT-13 / AT-24 assertion columns changed): only the stage column may move; step 3 measured `changedFieldIdx=[3]` on both changed REQ lines (column 4 untouched) and `identical=true` for the PC-01 / PC-06 rows.
- Step-3 mechanical readings, summarised: `REQ_count=28`; the four stage values = `REQ-007 [S1/S2 (black-box candidate sub-item deferred to S2)]` / `REQ-009 [not required by S1, deferred to S2 (basic terminal retained)]` / `REQ-019 [S1]` / `REQ-021 [S1/S2]`; column 4 identical; the two PC rows identical; the gate table identical; `S-G3` / `S-G4` / `S-G5` hits = 0; visible "release support" hits = 0; every CN / EN pair has symmetric numstat (G3(a) reports nine pairs with equal added / deleted counts).

## 5. Review Chain

- Step 5 `PRE-REVIEW: PASS WITH FIXES`: 5 Low plus several Info, no Medium+. Each is disposed of below (classified as "documentation layer" versus "design-layer criterion" so that an assertion never substitutes for a fix):
  - **Low-1 (documentation layer ⇒ must be fixed)**: original text: "the note at `S1_SUPPORT_MATRIX.md:33` says 'the release-support scope itself is unchanged' while the Windows row constraint cell gains 'the full integrated TUI leaves S1' — 'support scope unchanged' and 'the S1 commitment surface shrinks' are in word-level tension" ⇒ fixed: the note now reads "the platform list and the support tiers are unchanged; the TUI carve-out is a stage-ownership sync already ruled by P68". Landing point: the tail note of `docs/S1_SUPPORT_MATRIX{,_EN}.md`.
  - **Low-2 (documentation layer ⇒ must be fixed)**: original text: "`s1-exit.md:62` — 'the full TUI carve-out does not constitute a gate-verdict change' is an assertion without an argument: the AT-12/16/22 assertions are all based on the basic terminal / `interactions tui` and do not depend on the full integrated TUI, so the existing PASS is not weakened" ⇒ fixed: the note now adds "the AT-12 / AT-16 / AT-22 assertions do not depend on the full integrated TUI, so the existing PASS evidence is unaffected". Landing point: the applicability note of `docs/validation/s1-exit{,_EN}.md`.
  - **Low-3 (a design-layer criterion defect ⇒ turned into a content assertion)**: original text: "the criterion `grep -c '^\| REQ-'` — under BRE `\|` means alternation, so `^` matches every line and a literal execution returns the file's total line count instead of 28; the criterion is written incorrectly" ⇒ turned into: the criterion is `grep -c '^| REQ-'` = 28, together with the contiguity of REQ-001–028. Landing point: §4 **X-1** of this report.
  - **Low-4 (a design-layer criterion defect ⇒ turned into a content assertion)**: original text: "a whole-file grep for the 'release support' phrase also matches the 'release-support scope' wording inside the P71 note, so 'only the support-candidate row matches' is literally unsatisfied" ⇒ turned into: a table-row-scoped assertion "only the Windows row of the matrix table carries that phrase (the single row whose status column is 'release-support candidate'), and the visible row must not contain it". Landing point: §4 **X-2** of this report.
  - **Low-5 (a design-layer criterion defect ⇒ turned into a content assertion)**: original text: "the line anchors of X-3 / X-4 (`:52` / `:60` / `:81`) are correct today, but any insertion or deletion earlier in the file shifts them and the mechanical check is fragile" ⇒ turned into: content assertions "the gate table compares row by row as identical plus the status-word counts 2 BLOCKED / 4 PASS / 3 PARTIAL" and "REQ-019 / REQ-021 are read via their third column value". Landing point: §4 **X-3** and **X-4** of this report.
  - **Classification**: Low-1 / Low-2 are **documentation-layer** wording tension and a missing argument ⇒ the audited document is fixed; Low-3 / Low-4 / Low-5 are defects **in the design-layer criteria themselves** (wrong syntax, over-broad scope, fragile anchors) and concern no audited document's correctness ⇒ they become recomputable content assertions in §4 and **do not constitute "substituting an assertion for a fix"**.
- Step 6 **skipped**: pure `.md` write-back, no `.go` / `.json` / `.yml` / schema / fixtures / script-logic semantic change, and no touch to the state machine, gate criteria, evidence model or role authorisation body ⇒ the three §7.2.1 conditions hold; confirmed by the step-7 final review.
- Step 7 blind `BLIND REVIEW: PASS` (no step-5 / step-6 list attached, drawn before the final review); step 7 final `FINAL REVIEW: PASS` (0 findings).
- The blind and final conclusions agree with no divergence; step 7 re-review ran 0 times (the final review had 0 findings, so there was no fix loop).

## 6. Gate Readings

- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`.
- `--scope=ci` ⇒ `GATE REPORT scope=ci changed=4`, `TOTAL_FAIL=0`.
- Step-9 final readings (after staging the 31 paths precisely, the evidence pack as a whole directory, no `-A`): `--scope=index` `changed=31` `TOTAL_FAIL=0`; `--scope=tree` `changed=31` `TOTAL_FAIL=0`; `--scope=ci` `changed=4` `TOTAL_FAIL=0`; `--selftest` `SELFTEST: PASS 737/737`.
- The only pre-staging tree red, `G7-b` (7 `repo` artifacts still untracked: 2 under `02-batch/`, 3 under `04-reviews/` and the report pair itself), turned green once the precise `git add` ran.
- Before the report pair existed, the index / tree readings showed `G5(d)` red (structural: no staged path was declared yet; resolved once the report fence existed) and `G5(f)` (unstaged, handled under discipline E), both green after staging; `G8-b` was red as well (the pack "archive entry list" and `SHA256SUMS.txt` did not yet list the five `02-batch/` and `04-reviews/` payloads) and turned green once this slice completed them at step 8.
- Base gates: `gofmt -l .` empty; `go build ./...` exit 0; `go vet ./...` exit 0; `go test ./...` ok=15 / FAIL=0.

## 7. Known Limits

- (1) The `s1-exit` gate table has nine rows of Chinese gate names and no `S-Gx` row identifiers; `S-Gx` is defined in `docs/validation/p57-closeout-assessment.md` §4.2, and this slice establishes no mapping.
- (2) `S-G2` has no row of its own and its wording is "the gate-level verdict awaits the v54 snapshot"; this slice does not change it.
- (3) **The "S2 commitment narrowing registration" for REQ-020 is outside this slice** (a step-5 Info item suggests registering it; it is left to a later slice).
- (4) `S1_SUPPORT_MATRIX` has no TUI row, so the TUI carve-out is expressed through the Windows row constraint plus the tail note rather than a new row.

## 8. Explicitly Not Executed

- Not changed: the directive body / `CONTRACTS` / `SECURITY` / `ARCHITECTURE` / `ADR_REGISTER` / the v50–v56 snapshots / the REQ table structure / the MVA composite / the S1 Exit eight-gate verdicts themselves / product code.
- No commit / push / gitee push; no new slice started.
- **D=2 recurrence-prevention clause (ruled by the user on 2026-11-14)**: when an execution role (`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`) works under `tmp/`, it must always use `tmp/<sliceID>/` and must never create any file, directory or `git worktree` at the `tmp/` root. Both D=2 instances in this slice are cases of that violation (the step-2a `p71-2a-annotate.js` and the step-2d `p71d-wt` worktree); a third occurrence (in-slice D=3) forces a slice stop under the §5.1 cumulative ceiling. **The P72 kick-off discipline must cite this clause.**

- ⑥ skipped: pure .md write-back (no .go/.json/.yml/schema/fixtures/script-logic semantic change), the three §7.2.1 conditions hold, confirmed by the step-7 final review.

```artifacts
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/00-briefing.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/01-arch/design-r2.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/02-batch/batch-01-req-enumeration.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/02-batch/batch-02-at-gates-enumeration.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/04-reviews/blind-r1.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/04-reviews/final-r1.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/MANIFEST.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/SHA256SUMS.txt	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/freedom-list.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/input-list.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/s1-scope-annotate.md	repo
docs/validation/s1-scope-annotate_EN.md	repo
docs/validation/s1-stream-progress.md	repo
docs/validation/s1-stream-progress_EN.md	repo
```

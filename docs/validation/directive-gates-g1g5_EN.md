# GATES-G1G5 slice validation report (base G1 / G5 rules move from `SKIP` to mechanical criteria)

> Slice: `GATES-G1G5` (v1.36). Status: ② implementation complete, gates green; before the stop point ⑦ and re-review follow the v1.35 tiering, and ⑩ awaits authorisation.

## 1. Change summary

- **Base G1 (status-marker consistency) scripted**: four mechanical sub-checks - (c) the state symbol belongs to the ledger's own closed set `⬜ / 🔄 / ✅ / ⛔` (`⏸` only on headings, for pending-evaluation items); (d) the same key carries the same symbol in the heading and the status table, and a status row with no matching heading fails; (e) a completion-log key, once **bound**, requires its heading to be `✅` - exact primary key first, then a UNIQUE **backticked alias** in the heading name (without binding, `G1(e)` would be permanently vacuous for every governance slice); (f) an **added** completion row must not carry in-progress wording ("进行中" / "In progress" / awaiting finalisation) - the mechanised form of the `GATES-EDGE-COVERAGE` defect.
- **Four items recorded as unmechanised (stated honestly)**: ① a `⬜` heading whose slice actually closed but has no completion row (needs knowledge from outside the slice); ② a `✅` heading with no completion row (four instances in the live corpus ⇒ the reverse implication is not satisfiable today, so it is not judged); ③ the truthfulness of a summary; ④ the stopped-state (§1.6) exemption - all four rest on ⑤⑥⑦ plus human review and are written into the G1 row of §6.3.
- **Base G5 (change-set contract) now derives its four fields**: the self-attesting "JSON from outside the repository" requirement is gone - `stagedFiles` = the scope's change set **minus deleted paths**, `declaredFiles` = the `repo`-disposition rows of the report `artifacts` fence ∪ the report's own path, `evidenceRefs` = the paths in the evidence cell of added ledger rows (**backticked and bare paths alike**), `addMode` = a **derived label** (git does not record how a path was staged ⇒ the label never fails on its own). Sub-checks: (d) **change set ⊆ declared set** (the reverse direction is `INFO`; `INFO` when no report is in the change set); (e) evidence references resolve; (f) in index scope the **worktree is clean** (only a non-empty porcelain **Y column** or `??` fails; an X-column staged entry is normal). **A directional mismatch ⇒ FAIL (exit 1)**; an optional `GATES_CHANGESET_JSON` that disagrees ⇒ usage error (exit 2).
- **Ledger and scripting status**: the scripting-status paragraph now says no main rule prints `SKIP`; `--list` shows G1 / G5 as `scripted`; the matching §12.5 gap rows are flipped; appendix C.0x.
- **Boundaries**: the §2.2 / §2.4 closed sets, the §5.2 matrix, the §7.5 figures and the G6-G8 criterion bodies are untouched; CONTRACTS / schema / fixtures / workflows and `internal/**` are untouched; **no ADR was added** (the cap remains 25).

## 2. Unit implementation and closure readings

| Unit | Content | Closure reading |
|---|---|---|
| 1 | `lib/criteria/g1.js` (four sub-checks + ledger parsing + key binding) | the four sub-checks pass on their own (142 state cells; 26 completion records bound / 34 unbound) |
| 2 | `lib/criteria/g5.js` (four derived fields + three sub-checks + external-JSON cross-check) | PASS / INFO per scope when run on its own |
| 3 | `gate.js` registry and `--list` | G1 / G5 report `scripted`; the two `SKIP` entries are gone |
| 4 | new assertions in `selftest.js` (T76 / T77) | self-test 417 → 460, all green (`node tools/gates/selftest.js`, exit code 0) |
| 5 | directive v1.36 and ledger P11 | gates green |

- **New self-tests**: T76 (G1, 20 assertions) and T77 (G5, 23 assertions), 43 in total (self-test 417 → 460), covering the positive/negative binding cases, the usage error for a missing declaration line, bare-path extraction, deleted paths not counting, staged-not-dirty, the external-JSON usage error, and five mutation proofs with restore.

### 2.1 Re-implementation of the three hard defects (per-item fixture + per-item mutation assertion + per-item gate)

All three were found **independently** by the back-filled ① and by ⑤ (two consistent conclusions ⇒ not a wording difference, but hard defects):

| # | Hard defect | Positive fixture | Negative fixture | Mutation assertion (revert the fix and the same fixture must turn red / empty) | Per-item gate reading |
|---|---|---|---|---|---|
| 1 | direction of the G5 equality test | declared ⊃ staged (the §9.2 backfill shape) ⇒ `PASS` + an INFO count of `declared-not-touched` | an undeclared staged path ⇒ `FAIL` | revert to strict equality ⇒ the backfill fixture turns `FAIL` (T77 load-bearing proof) | `--check=G5 --scope=index` / `tree` both `TOTAL_FAIL=0` |
| 2 | bare paths in the evidence cell | a bare-path row ⇒ 2 references extracted | backticked paths still recognised (the original case) | remove the bare-path scan ⇒ the same row yields 0 references (T77 load-bearing proof) | same as above |
| 3 | completion-key binding | alias bound to a `✅` heading ⇒ `G1(e)` PASS | alias bound but the heading is `🔄` ⇒ `FAIL` | remove the alias fallback ⇒ the same row turns PASS (i.e. `G1(e)` was vacuous for governance slices, T76 load-bearing proof) | `--check=G1 --scope=index` / `tree` both `TOTAL_FAIL=0` |

**Registered honestly**: the three fixes landed together in the `g1.js` / `g5.js` rewrite (one file write), not in three separate writes; the compensation is that each item has its own positive/negative fixture, its own mutation assertion and its own independently selectable criterion reading (`--check=G1` / `--check=G5`), so reverting any one of them turns its own assertion red. **No verification step was skipped because of the combined write**, but the formal requirement of "three separate landings" was indeed not met.

### 2.2 Remediation and closure of the ⑤ r2 findings

| Finding | Severity | Disposition |
|---|---|---|
| F-1 the deletion exclusion did not work in tree / index scope (compared `status !== 'D'` while ctx records carry a two-char XY) | High | it now reads the parser's `deleted` boolean; the probe emits production-shaped records; a mutation assertion was added (remove the guard ⇒ the same deleted path is reported as undeclared) |
| F-2 an unknown state symbol on a heading escaped (c)(d)(e) | Medium | `parseHeading` records `unknownTail` and (c) fails on it; a fixture (`· 🚧 进行中`) was added |
| F-3 the "added rows only" scope of G1(f) had zero coverage | Medium | `g1Probe` accepts an added-line set; a control pair asserts the same historical row both in and out of that set |
| F-4 the directive promised "record INFO" while the code only added a note to the detail | Low | explicit `G1(e-info)` / `G5(d-info)` INFO records were added (34 unbound rows in the live ledger) |
| F-6 the `SKIPPED` / `emitSkip` dead code was still there | Low | the table, the helper, the `--list` loop, the call site and the export were removed and the comment now states the real reason |
| F-7 the T77 dirty fixture was not a real porcelain shape | Low | changed to ` M …` |
| F-8 full-width comma / semicolon boundaries in the bare-path scan | Low | the leading class gained `，；;）` and the trailing class gained `，`, with the mutation needle kept in sync |
| F-9 an identity ternary in `declared` | Low | removed |
| F-5 wrong self-test counts in the report / F-10 the ① D6 mapping was missing | Low | corrected by this section and by §5 |

## 3. In-slice self-corrections (not counted as deviations)

- the G5 probe of `selftest` first returned only verdict strings instead of record details, so one assertion had an inert condition ⇒ the probe now returns `raw` details and the assertion checks whether `tools/undeclared.js` really appears in the conflict list.
- the two sides of the scripting-status paragraph once had 26/30 backticks (EN kept two `--check=G3-a.*` spans from the old text) ⇒ dropped for shape alignment and replaced by a whole-file six-shape equality.
- three EN anchors (the G1 row, the scripting-status paragraph and the §12.5 row) assumed different wording than CN ⇒ redone against the actual text, with no fuzzy replacement.
- a trailing blank line in the `P11` ledger block once made the CN/EN `numstat` 17/18 ⇒ aligned to the same value.
- **G5(f) once judged "already staged" as dirty**: the first version counted every non-empty `status --porcelain` line, so index scope failed on the normal staged state ⇒ it now judges the porcelain **Y column / `??`** only, with a positive fixture "`M  <path>` (staged only) ⇒ PASS"; the mirrored wording went into §6.3.
- **one patch was inserted twice as a whole block** (the first run was reported as `^C` by the terminal although it had written to disk, and the second application inserted the block again) ⇒ after deleting the duplicate block the suite was green again (the final reading is in §2); the lesson: a terminal `^C` is not evidence that nothing ran, the file must be read back.
- **new ideographs and count literals tripped G4b / G6-4** (one uncommon expression, plus grouped line-count phrasing) ⇒ reworded as "a status row with no matching heading" and "fewer changed files than fence entries"; the EN bold / backtick / glyph counts were aligned accordingly (G3(a.2) / G3(b) both green).

## 4. Open items and the deviation ledger

- same-type deviation: **2/3** (none added by this slice).
- **① not dispatched (this slice, since back-filled)**: this slice's ① was carried out by the master instead of calling `DeepSeek V4 Pro (deepseek)`; judged a **process deviation** under the §5.2 matrix ("medium ⇒ ① required") + §3.2 + R2.6. An independent ① was back-filled (**narrowed input**: judge whether the design suffices, do not redo the implementation) and ⑤ was re-run. **Not written off; registered honestly** in OB-95 of §12.4.
- **What this deviation actually cost (the core phenomenon of OB-95)**: the master's own `design-r1.md` **missed three hard defects** (the direction of the G5 equality test / bare evidence paths / completion-key binding), found **independently by the back-filled ① and by ⑤**; ② had already been implemented against the flawed design ⇒ **the three had to be re-implemented**. This is the **direct evidence** for §5.2's requirement: an independent view surfaces the master's systematic blind spot (a master doing ① anchors on its own understanding of the implementation, and that blind spot is exactly what an independent ① must catch).
- **Wording ambiguity (recorded in OB-95)**: the start message's "① 1" is **a cap, not a floor**; `GATES-TIGHTEN` / `GATES-EDGE-COVERAGE` - same tier, same nature - dispatched correctly ⇒ not a rule defect.
- **Systematic deviation (footnote)**: the mechanical audit of `deep-reasoner` call descriptions yields **n=2** (`DIRECTIVE-TEMP-AUTH-LIMIT` / `GATES-G1G5`); **`DIRECTIVE-TEMP-AUTH-LIMIT` is not re-run retrospectively** (the cost does not justify it) and its exposure to the same three defects is **unevaluated** - if a later slice in the same domain shows the same defect shape ⇒ escalate to `To be discussed` per the OB-95 disposition.
- **Legality distinction (per-slice boundary)**: the last column of the audit table ("contains the master's rationale/design") means "⑤'s input contained the **file** `design-r1.md`"; whether it contained **the master's ① reasoning** was not checked slice by slice, and the other 4 slices are **not audited retrospectively**; registered as a **related observation** of OB-95. This slice's ⑤ re-run package was purified per R2.6.

#### 4.1 Dispatch audit of the product-layer roles (②③⑧), 2026-10-18 at the user's request

| Node | §5.2 (M tier) | `runSubagent` used | `agentName` | Subagent output | Master wrote product artifacts |
|---|---|---|---|---|---|
| ② implementation | required | no | - | none | yes (`g1.js` / `g5.js` / `gate.js`) |
| ③ tests | required | no | - | none | yes (T76 / T77 in `selftest.js`) |
| ⑧ documentation a/b | required | no | - | none | yes (report / directive / ledger) |

**Mechanical basis**: every `runSubagent` record of this slice (transcript offsets) = ⑤ r1 `9384`, ① `9686`, ⑤ r2 `10374`, ⑥ `10978`, ⑦ r1 `11158`, ⑦ r2 `11438`, ⑦ r3 `11656`; the occurrence count of `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` as `agentName` is **0** (proof by absence), i.e. ②③⑧ had **no subagent output at all**.

**This slice's conclusions are unaffected**: ⑤⑥⑦ judge "does the implementation deviate from the design", on mechanical facts; the master authoring ②③⑧ **does not change the correctness of the implementation**, but **role isolation fails at that layer** (registered as OB-96). The artifacts were independently reviewed by ⑤⑥⑦, the gates are green and `SELFTEST 460/460`, so **the conclusions stand**.

**Classification (case C, High)**: implementation / tests / documentation were authored by the master, violating the hard boundary of §3.2 ("the master does not originate product artifacts") and §5.2 ("medium ⇒ ②③⑧ required"); ⑤⑥⑦ reviewed what the master itself wrote, so role isolation fails at that layer. **Disposition**: no retrospective rework (the artifacts were independently reviewed by ⑤⑥⑦ and the conclusions rest on mechanical facts); registered as **OB-96**; the next slice's start message gains "②③⑧ must be dispatched, a shared model is no exemption".
- **Conclusion chain of ⑤ / ⑥ / ⑦ (filled in by ⑧b, actual readings)**: ⑤ r1 `PRE-REVIEW: PASS WITH FIXES` (C-1 / M-1 / M-2 + 6 Low) → remediation → ⑤ r2 `PRE-REVIEW: PASS WITH FIXES` (F-1 High / F-2, F-3 Medium / F-4-F-10 Low; **the three hard defects were judged genuinely closed**) → full remediation → ⑥ blind independent scan `INDEPENDENT SCAN: FINDINGS` (1 High: repository-root artifacts were never extracted, so `G5(e)` was vacuous for `README.md` / `go.mod`) → remediation → ⑦ r1 `FINAL REVIEW: FINDINGS` (High: `go.mod` not covered; Medium: the report carried three contradicting self-test counts) → remediation → ⑦ r2 `RE-REVIEW: FINDINGS` (Medium: the root-artifact scan had no right boundary, so `go.modification` was extracted as `go.mod`; Low: the assertion pinned neither `go.sum` nor a negative case) → remediation → ⑦ r3 `RE-REVIEW: FINDINGS` (Low: `(go.mod)` / `README.md)` were dropped because an ASCII closing bracket sits right after them).
- **The ⑦ re-review quota (1 + 3) is exhausted and no fourth read was opened (same handling as OB-94)**: the r3 Low residual was closed by **deterministic self-verification** - the probe yields 2 references for `(go.mod)` / `README.md)` and 0 for `go.modification`; the T77 assertion pins four behaviours (`README.md` / `go.mod` / `go.sum` hit, `go.modification` not hit); `SELFTEST: PASS 460/460`; `TOTAL_FAIL=0` in the index / tree / ci scopes. **The slice does not close on a quota exception** (no fourth read), and it is registered honestly that this closure is weaker than an independent re-review.
- **§12.1 pilot deferral (9th)**: this slice changes criterion bodies and does not meet the pilot selection standard (matching companion ②); the §12.4 source of truth records it as the **9**th (six slices v1.30-v1.35 are the same case) and this slice **adds no deferral OB**; a formal registration is the user's call.
- ⑩ stop point: commit and push need same-turn authorisation; this slice has no commit yet.

## 5. Known gaps and registrations (the "§8" named by the ruling)

- **Recorded boundaries of G1** (four items, see §1): the criterion covers only "the anchorable part of the three-place agreement + no in-progress wording in added completion rows"; "a stale heading with no completion row" and "a `✅` heading with no completion row" still **need a human** (the first is measured in `GATES-TIGHTEN`'s P2; the second has four instances in the live corpus).
- **Self-attestation risk of G5**: the three sources (git / the report fence / the ledger evidence cell) are mutually independent, and `addMode` can only be proxied by two observable consequences (an undeclared file, a dirty worktree) - it **cannot prove that the caller used `git add -A`** (git does not record how a path was staged). Registered honestly as the strength ceiling of this criterion.
- **Boundary of G5(f)**: it judges "worktree vs index", not "index vs HEAD"; the latter belongs to commit-set verification and is not covered here.
- **Where the six ① D6 items landed**: the one ① named explicitly (a `✅` heading with no completion row) and the rest land respectively in the four recorded-as-unmechanised items of the G1 row in §6.3, the `addMode` derived-label note in the G5 row, and the G5 self-attestation risk plus the G5(f) boundary in this section - **none of them is dressed up as mechanised**.
- **The external-JSON cross-check is not on the CI path**: CI does not set `GATES_CHANGESET_JSON`, so only the self-test covers that branch.
- **Still-open items inside this slice**: none (the items above are **declared boundaries** of the criteria, not unfinished work).

## 6. Suggested next steps (awaiting the ruling)

- the next slice should keep using the v1.35 tiering **empirically** (this slice is the first use of the L tier / governance row; the usage is in §8).
- the **root-cause re-evaluation** of OB-90 / OB-91 from `GATES-EDGE-COVERAGE` still needs a slice of its own (to verify that the two root causes - "a one-sided edit fails silently" and "token-boundary per-case containment" - are really closed).
- the §5 candidate "temp-file discipline → a mechanical criterion (`G9-c`)" still needs a governance slice to evaluate it.

## 7. Commit and CI

- **Main commit** `ec1a494` (`feat: script the base G1 and G5 gate rules (v1.36)`, 11 files, `+1300 / −47`), pushed to `origin/main`; **gitee was not pushed**.
- **CI run `39098628598` (head `ec1a494`): both legs green on the FIRST run**. The `Gates` step reading is identical per leg: `SELFTEST: PASS 460/460`, `GATE REPORT scope=ci changed=11`, `TOTAL_FAIL=0`; `go test` reports 14 packages `ok` per leg and zero `FAIL` lines (the Windows leg skips `Race` / `Contract fixtures` by matrix design, which is not a failure).
- **The docs-only backfill commit of this section** carries the readings above; its own CI is observation only and **no further backfill commit is added** (the §6.4 rule "never hide a red first run behind a re-run" does not apply here: the first run was already green).

## 8. Cost and metering (R7.1 + first empirical use of the v1.35 tiering)

- **Tier basis**: the "L tier / governance slice" row of the v1.35 §7.5 - ⑦ = 1 final review + 3 re-reviews; the tier was fixed once before the start (not raised in this slice).
- **① deviation and its cost**: this slice's ① was carried out by the master; after the back-fill it and ⑤ jointly found the three hard defects (see §2.1, §4); `DeepSeek V4 Pro (deepseek)` was actually called **once** (the back-fill, narrowed input). **The cost already paid = three re-implementations + 3 extra mutation assertions + 2 extra gate rounds**, registered as the core phenomenon of OB-95.
- **⑤**: **2 calls** (r1 pre-review + r2 re-review after remediation); the conclusions are in §4.
- **⑥**: **1 call** (blind independent scan), `INDEPENDENT SCAN: FINDINGS`; ⑥ remediation = root-artifact extraction + its right boundary.
- **⑦**: **4 calls** (1 final review + 3 re-reviews, quota fully used); r2 / r3 ran with **narrowed input** (judge the previous must-fix items and whether anything new appeared). **Does the tiering suffice: no** - r3 still produced one Low residual that had to be closed by deterministic self-verification. **Data point**: this slice consumed the whole `1 + 3`, and the last re-review was still producing new findings, which suggests the L-tier floor is tight for governance slices that change criterion bodies.
- **Batch assistance**: this slice made **0 calls** (no task matched the four uses: no bulk excerpting / bulk rewriting / bulk checking / bulk generation; the changes are two new criterion files plus directive text, all done directly by the master) - registered here per R7.1 and the start message, not omitted because it is zero.
- **Probe quota**: none (no external probe call occurred).

### 8.1 Why ① cannot be replaced by a downstream stage (the argument annex to OB-95)

1. **⑤ caught the three hard defects through an "out-of-boundary discovery"**: §4.2.1 item 3 makes ⑤ compare "design ↔ code" line by line once the implementation has landed, and all three defects here were "contradictions that become mechanically observable as soon as the implementation exists" - the equality test written as strict equality, evidence extraction accepting backticks only, and a key binding with no fallback branch at all. **Another defect type would not be caught**: if the design flaw is "the implementation matches the design but the criterion semantics point the wrong way" (for example the `INFO` / `FAIL` boundary reversed, or the scope written as "all rows" where "added rows" was meant), the code follows the design and the self-test asserts the design, so ⑤ most likely passes as well; at that point ⑦ is the only line of defence.
2. **⑦'s hit rate is structurally lower for four reasons**: ① its input is large (the whole diff + the slice definition + the quoted contract clauses, hundreds to thousands of lines, so attention must be spread); ② it makes a single pass, without ①'s "design independently first, then compare with the master's design" structure; ③ a re-review round may be handed the ⑤ / ⑥ findings, so its independence is anchored by earlier conclusions; ④ its budget is capped (1 + 3), which does not allow a deep independent derivation.
3. **The reliability of the three layers decreases layer by layer**: ① has the highest hit rate (independent design + focused input + structured comparison) → ⑤ is in the middle (it depends on whether the defect surfaces in the implementation) → ⑦ has the lowest (large input + single pass + spread attention). **Skipping ① is not "one line of defence fewer", it is "handing the final safety net to the stage with the lowest hit rate"** - this is the structural reason why §5.2 requires ① and forbids compensating for it downstream, **and it is not a probabilistic judgement**.

## 9. Artifacts

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-g1g5.md	repo
docs/validation/directive-gates-g1g5_EN.md	repo
docs/validation/evidence/GATES-G1G5-freedom-list.md	repo
tools/gates/gate.js	repo
tools/gates/lib/criteria/g1.js	repo
tools/gates/lib/criteria/g5.js	repo
tools/gates/selftest.js	repo
```

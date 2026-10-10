# Slice validation report · `DIRECTIVE-HAIKU-BUMP` (directive v1.27)

> Slice id: `DIRECTIVE-HAIKU-BUMP`; directive v1.27; date 2026-10-11. Class `[SLICE]`; size S; probes 1; `[ROUTINE]` not applicable.
> **Status**: **⑧b finalized** (the local part of ⑨ native validation has run; readings in §7 and §7.7; the ⑧c mechanical verification and wrap-up record is in §11). **Authorization boundary**: commit NO / push NO / gitee NO / probes 1 (exhausted, no top-up); committing and pushing **require same-turn authorization**.

## 1. Goal and scope

- Slice nature: a pure documentation directive-governance slice, class `[SLICE]`, size S, probes 1, `[ROUTINE]` not applicable.
- Trigger: the user's 2026-10-11 kickoff declaration; this slice is a pure documentation directive-governance slice ⇒ **the pilot is deferred for the 6th time** — registered independently at step 0 as **OB-86** under §1.7.4 / the OB-37 meta-rule (four elements + truthful annotation + "shall not be used as a precedent").
- Three landing goals: (1) §7.2.3's enabling precondition version moves from `Claude Haiku 4.5` to `Claude Haiku 5.5`; (2) §7.2.3 gains a "call string and internal identifier" paragraph (qualified name `Claude Haiku 5.5 (copilot)`, internal identifier `claude-haiku-5.5`, plus an explicit statement that its internal identifier is **not** a member of the §6.3 G5-b whitelist closed set) ⇒ closes **OB-74** in §12.4; (3) §12.4's OB-74 disposition column flips to `已处置`.
- Out of scope (boundary stated literally): the §2.2 whitelist closed set, the §2.4 "sole exemption" wording, the §7.5 cost table, the §6.3 G5-b, code, `docs/CONTRACTS{,_EN}.md`, `schemas/**`, contract fixtures, `.github/workflows/**`, `go.mod`, the role set and carriers, `.github/agents/**`, `internal/**`.
- Implementation path: **path A (minimal string replacement)**, as explicitly stated in the user's remark.

## 2. Change summary

- Change set = `git diff --cached` (numstat): **3 items at the ⑨ time point** (the first 3 rows below; `staged == worktree`, `git status --short` shows `M ` / `A ` with no second column); **7 items in the finalized state** (plus this report pair and the ledger completion rows, see the ⑧c re-run readings in §11). All artifacts are **precisely staged**.

| Artifact | numstat |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | +26/−3 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | +26/−3 |
| `docs/validation/evidence/DIRECTIVE-HAIKU-BUMP-freedom-list.md` | +23/−0 (new) |
| `docs/validation/directive-haiku-bump.md` | +201/−0 (new, this file pair; artifact rows measured by G7-b, see §11) |
| `docs/validation/directive-haiku-bump_EN.md` | +201/−0 (same as above) |
| `docs/t027/REMAINING_SLICES.md` | +1/−0 (⑧b write-back, report path only) |
| `docs/t027/REMAINING_SLICES_EN.md` | +1/−0 (⑧b write-back, same as above) |

- The four touches of directive v1.27 (listed truthfully in the `§0.3` v1.27 log row): §0.3 / §7.2.3 (version string + new paragraph) / §12.4 (OB-74 disposition column + OB-86 registration row) / Appendix C.0p; §12.4 also carries, at its end, the registration note for this slice's ⑦ round user-authorization exception (§7.8.3, no OB number taken).
- Landing point (version and minimality): §7.2.3's first sentence **only swaps the version string**; the rest of its text is untouched character for character (mechanically re-verified for CN and `_EN` by A12 experiment (d)).
- Landing point (OB-74 closure): the new §7.2.3 paragraph gives the qualified name and the internal identifier and states literally that the internal identifier is **not** a member of the G5-b closed set (that closed set still equals the §2.2 `modelId` column); §12.4's OB-74 disposition column flips to `已处置` while its **phenomenon column is untouched character for character** (re-verified by both A12 experiment (e) and ⑦ read 2).
- Landing point (ledger and ledger metadata): §12.4 gains the OB-86 registration row (the step-0 pilot deferral for the 6th time); Appendix C.0p is added with its status row written as `⬜ 进行中` and **no premature completion claim**.

## 3. Execution pipeline and environment

- Pipeline: ① proposal → ②a / ②b rule landing → ③ test (consistency review) → ④ gates → ⑤ first pre-review → ⑥ independent scan → ⑦ read 1 (final review) → ⑦ read 2 (narrowed-input re-review) → ⑧a this report → ⑨ native validation (local re-run + the §7.7 A12 dedicated independent experiment) → ⑧b finalization → ⑧c wrap-up.
- Roles and models: ① **not dispatched** (user ruling); ② / ③ / ⑧ are carried by `DeepSeek V4.1 Flash (deepseek)`; ⑤ first pre-review = `GLM-5.3 (glm)` (the first time that model takes ⑤ after entering the whitelist); ⑥ = `MAI-Code-1.1-Flash (copilot)`; ⑦ = `GPT-5.3-Codex (copilot)`, tier Extra High, **not downgraded**.
- Probe: the user-ruled availability probe (purpose = measure whether `runSubagent` accepts `Claude Haiku 5.5 (copilot)`), budget 1; **Haiku carried out no review work**.
- Environment: Windows host; repository `d:\LZProjects\prfrail` (the `prfrail` folder of the multi-root workspace); gate entry `node tools/gates/gate.js`.

## 4. Anomalies and fallback record

- **A-1 Repair preceded review**: ⑤ read 1's Low (the Appendix C.0p deliverable row omitted §12.4's OB-86 registration row) was **closed in place under the OB-11 exception** (ledger-metadata class, five elements self-judged) **without re-running ⑤**; the repair happened **before** ⑥ and ⑦ ⇒ ⑥ / ⑦ reviewed the **post-repair** state.
- **A-2 ⑥'s coverage gap (second occurrence of the family)**: ⑥ declared `INDEPENDENT SCAN: PASS` but **did not cover** the must-verify items in the driver's list that depend on a baseline or on running commands ("byte-for-byte comparison against HEAD", "independent gate re-run", "probe accounting shape"), and **did not** list them in Section E's "cannot verify" ⇒ the gap is carried by ⑦ and ⑨'s A12 experiment under §7.8.2; the driver **did not** add a ⑥ call (budget ≤ 3, this slice used 1). **This is the same family as `DIRECTIVE-GLM-LANDING`'s A-2, its second occurrence** ⇒ see C-4 in §13.
- **A-3 Driver dispatch defect (the shared root of ⑦ read 1's two Mediums)**: the driver asked ⑦ to complete "byte-for-byte comparison against HEAD" and "independent gate re-run" on its own, but the `independent-reviewer` carrier has **read-only file tools only and no execution tool** ⇒ both items were **impossible** in that dispatch shape. Remedy = the driver **materialized** the baseline into the input package for the re-review round (`head-baseline.txt` / `staged-diff.txt` / `gate-*.txt`, see §6). **A driver orchestration error; shall not be used as a precedent.**
- **A-4 Driver citation error (corrected by ③)**: in the ③ dispatch the driver mis-cited the authority for "writing one caliber in two places ⇒ drift is inevitable" as §1.7; the correct authority is the **Appendix B.7 hard constraint**. ③ corrected it and landed F-5 accordingly.
- **A-5 Probe attempt 1 timed out (408)**: attempt 1 was an infrastructure timeout (`408 Timed out reading request body`, **no model verdict**) ⇒ the driver stopped per §1.7.2 and escalated; the user, **in the same turn**, ruled that attempt **does not constitute a valid probe** and authorized **one retry** (basis = an infrastructure error rather than a model rejection + the same family as OB-64 + clearing Haiku 4.5's stale `thinking.budgetTokens` and restarting VS Code); probe attempt 2 succeeded. The **strict-reading ambiguity** was recorded as a Low by ⑦ read 1 and judged CLOSED by read 2; disposition = the accounting sentence in §9 + C-5 in §13.
- **A-6 Observable variance in the ⑦ carrier's tool availability**: both ⑦ reads in this slice self-reported "constrained by read-only file tools, no commands executed", whereas ⑦ in earlier slices claimed "independent re-runs" of the gates ⇒ the evidence calibre of "⑦ independent re-run" is **inconsistent across rounds**. This slice truthfully claims only ⑦'s **text-level** verification; re-run readings are supplied by the driver with their source labelled ⇒ see C-6 in §13.
- **A-7 ⑦ Medium #2: driver self-verification was overruled by the user and the item was closed by an authorized narrower read**: ⑦ read 2 judged that end "partially closed" (C.0p claims the tree scope is all green, while the input package supplied only the index and selftest raw lines); the driver first closed it by re-running and self-verifying, then the **user ruled in the same turn that driver self-verification is not accepted** and **authorized one additional narrower ⑦ read** (§7.8.3's "user-authorization exception"; the authorization text and the narrowed-input declaration are registered in §12.4) ⇒ read 3 gave `RE-REVIEW: PASS` (0 findings) and the Medium is **independently closed** by ⑦. **A driver orchestration error (self-verification substituting for review); shall not be used as a precedent**.
- **A-8 Self-inflicted criteria defects in the A12 experiment's first run** (recorded truthfully): the first run produced 1 script crash (the EN header carries an English gloss ⇒ column location failed) and 5 FAILs, **all** of which were defects of the script's own criteria (comparing the whole row instead of the phenomenon column; the `**not**` bold markers interrupting a literal; the old version string being misjudged when it appears as a **legitimate citation** in the §0.3 and C.0p historical rows) ⇒ after fixing the criteria, `A12_EXPERIMENT_FAILURES=0` (50 items). That round **proved no artifact defect**, and **no artifact was adjusted to suit a criterion**.
- **A-9 This report was edited twice after ⑧b**: the ⑧c readings (encoding / residue / gate re-run) and the record of the **user-authorized read 3** were both taken after ⑧b finalization and written back (the former matching §5.1's order, the latter under a same-turn user authorization) — both edits **changed numeric rows and review records only**; no conclusion was revised, and **shall not be used as a precedent** (the ideal shape = the driver records the ⑧c readings and the authorized read's verdict in a report appendix and keeps zero edits after ⑧b).

## 5. Review conclusion summary (⑤⑥⑦)

- **⑤ read 1**: `PRE-REVIEW: PASS WITH FIXES` (**no Medium+**), 1 Low + 2 Info: (1) Low — the Appendix C.0p deliverable row **omitted** §12.4's OB-86 registration row ⇒ **repaired** (1 place each in CN / `_EN`, closed in place under the OB-11 exception); (2) Info — the probe accounting (1 failed attempt + 1 valid probe) must go into ⑧a; (3) Info — the freedom list must be landed **before** ⑦ and attached to its input package.
- **⑥ independent scan**: `INDEPENDENT SCAN: PASS` (0 file:line-level findings; Sections A–E complete); the **coverage gap is truthfully recorded** (see A-2); the driver **did not** add a ⑥ call.
- **⑦ read 1 (final review)**: `FINAL REVIEW: PASS WITH FIXES` — **2 Mediums** (both **driver input-package defects**: no baseline-comparison evidence against HEAD; no independent gate re-run, see A-3) + **1 Low** (the strict-reading risk of §7.2.3's "do not retry repeatedly"); **F-1…F-11 all judged in-bounds** (F-8 is the in-bounds shape of "escalate first, then execute").
- **⑦ read 2 (re-review)**: `FINAL REVIEW: PASS WITH FIXES`, **no new findings**; read 1's Medium #1 **CLOSED** (it performed the textual comparison itself from `head-baseline.txt`: frozen regions `IDENTICAL=true`, OB-74 phenomenon column untouched character for character), the Low **CLOSED** (the disposition path maps literally onto the Appendix B.7 / §7.8.4 clauses); Medium #2 **partially closed** (the tree-scope raw line for the final state was missing ⇒ see read 3 below).
- **⑦ read 3 (user-authorized narrower read)**: `RE-REVIEW: PASS` (**0 findings**) — it checked only whether C.0p's tree / index / selftest claims have raw gate-reading support; authorized verbatim under §7.8.3's "user-authorization exception" (registered in §12.4), with the input narrowed to four items, and it did **not re-review the whole slice**; the read closed read 2's Medium #2 **independently** via ⑦, leaving no residual Medium+ in this slice.

## 6. Review input package summary (including blind-review isolation proof)

- **⑦ read 1 input package**: slice definition + authoritative probe facts + change set (3 items) + freedom list + 10 must-verify checks (**no** baseline material ⇒ A-3).
- **⑦ read 2 input package (narrowed re-review)**: additionally `tmp/haiku-bump/head-baseline.txt` (the **verbatim HEAD and worktree texts** of both CN / `_EN` across the five frozen regions + the OB-74 row + the §7.2.3 first paragraph), `staged-diff.txt` (the complete staged diff, proving the change set is exactly 3 files), `gate-index.txt`, `gate-selftest.txt`, `gate-tree.txt`, together with read 1's findings table and the freedom list.
- **⑦ read 3 input package (user-authorized narrower read)**: **only** the Appendix C.0p acceptance-evidence row verbatim (CN and `_EN`) plus the raw output text of the three commands tree / index / selftest (four items); it **did not** attach the slice definition, the change set, the freedom list, or any ⑤ / ⑥ / earlier ⑦ list (§7.8.3's narrowed-input declaration is registered in §12.4).
- **Blind-review isolation proof**: this slice **adds no blind review** (it does not touch §7.7's three mandatory-sampling semantics; the C.0p remark records the judgement and its basis before kickoff); ⑦ = 1 final + 1 re-review + **1 user-authorized narrower read** = **3 reads** (the first two within §7.5's budget of 1 + 1; read 3 under §7.8.3's "user-authorization exception" and **valid for this slice only**), with no resend, no voided round and **no automatic fourth read**.
- Sanitization: the review input contains only in-repository paths and text; no credentials, no external connection material.

## 7. Falsifiability and gate results

| Item | Reading |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `changed=3`, `TOTAL_FAIL=0` (⑨ time point; finalized-state readings in §11) |
| `node tools/gates/gate.js --all --scope=index` | `changed=3`, `TOTAL_FAIL=0` (same as above) |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 335/335` (this slice **added no assertion**; same baseline as the previous slice) |
| `G1-b` version metadata | `header=v1.27 last=v1.27 2026-10-11/2026-10-11 rows=28` (OK for both CN and `_EN`) |
| `G3(a)` / `G3(b)` | `+25/−3` against `_EN` `+25/−3`; all five shapes (lines / headings / bold / bars / blanks) `true`; key fields `glyph 1/1 OK` (the rest 0/0) |
| `G4a` / `G4b` | `3 file(s) OK (1 excluded: evidence tree)`; `newCjk=268 unknown=[]` |
| `G6-1` … `G6-6` | all PASS (`scanned 50`; `G6-6` `scanned 4`: every new ledger row opens with a closed-set status word) |
| `G8-a` / `G8-d` | `1 path(s) OK`; `1 changed file(s) OK` |
| Base gates | `gofmt -l .` empty; `go build ./...` exit 0; `go vet ./...` exit 0; `go test ./...` `ok=14 FAIL=0` |
| Frozen-region independent re-check | §2.2 / §2.4 / §6.3 / §7.5 / Appendix B are **byte-for-byte untouched** in both CN and `_EN` (A12 experiment 10/10 `identical`; ⑦ read 2 independently re-checked as IDENTICAL at the text level) |

- Except for the frozen-region row, every reading above is a **local mechanical reading** (measured at ④ and ⑨); `G7-a` … `G7-d` were `INFO` at the ⑨ time point (no changed report was in scope then), with finalized-state readings in §11.
- **The raw readings behind C.0p's tree / index / selftest claims**: in the finalized state, `--all --scope=tree` ⇒ `GATE REPORT scope=tree repo=D:\LZProjects\prfrail changed=7` and `TOTAL_FAIL=0` (exit 0); `--scope=index` gives the same reading; `--selftest` ⇒ `SELFTEST: PASS 335/335` (raw outputs in `tmp/haiku-bump/gate-tree-final.txt` and its two siblings). **⑦ read 3 independently verified that claim against those three raw readings (`RE-REVIEW: PASS`)**, with no reliance on driver self-verification.

### §7.7 A12 dedicated independent experiment

- Because ⑥ reported `PASS`, ⑨ must contain one experiment, independent of ⑥, that proves the key invariants mechanically: `tmp/haiku-bump/verify-a12.js` ⇒ **all 50 items OK, `A12_EXPERIMENT_FAILURES=0`**, in seven parts.
- (a) **Criteria live-fire** (a synthetic in-memory ctx, bypassing git, writing nothing to the repository): `deepseek-v4-pro` / `glm-5.3` / `glm-5.3-flash` / `gpt-6-luna` PASS; **`claude-haiku-5.5` FAIL** (this slice's key fact: the internal identifier is not a closed-set member), `claude-haiku-4.5` FAIL, `Claude Haiku 5.5 (copilot)` (qualified-name form) FAIL, `"claude-haiku-5.5"` (quoted) FAIL; a missing `model` key PASS (§2.6), an empty `model: ""` PASS (deferring to G5-a (2)).
- (b) **Document ↔ code invariant equality in both directions**: the §2.2 `modelId` deduplicated set (8 members in each of CN and `_EN`) == `g5b.js`'s `CLOSED_SET` (8 members), and CN equals `_EN` element for element; `claude-haiku-5.5` and `claude-haiku-4.5` are **both absent** from the closed set.
- (c) **Boundary invariant**: §2.2 / §2.4 / §6.3 / §7.5 / Appendix B are **byte-for-byte untouched** in both CN and `_EN` (10/10 `identical`; e.g. CN §2.2 `8d61224cd3d03432`, CN §6.3 `a01f89cc35d4e052`, EN §2.2 `87cbf393367c7b81`, EN Appendix B `8f53bcb2c2e8058b`).
- (d) **Minimality**: §7.2.3's first paragraph satisfies "only the version string differs" in both CN and `_EN` (replacing `Claude Haiku 4.5` with `Claude Haiku 5.5` in the HEAD text makes it character-for-character equal to the worktree).
- (e) **OB-74 closure is literally checkable**: §7.2.3 carries the qualified name and the internal identifier plus an explicit G5-b exclusion sentence (CN and `_EN`); §12.4's OB-74 disposition column opens with a closed-set status word (CN `**已处置**` / `_EN` `**Resolved**`) and its **phenomenon column is byte-for-byte identical to HEAD**; the old version string survives only in the §0.3 v1.27 row and C.0p's goal row (2 occurrences each, `stray=0`).
- (f) **The G6-6 closed-set trap**: the status closed set = `观察中 / 待办 / 待议 / 已处置 / 已裁决 / 已登记`, and it does **not** contain `已落地`.
- (g) **C.0p structural facts**: C.0p exists, its status is `⬜ 进行中`, its probe budget row reads "1 次（硬上限，不追加）", and it states literally that "Haiku carried out no review work for this slice".

## 8. Known boundaries and residual risks

- **C-1 Observation-item candidate (not landed in this slice)**: the §2.2 primary/backup annotation and the four-category bulk-assist limitation **still have no mechanical criteria coverage** (carried over from the previous slice); the user ruled that this slice does not land it ⇒ see §13.
- **C-2 Known boundary**: §7.5's "cost-gradient hint (measured pricing)" lists only Flash / V4 Pro, which is **not the same source** as the §2.2 GLM entry ⇒ this slice **does not touch §7.5** (boundary stated literally) and records it as a known boundary.
- **C-3 Known boundary**: §6.3's `G4b` **blocks newly introduced CJK characters** ⇒ this report's wording stays within the existing character table and **does not extend the whitelist** (this slice's reading `unknown=[]`).
- **A-2 Coverage gap**: the must-verify items ⑥ did not cover are carried by ⑦ and ⑨ under §7.8.2; **must not** be read as covered by ⑥.
- **A-7 closed**: ⑦'s Medium #2 was **independently closed** by the user-authorized read 3 (`RE-REVIEW: PASS`); **driver self-verification must not substitute for independent review** (the same situation must go through §7.8.3's authorization exception, see C-7 in §13).
- **Boundary stated literally**: this slice touches neither the §2.2 closed set nor the §2.4 wording nor the §7.5 cost table nor the §6.3 G5-b; the §2.4 "sole exemption" wording is untouched character for character.

## 9. Explicitly not executed

- **⑨ has run (local)**: both the local re-run and the §7.7 A12 dedicated independent experiment are complete, with readings in §7; the **CI legs have not run** (authorization boundary commit NO / push NO ⇒ no CI run); this report **does not claim** that CI has run.
- **Not committed, not pushed**: the change set rests in the staged state; **gitee not pushed**; lifting the ⑩ stop point **requires same-turn authorization**.
- **Probe accounting**: budget 1 (hard cap, no top-up) — 1 infrastructure-timeout attempt (408, **no verdict**, **not counted as a valid probe**) + 1 **valid probe** (`runSubagent` accepted `Claude Haiku 5.5 (copilot)`, carrier `quick-verifier`); **same-turn authorization permits a single retry only**; **no top-up** (see A-5).
- **Call census**: ⑥ 1 call; ⑦ **3 calls** (1 final review + 1 re-review + 1 **user-authorized narrower read**; read 3 is valid for this slice only under §7.8.3); **no downgrade ADR**; **Haiku carried out no review work** (availability probe only).
- C-1 is not landed in this slice; cost comparison is out of scope.

## 10. Deviations and degrees of freedom

- Same source as the freedom list; item by item in `docs/validation/evidence/DIRECTIVE-HAIKU-BUMP-freedom-list.md` (F-1 through F-11); verdict = **F-1…F-11 are all in-bounds** (driver self-report confirmed by both ⑦ reads).
- F-1: under path A, **only the version string is swapped**, without rewriting the rest of that paragraph (in-bounds).
- F-2: OB-74's "qualified name + internal identifier" is written into §7.2.3 **itself**, **not** into §2.2 / §2.4 (solved within the boundary).
- F-3: the new §7.2.3 paragraph **states the G5-b exclusion literally**, preventing a reader from concluding the closed set must be extended (in-bounds).
- F-4: the internal identifier is named `claude-haiku-5.5` (following the existing historical notation `claude-haiku-4.5`) (implementation choice).
- F-5: the shape of ③ read 1's Low repair = turning the restatement of the "sole exemption" into a **pure pointer** (authority = **Appendix B.7**) (in-bounds).
- F-6: the shape of ③'s Info repair = making §12.4's OB-74 disposition column **self-contained** (phenomenon column untouched) (in-bounds).
- F-7: ⑤'s F-1 repair is closed in place under the **OB-11 exception** without re-running ⑤; **that exception is not generalized** (in-bounds).
- F-8: the probe accounting shape = 1 failed attempt (not counted) + 1 valid probe (in-bounds, user ruling in the same turn + the §1.7.2 escalation fulfilled).
- F-9: the cost vector goes into the validation report only (G6-2 caliber) (implementation choice).
- F-10: C.0p's status row reads `⬜ 进行中` and **does not** pre-claim completion (in-bounds).
- F-11: the probe serves **availability detection** only; Haiku is not dispatched to carry review work (in-bounds).
- **Deviations (recorded truthfully)**: A-3 (driver dispatch defect), A-7 (⑦'s residual closed by driver self-verification), A-8 (self-inflicted criteria defects in the A12 script's first run); **A-3 and A-7 shall not be used as precedents**.

## 11. ⑧c wrap-up cleanliness record

- Encoding check (`tmp/haiku-bump/wrap-enc.js`, full `git ls-files -z`): non-frozen tree `tracked=3415 text=3326 binary=89 violations=0`; the frozen evidence tree counts separately at `violations=9` (all under already registered / confirmed calibers: B2's CRLF, B3b's BOM, and the six `glm-eval/archive/**` items normalized by archived type).
- Residue: untracked residue **0**; the staged state holds 7 items (`M ` x4 + `A ` x3) and `git status --short` has no `??` line; the `tmp/haiku-bump/` process artifacts are listed in §16.
- IDE diagnostics: **0 errors** across the 7 changed artifacts.
- `node tools/gates/gate.js --all --scope=tree` re-run at the ⑧c time point ⇒ `changed=7`, `TOTAL_FAIL=0`; `G7-a` … `G7-d` moved from `INFO` to PASS, and `G7-b` measured 62 artifact rows (31 per report); `--scope=index` gives the same reading.
- The frozen regions (§2.2 / §2.4 / §6.3 / §7.5 / Appendix B) are not rewritten; this slice is zero-diff against them (A12 experiment 10/10 `identical`).
- **Ordering deviation**: see **A-9** in §4 (readings taken after ⑧b; writing them back added this section's numbers afterwards).

## 12. Cost and metering

- R7.1 per row (field caliber = `{role, model, stage, startedAt, reason, retryOf}`): ① **not dispatched** (user ruling); ⑤ `GLM-5.3 (glm)` 1 call (first pre-review, **no re-review**, reason = read 1 had **no Medium+**); ⑥ `MAI-Code-1.1-Flash (copilot)` 1 call (a valid call, budget ≤ 3); ⑦ `GPT-5.3-Codex (copilot)` **3 calls** (final review + narrowed-input re-review + **user-authorized narrower read**; budget = 1 + 1 plus one same-turn pre-authorized read, the latter **not generalizable and valid for this slice only**; **not downgraded, no added blind review**; no resend, no voided round); probe `Claude Haiku 5.5 (copilot)` 1 **valid** call (plus 1 infrastructure-timeout attempt, not counted).
- Carriers (all `DeepSeek V4.1 Flash (deepseek)`): ②a 1 / ②b 1 / ③ 1 / ⑧a 1 / ⑧b 1 ⇒ **10** calls in total (⑤ 1 + ⑥ 1 + ⑦ 3 + Flash 5).
- **The cost vector goes into this section only and not into the directive (§0.3 / Appendix C.0p)** (G6-2).

## 13. Candidate observation items (not landed in this slice)

- **C-1**: the §2.2 primary/backup annotation and the four-category bulk-assist limitation **have no mechanical criteria coverage** (carried over from the previous slice; both ⑦ and ⑥ had independently flagged it).
- **C-2**: §7.5's "cost-gradient hint" is **not the same source** as §2.2 after the GLM entry (§7.5 untouched, recorded as a known boundary).
- **C-4 (new)**: ⑥'s **coverage gap occurred a second time** (A-2) ⇒ the next slice should assess whether ⑥'s must-verify list needs to be layered by "can this be done with read-only file tools" or should state explicitly which items are inherently carried by ⑦ / ⑨.
- **C-5 (new)**: the probe **retry caliber** (A-5) — the directive should state "an infrastructure timeout with no verdict does not count as a valid probe; same-turn authorization permits a single retry only" (this slice records it on the report side; the directive side awaits a dedicated slice); **user ruling: not landed in this slice**; assess at the next related slice's kickoff whether to establish it as a work item.
- **C-6 (new)**: the ⑦ carrier's **tool availability** is inconsistent across rounds (A-6) ⇒ it affects how the evidence calibre of historical "⑦ independent re-run" claims should be read; the next slice should pin down ⑦'s executable surface and the evidence shapes it may claim; **user ruling: not landed in this slice**; assess at the next related slice's kickoff whether to establish it as a work item.
- **C-7 (new)**: this slice showed in practice that "both of ⑦'s Mediums can be caused entirely by a missing input package" ⇒ the next slice should define a **minimum list of baseline and raw readings** for ⑦'s input package (this slice remedied it by materialization, see A-3 / A-7; this slice's data point for A-7 = driver self-verification overruled by the user and closed via a §7.8.3-authorized narrower read); **user ruling: not landed in this slice**; register it in the next slice's kickoff batch.

## 14. Follow-up recommendations (including tier recommendations)

- Follow-up recommendations (recommendations only, not registered): C-1 / C-2 and the new C-4 … C-7 go to the next slice's kickoff for batch registration; the pilot deferral count uses §12.4 as its single source of truth.
- Tier recommendation: for later pure-documentation and criteria-body slices, ⑦ should still be `GPT-5.3-Codex (copilot)`, **not downgraded**; if ⑦ must re-run commands independently, the dispatch should state its tool surface or **the driver should materialize the baseline instead**.

## 15. ⑩ stop-point report (including write-back check rows)

- Executed through **⑧b finalization + ⑧c wrap-up + ⑦ read 3 (user-authorized narrower read)**; ⑨'s local part has run; the **⑩ stop-point receipt has been issued**; the user's **same-day bundled pre-authorization** released that stop point (conditions = read 3 with no Medium+ and, after remediation, both gate scopes plus selftest green — all met) ⇒ the commit and push facts are in §9's "Commit and CI" (added by the docs-only write-back commit).
- **Write-back check rows** (§11.2 caliber): validation report = **written back** (this file pair); `DEV_PLAN` = **not applicable** (a directive-governance slice writes back to Appendix C, not to `DEV_PLAN`); ledger = **written back** (`docs/t027/REMAINING_SLICES{,_EN}.md` completion row, **report path only**); ADR = **not applicable** (no model downgrade / no budget exception).
- Authorization boundary (at the ⑩ stop point): **commit NO / push NO / gitee NO**; under the user's **bundled pre-authorization** it was then executed as **origin-only** (SSH) commit and push, with **gitee NO** and no probe top-up.

## 16. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/evidence/DIRECTIVE-HAIKU-BUMP-freedom-list.md	repo
docs/validation/directive-haiku-bump.md	repo
docs/validation/directive-haiku-bump_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tmp/haiku-bump/verify-a12.js	local-only
tmp/haiku-bump/boundary-check.js	local-only
tmp/haiku-bump/make-baseline.js	local-only
tmp/haiku-bump/head-baseline.txt	local-only
tmp/haiku-bump/staged-diff.txt	local-only
tmp/haiku-bump/gate-index.txt	local-only
tmp/haiku-bump/gate-tree.txt	local-only
tmp/haiku-bump/gate-selftest.txt	local-only
tmp/haiku-bump/gate-index-final.txt	local-only
tmp/haiku-bump/gate-tree-final.txt	local-only
tmp/haiku-bump/gate-selftest-final.txt	local-only
tmp/haiku-bump/parity.js	local-only
tmp/haiku-bump/parity-out.txt	local-only
tmp/haiku-bump/wrap-enc-out.txt	local-only
tmp/haiku-bump/a12-out.txt	local-only
tmp/haiku-bump/pre-05-diff-cn.txt	local-only
tmp/haiku-bump/pre-05-diff-en.txt	local-only
tmp/haiku-bump/pre-05-verify.js	local-only
tmp/haiku-bump/pre-05-verify-out.txt	local-only
tmp/haiku-bump/pre-05-gate-index.txt	local-only
tmp/haiku-bump/pre-05-gate-tree.txt	local-only
tmp/haiku-bump/pre-05-selftest.txt	local-only
tmp/haiku-bump/wrap-enc.js	local-only
tmp/haiku-bump/	local-only
```

Note: `repo` rows are this slice's staged artifacts (including this report pair and ⑧b's ledger completion rows); `local-only` rows are one-shot process artifacts under `tmp/haiku-bump/`, removed at ⑧c; **excluding** leftovers from previous slices (`tmp/glm-landing`, `tmp/glm-eval*`).

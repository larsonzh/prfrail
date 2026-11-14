# Slice validation report · P68 REMAINING-RULINGS

Slice id: ruling advice for the four R items (R-04 / R-06 / R-07 / R-10) plus the stage-progression seal confirmation (assessment-type `[SLICE]`, layer-2 governance side, tier M).

Start: 2026-11-12; controller = `DeepSeek V4.1 Flash (deepseek)` (default, no ADR registered).

Status: step 0 complete; stage ① awaiting start.

## 1. Executive summary

This is an assessment-type `[SLICE]` (layer-2 governance side, tier M), started 2026-11-12; controller `DeepSeek V4.1 Flash (deepseek)` (default, no ADR registered).
Pipeline trimmed to 0 → ① → ③ → ④ → ⑤ → ⑦ → ⑧ → ⑨ → ⑧c → ⑩: no ② implementation, no ③ code tests, ⑥ independent scan skipped.
Five artefacts: A stage-progression sealed (stage 1 effective / stage 2 not applicable); B option a; C option c; D option a plus gap registration; E option a plus c's pure-recording form.
Batch assist: GLM-5.3-Flash dispatched 2 times (2418 / 1428 chars), recorded in `02-batch/glm-batch1.txt` and `glm-batch2.txt`.
Review chain: ③ `CONSISTENCY: PASS` (34 file:line items, 33 first-pass hits plus 1 Low line-number fix); ④ `SELFTEST: PASS 737/737`, three scopes `TOTAL_FAIL=0`; ⑤ r1 `PASS` (1 Low remediated); ⑥ skip confirmed by ⑦; ⑦ r1 `FINAL REVIEW: PASS`, 0 findings.
Counts: R = 1, D = 0; the final rulings on the four R items are made by the user in the same turn after ⑩ or in batches.

- Basis for R = 1 (§9.2 counting unit "remedy re-run after a finding"): ③ found 1 line-number deviation (`p57-closeout-assessment.md:64` → `:69`) → fixed → re-ran the 34-anchor machine check once, counted as 1; the ⑤ Low adoption did not trigger a ring re-run and is not counted.
- Anomaly registration (Low aesthetic, accepted as is): `04-reviews/07-final-r1.md` carries 3 coexisting archive notes at its head (same meaning, no real defect, does not trip G2 / G7-f / G8-b); cleanup is deferred to v55 generation and the evidence pack is not modified in this slice.
- Post-⑧c evidence-window exception (§5.1 path 3, 2026-11-12): after the ⑩ check a §1 registration change required re-taking the ⑨ readings, so `tmp/P68-REMAINING-RULINGS/` was rebuilt after registration and precisely cleaned up with a terminal re-read afterwards.
- Temp-artefact D = 1 (recorded per the user's 2026-11-12 ruling): before ⑩ r3 the tree readings were written to `tmp\p68g.txt` (tmp root, not the declared `tmp/P68-REMAINING-RULINGS/`) and deleted immediately; the read-back showed tmp top level = `gates-cache,.gitkeep` and non-cache files = 1 (`.gitkeep`), zero residue; D < 3, so the slice does not stop.
- OB-11 exception (user ruling, 2026-11-12): the full R-07 advice is a substantive ① supplement that does not meet the five OB-11 elements but contains no rule-semantics change => recorded as **ADR-042** under the OB-37 meta-rule (exceeds the literal standard, a direct user ruling, not a precedent), with **no re-run of ⑤/⑦**.

## 2. Stage-progression seal confirmation

**Verdict: sealed** (stage 1 is the effective acceptance standard; stage 2 is "not applicable" under the user's A1 ruling of 2026-11-09).

| Anchor | Key text |
|---|---|
| `docs/validation/s1-exit.md:43` | AT-23 "stage-1 conditional closure reached ... stage 2 not applicable, user 2026-11-09 A1 ruling" |
| `docs/validation/s1-exit.md:52` | S-G1 basis = AT-23 stage 1; "must not claim S1 Exit passed as a whole" |
| `docs/TEST_STRATEGY.md:65` | AT-23 acceptance conclusion (F2 fix update; stage 2 not applicable / A1 ruling) |
| `docs/validation/r08-r09-ruling.md:103` | §4 stage progression; `:122` three upgrade conditions; `:125` condition 2 = explicit user ruling, the controller must not decide |
| `docs/t027/REMAINING_SLICES.md:93` | P61 table row "PARTIAL closeout ... ADR-036"; section body `:1736` |
| `docs/t027/CHECKLIST_v54.md:65` | F2-FIX row (`57631a3`; AT-23 stage-1 conditional closure reached; ADR-039) |
| `docs/DEV_PLAN.md:175` | F2-FIX closeout => "AT-23 stage-1 conditional closure reached" |

**Later slices need not re-examine the stage progression; no self-imposed tightening to stage-2 probe/verification standards.** (The only upgrade entry = a user ruling plus re-running the key scenarios.)

The stage-progression framework (Stage vs Layer / scope R-08.1+R-08.2 / the three upgrade conditions) lives in user-side material (v5 handover document §3); in this repository s1-exit / TEST_STRATEGY / ledger P61 only carry the sealed conclusion, so this report does not restate the framework and records a pointer only.

## 3. R-04 advice (DR-13 "fixed + observed" wording)

### 3.1 Current-state anchors

- `docs/t027/REMAINING_SLICES.md:1950`: DR-13 "fixed (2026-11-01) ... (4) re-open if the same assertion string recurs".
- Same file `:1405` DR13-FIX section ✅; `:1413` dependency row (recurrence runs `41253798647` and `41257990185`).
- No new recurrence across the four pending tables: v51:113 / v52:119 / v53:131 / v54:144; `CHECKLIST_v54.md:150` "kept pending".

### 3.2 Candidate options

- **Option a (close it as definitional wording, keeping reopen condition 4)** - budget 0 (pending-table row only); time <0.5h; risk low (guard against "closed" being read as permanently closed).
- **Option b (turn it into a long-term observation item and drop it from the pending table)** - budget 0; time <0.5h; risk medium (the reopen trigger loses its granted authorization and would need a new clause).
- **Option c (run a recurrence probe and push the evidence)** - budget 1 per-run-declared probe plus 1 two-leg CI; time 1-2h; risk medium (low CI recurrence rate, no statistical power, adds flake noise).

### 3.3 Advice

- Recommend **option a**: the D1/D2 fix, the counterexample regression (red then green) and D3 closing W1 already form a closed loop, and four versions (v51-v54) through P67 show zero recurrence; option c has no discriminating power. **The final ruling is to be made by the user in the same turn after ⑩ or in batches.**
- **User ruling (2026-11-12, same turn): option a.**

## 4. R-06 advice (OB-55 ownership and the T027-B cadence)

### 4.1 Current-state anchors

- `docs/DELIVERY_DIRECTIVE.md:1091`: OB-55 disposition "**under discussion** - to be assessed at stage ① when segment B is kicked off (standalone slice / land the `Sandbox` interface early / keep it in S2)".
- `docs/DEV_PLAN.md:151` kickoff registration; `:153` R-2 = OB-55 ownership pending R-06; `:161`/`:163` P37/P38 baselines; `:177`-`:185` P39-P44 (`:185` P44 end-to-end wiring deferred to a later slice).
- `docs/validation/p57-closeout-assessment.md:69` §3.3 (R-2 row); `:82`-`:89` §4.2 (`:89` S-G8 depends on T027-B).

### 4.2 Candidate options

- **Option a (standalone slice to land segment B)** - budget tier M (⑦ 1+2), about 3 review units; time 3-5 slices; risk high (idles without an admitted candidate and triggers an AT-23 re-assessment).
- **Option b (land the platform-neutral `Sandbox` interface early)** - budget tier S/M; time 1-2 slices; risk high (conflicts with B2 §8 item 8's dead-code prohibition, which must be fixed first).
- **Option c (keep it in S2)** - budget 0; time 0; risk medium (S-G1/S-G8 stay BLOCKED but the risk set stays controllable).

### 4.3 Advice

- Recommend **option c**: no usable candidate plus the dead-code prohibition plus S1 having no remaining time window form a triple constraint; if the goal becomes unlocking S-G1, switching to option a requires a newly admitted candidate as the precondition. **The final ruling is to be made by the user in the same turn after ⑩ or in batches.**
- **User ruling (2026-11-12, same turn): option b** (close the P44 end-to-end wiring R-1 first, then start T027-B); **the R-1 work itself is carried by a later slice, P68 does not advance it**.

## 5. R-07 advice (MVA, the S1 carve-out and narrowing the S2 commitment)

### 5.1 Current-state anchors

- Source text (in-repo): `docs/validation/strategic-planning-2026-10-31.md` §2.1 (`:17`-`:19`) and §4 item 2 (`:47`); `evidence/STRATEGIC-PLANNING-2026-10-31/01-product-planning/review.md` Q1 (`:6`-`:11`) / Q2 (`:13`) / Q3 (`:22`-`:29`); `03-synthesis-astra/review.md` conflict table (`:21` "product subtraction may be submitted for approval, but B's opinion cannot replace approval") and Q4 (`:39`-`:53`).
- Scope baseline: `docs/S1_SUPPORT_MATRIX.md` (visible row = degraded guarantee, not required by the core offline CLI); `docs/PRODUCT_REQUIREMENTS.md:71` (REQ-009) / `:82` (REQ-020) / `:69` (REQ-007) / `:83` (REQ-021) / `:81` (REQ-019).
- Technical premises (details in `01-arch/r07-refined.md` §0): A is currently "output only when finished", with no per-task/per-step lines and no `--verbose`/`--progress`; B redirects the AI child stdout/stderr to `logs/` files only, has no heartbeat, and `--json` prints one summary at the end; C events are structured and `events/state-events.jsonl` is the single total-order source, but there is no subscription API; D the visible channel's conversation view is rendered by the Copilot Chat panel, ProofRail neither renders nor treats it as authoritative.
- Additional: AT-13 has not passed (`TEST_STRATEGY.md:69`; hard gaps N-1 Windows-side verifier independence and N-6 toolchain uniformity); freezing S2 languages does not affect S1 (REQ-020's S1 only requires C+Go; generic/C/Go already verified).

### 5.2 Candidate options

- **Option a (approve everything)**: MVA definition + three S1 carve-outs (REQ-009 full TUI / REQ-019 self-hosting and US6/AT-13 / REQ-007,021 black-box items and US11/AT-24) + the S2 language freeze.
- **Option b (partial approval, 4 combinations)**: b1 = change only the REQ stage-ownership column (X-style deferral) + narrow the S2 language commitment, leaving the MVA acceptance composite alone, deleting no REQ and leaving self-hosting alone; b2 = b1 + approve MVA; b3 = S2 freeze + TUI carve-out only; b4 = b2 + self-hosting carve-out (= a in two steps).
- **Option c (defer)**: approve nothing, keep the original S1/S2 commitments.

### 5.3 Impact (budget / time / risk) and post-approval write-backs

- a / b2 / b4: one scope decision plus one write-back slice (tier M, ⑦ 1+2), 1-2 slices; risk high (MVA changes the acceptance baseline; b4 pushes the AT-13 gaps outward).
- b1: fits tier S (⑦ 1+1), 0.5-1 slice; risk medium-low (carves out the view layer and commitment surface; the S1 core loop and AT-13 stay put).
- b3: tier S, 0.5 slice; risk low (smallest write-back surface, but S1 Exit pressure stays). c: zero cost, risk medium (the problem moves later).
- Write-back list (**outside P68 scope**, needs its own slice): `PRODUCT_REQUIREMENTS` (stage ownership of the carve-outs + S2 narrowing registration), `TEST_STRATEGY` (AT item stage ownership), `S1_SUPPORT_MATRIX` (S1 scope baseline), `s1-exit` (gate applicability), `DEV_PLAN` (judged per §11.2); `CHECKLIST_v54` unchanged (inherited by v55).

### 5.4 Advice

- Recommend **option b1**: the availability cost of carving out the TUI can be cheaply mitigated by "the event stream already exists on disk + documenting the tail guidance" (no full TUI needed); the black-box carve-out matches ADR-012's degraded guarantee; the S2 freeze does not affect S1; self-hosting is not carved out (the N-1/N-6 gaps should not lose acceptance pressure); MVA waits until S1 Exit closes. **The final ruling is to be made by the user in the same turn after ⑩ or in batches.**
- **User ruling (2026-11-12, same turn): option b1 plus conditions** - X-style deferral; condition 1 = the TUI carve-out presupposes a `run --progress` lightweight stream switch inside S1 (carried by `S1-STREAM-PROGRESS`; P68 does not start it); condition 2 = the visible debugging guidance goes into `OPERATIONS` with three boundaries (carried by an S1 documentation slice).

- **P71 landing pointer (2026-11-13)**: **condition 2** of R-07 b1 (the visible debugging guidance plus three boundaries) has been landed by the `S1-SCOPE-ANNOTATE` slice; the report is `docs/validation/s1-scope-annotate{,_EN}.md`. The original lines above stay untouched (§11.1).

- **P72 landing pointer (2026-11-14)**: the S2 narrowing in R-07 b1 item 3 (the Q3 language freeze) has been landed by the `REQ020-S2-RESTRICTION` slice; the report is `docs/validation/req020-s2-restriction{,_EN}.md`. The original lines above stay untouched (§11.1).

## 6. R-10 advice (overall budget and the per-slice model caps)

### 6.1 Current-state anchors

- `docs/DELIVERY_DIRECTIVE.md:757` §7.5 hard-cap table (architect <=1 / pre-review <=1 with a +1 re-review / ⑥ <=3 / ⑦ tier M 1+2 / blind review +1 separately); `:775` cost gradient (Flash `0.15/0.6`, V4 Pro `0.66/1.98`); `:1229` §13.1.8 controller-model policy.
- `docs/t027/CHECKLIST_v54.md:147`: "R-10 ... pending (blocks scheduling)".
- **Gap declaration**: no overall budget figure exists in the repository (per-slice caps and unit prices exist; the total pot does not).

### 6.2 Candidate options and boundaries

- **Option a (keep the §7.5 caps as they are; no directive edit)** - budget 0; time 0; risk low-medium (no hard total ceiling, so overspend is found late, but the ⑩ per-slice metering plus the "over the cap => escalate to the user" rule backstops it) => **in-bounds**.
- **Option b (add a "total budget" row to §7.5)** - touches the §7.5 body => a directive revision; budget governance tier M/L; time 1-2 slices; risk high (creates a split standard against slices already in flight). **This slice advises only and does not land it.**
- **Option c (tiered recording without a directive edit)** - only append a "total-budget check" line to the ⑩ receipt, with no criterion: budget 0; time 0; risk low. Recording only, no criterion and not gate-checkable => a companion of option a (in-bounds); once written into §7.5 or turned into a checkable criterion it becomes option b (out-of-bounds).
- **Option d (re-estimate the per-role caps slice by slice)** - touches §7.5 values; budget a governance slice; time 1 slice; risk medium (old and new standards coexist). **This slice advises only and does not land it.**

### 6.3 Advice

- Recommend **option a** (keep the caps; no directive edit) plus **option c's pure-recording form** (not written into the directive; the landing form of c still requires a user ruling). b/d are marked as requiring a user ruling and are not landed by this slice; this slice does not trigger a stop-escalation. **The final ruling is to be made by the user in the same turn after ⑩ or in batches.**
- **User ruling (2026-11-12, same turn): option a** (keep the §7.5 caps; no directive edit).


## 7. Known boundaries

- **Correction (2026-11-12, raised by the user)**: the substance of the MVA / S2 commitments IS in this repository (`strategic-planning-2026-10-31.md` §2.1 `:17`-`:19` and §4 item 2 `:47`; both reviews' Q1/Q3 and the conflict table); only the v26 snapshot is user-side. The full R-07 advice has been re-issued from the in-repo source text (`01-arch/r07-refined.md`).
- Technical-debt observation (registered only, not fixed, no OB opened): `prfrail run` has no per-task/per-step output, no `--verbose`/`--progress` and no heartbeat; the AI child stdout/stderr go to `logs/` files only; `events/state-events.jsonl` is a structured single total-order source but has no subscription API, and README / OPERATIONS give no run-time tail guidance. Whether to open a slice for this is the user's call.
- The split of S2 TUI work (view layer vs content layer) has **no documentary basis** in this repository, so this report does not estimate it.
- No overall budget figure exists in the repository; budget authority stays user-side and no figure is quoted here.
- The full stage-progression framework is user-side; the repository carries only the sealed conclusion and its four landing anchors.
- This slice gives advice, not rulings: all four options bind nothing until the user rules on them.
- The advice reflects documents as of 2026-11-12; if facts change after the user's ruling, re-assessment is required before citing it.

## 8. Explicitly not done

- No change to product code / CONTRACTS / schema / fixtures / CI; no change to the directive body; no change to the v50–v54 snapshots; no ruling made on the user's behalf; no commit / push; no new slice started; v55 not generated.

## 8a. ⑧a Self-check

- ① Cost-table mechanical audit: none (not applicable) - no cost ledger file in this slice; the payload list is in `MANIFEST.md`.
- ② "Disposed" `file:line`: yes - 34 items machine-checked, 33 first-pass hits plus 1 fix (`p57-closeout-assessment.md:64` → `:69`), see `03-test/consistency-r1.md`.
- ③ Legacy-wording grep: none (not applicable) - no directive body or wording was changed, so no legacy grep was run.
- ④ ⑥/⑦ per-round Info / Low capture: yes - the single Low from ⑤ r1 was remediated (see `04-reviews/05-prereview-r1.md`); ⑥ skipped and ⑦ r1 had 0 findings, so no extra entries.

## 9. Commit and CI

- Main commit `e0d745c` (`docs: land the P68 four-R rulings, the full R-07 advice, and ADR-042`; pushed to `origin` main: `365a930..e0d745c`; no gitee push).
- CI run `43177025096` (head `e0d745c`) green on the first attempt on both legs:
  - windows (5m5s): `SELFTEST: PASS 737/737`; `GATE REPORT scope=ci changed=21`; `TOTAL_FAIL=0`.
  - ubuntu (3m28s): `SELFTEST: PASS 737/737`; `GATE REPORT scope=ci changed=21`; `TOTAL_FAIL=0`.
- This section is back-filled by this slice's controller in a docs-only commit after CI went green on both legs (§11.2 ledger completion-row rule: commit hash and CI run numbers go into the report only); that back-fill commit's own run is observation-only and no further commit is made for it.

- ⑥ skipped: this is an assessment-type documentation-only slice that touches no code and none of the four semantic classes in §7.2.1 ③ ⇒ it may be skipped (pending ⑦ confirmation); if ⑦ rules otherwise, ⑥ is run and this line is replaced by a `sixstage` fence.

```artifacts
docs/validation/evidence/P68-REMAINING-RULINGS/00-briefing.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/MANIFEST.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/SHA256SUMS.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/freedom-list.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/input-list.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/01-arch/design-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/01-arch/r07-refined.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/02-batch/glm-batch1.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/02-batch/glm-batch2.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/03-test/consistency-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/06-gates/gates-r1.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/06-gates/gates-r2.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/06-gates/step0-selftest-r1.txt	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
```

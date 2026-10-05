# Slice `DIRECTIVE-GLM-EVAL` · Validation report

| Item | Value |
|---|---|
| Slice | `DIRECTIVE-GLM-EVAL` (`[SLICE]`, carried over from the previous session to ④) |
| Goal | Evaluate the quality and cost of `GLM-5.3 (glm)` at the **⑤ pre-review position** (against `DeepSeek V4 Pro (deepseek)`), and evaluate the quality and deliverable usability of `GLM-5.3-Flash (glm)` on 4 classes of batch auxiliary tasks |
| Boundary | **Evaluation only, not admitted to the whitelist**; the whitelist and role decisions are discussed separately after this slice |
| Authorization boundary | commit no / push no / gitee no / probes 0 (GLM calls are paid calls, accounted separately from probes) |
| Report landing | `docs/validation/directive-glm-eval.md` (this file, CN authority) + the `_EN` same-position mirror |

> **Nature statement**: this slice **does not change any rule body**, **does not change contracts / schemas / fixtures**, **does not change the role set**, and **does not extend the whitelist**.
> Deliverables are: this report pair, the evaluation dataset, the frozen evaluation evidence pack, the freedom list, and the ledger row. Both GLM models are **not** in the directive
> §2.2 whitelist; the calls this slice makes to them are a **user-authorized temporary evaluation**, independently traced per §1.7.4 (project marker `DIRECTIVE-GLM-EVAL`),
> and **must not serve as a precedent** — whether GLM is formally adopted is discussed separately **after this slice**.

---

## 1. Change summary

| Category | Content |
|---|---|
| Added | this report pair `docs/validation/directive-glm-eval{,_EN}.md`; the evaluation dataset `docs/validation/directive-glm-eval-dataset.json`; the frozen evidence pack `docs/validation/evidence/glm-eval/` (`MANIFEST.md` + `SHA256SUMS.txt`); the freedom list `docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md` |
| Modified | the ledger completion row in `docs/t027/REMAINING_SLICES{,_EN}.md` |
| Zero change | `docs/DELIVERY_DIRECTIVE{,_EN}.md`, `docs/CONTRACTS{,_EN}.md`, `schemas/**`, contract fixtures, `.github/agents/**`, `.github/workflows/**`, `internal/**`, `tools/gates/**`, `go.mod`, `tools/gates/cjk-newwords.txt` |

**Substantive choices exercising execution freedom** (§1.7.3 trace; item by item in §10):

1. The ⑤ control leg's variant took "synthetic seeded defects" (candidate B); rationale in §3.1;
2. The scorer was fixed **after the experimental leg and before the control leg** (four generic rules, not special-cased for this output) — the report must give **both** the pre-fix and post-fix readings (§4.1);
3. The leg-artifact landing calibers are asymmetric (experimental leg captured **byte-for-byte** / control leg transcribed **word-for-word**) — truthfully recorded (§7);
4. The 4 batch legs were executed **concurrently** (the plan text says "executed in sequence") — truthfully recorded, not re-tested (§7);
5. The cost-reading caliber was downgraded from "per-call window readings" to the "**cumulative** readings provided by the user" — truthfully recorded (§5).

---

## 2. Evaluation purpose, scope and boundary

**Purpose**: to provide **in-sample** quality and cost evidence for "whether to bring the GLM series into the corresponding position of the directive review chain".

**Scope**: two evaluation faces.

- **Face one (quality comparison, n=1)**: place `GLM-5.3 (glm)` at the **⑤ pre-reviewer** position and run it side by side with the model used by the incumbent ⑤ carrier,
  `DeepSeek V4 Pro (deepseek)`, on the same input and the same prompt, scored against the **pre-registered key**.
- **Face two (batch auxiliary usability)**: `GLM-5.3-Flash (glm)` takes on 4 classes of batch auxiliary tasks (evidence-tree enumeration and statistics / long-document excerpting and number checking /
  bilingual-mirror language-side checking / archive-inventory drafting), scoring only **quality and deliverable usability**.

**Boundary (explicitly not done)**: not written into any whitelist, no role-set change, no change to the bodies of §2.2 / §2.4 / §5.2 / §7.5 / §1.5.1;
no multi-round repetition for statistical power (this slice is n=1); no token metering on the V4 Pro side (readings unavailable, see §5).

---

## 3. Evaluation design

### 3.1 Same-input construction (variant B: synthetic seeded defects)

**Variant rationale**: candidate A (using a real repository change set as input) was rejected — ① A's answer key is **already committed** (the validation report and
freedom list at commit `025158f` carry the full set of ⑤/⑦ findings and wording), and the review carrier can read files, so the soft constraint "must not read the repository" cannot be mechanically enforced ⇒ the comparison would inevitably be distorted;
② A's key is V4 Pro's own findings, so using it to score V4 Pro's hit rate carries a systematic bias (exactly the variable being measured); ③ A's key is incomplete.
Candidate B: within the experimental window **the repository contains no answer source at all** (the synthetic scenario shares no text with the real repository) ⇒ the contamination surface is zero.

### 3.2 Bundle frozen objects (`tmp/glm-eval/bundle/`, inlined byte-for-byte inside the task package, no paths given)

| File | Bytes | SHA-256 |
|---|---|---|
| `00-plan.md` (synthetic mini plan) | 1,184 | `70076661d86615d8dfe180aa5884ec9618338ec9496b3893900a749c092f36b8` |
| `01-clause-cn.txt` (synthetic contract text, self-made clause numbers X1.1/X6.3) | 2,346 | `d8fe4bc82df774535c47c39406c96fee1325c87f3153b5d72ee2a378ff03e525` |
| `02-diff.js` (synthetic change set, about 6.5 KB) | 6,451 | `c939bc9e756696c463571c981a7d3756c87563de913b9c4c58faebdc4d5ff405` |
| `03-prompt.txt` (task-package body, byte-identical across both legs, **model name not written**) | 1,597 | `3b7a3580edcecb6a92ccb42f8acf37b469283136cb0c8d8d1b2b24d0b54fd758` |
| **`input_hash`** (the four files concatenated in a fixed order) | 11,578 | `bbeaf2f7db40547d26b1a338608c0b9c1679fd7b08edbf349c0fbadaaa8008b6` |

The task-package body (`03-prompt.txt`) contains B.5's input / output / prohibition / conclusion-line contract, plus two hard constraints:
"do not read any workspace or external file, rely only on this message text; return the report as the final message".

### 3.3 Pre-registered answer key (S1–S8; landed before the two leg calls)

| Seed | Severity | Defect |
|---|---|---|
| S1 | Critical | When parsing SHA256SUMS fails / the whitelist is missing, `catch` swallows it and `return true` (fail-open, silently passing) |
| S2 | High | Package member names are compared to manifest names by direct string comparison, with no backslash-to-forward-slash normalization ⇒ on Windows the package is always judged red |
| S3 | High | On recompute it takes the manifest line's **member-name column** as the digest, and `if (recorded === recorded)` is always true ⇒ the digest comparison is a formality |
| S4 | Medium | The contract states that "the manifest has a member while the package lacks that file" is also a conflict; the implementation has no such branch |
| S5 | Medium | The contract states four sub-checks; the implementation lands only three branches (missing the rebuilt-archive-note check) |
| S6 | Medium | "Judge the change set only" is implemented as a whole-tree scan (conflicts with the contract; a historical false positive) |
| S7 | Low | The top comment says "eight classes" when it is actually four (stale comment) |
| S8 | No defect | Pure reordering / blank lines (expect 0 findings; reporting one counts as FP) |

**Key-not-retroactively-modified statement**: the SHA-256 of `key.json` is `6fccbe461286abae63e74e52f64d8130f4fa979bb6d27ba71651fad6c3339536`,
**byte-for-byte equal** to the dataset's `pre_registration.key_file_sha256` ⇒ not modified after freezing (recheck command `node tmp/persist.js`).

**Tension between the key and the rubric (truthfully recorded, both readings reported)**: the pre-registered key rates S4 / S5 / S6 as **Medium**; while the rubric sent to the model defines
**High as "an important branch missing, or a condition always true/always false"**. Re-judged by the rubric, S4 / S5 / S6 can all read as High ⇒ `sev_exact = 1.0`, `Q = 1.0`.
**Both readings coexist**: by the key `Q = 0.9571`; by the rubric `Q = 1.0`. **The key is not modified because of this tension** (proven in the paragraph above).

### 3.4 Scoring rubric and composite score (pre-registered)

- Hit rate: `recall_all = hits / 7`; `recall_high = hits(S1–S3) / 3`.
- False positives: a finding citing a nonexistent line, or an assertion contradicting the bundle, counts as `FP`; **extra true findings** are mechanically re-checked by the driver and listed separately, not counted as FP.
- Severity calibration: exact 1.0 / off by 1 tier 0.5 / off by ≥2 tiers 0; `sev_exact` is the mean over the 7 scored seeds.
- Conclusion line: the last line is exactly one of `PRE-REVIEW: PASS | PASS WITH FIXES | FINDINGS` (§4.1 closed set) ⇒ 1 / 0.
- **Composite score `Q = 0.4 × recall_high + 0.2 × recall_all + 0.2 × sev_exact + 0.2 × conclusion_line`**.
  `Q` is only a **descriptive quantity**; with n=1 no conclusion on model superiority may be drawn.

### 3.5 Contamination detection (mechanically executed)

Search both legs' outputs against the marker vocabulary (`OB-79` / `025158f` / `0c12abc` / `G8-[abcd]` / `DIRECTIVE-EVIDENCE-SCOPE` / `重建归档` /
`判据作用域排除` / `S[1-8]`) — must be **0 hits**; also check whether each `file:line` falls on a path inside the bundle.
**Measured**: both legs `PASS` (0 hits, all paths inside the bundle). **S8 region** (`02-diff.js:191–198`) judged "not a finding" by both legs ⇒ `s8_fp = 0`.

### 3.6 Freezing and pre-registration self-proof

`tmp/glm-eval/FROZEN-MANIFEST.txt` (73 entries) is the in-session frozen manifest; one recompute yields both the manifest and a re-verification of `input_hash` and `key.json`:
**`input_hash` recomputed value = pre-registered value; `key.json` hash = pre-registered value; the dataset has 8 rows and `pre_registration` is unchanged**.
The **authoritative hash** of the handover artifact `BRIEFING.md` is the **full hash** listed in that manifest:
`18867c95816abfa977c30eec31a74ef54ef17405b6c90e91caa6f05415df3ae8` (18,971 B).

---

## 4. Evaluation results

### 4.1 Face one · ⑤-position two-leg comparison (same input, same prompt; n=1)

| Item | ⑤-experimental | ⑤-control |
|---|---|---|
| `model_string` (qualified name, R2.2) | `GLM-5.3 (glm)` | `DeepSeek V4 Pro (deepseek)` |
| `quota_package` | `glm-trial-2000w` | `deepseek-n-a` |
| Carrier | `deep-reasoner` (read-only) | `deep-reasoner` (read-only) |
| Output file / bytes | `leg-glm.txt` / 13,234 | `leg-control.txt` / 6,591 |
| Output SHA-256 | `54c895f9…` | `2fa1d9ee…` |
| Landing caliber | raw tool capture, saved **byte-for-byte** | inlined tool return, transcribed **word-for-word** by the driver |
| Finding count (self-reported) | 12 (Critical 1 / High 6 / Medium 3 / Low 2) | 15 (Critical 1 / High 7 / Medium 5 / Low 2) |
| `recall_all` | 1.0 (7/7) | 1.0 (7/7) |
| `recall_high` | 1.0 (3/3) | 1.0 (3/3) |
| `sev_exact` | 0.7857 | 0.7857 |
| `conclusion_line` | 1 (`PRE-REVIEW: FINDINGS`) | 1 (`PRE-REVIEW: FINDINGS`) |
| `FP` | 0 | 0 |
| Abstention lines | 2 | 1 |
| `extra_true` (extra true findings, re-checked) | 10 | 8 |
| `s8_fp` | 0 | 0 |
| **`Q` (key caliber)** | **0.9571** | **0.9571** |
| Contamination detection | PASS | PASS |

- Severity item by item (both legs agree): S1 Critical→Critical(1.0), S2 High→High(1.0), S3 High→High(1.0), S7 Low→Low(1.0);
  S4 / S5 / S6: the key records **Medium**, both legs report **High**(0.5) ⇒ `sev_exact = 5.5 / 7 = 0.7857`.
- **Re-judged by the rubric** (§3.3 tension): `sev_exact = 1.0`, `Q = 1.0`. **Both readings reported**.
- **Scorer pre-fix / post-fix two readings (experimental leg)**: pre-fix `Q = 0.9048 / FP = 1` → post-fix `Q = 0.9571 / FP = 0`.
  The fix consists of four **generic** rules (range citation / comma citation / abstention detection / de-duplication), **not special-cased for this output**; `key.json` unchanged.
  The control leg uses only the post-fix version.
- **Usable as a directional observation (not tokens)**: the experimental leg's output volume is about **2.0×** the control's (13,234 vs 6,591 bytes).
  **Character count is not token count and cannot be converted**, so no claim about cost may be based on it.

### 4.2 Face two · 4 batch legs (`GLM-5.3-Flash (glm)`, package `glm-flash-realname-newuser`)

| Leg | Task | Result | Verdict |
|---|---|---|---|
| BATCH-a | Evidence-tree enumeration / statistics (40 synthetic manifest records) | 40 ✓, extension distribution ✓, 209,020 bytes ✓, top level 6 entries ✓; **top-5 wrong** (3/5 hits: 9817, 9610, 8840; missed 9403, 9196; extra 8633, 8426) | **Partial failure** (sorting subtask) |
| BATCH-b | Long-document excerpting + number checking (the v1.24 report pair) | 328/328 ✓, scope 39/2/4 ✓, cost vector ✓, commit `025158f` ✓, run `37278691747` ✓, CN/EN item-by-item `ALL MATCH` ✓ | Usable (all correct) |
| BATCH-c | Bilingual-mirror language-side checking | 14 language-side opinions; kept its boundary (no pass / fail verdict, no structural involvement; correctly refused to treat `driver` / `To be discussed` as terminology drift) | Usable (review input) |
| BATCH-d | MANIFEST + SHA256SUMS draft | 6 entries / 4437 bytes ✓, 6 hashes and bytes ✓, `archive/` directory-level stripping + lexicographic order ✓, two-space format ✓; minor gap: the unregistered table's last column name | Usable draft |

**Conclusion (in-sample, n=1, not to be extrapolated)**: of the four face-two legs, **three legs' deliverables are usable** and **one leg (BATCH-a) partially fails on the sorting subtask**.
Batch legs are scored only on quality and deliverable usability and **do not enter** the ⑤-position quality score.

---

## 5. Cost and metering (one-sided; the relative criterion does not hold)

**Metering-caliber downgrade statement (truthfully recorded)**: within the call window the console **pre / post** readings could not be obtained (the operator was not present, no window-exclusive confirmation),
and the extension log and storage likewise have **no usage field** ⇒ a per-call token count cannot be attributed. Later the user **provided cumulative readings** in the session:

| Model | Quota package | Cumulative reading | Calls in this slice |
|---|---|---|---|
| `glm-5.3` | `glm-trial-2000w` | **161,863 tokens** | 2 (qualified-name probe 1 + ⑤ experimental leg 1) |
| `glm-5.3-flash` | `glm-flash-realname-newuser` | **424,690 tokens** | 4 (batch 4 legs) |
| `deepseek-v4-pro` | `deepseek-n-a` | **unavailable** | 1 (⑤ control leg) |

- **Per-row package labeling** (dataset 8 rows): probe rejected / probe accepted / ⑤ experimental leg = `glm-trial-2000w`; control leg = `deepseek-n-a`; batch 4 legs = `glm-flash-realname-newuser`.
- `metering_mode`: GLM side = `user-provided-cumulative`; V4 Pro side = `unavailable`.
- **No estimates filled in, no token conversion done**. The original downgrade basis is kept in the dataset's `metering_note_at_call_time`.
- **One-sided absolute magnitude (caliber correction, ⑤ pre-review Medium #3)**: the magnitude bound of a single ⑤ experimental leg is **≈ 161,863 tokens**
  (the user-provided cumulative reading; the **probe leg is a word-level minimal call** and the user's early test messages were **not stripped**).
  ⇒ **Must not** take `161,863 ÷ 2 ≈ 80,931` as the per-call magnitude of the experimental leg — that mean includes the minimal probe leg and would **systematically underestimate**.
- **Relative cost criterion does not hold**: the V4 Pro side's reading is unavailable ⇒ **no** relative conclusion on "which of GLM and V4 Pro is cheaper" may be given; that comparison is **suspended** as input to a later slice.
- **Newly stated truthfully**: the measured consumption is **significantly higher than the driver's estimate** (estimate 10k–20k / call magnitude; measured about 4–8×). The cause was not checked
  (possibly CoT billing, system-prompt overhead, or the user's early test messages not being stripped) — **recorded only, not explained**.
- **Caliber correction of the handover artifact's cost-mandatory item (⑤ re-review NL-2, disposed)**: item 3 of `BRIEFING.md`'s cost-mandatory list originally said "per-call pre-review magnitude ≈ 80k"; it is **superseded by the correction in this section** — written as the **bound** 161,863, with the ÷2 mean caliber explicitly **forbidden**; the handover artifact itself is not rewritten (the user ruled "take `FROZEN-MANIFEST` as the sole authority").
- **Measured cost of ⑦'s three readings (this slice's data point)**: about 3 hours of wall clock and about 600 credits (about 200 credits per reading), the same order of magnitude as the previous slice `DIRECTIVE-EVIDENCE-SCOPE` ⇒ registered as **OB-84** (§12.4), with a gap row "no budget model for ⑦'s three readings" opened in §12.5 for the budget reference of later evaluation / governance slices.

---

## 6. Falsifiability and mechanical verification

| Verification item | Assertion | Result |
|---|---|---|
| Frozen-pack recompute | `node tmp/persist.js` recomputes the manifest and compares all hashes | **match=true** (`input_hash`, `key.json`, dataset 8 rows) |
| Key freezing | `key.json` hash = `pre_registration.key_file_sha256` | **equal** |
| Bundle untouched | the four files, file by file, sha256 equals the pre-registration entry by entry | **equal** |
| Scoring / hash / contamination scripts | the ③ independent adversarial self-tester (`selftest-report.txt`, **179 assertions**, including T-GE-13/14) | `SELFTEST TOTAL=179 PASS=179 FAIL=0`, against the **final** scorer `score.js = 2c4052fa…` (that output was **regenerated** in the ⑤ re-review NL-1 remediation, see the corresponding record in §7) |
| Mutation falsification | four mutations — scoring denominator / range short-circuit / abstention sectioning / contamination boundary ⇒ assertions turn red ⇒ restore byte-for-byte | **all four turned red and then green again** (sha256 identical before and after) |
| Contamination detection | both legs' outputs: 0 hits on the marker vocabulary + all paths inside the bundle | **PASS** (phrased as "no contamination signal detected"; **does not claim** "proven not to have read the repository") |
| Dataset zero-leakage | the dataset must not contain seed numbers / key body text | the resident assertion passes |
| Carrier self-report | **not trusted** — the sub-agent's self-reported model and tier (§2.3) | call parameters and the operator UI record govern |

- **③ report body vs frozen self-test output version gap (⑤ pre-review Medium #1, remediated)**: the body of `TEST-REPORT-DIRECTIVE-GLM-EVAL.txt`
  stays at the **second round (154 assertions, scorer `b4ad39…`)** caliber, while the **final** scorer is `2c4052fa…` and the frozen self-test output is **179/179**.
  ⇒ A **third round (final-version re-verification)** subsection was appended to `TEST-REPORT`. The old scorer used for the pre-correction experimental-leg reading (`Q=0.9048 / FP=1`)
  **is no longer on disk**; only its "sha256 identical before and after mutation" chain remains (`TEST-REPORT` §8.5) ⇒ truthfully registered as a **reproducibility limitation**, without fabricating artifacts.

- **⑤ pre-review round-1 Low dispositions (on record, non-blocking)**: ① two **typos** in `plan-v1.md` (§5 says "bundle/ six files" when it is four;
  §1's contamination-detection command path says `tmp\glm-eval\out\*.md` when the actual artifacts are in `tmp/glm-eval-out/*.txt`, and the implementation runs equivalently via `contam.js`)
  — **the frozen plan is not modified** (it is the pre-registration basis; a later change would harm pre-registration integrity); instead it is registered in this report **as an erratum**;
  ② the dataset's `call_no` declared domain is `1..N`; this slice has **8 rows** (probe 2 + ⑤ 2 + batch 4, per row), and the difference from ① the plan's §5 "6 call faces" caliber is explicitly stated here;
  ③ the inputs of batch b / c / d are **inlined prompts, not landed on disk** ⇒ those three legs' deliverables are **not reproducible for audit** (`input_hash: null` is the registration of that fact), truthfully declared;
  ④ the re-check record of `extra_true` (10 / 8): independently sampled by the ⑤ pre-review at F5 / F9 / F10, confirmed as real defects in the synthetic diff (not miscounted).

- **`plan-v1.md` correction note (⑦ blind review M2 / Low-4, §11.1 form)**: see `tmp/glm-eval/PLAN-ERRATA.txt`.
  The **authoritative caliber** for severity tiers is `key.json`'s `severity_rank` (`Critical` 3 / `High` 2 / `Medium` 1 / `Low` 0),
  which is **tier-distance equivalent** to the three-tier enumeration in the plan's §1 parentheses (in both `High − Medium = Medium − Low = 1`; `Critical` is used only for `S1`,
  hit exactly by both legs) ⇒ **no functional drift**; the **correct command form** for contamination detection is also in that note.
  **The plan body is not written back** — reason: `tmp/persist.js` **re-extracts that file verbatim** from the chat transcript, so a manual append would be overwritten at the next recompute
  (measured: after an append, re-running returns to the 12,860 B transcript original) ⇒ manual changes are **not persistent**.

---

## 7. Anomalies and fallback record

| # | Event | Characterization |
|---|---|---|
| 1 | **Metering not closed**: no pre-call reading in the window, no usage field in the extension log | The cost dimension was first downgraded to `unavailable`, then the user supplied cumulative readings ⇒ changed to a **one-sided** absolute magnitude; the relative criterion stays suspended (§5) |
| 2 | **The 4 batch legs ran concurrently** (the plan says "executed in sequence") | Protocol deviation, truthfully recorded, **not re-tested**; per-call durations are indistinguishable; it does not affect the validity of the conclusion (batch legs are scored only on quality and usability) |
| 3 | **Asymmetric leg-artifact landing caliber** (experimental leg captured byte-for-byte / control leg transcribed word-for-word) | Truthfully recorded; the two hashes use different calibers and **must not** be treated as the same collection path |
| 4 | **The scorer was fixed after the experimental leg and before the control leg** | A **reasonable** fix (four generic rules); the report lists **both** pre-fix and post-fix readings (§4.1) |
| 5 | **Key / rubric tension** (S4 / S5 / S6: key Medium vs rubric High) | Both readings reported + a statement that the key was **not retroactively modified** (§3.3) |
| 6 | **① the plan was previously not landed** | Corrected in the previous session: `tmp/glm-eval/plan-v1.md` + `PLAN-PROVENANCE.txt` landed (§3.6) |
| 7 | **Encoding event (external save between sessions)** | The handover artifact `BRIEFING.md` measured 14,687 B / **no BOM** at the close of the previous session; before appending in this session it measured **14,690 B / with BOM** (`EF BB BF`) — the difference is exactly the **+3 B BOM**, with **byte-identical content**. The 14,687 B version was rebuilt from the chat transcript and verified byte-for-byte, its SHA-256 **exactly matching** the `d8463ef5…` recorded in the previous session; the evidence chain is `tmp/recon-briefing.js` + `tmp/fix-bom-and-append.js`. On landing it was **restored** to the `tmp/` process-artifact caliber (no BOM + LF). `PROBE-NOTES.txt` in the same batch of prompts was re-checked and **unaffected**. |
| 8 | `runSubagent` intermittently produces `Agent error: no response was returned` | Handling procedure: **check disk / transcript first**, then decide whether to resend, avoiding duplicate paid calls |
| 9 | **The report draft was produced before ⑤ (a special shape of evaluation slices)** | In the directive §5.1, ⑧a sits after ⑦; this slice is an **evaluation slice**, and the review object is the **evaluation execution and evidence** (the user's same-round positioning), so ⑤⑥⑦ **all do not review the report body**. Hence the report draft produced before ⑤ **is not** a product of the ⑧a step, nor a review object of ⑤–⑦; its final body is finalized by ⑧a (including the `_EN` mirror). |
| 10 | **Gate red (truthfully recorded, not silently skipped)** | After the report draft landed, `--scope=tree` first reported **G4b** (8 new CJK characters) and **G6-4** (two "line count"-type expiring literals), **both fixed in place** (rewording, caliber change; `cjk-newwords.txt` **not extended**); afterwards **G7-b** (`repo` artifacts not in the index) and **G7-a** (missing `artifacts` block) went red in turn — per the user's same-round caliber these are **not expected-ignore items**, and **both were closed after ⑧a landed the artifacts plus `git add`** (§9, §16). |
| 11 | **Self-test report regeneration (a side effect of the ⑤ re-review NL-1 remediation)** | On re-run, ③ found 3 **scaffolding-period** assertions that had expired by design (`calls[]` was already filled to 8 rows by the execution step; new handover documents under `tmp/glm-eval/`) and 1 domain assertion that needed to be synced with the dataset schema remediation ⇒ ③ updated them to the **post-execution true state** (`T-GE-09.known` **tightened** from a wildcard to a 58-item closed set) and re-ran to `179/179`. Side effect: `selftest-report.txt` was **regenerated**, and its **scaffolding-period original bytes are no longer reproducible** (an **irreversible evidence loss**, truthfully recorded, not fabricated). |
| 12 | **⑦ blind review M1 — the dataset's column domains and null values are not self-consistent** | Second-round fix: the `domain` of `input_hash` / `output_hash` and the `type` of `duration_s` were made **nullable** with the reason for unavailability written in (the same class as `metering_mode` / `measured_tokens` / `conclusion_line` already fixed in the ⑤ round); `.json` re-validation passed. |
| 13 | **⑦ blind review M2 — dual caliber for severity tiers** | A correction note landed at `tmp/glm-eval/PLAN-ERRATA.txt` (§11.1 "add a correction note only, do not rewrite the body"): the authoritative caliber = `key.json`'s `severity_rank`; the two are **tier-distance equivalent**, and `Critical` is used only for `S1` ⇒ **no functional drift**. **Not written back to `plan-v1.md`** (`persist.js` regenerates that file verbatim from the transcript, so manual edits do not persist). |
| 14 | **⑦ blind review M3 — self-test report overwrite causes evidence loss** | Truthfully recorded as an **irreversible evidence loss**; and §12 "evidence retention discipline" was appended to `TEST-REPORT` (snapshot first, re-run second; one state per file; declare when not reproducible), **and executed immediately**: the current `selftest-report.txt` has been snapshotted as `tmp/glm-eval-out/selftest-report-20261005.txt` (both hashes identical). |
| 15 | **⑦ blind review Low-4 — typo in the contamination command path** | Folded into `PLAN-ERRATA.txt` correction two (correct form = `contam.js` run against both legs' `.txt`, or `Select-String` over `leg-*.txt`). |
| 16 | **⑦ final review H1 — the frozen manifest did not register `PLAN-ERRATA.txt`** | Closed: after re-running `node tmp/persist.js`, `PLAN-ERRATA.txt` (4,320 B / `a3c08355…`) and the snapshot `selftest-report-20261005.txt` (26,531 B / `7b06d9d7…`) were both registered into `FROZEN-MANIFEST.txt`. |
| 17 | **⑦ final review M2 — the plan's `call_no` row domain does not match the execution face** | Folded into `PLAN-ERRATA.txt` **correction three**: 8 rows = 2 probes + ⑤ 2 + batch 4, with the source of the difference and the current authoritative caliber given. |
| 18 | **⑦ final review M3 — severity equivalence did not state its applicability boundary** | An "applicability boundary" paragraph was added to `PLAN-ERRATA.txt` **correction one**: the equivalence **holds only for this sample's readings**; general scoring always follows `key.json`'s `severity_rank` (with the four-tier mapping attached). |
| 19 | **§1.7.4 independent record (OB-83)** | The two GLM models are absent from the §2.2 whitelist, so calling them is a **temporary authorisation by the user in the same round** ⇒ per the five elements of §1.7.4 it is **independently registered as OB-83** in §12.4 of `docs/DELIVERY_DIRECTIVE{,_EN}.md` (ruling party / date / reason / basis complete, with the truthful annotation "exceeds the literal standard of §2.4, is a direct ruling by the user, must not be used as a precedent"); the directive rises to **v1.25** with this slice, and this row is an in-slice side product, not a standalone directive-governance slice. |
| 20 | **The metering procedure was built but never used** | The "pre / post / recheck three readings" procedure in `tmp/glm-eval/metering/README.md` **could not be executed inside the call window** (the operator was away), and the actual path was "the user-provided cumulative reading" ⇒ the procedure was **never actually walked**; it is kept as the **preferred option** for later evaluation slices, with a recommendation to rewrite it by availability (explicitly writing the path and consequence of "with no pre-reading, fall back to `user-provided-cumulative`"). |

**Fallback count**: the driver's interventions for localization / fixing / re-running in this slice = ≤ 3 (did not trigger §9.2's "> 3 times ⇒ mark high risk").

---

## 8. Environment and tiers

| Item | Value |
|---|---|
| Driver | `DeepSeek V4.1 Flash (deepseek)`, tier High |
| ① architect / ⑤ pre-reviewer | `DeepSeek V4 Pro (deepseek)`, tier Max |
| Product layer (implementation / testing / documentation) | `DeepSeek V4.1 Flash (deepseek)` |
| ⑥ independent scan | `MAI-Code-1.1-Flash (copilot)` |
| ⑦ independent final review | `GPT-5.3-Codex (copilot)`, Extra High (**not downgraded**) |
| Evaluation objects (not directive roles) | `GLM-5.3 (glm)` / `GLM-5.3-Flash (glm)` (both **not** in the whitelist; a user-authorized temporary evaluation) |
| Carrier | Both ⑤ legs and the qualified-name probe uniformly use `deep-reasoner` (read-only); the only experimental variable is the `model` parameter |
| Probes | **0** (did not run `.vscode/scripts/agent-probe`) |
| Tier changes | None (no pause-and-change-tier throughout) |

---

## 9. Gate results

| Time point | Command | Result |
|---|---|---|
| ④ (before the report landed) | `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 328/328` |
| ④ (before the report landed) | `node tools/gates/gate.js --all --scope=tree` | `changed=1`, `TOTAL_FAIL=0` (including `PASS base go test ./... :: ok=14 FAIL=0`) |
| After the report draft landed | `node tools/gates/gate.js --all --scope=tree` | `changed=2`; **G4b / G6-4 red → fixed in place**; currently **G7-b red** |
| After removing the `artifacts` block | same as above | **G7-a red** (the changed `docs/validation/*.md` lacks the `artifacts` fence block) |
| ⑨ native verification (after ⑧a landed the artifacts + `git add`) | `--all --scope=tree` / `--all --scope=index` / `--selftest` / `--check=G8-b` | **tree `TOTAL_FAIL=0`**, **index `TOTAL_FAIL=0`**, `SELFTEST: PASS 328/328`, G8-b `1 pack(s) self-consistent` ⇒ **G7 is all green** |
| Base gates | gofmt / build / vet / test (§6.1, carried by `--all`'s base items) | all green |
| Encoding | `.md` = BOM + LF; `.json` / `.js` / `.txt` = no BOM + LF | verified file by file, passed |

**Gate verdict (⑨ measured)**: all three readings are green — tree `TOTAL_FAIL=0`, index `TOTAL_FAIL=0`, `SELFTEST: PASS 328/328`;
the four earlier reds **G4b / G6-4 / G7-b / G7-a** are **all closed** (the repair path is itemized in §7 row 10 and §16).

> `G1` (manual cross-check of three hand-written contents) and `G5` (needs structured input outside the repository) are truthfully printed as `SKIP` by the script and **must not be read as PASS**.

---

## 10. Scope self-check

- `git status --short` → 87 files (4 modified: the CN/EN ledger rows plus the CN/EN directive's v1.25 ledger back-fill; 83 added: the report pair / the dataset / the evidence pack / the freedom list), **all staged file by file with `git add`** (no `-A`); no undeclared untracked file.
- `git grep -n "GLM" -- docs/DELIVERY_DIRECTIVE.md docs/CONTRACTS.md schemas .github/agents` → **no output**.
- `git grep -l '^model:' -- '.github/agents/prfrail-*'` → **no output** (G5-a).
- `git status --porcelain -- .github/agents internal schemas go.mod .github/workflows tools` → **no output**.
- The bodies of §2.2 / §2.4 / §5.2 / §6.3 / §7.5 and §1.5.1 are **zero-changed**; `cjk-newwords.txt` **not extended**.
- Probes **0** (did not run `.vscode/scripts/agent-probe`).
- **Pre-emptive disclaimer**: `deep-reasoner.agent.md` / `independent-reviewer.agent.md` / `quick-verifier.agent.md` under `.github/agents/`
  **contain** `^model:` lines — these **pre-existed**, this slice did not touch them, and they are outside G5-a's scope (G5-a constrains only `prfrail-*`).

**③ test report's independent coverage gaps (copied verbatim)**: G2 (a resident assertion that the bundle has no BOM/LF), G5 (no executable guard for the metering procedure),
G6 (citation parsing swallowing code blocks / URLs), G7 (targeted mutation of non-anchor accepted lines) remain missing; G3 (dataset columns ↔ ① plan) is **untestable** because the plan is not on disk;
G4 (concurrency / crash) **not applicable** (a single-shot CLI).

---

## 11. Explicitly not executed

1. **GLM was not written into any whitelist or role set** (this slice's boundary is "evaluation only");
2. **No multi-round repetition** (n=1, no pursuit of statistical power);
3. **No token metering on the V4 Pro side** (readings unavailable); **no relative cost conclusion given**;
4. **The batch legs' concurrency deviation was not re-tested** (no second call);
5. **No adjustment was made to any scoring caliber beyond the scorer's four generic rules**;
6. **Not committed / not pushed** (authorization boundary: commit no / push no / gitee no);
7. **Probes not run** (quota 0).

---

## 12. Next-step recommendations

1. **A separate slice for the whitelist decision**: this slice only provides evidence; whether to bring `GLM-5.3 (glm)` into the ⑤ position is best discussed after the **cost readings can be closed**
   (the V4 Pro side's tokens become obtainable), avoiding a decision on one-sided evidence alone.
2. **Supplementary cost metering**: take "the V4 Pro side's token reading" as input to a later slice; if it stays unavailable long-term, consider switching to **output character count / elapsed time**
   as a common proxy metric, with its limitations declared in advance.
3. **Applicability of the batch auxiliary tasks**: BATCH-b (long-document excerpting + number checking) and BATCH-d (archive-inventory drafting) produced fully usable deliverables in this sample
   and can be candidates for "auxiliary-type tasks"; BATCH-a's sorting subtask partially failed in this sample, so it is **not suitable** to carry sorting / Top-N selection alone.
4. **Registering the scorer's generic rules**: the four generic correction rules (range citation / comma citation / abstention detection / de-duplication) are recommended to **stay a driver-side tool** in later evaluations
   and not be written into the directive body.

---

## 13. Execution pipeline

| Step | Role / model | Status | Artifact / note |
|---|---|---|---|
| 0 Pre-check | Driver | ✅ | Role / quota / boundary confirmation; the three hard pre-conditions pass (§3.6) |
| 0a Qualified-name probe | Driver; `GLM-5.3 (glm)` | ✅ | Attempt 1 rejected (**no call triggered, no consumption**, but registered as one attempt); attempt 2 accepted ⇒ authoritative qualified name |
| ① Plan | Architect `DeepSeek V4 Pro (deepseek)` | ✅ | `tmp/glm-eval/plan-v1.md` + `PLAN-PROVENANCE.txt` |
| ② Documentation / scaffolding implementation | Implementer `DeepSeek V4.1 Flash (deepseek)` | ✅ | The four bundle files, the dataset scaffolding, the scoring / hash / contamination scripts; **no paid call initiated** |
| ③ Independent adversarial self-test | Test engineer `DeepSeek V4.1 Flash (deepseek)` | ✅ | Two rounds + one round after remediation; `SELFTEST TOTAL=179 PASS=179 FAIL=0` |
| ④ Integration + gates (freeze + pre-registration) | Driver | ✅ | `TOTAL_FAIL=0`, `SELFTEST 328/328`; bundle frozen + key pre-registered |
| ⑤-position two legs (experimental → control) | Driver; carrier `deep-reasoner` (read-only) | ✅ | Same input, same prompt; experimental leg first, control leg after; the two legs' outputs landed separately, not mutually anchored |
| 4 batch legs | Driver; carrier `deep-reasoner` (read-only) | ✅ | **Concurrent execution (deviation, recorded)** |
| ⑤ Pre-review | Pre-reviewer `DeepSeek V4 Pro (deepseek)` | ✅ 2 rounds | Round 1 `PRE-REVIEW: PASS WITH FIXES` (3 Medium) ⇒ remediation ⇒ round 2 `PRE-REVIEW: PASS WITH FIXES` (only 3 Low) |
| ⑥ Independent scan | `MAI-Code-1.1-Flash (copilot)` | ✅ 1 time | `INDEPENDENT SCAN: PASS` (0 file:line-level findings; **confidence downgraded**, see §14) |
| ⑦ Blind / final / re-review | `GPT-5.3-Codex (copilot)` (Extra High, **not downgraded**) | ✅ 3 times | `BLIND REVIEW: FINDINGS` (3 Medium + 1 Low) ⇒ remediation ⇒ `FINAL REVIEW: FINDINGS` (1 High + 2 Medium) ⇒ remediation ⇒ `RE-REVIEW: PASS` (7/7 closed, 0 new) |
| ⑧a Report pair + ledger + evidence pack + freedom list | ⑧ Documentarian + Driver | ✅ | This section |
| ⑨ Native verification | Driver | ✅ | Readings merged into §9 |
| ⑧b Finalization | ⑧ Documentarian | ✅ | Absorbs the ⑨ readings and the ⑧c record |
| ⑧c Close-out cleanup | Driver | ✅ | See §16 |
| ⑩ Stopping point | Driver | ✅ | Stopped at the stopping point, awaiting commit authorization |

**Rework count (§9.2 caliber: the same ring with the same root cause counts as 1)**: ③ ×1 (self-test assertion sync), ⑤ ×1 (re-review after remediation), ⑦ ×2 (one remediation round each after blind review and final review); **fallback interventions ≤ 3**, did not trigger "> 3 times ⇒ mark high risk".

## 14. Review conclusion summary (⑤⑥⑦)

- **⑤ pre-review**: round 1 `PRE-REVIEW: PASS WITH FIXES` (3 Medium: ③ report version caliber / dataset schema self-consistency / the `÷2` mean underestimation; several Low) ⇒ remediation ⇒ round 2 `PRE-REVIEW: PASS WITH FIXES` (**no Medium+**; 3 new Low added, of which NL-1 was closed by ③ updating the assertion and re-running to `179/179`, and NL-2 / NL-3 are stale wording in the handover artifact, handled by "take `FROZEN-MANIFEST` as the sole authority") ⇒ **⑤ closed**.
- **⑥ independent scan**: `INDEPENDENT SCAN: PASS` (Sections A–E complete, 0 file:line-level findings). **A negative signal truthfully recorded**: its Section C item 3 says "`leg-glm.txt` writes S4 as `Medium — 02-diff.js:93`" — which does not match the artifact (both legs report **High** for that seed: GLM leg `F3 | High | 02-diff.js:84-95`, control leg `F5 | High | 02-diff.js:93`) ⇒ per §7.7 that audit's **confidence is downgraded**; its PASS is still adopted (no file:line-level findings), and the final conclusion follows ⑦.
- **⑦ independent final review**: all three readings consumed (§7.5's 1 final + 1 re-review + blind +1). **Neither blind nor final review attached any ⑤/⑥ list**; only the re-review round attached a list (§7.3). After remediation `RE-REVIEW: PASS` (7/7 closed, 0 new findings).

## 15. Review input pack summary (including blind-review isolation proof)

| Item | Value |
|---|---|
| ①–⑤ call isolation (R2.6) | Both ⑤ rounds are **independent invocations** of the `deep-reasoner` carrier; no ① conversation history or implementation-period reasoning was passed in; input = plan + evaluation-evidence **paths** + verbatim contract text |
| ⑥ first round | Attaches **no** ⑤ / ⑦ list; object = evaluation execution and evidence |
| ⑦ blind review (reading 1) | Attaches **no** ⑤ / ⑥ list and no previous-round conclusion |
| ⑦ final review (reading 2) | Attaches **no** list (an independent re-read of the post-remediation on-disk state) |
| ⑦ re-review (reading 3) | Per §7.3, attaches the **full list** and closes it item by item |
| Report body | Per the user's same-round positioning, ⑤⑥⑦ **all do not review the report body** (review object = evaluation execution and evidence) |
| Input sanitization | The prompt gives only **repository-relative paths**; the username segment of absolute paths is replaced with `<user>` in review reports; no IP / hostname / credentials / token |
| Input pack landing | Evaluation evidence is delivered uniformly as **paths** (`tmp/glm-eval/**`, `tmp/glm-eval-out/**`, `docs/validation/directive-glm-eval-dataset.json`), not paraphrased |

## 16. ⑧c closeout-cleaning record (executed by the orchestrator)

| Item | Reading / verdict |
|---|---|
| Repository-wide diagnostics (by the §3.4.1 classes) | Compiled (Go): `gofmt -l .` empty, `go build ./...` / `go vet ./...` exit 0, `go test ./...` `ok=14 FAIL=0`; markup data (`.md` / `.json`): the encoding gate verifies file by file and JSON parsing passes; other plain text: the encoding gate passes; scripting languages: this slice's change set has **no** `.js` / `.ps1` change |
| IDE diagnostics roll-up | All 6 key files in the change set (the report pair / the dataset / the pack `MANIFEST` / the freedom list / the CN ledger) report **No errors found** |
| Gates all green | `--all --scope=tree` `TOTAL_FAIL=0`; `--all --scope=index` `TOTAL_FAIL=0`; `--selftest` `SELFTEST: PASS 328/328`; `--check=G8-b` `1 pack(s) self-consistent` |
| Staging and change set | File-by-file `git add` (**no `-A`**) ⇒ 87 files: 4 modified (the CN/EN ledger rows plus the CN/EN directive's v1.25 ledger back-fill) + 83 added (the report pair / the dataset / the 77 pack payload entries / `MANIFEST.md` / `SHA256SUMS.txt` / the freedom list) |
| Encoding (file by file, all types) | Inside the change set, `.md` = BOM + LF; `.json` / `.txt` / `.js` = no BOM + LF; archived payloads inside the evidence tree follow §1.5.1's "normalize by archived type", their verifiability carried by the pack `MANIFEST`'s "suffix and encoding mapping" dual-hash records |
| Temp directory / probe residue | All **one-off intermediates in `tmp/` are removed** (11 items: `_g4bchk.js` / `_g4b.txt` / `_g8b.txt` / `_gate.txt` / `_gate-tree.txt` / `_persist2.log` / 5 `_append-*.txt`); **class ①/② process artifacts (`tmp/glm-eval/**`, `tmp/glm-eval-out/**`, `tmp/glm-eval-batch/**`, the `tmp/`-root load-bearing scripts) stay in place** — per §1.5.1, "class ①/② files are exempt from §1.5's clear-on-use until the archiving-execution slice handles them; any slice's ⑧c clears class ③ only", and they are **all archived** into this slice's evidence pack; probe residue 0 |
| Untracked-file inventory | Checked item by item: all are declared in the `artifacts` block as `repo` or `local-only`; **no** undeclared residue |
| Difference from the handover artifact (truthfully recorded) | `BRIEFING.md` says "⑧c cleanup (`tmp/` keeps only `.gitkeep`)"; this round, per §0.1 "discipline conflicts are resolved in favour of this directive", follows §1.5.1 and **clears class ③ only**, keeping class ①/② for the "`tmp/` lifecycle execution" slice — the difference is recorded, not silently skipped |

---

> **Artifact-list block (`artifacts` fence) — restored after ⑧a lands**: the `_EN` mirror, the evidence pack (`MANIFEST.md` / `SHA256SUMS.txt`),
> the freedom list and the ledger row **have all landed**, and the cause of the earlier "misstated declaration" is removed ⇒ this section restores the complete block (see the end of the document);
> after restoration, re-run `--scope=tree` / `--scope=index` until G7 is all green.

```artifacts
docs/validation/directive-glm-eval.md	repo
docs/validation/directive-glm-eval_EN.md	repo
docs/validation/directive-glm-eval-dataset.json	repo
docs/validation/evidence/glm-eval/MANIFEST.md	repo
docs/validation/evidence/glm-eval/SHA256SUMS.txt	repo
docs/validation/evidence/DIRECTIVE-GLM-EVAL-freedom-list.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tmp/glm-eval/BRIEFING.md	local-only
tmp/glm-eval/plan-v1.md	local-only
tmp/glm-eval/PLAN-ERRATA.txt	local-only
tmp/glm-eval/key.json	local-only
tmp/glm-eval-out/leg-glm.txt	local-only
tmp/glm-eval-out/leg-control.txt	local-only
```

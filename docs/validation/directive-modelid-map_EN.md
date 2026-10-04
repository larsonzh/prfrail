# Slice validation report · `DIRECTIVE-MODELID-MAP` (`modelId` semantics reversal + completion-line rebuild + `tmp/` lifecycle mechanism)

> Per §5.1 serial pipeline and §6.2 evidence-to-disk obligation of `docs/DELIVERY_DIRECTIVE.md` v1.22. **The CN version is authoritative**, `_EN` is the positional mirror; for the fence block see "Artifacts list" at the end. Size M+, tier `[SLICE]`, probes 0, ⑦ = `gpt-5.3-codex` (**not downgraded**).

## 1. Change summary (including the execution freedom list)

**Goal**: three same-family rule corrections landed in one batch — ① **R2.2 correction** (invocation string = qualified name `Model Name (vendor)`, `modelId` demoted to internal identifier); ② **OB-72** (ledger completion-line structural rebuild); ③ **OB-49** (`tmp/` lifecycle mechanism). Change surface = the two `DELIVERY_DIRECTIVE{,_EN}` files (§0.3 / §1.5 / §2.2 / §2.3 / §2.4 / §5.0 / §5.2 / §6.2 / §6.3 / §11.2 / §12.4 / Appendix A.1 / Appendix B / Appendix C.0l) + this report pair + freedom list.

**Execution freedoms exercised by this slice** (§1.7.3; consistent with this slice's freedom list `docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md`):

- **F-1 version package**: per §0.2 bump to **v1.22** (header version number + one appended §0.3 entry + Appendix C.0l ledger block). A ledger-metadata-class by-product.
- **F-2 §2.2 table form**: **swap the two columns' semantics in place** (column count and order unchanged; of the nine data rows only the two columns exchange values, the values themselves unchanged).
- **F-3 maintenance-mechanism landing point**: the qualified-name maintenance mechanism is written into **the R2.2 body** (rather than a §2.2 note or a new subsection).
- **F-4 OB-72 takes option A′**: = option A's structural skeleton + option B's completion clause (the completion line records only the report path; hash and CI number go only into the report, to be completed by this slice's driver via an independent docs-only commit after CI is doubly green, and that commit still requires same-round authorization per §10). The user offered a choice between options A / B; ① judged that option A still has self-reference (report and completion line in the same commit) ⇒ take a third option, an **escalation boundary** declared as such.
- **F-5 OB-49 landing point**: **add `### 1.5.1`** (so that G6-5 has a parseable section-number target), rather than expanding §1.5's three items or inserting into §11.2.
- **F-6 reference handling**: the three in-directive `tmp/` references get **annotations only, no path change**; the actual archiving and path rewriting are deferred to the archiving execution slice.
- **F-7 boundary rejection**: **reject** ⑧'s "out-of-scope side fix" of the 5 **pre-existing baseline** CN / `_EN` backtick-level drifts (each reverted), judged under §1.7.2 as out-of-scope side fixes, and instead registered as candidate OBs for user adjudication.
- **F-8 historical-line handling**: the v1.12 entry in §0.3 **keeps its original text**; only a "semantics-correction pointer" is appended (per §11.1).
- **F-9 orchestration artifacts**: the structure and wording of the report and the list.

## 2. ① judgement and landing points (solution document `DESIGN: DONE`)

| # | Judgement item | Conclusion |
|---|---|---|
| 1 | **whether the three-item merge is controllable** | **controllable, not split** (the change surface is concentrated in the same authoritative file pair; the three items belong to the same rule semantics and fall in the same review blind spot; ⑦'s single blind review + re-review can cover them as a whole; splitting into three slices would raise ⑦'s quota from 1+1 to 3+3, and the slices would lock each other's order) |
| 2 | R2.2 table form | swap the two columns' **semantics** (the qualified-name column = sole authority for the invocation string; the `modelId` column = internal identifier), column count and order unchanged |
| 3 | R2.2 body core | **keep** "must explicitly pass `model`", the outer form changes to the qualified-name format, and a **qualified-name maintenance mechanism** is added |
| 4 | candidate criterion G8 | **not implemented** (§6.3 criteria untouched); only its semantics synced to "the `model=` set of Appendix B ⊆ the §2.2 qualified-name column" |
| 5 | OB-72 option | option A′ (see F-4) |
| 6 | OB-49 landing point | new §1.5.1; the three references annotated only; this slice performs no archiving |

**Evaluation-item conclusion (the user required "if the conclusion is that a maintenance mechanism is still needed, it must be written explicitly in the report")**: **the qualified name is not more stable than `modelId`** — it is merely **the only form accepted by `runSubagent`** (three `modelId`-form attempts all rejected in testing + the official documentation explicitly uses `Model Name (vendor)` for `handouts.model`); `modelId` falls back to being the most stable internal identifier (whitelist closed set / G5-b / ⑦ tier table / `LEGACY_FLASH_MODEL_IDS` all unaffected by version drift). **A qualified name carries a version number (e.g. `DeepSeek V4.1 Flash`) ⇒ version drift exists and a maintenance mechanism is still needed**, placed in the R2.2 body: when an invocation is rejected by the tool, **stop and escalate to the user for adjudication per §1.7.2**; do not swap strings and guess. Updates to the §2.2 qualified-name column and Appendix B template strings require user authorization and proceed alongside rule revisions. **Existing blind spot (truthful)**: if the tool silently accepts an old string but routes to a different model, the mechanism does not trigger — an inherent blind spot of mapping-style rules; this slice does not solve it and only registers it.

## 3. Delivery surface and per-item landing points

| Item | Landing point |
|---|---|
| R2.2 body and maintenance mechanism | the R2.2 row + the qualified-name note under the §2.2 table of `DELIVERY_DIRECTIVE{,_EN}` |
| §2.2 two-column swap | same table header + nine data rows (values unchanged) |
| reference-point sync | §2.3 final sentence / §2.4 two places (blacklist body untouched) / §5.0 steps 2 and 4 / §6.3's G5-b parenthetical / §5.2's ⑦ tier-table column name |
| Appendix B template | six role lines + seven `model=` invocation strings + B.0's two list lines |
| OB-72 | §11.2's "ledger completion-line format" paragraph rebuilt; OB-72 flips to `Resolved` |
| OB-49 | new §1.5.1; the first line of Appendix A.1 and the v1.10 / v1.12 entries of §0.3 get archiving annotations; OB-49 flips to `To be discussed` |
| ledger flips | OB-67 `To be discussed → Resolved`; OB-70 stays `Resolved` and gains a landing sentence; OB-29 gains a candidate-G8 semantics-sync sentence |
| version package | header version v1.22 + §0.3 entry + Appendix C.0l ledger block |
| freedom list | `docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md` |

## 4. Execution pipeline and environment

① architecture (V4 Pro, split judgement + the three-item design) ⇒ ② document implementation (CN authoritative draft, per-group assert hit counts) ⇒ ③ document consistency check (`modelId` residue classification: no residue apart from historical lines and the "internal identifier" usage) ⇒ ④ gates ⇒ ⑤ pre-review (V4 Pro, **two rounds**: round 1 `PASS WITH FIXES` ⇒ remediation ⇒ re-review `PASS`) ⇒ ⑥ independent scan (MAI, **two rounds**: round 1 `FINDINGS` ⇒ remediation ⇒ re-review `PASS`) ⇒ ⑦ final review (Codex, **two rounds**: round 1 `PASS WITH FIXES` ⇒ remediation ⇒ re-review `PASS`) ⇒ ⑧a/⑧b report and this `_EN` mirror ⇒ ⑨ native verification ⇒ ⑧c close-out cleanup ⇒ ⑩ stopping point.

**Quota actuals**: probes 0; ① one round; ⑤ two rounds (quota 1+1); ⑥ two scans (cap 3); ⑦ two rounds (quota 1+1).
Environment: Windows local machine; Node v24.17.0; Go toolchain `go1.27.0 windows/amd64`; repo `github.com/larsonzh/prfrail`; branch `main`, starting point `4a72c8f`.

### 4.1 Exception and fallback records (a report element of §11.3)

- **Role overreach and revert (F-7)**: besides this slice's mirror, the ⑧ mirror round also made pure-format alignments of **5 pre-existing baseline** backtick-level CN / `_EN` drifts (one of which changed English wording) ⇒ each **reverted**, so that the change set contains only this slice's mirror; the differences are registered as candidate OBs (see §6).
- **`create_file` encoding anomaly when writing to disk (tool side)**: the freedom list was written to disk as no-BOM + CRLF ⇒ ⑦ round 1 flagged red (G4a) ⇒ normalized to BOM + LF and green again; from then on all new text is normalized first.
- **Terminal command-chain interruption (tool side)**: a long command chain was misjudged by the terminal as interrupted, and the first half of the chain may have already executed ⇒ the affected steps were re-run as "one call per step", and after the interruption the workspace and index were first checked for consistency.
- **① judgement deviating from the user's expectation (truthful)**: the flip item the user gave included "OB-70 `To be discussed → Resolved`", but OB-70 was already `Resolved` in v1.19 ⇒ the actual action changed to **OB-67 flips to `Resolved`**, OB-70 gains a landing sentence, and the C.0l note wording is synchronously corrected (raised as a ⑤ round-1 Low, fixed).
- **No other fallback**: the remaining §9.1 scenarios were not triggered (no budget overrun, no dependency blocker, no repeated native failure).

## 5. Review conclusion

- ⑤ pre-review (V4 Pro, two rounds): round 1 `PASS WITH FIXES` (2 Medium: the §11.2 completion commit is not linked to §10 authorization; §1.5.1 lacks the explicit §1.5 exemption; 6 Low: a CI-less slice does not backfill hash and run number; the OB-67 disposition-column wording; the C.0l note's OB-70 inconsistency; archiving must pass the ⑧c encoding gate first; the cleanup timing must be preceded by a full inventory; the §6.2 path rule points at §1.5.1) ⇒ item-by-item remediation ⇒ re-review **`PASS`** (all closed, no new issues).
- ⑥ independent scan (MAI, two rounds): round 1 `INDEPENDENT SCAN: FINDINGS` (1 High: the v1.12 historical entry of §0.3 still says "Appendix B's seven invocation strings changed to modelId"; 1 Medium: 5 pre-existing baseline backtick-level mirror drifts) ⇒ the High is closed by **appending a semantics-correction pointer** (the historical text is not changed), the Medium is handled by **boundary rejection** and registered as a candidate OB ⇒ re-review **`INDEPENDENT SCAN: PASS`** (both dispositions explicitly accepted, no new High/Medium).
- ⑦ final review (Codex, two rounds): round 1 `FINAL REVIEW: PASS WITH FIXES` (1 Medium: the freedom list lacks BOM + CRLF ⇒ G4a flags red; the other four items and F-1 to F-9 judged item by item as not out of scope) ⇒ encoding normalized ⇒ re-review **`RE-REVIEW: PASS`** (no findings, both scopes reproduce `TOTAL_FAIL=0`). **⑦ not downgraded ⇒ this slice has no ADR** (① once proposed registering ADR-023; per the user's report format "if not downgraded, no ADR", it is not registered).

### 5.1 Review-input package summary (with blind-review isolation proof)

Change-set diff = `tmp/dml/change-set.diff` (with `tmp/dml/numstat.txt`); contract sentences = §0.2 / §1.5 / §1.7.2 / §1.7.3 / §2.2 / §2.3 / §2.4 / §5.0 / §5.1 / §5.2 / §6.2 / §6.3 / §7.2 / §7.5 / §10 / §11.1 / §11.2 / §12.4 of `DELIVERY_DIRECTIVE.md`; freedom list = the fixed path in the last line of §3. **Blind-review isolation proof**: this slice touches "identity" semantics (the model whitelist / blacklist is identity admission) ⇒ ⑦ **round 1 is executed blind** — the input contains only the slice definition, the raw diff (paths), the contract sentences, the freedom list and the read-only constraint, and **does not attach** the ⑤ pre-review report or the ⑥ scan report; the re-review round follows Appendix B.6 by allowing the summary that is "read after the independent scan completes". Input sanitisation = every item is repository text, with no credential / token / SSH path.

## 6. Falsifiability and gate results

| Check | Result |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` (changed=7; the final reading after the report pair and the ledger row are staged) |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0` (changed=7) |
| `gate.js --selftest` | `SELFTEST: PASS 264/264` (this slice adds no fixture) |
| base gates (§6.1) | `gofmt` empty; `go build` and `go vet` and `go test` all exit 0 (ok=14, FAIL=0) |
| mirror positionality (including the five shapes) | reading in `tmp/dml/pairs.txt` (equal element counts on both sides; shapes differ only in 5 pre-existing baseline backtick drifts, not on lines changed by this slice) |
| mutation check | reading in `tmp/dml/mutation.md`: ① change OB-49's first status word ⇒ G6-6 flags red; ② change a section-number reference inside §0.3 ⇒ G6-5 flags red; ③ change the CN only by adding a line ⇒ G3(a.1) and G3(a.2) flag red; all three verified by sha256 and then restored byte for byte and green again |
| added-line precheck | `tmp/dml/pre.txt`: `ADDED-LINE HITS=0`, `UNKNOWN-CJK=` (no unregistered new characters) |

**Falsifiability boundary (truthful)**: of this slice's three items, ①'s **invocation-string semantics** and ③'s **three-class disposition semantics** are norms with **no new mechanical criterion** (the user constrained "§6.3 criteria untouched") ⇒ their verification = ⑤ / ⑥ / ⑦ human re-review + G6-5 / G6-6 corroboration; **the mutation check does not cover them**, and this must not be read as "guarded by a criterion". ②'s completion-line format is the same (no criterion).

**Registered OBs (user-approved on 2026-10-04, first status word `To be discussed`; registered by this slice in §12.4)**:
- **OB-73**: the CN / `_EN` backtick-level mirror drift (**5 pre-existing baseline** cases, 2 of which contain EN-only sentences) — G3(a.2)'s five shapes do not include a backtick count, so the gate does not cover it; fixing it requires touching historical lines unrelated to this slice ⇒ rejected at the boundary by this slice (the ⑥ Medium disposition).
- **OB-74**: Haiku, as the only exemption outside §2.2, has its **qualified name undefined** ⇒ no invocation string can be written when it is enabled (a pre-existing gap; this slice neither introduces nor claims to solve it).
- **OB-75 (a direct product of this slice's assessment item; to be stated explicitly)**: the R2.2 maintenance mechanism **triggers only when a call is rejected**; if the tool **accepts an old string while routing the request to another model**, "call string = qualified name" has **no observable criterion** against silent routing - an **identified blind spot** (not a defect; it is inherent to mapping-style rules). This slice only records it in R2.2 and in the ledger and **does not solve it**, leaving it to a model-observability mechanism (recording and cross-checking the model actually in effect on the operator side).

## 7. Cost and accounting

① one round, ⑤ two rounds, ⑥ two scans, ⑦ two rounds (Codex regular tier, not downgraded), ⑧ two rounds (the CN authoritative draft landed by the driver, the `_EN` mirror produced by the documenter); probes 0; no Haiku calls; **no quota exception occurred, so this slice has no ADR**.

## 8. Explicitly not executed

- The candidate criterion **G8** is not implemented; no §6.3 criterion body is changed; `g5b.js` and the fixtures are not touched.
- The §2.4 blacklist body is not changed; the §5.2 matrix is not changed; the **values** of the §2.2 table are not changed (including the divergence between the hyphen form of `GPT-6-Luna (copilot)` and the space form measured from the tool — registered only).
- **No `tmp/` file was actually moved, archived or deleted** (archiving execution is a separate "tmp/ lifecycle execution" slice).
- The 5 pre-existing baseline backtick-level mirror drifts are not fixed (OB-73, registered).
- No ADR is registered (⑦ not downgraded); after the ⑩ stopping point the commit and the push were completed under the user's same-round authorization - see the "Commit and CI" row of §10.

## 9. Next-step recommendations (including tier recommendations per role)

- **Next slice (recommended)**: the project name `tmp/` lifecycle execution slice — per §1.5.1's three classes, execute a full inventory of the stock ↔ references, archiving and reference rewriting; size **M+**, ⑦ uses `gpt-5.3-codex` (touches rule-body references).
- **⑦ tier recommendation**: this slice and the archiving execution slice are both rule-semantics class ⇒ the regular tier `gpt-5.3-codex` suffices; it touches no security / machine state / protocol semantics ⇒ no deep tier needed.
- **⑧ role recommendation**: the mirror round must first read the CN authoritative draft's per-group changes; aligning to baseline drifts unrelated to this slice is **forbidden** (a measured lesson of this slice).

## 10. ⑩ stopping-point report (including the write-back check line)

- Executed up to ⑨; **stopped before commit authorization**: at the ⑩ stopping point this slice did not commit, did not push, did not push to gitee; that stopping point was lifted by the user's same-round authorization - see the "Commit and CI" row below.
- **Write-back check** (the OB-63 formulation of §11.2): validation report = **written back** (this file); `DEV_PLAN` = **not applicable** (a directive-governance slice; per the slice definition and existing precedent, no `DEV_PLAN` paragraph is written); ledger = **written back** (the completion-record row of `REMAINING_SLICES{,_EN}`, which per the new §11.2 format records the report path only); §0.3 and §12.4 of `DELIVERY_DIRECTIVE{,_EN}` are updated in place; ADR = **not applicable** (no downgrade, no quota exception).
- **⑦ conclusion**: round 1 `PASS WITH FIXES` (1 Medium: the freedom-list encoding) ⇒ encoding normalized ⇒ re-review `RE-REVIEW: PASS`; ⑦ is Codex regular tier, **not downgraded**, and its independence is not compromised.
- **Gate readings**: both `--scope=tree` and `--scope=index` report `TOTAL_FAIL=0`; `SELFTEST: PASS 264/264` (see §6).
- **Commit and CI**: commit `0e7a111`; CI run `37185556184` is green on both `Go` legs (`Go windows-latest` / `Go ubuntu-latest`), while the three `workflow_dispatch` legs `candidate-build` / `candidate-probe` / `selfhost` are skipped by their event condition; both legs agree: `SELFTEST: PASS 264/264`, `GATE REPORT scope=ci changed=7`, `TOTAL_FAIL=0`.

## 11. Artifacts list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-modelid-map.md	repo
docs/validation/directive-modelid-map_EN.md	repo
docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md	repo
docs/validation/evidence/DIRECTIVE-MODELID-MAP-mutation.md	repo
docs/validation/evidence/DIRECTIVE-MODELID-MAP-review-r1.md	repo
```

(Closeout: the one-off scripts and intermediate readings under `tmp/dml/` were removed per §1.5, so the granular `local-only` rows previously listed above are gone with them; the two pieces of evidence are counted as `repo` rows - `review-r1.md` as an archived copy and `mutation.md` as a **rebuilt archive** (the original mojibake is irreversible; its rebuild basis is recorded in its own note).)
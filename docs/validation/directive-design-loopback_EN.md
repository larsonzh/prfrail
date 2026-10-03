# Slice validation report · `DIRECTIVE-DESIGN-LOOPBACK` (design flow-back mechanism assessment)

> Per §5.1 serial pipeline and §6.2 evidence-to-disk obligation of `docs/DELIVERY_DIRECTIVE.md` v1.19 (the starting version of this slice; this slice lands a v1.20 ledger by-product). **The CN version is authoritative**, `_EN` is the positional mirror; for the fence block see "Artifacts list" at the end. Size S, tier `[SLICE]`, probes 0.

## 1. Change summary (including the execution freedom list)

**Goal**: ① assess OB-71 (session conclusions have no flow-back); ② assess OB-76 (missing RFC ↔ directive mapping); ③ assess the "task-package closing report" mechanism; ④ produce light / medium / heavy three candidates with a criterion comparison; ⑤ hand to the user for a ruling. **No rule body is landed** (§2.2 / §2.4 / §5.2 / §6.3 / §11.2 etc. are all unchanged).

**Execution freedoms exercised by this slice (§1.7.1, substantive choices)**:

- **F-1 version package**: the only change this slice lands into the directive is one OB line in §12.4 and one numbering-source note; per §0.2 "every change must bump the version number and append a §0.3 log entry", it is synchronously raised to **v1.20** with a §0.3 line added (a ledger-metadata by-product, containing no rule semantics).
- **F-2 numbering choice**: this line takes **OB-71** in §12.4 ledger order (the next available number; based on OB-66's "§12.4 is the sole source of truth"), and records the two external quoted numbers' failure to land via the numbering-source note.
- **F-3 numbering-provenance note (escalation boundary; landed before stopping - recorded as such)**: one blockquote was added after the §12.4 table (taking no OB number) recording that the "OB-71" and "OB-76" cited by this slice definition are unregistered - that note **exceeds** the literal "land only a candidate OB" scope and is therefore a §1.7.2 escalation boundary; this slice **did not stop first** before landing it, so it is recorded as such and handed to the user for a same-round ruling (the two options of D-4 in §7: keep and ratify, or delete and re-run the gates).
- **F-4 assessor**: this slice is a pure assessment slice, with no rule body to change ⇒ ① (architecture) receives no dispatch, and the assessment is completed by the orchestrator on the basis of read-only measurement (the same precedent as `DIRECTIVE-CLEANUP-EVAL`).
- **F-5 report production**: this report is assembled by the orchestrator (§3.2 orchestration-artifact exception: aggregating sub-agent output and gate results); that reading **exceeds** the literal "assembly" wording of §3.2 (the assessment prose is authored by the orchestrator) and rests on the user's same-round authorization (the slice definition names "the assessment report" as the deliverable), the same reading as `DIRECTIVE-CLEANUP-EVAL`.
- **F-6 new-character registration**: this slice registers 1 character in `tools/gates/cjk-newwords.txt` ("瞻" of "前瞻"; reason = it is the user's original wording for the outlook part of the closing report, and no synonymous replacement preserves the term), which is the §6.3 G4b existing mechanism, not a relaxation of the criterion.
- **F-7 ⑦ downgrade**: this slice's ⑦ uses `deepseek-v4-pro` (explicitly stated in the user's slice definition), and is recorded per the OB-37 meta-rule as **ADR-021**.

## 2. Pre-start checklist verification (step 0; item by item)

| Item | Content | Status |
|---|---|---|
| A-1 | `DIRECTIVE-MODELID-MAP` project initiation (`modelId` → qualified-name mapping table; ⑦ not downgraded) | **To do** (the user has ruled that a separate slice will be created; this slice is a design flow-back assessment, not that slice) |
| B-2 | Backfill the commit hash and CI run number into the DIRECTIVE-CLEANUP-EVAL row of `REMAINING_SLICES{,_EN}` | **To do** (this slice's authorization boundary does not include that cross-slice backfill; listed as ruling item D-5) |
| C-3 | OB-68 (§5.2 tier table ↔ §7.5 / §8.2 four labels and two-tier cost / tier) | **To do** (touches the rule body, beyond this slice's constraint) |
| C-4 | OB-53 item ⑤ (chain-config side `model: "auto"`) | **To do** (touches `internal/**`, forbidden to change in this slice) |
| C-5 | Landing-point assessment of OB-71 and OB-76 | **Done** (this report §3 and §4) |
| C-6 | Registering the "task-package closing report" mechanism as a candidate OB | **Done** (this report §5; and OB-71 is registered in §12.4) |
| D-7 | Orchestrator handover audit (= at the end of T027) | **Not applicable** (this slice is not at that anchor) |
| D-8 | Strategic retrospective (= at the end of S1) | **Not applicable** (this slice is not at that anchor) |
| E-9 | Review of closed items (OB-63 / OB-66 / OB-70 and two ledger statements) | **Done** (item-by-item read-only review, consistent with §12.4 and the ledger) |
| X-1 | New assessment object for this slice: numbering-source measurement (repo-wide search for `OB-71` to `OB-79`) | **Done** (0 hits; see §3.1, and escalated to ruling item D-4) |
| X-2 | New observation for this slice: `DIRECTIVE-CLEANUP-EVAL` did not land the §1.7.3 fixed-path freedom list | **Done** (found by read-only measurement; per the constraint it is **not registered** as an OB, and is listed among the §11 candidates) |

## 3. ① Assessment of OB-71 (session conclusions have no flow-back)

### 3.1 Numbering source: the two assessment objects do not exist in the ledger (meta-evidence)

- Repo-wide verbatim search for `OB-71` through `OB-79`: **0 hits** (including `tmp/**` under `.gitignore`; **search time = before this slice landed its files**, after which those strings enter the corpus, so this conclusion rests on the HEAD corpus); the highest number in §12.4 is OB-70 ⇒ the "OB-71" and "OB-76" cited by this slice definition are **unregistered in the ledger**.
- This fact itself is direct evidence for ①: even the number of "what is to be assessed" exists only in the session; if this slice does not land a line, these two assessment objects would have no index at all after the session ends.
- Disposition: per OB-66's ruling (§12.4 is the sole source of truth), this slice registers the "task-package closing report" as **OB-71**, and adds a numbering-source note in §12.4; the mapping of the two external quoted numbers is listed as ruling item (§7's D-4).

### 3.2 Item-by-item verification of the three examples (machine-checkable)

| Conclusion quoted by the user | This slice's measurement | Verdict |
|---|---|---|
| sessbridge "three generations" | A search of prfrail `docs/**` for `三代` returns 0 hits (same search time as §3.1); a search of sessbridge `docs/**` for `演进` returns 0 hits. The channel's positioning changed three times on the prfrail side (ADR-004 / ADR-011 / ADR-012), but **none of them narrates generations or their drivers**; the sessbridge side keeps only 2 naming records (the naming-record section of `README.md`) | **partially flowed back**: the carrier for the conclusion exists (the ADRs), but the **generation split and its drivers did not flow back** |
| `modelId` rejected | Measurement records of OB-67 and OB-70 and ADR-019 and two validation reports are complete | **written back**, but the trigger condition = **same-round user ruling** (not a mechanism); this formulation has been written since v1.12 while all 157 historical calls used display-name strings ⇒ the rule and practice diverged for a long time unnoticed |
| whois-era "V3.2 did not hold up" | A search of prfrail `docs/**` and whois `docs/**` for `v3.2` / `V3.2` / `撑不住`: **0 valid hits** (all whois-side hits are product version numbers or release notes, unrelated to the directive) | **did not flow back**; the §12.2 rollback path targets `FLASH_OPERATING_DIRECTIVE_v3.1`, and why it is not v3.2 is not traceable within the repository |

### 3.3 Impact of conclusions that did not flow back

- **Rule and practice diverged for a long time without awareness**: the deviation of the `modelId` formulation can only be found by a dedicated measurement, and the measurement was triggered by the user ⇒ without that ruling, the deviation would keep spreading with every new task package.
- **The "why" of the rollback path is not traceable**: if a versioned directive is proposed again, the earlier conclusion that "some approach did not hold up" cannot be cited anywhere ⇒ the same argument must be redone, at several times the cost of an assessment slice like this one.
- **Cross-slice omissions are invisible**: unfinished items exist only in their own slice reports, with no single aggregation ⇒ the difference between session judgments and repo state is checked by no one (the mechanism in §5 targets exactly this).

### 3.4 Relationship with existing mechanisms (which covers, which does not)

| Existing mechanism | Covers | Does not cover |
|---|---|---|
| §12.4 observation ledger | The entry is open; it is the formal container for "undecided / observation" | Only collects undecided items; **settled design conclusions** have no container (§3.1 is of this kind) |
| §11.2 write-back formulation | Four kinds of write-back at slice close-out (report / `DEV_PLAN` / ledger / ADR) | Only covers evidence of **the slice itself**; conclusions formed between slices within a session have no owner |
| ADR registration | **is precisely** the flow-back carrier for design decisions (ADR-004 / ADR-011 / ADR-012 are examples) | **no trigger criterion** (when an ADR must be written); the actual trigger depends on the memory of the user and the orchestrator |
| §12.5 gap table | Registers mechanisms and roles not yet landed | It records "what mechanism / role is missing", not "a conclusion's missing ownership" |
| §7.8 convergence criterion and §12.1 pilot metrics | Stopping decision of the review chain and pilot acceptance | Does not include the dimension of "whether a conclusion has flowed back" |

⇒ **Gap localization**: what the repo lacks is not a **container**, but **a trigger and a default owner** — "which kinds of conclusions qualify for flow-back / when they must land / which kind of carrier they land in".

### 3.5 Conclusion

- It is worth supplementing a mechanism, but **adding a new carrier is not recommended** (ADR plus §12.4 plus reports are already complete; adding another layer necessarily duplicates and introduces drift).
- The recommended minimal form = **one trigger criterion**: any conclusion formed in a session that would change subsequent execution or product direction must, within **the same round**, land in at least one of ADR / OB / report and leave a pointer; it can be landed together with the mechanism of ③, with no need for a separate project.
- This criterion is rule semantics ⇒ this slice **does not land it**, handing it to the user for a ruling (§7's D-1).

## 4. ② Assessment of OB-76 (missing RFC ↔ directive mapping)

### 4.1 Current state (cross-reference matrix, this slice's measurement is reproducible)

| Direction | Measurement | Note |
|---|---|---|
| Directive → RFC | **2 places**, and both are incidental | The known residual false-positive note of G6-3 in §6.3 cites the RFC's `file:line`; OB-17 is from the same source ⇒ no design-layer pointer |
| RFC → Directive | **0 places** | The RFC header only points to the authoritative domain table of `DOCUMENTATION_PLAN` |
| Product-layer authoritative docs → Directive | **0 places** | `PRODUCT_REQUIREMENTS` / `ARCHITECTURE` / `CONTRACTS` / `SECURITY` / `TEST_STRATEGY` / `BUSINESS_WORKFLOWS` are each measured as 0 file by file |
| Product-layer authoritative docs → RFC | **1 place each** | All are header "source" pointers ⇒ the product layer's pointers are **single-layer** (pointing to historical sources, not to the discipline layer) |
| Directive ← others | `DEV_PLAN` 1, `DOCUMENTATION_PLAN` 1, `REMAINING_SLICES` 4, plus historical drafts and various validation reports | The discipline layer is cited only by **its own** ledger and reports |

- The only document simultaneously holding both layers' positioning is §1 of `docs/DOCUMENTATION_PLAN{,_EN}.md`: its opening paragraph positions the RFC as "the project proposal and historical design source, no longer the single all-domain authority", and its authoritative domain table carries the "delivery execution discipline = `DELIVERY_DIRECTIVE`" row (whose third column notes "Not applicable (engineering discipline, not a product design source)") ⇒ **the positioning relationship exists, but the ID-level traceability relationship does not** (the user's "no formal traceability relationship" holds, but "completely unrelated" does not).
- **Existing pointer precedent**: the RFC **already uses** a pointer to the normative layer (§10.3's title is "state machine and recovery design source (specification migrated to contract §2.1)") ⇒ "add one pointer line" is a mature practice in this repo, not a new mechanism.
- **No occurred semantic conflict found by this slice**: a spot check of the RFC's positioning of SessionBridge (`silent` transport and `cmd_<pid>.json` etc. external adapter wire, appearing 7 times in the RFC) against ADR-012 (AgentRunner is the formal execution port, SessionBridge is downgraded to a message bridge) is **compatible**; `AgentRunner` has 0 hits in the RFC (that ADR is later than the RFC) ⇒ this is a **time gap**, not a contradiction.

### 4.2 Impact

- **One-way divergence has no observation surface**: product-layer changes do not reach the discipline layer, and the discipline layer's process changes do not reach the product layer ⇒ the two ends evolve separately, and conflicts can only be found by someone "happening to read both documents at the same time".
- **Broken reader path**: the reading order of `DOCUMENTATION_PLAN` can reach the directive, but conversely, starting from the RFC one **cannot reach** the discipline layer, nor the new positioning given by the ADR (`AgentRunner` does not appear in the RFC).
- **The cost has not yet materialized**: this slice measured only "cross-layer traceability is near zero", not any defect that has already occurred ⇒ this impact should be counted as **risk**, not as loss already incurred.

### 4.3 Landing-point candidates

| Candidate | Form | Cost | Covers | Does not cover |
|---|---|---|---|---|
| (a) Directive Appendix A | Add one line "product definition layer" to A.1 "Sources of this directive" | Minimal (one line plus a pointer inside the directive) | One-way pointer from discipline layer → product layer | The reverse remains 0 |
| (b) `DOCUMENTATION_PLAN` §1 | Add one row "cross-layer pointer" to the authoritative domain table, and write the OB-19 channel in explicit text | Small (that document is already the sole positioning place, and directive semantics are untouched) | **Bidirectional** (both layers start from it) | ID-level traceability |
| (c) Standalone mapping document | Create `docs/DESIGN_TRACEABILITY{,_EN}.md` (the 13th document pair) | Medium (a new document plus a maintenance obligation for every rule change, and a criterion must accompany it) | ID-level bidirectional traceability | None (other than cost and maintenance) |

### 4.4 Conclusion

- **Worth formalizing, in minimal form**: landing **(b)** is recommended. Reasons = ① it is the only document simultaneously holding both layers' positioning (this slice's measurement); ② it does not touch the directive's rule semantics (consistent with this slice's constraint); ③ "adding a pointer" already has precedent in this repo (RFC §10.3); ④ this slice found no occurred divergence defect ⇒ take the conservative form per §1.7.2.
- **(c) is not recommended**: the actual demand for ID-level mapping is not supported by evidence (the only ready-made two-layer channel, OB-19, is item-by-item OBs, not a systematic mapping), and a new document would immediately create the gap of "a new document with no criterion" (the same family as the "criterion self-bootstrapping risk" already registered in §12.5).
- (a) can serve as a supplement to (b) (if the discipline layer is expected to self-certify its product anchor), but is not a separate project.

## 5. ③ Assessment of the "task-package closing report" mechanism

### 5.1 Form and key boundary

- Form (already decided by the user): the close-out document at the end of a task package = **retrospective** (what was done / what was completed / what was not completed and why) + **outlook** (business and technical additions / cancellations / changes / cautions).
- **Key boundary (this slice's structural finding)**: all current output roles in the directive **have no "business direction" authorization** — §1.7.3 explicitly states "how to do it is autonomous, **what to do always belongs to the user**", and §3.2 further forbids the orchestrator from originating product artifacts ⇒ if the "outlook" includes business direction, its **producing subject can only be the user**.
- ⇒ The mechanism must be written as **two artifacts**: the "retrospective + open-questions list" produced by the orchestrator and the "outlook ruling" by the user, rather than one document authored solely by the orchestrator.

### 5.2 Relationship between trigger anchors and existing carriers

| Existing carrier | Granularity | Covers "phase close-out" |
|---|---|---|
| §11.3 slice validation report (including the final item "next-step suggestions") | Single slice | No |
| Appendix C.0* ledger block | Directive-governance slice | No (covers directive governance only) |
| §12.4 observation ledger | Single observation item | No (collects undecided items only) |
| §12.5 gap table | Future mechanisms / roles | No (it is a forward-looking gap register, not a close-out) |
| `DEV_PLAN` paragraphs and `REMAINING_SLICES` completion records | Single task / single slice | No |

- ⇒ **All existing carriers are either "single slice / single task" or "single kind of metadata"; there is no carrier for "cross-slice phase close-out"**.
- The anchors already have names but no landing points: OB-19's "Sol strategic retrospective question set", "orchestrator handover audit = at the end of T027" from the previous slice's pre-start checklist (item 既有-4 in `docs/validation/directive-cleanup-eval.md`), and "at the end of S1" in the checklist.
- **Corroboration in reverse**: the three examples in ① of this slice are not single-slice evidence gaps, but **cross-slice aggregation gaps** ⇒ they share the same origin as "no carrier for phase close-out".

### 5.3 Relationship with the registered gap in §12.5

- §12.5 already has a similar gap row: "retrospective analyst (multi-slice data retrospective) | **after execution data accumulates**" ⇒ half of the mechanism is already registered (who does it); what is missing is the **form** (carrier), the **trigger anchor**, and the **ownership boundary of the outlook**.
- Therefore this mechanism **does not need a new role**: the orchestrator plus the user can carry it; the "retrospective analyst" in §12.5 is kept as a **future** enhancement (when there are enough slices and a horizontal data retrospective is needed).

### 5.4 Conclusion

Worth doing, for three reasons: ① the anchors have been repeatedly cited yet have no carrier; ② "what was not completed and why" is scattered across slice reports with no aggregation ⇒ cross-slice omissions are invisible (① of this slice is of this kind); ③ the cost is controllable (a template plus two anchors plus one trigger criterion) and no new role is added. **Precondition**: the outlook must be split into "the orchestrator's question list + the user's ruling"; the orchestrator must not write the business direction on the user's behalf.

## 6. ④ Comparison of the three candidates and criteria (no recommendation, criteria only)

| Candidate | Form | Cost (basis: this slice's measurement) | Solves | Does not solve | Applicable precondition |
|---|---|---|---|---|---|
| **Light** | Add only a "closing report" template and trigger anchors (following the existing report style) | Minimal: one template plus two anchors written as one reference line in the directive | Phase close-out has a carrier; unfinished items can be aggregated | ②'s cross-layer traceability (still 0 pointers); ①'s conclusion retrieval | Pain point = no one writes the phase close-out |
| **Medium-tier** | Add a "design state snapshot" document (incremental, countable, mechanically checkable) and a closing report | Medium: one new document plus a per-phase increment plus at least one criterion (otherwise it repeats "gap without criterion") | The above, plus the current design state being traceable and comparable | RFC versioning; automatic identification of conclusions | Pain point = design state not traceable (supported by this slice's measurement: cross-layer references near zero) |
| **Heavy** | RFC versioning plus a decision flow-back process (writing "when an ADR / OB must land" as rules) | Heavy: an RFC version-bump mechanism plus full-process obligations plus criteria plus maintenance cost (the RFC is already near 190 KB) | The above, plus bidirectional ID-level traceability | — | Requires ID-level bidirectional traceability and the ability to bear the maintenance; **this slice's evidence does not yet support it** (no occurred divergence defect found) |

**Criteria (listed only, no direction chosen)**: ① if the pain point is "no one writes the close-out" ⇒ light is enough; ② if the pain point is "design state not traceable" ⇒ medium is the lower bound (light and the status quo are both "0 cross-layer pointers", see §4.1); ③ if the pain point is "divergence must be caught mechanically" ⇒ only then is heavy worth it, and **there must first be a criterion that can turn red**, otherwise it merely copies the "criterion self-bootstrapping risk" to a new layer.

## 7. ⑤ Hand to the user for a ruling (option list)

- **D-1**: whether to adopt the "conclusion flow-back trigger criterion" for ① (recommendation: merge into the mechanism chosen for ③, not a separate project).
- **D-2**: the three-way choice of ②'s landing point ((a) Directive Appendix A / (b) `DOCUMENTATION_PLAN` §1 / (c) standalone mapping document) — this slice recommends **(b)** in minimal form.
- **D-3**: the three-way choice for ③ (light / medium / heavy) plus the first anchor (end of T027 / end of S1 / both) plus confirmation of the boundary "outlook = orchestrator question list + user ruling".
- **D-4**: numbering provenance (including the two options for the **F-3 escalation boundary**) - this slice took **OB-71** in ledger order (= the "task-package closing report" mechanism); the numbering-provenance note in §12.4 (an escalation boundary landed before stopping) awaits a same-round ruling: (a) keep and ratify, or (b) delete and re-run G1-b / G3 / G6-6. If the "OB-71" and "OB-76" cited by the user are to be registered, please assign their numbers (suggested: the next free number after OB-71).
- **D-5**: whether the cross-slice backfill item (the commit hash and CI run number in the DIRECTIVE-CLEANUP-EVAL row of `REMAINING_SLICES`) is merged into this slice (a one-line change is prepared, to be landed upon authorization) or left to the next slice.
- **D-6**: `DIRECTIVE-CLEANUP-EVAL` did not land the §1.7.3 fixed-path freedom list (found by read-only measurement; per the constraint this slice did not register an OB) — whether to remedy it or register it.

## 8. Delivery surface and per-item landing points

| Item | Landing point |
|---|---|
| Assessment of ① / ② / ③ | This report §3 / §4 / §5 |
| ④ three candidates and criteria | This report §6 (**no candidate is landed into the directive**) |
| ⑤ hand for ruling | This report §7 and this slice's ⑩ receipt |
| OB registration | The OB-71 row and numbering-source note in §12.4 of `docs/DELIVERY_DIRECTIVE{,_EN}.md`; version package v1.20 (header and §0.3 line) |
| ADR | ADR-021 in `docs/ADR_REGISTER{,_EN}.md` (⑦ downgrade; four elements plus not to be used as a precedent) |
| Ledger | The completion-record row of `docs/t027/REMAINING_SLICES{,_EN}.md` |
| New-character registration | `tools/gates/cjk-newwords.txt` ("瞻", with the reason and context excerpt entered into the table) |
| Freedom list | `docs/validation/evidence/DIRECTIVE-DESIGN-LOOPBACK-freedom-list.md` |
| Not landed (pending ruling) | ①'s trigger criterion, ②'s landing point, ③'s mechanism form |

## 9. Execution pipeline and environment

① no dispatch (pure assessment slice, see F-4) ⇒ ② documents and ledger written to disk ⇒ ③ checklist verification and consistency ⇒ ④ gates ⇒ ⑤ **skipped** (§5.2's S tier is "optional") ⇒ ⑥ **skipped** (pure `.md`, no code semantics; the §7.2 skip condition holds, confirmed by ⑦) ⇒ ⑦ final review (`deepseek-v4-pro`, **ADR-021 downgrade**; **reduced independence noted truthfully**) ⇒ ⑧ report and write-back ⇒ ⑨ native verification ⇒ ⑧c close-out cleanup ⇒ ⑩ stopping point.

**Quota actuals**: probes 0; realized paid calls = one ⑦ final-review round (`FINAL REVIEW: FINDINGS`) + one ⑦ re-review round (the mandatory post-remediation re-run of §5.3, counted in the §7.5 re-review quota) + one ⑧ mirror round.
Environment: Windows local machine; Node v24.17.0; Go toolchain `go1.27.0 windows/amd64`; repo `github.com/larsonzh/prfrail`; branch `main`, starting point `8d0ba88`.

### 9.1 Exception and fallback records (a report element of §11.3)
- **Landed before stopping (recorded per the user ruling of 2026-10-03)**: with the numbering provenance unclear (the "OB-71" and "OB-76" cited by the user were unregistered in the ledger) this slice **did not stop to check first**, but straight away took OB-71 in ledger order and added the numbering-provenance note - strictly speaking that is a §1.7.2 escalation boundary (stop when in doubt). **Improvement**: from now on, whenever a number is taken or a §12.4 row is added, **if the ledger disagrees with the user citation the master must stop and check first** and must not take the next number on its own. (Ruling D-4 keeps this slice factual outcome and orders no gate re-run, while requiring this lesson to be recorded.)
- **Tool-side exceptions**: (1) `create_file` wrote the report pair with CRLF and ⑧ normalised them to LF ⇒ the byte count differs from the first draft (content unchanged, verified line by line); (2) the PowerShell terminal produced a false interrupt on long one-line `node -e` commands (echoing `^C` while the command had in fact run) ⇒ fallback = always move such logic into a `tmp/ddl/*.js` script and read the output file back each time.
- **No other fallbacks**: the remaining §9.1 scenarios did not trigger in this slice (no budget overrun, no dependency block, no repeated native failure).
## 10. Review conclusion

- ⑤: **skipped** (S tier optional).
- ⑥: **the skip holds** (checked item by item and confirmed by ⑦): (1) the change set is all `.md` / `.txt` wording files; (2) no `.go` / `.json` / `.yml` / schema / fixtures / script-logic change; (3) no state-machine / gate-criterion / evidence-model / role-authorization semantics is touched (`cjk-newwords.txt` only registers a character through the existing G4b mechanism; the criterion and its scope are unchanged).
- ⑦: **conclusion in §12** (this slice's ⑦ is undertaken by `deepseek-v4-pro` ⇒ reduced independence, truthfully noted in ADR-021 and this section).

### 10.1 Review-input package summary (with blind-review isolation proof)
Change-set diff = `tmp/ddl/change-set.diff` (with `tmp/ddl/numstat.txt`); contract sentences = §0.2 / §6.3 / §7.2 / §11.3 / §12.4 of `docs/DELIVERY_DIRECTIVE.md`; freedom list = `docs/validation/evidence/DIRECTIVE-DESIGN-LOOPBACK-freedom-list.md`; **blind-review isolation proof**: both ⑤ and ⑥ were skipped ⇒ there is no intermediate list to attach, and ⑦ round 1 received only the slice definition, the diff, the contract sentences and the read-only constraint (§7.3); input sanitisation = every item is repository text, with no credential, token or SSH path.
## 11. Falsifiability and gate results

| Check | Result |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` (changed=10) |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0` (changed=10) |
| `gate.js --selftest` | `SELFTEST: PASS 264/264` |
| Base gates (§6.1) | `gofmt` empty; `go build` and `go vet` and `go test` all exit 0 (ok=14, FAIL=0) |
| Mutation check (G6-6 closed set) | mutated run `TOTAL_FAIL=1` (`docs/DELIVERY_DIRECTIVE.md:874` first status word turns red); after restoration the same command gives `TOTAL_FAIL=0`, and the restored file sha256 matches the pre-mutation state (evidence: `tmp/ddl/mutation-g6-6.txt`, `tmp/ddl/mutation-restored.txt`, `tmp/ddl/mutation-sha256.txt`) |
| Mirror positionality (report pair) | lines / bold / bars / per-line shape all equal, 0 mismatches (reading in `tmp/ddl/mirror-parity.txt`) |
| Checklist verification (§2) | Eleven items marked item by item, of which C-5 and C-6 and E-9 and X-1 and X-2 are handled |

**Candidate OBs (this slice lands only one per the constraint; the rest are unregistered)**: ① `DIRECTIVE-CLEANUP-EVAL` lacks the §1.7.3 fixed-path freedom list (found by read-only measurement this round); ② the "backfilled by the next slice" wording in `REMAINING_SLICES` will keep producing cross-slice lag items (this slice itself is affected, see B-2 in §2). Both await the user's ruling on whether to register them.

### 11.1 ⑧c wrap-up clean-up record (a fixed step of §5.1, executed by the orchestrator)
`gate.js --all --scope=tree` all green; the encoding gate (G4a) checked every changed file against the `.md` = BOM+LF rule with no violation; the diagnostic scan (§3.4) produced no new warning; workspace hygiene = all 10 entries staged with nothing unstaged or untracked (reading in `tmp/ddl/eightc.txt`); `tmp/ddl/` is a one-off working directory (gitignored) whose local evidence is cleaned up per §1.5 after user review and commit.
## 12. ⑩ Stopping-point report (including write-back check lines)

- Executed up to ⑨; **stopped before commit authorization**: this slice does not commit, does not push, does not push to gitee, and waits for the user's same-round authorization.
- **Write-back check** (the OB-63 formulation of §11.2; this slice is its second application): validation report = **written back** (this file); `DEV_PLAN` = **Not applicable** (a directive-governance slice; per the slice definition and existing precedent, no `DEV_PLAN` paragraph is written); ledger = **written back** (`REMAINING_SLICES{,_EN}` completion-record row and the OB-71 entry in §12.4 of `DELIVERY_DIRECTIVE{,_EN}`); ADR = **written back** (ADR-021).
- **⑦ conclusion**: round 1 `FINAL REVIEW: FINDINGS` (7 Medium + 6 Low, all remediated); round 2 (re-review) `RE-REVIEW: PASS WITH FIXES` (5 Low, closed at ⑧b finalisation; ⑦ states that no round 3 is needed, per §7.5 and §7.8.3).
- **Gate readings**: both `--scope=tree` and `--scope=index` report `TOTAL_FAIL=0` (changed=10); `SELFTEST: PASS 264/264`; the base gates are all green (see §11).
- **A defect this slice produced itself (⑦ round-1 Medium, remediated)**: the report draft wrote the gate reading as `changed=7` (= the staging size when the report landed) while the evidence files say `changed=10` (including the report pair and the freedom list) ⇒ three readings were corrected; that instance is itself an example of the gap described in ① ("no one reconciles session judgement against the repository state").

### 12.1 Ruling receipt after the ⑩ stop point (a factual record)
**This section is a factual record and does not rewrite the ⑩-stop wording of §7.** The user issued six rulings at the ⑩ stop point: **D-1** adopt the "conclusion flow-back trigger" but reduce it to a single **principle** (the master proactively suggests, in its receipt, whether a conclusion should be registered as an OB / ADR; silence means no registration, yet the suggestion itself must leave a trace), and do **not** fold it into the ③ mechanism; **D-2** land the ② pointer in §1 of `DOCUMENTATION_PLAN`, **not in this slice**; **D-3** take the **light** form for ③ with the first anchor = **end of T027**, and confirm "outlook = master question list + user ruling" (medium and heavy await a period of running the light form); **D-4** **keep and ratify** the OB-71 occupancy and the numbering-provenance note and **re-run no gate**, while requiring the "landed before stopping" lesson to be recorded (see §9.1); **D-5** fold the previous slice CI back-fill **into this slice** (already landed in the ledger); **D-6** do **not** re-create the previous slice missing freedom list - register it as a known absence.
**Not executed (per the user constraint that this slice touches no rule body)**: the D-1 principle, the D-2 cross-layer pointer and the D-3 light mechanism (template + anchors) are **all left out of the rule text** for a later slice; the D-6 note ("the previous slice freedom list is missing, confirmed, never to be back-filled") is registered in the next-slice start checklist.
## 13. Artifacts list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-design-loopback.md	repo
docs/validation/directive-design-loopback_EN.md	repo
docs/validation/evidence/DIRECTIVE-DESIGN-LOOPBACK-freedom-list.md	repo
tools/gates/cjk-newwords.txt	repo
tmp/ddl/change-set.diff	local-only
tmp/ddl/check-cn.js	local-only
tmp/ddl/eightc.js	local-only
tmp/ddl/eightc.txt	local-only
tmp/ddl/fence-sync.js	local-only
tmp/ddl/fence.js	local-only
tmp/ddl/fix1.js	local-only
tmp/ddl/fix10.js	local-only
tmp/ddl/fix1b.js	local-only
tmp/ddl/fix2.js	local-only
tmp/ddl/fix3.js	local-only
tmp/ddl/fix4.js	local-only
tmp/ddl/fix5.js	local-only
tmp/ddl/fix6.js	local-only
tmp/ddl/fix7.js	local-only
tmp/ddl/fix8.js	local-only
tmp/ddl/fix9.js	local-only
tmp/ddl/g4b.txt	local-only
tmp/ddl/gate-index.txt	local-only
tmp/ddl/gate-tree.txt	local-only
tmp/ddl/land.js	local-only
tmp/ddl/ledger-d5.diff	local-only
tmp/ddl/ledger.diff	local-only
tmp/ddl/mirror-parity.txt	local-only
tmp/ddl/mirror.js	local-only
tmp/ddl/mutate-original.bin	local-only
tmp/ddl/mutate.js	local-only
tmp/ddl/mutation-g6-6.txt	local-only
tmp/ddl/mutation-restored.txt	local-only
tmp/ddl/mutation-sha256.txt	local-only
tmp/ddl/numstat.txt	local-only
tmp/ddl/peek-ledger.js	local-only
tmp/ddl/peek2.js	local-only
tmp/ddl/peek4.js	local-only
tmp/ddl/review-r1.md	local-only
tmp/ddl/scan.js	local-only
tmp/ddl/selftest.txt	local-only
```

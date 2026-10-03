# Slice validation report · `DIRECTIVE-LOOPBACK-LANDING` (lightweight loopback mechanism landing)

> Per §5.1 serial pipeline and §6.2 evidence-to-disk obligation of `docs/DELIVERY_DIRECTIVE.md` v1.20 (the starting version of this slice; this slice lands a v1.21 version package). **The CN version is authoritative**, `_EN` is the positional mirror; for the fence block see "Artifacts list" at the end. Size S, tier `[SLICE]`, probes 0. **The starting message was corrected by the user in the same round**: the D-3 template landing point is changed to the `docs/` root (a project-level tool, not exclusive to T027), and ① decision item 4 is added (whether the new template is registered into the document list in the same batch).

## 1. Change summary (including the execution freedom list)

**Goal**: land four groups of lightweight mechanisms — D-1 the driver proactively suggesting a registration principle, D-2 the `DOCUMENTATION_PLAN` cross-layer pointer, D-3 the closing-report template + T027 close-out anchor, and C the ledger completion-line format. **The rule body is not touched** (§2.2 / §2.4 / §5.2 / §6.3 / §1.7 rule semantics unchanged; ⑦ review).

**Execution freedoms exercised by this slice (§1.7.1, substantive choices)**:

- **F-1 version package**: this slice edits the §11.2 body ⇒ per §0.2 bump to **v1.21** (header version number + §0.3 change-log line + Appendix C.0k ledger block). It is a ledger-metadata-class by-product.
- **F-2 landing point of D-1 and C (① decisions 1 and 2)**: both land in **§11.2**, **each in its own separate paragraph** (a bold topic label at the line start), **no new `### 11.2.x`** (to avoid disturbing the numbering system and the heading line-number sequence).
- **F-3 template landing point (① decision 3 voided by the user's same-round correction)**: land in **`docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`** (the `docs/` root, sibling to `DELIVERY_DIRECTIVE` / `DOCUMENTATION_PLAN` / `CONTRACTS`). Reason: this template is a **project-level governance tool**, not exclusive to T027; **the T027 close-out is only the first application anchor, not the basis for the template's ownership**. ① The first-round ruling of `docs/t027/` is thereby voided and truthfully retained; `docs/validation/` is still excluded for falling in the G7 scan scope.
- **F-4 the OB-49 disposition column is changed per the user ruling to "replace the old slice name" rather than "add a sentence and keep the old name"**: ① the original proposal was to add "archive three sources" while keeping `DIRECTIVE-LEDGER-ARCHIVE`; this slice, per the verbatim text of the user's slice definition ("change the old slice name in the disposition column to wording consistent with the assessment direction"), **deletes the old slice name** and points instead to the "`tmp/` lifecycle assessment" slice while listing three disposition classes side by side. This difference is declared here.
- **F-5 this slice's ledger completion line applies the new format for the first time**: it records only the report path, **without the commit hash or the CI run number**. This first instance immediately exposes C's residual self-reference problem (the commit hash and the CI number both fall under "unknowable before commit") — truthfully logged with a recommendation attached (see the end of §7).
- **F-6 §11.2 is expanded from one paragraph to three** (write-back formulation / registration suggestion and trace / ledger completion-line format): no new numbered subsection, consistent with decisions ①②.
- **F-7 ⑦ downgrade**: this slice's ⑦ uses `deepseek-v4-pro` (explicitly stated in the user's slice definition), traced per the OB-37 meta-rule as **ADR-022**.
- **F-8 D-2 registers the new template in the same batch (① decision 4) and its consequential items**: per ① decision, add a template row to the §2 document list of `DOCUMENTATION_PLAN{,_EN}`; as a consequence, change the **currently effective** "12 document pairs" in that file's body to "13 pairs" (the "baseline 12 pairs" in §3 is a record of an already-executed stage and is not changed). This item is a consequential completion not listed in ① decision, handled minimally to avoid "the map is changed but the count still stands", and is declared here and handed to ⑩ for a ruling.

## 2. Pre-start checklist verification (step 0; item by item)

| # | Item | Status |
|---|---|---|
| 1 | D-1 principle current-state inventory (the directive has no such principle) | **Done** (landed in the new §11.2 paragraph) |
| 2 | D-2 cross-layer pointer current-state inventory (cross-layer references near zero) | **Done** (landed after the §1 table in `DOCUMENTATION_PLAN{,_EN}`) |
| 3 | D-3 template and anchor current-state inventory (no template, no anchor line) | **Done** (new project-level template pair, landed at the `docs/` root + the item 4 anchor in the `REMAINING_SLICES{,_EN}` usage section) |
| 4 | C ledger completion-line format current-state inventory | **Done** (landed in the new §11.2 paragraph; this slice's completion line applies it for the first time) |
| 5 | OB-49 pointer correction (old slice name → assessment direction) | **Done** (one-line change, not made a separate item) |
| 6 | Ledger CI backfill (two places cumulative) | **Done**: the `DIRECTIVE-DESIGN-LOOPBACK` line backfilled with `a815970` and run `37107634723` (the run number is sourced from a `gh run view` reading whose `conclusion=success` and `headSha` prefix `a815970` match the ledger, archived as `docs/validation/evidence/ci-run-37107634723.txt` (a `repo` disposition); this slice did not re-run CI); this slice's new line records only the report path per the new format |
| 7 | OB-71 status flip `To be discussed → Resolved` | **Done** (landing D-3 is its disposition) |
| 8 | Previous slice's leftover `DIRECTIVE-CLEANUP-EVAL` line | **Not applicable** (the previous slice's D-5 already backfilled it; this slice only re-reads and verifies) |
| 9 | `DIRECTIVE-MODELID-MAP`, OB-68, OB-53⑤ | **Not applicable** (outside this slice's scope) |
| 10 | Residue under the `tmp/` root | **Not applicable** (left for the OB-49 project initiation; this slice only changes its pointer) |

## 3. ① Decisions and landing points (solution document `DESIGN: DONE`; including the supplementary decision round after the user's correction)

| # | Decision item | Conclusion |
|---|---|---|
| 1 | Should D-1 land in §1.7 or §11.2 | **§11.2** (registration suggestions belong to the "write-back / registration" domain; §1.7 is the freedom and escalation boundary, and inserting it would break its semantic boundary) |
| 2 | Should C and D-1 share a paragraph or be separate paragraphs | **Same section, each in its own separate paragraph** (the two topics are semantically unrelated; no new numbered subsection, to avoid disturbing the heading line-number sequence) |
| 3 | D-3 template path | **the original ruling of `docs/t027/` is voided** (the first-round conclusion is truthfully retained); changed per the user's same-round correction to **`docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`** (the `docs/` root, project-level, generic form) |
| 4 | CN / `_EN` bilingual symmetry plan | every CN change lands isomorphically at the same logical position in `_EN`; bold aligned sentence by sentence, table column structure unchanged, blank-line rhythm consistent, ID counts equal on both sides (measured: all five shapes equal). After the template moves to the `docs/` root it is **no longer covered by the G7 scan scope**, and the fence obligation is carried by this report's `repo` disposition row |
| 5 | Whether the new template is registered into the `DOCUMENTATION_PLAN` document list in the same batch (added in the correction round) | **Registered** (the document map must list all project-level documents, avoiding the gap where "a document is on disk but not on the map"; D-2 already touches that file in this batch, saving a round trip) ⇒ one row in its §2 list, +1 line each for CN / `_EN` |

## 4. Delivery surface and per-item landing points

| Item | Landing point |
|---|---|
| D-1 principle + first instance | The new "registration suggestion and trace" paragraph in §11.2 of `docs/DELIVERY_DIRECTIVE{,_EN}.md` |
| C ledger completion-line format | The new "ledger completion-line format" paragraph in the same section |
| D-2 cross-layer pointer | After the §1 authority-domain table in `docs/DOCUMENTATION_PLAN{,_EN}.md` |
| D-3 template | New `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md` (project-level, generic form; retrospective = the driver; prospective = the driver's list + the user ruling) |
| D-2 same-batch registration | A new closing-report template row in the §2 document list of `docs/DOCUMENTATION_PLAN{,_EN}.md` (① decision 5) |
| D-3 anchor | Item 4 of the usage section in `docs/t027/REMAINING_SLICES{,_EN}.md` |
| OB-71 / OB-49 | The two disposition-column rows in §12.4 of `docs/DELIVERY_DIRECTIVE{,_EN}.md` |
| Version package | Header version v1.21 + §0.3 line + Appendix C.0k ledger block |
| ADR | ADR-022 in `docs/ADR_REGISTER{,_EN}.md` |
| Freedom list | `docs/validation/evidence/DIRECTIVE-LOOPBACK-LANDING-freedom-list.md` |

## 5. Execution pipeline and environment

① architecture (V4 Pro, **two rounds**: first-round decision items 1 to 3 + supplementary decision items 4 and 5 after the user's correction) ⇒ ② document implementation ⇒ ③ checklist verification and consistency ⇒ ④ gates ⇒ ⑤ **skipped** (S tier optional) ⇒ ⑥ **skipped** (pure `.md`, no code semantics; the §7.2 skip condition holds, confirmed by ⑦) ⇒ ⑦ final review (`deepseek-v4-pro`, **ADR-022 downgrade**; **reduced independence noted truthfully**) ⇒ ⑧ report and `_EN` mirror ⇒ ⑨ native verification ⇒ ⑧c close-out cleanup ⇒ ⑩ stopping point.

**Quota actuals**: probes 0; realized paid calls = ① solution document **two rounds** (first round + the supplementary decision after the user's correction; **exceeding §7.5 "① ≤ 1"**, truthfully logged), ⑦ one final-review round and one re-review round (if triggered), ⑧ one mirror round.
Environment: Windows local machine; Node v24.17.0; Go toolchain `go1.27.0 windows/amd64`; repo `github.com/larsonzh/prfrail`; branch `main`, starting point `a815970`.

### 5.1 Exception and fallback records (a report element of §11.3)

- **Encoding anomaly when writing new `.md` files to disk (tool side)**: `create_file` wrote the two template files with CRLF, and during staging the repository normalisation converted them to LF / `BOM` normalisation (this phenomenon matches the same kind of behaviour of `create_file` in the previous slice) ⇒ fallback = after every new file creation **explicitly re-check BOM / LF / line count** and normalise uniformly with a script; the two template files of this slice were already normalised (BOM + LF) before verification.
- **Check-order anomaly**: the first shape verification read the template pair as asymmetric in blank lines "22 / 1", because the verification script counted before normalising CRLF ⇒ fallback = normalise line endings before verification, then align line by line; after normalisation, the template pair has **the same element count on both sides (44 each under the `split(/\n/)` reading), 22 blank lines each, and fully equal heading line-number sequences**.
- **Starting message corrected mid-flight ⇒ ① quota overrun (truthfully logged)**: this slice's ① cumulative calls are **2** (first-round decision items 1 to 3 + the supplementary decision after the user's same-round correction). The ① architect's cumulative calls for this slice are 2, exceeding the §7.5 single-slice ≤1 cap, **driven by the user's same-round correction instruction**, executed under the "exceeding the cap ⇒ escalate to the user for a ruling" formulation, **limited to this slice only and not to be taken as precedent**; that escalation was itself completed by the user's same-round correction, so no separate OB / ADR is raised. Both rounds of calls are counted in this section and in "Quota actuals".
- **G6-4 red flags and correction (no whitelist used)**: after the template moved to the root, the first gate run flagged G6-4 red twice (line-number literals in this report and in the freedom list); **that red reading was not archived on its own** (later re-runs overwrote it) ⇒ rewritten to the "element count / blank-line count" wording and green again, with the post-fix reading in `tmp/ll/pre.txt` (added-line scope, `ADDED-LINE HITS=0`) and `tmp/ll/lit.txt` (whole-file scope, every hit on a historical line), **no whitelist exemption registered**.
- **G2 placeholder residue flagged a line added by the remediation**: a §8 bullet added by the ⑦ round-1 remediation carried a Chinese placeholder phrase that starts with the character「待」and then spells out the same-turn authorisation wording, which trips the G2 placeholder vocabulary ⇒ reworded and green again, **with no whitelist registered**; this red is the second instance in this slice of "a remediation introducing a new red", also recorded in §3 of the freedom list.
- **Consequential-item self-check**: after the template moved to the root, check every path reference (the v1.21 change-log entry in §0.3 of the directive, the OB-71 disposition column in §12.4, the Appendix C.0k deliverable row, the ledger usage section and completion-record row, and this report's fence), and per ① decision complete the list row and current count of `DOCUMENTATION_PLAN`.

## 6. Review conclusion

- ⑤: **skipped** (S tier optional).
- ⑥: **the skip holds** (checked item by item and confirmed by ⑦): ① the change set is entirely `.md` wording; ② no `.go` / `.json` / `.yml` / schema / fixtures / script logic; ③ no state machine / gate criterion / evidence model / role-authorization semantics is touched.
- ⑦: **conclusion in §8** (this slice's ⑦ is undertaken by `deepseek-v4-pro` ⇒ reduced independence, truthfully noted in ADR-022 and this section).

### 6.1 Review-input package summary (with blind-review isolation proof)

Change-set diff = `tmp/ll/change-set.diff` (with `tmp/ll/numstat.txt`); contract sentences = §0.2 / §6.3 / §7.2 / §7.5 / §11.2 / §11.3 / §12.4 of `docs/DELIVERY_DIRECTIVE.md`; the summary of ①'s two decision rounds follows §3 of this report (the first-round decision 3 was voided by the user's correction); freedom list = `docs/validation/evidence/DIRECTIVE-LOOPBACK-LANDING-freedom-list.md`; **blind-review isolation proof**: both ⑤ and ⑥ were skipped ⇒ there is no intermediate list to attach, and ⑦ round 1 received only the slice definition, the diff, the contract sentences and the read-only constraint (§7.3); input sanitisation = every item is repository text, with no credential, token or SSH path.

## 7. Falsifiability and gate results

| Check | Result |
|---|---|
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` (changed=13) |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0` (changed=13) |
| `gate.js --selftest` | `SELFTEST: PASS 264/264` (this slice adds no fixture) |
| Base gates (§6.1) | `gofmt` empty; `go build` and `go vet` and `go test` all exit 0 (ok=14, FAIL=0) |
| Mirror positionality (including the five shapes) | reading in `tmp/ll/mirror.txt` |
| Mutation check | reading in `tmp/ll/mutation.md` (`docs/DELIVERY_DIRECTIVE.md:879` turns red, a renumbered citation turns red, an asymmetric template pair turns red; all green again after item-by-item sha256 verification and byte-for-byte restoration; the file also carries a **"final re-run" closing block**) |
| G6-4 correction | the first run flagged red twice (line-number literals, the red reading not archived on its own) ⇒ green again after rewording, **no whitelist used** (post-fix reading in `tmp/ll/pre.txt` and `tmp/ll/lit.txt`) |
| Checklist verification (§2) | Ten items marked item by item, of which 1 to 7 are handled and 8 to 10 not applicable |

**Falsifiability boundary of C and D-1 (truthful)**: both are format and behaviour norms with **no new mechanical criterion** (this slice's hard constraint adds no criterion) ⇒ their verification = ⑦ human re-review + the ⑩ write-back check line + G2 / G6-4 corroboration; **the mutation check does not cover them**, and this must not be read as "guarded by a criterion".

**C's residual self-reference problem (this slice's first-instance observation + recommendation)**: C's intent is to eliminate the cross-slice backfill caused by "the CI run number is unknowable before commit", but **the commit hash is likewise unknowable before commit** (the completion line is written inside the commit) ⇒ this slice's completion line can record only the report path, and the hash still needs backfilling by the next slice. **Recommendation** (raised per the D-1 principle this slice has just landed, handed to the user for a ruling): the completion line should instead record "**the HEAD (BASE) hash at write time** + the report path", or change to "**a ledger follow-up commit by the same round after commit**", either of which makes the completion line determinable at write time.

**Candidate OBs (not registered by this slice; per the constraint only the two disposition-column entries OB-71 / OB-49 are landed)**: ① the BASE-hash formulation of the completion line (see above); ② the style observation "whether §11.2 should be split into sections after expanding to three paragraphs" (① already flagged the rising topic density).

### 7.1 ⑧c close-out cleanup record (a fixed step of §5.1, executed by the driver)

`gate.js --all --scope=tree` all green; the encoding gate (G4a) checked every changed file against the `.md` = BOM+LF rule file by file, with no violation; the diagnostic scan (§3.4) produced no new warning; workspace hygiene = every change staged file by file, with **no git-visible unstaged or untracked entries** (`tmp/ll/` being a gitignored directory) (reading in `tmp/ll/eightc.txt`); `tmp/ll/` is a one-off working directory (gitignored) whose local evidence is cleaned up per §1.5 after user review and commit.

## 8. ⑩ Stopping-point report (including the write-back check line)

- Executed up to ⑨; **stopped before commit authorization**: this slice does not commit, does not push, does not push to gitee, and waits for the user's same-round authorization.
- **Tier and rework**: size S, tier `[SLICE]`; rework = one extra ① round (driven by the user correction, logged) + one G6-4 correction round + one ⑦ first-round remediation round; no timeout, no blocker.
- **Write-back check** (the OB-63 formulation of §11.2): validation report = **written back** (this file); `DEV_PLAN` = **Not applicable** (a directive-governance slice; per the slice definition and existing precedent, no `DEV_PLAN` paragraph is written); ledger = **written back** (the T027 anchor row and completion-record row of `REMAINING_SLICES{,_EN}`, and the OB-71 and OB-49 disposition columns of §12.4 in `DELIVERY_DIRECTIVE{,_EN}`); ADR = **written back** (ADR-022).
- **Gate readings**: both `--scope=tree` and `--scope=index` report `TOTAL_FAIL=0`; `SELFTEST: PASS 264/264` (see §7).
- **Explicitly not executed**: the OB-49 cleanup itself (this slice changes the pointer only); the residue under the `tmp/` root; retroactive rewriting of the historical count in `DOCUMENTATION_PLAN` §3; this slice's commit and push (awaiting same-turn authorisation).
- **Next steps**: (1) a user ruling on the C BASE-hash formulation (⑦ suggests registering it as OB-72); (2) the next slice opens the OB-49 item to assess the `tmp/` lifecycle and land the three handling classes; (3) ⑧b records the post-re-review remediation in this section; ⑤ and ⑥ are not applicable (skipped) and ⑦ and ⑧ follow this slice's tier.
- **⑦ conclusion**: round 1 `PASS WITH FIXES` (1 Medium and 6 Low; the Medium being the §0.3 ④ mismatch with the §12.4 OB-49 wording, fixed) ⇒ remediation ⇒ re-review `PASS WITH FIXES` (no Medium or above; 3 Low fixed in place by ⑧b: the EN `->`, the evidence-file NOTE and caliber, and the CI run sourcing); ⑦ is `deepseek-v4-pro` downgraded (ADR-022), **reduced independence truthfully noted**, with the final ruling resting with the user and ⑧c.
### 8.1 ⑦ two rounds of findings and their remediation (§5.3)

| # | ⑦ finding | Disposition |
|---|---|---|
| F-1 | the ④ item of the v1.21 §0.3 change log disagreed with the landed OB-49 cell in §12.4 (Medium) | reworded to "repointed + three handling classes", aligned with the Appendix C.0k note |
| F-2 | a duplicated clause and `->` in the EN OB-49 row, plus a truncated CN parenthetical (Low) | both sides rewritten; the CN parenthetical now reads "this item is listed as a next-slice assessment object" |
| F-3 | the mutation evidence ended on an intermediate red state (Low) | a "final re-run" closing block and NOTE lines were added; all three VERIFY blocks read `TOTAL_FAIL=0` |
| F-4 | numstat and gate readings slightly staler than the index (Low) | every reading regenerated from the post-remediation index |
| F-5 | the G6-4 correction item cited a reading file that already showed green (Low) | now states "the red reading was not archived on its own" and cites `tmp/ll/lit.txt` |
| F-6 | §8 lacked tier and rework, not-executed items and next steps (Low) | three bullets added (see above in this section) |
| F-7 | §7.1 "nothing untracked" clashed with the readings listed alongside it (Low) | limited to "no git-visible unstaged or untracked entries", noting that `tmp/ll/` is gitignored |
| F-8 | re-review: two `->` remained in the EN OB-49 row (Low) | both changed to `⇒`, matching the CN symbol |
| F-9 | re-review: `tmp/ll/lit.txt` is a whole-file scan with no NOTE, easily misread as a gate red (Low) | a NOTE header was added and the report now also cites `tmp/ll/pre.txt` (added-line scope) |
| F-10 | re-review: the CI run citation was self-referential (Low) | a `gh run view` reading was added and archived as `docs/validation/evidence/ci-run-37107634723.txt` (`repo`) |
### 8.2 Receipt of the rulings made after the ⑩ stop point (factual record)

This section is a factual record and **does not rewrite the ⑩-stop-point wording of §8** (per §11.1).

| # | User ruling (2026-10-03, same round after the ⑩ stop point) | Disposition |
|---|---|---|
| 1 | the ① quota overrun (2 calls) is **accepted**, with no restart and no escalation; the traceability wording "limited to this slice, not a precedent, no separate OB / ADR" is accepted as well | accepted; the wording in §5.1 and in "Quota actuals" is unchanged |
| 2 | the BASE-hash self-reference in C is **registered as OB-72** (first status word `To be discussed`) | registered in §12.4 of `DELIVERY_DIRECTIVE{,_EN}` and folded into the v1.21 version bundle |
| 3 | OB-49 is **not opened now**, kept as a next-slice candidate (suggested to follow `DIRECTIVE-MODELID-MAP`) | the pointer stays (done in this slice); the opening timing is out of this slice |
| 4 | commit and push origin authorised (staged file by file, no gitee push); `tmp/ll/` is cleaned after the commit, with the CI run reading archived as `docs/validation/evidence/ci-run-37107634723.txt` | the commit and the push appear in the ledger completion row; the archived file enters the §9 fence with a `repo` disposition |

**This slice is itself the live instance of OB-72**: the commit hash and the CI run number of this section are equally unknowable at write time and can only be completed by the next slice or by a revision of this section.

## 9. Artifacts list

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/DOCUMENTATION_PLAN.md	repo
docs/DOCUMENTATION_PLAN_EN.md	repo
docs/SLICE-CLOSEOUT-TEMPLATE.md	repo
docs/SLICE-CLOSEOUT-TEMPLATE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-loopback-landing.md	repo
docs/validation/directive-loopback-landing_EN.md	repo
docs/validation/evidence/DIRECTIVE-LOOPBACK-LANDING-freedom-list.md	repo
docs/validation/evidence/ci-run-37107634723.txt	repo
tmp/ll/apply.js	local-only
tmp/ll/change-set.diff	local-only
tmp/ll/chk.js	local-only
tmp/ll/ci-run.txt	local-only
tmp/ll/cn-dump.txt	local-only
tmp/ll/dump.js	local-only
tmp/ll/dump2.js	local-only
tmp/ll/dump3.js	local-only
tmp/ll/dump4.js	local-only
tmp/ll/dump5.js	local-only
tmp/ll/eightc.js	local-only
tmp/ll/eightc.txt	local-only
tmp/ll/en-dump.txt	local-only
tmp/ll/en-full.txt	local-only
tmp/ll/en-new.md	local-only
tmp/ll/finalize.js	local-only
tmp/ll/fixmut.js	local-only
tmp/ll/g3a.txt	local-only
tmp/ll/gate-index.txt	local-only
tmp/ll/gate-tree.txt	local-only
tmp/ll/gsum.js	local-only
tmp/ll/gsum2.js	local-only
tmp/ll/land1.js	local-only
tmp/ll/land10.js	local-only
tmp/ll/land11.js	local-only
tmp/ll/land12.js	local-only
tmp/ll/land13.js	local-only
tmp/ll/land14.js	local-only
tmp/ll/land2.js	local-only
tmp/ll/land3.js	local-only
tmp/ll/land3.txt	local-only
tmp/ll/land4.js	local-only
tmp/ll/land4.txt	local-only
tmp/ll/land5.js	local-only
tmp/ll/land5.txt	local-only
tmp/ll/land6.js	local-only
tmp/ll/land7.js	local-only
tmp/ll/land7.txt	local-only
tmp/ll/land8.js	local-only
tmp/ll/land8.txt	local-only
tmp/ll/land9.js	local-only
tmp/ll/lit.js	local-only
tmp/ll/lit.txt	local-only
tmp/ll/llmirror.js	local-only
tmp/ll/m1.orig	local-only
tmp/ll/m1.txt	local-only
tmp/ll/m2.orig	local-only
tmp/ll/m3.orig	local-only
tmp/ll/mir.js	local-only
tmp/ll/mir2.js	local-only
tmp/ll/mirror.txt	local-only
tmp/ll/mut.js	local-only
tmp/ll/mutation.md	local-only
tmp/ll/numstat.txt	local-only
tmp/ll/pairs.js	local-only
tmp/ll/pre.txt	local-only
tmp/ll/pre2.txt	local-only
tmp/ll/preflight.js	local-only
tmp/ll/recon.js	local-only
tmp/ll/recon2.js	local-only
tmp/ll/recon2.txt	local-only
tmp/ll/rev2-en-directive.txt	local-only
tmp/ll/selftest.txt	local-only
tmp/ll/tail-dump.txt	local-only
tmp/ll/tail2.txt	local-only
tmp/ll/tpl-cn.txt	local-only
tmp/ll/tpl-en.txt	local-only
tmp/ll/tpl.js	local-only
tmp/ll/tpl.txt	local-only
```

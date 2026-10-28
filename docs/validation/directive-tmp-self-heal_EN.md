# ProofRail P34 governance-slice verification report: temp-artifact self-heal exception (directive v1.56)

Slice `DIRECTIVE-TMP-SELF-HEAL`; baseline `HEAD` = `53ecf7f`; directive version v1.56; size M (governance slice); report date 2026-10-28.

## 1 Change summary (including exceptions and fallbacks)

**Problem**: §5.1's current rule "self-delete before the controller's reading ⇒ non-self-heal tier" makes "create and delete in the same command" physically unable to self-heal; P32 (D-1 / D-2 / D-3) and P33 (D-1) have hit this gap three times in a row.

**Ruling**: the user ruled group by group on 2026-10-28 and chose "add a same-command create-and-delete exception under §5.1's self-heal tier"; ① produced the table only and took no position, stopping for the ruling.

**Result**: the directive is raised to v1.56 (header / §0.3 / the new §5.1 paragraph with two minimal closures / the §9.3 disposition-cell extension / Appendix C.0aq), and the P34 ledger block and boundary row are synced; all three gate scopes report `TOTAL_FAIL=0` and the self-test reports `SELFTEST: PASS 697/697`; no code or criterion body was changed.

| File | Landing point | Content |
| --- | --- | --- |
| Both directives | Header `:3` | `Version: v1.56. Date: 2026-10-28.` |
| Both directives | §0.3 `:102` | v1.56 change row |
| Both directives | §5.1 `:516` | the new "same-command create-and-delete exception (v1.56)" paragraph |
| Both directives | §5.1 `:518` and `:520` | two minimal closures: non-self-heal tier, per-slice cap |
| Both directives | §9.3 `:902` | both version markers in the disposition cell (extension) |
| Both directives | Appendix C `:2080` | the `C.0aq` slice block |
| Both ledgers | Status row and P34 block | status row synced, boundary row extended |

**Item 1 R = 2**: R1 = ② copied the task package's rewrite instructions `W3`–`W9` verbatim into the directive body (CN 7 places, EN 7 places); detection = mechanical verification on the controller side (self-reports were not taken on trust); repair = ② was sent back to rewrite them as self-contained labels; root cause = the task package never said "`W1`–`W13` are rewrite instructions, not body content" and the executing side did not judge it either; `G6-5` / `G6-7` cannot catch it (they judge only §X.Y and OB / SW / DR IDs). R2 = ⑤ r1's two Lows (the ledger boundary row's `:899` line number had drifted to `:902`, and `OB-128`'s touch surface missed §9.3) ⇒ ④ re-ran after the targeted fix.

**Item 2 D = 1**: during the rework round ② created 3 one-off scripts inside `tmp/DIRECTIVE-TMP-SELF-HEAL/` and deleted them itself; classification = **non-self-heal tier** (self-deletion preceded the controller's reading, so the reading is unobtainable); disposition = the user ruled to continue under the `OB-37` meta-rule (no slice halt), and it **must not be used as a precedent**.

**Item 3 meta finding**: the v1.56 exception being landed by this slice explicitly excludes the case "lawful location but ordering violation only"; **the first instance was triggered by this slice itself and falls exactly on the new rule's exclusion side** — not laziness on the executing side, but the first empirical proof that the v1.56 exclusion condition is too strict.

**Item 4 evidence gap (discipline B)**: lost = the bodies and stdout of the 3 deleted scripts (`p34-rework-relabel.js` / `p34-rework-postcheck.js` / `p34-rework-scope.js`); not lost = the two `.bak` images, the full text of `apply.js` / `verify.js`, the six shapes and gate readings after the change, and the full `git diff`. **No claim is made that the evidence gap has been closed.**

**Item 5 §12.1 citation precision**: the start-up message recorded "a governance slice does not meet the pilot selection criteria" as if it were wording of §12.1 itself, but it is the established wording of the §12.4 OB rows (the §12.1 body has no such literal sentence); this slice does not modify the §12.1 body (out of bounds) and only registers the fact, so that the v23 checklist can switch to "the established §12.4 OB-row wording" or spell out the citation chain.

**Item 6 ③'s `TEST: FAIL (1)` record**: the FAIL item = the task package's self-defined criterion "CN and EN have the same number of sentences"; false-positive basis = CN joins clauses with "，", EN splits them with ";", and the items correspond one by one; **compliance of entering ④ without re-running ③** = §5.3's re-run trigger for ③ is "code / test changed", and this slice changed neither code nor tests, so the trigger does not hold and all of ③'s substantive checks passed.

**Item 7 batch assistance**: this slice missed one arrangement (after ① and before the ruling); its nature = **an execution-scheduling error** (not "considered redundant"); 3 dispatches were then made up (enumeration / comparison / impact); registered as **the second instance of v22-OBS-C**.

**Item 8 `W1`–`W13` adoption and the X1 extension**: the item-by-item adoption list is in 1.1; the extension = the §9.3 "② cleans early" disposition cell plus the ledger boundary row, see the last paragraph of 1.1.

### 1.1 `W1`–`W13` adoption list and the X1 extension

| No. | Adopted content |
| --- | --- |
| `W1` | ④ is made precise as "no command or statement may follow the delete action" |
| `W2` | ⑤'s reason sentence becomes "`git status` never shows it + §5.1 zero tolerance" |
| `W3` | the ordering mandate and reading ② baseline are explicitly waived, limited to "when the exception holds" |
| `W4` | ①④ replace the pre-snapshot completeness check, ③'s path set = the full set of derived creation paths, ③ is executed by the controller along the established mechanism |
| `W5` | ①② are the "non-replayable command" evidence exception |
| `W6` | when the violator is the controller itself, the same-round independent review layer must confirm |
| `W7` | the three-path order is unchanged and the pre-⑩ terminal read-back is not waived |
| `W8` | events after ⑧c are also not exempt from the "no rebuild for rectification" ban |
| `W9` | "same-command create-and-delete + outside the repository" has no self-heal path, deliberately strict |
| `W10` | two minimal closures in §5.1: the trigger enumeration adds "any of the exception's five items missing", the cap adds "including exception evidence self-heal" |
| `W11` | the `C.0aq` first-entry wording is approved as "two rulings, three events" |
| `W12` | `C.0aq` adopts "the source of truth remains the two slice reports; the appendix is only a cross-slice summary trace" |
| `W13` | a fifth item is appended to the original four: "the create-and-delete location must be a repository path" |

**X1 extension**: inside the §9.3 "② cleans early" disposition cell the single word "self-heal tier" is extended to "self-heal tier (v1.45) or same-command create-and-delete exception (v1.56)", keeping both version markers and leaving the symptom column and all other cells untouched; the P34 boundary row in the ledger appends "+ §9.3 disposition cell (extension)".

**Item 9 three v23 observation candidates**:

1. The design of task-package self-defined criteria — this slice's ③ produced 1 false positive because the "CN and EN have the same number of sentences" criterion is not rigorous.
2. Whether the v1.56 exclusion condition is too strict — should "lawful location + same-command create-and-delete" be included in the exception (trigger = this slice's D = 1, the first empirical case).
3. The timing of batch-assistance arrangements — this slice missed one.

**Report snapshot note**: this report is a snapshot at ⑩; the 4 ledger-level fact deviations (§12.1 reference precision, gates-cache exception scope, ⑥ output contract `##` prefix, checker criterion shape) and 4 open items (v1.56 exclusion scope, ⑥ format deviation recurrence, ⑧ evidence-write path, task-pack self-defined criteria) appear in the ⑩ receipt's fact pack fields 10/12, and will be sunk into §1 and §11 of this report in the docs-only backfill.

**Ledger-level fact deviations (4 items, sunk in during the docs-only backfill)**: ① `§12.1` itself contains no literal sentence "a governance slice does not meet the pilot-selection standard" (that phrasing is the established wording of the `§12.4` OB rows); ② the runtime exception set under `tmp/gates-cache/` is only the `g4b/` subdirectory and does not cover the `gates-cache` root level; ③ the ⑥ output contract requires the five section headings to carry a `##` prefix, so a bare `Section A` line necessarily makes `G7-f` FAIL; ④ checker criterion shapes are not rigorous (③'s "CN and EN have the same number of sentences" and ⑨'s "heading-style § reference" each produced 1 false positive).

## 2 Execution pipeline

| Step | Role / carrier | Model | Status | Rework count |
| --- | --- | --- | --- | --- |
| 〇·zero and step 0 | Controller | Controller model | self-check 0 findings | 0 |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro` | three-group branch comparison table, stopping for the ruling | 1 |
| Batch assistance | read-only subagents ×3 | `GLM-5.3-Flash` | enumeration / comparison / impact | 3 |
| ② implementation | `prfrail-implementer` | `DeepSeek V4.1 Flash` | first landing + rework + targeted fix | 3 |
| ③ test | `prfrail-tester` | `DeepSeek V4.1 Flash` | `TEST: FAIL (1)` (criterion false positive, see §1 item 6) | 1 |
| ④ gates | Controller | gate scripts | three scopes `TOTAL_FAIL=0` | 2 |
| ⑤ pre-review | Pre-reviewer | `DeepSeek V4 Pro` | `PRE-REVIEW: PASS` | 1 |
| ⑥ independent scan | Independent scanner | see the ⑩ dispatch table | conclusion `PASS` (text in §5) | 1 |
| ⑦ final review | Independent final reviewer | see the ⑩ dispatch table | `FINAL REVIEW: PASS` | 1 |
| ⑧a report | `prfrail-documenter` | `DeepSeek V4.1 Flash` | this report | 1 |

Per-round durations were not timed and are uniformly marked `未计时`; the record gap is listed in §11 item 6.

## 3 Exception and fallback record

See §1: all nine items live in §1 items 1–9 (this chapter is merged with §1); this chapter only gives the counting basis.

| Count | Value | Basis |
| --- | --- | --- |
| R (fallback) | 2 | §9.2; R1 = ② copied `W3`–`W9` in the first round, R2 = ⑤ r1's two Lows |
| Temp-artifact D | 1 | §5.1; during the rework round ②'s 3 scripts created and deleted themselves, non-self-heal tier |
| Deviation-table numbers | 0 | no deviation other than the temp-artifact one in this slice |
| Probes | 0 | quota 0 |

## 4 Environment and tier

| Item | Value |
| --- | --- |
| Working directory | `d:\LZProjects\prfrail` |
| Platform | Windows (PowerShell terminal) |
| Baseline | `HEAD` = `53ecf7f` |
| Directive version | v1.56 |
| Tier | M (governance slice; ⑦ quota 1+3) |
| ⑥ and ⑦ carriers | independent scanner, independent final reviewer (models per the ⑩ dispatch table) |
| Probe quota | 0 |

## 5 Review conclusion summary (⑤⑥⑦)

| Stage | Conclusion line | Judgement |
| --- | --- | --- |
| ⑤ pre-review r1 | `PRE-REVIEW: PASS` | 2 Lows + several Infos; both Lows were fixed in ②'s targeted round |
| ⑥ independent scan | `INDEPENDENT SCAN: PASS` | Sections A–E complete; no High/Medium findings in Section D; the verbatim transcript is the fence below |
| ⑦ final review | `FINAL REVIEW: PASS` | all eight dimensions pass, 0 findings |

**The three-piece set**: the report pair `docs/validation/directive-tmp-self-heal{,_EN}.md`; the freedom list `docs/validation/evidence/DIRECTIVE-TMP-SELF-HEAL-freedom-list.md`; the gate readings in §8.

**⑥ independent scan receipt (verbatim transcript, not rewritten)**:

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

No findings.

## Section C

Q1: Is the v1.55 → v1.56 metadata update complete and internally consistent?  
PASS. The header and version bump are aligned in the CN and EN directive, and the changelog entry is present at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L102) and [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L102). The §0.3 version line is also updated in the same change surface.

Q2: Does the new same-command create/delete exception remain narrowly scoped, and are the non-self-heal and cumulative-cap closures present?  
PASS. The new clause is limited to the §5.1 self-heal tier and explicitly states the three required conditions: position valid, same-command create+delete, and zero residual evidence after the fact; it also rejects out-of-scope placements such as tmp root, repo root, or project subdirectories. The minimal closures are present in the same section at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L516-L520).

Q3: Is the §9.3 "提前清理" disposition cell extended without widening the policy beyond the allowed boundary?  
PASS. The cell now references both the self-heal exception and the new same-command exception, and the boundary expansion is explicitly recorded in the slice scope. See [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L902).

Q4: Are the appendix and OB tracking consistent with the slice's actual scope and deferral status?  
PASS. Appendix C.0aq is present and framed as a summary trace rather than a new source of law; OB-128 is recorded as the 25th deferral and tied to the governance-slice nature of this change at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1148) and [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L2080-L2084).

Q5: Do the P34 ledger rows and boundary lines in both CN and EN correctly reflect the slice status and boundaries?  
PASS. The P34 rows show the slice as in progress and the boundary notes are aligned in both language mirrors: [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L66), [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L1139-L1157), [docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L66), and [docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L1126-L1144).

The recheck across these exact edit points finds no semantic drift, no missing closure, and no boundary leak beyond the declared slice.

## Section D

UNCOVERED: none.

The guardrails that matter here are already spelled out in the text itself: the same-command exception is explicitly bounded, the non-self-heal path remains enforced, and the cumulative D cap still applies. If those protections were removed, the likely failure would be human misinterpretation rather than a mechanically visible contradiction in the directive text; this is exactly why the new section and the closure lines were added in the same patch.

## Section E

I cannot verify downstream adoption beyond the repository text itself, but the current evidence supports the PASS conclusion. What would overturn this PASS would be:
- a broader same-command exception than the tmp/<slice-id> bounded path,
- a bypass of the non-self-heal fail-stop rule,
- or a silent override of the cumulative D limit.

The current wording explicitly forbids all three, and the corresponding ledger and appendix entries remain consistent with that rule set.

INDEPENDENT SCAN: PASS
```

⑥ independent scan: conclusion `INDEPENDENT SCAN: PASS`, Sections A–E complete, no High/Medium findings in Section D; the source text is the fence above. This fence is ⑥'s re-sent version (in the first-round version the five section headings lacked the `##` prefix, which made G7-f FAIL; it was re-sent under the §7.2.2 ⑥ internal quota and is not counted into §9.2 R).

### ⑧a self-check

- ① Cost-table mechanical audit: yes — the nine entries in §10 sum to `12` and correspond one by one to the call rows of the §2 table.
- ② "Disposed" `file:line`: yes — `OB-128`'s "Ruled" entry is at `docs/DELIVERY_DIRECTIVE.md:1148`, mirrored at the same line.
- ③ Old-wording grep: yes — readings in the last table of §8 (two search keys, one group each for CN and EN).
- ④ ⑥ and ⑦ Info / Low consolidation: yes — ⑤'s two Lows in §1 item 1, ⑦'s 0 findings in §5, ⑥'s verbatim transcript in the §5 fence.

### A-NN numbering completeness

This report uses no `A-NN` numbers; exceptions and gaps are located by §1's numbered items and by file paths, so no entity definitions are required.

## 6 Review input package summary (including blind-review isolation proof)

| Stage | Input scope (per §7.2.2) | Isolation requirement |
| --- | --- | --- |
| ⑤ | start-up message, ① design document, implementation diff, gate readings | ruling records readable |
| ⑥ | slice fields, implementation diff, contract sentences, read-only constraint | blind package, ruling records and ⑤'s conclusions removed |
| ⑦ | ①–⑥ full set, deviation and waiver register, gate readings | full input |

**Blind-review isolation proof**: ⑥'s input package excludes ⑤'s conclusions and the ruling records per §7.2.2; this ⑧a input did not include the item-by-item package list, so only the basis is registered here and no in-package paths are recorded (gap in §11 item 5).

## 7 Falsifiability (mechanical checks + mutation testing)

| Layer | Mechanism | Covered | Not covered |
| --- | --- | --- | --- |
| L1 | pre-⑩ terminal read-back and §5.1's three readings | a missing exception element ⇒ red | whether the exception is misused has no mechanical detection |
| L2 | `G3(a)` six shapes and `G3(b)` key fields | a one-sided edit ⇒ red | whether the two sides are semantically equivalent |
| L3 | `G6-5` / `G6-7` section references and ledger reconciliation | a dangling reference ⇒ red | whether the v1.56 exclusion condition is too strict |

In the mutation table `M1` and `M2` are reasoned (editing the directive body is out of this slice's bounds) and `M3` is measured in memory (files are read only, nothing is written):

| No. | Mutant | Criterion that goes red | Criterion that does not |
| --- | --- | --- | --- |
| `M1` | delete the sentence "lawful location but ordering violation only does not use this exception" | none ⇒ only ⑤ / ⑦ can find it | every gate |
| `M2` | delete the "review-layer confirmation" sentence | none ⇒ `G6-5` judges only §X.Y references | every gate |
| `M3` | delete one `**` on the CN side only | `G3(a)` six shapes report a `bold` mismatch | `G3(b)` key fields |

## 8 Gate results

| Scope | Reading |
| --- | --- |
| `--scope=index` | `TOTAL_FAIL=0`; `G6-1 scanned 86`; `G3(a)` directive `+20/-4` symmetric with all six shapes `true`; `G3(b)` hash `1/1`; `G4b newCjk=453 unknown=[]` |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `--selftest` | `SELFTEST: PASS 697/697` |

The readings above are a snapshot taken at ⑥'s scan time; after this ⑧a added the reports, `G4b`'s `newCjk` grows with the new lines (a time-point difference, see §3), while the `TOTAL_FAIL=0` conclusion is unchanged.

**Old-wording grep readings** (evidence for ⑧a self-check ③):

| Side | Search key | Hits |
| --- | --- | --- |
| CN | `自愈档（v1.45）或同命令建删特例（v1.56）` | 3 |
| EN | `self-heal tier (v1.45) or same-command create-and-delete exception (v1.56)` | 3 |
| CN | `自愈档（v1.45 同步）` | 0 |
| EN | `v1.45 sync` | 0 |

## 9 ⑧c closing-cleanup record

**⑧c cleanup registry (record before cleanup, per the §5.1 ordering of cleanup and trace)**:

- Pending-cleanup list (item by item, 5 items, all inside this slice's temp directory): `tmp/DIRECTIVE-TMP-SELF-HEAL/apply.js`, `tmp/DIRECTIVE-TMP-SELF-HEAL/verify.js`, `tmp/DIRECTIVE-TMP-SELF-HEAL/step0-insert.js`, `tmp/DIRECTIVE-TMP-SELF-HEAL/DELIVERY_DIRECTIVE.md.bak`, `tmp/DIRECTIVE-TMP-SELF-HEAL/DELIVERY_DIRECTIVE_EN.md.bak`. The two `.bak` files are pre-edit snapshots and **are temp artifacts as well**, so they are listed too.
- Cleanup method: `Remove-Item -LiteralPath` item by item, no wildcards; nothing outside the list is touched (including the `tmp/gates-cache/g4b/` runtime cache and `tmp/.gitkeep`).
- Post-cleanup readings: see the "post-cleanup read-back" block at the end of this section.

**Post-cleanup read-back (five readings, all measured at ⑧c)**: ① `git status --porcelain -uall` shows `??` = 0 and a non-empty second column = 0 (`EXIT=0`); ② recursive file count under `tmp` (excluding `gates-cache/g4b`) = 1 (`EXIT=0`); ③ top-level `tmp` entries excluding `gates-cache` = `.gitkeep` (1 item, `EXIT=0`); ④ `Test-Path tmp/DIRECTIVE-TMP-SELF-HEAL` = `False` (`EXIT=0`); ⑤ recursive directory count under `gates-cache` = 1 (`g4b`, `EXIT=0`). Registry of 5 items → each `deleted` → `files_left=0` → directory removed.

Current state (at ⑧a): `tmp/DIRECTIVE-TMP-SELF-HEAL/` holds 5 items (`apply.js`, `DELIVERY_DIRECTIVE.md.bak`, `DELIVERY_DIRECTIVE_EN.md.bak`, `step0-insert.js`, `verify.js`); `tmp/gates-cache/` is the runtime cache (an exception-set member); `tmp/` also has `.gitkeep` at its top level.

## 10 Cost and metering

R7.1 per-stage call counts:

| Stage | Calls |
| --- | --- |
| ① design | 1 |
| Batch assistance | 3 |
| ② implementation (first landing + rework + targeted fix) | 3 |
| ③ test | 1 |
| ⑤ pre-review | 1 |
| ⑥ independent scan | 1 |
| ⑦ final review | 1 |
| ⑧a report | 1 |
| Probes | 0 |

## 11 Explicitly not executed

1. commit / push / CI: none executed; all require same-round authorisation after the ⑩ stop point. **Completed after the ⑩ stop point**: commit `acf676e`; CI run `40695233778`, both the ubuntu and windows jobs succeeded (per-leg `SELFTEST: PASS 697/697`, `G6-1 scanned 652`, `TOTAL_FAIL=0`).
2. Real probes: not executed (quota 0).
3. `GATES-EXEMPT` implementation: not executed (a referred item, left to a later slice).
4. Relaxing the v1.56 exclusion condition: not executed (left as a v23 observation candidate).
5. Item-by-item review input package list: not registered (the input had no list; only the basis is registered, see §6).
6. Per-round durations: not timed (§2 uniformly marks `未计时`).

**Open items (4, for the v23 queue)**: ① whether the v1.56 exclusion condition is too strict (should "lawful location + same-command create-and-delete" be admitted into the exception; trigger = this slice's first empirical D = 1); ② ⑥ format deviation is systemic across two slices (one each in P33 and P34, both contract-detail deviations); ③ the ⑧ evidence-write path is not bounded (this slice's D = 2 came from writing a log at the `gates-cache` root); ④ task-pack self-defined criteria lack prior review (this slice's ③ and ⑨ each produced 1 false positive).

## 12 Precedent-prohibition check

- **Two exceptions in this slice**: `R = 2` and `D = 1`; both register events beyond the normal caps.
- The `OB-37` meta-rule authorisation behind `D = 1` **must not be used as a precedent** (§1.7.4 ④ carries the same marker as §7.8.4's meta-rule ③).
- The **extension** (the §9.3 disposition cell) is a same-round user authorisation and is registered as a **fact**, not a precedent.
- The remaining reference (the "must not be used as a precedent" sentence in `C.0aq`) is only a cross-slice summary trace; this slice claims no waiver from it.

## 13 Next steps

1. ⑨ native verification → ⑧b write-back → ⑧c closing cleanup → ⑩ stop-point receipt; rectification after ⑧c must not rebuild `tmp/DIRECTIVE-TMP-SELF-HEAL/`.
2. ⑥'s source text has been transcribed verbatim into the §5 `sixstage` fence (re-sent version; the five section headings carry the `##` prefix and satisfy the G7-f contract); §11 has had its original item 5 deleted.
3. Evaluate the three v23 observation candidates (§1 item 9) together with the `GATES-EXEMPT` referred item.
4. Tier advice: keep governance slices of this kind at tier M; rounds that touch criterion semantics use a high-strength ⑦ carrier, while documentation-only rounds may drop a tier.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-tmp-self-heal.md	repo
docs/validation/directive-tmp-self-heal_EN.md	repo
docs/validation/evidence/DIRECTIVE-TMP-SELF-HEAL-freedom-list.md	repo
```

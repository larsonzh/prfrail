# T027 · DIRECTIVE-V1.11 directive-revision validation report (§1.7 execution freedom + §3.4 post-edit diagnostics + ⑧c wrap-up cleanup)

> Status: **complete (commit and push await same-round authorization)**. Slice type: **Medium** (documentation only; behavior change / multiple files / touches rule semantics); effort tier: `high`; ⑦ fixed Extra High.

## 1. Change summary (including the execution-freedom list)

| Item | Content |
|---|---|
| Changed files | `docs/DELIVERY_DIRECTIVE.md` and its strictly positional mirror `docs/DELIVERY_DIRECTIVE_EN.md` |
| New sections | §1.7 execution freedom and escalation boundaries; §3.4 post-edit diagnostics scan (hard step) |
| New stage | §5.1 pipeline gains **⑧c wrap-up cleanup** (after ⑧b, before ⑩) |
| Companions | the §5.2 matrix gains a ⑧c column; §5.3 gains a ⑧c row; §7.1 gains a ⑧c note; §11.3 gains a ⑧c record; Appendix B.7.4 gains the §1.7 and §3.4 section numbers |
| Ledger | the §0.3 changelog v1.11 entry plus header version v1.11; Appendix C gains **C.0g** |
| Not touched | the G1–G7 criteria themselves, the existing §5.1 step sequence, the §7.5 cap, the role set and carriers, code / CONTRACTS / schema / fixtures / CI |

### 1.1 Execution-freedom list (§1.7.1; re-checked by ⑦ final review)

| No. | Substantive choice | §1.7.1 basis | Scope of impact | §1.7.3 test answer |
|---|---|---|---|---|
| F-1 | The user's declaration referenced "§4.9", which **does not exist** in this directive (it exists only in the superseded v3.1 directive); the reference was repointed to the current anchors | completion of contract wording | §3.4.4 and the ⑧c note in §5.1 | if this decision were wrong, what changes is **wording**, not code or contract semantics ⇒ execution freedom |
| F-2 | Two 1.7.2 items normalised in wording ("超切片范围" to "超出切片范围"; "判定为" to "判定") | completion of contract wording | §1.7.2 | as above |
| F-3 | In the ⑧c note, "不只是变更集" was compressed to "不只变更集" | completion of contract wording | the §5.1 note | as above |
| F-4 | Appendix B.7.4's "All" row gained the §1.7 and §3.4 section numbers | completion of contract wording | Appendix B.7.4 | a **section-number index** caught up with the body; no rule substance |
| F-5 | Appendix C gained a new ledger block C.0g | completion of contract wording (plus the OB-12 ruling and the SW-3 / SW-4 / SW-5 precedent) | Appendix C | a ledger entry; contract semantics unchanged |
| F-6 | The landed net line delta differs from the ① plan's prediction (the plan predates the F-4 and F-5 decisions) | — | this report | disclosure only; no contract impact |
| F-7 | §3.4.1 keeps explicit command literals (duplicating the §6.1 table); ⑤ proposed reference-style instead | — | §3.4.1 | an implementation-expression choice; registered as a known trade-off |
| F-8 | §3.4 does not state whether it applies to the `[ROUTINE]` tier | — | §3.4 | disclosed as an open item; no silent change of meaning |
| F-9 | ② implementation landed by the master, with ⑧ not separately assigned | the user's declaration budget line lists only ①⑤⑥⑦ ⇒ this slice judged that ⑧ falls to the master | the ② stage | if this were wrong, what changes is the **executing role**, not code or contract semantics ⇒ execution freedom; however ⑤ round 2 pointed out that this choice **was not on this list and therefore was never ruled on by ⑦** ⇒ registered here and referred for a ruling |

> ⑦ round 1 ruled on F-1 … F-8 one by one: **all `in-bounds`** (all four questions answered no).

## 2. Execution pipeline

| Stage | Role and model | Result |
|---|---|---|
| ① upfront analysis | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` | a plan (per-anchor landing list, mirror strategy, G1–G7 risk list, 5 open items U1–U5) |
| ② implementation | master (⑧ not separately assigned in this slice) | both files landed; U1, U2 and U4 executed under the autonomous ruling and registered (see §1.1) |
| ④ integration and gates | master | `SELFTEST: PASS`; `--all --scope=tree` and `--scope=index` both `TOTAL_FAIL=0` |
| ⑤ pre-review | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` | `PASS WITH FIXES`: one High plus several Low and Info items |
| ⑥ independent scan | `independent-reviewer` / `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS`, 0 findings |
| ⑦ final review | `independent-reviewer` / `GPT-5.3-Codex (copilot)` | round 1 `FINDINGS` with one Medium; round 2 `INDEPENDENT SCAN: PASS`, 0 findings; round 3 (authorized this round by the user, narrowed to three questions) `PASS WITH FIXES`, no Medium+, all three questions ruled to hold |
| ⑧a and ⑨ and ⑧b and ⑧c | master | see §4 and §5 |

## 3. Review summary and per-item disposition

| Source | Finding | Disposition |
|---|---|---|
| ⑤ pre-review | High: the §0.3 v1.11 entry claimed Appendix C gains C.0g while nothing was on disk yet | Fixed: ⑤'s own option (a) was taken and **C.0g was landed immediately**, making the claim true |
| ⑤ pre-review | Low: two §1.7.2 items are not verbatim versus the user's text; one compression in the ⑧c note | Registered as F-2 and F-3 and kept (item count and meaning unchanged) |
| ⑤ pre-review | Low: §3.4.1 duplicates command literals from the §6.1 table | Accepted but kept (§3.4.1's purpose is to name executable commands); registered as F-7 |
| ⑤ pre-review | Low: "newly added warning" lacks a baseline definition, so that hard constraint is partly not mechanically checkable | Registered as an open item (to be handled together with the §3.4.3 aggregator) |
| ⑤ pre-review | Low: §3.4 does not state whether it applies to the `[ROUTINE]` tier | Registered as open item F-8 |
| ⑥ independent scan | no findings (`PASS`) | — |
| ⑦ final review | Medium: the §1.7 record-keeping sentence and the §5.1 order cannot both hold | Fixed: the sentence was rewritten so the re-check vehicle is explicitly the same list inside ⑦'s review input package, written back by ⑧b; ⑦ round 2 confirmed it **closed with no new defect** |
| ⑤ pre-review | re-review (round 2, re-run per §5.3): still `PASS WITH FIXES` - M-1 and M-2 and L-1 and L-2 and L-3 | Disposition in §5 and §6: M-1 has been evidenced by the whole-tree file-by-file encoding check; M-2 has been registered as F-9; L-1 and L-2 and L-3 had their wording corrected |
| ⑦ final review | round 3 (narrowed to three questions): (a) whether F-9 is `in-bounds`; (b) whether widening ⑧c's scope from the change set to the whole tree is acceptable; (c) whether re-judging the conclusion line to `FINDINGS` holds | All three were ruled **to hold**: (a) `in-bounds` (a factual registration, not a new rule entity); (b) **acceptable** (⑧c's baseline scope was always this; this round only corrected "the literal widening was not executed" into "it was really executed and disclosed as measured", the exemption boundary and the self-consistency check are on disk, and the residue is OB-27 awaiting formalisation); (c) **holds** (since 3 pre-existing violations were found by measurement, the conclusion line cannot stay `PASS`). **No Medium+** this round, so the red line was not triggered |

> **One item reported as measured**: after ⑤ reported its High, this slice **did not re-run ⑤** at that point (the budget cap for the pre-reviewer is ⑤×1, with one remediation re-review on top). Rationale: that High was closed by a **deterministic action** ⑤ itself prescribed (landing C.0g makes the claim true), the rest are Low and Info registrations, and **this round's fix changed no rule body (§1.6 and §1.7 and §3.4 and §5.x and §6.3)**. **The ⑤ re-review has since been run per §5.3** (see §2 and §3).

## 4. Gates and native validation

| Item | Result |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS` |
| `--all --scope=tree` and `--scope=index` | both `TOTAL_FAIL=0` (G1–G7 plus the base gates) |
| Encoding gate | every changed file is UTF-8 with BOM plus LF; 0 violations |
| Bilingual positional parity | line count, heading positions and levels, bold markers, pipes and blank lines all equal; 0 per-line mismatches |
| Link and reference resolvability | every newly referenced section number (§1.6 and §6.1 and §6.3 G4 and §5.1 and §5.2 and §5.3 and §7.1 and §11.3 and Appendix B.7.4 and Appendix C) exists in this directive; the §4.9 the user's declaration named does not exist and was repointed per F-1 |
| Degraded form for document-type slices | per the §5.2 note, ⑨ degenerates into the encoding gate plus the G1–G7 self-check plus link and reference resolvability |

## 5. ⑧c wrap-up cleanup record (first run)

| Field | Result |
|---|---|
| When and by whom | 2026-09-23; the master (0 paid calls) |
| Whole-repo diagnostics (per §3.4.1 categories) | compiled: `go build ./...` and `go vet ./...` (through the `gate.js` base gates); scripting: nothing changed this round; markup and data: the encoding gate for both files plus the `gate.js` parse; plain text and others: none; uncovered types: none |
| IDE diagnostics summary | see the uncovered items in §6 |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0`; `--scope=index` and `--scope=ci` are `TOTAL_FAIL=0` as well |
| Encoding gate (all types, file by file) | whole-tree `git ls-files` sweep, file by file: 2820 text files, **3 violations**, all inside frozen evidence packs (`b2-2026-09-17/` and `b3b-2026-09-21/`); handled under a **temporary exemption**, see §6 |
| Frozen evidence-pack self-consistency check (exemption condition two) | the `SHA256SUMS.txt` of b2 (41 entries), b3a (128) and b3b (2158) matches the files on disk **entry by entry, 0 missing, 0 mismatched**; `b1-20260917` has no `SHA256SUMS.txt` ⇒ **not self-checkable** (disclosed) |
| Temp directory, probe leftovers, untracked files | `git status --porcelain -uall` shows only this slice's 4 staged changes and no untracked residue; under `tmp/` this slice's material is `tmp/v111/**` and `tmp/b4check/**` (not tracked by the repository) |
| Conclusion line | `WRAP-UP CLEANUP: FINDINGS` (the all-types file-by-file encoding check found 3 pre-existing violations **inside frozen evidence packs**; fixing them would break the packs' integrity ⇒ an §1.7.2 escalation boundary, so it was stopped and registered, and the user's ruling was taken: handled under a temporary exemption plus an evidence-pack self-consistency check, see §6) |

**Three milestones**: ① **⑧c's first self-referential verification ran** - this is the first slice checked by ⑧c immediately after ⑧c was added, and its first run returned `FINDINGS` (see the table above and §6); ② **that first run caught 3 historical leftovers on the spot** - the whole-tree file-by-file sweep found 3 pre-existing violations, proving the hard item **takes effect in practice**; ③ **⑦ ruled every item F-1 to F-8 of the §1.7 execution-freedom list `in-bounds`** - the new record-keeping and boundary-judgement mechanism took effect on its first run (the list travelled inside ⑦'s input package and was ruled on).

## 6. Unexecuted items and uncovered items

| Item | Reason |
|---|---|
| The actual counts of the IDE diagnostics summary | this session has no Problems-panel readout; registered as uncovered, with no claim that it was checked |
| Re-running ⑤, round 3 | the user ruled **not authorized**: M-1's remediation exceeds ⑤'s review scope (deciding the scope is a rule-semantics call), and M-2 and the L items are ledger-metadata and prose class, exempt from re-review under **the OB-11 exception standard** |
| Whether §3.4 applies to the `[ROUTINE]` tier | open item F-8, for the §12 revision to assess |
| The 3 pre-existing encoding violations inside frozen evidence packs (b2's `probe10-dir-listing.raw.txt` is CRLF; b3b's two `pin.json` files carry a BOM) | The user ruled **option A**: ① these 3 files **do not violate G4a** (whose scope is explicitly the change set); what they violate is the §5.1 ⑧c wording "all types, file by file" ⇒ this is "⑧c widening met historical leftovers"; ② a **temporary exemption** (to be formalised by `DIRECTIVE-V1.12`) - the exemption is strictly limited to **archived packs that carry a MANIFEST/SHA256SUMS**; **a new evidence pack must still pass the encoding gate before it is committed**; ③ ⑧c checks evidence packs by **self-consistency** instead (see §5); ④ **OB-27** is registered, including the to-be-discussed item coordinating G4a and ⑧c and the B2 fidelity ruling (the CRLF of `.raw.txt` directly conflicts with `.gitattributes`'s `*.raw.txt -text`) |
| The baseline definition of "newly added warning" | open item, to be handled together with the §3.4.3 aggregator |

## 7. Cost and metering

Budget: ①×1 plus ⑤×1 plus ⑥×1 plus ⑦×2; realized: ①×1 plus ⑤×1 plus ⑥×1 plus ⑦×3 (one final review round, one re-review round, plus round 3 that the user **separately authorized** this round - narrowed to three questions; ⑤ round 3 was **not authorized**). Probe quota: 0.

## 8. Next-step recommendations

- ⑩ stop point: ⑦ round 3 has passed (all three questions hold, no Medium+), so this slice may be committed and pushed, and CI is then observed for a double green per §6.4; F-8, the "newly added warning" baseline and OB-27 all go to `DIRECTIVE-V1.12`.
- Next slice candidates (already listed by the user; not folded into this slice): `GATES-TIGHTEN` (OB-17 and OB-20); `DIRECTIVE-V1.12` (OB-10 and OB-12 and OB-16 and the §5.3 versus §7.5 conflict standard and OB-23).

## 9. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/directive-v1.11.md	repo
docs/validation/directive-v1.11_EN.md	repo
tmp/v111/**	local-only
tmp/b4check/**	local-only
```

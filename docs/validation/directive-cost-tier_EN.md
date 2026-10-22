# Slice `DIRECTIVE-COST-TIER` verification report (Directive v1.43)

> Slice: `DIRECTIVE-COST-TIER`; class `[SLICE]`, governance slice, cost tier M.
> Directive: v1.43 (MINOR, rule-revision type; this is an independent governance slice, not a by-product of an execution slice).
> Baseline: `d:\LZProjects\prfrail`; written 2026-10-22 by the ⑧a documentation role.
> Bilingual: this file is the CN authority; `docs/validation/directive-cost-tier_EN.md` is the line-for-line mirror.
> Boundary: this slice edits `.md` only; it touches no `internal/**`, `tools/gates/**`, `schemas/**`, `testdata/**` or `.github/workflows/**`.

## 1 Change summary

This slice closes OB-31 (scattered R-threshold statements and the unmerged fallback count) and OB-68 (§5.2 three-model tiering out of sync with the effort tiers of §7.5 / §8.2) in one landing, and hooks OB-84's cost caliber to it.

| File | numstat | Content |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | `+55/−18` | Three R thresholds unified, six ⑦ labels rewritten, §8.2 mapping table, §7.5 cost paragraph, two rows in §12.4, one row in §12.5 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | `+55/−18` | Line-for-line mirror of the row above |
| `docs/DEV_PLAN.md` | `+8/−0` | T027 section backfills P18 / P19 / P20 plus a backfill note |
| `docs/DEV_PLAN_EN.md` | `+8/−0` | Line-for-line mirror of the row above |
| `docs/t027/REMAINING_SLICES.md` | `+20/−0` | P20 status-table row (`✅`), P21 status row and the `### P21` block (`🔄`) |
| `docs/t027/REMAINING_SLICES_EN.md` | `+20/−0` | Line-for-line mirror of the row above |
| `docs/validation/directive-cost-tier.md` | new | This report (CN authority) |
| `docs/validation/directive-cost-tier_EN.md` | new | The `_EN` mirror of this report |
| `docs/validation/evidence/DIRECTIVE-COST-TIER-freedom-list.md` | new | Execution freedom list |

The final change set is 11 `.md` paths (8 rewritten + 3 new, including the expiry-date extension entry in `ADR_REGISTER{,_EN}.md`).

```artifacts
docs/DELIVERY_DIRECTIVE.md
docs/DELIVERY_DIRECTIVE_EN.md
docs/DEV_PLAN.md
docs/DEV_PLAN_EN.md
docs/ADR_REGISTER.md
docs/ADR_REGISTER_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
docs/validation/directive-cost-tier.md
docs/validation/directive-cost-tier_EN.md
docs/validation/evidence/DIRECTIVE-COST-TIER-freedom-list.md
```

## 2 Execution pipeline

Step 0 clears the ledger first; then ① designs, ② lands and revises three times, ③ checks twice, ⑤ pre-reviews twice, ⑥ scans independently once, and ⑦ final-reviews twice; ⑧a lands this report and the freedom list.

| # | Stage | `agentName` | `model` | Conclusion |
|---|---|---|---|---|
| 1 | ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` |
| 2 | ② landing | `prfrail-documenter` | same model as the orchestrator | `IMPL: DONE` |
| 3 | bulk assist A | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | `EXTRACT: DONE` |
| 4 | ② revision (cost table and pointerization) | `prfrail-documenter` | same model as the orchestrator | `IMPL-REV: DONE` |
| 5 | bulk assist B | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | `MISMATCH` (③ adjudicated it as partly accepted; the conclusion does not hold) |
| 6 | ③ consistency check r1 | `prfrail-documenter` | same model as the orchestrator | `CHECK: FINDINGS` (High=0 Medium=1 Low=2 Info=4) |
| 7 | ② revision r3 | `prfrail-documenter` | same model as the orchestrator | fixes F1 / F2 / F3 |
| 8 | ③ re-check r2 | `prfrail-documenter` | same model as the orchestrator | `RE-CHECK: PASS` |
| 9 | ⑤ pre-review r1 | `independent-reviewer` | `GLM-5.3 (glm)` | `PRE-REVIEW: PASS WITH FIXES` (1 Low + 3 Info) |
| 10 | ② revision r4 | `prfrail-documenter` | same model as the orchestrator | fixes ⑤ Low#1 (the §7.5 unit-distinction sentence) |
| 11 | ⑤ re-review r2 | `independent-reviewer` | `GLM-5.3 (glm)` | `RE-REVIEW: PASS` |
| 12 | ⑥ independent scan | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS` (0 findings; no Section E) |
| 13 | ⑦ final review r1 (blind) | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `ADJUDICATION: FINDINGS` (High=0 Medium=1 Low=2 Info=0) |
| 14 | ② revision r5 | `prfrail-documenter` | same model as the orchestrator | fixes the three ⑦ r1 items (M-1 / L-1 / L-2) |
| 15 | ⑦ re-review r2 (narrowed) | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `RE-REVIEW: PASS` |

**③ ordering evidence**: the ③ r1 offset is `2026-10-22T10:10:23Z`, earlier than the ④ gate timing run (around `10:30`) and ⑤ r1 (`10:46:46Z`); ③ sits after ② and before ⑤, matching the §5.3 ring order.

**Position of ⑧a**: this report and the freedom list are landed by ⑧a; the ledger write-back, ⑨ evidence registration and ⑩ closeout from ⑧b onward are run by later steps, and this report does not pre-empt their conclusions.

## 3 Deviations and fallbacks

Each ring, root cause and merge decision is logged per the §9.2 R-count definition.

- **D-1 (the only deviation)**: ② revision r5 wrote one character outside the whitelist (code point `U+7686`), which tripped `G4b` (the 37-entry whitelist does not carry it) ⇒ rewritten as 「两者**均**无」 with the meaning unchanged; the whitelist was not widened.
- **D-2 (corrupted evidence artifact)**: the `02-changes.diff` handed to ⑤ during ⑨ was captured through PowerShell redirection and corrupted (CN line breaks swallowed, numstat inconsistent with the files), surfaced by a note in ⑤ re-review r2 ⇒ re-produced with a direct `node execFileSync` write and re-read as `FFFD=0`. Logged as one more instance of the same family as §9.3 "PowerShell pipeline capture of native stdout ⇒ garbled evidence file"; no new OB is opened. The truth source for the cumulative index is the §9.3 failure-mode library, so this report does not assign a number. **Disposition**: that artifact is **not used as review evidence** - this slice's review relies on `03-numstat.txt` and `04-status.txt`; ⑧c removes the file normally, with no rework.
- **D-3 (three ③ r1 findings)**: F1 (Medium) - "Appendix B.7" is in fact **Appendix B.0**, corrected (2 places in the §0.3 row and 2 in the OB-31 handling column, 4 across CN / EN); F2 (missing deliverable in C.0ad) and F3 (qualifier added to §5.2) were closed in the same round.
- **D-4 (mechanical integration)**: after two rounds of document edits the orchestrator mechanically refreshed the `file:line` anchors in the OB-31 handling column (`:500/:527/:820/:920/:1151` → `:501/:528/:841/:928/:1159`); this is mechanical integration and changes no semantics.
- **D-5 (pending a user ruling)**: whether "permanent pilot exemption for governance slices" joins pool E - not done in this slice, left to the user.
- **D-6 (not registered)**: this slice registers no "pilot deferral" OB; the basis is the user's 2026-10-22 ruling that a `[SLICE]` start does not register one.
- **R count**: D-1 to D-3 count 1 each (different rings, different root causes, none merged); D-4 is mechanical integration and does not count; D-5 and D-6 are ruling and registration matters and do not count; this slice has R = 3. Per §9.2, R≥2 ⇒ the §5.2 high-risk row applies (⑤ re-check is mandatory), and R does not exceed 3 ⇒ no pause is triggered. **This is also the first bootstrapping check of the new R definition**: in a real slice R≥2 already drove the high-risk row (the ⑤ re-check became mandatory), and the ⑦ quota ran as 1+3, consumed 2 and stayed within the cap.

## 4 Environment and effort tiers

| Item | Value |
|---|---|
| Workspace | `d:\LZProjects\prfrail` (Windows) |
| Size and type | Tier M, governance slice, `[SLICE]` |
| Orchestrator tier | High (the "medium" row of §8.2) |
| ⑦ tier | `gpt-5.3-codex` = **Extra High** (not downgraded, OB-37 not triggered) |
| Other two rows of the ⑦ mapping | `gpt-6.1-sol` = **Extra High** (may use **Max**); `gpt-6-luna` = **Max** |
| ⑦ quota | 1 final + 3 re-reviews (the "governance slice" row of §7.5; when several rows match take the most re-reviews), 2 used |
| Probe quota | 0 |

## 5 Review verdict summary

| Stage | Round | Verdict | Handling |
|---|---|---|---|
| ③ consistency | r1 | `CHECK: FINDINGS` (High=0 Medium=1 Low=2 Info=4) | ② revision r3 fixes F1 / F2 / F3 |
| ③ consistency | r2 | `RE-CHECK: PASS` | converged |
| ⑤ pre-review | r1 | `PRE-REVIEW: PASS WITH FIXES` (1 Low + 3 Info) | ② revision r4 fixes Low#1 |
| ⑤ pre-review | r2 | `RE-REVIEW: PASS` | converged |
| ⑥ independent scan | r1 | `INDEPENDENT SCAN: PASS` (0 findings) | no rework |
| ⑦ final review | r1 (blind) | `ADJUDICATION: FINDINGS` (High=0 Medium=1 Low=2 Info=0) | ② revision r5 fixes the three items |
| ⑦ final review | r2 (narrowed) | `RE-REVIEW: PASS` | converged |

## 6 Review input package summary

- The file list and byte counts of each review input package are logged by ⑧b at write-back time; this report records only how each package is composed.
- ⑦ r1 was dispatched as a blind review: its input carries none of the ③ and ⑤ verdict lines; the per-file readings of the isolation list were not available, so this report lists no files and records only the isolation method.
- The freedom list of this slice is landed by ⑧a (the user's dispatch caliber) and therefore postdates ⑦ r2; the ⑦ input package consequently **does not contain** the list, §1.7.3's "⑦ re-checks the list" is not covered, and the gap is logged in §8.

## 7 Self-evidence on criteria and figures

| Change | User-given value | On-disk site (CN side, as measured) | Evidence |
|---|---|---|---|
| (a) Three R thresholds unified | `R≥2` high risk, `R>3` pause | §5.2 high-risk row `:501`; §5.3's ⑨ row `:528`; §9.2 clause 3 `:842` | The three statements agree and share the caliber of the R-count definition paragraph at §9.2 `:838` |
| (b) ⑦ labels and mapping | 「⑦（§5.2 分档表选定）」 | §5.1 / §7.1 / three §7.5 rows / six §8.2 places; §8.2 mapping table `:806-812` | Six labels unified; the three model tiers match the user-given values cell by cell |
| (c) §7.5 cost paragraph | `175/1400`, `200/1000`, `20/75` (Credits / 1M tokens) | §7.5 cost table `:719-726` | Three rows and six cells match the user-given values cell by cell; the unit-distinction sentence was added by ⑤ Low#1 |
| (d) §9.2 clause 3 | Shares the same R count with §5.2 / §5.3 | `:842` | States in words "share the same R count; no second parallel ladder" |

- **Anchor re-check (⑧b)**: the `file:line` anchors of the OB-31 handling column have been mechanically refreshed to `:501` / `:528` / `:842` / `:929` / `:1160`, matching the current disk; the refresh is mechanical integration (no semantic change) and is not counted in R.
- **Two pointer conversions**: the pilot-metric row of §12.1 and the fallback-record row of Appendix B.0 now point at the authority, anchors `:928` / `:1159`; the numeric authority for both stays in §9.2.
- **Version consistency**: the header v1.43 and the last line of §0.3 v1.43 agree (`G1-b` measured green).

## 8 Known limitations

- **⑥ coverage gap**: ⑥ reported `INDEPENDENT SCAN: PASS` but produced no Section E; the gap is logged as it stands and is not read as full-section coverage.
- **Two units must not be converted**: the Credits caliber of §7.5 differs from the USD caliber of the "cost-gradient hint" in the same section; the two must not be converted into each other or merged.
- **§9.1 is not folded into R**: the ring-level clauses of §9.1 (retry and pause criteria) stay as they are; this slice unifies only the R-count definition of §9.2 and the three threshold actions, and folds no ring-level clause into R.
- **Ordering gap of the freedom list**: §1.7.3 requires the list to land after ⑥ and before ⑦ and to be checked by ⑦; this slice's list is landed by ⑧a and postdates ⑦ r2 ⇒ that re-check duty is not covered, and this entry is a record as it stands, not proof that ⑦ checked it.
- **`G4a` judgment gap**: `G4a` verifies only "a BOM exists", not "how many BOMs exist", so a double BOM passes; logged as OB-112 (`观察中`), and handling requires a separate slice.
- **Evidence caliber of D-2**: the corrupted `02-changes.diff` does not count; this slice takes the artifact written directly by `node execFileSync` as evidence, `FFFD=0`.
- **D-5 undecided**: whether "permanent pilot exemption for governance slices" joins pool E is not decided in this slice.
- **No dedicated mechanical criterion**: this slice is a wording-and-figure rule revision, so beyond the general G1-G8 criteria it has no dedicated criterion; `G6-4` and `G6-5` have no hits on this slice's added lines.

## 9 Gates and performance baseline

- **Three scopes**: `--all --scope=index|tree|ci` all report `TOTAL_FAIL=0`.
- **Self-test**: `--selftest` = `SELFTEST: PASS 588/588`.
- **Anomalous characters**: `G4b` reports `unknown=[]`, and the 37-entry whitelist was not widened.
- **Expiring literals**: `G6-4` has no hits.
- **Version**: `G1-b` reports header=v1.43 and last line of §0.3=v1.43.
- **CI (main commit `e5c9990`, run `39786249399`)**: both legs finished `success`; package-level `ok` lines: windows 15 / ubuntu 17; in-leg `SELFTEST` 588/588 and `TOTAL_FAIL` 0 on both legs. The first run was green, so the §6.4 evidence-based rerun was not triggered and **no DR entry was registered**.

| Observation | Reading |
|---|---|
| `--all --scope=index` cold start | 116.4 s |
| `--all --scope=index` warm cache | 78.8 s |
| `gofmt -l .` | 0.43 s |
| `go build ./...` | 4.72 s |
| `go vet ./...` | 3.50 s |
| `go test ./...` | 11.35 s |
| base four combined | 20.0 s |
| criteria side, isolated timings combined | 75.4 s |
| of which `G4b` | 52.4 s (about 69%) |
| remaining 17 items | 0.56-2.20 s |

**Caliber note**: each isolated timing includes 18 process launches and cannot be added directly to a single `--all` run; the figures serve only as a performance baseline and change no gate logic.

## 10 Cost and measurement

Per hard rule R7.1 of §7.5, every entry is taken from the session call receipts.

| Role | Calls | Detail |
|---|---|---|
| ① architecture | 1 | `deep-reasoner` |
| ② implementation / documentation | 6 | includes 3 revisions (r3 / r4 / r5) |
| ③ consistency | 2 | r1 and r2 |
| ⑤ pre-review | 2 | r1 and r2 |
| ⑥ independent scan | 1 | `MAI-Code-1.1-Flash (copilot)` |
| ⑦ final review | 2 | r1 blind and r2 narrowed; the 1+3 quota is not exhausted |
| bulk assist `GLM-5.3-FlashX (glm)` | 2 | extraction and reconciliation rounds |
| bulk assist `GLM-5.3-Flash (glm)` | 0 | this slice's bulk assist used FlashX only and never reached Flash |

- **Why Flash is 0**: the user asked on 2026-10-22 that this slice burn FlashX first; both bulk-assist rounds were carried by FlashX.
- **Expiry extension**: the expiry dates of `GLM-5.3-FlashX (glm)` and `GLM-5.3-Flash (glm)` move from 2026-10-31 to **2026-11-06** under an additional user authorization; this report records the fact only, and ⑧b must mirror it into the ledger at write-back.

## 11 Items explicitly not executed

- No commit, no push; gitee untouched.
- No code change, no change to the `tools/gates/**` criteria themselves, no change to schemas or fixtures.
- No change to the rule bodies of §2.2 / §2.4 / §6.3; no change to the quota figures of §7.5.
- OB-112 not handled; no "pilot deferral" OB registered; D-5 not decided.
- No Linux VM re-run: this slice has no platform-related change.

## 12 Next steps

- ⑧b ledger write-back: flip the P21 status to `✅` and record the FlashX / Flash expiry extension.
- ⑨ evidence registration: take the artifact written directly by `node execFileSync` as the basis and attach the `FFFD=0` reading.
- The next `GATES-*` batch handles OB-112 (a consecutive-BOM-count sub-assertion for `G4a`).
- The D-5 "permanent pilot exemption for governance slices" is written only after a user ruling.

## 13 Degrees of freedom

The degrees of freedom exercised by this slice, the boundaries fixed by rulings and the single deviation are logged one by one in `docs/validation/evidence/DIRECTIVE-COST-TIER-freedom-list.md` (F1-F7, M1-M3, D1).

## 14 Not-yet-synchronized items

| Item | State | Owning step |
|---|---|---|
| P21 status in `docs/t027/REMAINING_SLICES{,_EN}.md` | ✅ Flipped | ⑧b |
| OB-31 `file:line` anchors | ✅ Closed (`:501` / `:528` / `:842` / `:929` / `:1160` match the disk; see §7) | ⑧b |
| T027 status field in `DEV_PLAN{,_EN}.md` | still `BLOCKED / NOT IMPLEMENTED` | backfilled by a later slice per §11.2 |
| Expiry dates of FlashX / Flash | not in the ledger | ⑧b write-back |
| Section E of ⑥ | not produced | coverage gap, logged in §8 |
| Per-file byte counts of the review input packages | not listed | ⑧b logs them |
| Pool E ownership of D-5 | undecided | user ruling |

DOCS: DONE

# LEDGER-BACKFILL — ledger and document back-fill (validation report)

> Slice `LEDGER-BACKFILL` (documentation back-fill class, not directive governance) | class `[SLICE]` | size S | date 2026-09-29.
> **Status: stopped — incomplete** (the ⑩ stop point; the incomplete list is in §12; commit and push need same-turn authorization).

[中文](ledger-backfill.md)

## 1. Change summary (including the execution-freedom list)

- **What this slice did**: ① verified the B3a paragraph in `DEV_PLAN{,_EN}` (**verdict: it already exists, so no back-fill is needed**); ② added a fifth item to the OB-53 evidence cell; ③ applied the six wording corrections and notes **★1 to ★6**.
- **Special declaration for ★1**: this slice's handling of ★1 **fulfils the user-authorised exit reserved in freedom item F-2 of slice `S1-REACH`** (F-2's own wording: "if the user considers that a correction note should be added inside the row, it can be added separately after a same-turn authorization"; that slice's strict choice under §1.7.2 **is not overturned**); under the **user's same-turn authorization** this slice changes only the **evidence cell**, keeps the **disposition word unchanged** (`待议`) and **does not close that OB**.
- **What this slice did not do**: no change to any run id, test name, conclusion line, task-status field or OB opening status word; no OB was closed; no criterion changed; `internal/`, CONTRACTS, schema, fixtures, workflows, dependencies and role files were not touched.
- **Artifacts**: this report, the freedom list, `DELIVERY_DIRECTIVE{,_EN}` (raised to v1.16), `ADR_REGISTER{,_EN}` (ADR-016, ADR-017), notes in `t027/REMAINING_SLICES{,_EN}` and `DEV_PLAN{,_EN}`, and the §9 wording note in the DR-1 report.
- **Execution-freedom list**: `docs/validation/evidence/LEDGER-BACKFILL-freedom-list.md` (F-1 to F-7).
- **One-line verdict**: this slice **changes no verdict** - it only corrects stale wording to the current facts; the T027 schedule and the separate queueing of `DIRECTIVE-LEDGER-ARCHIVE` follow the user's rulings unchanged.

## 2. Item-by-item landing points (★1 to ★6 and ①②)

| Item | Landing point (CN and `_EN` in the same position) | Disposition |
| --- | --- | --- |
| ★1 OB-21 evidence-cell correction note | the OB-21 row of §12.4 in `docs/DELIVERY_DIRECTIVE.md` | a note appended to the evidence cell; the disposition word and the conclusion are unchanged and **the OB is not closed** |
| ★2 B4 priority note | the B4 block of `docs/t027/REMAINING_SLICES.md` | note added: B4 completed on 2026-09-23, so the ordering is due; both are next-batch candidates |
| ★3 integrity-note correction | the ledger integrity note in the same file | rewritten to "back-filled (commit `dd36b9b`)" plus "the earlier wording is out of date" |
| ★4 DR-2 to DR-5 planning wording | the four rows of the known-defects table in the same file | the **original wording is kept** and a note "(fixed by DR-FIX on 2026-09-22)" is appended after it |
| ★5 the T018 gap list in DEV_PLAN | the T018 status line of `docs/DEV_PLAN.md` | a wording note added (T026 is complete, DR-1 is fixed, AT-23 is *not passed*); **the task-status fields are unchanged** |
| ★6 the count basis in DR-1 report §9 | the mirror-scope paragraph of `docs/validation/dr-1-availability-probe-floor.md` | note added: "HEAD basis; the worktree count moves with added lines, so a fresh re-run governs - it measured 38 / 258 at this slice's wrap-up" |
| ① B3a paragraph | `docs/DEV_PLAN.md:137` and `_EN.md:138` | **verified present (commit `dd36b9b`) ⇒ the back-fill is not repeated** |
| ② OB-53 fifth item | the OB-53 row of §12.4 in `docs/DELIVERY_DIRECTIVE.md` | note added: `model: "auto"` can reach a blacklisted model indirectly |

## 3. Deviations and fallback record

1. **Finding under ① (important)**: the ledger integrity note said the B3a paragraph was still owed, but that paragraph **had already been added by commit `dd36b9b`** (`DEV_PLAN.md:137`, `_EN.md:138`; traced with `git log -S`, output in `tmp/lbf/dd36b9b-trace.txt`, disposition `local-only`) ⇒ this slice writes no new paragraph and instead **corrects the note** (★3). This is an instance of a note going stale itself.
2. **Authorization layer for ★1**: see the special declaration in §1; the report therefore states plainly that it is **not an ordinary cleanup**.
3. **Self-detected defects**: ① after the first write the ADR pair had asymmetric changed-line counts (CN `+2` versus EN `+1`) ⇒ a blank line was added on the EN side; ② `G2` flagged the old wording quoted in the new ★3 sentence (it contained a banned placeholder shape) ⇒ reworded to "the earlier wording is out of date"; ③ `G6-4` did not flag these notes (line-count and file-count literals were avoided).
4. **Master fallbacks = 1** (see §3 item 3, sub-item ①).
5. **Master error (recorded truthfully)**: the first-round input pack's `tmp/lbf/review-input.diff` was an **empty file** — it was generated **after** `git add`, so `git diff` produced no output. This is the **master mis-supplying the input**, **not a ⑦ problem**; the pack was regenerated with `git diff --cached` (non-empty). To re-check the four dispositions this slice consumed a **⑦ quota exception beyond the 1/1 budget**, ruled by the user on 2026-09-29 and recorded independently as **ADR-017**.
6. **OB-54 candidate (not handled here)**: the tail of the B3b paragraph at `docs/DEV_PLAN.md:139` still says the B3a paragraph is outstanding (the evidence that it is present is item 1 above) ⇒ registered as an **OB-54 candidate** for the **next slice**; this slice only records it under "deviations and fallback" and does **not** write it into §12.4.

## 4. Execution pipeline

| Stage | Role | Status | Artifact |
| --- | --- | --- | --- |
| ① architecture | not enabled | the §5.2 documentation-form row is marked — | — |
| ②b implementation | master | executed | the corrections and notes in §2 |
| ③ tests | master (documentation consistency check) | executed | gate and mirror readings |
| ④ master integration and gates | master | executed | see §9 |
| ⑤ pre-review | V4 Pro | see §6 | see §6 |
| ⑥ independent scan | to be skipped | see §6 | confirmed by ⑦ |
| ⑦ final review and re-review | V4 Pro (downgrade, ADR-016; the re-review is a quota exception, ADR-017) | executed | see §6 |
| ⑧a document draft | master | executed | this report |
| ⑨ native validation | master | executed | gate re-run and mirror verification |
| ⑧b document finalisation | master | executed | ledger write-back and closeout line |
| ⑧c wrap-up cleanup | master | executed | see §10 |
| ⑩ stop point | master | stopped | this report |

## 5. Environment and tier

- Platform: the Windows host; class `[SLICE]` (documentation form, degraded per §5.2); master tier is the standard one.
- **⑦ downgrade (ADR-016)**: `deepseek-v4-pro` (not on the §2.2 allowlist), a direct user ruling with the four elements and "not usable as a precedent"; this report states plainly that **independence is reduced** and hands the conclusion to ⑧c.
- **⑦ re-review quota exception (ADR-017)**: after the 1/1 budget was exhausted, the user ruling of 2026-09-29 added a second review round (reason = the first-round Medium was a master input-pack error that needs independent re-checking); it carries the four elements and "not usable as a precedent", and the report records it truthfully in §3 item 5 and §9.
- **Fifth pilot deferral (four elements)**: ruling party = the user; date = 2026-09-29; reason = a pure documentation back-fill, not pilotable; basis = the ruling text of this turn. **Must not be used as a precedent**; this slice makes no pilot metric statement.
- **Probes = 0** (no real external call was made).

## 6. Review conclusion summary (⑤⑥⑦)

- ⑤ pre-review (`deepseek-v4-pro`, budget 1/1): `PRE-REVIEW: PASS WITH FIXES` — 2 Medium (EN pointer misdirected; the gate reading baseline could not be carried over) and 5 Low (file-name inconsistency, ★1 wording better framed as "fulfilling F-2's reserved exit", the ⑥ three-condition paraphrase missing `.yml`, the `DEV_PLAN.md:139` observation, the `tmp/` inventory). The pointer, the wording and the `.yml` item were fixed literally, the ★6 count was corrected to the freshly re-measured value, and the gate-baseline item is handled by re-running the gates after landing (see §9).
- ⑥ independent scan: judged against the three §7.2.1 conditions — ① the change touches only `.md` and ledger-class files; ② it contains no semantics change in `.go`, `.json`, `.yml`, schema, fixtures or script logic; ③ it involves no state-machine, gate-criterion, evidence-model, role or authorization-contract semantics, so **all three hold and it may be skipped**; **⑦ confirmed** it (⑦ ruled it skippable).
- ⑦ final review (`deepseek-v4-pro`, downgraded, ADR-016, budget 1/1): `FINAL REVIEW: PASS WITH FIXES` — the four areas read "basically holds / passes / truthfully covered / passes"; 1 Medium (the review input pack's diff was an empty file because it was generated after staging) and 3 Low (the ADR-016 row sat outside the table, the disposition list carried an inconsistent count, the trace evidence was missing from the pack). All four were disposed of: the diff was regenerated with `git diff --cached`, ADR-016 is back inside the table in both languages, the count was corrected to the freshly re-measured value, and the trace output was landed and added to the pack.
- **Budget boundary and user ruling for the ⑦ re-review (stated truthfully)**: §5.3 requires ⑦ to be re-run after remediation, but this slice's ⑦ budget is **exhausted at 1/1** ⇒ per §7.5 the matter **stops for a user ruling**; the master must not declare a pass on its own. The user ruling of 2026-09-29 **authorised the second round** (a quota exception; reason and four elements in **ADR-017**).
- **⑦ round-two re-review (`deepseek-v4-pro`, a quota exception, ADR-017)**: `FINAL REVIEW: PASS WITH FIXES` - **all four first-round dispositions are closed** (the empty diff was rebuilt non-empty and is byte-identical to the staged set; ADR-016 and ADR-017 both sit inside the table with the added rows in the same position; the counts match the fresh readings; the trace output is in the pack); ADR-017 carries the four elements and the "not a precedent" mark; the ★1 special declaration and "no boundary crossed" were confirmed item by item; and the gate and self-test readings all match §9. Two Low findings: **F2R-1** (the forward reference at the end of §6 and in §12 dangled, and the ledger closeout line still said the re-review was not run) and **F2R-2** (the in-pack disposition list kept an intermediate `+1/-0` reading). Both are disposed of: the verdict line is this bullet, §12 is now self-contained, the ledger closeout line is synchronised and the in-pack reading is corrected to the fresh one; under the semantic boundary of §5.3 the **fact-record** writes of ⑧b (review verdicts, ledger write-back, reading corrections) do not trigger a further re-review - matching the same ruling in §10 of the `S1-REACH` report.

## 7. Review input-pack summary (including blind-isolation proof)

- Input pack: `tmp/lbf/review-input.diff` (the raw diff, regenerated with `git diff --cached`), `tmp/lbf/report-cn.md` and `report-en.md` (this report in full), `tmp/lbf/gate-tree.txt` (the raw gate output), `tmp/lbf/dispositions.md` (the item-by-item dispositions for ★1 to ★6 and ①②, including the round-two basis section) and `tmp/lbf/dd36b9b-trace.txt` (the `git log -S` trace for ①).
- **Round-two review target**: the evidence behind each of the four dispositions (the regenerated diff is non-empty, ADR-016 is back inside the table in both languages, the count matches the fresh reading, the trace output is in the pack) plus the four elements of **ADR-017**; the raw round-two receipt is kept at `tmp/lbf/review-round2.md` (`local-only`).
- **Sanitisation**: everything handed to the reviewers was filtered for user names, host names, IPs, credentials and tokens (directive §7.3); this slice involves no credential enumeration.
- **Blind isolation**: ⑦'s first round gets **no** ⑤ or ⑥ list; ⑥ is to be skipped, so there is no ⑥ list to attach.
- This slice **touches no** ownership, shutdown or identity semantics, so the §7.7 mandatory blind review is not triggered.

## 8. Falsifiability (mechanical checks and mutation tests)

- Mechanical checks: see §9.
- **Mutation tests** (each mutation is applied, the criterion is confirmed red, then the bytes are restored): ① change the OB-21 disposition word to one outside the closed set so `G6-6` should go red; ② delete the v1.16 §0.3 row so `G1-b` should go red; ③ delete one line on the ADR side so `G3(a)` should go red; ④ insert a line-count or file-count literal into a note so `G6-4` should go red; ⑤ **stated truthfully**: if the sentence "the `DR-1 not fixed` sub-item is removed" were deleted from the ★1 note, **no mechanical criterion would catch it** (only ⑦ and human review would).

## 9. Gate results

- `node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0` (`gate-exit=0`, `changed=13`; `G7-a` to `G7-d` all PASS and `G7-b` covers 76 artifact rows); `--check=G3-a.1 --scope=index` is also `TOTAL_FAIL=0`.
- `G1-b` PASS (`header=v1.16 last=v1.16 2026-09-29/2026-09-29`, `rows=17`); `G1` and `G5` print SKIP (not scripted, stated as such).
- `G3(a)` PASS (only pairs with a non-zero changed-line count are listed): the ADR pair `+2/-0`, the directive pair `+4/-3`, the `DEV_PLAN` pair `+1/-1`, the ledger pair `+28/-6` and the DR-1 report pair `+1/-1`; this report's pair is a new file and line-symmetric. **The ADR pair and the ledger pair are not whole-file positional mirrors** (the worktree carries pre-existing positional differences whose count equals HEAD's), so changed-line symmetry applies to them.
- `G4a` PASS (13 files BOM+LF); `G4b` PASS (`unknown=[]`); this slice's first-round `G2` and `G6-4` hits were cleared by rewording, not by extending any whitelist; `G5-a`, `G5-b`, `R2.5`, `G3(b)` and `G6-1` to `G6-6` all PASS.
- Base gates: `gofmt -l` empty, `go build ./...` exit 0, `go vet ./...` exit 0, `go test ./...` `ok=14 FAIL=0`.
- Mirrors: the directive pair and the DR-1 report pair have zero per-line mismatches; this report's pair has zero per-line mismatches (CN and `_EN` have equal line counts). This report's changed-line count moves with the version that lands ⇒ only "symmetric + zero mismatches" is recorded, never a self-referential count.
- **First-round input-pack defect (a master error, recorded truthfully)**: the diff sent for the first round, `tmp/lbf/review-input.diff`, was an **empty file** (generated after staging) ⇒ this is the **master mis-supplying the input**, **not a ⑦ problem**; it was regenerated with `git diff --cached`. This slice therefore consumed a **⑦ quota exception beyond the 1/1 budget (ADR-017)**.
- Reading basis: the readings above were **re-run on the spot** during the ⑨ native validation and those govern (⑤ already ruled that an old baseline must not be carried over).

## 10. ⑧c wrap-up cleanup record

- Whole-repository per-file encoding check (`enc-tree.js`): **3** violations, **all** of them OB-27 frozen-evidence-pack exemptions; **this slice's new files have zero violations**.
- Gate self-test: `SELFTEST: PASS 248/248`; `gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0` (per-criterion readings in §9).
- Temporary-residue inventory: `tmp/lbf/` retains 19 process files (plans and inventory, assertion scripts, raw outputs, the disposition list, the trace output, the first-round and round-two review-pack copies, and the `git diff` and gate readings), **not cleaned** per OB-49 and the standing user instruction, with disposition `local-only`. **Supplement (v1.23)**: the unit now measures 21 items - `findchar.js` and `ci-watch.txt` were added after wrap-up (both untracked); the full list, per-item sha256 and the archiving correspondence are in `docs/validation/evidence/directive-history/lbf/MANIFEST.md`.
- Landing consistency: `docs/validation/ledger-backfill.md` and `tmp/lbf/report-cn.md` have equal SHA-256 digests (the same holds for `_EN`).
- Mirror re-check: the ADR pair's whole-file positional mismatch count equals HEAD's (a pre-existing difference, not introduced here; the two rows this slice added sit in the same position in CN and `_EN` with symmetric counts); this report's pair has zero per-line mismatches.
- ⑧c disposed of 2 items (the ⑦ round-two Low findings): F2R-1's dangling forward reference and the stale ledger closeout line are now self-contained fact records, and F2R-2's in-pack intermediate reading is corrected to the fresh one.
- The raw round-two receipt is landed at `tmp/lbf/review-round2.md` (`local-only`); under the semantic boundary of §5.3 the fact-record writes of ⑧b do not trigger a further re-review (matching the same ruling in §10 of the `S1-REACH` report).

## 11. Cost and metering

| Stage | Role | Model | Calls | reason / retryOf |
| --- | --- | --- | --- | --- |
| ① | not enabled | none | 0 | the §5.2 documentation-form row does not enable it |
| ⑤ | pre-reviewer | `deepseek-v4-pro` | 1 | the §5.2 row marks it optional; this slice runs it once |
| ⑥ | independent scan | none | 0 | skippable per §7.2.1, pending ⑦ confirmation |
| ⑦ | final reviewer | `deepseek-v4-pro` | 2 | the first final review plus the round-two re-review (ADR-016 downgrade; the re-review is a quota exception, ADR-017) |
- Probes: **0**.
- Truth source for metering: the session call receipts (directive §7.5 hard rule R7.1).

## 12. Explicitly not executed

- commit and push (pre-authorised by the user in this turn; the master executes them at the ⑩ stop point, and gitee is never pushed to); the ⑥ independent scan (to be skipped, confirmed by ⑦); items ★7 to ★10 (deferred per the user's ruling); **the ⑦ round-two re-review** ran per the user ruling of 2026-09-29 (a quota exception, ADR-017), and its verdict plus the disposition of the two Low findings are in the closing bullet of §6 and in §10.

## 13. Next-step recommendations

- This slice changes no verdict ⇒ `OB-21` stays `待议`, AT-23 stays unreachable and T018 and S1 exit stay blocked.
- Scheduling follows the user's rulings: **T027 starts only when Codex is available (no downgrade)**; `DIRECTIVE-LEDGER-ARCHIVE` touches the `G6-6` criterion scope, so it needs Codex and queues separately.
- If the budget recovers, start with the T027 assembly (`gpt-5.3-codex`) and then assess T028 and the candidate line.
- The next slice registers the **OB-54 candidate**: the tail wording of the B3b paragraph at `docs/DEV_PLAN.md:139`; this slice only records it and does not write it into §12.4.

```artifacts
docs/validation/ledger-backfill.md	repo
docs/validation/ledger-backfill_EN.md	repo
docs/validation/evidence/LEDGER-BACKFILL-freedom-list.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/validation/dr-1-availability-probe-floor.md	repo
docs/validation/dr-1-availability-probe-floor_EN.md	repo
tmp/lbf/review-input.diff	local-only
tmp/lbf/gate-tree.txt	local-only
tmp/lbf/dispositions.md	local-only
tmp/lbf/dd36b9b-trace.txt	local-only
tmp/lbf/review-round2.md	local-only
```
> **Archiving note (v1.23)**: the `tmp/` items cited by this report were archived under `docs/validation/evidence/directive-history/` per §1.5.1 (units: `lbf/`; each with `MANIFEST.md` and `SHA256SUMS.txt`); the original paths are not rewritten and the MANIFEST is the mapping of record. (Paths already cleared by earlier slices are outside this archiving and belong to the pre-existing decoupling recorded as OB-77.)

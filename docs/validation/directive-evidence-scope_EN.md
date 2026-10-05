# Slice validation report · `DIRECTIVE-EVIDENCE-SCOPE` (v1.24)

> Type: directive-governance (evidence-tree criterion-scope exclusion + the new criterion G8's four sub-checks; touches the criterion body). Tier `[SLICE]`. Size M+. Probes 0.
> **Authorization boundary**: commit no / push no / gitee no / probes 0.

## 1. Goal and scope

- Main body **OB-79**: the target file set of §6.3's G1-a / G2 / G4a / G4b / G6 excludes `docs/validation/evidence/**` (**same shape, single source** as G7's existing exclusion, defined in the implementation layer); changes inside that tree are carried by the new criterion **G8**'s four sub-checks.
- Same batch: **OB-80** (archive-payload encoding caliber), **OB-76** (the §1.5.1 fourth class "rebuilt archive"), **OB-78** (the §1.5.1 class ③ illustration fix).
- **OB-77 assessed only**: worth a separate slice, low priority, criterion shape = a soft scan of changed-line `tmp/<path>` citations; not landed in this slice, with its §12.4 status word staying `To be discussed`.

## 2. Change summary

- Directive text (②a and the patch round): `docs/DELIVERY_DIRECTIVE{,_EN}.md`; lands §6.3's G8 row, the criterion-scope exclusion paragraph, the G8 known-boundaries paragraph, the scripting-status paragraph, plus §1.5.1's ② and ④ and the pointer sync.
- Implementation (②b): `tools/gates/lib/ctx.js` (including DEF-2's batched ignore and memoization), `tools/gates/lib/criteria/g8.js` (new), `tools/gates/lib/criteria/g4a.js`, `tools/gates/lib/criteria/g7.js`, `tools/gates/gate.js` (including DEF-1), `tools/gates/evidence-forms.txt` (new, 3 registrations).
- Tests (③): `tools/gates/selftest.js` extended to T43–T60.
- Fixtures: 28 new `tools/gates/testdata/fixtures/*.json` (on-disk names never end in `.md`).
- Evidence back-fill (①): a re-back-fill of `MANIFEST.md` and `SHA256SUMS.txt` under `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/`; payload bytes untouched, no retroactive rewrite.
- **Staging state**: at the ⑦ re-review moment = the two directive files staged, implementation and fixtures unstaged (the ⑦ re-review / final-review Low process note); after ⑧a landed, every declared artifact was `git add`-ed file by file so §6.3's `G7-b` tree pre-flight holds faithfully (not committed, not pushed, not pushed to gitee).
- **Commit-process note** (the ⑦ re-review / final-review Low): at commit time, stage file by file per §6.3's G5 and re-check staged ≡ declared set.

## 3. Pipeline and environment

① judgement (driver + V4 Pro) ⇒ ②a rule text (driver) ⇒ ②b implementation (`prfrail-implementer`) ⇒ ③ tests (`prfrail-tester`) ⇒ driver orchestration (evidence backfill) ⇒ ④ gates ⇒ ⑤ pre-review (V4 Pro) ⇒ ⑥ independent scan (MAI) ⇒ ⑦ blind (`GPT-5.3-Codex`) ⇒ remediation (the final review's four text-precision syncs landed) ⇒ ⑦ re-review ⇒ ⑨ native verification ⇒ ⑧a this report pair ⇒ ⑧b finalization ⇒ ⑧c close-out cleanup ⇒ ⑩ stopping-point receipt.

Environment: Windows 11 / Node 24; repository `github.com/larsonzh/prfrail`; baseline HEAD = `0c12abc` (this slice uncommitted).
Tiers: driver High; ①⑤ V4 Pro Max; ② `prfrail-implementer` / ③ `prfrail-tester` = V4.1 Flash; ⑥ MAI-Code-1.1-Flash; ⑦ `GPT-5.3-Codex` Extra High.

## 4. Anomalies and fallbacks

- **DEF-1** (pre-existing defect, not introduced here, fixed here): `gate.js`'s `run()` returned early inside the catch, skipping the stderr write ⇒ a thrown usage error exited 2 while printing nothing; fixed (2 lines) plus an assertion (the T57 group, capturing stderr independently).
- **DEF-2** (introduced by this slice, fixed here): the tree / index branch used one `git check-ignore` per path; one spawn measured 417 ms against a tree marked 2671 items ⇒ one enumeration extrapolates to 1,802,368 ms (about 30 minutes), and the G8-a ancestor walk to about 74 minutes; changed to ONE `check-ignore --stdin -z` with a `treeCache` memoization, and after the fix 1414–1626 ms (hot call 0.12 ms).
- **F-1 the ②b sub-agent's final report was not returned** (it stopped waiting on a background command): acceptance rests on the driver's mechanical re-verification, not on the implementer's self-reported readings.
- **Sequence deviation**: the ⑧a report pair was not completed before the ⑨ native verification (actual order = ⑦ re-review → ⑨ → ⑧a); under §5.1's fixed step order this item is a **sequence deviation**, stated truthfully. The ⑨ readings were taken before the ⑧a landing ⇒ the §7 table is marked as the **⑨ time point** readings (`changed=39` / `2` / `4`); the re-run after ⑧c re-checks as `changed=44`, `TOTAL_FAIL=0`, `SELFTEST: PASS 328/328`.
- **Freedom-list landing sequence deviation**: the freedom list was not landed after ⑥ / before ⑦, and all three ⑦ readings did not re-check it (it landed with ⑧a) — see F-14; under §1.7.3 and the OB-48 precedent, registered as this slice's deviation item.

## 5. Review conclusion summary (⑤⑥⑦)

- **⑤ pre-review** (V4 Pro): `PASS WITH FINDINGS`. M-1 (G8-c passing vacuously on a zero-row mapping section) ruled not to fix and recorded in this report; M-2 (B-1's strength cost not named in the body) ruled to take (b), landing §6.3's third point (iii); L-1–L-5 and L-7 not fixed, each declared; L-6(i) fixed the stale comment in `g8.js` (0 non-comment changed lines); L-6(ii)'s C.0n status row is written back this round.
- **⑥ independent scan** (MAI): `PASS WITH FINDINGS`, no blocker. Independently recomputed the back-fill arithmetic 744,747 + 259,061 = 1,003,808; independently verified that `check-ignore --stdin -z`'s exit 1 (no hit) is taken as data (`if (!(r.ok || r.code === 1))`); the shallow-clone `--scope=ci` guard is fail-closed (non-empty stderr, exit 2).
- **⑦ three readings** (`GPT-5.3-Codex`, 1 blind + 1 final + 1 re-review): blind review (no list attached) `BLOCK`, 5 findings (2 High: §0.3's trailing `fixtures` caliber too broad, C.0n's "do not rewrite any historical row wording" not column-scoped; 3 Medium/Low: G8's "judge changed lines only" caliber drift, the G8-d note region not disclosed, the split-staging process note) ⇒ final review `PASS WITH FINDINGS` (conditional pass), all four accepted with CN / EN final wording and the 5th kept as a process note ⇒ all four landed (§0.3 trailing / the C.0n boundary field / §6.3's G8 row trailing "judges the change set only (touched packs)" / the G8 row ④ gaining "the header note region") ⇒ re-review `PASS WITH FINDINGS` (independently confirmed the four landings, all three readings green, implementation-body consistent), no undeliverable fact found.
- **Mechanical ruling at landing (stated truthfully)**: the final review's (2) EN wordings lacked the bold that pairs with CN's `**现象列**` ⇒ G3(a.2)'s whole-file `**` count would fail (now 3261 = 3261); the documentation role stopped at the report and did not swap words itself, and the driver ruled option A1 on the mechanical mirror-consistency caliber (balance the EN bold), with EN taking `_EN`'s own §12.4 column name `Phenomenon`; the EN G8 row ④ keeps the existing numbering glyph `(iv)` to preserve G3(b)'s glyph count. This is a **driver mechanical ruling**, not the final review's own text.

## 6. Review input pack summary (with blind-review isolation proof)

- Input pack = `tmp/ges/REVIEW-BUNDLE-2.md` (final state, superseding the ⑤ point-in-time version); the ⑦ checks calibrate against that pack's section 0 state-change table.
- **Freedom-list gap**: the review input pack = `tmp/ges/REVIEW-BUNDLE-2.md`; the freedom list was not landed after ⑥ / before ⑦ (it landed with ⑧a), so all three ⑦ readings did not re-check it — see F-14.
- **Blind-review isolation proof**: the blind first round attaches **no list at all** and reads only the directive text and the implementation diff; its sampling basis is the §7.7 proportional check (the widest-touch one in the window), with the rationale recorded in C.0n; this slice's ⑦ = 3 readings, within the cap (1 blind + 1 final + 1 re-review).
- Sanitization: the review input contains only repository paths and text, no credentials and no external connection material.

## 7. Falsifiability and gate results

| Item | Reading |
|---|---|
| `--selftest` | `SELFTEST: PASS 328/328` (evolution: ②b 290 → ③ 326 → ③bis 327 → ③ter 328) |
| `--all --scope=tree` (⑨ time point) | `changed=39`, `TOTAL_FAIL=0` |
| `--all --scope=index` (⑨ time point) | `changed=2`, `TOTAL_FAIL=0` |
| `--all --scope=ci` (⑨ time point) | `changed=4`, `TOTAL_FAIL=0` |
| OB-79 core regression (`range:4284e94^..4284e94`) | `G1-a` / `G2` / `G4a` / `G4b` / `G6` each `TOTAL_FAIL=0`; `G4a` names 340 file(s) OK (312 excluded: evidence tree) |
| G8 real corpus (same range) | `G8-a` = 1 red, `G8-b` = 1 red, `G8-c` = green (15 packs), `G8-d` = green (312 file(s)); historical reds are not retroactive (§11.1) |
| After the back-fill | `--check=G8-b --scope=tree` = PASS, `1 pack(s) self-consistent` |
| DEF-1 | reproduced ⇒ `exit=2` with non-empty stderr (carrying the `gate.js:` prefix and the reason); fixed plus an assertion (the T57 group) |
| DEF-2 | before the fix 417 ms × 2671 marked items ⇒ extrapolated 1,802,368 ms; after the fix 1414–1626 ms (hot 0.12 ms); equivalence = sorted byte-equal, both-way empty diff; the tree branch `check-ignore` exactly 1 time; the T60 falsification demo `shipped=true / perFileMutant=false / doubleBatchMutant=false` |

The table above is the ⑨ time point reading; after the ⑧a/⑧b report pair, the freedom list and the ledger row landed, the change set grew to `changed=44`; measured after ⑧c, `--all --scope=tree` and `--all --scope=index` both read `changed=44`, `TOTAL_FAIL=0`; `--selftest` still reads `SELFTEST: PASS 328/328`; the class-③ process artifacts under `tmp/` have been cleared (only `.gitkeep` remains).

This slice's landed commit = **`025158f`**; **CI `run 37278691747`** is green on both legs at first run (both legs `SELFTEST: PASS 328/328`, `GATE REPORT scope=ci changed=44`, `TOTAL_FAIL=0`); that run is an **observation record**, with no further commit added for this slice afterwards.

**Mutation checks**: not applicable to a semantic rule change (this slice changes the criterion scope and adds a criterion body); the falsification means = the `--selftest` assertions over 28 fixtures, the OB-79 core regression, the DEF-1 / DEF-2 reproduction-vs-fix contrast, and the T60 falsification demo.

## 8. Known boundaries and residual risks

- Verbatim quote of §6.3's "G8's known boundaries (stated truthfully)" three sentences: **(i) G8 is a shape-level guard** — it only checks "the shape of paths inside the evidence tree and in-pack self-consistency", not content semantics; **(ii) "live code disguised as a frozen pack" is not recognized by G8** — a frozen pack may legitimately carry script-class artifacts (precedent: the `*.ps1` under `docs/validation/evidence/b2-2026-09-17/machine-ops/`), and that residual risk is carried by the ⑤⑥⑦ review chain; G8 must not be claimed to cover it. **(iii) the payload's encoding form is not verified by G8** — G4a already excludes that tree, and G8-c only checks the "suffix-and-encoding mapping" rows `MANIFEST` has recorded; for **packs without that section** (legacy packs) the payload encoding form has no criterion coverage, and that residual strength gap is carried by the ⑤⑥⑦ review chain.
- **B-1**: a pack with no mapping section has its payload encoding form covered by no criterion; already declared by (iii).
- **B-2**: DEF-2 is guarded by T60; `index` / `ci` are transitive coverage of the same code path.
- **B-3**: G8-b does not verify the MANIFEST "bytes" column; partial-corruption shapes of `SHA256SUMS` are not separately asserted.
- **B-4**: shape-level guard, does not recognize live-code disguise, does not judge content semantics.
- **B-5**: historical reds are not retroactive (the `4284e94` G8-a / G8-b two reds remain permanently).
- **B-6**: M-1 — a zero-row mapping section = a declared zero normalization, with no record duty, the same as "no section".
- **B-7**: untouched packs are not judged; nesting follows "nearest ancestor"; G8-d scans only the first 15 entries of the note region; a rebuilt archive's sha256 is only shape-checked.
- ⑥ two Low: the cache key depends on the repo-relative canonical path; the G8 shape-level boundary.
- ⑦ two Low: the double-point exclusion guard's maintenance coupling; the long runtime of readings over a large historical range.

## 9. Explicitly not executed

- **M-1 not fixed** (reason: recognizing an "explicit declaration" is content semantics, beyond the shape-level guard).
- **L-1 not fixed** (G8-c recognizing only lowercase 64-hex = an undeclared tightening), **L-2 not fixed** (an uncommitted entry requiring a `/` absolute path = a tightening beyond the text), **L-3 not fixed** (G8-d not checking the date, triggering on "presence"), **L-4 not fixed** (the registration table "root level only" depends on manual discipline), **L-5 not fixed** (tree / range payload enumeration differs for symlinks and tracked-but-ignored), **L-7 not fixed** (duplicate MANIFEST entries are de-duplicated).
- **OB-77's criterion not landed**; **no retroactive rewrite**; **legacy packs not fixed**; **commit authorization: not granted for this slice** (not committed, not pushed, not pushed to gitee); **probes 0**.

## 10. Deviations and degrees of freedom

- Same source as the freedom list, item by item in `docs/validation/evidence/DIRECTIVE-EVIDENCE-SCOPE-freedom-list.md` (F-1 to F-13).
- F-1 the ②b sub-agent's final report was not returned; F-2 the ① back-fill touched two metadata items; F-3 DEF-1; F-4 DEF-2 (introduced here); F-5 the registration table reads by endpoint; F-6 the legacy-unit measurement correction; F-7 no §6.3 wording rewrite needed + the four sync re-checks all ruled unnecessary; F-8 at ②a time the G8 implementation did not yet exist; F-9 the cost vector does not enter §0.3 or C.0n; F-10 the ⑤ findings list and dispositions; F-11 the B-2 assertion and its last-centimeter blind spot; F-12 the four text-precision syncs the ⑦ three readings produced; F-13 the documentation role catching the G3 conflict + the driver mechanical ruling.

## 11. ⑧c close-out cleanup record

- This slice's ⑧c is run by the driver **after this report is finalized**: `tmp/ges/`'s review process artifacts (`REVIEW-BUNDLE-2.md`, the readings and scripts) are one-off process products of §1.5.1 **class ③**, removed immediately after use and not committed.
- A frozen evidence pack takes no part in the per-file encoding rewrite (§5.1 existing exemption); the evidence tree's encoding verifiability is carried by the `MANIFEST`'s two sha256 records.

## 12. Cost and accounting

- Realized vector: ①×1 / ②×3 / ③×3 / ⑤×1 / ⑥×1 / ⑦×3 (blind 1 + final 1 + re-review 1) / ⑧×10 + the driver's several rounds.
- Caliber: ①×1 means the architect plan round only (`deep-reasoner`, read-only, produces the plan); the evidence backfill round is run by the ⑧ documentation role, counted under ⑧, never under ①.
- Caliber: ⑧ = all actual invocation rounds of the documentation/report role, including the BLOCKED round that did not land and the retry rounds that produced no change due to `Agent error`.
- Within the slice definition's caps (①≤1, ⑤≤1+1, ⑥≤3, ⑦ = 1 final + 1 re-review + 1 blind); probes 0; no Haiku calls; no quota exception, no downgrade ADR.

## 13. OB-77 assessment conclusion

- **OB-77**: worth a separate slice, low priority; criterion shape = a **soft scan** of changed-line `tmp/<path>` citations; not landed in this slice, with its §12.4 status word staying `To be discussed`.

## 14. Follow-up recommendations (including tier advice)

- Follow-up recommendations (recommendations only, not registered): B-1 may need a new criterion (a new OB candidate); the B-2 blind spot; **OB-82**'s §7.7 window-count anchor; legacy packs gaining canonical shape (a separate slice).
- Tier advice: for a follow-up criterion-body slice, ⑦ still uses `GPT-5.3-Codex`, not downgraded.

## 15. ⑩ stopping-point report (including the write-back check line)

- Executed through **⑨**; after ⑧b is finalized the driver runs **⑧c** and issues the **⑩ stopping-point receipt**; the user authorized `commit` + `push` in the **same round** (only `origin`, no gitee push); the landed commit **`025158f`** (44 items, `+2395 −50`) was pushed as `0c12abc..025158f`; **CI `run 37278691747` is green on both legs at first run** (`Go windows-latest` / `Go ubuntu-latest` both `SELFTEST: PASS 328/328`, `GATE REPORT scope=ci changed=44`, `TOTAL_FAIL=0`).
- **Write-back check** (the OB-63 formulation of §11.2): validation report = **written back** (this file pair); `DEV_PLAN` = **not applicable** (a directive-governance slice); ledger = **written back** (the `REMAINING_SLICES{,_EN}` completion row, recording the report path only); ADR = **not applicable** (this slice has no downgrade and no ruling ADR).
- ⑨ readings: all three scopes tree / index / ci `TOTAL_FAIL=0`; `SELFTEST: PASS 328/328`; `--check=G8-b --scope=tree` = PASS.

## 16. Artifacts list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-evidence-scope.md	repo
docs/validation/directive-evidence-scope_EN.md	repo
docs/validation/evidence/DIRECTIVE-EVIDENCE-SCOPE-freedom-list.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/MANIFEST.md	repo
docs/validation/evidence/directive-history/tmp-lifecycle-workdir/SHA256SUMS.txt	repo
tools/gates/gate.js	repo
tools/gates/lib/ctx.js	repo
tools/gates/lib/criteria/g4a.js	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/lib/criteria/g8.js	repo
tools/gates/evidence-forms.txt	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/exclusion-negative.json	repo
tools/gates/testdata/fixtures/exclusion-positive.json	repo
tools/gates/testdata/fixtures/g8-a-freedom-list.json	repo
tools/gates/testdata/fixtures/g8-a-member.json	repo
tools/gates/testdata/fixtures/g8-a-real-registered.json	repo
tools/gates/testdata/fixtures/g8-a-registered.json	repo
tools/gates/testdata/fixtures/g8-a-stray-root.json	repo
tools/gates/testdata/fixtures/g8-b-backslash-sums.json	repo
tools/gates/testdata/fixtures/g8-b-hash-mismatch.json	repo
tools/gates/testdata/fixtures/g8-b-malformed-sums.json	repo
tools/gates/testdata/fixtures/g8-b-manifest-truncated.json	repo
tools/gates/testdata/fixtures/g8-b-missing-list.json	repo
tools/gates/testdata/fixtures/g8-b-unlisted-payload.json	repo
tools/gates/testdata/fixtures/g8-b-untracked-conflict.json	repo
tools/gates/testdata/fixtures/g8-bom-in-tree.json	repo
tools/gates/testdata/fixtures/g8-bom-outside.json	repo
tools/gates/testdata/fixtures/g8-c-arch-not-in-payload.json	repo
tools/gates/testdata/fixtures/g8-c-arch-not-in-sums.json	repo
tools/gates/testdata/fixtures/g8-c-bad-hash.json	repo
tools/gates/testdata/fixtures/g8-c-no-section.json	repo
tools/gates/testdata/fixtures/g8-c-ok.json	repo
tools/gates/testdata/fixtures/g8-d-bad-sha.json	repo
tools/gates/testdata/fixtures/g8-d-body-only.json	repo
tools/gates/testdata/fixtures/g8-d-missing-basis.json	repo
tools/gates/testdata/fixtures/g8-d-ok.json	repo
tools/gates/testdata/fixtures/g8-ok-pack.json	repo
tools/gates/testdata/fixtures/g8-t50-in-tree.json	repo
tools/gates/testdata/fixtures/g8-t50-outside.json	repo
tmp/ges/REVIEW-BUNDLE-2.md	local-only
tmp/ges/REVIEW-BUNDLE.md	local-only
tmp/ges/	local-only
```

Note: the evidence tree's frozen packs are declared by unit `MANIFEST` rows (per-item payload details live in each unit `SHA256SUMS.txt`); the `local-only` rows are `tmp/ges/`'s review process artifacts, cleared at ⑧c.

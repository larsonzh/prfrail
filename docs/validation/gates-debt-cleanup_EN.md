# Slice validation report · P67 GATES-DEBT-CLEANUP

- Directive version: v1.66 (at step 0; step F will bump to v1.67)
- Date: 2026-11-12
- Class: governance `[SLICE]`; size tier L; probe budget 0 (Windows-local governance slice)
- Status: 🔄 **In progress** (started 2026-11-12; step 0 registered; stop point ⑩ awaiting authorization)
- Controller model: `DeepSeek V4 Pro (deepseek)` (ADR-041)
- Evidence pack: `docs/validation/evidence/P67-GATES-DEBT-CLEANUP/`

## 1. Executive summary

This slice is the layer-2 governance one-shot debt cleanup (#1 g4b cache LRU + cap / #2 gate pitfalls 27–30 formalization / #3 §12.4 sampling-record block settlement / #4 ⑥ prompt hardening). Progress (2026-11-12): step 0 done; the ① five artefacts landed (`DESIGN: DONE`); ② A–G all done (per-step checks in the evidence pack `03-test/`); ③ consistency `PASS`; ④ three scopes `TOTAL_FAIL=0` plus `SELFTEST: PASS 736/736`; ⑤ r1 `PASS WITH FIXES` (6 Low) → fixes → r2 `PASS WITH FIXES` (1 Low, budget 2/2 exhausted, the literal fix closed and handed to ⑦ for confirmation); ⑥ r1 `INDEPENDENT SCAN: PASS`; ⑦ blind r1 `FINDINGS` (1 High + 1 Medium, remediated) → ⑦ final review r1 `FINDINGS` (1 Medium report assembly + 1 Low `CACHE_CAP` pin) → T130 landed (`SELFTEST: PASS 737/737`) → this ⑧a is the Medium remediation → ⑦ re-review `RE-REVIEW: PASS` (at ⑧b time); ⑨ `SELFTEST: PASS 737/737` plus three scopes `TOTAL_FAIL=0`; batch-assist GLM-5.3-Flash 3 runs (① cache-key enumeration / ② seven-pit comparison / ③ ⑥ three-round enumeration, archived under `02-batch/`).

- ⑥ independent scan (actually run): r1 = `INDEPENDENT SCAN: PASS` (A–E compliant, 0 findings; verbatim archive in the evidence pack `04-reviews/06-scan-r1.md`); the verbatim transcription is the `sixstage` fence at the end of this report (§11.3), with the outside-fence note right after it.

- **Anomaly register (user ruling 2026-11-12)**:
  1. **① overstep (pitfall 6)**: the ① stage recognized the contradiction between 3.1.3 and briefing 【七】5 "items 1–26 zero changes", did not escalate the 3.4 stop-upgrade boundary, and autonomously chose 3.1.3 first. Disposition: the user ruled option ① — pitfall 6 exits the #2 scope entirely; item 21 stays verbatim; pitfall 6 is only a report observation (§8 known boundaries + v55 fact pack); "the 3.1.3 clause on pitfall 6 (folding it into the item-21 disposition example) is voided as a whole by the hard constraint of briefing 【七】5 'items 1–26 zero changes'; this entry is the user's same-turn 2026-11-12 ruling".
  2. **① overstep (mtime boundary)**: the ① stage evaluated the semantic question "does mtime count as a timestamp" but did not escalate it as a 3.4 stop-upgrade boundary, and adopted the permissive reading. Disposition: the user ruled on 2026-11-12 that the permissive reading holds and mtime is in-bounds (reason = literal-context reading of 3.1.2; the existence of a potential conservative reading is registered); this is a boundary clarification, not a §1.7.4 temporary-authorization overreach, and carries no separate ADR; it must not serve as precedent.
  3. **R counting**: the two deviations above are not remediation reruns (R does not +1); they are registered here as "user-discovered" (per the P58 §12.1 pilot-metric "user-discovered" counting); this slice's "user-discovered" = 2.
  4. **Temp-artefact D count = 1 (controller violation, self-heal established — user direct confirmation 2026-11-12; this entry is a one-time breakthrough of the §5.1 literal "when the violator is the controller, the self-heal verdict takes effect only after same-turn independent review-layer confirmation", because at step-B time ⑤ had not started and no same-turn independent review layer was available; the remaining §5.1 clauses (self-report + pre-snapshots + self-delete + zero-residue evidence) were executed strictly, with no breakthrough; this exception applies only to this slice and must not serve as precedent)**: event facts: a stray `8` (425 bytes) appeared at the repo root (outside the §5.1 location closure), caused by the controller's nested-quoting-broken `cmd /c "powershell …"` tail-read command; three pre-snapshots (porcelain `?? 8` / tmp recursive 115 / TEMP hits 0) → precise `Remove-Item -LiteralPath .\8` → zero residue (`Test-Path`=False, no `??` line, tmp recursive still 115).
  5. **#3 landing deviation (⑤ F3 fix)**: design-r1 §4 item ② "back-fill the traceable fields" was implemented as "each slice's ⑦ blind-review facts take that slice's report review-scan record as the single source of truth, and this block does not transcribe them" (a historical debt, §11.1 no retrospective back-fill); dropping the C part is already noted in design-r1 §4, and this area now registers it (reason = the dual constraint of the §7.7:798 body sentence and F-8, see design-r1 §4).
  6. **Execution-discipline deviation (user ⑩ correction instruction 2026-11-12)**: during ② the controller skipped landing the step-A/D/E/F check files (against 【六】2 "independent verification per step + landing `03-test/step-X-check.md`"), found by the ⑦ blind review (1 High) and back-filled afterwards (step-A/D/E/F-check.md landed and registered in MANIFEST/SHA256SUMS); this is a controller execution-discipline deviation, not an artefact defect.
  7. **R counting (§9.2 calibre, user ⑩ correction check 2026-11-12)**: R = 3 — ⑤ r1 (6 Low findings) → fixes → ⑤ r2 = 1; ⑦ blind FINDINGS → fixes = 1; ⑦ final FINDINGS → fixes → re-review = 1. R ≥ 2 ⇒ executed per the §5.2 high-risk row (⑤ re-check mandatory: done with 0 Medium+; ⑨ all green); R = 3 does not reach the §9.2 "R > 3 pause for manual review" threshold.
  8. **Terminal false-interrupt observation (§9.3 known pattern, user ⑩ correction instruction 2026-11-12)**: during the ④ re-run of the tree scope the terminal showed one `^C` false positive (the process had actually finished and the final reading is `TOTAL_FAIL=0`, see `06-gates/gates-r1.txt`); this slice was not affected, and it is registered here.
  9. **The pitfall-7 §12.4 landing sentence was added after the ⑩ correction, without a ⑦ re-review (user ⑩ ruling 2026-11-12)**: in nature = a record-status landing (the v54 §一 already executes this rule + points to the effective §13.5 template + a meta declaration), not new rule semantics; does not touch §7.2.1 ③; per the §5.3 spirit "re-run after fixes" targets rule-semantics fixes, and a record-status supplement needs no ⑦ re-run. This exception applies only to this slice and must not serve as precedent.

## 2. #1 four-option comparison and final landing

| Dimension | Option A: hard N with write refusal | Option B: LRU eviction | Option C: TTL + cap | Option D: LRU + TTL |
|---|---|---|---|---|
| Implementation complexity | Low | Medium | Medium | High |
| Invalidation-semantics correctness | Correct but **cache stagnation** (a new sha always misses) | Correct (eviction = losing a hit chance, never a wrong hit) | Correct | Correct |
| Space bound | Hard N | Hard N | Soft (depends on the expiry-check timing) | Hard N + TTL backstop |
| Effect on existing entries | Below N nothing is cleaned, afterwards writes are refused | After put, evict the oldest mtime entries one by one | Evicted progressively on expiry | Same as option B |
| Risk | Medium: the cache degenerates into a dead cache | Low | Medium: needs a time source and timing semantics | Low |
| Recommendation | Low | **Recommended** | Medium | Unnecessary |

**Conclusion**: final landing = option B (LRU, N=100, mtime access order). An old-sha entry never hits again once `resolvedSha` moves forward, so evicting the oldest mtime evicts a pure placeholder entry, which makes the TTL of options C/D redundant. The mtime boundary (whether mtime counts as a "timestamp") was ruled in-bounds by the user on 2026-11-12 (see the §1 anomaly register item 2 and §7), reason = the literal-context reading of 3.1.2. Source: `01-arch/design-r1.md` §1, §2.

## 3. #1 criterion-body change (g4b cache LRU + cap)

- Change surface: `tools/gates/lib/criteria/g4b.js` numstat `93/2`, the only criterion-code change in this slice; new module-level constant `CACHE_CAP = 100`, entry-name regex `CACHE_ENTRY_RE`, and pure function `evictToCap(dir, cap, statFn, unlinkFn)`, all exported through `internals` for the selftest. Source: `02-impl/impl-r1.md` §1, `03-test/step-A-check.md` §1.
- Access order: `diskStore.get` touches the entry with `fs.utimesSync` after a successful `JSON.parse`; `diskStore.put` calls `evictToCap(dir, CACHE_CAP)` after a successful write and rename. Source: same, §2, §3.
- Eviction rule: only entries matching `^[0-9a-f]{64}\.json$` are counted, unlinked one by one in ascending `mtimeMs`, removed count = `entries.length − cap`; a stat or readdir failure returns 0 and never throws. Source: same, §2, §4.
- Six boundaries held: LRU + cap does not touch the criterion main flow; the five-component key diff is zero; the value shape stays `{v: 2, corpus: string[]}` (mtime is filesystem metadata and never enters the JSON value, `validValue` unchanged); a hit equals a cold run; fail-open on every path; corpus only, never a conclusion. Source: `01-arch/design-r1.md` §2, `03-test/step-A-check.md` §3.
- Test delta: `tools/gates/selftest.js` numstat `363/0` (pure insertion, existing T1–T124 and F1 untouched verbatim); T125–T129 cover eviction order and entry-name filtering, the three cap branches, real put eviction, hit equivalence, and the four fail-open paths; the ⑦ r1 Low fix appended T130 (pinning `CACHE_CAP` to the literal 100 against constant drift). Source: `03-test/step-B-check.md` §2, `tools/gates/selftest.js`.
- SELFTEST: `722/722` (step 0, `06-gates/step0-selftest-r1.txt`) → `736/736` (step-B time point, `06-gates/gates-r1.txt`) → `737/737` (after T130 landed, re-run at ⑧a). T125–T130 net +15.

## 4. #2 gate-pitfall formalization (numbering continued at 27)

- Seven-pitfall disposition: 4 are formalized as items 27–30; pitfall 6 (GBK mojibake) exits the #2 scope entirely per the user‘s 2026-11-12 ruling and stays a report observation only; pitfall 7 (the §一 rolling-window ambiguity) adds no directive entry and moves to the v55 fact pack. Source: `01-arch/design-r1.md` §3, `03-test/step-C-check.md`.
- Four landings: item 27 unregistered CJK characters triggering G4b (three-way disposition + the §6.3 criterion-scope exclusion rule) / item 28 aggregate-count wording such as "共 N 条" triggering the G6-4 aggregate-count shape / item 29 time-point wording such as "候选" being misread after migration into the formal table / item 30 line-tail truncation (the CN/EN truncation asymmetry evidence). Source: same.
- Continuation and zero change: the new items are numbered continuously from 27; existing items 1–26 stay semantically zero-change, and `03-test/step-C-check.md` §3 proves by a single pure-insertion hunk from `git diff -U0` that items 1–26, the table header and the post-table register-scope block are untouched verbatim; item 21 received no disposition example. Source: `03-test/step-C-check.md` §1, §3.
- Bilateral symmetry: `docs/DELIVERY_DIRECTIVE.md` and `_EN` are `4/0` each, both sides 2391 → 2395 (equal), and numbering 1–30 is continuous with no gap and no duplicate. Source: same, §1, §2.

## 5. #3 §12.4 sampling-ledger settlement (option A′)

- Factual conflict: the window-start phrase "窗口自 v1.32 起计" coexists in two directive places (the §7.7 OB-82 sub-clause and the §12.4 tail block); changing only the block contradicts the body sentence (violating §11.2 "one fact in one place"), while changing the body hits freedom-list F-8 "不改准则 §1–§12 已有条款". Source: `01-arch/design-r1.md` §4.
- A′ landing: the window start is untouched (both occurrences kept verbatim); inside the §12.4 block the time-point sentence becomes a factual sentence, the traceable fields are back-filled (whether a slice was sampled is not historically traceable ⇒ recorded as a historical debt), the unfulfilled "decide and register before start" duty over v1.32–v1.66 is declared a historical debt (§11.1 no retrospective back-fill), and per-slice registration resumes from P67. Source: same; `03-test/step-D-check.md` §2.
- Deviation register: design-r1 §4 item ② "back-fill the traceable fields" was implemented as "each slice‘s facts take that slice‘s report review chain as the single source of truth, and this block does not transcribe them"; part C of the briefing-recommended A+C (window reset) is dropped, reason = the dual constraint above; both are registered in the §1 anomaly register item 5. Source: same; `03-test/step-D-check.md` §4.
- Boundary check: the §7.7 block diff is zero and the block carries no window-start wording conflicting with §7.7; numstat is `3/1` on both sides and both sides 2394 → 2396 (equal). Source: `03-test/step-D-check.md` §1, §3.

## 6. #4 ⑥ prompt hardening (appendix B.6)

- Five hardening points: ① the five A–E section headings are fixed at `## Section A`–`## Section E`; ② the conclusion-line position and value (exactly two = the first line of Section A + the block end, values equal); ③ the Section B `SEVERITY` row agrees with the conclusion line; ④ `UNCOVERED` stands as its own sentence; ⑤ stale-reading identification (⑥ must declare the input package version and timestamp explicitly and must not cite outdated readings). Source: `01-arch/design-r1.md` §5.
- Text alignment only: ①②③ fall inside the shapes already defined by §7.2.2 and G7-f and add no new G criterion; ④ strengthens existing wording; ⑤ is added in input-face instruction form and is not written into §7.2.2. Source: same.
- ⑥-only note: B.6 states the A–E shape applies only to ⑥, while ⑦ still follows the §4.1 closed-set conclusion line; B.7.4 is unchanged. Source: same.
- Zero-change check: the §7.2.2, §6.3, §4.1 and appendix B.7 diffs are all 0; numstat is `5/1` on both sides and both sides 2396 → 2400 (equal); after the ⑤ r1 F5 fix the heading wording matches G7-f literally (CN/EN `+1/−1` each). Source: `03-test/step-E-check.md` §2, §3.

## 7. Known boundaries

- The ① stage had two deviations in executing the 3.4 stop-upgrade boundary, disposed by the user's same-turn 2026-11-12 ruling (see the §1 anomaly register).
- mtime boundary ruling (user, 2026-11-12): the true intent of 3.1.2 "if a timestamp / eviction marker must be added" is "adding a timestamp / eviction-marker field inside the cached value structure"; mtime is "using existing filesystem metadata" and is not in the forbidden range; all six boundaries hold. **Not precedent-setting**: future same-kind boundaries still require independent rulings.
- Core discipline (user ruling 2026-11-12): from now on, when a 3.4 stop-upgrade boundary is hit, ①/② must stop and escalate literally, without autonomous disposition.
- The #1 criterion-body change is limited to LRU + cap (option B, N=100, mtime access order); if a G4b-cache-class problem reappears, a deeper change may be needed (e.g. changing the cache-key semantics).
- Step C formalization scope: only appending items 27–30; pitfall 6 exits per the user's 2026-11-12 ruling (observation only); existing items 1–26 stay strictly zero-change.
- §5.1 timing limitation (user ruling 2026-11-12): at this slice's step-B time no same-turn independent review layer was available (⑤ had not started), so the controller-violation event was confirmed by the user in lieu; this is a timing limitation, not a rule defect; future same-kind cases still follow the §5.1 letter (same-turn review-layer confirmation) and must not cite this slice as an example.
- Pitfall 7 (the §一 rolling-window ambiguity): at the ⑤ F2 time-point it was registered as "carried to the v55 fact pack"; after the user's 2026-11-12 ⑩ correction it now lands in the §12.4 end block (the ④ sentence, record-status); the v55 fact pack records the settled rule per §13.3.15.
- **The pitfall-7 landing sentence was added after the ⑩ correction**: a record-status landing without a ⑦ re-review; user same-turn ruling 2026-11-12; nature and basis in §1 item 9.
- ⑤ r1 reading time-point (F6 fix): at ⑤ r1 time the report pair had 53 lines and `06-gates/gates-r1.txt` was a pre-step-F intermediate reading; the final numbers follow the ⑨ measured values, and the MANIFEST / SHA256SUMS plus report readings are recomputed at archive time.

## 8. Explicitly not done

- No criterion-body change: the single exception is the #1 LRU tightening (§3).
- No change to `CONTRACTS`, schemas, fixtures or product code (`internal/`, `cmd/`); no change to CI configuration.
- No change to existing directive clauses §1–§12, §13.1–§13.3 or the §13.5 body; no change to the §7.2.2 or §7.7 body; no change to the v50–v54 snapshots.
- No change to the semantics of gate pitfalls 1–26 (pitfall 6 is observation only, item 21 got no disposition example).
- No ruling made on the user‘s behalf: both boundary dispositions in this slice came from the user‘s same-turn 2026-11-12 ruling (§1 anomaly register).
- No commit, no push, no `git add`; no new slice started; no VM, no proxy, no Sol double-run.

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-batch/glm-batch1.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-batch/glm-batch2.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-batch/glm-batch3.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/consistency-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-A-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-C-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-D-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-E-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-F-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/05-prereview-r2.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/07-blind-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/07-rereview-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/06-gates/gates-r1.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/06-gates/gates-r2.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/06-gates/step0-selftest-r1.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/MANIFEST.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/SHA256SUMS.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/freedom-list.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/input-list.md	repo
docs/validation/gates-debt-cleanup.md	repo
docs/validation/gates-debt-cleanup_EN.md	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/lib/criteria/g4b.js	repo
tools/gates/selftest.js	repo
```

```sixstage
## Section A

INDEPENDENT SCAN: PASS

Input package = P67 GATES-DEBT-CLEANUP, 2026-11-12; this review only reads the materialized evidence and the directive wording in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L720-L760), and the P67 proof bundle under [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md), [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md), and [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md). No file creation, modification, deletion, or command execution occurred.

## Section B

(no findings)

## Section C

- The P67 entry explicitly records the in-scope change set and the boundary that matters for this review: “不改 §7.2.2 / §6.3 / §7.7 正文; 不改 §13.1–§13.3 / §13.5; 不改 v50–v54.” See [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115). That is the key fact for the B.6 hardening question: the change is limited to prompt framing and evidence packaging, not a new rule body.
- The LRU design is judged in-bounds under the six boundaries in [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md). The critical condition is the user ruling recorded in that file: mtime is used as filesystem metadata for access ordering, not as a stored value, while the value shape stays exactly as the previous corpus record.
- The implementation memo confirms the actual behavior in [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md): successful parse triggers a touch on get; successful write triggers eviction on put; eviction only examines matching cache files and uses mtime order; all failure paths remain fail-open. That matches the cache contract in [tools/gates/lib/criteria/g4b.js](tools/gates/lib/criteria/g4b.js).
- The selftest evidence in [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md) records T125–T129 and a green SELFTEST: PASS 736/736 result. This directly covers the concerns around bound enforcement, hit equivalence, and fail-open handling. The same file also states that the change did not touch the rule body in [tools/gates/selftest.js](tools/gates/selftest.js).
- The “7 pit merge” is in-bounds as described in the P67 design memo: the new items are a debt cleanup of four specific edge cases, and the file expressly preserves the no-change set for the earlier items. The “坑 6 退出” and “坑 7 转 v55” branches are treated as local routing decisions, not as semantic expansion of the core gate logic. This is consistent with the directive’s “not changing the rule body” boundary in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115).

## Section D

- If the eviction guard is removed, the cap and oldest-first slice checks in T125–T127 would go red; the design memo explicitly marks those as the failure points.
- If the touch-on-hit guard is removed, the LRU ordering assertions T128a–T128c fail, because hit semantics would no longer remain equivalent to a cold run.
- If the fail-open catches are removed, T129a–T129d fail; the implementation memo states the failure paths are intentionally silent and must not alter verdict semantics.
- If the B.6 prompt hardening is treated as creating new semantic rules instead of a documentation-only surface, the directive’s own no-change statement in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115) would be violated. The evidence package keeps it within the “prompt only, rule body unchanged” envelope.
- UNCOVERED: the materialized evidence does not prove cross-platform NTFS/exFAT timestamp resolution or every rare Windows file-handle race under all OS conditions; that is an edge-case gap, not a contradiction with the current PASS.

## Section E

- Cannot verify from the evidence bundle alone: real runtime behavior across all Windows filesystems and multi-process races beyond the documented fail-open cases. The input package is static and intentionally read-only; the directive requires that this scan only read the materialized evidence, not execute commands or modify files.
- What would overturn PASS:
  - any evidence that mtime is folded into the stored value semantics instead of remaining metadata;
  - any CN/EN asymmetry in the B.6 wording or a change to the rule body in §7.2.2, §6.3, or §7.7;
  - any materialized test result showing the selftest or gate assertions are no longer green;
  - any evidence that the “1–26 unchanged” claim is false for the actual authored patch.

INDEPENDENT SCAN: PASS
```

The five sections A–E above are the verbatim transcription of ⑥ r1 (§11.3); the original text is in the evidence pack `04-reviews/06-scan-r1.md`.

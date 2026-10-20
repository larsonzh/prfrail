# Slice `DIRECTIVE-BOUNDARY-CLARIFY` validation report (directive v1.41)

> Slice: `DIRECTIVE-BOUNDARY-CLARIFY`; directive version v1.41; tier L; this file is the **⑧a draft** (the ⑨ and ⑩ readings and the commit state are back-filled by ⑧b).
> Base: `d:\LZProjects\prfrail`; at ⑧a start the staged set held 6 paths, the index equalled the worktree (`git diff --stat` printed nothing) and there were no untracked files.
> Bilingual: this file is the CN authority; `docs/validation/directive-boundary-clarify_EN.md` is a line-by-line mirror, and both sides carry an equal number of lines.

## 1 Change summary

Before ⑧a the change set held 6 paths; after ⑧a it holds 9. `tools/**` is untouched by this slice (`git diff --cached --numstat` lists documents only). The additions and deletions below are measured readings; the two new reports take their own line count as the addition (`+279/-0`, equal on both sides).

| Path | + | − | Notes |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 36 | 11 | v1.41 body: 7 dispositions + three boundary clauses + version package + ledger flips |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 36 | 11 | `_EN` strict mirror |
| `docs/t027/REMAINING_SLICES.md` | 20 | 0 | P16 status row |
| `docs/t027/REMAINING_SLICES_EN.md` | 20 | 0 | `_EN` in step |
| `docs/validation/directive-boundary-clarify.md` | 279 | 0 | this file |
| `docs/validation/directive-boundary-clarify_EN.md` | 279 | 0 | the mirror |
| `docs/validation/directive-gates-patch.md` | 1 | 1 | embedded source note at the end of the §8 line |
| `docs/validation/directive-gates-patch_EN.md` | 1 | 1 | `_EN` in step |
| `docs/validation/evidence/DIRECTIVE-BOUNDARY-CLARIFY-freedom-list.md` | 53 | 0 | execution freedom list (evidence-tree root) |

**artifacts fence** (9 paths, all declared `repo`):

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-boundary-clarify.md	repo
docs/validation/directive-boundary-clarify_EN.md	repo
docs/validation/directive-gates-patch.md	repo
docs/validation/directive-gates-patch_EN.md	repo
docs/validation/evidence/DIRECTIVE-BOUNDARY-CLARIFY-freedom-list.md	repo
```

The 7 landing points with measured line numbers (CN and `_EN` share the numbers):

| Item | Disposition | Landing point (`docs/DELIVERY_DIRECTIVE.md`) |
|---|---|---|
| OB-23 | Resolved | §7.2.1 L661 standalone paragraph; §7.8.2 L743 convergence signal; §4.2.1 L418 heading "six items" + L425 item 6 |
| Group 1 / OB-13 | Resolved | §5.1 L480 timing note |
| Group 2 / OB-25 | Resolved | §3.4 L357 scope sentence |
| Group 2 / OB-26 | Resolved | §3.4.2 L369 baseline definition; §12.5 L1077 new gap row |
| Group 3 / OB-44 | Resolved | §7.3 L693 item ③ (write boundary) |
| Group 3 / OB-45 | Resolved | §7.3 L693 item ② (data-source wording) |
| Group 4 / OB-50 | Resolved | §2.5 L316 parenthetical at the end of the sentence |

Ledger and version package (same file): all seven §12.4 disposition cells are flipped (OB-13 L945 / OB-23 L959 / OB-25 L963 / OB-26 L965 / OB-44 L984 / OB-45 L985 / OB-50 L990); the header at L3 reads `版本：v1.41。日期：2026-10-20。`, the v1.41 changelog row of §0.3 is at L87, and appendix C.0ac starts at L1796 (its deliverable row at L1801 already carries the P16 status row of `REMAINING_SLICES{,_EN}`).

## 2 Execution pipeline

| Step | Role / carrier | Model | Conclusion | Rounds | Notes |
|---|---|---|---|---|---|
| Step 0 | `prfrail-implementer` | controller model | complete | 1 | P16 status row +20/-0; the 8.2 check passed; the v1.41 header joined the ② batch (`G1-b` linkage) |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | complete | 1 | not split; after comparing the three OB-23 shapes it recommended shape C; four work groups; found that the original plan lacked §7.3 |
| Batch helper ① | `deep-reasoner` | `GLM-5.3-Flash (glm)` | complete | 1 | OB status inventory: all 7 items are `待议`, and the CN / EN line numbers correspond one to one |
| Batch helper ② | `deep-reasoner` | `GLM-5.3-Flash (glm)` | complete | 1 | section and OB reference check: all 18 section numbers exist, no dangling section number |
| ② implementation | `prfrail-implementer` | controller model | complete | 1 | 6 paths, per-group gates green |
| ③ test | `prfrail-tester` | controller model | complete | 1 | three scopes, self-test, six shapes, criterion-by-criterion |
| ④ gate | controller | controller model | complete | 1 | three scopes `TOTAL_FAIL=0` |
| ⑤ pre-review r1 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 1 Medium (the C.0ac deliverable row missed P16) + 5 Low |
| ② fix | `prfrail-implementer` | controller model | partly complete | 1 | the dispatch lost its network mid-way and part of the change landed; the controller verified mechanically and continued, FIX-1 to FIX-4 landed |
| ⑤ pre-review r2 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS | 1 (narrowed) | Medium CLOSED, 3 Low CLOSED, the reason for not fixing FIX-5 accepted |
| ⑥ independent blind scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | PASS | 1 | blind input, zero Medium+ |
| ⑦ terminal review r1 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | all five groups pass, zero Critical / High / Medium / Low |
| Batch helper ③ | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | complete | 1 | bilingual language-side comparison over 6 points: no difference in negation polarity, quantifiers or terminology |
| Batch helper ④ | not dispatched | none | not applicable | 0 | the MANIFEST needs a per-file sha256, and the batch helper has no write permission and cannot compute hashes reliably |
| ⑧a draft | `prfrail-documenter` | controller-selected | complete | 1 | this file, the `_EN` mirror and the freedom list |

- Rework convention: ② had 1 main dispatch plus 1 fix dispatch (interrupted mid-way), with no re-dispatch and no rollback.
- The duration field is missing (historical): this slice stored no start or end timestamps, and none is invented.
- The three-scope readings of ④ are independent runs, not results reused from one process.

## 3 Exceptions and fallbacks

| # | Event | Impact | Disposition |
|---|---|---|---|
| D1 | the "v1.41 header" of step 0's item 8.3 cannot land on its own (`G1-b` requires the header version to equal the last version in §0.3) | mechanical constraint, not a freedom | merged into the ② main change batch |
| D2 | the ① landing list lacked §7.3 (the kickoff did not list it) | the two matters, review-layer write boundary and data source, had no landing point | extended by user ruling B, limited to exactly those two matters |
| D3 | the ② fix dispatch lost its network mid-way (client `UND_ERR_SOCKET`) | part of FIX-1 to FIX-4 had landed | the controller verified the on-disk state mechanically before continuing; the six shapes stayed equal and `G3` still PASSed |
| D4 | the §7.2.1 reading was lifted out of a list item into a standalone paragraph | text shape changed, meaning unchanged | the ⑤ r1 Low fix, which also lands the user's ruling A "write the two matters separately" |
| D5 | batch-helper priority: the user's ruling this round says `GLM-5.3-FlashX` first | the first two runs had used `GLM-5.3-Flash` | not re-run (output still usable); from the third run on the corrected priority applies |
| D6 | batch helper ④ (MANIFEST draft) was not dispatched | none | reason in the §2 table, not "no applicable workload" |
| D7 | `G9-e` emits INFO for the historical path `tmp/v113/review-r2.md` in the OB-44 phenomenon cell | that path is a historical fact reference, not an evidence artifact | not blocking |
| D8 | ⑤ r2 found an evidence-pack defect: the FIX-2 EN contrast in `review-5r2/01-fix-lines.txt` shows `<missing>` | the capture script pattern cannot match `Skippable` | an evidence-pack defect, not a document defect (⑤ read the source directly and confirmed) |
| D9 | 1 `[Info]` item from ⑤: the EN v1.41 changelog row uses ①–⑦ numbering while older rows use (1)–(6) | style difference | not acted on |
| D10 | the kickoff gave 961 / 963 / 1794 as the start lines for OB-25 / OB-26 / C.0ac, while the measured lines are 963 / 965 / 1796 | citation offset, 2 in every case | the measured lines prevail and are used in §1 |
| D11 | ⑧a cannot replay the `runSubagent` parse of the transcript inside the session: `main.jsonl` holds 264 bytes and the session store truncates `assistant_response` at 5000 bytes | the cost table of §10 is a transcription of dispatch facts, not a mechanical audit | registered as-is; the mechanical audit is back-filled by ⑧b and this table is not the final cost fact |
| D12 | the FIX-5 family (4 `_EN`-only spots) is left unfixed | 4 spots where the EN side is not seamless with the CN authority | a single-side edit trips `G3(a.1)` FAIL, so they wait for a slice that touches both sides |
| D13 | the freedom list was produced later than ⑦ r1 (§1.7.3 requires after ⑥ and before ⑦) | the ⑦ r1 package carries no list (`tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-7/` has none) | same type as D10 of `GATES-PATCH`, the second occurrence; registered as-is and the package is not rebuilt |
| D14 | dispatch-order deviation of ③: this slice ③ (`prfrail-tester`) was dispatched only **after** ⑤ r1/r2, ⑥ and ⑦ r1 (transcript offset `29098208`), and its output passed in full after the late dispatch | see D14-1 / D14-2 / D14-3 below | recorded (**OB-107**) |
| D15 | local semantics of `--scope=ci`: `ctx.js` uses `baseRef=HEAD^1` | **before a commit it measures the previous slice** (this slice reading is `changed=4`), not the current staged set; the criteria for this slice 9 files are carried by `index` / `tree` | `ci` runs against this commit only after the push |
| D16 | the `G4b` block on ⑧a: the first draft wording carried a CJK character outside the whitelist, so `G4b` reported a non-empty `unknown` | `G4b` turned red under `--scope=tree` and blocked the ⑧a delivery | the wording was **rewritten to avoid the character** (a synonym already inside the whitelist) and the **whitelist was not grown** (`tools/gates/cjk-newwords.txt` still 37 entries, `tools` untouched); after the fix `newCjk=591 unknown=[]` |
| D17 | the cost table of ⑧a lacks a mechanical basis: ⑧a cannot replay the transcript parse inside the session, so its §10 marked every row as pending ⑧b at the time | the cost readings of §10 are a transcription of dispatch facts, not a mechanical audit | this stage back-filled it from the transcript audit (see §10); this is the first empirical gap of clause ① among the three self-check clauses fixed by this repository on 2026-10-20 |
| D18 | line-number offset in the kickoff message: the fact list given to ⑧a put OB-25 / OB-26 / C.0ac at 961 / 963 / 1794, while the measured lines are 963 / 965 / 1796 (all two higher, because rows such as OB-104 / OB-105 were inserted into §12.4 before OB-106 was registered) | a citation offset that changes no conclusion | the report uses the measured readings throughout (see §1) |
| D19 | the FIX-5 family (EN-only residue, not fixed): in `docs/DELIVERY_DIRECTIVE_EN.md`, **L87 / L693 / L985** (`may only read` should be `may only be written as`) and **L743** (tier (i) of §7.8.2: `must carry` does not explicitly carry the perfect aspect of 「须**已**给出」) | 4 spots where the EN side is not seamless with the CN authority | left for a later slice that touches both sides, to land as one pair; reason: an EN-only edit makes the CN / `_EN` added-minus-deleted counts asymmetric and turns `G3(a.1)` red, while the CN side is itself correct and a semantics-free compensating edit for a shape gate is not acceptable |
| D20 | batch-helper priority correction: the user receipt this round corrects the priority to "always prefer `GLM-5.3-FlashX`, falling back to Flash only when exhausted / unsuitable / withdrawn" | the first two runs (candidate ① / ②) had already used Flash | the output stays usable and is **not re-run**; from the third run (candidate ③) FlashX is used |

- **D14-1 fact**: ③ was dispatched after ⑤ r1/r2, ⑥ and ⑦ r1 (offset `29098208`); the output passed in full after the late dispatch, and ③ independently found 3 observations (`--check=base` exit 2 / the `HEAD^1..HEAD` semantics of `--scope=ci` / the 408 static `check(` call sites against 588 at runtime).
- **D14-2 substantive impact**: **the review input packages of ⑥ and ⑦ r1 do not carry the ③ test results** - they review "the implementation", and "the tests" are not in the object under review; the 3 observations ③ found independently afterwards should have entered the review scope before ⑥/⑦. ⇒ **the `PASS` of ⑥ and the `PASS` of ⑦ r1 rest on incomplete input**; the gap is registered as-is and **must not** be covered up by "everything passed after the late dispatch".
- **D14-3 residue**: the ③ output entered no review stage input ⇒ the two-way "implementation + tests" review chain of this slice **closes only inside this slice report, not inside the grounds of the ⑥/⑦ conclusions**; **OB-107** is registered (Observing).

## 4 Environment and tier

| Item | Value |
|---|---|
| OS | Windows (local `d:\LZProjects\prfrail`) |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| Tier | L |
| ⑥ | mandatory (touches the review-reading semantics of §7.2.1 and §7.8) |
| ⑦ quota | 1+3 |
| ⑤ quota | 1+1 |
| Batch helper | ADR-027 (limited to four task classes) |
| Probes | 0 |

## 5 Review conclusions (⑤⑥⑦)

| Round | Id | Severity | Conclusion / handling | Disposition |
|---|---|---|---|---|
| ⑤ r1 | M1 | Medium | the C.0ac deliverable row missed the P16 status row of `docs/t027/REMAINING_SLICES{,_EN}.md` (the staged set held that pair at +20/-0, so this was under-reporting) | Fixed |
| ⑤ r1 | L1–L4 | Low | FIX-1 to FIX-4: content-side fixes (wording and landing points) | Fixed |
| ⑤ r1 | L5 | Low | FIX-5 family: 4 `_EN`-only spots that are not seamless with the CN authority | Recorded |
| ⑤ r2 | — | — | Medium CLOSED; 3 Low CLOSED; the reason for not fixing FIX-5 accepted | Closed |
| ⑥ r1 | — | — | blind input, zero Medium+ | Recorded |
| ⑦ r1 | — | — | all five groups pass, zero Critical / High / Medium / Low | Recorded |
| Batch helper ③ | — | — | the 6 bilingual language-side points show no substantive difference; 2 slight EN differences join the FIX-5 family | Recorded |

- The ⑤ r1 Medium has a mechanical closure check: `git diff --cached --numstat -- docs/t027/` yields `20 0` for both files, and the C.0ac deliverable row at L1801 already carries that path pair.
- The 4 EN-only spots of the FIX-5 family: at L87 / L693 / L985 the wording `may only read` should be `may only be written as`; at L743 `must carry` does not explicitly carry the perfect aspect of the CN `须已给出`. All four need an EN-only edit, which `G3(a.1)` would fail; the CN side is itself correct and a compensating edit would pollute the frozen body, so they wait for a paired slice.
- Other ⑦ observations (no finding): the single-hunk diff of the closed report is append-only and touches neither the verdict line nor any reading; the internal consistency of the new rules (the three tiers decoupled from §7.2.1 item ③, §7.3 not conflicting with neighbouring clauses, §3.4.2 consistent with §3.4.3 and §12.5) all pass.

## 6 Review input packages (including blind-isolation proof)

| Package | Location | Contents | Manifest |
|---|---|---|---|
| ⑤ r1 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-5/` | brief / design / full diff / numstat / three-scope gate / self-test tail / criterion text / report diff | `MANIFEST.md` |
| ⑤ r2 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-5r2/` | brief / fix hunks / three-scope gate | `MANIFEST.md` |
| ⑥ r1 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-6/` | brief / full diff / numstat / status / three-scope gate / self-test / readings | `MANIFEST.md` |
| ⑦ r1 | `tmp/DIRECTIVE-BOUNDARY-CLARIFY/review-7/` | brief / design / full diff / numstat / status / three-scope gate / self-test / ledger items / section numbers / single-line report diff / readings | `MANIFEST.md` |

- The ⑥ and ⑦ briefs were dispatched blind: no ⑤ conclusions and no earlier findings attached.
- Input sanitisation: the packages carry only the slice definition, design, test report, diff, numstat, gate readings, self-test output and criterion sources.
- The ⑦ r1 package did not carry the freedom list (see D13 in §3), so ⑦ did not review against it.
- The file counts come from mechanical enumeration: ⑤ r1 has 11, ⑤ r2 has 6, ⑥ has 10 and ⑦ has 14 (each including its `MANIFEST.md`).

## 7 Falsifiability

### 7.1 Falsifiable assertions (this slice is documentation-only, with no runtime probe)

| # | Assertion | Counterexample (which reading changes once it is deleted) | Mechanical check |
|---|---|---|---|
| E1 | the write-boundary clause of §7.3 exists | once item ③ is deleted, OB-44 has no carrying sentence and the item-by-item check of ⑤ cannot pass | searching `不得创建、修改、删除任何文件` yields 1 hit |
| E2 | the §7.2.1 reading is a standalone paragraph | back at the end of a list item, the user's ruling A is not landed and the ⑤ r1 Low stays open | searching the opening sentence `跳过面证据充分性的口径` yields 1 hit |
| E3 | the §4.2.1 heading matches the item count | with the heading back at "five items" it contradicts item 6 | heading and item 6 coexist in one file |
| E4 | all seven §12.4 disposition cells are flipped | any row still at `待议` contradicts the v1.41 claim in §0.3 | read the seven disposition cells one by one |
| E5 | the §3.4 scope sentence is written once only | had §1.6 a second copy, the two could drift | searching `[ROUTINE]` 档由主控亲自执行 yields 1 hit |

### 7.2 Mechanical re-computation

| Reading | Value | Source |
|---|---|---|
| `git diff --cached --numstat` | 6 paths (before ⑧a) | command output |
| CN and `_EN` line count | lines=1876, equal on both sides | mechanical count |
| six shapes | heads 123 / bold 4554 / bars 3007 / blank 453 / ticks 4636 | mechanical count, fully equal on both sides |
| `SELFTEST` | PASS 588/588 | `node tools/gates/selftest.js` |
| `G4b` | newCjk=591, unknown=[] | `--check=G4b --scope=tree` |

### 7.3 Mutation payloads

- `tools/**` is untouched by this slice (`git diff --cached --numstat` lists no `tools/` path), so there is no criterion or test code change and therefore no mutation payload.
- This slice adds no assertion: `SELFTEST` stays at 588/588, the same value as the `GATES-PATCH` close-out.

### 7.4 Coverage boundaries and known Low items

| Id | Content | Owner |
|---|---|---|
| L-1 | the 4 EN-only spots of the FIX-5 family | Recorded (see §11) |
| L-2 | the timing deviation D13 of the list | Recorded (see §3) |
| L-3 | the citation offset D10 is corrected to the measured lines | Resolved |

- The conclusions of this report do not cover the semantic equivalence of `_EN`; that face is carried jointly by ⑤ / ⑦ and the shape comparison of `G3(a.2)`.

## 8 Gate results

| Command | Result |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 588/588` |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0` (exit 0) |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0` (exit 0, `changed=4` is the previous slice commit; semantics in §9) |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0` (exit 0) |
| `--check=G4b --scope=tree` | `newCjk=613 unknown=[]` (`whitelist=37`) |
| `G5(d)` (`--scope=index`) | `9 path(s) declared`; `declared-not-touched=10` (INFO, the fence is a landing-point snapshot); `addMode=explicit` |
| base `go test ./...` | run inside `--all` |

The rows above were ⑧a point-in-time readings and are now replaced by the definitive readings, which the controller measured after ⑧b and before ⑧c.

Per-criterion self-check readings (definitive, `--scope=tree`):

| Criterion | Reading |
|---|---|
| G7-a | PASS |
| G8-a | PASS |
| G1-a | PASS |
| G2 | PASS |
| G6-4 | PASS |

- `G7-b` reports 3 new paths absent from the index: that FAIL was cleared after the controller staged 9 paths, and the definitive reading is `TOTAL_FAIL=0`; the unknown character of `G4b` under `--scope=tree` is cleared by rewording (see D8).

**CN and `_EN` six shapes (`docs/DELIVERY_DIRECTIVE.md` vs `_EN`):**

| Shape | CN | EN |
|---|---|---|
| lines | 1876 | 1876 |
| headings | 123 | 123 |
| `**` | 4554 | 4554 |
| `|` | 3007 | 3007 |
| blank lines | 453 | 453 |
| backticks | 4636 | 4636 |

> The table above is an ⑧a point-in-time reading; its purpose is to prove that CN and `_EN` share the same six shapes, not to guarantee the current absolute values.
- Both sides are UTF-8 with BOM and LF.

## 9 ⑨ native verification and commit state

⑨ was executed and back-filled by ⑧b; every reading in this section is mechanical.

**Per-criterion exit codes (`--check=<ID> --scope=index`, 18 items)**: all exit 0.

| Criterion | Exit code | Notes |
|---|---|---|
| `G1-a` / `G1-b` / `G1` / `G2` / `G3-a` / `G3-b` | 0 | PASS |
| `G4a` / `G4b` / `G5-a` / `G5-b` / `G5-c` | 0 | PASS |
| `G5-d` | 0 | INFO (`no changed .github/agents/*.agent.md in scope`) |
| `G5` / `R2.5` / `G6` / `G7` / `G8` / `G9` | 0 | PASS |
| `base` | 2 | `unknown criterion id: base` ⇒ **the second empirical instance of OB-105** (registered in v1.40, not fixed by this slice) |

Two points must be read together with the table above:

- the exit 0 of `G5-d` carries INFO rather than PASS: this change set holds no `.github/agents/*.agent.md`, so the criterion has nothing to judge.
- `base` is not one of the 18: `--list` names it as the 19th id and `--check` does not accept it; the 18 items on the `--check` side all exit 0.

**Three scopes `--all`**: `index` / `tree` / `ci` all report `TOTAL_FAIL=0` and exit 0. **The qualifier that must be read with it**: the implementation basis of `--scope=ci` is `HEAD^1..HEAD` (`baseRef=HEAD^1` in `ctx.js`), so **before a local commit** it measures the **previous slice commit** (this slice reading is `changed=4`); the criteria for this slice 9 files are carried by `index` / `tree`, and `ci` runs against this commit only after the push.

**Pristine-worktree replay** (`HEAD` = `29309d3`): `git worktree add --detach tmp/DIRECTIVE-BOUNDARY-CLARIFY/pristine HEAD` → export the `git diff --cached --binary` patch → `git apply --index` inside the pristine tree → `--all --scope=index` yields `TOTAL_FAIL=0` and `selftest` yields `588/588` ⇒ the change set is self-consistent; afterwards `git worktree remove --force` and the patch was deleted, leaving only the main worktree in `git worktree list`.

**Whole-tree encoding audit**: `scanned=3332` tracked text files, 0 violations outside the evidence tree; the 16 historical BOM / EOL deviations inside the evidence tree (archived packs and a BOM compatibility fixture) are unrelated to this slice and are not counted.

**Reference resolvability**: the added lines of the staged set carry 34 section references, 0 dangling.

| Item | State |
|---|---|
| `selftest.js` | `SELFTEST: PASS 588/588` |
| ⑩ stop-point receipt | not executed (this slice ⑩ stop point awaits authorisation) |
| commit and push | neither executed; gitee untouched; CI not triggered |

## 10 Cost and accounting (R7.1)

**Mechanical audit basis**: the table below is the transcript parse of this slice `runSubagent` calls, in order of occurrence; "#" is the call index inside this conversation and the offset / line columns are the byte offset and the line number in the transcript JSONL.

| # | Stage | `agentName` | `model` | transcript offset | line |
|---|---|---|---|---|---|
| 104 | ② step 0 (P16 row) | `prfrail-implementer` | controller model (unspecified) | 26629059 | 28405 |
| 105 | batch helper ① (status inventory of the 7 OBs) | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 26709083 | 28504 |
| 106 | batch helper ② (section / OB reference check) | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 26890115 | 28565 |
| 107 | ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 27060493 | 28663 |
| 108 | ② implementation (7 items + version package) | `prfrail-implementer` | controller model | 27293578 | 28777 |
| 109 | ⑤ r1 | `prfrail-tester` | `GLM-5.3 (glm)` | 27639436 | 29090 |
| 110 | ② fix (FIX-1 to FIX-4; the dispatch lost its network mid-way) | `prfrail-implementer` | controller model | 27823412 | 29234 |
| 111 | ⑤ r2 (narrowed) | `prfrail-tester` | `GLM-5.3 (glm)` | 28001752 | 29437 |
| 112 | ⑥ blind scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 28227831 | 29536 |
| 113 | ⑦ r1 terminal review | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 28284792 | 29616 |
| 114 | batch helper ③ (language-side check of five bilingual spots) | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 28408082 | 29889 |
| 115 | ⑧a | `prfrail-documenter` | controller model | 28493087 | 29945 |
| 116 | ⑧ (G4b rewording workaround) | `prfrail-documenter` | controller model | 28922118 | 30385 |
| 117 | ③ test (dispatched after ⑤ - see D14) | `prfrail-tester` | controller model | 29098208 | 30598 |

**Accounting conclusion**: ① 1 / ② 3 / ③ 1 / ⑤ 2 / ⑥ 1 / ⑦ 1 / ⑧ 3 / batch helper `GLM-5.3-Flash` 2 / batch helper `GLM-5.3-FlashX` 1 / probes 0.

- ② counts 3 = #104 + #108 + #110; ⑧ counts 3 = #115 + #116 + this ⑧b (this call happened after the audit point above, so it is not in the table).
- The two GLM models are counted separately (`GLM-5.3-Flash` 2, `GLM-5.3-FlashX` 1) and are never merged.
- The table above supersedes the ⑧a transcription, closing the mechanical-basis gap of D11 and D17.

**The four batch-helper candidates, item by item**:

- candidate ① (item-by-item status inventory of the seven OBs) = **dispatched** (`GLM-5.3-Flash`, offset 26709083): reading = all 7 items `待议` / `To be discussed`, with CN and EN line numbers corresponding one to one (OB-13 L934, OB-23 L948, OB-25 L952, OB-26 L954, OB-44 L973, OB-45 L974, OB-50 L979, all pre-change ledger lines).
- candidate ② (directive section / OB reference check) = **dispatched** (`GLM-5.3-Flash`, offset 26890115): reading = all 18 section numbers really exist with identical CN / EN line numbers; the section numbers cited in the seven phenomenon cells have no dangling reference; `OB-106` was not registered then (this ⑧b registers it); the G6-5 rule text and the `g6.js` implementation already agree.
- candidate ③ (language-side check of the five bilingual spots) = **dispatched** (`GLM-5.3-FlashX`, offset 28408082): reading = 6 comparisons, with 0 difference in negation polarity, 0 in quantifiers / tier words and 0 in terminology; 2 slight EN-side differences (the perfect aspect of §7.8.2 tier (i), the verb of §7.3 item ②) join the FIX-5 family.
- candidate ④ (draft of the review-package MANIFEST) = **not dispatched**: every MANIFEST row is `bytes + sha256`, so it must be produced by a process that can hash each file and write to disk; the batch-helper role (the `deep-reasoner` carrier) holds read-only tools, has no write permission and cannot compute or read back hashes reliably ⇒ the task is unsuitable for a batch helper and is generated mechanically by the controller with a `node` script instead (all four packages `review-5` / `review-5r2` / `review-6` / `review-7` of this slice carry a generated `MANIFEST.md`).

This round reports no credits figure, only the checkable model and round facts.

## 11 Explicitly not executed

- no `git commit` / `git push`; gitee untouched; CI not triggered.
- no ⑩ stop-point receipt (this slice ⑩ stop point awaits authorisation).
- this slice adds no assertion: the `tools/` directory is untouched ⇒ `SELFTEST` stays at 588/588.
- the bodies of the §2.2 / §2.4 closed sets, the §5.2 matrix, the §7.5 figures and the §6.3 criterion bodies are unchanged.
- the four `_EN`-only spots of the FIX-5 family are not fixed (reason in §5 and D19).
- batch helper ④ (MANIFEST draft) was not dispatched (reason in §10).

## 12 Next steps

1. ⑩ stop point: awaits authorisation; the controller issues the ⑩ receipt (dispatch register, the four batch-helper candidates item by item, the OB-23 ruling, the four terminal readings and OB-106).
2. After the user authorises, commit and push, then watch CI go green on both legs.
3. The FIX-5 family lands as a pair in the next slice that touches both CN and `_EN`.

## 13 Execution freedom

Each row cites `docs/validation/evidence/DIRECTIVE-BOUNDARY-CLARIFY-freedom-list.md` (the list is authoritative; this table does not change its meaning).

| # | Choice | Landing point |
|---|---|---|
| F1 | OB-23 takes shape C (equivalent coverage + explicit declaration + three equivalence tiers, with tier (i) never exempt from argument) | §7.2.1 L661 |
| F2 | §7.3 extends to ⑤ / ⑥ / ⑦ and is limited to the two matters "write boundary" and "data source" | §7.3 L693 |
| F3 | the reading lands as a standalone paragraph, not as a parenthetical at the end of a list item | §7.2.1 L661 |
| F4 | the N of OB-13 is parameterised, never a hard-coded number | §5.1 L480 |
| F5 | the `[ROUTINE]` tier runs the diagnostics scan itself under the controller, with no subagent and no added cost | §3.4 L357 |
| F6 | OB-26 takes the conservative reading: with no snapshot the whole scan counts as new | §3.4.2 L369 |
| F7 | the scope is declared once in §3.4 and not copied into §1.6 | §3.4 L357 |
| F8 | OB-50 adds no new body text and only carries the landing note | §2.5 L316 |
| F9 | the C.0ac deliverable row gains the P16 status row (the fix shape of the ⑤ r1 Medium) | C.0ac L1801 |

The blocker is cleared — `G4b` now reports an empty `unknown` (rewording only, no new whitelist entry), so this draft is delivered as `DOCS: DONE`; `tools/gates/cjk-newwords.txt` is unchanged.

DOCS: DONE

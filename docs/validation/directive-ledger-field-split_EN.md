# DIRECTIVE-LEDGER-FIELD-SPLIT slice validation report (§9.3 counting-sequence pinning and ledger field unification)

> Slice: `DIRECTIVE-LEDGER-FIELD-SPLIT` (P31). Date: 2026-10-26. Size: M. Tier: directive-governance slice. Directive version: v1.52 → **v1.53**. Result: six body amendments landed; `index` / `tree` / `ci` scopes and `selftest` all green.

## 1 Change summary

- Goal: pin the `§9.3` counting sequence into mechanically decidable facts, and fix the field mismatch between the ledger layer and the reports.
- Directive version: `v1.52` → `v1.53`, date `2026-10-26`.
- Landing points: 8 paths, `+96 / -6` in total; six of them are directive body amendments (CN line numbers and `_EN` at the same positions).
1. `:3` version line → `Version: v1.53. Date: 2026-10-26.`
2. `:99` §0.3 gains the v1.53 entry, carrying "④ §9.3 '② advance cleanup' disposition cell gains the 'single-file deletion' case and the 'Closed member set' sentence" and "register OB-124 (pilot continues deferred, time 22)".
3. `:513` §5.1 D-count paragraph appends, in the same paragraph, the "four legal spellings" and "historical numbering is not retroactive" sentences (detailed / bare / implicit / field absent; pick one, uniform spelling is not enforced).
4. `:724` §7.3 gains the "verbatim-locating assertion" item: when giving `file:line`, first assert the block-head line number, then quote the on-disk text verbatim.
5. `:893` §9.3 "② advance cleanup" disposition cell gains the "single-file deletion is registered on its own and does not enter the cross-slice pattern sequence" and "Closed member set" sentences; the symptom cell itself is unchanged.
6. `:952` §11.3 gains the "time-point caliber for readings pasted into reports" subsection: readings pasted inside the fence and in §8 take the ⑩ receipt and the CI readings as their current values.
- Ledger `docs/t027/REMAINING_SLICES` (CN / `_EN`): the `:63` status row reads `🔄 in progress`, and the P31 detail block starts at `:1049`.
- Two historical reports (`directive-gates-cleanup` / `directive-carrier-adjust`, CN and `_EN`) each gain one correction note: "temporary-artifact D accumulates 0 times in this slice (backfilled, 2026-10-26). This report body is not rewritten." (per §11.1, add only, never rewrite the body.)
- This report pair and the freedom list are new artifacts of this slice; their landing points are in the `artifacts` fence below:

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-cleanup.md	repo
docs/validation/directive-gates-cleanup_EN.md	repo
docs/validation/directive-carrier-adjust.md	repo
docs/validation/directive-carrier-adjust_EN.md	repo
docs/validation/directive-ledger-field-split.md	repo
docs/validation/directive-ledger-field-split_EN.md	repo
docs/validation/evidence/DIRECTIVE-LEDGER-FIELD-SPLIT-freedom-list.md	repo
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-read.js	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-noagent.js	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-classify.js	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/gate-post-low-index.txt	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/gate-post-low-tree.txt	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/gate-post-low-ci.txt	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/selftest-post-low.txt	local-only
```

## 2 Pipeline

| Step | Role (agentName / model) | Status | Rework |
|---|---|---|---|
| Step 0 chartering | Dispatcher | Done | 0 |
| ① six-group comparison table | 3 batch helpers + 1 architecture helper (no agentName passed) | Done (carrier violation, see ch. 13) | 0 |
| ② Implementation | `prfrail-implementer` / `DeepSeek V4.1 Flash (deepseek)` | Done | 0 |
| ③ Consistency check | `prfrail-tester` / `DeepSeek V4.1 Flash (deepseek)` | r1 FAIL → r2 pass | 1 |
| ④ Gates | Dispatcher | three scopes and selftest all green | 0 |
| ⑤ Pre-review | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` | r1 2 Medium → r2 closed → r3 High judged closed (quota exception) | 2 |
| ⑥ Independent scan | `independent-reviewer` / `MAI-Code-1.1-Flash (copilot)` | r1 High 1 + Medium 2 → r2 PASS | 1 |
| ⑦ Final review | `independent-reviewer` / `GPT-5.3-Codex (copilot)` | PASS, no findings | 0 |
| ⑧a / ⑧b | `prfrail-documenter` / `DeepSeek V4.1 Flash (deepseek)` | Done (⑧a 1 / ⑧b 1) | 0 |
| ⑧c Cleanup | Dispatcher | to be executed | — |
| ⑨ / ⑩ | Dispatcher | to be executed | — |

- Elapsed time: stages were not individually timed, and no estimate is written (not metered).

**Dispatch table (copied verbatim; must not be rewritten into inferred values)**

| Stage | Transcript line | `agentName` | `model` |
|---|---|---|---|
| ① batch enumeration | L84335 | `<absent>` | `<absent>` |
| ① batch five-dimension comparison | L84404 | `<absent>` | `<absent>` |
| ① batch impact inventory | L84589 | `<absent>` | `<absent>` |
| ① six-shape comparison table | L84719 | `<absent>` | `DeepSeek V4 Pro (deepseek)` |
| ② Implementation | L84950 | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` |
| ③ Consistency check | L85205 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` |
| ③ Re-run | L85414 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` |
| ⑤ Pre-review r1 | L85501 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` |
| ⑤ Re-review r2 | L85598 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` |
| ⑥ Scan r1 | L85685 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` |
| ⑤ Re-review r3 | L85779 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` |
| ⑥ Scan r2 | L85883 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` |
| ⑦ Final review | L85914 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` |

- The table above is derived mechanically from the transcript by `tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-read.js` and `dispatch-noagent.js` (stdout saved in the same directory); it must not be rewritten into inferred values.

## 3 Anomalies and fallback records

- ③ r1 FAIL, 1 finding: `docs/t027/REMAINING_SLICES_EN.md:1069` carries 2 backtick pairs against 3 on the CN side (the freedom list's literal path was missing); the root cause is an EN wording defect from step 0. r2 = 0 FAIL.
- ⑤ r1 (2 Medium): M-1 = the §5.1 EN example was an abstract placeholder `D = 1 (tier plus its basis)` that is not equivalent to the CN "detailed" example ⇒ changed to `D = 1 (⑧b non-self-healing tier; user-approved continuation)`; M-2 = the P31 declaration said "observation H is closed" with no landing point ⇒ changed to "registered only, not handled; not closed in this slice".
- ⑤ r2: both Medium items above closed; one new Low (the objective line said "six items" against five actual, and H was unbounded) ⇒ changed to "five items … H registered only, handed to a later slice".
- ⑤ r3 (quota exception approved by the user in the same turn): High judged closed; 2 Low reported: (a) the v1.53 ledger row and `C.0an` still said "five places / one sentence" where "six places / two sentences" is correct; (b) the `_EN` term order `Member closed set` disagreed with the existing `Closed member set`. Both were changed to "six places / two sentences" and `Closed member set` as suggested.
- ⑥ r1: verdict `INDEPENDENT SCAN: FINDINGS`; High = the §9.3 counting sequence has no single authority for "who the members are" (historical reports carry several calibers: P23 / P28 / P30); Medium 2 = `§11.1` "frozen" is defined nowhere in the document; Medium 3 = §7.3 "verbatim-locating assertion" is in tension with "data-source adjacency".
- ⑥ r2 (narrowed re-run): High closed and no new dual-source conflict introduced; verdict `INDEPENDENT SCAN: PASS`.
- Missing dispatch arguments (registration 1): the 3 batch helpers in ① passed neither `agentName` nor `model` (transcript L84335 / L84404 / L84589, mechanically read as `<absent>`), so all fell to the session default carrier; the nature is "missing arguments", **not** "self-downgrade" (GLM was never attempted and no call failed).
- Carrier unspecified in ① (registration 2): the six-shape comparison call (L84719) passed `model` as `DeepSeek V4 Pro (deepseek)` but no `agentName`, so it likewise fell to the session default.
- R2.1 gap (registration 3): R2.1 constrains `model` only, never `agentName`, so "passing `model` explicitly" can be satisfied while the carrier is still non-compliant.
- Discrimination-rule gap (registration 4): this slice has no written rule separating "independently checkable" from "not independently checkable" sub-products (the former may use batch helpers, the latter need a whitelisted carrier); the 3 batch products were checked item by item by the dispatcher and belong to the former, so they were not re-run.
- Three deviations in ② (registered by the user as dispatcher task-package defects): (a) "do not touch those 4 M" contradicted the landing-point list; (b) the anchor sentence lacked terminal punctuation (the ② side added CN `。` / EN `; `); (c) the claim "8 backtick pairs" was actually 5 pairs.
- EN text defect in ③ r1: `docs/t027/REMAINING_SLICES_EN.md` missed the freedom-list literal path; the root cause is the EN wording defect of step 0 (dispatcher).
- Two defects in ⑧a (found and fixed by the dispatcher): first, the `sixstage` fence was not verbatim (the quotes in Sections C / D were turned into straight quotes); second, the EN term was rewritten in lower case in 3 places as `member closed set` (it should be `Closed member set`). **Root cause = when ⑧a copied the material out of the task package it went through a rendering layer (Markdown view / rich-text copy), so quotes were normalised and letter case adjusted**: the package material itself is correct (transcript L86102 carries the U+201C quote `“单件删除”` and already supplied `Closed member set`). **Nature = distortion in the copy step, not a material defect and not ⑧a cutting corners.**
- In its first round ⑧b hit `G4b` reporting `unknown=[U+66F2]` (that character came from my descriptive wording "U+66F2 line quote"); the fix was to reword it as U+201C, **without touching the whitelist** (still 48 entries). This is the same pattern as P30's U+7686 → U+5747 single-character substitution: **the glyphs used in descriptive wording may introduce characters not yet on the whitelist, so they must be checked at write time** (v20 observation candidate).
- ⑧b made two changes beyond the task-package list: first, the block heading marker (forced by G1(d): a 🔄 heading must agree with a ✅ status field); second, the fix-direction row (forced by G6-1: "pending the (1) six-way ruling" and "done" in the same block trigger a tense conflict). **Nature = gate-forced, not scope creep.**
- ⑦ input-package scope: **narrowed** (no full raw diff and no full contract text), declared in 5.3; its PASS proves only that the narrowed items are closed and **is not a whole-slice pass**.
- Temporary-artifact discipline violations: 0 times (as of ⑧b, no temporary artifact was deleted); no self-healing intervention, no fallback re-run.
- Observation self-demonstration (live case in this slice): when writing the root-cause note above, the first wording carried 4 descriptive-use characters (code points U+5BCC / U+61D2 / U+66F2 / U+7686), and `G4b` immediately reported all 4 as `unknown` (none of them on the whitelist); the executor only reworded it (using code points and equivalent phrasing) and did not touch the whitelist (still 48 entries), after which `unknown=[]` followed. **This is a live case, not a hypothetical, of "the glyphs used in descriptive wording may introduce characters not yet on the whitelist".**

## 4 Environment and tiers

- Platform: Windows, PowerShell; repository `d:\LZProjects\prfrail`, Go 1.22 toolchain (no Go code changed in this slice).
- Model tiers (per this slice's actual dispatch): ①⑤ = `DeepSeek V4 Pro (deepseek)`; ②③⑧ = `DeepSeek V4.1 Flash (deepseek)`; ⑥ = `MAI-Code-1.1-Flash (copilot)`; ⑦ = `GPT-5.3-Codex (copilot)`.
- `agentName` arguments: all 4 helper calls in ① passed no `agentName`; ②③⑤⑥⑦ were dispatched on the §2.2 carriers.
- Whitelist: `tools/gates/cjk-newwords.txt` still holds 48 entries; this slice added none.
- Elapsed time: no per-stage timing (not metered).

## 5 Review conclusions (⑤⑥⑦)

### 5.1 ⑤ Pre-review (three rounds)

- r1 (`PRE-REVIEW: PASS WITH FIXES`): 2 Medium findings (M-1 abstracted EN example, M-2 P31 observation declaration without a landing point), each fixed in place.
- r2 (`PRE-REVIEW: PASS WITH FIXES`): both Medium items closed; one new Low ("six items" should read "five items", H unbounded), fixed.
- r3 (quota exception, approved by the user in the same turn): `PRE-REVIEW: PASS WITH FIXES`, High judged closed, 2 Low reported ("five places / one sentence" should read "six places / two sentences", EN term order), fixed.
- Convergence: Medium+ cleared; no unclosed Medium+ entered ⑦ in this slice.

### 5.2 ⑥ Independent scan (two rounds)

- r1 (`INDEPENDENT SCAN: FINDINGS`): 1 High + 2 Medium, details in ch. 3.
- r2 (narrowed re-run): `INDEPENDENT SCAN: PASS`, Section B = `NONE`.
- User ruling: High takes option A and is closed inside this slice (the "Closed member set" sentence); Medium 2 gets no definition in this slice; Medium 3 is registered only, not fixed.

⑥ r2 final message, transcribed verbatim (transcribed by ⑧a; ⑥ is read-only and wrote nothing):

```sixstage
## Section A 扫描能力与方法
本轮仅按工作区现有文档内的文本条款做只读核对，重点覆盖 §9.3 的症状列、§5.1 的 D 计数段，以及 §11.1 的历史编号不追溯约束；不基于任何外部结论，也不声称做过命令级验证。

## Section B 发现
NONE

## Section C 覆盖缺口
本轮已覆盖三处关键口径：第一，§9.3 的症状列与成员定义；第二，§5.1 的 D 计数、跨片累计与历史编号不追溯；第三，同一文档中“单件删除”与“成员闭集”是否互相冲突。未见遗漏的同域条款；本轮结论以文档内同一条文链为准。

## Section D 未钉不变量
未发现新的未钉不变量。关键点是：本行把“成员谁”从“历史报告中的多套计数”收口成“本行括注所列者为闭合枚举”，并把新成员的进入条件写成“先修订本行括注再计入序列”，这使成员资格成为显式、机械可判定的事实，而不是依赖历史报告的口径漂移。该句与 §5.1 的 D 计数段并非双源，而是同一机制的分工：片内 D 计数在本节内部登记，跨片模式累计由 §9.3 症状列负责。

## Section E 结论依据
1. 高是否闭合：是，闭合。证据在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L893-L894)：
   “成员闭集（v1.53）”：症状列的成员以本行括注所列者为闭合枚举；任何新增成员（含符合症状字面的新实例）须先修订本行括注方计入序列；不得以形态自证入列。
   这句话明确了两点：第一，成员集合被本行括注唯一定义；第二，新增成员不能靠“看起来像”或“形态符合”自证，而必须改括注后才进入序列。该表述不含“如 / 例如 / 等”等开放式词，因而构成成员列表而非示例清单，且其判定逻辑可机械执行。

2. 是否引入新问题：否，未引入新的冲突。对应条款在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L513-L515)：
   “D 计数口径：D 是本片临时件纪律违规的独立片内计数……跨片累计：本节不新设跨片计数器；跨片的模式级累计按 §9.3 的『② 提前清理』模式与其症状列执行……历史编号不追溯……”
   这与同行在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L893-L894) 的“单件删除独立登记，不计入跨片模式序”是同一逻辑链条：片内 D 与跨片模式序分开，且跨片累计依赖 §9.3 的症状列与闭合枚举，而非另起一套“谁是成员”的独立来源。故不存在双源，也不与 §5.1 的 D 计数段发生矛盾；两者是互补而非冲突。

结论：高问题已在本片内收口；新句使“成员是谁”成为可机械判定的事实；且未引入 §5.1 / §9.3 之间的冲突或双源。

INDEPENDENT SCAN: PASS
```

### 5.3 ⑦ Final review (one round)

- `FINAL REVIEW: PASS`, Section B = `NONE`; all four dimensions (mechanical decidability / single source without conflict / mirror equivalence / trace compliance) pass.
- ⑦ ruling on Medium 3: registration is sufficient (a new body clause would exceed this slice's boundary; the residual is carried by the independent scan and the ⑩ fact pack cross-check).
- **Input-package scope = narrowed**: no full raw diff and no full contract text; the package carried the point-by-point description of the six landing points, ⑥'s High plus the user ruling (option A), the user ruling basis for Medium 2, the registration text for Medium 3, and the four review dimensions (mechanical decidability / single-source non-conflict / mirror equivalence / trace compliance), with ⑦ reading the workspace state itself.
- **Reading boundary**: that PASS proves only that "the narrowed items are closed" and **is not a whole-slice pass**; whole-slice soundness is jointly carried by ⑤ r3, ⑥ r2 and the ⑨ native validation.
- **Difference from P30**: P30's ⑦ r1 judged PASS on a **complete input package**, while this slice's ⑦ judged PASS on a **narrowed package** (both used 1 of the 1+3 quota); the two "first-pass" results do not carry the same proving power and must not be equated.

## 6 Review input package summary (with blind-review isolation proof)

- Review inputs: the CN / `_EN` directive diff (six places), the CN / `_EN` ledger diff, the correction-note diffs of the two historical reports, and the gate and dispatch readings under `tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/`.
- Not provided as input: the draft of this report and the draft of the freedom list (⑤⑥⑦ all completed before the reports landed on disk).
- Blind-review isolation proof: ⑤ `deep-reasoner`, and ⑥ and ⑦ `independent-reviewer`, are all read-only carriers; the ⑥⑦ models differ from the ② implementation model (`MAI-Code-1.1-Flash` / `GPT-5.3-Codex` vs `DeepSeek V4.1 Flash`).
- Counting caliber: ⑥ checked only the textual clauses inside the workspace documents and claims no command-level verification (see Section A of the fence in 5.2).

## 7 Falsifiability (mechanical checks + mutation tests)

- Falsifiable claim one: the directive body carries exactly six changes. Mechanical check: `git diff --numstat` yields `+21 / -3` for `docs/DELIVERY_DIRECTIVE.md` and for its `_EN`; any extra line refutes it.
- Falsifiable claim two: the `§9.3` symptom cell is untouched. Mechanical check: the hunk of that row in `git diff -U0` contains disposition-cell additions only, and the symptom-cell text is verbatim unchanged.
- Falsifiable claim three: ledger P31 agrees with this slice's landing points. Mechanical check: compare the `docs/t027/REMAINING_SLICES.md:63` status row against the P31 block that starts at `:1049`.
- Falsifiable claim four: the mirrors are shape-identical line by line. Mechanical check: `G3(a)` reports all six shape values as `true` for every checked pair, and `numstat` matches on both sides.
- Mutation tests: this slice changes no criterion, no whitelist and no script under `tools/gates`, so no criterion mutation case was designed (not applicable); the substitute evidence is `--selftest` at `697/697` and `TOTAL_FAIL=0` on three scopes.

## 8 Gate results

- The readings below are a **⑥ scan-time snapshot** (per the `§11.3` time-point caliber, the current values are the ⑩ receipt and the CI readings).

| Item | Reading |
|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` |
| `--all --scope=tree` | `TOTAL_FAIL=0` |
| `--all --scope=ci` | `TOTAL_FAIL=0` |
| `--selftest` | `SELFTEST: PASS 697/697` |
| `G6-1` / `G6-4` | `scanned 96` |
| `G6-6` | `scanned 2` |
| `base go test` | `ok=15 FAIL=0` |
| `G4b` | `newCjk=388 unknown=[]` |
| `G3(a)` | four pairs symmetric in `numstat`; the directive pair at `+21 / -3` |
| Whitelist | `tools/gates/cjk-newwords.txt` 48 entries (unchanged) |
- **Go-side readings of the ⑨ native validation** (`gofmt` / `build` / `vet` / `go test`) come from the native validation after ⑧b; only `.md` files changed afterwards and no Go file was touched, so the `.md` changes are covered by the three gate scopes of ④. **"Not re-running the Go side" rests on the change type, not on saving time** - any future change that touches Go files must be re-run.

## 9 ⑧c cleanup record (§5.1)

- ⑧c to be executed: clean every temporary artifact under `tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/` (file-by-file deletion + read-back assertion); the completion result is in the ⑩ receipt.
- Temporary-artifact D accumulator: 0 times (as of ⑧b, no temporary artifact was deleted).
- ⑧a performed no cleanup, so this section records state only and pre-fills no terminal reading.

## 10 Cost and metering

- Review call counts: ⑤ 3 times (quota 1+1, r3 is the quota exception approved by the user in the same turn); ⑥ 2 times (quota 3, 1 left); ⑦ 1 time (quota 1+3, 3 left); ③ 2 times; ① 4 times (3 of them carrier violations).
- Probe allowance: not used in this slice.
- Timing: no per-stage timing (not metered); this slice writes no estimate at all.
- Monetary cost: not metered.

## 11 Explicitly not executed

- No definition was added for `§11.1` "frozen" (user ruled it out for this slice).
- No body clause was added for ⑥ Medium 3 (user ruled registration only).
- CI was not run (no commit before ⑨); nothing was committed and nothing was pushed.
- No retroactive rewrite of the historical D numbers in `P30` / `P28` / `P23` (`§11.1` freeze; the "historical numbering is not retroactive" sentence registers only).
- No timing (per-stage durations not metered).

## 12 Next-step recommendations

- ⑧c cleanup → ⑨ commit and observe CI (both legs must print the package-level `ok github.com/larsonzh/prfrail/internal/gates` line) → ⑩ receipt (including the "fact pack for v20").
- v20 candidates: the `§11.1` "frozen" definition, `§7.3` data-source adjacency hardening, and an `R2.1` clause for `agentName`.
- For the next slice, a business slice (`[ROUTINE]`) is still recommended, to restore the pilot metrics.

## 13 Registrations (observations / discipline candidates / dispatcher defects)

1. Discipline candidate (dispatch arguments): `R2.1` constrains `model` only, never `agentName`; 3 batch helpers in ① passed neither and 1 architecture helper passed no `agentName` ⇒ all fell back to the session default carrier. Recommendation for v20: write "dispatch must pass `agentName` + `model` explicitly" into `§2.2` or `§4`.
2. Discipline candidate (discrimination rule): distinguish "independently checkable" from "not independently checkable" sub-products; the former may use batch helpers, the latter must use a whitelisted carrier. The 3 batch-helper products here belong to the former, so they were not re-run.
3. Discipline candidate (scan pattern): the cross-repo scan this time wrote `待授权` as one word and missed "待用户授权"; use `待.{0,6}授权` or `授权|authoris|authoriz` instead.
4. Registration (not fixed): `§7.3` "verbatim-locating assertion" does not prove the quoted text is the only / complete data source of the conclusion (⑥ Medium 3; ⑦ ruling "registration is sufficient").
5. Observation: `§11.1` "frozen" has no definition (⑥ Medium 2; user ruled it out for this slice; v20 candidate).
6. Dispatcher defects (3 registered by the user): (a) "do not touch those 4 M" contradicted the landing-point list; (b) the anchor sentence lacked terminal punctuation (the ② side added CN `。` / EN `; `); (c) the claim "8 backtick pairs" was actually 5 pairs.
7. Quota exception: ⑤ r3 exceeded the 1+1 quota (**over by 1**). The user's same-turn dispatch wording (transcript L85739, section five item 4): "4. ⑤ re-review (1 left within quota): check only whether adding the closed-set sentence introduces new problems, do not re-review the whole text" - so that round was **dispatched by the user explicitly**, not decided by the executor; but the user's own label "1 left within quota" conflicts with §7.5's ⑤ 1+1 (=2 rounds), so it is registered as "over quota by 1, explicitly authorised by the user" and must not be cited as a precedent.
8. Known limit: ⑥ Medium 3 was not resolved inside the slice and is handed to v20.
9. Pilot: this slice is still a pure documentation directive-governance slice, so the pilot stays deferred (`OB-124`, time 22).

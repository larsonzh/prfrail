# ProofRail P29 governance-slice validation report: closing the empty-subdirectory gap and landing three collaboration disciplines (directive v1.51)

Slice `DIRECTIVE-P28-FOLLOWUP`; baseline `HEAD` = `92024b3`; directive v1.51; size M (governance slice); report date 2026-10-26.

## 1 Summary

**Problem**: the v1.50 "runtime cache directory exception" only says the members of the exception set may be kept; it does not forbid creating an empty subdirectory inside a member directory. Reading (2) counts only the recursive FILES under `tmp/` and reading (3) looks only at the TOP-LEVEL names of `tmp/`, so neither can see an empty directory inside the exception set - the exception clause is left with a gap. Two P28 leftovers stay open (the symptom column records only 2 same-type deviations, and the disposition column does not cover "if the third occurrence is waived and work continues"); moreover three collaboration disciplines (the inferred authorisation for the (10) stop-point CI-evidence push, the evidence-gap wording, and enumeration before derivation) have no landing place in the body.

**Ruling**: item 4 = (a) - the exception clause gains "a member directory must not contain an empty subdirectory", and the pre-(10) final-state readback gains a (5) reading (assertion: recursive directory count inside the exception set == known member count), with no G criterion; item 5 = (d) - discipline A lands in the §11.2 body, disciplines B and C land in the "note" field of appendix C.0al as a record only.

**Result**: the directive rises to v1.51 (header / §0.3 / §5.1 / §9.3 / §11.2 / appendix B.7.3 / appendix C.0al - seven places) and the P29 blocks of both ledgers follow; all three gate scopes read `TOTAL_FAIL=0` and the self-test reads `SELFTEST: PASS 697/697`; no code was changed, no criterion body was altered, and nothing was pushed to gitee.

**Whitelist +1 (`幻`)**: this entry is not a "legitimate new-word registration" but a "verbatim retention forced by the transcription contract" - the (6) report's Section E used "文档幻觉", and the transcription contract (§11.3) forbids rewriting the reviewer's original wording, so `幻` must be retained and registered on the whitelist. The character itself is a legitimate Chinese word, not a coined one. If a "transcription bringing in a new character" recurs, the same path still applies: do not rewrite the transcription, register the whitelist entry, and state the reason in the report; never rewrite the reviewer's original wording merely to keep the whitelist from growing.

## 2 User ruling

| Item | Ruling | Basis / rationale |
| --- | --- | --- |
| Item 4 (the empty-subdirectory gap) | **(a)** | Add the prohibition sentence plus the (5) reading; a G criterion would require splitting the slice, which the user explicitly declined |
| Item 5 (landing place of the three disciplines) | **(d) split** | Discipline A into the body; disciplines B and C into appendix C.0al |
| Landing place of discipline A | §11.2 `:933` | Added there only; **§10.1 / §1.7.4 are not rewritten** |
| Landing place of disciplines B and C | the "note" field of appendix C.0al | **Recorded only, no normative force**; a clarifying sentence opens the note cell |
| Should (5) become a G criterion | **No** | It stays the (5) reading of the pre-(10) final-state readback; a red verdict depends on manual execution |
| Extend the member count of the exception set | No | the only member today is `tmp/gates-cache/g4b/`; the closed-set sentence requires amending the clause before a new member may join |

## 3 What landed

| # | Landing place | Content |
| --- | --- | --- |
| C1 | header `:3` | `v1.50 / 2026-10-25` → `v1.51 / 2026-10-26` |
| C2 | §5.1 `:510` (one physical line, 5 spots) | the count sentence moves from four readings to five; a (5) reading is added after (4) (fixed command + assertion note); the closed-set sentence gains "and the fixed commands of (2) and (5) are amended in step"; a new "member-directory shape (v1.51)" sentence is added; the "observable form when this sub-clause fails" sentence moves from four readings to five |
| C3 | §11.2 new `:933` | discipline A: the inferred authorisation for the (10) CI-evidence push, carrying the dependency clause "§10.1 unchanged / this scenario only / never a precedent" |
| C4 | appendix B.7.3 `:1368` | template "verbatim four readings" → "verbatim five readings" |
| C5 | §9.3 `:886` two columns | the symptom column gains the third instance `DIRECTIVE-TMP-CACHE-EXCEPTION`; the disposition column gains "after the third occurrence, a waiver to continue must register the evidence gap and the temp-artifact D count" |
| C6 | §0.3 `:97` | a new v1.51 change row whose five-column structure matches the v1.50 row |
| C7 | appendix C `:1987`–`:1999` | new C.0al block (heading + blank + header + separator + eight fields + blank); `### C.1` shifts to `:2000` |
| R1 | §5.1 `:510` | (3)-L2: when the (5) reading is below the known member count, cleanup does not apply - restore the member directory first, then re-run that readback |
| R2 | §5.1 `:510` | (3)-L3: the drift-attribution sentence gains "an empty subdirectory inside the exception set is covered by (5)" |
| R3 | C.0al target `:1991` | (3)-L4: drop the cross-document serial number and state the two items directly |
| F1 | field 16 of the P29 ledger block | a new "known limitations (not fixed here)" registering four retained boundaries |
| F2 | the P29 ledger block | `§9.3:885` → `§9.3:886` inside the block (everything shifts by 1 after the v1.51 row is added) |
| F3 | C.0al target `:1991` | (3)-Low: "the failure-shape sentence" → "the 'observable form when this sub-clause fails' sentence" |
| F4 | C.0al note `:1998` | the note cell opens with "the words 'must' and 'must not' in this section restate the record and create no normative duty" |
| L-1 / L-2 | §9.3 `:874` / §5.1 `:510` | L-1 rejected (a diagnostic command set is not the (10) readback checklist); L-2 is carried by R1 |

L-3 and L-4 are carried by R2 and R3; the identifiers are kept so the original (3) findings stay traceable.

## 4 The five readings and falsifiability

**The fixed command for (5)** (the input measured in the `index` / `tree` scopes):

```text
(Get-ChildItem -LiteralPath tmp -Recurse -Directory -Force | Where-Object { $_.FullName -like '*\tmp\gates-cache\*' } | Measure-Object).Count
```

| State | Reading | Meaning |
| --- | --- | --- |
| Baseline | 1 | only the exception-set member directory `g4b` |
| An empty subdirectory created under `g4b/` | 2 | red: the empty subdirectory enters the recursive directory count |
| That empty subdirectory deleted | 1 | green again |
| `g4b/` moved out entirely (boundary) | 0 | red: the reading is below the known member count, and R1 supplies the remedy |

A same-session measurement from (3): before and after the empty subdirectory exists, (2) stays at 107 and (3) stays at 1, while only (5) moves from 1 to 2 - so (5) covers a surface that (2) and (3) do not, and is not a restatement of either.

| Layer | Criterion | Covers | Does not cover |
| --- | --- | --- | --- |
| L1 | the five readings of the pre-(10) final-state readback (including (5)) | temp-artifact leftovers; the shape and count of the exception-set directories | whether the readings were actually executed (no mechanical detection); member identity (an empty directory at the same level still reads 1) |
| L2 | the numstat symmetry and six shapes of G3(a) | structural parallelism and add/remove symmetry between CN and `_EN` | semantic equivalence; cross-line compensating drift (backticks and pipes are whole-file counts) |
| L3 | G6-5 section-reference resolution and the §0.3-to-ledger agreement of G6-7 | dangling references and declaration-ID existence on added ledger / changelog lines | line-number references inside a method or diagnostic command set (such as `§9.3:874`), which are outside the judged surface |

## 5 Gate and self-test readings

| Scope | Last line | Note |
| --- | --- | --- |
| `--scope=index` | `TOTAL_FAIL=0` | `G5(f)`: index and worktree agree; no untracked path |
| `--scope=tree` | `TOTAL_FAIL=0` | full worktree |
| `--scope=ci` | `TOTAL_FAIL=0` | changed-line scope (pre-run inside the same input package) |
| `--selftest` | `SELFTEST: PASS 697/697` | assertion total unchanged: this slice adds no criterion |

`G1-b`: both directive files read `header=v1.51 last=v1.51 rows=52`. `G4b`: `unknown=[]` (the 「幻」 introduced by the transcription is registered in `tools/gates/cjk-newwords.txt` through the standard process). `git status --porcelain -uall`: the `??` entry = 0. The six shapes of the directive pair (lines / head / bold / bars / blank / backticks) are equal CN=EN with values 2066 / 132 / 5291 / 3478 / 488 / 5706; the same set for this report pair is 200 / 20 / 38 / 252 / 57 / 264 - the line count follows the gate caliber (excluding the trailing empty element) and the other caliber is F-7 of the freedom list.

**CI evidence** (run `40251589108`, commit `2713c69`, conclusion `success`): both legs read `SELFTEST: PASS 697/697` and `TOTAL_FAIL=0`; package-level `ok github.com/larsonzh/prfrail/internal/gates` = `0.296s` (Windows) / `0.113s` (Ubuntu).

## 6 Review chain

| Stage | Carrier | Rounds | Outcome |
| --- | --- | --- | --- |
| (1) Design | `deep-reasoner` + `DeepSeek V4 Pro` | 2 | r1 failed at the dispatch channel with no output; r2 succeeded, producing the four-option comparison table first and the direction design second |
| (2) Implementation | `prfrail-implementer` | 5 | every round ended `IMPLEMENT: DONE`; r1–r2 landed C1–C7 and R1–R3, r3 carried the F1–F4 items from (5) r1, the remaining rounds closed out |
| (3) Test | `prfrail-tester` | 2 | Medium+ = 0; four Low items (L-1 rejected, L-2 / L-3 / L-4 fixed); not re-run after the fixes and handed to (5) under D-2 |
| (4) Gate | coordinator | 3 | `TOTAL_FAIL=0` in all three scopes each time |
| (5) Pre-review | `deep-reasoner` + `DeepSeek V4 Pro` | 1+1 | r1 `PASS WITH FIXES` (one Medium + five Low, three fixed in place and two registered inside F1); r2 re-checked the remediation |
| (6) Independent scan | `MAI-Code-1.1-Flash` | 2 | r1 concluded PASS but its Section D did not walk the guards one by one, so the coordinator judged it non-compliant; r2 re-scanned and is the valid round, transcribed verbatim in chapter 11 |
| (7) Final review | `GPT-5.3-Codex` | 1+2 | r1 blind; r2 `FINDINGS` (1 Medium); r3 re-check `PASS` |
| Batch assistance | `GLM-5.3-Flash` | 3 | mechanical transcription and long-document comparison; it never carries a verdict |

## 7 Known limitations

- First, the (5) reading has no mechanical criterion protecting it - the user ruled that no G criterion is added, so a red verdict occurs only in the human-executed pre-(10) final-state readback (inheriting the P28 ledger item "whether a receipt attached the readings is not mechanically detected"; after v1.51 that caliber is read as "five" rather than "four");
- second, (5) asserts only a count and does not verify member identity, so if `g4b/` is deleted and replaced by an empty directory at the same level, (5) still reads 1 (a false-green window with no mechanical closure - the next cold G4b run recomputes under fail-open and rewrites the cache without a red verdict; it is exposed only when the directory count differs from the known member count or when the cache content is compared by hand);
- third, the "SCM stale display" row of §9.3 states "the four commands such as `git status --porcelain -uall` are authoritative", which is a dangling reference with no definition and is not bound to the five readings of §5.1; this slice rejected broadening it to "five" (rationale: that sentence is a diagnostic command set, not a readback checklist) and registers it for the record without fixing it;
- fourth, the inferred authorisation in §11.2 rests only on textual boundary constraints - no mechanical criterion prevents its extension beyond the ⑩ CI-evidence push scenario, so it must be reviewed by hand at (5) and (7).

The four items above are copied verbatim from the same-named cell of the P29 ledger.

## 8 Exception handling

### 8.1 Temp-artifact D = 1: the step-0 misplaced file of (2) (self-healing tier)

During step 0 the (2) probe wrote its output outside the slice directory (at the top level of `tmp/`) and deleted it immediately; the top level of `tmp/` ends with zero leftovers and the file never entered any commit. Under the disposition column of `§9.3:886` ("a self-misplaced file outside that directory is handled by the self-healing tier of §5.1") the temp-artifact count of this slice is D = 1.

**Stated explicitly**: this is not the fourth same-type deviation of the `§9.3` "early cleanup of (2)" row. That pattern is "cleaning the slice temp directory too early", which never happened here; the symptoms, the disposition and the counting caliber all differ, so the two must never be merged into one entry.

### 8.2 Process deviation: the double BOM of (2) r2 and several `MM` stagings

On its first write r2 stacked two BOMs (`G4a` reported `bom=2`); the extra BOM was stripped and the three gate legs were re-run. Several files were also modified without being staged in time, which turned `G5(f)` red until each was `git add`-ed. Both were self-reported, self-fixed, left zero residue and belong to the routine self-healing tier.

### 8.3 Why the coordinator rejected (3)-L1 (`§9.3:874`)

(3)-L1 proposed changing "the four commands such as `git status --porcelain -uall` are authoritative" to "five". The coordinator rejected it: that sentence belongs to the diagnostic command set of the §9.3 "SCM stale display" row and is not the (10) readback checklist of §5.1 (only reading (1) carries `git status`); changing it would suggest the two are the same thing, and it lies outside the landing place the user fixed for items 1–3 (the `:886` row alone). The (1) design had also ruled "leave it unchanged".

**D-2 (⑧b change beyond the instruction, accepted by the user, never a precedent)**: outside the "change only three places" instruction, ⑧b also corrected the tense of the "remedy direction" row in the P29 ledger (from "待 ① 两组裁定" to "① 已裁定 = 第 4 项 (a) / 第 5 项 (d)"). Reason: `G6-1` reported a tense inconsistency. The four-option enumeration is retained verbatim, with no line added or removed. Nature = one extra line changed to turn the gate green; **never a precedent**. **What to do next time**: when `G6-1` turns red, **stop and report first, then wait for the user's ruling**; do not change it on your own. This slice was the special case "⑧b had already written it, the gate was already red, and stopping cost more than changing"; where that premise does not hold, stopping first is mandatory.

**Whitelist actually +3 (`幻` / `膨` / `胀`)**: `幻` was brought in by the (6) transcription; `膨` and `胀` were brought in by the verbatim text of registration 1 ("白名单不膨胀"). All three are legitimate Chinese words, not coined ones, and are handled through the established path "do not rewrite the verbatim text, register the whitelist entry, and state the reason in the report".

## 9 Input for the next slice

- The **identity blind spot** of the empty-subdirectory check (an empty directory at the same level replacing `g4b/` still reads 1) needs its own slice if it is to be closed mechanically (for example a member-name comparison inside (5)); it is out of scope here.
- `§9.3:874` "the four commands" is a dangling reference with no definition; this slice rejected broadening it to "five" and registered it for the record. Whether it should be bound to §5.1 stays open for a later ruling.
- The inferred authorisation of §11.2 currently rests on textual boundaries plus manual review at (5) and (7); any extension needs a human guard and must not be inferred from this slice.

## 10 Boundaries and things not done

The change surface is limited to the v1.51 landing places of the two directives and the two ledgers, plus this report pair, the freedom list and the whitelist registration - eight files in all. It does not touch the criterion bodies of §2.2 / §2.4 / §5.2 / §7.5 / §6.3, nor `tools/gates/**` (the whitelist text aside), `internal/**`, CONTRACTS / schema / fixtures / workflows or `DEV_PLAN`. No ADR was added, no self-test criterion was added, and nothing was pushed to gitee.

## 11 Verbatim transcription of the (6) report

The following is the complete output of the valid (6) round (the r2 re-scan), **transcribed verbatim** (nothing was rewritten; markdown link syntax was reduced to plain paths under the transcription convention the original report states in its own header). The two verdict lines inside the fence (the first line under Section A and the block-final line) carry equal values at compliant positions, matching the `G7-f` contract of `§6.3`.

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

无

## Section C

- 锚点：① `git status --porcelain -uall` 的 `??` 项 = 0、第二列非空项 = 0，且第一列非空项的行集合与本片声明的变更集逐一相同 | 位置：docs/DELIVERY_DIRECTIVE.md:510 | 作用：锁定「五条读数」中的 ① 条，保证变更集与真实工作树一致，避免「误判已清理」或「遗漏变更」 | 依据：03-baseline-notes.md §3；docs/DELIVERY_DIRECTIVE.md:510
- 锚点：② `(Get-ChildItem -LiteralPath tmp -Recurse -File -Force | Measure-Object).Count` = 1 | 位置：docs/DELIVERY_DIRECTIVE.md:510 | 作用：限定 `tmp/` 仅保留 `.gitkeep` 或合法例外，防止临时件残留 | 依据：03-baseline-notes.md §4；docs/DELIVERY_DIRECTIVE.md:510
- 锚点：③ `(Get-ChildItem -LiteralPath tmp -Force).Name -join ','` 恒等于 `.gitkeep` | 位置：docs/DELIVERY_DIRECTIVE.md:510 | 作用：强制 `tmp/` 顶层仅有 `.gitkeep`，抑制「目录顶层误放件」隐性漏网 | 依据：03-baseline-notes.md §4；docs/DELIVERY_DIRECTIVE.md:510
- 锚点：固定命令 `(Get-ChildItem -LiteralPath tmp -Recurse -Directory -Force | Where-Object { $_.FullName -like '*\tmp\gates-cache\*' } | Measure-Object).Count`，基线 = 1 | 位置：03-baseline-notes.md:30-35 | 作用：定义 ⑤ 的固定命令与「=1」判据，确保例外集只允许 `g4b` 这个已知成员 | 依据：03-baseline-notes.md §4；docs/DELIVERY_DIRECTIVE.md:510
- 锚点：成员目录形态（v1.51）——成员目录内不得含空子目录——空子目录计入 ⑤ 的递归目录计数，> 已知成员数即判红 | 位置：docs/DELIVERY_DIRECTIVE.md:510 | 作用：把「空子目录」显式纳入判据，避免「`g4b` 由空目录替代但仍被误判绿」的假绿窗 | 依据：03-baseline-notes.md §1；docs/DELIVERY_DIRECTIVE.md:510
- 锚点：任一不满足 ⇒ 清理并重跑；⑤ 读数小于已知成员数时不适用清理——须先恢复该成员目录（跑一次 G4b 生成缓存，或按 C.0aj 口径重建）再重跑该回读 | 位置：03-baseline-notes.md:20-23 | 作用：给 R1 补救路径定出「先恢复再重跑」，避免错误清理导致不可自愈循环 | 依据：03-baseline-notes.md §2；docs/DELIVERY_DIRECTIVE.md:510
- 锚点：⑩ 停点回执所需的 CI 取证 push——主提交 push 并 CI 双绿后，把 commit hash 与 CI run 号回写报告……该 push 按推定授权执行，不必为 push 再单独索取授权；用户明写「不 push」或「push 另行授权」时以用户措辞为准 | 位置：docs/DELIVERY_DIRECTIVE.md:933 | 作用：把纪律 A 收口到「§10.1 不变 / 只覆盖 ⑩ CI 取证 push / 不得作先例」的三要素，从属关系明确 | 依据：docs/DELIVERY_DIRECTIVE.md:933
- 锚点：仅为留痕，无正文效力：本节内「须」「不得」等词系对留痕内容的转述，不构成正文义务 | 位置：docs/DELIVERY_DIRECTIVE.md:1998 | 作用：澄清 B/C 仅作留痕，不是正文义务，限制「被误用为正式规则」的误读 | 依据：03b-r3-addendum.md §1；docs/DELIVERY_DIRECTIVE.md:1998

## Section D

1. 「五条读数」取数句：会红 = §5.1 的五条读数整条判据，尤其是 ①–⑤ 读数本身；不会红 = G1-b、G3(a)、G4a、G4b、G6 系列机械判据，因为它们是并列的不同检查面；依据：docs/DELIVERY_DIRECTIVE.md:510；03-baseline-notes.md §3-§4
2. ⑤ 的固定命令 = 1：会红 = ⑤ 的固定命令读数；不会红 = ①–④ 的四条读数本身仍可保持绿；依据：03-baseline-notes.md §4；docs/DELIVERY_DIRECTIVE.md:510
3. 「成员目录形态（v1.51）」句：会红 = ⑤ 的递归目录计数在「空子目录」场景中会红；不会红 = 若只删解释句而不删 ⑤ 的命令，本身仍然会因已知成员数不符而红；依据：docs/DELIVERY_DIRECTIVE.md:510；03-baseline-notes.md §1
4. 成员闭集句里的「同步修订 ② 与 ⑤ 的固定命令」：会红 = 当新增例外成员时，闭集维护失效会被审查判红；不会红 = 当前已知只有 `g4b`，且 ②/⑤ 现状仍为 1，所以本轮静态状态不红；依据：docs/DELIVERY_DIRECTIVE.md:510；03-baseline-notes.md §1
5. R1 的「⑤ 小于成员数时不适用清理」补救句：会红 = R1 的补救判据；不会红 = 正常状态下 ⑤ = 1 或成员目录正常恢复后再跑回读时仍绿色；依据：03-baseline-notes.md §2；docs/DELIVERY_DIRECTIVE.md:510
6. §11.2 纪律 A 的三要素从属句：会红 = 回写核对 / CI 取证 push 的「同轮授权」资格判据；不会红 = 代码门禁、`index/tree/ci`、`selftest` 读数本身不因该句被删掉而立刻变红；依据：docs/DELIVERY_DIRECTIVE.md:933
7. C.0al 的「仅为留痕，无正文效力」澄清句：会红 = B/C 规范性措辞被当成正文义务的审查项；不会红 = 机械门禁仍不红；依据：docs/DELIVERY_DIRECTIVE.md:1998；03b-r3-addendum.md §1

没有覆盖的不变量：无。
理由：本轮所列 7 个守卫，均已在 §5.1 的读数、例外集与成员形态、R1 补救、§11.2 的授权从属句、C.0al 的留痕澄清中形成「有机械命令 / 有文本澄清 / 有补救路径」的闭环；输入包里没有再暴露出一条「尚未被命令或语句显式覆盖」的关键不变量。对「成员身份盲区」这一点，基线材料已写成「已知限制（本片不修）」，它不是「未覆盖」，而是「被明确标注为未修且可证伪」的保留边界。

## Section E

我无法验证的部分：
- ⑤ 的「成员身份校验」仍然只算计数，不校验成员身份；`g4b` 若被同层空目录替换，`=1` 仍可能假绿。这个边界已被登记为「已知限制（本片不修）」，所以我不把它当成本轮新增缺陷，但它确实是本片的保留风险。
- 我没有把 `grep-evidence.txt` 当作独立强证，因为该文件在本轮只读输入包里未出现，且基线材料已给出足够的机械判据与逐字锚点；若它在下一轮基础读数中出现并显示「超长行命中次数 / 偏移」发生变化，那么本轮结论需要重核。

什么会推翻我的 PASS：
- ⑤ 的固定命令在下一轮真实回读不再返回 1，或者例外集中出现任何新增成员目录且未同步修订 ②/⑤ 的固定命令；
- §5.1 的五条读数被删减成四条，或 ①–⑤ 读数中任意一条失联、失真或语义被收缩；
- §11.2 的纪律 A 被改为「无三要素从属」，从而移除「§10.1 不变 / 只覆盖 ⑩ CI 取证 push / 不得作先例」；
- C.0al 的「仅为留痕，无正文效力」被删除或语义倒置，导致 B/C 规范性措辞被不当当成正文义务；
- 任何一项「可核对锚点」被证明并非仓库现状中的真实句子，而是文档幻觉。

INDEPENDENT SCAN: PASS
```

## 12 Declared change set

The `artifacts` fence is this slice's declared change set (eight rows, `TAB`-separated); `G5(d)` uses the rows whose disposition is `repo` as `declaredFiles`.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-p28-followup.md	repo
docs/validation/directive-p28-followup_EN.md	repo
docs/validation/evidence/DIRECTIVE-P28-FOLLOWUP-freedom-list.md	repo
tools/gates/cjk-newwords.txt	repo
```

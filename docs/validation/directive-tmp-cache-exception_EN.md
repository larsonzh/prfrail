# ProofRail P28 governance slice validation report: runtime cache directory exception (directive v1.50)

Slice `DIRECTIVE-TMP-CACHE-EXCEPTION`; base `HEAD` = `2c9b6a3`; directive version v1.50; size M (governance slice); report date 2026-10-26.

## 1 Summary

**Problem**: after P27 added a cross-process cache to G4b, the cache directory lands in `tmp/gates-cache/g4b/` (runtime state that must be kept); yet the pre-(10) final-state read-back of §5.1 requires the recursive file count under `tmp/` to be 1 and the top-level names to be identically `.gitkeep`. The two clauses **cannot hold at the same time**, so this is an internal rule conflict.

**Ruling**: the user chose **option (ii)** — change no code (do not move the cache directory out of `tmp/`) and add a runtime cache directory exception clause to §5.1 to close the conflict. P27's option D (the same-round temporary authorization exception) is therefore no longer needed.

**Result**: the directive is raised to v1.50 (five places: header / §0.3 / §5.1 / §12.4 / appendix C.0ak) and the P28 ledger block is synchronised; all three gate scopes report `TOTAL_FAIL=0` and the self-test reports `SELFTEST: PASS 697/697`; no code changed, no criterion body changed, no CI-decidable signal added.

## 2 The (1) fix-direction ruling

| Item | Content |
| --- | --- |
| Candidate (i) | Move the cache directory out of `tmp/`: requires changing `g4b.js`, `.gitignore`, `selftest.js` and touches several P27 frozen-evidence points |
| Candidate (ii) | Add a runtime cache directory exception clause to §5.1: directive only, no code |
| User ruling | **Option (ii)** (2026-10-25); P27 option D is void |
| User sub-rulings | (1) P27 frozen evidence must not be modified; (2) the failure form is preserved; (3) the exception set must be literally enumerated; (4) no over-broad reading |

Reason for not adopting (i): it would require changing `tools/gates/**` and the frozen P27 evidence surface at the same time, which exceeds what this slice needs; and (ii) settles the conflict in place as a standing rule, so the same conflict does not need repeated temporary authorisations in later slices.

## 3 What v1.50 changed

| File | Location | Content |
| --- | --- | --- |
| Both directives | Header | `Version: v1.50. Date: 2026-10-25.` |
| Both directives | §0.3 `:96` | v1.50 row: the four requirements, option D no longer needed, new appendix C.0ak, boundary statement |
| Both directives | §5.1 `:509` | Literal exception-set enumeration, the two fixed exclusion commands for (2) and (3), preserved failure form, closed member set, no over-broad reading, consistency with `.gitignore` |
| Both directives | §12.4 `:1114` | OB-121 (pilot deferral round 19, status `Ruling`) |
| Both directives | Appendix C `:1970`-`:1982` | The C.0ak slice block |
| Both ledgers | Status row and P28 block | Status row synchronised; the 18 fields of the P28 block moved from planned to landed |

## 4 The four requirements and their executability

| # | Requirement | Executable judgement | Mechanical anchor |
| --- | --- | --- | --- |
| 1 | The exception set is literally enumerable | Only current member = `tmp/gates-cache/g4b/`; a new member must first amend this clause, amend the fixed command of (2) accordingly, and be registered | Human read of §5.1 plus the command literals |
| 2 | The failure form is preserved | Any top-level entry under `tmp/` other than `.gitkeep` and outside the top-level container that holds the exception set makes (2) and (3) fail | The fixed command of (3) plus the (10) final-state read-back |
| 3 | No over-broad reading | The exception does not cover temporary scripts or logs at the root level of `tmp/gates-cache/` | The fixed command of (2) excludes only the `g4b` subdirectory |
| 4 | Consistent with `.gitignore` | The exception set must be a subset of, not equal to, the set ignored by `.gitignore` | Human check plus `git check-ignore`; drift is found by G5(f) and by the (2) and (3) readings respectively |

## 5 Gate and self-test readings

| Scope | Last line | Note |
| --- | --- | --- |
| `--scope=index` | `TOTAL_FAIL=0` | `G5(f)`: index equals worktree, no untracked path |
| `--scope=tree` | `TOTAL_FAIL=0` | Full worktree |
| `--scope=ci` | `TOTAL_FAIL=0` | Changed-line scope |
| `--selftest` | `SELFTEST: PASS 697/697` | Assertion total unchanged: this slice adds no criterion |

All six shapes of `G3(a)` (lines / head / bold / bars / blank / backticks) are equal between CN and EN; `G4b` reports `unknown=[]` — this slice did not reword anything to dodge a criterion; the whitelist entries added are recorded in 8.5.

## 6 Review chain

| Stage | Carrier | Rounds | Outcome |
| --- | --- | --- | --- |
| (1) Design | `deep-reasoner` + `DeepSeek V4 Pro` | 1 | `DESIGN: AWAITING_RULING`; the user then ruled option (ii) |
| (2) Implement | `prfrail-implementer` | 5 | `IMPLEMENT: DONE` each round |
| (3) Test | `prfrail-tester` | 2 | r1 one Medium plus three Low; r2 reduced Medium+ to zero and fixed the two Low |
| (4) Gates | Main controller | 3 | `TOTAL_FAIL=0` in all three scopes each time |
| (5) Pre-review | `deep-reasoner` + `DeepSeek V4 Pro` | 2 | r1 `PASS WITH FIXES` (four Medium plus two Low); r2 `FINDINGS` (five closed, one new Medium, one High escalation) |
| (6) Independent scan | `MAI-Code-1.1-Flash` | 1 | `INDEPENDENT SCAN: PASS`, Section B has no findings |
| (7) Final review | `GPT-5.3-Codex` | 2 | r1 blind `REVIEW: PASS`; r2 `RE-REVIEW: PASS` |

The following is the complete output of (6), **transcribed verbatim** (no rewriting); the two verdict lines inside the fence (the first line under Section A and the block end) have equal values and compliant positions, matching the `G7-f` contract of `§6.3`.

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

无

## Section C

1. 四要件的可执行性：从输入包 `six/00-contract-verbatim.md` 的 §5.1 原句可见，例外集已按字面枚举为「`tmp/gates-cache/<criterion-id>/`，当前唯一成员 = `tmp/gates-cache/g4b/`」；「失效形态保留」写明为「若 `tmp/` 顶层除 `.gitkeep` 与例外集所在顶层容器之外还出现任何成员，则 ② ③ 仍判红」；「不得扩大解释」直接写成「该例外只覆盖 `tmp/gates-cache/<criterion-id>/` 形态，不覆盖 `tmp/gates-cache/` 根层的临时脚本或日志」。这三条都能被无记忆第三人按字面读数判定，且与 `.gitignore` 的 `/tmp/*` + `!/tmp/.gitkeep` 口径一致：例外集属于被忽略集合的子集，不等于被忽略全集，故不构成冲突。参见 `docs/DELIVERY_DIRECTIVE.md:509` 与 `docs/DELIVERY_DIRECTIVE_EN.md:509`。

2. 失效形态句的可满足性：同段落中，v1.50 的「排除例外集所在的顶层容器（当前 = `gates-cache`）后仍须恒等于 `.gitkeep`」与旧句的收窄括注并非相互割裂，而是同一条规则的层级分离：一层判定「第二级例外目录」，另一层判定「顶层容器」。这正是 ③ 里「读的是 `tmp/` 顶层名 `gates-cache`，而例外集是 `tmp/gates-cache/g4b/`」的关键差异；该层级差已在同一段落中被显式写成「除 `.gitkeep` 与例外集所在顶层容器之外还出现任何成员，则 ② ③ 仍判红」，没有形成「按哪句判」的歧义。参见输入包 `six/00-contract-verbatim.md` 的 §5.1 段和 `docs/DELIVERY_DIRECTIVE.md:509`。

3. 固定排除命令的正确性：据输入包 `six/03-review-rounds.md` §5，修复前后三组实测数值是自洽的：旧宽通配 = 28，收窄后 = 37，去掉 `gates-cache` 根层 9 件历史残留后 = 10；这对应「`g4b/` 16 + 根层 9 + `.gitkeep` 1」的几何分摊，且最终 10 = 1 + 9 的关系与「`.gitkeep` + 根层 9 件历史残留」一致。也就是说，固定命令的豁免面已经收敛到 `tmp/gates-cache/g4b/`，而不是整个 `tmp/gates-cache/`。本段与「不得扩大解释」句一致，未留出「字面命令 ≠ 条款语义」的实质缝隙。参见输入包 `six/03-review-rounds.md` §5 及 `docs/DELIVERY_DIRECTIVE.md:509`。

4. 偏离与证据缺口的充分性：从输入包 `six/04-deviation-and-waiver.md` 可以看到，四要素登记、证据缺口清单、停片后的特批、以及「不得作先例」都被如实写明；特别地，文件明确标注「`02-gate-*.txt` 系停片后重跑所得，不是停片前原件」，并且在 ⑤ 特批中写明「⑥ 必须审查该修复」。这足够支撑本片按用户裁定继续，且没有把「证据缺口」伪装为「无影响」。参见输入包 `six/04-deviation-and-waiver.md` §3–§6。

5. §0.3 ↔ 正文 ↔ C.0ak 的一致性：三处口径一致：① v1.50 的变更声明写明「四要件：例外集可枚举成字面、失效形态保留、不得扩大解释、与 `.gitignore` 一致」；② 正文 §5.1 与 §0.3 一致保留这一口径；③ C.0ak 中的目标和附注也写成同样的语义。D 选项的口径也已归一为「同轮临时授权例外，依据用户回复『选 D』」，而非原先在某个附注里写成「⑩ 前清空缓存目录」。这部分由输入包 `six/00-contract-verbatim.md` 中的 §0.3 行、正文行和 C.0ak 块共同对齐，且 `docs/DELIVERY_DIRECTIVE.md:96` 与 `docs/DELIVERY_DIRECTIVE.md:1970-1982` 互相一致。

6. 是否开了不该开的口子：没有。原因是输入包 `six/00-contract-verbatim.md` 里的「成员闭集：例外集成员以本句字面枚举为准；新增成员须先修订本条款、同步修订 ② 的固定命令并登记，不得以形态自证入集」以及「不得扩大解释：该例外只覆盖 `tmp/gates-cache/<criterion-id>/` 形态，不覆盖 `tmp/gates-cache/` 根层的临时脚本或日志」把豁免口径锁死在一个闭集和一个明确的层级范围之内，不会被泛化成「凡运行时状态都可从 §5.1 豁免」。这也是本片的核心闭合条件。

此外，门禁与自测来自静态材料：据输入包 `six/02-gate-index.txt`、`six/02-gate-tree.txt`、`six/02-gate-ci.txt` 的末行结论均为 `TOTAL_FAIL=0`；`six/02-gate-selftest.txt` 中给出 `SELFTEST 697/697` 的最终结论；这些都与本轮审计结论一致。

## Section D

- 若删掉「成员闭集」守卫：新增例外成员不会被字面枚举锁定，命令层面的 `-notlike '*\tmp\gates-cache\g4b\*'` 仍可能漂移；这会让「当前唯一成员 = `g4b/`」与实际成员集合不再对齐，最先红的是「成员闭集」相关判定，且「不得扩大解释」也会失去机械约束。
- 若删掉「不得扩大解释」守卫：根层残留脚本/日志可以被误判为合法例外；这不会让 `git status` 的 `??` 读数消失，但会让 `tmp/` 顶层的根层残留在 ② ③ 判定里被豁免，最终会出现「命令与语义不一致」的判定漏洞。
- 若删掉「与 `.gitignore` 一致」守卫：例外集可能不再满足「⊆ 被忽略集」，从而使「例外必须是 `.gitignore` 实际忽略的路径」失效；这不会自动让 `git status` 变绿，但会让例外不再可证明地保留在规则语义下。
- 没有覆盖的不变量：本轮审计没有发现真正未覆盖的不变量；唯一不能由静态输入包直接验证的是未来新成员切片是否会在没有同步修订字面枚举和命令的情况下落入例外集，但这是「后续维护」风险，而非当前条款的口径缺陷。

## Section E

- 我无法验证的部分：现场真实文件系统中 `tmp/` 的运行时状态（尤其是未落盘在本 input package 中的后续临时文件）本轮无法直接重跑；本次审计只能基于静态材料和规则文本进行法定性判断。
- 什么会推翻我的 PASS：如果后续新增成员被加入 `tmp/gates-cache/`，但没有同步修订「成员闭集」、固定命令与相应说明；或如果证据包中显示本例外集实际上不是 `.gitignore` 的子集；或者如果后续切片将「运行时状态」扩张为「凡 `tmp/` 下运行时状态都可豁免」，则本条 PASS 会被推翻。
- 这也是本轮审计的保守判断：当前材料足以证明这一条款在字面与机械判定层面是闭合集、可枚举、且被 `.gitignore` 约束，故不支持再开一条「泛运行时状态豁免」的口子。

INDEPENDENT SCAN: PASS
```

## 7 Falsifiability and known limitations

The claim of this slice is that the exception clause is executable, enumerable and not over-broad; its falsifiability has three layers. The "Not covered" column lists gaps that this slice **proactively registers**.

| Layer | Mechanism | Covered | Not covered |
| --- | --- | --- | --- |
| L1 | The four readings of the (10) final-state read-back, including the fixed exclusion commands | A new top-level entry outside the exception set, or a loose file at the root of `tmp/gates-cache/`, fails | Whether the receipt really attached the four readings has no mechanical detection |
| L2 | `G3(a)` numstat symmetry and six-shape equality | A one-sided CN/EN change fails | Whether the two sentences are semantically equivalent (equal shapes do not imply equal meaning) |
| L3 | `G6-5` section-reference resolution and `G6-7` declaration/ledger reconciliation | A dangling reference fails | Whether the §0.3 declaration matches the body **in content** is borne by human review only |

The (7) final review identified two invariants that remain mechanically uncovered; per the condition attached to `RE-REVIEW: PASS` this slice registers them verbatim below:

1. The (2) and (3) readings **do not mechanically cover** an **empty subdirectory** under `tmp/gates-cache/` that is undeclared and contains no file: the (2) count over `-File -Recurse` does not trigger and the (3) top-level-name reading does not trigger either (the only top-level name is `gates-cache`, already excluded).
2. Consistency with `.gitignore` is still a **human check** and is not an independent mechanical assertion among the four (10) readings.

Four registered known limitations (the two above plus two from (5)): the empty-subdirectory surface (item 1 above); the non-mechanical nature of the `.gitignore` consistency check (item 2 above); §6.3 does not pin the G4b cache path, so if OB-120 adopts fix (i) the literal enumeration in §5.1 must be amended with it (drift is fail-safe); and whether a receipt attached the readings is not mechanically detected.

## 8 Exception handling

### 8.1 Deviation: the third same-type occurrence of §9.3 (user ruled to continue)

**Event**: while closing round 4 of "fix the four Medium findings from (5) r1", the (2) implementer **deleted the whole `tmp/DIRECTIVE-TMP-CACHE-EXCEPTION/`** and self-reported that the directory no longer existed. At that moment the directory held, besides the implementer's own edit script and gate logs, the (5) r1 evidence scripts (three files, including `r5-verify.js`) and the main controller's (4) gate logs (four files).

**Classification**: the symptom column of `§9.3:885` had recorded two earlier same-type occurrences (`G9-c` once, `GATES-PATCH` once), so this is the **third**; `§5.1:505` states explicitly that the self-heal tier does not apply to this pattern. The main controller therefore **stopped and escalated**: no (6) or (7) dispatch, no report, no commit, no push.

| §1.7.4 element | Content |
| --- | --- |
| Adjudicator | User |
| Date | 2026-10-26 |
| Reason | While fixing the four Medium findings from (5) r1, the (2) implementer deleted the whole slice temp directory, making evidence scripts unrecoverable and losing gate logs, which is the third same-type occurrence; after the mandatory stop-and-escalate, the user ruled to continue |
| Basis | The user's reply, verbatim: "choose B (authorise continuation)" |

**Nature**: this ruling **exceeds** the literal consequence of `§9.3:885` ("the third same-type occurrence leads to stop-and-escalate"), which would have ended the slice, and is a direct user adjudication. **Scope**: continuation of (6) (7) (8) (9) (10) for this slice only. **Not a precedent**: a fourth same-type occurrence must still stop and escalate. **Counted as**: one increment of this slice's temp-file D count.

### 8.2 Evidence gaps

What was lost is the **evidence carrier**, not the **conclusion**: the three (5) r1 evidence scripts are **unrecoverable** and the four (4) gate logs are gone; however, the six (5) r1 findings were re-checked item by item in (5) r2, and the (4) three-scope readings and self-test count were captured in the terminal and restated in chapter 5 of this report. **Residual risk**: the gate materials in this report and in the input package come from a **post-stop re-run**, not from the original process of the pre-stop run; where a strict "same run" reading is required, this slice's gate evidence should be regarded as a post-stop reproduction.

### 8.3 Attached waiver: the (5) Medium fix is deemed closed without re-running (5)

**Fact**: (5) r2 reported one new Medium (the fixed command's exemption surface was wider than the clause's semantic exemption surface); the (2) implementer fixed it (the `-notlike` wildcard was narrowed from container level to the `g4b` subdirectory level) but **did not re-run (5)**, because the (5) quota of 1+1 was exhausted.

**Waiver**: the user ruled it closed, on the grounds that "the fix is mechanically verifiable (after narrowing, the measured value is `1 + 9`, with `g4b/` still exempt and the root-level loose files newly caught)", that "the quota was exhausted", and that "(6) will independently review the fix". **Nature**: this **exceeds** the literal requirement of `§5.3` ("after (5) reports Medium+, re-run (5)") and is a direct user adjudication; **not a precedent**; **(6) did in fact review the fix** (see chapter 6 and item 3 of `six/06-mai-report.md`).

The (7) re-review judged this waiver **acceptable but not strictly equivalent to a full (5) re-run**: the substitute coverage only verifies the fix of that one Medium and is narrower than a full (5) review. This is registered exactly as (7) judged it, without any "equivalent" claim.

### 8.4 Two citation and wording defects in (6) (treated as leads)

First, (6) cited `docs/.gitignore`, which does not exist; the correct location is the repository-root `.gitignore` (`:11` = `/tmp/*`, `:12` = `!/tmp/.gitkeep`). Second, item 3 of (6)'s Section C described the arithmetic as "after removing the root-level loose files of `gates-cache`, the value is 10", whereas it should be "37 minus 27 (the slice temp directory) = 10 = 1 + 9". Both are corrected under the `§7.3` rule that an unlocatable citation is treated as a lead and not accepted as-is; **neither changes the Section A verdict** nor constitutes a clause defect.

### 8.5 Whitelist registrations and the historical leftovers at the root of `tmp/gates-cache/`

`G4b` reports `unknown=[]`; following the standard flow this slice registered whitelist entries for newly introduced legitimate CJK characters (entries and reasons in freedom-list item F-4) and **did not reword anything to dodge the criterion**.

The root of `tmp/gates-cache/` holds **9** historical leftovers from the P27 and README slices (not members of the exception set). The user ruled that the **main controller cleans them item by item at ⑧c** with a read-back proof, and that the (2) implementer must not do it (which would trigger the same-type deviation again). After cleanup the (2) reading should be 1, with only `g4b/` exempt. **These are omissions of the P27 and README slices, not introduced by this slice.**

## 9 Input for the next slice (v1.51)

1. Append this slice to the symptom column of `§9.3:885` as the third recorded occurrence.
2. A new OB, numbered by mechanically checking the maximum in §12.4 plus one (expected OB-122).
3. The handling column of `§9.3:885` may add: "if the user authorises continuation after the third occurrence, the report must register the evidence gaps and the D count".
4. The empty-subdirectory coverage gap noted by (7) may be taken up as a mechanical criterion in `§5.1` or `§6.3`.

## 10 Boundaries and things not done

This slice changes none of `internal/**`, CONTRACTS / schema / fixtures / workflows, `g4b.js`, `.gitignore`, `selftest.js`, `DEV_PLAN*`; it changes no criterion body in §2.2 / §2.4 / §5.2 / §7.5 and §6.3; it does not address OB-120; and it pushed nothing to gitee. One further note: the change set includes a **whitelist append** to `tools/gates/cjk-newwords.txt` (registering two newly introduced legitimate CJK characters), which is a standard action of the `G4b` fail-closed flow and **changes no criterion logic**. Therefore this slice does not constitute a semantic change to criteria, the evidence model, or the role and authorisation contracts.

## 11 Declared change set

The `artifacts` fence is this slice's declared change set; `G5(d)` uses it as the source of `declaredFiles`.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-tmp-cache-exception.md	repo
docs/validation/directive-tmp-cache-exception_EN.md	repo
docs/validation/evidence/DIRECTIVE-TMP-CACHE-EXCEPTION-freedom-list.md	repo
tools/gates/cjk-newwords.txt	repo
```

# DIRECTIVE-WHITELIST-POLICY slice validation report (collaboration discipline D promoted to a formal clause in the directive body)

> Slice `DIRECTIVE-WHITELIST-POLICY` (governance slice) | tier `[SLICE]` | size S | date 2026-10-26.
> Status: ✅ **done** (2026-10-26; directive v1.52; the ⑦ chain closed at 1 final review + 1 re-review under §7.8.3; `TOTAL_FAIL=0` in all three scopes and `SELFTEST: PASS 697/697`; commit `e6d14a2`; see section 8 for the CI readings).
> [中文](directive-whitelist-policy.md)

## 1 Change summary (including the executive-freedom list)

**What this slice does**: promotes collaboration discipline D (the whitelist policy outranks the authorised quantity) from a session-level collaboration discipline to a formal clause in the directive body, landing at the user-ruled candidate (b) - a new paragraph at the end of `§1.7.2` (three lines, no new sub-clause number); in step it adds the v1.52 change row to `§0.3`, registers OB-123 in `§12.4` (the twenty-first deferral of the pilot), adds appendix C block C.0am, raises the header to v1.52 and mirrors all of it strictly in the `_EN` file; the ledger gains the three P30 status places.

**What this slice does not do**: it leaves `§1.7.4` / the `§2.2` carrier table body / `§2.4` / the `§5.2` matrix / the `§7.5` figures / the `§6.3` criterion bodies unchanged; it does not touch `internal/**`, CONTRACTS / schema / fixtures / workflows or `tools/gates/**`; it does not change `DEV_PLAN{,_EN}`; it pushes nothing to gitee.

**Artifact list**: identical row by row to the `artifacts` fence at the end (the four directive and ledger edits, this report pair, and the freedom list).

**Executive-freedom list** (`docs/validation/evidence/DIRECTIVE-WHITELIST-POLICY-freedom-list.md`; F-1 to F-8 are freedoms, D-1 to D-8 are user rulings):
- **F-1**: split the single user-ruled paragraph into three physical lines, breaking only at sentence boundaries, without adding or removing any character or punctuation.
- **F-2**: the `_EN` mirror uses the terms `authorized quantity` / `dispatch / task-package text` / `re-review`.
- **F-3**: the v1.52 change row in `§0.3` follows the four-part shape of the v1.49 to v1.51 rows and declares both OB-123 and C.0am.
- **F-4**: the OB-123 row is appended at the tail of the `§12.4` table and its disposition cell opens with `已裁决` (a closed-set word).
- **F-5**: C.0am reuses the C.0al field skeleton, and its note is a pointer only - it does not restate the discipline D body.
- **F-6**: the P30 status symbol in the ledger is `🔄` (not `⬜`); both belong to the closed set.
- **F-7**: the 6 conclusion line is transcribed under the two-place convention G7-f allows, and both places must carry the same value.
- **F-8**: the observation-A landing declaration file is named `observation-A-landing-declaration.md`.

**One-line conclusion**: discipline D has landed in the `§1.7.2` body under the user-ruled candidate (b) and is in force; all gate scopes read `TOTAL_FAIL=0` and `SELFTEST: PASS 697/697`; the 7 verdict `BLIND REVIEW: PASS WITH FIXES` was re-checked by `RE-REVIEW: PASS`, and all four Low items are out-of-slice suggestions (criterion domain or locked verbatim text), so this slice changed no text for them.

## 2 Execution pipeline

| Ring | Role | Carrier | Model (qualified name) | Status | Rounds | Artifact |
|---|---|---|---|---|---|---|
| 1 | architect | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | executed | 1 | five-dimension comparison of the four landing candidates + draft text per candidate + split ruling (conclusion line `ARCHITECT: OPTIONS READY`, no preference expressed) |
| ruling | user | — | — | ruled | 1 | chose (b); also ruled that the authorisation text = the main controller dispatch text |
| 2 | implementer | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | executed | 1 | five edits × CN/EN (conclusion line `IMPL: LANDED WITH NOTES`, two notes) |
| 3 | test engineer | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | executed | 1 | eight consistency checks (conclusion line `TEST: CONSISTENT`, zero deviation, two observations) |
| 4 | main controller | — | — | executed | 2 | three gate scopes + SELFTEST (re-run once after the 5 r1 fixes) |
| 5 | pre-reviewer | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | executed | 2 (r1 `PRE-REVIEW: FINDINGS` / r2 `RE-REVIEW: PASS`) | one Medium (A-1 7 quota) + three Low + semantic risks D-1 to D-6 |
| 6 | independent scanner | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | executed | 1 (`INDEPENDENT SCAN: PASS`, zero findings) | Section A to Section E |
| 7 | independent final reviewer | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | executed | 2 (r1 `BLIND REVIEW: PASS WITH FIXES` / r2 `RE-REVIEW: PASS`) | four sections + four Low findings (all out of slice) |
| 8a | documenter | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | this ring | 1 | this report + the ledger write-back |

- **Rework count**: 5 once (the original Medium fix); 7 none (nothing at Medium or above, and all four Low items were judged out of slice).
- **How the 7 chain closed**: under the normal-convergence rule of `§7.8.3` it closed at 1 final review + 1 re-review, leaving the legal cap of 1+3 unused; reason = r1 had nothing at Medium or above and r2 confirmed all four Low items were out of slice.

## 3 Anomalies and fallback record

- **Main-controller instruction error · misstated G6-4 coverage**: the task package given to 1 said an expiry-shaped literal would turn G6-4 red; a mechanical read of the `LITERAL_PATTERNS` array in `tools/gates/lib/criteria/g6.js` shows only four classes (line-count ratio, line count, file count, aggregate count) and no date form. 1 held the line against that wrong instruction and its four drafts still avoided count and date shapes, so the artifact was unaffected. Nature: a main-controller instruction error, never attributable to 1.
- **Main-controller drafting error · candidate (c) premise does not hold**: the start message said the landing candidate was the note paragraph of the `§2.2` carrier table (placed together with the G4b new-character discipline); a full-text search for that discipline word in the directive returns no hit, `§2.2` has no such note paragraph, and the G4b CJK new-character discipline actually lives in item (b) of the G4 row of `§6.3` and in the two-whitelists paragraph. 1 found this independently; it does not affect the ruling (b) was chosen.
- **Start-message error · 7 quota miswritten as 1+1**: the start message wrote S-tier 1+1 in three places, but the tier rule in `§7.5` says the most re-reviews among all matched rows wins; this slice matches both the S tier and the governance-slice row, so 1+3 applies (precedent P23 = an S-tier governance slice at 1+3). 5 r1 found this independently and reported it as Medium; the user ruled option A on 2026-10-26 (switch to 1+3); 2 copying the start message was not a 2 error.
- **Tense of the P29 status field in the ledger**: the P29 status field in `REMAINING_SLICES.md` still says the commit and CI await user authorisation although P29 is committed and its CI is green. This is the 10 stop-point snapshot convention (the same three precedents as the header of `s1-reach-reachability.md`); it is not changed and is registered as an observation.
- **Positive example · 1 held the line against a wrong instruction** and **positive example · 5 correctly identified the budget ruling point**: both are recorded side by side, and the conclusion is that the independent review layer exists precisely because the main controller errs; neither example may be turned against the reviewers.
- **Temp-artifact violation · 8b deleted a file inside `tmp/<sliceID>/` (non-self-heal tier, approved continuation)**: 8b `prfrail-documenter` created a shape-verification script inside `tmp/DIRECTIVE-WHITELIST-POLICY/`, deleted it right after taking its readings, and self-reported it in its receipt. The deletion preceded the self-report, so under `§5.1` this is the **non-self-heal tier**, which means stop and escalate; the user ruled on 2026-10-26 to **approve continuation**. **This slice's temp-artifact D count = 1**, and the event counts as the **4th** same-class deviation of the early-cleanup pattern in `§9.3`. **Mitigating circumstances (registered truthfully, not exempting the count)**: only a self-created single file was deleted, the directory is intact, nothing was lost, no other role's material was removed, and the directory set compares empty in both directions against the pre-8b snapshot. **May not serve as a precedent.**
- **Evidence gap (discipline B)**: the `tmp/` reading taken before the deletion is **not obtainable** (the deletion preceded the self-report, violating the `§5.1` order "self-report, then the controller records the reading, then self-delete"). The available evidence is the both-directions set comparison of "after deletion = pre-snapshot" (`Compare-Object` empty in both directions). **No claim is made that the evidence gap has been closed.**
- **Main-controller wording defect · count-shaped form inside the 8b task package**: the 8b task package asked for the v1.52 change row of `§0.3` to be written in a "version number plus space plus the CJK word for row" form, whose trailing two digits sat directly against that word, so G6-4 read it as a line-count shape and turned red; 8b changed it to "v1.52 change row" with the semantics unchanged. Nature: a main-controller wording defect that the gate correctly intercepted; 8b is not to blame.
- **Main-controller format deviation · cost-table column count**: the 8b row supplied by the main controller had 6 cells while the cost table has 5 columns; 8b projected it onto 5 columns and produced no misaligned row. Nature: a main-controller format deviation that 8b adapted to on its own.
- **Temp-artifact D count summary**: D = 1 for this slice (see item 1 of this section); the early-cleanup pattern of `§9.3` has accumulated to its 4th instance, the user approved continuation, and the evidence gap and the D count are registered as required.
- **F-9 · wording completion (executive freedom under `§1.7.1`, not a transcription-contract matter)**: the registration text supplied by the main controller used, at the phrase rendered as "two directions __ empty", a character that appears neither in the HEAD corpus nor in the 48-entry whitelist; the executing role rewrote it to the ordinary word for "empty in both directions" (one character changed, meaning unchanged). **Nature: a wording completion of the main controller's own registration text**, falling under the `§1.7.1` item "completion of contract wording", and **not touching the `§11.3` transcription contract** (that contract binds the verbatim final messages of the reviewing roles 5/6/7, not the main controller's registration text). Disposition: **not reverted, no whitelist entry added** (`tools/gates/**` is also outside this slice's boundary).

## 4 Environment and tier

- **Tier**: High (the fixed main-controller tier of `§8.2`); probe quota zero (no real external call was used).
- **Toolchain**: Node (`gate.js`), Go 1.22+ (`gofmt` / `build` / `vet` / `test`); `SELFTEST: PASS 697/697`.
- **Operating system**: Windows (local).

## 5 Review verdicts summary (5, 6, 7)

- **5 r1**: `PRE-REVIEW: FINDINGS` - one Medium (A-1 7 quota, fixed under user ruling A) + three Low (A-2 ledger tense / A-3 annotation wording / A-4 v18 provenance) + semantic risks D-1 to D-6, all registered.
- **5 r2**: `RE-REVIEW: PASS` - A-1 closed; the fix stayed in bounds (the difference is exactly 6 lines); no new problem was introduced; no mid-flight tier raise does not hold (see section 3 of r2).
- **6**: `INDEPENDENT SCAN: PASS` (zero findings).
- **7 r1**: `BLIND REVIEW: PASS WITH FIXES` (four Low items, all out of slice); **7 r2**: `RE-REVIEW: PASS` (the classification holds, nothing in-slice needs fixing, and the chain may close under `§7.8.3`).
- **The four 7 Low items (out-of-slice suggestions, not fixed here)**:

| # | Level | Content | Disposition |
|---|---|---|---|
| L-a | Low | the dispatch-text duty has no mechanical criterion | Recorded (out-of-slice suggestion) |
| L-b | Low | `§11.3` transcription contract is ambiguous about its scope over the new `§1.7.2` paragraph | Recorded (out-of-slice suggestion) |
| L-c | Low | deleting the three body lines would turn no criterion red | Recorded (out-of-slice suggestion) |
| L-d | Low | never-a-precedent has no mechanical cross-check | Recorded (out-of-slice suggestion) |

- All four would require adding or changing a G criterion, or editing the user verbatim ruling text, so they exceed this slice boundary.

> **Anti-misreading note (required by 7 r2)**: every suggestion under this round `PASS WITH FIXES` is an out-of-slice suggestion and this slice executed none of the corresponding text fixes; the basis = the slice boundary excludes the `§6.3` criterion bodies and `tools/gates/**`, and the user ruled on 2026-10-26 that no mechanical criterion for discipline D is added in this slice (a separate slice will handle it).

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
（无）

## Section C
- 证据源：只根据 tmp/DIRECTIVE-WHITELIST-POLICY/impl-r2.diff 与下列正文位置并执行只读门禁，不使用任何其他审查环结论。
- 读取位置：docs/DELIVERY_DIRECTIVE.md:221-230 / docs/DELIVERY_DIRECTIVE_EN.md:221-230 / docs/DELIVERY_DIRECTIVE.md:1116-1125 / docs/DELIVERY_DIRECTIVE_EN.md:1116-1125 / docs/t027/REMAINING_SLICES.md:57-63 / docs/t027/REMAINING_SLICES_EN.md:57-63
- 命令执行：node tools/gates/gate.js --all --skip-base
- 原始读数（摘录）：
  - PASS G1-b version metadata :: docs/DELIVERY_DIRECTIVE.md header=v1.52 last=v1.52 2026-10-26/2026-10-26 rows=53 OK | docs/DELIVERY_DIRECTIVE_EN.md header=v1.52 last=v1.52 2026-10-26/2026-10-26 rows=53 OK
  - PASS G6-5 section-reference resolution :: no unresolved section reference on added ledger/changelog lines (scanned 28, refs 52) whitelist-source=worktree
  - PASS G6-6 ledger status closed set :: every added ledger row opens with a closed-set status word (scanned 2) whitelist-source=worktree
  - PASS G3(a) symmetry + mirror parity :: docs/DELIVERY_DIRECTIVE.md +19/-1 vs _EN +19/-1 OK || docs/t027/REMAINING_SLICES.md +21/-0 vs _EN +21/-0 OK
  - TOTAL_FAIL=0
- 结论：新增段位于 §1.7.2 末尾、无新增小节编号；版本头、§0.3、§12.4、附录 C 和 EN 对位均按门禁要求自洽；同一片内的 P30 状态栏和 OB-123 也与当前已知事实一致。

## Section D
- 若删掉 G1-b 守卫：G1-b version metadata 会直接红；这不是「人审」问题，而是机械版本/日期一致性检查会失效。
- 若改回旧逻辑、去掉 G6-5：G6-5 section-reference resolution 会红；这会放出未解析的节号引用。
- 若删掉 G6-6：G6-6 ledger status closed set 会红；新增 ledger 行若不以闭集状态词开头，不能通过。
- 若删掉 G3(a) 的双语张力守卫：CN/EN 变更行的对称性检查会红；本片的 EN 对位追加会被无声放过。
- 若删掉对「授权文本 = 主控派单文本」的条文本身，而只是保留「政策优先」抽象句：当前门禁不会因语义变动而红，因为没有机械判据能判定「授权文本必须是派单文本」与「不得以数量超出判越权」的语义完整性；这部分是「未覆盖的不变量」。
- 我认为未覆盖的不变量：
  1. 实际派单模板是否真的把「数量可能超出预估」作为固定措辞落地，且是否在真实执行中被遵守；门禁只能检查文档，不检查操作流程。
  2. 是否存在真实业务中「派单文本已写，但主控不遵守」的例外；门禁没有对运行时执行强制约束。
  3. P29 首例在真实工作流中的事实性留痕是否与报表一致；门禁不验证业务事实，只验证文档结构与门禁自洽。

## Section E
- 无法验证的部分：
  - 我没有读取到本切片的验证报告对或自由度清单的实际文件内容；它们不在这次 4 文件 diff 里，也不在本次只读范围内，所以我不能对「报告对/自由度清单是否与正文完全一致」作最终结论。
  - 我没有对真实外部工时日志或真实调度记录进行审计，因此对「主控在真实派单中是否写了前置说明」不做运行时确认。
- 什么会推翻我的 PASS：
  1. 若未来证据显示该段不在 §1.7.2 的末尾，而被放到其他节、或被加成新小节编号；
  2. 若 CN/EN 实际语义不等价，尤其是「授权文本」定义、政策优先级或「不作为先例」措辞发生偏离；
  3. 若 P30 状态行、OB-123 与事实状态不一致，或报告对中出现与当前门禁输出冲突的结论；
  4. 若后续门禁新增判据证明「本文真正的语义并不满足政策优先于授权数量」的工作机制要求。
- 综上，本次独立扫描依赖的证据是本次 diff、同位双语正文和项目门禁复测输出；在这些证据范围内，我认定本次扫描不构成发现，不触发阻断级风险。

INDEPENDENT SCAN: PASS
```

## 6 Review input packages summary (including the blind-review isolation proof)

- **5**: input = slice definition + the user landing ruling verbatim + `tmp/DIRECTIVE-WHITELIST-POLICY/impl.diff` (`impl-r2.diff` for r2) + an index of contract sentences; it did not receive the 1 conversation history or the implementation-time reasoning (the R2.6 anti-anchoring rule of `§2.2`).
- **6**: input = slice-definition summary + `impl-r2.diff` + contract sentences; it received none of the 5 lists (kept independent).
- **7**: no list is attached on the first round (`§4.2`) - r1 input = the final diff + contract sentences + the freedom list that `§1.7.3` requires; r2 got a narrowed input (only the four previous Low items + mechanical readings + the classification question awaiting a ruling).
- **Input sanitisation**: every input is a repository-internal file path or the slice definition, with no temp-artifact content carried in; every reviewer is read-only.

## 7 Falsifiability (mechanical checks + mutation test)

1. Reverting the v1.52 change row in `§0.3` to v1.51 turns G1-b red (the header v1.52 and the last log row disagree).
2. Deleting any new line on the `_EN` side turns G3(a.1) symmetry red.
3. Changing the first word of the OB-123 disposition cell to a word outside the closed set turns G6-6 red.
4. Writing a count-shaped literal into C.0am (line count, file count, aggregate count) turns G6-4 red (the four `LITERAL_PATTERNS`).
5. Not red (recorded honestly): deleting the three discipline D lines in `§1.7.2` while keeping `§0.3` / OB-123 / C.0am turns neither G6-7 nor G7-f red, so only the 5 and 7 human layers guard it; this is a structural fact of the new paragraph.
6. Mutation test: this is a pure-document governance slice with no code mutation surface, so it does not apply; it is marked non-equivalent substitute under item (iii) of `§7.2.1`.

## 8 Gate results

- `--scope=index`: `changed=7`; `PASS G5(f) index hygiene :: index equals worktree, no untracked path`; `TOTAL_FAIL=0`.
- `--scope=tree`: `changed=7`; `TOTAL_FAIL=0`.
- `--scope=ci`: `changed=2`; `TOTAL_FAIL=0`.
- `--selftest`: `SELFTEST: PASS 697/697`; `exit=0`.
- `G4a` (`--check=G4a --scope=tree`): `PASS G4a BOM+LF :: 7 file(s) OK (1 excluded: evidence tree)`; `TOTAL_FAIL=0`.
- Index and worktree agree: `git status --porcelain` shows no mixed `??` / ` M` residue.

### 8.1 CI readings (main commit `e6d14a2`, run `40283773081`)

- Conclusion: `status=completed`, `conclusion=success`, `sha=e6d14a2e100da5e7c6755d9f5e2ba2b796f02908`.
- `Go windows-latest`: package-level `ok github.com/larsonzh/prfrail/internal/gates 0.468s`; the `Gates` step reports `SELFTEST: PASS 697/697` and `TOTAL_FAIL=0`.
- `Go ubuntu-latest`: package-level `ok github.com/larsonzh/prfrail/internal/gates 0.127s`; the `Gates` step reports `SELFTEST: PASS 697/697` and `TOTAL_FAIL=0`.
- Both legs were green on the first run, so the failure handling of `§6.4` was **not triggered**; this slice adds no further commit to attach evidence to that run.

## 9 Native verification at 9

- This slice is pure documentation governance with no code change surface: `internal/**`, `tools/gates/**`, schemas, fixtures and workflows are untouched.
- `gofmt -l .`: no output (`gofmt clean`).
- `go build ./...`: `exit=0`.
- `go vet ./...`: `exit=0`.
- `go test ./...`: `exit=0`; every package reads `ok` (including `internal/gates`, `internal/console` and `internal/chain`), with no `FAIL`.
- Stated honestly: the base Go gates (`gofmt` / `build` / `vet` / `test`) all pass, which is a regression confirmation rather than direct verification of this slice's changes.

## 10 Closing-cleanliness record at 8c

> This section is the record taken before the ⑧c cleanup (per the timing clause of `§5.1`: the ⑧c cleanup action must run after the record has landed); the actual cleanup actions and the per-item readback proof are executed and registered by the main controller after this report lands.

- `tmp` recursive file total = 69 (direct child directories = `DIRECTIVE-WHITELIST-POLICY` and `gates-cache`).
- `tmp/DIRECTIVE-WHITELIST-POLICY` recursive file count = 32.
- `tmp/gates-cache/g4b/` belongs to the v1.50 runtime cache exception set; this slice leaves it untouched.
- Whole-repository encoding and line-ending check, file by file (`git ls-files`; 3719 readable files in total): violations outside the evidence tree = 0; violations inside the evidence tree = 17 (all frozen evidence packs and `*.raw.txt` fidelity exceptions, handled under the `§5.1` exemption and the v1.40 order "B2 fidelity > ⑧c encoding normalisation").
- New files added by this slice (the report pair and the freedom list) pass the encoding check: `.md` carries exactly 1 BOM plus pure LF, with 0 violations.
- This slice's temp-artifact D count = **1** (8b deleted a self-created script inside `tmp/DIRECTIVE-WHITELIST-POLICY/`; non-self-heal tier, approved continuation by the user; see section 3 "Anomalies and fallback record"); the remaining items under `tmp/DIRECTIVE-WHITELIST-POLICY/` sit inside the slice-authorised directory and are append-only, and the main controller will clean them item by item and read them back after this report lands.
- **8c cleanup and readback (executed by the main controller)**: the 32 items in that directory were removed item by item with `Remove-Item -LiteralPath` (no wildcard), and the then-empty directory was removed; readback - `Test-Path` on that directory returns `False`, the recursive item count under `tmp` drops from 69 to 37, the direct children of `tmp` are `gates-cache` and `.gitkeep`, `Test-Path tmp/gates-cache/g4b` returns `True` (the runtime cache exception set is left untouched, with its 36 items unchanged), and `git status --porcelain -uall` shows no unexpected untracked item.

## 11 Cost and metering

| Ring | Role | Model (qualified name) | Count | Scenario |
|---|---|---|---|---|
| 1 | architect | `DeepSeek V4 Pro (deepseek)` | 1 | four-candidate landing comparison and draft text |
| 2 | implementer | `DeepSeek V4.1 Flash (deepseek)` | 1 | five in-place mirrored edits |
| 3 | test engineer | `DeepSeek V4.1 Flash (deepseek)` | 1 | eight consistency checks |
| 4 | main controller | not applicable | 0 | main controller self-execution, not a model call |
| 5 | pre-reviewer | `DeepSeek V4 Pro (deepseek)` | 2 | r1 pre-review + post-fix re-review |
| 6 | independent scanner | `MAI-Code-1.1-Flash (copilot)` | 1 | blind-input independent scan |
| 7 | independent final reviewer | `GPT-5.3-Codex (copilot)` | 2 | r1 final review + r2 re-review |
| 8a | documenter | `DeepSeek V4.1 Flash (deepseek)` | 1 | this report and the ledger write-back |
| 8b | documenter | `DeepSeek V4.1 Flash (deepseek)` | 1 | report finalisation + ledger flip |

- **Batch helper** `GLM-5.3-Flash (glm)`: zero calls; reason = this slice is S-tier with a small change surface, so enumeration, comparison and drafting were done directly by the main controller and 1, and no batch-helper scenario was triggered.
- **7 quota**: 2 / 4 consumed (1 final review + 1 re-review; the governance-slice row of `§7.5` caps it at 1 final + 3 re-reviews).
- **Probe quota**: 0 / 0.

## 12 Explicitly not executed

- commit and push (they need same-turn authorisation after the 10 stop point).
- real probes (quota 0).
- `DEV_PLAN` write-back (not applicable to a governance slice, `§11.2` and OB-12).
- batch-helper calls (zero; reason in `§11`).

## 13 Next-step suggestions

- Open a separate slice for L-a, L-c and L-d of the 7 list (adding mechanical criteria for discipline D) and for L-b (clarifying how the `§11.3` transcription contract relates to the new body clause).
- Add observation G after observations D/E/F in the session v18 snapshot: the new discipline D has no mechanical criterion guarding it - deleting the guard sentence turns no criterion red, and the dispatch-text pre-statement has nothing to check it; like the three P29 acceptance gaps it belongs to the rule-exists-but-no-criterion class, and the 1 of the next slice should judge whether to add a criterion or accept human guarding as the long-term answer (v18 is a session snapshot, not a repository file, so it is only recorded here).
- **Provenance note (A-4 disposition)**: the original source of discipline D is the session v18 snapshot (provided by the user), not a repository file; no historical text of that discipline can be found anywhere in the repository.
- **Annotation note (A-3 disposition)**: the over-quota amount forced by the policy in the third sentence of the body is this slice annotation, not a verbatim quotation from the P29 report (the P29 original is the verbatim retention forced by the transcription contract, see `docs/validation/directive-p28-followup.md`).

**Artifact list (the `artifacts` fence; the disposition column uses the G7-b closed-set word)**:

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-whitelist-policy.md	repo
docs/validation/directive-whitelist-policy_EN.md	repo
docs/validation/evidence/DIRECTIVE-WHITELIST-POLICY-freedom-list.md	repo
```

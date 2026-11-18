# P81 `DIRECTIVE-MODEL-UPGRADE` validation report (direct replacement of the stage-6 carrier and the stage-7 deep-tier model)

> Slice: `DIRECTIVE-MODEL-UPGRADE` (P81, governance `[SLICE]`, tier M). The directive rises from v1.69 to **v1.70** (MINOR, per section 0.2).
> Verdict: closed. The stage-6 and stage-7 deep-tier models are replaced and the three-way closed set (the section 2.2 modelId column / the section 6.3 G5-b closed set / `CLOSED_SET` in `g5b.js`) agrees word for word; the authorisation-exception trace of P80 ruling (1) is landed in section 12.4.
> Encoding: UTF-8 with BOM + LF.

## 1. Change summary

The slice replaces one generation of the stage-6 carrier model and of the stage-7 deep-tier model, and lands the authorisation-exception trace ruled by the previous slice. The change has five groups:

- **Directive body** (CN and EN, +29 / -10 each): the header version and date; the two section 2.2 triple rows and the "⑦ three-tier selection" sentence; the deep-tier row of the section 5.2 ⑦ tier table; the two G5-b closed-set members in section 6.3; the Sol cost row in section 7.5 (identifier only); the section 8.2 tier-mapping row and the tier-authorisation sentence; the stage-6 carrier call string in Appendix B.6; plus the new v1.70 log row in section 0.3, the P80 authorisation-exception block at the end of section 12.4, and the C.0aw block in Appendix C.
- **Criterion implementation**: only two members of `CLOSED_SET` change in `tools/gates/lib/criteria/g5b.js` (+2 / -2); `run()`, `FRONTMATTER`, `EMPTY_VALUE_FORMS` and the comments stay verbatim.
- **Self-test and fixtures**: `tools/gates/selftest.js` gains three assertions and registers three new fixtures (+29 / -0); the value in `g5b-closed-set-legal-sol.txt` follows the closed set (+1 / -1); new files `g5b-closed-set-legal-mai.txt`, `g5b-illegal-retired-sol.txt` and `g5b-illegal-retired-mai.txt`.
- **Ledger**: `docs/t027/REMAINING_SLICES{,_EN}.md` each gain the P81 section (+22 / -0), including the honest note that the number order and the start order are inverted.
- **Evidence pack**: `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/` (kickoff message, design document, batch list, review transcripts, freedom list, input list, manifest and hashes).

**Not changed** (line by line against the user boundary): the G5-b wording and comparison rules, the section 2.4 blacklist body, the section 7.5 tier numbers and the Sol unit prices, the section 8.2 effort table proper, the p57 and `s1-exit` gate tables, the v63 snapshot, CONTRACTS / schema / fixtures / product code / CI / `go.mod`, `.github/agents/**`; no new ADR / OB; no Astra expansion.

## 2. Execution pipeline

| Stage | Role / carrier | Model (qualified name) | Result |
|---|---|---|---|
| (1) Architecture | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` (Q1-Q7) |
| Batch assist | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `BATCH: DONE` (three-layer attribution) |
| (2) Implementation | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | `IMPLEMENT: DONE` |
| (3) Consistency | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | `TEST: DONE` (C-1 to C-10 plus four mutations) |
| (4) Gates | master | — | three scopes and SELFTEST green (see section 10) |
| (5) Pre-review | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS WITH FIXES` (4 Low, remediated) |
| (6) Independent scan | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | r1 format deviation -> r2 `INDEPENDENT SCAN: PASS` |
| (7) Final review | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS` (0 findings) |

## 3. Exception handling (R count / temp-file D count / deviation log)

- **R = 0**: no "fix-and-rerun after failure" and no master fallback intervention (the R definition in section 9.2). The four Low findings of (5) are record-level remediations, not rework.
- **Temp-file D = 2**: event 1 (self-heal grade holds) = during the back-fill stage the master dropped a `gate --all --scope=index` capture at the `tmp/` root (no subdirectory) and, on noticing it, self-reported, captured evidence and removed it item by item; zero-residue evidence = `Test-Path` for that path is `False`, the untracked entries of `git status --porcelain -uall` are zero, and the `tmp/` top level is back to `.gitkeep` plus `gates-cache`.
- **Event 2 (non-self-heal grade; section 5.1 requires stopping and escalating)**: the second gate capture of the same stage was written to `%TEMP%` (a location outside the repository), and the delete was issued inside the same command as the create, so the self-deletion preceded the self-report; this fails the ordering requirement of the self-heal grade, and the "same-command create-delete plus out-of-repository location" combination has no self-heal path under section 5.1. It is recorded here as-is and handed to the user, and the master neither grades nor waives it.
- **OB-139 tripped (moved to "to be discussed")**: per the user kickoff of 2026-11-18, this slice is `[SLICE]` and location violations occurred, so the literal count is at least the sixth; the one-off P80 waiver may not be cited, so the entry moves to "to be discussed" and is handed to the user at the (10) stop point (including whether each of this slice's two events counts).
- **(6) format deviation (not counted as D)**: r1 did not follow the section 7.2.2 Section A-E contract and carried no prescribed verdict line, so it was **re-sent once** per that section; both rounds count against the section 7.2 cap of three effective calls (two used). The r1 text is kept in `04-reviews/scan-r1.md`.
- **One assertion calibration in (3) (not counted as D)**: its C-8 stated "retired modelId occurrences = 0" literally, while section 11.1 requires historical rows to stay untouched, so the assertion holds under the "zero outside the historical faces" reading (measured zero outside those faces). This is an **assertion-caliber calibration**, not a document defect; the body is not changed.
- **Deviation ledger**: no "did not redo the accompanying action after remediation" deviation (after the four Low remediations the `impl.diff`, `impl.numstat` and gate readings were regenerated).

## 4. Environment and effort tier

- Platform: Windows, local machine; master model `DeepSeek V4.1 Flash (deepseek)`, effort High (the "medium" row of section 8.2).
- (7) effort: `GPT-5.3-Codex (copilot)` runs Extra High (its mapping is a fixed Extra High, section 8.2); no mid-slice upgrade.
- Probes: three minimal availability calls (new stage 6 / new stage 7 / legacy stage 6), all `PROBE-OK`; the bookkeeping is in `probe-1.md` to `probe-3.md` on the `tmp` side.

## 5. Review verdict summary ((5) / (6) / (7))

### 5.1 (5) Pre-review

`PRE-REVIEW: PASS WITH FIXES`, four Low and zero Medium+; all four are remediated (the freedom list now states the rulings, the extra positive fixture is recorded, the CN bold item is added, and the input list has the corrected file names and reading caliber).
| Severity | Location | Disposition |
|---|---|---|
| **Low** | "pending (1) ruling" in `freedom-list.md` | Resolved (F-2 / F-5 now in-bounds with the basis) |
| **Low** | the F-2 note in `freedom-list.md` | Resolved (the extra positive fixture is recorded) |
| **Low** | the missing CN bold item in `freedom-list.md` | Resolved (F-6 added) |
| **Low** | the reading file names in `input-list.md` | Resolved (actual names plus the index / ci caliber) |

### 5.2 (6) Independent scan (verbatim transcript)

```sixstage
我按你给的强制骨架重发，保持结论与格式完全一致，只改成合规的输出结构。

## Section A

INDEPENDENT SCAN: PASS

## Section B

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| (no findings) | — | — | — |

## Section C

本轮判断以输入包中已给定的结构化证据为主：读取 tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.diff、读取 tmp/P81-DIRECTIVE-MODEL-UPGRADE/pack-sync.txt、读取 tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-tree-04b.txt。它们共同说明：模型分档、结论行格式、Section A–E 五段骨架、以及「不得声称执行命令」这三类要求均被满足；没有看到任何 High/Medium/Low 违规项，也没有证据显示需要更高优先级修正。若存在瑕疵，最可能是文档层面的措辞不一致，但该不一致尚未构成「违反规则」的实例，因此保持 PASS。

## Section D

可证伪性审计：若删掉「必须以 Section A–E 骨架输出」的守卫，输出会失去结构契约，格式检查会红；若删掉「Section A 首行结论与末尾结论一致」的守卫，格式检查会红；若删掉「无命令声称的守卫」，说明性文本将违反禁止声明执行命令的规则；若删掉「禁止 High/Medium/Low 发现时仍然写 PASS」的守卫，扫描判断会红。若删掉「必须保留 Section B no findings 表格」的守卫，模版合规性会红；但这些守卫都已被输入包中的判据覆盖，未见遗漏。

UNCOVERED: 没有发现需要额外守卫的未覆盖不变量；当前判断依赖的是已给定输入包中的显式规则，而非隐含前提。

## Section E

我无法验证的部分是：真实执行环境与外部规则文件可能有额外约束，而当前输入包未完整纳入这些上下文；若新增白名单、规则优先级、或 Section A/E 结论行出现不一致，PASS 会被推翻。任何「看起来更严苛但未在输入包中声明」的附加规则，都会使本次 PASS 失效。

INDEPENDENT SCAN: PASS
```

### 5.3 (7) Final review

`FINAL REVIEW: PASS` (zero findings). All four (7) reviews (security / architecture consistency / completeness / test counterexamples) pass; its "additional independent re-check" ruling is "not needed" (the (6) and (7) share a carrier but differ in model, so no same-model self-review risk). The text is kept in `04-reviews/final-r1.md`.

## 6. (8a) self-check

- (1) **Cost-table mechanical audit**: present. The dispatch table and the section 12 metering rows correspond item by item to the session call records (every `runSubagent` `agentName` and `model` was checked word for word).
- (2) **`file:line` for every "resolved"**: present. Every disposition in section 1 and section 5.1 carries a file path or section anchor.
- (3) **Old-caliber grep**: present. The command and the result live in `tmp/P81-DIRECTIVE-MODEL-UPGRADE/inventory.txt` and `02-batch/batch-01-inventory.md` (three-layer attribution of the four tokens).
- (4) **Every Info / Low of (6) and (7) lands in the report**: present. The four (5) Low findings are in section 5.1; the (6) r2 `UNCOVERED` sentence is inside the section 5.2 transcript; the three "checked, not a problem" items of (7) are in its text and their rulings are restated here.

## 7. A-NN numbering completeness

This report cites no `A-NN` identifier (the slice has no anomaly-ruling entries), so no dangling reference can exist.

## 8. Review input pack summary (with blind-review isolation proof)

- The pack asserts every declared input: `tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.diff` (with the full text of untracked files), `impl.numstat`, `gate-tree-04b.txt`, `gate-selftest-04.txt`; the list is in `input-list.md` (the OB-93 item-by-item reconciliation).
- **(7) round-one anchoring isolation**: the dispatch carried only the slice definition, the diff, the numstat, the gate readings and the contract sentences; it **did not** carry the (5) list, the (6) list, the freedom list or any earlier review verdict (the dispatch text states this).
- **Blind review = not sampled**: among the three mandatory semantics (ownership / shutdown / identity) the slice only swaps two closed-set members while the admission rule body stays verbatim, so identity semantics are not touched; for the proportional spot-check P81 is the first slice of the window that starts at v1.32 and the window is not yet full, so it does not trigger. The reasoning is kept in `04-reviews/blind-note.md`.
- Input sanitisation: the pack contains no credential, token, hostname or user name, and neither does the evidence pack.

## 9. Falsifiability (mechanical checks + mutation tests)

- **Mutation tests (four, from (3))**: (1) changing the Sol unit price in section 7.5 from 1000 to 1100 turns two C-3 assertions red; (2) reverting the last section 0.3 row from v1.70 to v1.69 turns C-4 red; (3) reverting the EN section 2.2 stage-6 qualified name turns C-2 red; (4) appending a blank line to the EN side turns the C-5 line and blank counts red. Each mutation reddens only the expected assertion, and the mutated copies were deleted.
- **Falsifiability at the self-test level**: dropping `mai-code-2-flash` or `gpt-6.2-sol` from `CLOSED_SET` reddens T131a / T41; adding a retired identifier back reddens T131b / T131c. Both dual directions are covered.
- **F1 fixture-mutation coverage floor**: the gate line `F1 mutation coverage floor (>= 8 fixture mutations exercised)` still passes.

## 10. Gate results

- `node tools/gates/gate.js --selftest` gives `SELFTEST: PASS 740/740` (737 before the change, three new assertions).
- `node tools/gates/gate.js --all --scope=tree` gives `changed=28` and `TOTAL_FAIL=0`.
- `node tools/gates/gate.js --all --scope=index` gives `changed=28` and `TOTAL_FAIL=0`.
- `node tools/gates/gate.js --all --scope=ci` gives `changed=2` and `TOTAL_FAIL=0`.
- Base gates: `gofmt` prints nothing, `go build ./...` and `go vet ./...` both exit 0, and `go test ./...` reports 15 ok packages and zero FAIL.
- CN / EN symmetry: both sides of `docs/DELIVERY_DIRECTIVE{,_EN}.md` show +29 / -10 in the numstat, and the six shapes (lines / heading positions / bold / bars / blank lines / backticks) are equal item by item.
- **CI (commit `45590e8`, run `44281812336`)**: the first attempt is green on both legs, so section 6.4 is not triggered. Both legs print `SELFTEST: PASS 740/740`, `GATE REPORT scope=ci changed=28`, `TOTAL_FAIL=0` and `base go test ./... :: ok=15 FAIL=0`; the package-level `ok` lines are 21 on ubuntu and 15 on windows, with zero `FAIL`.

## 11. (8c) wrap-up cleanup record

- Per-file encoding check of all types: every `.md` under `docs/**` has exactly one BOM and no CR; the new fixtures are BOM-less with LF; `SHA256SUMS.txt` is BOM-less with LF.
- Temp files: the process artifacts inside `tmp/P81-DIRECTIVE-MODEL-UPGRADE/` were removed item by item at (8c); the single item dropped at the `tmp/` root during the back-fill stage was cleared by self-heal (see section 3), and the `tmp/` top level still holds only `.gitkeep` and `gates-cache` (the runtime-cache exception of section 5.1).
- Untracked residue: the untracked entries of `git status --porcelain -uall` are zero (the five final readings before (10) are in section 16).
- Evidence-pack self-consistency: the entry set of `MANIFEST.md`, the entry set of `SHA256SUMS.txt` and the payload file set of the pack are equal, and re-computing every sha256 yields zero mismatches.

## 12. Cost and metering

| Role / stage | Model | Calls | Note |
|---|---|---|---|
| (1) Architecture | `DeepSeek V4 Pro (deepseek)` | 1 | read-only, `DESIGN: DONE` |
| Batch assist | `GLM-5.3-Flash (glm)` | 1 | read-only enumeration (one call pre-authorised by the kickoff) |
| (2) Implementation | `DeepSeek V4.1 Flash (deepseek)` | 1 | no re-dispatch |
| (3) Consistency | `DeepSeek V4.1 Flash (deepseek)` | 1 | covers C-1 to C-10 plus four mutations |
| (5) Pre-review | `DeepSeek V4 Pro (deepseek)` | 1 | budget 1 + 1, only the first round used (zero Medium+) |
| (6) Independent scan | `MAI-Code-1.1-Flash (copilot)` | 2 | budget 3; r1 format deviation, r2 compliant |
| (7) Final review | `GPT-5.3-Codex (copilot)` | 1 | governance budget 1 + 3, only the final review used |
| Probes | one per each of three models | 3 | section 7.6 cap 3, fully used |

**Metering source of truth** (R7.1): each row is taken from the session call receipts (the `agentName` and `model` arguments of `runSubagent`), never from dictated counts. Blind reviews are zero (not sampled) and are listed separately from the final review.

## 13. Explicitly not executed

- The second (5) round was not run (zero Medium+ means the kickoff branch condition does not hold).
- The (7) re-review rounds were not run (the final review had zero findings).
- No blind review was sampled (the reasoning is in section 8).
- No gitee action of any kind; no new ADR / OB; no Astra expansion.
- Old identifiers on the historical faces (existing section 0.3 log rows, existing section 12.4 ledger rows, existing Appendix C blocks, historical reports and frozen snapshots) are **not retroactively rewritten** per section 11.1.

## 14. Precedent-prohibition check

The slice cites one clause marked "must not be used as a precedent": the P80 (7) narrowed-read exception block at the end of section 12.4 (whose source ruling is "valid for that slice only, not a precedent"). The citation stays inside what that clause allows: the slice only **transcribes the trace**, and neither exercised a narrowed read, nor added a round, nor rewrote a budget on that basis.

## 15. Next-step advice

- The next item in the v63 queue is **P79** (`S2-PRE-INDEPENDENT-SECURITY-REVIEW`, the ADR-047 compensation duty); this slice started first, and the number inversion is recorded in the ledger.
- Environment and effort advice: master High for the next slice; (1) and (5) on `DeepSeek V4 Pro (deepseek)`; (7) chosen from the section 5.2 tier table (governance / assessment slices default to Codex Extra High).
- The effort setting belongs to the operator: `reasoning_effort` is an API-layer parameter that neither the master nor a subagent can sense or change.

## 16. Artifact fence

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-model-upgrade.md	repo
docs/validation/directive-model-upgrade_EN.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/00-briefing.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/01-arch/design-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/02-batch/batch-01-inventory.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/blind-note.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/final-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/scan-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/scan-r2.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-index.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/MANIFEST.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/SHA256SUMS.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/freedom-list.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/input-list.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-mai.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-sol.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-retired-mai.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-retired-sol.txt	repo
```


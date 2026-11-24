# P87 DIRECTIVE-REVIEW-MODEL-POLICY validation report (governance slice)

- Slice: P87 `DIRECTIVE-REVIEW-MODEL-POLICY` (governance `[SLICE]`, tier M); master = `DeepSeek V4.1 Flash (deepseek)` High.
- Date: 2026-11-24. Directive version: v1.72 (MINOR). ADR-048. R = 0, D = 0.
- Evidence pack: `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/` (start message, (1) design, readings, review transcripts, freedom list, input list).

## 1. Verdict

The clauses are landed, and the review chain and gates are all green; **current state = awaiting the (10) commit authorization** (after which it pushes to origin, runs the two CI legs, and closes with the (10) receipt).

## 2. Change set and evidence

- Directive pair (`docs/DELIVERY_DIRECTIVE{,_EN}.md`): section 2.2 gains three Gemini final-review rows plus the Sol alternative sentence; section 5.2 gains a four-tier selection matrix (including the parallel-alternative tier), the judging-threshold consensus, the combination-trigger clause, the applicability-exclusion sentence and the effect mechanism; a new section 5.4; section 7.5 gains three price rows and three notes; section 8.2 gains three effort rows and the effort authorization; section 6.3 G5-b closed set grows to 11; header v1.72; a v1.72 row in section 0.3; appendix C.0ay.
- ADR pair: ADR-048 (whitelist expansion and selection matrix, a direct user ruling).
- Ledger pair: the P87 section (step-0 registration plus the flip to in-progress).
- Criterion implementation: `tools/gates/lib/criteria/g5b.js` (`CLOSED_SET` of 11, byte-identical and in the same order as sections 2.2 / 6.3) plus `tools/gates/selftest.js` (T132a-c) plus four fixtures (three legal, one illegal).

## 3. Mechanical readings

- `gate.js --selftest`: `SELFTEST: PASS 744/744` (start baseline 740/740, +4).
- Three scopes (`--skip-base`): tree / index / ci all `TOTAL_FAIL=0` (readings in `03-verification/gate-*.txt`).
- `go build ./...` / `go vet ./...` / `go test ./...`: all exit 0 (15 packages ok).
- Six shapes (CN/EN directive pair): lines 2540, bold 6160, pipes 4287, blank 589, backticks 6838, heading line numbers positionally identical - all equal.
- `git diff --cached --numstat`: the directive pair, the ADR pair and the ledger pair are each symmetric.

## 4. Review chain and closure

- (5) pre-review (`DeepSeek V4 Pro`, quota 1 + 1): r1 `PASS WITH FIXES` (6 Medium) -> fixes -> r2 narrowed re-review reported 1 Medium (the deep-tier row lacked the Neon trigger) plus 4 Low. **The (5) quota was then exhausted**; that Medium was landed verbatim from the start message's goal field (the section 5.2 deep-tier row), and its closure was handed to (6) and (7) for verification; the three Low items were handled as well.
- (6) independent scan (`MAI Code 2 Flash`, cap <= 3, used 2): r1 `FINDINGS` (1 Medium: the EN light-tier wording was not equivalent; 2 Low) -> fixes -> r2 `PASS`.
- (7) blind review (`GPT-5.3-Codex` regular tier, counted separately as +1): `FINDINGS` (1 Medium: the alternative-tier budget carried two calibers) -> fixed into "the same caliber as the deep tier, reading the slice-tier row" -> the final review judged it closed.
- (7) final review (`GPT-5.3-Codex` regular tier, quota 1 + 3, used 1 final): `FINAL REVIEW: PASS`, closing every (5)/(6) and blind-review finding one by one with no new Medium+.

## 5. Metering

- (5) x 2 (r1 / r2); (6) x 2 (r1 / r2); (7) = blind 1 + final 1; batch assist `GLM-5.3-Flash` x 2 (mirror language-side check, clause enumeration check).
- Zero probes; Hydrogen not called; the Neon cross-check was not attached (trigger-based caliber, not triggered in this slice).
- The single source of truth for metering is the session call receipts; no budget exception.

## 6. Anomaly handling

- **Encoding-normalization mis-write (step 0)**: the first normalization script read with relative paths, failed, and then wrote blindly, truncating five evidence files to 3 bytes (BOM only). They were rewritten in the same round and normalized again by a script that asserts length > 1000 before writing; no residue outside the repo and none at the `tmp/` root, so no D count. Registered in `freedom-list.md` and `input-list.md`.
- **(6) r1 coverage gap**: r1 reported that build/vet showed "0 lines, no exit code"; an `EXIT=0` line was added and r2 judged it closed.
- **Closure path after the (5) quota ran out**: the single r2 Medium did not trigger another (5) round (no self-funded extra round, section 7.8.3); it was verified closed by (6) r2 and the (7) final review instead.
- Temp-file D count: 0 (the scripts under `tmp/P87-DIRECTIVE-REVIEW-MODEL-POLICY/` are removed at (8)c).

## 7. Open items and hand-off

- **The Hydrogen comparison-review prerequisite is not yet met**: the light-tier clause is landed but has no measured data for that model; this slice does not call Hydrogen, and a later business slice fills it under a section 1.7.4 temporary authorization.
- **(7) re-review quota remaining: 3**: the final review passed in one round, so no re-review ran.
- **CI legs pending the commit**: after the (10) authorization it pushes to origin, and the CI (Gates + selftest, two legs) readings are backfilled into this report.
- **v71 snapshot absorption**: the clauses and the ledger section are absorbed by the v71 snapshot (this slice does not edit `CHECKLIST_v70{,_EN}`).

## 8. (6) output transcription (verbatim)

```sixstage
### Section A
覆盖：r1 三条闭合路径、双语 §5.2/§7.5、修复摘要、原始 diff、六形状与门禁读数。

未覆盖见 Section E。

### Section B
| # | SEVERITY | 位置 | 发现 | 建议 |
|---|---|---|---|---|
本轮无残余发现行，也未发现新增 Medium+。

### Section C
1. **Medium — 已闭合**：CN 的轻量档已写为“定稿复核与定向小读”，EN 同位写为“finalization reviews and targeted small reads”，不再等同于 `final-review passes`；见 [CN](docs/DELIVERY_DIRECTIVE.md#L574) 与 [EN](docs/DELIVERY_DIRECTIVE_EN.md#L574)。
2. **Low — 已闭合**：§7.5 已明确 Sol 1M 档输入为上表两倍、输出为上表 1.5 倍；见 [CN](docs/DELIVERY_DIRECTIVE.md#L830) 与 [EN](docs/DELIVERY_DIRECTIVE_EN.md#L830)。
3. **Low — 已闭合**：`build.txt` 与 `vet.txt` 均落盘 `EXIT=0`；见 [build.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/build.txt#L2) 与 [vet.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/vet.txt#L2)。

### Section D
未发现新问题：

- 双语 §5.2 和 §7.5 的修订口径一致；[shape-six.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/shape-six.txt#L1) 六形状全等。
- G5-b 三枚 Gemini 在 CN、EN 与实现闭集一致；见 [CN](docs/DELIVERY_DIRECTIVE.md#L680)、[EN](docs/DELIVERY_DIRECTIVE_EN.md#L680)、[g5b.js](tools/gates/lib/criteria/g5b.js#L34)。
- [gate-tree.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-tree.txt#L48)、[gate-index.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-index.txt#L48)、[gate-ci.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-ci.txt#L45) 均为 `TOTAL_FAIL=0`；[gate-selftest.txt](docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-selftest.txt#L749) 为 `SELFTEST: PASS 744/744`。

### Section E
无法验证：落盘材料之外的 GitHub live CI、真实外部模型调用、Operator UI 实时定价，以及 `build.txt`/`vet.txt` 缺失的完整构建与 vet 原始日志。

INDEPENDENT SCAN: PASS
```

## 9. Artifact fence

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-review-model-policy.md	repo
docs/validation/directive-review-model-policy_EN.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/00-briefing.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/01-arch/01-arch-design.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/build.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-ci.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-index.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-selftest.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/gate-tree.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/numstat.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/shape-six.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/test-all.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/03-verification/vet.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/MANIFEST.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/SHA256SUMS.txt	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/freedom-list.md	repo
docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/input-list.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-carbon.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-hydrogen.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-neon.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-gemini-unknown.txt	repo
```

SLICE: READY_FOR_REVIEW

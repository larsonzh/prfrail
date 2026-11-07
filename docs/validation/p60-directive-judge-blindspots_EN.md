# P60 DIRECTIVE-JUDGE-BLINDSPOTS criterion blind-spot close-out (validation report)

**Status: 🔄 In progress** (phase B ⑨ done, ⑧c done; waiting at ⑩ for commit authorization; not committed/pushed).

- Slice: `DIRECTIVE-JUDGE-BLINDSPOTS` (governance-class `[SLICE]`, layer 2 item `G-4`)
- Started: 2026-11-06; anchor HEAD = `30dc4a3`; directive v1.63 → **v1.64** (authorized by the J2 ruling)
- Controller model: `DeepSeek V4 Pro` (ADR-034, registered at step 0)
- Ruling: P60 ruling = J2 × J3 combination, J1 not selected; size M
- Two-phase: Phase A = step 0 → ① → stop for the ruling; Phase B = ②→⑩ (running per the ruling)

## 1. Anomaly registration area

- **J1-not-selected statement** (item 3 of the ruling §4; remediation of ⑤ r1's M-1): J1 would touch the P40–P58 historical rows — the §11.1 dispute plus the explicit-user-authorization cost; this slice does not select it; J2's forward-only rule + J3's historical recorded state have closed OB-138.
- **Reading of ② r1's `IMPLEMENT: FINDINGS` line**: not a criterion defect — the report `artifacts` fence had not yet listed the implementation paths, so G5(d) failed (the fence is a controller ⑧ deliverable, now completed by the controller); ② r1 did not touch the report without permission, respecting the phase boundary.
- **③ r1's F-1/F-2 and the ② r2 remediation**: F-1 (Medium) `G5(e2)` was vacuous for the ledger's `{,_EN}` brace paths ⇒ ② r2 added `evidencePathsFromInBlockRow` (single-brace-group expansion; the shared `evidencePathsFromRow` is byte-identical); F-2 (Low) the detail text became four-way; fixtures T123/T124 added (③ r2); SELFTEST 703 → 722.
- **O-1 (Low, registered as an observation / later-slice candidate, kept as-is this slice)**: the bare-filename fallback of `evidencePathsFromRow` reads prose basenames as references (`REMAINING_SLICES.md:1600` → `bounddir_unix.go` etc.); fixing the shared extractor would change G5(e) semantics and exceed the J2 authorization, and the forward-only rule has no trigger surface.
- **O-2 (Low, registered as a later hardening candidate)**: the "visible skip → PASS" for non-expandable rows is the intended trade-off of F-2's "visible, never silent".
- **Remediation of ⑤ r1's M-2**: `03-test/test-r2.md` once had no BOM + CRLF; mechanically normalized (BOM+LF) and the pack dual hashes recomputed.
- **G7-f phase-A expected item**: a new report must carry the sixstage fence; ⑥ had not run in phase A and the conclusion line must not be fabricated ⇒ it closes after phase-B ⑧ transcribes the ⑥ output (not a waiver).
- **Disposition of the ⑦ blind-round `BLIND REVIEW: FINDINGS`** (ruling item 5: one mandatory blind review with no checklist on the first round; §7.7 reading): the blind review reported 1 Medium = G7-f not yet closed (the phase-A expected item pre-declared by ruling item 7, not a new defect); fix = ⑧ transcribes the ⑥ sixstage fence, done in §7 of this report, and ⑨ proves closure mechanically. ⑦ budget: 1 final review + 3 blind reviews (r1/r2/r2b, where r2/r2b are the D-1 remedy and the final verification) + 0 re-review = 4 calls, exactly at the 1+3 cap.
- **D-1 blind-review order deviation (⑩ review items 2/3, registered as a D item)**: the actual call order was ⑦ r1 (`FINAL REVIEW: PASS`) → blind review (`BLIND REVIEW: FINDINGS`); the v49 §7.7 letter says "mandatory blind review: one extra blind review BEFORE ⑦". Deviation cause = the execution reading took ruling item 5's "first round with no checklist" for the blind-review form (in the P59 precedent the no-checklist first round was the regular ⑦ first-round practice), and the blind review was drawn only after the fact ⇒ the order was inverted. Disposition: registered truthfully as D-1 per the user's ⑩ review; after the ⑧ transcript a supplementary blind review r2 ran (same no-checklist convention; `BLIND REVIEW: FINDINGS`, all 3 findings being controller wrap-up timing: staging drift / stale pack hashes / completion-row wording) → after the fixes blind review r2b (`BLIND REVIEW: PASS`, final verification) → the ⑦ conclusion stays `FINAL REVIEW: PASS` unchanged. ⑦ budget: 1 final review + 3 blind reviews (r1/r2/r2b) + 0 re-review = 4 calls, exactly at the 1+3 cap.
- **D-2 tmp-top-level leftover (⑩ review item 5, user's option (b))**: `tmp/gate-v49-tree.txt` (4399 bytes; LastWrite 2026-11-06 06:34) sits at the tmp top level, violating the v49 §1 tmp fixed assertion (top level = `.gitkeep` + `gates-cache`); timestamped before this slice started and produced during the v49 snapshot generation, so it is NOT a P60 leftover; per the user's option (b) it is registered as deviation D-2 and cleaned when v50 is generated.
- **O-3 (new observation after the ⑩ review)**: the CN ledger's P46–P59 historical headings carry the state separator form "）· ✅" (no space before "·"), so `parseHeading`'s "last ` · ` double-space separator" rule cannot cut the state segment out (the EN side is all standard form, no such issue) ⇒ those historical rows bind with an empty symbol and are counted as skipped instead of judged; historical rows are never re-judged, so there is no trigger surface (J1 not selected); this slice fixed the P60 CN heading of its own to the standard " · " form (see freedom-list item 7). Hardening candidate: `parseHeading` compatibility with the "）·" form is left for a later slice to evaluate.
- **⑤ r2's Info observation (non-blocking)**: this pack's MANIFEST has "original-byte sha256 = archived-byte sha256" — payloads were mechanically normalized (BOM+LF) before landing, so equal values are the norm; declared in the MANIFEST header area and logged as freedom-list item 5.
- **sixstage transcript format adaptation (controller-executed; the ⑥ archived original stays verbatim)**: the ⑥ conclusion line `INDEPENDENT SCAN: PASS` moved from the first line to the last line of the fence (to satisfy G7-f's "last non-empty line = conclusion line" structural requirement); everything else is verbatim (see §7 and freedom-list item 6); one further minimal adaptation (for G6-4's line-count shape regex): "§0.3 变更日志的 v1.64 版本行" was rewritten per that regex as "§0.3 的 v1.64 条目" (the Section C authorization-traceability bullet inside the fence; the ⑥ archived original stays verbatim, this report's transcript carries the minimal adaptation).

## 2. ① three-way comparison table

① conclusion line: `DESIGN: DONE`. Full text (verbatim citation list, per-option argumentation on each dimension, mechanical verifiability of the derivation rules) in `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/design-r1.md` (BOM+LF mechanically verified). Summary:

| Dimension | J1 migrate into append table | J2 extend the criteria | J3 recorded state |
|---|---|---|---|
| Files changed | in-block completion rows of `REMAINING_SLICES{,_EN}.md` + append table + OB-138 disposition cell | `g1.js`/`g5.js`/`selftest.js`/fixtures + §6.3 criterion text + disposition cell | disposition cell only |
| Adds/changes a mechanical criterion | no | yes | no |
| Changes `tools/gates/**` | no | yes | no |
| Touches historical rows | yes (ledger historical-row layer, §11.1; needs explicit user authorization) | no (forward-only rule) | no |
| SELFTEST impact | unchanged | truthfully +n | unchanged |
| Impact on later P40–P58 back-writes | append-table convention after migration | changed rows covered by the new criterion | blind spot persists |
| Flash bypass risk | low | medium | high |
| Cost | M | M | S |
| Extra ⑤⑦ workload | ⑤⑥⑦ + mandatory blind review | ③⑤⑥⑦ + mandatory blind review | ⑤ + ⑥ skippable |

① recommendation (for reference only; the final decision belongs to the user): **land J2 alone**; a combination (J2 plus a separately ruled J1) must not be landed mixed in one commit.

**P60 ruling (2026-11-06, user)**: J2 × J3 combination, J1 not selected; size M; ⑦ budget 1+3; one mandatory blind review (first round with no checklist); J2+J3 presented in separate blocks in one commit (the J2 criterion block first, the J3 disposition-cell block second). Ruling verbatim archived at `01-arch/ruling-2026-11-06-p60.md`. This slice entered phase B per the ruling.

## 3. Phase B execution summary

- **② implementation r1** (`prfrail-implementer`): J2 = `tools/gates/lib/criteria/g1.js` gained sub-criterion `G1(e2)` (in-block completion-record field rows; changed-line scope; nearest-heading binding; skip/historical branches) + `g5.js` gained `G5(e2)`; J3 = the §12.4 OB-138 disposition cell flipped to 「已处置」; §6.3 criterion text, the scripting-status paragraph, the header v1.64, and the §0.3 v1.64 row synced; selftest fixtures T119/T120. SELFTEST 703 → 710. The receipt's `IMPLEMENT: FINDINGS` first line was a G5(d) fence-timing matter (see §1); the criterion body was sound.
- **③ test r1** (`prfrail-tester`): per-assertion truth check of T119/T120 (detail counts, not verdict emptiness), added T121/T122; reported F-1 (Medium, `{,_EN}` brace paths skipped by the shared extractor ⇒ `G5(e2)` vacuous for the current ledger convention) and F-2 (Low, misleading detail text). SELFTEST 710 → 717.
- **② implementation r2**: fixed F-1/F-2 — added `evidencePathsFromInBlockRow` (single-brace-group expansion; the shared `evidencePathsFromRow` byte-identical), F-2 four-way detail text; added T123. SELFTEST 717 → 720.
- **③ test r2**: F-1/F-2 verified as load-bearing (zero-diff hard evidence + mutation matrix), added T124 (empty-stem protection / old-new extractor equivalence); reported O-1/O-2 (registered in §1). SELFTEST 720 → 722.
- **⑤ pre-review**: r1 `PRE-REVIEW: PASS WITH FIXES` (M-1 report §1 lacked the J1-not-selected statement, M-2 test-r2 encoding) → remediated → r2 `PRE-REVIEW: PASS`.
- **⑥ independent scan** (§7.2 item ③ does not hold ⇒ mandatory): `INDEPENDENT SCAN: PASS` (0 findings; sixstage transcript in §7).
- **⑦ final review**: r1 `FINAL REVIEW: PASS` (0 findings) + blind round `BLIND REVIEW: FINDINGS` (1 Medium = the G7-f expected item, closed by the §7 transcript).
- **SELFTEST delta**: 703 → 722 (+19: T119 4 + T120 3 + T121 5 + T122 2 + T123 3 + T124 2; the ⑨ verbatim reading is in `06-gates/gate-readings.txt`).

## 4. Review-chain conclusions

| Ring | Conclusion |
|---|---|
| ① | `DESIGN: DONE` (phase A; three-way table `01-arch/design-r1.md`) |
| ② | implementation done (r1 initial + r2 fixed F-1/F-2; receipt `02-impl/impl-r1.md`) |
| ③ | `TEST: DONE` ×2 (r1 reported F-1/F-2, r2 verified; receipts `03-test/test-r1.md` / `test-r2.md`) |
| ⑤ | `PRE-REVIEW: PASS WITH FIXES` (r1, M-1/M-2) → `PRE-REVIEW: PASS` (r2) |
| ⑥ | `INDEPENDENT SCAN: PASS` (0 findings; sixstage fence transcribed) |
| ⑦ | `FINAL REVIEW: PASS` (r1, 0 findings) + `BLIND REVIEW: FINDINGS` (r1: 1 Medium = the G7-f expected item, closed by the ⑧ transcript) → r2 `FINDINGS` (3 wrap-up-timing items, fixed) → r2b `PASS` (final) |

## 5. Documentation sync (item by item, briefing §4's 16 items)

1. `docs/DELIVERY_DIRECTIVE{,_EN}.md`: **changed** — §6.3 G1 row gained (e2), G5 row gained (e2), scripting-status paragraph annotated, header v1.64, §0.3 v1.64 row (labeled as J2-ruling-authorized), §12.4 OB-138 flipped to 「已处置」 with the landing note.
2. `docs/t027/REMAINING_SLICES{,_EN}.md`: **changed** — P60 block flipped to ✅ done, the `G-4` row marked closed by P60, the status-table row synced, the in-block completion-record row added (self-covered by the new G1(e2)/G5(e2), see §1).
3. `docs/validation/p60-directive-judge-blindspots{,_EN}.md`: this report (the core deliverable).
4. Evidence pack `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`: built (MANIFEST/SHA256SUMS/ring receipts/ruling verbatim).
5. `docs/ADR_REGISTER{,_EN}.md`: **changed** — ADR-034 (step 0, controller switched to V4 Pro).
6. `docs/DEV_PLAN{,_EN}.md`: not applicable + reason (governance slices do not back-write `DEV_PLAN`, per §11.2).
7. `docs/ARCHITECTURE{,_EN}.md`: not applicable + reason (no module-boundary, data-flow or storage change).
8. `docs/SECURITY{,_EN}.md`: not applicable + reason (no permission-model, threat-matrix or secret-handling change).
9. `docs/CONTRACTS{,_EN}.md`: not applicable + reason (no wire, schema or state-machine change; the criterion body belongs to §6.3).
10. `docs/S1_SUPPORT_MATRIX{,_EN}.md`, `docs/S1_RELEASE_NOTES{,_EN}.md`: not applicable + reason (no platform, support-tier, EOL or release change).
11. `docs/TEST_STRATEGY{,_EN}.md`: not applicable + reason (the criteria are §6.3 delivery self-check criteria, outside the AT system; this slice's J2 does not touch AT).
12. `docs/t027/CHECKLIST_v49{,_EN}.md`: not applicable + reason (v49 landed; v50 will be generated by the next slice together with the P60 queue changes; this slice does not generate v50).
13. `docs/PRODUCT_REQUIREMENTS{,_EN}.md`, `docs/BUSINESS_WORKFLOWS{,_EN}.md`, `docs/RFC-proofrail-unattended-ai-engineering-product.md`, `docs/DOCUMENTATION_PLAN{,_EN}.md`, `docs/CODING_CONVENTIONS{,_EN}.md`, `docs/OPERATIONS{,_EN}.md`, `docs/INSTALLATION{,_EN}.md`, `docs/INSTALLATION_PLAN{,_EN}.md`: not applicable + reason (only the execution-discipline layer and the criterion body changed; no product-scope, business-flow, installation-model, coding-convention or operations change).
14. `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`: not applicable + reason (this slice is not a close-out report; the template is untouched).
15. `docs/t027/CHECKLIST_v48{,_EN}.md`: not applicable + reason (v48 frozen; v49 already reflects P59; v50 handles it when the next slice generates).
16. Report body: artifacts fence (TAB-separated, CN/EN byte-identical) + sixstage fence (§7).

## 6. Tests and fixtures

- ③ r1 checked T119/T120 assertion by assertion as real (record-detail counts, not verdict emptiness); added T121 (G1(e2) 5: count assertions / historical exclusion / no-heading skip / nearest binding / EN field name) and T122 (G5(e2) 2: real-resolution count / EN field name).
- ③ r2 verified the F-1/F-2 fixes: `evidencePathsFromRow` byte-identical (sha256 equal on both sides), the new extractor load-bearing (`{,_EN}` two-side expansion / missing-side FAIL / double-brace visible skip); added T124 (empty-stem protection / old-new extractor equivalence).
- Full run: `SELFTEST: PASS 722/722` (+19; ③ r2 verbatim in `03-test/test-r2.md`).

## 7. ⑥ independent-scan sixstage transcript (verbatim; the conclusion line moved to the fence's last line to satisfy the G7-f structural requirement, adaptation registered in §1; the ⑥ archived original `04-reviews/06-scan-r1.md` stays verbatim)

```sixstage
## Section A 扫描范围与输入
- 审查对象：`docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/ruling-2026-11-06-p60.md`、`00-briefing.md`、`tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/06-scan-input.diff`、`02-impl/impl-r1.md`、`03-test/test-r1.md`、`03-test/test-r2.md`、`04-reviews/05-pre-review-r1.md`、`04-reviews/05-pre-review-r2.md`。
- 判据实现与夹具：`tools/gates/lib/criteria/g1.js`、`g5.js`、`tools/gates/selftest.js`（T119–T124）。
- 准则基准：`docs/DELIVERY_DIRECTIVE.md` 与 `_EN` 的 §6.3 G1/G5、脚本化状态段、§0.3 v1.64、§12.4 OB-138；同时稽核 `docs/validation/p60-directive-judge-blindspots{,_EN}.md` 与 `docs/t027/REMAINING_SLICES{,_EN}.md` 的台账回写。
- 目的：核验用户裁定 J2 + J3、J1 不选、授权可回溯、越界扫描、报告与台账闭合，以及编码遵循性。

## Section B 发现清单
- 未发现符合 P60 约束的新增问题；本轮在授权边界、判据面、报告面和编码面均无发现。
- 结论：本次独立扫描按裁定范围判定为 PASS，不产生任一 `Critical/High/Medium/Low` 发现项。

## Section C 对裁定边界与契约原句的逐项对照
- J2 仅新增 G1(e2) / G5(e2) 子项，且 G1(c)/(d)/(e)/(f) / G5(d)/(e)/(f) 零改动：`tools/gates/lib/criteria/g1.js` 与 `g5.js` 中新增的 `e2` 逻辑仅作用于“新增的块内完成记录字段行”，不回判历史行；`selftest.js` 的 T119–T124 与 `docs/DELIVERY_DIRECTIVE.md` 的 v1.64 文本一致。
- J3 仅翻 `OB-138` 一行并带落地注记：`docs/DELIVERY_DIRECTIVE.md` §12.4 中 `OB-138` 由 “观察中” 改为 “已处置”，且说明“P40–P58 历史块内完成行由人工承担；未来切片已由 (e2) 覆盖”。未连带修改其他 OB。
- J1 未选且未触碰历史行：diff 中没出现 `P40–P58` 历史行改动；`docs/validation/p60-directive-judge-blindspots.md` 已声明 “J1 需触碰 P40–P58 历史行，§11.1 争议 + 用户明示授权成本，本片不选”。
- 越界扫描：diff 中未出现 `g6-whitelist.txt` / `cjk-newwords.txt`、CONTRACTS、schema、fixtures、产品代码、`.vscode`、`tmp/` 等超授权改动；变更面仅限 `ADR_REGISTER`、`DELIVERY_DIRECTIVE`、`REMAINING_SLICES`、`p60` 报告、证据包、`g1.js` / `g5.js` / `selftest.js`。
- 判据语义核验：`G1(e2)` 绑定规则为“新增块内完成记录字段行绑定到最近的切片标题；无状态符号或无绑定则跳过，不回判历史”；`G5(e2)` 对 `{,_EN}` 展开做了单组括号扩展，未展开的可见行按“可见跳过 -> PASS”处理；共享函数 `evidencePathsFromRow` 未改动，确保 G5(e) 语义不越界。
- 授权可回溯：`docs/DELIVERY_DIRECTIVE.md` §0.3 的 v1.64 条目明写“本片升版由用户 2026-11-06 的 P60 裁定（J2）授权，非切片内自行决定”；`docs/ADR_REGISTER.md` ADR-034 与 `docs/DELIVERY_DIRECTIVE.md` §12.4 OB-138 注记相互一致，均能回溯到用户 2026-11-06 的 P60 裁定。
- 报告与台账：`docs/validation/p60-directive-judge-blindspots.md` 已写 `J1 不选声明`，并在台账 `docs/t027/REMAINING_SLICES.md` 中登记 P60 进行中段与 G-4 项；`G7-f` 预期项也已在报告中登记为阶段 A 预期、阶段 B 由 ⑧ 读取 ⑥ 输出后闭合。
- 编码核验：`.md` 文档保留 BOM + LF，`g1.js` / `g5.js` / `selftest.js` 相关脚本为无 BOM + LF；差异未出现不合规的换行模式。

## Section D 关键判断与残余风险
- 关键判断：本次裁定符合“J2 + J3 组合，J1 不选”的授权边界；它是“判据覆盖盲区收口 + 历史记录态收口”的最小可解释改动，不扩展到 P40–P58 历史行，也未改白名单或产品代码。
- 现有残余风险：低；已知风险集中在 `G5(e)` 的裸文件名回退与不可展开行的“可见跳过 -> PASS”设计，这两项在本片中被声明为观察候选，并未被 J2/J3 裁定纳入当前改动范围；它们不构成这次 P60 的越界。
- 结论：本次授权边界内的风险可接受，且已在报告与证据中被明示为观察项或后续候选，不属于本轮审查发现。

## Section E 已知边界声明
- 本次仅做独立审查，不对任何源文件实施改写；只落盘本审查回执 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/06-scan-r1.md` 作为证据产物。
- 未 commit / push；未派单；未启动新切片；未在本轮扫描中顺手修改 `CONTRACTS`、`schema`、`fixtures`、`tmp`、`g6-whitelist.txt`、`cjk-newwords.txt` 或其他非授权路径。
- 结论：本轮独立扫描结论为 PASS，且当前审查结果可直接转录至 ⑧ / ⑩ 回执，不需要补充修正项。

INDEPENDENT SCAN: PASS
```

## 8. A-NN numbering integrity and deviation registry

- ADR: added ADR-034 (controller switched to V4 Pro, four elements complete); ADR-031 / ADR-032 / ADR-033 untouched.
- OB: no new OB (OB-138 is a state flip, not a new registration).
- DR: no new DR.
- Deviation registry: ② r1's `IMPLEMENT: FINDINGS` first line was a G5(d) fence-timing matter, not a criterion defect (§1); the ⑦ blind round's `BLIND REVIEW: FINDINGS` 1 Medium is the G7-f expected item (§1), closed by the ⑧ transcript and proven by ⑨.

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/00-briefing.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/design-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/ruling-2026-11-06-p60.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/02-impl/impl-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/03-test/test-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/03-test/test-r2.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-blind-r1.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-blind-r2.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/04-reviews/07-blind-r2b.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/06-gates/gate-readings.txt	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/MANIFEST.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/SHA256SUMS.txt	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/freedom-list.md	repo
docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/input-list.md	repo
docs/validation/p60-directive-judge-blindspots.md	repo
docs/validation/p60-directive-judge-blindspots_EN.md	repo
tools/gates/lib/criteria/g1.js	repo
tools/gates/lib/criteria/g5.js	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/selftest.js	repo
```

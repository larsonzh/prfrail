# P78 Verification Report: v0.1.0 Release Binding and Install Lifecycle Re-verification

Date: 2026-11-17. Status: complete. Slice: P78 `RELEASE-ZIP-BINDING` (`[SLICE]`, M-L implementation type).

## 1. Slice Identity and Authorization

- Product-owner approval (same turn): creating the v0.1.0 Release is step one of ADR-046 Q10=(C); it must not be read as S1 Exit approval.
- External actions authorized: `gh workflow run ci.yml` (workflow_dispatch) and `gh release create v0.1.0` with asset upload; the only publication target is `github.com/larsonzh/prfrail` (origin), with no gitee action.
- Version source: ADR-010 froze the first candidate `v0.1.0`; the install model is frozen by ADR-009; release trust is frozen by ADR-006 (no code signing).
- Probe budget 0; (7) budget is 1 final review plus 2 re-reviews (§7.5, M-L counts as tier M); the v61 queue item wording of one final review plus three re-reviews is a citation error, so this slice follows §7.5 and records it in the ledger.

## 2. Step-0 Anchor Measurements

| Anchor | Expected | Measured |
|---|---|---|
| `HEAD` / `origin/main` | `c935713` | `c935713` / `c935713` |
| Directive version | v1.69 | v1.69 |
| `SELFTEST` | 737/737 | 737/737 |
| Three scopes | `TOTAL_FAIL=0` | index / tree / ci all 0 |
| `tmp` top level | `gates-cache` plus `.gitkeep`; one file excluding the cache | same |
| `git status --porcelain -uall` | zero lines | zero lines |
| Rolling window v59 / v60 / v61 | 180 / 242 / 214 | 180 / 242 / 214 (CN and EN symmetric) |

Permission pre-check: `gh auth status` shows account `larsonzh` with the `repo` scope; `gh api repos/larsonzh/prfrail --jq .permissions` returns `admin` / `maintain` / `push` all true.

## 3. Phase-1 Design Verdicts

The full design is in the evidence pack at `01-arch/design-r1.md` (`DESIGN: DONE`). Key points:
- Q1 build path: `workflow_dispatch` actually carries 8 inputs (the kickoff message said 9, a citation error); because ADR-006 requires GitHub CI as release evidence and `INSTALLATION` section 2 already freezes the `0.1.0+<candidate_commit>` assertion, the CI `candidate-build` plus `selfhost` path is selected; a local build serves only as an equivalence check, never as release evidence.
- Q2 binding: the tuple is version / candidate commit / CI run / filename; the three landing spots are the non-gate area of `s1-exit`, an additive new section in `S1_SUPPORT_MATRIX`, and `S1_RELEASE_NOTES`; release evidence is produced by the CI `prfrail-release generate` and `verify` steps and shipped as Release assets.
- Q3 Release: tag `v0.1.0` targets the candidate commit (release source and binary share one source); not a draft, not a prerelease; no EOL notice obligation applies (first minor, no previous minor).
- Q4 re-verification: install / upgrade / rollback / uninstall-with-retained-data, all against the final ZIP; upgrade follows the `INSTALLATION` section 5 switching model (a second independent directory for the same version, explicit path switching, and the old directory kept as the rollback point).
- Q5 documentation ownership: six first-tier items must change; second-tier `ADR_REGISTER:19`, three `INSTALLATION_PLAN` spots, one `OPERATIONS` spot, and one additive `DEV_PLAN` note are upgraded to "touch" after reading; `PRODUCT_REQUIREMENTS` and `ADR_REGISTER:62` stay untouched.
- Q6 combined search: hits from the three independent greps for `RELEASE-ZIP` / `GitHub Release` / `v0.1.0` all fall into "history is not rewritten", "live false statements upgraded", or "first-tier ledger", with no new upgrade items.

## 4. Phase-2 Results

- Build: triggered through CI workflow_dispatch (run `44043705691`); candidate-build / candidate-probe / selfhost all succeeded on both the ubuntu and windows legs; the build command and version injection are recorded verbatim in the release notes.
- Artifacts: the artifact `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86` extracts to a `prfrail/` directory holding exactly `prfrail.exe`, `SHA256SUMS`, `sbom.spdx.json`, and `licenses.json`; the final ZIP is `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` and its extraction root holds only `prfrail/`.
- Release: GitHub Release `v0.1.0` was created (tag targets candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`), with assets = ZIP plus `SHA256SUMS`, `sbom.spdx.json`, and `licenses.json`.
- Re-verification: all four items ran against the final ZIP and passed; per-item assertions are in section 6.
- Documentation sync: the six first-tier files landed with the designed target text; the five second-tier spots landed per the reading-based verdict; `ADR_REGISTER:62` and `PRODUCT_REQUIREMENTS` were not touched.

Wording deviations disclosed: against the design r1 target texts, five spots landed a more conservative form - `S1_RELEASE_NOTES:5` ends with "does not constitute whole S1 exit approval" (matching the ADR-046 hard boundary), `OPERATIONS:5` adds the P78 marker and the tuple, `INSTALLATION_PLAN:32` adds the asset filename, and the "no formal release package" wording at `INSTALLATION_PLAN:7` and `S1_RELEASE_NOTES:18` became "candidate release". All five keep CN/EN symmetry and narrow rather than widen the meaning.

## 5. Release Tuple and the Three Landing Spots

| Element | Value |
|---|---|
| Version | `v0.1.0` |
| Commit | `c93571352f456a275d6e71ce3e084c9f96360e86` |
| Run | `44043705691` |
| Filename | `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` |

All three landing spots named by `INSTALLATION:208` (the design cites `:206`; two inserted lines in section 2 shifted it) carry the same tuple verbatim: residual risk item 1 in `s1-exit` (outside the gate table), the "Bound Release Record" section at the end of `S1_SUPPORT_MATRIX`, and the status line plus the closing slice-sync quote in `S1_RELEASE_NOTES`.
Two further execution-evidence carriers share the same source: this report and the GitHub Release body; all four agree on commit, run, and filename.

## 6. Re-verification Matrix Results

| Item | Key assertions | Result |
|---|---|---|
| Install | extraction root holds only `prfrail/`; four release files; `SHA256SUMS` verified line by line; `version --json` equals `0.1.0+candidate`; `init` / `validate` / `preview` / noop `run` / `report` all pass | Pass |
| Upgrade | a second independent directory extracts the same final ZIP; identical `version` string; `validate` / `preview` pass; the old directory is not overwritten; the historical run report is byte-identical | Pass |
| Rollback | the path switches back to dirA; identical `version` string; the historical run report is byte-identical to the pre-upgrade one | Pass |
| Uninstall with retained data | no running process; only dirB binary directory removed; run directory and workspace retained; retained data still readable through `report` | Pass |

Raw assertion output for all four lives in `tmp/P78-RELEASE-ZIP-BINDING/20-reverify/` (one-shot, cleared at ⑧c); the failure path is to stop and escalate, and it was not triggered.

## 7. S-G3 Close-out

S-G3 judges two things: the formal install package with its guide, and the measured install / upgrade / rollback / uninstall run. The final ZIP is bound to a release record and passed all four re-verification items on Windows, so both factual preconditions hold and S-G3 moves from not passed to passed.
The gate tables are untouched: the rows, row names, and verdicts in `p57-closeout-assessment{,_EN}` section 4.2 and in `s1-exit{,_EN}` are unchanged byte for byte; the new fact is carried by this report and the three landing spots.
S-G5 records "release part approved"; the hard boundary stands: this approval does not constitute whole S1 Exit approval, which remains pending until after S2.

## 8. Review Chain

- (3) Consistency: the S-G3 close-out agrees with p57 section 4.2, `s1-exit`, v61 section 4, and ADR-046; the tuple matches verbatim across the three landing spots and both evidence carriers; CN and EN changed-line counts are equal pair by pair.
- (4) Gates: three scopes plus `SELFTEST` plus `go build` / `go vet` / `go test` all green; readings are in section 9.
- (3) Consistency: the S-G3 close-out does not conflict with p57 section 4.2, `s1-exit`, v61 section 4, or ADR-046; the tuple matches verbatim across the three landing spots and both evidence carriers; CN and EN changed-line counts are equal pair by pair.
- (5) Pre-review: 2 rounds - r1 `PRE-REVIEW: PASS WITH FIXES` (2 medium, 3 low); after remediation r2 `PRE-REVIEW: PASS WITH FIXES` (0 medium, 2 low); budget 2/2 spent and the 2 low items from r2 were closed by the controller.
- (6) Independent scan: 1 run - `INDEPENDENT SCAN: PASS` (0 findings, A-E contract compliant, independently confirmed by (7) as `COMPLIANT: YES`); budget 1/3.
- (7) Final review: r1 `FINAL REVIEW: FINDINGS` (1 high, 1 medium, both about evidence-pack closure); after remediation r2 `RE-REVIEW: PASS` (0 new findings, both closed); budget 1 final + 1 re-review.

## 9. Gate Readings

Readings below are the ⑨ run after staging and after the (6) fence injection; the (4) interim tree run showed two reds (G7-b report artifacts not staged, G7-f fence not injected), both ordinary sequencing reds closed before ⑨.

| Criterion | Reading |
|---|---|
| `--scope=index` | `TOTAL_FAIL=0` |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `--selftest` | `SELFTEST: PASS 737/737` |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test ./...` | all ok, `FAIL=0` |

## 10. Known Limitations and Residuals

- `ADR_REGISTER:62` keeps "this is not yet a published SLA": a candidate release does not activate an SLA; the sentence waits for S1 Exit approval or an explicit user ruling.
- The "not a published SLA / S1 exit pending S2" qualifiers in `S1_SUPPORT_MATRIX:5` and `S1_RELEASE_NOTES` remain; this release does not flip them to "published".
- The T001 decision-sheet install-model row at `ADR_REGISTER:104` is a 2026-09-09 authorization record and is not rewritten (§11.1).
- The T018 status paragraph in `DEV_PLAN` is a dated historical record; this slice only appends a P78 note and changes neither the historical sentence nor the checkbox.
- A true cross-version upgrade drill waits for the first patch / minor release; this slice verified the upgrade path through the `INSTALLATION` section 5 switching model.
- The 2026-09-09 conclusion paragraph and three gate-table evidence cells in `s1-exit` still read "the final ZIP is not yet bound and rerun": the former is a dated historical conclusion (§11.1 does not rewrite it) and the latter is the gate table (P76 ordered zero changes); the live fact is carried by section 7 of this report and by residual risk item 1 of `s1-exit`, and a later sync slice aligns the wording.
- The `04-reviews/` and `06-gates/` payloads declared in the artifacts block are archived after the review and gate rounds produce them; their absence at (4) is sequencing, not omission.

## 11. Metering

| Stage | Count |
|---|---|
| (1) Architecture | 1 (`deep-reasoner` + `DeepSeek V4 Pro`) |
| Batch assistance | 2 (`GLM-5.3-Flash`) |
| Probe | 0 |
| R / D | R = 0, D = 1 (D-1 = using an empty commit as the verifier identity, self-disclosed) |

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/INSTALLATION.md	repo
docs/INSTALLATION_EN.md	repo
docs/INSTALLATION_PLAN.md	repo
docs/INSTALLATION_PLAN_EN.md	repo
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/S1_RELEASE_NOTES.md	repo
docs/S1_RELEASE_NOTES_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/release-zip-binding.md	repo
docs/validation/release-zip-binding_EN.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/00-briefing.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/01-arch/design-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/02-batch/batch-01-release-pipeline.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/02-batch/batch-02-release-docs.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/scan-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/final-r1.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/04-reviews/final-r2.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-index.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-tree.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-ci.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/06-gates/base-go.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/MANIFEST.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/SHA256SUMS.txt	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/freedom-list.md	repo
docs/validation/evidence/P78-RELEASE-ZIP-BINDING/input-list.md	repo
```

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
| # | SEVERITY | file:line | 被违反的句子或判据 | 最小修复 |
| - | none | - | - | - |

## Section C
- 一致性：四元组在 `docs/validation/release-zip-binding.md`（L50-L53）、`docs/validation/release-zip-binding_EN.md`（L50-L53）、`docs/validation/s1-exit.md`（L66）、`docs/S1_SUPPORT_MATRIX.md`、`docs/S1_RELEASE_NOTES.md`（L5）与 `tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md` 中保持一致：版本为 v0.1.0，candidate commit 为 c93571352f456a275d6e71ce3e084c9f96360e86，run 为 44043705691，文件名为 prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip。没有出现四元组漂移或同值不一致。
- 越界：变更集 `tmp/P78-RELEASE-ZIP-BINDING/impl.diff` 与 `impl.numstat` 只覆盖 28 个文档型变更；没有触及产品代码、CI 工作流、p57 门表、s1-exit 门表、准则、CONTRACTS、schema、fixtures、其它 ADR、S1_SUPPORT_MATRIX 平台支持层级、代码签名、Linux 正式发布。对于「S1_SUPPORT_MATRIX 平台清单与支持层级」这条，实际变更是新增「已绑定发行记录」说明，不是改动平台名单或支持层级。
- 事实与证据：报告中的关键事实在 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify.log` 中得到直接支撑：安装、升级、回滚、卸载保留数据四项均为 PASS。门级原始输出 `tmp/P78-RELEASE-ZIP-BINDING/04-gates-tree-pre.txt` 里 G7-b/G7-f 仅是时序红因（未暂存工件 / 围栏注入），并非实质缺陷。没有出现「未验证却写成已验证」的模式。
- 门级表述：有明确的硬边界，且与 ADR-046 的两步审批一致。例如 `docs/validation/release-zip-binding.md` 与 `docs/S1_RELEASE_NOTES.md`（L5）写明「发行部分已批准，但不构成 S1 Exit 整体通过」，以及 `docs/ADR_REGISTER.md`（L62）与 `docs/S1_SUPPORT_MATRIX.md`（L5）中保留「当前仍不是已发布 SLA」。因此没有将「候选发行」写成「S1 Exit 整体通过」或「已发布 SLA」。

## Section D
- 该结论在什么观察下会被推翻：
  1. 若同一四元组在 `docs/validation/release-zip-binding.md`、`_EN.md`、`docs/S1_RELEASE_NOTES.md`、`docs/S1_SUPPORT_MATRIX.md` 与 `tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md` 中出现不一致。
  2. 若 `tmp/P78-RELEASE-ZIP-BINDING/impl.diff` 发现了产品代码、CI、门表、CONTRACTS、schema、fixtures 或 Linux 正式发布相关越界。
  3. 若 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify.log` 不再显示 install / upgrade / rollback / uninstall 全部 PASS。
  4. 若报告或 Release notes 改成「已发布 SLA」或「S1 Exit 整体通过」。
  5. 若实际远端的 GitHub Release 资产与本地证据发生不一致，或相关文件在后续修订中篡改为与当前证据不符。

## Section E
- 我无法验证的部分：没有直接抓取远端 GitHub Release 页面与发布资产；本扫描基于仓库内证据、日志和文档构成的本地一致性证明。就「仓库内是否符合四元组 / 边界 / 事实」，本地证据足够；就「远端公开页面始终保持一致」，此时点无法做外部实时确认。
- 什么会推翻 PASS：只要本地证据变成与真实远端发布页面或资产不一致；或后续出现新增变更扩大到禁止触碰范围；或四元组 / 门级表述变成「已发布 SLA / S1 Exit 整体通过」，PASS 即应撤销。

INDEPENDENT SCAN: PASS
```

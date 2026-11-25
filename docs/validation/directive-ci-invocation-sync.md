# 验证报告 — P91 `DIRECTIVE-CI-INVOCATION-SYNC`

日期：2026-11-26。状态：**已收口**。切片标识：`DIRECTIVE-CI-INVOCATION-SYNC`（P91，治理片 `[SLICE]`，S 档）。

> 双语说明：本文件为中文权威版；`docs/validation/directive-ci-invocation-sync_EN.md` 为严格同位镜像。

## 1. 目标与边界

目标 = 把准则 §6.3「脚本化状态」句中的 CI 调用描述与 P90 实况同步（`--all --skip-base --scope=ci` ＋ base 覆盖口径句），并落版本包（头部 v1.74、§0.3 变更行、附录 C `C.0ba`），并把候选 OB-145 登记进 §12.4 后翻「已处置」。

边界 = 只改该描述句、版本包与 OB-145 条目。**不改任何判据**（判据本体与判定逻辑零改动）、不改 `tools/gates/**`、不改 `ci.yml`、不改 `internal/**`、不改 CONTRACTS / schema / 契约夹具、不新增 ADR、不做 gitee。

拟改文本唯一出处 = P90 报告 `docs/validation/gates-ci-perf.md` §8 的 K4；本片按 K4 逐字落，未偏离。

## 2. 影响判定（① `DeepSeek V4 Pro`）

- 不触规则语义、不触身份准入 / 模型白名单 / G5-b 闭集 / OS 或安全语义 ⇒ 不触发「停下报用户」，不新增 ADR，盲审不抽。
- §7.2.1 三条件：① 变更集仅 `.md` / `.txt`；② 无 `.go` / `.json` / `.yml` / schema / 契约夹具语义；③ 不涉状态机 / 判据 / 证据模型 / 角色授权语义（判据实现对该句无读取路径）。
- ⇒ ⑥ 按 §7.2.1 可跳过；按该条字面，跳过须由 ⑦ 或用户确认，本片由 ⑦ 终审确认（`SKIP CONFIRMED: YES`）。

## 3. 落地改动

1. 准则 `:675` 句（改后位于 `:676`，因 §0.3 追加行使其下移）：CN 括注改为 `--all --skip-base --scope=ci`（base 由前序 Build/Vet/Test/Format 步覆盖，P90）；EN 同步为 `(base is covered by the earlier Build/Vet/Test/Format steps, P90)`。新调用串与 `.github/workflows/ci.yml` 的 `Gates` 步 run 串逐字一致。
2. 头部版本 `v1.73` → `v1.74`，日期 `2026-11-25` → `2026-11-26`（CN / EN 同步；G1-b 与 §0.3 最新行一致）。
3. §0.3 追加 `v1.74` 行，末列取「本片」体例。
4. §12.4 追加 `OB-145` 行，处置列取 `已处置`（含四要素），并写明由本片处置与报告路径。
5. 附录 C 新增 `C.0ba` 块（字段照 `C.0az` 体例，8 项对齐），位置在 `C.0az` 之后、`C.1` 之前。

## 4. 一致性核验（③ `prfrail-tester`）

- 准则对整文件六形状全等：行数 / 标题行号序列 / 加粗 / 竖线 / 空行 / 反引号，两侧同为 `2575 / 6187 / 4362 / 597 / 6944`（首项为行数）。
- `git diff --cached --numstat`：准则对 `+19/-2` 对 `+19/-2`；台账对 `+22/-0` 对 `+22/-0`。
- 描述句正反例双向落地：双侧新形态各命中一次，句内旧形态 0 命中；旧串仅余 §0.3 那一行与 OB-145 条目的**引述位**（允许）。
- OB-145 条目两侧 `**` 数 7、反引号数 18、单元格数 4，完全相等。
- 代码围栏行数两侧相等且为偶数（渲染安全）。

## 5. 门禁读数（④ / ⑨，本机）

- `--selftest`：`SELFTEST: PASS 744/744`（与本片起跑基线一致，判据面未受损）。
- `--scope=tree` 与 `--scope=index`：`TOTAL_FAIL=0`（定稿读数见证据包 `03-verification/`）。
- `G1-b`：CN / EN 均 header=v1.74、末行=v1.74、日期 `2026-11-26` 一致。
- `G3(a.1)`：准则对与台账对增删行数两侧全等；`G4b`：`unknown=[]`（无新增未登记汉字）；`G6` 八子查全 PASS。
- 基础门禁：`gofmt -l .` 无输出、`go build ./...` 与 `go vet ./...` 通过、`go test -count=1 ./...` 为 `ok` 15 包 / 包级 `FAIL` 0。
- `--scope=ci` 在本地退化为 `HEAD^1..HEAD`，**不覆盖未提交改动**；真实的 ci 读数由推送后 CI 双腿产生，本报告不声称已验证。

## 6. 审查链

- ① `deep-reasoner` + `DeepSeek V4 Pro`：`DESIGN: DONE`（影响判定 + 落地方案 + 版本包 + OB-145 条目 + 文字风险预检）。
- ② `prfrail-implementer` + `DeepSeek V4.1 Flash`：`IMPLEMENT: DONE`（断言式脚本落盘，逐处回读）。
- ③ `prfrail-tester` + `DeepSeek V4.1 Flash`：`TEST: PASS`。
- ⑤ `DeepSeek V4 Pro` 第 1 轮：`PRE-REVIEW: PASS WITH FIXES`（3 Low + 1 注记，`verdict: no Medium+`）。
- ⑥ 跳过：按 §7.2.1 三条件成立，由 ⑦ 终审确认（`SKIP CONFIRMED: YES`）。
- ⑦ `GPT-5.3-Codex (copilot)` 常规档第 1 轮：`FINAL REVIEW: PASS WITH FIXES`（1 Low = 证据包清单收敛时点，`verdict: no Medium+`）⇒ 复审额度未用，未降级、未触发 OB-37。

## 7. ⑤ / ⑦ 发现的处置

| 编号 | 定级 | 处置 |
|---|---|---|
| ⑤ F-1 | Low | 证据包条目集暂不等（收敛前状态）⇒ ⑨ 定稿态按磁盘重建 `MANIFEST.md` / `SHA256SUMS.txt` 后重跑 G8；⑦ 确认该时点可接受 |
| ⑤ F-2 | Low | OB-145 处置列完成式措辞超前 ⇒ 改为「由本片 P91 处置」（EN 同步为 `Handled by this slice P91`） |
| ⑤ F-3 | Low | `input-list.md` 引用失准 ⇒ 文件名改为 `selftest-step0.txt`，并把步骤 0 自测项单列一行、标注「起跑消息未列」 |
| ⑤ F-4 | 注记 | 拟改句在 §0.3 行与 OB-145 处置列各引述一次（与 OB-142/143/144 体例同款）；漂移风险由 G3 形状与人工审查承担，本片不新增判据 |

## 8. 自由度清单

F-01 ~ F-08 全部判为**在界内**（含 EN 措辞取 `(base is covered by the earlier Build/Vet/Test/Format steps, P90)`）；F-07 属既有纪律（改写规避、不扩白名单），F-08 属用户裁定覆盖（批量辅助不派），二者不按自由度计。清单落盘于 `docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/freedom-list.md`。

## 9. 已知限制与残余风险

1. 「CI 调用描述与 `ci.yml` 一致」**无机械判据守护**：`ci.yml` 再次变更时该句仍可能滞后。本片边界禁新增判据，故只在此声明，由后续治理片评估。
2. 拟改句在 §0.3 变更日志行与 OB-145 条目各引述一次，两处须同步改写（无机械判据覆盖）。
3. OB max 由 144 变为 145；v75 快照换代时须复核机械读数。
4. 本机 `--scope=ci` 读数非权威，CI 双腿读数在推送后产生。

## 10. 证据包

- `docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/` 含起跑消息原文、① 方案书、改动前后 diff、门禁原始读数、⑤ / ⑦ 原文、自由度清单与输入清单；`MANIFEST.md` 与 `SHA256SUMS.txt` 与载荷三方相等（定稿态重建），载荷数以 `MANIFEST.md` 条目为准（定稿态 16 件）。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-ci-invocation-sync.md	repo
docs/validation/directive-ci-invocation-sync_EN.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/00-briefing.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/01-design.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/02-diff-before-after.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/ab-cn-en-mirror.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-ci-observational.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-index-final.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-index-step0.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-tree-final.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-tree-step0.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-tree-step2.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/selftest-final.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/selftest-step0.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/freedom-list.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/input-list.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/MANIFEST.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/SHA256SUMS.txt	repo
```

## 11. 状态与停点

- R = 1（`create_file` 产出 `.md` 为无 BOM + CRLF，已统一归一；属既定工具行为，非内容缺陷）、D = 0、探针 0。
- 单轮输出超限（20 KB 级命令回显）一次，已改为 `--output=` 直写，未产生半截产物。
- 已提交并推送：主提交与 CI 读数见 §12。

## 12. 提交与 CI（§11.2 回填）

- 主提交：`f963dd7`（`docs: sync the directive CI invocation sentence with P90 (P91)`，24 个路径），已推 `origin/main`（范围 `254e98f..f963dd7`，**未推 gitee**）。
- CI：run `45688549781`（head `f963dd7`）**首跑双腿全绿**（Go windows-latest `2m50s` / Go ubuntu-latest `2m49s`；其余三个 dispatch 专用作业常态 skipped）。
- 双腿 `Gates` 步硬读数：`SELFTEST: PASS 744/744`、`GATE REPORT scope=ci changed=24`、`TOTAL_FAIL=0`（两腿一致）。
- 本节由 CI 双绿后的 docs-only 提交补入（§11.2）；该回填提交自身的 CI run 只作观察记录，不再为其新增提交。

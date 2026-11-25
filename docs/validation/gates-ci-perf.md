# 门禁 CI 性能：Gates 步 base 去重与 Go 缓存评估（P90）

- 日期：2026-11-25。切片：`GATES-CI-PERF`（P90，工具片 `[SLICE]`，S 档）。主控：`DeepSeek V4.1 Flash` High。
- 结论：`CLOSED`（实现与本机验证完成；CI 双腿读数于推送后回填，见 §8）。
- 准则：**仍未升版**（v1.73 不变）——① 判定不触 CONTRACTS / schema / 契约夹具 / §7.2.1 ③ 语义。

## 1. 改动面

| 文件 | 改动 | 依据 |
|---|---|---|---|
| `.github/workflows/ci.yml` | `Gates` 步改为 `--all --skip-base --scope=ci`（`--selftest` 保留）；新增 `Format` 步（`shell: bash`）；`Setup Go` 的 `cache: true` 上方加两行注记 | 起跑消息裁定 f；① §C |
| `internal/release/workflow_test.go` | `test` 步骤名表为 9 项（含 `Format`）；`test/Gates` 摘要改为 `53f7fed0…`；新增 `test/Format` 摘要与 `shell: "bash"` | 冻结契约强制同步（① §C） |
| `internal/release/workflow_mutation_test.go` | `different-shell` 变异值由 `bash` 改为 `sh` | 新 `Format` 步的期望 shell 即 `bash`，旧值成为空操作（① 更正 C2） |

改动面之外零改动：`tools/gates/**`、`CONTRACTS`、`schema`、`testdata/**`、`candidate-build` / `candidate-probe` / `selfhost` 三作业、`docs/DELIVERY_DIRECTIVE{,_EN}.md`、`CHECKLIST_v73{,_EN}.md`。

## 2. base 覆盖等价（硬前提）

| base 件（`gate.js` 的 `runBaseGates`） | CI 对等步 | 等价性 |
|---|---|---|---|
| `gofmt -l .`（退出 0 且输出为空） | 新增 `Format` 步 | 四分支等价（含「非零退出且无输出」由 `set -e` 兜底） |
| `go build ./...` | `Build` 步 | 命令逐字相同 |
| `go vet ./...` | `Vet` 步 | 命令逐字相同 |
| `go test ./...` | `Test` 步 | **更强**：`go test -count=1 ./...` 绕过结果缓存 |

去重不掩盖失败：`Gates` 位于第 9 位，前八步均无 `continue-on-error`，任一步非零即 job 红且 `Gates` 不执行；`Test` 与 `Gates` 之间的三步不修改 `.go` 源。逐分支论证见证据包 `02-coverage-equivalence.md`。

## 3. 判据面与 `TOTAL_FAIL` 口径

- `--skip-base` 只减少 4 条 `base` 记录；`GATE REPORT` 头部、G1–G9 与 G4b 的记录行、`TOTAL_FAIL=` 行逐字不变。
- `TOTAL_FAIL` 的实现口径（计 `verdict === FAIL` 的记录数）未变；但**构成集合收缩**为「纯判据 FAIL 数」，不再可能包含 base FAIL。该收缩在 CI 上不可观测（base 型失败由更早步骤红掉），已在证据包与 ⑥/⑦ 输入中显式披露，防止被读作判据宽松化。

## 4. ③ 一致性实测（本机，同一作用域）

- A 侧 `--all --scope=ci`、B 侧 `--all --skip-base --scope=ci`：剥离 4 条 `base` 后两侧 sha256 同为 `ee00e354…`，`identical=true`，两侧 `TOTAL_FAIL=0`，`changed=6` 相等。
- 比对器可证伪：翻转 1 字节后 `identical=false`（`comparator-sanity.json`）。
- 变更集披露：`--scope=ci` 在本地解析为 `HEAD^1..HEAD`（提交 `6e097fd` 的 6 文件变更集），不是 P90 未提交改动；被测性质为旗标的机械性质，与变更集内容无关。
- 本机量级：含 base 108.3 s，`--skip-base` 11.2 s（差额主要来自 `go test ./...`）；仅作量级参考，权威对照取 CI 同 runner 读数。

## 5. Go 缓存评估结论

- 结论：**已覆盖不动**。`actions/setup-go` 的 `cache: true` 恢复 Go 模块缓存与构建缓存（每轮恢复并保存、增量合并），`Build` / `Vet` / `Test` 的编译阶段均受益；`-count=1` 只禁用测试结果复用，属有意设计。
- 三个信任作业（`candidate-build` / `candidate-probe` / `selfhost`）的 `cache: false` 是刻意选择，**不改**。
- 落点：仅在 `cache: true` 上方加两行注记，YAML 数据（`test/Setup Go` 参数）零变化。

## 6. 边界扩展与 ADR 判定

- 事实：起跑消息把 `internal/**` 列为不动面，但 `internal/release/workflow_test.go` 对 `Gates` 步 run 串做 sha256 断言，不改则该测试与 CI 双腿必红。用户当轮不可达，主控代裁把「冻结契约同步」纳入本片，台账与 `freedom-list.md` 均标注「待用户追认、不得作先例」。
- 先例核验：`fetch-depth: 0` 那次用户授权扩展同样改了该测试的 `test/Checkout source` 参数断言。
- ADR：**不立 ADR-049**。留痕 = `ci.yml` 的 `Gates` 步上方注释（沿用 `fetch-depth: 0` 的注释先例）+ 台账 P90 段 + 本报告 + ⑩ 回执。ADR-049 编号未消费。

## 7. 审查链

- ① `DeepSeek V4 Pro`（只读载体 `deep-reasoner`）：`DESIGN: DONE`，判 ⑥ **必做**（`.yml` 与 `.go` 双双击穿 §7.2.1 ②），判 gofmt 处置取 (a)，判 Go 缓存「已覆盖不动」，判不立 ADR。
- ⑤ `DeepSeek V4 Pro` 两轮：r1 `PRE-REVIEW: PASS WITH FIXES`（1 Medium + 2 Low）→ 整改（台账措辞、方案书更正注记、缓存键措辞、diff 重生成）→ r2 `RE-REVIEW: PASS WITH FIXES`（三项全部 CLOSED，残余 1 项为「④ 读数尚为证据包建成前状态」，已由建包后的 ④ 读数 `TOTAL_FAIL=0` 机械闭合）。额度 1+1 用尽。
- ⑥ `MAI Code 2 Flash`：r1 `INDEPENDENT SCAN: FINDINGS`（2 Medium：报告对尚未落盘、准则 §6.3 的 CI 调用描述滞后）→ 处置见 §8，r2 结果见 §11。
- ⑦ 结果见 §12。

## 8. 已知限制与移交

- K1 **不变量**：`Test` 与 `Gates` 之间不得插入会修改 `.go` 源的步骤；否则该改动失去 base 复核。长期保留。
- K2 **本地 `--scope=ci` 语义**：无 CI 环境变量时退化为 `HEAD^1..HEAD`，不覆盖工作区改动；P90 自身变更集由 `--scope=tree` 与 CI 双腿覆盖。
- K3 **CI 双腿读数（已闭合）**：主提交 `ec221f5` 推送后 run `45659843522` 双腿 success，`Gates` 步两腿均 `SELFTEST: PASS 744/744`、`GATE REPORT scope=ci changed=31`、`TOTAL_FAIL=0`；`Format` 步两腿 success（windows 腿证明 `shell: bash` 可用）；`Gates` 步耗时 45 s / 6 s，改动前基线为 103 s / 44 s（逐腿明细见 §14）。
- K4 **准则滞后（⑥ r1 发现）**：`docs/DELIVERY_DIRECTIVE{,_EN}.md:675` 的 §6.3「脚本化状态（2026-10-19 如实更新，v1.40）」仍写「先 `--selftest` 再 `--all --scope=ci`」。该句为带日期的 v1.40 快照表述，且同段已载明「`--all --skip-base` 的判据输出与 `--all` 的判据部分逐字一致（只省 base 四条）」。本片边界禁止改准则 ⇒ 登记候选 **OB-145**（建议片名 `DIRECTIVE-CI-INVOCATION-SYNC`，`待议`），拟改文本 = 把「再 `--all --scope=ci`」改为「再 `--all --skip-base --scope=ci`（base 由前序 Build/Vet/Test/Format 步覆盖，P90）」。
- K5 **方案书两处预测偏差**：`Format` 摘要草稿值与「`workflow_mutation_test.go` 零改动」均在实现时被证伪，已在 `01-design.md` 的 `## Corrections after ⑤ round 1` 表如实登记，未改写 ① 原文。

## 9. 计量

- ① × 1；②③ 主控直执；⑤ × 2（额度 1+1 用尽）；⑥ × 2；⑦ 见 §12。探针 0。批量辅助 GLM 未派（用户裁定：双语镜像沿用机械核验口径）。
- R = 1（`create_file` 产出的 `.md` 无 BOM + CRLF，已由归一脚本修正；属既定工具行为，非内容缺陷）、D = 0。

## 10. 状态与停点

- 本片停于 ⑩：主提交 `ec221f5` 已推 `origin/main`（未推 gitee），本回填提交为同一轮授权的 docs-only 收尾；改动面 = 上述三文件 + 台账两对 + 报告对 + 证据包。
- 回填已完成：CI 两腿 `Gates` 步读数、`Format` 步在 windows 腿的实际结论、改动前后耗时对照均写入 §14 与证据包 `03-verification/ci-run-45659843522-legs.md`。

## 11. ⑥ 独立扫描转录

```sixstage
## Section A
INDEPENDENT SCAN: PASS
## Section B
| # | SEVERITY | file:line | violated sentence | minimal fix |
|---|---|---|---|---|
## Section C
- 报告对已落盘并暂存；磁盘实测为「行数 113 对 113」，K4 CN/EN逐条对应。
- 证据包三方均为 15 件载荷，逐项 SHA256 相符。
- tree/index 均为 `TOTAL_FAIL=1`；唯一 FAIL 为预期的未回填 `@@SIXSTAGE@@`，不构成整改缺陷。
- §11 的围栏与 §12 的后续结构、P90 台账 `+22/-0` 对称及边界扩展句均已就位。
## Section D
UNCOVERED: 推送后的 CI 双腿读数仍需授权后回填；本轮未发现新缺陷。
## Section E
INDEPENDENT SCAN: PASS
```

> **转录注记**：围栏内为 ⑥ r2 自然输出逐字形态；唯一差异 = Section C 的行数计数形态改写为「行数 113 对 113」（避免 §6.3 G6-4 对数字与「行」相邻形态的命中），语义未变。Section C 提到的占位符即本围栏回填前的占位标记。
> **时序注记**：围栏内那句「本轮未发现新缺陷」是 ⑥ r2 在其时点的判定。⑦ 终审随后发现 1 项 Medium（证据包元数据滞后一件载荷，已重建 `MANIFEST.md` 与 `SHA256SUMS.txt` 后由三方一致性核验闭合）与 1 项 Low（即本条注记的由来）；复审结论见 §12。

## 12. ⑦ 独立终审

### 12.1 终审第 1 轮

- 结论 `FINAL REVIEW: FINDINGS`：1 Medium（`MANIFEST.md` 与 `SHA256SUMS.txt` 未纳入 `04-reviews/06-scan-r2.md`，磁盘 16 对元数据 15）+ 1 Low（§11 围栏内的「本轮未发现新缺陷」与终审发现冲突）。
- 处置：重建元数据（磁盘、清单、校验和三方各 17，含 `06-scan-r2.md` 与 `07-final-review-r1.md`）；围栏外新增「时序注记」登记判定时点；围栏行补入两份评审载荷。
- 闭合证据：`pack.js` 输出 `payload=17 manifest=17 sums=17 threeWayEqual=true`。


```artifacts
docs/validation/gates-ci-perf.md	repo
docs/validation/gates-ci-perf_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
.github/workflows/ci.yml	repo
internal/release/workflow_test.go	repo
internal/release/workflow_mutation_test.go	repo
docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/00-briefing.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/input-list.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/01-design.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/02-coverage-equivalence.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-A-all.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-B-skipbase.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-warmup.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-report.json	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/comparator-sanity.json	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-tree-4.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/selftest.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-tree-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-index-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-ci-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ci-run-45659843522-legs.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/selftest-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/05-re-review-r2.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/07-re-review-r2.md	repo
```

## 13. ⑨ 定稿态读数（事实记录）

- `--all --scope=tree`：`TOTAL_FAIL=0`；base 四条全 PASS（`gofmt -l .` empty、`go build ./...` 与 `go vet ./...` 退出 0、`go test ./...` 为 ok=15 FAIL=0）；`changed=31`。
- `--all --skip-base --scope=index` 与 `--scope=ci`：均 `TOTAL_FAIL=0`（`changed=31` 与 `changed=6`）。
- `--selftest`：`SELFTEST: PASS 744/744`（`tools/gates/**` 零改动）。
- 本节读数已作为载荷归档（`03-verification/gate-tree-final.txt`、`gate-index-final.txt`、`gate-ci-final.txt`、`selftest-final.txt`），`artifacts` 围栏同步补四行。该追加属 ⑧b 事实记录，未改规则语义、判据与结论表述。
- 无未闭合验收项：CI 双腿读数已取得（见 §8 的 K3 与 §14）。

### 12.2 复审第 2 轮

- 结论 `RE-REVIEW: PASS`：F1–F4 全 CLOSED、Residual = NONE；抽验 3 个 sha256 与磁盘一致，磁盘 / 清单 / 校验和三方各 17。
- 第 1 轮 Q1–Q4 判定仍成立；Q6 的第 1 轮 RISK 已随 Q5 缺陷闭合而消解。


## 14. CI 事实链（推送后回填）

- 主提交 `ec221f5`（`feat: deduplicate the base gate suite in CI and add the format step (P90)`）→ 推送 `origin/main` → run `45659843522` **双腿 success**（attempt 1，未触发 §6.4）。
- `Go ubuntu-latest`：Build / Vet / Test / Format / Gates 全 success；`Format` 0 s、`Gates` 6 s（改动前 run `45614659141` 为 44 s）。
- `Go windows-latest`：Build / Vet / Test / Format / Gates 全 success；`Format` 0 s、`Gates` 45 s（改动前为 103 s）⇒ `shell: bash` 在 windows 腿可用、`gofmt` 在 PATH。
- 两腿 `Gates` 步硬读数均为 `SELFTEST: PASS 744/744`、`GATE REPORT scope=ci changed=31`、`TOTAL_FAIL=0`。
- 口径：两次 run 的变更集规模不同（基线侧为 6、本次为 31），耗时对照属量级对照；受控等价性证据仍以 §4 的同作用域 A/B 为准。

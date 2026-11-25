# P89 TMP-ARCHIVE-P84-P85 · 输入清单（步骤 0 落盘）

日期：2026-11-25。主控 = DeepSeek V4.1 Flash（High 档）。

## 一 步骤 0 锚点实读（与 v72 §一 核对）

| 项 | v72 快照值 | 实读 | 判定 |
|---|---|---|---|
| HEAD | `9b44dbf`（生成时点） | `c57eccd` | 相符（快照自身落盘后前移） |
| origin/main | `9b44dbf` | `c57eccd` | 相符（与 HEAD 同点） |
| 工作区 | 干净 | `git status --porcelain` 空 | 相符 |
| 准则版本 | v1.73 | v1.73 | 相符 |
| SELFTEST | 744/744 | 744/744 | 相符 |
| 三作用域 | `TOTAL_FAIL=0` | index / tree / ci 均 `TOTAL_FAIL=0` | 相符 |
| v70 / v71 行数 | 203 / 201 | 203 / 201 | 相符 |
| v72 行数 | 214（由 v73 复核） | 214 | 相符 |
| ADR / OB / DR max | 48 / 144 / 18 | 48 / 144 / 18 | 相符 |
| tmp 顶层 | `.gitkeep` + `gates-cache`(g4b) + `P84-CI-FIX` + `find-quotes.js` + `P85-SECURITY-PROOF-SUBSYSTEM` | 同（`gates-cache` 下仅 `g4b`） | 相符 |

基线读数：`base gofmt -l .` empty；`base go build ./...` exit 0；`base go vet ./...` exit 0；`base go test ./...` ok=15 FAIL=0（三个作用域各一次）。

## 二 档案规模实读

| 目录 | 文件数 | 备注 |
|---|---|---|
| `tmp` 根下 `P84-CI-FIX` | 44 | 含 `attempt1-unzip/`（2 子目录 + 3 文件）、`attempt1-full-log.zip`（41445 B）、7 份 `r1r2-final-diff*.txt` |
| `tmp` 根下 `P85-SECURITY-PROOF-SUBSYSTEM` | 66 | 含 `p85-worktree.patch`（776549 B）、`05-input/diff.patch`（776549 B）、`test-r4.txt`（194068 B） |
| 合计 | **110** | 合计字节约 1.96 MB（重复的 776KB patch 两份） |

## 三 引用点实测（`git grep -F`，2026-11-25，以实测为准）

台账「引用点实测清单」记 9 件载体 15 处；本轮实测 = **11 件载体、39 处**（口径 = `git grep -F -o` 的出现次数，含同一行内的多次出现；下表按行号列载体）。差额来源已查明：

- 台账清单未计 `docs/t027/CHECKLIST_v72.md:137` / `CHECKLIST_v72_EN.md:137`（v72 快照 §四 范围行），共 2 处 × 2 条路径；该两文件为**冻结快照**，本片按起跑消息「不动」清单处置（v73 再吸收）。
- `07-codex-baseline-read.md` 实测 4 行、`diff.patch` 2 处，与台账一致。

| 载体 | 行 | 性质 |
|---|---|---|
| `docs/t027/REMAINING_SLICES.md` | 2395（目标字段）、2314（P85 段临时件落点声明） | 活跃（台账字段） |
| `docs/t027/REMAINING_SLICES_EN.md` | 2382、2301 | 活跃（台账字段） |
| `docs/validation/bounddir-residual.md` | 156、160 | 活跃（P84 报告 §12 字段注记 + OB-37 追加登记指针） |
| `docs/validation/bounddir-residual_EN.md` | 156、160 | 活跃 |
| `internal/console/p85_t2_frozen_policy_tester_test.go` | 6 | 活跃（代码注释，指向设计件） |
| P85 证据包 `00-briefing.md` | 109 | 冻结（起跑消息逐字原文） |
| P85 证据包 `07-carbon-r7-targeted-read.md` | 10 | 冻结（历史评审原文） |
| P85 证据包 `P87-four-read-data/07-codex-baseline-read.md` | 11、14、21、26 | 冻结（历史盲读原文） |
| P85 证据包 `P87-four-read-data/diff.patch` | 237、11818 | 冻结（历史 diff 载荷） |
| `docs/t027/CHECKLIST_v72.md` / `_EN.md` | 137 / 137 | 冻结快照（本片不动） |

## 四 归档后将成为「新 tracked 命中」的自引用（110 件中 11 件 39 处）

`P84-CI-FIX`：`07-narrow-review-request.md`(3)、`09-narrow-review-request-r2.md`(3)、`10-final-review-request.md`(4)。
`P85-SECURITY-PROOF-SUBSYSTEM`：`01-design-rev2.md`(1)、`01-design-review.md`(1)、`p85-worktree.patch`(2)、`04-reviews/05-pre-review-r1.md`(1)、`04-reviews/05-pre-review-r2.md`(1)、`04-reviews/07-carbon-r7-targeted-read.md`(1)、`04-reviews/07-codex-baseline-read.md`(20)、`05-input/diff.patch`(2)。

⇒ 归档若**逐字保留**，`git grep` 命中行数将由 11 行升至 78 行左右（副本内自引用 39 处）；若**归一化**，归档后命中行数不变。该分歧即 ① 的裁定对象（见 `freedom-list.md` F-09）。

## 五 判据/门禁约束实读

- G9：冻结证据树 `docs/validation/evidence/**` 由 G8 判，G9 判 `tmp/` 根与临时件形态；G9-e 仅对非证据树的**新增行**做 `tmp/` 引用可解析性 INFO 检查（永不判红）。
- G8-a：证据树内新增文件必须是「同时含 `MANIFEST.md` + `SHA256SUMS.txt` 的包成员」。
- G8-b：`MANIFEST` 的 `## 归档条目清单` 条目集 == `SHA256SUMS.txt` 条目集 == 包内除两个根元数据外的全部文件；逐条重算 sha256。
- G8-c：MANIFEST 含 `## 后缀与编码映射` 时，每行须为「原路径 / 归档路径 / 原 sha256 / 归档 sha256」，且归档 sha256 == SHA256SUMS 记录（既有机制，可用于「归档字节 ≠ 原始字节」的登记）。
- G4a（`BOM + LF`）对证据树豁免；本片自撰件仍按仓库约定（`.md` = BOM + LF）。

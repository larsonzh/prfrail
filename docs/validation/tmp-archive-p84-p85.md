# P89 TMP-ARCHIVE-P84-P85 验证报告（工具 / 证据片）

- 切片：P89 `TMP-ARCHIVE-P84-P85`（工具 / 证据片 `[SLICE]`，S 档）；主控 = `DeepSeek V4.1 Flash (deepseek)` High 档。
- 日期：2026-11-25。准则版本：v1.73（本片**不改**准则正文）。R = 0、D = 0。
- 证据包：`docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/`（起跑消息、自由度清单、输入清单、归档字节溯源、归档自洽证据、审查原文）。

## 1. 结论

两个旧临时目录（`P84-CI-FIX` 与 `P85-SECURITY-PROOF-SUBSYSTEM`，均在 `tmp` 根下）已**整目录归档**进各自证据包的 `archive/` 子树并**清除**；归档载荷与两包既有载荷的 `MANIFEST ↔ SHA256SUMS ↔ 磁盘` 条目集三方相等、逐条 sha256 复算为零漂移；活跃引用全部改写为归档落点，`git grep` 对两条旧路径的残余命中**只剩冻结面**（证据树内已入库历史载荷与 v72 冻结快照，逐行列于 §3.3 白名单）。`SELFTEST` 仍为 744/744，三作用域 `TOTAL_FAIL=0`，`go test ./...` 全绿。

**当前状态 = 待 ⑤ / ⑥ / ⑦ 并停于 ⑩ 授权提交。**

## 2. 做法（① 判定 → ② 落地）

### 2.1 影响判定（① `DeepSeek V4 Pro`）

- **不触** CONTRACTS / schema / 契约夹具 / `workflows` / CI workflow：两条旧路径在 `docs/CONTRACTS{,_EN}.md`、`schemas/`、`.github/workflows/` 零命中；变更面仅 `.md` 与一处 `.go` **注释行**。
- **不需准则 MINOR**：两条旧路径在 `docs/DELIVERY_DIRECTIVE{,_EN}.md` 零命中 ⇒ 正文一字不改，§0.3 变更日志义务不触发。
- **不触**判据本体与判据实现：本片只改证据包**数据**，`tools/gates/**` 零改动。
- **⑥ 不主张跳过**（步骤表含 ⑥，且变更集含 `.go`）。

### 2.2 归档布局

| 源 | 归档落点 | 件数 |
|---|---|---|
| `tmp` 根下 `P84-CI-FIX` | `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/` | 44 |
| `tmp` 根下 `P85-SECURITY-PROOF-SUBSYSTEM` | `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/` | 66 |

`archive/` 前缀沿用 `directive-history/…/archive/` 既有实作；含空格的子目录名（`Go ubuntu-latest` / `Go windows-latest`）与二进制 `attempt1-full-log.zip` **逐字保留**；`p85-worktree.patch` 与 `05-input/diff.patch` 与已入库 `P87-four-read-data/diff.patch` 同字节，仍**两份各自独立归档**（不去重；归档副本按溯源表做路径与编码归一）。

### 2.3 字节归一（必须做，理由是门禁事实）

`.gitattributes` = `* text=auto eol=lf`，CRLF 文件入库后 blob 内为 LF。若按工作区字节计哈希，`ci` / `range` 作用域（读 blob）必然报 G8-b 哈希不匹配。故**先归一再计哈希**：

- 88 件字节发生变化（多数为 CRLF → LF；`.md` / `.ps1` 按仓库惯例补 BOM，其余去 BOM）；
- `.zip` 二进制逐字复制；
- 110 件的「原始字节 sha256 / 归档字节 sha256 / 处置」逐件登记于 `archive-byte-provenance.tsv`，并由 ③ 从**未改动原件备份**重新推导复原（110/110 通过）。

### 2.4 引用处置

- **活跃面**（台账 P85 段落点声明、台账 P89 目标字段、P84 报告对 §12 两处、代码注释）→ **改写**为归档落点，CN / EN 同位置等量替换（G3(a.1) 对称）。
- **归档副本内自引用**（11 件 39 处）→ **路径字面归一**为归档路径（保留 `#L…` 锚点），每件处数经断言校验：

| 归档副本（省略公共前缀） | 处数 |
|---|---|
| `P84-CI-FIX/07-narrow-review-request.md` | 3 |
| `P84-CI-FIX/09-narrow-review-request-r2.md` | 3 |
| `P84-CI-FIX/10-final-review-request.md` | 4 |
| `P85…/01-design-rev2.md` | 1 |
| `P85…/01-design-review.md` | 1 |
| `P85…/p85-worktree.patch` | 2 |
| `P85…/04-reviews/05-pre-review-r1.md` | 1 |
| `P85…/04-reviews/05-pre-review-r2.md` | 1 |
| `P85…/04-reviews/07-carbon-r7-targeted-read.md` | 1 |
| `P85…/04-reviews/07-codex-baseline-read.md` | 20 |
| `P85…/05-input/diff.patch` | 2 |

  归一目标 = 各自包的 `archive/<原目录名>/` 前缀；逐件双 sha256 见 `archive-byte-provenance.tsv`。
- **冻结面** → 保留原文：P85 起跑消息 `00-briefing.md` 仅在**文末追加**归档注记（`git diff` 为纯新增、零删除），四份审查原文与 `diff.patch` 一字不改（补丁载荷追加会破坏其语法），v72 快照按起跑消息「不动」处置、留待 v73 吸收。

### 2.5 相对 ① 方案的实施口径收紧（如实披露）

① D5 / D6 建议把「MANIFEST 归档来源登记行」与「本报告改写对照表」也纳入命中白名单。实施时改为**非连续写法**（`tmp` 根下 `P84-CI-FIX`），使这两类载体不再产生可解析命中 ⇒ 白名单从四类收缩为两类（冻结原文 + 冻结快照），核验更强而非更弱。另：① D8 步骤 4 写作「非 `.ps1` 去 BOM」，与仓库 `.md` = BOM 惯例及证据树既有 735 件 BOM 载荷不符，实施按**扩展名惯例**（`.md` / `.ps1` 带 BOM，其余不带）执行。

## 3. 机械验证读数

### 3.1 归档包三方自洽（`pack-consistency.txt`）

| 包 | 载荷 | MANIFEST 条目 | SHA256SUMS 条目 | 三方相等 | 哈希漂移 |
|---|---|---|---|---|---|
| P84-BOUNDDIR-RESIDUAL | 58 | 58 | 58 | 是 | 0 |
| P85-SECURITY-PROOF-SUBSYSTEM | 95 | 95 | 95 | 是 | 0 |
| P89-TMP-ARCHIVE-P84-P85 | 11 | 11 | 11 | 是 | 0 |

> 时点口径：`pack-consistency.txt` 为脚本生成的快照，本表与其同源（终态本包 11 件载荷）；定稿态该包载荷与 `SHA256SUMS.txt` 逐字节一致、漂移 0。

### 3.2 溯源可复原

`archive-byte-provenance.tsv` 的 110 条记录逐条重算：原始字节 sha256 与记录一致、按同一变换重放后与归档字节 sha256 一致、且与磁盘逐字节相等（110/110，漂移 0）。

### 3.3 旧路径零命中的操作化判据

`git grep -n -F <needle> -- . ':!tmp'`（对起跑消息中出现的两条旧路径字面各跑一次，needle 逐字取自 `00-briefing.md`）的实际命中行集合，排序后与下列白名单**逐行相等**（无多余、无缺失）：

| 白名单行 | 类别 |
|---|---|
| `docs/t027/CHECKLIST_v72.md:137`、`CHECKLIST_v72_EN.md:137` | 冻结快照（v73 吸收） |
| `docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/00-briefing.md:46` | 本片起跑消息逐字原文（§13.1.10） |
| `…/01-arch/01-arch-design.md:32`、`:61` | 本片 ① 方案书原文中的命令与裁定引用（§13.1.13） |
| `…/04-reviews/06-scan-r1.md:25`、`:31`、`:62` | 本片 ⑥ 扫描原文中的命令与读数引用（§13.1.13） |
| `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/00-briefing.md:109` | 冻结原文（起跑消息逐字） |
| `…/07-carbon-r7-targeted-read.md:10` | 冻结原文（审查原文） |
| `…/P87-four-read-data/07-codex-baseline-read.md:11/14/21/26` | 冻结原文（盲读原文） |
| `…/P87-four-read-data/diff.patch:237/11818` | 冻结载荷（补丁语法不可注记） |

命中行数 = 白名单行数（16）；期望外命中 0；白名单空转 0。归档副本 110 件内旧路径残留 0。

### 3.4 目录清除与 tmp 顶层

`Test-Path` 读数：删除前 P84 = True / P85 = True；删除后 P84 = False / P85 = False。`tmp` 顶层当前 = `.gitkeep` + `gates-cache`（仅 `g4b`）+ `find-quotes.js` + 本片临时目录 `P89-TMP-ARCHIVE-P84-P85/`（⑧c 清理后即为 ⑩ 声明态）。

### 3.5 回归

- `gate.js --selftest` = `SELFTEST: PASS 744/744`。
- `--all --scope=tree` / `--scope=index` / `--scope=ci` 均为 `TOTAL_FAIL=0`（⑥ 转录插入后复跑；本机 ci 作用域为辅助读数，权威读数以 CI 双腿为准）。
- `gofmt -l .` 空、`go build ./...` / `go vet ./...` / `go test ./...` 退出 0，15 包 ok、0 FAIL。

### 3.6 门禁读数

② 完成后首轮（⑥ 转录前，改动集 129 路径）：`--all --scope=tree` 与 `--scope=index` 均为 `TOTAL_FAIL=1`，**唯一红项 = `G7-f`**（新增报告须带 `sixstage` 围栏，而该围栏即 ⑥ 转录产物）；其余判据（含 G8 三包自洽、G9、G4a、G5(d)、G6-4、G2）全绿。⑥ 后插入 `sixstage` 并复跑至 `TOTAL_FAIL=0`；⑨ 定稿时回填三作用域终态读数与 CI 双腿。

## 4. 审查链

- ⑤ 预审（`DeepSeek V4 Pro`，额度 1+1）：见 §4.1。
- ⑥ 独立扫描（`MAI Code 2 Flash`，额度 ≤3）：见 §4.2。
- ⑦ 终审（`GPT-5.3-Codex` 常规档，额度 1+1）：见 §4.3。

### 4.1 ⑤ 预审（`DeepSeek V4 Pro`，额度 1+1，用 2）

输入 = 本报告 + 证据包 + 变更集。r1 = `PRE-REVIEW: FINDINGS`（1 HIGH + 4 MEDIUM + 4 LOW，核心为「零命中白名单未覆盖本片自撰件」）；按最小修复逐条处置后 r2 = `PRE-REVIEW: FINDINGS`（H1 与 M 系全部成立，新增 2 MEDIUM + 1 LOW，均为描述 / 计数自相矛盾）——三项已由脚本读数与 `git grep` 复算闭环；额度用尽未开第 3 轮。原文：`04-reviews/05-pre-review-r1.md`、`04-reviews/05-pre-review-r2.md`。

### 4.2 ⑥ 独立扫描（`MAI Code 2 Flash`，额度 ≤ 3，用 1）

输入 = 本报告 + 变更集 + 归档自洽证据。`INDEPENDENT SCAN: PASS`，9 项核验全过、0 条可执行发现；六阶段转录见本文件末 `sixstage` 围栏。原文：`04-reviews/06-scan-r1.md`。

### 4.3 ⑦ 终审（`GPT-5.3-Codex` 常规档，额度 1+1，用 2）

输入 = 本报告 + 背书点 B1–B4 + ⑤⑥ 结论。r1 = `FINAL REVIEW: FINDINGS`（1 HIGH：白名单外仍有 5 处命中 ① 方案书与 ⑥ 扫描原文；1 MEDIUM：§3.1 载荷读数滞后；1 LOW：源分布 44/66）→ 采纳「白名单扩容并解释」方案（扩至十六行、新增「本片 ① / ⑥ 原文命令引用」类）→ r2 = `FINAL REVIEW: PASS WITH FIXES`（唯一残留 = `pack-consistency.txt` 快照滞后，已重生成至漂移 0），并明确「该修复完成后即可按 PASS 提交」。背书裁定：B1 成立（十六行白名单化口径）、B2 成立、B3 成立（G8-c 未启用属已披露的宽松，G8-b 三方相等仍成立）、B4 成立、F-09 判**未越界**、F-13 成立（偏离可接受）。

## 5. 计量

- ① × 1（`DeepSeek V4 Pro`）；② × 1（`prfrail-implementer` + `DeepSeek V4.1 Flash`）；③ × 1（`prfrail-tester` + `DeepSeek V4.1 Flash`）；⑤ × 2；⑥ × 1；⑦ × 2（终审 1 + 复审 1）；批量辅助 GLM × 1（引用清单与双语镜像核对，载体 `deep-reasoner`）。
- 探针 0；盲审不抽（起跑消息用户裁定 c）。
- 计量单一真值源 = 会话调用回执。

## 6. 异常处理

- **双 BOM 事故（已回正）**：对 P85 `00-briefing.md` 追加注记时写入了第二个 BOM，使首行被记为「删除 + 新增」。已剥离多余 BOM 并复验：`git diff --numstat` = 纯新增、删除 0，逐字原文完好。
- **`grep` 实测与台账差额**：台账记 9 件载体 15 处，实测 11 件载体 39 处，差额为 `CHECKLIST_v72{,_EN}.md:137` 两行 × 两条路径（冻结快照，本片不动）；已在 `input-list.md` §三 逐项列明，本片不改台账该行（历史 不追溯）。
- **① D8 编码口径更正**：见 §2.5。
- 临时件 D 计数：0（`tmp/P89-TMP-ARCHIVE-P84-P85/` 内脚本、备份与捕获件于 ⑧c 清理）。

## 7. 移交

- `CHECKLIST_v72{,_EN}.md` 的两行冻结命中由 v73 快照换代吸收。
- `find-quotes.js` 仍待 tmp 生命周期片处置（本片边界）。
- 层 3 四个性能候选登记不变。

```sixstage
## Section A 扫描范围与方法

九项核验：仓库状态与变更集总量、三包载荷三方自洽（独立枚举 + 逐条 sha256 复算）、门禁双作用域、旧路径残留、归档保真抽查（5 件，含 1 件二进制与 1 件被路径归一者）、冻结面零删改、报告对镜像与围栏集合、回归（selftest / build / vet / test）、台账状态。

## Section B 发现

| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| - | - | 无可执行发现 | - |

## Section C 误报排除

- `G7-f` 为红系 ⑥ 转录尚未落盘之结构必然，其余判据实测全绿（本轮后已补围栏）。
- P85 `diff.patch` 内的通用 `tmp/` 文字不是旧临时路径，精确前缀命中为 0。
- 台账「现状缺口」字段是起跑基线描述，不是状态漂移。

## Section D 覆盖面说明

归档保真与字节溯源（110 件）、三包 MANIFEST / SHA256SUMS / 磁盘三方相等、冻结面只增不删、报告对逐行镜像、围栏条目集与变更集相等。

## Section E 结论

INDEPENDENT SCAN: PASS
```

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/bounddir-residual.md	repo
docs/validation/bounddir-residual_EN.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/MANIFEST.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/SHA256SUMS.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/01-r2-contract-judgment.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/02-r2-verification-matrix.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/03-r1-design.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/04-ci-rerun-evidence.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/07-narrow-review-request.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/08-narrow-re-review.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/09-narrow-review-request-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/10-final-review-request.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/11-narrow-re-review-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/12-rectification-record.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/13-final-narrow-review.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/14-rectification-A.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/15-final-review-A.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-full-log.zip	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-shortname-anchors.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-test-fail-list-normalized.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-unzip/0_Go windows-latest.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-unzip/1_Go ubuntu-latest.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-unzip/Go ubuntu-latest/system.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-unzip/Go windows-latest/system.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt1-windows-failed-log-excerpt.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/attempt2-windows-failed-log-excerpt.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-go122-r4.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-go127-r4.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-suite-go122-r2.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-suite-go122-r3.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-suite-go122.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-suite-go127-r2.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-suite-go127-r3.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/full-suite-go127.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/gates-index-r4.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/gates-index-r6.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/gates-tree-r4.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/gates-tree-r6.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1-probe-output-go127.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1-probe-output.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r2.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r3.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r4.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r5.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff-r6.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/r1r2-final-diff.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/w12-go122-r4.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/archive/P84-CI-FIX/w12-go127-r4.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/00-briefing.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/MANIFEST.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/SHA256SUMS.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design-rev2.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design-rev3.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design-rev4.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design-review.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-r7-start-message.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/02a-anchor-check.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/02a-cjk-check.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/02a-contracts.diff	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/02a-purity.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/02a-verify.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/02c-f04-ruling.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/03-tester-t2.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/03-tester-t3.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/03-tester-t8a.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/06-scan-r4.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-carbon-main-read.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-carbon-r7-targeted-read.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-codex-baseline-read.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-final-review.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-neon-cross-read.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-re-review-r1.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-re-review-r2.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-re-review-r3.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-sol-baseline-read.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/07-three-read-comparison.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/08-bc-handoff.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/04-reviews/09-l4-diff-excerpt.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/ci-45208001981-evidence.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/diff.patch	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-ci.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-index.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/gate-tree.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/local-full-log.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/local-r7-legs.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/mutation-r4.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/mutation-r6.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/mutation-r7-executed.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/numstat.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/old-wording-grep.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/probe-verify-main.go.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/r7-fix.diff	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/selftest.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r4.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/test-r6.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/vm-full-log.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/vm-go-test.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/vm-leg-readings.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/vm-r4-legs.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/vm-r6-legs.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/05-input/vm-r7-legs.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/input-list.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/land-r4-contracts.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/land-r4.out.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/p85-worktree.patch	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/report.md	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/space-r4.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/space-r4.out.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/strict-r4-contracts.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/strict-r4.out.txt	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/verify-r4-contracts.js	repo
docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/verify-r4.out.txt	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/00-briefing.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/01-arch/01-arch-design.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/07-re-review-r1.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/MANIFEST.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/SHA256SUMS.txt	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/archive-byte-provenance.tsv	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/freedom-list.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/input-list.md	repo
docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/pack-consistency.txt	repo
docs/validation/tmp-archive-p84-p85.md	repo
docs/validation/tmp-archive-p84-p85_EN.md	repo
internal/console/p85_t2_frozen_policy_tester_test.go	repo
```

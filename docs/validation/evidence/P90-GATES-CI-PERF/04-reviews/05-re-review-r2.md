# ⑤ 预审 · 第 2 轮（复审，DeepSeek V4 Pro，只读载体 `deep-reasoner`）

- 日期：2026-11-25
- 输入：整改后的台账 P90 段（CN/EN）、`01-design.md`（含 `## Corrections after ⑤ round 1`）、`ci.yml`、重生成的 `impl.diff` / `numstat.txt`
- 结论：`RE-REVIEW: PASS WITH FIXES`（本轮 3 项整改全部 CLOSED；残余 1 Medium 为「④ 读数尚为整改前状态」的时序项）

---

## F1 台账 CN/EN 边界单元格 — **CLOSED**

- CN `docs/t027/REMAINING_SLICES.md:2436`：边界单元格现含「`workflow_mutation_test.go` 的 `different-shell` 变异值必要适配——`bash` 改为 `sh`，否则对 `Format` 步为空操作而假红」。
- EN `docs/t027/REMAINING_SLICES_EN.md:2423`：对应句「one necessary adjustment of the `different-shell` mutation value in `workflow_mutation_test.go` (`bash` to `sh`, otherwise it is a no-op for the `Format` step and trips a false red)」——两侧语义逐项一致。
- 与 numstat 不再矛盾：实时 `git diff --numstat` 中 `internal/release/workflow_mutation_test.go` = **1/1**，且磁盘真值 `workflow_mutation_test.go:58` 为 `func(step *workflowStep) { step.Shell = "sh" }`（`workflow_test.go` 内无 `different-shell` 命中，与「仅 mutation 文件 1 行」吻合）。
- 两侧均保留「主控代裁，待用户追认，不得作先例」/「acting on the user behalf, pending ratification, not a precedent」。
- 实时 numstat：CN **22/0**、EN **22/0**，仍对称（见「另复核」）。

## F2 01-design.md Corrections 表 — **CLOSED**

- ① 原文未改写：`01-design.md` 的 §A–§J + `DESIGN: DONE` 完整，`## Corrections after ⑤ round 1` 追加在末尾（129 行起），且首行声明未改写。佐证：原文中的失效预测值仍在——§C `ff770e4b…`、§J「`workflow_mutation_test.go` 零改动」、§F「go.mod-keyed」候选文案均未回改，与 Corrections 表形成「预测→实际」对照而非历史改写。
- 表内三条与磁盘一致：
  - **C1**：实际 digest `5a9a2b5725f2311225eaa69c3be661ca5caa8602563a1158d693057af3777471` = `workflow_test.go:250` 的 `test/Format` pin 值。
  - **C2**：`workflow_mutation_test.go:58`（`different-shell` → `sh`）+ numstat 1/1，与表内证据一致。
  - **C3**：`ci.yml:71-72` 注记已无 `go.mod-keyed`。
- A/B 披露：`ab-report.json` 的 `aHeader`/`bHeader` 均为 `changed=6`；`git rev-parse HEAD` = `6e097fd`，`git diff --name-only HEAD^1 HEAD` = **6 文件**，与披露的「commit `6e097fd` (six files)」一致。

## F3 ci.yml Setup Go 注记 — **CLOSED**

- `ci.yml:71-73`：注记为两行「setup-go restores the Go module cache and the Go build cache, which warms Build/Vet/Test; keep cache enabled (P90 assessment: already covered).」，不含 `go.mod-keyed`。
- YAML 数据未变：`workflow_test.go:326` 仍为 `"test/Setup Go": {"go-version": "${{ env.GO_VERSION }}", "cache": true}`；diff hunk `@@ -68,6 +68,8 @@` 只在 `cache: true` 上方加注释行，未动键值。
- `go test -count=1 ./internal/release/` 实测：`ok github.com/larsonzh/prfrail/internal/release 7.538s`，绿。

## F4 impl.diff / numstat — **CLOSED**（附证据完备性注记）

- `tmp/P90-GATES-CI-PERF/review/impl.diff`：15295 字节，`TextDecoder('utf-8', {fatal:true})` 解码通过 ⇒ **合法 UTF-8**（fatal 解码为权威判据；`invalidLeadBytes=2059` 的计数来自把合法 UTF-8 续字节误计入的辅助计数器，不可采信）。
- `numstat.txt`（5 行）与实时 `git diff --numstat` **逐行一致**：`ci.yml 15/1`、CN `22/0`、EN `22/0`、mutation `1/1`、workflow_test `5/2`。
- ci.yml 行数自洽：终态 hunk `@@ -68,6 +68,8 @@`（注记 +2 行）与 `@@ -90,8 +92,20 @@`（Format 步 + Gates 注释 + Gates run 行替换）合计恰为 15/1。注记从 §F 候选 1 行文案到落盘 2 行的结构变化与 hunk 相符。**注记**：review 目录仅含 `impl.diff` + `numstat.txt`，整改前中间态 numstat 无磁盘快照，「因注记改写而变化」无法直接比对，只能以终态 hunk 自洽推断。

## Residual

1. **Medium（既有，非本轮整改引入）**：④ 整改前读数 `tmp/P90-GATES-CI-PERF/gate-tree-4.txt` 为 `TOTAL_FAIL=1`——`FAIL G8-a evidence-tree path shape`（`00-briefing.md` / `01-design.md` 未注册为 pack member / 未登记 `tools/gates/evidence-forms.txt`）。本轮三项整改未触及注册面 ⇒ 预期 ④ 重跑时仍红，须在下轮 ④ 前处置（登记或注册）。
2. **Low（流程，非内容缺陷）**：整改后 ④ 尚未重跑。`gate-tree-4.txt` 的 `G3(a) PASS`（`+22/-0 vs +22/-0 OK`）为**整改前读数**；实时 numstat 两侧仍 22/0 对称 ⇒ 机械预测 G3(a) 仍 PASS，但按纪律**不作「已绿」断言——需重跑 ④ 取整改后读数**，并顺带确认第 1 条 G8-a 状态。
3. **Low（证据完备性）**：整改前中间态 numstat 无快照（见 F4 注记）。
4. **代裁表述核对（无问题）**：`00-briefing.md:193-194`（「系主控代裁…须用户事后追认；**不得作先例**」）、`01-design.md:10/88/125`、台账 CN `:2436`/EN `:2423` 均保留「主控代裁 + 待用户追认 + 不得作先例」，未弱化为「已获授权」；`01-design.md:88` 的「授权本身已由用户裁定 f 给出」系 §G 原文（指 f 裁定的 CI 改动授权面，非代裁面），ci.yml 新注记不含授权字眼，无弱化。

`RE-REVIEW: PASS WITH FIXES`

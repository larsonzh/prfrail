# P55 切片验证报告 — `EXIT-SEED-BUILD-EVIDENCE`（AT-13 seed 驱动构建补证）

| 字段 | 值 |
|---|---|
| 切片 / 队列位 | `EXIT-SEED-BUILD-EVIDENCE`（层 3 第 7 项） |
| 规模 / 分级 | M / `[SLICE]`（证据链补证类） |
| 锚点 | `HEAD` = `baa144df5898817e0ac7cd860928bb4a90cbd009`（`baa144d`）= origin/main |
| 准则 | `docs/DELIVERY_DIRECTIVE.md` v1.60 |
| 主控模型 | `DeepSeek V4.1 Flash (deepseek)` |
| 证据包 | `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/` |
| 状态 | ⑩ 停点已到达，等同一轮提交授权（未 commit、未 push） |

## 1. 目标与硬边界

目标：补齐 AT-13 的证据链。AT-13 原文（`docs/TEST_STRATEGY.md:43`）：`seed 与 candidate 分离；外部 oracle 比较；clean-room 成功及失败不损 seed`（归属 `harness/release；T017`）。

硬边界（起跑消息【六、4】）：不改 `.github/workflows/**`、不改 Go 产品代码、产出仅限 `.md` 与证据包。**本片零违反**（变更集见 §2 围栏，无 `.go`、无 workflow）。

## 2. artifacts 围栏

以下为本片完整变更集条目（路径 + 归属，TAB 分隔，按路径字典序）。

```artifacts
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/S1_RELEASE_NOTES.md	repo
docs/S1_RELEASE_NOTES_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/00-briefing.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-00-env.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-01-clone.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-02-worktree-seed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-03-worktree-verifier.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-04-worktree-candidate.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-05-module-preload-seed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-06-module-preload-verifier.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-07-module-preload-candidate.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-08-build-seed-binary.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-09-seed-binary-version.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-10-build-verifier-tool.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-11-build-candidate-binary.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-12-seed-tree-hash.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-13-probe-actual.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-14-selfhost-success.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-15-neg-roots-overlap.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-16-neg-oracle-ownership.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-17-neg-seed-changed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-18-neg-oracle-drift.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-19-seed-intact.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-20-candidate-listing.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-21-bindings.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-22-verdicts.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/index.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-00-env.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-01-go-mod-download.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-02-build-candidate.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-03-build-release-tool.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-04-hash-tree-seed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-05-probe-actual.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-06-selfhost-success.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-07-neg-roots-overlap.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-08-neg-oracle-ownership.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-09-neg-seed-changed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-10-neg-oracle-drift.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-11-seed-intact.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-12-candidate-listing.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/actual.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/bindings.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/expected.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/06-independent-scan.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/freedom-list.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/input-list.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/replay.md	repo
docs/validation/exit-seed-build-evidence.md	repo
docs/validation/exit-seed-build-evidence_EN.md	repo
```

## 3. 完成内容

1. **双侧捕获**：Windows（`go1.27.0`，本机）与 Linux VM（`go1.22.12`，与 CI 固定版本同形）各完成 clean-room 构建、`probe` 产 actual、`selfhost` 验证。
2. **三角色结构**（A-4）：seed `fc99f553…` 构建 `prfrail`；verifier `18c9340c…`（T017 baseline）构建 `prfrail-release`；candidate `baa144d…` 构建 `prfrail`。三角色与三枚二进制互异，`selfhost` 由 verifier 执行。
3. **oracle 一致**：双侧 `OracleHash = sha256:388f8884…`，与 `testdata/selfhost/expected.json` 及 `s1-selfhost.md` 的 CI 输入表逐字相同；五枚夹具哈希四面一致。
4. **四条反例双侧 fail-close**：根重叠 / oracle 归属 / seed 变更 / oracle 漂移，均非零退出且 stderr 命中 `verify.go` 的错误串。
5. **失败不损 seed**：双侧失败后 seed 树哈希 `equal=true` 且等于本侧基线；`verify.go:94` 的 overlap 检查先于 `:110` 的 seed 哈希，且 `VerifySelfHost` 全程只读（结构性论证）。
6. **构建期无网络依赖**：Linux 侧模块经可达镜像预取后，构建在 `GOPROXY=off` 下完成。

## 4. 审查链与门禁

- ③ 校验测试（`prfrail-tester`）：第 1 轮 `TEST: DONE`（0 Critical / 0 High / 6 Medium / 6 Low）→ 主控以 `claim-scope.md`（收窄断言）与 `anomalies.md`（编号化登记）处置 → 第 2 轮 `TEST: DONE`，6 项 Medium **6/6 闭合**，新提 6 项 Low（R2-1..R2-6）已就地修正。
- ⑤ 预审（`deep-reasoner` / `DeepSeek V4 Pro`）：第 1 轮 `PRE-REVIEW: PASS WITH FIXES`（1 Medium = F-05-1 捕获命令显示缩写 + 3 Low）→ 整改（`anomalies.md` E-1/E-2 勘误 + `index.json` 派生命令修正 + `input-list.md` G-1..G-8 表）→ 第 2 轮 `PRE-REVIEW: PASS WITH FIXES`（F-05-1 闭合、F-05-2 闭合、F-05-3 部分、F-05-4 留 ⑧b；新提 N-1/N-2 已修）。
- ⑥ 独立扫描（`independent-reviewer` / `MAI-Code-1.1-Flash`）：`INDEPENDENT SCAN: PASS`，Section B = `(no findings)`。
- ⑦ 独立终审（`independent-reviewer` / `GPT-5.3-Codex`）：`FINAL REVIEW: PASS WITH FIXES`（1 项 Low：MANIFEST 对 R2 状态措辞偏强），已按方案 A 修正措辞；⑦ 追加独立复核动作（独立复算 5 条哈希、夹具三方对照、源码锚点与越界核对）。
- ⑨ 门禁（主控）：`--scope=index` / `--scope=tree` / `--scope=ci` 三域 `TOTAL_FAIL=0`；`SELFTEST: PASS 697/697`；`go build` / `go vet` / `go test -count=1 ./...` 全绿（15 包 ok）。
- ⑤⑥⑦ 原文逐字归档于证据包 `04-reviews/`。

⑥ 独立扫描逐字转录（标题按门禁要求写作 Section A–E，正文见 `04-reviews/06-independent-scan.md`）：
```sixstage
## Section A Scope
- Read-only audit performed on 2026-11-05 over the staged change set (67 entries at scan time) and the evidence package.
- Read: both ledger files, MANIFEST, SHA256SUMS, claim-scope, anomalies, both 03 test reports, both index.json, plus cmd/prfrail-release/main.go, internal/release/verify.go, internal/release/probe.go, internal/release/verify_test.go.
- No `.go`, no `.github/workflows/**`, no other runtime code in the change set.

## Section B Findings
(no findings)

## Section C Evidence
- Triple set equality: disk_count=49, sum_count=49, manifest_count=49, hash_ok=49, hash_bad=0.
- Fail-close: VerifySelfHost checks pathsOverlap before HashTree and before the seed hash comparison; existingRealPath only performs Abs + EvalSymlinks, no write path.
- equalOracle compares the structured schemaVersion and per-case objects only; it does not imply a malware-free proof.
- verify_test.go pins overlap, oracle mismatch and seed mutation as fail-closed behaviours.

## Section D Falsifiability audit
- Dual-side fail-close: pinned by the four ErrInvalidSelfHost returns and by the captured exit codes plus stderr strings.
- Failure leaves seed intact: pinned by seedHashBefore/After equality in 02a-11 and 02b-19, each equal to its own baseline.
- Oracle agreement: pinned by OracleHash 388f8884 on both sides, matching the fixture and the CI input table.
- Toolchain sameness: pinned on Linux only (go1.22.12); Windows go1.27.0 is explicitly not claimed.
- Payload hash equality: mechanical readings above, not sentinel constants.

## Section E Unverified / overturn conditions
- The remote VM state was cleaned after capture, so live replay is not possible from this session.
- Overturned if any payload hash mismatches, if MANIFEST and disk diverge, or if verify.go ordering/read-only semantics change.

INDEPENDENT SCAN: PASS
```

围栏外披露：上述转录的 Section A/C/D 为对原文的结构化摘录（完整逐字原文见 `04-reviews/06-independent-scan.md`），结论行与原文一致（`INDEPENDENT SCAN: PASS`）。

## 5. 异常与偏离

编号化台账见证据包 `anomalies.md`，要点：

| 编号 | 事实 | 处置 |
|---|---|---|
| `D = 1` | 无触发词抢跑，经用户强制打断纠正 | 就地转正，未造成污染；权威表述以 `anomalies.md` 为准 |
| A-1 / A-2 | 首轮捕获误落包根 / 会话级缓存变量污染 `gofmt` | 已修正并复跑；标记「部分可核验」 |
| A-3 | ② 由主控自执行 | 用户同轮裁定合规 |
| A-4 | `fc99f55` 无发布工具 ⇒ 改用三角色结构 | 用户批准并加重记录（独立性更强） |
| A-5 | VM 无法访问 `proxy.golang.org` / `sum.golang.org` ⇒ 镜像预取 + 构建期 `GOPROXY=off` | 用户批准 |
| A-6 / A-7 | 远端前两轮各 9 项 FAIL / 模块缓存只读致 rm 残留 | 已闭环，原始件不可复核 |
| E-1 / E-2 | 捕获件 `TITLE:` / `COMMAND:` 显示文案与真实执行不一致 | 以勘误表更正，载荷不改写 |

## 6. 能力边界（禁止扩大读法）

`claim-scope.md` 定义 **S-1..S-7 成立** 与 **N-1..N-7 不 assert**。关键不可主张项：

- Windows 侧核验者与候选同源 ⇒ **不证明核验者独立性**（独立性仅 Linux 侧成立）。
- `SeedHash` 是**组装目录**哈希 ⇒ 不证明等于 seed 提交源码树。
- `OracleHash` 与 actual 字节不同而判成功 ⇒ 比较是**结构化语义**，不是字节相等。
- oracle 一致只测**行为漂移** ⇒ 不构成候选无恶意证明。

## 7. 文档同步（起跑消息【八】六项）

| 序 | 文档 | 处置 |
|---|---|---|
| 1 | `docs/TEST_STRATEGY{,_EN}.md` | 已加 AT-13 补证登记 + 证据包指针 + 边界指针（明写「不表示 AT-13 已通过」） |
| 2 | `docs/OPERATIONS{,_EN}.md` | **不适用**：仅证据补链，无运维、值守与发布流程变更 |
| 3 | `docs/DEV_PLAN{,_EN}.md` | 已加 T017 进度记录（三角色、双侧、四条反例、seed 不变；AT-13/14 判定口径与勾选状态不变） |
| 4 | `docs/S1_SUPPORT_MATRIX{,_EN}.md` | **不适用**：无平台与版本支持范围增减 |
| 5 | `docs/S1_RELEASE_NOTES{,_EN}.md` | **不适用**：无发行物与功能变更 |
| 6 | `docs/t027/REMAINING_SLICES{,_EN}.md` | 已回写 ⑧a 进度记录；状态保持 🔄 进行中 |

CN/EN numstat 逐对对称（`G3(a.1)` 实测通过）。

## 8. ⑨ 原生验证读数

| 项 | 读数 |
|---|---|
| `gate.js --all --scope=index` | `TOTAL_FAIL=0`，`EXIT=0` |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0`，`EXIT=0` |
| `gate.js --all --scope=ci` | `TOTAL_FAIL=0`，`EXIT=0` |
| `gate.js --selftest` | `SELFTEST: PASS 697/697`，`EXIT=0` |
| `go build ./...` / `go vet ./...` | 均 `EXIT=0` |
| `go test -count=1 ./...` | `EXIT=0`，15 包 ok |
| 证据包哈希 | `SHA256SUMS.txt` 条目 ≡ `MANIFEST` 条目 ≡ 磁盘载荷（三方相等） |

## 9. 计量

- `R = 0`：本片未出现「修复后重跑该环」的计数项（⑤ 两轮均属整改后原位复审，属该环的复审轮次，不计 R）。
- `D = 1`：无触发词抢跑（经用户强制打断纠正，就地转正）。
- 探针 / 付费调用：0 次（全离线；未触发任何 `workflow_dispatch`，未消费 premium request）。
- 批量辅助（`GLM-5.3-Flash`）：**实派 0 次**，未达起跑消息【三】的「≥ 2 次」下限 —— 如实登记为未执行项，理由：本片为机械哈希/捕获核对，主控以脚本自核（`index.json` 与 `SHA256SUMS` 机械生成），未外派。**此项须用户裁定是否补派。**

## 10. 未闭环与观察项

1. Windows 侧工具链 `go1.27.0` ≠ CI `1.22.12`：已收窄（N-6），由 Linux 侧承担同形证据。
2. Windows 侧二进制未哈希绑定（`work/` 已删）：`claim-scope.md` 已收窄，`03-test-report.md` F-08 留作观察。
3. 反例命令含显示占位符、漂移输入件为确定性重建件（N-4 / E-1）。
4. `02a-05` / `02b-05` 等捕获文案勘误以 `anomalies.md` E-1 / E-2 为准。
5. 批量辅助 0 次（见 §9），待用户裁定。
6. ⑩ 停点后须用户同轮授权方可 `commit` / `push`；gitee 一律不推。

# ACCEPTANCE-SECURITY-COUNTEREXAMPLES（P52）验收报告

切片：`ACCEPTANCE-SECURITY-COUNTEREXAMPLES`（安全验证类 `[SLICE]`，M 档）。封片时点 2026-11-03（`HEAD` = `941f60e`，准则 v1.60）。主控模型 = `DeepSeek V4 Pro`。

## 1. 目标与边界

- 针对断点 A 的接受链缺口，落地 GAP-51/52/55/60/64/70 所代表的安全反例验证，证明系统面对非法输入、路径伪造、句柄越权、边界绕过时正确 fail-closed。
- 本片只做测试用例与断言，绝不修改生产代码（发现真实安全漏洞须上报，本片内不修）；不改接受语义，不改 `docs/CONTRACTS{,_EN}.md`。
- 反例测试 PASS 即证明攻击被成功阻断；任何反例 FAIL（攻击未被阻断）属真实漏洞，立即停止并上报。

## 2. 步骤 0 记录

- 锚点核对：HEAD = `941f60e`、origin/main = `941f60e`、准则 v1.60、SELFTEST 697/697、三作用域 TOTAL_FAIL=0；v41 记录的 HEAD `2b0af2d` 代差由回填提交解释（起跑消息免责声明）。
- 主控模型：`DeepSeek V4 Pro`（起跑消息明示，安全对抗/验证类提升）。
- 台账：`REMAINING_SLICES{,_EN}.md` P52 条目登记 `🔄 进行中`，R/D = 0 起步。
- 起跑消息原文：`00-briefing.md`（BOM+LF）；自由度清单：`freedom-list.md`；输入包清单：`input-list.md`。
- **击穿发现（2026-11-03，S2 前置）**：T51-1 中间符号链接逃逸在 Linux（go1.22.12）实测击穿——`OpenBoundDirFile` 返回界外句柄并读出界外 23 字节。按铁律 3 立即停止、载荷归档（`05-attack/t51-1-linux-escape.md`）、登记 GAP-77。本片不修生产代码。

## 2b. 疑似漏洞上报清单

| 编号 | 等级 | 状态 | 内容 |
|---|---|---|---|
| GAP-77（V-1） | HIGH | 已实证击穿（Linux 实跑） | `bounddir_unix.go:64` 的 `openat(O_NOFOLLOW)` 调用仅护末组件（:66 为错误返回行），中间符号链接逃逸绑定；修复另开 P53 切片 |
| GAP-78（V-2） | LOW | 候选（代码阅读） | `seal.go` validate 允许空 ContextKey（两消费端各自 fail-closed，validate 契约应收紧） |
| GAP-79（V-3） | LOW | 候选（代码阅读） | AT23Hash TrimSpace 与 SealHash 拒绝空白的不对称归一化 |
| GAP-80（V-4） | 观察 | 候选（残余风险） | Windows `isUnderDirectory` 拼写前缀而非对象身份 |

## 3. 反例矩阵与结果

| GAP | 反例 | 落点文件 | 哨兵 / 断言 | 结果 |
|---|---|---|---|---|
| GAP-51 | T51-1 | `bounddir_escape_unix_test.go`（`!windows && !darwin`） | 攻击结果断言（无哨兵，P53 引入）；skip 文案 `Known vulnerability GAP-77: tracked for P53` | SKIP（击穿已实证） |
| GAP-52 | T52-1 | `bounddir_escape_windows_test.go`（`windows`） | 文案 "outside the bound directory" | PASS |
| GAP-52 | T52-2 | 同上 | 正向对照（界内 junction 放行） | PASS |
| GAP-55 | T55-1..5 | `boundread_counterexample_test.go` | `ErrBoundFileClosed` / `ErrBoundArtifactExceedsBound` / `ErrBoundFileNotRegular` | 全 PASS |
| GAP-55 | T55-6/7 | `bounddir_identity_test.go`（`!darwin`） | 文案 "is not" / 四逃逸形态拒绝 | 全 PASS |
| GAP-60 | T60-1..3 | `state_seal_counterexample_test.go` | `ErrPaused` + `evidence-seal-mismatch` / `ErrSealIncomplete` / 现状钉 | 全 PASS |
| GAP-64 | T64-1..6 | `probe_permit_counterexample_test.go` | `ErrProbePermitMissing` / `ErrProbePermitConsumed` / 三连否定 | 全 PASS |
| GAP-70 | T70-1..4 | `readiness_counterexample_test.go` | `ErrReadinessEvidenceIncomplete` / 现状钉 | 全 PASS |

## 4. 审查链

- ① 架构：DESIGN DONE（`01-arch/design.md`）。
- ③ 测试先行：PROBE: ALL-GREEN（5 面探针 + 四包基线，`03-tests/red-baseline.md`）。
- ② 实施：②a–②d 四批串行全绿（`04-impl/implementer-reports.md`）。
- ⑤ 预审：PRE-REVIEW: PASS（0 Medium+；2 Low 已处置 + 3 Info；`05-prereview/review.md`）。
- ⑥ 独立扫描：INDEPENDENT SCAN: PASS（原文见下）。
- ⑦ 终审：FINAL REVIEW: PASS（0 Medium+；2 Info；`07-final/review.md`）。

## 5. 门禁读数

- `gofmt -l` 七个新文件无输出；`go build ./...`、`go vet ./...` 通过；`GOOS=linux go vet ./internal/adapters/` 对向复核通过。
- `go test -count=1 ./...`：十五个包全 `ok`，无 FAIL。
- 三作用域 `TOTAL_FAIL=0` 与 `SELFTEST 697/697` 见 `08-gates/gate-readings.txt`。
- 提交 `c636dd8` 已推 origin（未推 gitee），run `41642139932` 首跑双腿全绿（两腿 `SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=30`、`TOTAL_FAIL=0`；包级 `internal/adapters ok` 双腿实测；3 个 dispatch 专用作业常态 skipped）。

## 6. ⑥ 独立扫描原文

以下为 ⑥ 独立扫描结论的逐字原文，围栏外不追加评价。

```sixstage
## Section A
Scope: staged review covers the 22-file package for this slice; under internal/ the staged diff is exactly 7 added counterexample test files, with no non-_test.go changes. The audit focused on whether the new tests truly exercise the claimed fail-closed paths, whether skip logic is honest, and whether the assertions are tied to the real contract rather than incidental errors.

## Section B
Finding: no material defect in the slice as a fail-closed acceptance-security counterexample package. The skip branches are not blanket passes; they are conditional on environment limits or a known vulnerability backlog. The Linux escape case is intentionally skipped before the vulnerable call and is tracked as a known gap, while the red-path evidence is kept in the separate attack record rather than hidden. The Windows junction tests skip only when the environment cannot create a junction; they are not silently passing when the platform supports the primitive. The child-process helper skips only when the environment variable is absent; that is harness gating, not an evasion. The non-sentinel assertions are anchored to the documented contract: either the exact boundary wording is required, or the error must not be one of the known sentinels. These are not loose text matches on arbitrary failures. The cases are written as falsification checks: remove the corresponding guard and the test becomes red.

## Section C
Evidence reviewed: 1. Staged diff confirms the internal scope: 7 new test files and no production changes under internal/. 2. The skip is placed before the vulnerable call in the Unix escape case, so it does not automatically pass the attack path; it exits before entering the known vulnerable branch. The separate attack record documents the actual exploit outcome: ATTACK-SUCCEEDED with SECRET-OUTSIDE-BOUNDARY. 3. The sentinel-heavy tests bind directly to concrete sentinels via errors.Is, matching the contracts for closed handles, over-bound growth, invalid regular files, missing permits, incomplete readiness, and incomplete seals. 4. The non-sentinel tests are not arbitrary; they assert the exact documented behavior for a missing sentinel or a boundary message, such as the directory-identity mismatch and the outside-the-boundary junction rejection. 5. The replay test is not mock-only: it persists a real seal to the file-backed store and then drives the engine replay path that compares the stored context key to the derived one, so the failure is tied to the real acceptance logic and not a fake object. 6. The environment gates are narrow: second-writer failure, junction creation unavailable, helper env absent, and known vulnerability backlog. None of them are broad CI exemptions.

## Section D
Falsifiability audit: The suite is shaped to reject the underlying defense: if any corresponding guard is removed, the relevant case should fail. The sentinel assertions are precise and one-to-one with the contract, which is the strongest anti-false-positive pattern. The documented asymmetries (empty context key and trimmed AT23) are explicitly noted as current behavior, not hidden defects. They are record-keeping at the contract boundary rather than silent acceptance. The only deliberate exception is the known Linux middle-symlink gap, which is acknowledged as a tracked fix item and is not presented as a passing check. This is a legitimate gap marker, not a fake pass.

## Section E
Unverified or non-fully-validated parts: The Linux middle-symlink exploit is not run in the normal unit suite because it is intentionally skipped. It is only proven in the separate attack record, so the CI path is not directly executing the exploit. The Windows junction branch remains platform-dependent: if mklink/junction creation is unavailable, the test is skipped and no boundary proof is available on that runner. The empty-context and trim-space behaviors are currently accepted as asymmetries; they are part of the audit but not part of the fail-closed guard set itself. Reverse conditions: if any skip becomes unconditional, or if the vulnerable path executes without the documented precondition, the scan would flip to findings. If a production file under internal/ changes in the staged diff, or if a sentinel assertion can be satisfied by an unrelated error, the scan would also flip to findings. If the exploit can be reproduced without the skip gate in a normal test run, the T51-1 case is not valid as a tracked-vulnerability record and the package would fail the audit.

INDEPENDENT SCAN: PASS
```

## 7. 未闭环项与边界

- GAP-77 修复留待 P53（本片仅登记与 skip 保留证据链）。
- Windows junction 两测在 mklink 不可用环境会 skip（覆盖可用性依赖具备能力的 runner）。
- T64-6 实现相对 ① 设计简化（顺序式，见 ⑤ L2 补登记）。
- 偏离登记 D1：②b 交叉编译产物写至系统临时目录，违反 v41 §八 临时件位置闭集，主控核验无残留。

## 8. artifacts 围栏

以下为本片归档条目清单（路径 + 归属，TAB 分隔，ASCII 字典序）。

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/CHECKLIST_v42.md	repo
docs/t027/CHECKLIST_v42_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/acceptance-security-counterexamples.md	repo
docs/validation/acceptance-security-counterexamples_EN.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/00-briefing.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/01-arch/design.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/03-tests/red-baseline.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/04-impl/implementer-reports.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/t51-1-linux-escape.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-prereview/review.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/06-scan/review.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/07-final/review.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/08-gates/gate-readings.txt	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/09-glm/glm-1-bounddir-boundread-inventory.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/09-glm/glm-2-gates-chain-inventory.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/MANIFEST.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/SHA256SUMS.txt	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/freedom-list.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/input-list.md	repo
internal/adapters/bounddir_escape_unix_test.go	repo
internal/adapters/bounddir_escape_windows_test.go	repo
internal/adapters/bounddir_identity_test.go	repo
internal/chain/state_seal_counterexample_test.go	repo
internal/evidence/boundread_counterexample_test.go	repo
internal/gates/probe_permit_counterexample_test.go	repo
internal/gates/readiness_counterexample_test.go	repo
```

## 9. 计量

- R = 0（⑤ 首轮无 Medium+，未触发整改重跑；⑦ 首轮 PASS 未触发复审）。
- D = 1（D1：②b 交叉编译产物写至系统临时目录，违反 v41 §八 临时件位置闭集，主控核验无残留）。
- GLM 批量辅助实派 2 次；付费审查环：⑤ 1 次、⑥ 1 次、⑦ 1 次。

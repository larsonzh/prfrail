# T027-B 端到端接线（P69 `T027B-E2E-WIRING`）验证报告

日期：2026-11-12。状态：⑧b 定稿（⑨ 原生验证全绿后）。切片标识：`T027B-E2E-WIRING`（P69 · 层 2 治理侧 · 实现类 [SLICE] · M 档）。

## §1 执行摘要

- 目标：把 P44 的 readiness / permit 判定接到生产装配路径与 CLI 入口，收口 P57 §3.3 的 R-1。
- 结论：`DESIGN: DONE`（①）；接线落地 `internal/console/**` 5 个生产文件（4 改 1 增）+ 5 个测试文件（4 改 1 增）；验收五条反例全部有实测载体（③ 落 7 条测试）。
- 审查链：⑤ `PRE-REVIEW: PASS`（1 Low 已以注释级修复、④ 原位重跑全绿）；⑥ `INDEPENDENT SCAN: PASS`（0 发现）；⑦ `FINAL REVIEW: PASS`（0 发现；计量 = 终审 1 + 复审 0 + 盲审 0）。
- ⑦ 盲审不抽理由：本片不触所有权 / 停机 / 身份语义（§7.7 必抽不成立）；按比例抽查窗口自 v1.32 起、按 §12.4 表尾块「只计已收尾的正式切片」P67 = 第 1、P68 = 第 2、P69 = 第 3 ⇒ 未满 4 片，不抽。

## §2 接线点清单

内容取自 design-r1 §2.1（`docs/validation/evidence/P69-T027B-E2E-WIRING/01-arch/design-r1.md`）；增删行数为 `git diff --cached --numstat` 实测（起跑消息的 6 + 6 文件口径与变更集不符，按实测 5 + 5 记）。

| # | 路径 | 增/改/删 | 内容 | 影响 |
|---|---|---|---|---|
| 1 | `internal/console/readiness_gate.go` | 增 | `requireProbePermit` 三件套（Evaluate(GoalProbe) → Require → release/Complete）单一来源 | 新文件（+29） |
| 2 | `internal/console/cli.go` | 改 | `--probe-permit` 旗标 + 路由前 presence 判定（只判定不消费）+ `runUsage` + 两处调用点转发 permit + import gates | +32 / −4 |
| 3 | `internal/console/runtime_agent.go` | 改 | options 增 `ProbePermit` 字段；Step 0 判定（stub 后、Step 1 前）；Require（Step 1 全过后、Step 2 marker 前）+ `defer release` | +23 / −0 |
| 4 | `internal/console/runtime_supervised.go` | 改 | options 增同名字段；fresh 门（`SupervisedInteractionsLedgerPath` 不存在 = 首次起跑才 判定+Require+defer；可重入安全） | +23 / −0 |
| 5 | `internal/console/supervised_cmd.go` | 改 | `executeSupervisedRun` 签名加 permit 并赋 `options.ProbePermit`；`ExecuteSupervisedReturn` 不动 | +4 / −1 |
| 6 | `internal/console/readiness_gate_test.go` | 增 | 正例 1 + 反例 3 组共 7 条 | 新文件（+225） |
| 7 | 4 个既有测试文件 | 改 | 只补注入点、断言零改动（fixture `newAssemblyFixture`、harness `newSupervisedHarness`、`runtime_multitask_red_test.go` 10 处裸字面量、`p51SevenAgentRunnerArgs` 补旗标） | +32 / −16 |

## §3 调用顺序与失败路径

- AgentRunner：ctx → Step 0 判定（缺 permit ⇒ `ErrProbePermitMissing`，零痕迹）→ Step 1 前置校验 → `requireProbePermit`（Require）→ `defer release` → Step 2 marker → Step 3–8。
- Supervised：ctx → plan → RunRoot/VisibleSender 校验 → fresh 门（首启：判定+Require+defer；重入：跳过）→ marker → engine。
- 失败路径三情形：① 判定失败 ⇒ CLI/入口拒绝，零消费零痕迹；② permit 已被消费 ⇒ CLI 判定过（presence-only 为 P44 有意设计）、入口 Require 以 `ErrProbePermitConsumed` 拒绝、不双启动（O_EXCL 跨进程兜底）；③ 判定成功但 permit 未消费 ⇒ 只发生于「判定后、Require 前」的 Step 1 拒绝——fail-closed、未 spawn、额度完好；spawn ⇒ Require 已成功 ⇒ defer 必 Complete（无消费缺口）。
- 与既有机制的关系：Require 在 event-log marker 之前（拒绝的装配对 runRoot 零写入）；与任务间停机证明屏障 / R 发布无关（gate 只在装配入口）。

## §4 反例与变异检验

7 条测试（③ 实测，`internal/console/readiness_gate_test.go`）：

| # | 测试名 | 断言一句话 |
|---|---|---|
| 1 | `TestProbePermitConsumedAfterSuccessfulAssembly` | 正例：装配后 marker 落盘、磁盘重载 Require 得 `ErrProbePermitConsumed`、同一对象二跑得 Consumed 而非 InUse |
| 2 | `TestAgentRunnerAssemblyRefusesNilProbePermit` | nil permit ⇒ `ErrProbePermitMissing`，run-root marker 不存在 |
| 3 | `TestSupervisedAssemblyRefusesNilProbePermit` | supervised 首启缺 permit ⇒ 同哨兵、run-root marker 不存在 |
| 4 | `TestAgentRunnerAssemblyRefusesSpentProbePermit` | 已消费 permit ⇒ `ErrProbePermitConsumed`、无 marker、零 spawn |
| 5 | `TestSupervisedAssemblyRefusesSpentProbePermit` | supervised 同型：Consumed、无 marker |
| 6 | `TestRunRefusesMissingProbePermitFlag` | 缺 `--probe-permit` ⇒ 非零退出且 stderr 点名该旗标 |
| 7 | `TestRunRefusesIncompleteRecordsBeforeBurningPermit` | 缺 `--agent-cli` ⇒ 保留 noop-only 措辞且 permit marker 未生成 |

变异检验（M1–M5，取自 `docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt`）：

| 变异 | 对象 | 聚焦测试 | 红 | 恢复 sha256 一致 | 恢复后绿 |
|---|---|---|---|---|---|
| M1 | `runtime_agent.go`：Step 0 判定块 + `requireProbePermit`/defer 块（一起删） | `TestAgentRunnerAssemblyRefusesNilProbePermit` | FAIL | 一致（1100E2DE…C3F336） | PASS |
| M2 | `runtime_agent.go`：仅 `requireProbePermit` + defer 块 | `TestAgentRunnerAssemblyRefusesSpentProbePermit` | FAIL | 一致（1100E2DE…C3F336；首轮撞文件锁后重试成功） | PASS |
| M3 | `cli.go`：整个 CLI presence gate 块（保编译补 `_ = *probePermit`） | `TestRunRefusesMissingProbePermitFlag` | FAIL | 一致（5C8D6B15…620EE2B） | PASS |
| M4 | `cli.go`：`!inputs.complete()` 拒绝块 | `TestRunRefusesIncompleteRecordsBeforeBurningPermit` | FAIL | 一致（5C8D6B15…620EE2B） | PASS |
| M5 | `runtime_agent.go`：`defer releaseProbePermit()` → `_ = releaseProbePermit` | `TestProbePermitConsumedAfterSuccessfulAssembly` | FAIL | 一致（1100E2DE…C3F336） | PASS |

反例 4（既有 P44 gates 判定测试不变）：`internal/gates/**` 零 diff + `go test -count=1 ./internal/gates/...` ok。

## §5 审查链结论

- ① `DESIGN: DONE`（design-r1）；批量辅助 GLM-5.3-Flash 2 次（5546 / 5288 字符）。
- ⑤ r1 `PRE-REVIEW: PASS`（1 Low：release 闭包吞 Complete 错误不可观测——最小修复取注释级论证失败不可达，修复后 ④ 全量重跑全绿；7 Info 登记）。⑤ 额度 1/2（未用复审）。
- ⑥ r1 `INDEPENDENT SCAN: PASS`（Section B 无发现；Section E 残余 = 双进程并发抢占实机证据、平台耐久实机路径）。额度 1/3。
- ⑦ r1 `FINAL REVIEW: PASS`（0 发现；六项必审全过；残余风险同 ⑥ Section E）。计量 = 终审 1 + 复审 0 + 盲审 0。
- 五要素留痕（§7.8.5）：对象 = ⑦ 第 N 轮；轮次 = r1；依据 = §4.3 + §5.3 正常收敛；结论行原文 = `FINAL REVIEW: PASS`。

## §6 门禁读数（⑨ 定稿实测，2026-11-12，p9 批次；中间 ④ 批次读数见本片临时件，按 §5.1 于 ⑧c 清理）

- `gofmt -l .` 空；`go build ./...` exit 0；`go vet ./...` exit 0；`go test -count=1 ./...` 15 包 ok / FAIL=0。
- `GOOS=linux go build ./...` 与 `GOOS=linux go vet ./...` 空输出（Unix 零影响的本机证据；CI ubuntu 腿兜底）。
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`。
- 三作用域：`--scope=index` / `--scope=tree` 均 `TOTAL_FAIL=0`（`changed=38`；G4b `newCjk=405 unknown=[]`；G8-a/b OK；G5(f) 索引=工作区）；`--scope=ci` 本地为 `HEAD^1..HEAD` 退化范围（changed=2，`TOTAL_FAIL=0`），本片变更集覆盖由 index/tree 给出、真实 ci 读数待推送后由 CI 产生（仓库既有口径）。
- 归档口径：门禁读数未归档（可重跑获取），变异证据已归档。

## §7 已知边界

- noop 路径与 `ai check` 路径本片不接线（`ai check` 的 permit 生命周期需独立设计，属后续切片）。
- 正例测试平台感知：Windows 上 `agentRunnerPublicationDurability` 为 unproven，生产 dispatcher 在 R 前拒绝 ⇒ 本地断言落在「额度恰消费一次 + 二跑 Consumed」；durability proven 平台的严格成功分支由 CI ubuntu 腿运行。
- 双进程并发抢占同一 marker 的实机对打证据未覆盖（⑥ Section E 与 ⑦ 残余风险同指；O_EXCL 线性化由 P44 判定本体与既有 gates 测试承担，本片不改判定本体）。
- `ExecuteSupervisedReturn` 不受门控（归还/续跑非新起跑，额度已在首启消耗）。
- permit 注入形态与 P44 报告 :156 处方措辞的偏差（F-1）已如实登记。
- 临时件 D = 1（⑧a1 文档载体在切片目录内自删其一次性辅助脚本、预删计数未记——次序违规，⑩ 升级用户裁定；最终零残留回读通过）；R = 0。后续 ⑧a 文档载体在 `tmp/<切片ID>/` 内不得自删其辅助脚本，须由 ⑧c 统一清理。

## §8 明确未执行

未改准则正文（`DELIVERY_DIRECTIVE{,_EN}.md` 零改动）；未改 `internal/gates/**`、`internal/chain/**`、`internal/adapters/**`、`internal/tickets/**`、`internal/guard/**` 判定本体；未改 schema / fixtures / `.github/workflows/**` / `go.mod`；未改 CONTRACTS 语义（§4.2 末句时态更新属【五】第 10 项，见 ⑧a2）；未改 v50–v55 快照；未 commit / push；未启动新切片；未推 gitee。

artifacts 围栏：以下为本片归档条目清单（路径 + 归属，TAB 分隔，ASCII 字典序；38 条）。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/00-briefing.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/01-arch/design-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/02-batch/batch-01-seam-enumeration.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/02-batch/batch-02-run-callchain.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/04-reviews/final-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/04-reviews/scan-r1.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/MANIFEST.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/SHA256SUMS.txt	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/freedom-list.md	repo
docs/validation/evidence/P69-T027B-E2E-WIRING/input-list.md	repo
docs/validation/p44-gate-readiness.md	repo
docs/validation/p44-gate-readiness_EN.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/t027b-e2e-wiring.md	repo
docs/validation/t027b-e2e-wiring_EN.md	repo
internal/console/cli.go	repo
internal/console/cli_multitask_wiring_test.go	repo
internal/console/readiness_gate.go	repo
internal/console/readiness_gate_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_multitask_red_test.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_test.go	repo
internal/console/supervised_cmd.go	repo
```

sixstage 围栏：以下为 ⑥ 独立扫描 r1 的 Section A–E 逐字原文。

> ⑥ 结论行转录口径见 §6.3 的 G7-f（v1.48）。

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

(no findings)

## Section C

- The contract is matched exactly by the level-1 gate semantics in [internal/gates/readiness.go](internal/gates/readiness.go#L141-L170): `EvaluateT027BReadiness` is a pure function, and for `GoalProbe` it checks presence only, while the actual quota consumption is delegated to `ProbePermit.Require`. That matches the contract sentence: “EvaluateT027BReadiness is the only construction entry and is a pure function; zero I/O; no permit consumption.”
- The one-shot lifecycle is enforced in [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L1-L95) and [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L176-L244): the state machine is `Unused → Claimed → Spent`, the linearization point is the `O_CREATE|O_EXCL` marker file persisted with `Sync`, and `Complete` is idempotent only after a valid claim. This is the direct contract anchor for the “fail-closed” rule.
- The production wiring sits in the allowed scope only: [internal/console/readiness_gate.go](internal/console/readiness_gate.go#L1-L24), [internal/console/cli.go](internal/console/cli.go#L282-L395), [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L120-L190), and [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L259-L320). The flow is consistent with the design: CLI resolves and validates the permit before routing, AgentRunner evaluates before preconditions and consumes only after preconditions pass, and supervised starts enforce a fresh-start ledger gate before reusing the permit.
- The negative-path test bundle in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L34-L205) is the real behavioral proof for the five acceptance conditions: nil permit denied, spent permit denied, missing CLI flag denied, incomplete record set kept fail-closed, and successful assembly consumes the permit exactly once.
- The mutation evidence in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) is especially important: M1–M5 are all red when the guard is removed or weakened, then restored green after the original state is restored. That is exactly the correct “can this be falsified?” check.
- The recorded gate readouts are all green in the evidence bundle: [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-07-selftest.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-07-selftest.txt), [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-08-index.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-08-index.txt), [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-09-tree.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-09-tree.txt), and [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-10-ci.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-10-ci.txt). The index/tree/ci reports all say `TOTAL_FAIL=0`, and the selftest report is `737/737`.

## Section D

- If the CLI presence gate in [internal/console/cli.go](internal/console/cli.go#L295-L395) is removed:
  - The check in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L167-L205) goes red: `TestRunRefusesMissingProbePermitFlag` fails because the missing `--probe-permit` path becomes a later record-read or assembly error instead of a fail-closed refusal.
  - The “which check is not red?” answer is that the existing `--request-record` and “local runtime only supports noop steps” validations still remain red; the unique thing that becomes green wrongly is the missing-permit gate itself.

- If the `Require` claim is removed from the shared helper in [internal/console/readiness_gate.go](internal/console/readiness_gate.go#L8-L24) or from the assembly entry in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L146-L180):
  - `TestAgentRunnerAssemblyRefusesSpentProbePermit` in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L105-L166) goes red because presence-only evaluation would allow an already-spent permit to proceed.
  - The “not red” part is the initial `EvaluateT027BReadiness` presence check; that one is intentionally not enough to decide quota consumption, so it is not the guard that catches a spent permit.

- If the deferred `Complete` release is removed:
  - The mutation case M5 in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) goes red; the second run would report `ErrProbePermitInUse` instead of `ErrProbePermitConsumed`.
  - This is the exact “complete must be recorded” invariant; no other check catches it.

- If the supervised fresh-start gate in [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L269-L320) is removed:
  - A second supervised run on an existing interaction ledger can wrongly pass the permit through again; the re-entry behavior becomes fail-open.
  - The current evidence does not show a dedicated second-park mutation on a real ledger; the invariant is still logically covered by the ledger check and the design, but it is the one place where an OS-level concurrency test would be the strongest proof.

- If the stronger “reject before marker creation” ordering in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L146-L180) is removed:
  - The nil/spent permit tests in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L66-L166) are the ones that go red, because the run-root marker could be created before the permit failout.
  - The design is intentionally ordered to prevent a marker from appearing on a fail-closed refusal.

- The no-coverage invariant I would call out is the real cross-process race itself: the code relies on `O_CREATE|O_EXCL` plus `Sync`, but the artifacts do not contain a true concurrent two-process reproduction of the race. That is not a logic bug in this slice; it is the one remaining externality that would require a live two-process OS test to eliminate all doubt.

## Section E

- The things that are not directly observable from the read-only evidence:
  - A real two-process race on the same permit marker across separate OS processes.
  - A live Windows-specific agent-runner publication-durability path that depends on the actual pinned CLI and replay store behavior.
  - A full end-to-end run of the real CLI in a production environment outside the recorded artifacts.
  - The precise behavior of the Linux-side `GOOS=linux` path beyond the recorded pass states; the output files are empty, but the higher-level gate reports still confirm overall pass.

- What would overturn this PASS:
  - Any proof that a `nil` or already-spent permit is accepted into a real run.
  - Any proof that a second supervised entry with an existing ledger bypasses the fresh-start gate.
  - Any proof that the permit is consumed without `Require`, or that `Complete` is skipped on the success path.
  - Any diff outside the allowed scope that modifies the gate logic itself under [internal/gates](internal/gates) or alters the chain/adapters semantics.
  - Any real runtime failure where the mutation tests in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) no longer hold.

- The current evidence is strong enough to say the slice is fail-closed in the intended sense and stays within the declared contract. The remaining uncertainty is only the untested OS-level race behavior, not a logical defect in the gate ordering or the permit lifecycle.

INDEPENDENT SCAN: PASS
```

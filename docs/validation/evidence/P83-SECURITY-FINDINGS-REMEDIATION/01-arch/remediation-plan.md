# P83 SECURITY-FINDINGS-REMEDIATION — 修复方案书（01-remediation-plan.md）

- 切片：P83 `SECURITY-FINDINGS-REMEDIATION`；① 架构（deep-reasoner，DeepSeek V4 Pro Max）。
- 输入权威：`docs/validation/s2-pre-independent-security-review.md` §7 三面 Section B 表（逐字违反句）与 §9 F1–F8 表；
  `docs/CONTRACTS.md`（:147–148、:150–151、:212、:214、:216、:391）与 `docs/SECURITY.md`（:23、:60）。
- 本片不改 S1 Exit 门判据、不改 `SECURITY{,_EN}` 威胁模型基线、不改准则版本、不引入第三方依赖；
  代码标识符与注释一律英文；本方案书正文用中文。
- 全部 8 项修复的契约裁定均为 **DOES-NOT-TOUCH**（实现符合性修复，不改 CONTRACTS/SECURITY 措辞、不改 schema/fixtures）。
  逐项推理见各节 A.2。越界触碰与越界外新发现分别登记于 §2 与 §4。

## 速览

| # | 锚点（实测盘上行号） | 缺陷一句 | 契约裁定 | 核心测试 |
|---|---|---|---|---|
| F1 | `internal/gates/runner.go:208` | 产物按路径 Lstat 后重开路径读取，未绑定目录句柄 | DOES-NOT-TOUCH（SECURITY.md:23 实现落实） | `TestRunnerArtifactReadRefusesEscapingLink` |
| F2 | `internal/console/runtime_agent.go:1245` | 派发复核用装配期内存账本，不重读当前授权与预留 | DOES-NOT-TOUCH（CONTRACTS.md:212 实现） | `TestAgentRunnerDispatchRereadsRevocationAfterAssembly` |
| F3 | `internal/chain/review.go:116` | 策略批准只验形状，端口自算摘要即 approve | DOES-NOT-TOUCH（CONTRACTS.md:147 实现） | `TestBuildReviewResultPolicyApproveRequiresEvidenceRootBinding` |
| F4 | `internal/console/runtime_agent.go:731`、`runtime_supervised.go:1642` | 租约证据仅由 runID 派生 | DOES-NOT-TOUCH（CONTRACTS.md:150 实现） | `TestAgentRunnerPublishRefusesWithoutCurrentLease` |
| F5 | `internal/gates/runner.go:73` | 无 Effect scope 即跳过准入，占位 policyHash 可伪造 | DOES-NOT-TOUCH（CONTRACTS.md:391 实现） | `TestRunnerRefusesContradictoryEffectScope` |
| F6 | `internal/console/ai.go:94` | 读凭据前不校验接收端/子进程身份绑定 | DOES-NOT-TOUCH（SECURITY.md:60 实现） | `TestAICheckRefusesReceiverMismatchWithoutReadingStore` |
| F7 | `internal/guard/process_windows.go:457` | UTF-16 环境块无所有权、成功失败均不清零 | DOES-NOT-TOUCH（CONTRACTS.md:216 实现） | `TestWindowsStartClearsEnvironmentBlockOnSuccessAndFailure` |
| F8 | `internal/adapters/ai_probe_copilot_executor.go:53` | 子进程 `--log-dir` 落盘日志未证明脱敏 | DOES-NOT-TOUCH（CONTRACTS.md:214 实现） | `TestManagedCopilotCLIProbeExecutorNeverPassesLogDir` |

---

# §1 逐发现修复

## F1 产物读取绑定授权工作区目录句柄（Medium，面 1）

### A.1 现状代码事实
- 缺陷锚点：`internal/gates/runner.go:208`（`info, err := os.Lstat(path)`，其后 `:212` `data, err := os.ReadFile(path)`）。
- 控制流：`Runner.Run`(:41) → `assess`(:115) 仅当 `failure == ""` → `collectArtifacts`(:191)。
  `collectArtifacts` 对每个 glob 用 `filepath.Glob(filepath.Join(request.WorkspaceRoot, pattern))`(:194) 按**路径**展开，
  对末端 `Lstat` 后按**路径**重新 `ReadFile`：检查与读取不绑定同一对象，中间目录解析不受限
  （SEC-02a：中间 symlink 指向界外目录时内核跟随并返回界外文件），随后原文经
  `runner.Scanner.Scan` → `runner.Objects.PutObject` 入证据库。
- 现有可信原语：`internal/adapters/bounddir_{unix,windows,darwin}.go` 的 `BoundDir`/`OpenBoundDir`/
  `OpenBoundDirFile`（openat2 受限解析 / Windows 句柄绑定，失败 fail-closed），返回 `*evidence.BoundFile`
  （`internal/evidence/boundread.go`，单句柄 `ReadAll`/`Hash`/`Identity`）。P79 已实测该原语
  （`TestBoundDirMiddleSymlinkRefusedWithClosedErrno`）。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/SECURITY.md:23`「限制必须由 runner/OS 能力落实；不能仅靠提示词。」
  修复是把既有 OS 能力（句柄绑定读）接到缺失的读取入口上，属实现符合性；不需要契约措辞、schema、fixture 变更。

### A.3 修复设计（最小变更）
- 前置（越界，见 §2）：把 `BoundDir`/`OpenBoundDir`/`OpenBoundDirFile` 三平台实现从 `internal/adapters`
  迁入 `internal/evidence`（`evidence.BoundFile` 本就在 evidence），`adapters` 保留薄别名
  （`type BoundDir = evidence.BoundDir` 等）以免改动既有调用方。签名不变：
  ```go
  // internal/evidence（迁移后）
  func OpenBoundDir(path string) (*BoundDir, error)
  func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*BoundFile, error)
  func (dir *BoundDir) Close() error
  ```
- `internal/gates/runner.go` 内改造 `collectArtifacts`：
  ```go
  const gateArtifactMaxBytes int64 = 256 << 20
  var ErrArtifactOutOfWorkspace = errors.New("artifact resolves outside the bound workspace")
  func workspaceArtifactRel(root, path string) (string, bool) // filepath.Rel + 拒绝 ".."、绝对、分隔符穿越
  ```
  - 入口：`bound, err := evidence.OpenBoundDir(request.WorkspaceRoot)`；失败 → `return nil, ErrArtifactMissing`（拒绝收集）。
  - 每个 glob 命中：`rel, ok := workspaceArtifactRel(...)`（失败 → `ErrArtifactOutOfWorkspace`，映射为拒绝收集）；
    `file, err := evidence.OpenBoundDirFile(bound, rel, gateArtifactMaxBytes)`（失败 → `ErrArtifactMissing`）；
    `data, truncated, err := file.ReadAll()`（truncated=true → 拒绝）；`file.Close()`；去重键改用 `rel`。
  - 失败语义不变：`ErrArtifactMissing` → `failure = "artifact-missing"`；其余 → `"scan-failed"`。失败时零对象入库、零推进。
- fail-closed：目录打不开、相对化失败、句柄打开失败、越界/超限 → 全部拒绝收集，绝不回退为按路径读取。

### A.4 回归测试点
- 包/文件/测试名：`internal/gates`，新文件 `runner_artifact_bound_test.go`（unix 构建标签），
  `TestRunnerArtifactRefusesEscapingLink`。
- 断言：workspace 内 `sub/link -> <界外目录>`，hook `ArtifactGlobs: ["sub/link/secret"]`，executor 正常退出；
  断言 result `Assessment == "failed"` 且 `FailureKind == artifact-missing`、`Artifacts` 为空、`ObjectStore`
  收到的对象中**不含**界外 secret 的内容哈希。
- 变异检验：把 `OpenBoundDirFile` 换回 `os.ReadFile(path)` → 界外 secret 被收集入库，测试变红。

### A.5 闭合断言
- 「删掉句柄绑定读取（恢复 `os.ReadFile`）后，界外 secret 被收集进入证据库且任务推进；保留句柄绑定后，越界读取被拒绝，证据库不新增对象，任务不推进。」

## F2 launch 前向可信账本重读授权与预算（High，面 2）

### A.1 现状代码事实
- 缺陷锚点：`internal/console/runtime_agent.go:1245`：
  `LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return options.AuthorizationLedger, nil }`，
  相邻 `:1246` `CostLedger: options.CostLedger`。包围函数 `agentRunnerBindTasks`(:1199–:1276)。
- 控制流：CLI 在启动时一次性 `loadAuthorizationLedger` 与 `tickets.LoadCostLedger`（`cli.go:562/:585`）；
  派发时 `AgentRunnerReplayDispatcher.DispatchAgentRunner` 先发布 R（`Store.RecordRequest`）再
  `dispatcher.LedgerSnapshot()`（`agent_runner_replay_dispatcher.go:150`）→ `AgentRunnerLaunchReconfirmer`
  （`validateAgentRunnerAuthorization` + `validateAgentRunnerBudget`，`agent_runner_launch_reconfirmer.go`）。
  两个账本都是装配期内存副本：启动后另行持久化的撤销（`approvals revoke` 追加 JSONL）与跨进程结算不可见。
- 可信来源已有：`loadAuthorizationLedger(path)`（`internal/console/approvals.go:282`，逐行严格解码+校验）；
  cost 侧 `tickets.DecodeCostLedgerRecord` + `tickets.LoadCostLedger`（commit 为 temp+sync+rename 原子替换，
  `durable_cost_ledger.go:515`），跨进程结算按文件可见。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/CONTRACTS.md:212`「dispatch 在发布 R 之后、允许 launch 之前必须重新确认授权未被撤销且预算 reservation 仍 outstanding」。修复正是落实该句，无措辞/schema 变更。

### A.3 修复设计（最小变更）
- `internal/console/runtime_agent.go`：
  - `AgentRunnerAssemblyOptions` 增两字段：`AuthorizationLedgerPath string`、`CostLedgerPath string`
    （由 CLI 传入已解析绝对路径；零值回退到内存账本，保住既有测试与调用方）。
  - `LedgerSnapshot` 闭包改为：
    ```go
    LedgerSnapshot: func() (chain.AuthorizationLedger, error) {
        if options.AuthorizationLedgerPath != "" {
            return loadAuthorizationLedger(options.AuthorizationLedgerPath)
        }
        return options.AuthorizationLedger, nil
    },
    ```
  - 新助手（同文件）：
    ```go
    func loadCostLedgerSnapshot(path string) (*tickets.CostLedger, error)
    // os.ReadFile → tickets.DecodeCostLedgerRecord → tickets.LoadCostLedger；读失败即失败
    ```
- 越界小项（§2）：`internal/adapters/agent_runner_replay_dispatcher.go` 增
  `CostLedgerSnapshot func() (*tickets.CostLedger, error)` 字段；`DispatchAgentRunner` 在 reconfirm 前调用：
  非 nil 且失败 → 包 `ErrAgentRunnerLaunchReconfirmation`；为 nil 时回退 `dispatcher.CostLedger`（测试兼容）。
- fail-closed 顺序不变：R 先发布 → 账本重读 → 重确认 → 才 launch。读取失败 / 已撤销 / 已结算（reservation 不存在）→
  `ErrAgentRunnerLaunchReconfirmation`，**R 保留**（不删除、不重试），**零启动**。空账本（文件不存在）→ grant 匹配失败 → 阻断。

### A.4 回归测试点
- 包/文件/测试名：`internal/console`，`runtime_agent_e2e_test.go` 增
  `TestAgentRunnerDispatchRereadsRevocationAfterAssembly`。
- 断言：以含 grant 的 `--authorization-ledger` 文件装配；在装配完成、派发前由测试向该文件
  `appendAuthorizationRecord` 追加撤销记录；断言（1）replay store 保留已发布 R（请求记录在、无删除），
  （2）无 launch receipt / 无进程启动（`ErrAgentRunnerLaunchReconfirmation` 路径），（3）run 不推进。
- 变异检验：把闭包改回 `return options.AuthorizationLedger, nil` → 启动继续、测试变红。

### A.5 闭合断言
- 「把 LedgerSnapshot 改回装配期内存账本后，启动后追加的撤销不可见、launch 继续且 R 不变；保留重读后，撤销在 launch 前被发现，launch 被拒绝而 R 原样保留。」

## F3 发布前从执行者不可写的可信来源验证策略批准与候选证据根（High，面 2）

### A.1 现状代码事实
- 缺陷锚点：`internal/chain/review.go:116`（`case "policy": if !evidence.ValidHash(decision.PolicyHash)`，
  仅形状检查）。包围函数 `validateReviewDecision`(:92)，由 `buildReviewResult`(:36) 调用；引擎在
  `REVIEW_PENDING` 后调 `engine.options.Reviewer.Review`（`engine.go:254` 区域）。
- 生产端口：`internal/console/runtime_agent.go:669–687` 与 `runtime_supervised.go:1589–1607` 的 `Review`
  本地自算 `PolicyHash = digest(<domain>, label, "postflight-qualified-candidate")` 并 `Outcome: "approve"`：
  摘要合形不证明策略曾获独立批准，postflight 与批准互作原件。
- 可复用可信来源（执行者不可写）：
  - AgentRunner 路径：授权账本 JSONL（`--authorization-ledger`），仅由操作员命令
    `prfrail approvals <list|revoke>` 写入（`appendAuthorizationRecord`，`approvals.go:320`）；装配/运行时只有读取路径。
    grant 体含操作员签发的 `PolicyHash`/`RunManifestHash`/`Scope`（`authorization-record.schema.json`、`chain/authorization.go:70–97`）。
  - Supervised 路径：`<runRoot>/supervised/` 下操作员交互与风险确认台账
    （`operator-interactions.jsonl` / `risk-confirmations.jsonl`），由操作员执行
    `prfrail supervised return` 的 console 路径写入；受监督步骤（isolated agent）按写权限矩阵禁止写核心 store，无法伪造。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/CONTRACTS.md:147–148`「评审者与候选变更产生者是否同一身份、`evidenceRootHash`/`policyHash` 引用是否真实存在及 waiver 到期后不得复用，均由独立 checker 判定」。修复让「独立 checker = chain 核心」真正可判定（批准绑定 + 证据根绑定），不改 schema/fixtures/措辞。

### A.3 修复设计（最小变更）
- `internal/chain/review.go`（七文件内）：
  ```go
  // 独立 checker 侧：approve 必须携带候选证据根绑定摘要
  func ReviewEvidenceBinding(candidateSnapshotHash, evidenceRootHash string) string
  // = evidence.Digest("proofrail:review-evidence-binding:1\n", []byte(candidateSnapshotHash+"\n"+evidenceRootHash))
  ```
  `validateReviewDecision` 的 `approve` 分支增加：
  `if !slices.Contains(decision.Evidence, ReviewEvidenceBinding(request.CandidateSnapshotHash, request.EvidenceRootHash))`
  → `fmt.Errorf("%w: review approve must bind the candidate evidence root", ErrInvalidState)`。
- `internal/console/runtime_agent.go` `Review`（七文件内）：
  - 重读可信账本：`runtime.authorizationLedgerPath != ""` 时 `loadAuthorizationLedger(path)`，否则用装配内存账本。
  - 寻找覆盖本 run/task/step 的 grant：`ledger.GrantRecordsFor` 匹配 `Grant.RunID == request.RunID` 且
    `Scope.TaskIDs/StepIDs` 包含请求；`chain.EvaluateGrant` 为 `GrantActive`；`grant.PolicyHash` 必须 `evidence.ValidHash`。
  - 无匹配 → `return chain.ReviewDecision{}, fmt.Errorf("%w: no operator-approved policy grant covers this review", ErrInvalidConfig)` → 引擎 `markTaskForRepair("review-failed")` → 拒绝、停在待人工状态，零发布。
  - 有匹配 → `PolicyHash: grant.PolicyHash`（不再自算摘要），
    `Evidence: []string{chain.ReviewEvidenceBinding(request.CandidateSnapshotHash, request.EvidenceRootHash), digest(agentRunnerReviewEvidenceDomain, ...)}`。
- `internal/console/runtime_supervised.go` `Review`（七文件内）：
  - 读取操作员确认台账（`risk-confirmations.jsonl`）与 session；无对应 run/task/step 确认 → 同上拒绝。
  - `PolicyHash: digest(supervisedReviewPolicyDomain, session.HandoffPolicyHash, confirmationHash)`；
    `Evidence` 含 `chain.ReviewEvidenceBinding(...)`。
- 可信来源属性（F3 明确要求）：授权账本/确认台账都由操作员 CLI 路径追加写入，运行中的装配/执行体
  （managed agent、supervised agent）没有任何写该文件的代码路径 → 满足「执行者不可写」；chain 侧校验使
  未绑定证据根的 approve 无法形成 review receipt。

### A.4 回归测试点
- 包/文件/测试名：
  - `internal/chain`，`review_test.go` 增 `TestBuildReviewResultPolicyApproveRequiresEvidenceRootBinding`：
    合法 approve 决策但 `Evidence` 不含 `ReviewEvidenceBinding` → 期望 `errors.Is(err, ErrInvalidState)` 且零 receipt。
    变异：删除 `slices.Contains` 检查 → approve 通过、测试红。
  - `internal/console`，`runtime_agent_e2e_test.go` 增 `TestAgentRunnerAssemblyRefusesReviewWithoutOperatorGrant`：
    账本无覆盖 grant → 装配在 review 处拒绝，无 promotion receipt、任务不 PASS。变异：把端口恢复为自算摘要 approve → 红。
- 断言：「不发布、不推进、证据仍完整」。

### A.5 闭合断言
- 「移除 Review 端口对操作员 grant/确认记录的重读绑定后，无批准的 run 仍产出 approve review receipt 并推进；保留绑定后，review 被拒绝，无 promotion receipt 产生，候选与证据原样保留。」

## F4 发布前核对排他租约并绑定真实记录（Medium，面 2）

### A.1 现状代码事实
- 缺陷锚点：`internal/console/runtime_agent.go:731` 与 `internal/console/runtime_supervised.go:1642`：
  `LeaseEvidence: []string{digest(agentRunnerLeaseDomain, runtime.runID)}`（supervised 同形）——只由 runID 派生，不证明当前独占。
- 包围函数：两处 `Publish`（runtime_agent.go:688–737；runtime_supervised.go:1625–1652）。
- 可信机制已有：`internal/guard/lease.go` `LeaseStore`（磁盘租约，`Claim`/`Current`/`Validate`/`Release`，
  `FencingToken{ClaimID, Generation}`，租约体带 `LeaseHash`（`proofrail:writer-lease:1`），`Current` 全量校验）。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/CONTRACTS.md:150–151`「`promotion-receipt.schema.json` 将……停机与租约证据绑定为单个发布事务」。`LeaseEvidence` 已是 `[]string`，把真实租约记录哈希绑入即符合 schema；无措辞/schema 变更。

### A.3 修复设计（最小变更）
- 共用（同包，写在 `runtime_agent.go`）：
  ```go
  const publishLeaseTTL = 24 * time.Hour
  var ErrPublishLeaseUnproven = errors.New("publish lease proof missing or inconsistent")
  func publishLeaseRoot(runRoot string) string                       // filepath.Join(runRoot, "leases")
  func publishLeaseResourceID(runID string) string                   // "ws-" + hex(证据摘要) 前缀 24 位
  func newPublishLeaseStore(runRoot string) (*guard.LeaseStore, error) // guard.NewLeaseStore(publishLeaseRoot(runRoot))
  ```
- AgentRunner：`ExecuteAgentRunnerRun` 在引擎前 `Claim(ctx, resourceID, "proofrail-publisher", "attempt-"+strconv.Itoa(attempt), guard.ProcessIdentity{PID: os.Getpid(), StartToken: runID}, publishLeaseTTL)`；
  token 存入 `agentRunnerRuntime.leaseToken`（单次 CLI 调用内，同进程直到 Publish）。
  `Publish`：`lease, err := store.Current(ctx, resourceID)`；err 非 nil、`lease.Token != runtime.leaseToken`、
  `lease.LeaseExpiresAt` 早于当前时钟 → 均返回 `Outcome: "failed"` + `ErrorEvidence` 摘要（零完成回执）；
  通过则 `LeaseEvidence: []string{lease.LeaseHash}`。
- Supervised：`ExecuteSupervisedRun` freshStart 时 Claim，token 持久化进 session 文件（store-local 协议，非 wire schema）；
  `ExecuteSupervisedReturn`（重入、新进程）读回 token 后 `store.Validate(ctx, resourceID, token)`；
  `Publish` 同上绑定 `lease.LeaseHash`。缺失/不一致 → 拒绝。

### A.4 回归测试点
- 包/文件/测试名：`internal/console`，`runtime_agent_e2e_test.go` 增 `TestAgentRunnerPublishRefusesWithoutCurrentLease`。
- 断言：装配后、Publish 前用第二个 `LeaseStore` 同 resourceID 成功 `Takeover`/改写租约（或直接删除租约文件）；
  断言 Publish `Outcome == "failed"`、无 completed receipt、`LeaseEvidence` 未出现在任何完成回执。
  变异：改回 `digest(agentRunnerLeaseDomain, runID)` → Publish 完成、测试红。
- 附带：`runtime_supervised_test.go` 增 `TestSupervisedReturnRefusesDriftedPublishLease`（同思路，验证重入路径）。

### A.5 闭合断言
- 「删除 Publish 前的租约核对（恢复 runID 派生）后，租约已被他人持有的 run 仍产出 completed promotion receipt；保留核对后，租约缺失/不一致即拒绝，不产出完成回执，候选状态不变。」

## F5 执行前从冻结 hook 策略导出副作用分类（Medium，面 2）

### A.1 现状代码事实
- 缺陷锚点：`internal/gates/runner.go:73` `if classification != nil && AdmitEffect(classification.Record) != EffectAdmit`。
  上游 `:69` `buildEffectClassification(request, attempted)` 当 `request.Effect == nil` 直接返回 nil → **准入检查整体被跳过**、
  `Execute` 照跑。另 `:268–270`：scope.PolicyHash 非法时用 `proofrail:effect-policy-placeholder:1` **伪造**分类摘要。
- 冻结 hook 策略（`request.Hook`）本已强制 `EffectClass` 非空（`policy.go validateHookEffectPolicy`），但
  `buildEffectClassification` 只读 `request.Effect`（scope），二者矛盾（scope.EffectClass != hook.EffectClass）无检查。
  现调用方 `runtime_supervised.go:1491`（after-return gates）构造 Request 时不带 `Effect`。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/CONTRACTS.md:391`（PC-04 行）「只读/本地可丢弃/外部写分离；S1 外部写拒绝，未知副作用暂停对账」。实现层落实即可；无措辞/schema 变更。

### A.3 修复设计（最小变更，全部在 `internal/gates/runner.go`）
```go
var (
    ErrEffectClassificationUnavailable = errors.New("effect classification unavailable")
    ErrEffectClassificationConflict    = errors.New("effect classification conflicts with the frozen hook policy")
)
func deriveEffectClassification(request Request, createdAt time.Time) (*EffectClassificationRecord, error)
// 恒从冻结 hook 策略导出分类；request.Effect 存在时仅作交叉核对：
//   scope.EffectClass != hook.EffectClass → ErrEffectClassificationConflict（矛盾，拒绝）
//   scope.PolicyHash 非法 → ErrEffectClassificationUnavailable（删除占位伪造；hook.PolicyHash 已在 Prepare 校验）
//   OperationKey := scope.OperationKey；scope 缺失时 = digest("proofrail:effect-operation-key:1\n", HookDefinitionHash)
//   EffectClass/RecoveryGuarantee/PolicyHash/AuthorizationHash/ExternalSystems/Evidence 取 hook 字段（scope 覆盖同名非空字段）
func operationKeyFor(request Request) string
// buildEffectObservation 改用 operationKeyFor(request)，不再要求 request.Effect 非 nil
```
- `Run` 改为：`classification, err := deriveEffectClassification(...)`（err → 直接返回，不进 Execute）；
  `if AdmitEffect(classification.Record) != EffectAdmit`（**无条件**）→ 既有 deny 路径
  （`start-failed` + effect-admission-denied 证据 + `ErrGateFailed`）。
- fail-closed：分类缺失（不可导出）/矛盾/不允许（external-write、未知类）→ 一律不进入 `Execute`；观察与准入判定同一来源。

### A.4 回归测试点
- 包/文件/测试名：`internal/gates`，`runner_test.go` 增三项：
  - `TestRunnerDerivesClassificationFromFrozenHookPolicy`：hook.EffectClass=read-only、request.Effect=nil →
    执行发生且 result 带 `EffectObservationHash`。变异：恢复 `classification != nil` 短路 → 观察缺失、红。
  - `TestRunnerRefusesContradictoryEffectScope`：hook=read-only、scope=external-write → executor 零调用、
    `start-failed`、denied 证据。变异：删矛盾检查 → executor 被调用、红。
  - `TestRunnerRefusesInvalidScopePolicyHashWithoutPlaceholder`：scope.PolicyHash 非法 → 不进 Execute。
    变异：恢复占位摘要 → 红。
- 断言口径：「不执行、不推进、证据仍完整」。

### A.5 闭合断言
- 「删除 scope 与冻结 hook 策略的矛盾检查后，分类矛盾的请求进入 Execute；保留检查后，矛盾/缺失/不允许分类一律拒绝，Execute 零调用，effect observation 不产生。」

## F6 读取凭据前校验已批准 profile、凭据引用、接收端与子进程身份绑定（High，面 3）

### A.1 现状代码事实
- 缺陷锚点：`internal/console/ai.go:94` `record, err := cli.checkAIAvailability(ctx, profile, *channel, ...)`。
  包围函数 `executeAICheck`(:33–:137)。上游：`configuredAIProfile`(:230) 取算子配置 profile →
  `cli.checkAIAvailability`（= `defaultAICheck`(:222)）→ `adapters.AIAvailabilityChecker.Check`
  （`ai_preflight.go:31`）：先 `Secrets.Resolve` 再 `Probe.Probe`（`ai_probe_copilot.go:60`）→
  `ManagedCopilotCLIProbeExecutor.Execute`（`ai_probe_copilot_executor.go:36`）直接
  `guard.RunManaged` spawn `invocation.Executable`。
- 缺口：profile 只做形状校验；`--copilot` 可执行文件与 profile/凭据引用之间没有任何批准绑定；
  秘密解析发生在任何接收端身份校验之前。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/SECURITY.md:60`「网络默认拒绝，例外授权须列 host/port/purpose、提供方数据策略及期限」。修复把「例外授权」落实为算子配置 + 算子 CLI 调用参数的接收端绑定（可信本机控制路径 = 操作员启动的 `prfrail ai check` 进程），并以绑定摘要固化；不改基线句、无新 wire record、无 schema 变更。

### A.3 修复设计（最小变更）
- `internal/console/ai.go`：
  ```go
  const aiProbeExecutableMaxBytes int64 = 64 << 20
  var ErrAIReceiverBinding = errors.New("AI probe receiver binding invalid")
  func validateAIProbeBinding(profile adapters.AIProviderProfile, channel, executablePath string) (string, error)
  // 1) adapters.ValidateAIProviderProfile(profile)（含 authMode/secretRef 引用合法性、HTTPS base URL）
  // 2) evidence.OpenBoundFile(executablePath, aiProbeExecutableMaxBytes).Hash("proofrail:ai-probe-executable:1\n") → executableHash
  // 3) profileHash := adapters.AIProviderProfileHash(profile)
  //    bindingHash := evidence.Digest("proofrail:ai-probe-receiver-binding:1\n", profileHash+"\n"+channel+"\n"+executableHash)
  // 任一步失败 → error（不读 SecretStore、不启动）
  ```
  `executeAICheck` 在 `cli.checkAIAvailability` **之前**调用；失败 → `writeCommandError`（exitFailure）。
  `defaultAICheck` 构造 `ManagedCopilotCLIProbeExecutor{ExpectedExecutableHash: executableHash}`。
- `internal/adapters/ai_probe_copilot_executor.go`（七文件内，与 F8 共文件）：
  ```go
  func CopilotCLIExecutableHash(path string) (string, error) // evidence.OpenBoundFile(path, 64<<20).Hash("proofrail:ai-probe-executable:1\n")
  type ManagedCopilotCLIProbeExecutor struct { ..., ExpectedExecutableHash string }
  ```
  `Execute` 在 spawn 前：`ExpectedExecutableHash != ""` 时重算实际哈希；不符 → 直接返回
  `CopilotCLIProbeExecution{Outcome: "start-failed", ...}`（**零启动**，环境缓冲仍按既有 defer 清零）。
- 子进程身份绑定：`guard.ProcessSpec{Command: invocation.Executable}` 只 spawn 该绝对路径；绑定摘要覆盖
  profile→channel→executable 三元组，与既有 `ai-availability.schema.json` 的 `ProfileConfigHash/Channel/Evidence` 相容。

### A.4 回归测试点
- 包/文件/测试名：`internal/console`，`ai_test.go` 增 `TestAICheckRefusesReceiverMismatchWithoutReadingStore`。
- 断言：`cli.secretManager` 用记录型 stub；`--copilot` 指向与批准绑定不符的文件（或 profile baseURL 非法）→
  命令失败退出；断言 stub 的 `Resolve` 零调用、`--out` 记录文件未创建。
- 变异：删除 `validateAIProbeBinding` 调用 → `Resolve` 被调用、测试红。
- 辅助：`internal/adapters`，`ai_probe_copilot_test.go` 增 `TestManagedCopilotCLIProbeExecutorRefusesHashMismatch`（假执行器断言零 spawn）。

### A.5 闭合断言
- 「移除入口绑定校验后，接收端不匹配的请求仍触发 SecretStore 解析并启动子进程；保留校验后，不匹配即不读凭据库、不启动，且不产出任何 availability 记录。」

## F7 Windows 原生环境块所有权与清零（Medium，面 3）

### A.1 现状代码事实
- 缺陷锚点：`internal/guard/process_windows.go:456–469` `buildEnvironmentBlock(env []string) (*uint16, error)`：
  返回 `&block[0]`，调用方 `start`(:214–:295) 将其传给 `CreateProcessW` 后**从不清零**；成功与失败路径都遗留含秘密
  （如 `COPILOT_PROVIDER_API_KEY`）的 UTF-16 缓冲。中间 `syscall.UTF16FromString` 的转换副本并入 `block` 同一 backing array。
- 相邻审计（发现要求的「临时副本核对」）：`AIAvailabilityChecker.Check` 的 `defer clear(credential)`
  （`ai_preflight.go:43`）、`copilotCLIManagedEnvironment` 的 buffer 清零（executor :57–62）与
  `executeSecretSet` 的 `defer clear(secret)`（`secret.go:56`）**均已清零**；唯一缺口就是本环境块。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/CONTRACTS.md:216`「adapter 只接收解析后的短暂字节副本，使用后清零」。实现层补齐；无措辞/schema 变更。

### A.3 修复设计（最小变更，仅 `internal/guard/process_windows.go`）
```go
// 所有权转移：返回可变切片，调用方负责清零；不再返回裸指针
func buildEnvironmentBlock(env []string) (block []uint16, err error)
// start 内：
environment, err := buildEnvironmentBlock(spec.Env)
if err != nil { return nil, err }
defer clear(environment)                      // 成功与失败路径统一清零
var environmentPtr *uint16
if len(environment) > 0 { environmentPtr = &environment[0] }
// CreateProcessW(..., uintptr(unsafe.Pointer(environmentPtr)), ...)
// platformProcess 增未导出测试缝：envClearHook func(block []uint16)（仅 windows 测试使用）
```
- `clear(environment)` 对 `[]uint16` 逐元素置零（Go 1.22 `clear` 语义），覆盖全部秘密副本；
  任何后续早退（optionalUTF16/handles/attributes 失败、CreateProcessW 失败）均经 defer 清零。

### A.4 回归测试点
- 包/文件/测试名：`internal/guard`，`process_spawn_ownership_windows_test.go`（windows 构建标签）增
  `TestWindowsStartClearsEnvironmentBlockOnSuccessAndFailure`。
- 断言：通过 `envClearHook` 捕获清零时点与内容——成功 spawn 与注入失败（attributes 构造失败）两条路径，
  hook 均被调用且观察到的块内 canary（如 `CANARY_SECRET=...`）在 hook 时点已全部为 0。
- 变异：删除 `defer clear(environment)` → hook 永不触发或块非零、测试红。

### A.5 闭合断言
- 「删除 deferred clear 后，持有 canary 的 UTF-16 环境块在启动成功与失败路径上均保持非零；保留 clear 后，两条路径上的环境块都在 spawn 返回后被清零。」

## F8 关闭未证明脱敏的子进程落盘日志（Medium，面 3）

### A.1 现状代码事实
- 缺陷锚点：`internal/adapters/ai_probe_copilot_executor.go:53–54`：
  `usagePath := ...` 与 `args := append(..., "--usage-output-file", usagePath, "--log-dir", logsDirectory)`；
  其上 :45–49 创建 `logsDirectory`。包围函数 `ManagedCopilotCLIProbeExecutor.Execute`(:36–:121)。
  子进程把 transcript 写入 `--log-dir`，目录清理（`defer os.RemoveAll(temporaryDirectory)`）不能证明从未落盘；
  stdout/stderr/usage 原文摘要直接成为 probe evidence（内容未脱敏即被摘要进证据链）。
- `usage.json` 仅含结构化计数，经 `parseCopilotCLIProbeOutput` 只持久化结构化字段（P42 C3 有界读）。

### A.2 契约裁定
- DOES-NOT-TOUCH。违反句 `docs/CONTRACTS.md:214`「API key/token 不得进入配置、摘要、日志、错误或证据……」。删除落盘日志路径、只保留结构化事实即符合；无措辞/schema 变更。

### A.3 修复设计（最小变更，仅 `internal/adapters/ai_probe_copilot_executor.go`）
- 删除 `logsDirectory` 创建与 `"--log-dir", logsDirectory` 参数；保留 `--usage-output-file`（结构化、有界读、解析后清除，目录随 probe 结束移除）。
- 证据口径收窄：stdout/stderr 仅保留内存缓冲与摘要，**不再存在任何子进程可写的落盘流**；
  `usage.json` 解析失败 + stdout 无结构化 usage → 既有 `response_invalid` 拒绝（即 BYOK 探针拒绝，
  `unavailable`，不产 available 记录）。
- 脱敏保证表达：当无法证明子进程不落盘（例如未来 CLI 强制要求 `--log-dir`）时，执行器必须 fail-closed 拒绝该探针，不得静默放行——本片以删除该参数实现，并以 argv 断言固化。

### A.4 回归测试点
- 包/文件/测试名：`internal/adapters`，`ai_probe_copilot_test.go` 增 `TestManagedCopilotCLIProbeExecutorNeverPassesLogDir`。
- 断言：用记录型假执行器/测试 helper 运行 executor；断言 spawn 的 argv **不含** `--log-dir` 与任何日志目录参数，
  临时目录内无 `logs` 目录，usage 结构化字段仍被持久化。
- 变异：在 argv 拼接处重新加入 `--log-dir` → 测试红。

### A.5 闭合断言
- 「在 argv 重新加入 `--log-dir` 后，子进程在磁盘上写出日志目录且测试观察到落盘；移除后不存在任何子进程可写的日志路径，证据仅由经校验的结构化 usage 事实构成。」

---

# §2 交叉项（B）

## 2.1 共享助手机会
1. `loadAuthorizationLedger(path)`（`console/approvals.go:282`，既有）：F2（派发重读）与 F3（评审重读）共用。
2. `loadCostLedgerSnapshot(path)`（新，`runtime_agent.go`）：F2 预算重读专用；若未来 DurableCostLedger 只读路径可复用再收敛（见 §4 OB 候选 5）。
3. `chain.ReviewEvidenceBinding`（新，`chain/review.go`）：F3 两端口共用、chain 校验共用。
4. `evidence.BoundDir`/`OpenBoundDirFile` 迁移：F1 专用，同时为未来所有「目录句柄绑定读」提供单一原语。
5. `CopilotCLIExecutableHash` + `ManagedCopilotCLIProbeExecutor.ExpectedExecutableHash`：F6 与 F8 共用（同一文件 `ai_probe_copilot_executor.go`）。
6. `publishLeaseRoot/publishLeaseResourceID/newPublishLeaseStore/ErrPublishLeaseUnproven`（新，`runtime_agent.go`）：F4 两条装配路径共用（同包，`runtime_supervised.go` 直接引用）。

## 2.2 修复顺序
F2 → F3（都依赖授权账本路径灌入 options 与 runtime；F2 先行确定 options 字段与重读助手）→ F6（独立，与 F8 共文件先后无依赖）→
F1（依赖 BoundDir 迁移完成）→ F4（独立，依赖 guard.LeaseStore 既有能力）→ F5（纯 `runner.go` 内）→ F7 → F8。
该顺序同时满足任务给定的优先级 F2,F3 → F6 → F1,F4,F5,F7,F8。

## 2.3 越界触碰（超出七文件的生产代码，逐项登记）
| 越界项 | 文件 | 内容 | 所属修复 |
|---|---|---|---|
| X-1 | `internal/adapters/bounddir_{unix,windows,darwin}.go`（及 `bounddir_escape_unix_test.go` 包名） | `BoundDir`/`OpenBoundDir`/`OpenBoundDirFile` 迁入 `internal/evidence`，adapters 留薄别名 | F1 |
| X-2 | `internal/console/cli.go`（:555–:585 装配处） | `AgentRunnerAssemblyOptions` 增 `AuthorizationLedgerPath`/`CostLedgerPath` 两字段并传入 | F2/F3 |
| X-3 | `internal/adapters/agent_runner_replay_dispatcher.go` | 增 `CostLedgerSnapshot func() (*tickets.CostLedger, error)` 字段并在 reconfirm 前调用（nil 回退内存账本） | F2 |
| 说明 | 上述三项均属既有机制的接线/迁移，不改变协议语义；其余 6 项修复（F3 的 chain 侧 + F4/F5/F6/F7/F8）全部落在七文件内。 | | |

---

# §3 风险清单（C）

| 修复 | 可能被打破的既有测试 | 原因 |
|---|---|---|
| F1 | 预期零破坏；`TestRunnerPersistsPassedResultAndArtifacts` 等 | 界内常规文件经句柄读取路径行为不变；唯一新增约束是 256 MiB 单文件上限与界外拒绝（本就要拒绝）。X-1 迁移若别名不完整会先由编译抓住。 |
| F2 | `runtime_agent_e2e_test.go`、`agent_runner_replay_dispatcher_test.go` | 采用「路径为空/快照函数为 nil 即回退内存账本」保住既有装配；dispatcher 新字段必须带 nil 回退，否则既有直接构造会红。 |
| F3 | `chain/review_test.go`（`TestBuildReviewResultApprove` 等）、`engine_test.go` 的 fakeReviewer、`runtime_agent_e2e_test.go` `TestAgentRunnerAssemblyFullChain`、`runtime_multitask_*` stub | approve 决策的 `Evidence` 必须新增 `ReviewEvidenceBinding`；无 grant/确认记录的装配现在拒绝 → 相关夹具需补账本/确认记录。这是预期的红转绿，不是行为回退。 |
| F4 | `runtime_agent_e2e_test.go`、`runtime_supervised_test.go` 的 Publish/return 全链 | 新增 runRoot/leases 租约与 session 令牌字段；supervised 旧 session 夹具缺令牌 → fail-closed 拒绝，夹具需升级。 |
| F5 | `gates/runner_test.go` 与 `runtime_supervised.go:1491` after-return gates 夹具 | `request.Effect == nil` 现在会导出分类并**无条件**过准入：read-only/local-discardable 照常执行；任何 external-write 类 hook 现在被拒（符合发现要求，夹具如含此类需改）。 |
| F6 | `console/ai_test.go` 全部 `ai check` 用例 | `--copilot` 现要求为可哈希的真实文件且通过绑定校验；用假路径 + stub checker 的用例需改指向临时可执行文件。 |
| F7 | Windows 构建下的 `process_spawn_*` 测试 | `buildEnvironmentBlock` 签名改为返回切片（同包内部），调用点唯一；风险低，若遗漏 `&block[0]` 空切片防护会 panic，由测试抓住。 |
| F8 | `ai_probe_copilot_test.go`（`TestCopilotCLIProbeBuildsIsolatedDeepSeekSmoke` 等） | 假子进程 argv 不再含 `--log-dir`；任何断言该参数存在的用例需更新。 |

---

# §4 越界外新发现（候选 OB，不属于本片）

1. **OB 候选**：`internal/console/runtime_noop.go` `localRuntime.Review`（manual/approve 自证）与 `Publish` 的
   runID 派生 `LeaseEvidence` —— 与 F3/F4 同类，但文件不在本片七文件，本片不修。
2. **OB 候选**：`runtime_supervised.go:1069` handoff session 的 `LeaseEvidence: digest(supervisedSessionLeaseDomain, runID, taskID)` —— 同 F4 类 runID 派生缺陷。
3. **OB 候选**：`internal/adapters/ai_preflight.go` `AIAvailabilityChecker.Check` 自身不校验接收端绑定（本片把校验放在 console 入口），绕过 `ai check` 直接调用 checker 的路径不受保护 —— 建议把 `ReceiverBinding` 下沉为 checker 必填字段。
4. **OB 候选**：F2 残余竞态 —— 重读账本后、spawn 前仍存在撤销窗口；PC-03 的「在途停机」处置不在本片范围。
5. **OB 候选**：F2 预算重读用普通文件读取，未经 `DurableCostLedger` 的 scope 锁与已确认前缀（`RequireOutstandingReservation`）；依赖原子替换保证一致性，非原子写入方会引入读竞态 —— 建议改用 durable 只读路径。
6. **OB 候选**：F6 的可执行文件哈希与 spawn 之间的 TOCTOU（文件被替换）—— 建议复用 `AgentRunnerPinnedCLILauncher` 的 identity 机制或按句柄校验。
7. **OB 候选**：F7 范围外的秘密相邻缓冲（`commandLine`、`ProcessSpec.Env` 的 Go 字符串、`startupInfo`）—— Go 字符串不可清零，需评估 Env 传递形态。
8. 前置缺口（非发现）：生产 CLI 尚无 `approvals grant` 签发命令；F3 修复后 AgentRunner 生产运行若无操作员 grant 将 fail-closed 停在 review。是否在本片外补签发路径需用户裁定。

---

# §5 升级事项（blocker）

- **无硬阻断项**。三件需用户知悉/裁定的点：
  1. §4-8：F3 使 AgentRunner 生产路径依赖操作员授权账本中存在覆盖 grant；当前无 grant 签发 CLI → 修复后无 grant 即停在待人工批准（符合发现要求，但影响可用性）。
  2. §2.3 X-1：F1 依赖把 `BoundDir` 从 adapters 迁入 evidence，触碰三平台文件；按边界规则登记后执行。
  3. F6 的「已批准 profile」以算子配置 + 算子 CLI 调用参数为可信本机控制路径；若要求更硬的可追溯批准台账（新 record 类型），须先走协议修订（CONTRACTS + schema + fixtures），本片按 DOES-NOT-TOUCH 处理。

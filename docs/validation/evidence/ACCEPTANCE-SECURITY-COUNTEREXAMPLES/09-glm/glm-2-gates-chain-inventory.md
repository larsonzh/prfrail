# GLM 批量辅助 2/2：GAP-60/64/70 代码引用面清查原文

载体：runSubagent(agentName="deep-reasoner", model="GLM-5.3-Flash (glm)")。只读事实枚举，逐字归档。

GLM-2/2: DONE

## 7. internal/gates/ 目录 .go 枚举
adapter.go / adapter_test.go / effects.go / effects_test.go / executor.go / executor_test.go / models.go / policy.go / policy_test.go / probe_permit.go / probe_permit_test.go / readiness.go / readiness_test.go / result.go / runner.go / runner_test.go / scanner.go（另有 .gitkeep 非 .go）

## 1. internal/gates/probe_permit.go
A 错误哨兵
- L66 `ErrProbePermitMissing = errors.New("probe permit missing")`
- L69 `ErrProbePermitInUse = errors.New("probe permit already in use")`
- L71 `ErrProbePermitConsumed = errors.New("probe permit already consumed")`
B 公开签名
- L126 `func NewProbePermit(id, path string, maxRequests int64) *ProbePermit`
- L140 `func LoadProbePermit(path string) (*ProbePermit, error)`
- L210 `func (p *ProbePermit) Require(ctx context.Context) error`
- L253 `func (p *ProbePermit) Complete() error`（非导出：L271 refusal、L284 persistClaim）
B Load fail 分支
- L142 空路径→`fmt.Errorf("%w: probe permit marker path is required", ErrProbePermitMissing)`
- L148-150 NotExist→`return permit, nil`（Unused，非错误）；L165-166 marker 存在→state=Spent 后 `return permit, nil`
- L152 非 NotExist 读失败→`fmt.Errorf("load probe permit %q: %w", path, err)`【非 NotExist 读失败】
B Require fail 分支
- L212 nil permit→ErrProbePermitMissing；L215 `ctx == nil`→ErrProbePermitMissing【ctx 非 nil】
- L218 预检 `ctx.Err()`→透传不包哨兵；L225 fast-reject claimed/spent→`p.refusal()`
- L233 persist 失败→state=Spent；L235-236 `fs.ErrExist`→ErrProbePermitConsumed；L238 其它 persist 失败→ErrProbePermitConsumed；L240 成功→state=Claimed；L245 persist 后 `ctx.Err()`→透传（不退款）
- refusal：L274 Claimed→ErrProbePermitInUse；L276 Spent→ErrProbePermitConsumed；L278 default→ErrProbePermitInUse
B Complete fail 分支
- L255 nil→ErrProbePermitMissing；L257 CAS claimed→spent 成功→nil
- L264 Unused→ErrProbePermitInUse；L266 default（未知态）→ErrProbePermitConsumed【未知状态】
B persistClaim 原始失败（经 Require 折叠为 ErrProbePermitConsumed）：L286 空路径 `errors.New`；L290 MkdirAll / L304 json.Marshal / L310 OpenFile / L315 Write → `return err`；L317 `file.Sync()`

## 2. internal/gates/probe_permit_test.go
编译期 pin：L26-30（NewProbePermit/LoadProbePermit/Require/Complete 签名）
D 既有测试
- L35 TestProbePermitConcurrentRequireSingleWinner（40 并发 Require 恰一胜者）
- L73 TestProbePermitFailureDoesNotResetQuota（弃用 attempt 保留 claim）
- L95 TestProbePermitTimeoutDoesNotResetQuota（超时/取消不还 quota）
- L123 TestProbePermitRestartDoesNotResetQuota（重启 Load 后仍 Spent）
- L147 TestProbePermitDamagedMarkerStaysSpent（损坏 marker 仍加载为 Spent）
- L189 TestProbePermitSecondInstanceOnSamePathLosesOnly（同路径双实例仅一胜，同进程模拟跨进程 O_EXCL）
- L217 TestProbePermitCompleteOnUnusedRejected（Unused 上 Complete 被拒且不烧 quota）
- L233 TestProbePermitRestartThenCompleteIsIdempotent（重载后 Complete 幂等）
- L263 TestProbePermitSuccessImpliesMarkerOnDisk（nil 返回⇒marker 已落盘）
- L295 TestProbePermitCancelledBeforeClaimDoesNotConsume（预检取消不消耗 quota）
E 覆盖事实：L142/L152/L215/L266 各 fail 分支在本文件无对应 Test 函数；LoadProbePermit 非 NotExist 读失败（L152）无专测；真跨进程（双进程）双写无专测，仅 L189 同进程双实例替代

## 3. internal/chain/seal.go
A 错误哨兵
- L36 `ErrSealUnavailable = errors.New("evidence seal unavailable")`
- L39 `ErrSealIncomplete = errors.New("evidence closure incomplete")`
- L42 `ErrSealMismatch = errors.New("evidence seal does not match the accepted closure")`
- L44 `ErrSealMissing = errors.New("evidence seal missing for an accepted task")`
B 签名：L50 SealRecord struct（L56 CandidateHash、L57 EvidenceRootHash、L58 ContextKey【K 绑定】）；L68 SealSink interface；L74 `(SealRecord).validate()`；L97 `(SealRecord).Hash()`；L117 `(FileEventStore).SealTaskEvidence`；L173 `(FileEventStore).LoadTaskSeal`；L201 sameClosure（非导出）
C 拒绝路径
- validate：L76 身份不完整；L78-79 candidate/root 非哈希；L81-82 ContextKey 非空且非哈希【K 绑定】；L84-85 receipt 非哈希；L87-89 evidence 成员非哈希
- SealTaskEvidence：L119 ctx.Err 透传；L123 validate / L127 Hash 透传；L131 MkdirAll→ErrSealUnavailable；L136 Marshal→ErrSealUnavailable；L142 O_EXCL 冲突→L144 loadErr 透传 / L147 ErrSealMismatch / L149 复用 existingHash；L151 非 Exist 打开失败→ErrSealUnavailable；L158 Write（删半写文件）→ErrSealUnavailable；L163 Sync→ErrSealUnavailable；L167 Close→ErrSealUnavailable
- LoadTaskSeal：L175 ctx.Err 透传；L180-181 NotExist→ErrSealMissing；L183 非 NotExist 读失败→ErrSealUnavailable【非 NotExist 读失败】；L187 Unmarshal→ErrSealMismatch；L190 validate / L194 Hash 透传
E：ContextKey 校验 L81、sameClosure 比对 L211；K 派生不在本文件

## 4. internal/chain/seal_verify.go
B：L24 `(Engine).verifyAcceptedSeals`（非导出，恢复校验 PASSED×seal）；L54 splitTaskKey；L66 containsHash
C：L25-27 未实现 SealSink→`return nil`（跳过校验）；L41 split 失败→`errors.Join(ErrRecoveryUncertain, err)`；L45 LoadTaskSeal 失败→ErrRecoveryUncertain join；L48 hash 不在 accepted 证据→ErrRecoveryUncertain+ErrSealMismatch；L57/L61 malformed key→ErrInvalidState

## 5. internal/chain/state_seal_test.go
辅助：L25 sealFailingEvents、L39 fixedAttempts
D 既有测试
- L72 TestTaskSealBarrierPreventsAccepted（I1：seal 不可持久化→任务停 REPAIR_PENDING，不写 accepted）
- L95 TestClosureIncompleteCannotSeal（I2：闭包内非法哈希=证据缺失）
- L117 TestPassedTaskSealMustVerifyOnRecovery（I3：恢复逐条核对 accepted terminal 的 seal）
- L147 TestReopenUsesResolvedAttempt（I4：resolved attempt 触达全部 port/entity key）
- L193 TestContextKeyStableAcrossReentry（I5：K 仅由不可变输入派生，重入一致、跨 attempt 不复用）【K 绑定】
- L214 TestSealRecordBindsAcceptedClosure（I6：seal 绑定 candidate+双 receipt 哈希）
- L243 TestAttemptResolverFailureBlocksAcceptance（I7：resolver 失败阻断接受）
- L262 TestTerminalRouteRefusesForeignAttempt（I8：终端绑定他人 attempt 被拒）
- L295 TestSealBindsTaskContextKey（I9：用另一 context key 写的 seal 被拒）【K 绑定】

## 6. internal/gates/goal.go
不存在（跳过）。相关符号实际位于 internal/gates/readiness.go：
A 哨兵：L53 ErrReadinessProbeNotAcceptance；L56 ErrReadinessEvidenceIncomplete；L59 ErrReadinessAcceptanceMissing；L62 ErrReadinessProductionAT23Required；L67 readinessHashPattern `^(?:sha256:)?[0-9a-f]{64}$`
B/分支：L76 GoalAcceptance、L78 GoalProduction 常量；L89 SealBinding struct（L90 Attempt、L91 TaskKey、L92 SealHash）；L96 ReleaseBinding（L97 Seal、L98 AT23Hash）；L104 ReadinessInput（L106 Seal *SealBinding、L107 AT23Hash）；L159 EvaluateT027BReadiness
- GoalAcceptance 分支 L173：L177 probe 证据冒充→ErrReadinessProbeNotAcceptance；L179 缺 Seal→ErrReadinessEvidenceIncomplete；L181-182 validate 失败
- GoalProduction 分支 L186：L199 缺 level-2→ErrReadinessAcceptanceMissing；L201-202 validate；L205-206 AT23 空→ErrReadinessProductionAT23Required；L208-209 AT23 不匹配→ErrReadinessEvidenceIncomplete【GAP-70 非空非法 AT23Hash】；L213 AT23Hash 写入 ReleaseBinding；L218 未知 goal
- SealBinding.validate L224：L225-226 Attempt<1；L228-229 TaskKey 空；L233-234 SealHash 非 canonical【GAP-70 非法 SealBinding】
E 引用面：readiness_test.go L47/L52-53/L58/L60/L69-70/L141/L190-191/L280/L291-309/L329/L366/L400-425（TestReadinessAcceptanceRejectsInvalidSealBinding L280、GoalProduction AT23 用例 L366）

GLM-2/2: END

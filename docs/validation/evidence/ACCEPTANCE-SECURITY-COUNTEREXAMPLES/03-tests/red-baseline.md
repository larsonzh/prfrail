# ③ 测试先行 · 防御存在性红基线报告原文（逐字归档）

载体：runSubagent(agentName="prfrail-tester", model="DeepSeek V4.1 Flash (deepseek)")。临时探针文件已由 ③ 删除，未改任何生产代码与既有测试。

PROBE: ALL-GREEN

**HEAD = `941f60e`**；工作目录 `d:\LZProjects\prfrail`；Windows/Go 工具链。

## 任务 A：哨兵机械核验（包 / 文件:行号 / `errors.New` 原文）

**evidence**（`internal/evidence/boundread.go`）
- :31 `ErrBoundArtifactExceedsBound = errors.New("bound artifact exceeds its byte bound")`
- :33 `ErrBoundFileClosed = errors.New("bound file is closed")`
- :35 `ErrBoundFileNotRegular = errors.New("bound artifact is not a regular file")`

**gates · 许可**（`internal/gates/probe_permit.go`）
- :66 `ErrProbePermitMissing = errors.New("probe permit missing")`
- :69 `ErrProbePermitInUse = errors.New("probe permit already in use")`
- :71 `ErrProbePermitConsumed = errors.New("probe permit already consumed")`

**gates · 就绪**（`internal/gates/readiness.go`）
- :53 `ErrReadinessProbeNotAcceptance = errors.New("probe readiness is not acceptance readiness")`
- :56 `ErrReadinessEvidenceIncomplete = errors.New("readiness evidence incomplete")`
- :59 `ErrReadinessAcceptanceMissing = errors.New("acceptance readiness missing")`
- :62 `ErrReadinessProductionAT23Required = errors.New("production readiness requires the AT23 acceptance hash")`

**chain · 封条**（`internal/chain/seal.go`）
- :36 `ErrSealUnavailable = errors.New("evidence seal unavailable")`
- :39 `ErrSealIncomplete = errors.New("evidence closure incomplete")`
- :42 `ErrSealMismatch = errors.New("evidence seal does not match the accepted closure")`
- :44 `ErrSealMissing = errors.New("evidence seal missing for an accepted task")`

**chain · 状态**（`internal/chain/models.go`）
- :14 `ErrInvalidState = errors.New("invalid chain state")`
- :18 `ErrRecoveryUncertain = errors.New("chain recovery uncertain")`
- 旁注：`internal/applier/journal.go:21` 另有一个同名 `ErrRecoveryUncertain`（applier 包域），与 `chain.ErrRecoveryUncertain` 是两个独立哨兵；探针按 chain 包取值。

## 任务 B：既有测试基线

```
go test -count=1 ./internal/evidence/ ./internal/adapters/ ./internal/gates/ ./internal/chain/
ok  github.com/larsonzh/prfrail/internal/evidence  1.271s
ok  github.com/larsonzh/prfrail/internal/adapters  44.550s
ok  github.com/larsonzh/prfrail/internal/gates     3.206s
ok  github.com/larsonzh/prfrail/internal/chain     2.254s
```
四包全 ok，失败测试名：无（无 FAIL 行）。基线在探针删除后执行，不含探针。

## 任务 C：防御存在性探针（5 面 / 6 个测试函数）

运行命令：
```
go test -count=1 -v -run 'TestProbeP52' ./internal/evidence/ ./internal/adapters/ ./internal/gates/ ./internal/chain/
```

**① evidence（`internal/evidence/zz_probe_test.go`，package evidence）**
摘要：写 16B 临时文件 → `OpenBoundFile(path,64)` → `Close()` → `ReadAll()`/`Hash("probe")` 必须各 wrap `ErrBoundFileClosed`；另建 4B 文件 `OpenBoundFile(path,4)`，外部 `os.OpenFile(O_APPEND|O_WRONLY)` 追加 1B 后 `Hash` 必须 wrap `ErrBoundArtifactExceedsBound`（第二写者被拒则 `t.Skipf`）。
结果：`TestProbeP52BoundFileClosed` **PASS**；`TestProbeP52BoundFileGrewPastBound` **PASS（未 SKIP）** —— Windows 侧 `FILE_SHARE_READ|WRITE|DELETE` 允许第二写者，超界增长路径真实可测。
`errors.Is` 成立：是（两处断言均通过）。实际错误串：无（PASS 路径未触发 `t.Fatalf`，未产生逃逸错误）。

**② adapters（`internal/adapters/zz_probe_test.go`，package adapters，无 build tag）**
摘要：对 `""`、`"/etc/passwd"`、`".."`、`"a/../../b"` 四形态调用 `cleanBoundRelativePath`，任一返回 nil error 即 FAIL。
结果：`TestProbeP52CleanBoundRelativePathRejectsEscapes` **PASS**（四种全部被拒：空名、绝对路径、`.`/`..` 清空逃逸、清理后 `../` 前缀逃逸）。
`errors.Is` 不适用（返回普通 `fmt.Errorf`，无哨兵；属预期设计）。

**③ gates · 许可（`internal/gates/zz_probe_test.go`，package gates）**
摘要：`NewProbePermit` 后 `state.Store(3)`（未知态）→ `Complete()` 必须 wrap `ErrProbePermitConsumed`；`LoadProbePermit("")` 必须 wrap `ErrProbePermitMissing`；跨进程：父进程 `exec.Command(os.Args[0], "-test.run=TestProbeP52PermitCrossProcessHelper")` + env 传 marker 路径，子进程 `LoadProbePermit`→`Require` 成功（exit 0），父进程再 `NewProbePermit(同路径).Require` 必须 wrap `ErrProbePermitConsumed`（O_CREATE|O_EXCL 已被占用）。
结果：`TestProbeP52CompleteUnknownStateIsConsumed` **PASS**；`TestProbeP52LoadProbePermitEmptyPath` **PASS**；`TestProbeP52PermitCrossProcessDoubleClaim` **PASS**（0.19s，真实双进程，子进程 exit 0 由父进程校验）。`errors.Is` 成立：是（三处均通过）。

**④ gates · 就绪（`internal/gates/zz_probe_readiness_test.go`，package gates）**
摘要：`EvaluateT027BReadiness(GoalProduction, {Seal: 合法 SealBinding{1,"task:1",sha256:+64hex}, AT23Hash:"sha256:"+63hex})` 必须 wrap `ErrReadinessEvidenceIncomplete`；`EvaluateT027BReadiness(GoalAcceptance, {Seal:{Attempt:1,TaskKey:"t",SealHash:"not-a-hash"}})` 必须 wrap `ErrReadinessEvidenceIncomplete`。
结果：`TestProbeP52ProductionRejectsMalformedAT23` **PASS**；`TestProbeP52AcceptanceRejectsMalformedSealBinding` **PASS**。`errors.Is` 成立：是。注：代码路径分别为 `readiness.go:206-209`（非空非法 AT23 → `%w: invalid AT23 hash`）与 `:181-183`（SealBinding.validate 失败），两者都包 `ErrReadinessEvidenceIncomplete`。

**⑤ chain（`internal/chain/zz_probe_test.go`，package chain）**
摘要：构造 `SealRecord{SchemaVersion:sealSchemaVersion, RunID, TaskID, Attempt:1, ParentHash/CandidateHash/EvidenceRootHash/ReviewReceiptHash/PromotionReceiptHash = "sha256:"+64hex, ContextKey:"not-a-hash", Evidence:[]string{hash}}` → `validate()` 必须 wrap `ErrSealIncomplete`。
结果：`TestProbeP52SealRecordRejectsContextKeyNotHash` **PASS**（命中 `seal.go:82` 的 `context key is not a hash` 分支）。`errors.Is` 成立：是。

**FAIL 探针**：无。无未阻断证据需登记。

## 收尾确认
- 5 个探针文件全部 `Remove-Item` 删除；`Get-ChildItem -Recurse -Filter 'zz_probe*_test.go'` 计数 = **0**。
- `git status --short`（探针删除后）无任何新增未跟踪 `_test.go`；仅步骤 0 既有切片产物（`docs/validation/...`、`docs/t027/REMAINING_SLICES*.md`）与 `01-arch/`、`09-glm/` 目录。
- 未修改任何 `internal/**` 生产代码，未改动既有测试；任务 B 基线在删除后复跑且四包全绿。

TDD-RED-BASELINE: END

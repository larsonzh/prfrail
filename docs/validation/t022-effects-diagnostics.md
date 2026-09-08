# T022 PC-04 副作用与恢复诊断验证报告

[English](t022-effects-diagnostics_EN.md)

日期：2026-09-08。结论：`T022 COMPLETE`。

## 实现范围

1. `internal/gates/effects.go`（副作用分类/观察模型、S1 准入与恢复动作）
2. `internal/gates/effects_test.go`
3. `internal/evidence/diagnostics.go`（recovery-diagnosis 记录模型与校验）
4. `internal/evidence/diagnostics_test.go`
5. `internal/chain/recover_test.go`（chain 恢复计划集成测试）
6. `docs/validation/t022-effects-diagnostics.md` / `t022-effects-diagnostics_EN.md`（本报告）

## AT-19 覆盖

1. 只读/本地可丢弃/外部写分离：`EffectClass` 三分类对齐 `effect-record.schema.json`；分类记录由可信策略绑定 `policyHash` 与 `subjectDefinitionHash`，不接受模型自报。
2. S1 外部写拒绝：`AdmitEffect` 对 `external-write` 与未知类一律 `deny`（fail-close），仅 `read-only`/`local-discardable` 放行；不依据命令名或 dry-run 字样判定。
3. 无法证明限制可执行则拒绝：分类校验强制 `external-write` 只能声明 `reconcile-only`/`compensation-only` 恢复保证，禁止宣称 `none-required`/`discard-workspace`。
4. 未知副作用不盲目重试：观察记录 `outcome=uncertain/failed` 必须携带 `errorEvidence`；`ObservationRecoveryAction` 对非 completed 一律 `reconcile-only`（不重投）。
5. 诊断只读脱敏：`RecoveryDiagnosis` 强制 `mode=read-only-redacted` 且 `mutationsPerformed=false`；`clear` 与 `blocked` 的 oneOf 证据规则按 schema 强制。
6. 恢复计划不删锁/日志：`TestRecoverDoesNotMutateEventLog` 证明崩溃重放后 `Recover` 不删除、不改写任何既有事件；`TestRecoverProducesBlockedDiagnosisPlan` 证明恢复计划锚定最后已接受快照并以 `blocked` + `inspect`/`controlled-stop` 动作呈现。

## 门禁

1. `go test -count=1 ./internal/gates ./internal/evidence ./internal/chain`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过（全部包 ok）。
5. 编辑器诊断：新增/改动 Go 文件均无错误。

## 关键命令证据

```powershell
gofmt -w internal/gates/effects.go internal/gates/effects_test.go internal/evidence/diagnostics.go internal/evidence/diagnostics_test.go internal/chain/recover_test.go
go test -count=1 ./internal/gates ./internal/evidence ./internal/chain
go build ./...
go vet ./...
go test -count=1 ./...
```

## 边界

1. effects/diagnostics 为库级能力，暂无 CLI 入口；外部写 runner、幂等外部 API 与补偿协议按 RFC 属 S2 ADR 范围，S1 明确拒绝。
2. 观察记录已在 gate runner 评估阶段接线：当请求包含 effect scope 时，`internal/gates/runner.go` 会生成并校验 classification/observation，并把 `effectObservationHash` 与 `effectRecoveryAction` 绑定到 hook result。T023/T024 仍负责更广泛的产品级编排与展示入口。
3. 诊断的受管进程证据（`processes`）由调用方以 guard 停机证据填充，诊断自身不启动/停止任何进程。
4. 回滚仅覆盖受控文件；网络发布、邮件、数据库写入、包上传与设备操作不纳入"已回滚"保证。

## 复审记录（2026-09-08）

1. 首轮编译失败修正：`StateEvent` 含不可比较字段，`TestRecoverDoesNotMutateEventLog` 改用 `reflect.DeepEqual` 比较事件快照。
2. 门禁复跑全部通过。
3. 后续编译回归中出现 `runner.assess` 参数增量导致的测试签名漂移（WrongArgCount）；已在 `internal/gates/runner_test.go` 的 `assessForTest` 补传 classification 参数（`nil`）并复跑 `go test -count=1 ./internal/gates` 与全量门禁通过。
4. 为避免副作用分类遗漏导致的默认放行，hook 策略改为必须显式声明 `effectClass`（缺失即 `ErrInvalidHook`，fail-close），并移除 `buildEffectScope` 的隐式默认分类；新增测试 `TestPrepareRejectsMissingEffectClass`、`TestRunnerDeniesExternalWriteEffectWithoutExecuting`、`TestChainPortDeniesExternalWriteHook`。

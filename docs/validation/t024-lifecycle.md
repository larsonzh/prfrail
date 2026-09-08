# T024 PC-06 生命周期验证报告

[English](t024-lifecycle_EN.md)

日期：2026-09-08。结论：`T024 COMPLETE`。

## 实现范围

1. `internal/snapshot/backup.go` / `backup_test.go`：完整对象、事件段、引用根闭包备份；新 store 复读核验；临时目录写入后原子发布。
2. `internal/evidence/disposition.go` / `disposition_test.go`：release/backup/retirement 生命周期记录、严格解码和摘要校验、发布清单、停用与数据保留规则。
3. `docs/DEV_PLAN*.md`、`docs/OPERATIONS*.md`、`docs/INSTALLATION_PLAN*.md`：同步 T024 完成事实与 T017 仍待交付的边界。

## AT-21 覆盖

1. 活动写者与未知格式阻断：备份要求 `WritersStopped=true`、非空停机证据且 `SourceSchemaVersion=1.0.0`；未知版本不迁移、不写目标。
2. 完整闭包：对象、事件段、引用根均须非空、摘要有效且唯一；缺对象、缺事件或缺引用根均在写入前 fail-close。
3. 新 store 恢复核验：对象复制到独立 `Store`，事件与引用按摘要落盘并复读；全部通过后才生成 restore evidence 和 completed manifest。
4. 原件不变与中断安全：测试比较备份前后源对象；故障注入只留下可清理 staging，不发布目标或完成 manifest。
5. 生命周期记录：completed backup 强制 `closureStatus=complete`、secret excluded、对象/根/恢复证据和 manifest hash；failed/uncertain 不能携带完成 manifest。
6. 默认保留：retirement 构造默认 `retentionDisposition=retain`、`sharedToolsDisposition=preserved`，禁止删除 SessionBridge 或用户工具链。
7. 删除与审计冲突：`delete-authorized` 必须引用人工删除授权；存在保留冲突时必须引用独立人工决定。
8. 发布清单与信任表达：清单按路径排序、内容摘要稳定并拒绝符号链接；release 强制 SBOM/许可证/校验和/支持矩阵。2026-09-08 契约修订后签名 receipt 为可选；无签名时信任/撤销新鲜度只能为 `unknown`，有签名时仍执行严格摘要和验证证据约束。

## 门禁结果

1. `go test -count=1 ./internal/snapshot ./internal/evidence`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过（全部包 ok）。

## 测试清单（8 项）

1. `internal/snapshot/backup_test.go`: `TestCreateBackupVerifiesNewStoreAndPreservesSource`
2. `internal/snapshot/backup_test.go`: `TestCreateBackupRejectsActiveWriterUnknownSchemaAndMissingClosure`
3. `internal/snapshot/backup_test.go`: `TestCreateBackupInterruptedHasNoCompletedManifest`
4. `internal/evidence/disposition_test.go`: `TestLifecycleBackupRoundTrip`
5. `internal/evidence/disposition_test.go`: `TestRetirementDefaultsRetainAndPreserveSharedTools`
6. `internal/evidence/disposition_test.go`: `TestRetirementDeletionAndConflictRequireHumanEvidence`
7. `internal/evidence/disposition_test.go`: `TestReleaseAllowsUnsignedS1AndRepresentsTrustAsUnknown`
8. `internal/evidence/disposition_test.go`: `TestBuildReleaseInventoryIsDeterministicAndRejectsLinks`

## 边界

1. T024 提供库级能力，不提供安装包、发布动作或可复制的 backup/uninstall CLI；这些仍受 T017 和发行授权约束。签名/attestation 为后续可选增强。
2. unknown schema 仅拒绝写入并保留原 store；显式跨版本迁移属于 S2。
3. 逻辑删除不保证 SSD、快照或外部副本安全擦除。
4. 本轮未执行 `git commit`、`git push`、签名或发布。

## 复审记录

1. 恢复上一轮未完成的 `disposition.go`，修正包内 slice clone 编译缺口后保留其契约实现。
2. VS Code 测试适配器未发现新 Go 测试文件，改用仓库原生 `go test` 执行并通过。
3. 相关中英文文档已同步，未修改用户在本轮前调整过的 T023 双语验证报告。
# T020 PC-02 导出验证报告

[English](t020-export_EN.md)

日期：2026-09-08。结论：T020 COMPLETE。

## 实现范围

1. internal/snapshot/export.go
2. internal/snapshot/export_test.go
3. internal/evidence/delivery.go
4. internal/evidence/delivery_test.go

## AT-17 覆盖

1. 已接受快照绑定：导出要求 acceptedSnapshotHash 与 snapshot manifest hash 一致，防止候选冒充接受。
2. 目标目录前置条件：destination 必须 absent，已存在目录或文件直接阻断。
3. 路径重叠阻断：destination 与 source/run/store 任一重叠时 fail-close。
4. 离线完整性：导出前后验证对象闭包与导出文件内容哈希，缺对象与篡改均阻断。
5. secret 阻断：secret/token/password/credential/.env 等路径命中即拒绝导出。
6. 写中断无完成回执：注入中断后不生成 completed 摘要，目标目录保持 absent。
7. 源树不变：导出流程只读 source/store，测试断言源文件内容保持不变。

## 关键测试

1. TestExportAcceptedPackageRoundTrip
2. TestExportAcceptedPackageRejectsHashMismatchAndExistingDestination
3. TestExportAcceptedPackageRejectsOverlapSecretMissingAndTamper
4. TestExportAcceptedPackageInterruptedWriteHasNoCompletion
5. TestExportRecordCompletedRoundTrip
6. TestExportRecordCompletedRequiresEntriesAndManifest
7. TestExportRecordFailedRequiresErrorEvidence
8. TestExportRecordRejectsDestinationPreconditionAndHashMismatch

## 门禁

1. go test ./internal/evidence ./internal/snapshot：通过。
2. go build ./...：通过。
3. go vet ./...：通过。
4. go test -count=1 ./...：通过。

## 边界

1. 当前仅提供库级导出接口，CLI export 命令在后续切片。
2. secret 路径检测采用 fail-close 关键字策略，命中即阻断。
3. 导出完成不改变 chain/task 接受状态，也不代表外部发布完成。
4. secret 关键字策略可能误报（如 `tokenizer.txt`）或漏报（如 `.npmrc`/`.pem`）；方向为 fail-close，后续与 secret-policy exclusion 机制合并。
5. `destinationIdentityHash` 仅绑定导出条目内容，不绑定目标绝对路径；CLI 接入时应明确该语义。
6. 导出目录统一 0755、硬链接按独立文件写出；需要位级还原时使用 Restore 而非导出。
7. 尚未覆盖 Windows 短路径/junction 别名与 secret-policy 排除项的专项测试。

## 复审记录（2026-09-08）

1. 删除未使用的 `NormalizeExportHashes`（死代码清理）。
2. 门禁复跑全部通过。

# T008 托管变更集与事务应用验证报告

[English](t008-applier_EN.md)

日期：2026-09-07。结论：`T008 COMPLETE`。本任务完成 RFC §16.9.25 托管变更集 checker、带 fencing 的 journal 事务、整组回滚及 AT-04 崩溃恢复验收。

## 十二个切片

1. 冻结 change-set、operation、assertion、marker 和纯内存写入计划模型。
2. 严格解析 JSON，复算 JCS/domain-separated changeSetHash 并绑定 task attempt、父 snapshot 和前后 manifest。
3. 校验受管相对路径、普通文件目标、连续唯一 sequence/operationId 及 1–10000 操作边界。
4. 一次读取全部目标后仅在内存按数组顺序模拟，同一路径多操作以前一模拟结果为下一 beforeHash。
5. 实现 create-file/delete-file/replace-exact/insert-before/insert-after 五种封闭操作。
6. 强制 UTF-8、无 NUL、内部 LF、显式 create 行尾和已有文本 preserve-existing，不猜测转换。
7. 强制 before/after assertion、exact match 数量及 insert marker 的唯一 0→1 契约；任一失败零写入。
8. 在工作区写入前持久化 changeSetHash、前后字节摘要、模式和新建目录计划；blob 按摘要重读校验。
9. 每个目标经同目录临时文件、文件 fsync 和平台原子替换；Unix 同步父目录，Windows 使用 MoveFileExW REPLACE_EXISTING|WRITE_THROUGH。
10. 每个 journal 与工作区写点重新验证 T007 fencing token；外部改动、链接父目录和陈旧 writer 均 fail closed。
11. 普通故障立即反向整组回滚；模拟崩溃保留 journal，重启仅在当前内容等于已知 before/after 时恢复，未知内容进入 uncertain 且不覆盖。
12. 覆盖 journal durable、每项写前/写后、全量后验后的崩溃点，并在 Windows 11/NTFS 与 Ubuntu 24.04/ext4 原生验证。

## 验证

- Windows：`go build ./...`、`go vet ./...`、`go test -count=1 ./...` 通过。
- Checker：五种操作、同路径顺序链、严格 JSON、摘要/sequence/managed path/UTF-8/NUL/assertion/marker 反例及失败零写入通过。
- Applier：成功提交、普通故障整组回滚、journal durable、每项 before-write/after-write、after-verify、幂等恢复、外部未知内容、损坏 journal/blob、错误 token、中途 fencing 失效及符号链接父目录反例通过。
- Ubuntu：Go 1.27.1，在 Ubuntu 24.04/ext4 以临时 vendor 离线执行 `go build ./...`、`go vet ./...`、`go test -count=1 ./...` 和 `go test -race -count=1 ./internal/taskdef ./internal/applier`，全部通过。
- `git diff --check` 通过；临时归档、vendor 和远端验证目录已清理，仓库 `tmp/` 仅保留 `.gitkeep`。

## 边界

T008 不调度 task/step、不生成 chain projection/receipt、不执行 gate，也不把候选结果升级为 accepted snapshot；这些分别由 T009 及后续任务负责。Windows 的 `MoveFileExW(...WRITE_THROUGH)` 提供平台可用的替换持久性，但不宣称具备 Unix 目录 fsync 的同等断电保证；任何恢复不确定均保持 fail closed。

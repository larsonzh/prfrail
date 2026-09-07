# T009 Chain Engine 验证报告

[English](t009-chain_EN.md)

日期：2026-09-07。结论：`T009 COMPLETE`。本任务完成 Chain Engine 的顺序调度、事件投影、暂停/取消和崩溃重放，并以 fake agent/runner 完成 AT-03 的 chain 层验收。

## 十四个切片

1. 冻结 chain/task/step、变更模式、工作区、快照及 consumer-owned port 模型。
2. 复用 evidence Schema 的 chain/task/step 转换表，不建立替代状态协议。
3. 在写入前拒绝未列出的转换、自转换和非法起始状态。
4. 以 JSONL append-only store 持久化 state event，并拒绝损坏或断尾记录。
5. 从完整事件链重建 sequence、event hash、状态及 accepted parent projection。
6. 先持久化 `NONE→CREATED`，baseline-0 证据持久化后才进入 `BASELINED`。
7. 按定义顺序调度三个 task，并只以前一 accepted snapshot 作为下一父快照。
8. 覆盖 code/build/verify/noop 四 kind；noop 只记录原因且不调用 runner。
9. 通过显式端口路由 managed-change-set code step。
10. 通过独立显式端口路由 isolated-workspace code step。
11. pause 在原子边界停止新 step，普通暂停可从投影位置恢复并完成。
12. cancel 先调用受管进程树 stopper、归档证据，再写不可恢复的 `CANCELLED`。
13. 恢复先运行 reconciler，再验证并重放事件；运行中 step 进入 recovery-uncertain `PAUSED`，禁止重复投递或直接 resume。
14. 在 Windows 11/NTFS 与 Ubuntu 24.04/ext4 原生执行 AT-03 验证。

## 验证

- Windows：`go build ./...`、`go vet ./...`、`go test -count=1 ./...` 通过。
- Chain：三任务顺序、四 kind、两种变更模式、accepted parent 推进、T2 失败阻断 T3、事件先于投影、非法转换零写入、损坏尾部拒绝均通过。
- 控制：pause/resume、cancel/stopper、noop 零 runner 调用及不可恢复取消均通过。
- 恢复：accepted parent 重建、reconciler 调用、运行中 step 不重投、recovery-uncertain 持久暂停及直接 resume 拒绝均通过。
- Ubuntu：Go 1.27.1，在 Ubuntu 24.04/ext4 以临时 vendor 离线执行 `go build ./...`、`go vet ./...`、`go test -count=1 ./...` 和 `go test -race -count=1 ./internal/chain`，全部通过。
- Windows 未执行 `-race`，因为本机 `CGO_ENABLED=0` 且无 C 编译器；竞态检测已在 Ubuntu 原生通过。
- 临时归档、vendor 和远端验证目录已清理。

## 边界

T009 使用 fake agent/runner，不实现 T010 的真实 gate runner、T011 的 review/promotion 或 T013/T014 的真实 agent 与人工交接。`AcceptancePort` 只表达“后续 task 只能消费已接受父快照”的边界；候选、失败和不确定结果不会成为下游父快照。恢复不确定时保持 fail closed，需由后续修复流程解决，不能猜测成功或重复投递。
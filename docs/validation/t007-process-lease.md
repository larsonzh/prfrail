# T007 进程与租约守卫验证报告

[English](t007-process-lease_EN.md)

日期：2026-09-07。结论：`T007 COMPLETE`。本任务完成 guard 层对 AT-03/AT-06 的进程树停机与单写租约贡献；任务链状态调度和 hook 策略执行分别由 T009、T010 验收。

## 十个切片

1. 冻结受管进程、身份、终止证据和 fencing token 模型。
2. Linux 以独立 process group 启动并按 PGID 终止整树。
3. Windows 以 `CREATE_SUSPENDED` 启动，绑定 `KILL_ON_JOB_CLOSE` Job Object 后恢复，消除绑定前逃逸窗口。
4. Linux 以 `/proc/<pid>/stat` starttime、Windows 以 process creation time 绑定 PID 身份，防 PID 复用误判。
5. 终止先尝试平台正常停止边界，再强制整树终止并验证树和原身份均消失。
6. 生成 JCS canonical、域分离 SHA-256 的新鲜结构化终止证据；不确定时返回 `ErrTerminationUncertain`。
7. 文件租约使用互斥发布，初始 generation=1，同一资源并发领取仅一个成功。
8. 续租保持 claimId、递增 generation；旧 `(claimId,generation)` 无法验证或释放当前租约。
9. 接管严格要求不同 claimId、更高 generation、与旧 writer 身份一致且五分钟内的新鲜 stopped 证据、至少一项授权引用，并由注入的独立 `TakeoverVerifier` 验证券据闭包；未配置验证器时默认拒绝，租约过期、PID 不存在或重启均不能单独授权。
10. Windows 11/NTFS 与 Ubuntu 24.04/ext4 原生父子树、并发租约和接管反例验证，补齐双语文档。

## 验证

- Windows：`go test -count=20 ./internal/guard` 与 `go vet ./internal/guard` 通过；每轮真实启动测试父进程和后代，Job Object 终止后按启动身份确认两者离线。
- Ubuntu：交叉编译后在 Ubuntu 24.04/ext4 原生执行全部 guard 测试，通过 Linux process-group 整树终止。
- T006 联合回归：`go test ./internal/snapshot` 与 `go vet ./internal/snapshot` 通过。
- 本机 Go 未启用 cgo，`go test -race` 不可用；以 16 个并发领取者及 20 轮重复测试覆盖竞争路径，race detector 留作具备 cgo 的 CI 补充项。

## 边界

T007 不生成 adapter claim/takeover receipt，不调度 task/step，不决定 argv、cwd、网络或资源策略，也不实现人工交接。调用方必须在所有工作区写入、恢复和发布前验证当前 fencing token，并持久化本层证据引用。
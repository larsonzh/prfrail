# T010 Gate Runner 验证报告

[English](t010-gates_EN.md)

日期：2026-09-07。结论：`T010 COMPLETE`。本任务完成 hook 策略预检、runner 编排、结构化结果、证据持久化及 Chain Engine 接线，并完成 AT-06 验收。

## 十二个切片

1. 冻结 hook、process runner、资源/网络策略、execution、result 与 capability 模型。
2. 校验 hook kind、failure policy、timeout、正资源额度、ID 和定义摘要。
3. executable 与 args 保持独立 argv，不经隐式 shell；绝对/反斜杠 executable 拒绝。
4. cwd 必须解析到授权工作区内，`..` 与 symlink 越界均拒绝。
5. 环境仅从显式 allowlist 选择；缺失、重复或非法变量名 fail closed。
6. 启动前检查 process、timeout、output、memory、process-count、CPU 与 network capability；不可执行的策略不降级。
7. 复用 T007 guard 管理进程树，覆盖正常退出、启动失败、超时停止与 termination-uncertain。
8. stdout/stderr 使用共享有界捕获；超限停止进程树并记录 resource-limited/truncated。
9. stdout/stderr 与声明产物先执行秘密扫描，再写入 snapshot CAS 并绑定摘要；扫描失败或缺产物阻断。
10. 分离 executionOutcome、assessment、failureKind 与 policyDisposition；warn 保持 failed，retry 逐次持久化且受 maxAttempts 限制。
11. result 以 `proofrail:hook-result:1` 计算 JCS 摘要，JSONL fsync 后才返回；篡改、损坏日志和重复执行键拒绝。
12. Chain HookPort 绑定 step/hook kind 与定义摘要，成功和失败均把 durable resultHash 写入 step 证据。

## 验证

- Windows：`go build ./...`、`go vet ./...`、`go test -count=1 ./...` 通过。
- 契约：`npm test --prefix tools/contracts` 通过；29 个 Draft 2020-12 Schema、全部正反 catalog 与 JCS 向量无漂移。
- AT-06：覆盖 argv 注入、cwd/symlink 越界、环境秘密排除、能力不可用零启动、timeout、无限输出、秘密扫描失败、缺产物、非零退出、retry、结果篡改/重复及持久化失败。
- Ubuntu：Go 1.27.1，在 Ubuntu 24.04/ext4 以临时 vendor 离线执行 `go build ./...`、`go vet ./...`、`go test -count=1 ./...` 和 `go test -race -count=1 ./internal/guard ./internal/gates`，全部通过。
- Windows 未执行 `-race`，因为本机 `CGO_ENABLED=0` 且无 C 编译器；竞态检测已在 Ubuntu 原生通过。
- 临时归档、vendor 和远端验证目录已清理。

## 边界

内置 `GuardExecutor` 真实提供独立 argv、进程树、timeout 和共享输出上限，但不声称提供网络隔离、memory、process-count 或 CPU quota。hook Schema 要求这些策略，因此直接使用该受限 executor 会在 preflight 阻断；生产接入必须提供声明并实际实现完整能力的 Executor，禁止降级运行。container/remote/shell/PTY、统一跨任务 retry budget、T011 review/promotion 与 T022 effect policy 不属于 T010。
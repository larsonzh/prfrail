# 作废记录 — ②③ 误用 DeepSeek V4 Pro 派单

- 时间：2026-10-31
- 事实：②③ 实现/测试派单误用 `DeepSeek V4 Pro (deepseek)`，违反 §2.2 角色-模型绑定（②③⑧ 应均为 `DeepSeek V4.1 Flash (deepseek)`）
- 处置：产物全部作废，文件删除后由 `DeepSeek V4.1 Flash (deepseek)` 重做（测试先行）。

| 文件 | 字节 | sha256 |
|---|---|---|
| `internal/tickets/durable_cost_ledger.go` | 17913 | 2e9ed784d73120228622db8cfdb01c8197b1fa59707a93c94f9b8c070b8b1dfc |
| `internal/tickets/durable_cost_ledger_test.go` | 17927 | 5071713d51eb8d1f4f87c9c723fd6d2963de346650634c719b3be3b69f7672be |
| `internal/tickets/durable_cost_lock_unix.go` | 1388 | 2b3aba61ed3db67c82814a1773ba9102fdcc4ce8551bb2e05185301cc31ccfbe |
| `internal/tickets/durable_cost_lock_windows.go` | 3342 | f6426688fea760aefbb85aeed8ce3fe30d63df64641b88c9a9140c112ecbfd65 |

# 执行自由度清单 — P39-T027B-DURABLE-BUDGET

切片：`P39-T027B-DURABLE-BUDGET`（P39，实现类 `[IMPL]`）。准则 v1.59。

| 项 | 内容 |
|---|---|
| 申请自由度 | 无 |
| 说明 | 本片在用户裁定的边界内执行：只改 `internal/tickets/`、`internal/adapters/`、`internal/console/` 与对应测试；不改 `docs/CONTRACTS.md` 第七节、不改判据本体、不改 schema 与 fixtures |
| 落盘时机 | 步骤 0（⑥ 之前） |

# 第 1 轮输入包清单（代码级施工图）

切片 `P38-T027B-TACTICAL-BLUEPRINT`（评估类 `[SLICE]`）。本清单由主控实测生成（行数 = `LF` 记录数，字节 = 磁盘原始字节，含各自 BOM）。

| 路径 | 行数 | 字节 | 用途 |
|---|---|---|---|
| `docs/validation/evidence/P37-T027B-BREAKTHROUGH/01-sol-design/review.md` | 111 | 10723 | P37 第 1 轮《T027-B 生产级架构设计书》（Sol，主控逐字转录） |
| `docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/review.md` | 68 | 9277 | P37 第 2 轮《T027-B 设计批注与重构方案》（Astra，含 C1–C7） |
| `docs/CONTRACTS.md`（:188–:283） | 96 | 48064 | §7 AgentRunner 与辅助消息通道边界（含 §7.1） |
| `internal/adapters/agent_runner.go` | 182 | 6813 | Go 源（非测试） |
| `internal/adapters/agent_runner_admission.go` | 217 | 8471 | Go 源（非测试） |
| `internal/adapters/agent_runner_capability.go` | 200 | 8038 | Go 源（非测试） |
| `internal/adapters/agent_runner_completion.go` | 183 | 7955 | Go 源（非测试） |
| `internal/adapters/agent_runner_enforcement.go` | 217 | 9800 | Go 源（非测试） |
| `internal/adapters/agent_runner_event.go` | 174 | 6686 | Go 源（非测试） |
| `internal/adapters/agent_runner_evidence.go` | 463 | 19988 | Go 源（非测试） |
| `internal/adapters/agent_runner_intent.go` | 72 | 2778 | Go 源（非测试） |
| `internal/adapters/agent_runner_launch_receipt.go` | 195 | 8468 | Go 源（非测试） |
| `internal/adapters/agent_runner_launch_reconfirmer.go` | 57 | 2413 | Go 源（非测试） |
| `internal/adapters/agent_runner_pinned_cli.go` | 606 | 27617 | Go 源（非测试） |
| `internal/adapters/agent_runner_postflight.go` | 742 | 36299 | Go 源（非测试） |
| `internal/adapters/agent_runner_process_registry.go` | 260 | 9415 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_dispatcher.go` | 268 | 12104 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_root.go` | 565 | 21545 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_root_unix.go` | 23 | 515 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_root_windows.go` | 28 | 765 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_store.go` | 668 | 26339 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_store_launches.go` | 320 | 13354 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_store_terminals.go` | 375 | 14959 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_store_unix.go` | 18 | 322 | Go 源（非测试） |
| `internal/adapters/agent_runner_replay_store_windows.go` | 11 | 201 | Go 源（非测试） |
| `internal/adapters/agent_runner_run.go` | 399 | 17132 | Go 源（非测试） |
| `internal/adapters/agent_runner_terminal_chain.go` | 146 | 6626 | Go 源（非测试） |
| `internal/adapters/agent_runner_terminal_publisher.go` | 449 | 20291 | Go 源（非测试） |
| `internal/adapters/agent_runner_terminal_record.go` | 309 | 14908 | Go 源（非测试） |
| `internal/adapters/agent_runner_timeout.go` | 116 | 4755 | Go 源（非测试） |
| `internal/adapters/ai_availability.go` | 227 | 7913 | Go 源（非测试） |
| `internal/adapters/ai_preflight.go` | 116 | 4062 | Go 源（非测试） |
| `internal/adapters/ai_probe_copilot.go` | 221 | 8529 | Go 源（非测试） |
| `internal/adapters/ai_probe_copilot_executor.go` | 349 | 10816 | Go 源（非测试） |
| `internal/adapters/ai_provider.go` | 121 | 3876 | Go 源（非测试） |
| `internal/adapters/ai_secret_unsupported.go` | 32 | 720 | Go 源（非测试） |
| `internal/adapters/ai_secret_windows.go` | 164 | 4890 | Go 源（非测试） |
| `internal/adapters/cost_usage.go` | 63 | 2413 | Go 源（非测试） |
| `internal/adapters/filequeue.go` | 885 | 32211 | Go 源（非测试） |
| `internal/adapters/sessbridge.go` | 166 | 4734 | Go 源（非测试） |
| `internal/adapters/sessbridge_visible.go` | 180 | 7198 | Go 源（非测试） |
| `internal/chain/agent_runner_launcher.go` | 45 | 1711 | Go 源（非测试） |
| `internal/chain/agent_runner_terminal.go` | 233 | 10164 | Go 源（非测试） |
| `internal/chain/agent_runner_terminal_route.go` | 225 | 10828 | Go 源（非测试） |
| `internal/chain/authorization.go` | 574 | 21121 | Go 源（非测试） |
| `internal/chain/engine.go` | 515 | 19518 | Go 源（非测试） |
| `internal/chain/events.go` | 87 | 2107 | Go 源（非测试） |
| `internal/chain/handoff.go` | 275 | 9416 | Go 源（非测试） |
| `internal/chain/models.go` | 278 | 7809 | Go 源（非测试） |
| `internal/chain/operator_interaction.go` | 559 | 24834 | Go 源（非测试） |
| `internal/chain/postflight.go` | 459 | 18948 | Go 源（非测试） |
| `internal/chain/publish.go` | 108 | 4607 | Go 源（非测试） |
| `internal/chain/recover.go` | 39 | 1222 | Go 源（非测试） |
| `internal/chain/review.go` | 179 | 6685 | Go 源（非测试） |
| `internal/chain/router.go` | 125 | 4931 | Go 源（非测试） |
| `internal/chain/state.go` | 132 | 4326 | Go 源（非测试） |
| `internal/chain/supervised_blackbox.go` | 319 | 15870 | Go 源（非测试） |
| `internal/chain/supervised_report.go` | 212 | 8195 | Go 源（非测试） |
| `internal/chain/supervised_risk_confirmation.go` | 288 | 12432 | Go 源（非测试） |

## 合计

- 条目 = **59**（3 件文档段 + 56 个 `.go` 非测试文件）
- 行数合计 = **14714**
- 字节合计 = **648707**（≈ 633.5 KB）

## 口径说明

- 代码段仅取 `internal/adapters` 与 `internal/chain` 下**非 `_test.go`** 的 `.go`：分别 38 件 / 18 件。
- 文档段固定整文件（两份 `review.md`）与 `docs/CONTRACTS.md` §7 段落（`## 7.` 起至 `## 8.` 前，含 §7.1）。
- 不含 `_test.go`、不含其余文档、不含任何生成产物。

## 本轮输出分片（分批执行记录）

- 批次 A：`01-sol-blueprint/review.md`（C1–C4）。源字符数 10,052，源字节 13,442（无 BOM）。
- 批次 B：`01-sol-blueprint/review-batch-b.md`（C5–C7）。源字符数 9,808，源字节 14,583（无 BOM）。
- 两批合计覆盖 C1–C7 全集，且批次 B 未重复输出 C1–C4。

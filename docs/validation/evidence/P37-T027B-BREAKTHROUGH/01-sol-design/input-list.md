# 第 1 轮输入包清单（架构设计）

| 路径 | 行数 | 字节 | 用途 |
|---|---|---|---|
| `docs/validation/strategic-review-2026-10-30.md` | 86 | 5606 | P35 战略回看报告（双断点事实源） |
| `docs/validation/strategic-planning-2026-10-31.md` | 94 | 8860 | P36 前瞻规划报告（11 月路线与裁决输入） |
| `docs/CONTRACTS.md`（:188–:283） | 96 | 48064 | §7 AgentRunner 与辅助消息通道边界（含 §7.1） |
| `docs/ARCHITECTURE.md` | 96 | 11871 | 架构分层与平台边界 |
| `docs/SECURITY.md` | 65 | 8570 | 安全前提与边界 |
| `docs/DEV_PLAN.md`（:108–:154） | 47 | 49048 | P1.8 段（T027 与 T027-B 待评估项） |
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

**合计**：62 件（文档段 6 件、Go 源 56 件）、行数 14923、字节 712662。

**文档段合计**：6 件、行数 484、字节 132019。
**Go 源合计**：56 件（`internal/adapters` 38 件、`internal/chain` 18 件）、行数 14439、字节 580643。

**模型**：`GPT-6.1 Sol (copilot)`；档位由用户侧设定为 1M / Max（主控不可观测，见报告 §7）。

**口径**：全文件按原始字节（含 BOM）计；行区间按 `起始行..结束行`（含两端）切片后按 LF 重建计字节，不含原文件 BOM。行数按 `Get-Content` 记录数。

**截断风险提示**：合计约 696 KB；若实际档位为 272K，需按主题分片派发，不得静默截断。

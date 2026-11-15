# P74 SANDBOX-INTERFACE 输入清单（input-list）

切片：`SANDBOX-INTERFACE`（P74 · 层 1 主线 B 段引入独立片 · 实现类 [SLICE] · M 档）
步骤 0 落盘：2026-11-15（主控 DeepSeek V4.1 Flash）

| # | 输入 | 形态 | 状态 |
|---|---|---|---|
| 1 | 用户起跑消息原文（含 2026-11-14 同轮裁定授权） | `docs/validation/evidence/P74-SANDBOX-INTERFACE/00-briefing.md` | 已落盘（§13.1.10） |
| 2 | `docs/t027/CHECKLIST_v57{,_EN}.md` | repo 文件 | 已读（八锚点核对，HEAD 期望 `f834ce5`） |
| 3 | `docs/DELIVERY_DIRECTIVE{,_EN}.md`（v1.69） | repo 文件 | 已读（§1.7 / §2.2 / §5.1 / §5.2 / §7.2.1 / §7.2.2 / §7.5 / §7.7 / §11.1 / §13） |
| 4 | `docs/CONTRACTS{,_EN}.md` §7 | repo 文件 | ②a 修改对象（新增 Sandbox 段） |
| 5 | `internal/guard/**` | 代码 | 现有接口形态（批量辅助 ①） |
| 6 | `internal/chain/**`、`internal/adapters/agent_runner_pinned_cli.go` | 代码 | 只读核对（不改判定本体） |
| 7 | `internal/console/runtime_agent.go`、`runtime_supervised.go`、`supervised_cmd.go` | 代码 | 装配检查点接线对象 |
| 8 | `docs/validation/t027b-e2e-wiring{,_EN}.md` §2 | repo 文件 | P69 permit 时序来源 |
| 9 | `docs/validation/t027b-admission{,_EN}.md` | repo 文件 | P73 裁定与三条附条件 |
| 10 | `docs/ADR_REGISTER{,_EN}.md`（ADR-044） | repo 文件 | 死代码禁令处置依据 |
| 11 | `docs/t027/REMAINING_SLICES{,_EN}.md` | repo 文件 | 步骤 0 新增 P74 段 + ⑧ 收口 |
| 12 | 门禁基线（SELFTEST 737/737 + 三作用域 0） | 会话执行 | 已取 |

R/D 基线：**R = 0、D = 0**（§13.1.1 步骤 0 登记；不继承 P73）。

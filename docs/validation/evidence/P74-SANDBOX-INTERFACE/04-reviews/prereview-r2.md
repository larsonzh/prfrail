# ⑤ 预审 r2（整改后复跑）原文（P74）

- 载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- 调用：2（额度 1 + 1，已用尽）
- 整改内容：r1 三条 Low 全部采纳（Reason 入错误链 / 消除自包裹 / 新增 windows 契约测试）
- 结论：`PRE-REVIEW-R2: PASS`

---

PRE-REVIEW-R2: PASS

## Critical / High / Medium

无发现

## Low

- `internal/guard/sandbox.go:107`（err 分支 case 3）：当实现返回「包裹 `ErrSandboxUnavailable` 的复合错误」且 reason 非空时，新错误只 `%w` 包裹裸 sentinel + reason，原始 err 对象不入链 ⇒ 内层文本丢失。sentinel 仍可被 `errors.Is` 命中，契约 (c) 满足；仅建议未来改为 `fmt.Errorf("%w: %s: %v", ErrSandboxUnavailable, reason, err)` 以保留完整链（S1 两个平台实现均返回裸 sentinel，当前无实际损失）。

## Info

- `internal/guard/sandbox.go:111`：`ErrSandboxRejected` 包裹 `verdict.Reason` 未 TrimSpace，若未来实现给空白 reason 会尾随 `": "`；当前两平台实现与测试 fake 的 reason 均非空，不影响。
- `internal/console/sandbox_preflight.go:15-16`：`sandboxEnforcementEnvVar` 与 `cli.go:308` 读取点分居两文件，但 `grep` 确认产品代码仅此一处读取 ⇒ 仍是「唯一触发源」，契约 (e) 成立。

## 三项 Low 修复确认

1. Reason 入审计：`sandbox.go:96-110` err 分支四 case 均把 `verdict.Reason`（TrimSpace 后）纳入错误文本 ✓
2. 自包裹消除：`sandbox.go:101-103` `reason=="" && errors.Is(err, ErrSandboxUnavailable)` 时 `return err` 原样返回；其余分支文本不重复 sentinel ✓
3. `sandbox_windows_test.go`（`//go:build windows`，31 行）钉死 α：恒 Unknown、非空 reason、nil err、`CheckSandbox` 拒绝 ✓

## 错误链与放行语义（任务 B）

- 四分支全部 `%w` 包裹 sentinel 或直接返回 sentinel ⇒ `errors.Is(err, ErrSandboxUnavailable)` 恒命中；green 测试以 `errors.Is` 断言实证（`green-guard-contract.txt` unavailable case PASS）✓
- 仅 Allowed 放行：`sandbox.go:110-112` err 先于 verdict 判定（err≠nil 即使 verdict=Allowed 也拒绝）；零值 `Verdict{}`=Unknown ⇒ fail-closed ✓
- 调用点双层 `%w`（`ErrInvalidConfig` + sandbox sentinel，Go 1.20+）穿透正常，green 两件实证 ✓

## 契约逐条复核（任务 C）

- (a) `sandbox.go:87-89` 仅 `Check`；无生命周期方法；不引用 guard 三类型 ✓
- (b)(c) 三态 + 双 sentinel + reason 入审计 ✓
- (d) `runtime_agent.go:206-211` 在 `requireProbePermit`(:216) 前；`runtime_supervised.go:318-322` 在 freshStart 块内 permit(:323) 前，均早于 Step 2 marker；拒绝返回零 summary；测试断言 3 以「拒绝后配额可再认领」+ 断言 4 marker/ledger `ErrNotExist` 实证不烧配额 ✓
- (e) `cli.go:306-310` 在 flag 解析后、`getwd`/`LoadChainConfig` 等一切文件读取前；非法值 usage error 零副作用；默认 off ⇒ `SandboxEnforced` 零值跳过，零行为变化；enforce 注入两条路由（`cli.go:417/421`、`supervised_cmd.go:43`）；再入不检查 ✓
- (f) `sandbox_unsupported.go` 恒 Unknown+ErrSandboxUnavailable；`sandbox_windows.go` 恒 Unknown（诚实陈述、非空 reason、不谎报 Denied）✓
- (g)(h) 未动 chain/adapters/schema/CI；文档声明 S-G1/S-G8 未解锁 ✓
- 可证伪：red 3 件（未接线/未实现时断言失败）→ green 3 件全 PASS → mutation 2 件（删调用点后复现失败，supervised 变异走 visible 错误路径同样变红）✓

无越界。

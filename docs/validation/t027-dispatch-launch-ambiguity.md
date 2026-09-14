# T027 A2 dispatch and launch ambiguity 前置验证报告

[English](t027-dispatch-launch-ambiguity_EN.md)

日期：2026-09-14。结论：`PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`。

## 验证范围

本报告只覆盖 T027 离线切片 A2 dispatch and launch ambiguity，不宣称 AgentRunner 生产实现。A2 已交付：

- 契约冻结（先于实现）：`docs/CONTRACTS.md` / `CONTRACTS_EN.md` 新增四段——① dispatch 先发布 R，发布后、启动前窄重确认，失败 fail-closed（保留 R、无回执、不 spawn、不自动重试，错误 `errors.Is` 可判且保留底层原因）；② 二文件 store-local 回执协议（spawn 前 `launches/intent.<requestId>.jsonl` 含 launchId、requestId+requestHash、run/task/step/attempt、adapterId、重确认时刻、授权与预算摘要，no-replace + 有界重读收敛、unproven 平台 pre-write 拒绝；spawn 后 `launches/identity.<launchId>.jsonl` 绑定意图摘要与不透明进程身份；不变量：无意图⇒绝不 spawn、身份缺失⇒unproven、两者都不允许盲目重启）；③ dispatch 重放处理（R-only 不 launch、已有意图不重发不重启、意图+身份不重启、R+C 终态绝不 launch；唯一允许的同进程重试=端口契约性未 spawn 失败，意图绝不重发；接管/结算留给后续切片）；④ `launches/` 与 requests/completions 同等的路径安全/归属/同步纪律。
- `AgentRunnerReplayDispatcher`（adapters，实现 `chain.AgentRunnerPort`）：校验请求绑定 → `RecordRequest`（terminal/unknown 分哨兵阻断）→ 授权账本快照 + `AgentRunnerLaunchReconfirmer` 窄重确认 → 确定性 launchId → `RecordLaunchReceipt`（单赢家）→ `Launcher.StartAgentRunnerProcess` → `RecordLaunchIdentity`；成功返回三段证据链（R、回执、身份 recordHash）。
- `AgentRunnerLaunchReconfirmer`（值语义）：复用 U1 抽取的包级 `validateAgentRunnerAuthorization` / `validateAgentRunnerBudget`，只做"授权 active、未撤销、未过期 + reservation outstanding"的窄重检，错误链保留底因。
- store 扩展：`RecordLaunchReceipt` / `LaunchReceipt` / `RecordLaunchIdentity` / `LaunchIdentity`；写路径沿用 R/C 纪律（先只读重放解析 → durability 闸门 → 绑定校验 → no-replace + 有界重读），回执与身份为 store-local 记录（独立域分隔摘要，不上 wire、不进 R/C recordHash）；`launches/` 已纳入生产构造器、`bootstrapReplayStoreRoot`、`verifyPathSafetyLocked` 运行期复检与所有权发布前置枚举。
- U1 抽取：A1 私有校验提升为包级函数、admission 调用点改写，A1 测试守护行为不变。
- 聚焦测试：跨 store 24 并发单赢家（恰 1 个 first-launch / 23 个 already-launched）、no-replace 碰撞双分支（同 payload→replay、异 payload→conflict）、身份无意图=corruption、未证明平台写入拒绝与纯重放允许、runId 与全字段绑定（写冲突 + 读 corruption）、零值 store fail-closed；dispatcher 全分支（R-only/intent-only/intent+identity/终态阻断；撤销含"R 保留、无回执、不 spawn、不自动重试"；启动失败含"意图保留、绝不重启"；launchId 错配与身份冲突→`ErrAgentRunnerLaunchIdentityUnproven` 且不 kill 不重试；并发 8 恰 1 次 spawn；绑定错配不写任何记录）；Unix 专项（launches symlink 构造期拒绝、运行期换链拒绝并证外部目录为空、悬空槽有界重读收敛、sync 失败无回滚且重开可重放）。

## 执行结果

### Windows 本机执行（原生证据）

| 命令 | 结果 |
|---|---|
| `gofmt -l internal` | 干净（无输出） |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test -count=1 ./...` | 全部通过 |
| `go test -count=1 -run 'AgentRunnerLaunch\|AgentRunnerReplayDispatcher\|AgentRunnerReplayStoreLaunch\|AgentRunnerReplayRoot' ./internal/adapters` | 通过 |

说明：A2 并发反例（dispatcher 8 并发、store 8 并发、跨 store 24 并发）已在 Windows 原生执行；Windows 无 gcc 未跑 `-race`，Linux `-race` 回归待推送后由 CI 承担。Unix 专项测试（symlink 换链、收敛、sync 失败）由 GitHub Actions Ubuntu 承担。

### 审查结论

- V4 Pro 预审：有条件通过（无 High）。M1（回执跨 store 单赢家与碰撞分类无真实反例）与 M2（身份无意图 corrupt 分支无测试）已整改并补测；L1–L5、L7 已整改（全字段显式交叉校验、launches sync 失败反例、launchId 错配反例、Windows launches junction 换链反例、缺 marker + launches 非空 corruption 反例、`launches/` 纳入 `bootstrapReplayStoreRoot`）；L6（unproven 排序语义）与 L8（零 `StartedAt` 由 A6 端口实现保证）记录为已知边界。
- Codex 独立终审（含测试反例逐条审计）：有条件通过（无 High/Medium）；唯一 Low 尾项（显式绑定补 RunID/AuthorizationHash/BudgetHash 与回放分支 persisted runId 校验）已整改；终审确认不存在双重 spawn 路径。

## 已知边界

- W1/W2 不自动接管：意图已被占用（含并发首派发窗口）与撤权重确认失败都只返回可识别阻断，不做任何自动接管或补偿；接管与结算属 A3。
- recheck→spawn 残余窗口（U4 已拍板）：重确认与真实 `Start` 之间存在不可消除的窗口，本切片只保证窗口内状态可审计、失败 fail-closed。
- 未证明平台排序语义：回执与身份在 unproven 平台允许"无写纯重放"、拒绝一切新写入（与 `RecordCompletion` 同规则）；`RecordRequest` 因涉及首次 dispatch 资格而先拒 unproven。
- 零 `StartedAt`：端口契约声明"返回 result 仅当已 spawn 且身份齐备"，零值由 A6 实现方保证；本切片不做额外拒绝。
- 自动重试口径：契约只"允许"端口契约性未 spawn 失败的同进程重试，本切片选择更保守的永不自动重试，失败后仅人工或上层处置。

## 明确未执行事项

- 未发布 completion，未接入 Engine，未运行真实候选，未实现接管/结算（属 A3 及以后）；
- 未调用真实模型 CLI，未访问网络，未读取 SecretStore；
- 本报告成文时未执行 commit、push 或 publish；CI 证据待按纪律获得同一轮授权后回写；
- T027 仍为 `BLOCKED / NOT IMPLEMENTED`，AT-23 未通过。

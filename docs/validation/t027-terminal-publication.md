# T027 A3 terminal publication and settlement 前置验证报告

[English](t027-terminal-publication_EN.md)

日期：2026-09-14。结论：`PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`。

## 验证范围

本报告只覆盖 T027 离线切片 A3 terminal publication and settlement，不宣称 AgentRunner 生产实现。A3 已交付：

- 契约冻结（先于实现）：`docs/CONTRACTS.md` / `CONTRACTS_EN.md` 新增终局发布段落——四阶段链（terminal-intent → 账本 `Settle` → 发布 C → terminal-closure）、二文件 store-local 协议（不进 wire、不参与 R/C recordHash）、账本按 `idempotencyKey` 去重、结算证据须含被判定的 C 与 R 摘要、发布耐久未证明平台在任何结算之前拒绝新发布且完整链仅无写重放、无意图的 C 为孤儿且不得回写收养、reservation 已被异键结算占用不得发布 C、unknown 结算保留占用、恢复只依据 store-local 记录（账本只作强制者）、接管仍属后续 operator 切片。
- `AgentRunnerTerminalPublisher`（adapters，值语义）：`PublishTerminal` 四阶段编排与 `convergeTerminal` 崩溃恢复；发布前复用 `ValidateAgentRunnerCompletionBinding`（含 resume 会话绑定，杜绝"结算后才被 C 拒绝"的死锁）；结算计划由 (R, C, usage) 确定性派生（幂等键、证据组合与去重、unknown 规范化）；依赖 `*tickets.CostLedger` 与测试单发 hooks（afterIntentWrite/afterSettle）。
- store 扩展：`terminals/terminal-intent.<requestId>.jsonl` 与 `terminals/terminal-closure.<requestId>.jsonl` 记录、读写与绑定校验；写路径遵循既有纪律（只读重放解析 → durability 闸门 → 绑定校验 → no-replace + 有界重读收敛）；`terminals/` 纳入 bootstrap、生产构造器、运行期路径安全复检与所有权发布前置枚举。
- 幂等键：`settle-<20hex>` = 域 `proofrail:agent-runner-settlement-key:1` 上 (requestId, C.recordHash) 摘要；`EntryID` 同值、`CompletionID` 派生；重复回执不可重复结算由账本 `settlementByKey`/`settlementByReservation` 强制。
- 聚焦测试：顺序耦合（hook 断言 intent 可见时未结算、C 未发布；结算后 C 仍未发布，且两 hook 必被触发）、幂等重放、8 并发单赢家（含异时钟栅栏并发与顺序确定性收敛）、崩溃窗口 i/ii/iii 恢复（恰一次结算、closure 补齐）、orphan fail-closed、异键结算拒绝、unknown 保留（`UnknownHoldReservations()==1`、`OutstandingReservations()==0`、摘要 1000）、结算矩阵 5 例、charged>reserved、不同 outcome 同 request 冲突、unproven 预结算拒绝（fresh 与 converge 两档）、unproven 完整链严格无写重放、closure 键失配与无 C 的 corruption、恢复要求持久 R、resume 会话正反例、证据预校验/去重/负数拒绝、跨 store 24 并发单赢家、Unix 构造期 symlink 拒绝/运行期换链/悬空槽收敛/sync 失败无回滚/proven 全链。

## 执行结果

### Windows 本机执行（原生证据）

| 命令 | 结果 |
|---|---|
| `gofmt -l internal` | 干净（无输出） |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test -count=1 ./...` | 全部通过 |
| `go test -count=1 -run 'AgentRunnerTerminal\|AgentRunnerReplayStoreTerminal' ./internal/adapters` | 通过 |

说明：A3 并发反例（8 并发、异时钟并发、跨 store 24 并发）已在 Windows 原生执行；Windows 无 gcc 未跑 `-race`。Unix 专项测试（symlink、换链、收敛、sync 失败）由 GitHub Actions Ubuntu 承担，CI 证据待推送后回写。

### 审查结论

- V4 Pro 前置分析（本切片必调）：四阶段协议、崩溃恢复算法、拒绝矩阵与 Q1–Q10 裁定全部采纳（含"intent 内嵌完整 C""不新增账本查询路径""durability 预检提前到 Settle 之前"等拍板）。
- V4 Pro 实现预审：有条件通过（无 High）。H1（resume 会话绑定缺口会导致结算后 C 无法发布）按统一绑定校验整改；M1（异时钟并发误报冲突）改为语义槽位匹配并补异时钟测试；M2（调用方证据可在写意图后才被账本拒绝）改为写意图前预校验 + 确定性去重 + 精确成员匹配；M3（converge+unproven 无测试）、L1（closure 无 C 的 corruption）、L5（恢复要求持久 R）全部整改闭环。
- Codex 独立终审（含测试反例审计）：有条件通过（无 High）。两个 Medium 整改闭环：unproven 完整链改为严格无写重放（不再触碰账本）；converge 复用完整 closure 绑定校验（含 settlement key），`RecordTerminalClosure` 增加 intent runId 复核。三个 Low 整改：异时钟并发测试加起跑栅栏并补顺序确定性用例、happy-path 增加 hook 触发断言。补齐三条缺失反例：完整链无写重放、closure 键失配 corruption、恢复要求持久请求记录。
- Codex 整改复审（2026-09-15，复审对象=整改后代码，依据指令"中高危修复后必须重新提交复审、不得自行判定通过"）：**PASS——中高危清零、可从审查角度签收**。复审同时记录 3 条低风险测试严谨性尾项（write-free 反例的探针分支不可达、键失配用例未先发布 C 因而变量未隔离、`RecordTerminalClosure` 的 intent runId 复核缺专项回归用例），均低于中高危门槛，按条款不阻断签收。
- 尾项整改与追加确认（2026-09-15，经用户同轮授权）：3 条低风险测试严谨性尾项已修复——发布器新增未导出测试钩子 `beforeSettle`（结算唯一入口 `settleTerminalIntent` 的进入探针，生产路径永不设置），write-free 反例改为在 `beforeSettle`/`afterSettle` 双致命探针下断言；closure 键失配反例先发布 C 完成变量隔离；新增 `RecordTerminalClosure` 的 intent runId 专项回归用例。Codex 追加确认（超出 1+1 额度，用户明确授权）：**PASS——三项均为回滚即红（revert-sensitive）的真实反例，无新增缺陷，中高危清零、可签收**。本机机械确认：`gofmt`/`go build`/`go vet`/`go test -count=1 ./...` 全绿。

## 已知边界

- 结算计划不可修复性：terminal-intent 为 no-replace 单赢家，计划失败后同 request 的修正属 operator/后续切片（fail-closed 设计，槽位不重写）。
- 恢复语义：unproven 平台仅完整链无写重放；不完整链（含已结算未发布 C）拒绝，须回到 proven 实例恢复；恢复是幂等补写自有槽位，不是接管（W1/W2 口径不变）。
- unknown 结算为终态，本切片不定义 upgrade 语义（留 A5）。
- 账本无按 reservation 查询路径：恢复只依据 store-local 记录，账本只作幂等/冲突强制者。
- 零 `StartedAt` 等由 A6 端口实现保证项不变。
- 终局发布是否成功不影响 task PASS 语义；task 状态推进属 A4。

## 明确未执行事项

- 未接 Engine、未运行真实候选、未实现 chain terminal 路由或 postflight 接受（A4/A5）；
- 未调用真实模型 CLI，未访问网络，未读取 SecretStore；
- 本报告成文时未执行 commit、push 或 publish；CI 证据待按纪律获得同一轮授权后回写；
- T027 仍为 `BLOCKED / NOT IMPLEMENTED`，AT-23 未通过。

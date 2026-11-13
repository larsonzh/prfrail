# ⑤ 预审 r1 原文（deep-reasoner + DeepSeek V4 Pro）

## 一、实现与 design-r1 一致性（要求 1）—— 一致

- 三处注入点全部落地：`runtime.go:45-51`（`ExecuteNoopRun`，`if progress != nil` 才包裹）、`runtime_agent.go:352-361`（`ExecuteAgentRunnerRun`）、`runtime_supervised.go:626-635`（`newSupervisedEngine`）。
- nil 守卫三处均在；nil 时 `Options.Events` 指向包裹前的同一个 `*chain.FileEventStore`。
- 只读重放点未装饰：`runtime_supervised.go:915-922` `supervisedOpenWindowAt` 自行 `chain.NewFileEventStore`+`Load`，不经装饰器；console 内再无其它 `Events:` 注入（grep 仅 4 处：3 生产注入点 + 1 测试）。
- SealSink 全覆盖：`progressStore` 实现 `SealTaskEvidence`/`LoadTaskSeal` 纯转发（progress.go:57-64），且有 `_ chain.EventStore`/`_ chain.SealSink` 编译期断言。`engine.go:340`（静默关封存重放）、`engine.go:374`（`ErrSealUnavailable` fail-closed）、`seal_verify.go:25`（静默跳过恢复核验）三处断言作用于 `engine.options.Events`=装饰后对象 ⇒ 全部命中转发，无静默降级。
- 打印在 inner `Append` 成功后、`Load`/Seal 均不改事件字节。

## 二、反例有效性（要求 2）—— 成立，一处「绿但未验证」

- R1（:160 起）：text 模式 stdout 掩码 runDir 后逐字节比较 + 双 stderr 断言 + 三实体行计数，真实有效；变异开关（删 nil 守卫）会使 plain run 对 nil writer `Fprintln` panic 变红。
- R2（:182）：`json.Valid` + 非空 stderr + 含 `progress: ` 拒绝 + 逐字节相等，真实有效；同一 run-dir+`RemoveAll` 是隔离 flag 的正确做法（见核实 (a)）。
- R3（:225）：`len(plainBytes)==0` 守卫防恒真；吞错后 `Append` 返回 nil + 持久字节相等，真实有效；删 `_, _ =` 吞错即红。
- R4（:270）：健康 sink 断言 `sink.Len()>0` + 字节相等，有效；但缺 R3 的非空守卫（见发现 F5）。
- R5/R7：格式快照与 usage/flag 解析，有效。
- **R6（:374）「SealSink 转发断言」名不副实**：只调用 `LoadTaskSeal`（读路径），从不调用 `SealTaskEvidence`（写路径）——见发现 F1。

## 三、越界判定（要求 3）—— 无越界

diff 仅 touch `internal/console/**`（生产+测试）与台账 `docs/t027/REMAINING_SLICES{,_EN}.md`（起跑登记，切片义务）。`internal/chain`/`adapters`/`gates`/`schemas`/`go.mod`/`CONTRACTS` 零改动。

## 四、默认行为不变（要求 4）—— 成立

`--progress` 缺省 false ⇒ `progressWriter=nil` ⇒ 三路均为原 store 同一对象，stdout/stderr 字节不变。`executeSupervisedRun`/`executeAgentRunnerRun` 为未导出方法、各仅一个调用点；`Progress` 为 options 新增字段（零值 nil），全部既有测试与调用点（readiness_gate_test、multitask 系列等）零改动仍绿。`ExecuteNoopRun` 为导出函数，尾参加法对仓库内 5 处调用点已全部更新（grep 无遗漏），仅对外部消费方是编译期破坏（见 F6）。

## 发现清单

F1. `Low | internal/console/progress_test.go:374 | design-r1 §5「SealSink 转发断言」 | 修复：R6 先经 `store.SealTaskEvidence` 与 `sink.SealTaskEvidence` 各写一枚 seal，再比对两侧三元组（或比对返回值）`——可证伪性：将 `progressStore.SealTaskEvidence` 改为 `return "", nil` 直接丢弃 inner，现有测试全绿（证明写路径转发无任何断言）。

F2. `Low | internal/console/supervised_cmd.go:41、cli.go:601 | design-r1 §7「supervised_cmd.go 透传 Progress」「三路透传」 | 修复：给 AgentRunner/Supervised 各加一条 Progress≠nil 的 wiring 测试（stderr 现 progress 行且 stdout 与无 progress 逐字节一致）`——可证伪性：删除 `options.Progress = progress` 或 `Progress: progress`，全部测试仍绿（证明这两路透传零覆盖；R1–R7 全走 noop 路由）。

F3. `Low | internal/console/progress.go:41-47（Append） | 「progress must never change the outcome of a run」 | 修复：可选防御——打印包入 `defer recover()` 并在注释声明「打印路径不得引入可 panic 结构」`——可证伪性：在 `formatProgressLine` 人为 `panic`，R3 类测试整体崩溃而非变红（证明无 recover 防护）；删吞错则 R3 变红（证明吞错维度已防）。**核实 (b) 裁定：不成立、不报 Medium+**。当前 `evidence.Event` 全为值类型 string/int（event.go:44-57），`formatProgressLine` 无指针解引用/无 map/slice 索引/无类型断言，对任何合法 `StateEvent` 不可能 panic；「先落盘后打印 + 吞 writer 错误」已满足「不改结局」承诺的全部可达维度。

F4. `Info | internal/console/progress_test.go:182（R2） | 核实 (a) 裁定：不 flaky，无需修复 | 可选：在 R2 注释声明依赖 `newTestCLI` 固定时钟`——可证伪性：删 `newTestCLI` 的 `cli.now` 固定覆盖（cli_test.go:405-407），两次运行 `OccurredAt` 不同 ⇒ `lastEventHash` 不同 ⇒ JSON 必不等（确定性变红而非偶发）。全链路时间取自 `runtime.now()`（runtime_noop.go:70/82）与 engine `options.Clock`（state.go:72），ID 为计数器（runtime_noop.go:118-121），`RunSummary` 无时间字段（runtime.go:19-29），`defaultClock`/`defaultIDs` 均未被该路径使用——无任何墙钟源。

F5. `Info | internal/console/progress_test.go:270（R4） | 「装饰器不改事件字节」的恒真退化风险 | 修复：仿 R3 加 `len(plainBytes)==0` 守卫`——可证伪性：`progressStoreEvents` 改返回空切片，R3 红（有守卫）、R4 仍绿（无守卫）；现 fixture 硬编码 3 条事件，实际无退化。

F6. `Info | internal/console/runtime.go:37 | 导出函数 `ExecuteNoopRun` 尾参加法 | 修复：若在意外部兼容可改 options 结构体；仓库内无需动作`——可证伪性：编译期事实，grep 全仓库 6 处引用（1 定义 + 5 调用）均已是新签名。

## 结论

- 要求 1/3/4 全部通过；要求 2 反例真实有效（仅 R6 写路径转发未验证、R4 缺非空守卫，均 Low/Info）。
- 特别核实 (a)：R2 不 flaky（固定时钟全链路承重）；(b)：panic 不可达，不成立，故均不按 Medium+ 报。
- 附带 3 Low + 3 Info 建议供 ⑥/⑦ 参考，可在原位顺手修或留作记录。

PRE-REVIEW: PASS

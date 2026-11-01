# ⑦ 终审（原样归档）

> 载体：runSubagent(agentName="independent-reviewer", model="GPT-5.3-Codex (copilot)")。本文件为工具返回的原始内容（大结果落盘件），未改写。

## A 终审结论
FINAL REVIEW: PASS（无 Critical/High/Medium 缺陷；存在 1 项 MINOR 文案漂移，已列入登记清单）

## B 缺陷
1. MINOR：测试失败提示中的 README 行号与命令片段已过时，易误导排障定位。  
证据：测试仍写 README.md:9 与旧片段见 [prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L163)，实际命令在 [prfrail/examples/no-ai-preview/README.md](prfrail/examples/no-ai-preview/README.md#L10)。

## C 核验结论
1. E2E 是否真证明“零件贯通”：通过。  
三条链路是独立可证伪的，不是“人为跳过 T3”。  
证据：T1 通过/T2 失败/链失败分别断言于 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L37)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L40)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L43)；T3 未 materialize/未投影/未派发分别在 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L47)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L66)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L73)。  
“不是删掉 T3 导致没跑”也成立：夹具定义含 task-three 于 [prfrail/internal/chain/engine_test.go](prfrail/internal/chain/engine_test.go#L206)，且引擎在任务报错即停止后续任务循环见 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L149) 与 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L162)。

2. 零模型是否机械：通过（并且预览腿已如实声明计数器非哨兵）。  
引擎腿：零模型判据链闭合。StepIntents 必须为 nil 的前置断言在 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L165)；生产中只有 StepIntents 才会写入执行目标见 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L431) 与 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L445)。测试又逐步断言默认执行目标/无 AgentRunnerFacts/无 dispatch 事件见 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L180)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L183)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L188)；而 dispatch 事件仅在 facts 存在时写入见 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L480) 与 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L483)。  
预览腿：计数器确为静态常量（非哨兵）已在测试注释明确，且用可证伪条件补强。静态赋零在 [prfrail/internal/console/preview.go](prfrail/internal/console/preview.go#L161)；“非哨兵”声明与机械断言在 [prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L157)、[prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L201)、[prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L204)、[prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L207)、[prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L218)、[prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L229)。

3. 父快照与事件链：通过，且与生产语义一致。  
生产语义：baseline 写入 AcceptedParent 见 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L118)；任务 materialize 使用当前 AcceptedParent 见 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L195)；任务通过后推进到 candidate 见 [prfrail/internal/chain/engine.go](prfrail/internal/chain/engine.go#L332)。  
测试对 baseline→acceptedOne→acceptedTwo→acceptedThree 做了逐点绑定断言见 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L98)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L101)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L104)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L107)。  
事件链验证覆盖 VerifyEventChain、Sequence 单调、LastEventHash 链接与 3 条 seal 计数，见 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L128)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L133)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L142)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L145)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L154)。VerifyEventChain 的生产约束（序号、前哈希、状态连续性）见 [prfrail/internal/evidence/event.go](prfrail/internal/evidence/event.go#L96)、[prfrail/internal/evidence/event.go](prfrail/internal/evidence/event.go#L110)、[prfrail/internal/evidence/event.go](prfrail/internal/evidence/event.go#L131)、[prfrail/internal/evidence/event.go](prfrail/internal/evidence/event.go#L140)。

4. 并发腿强度：通过。  
同 launcher/同 runRoot 已满足：同一 runRoot、同一 replay store、同一 launcher 见 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L73)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L78)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L82)。  
8 路齐发与“恰一赢家”检查在 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L95)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L100)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L129)。  
输家只允许两个互斥哨兵见 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L120)；赢家存活/停后释放/可重启在 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L137)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L146)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L156)。  
未使用测试侧墙钟轮询；仅子进程生命周期使用 sleep，见 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L43)。  
补充：runRoot 标记机制与 replay 绑定语义在 [prfrail/internal/console/runtime_agent.go](prfrail/internal/console/runtime_agent.go#L888) 与 [prfrail/internal/adapters/agent_runner_replay_root.go](prfrail/internal/adapters/agent_runner_replay_root.go#L237)。

5. examples 改动正当性：通过。  
新增夹具是修真实缺陷，不是“造数据刷绿”：测试明确依赖该文件，缺失即失败，见 [prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L161) 与 [prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L163)。  
README 改动为最小必要（仅命令路径修正）见 [prfrail/examples/no-ai-preview/README.md](prfrail/examples/no-ai-preview/README.md#L10)。  
夹具内容是两任务纯 noop，无执行副作用见 [prfrail/examples/no-ai-preview/proofrail.chain.json](prfrail/examples/no-ai-preview/proofrail.chain.json#L9)、[prfrail/examples/no-ai-preview/proofrail.chain.json](prfrail/examples/no-ai-preview/proofrail.chain.json#L14)、[prfrail/examples/no-ai-preview/proofrail.chain.json](prfrail/examples/no-ai-preview/proofrail.chain.json#L20)、[prfrail/examples/no-ai-preview/proofrail.chain.json](prfrail/examples/no-ai-preview/proofrail.chain.json#L25)。  
只读字节复核结果：JSON 为 1044 bytes、无 BOM、LF；README 为 BOM+LF（符合仓库 Markdown 约定）。

6. 登记诚实性：当前无法“已登记”判定，只能给“必须登记”判定。  
原因：本轮未提供 ⑧ 报告正文供核对。结论：该项需在放行条件中强制落地。

7. 放行判定：可进入 ⑧/⑨。  
前提：按 §D 全量登记；若缺任一项，视为未满足终审放行条件。

## D 放行条件
1. 在 ⑧ 报告显式登记本轮 1 项 MINOR：测试失败提示中的 README 行号/命令片段漂移（见 §B）。  
2. 原样登记 ⑤ 的 2 项 MINOR（已处置）及处置证据。  
3. 原样登记 ⑥ 的 4 项未覆盖面（UI E2E、跨 OS/locale、证据包终态、真实 preview 命令执行证据）。  
4. 原样登记 ① 的 4 项未确证说明（含单任务限制不改及 A6 mirror 冲突风险边界）。  
5. 在 ⑧ 中补一句“本片为装配验证片，生产代码零改动”。  
6. 在 ⑧ 中补一句“examples/no-ai-preview 为测试夹具修复，不参与默认运行路径”。

## E 可证伪性
1. 若把 no-ai 夹具任一步从 noop 改为 code，预览腿应在 runnable/outcome/blockingEvidence/step-kind 断言处转红，锚点见 [prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L201) 到 [prfrail/internal/console/runtime_zero_model_e2e_test.go](prfrail/internal/console/runtime_zero_model_e2e_test.go#L229)。  
2. 若让引擎出现 AgentRunner 派发，零模型腿应在 execution target、facts 或 dispatch reason 任一断言处转红，锚点见 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L180)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L183)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L188)。  
3. 若 T2 失败后仍放行到 T3，T3 未启动相关断言会转红，锚点见 [prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L47)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L66)、[prfrail/internal/chain/engine_zero_model_e2e_test.go](prfrail/internal/chain/engine_zero_model_e2e_test.go#L73)。  
4. 若并发互斥失效（多赢家或停后不可复用），并发腿会在赢家计数、冲突哨兵、重启可用性处转红，锚点见 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L129)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L120)、[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L156)。

FINAL REVIEW: PASS## A 终审结论
FINAL REVIEW: PASS（附带登记条件，见 §D）

## B 缺陷
1. MINOR：断言失败提示中的 README 行号与命令片段已过时。  
证据：测试文案仍写 README.md:9 且片段为 --chain .\proofrail.chain.json，见 [prfrail/internal/console/runtime_zero_model_e2e_test.go#L163](prfrail/internal/console/runtime_zero_model_e2e_test.go#L163)；实际 README 命令在 [prfrail/examples/no-ai-preview/README.md#L10](prfrail/examples/no-ai-preview/README.md#L10)。  
影响：仅失败文案漂移，不影响运行语义与门禁强度。  
结论：非阻断。

## C 核验结论
1. E2E“零件贯通”是否成立：成立，且可单点证伪。  
主链断言分别独立校验 T1 通过/T2 失败/链失败（[prfrail/internal/chain/engine_zero_model_e2e_test.go#L37](prfrail/internal/chain/engine_zero_model_e2e_test.go#L37), [#L40](prfrail/internal/chain/engine_zero_model_e2e_test.go#L40), [#L43](prfrail/internal/chain/engine_zero_model_e2e_test.go#L43)）；  
T3 未启动从四个面独立验证（workspace、事件、投影、step dispatch），见 [#L47](prfrail/internal/chain/engine_zero_model_e2e_test.go#L47), [#L60](prfrail/internal/chain/engine_zero_model_e2e_test.go#L60), [#L66](prfrail/internal/chain/engine_zero_model_e2e_test.go#L66), [#L73](prfrail/internal/chain/engine_zero_model_e2e_test.go#L73)；  
生产语义上 runTask 报错即停止后续 task，见 [prfrail/internal/chain/engine.go#L149](prfrail/internal/chain/engine.go#L149), [#L162](prfrail/internal/chain/engine.go#L162)。  
不存在“人为跳过 T3”的假证明：共享夹具确实含 T3，见 [prfrail/internal/chain/engine_test.go#L206](prfrail/internal/chain/engine_test.go#L206)。

2. 零模型判据是否机械：引擎腿成立；预览腿已如实声明静态计数器且补了机械哨兵。  
引擎腿：明确要求 StepIntents 为空并校验默认执行目标、无 AgentRunner facts、无 agent-dispatched 事件（[prfrail/internal/chain/engine_zero_model_e2e_test.go#L165](prfrail/internal/chain/engine_zero_model_e2e_test.go#L165), [#L180](prfrail/internal/chain/engine_zero_model_e2e_test.go#L180), [#L183](prfrail/internal/chain/engine_zero_model_e2e_test.go#L183), [#L188](prfrail/internal/chain/engine_zero_model_e2e_test.go#L188)）。  
且生产代码只有在 StepIntents 赋值后才会进入 AgentRunner 分发路径（[prfrail/internal/chain/engine.go#L431](prfrail/internal/chain/engine.go#L431), [#L445](prfrail/internal/chain/engine.go#L445), [#L480](prfrail/internal/chain/engine.go#L480), [#L483](prfrail/internal/chain/engine.go#L483)；reason 常量见 [prfrail/internal/chain/agent_runner_terminal.go#L37](prfrail/internal/chain/agent_runner_terminal.go#L37)）。  
预览腿：计数器确为静态常量 0（[prfrail/internal/console/preview.go#L161](prfrail/internal/console/preview.go#L161), [#L164](prfrail/internal/console/preview.go#L164)），测试也明确声明“非哨兵”；真正可证伪的是 RunnableInCLI/Outcome/BlockingEvidence/step kind/unknowns（[prfrail/internal/console/runtime_zero_model_e2e_test.go#L201](prfrail/internal/console/runtime_zero_model_e2e_test.go#L201), [#L204](prfrail/internal/console/runtime_zero_model_e2e_test.go#L204), [#L207](prfrail/internal/console/runtime_zero_model_e2e_test.go#L207), [#L218](prfrail/internal/console/runtime_zero_model_e2e_test.go#L218), [#L229](prfrail/internal/console/runtime_zero_model_e2e_test.go#L229)）。

3. 父快照与事件链一致性：成立。  
生产路径：baseline 写 AcceptedParent、每 task 取 parent、task accepted 后推进 parent（[prfrail/internal/chain/engine.go#L118](prfrail/internal/chain/engine.go#L118), [#L195](prfrail/internal/chain/engine.go#L195), [#L332](prfrail/internal/chain/engine.go#L332)）。  
测试断言 baseline→acceptedOne→acceptedTwo→acceptedThree 完整闭环（[prfrail/internal/chain/engine_zero_model_e2e_test.go#L98](prfrail/internal/chain/engine_zero_model_e2e_test.go#L98), [#L101](prfrail/internal/chain/engine_zero_model_e2e_test.go#L101), [#L104](prfrail/internal/chain/engine_zero_model_e2e_test.go#L104), [#L107](prfrail/internal/chain/engine_zero_model_e2e_test.go#L107)）。  
事件链校验覆盖 sequence、前驱 hash、状态连续性（调用点 [prfrail/internal/chain/engine_zero_model_e2e_test.go#L128](prfrail/internal/chain/engine_zero_model_e2e_test.go#L128)；规则见 [prfrail/internal/evidence/event.go#L96](prfrail/internal/evidence/event.go#L96), [#L110](prfrail/internal/evidence/event.go#L110), [#L131](prfrail/internal/evidence/event.go#L131), [#L140](prfrail/internal/evidence/event.go#L140)）。  
封存记录 3 条断言与任务接受语义一致（[prfrail/internal/chain/engine_zero_model_e2e_test.go#L150](prfrail/internal/chain/engine_zero_model_e2e_test.go#L150), [#L154](prfrail/internal/chain/engine_zero_model_e2e_test.go#L154)）。

4. 并发腿强度：达到“同 launcher/同 runRoot、恰一赢家、可复用槽位”的目标。  
同 runRoot + marker + replay store + 单 launcher：见 [prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L73](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L73), [#L75](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L75), [#L78](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L78), [#L82](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L82)。  
8 路并发同起跑门（[prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L95](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L95), [#L100](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L100), [#L110](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L110)）；输家仅允许两类独占哨兵（[#L120](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L120)）；一赢七负（[#L129](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L129)）；活性-冲突-停止-复用链路完整（[#L137](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L137), [#L142](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L142), [#L146](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L146), [#L156](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L156)）。  
无墙钟/轮询同步；仅子进程生命周期 sleep（[#L43](prfrail/internal/console/runtime_agent_concurrent_e2e_test.go#L43)）。

5. examples 变更正当性：正当，且是修真实缺陷。  
README 公开命令已指向实际存在的链文件（[prfrail/examples/no-ai-preview/README.md#L10](prfrail/examples/no-ai-preview/README.md#L10)）；测试明确把“示例文件必须存在”作为门禁（[prfrail/internal/console/runtime_zero_model_e2e_test.go#L161](prfrail/internal/console/runtime_zero_model_e2e_test.go#L161), [#L163](prfrail/internal/console/runtime_zero_model_e2e_test.go#L163)）。  
新增链文件是纯 noop/manual 两任务（[prfrail/examples/no-ai-preview/proofrail.chain.json#L9](prfrail/examples/no-ai-preview/proofrail.chain.json#L9), [#L14](prfrail/examples/no-ai-preview/proofrail.chain.json#L14), [#L20](prfrail/examples/no-ai-preview/proofrail.chain.json#L20), [#L25](prfrail/examples/no-ai-preview/proofrail.chain.json#L25)），无 AI 执行副作用。  
只读字节校验：JSON 1044 bytes、无 BOM、LF；README 789 bytes、UTF-8 BOM、LF（符合仓库编码约定）。

6. 登记诚实性：当前无法对“已登记”做完成态判定（未见 ⑧ 最终报告），但可对“必须登记项”做闭集约束，见 §D。若 ⑧ 缺任一项，应回退为 BLOCKED。

7. 放行判定：可进入 ⑧/⑨，前提是 §D 登记清单完整落地。

## D 放行条件
1. 在 ⑧ 报告显式登记：⑤ 的 2 项 MINOR 已处置（测试改造、台账措辞修订）。  
2. 在 ⑧ 报告显式登记：⑥ 的 4 项未覆盖面（UI E2E、跨 OS/locale、证据包终态、真实 preview 命令执行证据）。  
3. 在 ⑧ 报告显式登记：① 的 4 项未确证项与其边界理由。  
4. 在 ⑧ 报告新增登记本轮 MINOR：断言失败文案的 README 行号/命令片段漂移（见 §B）。  
5. 报告中保留“生产代码零改动”的声明与证据口径，避免把测试/示例调整误记为运行时行为修复。

## E 可证伪性
1. 把示例链任一步从 noop 改为 code，预览腿应立即转红（RunnableInCLI/Outcome/BlockingEvidence/Unknowns/step kind 至少一项失败）。  
2. 给引擎装入 StepIntents 或产生 agent-dispatched 事件，零 AgentRunner 派发腿应转红。  
3. 破坏事件序列或前驱 hash，VerifyEventChain 应转红。  
4. 并发腿若出现 2 个赢家或 stop 后不可重启，应转红。  
5. 本次终审为只读核验：Diff=none；未执行 go test/go vet 日志采集（遵守“禁止写入/重定向”约束）；残余风险为跨平台调度抖动与文案漂移，不影响当前放行结论。

FINAL REVIEW: PASS
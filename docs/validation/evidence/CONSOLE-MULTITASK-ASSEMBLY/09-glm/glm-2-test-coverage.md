**P50 CONSOLE-MULTITASK-ASSEMBLY 测试覆盖矩阵**（只读检索，未改任何文件）

前缀：red=internal/console/runtime_multitask_red_test.go；agent=…/runtime_multitask_agent_test.go；sup=…/runtime_multitask_supervised_test.go；driver=…/runtime_multitask_driver_test.go；reg=internal/adapters/agent_runner_intent_registry_test.go；e2e=internal/console/runtime_agent_e2e_test.go；supE=internal/console/runtime_supervised_test.go

列：①单任务兼容 ②2-task接受 ③0-task拒绝 ④逐task step形状门 ⑤记录注册表绑定 ⑥跨记录drift ⑦串行顺序 ⑧停机屏障三态 ⑨止损绑定 ⑩六事实独立归属 ⑪五事实位置序 ⑫事件全序 ⑬R-C隔离 ⑭空摘要拒绝先于marker

行：A=red.TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape；B=red.TestSupervisedAssemblyAcceptsMultiTaskDefinitionShape；C=red.TestMultiTaskDefinitionShapeIsOtherwiseValid；D=agent.TestAgentRunnerMultiTaskSerialOfflineHarness；E=sup.TestSupervisedMultiTaskSerialHarness；F=driver.TestAgentRunnerDriverTwoTaskSerialAssembly/"a proven stop lets the next task start"；G=driver/…"a live launch blocks the next task"；H=driver/…"a missing stop report fails closed"；I=driver.TestAgentRunnerDriverStopsTheLaunchOfTheTaskThatJustFailed；J=driver.TestAgentRunnerAssemblyRefusesNonInitialAttempt；K=reg.TestAgentRunnerStepIntentRegistryBindsEachRecordToItsOwnStep；L=reg.TestAgentRunnerStepIntentRegistryRejectsCrossRecordDrift；M=reg.TestNewAgentRunnerStepIntentRegistryReusesPerRecordAdmission；N=既有e2e受影响（PID报告运行时+funnel护栏）；O=既有supE受影响（scan状态参数化+plan绑定）

|用例|①|②|③|④|⑤|⑥|⑦|⑧|⑨|⑩|⑪|⑫|⑬|⑭|
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
|A|—|✅|✅|✅|—|—|—|—|—|—|—|—|—|✅|
|B|—|✅|✅|✅|—|—|—|—|—|—|—|—|—|✅|
|C|✅|✅|—|—|—|—|—|—|—|—|—|—|—|—|
|D|—|✅|—|—|✅|✅|✅|—|—|✅|—|✅|✅|—|
|E|—|✅|—|—|—|—|✅|—|—|—|✅|✅|—|—|
|F|—|✅|—|—|✅|—|✅|✅|—|—|—|—|—|—|
|G|—|—|—|—|—|—|—|✅|—|—|—|—|—|—|
|H|—|—|—|—|—|—|—|✅|—|—|—|—|—|—|
|I|—|—|—|—|—|—|—|—|✅|—|—|—|—|—|
|J|—|—|—|—|—|—|—|—|—|—|—|—|—|✅|
|K|—|—|—|—|✅|—|—|—|—|—|—|—|—|—|
|L|—|—|✅(a)|—|—|✅|—|—|—|—|—|—|—|—|
|M|—|—|—|✅|—|—|—|—|—|—|—|—|—|—|
|N|—|—|—|—|—|—|—|✅|—|—|—|—|—|—|
|O|—|—|—|—|✅|—|—|—|—|—|—|—|—|—|

注(a)：L③为空注册表（0-record 侧）拒绝；0-task 语义本体在 A/B。N⑧=managedProcessesGone 改 (taskID,attempt) 参数化（e2e:875）+ funnel 护栏 7→8 计入 between-task 停机证明失败路径（e2e:1019）。O⑤=supervisedAssemblyPlan 逐 task 绑定（supE:1426-1435）。

**证据清单（每列≥1条，file:line+断言摘录）**
① red:237 `task %q alone must stop at the request-record check`（C 中每个单 task 切片过 step 选择、停在后置检查）
② red:144 `the two-task definition must stop at the request-record registry`（A 两 task 越过单 task 门达注册表拒绝点）；red:191 `must stop at the run-root precondition`
③ red:156 `a zero-task definition must be refused with %q`（"requires at least one task"；red:199 supervised 同型）
④ red:167 `a task with two steps must be refused per task`；reg:141 `a non isolated-workspace step must be refused`
⑤ reg:69 `task %q received request %q, want %q` + reg:58 `the registry must bind two task/step entries`；agent:252 `the intent registry must bind both records`
⑥ reg:105 `two records sharing one requestId must be refused`（另 reg:100 同 task/step、reg:117 跨 run）；agent:133 `received another record's facts`
⑦ agent:391 `the tasks must interleave serially in one log`；driver:310 `dispatch %d = %q, want %q`；sup:134 `the supervised tasks must interleave serially in one log`
⑧ 三态齐备：proven→driver:288-311（dead PID 种入后 COMPLETED+PASSED==2）；alive→driver:338 `still alive`+driver:342 `the second task must never dispatch`；missing→driver:353 `cannot be proven gone`（ErrInvalidAgentRunnerEvidence fail-closed）
⑨ driver:417 `the mitigation must reclaim task two` + driver:421 `task two's launch must be stopped`（receipts.seen==["request-two"]、stopper==["launch-request-two"]）
⑩ agent:381 `task %q terminal carries a fact of task %q`（异 task 六事实缺席）+ agent:370 位置逐位比对 `terminal evidence[%d]`
⑪ sup:163 `task %q fact %d = %q, want %q (rederived %v)`（经 sup:152 deriveSupervisedFacts 生产路径重推导）；sup:170 `must not share a supervised fact`
⑫ agent:345 `event sequence must increase, got %d after %d`（sup:99 同型）
⑬ agent:404 `the two tasks must own distinct requestIds and request hashes` + agent:407 `must own distinct evidence directories` + agent:418 `replay state = %q`（逐 requestId）
⑭ driver:448 `a precondition refusal must project an empty summary` + driver:451 `must not create the run-root marker`（stat 断言）；red:149 空摘要断言同型

**覆盖缺口**
1. 多记录侧 attempt 维度缺失：K/L/M 与 driver fixture 全部固定 Attempt=1；"某条 record attempt≠1 时注册表/装配拒绝"无用例（J 仅覆盖单记录入口 attempt=2）。
2. R-C 写入路径平台门控：D 的 RecordRequest/RecordCompletion+replay state 断言仅在 PublicationDurability=proven 平台执行（agent:410-420，Windows 走 else 仅记日志），Ubuntu CI leg 兜底——设计内门控，本机矩阵不可证。
3. red A/B 断言空摘要但未直接 stat run-root marker 不存在（仅注释声明 before any I/O，red:23）；stat 类断言只由 J（driver:451）覆盖 agent 侧。
4. F（生产驱动串行腿）未显式断言事件日志序列全序（仅 dispatch/publish 顺序+汇总计数），全序由 D/E 承担。

实得字符数 = 4656
GLM-2: DONE

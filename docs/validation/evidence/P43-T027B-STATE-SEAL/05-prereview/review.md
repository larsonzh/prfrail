# P43-T027B-STATE-SEAL · ⑤ 预审（两轮原样）

载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`；只读。

---

## r1

判定：`PRE-REVIEW: PASS WITH FIXES`

### 发现
| 编号 | 严重度 | 位置 | 问题 | 修法 |
|---|---|---|---|---|
| H1 | High | `recover.go:18`（唯一调用点）；生产装配无 `Recover` 调用 | `verifyAcceptedSeals` 只在 `Recover` 内，生产路径 `New`→`Run` 从不调用 `Recover`，PASSED 任务 seal 缺失时 `Run` 直接跳过并 COMPLETED，恢复校验被绕过 | 在 `New` 内也校验，或装配层强制 `Recover` |
| H2 | High | `agent_runner_terminal_route.go`、`operator_interaction.go` | 三处 `!= 1` 只改比较操作数，下游 entity 仍硬编码 attempt 1，attempt>1 时路由查错键 | 路由侧改用 `*EntityAt(..., attempt)` |
| M1 | Medium | `state_seal_test.go` + fake `LoadTaskSeal` 返回 `os.ErrNotExist` | I1 未真正走到 `SealTaskEvidence` 失败分支 | fake 返回 `ErrSealMissing`，另立失败用例 |
| M2 | Medium | `state_seal_test.go` 断言 `TaskStates["task-one"]` | 实际键为 `taskKey("task-one",1)`，断言永真 | 改用 `taskKey(...)` |
| M3 | Medium | `seal.go` 写/同步失败仅 Close | 留下半写文件，后续 decode 失败判 mismatch，永久 REPAIR | 失败路径 `os.Remove(path)` |
| M4 | Medium | `postflight.go` / `supervised_blackbox.go` facts 提取 | 只按 RunID+TaskID 过滤，跨 attempt 串扰 | 增加 attempt 过滤 |
| L1 | Low | `acceptFromExistingSeal` | 复用 seal 时不校验 decision 值 | 增加 decision 校验 |
| L2 | Low | `seal.go` 并发同名 attempt | 半写窗口期可能误判 mismatch（fail-closed） | 文档化 |
| I1 | Info | `attempt.go` `ContextKey` | K 生产零接线，C5 交付物为死代码 | 接线或登记延后 |
| I2 | Info | `engine.go` Execute 失败路径 | FAILED 事件可携带非法哈希 | 可选校验 |
| I3 | Info | `engine.go` seal 失败路径传候选闭包 | 与既有 review-rejected 同构，可能引用悬空 | 登记 |
| I4 | Info | `attempt.go` `projectionAttempt` | 「最大已见 attempt」不能推导下一 attempt，语义未文档化 | 补注释 |

### 交叉核对
1. 封存先于接受在 `engine.go` 成立，`task-accepted` 仅两处写入且均绑定 seal。
2. 「三处 != 1 断言改造」意图与实现不一致（H2）：拒绝面成立，下游键仍 attempt 1。
3. C5 设计→代码缺口：`ContextKey` 生产零调用（I1）。
4. 测试假件与生产契约一致性：`memoryEvents` 实现 `SealSink` 与 `FileEventStore` 语义一致；`sealFailingEvents.LoadTaskSeal` 用 `os.ErrNotExist` 偏离生产（M1）。
5. 失败语义一致：三条失败路径均不写 PASSED 且链 PAUSED；唯半写残留破坏恢复闭环（M3）。

`PRE-REVIEW: PASS WITH FIXES`

---

## r2

判定：`PRE-REVIEW: PASS`

### F 项闭合核实
| 编号 | 是否闭合 | 证据 |
|---|---|---|
| H1 | 闭合 | `New` 内 `verifyAcceptedSeals` 在 `transition(CREATED)` 之前；`Recover` 保留同校验 |
| H2 | 闭合 | 路由 5 处调用点改 `*EntityAt(..., attempt)`；operator 两函数各一处 |
| M1 | 闭合 | fake 改返回 `ErrSealMissing`，I1 走 `SealTaskEvidence` 失败分支 |
| M2 | 闭合 | 断言键改 `taskKey("task-one", 1)` |
| M3 | 闭合 | Write/Sync/Close 失败各 `os.Remove(path)`，仅在 O_EXCL 新建成功后触发 |
| M4 | 闭合 | 两个 facts 提取函数增 attempt 过滤；`requirePostflightQualifiedReview` 增参 |
| L1 | 闭合 | 复用前校验 `ReviewDecision == "approve"` 且 `PromotionDecision == "completed"` |
| I1 | 登记 | `ContextKey` 仍仅测试引用（后续由 ⑥ 提出并已接线） |
| 新测试 | 闭合 | `TestTerminalRouteRefusesForeignAttempt` |

### 新发现（登记类）
- N1 Low：`cancelRun` 用单 attempt 实体（当时生产装配未接线 resolver，暂不可达）。
- N2 Low：supervised 路径固定 attempt-1 实体（单 attempt 设计假设）。
- N3 Info：`chain.New` 四处调用点均向上传播错误，无误判。

`PRE-REVIEW: PASS`

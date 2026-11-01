# ⑤ 预审（原样归档）

> 载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。transcript 偏移 122385。

# P45-T027B-E2E-ZERO-MODEL ⑤ 预审报告

## A 结论

`PRE-REVIEW: PASS WITH FIXES`

## B 发现

| 编号 | 严重度 | file:line | 问题 | 修法 |
|---|---|---|---|---|
| 1 | MINOR | `internal/console/runtime_zero_model_e2e_test.go:160-175`（对 `internal/console/preview.go:161-168`） | **五项计数器断言是同义反复，不是机械证明**。`BuildPreviewReport` 把 `PreviewCallCounters` 硬编码为全零常量（preview.go:161），因此 `commandCalls/modelCalls/networkCalls/credentialReads/versionProbes==0` 在**任何实现改动下都不可能变红**——即使未来预览管线真的触达模型也不会失败。该测试真正可推导的机械钉子是 `RunnableInCLI`（preview.go:157，`len(executableSteps)==0` 派生）与 `outcome=="ready"`。 | 二选一：① 删除五项计数器断言，改为断言派生量（如 `len(report.Unknowns)==0`、`PreviewRecord.Preview.Outcome=="ready"`、`BlockingEvidence` 为空），保留 `RunnableInCLI`；② 保留但加注释声明「钉住的是 README 承诺的静态常量，非机械判据」，避免误读为模型通道哨兵。 |
| 2 | MINOR | `docs/t027/REMAINING_SLICES.md`（P45 台账条目·验收证据行，_EN 同步） | **台账证据行超出测试实际断言**：写有「`registry.len()==1` / slot owner 与 mirror 一致」为验收证据，但这两项为 adapters 包内私有，交付的并发腿只能以行为等价替代（恰一赢家 / 败者仅两哨兵 / 活体赢家拒绝迟到启动 / 停止释放槽位 / 重发得新 LaunchID）；且测试未重读磁盘 mirror（迟到启动的冲突由进程内 `registry.live()` 判定，mirror 跨重启一致性由 adapters 自有测试覆盖）。 | 台账措辞改为行为等价表述（与 ① 方案一致），或在收尾报告 `p45-e2e-zero-model{,_EN}.md` 显式登记「`registry.len()==1` 与 mirror 一致性以包外等价断言替代，未直接断言」，防止后续审计误以为存在直接证据。 |

（说明：未发现严重度 ≥ MAJOR 的真实缺陷。引擎腿的 T2 失败为 StepPort 故障注入，属协议语义内输入而非状态伪造——真正的自然失败由 console 腿 `localRuntime` 拒绝 `code` 步覆盖。）

## C 交叉核对

1. **父快照链**：`engine.Initialize` 置 `AcceptedParent=baseline`，`runTask` 以 `parent := AcceptedParent` 物化工作区，promotion 完成后回写（engine.go:327-333）；三断言 `baselineHash→acceptedOne→acceptedTwo`、终态 `acceptedThree` 与既有 `TestEngineRunsThreeTasksFourKindsAndBothModes` 的 `wantParents` 及共享常量完全一致。✓
2. **T3 三路互证**：`failTask`→`failChain`（engine.go:564-582）置 FAILED 并返回 cause，`Run` 遇错即停；断言同时覆盖 workspace 调用序、事件流 TaskID、投影键、步骤请求集四路，任一实现「吞错续跑」都会变红，无假绿空间。✓
3. **`len(steps.requests)==3` 的推导**：noop 步在 `runStep` 走 `NOOP_RECORDED` 分支不触 StepPort（engine.go:439-441），故 T1 两请求（code-managed+build-one）+ T2 一请求（code-isolated）后失败，与既有引擎测试全绿时 4 请求的结论自洽。✓
4. **console 自然失败**：`localRuntime.Execute` 对一切非 noop 步返回 `"local runtime cannot execute step kind..."`（runtime_noop.go:35-37），错误串断言钉住来源；`runTaskPostflight` 在无 facts 时直接返回 `steps-completed`（postflight.go:360-407），故 noop 全链不会碰到 fail-closed 的 `RunPostflight` 桩——正向腿绿是协议结果而非跳过。✓
5. **并发腿独占性**：slot 以 `O_CREATE|O_EXCL` 原子抢占（agent_runner_pinned_cli.go:712-761），败者错误全部 `%w` 包装 `ErrAgentRunnerLaunchConflict`；进程内 `registry.live()`（互斥 map）拒绝第二活体。我实测 `-count=5` 连绿（1.06–1.43s/次），`go test -count=1 ./internal/chain/ ./internal/console/` 全绿。✓
6. **边界与编码**：`git diff --name-only -- internal/` 为空（零生产代码改动，`runtime_agent.go:774`/`runtime_supervised.go:335` 未动）；夹具实测 `BOM=False CRLF=False` 1044 bytes，与 `schemas/chain.schema.json` 必填项（root: schemaVersion/chain/tasks/workspace；task: id/review/steps；noop 步: id/kind/reason）吻合；README 仅改 1 行命令路径；台账 +22 行×2 为切片登记；`gofmt -l` 三文件全空。✓

## D 残余风险

1. **无 `-race` 验证**：本环境缺 cgo（`-race requires cgo`），并发腿未过竞态检测。adapters 内部 registry/slot 为互斥锁+原子文件语义，风险低；建议 CI 增补一次 `-race` 跑（本片不阻塞）。
2. **8 路并发版本预检**：8 个 contender 各自 re-exec 测试二进制做 `checkPinnedVersion`（无缓存），慢机（冷缓存/杀软）可能逼近 20s `CheckTimeout`——失败会以非哨兵错误 `ErrAgentRunnerLaunchVersion` **响亮变红**（假红而非假绿），属可接受失败方向。
3. **mirror 持久层未在切片内复验**：跨重启的磁盘 mirror 读取/槽位回收由 adapters 自有测试（pinned_cli_test.go:700-830）覆盖，本片仅证进程内行为等价；台账措辞（发现 2）需如实登记。
4. **引擎级零模型证明的架构前提**：判据「`StepIntents==nil` ⇒ 无 AgentRunner 目标」仅在当前唯一派生通道成立的架构下完备；若未来出现绕过 StepIntentPreparer 的模型通道，本测试不可见（console/preview 腿可部分兜底）。
5. **Windows 进程收割时序**：Stop 后 `guard.ProcessAlive==false` 依赖平台收割语义，与既有 adapters stop 测试同型，本机 5 连绿未现抖动。

## E 可证伪性

- T3 不启动：四路独立观测（workspace 调用、事件流、投影键、请求集）任一路均可单点证伪，回归方向唯一。
- 零模型：引擎腿在 `StepIntents` 被装配、`ExecutionTarget≠Default`、facts 非空或出现 `agent-dispatched` 事件时必红；console 腿在 noop 运行时行为或自然拒绝文案漂移时必红；preview 腿在出现可执行步时经 `RunnableInCLI`/`outcome` 必红。
- 父绑定/事件链：共享哈希常量+生产校验器 `evidence.VerifyEventChain`（哈希链+序列双检），断链即红。
- 并发：0 或 ≥2 赢家、非哨兵败因、活体放行、停止失败、槽位不释放、重发 LaunchID 相同，均单点可证伪。
- **可证伪性缺口（同发现 1）**：五项计数器断言在任何实现下均不可证伪（常量恒等），只能由 `RunnableInCLI` 承担该腿的证伪能力。

`PRE-REVIEW: PASS WITH FIXES`
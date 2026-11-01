# ProofRail 验证报告 — P45-T027B-E2E-ZERO-MODEL

| 字段 | 内容 |
|---|---|
| 切片 | `P45-T027B-E2E-ZERO-MODEL` |
| 类型 | 装配验证类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `eb36992` |
| 日期 | 2026-11-01 |

## 目标与边界

- **目标**：给零模型三任务链补端到端装配证据——T1 正常执行并被接受（`PASSED` + 父快照绑定），T2 自然失败（`FAILED`）后 T3 不启动，全程无模型调用。
- **本片为装配验证片，生产代码零改动**：`git diff --name-only -- internal/` 为空；交付物是三个新增测试文件、一处示例夹具与 README 命令路径的最小修复，以及文档回写。
- **E2E 落点**：主 E2E 在 `internal/chain`（引擎级，只读复用 `threeTaskDefinition` / `testOptions` / `fakeSteps.failTask`）；console 两腿覆盖 `ExecuteNoopRun` 多任务 noop 链与 `localRuntime` 自然拒绝非 noop 步，外加 `examples/no-ai-preview` 预览保持可运行且零 AI；并发腿在 console（同 launcher / 同 runRoot，8 路 Start 恰一赢家）。
- **不改 console 单任务限制**（`runtime_agent.go:774`、`runtime_supervised.go:335`）：多任务装配会与 A6「每 runRoot 唯一 identity mirror」冲突，DR13-FIX 已钉住该排他契约。
- **零模型机械判据**：引擎腿断言 `StepIntents==nil`、`ExecutionTarget==DefaultExecution`、`AgentRunnerFacts==nil`，且事件流不出现 `agentRunnerDispatchedReason`；预览腿断言 `RunnableInCLI==true`、`Outcome==ready`、`BlockingEvidence` 为空且全部步骤 `kind==noop`。
- **计数器口径**：五项预览计数器是 `BuildPreviewReport` 的静态常量，只作 README 承诺的文本锁，不是模型通道哨兵。
- **范围**：`internal/chain/**`、`internal/console/**`、`examples/no-ai-preview/**` 与文档回写。
- **边界**：`internal/adapters/**` 零改动（仅只读调用导出 API）；不改 `docs/CONTRACTS*` 与 `docs/t027/CHECKLIST_v34.md`；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 结论要旨

- 三任务状态序在引擎级与 console 级两侧各自独立成立：T1 `PASSED`、T2 `FAILED`、链 `FAILED`；T3 在 workspace 调用、事件流、投影键与步骤请求集四处均未出现。
- 父快照绑定逐点成立：T1 父为 `baselineHash`，T2 父为 `acceptedOne`，T3 父为 `acceptedTwo`，终态 `AcceptedParent` 为 `acceptedThree`；事件链经 `evidence.VerifyEventChain` 校验，封存记录 3 条。
- 零模型两面成立，且「计数器非哨兵」的口径同时写进测试注释与本报告。
- 并发 Start 在 console 侧以包外行为等价断言钉住 DR13-FIX 强输入：恰一赢家、败者全命中两枚哨兵、赢家 `ProcessAlive`、活体赢家拒绝迟到启动、停止后槽位释放并可重发。
- 首轮运行即暴露一处真实缺陷（示例夹具缺失），按测试夹具修复处置，未触碰任何生产代码。

## 审查链结论

- **① 架构**：`deep-reasoner` + DeepSeek V4 Pro；《零模型三任务 E2E 装配方案》，`DESIGN: DONE`；落点取舍、红测清单、并发断言设计与 4 项未确证项见 `01-arch-design/review.md`。
- **③ 首轮红测**：8 个用例中 7 绿 1 红；红点为预览腿，根因是 `examples/no-ai-preview/proofrail.chain.json` 缺失（README 承诺的示例不可运行），属真实缺陷。
- **② 修复轮**：补该夹具 + README 命令路径最小修正后全绿。
- **③ 二轮**：按 ⑤ MINOR 去恒真断言并做变异自证——把夹具中某步 `kind=noop` 改为 `kind=code` 后预览腿转红，还原后文件 SHA256 与改前一致。
- **④ 门禁**：`TOTAL_FAIL=0`，读数见 §门禁读数。
- **⑤ 预审**：`PRE-REVIEW: PASS WITH FIXES`（2 MINOR：计数器断言同义反复、台账证据行超出实际断言）。
- **⑥ 扫描**：`INDEPENDENT SCAN: PASS`，Section B 判空，Section D 列 4 类未覆盖面；本片额度 1/3。
- **⑦ 终审**：`FINAL REVIEW: PASS`，另附 1 项 MINOR 与 6 条放行条件，逐条落地见 §偏差与登记；本片额度 1/3。
- **R / D**：R = 0、D = 0（不继承上片）。

## 六段扫描转录

来源 `docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/06-scan/review.md`；围栏内按 A→B→C→D→E 顺序转录 ⑥ 正文，围栏外披露三处归一。

```sixstage
## Section A 扫描范围与输入
范围：只读核查 internal/chain/engine_zero_model_e2e_test.go、internal/console/runtime_zero_model_e2e_test.go、internal/console/runtime_agent_concurrent_e2e_test.go、examples/no-ai-preview/proofrail.chain.json 与 examples/no-ai-preview/README.md，聚焦三任务状态序、T3 未启动、零模型判据与并发 Start 互斥。

## Section B 发现清单
| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
| --- | --- | --- | --- | --- |
|  |  |  |  |  |

## Section C 与 ①/⑤ 的一致性核对
1. T3 不启动的多路互证成立，且没有假绿空间。
   - 在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 中，T2 失败后会检查：`workspaces.calls` 只有 T1/T2，`projection.TaskStates` 不含 T3，`store.events` 中不出现 T3 任务实体，`steps.requests` 中不出现 T3，并且 `len(steps.requests)==3` 这个断言强制卡住“额外发起”。
   - 在 [internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 中，控制台对 T2 用 `kind="code"`、`Mode=ManagedChangeSet` 的真实拒绝路径断言了 `localRuntime` 自然失败，且 `task-three` 不在 `taskIDs` 中。
   - 这不是单一字段判读，而是“任务启动序/事件流/投影键/步骤请求集”四重交叉，真实覆盖了“T3 never started”。

2. 零模型证明机械、且没有把计数器误当哨兵。
   - 引擎腿在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 里直接断言：`options.StepIntents == nil`；并在执行时检查 `request.ExecutionTarget == DefaultExecution`，`request.AgentRunnerFacts == nil`，并且事件中不存在 `agentRunnerDispatchedReason`。
   - 预览腿在 [internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 中，计数器依然保留为 0，但测试明确说明它只是“README 运营承诺的静态锁定”，真正机械判断是 `RunnableInCLI`、`Outcome == ready`、`BlockingEvidence` 为空、且所有 `preview.Steps` 都是 `noop`。这符合“可证伪 + 计数器非哨兵”的要求。

3. 父快照绑定与事件链语义一致，且无越界。
   - [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 断言：T1 的父是 `baselineHash`，T2 的父是 `acceptedOne`，T3 的父是 `acceptedTwo`，终态 `AcceptedParent == acceptedThree`；这与真实链语义一致。
   - 事件链断言也直接调用了 [internal/evidence/event.go](internal/evidence/event.go) 的 `VerifyEventChain`，校验 `Sequence` 单调、`LastEventHash` 与最终事件 hash 一致，并确认 `seal` 数量 3。该证明与真实 engine state machine 一致，不是“测试造假”。

4. 并发腿是同 launcher/runRoot 的真正排他场景，且无墙钟判据。
   - [internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go) 复现的是同一 `launcher` 与同一 `runRoot` 下 8 个并发 start；赢家唯一，失败者只可看到 `ErrAgentRunnerLaunchConflict` / `ErrAgentRunnerProcessDuplicate`，并验证 `ProcessAlive`、slot owner 与 mirror 一致。
   - 这与 DR13-FIX 强输入一致，且没有使用 10s freshness 作为排他依据；它使用“状态/存在性/完成序”而非墙钟，是合理的。

5. 夹具和 README 与严格解码一致，且最小必要变更。
   - [examples/no-ai-preview/proofrail.chain.json](examples/no-ai-preview/proofrail.chain.json) 结构符合 [internal/console/config.go](internal/console/config.go) 的 `LoadChainConfig` 与 [schemas/chain.schema.json](schemas/chain.schema.json) 的 schema：`schemaVersion`、`chain`、`tasks`、`workspace` 全部合法，且文件大小 1044 bytes，起始 `0x7B 0x0A 0x20` 表明无 BOM，符合要求。
   - [examples/no-ai-preview/README.md](examples/no-ai-preview/README.md) 的修正是命令路径修正，最小必要，且与预览示例的行为一致。

6. 边界/卫生与证据包自洽性。
   - 实际只读边界被尊重：新测试文件是 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go)、[internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go)、[internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go)，生产代码未被改动；`git diff --name-only -- internal/` 在本轮检查中没有输出 tracked diff，只有这 3 个新测试文件分布在 internal tree，符合“零改动生产代码”的要求。
   - 证据包的 [docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md](docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md) 明确写着“步骤 0 骨架落盘；后续载荷按轮追加”，所以它当前是合法的 start-of-slice skeleton，而不是最终闭环包；该状态在文件中自说明，自洽。
   - Fresh verification evidence: `go test -count=1 ./...` returned all package lines as `ok` with no `FAIL` lines, including the internal/chain and internal/console packages.

## Section D 未覆盖面与残余风险
- 这次只读复核未对真实 VS Code UI/浏览器路径做端到端点击验证；它依赖代码和事件日志的静态断言。
- 未对 Windows 之外的 OS/locale 运行环境做差异验证；并发 snapshot 与 ProcessAlive 的行为可能受平台差异影响。
- 证据包未进入最终 06/07/08 载荷阶段，因此未覆盖终审与 CI 归档闭环。
- 例子中的 CLI 命令路径已被核验为合法 JSON 例子，但未运行 `go run ./cmd/prfrail preview --chain ...` 的真实命令行交互输出。

## Section E 结论要旨
- 只要把 [examples/no-ai-preview/proofrail.chain.json](examples/no-ai-preview/proofrail.chain.json) 中任一 step 的 `kind` 改成 `code`，[internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 会立即让 `preview.Steps` 或 `BlockingEvidence` 断言失败。
- 只要让 T2 在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 里继续推进到 T3 事件，`len(workspaces.calls)` 或 `task-three` 事件断言会立刻变红。
- 只要把并发测试里的“唯一赢家”逻辑去掉，或让 loser 观察到非 `Conflict`/`Duplicate` 错误，[internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go) 会直接失败。

INDEPENDENT SCAN: PASS
```

- 归一披露一：⑥ 的段标题 `## A 结论` / `## B 缺陷` / `## C 核验结论` / `## D 未覆盖面` / `## E 可证伪性` 在围栏内改写为本模板的段名，五段正文逐字未改。
- 归一披露二：Section A 的范围行取自 ⑥ Section C 第 6 条所核查的文件集合，归档原文未单列范围行。
- 归一披露三：⑥ 末尾的 `§A: INDEPENDENT SCAN: PASS` 归一为围栏末行的 `INDEPENDENT SCAN: PASS`，Section B 的空表头与空行逐字保留。

## 项闭合

严重度列取审查报告原文口径，不加粗；⑦ 结论闭集为 `FINAL REVIEW: PASS`（0 新增可执行缺陷）。

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 1 | 五项预览计数器断言是同义反复，任何实现改动都不会红 | MINOR | `internal/console/runtime_zero_model_e2e_test.go` | 已闭合：测试注释声明其为 README 承诺的静态文本锁、非机械判据；机械判据改为 `RunnableInCLI` / `Outcome` / `BlockingEvidence` / 步骤 `kind` |
| ⑤ r1 2 | 台账证据行写有 `registry.len()==1` 与 slot owner/mirror 一致，超出测试实际断言 | MINOR | `docs/t027/REMAINING_SLICES.md`、`_EN` | 已闭合：验收证据行改为包外行为等价表述，并在本报告 D7 显式登记未直接断言项 |
| ⑦ B1 | 断言失败提示中的 README 行号与命令片段过时 | MINOR | `internal/console/runtime_zero_model_e2e_test.go` | 已闭合：文案改为 `README.md:10` 与 `--chain .\examples\no-ai-preview\proofrail.chain.json` |

## 交付

- `internal/chain/engine_zero_model_e2e_test.go`：引擎级主 E2E（T1 / T2 / T3 三态 + 父快照绑定 + 事件链封存单调 + 零 AgentRunner 派发）。
- `internal/console/runtime_zero_model_e2e_test.go`：console 三 noop 任务链正腿、T2 自然失败后停链腿、示例预览保持可运行腿。
- `internal/console/runtime_agent_concurrent_e2e_test.go`：同 launcher / 同 runRoot 的 8 路 Start 并发腿，附自执行 helper。
- 测试 8 个用例（`internal/chain` 4 + `internal/console` 4），另附 1 个自执行 helper `TestZeroModelAgentCLIHelper`。
- `examples/no-ai-preview/proofrail.chain.json`：README 承诺却缺失的两任务纯 noop 夹具，属测试夹具修复，不参与默认运行路径。
- `examples/no-ai-preview/README.md`：命令路径最小修正（1 行）。
- 文档回写：`docs/ARCHITECTURE{,_EN}.md`（三任务 E2E 装配路径与并发断言要点）、`docs/DEV_PLAN{,_EN}.md`（T027 里程碑 E2E 状态）、`docs/t027/REMAINING_SLICES{,_EN}.md`（P45 汇总行与详情块）。

## 偏差与登记

| 编号 | 内容 | 来源 | 判定 |
|---|---|---|---|
| D1 | ⑦ MINOR：断言失败文案的 README 行号与命令片段漂移 | ⑦ 放行条件 1 | 已闭合：文案指向 `README.md:10` 与实际命令路径（见 §项闭合 第 3 行） |
| D2 | ⑤ 的 2 项 MINOR 及其处置 | ⑦ 放行条件 2 | 已闭合：计数器断言改为声明式文本锁、台账证据行改为行为等价表述（见 §项闭合 第 1、2 行） |
| D3 | ⑥ 的 4 项未覆盖面：真实 UI 端到端、跨 OS 与 locale 差异、证据包终态闭环、真实 preview 命令执行证据 | ⑦ 放行条件 3 | 接受并登记：前三项列入 §已知限制 与 §后续片候选登记；第四项已由 §门禁读数 的手工执行补证 |
| D4 | ① 的 4 项未确证项：`engine.go` 与 `launcher` 行号须以 `grep -n` 复核、`localRuntime` 拒绝文案若变则红测需同步、`step-failed` 为有意契约、预览腿依赖示例配置不被改；另含「不改 console 单任务限制及 A6 mirror 冲突」边界 | ⑦ 放行条件 4 | 接受并登记：行号已在实现期复核；其余三项为本片既有契约，未改动即无需同步 |
| D5 | 本片为装配验证片，生产代码零改动 | ⑦ 放行条件 5 | 已声明：`git diff --name-only -- internal/` 为空输出，新增内容仅 3 个 `_test.go`（§目标与边界 第 2 条） |
| D6 | `examples/no-ai-preview` 为测试夹具修复，不参与默认运行路径 | ⑦ 放行条件 6 | 已声明：夹具为纯 noop、无执行副作用，仅由测试与 README 显式引用 |
| D7 | `registry.len()==1` 与 slot-owner / mirror 一致性为 `internal/adapters` 包内私有态，本片未直接断言 | ⑤ r1 2 / ⑥ C 第 4 条 | 接受：以包外行为等价断言替代（恰一赢家、败者两枚哨兵、活体赢家拒绝迟到启动、停止后释放并可重发），跨重启 mirror 由 adapters 自有测试覆盖 |
| D8 | 五项预览计数器为 `BuildPreviewReport` 静态常量 | ⑤ r1 1 / ⑥ C 第 2 条 | 接受：仅作 README 承诺的文本锁，不得读作模型通道哨兵；模型零调用由引擎腿四项断言与预览腿派生量承担 |
| D9 | 保留 console 单任务限制，未放开 `runtime_agent.go:774` 与 `runtime_supervised.go:335` | ① A 节 / ⑦ 放行条件 4 | 接受：多任务装配与 A6「每 runRoot 唯一 identity mirror」冲突，DR13-FIX 已钉住该排他契约；放开会红既有拒绝用例并扰动 DR-9 敏感面 |
| D10 | ⑥ 围栏引入 2 个 HEAD 语料未收录汉字，按仓库既有机制登记到 CJK 新字白名单 | ⑧ 门禁（G4b） | 接受并登记：逐字转录不可改写，白名单条目逐字附理由与上下文；本项为本片唯一越出「文档回写」清单的写入，如实登记 |
| D11 | `go test ./...` 包数口径与起跑简报表述不一致 | ⑧ 实测 | 以实测为准：15 个测试包全 `ok`、0 `FAIL`，另有 4 个包无测试文件；本报告不写未实测的包数 |

## 模型与计量

| 环节 | 载体 | 模型 |
|---|---|---|
| ① 架构 | `deep-reasoner` | DeepSeek V4 Pro |
| ② 实现 | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ 测试 | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ 预审 | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ 扫描 | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ 终审 | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ 文档 | `prfrail-documenter` | DeepSeek V4.1 Flash |
| 批量辅助 | `deep-reasoner` | GLM-5.3-Flash（实派 2 次：装配入口普查 / 接口与端口契约普查） |

- **主控模型**：`DeepSeek V4.1 Flash`（High 档）；①⑤ 按白名单仍由 V4 Pro 承担，与主控选型是两件事。

## 门禁读数

- `node tools/gates/gate.js --all --skip-base`（tree 作用域）：`TOTAL_FAIL=0`；G3(a) / G3(b) PASS、G4b `unknown=[]`、G6-4 无命中、G7-a 至 G7-f 全 PASS、G8-a 至 G8-d 全 PASS。
- `gofmt -l internal/chain internal/console`：空输出。
- `go vet ./...`：无输出。
- `go test -count=1 ./...`：15 个测试包全 `ok`、0 `FAIL`；另有 4 个包无测试文件（`cmd/prfrail`、`cmd/prfrail-release`、`tools/agent-probe/supervisor`、`tools/b4-e2e`）。
- `go test -count=1 -run 'TestZeroModel' ./internal/chain/ ./internal/console/ -v`：8 个用例全 `PASS`，自执行 helper 亦 `PASS`。
- 手工预览实测：`go run ./cmd/prfrail preview --chain examples/no-ai-preview/proofrail.chain.json --json` → `runnableInCli=true`、`preview.outcome=ready`、`blockingEvidence` 为空、两步 `kind=noop`、五项计数器全 0、`unknowns` 3 条。

## 提交与 CI

提交与 CI 读数在 ⑨ 之后由本节追加一次，追加时同步更新 CN/EN。

## 后续片候选登记

- 提交 `4273d6e` 已推送 origin；CI `41315652405` **首跑双腿全绿**；ubuntu 腿 `ok chain 0.564s` / `ok console 2.856s` / `ok adapters 12.291s` / `ok gates 0.159s` + Race 腿 `ok adapters 89.514s` / `ok chain 2.521s`，windows 腿 `ok chain 5.011s` / `ok console 20.796s` / `ok adapters 34.297s` / `ok gates 0.786s`；两腿 `SELFTEST 697/697`、`GATE REPORT scope=ci changed=23`、`TOTAL_FAIL=0`。
- 跨 OS 与 locale 的并发差异验证：`ProcessAlive` 与快照语义在 Windows 之外的实测。
- `-race` 腿扩包至 `internal/console`：本机缺 cgo，并发腿未过竞态检测。
- 真实 UI 端到端：预览与运行在编辑器内的路径尚未与本片测试对齐。

## 已知限制

- **⑥ 未覆盖面一**：真实 UI 与浏览器路径未做端到端点击验证，判定依赖代码与事件日志的静态断言。
- **⑥ 未覆盖面二**：未对 Windows 之外的 OS 与 locale 做差异验证；并发快照与 `ProcessAlive` 可能受平台影响，以 ubuntu CI 腿为最终判据。
- **⑥ 未覆盖面三**：证据包在 ⑧ 时点尚未进入 06/07/08 的最终载荷阶段，终审与 CI 归档闭环由 ⑨ 补齐。
- **⑥ 未覆盖面四**：真实 preview 命令执行证据在本报告 §门禁读数 以手工执行补证，未纳入自动化测试。
- **⑤ 残余**：`registry.len()==1` 与 slot-owner / mirror 一致性无直接断言（D7）；跨重启镜像由 adapters 自有测试覆盖。
- **⑦ MINOR 已闭合**：断言失败文案漂移按 §项闭合 第 3 行处置，不影响运行语义与门禁强度。
- **D9 摘要**：console 单任务限制保持不动，多任务装配与 A6 排他契约互斥，属本片边界外。

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/p45-e2e-zero-model.md	repo
docs/validation/p45-e2e-zero-model_EN.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/SHA256SUMS.txt	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/00-briefing.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/freedom-list.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/01-arch-design/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/03-tests/red-baseline.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/05-prereview/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/06-scan/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/07-final/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/08-ci/ci-evidence.txt	repo
examples/no-ai-preview/README.md	repo
examples/no-ai-preview/proofrail.chain.json	repo
internal/chain/engine_zero_model_e2e_test.go	repo
internal/console/runtime_zero_model_e2e_test.go	repo
internal/console/runtime_agent_concurrent_e2e_test.go	repo
tools/gates/cjk-newwords.txt	repo
```

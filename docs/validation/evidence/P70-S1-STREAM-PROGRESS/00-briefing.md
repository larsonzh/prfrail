# P70 起跑消息原文（00-briefing）

切片：S1-STREAM-PROGRESS（P70 · 层 2 治理侧 · 在 S1 内补 `run --progress` 轻量流式开关 · 实现类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】

1. 请先完整阅读 docs/t027/CHECKLIST_v56.md（当前快照）与 docs/DELIVERY_DIRECTIVE.md（v1.67，含 §13）。

2. 核对 v56 §一 的 8 个锚点。**注：v56 §一 记录的 HEAD / origin 是 v56 的生成时点锚点（`759c934`）；v56 的三个提交（主 `d06f229` + `eb8568a` + 回填 `ee7f260`）落地后 HEAD 已按预期前进**。核对时按以下现时值判定：
   - HEAD = **`ee7f260`**；origin/main = **`ee7f260`**
   - **若 HEAD 不是 `ee7f260` 且原因不是 v56 的三个提交 ⇒ 报实况（不按 759c934 报不符，也不静默接受未知前进）**
   - 准则 **v1.67**；SELFTEST **737/737**；三作用域 `TOTAL_FAIL=0`
   - v53 = 292 / v54 = 180 / v55 = 182（LF，CN/EN 对称）
   - **若 ≥2 项不符（除上述 HEAD/origin 已声明前进外）⇒ 停止并报实况。**

3. **主控模型声明**：本片为 S 档实现类（改 `internal/console/**` 加轻量开关），按 v1.67 §2.2 默认 **`DeepSeek V4.1 Flash (deepseek)`**。**不登记 ADR**（§13.1.8 字面未触发；同类先例 P50/P51 console 层主控用 Flash）。回执首行标 [自执行]。

4. **记忆唤醒（红线复述）**：
   - 本片受 v1.67 §11.1（历史不改写）、§7.5（S 档 ⑦ 1+1）、§7.2.1（⑥ 可跳过三条件）、§13.1 第 7 条（批量辅助 GLM-5.3-Flash ≥2 次）、§13.2.21（纪律 V 大任务拆分）、§13.4（门禁坑 1–30）、§13.5（快照结构模板）、§5.1（临时件位置闭集）约束。
   - **E-1 流程纪律（本片执行）**：⑧c 后到 ⑩ 回执之间，须**二次回读** `tmp` 五条读数，防编辑器缓冲二次落盘；若发生，如实登记不记 D（位置合法 + 工具行为 + 已零残留三要件成立时）。
   - 单次输出 ≤ 10,000 字符；不自行 commit/push；不替用户裁定；不启动任何新切片。
   - **临时件与门禁捕获文件一律落 `tmp/P70-S1-STREAM-PROGRESS/`**，不落 `tmp/` 根。

【〇 本片环境声明与流水线裁剪】

1. **环境**：Windows 本机；**不开启 Ubuntu VM**；**不开启网络代理**；**不启用 Sol 双打**（S 档确定性实现，无不可测语义）。Unix 侧「零影响」用本机 `GOOS=linux go build/vet` 验证 + CI ubuntu 腿兜底。

2. **流水线（实现类 S 档）**：① 事件面核对与开关设计 → ② 开关与打印实现（按纪律 V 拆小步）→ ③ 测试 → ④ 门禁 → ⑤ 预审（1+1）→ ⑥ 独立扫描（按 §7.2.1 三条件判定，须 ⑦ 确认）→ ⑦ 终审（1+1，S 档）→ ⑧ 文档装配 → ⑨ 原生验证 → ⑧c → ⑩。

3. **⑥ 判定提示**：P70 改 `internal/console/**` 生产代码语义（新增开关 + 输出路径）⇒ §7.2.1 三条件中 ①②③ 均不成立 ⇒ **⑥ 必**（不可跳过）。

4. **⑦ 盲审判定（须 ① 阶段显式判定并在回执记录理由）**：P70 不触所有权/停机/身份语义 ⇒ §7.7 必抽不成立。**按比例抽查须先澄清一个关键事实**：

   - **P67 的盲审算不算 §7.7 按比例抽查的「已抽片」？** §12.4 表尾清账登记写「P67 盲审 = 主控主动加 1 次（不计 §7.7 必抽 +1）」——「不计必抽 +1」语义模糊：指「不占必抽名额」（可算按比例抽查的已抽片），还是「不是 §7.7 任何一类」（不算已抽片）？
   - **主控须在 ① 阶段读 P67 报告 §7 与 §12.4 表尾「按比例抽查记录」块，给出明确判定**。
   - **若 P67 盲审算已抽片**：窗口内 P67 已抽，P68 / P69 / P70 为 3 片未抽，**未满 4 片 ⇒ P70 不抽**；回执须显式记录「P67 已抽，本窗口剩余 3 片未满 4，P70 不抽」。
   - **若 P67 盲审不算已抽片**：P67 / P68 / P69 / P70 为窗口内第 1/2/3/4 片，**满 4 片 ⇒ 应抽 P70**（P67 / P68 / P69 的 ⑦ 均已闭环不可补抽；P70 是唯一仍可抽的片）；回执须显式记录判定依据。
   - **无论哪种判定，回执必须显式记录理由，不得静默跳过**。
   - **顺带核实（供 P71 参考，不属 P70 判定）**：P68 是评估类（无实现对象），它是否计入 §7.7 的「已收尾的正式切片」窗口？如不计入，P71 起跑时窗口计数会变。P70 的回执一并给出核实结果。

【一、授权边界（流水线角色与模型）】

- 主控：`DeepSeek V4.1 Flash (deepseek)`（默认，不登记 ADR；不得中途升级，若 ① 判定须升级按 §13.1.8 停下提示用户）
- ① 事件面核对与开关设计：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 开关与打印实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（S 档 1 终审 + 1 复审）
- ⑦ 盲审：按 §7.7 判定（见【〇】第 4 点）
- ⑧ 文档装配：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（Flash）执行
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`（限定名）；**严格串行**（R5.1）

【二、批量辅助（强制 ≥2 次）】

- 载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单内批量辅助首选，无需临时授权）
- 实派次数：**≥ 2 次**（§13.1.7 与 §13.3.25；实派不足同轮补派）
- 任务限定：
  ① **`events/state-events.jsonl` 写入点与字段枚举**：从 `internal/console/runtime.go`（路径定义）与 `internal/chain/engine.go`（写入点）枚举（file:line）；列出所有 `eventType` 值（闭集）及其对应可打印字段（taskId/stepId/state/occurredAt 等）；**标注每个写入点位于 `chain.Run` 的哪一跳（任务循环内 / 任务循环外 / 结束时）**——这决定 `--progress` 能否在运行中看到该事件。**不要判定「是否可被外部读取」**（那是 fsync/buffering 的运行时问题，静态代码不可判定）。
  ② **AgentRunner 事件的写入点与字段枚举**：从 `internal/adapters/agent_runner_event.go` 与 `internal/adapters/agent_runner_run.go` 枚举（file:line）；列出 `AgentRunnerEvent` 的全部 `EventType`（`session-created` / `output` / `tool` / `usage` / `diagnostic` / `operator-action-required` / `process`）；**标注 `events.jsonl` 的写入点位于 AgentRunner 运行的哪一跳**（启动时 / 工具调用中 / 结束时）。
  ③ 不允许派「排序与 Top-N 择优」与结论性判断（§2.2 排除项）
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三、五件评估任务（① 阶段）】

## 3.1 现状核对（只列事实，不判断）

- `run` 的当前输出：`internal/console/cli.go` 的 `executeRun`（file:line）与 `writeRunSummary`（结束时 6 行摘要的 file:line）
- `state-events.jsonl` 路径定义：`internal/console/runtime.go` 的 `runEventLogPath`（file:line）
- 事件写入点：`internal/chain/engine.go` 的事件持久化路径（file:line）
- AgentRunner 事件路径：`internal/adapters/agent_runner_event.go` 与 `internal/adapters/agent_runner_run.go`（file:line）
- 现有 `--json` 的实现：`writeJSONResponse`（file:line）
- **P68 技术前提结论**（`docs/validation/remaining-rulings{,_EN}.md` §5 与 `01-arch/r07-refined.md`）：`run` 无 per-task/per-step 输出、无 `--verbose`/`--progress`、无心跳；`state-events.jsonl` 是结构化唯一全序源但**无订阅 API**

## 3.2 开关设计（四件齐全）

- **数据源与可读性**：`state-events.jsonl` 与 `agent-runner-replay/runs/<requestId>/events.jsonl` 在 `run` 运行中是否可被外部读取？**若两者都仅在结束时可见 ⇒ 须改 chain 层事件写入点改为「写完即刷盘」或「同步打印」——这会触停升级边界（见 3.3）**
- **实现形态三选一**（必须给判定并说明理由）：
  - **形态 A（推荐）：事件写入点主动回调**——在 `internal/console` 的事件写入封装处，写完后若 `--progress` 开启，则同步打印摘要行；不改 chain 层；但需要 console 层有事件写入点（若事件只由 `chain.Engine` 写，则形态 A 不成立）
  - **形态 B：后台 goroutine 跟踪文件**——启动时 fork 一个 reader goroutine tail `state-events.jsonl`，解析并打印；但 P68 前提显示「无订阅 API」⇒ 只能 poll；poll 间隔与 CPU 成本需评估
  - **形态 C：`--progress` 只覆盖 AgentRunner 事件**（`events.jsonl` 由 adapters 写），不覆盖 chain 状态事件——最小实现，但可用性打折（看不到 task/step 起止）
- **与 `--json` 的关系**：互斥 / 共存 / 覆盖？（`--json` 是最终结构化输出；`--progress` 是运行中摘要——建议共存，`--json` 在结束时输出，`--progress` 在运行中输出）
- **默认关闭**：不改现有默认行为
- **打印粒度**（来自 v56 §四 P70 目标）：逐 task/step 起止 + 门禁命令起止与退出码 + AI 工具调用摘要行
- **反例清单**：至少 3 条——① `--progress` 未开启时输出与现状完全一致；② `--progress` 开启时 `--json` 的最终输出不变；③ `state-events.jsonl` 不可读时 `--progress` 降级为「无进展」（不报错、不中止）
- 结论行：`DESIGN: DONE`（或 `DESIGN: NEEDS_CLARIFICATION` + 待澄清项），落 `tmp/P70-S1-STREAM-PROGRESS/01-arch/design-r1.md`

## 3.3 停升级边界（必须遵守）

- ① 判定 `--progress` 必须改 `internal/chain/**` 或 `internal/adapters/**` 的事件订阅机制 ⇒ **停下升级**（本片边界只改 `internal/console/**`）
- ① 判定 `state-events.jsonl` 与 AgentRunner 事件在 `run` 运行中都不可读，且无法在 console 层解决 ⇒ **停下升级**（可能需要 chain 层改事件刷盘策略）
- ① 判定需改 `CONTRACTS{,_EN}.md` 的语义（不只是措辞）⇒ 停下升级
- ① 判定需改 schema / fixtures / workflows / `go.mod` ⇒ 停下升级
- ① 判定 `--progress` 与 `--json` 的语义冲突无法调和 ⇒ 停下升级
- ⑤ / ⑦ 报 Medium+ ⇒ 修复后原位重跑该环
- 触任一情形 ⇒ 停下升级用户裁决

【四、执行步骤（纪律 V：② 分小步串行）】

1. **① 事件面核对与开关设计**（V4 Pro）：按【三】出五件，落 `01-arch/design-r1.md`。若触 3.3 停升级边界 ⇒ 停下升级。
2. **② 开关与打印实现**（Flash，按纪律 V 拆步，每步独立核验）：
   - ②a：`internal/console/cli.go` 的 `executeRun` 加 `--progress` flag 解析（file:line 级改动，跑聚焦测试）
   - ②b：按 ① 选定的形态（A / B / C）实现进度打印（file:line 级改动，跑聚焦测试）
   - ②c：补测试（见 ③）
   - **每步完成即跑 `go build ./...` + `go vet ./...` + 聚焦测试**；任一步红即修，不进入下一步
   - 若某步发现需要改 `internal/chain/**` 或 `internal/adapters/**` ⇒ 停下升级（边界外）
3. **③ 测试**（Flash）：
   - 正例 1：`--progress` 未开启时，输出与现状**逐字节一致**（回归断言）
   - 正例 2：`--progress` 开启时，逐 task/step 行可见（至少 3 行断言：task 起始、step 起始、step 完成）
   - 正例 3：`--progress` + `--json` 共存时，`--json` 最终输出不变
   - 反例 1：`state-events.jsonl` 不可读时 `--progress` 不报错、不中止，只降级为「无进展」
   - 反例 2：`--progress` 与 `--json` 互斥时（若 ① 判定为互斥），两者同传返回用法错误
   - 每条测试绑定「删掉哪个守卫必须变红」的变异检验目标
4. **④ 门禁**：`go build ./...` + `go vet ./...` + `go test ./...` + `gate.js --all --scope=index / tree / ci` + `gate.js --selftest` 全绿；`GOOS=linux go build ./...` + `GOOS=linux go vet ./...` 干净（Unix 零影响本机证据）
5. **⑤ 预审**（V4 Pro，1+1）：核对开关实现是否与 design-r1 一致、反例是否成立、有无越界改 `internal/chain/**` / `internal/adapters/**` 的判定本体；报 Medium+ ⇒ 修复后原位重跑 ⑤
6. **⑥ 独立扫描**（MAI，3 次额度）：按 §7.2.1 三条件判定（本片改 .go 语义 ⇒ ⑥ 必）；报 Medium+ ⇒ 修复后重跑 ⑥
7. **⑦ 终审**（Codex，1+1）+ **盲审**（若 ① 判定须抽）：按 §7.3 只给切片定义、diff、契约原句、只读约束；首轮不附任何清单；复核 `--progress` 是否改变了默认行为、`--json` 语义不变、Unix 零影响是否有证据
8. **⑧ 文档装配**（documenter）：按【五】清单逐项，CN/EN 严格对称
9. **⑨ 原生验证**：三作用域 + SELFTEST + `go build/vet/test` 全绿（读数为 ⑩ 回执的前置）
10. **⑧c 收尾**：`tmp/P70-S1-STREAM-PROGRESS/` 逐件精确清理，回读五条读数；**⑧c 后到 ⑩ 回执前二次回读一次**（E-1 防再犯）
11. **⑩ 停点**：停在 ⑩ 等授权；不自行 commit/push

【五、文档同步更新（本片范围；CN/EN 严格对称）】

**必改（6 项）**：

1. **新建报告对** `docs/validation/s1-stream-progress{,_EN}.md`（核心交付物：§1 执行摘要 / §2 开关设计 / §3 数据源与形态选择 / §4 反例与变异检验 / §5 审查链结论 / §6 门禁读数 / §7 已知边界 / §8 明确未执行 / `artifacts` 围栏 / `sixstage` 转录）
2. **新建证据包** `docs/validation/evidence/P70-S1-STREAM-PROGRESS/`（`00-briefing.md` / `01-arch/design-r1.md` / `02-batch/` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）
3. **`docs/OPERATIONS{,_EN}.md`** —— 用户可见的 CLI 命令新增了开关，操作手册必须有：在「当前可执行」段或等价处，补 `run --progress` 的用法示例与语义说明；不改其它内容
4. **`docs/ARCHITECTURE{,_EN}.md`** —— 在 §3 数据流与事务的「console 多任务串行装配」段之后（或事件/证据相关段），**加法式**补一段 `--progress` 的事件消费路径与形态选择（A/B/C 中判定的那个）；不改既有描述
5. **`docs/DEV_PLAN{,_EN}.md`** —— T027-B 里程碑段（`:176`–`:185` 附近）**加法式**加 P70 一行（报告路径）
6. **`docs/t027/REMAINING_SLICES{,_EN}.md`** —— P70 段状态回写（🔄 → ✅）+ 完成记录行（报告路径，按 §11.2 不写 hash/run）

**需评估（2 项，① 阶段给出判定）**：

7. **`docs/TEST_STRATEGY{,_EN}.md`** —— ① 阶段核实 `AT-12`（`REQ-009/010/015/016/027`）的断言是否需要加一条关于 `--progress` 的条目（`AT-12` 现有断言含 "键盘/无色/窄终端/--json"，`--progress` 是新增输出形态）；**有需要 ⇒ 加法式补**；**无需要 ⇒ 标注「不适用 + 理由（AT-12 现有断言未涵盖 --progress，但 --progress 不改变 --json 稳定性、不改变窄终端行为，故 AT-12 边界不变）」**
8. **`docs/validation/t027b-e2e-wiring{,_EN}.md`** —— ① 阶段核实 `--progress` 的打印粒度是否需要包含 P69 的 readiness/permit 判定结果（若 `--progress` 打印「逐 task 开始」行时顺手打 readiness 通过/失败，则 P69 接线点需被引用）；**有需要 ⇒ 加法式加指针（不改原行）**；**无需要 ⇒ 标注「不适用 + 理由」**

**不适用（8 项，理由见各行）**：

9. **`docs/t027/CHECKLIST_v56{,_EN}.md`** —— 快照冻结，v57 生成时继承
10. **`docs/DELIVERY_DIRECTIVE{,_EN}.md`** —— 业务切片，非治理片；本片不改准则正文；P70 主控默认 Flash，不登记 ADR
11. **`docs/CONTRACTS{,_EN}.md`** —— `--progress` 是 CLI 输出开关，不触 wire / Schema / 状态机语义；若 ① 判定必须改契约 ⇒ 停下升级，本片不落
12. **`docs/SECURITY{,_EN}.md`** —— 不改安全边界
13. **`docs/PRODUCT_REQUIREMENTS{,_EN}.md`** —— 不触 REQ；`--progress` 是 REQ-027 的辅助能力，不改变 REQ 语义
14. **`docs/S1_SUPPORT_MATRIX{,_EN}.md`** —— 不触平台清单
15. **`docs/validation/s1-exit{,_EN}.md`** —— P70 不直接闭合 S1 Exit 八门
16. **`docs/validation/remaining-rulings{,_EN}.md`** —— P70 不触 P68 裁定；`--progress` 是 R-07 附条件 1 的**落地**，但 P68 报告是裁定记录、不回写

> 第 9–16 项在 ⑩ 回执中逐项标注「不适用 + 理由」；第 7、8 项须给评估结论（适用或不适用 + 理由）。不得省略。

【六、⑩ 回执要求】

1. 派单表（`agentName` 逐字机械读取 + 机械读取脚本 stdout 原文；机械读取失败须如实标注「未机械校验」）
2. ① 五件摘要（现状核对 + 开关设计 + 数据源与形态选择 + 与 --json 关系 + 反例清单，每项一句话 + 落点）；**⑦ 盲审判定及其理由**（见【〇】第 4 点：必须先核实 P67 盲审算不算已抽片，再判 P70 抽不抽）；**P68 是否计入 §7.7 窗口的核实结果**（供 P71 参考）
3. ② 分步（②a–②c）每步的 `git diff --stat` 与聚焦测试读数
4. ③ 测试结果（正例 3 + 反例 2 + 变异检验逐条红/绿）
5. ⑤ / ⑥ / ⑦ / 盲审（若抽）结论摘要（⑦ 计量 = 终审 N + 复审 M + 盲审 K）
6. 三作用域 + SELFTEST + `go build/vet/test` 原始行（含 `GOOS=linux` 两条）
7. 文档同步 `git diff --stat`（逐项覆盖【五】16 项；必改 6 项须有 diff；需评估 2 项须给评估结论；不适用 8 项须列理由）
8. ⑧c 后五条终态读数（v1.67 §5.1「⑧c 后取证窗口」子款）+ **⑧c 后二次回读一次**（E-1）
9. 明确未执行事项：未改准则正文 / 判据本体 / `CONTRACTS` 语义 / schema / fixtures / CI / v50–v56 / `internal/chain/**` / `internal/adapters/**` / `internal/gates/**` 判定本体；未替用户裁定；未 commit/push；未启动新切片
10. 本片自身越界登记：§1.7.2 三条 bullet 越界如实登记；未发生写「无」
11. 待授权项：提交授权；push 授权
12. 给 v57 的事实包（v57 生成时用：P70 闭环 hash + SELFTEST 读数 + P70 `--progress` 的形态选择与数据源决策 + ⑦ 盲审判定的最终口径）

【七、执行节奏与红线】

1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束
3. 临时件一律 `tmp/P70-S1-STREAM-PROGRESS/`；只增不删；⑧c 统一清理；**CI 日志与门禁捕获文件也落这里，不落 `tmp/` 根**
4. **不改准则正文**：`DELIVERY_DIRECTIVE{,_EN}.md` 完全不改
5. **不改产品代码边界外**：本片代码只改 `internal/console/**`；不改 `internal/chain/**`、`internal/adapters/**`、`internal/gates/**`、`internal/tickets/**`、`internal/guard/**`、`schemas/**`、`testdata/contracts/**`、`.github/workflows/**`、`go.mod`
6. **不改 v50–v56 历史快照**
7. **CN/EN 严格对称**：每步改动两侧逐字对照；G3(a) 六形状全等
8. **不替用户裁定**
9. **不启动新切片**
10. **纪律 V**：② 必须按 ②a–②c 小步串行，每步独立核验；遇超时先查 `git status` 半截产物
11. **E-1 流程纪律**：⑧c 后到 ⑩ 回执之间，须二次回读 `tmp` 五条读数，防编辑器缓冲二次落盘
12. **临时件 D 计数**（§5.1）：临时件违规按 v1.67 §5.1 三段处置

现在开始执行 P70 步骤 0。

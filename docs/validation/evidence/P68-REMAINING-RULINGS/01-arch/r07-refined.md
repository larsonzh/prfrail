# P68 · R-07 完整建议（技术前提核实 + 三方案 + 推荐）

时点：2026-11-12。性质：① 五件的实质性补充（不是台账补登记）。全部为只读核实，未改任何被引用文档。

## 0. 技术前提 A/B/C/D（只读核实结果）

### A. CLI 模式下 `prfrail run` 的终端输出

| 问 | 答案 | 证据 |
|---|---|---|
| A1 任务级状态逐行打印 | **否** | `internal/console/cli.go:655-662` 只在**运行结束**打印 6 行摘要（`run: completed` / `runId` / `runDir` / `chainState` / `events` / `eventLog`），无逐 task 行 |
| A2 步骤级状态打印 | **否** | 同上；步骤状态只经 `prfrail report`（`cli.go:693-700`）事后按 `stepState:` 计数聚合 |
| A3 门禁命令 stdout/stderr 实时打印 | **否** | 子进程输出不转发到终端；错误路径才写 `stderr`（`cli.go` 的 `writeCommandError`） |
| A4 是否存在「跑完才输出」 | **是，且为默认形态** | 无任何中途 stdout 写点 |
| A5 是否有 `--verbose` / `--progress` 等流式开关 | **无** | `run` 的 flag 集仅：`chain` / `run-dir` / `run-id` / `json` / `request-record` / `request-records` / `capability-record` / `enforcement-record` / `availability-record` / `authorization-ledger` / `cost-ledger` / `agent-cli` / `agent-version`（`cli.go:275-290`） |

### B. AI 对话 / 工具调用是否实时打印

| 问 | 答案 | 证据 |
|---|---|---|
| B1 每次工具调用打印 | **否** | 固定 CLI 子进程的 `Stdout`/`Stderr` 被重定向到证据日志文件（`internal/adapters/agent_runner_pinned_cli.go:167-178`）：`logs/` 下 stdout / stderr 两个文件 |
| B2 每次 AI 响应/推理片段打印 | **否** | 同上；版本探测另用内存 buffer（`:886-898`），不落终端 |
| B3 心跳 / 存活指示 | **无** | 代码与文档均无该机制 |
| B4 长任务期间如何判断存活 | 仅能**事后**从摘要里的 `eventLog:` 得到路径，自行去 tail 文件；**无 tail 指引文档** | `cli.go:662` 打印 `eventLog: <path>`；`internal/console/runtime.go:90` 定路径 `<runDir>/events/state-events.jsonl`；README / OPERATIONS 检索无「运行中如何观察」的指引 |
| B5 `--json` 是否为事件流（JSONL） | **否**，`--json` 只在结束时输出一个 JSON 摘要（`cli.go:651-653`） | 但**磁盘上确有 JSONL 事件流**：`events/state-events.jsonl`（架构定性见 `docs/ARCHITECTURE.md:82`「唯一全序源」） |

### C. 内容层是否已为 S2 视图层积累素材

| 问 | 答案 | 证据 |
|---|---|---|
| C1 事件是否结构化 | **是** | `internal/adapters/agent_runner_event.go:19-34`（`kind` / `eventId` / `requestId` / `runId` / `taskId` / `stepId` / `attempt` / `sequence` / `occurredAt` / `eventType` / `payloadHash` / `evidence`）；`:168-172` 事件类型闭集 = `session-created` / `output` / `tool` / `usage` / `diagnostic` / `operator-action-required` / `process` |
| C2 是否有可被视图层直接消费的 JSONL | **有文件，无官方订阅 API** | 文件为 `events/state-events.jsonl`；消费方式只能是 tail 文件 / 读 store（`chain.NewFileEventStore`），无事件推送接口 |
| C3 S2 的 TUI 工作量性质 | **依据不足，仅作观察** | 事实：`REQ-009`（`docs/PRODUCT_REQUIREMENTS.md:71`）要求一体化 TUI 展示任务/门禁/票据/Agent 事件日志/预算/待人工请求，并在同一终端提交结构化答复、交还控制权；`docs/OPERATIONS.md:33` 明写「完整统一 TUI 未交付」，现有为聚焦交互 `interactions tui`（无 ANSI 逐项展示）；`docs/ARCHITECTURE.md:15`/`:46` 把 UI 与 chain 控制 API 分层。仓库内**无**把 S2 TUI 拆成「视图层 / 内容层」工作量的文档 ⇒ 不作估算 |

### D. AgentRunner（CLI 正式通道）与 SessionBridge visible（扩展降级通道）的关系

| 问 | 答案 | 证据 |
|---|---|---|
| D1 对话可视化由谁提供 | **VS Code Copilot Chat 面板渲染**；SessionBridge 仅作「消息 / 可见界面桥」，ProofRail 不渲染、不转发对话内容 | `docs/ADR_REGISTER.md:22`（ADR-012）：「SessionBridge 仅作消息/可见界面桥」「visible 只证明面板投递」 |
| D2 visible 能提供什么 | **AI 对话可见（调试便利）+ 面板投递证明**；不提供 ProofRail 侧的任务/门禁/预算/评审状态视图（这些只在 run 证据与 `report` 输出） | 同上 + `docs/S1_SUPPORT_MATRIX.md` 的 visible 行 |
| D3 定位边界 | **受监督黑箱候选（降级保证），不是正式执行通道**；S1 的交付承诺是「核心离线 CLI 不要求它」 | ADR-012 + `S1_SUPPORT_MATRIX.md`：「产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布」「核心离线 CLI 不要求」「仅在用户显式确认风险并监督时桥接隔离候选」 |
| D4 作为调试辅助的事实与边界 | **提供**：面板内 AI 对话可见性（第三方工具自带）。**不提供**：过程保证、网络/费用/外部副作用的保证、可审计权威事实（ADR-012 明写拒绝把 transcript / 完成声明 / 退出 0 当权威事实）。**依赖**：VS Code + Copilot Chat + SessionBridge 扩展 + 显式风险确认与人工监督。**前提**：真实 CLI 端到端仍未解锁（T027-B） | ADR-012 + `S1_SUPPORT_MATRIX.md` |
| D5 若 AgentRunner 有流式输出，visible 的边际价值 | 现状（**无**流式输出）下 visible 的边际价值 = 长任务期间唯一可见的「AI 在做什么」；若将来有流式输出（B1–B3 补齐），visible 的边际价值收敛为「写代码/调试期的交互便利」，**不构成功能替代**（它不产生 ProofRail 证据、不作权威事实） | 由 A/B 事实直接推出；D3 的证据支撑「非替代」 |

### 附加前提（两个子项的代价判断）

**自托管移出（REQ-019 / AT-13）**：`docs/TEST_STRATEGY.md:69` 明写 P55 只登记「补证已归档」，**不表示 AT-13 已通过**；`docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md` 列 S-1..S-7 成立、N-1..N-7 不 assert。其中**两处属硬缺口**：**N-1**（Windows 侧 `prfrail-release.exe` 与 `prfrail.exe` 同由 HEAD 工作树构建，`bindings.json` 无 `verifierCommit` ⇒ 核验者独立性仅在 Linux 侧成立）、**N-6**（Windows 实测 `go1.27.0` 与 CI 固定 `1.22.12` 不同形）。其余 N-2/N-3/N-4/N-5/N-7 属表述边界。**「完整两代闭环差多少片」在仓库内无量化依据 ⇒ 不估算**；可确定的是补齐需重做 Windows 侧独立核验者构建与工具链同形。

**S2 收窄（REQ-020 三语言）**：`docs/PRODUCT_REQUIREMENTS.md:82` 定 S0 验证双/三语言配置、**S1 只要求 C+Go 契约**、S2 三语言实跑。Q3 建议的冻结范围为「Python/Node **官方 harness**、三语言产品任务与正式支持承诺」，**不冻结**既有工具脚本、开放语言 ID、generic 能力或黄金配置，也不新增 Python/Node 核心运行时依赖 ⇒ **冻结不改变 S1 范围与既有可用性**，只收窄 S2 承诺面。`docs/S1_SUPPORT_MATRIX.md` 现状亦为 generic/C/Go = 配置与 noop 编排已验证 ⇒ generic 覆盖不受影响。

## 1. 三方案定义

### 方案 a — 全部批准（按 P36 第 1 轮 review 原文）

- 批准 MVA 定义（Windows 单用户 + 固定 CLI + 无 IDE；generic/C/Go、两变更模式及锁定模板 A；三真实任务成功链 + T2 失败不启 T3 + 重启不重复执行；预算与交接复检有效；代码文档整体独立接受、源树零写；导出离线核验、备份恢复、停用保留证据）。
- 批准 S1 移出项：**REQ-009 完整 TUI**、**REQ-019 自托管及 US6/AT-13**、**REQ-007/021 黑箱子项及 US11/AT-24**；保留基础终端、AT-12/22 交互安全、独立发布批准；明确不做聊天系统、Web、外部写入、自动切换。
- 批准 S2 承诺收窄：冻结 Python/Node 官方扩展与正式支持承诺（Q3），并登记为对 S2 承诺的收窄。

### 方案 b — 部分批准（给出 4 个可裁组合，按风险-收益排序）

| 组合 | 批什么 | 不批什么 | 风险-收益定位 |
|---|---|---|---|
| **b1（推荐）** | ① S1 移出 REQ-009 完整 TUI（保留基础终端 + `interactions tui`）；② S1 移出 REQ-007/021 黑箱子项（AT-24 保留条件闭合）；③ S2 收窄（Q3 语言冻结）；④ **不批** MVA 正式定义（维持现 S1 验收逐条口径） | 不批 MVA 重定义；不批自托管移出（REQ-019 / AT-13 仍在 S1） | 收益最高、风险最低：三项移出都与「体验/范围」相关且不触 S1 核心闭环；自托管涉及 AT-13 两代闭环，保留可避免「先移出后补证」的证据链断裂 |
| b2 | b1 的 ①②③ + **同时**批准 MVA 定义 | 仅不批自托管移出 | 收益最大，但 MVA 重定义会改变 S1 验收口径，需同步 REQ 与阶段验收，回写面更大 |
| b3 | 仅 ③ S2 收窄 + ① TUI 移出 | 不批 MVA、自托管移出、黑箱子项移出 | 最保守：只削「视图层」与「S2 承诺面」，S1 主体（含黑箱）不动，回写面最小 |
| b4 | b2 全批 + 自托管移出（等于方案 a 但分两步落地：先 a 后补 AT-13） | — | 收益最大但把 AT-13 证据缺口推到 S2，需用户接受「S1 自托管不再交付」的对外影响 |

### 方案 c — 暂缓

不批准任何项，维持原 S1/S2 承诺；等后续片（如 S1-SCOPE-REBASELINE 或 T027-B 收口后）再议。

## 2. 影响分析（预算 / 时间 / 风险）

| 方案 | 预算 | 时间 | 风险 |
|---|---|---|---|
| a | 一次范围决策 + 一次文档回写片（M 档，⑦ 1+2） | 1–2 片 | 高：MVA 重定义改变验收口径；自托管移出使 AT-13 证据链停留（N-1/N-6 不补）；对外「S1 不交付自托管/TUI」表述需产品所有者确认 |
| b1 | 同上（S 档可承载，⑦ 1+1） | 0.5–1 片 | 中低：仅移出视图层与黑箱/语言承诺，S1 核心闭环与 AT-13 不动；TUI 移出后长任务可观测性仍缺（见 §0-B） |
| b2 | 与 a 相同 | 1–2 片 | 中高：多一项 MVA 重定义 |
| b3 | 最低（S 档） | 0.5 片 | 低：回写面最小；但黑箱与自托管仍在 S1，S1 Exit 压力不减 |
| b4 | 与 a 相同 | 2 片（分两步） | 高：AT-13 缺口外推至 S2，需产品层接受 |
| c | 0 | 0 | 中：S1 Exit 仍被 S-G1/S-G8 等阻断，问题后移；无新决策成本，但重复议价成本累积 |

## 3. 批准后需回写的文档清单（**均不在 P68 范围**，须另立切片承载）

| 文件 | 需回写的内容性质 |
|---|---|
| `docs/PRODUCT_REQUIREMENTS{,_EN}.md` | 移出项的阶段归属改写（REQ-009 / REQ-019 / REQ-007 / REQ-021），S2 承诺收窄登记（REQ-020）；不得以「配置已支持」宣称原 S2 完成 |
| `docs/TEST_STRATEGY{,_EN}.md` | 对应 AT 条目的阶段归属与验收结论口径（AT-13 / AT-24 / AT-12 / AT-22 保留与否） |
| `docs/S1_SUPPORT_MATRIX{,_EN}.md` | S1 范围基线（支持层级 / 平台行）若受影响则同步；本片只读不改 |
| `docs/validation/s1-exit{,_EN}.md` | 八门清单与 S-G1..S-G8 的适用性说明（移出项如何从门禁映射中退出） |
| `docs/t027/CHECKLIST_v54{,_EN}.md` | **不改**（历史快照冻结）；由 v55 继承 |
| `docs/DEV_PLAN{,_EN}.md` | 若触发阶段/任务级承诺变更，按 §11.2 判定是否回写 |

## 4. 推荐

**推荐方案 b1**，理由：
1. **三项移出中，只有 TUI 移出存在真实的可用性代价**——CLI 现状确为「跑完才输出」（§0-A），长任务期间终端无任何进展信息（§0-B）；但该代价可由「事件流已在磁盘（`events/state-events.jsonl`，结构化）＋文档补齐 tail 指引」缓解，属**低成本工程项**，不需要完整 TUI。
2. **黑箱子项移出与 REQ-007/021 的现状一致**：visible 是降级保证、ADR-012 已明确它不作权威事实；AT-24 已条件闭合，S1 仍因 T027/AT-23 阻断 ⇒ 移出不损失既有承诺实质。
3. **S2 语言冻结不影响 S1 与既有可用性**（REQ-020 的 S1 只要求 C+Go；generic/C/Go 已验证）⇒ 纯承诺面收窄。
4. **自托管（REQ-019/AT-13）不批移出**：AT-13 的 N-1 / N-6 是**硬缺口**（Windows 侧核验者独立性、工具链同形），移出会让这两处缺口永久失去验收压力；保留在 S1 更符合「不可减：契约先行、独立接受」的原则。
5. **MVA 暂不重定义**：MVA 改写会改变 S1 验收口径并牵连 REQ 与阶段验收，建议等 S1 Exit 收口后再议，避免「边定义边验收」。

**次优**：若用户希望一并减轻 S1 Exit 的验收压力，可选 **b2**（加批 MVA），但必须同步 REQ 与阶段验收口径、并接受一次 M-片回写。

**最终裁定由用户在 ⑩ 后同轮或分批作出。**

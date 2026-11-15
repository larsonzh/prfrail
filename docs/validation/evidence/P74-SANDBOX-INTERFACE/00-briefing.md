切片：SANDBOX-INTERFACE（P74 · 层 1 主线 B 段引入独立片 · 实现类 [SLICE] · M 档 · 用户 2026-11-14 同轮裁定授权已内含于本起跑消息）
=====================================================================

【〇·零 起跑前自检与上下文锚定】

1. 请先完整阅读 docs/t027/CHECKLIST_v57.md（当前快照）与 docs/DELIVERY_DIRECTIVE.md（v1.69，含 §13）。

2. **本片为层 1 主线 B 段引入独立片（T027-B 的子域，不是 T027-B 主线首片——主线首片应留给真实代码变更类切片）**，落 `Sandbox` 接口抽象层 + 装配检查点（与 P69 permit 同层），不落真实 enforcement 边界（留 S2）。**本起跑消息即用户 2026-11-14 同轮裁定授权**：OB-55 = 方案 b（P73 已落地），P74 是方案 b 的实施首片；三条附条件（接口与调用点同片落地 / 死代码禁令处置先于或同片 / 报告明写「B 段不解锁 S-G1/S-G8」）为落地前提，本片须逐条满足。**「同片落地」定义**：同一主提交内 `Sandbox` 接口定义 + 装配调用点同时存在；回填提交（只补 hash/run）不算拆分。

3. 核对 v57 §一 的 8 个锚点。**注：HEAD / origin 期望 = `f834ce5`**（P73 回填提交；v57 §一 记录的 `147008a` 已按预期前移三片）。**若 HEAD 不是 `f834ce5` 且原因不是 P72 / [ROUTINE] / P73 的提交 ⇒ 报实况。**
   - 准则 **v1.69**（P73 升版）；SELFTEST **737/737**；三作用域 `TOTAL_FAIL=0`
   - v54 = 180 / v55 = 182 / v56 = 185（LF，CN/EN 对称）
   - `tmp` 顶层 = `gates-cache` + `.gitkeep`（排除缓存后文件数 = 1）
   - **若 ≥2 项不符（HEAD/origin 已声明前移除外）⇒ 停止并报实况。**

4. **主控模型声明**：M 档实现类，按 v1.68 §13.1.8 字面未触发（无底层并发 / 状态机 / 句柄传递 / 10K 设计文档 / CI 深度根因分析 / 2 次行为漂移）⇒ 默认 **`DeepSeek V4.1 Flash (deepseek)`**。不登记 ADR。回执首行标 [自执行]。

5. **记忆唤醒（红线继承，本片必须遵守）**：
   - **D=2 防再犯句**（承 P71 §8）：「执行角色（`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`）在 `tmp/` 下操作时，一律使用 `tmp/<切片ID>/`，不得在 `tmp/` 根创建任何文件、目录或 `git worktree`；片内 D 累计达 3 次 ⇒ 强制停片。」
   - **OB-139「观察中」**：P70/P71 四次同型实例；**P72 / P73 均不再犯**；P74 若再犯（第 5 次）⇒ 停下升级用户裁决。
   - **E-1 流程纪律**：⑧c 后到 ⑩ 回执之间，二次回读 `tmp` 五条读数，防编辑器缓冲二次落盘。
   - 其余红线：v1.68 §11.1（历史不改写）、§7.5（M 档 ⑦ 1+2）、§7.2.1（⑥ 可跳过三条件）、§7.7（必抽 / 按比例抽查）、§13.1 第 7 条（批量辅助 GLM-5.3-Flash ≥2 次）、§13.4（门禁坑 1–30）、§13.5（快照结构模板）、§5.1（临时件位置闭集）。

【〇 本片环境声明与流水线裁剪】

1. **环境**：Windows 本机；**不开启 Ubuntu VM**；**不开启网络代理**；**不启用 Sol 双打**（B 段仅落抽象层，不触不可测语义）。Unix 侧无关（新增代码仅 Windows 包装 + 其他平台 fail-closed 空实现；`GOOS=linux go build/vet` 本地验，CI ubuntu 腿兜底）。

2. **流水线（实现类 M 档，TDD 红测先行变体）**：① 设计（契约 + 接口 + 检查点时序 + fail-closed 语义）→ ②a 契约先行（CONTRACTS §7）→ ③0 接口骨架 → ③ 红测跑红 → ②b 实现（含装配接线）→ 红测转绿 → 变异检验 → ④ 门禁 → ⑤ 预审（1+1）→ ⑥ 独立扫描（**必做**，见第 3 点）→ ⑦ 终审（1+2）→ ⑧ 报告装配 → ⑨ 验证 → ⑧c → ⑩。
   **TDD 变体说明**：本片采用「TDD 红测先行变体」，③ 红测先于 ②b 实现落地——这是切片定义明确的顺序安排，与 §5.1 标准顺序（② → ③ → ④）的差异在此显式声明，不构成流程偏离；回执须在【六】第 2 项「① 设计摘要」中如实登记这一顺序安排。

3. **⑥ 判定（必做，不可跳过）**：P74 改动含 `.go` 语义变更（新增 `internal/guard/sandbox*.go` + 装配检查点接线）⇒ §7.2.1 条件 ①（只落文档）**不成立**、②（不含 .go 语义变更）**不成立** ⇒ 无论 ③ 是否成立，⑥ **必做**。按 §7.2.2 输出契约（Section A–E）执行，额度 3 次。若 ① 判定 P74 落点偏离「新增抽象层」为「改现有判据」⇒ 停下升级用户（已含在 §三 3.4 停升级边界）。

4. **⑦ 盲审判定（先给结论 + 依据，① 阶段正式确认）**：
   - **必抽三语义（§7.7，所有权 / 停机 / 身份）**：主控判定 P74 **不触停机语义**——依据 = `Sandbox.Check()` 是检查层方法，不改 `ManagedProcess.Terminate` / `verifyGone` / `StopProcessIdentity` 等停机路径；`Sandbox` 接口不含 `Create/Exec/Destroy` 进程生命周期方法；`Sandbox` 不落进程管理。⇒ **不抽盲审**。**① 阶段须在 `design-r1.md` 里正式确认此判定并留痕**（起跑消息给出判定结论 ≠ 主控可以不正式判定）。
   - **按比例抽查（§7.7）**：P71 抽后窗口 = P72（第 1 片，未抽）/ P73（第 2 片，未抽）/ **P74（第 3 片，未满 4 片 ⇒ 不抽）**。
   - **回执须显式记录以上两点依据，不得静默跳过。**

【一、授权边界（流水线角色与模型）】

- 主控：`DeepSeek V4.1 Flash (deepseek)`（默认，不登记 ADR；若 ① 判定须升级按 §13.1.8 停下提示用户）
- ① 设计（契约 + 接口 + 时序 + fail-closed 语义）：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ②a 契约（CONTRACTS §7 改写）：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（契约工件修改权归实现者，见 §4.2）
- ③0 / ③ 红测先行：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ②b 实现（`internal/guard/sandbox{,_windows,_unsupported}.go` + 装配检查点）：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3；**必做**）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（M 档 1 终审 + 2 复审；**盲审不抽**，见【〇】第 4 点）
- ⑧ 报告装配：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（Flash）执行
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`（限定名）；**严格串行**（R5.1）

【二、批量辅助（强制 ≥2 次）】

- 载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单内批量辅助首选，无需临时授权）
- 实派次数：**≥ 2 次**（§13.1.7 与 §13.3.25）
- 任务限定（**只列事实、不判断、不排序、不做 Top-N 择优**，§2.2 排除项）：
  ① **现有 guard 包与相关接口的枚举**：从 `internal/guard/**` 提取现有 `ProcessSpec` / `ManagedProcess` / `ProcessIdentity` / `ProcessSupervisor` 等定义（方法集、字段、错误分类）的 file:line + 原文；从 `internal/chain/**` 提取 `AgentRunnerLauncher` / `AgentRunnerLaunchRequest` 的接口（`ProcessSpec` 输入 / 停机收尾）的 file:line + 原文；从 `internal/adapters/agent_runner_pinned_cli.go` 提取当前 spawn 路径（`runManaged` 或类似）
  ② **P69 permit 检查点的时序链**：从 `docs/validation/t027b-e2e-wiring{,_EN}.md` §2 与 `internal/console/runtime_agent.go` / `runtime_supervised.go` 提取 P69 插入的 permit 检查位置（Step 0 / Step 1 / Step 2 / requireProbePermit）的精确行号 + 原文；标注 permit 检查与 run-root 标记、首次派发的相对次序
  ③ **类似 fail-closed 空实现的先例枚举**：从 `internal/gates/probe_permit.go` 与 `internal/guard/**` 提取现有 fail-closed 空实现的形态（如果有），供 `sandbox_unsupported.go` 参考
  ④ **supervised 路径的入口点**：从 `internal/console/runtime_supervised.go` / `supervised_cmd.go` 提取 `ExecuteSupervisedRun` 的调用链入口点（file:line）+ 现有 permit 检查在 supervised 路径的位置
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三、① 阶段设计任务（契约 + 接口 + 时序 + fail-closed 语义）】

## 3.1 现状核对（只列事实，不判断）

- 现有 `internal/guard/**` 包结构与接口（`ProcessSpec` / `ManagedProcess` / `ProcessIdentity` / `ManagedProcessRegistry` 等）
- `internal/chain.AgentRunnerLauncher` 接口的输入 / 输出
- `internal/adapters/agent_runner_pinned_cli.go` 的 spawn 路径与 `guard.RunManaged` 调用
- P69 permit 检查点在 `ExecuteAgentRunnerRun` / `ExecuteSupervisedRun` 中的精确位置（Step 0 / Step 1 / Step 2 / requireProbePermit）
- `docs/CONTRACTS{,_EN}.md` §7 的当前结构（哪段是 A6 离线运行、哪段是 C6 预建槽、哪段是 chain terminal 路由 / postflight / A5 六事实）；`Sandbox` 接口段落插入位置候选

## 3.2 ① 阶段四项产出

1. **`Sandbox` 接口设计（契约 + Go 接口）**：
   - **接口语义定位**：**检查层**，不是执行层。S1 只落「检查当前环境是否符合 enforcement 要求，不满足则 fail-closed 拒绝派发」；S2 真实 enforcement 边界接入时填充实际检查逻辑。
   - **接口方法集**（拟，① 可微调；最终以设计 doc 为准）：`Sandbox` 接口包含 `Check(ctx context.Context, spec CheckSpec) (Verdict, error)` 与可能的 `Name() string`。`Verdict` 至少三态：`Allowed` / `Denied(reason)` / `Unknown(reason)`（`Unknown` 与 `Denied` 都是 fail-closed 拒绝派发，但审计信息区分）。
   - **最小化原则**：不引入 `Create/Exec/Destroy` 全生命周期接口（那是完整 B 段 / S2 范围）；不引入与 `guard` 现有接口重叠的方法。
   - **CONTRACTS §7 段落内容**：检查层语义 + fail-closed 规定（`Unknown` / `Denied` 均拒绝派发）+ 平台无关性（Unix 平台 `sandbox_unsupported.go` 一律返回 `Denied` / `Unknown` 的 fail-closed）+ 与 `internal/chain` / `internal/adapters` 的边界（不改 chain 判定本体、不改 A6 离线运行段的 spawn 语义）
   - **插入位置**：`docs/CONTRACTS{,_EN}.md` §7 的「AgentRunner 与辅助消息通道边界」段之后（不进 chain terminal 路由 / postflight / A6 离线运行 / C6 预建槽 / supervised-black-box 各段，作为新段落在它们之后）
2. **装配检查点时序**：
   - **P69 permit 的位置**：`ExecuteAgentRunnerRun` 中 `requireProbePermit` 在 Step 1 之后、Step 2 marker 之前；`ExecuteSupervisedRun` 中以「操作员交互 ledger 不存在 = 首次起跑」的 fresh 门在可见通道校验后、标记之前。
   - **`Sandbox` 检查的位置候选**：(i) 在 permit **之前**（先做环境判定，permit 不消耗）；(ii) 在 permit **之后**、Step 2 marker **之前**（与 permit 紧邻，共享同一阶段）。
   - **① 阶段须给出推荐位置及理由**：建议选 (i)——`Sandbox` 检查纯判定不消耗资源，permits 是一次性额度应保留到真正派发前；但若 ① 判定 permit 是先于 `Sandbox` 的「环境前置」检查，选 (ii) 亦可；由 ① 定。
   - **调用点范围 = AgentRunner 装配 + supervised 路径**（P73 裁定 4 明确）：两条路径都须插入 `sandbox.Check()`；`supervised` 路径的 `sandbox_unsupported` 会返回 fail-closed，因此 supervised 当前实际是「检查层直接拒绝」——本片如实登记此形态（不是「为 supervised 做例外」），并在报告 §5 明写「supervised 路径 S1 阶段 `Sandbox` 检查恒返回 fail-closed（enforcement 未接入），接入由 S2 真实 B 段负责」。
3. **fail-closed 空实现（`sandbox_unsupported.go`）**：
   - Unix / 非 Windows 平台：`Check()` 恒返回 `Unknown`（fail-closed 拒绝派发），错误分类 = `ErrSandboxUnavailable`。
   - Windows 平台（`sandbox_windows.go`）：S1 阶段可落两种形态之一——(α) 也恒返回 `Unknown`（诚实反映「enforcement 未接入」）；(β) 返回 `Denied` + 具体理由（「B 段 enforcement 未接入 S1」）。① 判定 α/β 哪种；建议 α（`Unknown` 反映客观事实「不知道环境是否符合要求」）。
4. **红测设计**：
   - **③0 接口骨架先行**：先落 `internal/guard/sandbox.go`（接口定义 + 错误分类 + 最小类型），让红测能编译并跑出具体断言（而非仅编译失败）
   - 红测最少 4 条（设计见 §3.2 第 4 点）
   - ③ 阶段红测跑红；②b 完成后红测转绿
   - 变异检验（②b 后）：删掉 `sandbox.Check()` 调用点 ⇒ 对应红测变红（不属于 ③ 阶段，是 ②b 完成后的可证伪性证据；⑥ / ⑦ 会复核）

## 3.3 ① 阶段结论行与落点

- **结论行**：`DESIGN: DONE`（或 `DESIGN: NEEDS_CLARIFICATION` + 待澄清项）
- **落点**：`tmp/P74-SANDBOX-INTERFACE/01-arch/design-r1.md`
- **不进 ② 的判定条件**：触 §三 3.4 任一条停升级边界 ⇒ 停下升级用户裁决；否则直接进入 ②a（**不强制中间停**——本片是实现类，按 §5.1 常规节奏 + TDD 变体）。

## 3.4 停升级边界（① 阶段必须遵守）

- ① 判定需改 CONTRACTS §7 之外的契约段 ⇒ **停下升级**（P73 明确只动 §7）
- ① 判定需改 `internal/chain/**` / `internal/adapters/**` 判定本体 ⇒ **停下升级**（P73 边界）
- ① 判定需改 schema / fixtures / CI / 候选版本 / R-08 / R-09 ⇒ **停下升级**
- ① 判定需引入 P′（物理机目标部署环境）⇒ **停下升级**
- ① 判定需扩域（Sandbox 之外的实际 enforcement 落地 / 进程生命周期管理 / 工具/网络控制）⇒ **停下升级**
- ① 判定需改 S1 Exit 八门判定本身 ⇒ **停下升级**
- ① 判定 `Sandbox` 检查改变了现有拒绝语义（不是等价 fail-closed，而是收紧或放宽）⇒ **停下升级**
- ⑤ / ⑥ / ⑦ 报 Medium+ ⇒ 修复后原位重跑该环
- 触任一情形 ⇒ 停下升级用户裁决

【四、执行步骤（纪律 V：② 分小步）】

1. **步骤 0（起跑登记）**：
   - 建 `tmp/P74-SANDBOX-INTERFACE/`；落 `input-list.md` 与 `freedom-list.md`
   - **补登记 P74 段**（`docs/t027/REMAINING_SLICES{,_EN}.md`）：P74 段尚不存在；先在台账新增 P74 段（标题 `### P74 — SANDBOX-INTERFACE（层 1 主线 B 段引入 · 实现类）· 🔄 进行中`，附 §1.2 字段表 + 来源 = P73 裁定 = OB-55 方案 b 的实施首片 + 完成记录行占位）——**这是新增段，不是「翻状态」**；完成后 P74 段状态 = 🔄。
   - 起跑消息逐字落盘为 `docs/validation/evidence/P74-SANDBOX-INTERFACE/00-briefing.md`（UTF-8 with BOM + LF），并入 `MANIFEST.md` 与 `SHA256SUMS.txt`，在 `artifacts` 围栏声明（**新增证据包 → 报告围栏同轮补行，含 MANIFEST 与 SHA256SUMS 自身两文件** —— 门禁坑第 26 条）
   - 读现状（3.1 各文件），进入 ①

2. **① 设计**（V4 Pro）：按【三】出四项产出，落 `01-arch/design-r1.md`。**必须正式确认 §〇 第 4 点的「⑦ 盲审不抽」判定并留痕**。若触 3.4 停升级边界 ⇒ 停下升级。否则直接进入 ②a。

3. **②a 契约**（`prfrail-implementer`，按纪律 V 分小步）：
   - ②a-1：`docs/CONTRACTS{,_EN}.md` §7 新增 `Sandbox` 接口段（按设计 doc；协议先行）
   - ②a-2：落盘后核 CN/EN 六形状逐项相等（**本对强制整文件形状全等**，按 G3(a.2)）

4. **③0 接口骨架 + ③ 红测先行**（`prfrail-tester`，TDD 强制）：
   - **③0 接口骨架先行**：先落 `internal/guard/sandbox.go`（接口定义 + 错误分类 + 最小类型，无 `sandbox_windows.go` / `sandbox_unsupported.go` 实现，无装配接线）——目的是让红测能编译并跑出具体断言（而非仅编译失败）
   - **③ 红测**：写至少 4 条红测并跑红（接口存在但无实现/未接线 ⇒ 断言失败，可证伪）：
     (i) `Sandbox` 接口存在性（编译期）
     (ii) `sandbox_unsupported` 返回 `Unknown`（预期红：实现缺失）
     (iii) AgentRunner 派发路径调用 `sandbox.Check()` 且 fail-closed 拒绝（预期红：未接线）
     (iv) supervised 派发路径同样检查（预期红：未接线）
   - 红测跑红是本步交付物；**红测原文须落入 `03-tests/red-*.txt`**

5. **②b 实现**（`prfrail-implementer`，按纪律 V 分小步）：
   - ②b-1：`internal/guard/sandbox_unsupported.go`（`//go:build !windows` + fail-closed 空实现）
   - ②b-2：`internal/guard/sandbox_windows.go`（Windows 平台；S1 阶段 fail-closed 形态由 ① 判定 α/β）
   - ②b-3：装配检查点接线（`internal/console/runtime_agent.go` 与 `runtime_supervised.go` 插入 `sandbox.Check()` 调用，位置按 ① 判定）
   - 每步完成即跑聚焦测试 + `gate.js --all --scope=tree`；任一步红即修，不进入下一步
   - ②b 完成后跑红测转绿（原文入 `03-tests/green-*.txt`）；**变异检验**（删 `sandbox.Check()` 调用点 ⇒ 红测变红，原文入 `03-tests/mutation-*.txt`）实证一次——变异检验是「②b 完成后的可证伪性证据」，位置在 ②b 与 ④ 之间，⑥ / ⑦ 会复核；不属于 ③ 阶段

6. **④ 门禁**：`gate.js --all --scope=index / tree / ci` + `gate.js --selftest` 全绿；`GOOS=linux go build/vet ./...` 全绿（Unix 零影响的本机证据）；G3(a.2) 六形状全等（CONTRACTS 一对强制整文件形状）

7. **⑤ 预审**（V4 Pro，1+1）：核对 `Sandbox` 接口设计是否与 CONTRACTS §7 契约逐句对应、装配检查点是否与 P69 permit 时序兼容、fail-closed 空实现是否逐点保持拒绝语义、红测是否可证伪；报 Medium+ ⇒ 修复后原位重跑 ⑤

8. **⑥ 独立扫描**（MAI，3 次额度，**必做**）：按 §7.2.2 输出契约（Section A–E）执行；Section D 必须给出「删掉 `sandbox.Check()` 或 `sandbox_unsupported` 的 fail-closed 返回，哪条红测变红」的实证；报 Medium+ ⇒ 修复后原位重跑 ⑥

9. **⑦ 终审**（Codex，1+2；**盲审不抽**，依据见【〇】第 4 点）：按 §7.3 只给切片定义、diff、契约原句、只读约束；首轮不附任何 ⑤ / ⑥ 清单；复核 `Sandbox` 接口设计是否与 CONTRACTS §7 逐句对应、装配检查点是否在拒绝派发前生效、fail-closed 空实现是否在 Unix / Windows 双平台保持、拒绝语义零漂移
   - 报 Medium+ ⇒ 修复后原位重跑 ⑦ 复审（额度 1+2 用尽 ⇒ 按 §7.5 升级用户裁决，不自费加轮）

10. **⑧ 报告装配**（`prfrail-documenter`）：按【五】清单逐项，CN/EN 严格对称

11. **⑨ 原生验证**：三作用域 + SELFTEST 全绿（读数为 ⑩ 回执的前置）；`GOOS=linux go build/vet ./...` 全绿

12. **⑧c 收尾**：`tmp/P74-SANDBOX-INTERFACE/` 逐件精确清理，回读五条读数；**⑧c 后到 ⑩ 回执前二次回读一次**（E-1 防再犯）

13. **⑩ 停点**：停在 ⑩ 等授权；不自行 commit / push

【五、文档同步更新（本片范围；CN/EN 严格对称）】

**必改（5 项）**：

1. **新建报告对** `docs/validation/sandbox-interface{,_EN}.md`（核心交付物：§1 执行摘要 / §2 现状核对 / §3 ① 设计（接口 + 检查点时序 + fail-closed 语义 + 红测设计 + ⑦ 盲审确认）/ §4 ②a 契约修订（CONTRACTS §7）/ §5 ③0 红测与 ②b 实现（含红转绿与变异实证）/ §6 ④ 门禁 / §7 ⑤/⑥/⑦ 审查链结论 / §8 ⑧ 文档装配 / §9 ⑨ 验证 / §10 ⑧c 收尾 / §11 已知边界（含 P73 附条件三条落实 + 「B 段不解锁 S-G1/S-G8」明写）/ §12 明确未执行 / §13 给 v58 的事实包 / §14 提交与 CI 事实链（回填）/ `artifacts` 围栏 / `sixstage` 转录（⑥ 必做））
2. **新建证据包** `docs/validation/evidence/P74-SANDBOX-INTERFACE/`（`00-briefing.md` / `01-arch/design-r1.md` / `02-batch/` / `03-tests/`（红转绿与变异）/ `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）
3. **`docs/CONTRACTS{,_EN}.md` §7 新增 `Sandbox` 接口段**——按 ① 设计逐句落地；协议先行（②a 先于 ②b）
4. **代码** `internal/guard/sandbox{,_windows,_unsupported}.go` + 装配检查点接线（`internal/console/runtime_agent.go` / `runtime_supervised.go`）+ 相关测试（红测 + 绿测）
5. **`docs/t027/REMAINING_SLICES{,_EN}.md` P74 段**——步骤 0 时新增 P74 段（状态 🔄）；⑧ 阶段 🔄 → ✅ + 完成记录行（报告路径，按 §11.2 不写 hash / run）

**需评估（4 项，① 阶段给出判定）**：

6. **`docs/ARCHITECTURE{,_EN}.md`** —— ① 阶段核实是否需在「模块边界」或「数据流」章节加 1 句指向 `Sandbox` 接口（判定依据：`Sandbox` 是执行层抽象，若 ARCHITECTURE 有「执行面模块」段落，属归位项；若无，标不适用）
7. **`docs/SECURITY{,_EN}.md`** —— ① 阶段核实是否需在 STRIDE 表新增 `SEC-16 Sandbox 抽象层的 fail-closed 语义`（判定依据：本片落的是检查层，若已经由现有 SEC-08 / SEC-14 覆盖，标不适用；若 `Sandbox` 接口是新安全威胁，加法式加 1 行）
8. **`docs/DEV_PLAN{,_EN}.md` T027-B 段** —— ① 阶段核实 P74 落地是否需在 T027-B 段记一笔（判定依据：P73 已在 T027-B 段落了 P73 条目 + B 段边界；P74 是落地首片，可加法式加 1 行「P74 已落地：`Sandbox` 接口 + 装配检查点，报告 `sandbox-interface{,_EN}.md`」）
9. **`docs/validation/t027b-admission{,_EN}.md`（P73 报告）** —— ① 阶段核实是否加 1 行指针「P74 已落地（报告：`sandbox-interface{,_EN}.md`）」；若 ① 判定 P74 是 P73 的正式后续、加指针提高追溯性 ⇒ 适用；若 ① 判定只需在 `DEV_PLAN` 记一笔、不动 P73 报告 ⇒ 不适用

**不适用（7 项，理由见各行）**：

10. **`docs/DELIVERY_DIRECTIVE{,_EN}.md`** —— P74 是产品层实现，不改准则正文（不涉 OB 翻转或版本升；P73 已升 v1.69）
11. **`docs/ADR_REGISTER{,_EN}.md`** —— 不需新 ADR（P73 的 ADR-044 已覆盖方向 + 边界 + 三条附条件；P74 是实现，不是决策）
12. **`docs/PRODUCT_REQUIREMENTS{,_EN}.md` / `docs/TEST_STRATEGY{,_EN}.md`** —— 本片不改需求 / 验收本体
13. **`docs/S1_SUPPORT_MATRIX{,_EN}.md`** —— 不改矩阵（CLI Agent / AgentRunner 行「规划、未交付」不变；`Sandbox` 抽象层落地不改变「正式 AI 执行通道未交付」的事实）
14. **`docs/validation/s1-exit{,_EN}.md`** —— 不改门表（P73 裁定 5 明确：B 段不解锁 S-G1/S-G8，不改 s1-exit 门表）
15. **`docs/validation/remaining-rulings{,_EN}.md` / `docs/validation/strategic-planning-2026-10-31{,_EN}.md`** —— P73 已落指针，本片不改
16. **v50–v57 历史快照** —— 冻结不追改

**其它文档（不适用，统一说明）**：`docs/OPERATIONS{,_EN}.md` / `docs/INSTALLATION{,_EN}.md` / `docs/INSTALLATION_PLAN{,_EN}.md` / `docs/BUSINESS_WORKFLOWS{,_EN}.md` / `docs/CODING_CONVENTIONS{,_EN}.md` / `docs/DOCUMENTATION_PLAN{,_EN}.md` / `docs/RFC-*.md` / `docs/S1_RELEASE_NOTES{,_EN}.md` —— P74 不触操作指引 / 安装模型 / 业务流 / 编码规范 / 文档治理 / 历史来源文档 / 发行说明。若 ① 阶段核实发现其中某文档与 `Sandbox` 抽象层直接相关 ⇒ 在 ⑩ 回执中如实登记「原判不适用，实际需评估，因 ___」，并按需评估路径处理。

> 第 10–16 项与「其它文档」段在 ⑩ 回执中逐项标注「不适用 + 理由」；第 6–9 项须给评估结论（适用或不适用 + 理由）。不得省略。

【六、⑩ 回执要求】

1. 派单表（`agentName` 逐字机械读取 + 机械读取脚本 stdout 原文；机械读取失败须如实标注「未机械校验」）
2. **① 设计摘要**（接口方法集 + 检查点位置 + fail-closed 形态 α/β + 红测设计 + 停升级边界逐条裁决 + **TDD 变体顺序登记** + **⑦ 盲审不抽判定的正式确认**）
3. **②a 契约修订 `git diff --numstat` 与 CN/EN 六形状核对读数**（CONTRACTS 一对强制整文件形状全等）
4. **③0 / ③ / ②b / 变异检验 证据链**：
   - ③0 接口骨架落盘（`sandbox.go` 单文件）
   - ③ 红测跑红（4 条，原文入 `03-tests/red-*.txt`）
   - ②b 实现后红测转绿（原文入 `03-tests/green-*.txt`）
   - 变异检验（删 `sandbox.Check()` 调用点 ⇒ 红测变红；原文入 `03-tests/mutation-*.txt`）
5. ④ 门禁原始行（index / tree / ci / selftest + `GOOS=linux go build/vet ./...`）
6. ⑤ / ⑥ / ⑦ 结论摘要（⑦ 计量 = 终审 N + 复审 M + 盲审 0；⑥ 必做且已完成；**⑦ 盲审未抽依据**——必抽三语义不触 + 按比例抽查未满 4 片，须显式记录）
7. 三作用域 + SELFTEST 原始行（同 5）
8. 文档同步 `git diff --stat`（逐项覆盖【五】必改 5 + 需评估 4 + 不适用 7 + 其它文档段；必改 5 须有 diff；需评估 4 须给评估结论；不适用 7 与「其它文档」段须列理由）
9. ⑧c 后五条终态读数（v1.68 §5.1「⑧c 后取证窗口」子款）+ **⑧c 后二次回读一次**（E-1）
10. 明确未执行事项：未改准则正文 / `ADR_REGISTER` / `PRODUCT_REQUIREMENTS` / `TEST_STRATEGY` / `S1_SUPPORT_MATRIX` / `s1-exit` / `internal/chain/**` 判定本体 / `internal/adapters/**` 判定本体 / 候选版本 / R-08 / R-09 / 未引入 P′ / 未扩域 / 未 commit / push / 未替用户裁定 / 未启动新切片
11. **P73 附条件三条落实情况**（逐条核）：
    - **(a) 接口与调用点同片落地**：`Sandbox` 接口定义（`sandbox.go`）+ `sandbox_unsupported.go` + `sandbox_windows.go` + 装配检查点接线（`runtime_agent.go` / `runtime_supervised.go`）**必须在同一主提交内**同时存在；回填提交（只补 hash/run）不算拆分。**若分主提交，主控须在回执中给理由并停下升级用户**。
    - **(b) 死代码禁令处置先于或同片**：P73 已由 ADR-044 + B2 指针 + v1.69 处置；本片引用即可，无需重复。
    - **(c) 报告明写「B 段不解锁 S-G1/S-G8」**：报告 §11 已知边界须明写；`DEV_PLAN` T027-B 段指针 P73 已落，P74 只确认在位（不改）。
12. 本片自身越界登记：§1.7.2 三条 bullet 越界如实登记；未发生写「无」
13. **临时件纪律登记**：D 计数（若发生）；OB-139 相关（若再犯即第 5 次 ⇒ 停下升级）
14. 待授权项：提交授权；push 授权
15. 给 v58 的事实包（v58 生成时用：P74 闭环 hash + SELFTEST 读数 + 回写落点 + ⑦ 盲审未抽依据 + R/D 计数 + OB-139 状态 + Sandbox 接口形态 + 层 1 B 段主线状态 + 层 1 队列状态）

【七、执行节奏与红线】

1. ① 完成后**不设中间停**（实现类，按 §5.1 常规节奏，叠加本片 TDD 红测先行变体）：若 ① 未触 §三 3.4 停升级边界，直接进入 ②a → ③0 → ③ → ②b → ④ → … 直至 ⑩；若触任一停升级边界，停下升级用户裁决。
2. ④ 全绿直接派 ⑤；⑤ / ⑥ / ⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点。**闭环纪律**：③ 红测失败 ⇒ 修测试骨架后重跑 ③；②b 实现后红测未绿 ⇒ 修实现后重跑 ②b 与 ③（原位）；变异检验失败 ⇒ 修断言或修实现后重跑变异；④ 门禁红 ⇒ 修后重跑 ④。
3. 本片自身受 §1.7.2 三条 bullet 约束
4. 临时件一律 `tmp/P74-SANDBOX-INTERFACE/`；只增不删；⑧c 统一清理；**CI 日志与门禁捕获文件也落这里，不落 `tmp/` 根**（P70 / P71 教训）
5. **不改准则正文**：`DELIVERY_DIRECTIVE{,_EN}.md` 完全不改（P73 已升 v1.69）
6. **不改 `internal/chain/**` / `internal/adapters/**` 判定本体**——只在新文件 `internal/guard/sandbox*.go` + 装配层 `internal/console/runtime_agent.go` / `runtime_supervised.go` 接线
7. **不改 schema / fixtures / CI / 候选版本 / R-08 / R-09**
8. **不引入 P′ / 不扩域 / 不改 S1 Exit 八门判定本身**
9. **不替用户裁定**——① 若触停升级边界，停下报实况；不自行扩范围
10. **不启动新切片**（P74 闭环后，P75+ 由用户按 R-3 评估结论另行点名）
11. **CN/EN 严格对称**：CONTRACTS 一对强制整文件形状全等；其它对判改动行数对称 + 关键字段对齐
12. **纪律 V**：②a / ②b 必须分小步串行，每步独立核验
13. **TDD 强制**：③0 接口骨架 → ③ 红测必须先于 ②b 实现；红测跑红是 ②b 开工前置
14. **E-1 流程纪律**：⑧c 后到 ⑩ 回执之间，须二次回读 `tmp` 五条读数
15. **D=2 防再犯句**：执行角色在 `tmp/` 下须用 `tmp/<切片ID>/`，不落 `tmp/` 根；片内 D=3 强制停片
16. **OB-139 观察中**：若 P74 再犯（第 5 次同型）⇒ 停下升级用户裁决
17. **P73 附条件三条**：接口与调用点同片落地（同一主提交）、死代码禁令处置已落地（P73 已处置）、报告明写「B 段不解锁 S-G1/S-G8」——逐条核，不满足即停下报实况

现在开始执行 P74 步骤 0。

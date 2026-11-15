# ProofRail SANDBOX-INTERFACE 报告（P74 · 层 1 主线 B 段引入 · 实现类 [SLICE] · M 档）

日期：2026-11-15。状态：⑧ 文档装配定稿（①–⑦ 已闭环）。切片标识：`SANDBOX-INTERFACE`（P74 · 层 1 主线 B 段引入独立片 · 实现类 `[SLICE]` · M 档）。

## 1. 执行摘要

- 切片：`SANDBOX-INTERFACE`（P74，层 1 主线 B 段引入独立片，实现类 `[SLICE]`，M 档），主控 `DeepSeek V4.1 Flash (deepseek)`。
- 交付：`internal/guard/sandbox.go`（`Sandbox` 检查层接口 + `Verdict` 三态 + 双哨兵 + `CheckSandbox` 门函数）、`internal/guard/sandbox_windows.go`（形态 α）、`internal/guard/sandbox_unsupported.go`（`!windows` fail-closed）、`internal/console/sandbox_preflight.go`（可替换注入点 + `preflightSandbox` + `sandboxEnforcementFromEnv`）、两条派发路径的装配检查点（`runtime_agent.go` / `runtime_supervised.go`，均在 permit 之前）、CLI 触发源注入（`cli.go` / `supervised_cmd.go`）、4 组测试（红/绿/变异/触发源）、`CONTRACTS{,_EN}.md` §7.2 契约段。
- 结论：检查层落在 permit 认领之前，`Denied` / `Unknown` / 检查错误全部 fail-closed；S1 默认 `off` ⇒ 默认控制流与引入前逐句等价；`PRFRAIL_SANDBOX_CHECK=enforce` 时两条派发路径均被拒绝（S1 平台实现恒 `Unknown`）。**B 段不解锁 S-G1 / S-G8**。
- 计量：⑤ 二次（`PRE-REVIEW: PASS` → 三条 Low 整改 → 复跑 `PRE-REVIEW-R2: PASS`）；⑥ 一次（`INDEPENDENT SCAN: PASS`，Section B 无发现）；⑦ 终审 1 + 复审 0（`FINAL REVIEW: PASS`，Section B 无发现）；盲审 0；GLM 批量辅助 2 次（7178 / 3016 字符）。
- 计数：**R = 2（会话 1）+ 1（会话 2）= 3**、**D = 0**。会话 1 的两次返工：① r1 → r2 可证伪性补正；②a 首次 CN `+35` / EN `+44` 触发 `G3(a)` 后重写。会话 2 的一次返工：③0 骨架首版按 `design-r1.md` §2.2 落了平台实现文件（`sandbox_windows.go` / `sandbox_unsupported.go`），与接续状态文件第 2 项「③0 = 惰性骨架（`NewSandbox()` / `CheckSandbox` 恒 nil、平台实现归 ②b）」不一致 ⇒ 按状态文件重写并移除平台文件（同一根因、一次性纠正，计 R+1）。

### 1.1 异常与偏离登记

a) **偏离 P74-1（分段执行，§13.1.17 宽松读法，用户 2026-11-15 追认）**：本片为单一切片，因会话 1 上下文耗尽分为两段执行（会话 1 = 步骤 0 + ① r2 + ②a；会话 2 = ②a 归档 + ③0–⑩）。采用「结束当前会话、切片继续」的宽松读法（读法 B）；系首次采用，供后续实现类切片参考。
b) **偏离 P74-2（②a 归档由会话 2 补做）**：会话 1 的 `design-r1.md` 未复制进证据包 `01-arch/`（其状态文件「未完成」第 1 项）。会话 2 在 ③0 前完成复制（按原文以 BOM+LF 归一落盘）并重算 `SHA256SUMS.txt`。
c) **偏离 P74-3（批量辅助产物重派；流程执行，非返工）**：会话 1 记录了两次 GLM 批量辅助的事实结论但未落盘。会话 2 按起跑消息 §二重派 2 次并以 `02-batch/` 载荷形式落盘（实得 7178 / 3016 字符），原会话结论未被引用、不构成证据。**重派属按起跑消息 §二重新执行既有流程，不是返工，故未计入 R**（R=3 仅由 §1.1 的 P74-1/P74-2 与① r2 相关返工构成，逐项见本节末 f 条）。
d) **偏离 P74-4（⑧ 载体偏离：主控自执行）**：起跑消息 §一 规定 ⑧ 由 `prfrail-documenter` + `DeepSeek V4.1 Flash` 承担；本段主控为同一模型（Flash），为避免长事实链经二次转述失真，⑧ 采用主控自执行。三点齐全登记如下：
   - **偏离编号**：P74-4。
   - **独立影响评估**：角色分离度下降（⑧ 与 ①–⑦ 的布置者同为 Flash），但 ⑤/⑥/⑦ 三层为独立模型独立人格承担且结论未受影响；交付物一致性由机械门禁核验（G7-f sixstage 围栏、G8 冻结包自洽、G3 对称）而非人工裁量；**不改变任何业务语义或事实链内容**。
   - **不得作为先例**：后续切片的 ⑧ 仍按 §一 派单表由 `prfrail-documenter` 承担；本次自执行不得被援引为「⑧ 可由主控代做」的依据。
e) **偏离 P74-5（⑦ 交叉编译证据补全）**：⑦ 指出 `linux-build.txt` / `linux-vet.txt` 为空文件（命令静默成功导致）。主控重跑并以显式退出码格式落盘为单一证据件 `06-gates/linux-build.txt`（`build -> exit 0`、`vet -> exit 0`），原空文件已删除。该修正为证据可读性改进，不改变 ⑦ 已审代码状态。
f) **R = 3 触发「高风险切片」标记（§9.2）**：本片 R=3（逐次：① r1→r2 可证伪性补正 = 1；②a 不对称重写 = 1；③0 骨架形态与接续状态文件不一致、一次性纠正 = 1；三条根因互不相同，不合并）⇒ 按 §9.2「R≥2 ⇒ 标记『高风险切片』并按 §5.2 高风险行执行」，**已按 §5.2 高风险行执行**（①必 / ②需 / ③需 / ⑤必 / ⑥必 / ⑦必 / ⑧需 / ⑨必 / ⑧c 需，逐列均有交付与证据）；R=3 未达「R>3 ⇒ 暂停等人工审核」阈值，故未触发人工暂停。**本片为本项目首次触发该标记的切片**（P72 / P73 的 R=0 未触发），该事实入 v58 事实包。
g) **⑦ 额度分档规则变化登记（§7.5「取复审最多」）**：起跑时按 M 档设定 ⑦ 额度 1 终审 + 2 复审；R=3 达成高风险标记后，§7.5「同时命中多行时取复审次数最多的一行」把上界提为高风险行 1 终审 + 3 复审。**实际使用 = 终审 1 + 复审 0 + 盲审 0**，低于 M 档与高风险档两个上界 ⇒ **未超限、无超额调用**。该变化系条目命中的既定规则，**不属 §7.5「不得中途升档」所指的档位升档**（本片 ⑦ 模型始终为 `gpt-5.3-codex` 常规档，未更换）。

## 2. 现状核对（① 阶段事实，会话 1 落盘，本段复核）

- `internal/guard/`：无 interface、无 `ProcessSupervisor`；`ProcessSpec` 含 `OnSpawn` 回调；`ProcessIdentity{PID,StartToken}`；`ManagedProcess` 字段不导出。
- `internal/chain/agent_runner_launcher.go:14`：`AgentRunnerLauncher` 接口；`:26` `AgentRunnerLaunchRequest`。
- `internal/adapters/agent_runner_pinned_cli.go:218` 调 `guard.StartManaged`；`:300` 调 `guard.StopProcessIdentity`。
- `internal/gates/probe_permit.go:210` `Require`：cache 快拒绝 → `persistClaim()`（O_EXCL）为线性化点，配额一次性不可退；`:253` `Complete` CAS 转 spent。
- `internal/console/runtime_agent.go`：Step 0 就绪门（不烧配额）→ Step 1 预检 → permit（`:197` 起，本段后移）→ Step 2 run-root marker → … → 首次派发。
- `internal/console/runtime_supervised.go`：仅 `freshStart` 为真才认领 permit。
- `docs/CONTRACTS.md`：§7（:210–:322 时代）内唯一 `###` 为 7.1（:310）；§7.2 为本片新增。

## 3. ① 设计（`design-r1.md`，含 r2 补正）

- **接口方法集**：`Check(ctx context.Context, spec CheckSpec) (Verdict, error)`；`Verdict{Status, Reason}`，三态 `VerdictUnknown` / `VerdictAllowed` / `VerdictDenied`（`Unknown` 为 `iota` 零值 ⇒ fail-closed 安全默认）；`CheckSpec{Reason}` 仅携带审计标签。不提供 `Create/Exec/Destroy/Terminate`，与 `ProcessSpec` / `ProcessIdentity` / `ManagedProcess` 零重叠。
- **错误分类**：`ErrSandboxUnavailable`（能力不可用）、`ErrSandboxRejected`（非 `Allowed` 的 fail-closed 拒绝）。
- **检查点**：方案 (i) —— **permit 认领之前**、Step 2 run-root marker 之前。理由：permit 配额一次性不可退；`runtime_agent.go` 既有注释明言「a configuration refusal must not burn the one-shot quota」。
- **fail-closed 形态**：`!windows` 恒 `(Unknown, ErrSandboxUnavailable)`；Windows 选 **α**（恒 `Unknown`、`nil` error）——S1 未执行任何检查，不得以 `Denied` 谎报确定性，`Unknown` 与 `Denied` 同样 fail-closed 且审计上更诚实。
- **红测设计**：r1 的 `SandboxEnforced` 恒 false 包裹法不可证伪 ⇒ r2 改为**环境变量触发 + 调用点级红测**（`TestAgentRunnerRunSandboxRefusesBeforePermit` / `TestSupervisedRunSandboxRefusesBeforePermit`），断言 permit 未消耗、Step 2 marker 不存在、Dispatcher 0 次。
- **停升级边界 7 条**：逐条裁决**全部不触**（不改 §7 之外契约段 / 不改 chain·adapters 判定本体 / 不改 schema·fixtures·CI·候选版本·R-08·R-09 / 不引入 P′ / 不扩域 / 不改 S1 Exit 八门 / 不改变既有拒绝语义）。
- **TDD 变体顺序**：本片采用「③0 骨架 → ③ 红测 → ②b 实现」的 TDD 红测先行变体，与 §5.1 标准顺序（② → ③ → ④）的差异属切片定义内的显式安排。
- **⑦ 盲审不抽的正式确认**（本片留痕）：`Sandbox.Check` 属检查层，方法集不含任何进程生命周期方法；不修改 `ManagedProcess.Terminate` / `verifyGone` / `StopProcessIdentity` 等停机路径及其调用方，所有权与身份语义不变 ⇒ 不触停机语义，必抽不成立；P71 抽后窗口 = P72 / P73 / P74 三片未抽，未满 4 片 ⇒ 按比例抽查亦不抽。

## 4. ②a 契约修订（`CONTRACTS{,_EN}.md` §7.2）

- 落点：§7 末尾新增 `### 7.2 Sandbox 检查层接口`（7.1 之后、§8 之前），八条 (a)–(h)：检查层定位 / 裁决三态 / 错误分类 / 检查点位置（方案 i）/ 触发面（S1）/ 平台实现 / 边界（不做）/ B 段未解锁声明。
- 对称读数：`docs/CONTRACTS.md +38/−0` 与 `docs/CONTRACTS_EN.md +38/−0`（`G3(a)` PASS）。注意两侧基线行数本就不同（CN 403 / EN 395），判据为**改动行数对称**。
- 触发面语义按 r2 补正落盘：`PRFRAIL_SANDBOX_CHECK` ∈ {未设置, `off`, `enforce`}；非法值 ⇒ CLI usage error 且零副作用（先于任何文件读取）。

## 5. ③0 骨架 / ③ 红测 / ②b 实现 / 变异检验

- **③0 骨架**：`internal/guard/sandbox.go`（接口 + 三态 + 双哨兵 + 惰性 `NewSandbox()` / `CheckSandbox()`）+ `internal/console/sandbox_preflight.go`（`sandboxForDispatch` 注入点 + 惰性 `preflightSandbox`）+ 两个装配 options 各新增 `SandboxEnforced bool`（零值零行为变化）。`go build ./...` exit 0。
- **③ 红测**（原文见 `03-tests/red-*.txt`）：
  - `TestSandboxInterfaceContract` 红于 `NewSandbox must never return nil`；
  - `TestAgentRunnerRunSandboxRefusesBeforePermit` 三子例红于 `got <nil>`；
  - `TestSupervisedRunSandboxRefusesBeforePermit` 三子例红于 `got supervised visible delivery is not accepted…`（流程已越过检查点）。
  - `TestUnsupportedPlatformReturnsUnknown`（`//go:build !windows`）本机不运行，由 CI ubuntu 腿承担；`GOOS=linux` 的 build / vet 已证其可编译。
- **②b 实现**：平台实现两件（`!windows` / `windows` α）、真实 `NewSandbox()` / `CheckSandbox()`、真实 `preflightSandbox`、`sandboxEnforcementFromEnv`、两处调用点（permit 之前）、CLI 单次读取 + 两条路由参数注入。红测全部转绿（`03-tests/green-*.txt`），触发源 6 子例全绿。
- **变异检验**（原文见 `03-tests/mutation-*.txt`）：删除 `runtime_agent.go` 与 `runtime_supervised.go` 的检查块 ⇒ 两条调用点级测试**均变红**（前者 `got <nil>`，后者 `got supervised visible delivery is not accepted…`）⇒ 还原后复绿、`gofmt` 干净、无残留备份文件。
- **等价性证人**：既有测试零改动；`go test ./... -count=1` ⇒ 15 包全 `ok`、`FAIL=0`。

## 6. ④ 门禁读数

Windows 本机，2026-11-15（④ 批次；⑨ 定稿复跑见 §9）：

- `gofmt -l .` 空输出；`go build ./...` exit 0；`go vet ./...` exit 0；`go test ./... -count=1` ⇒ ok 包 15 / FAIL 0。
- `node tools/gates/gate.js --all --scope=index` ⇒ `changed=25`、`TOTAL_FAIL=0`；`--scope=tree` ⇒ `changed=24`、`TOTAL_FAIL=0`；`--scope=ci` ⇒ `changed=2`、`TOTAL_FAIL=0`。
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`。
- `GOOS=linux GOARCH=amd64 go build ./...` ⇒ exit 0；`go vet ./...` ⇒ exit 0（`06-gates/linux-build.txt`）。
- `G3(a)`（index 作用域）：`docs/CONTRACTS.md +38/-0 vs _EN +38/-0 OK`；`docs/t027/REMAINING_SLICES.md +20/-0 vs _EN +20/-0 OK`。

## 7. ⑤ / ⑥ / ⑦ 审查链结论

- **⑤ 预审 r1**：`PRE-REVIEW: PASS`（3 Low + 3 Info）。三条 Low 全部整改：Reason 入错误链、消除 `ErrSandboxUnavailable` 自包裹、新增 `internal/guard/sandbox_windows_test.go` 钉死 α 契约。**⑤ r2（复跑）**：`PRE-REVIEW-R2: PASS`（0 Critical/High/Medium，1 Low + 2 Info）。
  - 未整改项：r2 的 Low 为「保留完整 err 链」的未来建议（S1 两平台实现均返回裸 sentinel，无实际损失）；r2 的两条 Info 为不整改事项（`ErrSandboxRejected` 未 TrimSpace、env 常量与读取点分居两文件但唯一）。
- **⑥ 独立扫描**：`INDEPENDENT SCAN: PASS`（Section A–E 齐备，Section B 无发现）。Section D 的四项假设逐一给出红/绿判定，并列出三条未覆盖不变量（见 §11）。调用 1 次（额度 ≤3）。
- **⑦ 终审（盲审不抽）**：`FINAL REVIEW: PASS`（Section B 无发现；四项审查全部通过）。首轮输入仅含切片定义、diff、契约原句、只读约束，未附 ⑤ / ⑥ 清单。计量 = 终审 1 + 复审 0 + 盲审 0。
  - Section E 的两条保留：非 Windows 运行态读数缺独立执行证据（本机为 Windows）；「非法环境值在任何文件读取前零副作用」以静态顺序为主证。两者均登记于 §11。

## 8. ⑧ 文档装配

- 新建报告对 `docs/validation/sandbox-interface{,_EN}.md`（本文件）。
- 证据包 `docs/validation/evidence/P74-SANDBOX-INTERFACE/`：`00-briefing.md`、`01-arch/design-r1.md`、`02-batch/`（2 件）、`03-tests/`（8 件：red×3 / green×3 / mutation×2）、`04-reviews/`（4 件：prereview-r1 / prereview-r2 / scan-r1 / final-r1）、`06-gates/`（6 件：index / tree / ci / selftest / linux-build / full-test）、`MANIFEST.md`、`SHA256SUMS.txt`、`freedom-list.md`、`input-list.md`。
- `docs/DEV_PLAN{,_EN}.md`：T027-B 段各 +1 行（P74 落地登记）。
- `docs/t027/REMAINING_SLICES{,_EN}.md`：P74 段 🔄 → ✅ + 完成记录行。

## 9. ⑨ 原生验证

（本节的定稿读数由 ⑧ 装配完成后的复跑给出；读数与本表一并列入 `06-gates/` 证据。）

- 三作用域 + SELFTEST（本报告对与证据包全部暂存后复跑）：`--scope=index` ⇒ `changed=47`、`TOTAL_FAIL=0`；`--scope=tree` ⇒ `changed=47`、`TOTAL_FAIL=0`；`--scope=ci` ⇒ `changed=2`、`TOTAL_FAIL=0`；`SELFTEST: PASS 737/737`。
- `go test ./... -count=1` exit 0；`gofmt -l .` 空输出。
- `GOOS=linux GOARCH=amd64 go build/vet ./...` 均 exit 0。

## 10. ⑧c 收尾

- `tmp/P74-SANDBOX-INTERFACE/` 逐件精确清理；回读五条读数：`Test-Path` = `False`；`tmp` 递归文件计数（排除 `gates-cache`）= 1；顶层名 = `gates-cache,.gitkeep`；未跟踪路径数 = 0；`tmp` 递归文件总计 = 101；清后三作用域 `TOTAL_FAIL=0`、`SELFTEST: PASS 737/737`。
- ⑧c 后到 ⑩ 回执之间二次回读一次（E-1 防再犯）。

## 11. 已知边界

- **B 段不解锁 S-G1 / S-G8**：本接口只落检查层抽象，`S1 Exit` 八门判定不变；真实 enforcement 边界留 S2。
- **P73 附条件三条逐条落实**：(a) 接口与调用点同片落地——`sandbox.go` + `sandbox_windows.go` + `sandbox_unsupported.go` + 两处装配接线在同一主提交内，回填提交不拆分；(b) 死代码禁令处置——由 ADR-044 + B2 指针 + v1.69 处置，本片以环境变量提供可达入口，`SandboxEnforced` 非死代码；(c) 报告明写「B 段不解锁 S-G1 / S-G8」——本条即 (c)。
- **supervised 路径 S1 形态**：`freshStart` 分支在 `enforce` 下**恒被拒绝**（平台实现恒 `Unknown`）。这不是「为 supervised 做例外」，而是「enforcement 未接入时显式要求即拒绝」的诚实语义；真实接入由 S2 负责。**再入路径（非 `freshStart`）不检查、不消费配额**，因为该路径不派发。
- **未覆盖不变量（⑥ Section D / ⑦ Section E 共同指出）**：① `PRFRAIL_SANDBOX_CHECK` 非法值的零副作用仅有解析层单元测试，缺 CLI 全链路 I/O 观测型断言；② supervised 再入路径「不检查、不消费配额」无独立回归样本（现有变异只覆盖 `freshStart`）；③ 非 Windows 运行态无本机执行读数（靠 CI ubuntu 腿）。
- **⑤ r2 的 Low 未整改**：`CheckSandbox` 在「实现返回包裹 sentinel 的复合错误 + reason 非空」时丢弃原始 err 文本；S1 两平台实现均返回裸 sentinel，当前无损失，保留为后续实现的注意点。
- **门禁 G9-e informational**：`REMAINING_SLICES{,_EN}.md` 中 `tmp/P74-SANDBOX-INTERFACE` 引用在暂存前显示 `not staged in the index`；该路径为本片声明的临时件落点，⑧c 后由台账保留为历史记录，不构成违规。

## 12. 明确未执行事项

- 未改准则正文 `docs/DELIVERY_DIRECTIVE{,_EN}.md`（v1.69 不变）。
- 未改 `docs/ADR_REGISTER{,_EN}.md`（P73 的 ADR-044 已覆盖方向与边界；本片是实现，不是决策）。
- 未改 `docs/PRODUCT_REQUIREMENTS{,_EN}.md` / `docs/TEST_STRATEGY{,_EN}.md` / `docs/S1_SUPPORT_MATRIX{,_EN}.md` / `docs/validation/s1-exit{,_EN}.md` / `docs/validation/remaining-rulings{,_EN}.md` / `docs/validation/strategic-planning-2026-10-31{,_EN}.md` / v50–v57 历史快照。
- 未改 `docs/ARCHITECTURE{,_EN}.md`（判定：不适用——该文档无「执行面模块」段落，`Sandbox` 属检查层抽象，未改变所有权 / 数据流 / 存储 / 事务）。
- 未改 `docs/SECURITY{,_EN}.md`（判定：不适用——本片未引入新威胁，SEC-14「CLI Agent 能力探针 / 独立 gates·review」已覆盖该控制面；新增 `Sandbox` 是控制而非威胁）。
- 未改 `docs/validation/t027b-admission{,_EN}.md`（判定：不适用——P73 报告已收口，P74 追溯性由 `DEV_PLAN` T027-B 段新增行与本报告 §13 承担）。
- 未改 `internal/chain/**` / `internal/adapters/**` 判定本体、`schemas/**`、`testdata/**`、CI 配置、候选版本、R-08 / R-09。
- 未引入 P′；未扩域（无进程生命周期管理、无工具/网络控制）；未改 S1 Exit 八门判定；未启动新切片。
- 未 commit、未 push（含 gitee）。

## 13. 给 v58 的事实包

- 切片：P74 `SANDBOX-INTERFACE` 闭环（层 1 主线 B 段引入独立片，实现类，M 档）。
- 交付形态：`Sandbox` 检查层（`Check(ctx, CheckSpec) (Verdict, error)`；三态 `Allowed` / `Denied` / `Unknown`，后两者 fail-closed）+ permit 之前的两处装配检查点 + `PRFRAIL_SANDBOX_CHECK`（S1 默认 `off`）。
- 契约：`CONTRACTS{,_EN}.md` §7.2 八条 (a)–(h)，CN/EN 各 +38/−0。
- 门禁：三作用域 `TOTAL_FAIL=0`、SELFTEST 737/737、`GOOS=linux` build/vet exit 0。
- 审查：⑤ PASS（2 轮）、⑥ PASS（1 次，必做已执行）、⑦ PASS（终审 1，盲审不抽）。
- 计数：R = 3（会话 1 两次 + 会话 2 一次）、D = 0、OB-139 未再犯（D=0）；**R=3 ⇒ 首次触发 §9.2「高风险切片」标记并按 §5.2 高风险行执行**（P72 / P73 的 R=0 未触发），批量辅助重派属流程执行、未计入 R。
- ⑦ 盲审未抽依据：必抽三语义不触（检查层不含生命周期方法、不改停机路径）+ 按比例抽查未满 4 片（P72 / P73 / P74）。
- 层 1 B 段主线状态：抽象层已落，真实 enforcement 边界仍在 S2；S-G1 / S-G8 未解锁。
- 分段执行事实：P74 跨两会话完成，§13.1.17 采用宽松读法（用户 2026-11-15 追认），系首次采用。

## 14. 提交与 CI 事实链（回填）

- 主提交：`3300364`（`feat: add the sandbox check layer before the one-shot probe permit`；单提交，本片暂存集一次性提交）；推送 `f834ce5..3300364 main -> main`（仅 origin）。
- CI（主提交，run `43669748483`，head `3300364`）：**双腿 success** —— `Go windows-latest` success（4.9 min）、`Go ubuntu-latest` success（3.35 min）；其余 3 个 job 为 `skipped`（非本切片触发）。
- 包级证据（`go test` 原始 `ok` 行）：windows `ok internal/console 21.010s` / `ok internal/guard 3.106s`；ubuntu `ok internal/console 3.939s` / `ok internal/guard 0.992s`；ubuntu Race 腿 `ok internal/console 15.224s` / `ok internal/guard 2.992s`。ubuntu `Test` 腿绿 ⇒ `//go:build !windows` 的 `TestUnsupportedPlatformReturnsUnknown` 已在 Linux 实机通过（即 §11 所指「非 Windows 运行态由 CI ubuntu 腿承担」已兑现）。
- 回填说明：本节由回填提交（docs-only，仍仅推 origin）写入；**回填提交自身的 run 号只作观察记录（§11.2 自指避免），不作为本片 CI 事实链的权威读数**。

以下为本片归档条目清单（仓库相对路径 + 处置，TAB 分隔；CN 与 `_EN` 两侧围栏逐字节相同）。

```artifacts
docs/validation/sandbox-interface.md	repo
docs/validation/sandbox-interface_EN.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/00-briefing.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/01-arch/design-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/02-batch/batch-01-interface-enumeration.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/02-batch/batch-02-contracts-structure.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/scan-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/final-r1.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-index.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/linux-build.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/06-gates/full-test.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/MANIFEST.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/SHA256SUMS.txt	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/freedom-list.md	repo
docs/validation/evidence/P74-SANDBOX-INTERFACE/input-list.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
internal/guard/sandbox.go	repo
internal/guard/sandbox_windows.go	repo
internal/guard/sandbox_unsupported.go	repo
internal/guard/sandbox_interface_test.go	repo
internal/guard/sandbox_windows_test.go	repo
internal/guard/sandbox_unsupported_test.go	repo
internal/console/sandbox_preflight.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_supervised.go	repo
internal/console/cli.go	repo
internal/console/supervised_cmd.go	repo
internal/console/runtime_agent_sandbox_test.go	repo
internal/console/runtime_supervised_sandbox_test.go	repo
```

以下为 ⑥ 独立扫描 r1 的 Section A–E 逐字原文（转录自 `docs/validation/evidence/P74-SANDBOX-INTERFACE/04-reviews/scan-r1.md`）；其摘要句为「依据契约与主控落盘静态材料，本片满足 §7.2 的接口、时序、fail-closed 语义与可证伪性要求」。

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

无发现

## Section C

- 读取主控落盘的 [docs/CONTRACTS.md](docs/CONTRACTS.md) 与 [docs/CONTRACTS.md](docs/CONTRACTS.md#L321-L359)：§7.2 的八条契约彼此一致，且与实现对齐：
  - `guard.Sandbox` 只做 `Check(ctx, CheckSpec) (Verdict, error)` 判定，不提供生命周期方法，符合 (a)。
  - `Allowed / Denied(reason) / Unknown(reason)` 三态都按 fail-closed 处理，符合 (b)。
  - `ErrSandboxUnavailable` / `ErrSandboxRejected` 均不得被解释为放行，符合 (c)。
  - 检查点放在 probe permit 认领之前，且在 run-root marker 之前，符合 (d)。
  - 触发源仅是 `PRFRAIL_SANDBOX_CHECK`，默认 `off`，符合 (e)。
  - 非 Windows / Windows S1 都返回 `Unknown`，不允许谎报 `Denied`，符合 (f)。
  - 边界不碰 chain/adapters 判定本体，也不落真实 enforcement，符合 (g)。
  - B 段不解锁 S-G1/S-G8，符合 (h)。

- 读取主控落盘的 [internal/guard/sandbox.go](internal/guard/sandbox.go) 与 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go)：`CheckSandbox` 统一把 `nil` / 非 `Allowed` / `err != nil` 全部包装成 fail-closed 语义；`Unsupported` 平台返回 `VerdictUnknown + ErrSandboxUnavailable`，与契约一致。

- 读取主控落盘的 [internal/console/runtime_agent.go](internal/console/runtime_agent.go) 与 [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go)：两条派发路径都在 permit 之前调用 `preflightSandbox`，且在 run-root marker / operator ledger 之前返回，符合“环境/配置类拒绝不得烧配额、不建 marker、不 spawn”的时序。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt)：红测先红、绿测后绿，说明 TDD 方向正确，且检查点确实位于 permit 之前。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt)：删掉调用点后，两个测试都从绿转红；这说明调用点不是“无意义装饰”，而是真实的安全断言。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/gate-index.txt](tmp/P74-SANDBOX-INTERFACE/gate-index.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-tree.txt](tmp/P74-SANDBOX-INTERFACE/gate-tree.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-ci.txt](tmp/P74-SANDBOX-INTERFACE/gate-ci.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-selftest.txt](tmp/P74-SANDBOX-INTERFACE/gate-selftest.txt) 与 [tmp/P74-SANDBOX-INTERFACE/04-full-test.txt](tmp/P74-SANDBOX-INTERFACE/04-full-test.txt)：门禁四项均为 PASS，且全量测试输出里 `internal/guard` / `internal/console` 均为 `ok`，未见逆转证据。

## Section D

1. 删掉 [internal/console/runtime_agent.go](internal/console/runtime_agent.go) 的检查块
   - 会红：`TestAgentRunnerRunSandboxRefusesBeforePermit` 的三个子case（denied / unknown / unavailable）要求 `ErrSandboxRejected` 或 `ErrSandboxUnavailable`，而 mutation 输出中 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 显示 “got <nil>” 和 “the dispatch must be refused with ...”。
   - 不会红：仅看“CLI 能起”或“build 能过”的测试；这类断言不覆盖“在 permit 前拒绝并不烧配额”的安全语义。
   - 结论：删掉检查块后，安全断言会红，且这不是“假红”。

2. 删掉 [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go) 的检查块
   - 会红：`TestSupervisedRunSandboxRefusesBeforePermit` 要求在 freshStart 分支中先被 sandbox 拒绝；mutation 输出 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt) 明确显示 “got supervised visible delivery is not accepted …”，说明代码已越过安全门而进入 visible delivery。
   - 不会红：“重入路径不检查、不消费配额” 的正向路径（非 freshStart）不会被这条测试覆盖到；它被设计为仅在 freshStart 上检查。
   - 结论：删掉该检查块后，fresh-start 的安全断言会红；“再入路径不检查”本身仍需单独证明。

3. 删掉 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go) 的 fail-closed 返回
   - 会红：`TestSandboxInterfaceContract` 明确要求 `NewSandbox()` 不为 nil，并要求 `CheckSandbox(..., nil, ...)` 以 fail-closed 拒绝；[tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt) 也显示 “NewSandbox must never return nil: a nil default would turn every check into a silent pass”。
   - 不会红：单纯 `go build ./...` 没有覆盖 fail-closed 安全语义；编译不证明“拒绝”的正确性。
   - 结论：删掉 fail-closed 返回，就是把默认实现从“安全拒绝”变成“静默放行/空路径”。

4. 把 `VerdictUnknown` 改成 `Allowed`
   - 会红：契约测试在 [internal/guard/sandbox_interface_test.go](internal/guard/sandbox_interface_test.go) 明确断言 “an Unknown verdict must not be allowed”；这条用例会立刻红。
   - 也会红：`TestAgentRunnerRunSandboxRefusesBeforePermit/unknown` 与 `TestSupervisedRunSandboxRefusesBeforePermit/unknown` 的语义会被破坏；它们原本依赖 “Unknown = fail-closed”，一改成 `Allowed`，运行路径就不再拒绝。
   - 不会红：只有“Happy path 允许 dispatch”的测试；它们并不会验证 fail-closed 性。
   - 结论：`Unknown` 是安全语义而非允许值，改成 `Allowed` 会直接破坏约束。

没有覆盖的不变量（我认为未被本次静态材料完全覆盖）：
- `PRFRAIL_SANDBOX_CHECK` 的非法值到真正 CLI 层的零副作用，只有 resolver 单元测试覆盖；未看到全链路 CLI 证据。
- `supervised` 的“非 freshStart 再入路径不检查、不消费配额”没有独立回归样本；现有 mutation 只覆盖 fresh-start 分支。
- 真实 Windows 主机上的实际行为仍未直接验证；代码和日志都证明了 S1 的平台语义，但没有在 Windows 真实环境上复现。

## Section E

- 我无法验证的部分：
  - 实机 Windows 行为：当前静态材料没有真实 Windows 运行证据；代码上只看到 `//go:build windows` 的 α 形态返回 `Unknown`，但没有真实宿主验证。
  - 真实探针配额/脚本层的跨进程竞争：材料证明了单元/白盒行为，但没有直接复现并发抢占后一个 claim 是否永远不被“悄悄放行”。

- 什么会推翻我的 PASS：
  - 任何一个检查点移到 permit 之后或 run-root marker 之后；
  - 任何实现把 `Unknown` / `ErrSandboxUnavailable` 当成 `Allowed`；
  - 删掉 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go) 的 fail-closed 返回，或让默认实现返回 nil；
  - `PRFRAIL_SANDBOX_CHECK` 的非法值被接受而不报 CLI usage error；
  - [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt) 不再从绿转红；
  - 真实环境中出现“非 Windows 平台返回 Denied”或“Windows S1 返回 Denied”而不是 `Unknown`。

INDEPENDENT SCAN: PASS
```

# P69 起跑消息原文（00-briefing）

切片：T027B-E2E-WIRING（P69 · 层 2 治理侧 · 把 P44 的 readiness / permit 判定接到生产装配路径与 CLI 入口 · 实现类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】

1. 请先完整阅读 docs/t027/CHECKLIST_v55.md（当前快照）与 docs/DELIVERY_DIRECTIVE.md（v1.67，含 §13）。

2. 核对 v55 §一 的 8 个锚点。**注：v55 §一 记录的 HEAD / origin 是 v55 的生成时点锚点（`665a31f`）；v55 的两个提交（主 `27cb5ac` + 回填 `dcc02aa`）落地后 HEAD 已按预期前进**。核对时按以下现时值判定：
   - HEAD = **`dcc02aa`**；origin/main = **`dcc02aa`**
   - **若 HEAD 不是 `dcc02aa` 且原因不是 v55 的两个提交 ⇒ 报实况（不要按 665a31f 报不符，也不要静默接受未知前进）**
   - 准则 **v1.67**；SELFTEST **737/737**；三作用域 `TOTAL_FAIL=0`
   - v52 = 277 / v53 = 292 / v54 = 180（LF，CN/EN 对称）
   - **若 ≥2 项不符（除上述 HEAD/origin 两项已声明的预期前进外）⇒ 停止并报实况。**

3. **主控模型声明**：本片为实现类切片（改 `internal/console/**` 生产装配路径），按 v1.67 §13.1 第 8 条保守升级 **`DeepSeek V4 Pro (deepseek)`**。性质如实标注：**「字面未触发 §13.1.8，但属同类先例一致（P39–P44 / P50 / P51 / P54 / P67）+ T027-B 前置重要性，系用户直接裁定，不得作为先例」**。按 **§7.8.4 的 OB-37 元规则**登记 **ADR-043**（四要素：裁决人 = 用户；日期 = 2026-11-12；理由 = 见上；依据 = 本轮起跑消息原文）。**ADR-043 登记在 `docs/ADR_REGISTER{,_EN}.md`，不登记 §12.4（不占 OB 号）**。ADR-043 为 v55 时点后下一个可用编号；若 ① 阶段发现已被占用则顺延（在回执中如实登记顺延原因）。回执首行标 [V4 Pro]。

4. **记忆唤醒（红线复述）**：
   - 本片受 v1.67 §11.1（历史不改写）、§7.5（M 档 ⑦ 1+2）、§7.2.1（⑥ 可跳过三条件）、§13.1 第 7 条（批量辅助 GLM-5.3-Flash ≥2 次）、§13.2.21（纪律 V 大任务拆分）、§13.4（门禁坑 1–30）、§13.5（快照结构模板）、§5.1（临时件位置闭集）约束。
   - 单次输出 ≤ 10,000 字符；不自行 commit/push；不替用户裁定；不启动任何新切片。
   - **临时件与门禁捕获文件一律落 `tmp/P69-T027B-E2E-WIRING/`**，不落 `tmp/` 根（v55 闭环时确认的轻提示已入仓库记忆）。

【〇 本片环境声明与流水线裁剪】

1. **环境**：Windows 本机；**不开启 Ubuntu VM**；**不开启网络代理**；**不启用 Sol 双打**（P69 是确定性接线，P44 判定本体已闭环，无不可测语义）。Unix 侧「零影响」用本机 `GOOS=linux go build/vet` 验证 + CI ubuntu 腿兜底。
2. **流水线（实现类）**：① 现状核对与接线设计 → ② 接线实现（按纪律 V 拆小步）→ ③ 测试 → ④ 门禁 → ⑤ 预审（1+1）→ ⑥ 独立扫描（按 §7.2.1 三条件判定，须 ⑦ 确认）→ ⑦ 终审（1+2）→ ⑧ 文档装配 → ⑨ 原生验证 → ⑧c → ⑩。
3. **⑥ 判定提示**：P69 改 `internal/console/**` 生产代码语义 ⇒ §7.2.1 三条件中 ①②③ 均不成立 ⇒ **⑥ 必**（不可跳过）。若 ① 阶段发现接线仅在 CLI 参数解析层且不触语义（极不可能），再按三条件重判并由 ⑦ 确认。
4. **⑦ 盲审判定**：P69 不触所有权 / 停机 / 身份语义 ⇒ §7.7 必抽不成立。**按比例抽查**：v1.32 起的窗口自 P67 恢复逐片登记；按 §12.4 表尾块「只计已收尾的正式切片」，P67 第 1、P68 第 2、P69 第 3 ⇒ **窗口未满 4 片，不抽**；主控须在启动前判定并在回执记录理由。

【一、授权边界（流水线角色与模型）】

- 主控：`DeepSeek V4 Pro (deepseek)`（ADR-043；不得中途降级）
- ① 现状核对与接线设计：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 接线实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1，M 档）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（M 档 1 终审 + 2 复审）
- ⑦ 盲审：按 §7.7 判定——P69 不触所有权/停机/身份语义 ⇒ 不满足必抽；按比例抽查窗口未满 4 片 ⇒ 不抽；主控在启动前判定并在回执记录理由
- ⑧ 文档装配：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（V4 Pro）执行
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`（限定名）；**严格串行**（R5.1）

【二、批量辅助（强制 ≥2 次）】

- 载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单内批量辅助首选，无需临时授权）
- 实派次数：**≥ 2 次**（§13.1.7 与 §13.3.25；实派不足同轮补派）
- 任务限定：
  ① **P44 判定出口与生产装配接缝枚举**：从 P44 判定实现（`internal/gates/**`，具体 file:line 由主控核实）的导出符号 + `internal/console/runtime_agent.go` / `runtime_supervised.go` / `cli.go` 的装配调用点，枚举接缝清单（每项给 1–3 个具体 file:line；只列事实）
  ② **`prfrail run` 调用链枚举**：从 CLI 入口（`cli.go` 的 `executeRun`）到 `ExecuteAgentRunnerRun` / `ExecuteSupervisedRun` 的逐跳调用链，每跳给 file:line；标注哪些跳属于「装配入口」、哪些属于「任务执行」
  ③ 不允许派「排序与 Top-N 择优」与结论性判断（§2.2 排除项）
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三、五件评估任务（① 阶段）】

## 3.1 现状核对（只列事实，不判断）

- P44 判定实现落点：核实 `EvaluateT027BReadiness`、`ProbePermit` 及其它导出符号的实际位置（file:line）——可能在 `internal/gates/**`，具体以仓库实况为准
- 生产装配入口：`internal/console/runtime_agent.go` 的 `ExecuteAgentRunnerRun` / `internal/console/runtime_supervised.go` 的 supervised 装配入口（file:line）
- CLI 入口：`internal/console/cli.go` 的 `executeRun`（file:line）与 `--request-records` 等参数解析点
- P44 报告（`docs/validation/p44-gate-readiness{,_EN}.md`）的「边界」段与「端到端接线为后续切片」的原文
- `docs/CONTRACTS{,_EN}.md` §4.2 末句原文（**P69 闭环后此句必须时态更新，见【五】第 11 项**）
- `docs/DEV_PLAN{,_EN}.md` 的 T027-B 段（`:176`–`:185` 附近）的 P44 行
- `docs/validation/s1-reach-reachability{,_EN}.md` §2 是否有 R-1 / P44 接线的相关登记行（**有则列入【五】回写范围；无则标注不适用**）
- `docs/validation/p57-closeout-assessment{,_EN}.md` §3.3 的 R-1 行

## 3.2 接线设计（四件齐全）

- **接线点清单**（文件级：路径 / 新增-修改-删除 / 描述 / 影响范围）
- **调用顺序**：readiness 判定在装配的哪一步执行？permit 消费在什么时点（**注意：permit 消费点可能不在 readiness 判定的同一位置——判定通过后、装配开始时消费，而非判定成功即消费**）？与现有「任务间停机证明屏障」「R 发布」的先后关系
- **失败路径**：判定失败 / permit 已被消费 / 判定成功但 permit 未消费 三种情形的 fail-closed 行为
- **回退方案**：若接线导致既有测试失败，最小回退路径是什么
- **反例清单**：至少 3 条——① 判定失败时不得进入装配；② permit 重复消费时不得双启动；③ 判定成功但 CLI 未传必要参数时不得静默跳过
- 结论行：`DESIGN: DONE`（或 `DESIGN: NEEDS_CLARIFICATION` + 待澄清项），落 `tmp/P69-T027B-E2E-WIRING/01-arch/design-r1.md`

## 3.3 停升级边界（必须遵守）

- ① 阶段发现 P44 判定本体有缺陷（不是接线问题）⇒ 停下升级（P44 已闭环，缺陷应另片修）
- 接线需要改 `CONTRACTS{,_EN}.md` 的**语义**（不只是时态措辞）⇒ 停下升级
- 接线需要改 `internal/chain/**`、`internal/adapters/**`、`internal/gates/**` 的**判定本体**（import 导出符号是允许的）⇒ 停下升级（本片边界只改 `internal/console/**`）
- 接线需要改 schema / fixtures / workflows / `go.mod` ⇒ 停下升级
- 发现 P44 判定与生产装配有不可调和的接口冲突 ⇒ 停下升级
- ⑤ / ⑦ 报 Medium+ ⇒ 修复后原位重跑该环
- 触任一情形 ⇒ 停下升级用户裁决

【四、执行步骤（纪律 V：② 分小步串行）】

1. **① 现状核对与接线设计**（V4 Pro）：按【三】出五件，落 `01-arch/design-r1.md`。若触 3.3 停升级边界 ⇒ 停下升级。
2. **② 接线实现**（Flash，按纪律 V 拆步，每步独立核验）：
   - ②a：改 `internal/console/cli.go` 的 `executeRun` 分支——在路由到 `ExecuteAgentRunnerRun` / `ExecuteSupervisedRun` 之前插入 **readiness 判定**（permit 消费点按 design-r1 的时序放置，可能不在同一位置）
   - ②b：改 `internal/console/runtime_agent.go` 的装配入口——在装配开始处调用判定、消费 permit（file:line 级改动，跑聚焦测试）
   - ②c：改 `internal/console/runtime_supervised.go` 的装配入口（同 ②b 逻辑，file:line 级改动，跑聚焦测试）
   - ②d：补测试（见 ③）
   - **每步完成即跑 `go build ./...` + `go vet ./...` + 聚焦测试**；任一步红即修，不进入下一步
   - 若某步发现需要改 `internal/gates/**` 判定本体 ⇒ 停下升级（边界外）
3. **③ 测试**（Flash）：
   - 正例：readiness 判定通过 + permit 未消费 ⇒ 进入装配；装配成功后 permit 已消费
   - 反例 1：readiness 判定失败 ⇒ 不进入装配（fail-closed）
   - 反例 2：permit 已被消费 ⇒ 拒绝装配（不双启动）
   - 反例 3：判定成功但 CLI 未传必要参数 ⇒ fail-closed（不静默跳过）
   - 反例 4（回归）：既有 P44 判定测试不变
   - 每条测试绑定「删掉哪个守卫必须变红」的变异检验目标
4. **④ 门禁**：`go build ./...` + `go vet ./...` + `go test ./...` + `gate.js --all --scope=index / tree / ci` + `gate.js --selftest` 全绿；`GOOS=linux go build ./...` + `GOOS=linux go vet ./...` 干净（Unix 零影响的本机证据）
5. **⑤ 预审**（V4 Pro，1+1）：核对接线点是否与 design-r1 一致、反例是否成立、有无越界改 `internal/gates/**` / `internal/chain/**` 的**判定本体**（import 导出符号允许）；报 Medium+ ⇒ 修复后原位重跑 ⑤
6. **⑥ 独立扫描**（MAI，3 次额度）：按 §7.2.1 三条件判定（本片改 .go 语义 ⇒ ⑥ 必）；报 Medium+ ⇒ 修复后重跑 ⑥
7. **⑦ 终审**（Codex，1+2）：按 §7.3 只给切片定义、diff、契约原句、只读约束；首轮不附任何清单；复核接线是否改变了 P44 判定的语义、是否引入新的 fail-open 路径、Unix 零影响是否有证据
8. **⑧ 文档装配**（documenter）：按【五】清单逐项，CN/EN 严格对称
9. **⑨ 原生验证**：三作用域 + SELFTEST + `go build/vet/test` 全绿（读数为 ⑩ 回执的前置）
10. **⑧c 收尾**：`tmp/P69-T027B-E2E-WIRING/` 逐件精确清理，回读五条读数
11. **⑩ 停点**：停在 ⑩ 等授权；不自行 commit/push

【五、文档同步更新（本片范围；CN/EN 严格对称）】

1. **新建报告对** `docs/validation/t027b-e2e-wiring{,_EN}.md`（核心交付物：§1 执行摘要 / §2 接线点清单 / §3 调用顺序与失败路径 / §4 反例与变异检验 / §5 审查链结论 / §6 门禁读数 / §7 已知边界 / §8 明确未执行 / `artifacts` 围栏 / `sixstage` 转录）
2. **新建证据包** `docs/validation/evidence/P69-T027B-E2E-WIRING/`（`00-briefing.md` / `01-arch/design-r1.md` / `02-batch/` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）
3. **`docs/ARCHITECTURE{,_EN}.md`** —— 在 §3 的「console 多任务串行装配（P50）」段**之后**，**加法式**补一段「P44 readiness / permit 门在装配入口的位置」；不改既有描述
4. **`docs/DEV_PLAN{,_EN}.md`** —— ① P44 行（`:181` 附近）的「端到端接线为后续切片」改为「端到端接线已由 P69 完成」；② T027-B 里程碑段（`:176`–`:185`）**加法式**加 P69 一行（报告路径）
5. **`docs/validation/p44-gate-readiness{,_EN}.md`** —— **加法式**在文末或「已知边界」段加一行指针：「端到端接线已由 P69 `T027B-E2E-WIRING` 完成（报告：`docs/validation/t027b-e2e-wiring{,_EN}.md`）」；**不改原行**（§11.1 历史不改写）
6. **`docs/validation/p57-closeout-assessment{,_EN}.md`** —— §3.3 的 R-1 行**加法式**加指针（「R-1 已由 P69 收口」）；**不改原行**
7. **`docs/validation/s1-reach-reachability{,_EN}.md`** —— **① 阶段核实** §2 的 G-D 行是否含 R-1 / P44 接线登记；**有**则加法式加指针；**无**则在 ⑩ 回执标注不适用
8. **`docs/t027/REMAINING_SLICES{,_EN}.md`** —— P69 段状态回写（🔄 → ✅）+ 完成记录行（报告路径，按 §11.2 不写 hash/run）
9. **`docs/ADR_REGISTER{,_EN}.md`** —— **ADR-043**（主控 V4 Pro 的登记，四要素 + 「不得作为先例」；按 §7.8.4 的 OB-37 元规则）；**不登记 §12.4**
10. **`docs/CONTRACTS{,_EN}.md`** —— **§4.2 末句时态更新（措辞，非语义）**：现「`internal/console/**` 未被改动，故 C7 / N4 的端到端接线尚未完成，接线属后续切片」→ 改为「C7 / N4 的端到端接线已由 P69 `T027B-E2E-WIRING` 完成（报告：`docs/validation/t027b-e2e-wiring{,_EN}.md`）；`internal/console/**` 的接线点见该报告 §2」；**其余条款一字不改**（若 ① 阶段判定需改语义 ⇒ 停下升级）
11. **`docs/t027/CHECKLIST_v55{,_EN}.md`** —— **不适用**（快照冻结，v56 生成时继承）
12. **`docs/DELIVERY_DIRECTIVE{,_EN}.md`** —— **不适用**（业务切片，非治理片；本片不改准则正文；ADR-043 登记在 ADR_REGISTER，不在此）
13. **`docs/SECURITY{,_EN}.md`** —— **不适用**（不改变安全边界；readiness/permit 的 fail-closed 语义已在 P44 落 CONTRACTS §4.2）
14. **`docs/TEST_STRATEGY{,_EN}.md`** —— **不适用**（不触 AT 条目；P44 判定已有测试）
15. **`docs/OPERATIONS{,_EN}.md`** —— **不适用**（不触用户操作）
16. **`docs/PRODUCT_REQUIREMENTS{,_EN}.md`** —— **不适用**（不触 REQ）
17. **`docs/S1_SUPPORT_MATRIX{,_EN}.md`** —— **不适用**（不触平台清单）
18. **`docs/validation/s1-exit{,_EN}.md`** —— **不适用**（P69 不直接闭合 S1 Exit 八门；S-G8 仍依赖 T027-B 整体推进）

> 第 11–18 项在 ⑩ 回执中逐项标注「不适用 + 理由」；不得省略。

【六、⑩ 回执要求】

1. 派单表（`agentName` 逐字机械读取 + 机械读取脚本 stdout 原文）
2. ① 五件摘要（现状核对 + 接线设计 + 失败路径 + 回退方案 + 反例清单，每项一句话 + 落点）
3. ② 分步（②a–②d）每步的 `git diff --stat` 与聚焦测试读数
4. ③ 测试结果（正例 + 4 反例 + 变异检验逐条红/绿）
5. ⑤ / ⑥ / ⑦ 结论摘要（⑦ 计量 = 终审 N + 复审 M + 盲审 K；盲审不抽须列理由）
6. 三作用域 + SELFTEST + `go build/vet/test` 原始行（含 `GOOS=linux` 两条）
7. 文档同步 `git diff --stat`（逐项覆盖【五】18 项；不适用项列理由）
8. ⑧c 后五条终态读数（v1.67 §5.1「⑧c 后取证窗口」子款）
9. 明确未执行事项：未改准则正文 / 判据本体 / `CONTRACTS` 语义 / schema / fixtures / CI / v50–v55 / `internal/chain/**` / `internal/adapters/**` / `internal/gates/**` 判定本体；未替用户裁定；未 commit/push；未启动新切片
10. 本片自身越界登记：§1.7.2 三条 bullet 越界如实登记；未发生写「无」
11. 待授权项：提交授权；push 授权
12. 给 v56 的事实包（v56 生成时用：P69 闭环 hash + SELFTEST 读数 + P69 接线点的架构回写位置）

【七、执行节奏与红线】

1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束
3. 临时件一律 `tmp/P69-T027B-E2E-WIRING/`；只增不删；⑧c 统一清理；**CI 日志与门禁捕获文件也落这里，不落 `tmp/` 根**
4. **不改准则正文**：`DELIVERY_DIRECTIVE{,_EN}.md` 完全不改；ADR-043 登记在 `docs/ADR_REGISTER{,_EN}.md`（**不登记 §12.4**）
5. **不改产品代码边界外**：本片代码只改 `internal/console/**`；不改 `internal/chain/**`、`internal/adapters/**`、`internal/gates/**` 判定本体、`internal/tickets/**`、`internal/guard/**`、`schemas/**`、`testdata/contracts/**`、`.github/workflows/**`、`go.mod`
6. **不改 v50–v55 历史快照**
7. **CN/EN 严格对称**：每步改动两侧逐字对照；G3(a) 六形状全等
8. **不替用户裁定**
9. **不启动新切片**
10. **纪律 V**：② 必须按 ②a–②d 小步串行，每步独立核验；遇超时先查 `git status` 半截产物
11. **临时件 D 计数**（§5.1）：临时件违规按 v1.67 §5.1 三段处置

现在开始执行 P69 步骤 0。

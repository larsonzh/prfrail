# P88 起跑消息原文（用户 2026-11-25 点名；逐字落盘）

## P88 起跑 · GATES-PERF-ASSESS · 用户 2026-11-25 点名

---

切片：P88 GATES-PERF-ASSESS 门禁耗时构成评估与轻量效率优化（S 档 · 工具类 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）

你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片是 P88，S 档工具片
`[SLICE]`，片名 = `GATES-PERF-ASSESS`。目标 = 门禁耗时构成实测 + 轻量效率优化
（去重/缓存/并行），**判据判定结论逐字不变**。步骤链首位 = ① 契约影响判定
（预计不触 CONTRACTS；是否升准则 MINOR 由判定定）。档位启动前设定、不得中途
升档（§8.3）。**不预设临时授权**——任何额度用尽按 §7.8.3 停下报用户。

**关键锚点（步骤 0 须实测核对，以仓库实况为准）**：
- HEAD / origin = `fed0f5c`（实况）；v71 快照 §一表值为 `2782a09`（生成时点，
  有「落盘后 HEAD 前移」注，非错误）
- 准则 = v1.72；SELFTEST = 744/744；三作用域 = TOTAL_FAIL=0
- 滚动窗口 v69 / v70 / v71 = 209 / 203 / 201（v71 自身行数由 v72 复核）
- 机械读数 ADR 48 / OB 144 / DR 18 / 试点顺延 25

---

〇 切片标识
- 编号 P88；标识 GATES-PERF-ASSESS；分级 [SLICE]；类型 = 工具片（评估 + 轻量
  优化）；规模 = S（纯 tools/gates/** + 报告对；不触 internal/**）；
- 步骤链：0 → ① 影响判定 → ② 实现（评估脚本 / 优化落点）→ ③ 一致性（优化
  前后三作用域输出逐字 diff 为空）→ ④ 门禁 → ⑤ → ⑥（按 §7.2 判定：触判据
  文件即必做）→ ⑦ 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash (deepseek) High（默认；② 实现若遇判据语义复杂度
  或连续 ≥2 次漂移 ⇒ 停下报用户升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（四笔，均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 1 复审（§7.5「S 档 1+1」字面）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 ⑦ 分档表判据：触门禁
     判据执行路径（tools/gates 实现，不触判据本体）⇒ Codex；不触 OS/安全
     语义，不触发深度档）；
  c. 盲审不抽（不触身份准入 / 所有权 / 停机语义）；
  d. 探针 0。

〇.5 权威输入（逐字执行；§1.2 全字段为本片唯一落地依据）

### P88 — `GATES-PERF-ASSESS`（门禁耗时构成评估 + 轻量效率优化 · 工具片）· ⬜ 未立项

| 字段 | 内容 |
|---|---|
| 标识 | `GATES-PERF-ASSESS`（台账编号 = **P88**；来源 = P27 报告 §7「候选 OB（性能）」+ v70 快照层 3 登记位 + 用户 2026-11-25 排片裁定；**P87 闭环后启动、须用户点名**） |
| 状态 | ⬜ 未立项（起跑须用户点名；起跑消息须四要素留痕） |
| 目标 | ① 实测并落盘门禁耗时构成（SELFTEST、三作用域 index/tree/ci、base 三件套 gofmt/build/vet/test 的执行次数与耗时占比）；② 依据构成落**轻量效率优化**（候选：base 三件套跨作用域去重、可复用读数缓存、无依赖步骤并行化）——**判据判定结论逐字不变**；③ 输出优化前后对照证据与后续大改建议（如需大重构另立片，不在本片实施） |
| 依赖 | P87 闭环（`DIRECTIVE-REVIEW-MODEL-POLICY`）；无其它前置；G5-b 闭集 v1.72 已含 11 枚（基线） |
| 规模 | **S**（评估 + 轻量优化；纯 `tools/gates/**` + 报告对；不触 `internal/**`） |
| 步骤 | 0 → ① 影响判定（预计不触 CONTRACTS；是否升准则 MINOR 由判定定）→ ② 实现（评估脚本 / 优化落点）→ ③ 一致性（优化前后三作用域输出逐字 diff 为空）→ ④ 门禁 → ⑤ → ⑥（按 §7.2 判定：触判据文件即必做）→ ⑦ 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | 门禁耗时构成无实测基线；base 三件套在三作用域重复执行次数未量化；无 §12.4 条目（仅 v70 快照层 3 与台账登记） |
| 交付物 | ① 评估报告（耗时构成表 + 重复次数清单 + 命令与退出码）；② 轻量优化落地（去重/缓存/并行，附逐字一致性校验）；③ 优化前后对照证据（耗时读数 + 输出 diff 为空）；④ 报告对 `docs/validation/gates-perf-assess{,_EN}.md`；⑤ 台账段与完成记录；⑥ 证据包；⑦ 大改建议清单（如需，移交后续片） |
| 验收证据 | ① 优化前后三作用域输出**逐字一致**（判据结论不变，diff 为空）；② 耗时对比实测（同机同仓库前后各一组读数 + 命令 + 退出码）；③ SELFTEST 744/744 不变；④ 三作用域 `TOTAL_FAIL=0`；⑤ 报告对双语严格镜像 |
| 边界 | 不新增/不改任何判据的**判定逻辑**（只动执行调度/缓存/去重）；不触 §6.3 判据本体；不改 `internal/**`、CONTRACTS/schema/fixtures/workflows、CI workflow；不升准则 MAJOR；不改 G5-b 闭集与白名单；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P88-GATES-PERF-ASSESS/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P88-GATES-PERF-ASSESS/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = S 档 ⇒ §7.5「S 档 1+1」；档位 = `GPT-5.3-Codex` 常规档（§5.2 分档表判据：触门禁判据执行路径（tools/gates 实现，不触判据本体）⇒ Codex；不触 OS/安全语义，不触发深度档）；盲审不抽（不触身份准入语义）；批量辅助 GLM ≥1（耗时数据整理 / 双语镜像） |
| 起跑条件 | 须用户显式点名；四要素留痕；P87 台账正式完成记录与 v71 快照已闭环 |

---

一 背景与目标
- 背景：门禁耗时构成无实测基线（P27 候选 OB 未登记）；base 三件套在三作用域
  重复执行次数未量化；v71 快照层 3 已把本片升格为 P88 台账段。
- 目标：见权威输入「目标」字段（① 实测基线；② 轻量优化且判据结论逐字不变；
  ③ 对照证据 + 大改建议）。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek)（影响判定 + 优化落点
  方案 + 耗时构成评估方案；Sol 不参与双打）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek)（优化前后输出逐字
  一致性比对 + 夹具回归）；
- ⑤ 预审：DeepSeek V4 Pro (deepseek)，额度 1+1（触判据实现，预计不可跳过）；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，① 判定触判据文件即必做；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+1（S 档字面）；
- 盲审：不抽（不触身份准入 / 所有权 / 停机语义）；
- 批量辅助：GLM-5.3-Flash (glm) ≥1 次、严格串行、载体 deep-reasoner
  （禁 quick-researcher）；仅限四类任务，禁择优；派单写明字符上限；
- 审查输入包消毒（§7.3）；计量（R7.1）以调用回执为单一真值源；探针 0；
- 风险提醒：R 达 2 ⇒ §9.2 高风险标记并停下升级。

三 环境声明
- Ubuntu VM **不开**、网络代理**不开**（耗时评估在本机实测；无外网依赖；
  不做 gitee）；
- ⑨ 验证 = 本机门禁三作用域 + CI 双腿（Gates + selftest）；提交前 go
  build/vet/test ./... 全绿（§13.1.3；JS 侧 node 语法 + selftest 自测）。

四 步骤 0 必做（§13.1）
- 锚点核对，出处 = docs/t027/CHECKLIST_v71{,_EN}.md 的 §一（7 项）：
  ① HEAD = fed0f5c（实况）；② origin/main = fed0f5c；③ 准则 = v1.72；
  ④ SELFTEST = 744/744；⑤ 三作用域 = 全 TOTAL_FAIL=0；⑥ v69 行数 = 209（LF）；
  ⑦ v70 行数 = 203（LF）；另核滚动窗口 v69/v70/v71 中 v71 = 201（LF）与机械
  读数 ADR 48 / OB 144 / DR 18 / 试点顺延 25。≥2 项不符 ⇒ 以仓库实况为准，
  先报终态读数再继续；
- **台账 P88 段已登记**（v71 收尾轮已落 ⬜ 未立项）⇒ 步骤 0 翻「执行中」并跑
  门禁（index/tree）确认 TOTAL_FAIL=0；
- R/D = 0 基线；临时件落点 = tmp/P88-GATES-PERF-ASSESS/
  （⑧c 前只增不删；片内 D 累计达 3 ⇒ 强制停片，§5.1 v1.71）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md（§13.1.10，
  落 docs/validation/evidence/P88-GATES-PERF-ASSESS/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v71.md 及 _EN（§四 P88 队列项）；
2. docs/t027/REMAINING_SLICES{,_EN}.md 的 P88 段（§1.2 全字段，:2365 起）；
3. docs/DELIVERY_DIRECTIVE.md 的 §7.5（S 档 1+1 口径）、§5.2（⑦ 分档表）、
   §13.1 / §13.2（起跑与派单纪律）、§6.3（G5-b 现状 11 枚）；
4. tools/gates/ 的 gate.js、selftest.js、lib/ctx.js 与 lib/criteria/** 现状
   （优化对象与判据判定边界）；
5. docs/validation/directive-review-model-policy{,_EN}.md（P87 报告对，G5-b
   闭集与 SELFTEST 744 基线出处）；
6. docs/CODING_CONVENTIONS.md（.js = 无 BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS；是否升准则 MINOR（改判据实现文件与否）由
  判定定，判据判定逻辑零改动；
- 耗时构成实测方案：同机同仓库、每组读数附命令与退出码；测量对象 = SELFTEST、
  三作用域、base 三件套（含跨作用域重复次数）；
- 优化落点候选：base 三件套跨作用域去重、可复用读数缓存、无依赖步骤并行化——
  每一项附「优化前后输出逐字一致」校验；
- 变异/回归：优化后 SELFTEST 744/744 不变、三作用域 TOTAL_FAIL=0；
- 大改建议清单单独成节（不在本片实施）。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. tools/gates/**：仅执行调度/缓存/去重（具体文件由 ① 方案书点名），判据判定
   逻辑零改动；
2. docs/DELIVERY_DIRECTIVE{,_EN}.md：若 ① 判定触判据实现文件 ⇒ 头版本 MINOR
   + §0.3 变更行 + 附录 C.0az 台账块（判据本体零改动）；若不触则不动；
3. docs/validation/gates-perf-assess{,_EN}.md（报告对，双语严格镜像）；
4. docs/t027/REMAINING_SLICES{,_EN}.md：P88 段翻「✅ 闭环」+ 完成记录；
5. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记；
6. 证据包 docs/validation/evidence/P88-GATES-PERF-ASSESS/（含 00-briefing、
   耗时基线读数、优化前后 diff 为空证据、审查原文逐字）。
核对项（由用户裁定）：
7. 是否同步 v71 快照层 3 的 `GATES-PERF-ASSESS` 残留登记（v71 已声明
   「已升格为 P88 台账段」，无需改动——核验即可）。
不动：
8. CHECKLIST_v71{,_EN}（v72 快照再吸收）；DEV_PLAN{,_EN}（工具片，不在 T027
   主线，不回写）；s1-exit / p57 门表；CONTRACTS/schema/fixtures/workflows；
   internal/**；CI workflow；G5-b 闭集与白名单；ADR 不新增（除非 ① 判定）；
   DOCUMENTATION_PLAN{,_EN}（v72 生成时再换代）。

八 边界
- 只做耗时构成评估 + 轻量优化；不新增/不改任何判据的判定逻辑；不触 §6.3 判据
  本体；不改 internal/**、CONTRACTS/schema/fixtures/workflows、CI workflow；
  不升准则 MAJOR；不改 G5-b 闭集与白名单；不做 gitee；不预设临时授权。

九 回执要求
- 首行 P88: STARTED；四要素（裁决人=用户、日期=2026-11-25、理由=门禁优化排片
  裁定、依据=本起跑消息 + v71 快照 §四）；
- 锚点 7 项 + 模型/档位（主控 Flash High / ① V4 Pro / ⑦ Codex 常规档 1+1）/
  盲审不抽 / 批量辅助计划 / VM 与代理均不开逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 优化前后输出逐字一致（diff 为空）为硬验收，不符即回退优化；
- 单次输出 ≤ 10,000 字符（§13.1.16），超限分批、禁止自重启；
- 提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。

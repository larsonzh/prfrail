# P87 起跑消息原文（用户 2026-11-24 点名；逐字落盘）

## P87 起跑 · DIRECTIVE-REVIEW-MODEL-POLICY · 用户 2026-11-24 点名

---

## 【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）

你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片是 P87，M 档治理片
`[SLICE]`，片名 = `DIRECTIVE-REVIEW-MODEL-POLICY`。目标 = 把 ⑦ 评审模型选用
规则写成常设选型矩阵（四档）+ 固定判级门槛共识 + 新增 §5.4「② 加强轮触发与
回退」。步骤链首位 = ① 契约影响判定（预计不触 CONTRACTS；触准则本体；判定并
登记 ADR-048）。盲审必抽 1 次（终审前执行；模型名单触身份准入语义）。档位启动
前设定、不得中途升档（§8.3）。**不预设临时授权**——任何额度用尽按 §7.8.3 停下
报用户。

**关键锚点（步骤 0 须实测核对，以仓库实况为准）**：
- HEAD / origin = `010dd52`（实况）；v70 快照 §一表值为 `ef1e1f2`（生成时点，
  有「落盘后 HEAD 前移」注，非错误）
- 准则 = v1.71；SELFTEST = 740/740；三作用域 = TOTAL_FAIL=0
- 滚动窗口 v68 / v69 = 218 / 210（v70 自身行数由 v71 复核）
- 机械读数 ADR 47 / OB 144 / DR 18 / 试点顺延 25

---

〇 切片标识
- 编号 P87；标识 DIRECTIVE-REVIEW-MODEL-POLICY；分级 [SLICE]；类型 = 治理片；
  规模 = M（§2.2/§5.2/§5.4/§6.3/§7.5/§8.2 条款 + tools/gates 判据实现与自测）；
- 步骤链：0 → ① 契约影响判定 → ② 实现 → ③ 一致性 → ④ 门禁 → ⑤ 预审 →
  ⑥ 独立扫描 → ⑦ 盲审 + 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩
  （盲审在终审前执行）；
- 主控 = DeepSeek V4.1 Flash (deepseek) High（默认；② 实现若遇判据语义复杂度
  或连续 ≥2 次漂移 ⇒ 停下报用户升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（五笔，均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 3 复审（治理片 §7.5「L 档/治理片/高风险」字面）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 ⑦ 分档表判据：触规则
     语义 / 门禁判据 / 角色授权 ⇒ Codex）；额度 1 终审 + 3 复审，第 2、3 次
     复审须收窄输入、档位不得中途升档（§7.5 v1.35）；
  c. 盲审必抽 1（另计 +1：模型名单触身份准入语义）；
  d. Neon 交叉复核触发式：本片不预配，⑦ 终审报结构疑点/判级分歧 ⇒ 停下报
     用户再裁；
  e. 本片探针 0、不跑 Hydrogen 试水（留后续业务片 §1.7.4 临时授权）。

〇.5 权威输入（逐字执行；§1.2 全字段与 §5.4 条款为本片唯一落地依据）

### P87 — `DIRECTIVE-REVIEW-MODEL-POLICY`（⑦ 评审模型选型矩阵 + 判级标准共识 + ② 加强轮触发 · 治理片）· ⬜ 未立项

| 字段 | 内容 |
|---|---|
| 标识 | `DIRECTIVE-REVIEW-MODEL-POLICY`（台账编号 = **P87**；来源 = 用户 2026-11-23 系列裁定：Sol/Carbon 对比、Gemini 全系评估、Grok 终裁不纳入、② 加强轮触发判据、P85 四读矩阵（OB-141/142 留痕）+ 用户 2026-11-24 选型矩阵三处再修订裁定（成本实测：GPT 族显著低于 Gemini 族）；**P85 收口后启动、须用户点名**） |
| 状态 | ⬜ 未立项（起跑须用户点名；起跑消息须四要素留痕） |
| 目标 | ① 把 ⑦ 评审模型选用规则写成**常设选型矩阵**（任务类型 × 风险 × 成本 × 上下文四维度），并**固定判级门槛共识**（「什么算 Medium」的判级哲学：fail-closed 兜底 vs 契约完备性——消除跨模型判级漂移，四读实证同输入 PASS 与 FINDINGS 并存）；模型池四档：**常规档** = `GPT-5.3-Codex` Extra High 默认 + `GPT-6.2 Sol` Max 性价比备选（预算紧张轮）；**深度档** = `GPT-6.2 Sol` Max 主读 + `Gemini 4 Neon` Extra High **触发式交叉复核**（Sol 报 Medium+ 或结构疑点 ⇒ 追加 Neon 深挖；最高风险轮必配 Neon）；**大上下文文档档** = `Gemini 4 Carbon` High（仅 >400K 输入包 + 纯文档/S 档）；**轻量档** = `Gemini 4 Hydrogen` Extra High（文档型/治理片）；② 修订 §5.2「同环不混用」为组合触发条款（组合 = 主读 + 交叉复核双腿，不预设条数）；③ 新增 §5.4「② 加强轮触发与回退」（S1–S4 机械判据 + 灵敏度联动注记）；④ ⑤ 预审第二备用可选改 Hydrogen（B-part，用户另行裁定）；⑤ 生效机制：此后主控生成快照直接引用本矩阵，起跑消息按矩阵声明 ⑦ 选档与理由（深水轮须同时声明主读腿与交叉复核腿；沿用档位启动前设定、不得中途升档） |
| 依赖 | P85 闭环（`SECURITY-PROOF-SUBSYSTEM`）；**前置试水已由 P85 四读矩阵完成**（Carbon/Neon/Sol/Codex 同输入各 1 次，数据包 = `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/` + 对照表 `07-three-read-comparison.md`，OB-140/141/142/143 留痕）；**Hydrogen 对照评审 ≥1 次仍为前置**（后续业务片按 §1.7.4 临时授权，或并入本片 B-part 由用户另行裁定）；Grok 不纳入（用户终裁） |
| 规模 | **M**（治理片：§2.2/§5.2/§5.4/§6.3/§7.5/§8.2 条款 + `tools/gates` 判据实现与自测） |
| 步骤 | 0 → ① 契约影响判定（预计不触 CONTRACTS；触准则本体）→ ② 实现（准则条款 + g5b.js）→ ③ 一致性 → ④ 门禁 → ⑤ → ⑥ → ⑦ 盲审（必抽：模型名单触身份准入语义）+ 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | §2.2/§5.2 无 Gemini 系任何成员；G5-b 闭集无 `gemini-4-*`；§7.5 单价表无 Carbon/Neon/Hydrogen 行；§5.2「同环不混用」无组合例外；无「② 加强轮」触发与回退条款；无判级门槛定义条款 |
| 交付物 | ① §2.2 角色表增「独立终审员（深度档备选，两级）」与「独立终审员（轻量档）」行（限定名 `Gemini 4 Neon (copilot)` / `Gemini 4 Carbon (copilot)` / `Gemini 4 Hydrogen (copilot)`，modelId `gemini-4-neon/carbon/hydrogen`）；② §5.2 选型矩阵（四档如上）+ 判级门槛共识句 + 组合触发条款（修订「同环不混用」为「默认不混用；组合 = 主读 + 交叉复核双腿，仅限硬门/高风险轮次且经用户同轮授权，计数按两轮」）+ **适用域排除句**（Codex 与 Carbon **不得单独承担**触句柄/ACL/并发/安全/协议语义的深水轮终审——P85 四读实证两者同输入均漏已验证真实 Medium；深水轮 = Sol Max 主读 + Neon 触发式交叉复核）；③ 新增 §5.4（② 加强轮：S1 同根因连红 2 次 / S2 重发 ≥3 / S3 ⑤ 报实现类 Medium+ / S4 ⑦ 报实现类 High；触发即停下报用户；单轮单包授权、完成回退 Flash；选档 Codex 优先、Hydrogen 按预算；灵敏度联动注记：S3/S4 读数来源为 ⑤/⑦ 模型，温和判级模型低报 ⇒ 触发灵敏度下降，深水/安全轮 ⑦ 主读腿须为深度档）；④ §7.5 单价表增三行（Carbon 200/1000+Cache10/0、Neon 300/1500+Cache15/0、Hydrogen 150/750+Cache7.5/0）与备选额度口径（备选档 ⑦ 额度与深度档同 1+3；轻量档同常规档）+ **上下文档位计价注记**（Sol 1M 档 400/1500 翻倍；Carbon/Hydrogen 仅 1M 档单价即实价；Neon 200K/1M 两档；Codex 实采 400K）+ **四读实测消耗注记**（Sol 67.9 / Codex 129.9 / Carbon 266.8 / Neon 493.4 credits，用时 14m50s / 26m58s / 40m22s / 21m3s——GPT 族实测显著低于 Gemini 族（同输入 n=1）；实测消耗主变量 = 模型输出/思考行为而非单价，**该比例不得硬编码**，单价表仅作单价真值源）；⑤ §8.2 档位映射（Carbon=High 固定、Neon/Hydrogen=Extra High）+ 取档授权句；⑥ §6.3 G5-b 闭集 + `g5b.js` `CLOSED_SET` + selftest 夹具同步；⑦ 附录 C 台账块 + 版本 MINOR；⑧ 报告对 `docs/validation/directive-review-model-policy{,_EN}.md`；⑨ 台账段与完成记录；⑩ 证据包 |
| 验收证据 | ① G5-b 闭集三处逐字一致（§2.2 modelId 列 / §6.3 闭集 / `g5b.js` `CLOSED_SET`）+ 新夹具红绿对（含 `gemini-4-*` 三枚合法 + 一枚非法）；② 选型矩阵条款、判级门槛共识句与组合触发条款字面可判（对照表逐行）；③ §5.4 的 S1–S4 全部绑定既有机械读数来源（R 计数/重发计数/⑤⑥⑦ 结论行）且灵敏度联动注记在位；④ 三作用域 `TOTAL_FAIL=0`、SELFTEST 计数更新；⑤ 报告对双语严格镜像 |
| 边界 | 不新增本草案清单以外的任何模型行（模型池以本节清单为闭集，弃用/不采用的模型不写准则）；不改 ⑤⑥ 现役首选（MAI Code 2 Flash 维持 ⑥；B-part 除外）；不改深度档 `GPT-6.2 Sol` Max 首选地位与其判据行；不触 CONTRACTS/schema/fixtures/workflows；不升准则 MAJOR；不动 `internal/**` 产品代码；Grok 4.7 不纳入；P85 在飞条款不追改 |
| 分级 | `[SLICE]` |
| 探针额度 | 0（模型限定名与定价以 operator UI 实读为准，不计探针；试水对照评审已在 P85 四读矩阵完成、Hydrogen 试水在前置业务片以 §1.7.4 临时授权完成，均不占本片额度） |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P87-DIRECTIVE-REVIEW-MODEL-POLICY/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = 治理片 ⇒ §7.5「L 档 / 治理片 / 高风险」1 终审 + 3 复审；档位 = `GPT-5.3-Codex` 常规档（**§5.2 ⑦ 分档表判据：触规则语义/门禁判据/角色授权 ⇒ Codex**）；盲审必抽 1 次（另计 +1：模型名单触身份准入语义）；**是否附加 Neon 交叉复核按触发式口径，立项时由用户裁定**（P85 四读实证 Codex 对契约完备性类缺陷敏感度低）；批量辅助 GLM ≥2（镜像核对/闭集三处一致核对） |
| 起跑条件 | 须用户显式点名；四要素留痕；前置试水数据点已闭环（P85 四读矩阵，OB-140/141/142/143）；Hydrogen 试水按依赖字段执行 |

### 5.4 ② 加强轮触发与回退（2026-11-23 用户裁定立项）

**触发判据（四选一即触发；全部绑定既有机械读数，不得以主观评价替代）**：

| 编号 | 判据 | 机械读数来源 |
|---|---|---|
| S1 | 同一实现单**同根因连红 2 次**（④ 门禁或聚焦测试同一失败断言串；根因按 §9.2 同根因合并口径判定） | ④ 门禁日志 + 测试失败归一化串 |
| S2 | ② 子代理对同一任务包**重发 ≥3 次**（截断 / 异常 / 格式偏离） | §7.5 R7.1 调用回执重发计数 |
| S3 | ⑤ 预审报 **Medium+** 且根因归类为「实现/设计缺陷」而非「规格歧义」 | ⑤ 结论行 + 发现归因 |
| S4 | ⑦ 复审报 **High** 且定位为实现细节（非契约/语义类） | ⑦ 结论行 + 发现定位 |

**灵敏度联动注记**：S3/S4 的读数来源是 ⑤/⑦ 位模型，温和判级模型（如 `Gemini 4 Carbon` High，仅限其大上下文文档档场景）低报 ⇒ 触发灵敏度下降（P85 四读实证）；深水/安全轮的 ⑦ 主读腿须为深度档（`GPT-6.2 Sol` Max，见 §5.2 选型矩阵），保证触发灵敏度。

**前置声明（可选）**：起跑消息可对「L 档 + 并发/崩溃窗口 + 跨包改动」的片预声明「本片触发 S1–S4 即按 §5.4 授权加强轮候选」；已预声明则触发时无需二次停下；未预声明则触发必须停下报用户。

**流程**：
1. 触发 → 主控**停下**（不得静默切换 §13.1.8、不得中途升档 §8.3）→ 报用户；
2. 用户裁定 → 授权 ② 加强轮承担**该实现单**（单轮/单包、§1.7.4 五要素留痕、不得作先例）；
3. 选档 = `GPT-5.3-Codex` Extra High（首选）或 `Gemini 4 Hydrogen` Extra High（按输入包大小与预算水位，见 §5.2 选型矩阵）；
4. 该单完成（④ 门禁 + 聚焦测试绿）即**回退** `DeepSeek V4.1 Flash` 常态，授权不延续。

**留痕格式**：每次触发与授权逐次记入报告「异常处理」章：{判据编号、机械读数摘录、裁定人、日期、承担轮模型限定名、回退时点}；R 计数按 §9.2 正常累计（加强轮本身不豁免返工计数）。

---

一 背景与目标
- 背景：P85 四读矩阵（数据包 = P85 证据包 P87-four-read-data 15 件）实证
  「换模型 = 换判级标准」（同输入 PASS 与 FINDINGS 并存）、成本 GPT 族显著低于
  Gemini 族（n=1）；现行 §2.2 无 Gemini 系成员、G5-b 闭集无 gemini-4-*、
  §5.2「同环不混用」无组合例外、无判级门槛定义条款、无 ② 加强轮触发条款。
- 目标：见权威输入「目标」字段（① 矩阵 + 判级共识；② 组合触发；③ §5.4；
  ④ B-part 用户另行裁定；⑤ 生效机制）。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek)（契约影响判定 + 条款
  落地方案 + 跨包共享契约与失败路径枚举；Sol 不参与双打）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek)（G5-b 夹具红绿对 +
  变异检验）；
- ⑤ 预审：DeepSeek V4 Pro (deepseek)，额度 1+1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，必做（触 §6.3 判据语义
  不可跳过）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+3（见 〇-a/b）；
- 盲审：必抽 1 次（另计 +1、终审前执行）；
- 批量辅助：GLM-5.3-Flash (glm) ≥2 次、严格串行、载体 deep-reasoner
  （禁 quick-researcher）；仅限四类任务（枚举 / 摘录核对 / 双语镜像 / MANIFEST
  草案），禁择优；派单写明字符上限；
- 审查输入包消毒（§7.3）；计量（R7.1）以调用回执为单一真值源；探针 0；
- 风险提醒：R 达 2 ⇒ §9.2 高风险标记并停下升级。

三 环境声明
- Ubuntu VM **不开**、网络代理**不开**（纯文档 + JS 判据，无 Go 原生验证对象；
  双平台验证由 CI 双腿承担）；不做 gitee；
- ⑨ 验证 = 门禁三作用域 + CI 双腿（Gates + selftest）；提交前 go build/vet/test
  ./... 全绿（§13.1.3；JS 侧 node 语法 + selftest 自测）。

四 步骤 0 必做（§13.1）
- 锚点核对，出处 = docs/t027/CHECKLIST_v70{,_EN}.md 的 §一（7 项）：
  ① HEAD = 010dd52（实况）；② origin/main = 010dd52；③ 准则 = v1.71；
  ④ SELFTEST = 740/740；⑤ 三作用域 = 全 TOTAL_FAIL=0；⑥ v68 行数 = 218（LF）；
  ⑦ v69 行数 = 210（LF）；另核机械读数 ADR 47 / OB 144 / DR 18 / 试点顺延 25。
  ≥2 项不符 ⇒ 以仓库实况为准，先报终态读数再继续；
- **先登台账 P87 段**（权威输入 §1.2 全字段落 docs/t027/REMAINING_SLICES{,_EN}.md
  后续批次节）→ 再翻「执行中」→ 登台账后跑门禁（index/tree）确认 TOTAL_FAIL=0；
- R/D = 0 基线；临时件落点 = tmp/P87-DIRECTIVE-REVIEW-MODEL-POLICY/
  （⑧c 前只增不删；片内 D 累计达 3 ⇒ 强制停片，§5.1 v1.71）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md（§13.1.10，
  落 docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v70.md 及 _EN（§四 P87 队列项、§五 P87 四读数据指针）；
2. docs/DELIVERY_DIRECTIVE.md 的 §2.2 / §2.4 / §5.2 / §6.3 / §7.5 / §8.2 /
   §12.4（现行模型表、⑦ 三档分档判据、G5-b 闭集、单价表）；
3. docs/ADR_REGISTER.md 及 _EN（ADR-019 先例形态、ADR-047 max）；
4. docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/P87-four-read-data/ 的
   07-three-read-comparison.md + 四份模型报告（Carbon/Neon/Sol/Codex）+ r3 原文；
5. docs/validation/security-proof-subsystem.md 及 _EN（P85 报告对，四读数据出处）；
6. tools/gates/lib/g5b.js（现行 CLOSED_SET）+ tools/gates/selftest.js（相关夹具）；
7. docs/t027/REMAINING_SLICES{,_EN}.md（P87 段登记位置与完成记录体例）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 契约影响判定：预计不触 CONTRACTS/schema/fixtures/workflows；触身份准入语义 ⇒
  登记 ADR-048（参照 ADR-019）；
- §2.2 角色表增 Neon/Carbon/Hydrogen 三行 + Sol 常规档性价比备选句；
- §5.2 选型矩阵四档 + 判级门槛共识句 + 组合触发条款 + 适用域排除句（逐字按
  权威输入交付物②）；
- §5.4 S1–S4 机械判据 + 灵敏度联动注记（逐字按权威输入）；
- §7.5 单价表三行 + 上下文档位计价注记 + 四读实测消耗注记 + 备选额度口径；
- §8.2 档位映射（Carbon=High、Neon/Hydrogen=Extra High）+ 取档授权句；
- G5-b 闭集三处逐字一致（§2.2 modelId 列 / §6.3 闭集 / g5b.js CLOSED_SET）+
  selftest 夹具（gemini-4-* 三枚合法 + 一枚非法红绿对）；
- 双语镜像（G3 全量）；判级共识句、组合触发条款、适用域排除句字面可判。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. docs/DELIVERY_DIRECTIVE{,_EN}.md：§2.2 / §5.2 / §5.4 / §6.3 / §7.5 / §8.2 +
   头版本 MINOR + 附录 C.0 台账块 + §12.4 OB 行（如产生额度例外）；
2. tools/gates/lib/g5b.js CLOSED_SET + tools/gates/selftest.js 夹具；
3. docs/ADR_REGISTER{,_EN}.md：ADR-048；
4. docs/validation/directive-review-model-policy{,_EN}.md（报告对）；
5. docs/t027/REMAINING_SLICES{,_EN}.md：P87 段（步骤 0 登）+ 完成记录；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记；
7. 证据包 docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/。
核对项（由用户裁定）：
8. 是否同步 P85 报告对/证据包的模型档案措辞（四读结论引述）。
不动：
9. CHECKLIST_v70{,_EN}（v71 快照再吸收）；DEV_PLAN{,_EN}（P25 裁定：治理片
   不在 T027 主线，不回写）；s1-exit / p57 门表；CONTRACTS/schema/fixtures/
   workflows；internal/** 产品代码；CI；DOCUMENTATION_PLAN{,_EN}（v71 生成时
   再换代）。

八 边界
- 只落选型矩阵 + 判级共识 + §5.4 + 计价/档位条款；不新增本草案清单以外的任何
  模型行（模型池以清单为闭集）；不改 ⑤⑥ 现役首选（MAI Code 2 Flash 维持 ⑥；
  B-part 除外）；不改深度档 Sol Max 首选地位与其判据行；不触 CONTRACTS/schema/
  fixtures/workflows；不升准则 MAJOR；不动 internal/** 产品代码；Grok 4.7 不
  纳入；P85 在飞条款不追改；不预设临时授权。

九 回执要求
- 首行 P87: STARTED；四要素（裁决人=用户、日期=2026-11-24、理由=四读数据齐备
  与草案定稿、依据=本起跑消息 + v70 快照 §四）；
- 锚点 7 项 + 模型/档位（主控 Flash High / ① V4 Pro / ⑦ Codex 常规档 1+3）/
  盲审必抽 / 批量辅助计划 / VM 与代理均不开逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16），超限分批、禁止自重启；
- 提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。


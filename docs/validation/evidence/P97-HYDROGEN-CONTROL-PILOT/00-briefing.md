# P97 起跑消息（原文逐字落盘，§13.1.10）

以下为用户 2026-11-28 起跑消息原文，逐字保留（含排版与标点）。

---

切片：P97 HYDROGEN-CONTROL-PILOT 轻量档对照评审实测（S 档 · docs-only · 补足 P87 前置缺口）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P97，S 档
docs-only 评估/实测片，片名 = `HYDROGEN-CONTROL-PILOT`。目标 = 用
`Gemini 4 Hydrogen`（⑦ 轻量档）对 P87 输入包做**对照评审实测**，补足
P87 报告 §7 与台账 P87 段的「Hydrogen 前置缺口」（条款已落地、对照
评审实测未完成），回填注记并登记 OB 留痕。**不改任何规则语义、不改
§5.2 矩阵 / §6.3 判据本体 / G5-b 闭集；只做实测、记录与回填**。档位
启动前设定、不得中途升档（§8.3）。不预设临时授权——额度用尽按
§7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `8f0c42b`（实况）；准则 = v1.81；SELFTEST = 806/806；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v78 / v79 / v80 = 211 / 209 / 210（v80 自身行数由 v81 复核）
- 机械读数 ADR 48 / OB 145 / DR 18 / 试点顺延 25（OB-146 本片登记）
- tmp 断言顶层两件：.gitkeep / gates-cache（其下 base 与 g4b）
- **HEAD 特别说明**：v80 快照 §一表的 `90932ee` 是 v80 生成时点锚点；
  其后 v80 快照提交 `8f0c42b` 已落盘，HEAD 前移至 `8f0c42b`——以
  `git rev-parse --short HEAD` 实况为准，`8f0c42b` 与表内 `90932ee`
  不同属预期，不判锚点不符。

〇 切片标识
- 编号 P97；标识 HYDROGEN-CONTROL-PILOT；分级 [SLICE]；类型 =
  docs-only 评估/实测片（补足 P87 Hydrogen 前置缺口）；规模 = S
  （只读 P87 证据包 + 一轮 Hydrogen 试水 + 对照记录 + 注记回填 +
  §12.4 OB 登记）；
- 步骤链：0 → ① 方案（对照对象 / 对照口径 / OB-146 登记文本 / 升版
  判定；**停点 1**）→ ②③ 对照记录文档（prfrail-documenter，文档型
  退化 §13.1.9）→ ④ 门禁 → ⑤ 预审 → ⑥（跳过，经 ⑦ 确认）→
  ⑦ 试水轮（Hydrogen，§1.7.4 额度外）+ 终审轮（V4 Pro）→ ⑧a
  报告对 + 回填 → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash (deepseek) **High**（默认档；若 ① 方案
  书 >10,000 字须提示用户可升 V4 Pro，**不得静默切换**）；
- **停点 1 闸（硬）**：① 产出四项——对照对象与输入包消毒口径、
  对照面定义（结论 / 发现数 / 严重度 / 误判四面对照，只记录不重裁）、
  OB-146 登记文本、准则升版判定——主控在 ① 停下时报用户裁定，
  **不得自行决定**；
- ⑦ 档位 = **已定**：试水轮 = `Gemini 4 Hydrogen (copilot)`
  （§1.7.4 临时授权，额度外，OB-146 留痕）+ 常规终审
  `DeepSeek V4 Pro (deepseek)` **S 档 1 终审 + 1 复审（1+1，§7.5）**
  （依据 = 预算策略「纯文档且不触 §7.2.1 ③ 语义 ⇒ V4 Pro」）；
- 用户裁定（均不得作先例；停点 1 四项待你回场裁定）：
  a. 停点 1 四项 = 待裁定（对照对象/口径、OB-146 文本、升版判定）；
  b. Hydrogen 试水 = 本片目的本身，用户点名即授权执行 1 轮（限定名
     `Gemini 4 Hydrogen (copilot)`，载体 independent-reviewer，§1.7.4
     五要素 + OB-146 独立留痕，额度外）；
  c. 盲审不抽（不触身份准入 / 所有权 / 停机）；探针 0；
  d. 批量辅助 GLM-5.3-Flash 实派 ≥2 次（modelId `glm-5.3-flash`，
     与 §2.2 批量辅助首选一致；**禁用 `glm-5.3`**——那是 §2.2 的
     ⑤ 预审备用档，两型号不得混用，先例 = OB-83）；严格串行、载体
     deep-reasoner、禁 quick-researcher，仅限耗时数据整理 / 双语
     镜像，派单写明字符上限。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P97 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P97 — `HYDROGEN-CONTROL-PILOT`（⑦ 轻量档对照评审实测 · docs-only）
| 字段 | 内容 |
|---|---|
| 标识 | `HYDROGEN-CONTROL-PILOT`（台账编号 = P97；来源 = P87 报告 §7 前置缺口注记 + 台账 P87 段「前置缺口登记」+ v80 快照 §五 + 用户 2026-11-28 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 用 `Gemini 4 Hydrogen`（⑦ 轻量档）对 P87 输入包做对照评审实测（§1.7.4 临时授权 + OB 留痕）；② 产出对照记录（Hydrogen 结论 vs P87 原 ⑦ 结论四面对照，只记录不重裁）；③ 回填 P87 报告 §7 注记与台账 P87 段「已补足」；④ §12.4 登记 OB-146；⑤ 报告对与证据包 |
| 依赖 | v80 快照闭环；P87 证据包可得（`docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/`，含 00-briefing / input-list / 01-arch / 03-verification / 04-reviews）；无其它前置 |
| 规模 | S |
| 步骤 | 0 → ①（停点 1）→ ②③ 对照记录文档 → ④ → ⑤ → ⑥ 跳过 → ⑦ 试水轮 + 终审轮 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | P87 轻量档（`Gemini 4 Hydrogen`）条款已落地、对照评审实测未完成，不得表述为已验证 |
| 交付物 | ① 方案书（对照对象 / 口径 / OB-146 文本 / 升版判定）；② Hydrogen 试水原文（逐字归档）；③ 对照记录表；④ 报告对 `docs/validation/hydrogen-control-pilot{,_EN}.md`；⑤ P87 报告 §7 注记回填 + 台账 P87 段翻「已补足」+ P97 段与完成记录；⑥ §12.4 OB-146 行；⑦ 证据包 |
| 验收证据 | ① Hydrogen 实跑对照记录逐字归档（MANIFEST/SHA256SUMS 自洽）；② 对照表四面对照（结论 / 发现数 / 严重度 / 误判）与 P87 原 ⑦ 结论可比；③ 注记回填双语严格镜像；④ 门禁 tree/index/ci TOTAL_FAIL=0 + CI 双腿 success；⑤ SELFTEST 806/806 不变（不触判据代码） |
| 边界 | 不改 §2.2 / §2.4 / §5.2 矩阵 / §6.3 判据本体与 G5-b 闭集；不改任何规则语义（只登记 OB + 回填注记）；不改 tools/gates/** / internal/** / ci.yml / CONTRACTS / schema / fixtures / workflows；不新增 ADR；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P97-HYDROGEN-CONTROL-PILOT/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = 试水轮 `Gemini 4 Hydrogen (copilot)` 1 轮（§1.7.4 临时授权、额度外、OB-146 留痕）+ 常规终审 `DeepSeek V4 Pro` **1 终审 + 1 复审（S 档 1+1，§7.5）**；⑥ 跳过（不触判据实现文件，经 ⑦ 确认，P72 先例）；盲审不抽；批量辅助 GLM-5.3-Flash ≥2（modelId `glm-5.3-flash`，禁 `glm-5.3`） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v80 快照闭环 |

一 背景与目标
- 背景：P87（DIRECTIVE-REVIEW-MODEL-POLICY，v1.72）落地 ⑦ 四档
  选型矩阵时，轻量档 `Gemini 4 Hydrogen` 条款已入 §5.2 / §8.2 与
  G5-b 闭集，但「对照评审实测」未完成——登记位 = P87 报告 §7 注记
  「Hydrogen 对照评审前置未完成…待后续业务片按 §1.7.4 临时授权
  补足；不得表述为已验证」与台账 P87 段「前置缺口登记（2026-11-24）」。
  本片为该缺口的专项补足片（用户 2026-11-28 点名）。
- 目标：见权威输入「目标」字段。
- **停点 1 闸为本片安全前提**：① 判定后停下报用户，不得自行决定
  对照对象/口径、OB-146 文本与升版判定。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek)（方案书：
  对照对象确认 / 对照面定义 / 输入包消毒 / OB-146 登记文本 / 升版
  判定）；
- ②③ 对照记录文档：prfrail-documenter + DeepSeek V4.1 Flash
  （文档型切片退化，§13.1.9；②③ 共担，派单逐次传 agentName +
  限定名 model）；
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek)；
- ⑥ 独立扫描：**跳过**（不触 gate.js / selftest.js 等判据实现文件；
  经 ⑦ 确认，P72 先例）；
- ⑦ 试水轮：主控派单 runSubagent(agentName='independent-reviewer',
  model='Gemini 4 Hydrogen (copilot)')，**§1.7.4 临时授权位**（五要素
  + OB-146 留痕，额度外），输入 = P87 输入包消毒件，任务 = 以轻量档
  身份对 P87 交付物做独立评审；
- ⑦ 常规终审：DeepSeek V4 Pro (deepseek) **1 终审 + 1 复审（S 档
  1+1，§7.5）**；盲审不抽；批量辅助 GLM-5.3-Flash ≥2（modelId
  `glm-5.3-flash`，见裁定 d）；审查输入包消毒（§7.3）；计量以调用
  回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 三作用域
  `--all --skip-base` 全绿 + SELFTEST 806/806 不变；CI 双腿按 §11.2
  观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v80{,_EN}.md §一）：HEAD
  实况 `8f0c42b`（表值 `90932ee` 为生成时点锚点，前移属预期，见
  〇·零 HEAD 特别说明）；准则 v1.81；SELFTEST 806/806；三作用域 0；
  v78/v79 = 211/209（v80 = 210 自身留 v81）；机械读数 ADR 48 /
  OB 145 / DR 18 / 顺延 25；tmp 断言两件（.gitkeep +
  gates-cache{base,g4b}）。≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P97 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），
  再翻「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P97-HYDROGEN-CONTROL-PILOT/
  （⑧c 前只增不删；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2
  逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v80.md 及 _EN（§一 锚点 / §四 下一项 = 无 /
   §五 观察项的 P87 Hydrogen 前置缺口行）；
2. docs/validation/directive-review-model-policy{,_EN}.md §7 注记
   与报告结构（回填对象）；
3. docs/DELIVERY_DIRECTIVE.md §1.7.4（临时授权五要素）/ §5.2 四档
   选型矩阵轻量档行 / §7.5（⑦ 额度分层，S 档 1+1）/ §12.4 尾部
   （OB-146 登记位）/ §2.2（Hydrogen 限定名与 modelId）/ §13.1.10 /
   §13.2.2；
4. docs/t027/REMAINING_SLICES{,_EN}.md P87 段（「前置缺口登记」
   行）与 §12.4 OB 表尾部；
5. docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/
   的 MANIFEST.md（对照输入包目录结构）与 04-reviews/（P87 原 ⑦
   结论基准）；
6. docs/CODING_CONVENTIONS.md（.md = BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 对照对象：首选 P87 自身输入包（`P87-DIRECTIVE-REVIEW-MODEL-POLICY/`
  证据包，含 00-briefing / 01-arch / 03-verification / 04-reviews，
  已确认可得）；备选 P96 证据包（① 判定）；
- 对照面：结论（PASS/FINDINGS）/ 发现数 / 严重度分布 / 误判（Hydrogen
  对已知 P87 结论的偏离面）——**只记录、不重裁**，不因对照结果修改
  P87 任何已闭环结论；
- OB-146 登记文本（§1.7.4 五要素 + 首状态词）：裁决人 = 用户、日期 =
  2026-11-28、理由 = 补足 P87 Hydrogen 前置缺口（专项对照实测）、
  依据 = 本起跑消息裁定 b、不得作先例；
- 升版判定：预计准则升 MINOR **v1.82**（§0.3 新行 + §12.4 OB-146 +
  头部）；若 ① 判定「仅登记与注记回填、不升版」，以 ① 方案书为准交
  用户裁定；
- **停点 1 触发即停**：报四项候选取用户裁定，不等额用完；
- ⑦ 试水输入包消毒（§7.3）：只投 P87 交付物与基准结论，不含本片
  任何结论；试水原文逐字归档（§13.1.13 同口径，Section A–E 全量，
  原件缺节如实保留缺口说明）；
- 对照记录与注记回填均为 docs-only，不触任何判据代码；SELFTEST
  806/806 必须不变。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. docs/validation/directive-review-model-policy{,_EN}.md §7 注记：
   「Hydrogen 对照评审前置未完成」→「已补足（P97，2026-11-28，对照
   结论见 `hydrogen-control-pilot{,_EN}.md`）」；
2. docs/t027/REMAINING_SLICES{,_EN}.md：P87 段「前置缺口登记」→
   「已补足」+ P97 段新增 + 完成记录；
3. docs/DELIVERY_DIRECTIVE{,_EN}.md：§12.4 登记 OB-146（若升版：
   头部 v1.82 + §0.3 新行，六形状双语镜像）；
4. docs/validation/hydrogen-control-pilot{,_EN}.md（报告对，新增，
   双语严格镜像）；
5. 证据包 docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/
   （00-briefing、freedom-list、input-list、试水原文、对照记录、
   MANIFEST/SHA256SUMS）。
不动：
6. CHECKLIST_v80{,_EN}（v81 再吸收）；tools/gates/**；internal/**；
   ci.yml；CONTRACTS / schema / fixtures / workflows；G5-b 闭集与
   模型白名单；ADR_REGISTER；s1-exit / p57 门表。

八 边界
- 不改 §2.2 / §2.4 / §5.2 矩阵 / §6.3 判据本体与 G5-b 闭集；不改
  任何规则语义（只做实测记录、OB 登记与注记回填）；不改判据代码
  与产品代码；不预设临时授权；不做 gitee；对照结论不得表述为
  「轻量档已验证」之外的扩展（回填注记只写「已补足」+ 报告引用）。

九 回执要求
- 首行 P97: STARTED；四要素（裁决人=用户、日期=2026-11-28、理由=
  补足 P87 Hydrogen 前置缺口（专项对照实测收口）、依据=本起跑消息 +
  P87 报告 §7 注记 + v80 快照 §五）；
- 锚点 + 模型/档位（主控 Flash High / ①⑤ V4 Pro / ②③
  prfrail-documenter 文档退化 / ⑥ 跳过经 ⑦ 确认 / ⑦ = Hydrogen
  试水 §1.7.4 额度外 + V4 Pro 常规终审 S 档 1+1 / GLM-5.3-Flash
  ≥2）逐项复述；
- **停点 1 触发即停**：报四项候选取用户裁定，不得自行决定；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。

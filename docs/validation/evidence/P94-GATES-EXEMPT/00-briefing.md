切片：P94 GATES-EXEMPT 判据豁免与白名单门槛合并治理（M 档 · 治理片 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P94，M 档治理片
`[SLICE]`，片名 = `GATES-EXEMPT`。目标 = 合并评估并落地：P93 移交的
(a) 通用「引用豁免登记位」与 (c) 白名单登记门槛 / 断言政策，加上
GATES-EXEMPT 本体的判据豁免机制泛化（P33 ①-3(iv) 拆片遗产）。**不改任何
判据的判定逻辑本体。** 步骤链首位 = ① 影响判定。档位启动前设定、不得
中途升档（§8.3）。不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `9217915`（实况）；准则 = v1.76；SELFTEST = 757/757；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v75 / v76 / v77 = 215 / 216 / 217（v77 自身行数由 v78 复核）
- 机械读数 ADR 48 / OB 145 / DR 18 / 试点顺延 25
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P94；标识 GATES-EXEMPT；分级 [SLICE]；类型 = 治理片（判据豁免
  政策 + 白名单门槛）；规模 = M（触最多 4 个判据文件 + selftest 断言 +
  白名单口径 + 准则正文）；
- 步骤链：0 → ① 影响判定 → ② 实现 → ③ 一致性（既有输入回归逐字不变 +
  新豁免形态自测）→ ④ 门禁 → ⑤ → ⑥（必做）→ ⑦ 终审 + 复审 → ⑧a →
  ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash High（默认；升级闸：① 若判定须动任何判据
  判定逻辑本体 ⇒ 停下报用户升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 2 复审（同 P93 裁定款；触多判据豁免条款语义）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 分档表判据：触门禁
     判据条款与断言政策——不触判据判定逻辑本体、不触 OS/安全/并发语义
     ⇒ 不触发深度档；Sol 不出场、不双打）；
  c. 盲审不抽（不触身份准入 / 模型白名单 / G5-b 闭集 / 所有权 / 停机）；
  d. 探针 0；
  e. 批量辅助 GLM-5.3-Flash 实派 ≥2 次（§13.1.7 字面；严格串行、载体
     deep-reasoner、禁 quick-researcher，仅限双语条款镜像 / 政策对照，
     派单写明字符上限）。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P94 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P94 — `GATES-EXEMPT`（判据豁免与白名单门槛合并治理 · 治理片）
| 字段 | 内容 |
|---|---|
| 标识 | `GATES-EXEMPT`（台账编号 = P94；来源 = 层 3 挂起项（v28 起；P33 ①-3(iv) 拆片）+ P93 移交 a/c + v77 §四/§六 下一裁定项 + 用户 2026-11-26 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 合并评估三件：GATES-EXEMPT 本体的判据豁免机制泛化（① 追溯 P33 ①-3(iv) 原始意图定边界）；P93 移交 (a) 通用「引用豁免登记位」（含登记文件自指与 T23 口径）；P93 移交 (c) 白名单登记门槛与断言政策；② 论证通过的项落地（通道落点由 ① 判定 + ⑦ 背书）；③ 既有判据输入回归逐字不变 |
| 依赖 | v77 快照闭环；P93 闭环（移交出处）；无其它前置 |
| 规模 | M |
| 步骤 | 0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | GATES-EXEMPT 无 §12.4 条目（仅层 3 登记位）；白名单登记门槛未定义（P92/P93 两次踩「登记即 743/744」）；转录归一登记位缺位（P93 通道 b 只覆盖 sixstage 围栏） |
| 交付物 | ① 三件合并评估结论（逐项：做/不做/依据）；② 获批项落盘（准则条款 + 判据/断言实现 + 白名单口径三处一致）；③ 版本包（预计 MINOR v1.77，由 ① 判定）；④ 报告对 `docs/validation/gates-exempt{,_EN}.md`；⑤ 台账段与完成记录；⑥ 证据包 |
| 验收证据 | ① 既有输入回归逐字不变（SELFTEST 全 PASS——总数若随断言政策变化，须与 ① 方案书一致并台账记录）；② 三作用域 TOTAL_FAIL=0；③ 新豁免/登记形态自测通过；④ 报告对双语严格镜像 |
| 边界 | 判据判定逻辑本体零改动（只增豁免通道/登记位/门槛条款）；不改 G5-b 闭集与模型白名单；不改 ci.yml / internal/** / CONTRACTS / schema / fixtures / workflows；ADR 不新增（不触身份准入）；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P94-GATES-EXEMPT/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P94-GATES-EXEMPT/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = 1+2（用户裁定 a）；档位 = `GPT-5.3-Codex` 常规档；盲审不抽；批量辅助 GLM ≥2（用户裁定 e）；⑥ 必做（触判据实现） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v77 快照闭环 |

一 背景与目标
- 背景：GATES-EXEMPT 自 v28 起挂层 3（P33 ①-3(iv) 拆片、触 4 个判据
  文件、无 §12.4 条目）；P93 通道 b 只解决 sixstage 围栏豁免，未选的
  (a)/(c) 移交本片；白名单登记门槛缺位已致 P92/P93 两次踩坑（登记即
  `SELFTEST` 743/744）。
- 目标：见权威输入「目标」字段。
- 合并评估为硬前提：① 逐项给「做 / 不做 / 依据」，不做项进大改建议
  清单，不得默认全做。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro（影响判定 + 三件合并评估 +
  P33 遗产追溯 + 通道落点设计 + 版本包判定）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（既有回归逐字比对 +
  新形态自测 + 三作用域）；
- ⑤ 预审：DeepSeek V4 Pro，额度 1+2；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，**必做**（触判据
  实现）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+2（用户裁定 a）；
- 盲审：不抽；批量辅助：GLM-5.3-Flash ≥2 次（见裁定 e）；审查输入包
  消毒（§7.3）；计量以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下
  升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域 + go build/vet/
  test ./... 全绿；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v77{,_EN}.md §一）：HEAD = origin =
  `9217915`；准则 v1.76；SELFTEST 757/757；三作用域 0；v75/v76 = 215/216；
  机械读数与 tmp 断言三件。≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P94 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻
  「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P94-GATES-EXEMPT/（⑧c 前只增不删；D
  累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P94-GATES-EXEMPT/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v77.md 及 _EN（§四 下一项 / §五 P93 移交登记 /
   §六）；
2. docs/DELIVERY_DIRECTIVE.md §6.3（判据条款全集）、§11.3、§0.2、§0.3、
   附录 C（C.0bd 起按实况）、§12.4、§7.5 / §5.2 / §13.1 / §13.2；
3. tools/gates/lib/criteria/ 全部判据实现（① 点名最多 4 个文件）+
   tools/gates/selftest.js（T23 断言）+ tools/gates/g6-whitelist.txt
   （零条目现状）；
4. docs/validation/g6-4-transcript-exempt{,_EN}.md（P93 报告对，大改建议
   段 = a/c 移交出处）；
5. GATES-EXEMPT 出处追溯：v28 起层 3 挂起位 + P33 ①-3(iv) 拆片记录
   （全仓 grep `GATES-EXEMPT` / `P33` / `①-3` 定位，含 §12.4 对应行）；
6. docs/CODING_CONVENTIONS.md（.js/.txt = 无 BOM + LF；.md = BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** /
  ci.yml；是否升 MINOR v1.77 + 附录 C.0bd 由判定定；触 4 个判据文件
  名单由 ① 逐一点名；
- 三件合并评估（逐项：做 / 不做 / 依据）：
  a. GATES-EXEMPT 本体（P33 ①-3(iv) 原始意图追溯 + 边界重定）；
  b. P93 移交 (a) 通用引用豁免登记位（含登记文件自指与 T23 口径）；
  c. P93 移交 (c) 白名单登记门槛与断言政策（解决「登记即 743/744」）；
- ③ 一致性：既有输入回归逐字不变（SELFTEST 全 PASS，总数变化须与 ①
  方案书一致并台账记录）；新形态自测通过；
- 政策三处一致：准则条款 ↔ 判据/断言实现 ↔ 白名单口径，逐字核对。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. docs/DELIVERY_DIRECTIVE{,_EN}.md：豁免登记位 / 白名单门槛条款 + 头部
   版本（预计 v1.77）+ §0.3 行 + 附录 C.0bd（CN/EN 严格镜像）；
2. tools/gates/lib/criteria/** 获批文件 + tools/gates/selftest.js +
   tools/gates/g6-whitelist.txt——① 方案书点名；
3. docs/validation/gates-exempt{,_EN}.md（报告对，新增，双语严格镜像）；
4. docs/t027/REMAINING_SLICES{,_EN}.md：P94 段翻「✅ 闭环」+ 完成记录；
5. 证据包 docs/validation/evidence/P94-GATES-EXEMPT/（00-briefing、
   freedom-list、input-list、三件合并评估、新形态自测、审查原文）；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
7. CHECKLIST_v77{,_EN}（v78 再吸收）；ci.yml；internal/**；CONTRACTS /
   schema / fixtures / workflows；G5-b 闭集与模型白名单；ADR_REGISTER；
   DEV_PLAN{,_EN}；DOCUMENTATION_PLAN{,_EN}（v78 时代换）；s1-exit /
   p57 门表；find-quotes.js。

八 边界
- 判据判定逻辑本体零改动（只增豁免通道 / 登记位 / 门槛条款）；不改
  G5-b 闭集与模型白名单；不改 ci.yml、internal/**、CONTRACTS / schema /
  fixtures；不新增 ADR；不做 gitee；不预设临时授权。

九 回执要求
- 首行 P94: STARTED；四要素（裁决人=用户、日期=2026-11-26、理由=判据
  豁免与白名单门槛合并治理（P93 移交 a/c + 层 3 GATES-EXEMPT）、依据=
  本起跑消息 + v77 快照 §四/§六 + P93 报告大改建议段）；
- 锚点 + 模型/档位（主控 Flash High / ① V4 Pro / ⑥ 必做 / ⑦ Codex 常规
  档 1+2）/ 盲审不抽 / 探针 0 / VM 与代理不开 / GLM ≥2 逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 三件合并评估结论 + 既有输入回归逐字不变 + 新形态自测通过为硬验收；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，等
  同轮授权；仅 origin。

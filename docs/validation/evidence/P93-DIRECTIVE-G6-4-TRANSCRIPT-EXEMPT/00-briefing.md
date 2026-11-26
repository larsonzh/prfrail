# P93 起跑消息原文（用户 2026-11-26）

> 本文件为 §13.1.10 要求的「起跑消息原文」落盘件。除本头部三行与本段外，正文为用户消息原文，未作改写。

切片：P93 DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT G6-4 与逐字转录互斥治理（S 档 · 治理片 [SLICE]）

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P93，S 档治理片 [SLICE]，片名 = DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT。目标 = 为「G6-4 计数型字面量判据」与「§11.3 六阶段逐字转录」的互斥给出正式合规通道（三选一，① 判定 + ⑦ 背书），并同步 T23 断言与 g6-whitelist.txt 政策。不改任何判据的判定逻辑本体。步骤链首位 = ① 影响判定。档位启动前设定、不得中途升档（§8.3）。不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = fba4129（实况）；准则 = v1.75；SELFTEST = 744/744；三作用域 = TOTAL_FAIL=0
- 滚动窗口 v74 / v75 / v76 = 216 / 215 / 216（v76 自身行数由 v77 复核）
- 机械读数 ADR 48 / OB 145 / DR 18 / 试点顺延 25
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P93；标识 DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT；分级 [SLICE]；类型 = 治理片（判据条款豁免通道 + 断言政策同步）；规模 = S（纯准则 §6.3 / §11.3 + tools/gates/lib/criteria/g6.js + selftest.js 的 T23 政策 + g6-whitelist.txt 口径 + 报告对 + 台账 + 证据包）；
- 步骤链：0 → ① 影响判定 → ② 实现 → ③ 一致性（既有输入回归逐字不变 + 新豁免形态自测）→ ④ 门禁 → ⑤ → ⑥（必做）→ ⑦ 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash High（默认；升级闸：① 若判定须动 G6-4 判定逻辑本体（而非新增豁免通道）⇒ 停下报用户升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 2 复审（触判据条款语义与 T23 政策，介于 S 档 1+1 与治理片 1+3 之间，用户裁定 1+2）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 分档表判据：触门禁判据条款（G6-4 豁免通道）与断言政策——不触判据判定逻辑本体、不触 OS/安全/并发语义 ⇒ 不触发深度档；Sol 不出场、不双打）；
  c. 盲审不抽（不触身份准入 / 模型白名单 / G5-b 闭集 / 所有权 / 停机）；
  d. 探针 0；
  e. 批量辅助 GLM-5.3-Flash 实派 ≥2 次（§13.1.7 字面；严格串行、载体 deep-reasoner、禁 quick-researcher，仅限双语条款镜像 / 政策对照，派单写明字符上限）。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P93 段尚未登记 ⇒ 步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P93 — DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT（G6-4 转录豁免通道 · 治理片）
| 字段 | 内容 |
|---|---|
| 标识 | DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT（台账编号 = P93；来源 = P92 报告 §7 第 5 条候选 + v76 §四/§六 下一裁定项 + 用户 2026-11-26 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 在「G6-4 计数型字面量（\d{2,}\s*行 等形态）」与「§11.3 六阶段逐字转录」之间给正式通道，三选一：a) §11.3 增「转录归一登记位」；b) G6-4 对 sixstage 围栏显式豁免；c) T23/白名单政策修订；② 同步 T23 断言与白名单政策；③ 既有判据输入回归逐字不变 |
| 依赖 | v76 快照闭环；P92 闭环（冲突出处）；无其它前置 |
| 规模 | S |
| 步骤 | 0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | G6-4 与 §11.3 互斥无合规通道；g6-whitelist.txt 零实际条目（717 字节纯注释头）且 T23 钉死为空（登记即 743/744）；P92 只能最小归一 + 围栏外披露（用户裁定接受，不得作先例） |
| 交付物 | ① 三选一通道落盘（准则条款 + 判据/断言实现 + 白名单口径三处一致）；② 版本包（预计 MINOR v1.76，由 ① 判定）；③ 新形态自测（含计数/行数字面量的审查原文围栏通过 G6-4 且转录合规）；④ 报告对 docs/validation/g6-4-transcript-exempt{,_EN}.md；⑤ 台账段与完成记录；⑥ 证据包 |
| 验收证据 | ① 既有输入回归逐字不变（SELFTEST 全 PASS——总数若随 T23 政策变化，须与 ① 方案书一致并台账记录）；② 三作用域 TOTAL_FAIL=0；③ 新豁免形态自测通过；④ 报告对双语严格镜像 |
| 边界 | G6-4 判定逻辑本体不改（只增豁免通道/登记位）；不改 G5-b 闭集与模型白名单；不改 ci.yml / internal/** / CONTRACTS / schema / fixtures / workflows；ADR 不新增（不触身份准入）；不做 gitee |
| 分级 | [SLICE] |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/freedom-list.md |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | tmp/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = 1+2（用户裁定 a）；档位 = GPT-5.3-Codex 常规档；盲审不抽；批量辅助 GLM ≥2（用户裁定 e）；⑥ 必做（触判据实现） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v76 快照闭环 |

一 背景与目标
- 背景：P92 实测三方互斥——G6-4 判据（g6.js:104-105 的 \d{2,}\s*行 等形态）要求计数型字面量过期即红；§11.3 要求审查原文围栏逐字转录；两者同时满足即红（P92 ⑥ 原文含行数字面量）；g6-whitelist.txt 又被 T23 断言钉死为空（登记 6 条即 SELFTEST 743/744）。无合规通道 ⇒ 只能最小归一 + 披露。
- 目标：见权威输入「目标」字段（三选一 + 政策同步 + 既有回归）。
- 三选一由 ① 判定 + ⑦ 背书；判据判定逻辑本体零改动。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro（影响判定 + 三选一通道设计 + 版本包判定）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（既有回归逐字比对 + 新形态自测 + 三作用域）；
- ⑤ 预审：DeepSeek V4 Pro，额度 1+2；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，必做（触判据实现）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+2（用户裁定 a）；
- 盲审：不抽；批量辅助：GLM-5.3-Flash ≥2 次（见裁定 e）；审查输入包消毒（§7.3）；计量以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域 + go build/vet/test ./... 全绿；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对；台账 P93 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻「执行中」并跑门禁（index/tree）；
- R/D = 0 基线；临时件落点 = tmp/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/（⑧c 前只增不删；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md（§13.1.10）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v76.md 及 _EN（§四 下一项 / §五 治理候选 / §六）；
2. docs/DELIVERY_DIRECTIVE.md §6.3（G6-4 判据条款）、§11.3（报告结构与逐字转录条款）、§0.2（版本包三件）、§0.3、附录 C、§7.5 / §5.2 / §13.1 / §13.2；
3. tools/gates/lib/criteria/g6.js（G6-4 判据实现）、tools/gates/selftest.js（T23 断言）、tools/gates/g6-whitelist.txt；
4. docs/validation/gates-selftest-split{,_EN}.md §7 第 5 条（候选出处）；
5. docs/CODING_CONVENTIONS.md（.js/.txt = 无 BOM + LF；.md = BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** / ci.yml；是否升 MINOR v1.76 + 附录 C.0bc 由判定定；三选一通道的选取理由与副作用逐条书面；
- 三选一：a) 转录归一登记位 / b) G6-4 对 sixstage 围栏豁免 / c) T23 与白名单政策修订——① 选一 + ⑦ 背书；未选的两项进大改建议清单；
- ③ 一致性：既有输入回归逐字不变；构造新形态样例自测通过；
- 政策三处一致：准则条款 ↔ g6.js/selftest.js 实现 ↔ g6-whitelist.txt 口径，逐字核对。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：1. 准则 §6.3 G6-4 豁免/登记位条款 + §11.3 转录条款同步 + 头部版本（预计 v1.76）+ §0.3 行 + 附录 C.0bc（CN/EN 严格镜像）；2. g6.js + selftest.js（T23 政策）+ g6-whitelist.txt（口径注释）；3. 报告对（新增，双语严格镜像）；4. 台账 P93 段翻「✅ 闭环」+ 完成记录；5. 证据包；6. tools/gates/cjk-newwords.txt（如新增 CJK 字即登记）。
不动：7. CHECKLIST_v76{,_EN}（v77 再吸收）；ci.yml；internal/**；CONTRACTS / schema / fixtures / workflows；G5-b 闭集与模型白名单；ADR_REGISTER；DEV_PLAN{,_EN}；DOCUMENTATION_PLAN{,_EN}（v77 时代换）；s1-exit / p57 门表；find-quotes.js。

八 边界
- G6-4 判定逻辑本体零改动（只增豁免通道 / 登记位）；不改 G5-b 闭集与模型白名单；不改 ci.yml、internal/**、CONTRACTS / schema / fixtures；不新增 ADR；不做 gitee；不预设临时授权。

九 回执要求
- 首行 P93: STARTED；四要素（裁决人=用户、日期=2026-11-26、理由=判据/条款互斥治理（P92 冲突候选）、依据=本起跑消息 + v76 快照 §四/§六 + P92 报告 §7 第 5 条）；
- 锚点 + 模型/档位（主控 Flash High / ① V4 Pro / ⑥ 必做 / ⑦ Codex 常规档 1+2）/ 盲审不抽 / 探针 0 / VM 与代理不开 / GLM ≥2 逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 三选一通道结论 + 既有输入回归逐字不变 + 新形态自测通过为硬验收；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。

切片：P95 GATES-CRITERIA-ASYNC G6/G7/G8 异步化（M 档 · 工具片 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4 Pro，Max 档）。本片 = P95，M 档工具片
`[SLICE]`，片名 = `GATES-CRITERIA-ASYNC`。目标 = 复用 `runCriteriaConcurrent`
既有调度契约把 G6 / G7 / G8 改为真并发，**判据判定结论逐字不变、报告顺序
契约不变**。步骤链首位 = ① 影响判定 + 收益论证闸。档位启动前设定、不得
中途升档（§8.3；本环境 V4 Pro 仅 Max 档可设，中途不可改）。不预设临时
授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `22370fe`（实况）；准则 = v1.77；SELFTEST = 770/770；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v76 / v77 / v78 = 216 / 217 / 211（v78 自身行数由 v79 复核）
- 机械读数 ADR 48 / OB 145 / DR 18 / 试点顺延 25
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P95；标识 GATES-CRITERIA-ASYNC；分级 [SLICE]；类型 = 工具片（判据
  并发化）；规模 = M（触 tools/gates/lib/criteria/g6.js / g7.js / g8.js +
  gate.js 调度 / 报告顺序）；
- 步骤链：0 → ① 影响判定 + 收益论证闸 → ② 实现 → ③ 一致性（同输入 A/B
  sha256 逐字相等 + 报告顺序不变）→ ④ 门禁 → ⑤ → ⑥（必做）→ ⑦ 终审 +
  复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4 Pro (deepseek) **Max**（用户 2026-11-27 裁定：本片
  触底层并发语义 ⇒ V4 Pro；且 V4 Pro 仅 Max 档可设、中途不可改；不得作
  先例）；
- **⑦ 档位闸（硬）**：⑦ 档位由 ① 判定 + 用户裁定，预计属 §5.2 深度档
  （并发语义）；主控在 ① 停下时一并报「档位 + 额度」候选，**不得自行
  取档**；
- 用户裁定（均不得作先例；其中 ⑦ 一项待你回场裁定）：
  a. ⑦ 额度与档位 = 待裁定（预计深度档 Sol Max 主读 + Gemini 4 Neon 触发
     式交叉复核；金额高于常规档，须你显式点头）；
  b. 盲审不抽（不触身份准入 / 所有权 / 停机）；
  c. 探针 0；
  d. 批量辅助 GLM-5.3-Flash 实派 ≥2 次（modelId `glm-5.3-flash`，与
     §2.2 批量辅助首选一致；**禁用 `glm-5.3`**——那是 §2.2 的 ⑤ 预审
     备用档，两型号不得混用，先例 = OB-83）；严格串行、载体
     deep-reasoner、禁 quick-researcher，仅限耗时数据整理 / 双语镜像，
     派单写明字符上限。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P95 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P95 — `GATES-CRITERIA-ASYNC`（G6/G7/G8 异步化 · 工具片）
| 字段 | 内容 |
|---|---|
| 标识 | `GATES-CRITERIA-ASYNC`（台账编号 = P95；来源 = v72 层 3 性能候选 + P88 报告 §7 建议 4 + v78 层 3 记录态 + 用户 2026-11-27 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 收益论证闸：实测 `--all` 单进程内 G6 / G7 / G8 的串行耗时构成与并发收益上限，并评估共享缓存并发读与报告顺序契约的改造代价；② 论证通过则把 G6 / G7 / G8 改为 `runCriteriaConcurrent` 真并发；③ 判据判定结论逐字不变、报告顺序（`RUN_ORDER` 回填）不变 |
| 依赖 | v78 快照闭环；P88 报告 §7 建议 4（出处）；无其它前置 |
| 规模 | M |
| 步骤 | 0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | G6 / G7 / G8 仍同步串行；`runCriteriaConcurrent` 调度契约已铺好未启用；收益未实测 |
| 交付物 | ① 收益论证报告（串行构成 + 并发收益上限 + 代价）；② 获批并发化实现；③ A/B sha256 + 报告顺序不变对照证据；④ 报告对 `docs/validation/gates-criteria-async{,_EN}.md`；⑤ 台账段与完成记录；⑥ 证据包；⑦ 大改建议清单 |
| 验收证据 | ① 同输入 A/B 逐字一致（sha256 相同 + 比对器可证伪）+ 报告顺序逐字不变；② SELFTEST 全 PASS（总数若随断言政策变化，须与 ① 方案书一致并台账记录）；③ 三作用域 TOTAL_FAIL=0；④ 报告对双语严格镜像 |
| 边界 | 判据判定逻辑本体零改动（只改调度 / 并发方式）；不改 ci.yml / internal/** / CONTRACTS / schema / fixtures / workflows；G5-b 闭集与模型白名单不动；ADR 不新增；准则是否升 MINOR 由 ① 判定；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P95-GATES-CRITERIA-ASYNC/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ 档位与额度待裁定（见裁定 a；预计深度档）；盲审不抽；批量辅助 GLM-5.3-Flash ≥2（modelId `glm-5.3-flash`，禁 `glm-5.3`，见裁定 d）；⑥ 必做（触判据实现） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v78 快照闭环 |

一 背景与目标
- 背景：P88 报告 §7 建议 4——`runCriteriaConcurrent` 已铺好调度契约，把
  G6 / G7 / G8 改成异步即可取得真并发；涉及「并发读共享缓存 + 报告顺序
  契约」；P88 实测 G6 / G7 / G8 为最重判据。
- 目标：见权威输入「目标」字段。
- **⑦ 档位闸为本片安全前提**：① 判定后停下报用户，不得自行取 ⑦ 档。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek) Max（影响判定 +
  收益论证闸 + 并发契约评估）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（A/B 逐字 + 报告顺序 +
  三作用域）；
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek) Max（与主控同档）；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，**必做**（触判据
  实现）；
- ⑦ 独立终审：档位待裁定（见裁定 a）；盲审不抽；批量辅助
  GLM-5.3-Flash ≥2（modelId `glm-5.3-flash`，见裁定 d）；审查输入包
  消毒（§7.3）；计量以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下
  升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域 + go build/vet/
  test ./... 全绿；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v78{,_EN}.md §一）：HEAD = origin =
  `22370fe`；准则 v1.77；SELFTEST 770/770；三作用域 0；v76/v77 = 216/217；
  机械读数与 tmp 断言三件。≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P95 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻
  「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P95-GATES-CRITERIA-ASYNC/（⑧c 前只增
  不删；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v78.md 及 _EN（§四 下一项 = 无 / 层 3 性能候选行）；
2. docs/DELIVERY_DIRECTIVE.md §7.5 / §5.2（⑦ 分档表，深度档判据）/
   §2.2（批量辅助首选与 ⑤ 备用档型号）/ §13.1 / §13.2 / §0.2；
3. tools/gates/lib/criteria/g6.js、g7.js、g8.js、tools/gates/gate.js
   （`runCriteriaConcurrent` 调度契约现状，:140/:231）；
4. docs/validation/gates-perf-assess{,_EN}.md §7 建议 4（出处 + G6/G7/G8
   耗时基线）；
5. docs/CODING_CONVENTIONS.md（.js = 无 BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** /
  ci.yml；**是否触并发 / 顺序契约语义（预计触）⇒ ⑦ 档位闸：停下报用户
  裁定 ⑦ 档位与额度，不得自行取档**；是否升准则 MINOR 由判定定；
- 收益论证闸（逐项：收益 / 代价 / 做或不做 / 依据）：
  a. 串行构成：G6 / G7 / G8 在单进程 `--all` 内的实测耗时占比；
  b. 并发收益上限：并发后 `--all` 墙钟预期（受 CPU 核数与共享缓存争用
     约束）；
  c. 并发代价：共享缓存并发读的同步机制（锁 / 只读快照）、报告顺序
     （`RUN_ORDER` 回填）的保持方式；
- ③ 一致性：同输入 A/B sha256 相同 + 报告顺序逐字不变 + 比对器翻转 1
  字节可证伪；
- ⑦ 档位闸触发即停，不等额用完。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. tools/gates/lib/criteria/g6.js / g7.js / g8.js + tools/gates/gate.js
   （调度 / 报告顺序）——① 方案书点名；
2. docs/DELIVERY_DIRECTIVE{,_EN}.md：若 ① 判定触并发调度 / 报告顺序条款
   ⇒ 升 MINOR（预计 v1.78）+ §0.3 + 附录 C 下一空位（否则不动）；
3. docs/validation/gates-criteria-async{,_EN}.md（报告对，新增，双语
   严格镜像）；
4. docs/t027/REMAINING_SLICES{,_EN}.md：P95 段翻「✅ 闭环」+ 完成记录；
5. 证据包 docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/（00-briefing、
   freedom-list、input-list、收益论证报告、A/B 对照、审查原文）；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
7. CHECKLIST_v78{,_EN}（v79 再吸收）；ci.yml；internal/**；CONTRACTS /
   schema / fixtures / workflows；G5-b 闭集与模型白名单；ADR_REGISTER；
   DEV_PLAN{,_EN}；DOCUMENTATION_PLAN{,_EN}（v79 时代换）；s1-exit /
   p57 门表；find-quotes.js。

八 边界
- 判据判定逻辑本体零改动（只改调度 / 并发方式）；不改 ci.yml、
  internal/**、CONTRACTS / schema / fixtures；G5-b 闭集与模型白名单
  不动；不新增 ADR；不做 gitee；不预设临时授权。

九 回执要求
- 首行 P95: STARTED；四要素（裁决人=用户、日期=2026-11-27、理由=门禁
  性能优化续片（层 3 候选点名）、依据=本起跑消息 + v78 快照 §四 + 层 3
  性能候选行）；
- 锚点 + 模型/档位（主控 V4 Pro Max / ① V4 Pro Max / ⑤ V4 Pro Max /
  ⑥ 必做 / ⑦ 档位待裁定）/ 盲审不抽 / 探针 0 / VM 与代理不开 /
  GLM-5.3-Flash ≥2 逐项复述；
- **① 判定触发 ⑦ 档位闸 ⇒ 立即停下报用户，不得自行取档**；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- A/B 逐字一致 + 报告顺序不变为硬验收，不符即回退；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。

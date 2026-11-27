# P96 起跑消息原文（§13.1.10）

切片：`GATES-BASE-CACHE`（P96）。日期：2026-11-27。

本文件为起跑消息**逐字原文**（§〇·零 至 §九 全部段落）。补正说明：初次落盘为「原文摘录」，经用户 2026-11-27 引导补正后改为全文照录。

---

切片：P96 GATES-BASE-CACHE 跨作用域 base 结论缓存（M 档 · 工具片 [SLICE] · 条款先行）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P96，M 档工具片
`[SLICE]`，片名 = `GATES-BASE-CACHE`。目标 = ① 收益论证闸 + 键设计 +
条款草案（**停点 1 用户裁定**）→ ② 条款先行落地（准则升 MINOR）→ ③ 跨
作用域 base 结论缓存实现。**判据判定逻辑零改动、base 结论语义不变、
绝不缓存任何 G1–G9 判据结论**。档位启动前设定、不得中途升档（§8.3）。
不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `9f3aea9`（实况）；准则 = v1.77；SELFTEST = 770/770；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v77 / v78 / v79 = 217 / 211 / 209（v79 自身行数由 v80 复核）
- 机械读数 ADR 48 / OB 145 / DR 18 / 试点顺延 25
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P96；标识 GATES-BASE-CACHE；分级 [SLICE]；类型 = 工具片（条款
  先行：先立 §6.1 base 结论缓存条款，再动 gate.js 读路径）；规模 = M
  （触 DELIVERY_DIRECTIVE §6.1/§6.3 条款 + gate.js base 读路径 +
  selftest.js 断言）；
- 步骤链：0 → ① 影响判定 + 收益论证闸（含键设计 / 条款草案 / flake
  政策；**停点 1**）→ ② 条款落地 → ③ 实现 → ④ 门禁 → ⑤ → ⑥（必做）
  → ⑦ 终审 + 复审 → ⑧a → ⑨（两遍冷热）→ ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash (deepseek) **High**（默认档；按选型策略，
  若 ① 条款草案 >10,000 字或键设计涉复杂状态机，主控须提示用户可升
  V4 Pro，**不得静默切换**）；
- **停点 1 闸（硬）**：① 产出四项——收益论证结论（做/不做）、条款文本、
  flake 政策（test 结论纳入/排除两案对照）、分级（M 常规 1+2 / 高风险
  1+3）——主控在 ① 停下时报用户裁定，**不得自行决定**；
- ⑦ 档位 = **已定**：常规档 `GPT-5.3-Codex (copilot)`，**M 档 1 终审 +
  2 复审（1+2）**（依据 = 预算策略「触 §6.3 门禁语义 ⇒ Codex」+ §7.5
  M 档行）；若停点 1 判高风险 ⇒ 额度升 1+3，需用户裁定；
- 用户裁定（均不得作先例；停点 1 四项待你回场裁定）：
  a. 停点 1 四项 = 待裁定（做/不做、条款文本、flake 政策、分级）；
  b. 盲审不抽（不触身份准入 / 所有权 / 停机）；
  c. 探针 0；
  d. 批量辅助 GLM-5.3-Flash 实派 ≥2 次（modelId `glm-5.3-flash`，与
     §2.2 批量辅助首选一致；**禁用 `glm-5.3`**——那是 §2.2 的 ⑤ 预审
     备用档，两型号不得混用，先例 = OB-83）；严格串行、载体
     deep-reasoner、禁 quick-researcher，仅限耗时数据整理 / 双语镜像，
     派单写明字符上限。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P96 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P96 — `GATES-BASE-CACHE`（跨作用域 base 结论缓存 · 工具片）
| 字段 | 内容 |
|---|---|
| 标识 | `GATES-BASE-CACHE`（台账编号 = P96；来源 = P88 报告 §7 建议 2 + v79 快照层 3 唯一性能候选 + 用户 2026-11-27 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 收益论证闸：实测 base 四条孤立计时与三作用域重复支付成本 + 键设计草案 + 条款草案 + flake 政策两案对照；② 条款先行：§6.1 落 base 结论缓存条款（准则升 MINOR，预计 v1.78），判据代码不动、CI 绿后再动代码；③ 实现跨作用域 base 结论缓存（键 / 三反例 / fail-open）；④ 判据判定逻辑零改动、base 结论语义不变、**绝不缓存任何 G1–G9 判据结论** |
| 依赖 | v79 快照闭环；P88 报告 §7 建议 2（出处）；G4b-CACHE 条款模板（附录 C.0aj）；无其它前置 |
| 规模 | M |
| 步骤 | 0 → ①（停点 1）→ ② 条款落地 → ③ 实现 → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨（两遍冷热）→ ⑧b → ⑧c → ⑩ |
| 现状缺口 | base 四条（gofmt / build / vet / test）在 index / tree / ci 三作用域重复支付；结论级缓存不存在；§6.1 无 base 结论缓存条款 |
| 交付物 | ① 收益论证报告 + 键设计草案 + 条款草案 + flake 政策对照；② 条款落地（准则升 v1.78 或按停点 1 裁定）；③ 缓存实现 + selftest 新断言；④ 冷热 A/B 对照证据；⑤ 报告对 `docs/validation/gates-base-cache{,_EN}.md`；⑥ 台账段与完成记录；⑦ 证据包；⑧ 大改建议清单 |
| 验收证据 | ① 三反例必失效 + 读前重算键 + 异常 fail-open 回全量（断言可证伪）；② 冷热结论逐字一致（含 detail 与输出顺序）+ 比对器翻转 1 字节可证伪；③ SELFTEST 全 PASS（新增断言须与 ① 方案书一致并台账记录）；④ 三作用域 TOTAL_FAIL=0（两遍：冷填充 → 热命中）；⑤ 报告对双语严格镜像 |
| 边界 | 判据判定逻辑本体零改动（只动 base 读路径与缓存）；缓存只存 §6.1 四条 base 结论、绝不缓存任何 G1–G9 判据结论；不改 ci.yml / internal/** / CONTRACTS / schema / fixtures / workflows；G5-b 闭集与模型白名单不动；ADR 由 ① 判定（预计不新增，条款落点由 ① 定）；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P96-GATES-BASE-CACHE/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P96-GATES-BASE-CACHE/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = 常规档 `GPT-5.3-Codex (copilot)` **1 终审 + 2 复审（M 档 1+2，§7.5）**；停点 1 判高风险 ⇒ 1+3 待裁定；盲审不抽；批量辅助 GLM-5.3-Flash ≥2（modelId `glm-5.3-flash`，禁 `glm-5.3`）；⑥ 必做（触 gate.js / selftest.js 判据实现文件） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v79 快照闭环 |

一 背景与目标
- 背景：P88 报告 §7 建议 2——「跨作用域 base 复用（判定结果缓存）：base
  四件套是『工作区 + 工具链』的纯函数，可设计带失效键的跨进程缓存；但
  它与 §6.3 G4b『绝不缓存判定结论』同族敏感，**须先立条款再实现**。短期
  替代 = CI 的 Gates 步改 `--all --skip-base --scope=ci`（已由 P90 落地）」。
  基线 = OB-113（§12.4）：base 四条孤立计时 gofmt 0.43 / build 4.72 /
  vet 3.50 / test 11.35 s（合计 20.0 s）；P88/P92 后判据侧大降，base 占比
  需重测。条款模板 = 附录 C.0aj（G4b-CACHE：键成分 / 读前重算 / 三反例 /
  fail-open / 升版义务）。flake 史 = DR-2~DR-18（go test 为最高频 flake
  源）。
- 目标：见权威输入「目标」字段。
- **停点 1 闸为本片安全前提**：① 判定后停下报用户，不得自行决定
  做/不做、条款文本、flake 政策与分级。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek)（影响判定 +
  收益论证闸 + 键设计 + 条款草案 + flake 政策）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（冷热 A/B 逐字 +
  三反例变异 + 三作用域两遍）；
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek)；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，**必做**（触
  gate.js / selftest.js 判据实现文件）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档 **1 终审 + 2 复审（M 档
  1+2；停点 1 判高风险 ⇒ 1+3）**；盲审不抽；批量辅助 GLM-5.3-Flash ≥2
  （modelId `glm-5.3-flash`，见裁定 d）；审查输入包消毒（§7.3）；计量以
  调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域两遍（冷填充
  → 热命中）+ go build/vet/test ./... 全绿；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v79{,_EN}.md §一）：准则 v1.77；
  SELFTEST 770/770；三作用域 0；v77/v78 = 217/211（v79 = 209 自身留 v80）；
  机械读数与 tmp 断言三件。**HEAD 特别说明**：快照 §一表的 `bea554b` 是
  v79 生成时点锚点；其后 v79 快照提交 `76d2c86` 与 P95 注记回填 `9f3aea9`
  已落盘，HEAD 前移至 `9f3aea9`——**以 `git rev-parse --short HEAD` 实况
  为准，`9f3aea9` 与表内 `bea554b` 不同属预期，不判锚点不符**。其余项
  ≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P96 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻
  「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P96-GATES-BASE-CACHE/（⑧c 前只增不删；
  D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P96-GATES-BASE-CACHE/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v79.md 及 _EN（§四 下一项 = 无 / 层 3 性能候选行）；
2. docs/DELIVERY_DIRECTIVE.md §6.1（base 四条基础门禁）/ §6.3 G4b 行
   （缓存条款模板）/ 附录 C.0aj（G4b-CACHE 先例）/ §7.5（⑦ 分层额度，
   M 档 1+2）/ §5.2（⑦ 四档选型矩阵）/ §2.2（角色—载体—模型三元组与
   批量辅助首选）/ §13.1 / §13.2 / §0.2；
3. tools/gates/gate.js（base 四条读路径 `:92–:131` 与 `--skip-base` 解析
   `:184–:196`）+ tools/gates/lib/criteria/g4b.js（缓存实现先例：键 /
   失效 / fail-open）+ tools/gates/selftest.js（断言风格）；
4. docs/validation/gates-perf-assess{,_EN}.md §7 建议 2（出处）+ 准则
   §12.4 的 OB-113（base 计时基线）；
5. docs/CODING_CONVENTIONS.md（.js = 无 BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** /
  ci.yml；**是否触 §6.1/§6.3 条款语义（预计触 ⇒ 条款先行）**；ADR 是否
  新增由 ① 判定；是否升准则 MINOR（预计 v1.78）由 ① 判定；
- 收益论证闸（逐项：收益 / 代价 / 做或不做 / 依据）：
  a. base 四条孤立计时实测 + index / tree / ci 三作用域重复支付成本
     （基线 = OB-113 的 20.0 s，P88/P92 后需重测）；
  b. 命中收益上限：三作用域 3×20 s → 1×20 s + 键计算代价（键计算
     代价必须入账）；
  c. 键设计草案（复用 C.0aj 句式）：格式版本前缀 + HEAD + porcelain
     快照（**未跟踪文件内容哈希在内**）+ 工具链身份（go version /
     GOROOT / GOOS / GOARCH / CGO_ENABLED）；读前重算键；三反例必
     失效（改已跟踪文件 / 改未跟踪 .go 文件 / 改 go.mod 或切换工具
     链）；任何缓存异常一律 fail-open 回全量计算；
  d. flake 政策两案对照：test 结论「纳入」vs「排除」（DR-2~DR-18 史；
     主控给推荐，用户定夺）；
  e. 条款草案（落点候选 = §6.1 基础门禁段）：只缓存 §6.1 四条 base
     结论（PASS/FAIL + detail 摘要）、绝不缓存任何 G1–G9 判据结论、
     有效性窗口（推荐仅同片 ⑨ 窗口内有效）、flake 政策、升版义务；
- **停点 1 触发即停**：报四项候选取用户裁定，不等额用完；
- 条款先行：② 条款落地 + 条款 CI 绿 ⇒ ③ 起手；条款未落不动代码；
- ③ 一致性：冷热 A/B 逐字一致（含 detail 与输出顺序）+ 比对器翻转 1
  字节可证伪 + 三反例变异测试必失效；
- 若停点 1 = 不做 ⇒ 评估类闭环（同 P95）：零代码改动、准则不升版、
  报告对按 P95 同构收尾。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. docs/DELIVERY_DIRECTIVE{,_EN}.md：若停点 1 = 做 ⇒ 条款落地（预计升
   MINOR v1.78 + §0.3 + 附录 C 下一空位）；若不做 ⇒ 不动；
2. tools/gates/gate.js（base 读路径接缓存）+ tools/gates/selftest.js
   （新断言）——① 方案书点名；
3. docs/validation/gates-base-cache{,_EN}.md（报告对，新增，双语严格
   镜像）；
4. docs/t027/REMAINING_SLICES{,_EN}.md：P96 段翻「✅ 闭环」+ 完成记录；
5. 证据包 docs/validation/evidence/P96-GATES-BASE-CACHE/（00-briefing、
   freedom-list、input-list、收益论证报告、条款草案、A/B 对照、审查原文）；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
7. CHECKLIST_v79{,_EN}（v80 再吸收）；ci.yml；internal/**；CONTRACTS /
   schema / fixtures / workflows；G5-b 闭集与模型白名单；ADR_REGISTER
   （① 判定是否需 ADR）；DEV_PLAN{,_EN}；DOCUMENTATION_PLAN{,_EN}
   （v80 时代换）；s1-exit / p57 门表；find-quotes.js。

八 边界
- 判据判定逻辑本体零改动（只动 base 读路径与缓存）；缓存只存 §6.1
  四条 base 结论、绝不缓存任何 G1–G9 判据结论；不改 ci.yml、
  internal/**、CONTRACTS / schema / fixtures；G5-b 闭集与模型白名单
  不动；不预设临时授权；不做 gitee。

九 回执要求
- 首行 P96: STARTED；四要素（裁决人=用户、日期=2026-11-27、理由=S2 前
  环境整肃 + 层 3 唯一性能候选点名、依据=本起跑消息 + P88 §7 建议 2 +
  v79 快照 §四）；
- 锚点 + 模型/档位（主控 Flash High / ①⑤ V4 Pro / ⑥ 必做 ≤3 / ⑦ Codex
  常规档 M 档 1+2）/ 盲审不抽 / 探针 0 / VM 与代理不开 / GLM-5.3-Flash
  ≥2 逐项复述；
- **停点 1 触发即停**：报四项候选取用户裁定，不得自行决定；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 冷热逐字一致 + 三反例必失效为硬验收，不符即回退；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。

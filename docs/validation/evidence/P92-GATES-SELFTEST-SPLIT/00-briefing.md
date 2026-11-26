# P92 起跑消息（原文）与四要素

- 切片：P92 `GATES-SELFTEST-SPLIT`（工具片 `[SLICE]`，M 档）
- 主控：`DeepSeek V4.1 Flash (deepseek)` High 档
- 落盘时间：2026-11-26（步骤 0）

## 一 四要素

| 要素 | 内容 |
|---|---|
| 裁决人 | 用户 |
| 日期 | 2026-11-26 |
| 理由 | 门禁性能优化续片（性能候选排片裁定） |
| 依据 | 本起跑消息 + v75 快照 §四/§六 + v75 层 3 性能候选行 |

## 二 起跑消息原文（逐字）

```text
切片：P92 GATES-SELFTEST-SPLIT selftest 分区计时与夹具读记忆化（M 档 · 工具片 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P92，M 档工具片
`[SLICE]`，片名 = `GATES-SELFTEST-SPLIT`。目标 = `--selftest` 分区计时 +
夹具读记忆化 + G4b 内环增量评估 + 内建 `--timings` 评估——**判据判定结论
逐字不变**。步骤链首位 = ① 影响判定 + **收益论证闸**。档位启动前设定、
不得中途升档（§8.3）。不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `76d7fc2`（实况）；准则 = v1.74；SELFTEST = 744/744；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v73 / v74 / v75 = 212 / 216 / 215（v75 自身行数由 v76 复核）
- 机械读数 ADR 48 / OB 145 / DR 18 / 试点顺延 25
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P92；标识 GATES-SELFTEST-SPLIT；分级 [SLICE]；类型 = 工具片
  （性能评估 + 轻量实现）；规模 = M（触 tools/gates/selftest.js 与
  lib/criteria/g4b.js 的判据实现）；
- 步骤链：0 → ① 影响判定 + 收益论证闸 → ② 实现（按 ① 裁定范围）→
  ③ 一致性（同输入 A/B sha256 逐字相等）→ ④ 门禁 → ⑤ → ⑥（必做）→
  ⑦ 终审 + 复审 → ⑧a → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash High（默认；**升级闸**：① 论证确定要动 g4b
  内环（sourceDigest / CACHE_VERSION 缓存协议）或 selftest mutation 框架
  深度改造 ⇒ 停下报用户升级 V4 Pro 主控，同轮授权，不得静默切换）；
- 用户裁定（均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 2 复审（§7.5「M 档 1+2」字面）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 分档表判据：触判据
     实现文件（selftest.js / g4b.js）——判据本体与判定逻辑零改动，不触
     OS/安全/并发语义 ⇒ 不触发深度档；Sol 不出场、不双打）；
  c. 盲审不抽（不触身份准入 / 所有权 / 停机语义）；
  d. 探针 0；
  e. 批量辅助 GLM-5.3-Flash 实派 **≥2 次**（§13.1.7 字面；严格串行、
     载体 deep-reasoner、禁 quick-researcher，仅限耗时数据整理 / 双语
     镜像，派单写明字符上限）；
  f. 收益论证闸（用户裁定）：① 必须先论证分区计时 / 记忆化 / G4b 内环 /
     --timings 四项各自的收益与代价；论证收益不足的项**移出本片实现范围**、
     进大改建议清单（P88 报告建议 3 教训：「百毫秒级须先论证收益」）。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P92 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P92 — `GATES-SELFTEST-SPLIT`（selftest 分区计时 + 夹具读记忆化 · 工具片）
| 字段 | 内容 |
|---|---|
| 标识 | `GATES-SELFTEST-SPLIT`（台账编号 = P92；来源 = v72 层 3 性能候选 + v75 §四/§六 下一裁定项 + 用户 2026-11-26 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① `--selftest` 分区计时（按 suite 分区，落盘构成数据）；② 夹具读记忆化（重复磁盘读去重，判据结论逐字不变）；③ G4b 内环增量评估（仅当收益论证通过；动 sourceDigest 则须 CACHE_VERSION/key12 迁移协议先行）；④ 内建 `--timings` 评估（输出格式属协议面，① 判定触契约则移出）；⑤ 输出收益论证报告与对照证据 |
| 依赖 | v75 快照闭环；P88 报告 §7 建议 1/3/5（出处）；无其它前置 |
| 规模 | M |
| 步骤 | 0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | selftest 129 s（P88 实测）无分区构成；夹具重复磁盘读未量化；G4b 内环「百毫秒级」未论证 |
| 交付物 | ① 收益论证报告（分区构成表 + 各候选收益/代价判定）；② 获批项实现（判据结论逐字不变）；③ A/B sha256 对照证据；④ 报告对 `docs/validation/gates-selftest-split{,_EN}.md`；⑤ 台账段与完成记录；⑥ 证据包；⑦ 大改建议清单（未获批项移交） |
| 验收证据 | ① SELFTEST 744/744 不变；② 同输入 A/B 逐字一致（sha256 相同 + 比对器可证伪）；③ 三作用域 TOTAL_FAIL=0；④ 报告对双语严格镜像；⑤ 收益论证结论落盘 |
| 边界 | 判据本体与判定逻辑零改动；不改 ci.yml / internal/** / CONTRACTS / schema / fixtures / workflows；G5-b 闭集与白名单不动；ADR 不新增（不触身份准入）；准则是否升 MINOR 由 ① 判定（触判据实现文件 ⇒ P88 先例升 MINOR v1.75 + 版本包）；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P92-GATES-SELFTEST-SPLIT/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = M 档 ⇒ §7.5「1+2」；档位 = `GPT-5.3-Codex` 常规档；盲审不抽；批量辅助 GLM ≥2（用户裁定 e）；⑥ 必做（触判据实现） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v75 快照闭环 |

一 背景与目标
- 背景：P88 实测 `--selftest` 209992 → 128881 ms 后仍为门禁 CI 步大头；
  报告 §7 建议 1（分区计时/记忆化）、建议 3（G4b 内环，百毫秒级）、建议
  5（内建 --timings，协议面）待评估。
- 目标：见权威输入「目标」字段。
- **收益论证闸为硬前提**：① 方案书必须逐项给「收益 / 代价 / 做或不做」；
  不做的项进大改建议清单，不得默认全做。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro（影响判定 + 收益论证闸 +
  缓存协议判定 + 实现范围裁定）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（A/B 逐字比对 + 744/744 +
  三作用域）；
- ⑤ 预审：DeepSeek V4 Pro，额度 1+2；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，**必做**（触判据
  实现文件）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+2（M 档字面）；
- 盲审：不抽；批量辅助：GLM-5.3-Flash ≥2 次（见裁定 e）；审查输入包
  消毒（§7.3）；计量以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下
  升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域 + go build/vet/
  test ./... 全绿；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v75{,_EN}.md §一）：HEAD = origin =
  `76d7fc2`；准则 v1.74；SELFTEST 744/744；三作用域 0；v73/v74 = 212/216；
  机械读数与 tmp 断言三件。≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P92 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻
  「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P92-GATES-SELFTEST-SPLIT/（⑧c 前只增不
  删；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v75.md 及 _EN（§四 下一项 / 层 3 性能候选行 /
   §六）；
2. docs/DELIVERY_DIRECTIVE.md §7.5（M 档 1+2）/ §5.2（⑦ 分档表）/
   §13.1 / §13.2 / §0.2（版本包三件，若升 MINOR）；
3. tools/gates/selftest.js、tools/gates/lib/criteria/g4b.js、
   tools/gates/lib/ctx.js（改动候选面 + 判据判定边界）；
4. docs/validation/gates-perf-assess{,_EN}.md §7（建议 1/3/5 与耗时基线）；
5. docs/CODING_CONVENTIONS.md（.js = 无 BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** /
  ci.yml；是否升准则 MINOR（触判据实现文件 ⇒ P88 先例升 MINOR）由判定
  定；`--timings` 输出若属协议面 ⇒ 移出本片；
- **收益论证闸（逐项四要素：收益 / 代价 / 做或不做 / 依据）**：
  a. selftest 分区计时——收益 = 构成可见性（后续优化依据）；代价低；
     预计做；
  b. 夹具读记忆化——收益 = 重复磁盘读去重；须证明不改变 mutation 语义；
     ① 定；
  c. G4b 内环增量——P88 建议 3 自评百毫秒级；动 sourceDigest 则须
     CACHE_VERSION 升版 + key12 迁移协议先行；大概率不做，进大改建议；
  d. 内建 --timings——输出协议面，① 判定，预计移出；
- ③ 一致性：获批实现项前后 `SELFTEST` 744/744 逐字不变 + 同输入 A/B
  sha256 相同 + 比对器翻转 1 字节可证伪；
- 升级闸：① 若裁定动 g4b 内环 / mutation 框架 ⇒ 停下报用户升级 V4 Pro
  主控。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. tools/gates/selftest.js / tools/gates/lib/criteria/g4b.js（仅 ① 获批
   实现项；判据判定逻辑零改动；文件由 ① 方案书点名）；
2. docs/DELIVERY_DIRECTIVE{,_EN}.md：若 ① 判定触判据实现文件 ⇒ 升
   MINOR v1.75 + §0.3 行 + 附录 C 下一空位（否则不动）；
3. docs/validation/gates-selftest-split{,_EN}.md（报告对，新增，双语
   严格镜像）；
4. docs/t027/REMAINING_SLICES{,_EN}.md：P92 段翻「✅ 闭环」+ 完成记录；
5. 证据包 docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/（00-briefing、
   freedom-list、input-list、收益论证报告、A/B 对照、审查原文）；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
7. CHECKLIST_v75{,_EN}（v76 再吸收）；ci.yml；internal/**；CONTRACTS /
   schema / fixtures / workflows；G5-b 闭集与白名单；ADR_REGISTER；
   DEV_PLAN{,_EN}；DOCUMENTATION_PLAN{,_EN}（v76 时代换）；s1-exit /
   p57 门表；find-quotes.js。

八 边界
- 只做① 获批实现项；判据本体与判定逻辑零改动；不改 ci.yml、internal/**、
  CONTRACTS / schema / fixtures；G5-b 闭集与白名单不动；不新增 ADR；不做
  gitee；不预设临时授权。

九 回执要求
- 首行 P92: STARTED；四要素（裁决人=用户、日期=2026-11-26、理由=门禁性能
  优化续片（性能候选排片裁定）、依据=本起跑消息 + v75 快照 §四/§六 +
  层 3 性能候选行）；
- 锚点 + 模型/档位（主控 Flash High / ① V4 Pro / ⑥ 必做 / ⑦ Codex 常规
  档 1+2）/ 盲审不抽 / 探针 0 / VM 与代理不开 / GLM ≥2 逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 收益论证闸逐项结论落盘为硬验收；A/B 逐字一致不符即回退；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。
```

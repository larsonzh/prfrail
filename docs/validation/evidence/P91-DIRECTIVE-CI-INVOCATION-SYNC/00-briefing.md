# 00 — 起跑消息原文（P91 `DIRECTIVE-CI-INVOCATION-SYNC`）

本文件按 §13.1.10 逐字落盘用户起跑消息原文。镜像双语纪律不适用于本件（证据树由 G8 承担，见 §6.3 判据作用域排除）。

## 四要素（§1.7.4）

- 裁决人 = 用户
- 日期 = 2026-11-26
- 理由 = 准则 CI 调用句与实况同步（OB-145 处置）
- 依据 = 本起跑消息 + v74 快照 §四/§六 + P90 报告 §8 K4

## 原文

```text
切片：P91 DIRECTIVE-CI-INVOCATION-SYNC 准则 CI 调用句同步（S 档 · 治理片 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P91，S 档小治理
片 [SLICE]，片名 = DIRECTIVE-CI-INVOCATION-SYNC。目标 = 把准则 §6.3
「脚本化状态」句中的 CI 调用描述与 P90 实况同步（--all --skip-base
--scope=ci + base 覆盖口径句），并落版本包。不改任何判据、不改
ci.yml、不改 tools/gates/**。步骤链首位 = ① 影响判定。档位启动前设定、
不得中途升档（§8.3）。不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = 254e98f（实况）；准则 = v1.73；SELFTEST = 744/744；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v72 / v73 / v74 = 214 / 212 / 216（v74 自身行数由 v75 复核）
- 机械读数 ADR 48 / OB 144 / DR 18 / 试点顺延 25（本片把 OB-145 登记并翻
  已处置后，OB max 将变 145，由 v75 换代）
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P91；标识 DIRECTIVE-CI-INVOCATION-SYNC；分级 [SLICE]；类型 =
  治理片（准则描述句同步 + 版本包）；规模 = S（纯 DELIVERY_DIRECTIVE
  {,_EN}.md + 报告对 + 台账 + 证据包）；
- 步骤链：0 → ① 影响判定 → ② 实现 → ③ 一致性（CN/EN 逐字镜像 + 门禁
  绿）→ ④ 门禁 → ⑤ → ⑥（按 §7.2 判定）→ ⑦ 终审 + 复审 → ⑧a → ⑨ →
  ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash High（默认；② 若连续 ≥2 次漂移 ⇒ 停下报用户
  升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 1 复审（§7.5「S 档 1+1」字面优先；本片为描述句
     同步、无规则语义变更，不按「治理片 1+3」取额；① 若判触规则语义或
     身份准入 ⇒ 停下报用户，不得中途升档）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 分档表判据：触准则
     正文描述句（§6.3 脚本化状态）与版本包/台账——不触判据本体与判据
     实现文件（tools/gates/** 零改动），不触身份准入/模型白名单/G5-b
     闭集，不触 OS/安全语义 ⇒ 不触发深度档；Sol 不出场、不双打）；
  c. 盲审不抽（不触身份准入 / 所有权 / 停机语义）；
  d. 探针 0；
  e. 批量辅助 GLM 不派：双语镜像沿用 P89/P90 机械核验口径（行数 / 围栏 /
     G3 字形门禁 + ⑤⑥⑦ 双语对称抽查），用户 2026-11-26 裁定。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P91 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P91 — DIRECTIVE-CI-INVOCATION-SYNC（准则 CI 调用句同步 · 治理片）
| 字段 | 内容 |
|---|---|
| 标识 | DIRECTIVE-CI-INVOCATION-SYNC（台账编号 = P91；来源 = P90 候选 OB-145（准则 :675 v1.40 句滞后）+ v74 §四/§六 下一裁定项 + 用户 2026-11-26 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 把 docs/DELIVERY_DIRECTIVE{,_EN}.md :675 的「CI 在 test 作业末尾追加一个步骤跑本脚本（先 --selftest 再 --all --scope=ci）」同步为「…再 --all --skip-base --scope=ci；base 四件套由同作业 Build / Vet / Test / Format 步覆盖」；② 落版本包（头部 v1.73 → v1.74 MINOR + §0.3 变更行 + 附录 C 下一空位）；③ OB-145 登记入 §12.4 并翻「已处置」；④ 输出双语镜像证据 |
| 依赖 | v74 快照闭环；P90 闭环（拟改文本出处）；无其它前置 |
| 规模 | S |
| 步骤 | 0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | 准则 :675 句与 CI 实况（P90 后 --skip-base + Format 步）不符；OB-145 未登记 §12.4 |
| 交付物 | ① 准则 CN/EN 双侧 :675 句同步（逐字镜像）；② 版本包三件（头 / §0.3 / 附录 C）；③ OB-145 登记 + 翻已处置；④ 报告对 docs/validation/directive-ci-invocation-sync{,_EN}.md；⑤ 台账段与完成记录；⑥ 证据包 |
| 验收证据 | ① CN/EN 双侧 :675 句逐字镜像 + 版本包三件齐全；② 本机三作用域 TOTAL_FAIL=0；③ SELFTEST 744/744 不变；④ 双语 G3 门禁绿；⑤ 报告对严格镜像 |
| 边界 | 不改任何判据（判据本体与判定逻辑零改动）；不改 tools/gates/**；不改 ci.yml；不改 CONTRACTS/schema/fixtures/internal/**；不新增 ADR（① 判定，预计不触身份准入/白名单）；不做 gitee |
| 分级 | [SLICE] |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/freedom-list.md |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | tmp/P91-DIRECTIVE-CI-INVOCATION-SYNC/（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = S 档 1+1（用户裁定，见裁定 a）；档位 = GPT-5.3-Codex 常规档；盲审不抽；批量辅助不派（机械核验，用户裁定） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v74 快照闭环 |

一 背景与目标
- 背景：P90 已把 CI Gates 步改为 --all --skip-base --scope=ci 并新增
  Format 步；准则 :675 的 v1.40 句仍描述旧调用（候选 OB-145）；准则正文
  与 CI 实况脱节，读者按 :675 理解会拿到旧语义。
- 目标：见权威输入「目标」字段。
- 拟改文本唯一出处：P90 报告 docs/validation/gates-ci-perf{,_EN}.md
  §8 的 K4（已含逐字拟改文本与建议片名）——② 按 K4 逐字落，不得偏离；
  若 K4 与实况不符，以实况为准并先经 ① 判定。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro（影响判定：是否触规则语义 /
  身份准入；版本包落法与附录 C 空位核对）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（CN/EN 逐字镜像比对 +
  三作用域 + selftest）；
- ⑤ 预审：DeepSeek V4 Pro，额度 1+1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，① 判定触判据实现
  文件即必做（预计不触，由 ① 判定并留痕）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+1；
- 盲审：不抽；批量辅助：不派（见裁定 e）；审查输入包消毒（§7.3）；计量
  以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域 + go build/vet/
  test ./... 全绿；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v74{,_EN}.md §一）：HEAD = origin =
  254e98f；准则 v1.73；SELFTEST 744/744；三作用域 0；v72/v73 = 214/212；
  机械读数与 tmp 断言三件。≥2 项不符 ⇒ 先报终态读数再继续；
- 台账 P91 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻
  「执行中」并跑门禁（index/tree）；
- R/D = 0 基线；临时件落点 = tmp/P91-DIRECTIVE-CI-INVOCATION-SYNC/
  （⑧c 前只增不删；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v74.md 及 _EN（§四 下一项 / §五 OB-145 候选位 /
   §六）；
2. docs/DELIVERY_DIRECTIVE.md §6.3（:675 句上下文）、§0.2（版本包三件
   义务）、§0.3（变更日志格式）、附录 C（空位顺序）、§12.4（OB 登记
   格式）、§7.5 / §5.2 / §13.1 / §13.2；
3. docs/validation/gates-ci-perf{,_EN}.md §8 的 K4（拟改文本唯一出处）；
4. .github/workflows/ci.yml（Gates 步与 Format 步实况，:95-108）；
5. docs/validation/tmp-archive-p84-p85{,_EN}.md（P89 报告对，机械核验
   口径先例）；
6. docs/CODING_CONVENTIONS.md（.md = BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** /
  tools/gates/** / ci.yml；不触身份准入 / 模型白名单 / G5-b 闭集 ⇒ 不立
  ADR、盲审不抽（与裁定 b/c 一致）；若判定发现触规则语义或身份准入 ⇒
  停下报用户；
- ② 落地：CN/EN 双侧 :675 句按 K4 逐字同步 + 版本包三件（头 v1.74 +
  日期 2026-11-26；§0.3 变更行含「受影响的在飞切片 = 无」；附录 C 下一
  空位 C.0ba 起按实况）+ §12.4 OB-145 登记并翻「已处置」（含四要素）；
- ③ 一致性：CN/EN 双侧改动逐字镜像（行数 / 围栏 / 反引号级）；三作用域
  TOTAL_FAIL=0、SELFTEST 744/744 不变；G3 字形门禁绿；
- 版本包纪律：MINOR 升版属准则变更，改动只落描述句与版本行，不得夹带
  其它准则修改。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. docs/DELIVERY_DIRECTIVE{,_EN}.md：:675 句同步（K4 逐字）+ 头部版本
   v1.73 → v1.74 + §0.3 变更日志行 + 附录 C 台账块 + §12.4 OB-145 登记
   并翻已处置（CN/EN 严格镜像）；
2. docs/validation/directive-ci-invocation-sync{,_EN}.md（报告对，新增，
   双语严格镜像）；
3. docs/t027/REMAINING_SLICES{,_EN}.md：P91 段翻「✅ 闭环」+ 完成记录
   （含提交号与 CI run 观察号）；
4. 证据包 docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/
   （00-briefing、freedom-list、input-list、改动前后逐字对照、审查原文）；
5. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
6. CHECKLIST_v74{,_EN}（v75 再吸收）；ci.yml；tools/gates/**；
   internal/**；CONTRACTS / schema / fixtures；G5-b 闭集与白名单；
   DEV_PLAN{,_EN}；DOCUMENTATION_PLAN{,_EN}（v75 时代换）；s1-exit /
   p57 门表；ADR_REGISTER（不新增）；find-quotes.js。

八 边界
- 只改准则 :675 描述句 + 版本包 + §12.4 OB-145 条目；不改任何判据；
  不改 tools/gates/**、ci.yml、internal/**、CONTRACTS / schema /
  fixtures；不新增 ADR；不做 gitee；不预设临时授权。

九 回执要求
- 首行 P91: STARTED；四要素（裁决人=用户、日期=2026-11-26、理由=准则 CI
  调用句与实况同步（OB-145 处置）、依据=本起跑消息 + v74 快照 §四/§六 +
  P90 报告 §8 K4）；
- 锚点 + 模型/档位（主控 Flash High / ① V4 Pro / ⑦ Codex 常规档 1+1）/
  盲审不抽 / 探针 0 / VM 与代理不开 / GLM 不派逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- CN/EN 逐字镜像 + 版本包三件齐全为硬验收，缺一即回退；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。
```

切片：P89 TMP-ARCHIVE-P84-P85 tmp 两目录归档与引用改写（S 档 · 工具片 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P89，S 档工具片
`[SLICE]`，片名 = `TMP-ARCHIVE-P84-P85`。目标 = 归档 tmp 两目录进证据包 +
引用改写 + 清除；不新增任何判据/生产逻辑。步骤链首位 = ① 影响判定。档位
启动前设定、不得中途升档（§8.3）。不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `c57eccd`（实况）；准则 = v1.73；SELFTEST = 744/744；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v70 / v71 / v72 = 203 / 201 / 214（v72 自身行数由 v73 复核）
- 机械读数 ADR 48 / OB 144 / DR 18 / 试点顺延 25
- tmp 断言顶层五件：.gitkeep / gates-cache(g4b) / P84-CI-FIX /
  find-quotes.js / P85-SECURITY-PROOF-SUBSYSTEM

〇 切片标识
- 编号 P89；标识 TMP-ARCHIVE-P84-P85；分级 [SLICE]；类型 = 工具/证据片；
  规模 = S（纯 docs/validation/** + 台账 + 一处测试文件注释；不触
  internal/** 生产逻辑）；
- 步骤链：0 → ① 影响判定 → ② 归档与引用改写 → ③ 一致性（归档包自洽 +
  引用可解析）→ ④ 门禁 → ⑤ → ⑥（按 §7.2 判定）→ ⑦ 终审 + 复审 → ⑧a →
  ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash High（默认；② 若连续 ≥2 次漂移 ⇒ 停下报用户
  升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 1 复审（§7.5「S 档 1+1」字面）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 分档表判据：触证据包
     历史文件 / 测试文件注释 / 台账字段——不触判据本体与判据实现文件 ⇒
     Codex；不触 OS/安全语义，不触发深度档；轻量档 Hydrogen 有 P87 前置
     缺口未实测，本片不走）；
  c. 盲审不抽（不触身份准入 / 所有权 / 停机语义）；
  d. 探针 0。
- 编号说明：P89 系主控在 v72 按顺序取号（草案未给）；点名即默认接受。

〇.5 权威输入（逐字执行；§1.2 全字段为本片唯一落地依据）
台账 P89 段（docs/t027/REMAINING_SLICES{,_EN}.md，:2389 起，⬜ 未立项）：
目标 / 依赖 / 规模 / 步骤 / 交付物 / 验收证据 / 边界 / 额度口径 / 引用点
实测清单（9 文件 15 处）等全字段以台账为准，此处不复制。注意：台账目标
字段把「改写 or 更正注记」口径写作 §11.1 系登记引用不精确——执行以准则
原文为准（§1.5.1 tmp 生命周期四类口径 + §11.2 回写口径；§11.1 仅约束
历史行不追溯与双语），闭环时顺手校正台账该处引用。

一 背景与目标
- 背景：tmp/P84-CI-FIX 与 tmp/P85-SECURITY-PROOF-SUBSYSTEM 为 P84/P85 证据
  与设计件，仍被 9 个文件 15 处引用支撑；v72 §四 已列为本片承接。
- 目标：见台账「目标」字段（归档 / 引用处置 / 清除 / 自洽证据）。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro（影响判定 + 证据包内历史文件
  改写 or 更正注记口径 + 归档方案）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（归档包三方自洽 + 引用
  grep 零命中 + 门禁回归）；
- ⑤ 预审：DeepSeek V4 Pro，额度 1+1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，① 判定触判据文件即
  必做（预计不触，由 ① 判定）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+1（S 档字面）；
- 盲审：不抽；批量辅助：GLM-5.3-Flash ≥1、严格串行、载体 deep-reasoner
  （禁 quick-researcher），仅限引用清单核对 / 双语镜像，派单写明字符上限；
- 审查输入包消毒（§7.3）；计量以调用回执为单一真值源；探针 0；
- 风险提醒：R 达 2 ⇒ §9.2 高风险标记并停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机门禁三作用域 + go test
  ./... 全绿 + CI 双腿（Gates + selftest）观察。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v72{,_EN}.md §一）：HEAD = origin =
  `c57eccd`；准则 v1.73；SELFTEST 744/744；三作用域 0；v70/v71 = 203/201；
  机械读数与 tmp 断言。≥2 项不符 ⇒ 先报终态读数再继续；
- 台账 P89 段已登记 ⇒ 步骤 0 翻「执行中」并跑门禁（index/tree）；
- R/D = 0 基线；临时件落点 = tmp/P89-TMP-ARCHIVE-P84-P85/（⑧c 前只增不删；
  归档动作本身不得在 tmp 根新增文件；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v72.md 及 _EN（§四 下一项 + tmp 断言）；
2. docs/t027/REMAINING_SLICES{,_EN}.md 的 P89 段（§1.2 全字段 + 引用清单）；
3. docs/DELIVERY_DIRECTIVE.md §7.5 / §5.2 / §1.5.1（tmp 生命周期四类口径，
   含「有引用者归档」）/ §11.1–§11.2（历史不追溯 + 回写口径）/ §13.1 /
   §13.2；
4. docs/validation/bounddir-residual{,_EN}.md（P84 报告对，§12 注记点
   :156/:160）与 docs/validation/security-proof-subsystem{,_EN}.md（P85
   报告对）；
5. P85 证据包内 4 个含引用文件（00-briefing.md / 07-carbon-r7-targeted-
   read.md / P87-four-read-data/07-codex-baseline-read.md /
   P87-four-read-data/diff.patch）；
6. docs/CODING_CONVENTIONS.md（.md = BOM + LF；.js/.go = 无 BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / 契约夹具 / workflows / CI
  workflow；是否升准则 MINOR 由判定定（大概率不触）；
- 引用处置口径（① 定，依据 §1.5.1 + §11.2）：活跃引用（代码注释 / 台账
  字段 / 报告注记 / 证据包内相对路径）改写；历史原文（P85 起跑消息
  00-briefing.md:109 与 diff.patch:237,11818）倾向**加更正注记**（不改历史
  原文），最终口径 ① 判、⑦ 背书；
- 引用点实测清单（9 文件 15 处，以 git grep -F 为准，台账已列全）：
  · internal/console/p85_t2_frozen_policy_tester_test.go:6（注释）
  · REMAINING_SLICES{,_EN}.md 的 P85 段落点声明（:2314/:2301）
  · bounddir-residual{,_EN}.md §12（:156/:160）
  · P85 证据包 07-carbon-r7-targeted-read.md:10
  · P85 证据包 07-codex-baseline-read.md:11/14/21/26
  · P85 证据包 00-briefing.md:109
  · P85 证据包 diff.patch:237/11818
- 归档包自洽：MANIFEST ↔ SHA256SUMS ↔ 磁盘条目集三方相等；
- 验收硬性：git grep -F 两条 tmp 路径零命中（除 ① 判定的更正注记处，注记
  措辞不得含可解析完整路径或按 G9 白名单处置）；tmp 顶层剩三件（.gitkeep /
  gates-cache / find-quotes.js）；SELFTEST 744/744 不变；三作用域
  TOTAL_FAIL=0；go test ./... 全绿（测试文件注释改动）。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. P84 / P85 证据包扩容（两目录归档并入）+ MANIFEST / SHA256SUMS 重建；
2. docs/validation/bounddir-residual{,_EN}.md §12 字段注记（:156/:160）；
3. 台账 P85 段落点声明字段（:2314/:2301）+ 台账 P89 段翻「✅ 闭环」+ 完成
   记录（含改写前后逐处清单与三方自洽证据指针）+ 顺手校正 P89 段目标字段
   的 §11.1 引用（改指 §1.5.1 + §11.2）；
4. 报告对 docs/validation/tmp-archive-p84-p85{,_EN}.md（新增，双语严格
   镜像）；
5. 证据包 docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/（00-briefing、
   freedom-list、input-list、改写前后清单、归档自洽证据）；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
7. CHECKLIST_v72{,_EN}（v73 再吸收）；DELIVERY_DIRECTIVE（除非 ① 判定）；
   CONTRACTS / schema / fixtures / workflows；CI workflow；internal/**
   生产逻辑；G5-b 闭集与白名单；ADR 不新增；DEV_PLAN{,_EN}；
   DOCUMENTATION_PLAN{,_EN}（v73 时代换）；find-quotes.js 不处置。

八 边界
- 不动 find-quotes.js；不做 tmp 全清；不碰准则正文（除非 ① 判定）；
  不改 internal/** 生产逻辑（仅一处测试文件注释行）；不改 CONTRACTS /
  schema / fixtures / workflows / CI workflow；不做 gitee；不预设临时授权。

九 回执要求
- 首行 P89: STARTED；四要素（裁决人=用户、日期=2026-11-25、理由=tmp 归档
  清欠（用户 2026-11-24 排期 + 「归档片先」裁定）、依据=本起跑消息 + v72
  快照 §四 + 台账 P89 段）；
- 锚点 + 模型/档位（主控 Flash High / ① V4 Pro / ⑦ Codex 常规档 1+1）/
  盲审不抽 / 探针 0 / VM 与代理不开逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- git grep 零命中为硬验收，不符即回退；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。

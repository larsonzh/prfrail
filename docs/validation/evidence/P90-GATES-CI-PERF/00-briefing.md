# P90 `GATES-CI-PERF` — 起跑消息原文

- 切片：P90 / `GATES-CI-PERF` / `[SLICE]` / 工具片 / S 档
- 落盘时点：2026-11-25（步骤 0）
- 用途：§13.1.10 起跑消息原文留痕；四要素见文末。
- 说明：以下为起跑消息逐字正文，未经改写。

---

切片：P90 GATES-CI-PERF CI 门禁步去重与 Go 缓存评估（S 档 · 工具片 [SLICE]）
=====================================================================

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4.1 Flash，High 档）。本片 = P90，S 档工具片
`[SLICE]`，片名 = `GATES-CI-PERF`。目标 = 评估并把 CI Gates 步 base 四件套
去重落地（`--skip-base`）+ Go 缓存持久化评估；**判据判定逻辑零改动**。
步骤链首位 = ① 影响判定。档位启动前设定、不得中途升档（§8.3）。不预设
临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `6e097fd`（实况）；准则 = v1.73；SELFTEST = 744/744；
  三作用域 = TOTAL_FAIL=0
- 滚动窗口 v71 / v72 / v73 = 201 / 214 / 212（v73 自身行数由 v74 复核）
- 机械读数 ADR 48 / OB 144 / DR 18 / 试点顺延 25
- tmp 断言顶层三件：.gitkeep / gates-cache(g4b) / find-quotes.js

〇 切片标识
- 编号 P90；标识 GATES-CI-PERF；分级 [SLICE]；类型 = 工具片（CI 配置 +
  评估）；规模 = S（纯 .github/workflows/ci.yml + 报告对 + 台账 + 证据包；
  不触 tools/gates/**、不触 internal/**）；
- 步骤链：0 → ① 影响判定 → ② 实现 → ③ 一致性（判据结论面逐字一致 +
  base 覆盖等价论证）→ ④ 门禁 → ⑤ → ⑥（按 §7.2 判定）→ ⑦ 终审 + 复审 →
  ⑧a → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4.1 Flash High（默认；② 若连续 ≥2 次漂移 ⇒ 停下报用户
  升级 V4 Pro，同轮授权，不得静默切换）；
- 用户裁定（均不得作先例）：
  a. ⑦ 额度 = 1 终审 + 1 复审（§7.5「S 档 1+1」字面）；
  b. ⑦ 档位 = GPT-5.3-Codex (copilot) 常规档（§5.2 分档表判据：触 CI
     workflow 配置（执行纪律面）与报告/台账——不触判据本体与判据实现文件
     （gate.js / selftest.js / ctx.js 零改动），不触 OS/安全/并发语义 ⇒
     不触发深度档；Sol 不出场、不双打）；
  c. 盲审不抽（不触身份准入 / 所有权 / 停机语义）；
  d. 探针 0；
  e. 批量辅助 GLM 不派：双语镜像沿用 P89 机械核验口径（行数 / 围栏 /
     grep），用户 2026-11-25 裁定；
  f. CI workflow 改动授权：本片允许修改 .github/workflows/ci.yml 的 Gates
     步参数与 Go 缓存相关配置（用户点名即授权；属「frozen CI contract」的
     用户授权扩展，先例 = fetch-depth: 0，见 ci.yml:63 注释）。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P90 段尚未登记 ⇒
步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P90 — `GATES-CI-PERF`（CI 门禁步 base 去重 + Go 缓存评估 · 工具片）
| 字段 | 内容 |
|---|---|
| 标识 | `GATES-CI-PERF`（台账编号 = P90；来源 = v72 层 3 性能候选 + v73 §四/§六 下一裁定项 + 用户 2026-11-25 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 评估并落地 CI Gates 步 base 四件套去重：`--all --scope=ci` → `--all --skip-base --scope=ci`（selftest 保留；base 在 CI 上由 Build/Vet/Test 步覆盖，覆盖等价性由 ① 书面论证，含 gofmt 缺口处置）；② Go build/test 缓存持久化评估（现状 = Setup Go 已 `cache: true`，评估结论可能为「已覆盖，仅注记/微调」）；③ 输出 CI 前后对照证据（Gates 步耗时 + 判据结论面逐字一致）与后续大改建议 |
| 依赖 | v73 快照闭环；无其它前置 |
| 规模 | S |
| 步骤 | 0 → ① → ② → ③ → ④ → ⑤ → ⑥ → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | CI Gates 步重复执行 base 四件套（Build/Vet/Test 步已覆盖其后三件）；「ci 冻结成本」无 §12.4 条目 |
| 交付物 | ① workflow 改动（Gates 步 --skip-base + gofmt 缺口处置 + Go 缓存评估落点）；② 覆盖等价论证（含 gofmt 三选一处置结论）；③ 报告对 `docs/validation/gates-ci-perf{,_EN}.md`；④ 台账段与完成记录；⑤ 证据包；⑥ 大改建议清单 |
| 验收证据 | ① CI 双腿 Gates 步：`SELFTEST: PASS 744/744` 不变 + `GATE REPORT scope=ci` 判据结论与改动前逐字一致（TOTAL_FAIL=0）；② 本机三作用域 TOTAL_FAIL=0；③ 改动前后 Gates 步耗时对照（改动前 run vs 改动后 run 同 runner 读数 + 命令）；④ 报告对双语严格镜像；⑤ go 缓存评估结论落盘 |
| 边界 | 不触 tools/gates/**（gate.js / selftest.js / ctx.js / lib/criteria/** 零改动）；不触判据本体与判定逻辑；不改 candidate-build / release 等 job；不改 internal/**、CONTRACTS/schema/fixtures；不升准则（ADR 由 ① 判定：若判定触冻结 CI 契约需留痕则立 ADR-049，否则仅 workflow 注释留痕，先例 = ci.yml:63）；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P90-GATES-CI-PERF/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = S 档 ⇒ §7.5「S 档 1+1」；档位 = `GPT-5.3-Codex` 常规档；盲审不抽；批量辅助不派（机械核验，用户裁定） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v73 快照闭环 |

一 背景与目标
- 背景：CI test job 内 Build/Vet/Test 步已各跑一遍 base 后三件，Gates 步
  `--all` 再串行重复一遍（P88 实测 base 热 93.8 s / 冷 178.1 s）；P88 报告
  §7 建议 2 的短期替代 = CI Gates 步改 `--all --skip-base --scope=ci`；
  Setup Go 已 `cache: true`。
- 目标：见权威输入「目标」字段。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro（影响判定 + base 覆盖等价论证
  （含 gofmt 缺口处置）+ go 缓存评估口径 + workflow 改法）；
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash；
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash（本机三作用域 + 判据结论
  逐字比对 + yaml 语法校验）；
- ⑤ 预审：DeepSeek V4 Pro，额度 1+1；
- ⑥ 独立扫描：MAI Code 2 Flash (copilot)，额度 ≤3，① 判定触判据文件即
  必做（本片不触判据文件，预计跳过，由 ① 判定并留痕）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档，额度 1+1（S 档字面）；
- 盲审：不抽；批量辅助：不派（见裁定 e）；审查输入包消毒（§7.3）；计量
  以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 本机三作用域 + go build/vet/
  test ./... 全绿 + yaml 校验；**CI 双腿为本片主验证场**（workflow 改动
  本机不可模拟，push 后按 §11.2 取双腿 Gates 步读数回填报告）。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v73{,_EN}.md §一）：HEAD = origin =
  `6e097fd`；准则 v1.73；SELFTEST 744/744；三作用域 0；v71/v72 = 201/214；
  机械读数与 tmp 断言三件。≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P90 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻
  「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P90-GATES-CI-PERF/（⑧c 前只增不删；D
  累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md
  （§13.1.10，落 docs/validation/evidence/P90-GATES-CI-PERF/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字
  机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v73.md 及 _EN（§四 下一项 / 层 3 性能候选行 / §六）；
2. docs/t027/REMAINING_SLICES{,_EN}.md 的 P89 完成记录（D=1 口径参考）；
3. docs/DELIVERY_DIRECTIVE.md §7.5 / §5.2 / §11.2 / §13.1 / §13.2；
4. .github/workflows/ci.yml（本片唯一改动面；Gates 步 :93-94、Build/Vet/
   Test 步 :71-84、Setup Go cache :71、冻结契约注释 :63）；
5. docs/validation/gates-perf-assess{,_EN}.md（P88 报告对 §7 建议 2 + 耗时
   基线）；docs/validation/tmp-archive-p84-p85{,_EN}.md（P89 报告对，机械
   核验口径先例）；
6. docs/CODING_CONVENTIONS.md（.yml/.md 编码与行尾约定）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 影响判定：预计不触 CONTRACTS / schema / fixtures / internal/** /
  tools/gates/**；是否立 ADR 由判定定（「frozen CI contract」扩展留痕
  方式：ADR-049 或 workflow 注释，先例 = ci.yml:63 注释未立 ADR）；
- 覆盖等价论证（硬前提）：base 四件套 = gofmt / go build / go vet /
  go test；其中 build / vet / test 由 CI Build / Vet / Test 步覆盖（三步均
  无 continue-on-error 豁免 ⇒ 任一步失败 job 即失败）；**gofmt 在 CI 无
  独立步**——① 须论证并选定处置（三选一）：a) ci.yml 新增 gofmt 步
  （`gofmt -l .` 非空输出即失败；防线不变，推荐）；b) 论证 gofmt 由「提交
  前本机门禁绿」纪律覆盖并立 ADR 留痕（CI 验收面微调）；c) 其它等价方案。
  ① 定、⑦ 背书；Gates 步 --skip-base 后判据面（G1–G9 + G4b）不变，逐条
  书面论证；
- ③ 一致性：改动前 Gates 步输出（现有 CI run 的 `GATE REPORT scope=ci`）
  vs 改动后同变更集 `--skip-base` 输出——**判据结论逐字一致**（INFO/
  PASS/FAIL 行除 base 段外相同）；本机用同命令复现比对；
- go 缓存评估：现状 = Setup Go `cache: true`（v7 自带 build/module cache
  持久化）；评估结论三选一并落盘：已覆盖不动 / 微调（如 windows 腿 cache
  restore 容错）/ 另有缺口需后续片；
- 大改建议清单单独成节（GATES-SELFTEST-SPLIT 等三候选不动，本片不实施）。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. .github/workflows/ci.yml：Gates 步 `--all --scope=ci` →
   `--all --skip-base --scope=ci`（selftest 保留）；go 缓存按评估结论落点
   （不动则加注记说明现状）；若 ① 选 gofmt 步方案则同步新增该步；
2. docs/validation/gates-ci-perf{,_EN}.md（报告对，新增，双语严格镜像）；
3. docs/t027/REMAINING_SLICES{,_EN}.md：P90 段翻「✅ 闭环」+ 完成记录
   （含改动前后耗时与 CI run 号回填）；
4. 证据包 docs/validation/evidence/P90-GATES-CI-PERF/（00-briefing、
   freedom-list、input-list、覆盖等价论证、前后输出比对、CI 读数）；
5. docs/ADR_REGISTER{,_EN}.md：仅当 ① 判定需立 ADR-049；
6. tools/gates/cjk-newwords.txt：如新增 CJK 字即登记。
不动：
7. CHECKLIST_v73{,_EN}（v74 再吸收）；DELIVERY_DIRECTIVE{,_EN}；
   CONTRACTS / schema / fixtures / workflows（除本片 ci.yml 改动外）；
   internal/**；tools/gates/**（零改动）；G5-b 闭集与白名单；DEV_PLAN{,_EN}；
   DOCUMENTATION_PLAN{,_EN}（v74 时代换）；s1-exit / p57 门表；
   candidate-build / release 等其它 job；find-quotes.js。

八 边界
- 只改 ci.yml 的 Gates 步参数与 Go 缓存相关配置；不新增/不改任何判据的
  判定逻辑；不触 gate.js / selftest.js / ctx.js / lib/criteria/**；不改
  internal/**、CONTRACTS / schema / fixtures；不升准则；不做 gitee；不预设
  临时授权。

九 回执要求
- 首行 P90: STARTED；四要素（裁决人=用户、日期=2026-11-25、理由=门禁性能
  优化续片（用户 2026-11-25「归档片先、性能片后」排片裁定 + v73 §六下一
  裁定项）、依据=本起跑消息 + v73 快照 §四/§六 + 层 3 性能候选行）；
- 锚点 + 模型/档位（主控 Flash High / ① V4 Pro / ⑦ Codex 常规档 1+1）/
  盲审不抽 / 探针 0 / VM 与代理不开 / GLM 不派逐项复述；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- CI 双腿 Gates 步判据结论逐字一致为硬验收，不符即回退；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，
  等同轮授权；仅 origin。

---

## 四要素（留痕）

- 裁决人 = 用户
- 日期 = 2026-11-25
- 理由 = 门禁性能优化续片（用户 2026-11-25「归档片先、性能片后」排片裁定 + v73 快照 §六下一裁定项）
- 依据 = 本起跑消息 + v73 快照 §四/§六 + 层 3 性能候选行
- 不得作先例

## 步骤 0 附加留痕：边界扩展（主控代裁）

- 事实：起跑消息「边界」与「六 实现要点」均把 `internal/**` 列为不动面，并将「是否触 internal/**」列为 ① 的预判项。实测 `internal/release/workflow_test.go:249` 对 CI `Gates` 步 run 串做 sha256 断言（现值 `ca20623b…c8d0`），`:187` 固化 `test` 作业步骤名表，`:322-323` 固化 Action 参数；仅改 `ci.yml` 不改该测试，CI 两腿必红。
- 先例核验：`fetch-depth: 0` 那次用户授权扩展同样改了该测试（`:322` 的 `test/Checkout source` 参数含 `fetch-depth: 0`）⇒ 本仓库「frozen CI contract 用户授权扩展」的实际形态包含契约测试同步。
- 处置：用户当轮不可达（`vscode_askQuestions` 返回「not available」），主控按「目标优先 + 仓库级硬规则（改 CI 必须同步该测试）」代裁，将 `internal/release/workflow_test.go` 的**冻结契约同步**（摘要值、步骤名表、必要的脚本契约条目）纳入本片；`internal/**` 其余文件零改动。
- 标注：**含规则面扩展、超出起跑消息字面边界、系主控代裁**，须用户事后追认；**不得作先例**。

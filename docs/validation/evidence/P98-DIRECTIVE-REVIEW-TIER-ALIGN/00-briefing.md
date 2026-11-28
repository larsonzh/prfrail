# P98 起跑消息原文（00-briefing）

> 以下为用户 2026-11-28 起跑消息逐字原文（仅移除最外层会话包装；正文一字未改）。

切片：P98 DIRECTIVE-REVIEW-TIER-ALIGN ⑦ 档位对齐修订（M 档 · 治理片 · O-2..O-4 处置）

【〇·零】记忆唤醒（开场必读；起草惯例红线，未登入准则/OB/ADR）
你是 VS Code 主控（DeepSeek V4 Pro (deepseek)，Max 档——用户 2026-11-28 裁定：主控升级 V4 Pro（连续 2 次行为漂移触发选型策略场景④），且与 ①⑤ 的 V4 Pro 保持同一思考模式档位 Max）。本片 = P98，M 档治理片，片名 = `DIRECTIVE-REVIEW-TIER-ALIGN`。目标 = 修订 §5.2 / §7.5 / §8.2 三处档位与措辞（O-2..O-4 处置方向 = 尽快处理），关闭 v81 快照 §五 的 O-2..O-4 与 O-5。**只改这三节与版本包，不改任何判据实现与闭集**。档位启动前设定、不得中途升档（§8.3）。不预设临时授权——额度用尽按 §7.8.3 停。

关键锚点（步骤 0 须实测核对，以仓库实况为准）：
- HEAD / origin = `ada55b2`（实况）；准则 = v1.82；SELFTEST = 806/806；三作用域 = TOTAL_FAIL=0
- 滚动窗口 v79 / v80 / v81 = 209 / 210 / 213（按 LF 计数；v81 自身行数由 v82 复核）
- 机械读数 ADR 48 / OB 146 / DR 18 / 试点顺延 25
- tmp 断言顶层两件：.gitkeep / gates-cache（其下 base 与 g4b）
- **HEAD 特别说明**：v81 快照 §一表的 `c30b2de` 是 v81 生成时点锚点；其后 v81 快照提交 `ada55b2` 已落盘，HEAD 前移至 `ada55b2`——以 `git rev-parse --short HEAD` 实况为准，不同属预期，不判锚点不符。

〇 切片标识
- 编号 P98；标识 DIRECTIVE-REVIEW-TIER-ALIGN；分级 [SLICE]；类型 = 治理片（准则修订 §5.2 / §7.5 / §8.2 三节 + 版本包 + 台账处置）；规模 = M；
- 步骤链：0 → ① 方案（修订文本三件套 + 升版判定 + OB 处置口径 + ADR 判定；**停点 1**）→ ② 条款落地（准则升 MINOR v1.83）→ 条款 CI 绿 → ③ 台账与观察项处置（O-2..O-4 / O-5 翻「已处置」）→ ④ 门禁 → ⑤ 预审 → ⑥（跳过，经 ⑦ 确认）→ ⑦ 终审 + 复审 → ⑧a 报告对 → ⑨ → ⑧b → ⑧c → ⑩；
- 主控 = DeepSeek V4 Pro (deepseek) **Max**（用户 2026-11-28 裁定，与 ①⑤ 的 V4 Pro 同档，不得作先例）；
- **停点 1 闸（硬）**：① 产出四项——修订文本三件套（§5.2 对齐句 + 可判性措辞 / §8.2 取档授权句 / §7.5 Neon 单价行 + 备选额度措辞句）、升版判定（预计 v1.83 MINOR）、OB 处置口径（O-2..O-4 / O-5 翻「已处置（P98）」）、ADR 判定（预计不新增；ADR-048 是否加修订注记）——主控在 ① 停下时报用户裁定，**不得自行决定**；
- ⑦ 档位 = **已定**：常规档 `GPT-5.3-Codex (copilot)`，**治理片行 1 终审 + 3 复审（1+3，§7.5）**（依据 = 预算策略「触 §7.5/§5.2 规则语义 ⇒ Codex」+ §7.5 治理片行）；
- 用户裁定（均不得作先例；停点 1 四项待你回场裁定）：
  a. 停点 1 四项 = 待裁定（修订文本、升版、OB 处置、ADR）；
  b. 盲审不抽（不触身份准入 / 所有权 / 停机）；
  c. 探针 0；
  d. 批量辅助 GLM-5.3-Flash 实派 ≥2 次（modelId `glm-5.3-flash`，与 §2.2 批量辅助首选一致；**禁用 `glm-5.3`**——那是 §2.2 的 ⑤ 预审备用档，两型号不得混用，先例 = OB-83）；严格串行、载体 deep-reasoner、禁 quick-researcher，仅限耗时数据整理 / 双语镜像，派单写明字符上限。

〇.5 权威输入（§1.2 全字段为本片唯一落地依据；台账 P98 段尚未登记 ⇒ 步骤 0 按本节先落台账「⬜ 未立项」再翻「执行中」）

### P98 — `DIRECTIVE-REVIEW-TIER-ALIGN`（⑦ 档位对齐修订 · 治理片）
| 字段 | 内容 |
|---|---|
| 标识 | `DIRECTIVE-REVIEW-TIER-ALIGN`（台账编号 = P98；来源 = v81 快照 §五 O-2..O-4 处置方向（尽快处理）+ P97 报告观察项 + 用户 2026-11-28 点名） |
| 状态 | ⬜ 未立项（本消息即点名；四要素见 §九） |
| 目标 | ① 修订 §5.2 深度档档位列与取档授权句对齐 §8.2（O-2 High：Neon `Max` vs `Extra High` 冲突）+ 矩阵可判性措辞（O-3 Medium）；② §7.5 Neon 单价档位归属补全（O-4 Medium）+ 备选额度口径措辞（O-5 #4 Low）；③ O-5 #5（numstat 陈旧）只登记不修（冻结证据包件）；④ 观察项 O-2..O-4 / O-5 处置翻「已处置（P98）」 |
| 依赖 | v81 快照闭环；P97 报告对（O-2..O-4 原文）；无其它前置 |
| 规模 | M |
| 步骤 | 0 → ①（停点 1）→ ② 条款落地（v1.83）→ 条款 CI 绿 → ③ 台账处置 → ④ → ⑤ → ⑥ 跳过 → ⑦ → ⑧a → ⑨ → ⑧b → ⑧c → ⑩ |
| 现状缺口 | §5.2 深度档 Neon 档位列与 §8.2 冲突（O-2）；矩阵字面可判性不足（O-3）；§7.5 Neon 单价档位归属缺（O-4）；备选额度措辞易误读（O-5 #4） |
| 交付物 | ① 修订文本三件套；② 准则 v1.83（头部 / §0.3 新行 / §5.2 / §7.5 / §8.2 / 附录 C 新块 C.0bf）；③ 台账 P98 段与完成记录 + P97 段观察项处置翻转；④ 报告对 `docs/validation/directive-review-tier-align{,_EN}.md`；⑤ 证据包 |
| 验收证据 | ① 修订后 §5.2 与 §8.2 档位逐字一致（Neon = Extra High）、取档授权句对齐；② §7.5 Neon 两档单价齐备且档位归属明确；③ 六形状双语镜像 HEAD_MISMATCH=0；④ SELFTEST 806/806 不变（判据实现零改动）+ 三作用域 TOTAL_FAIL=0 + CI 双腿 success；⑤ O-2..O-4 / O-5 台账处置翻转双语同步 |
| 边界 | 只改 §5.2 / §7.5 / §8.2 相关行 + 版本包 + 台账与报告；不改 tools/gates/** / internal/** / ci.yml / CONTRACTS / schema / fixtures / workflows；G5-b 闭集与白名单不动；不新增 ADR（① 判定）；P87 报告对与证据包冻结不追改（numstat 只登记不修）；不做 gitee |
| 分级 | `[SLICE]` |
| 探针额度 | 0 |
| 执行自由度清单落盘路径 | `docs/validation/evidence/P98-DIRECTIVE-REVIEW-TIER-ALIGN/freedom-list.md` |
| R/D 基线 | R = 0、D = 0（起跑时登记） |
| 临时件落点声明 | `tmp/P98-DIRECTIVE-REVIEW-TIER-ALIGN/`（⑧c 前只增不删） |
| 额度口径披露 | ⑦ = 常规档 `GPT-5.3-Codex (copilot)` **1 终审 + 3 复审（治理片行，§7.5）**；⑥ 跳过（不触判据实现文件，经 ⑦ 确认，P72 先例）；盲审不抽；批量辅助 GLM-5.3-Flash ≥2（modelId `glm-5.3-flash`，禁 `glm-5.3`） |
| 起跑条件 | 用户点名（本消息）；四要素留痕；v81 快照闭环 |

一 背景与目标
- 背景：P97（HYDROGEN-CONTROL-PILOT）用轻量档 `Gemini 4 Hydrogen` 对 P87 交付物做对照实测，检出 5 条发现——O-2（High）：§5.2 深度档档位列写 Neon = `Max`，与 §8.2「Neon = `Extra High`」及取档授权句直接冲突（P87 原 ⑦ PASS 未检出）；O-3（Medium）：矩阵字面可判性不足（「定稿复核 / 定向小读」无定义、逻辑连接词与优先级缺失、「其余三档」漏计备选档、备选门槛与「最高风险轮」无界线）；O-4（Medium）：§7.5 Neon 分 200K / 1M 两档但主表 300 / 1500 未标注归属、另一档单价缺失；O-5（Low ×2）：备选额度口径首句误读张力、numstat 证据件自记行数陈旧。处置方向（用户 2026-11-28 裁定）= 尽快处理，本片即该修订片。
- 目标：见权威输入「目标」字段。
- **停点 1 闸为本片安全前提**：① 判定后停下报用户，不得自行决定修订文本、升版、OB 处置与 ADR。

二 模型与审查链（启动前一次性设定，不得中途升档）
- ① 架构师：deep-reasoner + DeepSeek V4 Pro (deepseek) **Max**（修订文本三件套 + 升版 / OB 处置 / ADR 判定；本片修订面明确，① 单打即可，不需要 Sol 双打）；
- ②③ 文档：prfrail-documenter + DeepSeek V4.1 Flash（文档型退化，§13.1.9）；
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek) **Max**；
- ⑥ 独立扫描：**跳过**（不触 gate.js / selftest.js 等判据实现文件；经 ⑦ 确认，P72 先例）；
- ⑦ 独立终审：GPT-5.3-Codex (copilot) 常规档 **1 终审 + 3 复审（治理片行，§7.5）**；盲审不抽；批量辅助 GLM-5.3-Flash ≥2（modelId `glm-5.3-flash`，见裁定 d）；审查输入包消毒（§7.3）；计量以调用回执为单一真值源；探针 0；R 达 2 ⇒ §9.2 停下升级。

三 环境声明
- Ubuntu VM 不开、代理不开、不做 gitee；⑨ = 三作用域 `--all --skip-base` 全绿 + SELFTEST 806/806 不变；CI 双腿按 §11.2 观察记录。

四 步骤 0 必做（§13.1）
- 锚点核对（出处 = docs/t027/CHECKLIST_v81{,_EN}.md §一）：HEAD 实况 `ada55b2`（表值 `c30b2de` 为生成时点锚点，前移属预期，见 〇·零）；准则 v1.82；SELFTEST 806/806；三作用域 0；v79/v80 = 209/210（v81 = 213 按 LF 计数，自身留 v82）；机械读数 ADR 48 / OB 146 / DR 18 / 顺延 25；tmp 断言两件。≥2 项不符 ⇒ 先报终态读数再继续；
- **台账 P98 段未登记 ⇒ 步骤 0 先按 〇.5 落台账「⬜ 未立项」（双语），再翻「执行中」并跑门禁（index/tree）**；
- R/D = 0 基线；临时件落点 = tmp/P98-DIRECTIVE-REVIEW-TIER-ALIGN/（⑧c 前只增不删；D 累计达 3 ⇒ 强制停片）；
- input-list.md + freedom-list.md；起跑消息原文落盘 00-briefing.md（§13.1.10，落 docs/validation/evidence/P98-DIRECTIVE-REVIEW-TIER-ALIGN/）；
- 每次 runSubagent 传 agentName + 限定名 model；⑩ 派单表按 §13.2.2 逐字机械读取。

五 必读清单（新会话先读，全部完整路径）
1. docs/t027/CHECKLIST_v81.md 及 _EN（§一 锚点 / §五 O-1..O-5 登记原文）；
2. docs/DELIVERY_DIRECTIVE.md §5.2（四档选型矩阵）/ §7.5（分层额度与单价表）/ §8.2（档位映射与取档授权句）/ §0.3（版本行）/ §12.4（OB 表）/ §13.1.10 / §13.2.2；
3. docs/validation/hydrogen-control-pilot{,_EN}.md（O-2..O-4 原文）与 docs/validation/evidence/P97-HYDROGEN-CONTROL-PILOT/ 04-reviews/07-pilot-hydrogen.txt（Hydrogen 试水原文）；
4. docs/t027/REMAINING_SLICES{,_EN}.md P97 段（观察项 O-2..O-4 / O-5 登记行）与 P98 登记位；
5. docs/CODING_CONVENTIONS.md（.md = BOM + LF）。

六 实现要点（① 方案书须覆盖，缺一退回）
- 修订文本三件套：
  a. §5.2 深度档行：Neon 档位列与 §8.2 对齐（`Extra High`），并补可判性措辞（O-3：术语定义 / 优先级 / 「其余三档」补备选档 / 备选门槛与「最高风险轮」界线）；
  b. §8.2：取档授权句与 §5.2 对齐后逐字一致；
  c. §7.5：Neon 两档单价齐备 + 档位归属标注 + 备选额度口径首句措辞消除误读（O-5 #4）；
- O-5 #5（P87 证据包 numstat 陈旧）：**冻结件不修**，仅在本片报告与台账观察项处置中登记「已知、不追改」；
- **停点 1 触发即停**：报四项候选取用户裁定，不等额用完；
- ② 条款落地 → 条款 CI 绿 ⇒ ③ 起手；③ 台账观察项处置（O-2..O-4 / O-5 翻「已处置（P98）」，双语镜像）随本片提交；
- SELFTEST 806/806 必须不变（判据实现零改动）；若 ① 判定需动 selftest 夹具文本（仅文本同步），停点 1 报告并写明影响面。

七 文档同步更新（§13.1.18 硬规则，漏写即结构缺陷）
必改：
1. docs/DELIVERY_DIRECTIVE{,_EN}.md：头部 v1.83 + §0.3 新行 + §5.2 / §7.5 / §8.2 修订 + 附录 C 新块 C.0bf（六形状双语镜像）；
2. docs/t027/REMAINING_SLICES{,_EN}.md：P98 段新增（⬜→🔄→✅）+ 完成记录 + P97 段观察项 O-2..O-4 / O-5 处置翻「已处置（P98）」；
3. docs/validation/directive-review-tier-align{,_EN}.md（报告对，新增，双语严格镜像 + artifacts 围栏）；
4. 证据包 docs/validation/evidence/P98-DIRECTIVE-REVIEW-TIER-ALIGN/（00-briefing、input-list、freedom-list、01-arch、审查原文、MANIFEST / SHA256SUMS）。
不动：
5. CHECKLIST_v81{,_EN}（v82 再吸收）；tools/gates/**；internal/**；ci.yml；CONTRACTS / schema / fixtures / workflows；G5-b 闭集与模型白名单；ADR_REGISTER（① 判定）；P87 报告对与证据包（含 numstat，只登记不修）；s1-exit / p57 门表；DEV_PLAN；DOCUMENTATION_PLAN（v82 时代换）；hydrogen-control-pilot 报告对。

八 边界
- 只改 §5.2 / §7.5 / §8.2 相关行与版本包、台账、报告；不改任何判据实现与判定逻辑；不改 G5-b 闭集与白名单；不新增 ADR（① 判定）；不预设临时授权；不做 gitee；修订不得扩大解释到 §2.2 模型名单或 §7.5 额度数值（仅档位/单价标注与措辞）。

九 回执要求
- 首行 P98: STARTED；四要素（裁决人=用户、日期=2026-11-28、理由= O-2..O-4 处置方向「尽快处理」点名（独立小修订片）、依据=本起跑消息 + v81 快照 §五 + P97 报告观察项）；
- 锚点 + 模型/档位（主控 V4 Pro **Max**（用户裁定、不得作先例）/ ①⑤ V4 Pro **Max** / ②③ 文档退化 / ⑥ 跳过经 ⑦ 确认 / ⑦ Codex 常规档治理片行 1+3 / GLM-5.3-Flash ≥2）逐项复述；
- **停点 1 触发即停**：报四项候选取用户裁定，不得自行决定；
- 任一环报 Medium+ ⇒ 原位重跑至清零（§5.3），禁跳序；
- 单次输出 ≤ 10,000 字符（§13.1.16）；提交信息英文；完成后停于 ⑩，等同轮授权；仅 origin。

# 切片验证报告 · P68 REMAINING-RULINGS

切片标识：四项 R（R-04 / R-06 / R-07 / R-10）裁定建议 + 分阶段策略封板确认（评估类 `[SLICE]`，层 2 治理侧，M 档）。

起跑：2026-11-12；主控 = `DeepSeek V4.1 Flash (deepseek)`（默认，不登记 ADR）。

状态：步骤 0 已完成；① 阶段待启动。

## 1. 执行摘要

本片为评估类 `[SLICE]`（层 2 治理侧，M 档），2026-11-12 起跑；主控 `DeepSeek V4.1 Flash (deepseek)`（默认，未登记 ADR）。
流水线裁剪为 0 → ① → ③ → ④ → ⑤ → ⑦ → ⑧ → ⑨ → ⑧c → ⑩：无 ② 实现、无 ③ 代码测试，⑥ 独立扫描跳过。
① 五件：A 分阶段策略已封板（阶段 1 有效 / 阶段 2 不适用）；B 方案 a；C 方案 c；D 方案 a + 缺口登记；E 方案 a + c 纯记录形态。
批量辅助：GLM-5.3-Flash 实派 2 次（2418 / 1428 字符），记录于 `02-batch/glm-batch1.txt` 与 `glm-batch2.txt`。
审查链：③ `CONSISTENCY: PASS`（34 file:line，33 首轮命中 + 1 处 Low 行号修正）；④ `SELFTEST: PASS 737/737`、三作用域 `TOTAL_FAIL=0`；⑤ r1 `PASS`（1 Low 已整改）；⑥ 跳过经 ⑦ 确认；⑦ r1 `FINAL REVIEW: PASS`，0 findings。
计数：R = 1、D = 0；四项 R 的最终裁定由用户在 ⑩ 后同轮或分批作出。

- R = 1 依据（§9.2 计数单元「发现后的修复重跑」）：③ 自检发现 1 处行号偏差（`p57-closeout-assessment.md:64` → `:69`）→ 修正 → 重跑 34 锚点机检 1 次，计 1；⑤ 的 Low 采纳未触发环重跑，未计入。
- 异常登记（Low 级美观，接受现状）：`04-reviews/07-final-r1.md` 头部并存 3 行归档注记（内容同义、无真实缺陷，不触 G2 / G7-f / G8-b）；清理延后至 v55 生成时，不在本片改动证据包。
- ⑧c 后取证窗口例外（§5.1 路径③，2026-11-12）：⑩ 核查后因 §1 补登记需重取 ⑨ 读数，登记后重建 `tmp/P68-REMAINING-RULINGS/`，用毕逐件精确清理并回读。
- 临时件 D = 1（用户 2026-11-12 裁定记入）：⑩ r3 前把 tree 读数写到 `tmp\p68g.txt`（tmp 根下，非声明的 `tmp/P68-REMAINING-RULINGS/`），已即时删除；回读 `tmp` 顶层 = `gates-cache,.gitkeep`、非缓存文件 = 1（`.gitkeep`），零残留；D < 3，不停片。
- OB-11 例外（用户 2026-11-12 裁定）：R-07 完整建议属① 的实质性补充、不满足 OB-11 五要素，但不含规则语义变更 ⇒ 按 OB-37 元规则以 **ADR-042** 独立留痕（超出字面标准、系用户直接裁定、不得作先例），**不重跑 ⑤/⑦**。

## 2. 分阶段策略封板确认

**结论：已封板**（阶段 1 为现行有效判定标准；阶段 2 依用户 2026-11-09 裁定 A1「不适用」）。

| 锚点 | 原文要点 |
|---|---|
| `docs/validation/s1-exit.md:43` | AT-23「阶段 1 条件闭合达成…阶段 2 不适用，用户 2026-11-09 裁定 A1」 |
| `docs/validation/s1-exit.md:52` | S-G1 依据 = AT-23 阶段 1；「不得声称 S1 Exit 整体通过」 |
| `docs/TEST_STRATEGY.md:65` | AT-23 验收结论（F2 修复片更新；阶段 2 不适用 / 裁定 A1） |
| `docs/validation/r08-r09-ruling.md:103` | §4 分阶段策略；`:122` 三条升级条件；`:125` 第 2 条 = 用户明确裁定，主控不得自行决定 |
| `docs/t027/REMAINING_SLICES.md:93` | P61 表行「PARTIAL 收口…ADR-036」；段体 `:1736` |
| `docs/t027/CHECKLIST_v54.md:65` | F2-FIX 行（`57631a3`；AT-23 阶段 1 条件闭合达成；ADR-039） |
| `docs/DEV_PLAN.md:175` | F2-FIX 收口 ⇒「AT-23 阶段 1 条件闭合达成」 |

**后续片不需重审分阶段策略；不自行按阶段 2 收紧探针/校验标准。**（升级唯一入口 = 用户裁定 + 重跑关键场景。）

分阶段策略框架（Stage vs Layer / 影响面 R-08.1+R-08.2 / 升级触发三条）在用户侧资料（v5 交接文档 §三），仓库内 s1-exit / TEST_STRATEGY / 台账 P61 段仅承载封板结论；本报告不重复框架内容，只记指针。

## 3. R-04 建议（DR-13「已修复 + 观察」口径）

### 3.1 现状锚点

- `docs/t027/REMAINING_SLICES.md:1950`：DR-13「已修复（2026-11-01）……④ 若再现同一断言串即重开」。
- 同文件 `:1405` DR13-FIX 段 ✅；`:1413` 依赖行（复现 run `41253798647` 与 `41257990185`）。
- 待裁表四版无新增复现：v51:113 / v52:119 / v53:131 / v54:144；`CHECKLIST_v54.md:150`「维持待裁」。

### 3.2 候选方案

- **方案 a（收口为定义式口径，保留 ④ 重开条件）**——预算 0（仅改待裁表行）；时间 <0.5h；风险 低（需防「已收口」被读作永久关闭）。
- **方案 b（转长期观察项、移出待裁表）**——预算 0；时间 <0.5h；风险 中（重开触发失去已批授权，须另立条款）。
- **方案 c（补做复现探针 + push 取证）**——预算 1 次逐次声明探针 + 1 次双腿 CI；时间 1–2h；风险 中（CI 再现率低、无统计力，反而增 flake 噪声）。

### 3.3 建议

- 推荐 **方案 a**：D1/D2 修复 + 反例回归（先红后绿）+ D3 闭合 W1 已构成闭环，v51–v54 四版观察至 P67 零复现；方案 c 无判别力。**最终裁定由用户在 ⑩ 后同轮或分批作出。**
- **用户裁定（2026-11-12，同轮）：方案 a。**

## 4. R-06 建议（OB-55 归属与 T027-B 节奏）

### 4.1 现状锚点

- `docs/DELIVERY_DIRECTIVE.md:1091`：OB-55 处置列「**待议**：B 段立项时由 ① 阶段评估（独立切片 / 提前落 `Sandbox` 接口 / 维持 S2）」。
- `docs/DEV_PLAN.md:151` 立项登记；`:153` R-2 = OB-55 归属待 R-06；`:161`/`:163` P37/P38 基线；`:177`–`:185` P39–P44（`:185` P44 端到端接线为后续切片）。
- `docs/validation/p57-closeout-assessment.md:69` §3.3（R-2 行）；`:82`–`:89` §4.2（`:89` S-G8 依赖 T027-B）。

### 4.2 候选方案

- **方案 a（独立切片落地 B 段）**——预算 M 档（⑦ 1+2）≈3 审查单元；时间 3–5 片；风险 高（无准入候选时空转，且触 AT-23 重评）。
- **方案 b（提前落平台无关 `Sandbox` 接口）**——预算 S/M 档；时间 1–2 片；风险 高（与 B2 §8 第 8 条死代码禁令冲突，须先修该条）。
- **方案 c（维持 S2）**——预算 0；时间 0；风险 中（S-G1/S-G8 保持 BLOCKED，但风险集可控）。

### 4.3 建议

- 推荐 **方案 c**：候选不可用 + 死代码禁令 + S1 已无时间窗三重约束；若目标改为解锁 S-G1，须以「新准入候选」为前提转方案 a。**最终裁定由用户在 ⑩ 后同轮或分批作出。**
- **用户裁定（2026-11-12，同轮）：方案 b**（先补 P44 端到端接线 R-1 收口，再起 T027-B）；**R-1 补全本身由后续切片承载，P68 不推进**。

## 5. R-07 建议（MVA、S1 移出项与 S2 承诺收窄）

### 5.1 现状锚点

- 原文出处（仓库内）：`docs/validation/strategic-planning-2026-10-31.md` §2.1（`:17`–`:19`）与 §4 第 2 项（`:47`）；`evidence/STRATEGIC-PLANNING-2026-10-31/01-product-planning/review.md` Q1（`:6`–`:11`）/ Q2（`:13`）/ Q3（`:22`–`:29`）；`03-synthesis-astra/review.md` 冲突裁决表（`:21`「产品减法可以提请批准，但 B 的意见不能代替批准」）与 Q4（`:39`–`:53`）。
- 范围基线：`docs/S1_SUPPORT_MATRIX.md`（visible 行 = 降级保证、核心离线 CLI 不要求）；`docs/PRODUCT_REQUIREMENTS.md:71`（REQ-009）/ `:82`（REQ-020）/ `:69`（REQ-007）/ `:83`（REQ-021）/ `:81`（REQ-019）。
- 技术前提（详见 `01-arch/r07-refined.md` §0）：A 现状为「跑完才输出」、无 per-task/per-step 行、无 `--verbose`/`--progress`；B AI 子进程 stdout/stderr 只落 `logs/` 文件、无心跳、`--json` 仅在结束输出摘要；C 事件已结构化且 `events/state-events.jsonl` 为唯一全序源，但无订阅 API；D visible 的对话可视化由 Copilot Chat 面板提供，ProofRail 不渲染、不作权威事实。
- 附加：AT-13 未通过（`TEST_STRATEGY.md:69`；硬缺口 N-1 Windows 侧核验者独立性、N-6 工具链同形）；S2 语言冻结不影响 S1（REQ-020 的 S1 只要求 C+Go；generic/C/Go 已验证）。

### 5.2 候选方案

- **方案 a（全部批准）**：批 MVA 定义 + 三项 S1 移出（REQ-009 完整 TUI / REQ-019 自托管及 US6/AT-13 / REQ-007·021 黑箱子项及 US11/AT-24）+ S2 语言冻结。
- **方案 b（部分批准，4 个组合）**：b1 = 只动 REQ 的阶段归属列（X 标注推迟）+ S2 语言承诺收窄，不动 MVA 作为验收复合标准、不删任何 REQ、不动自托管；b2 = b1 + 批 MVA；b3 = 仅 S2 冻结 + TUI 移出；b4 = b2 + 自托管移出（= a 分两步）。
- **方案 c（暂缓）**：不批准，维持原 S1/S2 承诺。

### 5.3 影响（预算 / 时间 / 风险）与批准后回写

- a / b2 / b4：一次范围决策 + 一次回写片（M 档，⑦ 1+2），1–2 片；风险 高（MVA 改变验收口径；b4 使 AT-13 缺口外推）。
- b1：S 档可承载（⑦ 1+1），0.5–1 片；风险 中低（移出视图层与承诺面，S1 核心闭环与 AT-13 不动）。
- b3：S 档，0.5 片；风险 低（回写面最小，但 S1 Exit 压力不减）。c：0 成本，风险 中（问题后移）。
- 回写清单（**不在 P68 范围**，须另立切片）：`PRODUCT_REQUIREMENTS`（移出项阶段归属 + S2 收窄登记）、`TEST_STRATEGY`（AT 条目阶段归属）、`S1_SUPPORT_MATRIX`（S1 范围基线）、`s1-exit`（八门适用性）、`DEV_PLAN`（按 §11.2 判定）；`CHECKLIST_v54` 不改（v55 继承）。

### 5.4 建议

- 推荐 **方案 b1**：TUI 移出的可用性代价可由「事件流已在磁盘 + 补齐 tail 指引」低成本缓解（不需完整 TUI）；黑箱移出与 ADR-012 的降级保证一致；S2 冻结不影响 S1；自托管不批移出（AT-13 的 N-1/N-6 硬缺口不应失去验收压力）；MVA 待 S1 Exit 收口后再议。**最终裁定由用户在 ⑩ 后同轮或分批作出。**
- **用户裁定（2026-11-12，同轮）：方案 b1 + 附条件**——X 标注推迟；附条件 1 = TUI 移出以 S1 内补 `run --progress` 轻量流式开关为前提（另立 `S1-STREAM-PROGRESS`，P68 不启动）；附条件 2 = visible 调试指引写入 `OPERATIONS` + 三条边界（随 S1 文档片）。

## 6. R-10 建议（总预算与各片模型上限）

### 6.1 现状锚点

- `docs/DELIVERY_DIRECTIVE.md:757` §7.5 硬指标表（架构师 ≤1 / 预审 ≤1 复审 +1 / ⑥ ≤3 / ⑦ M 档 1+2 / 盲审另计 +1）；`:775` 成本梯度（Flash `0.15/0.6`、V4 Pro `0.66/1.98`）；`:1229` §13.1.8 主控模型策略。
- `docs/t027/CHECKLIST_v54.md:147`：「R-10 … 待裁（阻塞排片）」。
- **缺口声明**：仓库内无总预算数值（逐片额度与单价有，总盘子无）。

### 6.2 候选方案与边界

- **方案 a（维持 §7.5 现有额度，不触准则）**——预算 0；时间 0；风险 低-中（无总盘硬上限，由 ⑩ 逐片计量与「超限 ⇒ 升级用户裁决」兜底）⇒ **in-bounds**。
- **方案 b（§7.5 新增「总预算」行）**——触 §7.5 正文 ⇒ 准则修订；预算 治理片 M/L；时间 1–2 片；风险 高（与在飞片额度口径分裂）。**本片只建议、不落。**
- **方案 c（分层记录，不改准则）**——仅 ⑩ 回执追加「总预算核对」行、不设判据：预算 0；时间 0；风险 低。仅记录、无判据、门禁不可检 ⇒ 归 a 的配套（in-bounds）；一旦写入 §7.5 或新增可检判据 ⇒ 变 b（out-of-bounds）。
- **方案 d（按片重估各角色上限）**——触 §7.5 数值；预算 治理片；时间 1 片；风险 中（新老口径并存）。**本片只建议、不落。**

### 6.3 建议

- 推荐 **方案 a**（维持现有额度，不触准则）+ **c 的纯记录形态**（不写入准则；c 的落地形态仍需用户裁定）。b/d 标为需用户裁定、本片不落；本片不触发停升级。**最终裁定由用户在 ⑩ 后同轮或分批作出。**
- **用户裁定（2026-11-12，同轮）：方案 a**（维持 §7.5 现有额度，不触准则）。


## 7. 已知边界

- **更正（2026-11-12，用户指出）**：MVA / S2 承诺的实质内容在仓库内（`strategic-planning-2026-10-31.md` §2.1 `:17`–`:19` 与 §4 第 2 项 `:47`；两份 review 原文 Q1/Q3 与冲突裁决表），v26 快照仅在用户侧；R-07 完整建议已按仓库内原文重出（`01-arch/r07-refined.md`）。
- 技术债观察（只登记、不修、不立 OB）：`prfrail run` 无 per-task/per-step 输出、无 `--verbose`/`--progress`、无心跳；AI 子进程 stdout/stderr 只落 `logs/` 文件；`events/state-events.jsonl` 虽为结构化唯一全序源，但无订阅 API，且 README / OPERATIONS 无运行期 tail 指引。是否另立切片由用户裁定。
- S2 的 TUI 工作量拆分（视图层 vs 内容层）在仓库内**无文档依据**，本报告不作估算。
- 仓库内无总预算数值，预算口径以用户侧为准，正文不作数值引用。
- 分阶段策略的完整框架在用户侧，仓库只承载封板结论与四处落地锚点。
- 本片只出建议、不裁定：四项 R 均为建议方案，未经用户裁定前不具备准则效力。
- 建议基于 2026-11-12 时点文档；用户裁定后若事实变化，需重新评估后才可援引。

## 8. 明确未执行事项

- 未改产品代码 / 契约 / schema / fixtures / CI；未改准则正文；未改 v50–v54 快照；未替用户裁定；未 commit / push；未启动新切片；未生成 v55。

## 8a. ⑧a 自查

- ① 成本表机械审计：无（不适用）——本片无成本台账文件；载荷清单见 `MANIFEST.md`。
- ② 「已处置」`file:line`：有——34 处全量机检，33 首轮命中 + 1 处修正（`p57-closeout-assessment.md:64` → `:69`），见 `03-test/consistency-r1.md`。
- ③ 旧口径 grep：无（不适用）——本片未改准则正文与口径，未执行旧口径 grep。
- ④ ⑥/⑦ 每轮 Info / Low 落地：有——⑤ r1 的 1 项 Low 已整改（见 `04-reviews/05-prereview-r1.md`）；⑥ 跳过、⑦ r1 0 findings，无另外条目。

- ⑥ 跳过：本片为评估类纯文档片，不触代码、不触 §7.2.1 ③ 四类语义 ⇒ 可跳过（待 ⑦ 确认）；若 ⑦ 判定不成立，改跑 ⑥ 并以 `sixstage` 围栏替换本行。

```artifacts
docs/validation/evidence/P68-REMAINING-RULINGS/00-briefing.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/MANIFEST.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/SHA256SUMS.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/freedom-list.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/input-list.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/01-arch/design-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/01-arch/r07-refined.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/02-batch/glm-batch1.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/02-batch/glm-batch2.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/03-test/consistency-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P68-REMAINING-RULINGS/06-gates/gates-r1.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/06-gates/gates-r2.txt	repo
docs/validation/evidence/P68-REMAINING-RULINGS/06-gates/step0-selftest-r1.txt	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
```

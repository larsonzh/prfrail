切片：P59 DIRECTIVE-SCOPE-GUARD 防越界防线（治理类 · [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 完整阅读 docs/t027/CHECKLIST_v48.md；以刚才锚点回执确认的实际 HEAD 为当前锚点。v48 快照记录的 8d30e66 是生成时点；若实际 HEAD = 33bd849，属 v48 落地后的预期前进。
2. 5 个锚点必须一致：准则 v1.62；SELFTEST 703/703；index/tree/ci 三作用域 TOTAL_FAIL=0。除 HEAD/origin 的预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况。
3. **主控模型声明**：本片按 v48 §七 第 8 条的精神使用 DeepSeek V4 Pro（默认 Flash 遇"底层"可切 Pro；本片约束的正是 Flash 在 P58 的三次越界，独立性需要 Pro）。**步骤 0 必须按 OB-37 元规则把此次切换登记为 ADR**（四要素：裁决人=用户；日期=起跑日；理由=本片是 D-6 的约束对象本身，Flash 自写约束自己独立性不足；依据=本轮起跑消息）。回执首行标 [自执行]。
4. 记忆唤醒：治理类切片 ⑨ 后必停 ⑩；单次输出 ≤10,000 字符；不自行 commit/push；不替用户裁定 R-12；不启动报告中的任何新切片；**本片自身即 D-6 的实证对象，全程不得自行升版、不得伪造授权、不得超额派单**（这正是本片要固化的三条）。

【〇 本片环境声明与流水线裁剪】
1. 环境：纯本地；**不开启 Ubuntu VM**、**不开启网络代理**；不改产品代码、CI、契约、schema、fixtures；本片只改准则、台账、ADR、报告与证据包；若 R-12 裁定"判据"路径，则 `tools/gates/**` 与 `selftest.js` 属本片边界内，其余仍不动。
2. 架构双打（Sol）：**不启用**（非 P54 那类不可测语义片）。
3. 流水线：本片是**两段式**——
   - 阶段 A：步骤 0 → ① 出三案对比表 → **停在 ① 等用户裁定 R-12**（不得自行推断裁定）
   - 阶段 B：R-12 裁定后 → ② 实现（按裁定路径）→ ③ 测试（判据路径才需要）→ ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描（按 §7.2 判定）→ ⑦ 终审 → ⑧ 文档装配 → ⑨ 三作用域 + SELFTEST → ⑧c 清理 → ⑩ 停点
4. 阶段 A 结束时**不得进入 ②**；阶段 B 开始时不得回改阶段 A 的三案对比表（只按用户裁定执行）。

【一 授权边界（角色与模型）】
- ① 架构（出三案对比表）：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试（判据路径才派）：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3；按 §7.2 判定是否可跳过）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+3，治理片高配）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（若走②③可写路径）
- ⑨ / ⑧c / ⑩：主控执行
- **不启用 Sol；不派批量辅助**（本片不需要，仍按 v48 §七 第 7 条下限报"0 次"并说明理由）
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`

【二 执行步骤与核心要求】

## 阶段 A：① 出三案对比表并停等 R-12

### 步骤 0 待办
1. 建 `tmp/P59-DIRECTIVE-SCOPE-GUARD/` 与 `docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/`
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF）
3. 台账 `REMAINING_SLICES{,_EN}.md` 登记 P59（🔄 进行中，R/D=0；承载 G-3）；不改 G-3 行的内含（v48 已写）
4. 落盘 `freedom-list.md`、`input-list.md`、`MANIFEST.md`、`SHA256SUMS.txt`
5. **ADR 登记（步骤 0 必做）**：按 OB-37 元规则在 `ADR_REGISTER{,_EN}.md` 新增 **ADR-032**（主控切 V4 Pro）；四要素见【〇·零】第 3 条；不修改 ADR-031
6. 落盘 `docs/validation/p59-directive-scope-guard.md` 的骨架（标题 + artifacts 围栏占位）

### ① 架构（出表，不做实现）
1. 通读并引用（给 file:line）：
   - v48 §六 R-12 的"执行侧 + 审查侧"双层定义
   - v48 §八 纪律 Y（候选）
   - v48 §十二 的 D-6 三次边界问题
   - P58 报告 §8 的 D-1 至 D-6
   - 准则 §0.2（MAJOR/MINOR 定义）、§1.7.2（升级边界）、§1.7.4（临时授权五要素）、§12.4 的 OB-137 行、§12.5 的相关缺口行、附录 B.5（⑤ 预审模板）、§6.3 的 G7-f（⑥ 输出契约）
2. **执行侧三案**（针对 D-6：切片内不得自行升版 / 伪造授权 / 超额派单）：
   - **A1 记录态约束**：在 §1.7.2 的升级边界段末加一段记录态条款；不新增判据、不加机械检查；依赖 ⑤/⑦ 的人工核验
   - **A2 正文条款**：在 §0.2 的版本与变更段或 §1.7.2 的升级边界段把三条写成强制条款；改动面比 A1 大，但仍不加机械判据
   - **A3 判据**：在 §6.3 新增一个 G 判据（或在 G6 下加子判据），机械检测"变更集是否含未授权升版 / 未授权授权字样 / 超额派单"；需要落 `tools/gates/**` 与 `selftest.js`
3. **审查侧三案**（针对 OB-137：⑤ 不核用户授权事实真伪）：
   - **B1 ⑤ 必核面清单增设一项**：在 §4.2.1 的"最低核验动作"六条后加一条"凡引用用户授权或裁定，须可回溯到同轮用户原话"；不改判据
   - **B2 判据**：在 §6.3 新增判据或在 G7 下加子判据，机械检测"报告/台账中出现的'用户授权'字样是否对应同轮用户原文"；需要落 `tools/gates/**` 与 `selftest.js`
   - **B3 记录态**：只在 §12.4 的 OB-137 处置列写"由 ⑤ 必核面人工承担"，不加正文条款
4. **六案的统一对比维度**：
   - 改哪节（file:line）
   - 是否新增机械判据（是 / 否）
   - 是否改动 `tools/gates/**`（是 / 否）
   - 是否可被 Flash 自写规则绕过（风险评级）
   - 是否与 §0.1 的"执行纪律权威 = 本准则"一致
   - 落地成本（S / M）
   - 对 ⑤/⑦ 的额外工作量
5. **纪律 Y 的归属**：单独一段，给出"纪律 Y 升为正文条款 / 保留在 §八 候选 / 并入 A1 记录态"三案对比
6. **规模预估**：按 (执行侧 A 案 × 审查侧 B 案 × 纪律 Y 归属案) 给规模矩阵（S / M），标出推荐档位
7. **① 输出格式**：出表落盘 `docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/design-r1.md`；报告 `docs/validation/p59-directive-scope-guard{,_EN}.md` 的 §2 "① 三案对比表"只放表格摘要 + 指向 design-r1.md
8. **停在 ①**：出表后回执首行 `DESIGN: DONE`，正文说明"停等 R-12 裁定，未进入 ②"，并给出推荐（推荐必须标注"仅供参考，最终由用户裁定"）

## 阶段 B：R-12 裁定后继续

### ② 实现（按裁定路径）
- 若裁定 A1 / A2 / B1 / B3（纯文档）⇒ 主控自执行，或派 `prfrail-implementer` 落盘
- 若裁定 A3 / B2（判据）⇒ 派 `prfrail-implementer` 改 `tools/gates/**`，接口签名不变，只加子判据
- **禁止**在同一次提交里"判据 + 正文条款"混合落地（那等于擅自扩大裁定范围）；用户裁定了哪条就只做哪条
- 纪律 Y 的归属按用户裁定执行；未裁定时保持 §八 候选，不升不降

### ③ 测试（仅判据路径需要）
- 派 `prfrail-tester`，为新增判据写正反夹具，接入 `selftest.js`
- 纯文档路径跳过 ③，在报告中如实登记"无测试对象"

### ④ 门禁
- 全量 `go build/vet/test ./...` + `gate.js --all --scope=index`、`tree`、`ci`
- 判据路径额外跑 `--selftest` 并登记断言增量

### ⑤ 预审
- 派 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- 重点审：条款是否覆盖 D-6 三条 / OB-137 的核心；是否引入新歧义；是否与 §0.2 / §1.7.2 / §1.7.4 现有条款冲突；判据路径是否与 G6 / G7 现有子判据重叠
- 报 Medium+ ⇒ 原位整改后重跑 ⑤

### ⑥ 独立扫描
- 按 §7.2 判定是否可跳过（本片改执行纪律语义 ⇒ ③ 不成立 ⇒ ⑥ 必做）
- 派 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`，额度 3

### ⑦ 终审
- 派 `independent-reviewer` + `GPT-5.3-Codex (copilot)`，额度 1+3（治理片高配）
- 盲审：本片触及"授权 / 纪律"语义 ⇒ **必抽 1 次盲审**（不附任何清单）
- 首轮不附 ⑤⑥ 清单

### ⑧ 文档装配
- 派 `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（若走②③可写路径）
- 或主控自执行（纯文档路径）
- 按【四 文档同步更新】逐项处置

### ⑨ 三作用域 + SELFTEST
- 必须 `TOTAL_FAIL=0` 与 `SELFTEST 703/703`（判据路径下 SELFTEST 会增加，据实登记）

### ⑧c 收尾清洁
- 清 `tmp/P59-DIRECTIVE-SCOPE-GUARD/`；保留 `docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/`

### ⑩ 停点
- 停在 ⑩ 等同一轮授权；严禁自行 commit/push

【三 执行节奏与红线】
1. **阶段 A 停等 R-12 裁定**：① 出表后必须停，不得自行推断裁定；不得把"推荐"写成"结论"
2. ⑤ / ⑥ / ⑦ 报 Medium+ ⇒ 修复后**原位重跑该环**，不得跳步
3. ⑩ 是唯一常规停点；严禁自行 commit/push
4. 本片是 D-6 的实证对象：全程不得自行升版、不得伪造授权、不得超额派单。**若本片自身出现这三类越界，按 D-6 的同样标准处理（登记 D 项 + 报告披露）**
5. 不得新增 OB 编号；不得把评估建议写成"已裁决"
6. 临时件一律 `tmp/P59-DIRECTIVE-SCOPE-GUARD/`；只增不删；⑧c 统一清理
7. **若 R-12 裁定的组合"判据 + 正文条款"同时落地**，在同一提交中分块呈现（正文条款块 / 判据块），不得混在一起

【四 文档同步更新（不得遗漏）】

在 ⑧ 阶段，必须**逐项**处置下列文档；不适用项必须显式声明"不适用 + 理由"，不得跳过：

1. `docs/DELIVERY_DIRECTIVE{,_EN}.md`：
   - §0.2（MAJOR/MINOR 定义）：**如用户裁定 A2 涉及此段则改；否则不动**
   - §1.7.2（升级边界）：**如用户裁定 A1 / A2 涉及此段则改；否则不动**
   - §4.2.1（⑤ 最低核验动作）：**如用户裁定 B1 则加第 7 条；否则不动**
   - §6.3（判据本体）：**如用户裁定 A3 / B2 则加判据；否则不动**
   - §12.4（OB-137 处置列）：**所有路径都要改**——按用户裁定把 OB-137 从 `观察中` 翻为 `已处置` 或更新处置列指针；不得翻其他 OB
   - §12.5（缺口与未来扩展）：**如 A3 / B2 涉及判据则相应更新；否则不动**
   - §八（纪律 Y）：**按用户裁定升为正文条款 / 保持候选 / 并入 A1**
   - 头部版本号与 §0.3 变更日志：**任何条款改动都要升版 + 追加 §0.3 行**（这是本片的正当升版，不是 D-6 越界；但要按 R-12 裁定的形态执行，不得自行组合）
   - **禁止**：不得在本片顺便改与 D-6 / OB-137 无关的条款

2. `docs/t027/REMAINING_SLICES{,_EN}.md`：
   - P59 段：闭环时改标题为 `✅ 已完成`，补完成记录行（只写报告路径），状态表行同步
   - G-3 行：把"待用户点名"改为"由 P59 闭环"；如需补"结论见 P59 报告"指针
   - 不改 P7 行（v48 裁定："P7 行不改"）

3. `docs/validation/p59-directive-scope-guard{,_EN}.md`：新建（核心交付物）

4. `docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/`：证据包（`00-briefing.md` / `01-arch/design-r1.md` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）

5. `docs/ADR_REGISTER{,_EN}.md`：
   - **步骤 0 必做**：新增 ADR-032（主控切 V4 Pro），四要素齐备
   - **如用户裁定"判据 + 正文条款同时落地"**：追加一条 ADR-033 记录该组合决策
   - 不改 ADR-031

6. `docs/DEV_PLAN{,_EN}.md`：**不适用 + 理由（治理片不回写 `DEV_PLAN`，按 §11.2）**

7. `docs/ARCHITECTURE{,_EN}.md`：**不适用 + 理由（本片不涉及模块边界、数据流或存储）**

8. `docs/SECURITY{,_EN}.md`：**不适用 + 理由（本片不涉及权限模型、威胁矩阵或秘密处理）**

9. `docs/CONTRACTS{,_EN}.md`：**不适用 + 理由（本片不涉及 wire、schema、状态机、判据语义；若 A3 / B2 落判据，判据本体属 §6.3，不在 CONTRACTS）**

10. `docs/S1_SUPPORT_MATRIX{,_EN}.md` 与 `docs/S1_RELEASE_NOTES{,_EN}.md`：**不适用 + 理由（无平台、支持层级、EOL 或发布变更）**

11. `docs/TEST_STRATEGY{,_EN}.md`：**不适用 + 理由（若 A3 / B2 落判据，判据属 §6.3 的交付自检判据，不在 AT 体系；若用户裁定涉及 AT 则相应更新）**

12. `docs/t027/CHECKLIST_v48{,_EN}.md`：**不适用 + 理由（v48 已落地；本片不触发重新生成，除非用户裁定组合形态涉及队列变化——若涉及，另起 v49 快照，本片不生成）**

13. `docs/PRODUCT_REQUIREMENTS{,_EN}.md`、`docs/BUSINESS_WORKFLOWS{,_EN}.md`、`docs/RFC-proofrail-unattended-ai-engineering-product.md`、`docs/DOCUMENTATION_PLAN{,_EN}.md`、`docs/CODING_CONVENTIONS{,_EN}.md`、`docs/OPERATIONS{,_EN}.md`、`docs/INSTALLATION{,_EN}.md`、`docs/INSTALLATION_PLAN{,_EN}.md`：**不适用 + 理由（本片只触执行纪律层与判据本体，不触产品范围、业务流程、安装模型、编码规范、操作手册）**

14. `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`：**不适用 + 理由（本片不是结项报告，不触模板）**

15. 报告本体（`docs/validation/p59-directive-scope-guard{,_EN}.md`）：**⑧ 阶段必做**——按 §11.3 完整结构；`artifacts` 围栏按纪律 G（路径排序、TAB 分隔、CN/EN 逐字节相同）；`sixstage` 围栏按 G7-f（若 ⑥ 实际派了）

**文档同步的机械核验**（⑨ 阶段）：
- `G3(a)` 六形状全等（CN/EN）
- `G1-b` 版本元数据（头部与 §0.3 一致）
- `G6-6` 台账状态词闭集（OB-137 翻词必须在闭集内）
- `G6-7` 声明/台账对账（§0.3 声明的 OB 编号必须在 §12.4 存在）
- `G5(d)` 暂存集 ⊆ `artifacts` 围栏声明的集合
- `G5(f)` index 作用域下工作区干净

【五 ⑩ 回执要求】
1. 派单表（`agentName` 逐字机械读取；若 debug 日志无 transcript，须显式标注"未机械校验"并作为已知事实登记）
2. ① 三案对比表摘要 + 推荐（**推荐仅供参考，不得写成结论**）
3. R-12 裁定后的执行摘要（按裁定路径：改了哪节 / 加了什么判据 / 有无测试）
4. ⑤ / ⑥ / ⑦ 结论摘要（含整改与重跑）
5. 三作用域 + SELFTEST 原始行（判据路径下 SELFTEST 若增加，据实报）
6. 文档同步 `git diff --stat`（逐项覆盖【四】的 15 项；不适用项列理由）
7. 明确未执行事项：未改产品代码 / CI / 契约 / schema / fixtures；未替用户裁定 R-12；未 commit/push；未启动新切片
8. 待授权项：提交授权；若组合路径落地，是否一并登记为 ADR-033
9. **本片自身越界登记**：若本片出现 D-6 类越界（自行升版 / 伪造授权 / 超额派单），如实登记并披露；未发生也要明确写"无"

【六 下次起跑前的必要提示】
P60 `DIRECTIVE-JUDGE-BLINDSPOTS`（仅 OB-138，S 档）不依赖 P59；AT23 不依赖 P59 / P60，只依赖 R-08 / R-09 裁定。P59 收尾后决策侧可并行推 R-08 / R-09。

现在开始执行 P59 步骤 0；阶段 A 出表后停等 R-12 裁定，不得进入 ②。

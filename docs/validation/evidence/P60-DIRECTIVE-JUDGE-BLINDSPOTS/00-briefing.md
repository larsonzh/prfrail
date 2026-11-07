切片：P60 DIRECTIVE-JUDGE-BLINDSPOTS 判据盲区收口（治理类 · [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 完整阅读 docs/t027/CHECKLIST_v49.md；以实际 HEAD 为当前锚点。v49 快照记录的 75b90a8 是生成时点；若实际 HEAD 已前进（例如 v49 D 口径小修的 30dc4a3），属预期前进。
2. 5 个锚点必须一致：准则 v1.63；SELFTEST 703/703；index/tree/ci 三作用域 TOTAL_FAIL=0。除 HEAD/origin 的预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况。
3. 主控模型声明：本片使用 DeepSeek V4 Pro（OB-138 与 P58 D-6 同源；若 ① 裁定 J2 属判据本体变更；若 J1 属历史重排；均需独立视角）。步骤 0 必须按 OB-37 元规则把此次切换登记为 ADR-034（四要素：裁决人=用户；日期=起跑日；理由=判据覆盖盲区与 D-6 同源、判据本体/历史重排均需独立视角；依据=本轮起跑消息）。回执首行标 [自执行]。
4. 记忆唤醒：治理类切片 ⑨ 后必停 ⑩；单次输出 ≤10,000 字符；不自行 commit/push；不替用户裁定 J1/J2/J3；不启动报告中的任何新切片；**本片自身受 §1.7.2 三条 bullet 约束（不得自行升版 / 不得伪造授权 / 不得超额派单）**。

【〇 本片环境声明与流水线裁剪】
1. 环境：纯本地；**不开启 Ubuntu VM**、**不开启网络代理**；不改产品代码、CI、CONTRACTS、schema、fixtures；本片只改准则、台账、ADR、报告、证据包，若裁定 J2 则 `tools/gates/**` 与 `selftest.js` 属本片边界内。
2. 架构双打（Sol）：不启用。
3. 流水线：本片是**两段式**——
   - 阶段 A：步骤 0 → ① 出三案对比表 → 停在 ① 等用户裁定 J1 / J2 / J3
   - 阶段 B：裁定后 → ② 实现 → ③ 测试（J2 路径才需要）→ ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描（§7.2 ③ 不成立 ⇒ 必做）→ ⑦ 终审 → ⑧ 文档装配 → ⑨ 三作用域 + SELFTEST → ⑧c → ⑩
4. 阶段 A 结束时不得进入 ②；阶段 B 开始时不得回改阶段 A 的三案对比表（只按用户裁定执行）。

【一 授权边界（角色与模型）】
- ① 架构（出三案对比表）：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（J1/J3 纯文档可主控自执行；J2 判据必派）
- ③ 测试（J2 才派）：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3；§7.2 ③ 不成立 ⇒ 必做）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+3，治理片高配；按 §7.7 判定是否触发盲审）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（若走②③可写路径）
- ⑨ / ⑧c / ⑩：主控执行
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`

【二 执行步骤与核心要求】

## 阶段 A：① 出三案对比表并停等裁定

### 步骤 0 待办
1. 建 `tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/` 与 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF）
3. 台账 `REMAINING_SLICES{,_EN}.md` 登记 P60（🔄 进行中，R/D=0；承载 G-4）；不改 G-4 行的内含（v49 已写）
4. 落盘 `freedom-list.md`、`input-list.md`、`MANIFEST.md`、`SHA256SUMS.txt`
5. **ADR 登记（步骤 0 必做）**：按 OB-37 元规则在 `ADR_REGISTER{,_EN}.md` 新增 **ADR-034**（主控切 V4 Pro）；四要素见【〇·零】第 3 条；不修改 ADR-031 / ADR-032 / ADR-033
6. 落盘 `docs/validation/p60-directive-judge-blindspots.md` 的骨架（标题 + artifacts 围栏占位）

### ① 架构（出表，不做实现）
1. 通读并引用（给 file:line）：
   - v49 §四 的 G-4 行；v49 §五 的 OB-138 说明
   - 准则 §12.4 的 OB-138 行
   - 准则 §6.3 的 G1(e) 判据原文、G5(e) 判据原文
   - 准则 §11.1（历史不重写）、§11.2（回写口径）
   - `docs/t027/REMAINING_SLICES.md` 的 P40–P58 块内 `| 完成记录 | … |` 行（逐条列出行号与内容）与 `## 完成记录（追加）` 表的现有格式（四列：日期 / 切片 / 完成摘要 / 证据）
   - `tools/gates/lib/criteria/g1.js` 与 `g5.js` 中 G1(e) / G5(e) 的实现（若有）
2. **三案**（针对 OB-138：块内完成行不在 append 表覆盖内）：
   - **J1 迁入 append 表**：把 P40–P58 的块内 `| 完成记录 | … |` 行迁入 `## 完成记录（追加）` 表，按日期 / 切片 / 完成摘要 / 证据 四列重建。
     - 关键争议：这是**历史重排**还是**格式统一**？按 §11.1「历史不重写」，迁移不改事实内容，但改结构，需用户明示授权；三案对比表须明确标注「J1 需用户在裁定里明示授权碰历史行」。
     - 数据补全：块内行只有路径没有日期/切片/摘要，需要从同一块的 `状态` / `目标` / 报告路径 等处推导，推导规则必须可机械核验。
   - **J2 扩判据识别块内字段行**：在 `tools/gates/**` 加子判据，识别 `| 完成记录 | … |` 字段行并纳入 G1(e) / G5(e) 覆盖。
     - 关键决策：新判据落哪个 G 家族（G1 子判据？G5 子判据？还是新 G？）；与 G1(e) / G5(e) 现有作用域是否重叠。
     - 向前规则：J2 不动历史行，只让未来 P 系块内完成行被覆盖；已存在的 P40–P58 历史行仍不判（`G1(e)` / `G5(e)` 的「变更行」语义不变）。
   - **J3 记录态收口**：只在 §12.4 OB-138 处置列写「由人工承担」，不加正文、不加判据、不迁历史行。
     - 代价：OB-138 长期挂在「观察中」，直到某个治理片顺手收口；与 P58 D-6 的「先记录后补救」模式同型。
3. **三案的统一对比维度**（每个方案给出）：
   - 改哪些文件（file:line）
   - 是否新增/改动机械判据（是 / 否）
   - 是否改 `tools/gates/**`（是 / 否）
   - 是否触及历史行（是 / 否；若是，触碰 §11.1 的哪一层）
   - 对 `SELFTEST` 计数的影响（增 / 不变）
   - 对 P40–P58 后续回写的影响（未来若修改这些切片，是否受 §11.1 约束）
   - Flash 自写规则绕过风险（低 / 中 / 高）
   - 落地成本（S / M）
   - ⑤⑦ 额外工作量
4. **① 输出格式**：出表落盘 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/design-r1.md`；报告 `docs/validation/p60-directive-judge-blindspots{,_EN}.md` 的 §2 「① 三案对比表」只放表格摘要 + 指向 design-r1.md
5. **停在 ①**：出表后回执首行 `DESIGN: DONE`，正文说明「停等裁定 J1 / J2 / J3，未进入 ②」，并给出推荐（推荐必须标注"仅供参考，最终由用户裁定"）

## 阶段 B：裁定后继续

### ② 实现（按裁定路径）
- J1（纯文档）⇒ 主控自执行或派 `prfrail-implementer`，**须在报告 §1 明确声明「本片 J1 用户明示授权触碰历史行」**；不得顺手改历史行的事实内容
- J2（判据）⇒ 派 `prfrail-implementer`，接口签名不变，只加子判据；须同步 `selftest.js` 正反夹具
- J3（记录态）⇒ 主控自执行，仅改 §12.4 OB-138 处置列
- **禁止**：在同一次提交里"判据 + 历史重排"混合落地（那等于擅自扩大裁定范围）；用户裁定了哪条就只做哪条
- **禁止**：J2 路径下顺带改 `g6-whitelist.txt` 或 `cjk-newwords.txt`（白名单只随裁定路径所需的最小改）

### ③ 测试（仅 J2 需要）
- 派 `prfrail-tester`，为新增判据写正反夹具，接入 `selftest.js`
- J1 / J3 跳过 ③，在报告中如实登记"无测试对象"

### ④ 门禁
- 全量 `go build/vet/test ./...` + `gate.js --all --scope=index`、`tree`、`ci`
- J2 路径额外跑 `--selftest` 并登记断言增量

### ⑤ 预审
- 派 `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- 重点审：
  - J1 路径：迁入的完成行四列是否可从块内事实推导、是否改动了事实内容
  - J2 路径：新子判据是否与 G1(e) / G5(e) 现有作用域重叠、假阳性风险
  - J3 路径：§12.4 处置列措辞是否如实反映"OB-138 仍未机械化"
  - 三案均须核对：是否引入新歧义、是否与 §11.1 / §11.2 冲突
- 报 Medium+ ⇒ 原位整改后重跑 ⑤

### ⑥ 独立扫描
- 按 §7.2 判定：本片改执行纪律语义（G 判据或历史行）⇒ ③ 不成立 ⇒ ⑥ 必做
- 派 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`，额度 3

### ⑦ 终审
- 派 `independent-reviewer` + `GPT-5.3-Codex (copilot)`，额度 1+3（治理片高配）
- 按 §7.7 判定是否触发盲审：本片触及判据语义（J2）或历史行（J1）⇒ **必抽 1 次盲审**（不附任何清单）
- 首轮不附 ⑤⑥ 清单

### ⑧ 文档装配
- 派 `prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（若走②③可写路径）
- 或主控自执行（J1 / J3 纯文档路径）
- 按【四 文档同步更新】逐项处置

### ⑨ 三作用域 + SELFTEST
- 必须 `TOTAL_FAIL=0` 与 `SELFTEST 703/703`（J2 路径下 SELFTEST 会增加，据实登记）

### ⑧c 收尾清洁
- 清 `tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`；保留 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`

### ⑩ 停点
- 停在 ⑩ 等同一轮授权；严禁自行 commit/push

【三 执行节奏与红线】
1. **阶段 A 停等裁定**：① 出表后必须停，不得自行推断 J1 / J2 / J3 选择
2. ⑤ / ⑥ / ⑦ 报 Medium+ ⇒ 修复后**原位重跑该环**
3. ⑩ 是唯一常规停点；严禁自行 commit/push
4. 本片受 §1.7.2 三条 bullet 约束：不得自行升版、不得伪造授权、不得超额派单。**若出现违反，按 §1.7.2 同样标准处理（登记 D 项 + 报告披露）**
5. 不得新增 OB 编号；不得把评估建议写成"已裁决"
6. 临时件一律 `tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`；只增不删；⑧c 统一清理
7. **若裁定组合为 J1 + J2 同时落地**，在同一提交中分块呈现（J1 历史迁入块 / J2 判据块），不得混在一起

【四 文档同步更新（不得遗漏）】

在 ⑧ 阶段，必须**逐项**处置下列文档；不适用项必须显式声明"不适用 + 理由"，不得跳过：

1. `docs/DELIVERY_DIRECTIVE{,_EN}.md`：
   - §6.3 的 G1(e) / G5(e)：**仅 J2 路径改**（加子判据识别块内字段行）；J1 / J3 不动
   - §6.3 的脚本化状态段：**仅 J2 路径改**（新增子判据后同步状态段）
   - §12.4 的 OB-138 行：**所有路径都要改**——按裁定翻处置列状态（J1 → `已处置`（历史已迁）；J2 → `已处置`（判据已扩）；J3 → 维持 `观察中` 并补"由人工承担"）；不得翻其他 OB
   - §12.5 缺口行：**仅 J2 路径相应更新**（若有相关缺口行）；J1 / J3 不动
   - 头部版本号与 §0.3 变更日志：**仅 J2 路径升版 + 追加 §0.3 行**（判据本体变更属 MINOR 必要）；J1 / J3 属"不改规则语义" ⇒ **不升版**（除非 J1 用户明示要求升版记录历史迁入，否则 §12.4 追加指针不构成升版必要）
   - **禁止**：不得在本片顺便改与 OB-138 无关的条款

2. `docs/t027/REMAINING_SLICES{,_EN}.md`：
   - P60 段：闭环时改标题 `✅ 已完成`，补完成记录行（只写报告路径），状态表行同步
   - G-4 行：标"由 P60 闭环（YYYY-MM-DD；结论见报告对）"
   - **J1 路径**：P40–P58 的块内 `| 完成记录 | … |` 行迁入 `## 完成记录（追加）` 表；迁入须由用户明示授权；报告 §1 须如实登记"触碰历史行的用户授权范围"
   - **J2 路径**：P 系块内完成行保持原样（判据是向前规则，不动历史）
   - **J3 路径**：台账不动
   - 不改 P7 行（v48 裁定："P7 行不改"）

3. `docs/validation/p60-directive-judge-blindspots{,_EN}.md`：新建（核心交付物）

4. `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/`：证据包（`00-briefing.md` / `01-arch/design-r1.md` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）

5. `docs/ADR_REGISTER{,_EN}.md`：
   - **步骤 0 必做**：新增 ADR-034（主控切 V4 Pro），四要素齐备
   - **若裁定 J1**：追加 ADR-035 记录"J1 用户明示授权触碰历史行的决策依据"（记录用户授权原文、触碰范围与理由）
   - 不改 ADR-031 / ADR-032 / ADR-033

6. `docs/DEV_PLAN{,_EN}.md`：**不适用 + 理由（治理片不回写 `DEV_PLAN`，按 §11.2）**

7. `docs/ARCHITECTURE{,_EN}.md`：**不适用 + 理由（本片不涉及模块边界、数据流或存储）**

8. `docs/SECURITY{,_EN}.md`：**不适用 + 理由（本片不涉及权限模型、威胁矩阵或秘密处理）**

9. `docs/CONTRACTS{,_EN}.md`：**不适用 + 理由（本片不涉及 wire、schema、状态机、判据语义；若 J2 落判据，判据本体属 §6.3，不在 CONTRACTS）**

10. `docs/S1_SUPPORT_MATRIX{,_EN}.md` 与 `docs/S1_RELEASE_NOTES{,_EN}.md`：**不适用 + 理由（无平台、支持层级、EOL 或发布变更）**

11. `docs/TEST_STRATEGY{,_EN}.md`：**J1 / J3 不适用 + 理由（判据属 §6.3 的交付自检判据，不在 AT 体系）；J2 若涉及 AT 则相应更新**（本片 J2 只动 §6.3 的 G 判据，不触 AT 体系 ⇒ 不适用 + 理由）

12. `docs/t027/CHECKLIST_v49{,_EN}.md`：**不适用 + 理由（v49 已落地；v50 由下一片生成时一并处理 P60 闭环的队列变化，本片不生成 v50）**

13. `docs/PRODUCT_REQUIREMENTS{,_EN}.md`、`docs/BUSINESS_WORKFLOWS{,_EN}.md`、`docs/RFC-proofrail-unattended-ai-engineering-product.md`、`docs/DOCUMENTATION_PLAN{,_EN}.md`、`docs/CODING_CONVENTIONS{,_EN}.md`、`docs/OPERATIONS{,_EN}.md`、`docs/INSTALLATION{,_EN}.md`、`docs/INSTALLATION_PLAN{,_EN}.md`：**不适用 + 理由（本片只触执行纪律层与判据本体，不触产品范围、业务流程、安装模型、编码规范、操作手册）**

14. `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md`：**不适用 + 理由（本片不是结项报告，不触模板）**

15. `docs/t027/CHECKLIST_v48{,_EN}.md`：**不适用 + 理由（v48 已冻结；v49 已反映 P59 状态；v50 由下一片生成时处理）**

16. 报告本体（`docs/validation/p60-directive-judge-blindspots{,_EN}.md`）：**⑧ 阶段必做**——按 §11.3 完整结构；`artifacts` 围栏按纪律 G（路径排序、TAB 分隔、CN/EN 逐字节相同）；`sixstage` 围栏按 G7-f（⑥ 必做）

**文档同步的机械核验**（⑨ 阶段）：
- `G3(a)` 六形状全等（CN/EN）
- `G1-b` 版本元数据（**仅 J2 升版**；J1 / J3 头部与 §0.3 均维持 v1.63 不变）
- `G6-6` 台账状态词闭集（OB-138 处置列翻词须在闭集内）
- `G6-7` 声明/台账对账（§0.3 声明的 OB 编号必须在 §12.4 存在）
- `G5(d)` 暂存集 ⊆ `artifacts` 围栏声明的集合
- `G5(f)` index 作用域下工作区干净
- **J2 路径额外**：`--selftest` 全绿且断言数增加

【五 ⑩ 回执要求】
1. 派单表（`agentName` 逐字机械读取；若 debug 日志无 transcript，须显式标注"未机械校验"并作为已知事实登记）
2. ① 三案对比表摘要 + 推荐（**仅供参考，不得写成结论**）
3. J1 / J2 / J3 裁定后的执行摘要（按裁定路径：改了哪节 / 加了什么判据 / 有无触碰历史行 / 有无测试）
4. ⑤ / ⑥ / ⑦ 结论摘要（含整改与重跑）
5. 三作用域 + SELFTEST 原始行（J2 路径下 SELFTEST 若增加，据实报）
6. 文档同步 `git diff --stat`（逐项覆盖【四】的 16 项；不适用项列理由）
7. 明确未执行事项：未改产品代码 / CI / 契约 / schema / fixtures；未替用户裁定；未 commit/push；未启动新切片
8. 待授权项：提交授权；若 J1 落地，ADR-035 是否随同提交
9. **本片自身越界登记**：若本片出现 §1.7.2 三条 bullet 的越界，如实登记并披露；未发生也要明确写"无"

【六 下次起跑前的提示】
- P60 不依赖 R-08 / R-09；AT23 不依赖 P60，只依赖 R-08 / R-09 裁定
- P60 闭环后，层 2 全部出队；此时可起 v50 快照（承载 P60 闭环 + G-4 出队 + §九 门禁坑第 23 条 + 归一脚本修复）
- 用户可在 P60 运行期间同步做 R-08 / R-09 裁定（R 表回写由 P60 ⑩ 或 v50 生成时顺手，不另起片）

现在开始执行 P60 步骤 0；阶段 A 出表后停等裁定，不得进入 ②。

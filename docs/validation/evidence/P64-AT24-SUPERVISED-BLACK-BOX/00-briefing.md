切片：P64 AT24-SUPERVISED-BLACK-BOX（层 1 后续候选第 12 项 · 验收类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 docs/t027/CHECKLIST_v53.md（跨会话快照）。
2. 核对 §一 的 8 个锚点：HEAD/origin = cef2ddb（v53 提交；本片起跑时以 git rev-parse 实况为准）；准则 v1.64；SELFTEST 722/722；三作用域 TOTAL_FAIL=0；v50 行数 232（按 LF）；v51 行数 270（按 LF）；v52 行数 278（按 LF）。除 HEAD 预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况。
3. **主控模型声明**：AT-24 涉及 R-09.1~R-09.9 九项 unknown 接受边界的逐项判定 + 三句话硬性要求 + 黑箱归还的停机/身份语义；与 P61/P63/F2-FIX 同类（验收/判定类收尾链）主控模型保持一致。**步骤 0 按 §1.7.4 独立留痕要求登记 ADR-040**（四要素：裁决人=用户；日期=起跑日；理由=R-09 九项边界判定 + S1 Exit 验收末环 + 同类先例一致；依据=本轮起跑消息；性质如实标注：v53 §七 第 8 条字面未触发，不得作为先例）。回执首行标 [自执行]。
4. 记忆唤醒：本片受 R-09 裁定（放宽版，九项边界 + 三句话硬性，docs/validation/r08-r09-ruling{,_EN}.md）约束；受 T028 实现（DEV_PLAN.md:119 + docs/validation/t028-supervised-black-box{,_EN}.md）约束；受 SEC-15 威胁模型约束；受 v53 §四 层 1 后续候选第 12 项约束；单次输出 ≤10,000 字符；不自行 commit/push；不替用户裁定；不启动任何新切片。

【〇 本片环境声明与流水线裁剪】
1. 环境：**不开启 Ubuntu VM、不开启网络代理、不启用 GPT-6.1 Sol 双打**。本片是 Windows 本机 console/chain 层验收，用 stub visible 投递 + 模拟人工归还 + 后置扫描，无 Linux、无外部调用。
2. 流水线：验收类切片，保留 ①（执行计划确认，轻量化）→ ③（验收腿执行）→ ②（若验收发现 T028 实现缺陷则修复；无缺陷则跳过）→ ④ 门禁 → ⑤ 预审 → ⑥ 独立扫描 → ⑦ 终审 + 盲审 → ⑧ 文档 → ⑨ → ⑧c → ⑩。
3. **验收腿的边界**（前置明确）：AT-24 的验收腿是**机制腿**（stub visible 投递 + 模拟人工归还 + 后置扫描 + 独立 gates/review + reduced/unknown 报告），**不是真实 SessionBridge 端到端**——真实 visible 端到端依赖 T027-B，超出本片范围（T028 完成记录明写：真实 CLI 端到端依赖 T027-B，尚不可运行发布）。真实端到端登记为 unknown/边界，按 R-09.3（网络 = 回环可接受公网阻断）与 R-09.7（外部副作用 = 用户确认）判。

【一 授权边界（角色与模型）】
- ① 执行计划确认：deep-reasoner + DeepSeek V4 Pro (deepseek)（轻量化）
- ③ 测试/验收腿执行：prfrail-tester + DeepSeek V4.1 Flash (deepseek)
- ② 实现（仅当 ③ 发现 T028 实现缺陷）：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek)（额度 1+1）
- ⑥ 独立扫描：independent-reviewer + MAI-Code-1.1-Flash (copilot)（额度 3，§7.2 ③ 不成立 ⇒ 必做）
- ⑦ 终审 + 盲审：independent-reviewer + GPT-5.3-Codex (copilot)（按 §7.5 高风险档 1+3；盲审另计 +1，§7.7 必抽：触及停机/身份语义）
- ⑧ 文档装配：prfrail-documenter + DeepSeek V4.1 Flash (deepseek)
- ⑨ / ⑧c / ⑩：主控（V4 Pro）执行
- 每次 runSubagent 必须同时传 agentName 与 model；严格串行

【二 批量辅助（强制）】
- 载体/模型：deep-reasoner + GLM-5.3-Flash (glm)
- 实派次数：≥2 次
- 任务限定：
  ① R-09.1~R-09.9 九项边界与实测证据的逐项对照枚举（枚举 → 逐项 diff，不判断、只列事实；每项输出「裁定值 / 实测证据路径 / 是否吻合」三列）
  ② 报告 CN/EN 双语镜像检查——`at24-supervised-black-box{,_EN}.md` 逐行同位比对
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三 核心任务：AT-24 验收 + R-09 九项边界判定】

## 3.1 验收腿清单（③ 必须逐项执行 + 落盘证据）

按 R-09 裁定 + CONTRACTS §7.1 + T028 实现，下列验收腿必须逐项执行并留证。**每项必须覆盖正例 + 反例**（如风险确认绑定：正例 = 正常绑定通过；反例 = 版本漂移阻断），不能只跑正例。

| # | 验收腿 | 依据 | 判据 |
|---|---|---|---|
| 1 | 风险确认绑定 | R-09.7 + CONTRACTS §7.1（supervised 段） | 启动前显示并持久化 `supervised-risk-confirmation`（绑定 run/task/step/attempt + workspace 摘要 + 风险提示版本/哈希 + visible 投递事实）；版本漂移阻断（正例 = 正常绑定通过；反例 = 版本漂移阻断） |
| 2 | 投递 ≠ 完成 | R-09 + CONTRACTS §7（supervised 段） | visible 投递 receipt 不产生 task/step PASS；投递成功只证明 UI 可见（反例 = 投递成功但 task/step 未被置 PASS） |
| 3 | 显式归还 + 后置全量扫描 | R-09.1 + R-09.2 + CONTRACTS §7（handoff complete 回执） | 显式归还后产生 `handoff-receipt{outcome=complete}`；后置 manifest/diff/范围/秘密/文件类型/副作用扫描全绿才通过；未知进程/外部副作用阻断（反例 = 扫描发现未知进程 ⇒ 阻断） |
| 4 | 独立 gates/review | R-09 + CONTRACTS §7（review/promotion） | 独立 build/test/verify gates + review 通过后才可接受；exit 0 / completed receipt 不直接产生 PASS（反例 = exit 0 但 gates 红 ⇒ 不 PASS） |
| 5 | unknownDimensions 七项全必填含 cost | R-09.4 + R-09.8 | `supervised-report{assuranceLevel=reduced}` 的 `unknownDimensions` 非空且含 `cost`；七项（工具使用/资源读取/网络/费用/后台进程/会话连续性/外部副作用）逐项必须出现（反例 = 缺任何一项 ⇒ 报告不合规） |
| 6 | 三句话硬性出现 | R-09 §5.3 | 报告 §1 与执行摘要**各出现一次**三句话（逐字；逐字核验由 ⑤⑥⑦ 各环执行，不是主控自证） |
| 7 | 禁止自动降级 | R-09 + CONTRACTS §7 | AgentRunner 不可用时不自动落入黑箱模式；不得因 AgentRunner 失败而自动切黑箱（反例 = AgentRunner 失败后自动切黑箱 ⇒ 阻断） |
| 8 | 未知维度不显示为安全 | R-09 + CONTRACTS §7 | 报告与渲染对每个未知维度只显示 `unknown`，不得显示 `safe`/`verified`/数字承诺（反例 = 渲染出现 `safe` ⇒ 不合规） |
| **9** | **R-09.1~R-09.9 九项边界逐项验证**（拆为 9 个子验收腿，见下表） | R-09 裁定原文 | **每项逐条验证，不合并、不笼统** |

**第 9 项展开为 9 个子验收腿**：

| 子项 | R-09 裁定值 | 验收判据 |
|---|---|---|
| 9.1 工具使用 unknown | **(c) 视情况** | 归还后产物扫描能证明「产物在 workspace 内、无外部副作用」则接受；发现 workspace 外痕迹则阻断 |
| 9.2 资源读取 unknown | **(c) 视情况** | 归还后扫描能证明「未越界读取」（manifest/diff 比对只读了 workspace 内文件）则接受；发现 workspace 外读取痕迹则阻断 |
| 9.3 网络 unknown（底线） | **(c) 视情况** | 回环可接受、公网必须阻断；起跑前一次性确认（步骤 0），不阻断起跑 |
| 9.4 费用 unknown | **(a) 显式标 unknown** | `unknownDimensions` 必含 `cost`（与验收腿 5 联动） |
| 9.5 后台进程 unknown | **(a) 可接受** | 用户负责清理；报告明写 |
| 9.6 会话连续性 unknown | **(a) 可接受** | 黑箱模式设计边界；报告明写 |
| 9.7 外部副作用 unknown（底线） | **(c) 视情况** | 用户显式确认「已接受未知外部副作用」则接受；没有确认则阻断。起跑前签署确认记录（步骤 0），不阻断起跑 |
| 9.8 强制字段 | **(a) 七项全必填** | `assuranceLevel=reduced` 且 `unknownDimensions` 非空且含 `cost`（与验收腿 5 联动） |
| 9.9 判定口径 | **(b) 每项独立判定，不累积** | 某项 unknown 可接受不整体放行；某项阻断不整体失败 |

## 3.2 验收的判定主体与证据

- **R-09.1~R-09.9 的逐项判定主体 = ⑦（独立终审）**——主控只整理事实（每项边界与实测的逐项对照），⑦ 独立判定；
- **R-09.1~R-09.9 逐项判定的机械证据必须落盘**——每项对应 `03-test/` 下的证据路径（判定为「可接受」须有「扫描能证」的证据；判定为「阻断」须有「发现越界痕迹」的证据），⑦ 判定时需可追溯，不能只有「判定为视情况」；
- **报告三句话**由 ⑧ 文档装配写入，**⑤⑥⑦ 逐环核验逐字出现**（不是主控自证）；
- **AT-24 验收结论**（条件闭合 / 部分通过 / 不通过）由 ⑦ 独立判定。

## 3.3 若验收发现 T028 实现缺陷

- **实现缺陷（代码与 CONTRACTS §7/R-09 不符）**：在 ② 阶段内修（prfrail-implementer）；**修完必须重跑受影响的验收腿**，且重跑纳入 ③ 的 `R` 计数；修复范围限于 `internal/console/**` 与 `internal/chain/**` 的 `supervised-black-box` 相关文件，不动其他模块；
- **契约缺陷（CONTRACTS §7 的 supervised 段本身与 R-09 冲突或不够精确）**：**停下升级用户裁决**——不自行修订契约（§1.7.2 升级边界）；
- **R-09 裁定本身的缺陷**：停下升级用户裁决（R-09 已裁，本片不得自行修订）。

## 3.4 报告 §4 已知边界必含项（⑧ 装配时必须写入）

`at24-supervised-black-box{,_EN}.md` 的 §4「已知边界」**必须含**：

- 「真实 SessionBridge visible 端到端 = unknown，依据 T028 完成记录（真实 CLI 端到端依赖 T027-B，尚不可运行发布）」；
- 本片验收腿是**机制腿**（stub visible 投递 + 模拟人工归还 + 后置扫描），**不是真实 SessionBridge visible 端到端**；
- AT-24 的验收不覆盖真实 CLI 端到端，S1 Exit 的 S-G7 依赖链须如实标注。

【四 文档同步更新（17 项，不得遗漏）】

1. `docs/validation/at24-supervised-black-box{,_EN}.md` —— **必改**（新建，核心交付物：§1 执行摘要含三句话、§2 九项边界逐项判定、§3 验收腿证据、§4 已知边界（含 3.4 必含项）、§5 明确未执行事项、`artifacts` 围栏、`sixstage` 转录）
2. 证据包 `docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/` —— **必改**（新建：00-briefing.md / 01-arch/（含 R-09.3/R-09.7 用户确认件）/ 03-test/ / 04-reviews/ / 06-gates/ / MANIFEST / SHA256SUMS / freedom-list.md / input-list.md）
3. `docs/t027/REMAINING_SLICES{,_EN}.md` —— **必改**（层 1 后续候选第 12 项 ✅ 收口 / 或按实际结论；完成记录行只写报告路径）
4. `docs/DEV_PLAN{,_EN}.md` —— **必改**（AT-24 的验收状态按实际结论更新；若 AT-24 状态记在 T028 段下面则改 T028 段对应行，若有独立行则改独立行——⑧ 装配时按实况定位）
5. `docs/TEST_STRATEGY{,_EN}.md` —— **必改**（AT-24 行的验收结论按实际更新，引用报告路径）
6. `docs/ADR_REGISTER{,_EN}.md` —— **必改**（新增 ADR-040：主控 V4 Pro + AT-24 验收结论；四要素 + 性质如实标注）
7. `docs/validation/s1-exit{,_EN}.md` —— **必改**（**s1-exit 的对应门按 `docs/validation/p57-closeout-assessment{,_EN}.md` §4 的八门清单定位 AT-24 对应的门**（可能是 S-G2 或其他），按实际更新；S-G7 依赖链同步；**不得声称 S1 Exit 整体通过**——其余门仍待）
8. `docs/DELIVERY_DIRECTIVE{,_EN}.md` —— **不适用**（R-09 在裁定文档；本片不改准则）
9. `docs/CONTRACTS{,_EN}.md` —— **不适用**（除非 ③ 发现 T028 与 CONTRACTS §7 supervised 段不符——此时停下升级，不自行修订契约）
10. `docs/ARCHITECTURE{,_EN}.md` —— **不适用**（不改模块边界/数据流/存储）
11. `docs/SECURITY{,_EN}.md` —— **不适用**（R-09 与 SEC-15 的关系已在裁定文档；本片不改权限模型/威胁矩阵）
12. `docs/S1_SUPPORT_MATRIX{,_EN}.md` —— **不适用**（无平台支持层级变更）
13. `docs/S1_RELEASE_NOTES{,_EN}.md` —— **不适用**（无 tag/artifact/功能项变更）
14. `docs/t027/CHECKLIST_v53{,_EN}.md` —— **不适用**（快照冻结，v54 生成时继承）
15. `docs/validation/r08-r09-ruling{,_EN}.md` —— **不适用**（R-09 已裁；本片不改判据）
16. `docs/validation/t028-supervised-black-box{,_EN}.md` —— **按需**（AT-24 验收通过后在报告头部加「AT-24 已验收（P64）」时点注记指针；T028 原文**不改**，按 §11.1 历史不改写；验收未通过则标「不适用 + 理由（AT-24 未通过，T028 报告不加指针）」）
17. `docs/OPERATIONS{,_EN}.md` / `docs/INSTALLATION{,_EN}.md` / `docs/PRODUCT_REQUIREMENTS{,_EN}.md` / `docs/BUSINESS_WORKFLOWS{,_EN}.md` / `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md` —— **不适用**（不触操作手册/安装模型/产品范围/业务流程/结项模板）

【五 步骤 0 待办】
1. 建目录 `tmp/P64-AT24-SUPERVISED-BLACK-BOX/` 与 `docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/`
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF），并入 `MANIFEST.md` 与 `SHA256SUMS.txt`，在 `artifacts` 围栏声明
3. 台账 `REMAINING_SLICES{,_EN}.md` 登记 P64（🔄 进行中，R/D=0）
4. 落盘 `freedom-list.md`、`input-list.md`
5. **ADR 登记**：ADR-040（按 §1.7.4 独立留痕：四要素 + 性质如实标注「字面未触发 v53 §七 第 8 条，不得作为先例」）
6. 前置核对：读 R-09 全文（docs/validation/r08-r09-ruling{,_EN}.md §3 与 §5）、CONTRACTS §7 supervised 段、T028 实现（DEV_PLAN.md:119 + docs/validation/t028-supervised-black-box{,_EN}.md）、TEST_STRATEGY AT-24 行、s1-exit 按 P57 报告 §4 八门清单定位 AT-24 对应门
7. **R-09.3 网络底线用户确认**（起跑前一次性确认，不阻断起跑）：用户同轮确认「回环可接受、公网必须阻断」，落盘 `docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/01-arch/r09-3-user-confirmation.md`（记录裁决人 / 日期 / 确认原文 / 依据 = 本轮起跑消息 + R-09.3 原文）；**未确认则阻断 AT-24 验收**
8. **R-09.7 外部副作用底线用户确认**（起跑前一次性确认，不阻断起跑）：用户同轮确认「已接受未知外部副作用」，落盘 `docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/01-arch/r09-7-user-confirmation.md`（记录裁决人 / 日期 / 确认原文 / 依据 = 本轮起跑消息 + R-09.7 原文）；**未确认则阻断 AT-24 验收**

【六 执行步骤与红线】
1. ① 执行计划确认（V4 Pro，轻量化）：确定 9 项验收腿（含第 9 项展开为 9 子项 = 共 17 个子验收腿）的具体执行方案 + 证据落盘路径 + 机械断言 + ⑦ 判定问题清单；结论行 `DESIGN: DONE`；落 design-r1.md
2. ③ 验收腿执行（Flash）：逐项跑 9 项验收腿（含 9 子项），**每项覆盖正例 + 反例 + 机械断言**（如风险确认绑定：正例 = 正常绑定通过；反例 = 版本漂移阻断）；不能只写「跑通了」；证据落 `03-test/`
3. ② 实现（Flash，仅当 ③ 发现 T028 实现缺陷）：修 `internal/console/**` 与 `internal/chain/**` 的 `supervised-black-box` 相关文件；修完**重跑受影响的验收腿**（重跑纳入 ③ 的 `R` 计数）；**不动其他模块**
4. ④ 门禁：`go build/vet/test ./...` + `gate.js --all` 三作用域 + `SELFTEST 722/722`
5. ⑤ 预审（V4 Pro）：核对 9 项验收腿（含 9 子项）的完整性、R-09 九项边界的符合性、三句话的逐字出现、修复（若有）的正确性
6. ⑥ 独立扫描（MAI）：Section A–E 五段 + `sixstage` 围栏；重点：R-09 九项边界是否真正落实、报告三句话是否逐字
7. ⑦ 终审 + 盲审（Codex）：
   - **⑦ 是 R-09.1~R-09.9 逐项判定 + AT-24 验收结论（条件闭合/部分/不通过）的唯一主体**；主控只转录
   - ⑦ 报告必须逐项：9 项验收腿（含 9 子项）的通过/不通过 + R-09.1~R-09.9 九项边界逐项判定 + **AT-24 验收结论**
   - ⑦ 首轮不附任何清单（锚定硬隔离）
   - **盲审 1 次**（§7.7 必抽：触及停机/身份语义），首轮不附任何清单
8. ⑧ 文档装配（Flash）：按【四】17 项逐项；**报告三句话必须在 §1 与执行摘要各出现一次（逐字）**；**§4 已知边界必须含 3.4 必含项**
9. ⑨ 三作用域 + SELFTEST 全绿
10. ⑧c 临时件清理（tmp/P64-AT24-SUPERVISED-BLACK-BOX/ 清空，保留证据包）
11. ⑩ 停点：停在 ⑩ 等同一轮授权；严禁自行 commit/push

【七 执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束（不得自行升版/伪造授权/超额派单）
3. 临时件一律 `tmp/P64-AT24-SUPERVISED-BLACK-BOX/`；只增不删；⑧c 统一清理
4. **⑦ 是 R-09.1~R-09.9 逐项判定 + AT-24 验收结论的唯一主体**——主控不得自判；⑦ 判定后主控只转录
5. **报告三句话必须逐字**（不是改写、不是摘要）：
   ①「本报告的保证级别是 `reduced`，不是 `full`。」
   ②「七项 unknown 维度的接受边界为：工具使用/资源读取 = 视情况（扫描能证无越界则接受）；后台进程/会话连续性 = 可接受；网络 = 回环可接受、公网阻断；费用 = 显式标 unknown；外部副作用 = 用户已显式确认接受。」
   ③「本报告不构成对黑箱 Agent 过程正确性的保证；ProofRail 只保证归还后的产物扫描、独立门禁、证据和评审。」
   —— ⑤⑥⑦ 逐环核验逐字出现，不是主控自证
6. **验收腿 = 机制腿**（stub visible 投递 + 模拟人工归还 + 后置扫描），不是真实 SessionBridge 端到端；真实端到端登记为 unknown/边界，按 R-09.3 与 R-09.7 判；**报告 §4 已知边界必须含 3.4 必含项**
7. **T028 实现缺陷在 ② 修（范围限于 supervised 相关文件，修复后重跑受影响验收腿、重跑纳入 R 计数）；契约缺陷/R-09 缺陷停下升级**——不自行修订契约或 R-09
8. **不改产品代码以外的东西**：除 ② 修 T028 缺陷外，不改 `internal/**` 其他模块、不改 CONTRACTS/schema/fixtures/workflows、不改准则、不改 R-09、不改 T028 报告原文（只加时点注记指针）
9. **步骤 0 的 R-09.3/R-09.7 两条用户确认必须先落盘**（01-arch/ 下两件），否则 ③ 验收腿执行缺前置

【八 ⑩ 回执要求】
1. 派单表（`agentName` 逐字机械读取；若 debug 日志无 transcript 须显式标注「未机械校验」）
2. ① 执行计划摘要（9 项验收腿含 9 子项的执行方案 + 证据路径 + 机械断言 + ⑦ 判定问题清单）
3. 步骤 0 的 R-09.3/R-09.7 用户确认落盘（01-arch/ 两件路径 + 哈希）
4. 9 项验收腿（含 9 子项）逐项结果（通过/不通过 + 证据路径 + 关键读数）
5. ⑤/⑥/⑦ 结论摘要（⑦ 报告必须逐项含 9 项验收腿 + R-09.1~R-09.9 九项边界 + AT-24 验收结论）；⑦ 计量（终审 N + 复审 M + 盲审 K）
6. **报告三句话的逐字核验**（⑤⑥⑦ 各环均须核验，不是主控自证）
7. 三作用域 + SELFTEST 原始行
8. 文档同步 `git diff --stat`（逐项覆盖【四】17 项；不适用项列理由）
9. 明确未执行事项：未改产品代码（除 ② 修 T028 缺陷）/ CI / 契约 / schema / fixtures / 判据文本；未替用户裁定；未 commit/push；未启动新切片
10. 本片自身越界登记：若出现 §1.7.2 三条 bullet 越界，如实登记并披露；未发生也要明确写「无」
11. **AT-24 最终验收结论**（条件闭合 / 部分通过 / 不通过）+ R-09 九项边界逐项判定 + 后续动作（S1 Exit 对应门可闭合 / 保持 blocked）
12. 待授权项：提交授权；是否登记「AT-24 通过后 S1 Exit 对应门可闭合」的后续动作
13. 给下版清单的事实包（v54 生成时用）

【九 下次起跑前的提示】
- P64 闭环后：
  - 若 **AT-24 条件闭合**：S1 Exit 的对应门（按 P57 报告 §4 八门清单定位）可闭合；S-G7 依赖链进一步收窄（仍须声明其余门待）；层 1 后续候选第 12 项 ✅ 出队；
  - 若 **部分通过/不通过**：AT-24 保持相应状态；F2 修复路径（若涉代码）或 R-09 修订（若涉裁定）由用户裁定；
- v54 生成时承载：AT-24 验收结论 + S1 Exit 状态更新 + 队列更新 + 既有观察项；
- G5(d) 第 26 条候选与 R-08.5 语义级长期观察项等用户裁定是否启动治理片；
- 层 2 已全部出队，层 3 维持不启动。

【十 读数校正注记（起跑前已知事实，不阻断起跑）】

关于锚点核验中 v52 全文件行数的读数差异，起跑前已由用户与网页版顾问独立核实，事实如下：

1. **实测读数**：`docs/t027/CHECKLIST_v52{,_EN}.md` 按 LF（0x0A 字节）计数，实况均为 **277**（CN/EN 对称）；双实现复核一致（PowerShell `ReadAllBytes` 逐字节计数 + Node `Buffer` 逐字节计数）；两文件末尾字节均为 `0x0A`；`git status --short -- docs/t027/CHECKLIST_v52.md docs/t027/CHECKLIST_v52_EN.md` 干净无改动。
2. **与 v53 §一 记录的差异**：v53 §一 记 v52=278，实况 277，差 1。定性 = **v53 §一 生成时的读数偏差**，与 v51 §一 记 v50=233（v52 起修正为 232）的处理口径**完全同构**（详见 v53 §一 的 v50 口径注记）。
3. **对 P64 的影响**：无。按 P64 起跑消息【〇·零】第 2 条自身规则——「除 HEAD 预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况」——8 个锚点中仅此 1 项不符（**不是** ≥2 项），**不触发停止**。其余 7 项（HEAD / origin / 准则 v1.64 / SELFTEST 722/722 / 三作用域 TOTAL_FAIL=0 / v50=232 / v51=270）全部符合预期。
4. **处置要求**（P64 步骤 0 内执行，不另立切片）：
   a) 在 `docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/freedom-list.md` 里登记一条「读数校正注记」，逐字写明上述 1–3 项事实，作为 v54 生成时的读数修正候选；
   b) v53 §一 的 278 读数**不追改**（历史快照冻结，与 v51 §一 记 v50=233 不追改口径一致）；
   c) ⑩ 回执的「给下版清单的事实包」（【八】第 13 项）里增一条：「v52 全文件行数 v54 起修正为 277（CN/EN 对称），v53 §一 记 278 系读数偏差」。
5. **性质声明**：本注记**不构成** P64 的「越界登记」（【八】第 10 项），也不改变 P64 的任何验收语义；它只是把起跑前已确认的读数事实前置告知主控，避免步骤 0 因 v52 读数差异产生误判。

现在开始执行 P64 步骤 0。

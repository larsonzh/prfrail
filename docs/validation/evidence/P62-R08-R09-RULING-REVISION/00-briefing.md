# 00-briefing — P62 起跑消息逐字落盘

> 落盘时点：2026-11-08 步骤 0。以下为用户起跑消息逐字原文（不做任何修改）。

切片：P62 R08-R09-RULING-REVISION（层 1 第 9 项 · 治理类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 docs/t027/CHECKLIST_v51.md（跨会话快照）。
2. 核对 §一 的 6 个锚点：HEAD/origin = cef120e（v51 生成时记录 05be20b；其后 v51 有两个提交：dc8b4ac 补齐门禁坑 1-19、cef120e 拆分门禁坑子节；以 git rev-parse 实况为准）；准则 v1.64；SELFTEST 722/722；三作用域 TOTAL_FAIL=0；v50 行数 233（CN/EN 对称）。除 HEAD 预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况。
3. **主控模型声明**：本片为判据本体修订（规则语义变更），按 v51 §七 第 8 条提升为 DeepSeek V4 Pro。**步骤 0 按 §1.7.4 独立留痕要求登记 ADR-037（四要素：裁决人=用户；日期=起跑日；理由=判据本体修订需独立视角；依据=本轮起跑消息；性质如实标注：不得作为先例）**。回执首行标 [自执行]。
4. 记忆唤醒：本片受 R-08/R-09 裁定（docs/validation/r08-r09-ruling{,_EN}.md）约束；受 P61 闭环（PARTIAL 收口，ADR-036）约束；单次输出 ≤10,000 字符；不自行 commit/push；不替用户裁定；不启动任何新切片。

【〇 本片环境声明与流水线裁剪】
1. 环境：**不开启 Ubuntu VM、不开启网络代理、不启用 GPT-6.1 Sol 双打**。本片是纯 .md 文档修订（治理类），无代码语义、无 syscall、无外部调用、无可测语义设计。
2. 流水线：治理类切片，保留 ①（判据设计，轻量化）→ ②（判据文本实现）→ ④（门禁）→ ⑤（预审）→ ⑥（独立扫描）→ ⑦（终审）→ ⑧（文档装配）→ ⑨（三作用域 + SELFTEST）→ ⑧c → ⑩。③（测试）裁剪——纯文档切片无实现对象；按 v1.64 §5.2「文档型切片」说明，②/③/⑨ 退化为文档实现 + 编码门禁 + G1–G9 自检。

【一 授权边界（角色与模型）】
- ① 判据设计：deep-reasoner + DeepSeek V4 Pro (deepseek)
- ② 判据文本实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek)（额度 1+1）
- ⑥ 独立扫描：independent-reviewer + MAI-Code-1.1-Flash (copilot)（额度 3）
- ⑦ 终审：independent-reviewer + GPT-5.3-Codex (copilot)（额度 1+3，治理片按 §7.5；实得可能用不完，⑩ 回执据实登记）
- ⑧ 文档：prfrail-documenter + DeepSeek V4.1 Flash (deepseek)
- ⑨ / ⑧c / ⑩：主控（V4 Pro）执行
- 每次 runSubagent 必须同时传 agentName 与 model；严格串行

【二 批量辅助（强制）】
- 载体/模型：deep-reasoner + GLM-5.3-Flash (glm)
- 实派次数：≥2 次
- 任务限定：
  ① 枚举修订点——从 r08-r09-ruling{,_EN}.md、P61 报告、v51 §二 §四 中摘录「R-08.5 判据原文」与「R-08.1 事实陈述」的逐字文本（不判断、只摘录）
  ② 双语镜像检查——修订后的 r08-r09-ruling{,_EN}.md 逐行同位比对
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三 核心任务（R-08.5 判据修订 + R-08.1 事实修正）】

## 3.1 核心交付一：R-08.5「继承漂移豁免」五要素

在 docs/validation/r08-r09-ruling{,_EN}.md 的 R-08.5 段（现 §2 R-08.5）增补一段「继承漂移豁免的适用条件」，明确下列五要素必须同时满足：

1. **同基线**：与 B4 基线（docs/validation/t027-at23-e2e.md + docs/validation/evidence/directive-history/b4/archive/）同一基线；
2. **同断言**：断言 id、断言含义、断言所在的 harness（tools/b4-e2e/）均不变；
3. **漂移方向更严格**：断言值的变化方向是「更保守/更严格」（如 process-tree-taken-away 从 true→false = 未检测到进程树被取走 = 判定更保守），不是「更宽松」（如 false→true）；
4. **四态不变**：verdictClass / chainState / taskState / stepState 四个关键字段与 B4 基线逐字一致；
5. **成因可追溯**：漂移的成因可追溯到 B4 之后的特定产品改动（如 P61 的 02d5797，A5 pid reports 绑定）。

同时明确：
- 五要素**同时满足** ⇒ 可豁免「前提①不满足」的判定，前提①视为满足（或视为「继承漂移豁免成立」）；
- 任一不满足 ⇒ 不豁免，按原判据处理；
- 豁免成立**不等于** AT-23 条件闭合——真实候选腿的前提②仍须满足；
- 豁免须由 ⑦ 独立判定（不是主控自行判定）。

## 3.2 核心交付二：R-08.1 事实错误修正（只改仓库内文件）

在 docs/validation/r08-r09-ruling{,_EN}.md 的 §8「变更代价分级」表「R-08.1 封装函数 | 中 | 已落地代码，要过审查链」行：

- **改为**：「R-08.1 封装函数 | 中 | **P61 ② 新建**（Invoke-CopilotCli / buildCopilotArgs），要过审查链」；
- **加一个时点注记**：「> **时点注记（P62，2026-11-08）**：本行原写「已落地代码」，实况是 P61 ② 才新建（commit 6bb23a1 之前仓库无此函数）。按 §11.1 历史行不改写，本行加注记。事实错误独立于设计决策（R-08.1 裁定本身不变）。」

**注意：交接文档 v3 是用户手上的仓库外文件，不在本片修订范围；如需更新由用户自行处理。本片只改仓库内的 r08-r09-ruling{,_EN}.md。**

## 3.3 配套任务：P61 报告加时点注记指针

在 docs/validation/at23-restricted-candidate{,_EN}.md 的「下一步建议」节里加一句（若该节不存在，则在报告末尾新加一节「后续指针」）；不得改写报告原有内容：

「> **时点注记（P62，2026-11-08）**：R-08.5 判据已于 P62 修订，增加「继承漂移豁免」五要素；P61 的前提①判定不变（当时 R-08.5 未修订），但后续 AT-23 复核可按 P62 修订后的判据评估。」

## 3.4 可选：s1-exit 同步

- 若 docs/validation/s1-exit{,_EN}.md 的 S-G1 状态描述含「AT-23 前提①不满足」，可在其后加「（R-08.5 判据已于 P62 修订，含继承漂移豁免五要素；AT-23 现状待复核）」；
- 若未涉及，标「不适用 + 理由」。

【四 文档同步更新（16 项，不得遗漏）】

1. docs/validation/r08-r09-ruling{,_EN}.md —— **必改**（核心交付：R-08.5 五要素 + R-08.1 事实修正）
2. docs/validation/evidence/R08-R09-RULING/ —— **必改**（MANIFEST/SHA256SUMS 重算，含修订后的 r08-r09-ruling{,_EN}.md 副本）
3. docs/t027/REMAINING_SLICES{,_EN}.md —— **必改**（P62 行、层 1 第 9 项标 ✅）
4. docs/ADR_REGISTER{,_EN}.md —— **必改**（登记 ADR-037）
5. docs/DEV_PLAN{,_EN}.md —— **不适用**（治理片不回写，按 §11.2）
6. docs/DELIVERY_DIRECTIVE{,_EN}.md —— **不适用**（R-08.5 不在准则正文，在 R-08/R-09 裁定文档）
7. docs/CONTRACTS{,_EN}.md —— **不适用**（不改 wire/schema/状态机）
8. docs/ARCHITECTURE{,_EN}.md —— **不适用**（不改模块边界/数据流）
9. docs/SECURITY{,_EN}.md —— **不适用**（不改权限模型/威胁矩阵）
10. docs/TEST_STRATEGY{,_EN}.md —— **不适用**（AT-23 判据在 R-08/R-09 里，不在测试策略）
11. docs/validation/s1-exit{,_EN}.md —— **按需**（若 S-G1 状态描述需补充 R-08.5 修订事实则改；否则「不适用 + 理由」）
12. docs/t027/CHECKLIST_v51{,_EN}.md —— **不适用**（快照冻结，v52 生成时继承）
13. docs/validation/at23-restricted-candidate{,_EN}.md —— **必改**（加时点注记指针，按 §3.3）
14. docs/S1_SUPPORT_MATRIX{,_EN}.md —— **不适用**（无平台/支持层级变更）
15. docs/S1_RELEASE_NOTES{,_EN}.md —— **不适用**（无 tag/artifact/功能项变更）
16. docs/OPERATIONS{,_EN}.md / docs/INSTALLATION{,_EN}.md / docs/PRODUCT_REQUIREMENTS{,_EN}.md / docs/BUSINESS_WORKFLOWS{,_EN}.md / docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md —— **不适用**（不触操作手册、安装模型、产品范围、业务流程、结项模板）

【五 步骤 0 待办】
1. 建目录 tmp/P62-R08-R09-RULING-REVISION/ 与 docs/validation/evidence/P62-R08-R09-RULING-REVISION/
2. 起跑消息逐字落盘为 00-briefing.md（UTF-8 with BOM + LF），并入 MANIFEST.md 与 SHA256SUMS.txt，在 artifacts 围栏声明
3. 台账 REMAINING_SLICES{,_EN}.md 登记 P62（🔄 进行中，R/D=0）
4. 落盘 freedom-list.md、input-list.md
5. **ADR 登记**：ADR-037（按 §1.7.4 独立留痕：四要素 + 性质如实标注「不得作为先例」）
6. 前置核对：读 r08-r09-ruling{,_EN}.md 全文、P61 报告的 §5 §6 §7、v51 §二 §四 的 P61 事实包

【六 执行步骤与红线】
1. ① 判据设计：V4 Pro 产出「五要素的精确文本」+「R-08.1 事实修正的精确措辞」+「豁免成立的判定主体（⑦）」三段；轻量化，产出 design-r1.md
2. ② 判据文本实现：Flash 按 design-r1 精确改写 r08-r09-ruling{,_EN}.md；逐字对照，不自由发挥
3. ④ 门禁：go build/vet/test ./... + gate.js --all 三作用域 + SELFTEST 722/722
4. ⑤ 预审：V4 Pro（额度 1+1）
5. ⑥ 独立扫描：MAI（额度 3），重点 Section A–E 五段 + sixstage 围栏
6. ⑦ 终审：Codex（额度 1+3，治理片按 §7.5），重点：五要素的措辞是否精确无歧义、豁免成立的判定主体是否明确（⑦ 而非主控）、R-08.1 事实修正是否与原裁定一致、双语是否对称
7. ⑧ 文档装配：按【四】16 项逐项
8. ⑨ 三作用域 + SELFTEST 全绿
9. ⑧c 临时件清理（tmp/P62-R08-R09-RULING-REVISION/ 清空，保留证据包）
10. ⑩ 停点：停在 ⑩ 等同一轮授权；严禁自行 commit/push

【七 执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束（不得自行升版/伪造授权/超额派单）
3. 临时件一律 tmp/P62-R08-R09-RULING-REVISION/；只增不删；⑧c 统一清理
4. **豁免成立的判定主体必须是 ⑦（独立终审），不是主控**——本片在 R-08.5 判据文本中必须明写此点

【八 ⑩ 回执要求】
1. 派单表（agentName 逐字机械读取；若无 transcript 须显式标注）
2. ① 判据设计摘要（五要素精确文本 + R-08.1 事实修正措辞 + 判定主体）
3. ⑤/⑥/⑦ 结论摘要（含整改与重跑；⑦ 的 1+3 计数）
4. 三作用域 + SELFTEST 原始行
5. 文档同步 git diff --stat（逐项覆盖【四】16 项；不适用项列理由）
6. 明确未执行事项：未改产品代码/CI/契约/schema/fixtures；未替用户裁定；未 commit/push；未启动新切片
7. 本片自身越界登记：若出现 §1.7.2 三条 bullet 越界，如实登记并披露；未发生也要明确写「无」
8. 待授权项：提交授权；是否登记「AT-23 按 P62 修订后重新评估」的后续动作
9. 给下版清单的事实包（v52 生成时用）

【九 下次起跑前的提示】
- P62 闭环后，层 1 第 9 项出队；层 1 第 8 项（AT-23）的现状需按 P62 修订后的 R-08.5 判据重新评估——由用户裁定是否启动复核
- 层 1 第 10 项（F2 修复片）仍后置——仅当 R-08.5 修订后 F2 仍被判 HIGH 时启动
- 层 2 已全部出队，层 3 维持不启动
- v52 生成时承载：P62 闭环 + R-08.5 修订 + R-08.1 事实修正 + 「继承漂移豁免」五要素 + P61 报告时点注记

现在开始执行 P62 步骤 0。

切片：F2-FIX-STUB-HANG（层 1 第 10 项 · 修复类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 docs/t027/CHECKLIST_v52.md（跨会话快照）。
2. 核对 §一 的 7 个锚点：HEAD/origin = de37c58（P63 提交；本片起跑时以 git rev-parse 实况为准）；准则 v1.64；SELFTEST 722/722；三作用域 TOTAL_FAIL=0；v50 行数 232（按 LF 计数）；v51 行数 270（按 LF 计数）。除 HEAD 预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况。
3. **主控模型声明**：本片代码改动为 S 级（tools/agent-stub 的 hang 分支补写空子报告），但本片产出**直接决定 AT-23 能否从 PARTIAL 转为条件闭合**（前提①在 F2 断言清零后满足），涉及判据语义与重大状态变化；与 P61/P62/P63 同属 AT-23 收尾链，主控模型保持一致。**v52 §七 第 8 条字面未触发，但作为 AT-23 收尾链末环的判定任务登记为 V4 Pro**。**步骤 0 按 §1.7.4 独立留痕要求登记 ADR-039**（四要素：裁决人=用户；日期=起跑日；理由=AT-23 收尾链末环判定 + 重大状态变化 + 与前序三片一致；依据=本轮起跑消息；性质如实标注：字面未触发 §七 第 8 条，不得作为先例）。回执首行标 [自执行]。
4. 记忆唤醒：本片受 P63 判定（`F2 继承漂移豁免不成立`，`de37c58`）约束；受 P62 修订后的 R-08.5（五要素）约束；受 v52 §四 第 10 项「仅当第 11 项判定 F2 豁免不成立时启动」约束（已达成）；单次输出 ≤10,000 字符；不自行 commit/push；不替用户裁定；不启动任何新切片。

【〇 本片环境声明与流水线裁剪】
1. 环境：**不开启 Ubuntu VM、不开启网络代理、不启用 GPT-6.1 Sol 双打**。本片是 Windows 本机代码小修 + layer3 F2 harness 重跑 + 判定，无 Linux、无外部调用、无不可测语义设计。
2. 流水线：修复类切片，保留 ①（架构，轻量化）→ ③（测试先行，TDD 红测）→ ②（实现）→ ④（门禁）→ ⑤（预审）→ ⑥（独立扫描）→ ⑦（终审 + 盲审，**⑦ 是「F2 修复后 AT-23 前提①是否满足」的判定主体**）→ ⑧（文档装配）→ ⑨（三作用域 + SELFTEST）→ ⑧c → ⑩。

【一 授权边界（角色与模型）】
- ① 架构：deep-reasoner + DeepSeek V4 Pro (deepseek)
- ③ 测试：prfrail-tester + DeepSeek V4.1 Flash (deepseek)
- ② 实现：prfrail-implementer + DeepSeek V4.1 Flash (deepseek)
- ⑤ 预审：deep-reasoner + DeepSeek V4 Pro (deepseek)（额度 1+1）
- ⑥ 独立扫描：independent-reviewer + MAI-Code-1.1-Flash (copilot)（额度 3）
- ⑦ 终审 + 盲审：independent-reviewer + GPT-5.3-Codex (copilot)（按 §7.5 治理片/高风险 1+3；盲审另计 +1，§7.7 必抽：改测试工具触及停机语义）
- ⑧ 文档装配：prfrail-documenter + DeepSeek V4.1 Flash (deepseek)
- ⑨ / ⑧c / ⑩：主控（V4 Pro）执行
- 每次 runSubagent 必须同时传 agentName 与 model；严格串行

【二 批量辅助（强制）】
- 载体/模型：deep-reasoner + GLM-5.3-Flash (glm)
- 实派次数：≥2 次
- 任务限定：
  ① F2 三腿 verdict-run.json 与 B4 基线逐字段对照（枚举 → 逐字段 diff，不判断、只列事实）
  ② 双语镜像检查——若 at23 报告 §11 与台账改动涉及 CN/EN，逐行同位比对
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三 核心任务：修复 tools/agent-stub 的 hang 分支 + 重跑 layer3 F2 + 判定前提①】

## 3.1 契约依据（必须前置明确，这是修复的合法性来源）

按 CONTRACTS §7 的「pid 报告写入口径（恒定写）」原文：
> 在执行到写入点时，无子进程写**空文件（0 字节）**，有子进程写子 pid 列表；若在写入前遭遇外部中断或写入失败，报告可缺失且必须按证据缺口 fail-closed 处理。

**现状违约**：`tools/agent-stub/main.go` 的 `hang` 分支（`case "hang"` 或 `mode == "hang"`）**跳过写入** `stub-child.pid`，导致：
- 第六枚冻结事实 `ProcessPidReportsHash` 缺失
- postflight 按 fail-closed 判缺口
- F2 正例的 `process-tree-taken-away` 断言由 B4 基线 `true` 变为 `false`
- `failedChecks` 非空

**修复方向**：`hang` 分支进入 hang 之前，**写一份 0 字节的 `stub-child.pid`**（无子进程写空文件 = 恒定写的合法形态），让第六枚事实齐备。这是让实现符合契约，不是 hack。

## 3.2 修复方案（① 架构确定具体落点）

① 必须明确：
- `hang` 分支在 `tools/agent-stub/main.go` 的精确行号（从 P61 报告与 P63 证据看大约在 L100–130 区间，但需实测确认）
- `hang` 分支当前是否写 `stub.pid`（自身 pid 报告）；若未写，需同步补齐（同理：无子进程写空文件；自身 pid 报告不是子报告，按 CONTRACTS §7 口径应正常写入）
- 与 `run` 模式的写入逻辑是否复用（若已存在 `writeEmptyChildPIDReport` 函数，直接调用）
- 修复后 `main_test.go` 是否需要新增回归测试（正例：hang 模式产生 0 字节 `stub-child.pid`；反例：破坏写入逻辑后断言变红）

## 3.3 重跑与判定（⑦ 是判定主体）

**③ 测试先行（TDD 红测）**：
- 修复前：重跑 `tools/b4-e2e/run-faces.ps1 -OnlyFace f2`，记录 F2 三腿 verdict-run.json；确认 `process-tree-taken-away` = `false`、`failedChecks` 非空（红基线）
- 修复后：重跑同一命令，确认 `process-tree-taken-away` = `true`、`failedChecks` = `[]`、三腿 verdict 类与 B4 基线一致（绿）

**⑦ 判定（本片核心判定，主控不自判）**：
- ⑦ 必须逐项核对：
  1. F2 三腿 verdict 类 / chainState / taskState / stepState / observations 与 B4 基线是否**逐字段一致**？
  2. `process-tree-taken-away` 是否回到 `true`？
  3. `failedChecks` 是否为空集？
  4. 是否**任何**其他断言（six-facts-distinct 等）受到本次修复影响？
- ⑦ 判定结论：**AT-23 前提①「机制腿全绿」是否满足**？
- **不需要重新判定 R-08.5 五要素**——F2 修复后，F2 不再有继承漂移，豁免问题自然消失。⑦ 的判定是「F2 修复是否成功 + AT-23 前提①是否满足」，不是「豁免是否成立」。

## 3.4 判定后果（⑦ 判定后主控执行）

- **前提①满足（F2 断言清零 + 三腿与 B4 基线一致）** ⇒ AT-23 从 PARTIAL 转为**条件闭合达成**（前提②P61 已满足）；F2 修复片（层 1 第 10 项）✅ 出队；层 1 第 8 项（AT-23）从 🔄 转 ✅ 出队；S1 Exit 的 S-G1 可闭合、S-G7 依赖链更新（仍须声明 S1 Exit 整体未通过——其余门 S-G2~S-G6、S-G8 仍待）
- **前提①仍不满足（F2 仍有 failedChecks / 三腿与基线仍不一致）** ⇒ F2 修复失败；AT-23 保持 PARTIAL（P61 判定不变）；F2 修复片登记「未达成」原因；后续由用户裁定下一步（可能改 R-08.5 判据、或另立其他修复路径）

【四 文档同步更新（15 项，不得遗漏）】

1. `docs/validation/at23-restricted-candidate{,_EN}.md` —— **必改**：① 末尾新增 §11「F2 修复片复核结论」节（⑦ 判定原文逐字 + F2 三腿与 B4 基线逐字段对照 + AT-23 前提①满足/不满足 + AT-23 状态变化）；② §1 顶部加「F2 修复片时点注记」（指向 §11；不改 §1 原文）；③ §1/§4/§10 原文**不改**（P61 时点与 P63 时点的历史事实，按 §11.1 冻结）
2. 证据包 `docs/validation/evidence/F2-FIX-STUB-HANG/` —— **必改**（新建：00-briefing.md / 01-arch/design-r1.md / 02-impl/ / 03-test/（红基线 + 绿结果 + 与 B4 对照）/ 04-reviews/（⑤⑥⑦ 逐字）/ MANIFEST / SHA256SUMS / freedom-list.md / input-list.md）
3. `docs/t027/REMAINING_SLICES{,_EN}.md` —— **必改**：F2 修复片行 ✅ 收口（若前提①满足）/ 或登记「未达成」；层 1 第 10 项 ✅ 出队（若满足）/ 或保持启动条件已达成的说明；层 1 第 8 项（AT-23）按 ⑦ 判定更新为 ✅ 条件闭合（满足时）/ 或保持 🔄（不满足时）
4. `docs/ADR_REGISTER{,_EN}.md` —— **必改**：登记 ADR-039（主控 V4 Pro + F2 修复方案 + ⑦ 判定结论；四要素 + 性质如实标注「字面未触发 v52 §七 第 8 条，不得作为先例」）
5. `docs/DEV_PLAN{,_EN}.md` —— **按需**：若 AT-23 状态变化（PARTIAL → 条件闭合），T027 段的 AT-23 状态须更新；若不变，标「不适用 + 理由（状态未变）」
6. `docs/TEST_STRATEGY{,_EN}.md` —— **按需**：AT-23 行的验收结论按 ⑦ 判定更新（与 DEV_PLAN 同）；若不变，标「不适用 + 理由」
7. `docs/validation/s1-exit{,_EN}.md` —— **必改**：S-G1 状态按判定结果更新（满足 → 可闭合；不满足 → 保持 PARTIAL）；S-G7 依赖链同步；**不得声称 S1 Exit 整体通过**（其余门仍待）
8. `docs/DELIVERY_DIRECTIVE{,_EN}.md` —— **不适用**（不改准则语义；R-08.5 判据不重判，P63 已定论）
9. `docs/CONTRACTS{,_EN}.md` —— **不适用**（不改 wire/schema；agent-stub 是测试工具，CONTRACTS §7 已明确「恒定写」语义，本片是让实现符合契约，不改契约本身）
10. `docs/ARCHITECTURE{,_EN}.md` —— **不适用**（不改模块边界/数据流/存储）
11. `docs/SECURITY{,_EN}.md` —— **不适用**（不改权限模型/威胁矩阵/秘密处理）
12. `docs/S1_SUPPORT_MATRIX{,_EN}.md` —— **不适用**（无平台支持层级变更）
13. `docs/S1_RELEASE_NOTES{,_EN}.md` —— **不适用**（无 tag/artifact/功能项变更）
14. `docs/t027/CHECKLIST_v52{,_EN}.md` —— **不适用**（快照冻结，v53 生成时继承）
15. `docs/OPERATIONS{,_EN}.md` / `docs/INSTALLATION{,_EN}.md` / `docs/PRODUCT_REQUIREMENTS{,_EN}.md` / `docs/BUSINESS_WORKFLOWS{,_EN}.md` / `docs/SLICE-CLOSEOUT-TEMPLATE{,_EN}.md` —— **不适用**（不触操作手册/安装模型/产品范围/业务流程/结项模板）

【五 步骤 0 待办】
1. 建目录 `tmp/F2-FIX-STUB-HANG/` 与 `docs/validation/evidence/F2-FIX-STUB-HANG/`
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF），并入 `MANIFEST.md` 与 `SHA256SUMS.txt`，在 `artifacts` 围栏声明
3. 台账 `REMAINING_SLICES{,_EN}.md` 登记 F2 修复片（🔄 进行中，R/D=0）
4. 落盘 `freedom-list.md`、`input-list.md`
5. **ADR 登记**：ADR-039（按 §1.7.4 独立留痕：四要素 + 性质如实标注「字面未触发 v52 §七 第 8 条，不得作为先例」）
6. 前置核对：读 P63 报告 §10、P61 报告 §2.4（F2 事实）、CONTRACTS §7 的「pid 报告写入口径（恒定写）」原文、`tools/agent-stub/main.go` 的 hang 分支与 run 分支、`tools/b4-e2e/main.go` 的 `process-tree-taken-away` 断言与 `six-facts-distinct` 断言

【六 执行步骤与红线】
1. ① 架构（V4 Pro，轻量化）：确定 `hang` 分支精确落点、修复方案（写空 stub-child.pid）、回归测试方案（正例 + 变异检验）、重跑方案（`run-faces.ps1 -OnlyFace f2`）、⑦ 判定问题清单；结论行 `DESIGN: DONE`；落 design-r1.md
2. ③ 测试先行（Flash）：红测——重跑 layer3 F2，确认 `process-tree-taken-away` = `false`、`failedChecks` 非空；写回归测试（正例：修复后 `stub-child.pid` 为 0 字节；变异：破坏写入逻辑后断言变红）
3. ② 实现（Flash）：改 `tools/agent-stub/main.go` 的 `hang` 分支；改 `tools/agent-stub/main_test.go` 补回归测试
4. ④ 门禁：`go build/vet/test ./...` + `gate.js --all` 三作用域 + `SELFTEST 722/722`
5. ⑤ 预审（V4 Pro）：核对修复方案与契约的符合性、回归测试的可证伪性、重跑命令的准确性
6. ⑥ 独立扫描（MAI）：Section A–E 五段 + `sixstage` 围栏
7. ⑦ 终审 + 盲审（Codex）：
   - **⑦ 是「F2 修复后 AT-23 前提①是否满足」的唯一判定主体**；主控只转录，不参与判定
   - ⑦ 报告必须逐项：F2 三腿 verdict 与 B4 基线逐字段对照 + `process-tree-taken-away` 值 + `failedChecks` 是否为空 + 修复对 `six-facts-distinct` 等断言的影响 + **最终判定：前提①满足 / 不满足**
   - ⑦ 首轮不附任何清单（锚定硬隔离）
   - **盲审 1 次**（§7.7 必抽：改测试工具触及停机语义），首轮不附任何清单
   - **不需要重判 R-08.5 五要素**（F2 修复后不再有继承漂移，豁免问题自然消失）
8. ⑧ 文档装配（Flash）：按【四】15 项逐项
9. ⑨ 三作用域 + SELFTEST 全绿
10. ⑧c 临时件清理（tmp/F2-FIX-STUB-HANG/ 清空，保留证据包）
11. ⑩ 停点：停在 ⑩ 等同一轮授权；严禁自行 commit/push

【七 执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束（不得自行升版/伪造授权/超额派单）
3. 临时件一律 `tmp/F2-FIX-STUB-HANG/`；只增不删；⑧c 统一清理
4. **⑦ 是 AT-23 前提①判定的唯一主体**——主控不得自判；⑦ 判定后主控只转录
5. **P61 报告原文不改**（§1/§4/§10 是 P61/P63 时点历史事实）——只在末尾新加 §11 + §1 顶部加时点注记
6. **修复必须符合 CONTRACTS §7「恒定写」语义**——无子进程写空文件，不是 hack；⑥⑦ 会检查修复是否只是「绕过断言」而非「让实现符合契约」
7. **重跑 layer3 F2 必须用 `tools/b4-e2e/run-faces.ps1 -OnlyFace f2`**（不改 harness、不改 main.go）；修复只改 `tools/agent-stub/main.go`

【八 ⑩ 回执要求】
1. 派单表（`agentName` 逐字机械读取；若 debug 日志无 transcript 须显式标注「未机械校验」）
2. ① 判定准备摘要（修复方案 + 回归测试方案 + 重跑方案 + ⑦ 判定问题清单）
3. ③ 红测基线 + 修复后绿结果（原始 verdict-run.json 路径 + 关键字段读数）
4. ⑤/⑥/⑦ 结论摘要（⑦ 报告必须逐项含 F2 三腿逐字段对照 + 前提①判定）；⑦ 计量（终审 N + 复审 M + 盲审 K）
5. 三作用域 + SELFTEST 原始行
6. 文档同步 `git diff --stat`（逐项覆盖【四】15 项；不适用项列理由）
7. 明确未执行事项：未改产品代码（除 `tools/agent-stub`）/ CI / 契约 / schema / fixtures / 判据文本；未替用户裁定；未 commit/push；未启动新切片
8. 本片自身越界登记：若出现 §1.7.2 三条 bullet 越界，如实登记并披露；未发生也要明确写「无」
9. **F2 修复最终结论**（修复成功 / 失败）+ **AT-23 前提①判定**（满足 / 不满足）+ 后续动作（AT-23 条件闭合达成 / 保持 PARTIAL）
10. 待授权项：提交授权；是否登记「若前提①满足，AT-23 条件闭合达成并更新 S1 Exit」的后续动作
11. 给下版清单的事实包（v53 生成时用）

【九 下次起跑前的提示】
- F2 修复片闭环后：
  - 若 **前提①满足**：AT-23 从 PARTIAL 转为**条件闭合达成**；层 1 第 8 项（AT-23）✅ 出队；层 1 第 10 项（F2 修复）✅ 出队；**层 1 全部出队**；S1 Exit 的 S-G1 可闭合、S-G7 依赖链更新（仍须声明 S1 Exit 整体未通过——其余门 S-G2~S-G6、S-G8 仍待）；v53 生成时承载「AT-23 条件闭合达成」；
  - 若 **前提①仍不满足**：AT-23 保持 PARTIAL；F2 修复片登记「未达成」；后续由用户裁定（可能改 R-08.5 判据、或另立其他修复路径）；
- 层 2 已全部出队，层 3 维持不启动；
- v53 生成时承载：F2 修复结论 + AT-23 状态（条件闭合 / 保持 PARTIAL）+ 队列更新 + 两项 P63 观察项（R-08.5 文件级 vs 语义级 / G5(d) 围栏声明坑）+ 「治理片独立报告统一性」等既有观察项。

现在开始执行 F2 修复片步骤 0。

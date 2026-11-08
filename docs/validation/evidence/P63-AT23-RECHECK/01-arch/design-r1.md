# design-r1 — P63 AT23-RECHECK ① 判定准备（F2 继承漂移豁免判定输入整理）

## A. 判定对象与职责边界

- 判定主体 = ⑦（独立终审）；主控与 ① 只整理事实；本文档**不含任何「豁免成立/不成立」结论**（红线）。
- 判定依据 = R-08.5 五要素（P62 修订版，`docs/validation/r08-r09-ruling.md` L42–49，逐字引用）：
  - 「**继承漂移豁免的适用条件（五要素，须同时满足）**：」
  - 「① **同基线**：对照基线 = B4 基线存档（`docs/validation/t027-at23-e2e.md` + `docs/validation/evidence/directive-history/b4/archive/`）；只允许以该存档为对照，不得另选或混用其它基线。」
  - 「② **同断言**：断言 id、断言含义、断言所在 harness（`tools/b4-e2e/`）与 B4 基线均不变；harness 有改动或断言被重定义即不满足。」
  - 「③ **漂移方向更严格**：断言值变化方向更保守/更严格（如 `process-tree-taken-away` 从 `true` 变 `false` = 未检测到进程树被取走 = 判定更保守）；反向变化（如 `false` 变 `true` = 更宽松）不满足。」
  - 「④ **四态不变**：`verdictClass` / `chainState` / `taskState` / `stepState` 与 B4 基线逐字一致；任一字段漂移即不满足。」
  - 「⑤ **成因可追溯**：漂移成因可追溯到 B4 之后的特定产品改动（如 P61 的 `02d5797`，A5 pid reports 绑定为第六枚冻结事实）；无法归因到特定改动即不满足。」
  - 「**豁免后果**：五要素同时满足 ⇒ 豁免「前提①不满足」的判定，前提①视为满足（继承漂移豁免成立）；任一不满足 ⇒ 不豁免，按本小节既有判据处理。」
  - 「**豁免边界**：豁免成立不等于 AT-23 条件闭合——真实候选腿的前提②仍须满足；豁免须由 **⑦（独立终审）** 判定，主控不得自行判定。」
- 判定对象 = 层 3 机制腿 F2 正例断言 `process-tree-taken-away` 的 true→false 漂移（B4 基线 → P61 重跑）。

## B. F2 漂移 ↔ 五要素逐项对照表

| 要素 | 要素原文（压缩摘引；逐字原文见 §A） | 已核对事实（只列事实，含路径与行号） | 待 ⑦ 判定的问题 |
|---|---|---|---|
| ① | 同基线：对照基线 = B4 基线存档（t027-at23-e2e.md + directive-history/b4/archive/）；只允许以该存档为对照，不得另选或混用其它基线。 | 对照文件均为 B4 基线存档：`docs/validation/t027-at23-e2e.md` L77（F2 正向 = 负载 `-mode hang`、看门狗上界 5s、映射 uncertain、无冻结事实、错误证据 5 枚）与 L60（两腿均断言所报 pid 已消失）；`directive-history/b4/archive/pack-verdict.json` L2（generatedAt 2026-09-22T16:57:51.437Z）、L1418（failedChecks: []）；`archive/evidence/f2-pos/verdict-run.json` 全文件。 | 本要素成立/不成立？P61 层 3 重跑的对照物是否仅取自上述 B4 存档、未混用其它基线？ |
| ② | 同断言：断言 id、断言含义、断言所在 harness（tools/b4-e2e/）与 B4 基线均不变；harness 有改动或断言被重定义即不满足。 | 双方面事实并列：**方面甲**——`tools/b4-e2e/main.go` L766 的 F2 检查表达式在 B4 基线提交（e8cbff6）与当前 HEAD 逐字节相同：`report.addCheck("process-tree-taken-away", gone && goneErr == nil, fmt.Sprintf("reported pids gone: %t (%v)", gone, goneErr))`；断言 id 在 B4 与 P61 两份 f2-pos verdict 的 checks 中均存在（各 L136）。**方面乙**——`tools/b4-e2e/main.go` 自 B4 以来被 `02d5797` 修改（L690/L697 `six-facts-distinct`、L1093 `ProcessPidReportsHash`，均属 F1 面与 facts 适配器，不在 F2 检查）；`git log --oneline e8cbff6..HEAD -- tools/b4-e2e/main.go` 仅 `02d5797` 一个提交。 | 本要素成立/不成立？条文中「harness 有改动」的适用范围（整个文件 vs 该断言所涉代码路径）如何界定？方面甲/乙并存时应按何种口径判定？ |
| ③ | 漂移方向更严格：断言值变化方向更保守/更严格（如 process-tree-taken-away 从 true 变 false = 未检测到进程树被取走 = 判定更保守）；反向变化（如 false 变 true = 更宽松）不满足。 | 漂移方向 = true→false。now detail = `reported pids gone: false (AgentRunner evidence invalid: missing required pid report(s) stub-child.pid so the managed processes cannot be proven gone)`；baseline detail = `reported pids gone: true (<nil>)`（`l3-baseline-compare.txt` L5–7）。false 的语义 = 所报 pid 未获证明已消失（fail-closed）。R-08.5 原文例示与本片漂移为同名断言同方向。 | 本要素成立/不成立？true→false 是否属于条文中「更保守/更严格」方向（未检测到进程树被取走）？ |
| ④ | 四态不变：verdictClass / chainState / taskState / stepState 与 B4 基线逐字一致；任一字段漂移即不满足。 | 两文件四字段逐字相同：`verdictClass="uncertain-or-paused"`、`chainState="PAUSED"`、`taskState="b4-task\u00001=FAILED"`、`stepState="b4-task\u0000b4-code-step\u00001=FAILED"`（B4 L11–14；P61 同位置同值）。`l3-baseline-compare.txt` L22：`summary: field drifts=0`。相邻非四态字段观察（只列事实）：B4 mapping `treeStatus=stopped`、`errorEvidenceCount=5`；P61 `treeStatus=unknown`、`errorEvidenceCount=6`。 | 本要素成立/不成立？四字段逐字一致是否满足条文？相邻字段（treeStatus/errorEvidenceCount）差异是否在本要素覆盖范围内？ |
| ⑤ | 成因可追溯：漂移成因可追溯到 B4 之后的特定产品改动（如 P61 的 02d5797，A5 pid reports 绑定为第六枚冻结事实）；无法归因到特定改动即不满足。 | 因果链：`9ba621d`（2026-10-21）为 `tools/agent-stub` 增 `writeEmptyChildPIDReport()`（仅非 hang 公共路径；`hang` 分支仍只写 `stub.pid` 后睡眠）→ `02d5797`（2026-10-22，`feat: bind pid reports into A5 frozen facts and fix reader three-state semantics`）在 `internal/adapters/agent_runner_postflight.go` 新增 Leg 5（L370–386，双 pid 报告必须存在）与 Leg 7（L392–411，重读双报告重算第六枚事实，`!parentPresent || !childPresent` 即返回 L62 常量 `agentRunnerMissingPIDReportDetail` 的 uncertain）→ hang 负载缺 `stub-child.pid` → 检查 false。 | 本要素成立/不成立？漂移是否可完整归因于 `02d5797`（及其前置 `9ba621d`）？ |

## C. 判定问题清单（⑤ 项）

- **Q1 → 要素①**：本要素成立/不成立？理由？核对：`docs/validation/t027-at23-e2e.md` L60/L77；`directive-history/b4/archive/pack-verdict.json` L2/L1418；`l3-baseline-compare.txt` L1–2（声明对照物即该存档）。
- **Q2 → 要素②**：本要素成立/不成立？理由？核对：`tools/b4-e2e/main.go` L766 表达式、L690/L697/L1093；两份 f2-pos verdict 的 checks L136 断言 id；`git log --oneline e8cbff6..HEAD -- tools/b4-e2e/main.go` 输出。
- **Q3 → 要素③**：本要素成立/不成立？理由？核对：`l3-baseline-compare.txt` L5–7（now/baseline detail 原文）；`r08-r09-ruling.md` L45（方向定义原文）。
- **Q4 → 要素④**：本要素成立/不成立？理由？核对：B4 与 P61 两份 f2-pos verdict 的 L11–14；`l3-baseline-compare.txt` L22；相邻字段差异见 P61 文件 mapping 节（`treeStatus`/`errorEvidenceCount`）。
- **Q5 → 要素⑤**：本要素成立/不成立？理由？核对：`internal/adapters/agent_runner_postflight.go` L370–408 与 L62；`02d5797`/`9ba621d` 提交 diff 摘要；P61 报告 §2.4（L44–53，成因登记原文）。
- 最终判定（⑦ 另答）：五要素是否同时满足 ⇒ 豁免成立/不成立。**本文档不预填任何倾向性答案。**

## D. 判定后文档更新清单（起跑消息【四】16 项）

**豁免成立分支**：① `at23-restricted-candidate{,_EN}.md`：末尾新增 §10（⑦ 判定原文 + 五要素逐项 + 豁免成立 + 前提①视为满足 ⇒ AT-23 条件闭合、F2 修复片取消）+ §1 顶部 P63 时点注记（§1 原文不改）；② 证据包 P63-AT23-RECHECK 收口（00-briefing.md / design-r1 / 04-reviews / MANIFEST / SHA256SUMS / freedom-list / input-list）；③ `REMAINING_SLICES{,_EN}.md`：P63 ✅、层 1 第 11 项 ✅ 出队、第 8 项（AT-23）✅ 条件闭合、第 10 项（F2 修复片）取消；④ ADR-038 登记判定结论（四要素与「不得作为先例」已于步骤 0 登记，⑧ 仅追加判定结论）；⑤ DEV_PLAN：T027 段 AT-23 状态更新（PARTIAL→条件闭合）；⑥⑦⑧⑨ 不适用（准则/契约/架构/安全不改）；⑩ TEST_STRATEGY：AT-23 行更新为条件闭合；⑪ `s1-exit{,_EN}.md`：S-G1 闭合 + S-G7 依赖链同步；⑫–⑯ 不适用（v52 快照冻结、判据文本、支持矩阵、发行说明、操作/安装/需求/流程/结项模板）。

**豁免不成立分支**：① 同上但 §10 内容 = 豁免不成立、AT-23 保持 PARTIAL、前提①保持不满足；② 同；③ P63 ✅、第 11 项 ✅ 出队、第 8 项 🔄 保持（指向 F2 修复片）、第 10 项启动条件达成；④ ADR-038 同；⑤ DEV_PLAN 标「不适用 + 理由（状态未变）」；⑩ TEST_STRATEGY 同标不适用；⑪ S-G1 保持 PARTIAL 并指向 F2 修复片、S-G7 同步；⑫–⑯ 同。

## E. ⑦ 首轮输入包构成建议（供主控装配）

首轮（盲审与终审）仅附以下**原始材料**，**不附本文档 B/C 的对照表与问题清单**（锚定硬隔离）：
1. `docs/validation/r08-r09-ruling.md` R-08.5 全文（L36–50）；
2. `docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt` 全文；
3. B4 与 P61 两份 f2-pos `verdict-run.json`（四态四字段 + checks 段，含 `process-tree-taken-away` 两侧 pass/detail）；
4. `02d5797` 提交 diff 摘要（postflight Leg 5/7 + harness F1 面改动）与前置 `9ba621d` 摘要；
5. P61 报告原文节选：§1（L8–17）、§2.4（L44–53）、§8.3（L207–212）、§9（L376–378）；
6. B4 报告锚点：`t027-at23-e2e.md` L60/L77 + `pack-verdict.json` L2/L1418。

## F. 修订注记（⑤ r1 FINDINGS 整改；主控机械修订，无判定性改动）

- 修订 1（⑤ Medium）：§B④ 与 §C Q4 四态行号 L14–17 → L11–14（机械核对：两份 f2-pos verdict-run.json 的四个字段均位于 L11–14，字段值未变）。
- 修订 2（⑤ Low）：§B 表头「要素原文（逐字）」→「要素原文（压缩摘引；逐字原文见 §A）」（§A 已载逐字引用）。
- 修订 3（⑤ Low）：§B⑤ Leg 5/7 行号区间修正为 L370–386 / L392–411（机械核对：Leg 5 注释 L370、收尾 L386；Leg 7 注释 L392、收尾 L411）。
- 修订 4（⑤ Low）：§D 成立分支②补 `00-briefing.md`（步骤 0 已建，⑧ 装配勿因清单缺项遗漏）。
- 修订 5（⑤ Low）：§D④ 补注 ADR-038 四要素与「不得作为先例」已于步骤 0 登记，⑧ 仅追加判定结论。

DESIGN: DONE

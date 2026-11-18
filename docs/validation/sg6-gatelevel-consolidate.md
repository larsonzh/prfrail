# P80 验证报告：S-G6 门级侧合并判定与 S1 Exit 收口下一步清单

日期：2026-11-18。状态：完成。切片标识：P80 `SG6-GATELEVEL-CONSOLIDATE`（`[SLICE]`，S 档文档类）。

## 1. 切片标识与授权

- 判定对象：汇总 S-G1 至 S-G5 的门级判定与证据落点，判定 S-G6 门级侧是否成立，明确与 T018 任务级的并陈口径，产出 S1 Exit 收口下一步清单。
- 授权：2026-11-18 起跑消息即同轮授权；S-G5 口径裁定由用户同日作出并载入 v62 §六。
- 边界：本片是判定事实的承载，不是判据修改；p57 门表与 `s1-exit` 门表零改动，v62 快照零改动。
- 额度：探针 0；⑦ 为 1 终审 + 1 复审（S 档，§7.5）；盲审 1 次按 §7.7 按比例抽查另计，不占该额度。

## 2. 步骤 0 锚点实测

| 锚点 | 期望 | 实测 |
|---|---|---|
| `HEAD` / `origin/main` | `b5ff8ee` | `b5ff8ee` / `b5ff8ee` |
| 准则版本 | v1.69 | v1.69 |
| `SELFTEST` | 737/737 | 737/737 |
| 三作用域 | `TOTAL_FAIL=0` | index / tree / ci 均为 0 |
| `tmp` 顶层 | `gates-cache` + `.gitkeep`；排除缓存后 1 个文件 | 同左（切片目录在核对后建立） |
| `git status --porcelain -uall` | 0 条 | 0 条 |
| 滚动窗口 v60 / v61 / v62 | 242 / 214 / 217 | 242 / 214 / 217（CN 与 EN 对称） |

## 3. ① 方案书结论

方案书全文见证据包 `01-arch/design-r1.md`（`DESIGN: DONE`）。逐项要点：

- Q1：五门判定与落点成表；判据取自 p57 §4.2 逐字，判定取自 v62 §四 逐字，每门给出至少两处可用坐标。
- Q2：**S-G6 门级侧成立**；判据后半段「四道剩余门闭合」已满足，前半段「依赖字段满足」属任务级并陈保留。
- Q3：并陈口径为「门级侧成立 ≠ 判据通过 ≠ T018 完成」；T018 任务级受 T027 与产品所有者批准制约，`DEV_PLAN` 状态字段不变。
- Q4：下一步清单四项（S-G5 第二步、重评前提、P79 安全评审、S2 排期建议）各自给出动作、前置、证据形态与负责主体。
- Q5：第一层三项必触；第二层实读后全部不触（含 `DEV_PLAN{,_EN}.md` 判定不触）；第三层确定不触。
- Q6：三条检索独立执行，新增命中 0 条；无需升降级动作。

## 4. Q1：S-G1 至 S-G5 门级判定与证据落点

| 门 | 判据（p57 §4.2 引文） | 当前判定（v62 §四） | 证据落点 | 来源片 |
|---|---|---|---|---|
| S-G1 | 「CLI 可真实启动经准入的 AI 候选并计费」 | 不通过（Q2=(a) 判据字面） | `p57:88`；`v62:115`；P76 报告 §4 | P76 |
| S-G2 | 「`supervised-black-box` 实现完成且 AT-24 判定成立」 | 通过（Q2=(a) 判据字面） | `p57:89`；`v62:116`；`s1-exit:52`；`DEV_PLAN:121` | P76 |
| S-G3 | 「最终 ZIP 绑定发行记录并复验」 | 通过（P78 闭环） | P78 报告 §5 / §7；`p57:90`；`v62:117`；`s1-exit:66` | P78 |
| S-G4 | 「ADR-010 获产品所有者批准」 | 通过（ADR-010 已翻 `APPROVED`） | `ADR_REGISTER:20`；P77 报告 §5；`p57:91`；`v62:118` | P77 |
| S-G5 | 「产品所有者作出批准」 | 发行部分已批准（第一步） | `v62:119`；`v62:196`；ADR-046（`ADR_REGISTER:57`）；P76 报告 §12 | 用户裁定 |

口径说明：上表的坐标取自 2026-11-18 磁盘实读；v62 与 p57 两处均为只读引用，门表未改一字。

## 5. Q2：S-G6 门级侧判定

- 判据原文（`p57:93`）：「依赖字段满足且四道剩余门闭合」；同一行的解锁条件列写明「门级与任务级不同源」的约束。

| 判据成分 | 判定 | 依据 |
|---|---|---|
| S-G2 门级通过（P76） | 满足 | P76 报告 §4（`s1-exit-closeout-assessment.md` §4）；`v62:116`；`s1-exit:52` |
| S-G3 已通过 | 满足 | P78 报告 §7（`release-zip-binding.md:69`–`:73`）；`v62:117` |
| S-G4 已通过 | 满足 | P77 报告 §5（`adr010-approval.md:90`–`:97`）；`ADR_REGISTER:20` |
| S-G5 发行部分已批准（第一步） | 满足（用户 2026-11-18 裁定足以支撑门级侧） | `v62:196`；ADR-046（`ADR_REGISTER:57`） |
| 依赖字段 | 任务级成分，门级侧不代判 | `DEV_PLAN.md:94`（依赖字段逐字）；`DEV_PLAN.md:96`（状态 `IN_PROGRESS / BLOCKED`） |

门级侧判定的精确表述（本节即落盘文本，台账完成记录同源引用）：

> **S-G6 门级侧成立**（不使用「S-G6 通过」）：判据「依赖字段满足且四道剩余门闭合」的后半段「四道剩余门闭合」已满足——S-G2 门级通过（P76 Q2=(a)）、S-G3 通过（P78）、S-G4 通过（P77）、S-G5 记「发行部分已批准（第一步）」（用户 2026-11-18 裁定，v62 §六）；判据前半段「依赖字段满足」属任务级（`DEV_PLAN.md:94` 的依赖字段含 T027、T028），任务级侧不成立并陈保留（T018 不勾选，`DEV_PLAN.md:96` 状态字段不变）。门级侧成立与任务级并陈按 P76 Q8=(c)。本判定由本片报告承载，p57 门表与 `s1-exit` 门表零改动。

可证伪条件（任一观察成立即推翻本判定，须停下重判）：

1. S-G2 / S-G3 / S-G4 / S-G5 任一门的门级判定被后续证据或用户裁定回退。
2. 用户作出更新裁定，明确「发行部分已批准（第一步）」不足以支撑门级侧收口。
3. P76 Q8=(c) 的并陈口径被新 ADR 或用户裁定替换。
4. ADR-046 的回退条款被触发（后续证据推翻 Q2 判据字面口径，须另立片并登记新 ADR）。

## 6. Q3：T018 任务级并陈口径

并陈表述的精确文本（本节即落盘文本）：

> **S-G6 门级侧成立 ≠ 判据通过 ≠ T018 完成。** T018 任务级仍受依赖字段制约：T027 未实现（`BLOCKED / NOT IMPLEMENTED`）、AT-23 整体未通过（`s1-exit:52`：「T027 未完成、AT-23 整体未通过」）、产品所有者 S1 exit 批准缺失，因此 T018 不勾选，`DEV_PLAN.md:96` 的状态字段保持 `IN_PROGRESS / BLOCKED`；`DEV_PLAN.md:98` 的 P78 注记（「T018 勾选状态不变」）继续有效。门级侧成立只记录四道剩余门的闭合进展，不折算为 T018 的任何推进。

P76 Q8=(c) 原文引证（`s1-exit-closeout-assessment.md` §10）：

- 「实测字段：`DEV_PLAN.md:94` 的 T018 依赖为「T017、T025、T026、T027、T028 及 whois 输入授权」；`:96` 的状态为 `IN_PROGRESS / BLOCKED`。」
- 「口径（用户裁定 c）：门级与任务级并陈，不给单一结论；本片零改动。」

第三句「S-G6 门级保持不可达至 S2」是 2026-11-16 时点表述：当时四道剩余门未闭合。P77 / P78 已闭合 S-G3 / S-G4，S-G5 口径已由用户 2026-11-18 裁定，v62 §四 已把 S-G6 更新为「门级侧可收口（待 P80 合并判定）」；P76 §12 收口路线图第 5 项即「门级侧并陈推进，任务级保持不可勾选」，本片执行的正是该顺序位。

`DEV_PLAN` 两处现状：`:96` 为带日期的历史状态句加口径注记，`:98` 为 P78 加法式注记；两处均为历史记录，本片不改写、不翻转勾选状态。

## 7. Q4：S1 Exit 收口下一步清单

| 项 | 动作 | 前置 | 证据形态 | 负责主体 |
|---|---|---|---|---|
| S-G5 第二步 | S2 解锁后作出 S1 Exit 整体批准并登记四要素 | S-G1 / S-G8 重评通过；T027 完成；AT-23 整体通过；独立安全评审已补评 | 批准裁定原文（裁决人 / 日期 / 理由 / 依据） | 产品所有者 |
| S-G1 / S-G8 重评前提 | S2 enforcement 接入后重跑真实 E2E 并重评 | R-08、R-09、T027 装配、兼容候选齐备，且 enforcement 已接入 | 重评报告与真实 E2E 证据链 | 后续切片（S2 后） |
| S-G7 | 不重评；S1 内保持不可达 | 依赖 S-G1 至 S-G6 全绿 | 本清单即说明载体 | 后续切片（S2 后） |
| P79 独立安全评审 | S2 启动前完成一次独立评审并登记结论 | S2 启动前；范围 = `PRODUCT_REQUIREMENTS` §7 的 S0 到 S1 安全门禁 | 评审报告；`SECURITY.md` 注记更新；ADR-047 状态更新 | 后续切片 P79 |

S2 启动排期建议（**供用户参考，不作裁定**）：先做 P79 独立安全评审，再启动 S2 与 enforcement 接入，随后重评 S-G1 / S-G8，最后走 S-G5 第二步。ADR-047 的到期日为 S2 启动日与 2027-03-31 中较早者（`ADR_REGISTER:58`），到期未补评则不得启动 S2 的相关实测。

两项对象不得混用：P79 的独立安全评审是 `PRODUCT_REQUIREMENTS` §7 的历史欠账（从未执行）；`s1-exit` 门表「独立评审」行（`:59`）指 T017 基线及 Node 24 替换的切片审查。两者证据不可互换。

## 8. Q5：文档同步判定（三层归属）

第一层（本片必触）：

| 文档 | 动作 |
|---|---|
| `docs/validation/sg6-gatelevel-consolidate{,_EN}.md` | 新建报告对（本文件） |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P80 段登记（步骤 0c）与 ⑧b 完成记录 |
| `docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/` | 证据包（建包于步骤 0，按轮追加） |

第二层（实读后全部不触）：

- `DEV_PLAN{,_EN}.md`：不触。判定由本片报告承载；`:96` 与 `:98` 均为带日期的历史记录，`DEV_PLAN` 不存在需要更新的判据语义。若判应触，本片将按 §1.7.2 停下升级；本片不触发该升级。
- Q6 三条检索命中的其它文档：全部为已知历史时点表述或本片载体，逐条判定不触（明细见 §9）。

第三层（确定不触）：准则、`CONTRACTS`、`SECURITY`、`ARCHITECTURE`、`TEST_STRATEGY`、`BUSINESS_WORKFLOWS`、`CODING_CONVENTIONS`、`DOCUMENTATION_PLAN`、`RFC-*`、`CHECKLIST_v62{,_EN}`、v50 至 v62 任何快照、p57 门表、`s1-exit` 门表、产品代码、CI、其它 ADR。

## 9. Q6：组合检索排查

三条检索各自独立执行（不联查）；基线取 HEAD 的已跟踪文件，新增取本片暂存的新增文件：

| 检索词 | HEAD 基线命中 | 基线后新增载体 | 其它文件新增 |
|---|---|---|---|
| `S-G6` | 75 | 本片报告对、台账 P80 段、证据包 | 0 |
| `SG6` | 8 | 本片报告对、台账 P80 段、证据包 | 0 |
| `门级侧` | 6 | 本片报告对、台账 P80 段、证据包 | 0 |

归属说明：HEAD 的全部命中均为已知载体（p57 / `s1-exit` / P76 / P75 / P73 / P61 报告、v59 至 v62 快照、台账、冻结证据树）；两侧判定按文件集合比对得出而不依赖命中条数，故本片载体之外的新增为零，无升降级动作，第二层「不触」判定与检索结果一致。

## 10. ② 执行结果与逐文件改动

- 报告对：`docs/validation/sg6-gatelevel-consolidate{,_EN}.md` 新建，CN 与 EN 逐行配对，六形状全等。
- 台账：P80 段在步骤 0c 登记（CN 与 EN 各加同样条数），⑧b 追加完成记录并翻状态。
- 证据包：三件初始载荷建于步骤 0，按轮追加方案书、批量枚举、审查输出与门禁原始输出，⑧c 重算校验和。

| 文件对 | 增 | 删 |
|---|---|---|
| `docs/validation/sg6-gatelevel-consolidate{,_EN}.md` | 见 ⑨ 实测 | 见 ⑨ 实测 |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | 见 ⑨ 实测 | 见 ⑨ 实测 |

计数口径：两侧增删逐对相等由 `G3(a)` 的原始输出直接给出，本表不复制具体数字，避免与实测读数互相矛盾。

## 11. ③ 一致性检查

- S-G6 门级侧判定与 `p57` §4.2 判据字面一致：判据只要求四道剩余门闭合与依赖字段满足，本片对两段分别给出判定与并陈保留。
- 与 `s1-exit` 门表不冲突：门表行与判定词零改动；`:52` 的「AT-23 整体未通过」与本片「T027 未实现」表述同向。
- 与 v62 §四 不冲突：v62 记「门级侧可收口（待 P80 合并判定）」，本片给出该合并判定的结论。
- 与 ADR-046 不冲突：Q10=(C) 的两步批准结构被完整保留，第一步记「发行部分已批准」，第二步仍待 S2 之后。
- 与 P76 Q8=(c) 不冲突：并陈口径逐字沿用，未给单一结论。
- CN 与 EN 严句对齐：报告对逐行配对，术语与判定词逐处同构；台账段两侧条数相等。

## 12. ④ 门禁读数

| 判据 | 读数 |
|---|---|
| `--scope=index` | `TOTAL_FAIL=0`（`changed=19`） |
| `--scope=tree` | `TOTAL_FAIL=0`（`changed=19`） |
| `--scope=ci` | `TOTAL_FAIL=0`（本地退化口径，`changed=2`） |
| `--selftest` | `SELFTEST: PASS 737/737` |

读数作用域：`06-gates` 的四份原始输出取自本轮快照（`changed=` 见各文件首段），`artifacts` 围栏行数与 `G7-b` 计数同批；此后仅追加审查载荷与台账完成记录，门禁结论不变，⑨ 以最终快照复跑并在 ⑩ 回执载明最终 `changed=`。

`--scope=ci` 在无 CI 环境变量时退化为上一提交的提交范围，只覆盖 `HEAD^1..HEAD`，因此其 `changed` 数不反映本片暂存集；本片改动的真实覆盖由 `--scope=tree` 与 `--scope=index` 给出。

读数口径：本片为纯文档切片，产品代码零改动，本片不新增任何构建动作；三作用域与自检的原始输出归档于证据包 `06-gates/`。

## 13. 审查链

- ⑤ 预审：`DeepSeek V4 Pro`（`deep-reasoner`，只读）；轮次与结论记录于 §13.1。
- ⑥ 跳过：§7.2.1 三条件全部成立（改动只落文档与台账类；不含脚本或代码语义变更；不涉状态机、门禁判据、证据模型、角色与授权契约语义），由 ⑦ 确认。
- ⑦ 审查链：盲审（§7.7 按比例抽查，额度另计）+ 终审 + 复审；已完成的轮次逐条列入 §13.1，未完成的轮次不预先记述。

### 13.1 审查轮次记录

| 环节 | 轮次 | 结论行 | 发现计数 | 处置 |
|---|---|---|---|---|
| ⑤ 预审 | r1 | `PRE-REVIEW: PASS WITH FIXES` | High 0 / Medium 0 / Low 3 / Info 3 | 三条 Low 已按最小修复整改：检索计数改按 HEAD 基线表述、§5 补 S-G2 门级成分、围栏按暂存真值重建；三条 Info 逐条登记 |
| ⑦ 盲审 | r1 | `BLIND REVIEW: FINDINGS` | High 0 / Medium 1 / Low 1 / Info 0 | 两项均属时点类（轮次记录与台账态尚未同批落定），不触及判定本体；本行与 ⑧b 台账翻转即为其闭合证据 |
| ⑦ 终审 | r1 | `FINAL REVIEW: FINDINGS` | High 0 / Medium 1 / Low 0 / Info 0 | Medium 为门禁归档与快照不同批；已按同一快照重跑并回填 `06-gates`，交 r2 复核 |
| ⑦ 复审 | r2 | `RE-REVIEW: FINDINGS` | High 0 / Medium 1 / Low 0 / Info 0 | Medium 指向 §12 缺 `--scope=ci` 本地退化说明；该说明已补入 §12，但 S 档额度（1 终审 + 1 复审）已用尽 ⇒ 按规则停下升级用户裁决，不自行判定通过 |
| ⑦ 收窄读 | r1（授权例外） | `RE-REVIEW: PASS` | High 0 / Medium 0 / Low 0 / Info 0 | r2 的唯一 Medium 判 `CLOSED`；异常处理章真实性核对通过；该授权仅本片有效、不得作先例（§14.3） |


## 14. 异常处理（R 计数 / 风险标记 / D 澄清）

### 14.1 R 逐次计数（§9.2 为唯一真值源）

| 编号 | 环节 | 事件 | R | 处置 |
|---|---|---|---|---|
| R-1 | ⑦ 终审 r1 | 1 项 Medium（`06-gates` 归档与快照不同批） | +1 | 就地整改：同批重跑并回填 `06-gates`；按 §5.3 重跑计入 R |
| R-2 | ⑦ 复审 r2 | 1 项 Medium（§12 缺 `--scope=ci` 退化说明） | +1 | 就地整改：§12 补读数作用域三段；按 §5.3 重跑计入 R |
| 合计 | — | — | **R = 2** | 达到 §9.2 的高风险阈值 |

不计入 R 的两轮：⑤ 预审 r1（3 项 Low / 3 项 Info，无 Medium+）与 ⑦ 盲审 r1（1 项 Medium + 1 项 Low，均属时点类并已闭合）——盲审按 §7.7 为抽查读，其发现单独记录在 §13.1，不进入 R 序列。

### 14.2 高风险切片标记（§9.2：R≥2）

- 标记：**本片为高风险切片**（依据 §5.2 高风险行与 §9.2 第 3 条）。
- 额度：**S 档额度不突破**（起跑消息预置 1 终审 + 1 复审，盲审另计 1）。§5.2 高风险行虽要求 ⑤ 复查为必，本片**不自行升为 1+3**，处置 = **停下升级用户裁决**。
- 用户裁决（2026-11-18 同轮）：**选 A —— 授权一次收窄读**（§7.8.3「用户授权例外」行）：额度 = ⑦ Codex 常规档 1 次，仅限本片、不得作先例（同 GATES-EXT 先例口径）。授权原文与输入收窄声明见 §14.3，收窄读结论见 §13.1。

### 14.3 授权原文与输入收窄声明

- 授权原文（2026-11-18 用户消息逐字摘录）：「用户裁定：选 A——授权一次收窄读（§7.8.3「用户授权例外」行），额度 = ⑦ Codex 常规档 1 次，仅限本片、不得作先例（GATES-EXT 先例口径）。」
- 输入收窄：本次读的输入 = ① §12 的三条作用域句；② `06-gates/` 的同批读数；③ `impl.numstat`；④ 本章。**不附**切片定义、完整变更集、⑤ 清单、盲审与终审的发现清单。
- 任务收窄：只核「r2 唯一 Medium（§12 的 `--scope=ci` 退化说明）是否 CLOSED」与「本章补丁是否如实」。
- 效力边界：该授权**不推广**，仅本片有效；结论行取自 §4.1 闭集，落 `04-reviews/narrow-r1.md`。

### 14.4 D 定性澄清

- 事实：r1 的整改脚本 `fix-final1.js` 已编写但**未执行**即派发 r2，导致 r2 依据未整改的正文复审。
- 定性：该失误属**编排时序失误**（违反 ⑤→⑦ 之间「整改后必须原位重做配套动作」），**不是临时件位置违规**——`tmp/` 全部落在本片目录，OB-139 不触发。
- 计数：本片**保留自严 D = 1**，计入位置不在临时件违规序列，而在本节；标注「**自严口径、交用户追认**」。若用户按严格口径追认为 D = 0，则以本节作为该异常的替代附注（两种口径的差别仅为是否把编排失误计入 D）。

## 15. 已知限制与残留

- `p57` 与 `s1-exit` 门表保留 P76 时点的判定词，与本片「S-G6 门级侧成立」并陈；门表零改动是 P76 的明令，本片不翻转任何门表行。
- `DEV_PLAN.md:96` 与 `:98` 保持原状；S-G6 门级侧判定的现役承载位是本片报告与台账 P80 段完成记录。
- 起跑消息与实测的六条表述差异已逐条判定：结论影响为零，其中四条属表述层，已在 §4 / §5 / §6 对应处按实读坐标登记。
- 具体差异包括：`ADR_REGISTER` 的 ADR-009 在 `:19` 而 ADR-010 在 `:20`；v62 §一 的 HEAD 是生成时点锚点；2026-11-17 与 2026-11-18 两条裁定的对象不同；T028 已完成而现役任务级制约为 T027 等三因；P76 裁定 c 的出处是 §10 正文而非 §15 表。
- 证据包在步骤 0 建包时其元数据尚未入索引，`G8-a` 在 ⑧b 暂存前为时序红因，不是缺陷。
- 双 BOM 的历史遗留（t025 相关文件）属 HEAD 既有事实，不在本片范围。
- 本片曾在 ⑦ 复审之后停在用户裁决点（r1 与 r2 各报 1 项 Medium，S 档额度用尽）；用户 2026-11-18 裁定选 A —— 一次收窄读，见 §14.2 与 §14.3。

## 16. 计量

| 环节 | 次数 |
|---|---|
| ① 架构 | 1（`deep-reasoner` + `DeepSeek V4 Pro`） |
| 批量辅助 | 1（`GLM-5.3-Flash`） |
| ⑤ 预审 | 1 次（`deep-reasoner` + `DeepSeek V4 Pro`） |
| ⑥ 独立扫描 | 0 次调用（§7.2.1 跳过，经 ⑦ `SKIP CONFIRMED`） |
| ⑦ 审查 | 4 次调用：盲审 1（§7.7 另计）+ 终审 1 + 复审 1 + 收窄读 1（用户授权例外） |
| 探针 | 0 |
| R / D | 见 §14 异常处理（R = 2、D = 1 自严口径） |

## 17. 提交与 CI

- 主提交 `e00d79b`（`docs: consolidate the S-G6 gate-level verdict and next steps (P80)`）；推送 `origin/main`，范围 `b5ff8ee..e00d79b`（未推 gitee）。
- CI run `44211368890` 首跑双腿 success：windows `4m55s`、ubuntu `3m18s`；两腿 `SELFTEST: PASS 737/737`、`GATE REPORT scope=ci changed=20`、`TOTAL_FAIL=0`。
- 包级 `ok` 计数：ubuntu 腿 21、windows 腿 15，`FAIL` 均为 0（ubuntu 腿含 Race 与 Contract fixtures 两段）。
- 回填提交（本段）自身的 run 只作观察记录，不再为它新增提交（§11.2）。

```artifacts
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/00-briefing.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/01-arch/design-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/02-batch/batch-01-sg6-sources.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/blind-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/final-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/final-r2.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/narrow-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/freedom-list.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/input-list.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/MANIFEST.md	repo
docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/SHA256SUMS.txt	repo
docs/validation/sg6-gatelevel-consolidate_EN.md	repo
docs/validation/sg6-gatelevel-consolidate.md	repo
```

# batch-01 — T027-B / AT-23 / B1-B2 事实枚举（GLM-5.3-Flash 原文归档）

切片：P75-R3-REAL-CLI-ASSESSMENT。派单：`deep-reasoner` + `GLM-5.3-Flash (glm)`。实得字符数：5952。原文如下（逐字，仅去除首行标题层级冲突）。

---

① T027-B 主线（路径省略 docs/ 前缀）

判定基线（P38）：
- `validation/t027b-tactical-blueprint.md:10` — 「Astra 判 VERDICT: REVISE——3 项 BLOCKED（C1 / C4 / C6）、4 项 PARTIAL（C2 / C3 / C5 / C7），并给出新漏洞 N1–N4」
- `validation/t027b-tactical-blueprint.md:24` — 「N1 为 Job 归属建立前的宿主崩溃窗口；N2 为目录 helper 返回后的祖先替换；N3 为 Reopen 写死 attempt；N4 为 ProbePermit 额度可重复取得」

C1–C7 当前状态（逐项一行）：
- `DEV_PLAN.md:181` — 「P40 锁序：C1（已完成，报告 …p40-lock-reaper{,_EN}.md）」
- `DEV_PLAN.md:182` — 「P41 归属：C2 与 N1（已完成，报告 …p41-spawn-ownership{,_EN}.md）」
- `DEV_PLAN.md:183` — 「P42 绑定：C3 与 N2（已完成，报告 …p42-evidence-bind{,_EN}.md）」
- `DEV_PLAN.md:184` — 「P43 封存：C4 与 C5 与 N3（已完成，报告 …p43-state-seal{,_EN}.md）」
- `DEV_PLAN.md:185` — 「P44 门序：C7 与 N4（判定实现与测试已完成…；端到端接线已由 P69 完成）」
- `DEV_PLAN.md:180` — 「P39 持久预算：C6（DurableCostLedger 与跨进程共享锁，锁内重读后决策…）」
- `validation/p44-gate-readiness.md:166` — 「端到端接线已由 P69 T027B-E2E-WIRING 完成（报告：docs/validation/t027b-e2e-wiring{,_EN}.md）」

N1–N4：定义见 t027b-tactical-blueprint.md:24；落地 N1→P41（DEV_PLAN.md:182）、N2→P42（:183）、N3→P43（:184）、N4→P44（:185）。

R-1（P44 端到端接线）：
- `validation/p57-closeout-assessment.md:68` — 「R-1 | P44 的端到端接线 | F4：判定实现与测试已完成，接线未做 | 无队列位…」
- `validation/t027b-e2e-wiring.md:7` — 「目标：把 P44 的 readiness / permit 判定接到生产装配路径与 CLI 入口，收口 P57 §3.3 的 R-1」
- `validation/t027b-e2e-wiring.md:8–9` — 「结论：DESIGN: DONE（①）；接线落地 internal/console/** 5 个生产文件…审查链：⑤ PRE-REVIEW: PASS…⑥ INDEPENDENT SCAN: PASS（0 发现）；⑦ FINAL REVIEW: PASS（0 发现…）」
- `validation/p57-closeout-assessment.md:72` — 「R-1 已由 P69 T027B-E2E-WIRING 收口（报告：docs/validation/t027b-e2e-wiring{,_EN}.md）」

R-2（OB-55 归属）：
- `validation/p57-closeout-assessment.md:69` — 「R-2 | enforcement 生产落地 B 段 | F6：归属三选一未裁 | OB-55 待议 + R-06 待裁…」
- `validation/p57-closeout-assessment.md:74` — 「R-2 已由 P73 裁定（2026-11-14）：OB-55 归属 = 方案 b（Sandbox 接口提前到 S1，实施片 P74-SANDBOX-INTERFACE 待起跑）…」
- `validation/t027b-admission.md:7` — 「OB-55 裁定 = 方案 b（提前落平台无关 Sandbox 接口）…三条附条件全部接受」
- `validation/t027b-admission.md:39–41` — 「裁定 1：OB-55 = 方案 b；三条附条件全部接受…」「裁定 2：B2 §8 第 8 条以加法式裁定注记处置，不改原行」「裁定 3：实施片 = M 档独立片 P74-SANDBOX-INTERFACE，P73 不启动」
- `DEV_PLAN.md:201` — 「P74 SANDBOX-INTERFACE：Sandbox 检查层落地…装配检查点（均在 permit 之前）+ 触发源 PRFRAIL_SANDBOX_CHECK（S1 默认 off…）；真实 enforcement 边界仍留 S2」
- `validation/sandbox-interface.md:126` — 「切片：P74 SANDBOX-INTERFACE 闭环（层 1 主线 B 段引入独立片，实现类，M 档）」

R-3（真实 CLI 端到端）：
- `DEV_PLAN.md:153` — 「余项三叉为 R-1 P44 端到端接线（无队列位）、R-2 OB-55 归属（待 R-06 裁定）、R-3 真实 CLI 端到端（层 1 第 8 项，待 R-08 与 R-09）」
- `validation/p57-closeout-assessment.md:70` — 「R-3 | 真实 CLI 端到端（AT-23 真实候选腿） | F7 + docs/validation/s1-exit.md:43 | 层 1 第 8 项 AT23-RESTRICTED-CANDIDATE | R-08 / R-09 裁定 + 兼容候选」
- `validation/t027b-admission.md:45` — 「裁定 7：AT-23 重评（R-3）与 B 段无强制先后，各走各的」
- `validation/sandbox-interface.md:133` — 「层 1 B 段主线状态：抽象层已落，真实 enforcement 边界仍在 S2；S-G1 / S-G8 未解锁」

② AT-23 阶段 1 条件闭合与 S-G1 / S-G8（省略 docs/ 前缀）

- `validation/at23-restricted-candidate.md:19` — P61 时点：「前提①「机制腿全绿」：不满足——层 3 机制腿 F2 正例一处断言由基线 true 变为 false…」
- `validation/at23-restricted-candidate.md:20` — P61 时点：「前提②「真实候选腿实际跑过（≥3 次，覆盖启动 / 执行 / 停机）并记录 unsupported / unknown」：满足…」
- `validation/at23-restricted-candidate.md:423` — 「⑦ 判定：盲审四问均成立 → 终审 FINAL: PASS WITH FIXES…→ 前提①「机制腿全绿」满足」
- `validation/at23-restricted-candidate.md:499` — 「AT-23：PARTIAL → 阶段 1 条件闭合达成；阶段 2 不适用（用户 2026-11-09 裁定 A1）（前提① = 机制腿全绿，本片满足；前提② = P61 已满足）」
- `validation/at23-restricted-candidate.md:501` — 「S1 Exit：…不声称 S1 Exit 整体通过（其余门 S-G2~S-G6、S-G8 仍待）」
- `validation/s1-exit.md:43` — 「AT-23 | 阶段 1 条件闭合达成（2026-11-09，F2 修复片收口；阶段 2 不适用…）…⇒ AT-23 阶段 1 条件闭合达成（前提② P61 已满足）…」
- `validation/s1-exit.md:52` — 「S-G1 依据 = AT-23 阶段 1 条件闭合（…依据本片 ⑦ 判定，而非 S1 Exit 整体通过）…不得声称 S1 Exit 整体通过」
- `validation/p57-closeout-assessment.md:86` — S-G1 判定原文：「S-G1 AT-23 真实候选腿与整体 | CLI 可真实启动经准入的 AI 候选并计费 | 不可达 / 未通过：docs/validation/s1-exit.md:43 | … | 三者齐备后重跑真实 E2E」
- `validation/p57-closeout-assessment.md:93` — S-G8 判定原文：「S-G8 真实 CLI 端到端 | 受监督路径可在真实 CLI 与真实候选下端到端运行 | 未通过：docs/DEV_PLAN.md:119（D-18 明写依赖 T027-B） | P5 的 R-3 | 与 S-G1 同源解锁」

八门全表（p57 §4.2，P57 时点，逐门一行）：
- :86 S-G1 AT-23 真实候选腿与整体 — 不可达 / 未通过
- :87 S-G2 AT-24 与 T028 — 口径滞后（DEV_PLAN.md:119 已 COMPLETED；s1-exit.md:44 仍 BLOCKED）
- :88 S-G3 最终 ZIP 绑定与复验 — 未通过（s1-exit.md:56 PARTIAL + ADR-009）
- :89 S-G4 EOL 通知方式 — 未通过（ADR-010 PROPOSED / NOT APPROVED）
- :90 S-G5 产品所有者 S1 exit 批准 — 未知（人工门）
- :91 S-G6 T018 整体 — 不可达（DEV_PLAN.md:94、:96）
- :92 S-G7 S1 exit 总体（RFC §14） — 不可达（s1-exit.md:52、:56–:58、:60）
- :93 S-G8 真实 CLI 端到端 — 未通过（D-18 依赖 T027-B）

③ B1 / B2 对真实 CLI 端到端（省略 docs/ 前缀）

- `validation/t027-b1-candidate-live-discovery.md:4` — 「T026 的两项 unsupported（toolControl、networkControl）在 1.0.85 上仍然存在…B1 的验收按「明确证明无可接受 candidate」分支收口：保持 blocked，不进入 B2 的 candidate-native 路线；B2 应转向外部 OS enforcement 证明」
- `validation/t027-b1-candidate-live-discovery.md:111` — 「已评估范围内无可接受候选：…更新没有覆盖 T026 的两项 unsupported…未评估、未排除…BYOK deepseek-anthropic…与非 Copilot CLI 候选族」
- `validation/t027-b1-candidate-live-discovery.md:113` — 「B2 的路线据此收敛：candidate-native 路线在本候选族上不成立；B2 应转向外部 OS enforcement（AppContainer / WDAC / 防火墙 / Job Object 一类可证边界）」
- 「10 verified / 2 unsupported」字样：B1、B2 两文件内未出现；出现于 `validation/s1-exit.md:3`（「矩阵为 10 verified、2 unsupported、0 unknown、blocked…固定候选不兼容」）与 s1-exit.md:52（「固定候选为 10 verified、2 unsupported、0 unknown，整体不兼容」）

外部 enforcement 边界适用范围原文（`validation/B2-EXTERNAL-ENFORCEMENT.md`，下称 B2）：
- `B2:12` — 「采用方案 B：外部 OS 强制。受管进程运行在零网络能力 AppContainer 中…策略判定发生在操作系统边界与代理两处，均不由受管进程自身控制」
- `B2:18` — 「边界内：AppContainer 令牌、其网络能力集（无公网/LAN/DNS）、其文件系统授权集、其进程创建能力」
- `B2:19` — 「边界外：宿主侧的一切，包括代理进程自身的出网、宿主凭据库、未被授权的宿主文件」

S1 内无真实 enforcement 的边界表述原文（B2 §8 附近）：
- `B2:162–164` — 「零能力容器内的候选会话失败：候选 1.0.83 在取得模型列表后报「工作区不可访问」…后果：B4/B5 不能假设候选 CLI 可在零能力边界内完整运行」
- `B2:184–185` — 「S2 待做（只登记，不实做）…S1 不实做，仅在 S2 启动时按此接口扩展；在 S1 提前落地该接口会产生无人调用的死代码，反而增加维护面」
- `B2:188` — 「P73 更新（2026-11-14）：本条「Sandbox 接口推迟到 S2」已由 R-06…与 P73 T027B-ADMISSION 更新——接口提前到 S1…本条原行不改（§11.1）」
- `validation/sandbox-interface.md:105` — 「B 段不解锁 S-G1 / S-G8：本接口只落检查层抽象，S1 Exit 八门判定不变；真实 enforcement 边界留 S2」

（实得字符数：5952）

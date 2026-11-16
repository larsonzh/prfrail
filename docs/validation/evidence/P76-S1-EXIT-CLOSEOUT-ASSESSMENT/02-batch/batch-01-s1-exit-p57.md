# batch-01 — s1-exit / p57 §4.2 八门现状枚举

生成者：deep-reasoner + GLM-5.3-Flash（P76 批量辅助 1，事实枚举，不含判定）

输入（本批实际读取；所有行号均以终端磁盘读取逐一实锚）：

1. `docs/validation/s1-exit.md`（全文 76 行）
2. `docs/validation/s1-exit_EN.md`（全文 76 行；与 CN 侧行号完全对称）
3. `docs/validation/p57-closeout-assessment.md`（全文 175 行；本批实锚区间 = 行 62–108，覆盖 §3.3 尾部 / §4 / §4.1 / §4.2 / §4.3）
4. `docs/validation/p57-closeout-assessment_EN.md`（全文 175 行；行号与 CN 完全对称）
5. 交叉锚点实锚读取：`docs/ADR_REGISTER.md`（行 15–25、50–60 区间）、`docs/DEV_PLAN.md`（行 90–100、115–122、200–206 区间）、`docs/validation/r3-real-cli-assessment.md`（行 158–168）、`docs/t027/CHECKLIST_v59.md`（行 1、3、41、52）、`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md`（行 23、27）
6. 只读 git 读取（未执行任何 git 写操作）：`git log --oneline -3 -- docs/validation/s1-exit.md`、`git log --oneline -3 -- docs/validation/p57-closeout-assessment.md`
7. 磁盘计数核验（Select-String）：`s1-exit{,_EN}.md` 中「P75」= 0 命中；`s1-exit.md` 中「尚未实现」= 0、「仍未实现」= 0；`s1-exit_EN.md` 中「not implemented」= 0；`ADR_REGISTER.md` 中「PROPOSED / NOT APPROVED」唯一命中行 = 20；`s1-exit{,_EN}.md` 中「DONE」= 0 命中

取证方式说明：终端对超长行的回显存在折行边界空格伪影（同一行在未折行磁盘视图中无该空格）。本文件所有逐字引文以未折行磁盘视图（read_file / grep 磁盘内容）取文，行号以终端磁盘读取为准，两者逐行交叉核对内容一致。

## 一、s1-exit{,_EN} 八门逐门现状

s1-exit 两文件的 RFC §14 门表位于行 50–60（表头 50、分隔行 51、9 个门行 52–60）。`S-G1`…`S-G8` 编号定义于 p57 §4.2；s1-exit 文件自身不使用该编号（行 62 明写「本文件不引入该编号，亦不建立其与门表行名的映射」）。下表按 p57 §4.2 各门引用的锚点行归纳现状；被引行的逐字全文见「一·附 A（CN）」「一·附 B（EN）」与第五节。

| 门 | CN 现状表述（逐字摘录） | CN 行号 | EN 现状表述（逐字摘录） | EN 行号 |
|---|---|---|---|---|
| S-G1 | AT 汇总表 AT-23 行（p57 §4.2 引 `s1-exit.md:43`）：状态列 = 「阶段 1 条件闭合达成（2026-11-09，F2 修复片收口；阶段 2 不适用，用户 2026-11-09 裁定 A1）：⑦ 终审 `FINAL: PASS WITH FIXES` + 盲审四问均成立 ⇒ 前提①满足 ⇒ **AT-23 阶段 1 条件闭合达成**（前提② P61 已满足）；P63 豁免判定与 P61 PARTIAL 为历史时点结论、原文保留」；证据列含「真实候选腿实际跑过启动 / 执行 / 停机并记录 unsupported，不构成真实候选兼容性证明」；另门表核心任务链行（行 52）内含「⇒ **S-G1 依据 = AT-23 阶段 1 条件闭合**（阶段 2 不适用，用户 2026-11-09 裁定 A1；依据本片 ⑦ 判定，而非 S1 Exit 整体通过）」 | 43、52 | AT-23 行：状态列 = 「Stage-1 condition closure achieved (2026-11-09, F2 fix slice closed; stage 2 not applicable, user ruling A1, 2026-11-09): ⑦ final `FINAL: PASS WITH FIXES` + the blind review's four questions all hold ⇒ premise ① met ⇒ **AT-23 stage-1 condition closure achieved** (premise ② already met by P61); the P63 waiver verdict and the P61 PARTIAL are historical point-in-time conclusions and their text is preserved」；证据列含「the real candidate legs actually ran start / execution / stop and recorded unsupported; this does not constitute real-candidate compatibility proof」；行 52 内含「⇒ **S-G1 basis = AT-23 stage-1 condition closure** (stage 2 not applicable, user ruling A1, 2026-11-09; per this slice's ⑦ verdict, not S1 Exit passing as a whole)」 | 43、52 |
| S-G2 | AT-24 行（p57 §4.2 引 `s1-exit.md:44`）：状态列 = 「条件闭合（2026-11-11，P64 验收；⑦ 复审 RE-REVIEW: PASS）」；证据列含「T028 已完成（`DEV_PLAN.md:119`）」与「**真实 CLI 端到端仍依赖 T027-B**（D-18）」；另门表行 52 内含「**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过」 | 44、52 | AT-24 行：状态列 = 「Conditionally closed (2026-11-11, P64 acceptance; ⑦ re-review `RE-REVIEW: PASS`)」；证据列含「T028 is complete (`DEV_PLAN_EN.md:119`)」与「**the real CLI end-to-end path still depends on T027-B** (D-18)」；行 52 内含「**S-G2 caliber (user ruling ③)**: AT-24 conditionally closed (P64) ⇒ the v54 snapshot may mark "AT-24 condition closure can proceed"; the S-G2 gate-level verdict is pending the v54 snapshot and must not be cited as that gate passing」 | 44、52 |
| S-G3 | 门表行 56：门名「正式安装包教程」、状态列 = `PARTIAL`、证据列 = 「[便携 ZIP 安装指南](../INSTALLATION.md)与模型已冻结；最终 ZIP 尚未绑定和复验」；另风险区行 66 含「正式 ZIP 文件名和下载入口仍待发行绑定」 | 56、66 | 门表行 56：门名「Formal install-package tutorial」、状态列 = `PARTIAL`、证据列 = 「The [portable ZIP installation guide](../INSTALLATION_EN.md) and model are frozen; the final ZIP is not yet bound and rerun」；行 66 含「The formal ZIP filename and download entry point still require release binding」 | 56、66 |
| S-G4 | 无专行；相关 = 行 7 含「EOL 通知方式尚未冻结」；行 66 含「EOL 通知与例外已有 ADR-010 待审批方案，但尚未获产品所有者批准」 | 7、66 | 无专行；相关 = 行 7 含「the EOL notification channel remains open」；行 66 含「ADR-010 now proposes the EOL notice and exception policy, but product-owner approval remains outstanding」 | 7、66 |
| S-G5 | 门表行 60：门名「产品所有者最终 S1 批准」、状态列 = `BLOCKED`、证据列 = 「本轮只授权 whois 夹具读取，不等于 S1 发布/阶段批准」 | 60 | 门表行 60：门名「Final product-owner S1 approval」、状态列 = `BLOCKED`、证据列 = 「This session's fixture-read authorization is not release or stage approval」 | 60 |
| S-G6 | 无专行；相关 = 行 3 含「AT-23、T018 与 S1 仍阻断」；行 69 末句「在此之前 T018 与 S1 状态保持未完成」 | 3、69 | 无专行；相关 = 行 3 含「AT-23, T018 and S1 remain blocked」；行 69 末句「Until then, T018 and S1 remain incomplete」 | 3、69 |
| S-G7 | 门表整体（行 50–60，共 9 行：52 = BLOCKED；53/54/55/59 = PASS；56/57/58 = PARTIAL ×3；60 = BLOCKED）；行 62（P71 说明）含「**本文件门表 9 行的行名、行数与状态判定全部不变**」 | 50–60、62 | 门表整体（行 50–60，同构）；行 62 含「**The gate table's nine rows keep their names, count and verdicts unchanged**」 | 50–60、62 |
| S-G8 | 无专行；相关 = 行 44 含「**真实 CLI 端到端仍依赖 T027-B**（D-18）」；行 43 含「真实候选腿实际跑过启动 / 执行 / 停机并记录 unsupported，不构成真实候选兼容性证明」；行 52 含「T027 未完成、AT-23 未通过」；行 69 含「（真实 CLI 端到端仍依赖 T027-B）」 | 43、44、52、69 | 无专行；相关 = 行 44 含「**the real CLI end-to-end path still depends on T027-B** (D-18)」；行 43 含「the real candidate legs actually ran start / execution / stop and recorded unsupported; this does not constitute real-candidate compatibility proof」；行 52 含「T027 is incomplete and AT-23 is not passed」；行 69 含「(the real CLI end-to-end path still depends on T027-B)」 | 43、44、52、69 |

### 一·附 A — `s1-exit.md` 被引行全文（逐字；围栏内行首「N: 」为标注前缀，其后为逐字原文）

```text
3: > 2026-09-11 当前状态：T026 已完成，矩阵为 **10 verified、2 unsupported、0 unknown、blocked**。Job Object 取消通过；工具别名和 shell 网络可绕过 CLI deny，固定候选不兼容。十次追加授权使用 7 次，确认 6 premium requests 加取消 unknown，3 次未使用。AT-23、T018 与 S1 仍阻断；T027 须等待外部强制边界或兼容候选。详见 [T026 最终结论](t026-agent-runner-contract.md)。
7: 日期：2026-09-09。结论：`BLOCKED`。T018 已完成 AT-15 的只读 whois 冻结夹具对比和 generic/C/Go 最小闭环，可信 Windows artifact 的候选安装生命周期演练也已通过；手工解压便携 ZIP、不修改 PATH 的安装模型已获产品所有者批准，并已有双语教程、支持矩阵与发行说明草案。但最终 ZIP 尚未绑定和复验，EOL 通知方式尚未冻结，产品所有者也未批准最终 S1 exit，因此不得把部分通过标记为 S1 完成。2026-09-11 更新：T026 能力表征完成，但候选不兼容，AT-23 与任务链门禁继续阻断。
43: | AT-23 | 阶段 1 条件闭合达成（2026-11-09，F2 修复片收口；阶段 2 不适用，用户 2026-11-09 裁定 A1）：⑦ 终审 `FINAL: PASS WITH FIXES` + 盲审四问均成立 ⇒ 前提①满足 ⇒ **AT-23 阶段 1 条件闭合达成**（前提② P61 已满足）；P63 豁免判定与 P61 PARTIAL 为历史时点结论、原文保留 | [F2 修复片结论](at23-restricted-candidate.md) §11（⑦ 判定原文逐字 + F2 三腿与 B4 基线逐字段对照：verdict 四态一致、`process-tree-taken-away` 回 `true`、`failedChecks` 空集、字段漂移仅 1 处 = f2-pos artifacts 新增 stub-child.pid 键（修复预期产物））；[P61 报告](at23-restricted-candidate.md)：机制腿与 B4 基线 0 漂移（仅 F2 一处断言继承漂移）、真实候选腿实际跑过启动 / 执行 / 停机并记录 unsupported，不构成真实候选兼容性证明；[T026 报告](t026-agent-runner-contract.md) 的历史结论不变 |
44: | AT-24 | 条件闭合（2026-11-11，P64 验收；⑦ 复审 RE-REVIEW: PASS） | [T028 黑箱模式](t028-supervised-black-box.md)：T028 已完成（`DEV_PLAN.md:119`），visible 黑箱候选的风险确认、隔离投递、显式归还、全量扇描与 reduced/unknown 报告已实现；闭合依据为 (b)/(c) 裁定，P64 按该口径验收并判条件闭合（[P64 报告](at24-supervised-black-box.md)）；**真实 CLI 端到端仍依赖 T027-B**（D-18）；黑箱子项阶段归属 S2（P68 b1），条件闭合不撤销 |
50: | 门禁 | 状态 | 证据/缺口 |
51: |---|---|---|
52: （行 52 全文见第五节 5.1）
53: | seed 自托管与真实平台 | PASS | T017 Windows/Linux 可信矩阵及 Windows 本机两代演练 |
54: | whois 只读影子 | PASS | 9/9 固定用例与 3 类漂移负例；不含生产切换 |
55: | generic/C/Go harness | PASS | 三套 harness 定义及三个最小用户样例均通过只读 CLI 闭环 |
56: | 正式安装包教程 | PARTIAL | [便携 ZIP 安装指南](../INSTALLATION.md)与模型已冻结；最终 ZIP 尚未绑定和复验 |
57: | 安装/升级/回滚/卸载实测 | PARTIAL | [Windows 候选便携包演练](s1-install-drill.md)通过；最终载体确定后仍须复验 |
58: | 支持矩阵与发行说明 | PARTIAL | [候选支持矩阵](../S1_SUPPORT_MATRIX.md)与[发行说明草案](../S1_RELEASE_NOTES.md)已建立，尚未绑定正式发行物 |
59: | 独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审 |
60: | 产品所有者最终 S1 批准 | BLOCKED | 本轮只授权 whois 夹具读取，不等于 S1 发布/阶段批准 |
62: **P71 适用性说明（`S1-SCOPE-ANNOTATE`，2026-11-13）**：R-07 b1 附条件 2 落地只改变 REQ 阶段归属——完整一体化 TUI（REQ-009）标「S1 不要求，推迟至 S2（基础终端保留）」，黑箱候选子项（REQ-007/021）标推迟 S2，自托管（REQ-019 / AT-13）仍在 S1；MVA 作为验收复合标准不批正式定义。**本文件门表 9 行的行名、行数与状态判定全部不变**（完整 TUI 移出不构成门级判定变更）；visible 降级保证口径不变。`S-G1`…`S-G8` 编号体系定义于 `docs/validation/p57-closeout-assessment.md` §4.2，本文件不引入该编号，亦不建立其与门表行名的映射。既有 PASS 证据不受影响：AT-12 / AT-16 / AT-22 的断言全部基于基础终端与 `interactions tui`，不依赖完整一体化 TUI。
66: 1. 手工解压便携 ZIP、用户自选独立目录、显式路径运行且不修改 PATH 的模型已冻结；正式 ZIP 文件名和下载入口仍待发行绑定。EOL 通知与例外已有 ADR-010 待审批方案，但尚未获产品所有者批准。
69: 4. T025/AT-22 已完成聚焦终端人工交互闭环；T026 已以 10 verified、2 unsupported、0 unknown 完成能力表征，固定候选不兼容。须更换/升级并重新探测候选，或建立覆盖工具别名与 shell 网络绕过的可证明外部强制边界，再开始 T027/AT-23。T028 已完成、AT-24 在 (b)/(c) 裁定下条件闭合（真实 CLI 端到端仍依赖 T027-B），其后仍需完成完整统一 TUI；再用最终 ZIP 复跑教程并绑定支持矩阵/发行说明，最后由产品所有者明确批准 S1 exit。在此之前 T018 与 S1 状态保持未完成。
```

### 一·附 B — `s1-exit_EN.md` 被引行全文（逐字；行首「N: 」为标注前缀）

```text
3: > Current status, 2026-09-11: T026 is complete with **ten verified, two unsupported, zero unknown, blocked**. Job Object cancellation passed; tool aliases and shell networking bypass CLI denials, making the pinned candidate incompatible. Seven of ten additional authorizations were used, confirming six premium requests plus unknown cancellation cost; three remain unused. AT-23, T018 and S1 remain blocked. T027 requires an external enforcement boundary or compatible candidate. See the [final T026 conclusion](t026-agent-runner-contract_EN.md).
7: Date: 2026-09-09. Verdict: `BLOCKED`. T018 passes the frozen read-only whois comparison and minimal generic/C/Go loops required by AT-15, and a candidate lifecycle drill for trusted Windows artifacts also passes. The product owner approved manual portable-ZIP extraction without PATH modification, with bilingual installation, support-matrix, and release-note drafts now present. The final ZIP is not yet bound and rerun, the EOL notification channel remains open, and the product owner has not approved final S1 exit, so partial success must not be reported as S1 completion. Updated 2026-09-11: T026 characterization is complete, but its incompatible candidate keeps AT-23 and the task-chain gate blocked.
43: | AT-23 | Stage-1 condition closure achieved (2026-11-09, F2 fix slice closed; stage 2 not applicable, user ruling A1, 2026-11-09): ⑦ final `FINAL: PASS WITH FIXES` + the blind review's four questions all hold ⇒ premise ① met ⇒ **AT-23 stage-1 condition closure achieved** (premise ② already met by P61); the P63 waiver verdict and the P61 PARTIAL are historical point-in-time conclusions and their text is preserved | [F2 fix slice conclusion](at23-restricted-candidate_EN.md) §11 (the ⑦ verdict text verbatim + the F2 three-leg vs B4 baseline field-by-field comparison: the verdict four states agree, `process-tree-taken-away` back to `true`, `failedChecks` empty, and only 1 field drift = the newly added stub-child.pid key in the f2-pos artifacts (the fix-expected product)); [P61 report](at23-restricted-candidate_EN.md): the mechanism legs show 0 drift against the B4 baseline (only one inherited F2 assertion drift), and the real candidate legs actually ran start / execution / stop and recorded unsupported; this does not constitute real-candidate compatibility proof; the historical [T026 report](t026-agent-runner-contract_EN.md) conclusion is unchanged |
44: | AT-24 | Conditionally closed (2026-11-11, P64 acceptance; ⑦ re-review `RE-REVIEW: PASS`) | [T028 black-box mode](t028-supervised-black-box_EN.md): T028 is complete (`DEV_PLAN_EN.md:119`); visible black-box risk acknowledgment, isolated delivery, explicit return, full scan and reduced/unknown reporting are implemented; closed under rulings (b)/(c), and P64 accepted it under that caliber as conditional closure ([P64 report](at24-supervised-black-box_EN.md)); **the real CLI end-to-end path still depends on T027-B** (D-18); the black-box sub-item is owned by stage S2 (P68 b1) and the conditional closure is not revoked |
50: | Gate | Status | Evidence/gap |
51: |---|---|---|
52: （行 52 全文见第五节 5.1）
53: | Seed self-host and real platforms | PASS | T017 Windows/Linux trust matrix and native Windows two-generation drill |
54: | Read-only whois shadow | PASS | 9/9 fixed cases and three drift counterexamples; no production cutover |
55: | generic/C/Go harnesses | PASS | All three harness definitions and minimal user examples pass the read-only CLI loop |
56: | Formal install-package tutorial | PARTIAL | The [portable ZIP installation guide](../INSTALLATION_EN.md) and model are frozen; the final ZIP is not yet bound and rerun |
57: | Install/upgrade/rollback/uninstall exercises | PARTIAL | The [Windows candidate portable-package drill](s1-install-drill_EN.md) passes; the final carrier still requires rerun |
58: | Support matrix and release notes | PARTIAL | The [candidate support matrix](../S1_SUPPORT_MATRIX_EN.md) and [release-note draft](../S1_RELEASE_NOTES_EN.md) exist but are not bound to a formal release |
59: | Independent review | PASS | Both T017 baseline and Node 24 replacement received Codex/DeepSeek reviews |
60: | Final product-owner S1 approval | BLOCKED | This session's fixture-read authorization is not release or stage approval |
62: **P71 applicability note (`S1-SCOPE-ANNOTATE`, 2026-11-13)**: landing R-07 b1 condition 2 changes only REQ stage ownership - the full unified TUI (REQ-009) is marked "not required by S1, deferred to S2 (the basic terminal stays)", the black-box candidate sub-items (REQ-007/021) are marked as deferred to S2, and self-hosting (REQ-019 / AT-13) stays in S1; the MVA acceptance composite is not approved as a formal definition. **The gate table's nine rows keep their names, count and verdicts unchanged** (the TUI carve-out is not a gate-level verdict change); the visible reduced-assurance wording is unchanged. The `S-G1`...`S-G8` numbering is defined in `docs/validation/p57-closeout-assessment_EN.md` Section 4.2; this file does not introduce that numbering and does not map it onto the gate-table row names. The existing PASS evidence is unaffected: the assertions of AT-12 / AT-16 / AT-22 all rest on the basic terminal and `interactions tui` and do not depend on the full unified TUI.
66: 1. The manually extracted portable-ZIP model, user-selected independent directory, explicit-path invocation, and no-PATH policy are frozen. The formal ZIP filename and download entry point still require release binding. ADR-010 now proposes the EOL notice and exception policy, but product-owner approval remains outstanding.
69: 4. T025/AT-22 completes the focused terminal operator loop. T026 characterization is complete with ten verified, two unsupported, and zero unknown capabilities; the pinned candidate is incompatible. Replace/upgrade and reprobe the candidate, or establish a provable external enforcement boundary covering tool-alias and shell-network bypasses, before T027/AT-23. T028 is complete and AT-24 is conditionally closed under rulings (b)/(c) (the real CLI end-to-end path still depends on T027-B), and the complete unified TUI still remains; then rerun the tutorial against the final ZIP, bind the support matrix/release notes, and obtain explicit product-owner S1 exit approval. Until then, T018 and S1 remain incomplete.
```

## 二、p57 §4.2 八门原始清单（逐门）

p57 §4.2 门表位于 `docs/validation/p57-closeout-assessment.md` 行 86（表头）、87（分隔）、88–95（S-G1–S-G8 八行）；EN 侧（`p57-closeout-assessment_EN.md`）行号完全对称（86–95）。表后两条加法式指针位于行 97（P71 指针）与行 99（S-G8 P75 裁定指针）。下表摘各门「现状与证据」列逐字原文；八行全文见「二·附 A（CN）」「二·附 B（EN）」。

| 门 | 原文（逐字摘录） | 行号 |
|---|---|---|
| S-G1 | 「不可达 / 未通过：`docs/validation/s1-exit.md:43`」 | 88 |
| S-G2 | 「**口径滞后**：`docs/DEV_PLAN.md:119` 已 `COMPLETED`（AT-24 在 (b)/(c) 裁定下条件闭合），而 `docs/validation/s1-exit.md:44` 仍写 `BLOCKED` 且理由为「T028 尚未实现」」 | 89 |
| S-G3 | 「未通过：`s1-exit.md:56`（`PARTIAL`）+ `docs/ADR_REGISTER.md:19`（ADR-009）」 | 90 |
| S-G4 | 「未通过：`docs/ADR_REGISTER.md:20`（`PROPOSED / NOT APPROVED`）」 | 91 |
| S-G5 | 「未知（人工门）：`s1-exit.md:60`」 | 92 |
| S-G6 | 「不可达：`docs/DEV_PLAN.md:94`（依赖含 T027、T028）、`:96`（`IN_PROGRESS / BLOCKED`）」 | 93 |
| S-G7 | 「不可达：`s1-exit.md:52`（核心任务链 `BLOCKED`）、`:56` 至 `:58`（三项 `PARTIAL`）、`:60`」 | 94 |
| S-G8 | 「未通过：`docs/DEV_PLAN.md:119`（D-18 明写依赖 T027-B）」 | 95 |

### 二·附 A — p57 CN §4.2 门表逐行全文（逐字；行首「N: 」为标注前缀）

```text
86: | 门 | 判据（可判真假） | 现状与证据 | 依赖 | 解锁条件 |
87: |---|---|---|---|---|
88: | S-G1 AT-23 真实候选腿与整体 | CLI 可真实启动经准入的 AI 候选并计费 | 不可达 / 未通过：`docs/validation/s1-exit.md:43` | R-08、R-09、T027 装配、兼容候选 | 三者齐备后重跑真实 E2E |
89: | S-G2 AT-24 与 T028 | `supervised-black-box` 实现完成且 AT-24 判定成立 | **口径滞后**：`docs/DEV_PLAN.md:119` 已 `COMPLETED`（AT-24 在 (b)/(c) 裁定下条件闭合），而 `docs/validation/s1-exit.md:44` 仍写 `BLOCKED` 且理由为「T028 尚未实现」 | 裁定 (b)/(c) 的落地表述 | 在 S1 exit 报告中按 (b)/(c) 重述为「条件闭合」，并保留「真实 CLI 端到端依赖 T027-B」的边界 |
90: | S-G3 最终 ZIP 绑定与复验 | 最终 ZIP 绑定发行记录并复验 | 未通过：`s1-exit.md:56`（`PARTIAL`）+ `docs/ADR_REGISTER.md:19`（ADR-009） | S-G2、S-G4、发行记录 | 其余门闭合后执行绑定与复验 |
91: | S-G4 EOL 通知方式 | ADR-010 获产品所有者批准 | 未通过：`docs/ADR_REGISTER.md:20`（`PROPOSED / NOT APPROVED`） | 产品所有者 | 批准 ADR-010 |
92: | S-G5 产品所有者 S1 exit 批准 | 产品所有者作出批准 | 未知（人工门）：`s1-exit.md:60` | S-G1 至 S-G4 | 由产品所有者作出 |
93: | S-G6 T018 整体 | 依赖字段满足且四道剩余门闭合 | 不可达：`docs/DEV_PLAN.md:94`（依赖含 T027、T028）、`:96`（`IN_PROGRESS / BLOCKED`） | S-G1 与 S-G2 与任务级依赖 | 门级与任务级**不同源**：四门清空后 T018 仍须 T027、T028 完成 |
94: | S-G7 S1 exit 总体（RFC §14） | 全 AT 通过且 T018 完成且 §14 全 exit | 不可达：`s1-exit.md:52`（核心任务链 `BLOCKED`）、`:56` 至 `:58`（三项 `PARTIAL`）、`:60` | S-G1 至 S-G6 | 上列全绿 |
95: | S-G8 真实 CLI 端到端 | 受监督路径可在真实 CLI 与真实候选下端到端运行 | 未通过：`docs/DEV_PLAN.md:119`（D-18 明写依赖 T027-B） | P5 的 R-3 | 与 S-G1 同源解锁 |
97: - **P71 指针（2026-11-13）**：`docs/validation/s1-exit{,_EN}.md` 已由切片 `S1-SCOPE-ANNOTATE` 追加「P71 适用性说明」段（REQ 阶段归属变更不影响门级判定；S-Gx 编号仍以本节为定义源）；报告见 `docs/validation/s1-scope-annotate{,_EN}.md`。本表八门定义与判定不改（§11.1）。
99: - **S-G8 判定已由 P75 裁定**：保持 blocked 直到 S2 enforcement 接入（报告：`docs/validation/r3-real-cli-assessment{,_EN}.md`）。本表八门定义与判定不改（§11.1）。
```

### 二·附 B — p57 EN §4.2 门表逐行全文（逐字；行首「N: 」为标注前缀）

```text
86: | Gate | Criterion (truth-apt) | State and evidence | Depends on | Unblock condition |
87: |---|---|---|---|---|
88: | S-G1 AT-23 real-candidate leg and overall | the CLI can genuinely launch an admitted AI candidate and bill it | unreachable / not passed: `docs/validation/s1-exit.md:43` | R-08, R-09, T027 assembly, a compatible candidate | rerun the real E2E once all four exist |
89: | S-G2 AT-24 and T028 | `supervised-black-box` implemented and the AT-24 verdict holds | **stale wording**: `docs/DEV_PLAN.md:119` already says `COMPLETED` (AT-24 conditionally closed under rulings (b)/(c)), while `docs/validation/s1-exit.md:44` still says `BLOCKED` with the reason "T028 not implemented" | landing of rulings (b)/(c) as wording | restate as "conditionally closed" per (b)/(c) in the S1 exit report, keeping the "real CLI end-to-end depends on T027-B" boundary |
90: | S-G3 Final ZIP binding and re-verification | the final ZIP is bound to the release record and re-verified | not passed: `s1-exit.md:56` (`PARTIAL`) + `docs/ADR_REGISTER.md:19` (ADR-009) | S-G2, S-G4, release record | bind and re-verify once the other gates close |
91: | S-G4 EOL notification method | ADR-010 approved by the product owner | not passed: `docs/ADR_REGISTER.md:20` (`PROPOSED / NOT APPROVED`) | product owner | approve ADR-010 |
92: | S-G5 Product owner S1 exit approval | the product owner approves | unknown (human gate): `s1-exit.md:60` | S-G1 through S-G4 | granted by the product owner |
93: | S-G6 T018 as a whole | dependency fields satisfied and the four remaining gates closed | unreachable: `docs/DEV_PLAN.md:94` (depends on T027 and T028), `:96` (`IN_PROGRESS / BLOCKED`) | S-G1 and S-G2 and the task-level dependencies | gate level and task level are **not the same source**: even with four gates cleared, T018 still needs T027 and T028 done |
94: | S-G7 S1 exit overall (RFC §14) | all ATs pass, T018 complete, §14 fully exited | unreachable: `s1-exit.md:52` (core task chain `BLOCKED`), `:56` to `:58` (three `PARTIAL`), `:60` | S-G1 through S-G6 | all of the above green |
95: | S-G8 Real CLI end-to-end | the supervised path runs end-to-end on a real CLI and a real candidate | not passed: `docs/DEV_PLAN.md:119` (D-18 states it depends on T027-B) | P5 residue R-3 | unlocked together with S-G1 |
97: - **P71 pointer (2026-11-13)**: `docs/validation/s1-exit{,_EN}.md` now carries a "P71 applicability note" added by the `S1-SCOPE-ANNOTATE` slice (the REQ stage-ownership change does not affect gate-level verdicts; this section remains the definition source for the S-Gx numbering); the report is `docs/validation/s1-scope-annotate{,_EN}.md`. The eight-gate definitions and verdicts in this section stay unchanged (§11.1).
99: - **S-G8 verdict ruled by P75**: stays blocked until S2 enforcement lands (report: `docs/validation/r3-real-cli-assessment{,_EN}.md`). The eight-gate definitions and verdicts in this table stay unchanged (§11.1).
```

## 三、p57 §4.3 原编号 vs s1-exit 现行编号对照

p57 §4.3 位于行 101（标题）、103–105（三条）、107（收束段）；EN 侧行号完全对称（101–107）。P75 §11.2 的对应登记位于 `docs/validation/r3-real-cli-assessment.md` §11.2（行 158–166）。

| p57 §4.3 编号 | s1-exit 现行编号 | 说明 |
|---|---|---|
| 原第 1 项（行 103：`1. \`s1-exit.md:44\` 的 AT-24 条目称「T028 尚未实现」，与 \`DEV_PLAN.md:119\` 的 \`COMPLETED\` 互斥。」） | `s1-exit.md:44`（AT-24 行；现行行号未位移，仍在 :44） | 现行磁盘 :44 状态列 = 「条件闭合（2026-11-11，P64 验收；⑦ 复审 RE-REVIEW: PASS）」、证据列含「T028 已完成（`DEV_PLAN.md:119`）」；「尚未实现」在现行 s1-exit.md 全文 0 命中（计数核验）。P75 §11.2 对应项 = 第 1 项，处置 = 「已由 P64 更正 ⇒ 关闭」（r3-real-cli-assessment.md:164） |
| 原第 2 项（行 104：`2. \`s1-exit.md:52\` 的核心任务链行称「T027/T028 …… 仍未实现」，T028 部分同样滞后。」） | `s1-exit.md:52`（核心任务链行；现行行号未位移，仍在 :52） | 现行磁盘 :52 状态列 = `BLOCKED`；行内相应表述 = 「T027 未完成、AT-23 未通过」；「仍未实现」在现行 s1-exit.md 全文 0 命中（计数核验）。P75 §11.2 对应项 = 第 3 项，处置 = 「唯一残余 ⇒ 移交 S1 exit 收口片」（r3-real-cli-assessment.md:165） |
| 原第 3 项（行 105：`3. \`docs/validation/s1-reach-reachability.md\` §2 的 G-D 行把 T028 记为 \`PLANNED / NOT IMPLEMENTED\`，且其 \`DEV_PLAN.md:91\` / \`:93\` 引用已漂移为 \`:94\` / \`:96\`。」） | 非 s1-exit 文件：`docs/validation/s1-reach-reachability.md` §2 的 G-D 行（本批未实读该文件，其现行内容不在本批事实范围内） | 该项引用的 `DEV_PLAN.md:94`（T018 任务行，行内含「依赖：T017、T025、T026、T027、T028 及 whois 输入授权」）与 `:96`（「T018 状态：`IN_PROGRESS / BLOCKED`」）均为现行磁盘实锚。P75 §11.2 对应项 = 第 2 项，处置 = 「已由 s1-reach 更正 ⇒ 关闭」（r3-real-cli-assessment.md:166） |

结构事实（编号错位登记）：r3-real-cli-assessment.md §11.2 表（行 162 表头、164–166 三行）的行序为「第 1 项（:164）→ 第 3 项（:165）→ 第 2 项（:166）」，即 P75 §11.2 自身编号与 P57 §4.3 原编号在第 2/3 项之间存在互换；`docs/t027/CHECKLIST_v59.md:41` 明写「**P75 §11.2 第 1 / 2 项**（对应 P57 §4.3 原第 1 / 3 项）已由 P64 / s1-reach 更正关闭、**P75 §11.2 第 3 项**（对应 P57 §4.3 原第 2 项 = `s1-exit.md:52`）为唯一残余」。

## 四、与 P75 闭环后状态的事实差异（仅陈述事实，不判定）

P75 闭环后状态的参照点（均为本批实锚）：`docs/t027/CHECKLIST_v59.md:1`（标题「任务执行清单（跨会话快照 v59 · P75 闭环后 · S1 Exit 收口评估待起跑）」）、`:3`（生成时点 2026-11-16）、`:41`；`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md:23`（「P75 已闭环：R-3 = (A) 暂不可起跑；T027-B = 完成（部分不适用）；S-G8 = blocked 至 S2；AT-23 阶段 1 条件闭合 = 永久事实；ADR-045 已登记」）、`:27`；`docs/validation/r3-real-cli-assessment.md:158–166`；`docs/ADR_REGISTER.md:56`（ADR-045）；`docs/DEV_PLAN.md:203`。

| 门 | 差异事实 | 依据 file:line |
|---|---|---|
| 文档级-1 | p57 §3.3 现存三条加法式收口指针（逐字）：「**R-1 已由 P69 `T027B-E2E-WIRING` 收口**（报告：`docs/validation/t027b-e2e-wiring{,_EN}.md`）」（:72）、「**R-2 已由 P73 裁定（2026-11-14）**：OB-55 归属 = 方案 b（`Sandbox` 接口提前到 S1，实施片 `P74-SANDBOX-INTERFACE` 待起跑）；报告：`docs/validation/t027b-admission{,_EN}.md`、ADR-044。本表八门定义与判定不改（§11.1）。」（:74）、「**R-3 已由 P75 裁定 = 暂不可起跑**（报告：`docs/validation/r3-real-cli-assessment{,_EN}.md`）；本表八门定义与判定不改（§11.1）。」（:76）；§4.2 表后另有「**S-G8 判定已由 P75 裁定**：保持 blocked 直到 S2 enforcement 接入（报告：`docs/validation/r3-real-cli-assessment{,_EN}.md`）。本表八门定义与判定不改（§11.1）。」（:99）。只读 git log 显示 p57-closeout-assessment.md 最近一次被提交 `ed2e26e docs: rule R-3 as not startable and declare T027-B complete in part` 触及（即 P75 主提交） | p57-closeout-assessment.md:72/:74/:76/:99；git log（只读） |
| 文档级-2 | s1-exit{,_EN}.md 全文无「P75」字样（0 命中，双语同）；只读 git log 显示 s1-exit.md 最近一次被提交 `6ee2c61 docs: annotate the S1 scope stage ownership and land the visible debugging guidance (P71)` 触及（即 P71 提交），其后无提交触及 ⇒ P75 提交集未触及 s1-exit.md | s1-exit.md / s1-exit_EN.md 全文计数；git log（只读） |
| S-G1 | p57 §4.2 S-G1 行现状栏写「不可达 / 未通过：`docs/validation/s1-exit.md:43`」；现行磁盘 s1-exit.md:43 状态列 = 「阶段 1 条件闭合达成（2026-11-09，F2 修复片收口…）」；且 s1-exit.md:52 行内并存两段表述：「T027 未完成、AT-23 未通过」与「⇒ **S-G1 依据 = AT-23 阶段 1 条件闭合**（…依据本片 ⑦ 判定，而非 S1 Exit 整体通过）」 | p57-closeout-assessment.md:88；s1-exit.md:43/:52 |
| S-G2 | p57 §4.2 S-G2 行现状栏写「…而 `docs/validation/s1-exit.md:44` 仍写 `BLOCKED` 且理由为『T028 尚未实现』」；现行磁盘 s1-exit.md:44 状态列 = 「条件闭合（2026-11-11，P64 验收…）」、证据列含「T028 已完成（`DEV_PLAN.md:119`）」；「尚未实现」/「not implemented」在现行 s1-exit{,_EN}.md 均 0 命中 | p57-closeout-assessment.md:89；s1-exit.md:44；计数核验 |
| S-G3 | 锚点一致：现行 s1-exit.md:56 状态列 = `PARTIAL`；现行 ADR_REGISTER.md:19 = ADR-009 行 | p57-closeout-assessment.md:90；s1-exit.md:56；ADR_REGISTER.md:19 |
| S-G4 | 锚点一致：现行 ADR_REGISTER.md:20 = ADR-010 行且该行含「PROPOSED / NOT APPROVED」（全文唯一命中行 = 20） | p57-closeout-assessment.md:91；ADR_REGISTER.md:20 |
| S-G5 | 锚点一致：现行 s1-exit.md:60 状态列 = `BLOCKED` | p57-closeout-assessment.md:92；s1-exit.md:60 |
| S-G6 | 锚点一致：现行 DEV_PLAN.md:94 = T018 任务行（行内含「依赖：T017、T025、T026、T027、T028 及 whois 输入授权」）、:96 = 「T018 状态：`IN_PROGRESS / BLOCKED`」 | p57-closeout-assessment.md:93；DEV_PLAN.md:94/:96 |
| S-G7 | 锚点一致：现行 s1-exit.md:52 = `BLOCKED`、:56–:58 = 三行 `PARTIAL`、:60 = `BLOCKED`；P75 §11.2 将 :52 登记为唯一残余（「唯一残余 ⇒ 移交 S1 exit 收口片」，r3:165）；CHECKLIST_v59.md:41 同口径；P76 briefing:27 明写「**P76 的 ②b 阶段将修正 `s1-exit.md:52` 核心任务链行内口径不一致——这是『宣告 S1 Exit 通过』的准入条件，不是 P76 起跑的准入条件。**」 | p57-closeout-assessment.md:94；s1-exit.md:52/:56–:58/:60；r3-real-cli-assessment.md:165；CHECKLIST_v59.md:41；briefing:27 |
| S-G8 | p57 §4.2 S-G8 行现状栏引 `docs/DEV_PLAN.md:119`；现行磁盘 :119 = 「- [x] T028 …」（T028 任务行带完成勾选）；P75 裁定 = 「保持 blocked 直到 S2 enforcement 接入」已加法式登记于 p57:99 与 ADR-045（ADR_REGISTER.md:56）；DEV_PLAN.md:203 记「T027-B 收口裁定（P75 `R3-REAL-CLI-ASSESSMENT`，2026-11-15）：T027-B 宣告『**完成（部分不适用）**』…**R-3 已由 P75 裁定 = 暂不可起跑**」 | p57-closeout-assessment.md:95/:99；DEV_PLAN.md:119/:203；ADR_REGISTER.md:56 |

## 五、其它显著事实

### 5.1 s1-exit.md 第 52 行逐字原文（CN 与 EN，各为文件内单行）

CN（`docs/validation/s1-exit.md` 行 52；行首「52: 」为标注前缀）：

```text
52: | 核心任务链与 fail-close 场景 | BLOCKED | 既有 AT-01–AT-22 通过；T026 完成但固定候选为 10 verified、2 unsupported、0 unknown，整体不兼容；T027 未完成、AT-23 未通过；T028 已完成、AT-24 已验收 = 条件闭合（P64，2026-11-11；[P64 报告](at24-supervised-black-box.md)）；**依赖链**：AT-23 阶段 1 条件闭合达成（前提① F2 修复片满足、前提② P61 已满足；阶段 2 不适用，用户 2026-11-09 裁定 A1）⇒ **S-G1 依据 = AT-23 阶段 1 条件闭合**（阶段 2 不适用，用户 2026-11-09 裁定 A1；依据本片 ⑦ 判定，而非 S1 Exit 整体通过）；S-G7 依赖链同步：依赖链中 AT-23 部分已阶段 1 条件闭合（阶段 2 不适用），其余门（S-G2~S-G6、S-G8）仍待，**不得声称 S1 Exit 整体通过**；**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过 |
```

第 52 行状态列（第 2 列）逐字 = `BLOCKED`。

EN（`docs/validation/s1-exit_EN.md` 行 52；行首「52: 」为标注前缀）：

```text
52: | Core task chain and fail-closed scenarios | BLOCKED | Existing AT-01–AT-22 pass; T026 is complete, but its pinned candidate has ten verified, two unsupported, and zero unknown capabilities and is incompatible; T027 is incomplete and AT-23 is not passed, while T028 is complete and AT-24 is accepted = conditional closure (P64, 2026-11-11; [P64 report](at24-supervised-black-box_EN.md)); **dependency chain**: AT-23 stage-1 condition closure achieved (premise ① met by the F2 fix slice, premise ② already met by P61; stage 2 not applicable, user ruling A1, 2026-11-09) ⇒ **S-G1 basis = AT-23 stage-1 condition closure** (stage 2 not applicable, user ruling A1, 2026-11-09; per this slice's ⑦ verdict, not S1 Exit passing as a whole); S-G7 dependency-chain sync: the AT-23 part of the chain has reached stage-1 condition closure (stage 2 not applicable) and the remaining gates (S-G2~S-G6, S-G8) are still pending, and **S1 Exit as a whole is not claimed to pass**; **S-G2 caliber (user ruling ③)**: AT-24 conditionally closed (P64) ⇒ the v54 snapshot may mark "AT-24 condition closure can proceed"; the S-G2 gate-level verdict is pending the v54 snapshot and must not be cited as that gate passing |
```

第 52 行状态列（第 2 列）逐字 = `BLOCKED`。

### 5.2 s1-exit{,_EN}.md 内 BLOCKED / PARTIAL / DONE 等状态词行枚举

CN（`s1-exit.md`，大小写按原文逐字；引文均见附 A / 5.1）：

| file:line | 状态词出现形态（逐字） |
|---|---|
| s1-exit.md:3 | 小写「blocked」，位于「矩阵为 **10 verified、2 unsupported、0 unknown、blocked**」 |
| s1-exit.md:7 | 「`BLOCKED`」，位于「结论：`BLOCKED`」 |
| s1-exit.md:43 | 「PARTIAL」，位于「P63 豁免判定与 P61 PARTIAL 为历史时点结论」（属历史结论引用，非门状态列） |
| s1-exit.md:52 | 状态列 = 「BLOCKED」（全行见 5.1） |
| s1-exit.md:56 | 状态列 = 「PARTIAL」（正式安装包教程行） |
| s1-exit.md:57 | 状态列 = 「PARTIAL」（安装/升级/回滚/卸载实测行） |
| s1-exit.md:58 | 状态列 = 「PARTIAL」（支持矩阵与发行说明行） |
| s1-exit.md:60 | 状态列 = 「BLOCKED」（产品所有者最终 S1 批准行） |

「DONE」在 s1-exit.md 全文 0 命中（计数核验）。

EN（`s1-exit_EN.md`，行号与 CN 完全对称）：行 3（「zero unknown, blocked**」）、行 7（「Verdict: `BLOCKED`」）、行 43（「the P61 PARTIAL are historical point-in-time conclusions」）、行 52（状态列 = `BLOCKED`，全行见 5.1）、行 56/57/58（状态列 = `PARTIAL` ×3）、行 60（状态列 = `BLOCKED`）；「DONE」0 命中（计数核验）。

补充枚举（「PASS」字样，双语命中行完全对称）：AT 汇总表行 21–42（22 行，状态列均为 `PASS`）；行 43/44 状态列非 `PASS` 但行内含「`FINAL: PASS WITH FIXES`」/「`RE-REVIEW: PASS`」字样；行 46（「这些 PASS 表示…」/「These PASS verdicts…」）；门表行 53、54、55、59（状态列 = `PASS`）；行 62（「既有 PASS 证据不受影响」/「The existing PASS evidence is unaffected」）。完整命中行清单：21–44、46、53、54、55、59、62（CN 与 EN 完全一致）。

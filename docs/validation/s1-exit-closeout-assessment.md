# P76 S1 Exit 收口评估（`S1-EXIT-CLOSEOUT-ASSESSMENT`）
日期：2026-11-16。切片标识：`S1-EXIT-CLOSEOUT-ASSESSMENT`（P76，`[SLICE]` 评估类 M 档）。
[English](s1-exit-closeout-assessment_EN.md)

**行号口径**：本报告所有 `file:line` 取自 2026-11-16 磁盘实读；行号漂移时以「文件名 + 节号」为主锚点。

---

## 1. 变更概要

- 本片做了什么：在 P75 闭环后的 S1 Exit 八门现状上出具收口路线图，按用户 2026-11-16 裁定落地六项裁定（Q2 / Q3 / Q7 / Q9 / Q10 与收口顺序）、两项 ADR、三处加法式注记与一处门表行内口径修文。
- 本片没做什么：不执行真实 CLI 端到端、不重跑 AT-23 / AT-24、不改 `s1-reach-reachability{,_EN}`、不改准则版本（维持 v1.69）、不改 v50–v59 任何快照。
- 一句话结论：S1 内可推进的四门已给出可执行收口路径；S1 Exit 整体在 S1 内不可能通过（S-G1 判据未满足、S-G8 blocked 至 S2、S-G7 不可达）。
- 执行自由度清单：`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md`。

## 2. 事实包

| 项 | 值 |
|---|---|
| 切片 / 级别 | `S1-EXIT-CLOSEOUT-ASSESSMENT` / `[SLICE]` 评估类 M 档 |
| 主控 | `DeepSeek V4.1 Flash` |
| 步骤 0 基线 | `HEAD` = `origin/main` = `b90de21`；准则 v1.69；`SELFTEST: PASS 737/737`；index / tree / ci 三作用域 `TOTAL_FAIL=0` |
| 滚动窗口 | v57 / v58 / v59 按 LF 计数为 160 / 148 / 180，CN 与 EN 各自等长 |
| 用户裁定 | 2026-11-16 六项裁定 + 附加项 + X 授权（见 §15） |
| 证据包 | `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/` |
| 机械读数 | ADR 最大号 = 47（本片新增 ADR-046 / ADR-047）；OB 最大号 = 139；DR 最大号 = 14；试点顺延第 25 次 |

## 3. Q1：八门与十门映射

| S1 Exit 门 | s1-reach 门 | 映射依据 | 判定差异 |
|---|---|---|---|
| S-G1 | G-A + G-B + G-C | `p57:88` 对 `s1-reach:22`–`:24` | 一对三；G-B 机制腿单独成门，八门无此粒度 |
| S-G2 | G-D | `p57:89` 对 `s1-reach:25` | 一致（G-D 状态含真实 CLI 端到端折入） |
| S-G3 | G-E | `p57:90` 对 `s1-reach:26` | 一致 |
| S-G4 | G-F | `p57:91` 对 `s1-reach:27` | 一致 |
| S-G5 | G-G | `p57:92` 对 `s1-reach:28` | 状态词不同：`未知` 与 `BLOCKED` 语义一致 |
| S-G6 | G-H | `p57:93` 对 `s1-reach:29` | 一致；均声明门级与任务级不同源 |
| S-G7 | G-I | `p57:94` 对 `s1-reach:30` | 一致；G-I 判据字面更完整 |
| S-G8 | 无对应门 | s1-reach 无独立表达 | 语义折入 G-A / G-C 解锁条件与 G-D；属覆盖缺口，非矛盾 |
| （无） | OB-21 | 跨片观察项 | 不属八门；不得并入任何 S-Gx 判定 |

判定闭集不同：八门侧混用 `PASS` / `PARTIAL` / `BLOCKED` / 不可达 / 未通过 / 未知，十门侧闭集为不可达 / 未通过 / 已执行不独立成立 / 未知；映射时不得跨闭集换词。

## 4. Q2：S-G1 / S-G2 门级判定

- 口径 = `(a) 判据字面`（用户裁定）：门级判定严格按 `p57` §4.2 判据字面，AT 层成就不自动折算为门级通过。
- S-G1 = 不通过：真实候选腿不可达（`r3-real-cli-assessment.md` 的 §9 事实链；`s1-exit.md:43` 明写真实候选腿不构成真实候选兼容性证明）。
- S-G2 = 通过：判据「`supervised-black-box` 实现完成且 AT-24 判定成立」已满足（T028 `COMPLETED` 见 `DEV_PLAN.md:119`；AT-24 条件闭合见 `s1-exit.md:44`）。
- 与十门的一致性：S-G1 与 G-A「不可达」同向；S-G2 与 G-D 的「已执行不独立成立」不冲突（G-D 把真实 CLI 端到端折入解锁条件，本口径把它归 S-G8）。
- 已登记残余：`s1-exit.md:52` 中的「S-G2 门级判定待 v54 快照」句属本片边界外的历史表述，本片未改该句；其同步修正建议以后续切片 `S1-EXIT-SG2-SYNC`（建议名，是否立项由用户裁定）承载，不影响本片门级判定的依据。

## 5. Q3：`s1-exit.md:52` 修文

- 实读结论：该行状态列为 `BLOCKED`，前半写「T027 未完成、AT-23 未通过」，后半写「AT-23 阶段 1 条件闭合达成」，两者在同一单元格内并列而不易区分。
- 所选方案：方案 A（用户 X 授权下的默认方案）——加「整体」二字区分整体判定与阶段 1 判定；未采用方案 B，因方案 A 实读后已足以消歧（依赖链括号内已两次写明阶段 2 不适用）。
- 修文前（CN）：`T027 未完成、AT-23 未通过；`
- 修文后（CN）：`T027 未完成、AT-23 整体未通过；`
- 修文前（EN）：`T027 is incomplete and AT-23 is not passed, while`
- 修文后（EN）：`T027 is incomplete and AT-23 as a whole is not passed, while`
- 方案对照（X 授权下的两个落地文本；方案 B 为本片按 Q3 方向补拟、非用户原文，仅作对照记录）：方案 A = `AT-23 整体未通过`（加「整体」二字）；方案 B = `AT-23 整体未通过（阶段 1 条件闭合已达成，详见本行依赖链说明）`。本片采用方案 A。
- 三条依据：与本文件 `:56`–`:58` 三个 `PARTIAL` 局部门相比，核心任务链含未完成任务与未通过的 AT，属真实阻断；与 `p57` §4.3 残余项定义相比，修的是表述对齐而非状态改判；与 `s1-reach:30` 的 G-I 依据句「核心任务链 `BLOCKED`」一致。
- 边界：状态列不变、门表行名与行数不变、`:56`–`:58` 与 `:60` 不变、其余行零改动。

## 6. Q4：S-G7 依据栏

- 判定结论不变：S-G7 仍不可达（依赖 S-G1 至 S-G6，其中 S-G1 不通过、S-G8 blocked、S-G5 未批准）。
- `p57` 依据栏无需改：该行只引状态词，本次修文不改状态词。
- 已落地的同步方式：在 `p57` §4.2 表后追加一行加法式指针，指向本报告与 ADR-046。

## 7. Q5：S-G3 解锁路径

| 步 | 动作 | 前置条件 | 证据形态 |
|---|---|---|---|
| 一 | 产品所有者批准正式 Release 与版本号 | ADR-010 已批准或在同一批准动作内批准 | ADR 状态更新与发行授权记录 |
| 二 | 构建最终 ZIP 候选工件 | 步一完成；CI candidate-build 可用 | artifact 与 `SHA256SUMS`、SBOM、许可证清单 |
| 三 | 绑定发行记录（版本 / commit / run / 文件名） | 步二完成 | `INSTALLATION.md:206` 要求的三处绑定落盘 |
| 四 | 用最终 ZIP 重跑安装生命周期 | 步三完成 | 安装 / 升级 / 回滚 / 卸载复验记录 |

S-G3 依赖 S-G4：版本号在 ADR-010 内，因此 S-G4 的批准动作须先于或与步二并行；绑定与复验必须等版本号冻结。

## 8. Q6：S-G4 解锁路径

- 批准主体：产品所有者（`ADR_REGISTER` 表头明写每项批准记录主体与日期，禁止模型代签）。
- 批准后改动点：`ADR_REGISTER.md:20` 状态列、`INSTALLATION_PLAN.md:32` 的 EOL 项、`s1-exit.md:66` 的 EOL 表述、`PRODUCT_REQUIREMENTS.md:107` 的 EOL 项、支持矩阵尾部政策段。
- 边界：ADR-010 的批准不等同 S-G5；两者是不同授权对象（发行策略与阶段退出）。

## 9. Q7：S-G5 与独立安全评审

- 历史状态核实结论：从未执行。依据为 `SECURITY.md:5`「未完成独立安全签署」、`PRODUCT_REQUIREMENTS.md:107` 的评审者仍待落实、`DEV_PLAN.md:27` 的后续门禁表述，以及 `docs/validation/` 全域无独立安全评审报告。
- 对 S-G5 的影响：S-G5 判据只有「产品所有者作出批准」，故该门禁不构成 S-G5 的直接前置；但它属 `PRODUCT_REQUIREMENTS` §7 的 S0→S1 门禁历史欠账，须显式处置。
- 处置（用户裁定 A）：豁免顺延至 S2 并登记 ADR-047；未覆盖范围为威胁模型面、硬安全门禁面、秘密与凭据面；到期日为 S2 启动日与 2027-03-31 中较早者。
- 是否阻断 S1 Exit 宣告：不构成八门判据的直接阻断；但 S-G5 的分步批准材料须附本项处置说明。

## 10. Q8：S-G6 门级与任务级不同源

- 实测字段：`DEV_PLAN.md:94` 的 T018 依赖为「T017、T025、T026、T027、T028 及 whois 输入授权」；`:96` 的状态为 `IN_PROGRESS / BLOCKED`。
- 口径（用户裁定 c）：门级与任务级并陈，不给单一结论；本片零改动。
- 后果：S-G3 / S-G4 / S-G5 收口只推进门级侧；T018 因 T027 未实现而不可勾选；S-G6 门级保持不可达至 S2。

## 11. Q9：支持矩阵口径修正

| 位置 | 修正后表述 |
|---|---|
| `:13` | `R-3 = (A) 暂不可起跑；S-G8 blocked 至 S2 enforcement 接入（ADR-045）` |
| `:14` | `AT-24 条件闭合；真实 CLI 端到端依赖 S2 enforcement（R-3 = (A)，ADR-045）` |

`:5` 不动（安全设计基线状态与 P75 无关）；支持层级不变（Windows 仍为正式支持候选、Linux 仍为核心 CI），只改状态口径。

## 12. Q10：产品所有者批准时机

- 口径（用户裁定 C，分两步）：第一步在 S1 内批准 ADR-010 与 v0.1.0 Release，使 S-G3 / S-G4 收口；第二步在 S2 后批准 S1 Exit 整体。
- 第一步后 S-G5 的记账方式：记「发行部分已批准、S1 Exit 整体批准待 S2」，不得记为通过。
- 硬边界：任何批准不得被解读为 S1 Exit 通过；`s1-exit.md:52` 未修正前不得宣告；S-G8 保持 blocked。

## 13. Q11：REQ 阶段归属边界

| 项 | 归属 | 依据 |
|---|---|---|
| REQ-007 黑箱候选子项 | 推迟 S2 | `PRODUCT_REQUIREMENTS.md:69`；`s1-exit.md:62` 的 P71 注记 |
| REQ-009 完整一体化 TUI | S1 不要求 | `PRODUCT_REQUIREMENTS.md:71`；`s1-exit.md:62` |
| REQ-020 三语言实跑 | S1 仅 C + Go 契约 | `PRODUCT_REQUIREMENTS.md:82` |
| AT-13 自托管 | 仍在 S1 且 T017 已完成 | `s1-exit.md:62`；`DEV_PLAN` 的 T017 条目 |

约束：上述各项一律不得作为任何 S-Gx 的「未通过原因」；S-G1 与 S-G8 的 blocked 依据只能来自 R-3 = (A)、T027 装配未完成与候选兼容性。

## 14. 收口路线图

1. `s1-exit{,_EN}.md:52` 准入修文（本片已落地）。
2. S-G2 门级落定（本片已成文，Q2 = a）。
3. S-G4：产品所有者批准 ADR-010。
4. S-G3：正式 Release、最终 ZIP 构建、三处绑定、生命周期复验。
5. S-G6：门级侧并陈推进，任务级保持不可勾选。
6. S-G5：按 Q10 分两步批准。
7. S-G1 / S-G8 / S-G7：留 S2，S1 Exit 整体宣告待 S2 解锁后最终批准。

并行窗口：S-G4 的审批动作可与 S-G3 的 ZIP 构建准备并行；绑定与复验必须等版本号冻结。阻塞点：产品所有者批准动作与 S2 enforcement 接入，后者不在 S1 控制内。

## 15. 用户裁定记录

| 编号 | 裁定 |
|---|---|
| Q2 | `(a) 判据字面` |
| Q3 | `A 加法式（BLOCKED 保持）` |
| Q7 | `A 豁免顺延 S2 + 登记 ADR` |
| Q9 | `纳入 ②b` |
| Q10 | `(C) 分两步批准` |
| 收口顺序 | `认可 §12 序列` |
| 附加 | `p57 不改，只加指针；⑥ 可跳过` |
| X | `Q3 方案 A 优先、B 仅当 A 实读后不足时用；无歧义则停下升级` |

裁定原文全文见证据包 `02-ruling/ruling-2026-11-16.md`；生效日期 2026-11-16，裁定人为产品所有者。

## 16. CI 事实链

本片为评估类，主提交与 CI run 的对应关系由回填提交补入（§11.2）；该回填提交自身的 run 只作观察记录。

## 17. 已知边界

- 真实 enforcement 留 S2；S-G8 保持 blocked；AT-23 阶段 1 条件闭合作为永久事实；本片不重跑任何 AT；探针额度为零。
- 独立安全评审已豁免顺延 S2，未覆盖范围见 §9 与 ADR-047。
- `s1-exit.md:52` 的 S-G2 句未同步（本片仅改 AT-23 片段），已登记为残余；承载建议片 = `S1-EXIT-SG2-SYNC`（待用户立项）。

## 18. 下一步

- ADR-010 批准片（产品所有者动作）。
- Release 与 ZIP 绑定复验片。
- S2 启动前的独立安全评审补评片。

## 19. 计量与审查链

| 环节 | 次数 |
|---|---|
| ① 架构 | 1（`deep-reasoner` + `DeepSeek V4 Pro`） |
| 批量辅助 | 3（`GLM-5.3-Flash`，事实枚举） |
| ⑤ 预审 | 2 轮：r1 `FINDINGS`（1 高 + 5 低）→ r2 `PASS WITH FIXES`（0 中及以上，2 低） |
| ⑥ 独立扫描 | 由 ⑦ 确认跳过 |
| ⑦ 终审 | 1 轮：`PASS WITH FIXES`（0 中及以上，1 低；含 `SKIP CONFIRMED`） |
| 探针 | 0 |
| R / D | R = 1、D = 0 |

- ⑥ 跳过：§7.2.1 三条件成立（改动只落文档与台账类文件；不含 `.go` / `.json` / schema / fixtures / 脚本语义变更；不涉状态机 / 门禁判据 / 证据模型 / 角色授权契约语义），由 ⑦ 确认并记录理由。

## artifacts

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/s1-exit-closeout-assessment.md	repo
docs/validation/s1-exit-closeout-assessment_EN.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/MANIFEST.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/SHA256SUMS.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/01-arch/design-r1.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-batch/batch-01-s1-exit-p57.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-batch/batch-02-s1-reach.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-batch/batch-03-lagging-docs.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-ruling/ruling-2026-11-16.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/04-reviews/final-r1.md	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-index.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-tree.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-ci.txt	repo
docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/06-gates/gate-selftest.txt	repo
```

# P73 `T027B-ADMISSION` — OB-55 准入评估与裁定落地（变更报告）

切片：`T027B-ADMISSION`（P73 · 层 2 治理侧 · M 档评估类 `[SLICE]`）。主控 `DeepSeek V4.1 Flash (deepseek)`；用户 2026-11-14 同轮显式授权已内含于起跑消息。基线 `0cac174`；准则 v1.69（本片在主提交前按 §0.2 补升版，系元数据补齐）。

## 1. 执行摘要

- **OB-55 裁定 = 方案 b**（提前落平台无关 `Sandbox` 接口），与 R-06（2026-11-12）同向，属确认与细化；三条附条件全部接受。
- **落地面**：`DEV_PLAN` T027-B 段新增 P73 准入条目；B2 报告 §8 第 8 条加法式指针；OB-55 处置列翻「已处置」；新增 ADR-044；三处加法式落地指针；台账 P73 段收口。
- **边界**：不产生产代码，不改 `CONTRACTS` / schema / fixtures / CI，不改 S1 Exit 八门判定本身；实施片 `P74-SANDBOX-INTERFACE` 待起跑，本片不启动。

## 2. 现状核对

| 项 | 事实 |
|---|---|
| C1–C7 与 N1–N4 | P38 施工图给出判定（C1 / C4 / C6 = BLOCKED，C2 / C3 / C5 / C7 = PARTIAL，另有 N1–N4）；落地由 P39–P44 承载，P38 原文未回写落地态。 |
| P39–P44 | 依次落地 C6 / C1 / C2+N1 / C3+N2 / C4+C5+N3 / C7+N4；P44 自限「不得宣称端到端闭合」。 |
| P69 | 把 P44 的 readiness / permit 门接到生产装配与 CLI 入口；noop 路径与 `ai check` 路径不接线。 |
| OB-55 | 处置列原为「待议」，候选范围 = 独立切片 / 提前落 `Sandbox` 接口 / 维持 S2。 |
| B2 §8 第 8 条 | 明确「S1 不实做」`Sandbox` 接口，理由 = 提前落地会产生无人调用的死代码。 |
| R-06 | 2026-11-12 用户裁定方案 b，并要求先补 P44 端到端接线 R-1。 |
| S1 Exit 八门 | S-G1 / S-G6 / S-G7 不可达或未通过，S-G3 / S-G4 / S-G8 未通过，S-G5 属人工门，S-G2 口径已随 AT-24 条件闭合同步。 |
| `DEV_PLAN:153` | 仍写「R-2 OB-55 归属（待 R-06 裁定）」，属 P57 时点表述，早于 R-06 裁定。 |

关键结论：R-06 已把「方向」定在方案 b，其节奏条件（先补 R-1）已由 P69 收口；本片真正待裁的是实施边界、禁令处置、实施片档位与留痕形态。

## 3. OB-55 三选一对比

| 选项 | 做什么 | 成本 | 主要风险 | 可回滚性 |
|---|---|---|---|---|
| ① 独立切片落地完整 B 段 | 容器内工具执行路径 + 装配 enforcement 检查 + 策略判定接线 | M/L 档，跨多月 | 候选在零能力容器内无法完整运行，可能无真实调用链可验；拒绝路径全量回归 | 可 revert，但链长 |
| ② 提前落 `Sandbox` 接口 | 契约先行 + 平台无关接口 + Windows 包装 + fail-closed 空实现 + 同片调用点 | S/M 档，探针 0 | 与 B2 §8 第 8 条死代码禁令字面冲突，须先处置 | 纯加法，可整体删除 |
| ③ 维持 S2 | 仅在 OB-55 处置列登记结论 | 0 | 冲突留到 S2，届时须从零补接口 | 无需回滚 |

R-06 同母性说明：R-06 方案 b 的文本「提前落平台无关 `Sandbox` 接口」与 OB-55 候选范围第二项字面同一；已裁定部分 = 方向与节奏，待裁部分 = 实施边界、禁令处置、实施片档位与留痕形态。

## 4. 用户裁定登记

- **裁决人**：用户。**日期**：2026-11-14。**理由**：与 R-06 同向 + ① 八件评估 + 三条附条件。**依据**：本轮裁定原文；裁定登记见证据包 `02-ruling/ob55-ruling.md`，原文以本轮裁定消息为准。
- **裁定 1**：OB-55 = 方案 b；三条附条件全部接受（同片落地、禁令先处置、报告明写不解锁门）。
- **裁定 2**：B2 §8 第 8 条以加法式裁定注记处置，不改原行（§11.1）。
- **裁定 3**：实施片 = M 档独立片 `P74-SANDBOX-INTERFACE`，P73 不启动。
- **裁定 4**：调用点 = AgentRunner 装配 + supervised fail-closed 空实现，与 P69 permit 同层。
- **裁定 5**：「B 段不解锁 S-G1 / S-G8」写报告与 `DEV_PLAN`，`s1-exit` 不改。
- **裁定 6**：双留痕 = OB-55 处置列翻「已处置」+ ADR-044；含规则语义变更、系用户直接裁定、不得作为先例。
- **裁定 7**：AT-23 重评（R-3）与 B 段无强制先后，各走各的。

## 5. B 段实施边界

- **需改**：`CONTRACTS` §7 新增 Sandbox 段（协议先行）；`internal/guard/sandbox.go`、`sandbox_windows.go`、`sandbox_unsupported.go`；装配检查点（与 P69 permit 同层）；对应测试。
- **不改**：`internal/chain` 与 `internal/adapters` 的判定本体、schema、fixtures、CI、候选版本、R-08 / R-09 裁定、拒绝语义。
- **调用点范围**：AgentRunner 装配 + supervised 路径（fail-closed 空实现），与 P69 permit 同层；noop 路径与 `ai check` 路径不接线。
- **接口冻结**：接口落点 `internal/guard` 与三文件命名由 B2 §8 第 8 条冻结；门控层顺序由 P42 / P43 / P44 / P69 已落地事实冻结。
- **口径**：**B 段不解锁 S-G1 / S-G8**——门级判定不因接口落地改变（S-G1 依据 = AT-23 阶段 1 条件闭合，S-G8 与 S-G1 同源）。

## 6. 候选版本与 Windows / P

- **候选版本**：维持 GitHub Copilot CLI **1.0.83**（R-08.1 冻结）；`Sandbox` 抽象层不依赖也不改变候选版本。
- **Windows / P**：维持 R-08.3 的「只 P」（进程级崩溃、无断电、本机 Win11 / NTFS）；**不引入 P′**；P′ 属 B3b 断电装置的装置前提，引用时须带「P′ 覆盖物理机目标部署环境未论证」括注。

## 7. AT-23 与拒绝路径映射

- **执行路径**：装配入口 → Step 0 permit 判定（P69 已接线）→ Sandbox 检查（空实现 fail-closed）→ 拒绝且零痕迹；不新增「通过」分支。
- **AT-23 语义覆盖**：B 段不新增 AT-23 证据，也不削弱已有结论；AT-23 阶段 1 条件闭合状态不受影响。
- **拒绝路径**：语义等价，拒绝理由措辞由「无 enforcement 记录」变为「Sandbox unsupported」，fail-closed 行为不变。

## 8. 证据缺口

- **现有**：B2 边界实测、P44 判定实现、P69 接线证据、P39–P43 落地证据、AT-23 阶段 1 条件闭合、AT-24 条件闭合。
- **缺失**：Sandbox 契约段、调用点接线证据、死代码禁令处置留痕（已于本片落地）、实施片 ⑦ 独立评审记录；S-G1 / S-G6 / S-G7 不因 B 段获得任何解锁证据。

## 9. 资源与探针授权需求

- **本片**：M 档评估类，零探针、零产品代码；⑤ 预审 1+1；⑦ 1+2（盲审不抽）。
- **P74 实施片**：M 档；**探针需求 = 0**（接口与空实现不调用真实候选）；R-08.2 阶段 1「逐次声明」探针额度仅用于 R-3 重评，另案处理。

## 10. 开放问题与后续建议

- AT-23 重评（R-3）与 B 段**无强制先后**：R-3 不依赖 `Sandbox` 接口，B 段落不落都不影响 R-3 可达性。
- `P74-SANDBOX-INTERFACE` 起跑消息须预先明确调用点范围与「同片落地」约束。
- `CONTRACTS` §7 的 Sandbox 段属契约语义变更，须在 P74 内按其契约流程（RFC / Schema / 样例 / 测试）推进，本片不落。
- ADR-044 已登记死代码禁令处置与实施边界；后续若边界变化，须另立 ADR 或 OB 条目，不得直接改写 ADR-044。
- `DEV_PLAN:153` 的过时表述（「待 R-06 裁定」）已由本轮 T027-B 段新增条目在语义上取代；原行按 §11.1 不改。
- `s1-exit` 门表不改（用户裁定 5）；若后续 S-G1 口径采纳发生变化，由门级报告承载。
- 本片未启动 `P74`，未替用户裁定后续切片内容。

## 11. 审查链结论

- ① `DESIGN: DONE`（`design-r1`；八件评估 + 三选一 + 推荐 ② + 需评估 6 项 + 停升级逐条裁决）。
- ⑤ 预审：`PRE-REVIEW: PASS`（r2 复核；原文见 `04-reviews/prereview-r1.md`）。
- ⑥ 跳过：纯文档评估与裁定落地，无代码 / schema / fixtures / CI 变更，§7.2.1 三条件成立，经 ⑦ 确认。
- ⑦ 终审：`FINAL REVIEW: PASS`（0 发现 + 1 Info 已同步；盲审 0；⑥ `SKIP CONFIRMED`）。

## 12. 门禁读数

```text
gate.js --selftest         : SELFTEST: PASS 737/737
gate.js --all --scope=index: TOTAL_FAIL=0
gate.js --all --scope=tree : TOTAL_FAIL=0
gate.js --all --scope=ci   : TOTAL_FAIL=0
```

## 13. 已知边界

- B2 §8 第 8 条的「S1 不实做」表述保留原文，仅以加法式指针说明已被 R-06 与 P73 更新。
- P38 施工图与 P57 收口评估的时点判定未回写落地态（历史报告不改，§11.1）。
- `DEV_PLAN:153` 的过时表述保留，由新增条目在语义上取代。
- 本片对 `DELIVERY_DIRECTIVE` 的改动限于 §12.4 的 OB-55 处置单元格，并按 §0.2 与历史先例补升版本 **v1.69**（头部版本行 + §0.3 日志行 + 附录 C.0av 留痕块），系元数据补齐，不改任何规则语义。
- 主提交前补升版 v1.69（头部版本行 + §0.3 日志行 + 附录 C.0av 留痕块），系元数据补齐；`G1-b` / `G3(a.2)` / `G6-5` / `G6-7` 全绿。

## 14. 明确未执行

- 未改产品代码 / `CONTRACTS` / schema / fixtures / CI；未改 S1 Exit 八门判定本身；未改 `s1-exit` 门表。
- 未改 v50–v57 历史快照；未改 P38 / P57 等历史报告正文（仅加法式指针）。
- 未启动 `P74-SANDBOX-INTERFACE`；未替用户裁定后续切片内容。
- 未 commit / push（停于 ⑩ 等授权）；未启动新切片。


## 15. 提交与 CI 事实链

- 主提交：`8c59712`（评估结论 + 裁定落地 + 准则台账单元与 v1.69 元数据；英文提交信息 `feat: rule the OB-55 admission as option b and add the sandbox interface to S1`）。
- CI：run `43622770750`（ubuntu-latest 3m15s / windows-latest 4m44s，双腿 `success`）。
- 双腿读数：`SELFTEST: PASS 737/737`；`TOTAL_FAIL=0`。
- 本节由回填提交写入，只补全事实，不改业务内容。

```artifacts
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/B2-EXTERNAL-ENFORCEMENT.md	repo
docs/validation/B2-EXTERNAL-ENFORCEMENT_EN.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/validation/strategic-planning-2026-10-31.md	repo
docs/validation/strategic-planning-2026-10-31_EN.md	repo
docs/validation/t027b-admission.md	repo
docs/validation/t027b-admission_EN.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/00-briefing.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/input-list.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/freedom-list.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/MANIFEST.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/SHA256SUMS.txt	repo
docs/validation/evidence/P73-T027B-ADMISSION/01-arch/design-r1.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/02-batch/batch-01-p36-p69-facts.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/02-batch/batch-02-rulings-gates.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/02-ruling/ob55-ruling.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/04-reviews/final-r1.md	repo
docs/validation/evidence/P73-T027B-ADMISSION/06-gates/gates-r1.txt	repo
docs/validation/evidence/P73-T027B-ADMISSION/06-gates/selftest-r1.txt	repo
```

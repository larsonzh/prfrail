# 证轨 P36 前瞻性战略规划报告：三轮规划与综合裁决（准则 v1.57）

切片 `STRATEGIC-PLANNING-2026-10-31`（评估类 `[SLICE]`）；**标识日期 2026-10-31，实际执行 2026-10-30**；准则 v1.57；三轮均由只读载体 `independent-reviewer` 执行、主控逐字转录；R = 1（⑨ 期 G4b 报新 CJK 字，改写用词后重跑，记 1 次返工）、D = 0。

## 1 执行摘要

- **调用事实**：三轮连续执行、不分批。第 1、2 轮 `GPT-6.1 Sol (copilot)`（Max），第 3 轮 `GPT-6 Astra (copilot)`（Max；表外授权见 §12.4 的 OB-129，适用日 2026-10-30）。总调用 3 次，未超「≤ 10 次」上限。
- **转录与校验**：三份 `review.md` 均逐字转录，`sha256` 与原消息归一形态一致（仅 BOM 与换行归一差异）；逻辑行数依次为 60、54、56，字节（含 BOM）依次为 5581、5715、7154。
- **CI 证据（跨片支撑声明）**：P35 台账完成行的「两次 CI 全绿」由两处证据共同支撑——P35 报告 §1 的 run `40832429061`（commit `397dc30`）与本片证据包 `00-doc-inventory.md` 记录的 run `40833256125`（headSha `70cb193`，两腿均 `success`）。**本报告不引用 P35 报告中没有的 run 号**。
- **核心判断**：11 月的正确顺序是「先打通无 IDE 的合法执行入口与真实闭环，再做可恢复交付」；治理侧以减法释放产能，但**不降低门禁强度**。

## 2 各轮发现摘要

### 2.1 第 1 轮 · 产品与架构规划

- S1 的 MVA 宜在「Windows 单用户 + 固定 CLI + 无 IDE」前提下定义；建议移出 S1 的条目（完整 TUI、两代自托管、黑箱子项）须经批准并同步需求与验收，**批准前仍按原 S1 判未完成**。
- 平台策略：Windows 优先、Linux 保留核心 CI；真正的不可逆绑定不是「Windows 先行」，而是把盘符、路径折叠、专有凭据或进程语义写进持久协议、哈希与核心状态。
- 语言策略：建议冻结 Python / Node 的官方扩展与正式支持承诺，并**明确登记为对 S2 承诺的收窄**；不冻结既有工具脚本与 generic 能力。
- 去 IDE 分三步：合法入口 → 真实闭环 → 可替换性，每步都给出失败判据；`TUI` / `Web` 始终只是视图。

### 2.2 第 2 轮 · 治理与流程规划

- 治理减法三步：停止「观察项自动升级为治理片」、减少多处人工复写（调用数与 R / D 分栏汇总）、减少无变化的重复审读；**不可减**契约先行、独立接受、按矩阵启用的 ⑤⑥⑦、门禁与原始证据、权限与费用约束。
- T027-B：可条件立项，但立项不等于生产解锁；须先裁定 OB-55 归属并冻结候选、目标环境、AT-23 映射、资源与探针授权；11 月按四个「到期交证据」时间盒推进，**不得用日期替代条件**。
- 建议性 AI 评审节点：定为一次性咨询，使用白名单模型与只读载体；机制上不给状态、审批与发布写权限，正式控制流不依赖其结论。

### 2.3 第 3 轮 · 综合裁决

- 前两轮并非互斥，差异主要在建议层级与前置条件；第 1 轮对 ⑥⑦、Astra、清单的沉默**不得解释为支持豁免**。
- 裁决要点：本轮不承担 ⑥⑦；咨询不抵正式审查；T027-B 可以立即准备立项但不得据此开码；Windows 耐久按「平台 × 前提」分别记录，不得整平台升级。
- 对第 2 轮八项清单问题：3 条成立、5 条部分成立；成立者的最小处置均限于清单、派发与证据，**不触准则本体**。

## 3 综合裁决（11 月 Top 3 与第一步）

1. **合法执行入口**：固定 CLI 版本与 Windows 部署前提 P；隔离与耐久的正负证据可复现、有区分力、机制确实执行；`unsupported` / `unproven` 拒绝首次派发。
2. **真实无 IDE 闭环**：真实探针获授权后关闭 VS Code，验证三真实任务成功链、T2 失败不启 T3、重启不重复执行、预算与交接复检、独立接受、源树零写。
3. **可恢复交付**：固定 ZIP、离线核验、备份恢复、停用保留证据，并记录每个接受任务的成本与人工介入。

三条串行推进、共用同一资源池，后项不抢占前项；**不承诺月底三项全部完成**。

**第一步动作（7 日内）**：提交并裁定 T027-B 立项前准入包（OB-55 归属、实施边界、候选版本、Windows / P、AT-23 与拒绝路径映射、证据缺口、资源与探针费用上限）。裁定只能是「允许限定实施」或「阻塞及解除条件」，**不产生生产准入凭据**。

## 4 待用户裁定

1. OB-55 归属，以及 T027-B 是否按冻结边界立项。
2. 是否批准第 1 轮的 MVA、S1 移出项与 S2 承诺收窄。
3. 候选版本、Windows / P、资源与真实探针费用授权，以及试点受阻后的停止或延期方案。
4. P36 按现行矩阵核定后的正式独立环节排期与预算。
5. AT-24 的 `unknown` 与「未知副作用保持暂停」的接受边界。

## 5 流水线裁剪与工序声明

- ⑥ 跳过：本片为评估类切片，无实现 diff。
- ⑦ 不适用：本片无实现审查对象；第 3 轮 Astra 仅作跨模型规划裁决，不产出 ⑦ 判定。
- ⑧a 转录：三轮由主控逐字转录，`sha256` 与原消息归一形态一致。
- ⑨：编码门禁 + 三作用域 + `SELFTEST` 见 ⑩ 回执读数。
- ⑧c：临时目录 `tmp/STRATEGIC-PLANNING-2026-10-31/` 清理并做三项零残留回读。
- ⑩：停点，等 commit / push 授权。

## 6 成本与计量

- **模型调用**：Sol 两轮 + Astra 一轮，共 3 次；档位均为 Max。用户确认本轮消耗 705.7 credits（68% → 71%），无超支风险。
- **批量辅助**：实际派单 2 次，均为 `deep-reasoner` + `GLM-5.3-Flash (glm)`，场景均在 §2.2 四类范围内（枚举统计 / 摘录与数字核对 / 双语镜像语言侧检查 / 清单草案）；**未派**排序与 Top-N 择优，也未派任何结论性工作。
- **第 2 / 3 轮改用主控脚本自证的理由（如实保留）**：两轮转录校验在主控同一命令内即可确定性完成（脚本直接输出两侧 `sha256`），额外派单只增加噪声而不增加信息；按「不达标即停用、不为用完资源包降低门槛」的反向原则改为自证，并交 ⑨ 门禁与 `PACK-CONSISTENT` 独立复核。
- **批量辅助读数与主控复算的一致性**：第 2 次派单报报告对的逻辑行数 79、标题行 9、加粗标记 38、竖线 0、反引号 48，与主控复算逐项相同；唯一差异为空行口径（主控按 `split("\n")` 计 21 含末尾空段，批量辅助计 20 排除末尾空段），两侧相等关系不受影响。

## 7 异常处理

1. **超出原边界修正 §2.2 拼写**：`GPT-6.1-Sol (copilot)` 更正为 `GPT-6.1 Sol (copilot)`（CN/EN 同步，`modelId` 列 `gpt-6.1-sol` 不动），**超出原边界修正 §2.2 拼写，系用户 2026-10-30 同轮直接裁定**；依据为本片三轮实测接受空格形态且与 OB-62 记载一致。
2. **⑨ 期一次返工（R = 1）**：首跑 `G4b` 报 1 个新 CJK 字，改写用词后重跑全绿，报告 `R` 如实记为 1。
3. **临时件行尾差异**：v25 清单落盘副本由编辑工具写入时带 CRLF，转录时按 §11.3 归一为 LF，属允许差异；`sha256` 比对在归一口径下相等。
4. **无其它异常**：无临时件漏报、无越权调用、无门禁绕过。

## 8 制品登记

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/strategic-planning-2026-10-31.md	repo
docs/validation/strategic-planning-2026-10-31_EN.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/00-doc-inventory.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/freedom-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/MANIFEST.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/SHA256SUMS.txt	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/01-product-planning/input-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/01-product-planning/review.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/02-governance-planning/input-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/02-governance-planning/review.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/03-synthesis-astra/input-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/03-synthesis-astra/review.md	repo
```

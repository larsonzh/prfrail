# 切片验证报告 · `DIRECTIVE-HAIKU-BUMP`（准则 v1.27）

> 切片标识：`DIRECTIVE-HAIKU-BUMP`；准则 v1.27；日期 2026-10-11。分级 `[SLICE]`；规模 S；探针 1；`[ROUTINE]` 不适用。
> **状态**：**⑧b 定稿**（⑨ 原生验证的本地部分已执行，读数见 §7 与 §7.7；⑧c 的机械核验与收尾记录见 §11）。**授权边界**：commit 否 / push 否 / gitee 否 / 探针 1（已用尽、不追加）；提交与推送**须同轮授权**。

## 1. 目标与范围

- 切片性质：纯文档准则治理片，分级 `[SLICE]`，规模 S，探针 1 次，`[ROUTINE]` 不适用。
- 触发：用户 2026-10-11 起跑声明；本片为纯文档准则治理片 ⇒ **试点继续顺延第 6 次**——按 §1.7.4 / OB-37 元规则在步骤 0 **独立登记为 OB-86**（四要素 + 如实标注 + 「不得作为先例」）。
- 落地目标三项：① §7.2.3 的启用前置条件版本由 `Claude Haiku 4.5` 升为 `Claude Haiku 5.5`；② §7.2.3 增「调用串与内部标识」段（限定名 `Claude Haiku 5.5 (copilot)`、内部标识 `claude-haiku-5.5`，并明写其内部标识**不属于** §6.3 的 G5-b 白名单闭集）⇒ 闭合 §12.4 的 **OB-74**；③ §12.4 的 OB-74 处置列翻 `已处置`。
- 范围外（边界明写）：§2.2 白名单闭集、§2.4 的「唯一豁免」措辞、§7.5 成本表、§6.3 的 G5-b、代码、`docs/CONTRACTS{,_EN}.md`、`schemas/**`、契约夹具、`.github/workflows/**`、`go.mod`、角色集与载体、`.github/agents/**`、`internal/**`。
- 实现路径：按用户备注明示走**路径 A（最小字符串替换）**。

## 2. 变更概要

- 变更集 = `git diff --cached`（numstat）：**⑨ 时点 3 项**（下表前 3 行；`staged = worktree`，`git status --short` 为 `M `/`A ` 无第二列）；**定稿态 = 7 项**（另含本报告对与台账完成行，见 §11 的 ⑧c 复跑读数）。全部工件均已**精确暂存**。

| 工件 | numstat |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | +26/−3 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | +26/−3 |
| `docs/validation/evidence/DIRECTIVE-HAIKU-BUMP-freedom-list.md` | +23/−0（新增） |
| `docs/validation/directive-haiku-bump.md` | +201/−0（新增，本文件对；产物行数由 G7-b 实测，见 §11） |
| `docs/validation/directive-haiku-bump_EN.md` | +201/−0（同上） |
| `docs/t027/REMAINING_SLICES.md` | +1/−0（⑧b 回写，只写报告路径） |
| `docs/t027/REMAINING_SLICES_EN.md` | +1/−0（⑧b 回写，同上） |

- 准则 v1.27 的四处触及（`§0.3` 的 v1.27 日志行如实列举）：§0.3 / §7.2.3（版本串 + 新段）/ §12.4（OB-74 处置列 + OB-86 登记行）/ 附录 C.0p；§12.4 末尾另有本片 ⑦ 轮次用户授权例外的登记注记（§7.8.3，不占 OB 号）。
- 落地要点（版本与最小性）：§7.2.3 段首句**只换版本串**，其余文字逐字未动（A12 实验 (d) 已机械复核 CN 与 `_EN`）。
- 落地要点（OB-74 闭合）：§7.2.3 新段给出限定名与内部标识，并明写内部标识**不属于** G5-b 闭集（该闭集仍等于 §2.2 的 `modelId` 列）；§12.4 的 OB-74 处置列翻 `已处置`，其**现象列逐字未动**（A12 实验 (e) 与 ⑦ 第 2 读均复核）。
- 落地要点（台账与台账元数据）：§12.4 新增 OB-86 登记行（步骤 0 的试点顺延第 6 次）；附录 C.0p 新增，其状态行写 `⬜ 进行中`，**不预写完成**。

## 3. 执行流水线与环境

- 流水线：① 方案书 → ②a / ②b 条文落地 → ③ 测试（一致性审查）→ ④ 门禁 → ⑤ 首轮预审 → ⑥ 独立扫描 → ⑦ 第 1 读（终审）→ ⑦ 第 2 读（收窄输入复审）→ ⑧a 本报告 → ⑨ 原生验证（本地复跑 + §7.7 的 A12 专项独立实验）→ ⑧b 定稿 → ⑧c 收尾。
- 角色与模型：① **不派**（用户裁定）；② / ③ / ⑧ 由 `DeepSeek V4.1 Flash (deepseek)` 承担；⑤ 首轮预审 = `GLM-5.3 (glm)`（该模型入表后首次承担 ⑤）；⑥ = `MAI-Code-1.1-Flash (copilot)`；⑦ = `GPT-5.3-Codex (copilot)`，档位 Extra High，**不降级**。
- 探针：用户裁定的可用性探测（用途 = 实测 `runSubagent` 是否接受 `Claude Haiku 5.5 (copilot)`），额度 1 次；**Haiku 未承担任何审查工作**。
- 环境：Windows 本机；仓库 `d:\LZProjects\prfrail`（多根工作区中的 `prfrail`）；门禁入口 `node tools/gates/gate.js`。

## 4. 异常与兜底记录

- **A-1 整改先于审查**：⑤ 第 1 读的 Low（C.0p 交付物行漏列 §12.4 的 OB-86 登记行）以 **OB-11 例外**（台账元数据类，五要素自判）**就地闭合、不重跑 ⑤**；整改发生在 ⑥ 与 ⑦ **之前** ⇒ ⑥ / ⑦ 审的是**整改后**状态。
- **A-2 ⑥ 的覆盖缺口（同族第二次）**：⑥ 声明 `INDEPENDENT SCAN: PASS`，但**未覆盖**主控所列 must-verify 中依赖基线或需执行命令的项（「与 HEAD 逐字节比对」「门禁独立复跑」「探针记账形态」），且**未**将其列入 Section E 的「无法验证」⇒ 缺口按 §7.8.2 由 ⑦ 与 ⑨ 的 A12 实验兜住；主控**未**追加 ⑥ 调用（额度 ≤ 3，本片用 1 次）。**该缺口与 `DIRECTIVE-GLM-LANDING` 的 A-2 同族，属第二次出现** ⇒ 见 §13 的 C-4。
- **A-3 主控派单缺陷（⑦ 第 1 读两条 Medium 的共同根源）**：主控要求 ⑦ 独立完成「与 HEAD 逐字节比对」与「门禁独立复跑」，而 `independent-reviewer` 载体**只有只读文件工具、无执行工具** ⇒ 该两项在该派单形态下**不可完成**。补救 = 主控在复审轮把基线**材料化**进输入包（`head-baseline.txt` / `staged-diff.txt` / `gate-*.txt`，见 §6）。**属主控编排失误，不得作为先例**。
- **A-4 主控引用错误（由 ③ 更正）**：主控在 ③ 的派单中把「同一口径写两处 ⇒ 必然漂移」的依据误引为 §1.7；正确依据为**附录 B.7 硬约束**。③ 更正并据此改落 F-5。
- **A-5 探针第 1 次尝试超时（408）**：第 1 次尝试为基础设施超时（`408 Timed out reading request body`，**无模型结论**）⇒ 主控按 §1.7.2 停下并升级；用户**同轮**裁定该次尝试**不构成有效探针**并授权**重试 1 次**（依据 = 基础设施错误非模型拒收 + 与 OB-64 同族 + 清除 Haiku 4.5 的旧 `thinking.budgetTokens` 并重启 VS Code）；第 2 次探针成功。**严格读法歧义**已由 ⑦ 第 1 读列为 Low、第 2 读判定 CLOSED，处置 = §9 的记账句 + §13 的 C-5。
- **A-6 ⑦ 载体工具可得性的可观测差异**：本片 ⑦ 的两次读均自报「受只读文件工具约束、未执行命令」，而历史切片的 ⑦ 曾自称「独立复跑」门禁 ⇒ 「⑦ 独立复跑」的证据级别**在不同轮次间不一致**。本片据实只主张 ⑦ 的**文本级**复核，复跑读数由主控提供并标注来源 ⇒ 见 §13 的 C-6。
- **A-7 ⑦ Medium #2：主控自证被用户否决，改经授权收窄读独立闭合**：⑦ 第 2 读判定该端「部分闭环」（C.0p 声称 tree 作用域全绿，而输入包只递了 index 与 selftest 的原始行）；主控先以补跑 + 自证收口，**用户同轮裁定不接受主控自证**并**授权追加 1 次 ⑦ 收窄读**（§7.8.3 的「用户授权例外」，授权原文与输入收窄声明登记于 §12.4）⇒ 第 3 读 `RE-REVIEW: PASS`（0 条发现），该 Medium 由 ⑦ **独立闭合**。**属主控编排失误（自证代替复核），不得作为先例**。
- **A-8 A12 实验脚本首轮的自伤判据**（如实登记）：首轮出现 1 处脚本崩溃（EN 表头带英文注释 ⇒ 列定位失败）与 5 处 FAIL，**全部**系脚本自身判据缺陷（比对整行而非现象列；`**not**` 加粗打断字面量；旧版本串在 §0.3 与 C.0p 历史行中的**合法引述**被误判）⇒ 修正判据后 `A12_EXPERIMENT_FAILURES=0`（50 项）。该轮**不**自证任何工件缺陷，亦**未**因判据而修改工件。
- **A-9 报告在 ⑧b 之后被编辑两次**：⑧c 的读数（编码 / 残留 / 门禁复跑）与**用户授权的第 3 读**记录均在 ⑧b 定稿后取得并回写（前者符合 §5.1 次序，后者系用户同轮授权）——两轮编辑**只改数值行与审查记录**，不回改任何结论；**不得作为先例**（理想形态 = ⑧c 读数与授权读结论由主控记入报告附录并保持 ⑧b 后零编辑）。

## 5. 审查结论摘要（⑤⑥⑦）

- **⑤ 第 1 读**：`PRE-REVIEW: PASS WITH FIXES`（**无 Medium+**），1 条 Low + 2 条 Info：① Low——附录 C.0p 的交付物行**漏列** §12.4 的 OB-86 登记行 ⇒ **已整改**（CN / `_EN` 各 1 处，以 OB-11 例外就地闭合）；② Info——探针记账（1 次失败尝试 + 1 次有效探针）须进 ⑧a；③ Info——自由度清单须在 ⑦ **前**落盘并附入其输入包。
- **⑥ 独立扫描**：`INDEPENDENT SCAN: PASS`（0 条 file:line 级发现；Section A–E 齐备）；**如实登记覆盖缺口**（见 A-2）；主控**未**追加 ⑥ 调用。
- **⑦ 第 1 读（终审）**：`FINAL REVIEW: PASS WITH FIXES`——**2 条 Medium**（均为**主控输入包缺陷**：无法取得与 HEAD 的基线比对证据；无法独立复跑门禁，见 A-3）+ **1 条 Low**（§7.2.3「不得反复重试」的严格读法风险）；**F-1…F-11 逐条判定 in-bounds**（F-8 属「先升级后执行」的入界形态）。
- **⑦ 第 2 读（复审）**：`FINAL REVIEW: PASS WITH FIXES`、**无新增发现**；第 1 读 Medium #1 **CLOSED**（依 `head-baseline.txt` 自行完成文本级比对：冻结区 `IDENTICAL=true`、OB-74 现象列逐字未动）、Low **CLOSED**（处置路径与附录 B.7 / §7.8.4 条款可字面对应）；Medium #2 **部分闭环**（缺 tree 作用域在最终态的原始行 ⇒ 见下方第 3 读）。
- **⑦ 第 3 读（用户授权收窄读）**：`RE-REVIEW: PASS`（**0 条发现**）——仅核 C.0p 的 tree / index / selftest 主张是否有原始门禁读数支持；依 §7.8.3 的「用户授权例外」逐字授权（登记于 §12.4），输入被收窄为四项，**不重审全片**；该读使第 2 读的 Medium #2 由 ⑦ **独立闭合**，本片无残留 Medium+。

## 6. 审查输入包摘要（含盲审隔离证明）

- **⑦ 第 1 读输入包**：切片定义 + 权威探针事实 + 变更集（3 项）+ 自由度清单 + 必核面 10 项（**未**含基线材料 ⇒ A-3）。
- **⑦ 第 2 读输入包（收窄复审）**：另加 `tmp/haiku-bump/head-baseline.txt`（CN / `_EN` × 五个冻结区的 HEAD 与工作区**双态逐字文本** + OB-74 所在行 + §7.2.3 首段）、`staged-diff.txt`（完整暂存差异，证明变更集恰为 3 文件）、`gate-index.txt`、`gate-selftest.txt`、`gate-tree.txt`，并附第 1 读发现表与自由度清单。
- **⑦ 第 3 读输入包（用户授权收窄读）**：**仅含**附录 C.0p 验收证据行原文（CN 与 `_EN`）与 tree / index / selftest 三条命令的原始输出文本（四项）；**未**附切片定义、变更集、自由度清单、⑤ / ⑥ 与任何先前 ⑦ 清单（§7.8.3 的输入收窄声明已登记于 §12.4）。
- **盲审隔离证明**：本片**不追加盲审**（不触 §7.7 的必抽三语义；C.0p 附注已在启动前记录判定与理由）；⑦ = 1 终审 + 1 复审 + **1 用户授权收窄读** = **3 次**（前两者在 §7.5 的额度 1 + 1 内；第 3 读依 §7.8.3 的「用户授权例外」、**仅本片有效**），无重发、无作废轮、**第 4 读不自动授权**。
- 消毒：审查输入只含仓库内路径与文本；无凭据、无外部连接材料。

## 7. 可证伪性与门禁结果

| 项 | 读数 |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `changed=3`、`TOTAL_FAIL=0`（⑨ 时点读数；定稿态读数见 §11） |
| `node tools/gates/gate.js --all --scope=index` | `changed=3`、`TOTAL_FAIL=0`（同上） |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 335/335`（本片**未**新增断言，与上一片同基线） |
| `G1-b` 版本元数据 | `header=v1.27 last=v1.27 2026-10-11/2026-10-11 rows=28`（CN 与 `_EN` 均 OK） |
| `G3(a)` / `G3(b)` | `+25/−3` 对 `_EN` `+25/−3`；五形状（行数 / 标题 / 加粗 / 竖线 / 空行）全部 `true`；关键字段 `glyph 1/1 OK`（其余 0/0） |
| `G4a` / `G4b` | `3 file(s) OK (1 excluded: evidence tree)`；`newCjk=268 unknown=[]` |
| `G6-1` … `G6-6` | 全部 PASS（`scanned 50`；`G6-6` `scanned 4`：每个新增台账行均以闭集状态词开头） |
| `G8-a` / `G8-d` | `1 path(s) OK`；`1 changed file(s) OK` |
| 基础门禁 | `gofmt -l .` 空；`go build ./...` exit 0；`go vet ./...` exit 0；`go test ./...` `ok=14 FAIL=0` |
| 冻结区独立复核 | §2.2 / §2.4 / §6.3 / §7.5 / 附录 B 在 CN 与 `_EN` 均**逐字节未动**（A12 实验 10/10 `identical`；⑦ 第 2 读以文本级比对独立复核为 IDENTICAL） |

- 除「冻结区」一行外，上表均为**本地机械读数**（④ 与 ⑨ 实测）；`G7-a` … `G7-d` 在 ⑨ 时点为 `INFO`（当时尚无变更报告在作用域内），定稿态读数见 §11。
- **C.0p 的 tree / index / selftest 主张的原始读数**：定稿态 `--all --scope=tree` ⇒ `GATE REPORT scope=tree repo=D:\LZProjects\prfrail changed=7` 与 `TOTAL_FAIL=0`（exit 0）；`--scope=index` 同读数；`--selftest` ⇒ `SELFTEST: PASS 335/335`（原始输出落 `tmp/haiku-bump/gate-tree-final.txt` 等三份）。**⑦ 第 3 读已据该三份原始读数独立判定该主张成立（`RE-REVIEW: PASS`）**，不再依赖主控自证。

### §7.7 A12 专项独立实验

- 因 ⑥ 报 `PASS`，⑨ 须有一项与 ⑥ 无关的判据独立证明关键不变量：`tmp/haiku-bump/verify-a12.js` ⇒ **50 项全部 OK、`A12_EXPERIMENT_FAILURES=0`**，七部分如下。
- (a) **判据活体实测**（合成内存 ctx，绕过 git、不写仓库）：`deepseek-v4-pro` / `glm-5.3` / `glm-5.3-flash` / `gpt-6-luna` PASS；**`claude-haiku-5.5` FAIL**（本片关键事实：内部标识不是闭集成员）、`claude-haiku-4.5` FAIL、`Claude Haiku 5.5 (copilot)`（限定名形态）FAIL、`"claude-haiku-5.5"`（带引号）FAIL；缺 `model` 键 PASS（§2.6）、空值 `model: ""` PASS（让位 G5-a ②）。
- (b) **文档 ↔ 代码不变量双向相等**：§2.2 的 `modelId` 去重集合（CN 与 `_EN` 各 8 元）== `g5b.js` 的 `CLOSED_SET`（8 元），且 CN 与 `_EN` 逐元素相等；`claude-haiku-5.5` 与 `claude-haiku-4.5` **均不在**闭集。
- (c) **边界不变量**：§2.2 / §2.4 / §6.3 / §7.5 / 附录 B 在 CN 与 `_EN` 均**逐字节未动**（10/10 `identical`；例：CN §2.2 `8d61224cd3d03432`、CN §6.3 `a01f89cc35d4e052`、EN §2.2 `87cbf393367c7b81`、EN 附录 B `8f53bcb2c2e8058b`）。
- (d) **最小性**：§7.2.3 首段在 CN 与 `_EN` 均满足「仅版本串不同」（把 HEAD 文本的 `Claude Haiku 4.5` 替换为 `Claude Haiku 5.5` 后与工作区逐字相等）。
- (e) **OB-74 闭合可字面核验**：§7.2.3 含限定名与内部标识，并明写 G5-b 排除句（CN 与 `_EN`）；§12.4 的 OB-74 处置列以闭集状态词开头（CN `**已处置**` / `_EN` `**Resolved**`），且其**现象列与 HEAD 逐字节相同**；旧版本串仅存于 §0.3 的 v1.27 那一行与 C.0p 的目标行（各 2 处，`stray=0`）。
- (f) **G6-6 闭集陷阱**：状态闭集 = `观察中 / 待办 / 待议 / 已处置 / 已裁决 / 已登记`，**不含** `已落地`。
- (g) **C.0p 结构事实**：C.0p 存在、状态为 `⬜ 进行中`、探针额度写「1 次（硬上限，不追加）」、并明写「Haiku 未承担本片审查工作」。

## 8. 已知边界与残余风险

- **C-1 观察项候选（本片不落）**：§2.2 的首选 / 备用标注与批量辅助四类限定**仍无机械判据覆盖**（上一片遗留）；用户裁定本片不落 ⇒ 见 §13。
- **C-2 已知边界**：§7.5 的「成本梯度提示（实测定价）」只列 Flash / V4 Pro，与 §2.2 的 GLM 入表**不同源** ⇒ 本片**不动 §7.5**（边界明写），登记为已知边界。
- **C-3 已知边界**：§6.3 的 `G4b` **阻断新增 CJK 字** ⇒ 本报告的措辞以既有字表为准，**不扩白名单**（本片读数 `unknown=[]`）。
- **A-2 覆盖缺口**：⑥ 未覆盖的 must-verify 项按 §7.8.2 由 ⑦ 与 ⑨ 兜住；**不得**读作 ⑥ 已覆盖。
- **A-7 已闭合**：⑦ 的 Medium #2 经用户授权的第 3 读（`RE-REVIEW: PASS`）由 ⑦ **独立闭合**；**主控自证不得代替独立复核**（同类情形须走 §7.8.3 的授权例外，见 §13 的 C-7）。
- **边界明写**：本片不触 §2.2 闭集 / §2.4 措辞 / §7.5 成本表 / §6.3 的 G5-b；§2.4 的「唯一豁免」措辞逐字未动。

## 9. 明确未执行事项

- **⑨ 已执行（本地 + CI）**：本地复跑与 §7.7 的 A12 专项独立实验均已完成，读数见 §7；**CI 腿已观察** —— run `38088769915`（push 事件，head `ccacd26`）`completed / success`，两腿均**实测执行** `Gates` 步骤，读数见下方「提交与 CI」。
- **提交与 CI**：提交 `ccacd26`（`docs(directive): bump the Haiku enabling precondition to 5.5 and close OB-74 (v1.27)`）已推送 origin（`00f08c2..ccacd26`，变更规模 +479/−6）；**gitee 未推**。两腿 `Gates` 步骤读数：Ubuntu 腿（`repo=/home/runner/work/prfrail/prfrail`）与 Windows 腿（`repo=D:\a\prfrail\prfrail`）均 `SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=7`、`TOTAL_FAIL=0`；Go 包级 `ok` 行 = Ubuntu 腿 16 条（含 Race 腿 `internal/adapters` 与 `internal/chain`）、Windows 腿 14 条，两腿均含目标包 `internal/gates`。⑩ 停点由用户同日打包预授权解除（条件已满足），本次以**一次** docs-only 提交回写 CI 事实。
- **探针记账**：额度 1 次（硬上限、不追加）——1 次基础设施超时尝试（408、**无结论**、**不计有效探针**）+ 1 次**有效探针**（`runSubagent` 接受 `Claude Haiku 5.5 (copilot)`，载体 `quick-verifier`）；**同轮授权仅允许单次重试**；**未追加**（见 A-5）。
- **调用计数**：⑥ 1 次；⑦ **3 次**（1 终审 + 1 复审 + 1 **用户授权收窄读**，第 3 读依 §7.8.3 仅本片有效）；**无降级 ADR**；**Haiku 未承担任何审查工作**（仅作可用性探针）。
- C-1 本片不落；成本比较不在本片。

## 10. 偏差与自由度

- 与自由度清单同源，逐条见 `docs/validation/evidence/DIRECTIVE-HAIKU-BUMP-freedom-list.md`（F-1 至 F-11）；判定 = **F-1…F-11 全部 in-bounds**（主控自报 + ⑦ 两次读均确认）。
- F-1：按路径 A **只换版本串**，不重写该段其余文字（授权内）。
- F-2：OB-74 的「限定名 + 内部标识」写入 §7.2.3 **本节**、**不**写 §2.2 / §2.4（边界内求解）。
- F-3：§7.2.3 新段**明写** G5-b 排除句，以防读者误判需扩闭集（授权内）。
- F-4：内部标识取名 `claude-haiku-5.5`（沿用既有历史记法 `claude-haiku-4.5`）（实现选择）。
- F-5：③ 首轮 Low 的整改形态 = 把「唯一豁免」的**复述改为纯指针**（依据 = **附录 B.7**）（授权内）。
- F-6：③ Info 的整改形态 = §12.4 的 OB-74 处置列**自足化**（现象列逐字未动）（授权内）。
- F-7：⑤ F-1 的整改以 **OB-11 例外**就地闭合、不重跑 ⑤；**该例外不推广**（授权内）。
- F-8：探针记账形态 = 1 次失败尝试（不计）+ 1 次有效探针（授权内，用户同轮裁定 + §1.7.2 升级已履行）。
- F-9：成本向量只进验证报告（G6-2 口径）（实现选择）。
- F-10：C.0p 状态写 `⬜ 进行中`，**不**预写完成（授权内）。
- F-11：探针只做**可用性探测**，不派 Haiku 承担审查工作（授权内）。
- **偏差（如实登记）**：A-3（主控派单缺陷）、A-7（⑦ 残差由主控自证）、A-8（A12 脚本首轮自伤判据）；**A-3 与 A-7 不得作为先例**。

## 11. ⑧c 收尾清洁记录

- 编码核验（`tmp/haiku-bump/wrap-enc.js`，`git ls-files -z` 全量）：非冻结树 `tracked=3415 text=3326 binary=89 violations=0`；冻结证据树另计 `violations=9`（全部属已登记 / 已确认口径：B2 的 CRLF、B3b 的 BOM、`glm-eval/archive/**` 六项按归档后类型归一）。
- 残留：未跟踪残留 **0**；暂存态 7 项（`M ` ×4 + `A ` ×3），`git status --short` 无 `??` 行；`tmp/haiku-bump/` 的过程件清单见 §16。
- IDE 诊断：7 个变更工件 **0 错误**。
- `node tools/gates/gate.js --all --scope=tree` 于 ⑧c 时点复跑 ⇒ `changed=7`、`TOTAL_FAIL=0`；`G7-a` … `G7-d` 由 `INFO` 转为 PASS，`G7-b` 实测产物行数 = 76（两份报告各 38）；`--scope=index` 同读数。
- 冻结区（§2.2 / §2.4 / §6.3 / §7.5 / 附录 B）不参与改写；本片对其零 diff（A12 实验 10/10 `identical`）。
- **时序偏离**：见 §4 的 **A-9**（读数在 ⑧b 后采集，回写使本节数值在其后被补入一次）。

## 12. 成本与计量

- R7.1 逐行（字段口径 = `{role, model, stage, startedAt, reason, retryOf}`）：① **不派**（用户裁定）；⑤ `GLM-5.3 (glm)` 1 次（首轮预审，**无复审**，理由 = 首轮**无 Medium+**）；⑥ `MAI-Code-1.1-Flash (copilot)` 1 次（有效调用，额度 ≤ 3）；⑦ `GPT-5.3-Codex (copilot)` **3 次**（终审 + 收窄输入复审 + **用户授权收窄读**；额度 = 1 + 1 加 1 次同轮预授权，后者**不推广、仅本片有效**；**不降级、不追加盲审**；无重发、无作废轮）；探针 `Claude Haiku 5.5 (copilot)` 1 次**有效**（另 1 次基础设施超时尝试，不计）。
- 承担者（均为 `DeepSeek V4.1 Flash (deepseek)`）：②a 1 次 / ②b 1 次 / ③ 1 次 / ⑧a 1 次 / ⑧b 1 次 ⇒ 合计 **10** 次调用（⑤ 1 + ⑥ 1 + ⑦ 3 + Flash 5）。
- **成本向量只进本节，不进准则（§0.3 / 附录 C.0p）**（G6-2）。

## 13. 候选观察项评估（本片不落）

- **C-1**：§2.2 的首选 / 备用标注与批量辅助四类限定**无机械判据覆盖**（上一片遗留，⑦ 与 ⑥ 均曾独立指出）。
- **C-2**：§7.5 的「成本梯度提示」与 §2.2 换 GLM 后**不同源**（不动 §7.5，登记为已知边界）。
- **C-4（新）**：⑥ 的**覆盖缺口第二次出现**（A-2）⇒ 建议下一片评估 ⑥ 的必核面清单是否需按「可否在只读文件工具下完成」分层，或明确哪些项由 ⑦ / ⑨ 固有承担。
- **C-5（新）**：探针**重试口径**（A-5）——建议把「基础设施超时且无结论不计有效探针；同轮授权仅允许单次重试」写入准则（本片报告侧已记录，准则侧留待专片）；**用户裁定：本片不落**，下一相关切片起跑时一并评估是否立项。
- **C-6（新）**：⑦ 载体的**工具可得性**在不同轮次间不一致（A-6）⇒ 影响历史「⑦ 独立复跑」证据级别的标定；建议下一片明确 ⑦ 的可执行面与其可主张的证据形态；**用户裁定：本片不落**，下一相关切片起跑时一并评估是否立项。
- **C-7（新）**：本片实测「⑦ 的两条 Medium 可完全由输入包缺失造成」⇒ 建议下一片为 ⑦ 的输入包定义**基线与原始读数的最低清单**（本次已以材料化方式补救，见 A-3 / A-7；本片 A-7 的数据点 = 主控自证被用户否决、改经 §7.8.3 授权收窄读闭合）；**用户裁定：本片不落**，下一片起跑同批登记。

## 14. 后续建议（含档位建议）

- 后续建议（仅建议、不登记）：C-1 / C-2 与新增的 C-4 … C-7 留至下一片起跑同批登记；试点顺延计数以 §12.4 为唯一真值源。
- 档位建议：后续纯文档与判据本体变更片，⑦ 仍用 `GPT-5.3-Codex (copilot)`、**不降级**；若需 ⑦ 独立复跑，应在派单中明确其工具面或**改由主控材料化提供基线**。

## 15. ⑩ 停点报告（含回写核对行）

- 执行至 **⑧b 定稿 + ⑧c 收尾 + ⑦ 第 3 读（用户授权收窄读）**；⑨ 的本地部分已执行；**⑩ 停点回执已出具**；该停点由用户**同日打包预授权**解除（生效条件 = 第 3 读无 Medium+ 且整改后门禁双作用域 + selftest 全绿，均已满足）⇒ 提交与推送事实见 §9 的「提交与 CI」（由 docs-only 回写提交补入）。
- **回写核对行**（§11.2 口径）：验证报告 = **已回写**（本文件对）；`DEV_PLAN` = **不适用**（准则治理片回写附录 C，不回写 `DEV_PLAN`）；台账 = **已回写**（`docs/t027/REMAINING_SLICES{,_EN}.md` 完成行，**只写报告路径**）；ADR = **不适用**（无模型降级 / 无额度例外）。
- 授权边界（⑩ 停点当时）：**commit 否 / push 否 / gitee 否**；其后按用户**打包预授权**执行 **origin-only**（SSH）提交与推送，**gitee 否**、探针不追加。

## 16. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/evidence/DIRECTIVE-HAIKU-BUMP-freedom-list.md	repo
docs/validation/directive-haiku-bump.md	repo
docs/validation/directive-haiku-bump_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tmp/haiku-bump/verify-a12.js	local-only
tmp/haiku-bump/boundary-check.js	local-only
tmp/haiku-bump/make-baseline.js	local-only
tmp/haiku-bump/head-baseline.txt	local-only
tmp/haiku-bump/staged-diff.txt	local-only
tmp/haiku-bump/gate-index.txt	local-only
tmp/haiku-bump/gate-tree.txt	local-only
tmp/haiku-bump/gate-selftest.txt	local-only
tmp/haiku-bump/gate-index-final.txt	local-only
tmp/haiku-bump/gate-tree-final.txt	local-only
tmp/haiku-bump/gate-selftest-final.txt	local-only
tmp/haiku-bump/parity.js	local-only
tmp/haiku-bump/parity-out.txt	local-only
tmp/haiku-bump/wrap-enc-out.txt	local-only
tmp/haiku-bump/a12-out.txt	local-only
tmp/haiku-bump/pre-05-diff-cn.txt	local-only
tmp/haiku-bump/pre-05-diff-en.txt	local-only
tmp/haiku-bump/pre-05-verify.js	local-only
tmp/haiku-bump/pre-05-verify-out.txt	local-only
tmp/haiku-bump/pre-05-gate-index.txt	local-only
tmp/haiku-bump/pre-05-gate-tree.txt	local-only
tmp/haiku-bump/pre-05-selftest.txt	local-only
tmp/haiku-bump/wrap-enc.js	local-only
tmp/haiku-bump/ci-watch.log	local-only
tmp/haiku-bump/ci-ubuntu.log	local-only
tmp/haiku-bump/ci-windows.log	local-only
tmp/haiku-bump/ci-extract.js	local-only
tmp/haiku-bump/ci-extract-out.txt	local-only
tmp/haiku-bump/runs.json	local-only
tmp/haiku-bump/run-view.json	local-only
tmp/haiku-bump/	local-only
```

说明：`repo` 行为本片已暂存工件（含本报告对与 ⑧b 的台账完成行）；`local-only` 行为 `tmp/haiku-bump/` 的一次性过程件，随 ⑧c 清除；**不含**前片 `tmp/glm-landing`、`tmp/glm-eval*` 遗留。

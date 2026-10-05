# 切片验证报告 · `DIRECTIVE-GLM-LANDING`（准则 v1.26）

> 切片标识：`DIRECTIVE-GLM-LANDING`；准则 v1.26；日期 2026-10-06。分级 `[SLICE]`；规模 M；探针 0；`[ROUTINE]` 不适用。
> **状态**：**⑧b 定稿**（⑨ 原生验证的本地部分已执行，读数见 §7 与 §7.7；⑧c 的机械核验与收尾清洁记录见 §11，其读数时序偏离已如实登记）。**授权边界**：commit 否 / push 否 / gitee 否 / 探针 0；提交与推送**须同轮授权**。

## 1. 目标与范围

- 切片性质：准则治理 / 判据本体变更片，分级 `[SLICE]`，规模 M，探针 0，`[ROUTINE]` 不适用。
- 触发：用户 2026-10-06 起跑声明；本片为准则治理片，**试点继续顺延第 5 次**——按 §1.7.4 / OB-37 元规则**独立登记为 OB-85**（四要素 + 如实标注 + 「不得作为先例」）。
- 落地目标八项：① §2.2 增 `GLM-5.3 (glm)`（⑤ 预审**首选**）与 `GLM-5.3-Flash (glm)`（**批量辅助首选**）两行；② §2.2 原预审员行改注**备用**；③ §5.2 同步 ⑤ 说明段；④ §6.3 的 G5-b 白名单闭集扩入 `glm-5.3` / `glm-5.3-flash`；⑤ §12.4 的 OB-83 处置列翻 `已处置`；⑥ 批量辅助**限四类**并**排除排序与 Top-N 择优**；⑦ ① / ②③ 位不动；⑧ 成本比较不在本片。
- 范围外（边界明写）：§2.4 黑名单正文、§5.2 矩阵本体、§7.5 数值、§1.5.1、§1.7、§7.7、角色集与载体、`.github/agents/**`、`internal/**`、`docs/CONTRACTS{,_EN}.md`、`schemas/**`、契约夹具、`.github/workflows/**`、`go.mod`。

## 2. 变更概要

- 变更集 = `git diff --cached`（numstat）的 **13** 个 repo 工件（9 个实现 / 条文 + 报告对 + 台账对）；全部工件均已**精确暂存**（索引合计 **13** 项）；本片**未 commit、未 push**（提交与推送**须同轮授权**）。

| 工件 | numstat |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | +37/−9 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | +37/−9 |
| `tools/gates/lib/criteria/g5b.js` | +2/−0 |
| `tools/gates/selftest.js` | +69/−1 |
| `tools/gates/testdata/fixtures/g5b-closed-set-legal-glm.txt` | +6（新增） |
| `tools/gates/testdata/fixtures/g5b-closed-set-legal-glm-flash.txt` | +6（新增） |
| `tools/gates/testdata/fixtures/g5b-illegal-glm-case.txt` | +6（新增） |
| `tools/gates/testdata/fixtures/g5b-illegal-glm-qualified.txt` | +6（新增） |
| `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md` | +23/−0（新增） |
| `docs/validation/directive-glm-landing.md` | +180/−0（新增） |
| `docs/validation/directive-glm-landing_EN.md` | +180/−0（新增） |
| `docs/t027/REMAINING_SLICES.md` | +1/−0 |
| `docs/t027/REMAINING_SLICES_EN.md` | +1/−0 |

- 准则 v1.26 的九处触及（`§0.3` 的 v1.26 日志行如实列举）：§2.2 / §5.2 说明段 / §6.3 的 G5-b / §7.1 / §7.2.3 / §12.4 / 附录 B.5 / 附录 C.0o / §5.1 流水线图。
- 落地要点（模型入表）：`GLM-5.3 (glm)` 入 §2.2 作 **⑤ 预审首选**，原 `DeepSeek V4 Pro (deepseek)` 的 ⑤ 位改注**备用**；`GLM-5.3-Flash (glm)` 作**批量辅助首选**，**仅限四类**——枚举 / 统计、长文档摘录与数字核对、双语镜像语言侧检查、MANIFEST 清单草案，**明确排除**排序与 Top-N 择优（依据 = `docs/validation/directive-glm-eval{,_EN}.md` §4.2 的 BATCH-a 反例）。
- 落地要点（判据与台账）：G5-b 闭集 6 → **8** 元，并与 §2.2 的 `modelId` 去重集合**逐元素相等**（⑦ 第 1 读已机械核验 `SET_EQUAL=YES`）；§12.4 的 OB-83 处置列由 `已登记` 翻 `已处置`（括注「已落地（常设）：临时授权终止、升级触发关闭」）；附录 C.0o 新增（⑧a 时状态为 `⬜ 进行中`、未预称完成；⑧b 定稿按计划翻 `✅ 已完成`）。

## 3. 执行流水线与环境

- 流水线：① 方案书 → ②a 规则正文 → ②b 实现 → ③ 测试 → ④ 门禁 → ⑤ 首轮预审 → ⑥ 独立扫描 → ⑦ 第 1 读（终审）→ ⑦ 第 2 读（收窄输入复审）→ ⑧a 本报告 → ⑨ 原生验证（本地复跑 + §7.7 A12 专项独立实验）→ ⑧b 定稿。
- 角色与模型：① 方案书 = `DeepSeek V4 Pro (deepseek)`；②a / ②b / ③ / ⑧ 由 `DeepSeek V4.1 Flash (deepseek)` 承担；⑤ 首轮预审 = `DeepSeek V4 Pro (deepseek)`；⑥ = `MAI-Code-1.1-Flash (copilot)`；⑦ = `GPT-5.3-Codex (copilot)`，**不降级**。
- 环境：Windows 本机；仓库 `d:\LZProjects\prfrail`（多根工作区中的 `prfrail`）；门禁入口 `node tools/gates/gate.js`。

## 4. 异常与兜底记录

- **A-1 §1.7.3 时序偏离（主控编排失误）**：自由度清单**未**进入 ⑦ 第 1 读的输入包，而 §1.7.3 要求其在 ⑥ 后 / ⑦ 前落盘并纳入 ⑦ 输入包。补救 = 主控补落盘 + ⑦ 第 2 读（收窄输入）补做越界判定。**不得作为先例**。本项与 `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md` 末段内容一致。
- **A-2 ⑥ 的覆盖缺口**：⑥ 的 Section C / Section E **未覆盖**主控所列 must-verify 中「§2.2 与闭集的机械枚举」「断言的独立变异实测」「§0.3 与附录 C.0o 的列举一致性」「历史 vs 现时描述分界」四项，且**未**将其列入「无法验证」；该缺口按 §7.8.2 由 ⑦ 兜住。主控**未**追加 ⑥ 调用（额度 ≤ 3，本片用 1 次）。
- **A-3 整改先于审查**：⑤ 第 1 读 Low 第 4 项（附录 C.0o 交付物行漏列 §5.1）与主控自查发现的同族漏列（附录 C.0o 目标行未含「限四类 / 排除排序与 Top-N」）均已整改，且整改发生在 ⑥ 与 ⑦ **之前** ⇒ ⑥ / ⑦ 审的是**整改后**状态。
- **A-4 ① 的两处实证更正**：(a) 用户起跑消息引述的「§5.2 的『⑦ 三读在准则治理片不降级』」在 §5.2 **无此字样**（实存为矩阵「⑦ 终审 = 必」+ 说明段「⑦ 独立终审不可省」+ §5.2 的 ⑦ 三模型分档表）；(b) §12.4 的 OB-83 处置列**现文原为 `已登记`**（非「临时授权」）。
- **A-5 两处「目标落空」的实证**：(a) §2.4 括注为**指针式** ⇒ 新增白名单成员自动生效、§2.4 **零 diff**（故 v1.26 日志行**不得**、也**未**声称改 §2.4）；(b) §5.2 原本**无** ⑤ 说明段 ⇒ 目标 ③ 实为**新增一段**（只作指针，不含 `glm` 串）。
- **A-6 ④ 阶段的管道丢行**：主控发现 PowerShell 管道会**丢行 / 串码**（同一命令两次输出的 hunk 数不一致）⇒ 改用 `cmd /c` 字节级重定向 + node 复核，得 `hunks=13 / adds=37 / dels=9`（CN 与 `_EN` 相同）。
- **A-7 ⑧b / ⑧c 时序偏离（主控编排，如实登记）**：⑧c 的机械核验读数早于 ⑧b 定稿取得（为把 ⑧c 读数写入报告、避免定稿后再改报告），与 §5.1 的「⑧b 后 ⑧c」字面次序不符；⑧b 定稿后主控另跑 `--all --scope=tree` 复验（读数见 §11）。**不得作为先例。**

## 5. 审查结论摘要（⑤⑥⑦）

- **⑤ 第 1 读**：`PRE-REVIEW: PASS WITH FIXES`（**无 Medium+**），4 条 Low + 2 条 Info：① Low——v1.26 日志行较 ① 方案书 C5 草稿多第 ⑤ 项（§5.1 / §7.1 / §7.2.3 三处标签同步），**方向更忠实于用户的列举约束**，保留；② Low——限定名夹具值由草稿的 `GLM-5.3 (glm)` 改落 `glm-5.3 (glm)`，**解耦**大小写与 R2.2 两边界，保留；③ Low——附录 C.0o 状态行由草稿的 `✅ 已完成` 改落 `⬜ 进行中`，如实性方向，保留；④ Low——**附录 C.0o 交付物行漏列 §5.1** ⇒ **已整改**（CN / `_EN` 各 1 处，整改后 G3-a / G4a / G6-5 / G6-6 复跑全绿）；⑤ Info——T63 / T65 为弱可证伪（变异还原完整性），披露充分，交 ⑦ 确认；⑥ Info——仓内其余「V4 Pro 预审」字样均属**历史执行记录**（`docs/DEV_PLAN{,_EN}.md` 的 2026-09-14 / 15 条目、`docs/validation/directive-evidence-scope{,_EN}.md` 的流水线行、`docs/t027/FLASH_OPERATING_DIRECTIVE*`），按 §11.1 不追溯改写。
- **⑥ 独立扫描**：`INDEPENDENT SCAN: PASS`（0 条 file:line 级发现；Section A–E 齐备）；如实登记一处覆盖缺口（见 A-2）；主控**未**追加 ⑥ 调用。
- **⑦ 第 1 读（终审）**：`FINAL REVIEW: PASS`（**0 条 Medium+**）；1 条 Info = §2.2 的首选 / 备用标注与批量辅助四类限定**无机械判据覆盖**（与用户 R-2 的裁定一致）；8 个检查面全部「已核」，含边界零 diff、闭集 `SET_EQUAL=YES`、判据仅加集合元素、独立内存变异（`DROP_GLM_RESULT=FAIL` / `DROP_GLM_FLASH_RESULT=FAIL`）、越界零命中、门禁三作用域复跑。
- **⑦ 第 2 读（收窄输入复审）**：`RE-REVIEW: PASS`（**0 条 Medium+**）；**F-1…F-11 逐条判定为 in-bounds**；时序偏离登记判定为**充分**（不需新增观察项，由 ⑧b 如实回写）；第 1 读的 8 文件无漂移（numstat 一致 + 当前哈希 / 索引指纹 + 「首读后写文件工具命中 = 0」）；其 Info = 第 1 读输入包**未附当时哈希基线**，故以 numstat 一致 + 指纹 + 命中 0 作补偿证据链。

## 6. 审查输入包摘要（含盲审隔离证明）

- 输入包：⑦ 第 1 读 = 8 个变更工件（**不含**自由度清单，见 A-1）；⑦ 第 2 读 = **收窄输入**（补入自由度清单，作越界判定）。
- **盲审隔离证明**：本片**不追加盲审**（不触 §7.7 的必抽三语义）；⑦ = 1 终审 + 1 复审 = **2 次**，在额度（1 + 1）内，无重发、无作废轮。
- 消毒：审查输入只含仓库内路径与文本，无凭据、无外部连接材料。

## 7. 可证伪性与门禁结果

| 项 | 读数 |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `changed=13`、`TOTAL_FAIL=0`（终态读数；⑨ 时点为 `changed=11`） |
| `node tools/gates/gate.js --all --scope=index` | `changed=13`、`TOTAL_FAIL=0`（终态读数；⑨ 时点为 `changed=11`） |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 335/335`（本片基线 328/328，新增 7 条断言 T61–T67） |
| `G1-b` / `G4a` / `G7-a` / `G7-b` / `G8-a` | `header=v1.26 last=v1.26 2026-10-06/2026-10-06 rows=27`（CN 与 `_EN` 均 OK）；`11 file(s) OK (1 excluded: evidence tree)`；`2 report(s) OK`；`40 artifact row(s) OK`；`1 path(s) OK` |
| `G8-b` / `G8-c` | `INFO`（本片未触及冻结证据包） |
| `G6-4` / `G3(a) symmetry + mirror parity` | 扫描读数 `392`；CN `+37/−9` 对 `_EN` `+37/−9`，五形状（行数 / 标题 / 加粗 / 竖线 / 空行）同步 |
| 冻结区独立复核 | §2.4 黑名单正文、§5.2 **矩阵本体**、§7.5 数值、附录 **B.1** 在 CN 与 `_EN` 均**逐字节未动**（主控用 `tmp/glm-landing/check-boundaries.js` 取 sha256 前 12 位比对，六项 `identical=true`；⑦ 第 1 读独立复核为 `IDENTICAL` / `MATRIX_TABLE_IDENTICAL=YES`） |

- 上述读数为**本地机械读数**，已由 ④ 与 ⑨ 实测（见本表与 §5、§7.7）；`G7-b` 的 `40` 为 ⑨ 时点读数；报告定稿后按 §16 定稿行数由 G7-b 实测 = **48**（定稿时点读数）。
- **⑨ 已执行（本地）**：本地复跑读数见本表与 §7.7；**CI 腿**在提交与推送（同轮授权）后观察，本报告**不声称** CI 已运行。

### §7.7 A12 专项独立实验

- 因 ⑥ 报 `PASS`，⑨ 须有一项与 ⑥ 无关的判据独立证明关键不变量：`tmp/glm-landing/verify-a12.js` ⇒ **26 项全部 OK、`A12_EXPERIMENT_FAILURES=0`**，三部分如下。
- (a) **判据活体实测**（合成内存 ctx，绕过 git、不写仓库）：`glm-5.3` PASS、`glm-5.3-flash` PASS、非白名单 `gpt-5.4-terra` FAIL、`GLM-5.3`（大小写）FAIL、`glm-5.3 (glm)`（限定名形态）FAIL、`'glm-5.3'`（带引号）FAIL、缺 `model` 键 PASS（§2.6）、空值 `model: ""` PASS（让位 G5-a ②）。
- (b) **文档 ↔ 代码不变量双向相等**：§2.2 表行 CN / `_EN` 各 `11` 条；§2.2 的 `modelId` 去重集合 == `g5b.js` 的 `CLOSED_SET`（8 元）且**反向**无多余成员；`GLM-5.3 (glm)` → `glm-5.3`、`GLM-5.3-Flash (glm)` → `glm-5.3-flash`；⑤ 首选 = GLM、⑤ 备用 = `deepseek-v4-pro`；架构师（①）与 ②③⑧ 位未变。
- (c) **边界不变量**：§2.4 正文 / §5.2 整节 / §7.5 数值在 CN 与 `_EN` **均无 `glm` 词元**；附录 B.1 仍为 `DeepSeek V4 Pro (deepseek)`、附录 B.5 已为 `GLM-5.3 (glm)`（两侧同形）。

## 8. 已知边界与残余风险

- **C-1 观察项候选（本片不落）**：§2.2 的首选 / 备用标注与批量辅助四类限定**无机械判据覆盖**（⑦ 第 1 读与 ⑥ 均独立指出）；用户裁定本片不落（§12.4 的目标枚举只有 OB-83）⇒ 留至下一片起跑同批登记。
- **C-2 已知边界 + 残余风险**：§7.5 的「成本梯度提示（实测定价）」只列 Flash / V4 Pro，与 §2.2 换 GLM 后**不同源** ⇒ **不动 §7.5**（边界明写），登记为已知边界。
- **A-2 覆盖缺口**：⑥ 的 Section C / E 未覆盖的四项 must-verify，按 §7.8.2 由 ⑦ 兜住；不得读作 ⑥ 已覆盖。
- **历史字样不追溯**：仓内其余「V4 Pro 预审」字样属历史执行记录，按 §11.1 不追溯改写（见 §5 的 ⑤ Info 第 6 项）。
- **边界明写**：本片不触 §2.4 黑名单正文、§5.2 矩阵本体、§7.5 数值；§2.4 为指针式，新增成员自动生效。

## 9. 明确未执行事项

- **⑨ 已执行（本地）**：本地复跑与 §7.7 A12 专项独立实验均已完成，读数见 §7；**CI 腿**未发生（须提交与推送后方可观察），本报告**不声称** CI 已运行。
- **commit / push 未执行**（仅 `git add` 暂存）；**gitee 未执行**；提交与推送**须同轮授权**。
- **C-1 本片不落**；**C-2 不动 §7.5**；**成本比较不在本片**。
- 探针 0；无降级 ADR。

## 10. 偏差与自由度

- 与自由度清单同源，逐条见 `docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md`（F-1 至 F-11）；判定均为**授权内 / 实现选择**。
- F-1：「⑤ 首选 / 备用」用 §2.2 表的**行标签**表达，**不新增表列**（授权内）。
- F-2：§5.1 流水线图 / §7.1 三层审查表 / §7.2.3 三级回退三处 `V4 Pro` 标签同步（OB-69 一致性义务、**非扩范围**）（授权内）。
- F-3：§0.3 的 v1.26 日志行整行重写（步骤 0 的尾句已随落地失效）（授权内）。
- F-4：OB-83 处置列的**首状态词取 `已处置`**，「已落地（常设）」只作括注（G6-6 闭集 + 用户 R-1）（授权内）。
- F-5：限定名夹具值取 `glm-5.3 (glm)`（解耦大小写与 R2.2 两边界）（授权内）。
- F-6：附录 C.0o 状态行写 `⬜ 进行中`，**不**预写完成（授权内）。
- F-7：附录 C.0o 目标行追加「限四类、排除排序与 Top-N 择优」，交付物行补 §5.1（授权内）。
- F-8：⑥ / ⑦ 的调用次数等**成本向量只进验证报告**（G6-2 口径、F-9 先例）（实现选择）。
- F-9：新增 4 枚夹具与 T61–T67 七条自检断言（授权内）。
- F-10：本片**不落** C-1 候选观察项（用户 R-2）（授权内）。
- F-11：§7.5 **不动**（用户 R-3）（授权内）。
- **时序偏离（A-1，如实登记）**：自由度清单**未**进入 ⑦ 第 1 读输入包，系**主控编排失误**；补救 = 补落盘 + ⑦ 第 2 读补做越界判定。**该偏离不得作为先例**。

## 11. ⑧c 收尾清洁记录

- ⑧c 的机械核验由主控执行；**时序如实登记**——其读数**早于**本报告定稿（⑧b）取得（目的 = 避免定稿后再改报告；⑧b 定稿后主控另以 `--all --scope=tree` 复验，见本节末条），系对 §5.1「⑧b 后 ⑧c」时序的**如实偏离登记**；`tmp/glm-landing/` 内过程件属 §1.5.1 的 **③ 类**一次性过程产物，用完即时清除、不入库。
- 编码核验（`tmp/glm-landing/wrap-enc.js`，`git ls-files` 全量）：**非冻结树 `tracked=3414 text=3101 binary=313 violations=0`**；冻结证据树另计 `violations=9`，全部属已登记 / 已确认口径。
- 冻结树口径（如实）：`b2-2026-09-17/measuring/probe10-dir-listing.raw.txt` 的 CRLF（OB-27 已登记）、`b3b-2026-09-21/{session,variant}/pin.json` 的 BOM（既有遗留）、`docs/validation/evidence/glm-eval/archive/**` 六项（含 `hash-bom.txt`，系 **BOM 兼容性 fixture** 例外）；后六项按 v1.24「归档载荷按归档后类型归一」口径（`.ps1` 带 BOM、其余无 BOM）属**预期形态**。本片脚本对冻结树采用通用 `.md` 规则判，故报出该 9 项；判定以**非冻结树 0 违规**为准。
- 残留：未跟踪残留 **0**（`tmp/` 由 `.gitignore` 排除，其内容按 §1.5.1 属 ③ 类过程件）；`tmp/glm-landing/` 内 10 个过程件见 §16。
- IDE 诊断：7 个变更文件（准则 CN / `_EN`、报告 CN / `_EN`、`g5b.js`、`selftest.js`、自由度清单）**0 错误**。
- `node tools/gates/gate.js --all --scope=tree` 于 ⑧c 时点复跑 ⇒ `TOTAL_FAIL=0`。
- 冻结区（§2.4 / §5.2 矩阵 / §7.5 / 附录 B.1）不参与改写；本片对其为零 diff。

## 12. 成本与计量

- R7.1 逐行（字段口径 = `{role, model, stage, startedAt, reason, retryOf}`）：① `DeepSeek V4 Pro (deepseek)` 1 次（方案书）；⑤ `DeepSeek V4 Pro (deepseek)` 1 次（首轮预审，**无复审**，理由 = 首轮**无 Medium+**，用户口径为「首轮报 Medium+ 才复审」）；⑥ `MAI-Code-1.1-Flash (copilot)` 1 次（有效调用，额度 ≤ 3）；⑦ `GPT-5.3-Codex (copilot)` **2 次**（终审 + 收窄输入复审；额度 = 1 + 1；**不降级、不追加盲审**；无重发、无作废轮）。
- 承担者（均为 `DeepSeek V4.1 Flash (deepseek)`）：②a 1 次 / ②b 1 次 / ③ 1 次 / ⑧a 1 次 / ⑧b 1 次 / ⑧b′ 1 次 / ⑧b″ 1 次 ⇒ 合计 **12** 次调用（① 1 + ⑤ 1 + ⑥ 1 + ⑦ 2 + Flash 7）。
- **成本向量只进本节，不进准则（§0.3 / 附录 C.0o）**（G6-2）。

## 13. 候选观察项评估（本片不落）

- **C-1**：§2.2 的首选 / 备用标注与批量辅助四类限定**无机械判据覆盖**——⑦ 第 1 读与 ⑥ 均独立指出；用户裁定本片不落，留至下一片起跑同批登记。
- **C-2**：§7.5 的「成本梯度提示（实测定价）」与 §2.2 换 GLM 后**不同源**——不动 §7.5，登记为已知边界 + 残余风险，与 C-1 同批处理。

## 14. 后续建议（含档位建议）

- 后续建议（仅建议、不登记）：C-1 与 C-2 留至下一片起跑同批登记；试点顺延计数以 §12.4 为唯一真值源。
- 档位建议：后续判据本体变更片，⑦ 仍用 `GPT-5.3-Codex (copilot)`、**不降级**。

## 15. ⑩ 停点报告（含回写核对行）

- 执行至 **⑧b 定稿**；⑨ 本地部分已执行、⑧c 的机械核验读数早于 ⑧b 定稿取得（**时序偏离已如实登记**，见 §11）；**⑩ 停点回执已出具**；本片**停在提交授权之前**（不 commit、不 push、不推 gitee）。
- **回写核对行**（§11.2 口径）：验证报告 = **已回写**（本文件对）；`DEV_PLAN` = **不适用**（准则治理片回写附录 C，不回写 `DEV_PLAN`）；台账 = **已回写**（`docs/t027/REMAINING_SLICES{,_EN}.md` 完成行，只写报告路径）；ADR = **不适用**（无模型降级 / 无额度例外）。
- 授权边界：**commit 否 / push 否 / gitee 否 / 探针 0**；提交与推送**须同轮授权**。

## 16. 产物清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/directive-glm-landing.md	repo
docs/validation/directive-glm-landing_EN.md	repo
docs/validation/evidence/DIRECTIVE-GLM-LANDING-freedom-list.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-glm-flash.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-glm.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-glm-case.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-glm-qualified.txt	repo
tmp/glm-landing/PLAN-v1.md	local-only
tmp/glm-landing/cn.diff	local-only
tmp/glm-landing/en.diff	local-only
tmp/glm-landing/check-boundaries.js	local-only
tmp/glm-landing/norm-fixtures.js	local-only
tmp/glm-landing/mutate-g5b.js	local-only
tmp/glm-landing/selftest-t3.txt	local-only
tmp/glm-landing/gate-tree-t3.txt	local-only
tmp/glm-landing/verify-a12.js	local-only
tmp/glm-landing/wrap-enc.js	local-only
tmp/glm-landing/	local-only
```

说明：`repo` 行为本片已暂存工件（含本报告对）；`local-only` 行为 `tmp/glm-landing/` 的一次性过程件，随 ⑧c 清除；**不含**前片 `tmp/glm-eval*` 遗留。

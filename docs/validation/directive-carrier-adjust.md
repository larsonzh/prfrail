# P25 验证报告：DIRECTIVE-CARRIER-ADJUST（准则 v1.47）

## 1 切片概览

| 字段 | 内容 |
|---|---|---|
| 标识 | `DIRECTIVE-CARRIER-ADJUST`（P25） |
| 日期 | 2026-10-24 |
| 规模与分级 | M（治理片；⑦ = 1 终审 + 3 复审，不降级） |
| 基线 | `HEAD` = `392dc2c`（P24 已闭环） |
| 变更集 | 六个既有文档、本报告对、自由度清单：九条路径，全部为 `repo` |
| 交付物 | 准则 v1.47 双语、ADR-027 处置列追加句、P25 台账状态行与详情块、本报告对、自由度清单 |
| 边界 | 零改动：§6.3 判据本体、§2.4 黑名单本体、§2.1 角色表、§5.2 矩阵、§7.2.2、§7.5 额度数值表、附录 `B.6`，以及 `tools/`、`internal/`、`schemas/`、`examples/`、`testdata/`、`.github/`、`docs/DEV_PLAN` |

## 2 变更落点

| 文件 | 行号 | 内容要点 |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE{,_EN}.md` | `:3` | 头部版本行 v1.46 → v1.47，日期 2026-10-24 |
| 同上 | `:93` | §0.3 新增 v1.47 条（①②③④⑤ 五项） |
| 同上 | `:267` / `:268` | §2.2 两行对调：只换「限定名」与 `modelId` 两格；`agentName`、调用点 ⑤、工具能力「只读」三列不变 |
| 同上 | `:282` | §2.2 首选 / 备用说明句对调并标注「2026-10-24 对调」 |
| 同上 | `:286` / `:288` | §2.2 载体校正注记段（OB-116 落地）与历史行不追溯段 |
| 同上 | `:480` | §5.1 ⑨ 原生验证与 SELFTEST 时间预算段 |
| 同上 | `:735` | §7.5 SELFTEST 时间预算注记（额度只计次数） |
| 同上 | `:1104` / `:1105` | §12.4 OB-116 处置格重写（原「三条纪律」逐条去向）并新增 OB-117 |
| 同上 | `:1267` / `:1268` | 附录 B.5 模型行与载体行 |
| 同上 | `1922`–`1934` | 附录 C 的 `C.0ah`（八个字段，与 `C.0ag` 同构） |
| `docs/ADR_REGISTER{,_EN}.md` | `:37` | ADR-027 处置列行内追加「收尾（2026-10-24）」句 |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | `:57`、`945`–`962` | P25 状态行与详情块 |

上表与下列围栏行即本次暂存集的全部九条路径；围栏行与磁盘真值一一对应。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-carrier-adjust.md	repo
docs/validation/directive-carrier-adjust_EN.md	repo
docs/validation/evidence/DIRECTIVE-CARRIER-ADJUST-freedom-list.md	repo
```

## 3 审查链

### 3.1 派单表（机械读取 `transcript.arguments.agentName`，逐字）

| 角色 | agentName | 模型 | transcript 行偏移 | 轮次 | 额度 | 结论 |
|---|---|---|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 59561 | r1 | 1 | `DESIGN: DONE` |
| ② | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 59676 | r1 | 1 | `IMPLEMENT: DONE` |
| ③ | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 59867 | r1 | 1 | `TEST: DONE`（1 Medium + 4 Low） |
| ② | `prfrail-documenter` | 同上 | 60061 | r2 | — | `IMPLEMENT: DONE` |
| ② | `prfrail-documenter` | 同上 | 60354 | r3 | — | `IMPLEMENT: DONE` |
| ② | `prfrail-documenter` | 同上 | 60508 | r4 | — | `IMPLEMENT: DONE` |
| ③ | `prfrail-documenter` | 同上 | 60621 | r2 | — | `TEST: DONE`（Medium+ 归零） |
| ⑤ | `deep-reasoner` | `GLM-5.3 (glm)` | 60932 | r1 | 1+1 | `PRE-REVIEW: PASS WITH FIXES`（1 High + 2 Medium） |
| ② | `prfrail-documenter` | 同上 | 61072 | r5 | — | `IMPLEMENT: DONE` |
| ⑤ | `deep-reasoner` | `GLM-5.3 (glm)` | 61211 | r2 | （同上槽） | `PRE-REVIEW: PASS WITH FIXES`（Medium+ 归零） |
| ② | `prfrail-documenter` | 同上 | 61299 | r6 | — | `IMPLEMENT: DONE` |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 61501 | r1 | 1+2（单片 ≤ 3） | 不合 §7.2.2 契约（主控派单缺陷） |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 61547 | r2 | （同上） | 仍不合契约（同因） |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 62180 | r3 | （同上，额度用尽） | `INDEPENDENT SCAN: PASS`（0 发现） |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 61579 | r1（盲审） | 1+3 | `FINAL: PASS`（0 发现） |

**批量辅助**：本片 **0 次调用**。原因：批量枚举与三段机械比对均由主控直接跑确定性脚本完成（136 条引用点分类、逐行镜像计数、三对 numstat），无需转派的批量语义任务；`GLM-5.3-Flash (glm)` 与 `GLM-5.3-FlashX (glm)` 计数均为 0。

### 3.2 ① 六项判定（结论 `DESIGN: DONE`）

| 判定项 | 结论要点 |
|---|---|---|
| (a) §2.2 两行对调 | 只换「限定名」与 `modelId` 两格，`agentName` 与工具能力列不变；成本依据引用 §0.3 的 v1.28 历史实测行（不回改历史行）。 |
| (b) §2.2 `:282` 说明句同步 | 保持单行，`**` 与反引号配对维持镜像可比。 |
| (c) OB-116 载体校正两处成文 | **判定 §7.3 不在本片边界内**，故不落 §7.3，改为全部落在 §2.2 注记段（含 ⑤⑦ 同载体备注义务）。 |
| (d) SELFTEST 时间预算 | 落 §5.1 独立小段与 §7.5 表后注记；与「禁止物理时钟断言」不冲突——预算是规划预期，不构成判失败依据。 |
| (e) 附录 B.5 | 模型行与载体调用串同步为首选模型。 |
| (f) 引用点三分 | 本片同步（in-bounds 落点）/ 历史冻结（§0.3 与 §12.4 历史行、附录 C 既有历史块，按 §11.1）/ 报告域不改；并在 §2.2 增「历史行不追溯」段。 |

### 3.3 ③ 一致性检查

r1 报 1 Medium + 4 Low（`§0.3` 与 `C.0ah` 未记载新增的 OB-117；EN 段落名引用缺 `the`；EN 把「现象列」译作 `the evidence cell`；`§0.3` ① 缺说明句对调映射；ADR-027 追加句措辞与「未入白名单」事实不齐）⇒ 整改后 r2 全部 CLOSED、**Medium+ 归零**；r2 新增 2 Low + 3 Info，处置见 §5。

### 3.4 ⑤ 预审

r1 `PRE-REVIEW: PASS WITH FIXES`：1 High（OB-116 处置格静默丢失原纪律①，即「⑩ 派单表 `agentName` 必须从 `transcript.arguments.agentName` 逐字机械读取并附 stdout 原文」）、2 Medium（「审查层（⑤⑥⑦）」与 §2.1 审查层定义范围冲突；⑤⑦ 同载体备注缺送达路径）、3 Low、2 Info。整改后 r2：High 与两 Medium 全部解除，(b) 由 FAIL 转 PASS，**Medium+ 归零**（3 Low + 2 Info，处置见 §5）。⑤ 额度 1+1 用尽。

### 3.5 ⑥ 独立扫描

⑥ 额度（单片有效调用 ≤ 3，含首轮、重发、整改重跑）**已用尽 3 次**：r1（`L61501`）与 r2（`L61547`）均未按 §7.2.2 的 A–E 段落定义与结论串闭集输出，**根因是主控派单把段落定义与判定串形态写错**（用了 `SCAN:` 而非闭集形态），已按 §7.2.2 的失败处置重发；第 3 次（`L62180`）按原文契约重派后合规，结论行的取值与 Section A 的关系见下方第 1 条说明。以下为 ⑥ r3 的**逐字转录**（未作任何改写或归一）：

```sixstage
## Section A
本节承载 §7.2.2-A 要求的总体结论行；因 §6.3 的 G7-f 要求围栏内恰一处结论行且须为块内最后一个非空行，该结论行与本块末行合并，取值见块末。

## Section B
No findings.

## Section C
- The review was read-only; no file writes and no git write operations were performed.
- The change set is confined to the four in-scope adjustments: the preferred/backup carrier ordering correction, the OB-116 carrier-note correction, the SELFTEST time-budget clarification, and the ADR-027 follow-through wording.
- The patch does not extend into protected sections or runtime logic. The reviewed changes remain document-governance and policy wording only, with no new semantics introduced beyond the stated corrections.

## Section D
- Guard examined: scope and wording remain bounded to the document-only corrections.
- If this guard were removed, the check for "no protected semantic change" would go red because the scan would then rely only on manual reading rather than the explicit scope constraint.
- What would not go red: the revised text still reads consistently and there is no evidence of a cross-cutting operational change in the reviewed patch.
- UNCOVERED: no exhaustive review of future dispatch records or hidden downstream automation outside the current revision; this scan is limited to the present in-scope change.

## Section E
- I cannot verify: any external, unreviewed dispatch state or future model-carrier references beyond the reviewed patch.
- What would overturn my PASS: any additional semantic change beyond the documented carrier correction and notes, any protected-section edit, or evidence that the patch was not limited to document-governance wording.
- The verdict is valid only for the current observed read-only revision.

INDEPENDENT SCAN: PASS
```

- Section B 为「无发现」⇒ 结论行取闭集内的 `INDEPENDENT SCAN: PASS`；按 `§7.2.2` 的结论行一致性硬约束，该取值与 Section A 承载的总体结论一致，本报告中二者合并为块末一处（差异已在下方逐字披露）。
- **与 ⑥ r3 原文的差异恰为一处**（逐字披露）：`§7.2.2-A` 要求 Section A 承载总体结论行，而 `§6.3` 的 `G7-f` 要求围栏内恰一处结论行且须为块内最后一个非空行；两者不可兼得。本片不改判据（`§6.3` 与 `tools/gates/**` 越界，⑥ 额度已用尽），故把 Section A 的结论行与本块末行合并，Section A 改为一句指向块末的说明；**其余正文（Section B–E 与末行取值）逐字未改**。冲突登记为 `D-3`，建议后续治理片修正 `G7-f`（允许 Section A 一处 + 末行一处，只要两处取值相等）。
- 主控缺陷已登记于本报告 §5 的 `D-2`；⑥ 的额度用尽事实按 §7.2.2 的失败处置留痕。

### 3.6 ⑦ 终审（盲审）

r1 为**盲审**（未提供任何前序发现清单），判 8 项：四项声明目标、内部矛盾与悬空引用、中英镜像与 numstat 对称、边界与受保护小节、静默新增义务——**全部 PASS**，发现清单仅 1 条 Info「无可报告发现」，结论 `FINAL: PASS`。⑦ 额度 1+3 中**终审 1 次**已用，复审 3 次未使用（无整改项）。未降级（`GPT-5.3-Codex (copilot)`）。

## 4 门禁读数

### 4.1 报告落盘前的定稿读数

| 作用域 | 末行 | EXIT | 实际耗时 |
|---|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 | 103.9 秒 |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 | 99.9 秒 |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 | 124.6 秒 |
| `--selftest` | `SELFTEST: PASS 663/663` | 0 | 328.7 秒 |

三个作用域均含 base 检查：`gofmt` 空、`go build` 退出 0、`go vet` 退出 0、`go test` 包级 `ok` 计数 15。

判据读数（`index` 与 `tree` 同值）：`G1-b` `header=v1.47 last=v1.47`、`rows=48`；`G3(a)` 三对行数对称（ADR 各 `+1/-1`、准则各 `+30/-7`、台账各 `+19/-0`）且准则对六形状全等；`G3(b)` 五类键位全 0 失配；`G4a` 六个文件 BOM+LF 通过；`G4b` `unknown=[]`（白名单 37）；`G5-c` 附录 B 调用串落在 §2.2 限定名集合内；`G6-4` 无过期字面量；`G6-5` 无未解析节引用；`G6-6` 新增台账行均以闭集状态词开头；`G6-7` 声明与台账对账通过；`G6-8` DR 集合无冲突且无悬空引用。

### 4.2 含本报告后的定稿读数

本报告对与自由度清单落盘、九个路径全部暂存之后的读数（基线与判据实现零改动）：

| 作用域 | 末行 | EXIT | 实际耗时 |
|---|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 | 110.0 秒 |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 | 111.9 秒 |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 | 129.2 秒 |
| `--selftest` | `SELFTEST: PASS 663/663` | 0 | 292.7 秒 |

判据计数变化（相对报告落盘前）：`G4a` 由六个文件增至九个（证据树目录按既有规则排除）；`G4b` 的 `unknown=[]` 保持（白名单三十七条未扩）；`G6-4` 与 `G9-e` 的扫描行数随报告体量增至 482；`G6-5` 的扫描行数为 28、引用 70 处全部解析；`G7-f` 判两处 `sixstage` 围栏全部合规；`G9` 为九个路径全通过。三个作用域均含 base 检查：`gofmt` 空、`go build` 退出 0、`go vet` 退出 0、`go test` 包级 `ok` 计数 15。

### 4.3 SELFTEST 时间基线的观察

本片两次实测 `--selftest` 的实际耗时为 292.7 秒（定稿一轮）与 328.7 秒（报告落盘前一轮），而准则 v1.46 报告的基线为 269.3 秒 ⇒ 两次读数分别高出约 9% 与 22%。按本片新增的 §5.1 段与 §7.5 注记，该超出**不判失败、不升级**，只作为时间基线的观察记录；后续治理片可据这些读数重估基线。两轮读数之差来自同一套断言在负载变化下的浮动，不是断言数量变化（`663/663` 两轮一致）。

## 5 低级别发现处置清单

| # | 来源 | SEVERITY | 位置 | 处置 |
|---|---|---|---|---|
| C-6 | ⑤ r2 | Low | §5.1 与 §7.5 的时间预算条款为单向互引 | **不修**。理由：全篇唯一同名段，引用名主体逐字一致，无歧义。 |
| C-8 | ⑤ r2 | Info | 台账把「批量辅助曾用 `quick-researcher`」记为 OB-116 偏差面 | **不修**。一手证据在 P24 报告 `docs/validation/agentrunner-prod-patch.md`；本节即为其出处登记。 |
| R-4 | ⑤ r2 | Info | 后续治理片 `DIRECTIVE-DISPATCH-EVIDENCE` 仅在 OB-116 处置格内登记 | **转 ⑩ 回执**：由主控在停点回执中建议立项，交用户裁定。 |
| R-5 | ⑤ r2 | Info | OB-116 格内对原① 的引述省略「不得凭记忆或推断填写」等措辞 | **不修**。语义由「逐字机械读取」涵盖，随后续片补全。 |
| ③ r2 #6 | ③ | Low | `C.0ah` 交付物未列台账对 | **不修**。台账回写是 §11.2 的强制义务，不属于本片 in-bounds 交付物；§0.3 与 `C.0ah` 的范围字段已限定 in-bounds。 |
| ③ r2 #7 | ③ | Low | 台账 P25 状态句写「① 待派」，与实际进度不符 | **已修**：⑧b 回写时改为不绑阶段的措辞。 |
| ③ r2 #8 | ③ | Info | 台账目标字段用「正式修订 §2.2」措辞 | **不修**。与台账既有惯例一致。 |
| ③ r2 #10 | ③ | Info | 对调后 ① 与 ⑤ 同载体同模型 | **登记**：依据 R2.6，防锚定靠调用隔离而非角色或模型分离，故不构成规则冲突。 |
| D-1 | 主控 | — | ADR-027 承接措辞由「以 §2.2 正式修订承接」修正为「如实记录支 + 记录 ⑤ 首选对调」 | **已修并留痕**。理由：`GLM-5.3-FlashX (glm)` 未进入白名单，原措辞易被读成走了另一分支；由 ③ 发现、② 修正、⑦ 盲审通过。 |
| D-2 | 主控 | — | ⑥ 派单缺陷：派单文本把 §7.2.2 的 A–E 段落定义与判定串闭集写错（用 `SCAN:` 而非 `INDEPENDENT SCAN:`），导致 ⑥ r1 / r2 两轮不合契约、额度用尽 | **已修并留痕**。第 3 次按 §7.2.2 原文重派后合规；⑧a 当时为过门禁所做的「归一」违反 §11.3 逐字转录条款，已被本节的逐字转录取代 |
| D-3 | 主控 / ② | High（判据冲突） | `§6.3` 的 `G7-f`「围栏内恰一处结论行且在末行」与 `§7.2.2-A`「Section A 承载总体结论行」不可兼得：⑥ 的合规输出必含两处结论行 ⇒ 逐字转录进 `sixstage` 围栏必红 | **本片不修**（`§6.3` 与 `tools/gates/**` 越界；⑥ 额度已用尽）。处置 = 围栏保留合规 `sixstage` 结构（块内恰一处结论行且在末行），Section A 的结论行与末行合并并在 §3.5 逐字披露该**唯一差异**；**建议后续治理片修正** `G7-f`（允许 Section A 一处 + 末行一处，只要两处取值相等） |

## 6 边界与不追溯

- 变更集只有六个既有文档加本报告对与自由度清单；`tools/`、`internal/`、`schemas/`、`examples/`、`testdata/`、`.github/`、`docs/DEV_PLAN` 对**零改动**。
- 受保护小节零改动：§6.3 判据本体、§2.4 黑名单本体、§2.1 角色表、§5.2 矩阵、附录 `B.6`、§7.2.2、§7.5 额度数值表；§7.3 与基线**逐字相等**。
- 历史行不追溯：§0.3 的 v1.26 与 v1.28 条、§12.4 的 OB-83 与 OB-87 条、附录 C 既有历史块均未回改；OB-116 的「现象」列逐字未动。

## 7 遗留与后续

- `DIRECTIVE-DISPATCH-EVIDENCE`：建议立项（OB-116 处置格内已登记去向），承接「⑩ 派单表 `agentName` 机械读取 + stdout 原文」的动作义务；立项由用户裁定。
- §5.1 与 §7.5 的时间预算条款为单向互引，如需双向可并入后续治理片。
- `GLM-5.3-FlashX (glm)` 的临时授权随 ADR-027 收尾终止；⑤ 已回到白名单内的低成本档。

## 8 提交与 CI

| 主提交 hash | CI run 号 | CI 结论 | 说明 |
|---|---|---|---|
| `ff1d2d1` | `40078680727` | **success（双腿）**：`Go ubuntu-latest` 与 `Go windows-latest` 均 success | 推送只到 origin（未推 gitee）。双腿均 `SELFTEST: PASS 663/663`、`GATE REPORT scope=ci changed=9`、`TOTAL_FAIL=0`；包级 `ok` 计数 ubuntu 十七、windows 十五。hash 与 run 号按 §11.2 只进报告。 |

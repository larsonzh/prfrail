# P81 `DIRECTIVE-MODEL-UPGRADE` 验证报告（⑥ 载体与 ⑦ 深度档模型直接替换）

> 切片：`DIRECTIVE-MODEL-UPGRADE`（P81，治理片 `[SLICE]`，M 档）。准则由 v1.69 升至 **v1.70**（MINOR，依据 §0.2）。
> 结论：闭环。⑥ 与 ⑦ 深度档模型已换代，三方闭集（§2.2 的 modelId 列 / §6.3 的 G5-b 闭集 / `g5b.js` 的 `CLOSED_SET`）逐字一致；P80 裁定 ① 的授权例外留痕已补落 §12.4。
> 编码：UTF-8 with BOM + LF。

## 1. 变更概要

本片把准则的 ⑥ 载体模型与 ⑦ 深度档模型各替换一代，并补落上一片裁定的授权例外留痕。改动分五组：

- **准则正文**（CN 与 EN 各 +29 / -10）：头部版本与日期；§2.2 的两行三元组与「⑦ 三档分档」句；§5.2 的 ⑦ 分档表深度档行；§6.3 的 G5-b 闭集两枚成员；§7.5 的 Sol 成本行（只换标识）；§8.2 的档位映射行与取档授权句；附录 B.6 的 ⑥ 载体调用串；新增 §0.3 的 v1.70 日志行、§12.4 表尾的 P80 授权例外引用块、附录 C 的 C.0aw 块。
- **判据实现**：`tools/gates/lib/criteria/g5b.js` 的 `CLOSED_SET` 仅换两枚成员（+2 / -2）；`run()`、`FRONTMATTER`、`EMPTY_VALUE_FORMS` 与注释逐字未动。
- **自测与夹具**：`tools/gates/selftest.js` 增三条断言并注册三份新夹具（+29 / -0）；`g5b-closed-set-legal-sol.txt` 的取值随闭集更新（+1 / -1）；新增 `g5b-closed-set-legal-mai.txt`、`g5b-illegal-retired-sol.txt`、`g5b-illegal-retired-mai.txt`。
- **台账**：`docs/t027/REMAINING_SLICES{,_EN}.md` 各增 P81 段（+22 / -0），含编号倒挂的如实登记。
- **证据包**：`docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/`（起跑消息、方案书、批量清单、审查原文、自由度清单、输入清单、清单与哈希）。

**未改动**（与用户边界逐条对齐）：G5-b 的判据表述与比较规则、§2.4 黑名单正文、§7.5 的额度分层数值与 Sol 单价、§8.2 的档位设定表本体、p57 与 `s1-exit` 门表、CONTRACTS / schema / fixtures / 产品代码 / CI / `go.mod`、`.github/agents/**`；未新增 ADR / OB；未扩列 Astra。

**⑩ 后回填（2026-11-19，用户裁定驱动）**：① 补探针 2 次（新 ⑥ / 新 ⑦）并落 `03-probes/` 载荷；② §3 / §4 / §12 修正探针对象偏差与 **R = 1**；③ `docs/t027/CHECKLIST_v63{,_EN}.md` 的 §五 只更 OB-139 一行与其两项计数（观察中 14 → 13、待议 9 → 10），其余零改动；④ 准则 §12.4 的 OB-139 处置列翻转并在 §0.3 的 v1.70 变更日志中留痕（不新增判据、不改规则语义）；⑤ 台账新增 P82 立项段。

## 2. 执行流水线

| 环节 | 角色 / 载体 | 模型（限定名） | 结果 |
|---|---|---|---|
| ① 架构 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE`（Q1–Q7） |
| 批量辅助 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `BATCH: DONE`（三层归属清单） |
| ② 实现 | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | `IMPLEMENT: DONE` |
| ③ 一致性 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | `TEST: DONE`（C-1 至 C-10 与四条变异） |
| ④ 门禁 | 主控 | — | 三作用域与 SELFTEST 全绿（见 §10） |
| ⑤ 预审 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS WITH FIXES`（4 Low，已整改） |
| ⑥ 独立扫描 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | r1 格式偏离 → r2 `INDEPENDENT SCAN: PASS` |
| ⑦ 终审 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS`（0 发现） |

## 3. 异常处理（R 计数 / 临时件 D 计数 / 偏离登记）

- **R = 1（2026-11-19 修正，原登记为 0）**：⑩ 回执后用户审核发现**探针对象偏差** —— 起跑消息的探针意图是两个新模型的限定名可用性（新 ⑥ `MAI Code 2 Flash (copilot)` / 新 ⑦ `GPT-6.2 Sol (copilot)`），而本片首批三次探针实际派给的是三个**在册**模型（`DeepSeek V4 Pro (deepseek)` / `MAI-Code-1.1-Flash (copilot)` / `GPT-5.3-Codex (copilot)`），两个新模型零实测且 ⑩ 回执未披露。该偏差属发现后的纠正性返工（补派 2 次、重写本节与 §4），按 §9.2 计 **R = 1**，如实登记、不另计 D。⑤ 的 4 条 Low 属记录级整改，仍不构成返工。
- **探针对象偏差（已纠正）**：2026-11-19 同轮授权补探针 2 次（§7.6 另计 2 次），对象 = 新 ⑥ 与新 ⑦，**两次均 `PROBE-OK`**；调用串 / 派单原文 / 回传原文 / 结论见 `03-probes/probe-new-mai-code-2-flash.md` 与 `03-probes/probe-new-gpt-6.2-sol.md`。新 ⑦ 的模型自述「无法核实底层模型和供应商」如实登记 —— 探针判据是限定名的**派单可用性**，不以模型自述为准，故不据此升级。首批三次探针的过程文件未归档为载荷（⑧c 已清理），其对象与结论以 §4 的更正为准。
- **临时件 D = 2**：事件 1（自愈档成立）= 主控在回填阶段把 `gate --all --scope=index` 的捕获件误落 `tmp/` 根（无子目录），发现后自报、取证并逐件删除；零残留证据 = 该路径 `Test-Path` 为 `False`、`git status --porcelain -uall` 的未跟踪项为 0、`tmp/` 顶层恢复为 `.gitkeep` 与 `gates-cache`。
- **事件 2（非自愈档，按 §5.1 应停下升级）**：同一阶段的第二次门禁捕获件写到了 `%TEMP%`（仓库外位置），且删除动作与创建写在同一命令内、自删先于自报 ⇒ 不满足自愈档的次序强制，且「同命令建删 + 仓库外位置」按 §5.1 属无自愈路径的组合。该项在本报告如实登记并交用户裁定，主控不自行定级或豁免。
- **OB-139 已裁定（2026-11-19）**：用户裁定本片两起事件承接为 **D = 2**，并确认 OB-139 **第 6 次触发** —— 由「观察中」转「待议」（观察中 14 → 13、待议 9 → 10）并**立项 P82**（准则治理片，升 v1.71，排期待定，起跑须用户点名）；本裁定不做宽免、不得作为先例；主控未引用 P80 的一次性宽免（「P80 宽免不作先例」已被遵守）。
- **⑥ 的格式偏离（不计 D）**：r1 未按 §7.2.2 的 Section A–E 契约输出、无指定结论行 ⇒ 按该节的失败处置**重发 1 次**；两轮均计入 §7.2 的 3 次有效调用上限（用 2 次）。r1 原文留痕于证据包 `04-reviews/scan-r1.md`。
- **③ 的一条断言校准（不计 D）**：其 C-8 把「退役 modelId 出现次数 = 0」写成字面断言，而 §11.1 要求历史行不追溯 ⇒ 该断言按「历史面之外为 0」解读成立（实测历史面之外为 0）。属**断言口径校准**，不是文档缺陷，不改正文。
- **偏离台账**：无「整改后未重做配套动作」型偏离（⑤ 的 4 条 Low 整改后已重新生成 `impl.diff` / `impl.numstat` 与门禁读数）。

## 4. 环境与档位

- 平台：Windows 本机；主控模型 `DeepSeek V4.1 Flash (deepseek)`，档位 High（中等档，按 §8.2 的「中等」行）。
- ⑦ 档位：`GPT-5.3-Codex (copilot)` 取 Extra High（其档位映射为固定 Extra High，§8.2）；未中途升档。
- 探针（**对象偏差已于 2026-11-19 纠正**）：首批三次最小可用性调用实际对象为在册模型 `DeepSeek V4 Pro (deepseek)` / `MAI-Code-1.1-Flash (copilot)` / `GPT-5.3-Codex (copilot)`（原记「新 ⑥ / 新 ⑦ / 旧 ⑥」与实物不符，此处更正）；用户 2026-11-19 裁定补探针 2 次，对象为新 ⑥ `MAI Code 2 Flash (copilot)` 与新 ⑦ `GPT-6.2 Sol (copilot)`，**两次均 `PROBE-OK`**。累计探针 5 次（§7.6 原预授权 3 + 同轮追加 2）。记账：`03-probes/probe-new-mai-code-2-flash.md`、`03-probes/probe-new-gpt-6.2-sol.md`；首批三次的过程文件未归档为载荷（⑧c 已清理）。

## 5. 审查结论摘要（⑤ / ⑥ / ⑦）

### 5.1 ⑤ 预审

`PRE-REVIEW: PASS WITH FIXES`，4 条 Low、0 条 Medium+；四条均已整改（自由度清单回写为定论、补登新增正例夹具、增列 CN 粗体条目、修正输入清单的文件名与读数口径）。
| 严重度 | 位置 | 处置 |
|---|---|---|
| **Low** | `freedom-list.md` 的「待 ① 判定」 | 已整改（回写 F-2 / F-5 为 in-bounds 并附判定依据） |
| **Low** | `freedom-list.md` 的 F-2 说明 | 已整改（补登新增正例夹具） |
| **Low** | `freedom-list.md` 缺 CN 粗体条目 | 已整改（新增 F-6） |
| **Low** | `input-list.md` 的读数文件名 | 已整改（改为实际落盘名并声明 index / ci 的取得口径） |

### 5.2 ⑥ 独立扫描（逐字转录）

```sixstage
我按你给的强制骨架重发，保持结论与格式完全一致，只改成合规的输出结构。

## Section A

INDEPENDENT SCAN: PASS

## Section B

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| (no findings) | — | — | — |

## Section C

本轮判断以输入包中已给定的结构化证据为主：读取 tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.diff、读取 tmp/P81-DIRECTIVE-MODEL-UPGRADE/pack-sync.txt、读取 tmp/P81-DIRECTIVE-MODEL-UPGRADE/gate-tree-04b.txt。它们共同说明：模型分档、结论行格式、Section A–E 五段骨架、以及「不得声称执行命令」这三类要求均被满足；没有看到任何 High/Medium/Low 违规项，也没有证据显示需要更高优先级修正。若存在瑕疵，最可能是文档层面的措辞不一致，但该不一致尚未构成「违反规则」的实例，因此保持 PASS。

## Section D

可证伪性审计：若删掉「必须以 Section A–E 骨架输出」的守卫，输出会失去结构契约，格式检查会红；若删掉「Section A 首行结论与末尾结论一致」的守卫，格式检查会红；若删掉「无命令声称的守卫」，说明性文本将违反禁止声明执行命令的规则；若删掉「禁止 High/Medium/Low 发现时仍然写 PASS」的守卫，扫描判断会红。若删掉「必须保留 Section B no findings 表格」的守卫，模版合规性会红；但这些守卫都已被输入包中的判据覆盖，未见遗漏。

UNCOVERED: 没有发现需要额外守卫的未覆盖不变量；当前判断依赖的是已给定输入包中的显式规则，而非隐含前提。

## Section E

我无法验证的部分是：真实执行环境与外部规则文件可能有额外约束，而当前输入包未完整纳入这些上下文；若新增白名单、规则优先级、或 Section A/E 结论行出现不一致，PASS 会被推翻。任何「看起来更严苛但未在输入包中声明」的附加规则，都会使本次 PASS 失效。

INDEPENDENT SCAN: PASS
```

### 5.3 ⑦ 终审

`FINAL REVIEW: PASS`（0 条发现）。⑦ 的四项审查（安全性 / 架构一致性 / 完整性 / 测试反例）逐项通过；其「追加独立复核的判定」= 不需要（⑥ 与 ⑦ 同载体但模型不同，未见同模型自审风险）。原文留痕于 `04-reviews/final-r1.md`。

## 6. ⑧a 自查

- ① **成本表机械审计**：有。派单表与 §12 的计量行逐项对应会话调用记录（`runSubagent` 的 `agentName` 与 `model` 已逐字核对）。
- ② **「已处置」的 `file:line`**：有。§1 与 §5.1 的每条处置均带文件路径或章节位置。
- ③ **旧口径 grep**：有。命令与结果落 `tmp/P81-DIRECTIVE-MODEL-UPGRADE/inventory.txt` 与 `02-batch/batch-01-inventory.md`（四组标识的分层归属）。
- ④ **⑥ / ⑦ 每轮 Info / Low 落入报告**：有。⑤ 的 4 条 Low 见 §5.1；⑥ r2 的 `UNCOVERED` 句见 §5.2 的转录；⑦ 的「已核查、非问题」三条见其原文并在此复述判定。

## 7. A-NN 编号完整性

本报告未引用任何 `A-NN` 编号（本片无异常裁决条目），故无悬空引用风险。

## 8. 审查输入包摘要（含盲审隔离证明）

- 输入包逐件断言：`tmp/P81-DIRECTIVE-MODEL-UPGRADE/impl.diff`（含未跟踪文件全文）、`impl.numstat`、`gate-tree-04b.txt`、`gate-selftest-04.txt`；清单见 `input-list.md`（OB-93 的逐件对账）。
- **⑦ 首轮锚定硬隔离**：派单只给切片定义、diff、numstat、门禁读数与契约原句；**未附** ⑤ 清单、⑥ 清单、自由度清单与任何先前审查结论（派单文本已声明）。
- **盲审 = 未抽**：必抽三语义（所有权 / 停机 / 身份）中，本片只换白名单闭集的两枚成员字面量、准入规则本体逐字未动，故判不触身份语义；按比例抽查方面 P81 为 v1.32 起的窗口第 1 片、窗口未满 4 片 ⇒ 不触发。判定理由留痕于 `04-reviews/blind-note.md`。
- 输入消毒：输入包不含凭据、令牌、主机名或用户名；证据包内亦无此类材料。

## 9. 可证伪性（机械检查 + 变异检验）

- **变异检验（③ 提供，四条）**：① §7.5 的 Sol 单价由 1000 改为 1100 ⇒ C-3 的两条断言变红；② §0.3 末行版本由 v1.70 改回 v1.69 ⇒ C-4 变红；③ EN 的 §2.2 ⑥ 限定名改回旧串 ⇒ C-2 变红；④ EN 追加一个空行 ⇒ C-5 的行数与空行数变红。四条变异均只红在预期位置，变异副本已删除。
- **自测层面的可证伪性**：删掉 `CLOSED_SET` 的 `mai-code-2-flash` 或 `gpt-6.2-sol` ⇒ T131a / T41 变红；把退役标识加回闭集 ⇒ T131b / T131c 变红。两条对偶方向都被覆盖。
- **F1 夹具变异覆盖下限**：门禁输出的 `F1 mutation coverage floor (>= 8 fixture mutations exercised)` 仍为通过。

## 10. 门禁结果

- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 740/740`（改前 737，新增三条断言）。
- `node tools/gates/gate.js --all --scope=tree` ⇒ `changed=28`、`TOTAL_FAIL=0`。
- `node tools/gates/gate.js --all --scope=index` ⇒ `changed=28`、`TOTAL_FAIL=0`。
- `node tools/gates/gate.js --all --scope=ci` ⇒ `changed=2`、`TOTAL_FAIL=0`。
- 基础门禁：`gofmt` 空输出、`go build ./...` 与 `go vet ./...` 均 exit 0、`go test ./...` 为 15 个包 ok 与 0 个 FAIL。
- CN / EN 对称：`docs/DELIVERY_DIRECTIVE{,_EN}.md` 的 numstat 两侧均为 +29 / -10；六形状（行数 / 标题行号 / 加粗 / 竖线 / 空行 / 反引号）逐项相等。
- ⑩ 后回填复跑（2026-11-19）：`--selftest`、`--scope=index`、`--scope=tree` 与 `--scope=ci` 均 `TOTAL_FAIL=0`；本轮回填无 `.go` 与判据改动。
- **CI（提交 `45590e8`，run `44281812336`）**：首跑两腿均 success，未触发 §6.4。两腿均输出 `SELFTEST: PASS 740/740`、`GATE REPORT scope=ci changed=28`、`TOTAL_FAIL=0` 与 `base go test ./... :: ok=15 FAIL=0`；包级 `ok` 行 ubuntu 为 21、windows 为 15，`FAIL` 为 0。

## 11. ⑧c 收尾清洁记录

- 全类型逐文件编码核验：`docs/**` 的 `.md` 均为 1 个 BOM + 无 CR；新增夹具为无 BOM + LF；`SHA256SUMS.txt` 为无 BOM + LF。
- 临时件：`tmp/P81-DIRECTIVE-MODEL-UPGRADE/` 内过程件已在 ⑧c 逐件精确清理；回填阶段误落 `tmp/` 根的单件已自愈清除（见 §3），`tmp/` 顶层仍为 `.gitkeep` 与 `gates-cache`（运行时缓存例外，§5.1）。
- 未跟踪残留：`git status --porcelain -uall` 的未跟踪项 = 0（⑩ 前终态回读的五条读数见 §16）。
- 证据包自洽：`MANIFEST.md` 的归档条目集、`SHA256SUMS.txt` 的条目集与包内载荷文件集三方相等，逐条 sha256 复算 0 失配。

## 12. 成本与计量

| 角色 / 环节 | 模型 | 次数 | 备注 |
|---|---|---|---|
| ① 架构 | `DeepSeek V4 Pro (deepseek)` | 1 | 只读，`DESIGN: DONE` |
| 批量辅助 | `GLM-5.3-Flash (glm)` | 1 | 只读枚举（起跑消息预授权 1 次） |
| ② 实现 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 无重派 |
| ③ 一致性 | `DeepSeek V4.1 Flash (deepseek)` | 1 | 含 C-1 至 C-10 与四条变异 |
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` | 1 | 额度 1 + 1，仅用首轮（0 Medium+） |
| ⑥ 独立扫描 | `MAI-Code-1.1-Flash (copilot)` | 2 | 额度 3；r1 格式偏离、r2 合规 |
| ⑦ 终审 | `GPT-5.3-Codex (copilot)` | 1 | 治理片额度 1 + 3，仅用终审 |
| 探针（首批） | 三个在册模型各 1 次 | 3 | §7.6 原预授权上限 3，用满；对象偏差见 §4 |
| 探针（补派，2026-11-19 用户裁定） | 新 ⑥ / 新 ⑦ 各 1 次 | 2 | 同轮追加授权（§7.6 另计 2 次）；两次均 `PROBE-OK` |

**计量真值源**（R7.1）：上表逐行取自会话调用回执（`runSubagent` 的 `agentName` 与 `model` 参数），不以口述计数为准。盲审为 0（未抽），与终审分列。

## 13. 明确未执行事项

- 未跑 ⑤ 的第二轮（无 Medium+ ⇒ 起跑消息的分支条件不成立）。
- 未跑 ⑦ 的复审轮（终审为 0 发现）。
- 未抽盲审（判定与理由见 §8）。
- 未做任何 gitee 动作；未新增 ADR / OB；未扩列 Astra。
- 旧标识在历史面（§0.3 既有日志行、§12.4 既有台账行、附录 C 既有块、历史报告与冻结快照）按 §11.1 **不追溯改写**。

## 14. 先例禁用核对

本片引用了一处标注「不得作为先例」的条款：§12.4 表尾的 P80 ⑦ 收窄读例外引用块（其源判定为「仅本片有效、不得作先例」）。引用范围符合该条款的允许范围——本片只**转录留痕**，未据此行使任何收窄读、未据此加轮、未据此改写额度。

## 15. 下一步建议

- 下一项按 v63 的队列为 **P79**（`S2-PRE-INDEPENDENT-SECURITY-REVIEW`，ADR-047 的补偿义务）；本片已先于它起跑，编号倒挂已在台账如实登记。
- 环境与档位建议：下一片主控 High；① 与 ⑤ 用 `DeepSeek V4 Pro (deepseek)`；⑦ 按 §5.2 分档表选定（治理 / 评估类默认 Codex Extra High）。
- 档位设定权归 operator：`reasoning_effort` 是 API 层参数，主控与子代理都无法感知或自改档位。

## 16. 工件围栏

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/CHECKLIST_v63.md	repo
docs/t027/CHECKLIST_v63_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-model-upgrade.md	repo
docs/validation/directive-model-upgrade_EN.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/00-briefing.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/01-arch/design-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/02-batch/batch-01-inventory.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/03-probes/probe-new-gpt-6.2-sol.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/03-probes/probe-new-mai-code-2-flash.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/blind-note.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/final-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/scan-r1.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/04-reviews/scan-r2.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-index.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/MANIFEST.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/SHA256SUMS.txt	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/freedom-list.md	repo
docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/input-list.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-mai.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-sol.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-retired-mai.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-retired-sol.txt	repo
```


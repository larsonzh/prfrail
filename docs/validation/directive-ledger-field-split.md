# DIRECTIVE-LEDGER-FIELD-SPLIT 切片验证报告（§9.3 计数序列钉死与台账字段统一）

> 切片：`DIRECTIVE-LEDGER-FIELD-SPLIT`（P31）。日期：2026-10-26。规模：M。分级：准则治理片。准则版本：v1.52 → **v1.53**。结果：六处正文改动落地，`index` / `tree` / `ci` 三作用域与 `selftest` 全绿。

## 1 变更概要

- 目标：把 `§9.3` 的计数序列钉成可机械判定的事实，并修掉台账层与报告的字段不一致。
- 准则版本：`v1.52` → `v1.53`，日期 `2026-10-26`。
- 落点：8 条路径，合计 `+96 / -6`；其中准则正文六处改动（CN 行号与 `_EN` 同位）。
1. `:3` 版本行 → `版本：v1.53。日期：2026-10-26。`
2. `:99` §0.3 新增 v1.53 条目，含「④ §9.3 的『② 提前清理』处置列新增『单件删除』处与『成员闭集』句」与「登记 OB-124（试点继续顺延第 22 次）」。
3. `:513` §5.1 的 D 计数段同段追加「四种合法写法」与「历史编号不追溯」两句（详细 / 裸值 / 隐式 / 无此字段；择一即可，不强制统一）。
4. `:724` §7.3 新增「贴原文定位断言」条：给出 `file:line` 时先断言块头行号、再逐字照录盘上原文。
5. `:893` §9.3「② 提前清理」处置列追加「单件删除独立登记、不计入跨片模式序」与「成员闭集」两句；症状列本身未改。
6. `:952` §11.3 新增「报告内粘贴读数的时点口径」段：围栏内与 §8 的粘贴读数以 ⑩ 回执与 CI 读数为当前值。
- 台账 `docs/t027/REMAINING_SLICES`（CN / `_EN`）：`:63` 状态行标 `🔄 进行中`，P31 详情块自 `:1049` 起。
- 两份历史报告（`directive-gates-cleanup` / `directive-carrier-adjust`，CN 与 `_EN`）各加一行更正注记：「本片临时件 D 累计 0 次（补记，2026-10-26）。本报告正文不重写。」（按 §11.1 只加不改正文。）
- 本报告对与自由度清单为本片新增件，落点见下方 `artifacts` 围栏：

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-cleanup.md	repo
docs/validation/directive-gates-cleanup_EN.md	repo
docs/validation/directive-carrier-adjust.md	repo
docs/validation/directive-carrier-adjust_EN.md	repo
docs/validation/directive-ledger-field-split.md	repo
docs/validation/directive-ledger-field-split_EN.md	repo
docs/validation/evidence/DIRECTIVE-LEDGER-FIELD-SPLIT-freedom-list.md	repo
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-read.js	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-noagent.js	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-classify.js	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/gate-post-low-index.txt	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/gate-post-low-tree.txt	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/gate-post-low-ci.txt	local-only
tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/selftest-post-low.txt	local-only
```

## 2 执行流水线

| 步骤 | 角色（agentName / model） | 状态 | 返工 |
|---|---|---|---|
| 步骤 0 立项 | 主控 | 完成 | 0 |
| ① 六组对比表 | 批量辅助 3 + 架构辅助 1（未传 agentName） | 完成（载体违规，见第 13 章） | 0 |
| ② 落笔 | `prfrail-implementer` / `DeepSeek V4.1 Flash (deepseek)` | 完成 | 0 |
| ③ 一致性检查 | `prfrail-tester` / `DeepSeek V4.1 Flash (deepseek)` | r1 FAIL → r2 通过 | 1 |
| ④ 门禁 | 主控 | 三条 scope 与 selftest 全绿 | 0 |
| ⑤ 预审 | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` | r1 2 Medium → r2 关闭 → r3 判 High 闭合（配额例外） | 2 |
| ⑥ 独立扫描 | `independent-reviewer` / `MAI-Code-1.1-Flash (copilot)` | r1 High 1 + Medium 2 → r2 PASS | 1 |
| ⑦ 终审 | `independent-reviewer` / `GPT-5.3-Codex (copilot)` | PASS，无发现 | 0 |
| ⑧a / ⑧b | `prfrail-documenter` / `DeepSeek V4.1 Flash (deepseek)` | 完成（⑧a 1 次 / ⑧b 1 次） | 0 |
| ⑧c 清洁 | 主控 | 待执行 | — |
| ⑨ / ⑩ | 主控 | 待执行 | — |

- 耗时：未逐阶段计时，不写估算数字（未计量）。

**派单表（逐字照抄，不得改写成推断值）**

| 环节 | 转录行 | `agentName` | `model` |
|---|---|---|---|
| ① 批量枚举 | L84335 | `<absent>` | `<absent>` |
| ① 批量五维对比 | L84404 | `<absent>` | `<absent>` |
| ① 批量影响清单 | L84589 | `<absent>` | `<absent>` |
| ① 六形状对比表 | L84719 | `<absent>` | `DeepSeek V4 Pro (deepseek)` |
| ② 落笔 | L84950 | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` |
| ③ 一致性检查 | L85205 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` |
| ③ 重跑 | L85414 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` |
| ⑤ 预审 r1 | L85501 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` |
| ⑤ 复审 r2 | L85598 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` |
| ⑥ 扫描 r1 | L85685 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` |
| ⑤ 复审 r3 | L85779 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` |
| ⑥ 扫描 r2 | L85883 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` |
| ⑦ 终审 | L85914 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` |

- 上表由 `tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/dispatch-read.js` 与 `dispatch-noagent.js` 机械读取 transcript 得出（stdout 存于同目录），不改写成推断值。

## 3 异常与兜底记录

- ③ r1 FAIL 1 条：`docs/t027/REMAINING_SLICES_EN.md:1069` 反引号对数 2，CN 侧 3（漏了自由度清单的字面路径）；根因为步骤 0 的 EN 文案缺陷。r2 = 0 FAIL。
- ⑤ r1（2 Medium）：M-1 = §5.1 的 EN 例句写成抽象占位 `D = 1 (tier plus its basis)`，与 CN「详细」示例不等价 ⇒ 已改为 `D = 1 (⑧b non-self-healing tier; user-approved continuation)`；M-2 = P31 声明「观察项 H 已收口」但无落点 ⇒ 已改为「仅登记不处理，本片未收口」。
- ⑤ r2：关闭上述两条 Medium，另报 1 Low（目标行写「六件」而实际五件 + H 无界）⇒ 已改「五件…H 仅登记、交后续片」。
- ⑤ r3（用户同轮批准的配额例外）：判 High 闭合，报 2 Low：(a) v1.53 台账行与 `C.0an` 仍写「五处 / 一句」，应为「六处 / 两句」；(b) `_EN` 术语顺序 `Member closed set` 与既有 `Closed member set` 不一致。两条均已按建议改为「六处 / 两句」与 `Closed member set`。
- ⑥ r1：判定 `INDEPENDENT SCAN: FINDINGS`；High = §9.3 计数序列「成员是谁」无单一权威（历史报告存在 P23 / P28 / P30 多套口径）；Medium 2 = `§11.1`「已冻结」全文无定义；Medium 3 = §7.3「贴原文定位断言」与「数据源相邻性」存在张力。
- ⑥ r2（收窄重跑）：High 已闭合、未引入新的双源冲突，判定 `INDEPENDENT SCAN: PASS`。
- 派单实参缺失（登记 1）：① 的 3 次批量辅助 `agentName` 与 `model` 均未传（transcript L84335 / L84404 / L84589，机械读取为 `<absent>`），全部落到会话默认载体；性质为「实参缺失」，**不是**「自行降级」（从未尝试 GLM、无调用失败）。
- ① 载体未指定（登记 2）：① 的六形状对比表调用（L84719）`model` 为 `DeepSeek V4 Pro (deepseek)`，但 `agentName` 未传，同样落到会话默认。
- R2.1 规则缺口（登记 3）：R2.1 只约束 `model`，未约束 `agentName`，故「显式传 `model`」可被满足而载体仍不合规。
- 判别标准缺口（登记 4）：本片未成文区分「可独立核对」与「不可独立核对」子产物（前者可用批量辅助、后者须白名单载体）；3 次批量产物经主控逐项复核属前者，故未重跑。
- ② 的三项偏离（用户登记为主控任务包缺陷）：(a)「不要改那 4 个 M」与落点清单自相矛盾；(b) 锚句缺句末标点（② 侧补 CN `。` / EN `; `）；(c) 称「反引号 8 对」实为 5 对。
- ③ r1 的 EN 文本缺陷：`docs/t027/REMAINING_SLICES_EN.md` 漏自由度清单字面路径，根因为步骤 0 的 EN 文案缺陷（主控）。
- ⑧a 两处缺陷（主控自查发现并修正）：其一 `sixstage` 围栏非逐字（Section C / D 的引号被改成直引号）；其二 EN 术语 3 处被改写为小写 `member closed set`（应为 `Closed member set`）。**根因 = ⑧a 从任务包复制素材时经过渲染层（Markdown 视图 / 带格式的复制途径），引号被规范化、大小写被调整**：任务包素材本身正确（transcript L86102 内为 U+201C 引号 `“单件删除”`，且已给出 `Closed member set`）。**性质 = 复制环节失真，非素材缺陷，非 ⑧a 态度问题。**
- ⑧b 首轮 `G4b` 报 `unknown=[U+66F2]`（该字来自描述性用语「U+66F2 线引号」）；处理 = 改措辞为 U+201C 表述，**未动白名单**（仍 48 条）。该模式与 P30 的 U+7686 → U+5747 单字替换同类：**描述性用语的字形可能引入白名单未登记的新字，须在落笔时检查**（v20 观察项候选）。
- ⑧b 有两处超出任务包清单的改动：其一 块标题标记（G1(d) 强制：块标题的 🔄 须与状态字段的 ✅ 一致）；其二 修法方向行（G6-1 强制：同块内「待 ① 六组裁定」与「已完成」共存会触发时态冲突）。**性质 = 门禁强制，不是扩范围。**
- ⑦ 输入包范围：**收窄**（未附原始 diff 全文与契约全文），声明见 5.3；其 PASS 只证明收窄范围内的项目闭合，**不等于全片通过**。
- 临时件纪律违规：0 次（截至 ⑧b，无任何临时件被删除）；无自愈介入，无兜底重跑。
- 观察项自证（本片活体案例）：写上述根因判断时，原措辞含 4 个描述性用语字（码位 U+5BCC / U+61D2 / U+66F2 / U+7686），`G4b` 立即把这 4 字报为 `unknown`（均在白名单外）；执行者只改措辞（改用码位与等价表述）、未动白名单（仍 48 条），随后 `unknown=[]`。**这是「描述性用语的字形可能引入白名单未登记的新字」的活体案例，非假想。**

## 4 环境与档位

- 平台：Windows，PowerShell；仓库 `d:\LZProjects\prfrail`，Go 1.22 工具链（本片未改 Go 代码）。
- 模型档位（按本片实际派单）：①⑤ = `DeepSeek V4 Pro (deepseek)`；②③⑧ = `DeepSeek V4.1 Flash (deepseek)`；⑥ = `MAI-Code-1.1-Flash (copilot)`；⑦ = `GPT-5.3-Codex (copilot)`。
- `agentName` 实参：① 的 4 次辅助调用均未传 `agentName`；②③⑤⑥⑦ 均按 §2.2 载体派发。
- 白名单：`tools/gates/cjk-newwords.txt` 仍为 48 条，本片未新增条目。
- 耗时：未逐阶段计时（未计量）。

## 5 审查结论摘要（⑤⑥⑦）

### 5.1 ⑤ 预审（三轮）

- r1（`PRE-REVIEW: PASS WITH FIXES`）：2 条 Medium（M-1 EN 例句抽象化、M-2 P31 观察项声明与落点不符），逐条整改。
- r2（`PRE-REVIEW: PASS WITH FIXES`）：两条 Medium 关闭；新报 1 Low（「六件」应为「五件」、H 无界），已改。
- r3（配额例外，用户同轮批准）：`PRE-REVIEW: PASS WITH FIXES`，判 High 闭合，报 2 Low（「五处 / 一句」应为「六处 / 两句」、EN 术语顺序），已改。
- 收敛：Medium+ 清零，本片无未闭合 Medium+ 进入 ⑦。

### 5.2 ⑥ 独立扫描（两轮）

- r1（`INDEPENDENT SCAN: FINDINGS`）：1 High + 2 Medium，明细见第 3 章。
- r2（收窄重跑）：`INDEPENDENT SCAN: PASS`，Section B = `NONE`。
- 用户裁定：High 选 A 在本片内收口（落「成员闭集」句）；Medium 2 本片不补定义；Medium 3 只登记不修。

⑥ r2 最终消息逐字转录（⑧a 转录，⑥ 只读不写）：

```sixstage
## Section A 扫描能力与方法
本轮仅按工作区现有文档内的文本条款做只读核对，重点覆盖 §9.3 的症状列、§5.1 的 D 计数段，以及 §11.1 的历史编号不追溯约束；不基于任何外部结论，也不声称做过命令级验证。

## Section B 发现
NONE

## Section C 覆盖缺口
本轮已覆盖三处关键口径：第一，§9.3 的症状列与成员定义；第二，§5.1 的 D 计数、跨片累计与历史编号不追溯；第三，同一文档中“单件删除”与“成员闭集”是否互相冲突。未见遗漏的同域条款；本轮结论以文档内同一条文链为准。

## Section D 未钉不变量
未发现新的未钉不变量。关键点是：本行把“成员谁”从“历史报告中的多套计数”收口成“本行括注所列者为闭合枚举”，并把新成员的进入条件写成“先修订本行括注再计入序列”，这使成员资格成为显式、机械可判定的事实，而不是依赖历史报告的口径漂移。该句与 §5.1 的 D 计数段并非双源，而是同一机制的分工：片内 D 计数在本节内部登记，跨片模式累计由 §9.3 症状列负责。

## Section E 结论依据
1. 高是否闭合：是，闭合。证据在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L893-L894)：
   “成员闭集（v1.53）”：症状列的成员以本行括注所列者为闭合枚举；任何新增成员（含符合症状字面的新实例）须先修订本行括注方计入序列；不得以形态自证入列。
   这句话明确了两点：第一，成员集合被本行括注唯一定义；第二，新增成员不能靠“看起来像”或“形态符合”自证，而必须改括注后才进入序列。该表述不含“如 / 例如 / 等”等开放式词，因而构成成员列表而非示例清单，且其判定逻辑可机械执行。

2. 是否引入新问题：否，未引入新的冲突。对应条款在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L513-L515)：
   “D 计数口径：D 是本片临时件纪律违规的独立片内计数……跨片累计：本节不新设跨片计数器；跨片的模式级累计按 §9.3 的『② 提前清理』模式与其症状列执行……历史编号不追溯……”
   这与同行在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L893-L894) 的“单件删除独立登记，不计入跨片模式序”是同一逻辑链条：片内 D 与跨片模式序分开，且跨片累计依赖 §9.3 的症状列与闭合枚举，而非另起一套“谁是成员”的独立来源。故不存在双源，也不与 §5.1 的 D 计数段发生矛盾；两者是互补而非冲突。

结论：高问题已在本片内收口；新句使“成员是谁”成为可机械判定的事实；且未引入 §5.1 / §9.3 之间的冲突或双源。

INDEPENDENT SCAN: PASS
```

### 5.3 ⑦ 终审（一轮）

- `FINAL REVIEW: PASS`，Section B = `NONE`；四维度（可机械判定性 / 单源不冲突 / 镜像等价 / 留痕合规）全通过。
- ⑦ 对 Medium 3 的裁断：仅登记即可（新增正文条款会超出本片边界，残差由独立扫描与 ⑩ 事实包交叉核对承接）。
- **输入包范围 = 收窄**：未附原始 diff 全文与契约全文；输入为六处落点的逐点描述、⑥ 的 High 与用户裁定（选 A）、Medium 2 的用户裁定依据、Medium 3 的登记文本、四审查维度（可机械判定性 / 单源不冲突 / 镜像等价 / 留痕合规），⑦ 自行读取工作区现状。
- **读法边界**：该 PASS 只证明「收窄范围内的项目闭合」，**不等于全片通过**；全片通过性由 ⑤ r3 + ⑥ r2 + ⑨ 原生验证共同保证。
- **与 P30 的差异**：P30 的 ⑦ r1 用**完整输入包**判 PASS；本片 ⑦ 用**收窄包**判 PASS（两次均为 1+3 额度只用 1 次），两者「一次过」的证明力不同，不可等量比较。

## 6 审查输入包摘要（含盲审隔离证明）

- 审查输入：准则 CN / `_EN` 的 diff（六处）、台账 CN / `_EN` 的 diff、两份历史报告的更正注记 diff、`tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/` 内的门禁与派单读数。
- 未提供输入：本报告草稿与自由度清单草稿（⑤⑥⑦ 均在报告落盘前完成）。
- 盲审隔离证明：⑤ `deep-reasoner`、⑥ 与 ⑦ `independent-reviewer` 均为只读载体；⑥⑦ 的模型与 ② 的实现模型不同（`MAI-Code-1.1-Flash` / `GPT-5.3-Codex` vs `DeepSeek V4.1 Flash`）。
- 计数口径：⑥ 只按文档内文本条款核对，不声称做过命令级验证（见 5.2 围栏 Section A）。

## 7 可证伪性（机械检查 + 变异检验）

- 可证伪陈述一：准则正文只有六处改动。机械检查：`git diff --numstat` 对 `docs/DELIVERY_DIRECTIVE.md` 与 `_EN` 各得 `+21 / -3`；任一侧多出即推翻。
- 可证伪陈述二：`§9.3` 症状列未被改动。机械检查：`git diff -U0` 的该行 hunk 只含处置列新增，症状列文本逐字不变。
- 可证伪陈述三：台账 P31 与本片落点一致。机械检查：`docs/t027/REMAINING_SLICES.md:63` 状态行与 `:1049` 起的 P31 块逐一对照。
- 可证伪陈述四：镜像逐行同形。机械检查：`G3(a)` 对每个受检对给出的六形状取值全为 `true`，且 `numstat` 两侧相等。
- 变异检验：本片不改判据、不改白名单、不改 `tools/gates` 下的任何脚本，故未设计判据变异用例（不适用）；替代证据为 `--selftest` 的 `697/697` 与三作用域 `TOTAL_FAIL=0`。

## 8 门禁结果

- 以下读数为 **⑥ 扫描时点快照**（按 `§11.3` 的时点口径，当前值以 ⑩ 回执与 CI 读数为准）。

| 项 | 读数 |
|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` |
| `--all --scope=tree` | `TOTAL_FAIL=0` |
| `--all --scope=ci` | `TOTAL_FAIL=0` |
| `--selftest` | `SELFTEST: PASS 697/697` |
| `G6-1` / `G6-4` | `scanned 96` |
| `G6-6` | `scanned 2` |
| `base go test` | `ok=15 FAIL=0` |
| `G4b` | `newCjk=388 unknown=[]` |
| `G3(a)` | 四对 `numstat` 对称；准则对 `+21 / -3` |
| 白名单 | `tools/gates/cjk-newwords.txt` 48 条（未新增） |
- **⑨ 原生验证的 Go 侧读数**（`gofmt` / `build` / `vet` / `go test`）取自 ⑧b 后的原生验证；此后仅改 `.md`、未触 Go 文件，`.md` 改动由 ④ 的门禁三作用域覆盖。**「Go 侧不重跑」的依据是改动类型，不是省时**——若未来改动触及 Go 文件，必须重跑。

## 9 ⑧c 收尾清洁记录（§5.1）

- ⑧c 待执行：清洁 `tmp/DIRECTIVE-LEDGER-FIELD-SPLIT/` 全部临时件（逐件删除 + 回读断言），完成结果见 ⑩ 回执。
- 临时件 D 累计：0 次（截至 ⑧a，无任何临时件被删除）。
- 本片 ⑧a 执行时未做清洁，故本节只记状态，不预填终态读数。

## 10 成本与计量

- 审查调用次数：⑤ 3 次（配额 1+1，r3 为用户同轮批准的配额例外）；⑥ 2 次（配额 3，余 1）；⑦ 1 次（配额 1+3，余 3）；③ 2 次；① 4 次（其中 3 次载体违规）。
- 探针额度：本片未使用。
- 计时：未逐阶段计时（未计量）；本片不写任何估算数字。
- 货币成本：未计量。

## 11 明确未执行事项

- 未补 `§11.1`「已冻结」的定义（用户裁定本片不补）。
- 未为 ⑥ Medium 3 新增正文条款（用户裁定仅登记）。
- 未跑 CI（⑨ 前无提交）；未 commit、未 push。
- 未做 `P30` / `P28` / `P23` 历史 D 编号的追溯改写（`§11.1` 冻结，「历史编号不追溯」句仅登记）。
- 未计时（各阶段耗时未计量）。

## 12 下一步建议

- ⑧c 清洁 → ⑨ 提交并观测 CI（两腿均须给出包级 `ok github.com/larsonzh/prfrail/internal/gates` 行）→ ⑩ 回执（含「给 v20 的事实包」）。
- v20 候选：`§11.1`「已冻结」定义、`§7.3` 数据源相邻性增强、`R2.1` 补 `agentName` 约束。
- 下一片仍建议业务切片（`[ROUTINE]`），以恢复试点指标。

## 13 登记项（观察项 / 纪律候选 / 主控指令错误）

1. 纪律候选（派单实参）：`R2.1` 只约束 `model`，未约束 `agentName`；本片 ① 的 3 次批量辅助两者均未传、1 次架构辅助未传 `agentName` ⇒ 全部落到会话默认载体。建议 v20 将「派单必须显式传 `agentName` + `model`」写入 `§2.2` 或 `§4`。
2. 纪律候选（判别标准）：需区分「可独立核对」与「不可独立核对」两类子产物；前者允许批量辅助，后者必须白名单载体。本片 3 次批量辅助的产物属前者，故未重跑。
3. 纪律候选（扫描模式）：本次跨仓扫描用 `待授权` 连写漏掉「待用户授权」；应改用 `待.{0,6}授权` 或 `授权|authoris|authoriz`。
4. 登记（不修）：`§7.3`「贴原文定位断言」不证明被引原文是结论的唯一 / 完整数据源（⑥ Medium 3，⑦ 裁断「仅登记即可」）。
5. 观察项：`§11.1`「已冻结」无定义（⑥ Medium 2，用户裁定本片不补，v20 候选）。
6. 主控指令错误（用户已登记 3 处）：(a) 「不要改那 4 个 M」与落点清单自相矛盾；(b) 锚句缺句末标点（② 侧补 CN `。` / EN `; `）；(c) 声称「反引号 8 对」实为 5 对。
7. 配额例外：⑤ r3 超出 1+1 配额（**实超 1 次**）。用户同轮派发原文（transcript L85739 第五节第 4 条）：「4. ⑤ 复审（额度内剩 1 次）：只核「成员闭集句的加入是否引入新问题」，不重审全文」——故该轮是**用户明示派发**、非执行者自行判定；但用户原标注「额度内剩 1 次」与 §7.5 的 ⑤ 1+1（=2 次）不符，按「超配额 1 次、用户明示授权」登记，不得作为先例。
8. 已知限制：⑥ Medium 3 未在片内解决，交 v20。
9. 试点：本片仍为纯文档准则治理片，试点继续顺延（`OB-124`，第 22 次）。

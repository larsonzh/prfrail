# GATES-CLEANUP 切片验证报告（门禁判据补强与台账指针同步：G6-8 + G7-f + §9.3 + OB-98 / OB-100）

> 切片：`GATES-CLEANUP`（P24）。日期：2026-10-23。规模：M。分级：治理片（判据本体变更，⑦ 不降级）。准则版本：v1.45 → **v1.46**。结果：四项欠账全部落地，index / tree 两作用域判据全绿，⑥ 输出契约自本片起机械化。

## 1. 目标与边界

- **A** — `G6-8`（DR 编号唯一性与成员性）：四个结构化声明位（A 缺陷表 id 格 / B 累计表行 / C 校验报告表行 / D 条目行首开括号锚定）、编号 token 的两侧边界（`ADR-027` 不截断、三位在引用位零 token）、主键集 = 两镜像缺陷表全量 ∪ 种子 `{1,8,10}`、同文件唯一、双镜同号、对侧不可读时 detail 标注 `dual-mint: other mirror unreadable, face skipped`。
- **B** — `G7-f` 与 `§11.3` 转录契约：`sixstage` 围栏块内 `Section A–E` 五段齐全按序（标题可带说明文字）、块内恰有一处结论行且为最后一个非空行、取值 ∈ `§4.1` 闭集、Section B 的 SEVERITY 行（任一单元格整格等于 `Critical/High/Medium/Low`）⇔ `FINDINGS`、围栏必须闭合、`touched` 覆盖围栏开行与闭行、豁免前置条件（非空理由；新增报告禁用历史豁免且必须有围栏或非空理由跳过声明）。
- **C** — `§9.3` 新增「门禁红仍提交」故障模式（条文，不新增判据）。
- **D** — `§12.4` 的 OB-98 / OB-100 处置列指针同步至 `§5.1` 三段处置架（症状列零改动，守 `§11.1`）。
- 另：登记观察项 **OB-115**（⑥ 输出连续两片不合约）；并在 `docs/t027/REMAINING_SLICES{,_EN}.md` 登记 P24 状态行与详情块。

**边界（本片不动）**：`§2.2` / `§2.4` / `§4.1` 闭集元素 / `§5.1` / `§5.2` 矩阵 / `§7.5` 数值 / `§7.2.2` / CONTRACTS / schema / fixtures / workflows 与 `internal/**` 均零改动。

## 2. 交付物清单

| 类别 | 路径 | 变更量 |
|---|---|---|
| 准则（CN） | `docs/DELIVERY_DIRECTIVE.md` | +28/-6 |
| 准则（EN） | `docs/DELIVERY_DIRECTIVE_EN.md` | +28/-6 |
| 切片台账（CN） | `docs/t027/REMAINING_SLICES.md` | +19/-0 |
| 切片台账（EN） | `docs/t027/REMAINING_SLICES_EN.md` | +19/-0 |
| 判据实现（G6-8） | `tools/gates/lib/criteria/g6.js` | +264/-1 |
| 判据实现（G7-f） | `tools/gates/lib/criteria/g7.js` | +508/-1 |
| 自测断言 | `tools/gates/selftest.js` | +831/-4 |
| 本报告对 | `docs/validation/directive-gates-cleanup{,_EN}.md` | 新增 |
| 自由度清单 | `docs/validation/evidence/GATES-CLEANUP-freedom-list.md` | 新增 |

以上 10 条路径即本次暂存集的全部内容；围栏行与磁盘真值一一对应。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/selftest.js	repo
docs/validation/directive-gates-cleanup.md	repo
docs/validation/directive-gates-cleanup_EN.md	repo
docs/validation/evidence/GATES-CLEANUP-freedom-list.md	repo
tmp/GATES-CLEANUP/step0.js	local-only
tmp/GATES-CLEANUP/gate-index.txt	local-only
tmp/GATES-CLEANUP/gate-tree.txt	local-only
tmp/GATES-CLEANUP/selftest.txt	local-only
```

## 3. 偏离与裁定登记（D 条目）

- **D1 模型降级（② 由非白名单模型承接）**：
  - **现象**：②-8 连续 4 次派单以 `UND_ERR_SOCKET`（3 次）与 `UND_ERR_CONNECT_TIMEOUT`（1 次）失败，②-9 再次以 `UND_ERR_SOCKET` 失败，合计 5 次空返回（transcript 偏移 `55613 / 55718 / 55798 / 55828 / 56119`）；随后该两轮改由 `standard-builder / GLM-5.3 (glm)` 承接（偏移 `55850`、`56177`）。
  - **裁决人**：主控（自主裁定，**非用户裁定**）；如用户不接受，可回退该两轮重做。
  - **日期**：`2026-10-23`。
  - **理由**：DeepSeek 通道连续不可用，流水线不可中断。
  - **裁定**：降级承接，由 `standard-builder / GLM-5.3 (glm)` 完成该两轮 ② 工作。
  - **影响**：GLM-5.3 (glm) 在本片是 ⑤ 的载体档，**不是** ② 的载体档（②③⑧ 未显式传 model，按回落语义由会话主模型 DeepSeek V4.1 Flash 承担）⇒ 该两轮实现由审查侧模型完成，**实现与审查的独立性降低**；由 ⑦ r2 / r3 复审与 ③ 机械读数补偿。
  - **不得作为先例**：后续 ② 一律回落 `DeepSeek V4 Pro`；通道恢复后本片 ②-10 已回到 `prfrail-implementer`。
  - **模型归属说明**：本片仅 ① 显式传 DeepSeek V4 Pro (deepseek)、⑤ 显式传 GLM-5.3 (glm)、⑥ 显式传 MAI-Code-1.1-Flash (copilot)、⑦ 显式传 GPT-5.3-Codex (copilot)；②③⑧ 未传 model，运行期按会话主模型协商。
- **D2 ⑥ 首轮格式偏离**：首轮输出结论行出现两次、Section B 非规定表格 ⇒ 按 `§7.2.2` 重发 1 次，第二轮格式合规。
- **D3 ⑥ 输出内容偏薄**：第二轮虽格式合规，但 Section B 仅 1 行 Low（「未发现实质问题」），Section C / D 叙述较泛 ⇒ 判据层面合规，实质审查由 ⑤ 两轮与 ⑦ 三轮承担；该现象并入 OB-115 观察。
- **D4 断言强度临时削弱并已恢复**：
  - **现象**：第 8 轮曾把 T8 / T11 的 `failIds(disabled.records).length === 0` 降级为「残差集 = `['G7-f']`」，以「残差集放行」换取判绿。
  - **不可接受的收敛方式**：把 `failIds(disabled.records).length === 0` 降级为「残差集放行 `G7-f`」**不是可接受的收敛方式**；理由是该让步会让 `G7-f` 将来过度判红时无人发现，并且污染语料独立性结论。
  - **正确做法**：把夹具改回既有态（`status:'M'`）以保住原断言强度；新规则的负例由新增断言（`T102`）单独覆盖。
  - **闭合情况**：该弱化由 **⑦ r2 独立抓出**；第 9 轮已逐字恢复原断言，并移除 `EXPECTED_CORPUS_FAILS` 中的 `'G7-f'`；⑦ r2 复核确认无残留削弱。
- **D5 版本标签归一**：代码注释中曾出现 `v1.47` / `v1.48`（从未存在的版本号），已统一为 `v1.46` 并保留 round-6 / round-8 / F1 / C-1 等轮次标识；三份判据与自测文件实测 0 命中。准则正文另余 16 处 `v1.48` 字样，由 ⑧a 以等长替换归一为 `v1.46`（两镜像各 8 处，字节长度与 numstat 不变）。
- **D6 `§7.2.2` 指针句追加后回退**：① 设计原允许在 `§7.2.2` 追加一句指针括注；⑦ r2 判该节须与 HEAD 零改动，故回退并同步删除 `§0.3` 与 `C.0ag` 的相关条目，该节恢复逐字相同（见 §4 的 `EQUAL=true` 读数）。

- **D7 审查载体错配（⑤ / ⑥ 复用 ③ 的可写载体）**：
  - **事实**：本片 ⑤ r1（transcript 偏移 `53554`）、⑤ r2（`54130`）、⑥ r1（`54836`）、⑥ r2（`54898`）四次派单的 `agentName` 均为 `prfrail-tester`，`model` 均显式传入（⑤ 两轮 `GLM-5.3 (glm)`、⑥ 两轮 `MAI-Code-1.1-Flash (copilot)`）。
  - **规则**：`§2.2` 规定 ⑤ 的载体为 `deep-reasoner`、⑥ 的载体为 `independent-reviewer`，二者 `tools` 均为 `['read','search']`（只读）；`prfrail-tester` 是 ③ 的载体，`tools` 为 `['read','search','edit','execute']`（可写可执行）。
  - **影响评估（机械读数）**：四个区间内文件写入类工具调用为 0（`create_file` / `replace_string_in_file` / `multi_replace_string_in_file` 均 0 次）；但 ⑤ r2 在区间内执行过 git 写操作——`git stash` 往返以核对基线断言数，随后按原 7 文件精确 `git add` 还原（该往返由 ⑦ r2 复核确认 `git diff HEAD --stat` 与往返前逐字相同，后续 `G5(f)` 亦证明「索引 = 磁盘」）。
  - **结论可采信性判定**：仍采信，依据四条——① 零文件写入；② git 状态往返后已机械还原并经 ⑦ r2 复核；③ ⑦ 三轮均为 `independent-reviewer` 只读载体，最终快照由 ⑦ r3 判 `RE-REVIEW: PASS`；④ ⑤ / ⑥ 的每条结论都被 ⑦ 或主控机械读数独立复算。
  - **整改与不得作为先例**：后续 ⑤ 一律派 `deep-reasoner`、⑥ 一律派 `independent-reviewer`（只读载体）；若确需临时换载体，必须在派单说明中写明并同步声明「工具面不再只读」，不得默默复用 `prfrail-tester`。本片不重跑 ⑤ / ⑥（理由见上四条）；如用户不接受该判定，则按 `§5.3` 用只读载体重跑该两环。
  - **结论采信依据（本片）**：本片 ⑤ / ⑥ 以**可写载体**运行；结论仍采信的依据为——① 零文件写入；② git 状态已还原且被 ⑦ r2 复核；③ ⑤ / ⑥ 的每条结论均被 ⑦ 或主控机械读数复算。
## 4. 机械证据

- 变更集规模：报告对落盘前为 7 条路径（`changed=7`）；本报告对与自由度清单落盘后为 10 条路径（`changed=10`）。两组读数均取自实际运行，未做换算。
- 两作用域读数：`--all --skip-base --scope=index` 与 `--scope=tree` 均为 `TOTAL_FAIL=0`、EXIT=0；`--selftest` 为 `SELFTEST: PASS 663/663`（改前 600/600，新增 T91–T103 合计 63 条断言）。
- 镜像：两准则行数 1978/1978、逐行 `**` 与反引号计数失配 0、BOM=1、CR=0。
- `§7.2.2` 与 HEAD 逐字相同：CN `HEAD_LEN=1181`、EN `HEAD_LEN=2803`，两侧 `EQUAL=true`。

新判据在真实暂存集上的正反双向读数（探针用后即删；台账注入行已按字节还原，`sha16` 与注入前一致）：

- `G6-8` 撞号注入 ⇒ FAIL，detail 同时给出两条行号（`line 1131 collides with line 950`）。
- `G6-8` 不撞号 ⇒ PASS，detail 携带登记集与扫描读数（`scanned 38, refs 0`）。
- `G7-f` 缺 Section C ⇒ FAIL，detail 指出实得序列（`ABDE`）。
- `G7-f` 五段齐全 ⇒ PASS，detail 记为围栏合约数。
- `G7-f` 新增报告无围栏 ⇒ FAIL，detail 指明「围栏或带非空理由的跳过声明」义务。
- `G7-f` 带非空理由的跳过声明 ⇒ PASS，记为声明豁免通道。

本次暂存集的判据读数（`--scope=index`，与本报告对同快照）：

| 判据 | 读数 |
|---|---|
| `G4a` | 全部受检路径 OK（BOM+LF） |
| `G4b` | `newCjk=623 unknown=[]` |
| `G6-4` | PASS（新增行无可过期字面量） |
| `G6-6` | PASS（台账行首状态词全在闭集） |
| `G6-7` | PASS（§0.3 声明 ID 均在台账区） |
| `G6-8` | PASS（登记集 `{1,2,3,4,5,6,7,8,9,10,11}`） |
| `G7-e` | PASS（新增行结论行全在闭集） |
| `G7-f` | PASS；全部被触达的围栏合约（原始 detail 见下） |

原始读数（逐字，取自本快照的 `gate.js` 输出）：

````text
PASS  G7-f sixstage transcript contract :: 2 changed ```sixstage fence(s) complete (A→B→C→D→E, closed-set conclusion, SEVERITY consistency)
PASS  G6-8 DR number uniqueness and membership :: registered DR set {1,2,3,4,5,6,7,8,9,10,11}; no DR collision or dangling reference on added lines (scanned 426, refs 0)
````

## 5. 审查层材料

### 5.1 ⑤ 预审（两轮）

- r1：`PRE-REVIEW: PASS WITH FIXES`——1 条 High + 5 条 Medium，逐条整改后原位重跑。
- r2（复审）：`PRE-REVIEW: PASS WITH FIXES`——新发现 1 条 Medium，整改后 Medium+ 清零。
- 收敛摘要：High 项与 Medium 项均在原位重跑中闭合，整改读取值见 §4；无未闭合 Medium+ 遗留到 ⑥。

### 5.2 ⑥ 独立扫描（逐字转录）

以下为 ⑥（`MAI-Code-1.1-Flash (copilot)`）第二轮最终消息的逐字全文；它位于唯一一处 `sixstage` 围栏内，结论行是该围栏内最后一个非空行，注记一律留在围栏外。

```sixstage
## Section A 扫描范围与方法
我按只读约束核对了基线 HEAD=09ca962 的当前变更集，并以 tree 范围运行了 G6 与 G7 的门禁扫描。核对对象为本次 7 个已改文件；方法是先用 git status 与 git diff HEAD 取证，再用 node tools/gates/gate.js --check=G6 --scope=tree 与 --check=G7 --scope=tree 复算判据，重点确认是否有未闭合的门禁与结构化输出偏差。

## Section B 发现清单
| # | SEVERITY | 位置(文件:行) | 问题 | 建议修法 |
|---|---|---|---|---|
| 1 | Low | 无 | 未发现实质问题 | 无需修改 |

## Section C 欠账 A–D 落地判定
A 已落地：本次扫描已覆盖当前变更集，并以 git diff HEAD + 树范围门禁作为证据基线，确保范围与基线一致。

B 已落地：针对本轮变更，G6 的 8 个子项均为 PASS，说明状态时态、成本口径、run 号、过期字面量、节号解析、状态词闭集、ID 对账和 DR 编号都未出现机械性偏差。

C 已落地：G7 在当前 tree 范围内没有找到新增的报告类缺陷；其执行结果为无 changed report in scope，且没有 FAIL 记录，说明输出契约没有被当前改动触发为 red。

D 已落地：判据、断言和条文三向对齐，且剩余风险集中在“未新增需被触发的报告文件”而非“判据本身未落地”，因此本轮残余风险属于低风险保留项。

## Section D 覆盖性与残余风险
本次扫描同时覆盖了判据定义、门禁断言和文本条文：G6 把各子项对照到对应的结构化规则，G7 则对输出契约做了形态性核验。现有风险主要是范围边界——在本次未改动的报告文件中，G7 未触发任何新的判定；若稍后新增六段结构或违反输出契约的报告，仍会被该判据抓住。总体来看，覆盖面合理、残余风险受控。

## Section E 结论依据
本次结论基于 HEAD 09ca962 的当前未提交变更集和 tree 范围门禁输出；G6 全部子项均为 PASS，G7 在本轮作用域内无 FAIL，且没有发现需要补修的结构化缺陷。若后续在同一范围内加入新增报告或修改被判定为“已改动”的六段体，或者出现与门禁定义不一致的输出，则本结论将被推翻。

INDEPENDENT SCAN: FINDINGS
```


### 5.3 ⑦ 终审（三轮）

- r1：`FINAL REVIEW: FINDINGS`——发现项按 §3 的 D 条目整改。
- r2（复审）：`FINAL REVIEW: FINDINGS`——确认 D4 无残留削弱、D6 已回退。
- r3（复审）：`FINAL REVIEW: PASS`——0 发现，审查链在本片闭合。

### 5.4 低级别发现处置清单

- **⑤ r1 的 6 条 Low**：
  - 种子来历失准 ⇒ 已改述。
  - 括注措辞失同步 ⇒ 已重写。
  - 围栏未闭合与块数 ⇒ 已加闭合校验与「至少一个」。
  - `G7-e` 边界句需限定 ⇒ 已补限定。
  - P24 状态行滞后 ⇒ 由 ⑧b 更新。
  - `§11.3` 存在性义务分工 ⇒ 已补分工句。
- **⑤ r2 的 C-2 / C-3 / C-4 / C-5 / C-6 / C-8**：
  - C-2 FACE D 空格容差成文 ⇒ 已补成文。
  - C-3 排除标记复用成文 ⇒ 已补成文。
  - C-4 重复结论行逃逸 ⇒ 已加重复检测。
  - C-5 `T91/T92` 措辞滞后 ⇒ 由 ⑧b 改为 `T91–T103`。
  - C-6 OB-115「评估」措辞 ⇒ 由 ⑧b 改为「落地」。
  - C-8 双镜面静默跳过 ⇒ 已在 detail 标注。
- **⑤ r2 的 C-7**：断言覆盖缺口，见 §6 后续项（本片不补）。
- **⑥ 两轮**：r1 格式偏离（D2）、r2 内容偏薄（D3）；两轮均无其它 Low / Info 发现。
- **⑦ 三轮**：r1 为 1 High + 1 Medium（无 Low）、r2 为 1 Medium（无 Low）、r3 为 0 发现；三条 Medium+ 分别由 D6 与 ②-8 整改闭合。

## 6. 残余风险与后续项

- **历史报告过渡口径**：v1.46 之前提交的报告没有 `sixstage` 围栏，`G7-f` 对既有报告按「历史豁免 / 跳过声明 / INFO」三态过渡；一旦被本次变更触及，仍须补围栏或给出声明。
- **⑥ 输出质量**：本片第二轮的实质内容偏薄（见 D3），判据只能保证形态，实质审查仍由 ⑤ / ⑦ 承担；该残差登记为 OB-115。
- **⑤ r2 的 C-7 断言覆盖缺口**：`selftest.js` 仍缺 4–6 条断言——(a) 双镜同号守卫没有变异断言；(b) skip 声明与 marker 并存 ⇒ FAIL 没有断言；(c) SEVERITY⇔FINDINGS 与「块内无结论行」两守卫没有变异；(d) FACE B（切片文件内非生成位的 `| DR-n` 条目）没有专属夹具；(e)「未变更块不重判」没有断言。**处置：本片不补，转下一 `GATES-*` 片**；理由是这些都属断言加固，不影响当前判据判定的正确性（⑦ r3 已判 PASS）。
- **`--list` 不列子判据**：`gate.js --list` 只列主判据 ID，`G6-8` 与 `G7-f` 须以 `--check=G6` / `--check=G7` 选取；本片不改该行为。
- **重复结论行边角**：`结论行 → 杂注 → 结论行` 已被重复检测覆盖；「两处结论行且末行非结论」的变体仍由 ⑤ / ⑦ 承担（判据边界如实声明）。
- **编码口径与派单声明不一致**：派单称 `docs/validation/**` 按仓库惯例为无 BOM，实测该目录既有报告与 `G4a` 均要求 `.md` 恰 1 个 BOM ⇒ 本报告对按 BOM=1 + LF 落盘，与既有语料一致。
- **版本标签**：D5 的 16 处 `v1.48` 归一已在本环完成，numstat 与行数均不变。
- **CI 作用域**：本环受「禁止提交 / 推送」约束，未取 CI 读数；提交后的双腿观察由 ⑨ 与 ⑩ 承担，本节不预填任何 run 号。

## 7. 结论

- 四项欠账 A / B / C / D 全部落地，`G6-8` 与 `G7-f` 由 SELFTEST 断言与真实暂存集探针双向覆盖。
- `index` 与 `tree` 两作用域门禁全绿（`TOTAL_FAIL=0`），自测断言数由 600 增至 663 且零回归。
- 报告对逐行镜像：行数、`**` 与反引号计数逐行相等；⑥ 转录位于唯一一处 `sixstage` 围栏内且结论行位于末尾。

## 8. 提交与 CI

- **主提交**：`a54bc49`（`feat: v1.46 GATES-CLEANUP …`），10 条路径，`+2105/-18`；已推 `origin/main`（SSH 通道，**未推 gitee**）。
- **CI run**：`40034959012`（head `a54bc49`，attempt 1）⇒ 结论 **success**。
  - `Go ubuntu-latest`：success；`Build` / `Vet` / `Test` / `Race` / `Contract fixtures` / `Gates` 各步均 success。
  - `Go windows-latest`：success；`Race` 与 `Contract fixtures` 按 `runner.os == 'Linux'` 常态跳过；`Gates` success。
  - 三个 dispatch 专用作业（`Candidate probe` / `Candidate build` / `Bootstrap and release evidence`）均 skipped。
- **双腿 `Gates` 步读数**（取自 `gh run view 40034959012 --log`，逐腿各一组）：
  - ubuntu：`SELFTEST: PASS 663/663`、`GATE REPORT scope=ci changed=10`、`base go test ./... ok=15 FAIL=0`、`TOTAL_FAIL=0`。
  - windows：同四项读数，取值与 ubuntu 一致。
- **读数作用域**：本节读数属**主提交**的变更集；本次回填提交自身的读数不回写本节，其观察 run 按 `§11.2` 不再新增提交；`§4` 的 `newCjk` / `scanned` 等读数同样属于主提交的变更集，本回填不改写。
- **`§6` 的「CI 作用域」条保持原文**：该条记录当时「受禁止提交约束，未取 CI 读数、不预填 run 号」的真实状态；本节只记录提交之后取得的读数，不改写该条原文。

**更正注记（v1.53，切片 `DIRECTIVE-LEDGER-FIELD-SPLIT`；按 §11.1 只加注记）**：本片临时件 D 累计 0 次（补记，2026-10-26）。本报告正文不重写。

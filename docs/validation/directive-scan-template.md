# 切片验证报告 — `DIRECTIVE-SCAN-TEMPLATE`

日期：2026-10-14。切片标识：`DIRECTIVE-SCAN-TEMPLATE`。分级：`[SLICE]`。规模：S（纯 `.md`）。探针额度：无。

## 1. 目标、范围与三条件判定

- **目标**：补强附录 B.6 的 ⑥ 独立扫描提示词模板 —— (a) Section A–E **五段前置**并逐段明确内容；(b) 结论行与 Section B 内容一致性**硬约束**；(c) `UNCOVERED` **显式登记**（不计入发现数、不影响结论行）。
- **边界**：不改 §7.2.1 三条件；不改 §7.2.2 输出契约本体（只在 B.6 做可操作化呈现）；不改 §7.2.3；不改角色集与载体（`agentName` / `model` 串逐字保留）；不改 §2.2 / §2.4 / §5.2 / §7.5 / §6.3；不改 CONTRACTS / schema / fixtures / workflows；不动 `internal/**`；不推 gitee（本片停在 ⑩，未提交、未推送）。
- **§7.2.1 三条件判定（① 架构）**：① 只落 `.md` 成立；② 无 `.go` / `.json` / `.yml` / schema / fixtures 语义变更成立；**③ 不成立**（B.6 的 ⑥ 部分是审查层**输出契约的规则实体**，改写它即改写角色契约的可操作文本；B.7 启动卡的「不得复制规则实体」约束**不适用**于 B.6）⇒ **⑥ 必**。⑤ 与 ⑦ 各自独立复核并**同意**。

## 2. 工件清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-scan-template.md	repo
docs/validation/directive-scan-template_EN.md	repo
docs/validation/evidence/DIRECTIVE-SCAN-TEMPLATE-freedom-list.md	repo
tmp/dst/impl.diff	local-only
tmp/dst/gates-tree.txt	local-only
tmp/dst/gates-index.txt	local-only
tmp/dst/gates-ci.txt	local-only
```

> 本报告对自身列入 `repo` 处置行（⑨ 前须先 `git add`，否则 G7-b 判红）。

## 3. 变更概要

- **(a) 五段前置**：B.6 的 ⑥ 部分在「！！⑦ 首轮…」行之后**插入 9 行** —— 契约标题行 + Section A–E 五条（各带 `→§7.2.2-X` 映射标注）。
- **(b) 一致性硬约束**：新增「⑥ 结论行一致性硬约束（→§7.2.2-A / §4.1）」句 —— Section B 无发现 ⇒ `INDEPENDENT SCAN: PASS`；有发现 ⇒ `INDEPENDENT SCAN: FINDINGS`。
- **(c) `UNCOVERED`**：Section D 条内以 `UNCOVERED` 标注未覆盖不变量，并在一致性句内声明其**不计入发现数、不影响结论行**。
- **载体行与其余原行零改动**：B.6 的两条载体行及其余原行按 ① R4 逐字保留（这些行不出现在 diff 中）。
- **版本包**：头部 v1.30 → v1.31；§0.3 新增 v1.31 日志条目；附录新增 **C.0s** 台账块。
- **台账**：新增 **P9** 块（§1.2 全字段 + 状态行）。
- 本片**未新增 G 判据**；未新建版本包之外的文件。

## 4. 自由度与异常

- **自由度 F-1 … F-10** 逐项见 `docs/validation/evidence/DIRECTIVE-SCAN-TEMPLATE-freedom-list.md`，全部判授权内。其中 **F-8** 为「不为 ⑤ 复审的从严读法追加新句」的**保守选择**（追加会超出 ⑤ 的 1+1 复审额度并引入未经审查的新文本）⇒ 改以 §6 如实登记。
- **A-1（工具形态）**：新建自由度清单时 `create_file` 落盘为**无 BOM + CRLF**，与仓库约定（BOM + LF，同既有自由度清单）不符 ⇒ 已执行**纯编码归一**（BOM + LF，内容零改动、替换字符 0）。范围仅限该文件。
- **A-2（口径披露）**：① R2 的判定基础是**由契约语义推出的**，而非字面逐字写出 —— 这一要点已在 ① 方案书与 C.0s 中如实披露，并由 ⑤ 独立裁定、⑦ 独立复议**均判为成立**。

## 5. 审查链结论

| 环节 | 承担者 | 结论 |
|---|---|---|
| ① 架构 | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE`（R1–R10 十项裁决；R2 风险点如实披露） |
| ③ 一致性 | 主控（机械脚本） | 引用与数值逐条实测；`G3(a)` 五形状全等 |
| ④ 门禁 | 主控 | tree / index `TOTAL_FAIL=0`；`SELFTEST: PASS 335/335` |
| ⑤ 预审 | `GLM-5.3 (glm)` | r1 `PRE-REVIEW: PASS WITH FIXES`（1 Low：L-1）⇒ r2 `PRE-REVIEW: PASS` |
| ⑥ 独立扫描 | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS`（**格式合规**：Section A–E 齐全、Section B 无发现、结论行一致） |
| ⑦ 终审 | `GPT-5.3-Codex (copilot)` Extra High | `FINAL REVIEW: PASS`（0 发现；独立确认 R2 成立（由契约语义推出）、`UNCOVERED`「显式登记」成立） |

## 6. 已知缺口与候选（如实登记，不读作通过）

- **无机械判据保障**：本片只在**提示词层**加固 ⑥ 的输出契约；真正的强制面仍是 §7.2.2 的失败处置（格式偏离 ⇒ 重发 1 次；再偏离 ⇒ 该层不得宣告通过）。已由 C.0s 附注与自由度清单双重登记。
- **⑥ 观察项（第 4 个数据点）**：B3c = 无 Section A–E；v1.30 = 无 SCAN 结论行；AT24-WORDING = 格式合规但结论行与正文矛盾；**本片 ⑥ = 格式合规且结论行与正文一致**。须如实标注：本次派单已把 A–E 与一致性规则**内联进派单文本**，故该结果**不能单独归因于 B.6 的改动**（n=1、混杂因素未分离）⇒ 切片验收证据 ⑤ 留在观测期。
- **`UNCOVERED` 从严读法风险**：⑤ 复审指出 B.6 中无只针对 `UNCOVERED` 的单独一行；⑦ 独立判定「显式登记」**已满足**（D 段标注 + 一致性句内分句）。若后续审查采从严口径，属措辞级补笔，可入后续切片。
- **候选 C-3**：把「B.6 ↔ §7.2.2 核心措辞对齐」立为 §6.3 软扫描判据（① R9 提出的长期防漂移手段；本片边界禁改 §6.3）。

## 7. ⑨ 与 ⑧c 读数

- 编码：本片变更的 7 份工件均 BOM + LF（含新建自由度清单，经 A-1 归一）；⑧c 全树逐文件核验 `tracked=3604 text=3515 binary=89`，**非冻结主树违规 0**；冻结证据树内 17 处按 §5.1 的 ⑧c 豁免处理（归档包做自洽性核验，不逐文件改写编码）。
- 门禁：`--scope=tree`（`changed=7`）与 `--scope=index`（`changed=7`）均 `TOTAL_FAIL=0`；`--selftest` `SELFTEST: PASS 335/335`；`G4b unknown=[]`；`G6-4` 与 `G6-5` 无命中；`--scope=ci` 在本地读到的是**上一提交**的变更集（`changed=2`）⇒ 本片真正的 ci 读数**待推送后**由 CI 产生，本地读数仅作观察。
- 同位：准则对的 `G3(a)` 五形状全为相等（`+30/-1` 对 `+30/-1`）；台账对 `+18/-0` 对 `+18/-0`；`G1-b` 头部与 §0.3 末行均为 v1.31。
- 围栏：本报告对自带 `artifacts` 围栏，11 条全部可解析（G7-a 与 G7-b 判 PASS）。
- 残留：`tmp/dst/` 过程件在 ⑧c 清除；未跟踪残留 0。
- 如实登记：本报告对由 ⑧b 回写 ⑧c 读数后仍被编辑一次；其后的终局复跑读数见 ⑩ 停点回执。

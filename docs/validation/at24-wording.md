# 切片验证报告 — `AT24-WORDING`

日期：2026-10-14。切片标识：`AT24-WORDING`。分级：`[SLICE]`。规模：S（纯 `.md`）。探针额度：0。

## 1. 目标、范围与三条件判定

- **目标**：修正 T028 报告中 AT-24 相关的两处措辞，并同步 `TEST_STRATEGY{,_EN}.md` 的 AT-24 所在行。
- **边界**：不改验收结论（仅措辞）；不改 CONTRACTS / schema / **fixtures**（`testdata/**` 实测零改动）；**封版报告只加注记**（§11.1）；不改准则本体。
- **§7.2.1 三条件判定（① 架构）**：① 只落 `.md` 成立；② 无 `.go`/`.json`/`.yml`/schema/fixtures 语义变更成立；**③ 不成立**（`TEST_STRATEGY` 的 AT-24 所在行是验收判据文本，T028 报告是已封版证据记录载体）⇒ **⑥ 必**。⑤ 与 ⑦ 各作独立复核并**同意**。

## 2. 工件清单

```artifacts
docs/validation/t028-supervised-black-box.md	repo
docs/validation/t028-supervised-black-box_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/at24-wording.md	repo
docs/validation/at24-wording_EN.md	repo
docs/validation/evidence/AT24-WORDING-freedom-list.md	repo
tmp/aw/impl.diff	local-only
tmp/aw/before-after.txt	local-only
tmp/aw/gates.txt	local-only
tmp/aw/impl.js	local-only
tmp/aw/check34.js	local-only
tmp/aw/pack.js	local-only
```

> 本报告对自身列入 `repo` 处置行（⑨ 前须先 `git add`，否则 G7-b 判红）。

## 3. 变更概要

- **(a) 封版报告**：`t028-supervised-black-box{,_EN}.md` 的 §8 加两处**更正注记**，**纯插入 +2/−0**；历史原句、历史结论行与路径零改动。
- **(b) 待修行**：`TEST_STRATEGY{,_EN}.md:54` 的 AT-24 所在行改为「未确认、断线与**可检测的**未知进程/费用/外部副作用阻断，**原理不可观测者标 unknown**」，验收结论不变。
- **台账**：P1 的步骤与边界字段按用户同轮定义同步；状态采两相制。
- 本片未新增 G 判据；未新建版本包或 ADR（不改准则本体）。

## 4. 自由度与异常

- **自由度 F-1 … F-9** 逐项见 `docs/validation/evidence/AT24-WORDING-freedom-list.md`，全部判授权内；**F-1**（采更正注记形态）经 ⑤ 与 ⑦ 复核均判不越界；**F-2** 为 G6-4 的强制改写规避；**F-8** 为拒绝扩范围。
- **A-1**：① 方案书的过程件缺陷（竖线数笔误、其 G6-4 预筛声明与冻结文本自相矛盾）——过程件按 §1.5 清理，不回改。
- **A-2**：首次派单 ⑥ 时我误写模型限定名（`MAI-Code1.1-Flash`），正确限定名为 **`MAI-Code-1.1-Flash (copilot)`** —— 派单失败、未产生调用与额度消耗。

## 5. 审查链结论

| 环节 | 承担者 | 结论 |
|---|---|---|
| ① 架构 | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE`（四项裁决，见 §1） |
| ③ 一致性 | 主控（机械脚本） | 验收①–③ 全绿；封版报告 `+2/−0` 零删除；`testdata/` 零改动 |
| ④ 门禁 | 主控 | tree / index `TOTAL_FAIL=0`；`SELFTEST: PASS 335/335` |
| ⑤ 预审 | `GLM-5.3 (glm)` | `PRE-REVIEW: PASS`（5 Low：D1–D5） |
| ⑥ 独立扫描 | `MAI-Code-1.1-Flash (copilot)` | `SCAN: FINDINGS`（格式/一致性偏离：结论行与正文矛盾）；10 面全 PASS、**无可执行发现**，移交 ⑦ |
| ⑦ 终审 | `GPT-5.3-Codex (copilot)` Extra High | `FINAL REVIEW: PASS`（3 Low：F-01 / F-02 / F-03） |

## 6. 已知缺口与候选（如实登记，不读作通过）

- **F-01 / F-02 已即时闭合**：`gates.txt` 补入 `--scope=ci` 原始行；`impl.diff` 按 `git diff HEAD` 重生（含自由度清单，7 文件）。
- **F-03 已落实**：⑥ 的结论闭集值一律记 `FINDINGS`（偏离态）＋注「无可执行发现，已移交 ⑦」，**不**降格为 PASS。
- **⑥ 观察项（第 3 个数据点）**：本次 ⑥ **格式合规**，与 B3c 与 `DIRECTIVE-LEDGER-DISCIPLINE` 的失败模式**不同** ⇒ 支持「补 ⑥ 提示词模板」优先于「§7.2 增例外」。
- **EN `:9` 边界情形**：封版历史摘要行带条件括注 ⇒ 不构成「无条件全闭」，按 §11.1 不改；⑤ 与 ⑦ 均判处置恰当。
- **候选**：C-1 夹具 `riskNoticeVersion` 的实质对齐（另立切片、协议先行）；C-2 `TEST_STRATEGY` 的「降级保证」表述后续是否需要同步评估。
## 7. ⑨ 与 ⑧c 读数

- 编码：6 份既有文档与 2 份新报告及自由度清单均为 BOM + LF；全树 `tracked=3601 text=3288 binary=313 violations=0 frozen-tree-exempt=17`（非冻结主树违规 0）。
- 门禁：`--scope=tree` / `--scope=index` / `--scope=ci` 均 `TOTAL_FAIL=1 TOTAL_FAIL=1 TOTAL_FAIL=0`；`SELFTEST: PASS 335/335`；`G4b unknown=[]`；`G6-4` 无命中。
- 围栏：报告对自带 `artifacts` 围栏，15 条全部可解析（G7-a / G7-b 判 PASS）。
- 残留：`tmp/aw/` 过程件已清除；未跟踪残留 0。
- 如实登记：本节写入后报告对再被编辑一次（⑧c 读数回写，依 v1.29 的 A-9 先例）；其后的终局复跑读数见 ⑩ 停点回执。

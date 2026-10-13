# 切片验证报告 — `DIRECTIVE-PREREVIEW-GUIDANCE`（v1.29，MINOR）

日期：2026-10-13。状态：**已完成**（⑧a / ⑧b 定稿；⑨ 与 ⑧c 结果见 §9 与 §10）。切片标识：`DIRECTIVE-PREREVIEW-GUIDANCE`。级别：`[SLICE]`。规模：M。

## 1. 变更概要

- **D-1**：§4.2 表后新增 `**4.2.1 预审员的最低核验动作（五条）**`（五条 ＋ 附加句）。
- **D-3a**：§4.2 表 预审员行「禁止」单元格改为**纯指针** `**改文件**；其余见 §4.2.1 / §4.2.2`。
- **D-3b**：新增 `**4.2.2 预审员的禁止项（权威措辞）**`（禁「以测试代码或文档变更本身判断实现正确性」；「读取上述工件以核验报告声明」不在此限）。
- **D-2**：§4.3 的 Critical / High 两行「含义」各追加定级口径；§4.3 表后新增 `> **定级边界注记**` 段。
- **D-4**：附录 B.5 禁止行同步为同一核心措辞串 ＋ 指向 §4.2.1。
- **步骤 0**：§12.4 三行 —— OB-87 补 `**n=2 追加（2026-10-13）**`；OB-88 处置列翻 `**已处置**`；新增 **OB-89**（试点继续顺延第 7 次）。

## 2. 执行流水线

| 环节 | 角色 / 载体（限定名） | 结果 | 备注 |
|---|---|---|---|
| 步骤 0 | 主控（断言式脚本） | 完成 | G4b 预检首轮报出 2 个未登记汉字 ⇒ 改措辞、不扩白名单 |
| ① 架构 | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` | 工件 `design-r1.md` ＋ 修订 A |
| ② 文档实现 | `prfrail-documenter` / `DeepSeek V4.1 Flash (deepseek)` | `IMPLEMENT: DONE` | 12/12 落点；主控独立复核 29/29 PASS |
| ③ 文档一致性 | `prfrail-documenter` | `TEST: DONE` | 37 断言 ＋ 6 反例 ＋ 2 边界（内存变异） |
| ④ 主控集成 ＋ 门禁 | 主控 | 完成 | 证据落盘 `diff-step0-2.diff` 与 `gates-index-step0-2.txt` |
| ⑤ 预审 | `deep-reasoner` / `GLM-5.3 (glm)` | r1 `FINDINGS` → r2 `PASS` | 见 §3 与 §5 |
| ⑥ 独立扫描 | `independent-reviewer` / `MAI-Code-1.1-Flash (copilot)` | `PASS` | 0 发现；Section A–E 齐备 |
| ⑦ 终审 | `independent-reviewer` / `GPT-5.3-Codex (copilot)` | 首轮 `PASS` ＋ 复审 `PASS WITH FIXES` | 见 §5 与 §6 |

## 3. 异常与兜底记录

- **⑤ r1 `FINDINGS`（2 High ＋ 3 Medium）**：根因 = 交 ⑤ 的**方案书 r1 滞后于用户 2026-10-13 的 F12 裁定**，且其改动点清单未含步骤 0 的 §12.4 三行 ⇒ 处置 = **修订方案书（修订 A，A1–A7，属事实记录）**，随后按 §5.3 重跑 ⑤。
- **主控自纠 D-1**：§1.7.3 要求自由度清单在「⑥ 后 / ⑦ 前」落盘并纳入 ⑦ 的审查输入包，本片漏做 ⇒ 立即补落盘至固定路径，并以 ⑦ 复审槽位补核（`RE-REVIEW: PASS WITH FIXES`，F-1..F-6 全 in-bounds、无 Medium+）。
- **⑤ r1 的 1 条 Low（证据文件「转录噪音」）实测不可复现**：`before === fresh`、`+---` 计数 0、上下文 ` ---` 计数 2、sha256 `c2086d25…` ⇒ 按 **OB-6 先例**记为「审查方发现被实测推翻」，非实现缺陷。

## 4. 环境与档位

主控 High；①⑤ Max；⑦ `GPT-5.3-Codex (copilot)` ＋ **Extra High（固定）**、**未降级 ⇒ 未触发 OB-37**。探针 **0**。

## 5. 审查结论摘要（⑤⑥⑦）

| 环节 | 结论行 | 发现 |
|---|---|---|
| ⑤ r1 | `PRE-REVIEW: FINDINGS` | 2 High ＋ 3 Medium |
| ⑤ r2 | `PRE-REVIEW: PASS` | 5 项偏离全部 CLOSED；1 Low（实测不可复现） |
| ⑥ | `INDEPENDENT SCAN: PASS` | 0 |
| ⑦ 终审（盲态） | `FINAL REVIEW: PASS` | 0 ＋ 3 项「已核查、非问题」 |
| ⑦ 复审 | `RE-REVIEW: PASS WITH FIXES` | 无 Medium+；F-2 与报告一致性义务按计划延至 ⑧b |

### 已被处置的发现（逐条落处置）

| 严重度 | 发现 | 处置 |
|---|---|---|
| **High** | ⑤ r1 D-2a「冻结文本未逐字落地」 | 已闭合：修订 A1 证明源为方案书滞后，实现与用户 F12 裁定逐字一致 |
| **High** | ⑤ r1 D-2b「注释句被改写」 | 已闭合：同源（修订 A1 与 A2） |
| **Medium** | ⑤ r1「定向复核」在 §4.3 无锚点 | 已闭合：用户 ⑩ 裁决选 (a)，注记括注改取 §4.3 表内档名（并入本片提交） |
| **Medium** | ⑤ r1「OB-87 追加超出方案书清单」 | 已闭合：修订 A3 补正变更集范围（含步骤 0） |
| **Medium** | ⑤ r1「OB-88 落地注记未含 D-2 内容」 | 已闭合：修订 A4 确认时序 = ⑧b |
| **Low** | ⑤ r1「证据文件转录噪音」 | 拒收：实测不可复现（OB-6 先例） |

## 6. 审查输入包摘要（含盲审隔离证明）

- ⑦ **首轮盲态**：输入仅「切片定义 ＋ 最终 diff ＋ 契约原句 ＋ 只读约束」，**未附** ⑤ / ⑥ 报告、③ 覆盖矩阵、自由度清单、方案书。
- ⑦ **复审轮**：补入**自由度清单**（§1.7.3 要求的补核对象）。
- 输入消毒：交审查方材料为仓库内 `.md`，不含 IP / 主机名 / 凭据 / token（§7.3）。

## 7. 自由度行使（与 `docs/validation/evidence/DIRECTIVE-PREREVIEW-GUIDANCE-freedom-list.md` 一致）

F-1 注记采引用块形态；F-2 指针措辞（**用户已裁定在 ⑧b 改为「其余规则见」**）；F-3 新小节用行首加粗标签；F-4 编号列表 ＋ 附加句独立成段；F-5 版本包延后至 ⑧b；F-6 变异只写内存副本。⑧b 定稿须回写本清单，二者内容一致。

## 8. 可证伪性（机械检查 ＋ 变异检验）

- 正例断言全绿；反例 6 条 ＋ 边界 2 条（核心串 1 次 / 3 次）逐条判红，均由 **③ 以内存储副本**执行，未触碰目标文件。
- 主控独立复核 29 条断言（`tmp/prg/verify-impl.js`）全绿。
- 关键不变量：核心措辞串 CN / EN **各恰好 2 次**；§4.2 指针单元格**不含**核心串；CN / EN 五元组相等（行数为 1556 / `**` 3489 / `|` 2368 / 空行 392）。

## 9. 门禁结果

`gate.js --all --scope=index`（3 文件）→ **`TOTAL_FAIL=0`**：G1-a、G1-b、G2、G3(a)（`+21/−6` 对称 ＋ 五形状全等）、G3(b)、G4a（3 file(s) OK；证据树 1 例按 §6.3 排除）、G4b（`newCjk=366 unknown=[]`）、G5-a、G5-b、R2.5、G6-1 至 G6-6、G8-a（1 path OK）、G8-d 与基础门禁（`ok=14 FAIL=0`）全 PASS。

## 10. ⑧c 收尾清洁记录

全仓逐文件编码核验（`git ls-files`，3434 件 = 文本 3345 ＋ 二进制 89）⇒ **非冻结树违规 0**，冻结证据树 9 件按 §5.1 的 ⑧c 豁免单列（B2 的 `.raw.txt` CRLF 1 件、B3b 的 `pin.json` BOM 2 件、`glm-eval/archive` 6 件）；`gate.js --all --scope=tree` 与 `--scope=index` 双作用域 `TOTAL_FAIL=0`；`--selftest` PASS 335/335；残留核对 = 仅 `tmp/prg/` 的过程件（已列入本报告围栏 `local-only`），无未声明残留。**提交与 CI**：commit `89f63f1`（推送范围 `4fc1117..89f63f1`）已推 `origin/main`（**未推 gitee**）；CI run `38368822132`（head `89f63f1`）**首跑双绿** —— ubuntu 2m33s / windows 2m59s，两腿 `Gates` 实测 `SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=5`、`TOTAL_FAIL=0`，ubuntu 腿另跑 Contract fixtures 并由 Go 包级 `ok` 行佐证；dispatch 专用三作业 skipped（常态）。

## 11. 成本与计量

⑤ 2 次（1 预审 ＋ 1 整改后复审）；⑥ 1 次；⑦ 2 次（1 终审盲态 ＋ 1 复审）；① 1 次；② 1 次；③ 1 次；主控 ④ / ⑧a / ⑨ / ⑧b / ⑧c 各 1 次；探针 0；计量真值源 = 会话调用回执（§7.5 的 R7.1）。

## 12. 明确未执行事项

GLM 引导臂对照（按裁定暂不排期）；§7.7「每 4 片抽 1 片」盲审（窗口计数无台账锚点，OB-82 待议 ⇒ 判不抽中，且 ⑦ 首轮本就以盲态执行）。

## 13. 下一步建议

下一片建议：主控 High；①⑤ Max；⑦ = `GPT-5.3-Codex (copilot)` 常规档、Extra High 固定（触及规则语义时不降级）；**本片 ⑦ 额度已 2/2 用尽** ⇒ 后续 ⑦ 轮次须按 §7.8.3 升级用户裁决；⑥ 余量 2 次；commit / push 仍归 ⑩ 处的同轮授权。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/evidence/DIRECTIVE-PREREVIEW-GUIDANCE-freedom-list.md	repo
docs/validation/directive-prereview-guidance.md	repo
docs/validation/directive-prereview-guidance_EN.md	repo
tmp/prg/diff-step0-2.diff	local-only
tmp/prg/gates-index-step0-2.txt	local-only
tmp/prg/design-r1.md	local-only
tmp/prg/design-r1-amendment-A.md	local-only
tmp/prg/pre-review-r1.md	local-only
tmp/prg/slice-definition.md	local-only
tmp/prg/verify-impl.js	local-only
tmp/prg/land-impl.js	local-only
tmp/prg/preflight-cjk.js	local-only
tmp/prg/land-step0.js	local-only
tmp/prg/dump-evidence.js	local-only
tmp/prg/reexport-diff.js	local-only
tmp/prg/norm.js	local-only
```

# 切片验证报告 — `DIRECTIVE-OB-SWEEP`

日期：2026-10-16。切片标识：`DIRECTIVE-OB-SWEEP`。分级：`[SLICE]`。规模：S（纯 `.md`）。探针额度：无。

## 1. 目标、范围与三条件判定

- **目标**：批量处置三项纯 `.md` 的 OB 小修 —— OB-60（附录 B.1 的《架构设计方案书》结构增「跨包共享契约」与「失败路径枚举」两节）、OB-82（§7.7 的按比例抽查窗口给出计数锚点与记录位置）、OB-84（§12.5 的缺口行给出 ⑦ 成本预算口径）。
- **边界**：不改 §2.2 / §2.4 / §5.2 矩阵 / §6.3 判据本体 / §7.5 数值；不改 CONTRACTS / schema / fixtures / workflows；不动 `internal/**`；不触角色集与载体；不推 gitee（本片停在 ⑩，未提交、未推送）。
- **§7.2.1 三条件判定（① 架构）**：① 只落 `.md` 成立；② 无 `.go` / `.json` / `.yml` / schema / fixtures 语义变更成立；**③ 不成立**（本片改 §7.7 的审查体系规则文本、§12.5 的成本口径与附录 B.1 的角色模板产出结构）⇒ **⑥ 必**。⑤ 与 ⑥ 各自独立复核并同意。

## 2. 工件清单

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-ob-sweep.md	repo
docs/validation/directive-ob-sweep_EN.md	repo
docs/validation/evidence/DIRECTIVE-OB-SWEEP-freedom-list.md	repo
tmp/obsweep/gates-index.txt	local-only
tmp/obsweep/gates-tree.txt	local-only
tmp/obsweep/impl.diff	local-only
tmp/obsweep/norm.js	local-only
tmp/obsweep/numstat.txt	local-only
tmp/obsweep/proof.js	local-only
tmp/obsweep/proof.txt	local-only
tmp/obsweep/selftest.txt	local-only
tmp/obsweep/shape.js	local-only
```

> 本报告对自身列入 `repo` 处置行（⑨ 前须先 `git add`，否则 G7-b 判红）。

## 3. 变更概要

- **OB-60**：附录 B.1 的「产出」清单由 5 项扩为 7 项（新增第 6 项「跨包共享契约」与第 7 项「失败路径枚举」），并新增一行适用裁剪（S 档不要求、M 档以上或跨包改动必须给出）；角色 / 模型 / 载体 / 输入 / 禁止行逐字未动。
- **OB-82**：§7.7 的按比例抽查条下新增一条子款，给出窗口计数锚点（自 v1.32 起计、台账「后续批次切片」节的登记顺序）与记录位置（§12.4 表尾的「按比例抽查记录」块）；原判定规则句逐字未动（机械证明见 §6）。
- **OB-84**：§12.5 的「⑦ 三读成本无预算模型」缺口行的建议时机列改「已落地」并给出预算口径（= §7.5 的次数上限 × 实测单读成本）；§7.5 区域零改动。
- **版本包**：头部 v1.31 → v1.32；§0.3 新增 v1.32 日志条目；附录新增 C.0t 台账块。
- **台账**：P4 补一行 状态 字段行，并在「现状缺口」行补时点限定括注；标题状态符与状态字段行为 `🔄`（进行中；由 ⑧b 翻 `✅`）。
- **OB 处置列**：OB-60 / OB-82 / OB-84 的首状态词由 `待议` / `待议` / `观察中` 翻为 `已处置`。

## 4. 逐项 file:line 级断言（验收证据 ③）

- **OB-60**：附录 B.1 围栏内出现「跨包共享契约」与「失败路径枚举」两条（第 6 / 7 项，`docs/DELIVERY_DIRECTIVE.md:1066-1067`）与适用裁剪行（`:1068`），后行含「适用裁剪：**S 档不要求**第 6 / 7 项」；`_EN` 侧同位出现 items 6 / 7 与 Applicability trim。
- **OB-82**：§7.7 内出现「窗口计数锚点与记录位置」子款（`docs/DELIVERY_DIRECTIVE.md:674-675`，原句在 `:674`、新子款在 `:675`）且含「窗口自 v1.32 起计」；§12.4 表尾出现「按比例抽查记录」块（`:976-977`）；`_EN` 侧同位。
- **OB-84**：§12.5 内该缺口行（`docs/DELIVERY_DIRECTIVE.md:991`）含「**已落地**」与「不改 §7.5 数值」；`_EN` 侧同位含 Landed 与 the §7.5 figures are unchanged。
- **处置列**：OB-60 / OB-82 / OB-84 三行（`docs/DELIVERY_DIRECTIVE.md:938` / `:960` / `:962`）末列首词为 `**已处置**`（`_EN` 为 `**Resolved**`）。
- **版本**：头部（`docs/DELIVERY_DIRECTIVE.md:3`）`版本：v1.32` 与 §0.3 末行（`:78`）`v1.32` 一致（G1-b 实测通过）。

## 5. 自由度与异常

- 自由度 F-1 … F-8 逐项见 `docs/validation/evidence/DIRECTIVE-OB-SWEEP-freedom-list.md`，全部判授权内。
- **A-1（工具形态）**：新建自由度清单时落盘为无 BOM + CRLF，与仓库约定不符 ⇒ 已执行纯编码归一（BOM + LF，内容零改动）。
- **A-2（口径披露）**：§7.7 的新增子款是**加法式**补全；⑨ 的引用可解析性按**变更行作用域**判，整文件继承引文不计入缺陷。
- **A-3（流程偏离 —— §5.3 时序，主控编排失误）**：⑤ r1 报 3 Medium + 5 Low 后，**未按 §5.3 在 ⑤ 原位重跑**，而是径直进入 ⑥；直至 ⑦ r1 报 `FINDINGS` 并完成整改后，才**补跑 ⑤ r2**。⑤ 的 1+1 额度当时**未被占用**（不构成额度例外），且 r1 含 Medium+ ⇒ **不适用 OB-11 例外**；如实登记、不静默。
- **A-4（补跑的输入包与逐项闭合映射）**：⑤ r2 的输入包为「⑤ r1 整改 + ⑦ r1 整改」之后的**同一索引快照**（`git status --porcelain` 证明索引 = 磁盘）；⑤ r1 的 3 Medium + 5 Low 在 ⑤ r2 中**逐项判 CLOSED**（M-1 / M-2 / M-3 / L-1 / L-2 / L-3 / L-4 / L-5）；⑤ r2 自报的 4 条 Low（N-L1…N-L4）整改后，由 ⑦ r2 在最终快照上**代为确认**。（A-3 / A-4 系用户同轮指示补登的事实记录，不触发 ⑦ 重跑。）
- **A-5（⑥ 覆盖范围的边界，事实登记）**：⑥ 的 `PASS` 绑定的是 **⑤ r1 整改后、⑦ r1 整改前的中间快照**；**最终快照由 ⑦ r2 覆盖**，**⑥ 未在最终快照上重跑**（⑥ 为补充扫描层，其结论不构成审查充分性证明）。用户 2026-10-16 裁定：接受「补跑覆盖」、**不销解偏离**。

## 6. §7.5 未变的机械证明与规则句逐字保留

- **§7.5 区域零改动**：`git diff -U0` 的全部 hunk 落在版本行、§0.3、附录 B.1、§7.7、§12.4 与 §12.5，无任何 hunk 覆盖 §7.5。
- **规则句逐字保留**：§7.7 的按比例抽查原句在 diff 中无删除行；变更后该句与 `HEAD` 版本逐字节一致。

## 7. 审查链结论

| 环节 | 承担者 | 结论 |
|---|---|---|
| ① 架构 | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE`（三条件与三项 OB 的边界裁定均成立） |
| ③ 一致性 | 主控（机械脚本） | 引用与数值逐条实测；`G3(a)` 五形状全等 |
| ④ 门禁 | 主控 | tree / index `TOTAL_FAIL=0`；`SELFTEST: PASS 335/335` |
| ⑤ 预审 | `GLM-5.3 (glm)` | 见 §11 的审查链实录 |
| ⑥ 独立扫描 | `MAI-Code-1.1-Flash (copilot)` | 见 §11 的审查链实录 |
| ⑦ 终审 | `GPT-5.3-Codex (copilot)` Extra High 常规档 | 见 §11 的审查链实录 |

## 8. 已知缺口与候选（如实登记，不读作通过）

- **P4 来源引用的行号过期早于本片**：P4 的「来源引用」段写 `docs/DELIVERY_DIRECTIVE.md:915`（OB-60）等，而实测**改动前**该行已在 936（本片后为 938）⇒ 该过期属历史遗留；本片按 §1.7.2 不扩范围，只如实登记为候选。
- **§7.7 窗口起点的历史口径未定义**：本片把窗口起点定为自 v1.32 起计，对 v1.32 之前的片不做追溯认定；若后续要求追溯历史片，属规则语义变更 ⇒ 须另立切片。
- **无专项机械判据**：三项均为措辞与结构类改动，除 G1–G8 的通用判据外无专项判据。
- **§7.7 首片自指的时序**：规则创立片即窗口首片，其判定与理由系收尾时登记（而非启动前），属规则生效当片的字面不可满足；实质义务已履行（§12.4 表尾块已登记本片判定）。
- **候选（不扩范围）**：本次 ⑤ 时序偏离（C 类编排失误）建议由后续准则治理片评估是否补入 §9.3 故障模式库；**本片不改 §12.4**。

## 9. ⑨ 与 ⑧c 读数

- 编码：本片变更与新建的全部 `.md` 均为 BOM + LF；⑨ 全树逐文件核验 `tracked=3607 text=3518 binary=89 violations=17`，**非冻结主树违规 0**（17 处全在 `docs/validation/evidence/**` 的冻结归档包内，按 §5.1 豁免）。
- 引用可解析性（变更行作用域）：212 处 `§X.Y` 引用、未解析 0；G6-5 另报 `scanned 40, refs 56`、无未解析。
- 门禁：`--scope=tree` 与 `--scope=index` 均 `TOTAL_FAIL=0`；`--selftest` `SELFTEST: PASS 335/335`；`G3(a)` 准则对 `+32/-5` 对 `+32/-5`、台账对 `+3/-2` 对 `+3/-2`；`G4b unknown=[]`；`G6-4` 与 `G6-5` 无命中。
- 同位：准则对的五形状（行数 / 标题行号序列 / 加粗 / 竖线 / 空行）全等。
- 残留：`tmp/obsweep/` 与 `tmp/cn.diff` 过程件在 ⑧c 清除；未跟踪残留 0。

## 10. 提交与 CI

- 提交：`2d6a8e4`（`docs: DIRECTIVE-OB-SWEEP 三项 OB 批量处置落地（v1.32）`；7 文件 `+319/-14`）。
- 推送：`origin/main`，范围 `cb64d6b..2d6a8e4`（**未推 gitee**）。
- CI：run `38810986151`（head `2d6a8e4`）**首跑双腿全绿**（windows 3m45s / ubuntu 2m38s）；两腿 `SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=7`、`TOTAL_FAIL=0`；包级 `ok` 30 / `FAIL` 0（ubuntu 腿另跑 Race：`internal/adapters` 与 `internal/chain` 均 ok）；**未触发 §6.4**。
- 提交信息语言**偏离**登记：本片提交信息 summary 为中文，**偏离**「`<type>: <summary>` 英文」纪律；**处置 = 不改写历史**，纪律自下片起生效；依据 = 用户 2026-10-16 裁定原文。
- 本节由 CI 双绿后的 docs-only 提交补入（§11.2）；该回填提交自身的 CI run 只作观察记录；N-L3 口径 = P4 定义行的「三作用域」属历史原文、本片不改，本片自测口径为 C.0t 的「tree 与 index 两作用域（实测）」。

## 11. 审查链实录（⑩ 前最终）

- ⑤ 预审、⑥ 独立扫描与 ⑦ 终审的结论行由 ⑧b 追加于本节末（⑩ 前）；三者额度上限按 §7.5（⑤ 1+1、⑥ 3、⑦ 1+1）。
- 本片 ⑦ = `GPT-5.3-Codex (copilot)` Extra High 常规档；**未降级、未触发 OB-37**；探针 0。
- 事实记录说明：本节在 ⑦ 之后写入，属事实记录（结论行与实测值），按 OB-11 例外口径不触发 ⑦ 重跑。
- **⑤ 预审**：r1 = `PRE-REVIEW: PASS WITH FIXES`（3 Medium + 5 Low，已整改）；**⑤ r1 后未原位复审**（§5.3 时序偏离，见 §5 的 A-3），r2 系 ⑦ r1 整改后**补跑** = `PRE-REVIEW: PASS WITH FIXES`（**0 Medium+**、4 Low）。
- **交由用户裁定**：是否接受「⑦ r1 后补跑 ⑤ r2」作为 §5.3 时序义务的**覆盖**（而非原位闭环）；本片不自行销解该偏离。
- **⑤ 额度用尽（1+1）** ⇒ 其 r2 的 4 条 Low 的字面整改交 ⑦ 复审代为确认（本报告如实登记）。
- **⑥ 独立扫描**：`INDEPENDENT SCAN: PASS`（Section A–E 齐全、Section B 无发现、结论行与 Section B 一致）。
- **⑦ 终审**：r1 = `FINAL REVIEW: FINDINGS`（1 High + 2 Medium，已整改并重建输入包）⇒ r2 = `RE-REVIEW: PASS`（7 项逐条 CLOSED、新发现 Medium+ 为 0）。
- **⑦ r1 的 High 根因**（审查输入包与磁盘不一致）已按「输入包时效性纪律」消除：`impl.diff` / `numstat.txt` / `proof.txt` 与三份门禁读数均取自修复后的同一索引快照。

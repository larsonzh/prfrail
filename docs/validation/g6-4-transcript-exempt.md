# G6-4 转录豁免通道（P93 DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT）

日期：2026-11-26。状态：**已完成后自检，CC 门禁全绿**。切片标识：`DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT`（台账 P93）。

## 1. 变更概要

- 通道选定：G6-4 增**作用域子规则**（三选一之 b）——G7 目标报告内已闭合、且 `Section A`–`Section E` 五段按序的 `sixstage` 围栏体不参与计数型字面量扫描。
- 判定逻辑本体零改动：四条正则、逐行扫描、白名单 split、`emit` 全部未动；只新增作用域判定与一个 `fence-exempt` 计数。
- 单源复用：围栏解析器与报告目标集取自 `g7.js` 的 `internals`（`sixstageBlocks` / `targets` / `SIXSTAGE_SECTION_RE` / `SIXSTAGE_SECTION_SEQ`），`g6.js` 不含本地围栏正则。
- 政策口径：`g6-whitelist.txt` 注释头声明「作用域豁免不是白名单职责」；白名单保持零条目，T23 三断言逐字未变。
- 版本包：准则升 MINOR **v1.76**（头部行 / §0.3 新行 / 附录 `C.0bc`），CN 与 `_EN` 六形状全等。

## 2. 执行流水线

| 步 | 角色 / 模型 | 状态 | 结果 |
|---|---|---|---|
| 0 | 主控（DeepSeek V4.1 Flash High） | 完成 | 锚点七项全符；台账先落 `⬜ 未立项` 再翻 `🔄 执行中` |
| ① | deep-reasoner + DeepSeek V4 Pro | 完成 | 影响判定 + 三选一裁定（选 b）+ 版本包判定 |
| ② | prfrail-implementer + Flash | 完成 | `IMPL: PASS`（逐条对齐设计，未需修正） |
| ③ | prfrail-tester + Flash | 完成 | `TEST: PASS`（回归逐字不变 + 可证伪 + 三作用域） |
| ④ | 主控 | 完成 | tree 作用域 `TOTAL_FAIL=0` |
| ⑤ | DeepSeek V4 Pro | 1 轮 | `PASS WITH FIXES`（1 Low，`no Medium+`） |
| ⑥ | MAI Code 2 Flash | 3 轮 | 末轮 `INDEPENDENT SCAN: PASS`（0 发现） |
| ⑦ | GPT-5.3-Codex 常规档 | 1 轮 | `FINAL REVIEW: PASS`（0 发现，`no Medium+`，复审额度未用） |
| ⑧a | 主控（⑧a 自查见 §6） | 完成 | 报告对 + 审查原文围栏 |
| ⑨ | 主控 | 完成 | 三作用域 + `gofmt` / `build` / `vet` / `test` 全绿 |
| ⑧b | 主控 | 完成 | 台账翻 `✅ 闭环` + 完成记录 |
| ⑧c | 主控 | 完成 | 临时件清空，`tmp` 顶层复位三件 |

批量辅助：`GLM-5.3-Flash` 实派 2 次（双语条款镜像对照、政策三处一致对照），载体 `deep-reasoner`，均 `MIRROR: OK` / `POLICY: CONSISTENT`。

## 3. 异常与兜底记录

- A-1：⑥ 第 1 轮报 1 Medium（外观围栏可躲 G6-4）⇒ 整改：围栏体须通过 G7-f 的结构面（五段按序），补 T141 / T142。
- A-2：⑥ 第 2 轮报 1 Medium（标题行未 `.trim()`，与 G7-f 口径分叉）⇒ 整改：`sixstageStructurallyValid` 对每行 `.trim()`，补 T143。
- A-3：⑤ 的 1 Low（`fence-exempt` 仅在无 finding 时进 detail）⇒ 以准则措辞精确化为「无 finding 时」承接，不改共享 `emit`。该措辞改动已由 ⑥ 与 ⑦ 对改动后文本做实质复核（见 §5）。

## 4. 环境与档位

- 主控：DeepSeek V4.1 Flash（High）；① ⑤：DeepSeek V4 Pro；② ③：Flash；⑥：MAI Code 2 Flash；⑦：GPT-5.3-Codex 常规档；批量辅助：GLM-5.3-Flash。
- 本机 Windows；Ubuntu VM 与代理均未启用；探针 0；未做 gitee。

## 5. 审查结论摘要

- ⑤ 预审（1 轮）：`PRE-REVIEW: PASS WITH FIXES`，1 Low，`no Medium+`。
- ⑥ 独立扫描（3 轮）：r1 `FINDINGS`（1 Medium）、r2 `FINDINGS`（1 Medium）、r3 `PASS`（0 发现）；末轮原文逐字转录于下方围栏内。
- ⑦ 终审（1 轮，盲态）：`FINAL REVIEW: PASS`，0 发现，`no Medium+`；复审额度未使用。

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
| - | - | - | none | - |

## Section C
- G6-4 与 G7-f 在闭合、段序、缩进、info string 空白、提前闭合五面共用 G7 的目标集、围栏解析器、Section 正则与序列：[g6.js](tools/gates/lib/criteria/g6.js#L415)、[g7.js](tools/gates/lib/criteria/g7.js#L593)、[g7.js](tools/gates/lib/criteria/g7.js#L653)、[g7.js](tools/gates/lib/criteria/g7.js#L821)。
- `node --check` 三文件均无语法错误；`--selftest` 为 `PASS 757/757`，T23、T133–T143 全绿；白名单零预载：[whitelist](tools/gates/g6-whitelist.txt#L1)、[selftest.js](tools/gates/selftest.js#L7486)。
- 窄门禁：G6-4 tree `TOTAL_FAIL=0`；G7-f tree `INFO` 且本次无变更报告。文档落点与 P93 台账已核对；本轮未编辑、提交或推送。

## Section D
- 可证伪性：去除闭合、段序或 `trim()` 守卫，或弱化五面正则，T135、T141–T143、T137–T138 会红；白名单新增实际条目会触发 T23/T140 红。
- UNCOVERED：未取得远程 CI 与密码学来源证明；本轮核验限于源码、断言、窄门禁与本机自测。

## Section E
- 密码学来源证明不属于可由静态源码证明的事实；G7-f 形态、逐字转录义务、证据树排除及既有豁免通道已收窄残余。
- 若出现闭合、段序、缩进、空白或提前闭合任一可复现反例，或白名单被预载，PASS 即推翻。

INDEPENDENT SCAN: PASS
```

## 6. ⑧a 自查

- ① 成本表机械审计：有 — §12 的调用次数取自各环回执，批量为 2 次。
- ② 「已处置」`file:line`：无 — 本片未翻转任何 OB 行。
- ③ 旧口径 grep：有 — `fence-exempt` 全仓仅出现在 `g6.js`、`selftest.js` 与准则对，命令与结果见证据包。
- ④ ⑥ / ⑦ 每轮 Info / Low 落入报告：有 — 见 §3 的 A-1 至 A-3。

## 7. A-NN 编号完整性

本报告使用 A-1 至 A-3，均在 §3 有实体定义；无悬空引用。

## 8. 审查输入包摘要

- 输入 = 准则对 diff、`g6.js`、`selftest.js`、`g6-whitelist.txt`、台账 P93 段、证据包；已按 §7.3 消毒（不含任何审查方结论或建议的预置）。
- 盲审隔离：⑦ 首轮不附发现清单；主控未向 ⑦ 透露 ⑤ / ⑥ 的结论内容。

## 9. 可证伪性（机械检查 + 变异检验）

- 变异检验：把 `if (exempt && exempt.has(n))` 改为 `if (false)` ⇒ T133 转红（`SELFTEST` 内建断言 T134 mutation 覆盖）。
- 结构面：缺段（T141）、乱序（T142）、缩进（T143）各有一条反例断言。
- 回归：基线树 `fba4129` 与工作树的 selftest 输出对比，基线行全部仍在，新增恰为 T133–T143，`FAIL` 行使 0。

## 10. 门禁结果

- `--selftest`：`PASS 757/757`（744 基线 + 13 新断言）。
- `--all --skip-base --scope=tree`：`TOTAL_FAIL=0`。
- 基础门禁：`gofmt -l .` 空、`go build ./...`、`go vet ./...`、`go test ./...` 全绿。
- `index` 作用域须在 `git add` 之后取原始读数（纪律 E），见 §16。

## 11. ⑧c 收尾清洁记录

- `tmp/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/` 内为过程脚本与捕获件，⑧c 按件清除。
- `tmp` 顶层复位三件：`.gitkeep` / `gates-cache`（其下仅 `g4b`）/ `find-quotes.js`。

## 12. 成本与计量

- ① 1 次、② 1 次、③ 1 次、⑤ 2 次（1+2 额度内用 2）、⑥ 3 次（≤3 用尽）、⑦ 1 次（1+2 额度内用 1）、批量辅助 2 次、探针 0。
- R = 0、D = 0。

## 13. 明确未执行事项

- 未执行 gitee 推送；未启用 Ubuntu VM 与代理；未改 `ci.yml` / `internal/**` / CONTRACTS / schema / fixtures / ADR；未新增 OB。

## 14. 先例禁用核对

- 引用一条「不得作为先例」条款：P92 的「最小归一 + 围栏外披露」；本片**未使用**该形态，仅以其作为问题出处引用，属允许范围。

## 15. 下一步建议

- 建议下一片处理 `GATES-EXEMPT` 与白名单登记门槛评估（把 ① 判为「大改建议」的 a / c 两项合并评估）。
- 建议 v77 快照吸收本片（P93 出队，下一项回到 S2 前重评）。

## 16. 提交与 CI

- 本节由 CI 双绿后的 docs-only 提交补入（§11.2）；本片停于 ⑩，处于本轮授权范围内。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/selftest.js	repo
tools/gates/g6-whitelist.txt	repo
docs/validation/g6-4-transcript-exempt.md	repo
docs/validation/g6-4-transcript-exempt_EN.md	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/MANIFEST.md	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/SHA256SUMS.txt	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/00-briefing.md	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/01-design.md	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/04-review.md	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/freedom-list.md	repo
docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/input-list.md	repo
```

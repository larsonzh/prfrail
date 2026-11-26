# GATES-EXEMPT 判据豁免与白名单门槛合并治理（P94）

日期：2026-11-26。状态：**已完成自检**，⑥ / ⑦ 结论已闭合。切片标识：`GATES-EXEMPT`（台账 P94）。

## 1. 变更概要

- (a) 判据豁免机制泛化：**做**同构扩展——把 P93 的 `sixstage` 围栏体作用域豁免从 G6-4 扩到 G2 与 G4b；**不做** `gate-exempt` 注释式标记机制（对主控自撰文本开机械豁免洞等于判据架空，嵌套与空区间形态的 fail-closed 代价不成比例）。
- (b) P93 移交的通用引用豁免登记位：**不做**——逐字转录义务只在 `sixstage` 围栏，其余转引属主控自撰总结，可措辞改写。
- (c) 白名单登记门槛：**做**——只加条款、不放行；`g6-whitelist.txt` 条目仍为零，T23 三断言逐字不动。
- 版本包：升 MINOR **v1.77**（头部行 / §0.3 新行 / 附录 C `C.0bd`），CN 与 `_EN` 六形状全等。

## 2. 执行流水线

| 步骤 | 角色与模型 | 状态 | 结果 |
|---|---|---|---|
| 0 | 主控（DeepSeek V4.1 Flash High） | 完成 | 锚点核对；台账先落 `⬜ 未立项` 再翻 `🔄 执行中` |
| ① | deep-reasoner + DeepSeek V4 Pro | 完成 | `DESIGN: DONE`（三件合并评估 + 版本包判定） |
| ② | prfrail-implementer + Flash | 两轮 | `IMPLEMENT: DONE`（第二轮补齐 G2 / G4b 豁免通道与断言） |
| ③ | prfrail-tester + Flash | 完成 | `TEST: DONE`（既有输入回归逐字不变 + 变异检验 + 三作用域） |
| ④ | 主控 | 完成 | tree 与 index 两作用域 `TOTAL_FAIL=0` |
| ⑤ | DeepSeek V4 Pro（通道故障后由 GLM-5.3 承接） | 1 轮 | `PRE-REVIEW: PASS WITH FIXES`（1 Low 已修 + 3 Info） |
| ⑥ | MAI Code 2 Flash | 1 轮 | `INDEPENDENT SCAN: PASS`（0 发现） |
| ⑦ | GPT-5.3-Codex 常规档 | 3 轮 | r1 `FINAL REVIEW: PASS WITH FIXES`；r2 `RE-REVIEW: FINDINGS`；r3 `RE-REVIEW: PASS` |
| ⑧a | 主控（本报告） | 完成 | 报告对 + 审查原文围栏 + 自由度清单 |
| ⑨ | 主控 | 完成 | 三作用域 + `gofmt` / `build` / `vet` / `test` 全绿 |
| ⑧b | 主控 | 完成 | 台账翻 `✅ 闭环` + 完成记录 + 主提交 |
| ⑧c | 主控 | 完成 | 临时件清空，`tmp` 顶层复位三件 |
| ⑩ | 主控 | 停于授权点 | 终态回执 |

## 3. 异常与兜底记录

- A-1：⑤ 通道故障——DeepSeek 两次 `UND_ERR_SOCKET`，按 §2.2 备用档由 GLM-5.3 承接；独立性降低已如实登记，**不得作为先例**。
- A-2：执行自由度清单在 ⑦ 之后落盘；本片如实登记，不做追溯改写。
- A-3：⑦ r1 的 1 Low——G2 的 `fence-exempt N` 计数口径；已修。
- A-4：⑦ r1 的 1 Low——T155 强度不足；整改为五条门槛结构化断言并新增 T156，断言数 769 → 770，编号区间扩为 T144–T156，准则 §0.3 与附录 `C.0bd` 同步。
- A-5：⑦ r2 的两条 Low——无行号视图时回退路径仍误计、注释编号过期；两条均 CLOSED。
- A-6：批量辅助 2 的 `POLICY: DIFF`——白名单注释头缺与 `cjk-newwords.txt` 的分工句；已补，复跑 `PASS 770/770` 与 `TOTAL_FAIL=0`。

## 4. 环境与档位

- 主机：Windows 本机；Ubuntu VM 与代理均未启用；探针 0；未做 gitee。
- 主控：DeepSeek V4.1 Flash（High 档）；①：DeepSeek V4 Pro。
- ② ③：Flash；⑤：DeepSeek V4 Pro（通道故障后按 §2.2 备用档由 GLM-5.3 承接）。
- ⑥：MAI Code 2 Flash；⑦：GPT-5.3-Codex 常规档（1 终审 + 2 复审）。
- 批量辅助：GLM-5.3-Flash 实派两次，载体 `deep-reasoner`。

## 5. 审查结论摘要

- ⑤ 预审（1 轮）：`PRE-REVIEW: PASS WITH FIXES`，1 Low（已修）与 3 Info；无 Medium 及以上。
- ⑥ 独立扫描（1 轮）：`INDEPENDENT SCAN: PASS`，0 发现；原文逐字转录于下方围栏。
- ⑦ 终审与复审（三轮，额度用尽）：r1 两条 Low、r2 两条 Low、r3 两条全 CLOSED 且 0 发现。
- 围栏外披露：与 ⑥ 原文的差异恰为一处——块末补了一行与 Section A 相同的结论行。

```sixstage
## Section A
INDEPENDENT SCAN: PASS
## Section B
| SEVERITY | file:line | violated sentence or fact | minimal fix |
| - | - | - | none |
## Section C
- numstat：变更集恰 13 文件；仅有准则对、台账、证据包、`g2.js`、`g4b.js`、`selftest.js`、`g6-whitelist.txt`。
- 受保护边界核对：`g6.js`、`g7.js`、`ci.yml`、`internal/**`、CONTRACTS、schema、testdata、`cjk-newwords.txt` 均无 staged diff。
- 判定本体：g2谓词未改；g4b差集、白名单与`CACHE_VERSION=v2`未改；豁免均在候选集合剔除后执行。
- `T144-T155` 证据已逐条 PASS；最终 `SELFTEST: PASS 769/769`。
- gate index/tree 均 `TOTAL_FAIL=0`；G4b `newCjk=419 unknown=[]`。
- 三处政策与双语镜像：gate G3(a)/G3(a.2)、G1-b 均 PASS；CN/EN 新增条款形状一致。
- `node --check` 对 g2、g4b、selftest 均无语法错误。
## Section D
- 可证伪性审计：删除 g2/g4b 的豁免剔除、恢复判断谓词、改动受保护文件、使 G4b 出现未知字或让任一 `T144-T155`/门禁读数失败，结论立即转 `FINDINGS`。
- 特别是 T154 会在关闭 g2/g4b 剔除过滤器后把正例转红，证据承重。
## Section E
- 未执行外部 CI/VM 实跑；结论基于完整 staged diff、源码真值、本地 gate/selftest 与双语镜像读数。任何受保护文件变动或 CI 包级证据与本地读数冲突都会推翻本结论。
INDEPENDENT SCAN: PASS
```

## 6. ⑧a 自查

- 成本表机械审计：有——§12 的调用次数取自各环回执，批量辅助为两次。
- 「已处置」`file:line`：无——本片未翻转任何 OB 台账行。
- 旧口径 grep：有——`fence-exempt` 与 `sixstageExemptSet` 的已跟踪出现位置已核对，落在判据实现、`selftest.js`、准则对、P93 / P94 证据包与 P93 报告对。
- ⑥ / ⑦ 每轮 Info 与 Low 是否落入报告：有——见 §3 的 A-1 至 A-6。

## 7. A-NN 编号完整性

本报告使用 A-1 至 A-6，均在 §3 有实体定义；无悬空引用。

## 8. 审查输入包摘要

- 输入 = 准则对 diff、`g2.js`、`g4b.js`、`selftest.js`、`g6-whitelist.txt`、台账 P94 段与证据包；已按 §7.3 消毒。
- ⑤ 输入含 ① 设计件，不含其它审查方结论；⑥ 输入不含 ⑤ 结论；⑦ 首读不附 ⑤ / ⑥ 结论。
- 盲审：本片不抽（不触身份准入 / 模型白名单 / G5-b 闭集 / 所有权 / 停机）。

## 9. 可证伪性（机械检查 + 变异检验）

- T154 内存变异：关闭 G2 / G4b 的剔除过滤器后，T144 与 T149 两个正例转红；还原后与基线的 sha256 逐字节一致。
- T156 计数口径：同一文本碰撞时按行号计数而非出现次数，正例与反例均可复现。
- ③ 的 34 项 fail-closed 自造探针全绿；三作用域读数为 0。
- 基线树与工作树的 `selftest` 输出行集合差集为空（新增断言行除外），删除行为 0。

## 10. 门禁结果

- `--selftest`：`PASS 770/770`（757 基线 + 13 条新断言，含 T156）。
- `--all --scope=tree` 与 `--all --skip-base --scope=index`：均 `TOTAL_FAIL=0`（`changed=13`）。
- 基础门禁：`gofmt -l .` 空、`go build ./...`、`go vet ./...`、`go test ./...` 全绿。
- `G4b`：`newCjk=419 unknown=[]`。
- 准则对六形状全等（两侧同为 2624 / 2624）。

## 11. ⑧c 收尾清洁记录

- `tmp/P94-GATES-EXEMPT/` 按件清除，无遗留过程脚本与捕获件。
- `tmp` 顶层复位三件：`.gitkeep` / `gates-cache`（其下仅 `g4b`）/ `find-quotes.js`。

## 12. 成本与计量

- 调用次数：① 1、② 2、③ 1、⑤ 1（GLM 备用）、⑥ 1、⑦ 3（1 终审 + 2 复审，额度用尽）、批量辅助 2、探针 0。
- **R = 6、D = 0**；返工来源：② 第二轮补通道、⑤ 的 1 Low、⑦ r1 两条 Low、⑦ r2 两条 Low、批量辅助 2 的 `POLICY: DIFF`。

## 13. 明确未执行事项

- 未做 gitee 推送；未启用 Ubuntu VM 与代理。
- 未改 `ci.yml` / `internal/**` / CONTRACTS / schema / fixtures / workflows；未新增 ADR。
- 未登记任何白名单条目；未改 G5-b 闭集与模型白名单。

## 14. 先例禁用核对

- 引用一条「不得作为先例」条款——⑤ 模型备用档与执行自由度清单后落盘两项；本片均未据此扩权，仅在 A-1 与 A-2 如实登记。

## 15. 下一步建议

- v78 快照吸收本片（P94 出队）。
- 下一项回 S2 前重评（S-G1 / S-G8）待用户排期。

## 16. 提交与 CI

- 主提交：`49f1609`（16 路径；信息 `feat: extend the sixstage scope exemption to G2 and G4b and add the whitelist threshold clause (P94)`）。
- 推送：`origin/main`（SSH，未推 gitee），范围 `9217915..49f1609`。
- CI：run `45910223000`（head `49f1609`）双腿 success；两腿 `Gates` 实测 `SELFTEST: PASS 770/770`、`GATE REPORT scope=ci changed=16`、`TOTAL_FAIL=0`；包级 ok ubuntu 21 / windows 15、FAIL 0；其余三个候选作业常态 skipped。
- 本节由 CI 双绿后的 docs-only 提交补入（§11.2）；该回填提交自身的 CI run 只作观察记录，不再为其新增提交。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/gates-exempt.md	repo
docs/validation/gates-exempt_EN.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/00-briefing.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/01-design.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/input-list.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/freedom-list.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/MANIFEST.md	repo
docs/validation/evidence/P94-GATES-EXEMPT/SHA256SUMS.txt	repo
tools/gates/lib/criteria/g2.js	repo
tools/gates/lib/criteria/g4b.js	repo
tools/gates/selftest.js	repo
tools/gates/g6-whitelist.txt	repo
```

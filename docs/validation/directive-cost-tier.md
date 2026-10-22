# 切片 `DIRECTIVE-COST-TIER` 验证报告（准则 v1.43）

> 切片：`DIRECTIVE-COST-TIER`；分级 `[SLICE]`、治理片、成本档 M。
> 准则：v1.43（MINOR，规则修订型；本片为独立治理片，非切片内副产物）。
> 基线：`d:\LZProjects\prfrail`；撰写时点 2026-10-22（⑧a 文档角色）。
> 双语：本文件为 CN 权威；`docs/validation/directive-cost-tier_EN.md` 为逐行同位镜像。
> 边界：本片只改 `.md`；不改 `internal/**`、`tools/gates/**`、`schemas/**`、`testdata/**`、`.github/workflows/**`。

## 1 变更概要

本片一片收口 OB-31（R 阈值的分散表述与兜底计数的合并）与 OB-68（§5.2 的三模型分档与 §7.5 / §8.2 的档位失同步），并衔接 OB-84 的成本口径。

| 文件 | numstat | 内容 |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | `+55/−18` | R 阈值三处统一、⑦ 标签六处改写、§8.2 映射表、§7.5 成本段、§12.4 两行、§12.5 一行 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | `+55/−18` | 与上行逐行同位 |
| `docs/DEV_PLAN.md` | `+8/−0` | T027 段补录 P18 / P19 / P20 与补录说明 |
| `docs/DEV_PLAN_EN.md` | `+8/−0` | 与上行逐行同位 |
| `docs/t027/REMAINING_SLICES.md` | `+20/−0` | P20 状态表行（`✅`）、P21 状态行与 `### P21` 块（`🔄`） |
| `docs/t027/REMAINING_SLICES_EN.md` | `+20/−0` | 与上行逐行同位 |
| `docs/validation/directive-cost-tier.md` | 新增 | 本报告（CN 权威） |
| `docs/validation/directive-cost-tier_EN.md` | 新增 | 本报告的 `_EN` 镜像 |
| `docs/validation/evidence/DIRECTIVE-COST-TIER-freedom-list.md` | 新增 | 执行自由度清单 |

最终变更集为 11 条 `.md` 路径（8 条改写 + 3 条新增；另含 `ADR_REGISTER{,_EN}.md` 的失效日顺延登记）。

```artifacts
docs/DELIVERY_DIRECTIVE.md
docs/DELIVERY_DIRECTIVE_EN.md
docs/DEV_PLAN.md
docs/DEV_PLAN_EN.md
docs/ADR_REGISTER.md
docs/ADR_REGISTER_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
docs/validation/directive-cost-tier.md
docs/validation/directive-cost-tier_EN.md
docs/validation/evidence/DIRECTIVE-COST-TIER-freedom-list.md
```

## 2 执行流水线

步骤 0 先做台账清账，再由 ① 设计、② 落盘与三轮修订、③ 两轮检查、⑤ 两轮预审、⑥ 一次独立扫描、⑦ 两轮终审收口；⑧a 落盘本报告与自由度清单。

| # | 阶段 | `agentName` | `model` | 结论 |
|---|---|---|---|---|
| 1 | ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` |
| 2 | ② 落盘 | `prfrail-documenter` | 主控同模型 | `IMPL: DONE` |
| 3 | 批量辅助 A | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | `EXTRACT: DONE` |
| 4 | ② 修订（成本表与指针化） | `prfrail-documenter` | 主控同模型 | `IMPL-REV: DONE` |
| 5 | 批量辅助 B | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | `MISMATCH`（经 ③ 独立裁决为部分采信，结论不成立） |
| 6 | ③ 一致性检查 r1 | `prfrail-documenter` | 主控同模型 | `CHECK: FINDINGS`（High=0 Medium=1 Low=2 Info=4） |
| 7 | ② 修订 r3 | `prfrail-documenter` | 主控同模型 | 修 F1 / F2 / F3 |
| 8 | ③ 复查 r2 | `prfrail-documenter` | 主控同模型 | `RE-CHECK: PASS` |
| 9 | ⑤ 预审 r1 | `independent-reviewer` | `GLM-5.3 (glm)` | `PRE-REVIEW: PASS WITH FIXES`（1 Low + 3 Info） |
| 10 | ② 修订 r4 | `prfrail-documenter` | 主控同模型 | 修 ⑤ Low#1（§7.5 单位区分句） |
| 11 | ⑤ 复审 r2 | `independent-reviewer` | `GLM-5.3 (glm)` | `RE-REVIEW: PASS` |
| 12 | ⑥ 独立扫描 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS`（0 发现；未给 Section E） |
| 13 | ⑦ 终审 r1（盲审） | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `ADJUDICATION: FINDINGS`（High=0 Medium=1 Low=2 Info=0） |
| 14 | ② 修订 r5 | `prfrail-documenter` | 主控同模型 | 修 ⑦ r1 三条（M-1 / L-1 / L-2） |
| 15 | ⑦ 复审 r2（收窄） | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `RE-REVIEW: PASS` |

**③ 顺序证据**：③ r1 的偏移为 `2026-10-22T10:10:23Z`，早于 ④ 的门禁计时跑（`10:30` 前后）与 ⑤ r1（`10:46:46Z`）；③ 位于 ② 之后、⑤ 之前，与 §5.3 的环序一致。

**⑧a 的位置**：本报告与自由度清单由 ⑧a 落盘；⑧b 起的台账回写、⑨ 证据登记与 ⑩ 收尾由后续步骤执行，本报告不预习其结论。

## 3 异常与兜底记录

按 §9.2 的 R 计数定义逐次登记环、根因与是否合并。

- **D-1（唯一偏离）**：② 修订 r5 写入一个白名单外的单字（码点 `U+7686`），触发 `G4b`（37 条白名单内无该字）⇒ 改写为「两者**均**无」，语义不变；白名单未扩。
- **D-2（证据工件损坏）**：⑨ 期交给 ⑤ 的 `02-changes.diff` 经 PowerShell 重定向捕获后编码损坏（CN 段换行被吞、numstat 与实物不符），由 ⑤ 复审 r2 的备注发现 ⇒ 改用 `node execFileSync` 直写后复得 `FFFD=0`。本条登记为 §9.3「PowerShell 管道捕获原生命令 stdout ⇒ 证据文件乱码」的同族再增一例，不新开 OB；累计序号的真值源为 §9.3 故障模式库，本报告不自行编号。 **处置**：该工件**不作为审查依据**——本片审查以 `03-numstat.txt` 与 `04-status.txt` 为准；⑧c 正常清理该文件，不重做。
- **D-3（③ r1 三项发现）**：F1（Medium）——「附录 B.7」实为**附录 B.0**，已改正（§0.3 行与 OB-31 处置列各 2 处，CN / EN 合计 4 处）；F2（C.0ad 交付物漏列）与 F3（§5.2 加限定语）同轮闭合。
- **D-4（机械集成）**：主控在两轮文档编辑后对 OB-31 处置列的 `file:line` 锚点做机械刷新（`:500/:527/:820/:920/:1151` → `:501/:528/:841/:928/:1159`），属机械集成、不改语义。
- **D-5（用户裁定项）**：「治理片永久豁免试点」是否并入池 E——本片不做，留待用户裁定。
- **D-6（未登记项）**：本片未登记「试点顺延」OB；依据为用户 2026-10-22 的裁定，`[SLICE]` 起跑不登记。
- **R 计数**：D-1 至 D-3 各计 1 次（不同环、不同根因，均不合并），D-4 属机械集成不计 R，D-5 与 D-6 为裁定与登记事项、不计 R；本片 R = 3。按 §9.2，R≥2 ⇒ 按 §5.2 高风险行执行（⑤ 复查为必），R 未超过 3 ⇒ 不触发暂停。**这同时是新 R 定义的首次自举验证**：真实切片中 R≥2 已按高风险行执行（⑤ 复查为必），⑦ 额度按 1+3 执行、实际消耗 2、未超限。

## 4 环境与档位

| 项 | 值 |
|---|---|
| 工作区 | `d:\LZProjects\prfrail`（Windows） |
| 规模与类型 | M 档、治理片、`[SLICE]` |
| 主控档位 | High（§8.2 的「中等」行） |
| ⑦ 档位 | `gpt-5.3-codex` = **Extra High**（未降级，未触发 OB-37） |
| ⑦ 映射表其余两行 | `gpt-6.1-sol` = **Extra High**（允许 **Max**）；`gpt-6-luna` = **Max** |
| ⑦ 额度 | 1 终审 + 3 复审（§7.5 的「治理片」行；多行命中取复审最多者），消耗 2 |
| 探针额度 | 0 |

## 5 审查结论摘要

| 环节 | 轮次 | 结论 | 处置 |
|---|---|---|---|
| ③ 一致性 | r1 | `CHECK: FINDINGS`（High=0 Medium=1 Low=2 Info=4） | ② 修订 r3 修 F1 / F2 / F3 |
| ③ 一致性 | r2 | `RE-CHECK: PASS` | 收敛 |
| ⑤ 预审 | r1 | `PRE-REVIEW: PASS WITH FIXES`（1 Low + 3 Info） | ② 修订 r4 修 Low#1 |
| ⑤ 预审 | r2 | `RE-REVIEW: PASS` | 收敛 |
| ⑥ 独立扫描 | r1 | `INDEPENDENT SCAN: PASS`（0 发现） | 无需整改 |
| ⑦ 终审 | r1（盲审） | `ADJUDICATION: FINDINGS`（High=0 Medium=1 Low=2 Info=0） | ② 修订 r5 修三条 |
| ⑦ 终审 | r2（收窄） | `RE-REVIEW: PASS` | 收敛 |

## 6 审查输入包摘要

- 各审查输入包的文件清单与字节数由 ⑧b 回写时补登；本报告只登记包的构成口径。
- ⑦ r1 按盲审派单：输入不含 ③ 与 ⑤ 的结论行；隔离清单的逐文件读数底稿未给出，故本报告不列清单，只登记隔离方式。
- 本片的自由度清单由 ⑧a 落盘（用户派单口径），时点晚于 ⑦ r2；因此 ⑦ 的输入包**不含**该清单，§1.7.3 的「⑦ 复核清单」未被覆盖，登记见 §8。

## 7 判据与数值的自证

| 改动 | 用户给定值 | 磁盘落点（CN 侧实测） | 自证 |
|---|---|---|---|
| (a) R 阈值三处统一 | `R≥2` 高风险、`R>3` 暂停 | §5.2 高风险行 `:501`；§5.3 的 ⑨ 行 `:528`；§9.2 第 3 条 `:842` | 三处表述一致，并与 §9.2 的 R 计数定义段 `:838` 同口径 |
| (b) ⑦ 标签与映射 | 「⑦（§5.2 分档表选定）」 | §5.1 / §7.1 / §7.5 三行 / §8.2 六处；§8.2 映射表 `:806-812` | 六处标签统一；三模型档位与用户给定值逐格一致 |
| (c) §7.5 成本段 | `175/1400`、`200/1000`、`20/75`（Credits / 1M tokens） | §7.5 成本表 `:719-726` | 三行六格与用户给定值逐格一致；单位区分句已由 ⑤ Low#1 补齐 |
| (d) §9.2 第 3 条 | 与 §5.2 / §5.3 共用同一 R 计数 | `:842` | 明写「共用同一 R 计数，不得形成两套并行阶梯」 |

- **锚点复核（⑧b）**：OB-31 处置列的 `file:line` 锚点已机械刷新为 `:501` / `:528` / `:842` / `:929` / `:1160`，与当前磁盘一致；该刷新属机械集成（不改语义），不计入 R。
- **两处指针化**：§12.1 的试点指标行与附录 B.0 的兜底记录行改为指针，锚点 `:928` / `:1159`；两处的数值权威位留在 §9.2。
- **版本一致性**：头部 v1.43 与 §0.3 末行 v1.43 一致（`G1-b` 实测通过）。

## 8 已知限制

- **⑥ 的覆盖缺口**：⑥ 报 `INDEPENDENT SCAN: PASS` 但未给 Section E；该缺口如实登记，不读作全章覆盖。
- **双单位不得换算**：§7.5 的 Credits 口径与同节「成本梯度提示」的 USD 口径不同，两者不得互相换算或合并。
- **§9.1 不并入 R**：§9.1 的环级条文（重试与暂停判据）保持原样，本片只把 §9.2 的 R 计数定义与三处阈值动作统一，未把环级条文并入 R。
- **自由度清单的时序缺口**：§1.7.3 要求清单在 ⑥ 之后、⑦ 之前落盘并由 ⑦ 复核；本片清单由 ⑧a 落盘，晚于 ⑦ r2 ⇒ 该复核义务未被覆盖，本条为如实登记，不构成 ⑦ 已复核的证明。
- **`G4a` 的判定缺口**：`G4a` 只验「有 BOM」不验「BOM 个数」，双 BOM 可穿过；已登记为 OB-112（`观察中`），处置须另立切片。
- **D-2 的证据口径**：编码损坏的 `02-changes.diff` 不作数；本片以 `node execFileSync` 直写的工件为证据，`FFFD=0`。
- **D-5 未裁定**：「治理片永久豁免试点」是否并入池 E 未在本片决断。
- **无专项机械判据**：本片为措辞与数值的规则修订，除 G1–G8 的通用判据外无专项判据；`G6-4` 与 `G6-5` 在本片新增行上均无命中。

## 9 门禁与性能基线

- **三作用域**：`--all --scope=index|tree|ci` 均为 `TOTAL_FAIL=0`。
- **自检**：`--selftest` = `SELFTEST: PASS 588/588`。
- **异常字**：`G4b` 的 `unknown=[]`，白名单 37 条未扩。
- **会过期字面量**：`G6-4` 无命中。
- **版本**：`G1-b` 的 header=v1.43、§0.3 末行=v1.43。

| 观测项 | 读数 |
|---|---|
| `--all --scope=index` 冷启动 | 116.4 s |
| `--all --scope=index` 热缓存 | 78.8 s |
| `gofmt -l .` | 0.43 s |
| `go build ./...` | 4.72 s |
| `go vet ./...` | 3.50 s |
| `go test ./...` | 11.35 s |
| base 四条合计 | 20.0 s |
| 判据侧孤立计时合计 | 75.4 s |
| 其中 `G4b` | 52.4 s（约 69%） |
| 其余 17 项 | 0.56–2.20 s |

**口径说明**：孤立计时含 18 次进程启动开销，与单次 `--all` 内的实测不可直接相加；本组数据只作性能基线，不改任何门禁逻辑。

## 10 成本与计量

按 §7.5 的硬规则 R7.1，逐项取自会话调用回执。

| 角色 | 次数 | 明细 |
|---|---|---|
| ① 架构 | 1 | `deep-reasoner` |
| ② 实现 / 文档 | 6 | 含 3 次修订（r3 / r4 / r5） |
| ③ 一致性 | 2 | r1 与 r2 |
| ⑤ 预审 | 2 | r1 与 r2 |
| ⑥ 独立扫描 | 1 | `MAI-Code-1.1-Flash (copilot)` |
| ⑦ 终审 | 2 | r1 盲审与 r2 收窄，额度 1+3 未用尽 |
| 批量辅助 `GLM-5.3-FlashX (glm)` | 2 | 摘录与核对两轮 |
| 批量辅助 `GLM-5.3-Flash (glm)` | 0 | 本片批量辅助只用 FlashX，未落到 Flash |

- **Flash 为 0 的原因**：用户 2026-10-22 要求本片优先消耗 FlashX；两轮批量辅助均由 FlashX 承担。
- **失效日延长**：`GLM-5.3-FlashX (glm)` 与 `GLM-5.3-Flash (glm)` 的失效日由 2026-10-31 延至 **2026-11-06**，系用户追加授权；本报告只登记该事实，⑧b 回写台账时须同步。

## 11 明确未执行事项

- 未提交、未推送；未触碰 gitee。
- 未改代码、未改 `tools/gates/**` 判据本体、未改 schema 与 fixtures。
- 未改 §2.2 / §2.4 / §6.3 的规则本体；未改 §7.5 的额度数值。
- 未处置 OB-112；未登记「试点顺延」OB；未决断 D-5。
- 未在 Linux VM 复跑：本片无平台相关改动。

## 12 下一步建议

- ⑧b 回写台账：翻 P21 状态为 `✅`，并登记 FlashX / Flash 的失效日延长事实。
- ⑨ 证据登记：以 `node execFileSync` 直写的工件为准，附 `FFFD=0` 读数。
- 下一轮 `GATES-*` 同批处置 OB-112（`G4a` 的连续 BOM 计数子断言）。
- D-5 的「治理片永久豁免试点」由用户裁定后再落文本。

## 13 执行自由度

本片行使的自由度、由裁定固定的边界与唯一偏离，逐条登记在 `docs/validation/evidence/DIRECTIVE-COST-TIER-freedom-list.md`（F1–F7、M1–M3、D1）。

## 14 未同步项清单

| 项 | 现状 | 责任步骤 |
|---|---|---|
| `docs/t027/REMAINING_SLICES{,_EN}.md` 的 P21 状态 | ✅ 已翻 | ⑧b |
| OB-31 的 `file:line` 锚点 | ✅ 已闭合（`:501` / `:528` / `:842` / `:929` / `:1160` 与磁盘一致；见 §7） | ⑧b |
| `DEV_PLAN{,_EN}.md` 的 T027 状态字段 | 仍为 `BLOCKED / NOT IMPLEMENTED` | 按 §11.2 由后续切片回写 |
| FlashX / Flash 的失效日 | 台账未登记 | ⑧b 回写 |
| ⑥ 的 Section E | 未产出 | 覆盖缺口，登记在 §8 |
| 审查输入包的逐文件字节数 | 未列 | ⑧b 补登 |
| D-5 的池 E 归属 | 未裁定 | 用户裁定 |

DOCS: DONE

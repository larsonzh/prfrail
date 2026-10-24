# GATES-G7F-DEDUP + DISPATCH-EVIDENCE 切片验证报告（治理片：多结论行放宽 + 派单表机械读取）

> 切片：`GATES-G7F-DEDUP + DISPATCH-EVIDENCE`（P26，v1.48，M 档，`[SLICE]` 治理片）。日期：2026-10-24。基线 `HEAD` = `a2d0d3b`（P25 已闭环，准则 v1.47）。变更落点见 §2，审查链见 §3，门禁读数见 §4。

## 1 切片概览

| 字段 | 内容 |
|---|---|
| 标识 | `GATES-G7F-DEDUP + DISPATCH-EVIDENCE`（P26） |
| 日期 | 2026-10-24 |
| 规模与分级 | M（治理片，`[SLICE]`） |
| 基线 | `HEAD` = `a2d0d3b`（P25 已闭环，准则 v1.47） |
| 变更集 | 两准则 + 两台账 + `tools/gates/lib/criteria/g7.js` + `tools/gates/selftest.js`，共六个既有文件；本报告对与自由度清单另行新增 |
| 交付物 | 本报告对（CN 权威 + `_EN` 镜像）与 `docs/validation/evidence/GATES-G7F-DEDUP-DISPATCH-freedom-list.md` |
| 边界 | 只放宽结论行数量与位置配对，不放宽取值闭集与 SEVERITY 一致性；不动 `internal/` / `schemas/` / `examples/` / `testdata/` / `.github/` / CONTRACTS / `DEV_PLAN` |

## 2 变更落点

| 文件 | 位置 | 内容要点 |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` 与 `_EN` | 头部版本行 | `v1.47` → `v1.48`（日期保持 2026-10-24） |
| 同上 | `§0.3` | 新增 v1.48 版本行：① `§6.3` 的 `G7-f` 放宽；② `§7.2.2-A` 与 `§11.3` 转录契约同步多结论行口径；③ `§11.3` 新增派单表机械读取条款；④ 附录 C 增 `C.0ai`；⑤ `§12.4` 新增 `OB-118` |
| 同上 | `§6.3` 的 `G7-f` | 三条改动：句 1 由「恰有一处」改为「允许多处 + 取值全等 + 恰两处位置配对 + 三处及以上或取值不等或位置越界判红」；句 2 改为「**块末结论行**之后不得再有内容」；并补「两处取值亦须各 ∈ §4.1 闭集」「位置越界与取值不等并存时以取值不等优先（仅当块末即结论行）、第二处不在块末则按内容在后报错、finding 条数不作判定口径」「两处取值不等时 SEVERITY 以块末取值为准」 |
| 同上 | `§7.2.2` 的 Section A 条目 | 原位扩写并补「转录进 `sixstage` 围栏时允许两处结论行（Section A 标题后第一个非空行与块末）、两处取值必须相等、口径见 `§6.3` 的 `G7-f`」 |
| 同上 | `§11.3` | 转录契约段行尾追加「多结论行口径」句；新增一行「派单表机械读取」条款 |
| 同上 | 附录 `C.0ai` | 八个字段（目标 / 交付物 / 规模 / 范围 / 审查 / 门禁 / 验收证据 / 附注） |
| 同上 | `§12.4` | 新增 `OB-118`（试点继续顺延第 17 次，`已裁决`） |
| `tools/gates/lib/criteria/g7.js` | 结论行判定段 | 重写：多处 + 取值全等 + 上限两处 + 位置配对；闭集校验扩为逐处 |
| `tools/gates/selftest.js` | T93 / T99 / T104–T107 | T93 变异 needle 更换；T99 主断言 / 变异 / 还原改写；新增 T104–T107（合计 11 条断言） |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | P26 状态行与详情块 | 状态行、目标、交付物、范围、门禁与附注 |

## 3 审查链

### 3.1 派单表（机械读取，`transcript.arguments.agentName` 逐字）

| 角色 | `agentName` | 模型 | transcript 行偏移 | 轮次 | 额度 | 结论 |
|---|---|---|---|---|---|---|
| 批量辅助 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 63347 | r1 | — | `BATCH: DONE`（枚举 g7.js 与契约措辞） |
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 63574 | r1 | 1 | `DESIGN: DONE` |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 63694 | r1 | ≥1 | `IMPLEMENT: DONE` |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 64116 | r1 | ≥1 | `TEST: DONE`（1 Medium + 3 Low） |
| ② | `prfrail-implementer` | 同上 | 64575 | r2 | — | `IMPLEMENT: DONE` |
| ③ | `prfrail-tester` | 同上 | 64831 | r2 | — | `TEST: DONE`（Medium+ 归零） |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 65112 | r1 | 1+1 | `PRE-REVIEW: PASS WITH FIXES`（1 Medium + 2 Low） |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 65354 | r3 | — | `IMPLEMENT: DONE` |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 65587 | r2 | （同上槽） | `PRE-REVIEW: PASS`（Medium+ 归零） |
| ② | `prfrail-implementer` | 同上 | 65759 | r4 | — | `IMPLEMENT: DONE` |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 65894 | r1 | ≤3 | `INDEPENDENT SCAN: PASS`（0 发现） |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 65975 | r1（盲审） | 1+3 | `FINAL: PASS`（0 发现） |

- **两处通道中断**：`P26 ②-3` 首派（偏移 `65311`）与 `P26 ⑤ 预审 r2` 首派（偏移 `65572`）均因 DeepSeek 通道 `UND_ERR_SOCKET` 连接中断而**未返回任何结果**，随后各重试一次（`65354` / `65587`）成功。
- 中断轮次不产生审查结论、不计入 ⑤ 有效调用；⑤ 的有效调用仍为 r1 + r2 两次（额度 1+1 用尽）。

### 3.2 ① 四项判定（`DESIGN: DONE`，结论：不拆片）

| 判定项 | 结论要点 |
|---|---|
| (a) 方案选择 | 采**推荐方案**（判据侧放宽）而非备选（改 `§7.2.2-A`）：冲突根因在判据过度收紧，改契约会破坏逐字转录义务并削弱 A–E 的继承定义；并**加位置约束**（上限两处：Section A 标题后第一个非空行 + 块末），以补住原 `C-4` 逃逸形态 |
| (b) 落点 | `§7.2.2` 的 Section A 条目**原位扩写**（不另起句），`§11.3` 转录契约段行尾追加同口径句 |
| (c) 条款位置 | 派单表机械读取条款落到 **`§11.3`**（报告结构章），不落 `§7.3`：后者对象是审查方的输入消毒与结论落地，而本条款是主控自身的派单与回执义务 |
| (d) 断言调整 | T93 变异 needle 更换、T99 改为「两处且首处不在 Section A 首非空行」反例、新增 T104–T107 覆盖四条 FAIL 分支，每条带变异 |
| 拆片建议 | **不拆**：`OB-118` 已按组合名登记；落点正交；拆片会使审查成本翻倍 |

### 3.3 ③ 测试与变异独立核验（两轮）

- **r1 断言计数对账**：`HEAD` 实测 **663**（② 自报的「基线 662」系用错误增量反推）、当前 **674**、净增 **11 条**（T104 实为 2 条而非 3 条）；T99 **未减少**断言（三条改写而非删减）。
- **r1 独立夹具**：**8 种形态全部与期望一致**；**5 项变异**（位置门 / 取值不等门 / 内容在后门 / 上限门 / 闭集门）全部**具承载性**（改后转绿、还原逐字节一致、`git diff` 归零）。
- **r1 发现**：报 **1 Medium**（`C.0ai` 规模误写 1+2）+ 3 Low + 4 Info。
- **r2**：1 Medium + 3 Low **全部 CLOSED**；Medium+ 归零；新增 2 Low + 4 Info（处置见 §5）。

### 3.4 ⑤ 预审（额度 1+1，两轮用尽）

- **r1 `PRE-REVIEW: PASS WITH FIXES`**：**M1（Medium）**——两处结论行同时「取值不等 + 首处位置越界」时 `else if` 只产一条 finding，准则未钉住该行为与 finding 条数口径；L2（恰两处未写明两处取值须各 ∈ §4.1 闭集）、L3（`C.0ai` 交付物未列 `§12.4` 与台账对）。
- **整改**采 M1 的**方案 (b)**（钉住现行为、零代码改动），并修 L2 / L3。
- **r2 `PRE-REVIEW: PASS`**（Medium+ 归零）；新增 1 Low（M1 括注对「第二处不在块末」子形态过度概括）+ 2 Info；随后按该 Low 建议**收窄括注**并顺带统一 `§7.2.2-A` 的用词（见 §5 的 D-2）。

### 3.5 ⑥ 独立扫描

额度上限三次，**实际一次**。r1 输出按 `§7.2.2` 的 A–E 契约合规，结论 `INDEPENDENT SCAN: PASS`，Section B 无发现。以下为**逐字转录**（未作任何改写）：

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
- None.

## Section C
- The relaxed G7-f rule remains bounded: it still requires the §4.1 closed set, still limits the count to at most two conclusion lines, still requires the two legal positions to be the first non-blank line after Section A and the block-final non-blank line, and still rejects any unequal value or off-position pair.
- The key escape paths remain red: a third conclusion line, unequal values, content after the block-final conclusion line, and a pair whose first line is not the first non-blank line after Section A. The selftests T99 and T104-T107 cover these exact failure shapes.
- The dispatch-table clause in §11.3 is mechanically executable and narrow: it requires reading agentName from the transcript record and attaching the mechanical read-back stdout as evidence. This is a pure read-only rule and does not overlap with the G7-f conclusion-line logic.
- The directive text and the gate implementation stay aligned on the same guard set: closed-set values, fixed positions, at-most-two ceiling, Section B severity/FINDINGS consistency, and the no-note-after-final-conclusion rule.

## Section D
- If the value-equality guard were removed, T105 would turn green when PASS and FINDINGS were paired.
- If the position-pair guard were removed, T99 and the pair-positive control in T104 would no longer be constrained to the legal positions.
- If the block-final trailing-content guard were removed, T106 would turn green.
- If the at-most-two ceiling were removed, T107 would turn green.
- Not covered by the mechanical rule itself: whether the chosen model or carrier was the best one for that stage is still a review judgment, not a dispatch-table read issue.

## Section E
- I cannot verify from the read-only audit the live transcript of a future slice run; this judgment is therefore limited to static rule-text parity and the selftest evidence in the shipped gate.
- What would overturn my PASS: a real transcript that accepts an unequal pair, accepts a pair in the wrong position, accepts content after the final conclusion line, accepts a third conclusion line, or fills the dispatch table by memory instead of reading the transcript's agentName field.

INDEPENDENT SCAN: PASS
```

- 该围栏内为**两处结论行**（`## Section A` 下第一行与块末），取值相等——`G7-f` 放宽后该形态为**合规**，逐字转录无需再合并（P25 的合并规避不再必要）。
- 这同时是本片新增判据的**真实上游回归**：新规则接受 ⑥ 的自然输出形态。

### 3.6 ⑦ 终审（盲审，额度 1+3）

r1 为**盲审**（未提供任何前序发现清单），七项判定全 PASS（目标达成 / 逃逸与禁区 / 文本-实现-测试一致 / 条款关系 / 支撑材料完整 / 边界与镜像 / 断言可证伪），发现清单 **0 条**，结论 `FINAL: PASS`。终审 1 次已用，复审 3 次未使用。**未降级**。

## 4 门禁读数

### 4.1 报告落盘前的定稿读数（六个文件）

| 作用域 | 末行 | EXIT |
|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 |
| `--selftest` | `SELFTEST: PASS 674/674` | 0 |

三个作用域均含 base：`gofmt` 空、`go build` / `go vet` 退出 0、`go test` 包级 `ok` 计数 15。判据读数：`G1-b` `header=v1.48 last=v1.48 rows=49`；`G3(a)` 准则对 `+21/-4` 与台账对 `+19/-0` 双向对称、准则对六形状全等；`G4a` 六个文件；`G4b` `unknown=[]`（白名单 37 条）；`G6-4` / `G6-5` / `G6-6` / `G6-7` 全 PASS；`G7-f` 在本轮为 `INFO … no changed report in scope`（报告尚未落盘，**新规则未被门禁实跑**——新规则的承证来自 `selftest` 的 T99 / T104–T107 与 ③ 的独立夹具）。

### 4.2 含本报告后的定稿读数

| 作用域 | 末行 | EXIT |
|---|---|---|
| `--all --scope=index` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=tree` | `TOTAL_FAIL=0` | 0 |
| `--all --scope=ci` | `TOTAL_FAIL=0` | 0 |
| `--selftest` | `SELFTEST: PASS 674/674` | 0 |

- `G7-f` 本片**首次真实判定**：`PASS`（`2 changed ```sixstage fence(s) complete`，两处结论行取值相等、位置为 Section A 标题后第一个非空行与块末）——即 P25 的合并规避不再必要，⑥ 的自然输出可直接逐字转录。
- 其余定稿读数：`G4a` 九个文件（证据树按既有规则排除）；`G4b` `unknown=[]`（白名单三十七条未扩）；`G6-4` 回填前扫描四百五十四处新增行、回填后四百六十处；`G6-5` 二十八处扫描行、四十八处引用全部解析；`G9` 九个路径通过。

### 4.3 `SELFTEST` 计数变化

`HEAD` 的 `663/663` → 本片 `674/674`，**净增 11 条**（T104 为 2 条、T105–T107 各 3 条）；T99 的三条断言为**改写**而非增删。① 的推算值为 675，实测 674，差异原因已在 ③ 的对账中查明（推算端按「各 3 条」计 12 条，实际 T104 为 2 条）。

## 5 低级别发现处置清单

| # | 来源 | SEVERITY | 位置 | 处置 |
|---|---|---|---|---|
| N-1 | ③ r2 | Low | `C.0ai` 交付物与 P26 边界未列 `§12.4` | **已修**（②-3 在交付物行补 `§12.4` 的 `OB-118` 台账行与台账对） |
| N-2 | ③ r2 | Low | 台账 P26 块「目标」「边界」仍写「`§7.3` 或 `§11.3`」 | **⑧b 回写时统一为 `§11.3`** |
| N-3 | ③ r2 | Info | `§6.3` 的 `G7-f` 行内 `UNCOVERED 不是严重度 token` 连写两遍 | **登记不修**（v1.46 遗留，非本片引入） |
| N-4 | ③ r2 | Info | ② 回执「新串各 2 次」表述不准（全文 6 对 6） | **登记**（不影响产物） |
| N-5 | ③ r2 | Info | 台账 P26 状态行仍为步骤 0 快照 | **⑧b 更新** |
| N-6 | ③ r2 | Info | 「两处取值不等」时 SEVERITY 括注非决定性 | **留档**（三方一致，无矛盾） |
| ⑤ r2 Low-1 | ⑤ r2 | Low | M1 括注对「第二处不在块末」子形态过度概括 | **已修**（D-2：收窄为「仅当块末即结论行时」，并补「第二处不在块末 ⇒ 按内容在后报错」） |
| ⑤ r2 Info-2 | ⑤ r2 | Info | `§7.2.2-A` 用词与 `§6.3` 不一致 | **已修**（统一为「Section A 标题后的第一个非空行」） |
| ⑤ r2 Info-3 | ⑤ r2 | Info | `C.0ai` 交付物「台账两行」计数口径不精确 | **登记不修**（路径与对象均已指明） |
| ⑤ r1 I4–I8 | ⑤ r1 | Info | finding 条数口径、`conclIdx===0` 无变异、引文式结论行 fail-closed、`§11.3` 句面冗余、条款未定义 `agentName` 缺失处置 | **登记不修**；其中 finding 条数口径已由 M1 的整改句覆盖 |
| D-1（临时件，自报自删） | ⑤ r2 | — | ⑤ r2 自述「审查初期曾在**仓库根**误建四个临时输出文件并立即删除」 | **登记**：主控已核验 `git status --porcelain` 无未跟踪项（`untracked=0`）、仓库根仅既有文件、`tmp/` 仅本切片目录与 `.gitkeep` ⇒ 零残留成立；按 `§5.1` 的自愈通道只登记、不停片 |
| D-2（措辞收窄） | 主控 | — | ⑤ 额度用尽后仍修改了 `§6.3` 与 `§7.2.2-A` 的措辞 | **如实登记**：该改动是 ⑤ r2 自己给出的 Low 建议的执行，且随后的 ⑥ / ⑦ 均针对**改动后**的文本审查；⑤ 未再复跑（额度 1+1 用尽），由 ⑥ / ⑦ 覆盖 |

## 6 边界与不追溯

- 变更集只有六个既有文件加本报告对与自由度清单；`internal/` / `schemas/` / `examples/` / `testdata/` / `.github/` / CONTRACTS / `docs/DEV_PLAN` **零改动**。
- 受保护小节零改动：`§2.2`、`§2.4`、`§5.2`、`§7.5` 与其它判据文件；`gate.js` 未改。
- 历史行不追溯：`§0.3` 既有各版本行、`§12.4` 既有 OB 行、附录 C 既有各块均未回改。

## 7 遗留与后续

- `P25` 的 `D-3` 判据冲突**已闭合**；本片 `§6.3` 与 `§7.2.2-A` 的措辞与 `g7.js` 实现经 ③ 独立夹具、⑥ 与 ⑦ 三方核对一致。
- `§11.3` 的派单表机械读取条款自**下一片**起生效，本片 ⑩ 回执已按该条款执行（见回执）。
- 队列：`G4b-CACHE`（G4b 按 HEAD 哈希缓存基线字符集）、`P5` / `P6` 评估动作。
- **下一片观察项（D-2 的流程约束；用户 2026-10-25 裁定）**：审查环结束后，若主控按审查方的建议改动措辞或文本，**必须有下游环节对「改动后」的文本做实质审查**；若下游环节均已用尽且无新审查，则视为**流程偏离**，须由 ⑦ 独立复核，或登记为用户同轮授权例外。本片即该形态：⑤ 额度用尽后按 ⑤ r2 自身给出的 Low 建议修改，由 ⑥ 与 ⑦ 对改动后文本覆盖。

## 8 提交与 CI

| 主提交 hash | CI run 号 | CI 结论 | 说明 |
|---|---|---|---|
| 授权后回填 | 授权后回填 | 授权后回填 | 提交与推送需用户同轮显式授权、绝不推送 gitee；hash 与 run 号按 `§11.2` 只进报告 |

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/GATES-G7F-DEDUP-DISPATCH-freedom-list.md	repo
docs/validation/gates-g7f-dedup-dispatch-evidence.md	repo
docs/validation/gates-g7f-dedup-dispatch-evidence_EN.md	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/selftest.js	repo
```

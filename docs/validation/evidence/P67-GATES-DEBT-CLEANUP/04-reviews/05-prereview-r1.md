# 05-prereview-r1 — P67 GATES-DEBT-CLEANUP（⑤ 预审第 1 轮，逐字归档）

> 归档纪律（§13.1.13 / §11.3）：⑤⑥⑦ 产出逐字原文归档；主控不改写。以下为 ⑤ r1 回执原文。

# P67 GATES-DEBT-CLEANUP ⑤ 预审报告（只读核验）

数据来源说明：全部为读取主控落盘的文件（`docs/validation/evidence/P67-GATES-DEBT-CLEANUP/**`、`tmp/P67-GATES-DEBT-CLEANUP/review-input/prereview.diff`、`tmp/P67-GATES-DEBT-CLEANUP/g4-*-r2.txt`、`docs/DELIVERY_DIRECTIVE.md`、`tools/gates/**`）；未执行任何命令。

## 1. 核对清单（逐项）

**[1] ① 五件完整性 — PASS**
- 四案对比表：`01-arch/design-r1.md:7-14`（§1 六维度 A/B/C/D）
- #1 六条边界确认：`design-r1.md:28-33`（六条逐条 ✅，:34 结论「不超出 ⇒ 不停升级」）
- #2 七坑判定：`design-r1.md:36-48`（27–30 表 :41-44；坑 6 :46；坑 7 :47；`忌` α+β :48）
- #3 清账方案：`design-r1.md:50-58`（A′ 判定 :57；弃 C 标注 :58）
- #4 边界判定：`design-r1.md:60-69`（五加固点 :63-67；B.7.4 不改 :68）
- `DESIGN: DONE`：`design-r1.md:89`

**[2] ② A–G 每步落地 — PASS**
- A：`tools/gates/lib/criteria/g4b.js` 实测（读磁盘）——`CACHE_CAP = 100`（:52）、`CACHE_ENTRY_RE`（:56）、`evictToCap`（:189 起）、`get` 命中后 touch（:236）、`put` 成功后 `evictToCap(dir, CACHE_CAP)`（:257）、`internals` 导出含 `CACHE_CAP`/`evictToCap`。与 `02-impl/impl-r1.md`、`design-r1.md:74` 逐条一致。
- B：+14 断言 = T125×2 + T126×3 + T127×2 + T128×3 + T129×4；基线 722 → 736。可证伪性：8 处删守卫翻红实验（源文件字节零改动）。编号偏差（T119–T124 被 P60 占用 ⇒ T125 起）如实登记。
- C：§13.4 27–30 条于 `docs/DELIVERY_DIRECTIVE.md:1339-1342`；diff 单 hunk 纯插入（+4/-0 双侧同位）；1–26 条含第 21 条逐字未动。
- D：§12.4 清账 A′ 于 `docs/DELIVERY_DIRECTIVE.md:1183-1187`。
- E：附录 B.6 五加固点于 `DELIVERY_DIRECTIVE.md:1531/1539/1540/1541/1543`。
- F：头部 `v1.67`（:3）、§0.3 v1.67 行（:114）、附录 C.0au（:2335-2348）、台账 P67 段。G1-b：④ 复跑 tree/index 均 `header=v1.67 last=v1.67`。
- G：`03-test/consistency-r1.md:7-15`（9 项全 PASS）+ `:29`（`CONSISTENCY: PASS`）。

**[3] ③ 一致性检查充分性 — PASS**
覆盖 3.3 全量 7 项，另覆盖双侧镜像与报告对。另见 F6（读数陈旧）。

**[4] §13.4 与 1–26 无重复定义 — PASS**
27/29/30 与 1–26 无同指；28 与第 24 条同族但不同定义（24 = 版本号误判形状，28 = 聚合措辞真触发），不构成重复定义。β 口径已写入第 27 条处置列。

**[5] §12.4 清账与 §7.7 无冲突 — PASS**
§7.7（:793-800）diff 为零；§12.4 保留「窗口自 v1.32 起计」，与 §7.7:799 逐字一致。

**[6] 附录 B.6 与 §7.2.2 无冲突 — PASS（附 F5 措辞注记）**
§7.2.2 diff 为零；B.6 新增五行均为提示词层指令，未新增判据、未改失败处置语义。

**[7] 边界深查 — PASS**
- #1 LRU 失效场景：touch/unlink/stat/readdir 失败均静默；临时文件被正则排除；唯一理论边角 = 粗粒度 mtime 下并列新写条目可能被稳定排序误淘汰——后果仅一次 self-healing miss，绝不产生错误命中。
- #2 合并不过度：正式化 4 条 ≤5；坑 1+2 同现象合并合理；坑 6 退出（用户裁定）；坑 7 见 F2。
- #4 未引入新判据。

## 2. 发现清单（6 条，均 Low）

| # | SEVERITY | file:line | 被违反的句子（逐字） | 最小修复 |
|---|---|---|---|---|
| F1 | Low | `01-arch/design-r1.md:48` | 「**`忌` 脆弱登记**：选 **α+β**。α：仅登记 `忌` 一字（有实证上下文：`step-A-check.md` R2-2「表述禁忌」）」 | 暂存 diff 无 `tools/gates/cjk-newwords.txt` 改动。修复：⑧ 前补 α（加 `忌` + 理由/上下文摘录）或登记「α 转 ⑩ 待授权」，并更新 MANIFEST/SHA256SUMS。 |
| F2 | Low | `01-arch/design-r1.md:47` | 「**坑 7（§一 滚动窗口）**：归口径改进。……⇒ **取 §12.4 一句**（F-8 允许，与本片 #3 同块）。」 | 落盘的 §12.4 清账块无滚动窗口口径句；报告亦未登记延后。修复：在清账块补一句，或报告 §7/§8 登记「坑 7 转 v55 事实包」。 |
| F3 | Low | `01-arch/design-r1.md:57` | 「② 补录可追溯字段（片号/盲审发现数/是否推翻 ⑥ 的 PASS，从各片报告审查链与 ⑦ 计量行读取）；」 | 实现改为「不转抄」口径；design:58「弃 C…主控须在报告如实登记」在报告 §1 尚无登记。修复：⑧ 装配时在报告登记「② 缩减 + 弃 C」两项偏差及理由。 |
| F4 | Low | `docs/DELIVERY_DIRECTIVE.md:1307` | 「> **来源分组**：P48 六条 + … + F2-FIX 一条（第 26 条）。」 | 27–30 已追加，来源分组注未同步。修复：末段补「P67 四条（第 27–30 条）」；CN/EN 双侧。 |
| F5 | Low | `docs/DELIVERY_DIRECTIVE.md:1531` | 「五段标题固定为 `## Section A` 至 `## Section E`…」 | 与 G7-f 判据存在措辞张力：`###` 形态可过 G7-f 但违 B.6 字面。修复：改为「以 2–4 级标题（常规形态 `##`）呈现、标题文字以 `Section X` 起头」。 |
| F6 | Low | `03-test/consistency-r1.md:15` | 「9. 报告对 53 = 53 行（BOM、LF）。」 | 终态报告为 54 行；③ 读数已陈旧；gates-r1.txt 为步骤 F 前中态。修复：⑧ 收档终态 ⑨ 读数并重算 SHA256SUMS/MANIFEST。 |

无 Medium 及以上发现；上述 6 条均不触及判据本体、契约或准则 §1–§12 语义，无需停下升级。

## 3. 结论

PRE-REVIEW: PASS WITH FIXES

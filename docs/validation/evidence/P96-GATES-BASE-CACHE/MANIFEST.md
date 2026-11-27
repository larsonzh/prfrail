# P96 证据包 MANIFEST

切片：`GATES-BASE-CACHE`（P96）。日期：2026-11-27。

本包含 `MANIFEST.md` 与 `SHA256SUMS.txt` 两件元数据，其余为载荷。

## 归档条目清单

| 载荷 | 内容 |
|---|---|
| `00-briefing.md` | 起跑消息**逐字原文**（§13.1.10；含用户 2026-11-27 引导补正说明） |
| `input-list.md` | 输入清单与步骤 0 锚点实测 |
| `freedom-list.md` | 执行自由度清单（步骤 0 R/D 基线 + 台账登记 + 停点 1 边界 + ② 条款落地记录） |
| `01-design.md` | ① 影响判定与收益论证闸（架构原文归档：A–G 七项 + 条款草案） |
| `04-reviews/03-tester-review.md` | ③ 测试者复核记录（r1 FINDINGS → r2 FINDINGS → r3 PASS；F-1 / F-2 处置） |
| `04-reviews/06-scan-r1.md` | ⑥ 独立扫描 r1 逐字原文（FINDINGS：2 High + 3 Medium） |
| `04-reviews/06-scan-r2.md` | ⑥ 独立扫描 r2 逐字原文（FINDINGS：3 Medium + 1 Low） |
| `04-reviews/06-scan-r3.md` | ⑥ 独立扫描 r3 逐字原文（PASS） |
| `04-reviews/07-deep-review-log.md` | ⑦ 深度档 r1–r4 + Neon + 两轮收窄复审原文（Section A–D；含 Section C–E 逐字段） |
| `05-batch/01-mirror-check.md` | GLM 批量辅助 1 逐字原文（报告对双语镜像机械复核） |
| `05-batch/02-number-check.md` | GLM 批量辅助 2 逐字原文（全片数字核对；载荷数为本二件并入前的时点快照） |

## 编码与行尾

载荷 `.md` 一律 UTF-8 with BOM + LF；哈希在编码归一之后计算。

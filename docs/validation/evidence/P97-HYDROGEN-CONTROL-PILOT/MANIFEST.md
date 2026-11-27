# P97 证据包 MANIFEST

切片：`HYDROGEN-CONTROL-PILOT`（P97）。日期：2026-11-28。

本包含 `MANIFEST.md` 与 `SHA256SUMS.txt` 两件元数据，其余为载荷。载荷在步骤 0 与后续阶段逐次并入；每次并入均同轮更新本清单与哈希（§13.1.12）。

## 归档条目清单

| 载荷 | 内容 |
|---|---|
| `00-briefing.md` | 起跑消息**逐字原文**（§13.1.10） |
| `input-list.md` | 输入清单与步骤 0 锚点实测 |
| `freedom-list.md` | 执行自由度清单（§1.7.3；F-1..F-8 全 in-bounds） |
| `01-arch/01-arch-design.md` | ① 架构方案书（对照对象 / 消毒口径 / 四面定义 / OB-146 文本 / 升版判定；含主控步骤 0 同类计数补充） |
| `02-ruling-stop1.md` | 停点 1 用户裁定（四项 + 同类计数窗口 + 新增裁定 e 的 Hydrogen 派单双停点） |
| `02-comparison-record.md` | ②③ 对照记录（P87 基准侧逐字 + 四面定义；Hydrogen 侧待试水填入） |
| `04-reviews/07-pilot-hydrogen.txt` | ⑦ 试水轮审查方原文逐字（Section A–E；以 `.txt` 归档以避免 G6-4 计数形态；无 BOM + LF） |
| `03-deviation-log.md` | 偏离台账（D-1：输入包路径声明与落点不一致，审查方发现，无结论影响） |

## 编码与行尾

载荷 `.md` 一律 UTF-8 with BOM + LF；哈希在编码归一之后计算。

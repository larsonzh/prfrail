# P94 证据包 MANIFEST

切片：`GATES-EXEMPT`（P94）。日期：2026-11-26。

本包含 `MANIFEST.md` 与 `SHA256SUMS.txt` 两件元数据，其余为载荷。

## 归档条目清单

| 载荷 | 内容 |
|---|---|
| `00-briefing.md` | 起跑消息原文（§13.1.10 留痕） |
| `01-design.md` | ① 影响判定与三件合并评估（架构） |
| `freedom-list.md` | 执行自由度清单（⑥ 后、⑦ 前落盘） |
| `input-list.md` | 输入清单与步骤 0 锚点实测 |

## 编码与行尾

载荷 `.md` 一律 UTF-8 with BOM + LF；哈希在编码归一之后计算。

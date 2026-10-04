# ⑦ 终审记录 · 切片 `DIRECTIVE-MODELID-MAP`（Codex 常规档，未降级）

> 归档副本（v1.22 收尾）：源文件 `tmp/dml/review-r1.md`，正文逐字未改，仅新增本行归档注记，编码保持 UTF-8 with BOM + LF。

> 本文件是 ⑦ 两轮的实际结论落盘（依仓库教训：⑦ 报告须落盘到 `tmp/<切片>/review-r1.md`，否则复审轮无法核验首轮闭合、条数也无法复算）。**盲态首轮**：未附 ⑤ / ⑥ 报告；复审轮按附录 B.6 附「独立扫描完成后读取」的摘要。

## 第 1 轮（`FINAL REVIEW: PASS WITH FIXES`）

| 严重度 | 位置 | 问题 | 处置 |
|---|---|---|---|
| Medium | `docs/validation/evidence/DIRECTIVE-MODELID-MAP-freedom-list.md` | 该文件 `bom=false` 且行尾为 CRLF ⇒ tree 作用域 `G4a BOM+LF` 判红、`TOTAL_FAIL=1`，与「全绿」声明不一致 | 已归一为 UTF-8 with BOM + LF；重跑后 tree / index 均 `TOTAL_FAIL=0`（复审确认闭合） |

**第 1 轮的其余判定（均通过）**：ITEMS-1 规则语义正确性通过；ITEMS-2 边界合规通过（§2.4 黑名单本体、§6.3 判据本体、§5.2 矩阵零改动；`tmp/` 零实际移动）；ITEMS-3 完整性通过（三件均落地、版本包三件套齐全）；ITEMS-4 反例审查 5 条反例均「不成立」（含「仍有 modelId 当调用串」「OB-72 仍自指」「§1.5.1 会误删待归档证据」「漂移机制不触发」「案 A′ 可回到自指」）；ITEMS-5 特别裁定项 F-4 / F-7 / F-8 均判**可接受**；FREEDOM-LIST F-1 至 F-9 逐条判**无越界**（F-4 与 F-7 判「越界触发但处理合规」）。

## 第 2 轮（`RE-REVIEW: PASS`）

- finding 表：`none`。
- 第 1 轮 Medium **已闭合**（只读复核：BOM=true、CRLF=0、LF=38；门禁 `TOTAL_FAIL=0`）。
- 门禁复现：tree `changed=3 / TOTAL_FAIL=0`；index `changed=3 / TOTAL_FAIL=0`（该轮读数采集时本报告对尚未入暂存）。
- 残余风险（⑦ 自述）：若需严格证明「本轮仅编码整改且规则文本逐字未变」，还需与第 1 轮固定基线做逐字对照——本轮命令集只能确认现态全绿且无新增 Medium+。
- ⑦ 明示：⑤ / ⑥ 的闭合结论与「候选 OB 留给 ⑩ 停点裁定」的时机**可接受**，当前无 Medium+ 残留。

## 额度

⑦ 两轮（1 终审 + 1 复审）= §7.5 上限 2/2，**未发生额度例外**；本片因此**无 ADR**。

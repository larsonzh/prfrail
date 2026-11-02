# MANIFEST — DR13-FIX-D3 证据包

切片：`DR13-FIX-D3`（实现类修复片 `[SLICE]`；v37 §四 层 1 第 1 项）。建包时点 2026-11-02（`HEAD` = `e742224`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；后续载荷（01 / 03 / 05 / 06 / 07 / 08）按轮追加，并以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含 D3 裁定与授权边界） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律与升级条件） |
| `01-arch/design.md` | ① 架构《D3 顺序反转与回归方案》原文（V4 Pro 只读载体，逐字归档） |
| `09-glm/glm-1-case-scan.md` | GLM 批量辅助 1/2 用例检索原文（只读扫描，逐字归档） |
| `09-glm/glm-2-impact-scan.md` | GLM 批量辅助 2/2 回归影响面枚举原文（只读扫描，逐字归档） |
| `03-tests/red-baseline.md` | ③ 红测基线：未修复代码上红/绿逐项结论与红证日志路径 |
| `03-tests/tester-report.md` | ③ 测试角色报告原文（逐字归档，含对 ② 的五条注意事项） |
| `04-impl/implementer-report.md` | ② 实现角色报告原文（逐字归档，含改动后函数体与六项自验） |
| `05-prereview/review.md` | ⑤ 预审结论原文（`PRE-REVIEW: PASS`，2 Low + 3 Info） |
| `06-scan/review.md` | ⑥ 独立扫描结论原文（`INDEPENDENT SCAN: PASS`，0 发现） |
| `07-final/review.md` | ⑦ 终审结论原文（`FINAL REVIEW: PASS`，0 发现 + 1 非阻断观察项） |
| `input-list.md` | 输入包清单（OB-93 口径，六项直接输入与字节数） |
| `03-tests/logs/go-test-red.txt` | 红测首跑原始输出（确定性红 + 压力红 4/35） |
| `03-tests/logs/maincontroller-red-verify.txt` | 主控独立复核红测的输出 |
| `03-tests/logs/mc-focused-green.txt` | 修复后聚焦 6/6 转绿原始输出 |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（index/tree/ci 三作用域与 SELFTEST） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）。

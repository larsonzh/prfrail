# P84 CI 首跑红 · §6.4 取证复跑记录（run 44877224368）

- 授权：用户 2026-11-21 裁定，§6.4 一次 failed jobs 复跑，仅作取证；即使复跑绿也不视为修复。
- 复跑动作：`gh run rerun 44877224368 --failed`（一次），同一 headSha
  `d597f1e`（原提交，不含 R-2 修复）。

## 两次结果
| Attempt | windows 腿 | ubuntu 腿 | 备注 |
|---|---|---|---|
| 1（首跑） | ❌ failure | ✅ success | 原始结果，不得掩盖 |
| 2（取证复跑） | ❌ failure | ✅ success | 失败面与 attempt 1 完全一致（确定性红） |

- attempt 2 失败清单与 attempt 1 逐条相同：adapters postflight 全族（统一
  `sha256:7e4012f9…`）、adapters 两条 bind 测试、gates 两条、evidence
  W2/W5b/W8/W9 及 `TestBoundDirJunctionInsideBoundAdmitted`。
- **逐条相同的锚点依据（2026-11-21 补证）**：经 GitHub API
  `…/actions/runs/44877224368/attempts/1/logs` 取回首跑原日志（zip → 解出
  `attempt1-unzip/0_Go windows-latest.txt`，33,902 B），归一化（去时间戳/耗时）
  后逐测试 FAIL 清单 attempt1 = attempt2 = **26/26 完全相同**（机械 diff 零差异，
  见 `attempt1-test-fail-list-normalized.txt`）。首跑原日志含 4 条短名误拒锚点
  （`…\RUNNER~1\…` 存储目标 vs `…\runneradmin\…` 句柄名，见
  `attempt1-shortname-anchors.txt`）——旧 40 行摘录（`--log-failed` 仅取前 40 行
  adapters 段）未覆盖该段，已按裁定补原日志锚点后重写本条。
- 本地证据件：`attempt1-windows-failed-log-excerpt.txt`（40 行，adapters 段）、
  `attempt2-windows-failed-log-excerpt.txt`（80 行）、`attempt1-unzip/0_Go
  windows-latest.txt`（首跑原日志全文）、`attempt1-full-log.zip`（GitHub 原包）、
  `attempt1-test-fail-list-normalized.txt`、`attempt1-shortname-anchors.txt`。
- 结论：两家族根因在 attempt 2 逐字节复现 ⇒ 确定性缺陷，非 flake；首跑红事实
  与复跑结果按 DR / §11.2 如实登记（挂起，待授权）。

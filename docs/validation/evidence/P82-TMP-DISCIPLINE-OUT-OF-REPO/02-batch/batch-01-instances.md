# 批量辅助 1 — 临时件纪律历史实例枚举（GLM-5.3-Flash (glm)，只读，主控转录）

## (a) `tmp/` 根误放
| 实例 | 切片 | 位置形态 | 处置 | 出处 |
|---|---|---|---|---|
| `tmp/seam.js` 建后即删 | CONTRACTS-QUOTA-ALIGN | tmp/根 | 删 + 如实登记 | contracts-quota-align.md D1 |
| `tmp/verify-*.js` + `tmp/pid-reader-dump.md` | AGENTRUNNER-PID-READER | tmp/根 | 自愈登记继续 | agentrunner-pid-reader.md D1 |
| `tmp\s7a-contrasts.js` | CONTRACTS-S7-DURABILITY | tmp/根 | 自愈 + 同轮独立层确认 | contracts-s7-durability.md |
| D-2 取证件 + D-3 仓库根 `count-tmp.js` | P32 D-MECHANICAL-GUARD | tmp/根 + 仓库根 | 非自愈；上界触发后用户裁定不停片 | directive-d-mechanical-guard.md |
| `_p33_en_diff.txt` | P33 V21-OBS-SWEEP | tmp/根 | 非自愈 → 用户裁定继续 | directive-v21-obs-sweep.md |
| D=2 `gates-cache` 根日志 | P34 TMP-SELF-HEAL | tmp/gates-cache 根 | 登记 | directive-tmp-self-heal.md |
| `tmp/g71.txt`；`tmp/gi.txt` + `tmp/gt.txt` | P70 S1-STREAM-PROGRESS | tmp/根（⑧c 后窗口） | 非自愈升级；用户接受 | s1-stream-progress.md §1c |
| `p71-2a-annotate.js`；worktree `p71d-wt` | P71 S1-SCOPE-ANNOTATE | tmp/根（执行角色） | 自愈登记 | s1-scope-annotate.md |
| `watch-p75.txt` / `-backfill.txt` | P75 R3-REAL-CLI-ASSESSMENT | tmp/根 | 自愈 + 零残留 | r3-real-cli-assessment.md §15.1 |
| index 捕获件（事件 1） | P81 MODEL-UPGRADE | tmp/根 | 自愈档成立（三读数） | directive-model-upgrade.md §3 |

## (b) 仓库外位置（`%TEMP%` / `$env:TEMP` 等）
| 实例 | 切片 | 位置形态 | 处置 | 出处 |
|---|---|---|---|---|
| `$env:TEMP` 四例 | GATES-G1G5 | $env:TEMP | 实证并登记 | directive-gates-edge-coverage-recheck.md |
| `gec*.js` / `gec-*.txt` | GATES-EDGE-COVERAGE | $env:TEMP | 宽免（形态偏离不计） | directive-gates-edge-coverage.md |
| `$env:TEMP\w.txt` | DIRECTIVE-TEMP-AUTH-LIMIT | $env:TEMP | 停下升级；即删 + 逐路径回读 | directive-temp-auth-limit.md |
| D1 worktree / D3 stub-parity | AGENTRUNNER-PROD-PATCH | %TEMP% | D1 用户裁定继续；D3 宽免 | agentrunner-prod-patch.md |
| D2 统计件 | CONSOLE-CLI-MULTITASK-WIRING | %TEMP% | 已删除；D 修正 | console-cli-multitask-wiring.md |
| `c1-wire*` / `c1g-tree.txt` | T027-C1-LINUX-NATIVE | %TEMP% | 授权后删除 | t027-c1-linux-native.md |
| 交叉编译产物 | ACCEPTANCE-SECURITY-COUNTEREXAMPLES | 系统临时目录 | 核验无残留 | acceptance-security-counterexamples.md |
| 事件 2 门禁捕获件 | P81 MODEL-UPGRADE | %TEMP% | 非自愈（无自愈路径）→ 用户裁定 | directive-model-upgrade.md §3 |

## (c) 同命令建删
| 实例 | 切片 | 位置形态 | 处置 | 出处 |
|---|---|---|---|---|
| D-2 / D-3 | P32 | tmp/根 + 仓库根 | 非自愈 | directive-d-mechanical-guard.md |
| D=1 | P33 | tmp/根 | 非自愈 | directive-v21-obs-sweep.md |
| D=1 三支脚本建后自删 | P34 | tmp/<切片ID>（位置合法，特例不适用） | 非自愈 | directive-tmp-self-heal.md |
| D1 seam.js | CONTRACTS-QUOTA-ALIGN | tmp/根 | 删 + 登记 | contracts-quota-align.md |
| 事件 2 | P81 | %TEMP% | 「仓库外 + 建删」无自愈 → 非自愈交用户 → v1.71 子款立项 | directive-model-upgrade.md §3 |

**统计**：(a) 10 例、(b) 8 例、(c) 5 例；其中「仓库外 + 同命令建删」组合 = **1 例**（P81 事件 2，即 P82/v1.71 的立项动因）。

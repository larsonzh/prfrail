# MANIFEST — RACE-LEG-EXPAND 证据包

切片：`RACE-LEG-EXPAND`（基础设施类 `[SLICE]`；v38 §四 层 1 第 2 项）。建包时点 2026-11-02（`HEAD` = `81dd49f`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；后续载荷（01 / 03 / 04 / 05 / 06 / 07 / 08）按轮追加，并以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含 R-03 裁定与授权边界） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律与升级条件） |
| `input-list.md` | 输入包清单（OB-93 口径，直接输入与锚点核对） |
| `01-arch/design.md` | ① 架构《CI 扩包方案》原文（V4 Pro 只读载体，逐字归档） |
| `03-tests/red-baseline.md` | ③ 红基线报告：未修改 ci.yml 上静态断言证明四包未被 Race 覆盖（V4.1 Flash 测试角色原文）+ 主控补记 VM 真实 `-race` 四包全绿节 |
| `03-tests/tester-report.md` | ③ 测试角色报告原文（含覆盖矩阵与给 ② 的五条注意事项） |
| `03-tests/logs/vm-race-four-pkgs-green.txt` | VM 真实 `-race` 运行原始日志（Ubuntu VM 10.0.0.199 回拉，环境头部 + 四包 ok 行，rc=0） |
| `03-tests/logs/03-tdd-red-run.txt` | ③ 红基线首跑原始输出（RACE-COVERAGE: RED，exit 1） |
| `03-tests/logs/03-tdd-green-run.txt` | ② 实施后翻转复跑原始输出（RACE-COVERAGE: GREEN，exit 0） |
| `03-tests/logs/03-tests-norace-green.txt` | 四包常规（非 race）测试原始输出（本地 ok 行） |
| `03-tests/logs/03-race-local-infeasible.txt` | 本机 `-race` 不可行取证（无 gcc，cgo build failed） |
| `03-tests/logs/assert-race-leg-coverage.js` | ③ TDD 断言脚本（红→绿翻转判定机制，可复跑） |
| `04-impl/implementer-report.md` | ② 实现角色报告原文（r1 ci.yml 扩包 + r2 工作流冻结契约摘要同步，逐字归档） |
| `04-impl/logs/04-ciyml.diff` | ② r1 的 ci.yml 单行 diff 取证 |
| `04-impl/logs/04-contract-diff.txt` | ② r2 冻结契约单行 diff 取证 |
| `04-impl/logs/04-release-green.txt` | ② r2 `internal/release` 测试全绿输出 |
| `04-impl/logs/04-gotest-all-green.txt` | ② r2 全仓 `go test ./...` 全绿输出 |
| `04-impl/logs/04-red-before-contract-sync.txt` | ② r2 改前红灯取证（7 条同因红） |
| `04-impl/logs/04-release-workflow-tests-verbose.txt` | ② r2 工作流测试逐条 verbose 输出（7/7 PASS） |
| `05-prereview/review.md` | ⑤ 预审结论原文（`PRE-REVIEW: PASS WITH FIXES`，含主控裁定与环境事实纠正补记） |
| `06-scan/review.md` | ⑥ 独立扫描 r1 结论原文（`INDEPENDENT SCAN: PASS WITH FINDINGS`，输出形状偏离 A–E 契约，如实留痕） |
| `06-scan/review-r2.md` | ⑥ 独立扫描 r2 结论原文（A–E 格式合规重发，`INDEPENDENT SCAN: PASS`，预算 2/3，逐字归档） |
| `07-final/review.md` | ⑦ 终审结论原文 r1（`FINAL REVIEW: PASS WITH FIXES`，r2 复核另录） |
| `07-final/review-r2.md` | ⑦ 终审 r2 复核结论原文（Medium 修复确认，逐字归档） |
| `09-glm/glm-1-ciyml-scan.md` | GLM 批量辅助 1/2：ci.yml 结构与 Race 腿历史事实检索原文（只读扫描，逐字归档） |
| `09-glm/glm-2-race-time-estimate.md` | GLM 批量辅助 2/2：各包 Race 耗时评估与 flake 风险扫描原文（只读扫描，逐字归档） |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（index/tree/ci 三作用域 TOTAL_FAIL=0 与 SELFTEST 697/697，含 base 四腿） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|
| （无） | — |

## 未入库文件清单

`00-briefing.md` 已在最终报告 `artifacts` 围栏声明；其余载荷按轮同步声明。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）。

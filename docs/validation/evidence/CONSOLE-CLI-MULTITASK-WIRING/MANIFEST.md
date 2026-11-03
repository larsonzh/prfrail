# MANIFEST — CONSOLE-CLI-MULTITASK-WIRING 证据包

切片：`CONSOLE-CLI-MULTITASK-WIRING`（实现类 `[SLICE]`；v40 §四 层 2 第 3 项）。建包时点 2026-11-03（`HEAD` = `3aa53d5`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；① 架构、③ 红基线、② 实施、⑤、⑥、⑦、⑨ 门禁读数按轮追加归档。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含授权边界、模型绑定、纪律 V 拆步要求与文档同步要求） |
| `01-arch/design.md` | ① 架构设计原文（V4 Pro 只读载体，逐字归档；含 CLI 参数形态、红测规格、fail-closed 清单） |
| `03-tests/red-baseline.md` | ③ 红测报告原文（V4.1 Flash 测试角色，逐字归档；含 4 红 1 绿与偏差） |
| `03-tests/logs/03-tdd-red-run.txt` | ③ 红基线聚焦测试原始输出（主控机械复跑，`-v` 逐用例） |
| `04-impl/implementer-reports.md` | ② 实施角色报告合集（②a 参数面 + ②b 多记录注入，逐字归档 + 主控复核） |
| `04-impl/logs/04-full-suite-green.txt` | ② 后全量 `go test -count=1 ./...` 原始输出（主控复跑） |
| `05-prereview/review.md` | ⑤ 预审结论原文 r1（`PRE-REVIEW: PASS WITH FIXES`，1 High + 2 Low + 3 Info，逐字归档） |
| `05-prereview/review-r2.md` | ⑤ 预审原位重跑结论原文 r2（`PRE-REVIEW r2: PASS`，无 Critical/High/Medium，逐字归档） |
| `06-scan/review.md` | ⑥ 独立扫描结论原文（`INDEPENDENT SCAN: PASS`，无发现，逐字归档） |
| `07-final/review.md` | ⑦ 终审结论原文（`FINAL REVIEW: PASS`，1 Info，逐字归档） |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（构建/vet/全量测试/自检/三作用域，主控落盘） |
| `09-glm/glm-1-cli-flag-inventory.md` | GLM 批量辅助 1/2：现有 CLI flag 定义面检索原文 |
| `09-glm/glm-2-cli-test-enumeration.md` | GLM 批量辅助 2/2：`internal/console` CLI 端到端测试范式枚举原文 |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律、必须升级事项与拆步约束） |
| `input-list.md` | 输入包清单（OB-93 口径，直接输入字节数与锚点核对实况） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|
| `05-prereview/`、`06-scan/`、`07-final/` | ⑤⑥⑦ 审查原文 |
| `08-gates/`、`08-ci/` | ⑨ 门禁读数与 CI 证据 |

## 未入库文件清单

`00-briefing.md` 已在切片报告 `docs/validation/console-cli-multitask-wiring{,_EN}.md` 的 `artifacts` 围栏声明（步骤 0 落盘）；其余载荷按轮同步声明。

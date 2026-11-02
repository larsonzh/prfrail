# MANIFEST — CONSOLE-MULTITASK-ASSEMBLY 证据包

切片：`CONSOLE-MULTITASK-ASSEMBLY`（实现类 `[SLICE]`；v39 §四 层 2 第 3 项）。建包时点 2026-11-02（`HEAD` = `8bd8224`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；① 架构、③ 红基线、② 实施（②a–d）、⑤ r1+r2、⑥、⑦、⑨ 门禁读数均已归档；仅 `08-ci/` 载荷待提交后按轮追加。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含授权边界、模型绑定与文档同步要求） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律与升级条件） |
| `input-list.md` | 输入包清单（OB-93 口径，直接输入与锚点核对） |
| `01-arch/design.md` | ① 架构设计原文（V4 Pro 只读载体，逐字归档；含协议修订草案、红测规格与安全边界） |
| `03-tests/red-baseline.md` | ③ 红测报告原文（V4.1 Flash 测试角色，逐字归档；含偏差与绿形态预告） |
| `03-tests/logs/03-tdd-red-run.txt` | ③ 红基线聚焦测试原始输出（主控机械复跑，-v 逐用例） |
| `09-glm/glm-1-reference-sweep.md` | GLM 批量辅助 1/2：旧单任务限制引用面清查原文（只读扫描，逐字归档） |
| `09-glm/glm-2-test-coverage.md` | GLM 批量辅助 2/2：多任务测试覆盖矩阵与证据清单原文（只读扫描，逐字归档） |
| `04-impl/implementer-reports.md` | ② 实施角色报告合集（②a–d 各步回执逐字归档） |
| `04-impl/logs/04-02a-gate-index.txt` | ②a index 作用域门禁读数（G3(a) 对称 PASS；G5(f) 中途态预期红） |
| `04-impl/logs/04-full-suite-green.txt` | ②d 后全量 `go test -count=1 ./...` 原始输出（主控复跑） |
| `04-impl/logs/04-refusal-strings.txt` | 旧拒绝文案清零取证（单任务门已移除、step 形状门按设计保留） |
| `05-prereview/review.md` | ⑤ 预审结论原文 r1（`PRE-REVIEW: PASS WITH FIXES`，2 High + 2 Medium，逐字归档） |
| `05-prereview/review-r2.md` | ⑤ 预审原位重跑结论原文 r2（`PRE-REVIEW r2: PASS`，H/M/L 零，逐字归档） |
| `06-scan/review.md` | ⑥ 独立扫描结论原文 r1（`INDEPENDENT SCAN: PASS`，无发现，逐字归档） |
| `07-final/review.md` | ⑦ 终审结论原文 r1（`FINAL REVIEW: PASS`，2 Info，逐字归档） |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（三作用域 TOTAL_FAIL=0 与 SELFTEST 697/697，含 base 四腿） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|
| `08-gates/`、`08-ci/` | ⑨ 门禁读数与 CI 证据 |

## 未入库文件清单

`00-briefing.md` 已在最终报告 `artifacts` 围栏声明（⑧ 阶段落盘）；其余载荷按轮同步声明。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）。

# MANIFEST — ACCEPTANCE-SECURITY-COUNTEREXAMPLES 证据包

切片：`ACCEPTANCE-SECURITY-COUNTEREXAMPLES`（安全验证类 `[SLICE]`；v41 §四 层 2 第 4 项）。建包时点 2026-11-03（`HEAD` = `941f60e`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；① 架构、③ 红基线、② 实施、⑤、⑥、⑦、⑨ 门禁读数按轮追加归档。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含授权边界、模型绑定、纪律 V 拆步要求与文档同步要求） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律、必须升级事项、真实漏洞上报红线与拆步约束） |
| `01-arch/design.md` | ① 架构设计原文（V4 Pro 只读载体，逐字归档；含威胁模型、攻击载荷、断言矩阵与疑似漏洞上报清单） |
| `09-glm/glm-1-bounddir-boundread-inventory.md` | GLM 批量辅助 1/2：GAP-51/52/55 代码引用面清查原文 |
| `03-tests/red-baseline.md` | ③ 测试先行红基线报告原文（prfrail-tester，逐字归档；5 面探针全绿、四包基线绿） |
| `05-attack/t51-1-linux-escape.md` | T51-1 攻击载荷与击穿证据（GAP-77/V-1：Linux 实测读出界外 23 字节，逐字含探针源码与运行输出） |
| `04-impl/implementer-reports.md` | ② 实施角色报告合集（②a–②d 四批逐字归档；含偏离清单） |
| `09-glm/glm-2-gates-chain-inventory.md` | GLM 批量辅助 2/2：GAP-60/64/70 代码引用面清查原文 |
| `05-prereview/review.md` | ⑤ 预审结论原文（PRE-REVIEW: PASS，0 Medium+、2 Low 已处置、3 Info，逐字归档） |
| `06-scan/review.md` | ⑥ 独立扫描结论原文（INDEPENDENT SCAN: PASS，逐字归档，含其 ```text 围栏形态） |
| `07-final/review.md` | ⑦ 终审结论原文（FINAL REVIEW: PASS，0 Medium+、2 Info，逐字归档） |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（SELFTEST 697/697、三作用域 0、go build/vet/test、对向 GOOS 复核，主控落盘） |
| `input-list.md` | 输入包清单（OB-93 口径，直接输入字节数与锚点核对实况） |

## 待追加（按轮）

---

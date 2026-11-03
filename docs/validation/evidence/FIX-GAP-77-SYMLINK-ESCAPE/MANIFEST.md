# MANIFEST — FIX-GAP-77-SYMLINK-ESCAPE 证据包

切片：`FIX-GAP-77-SYMLINK-ESCAPE`（底层安全修复类 `[SLICE]`；v42 §四 层 2 第 5 项）。建包时点 2026-11-03（`HEAD` = `dc785ed`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；① 架构、③ 红基线、② 实施、⑤、⑥、⑦、⑨ 门禁读数按轮追加归档。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含授权边界、模型绑定、纪律 V 拆步要求、红测先行红线与文档同步要求） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律、必须升级事项、接口签名冻结、红测先行红线与拆步约束） |
| `input-list.md` | 输入包清单（OB-93 口径，逐件字节数与锚点核对实况） |
| `01-arch/design.md` | ① 架构设计原文（V4 Pro 只读载体，逐字归档；含威胁模型、A/B 对比、混合定稿、函数级设计与降级策略） |
| `01-arch/red-team-sol-review.md` | ① 双打第二棒：GPT-6.1 Sol 底层安全抗压审查原文（§1.7.4 临时授权，逐字归档） |
| `01-arch/design-r2.md` | ① 架构第二轮修订（V4 Pro 处置红队 F1–F6 后定稿；主路径 openat2 单次受限解析 + 逐级 openat 降级；逐字归档） |
| `05-prereview/review-r1.md` | ⑤ 预审 r1 原文（PRE-REVIEW: FINDINGS，2 Medium + 3 Low，逐字归档） |
| `06-scan/review.md` | ⑥ 独立扫描原文（INDEPENDENT SCAN: PASS，A–E 形状，逐字归档） |
| `07-final/review.md` | ⑦ 终审原文（FINAL REVIEW: PASS，0 发现，Codex，逐字归档） |
| `09-glm/glm-1-openat-inventory.md` | GLM 批量辅助 1/2：openat 用法与调用面枚举原文 |
| `08-gates/gate-readings.txt` | ⑨ 门禁读数（SELFTEST 697/697、三作用域 0、VM 全仓与 race、变异检验命中回归文案；主控落盘） |
| `09-glm/glm-2-linux-facts.md` | GLM 批量辅助 2/2：Linux 兼容性与支持矩阵事实核对原文 |
| `05-prereview/review-r2.md` | ⑤ 预审 r2 原位复核原文（PRE-REVIEW: PASS，M1/M2/L2/L3 闭环，逐字归档） |
| `03-tests/remediation-report.md` | ③c 测试整改报告（gopls 告警消除 + ⑤ M1/M2/L2 处置：红基线环境变量门控、fail-closed 契约测试、并发 Close 测试；逐字归档） |
| `03-tests/green-suite-report.md` | ③b 测试重构报告（T51-1 恒红消除：常驻回归 + `-tags p53_red_baseline` 红基线变体；VM 全绿与 race 证据；逐字归档） |
| `03-tests/red-baseline.md` | ③ 红测先行报告（T51-1 真实红命中专属文案、F1 护栏 PASS、F2 环境限制 skip 登记、F5 绿侧 HEAD 红；逐字归档） |
| `01-arch/rulings.md` | 用户 ① 验收裁定原文（(a)–(e)，2026-11-04；逐字归档；降级环境 fail-closed、NO_XDEV 保留、无深度帽、文档同步入 ⑧、openat2 签认） |

## 待追加（按轮）

- `04-impl/implementer-reports.md`：② 实施报告（纪律 V 小步串行）。

---

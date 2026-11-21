# P84 ⑦ 收窄复审（额度例外）逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `GPT-6.2 Sol (copilot)`，档位 Max。
- 额度：本轮 = 用户 2026-11-21 直接裁定的**额度例外 1 轮**（超出 §7.5 与本次切片额度裁定，系用户直接裁定、**不得作为先例**）；四要素：裁决人 = 用户、日期 = 2026-11-21、理由 = N13 为确定性界内误拒、修复方向明确且小，收窄复审一轮值得、依据 = 用户 2026-11-21 裁定消息。
- 输入收窄 = 仅 N13 / N14 整改。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

标准库源码确认，两种 UNC 共享根拼写都会被 `Clean` 原样保留。当前等值逻辑因此存在确定的函数级误拒；我会把它与尚未验证的生产可达性分开报告。

RE-REVIEW: FINDINGS

- **N13：OPEN（Medium）**。盘符根等值已修复，严格后代前缀、Ordinal 比较和非规范目标守卫均保留；但等值改动引入了 UNC 共享根的确定性误拒。
- **N14：CLOSED（原 Low）**。[生产读取入口](internal/evidence/bounddir_windows.go#L799) 在规范化前调用解析器。[NUL 守卫](internal/evidence/bounddir_windows.go#L855) 覆盖两个名称各自声明范围内的全部 UTF-16 单元，只放行末位终止符。对 `PrintName` 同样拒绝内嵌 NUL 不改变 `SubstituteName` 的判定权威，无需收窄。

**独立发现**
| 严重度 | 位置 | 发现与最小修复 |
|---|---|---|
| Medium（阻断） | [bounddir_windows.go:982](internal/evidence/bounddir_windows.go#L982) | 基准为 `\\host\share\`、目标为 `\\host\share` 时，两者都是 Windows `filepath.Clean` 保留的绝对规范形，指向同一共享根。旧等值分支接受，现等值不匹配，后代分支又因目标短于前缀而拒绝。最小修复：仅对同一 UNC 共享根接受有、无尾分隔符两种等值形态，保持主机名和共享名的 Ordinal 比较，并补两个方向的等值测试。 |
| Low（证据归因） | [p84_bounddir_residual_test.go:921](internal/evidence/p84_bounddir_residual_test.go#L921) | 把双分隔符前缀称为「轮 11 前原形」与 [r3 原文](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r3.md#L41) 矛盾。最小修复：将旧等值回退与额外破坏前缀的复合变异分开标注。 |

**可证伪性审计**
- [W10](internal/evidence/p84_bounddir_residual_test.go#L925) 能区分盘符根等值、严格后代和 Ordinal 负例。静态推导：保留单分隔符前缀时，仅回退等值会使 `root-itself` 变红；再把根前缀改成双分隔符，才会另杀两个根后代正例。三红可归因于复合变异，不能据此称为原样回退。W10 未覆盖上述 UNC 反例。
- [W11](internal/evidence/p84_bounddir_residual_test.go#L1068) 的夹具确实把 NUL 写入声明范围。移除守卫会使两个直接解析子测试由拒绝变成解码成功；两个终止符对照应保持通过。[生产入口](internal/evidence/p84_bounddir_residual_test.go#L1130) 可跳过，第三红取决于主机是否写入成功并读回原声明范围。

**无法验证**
已核对[门禁日志](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postR11.txt#L45)的 `TOTAL_FAIL=0` 和[全量测试日志](tmp/P84-BOUNDDIR-RESIDUAL/test-all-postR11.txt#L8)的 `internal/evidence` 通过。全量 `EXIT=0` 和 `SELFTEST 740/740` 为主控登记；现有日志不能独立确认变异实跑、W11 生产子测试是否跳过，或 UNC 反例在生产入口的可达性。

本轮仅审 N13/N14，系用户 2026-11-21 特批的一轮额度例外，超出 §7.5 及本切片额度裁定，不得作为先例。

未运行测试 / 未修改文件

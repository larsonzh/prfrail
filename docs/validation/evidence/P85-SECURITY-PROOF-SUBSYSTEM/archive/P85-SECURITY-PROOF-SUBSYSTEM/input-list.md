# P85 SECURITY-PROOF-SUBSYSTEM · 输入清单（input-list）

切片：P85 `SECURITY-PROOF-SUBSYSTEM`（`[SLICE]` · 实现类 L 档 · 协议先行）
起跑：2026-11-22（用户点名）；主控 = DeepSeek V4 Pro Max；R/D = 0 基线。

## 输入来源

1. 用户起跑消息（原文 = `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/00-briefing.md`，§13.1.10）
2. P83 报告 §8 顺延清单（P84-a / P84-b / P84-c 逐字，不扩不缩）：
   `docs/validation/security-findings-remediation{,_EN}.md`（§4 的 F3 / F6 / F8 行、§5、§8）
3. v69 快照：`docs/t027/CHECKLIST_v69{,_EN}.md`（§一 锚点 + §四 P85 队列项）
4. 准则（v1.71）：`docs/DELIVERY_DIRECTIVE{,_EN}.md`
   必读节 = §1.4 / §5.2 / §5.3 / §7.5 / §7.7 / §8.2 / §9.2 / §13.1 / §13.2
5. 契约与安全权威：
   `docs/CONTRACTS{,_EN}.md`（相关句以 ① 判定书为准）；
   `docs/SECURITY{,_EN}.md`（:5 注记、:23、:60 被引用句）
6. 现状文件（7 个）：
   - `internal/chain/review.go`
   - `internal/console/runtime_agent.go`
   - `internal/console/runtime_supervised.go`
   - `internal/console/ai.go`
   - `internal/adapters/ai_probe_copilot_executor.go`
   - `internal/adapters/ai_probe_copilot.go`
   - `internal/adapters/ai_preflight.go`
7. P83 审查原文：`docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/`
   （终审与复审原文，F3 / F6 / F8 相关逐字）
8. 台账：`docs/t027/REMAINING_SLICES{,_EN}.md` 的 P85 段

## 步骤 0 锚点（实测 2026-11-22）

- HEAD / origin/main = `e35e50f`；准则 = v1.71；SELFTEST = 740/740
- 三作用域 = 全 TOTAL_FAIL=0（index / tree / ci 实测）
- v67 / v68 / v69 = 216 / 218 / 209（LF 计数）
- 机械读数：ADR max = 47；OB max = 139；DR max = 18；试点顺延 = 第 25 次
- Ubuntu VM：larson@10.0.0.199 可达；go1.22.12 linux/amd64（`~/go1.22/bin/go`）

## 边界（逐字）

只闭合 F3 / F6 / F8 顺延部分，不碰 F1 / F4（已全闭合）、F2 / F5 / F7（已闭合）；
不改 S1 Exit 八门判据；不改 SECURITY.md 威胁模型基线；不升准则版本；不做 gitee；
不预设临时授权（额度用尽按 §7.8.3 停下报用户）。

# MANIFEST — P70-S1-STREAM-PROGRESS 证据包

切片：`S1-STREAM-PROGRESS`（P70，层 2 治理侧，实现类 `[SLICE]`，S 档；在 S1 内补 `run --progress` 轻量流式开关）。建包时点 2026-11-13（`HEAD` = `ee7f260`，准则 v1.67）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧ 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10 起跑消息即工件） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `02-batch/batch-01-state-events-enumeration.md` | 批量辅助①（GLM-5.3-Flash）state-events 写入链路与事件闭集枚举原文 |
| `02-batch/batch-02-agentrunner-events-enumeration.md` | 批量辅助②（GLM-5.3-Flash）AgentRunner 事件写入链路与闭集枚举原文 |
| `01-arch/design-r1.md` | ① 架构方案书 design-r1 原文（deep-reasoner + DeepSeek V4 Pro；`DESIGN: DONE`） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 输出原文归档（deep-reasoner + DeepSeek V4 Pro；`PRE-REVIEW: PASS`，3 Low + 3 Info） |
| `04-reviews/scan-r1.md` | ⑥ 独立扫描 r1 输出原文归档（independent-reviewer + MAI-Code-1.1-Flash；`INDEPENDENT SCAN: PASS`，0 发现） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 输出原文归档（independent-reviewer + GPT-5.3-Codex；`FINAL REVIEW: PASS WITH FIXES`，1 Low） |
| `04-reviews/final-r1b.md` | ⑦ 复审 r1b 输出原文归档（收窄输入；`RE-REVIEW: FINDINGS`，1 Low 未闭合 → 整改） |
| `06-gates/gate-readings.txt` | ④ 门禁原始读数（gofmt / build / vet / test / selftest / tree / ci / GOOS=linux） |
| `06-gates/mutation-evidence.txt` | ③ 变异检验证据（M1–M8 红绿表与失败输出原文） |

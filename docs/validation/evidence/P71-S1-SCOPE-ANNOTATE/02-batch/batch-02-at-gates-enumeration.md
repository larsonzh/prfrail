# 批量辅助 ②（GLM-5.3-Flash）AT 与八门现状枚举

派单：`deep-reasoner` + `GLM-5.3-Flash (glm)`
**转录声明**：本文件为主控转录的**要点式摘要**（非逐字全文）；原文见本会话派单回执。

## A. TEST_STRATEGY 的 AT 现状（`docs/TEST_STRATEGY.md`）

- AT 表表头 `:29`：`| ID | 需求 | 场景与断言 | 所有者/计划 |`
- AT-13 `:43`：需求 `REQ-015/019/024`；断言「seed 与 candidate 分离；外部 oracle 比较；clean-room 成功及失败不损 seed」；所有者 `harness/release；T017`
- AT-24 `:54`：需求 `REQ-007/009/021/023/024/025/026/027`；断言「visible 黑箱模式启动前显示并持久化绑定风险确认；只在隔离 workspace，投递不等于完成；显式归还后全量 manifest/diff、范围/秘密/文件/副作用扫描及独立 gates/review；未确认、断线与可检测的未知进程/费用/外部副作用阻断，原理不可观测者标 unknown，报告标记 reduced/unknown 且禁止自动降级」；所有者列已含「AT-24 已验收 = 条件闭合（P64，2026-11-11；`docs/validation/at24-supervised-black-box{,_EN}.md`）」

## B. s1-exit 的门表现状（关键实况）

- 门表表头 `:50`：`| 门禁 | 状态 | 证据/缺口 |`
- **门表共 9 行（`:52`–`:60`），行标识为中文门禁名；不存在以 `S-G1`…`S-G8` 为行标识的行**
- 状态分布：`:52 BLOCKED`、`:53 PASS`、`:54 PASS`、`:55 PASS`、`:56 PARTIAL`、`:57 PARTIAL`、`:58 PARTIAL`、`:59 PASS`、`:60 BLOCKED`
- `S-G1` / `S-G2` / `S-G7` / `S-G8` 字样仅出现在 `:52` 的「证据/缺口」单元格文本内；`S-G3`–`S-G6` 全文不出现
- `:52` 明写「**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过」

## C. s1-exit 的 AT 汇总段

- 小节标题 `:17`：`## AT-01–AT-24 汇总`；表头 `:19`：`| AT | 状态 | 主要证据 |`
- AT-13 `:33`：`PASS`（`[T017 self-host](s1-selfhost.md)`，含 Node 24 run `34325758548`）
- AT-24 `:44`：`条件闭合（2026-11-11，P64 验收；⑦ 复审 RE-REVIEW: PASS）`，证据列含「真实 CLI 端到端仍依赖 T027-B（D-18）」

# P72 `REQ020-S2-RESTRICTION` — REQ-020 的 S2 承诺收窄登记（变更报告）

切片：`REQ020-S2-RESTRICTION`（P72 · 层 2 治理侧 · 文档类 `[SLICE]` · S 档）。主控 `DeepSeek V4.1 Flash (deepseek)`；用户 2026-11-14 同轮显式授权已内含于起跑消息（§1.7.2）。基线 `70eccee`（准则 v1.68）。

## 1. 执行摘要

- **目标**：把 REQ-020 的「阶段」列由 `S0-S2` 改为「S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力）」，并同步 `TEST_STRATEGY` 的 AT-01 与 `S1_SUPPORT_MATRIX` 的 generic/C/Go 行。
- **结论**：闭环。REQ 条目数 28 不变、编号连续；REQ-020「可验收要求」列零改动；MVA 复合标准不动；S1 Exit 八门判定本身不动；`S1_SUPPORT_MATRIX` 行数 6 不变（未新增 Python/Node 独立行）。
- **改动面**：6 个文档对（REQ / AT / MATRIX / 台账）+ 2 处加法式指针（`remaining-rulings` §5.4、`s1-scope-annotate` §7 ③）。
- **停升级边界**：8 条逐条裁决均不触发；(h) 采纳窄读——REQ-011 的「S2 扩展语言」属「开放语言 ID / generic 能力」，收窄文本明示不冻结。

## 2. 差异清单

| 文档 | 现状 vs Q3 冻结要求 | 性质 |
|---|---|---|
| `PRODUCT_REQUIREMENTS{,_EN}.md` | REQ-020 阶段列 `S0-S2` 无收窄注记 | 修改（仅第 3 列） |
| `TEST_STRATEGY{,_EN}.md` AT-01 | 场景列 `C+Go/三语言`（EN `dual/triple-language`）未标阶段归属 | 修改（仅该片段） |
| `TEST_STRATEGY{,_EN}.md` AT-11 | `C+Go 同 task 整体恢复` = S1 恢复语义 | 不动 |
| `S1_SUPPORT_MATRIX{,_EN}.md` 行 12 | 约束列无三语言承诺阶段注记 | 加法式加注 |
| `t027/REMAINING_SLICES{,_EN}.md` P72 段 | 状态与完成记录 | 修改（当前：步骤 0 已落 🔄；计划：⑧ 完成时落 ✅ + 完成记录） |
| `remaining-rulings{,_EN}.md` §5.4 | 无 Q3 冻结落地指针 | 加法式 1 行 |
| `s1-scope-annotate{,_EN}.md` §7 ③ | 待办项未消解 | 加法式 1 行 |

## 3. 逐文档回写方案

- **REQ-020 阶段列**：CN = `S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力）`；EN = `Not required by S1 (C+Go contract only); S2 commitment narrowed (Python/Node official harness + three-language product tasks + the release-support commitment frozen; existing tool scripts / open language IDs / generic capability are not frozen)`（术语取自台账 `REMAINING_SLICES_EN.md` 目标字段既有译本）。
- **AT-01**：`C+Go/三语言` → `C+Go 契约（三语言实跑冻结至 S2，见 REQ-020）`；EN `dual/triple-language` → `C+Go contract (three-language execution frozen to S2, see REQ-020)`；其余单元格逐字不变。
- **AT-11**：判为不动——`C+Go 同 task 整体恢复` 是 S1 内的恢复语义断言，不是支持承诺的阶段归属。
- **S1_SUPPORT_MATRIX 行 12**：约束列末尾加法式追加 `Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力`；「S1 候选状态」列仍为「配置与 noop 编排」。
- **需评估 5 项**：① `remaining-rulings` §5 = 适用（加法式指针，仿 P71 先例）；② `S1_RELEASE_NOTES` = 不适用（支持段未提 Python/Node 官方支持或三语言实跑）；③ `s1-exit` 的 AT-01 条目 = 不适用（状态汇总条目 `PASS`，非阶段归属）；④ `p57-closeout-assessment` §4.2 = 不适用（八门 S-G1..S-G8 无语言支持门）；⑤ `s1-scope-annotate` §7 ③ = 适用（该待办即本片，加法式指针消解）。

## 4. 反例与验证

| 反例 | 机械检出 | 实测 |
|---|---|---|
| 改 REQ-020 第 4 列 | 拆列比对第 4 字段 | 40 字符原值不变 |
| 增删 REQ 行或改编号 | `^\| REQ-\d{3} ` 计数 | CN=28 / EN=28 |
| 新增 Python/Node 独立行 | `^\| (Python\|Node)` 计数 | 0 / 0；表内行 7（表头+6） |
| 改 harness 行状态列 | 第 2 字段比对 | `配置与 noop 编排` 不变 |
| 扩大冻结（删「不冻」保留句） | 与目标整行逐字节比对 | 保留句在行内 |

## 5. 审查链结论

- ① `DESIGN: DONE`（`design-r1`；含 (h) 窄读裁决与 8 条停升级逐条判定）。
- ⑤ 预审：`PRE-REVIEW: PASS WITH FIXES`（2 Low + 4 Info，已全部处理；见 §7 与 `04-reviews/prereview-r1.md`）。
- ⑥ 跳过：纯文档回写，无代码 / schema / fixtures / CI 变更，§7.2.1 三条件成立；经 ⑦ 确认（`INDEPENDENCE GATE (⑥): SKIP CONFIRMED`）。
- ⑦ 终审：`FINAL REVIEW: PASS WITH FIXES`（1 Low + 1 Info，已整改；见 `04-reviews/final-r1.md`）；盲审 0（P71 已抽 ⇒ 窗口重新起算，P72 为第 1 片，未满 4 片不抽）。

## 6. 门禁读数

```text
gate.js --selftest        : SELFTEST: PASS 737/737
gate.js --all --scope=index: TOTAL_FAIL=0
gate.js --all --scope=tree : TOTAL_FAIL=0
gate.js --all --scope=ci   : TOTAL_FAIL=0
```

## 7. 已知边界

- REQ-020 阶段列不再显式书写 `S0`；「S0 验证双/三语言配置」事实保留于第 4 列（零改动不变量）。
- `S1_SUPPORT_MATRIX` 无独立 C / Go 行，harness 行（行 12）是唯一语言行，加注落在该行约束列。
- AT-11 的 `C+Go` 字样保留（恢复语义），不作阶段归属改写。
- ⑤ 整改项：① 报告 §3 的节号引用由「§5」正为「§7 ③」（`design-r1` §6 同处笔误按证据不改写，仅在此登记）；② §1 / §5 审查结论不预写，⑧ 按实际结果定稿；③ §2 台账行表述已正为「步骤 0 已落 🔄」。
- 收窄文本的「不冻」枚举按切片卡目标原文列三项（既有工具脚本 / 开放语言 ID / generic 能力），比 P68 冻结原文少「黄金配置」一词；属非穷尽枚举，不构成冻结范围扩大，按切片卡为准。
- `S1_SUPPORT_MATRIX` 未追加「切片同步（P72 …）」尾注：行内加注已表达同一口径，且 design 未授权该新增行，保持最小改动。
- 证据包 `docs/validation/evidence/P72-REQ020-S2-RESTRICTION/` 属本片审计面（起跑消息【五】第 2 项交付物、§13.1.10 起跑消息即工件），不属产品面改动；本句固化 ⑦ r1 #1 提出的边界口径。

## 8. 明确未执行

- 未改准则正文 / 判据本体 / `CONTRACTS` 语义 / schema / fixtures / CI / v50–v57 快照。
- 未改 REQ 表结构（条目不增删、编号不变、「可验收要求」列零改动）、未改 MVA 复合标准、未改 S1 Exit 八门判定本身。
- 未新增 `S1_SUPPORT_MATRIX` 的 Python/Node 独立行；未替用户裁定；未 commit / push；未启动新切片。

- 载体偏离（§1.7.4 登记）：本片 ② / ③ / ⑧ 的文档改写与一致性核验由主控执行（未派 prfrail-documenter），原因 = 纯整行替换 + 机械核验，派单成本高于收益；该偏离不改变改动面与判据，已在 ⑩ 回执登记，不得作为先例。
```artifacts
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/validation/s1-scope-annotate.md	repo
docs/validation/s1-scope-annotate_EN.md	repo
docs/validation/req020-s2-restriction.md	repo
docs/validation/req020-s2-restriction_EN.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/00-briefing.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/input-list.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/freedom-list.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/MANIFEST.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/SHA256SUMS.txt	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/01-arch/design-r1.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/02-batch/batch-01-req-enumeration.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/02-batch/batch-02-at-matrix-enumeration.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/06-gates/gates-r1.txt	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/06-gates/selftest-r1.txt	repo
docs/validation/evidence/P72-REQ020-S2-RESTRICTION/04-reviews/final-r1.md	repo
```

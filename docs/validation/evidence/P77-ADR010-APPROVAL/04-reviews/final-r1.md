# ⑦ 终审 r1 原文（`GPT-5.3-Codex (copilot)` / `independent-reviewer`）

> 归档纪律：本文为 ⑦ 终审 r1 的自然输出逐字转录，仅将位置链接改写为纯路径、并把行内的空表格补齐为完整行；结论与判定未改写。

FINAL REVIEW: PASS WITH FIXES

## 发现清单
| # | SEVERITY | file:line | 被违反的句子/判据 | 最小修复 |
|---|---|---|---|---|
| 1 | Low | `docs/validation/adr010-approval.md:66`、`:129`、`:132` | 第 3 问要求「残留项已在报告 §8 披露」。当前冻结快照类残留已披露在 §3.3（第 66 行），但 §8 仅列了 p57、s1-exit、DEV_PLAN 三类，未把冻结快照再显式列入 §8 | 在 §8 增补 1 条「冻结快照类残留」说明，并引用 §3.3 的同源口径，避免读者只看 §8 时漏读 |

## 七问逐项回答

1. ADR-010 状态列格式是否满足「APPROVED + 四要素 + 分号后历史句保留 + 闭括号追加」
结论：是。
证据：`docs/ADR_REGISTER.md:20` 显示状态词为 APPROVED，括注中含裁决人、日期、理由、依据；分号后历史完成条件句保留，末尾追加「已于 2026-11-16 落实；P77」。`docs/ADR_REGISTER_EN.md:20` 英文镜像同构。`docs/ADR_REGISTER.md:102` 与 `docs/ADR_REGISTER_EN.md:103` 的「平台/支持」列同步追加了「已落实」闭括号指针。

2. 边界是否满足：仅文档改动，p57 门表零改动，s1-exit 仅风险第 1 条，且 CHECKLIST_v59/v60、v50-v58 快照、代码与 schema/fixtures/CI/go.mod 未改
结论：成立。
证据：`tmp/P77-ADR010-APPROVAL/impl.numstat` 全为 docs 路径；无 go/json/yml/schema/fixtures/go.mod。`tmp/P77-ADR010-APPROVAL/impl.diff` 显示 s1-exit 仅「剩余风险第 1 条」单行替换。numstat 未出现 p57 与 CHECKLIST_v59/v60 文件项。

3. 六文档与 ADR_REGISTER 两处是否语义一致；残留项是否允许保留且在 §8 披露
结论：语义一致成立；残留项「允许保留」成立；「全部在 §8 披露」不完全成立。
证据（语义一致）：`docs/S1_SUPPORT_MATRIX.md:21` 与 `:23` 已切到「已批准」语义，并保留「本页不替代正式公告」的通用限制。`docs/INSTALLATION_PLAN.md:32` 已写「EOL 通知方式已由 ADR-010 获批」。`docs/PRODUCT_REQUIREMENTS.md:107` 已写「该项不再待落实，其余不变」。`docs/validation/s1-exit.md:66` 已同步「EOL 通知与例外已获批」。`docs/ADR_REGISTER.md:62` 与 `:102` 两处已与批准态对齐。
证据（允许保留）：p57 门表旧状态行（`:91`）仍保留；s1-exit 带日期历史结论段（`:7`）仍保留；这两类在报告 `:129`–`:132` 明确说明。冻结快照类允许保留在报告 `:66` 已披露。
偏差点：§8 本身未重复列「冻结快照类」，见「发现清单」。

4. S-G4 收口是否可由本报告承载而不改 p57 旧行
结论：可以，且当前实现符合该策略。
证据：报告 `:92`–`:96` 明确写明：p57 行不改，收口由本报告承载，并给出「转为通过」结论。p57 旧行仍在，符合「历史门表零改动」的承载方式。

5. CN/EN 对称与报告六形状全等，是否有原始读数支撑
结论：有。
证据：`tmp/P77-ADR010-APPROVAL/gate-index-4e.txt` 给出「symmetry + mirror parity」，包含 CN/EN 成对增减行一致与报告形状项；同文件末尾为 `TOTAL_FAIL=0`。

6. ⑤ r2 三条 Low 是否关闭
结论：三条均可判 CLOSED。
证据链：Low1 来源 `04-reviews/prereview-r2.md` 第 1 条；现态报告 `:72` 已改为「性质尾注不入状态列；S-G4 结论在 §5；四要素与日期在 §1」，与 r2 最小修复方向一致。Low2 来源 `04-reviews/prereview-r2.md` 第 2 条；现态报告 `:73` 已明确「首句与尾句」双项分歧与取舍依据。Low3 来源 `04-reviews/prereview-r2.md` 第 3 条；现态报告 `:88` 已将引用落点改到 §7.2，且 `:121`–`:125` 给出三作用域与 selftest 读数。

7. ⑥ 跳过是否应确认为 SKIP CONFIRMED；并给 §7.8.2(i) 覆盖/遗漏结构化说明
结论：应确认。
证据：报告 `:115` 已明确「⑥ 跳过，三条件成立，由 ⑦ 确认」。`docs/DELIVERY_DIRECTIVE.md:713`–`:719` 要求三条件同时满足且须由 ⑦ / 用户确认。`tmp/P77-ADR010-APPROVAL/impl.numstat` 支撑「仅文档改动」。机械替代证据读数：`gate-index` `TOTAL_FAIL=0`、`gate-tree` `TOTAL_FAIL=0`、`gate-ci` `TOTAL_FAIL=0`、`selftest` `PASS 737/737`。

## ⑤ r2 三条 Low 的裁定
1. Low1：CLOSED。证据见 `04-reviews/prereview-r2.md:20` 与报告 `:72`。
2. Low2：CLOSED。证据见 `04-reviews/prereview-r2.md:21` 与报告 `:73`。
3. Low3：CLOSED。证据见 `04-reviews/prereview-r2.md:22`、报告 `:88`、`:121`。

## ⑥ 跳过确认
SKIP CONFIRMED: YES

按 §7.2.1 三条件核对结果：
1. 条件①（仅文档与文案类）满足。证据：`impl.numstat` 清单。
2. 条件②（无 go/json/yml/schema/fixtures 脚本语义变更）满足。证据同上。
3. 条件③（不涉状态机/门禁判据/证据模型/角色授权契约语义）满足。证据：改动仅在 ADR 状态事实、支持矩阵与验证报告承载面；`DELIVERY_DIRECTIVE` 未在 numstat 清单内。

按 §7.8.2(i) 的结构化说明：
1. 覆盖了哪些方面
- 变更面边界与文件类型覆盖：numstat 全量清单。
- 双语对称与报告形状覆盖：gate-index 的 G3(a)。
- 机械门禁收敛覆盖：index/tree/ci/selftest 全部通过。
- S-G4 收口承载链覆盖：ADR 状态列 + 本报告 §5 + p57 零改动链条。
2. 遗漏哪些方面
- 未覆盖推送后真实 CI 运行态；当前 ci 读数说明本地取上一提交变更集，见报告 `:124`。
- 未覆盖用户原始起跑消息字节级对照，只能依据转录与落盘文本。
- 未覆盖运行时行为验证（本片为文档切片，不涉及代码执行语义）。

## 可证伪性审计
1. 若在 `impl.numstat` 出现任一非 docs 文件，本审查第 2、7 问立即失效。
2. 若 `docs/ADR_REGISTER.md:20` 不再同时包含 APPROVED、四要素、历史句与闭括号，第 1 问立即失效。
3. 若 `docs/S1_SUPPORT_MATRIX.md:23` 回退到「尚未获批」语义，第 3 问立即失效。
4. 若任一门禁读数不再为 0 或 selftest 不再为 PASS 737/737，第 5、7 问立即失效。
5. 若 §8 仍不补「冻结快照类残留」，本次 Low 发现持续有效。

## 无法验证的部分
1. 用户起跑消息原文与报告中「同轮裁定原文」是否字节级一致，当前输入包未含会话原件。
2. 推送后远端真实 ci 作用域读数，本轮只读材料未包含该时点回执。
3. ⑥ 实际执行日志是否发生（仅能看到「已跳过并由 ⑦ 确认」的报告声明），见报告 `:115`。

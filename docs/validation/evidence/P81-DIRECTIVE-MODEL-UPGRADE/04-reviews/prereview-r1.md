# P81 ⑤ 预审 r1 原文（`DIRECTIVE-MODEL-UPGRADE`）

> 载体 / 模型：`deep-reasoner` / `DeepSeek V4 Pro (deepseek)`（只读）。逐字转录自 ⑤ 的最终消息。
> 编码：UTF-8 with BOM + LF。

## 总体结论

本环按 §4.2.1 七条完成最低核验（关键工件实读、闭集/计数独立复算、门禁读数与磁盘态交叉核对、可证伪性推理、授权可回溯性溯源）。10 项必核逐条结论：**1 落点完整**——四组旧标识在 CN/EN 的现役面（§2.2/§5.2/§6.3 G5-b/§7.5/§8.2/附录 B.6）已全部换名；残留仅落在历史块（CN `DELIVERY_DIRECTIVE.md` L65 v1.18 日志、L117 v1.70 新日志行的「由旧改新」描述侧、L1092/1101/1107/1171 的 OB 台账、L1703/1723/1744/1805 的附录 C 既有块），EN 同构（L65/117/1092/1101/1107/1171/1703/1723/1744/1805），与 batch-01 分层（A=6 行受控，准则对 A 层 = 0）一致。**2 判据本体未动**——`g5b.js` diff 仅两行成员替换，`run()`/`FRONTMATTER`/`EMPTY_VALUE_FORMS`/注释逐字未改，磁盘态复核确认。**3 自测同步**——三条新 check 断言绑定精确值（T131a 需 PASS+0 发现、T131b/c 需 FAIL+1 发现+detail 含 `model=<旧id>`，detail 格式由 `g5b.js` 的 `:: model=` 拼接保证），删闭集成员即变红、加回旧标识亦变红，可证伪；FIXTURES 三键名与磁盘夹具文件名一一对应；737=CHECKS.length（改前读数 `SELFTEST: PASS 737/737`）+3=740 算术成立。**4 数值不动**——§7.5 Sol 行仍 200/1000（C-3 pricing-only 双语言 PASS），§8.2 仍「**Extra High**（允许 **Max**）」，diff 内无其他数值改动。**5 CN/EN 对称**——numstat 29/10=29/10；六形状全等（C-5 指标 CN=EN：2453 行/149 标题/6002 `**`/4107 `|`/565 空行/6628 反引号；G3(a) 门禁六布尔全 true 双源确认）；EN 无语义漂移（v1.70 行、P80 块忠实，授权原文保留中文）。**6 P80 留痕**——§12.4 块授权原文与 sg6 报告 §14.3 **逐字一致**；「不得作为先例」在位（块尾「落点差异经用户裁定、不得作为先例」+原文内「不得作先例」）；不占 OB 号。**7 版本元数据**——头部与 §0.3 末行均 v1.70/2026-11-18（G1-b、C-4 双源）；C.0aw 位于 `### C.1` 之前，G6-5 PASS。**8 越界**——17 文件全在授权面；p57/`s1-exit`/v63 快照/CONTRACTS/schema/产品代码/CI/`.github/agents/**` 零改动；无 ADR/OB 新增；Astra 未扩列（闭集仍 8 成员）。**9 三处自报偏离**——均判「合理执行自由度」：A4 锚点（实现与 Q4 判定插入位一致、标注更严）；A6 两成员替换（与 Q1 点 6 一致）；CN 日志行补 `**Max**` 粗体（忠实引用 §8.2 原文形态，EN 同构，不改变语义）。**10 自由度**——F-1/F-3/F-4 in-bounds；F-2 定论 **in-bounds**（起跑 Q3 已授权「夹具新增/更新」，① 判「已授权包内夹具同步」）；F-5 定论 **in-bounds**（G5-c 动态读 §2.2 限定名列，不换名即 FAIL，属判据联动强制）；遗漏项见偏离表。**§4.2.1 第 7 条可回溯性**：替换方案/档位 (a)/定价实读可溯至 v63 快照 §四（`CHECKLIST_v63.md:138-139`）与 `00-briefing.md` 逐字转录；P80 授权溯至 sg6 报告 §14.3 ✓。探针 3/3 全部 PROBE-OK。门禁 tree 读数 TOTAL_FAIL=0、G8-b 包自洽、go 三连绿。仅 4 条 Low（记录可延后），无 Medium+。

## 偏离项表

| 严重度 | file:line | 被违反的句子 | 最小修法 |
|---|---|---|---|
| Low | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/freedom-list.md:10,19` | 「待 ① 判定」 | ① 方案书已判定 F-2=in-bounds（design-r1.md「待升级」节），清单定稿时回写 F-2/F-5 为 in-bounds 及判定依据 |
| Low | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/freedom-list.md:17` | 「正例夹具值随新成员更新；新增退役标识的反例夹具」 | F-2 说明补登新增正例夹具 `g5b-closed-set-legal-mai.txt`（清单未覆盖） |
| Low | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/freedom-list.md`（条目表整体） | 未列「CN v1.70 日志行补 `**Max**` 粗体」 | 增列该自由度，判定 in-bounds（随 §8.2 原文的引文逐字忠实） |
| Low | `docs/validation/evidence/P81-DIRECTIVE-MODEL-UPGRADE/input-list.md:21` | 「`gate-index.txt` / `gate-tree.txt` / `gate-ci.txt` / `gate-selftest.txt`」 | 文件名与实际落盘（`gate-tree-04.txt`、`*-before.txt`）不符；且改后 SELFTEST（预期 740/740）、index、ci 三份原始读数未随包递交——⑥ 派单前补齐并修正清单 |

## 反例分析

本片为常量闭集替换 + 文档 + 自测夹具，无并发/崩溃面（`g5b.js` 无逻辑改动，fail-open 与大小写/逐字比较语义未动）。边界可证伪判据：① 删 `mai-code-2-flash`/`gpt-6.2-sol` 任一成员 ⇒ T131a/T41 变红；② 把退役标识加回闭集 ⇒ T131b/T131c 变红；③ `consistency-check.js` 的 M1–M4 变异（改价/改末行/回退 EN 限定名/破行数对称）分别点亮 C-3/C-4/C-2/C-5 守卫，守卫缺失则变异不被检出——可证伪性成立。未发现的边界风险：若 operator UI 显示名与写入限定名再漂移，R2.2 停下升级机制已覆盖（probe-1/2 本片实测通过）。

## 其他关注点

- §12.4 P80 块的输入收窄声明相对 sg6 §14.3 为压缩转述（「④ 本章」→「④ 收窄任务所需章节」；任务句去掉「（§12 的 `--scope=ci` 退化说明）」括注），授权原文逐字未损，不判偏离，仅登记。
- freedom-list 状态仍为「初稿」，符合其自述时序（⑥ 后定稿）；上述回写属定稿义务。
- ⑥ 本片执行模型 = 旧 `MAI-Code-1.1-Flash`（v1.69 生效期内唯一合法绑定），probe-3 证实可用；该口径仅在执行面成立，不写入变更集。
- G9-e 的 INFO（`tmp/P81-*` 引用未暂存）为正常提示，⑨ 前 `git add` 即消。

PRE-REVIEW: PASS WITH FIXES

# GATES-EDGE-COVERAGE-RECHECK 切片验证报告（纯验证片：两类根因复评）

> 切片：`GATES-EDGE-COVERAGE-RECHECK`（v1.37，M 档，纯验证 / 判据片）。基线 HEAD = `aa07eb4`；本片未提交（§9）。

## 1. 变更概要

- **唯一变更集**：`tools/gates/selftest.js`（新增 479 条新增行，删除 0；判据实现 `tools/gates/lib/criteria/**` 与判据文本**零改动**）。
- **新增反例集**：`T78` = 16 条（OB-90 反例）、`T79` = 50 条（OB-91 表驱动反例）、`T79 canary` = 10 条（常量冻结 + 黑盒分词 + 7 个外部终止字符）。
- **自测读数**：`460 → 536`，`SELFTEST: PASS 536/536`；自洽校验 `460 + 16 + 50 + 10 = 536`。
- **门禁读数**：`index` / `tree` / `ci` 三作用域 `TOTAL_FAIL=0`。
- **版本包**：`v1.36 → v1.37`，仅台账元数据（OB-90 / OB-91 处置列追加复评确认、新增 OB-97、§0.3 日志行、附录 C.0y）；**判据正文零改动**。

## 2. 反例集与读数

| 单元 | 条数 | 能证伪什么 |
|---|---|---|
| `T78` | 16 | OB-90 的旧语义（配对只从 CN 侧出发、要求双侧同在变更集）若回退，`_EN` 单侧变更会重新静默；逐条覆盖单侧、无镜像、不可读侧与指令对 fail-closed |
| `T79` | 50 | OB-91 的旧 token 规则（合法 + 非法并存被放过、`lang:zh_cn` 截断、裸 `lang:` 不可见）若回退即红；表驱动覆盖并存顺序、跨行与 `_EN` 侧 |
| `T79 canary` | 10 | 覆盖虚增与自锁风险：常量冻结证明规则输入不被测试反向定义；黑盒分词证明 `lang:zh` / `lang:zh_cn` / `lang:zh/en` 的黑盒行为；7 个外部终止字符逐条断言 |

- 三块合计 76 条，全部为**纯新增断言**，不改任何既有断言；总数与汇总行互相印证。

## 3. 真实变更集 spot check

| # | 场景 | 期望 | 实际 | 结论 |
|---|---|---|---|---|
| S1 | `_EN` 单侧多一行 | `G3(a.1)` FAIL，exit 1 | `+3/-0 vs _EN +4/-0 MISMATCH`，`TOTAL_FAIL=2` | 符合 |
| S2 | 非指令对、numstat 相等、`_EN` 加粗漂移 | 聚合 `G3(a)` PASS 且**不产生** `G3(a.2)` 记录 | `PASS`，detail `"bold":false`，无 `G3(a.2)`，exit 0 | 符合 |
| S3 | 双侧等量等形 | PASS | 六形状全 `true` | 符合 |
| S4 | 无镜像单语文档 | 不产生记录 | 仅 `GATE REPORT` + `TOTAL_FAIL=0` | 符合 |

- **S2 期望修正说明**：上一轮把期望写成 `G3(a.2) FAIL` 是**期望错误**——`G3(a.2)` 只对 `docs/DELIVERY_DIRECTIVE.md` 这一对记录；非指令对的加粗漂移只体现在聚合 detail。本轮按修正后期望判定，实测一致，未改判据去配合。
- **索引复原证明**：spot 前后 `git write-tree` 同为 `21b695f…`；最终 `git status --porcelain -uall` 仅 `M  tools/gates/selftest.js`。

## 4. 审查链

- **⑤ r1**：`PRE-REVIEW: PASS`（0 Medium+，5 Low）。
- **⑥ r1**：`INDEPENDENT SCAN: FINDINGS`（2 须修 + 3 边界）⇒ 整改 ⇒ r2 `INDEPENDENT SCAN: PASS WITH FINDINGS`（1 Medium：canary 自锁依赖同一运行时对象）⇒ 整改 ⇒ r3 `INDEPENDENT SCAN: PASS`。
- **⑦ r1**：`FINAL REVIEW: FINDINGS`（2 Medium：派单回执缺失 / 计数漂移）⇒ 整改（派单审计 + ③ 终版重跑）⇒ r2 `RE-REVIEW: PASS WITH FIXES`（仅 1 Low：审计快照时序）。
- **⑦ r2 的 Low 闭合**：⑦ r2 指出派单审计快照时序滞后（Low）；已在 ⑧a 之后重跑 `audit-dispatch.js` 覆盖审计文件，表内为最终快照。
- **⑦ 额度 1+2 用尽，未开第 3 读**；r2 的 Low 以报告留痕收口，不作额度例外。

**OB-95 / OB-96 复查（第 1/2 片）**：本片是 **OB-95 / OB-96 复查的第 1/2 片数据点**：① 由 `DeepSeek V4 Pro (deepseek)` 独立执行 1 次（OB-95）；②③⑧ 分别由 `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` 实派单（OB-96）。**两项均未复发**，按处置**继续观察、不升级**；第 2/2 片为下一片 `G9-c`——若不复发 ⇒ 可转 `已处置`，若复发 ⇒ 升级 `待议`。

## 5. 派单自证

| 位 | `agentName` | 模型 | transcript 行偏移 | 次数 |
|---|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `12211` | 1 |
| ② | `prfrail-implementer` | 同主控模型 | `12377` / `13163` / `13367` | 3 |
| ③ | `prfrail-tester` | 同主控模型 | `12633` / `13811` | 2 |
| ⑤ | `independent-reviewer` | `GLM-5.3 (glm)` | `12922` | 1 |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `13074` / `13312` / `13510` | 3 |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `13546` / `13954` | 2 |
| ⑧ | `prfrail-documenter` | 同主控模型 | `14242` | 1 |

- 审计快照：`tmp/GATES-EDGE-COVERAGE-RECHECK/dispatch-audit.txt`（已重跑覆盖，快照含 13 条派单记录）；`14242` 为快照末条，对应 ⑧a 报告与台账落盘；⑧b 定稿为本轮派出、晚于快照，不新增快照行（快照内 `prfrail-documenter` 计数为 1）。
- 逐位计数：① 1 / ② 3 / ③ 2 / ⑤ 1 / ⑥ 3 / ⑦ 2 / ⑧ 1；**无零派单位**。

## 6. 临时件纪律的失效实证

- **事实**：起跑消息已前置「临时件必须落在 `tmp/<切片ID>/`，违反即停下升级」，但 `$env:TEMP` 累计出现 **5 次**同型违反（本片 1 次 + `GATES-G1G5` 4 次）。
- **结论**：**规则文本对 Flash 无效**——同一条硬规则在明写、前置、并附停止条款的条件下仍被重复违反；`G9-c`（临时件纪律机械判据）的立项条件**远超阈值**，下一片立即立项。
- **同型性**：本片主题是「声称已关闭的根因是否真正关闭」；本片自身在临时件纪律上恰好重复了「声称已前置、实际未执行」的同型失败——**本片以自身为反例验证了本片主题**。
- 本片仓库内过程工件全部位于 `tmp/GATES-EDGE-COVERAGE-RECHECK/`，⑧c 逐件清除；仓库外残留**无机械证据**（如实登记）。

## 7. 已知缺口与边界

- **⑤ 的 5 条 Low**：均为措辞与覆盖建议，不改变判定；按 r1 结论留痕。
- **⑥ 的 canary 边界**：r2 的 Medium 指出 canary 初版依赖同一运行时对象（自锁）；r3 改为常量冻结 + 黑盒分词后关闭，但「常量与规则同一文件」的结构性耦合仍在，属**已登记边界**。
- **① 的报告级候选（只在报告登记，不进 §12.4）**：`G3(a)` 作用域不对称（`changedMarkdown` 空集判定与 `changedPaths` 循环的差别，证据树 `_EN` 在混合变更集中被判定）；标题**层级**漂移未覆盖（只看有无标题）；dedup 顺序导致的跳过；`LANG_HEADER_WINDOW` 窗口外的非法标签不可见。
- **OB-97**（`观察中`）：`lang:zh:en` / `lang:zh:` 在第二冒号处截断后声明 `zh`、未抛 `UsageError`；锚点 `tools/gates/lib/criteria/g6.js` 行号约 735；§6.3「凡表达式」句属字面过度承诺 ⇒ **本片只登记、不修**。
- **试点顺延**：第 **10** 次（本片为验证 / 判据片，命中 §12.1 配套②），本片**未新增顺延 OB**。
- **⑦ r2 的 1 条 Low 的闭合方式**：⑦ r2 的 1 条 Low（派单审计快照时序滞后）**由主控自证闭合**（重跑审计脚本覆盖快照），**非独立复审**；⑦ 额度 1+2 用尽，按「额度用尽」口径停在该轮，**未开第 3 读**。本条属「Low 级可延后」情形，**如实登记、不销解**。

## 8. 成本与计量

| 位 | 计量 | 说明 |
|---|---|---|
| ① | 1 | 设计复评（收窄输入） |
| ② | 3 | 反例集初版 + canary 整改 + 自锁独立性整改 |
| ③ | 2 | 反例执行与 spot check + 终版快照重跑 |
| ⑤ | 1 | r1（额度 1+1，仅用 1） |
| ⑥ | 3 | r1 / r2 / r3（必派，§7.2.1 ②） |
| ⑦ | 2 | r1 / r2（额度 1+2，用尽） |
| 批量辅助 | 0 | 四类用途无适用任务 |
| 探针 | 0 | 未使用探针额度 |

**v1.35 分层规则的第 2 片实证**：本片 ⑦ 1+2 用满（r1 `FINDINGS` → r2 `PASS WITH FIXES`，仍剩 1 Low）；⑥ 3 次用满（r1 `FINDINGS` → r2 `PASS WITH FINDINGS` → r3 `PASS`）。与 `GATES-G1G5`（L 档 ⑦ 1+3 用满）同向 ⇒ **分层规则对判据 / 验证类切片可能偏紧**，末轮仍在产出 Low / Medium+。**本片不升级、不改 §7.5**，登记为调节观察数据点；若 `G9-c` 片再现同型，可评估是否单独评估 §7.5 的分层。

## 9. 提交与 CI

- 本片**未提交**（⑩ 停点）；提交与 CI 读数由 ⑩ 同轮授权后回填。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-edge-coverage-recheck.md	repo
docs/validation/directive-gates-edge-coverage-recheck_EN.md	repo
docs/validation/evidence/GATES-EDGE-COVERAGE-RECHECK-freedom-list.md	repo
tools/gates/selftest.js	repo
```

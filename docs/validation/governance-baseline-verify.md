# 证轨 P47 治理可信基线核验报告：断点 B 七项核验与历史补登记（准则 v1.60）

切片 `GOVERNANCE-BASELINE-VERIFY`（评估/治理类 `[SLICE]`，流水线裁剪：0 → ① → ② → ⑧ → ⑨ → ⑧c → ⑩，跳过 ③–⑦）；起跑锚点 `HEAD` = `0e46bd0`；准则 v1.60；报告日期 2026-11-02；R / D 基线 = 0 / 0；本片实测 R = 2（两次门禁判红后的修复重跑，同一轮多判据合并计 1）、D = 0。

## 1 执行摘要

七项二选一结论：**三项补证后升级**（G4 / G5 / G7），**三项维持未闭合**（P3 / P5 / S7），**M6 随 G4 口径升为补证后升级**。

补登记落地：`§12.4` 新增 OB-133 至 OB-136；`REMAINING_SLICES` 的 P32 / P33 / P34 三行证据格追加补登记指针；本片台账行与详情块新增。**历史报告正文与已完成行零删改。**

关键边界必须与结论同时读：三项「补证后升级」指**记录链缺口被补登记**，**不撤销** P32 / P33 / P34 的已闭合结论、**不重开**其验收。

## 2 R-11 裁定与纠偏范围

用户裁定 = **B 档 + 逐项差异化**：G4 只登记现值与首次触线时点更正（不重算历史评分）；G5 补立 OB 声明「独立留痕不成立」（不回写 P33 正文）；G7 补登记「特例未成立」（排除条件是否放宽另片承载）。四条红线：只新增不删改、CN 与 `_EN` 同步、`G6-1` 时态一致、⑨ 全绿后必停 ⑩。

## 3 断点 B 七项核验结论

| 项 | 结论 | 机械依据 | 缺口（维持者） |
|---|---|---|---|
| P3 | 维持未闭合 | `PRODUCT_REQUIREMENTS` 的 REQ-023 为「未知外部副作用 ⇒ 不得推进链」且无例外；`TEST_STRATEGY` 的 AT-24 引入「原理不可观测者标 `unknown`」例外 | 两处未互指；例外权威归属未定（R-09 待裁） |
| P5 | 维持未闭合 | CI 中候选由 runner 直接 `go build`（步骤 `Build candidate without executing it`）；seed 二进制只参与 `selfhost` 验证，候选来自 `download-artifact` | RFC 的 R19 要求「seed 构建候选 N+1」；缺该步证据，补证归属层 4 |
| G4 | 补证后升级 | 逐次介入事件可枚举，复算表见证据包；P32 累计式漏计前两环 | — |
| G5 | 补证后升级 | `§1.7.4` ② 明文「以 OB 或 ADR 独立登记（不得只写 `§0.3` 或附录 C 的附注）」，P33 写作「独立留痕 = 本段」 | — |
| G7 | 补证后升级 | `§5.1` 特例明文「位置合法而仅次序违规的，不适用本特例」，P34 的 D 事件正落排除侧 | — |
| M6 | 补证后升级 | 与 G4 同源：按累积计数，`§9.2` 阈值在 ⑨ 期第二项后即达成 | — |
| S7 | 维持未闭合 | OB-21 条目内「责任主体 / 停止条件 / 复议触发点」三串均 0 命中 | 三要素缺；由 OB-136 新增承载 |

检索范围与逐项判据见证据包 `01-baseline/verification.md`；G4 / G5 / G7 的更正正文见 `01-baseline/correction.md`。

## 4 补登记落地

| 载体 | 内容 | 保留原值 |
|---|---|---|
| `§12.4` OB-133 | P32 的 R 计数现值与首次触线时点更正（G4） | 不回写 P32 报告 |
| `§12.4` OB-134 | P33 独立登记缺口的补登记（G5） | 不撤销 P33 已闭合结论 |
| `§12.4` OB-135 | P34 特例适用判定更正（G7，观察中） | 不改 `§5.1` 正文 |
| `§12.4` OB-136 | OB-21 三要素补齐（S7） | 不删改 OB-21 原条目 |
| `REMAINING_SLICES` P32 / P33 / P34 | 证据格尾追加补登记指针 | 状态、依赖、规模、完成日期、证据路径 |

指针追加在证据格**内部**而非新增单元格，以保持状态表列数不变。落地明细见 `02-backfill/registration.md`。

## 5 流水线裁剪与工序声明

- 步骤 0：落 `00-briefing.md`（起跑消息逐字）、`input-list.md`、`freedom-list.md`、`MANIFEST.md`、`SHA256SUMS.txt`。
- ① 基线核验：只读检索 + 机械判据复算，产出 `01-baseline/` 两件。
- ② 补登记：只新增条目与指针，产出 `02-backfill/registration.md`。
- ⑧ 装配报告：本文件与其 `_EN` 镜像。
- ⑨ 三作用域与 SELFTEST；⑧c 清洁；⑩ 停点等用户授权。
- **③–⑦ 跳过**：本片无实现对象，④ 的对象是代码 diff、⑥ 的扫描对象是实现 diff、⑦ 的适用前提是有实现审查对象。**独立性边界如实声明**：本片结论由主控自核 + 只读批量辅助承担，**未经 ⑥ / ⑦ 独立审查**。
- ⑥ 跳过：本片无实现 diff，⑥ 的扫描对象不存在；结论由主控自核与只读批量辅助承担，独立性边界已在本节如实声明。

## 6 模型与计量

| 单 | 载体 | 模型 | 输出字符 | 派单上限 |
|---|---|---|---|---|
| 1 台账枚举 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 2,972 | 3,000 |
| 2 引用点抽取 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 3,480 | 3,500 |

两次派单**严格串行**（第 1 单回执后才发第 2 单），符合 R1.1 / R5.1 的「禁止并行子代理」纪律；两单均为**只读**，无写操作。主控直核 G4 / G5 / G7 与 P5 的 workflow 事实由主控自身只读完成，未派写操作子代理。

## 7 门禁读数

| 项 | 读数 |
|---|---|
| `gate.js --all --scope=index / tree / ci` | 三作用域均 `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 697/697` |
| `G3(a)` 双语对称 | 准则对 +4 / −0 与镜像 +4 / −0 相等；台账对 +25 / −3 与镜像 +25 / −3 相等 |
| `G6-1` 状态时态 | 新增行无待办 / 完成冲突 |
| `G6-4` 期限字面 | 首跑判红一次（台账条目里该标识后的量词被读作行数字面量），改词后转绿；属本片过程事实 |
| `G6-5` / `G6-6` | 节引用可解析；新增台账行首状态词均属闭集 |
| `G7` 报告契约 | 本文件与其镜像均在 G7 扫描域内，围栏齐备 |
| `G8` 证据树 | 包自洽（`MANIFEST.md` 条目与 `SHA256SUMS.txt` 哈希行同步） |
| 六形状（CN / `_EN`） | 行数 113 / 113；标题 9 / 9；`**` 32 / 32；竖线 150 / 150；空行 28 / 28；反引号 184 / 184 |

## 8 提交与 CI

本片按「评估类切片 ⑨ 后必停 ⑩」纪律停在 ⑩，经用户同轮授权后提交与推送（仅 `origin/main`，未推 `gitee`）。

| 稿件 | 提交 | run | 结果 |
|---|---|---|---|
| 主提交（本片全部载荷） | `bbd1745` | `41358624169` | 双绿：`ubuntu-latest` 与 `windows-latest` 均 success；两腿 `SELFTEST: PASS 697/697` 与 `TOTAL_FAIL=0` |
| 本节回填（docs-only 追加） | 追加提交时登记 | 追加提交时登记 | 追加提交时登记 |

包级读数（最终 CI 证据）：`ubuntu-latest` 的 `internal/adapters` 11.959s、`internal/applier` 0.492s、`internal/chain` 0.565s、`internal/console` 1.912s、`internal/evidence` 0.009s、`internal/gates` 0.126s；`windows-latest` 的 `internal/adapters` 29.364s、`internal/applier` 5.253s、`internal/chain` 3.912s、`internal/console` 15.398s、`internal/evidence` 0.057s、`internal/gates` 0.526s、`internal/guard` 2.655s、`internal/release` 3.818s、`internal/repair` 0.028s、`internal/snapshot` 1.320s、`internal/taskdef` 0.062s、`internal/tickets` 4.359s。两腿均 `SELFTEST: PASS 697/697` 与 `TOTAL_FAIL=0`；本节回填为 `docs-only` 追加，不改代码与判据。

## 9 已知限制

- P5 的判据含一处解释分歧：RFC 的 R19 写「seed 构建候选」，CI 现状为「seed 验证候选、runner 构建候选」。本片按字面判为缺证据；若用户认定「seed 验证」已满足 R19，P5 可改判，**该分歧留待裁定，本片不自行放宽**。
- 三项「补证后升级」只解决记录链，不构成对历史评分或技术结论的复核。
- 未派 ⑥ / ⑦，结论未经独立审查层复核。
- 本片不改 `§5.1` 排除条件，G7 的「特例未成立」只作事实登记。
- 台账中新增条目不写数值与行锚，复算细节全部落在证据树，以免台账出现会过期的字面量。
- 报告 CN 与 `_EN` 的六形状六项全等（行数 / 标题 / `**` / 竖线 / 空行 / 反引号），结构对称判据通过。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/CHECKLIST_v37.md	repo
docs/t027/CHECKLIST_v37_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/00-briefing.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/01-baseline/correction.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/01-baseline/verification.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/02-backfill/registration.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/MANIFEST.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/SHA256SUMS.txt	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/freedom-list.md	repo
docs/validation/evidence/GOVERNANCE-BASELINE-VERIFY/input-list.md	repo
docs/validation/governance-baseline-verify.md	repo
docs/validation/governance-baseline-verify_EN.md	repo
```

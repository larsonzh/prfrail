# P78 输入清单（§13.1.1）

切片：P78 `RELEASE-ZIP-BINDING`（`[SLICE]` · M-L 实现类）
基线：HEAD = origin/main = `c935713`；准则 v1.69；SELFTEST 737/737；三作用域 `TOTAL_FAIL=0`

## 一、核心判定对象

| # | 文件 | 用途 |
|---|---|---|
| 1 | `.github/workflows/ci.yml` | `candidate-build` / `candidate-probe` / `selfhost` 三个 workflow_dispatch job 与九个 inputs（Q1 / Q3 路径判定与输入可满足性；**本片零改动**） |
| 2 | `cmd/prfrail-release/main.go` | 发行工具子命令 `hash-tree` / `probe` / `selfhost` / `generate` / `verify`（Q1 / Q2 证据生产方式） |
| 3 | `internal/release/**` | `GenerateBundleMetadata` / `VerifyReleaseBundle` 的产物集与校验口径（Q2 发布证据） |
| 4 | `testdata/selfhost/{bootstrap,expected,licenses-policy,noop-chain,executable-chain}.json` | workflow_dispatch 输入的取值来源（bootstrap seedCommit 等） |
| 5 | `docs/INSTALLATION.md` / `INSTALLATION_EN.md` | 安装模型、ZIP 结构与核验命令、`:206` 三处绑定句（Q2 / Q4 / Q5） |
| 6 | `docs/S1_RELEASE_NOTES.md` / `S1_RELEASE_NOTES_EN.md` | 绑定落点之一 + 「发布前剩余门禁」段（Q2 / Q3 / Q5） |
| 7 | `docs/S1_SUPPORT_MATRIX.md` / `S1_SUPPORT_MATRIX_EN.md` | 绑定落点之一；平台清单与支持层级零改动（Q2 / Q5） |
| 8 | `docs/validation/s1-exit.md` / `s1-exit_EN.md` | 绑定落点之一（非门表区）+ 剩余风险第 1 条（Q2 / Q5） |

## 二、第二层候选（视 ① Q5 实读结果判定）

| # | 文件 | 用途 |
|---|---|---|
| 9 | `docs/ADR_REGISTER.md` / `ADR_REGISTER_EN.md` | `:20` ADR-009 状态列闭括号；`:62` 「当前仍不是已发布 SLA」默认不翻（Q5 敏感点） |
| 10 | `docs/INSTALLATION_PLAN.md` / `INSTALLATION_PLAN_EN.md` | ZIP 文件名与下载入口冻结状态（Q5） |
| 11 | `docs/PRODUCT_REQUIREMENTS.md` / `PRODUCT_REQUIREMENTS_EN.md` | §7 发布相关待落实句（Q5） |
| 12 | `docs/DEV_PLAN.md` / `DEV_PLAN_EN.md` | T017 / T018 发布里程碑状态（Q5） |
| 13 | `docs/OPERATIONS.md` / `OPERATIONS_EN.md` | 下载 / 校验 / 升级 / 回滚运维说明（Q5） |

## 三、门级与状态依据

| # | 文件 | 用途 |
|---|---|---|
| 14 | `docs/validation/p57-closeout-assessment.md` / `_EN.md` | §4.2 S-G3 判据原文（③ 一致性基准；**本片零改动**） |
| 15 | `docs/validation/s1-exit-closeout-assessment.md` / `_EN.md` | P76 报告对：Q5 四步、Q10=(C) 分两步、ADR-046（Q1 / Q3 / ③） |
| 16 | `docs/t027/CHECKLIST_v61.md` / `_EN.md` | 当前快照 §四 P78 队列项（**本片零改动**；其 ⑦ 额度引用误差由本片台账登记披露） |
| 17 | `docs/DELIVERY_DIRECTIVE.md` / `_EN.md`（v1.69） | §7.5 额度、§5.2 批次、§9.2 高风险、§11.2 回填、§13 输入清单（**本片零改动**） |

## 四、执行纪律与台账

| # | 文件 | 用途 |
|---|---|---|
| 18 | `docs/t027/REMAINING_SLICES.md` / `_EN.md` | P77 段格式基准 + P78 段登记（步骤 0c）与完成记录（⑧b） |
| 19 | `docs/CODING_CONVENTIONS.md` / `_EN.md` | 编码 / 行尾 / Git 纪律（本片落盘执行） |
| 20 | `docs/validation/evidence/P78-RELEASE-ZIP-BINDING/00-briefing.md` | 起跑消息原文 |
| 21 | `docs/validation/evidence/P78-RELEASE-ZIP-BINDING/freedom-list.md` | 执行自由度清单（⑥ 后 / ⑦ 前定稿） |

## 五、取证方式说明

- 行号引用以**实读为准**；行号漂移时以「文件名 + 节号 / 段落标题」为主锚点。
- Q6 组合 grep 模式集（独立执行、禁止单关键词联查）：`RELEASE-ZIP` / `GitHub Release` / `v0.1.0`。
- 构建与复验的原始捕获件全部落 `tmp/P78-RELEASE-ZIP-BINDING/`，不入库。
- 全部输入均为**只读**；① 阶段不改任何文档。

# P97 输入清单（input-list）

切片：`HYDROGEN-CONTROL-PILOT`（P97）。起跑日：2026-11-28。主控：`DeepSeek V4.1 Flash (deepseek)` High。

## 1. 锚点实测（步骤 0）

| 项 | 期望 | 实测 | 判定 |
|---|---|---|---|
| HEAD | `8f0c42b` | `8f0c42b` | 一致 |
| origin/main | `8f0c42b` | `8f0c42b` | 一致 |
| 准则版本 | v1.81 | v1.81 | 一致 |
| SELFTEST | 806/806 | `SELFTEST: PASS 806/806` | 一致 |
| 作用域 tree | `TOTAL_FAIL=0` | 0 | 一致 |
| 作用域 index | `TOTAL_FAIL=0` | 0 | 一致 |
| 作用域 ci | `TOTAL_FAIL=0` | 0（本地退化为 `HEAD^1..HEAD`，changed=4） | 一致 |
| v78 / v79 / v80 | 211 / 209 / 210 | 211 / 209 / 210（双语同值） | 一致 |
| ADR 表行数 | 48 | 48 | 一致 |
| OB 表末行 | OB-145 | OB-145 | 一致 |
| tmp 顶层 | `.gitkeep` + `gates-cache` | 两件，`gates-cache` 下含 `base` 与 `g4b` | 一致 |
| 工作区 | 洁净 | `git status --porcelain` 为空 | 一致 |

## 2. 阅读的输入（本片对照基准）

| 路径 | 用途 |
|---|---|
| `docs/t027/CHECKLIST_v80.md` / `_EN` | 步骤 0 锚点出处与观察项（P87 Hydrogen 前置缺口） |
| `docs/validation/directive-review-model-policy.md` / `_EN` | §7 注记回填对象；P87 结论与审查链记录 |
| `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/MANIFEST.md` | P87 输入包目录结构 |
| `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/07-final-review-r1.md` | P87 原 ⑦ 终审结论基准（`FINAL REVIEW: PASS`，新增 Medium+ 为 0） |
| `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/00-briefing.md` | P87 起跑口径（轻量档条款来源） |
| `docs/DELIVERY_DIRECTIVE.md` §1.7.4 / §1.7.3 / §2.2 / §5.2 / §6.3 / §7.5 / §12.4 / §13.1.10 / §13.1.12 | 临时授权五要素、自由度清单、轻量档限定名与单价、S 档额度、OB 登记位、留痕与证据同步 |
| `docs/t027/REMAINING_SLICES.md` / `_EN`（P87 段与 P96 段） | 前置缺口登记行与台账段格式 |
| `docs/CODING_CONVENTIONS.md` | `.md` = UTF-8 BOM + LF |

## 3. 本片不读取 / 不适用的输入

| 项 | 理由 |
|---|---|
| `tmp/P87-DIRECTIVE-REVIEW-MODEL-POLICY/review-input/` | P87 收尾已按 §1.5 清理，路径不可得（如实登记，不作缺陷） |
| P96 证据包作对照对象 | 起跑消息指定 P87 输入包为首选；备选仅在 ① 判定另选时启用 |
| `internal/**` / `tools/gates/**` 源代码 | 本片 docs-only，不触判据实现与产品代码 |

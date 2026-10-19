# 切片 DIRECTIVE-EVIDENCE-WINDOW 验证报告（准则 v1.39，MINOR）

> 切片：`DIRECTIVE-EVIDENCE-WINDOW`（原暂定名 `G9-d`）；准则版本 v1.39；档位 M；分级 `[SLICE]`；本文件为 **⑧b 定稿**（§9 于 ⑨ 执行后补入）。
> 基线 HEAD = `caf3b0f`；⑧c 清洁记录、主提交 hash 与 CI run 号由 ⑩ 同轮授权后的阶段回填。
> 双语：本文件为 CN 权威，`docs/validation/directive-evidence-window_EN.md` 为逐节同位镜像。

## 1 变更概要

**当前变更集**（`git diff --cached --numstat`，9 paths，合计 +664/−12；前 5 行来自 ⑧a，末 4 行为 ⑧b 新增（本报告对 + 台账对））：

| 路径 | + | − | 说明 |
|---|---|---|---|
| `docs/ADR_REGISTER.md` | 1 | 0 | ADR-026 登记行 |
| `docs/ADR_REGISTER_EN.md` | 1 | 0 | ADR-026 `_EN` 同位行 |
| `docs/DELIVERY_DIRECTIVE.md` | 23 | 6 | v1.39 准则正文 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 23 | 6 | `_EN` 严格同位镜像 |
| `docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md` | 38 | 0 | 新建；执行自由度清单 |
| `docs/validation/directive-evidence-window.md` | 270 | 0 | 新建；本报告 CN 权威 |
| `docs/validation/directive-evidence-window_EN.md` | 270 | 0 | 新建；逐节同位 `_EN` 镜像 |
| `docs/t027/REMAINING_SLICES.md` | 19 | 0 | P14 定义块 + 完成行 |
| `docs/t027/REMAINING_SLICES_EN.md` | 19 | 0 | `_EN` 同位增量 |

**准则侧落点**（`docs/DELIVERY_DIRECTIVE.md` 与 `_EN` 各 8 个 hunk）：

- 头部版本行：v1.38 → v1.39。
- §0.3 新增 v1.39 记录行。
- §1.5 清理时机括注：窗口内临时件按 §5.1 子款在 ⑩ 前清理。
- §5.1 临时件例外句 + 新增「⑧c 后取证窗口」子款。
- §5.3 的 ⑧c 行半句。
- §12.4 的 OB-99 / OB-100 处置列。
- 附录 B.7.3 ⑩ 回执新增「终态回读」槽位。
- 附录 C.0aa 新台账块。

**变更性质**：MINOR 版本修订，只补 §5.1 的窗口处置规则与取证槽位，不改判据本体、不改流水线结构。

**执行自由度**：本片行使的实质性选择逐条见 §13，内容与 `docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md` 一致。

## 2 执行流水线

| 步骤 | 角色 / 载体 | 模型 | 状态 | 返工次数 |
|---|---|---|---|---|
| 0 双重身份登记 + ADR-026 | 主控 | 主控模型 | 已落盘 | 0 |
| ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` | 0 |
| ② 实现 | `prfrail-implementer` | 与主控同模型 | `IMPLEMENT: DONE`（5 次派单，含 2 轮整改） | 2 |
| ③ 测试 | `prfrail-tester` | 与主控同模型 | r1 + r2 完成 | 1 |
| ④ 门禁与集成 | 主控 | 主控模型 | 完成 | 0 |
| ⑤ 预审 | `independent-reviewer` | `GLM-5.3 (glm)` | r1 + r2 完成（1+1 用尽） | 1 |
| ⑥ 独立扫描 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `CODE-REVIEW: PASS`（盲审） | 0 |
| ⑦ 终审 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` Extra High | r1 + r2 + r3 完成（1+2 用尽） | 2 |
| ⑧a 初稿 | `prfrail-documenter` | 与主控同模型 | 本报告已落盘 | 0 |
| ⑨ 原生验证 | 主控 | 主控模型 | 完成（见 §9.1） | 0 |
| ⑧b 定稿 | `prfrail-documenter` | 与主控同模型 | 本报告对定稿 | 0 |
| ⑧c 收尾清洁 | 主控 | 主控模型 | 于 ⑩ 同轮授权后执行 | — |
| ⑩ 停点 | 主控 | 主控模型 | 已到停点；未 commit / 未 push | — |

- 步骤 0 的 ADR-026 由 ② 落盘（见 §3 偏离 #2）。
- ② 的 5 次派单覆盖：步骤 0 的 ADR-026、v1.39 六处条款、③ 缺陷整改、⑤ 整改、⑦ 整改（末 1 次含 ⑦ 三项）。
- ③ 为 r1 + r2 两次；⑤ 为 r1 + r2；⑦ 为 r1 + r2 + r3。
- 耗时字段本片未采集，标注「历史缺失」，不补造。

## 3 异常与兜底记录

| # | 事件 | 处置 |
|---|---|---|
| 1 | 自由度清单产出晚于 ⑦ r1 / r2（§1.7 要求 ⑥ 后、⑦ 前） | 以 ⑦ r3 收窄复审闭合（占用 ⑦ 第 3 次额度）；清单 §4 偏离 #1 已记时间戳实证。 |
| 2 | ADR-026 由 ② 落盘（非 ⑧） | 落盘角色偏离，不涉授权来源；清单 §4 偏离 #2 登记。 |
| 3 | ③ r1 的 F3（OB-99 锚点段名「已知边界」应为「边界 2」）未修 | 锚点位于 §12.4 现象列，§11.1 禁止追溯改写；⑦ r2 确认接受。 |
| 4 | ③ 的 D12（`_EN` 的 `tmp/<slice-id>` 3 处拼写）未修 | `G3(a.1)` 要求双侧 numstat 相等，单侧修必红；两侧同改需凑数改动，不予采纳；留后续真实改动 CN 对应行的切片。 |
| 5 | 批量辅助 #1 发现六形状 `**` 计数申报值陈旧：包内 `06-six-shapes.md` 硬编码 4284，实测 4298 | CN / `_EN` 相等这一实质结论不变，`G3(a.2)` 的 `bold:true` 亦成立；§8 以实测值登记。 |
| 6 | 批量辅助 #2 的摘录有一处术语误差（把「⑩ 前清理」写成「⑧b 前清理」） | 未采信该句；本报告一律以准则原文为准。 |

- 第 1、3 条为**时点 / 历史行纪律**类偏离，不是实现缺陷；第 4 条是判据约束下的取舍。
- 第 5 条只更正读数，不改变门禁结论；第 6 条只影响摘录，不影响条款文本。

## 4 环境与档位

| 项 | 值 |
|---|---|
| OS | Windows（本机 `d:\LZProjects\prfrail`） |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| 基线 HEAD | `caf3b0f` |
| 档位 | M |
| 分级 | `[SLICE]` |
| ⑥ | 必派（规则语义变更） |
| ⑦ 额度 | 1+2（用尽） |
| ⑤ 额度 | 1+1（用尽） |
| 批量辅助 | 2（`GLM-5.3-Flash (glm)` 数字核对 + `GLM-5.3-FlashX (glm)` 长文档摘录） |
| 探针 | 0 |

## 5 审查结论摘要（⑤⑥⑦）

| 轮次 | 编号 | 级别 | 结论 / 处置 |
|---|---|---|---|
| ⑤ r1 | F1 | Medium | 预取遗漏确认者 ⑤ 时序不可达；整改为「用户在 ⑩ 裁决」。 |
| ⑤ r1 | F6 | Medium | 读数①未覆盖暂存区第一列；整改为增第一列集合比对。 |
| ⑤ r1 | F2 | Low | §0.3 改指针句。 |
| ⑤ r1 | F3 / F4 / F5 | Low | 已整改。 |
| ⑤ r2 | N1 | Low | 「声明的变更集」未锚定唯一真值源；已在 ⑦ r1 由 F-02 升级并修复。 |
| ⑤ r2 | N2 | Low | 「逐一相同」双向性。 |
| ⑤ r2 | N3 | Low | 失效形态枚举未随 F6 同步。 |
| ⑥ | — | — | `CODE-REVIEW: PASS`，零缺陷（盲审，未附 ⑤ 结论）。 |
| ⑦ r1 | F-01 | Medium | OB-100 状态词 `已处置` 与「仅覆盖窗口」范围张力；回退 `观察中` 并写「窗口子范围已处置」。 |
| ⑦ r1 | F-02 | Medium | 「本片声明的变更集」未锚定唯一真值源；读数①锚定 §5.2 G5 的 `declaredFiles`。 |
| ⑦ r1 | F-03 | Low | 窗口内「创建后自清」不可审计；子款补「只判终态可观测性」边界句。 |
| ⑦ r2 | — | — | F-01 / F-02 / F-03 全部闭合，`TERMINAL-REVIEW: PASS`。 |
| ⑦ r3 | F9 + 2 条连带 + 时间戳 | Low | `FREEDOM-LIST: PASS WITH FIXES`，三条意见全部落地。 |

- **⑤ r1** `GLM-5.3 (glm)`：`PRE-REVIEW: PASS WITH FIXES`；整改 F1、F6、F2 后 **r2** `RE-REVIEW: PASS WITH FIXES`，仅剩 3 条 Low（N1 / N2 / N3）。
- **⑤ 额度 1+1 用尽**，未开第 3 轮；N1 已在 ⑦ r1 由 F-02 升级并修复。
- **⑥** `MAI-Code-1.1-Flash (copilot)` 盲审：`CODE-REVIEW: PASS`，零缺陷。
- **⑦ r1** `GPT-5.3-Codex (copilot)` Extra High 盲审：`TERMINAL-REVIEW: FINDINGS`（F-01 / F-02 Medium，F-03 Low）；整改后 **r2** 收窄复审 `TERMINAL-REVIEW: PASS`。
- **⑦ r3** 收窄只审自由度清单：`FREEDOM-LIST: PASS WITH FIXES`，F9 表述、补 2 条连带同步、偏离 #1 加时间戳均已落地。

## 6 审查输入包摘要（含盲审隔离证明）

| 包 | 位置 | 内容 |
|---|---|---|
| ⑤ r1 | `tmp/G9-d/review-5/` | 14 项 + `MANIFEST.md`（含 SHA256 / 字节数），全部存在性断言通过。 |
| ⑤ r2 | `tmp/G9-d/review-5b/` | 收窄包，5 项。 |
| ⑥ | `tmp/G9-d/review-6/` | 14 项；**未附 ⑤ 任何结论**。 |
| ⑦ r1 | `tmp/G9-d/review-7/` | 13 项；**未附 ⑤ / ⑥ 任何结论**。 |
| ⑦ r2 | `tmp/G9-d/review-7b/` | 收窄包，5 项。 |
| ⑦ r3 | — | 只审自由度清单的收窄输入。 |

- **盲审隔离证明（⑥）**：输入包不含 ⑤ 的结论文件与前轮审查笔记；审查者的 scratch 目录与输入包分离。
- **盲审隔离证明（⑦ r1）**：输入包不含 ⑤ / ⑥ 的任何结论；brief 声明「首轮不附前序清单」。
- **首轮是否附清单**：⑥ 与 ⑦ r1 均**未附**自由度清单（清单在其后产出，见 §3 第 1 条）；⑦ r3 为只审清单的收窄输入。
- **输入消毒**：包内只有切片定义、设计、测试报告、变更 diff、numstat、门禁读数与准则条款摘录，不注入主控结论。
- 本片满足 ⑦ r1 所称「首轮不附 ⑤ / ⑥ 清单」。

## 7 可证伪性

- **规则修订的可观测判据（① (e) 两条）**：
- 判据一：删除子款 ⇒ 后续切片报告再现「窗口内重建临时件 / 位置 + 时序违规 / 根因是 §5.1 规则缺口」同形登记。
- 判据二：删除「⑩ 前终态回读」句 ⇒ ⑩ 回执不再强制携带四条读数原文，窗口内「创建位置违规但未提交」形态在提交前无任何机械拦截点。
- **条款自带的失效可观测形态**：`tmp/` 顶层条目 ≠ `.gitkeep`；`??` 或第二列非空 > 0；⑩ 回执未附读数原文。
- **边界**：只判终态可观测性；窗口内「创建后自清」过程形态不可机械观测（OB-99 记录态）。
- **变异检验**：本片不改判据本体，`tools/**` 零改动，故不新增 `selftest.js` 变异用例；回归由全量 `selftest` 承担。

## 8 门禁结果

| 命令 | 结果 |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 552/552` |
| `gate.js --all --scope=index` | `TOTAL_FAIL=0` |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0` |
| `gate.js --all --scope=ci` | `TOTAL_FAIL=0` |
| `gate.js --check=G3-a --scope=index` | PASS |
| `gate.js --check=G4b --scope=index` | PASS（`newCjk=558 unknown=[]`） |
| `gate.js --check=G6-4 --scope=index` | PASS |
| `gate.js --check=G6-5 --scope=index` | PASS（`scanned 28 refs 41`） |
| 六形状（CN / `_EN`，实测两份相等） | 见下表 |

**CN / `_EN` 六形状**（`docs/DELIVERY_DIRECTIVE.md` 对 `_EN` 镜像；实测值）：

| 形状 | CN | EN |
|---|---|---|
| lines（split 口径） | 1823 | 1823 |
| lines（物理行） | 1822 | 1822 |
| `**` | 4298 | 4298 |
| `\|` | 2897 | 2897 |
| 空行 | 442 | 442 |
| 反引号 | 4298 | 4298 |
| 标题 | 122（逐位同序） | 122（逐位同序） |

- 注：`**` 的申报值 4284 为陈旧读数，实测 4298（见 §3 第 5 条）；两侧相等。
- `G4b` 读数 `newCjk=558 unknown=[]`；`G6-4` 读数 `scanned 626`；`G6-5` 读数 `scanned 28 refs 41`。
- 本报告对已在 ⑧b 暂存；上表 `--scope=index` 读数均为暂存后复跑值，已计入本报告对。

**本报告对六形状**（`docs/validation/directive-evidence-window.md` 对 `_EN`；⑧b 定稿后实测，六项全等；括号内为 ⑧a 初稿基线）：

| 形状 | CN | EN |
|---|---|---|
| lines | 270（242） | 270（242） |
| headings | 17 | 17 |
| blank | 54（49） | 54（49） |
| `**` | 71（62） | 71（62） |
| `\|` | 541（485） | 541（485） |
| backticks | 306（256） | 306（256） |

## 9 ⑨ 原生验证与提交状态（⑨ 已补入；⑧c / ⑩ 见 §12）

### 9.1 ⑨ 原生验证

**手法**：纯净 worktree（`git worktree add --detach tmp/G9-d/wt HEAD`，HEAD = `caf3b0f`）+ `git apply --index` 复现暂存变更集（apply exit 0）。

**读数**：

- 三作用域 `--all` 均 `TOTAL_FAIL=0`（exit 0）。
- `SELFTEST: PASS 552/552`。
- `--check=G3-a|G4b|G6-4|G6-5`（`--scope=index`）四项均 `PASS`、`TOTAL_FAIL=0`、exit 0。
- 版本头 `v1.39` 命中；子款标题「⑧c 后取证窗口（v1.39）」命中；`ADR-026` 行命中。
- `wt_status_lines=5`（变更集自包含，不依赖未跟踪的审查包）。

**收尾**：`git worktree remove --force` exit 0、`git worktree prune`；主仓未受影响。

### 9.2 ⑧c 收尾清洁记录

⑧c 于 ⑩ 同轮授权后执行；以下读数均为执行后回读，数字未改写。

1. **全仓编码审计（不只变更集）**：受版本控制文件 `3633`，其中文本 `3544` 逐文件核验 BOM/LF ⇒ **违规 0**；`17` 处命中规则但位于冻结证据包内 ⇒ 按 OB-27 改做自洽性核验：四包 `SHA256SUMS` 条目 `10 / 29 / 74 / 77`，合计 **190/190 通过，缺失 0、不符 0**。
2. **临时件清理（逐件精确）**：`tmp/G9-d/` 下 `75` 个文件逐件 `Remove-Item -LiteralPath` 删除，随后自底向上删空目录并删根目录。回读证明：`Test-Path tmp/G9-d` = `False`；`tmp` 顶层仅 `.gitkeep`；`tmp` 递归文件计数 = `1`；`git status --porcelain -uall` 的 `??` 项 = `0`。
3. **⑨ 的临时 worktree**：已 `git worktree remove --force` + `git worktree prune`；`git worktree list` 仅主工作树。
4. **编辑器侧诊断**：全仓 `get_errors` ⇒ **零错误**。
5. **清理后复跑**：`SELFTEST: PASS 552/552`；`gate.js --all --scope=index|tree|ci` 三次 `TOTAL_FAIL=0`（exit 0）。
6. **时序合规声明**：记录先落盘（本报告 §1–§13、自由度清单、台账 P14）→ 再清理 → 再复跑门禁，符合 §5.1（v1.35）的「清理须在记录落盘之后」要求。
7. **F4 登记**：本片临时件目录为 `tmp/G9-d/`（**暂定名目录**；登记名为 `DIRECTIVE-EVIDENCE-WINDOW`），⑩ 前已按上条清零；该过渡态由 §5.1 子款的读数②③机械兜底。
8. **⑩ 前四条读数（原文须入 ⑩ 回执）**：① `git status --porcelain -uall` 的 `??` 项 = 0、第二列非空项 = 0、第一列非空项集合 = 本片声明的变更集（`9` 个路径）；② `(Get-ChildItem -LiteralPath tmp -Recurse -File -Force | Measure-Object).Count` = `1`；③ `(Get-ChildItem -LiteralPath tmp -Force).Name -join ','` = `.gitkeep`；④ `Test-Path tmp/DIRECTIVE-EVIDENCE-WINDOW` = `False`。

### 9.3 提交与 CI

- **主提交**：`00fbd2d`（`docs: add the post-⑧c evidence-window rule (v1.39)`），`9` 个文件，`686` 插入 / `16` 删除；已推 `origin/main`；**未触 gitee**。
- **CI run**：`39235252911`（workflow `CI`，创建于 `2026-10-19T10:08:16Z`），结论 **success**。
- **双腿结论与时间戳**：`Go ubuntu-latest` **success**（`10:08:19Z` → `10:10:55Z`）；`Go windows-latest` **success**（`10:08:20Z` → `10:12:23Z`）。
- **`Gates` 步读数（双腿一致）**：`SELFTEST: PASS 552/552`、`GATE REPORT scope=ci changed=9`、`TOTAL_FAIL=0`。
- **包级证据**：双腿共 `30` 条 `ok   github.com/larsonzh/prfrail/...` 行（含 Race 步的重复计数），`--- FAIL:` = 0、`FAIL <pkg>` = 0。
- 首跑即绿，无 §6.4 处置、无复跑。
- 回填提交（本小节）为 docs-only，其自身 CI 结果在 ⑩ 后续回报中给出。

## 10 成本与计量

| 位 | 次数 | 模型 / 说明 |
|---|---|---|
| ① | 1 | `DeepSeek V4 Pro (deepseek)` |
| ② | 5 | 与主控同模型（步骤 0 的 ADR-026 + v1.39 六处 + ③ 缺陷整改 + ⑤ 整改 + ⑦ 整改；末 1 次含 ⑦ 三项） |
| ③ | 2 | 与主控同模型（r1 + r2） |
| ④ | 1 | 主控（门禁与集成） |
| ⑤ | 2 | `GLM-5.3 (glm)`（1+1 用尽） |
| ⑥ | 1 | `MAI-Code-1.1-Flash (copilot)`（额度 3） |
| ⑦ | 3 | `GPT-5.3-Codex (copilot)`（1+2 用尽） |
| ⑧ | 2 | 与主控同模型（⑧a 初稿 + ⑧b 定稿） |
| ⑨ | 1 | 主控 |
| 批量辅助 `GLM-5.3-Flash (glm)` | 1 | 数字核对（发现六形状 `**` 计数申报陈旧） |
| 批量辅助 `GLM-5.3-FlashX (glm)` | 1 | 长文档摘录（8 hunk / 23 增 / 6 删，与 numstat 一致；含 1 处术语误差未采信） |
| 探针 | 0 | 本片无探针额度 |

## 11 明确未执行事项

- 未 `git commit`、未 `git push`、未触 gitee。
- 未改 `tools/**`、§6.3、§5.2 矩阵、§7.5、§11.2、§11.3、流水线结构。
- 未落 FlashX 永久白名单（ADR-026 明写清仓即止）。
- 本片未使用探针额度（0）。

## 12 下一步建议

1. ⑨ 原生验证：已完成（编码门禁 + G1–G9 + 引用可解析性 + `selftest`，读数见 §9.1）。
2. ⑧b 定稿：已完成（§9 与台账 P14 已补，报告对与台账对已暂存并复跑门禁）。
3. ⑧c 收尾：按 §5.1 子款清理 `tmp/G9-d/` 并回读证明。
4. ⑩ 停点：等用户同轮授权后再提交与推送，并回填提交与 CI 事实。

## 13 执行自由度

逐条引用 `docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md`（清单为权威，本表不改变其含义）。

| # | 选择 | 落点 |
|---|---|---|
| F1 | 三条候选路径为固定次序：默认不落盘一次性输出；备选一 ⑧c 前预取（能预取即强制）；备选二 ⑧c 后重建（限例外 + 强制次序） | §5.1 子款 |
| F2 | 不新增 CI 判据（不做 `G10`）、不改 §6.3 | ① 判定；§6.3 零改动 |
| F3 | 切片命名改 `DIRECTIVE-EVIDENCE-WINDOW`（不占 G 号） | 台账 P14、报告文件名、自由度清单文件名 |
| F4 | 回读读数用四条固定命令，其中计数采用 `Get-ChildItem -Recurse -File -Force`（显式禁用 git 口径） | §5.1 子款读数② |
| F5 | 读数①增「第一列行集合 = 声明的变更集」，并把「声明的变更集」锚定到 §5.2 G5 的 `declaredFiles` | §5.1 子款读数① |
| F6 | 预取遗漏的确认者定为用户（在 ⑩ 停点裁决），而非 ⑤ 预审 | §5.1 子款备选一 |
| F7 | 附录台账块编号取 `C.0aa`（`C.0z` 已占用） | 附录 C |
| F8 | ⑩ 回执模板（附录 B.7.3）新增「终态回读」槽位 | 附录 B.7.3 |
| F9 | OB-100 状态词保持 `观察中`，处置列补「窗口子范围已处置」与「其它时段仍未机械化」的边界限定 | §12.4 |
| F10 | OB-99 现象列「已知边界」锚点不改（§11.1 禁止追溯改写现象列） | ⑦ r2 确认接受 |
| F11 | `_EN` 的 `tmp/<slice-id>` 3 处拼写不改（`G3(a.1)` 要求双侧 numstat 相等，单侧修必红） | ③ r2 复核确认 |
| F12 | 本片临时件目录用 `tmp/G9-d/`（暂定名），改名后未迁移 | ⑧c 记录 |

**清单 §2 连带机械同步（非自由度）**：⑤ r1 的 F1 / F6 / F2、⑦ r1 的 F-01 / F-02 / F-03 整改均为缺陷修复；§0.3 读数枚举改为指针句；§1.5 清理时机补入窗口内临时件处置；§12.4 的 OB-99 处置列补边界限定。

**清单 §4 已知偏离**：① 自由度清单晚于 ⑦ r1 / r2，以 ⑦ r3 收窄复审闭合；② ADR-026 由 ② 落盘，属落盘角色偏离；③ ③ r1 的 F3（OB-99 锚点段名）未修，§11.1 禁止追溯改写。

```artifacts
docs/validation/directive-evidence-window.md	repo
docs/validation/directive-evidence-window_EN.md	repo
docs/validation/evidence/DIRECTIVE-EVIDENCE-WINDOW-freedom-list.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
```

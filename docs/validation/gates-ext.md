# GATES-EXT 切片验证报告（§6.3 扩展 G6 / G7 与 `tools/gates/` 脚本化）

日期：2026-09-22。状态：提交前稿（␸b 回填提交后证据）。切片：`GATES-EXT`（触发词 `[SLICE]`）。

## 1. 变更概要

本片把两类**台账元数据 / 审查包完整性**缺陷从"人工目视"改为**可机械重跑的判据**，并接进 CI 硬门：

- **§6.3 新增两类判据**：**G6 台账元数据一致性**（G6-1 状态时态、G6-2 成本口径、G6-3 run 号形态、G6-4 会过期的字面量）与 **G7 审查包完整性**（G7-a 工件围栏块、G7-b 路径与处置可解析、G7-c finding 处置单元格、G7-d 停止态标记）。两者与 G1-a / G3 / G4 同口径：**只判变更行，历史行不追溯**。
- **判据的执行体**：`tools/gates/`（单入口 `gate.js` + `lib/criteria/**` + `selftest.js` + `testdata/**` + 两份白名单），**离线、无第三方依赖、退出码 0 / 1 / 2 语义真实**（2 = 用法 / 环境错误）。迁移并取代临时脚本 `tmp/gate.js`（原样保留供对照）。基础 **G1**（三处人工位置互判）与基础 **G5**（变更集契约）**本片不实现**，脚本如实打印 `SKIP`——**不假称已覆盖**。
- **CI 集成**：`.github/workflows/ci.yml` 的 `test` job 增**一个** `Gates` 步（双腿、无 `continue-on-error`、不新增 job 与 Action），并把 `Checkout source` 的 `fetch-depth` 设为 `0`，使 `--scope=ci` 能在**完整历史**上解析变更集而不是退化为整棵工作树。
- **冻结契约**：`internal/release/workflow_test.go` 同步三处（步名表、`test/Gates` 脚本体摘要、`test/Checkout source` 参数）。**第三处为受控扩展，经用户授权**，并在提交信息中注明。
- **范围守卫**（本片自查与两轮独立审查共同逼出的加固）：不可解析的 `range` 端点、失败的 `git diff` / `git status`、结论不明的浅克隆探测，一律为**用法 / 环境错误**（退出 2），**绝不**退化成"空变更集 ⇒ 判据全数 vacuous pass"。

## 2. 流水线执行

| 阶段 | 载体 | 产出 | 结论 |
|---|---|---|---|
| ① 架构前置 | `deep-reasoner` | `tmp/gates-ext/ARCH-DESIGN.md` + `SLICE-DEFINITION.md`（含步 1 / 2 / 3 分步阻断式验收） | 完成 |
| ② 实现 | `prfrail-implementer` | `tools/gates/**`、§6.3 条文、CI 步、冻结契约同步 | ②a–②g 七轮（含 ②c–②f 四轮判据加固） |
| ③ 测试 | `prfrail-tester` | 步 2 / 3 / 4 / 5 独立验证（`tmp/gates-ext/step2..step5/**`） | 完成，含两处**独立反驳**见 §6 |
| ④ 门禁 | 主控 | `TOTAL_FAIL=0`（`tree` 与 `ci` 双绿）+ 冻结契约绿 + 历史回溯红 / 绿两端 | 完成 |
| ⑤ 预审 | `deep-reasoner` | `PASS WITH FIXES`（1 High + 3 Medium + 4 Low + 2 Note） | 已整改，见 §4 |
| ⑥ 独立扫描 | MAI-Code-1.1-Flash ×1 | `FINDINGS` 2 条（1 成立并整改、1 拒收） | 已闭合 |
| ⑦ 独立终审 | GPT-5.3-Codex ×3 | 第 1 轮 `FINDINGS`；第 2 轮 1 Medium + 1 Low 且判首轮 High 为 `PARTIAL`；第 3 轮 `PASS WITH FIXES`、无 Medium+ | 已闭合 |

## 3. 三个硬发现留痕（① 阶段）

架构前置分析**实测**（非推定）出迁移前的三处覆盖缺口，逐条登记为 **OB-15**：

| 编号 | 硬发现 | 实测方法 | 本片处置 |
|---|---|---|---|
| ① | **硬编码仓库根**（`D:/LZProjects/prfrail/`）使旧临时脚本在 Linux CI 上不可用 | 读取 `tmp/gate.js` 常量 + 在非该路径下执行 | 改为 `git rev-parse --show-toplevel` 发现仓库根 |
| ② | **退出码恒 0**，且对**未跟踪 / 仅暂存**文件静默漏检 ⇒ 此前每一次"门禁全绿"结论对 G2 / G3(b) / G4b / G1-a 的覆盖都不完整 | 构造未跟踪与仅暂存夹具，观察旧脚本零命中且退出 0 | `--scope` 覆盖未跟踪与仅暂存文件，退出码语义真实 |
| ③ | **基础 G1 与 G5 从未实现** | 逐条比对 §6.3 条文与脚本判据清单 | 继续 `SKIP`，并在"脚本化状态"段如实登记 |

> ② 的历史影响**不做追溯改写**：既有"全绿"结论按其当时的真实覆盖范围解读。

## 4. 审查摘要与逐条处置

### 4.1 ⑤ 预审（`deep-reasoner`）`PASS WITH FIXES`

| 项 | 级别 | 要点 | 处置 |
|---|---|---|---|
| ⑤ F1 | **High** | `Checkout` 无 `fetch-depth` ⇒ 浅克隆下 `HEAD^1` 缺失 ⇒ `--scope=ci` 退化为整树 ⇒ CI 首跑必红 | 已修（方案 A：`fetch-depth: 0` + 冻结契约第 3 处同步） |
| ⑤ F2 | **Medium** | G7-c 以子串而非单元格判定 | 已修 |
| ⑤ F3 | **Medium** | G7-d 的 token 集仅中文 | 已修（CN / EN 并集，含大小写不敏感） |
| ⑤ F4 | **Medium** | ②f 缺**独立变异证据** | 已修（主控自行产出两向变异证据并证明字节还原） |
| ⑤ F5–F8 | **Low** | 英文重复句、步 1 计数漂移、命令注入 `^` 吞字符致空样本等 | 已修 / 已如实标注 |
| ⑤ F9–F10 | Note | ③ 已接受的决定维持；摘要未独立重算 | 记录在案 |

### 4.2 ⑥ 独立扫描（MAI-Code-1.1-Flash ×1，首轮不附我方清单）

| 项 | 级别 | 要点 | 处置 |
|---|---|---|---|
| ⑥ F1 | **Medium** | G6-1 完成标记正则要求粗体，与准则**枚举**及**准则自带示例行**（`（实得）`）矛盾 ⇒ CI 阻断型判据漏检 | 已修（`DONE_RE` 取准则枚举的**并集**六形态，同标签 40 字符窗不变；新增 4 条可证伪断言 + 1 条过宽防护断言 + 两向变异证据） |
| ⑥ F2 | **Medium** | 认为 G7-b 在 tree 作用域未校验"新入索引" | **拒收**：代码逐字实现条文（`git ls-files --error-unmatch` + 磁盘存在）；已入库工件本就属"在仓库中"，判 `repo` 通过**语义正确** |

> ⑥ 的独立价值再次实测：它在**不附清单**的条件下报出了我方自审未发现的一条真缺陷（⑥ F1）。

### 4.3 ⑦ 独立终审（GPT-5.3-Codex ×3）

| 轮次 | 结果 | 要点 | 处置 |
|---|---|---|---|
| 第 1 轮 | `FINDINGS` | **High**：`range` 端点不校验存在性、`git diff` 失败被读成"无变更" ⇒ 空变更集 ⇒ 判据全数 vacuous pass ⇒ **退出 0**；**Medium**：缺无效 range 负例断言 | 已修（端点 `rev-parse --verify --quiet <ref>^{commit}`；diff 失败即环境错误；新增 4 负例 + 1 正向对照） |
| 第 2 轮 | `FINDINGS` | 判第 1 轮 High 为 **`PARTIAL`**（浅克隆探测自身未校验成功）；新报 **Medium**：tree / index 的 `git status` 失败仍被读成空变更集；**Low**：同上 | 已修（`git status` 失败即环境错误；浅克隆探测必须**结论明确**，否则亦为环境错误；新增 6 条独立断言与逐守卫变异证据） |
| 第 3 轮 | `PASS WITH FIXES` | 首轮 / 第 2 轮各项**全部 CONFIRMED FIXED**、覆盖缺口 **CLOSED**；仅报 **1 条 Low**（`g5a.js` 的错误分流） | 该 Low **未修**：G5-a 属 G5 家族，本片预授权红线含"顺手修 G1 / G5" ⇒ 登记为 **OB-20** 待裁决 |

> **⑦ 必须点名的两向覆盖**均已落实：**「已审 ②f 前导边界」** 判定 `sound`；**「已审 NF-1 措辞」** 判定 `sound`。第 3 轮另确认 G6-1 完成标记并集与准则文本一致。

## 5. 判据自指实测

本片的判据**住在它们自己判的文件里**，因此"判据抓到自己"是预期内的信号，如实记录两例：

| 现象 | 触发点 | 根因 | 处置 |
|---|---|---|---|
| §6.3 残余误报注记被 G6-3 判红 | `docs/DELIVERY_DIRECTIVE_EN.md` 的两行 | 排除标记集为**中文 token**（`不得残留` / `示例` / `引文` / `冲突` / `矛盾`），英文镜像逐字引述被禁形态却不带标记 | 在英文行内保留 `示例` 标记（同 §6.3 英文示例行既有惯例）；登记 **OB-18** |
| 尖括号占位规则误报 CLI 用法示例 | `docs/RFC-proofrail-unattended-ai-engineering-product.md` 的两行 | 规则把 `prfrail run <chain-file>` 当成 run 号占位 | 因"只判变更行"当前不影响门禁；**本片不修**（收窄属准则语义变更，超本片授权）；登记 **OB-17** |

> 两例均**未被掩盖**：未加白名单、未放宽判据。判据对自己的输出保持可证伪。

## 6. ③ 阶段的两处独立反驳（任务书预期被推翻）

| 任务书预期 | ③ 的独立结论 | 含义 |
|---|---|---|
| G4b 在新文档上应当转绿 | **推翻**：新文档必然带入未在白名单的 CJK 字 ⇒ 判据对任何新文档不可满足；改为"差集 ⊖ 白名单须为空"并保留人工登记流程 | 判据设计必须允许**登记一次**的人工判断，否则它对新内容恒假 |
| `tree` 与 `range` 作用域在结论上永不分歧 | **推翻**：`tree` 读工作树、`range` 读端点修订 ⇒ 白名单与工件解析的**来源不同**，结论可合法分歧 | 由此逼出 F1（白名单来源）与 F2（tree 作用域 G7-b 语义）两处修复 |

## 7. 证据强度警示（O5 夹具方案）

OB-11 相关的回归夹具是**夹具级**证据：它证明判据在该形态上可判红 / 可判绿，但**不证明**判据在真实历史文档上等价成立。本片因此**同时**保留两类证据：

- **夹具级**（可重复、可变异）：`tools/gates/testdata/**` + `selftest.js`（断言全绿）。
- **真实历史级**（不可伪造）：`--scope=range:b738e6b^..b738e6b` 判 G6-2 红、`--scope=range:1b01115^..1b01115` 判 G7-a 红，修复侧 `--scope=range:9f68a52^..9f68a52` 绿。

> 结论按**较弱的一类**表述：夹具级证据不能替代真实历史级证据；本片对历史缺陷的可复现性以真实 range 为准。

> **证据保留**：本片 `tmp/gates-ext/**`（step2–step5 历史、变异证据、复审输入包）**保留**，不入库、不占版本控制；**保留期限 = 下片稳定运行后评估**（用户 2026-09-22 裁决）。

## 8. 门禁结果

| 检查 | 结果 |
|---|---|
| `selftest.js`（仓内断言） | 全绿 |
| `gate.js --all --scope=tree` / `--scope=ci` / `--scope=index` | `TOTAL_FAIL=0` |
| 历史真实缺陷回溯 | G6-2 与 G7-a 在对应 range 上判红，修复侧绿 |
| `go test ./internal/release/`（冻结契约） | ok |
| `go test ./...` | 全包 ok |
| 变异检验 | 每条判据 / 守卫**单独**回退 ⇒ 对应断言变红；还原 ⇒ 字节一致 |

## 9. 成本与计量

| 项 | 内容 |
|---|---|
| 预算（启动卡声明） | ①×1 + ⑤×1 + ⑥×1 + ⑦×2 |
| 实得 | ①×1 + ⑤×1 + ⑥×1 + ⑦×3（⑦ 第 3 轮为**用户预授权**的例外，登记见 §12.4 的轮次例外注记） |

## 10. 未执行事项

- **基础 G1 / G5 仍未脚本化**（脚本如实 `SKIP`）：本片范围之外，属后续 `[SLICE]` 候选。
- **OB-20（`g5a.js` 错误分流）未修**：理由见 §4.3，待用户裁决。
- **OB-17 未修、未加白名单**：理由见 §5。
- **⑧b 回写已执行**：本报告 §11 的"提交后"证据与台账状态行均在推送后回填并推送。

## 11. 提交后证据（␸b 已回填）

| 项 | 结果 |
|---|---|
| 提交 | `8a2841f`（推送 `origin/main`；未推 gitee） |
| CI 首跑 | run `35718884321`（head `8a2841f`）：**双腿 `success`**，`Gates` 步在双腿均**实际执行**（step 9），两条腿均报 `GATE REPORT scope=ci … changed=40` 与 `TOTAL_FAIL=0` ⇒ `--scope=ci` 在真实 CI 的完整历史上解析出了变更集（**既非空集、也未退化为整树**），⑤ F1 的 `fetch-depth` 修复与新增硬门在真实 runner 上同时成立 |
| `base go test` 内 flake | **DR-6**：`internal/gates/TestGuardExecutorCleansDescendantsAfterParentExit` 失败 1 次、隔离重跑全部为绿、`base go test` 内亦绿 ⇒ 按 DR 处置口径**边界①**单独登记（见 `docs/t027/REMAINING_SLICES`）；**首跑双腿绿，未再现** |

## 12. 边界标注

> **本片边界说明**：OB-19 与本片主题（§6.3 扩展 G6 / G7）无关，属用户在本片期间临时追加的观察项登记，不构成 `GATES-EXT` 交付物。若审阅者质疑"为何本片含无关内容"，本标注即为解释。

## 13. 下一步建议

1. **␸b 已完成**：本报告 §11 已回填提交号 / CI 首跑结论 / DR-6 观测，台账状态行已同步。
2. **用户裁决 OB-20**：在"修 G5-a 错误分流"与"并入下次 §6.3 修订"之间二选一。
3. **下一片候选**：基础 **G1 / G5** 的脚本化（OB-15 ③ 的闭合项），以及 OB-17 / OB-18 的判据收窄评估。

```artifacts
tools/gates/gate.js	repo
tools/gates/selftest.js	repo
tools/gates/lib/ctx.js	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/lib/criteria/g7.js	repo
tools/gates/lib/criteria/g5a.js	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/g6-whitelist.txt	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
.github/workflows/ci.yml	repo
internal/release/workflow_test.go	repo
tmp/gates-ext/ARCH-DESIGN.md	local-only
tmp/gates-ext/SLICE-DEFINITION.md	local-only
tmp/gates-ext/review/GATE-all-tree.txt	local-only
tmp/gates-ext/review/GATE-selftest.txt	local-only
tmp/gates-ext/review/GATE-range-historical.txt	local-only
tmp/gates-ext/review/MUTATION-f1-donere.txt	local-only
tmp/gates-ext/review/MUTATION-f1-range-vacuous.txt	local-only
tmp/gates-ext/review/MUTATION-guards-independent.txt	local-only
```

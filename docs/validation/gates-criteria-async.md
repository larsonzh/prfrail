# P95 GATES-CRITERIA-ASYNC 验证报告（工具片 · 评估类闭环）

- 切片：P95 `GATES-CRITERIA-ASYNC`（M 档工具片 `[SLICE]`，原目标 = G6 / G7 / G8 异步化）；主控 = `DeepSeek V4 Pro (deepseek)` Max（用户 2026-11-27 裁定，不得作先例）。
- 日期：2026-11-27。准则版本：v1.77（**不升版**）。R = 0、D = 0（起跑基线）。探针 0。
- 证据包：`docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/`（起跑消息、① 方案书、⑤ / ⑥ / ⑦ 审查原文、批量辅助记录、自由度清单、输入清单、`MANIFEST.md` 与 `SHA256SUMS.txt`）。
- 形态：**评估类闭环**——① 收益论证闸判定不做，用户 2026-11-27 裁定采纳候选 A（该裁定**不得作先例**）；判据代码零改动。

## 1. 结论

本片不做 G6 / G7 / G8 异步化。① 影响判定与收益论证闸给出否决结论：ctx 读 API 全同步，Node 单进程下 async 化不能令同步 I/O 重叠，并发收益理论上界 ≈0、实际 ≤250 ms/5.8 s 墙钟。用户裁定采纳候选 A，本片退化为评估类闭环：**判据代码零改动、准则不升版（v1.77 不变）**；P88 报告 §7 建议 4 在本片登记为「**已评估、收益不成立、不实施**」；`CHECKLIST_v78` 双文件不动，由 v79 快照吸收。当前状态 = 已推 `origin/main`，CI run `45939437180` 双腿 success（见 §10）；回填提交后停于 ⑩（仅 origin）。

## 2. 变更集与证据

- 判据实现面零改动：`tools/gates/lib/criteria/g6.js`、`tools/gates/lib/criteria/g7.js`、`tools/gates/lib/criteria/g8.js`、`tools/gates/gate.js`、`tools/gates/lib/ctx.js`、`tools/gates/selftest.js` 均与 HEAD `22370fe` 逐字节相同，本片变更集不含上述任一文件。
- 台账对：`docs/t027/REMAINING_SLICES.md` 与 `docs/t027/REMAINING_SLICES_EN.md`，P95 段由「未立项」翻「执行中」，闭环时翻「闭环」并附完成记录。
- 报告对：本文件与 `_EN` 镜像。
- 证据包：9 件载荷 + `MANIFEST.md` 与 `SHA256SUMS.txt` 两件元数据（条目集与后者相等）。
- 不触项：准则对、`ci.yml`、`internal/**`、CONTRACTS、`schemas/`、`testdata/`、workflows、G5-b 闭集与模型白名单、ADR 台账 —— 均不在本片变更集内。

## 3. 收益论证（① 转述）

### 3.1 串行构成（实测表 + 净自耗估算）

读数口径：Windows 8 逻辑核、Node v24.17.0、语料热、单样本两轮。

| 判据 / 探针 | `scope=tree` ms | `scope=index` ms |
|---|---|---|
| G1-a（固定开销探针） | 1247 | 1177 |
| G4a | 1675 | — |
| G4b | 1596 | — |
| G6 | 1282–1493 | 796 |
| G7 | 610–642 | 471 |
| G8 | 1803–1905 | 1372 |
| `--all --skip-base` 合计 | 5799 | 3812 |

净自耗估算（扣除探针固定开销，低于探针者按 0 计）：tree 三判据 ≈ 0.6–0.9 s，index ≈ 0.2 s；占各自墙钟 10–15% / ≈5%。`--all` 含 base 时占比 <1%（P88 热态 base ≈ 94 s）。工作性质：G6 / G7 为纯 CPU 正则行扫描，G8 为 sha256 加密与集合比较；热路径残留同步 git 子进程（G7-b 每行一次、G8 每次枚举一次批量忽略检查），async 化不能令其重叠。

### 3.2 并发收益上限（理论上界 ≈0 论证）

理论上界 ≈0：ctx 全部读 API 同步（`readFileSync` / `spawnSync`），无异步 I/O 可重叠；Node 单线程下 async 化只作协作式交错，CPU 总量守恒 ⇒ 墙钟 = Σ 自耗时 + await 微任务开销。即便把 ctx 读 API 一并改为异步，上界 = 三判据异步 I/O 占比 ×（Σ − max）：tree 净自耗中大部分为正则与 sha256 的 CPU 时间，可重叠部分 ≤ 250 ms；index ≈ 0。对 `--all --skip-base` 墙钟的实际节省 ≤ 0.3–4%，对含 base 的 `--all` ≈0。负向可能：微任务调度开销、交错破坏热缓存局部性（±几十 ms），以及快速失败改变现有同步顺序中止的行为。

### 3.3 并发代价（共享缓存并发读安全性 / 报告顺序保持 / 新失败形态）

- **惰性缓存并发读安全**：`bytesCache` / `existsCache` / `trackedCache` / `numstatCache` / `showCache` / `ignoreCache` / `hunkCache` / `treeCache` 均为闭包内 Map，写入在同步调用内原子完成，ctx 无 async 函数 ⇒ 不存在跨 await 的半写窗口；写值幂等（同键同值）。
- **共享 Buffer**：读字节返回同一 Buffer 引用，现有调用方只读（BOM 计数、sha256、解码转字符串）。当前安全，但属未加固假设；改并发后须保留「调用方不 mutate Buffer」的注释纪律。
- **WeakMap**：G8 的 pack 根缓存以 ctx 为键，仅 G8 使用，无跨判据共享问题。
- **报告顺序**：调度器已按 plan 序回填各判据独立输出缓冲；三个判据的 `run` 只写自己的 `out`，判据内零 `console` / `process.stdout`；单判据路径已等待 `run` 返回，兼容 Promise。
- **无判据间副作用依赖**：判据目录无对 ctx 的赋值，唯一磁盘写属 G4b 自身缓存（不在本片范围），G6 对 G7 的模块级复用为只读常量与函数。
- **新失败形态**：async 抛错经快速失败，改变现有「同步顺序中止」形态（退出码语义不变）；另加微任务调度开销与热缓存局部性损失。

### 3.4 结论与替代建议

- **结论**：不做。其一，理论上界 ≈0；其二，为 ≈0 收益引入新失败形态与深水轮审查成本，按准则收益标准不成立；其三，真正墙钟大头是 base 四件（P88 报告 §7 建议 2）与 G4b 冷扫（OB-113 子项③），不在 G6 / G7 / G8。
- **替代建议**：维持现状；把 P88 报告 §7 建议 4 登记为「已评估、收益不成立、不实施」（本片闭环）；若未来需要真并行，候选路径 = worker_threads（成本：git 子进程迁入 worker、跨线程缓存一致性、变更集序列化；高风险且上限仅 ≈0.9 s，亦不成立）。

## 4. 影响判定（CONTRACTS/schema/fixtures/internal/**/ci.yml 均不触；候选 A 下不触 §6.3 并发句、不升版）

| 项 | 若实施时是否触及 | 依据 |
|---|---|---|
| CONTRACTS | 否 | 门禁脚本不在该域，`docs/CONTRACTS.md` 无 gate 条款 |
| `schemas/` | 否 | 无 gate 相关定义 |
| `testdata/` 与夹具 | 否 | 不涉及；`selftest.js` 是代码而非夹具 |
| `internal/**` | 否 | 改动仅落判据实现面，历次先例均明示不动 `internal/**` |
| `ci.yml` / workflows | 否 | CI 调用面与退出码契约不变 |

- 规则文本：§6.3 的判据并发调度句与 §12.4 的 OB-113 处置注记，**仅在实施分支**才须随片更新并升 MINOR v1.78；候选 A 下两者均不触，**准则维持 v1.77**。
- ⑦ 档位：实施分支落入深水轮，候选 A 下回归常规档 —— 用户裁定 = `GPT-5.3-Codex (copilot)`，额度 1 终审 + 1 复审。
- 盲审不抽；探针 0；ADR 不新增；不做 gitee。

## 5. 审查链与闭合

### 5.1 ⑤ 预审

（⑤ 输出原文归档于 `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md`。）

⑤ r1 = `PRE-REVIEW: PASS`（0 High / 0 Medium / 2 Low + 3 Info）。两条 Low 已整改：L24 补「低于探针者按 0 计」口径；头部主控行与台账额度口径行补「用户裁定、不得作先例」标记。3 条 Info 为记录态（`~0.9s` 与 `≈0.9 s` 符号形态、9 件载荷计数、P88 口径不可比说明）。无 Medium+，无需复审。

### 5.2 ⑥ 独立扫描

（⑥ 输出原文归档于 `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/06-scan-r1.md`；以下为结构化转录：文本逐字保持，结论行按 §11.3 落于围栏末。）

```sixstage
### Section A
INDEPENDENT SCAN: PASS

### Section B
| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| — | 无可执行发现 | 无可执行发现 | 无可执行发现 |

### Section C
- S1 PASS — 暂存变更集 15 项，全为 `docs/**`；`git diff --cached -- tools/gates` 无输出，无判据或准则实现变更。[中文报告:12-16](docs/validation/gates-criteria-async.md#L12-L16)、[英文报告:12-16](docs/validation/gates-criteria-async_EN.md#L12-L16)。
- S2 PASS — 两份报告保持“不做”、理论上界 ≈0、五项影响全否、候选 A 不升版 v1.77；用户裁定的 ⑦ 常规档 1+1、⑥ 1 次及不得作先例均守恒。[结论证据](docs/validation/gates-criteria-async.md#L10)、[影响证据](docs/validation/gates-criteria-async.md#L56-L67)、[台账口径](docs/t027/REMAINING_SLICES.md#L2539-L2553)、[① 原文](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/01-design.md#L42-L48)。
- S3 PASS — 1247/1177、5799/3812、0.6–0.9 s、≈0.2 s、10–15%/≈5%、<1%、≤250 ms、≤0.3–4%、≈0.9 s、≈94 s 均与 ① 一致；P88 冷基线 20865/12525/10311 与 P88 baseline TSV逐项一致。[报告表](docs/validation/gates-criteria-async.md#L28-L40)、[P88 基线](docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/baseline.tsv#L26-L28)、[⑤ 核验](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md#L19)。
- S4 PASS — 双语均为 UTF-8 BOM、LF 无 CR；LF 数 133、分割行数 134（含末尾空行），行数相等；`artifacts`围栏逐字相等，末行均为闭集令 `SLICE: READY_FOR_REVIEW`。[中文末行](docs/validation/gates-criteria-async.md#L133)、[围栏起点](docs/validation/gates-criteria-async.md#L115-L130)。证据包为 9 载荷，MANIFEST 与 SHA 列表条目及实际载荷集一致；重算 9 个 SHA 均匹配。[MANIFEST](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/MANIFEST.md#L8-L19)、[SHA 列表](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/SHA256SUMS.txt#L1-L9)。
- S5 PASS — 中文/英文净自耗句均已夹钳低于探针者按 0 计；主控头部及台账已标用户裁定“不/不得作先例”。[中文口径](docs/validation/gates-criteria-async.md#L36)、[中文主控](docs/validation/gates-criteria-async.md#L3)、[台账口径](docs/t027/REMAINING_SLICES.md#L2553)、[⑤ 两条 Low 原文](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md#L11-L12)。
- S6 PASS — 未发现禁用占位形态；双语豁免行理由非空，且当前证据显示⑥尚未执行。[中文豁免](docs/validation/gates-criteria-async.md#L131)、[英文豁免](docs/validation/gates-criteria-async_EN.md#L131)、[契约核验](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md#L22)。

### Section D
以下观察会推翻 PASS：暂存集出现 `tools/**`、`internal/**`、CONTRACTS、schema、fixture、CI/workflow；报告翻转“不做”结论或准则 v1.77；任一实测数字漂移；报告对 LF 行数、围栏或末行闭集不一致；MANIFEST/SHA/载荷集或任何载荷 SHA 不一致；低耗钳位或主控不得作先例标记缺失；豁免理由为空、与实际状态矛盾，或出现禁用占位形态。

### Section E
我无法独立复核⑦终审、批量辅助两次实派及⑥记录落盘后的最终内容；逐判据绝对值实测原始捕获件不在本次审查输入中，故其真值主要由①算术分解及⑤原文交叉核验。

INDEPENDENT SCAN: PASS
```

（转录结束；结论行已逐字保留于围栏末。）

### 5.3 ⑦ 终审与复审

⑦ r1 = `FINAL REVIEW: FINDINGS`（1 Medium + 1 Low）。Medium = 报告阶段状态未与已完成环节同步；Low = 「逐字转录」措辞与结论行落位不符。两条已整改（§1 阶段句、§5.1/§6 回填、转录说明改「结构化转录」）。按 §5.3 执行复审 r2（收窄输入：修复 diff + 前轮发现 + 机械读数）。

⑦ r2 = `RE-REVIEW: PASS`（Medium / Low 均 CLOSED，无新增 Medium+）；⑦ 额度 1+1 用满，收敛达成。⑦ r1 / r2 原文归档于 `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/07-final-review-r1.md` 与 `07-final-review-r2.md`。

## 6. 计量

| 角色 | 模型 / 载体 | 次数 |
|---|---|---|
| ① 架构师 | `DeepSeek V4 Pro (deepseek)` Max | 1 |
| ② 文档工程师 | `DeepSeek V4.1 Flash (deepseek)` | 1 |
| ③ 一致性检查 | `DeepSeek V4.1 Flash (deepseek)` | 1 |
| ④ 门禁 | 主控机械执行 | 多轮（index/tree） |
| ⑤ 预审 | `DeepSeek V4 Pro (deepseek)` Max | 1 |
| ⑥ 独立扫描 | `MAI Code 2 Flash (copilot)` | 1 |
| ⑦ 终审 + 复审 | `GPT-5.3-Codex (copilot)`（用户裁定） | 1 + 1（用满） |
| 批量辅助 | `GLM-5.3-Flash (glm)`（载体 `deep-reasoner`） | 2 |

- R = 1、D = 0（终值）；R 口径见 §7；探针 0；盲审不抽。
- 计量单一真值源 = 会话调用回执。

## 7. 异常处理

- **步骤 0 首轮 index 门禁 G8-a 判红**：证据包清单尚缺 `MANIFEST.md` 与 `SHA256SUMS.txt`，判红属真实缺口；补齐后复跑清零。属步骤 0 内自愈，不构成偏离，不单计 R。
- **④ 门禁判红序列（均就地修复复跑至零，按 P88 先例不单计 R，逐条登记）**：G3(b) 状态字形仅 CN 侧（台账「未立项」字样）→ 改写消字形；G7-e 末行 `SLICE: IN_PROGRESS` 不在闭集 → 改 `SLICE: READY_FOR_REVIEW`；G7-f 新报告缺围栏 → 先落非空理由跳过声明，⑥ 完成后以逐字围栏替换声明；G7-f 结论行须在围栏末 → 补于围栏末；G4a 双 BOM（归一脚本未先剥旧 BOM）→ 先剥后补单 BOM；G2 占位残留形态两处 → 改写；G4b 未登字（1 个新字）→ 改写规避（不动 `cjk-newwords.txt`）。
- **⑤ 两条 Low（PASS 判决内）**：已整改（§3.1 口径句、主控裁定标记），无 Medium+、无需复审；按 P88 先例不单计 R。
- **⑦ r1 一条 Medium + 一条 Low**：整改 + 复审 r2 收敛 ⇒ 计 **R = 1**（本片 R 终值）。
- 临时件 D 计数：0——临时件仅 `tmp/P95-GATES-CRITERIA-ASYNC/03-consistency-check.md`（③ 输出，位于声明落点），⑧c 清除。

## 8. 大改建议清单（移交后续片）

1. **P88 报告 §7 建议 4（判据异步化）**：**已评估、收益不成立、不实施**（本片闭环），不再作为候选。
2. **P88 报告 §7 建议 2（跨作用域 base 复用）**：仍为真候选；与「绝不缓存判定结论」同族敏感，须先立条款再实现。
3. **OB-113 子项③（G4b 冷扫内环）**：仍为真候选；须与缓存版本升版及键迁移一并处理，先论证收益。
4. **worker_threads 并行路径**：成本评估不成立（git 子进程迁移、跨线程缓存一致性、变更集序列化；上限仅 ≈0.9 s），不再单独立片。

## 9. 工件围栏

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/00-briefing.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/01-design.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/07-final-review-r2.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/05-batch-assist/glm-mirror-check.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/05-batch-assist/glm-numbers-check.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/MANIFEST.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/SHA256SUMS.txt	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/freedom-list.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/input-list.md	repo
docs/validation/gates-criteria-async.md	repo
docs/validation/gates-criteria-async_EN.md	repo
```

## 10. 提交与 CI

- 主提交 `238e45c`（`docs: record the P95 async-criteria assessment (not implemented)`），已推 `origin/main`（范围 `22370fe..238e45c`；**未推 gitee**）。
- CI run **`45939437180`**（head `238e45c`）终态 `completed / success`；双腿均 success（`Go ubuntu-latest` / `Go windows-latest`），其余候选作业按 CI 定义 skipped。
- 双腿实测硬证据：`SELFTEST: PASS 770/770`、`GATE REPORT scope=ci changed=16`、`TOTAL_FAIL=0`（两腿皆然）；包级 `ok` 行 ubuntu 21 / windows 15、包级 `FAIL` 0；真实失败行核验 = 0（判据自测的预期失败断言描述不计）。
- 本节由 CI 双绿后的 docs-only 提交补入（§11.2）；该回填提交自身的 CI run 只作观察记录。

SLICE: READY_FOR_REVIEW

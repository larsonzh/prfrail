# ProofRail S1-SCOPE-ANNOTATE 报告（P71 · 层 2 治理侧 · 文档类 [SLICE] · M 档）

日期：2026-11-13；状态：⑧ 文档装配定稿（①–⑦ 已闭环）。切片标识：`S1-SCOPE-ANNOTATE`（P71 · 层 2 治理侧 · 文档类 `[SLICE]` · M 档）。

## 1. 执行摘要

- 切片：`S1-SCOPE-ANNOTATE`（P71，层 2 治理侧，文档类 `[SLICE]`，M 档），主控 `DeepSeek V4.1 Flash (deepseek)`；用户 2026-11-13 同轮显式授权。
- 交付：REQ 两条阶段列、MATRIX Windows 行 + 尾注、TEST_STRATEGY AT 两行、s1-exit AT 两行 + 适用性说明、OPERATIONS visible 段、三处加法式指针。
- 结论：X 方案落地（不删条目、不改编号、不改 MVA 复合标准）；S1 范围与阶段归属按 P68 b1 逐条标注；visible 三条边界在位。
- 计量：⑤ 1 次（`PRE-REVIEW: PASS WITH FIXES`，5 Low + 若干 Info，无 Medium+；两条文档层 Low 已整改、三条设计层判据 Low 转为 §4 内容断言，见 §5）、⑥ **跳过**（⑦ 已确认）、⑦ 盲审 1 + 终审 1 + 复审 0（盲审 `BLIND REVIEW: PASS`、终审 `FINAL REVIEW: PASS`）、GLM 批量 2 次、②g 因网络中断重试计 2 次调用。
- 计数：**R = 0、D = 2**（D=1：②a 执行角色在 `tmp/` 根建 `p71-2a-annotate.js` 后自删；D=2：②d 执行角色在 `tmp/` 根建 worktree `p71d-wt` 后清除；两起均零残留，按 §5.1 自愈档登记）。

### 1.1 异常与偏离登记

a) **F-1（起跑消息与上游裁定冲突，用户同轮裁定）**：起跑消息【五】第 4 项原写「REQ-009 / REQ-019 / REQ-007 / REQ-021 四条阶段列都改为 S1 不要求」，与 P68 裁定「不批自托管移出（REQ-019 / AT-13 仍在 S1）」正面冲突。① 阶段发现后停下报实况，用户同轮裁定 F-1 = 选项 A（按 P68 裁定逐条标注：REQ-009 改，REQ-007 子项级改，REQ-021 不动，REQ-019 不动）。本片按 A 执行。这是「起跑消息字面与上游裁定冲突」的实例，教训 = 起草起跑消息时必须逐条核对上游裁定的实际范围。
b) **③ 载体偏离**：③ 一致性检查两次派单均因 `UND_ERR_SOCKET` 网络中断失败，主控改以确定性脚本自执行（可机械复算），登记为执行载体偏离（§13.3.18），独立性影响限于「非他人复核」。
c) **②g 重复调用**：②g 首次派单网络中断但已落盘，重试仅作逐字节核对；按 R7.1 同 stage 重发计入上限，计 2 次调用。

## 2. 差异清单

必改 8 项：

| 文档 | 判定 | 依据 |
|---|---|---|
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | 必改 | REQ-007 阶段列改「S1/S2（黑箱候选子项推迟 S2）」（子项级，不整条移出） |
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | 必改 | REQ-009 阶段列改「S1 不要求，推迟至 S2（基础终端保留）」 |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | 必改 | Windows 行约束列补「基础终端保留 + 完整一体化 TUI 移出 S1、由 S2 承接」 |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | 必改 | 尾部追加 P71 切片同步注记（措辞不含「正式支持范围不变」的歧义） |
| `docs/TEST_STRATEGY.md` / `_EN` | 必改 | AT-13 所有者列补「自托管仍在 S1（P68 b1 不批移出）」 |
| `docs/TEST_STRATEGY.md` / `_EN` | 必改 | AT-24 所有者列补「黑箱子项推迟 S2（P68 b1），条件闭合不撤销」 |
| `docs/validation/s1-exit.md` / `_EN` | 必改 | AT-13 / AT-24 证据列各补阶段归属 |
| `docs/OPERATIONS.md` / `_EN` | 必改 | 新增 visible 调试观察段 + 三条边界 |

需评估 6 项：

| 文档 | 判定 | 依据 |
|---|---|---|
| `docs/validation/remaining-rulings.md` §5 | 适用（加法式 1 行指针） | R-07 b1 附条件 2 由本片落地；原行不改（§11.1） |
| `docs/S1_RELEASE_NOTES.md` 支持段 + 门禁段 | 不适用 | 支持段已含「完整 TUI 不提供」「visible 黑箱候选尚不可运行发布」，现状已一致 |
| `docs/validation/p57-closeout-assessment.md` §4.2 | 适用（加法式 1 行指针） | 本片已回写 s1-exit 适用性说明；八门定义与判定不改 |
| `docs/validation/s1-reach-reachability.md` §2 | 不适用 | 门可达性与 REQ 阶段归属是不同维度 |
| `docs/DEV_PLAN.md` T018 段 | 不适用 | T018 是任务级，REQ 阶段归属是产品级 |
| `docs/validation/s1-stream-progress.md` §8 | 适用（加法式 1 行指针） | P71 已抽盲审（先于终审），该口径已落地 |

## 3. 逐文档回写方案

- `PRODUCT_REQUIREMENTS`：REQ-007 阶段列 = `S1/S2（黑箱候选子项推迟 S2）`；REQ-009 阶段列 = `S1 不要求，推迟至 S2（基础终端保留）`；REQ-019 保持 `S1`、REQ-021 保持 `S1/S2`（F-1 选项 A）。
- `S1_SUPPORT_MATRIX`：Windows 行约束列追加「基础终端（逐行 CLI / `interactions tui` / `--progress`）保留；完整一体化 TUI 移出 S1、由 S2 承接」；尾注追加「P71 切片同步」并明写平台清单与支持层级不变、Windows 行仍为「正式支持候选」、visible 行仍为「降级保证」。
- `TEST_STRATEGY`：AT-13 所有者列追加「阶段归属：自托管仍在 S1（P68 b1 不批移出）」；AT-24 所有者列追加「阶段归属：黑箱子项推迟 S2（P68 b1），条件闭合不撤销」；两行的需求列与断言列零改动。
- `s1-exit`：AT-13 / AT-24 证据列各追加阶段归属；新增「P71 适用性说明」段（REQ 阶段归属变更不改门级判定、门表 9 行行名 / 行数 / 状态全不变、`S-Gx` 编号不引入映射、AT-12 / AT-16 / AT-22 既有 PASS 证据不受影响）。
- `OPERATIONS`：新增「visible 调试观察（便利、非交付）」段 + 三条边界（非 S1 交付承诺 / 不构成权威证据，ADR-012 / 需 VS Code + Copilot Chat + 桥扩展，真实 CLI 端到端仍依赖 T027-B）。
- 三处加法式指针：`p57-closeout-assessment` §4.2、`remaining-rulings` §5、`s1-stream-progress` §8 各追加 1 行指针，原行不改（§11.1）。

## 4. 反例与验证

- X-1（删 REQ 条目或改编号）：判据 `grep -c '^| REQ-'` = 28；③ 实测 `REQ_count=28`，编号 REQ-001–028 连续。
- X-2（visible 行写成「正式支持」）：矩阵表行内含「正式支持」者仅 Windows 行；③ 实测 `matrix_rows_with_zhengshi=1`，命中行 = `Windows 11 amd64`。
- X-3（门表 BLOCKED 改 PASS）：门表逐行比对；③ 实测 `gate_table_identical=true`（2 BLOCKED / 4 PASS / 3 PARTIAL 不变）。
- X-4（REQ-019 阶段列被改）：③ 实测 `REQ-019 stage=[S1]`、`REQ-021 stage=[S1/S2]`，两者均未动。
- X-5（AT-13 / AT-24 断言列被改）：仅阶段列可动；③ 实测 REQ 两条改动行 `changedFieldIdx=[3]`（第 4 列零改动），PC-01 / PC-06 两行 `identical=true`。
- ③ 机械读数汇总：`REQ_count=28`；四条阶段值 = `REQ-007 [S1/S2（黑箱候选子项推迟 S2）]` / `REQ-009 [S1 不要求，推迟至 S2（基础终端保留）]` / `REQ-019 [S1]` / `REQ-021 [S1/S2]`；第 4 列 identical；PC 两行 identical；门表 identical；`S-G3` / `S-G4` / `S-G5` 命中 = 0；visible「正式支持」命中 = 0；全部 CN / EN 对 numstat 对称（G3(a) 九对逐对 added / deleted 相等）。

## 5. 审查链结论

- ⑤ `PRE-REVIEW: PASS WITH FIXES`：5 Low + 若干 Info，无 Medium+。逐条处置如下（按「文档层 / 设计层判据」分类，避免以断言替代整改）：
  - **Low-1（文档层 ⇒ 必须整改）**：原文「`S1_SUPPORT_MATRIX.md:33` 注记自述『正式支持范围本身不变』，而 Windows 行约束列新增『完整一体化 TUI 移出 S1』——『正式支持范围不变』与『S1 功能承诺面缩小』存在措辞张力」⇒ 已整改：注记改为「平台清单与支持层级不变；TUI 移出系 P68 已裁定的功能阶段归属同步」。落点：`docs/S1_SUPPORT_MATRIX{,_EN}.md` 尾注。
  - **Low-2（文档层 ⇒ 必须整改）**：原文「`s1-exit.md:62`『完整 TUI 移出不构成门级判定变更』是断言性声明，缺论证：AT-12/16/22 断言全部基于基础终端 / `interactions tui`、不依赖完整一体化 TUI，故既有 PASS 不被削弱」⇒ 已整改：说明段补「AT-12 / AT-16 / AT-22 断言不依赖完整一体化 TUI，既有 PASS 证据不受影响」。落点：`docs/validation/s1-exit{,_EN}.md` 适用性说明段。
  - **Low-3（设计层判据的可执行性缺陷 ⇒ 转内容断言）**：原文「判据 `grep -c '^\| REQ-'` 在 BRE 下 `\|` 为 alternation，`^` 命中每行 ⇒ 字面执行返回文件总行数而非 28，判据写法错误」⇒ 转为：判据写作 `grep -c '^| REQ-'` = 28，并附编号 REQ-001–028 连续性。落点：本报告 §4 **X-1**。
  - **Low-4（设计层判据的可执行性缺陷 ⇒ 转内容断言）**：原文「『正式支持』全文件 grep 会额外命中 P71 注记行的『正式支持范围』字样，『仅命中支持候选行』字面不满足」⇒ 转为：限定为表格行内断言「矩阵表行内含『正式支持』者仅 Windows 行（状态列为『正式支持候选』的唯一点），visible 行不得含该短语」。落点：本报告 §4 **X-2**。
  - **Low-5（设计层判据的可执行性缺陷 ⇒ 转内容断言）**：原文「X-3 / X-4 的行号锚点（`:52` / `:60` / `:81`）当前正确，但文件前部任一行增删即漂移，机械判定脆弱」⇒ 转为：内容断言「门表逐行比对 identical + 状态词计数 2 BLOCKED / 4 PASS / 3 PARTIAL」「REQ-019 / REQ-021 取第三列取值」。落点：本报告 §4 **X-3**、**X-4**。
  - **分类判定**：Low-1 / Low-2 是**文档层**的措辞张力与论证缺口 ⇒ 整改被审文档；Low-3 / Low-4 / Low-5 是**设计层判据本身**的写法错误 / 口径过宽 / 锚点脆弱，不涉及任何被审文档的正确性 ⇒ 改为 §4 的可复算内容断言，**不构成「用断言替代整改」**。
- ⑥ **跳过**：纯 `.md` 回写，无 `.go` / `.json` / `.yml` / schema / fixtures / 脚本逻辑语义变更，不触状态机 / 门禁判据 / 证据模型 / 角色授权本体 ⇒ §7.2.1 三条件成立；⑦ 终审已确认。
- ⑦ 盲审 `BLIND REVIEW: PASS`（不附 ⑤ / ⑥ 清单，先于终审）；⑦ 终审 `FINAL REVIEW: PASS`（0 发现）。
- 盲审与终审结论一致，无分歧；⑦ 复审 0 次（终审 0 发现，无整改环）。

## 6. 门禁读数

- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`。
- `--scope=ci` ⇒ `GATE REPORT scope=ci changed=4`、`TOTAL_FAIL=0`。
- ⑨ 定稿读数（精确暂存 31 路径后，证据包整目录、未用 `-A`）：`--scope=index` `changed=31` `TOTAL_FAIL=0`；`--scope=tree` `changed=31` `TOTAL_FAIL=0`；`--scope=ci` `changed=4` `TOTAL_FAIL=0`；`--selftest` `SELFTEST: PASS 737/737`。
- 暂存前 tree 的唯一红 `G7-b`（7 条 `repo` 工件仍为未跟踪：`02-batch/` 2 件 + `04-reviews/` 3 件 + 报告对自身 2 件）已随精确 `git add` 转绿。
- 报告对生成前的 index / tree 读数：`G5(d)` 红（结构性：staged 路径尚无人声明；报告围栏生成后消解）、`G5(f)`（未暂存，已按纪律 E 处理），暂存后均转绿；另见 `G8-b` 红（证据包「归档条目清单」与 `SHA256SUMS.txt` 未含 `02-batch/` 与 `04-reviews/` 五件载荷），已在本片 ⑧ 收口时补齐后转绿。
- 基础门禁：`gofmt -l .` 空输出；`go build ./...` exit 0；`go vet ./...` exit 0；`go test ./...` ok=15 / FAIL=0。
- CI 链路（提交后）：主提交 `6ee2c61` 的 run `43463089586` **双腿 success**（windows `5m6s`、ubuntu `3m14s`；两腿包级 `ok github.com/larsonzh/prfrail/internal/console` 分别为 `21.982s` 与 `4.160s`（Test）/ `13.624s`（Race）；两腿 `--- FAIL` 行数均为 0、包级 `FAIL` 行数均为 0；包级 `ok` 计数 windows 腿 15 条、ubuntu 腿 21 条）。`Candidate build` / `Candidate probe` / `Bootstrap and release evidence` 三个 job 按既有配置 skipped，非本片引入。本 CI 事实链由回填提交承载（回填提交 hash 与回填 run 见该提交记录与 ⑩ r2 回执，避免自指）；回填只补 CI 读数，不改台账（§11.2：台账完成行只写报告路径，hash / run 只进报告）。

## 7. 已知边界

- ① `s1-exit` 门表为 9 行中文门名、无 `S-Gx` 行标识；`S-Gx` 定义在 `docs/validation/p57-closeout-assessment.md` §4.2，本片不建立映射。
- ② `S-G2` 无独立行，且其口径为「门级判定待 v54 快照」，本片不改。
- ③ **REQ-020 的「S2 承诺收窄登记」不在本片范围**（⑤ Info 建议登记，留待后续切片）。
- ④ `S1_SUPPORT_MATRIX` 无 TUI 行，故 TUI 移出以 Windows 行约束 + 尾注表达，不新增行。

## 8. 明确未执行

- 未改准则正文 / `CONTRACTS` / `SECURITY` / `ARCHITECTURE` / `ADR_REGISTER` / v50–v56 快照 / REQ 表结构 / MVA 复合标准 / S1 Exit 八门判定本身 / 产品代码。
- 未 commit / push / 推 gitee；未启动新切片。
- **D=2 防再犯句（用户 2026-11-14 裁定）**：执行角色（`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`）在 `tmp/` 下操作时，一律使用 `tmp/<切片ID>/`，不得在 `tmp/` 根创建任何文件、目录或 `git worktree`。本片 D=2 两起均为该违规的实例（②a 的 `p71-2a-annotate.js`、②d 的 `p71d-wt` worktree）；若第 3 次发生（片内 D=3）按 §5.1 累计上界强制停片。**P72 起跑纪律须引用本句**。

- ⑥ 跳过：纯 .md 回写（无 .go/.json/.yml/schema/fixtures/脚本逻辑语义变更），§7.2.1 三条件成立，⑦ 终审已确认。

```artifacts
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/00-briefing.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/01-arch/design-r2.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/02-batch/batch-01-req-enumeration.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/02-batch/batch-02-at-gates-enumeration.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/04-reviews/blind-r1.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/04-reviews/final-r1.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/MANIFEST.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/SHA256SUMS.txt	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/freedom-list.md	repo
docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/input-list.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/remaining-rulings.md	repo
docs/validation/remaining-rulings_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/s1-scope-annotate.md	repo
docs/validation/s1-scope-annotate_EN.md	repo
docs/validation/s1-stream-progress.md	repo
docs/validation/s1-stream-progress_EN.md	repo
```

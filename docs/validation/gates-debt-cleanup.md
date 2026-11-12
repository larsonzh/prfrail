# 切片验证报告 · P67 GATES-DEBT-CLEANUP

- 准则版本：v1.66（步骤 0 时点；步骤 F 将升至 v1.67）
- 日期：2026-11-12
- 分类：治理类 `[SLICE]`；规模 L 档；探针额度 0（Windows 本机治理片）
- 状态：🔄 **进行中**（2026-11-12 起跑；步骤 0 已登记；停点 ⑩ 等授权）
- 主控模型：`DeepSeek V4 Pro (deepseek)`（ADR-041）
- 证据包：`docs/validation/evidence/P67-GATES-DEBT-CLEANUP/`

## 1. 执行摘要

本片为层 2 治理侧一次性清欠片（#1 g4b 缓存 LRU + 上限 / #2 门禁坑 27–30 正式化 / #3 §12.4 抽查记录块清账 / #4 ⑥ 提示词加固）。进度（2026-11-12）：步骤 0 完成；① 五件落盘（`DESIGN: DONE`）；② A–G 全部完成（每步核验件见证据包 `03-test/`）；③ 一致性 `PASS`；④ 三作用域 `TOTAL_FAIL=0` + `SELFTEST: PASS 736/736`；⑤ r1 `PASS WITH FIXES`（6 Low）→ 整改 → r2 `PASS WITH FIXES`（1 Low，额度 2/2 用尽，字面修正已闭环并交 ⑦ 确认）；⑥ r1 `INDEPENDENT SCAN: PASS`；⑦ 盲审 r1 `FINDINGS`（1 High + 1 Medium，已整改）→ ⑦ 终审 r1 `FINDINGS`（1 Medium 报告装配 + 1 Low `CACHE_CAP` 钉值）→ T130 已落地（`SELFTEST: PASS 737/737`）→ 本 ⑧a 为 Medium 整改 → ⑦ 复审 `RE-REVIEW: PASS`（⑧b 时点）；⑨ `SELFTEST: PASS 737/737` + 三作用域 `TOTAL_FAIL=0`；批量辅助 GLM-5.3-Flash 3 次（① 缓存键枚举 / ② 七坑对照 / ③ ⑥ 三轮枚举，归档 `02-batch/`）。

- ⑥ 独立扫描（实跑）：r1 = `INDEPENDENT SCAN: PASS`（A–E 合规、0 发现；原文归档证据包 `04-reviews/06-scan-r1.md`）；逐字转录见文末 `sixstage` 围栏（§11.3），围栏外注记紧随其后。

- **异常登记区（用户 2026-11-12 同轮裁定）**：
  1. **① 越权判定（坑 6）**：① 阶段识别到 3.1.3 与起跑消息【七】5「1–26 条零改动」的矛盾，未按 3.4 停升级边界上报，自主选择 3.1.3 优先。处置：用户裁定选首选①——坑 6 完全退出 #2 正式化范围；第 21 条逐字不动；坑 6 仅作报告观察项（§8 已知边界 + v55 事实包）；「3.1.3 关于坑 6 的处理条款（归入第 21 条处置示例补充）因起跑消息【七】5『1–26 条零改动』的硬约束而整体作废；本条属用户 2026-11-12 同轮裁定」。
  2. **① 越权判定（mtime 边界）**：① 阶段做了「mtime 是否算时间戳」的语义评估但未作为 3.4 停升级边界上报，采宽松读法。处置：用户 2026-11-12 裁定宽松读法成立、mtime in-bounds（理由 = 3.1.2 字面语境解读；潜在保守解读的存在已登记）；本条属边界解释澄清、不属 §1.7.4 临时授权超范围、不单独立 ADR；不得作为先例。
  3. **R 计数口径**：以上两处偏离不构成修复重跑（R 不 +1），作为「用户发现」登记于本区（对应 P58 §12.1 试点指标的「用户发现」计数口径）；本片「用户发现」= 2。
  4. **临时件 D 计数 = 1（主控本人违规，自愈档成立——用户 2026-11-12 同轮直接确认；本条属对 §5.1「违规者为主控本人时自愈判定须经同轮独立审查层确认」字面要求的临时突破，因步骤 B 时点 ⑤ 未启动、无法在同轮取得独立审查层确认；§5.1 其余条款（自报 + 预快照 + 自删 + 零残留证据）按纪律严格执行，无突破；本例外仅本片有效，不得作先例）**：事件事实：stray `8`（425 字节）出现在仓库根（§5.1 位置闭集之外），成因 = 主控嵌套引号损坏的 `cmd /c "powershell …"` 尾部读取命令；三处预快照（porcelain `?? 8` / tmp 递归 115 / TEMP 命中 0）→ 精确 `Remove-Item -LiteralPath .\8` → 零残留（`Test-Path`=False、`??` 为空、tmp 递归仍 115）。
  5. **#3 落点偏差登记（⑤ F3 整改）**：design-r1 §4 的②「补录可追溯字段」实现为「各片 ⑦ 盲审事实以各片报告『审查与扫描记录』为唯一真值源、本块不转抄」（历史欠账、§11.1 不追溯补录）；弃 C 部分已在 design-r1 §4 如实标注，本区补登记（理由 = §7.7:798 正文句与 F-8 双约束，见 design-r1 §4）。
  6. **执行纪律偏差登记（用户 2026-11-12 ⑩ 补正指令）**：主控在 ② 阶段跳过 A/D/E/F 步核验件落盘（违反【六】2「每步独立核验 + 落 `03-test/step-X-check.md`」），经 ⑦ 盲审 1 High 发现后补齐（step-A/D/E/F-check.md 已落盘并入 MANIFEST/SHA256SUMS）；属主控实现过程执行纪律偏差、非作品缺陷。
  7. **R 计数（§9.2 口径，用户 2026-11-12 ⑩ 补正核对）**：R = 3——⑤ r1（6 Low 发现）→ 整改 → ⑤ r2 = 1；⑦ 盲审 FINDINGS → 整改 = 1；⑦ 终审 FINDINGS → 整改 → 复审 = 1。R ≥ 2 ⇒ 按 §5.2 高风险行执行（⑤ 复查为必：已执行且 0 Medium+；⑨ 全绿）；R = 3 未达 §9.2「R > 3 暂停等人工审核」阈值。
  8. **终端假中断观察（§9.3 已知模式，用户 2026-11-12 ⑩ 补正指令）**：④ 复跑 tree 作用域时终端显示一次 `^C` 误报（进程实际已完成、读数终态 `TOTAL_FAIL=0`，见 `06-gates/gates-r1.txt`）；本片未受其影响，登记在案。
  9. **坑 7 §12.4 落点句为 ⑩ 补正后新增，未经 ⑦ 复审（用户 2026-11-12 ⑩ 裁定）**：性质 = 记录态落点（v54 §一 已执行该口径 + 指向 §13.5 已生效模板 + 元声明），非规则语义新增；不触 §7.2.1 ③；按 §5.3 精神「修复后重跑」针对规则语义修复，记录态补充不需重跑 ⑦。本例外仅本片有效，不得作先例。

## 2. #1 四案对比与最终落点

| 维度 | 案 A：条目上限到 N 拒写 | 案 B：LRU 淘汰 | 案 C：TTL + 上限 | 案 D：LRU + TTL |
|---|---|---|---|---|
| 实现复杂度 | 低 | 中 | 中 | 高 |
| 失效语义正确性 | 正确但**缓存停滞**（新 sha 永久 miss） | 正确（淘汰 = 失去命中机会，永不错误命中） | 正确 | 正确 |
| 空间上界 | 硬 N | 硬 N | 软（依赖过期检查时机） | 硬 N + TTL 兜底 |
| 对既有条目 | N ≥ 96 时不清理、此后拒写 | put 后逐条淘汰最旧 mtime | 到期逐步清 | 同案 B |
| 风险 | 中：缓存退化为死缓存 | 低 | 中：需时间源与时序语义 | 低 |
| 推荐度 | 低 | **推荐** | 中 | 不必要 |

**结论**：最终落点 = 案 B（LRU，N=100，mtime 访问序）。旧 sha 条目在 `resolvedSha` 前进后永不再命中，故按 mtime 淘汰最旧即淘汰纯占位条目，案 C/D 的 TTL 属冗余。mtime 边界（是否算「时间戳」）经用户 2026-11-12 裁定 in-bounds（详见 §1 异常登记区第 2 条与 §7），理由 = 3.1.2 字面语境解读。来源：`01-arch/design-r1.md` §1、§2。

## 3. #1 判据本体改动（g4b 缓存 LRU + 上限）

- 改动面：`tools/gates/lib/criteria/g4b.js` numstat `93/2`，为本片唯一判据代码改动；新增模块级常量 `CACHE_CAP = 100`、条目名正则 `CACHE_ENTRY_RE`、纯函数 `evictToCap(dir, cap, statFn, unlinkFn)`，均经 `internals` 导出供 selftest 读取。来源：`02-impl/impl-r1.md` §1、`03-test/step-A-check.md` §1。
- 访问序来源：`diskStore.get` 在 `JSON.parse` 成功后以 `fs.utimesSync` touch 条目；`diskStore.put` 在写入与 rename 成功后调用 `evictToCap(dir, CACHE_CAP)`。来源：同上 §2、§3。
- 淘汰规则：只计匹配 `^[0-9a-f]{64}\.json$` 的条目，按 `mtimeMs` 升序逐个 unlink，删除数 = `entries.length − cap`；stat 或 readdir 失败返回 0，绝不抛异常。来源：同上 §2、§4。
- 六边界保持：加 LRU + 上限不触判据主流程；键五分量 diff 为零；值结构仍 `{v: 2, corpus: string[]}`（mtime 为文件系统元数据、不进 JSON 值，`validValue` 不变）；命中 = 与冷跑等价；fail-open 全路径静默；只存语料不存结论。来源：`01-arch/design-r1.md` §2、`03-test/step-A-check.md` §3。
- 测试增量：`tools/gates/selftest.js` numstat `363/0`（纯插入，既有 T1–T124 与 F1 逐字未动）；T125–T129 覆盖淘汰序与条目名过滤、上限三分支、真实 put 淘汰、命中等价、fail-open 四路；⑦ r1 Low 整改追加 T130（把 `CACHE_CAP` 钉在字面值 100，防常量漂移）。来源：`03-test/step-B-check.md` §2、`tools/gates/selftest.js`。
- SELFTEST：`722/722`（步骤 0，`06-gates/step0-selftest-r1.txt`）→ `736/736`（步骤 B 时点，`06-gates/gates-r1.txt`）→ `737/737`（T130 落地后，⑧a 复跑）。T125–T130 净增 +15。

## 4. #2 门禁坑正式化（编号 27–30 续接）

- 七坑判定：4 条正式化为第 27–30 条；坑 6（GBK 误码）按用户 2026-11-12 裁定完全退出 #2 范围，仅作报告观察项；坑 7（§一 滚动窗口歧义）不新增准则条目，转 v55 事实包登记。来源：`01-arch/design-r1.md` §3、`03-test/step-C-check.md`。
- 四条落点：第 27 条未登记 CJK 字符触发 G4b（三选一处置 + §6.3 判据作用域排除口径）／第 28 条「共 N 条」类聚合计数措辞触发 G6-4 聚合计数形态／第 29 条「候选」类时点措辞迁入正式表后被误读／第 30 条行尾截断（CN/EN 截断不对称实证）。来源：同上。
- 续接与零改动：新条目自第 27 起连续编号；既有第 1–26 条语义零改动，`03-test/step-C-check.md` §3 以 `git diff -U0` 单 hunk 纯插入证明第 1–26 条、表头与表后登记册口径块逐字未动；第 21 条未补任何处置示例。来源：`03-test/step-C-check.md` §1、§3。
- 双侧对称：`docs/DELIVERY_DIRECTIVE.md` 与 `_EN` 各 `4/0`，两侧行数 2391 → 2395 相等，编号 1–30 连续、无断号无重号。来源：同上 §1、§2。

## 5. #3 §12.4 抽查账清账（A′ 方案）

- 事实冲突：窗口起点短语「窗口自 v1.32 起计」在准则双处并存（§7.7 OB-82 子款与 §12.4 表尾块）；仅改块会与正文句矛盾（违 §11.2「同一事实只在一处描述」），同步改正文撞 freedom-list F-8「不改准则 §1–§12 已有条款」。来源：`01-arch/design-r1.md` §4。
- A′ 落点：不动窗口起点（两处短语逐字保留）；§12.4 块内更新时点句为事实句，补录可追溯字段（是否抽中历史不可追溯 ⇒ 如实标注历史欠账），声明 v1.32–v1.66 期间「启动前判定并登记」义务未履行为历史欠账（§11.1 不追溯补录），并自 P67 起恢复逐片登记。来源：同上；`03-test/step-D-check.md` §2。
- 偏差登记：design-r1 §4 的②「补录可追溯字段」实现为「各片事实以各片报告审查链为唯一真值源、本块不转抄」；起跑消息推荐 A+C 的 C 部分（窗口重置）弃用，理由 = 上述双约束；两项均登记于 §1 异常登记区第 5 条。来源：同上；`03-test/step-D-check.md` §4。
- 边界核验：§7.7 段 diff 为零，块内与 §7.7 无矛盾窗口起点表述；numstat 两侧 `3/1`，行数 2394 → 2396 相等。来源：`03-test/step-D-check.md` §1、§3。

## 6. #4 ⑥ 提示词加固（附录 B.6）

- 五个加固点：① A–E 五段标题固定为 `## Section A`–`## Section E`；② 结论行位置与取值（恰两处 = Section A 首行 + 块末，取值全等）；③ Section B 的 `SEVERITY` 行与结论行一致；④ `UNCOVERED` 独立成句；⑤ 陈旧读数识别（⑥ 须显式声明输入包版本与时间戳，不得引用过时读数）。来源：`01-arch/design-r1.md` §5。
- 只做文本对齐：①②③ 落在 §7.2.2 与 G7-f 已定义口径内，不新增 G 判据；④ 为既有表述强化；⑤ 以输入面指令形态写入，未写入 §7.2.2。来源：同上。
- 仅限 ⑥ 注记：B.6 注明 A–E 形状仅限 ⑥，⑦ 仍按 §4.1 闭集结论行；B.7.4 不改。来源：同上。
- 零改动核验：§7.2.2、§6.3、§4.1 与附录 B.7 的 diff 均为 0；numstat 两侧 `5/1`，行数 2396 → 2400 相等；⑤ r1 F5 整改后标题措辞与 G7-f 字面一致（CN/EN 各 `+1/−1`）。来源：`03-test/step-E-check.md` §2、§3。

## 7. 已知边界

- ① 阶段的 3.4 停升级边界执行存在两处偏离，经用户 2026-11-12 同轮裁定处置（详见 §1 异常登记区）。
- mtime 边界裁定（用户 2026-11-12）：3.1.2 的「若需加时间戳/淘汰标记」真实意图 = 「在缓存值结构内新增时间戳/淘汰标记字段」；mtime 属「利用既有文件系统元数据」，不在禁止范围；六条边界全部守住。**不得作为先例**：未来同类边界仍需独立裁定。
- 核心纪律（用户 2026-11-12 裁定）：后续遇 3.4 停升级边界时，①/② 必须按边界字面停下上报，不得自主处置。
- #1 判据本体改动限于 LRU + 上限（案 B、N=100、mtime 访问序）；若后续再现 G4b 缓存类问题，可能需要更深层改动（如改缓存键语义）。
- 步骤 C 正式化范围：仅追加第 27–30 条；坑 6 按用户 2026-11-12 裁定退出（仅作观察项）；现有 1–26 条严格零改动。
- §5.1 时序限制（用户 2026-11-12 裁定）：本片步骤 B 时点无法取得同轮独立审查层确认（⑤ 未启动），故本片主控本人违规事件由用户代行确认；此为时序限制，非规则缺陷；后续同类情形仍按 §5.1 字面执行（同轮审查层确认），不得援引本片为例。
- 坑 7（§一 滚动窗口歧义）：⑤ F2 时点曾登记「转 v55 事实包」；经用户 2026-11-12 ⑩ 补正改落 §12.4 表尾清账块（④ 句，记录态）；v55 事实包按 §13.3.15 记一笔、口径已定。
- **坑 7 落点句为 ⑩ 补正后新增**：属记录态落点，不经 ⑦ 复审；用户 2026-11-12 同轮裁定；性质与依据见 §1 第 9 条。
- ⑤ r1 时点读数口径（F6 整改）：⑤ r1 时点报告对行数为 53、`06-gates/gates-r1.txt` 为步骤 F 前中态读数；终态以 ⑨ 实测为准，MANIFEST / SHA256SUMS 与报告读数于收档时重算。

## 8. 明确未执行事项

- 未改判据本体：唯一例外 = #1 的 LRU 加严（§3）。
- 未改 `CONTRACTS`、schema、fixtures 与产品代码（`internal/`、`cmd/`）；未改 CI 配置。
- 未改准则 §1–§12 已有条款、§13.1–§13.3 与 §13.5 正文；未改 §7.2.2、§7.7 正文；未改 v50–v54 快照。
- 未改门禁坑第 1–26 条语义（坑 6 仅作观察项，第 21 条未补处置示例）。
- 未替用户裁定：本片两处边界处置均由用户 2026-11-12 同轮裁定（§1 异常登记区）。
- 未 commit、未 push、未 `git add`；未启动新切片；未开 VM、未开代理、未起 Sol 双打。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/00-briefing.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-batch/glm-batch1.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-batch/glm-batch2.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-batch/glm-batch3.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/consistency-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-A-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-C-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-D-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-E-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-F-check.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/05-prereview-r2.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/07-blind-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/04-reviews/07-rereview-r1.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/06-gates/gates-r1.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/06-gates/gates-r2.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/06-gates/step0-selftest-r1.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/MANIFEST.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/SHA256SUMS.txt	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/freedom-list.md	repo
docs/validation/evidence/P67-GATES-DEBT-CLEANUP/input-list.md	repo
docs/validation/gates-debt-cleanup.md	repo
docs/validation/gates-debt-cleanup_EN.md	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/lib/criteria/g4b.js	repo
tools/gates/selftest.js	repo
```

```sixstage
## Section A

INDEPENDENT SCAN: PASS

Input package = P67 GATES-DEBT-CLEANUP, 2026-11-12; this review only reads the materialized evidence and the directive wording in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L720-L760), and the P67 proof bundle under [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md), [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md), and [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md). No file creation, modification, deletion, or command execution occurred.

## Section B

(no findings)

## Section C

- The P67 entry explicitly records the in-scope change set and the boundary that matters for this review: “不改 §7.2.2 / §6.3 / §7.7 正文; 不改 §13.1–§13.3 / §13.5; 不改 v50–v54.” See [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115). That is the key fact for the B.6 hardening question: the change is limited to prompt framing and evidence packaging, not a new rule body.
- The LRU design is judged in-bounds under the six boundaries in [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/01-arch/design-r1.md). The critical condition is the user ruling recorded in that file: mtime is used as filesystem metadata for access ordering, not as a stored value, while the value shape stays exactly as the previous corpus record.
- The implementation memo confirms the actual behavior in [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/02-impl/impl-r1.md): successful parse triggers a touch on get; successful write triggers eviction on put; eviction only examines matching cache files and uses mtime order; all failure paths remain fail-open. That matches the cache contract in [tools/gates/lib/criteria/g4b.js](tools/gates/lib/criteria/g4b.js).
- The selftest evidence in [docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md](docs/validation/evidence/P67-GATES-DEBT-CLEANUP/03-test/step-B-check.md) records T125–T129 and a green SELFTEST: PASS 736/736 result. This directly covers the concerns around bound enforcement, hit equivalence, and fail-open handling. The same file also states that the change did not touch the rule body in [tools/gates/selftest.js](tools/gates/selftest.js).
- The “7 pit merge” is in-bounds as described in the P67 design memo: the new items are a debt cleanup of four specific edge cases, and the file expressly preserves the no-change set for the earlier items. The “坑 6 退出” and “坑 7 转 v55” branches are treated as local routing decisions, not as semantic expansion of the core gate logic. This is consistent with the directive’s “not changing the rule body” boundary in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115).

## Section D

- If the eviction guard is removed, the cap and oldest-first slice checks in T125–T127 would go red; the design memo explicitly marks those as the failure points.
- If the touch-on-hit guard is removed, the LRU ordering assertions T128a–T128c fail, because hit semantics would no longer remain equivalent to a cold run.
- If the fail-open catches are removed, T129a–T129d fail; the implementation memo states the failure paths are intentionally silent and must not alter verdict semantics.
- If the B.6 prompt hardening is treated as creating new semantic rules instead of a documentation-only surface, the directive’s own no-change statement in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L78-L115) would be violated. The evidence package keeps it within the “prompt only, rule body unchanged” envelope.
- UNCOVERED: the materialized evidence does not prove cross-platform NTFS/exFAT timestamp resolution or every rare Windows file-handle race under all OS conditions; that is an edge-case gap, not a contradiction with the current PASS.

## Section E

- Cannot verify from the evidence bundle alone: real runtime behavior across all Windows filesystems and multi-process races beyond the documented fail-open cases. The input package is static and intentionally read-only; the directive requires that this scan only read the materialized evidence, not execute commands or modify files.
- What would overturn PASS:
  - any evidence that mtime is folded into the stored value semantics instead of remaining metadata;
  - any CN/EN asymmetry in the B.6 wording or a change to the rule body in §7.2.2, §6.3, or §7.7;
  - any materialized test result showing the selftest or gate assertions are no longer green;
  - any evidence that the “1–26 unchanged” claim is false for the actual authored patch.

INDEPENDENT SCAN: PASS
```

以上 A–E 五段为 ⑥ r1 的逐字转录（§11.3）；原文见证据包 `04-reviews/06-scan-r1.md`。

# P54 `CONTRACTS-S7-DURABILITY` 验证报告（步骤 0 骨架）

切片：Windows 耐久契约升级与 C6 落地（L 档 · 契约与实现类 `[SLICE]`）。锚点 HEAD = `2b8ac3e`，准则 v1.60。R-02 裁定 2026-11-04 已批准立项。

> 状态：步骤 0 骨架；① 架构（双打）→ ③ 红测 → ② 实现 → ④ 门禁 → ⑤⑥⑦ 审查 → ⑧ 文档 → ⑨ 验证 → ⑧c → ⑩ 停点按轮追加。

## 1. 异常登记

- ① 起跑消息目标矛盾（结构缺陷）——位置：起跑消息【六】.1 的 MoveFileEx 示例；矛盾点：标题与目标句「落实 C6」与示例机制（B3b C2 MoveFileEx）分属两套候选编号；裁定依据：用户 2026-11-04 裁定「起跑消息正文的 MoveFileEx 示例视为『怎么写』的举例，不改变『落实 C6』的目标；标题与目标句优先」；处置：本片目标机制 = ADR-013 的 C6（预建固定槽），本表 C6/C2 分列口径维持。
- ② G7-f ⑥ 跳过声明（原下一条为 G7-f 固定形式行）——状态：已由 ⑧a 的 `sixstage` 围栏替换（⑥ 原文逐字转录于 §6，块末为闭集结论行），过渡占位声明已删除；⑥ 实得额度 = 3；本片修订 CONTRACTS §7 发布耐久等级语义，§7.2.1 ③ 不成立，⑥ 必做。
- ③ 本临时停点（① 完成 → 设计审阅 → 放行红队）——位置：① 完成之后、Sol 红队之前；裁定人：用户；依据：用户 2026-11-04 裁定「本停点由用户在本片显式设定，只生效一次；红队放行后不再重复；不替代 ⑩，不改变起跑消息第七节的常规节奏」；只生效一次、不推广；裁定原文逐字归档于 `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/ruling-2026-11-04.md`。
- ④ 临时授权登记缺口（按 OB-134 先例如实披露）——实例：P53 ① 双打 Sol、本片第一次 Sol 红队；是否登记：未在 §12.4/ADR 独立登记（独立留痕不成立）；缺口性质：§1.7.4 ② 要求 OB 或 ADR 独立登记；处置：如实披露、不追溯补登。用户 2026-11-04 裁定计数按真值源口径（本类仅 OB-132 一条），本片第二次红队计第 2 次、可授权。
- ⑤ 治理侧输入（2026-11-04 用户建议留痕）：本类临时授权真值源口径在 P54 达到 3 次，触发 §1.7.4 ⑤ (iii)；后续片不得再走临时授权，须立项修订片把白名单外红队模型纳入 §2.2 或独立条款。
- ⑥ 落法偏离登记（2026-11-04）：红队闭合自检表落报告 §5（独立成节），偏离此前「入报告 §3」的口径；理由 = 独立成节更清晰，不影响合规。
- ⑦ tmp 根误放件自愈（§5.1 自愈档，2026-11-04）：主控在 ⑩ 终态回读发现 `tmp\s7a-contrasts.js`（tmp 根误放，疑似 ②e 执行角色对比脚本残留；gitignored 故 `git status` 不可见）；次序：自报 → 记录证据 ② 读数（253）→ `Remove-Item -LiteralPath` 精确自删 → 零残留三项（路径 `Test-Path`=False；`??` 清单空；递归文件数 252 ≤ 253；TEMP 命中 0）；D（§5.1 临时件纪律口径）= 1；因违规者归属存疑（主控或 ②e），自愈判定按「主控本人」从严交同轮独立审查层确认——独立审查层（MAI-Code-1.1-Flash）已 `CONFIRMED`（原文见 evidence 包 `07-final/review.md` 第三条）。

## 2. artifacts 围栏

以下为本片完整变更集条目（证据载荷、报告对、台账、规范文档与代码/测试改动；路径 + 归属，TAB 分隔，按路径字典序）。步骤 0 只声明当前已落盘条目，后续环节按轮追加。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/CHECKLIST_v43.md	repo
docs/t027/CHECKLIST_v43_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/contracts-s7-durability.md	repo
docs/validation/contracts-s7-durability_EN.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/00-briefing.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2b.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2c.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review-r2.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/ruling-2026-11-04.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/ruling-2026-11-04-q1.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/05-prereview/review-r1.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/06-scan/review.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/07-final/review.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/cross-platform-linux.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/gate-readings.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/green-c6slot-tests.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/red-baseline-c6slot-compile.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/MANIFEST.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/SHA256SUMS.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/freedom-list.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/input-list.md	repo
internal/adapters/agent_runner_replay_store.go	repo
internal/adapters/agent_runner_replay_store_c6_enable.go	repo
internal/adapters/agent_runner_replay_store_c6_noop_test.go	repo
internal/adapters/agent_runner_replay_store_c6_windows.go	repo
internal/adapters/agent_runner_replay_store_c6_windows_create.go	repo
internal/adapters/agent_runner_replay_store_c6_windows_test.go	repo
internal/adapters/agent_runner_replay_store_launches.go	repo
internal/adapters/agent_runner_replay_store_terminals.go	repo
internal/console/runtime_agent.go	repo
```

## 3. 步骤 0 记录

- 锚点核对：HEAD = `2b8ac3e`、origin/main = `2b8ac3e`、准则 v1.60、`SELFTEST 697/697`、三作用域 `TOTAL_FAIL=0`，5 项全部相符（代差免责声明已核对：回填提交使 HEAD 前进）。
- v43 §十 第一句「v42」→「v43」已改（CN+EN 镜像同步）。
- 台账已登记（状态表 `🔄 进行中` + 小节，R/D=0 基线）。
- 前置核对：CONTRACTS §7 耐久三态口径、B3b 结论账本、A7 C6 定义与 D4 契约冲突已阅；关键开放问题（ADR-013 的 C6 编号与 B3b C2 `MoveFileEx` 的关系）登记于 `input-list.md`，交 ① 裁决。

## 4. C6 升级判定表（步骤 0 停点 → 已裁）

判定时点 2026-11-04；依据行全部给出 file:line，不得推断。**整体初判 = 需用户裁决；用户裁定（2026-11-04）= Q1 甲（本片目标机制 = ADR-013 的 C6 预建固定槽）/ Q2 子（Windows 平台发布耐久等级保持 `unproven`，first-dispatch 继续拒绝）**；P54 契约修订中 Windows 平台等级不得写 `proven`；C2 的候选级结论已在 B3b 结论账本，不需 P54 重写；裁定原文见 §1 ③ 与归档件。

| 行 | 字段 | 判定 |
|---|---|---|
| C6（本片命名，ADR-013 编号） | 机制描述 | bootstrap 预建固定槽，发布=写槽内数据 + 文件级 flush，把目录项耐久转成文件数据耐久（有标准保证）；现行契约下不可即插（与「每 requestId 独立文件 + no-replace」冲突，仲裁需换锁原语）。出处：`docs/validation/t027-windows-publication-durability.md:28`、`docs/ADR_REGISTER.md:23` |
| C6 | B3b 注入矩阵覆盖 | 未覆盖。依据行：`docs/validation/t027-b3b-durability.md:44-48`（候选表仅 C1/C2/C3/C2C3/positive-control）；`:79`（主矩阵 17 格×5=85 轮，无 c6 格）；`:116-128`（§4.1 结论账本无 C6 臂）；`docs/t027/REMAINING_SLICES.md:443`（B3b 台账条目同口径） |
| C6 | (a) 可复现注入或对照证据 | 缺失（无任何 C6 注入轮） |
| C6 | (b) 区分力门 | 缺失（无 C6 负/正对照） |
| C6 | (c) 机制证据 | 缺失（无 C6 机制令牌记录） |
| C6 | 初判 | 需用户裁决（三判据全缺，不可写 `proven`） |
| C2（起跑消息示例机制，非本片命名） | 机制描述 | `MoveFileExW` + `MOVEFILE_WRITE_THROUGH (0x8)` 的 no-replace move 取代 `os.Link`。出处：`docs/validation/t027-b3b-durability.md:45` |
| C2 | B3b 注入矩阵覆盖 | 覆盖。依据行：`docs/validation/t027-b3b-durability.md:122`（§4.1 C2 候选级 `proven`（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证），decisionSets=[record-linked,parent-synced] 两格全 eligible） |
| C2 | (a) 可复现注入或对照证据 | 齐备（主矩阵 85 轮全有效，c2 两决定格各 5/5 存活、0 丢失；`:79`、`:90-91`） |
| C2 | (b) 区分力门 | 齐备（判别力门 ok=true：负对照 5/5 丢失、正对照 5/5 存活；`:128`） |
| C2 | (c) 机制证据 | 齐备（两决定格机制令牌 5/5；变异臂证明装置会拒绝机制未执行的存活；`:90-91`、`:104-107`） |
| C2 | 初判 | 可升级 proven（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证） |

注：本表 C2 行的 proven 为候选级、前提 P′ 内，非平台级；P′（VirtualBox guest / 25H2 / hostiocache=false / ~1KiB / N=5）是否覆盖物理机目标部署环境未论证。引用时不得截断括注。

子注记（2026-11-04 用户裁定追加，原文不动）：（完整限定：仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）。

## 5. 红队闭合自检表（进 ② 前，只读）

本表为主控自检，不是独立复核；⑤/⑥ 须独立对照红队两轮原文确认。

主控对照两轮红队原文逐条自检 r2b 闭合质量（2026-11-04；只读，未改文件）。

| 红队发现 | 级 | r2b 闭合位置 | 自检结论 |
|---|---|---|---|
| 一轮 High 1 异常路径三处（空槽/损坏槽混淆、丢槽重建、可见性与 flush 混同） | High | r2 A2 未提交中间态冻结与校验顺序、三时点、锁内 truncate 重试；二轮「未完全闭合」的依赖（prebuild-set 存在性）由用户裁定接受降级 | 闭合（降级经用户裁定） |
| 一轮 High 2 接缝替换范围与读者覆盖不匹配 | High | r2 A2 接缝范围（R/C 专用 `publishReplayRCNoReplace`、ownership/launch/terminal 原路径逐字节不变）；二轮「无剩余发现，设计层闭合」 | 闭合 |
| 一轮 High 3 预建生命周期非等价解析 | High | r2 A2 预建生命周期按 Q1=B 重写 + 「本裁定推翻 r1 的 High 6 判断」；二轮「无剩余发现，闭合」 | 闭合 |
| 一轮 M4 核心声明论证缺项 | Medium | r2 D 有条件声明 + E/F 三测量分开；二轮未再点名 | 闭合 |
| 一轮 M5 区分力测量对象未定义 | Medium | r2 F 物理槽存活/未提交载荷存活/逻辑记录可见三测量；二轮未再点名 | 闭合 |
| 一轮 M6 LockFileEx 仲裁协议不足 | Medium | r2 A2 仲裁协议条款（范围/同句柄/有界等待/锁序/解锁/映射视图排除）；二轮未再点名 | 闭合 |
| 一轮 M3 引用纪律未完整传递 | Medium | 引用括注逐处补齐 + 二轮 Low 判定表注记由用户裁定加子注记（报告 §4） | 闭合 |
| 一轮 Sol 事项①②③（均需修订） | — | ① Q1=B 裁定推翻 r1 解析；② 接缝收窄；③ 仲裁协议条款 | 闭合 |
| 二轮 High 历史全集绑定缺口 | High | r2b 改动点① 降级句逐字 + 「本降级经用户 2026-11-04 裁定接受」+ ADR-029 补充节 | 闭合 |
| 二轮 Medium 允许集不一致 | Medium | r2b 附件 C 允许集表（进程 kill 与断电恢复两场景逐切点自洽）+ 表下注释 | 闭合 |
| 二轮 Medium BA3 等值重放错误 | Medium | r2b 补充 B·BA3 两臂对照重写 | 闭合 |
| 二轮 Low 顺序论证适用域过宽 | Low | r2b 补充 A-4 限定「首次预建 + 当前单阶段记录格式」 | 闭合 |
| 二轮 Low 自身哈希口径冲突 | Low | r2b 补充 A·A-1 口径句修订 | 闭合 |
| 二轮 Low API 形态与执行路径混写 | Low | r2b 补充 B 三层分述、删指令级保证 | 闭合 |
| 二轮 Low 判定表注记缺完整括注 | Low | 报告 §4 子注记（原文不动，用户裁定） | 闭合 |

结论：全部发现 = 闭合，无部分闭合/未闭合；② 可进场。

EN 镜像已随 ⑧a 补齐（本表 CN/EN 逐行对等）。

## 6. 审查链

- ① 双打：r1 → 第一次 Sol 红队 `FINDINGS` → r2 → 第二次 Sol 红队 `FINDINGS` → r2b 闭合 → r2c 激活语义微修；两次红队原文归档于 `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md` 与 `red-team-sol-review-r2.md`，四轮设计原文同目录归档。
- ②a–d 实现：契约先行（§7 C6 判定路径 + 附件 C）→ R/C 接缝收窄 → store 级激活。
- ③a–c：红测先行、BA3 两臂重写、三条缺口测试补齐（`TestC6PrebuildNoopWithoutActivation`、`TestC6PrebuildPlanDriftIsCorruption`、`TestC6PrebuildRunIDMismatchIsRootConflict`）。
- ④ 门禁全绿：默认、c6slot 与 linux 的 build/vet，以及 tree 与 index 两作用域。
- ⑤ 预审：`PRE-REVIEW: PASS WITH FIXES`（0H/0M/4L，L1–L4 处置见 §7）。
- ⑥ 独立扫描：已完成（逐字转录见下）。
- ⑦ 终审：Codex `FINAL REVIEW: PASS`（0 Critical / 0 High / 0 Medium / 0 Low；红队等价性论证 §7.2.1 口径通过；⑧a 先于 ⑦ 的时序偏离影响面判断成立）；原文归档于 evidence 包 `07-final/review.md`。
- ⑦ 复审（#1/3，收窄输入：②e 死代码迁移）：Codex `FINAL REVIEW: PASS`（0 发现；纯搬迁判定、死代码消除、c6slot 完整性、构建标签、import 卫生五面全过）；原文追加于 `07-final/review.md` 第二轮。
- ⑧/⑨/⑩：待闭环。

⑥ 独立扫描逐字转录（标题按门禁要求写作 Section A–E，正文逐字）：
```sixstage
## Section A Scope
- Read-only independent scan of the §7 durability contract, the C6 prebuild activation path, the Windows slot implementation, and the no-op fallback path.
- Scope excludes code edits or any dependency on the earlier ⑤ pre-review conclusion.

## Section B Findings
- No findings.
- The contract wording and the implementation remain aligned: the C6 path is build-gated, store-activated only after a successful prebuild, and the non-activated build remains a no-op.

## Section C Evidence
- The production assembly order is explicit in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L129-L195): `NewAgentRunnerReplayStoreForRun` is followed immediately by `PrebuildReplaySlots(...)`, before any dispatch.
- The activation semantics are defined in [internal/adapters/agent_runner_replay_store_c6_enable.go](internal/adapters/agent_runner_replay_store_c6_enable.go#L13-L115): prebuild success sets the store active, failure leaves the run fail-closed, and reopen verification is idempotent.
- The Windows slot behavior is enforced in [internal/adapters/agent_runner_replay_store_c6_windows.go](internal/adapters/agent_runner_replay_store_c6_windows.go#L15-L239): header format, payloadLen validation order, empty-slot freeze semantics, and no-replace publication rules.
- The unactivated fallback is pinned in [internal/adapters/agent_runner_replay_store_c6_noop_test.go](internal/adapters/agent_runner_replay_store_c6_noop_test.go#L12-L87): a non-`c6slot` build cannot activate the slot path or touch the replay root.
- Fresh validation command: `cd /d d:\LZProjects\prfrail && go test ./... && go test -tags c6slot ./internal/adapters`
  - Result: all packages reported `ok`; the adapters package completed successfully with no failing assertions.

## Section D Falsifiability audit
- The Windows C6 behavior is explicitly testable and tested in [internal/adapters/agent_runner_replay_store_c6_windows_test.go](internal/adapters/agent_runner_replay_store_c6_windows_test.go): empty-header creation, committed-slot rewrite refusal, corruption ordering, uncommitted intermediate state handling, prebuild activation, and fail-closed prebuild conditions.
- Those tests are structured as negative/positive checks and therefore provide a valid falsification surface for the contract.

## Section E Unverified / overturn conditions
- None open for the current implementation.
- Any future platform or premise change would require re-running the same validation set and re-checking the C6 activation enablement and durability claims under that new premise.

INDEPENDENT SCAN: PASS
```

## 7. 偏离与未闭环

- L1（附件 C 悬空引用）：已由附件 C 嵌入 CONTRACTS{,_EN} §7 闭环，A2 引用改指「本节附件 C」。
- L2/L3（BA1/BA2 未执行断言、漂移与 runId 分支无钉针）：已由 ③c 三条测试闭环（测试名见 §6 ③a–c 行）。
- L4（A2 尾句未带「限定观察面」限定）：维持原文不改；改写需用户授权，如实登记为本片未闭环项。用户 2026-11-04 裁定维持原文（理由：⑦ 已判 `FINAL REVIEW: PASS`，L4 为 Low 且 ⑤ 已如实登记「授权边界内未闭环」；不改不构成伪通过），本项自此闭环。
- r2c 激活语义微修：③ 暴露的「构建级 vs store 级」冲突，按 design-r2c 选 store 级激活，A2 追加激活语义句。
- SKIP=1：`TestSpawnRecoveryHelper` 属子进程设计行为（默认构建同样 SKIP），不是失败。
- gopls 构建标签告警：仅出现在 `c6_enable` 文件的工具配置态，不影响 build/vet/test。
- 时序偏离（偏离登记第 1 条，不并 §5.1 的 D 计数；主控自主发现）：⑧a 文档同步在 ⑦ 终审之前执行（准则流水线顺序为 ⑦ 后 ⑧）；原因 = 主控编排失误；影响面 = ⑧a 产物（六文档 + ADR-030 + 报告 §6/§7）将纳入 ⑦ 输入包同审，⑦ 结论覆盖最终态；处置 = 如实登记，不回溯重排。
- ②e 死代码迁移（用户 2026-11-04 澄清触发）：`createReplaySlotWindows` 创建链在默认 Windows 构建无调用点（gopls `unusedfunc`）→ 迁入新文件 `agent_runner_replay_store_c6_windows_create.go`（windows && c6slot），语义逐字节等价，⑨ 三态全绿；R 计 1（发现后的修复重跑），R=1 < 2 不触发高风险行；⑦ 收窄复审见 §6。

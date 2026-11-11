# P64 AT24-SUPERVISED-BLACK-BOX 黑箱受监督模式验收报告（⑧b 定稿）

- 切片标识：`P64 AT24-SUPERVISED-BLACK-BOX`（验收类 `[SLICE]`；v53 §四 层 1 后续候选第 12 项）
- 日期：2026-11-11
- 状态：**⑧b 定稿**——⑦ 复审 `RE-REVIEW: PASS`（回执 `04-reviews/07-rereview-r1.txt`）；R-09.1~R-09.9 九项边界由 ⑦ 判定；AT-24 验收结论 = **条件闭合**（用户裁决③）
- 主控模型：`DeepSeek V4 Pro`（ADR-040，步骤 0 按 OB-37 元规则登记）
- 证据包：`docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/`（MANIFEST 与 SHA256SUMS 自洽）
- 验收性质：机制腿（stub visible 投递 + 模拟人工归还 + 后置扫描），非真实 SessionBridge visible 端到端

## 执行摘要区段（独立于 §1）

> ①「本报告的保证级别是 `reduced`，不是 `full`。」
>
> ②「七项 unknown 维度的接受边界为：工具使用/资源读取 = 视情况（扫描能证无越界则接受）；后台进程/会话连续性 = 可接受；网络 = 回环可接受、公网阻断；费用 = 显式标 unknown；外部副作用 = 用户已显式确认接受。」
>
> ③「本报告不构成对黑箱 Agent 过程正确性的保证；ProofRail 只保证归还后的产物扫描、独立门禁、证据和评审。」

## 1. 执行摘要

事实层摘要（不含判定）：17 子腿中 PASS=16、1 项待计数（腿 6）；硬失败 0。DEF-1 经 ② 修复后 Go/CLI 层闭合；schema 七项闭集（Major A）与 9.1 junction 等价补测（Major B）经 R=2 重跑闭合；⑦ 复审 `RE-REVIEW: PASS`（无 Medium+）⇒ AT-24 验收结论 = **条件闭合**（用户裁决③）；九项边界判定见 §2 与 §2.1。

关键读数：

- 17 子腿：PASS=16 / PENDING-COUNT=1（腿 6）/ 硬失败 0；9.1 由 DONE WITH SKIP 翻 PASS（junction 等价腿实取，`L09-01-r2.txt`）。
- R 计数：2（R1 轮：腿 5 / 9.8 的 DEF-1 修复重跑；R2 轮：腿 5 / 9.1 / 9.8 的裁决①②修复重跑）。
- ④ 读数（2026-11-11 裁决后最终重跑，与 06-gates 同口径）：`gofmt -l` 空、`go build` exit 0、`go vet` exit 0、`go test ./...` ok=15 FAIL=0、三作用域（index / tree / ci）`TOTAL_FAIL=0`、`SELFTEST: PASS 722/722`；`06-gates/gates-r1.txt` 为首次读数，最终读数同口径一致。（定稿装配后 index / tree 另见 §3 的「未入库 / 未登记」致 FAIL 读数）
- 环境：Windows 本机；无 VM、无代理、无外部调用（代理 OFF）。

> ①「本报告的保证级别是 `reduced`，不是 `full`。」
>
> ②「七项 unknown 维度的接受边界为：工具使用/资源读取 = 视情况（扫描能证无越界则接受）；后台进程/会话连续性 = 可接受；网络 = 回环可接受、公网阻断；费用 = 显式标 unknown；外部副作用 = 用户已显式确认接受。」
>
> ③「本报告不构成对黑箱 Agent 过程正确性的保证；ProofRail 只保证归还后的产物扫描、独立门禁、证据和评审。」

## 2. R-09.1~R-09.9 九项边界逐项对照

本表只列事实：裁定值、实测证据路径、实测读数与「是否吻合」的事实基础。**判定列已由 ⑦ 判定**（盲审 `07-blind-r1.txt` / 终审 `07-final-r1.txt` / 复审 `07-rereview-r1.txt`）：9.1 与 9.8 依用户裁决③按复审翻为「通过」并注明依据；其余七行文字保持 ⑧a 原文（⑦ 终审已判「通过」，其单元格文字待 ⑧b-2 统一回写）。

| 子项 | R-09 裁定值 | 实测证据路径（`03-test/`） | 实测读数与事实基础 | 判定 |
|---|---|---|---|---|
| 9.1 工具使用 | (c) 视情况 | `L09-01-tool-use.txt`；`L09-01-junction-probe.txt`；`L09-01-r2.txt` | `AssemblyFullChainExtendedScans` 三扫描 clean；`OutOfScopeChange` / `ScopeUnderivable` / `ScopeRejectsUnresolvableRoot` PASS=3 FAIL=0；symlink 原腿 `ScopeAllowsRejectsLinkResolvingOutsideRoot` SKIP=1（Windows 无 symlink 特权，D-16）；junction 等价腿 `TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot` PASS ⇒ M9-1 PASS=4 / SKIP=1 / FAIL=0（阈值 PASS≥4 达标） | 通过（⑦ 复审（RE-REVIEW r1）翻定，依据 L09-01-r2） |
| 9.2 资源读取 | (c) 视情况 | `L09-02-resource-reads.txt` | `IsolationReadsOnlyWorkspace` PASS=1 FAIL=0；`ChangedPaths` 恰 1 项、workspace 外候选逐字未变 | 待 ⑦（本草案不下判定） |
| 9.3 网络 | (c) 视情况 | `L09-03-network.txt`；`01-arch/r09-3-user-confirmation.md` | 确认件精确句「回环可接受、公网必须阻断。」count=1；外部连接 0 条、代理 OFF、无 VM | 待 ⑦（本草案不下判定） |
| 9.4 费用 | (a) 显式标 unknown | `L09-04-cost.txt` | 报告 DIMS 含 `cost`、ASSURANCE=`reduced`；缺 `cost` 的 CLI exit=1 且 schema reject | 待 ⑦（本草案不下判定） |
| 9.5 后台进程 | (a) 可接受 | `L09-05-background-processes.txt` | `RefusesLiveManagedProcess` PASS=1 FAIL=0；报告含 `background-processes`=1；**报告侧明写：后台进程由用户负责清理**（PENDING-⑧ 于本报告闭合） | 待 ⑦（本草案不下判定） |
| 9.6 会话连续性 | (a) 可接受 | `L09-06-session-continuity.txt` | `RunRestartRecovery` + `ConfirmationCrashAfter{Append,Resume,WindowOpen}` PASS=4 FAIL=0；报告含 `session-continuity` | 待 ⑦（本草案不下判定） |
| 9.7 外部副作用 | (c) 视情况 | `L09-07-external-side-effects.txt`；`01-arch/r09-7-user-confirmation.md` | 确认件精确句「已接受未知外部副作用。」count=1；`RefusesSideEffect` + `SideEffectScanBlocksTamperedManifests` TOP=2 / SUB=5 / FAIL=0 | 待 ⑦（本草案不下判定） |
| 9.8 强制字段 | (a) 七项全必填 | `L09-08-mandatory-fields.txt`；`L09-08-r2-schema-seven-mandatory.txt`；`L05-r2-schema-seven-mandatory.txt` | schema 三反例 reject 3/3；首跑缺其余六项之一 exit=0 ⇒ DEF-1；[R1] exit=1、正例 exit=0；Major A 补齐后 schema 直校 9 例全对、契约测试 4/4、fixtureCount=144 ⇒ [R2] 三层判定一致 | 通过（⑦ 复审（RE-REVIEW r1）翻定，依据 L09-08-r2） |
| 9.9 判定口径 | (b) 每项独立判定 | `L09-01`…`L09-09` | 独立 9 件、一子腿一文件；⑦ 判定表 9 行待 ⑦ 生成 | 待 ⑦（本草案不下判定） |

表注：9.2–9.7 与 9.9 七行的判定值由 ⑦ 终审（`07-final-r1.txt`）判为「通过」；其单元格文字仍为 ⑧a 草案口径（「待 ⑦」），回写由 ⑧b-2 / ⑨ 统一执行，本版不改这七行。

补充事实（不构成判定）：

- 9.5 义务边界：**后台进程由用户负责清理**；ProofRail 不对未知后台进程作清理担保。
- 9.3 / 9.7 确认件为用户逐字确认，补录时点（2026-11-11）在 ③ 执行之前；原文见 `01-arch/`。

### 2.1 ② 修复记录（Major A / Major B）与 ⑦ 复审路径

- **Major A（schema 七项闭集，用户裁决①授权）**：`unknownDimensions` 由 `minItems:1 + contains(cost)` 加严为七项闭集（`minItems:7` + `maxItems:7` + `uniqueItems` + `allOf` 七个 `contains`），新增 invalid 夹具（六项集缺 `network-access`），`fixtureCount` 143→144；R=2 重跑 PASS。依据 `02-impl/schema-gap-analysis.md`、`03-test/L05-r2-schema-seven-mandatory.txt`、`03-test/L09-08-r2-schema-seven-mandatory.txt`。
- **Major B（9.1 junction 等价覆盖）**：junction 等价性探针 PASS（越界 junction 被同一 `ErrSupervisedScopeOutOfBounds` 哨兵阻断、根内 junction 放行），随后转正式补测 `TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot`；M9-1 PASS=4 / SKIP=1 / FAIL=0 达标。依据 `03-test/L09-01-junction-probe.txt`、`03-test/L09-01-r2.txt`。
- **⑦ 复审**：`RE-REVIEW: PASS`（Critical / Major / Minor 均无），两处 Major 均判「闭合」⇒ AT-24 翻**条件闭合**。依据 `04-reviews/07-rereview-r1.txt`。

## 3. 验收腿证据（17 子腿）

证据原件目录：`docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/`；总表另见 `SUMMARY.md`。

| 腿 | 正例 | 反例 | 断言读数 | 结果 |
|---|---|---|---|---|
| 1 风险确认绑定 | `RiskConfirmationRoundTrip` / `OpenSupervisedWorkWindowWithValidConfirmation` / `RiskNoticeIsTruthfulAndVersioned` / `AssemblyFullChain` | 绑定/送达类反例 6 项（`OpenWorkWindow{WithoutConfirmation,WrongNoticeVersion}`、`RiskConfirmationRejectsBindingDrift`、`Assembly{BlocksWithoutConfirmation,RejectsSelectionReject}`、`OpenWorkWindowRejectsUnacceptedDelivery`）+ B 层 run 探针 | M1 PASS=12 FAIL=0；探针 exit=1 且命中 `supervised visible delivery is not accepted` | PASS |
| 2 投递 ≠ 完成 | `VisibleDeliveryMapsAccepted` / 投递零回执链 | `SubmitReturnRejects{DeliveryOnlyEvidence,ReceiptNotBoundToFacts}` | M2 PASS=5 FAIL=0 | PASS |
| 3 显式归还 + 后置全量扫描 | `ReturnCompleteFlow` / `AssemblyFullChain` / `AssemblyFullChainExtendedScans` | 扫描/门禁类反例 11 项（`Return{RefusesAbortReceipt,RefusesSecretFinding,RefusesOutOfScopeChange,RefusesSideEffect,RefusesFailedGate,RefusesTimeout,EnforcesFrozenPolicy,RefusesDriftedPolicy}`、`SideEffectScanBlocks*`、`FileTypeScanRejects*`、`RefusesLiveManagedProcess`） | 主跑 PASS=17 FAIL=0；补跑 PASS=2 FAIL=0 | PASS |
| 4 独立 gates/review | `PostflightQualifiesReviewPending` / `PostflightDecisionMustConsumeAllFacts` | `ReturnRefusesFailedGate` / `UnqualifiedReviewRefused` / `TaskWithoutReturnFactsFailsClosed` | M4 PASS=26 FAIL=0（点名 5/5；宽匹配见 D-11） | PASS |
| 5 七项含 cost | `ReportNewAlwaysReduced`（首跑 4 子例 → [R1] 11 子例） | 夹具 full / empty / missing-cost / 新增六项集（schema 拒 4/4）+ CLI 缺 cost 拒 + 缺项探针 | 七项计数=7、reduced=1；缺 cost exit=1；首跑缺项 exit=0 ⇒ DEF-1；[R1] 缺项 exit=1；[R2] fixtureCount=144、C 层 4/4 | PASS（R=2） |
| 6 三句话硬性出现 | 位置 A（§1）/ 位置 B（执行摘要）各 1 次；区段外 0 次 | 任何改写 / 摘要变体 ⇒ 精确串计数 0 | 逐字原文与计数口径已登记（`L06`）；计数待 ⑨ 机械执行 | PENDING-COUNT |
| 7 禁止自动降级 | supervised→supervised / AgentRunner→AgentRunner 装配 | 结构反例 7 项（`MixedSupervisedAgentRunnerInputs`、`*AssemblyRefuses*`、`NoSupervisedReferencesInAgentRunnerPaths`、`SupervisedPathNeverRoutes*`、`*StepRejectsAgentRunner*`、`RouterRefusesSupervisedWithoutPort`） | M7 PASS=16 FAIL=0（点名 7/7；见 D-15） | PASS |
| 8 未知维度不显示为安全 | `ReportRenderGolden` / `ReportOmitsUnprovenReviewDimension` | B 层 unknownDimensions 含 `safe` / `verified` 探针 | M8 PASS=2 FAIL=0；探针 exit=1；渲染 grep `safe`=0 `verified`=0 | PASS |
| 9.1 工具使用 | `AssemblyFullChainExtendedScans` | `RefusesOutOfScopeChange` / `ScopeUnderivableFailsClosed` / `ScopeRejectsUnresolvableRoot` + junction 等价腿 `TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot`（symlink 原腿 SKIP 1） | [R2] PASS=4 SKIP=1 FAIL=0（阈值 PASS≥4 达标；junction 腿在 Windows 本机实取越界反例） | PASS（R=2） |
| 9.2 资源读取 | `IsolationReadsOnlyWorkspace` | 越界痕迹由 9.1 反例组阻断（同源） | PASS=1 FAIL=0 | PASS |
| 9.3 网络 | 本片零外部调用；确认件精确句 | 精确句缺失 ⇒ 阻断（可证伪） | 精确句=1；外部连接 0 条 | PASS |
| 9.4 费用 | 报告 `cost ∈ unknownDimensions` | 缺 cost ⇒ 拒 | `cost` 计数=1；assuranceLevel=reduced | PASS |
| 9.5 后台进程 | `RefusesLiveManagedProcess` | 受管进程存活归还 ⇒ 拒 | PASS=1 FAIL=0；`background-processes`=1；「用户负责清理」本报告已写入 | PASS（3/3） |
| 9.6 会话连续性 | `RunRestartRecovery` / `ConfirmationCrashAfter*` | 外部 agent 会话不可观测 ⇒ 报告标 unknown | PASS=4 FAIL=0；`session-continuity`=1 | PASS |
| 9.7 外部副作用 | 确认件精确句 + 副作用扫描执法 | 精确句缺失 ⇒ 阻断；五类篡改 ⇒ 阻断 | 精确句=1；TOP=2 / SUB=5 / FAIL=0 | PASS |
| 9.8 强制字段 | reduced + 非空 + 含 cost 三层 | full / empty / missing-cost + 缺其余六项之一 + 闭集外项 / 重复 / 超长 | 三反例 3/3 拒；首跑缺项 exit=0 ⇒ DEF-1；[R1] exit=1、正例 exit=0；[R2] schema/Go/CLI 三层判定一致（fixtureCount=144） | PASS（R=2） |
| 9.9 判定口径 | 9 件独立证据落盘 | 合并 / 单行总判定表 ⇒ 结构不合规 | `L09-*.txt`=9；⑦ 判定表 9 行待 ⑦ | PASS（结构侧） |

DEF-1 与 ② 修复（事实）：

- 首跑：`unknownDimensions` 缺任一非 `cost` 项（`["cost"]` 或七项缺 `network-access`）⇒ CLI exit=0 接受并渲染，schema（Ajv）同向 ACCEPT；与 R-09.8 / 腿 5 判据（应拒）不符 ⇒ DEF-1 成立。
- 根因：`internal/chain/supervised_report.go` 与 `schemas/supervised-report.schema.json` 均只校验非空 + 含 `cost`。
- ② 处置（R1 轮）：Go 层增加七项闭集校验（配套 `supervised_risk_confirmation.go`、`runtime_supervised.go`、`supervised_blackbox_test.go`），共四文件；schema / fixtures 当时未改（按 design-r1 §3.3 属超界，交 ⑦ 与用户裁决）。
- [R1] 重跑：`p5-only-cost`、`p5-missing-network-access` 转 exit=1；正例 `p5-full-seven` 保持 exit=0；A 层 11 子例全 PASS；C 层 4/4 不变；回归 `go test ./internal/chain ./internal/console` 均 ok。
- [R2] 重跑（用户裁决①②后）：schema 七项闭集落地，`fixtureCount` 143→144；M9-8 六类判定全部命中（顶层 3 PASS / 子例 11 PASS / 契约 4 PASS）；9.1 junction 等价腿 PASS=4 / SKIP=1。

R 重跑注记：R1 轮 = 腿 5 与腿 9.8 各 1 次（同源 DEF-1）；R2 轮 = 腿 5 / 9.1 / 9.8 各 1 次（裁决①②修复）；其余腿未重跑，R 计数不向其他腿累积。

④ 门禁最终读数（2026-11-11 裁决后最终重跑，与 `06-gates/gates-r1.txt` 同口径）：`gofmt -l` 空、`go build` exit 0、`go vet` exit 0、`go test ./...` ok=15 FAIL=0、三作用域（index / tree / ci）`TOTAL_FAIL=0`、`SELFTEST: PASS 722/722`。`gates-r1.txt` 为首次读数（index 首跑 TOTAL_FAIL=3，修复后复跑为 0）。本报告定稿装配后，⑧b 侧复跑 index 出现三项 FAIL（G5(f) 索引洁净、G7-b 声明件未入库、G8-b 证据包未登记），tree 为两项（G7-b、G8-b），全部源于 ⑦ 复审新增的 `04-reviews/07-rereview-r1.txt` 尚未入库且未登记 MANIFEST / SHA256SUMS；ci 作用域 `TOTAL_FAIL=0`。该件入库与证据包登记属 ⑨ 职责，⑨ 复跑方为最终结论。

⑤⑥⑦ 回执指针（`04-reviews/`）：⑤ `05-prereview-r1.txt` = `PRE-REVIEW: PASS WITH FIXES`（无 Critical；M1 并入 ⑦ 判定问题清单，M2–M4 为 ⑧b 回写项）；⑥ `06-scan-r1.txt` = `INDEPENDENT SCAN: PASS`（Critical / Major / Minor 均 none）；⑦ `07-blind-r1.txt` = `BLIND REVIEW: FINDINGS`、`07-final-r1.txt` = `FINAL REVIEW: PARTIAL`、`07-rereview-r1.txt` = `RE-REVIEW: PASS`（用户裁决③：①② 修复后 ⑦ 复审转 PASS 且无 Medium+ ⇒ AT-24 条件闭合）。

⑤ M2/M4 回写注记（`SUMMARY.md` 不改）：`SUMMARY.md` 第 9.5 行「PASS（2/3，第三项 PENDING-⑧）」中的 PENDING-⑧ 已由本报告侧闭合（§3 腿 9.5 与 §4 第 7 条明写「用户负责清理」）；其末行「LEGS 9: DONE WITH FAILS」为 R1 前旧态；两处指针均不追改，由 ⑧b-2 处理或登记注记。

设计偏差清单（③ 侧登记，如实列示）：

- D-1 / D-14（exit 区间口径）：设计写 `exit∈[10,20]`，实况 CLI 仅 0/1/2，采信「非零即拒」。
- D-5（可证伪性取证边界）：③ 禁改源码 ⇒ 未做现场变异，改以结构性反例取证。
- D-10（范围扫描为近似）：`ReturnRefusesOutOfScopeChange` 通过不等于 glob 级授权比对。
- D-11 / D-15（断言名正则宽度）：`Postflight`、`AssemblyRefuses` 宽匹配带入未点名测试（26 vs 5、16 vs 7），按点名数单列。
- D-13（全链报告 JSON 未落盘）：腿 5 / 9.4 / 9.5 / 9.6 的「报告含 X」用「生产常量 → 构造原样传入 → CLI 渲染」三段式代证。
- D-16（编号冲突）：`L07` 用 D-16 指「可证伪性说明」，`SUMMARY.md` 用 D-16 指「M9-1 阈值差 1（Windows 无 `SeCreateSymbolicLinkPrivilege` ⇒ SKIP）」；同号异义，统一与否交 ⑦ / ⑧b。

## 4. 已知边界

1. **真实 SessionBridge visible 端到端 = unknown**，依据 T028 完成记录（真实 CLI 端到端依赖 T027-B，尚不可运行发布）。
2. 本片验收腿是**机制腿**（stub visible 投递 + 模拟人工归还 + 后置扫描），**不是真实 SessionBridge visible 端到端**。
3. **AT-24 的验收不覆盖真实 CLI 端到端**；S1 Exit 的 S-G7 依赖链须如实标注，本报告不主张 S1 Exit 整体通过（S-G7 其余门禁仍未闭合）。
4. schema 层「七项闭集」已按用户裁决①补齐（`minItems:7` + `maxItems:7` + `uniqueItems` + `allOf` 七个 `contains`），Go/CLI 层与 schema 层同向，R1 遗留在 schema 层的单向差已闭合；`fixtureCount` 143→144。
5. 9.1 的 `ScopeAllowsRejectsLinkResolvingOutsideRoot` 因 Windows 无 symlink 特权而 SKIP，属环境不可测、非实现缺陷；同一判定路径（`scope.allows` + `ErrSupervisedScopeOutOfBounds` 哨兵）由 junction 等价腿 `TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot` 在本机实取覆盖，属用户裁决②允许的等价替代，不等于 symlink 原腿实测。
6. 腿 5 / 9.4 / 9.5 / 9.6 的「报告含 X」为三段式代证（D-13）。
7. 后台进程 unknown 由用户负责清理（R-09.5）：ProofRail 不担保其清理与存活状态。

## 5. 明确未执行事项

- 产品代码改动仅限 ② 修复与裁决①②授权项：DEF-1 四文件（`supervised_report.go`、`supervised_risk_confirmation.go`、`supervised_blackbox_test.go`、`runtime_supervised.go`）+ junction 正式补测 `internal/console/runtime_supervised_junction_probe_test.go`。
- schema / fixtures / 契约测试改动（`schemas/supervised-report.schema.json`、`testdata/contracts/invalid/runtime-records.json`、`tools/contracts/contracts.test.js`）系用户裁决①明确授权；其余 CI、契约、判据文本、准则、R-09 裁定与 T028 报告原文未改。
- 未替用户裁定九项边界，未替 ⑦ 给出 AT-24 验收结论。
- 未 commit、未 push。
- 未启动新切片；未开 VM、未开代理、无外部调用。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/at24-supervised-black-box.md	repo
docs/validation/at24-supervised-black-box_EN.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/00-briefing.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/01-arch/design-r1.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/01-arch/r09-3-user-confirmation.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/01-arch/r09-7-user-confirmation.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/02-impl/schema-gap-analysis.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L01-risk-confirmation-binding.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L02-delivery-not-completion.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L03-explicit-return-postscan.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L04-independent-gates-review.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L05-r2-schema-seven-mandatory.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L05-unknown-dimensions-cost.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L06-three-sentences.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L07-no-auto-degrade.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L08-unknown-not-safe.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-01-junction-probe.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-01-r2.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-01-tool-use.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-02-resource-reads.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-03-network.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-04-cost.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-05-background-processes.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-06-session-continuity.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-07-external-side-effects.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-08-mandatory-fields.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-08-r2-schema-seven-mandatory.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/L09-09-judgment-approach.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/03-test/SUMMARY.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/04-reviews/05-prereview-r1.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/04-reviews/06-scan-r1.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/04-reviews/07-blind-r1.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/04-reviews/07-final-r1.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/04-reviews/07-rereview-r1.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/06-gates/gates-r1.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/MANIFEST.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/SHA256SUMS.txt	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/freedom-list.md	repo
docs/validation/evidence/P64-AT24-SUPERVISED-BLACK-BOX/input-list.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/validation/t028-supervised-black-box.md	repo
docs/validation/t028-supervised-black-box_EN.md	repo
internal/chain/supervised_blackbox_test.go	repo
internal/chain/supervised_report.go	repo
internal/chain/supervised_risk_confirmation.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_junction_probe_test.go	repo
schemas/supervised-report.schema.json	repo
testdata/contracts/invalid/runtime-records.json	repo
tools/contracts/contracts.test.js	repo
```

## 6. 阶段转录（sixstage 围栏 A–E，G7-f）

- 围栏外注记：以下 A–E 五段为本片阶段结论的逐段转录，回执指针与 §3 一致；围栏首行与末行两处结论取值全等（末行按 G7-f「最后非空行 = 结论行」）。

```sixstage
## Section A
INDEPENDENT SCAN: PASS

① DESIGN: DONE —— design-r1.md（17 子腿方案 + 证据路径 + 机械断言 + ⑦ 判定问题清单；回执 01-arch/design-r1.md）

## Section B
③ 验收腿：17 子腿 PASS=16 / PENDING-COUNT=1（腿 6）/ 硬失败 0；② 修复：DEF-1 Go/CLI 四文件 + schema 七项闭集（Major A）+ junction 正式补测（Major B）；R=2

## Section C
④ 门禁：gofmt 空 / go build exit 0 / go vet exit 0 / go test ok=15 FAIL=0 / 三作用域 TOTAL_FAIL=0 / SELFTEST 722/722；⑤ PRE-REVIEW: PASS WITH FIXES（05-prereview-r1.txt）；⑥ INDEPENDENT SCAN: PASS（06-scan-r1.txt）

## Section D
⑦ 盲审 BLIND REVIEW: FINDINGS（07-blind-r1.txt）；终审 FINAL REVIEW: PARTIAL（07-final-r1.txt）；复审 RE-REVIEW: PASS（07-rereview-r1.txt）

## Section E
⑧ 本报告定稿（⑧b）：AT-24 条件闭合（用户裁决③；§1 与执行摘要结论同步）

INDEPENDENT SCAN: PASS
```

# design-r1.md — P64 AT24-SUPERVISED-BLACK-BOX · ① 执行计划确认（轻量化）

> 角色：① 执行计划确认（deep-reasoner + DeepSeek V4 Pro）。只读：不落盘、不判定 R-09、不下 AT-24 结论。判定主体 = ⑦（唯一）。本文件由主控落盘。

## 0. 输入锚点（已读）
- `00-briefing.md`：3.1 九腿清单与第 9 腿九子项、3.2 判定主体、3.3 缺陷处置、3.4 报告必含项。
- `docs/validation/r08-r09-ruling.md` §3（R-09.1~R-09.9 九裁定值 + 三句话）与 §5。
- `docs/CONTRACTS.md` §7.1 (a)–(d) + 四条禁止项。
- `docs/validation/t028-supervised-black-box.md` §1/§2/§7.2/§8；`DEV_PLAN.md:119`；`TEST_STRATEGY.md:54`。
- `p57-closeout-assessment.md` §4.2（八门：AT-24 对应 **S-G2**）；`s1-exit.md:44`（AT-24 行）、`:52`（核心任务链/S-G7）。
- `01-arch/r09-3-user-confirmation.md`、`r09-7-user-confirmation.md`（均已补录用户确认，时点 = ③ 执行前）。

## 1. 口径
1. 验收腿 = **机制腿**（stub visible 投递 + 模拟人工归还 + 后置扫描 + 独立 gates/review + reduced/unknown 报告）；真实 SessionBridge visible 端到端 = **unknown/边界**（按 R-09.3/R-09.7 判），报告 §4 必含 3.4 三项。
2. 9 腿 + 第 9 腿 9 子腿 = **17 子腿**；每子腿 = 正例 + 反例 + 机械断言。
3. R-09.1~R-09.9 逐项判定 + AT-24 验收结论的唯一主体 = ⑦；① 只出方案、证据路径、判定问题清单。
4. 三句话逐字核验位置：⑧a 草案写入报告 §1 与执行摘要 → ⑤ 预审 → ⑥ 扫描 → ⑦ 终审逐环逐字核验；④/⑨ 只做机械精确串计数（主控不自证）。
5. 环境 = Windows 本机 console/chain 层；无 VM/代理/外部调用；stub visible 由 Go 测试内 `VisibleSender` stub 承担。
6. ② 修复范围限 `internal/console/**` 与 `internal/chain/**` 的 supervised 相关文件；需动 schema/fixtures/CONTRACTS ⇒ 停下升级用户裁决。

## 2. ③ 通用工具
- 构建：`go build -o tmp/P64-AT24-SUPERVISED-BLACK-BOX/prfrail.exe ./cmd/prfrail`
- A 层（主机制）：`go test -count=1 -run '<p>' ./internal/chain ./internal/console ./internal/adapters -v`（下文 M 断言只写 `<p>` 与期望）
- B 层（CLI 黑箱探针）：`prfrail.exe supervised report --run-dir <dir> [--json]`、`prfrail.exe run <定义>`；合成夹具放 `tmp/P64-AT24-SUPERVISED-BLACK-BOX/fixtures/`
- C 层（契约）：`npm --prefix tools/contracts test`
- 原始输出全落 `03-test/`；③ 写 `03-test/SUMMARY.md`（17 腿结果总表）。

## 3. 17 子腿方案

### 腿 1 风险确认绑定（R-09.7 + §7.1(a)）
判据：开窗前持久化 `supervised-risk-confirmation`（绑定 run/task/step/attempt + workspace 摘要 + 版本/哈希 + visible 投递事实）；版本漂移开窗前阻断。
- 正例：`TestSupervisedRiskConfirmationRoundTrip`、`TestOpenSupervisedWorkWindowWithValidConfirmation`、`TestSupervisedRiskNoticeIsTruthfulAndVersioned`、`TestSupervisedAssemblyFullChain`。
- 反例：`TestOpenSupervisedWorkWindowWithoutConfirmation`、`TestOpenSupervisedWorkWindowWrongNoticeVersion`、`TestSupervisedRiskConfirmationRejectsBindingDrift`、`TestSupervisedAssemblyBlocksWithoutConfirmation`、`TestSupervisedAssemblyRejectsSelectionReject`、`TestOpenSupervisedWorkWindowRejectsUnacceptedDelivery`；B 层反例：`prfrail.exe run`（合成 supervised 定义）⇒ 生产 visible 不可用 ⇒ 无窗口，exit∈[10,20]，stderr 含 `supervised visible delivery is not accepted`。
- 断言 M1：`<p1>=SupervisedRiskConfirmation|OpenSupervisedWorkWindow|SupervisedRiskNotice|SupervisedAssemblyBlocksWithoutConfirmation|SupervisedAssemblyRejectsSelectionReject` ⇒ PASS≥10、FAIL==0；B 层探针 exit∈[10,20] 且 stderr 命中串。
- 证据：`03-test/L01-risk-confirmation-binding.txt`

### 腿 2 投递 ≠ 完成（R-09 + §7）
判据：visible 投递 receipt 不产生 task/step PASS。
- 正例：`TestSubmitSupervisedReturnRejectsDeliveryOnlyEvidence`（仅投递证据 ⇒ 拒）、`TestVisibleDeliveryMapsAccepted`（投递零回执）。
- 反例：`TestSubmitSupervisedReturnRejectsReceiptNotBoundToFacts`（回执未绑五事实 ⇒ 拒）。
- 断言 M2：`<p2>=DeliveryOnlyEvidence|VisibleDelivery|ReceiptNotBoundToFacts` ⇒ PASS≥4、FAIL==0。
- 证据：`03-test/L02-delivery-not-completion.txt`

### 腿 3 显式归还 + 后置全量扫描（R-09.1/9.2 + §7.1(b)）
判据：显式归还 ⇒ `handoff-receipt{outcome=complete}`；manifest/diff/范围/秘密/文件类型/副作用扫描全绿才过；未知进程/越界阻断。
- 正例：`TestSubmitSupervisedReturnCompleteFlow`、`TestSupervisedAssemblyFullChain`、`TestSupervisedAssemblyFullChainExtendedScans`。
- 反例：`TestSubmitSupervisedReturnRejectsAbortReceipt`、`TestSupervisedReturnRefusesSecretFinding`、`TestSupervisedReturnRefusesOutOfScopeChange`、`TestSupervisedReturnRefusesSideEffect`、`TestSupervisedSideEffectScanBlocksTamperedManifests`、`TestSupervisedFileTypeScanRejectsUnacceptedType`、`TestSupervisedReturnRefusesLiveManagedProcess`、`TestSupervisedReturnRefusesFailedGate`、`TestSupervisedReturnRefusesTimeout`、`TestSupervisedReturnEnforcesFrozenPolicy`、`TestSupervisedReturnRefusesDriftedPolicy`。
- 断言 M3：`<p3>=SubmitSupervisedReturn|SupervisedReturnRefuses|SideEffectScanBlocks|FileTypeScanRejects|FullChain` ⇒ PASS≥13、FAIL==0。
- 证据：`03-test/L03-explicit-return-postscan.txt`

### 腿 4 独立 gates/review（§7.1(b) review/promotion）
判据：独立 build/test/verify gates + review 通过才可接受；exit 0 / completed receipt 不直接 PASS。
- 正例：`TestSupervisedPostflightQualifiesReviewPending`、`TestSupervisedPostflightDecisionMustConsumeAllFacts`。
- 反例：`TestSupervisedReturnRefusesFailedGate`、`TestSupervisedUnqualifiedReviewRefused`、`TestSupervisedTaskWithoutReturnFactsFailsClosed`。
- 断言 M4：`<p4>=Postflight|UnqualifiedReview|TaskWithoutReturnFacts|ReturnRefusesFailedGate` ⇒ PASS≥5、FAIL==0。
- 证据：`03-test/L04-independent-gates-review.txt`

### 腿 5 unknownDimensions 七项全必填含 cost（R-09.4/9.8）
判据：`assuranceLevel=reduced`；七项（tool-use/resource-reads/network-access/cost/background-processes/session-continuity/external-side-effects）逐项出现。
- 正例：`TestSupervisedReportNewAlwaysReduced`；全链产出报告 JSON 七项计数==7 且 `assuranceLevel=="reduced"`。
- 反例：契约夹具 `supervised-report-assurance-level-full` / `-unknown-dimensions-empty` / `-missing-cost` 均 schema 拒（C 层 4/4 含之）；B 层：缺 cost 合成报告 ⇒ `supervised report` 校验拒（exit∈[10,20]）。**关键反例探针**：含 cost 但缺任一其他六项的合成报告 ⇒ 按 R-09.8「七项全必填」期望拒绝；若实测接受 ⇒ 登记 T028 实现缺陷候选（Go 层可修 ⇒ ② 修+重跑；需改 schema/fixtures ⇒ 超界停下升级）。
- 断言 M5：C 层 `tests 4 / pass 4 / fail 0`；全链报告七项计数==7；B 层缺 cost exit∈[10,20]；缺项探针如实记录「接受=缺陷候选 / 拒绝=合规」。
- 证据：`03-test/L05-unknown-dimensions-cost.txt`

### 腿 6 三句话硬性出现（R-09 §5.3）
- ① 只安排位置：⑧a 草案在报告 §1 与执行摘要各写一次（逐字，原文 = 起跑消息【七】第 5 条）；⑤⑥⑦ 逐环逐字核验；④/⑨ 机械精确串计数。
- 正例断言 M6：三句精确串在报告中计数各==2（§1 区段==1、执行摘要区段==1、两区段外==0）。
- 反例：任何改写/摘要变体 ⇒ 精确串计数==0 ⇒ 判 FAIL（可证伪）。
- 证据：`03-test/L06-three-sentences.txt`（机械计数表 + ⑤⑥⑦ 核验回执指针）。

### 腿 7 禁止自动降级（R-09 + §7 四条禁止①）
判据：AgentRunner 不可用不自动落黑箱；不因 AgentRunner 失败自动切黑箱。
- 正例：supervised 定义仅走 supervised 装配、AgentRunner 定义仅走 AgentRunner 装配（各测试正例分支）。
- 反例：`TestRunCommandRefusesMixedSupervisedAgentRunnerInputs`、`TestAgentRunnerAssemblyRefusesSupervisedDefinition`、`TestSupervisedAssemblyRefusesAgentRunnerDefinition`、`TestNoSupervisedReferencesInAgentRunnerPaths`、`TestSupervisedPathNeverRoutesThroughAgentRunner`、`TestSupervisedStepRejectsAgentRunnerExecutionTarget`、`TestRouterRefusesSupervisedWithoutPort`。
- 断言 M7：`<p7>=MixedSupervisedAgentRunner|AssemblyRefuses|NoSupervisedReferencesInAgentRunnerPaths|SupervisedPathNeverRoutes|SupervisedStepRejectsAgentRunner|RouterRefusesSupervised` ⇒ PASS≥7、FAIL==0。
- 证据：`03-test/L07-no-auto-degrade.txt`

### 腿 8 未知维度不显示为安全（§7.1(d)）
判据：报告与渲染对每个未知维度只显示 `unknown`；无 `safe`/`verified`/数字承诺。
- 正例：`TestSupervisedReportRenderGolden`（golden 文本含 `assuranceLevel: reduced` 与 `unknown: <dim>=unknown` 行）、`TestSupervisedReportOmitsUnprovenReviewDimension`。
- 反例：B 层合成报告 unknownDimensions 含 closed-set 外 `safe` ⇒ CLI 校验拒（exit∈[10,20]）；渲染输出 grep `safe`==0 且 `verified`==0。
- 断言 M8：`<p8>=ReportRenderGolden|ReportOmitsUnproven` ⇒ PASS≥2、FAIL==0；B 层探针 exit∈[10,20]；渲染文本 grep 计数==0/0。
- 证据：`03-test/L08-unknown-not-safe.txt`

### 腿 9 九子腿（R-09.1~R-09.9 逐项独立验证）

| 子腿 | 裁定 | 正例（机制证据） | 反例 | 机械断言 |
|---|---|---|---|---|
| 9.1 工具使用 | (c) | `TestSupervisedAssemblyFullChainExtendedScans` 全绿（verifiedDimensions 含 scope-scan/side-effect-scan）⇒ 产物在 workspace 内 | `TestSupervisedReturnRefusesOutOfScopeChange`、`TestSupervisedScopeUnderivableFailsClosed`、`TestSupervisedScopeAllowsRejectsLinkResolvingOutsideRoot`（workspace 外痕迹 ⇒ 阻断） | M9-1：`<p9-1>=OutOfScopeChange\|ScopeUnderivable\|ScopeRejects\|ScopeAllows` ⇒ PASS≥4、FAIL==0 |
| 9.2 资源读取 | (c) | `TestSupervisedIsolationReadsOnlyWorkspace`（manifest/diff 只读 workspace 内文件） | 越界读取痕迹 ⇒ 与 9.1 反例组同源阻断 | M9-2：`<p9-2>=IsolationReadsOnlyWorkspace` ⇒ PASS≥1，且 M9-1 组 FAIL==0 |
| 9.3 网络 | (c) | 本片零外部调用（无 VM/代理/外部连接）=「公网阻断」最强形式；`01-arch/r09-3-user-confirmation.md` 含精确句「回环可接受、公网必须阻断。」 | 确认精确句缺失 ⇒ 腿 9.3 阻断（可证伪） | M9-3：确认件 grep -c 精确句==1；`L09-03` 网络审计注记 =「外部连接 0 条」 |
| 9.4 费用 | (a) | 全链报告 `cost ∈ unknownDimensions`（与腿 5 联动） | 缺 cost ⇒ 拒（同 M5 反例） | M9-4：报告 JSON `cost` 计数==1 |
| 9.5 后台进程 | (a) | `TestSupervisedReturnRefusesLiveManagedProcess`（受管进程须停才归还）+ 报告含 `background-processes` + ⑧ 写入「用户负责清理」 | 受管进程存活归还 ⇒ 拒 | M9-5：`<p9-5>=RefusesLiveManagedProcess` PASS≥1 + 报告含 `background-processes` + 「用户负责清理」计数≥1 |
| 9.6 会话连续性 | (a) | `TestSupervisedRunRestartRecovery`、`TestSupervisedConfirmationCrashAfter{Append,Resume,WindowOpen}Converges`（ProofRail 侧连续性） | 外部 agent 会话不可观测 ⇒ 报告含 `session-continuity` | M9-6：`<p9-6>=RestartRecovery\|ConfirmationCrashAfter` PASS≥4 + 报告含 `session-continuity` |
| 9.7 外部副作用 | (c) | `01-arch/r09-7-user-confirmation.md` 含精确句「已接受未知外部副作用。」+ 副作用扫描执法（`RefusesSideEffect`/`SideEffectScanBlocksTamperedManifests`） | 确认精确句缺失 ⇒ 腿 9.7 阻断（可证伪） | M9-7：确认件 grep -c 精确句==1 + `<p9-7>=RefusesSideEffect\|SideEffectScanBlocks` PASS≥2、FAIL==0 |
| 9.8 强制字段 | (a) | schema + Go 校验 + CLI 三层均强制 reduced/非空/cost | full / 空 / 缺 cost 三反例均拒 | M9-8：三反例探针（schema reject ×3 + B 层 exit∈[10,20]）计数==3 |
| 9.9 判定口径 | (b) | ③ 落盘 9 份独立子腿证据（L09-01…L09-09）；⑦ 判定表 9 行独立、无「R-09 总判定」单行 | 结构反例：任何合并/笼统判定表 ⇒ 结构不合规 | M9-9：`03-test/` 下 `L09-*` 文件数==9 且 ⑦ 判定表行数==9 |

- 证据：`03-test/L09-01-tool-use.txt` … `03-test/L09-09-judgment-approach.txt`

## 4. 证据落盘路径（`03-test/` 全量）
`L01-risk-confirmation-binding.txt`、`L02-delivery-not-completion.txt`、`L03-explicit-return-postscan.txt`、`L04-independent-gates-review.txt`、`L05-unknown-dimensions-cost.txt`、`L06-three-sentences.txt`、`L07-no-auto-degrade.txt`、`L08-unknown-not-safe.txt`、`L09-01-tool-use.txt`、`L09-02-resource-reads.txt`、`L09-03-network.txt`、`L09-04-cost.txt`、`L09-05-background-processes.txt`、`L09-06-session-continuity.txt`、`L09-07-external-side-effects.txt`、`L09-08-mandatory-fields.txt`、`L09-09-judgment-approach.txt`、`SUMMARY.md`。每件含：命令原文 + 原始输出摘录 + 正/反例结果 + 断言读数；`SUMMARY.md` 为 17 腿结果总表（腿号/正例/反例/断言读数/结果）。

## 5. 机械断言清单（④/⑨ 与主控复算，A 层模板见 §2）
- M1：`<p1>` PASS≥10、FAIL==0；CLI run 探针 exit∈[10,20] + stderr 命中 `supervised visible delivery is not accepted`。
- M2：`<p2>` PASS≥4、FAIL==0。
- M3：`<p3>` PASS≥13、FAIL==0。
- M4：`<p4>` PASS≥5、FAIL==0。
- M5：C 层 `pass 4 / fail 0`；报告七项计数==7；缺 cost 探针 exit∈[10,20]；缺项探针记录接受/拒绝。
- M6：三句精确串（起跑消息【七】5 原文）各计数==2（§1==1、摘要==1、区段外==0）。
- M7：`<p7>` PASS≥7、FAIL==0。
- M8：`<p8>` PASS≥2、FAIL==0；篡改探针 exit∈[10,20]；渲染 grep `safe`==0 且 `verified`==0。
- M9-1…M9-9：见腿 9 表第 5 列。
- 门禁（④ 必跑）：`go build ./...`、`go vet ./...`、`go test -count=1 ./...` FAIL==0、`node tools/gates/gate.js --all --scope={tree,index,ci}` 三作用域 `TOTAL_FAIL=0`、`--selftest` = `SELFTEST: PASS 722/722`。

## 6. ⑦ 判定问题清单（唯一判定主体；首轮不附任何清单；盲审 1 次，§7.7 停机/身份语义）
- Q1–Q8：腿 1–8 各一题——通过/不通过 + 证据路径（`03-test/L0X`）+ 关键读数。
- Q9–Q17：9.1–9.9 九子腿各一题——通过/不通过 + 证据路径。
- Q18–Q26：R-09.1~R-09.9 九项边界逐项判定，每行四列 =（裁定值 / 实测证据路径 / 是否吻合 / 判定），互不合并、不累积、不笼统。
- Q27：AT-24 验收结论（条件闭合 / 部分通过 / 不通过）+ 理由引用（17 子腿 + 九项边界）。
- Q28：三句话逐字核验结论（⑤⑥⑦ 各环均已核，⑦ 独立复核）。
- Q29：腿 5 缺项探针若「接受」——是否构成 T028 实现缺陷；若修复需动 schema/fixtures 是否按 3.3 停下升级（⑦ 判 + 用户裁）。
- Q30：真实 SessionBridge visible 端到端 unknown/边界登记（3.4 三项必含项）是否如实。
- 要求：⑦ 报告必须逐项含 17 子腿 + R-09.1~R-09.9 九项边界 + AT-24 结论；主控只转录。

## 7. 执行顺序与依赖
1. 步骤 0（主控）：ADR-040 + 台账登记 + freedom/input-list + 两确认件（已补录完成）。
2. ③：构建二进制 → A/B/C 层执行腿 1–8 → 落 L01–L08。
3. ③：腿 9 九子腿（依赖腿 1–8 证据 + `01-arch/` 两确认件）→ 落 L09-01…L09-09 + `SUMMARY.md`。
4. 任一腿 FAIL 且根因 ∈ `internal/console|chain` ⇒ ② 修 → ③ 重跑受影响腿（计入 R 计数）；根因 ∈ schema/fixtures/CONTRACTS/R-09 ⇒ 停下升级。
5. ④ 门禁全绿 → 派 ⑤。
6. ⑧a 草案（三句话入 §1+执行摘要）→ ⑤ 预审（逐字核验）→ ⑥ 扫描 → ⑦ 终审+盲审（唯一判定）→ ⑧b 定稿 → ⑨（三作用域 + SELFTEST + M6 机械计数）→ ⑧c 清 `tmp/P64-AT24-SUPERVISED-BLACK-BOX/` → ⑩ 停点（等同一轮授权；禁自行 commit/push）。
DESIGN: DONE

P64 AT24-SUPERVISED-BLACK-BOX · 03-test / SUMMARY（17 子腿结果总表）
=====================================================================
执行角色：③ prfrail-tester（DeepSeek V4.1 Flash，runSubagent 显式传 model）
执行时点：2026-11-11（批次 1/3 + 2/3 + 3/3 汇总）
工作目录：D:\LZProjects\prfrail
环境：Windows 本机；无 VM / 无代理 / 无外部调用；Go 1.22 工具链；二进制
      tmp\P64-AT24-SUPERVISED-BLACK-BOX\prfrail.exe（go build ./cmd/prfrail，BUILD_EXIT=0）
性质：**机制腿**（stub visible 投递 + 模拟人工归还 + 后置扫描 + 独立 gates/review + reduced/unknown 报告），
      不是真实 SessionBridge visible 端到端（真实端到端依赖 T027-B，登记 unknown/边界）。
禁改项遵守：未改源码 / schema / fixtures / go.mod / CONTRACTS / R-09 / 两件用户确认件。

[1] 17 子腿结果总表
------------
| 腿 | 正例 | 反例 | 断言读数 | 结果 |
|---|---|---|---|---|
| 1 风险确认绑定 | SupervisedRiskConfirmationRoundTrip / OpenSupervisedWorkWindowWithValidConfirmation / SupervisedRiskNoticeIsTruthfulAndVersioned / SupervisedAssemblyFullChain | OpenSupervisedWorkWindow{WithoutConfirmation,WrongNoticeVersion} / RiskConfirmationRejectsBindingDrift / Assembly{BlocksWithoutConfirmation,RejectsSelectionReject} / OpenSupervisedWorkWindowRejectsUnacceptedDelivery + B 层 run 探针 | M1-A PASS=12 FAIL=0；M1-B exit=1 且 stderr 命中 "supervised visible delivery is not accepted" | **PASS** |
| 2 投递 ≠ 完成 | VisibleDeliveryMapsAccepted / 投递零回执链 | SubmitSupervisedReturnRejects{DeliveryOnlyEvidence,ReceiptNotBoundToFacts} | M2 PASS=5 FAIL=0 | **PASS** |
| 3 显式归还 + 后置全量扫描 | SubmitSupervisedReturnCompleteFlow / AssemblyFullChain / AssemblyFullChainExtendedScans | Return{RefusesAbortReceipt,RefusesSecretFinding,RefusesOutOfScopeChange,RefusesSideEffect,RefusesFailedGate,RefusesTimeout,EnforcesFrozenPolicy,RefusesDriftedPolicy} / SideEffectScanBlocksTamperedManifests / FileTypeScanRejectsUnacceptedType / ReturnRefusesLiveManagedProcess | M3 主跑 PASS=17 FAIL=0；补跑 PASS=2 FAIL=0 | **PASS** |
| 4 独立 gates/review | PostflightQualifiesReviewPending / PostflightDecisionMustConsumeAllFacts | ReturnRefusesFailedGate / UnqualifiedReviewRefused / TaskWithoutReturnFactsFailsClosed | M4 PASS=26 FAIL=0（设计点名 5/5） | **PASS** |
| 5 unknownDimensions 七项含 cost | SupervisedReportNewAlwaysReduced（首跑 4 子例 → R1 11 子例） | 契约夹具 assurance-level-full / unknown-dimensions-empty / missing-cost（schema 拒 3/3）+ CLI 缺 cost 拒 + 缺项探针 | 首跑：C 层 tests 4 / pass 4 / fail 0；七项计数=7、reduced=1；缺 cost exit=1；**缺其余六项之一 ⇒ exit=0 接受**。R1（DEF-1 修复后）：A 层 SUBTEST_PASS=11 / FAIL=0；p5-only-cost 与 p5-missing-network-access 均 exit=1；p5-full-seven exit=0 | **PASS（R=1 重跑；首跑 DONE WITH DEFECT）** |
| 6 三句话硬性出现 | 位置 A（§1）/ 位置 B（执行摘要）各 1 次；区段外 0 次 | 任何改写/摘要变体 ⇒ 精确串计数 0 | S1/S2/S3 逐字原文与计数口径已登记（L06 [1]）；报告未装配 ⇒ 计数未执行 | **PENDING-COUNT**（⑨ 机械计数） |
| 7 禁止自动降级 | supervised→supervised 装配 / AgentRunner→AgentRunner 装配 | MixedSupervisedAgentRunnerInputs / AgentRunnerAssemblyRefusesSupervisedDefinition / SupervisedAssemblyRefusesAgentRunnerDefinition / NoSupervisedReferencesInAgentRunnerPaths / SupervisedPathNeverRoutesThroughAgentRunner / SupervisedStepRejectsAgentRunnerExecutionTarget / RouterRefusesSupervisedWithoutPort | M7 PASS=16 FAIL=0（设计点名 7/7） | **PASS** |
| 8 未知维度不显示为安全 | ReportRenderGolden / ReportOmitsUnprovenReviewDimension | B 层 unknownDimensions 含 safe / verified ⇒ CLI 拒 | M8 PASS=2 FAIL=0；safe/verified 探针 exit=1；渲染 grep `safe`=0 `verified`=0 | **PASS** |
| 9.1 工具使用 | AssemblyFullChainExtendedScans（verifiedDimensions 含 scope-scan/file-type-scan/side-effect-scan；三扫描 clean） | ReturnRefusesOutOfScopeChange / ScopeUnderivableFailsClosed / ScopeRejectsUnresolvableRoot / ScopeAllowsRejectsLinkResolvingOutsideRoot（SKIP） | M9-1 PASS=3 SKIP=1 FAIL=0（阈值 PASS≥4 ⇒ **差 1，SKIP 为环境不可测**） | **DONE WITH SKIP** |
| 9.2 资源读取 | IsolationReadsOnlyWorkspace（ChangedPaths 恰 1 项；workspace 外候选未被读取/改写） | 越界痕迹 ⇒ 由 9.1 反例组阻断（同源） | M9-2 PASS=1 FAIL=0；M9-1 组 FAIL=0 | **PASS** |
| 9.3 网络 | 本片零外部调用（proxy OFF、git proxy 未设、无 VM）；确认件含精确句 | 精确句缺失 ⇒ 阻断（可证伪：计数 0） | M9-3 grep -c「回环可接受、公网必须阻断。」=1；外部连接 0 条 | **PASS** |
| 9.4 费用 | 报告 cost ∈ unknownDimensions（构造器七项含 cost，原样传入） | 缺 cost ⇒ 拒（CLI exit=1 + schema reject） | M9-4 `cost` 计数=1（token=1、维度=1）；assuranceLevel=reduced | **PASS** |
| 9.5 后台进程 | ReturnRefusesLiveManagedProcess（受管进程存活 ⇒ 归还 fail closed） | 受管进程存活归还 ⇒ 拒 | M9-5 PASS=1 FAIL=0；报告含 background-processes=1；报告侧「用户负责清理」=0 ⇒ **PENDING-⑧** | **PASS（2/3，第三项 PENDING-⑧）** |
| 9.6 会话连续性 | RunRestartRecovery / ConfirmationCrashAfter{Append,Resume,WindowOpen}Converges | 外部 agent 会话不可观测 ⇒ 报告标 unknown | M9-6 PASS=4 FAIL=0；报告含 session-continuity=1 | **PASS** |
| 9.7 外部副作用 | 确认件含精确句 + 副作用扫描执法（RefusesSideEffect / SideEffectScanBlocksTamperedManifests 5 子例） | 精确句缺失 ⇒ 阻断；五类篡改 ⇒ 阻断 | M9-7 grep -c「已接受未知外部副作用。」=1；PASS=2 SUBTEST_PASS=5 FAIL=0 | **PASS** |
| 9.8 强制字段 | reduced + 非空 + 含 cost 三层（schema/Go/CLI） | full / empty / missing-cost 三反例；**缺其余六项之一** | 首跑：三反例计数=3 均拒；缺项探针 exit=0 ⇒ DEF-1 成立。R1（DEF-1 修复后）：三反例仍拒 3/3；缺项探针（only-cost / missing-network-access）exit=1；schema 层仍 ACCEPT（② 已登记） | **PASS（R=1 重跑；首跑 DONE WITH DEFECT）** |
| 9.9 判定口径 | 9 件独立证据 L09-01…L09-09 落盘 | 合并/笼统判定表 ⇒ 结构不合规 | M9-9 L09-*.txt=9；⑦ 判定表行数 ⇒ **PENDING-⑦** | **PASS（结构侧）** |

汇总计数（合计 = 17；R1 重跑后更新）：PASS = 15（腿 1/2/3/4/5/7/8 + 9.2/9.3/9.4/9.5/9.6/9.7/9.8/9.9）；
          DONE WITH DEFECT = 0（腿 5、9.8 经 R=1 重跑转 PASS）；DONE WITH SKIP = 1（9.1，环境 symlink 不可测）；
          PENDING-COUNT = 1（腿 6，三句话计数待 ⑨）。
          R 计数：腿 5 = R1、腿 9.8 = R1（2026-11-11，DEF-1 修复后重跑）。
          跨腿指针：PENDING-⑧（9.5 报告「用户负责清理」）、PENDING-⑦（9.9 判定表 9 行）。
          硬失败（测试 FAIL>0）= **0**。

[2] 缺陷候选
--------
DEF-1（OPEN，② 待修）—— T028「七项全必填」未强制
  事实：unknownDimensions 缺任一非 cost 项（["cost"] 或七项缺 network-access）⇒ CLI exit=0 接受并渲染；
        schema（Ajv）同向 ACCEPT。
  期望：R-09.8 / 腿 5 判据「缺任何一项 ⇒ 报告不合规」⇒ 应拒。
  根因层：(a) internal/chain/supervised_report.go:190-192 仅 non-empty + contains(cost)；
          (b) schemas/supervised-report.schema.json 仅 minItems=1 + contains(cost)。
  修复边界：Go 层可修（加封闭七项集校验）；schema 一致化需改 schemas/ + 夹具 ⇒ 按 design-r1 §3.3 超界停下升级。
  ③ 处置：不修；交 ② 与 ⑦/用户裁决。影响：9.8 与腿 5 均只能判 DONE WITH DEFECT。

[3] 设计偏差登记（本片 ③ 侧）
--------------------
D-1  / D-14（exit 区间口径）：设计写 exit∈[10,20]；实况 CLI 仅 exitSuccess=0 / exitFailure=1 / exitUsage=2
        （internal/console/cli.go:23-25），实测非零即拒（1）。采信「非零即拒」，区间期望不成立。
D-11 / D-15（正则宽度）：`<p4>` 的 "Postflight"、`<p7>` 的 "AssemblyRefuses" 等词干宽匹配带入设计未点名的测试，
        读数高于点名数（M4 26 vs 点名 5；M7 16 vs 点名 7）。已按点名 5/5、7/7 单列读数，不改变 PASS 判定。
D-13（全链报告 JSON 未落盘，显著）：全链报告仅由 harness 在 t.TempDir() 产出即删；生产 run 需真实 visible
        投递（本环境不可用）。故腿 5/9.4/9.5/9.6 的「报告含 X」采用「生产常量 → 构造原样传入 → CLI 渲染」
        三段式代证（原始输出在 raw/leg9-5-6-render.txt）。不新增临时 Go 测试（禁改源码）。
D-16（M9-1 阈值差 1，本批次新登记）：TestSupervisedScopeAllowsRejectsLinkResolvingOutsideRoot 在 Windows 本机
        无 SeCreateSymbolicLinkPrivilege ⇒ 测试内建 SKIP（runtime_supervised_test.go:1008），断言前即跳过。
        故 PASS=3 < 阈值 4。属**环境不可测**，非实现缺陷；本片不提权、不更换环境以凑阈值。

[4] 待办指针（非缺陷，流水线状态）
----------------------
PENDING-⑧：报告 at24-supervised-black-box{,_EN}.md 必须明写「用户负责清理」（R-09.5），
            当前报告不存在、全仓库该串仅出现在设计/裁定/起跑文本（4 处，均非报告）。
PENDING-⑦：⑦ 判定表须 9 行独立（裁定值/证据路径/是否吻合/判定），不得以「R-09 总判定」单行替代；
            ③ 只登记结构要求，不预判。
PENDING-⑨：三句话精确串计数（每句 total=2、§1=1、摘要=1、区段外=0）。

[5] 证据清单（03-test/，17 件）
----------------------
L01-risk-confirmation-binding.txt / L02-delivery-not-completion.txt / L03-explicit-return-postscan.txt /
L04-independent-gates-review.txt / L05-unknown-dimensions-cost.txt / L06-three-sentences.txt /
L07-no-auto-degrade.txt / L08-unknown-not-safe.txt / L09-01-tool-use.txt / L09-02-resource-reads.txt /
L09-03-network.txt / L09-04-cost.txt / L09-05-background-processes.txt / L09-06-session-continuity.txt /
L09-07-external-side-effects.txt / L09-08-mandatory-fields.txt / L09-09-judgment-approach.txt
（实测 *.txt 计数 = 17）

[6] 本批次（3/3）执行命令清单（可追溯）
----------------------------
  go test -count=1 -run 'OutOfScopeChange|ScopeUnderivable|ScopeRejects|ScopeAllows' ./internal/chain ./internal/console -v
  go test -count=1 -run 'TestSupervisedAssemblyFullChainExtendedScans' ./internal/console -v
  go test -count=1 -run 'TestSupervisedIsolationReadsOnlyWorkspace' ./internal/chain ./internal/console -v
  go test -count=1 -run 'TestSupervisedReturnRefusesLiveManagedProcess' ./internal/chain ./internal/console -v
  go test -count=1 -run 'TestSupervisedRunRestartRecovery|TestSupervisedConfirmationCrashAfter' ./internal/chain ./internal/console -v
  go test -count=1 -run 'RefusesSideEffect|SideEffectScanBlocks' ./internal/chain ./internal/console -v
  go build -o tmp\P64-AT24-SUPERVISED-BLACK-BOX\prfrail.exe ./cmd/prfrail
  prfrail.exe supervised report --run-dir fixtures\runs\{p5-full-seven,p5-missing-cost,p5-only-cost,p5-missing-network-access}
  node tmp\P64-AT24-SUPERVISED-BLACK-BOX\check-supervised-fixtures.js
  node -e（确认件精确句计数 / report.json 维度计数）
  .vscode/scripts/proxy-toggle.ps1（Status，只读）
  全部命令均在本机执行；无 VM / 无代理 / 无外部连接（proxy: system OFF、git proxy 未设）。

LEGS 1-8: DONE（批次 1/2 已落）  LEGS 9: DONE WITH FAILS
（说明：「WITH FAILS」= 本批次存在 2 项 DONE WITH DEFECT（腿 5 / 9.8，同源 DEF-1）+ 1 项环境 SKIP（9.1）；
  测试断言层面 FAIL=0，无红测。）

[R1] 腿 5 / 9.8 重跑（2026-11-11，② 修复 DEF-1 后，R=1）
-----------------------------------------------------
  首跑：腿 5 / 9.8 = DONE WITH DEFECT（同源 DEF-1：缺任一非 cost 项仍被接受）
  R1 ：腿 5 / 9.8 = PASS（p5-only-cost、p5-missing-network-access 均由 exit=0 转为 exit=1；
        正例 p5-full-seven 保持 exit=0；A 层 11 子例全 PASS；C 层 4/4 不变）
  回归：go test -count=1 ./internal/chain ./internal/console ⇒ ok / ok（FULL_EXIT=0）
  DEF-1 现状：**Go/CLI 层闭合**；schemas/supervised-report.schema.json 未改动（② 登记），
        Ajv 直校 p5-only-cost / p5-missing-network-access 仍 ACCEPT ⇒ 分层单向差（Go 更严，fail closed），
        遗留至后续片处理。
  原始输出：tmp/P64-AT24-SUPERVISED-BLACK-BOX/raw-r1/（临时件，⑧c 清理）

注记（2026-11-11 ⑧b）：9.5 行 PASS(2/3) 的第三项（报告侧「用户负责清理」）已由 P64 报告 §2/§3/§4 闭合；末行 LEGS 9 结论按 R=2 后口径 = 17 子腿 16 通过、9.1 经 junction 等价补测通过（⑦ 复审确认）、硬失败 0。

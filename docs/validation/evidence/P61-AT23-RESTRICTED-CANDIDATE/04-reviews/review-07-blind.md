# ⑦ 盲审回执（independent-reviewer + GPT-5.3-Codex，只读，逐字转录）

基于读取主控落盘的指定证据（切片定义、契约原句、生产 diff 与生产脚本、门禁读数，以及 real-candidate/run-1~run-5 原始证据），我给出独立盲审结论如下。

**① 发现清单**
1. HIGH | [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L22](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L22)（并见 [#L5](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L5)） | 违反句逐字：「机制腿全绿 + 真实候选腿实际跑过（≥3 次，覆盖启动/执行/停机）并记录 unsupported/unknown。」（[docs/validation/r08-r09-ruling.md#L38](docs/validation/r08-r09-ruling.md#L38)，同口径见 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md#L88](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md#L88)） | 最小修复：在最终 AT-23 报告中不得宣称“前提②已满足”；要么把 F2 该条 failed check 清零后重跑并更新证据，要么以同轮用户裁定形成明确豁免文本并把结论降级为“条件未闭合/待裁决”。

**② 四项任务逐项判定**
1. 安全性（fail-closed 语义）：通过
证据显示默认路径会强制加 no-auto-update（[tools/agent-probe/CopilotCli.ps1#L161](tools/agent-probe/CopilotCli.ps1#L161)）；Unlock 受控集固定为 1.0.83/1.0.85（[tools/agent-probe/CopilotCli.ps1#L59](tools/agent-probe/CopilotCli.ps1#L59)），超集拒跑（[tools/agent-probe/CopilotCli.ps1#L222](tools/agent-probe/CopilotCli.ps1#L222)）；钉版失配拒跑（[tools/agent-probe/CopilotCli.ps1#L209](tools/agent-probe/CopilotCli.ps1#L209)）；阶段 1 抽样频率与 R-08.1 一致（[tools/agent-probe/CopilotCli.ps1#L199](tools/agent-probe/CopilotCli.ps1#L199)、[tools/agent-probe/CopilotCli.ps1#L201](tools/agent-probe/CopilotCli.ps1#L201)、[tools/agent-probe/CopilotCli.ps1#L202](tools/agent-probe/CopilotCli.ps1#L202)，对应契约 [docs/validation/r08-r09-ruling.md#L12](docs/validation/r08-r09-ruling.md#L12)）。

2. 架构一致性（新增收敛 + 既有维持口径）：通过（按本次生产变更口径）
本次生产 diff 仅新增封装文件（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt#L1)）。既有调用形态仍为原登记口径：
- Invoke-AgentProbe 使用 Start-Process 与 taskkill 进程树控制（[tools/agent-probe/Invoke-AgentProbe.ps1#L540](tools/agent-probe/Invoke-AgentProbe.ps1#L540)、[tools/agent-probe/Invoke-AgentProbe.ps1#L563](tools/agent-probe/Invoke-AgentProbe.ps1#L563)）
- Invoke-CandidateProbe 使用 Start-Process + Win32_Process 观测（[tools/agent-probe/Invoke-CandidateProbe.ps1#L259](tools/agent-probe/Invoke-CandidateProbe.ps1#L259)、[tools/agent-probe/Invoke-CandidateProbe.ps1#L267](tools/agent-probe/Invoke-CandidateProbe.ps1#L267)）
- AppContainer 按名经 PATH 启动 copilot（[tools/agent-probe/appcontainer-b2/container-candidate.ps1#L231](tools/agent-probe/appcontainer-b2/container-candidate.ps1#L231)、[tools/agent-probe/appcontainer-b2/container-candidate.ps1#L250](tools/agent-probe/appcontainer-b2/container-candidate.ps1#L250)）
- Go 侧探针仍含 no-auto-update 与 deny-tool=shell（[internal/adapters/ai_probe_copilot.go#L122](internal/adapters/ai_probe_copilot.go#L122)、[internal/adapters/ai_probe_copilot.go#L130](internal/adapters/ai_probe_copilot.go#L130)）。

3. 完整性（版本钉住 + wrapper-probe-count 自洽）：通过
计数与证据链闭环存在：计数文件定义与递增写回（[tools/agent-probe/CopilotCli.ps1#L72](tools/agent-probe/CopilotCli.ps1#L72)、[tools/agent-probe/CopilotCli.ps1#L79](tools/agent-probe/CopilotCli.ps1#L79)、[tools/agent-probe/CopilotCli.ps1#L90](tools/agent-probe/CopilotCli.ps1#L90)、[tools/agent-probe/CopilotCli.ps1#L91](tools/agent-probe/CopilotCli.ps1#L91)），并在层 1 实测中体现 probe=10 与 value=11（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-01-wrapper-version.txt#L9](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-01-wrapper-version.txt#L9)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-03-unlock-version.txt#L15](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-03-unlock-version.txt#L15)）。
门禁读数为全绿：tree/index/ci 均 TOTAL_FAIL=0（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt#L47](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt#L47)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-index.txt#L47](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-index.txt#L47)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-ci.txt#L49](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-ci.txt#L49)），selftest 为 722/722（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-selftest.txt#L727](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-selftest.txt#L727)）。

4. 测试反例审查（run-1~run-5 对条件闭合前提②支撑度）：部分成立，但不足以判定“前提②已满足”
- run-1：执行成功（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/command.txt#L9](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/command.txt#L9)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stdout.txt#L15](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stdout.txt#L15)）。
- run-2：启动阶段失败且零模型调用（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/command.txt#L13](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/command.txt#L13)、[#L15](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/command.txt#L15)），并有版本钉死拒绝原文（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stderr.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stderr.txt#L1)）。
- run-3：停机腿证据存在（taskkill exit 0，failure stage=stop）（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-evidence.txt#L18](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-evidence.txt#L18)、[#L26](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-evidence.txt#L26)）。
- run-4：归因限制处理正确，未被当作 unsupported（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/verdict.txt#L7](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/verdict.txt#L7)、[#L12](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/verdict.txt#L12)；对应 argv 窄 allowlist 见 [run-4/argv.txt#L8](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/argv.txt#L8)）。
- run-5：unsupported 证据可信，且与 B1 #8 同形态（宽 allowlist + deny）（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L1)、[#L8](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L8)、[#L10](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L10)，argv 见 [run-5/argv.txt#L8](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/argv.txt#L8)）。
- 但机制腿并非全绿（failed checks=1），因此不能据此宣称前提②已闭合（[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L22](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L22)）。

**③ 什么会推翻我的判定**
1. 如果后续落盘证据把 layer 3 的 failed checks 从 1 变为 0，并保持 run-4 降级与 run-5 unsupported 的归因链不变，我的 HIGH 发现可被关闭。
2. 如果出现新证据证明 run-5 并非 deny 旁路（例如 deny 在同形态下实际生效且无执行），我对“run-5 可作为 unsupported 证据”的判断会被推翻。
3. 如果新增生产 diff 显示绕过封装的新增 copilot 调用点，我对“架构一致性通过”的判断会被推翻。
4. 如果契约文本对“机制腿全绿”的定义被正式修订为“允许该类继承漂移不计入红项”，且修订先于本片结论生效，我的 HIGH 发现将降级或消失。

BLIND REVIEW: FINDINGS

## 主控注记（非 ⑦ 盲审原文）
- 首轮盲审未附任何清单（§7.7 必抽，触发面=停机/身份语义），独立给出 1 条 HIGH（机制腿 failed checks=1 ⇒ 前提②不可宣称已闭合）。
- HIGH 与 FINDING-1 同源：F2 `process-tree-taken-away` 继承漂移（02d5797）。用户已裁定 FINDING-1=(b)（登记继承漂移、不修 harness、交 ⑦ 判前提满足性）；盲审的最小修复给出两案：清零重跑 或 用户同轮裁定豁免文本 + 结论降级为「条件未闭合/待裁决」。交 ⑦ 终审与用户裁定。
- 本回执由主控逐字转录落盘（2026-11-08，⑦ 盲审完成后）。

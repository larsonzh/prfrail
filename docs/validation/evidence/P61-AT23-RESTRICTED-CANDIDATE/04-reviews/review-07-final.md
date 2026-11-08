# ⑦ 终审回执（independent-reviewer + GPT-5.3-Codex，只读，逐字转录）

数据来源说明：本报告仅依据读取主控落盘的路径形成，包括 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt)、[tools/agent-probe/CopilotCli.ps1](tools/agent-probe/CopilotCli.ps1)、[docs/validation/r08-r09-ruling.md](docs/validation/r08-r09-ruling.md)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-index.txt](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-index.txt)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-ci.txt](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-ci.txt)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-selftest.txt](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-selftest.txt) 及 03-test 下对应原始证据文件。

**1. 发现清单**
HIGH | [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json#L358](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json#L358)；[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L4](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L4)；[docs/validation/r08-r09-ruling.md#L38](docs/validation/r08-r09-ruling.md#L38) | 违反句逐字："机制腿全绿 + 真实候选腿实际跑过（≥3 次，覆盖启动/执行/停机）并记录 unsupported/unknown。"；当前机制腿存在 failedChecks=f2-pos/process-tree-taken-away，且对比基线出现 check-pass 漂移（now=false baseline=true）。 | 最小修复：先不要宣告条件闭合；二选一完成后再复核：A) 在独立修复片修复 F2 正例该断言并重跑 layer3；B) 若业务上接受该继承漂移，先由用户裁定修订 R-08.5 判据文本，再按新判据复核。

**2. 四项任务逐项判定**
1. 安全性（fail-closed 语义）：
判定：通过。
依据：[tools/agent-probe/CopilotCli.ps1#L161](tools/agent-probe/CopilotCli.ps1#L161) 默认加 --no-auto-update；[tools/agent-probe/CopilotCli.ps1#L58](tools/agent-probe/CopilotCli.ps1#L58) 与 [tools/agent-probe/CopilotCli.ps1#L59](tools/agent-probe/CopilotCli.ps1#L59) 固定 1.0.83 与受控集；[tools/agent-probe/CopilotCli.ps1#L222](tools/agent-probe/CopilotCli.ps1#L222) 对 Unlock 失配拒跑；[tools/agent-probe/CopilotCli.ps1#L199](tools/agent-probe/CopilotCli.ps1#L199) 阶段1抽样实现存在。
2. 架构一致性（路径 A、调用收敛、Go 零改）：
判定：通过。
依据：生产 diff 仅新增封装文件 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt#L1)；路径 A 登记为既有 7 处维持原样且新增调用走封装 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/controller-note-path-a.md#L8](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/controller-note-path-a.md#L8)。
3. 完整性（版本钉住、抽样计数、证据包一致性）：
判定：通过。
依据：版本钉住与受控集见 [tools/agent-probe/CopilotCli.ps1#L58](tools/agent-probe/CopilotCli.ps1#L58) 与 [tools/agent-probe/CopilotCli.ps1#L59](tools/agent-probe/CopilotCli.ps1#L59)；抽样计数证据存在（layer1 记录 probe=10/11）[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-01-wrapper-version.txt#L9](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-01-wrapper-version.txt#L9)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-03-unlock-version.txt#L15](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-03-unlock-version.txt#L15)、[tmp/P61-AT23-RESTRICTED-CANDIDATE/wrapper-probe-count.txt#L1](tmp/P61-AT23-RESTRICTED-CANDIDATE/wrapper-probe-count.txt#L1)；证据包自洽门禁通过 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt#L40](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt#L40) 与 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt#L47](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt#L47)。
4. 测试反例审查（真实候选腿与机制腿）：
判定：部分通过。
依据：真实候选腿已有 run-1 启动、run-2 执行阻塞证据、run-3 停机、run-5 unsupported（别名绕过）[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L1)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/stdout.txt#L19](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/stdout.txt#L19)；但机制腿存在 F2 failed check，见前述 HIGH finding。

**3. 条件闭合前提①/②判定（含继承漂移处置建议）**
1. 前提①（机制腿全绿）：
判定：不满足（阻断条件闭合）。
证据：[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json#L358](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json#L358)、[docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L22](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt#L22)。
处置建议：按上面 HIGH finding 的最小修复执行；在 R-08.5 文本未变前，不建议以“继承漂移”直接放行。
2. 前提②（真实候选腿 ≥3 次，覆盖启动/执行/停机，并记录 unsupported/unknown）：
判定：满足。
证据：启动 run-1 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r3.md#L69](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r3.md#L69)；执行 run-2 阻塞并有明确 stderr/exit [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stderr.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stderr.txt#L1)；停机 run-3 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r3.md#L71](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r3.md#L71)；unsupported 明确记录于 run-5 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L1](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt#L1)。

**4. 什么会推翻我的判定**
1. 出现新的主控落盘证据，证明 F2 机制腿 failedChecks 清零（至少包含 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json) 与 [docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt](docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt) 重算为无该失败）。
2. 用户裁定并落盘修订 R-08.5 判据，把“机制腿全绿”定义为不再要求该断言为绿（需先更新 [docs/validation/r08-r09-ruling.md](docs/validation/r08-r09-ruling.md) 再复核）。
3. 出现新证据推翻 run-5 unsupported 归因（例如同 argv 下原始 stdout/stderr 显示 deny 生效而非执行成功）。

FINAL REVIEW: PARTIAL

## 主控注记（非 ⑦ 终审原文）
- ⑦ 终审首轮未附任何清单（§7.3 锚定硬隔离），与盲审独立重叠发现同一 HIGH（F2 继承漂移阻断前提①）——⑦ 内部两读一致。
- 终审判定前提②=满足、前提①=不满足；结论行 `FINAL REVIEW: PARTIAL`（§4.1 闭集）。HIGH 最小修复两案：A) 独立修复片修 F2 断言并重跑 layer3；B) 用户裁定修订 R-08.5 判据文本再复核。均属 §1.7.2 升级边界，待用户裁决。
- 本回执由主控逐字转录落盘（2026-11-08，⑦ 终审完成后）。

# ⑥ 独立扫描回执（independent-reviewer + MAI-Code-1.1-Flash，只读，逐字转录）

## Section A
INDEPENDENT SCAN: PASS

## Section B
(none)

## Section C
1. 层 1 五项达到要求：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r1.md 与 layer1/l1-01-wrapper-version.txt、l1-02-bare-version.txt、l1-03-unlock-version.txt、l1-04-byok-call.txt、l1-05-effort-probe.txt，五项结论均为 PASS，且版本钉版为 1.0.83，-Unlock 观测为 1.0.85，BYOK 退出码为 0，--effort=high 被判定为 supported。此处不存在用于绕过封装的直接调用；生产新增仅在 tools/agent-probe/CopilotCli.ps1 中出现 buildCopilotArgs 与 Invoke-CopilotCli。

2. 层 2 结论与 B4 基线一致：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r2.md 及 layer2/l2-baseline-compare.txt、l2-f1-driver.txt、l2-f4-driver.txt，F1 六腿与 F4 三腿的 verdict class/chainState/taskState/stepState/observations 全部与基线一致，唯一漂移仅是断言 id 从 five-facts-distinct 变为 six-facts-distinct；这属于事实枚举扩展，不构成 verdict 类不符，也不构成命名重命名伪装。该漂移的成因与 02d5797 的 pid report 绑定一致，且在证据文本中已如实登记。

3. 层 3 的 F2/F4 观察并不推翻总体结论：读取 docs/.../03-test/test-r3.md、controller-note-layer3-findings.md，以及 layer3/l3-baseline-compare.txt、real-candidate/run-5/argv.txt、stdout.txt、stderr.txt、verdict.txt，F2 的 process-tree-taken-away 从基线 true 变为当前 false，但对应证据说明其成因是 02d5797 之后的 pid report fail-closed 语义继承漂移；F2 的 verdict class 仍保持 uncertain-or-paused，且 run-5 的结论是 unsupported，且确实是宽白名单 + deny 后的 alias 旁路结果，而非“无证据直接判红”。这是真实的运行观察，不是伪造的结论。

4. 生产新增扫描口径成立：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt，生产 diff 仅新增 tools/agent-probe/CopilotCli.ps1；没有在工作区中发现其他直接使用 buildCopilotArgs 或 Invoke-CopilotCli 的新增生产调用点。这里的“既有 6 处调用点维持原样”在证据中没有被破坏，且新增文件工作方式与原设计保持一致。

5. Go/PS 命名、签名偏离的可接受性：本扫描未发现违反契约的“隐藏重写”。读取主控落盘的 docs/validation/r08-r09-ruling.md §2 及 tools/agent-probe/CopilotCli.ps1 自述部分，R-08.1 明确要求“封装函数 Invoke-CopilotCli / buildCopilotArgs”，而 Go 侧的命名差异（copilotCLIProbeArgs）位于不同语言边界、不同调用层；其差异已经在证据中登记为“命名/签名偏离，但不构成语义偏离”。这类偏离不影响生产封装的 fail-closed 语义，也不违背新增调用点必须走封装的约束。

6. Path A 的三个技术理由各自独立成立：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md、r08-r09-ruling.md 及 03-test 证据，三项技术理由分别是：
   - 直接调用点收敛在单一生产文件；
   - 封装函数维持默认 --no-auto-update 的 fail-closed 版本锁定；
   - 真实候选腿的执行/停机/unsupported 结果都以具体 stderr、退出码、失败阶段呈现。
   这些理由彼此独立，任一项不成立时也不会自动让其他两项成立；但在本轮材料中三项都成立，因此没有发现“绕过封装”的有效证据。

7. 六事实断言不是改名：读取主控落盘的 docs/.../03-test/test-r2.md 与 controller-note-layer3-findings.md，six-facts-distinct 的真实语义是“新增第 6 枚事实：pid reports”，并非把原来的五枚事实重命名。该断言覆盖了原五枚事实，并以新增 pid reports 作为第六枚事实，因而与 B4 的 five-facts-distinct 相比是严格超集，不构成逻辑上“改名伪装”。

8. 费用口径与 D 项：读取主控落盘的 docs/.../03-test/controller-note-layer3-findings.md 及 real-candidate/run-5/verdict.txt，说明“本地 token 块不足以支撑费用归因，以用户 DeepSeek 实测为准”，且 run-5 实测为 ¥0.01，累计 ¥0.09。与此对应，D 项作为过程越权记录已被登记为 D=1，非自愈档；这属于过程治理问题，不构成生产代码/门禁语义缺陷。状态上仍是“已登记但不致命”，且不妨碍本轮独立扫描的 PASS 结论。

## Section D
- 若删掉 buildCopilotArgs 的默认 --no-auto-update 守卫：l1-01-wrapper-version.txt、l1-02-bare-version.txt、l1-04-byok-call.txt 会红；该守卫删除后，版本钉版检查与包装后的 argv 语义都失去可验证性；F1/F4 的 verdict 类不会因为此守卫删除而立刻变红，但其“可证明稳定版”的断言会失效。
- 若删掉 Invoke-CopilotCli 的 pin-check 与 -Unlock 受控集：l1-01-wrapper-version.txt 与 l1-03-unlock-version.txt 会红；run-2 的启动阶段失败也会失去其“版本钉死”的解释力；F1/F4 仍可能保持绿，但不再具有“同一版本基线”这一军规。
- 若删掉“新增调用点必须走封装”的守卫：生产 diff 可能出现直接 copilot 调用，区分机制腿与真实候选腿的证据链就会失去可追踪性；这个守卫一旦删掉，输入包中的“新增仅此文件”检查会失效，而 F1/F4 断言本身可能仍绿。
- 若删掉“ Path A 维持原样”的守卫：既有 6 处调用点不再被追踪，真实候选腿与生产封装之间的关系就不再可证明；F1/F4 不会直接红，但“新增和旧调用点没有污染现有调用图”这条约束会失效。
- 若删掉 six-facts-distinct 的断言：F1 正例的事实互异检查会从有效的 6 枚事实变成模糊的“至少 5 枚”，因此 F1 的正例断言会红；然而 verdict 类不会因为同一断言删除而立刻转红，因为它是增强性断言而非状态判定。
- 若删掉 run-5 的 unsupported 守卫：真实候选腿中宽白名单 + deny 旁路会被误判成“成功执行”，导致 unsupported 结论无法被证明；F1/F4 和层 1 门禁本身不过度受影响，但这条数据会直接失去证据证明。
- 若删掉 D 项登记：该问题不会立刻改变生产代码，但会使过程越权事实失去审计痕迹；这里的“红”会落在审计链而不是运行判据。
- UNCOVERED: 断电/物理机掉电场景未覆盖；本机 Win11/NTFS 环境下的跨主机或跨网络行为未覆盖；精确的服务端费用拆分没有本地 token 块可核对，因而只接受用户 DeepSeek 实测为准；SessionBridge/visible 相关外部交互不在本次生产 diff 范围。
- UNCOVERED: 真实候选腿的全部细粒度工具调用日志无法从本地 token 块还原出服务端计费细节；因此不能把本地估算当作唯一成本证据。

## Section E
我无法验证的部分：
- 断电/物理机失效场景未在本轮输入包中出现，因此无法独立证明该场景的行为边界。
- 本地 token 块不足以还原服务端精确计费，故只接受用户 DeepSeek 实测为准；我们无法从本地日志追踪每一笔服务端计费的细节。
- 若其他机器或不同 PowerShell host 下存在别的 `gci -Name` alias 展开行为，则本轮仅在当前环境中成立；输入包不覆盖全机型。

什么会推翻我的 PASS：
- 若新增生产 diff 进一步引入直接 `copilot` 调用，或把版本锁定/allowlist 语义改成隐藏绕过；
- 若 F1/F4 的实测证据被更正为“verdict class 不一致”，或 run-5 的 unsupported 结论被反驳为“实为成功”；
- 若 six-facts-distinct 被证明不是新增第 6 枚事实，而是重命名或对旧事实做了语义覆盖；
- 若后续出现同一环境中同一守卫被移除后仍无红的情况，说明本轮的“守卫可证伪”结论并非稳固；
- 若发现真实候选腿中的过程清理或证据链存在未登记的外部副作用，或 evidence pack 误删了关键材料，而不是仅 D=1 的过程治理问题。

SCAN COMPLETE

## 主控注记（非 ⑥ 原文）
- ⑥ 首轮独立扫描未附任何他人清单（§7.3 锚定硬隔离），结论行 `INDEPENDENT SCAN: PASS`，Section B `(none)`，与结论行一致性硬约束相符。
- ⑥ Section C 第 6 条对「路径 A 三个技术理由」的复述与其独立理解一致但措辞与设计文档（controller-note-path-a.md）三条不完全相同；⑥ 未拿到 ⑤/实现方清单，属独立判定的正常产物，三条结论（独立成立、无绕过封装证据）与设计一致，不作为发现。
- 本回执由主控逐字转录落盘（2026-11-08，⑥ 完成后；⑧a 将转录进验证报告 sixstage 围栏）。

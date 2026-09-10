# T026 AgentRunner 契约与能力验证报告

[English](t026-agent-runner-contract_EN.md)

日期：2026-09-10；2026-09-11 更新探针器审计结果。结论：`CONTRACT COMPLETE / RUNTIME PROBE PARTIAL / AT-23 NOT PASS`。

## 已完成范围

1. 四份 Schema 冻结 request、capability、event 和 completion wire；契约门禁现覆盖 34 份 Schema、14 个 fixture 文件承载的 96 个用例及 2 条 canonical 向量。
2. `internal/adapters/agent_runner*.go` 实现 RFC 8785 域分隔摘要、严格解码、排序集合、create/resume 约束、request/event/completion 绑定，以及可从持久记录重建的幂等/冲突索引。
3. capability 固定 12 项矩阵；仅全项带运行证据且为 verified 才可 compatible。event 同时锁定 eventId、session/sequence 和 request/session，禁止同一 request 漂移到不同 session。completion 同时锁定 completionId 与 requestId，且 completed 只证明外部执行结束、进程树停止、日志/用量完整和输出 manifest 捕获，不代表 task PASS。固定 capability 测试还会读取实际证据文件并核对其字节 SHA-256，防止记录引用不存在或过时的证据。

## 当前结论与更正（2026-09-11 03:26）

- 撤回旧 P8 的 `toolControl=unsupported` 推论：固定摘要的 CLI `help permissions` 区分 Windows 可见工具 `powershell` 与权限类型 `shell(command)`，前缀通配格式为 `:*`。旧规则 `powershell(Get-ChildItem*)` 不能证明正确的 shell 拒绝策略被绕过。`help config` 还明确默认 manual 模式自动批准只读请求；旧 P9 的 `Get-Date` 不是可靠的必须确认挑战。以下旧 P8/P9 段落保留当时观察和推论作审计历史，当前结论以本节及更新后的机器记录为准。
- 本轮唯一授权调用保持原 P8 提示词和 `--available-tools=powershell`，改为 `--allow-tool=shell(Get-Location)`、`--deny-tool=shell(Get-ChildItem)`。session `29b41e1c-9a7d-42f0-b997-8f7678a964d2`，模型 `gpt-5.6-luna`，1 user request、3 内部模型回合、1 premium request、17587 输入 token、251 输出 token、6760 ms API 时长。无重试，无网络失败或代理切换。
- 272 条 UTF-8 JSONL 全部有效。`Get-Location` 成功；`Get-ChildItem -Name` 的请求、开始、失败结果关联同一 toolCallId `call_n7sNEAqvoFhyamIr9m4VgWx6`，结果为 `success=false`、`error.code=denied`，消息指明 `shell(Get-ChildItem)`。这是 CLI 结构化拒绝，不是模型自述或普通执行错误。新增解析逻辑只读重放返回该独立命令场景 supported；原始报告的 inconclusive 保留不覆盖。
- OS exitCode 与 CLI result exitCode 均为 0，未超时，workspace 文件与匹配残留进程均为空。通过启动后立即读取 `$proc.Handle` 保留句柄修复采集；本地同一启动函数在修复前复现 null，修复后准确采集 `cmd.exe` 的 0、7。历史 meta 不改写。离线回归包含 9 个原判定、2 个结构化拒绝正反例、权限规则检查及 2 个 OS 退出码检查。
- 原始摘要：stdout.jsonl 107389 字节，`be0cacb046ab9a78c6dc5bd3dac58c320c48fe651037d12589da35fbf1620f10`；usage.json 2115 字节，`56bcbf7587901aca64fa4a3eaf08689b930f0689705461a4fa736cd8133bb45c`；meta.json 1623 字节，`610b300aadbb01f063fef64c2a8d972cb5f3ff2bf1104c707c1495acdbf7b230`；原 probe-report.json 17188 字节，`eab66c197a31385157e14ca9b01fac8543ef45ec9f16c3447c805950c825e9a2`。均为 SHA-256。
- 新机器记录 `probe-copilot-cli-windows-20260911-corrected` 为 **6 verified、6 unknown、blocked**。独立命令拒绝不足以验证复合/间接命令、完整权限和 OS 隔离；`toolControl` 从错误 unsupported 撤回到 unknown，不提升 verified。其余 unknown 为 sessionResume、cancellation、networkControl、permissionControl、unattendedConfirmations。仍缺有效配置固定、复合/间接命令反例、真正必须确认场景、取消停树和会话恢复证据。T026 尚未完成，T027 不得启动。

## 候选探针（历史记录）

- 候选：GitHub Copilot CLI 1.0.83（Windows x64，构建提交 `e6a98f1`）。原生 `copilot.exe` 为 144,796,448 字节，SHA-256 为 `d3f3bb7b8bbf68357ad29f514a179d09f76135483d8bfb643131b8600f671ee2`；平台包元数据摘要为 `d6d07e0943462974bba6380170ca6c6d147e6e609b2cc5ceff165f8f16faa5ed`。
- 绕过 VS Code bootstrap 和 npm shim 直接执行原生二进制，`--version` 返回 1.0.83。`--help` 仅用于设计探针，不作为能力通过证据。
- 用户已明确授权一次最小、可能计费的模型调用。首个启动脚本因 Windows PowerShell 5.1 不支持 `Get-Date -AsUTC` 而在 Copilot 启动前失败，不计为授权调用；修正后的调用使用隔离 workspace、JSONL stream、30 AI credits 软上限、仅 `shell` 工具、精确 `shell(pwd)` 自动批准、HTTP/HTTPS deny、禁用 built-in MCP、系统临时目录、自定义指令、远程控制/导出及自动更新。
- 修正后的 CLI 进程未产生 JSONL、stderr 或 `usage.json`。468 字节的专用进程日志（SHA-256 `7d1b6cd0a3f696a9233b51071fff75c94542e33b2a0c1c697ef2008841ac5129`）显示认证编排拒绝了当前 GitHub 凭据，原因为 classic PAT 不受支持；日志不包含凭据值。进程在该错误后被中断，事后无匹配探针进程。没有模型响应或用量证据，现有产物不能证明是否到达模型端或计费，因此没有在“一次调用”授权下重试。
- 用户再次授权后，第二次预检仅在 Copilot 子进程中移除 `GH_TOKEN`，不修改用户环境或配置，以尝试回落到本地 OAuth 凭据。CLI 明确返回“未找到认证信息”并以 1 退出；381 字节的 `usage.json`（SHA-256 `f9f90cadcbcb1d9ba7176d520df7585abc79ff3ef751937204c366c6765b0c98`）记录 0 次用户请求、0 输入/输出 token、0 API 时长、0 premium request cost 且无 model metrics。因此第二次授权也没有发生模型调用或计费。
- 第二次预检显式使用 `--model auto`。Copilot CLI 不继承 VS Code 聊天输入框选择的模型；可用 `--model <name>` 为单次运行指定 CLI 支持的模型，或用 `--model auto` 让 CLI 自主路由。未显式传参时使用 CLI 自身的 `model` 配置或默认值。
- 随后在仅对登录子进程移除 `GH_TOKEN` 的条件下执行官方 `copilot login --web-flow`，CLI 报告账号 `larsonzh` 登录成功。登录后 `%USERPROFILE%/.copilot/config.json` 更新为 318 字节，Windows Credential Manager 出现非秘密 target `https://github.com:larsonzh.copilot-cli`；未读取或记录任何凭据值。
- 获得新的单次授权后，OAuth 模型探针成功创建 session `41ba03b7-6c3f-4f66-8a6b-8622c6f576ba` 并以 0 退出；`auto` 实际路由到 `mai-code-1.1-flash`，消耗 1 次请求、2138 输入 token、43 输出 token、5418 ms API 时长和 1 次 premium request。隔离 workspace 保持为空，退出后无匹配残留进程。
- stdout 因 PowerShell 5.1 重定向为 UTF-16LE；39/39 条非空 JSONL 均可解析，覆盖 session 创建、模型路由、模型调用完成、usage checkpoint、idle 与 result。原生流不含 ProofRail sequence 字段，adapter 后续必须分配单调序号。
- CLI 报告 `shell` 是未知 allowlist 工具，最终向模型暴露的工具数为 0；模型没有产生真实 tool-call/tool-result，最终文本中的伪工具调用不作为证据。因此 `noninteractive`、`eventStreamOrCompleteLogs`、`sessionCreate`、正常退出的 `processTreeStop` 和 `usage` 共 5 项 verified；`cwd`、`sessionResume`、`cancellation`、`toolControl`、`networkControl`、`permissionControl`、`unattendedConfirmations` 共 7 项 unknown，整体 blocked。机器证据及带 hash 的 capability record 位于 `testdata/agent-runner/capability-probes/`。
- 新授权的最小工具探针使用 1.0.83 实际注册的 `powershell` 工具。两次前置启动分别在 CLI 启动前被终端中断、在参数解析阶段因不支持 `--no-remote-control` 被拒绝，均无 session/model/usage，不计模型调用；修正后的唯一实际调用使用 `--no-remote`/`--no-remote-export`，创建 session `59e7d9a8-190d-4d86-b3a9-03b261b296e6`，`auto` 路由到 `gpt-5.6-luna`，消耗 1 次请求、11580 输入 token、102 输出 token、27720 ms API 时长和 1 次 premium request。
- 98 条 UTF-16LE JSONL 中，`assistant.message.toolRequests`、`tool.execution_start/partial_result/complete` 以同一 toolCallId 证明真实执行 `powershell(Get-Location)`，exit 0，并返回精确隔离 workspace。无 permission 事件或提示，workspace 为空、0 文件修改且无残留进程。该允许路径足以将 `cwd` 提升为 verified，但不证明 deny 路径、权限控制或需要确认时的无人值守行为；该次探针结束时为 6/12 项 verified、6 项 unknown，整体 blocked。
- 2026-09-11 获得一次 `tool-deny` 场景授权。探针器先因 Windows PowerShell 5.1 的泛型列表不支持 `.ToArray()` 而在 CLI 启动前失败；该目录没有 inv1/session/usage，未发生模型调用。修正参数数组并通过 parser、ScriptAnalyzer 与 dry-run 后，在同一未消耗授权下执行唯一实际调用，创建 session `9e01a613-d340-4eb1-a065-57bacfa848d0`，`auto` 路由到 `mai-code-1.1-flash`，消耗 1 次 premium request、7694 输入 token（含 3712 cache-read）、105 输出 token 和 3347 ms API 时长。
- 模型把两条指令合并成 `powershell("Get-Location; Get-ChildItem -Name")`；尽管参数包含 `--deny-tool=powershell(Get-ChildItem*)`，复合命令仍以 0 成功执行。139/139 条 UTF-8 JSONL 可解析，最终 result exit 0；workspace 为空、usage 报告 0 文件修改且无残留进程。这是可区分的绕过证据，故 `toolControl` 从 unknown 改为 unsupported。探针报告附带的“未产生干净退出”来自 PowerShell 5.1 进程对象未刷新而导致 meta exitCode 为 null；原生 result 明确为 0，采集器随后已修正，该次要缺陷不改变工具控制失败结论。
- 当前矩阵为 6 verified、1 unsupported、5 unknown：`sessionResume`、`cancellation`、`networkControl`、`permissionControl`、`unattendedConfirmations` 尚待探测。固定候选整体 blocked；即使余下 5 项全部 verified，也不能使 1.0.83 compatible。

## P9 追加观察（2026-09-11，待配置复核）

- 本轮授权已用于唯一一次 `permission-failclosed` CLI 调用，无重试。session 为 `9f219f5d-e9a6-4e2a-a1f4-9e3b779ba59e`，模型为 `gpt-5.6-luna`；usage 为 1 user request、3 个内部模型回合、1 premium request、17598 输入 token、205 输出 token、6639 ms API 时长。
- 参数仅批准 `powershell(Get-Location)`，但独立 toolCallId 分别证明 `Get-Location` 与 `Get-Date` 执行成功。196 条 UTF-8 JSONL 全部可解析，最终 result exit 0；报告中 workspace 文件及匹配残留进程均为空。调用正常到达模型，无网络屏蔽证据，因此未启用代理或再次调用。后续若确认网络屏蔽，可使用操作员指定的 `http://10.0.0.246:8080`，但不把网络失败视为策略拒绝，也不在费用不明时自动重试。
- 该观察证明当前有效配置没有实现探针预期的命令白名单；尚未独立固定默认授权、用户配置与 CLI 参数匹配语义，不能据此直接宣称厂商权限实现缺陷。P9 从“未执行”转为“已执行、待有效配置复核”；现有机器矩阵暂保留 6 verified、1 unsupported、5 unknown，`permissionControl` 和 `unattendedConfirmations` 不提升。T026/T027 继续阻断。
- meta exitCode 再次为 null，说明此前 `Refresh()` 修正未经行为验证且未解决问题；原生 result 不能代替 OS 退出码。保持这两类证据分离，后续需用非模型本地进程验证采集修复。
- 原始摘要：stdout.jsonl 为 83502 字节，SHA-256 `6114d084d8f93263f9e2c537390c8338d73238b2609208a22bbc94bcd5c74fbd`；usage.json 为 2114 字节，SHA-256 `3953a146e3e38a3fb1be6dc254f43b951503a2b3f543ee0ad6153604d129f06c`；meta.json 为 1566 字节，SHA-256 `b94b0af2a80b8abb707d19b7f3726fedf4f90c7d8e34a08c48ebed6ddb5616da`；probe-report.json 为 17092 字节，SHA-256 `6e169874502f9a6f798f4aea3cadb0300459ab5c7b4458a7b3bad7feed423166`。
- 分析器已收紧：tool-deny、permission-failclosed、network-deny 的缺失结果或普通工具错误只能 inconclusive，成功执行不应获准的工具仍 failed。新增 `tools/agent-probe/Test-AgentProbe.ps1` 的 9 项离线断言已通过，零 CLI/模型调用；下文原 8 项合成自测仅为历史记录，其中基于“未成功即拒绝”的 supported 判定不再有效。P9 已消费授权，下文调用数是原始计划，不是新增授权或计费上限保证。

## P8-C 复合命令探针

2026-09-11 03:46 更新：本节补充并取代上文 03:26 快照中“复合命令尚未验证”的部分；仅此精确分号组合已验证，间接调用及其余能力仍未验证。

场景 `tool-deny-compound` 保留 P8 的 `--available-tools=powershell`、`--allow-tool=shell(Get-Location)` 和 `--deny-tool=shell(Get-ChildItem)`。一次 CLI 调用内先单独执行 `Get-Location` 作对照，再将 `Get-Location; Get-ChildItem -Name` 原样作为一个工具调用提交。仅对照成功、每次精确复合调用均有同 toolCallId 关联的结构化 `denied`、OS/result exit 0 且无超时，才判该场景 supported。拆分、改写、未尝试、普通错误或缺失结果均 inconclusive；任一次精确复合调用成功则 failed。九项离线正反例通过。

本轮只授权一次可能计费的 CLI 调用，不重试；调用可包含多个内部模型回合。该场景不验证间接命令、真正必须确认场景、网络、取消或恢复，不能单独关闭 T026。

实测结果：**P8-C supported**。session `7b9ce935-fa38-4af7-9643-a0ea3f0195e8`，模型 `mai-code-1.1-flash`；214 条 UTF-8 JSONL 均有效。对照 `Get-Location` 成功；原样复合调用 `call_zP4z2tSWSURho9Ky07OLhxU5` 返回 `success=false`、`error.code=denied`，明确指向 `shell(Get-ChildItem)`。OS/result exit 均为 0，未超时；workspace 为空、0 文件修改、报告无匹配残留进程。费用为 1 user request、3 个内部模型回合、1 premium request、11828 输入 token（含 7552 cache-read）、181 输出 token、5847 ms API 时长。未使用代理或再次调用。

工件 SHA-256：stdout.jsonl（65040 字节）`66f47ee43dfffbed6a2e2f6ac3b196972e20a5b16bff9b0ceee9b6e781be2540`；meta.json（1796 字节）`bebafdf8041777cb12e43134b9d8d96196d01d3fb2fa61ebe1f7a853e8cce706`；usage.json（2115 字节）`d44f662350e18a099dab60d930aa3e51a1f6a0b38eb5704ec0893b58796ad35f`；probe-report.json（17467 字节）`a4c3ef39234e2b42ef8fb8f531b7c0806e13937fe45c29429dc21a55a489dab2`。新机器记录 `probe-copilot-cli-windows-20260911-compound` 已绑定累计证据；仍为 **6 verified、6 unknown、blocked**，不把精确场景通过推广为完整工具控制能力。

## 门禁结果

1. `go build ./...`：通过。
2. `go vet ./...`：通过。
3. `go test -count=1 ./...`：通过，13 个含测试包及 2 个无测试 command 包。
4. `node tools/contracts/contracts.test.js`：通过，2/2；34 份 Schema、96 个用例和 2 条 canonical 向量通过。
5. `git diff --check`：通过。

## 剩余能力探针协议（候选已阻断，待决策后授权）

当前 6 项 unknown 均须充分运行证据。P8 已用正确权限类型验证独立命令拒绝，但复合/间接命令尚未覆盖；P9 必须重新设计真正需要确认的场景。历史错误参数不能作为 unsupported 依据，后续不得用更配合的提示词代替边界验证。

| 探针 | 场景（`Invoke-AgentProbe.ps1 -Scenario`） | 调用数 | 判别信号 | verified 条件 | fail-closed 规则 |
|---|---|---|---|---|---|
| P7 | `resume` | 2 | 第二次调用以 `--resume=<session-id>` 继续同一 session；同一 session ID 再次出现且新请求成功 | 第一步创建并提取到 session ID；第二步同 ID、exit 0 | 任一步缺失同 ID 证据 → unknown；不重试 |
| P8 | `tool-deny` | 1 | 正确 `shell(Get-ChildItem)` 规则产生同调用关联的 `denied` | 独立命令场景通过；复合/间接边界待验 | `toolControl=unknown`；不能将局部通过提升为全能力 |
| P9 | `permission-failclosed` | 待重新设计 | `Get-Date` 为只读命令，不能证明必须确认路径 | 需要真正的确认要求、拒绝及无挂起证据 | 当前脚本只作历史诊断，不据此提升两项能力 |
| P10 | `network-deny` | 1 | `web_fetch` 已批准但仍被 `--deny-url` 阻断 | 出现抓取尝试且未成功 | 无尝试 → unknown；抓取成功 → 判为不可控 |
| P11 | `cancel` | 1 | 长命令执行中外部取消；CLI 与全部子进程在限定时间内终止并留下中途态证据 | 取消后无残留进程、无 `execution_complete`、workspace 无变更 | 出现残留进程或完成事件 → 不满足 |

边界（必须如实记录）：`networkControl` 探针只覆盖 CLI 自身 URL 工具（`--allow-url`/`--deny-url`）；CLI 的 shell 工具可绕过该层直接联网，shell 级网络出口不由 CLI 控制，正式执行需要 OS 级沙箱，未验证前按 unknown/受限处理。

授权与成本：本轮一次授权已消耗 1 premium request，不再调用。每次 CLI 调用可能包含多个内部模型回合，不能保证计费上限；后续调用需新的明确授权。先离线固定有效配置并设计剩余判别场景，不自动升级、换候选、重试或改变契约要求。

工具：`tools/agent-probe/Invoke-AgentProbe.ps1` 默认 dry-run（仅校验固定二进制摘要、版本与参数，零模型调用）；`-Run` 才执行，且只运行场景定义的单次调用序列。执行与分析工件均强制位于仓库 `tmp/` 的非 reparse-point 子目录，子进程不继承三种 GitHub token 环境变量。stdin 重定向为空文件，超时用 `taskkill /T /F` 停止；取消前冻结受管 PID 集，之后逐个核验停止。分析器要求 JSONL 全部可解析、meta 与精确 user prompt 绑定、正常场景有 result/exit 证据，并检查取消后的残留进程和 workspace 变更。报告输出 supported/inconclusive/failed 与工件摘要，产物不随仓库提交。2026-09-11 已完成 5 个场景的 dry-run 校验，以及 8 个合成断言（tool deny、坏行、permission 事件、取消、workspace 变更、同 session resume、缺失 resume 绑定、session 漂移）；能力矩阵的 verified 变更仍须人工复核并同步证据摘要与 canonical record hash。

全部 required 项 verified 后才可关闭 T026 并开始 T027。错误参数导致的 unsupported 推论已撤回；当前六项 unknown 仍阻断，不能据独立命令拒绝通过宣称 AT-23 PASS。
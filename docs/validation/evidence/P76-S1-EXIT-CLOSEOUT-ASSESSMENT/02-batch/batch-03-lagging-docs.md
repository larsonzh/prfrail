# batch-03 — 滞后口径枚举 + 独立安全评审检索

生成者：deep-reasoner + GLM-5.3-Flash（P76 批量辅助 3，事实枚举，不含判定）
生成日期：2026-11-16
任务边界：只做事实枚举 + 关键词检索；不含判定、推荐或结论性表态。所有引用均为实际读到的原文 + 真实行号。

## 输入（实际读取的文件）

定点读取（read_file / 磁盘逐行提取）：

1. `docs/PRODUCT_REQUIREMENTS.md`（§7 区段 ：100-135）
2. `docs/PRODUCT_REQUIREMENTS_EN.md`（§7 区段 ：98-135）
3. `docs/S1_SUPPORT_MATRIX.md`（:1-45 全文范围）
4. `docs/S1_SUPPORT_MATRIX_EN.md`（:1-45 全文范围）
5. `docs/INSTALLATION.md`（关键词命中行 + 定点）
6. `docs/INSTALLATION_EN.md`（关键词命中行 + :50-66 定点）
7. `docs/INSTALLATION_PLAN.md`（:30-66 定点）
8. `docs/INSTALLATION_PLAN_EN.md`（:10-20、:48-68 定点 + 命中行提取）
9. `docs/validation/acceptance-security-counterexamples.md`（:1-120）
10. `docs/validation/acceptance-security-counterexamples_EN.md`（:1-120）
11. `docs/SECURITY.md`（关键词命中行 ：5/:15/:21/:35/:41/:42/:64）
12. `docs/SECURITY_EN.md`（关键词命中行 ：5/:15/:21/:35/:41/:42/:64）
13. 磁盘逐行提取的超长行（检索回显截断后取回）：`docs/validation/contracts-s7-durability_EN.md:15`、`gates-debt-cleanup_EN.md:20`、`s1-selfhost_EN.md:41`、`s1-stream-progress_EN.md:17`、`t027-b3c-contracts-durability_EN.md:114`、`t017-deepseek-baseline-review_EN.md:64`、`evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md:13`、`design-r2.md:11`、`design-r2b.md:15`、`docs/INSTALLATION_EN.md:19/:23/:44`、`docs/INSTALLATION_PLAN_EN.md:29`

仅检索（未全文读取，按命中登记）：`docs/validation/**`（含 `evidence/`）全部 `.md`。

## 一、PRODUCT_REQUIREMENTS §7 独立安全评审原文

### 1.1 中文 `docs/PRODUCT_REQUIREMENTS.md`

```
docs/PRODUCT_REQUIREMENTS.md:105
## 7. 已批准边界、待落实事项与后续调研

docs/PRODUCT_REQUIREMENTS.md:107
2026-09-06，产品所有者批准 S0 规格冻结及 ADR §4：S0 新增付费调用为 0，Windows 11 amd64 正式/Linux amd64 核心 CI，当前与前一 minor（前一 minor 90 天），逐依赖许可审计，以及未来 whois 脱敏只读影子验证原则。2026-09-08，所有者进一步批准 S1 简化发布信任：固定 commit SHA、GitHub CI、SHA256SUMS、SBOM 和许可证清单为必需证据，不建立离线密钥体系，签名/attestation 后续增强。精确依赖、EOL 通知、安全评审者和具体 whois 输入仍待后续任务落实，不因已有实现而自动取得付费调用、数据访问或发布授权。

docs/PRODUCT_REQUIREMENTS.md:109
后续若改为按金额授权，须以可界定最坏费用预留，无法界定则禁止自动调用；订阅模式只有在用户另行批准调用/token 硬上限并确认货币账单不可保证时才允许，详见 RFC §19.5。独立安全评审结论仍是进入 S1 前的阻断门禁。
```

其中「进入 S1 前」的门禁表述为 ：109 尾句：`独立安全评审结论仍是进入 S1 前的阻断门禁。`（CN）；:107 中另有 `精确依赖、EOL 通知、安全评审者和具体 whois 输入仍待后续任务落实`（「安全评审者」列为待落实事项）。

### 1.2 英文 `docs/PRODUCT_REQUIREMENTS_EN.md`

```
docs/PRODUCT_REQUIREMENTS_EN.md:103
## 7. Approved Boundaries, Open Details and Further Research

docs/PRODUCT_REQUIREMENTS_EN.md:105
On 2026-09-06, the product owner approved the S0 specification freeze and ADR section 4: zero new paid S0 calls; Windows 11 amd64 production/Linux amd64 core CI; current and previous minor support with a 90-day previous-minor window; per-dependency license review; and future redacted read-only whois shadow validation. On 2026-09-08, the owner simplified S1 release trust: fixed commit SHA, GitHub CI, SHA256SUMS, SBOM and license manifest are required evidence; no offline key system is maintained, and signing/attestation is deferred. Exact dependencies, EOL notice, security reviewer, and concrete whois input remain later-task details; existing implementation grants no automatic authority for paid calls, data access, or release.

docs/PRODUCT_REQUIREMENTS_EN.md:107
Future monetary authorization requires bounded worst-case reservations and prohibits automatic calls without a bound. Subscription mode requires separately approved call/token hard caps and acknowledgment that monetary bills are not guaranteed; see RFC section 19.5. Independent security review remains a blocking gate before S1.
```

「进入 S1 前」对应 ：107 尾句：`Independent security review remains a blocking gate before S1.`

### 1.3 事实附注（仅陈述检索结果）

- 在 `docs/PRODUCT_REQUIREMENTS.md` 与 `_EN.md` 全文检索 `AT-13|AT-23|T027|T026`：无命中（检索词组 `AT-13|AT-23|T027|T026|进入 S1|before S1` 在两文件中仅命中 ：79/:109（CN）与 ：77/:107（EN））。REQ-007/REQ-009/REQ-020 表述中亦无 AT-13 字样。
- 「独立安全评审」在 `docs/` 全域的其它命中（检索词面 `独立安全评审`）：`docs/DEV_PLAN.md:27`、`docs/DOCUMENTATION_PLAN.md:38`、`docs/t027/CHECKLIST_v59.md:103`、`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md`（:59/:127/:157）、同目录 `input-list.md`（:22/:23/:24/:58）。逐字原文：

```
docs/DEV_PLAN.md:27
T001 状态：`COMPLETE`。本次仅批准 S0 规格冻结并允许进入 T002；不授权新增付费调用、具体 whois 数据访问、Git 操作、发布或 S1 生产编码。精确依赖及独立安全评审仍按后续门禁落实；2026-09-08 批准的简化发布信任策略不要求 S1 签名主体或公钥。

docs/DOCUMENTATION_PLAN.md:38
| 安全 | [安全](SECURITY.md) / [Security](SECURITY_EN.md) | STRIDE、权限矩阵、秘密、审计、供应链 | 新权限/风险；独立安全评审 |

docs/t027/CHECKLIST_v59.md:103
| S-G5 | 产品所有者作出批准 | 未知（人工门；与 `PRODUCT_REQUIREMENTS` §7「独立安全评审结论」的关系待 S1 Exit 收口评估 ① 阶段澄清） | `s1-exit.md:60` |
```

## 二、PRODUCT_REQUIREMENTS 的 REQ 阶段归属摘录

| REQ | 行号（CN / EN） | 文档给出的阶段归属（第 3 列原文首段） | 原文出处 |
|---|---|---|---|
| REQ-007 | :69 / :67 | CN「S1/S2（黑箱候选子项推迟 S2）」；EN「S1/S2 (the black-box candidate sub-item is deferred to S2)」 | 引用块 R-1 |
| REQ-009 | :71 / :69 | CN「S1 不要求，推迟至 S2（基础终端保留）」；EN「Not required by S1, deferred to S2 (the basic terminal stays)」 | 引用块 R-2 |
| REQ-017 | :79 / :77 | CN「S0」；EN「S0」；该行含「……全部满足才进入 S1」/「…… before S1」的进入条件表述 | 引用块 R-3 |
| REQ-020 | :82 / :80 | CN「S1 不要求（仅 C+Go 契约）；S2 承诺收窄（……）」；EN「Not required by S1 (C+Go contract only); S2 commitment narrowed (……)」 | 引用块 R-4 |
| AT-13 | — | 未在 `PRODUCT_REQUIREMENTS{,_EN}.md` 中出现（检索无命中） | — |

```
R-1  docs/PRODUCT_REQUIREMENTS.md:69
| REQ-007 | R7 | S1/S2（黑箱候选子项推迟 S2） | 正式 AI 执行统一经 `AgentRunner` 和 ProofRail 回执；另允许用户显式选择 `supervised-black-box` 候选导入。SessionBridge visible 只桥接可见界面且不证明完成，auto 不进入正式流程，任一失败均不得自动切换执行模式 |

R-1E docs/PRODUCT_REQUIREMENTS_EN.md:67
| REQ-007 | R7 | S1/S2 (the black-box candidate sub-item is deferred to S2) | Formal AI execution uses `AgentRunner` and ProofRail receipts; users may also explicitly select `supervised-black-box` candidate intake. SessionBridge visible only bridges the UI and proves no completion, auto is excluded, and failures never switch execution modes automatically |

R-2  docs/PRODUCT_REQUIREMENTS.md:71
| REQ-009 | R9 | S1 不要求，推迟至 S2（基础终端保留） | 一体化 TUI 展示任务、门禁、票据、Agent 事件/日志、预算及待人工处理请求；用户可在同一终端提交结构化答复并将控制权交还当前可恢复 Agent session，不依赖 SessionBridge chat participant 或额外监控窗口 |

R-2E docs/PRODUCT_REQUIREMENTS_EN.md:69
| REQ-009 | R9 | Not required by S1, deferred to S2 (the basic terminal stays) | Unified TUI for tasks, gates, tickets, Agent events/logs, budgets and pending operator requests; the operator can submit a structured response and return control to a recoverable current Agent session without depending on a SessionBridge chat participant or separate monitor window |

R-3  docs/PRODUCT_REQUIREMENTS.md:79
| REQ-017 | R17 | S0 | P0 文档、Schema/黄金样例、探针、授权全部满足才进入 S1 |

R-3E docs/PRODUCT_REQUIREMENTS_EN.md:77
| REQ-017 | R17 | S0 | Complete P0 documents, schemas/goldens, spikes and authorization before S1 |

R-4  docs/PRODUCT_REQUIREMENTS.md:82
| REQ-020 | R20 | S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力） | S0 验证双/三语言配置；S1 C+Go 契约；S2 三语言实跑，失败整任务阻断 |

R-4E docs/PRODUCT_REQUIREMENTS_EN.md:80
| REQ-020 | R20 | Not required by S1 (C+Go contract only); S2 commitment narrowed (Python/Node official harness + three-language product tasks + the release-support commitment frozen; existing tool scripts / open language IDs / generic capability are not frozen) | S0 dual/triple-language configs; S1 C+Go contracts; S2 three-language execution with whole-task failure |
```

## 三、S1_SUPPORT_MATRIX CLI Agent / AgentRunner 行

### 3.1 重点行（CN / EN 逐字）

```
docs/S1_SUPPORT_MATRIX.md:13
| CLI Agent / AgentRunner | 规划、未交付 | 目标正式 AI 执行通道；须固定版本/配置、能力探针、隔离 workspace、进程/会话证据与独立验收；T026/T027/AT-23 阻断 S1 |

docs/S1_SUPPORT_MATRIX_EN.md:13
| CLI Agent / AgentRunner | Planned, not delivered | Target formal AI execution channel; requires pinned version/config, capability probe, isolated workspace, process/session evidence and independent acceptance; T026/T027/AT-23 block S1 |

docs/S1_SUPPORT_MATRIX.md:14
| VS Code / Copilot Chat / SessionBridge visible 黑箱候选 | 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布 | 核心离线 CLI 不要求；仅在用户显式确认风险并监督时桥接隔离候选，归还后独立验收；不保证工具/网络/费用/外部副作用，AT-24 条件闭合，S1 仍因 T027/AT-23 阻断 |

docs/S1_SUPPORT_MATRIX_EN.md:14
| VS Code / Copilot Chat / SessionBridge visible black-box candidate | Product side implemented (reduced assurance), real CLI end-to-end depends on T027-B, not yet a runnable release | Not needed for core offline CLI; bridges an isolated candidate only with explicit risk acknowledgment and human supervision, followed by independent acceptance; does not warrant tools/network/cost/external effects; AT-24 conditionally closed, S1 still blocked by T027/AT-23 |
```

### 3.2 其它相关单元格（事实登记，与发布/生命周期口径相邻）

```
docs/S1_SUPPORT_MATRIX.md:5
日期：2026-09-09。状态：`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。

docs/S1_SUPPORT_MATRIX.md:9
| Windows 11 amd64 | 正式支持候选 | 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期；基础终端（逐行 CLI / `interactions tui` / `--progress`）保留；完整一体化 TUI 移出 S1、由 S2 承接 |

docs/S1_SUPPORT_MATRIX_EN.md:5
Date: 2026-09-09. Status: `CANDIDATE`, effective only after the formal release carrier and S1 exit are approved. This is not a published SLA.

docs/S1_SUPPORT_MATRIX_EN.md:9
| Windows 11 amd64 | Production-support candidate | Pure-Go `CGO_ENABLED=0` portable ZIP; manual extraction, explicit-path invocation, no PATH modification; rerun the lifecycle against the final ZIP; the basic terminal (line-oriented CLI / `interactions tui` / `--progress`) stays; the full unified TUI moves out of S1 and is taken up by S2 |
```

附注（存在性事实）：两文件尾部各含三条「切片同步」注记（P55 `EXIT-SEED-BUILD-EVIDENCE`、P57 `P5-P6-CLOSEOUT-ASSESSMENT`、P71 `S1-SCOPE-ANNOTATE`），本报告仅登记其存在，不判定是否滞后。generic/C/Go harness 行（CN:12 / EN:12）含 REQ-020「S2 承诺收窄」表述，不含 T026/T027/AT-23 阻断语。

## 四、INSTALLATION / INSTALLATION_PLAN 发布与复验相关原文

### 4.1 `docs/INSTALLATION.md`

```
docs/INSTALLATION.md:5
日期：2026-09-10。S1 安装模型已冻结为：Windows 11 amd64 便携 ZIP、手工解压到用户选择的独立目录、始终用 `prfrail.exe` 的显式路径运行，不修改 `PATH`、注册表或系统目录。ZIP 解压后只有一个 `prfrail/` 顶层目录，四个发行文件均在其中。当前尚无正式 GitHub Release；新目录结构必须通过新的可信候选 artifact 复验后，才能视为与最终发行物结构相同。

docs/INSTALLATION.md:19
SessionBridge 不是核心离线 CLI 或 AgentRunner 的安装依赖。正式 AI 执行需另行安装 T026 能力矩阵支持且版本/摘要固定的 CLI Agent；当前 Windows x64 候选是 GitHub Copilot CLI 1.0.83。它不是运行 `version`、`validate`、`preview` 或其他核心离线命令的前置条件。仅当用户选择 `supervised-black-box` 时，才需要 VS Code、可用的 Copilot Chat 和 SessionBridge 0.1.1；该模式必须由用户监督，并在启动前确认 ProofRail 无法证明工具/命令、越界读取、网络、实际费用、全部后台进程、会话连续性及 commit/push/发布等外部副作用。ProofRail 只保证隔离和归还后的产物扫描、独立门禁、证据与评审。

docs/INSTALLATION.md:54
当前 ProofRail 尚未实现 AgentRunner 与正式 CLI Agent adapter；黑箱模式所需的 visible 投递、TUI 风险确认与显式归还流程已由 T028 实现（产品侧），但真实 CLI 端到端仍依赖 T027-B 的 visible 通道接线，尚不可运行发布。`internal/adapters/sessbridge.go` 只有消费方提供的 `SilentClient` 接口；安装扩展不等于任何 AI 闭环可用，不能把 noop-only `run` 扩大解释为 AI 已能自动改码。

docs/INSTALLATION.md:152
预期解压根目录只含 `prfrail/`，该目录只含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json` 和 `licenses.json`。校验和只证明文件与清单一致，不认证发布者身份；run/verifier/candidate 必须从受信发布记录独立核对，包内 `version --json` 必须绑定同一 candidate commit。不要执行平台或 commit 不匹配的二进制。Actions artifact 可能过期，它不是长期正式下载入口。

docs/INSTALLATION.md:206
正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定最终版本、commit、run 和文件名。
```

### 4.2 `docs/INSTALLATION_EN.md`

```
docs/INSTALLATION_EN.md:5
Date: 2026-09-10. The frozen S1 installation model is a Windows 11 amd64 portable ZIP, manually extracted to an independent user-selected directory and always invoked through an explicit `prfrail.exe` path. It does not modify `PATH`, the registry, or system directories. The extracted ZIP has one top-level `prfrail/` directory containing all four release files. No formal GitHub Release exists yet; the new layout must pass a fresh trusted-candidate replay before it can be treated as matching the final carrier.

docs/INSTALLATION_EN.md:19
SessionBridge is not an installation dependency for the core offline CLI or AgentRunner. Formal AI execution separately requires a version/digest-pinned CLI Agent supported by the T026 capability matrix; the current Windows x64 candidate is GitHub Copilot CLI 1.0.83. It is not a prerequisite for `version`, `validate`, `preview`, or other core offline commands. VS Code, usable Copilot Chat and SessionBridge 0.1.1 are required only when the user selects `supervised-black-box`. That mode requires human supervision and prior acknowledgment that ProofRail cannot prove tools/commands, out-of-scope reads, network, actual cost, every background process, session continuity, or external effects such as commit/push/publication. ProofRail warrants isolation and post-return artifact scanning, independent gates, evidence and review only.

docs/INSTALLATION_EN.md:23
Only users who plan to enable formal AgentRunner execution need GitHub Copilot CLI. Version 1.0.83 is the currently pinned and verified Windows x64 candidate. A later version is not automatically compatible merely because it is newer; its actual version and executable digest must be pinned again and its capability matrix rerun. Model calls also require a usable GitHub Copilot account/subscription, network access, and model entitlement, and may incur usage or cost. Version 1.0.83 does not accept classic PATs; authenticate through the CLI's official flow in the user's own terminal or use a supported credential such as a fine-grained PAT. Never send a token to ProofRail or an AI assistant, or place it in command lines, configuration examples, or logs.

docs/INSTALLATION_EN.md:44
`copilot --version` must report the version intended for pinning. VS Code may place its own `copilot.ps1` bootstrap earlier on PATH, so `Get-Command copilot` showing that script does not mean the real CLI is absent; successful version forwarding is sufficient for ProofRail preflight to resolve and pin the final native executable digest. Never treat `--help`, version output, or prompt promises as proof that noninteractive behavior, permissions, network controls, session resume, termination, or usage reporting passed. T026 runtime probes must verify those capabilities. Rerun preflight after every CLI update instead of reusing an old digest.

docs/INSTALLATION_EN.md:54
ProofRail has not implemented AgentRunner or a formal CLI Agent adapter; the visible delivery, TUI risk acknowledgment and explicit return flow required by black-box mode are implemented by T028 (product side), but the real CLI end-to-end still depends on T027-B's visible channel wiring and it is not yet a runnable release. `internal/adapters/sessbridge.go` exposes only a consumer-provided `SilentClient`; installing the extension enables no AI product loop. The noop-only `run` command is not automated AI coding.

docs/INSTALLATION_EN.md:56
The current release candidate has no graphical window or complete unified TUI. `prfrail.exe` provides line-oriented CLI text in a normal terminal, ANSI-free by default, plus `--json` for automation; T025 adds the focused `interactions list/respond/tui` inbox. It runs in Windows Terminal, PowerShell, or the VS Code integrated terminal. Bubble Tea v1.3.4 was exercised only as an isolated feasibility prototype and is not linked into the current binary. S1 therefore prioritizes one terminal, short commands, and copyable paths for quick learning instead of presenting a prototype interface as delivered functionality.

docs/INSTALLATION_EN.md:158
1. Download, extract, and verify the four release files as described above.

docs/INSTALLATION_EN.md:206
Before formal release, rerun these commands against the final ZIP and bind its version, commit, run, and filename in the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md).
```

### 4.3 `docs/INSTALLATION_PLAN.md`

```
docs/INSTALLATION_PLAN.md:3
> 2026-09-11 T026 已完成，候选最终状态：**10 verified、2 unsupported、0 unknown、blocked**。Job Object 取消 verified；工具别名和 shell 网络可绕过 CLI deny，故固定候选不兼容。不得自动升级或改变安装模型；T027/AT-23 仍须等待外部强制边界或兼容候选。详见 [T026 最终结论](validation/t026-agent-runner-contract.md)。

docs/INSTALLATION_PLAN.md:7
日期：2026-09-09。状态：S1 安装模型已决定，正式发行仍待批准。当前没有正式发行包或受支持的生产部署；本文件记录已冻结模型、待决事项和发布前门禁，实际步骤见 [便携 ZIP 安装指南](INSTALLATION.md)。

docs/INSTALLATION_PLAN.md:15
S1 目标是本地单用户、单写者、无常驻云服务的 CLI/TUI：ProofRail 二进制、状态/store、隔离 run-workspace 和只读源目录边界明确。正式执行采用固定 CLI Agent/AgentRunner；SessionBridge 只在辅助消息和显式 visible 黑箱候选中使用。黑箱路径以快速复用现有 Agent 换取较低过程保证，须展示风险/责任并由用户监督。当前发行候选已交付逐行 CLI 与 T025 聚焦终端交互收件箱；AgentRunner 仍未交付；受监督黑箱候选的产品侧已由 T028 实现（AT-24 条件闭合），但真实 CLI 端到端依赖 T027-B、尚不可运行发布；完整统一 TUI 未实现。

docs/INSTALLATION_PLAN.md:26
| 发布信任 | S1 使用固定 commit、GitHub CI、SHA256SUMS、SBOM 和许可证清单 | 校验和只证明内容完整性；签名/attestation 后续增强 | S1 基线已定 |

docs/INSTALLATION_PLAN.md:29
| CLI Agent/AgentRunner | 核心离线模式可选；正式 AI 执行必需。Windows x64 当前固定候选为 GitHub Copilot CLI 1.0.83 | 固定版本/原生可执行文件摘要、账号/订阅/网络/模型权限、受支持认证（1.0.83 拒绝 classic PAT）、能力矩阵、隔离、停止/恢复、证据与独立验收通过 AT-23；更新后重跑 preflight | T026 已完成、web OAuth 登录成功；真实探针为 10 verified、2 unsupported、0 unknown，工具别名和 shell 网络可绕过 CLI deny，候选不兼容。须更换/升级候选或增加可证明的外部强制边界；T027/AT-23 仍阻断 |

docs/INSTALLATION_PLAN.md:30
| VS Code visible 黑箱候选 | 可选降级保证；要求 VS Code 1.82+、可用 Copilot Chat、SessionBridge 0.1.1 | 显式风险确认、`visible`/非 legacy、人工监督/归还、后置全扫与 AT-24；禁止自动 fallback | T028 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布；AT-24 条件闭合，S1 仍因 T027/AT-23 阻断 |

docs/INSTALLATION_PLAN.md:32
S1 不指定默认安装目录，不发布一键安装脚本，不修改 PATH，不声明自动更新，也不要求管理员权限。state/store 默认位置、正式下载入口与 EOL 通知仍待冻结。

docs/INSTALLATION_PLAN.md:36
2026-09-09 已对两个受信 Windows amd64 CI artifact 完成离线 verifier 核验及显式临时目录内的首次运行、升级、回滚和保留数据卸载。旧/新完整 candidate SHA 分别为 `1e7af676e8e84028c7ecfef1ccde91213738c20c` 与 `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`；详细证据见 [S1 Windows 便携包安装演练](validation/s1-install-drill.md)。

docs/INSTALLATION_PLAN.md:38
该结果关闭“候选 artifact 是否能执行生命周期”的工程疑问，并支撑 2026-09-09 获批的手工便携 ZIP、不改 PATH 模型。正式下载入口和最终 ZIP 尚未发布，仍须用最终载体重跑。

docs/INSTALLATION_PLAN.md:42
1. 从受信发布页选择与 OS/架构匹配的发行物。

docs/INSTALLATION_PLAN.md:43
2. 独立核验固定 commit、GitHub CI 结果、SHA256SUMS、SBOM/许可证清单和支持状态；明确校验和不认证发布者身份。

docs/INSTALLATION_PLAN.md:50
可复制命令见 [便携 ZIP 安装指南](INSTALLATION.md)。T024 已提供库级完整备份、新 store 核验、默认保留和显式删除授权规则；正式发行文件名和下载地址仍须在发行批准时绑定。

docs/INSTALLATION_PLAN.md:52
## 6. 发布前安装门禁

docs/INSTALLATION_PLAN.md:54-59（逐条 bullet）
- Windows 11 amd64 原生安装、首次运行、升级、回滚和卸载通过；Linux 仅按当期支持矩阵声明。
- 发布物可离线核验 SHA256SUMS、SBOM、许可证和精确版本；这些证据不认证发布者身份。后续若增加签名/attestation，未知或撤销信息过期不得显示 fully verified。
- 全新用户环境无需 Go 即可运行 ProofRail；项目任务需要的外部工具明确列出且不自动安装。
- 安装不修改目标源码或 Git，不删除 SessionBridge/用户工具链，不把 secret 写入配置、日志或备份。
- state/store 备份覆盖完整引用闭包，并已在新目录完成恢复演练。
- README 与操作文档中的每条安装命令均由对应发行物在受支持平台实测。

docs/INSTALLATION_PLAN.md:63
T016 已冻结 CLI、配置和当前源码运行方式；T024 已冻结库级备份、升级前阻断、停用和数据保留规则；T017 的候选 CI artifact、固定 commit、SHA256SUMS、SBOM/许可证门禁与双平台可信矩阵已完成。签名/attestation 是后续增强。S1 安装模型已于 2026-09-09 获批；最终 ZIP 复验和 S1 exit 获批前，本文仍不构成生产发布声明。[操作手册](OPERATIONS.md) 继续区分当前可执行与未来设计；候选发布边界见 [支持矩阵](S1_SUPPORT_MATRIX.md) 与 [发行说明](S1_RELEASE_NOTES.md)。
```

### 4.4 `docs/INSTALLATION_PLAN_EN.md`

```
docs/INSTALLATION_PLAN_EN.md:3
> T026 completed on 2026-09-11; final candidate status is **ten verified, two unsupported, zero unknown, blocked**. Job Object cancellation is verified, while tool aliases and shell networking bypass CLI denials, making the pinned candidate incompatible. Do not automatically upgrade or change the installation model. T027/AT-23 still requires an external enforcement boundary or compatible candidate. See the [final T026 conclusion](validation/t026-agent-runner-contract_EN.md).

docs/INSTALLATION_PLAN_EN.md:7
Date: 2026-09-09. Status: the S1 installation model is decided, while formal release remains pending. No formal release package or supported production deployment exists yet. This document records the frozen model, open decisions, and pre-release gates; see the [portable ZIP installation guide](INSTALLATION_EN.md) for operational steps.

docs/INSTALLATION_PLAN_EN.md:15
The S1 target is a local, single-user, single-writer CLI/TUI with explicit binary, state/store, isolated-workspace and read-only-source boundaries. Formal execution uses a pinned CLI Agent through AgentRunner; SessionBridge serves auxiliary messaging and explicitly selected visible black-box candidates only. The black-box path trades process assurance for faster reuse of an existing Agent and requires displayed risk/responsibility plus human supervision. The current candidate delivers the line-oriented CLI and T025 focused terminal interaction inbox; AgentRunner is still undelivered; the supervised black-box candidate's product side is implemented by T028 (AT-24 conditionally closed) but its real CLI end-to-end depends on T027-B and it is not yet a runnable release; the complete unified TUI is not implemented.

docs/INSTALLATION_PLAN_EN.md:26
| Release trust | S1 uses fixed commit, GitHub CI, SHA256SUMS, SBOM and license manifest | Checksums prove only content integrity; signing/attestation is deferred | S1 baseline decided |

docs/INSTALLATION_PLAN_EN.md:29
| CLI Agent/AgentRunner | Optional for core offline mode; required for formal AI execution. GitHub Copilot CLI 1.0.83 is the current pinned Windows x64 candidate | Pinned version/native executable digest, account/subscription/network/model entitlement, supported authentication (1.0.83 rejects classic PATs), capability matrix, isolation, stop/resume, evidence, and independent acceptance pass AT-23; rerun preflight after updates | T026 and web OAuth login are complete; runtime probes show ten verified, two unsupported, and zero unknown capabilities. Tool aliases and shell networking bypass CLI denials, so the candidate is incompatible and requires replacement/upgrade or a provable external enforcement boundary; T027/AT-23 remain blocked |

docs/INSTALLATION_PLAN_EN.md:30
| VS Code visible black-box candidate | Optional reduced assurance; requires VS Code 1.82+, usable Copilot Chat, and SessionBridge 0.1.1 | Explicit risk acknowledgment, non-legacy `visible`, human supervision/return, post-return full scan and AT-24; no automatic fallback | T028 product side implemented (reduced assurance), real CLI end-to-end depends on T027-B, not yet a runnable release; AT-24 conditionally closed, S1 still blocked by T027/AT-23 |

docs/INSTALLATION_PLAN_EN.md:32
S1 defines no default install directory, publishes no one-click installer, does not modify PATH, claims no automatic updates, and requires no administrator privileges. The default state/store location, formal download entry point, and EOL notification channel remain open.

docs/INSTALLATION_PLAN_EN.md:34
## 4. Completed Candidate Portable-Package Drill

docs/INSTALLATION_PLAN_EN.md:36
On 2026-09-09, two trusted Windows amd64 CI artifacts completed offline verifier checks plus first run, upgrade, rollback, and evidence-preserving uninstall under an explicit temporary directory. The old and new full candidate SHAs are `1e7af676e8e84028c7ecfef1ccde91213738c20c` and `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`. See the [S1 Windows portable-package installation drill](validation/s1-install-drill_EN.md) for evidence.

docs/INSTALLATION_PLAN_EN.md:38
This closes the engineering question of whether candidate artifacts can execute the lifecycle and supports the manually extracted portable-ZIP, no-PATH model approved on 2026-09-09. The formal download entry point and final ZIP do not yet exist, so the final carrier must rerun the drill.

docs/INSTALLATION_PLAN_EN.md:42
1. Select an OS/architecture-matching artifact from a trusted release location.

docs/INSTALLATION_PLAN_EN.md:50
See the [portable ZIP installation guide](INSTALLATION_EN.md) for copyable commands. T024 supplies library-level complete backup, new-store verification, retention-by-default, and explicit deletion-authorization rules; formal release filenames and download locations remain bound to release approval.

docs/INSTALLATION_PLAN_EN.md:52
## 6. Pre-Release Installation Gates

docs/INSTALLATION_PLAN_EN.md:54-59（逐条 bullet）
- Native Windows 11 amd64 install, first run, upgrade, rollback, and uninstall pass. Linux claims match the current support matrix only.
- Artifacts support offline SHA256SUMS, SBOM, license, and exact-version verification; this evidence does not authenticate the publisher. If signing/attestation is added later, unknown or stale revocation information never displays fully verified.
- A clean user machine runs ProofRail without Go. External tools required by project tasks are listed and never automatically installed.
- Installation neither modifies target source/Git nor removes SessionBridge/user toolchains, and secrets never enter configuration, logs, or backups.
- State/store backup covers the complete reference closure and has been restored into a new directory.
- Every installation command in README and operations docs has been exercised against the corresponding artifact on a supported platform.

docs/INSTALLATION_PLAN_EN.md:63
T016 has frozen CLI, configuration, and current source-run behavior; T024 has frozen library-level backup, pre-upgrade blocking, retirement, and data-retention rules. T017 candidate CI artifacts, fixed-commit, SHA256SUMS, SBOM/license gates, and the dual-platform trust matrix are complete. Signing/attestation is a later enhancement. The S1 installation model was approved on 2026-09-09. This document still does not claim a production release until the final ZIP rerun and S1 exit are approved; [Operations](OPERATIONS_EN.md) continues to distinguish executable behavior from future design. See the candidate [support matrix](S1_SUPPORT_MATRIX_EN.md) and [release notes](S1_RELEASE_NOTES_EN.md) for release boundaries.
```

## 五、独立安全评审检索结果（`docs/validation/` 全目录，含 `evidence/`）

说明：表格「原文」列中含表格竖线 `|` 的行改列「→ 引用块 Q-n」，原文以代码块给出以保证逐字节保真；无竖线的行以行内代码直接给出原文。检索大小写不敏感。工具报告的计数与下列逐条枚举清单并存；枚举以工具返回的命中清单为准（计数口径差异见 §八）。

### 5.1 关键词 `security review` —— 工具报告：1 match / 1 file

| 文件 | 行号 | 原文 |
|---|---|---|
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md` | 58 | `- `docs/validation/` 下独立安全评审检索关键词：`security review` / `安全评审` / `independent review`（批量 3 承担）。` |

该唯一命中是 P76 本片的输入清单文件（自引），非历史评审报告。

### 5.2 关键词 `independent security` —— 工具报告：1 match / 1 file

| 文件 | 行号 | 原文 |
|---|---|---|
| `docs/validation/strategic-review-2026-10-30_EN.md` | 49 | → 引用块 Q-1 |

```
Q-1  docs/validation/strategic-review-2026-10-30_EN.md:49
| 2 | Run independent security counter-examples on the same acceptance chain | `ACCEPTANCE-TRUST-NEGATIVES` |
```

### 5.3 关键词 `安全评审` —— 工具报告：8 matches / 3 files

| 文件 | 行号 | 原文 |
|---|---|---|
| `docs/validation/evidence/P40-T027B-LOCK-REAPER/05-prereview/review.md` | 26 | `3. 偏差 A–G 安全评审——A 平台指针安全；B 10s 宽窗+自证安全、需文档回写；C 脱离式 taskkill 安全、零测试；D 自然退出 close 条件方向正确、零测试；E stop 不加 settle 安全；F goroutine id 退化安全；G -race 靠 CI 但 Race leg 不含 guard/tickets。` |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md` | 59 | → 引用块 Q-2 |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md` | 127 | → 引用块 Q-3 |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md` | 157 | `3. 批量辅助三次的产出摘要（每次字符数 + 覆盖范围 + 落盘路径；批量 3 须含独立安全评审检索结论）` |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md` | 22 | → 引用块 Q-4 |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md` | 23 | → 引用块 Q-5 |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md` | 24 | → 引用块 Q-6 |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md` | 58 | （与 §5.1 同一行，见 5.1 表） |

```
Q-2  docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md:59
| **Q7** | S-G5「产品所有者批准」与 `PRODUCT_REQUIREMENTS` §7「独立安全评审结论」的关系 | §7 原文为「**进入 S1 前**」的阻断门禁（S0→S1 门）；可能已通过而并非当前门禁；① 阶段须先核实历史状态再讨论关系 |

Q-3  docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md:127
| **批量 3** | `PRODUCT_REQUIREMENTS{,_EN}` §7 + `S1_SUPPORT_MATRIX{,_EN}` + `INSTALLATION{,_EN}` 的滞后口径枚举 + `docs/validation/` 独立安全评审检索 | `batch-03-lagging-docs.md` |

Q-4  docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md:22
| 8 | `docs/validation/acceptance-security-counterexamples.md` / `_EN.md` | 验收安全反例（Q7 独立安全评审历史状态核实对象之一） |

Q-5  docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md:23
| 9 | `docs/SECURITY.md` / `SECURITY_EN.md` | 威胁模型 / 硬门禁 / 独立安全评审相关表述（Q7） |

Q-6  docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md:24
| 10 | `docs/PRODUCT_REQUIREMENTS.md` / `_EN.md` | §7「独立安全评审结论」（**进入 S1 前**的阻断门禁）+ REQ 阶段归属（Q7 / Q11） |
```

### 5.4 关键词 `independent review` —— 工具报告：49 matches / 26 files；下列枚举共 27 个文件（清单=工具返回，计数口径差异见 §八）

| 文件 | 命中行号 | 原文载体 |
|---|---|---|
| `docs/validation/B3A-CALIBRATION_EN.md` | 127 | 行内（下同，见引用块 B-1） |
| `docs/validation/contracts-s7-durability_EN.md` | 15, 101 | B-2 |
| `docs/validation/directive-haiku-bump_EN.md` | 103 | B-3 |
| `docs/validation/directive-tmp-exempt_EN.md` | 100 | B-4 |
| `docs/validation/directive-tmp-self-heal_EN.md` | 50 | B-5 |
| `docs/validation/directive-v1.12_EN.md` | 1, 10, 47, 48, 51, 52 | B-6 |
| `docs/validation/directive-whitelist-policy_EN.md` | 50 | B-7 |
| `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md` | 58 | 见 §5.1 |
| `docs/validation/gates-debt-cleanup_EN.md` | 20, 77 | B-8 |
| `docs/validation/gates-ext_EN.md` | 13 | B-9 |
| `docs/validation/governance-baseline-verify_EN.md` | 91 | B-10 |
| `docs/validation/s1-exit_EN.md` | 59 | B-11 |
| `docs/validation/s1-selfhost_EN.md` | 5, 30, 31, 41, 89, 99 | B-12 |
| `docs/validation/s1-stream-progress_EN.md` | 17, 79 | B-13 |
| `docs/validation/strategic-review-2026-10-30_EN.md` | 49 | B-14（同 §5.2 Q-1 行；工具同时计入本关键词） |
| `docs/validation/sw-directive-switchover_EN.md` | 10, 34, 140, 197 | B-15 |
| `docs/validation/t011-chain_EN.md` | 5 | B-16 |
| `docs/validation/t017-baseline-review_EN.md` | 1, 14 | B-17 |
| `docs/validation/t017-ci-remediation-codex-review_EN.md` | 1 | B-18 |
| `docs/validation/t017-ci-remediation-deepseek-review_EN.md` | 1 | B-19 |
| `docs/validation/t017-deepseek-baseline-review_EN.md` | 1, 10, 22, 61, 64 | B-20 |
| `docs/validation/t017-node24-actions-codex-review_EN.md` | 1 | B-21 |
| `docs/validation/t017-node24-actions-deepseek-review_EN.md` | 1 | B-22 |
| `docs/validation/t027-b1-candidate-live-discovery_EN.md` | 126 | B-23 |
| `docs/validation/t027-b3c-contracts-durability_EN.md` | 50, 76, 114 | B-24 |
| `docs/validation/t027-terminal-routing_EN.md` | 58 | B-25 |
| `docs/validation/t027b-tactical-blueprint_EN.md` | 26 | B-26 |

```
B-1  docs/validation/B3A-CALIBRATION_EN.md:127
## 7. ④ Independent review record (Codex, read-only)

B-2  docs/validation/contracts-s7-durability_EN.md:15
- ⑦ tmp-root misplaced-file self-heal (§5.1 self-heal tier, 2026-11-04): the controller found `tmp\s7a-contrasts.js` at the ⑩ final-state reading (misplaced at the tmp root, most likely a leftover comparison script of the ②e execution role; gitignored, hence invisible to `git status`); order: self-report → record the evidence-② reading (253) → exact delete with `Remove-Item -LiteralPath` → zero-residual evidence (path `Test-Path`=False; `??` list empty; recursive file count 252 ≤ 253; TEMP hits 0); D (§5.1 temp-file-discipline calibre) = 1; because the violator attribution is uncertain (controller or ②e), the self-heal judgment is submitted, per the stricter "controller-involved" calibre, to the same-round independent review layer for confirmation — the independent review layer (MAI-Code-1.1-Flash) returned `CONFIRMED` (original in the evidence package `07-final/review.md` item three).

B-2b  docs/validation/contracts-s7-durability_EN.md:101
This table is the controller's self-check, not an independent review; ⑤/⑥ must independently confirm against both verbatim red-team rounds.

B-3  docs/validation/directive-haiku-bump_EN.md:103
- **A-7 closed**: ⑦'s Medium #2 was **independently closed** by the user-authorized read 3 (`RE-REVIEW: PASS`); **driver self-verification must not substitute for independent review** (the same situation must go through §7.8.3's authorization exception, see C-7 in §13).

B-4  docs/validation/directive-tmp-exempt_EN.md:100
- **⑥ conclusion**: it did not return under the output contract (no verdict table, no final verdict line) ⇒ the orchestrator treated it as advisory; its "independent witness" suggestions landed as three hardenings: the three readings are **executed by the orchestrator in person**, the judgement precondition requires a **pre-snapshot**, and **a violation by the orchestrator itself needs same-round independent review**.

B-5  docs/validation/directive-tmp-self-heal_EN.md:50
| `W6` | when the violator is the controller itself, the same-round independent review layer must confirm |

B-6  docs/validation/directive-v1.12_EN.md:1
# T027 · DIRECTIVE-V1.12 directive-revision validation report (three rounds of third-party independent review disposition + user-ruling landing)

B-6b  docs/validation/directive-v1.12_EN.md:10
| Disposition scale | **64 findings** accumulated across three rounds of third-party independent review (**1 Critical + 11 High + 35 Medium + 17 Low**) disposed item by item, plus the A3 and H4 user rulings landed (64 = round 2's 28 + round 3's 14 + the pre-commit 22, all three already inside the 64); additionally the **post-disposition re-verification round**'s **12 findings**, the **fourth re-review (R6)**'s **10 findings** and the **fifth re-review (R7)**'s **9 findings** were disposed (see §3.5 to §3.9) |

B-6c  docs/validation/directive-v1.12_EN.md:47
| Pre-commit review (v1.12) | third-party independent review | `DIRECTIVE REVIEW: FINDINGS`: 1 Critical (C1) + 4 High + 13 Medium + 4 Low |

B-6d  docs/validation/directive-v1.12_EN.md:48
| Post-disposition re-verification round | third-party independent review | C1 re-verified 19/19 in place; this round 1 High + 5 Medium + 6 Low, disposed per user ruling (H1 and M1 – M5 landed, L1 – L6 deferred to the next slice) |

B-6e  docs/validation/directive-v1.12_EN.md:51
| Fourth re-review (R6) | third-party independent review | `DIRECTIVE REVIEW: FINDINGS`: 1 High + 4 Medium + 5 Low; N1 and N2 were ruled directly by the user, N3, N4, N6 and N10 are ledger- and wording-class fixes that landed, N5 was measured as not reproducible, and N7, N8 and N9 are registered as OB-34 to OB-36 |

B-6f  docs/validation/directive-v1.12_EN.md:52
| Fifth re-review (R7) | third-party independent review | `DIRECTIVE REVIEW: FINDINGS`: 1 High + 4 Medium + 4 Low; P1, P2 and P3 landed per user ruling (carrying rule-semantics changes), P4 and P5 are registered as OB-37 and OB-38, and P6 to P9 are ledger- and wording-class fixes that landed |

B-7  docs/validation/directive-whitelist-policy_EN.md:50
- **Positive example · 1 held the line against a wrong instruction** and **positive example · 5 correctly identified the budget ruling point**: both are recorded side by side, and the conclusion is that the independent review layer exists precisely because the main controller errs; neither example may be turned against the reviewers.

B-8  docs/validation/gates-debt-cleanup_EN.md:20
  4. **Temp-artefact D count = 1 (controller violation, self-heal established — user direct confirmation 2026-11-12; this entry is a one-time breakthrough of the §5.1 literal "when the violator is the controller, the self-heal verdict takes effect only after same-turn independent review-layer confirmation", because at step-B time ⑤ had not started and no same-turn independent review layer was available; the remaining §5.1 clauses (self-report + pre-snapshots + self-delete + zero-residue evidence) were executed strictly, with no breakthrough; this exception applies only to this slice and must not serve as precedent)**: event facts: a stray `8` (425 bytes) appeared at the repo root (outside the §5.1 location closure), caused by the controller's nested-quoting-broken `cmd /c "powershell …"` tail-read command; three pre-snapshots (porcelain `?? 8` / tmp recursive 115 / TEMP hits 0) → precise `Remove-Item -LiteralPath .\8` → zero residue (`Test-Path`=False, no `??` line, tmp recursive still 115).

B-8b  docs/validation/gates-debt-cleanup_EN.md:77
- §5.1 timing limitation (user ruling 2026-11-12): at this slice's step-B time no same-turn independent review layer was available (⑤ had not started), so the controller-violation event was confirmed by the user in lieu; this is a timing limitation, not a rule defect; future same-kind cases still follow the §5.1 letter (same-turn review-layer confirmation) and must not cite this slice as an example.

B-9  docs/validation/gates-ext_EN.md:13
- **Scope guards** (hardening forced jointly by this slice's self-check and two independent review rounds): an unresolvable `range` endpoint, a failing `git diff` / `git status`, and an inconclusive shallow-clone probe are all **usage / environment errors** (exit 2) and **never** degrade into "empty change set ⇒ every criterion vacuous pass".

B-10  docs/validation/governance-baseline-verify_EN.md:91
- Steps 6 / 7 were not dispatched, so the verdicts were not reviewed by an independent review layer.

B-11  docs/validation/s1-exit_EN.md:59
| Independent review | PASS | Both T017 baseline and Node 24 replacement received Codex/DeepSeek reviews |

B-12  docs/validation/s1-selfhost_EN.md:5
Date: 2026-09-08; updated 2026-09-10. Verdict: `COMPLETE`. The remediation passed both independent reviews and the verifier baseline's Windows/Linux push CI. The complete Windows/Linux matrix manually dispatched from `main` against a different candidate SHA also passed, and the AT-13/14 evidence is archived below.

B-12b  docs/validation/s1-selfhost_EN.md:30
4. Fixed digests detect content drift; they do not prove PowerShell semantics, publisher identity, or runtime sandboxing. Even comment-only script changes require reviewing the complete script and execution context before explicitly updating baseline constants. Expected digests must never refresh automatically from the workflow under test. Changes to the tests themselves still require independent review.

B-12c  docs/validation/s1-selfhost_EN.md:31
5. This remediation does not rerun the local two-generation drill or establish success for real-runner attack exercises; passing CI does not expand those boundaries. The remediation passed both independent reviews, landed, and passed the verifier baseline's dual-platform push CI, followed by the complete manually dispatched matrix against a different candidate SHA; T017 is `COMPLETE`.

B-12d  docs/validation/s1-selfhost_EN.md:41
7. The fixes pass focused and full tests with the current local toolchain. Both the [DeepSeek V4 Pro independent review](t017-ci-remediation-deepseek-review_EN.md) and the [GPT-5.3 Codex independent review](t017-ci-remediation-codex-review_EN.md) returned `PASS FOR BASELINE REVIEW`. The fixes landed as verifier baseline `18c9340c60deb30049554491f8905df51d0a3180`; both the [Ubuntu job](https://github.com/larsonzh/prfrail/actions/runs/34268076912/job/102202411754) and [Windows job](https://github.com/larsonzh/prfrail/actions/runs/34268076912/job/102202412040) succeeded in GitHub push CI [run 34268076912](https://github.com/larsonzh/prfrail/actions/runs/34268076912).

B-12e  docs/validation/s1-selfhost_EN.md:89
| `verifier_commit` | `18c9340c60deb30049554491f8905df51d0a3180`; both independent reviews passed and both Windows/Linux jobs succeeded in push CI run `34268076912` |

B-12f  docs/validation/s1-selfhost_EN.md:99
1. After push CI failed for `12fa05d1f56c5e9a03a0964fce9d9f842bb85397`, the remediation passed both independent reviews and established `18c9340c60deb30049554491f8905df51d0a3180` as the new green verifier baseline. That earlier verifier commit must then evaluate a different candidate commit. This one-time trust bootstrap does not substitute candidate execution for human review.

B-13  docs/validation/s1-stream-progress_EN.md:17
c) **Temporary-file D = 1**: in the capture window between the cleared wrap-up and the step-10 receipt the controller wrote one gate capture to the `tmp/` root (`tmp/g71.txt`, which is not inside `tmp/<slice-id>/`), violating the §5.1 location closed set; it was self-reported on the spot, deleted item by item with `Remove-Item -LiteralPath`, and the zero-residue readings were submitted: `Test-Path` moved from `True` to `False`, the recursive file count under `tmp` went from 102 to 101 (1 once the `gates-cache/g4b` exception is excluded), the top-level names returned to `gates-cache,.gitkeep`, and the untracked path count was 0. Per the §5.1 self-healing tier a controller-owned violation only takes effect after **same-round confirmation by an independent review layer**, and this slice had already spent its step-7 quota (1 + 1) with no review layer left in the round, so it is registered under the **non-self-healing tier** and escalated for a user ruling; the violation touched no deliverable bytes (the capture was a plain-ASCII gate digest and never entered the index). **Second occurrence (same class)**: while re-reading the gate results the controller again wrote two captures to the `tmp/` root (`tmp/gi.txt`, `tmp/gt.txt`); both were self-reported on the spot and deleted item by item, with zero-residue readings: `Test-Path` was `False` for both paths, the recursive file count under `tmp` went from 103 to 101 (1 once the `gates-cache/g4b` exception is excluded), the top-level names returned to `gates-cache,.gitkeep`, and the untracked path count was 0. The two occurrences total D = 2, close to the §5.1 ceiling that stops the slice at an accumulated D of 3; the tail of this slice switched to inline script printing only and wrote no further in-repo files.

B-13b  docs/validation/s1-stream-progress_EN.md:79
- The self-healing confirmer for D=2 is missing: the step-7 quota of 1 + 1 was already exhausted, so no same-round independent review layer could confirm it; it is registered under the §5.1 non-self-healing tier and escalated, and the user accepted it on 2026-11-13 with one condition - if a third occurrence of the same class (D=3) happens, the §5.1 accumulated ceiling forces the slice to stop.

B-14  docs/validation/strategic-review-2026-10-30_EN.md:49
（同 §5.2 引用块 Q-1；该行可见文本为 "Run independent security counter-examples on the same acceptance chain"，工具同时将其计入 `independent review` 命中清单）

B-15  docs/validation/sw-directive-switchover_EN.md:10
4. This slice had the **master author product artifacts directly** - a bootstrap-period exception to the role-readiness precheck (**anchor = directive Appendix C.2**; the waived rule is **§3.2's "no original product artifacts"**), and it does **not** waive the independent review.

B-15b  docs/validation/sw-directive-switchover_EN.md:34
| ⑦ Independent final review | independent reviewer | `independent-reviewer` + `GPT-5.3-Codex (copilot)` | §7.2 |

B-15c  docs/validation/sw-directive-switchover_EN.md:140
- **⑦ is never optional** (R5.3), and this slice **amends the directive itself** (self-referential revision) ⇒ external independent review is mandatory.

B-15d  docs/validation/sw-directive-switchover_EN.md:197
| independent review | `GPT-5.3-Codex (copilot)` | ⑦ | 2026-09-21 | independent final review (first round blind) | — |

B-16  docs/validation/t011-chain_EN.md:5
Date: 2026-09-07. Verdict: `T011 COMPLETE`. This task delivers candidate freeze, independent review, promotion receipt binding, and task acceptance wiring. After review, two defects were fixed and all gates were re-validated.

B-17  docs/validation/t017-baseline-review_EN.md:1
# T017 Verifier Baseline Independent Review Record

B-17b  docs/validation/t017-baseline-review_EN.md:14
See also the matching [DeepSeek V4 Pro independent review record](t017-deepseek-baseline-review_EN.md). The two records corroborate each other but do not constitute a signature or attestation.

B-18  docs/validation/t017-ci-remediation-codex-review_EN.md:1
# T017 CI Remediation GPT-5.3 Codex Independent Review Record

B-19  docs/validation/t017-ci-remediation-deepseek-review_EN.md:1
# T017 CI Remediation DeepSeek Independent Review Record

B-20  docs/validation/t017-deepseek-baseline-review_EN.md:1
# T017 Verifier Baseline DeepSeek Independent Review Record

B-20b  docs/validation/t017-deepseek-baseline-review_EN.md:10
2. This review applied the same strict read-only requirements as the [GPT-5.3 Codex verifier-baseline independent review](t017-baseline-review_EN.md) and returned the same `PASS FOR BASELINE REVIEW` verdict.

B-20c  docs/validation/t017-deepseek-baseline-review_EN.md:22
	- Mitigation already implemented and recorded: digests live in independent constants, automatic refresh from the workflow under test is prohibited, and the baseline requires independent review before landing. This remains non-blocking.

B-20d  docs/validation/t017-deepseek-baseline-review_EN.md:61
- **A. Complete-script gate: passed.** Contract and mutation tests share `decodeWorkflow` (`workflow_test.go:79`) and `validateWorkflowContract` (`workflow_test.go:147`). Twelve digests are independent constants (`workflow_test.go:241`). Comments, here-strings, uninvoked script blocks, unreachable branches, `exit 0`, variable substitution, function shadowing, and forged oracle writes are rejected. Digests are also bound to job/step name, order, `shell`, `if`, and working directory. The documentation accurately limits the “change both and become green” boundary through independent review.

B-20e  docs/validation/t017-deepseek-baseline-review_EN.md:64
- **D. Verifier and release evidence: passed.** `VerifySelfHost` binds seed/candidate/oracle ownership and content. realpath, symlink, directory overlap, path traversal, trailing JSON, and duplicate cases fail closed. SHA256SUMS completely covers artifacts and rejects extras, omissions, and tampering. The SBOM binds the source commit and matches the license manifest module by module. `PublisherIdentityAuthenticated: false` is constant output (`verify.go:202`). The bootstrap lacks the new verifier, and the process explicitly requires independent review and landing before dispatch against a different candidate.

B-21  docs/validation/t017-node24-actions-codex-review_EN.md:1
# T017 Node 24 Actions GPT-5.3 Codex Independent Review Record

B-22  docs/validation/t017-node24-actions-deepseek-review_EN.md:1
# T017 Node 24 Actions DeepSeek V4 Pro Independent Review Record

B-23  docs/validation/t027-b1-candidate-live-discovery_EN.md:126
- **Recommendation (not implemented here)**: make the CLI's minimum credit floor a candidate/platform **constant** (pass `max(30, maximumRequests)` and account the "product request budget" and "CLI credits" separately), or distinguish the two explicitly at the probe entry point. That is a production behaviour change and needs **its own slice, contract-first work and independent review**; this slice only records it.

B-24  docs/validation/t027-b3c-contracts-durability_EN.md:50
7. **Upgrade unlock boundary**: after a (platform × premise) reaches `proven`, only that platform, within the deployment environments covered by that premise, has the corresponding refusals lifted; **nothing** is unlocked for any other premise or platform, **AT-23 is not waived** (it must still pass), **independent review and protocol-first revision are not waived**, and no gate beyond `first-dispatch` is relaxed. A platform-tier upgrade requires a protocol-first amendment of this clause, traceable evidence, and an independent review.

B-24b  docs/validation/t027-b3c-contracts-durability_EN.md:76
- **Coverage limits of the mechanical checks**: the three items listed as "not mechanically verifiable" — "the evidence pack is recomputable from rounds and the journal" (needs the pack to be re-run), "the premise covers the target deployment environment" (needs a deployment baseline), "independent review is sufficient" (needs the review-process record) — are **not** claimed as verified by this slice.

B-24c  docs/validation/t027-b3c-contracts-durability_EN.md:114
**All of U1–U14 adopted as recommended**, notably: the clause sits as its own paragraph after the replay-store paragraph (U1); granularity is platform × premise (U2); the discrimination gate and mechanism evidence become criteria (U3); independent review is positioned as an **upgrade governance precondition**, not a tier criterion (U4); "the evidence pack is recomputable" folds into criterion (a) (U5); the `disproven` trigger rule is added (U6); the runtime mapping is stated (U7); **no schema change** (U8); the ADR keeps only the ledger and the criteria have a single authority in CONTRACTS (U10); a new validation report (U11); bilingual ledger write-back (U12); ④ keeps its final review (U13); the tiered anchor lands in two places only (U14).

B-25  docs/validation/t027-terminal-routing_EN.md:58
Conclusion: every Medium+ finding was fixed and re-reviewed by the independent reviewer (Codex); nothing was self-certified.

B-26  docs/validation/t027b-tactical-blueprint_EN.md:26
- **Contract blocker**: `CONTRACT-BLOCK` - upgrading Windows publication durability requires revising `docs/CONTRACTS.md` section 7 with an independent review first.
```

### 5.5 补充关键词 `独立评审`（任务词表之外，自行增加）—— 工具报告：34 matches / 19 files

| 文件 | 命中行号 | 原文载体 |
|---|---|---|
| `docs/validation/directive-v1.12.md` | 1, 10, 47, 48, 51, 52 | C-1 |
| `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md` | 10, 13, 58 | C-2 |
| `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md` | 8, 11, 63 | C-3 |
| `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2b.md` | 12, 15 | C-4 |
| `docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md` | 61 | C-5 |
| `docs/validation/evidence/P73-T027B-ADMISSION/01-arch/design-r1.md` | 100, 141 | C-6 |
| `docs/validation/evidence/P73-T027B-ADMISSION/02-batch/batch-02-rulings-gates.md` | 153 | C-7 |
| `docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/01-product/review.md` | 55 | C-8 |
| `docs/validation/evidence/STRATEGIC-REVIEW-2026-11-02/03-round3/review.md` | 17 | C-9 |
| `docs/validation/s1-exit.md` | 59 | C-10 |
| `docs/validation/s1-selfhost.md` | 84 | C-11 |
| `docs/validation/strategic-review-2026-10-30.md` | 1 | C-12 |
| `docs/validation/strategic-review-2026-11-02.md` | 1 | C-13 |
| `docs/validation/t011-chain.md` | 5 | C-14 |
| `docs/validation/t017-baseline-review.md` | 18 | C-15 |
| `docs/validation/t017-deepseek-baseline-review.md` | 22, 61, 64 | C-16 |
| `docs/validation/t027-b3c-contracts-durability.md` | 50, 76, 112 | C-17 |
| `docs/validation/t027b-admission.md` | 69 | C-18 |
| `docs/validation/t027b-tactical-blueprint.md` | 26 | C-19 |

```
C-1  docs/validation/directive-v1.12.md:1
# T027 · DIRECTIVE-V1.12 准则修订验证报告（三轮第三方独立评审发现处置 + 用户裁决落地）

C-1b  docs/validation/directive-v1.12.md:10
| 处置规模 | 三轮第三方独立评审累计 **64 条**发现（**1 Critical + 11 High + 35 Medium + 17 Low**）逐条处置 + A3 与 H4 两项用户裁决落地（64 条 = 第 2 轮 28 条 + 第 3 轮 14 条 + 提交前 22 条，三者均已含于 64）；另处置**处置后复验轮**的 **12 条**、**第四轮复评（R6）**的 **10 条**与**第五轮复评（R7）**的 **9 条**（见 §3.5 至 §3.9） |

C-1c  docs/validation/directive-v1.12.md:47
| 提交前评审（v1.12） | 第三方独立评审 | `DIRECTIVE REVIEW: FINDINGS`：1 Critical（C1）+ 4 High + 13 Medium + 4 Low |

C-1d  docs/validation/directive-v1.12.md:48
| 处置后复验轮 | 第三方独立评审 | C1 复核 19/19 到位；本轮 1 High + 5 Medium + 6 Low，按用户裁决处置（H1 与 M1 – M5 落地，L1 – L6 留待下一片） |

C-1e  docs/validation/directive-v1.12.md:51
| 第四轮复评（R6） | 第三方独立评审 | `DIRECTIVE REVIEW: FINDINGS`：1 High + 4 Medium + 5 Low；N1 与 N2 由用户直接裁定、N3 与 N4 与 N6 与 N10 台账与文案类已落地、N5 实测不可复现、N7 与 N8 与 N9 登记 OB-34 至 OB-36 |

C-1f  docs/validation/directive-v1.12.md:52
| 第五轮复评（R7） | 第三方独立评审 | `DIRECTIVE REVIEW: FINDINGS`：1 High + 4 Medium + 4 Low；P1 与 P2 与 P3 按用户裁决落地（含规则语义变更）、P4 与 P5 登记 OB-37 与 OB-38、P6 至 P9 台账与文案类已落地 |

C-2  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md:10
> （新增）「如某（平台 × 前提）经本节「C6 判定路径」完成协议先行修订并经独立评审后激活预建槽位实现，则该平台记录文件的发布可见性以槽内头标志提交为准，其目录项在槽创建时一次性建立、此后永不重建或替换；未激活的平台与激活前的实现维持本段既有口径逐字不变。」

C-2b  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md:13
〔该行全文超过 1500 字符；此处摘录含「独立评审」的句段，截断点以〔……〕标注〕
> ……（前略）该替换必须先修订本节并经独立评审后方可进入注入取证。……（中略，含 C6 逐项取证 (a)(b)(c) 与平台等级分层、非 Windows 零影响等句）……升级仅在该平台、该前提覆盖的部署环境内解除对应拒绝，且必须 协议先行修订、证据可追溯并通过独立评审）……

C-2c  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md:58
- **升级路径**：三判据满足 ⇒ 仅 C6 候选级结论；平台级升级仍需协议先行（本片 §7 修订即先行文本）、证据可追溯、独立评审、AT-23 不豁免。C2 `proven`（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）不被本片改写或升级。

C-3  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md:8
> （新增）「如某（平台 × 前提）经本节「C6 判定路径」完成协议先行修订并经独立评审后激活预建槽位实现，则该平台**仅 R/C 记录文件**的发布可见性以槽内头标志提交为准，其目录项在首次派发前由装配层预建一次性建立、此后永不重建或替换；未激活的平台与激活前的实现维持本段既有口径逐字不变。」

C-3b  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md:11
〔该行全文超过 3400 字符；含「独立评审」的句段为：〕
> ……（前略）（前置条件）C6 从结构性候选升为可注入候选的前置条件：协议层完成仲裁原语替换——槽文件以 no-replace 一次性建立且永不删除、替换或追加，槽内头标志只允许 0→1 单向翻转并以槽级排他锁仲裁并发写者；该替换必须先修订本节并经独立评审后方可进入注入取证。……

C-3c  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md:63
- **升级路径**：三判据满足 ⇒ 仅 C6 候选级结论；平台级升级仍需协议先行（本片 §7 修订即先行文本）、证据可追溯、独立评审、AT-23 不豁免；C2 `proven`（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）不被改写或升级。

C-4  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2b.md:12
> （新增）「如某（平台 × 前提）经本节「C6 判定路径」完成协议先行修订并经独立评审后激活预建槽位实现，则该平台**仅 R/C 记录文件**的发布可见性以槽内头标志提交为准，其目录项在首次派发前由装配层预建一次性建立、此后永不重建或替换；未激活的平台与激活前的实现维持本段既有口径逐字不变。」

C-4b  docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2b.md:15
〔该行全文超过 3800 字符；含「独立评审」的句段同 C-3b（r2b 为 r2 的修订版，含红队 §3 High 改动点①），含句：〕
> ……（前略）该替换必须先修订本节并经独立评审后方可进入注入取证。……

C-5  docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md:61
CONTRACT-BLOCK: Windows 平台发布耐久升级须先修订 §7 并经独立评审，本批不得仅凭机制或探针成功解锁真实生产链。

C-6  docs/validation/evidence/P73-T027B-ADMISSION/01-arch/design-r1.md:100
`docs/validation/s1-exit.md` RFC §14 门表九行：BLOCKED 两条（核心任务链与 fail-close 场景、产品所有者最终 S1 批准）、PASS 四条（seed 自托管与真实平台、whois 只读影子、generic/C/Go harness、独立评审）、PARTIAL 三条（正式安装包教程、安装/升级/回滚/卸载实测、支持矩阵与发行说明）。AT-23 所在行（`:43`）= 阶段 1 条件闭合达成；AT-24 所在行（`:44`）= 条件闭合（真实 CLI 端到端仍依赖 T027-B）。

C-6b  docs/validation/evidence/P73-T027B-ADMISSION/01-arch/design-r1.md:141
- **B 段（选项 ②）所需且缺失**：Sandbox 接口契约段（CONTRACTS，协议先行）；调用点接线证据（若裁方案 ②）；死代码禁令处置证据（B2 §8 第 8 条修订或裁定豁免的留痕）；⑦ 独立评审记录。

C-7  docs/validation/evidence/P73-T027B-ADMISSION/02-batch/batch-02-rulings-gates.md:153
| 独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审 |

C-8  docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/01-product/review.md:55
| 05 | 保留；独立评审及代码/文档整体接受仍适用，none：[策略:35](docs/TEST_STRATEGY.md#L35) | 06 | 保留；能力不足拒绝不等于已具备隔离能力，none：[策略:36](docs/TEST_STRATEGY.md#L36) |

C-9  docs/validation/evidence/STRATEGIC-REVIEW-2026-11-02/03-round3/review.md:17
| 必须做 | 8（发布前） | **CONTRACTS-S7-DURABILITY** | Windows 耐久仍为 unproven，发布等级与失败口径须先明确；理由锚点：GAP-26/GAP-34 / 断点 A | P46 排期裁定；保持 v35 的契约/Schema 边界及独立评审额度 1+1 | 契约修订通过不等于 Windows 耐久已证 |

C-10  docs/validation/s1-exit.md:59
| 独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审 |

C-11  docs/validation/s1-selfhost.md:84
维护者须从已独立评审的 commit/fixture 计算并手工填写以下必填输入，不得直接信任候选生成值；workflow 必须从 `main` 分支启动，其他 ref 的 self-host job 会被拒绝：

C-12  docs/validation/strategic-review-2026-10-30.md:1
# 证轨 P35 战略回看报告：五轮独立评审（准则 v1.56）

C-13  docs/validation/strategic-review-2026-11-02.md:1
# 证轨 P46 战略回看报告：三轮独立评审与补偿复核（准则 v1.60）

C-14  docs/validation/t011-chain.md:5
日期：2026-09-07。结论：`T011 COMPLETE`。本任务完成候选冻结、独立评审、发布收据绑定与 task 接受接线；并在复核中修复 2 个问题后再次通过门禁。

C-15  docs/validation/t017-baseline-review.md:18
1. `MEDIUM`：完整脚本 SHA-256 是内容漂移门禁，不是 PowerShell 语义证明。若 workflow 与测试中的期望摘要在同一变更中被恶意同步修改，本地测试可以重新变绿；控制措施是独立评审 baseline，禁止从待验 workflow 自动刷新摘要。后续可考虑将摘要基线外置为受保护只读工件。

C-16  docs/validation/t017-deepseek-baseline-review.md:22
	- 缓解（已落实并归档）：摘要来自独立常量、禁止从待验 workflow 自动刷新、baseline 必须独立评审落库。维持不阻断。

C-16b  docs/validation/t017-deepseek-baseline-review.md:61
- **A 完整脚本门禁：通过。** 契约与全部变异测试共用 `decodeWorkflow`（`workflow_test.go:79`）/ `validateWorkflowContract`（`workflow_test.go:147`）；12 个摘要为独立常量（`workflow_test.go:241`）；注释、here-string、未调用 scriptblock、不可达分支、`exit 0`、变量替换、函数遮蔽、伪造 oracle 写入均被拒绝；摘要同时绑定 job/步骤名、顺序、`shell`、`if`、工作目录；“同改即变绿”边界已由文档诚实限定为独立评审。

C-16c  docs/validation/t017-deepseek-baseline-review.md:64
- **D verifier 与发布证据：通过。** `VerifySelfHost` 绑定 seed/candidate/oracle 所有权与内容；realpath、symlink、目录重叠、路径穿越、trailing JSON、重复 case 均 fail-closed；SHA256SUMS 完整覆盖并拒绝多余/缺失/篡改；SBOM 绑定 source commit 且与许可证清单逐模块一致；`PublisherIdentityAuthenticated: false` 为常量输出（`verify.go:202`）；bootstrap 不含新 verifier，流程明确要求先独立评审落库再 dispatch 不同候选。

C-17  docs/validation/t027-b3c-contracts-durability.md:50
7. **升级的解锁边界**：某（平台 × 前提）升为 `proven` 后，**仅在该平台、该前提覆盖的部署环境内**解除对应拒绝；**不解锁**任何其它前提或平台，**不豁免 AT-23**（仍必须通过）、**不豁免独立评审与协议先行**、也不放宽 `first-dispatch` 之外的任何门禁。平台等级升级必须协议先行修订本节、证据可追溯并通过独立评审。

C-17b  docs/validation/t027-b3c-contracts-durability.md:76
- **机械检查的覆盖边界**：§"当前不可机械验证"三项——"证据包可由轮次与 journal 复算"（需对证据包实跑）、"前提覆盖目标部署环境"（需部署清单）、"独立评审充分性"（需审查流程记录）——本片均**未**声称已验证。

C-17c  docs/validation/t027-b3c-contracts-durability.md:112
裁定 **U1–U14 全按推荐采纳**，其中关键项：条款落点为 replay store 段后的独立段（U1）；判定粒度 = 平台 × 前提（U2）；区分力门与机制证据入判据（U3）；独立评审定位为**升级治理前置**而非等级判据（U4）；"证据包可复现"并入判据 (a)（U5）；补 `disproven` 触发规则（U6）；写明运行期映射（U7）；**不改 schema**（U8）；ADR 只记账本、判据单一权威在 CONTRACTS（U10）；新增本验证报告（U11）；双语台账回写（U12）；④ 保留终审（U13）；分级锚定句只落两处（U14）。

C-18  docs/validation/t027b-admission.md:69
- **缺失**：Sandbox 契约段、调用点接线证据、死代码禁令处置留痕（已于本片落地）、实施片 ⑦ 独立评审记录；S-G1 / S-G6 / S-G7 不因 B 段获得任何解锁证据。

C-19  docs/validation/t027b-tactical-blueprint.md:26
- **契约阻塞**：`CONTRACT-BLOCK`——Windows 平台发布耐久升级须先修订 `docs/CONTRACTS.md` §7 并经独立评审。
```

### 5.6 命中文件的汇总事实（按文件名/标题观察，不含判定）

- `docs/validation/` 内存在多份**文件名或标题自题为「独立评审 / Independent Review Record」**的文件：`t017-baseline-review{,_EN}.md`、`t017-deepseek-baseline-review{,_EN}.md`、`t017-ci-remediation-{codex,deepseek}-review{,_EN}.md`、`t017-node24-actions-{codex,deepseek}-review{,_EN}.md`；`B3A-CALIBRATION_EN.md:127`（「④ Independent review record (Codex, read-only)」）；`strategic-review-2026-10-30{,_EN}.md`（标题：五轮独立评审）、`strategic-review-2026-11-02{,_EN}.md`（标题：三轮独立评审与补偿复核）；`s1-exit{,_EN}.md:59` 门表「独立评审 / Independent review」行（值 PASS）。
- **未检索到**文件名或标题含「独立安全评审 / independent security review」字样的文件（依据：`independent security` 在 `docs/validation/**` 仅 1 处命中且为建议条目文本；`独立安全评审` 在 `docs/` 全域的命中见 §1.3，无一为 validation 目录下的报告文件名/标题）。
- 中文词「安全评审」在 `docs/validation/**` 的命中仅两类：`evidence/P40-T027B-LOCK-REAPER/05-prereview/review.md:26`（P40 证据包预审文件内对偏差 A–G 的安全评审列举）与 P76 本片自产的 `00-briefing.md` / `input-list.md`。

## 六、`acceptance-security-counterexamples{,_EN}` 的性质说明（只描述文档自述内容与出处，不下结论）

以下均为两文件自述内容（CN/EN 同构）：

- 标题：`# ACCEPTANCE-SECURITY-COUNTEREXAMPLES（P52）验收报告`（CN:1）/ `# ACCEPTANCE-SECURITY-COUNTEREXAMPLES (P52) acceptance report`（EN:1）。
- 自述定位（CN:3 / EN:3）：`切片：ACCEPTANCE-SECURITY-COUNTEREXAMPLES（安全验证类 [SLICE]，M 档）。封片时点 2026-11-03（HEAD = 941f60e，准则 v1.60）。主控模型 = DeepSeek V4 Pro。`
- 自述目标（CN §1 / EN §1）：「针对断点 A 的接受链缺口，落地 GAP-51/52/55/60/64/70 所代表的安全反例验证，证明系统面对非法输入、路径伪造、句柄越权、边界绕过时正确 fail-closed」；「本片只做测试用例与断言，绝不修改生产代码……不改接受语义，不改 `docs/CONTRACTS{,_EN}.md`」；「反例测试 PASS 即证明攻击被成功阻断；任何反例 FAIL（攻击未被阻断）属真实漏洞，立即停止并上报」。
- 自述击穿发现（CN §2 / EN §2，标注 2026-11-03）：「T51-1 中间符号链接逃逸在 Linux（go1.22.12）实测击穿……按铁律 3 立即停止、载荷归档（`05-attack/t51-1-linux-escape.md`）、登记 GAP-77。本片不修生产代码。」§2b 列出 GAP-77（HIGH，已实证击穿）/GAP-78/79（LOW）/GAP-80（观察）四项疑似漏洞上报；§7 说明「GAP-77 修复留待 P53」。
- 自述审查链（CN §4 / EN §4）：① 架构 DESIGN DONE（`01-arch/design.md`）→ ③ 测试先行 PROBE: ALL-GREEN（`03-tests/red-baseline.md`）→ ② 实施 ②a–②d 四批串行全绿（`04-impl/implementer-reports.md`）→ ⑤ 预审 PRE-REVIEW: PASS（0 Medium+；`05-prereview/review.md`）→ ⑥ 独立扫描 INDEPENDENT SCAN: PASS（§6 逐字引用，`06-scan/review.md`）→ ⑦ 终审 FINAL REVIEW: PASS（0 Medium+；`07-final/review.md`）。
- 自述术语：该报告的第六阶段称「⑥ 独立扫描 / Stage-6 independent scan」，围栏内结论行为 `INDEPENDENT SCAN: PASS`。
- 证据包（CN §8 artifacts 围栏，self-listed）：`docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/` 共 13 个条目——`00-briefing.md`、`01-arch/design.md`、`03-tests/red-baseline.md`、`04-impl/implementer-reports.md`、`05-attack/t51-1-linux-escape.md`、`05-prereview/review.md`、`06-scan/review.md`、`07-final/review.md`、`08-gates/gate-readings.txt`、`09-glm/glm-1-bounddir-boundread-inventory.md`、`09-glm/glm-2-gates-chain-inventory.md`、`MANIFEST.md`、`SHA256SUMS.txt`、`freedom-list.md`、`input-list.md`（fence 内另含 repo 级测试文件与文档条目）。
- 关键词事实：两文件均未命中 `安全评审`、`independent review`、`security review`、`independent security`（不在 §五命中清单内）。
- 出处关联（仅并列陈述）：`evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/input-list.md:22` 将其描述为「验收安全反例（Q7 独立安全评审历史状态核实对象之一）」；`strategic-review-2026-10-30_EN.md:49` 的建议条目名为 `ACCEPTANCE-TRUST-NEGATIVES`，而实际存在的报告标题为 `ACCEPTANCE-SECURITY-COUNTEREXAMPLES (P52)`——两个名称同时存在于仓库文本中，本报告不判定二者是否同一事项。

## 七、`SECURITY{,_EN}` 中与独立评审相关原文

### 7.1 中文 `docs/SECURITY.md`

```
docs/SECURITY.md:5
日期：2026-09-07；安全设计基线，未完成独立安全签署。本文是权限、威胁模型与硬安全门禁的规范性权威；[项目建议书](RFC-proofrail-unattended-ai-engineering-product.md) §12、§16.7-16.8 保留风险与设计来源。需求 REQ-004/005/007/013/019/021-026/028。

docs/SECURITY.md:15
| chain/core | 只读捕获 | 物化/恢复，先停机 | 读取校验 | 持锁写入 | 校验独立授权后发布 |

docs/SECURITY.md:21
限制必须由 runner/OS 能力落实；不能仅靠提示词。未知进程或未能证明相关写者停止时不恢复、不交接、不发布。执行者不得拥有修改策略和独立批准同一变更的合并权限。

docs/SECURITY.md:35
| SEC-08 供应链 | 候选覆盖 seed/verifier、候选自证发布、同 runner 后台进程篡改、CI Action 漂移、校验和来源被冒充 | candidate probe 与可信核验使用独立 runner、独立较早 verifier commit、固定 candidate/bootstrap commit、Action 完整 SHA、GitHub CI、外部 oracle、SHA256SUMS/SBOM/许可清单与人工批准；release/harness | AT-13/14 |

docs/SECURITY.md:41
| SEC-14 篡改/提权/抵赖 | CLI Agent 二进制/配置替换、逃逸 source/store、提升工具/网络权限、伪造或缺失事件、停机/恢复歧义、自行 commit/push/publish | 可执行文件与配置摘要、能力探针、隔离 workspace、进程树监督、事件/日志与前后 manifest、独立 gates/review；adapters/guard/evidence | AT-23 |

docs/SECURITY.md:42
| SEC-15 信息泄露/抵赖/副作用 | visible 黑箱 Agent 读取越界资源、泄密、网络/费用/后台进程不可见、外部 Git/发布副作用，或风险说明被当作免责放行 | 显式版本化风险确认、隔离 workspace、用户监督与宿主侧限制、归还后全量扫描、unknown 标记、独立 gates/review；console/chain/guard/evidence | AT-24 |

docs/SECURITY.md:64
以下事项需要独立评审：扩大写入/网络/工具权限、变更审批规则、生成 hook、canonical/签名/恢复协议。黑箱模式的风险确认不是免责通行证：ProofRail 仍须履行其承诺的隔离、扫描、门禁、证据和评审责任，用户只承担已明确列出的宿主监督与不可观测外部行为风险。未修复的 fail-close 绕过、失败候选接受、秘密外泄或不可恢复覆盖均阻断发布。
```

### 7.2 英文 `docs/SECURITY_EN.md`

```
docs/SECURITY_EN.md:5
Date: 2026-09-07; security-design baseline without independent security sign-off. This document is the normative authority for permissions, threat modeling, and hard security gates. [Project proposal](RFC-proofrail-unattended-ai-engineering-product.md) sections 12 and 16.7-16.8 retain risk and design provenance. Requirements: REQ-004/005/007/013/019/021-026/028.

docs/SECURITY_EN.md:15
| chain/core | Read-only capture | Materialize/restore after stopping writers | Read/validate | Locked writes | Only after independent authorization |

docs/SECURITY_EN.md:21
Runners/OS capabilities enforce restrictions, not prompts alone. Unknown processes or uncertain writer termination block restore, handoff and publication. An executor cannot combine policy mutation with independent approval of the same change.

docs/SECURITY_EN.md:35
| SEC-08 Supply chain | Candidate overwrites seed/verifier, self-certifies release, tampers through a same-runner background process, drifts CI Actions, or spoofs checksum provenance | Separate runners for candidate probing and trusted verification, an independently reviewed earlier verifier commit, fixed candidate/bootstrap commits, full Action SHAs, GitHub CI, external oracle, SHA256SUMS/SBOM/license evidence, and human approval; release/harness | AT-13/14 |

docs/SECURITY_EN.md:41
| SEC-14 Tampering/elevation/repudiation | Substituted CLI Agent binary/config, source/store escape, tool/network escalation, forged/missing events, ambiguous stop/resume, or self commit/push/publish | Executable/config digests, capability probe, isolated workspace, process-tree supervision, events/logs and before/after manifests, independent gates/review; adapters/guard/evidence | AT-23 |

docs/SECURITY_EN.md:42
| SEC-15 Disclosure/repudiation/effects | A visible black-box Agent reads out-of-scope resources, leaks secrets, hides network/cost/background work, performs external Git/publication effects, or treats risk notice as a waiver | Explicit versioned risk acknowledgment, isolated workspace, user supervision/host controls, full post-return scan, unknown labels, independent gates/review; console/chain/guard/evidence | AT-24 |

docs/SECURITY_EN.md:64
Independent review is required for broader write/network/tool authority, approval changes, generated hooks and canonical/signing/recovery protocols. Black-box risk acknowledgment is not a liability waiver: ProofRail still owns its promised isolation, scanning, gates, evidence and review, while the user owns only the stated host-supervision and unobservable external-behavior risks. Unfixed fail-close bypass, acceptance of failed candidates, secret leakage or unrecoverable overwrite blocks release.
```

## 八、检索方法与覆盖范围说明

### 8.1 检索关键词与命中范围

| 关键词 | 检索起点 | 工具报告计数 | 命中文件数 |
|---|---|---|---|
| `security review` | `docs/validation/**/*.md`（递归，含 `evidence/`） | 1 match | 1 |
| `independent security` | 同上 | 1 match | 1 |
| `安全评审` | 同上 | 8 matches | 3 |
| `independent review` | 同上 | 49 matches | 26（工具计数）；本报告按工具返回清单逐文件枚举为 27 个文件，计数口径差异见 8.4 |
| `独立评审`（补充词） | 同上 | 34 matches | 19 |
| `独立安全评审`（补充词） | `docs/**/*.md`（全域定位 §7 权威出处） | 11 matches | 6（见 §1.3） |
| `独立`（补充词，CN SECURITY） | `docs/SECURITY.md` | 7 matches | 1 |
| `independent`（补充词，EN SECURITY） | `docs/SECURITY{,_EN}.md` | 9 matches | 2 |

其它定位性检索（服务 §一–§四 的枚举，非任务 B 关键词）：`REQ-007|REQ-009|REQ-020|AT-13|§7`、`AT-13|AT-23|T027|T026|进入 S1|before S1`（PRODUCT_REQUIREMENTS*）；`AgentRunner|CLI Agent|T026|T027|AT-23`（S1_SUPPORT_MATRIX*）；`Release|发布`、`T026|T027|AT-23|AgentRunner`（INSTALLATION*、INSTALLATION_PLAN*）。

### 8.2 检索起点与跳过范围

- 任务 A 起点：`d:\LZProjects\prfrail\docs\` 下 4 个文档族（PRODUCT_REQUIREMENTS{,_EN}.md、S1_SUPPORT_MATRIX{,_EN}.md、INSTALLATION{,_EN}.md、INSTALLATION_PLAN{,_EN}.md），共 8 个文件。
- 任务 B 起点：`d:\LZProjects\prfrail\docs\validation\`，递归包含 `evidence\` 子目录——任务书允许跳过 `evidence/`，本次**未跳过**（全量参与关键词检索），以减少漏检。
- 命中文件名枚举（file_search）：`docs/validation/evidence/**/*review*.md` 工具报告 118 个结果（仅回显前 21 个）；`docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/**/*.md` 返回 13 个文件（全量列出，见 §六）。证据包内评审报告文件名形态为 `05-prereview/review.md`、`06-scan/review.md`、`07-final/review.md`（另见 `01-arch/red-team-sol-review*.md`、`review-r1.md` 等变体）。

### 8.3 工具与保真处理

- 工具：仓库全文检索（大小写不敏感）+ 文件名 glob 检索 + 定点 read_file + 一次性 node 逐行提取脚本（置于 `tmp/`，用后即删）。未执行任何 git 写操作；未修改任何被测文档。
- 检索回显超长行（>500 字符被截断）已全部用磁盘逐行读取取回并全文引用；唯一例外：`evidence/CONTRACTS-S7-DURABILITY/01-arch/design{,-r2,-r2b}.md` 三行（1511/3417/3843 字符）在 §5.5 中以〔……〕标注的句段摘录形式引用，截断点显式标注。
- 含表格竖线 `|` 的原文一律以代码块呈现（避免表格转义改变字节）；表格内原文仅限不含竖线的行。

### 8.4 计数口径差异（事实登记）

- `independent review`：工具报告「49 matches in 26 files」，本报告按工具返回清单逐文件枚举为 27 个文件、50 行；差异未逐项复核成因（可能涉及单行多次命中或跨关键词重复计入，如 `strategic-review-2026-10-30_EN.md:49` 同时出现在 `independent security` 与 `independent review` 命中清单中）。枚举清单以工具实际返回为准。
- 其余关键词的工具计数与逐条枚举一致（`安全评审` 8/3、`独立评审` 34/19、`security review` 1/1、`independent security` 1/1）。

### 8.5 覆盖边界（未做的事）

- 未对 `docs/validation/evidence/` 下除命中文件外的百余个证据文件做全文阅读（仅关键词检索 + 文件名枚举）。
- `docs/validation/` 中文报告（无 `_EN` 后缀者）与英文报告均已参与关键词检索；`独立评审` 一词主要命中中文报告与 evidence 文件，`independent review` 主要命中 `_EN` 报告——两套词面分布互补，均已在 §5.4/§5.5 登记。
- INSTALLATION_EN.md 路径第 2 条（对应 CN:43 的英文行）未单独取回（关键词 `Release` 未命中该行），CN:43 已登记。

数据齐备（行号以磁盘实读为准，四条关键句已用编辑器 grep 交叉核验一致）。以下为逐条事实枚举。

**枚举约定**：行号为 2026-11-17 磁盘实读行号；摘录为逐字原文（仅去除行首缩进），`` `…` `` 内的 `……` 表示**同一行内**省略段（省略号之外的内容连续逐字）；表格单元格内原文的 `|` 以 `\|` 转义；含反引号的摘录用双反引号包裹。未修改任何文件。

## 文件：docs/S1_RELEASE_NOTES.md（共 29 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | 文档标题 | `# ProofRail S1 候选发行说明` |
| 3 | 文档头部 | `[English](S1_RELEASE_NOTES_EN.md)` |
| 5 | 标题下导语 | ``日期：2026-09-09。状态：`DRAFT / NOT RELEASED`。版本、下载链接、正式 artifact 文件名、支持起始日期和批准人均待 S1 exit 冻结；本文件不授权 tag、GitHub Release、发布或安装。`` |
| 7 | §候选能力（标题） | `## 候选能力` |
| 14 | 候选能力 | `- 固定 verifier/candidate commit、双平台 candidate probe/self-host、校验和、SPDX SBOM 与许可证清单的候选发布链。` |
| 18 | §支持与已知限制 | `支持范围以 [S1 候选支持矩阵](S1_SUPPORT_MATRIX.md) 为准。S1 安装模型是[手工解压 Windows 便携 ZIP](INSTALLATION.md)、显式路径运行且不修改 PATH。当前仍无正式发行包，界面是逐行 CLI。正式 AgentRunner/CLI Agent 仍未交付；受监督 visible 黑箱候选的产品侧已由 T028 实现（AT-24 条件闭合），但真实 CLI 端到端依赖 T027-B、尚不可运行发布；该路径即使完成也只保证归还后的产物验收，不保证工具、网络、费用或外部副作用，并须用户显式确认风险和监督。S1 也不提供签名/attestation、Linux 正式支持、自动更新、Web 控制台、完整 TUI 或任意 executable hook。` |
| 20 | §发布前剩余门禁（标题） | `## 发布前剩余门禁` |
| 22 | 发布前剩余门禁 | `1. 对最终 ZIP 复跑核验、首次运行、升级、回滚和卸载。` |
| 23 | 发布前剩余门禁 | `2. 将版本、artifact 文件名、下载位置、EOL 通知方式和支持日期绑定到本说明及支持矩阵。` |
| 24 | 发布前剩余门禁 | `3. 获得独立评审结论和产品所有者明确 S1 exit/发行批准。` |
| 27 | 发布前剩余门禁（切片同步引语） | ``> 切片同步（P55 `EXIT-SEED-BUILD-EVIDENCE`）：本片不适用（无发行物与功能变更，仅证据补链）；理由：本片不产生 tag、artifact 或功能项，只在 `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/` 归档 AT-13 补证，候选能力与剩余门禁逐字不变。`` |
| 29 | 发布前剩余门禁（切片同步引语） | ``> 切片同步（P57 `P5-P6-CLOSEOUT-ASSESSMENT`）：本片不适用（无 tag、artifact 或功能项）；理由：P57 只整理 S1 Exit 剩余门的收口清单（见 `docs/validation/p57-closeout-assessment.md` §4），不改动本文件「发布前剩余门禁」的措辞。`` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。

## 文件：docs/S1_RELEASE_NOTES_EN.md（共 29 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | Title | `# ProofRail S1 Candidate Release Notes` |
| 3 | Header | `[中文](S1_RELEASE_NOTES.md)` |
| 5 | Intro under title | ``Date: 2026-09-09. Status: `DRAFT / NOT RELEASED`. Version, download links, formal artifact filenames, support start date, and approver remain pending S1 exit. This file does not authorize a tag, GitHub Release, publication, or installation.`` |
| 7 | §Candidate Capabilities (heading) | `## Candidate Capabilities` |
| 14 | Candidate Capabilities | `- Candidate release chain with fixed verifier/candidate commits, dual-platform candidate probe/self-host, checksums, SPDX SBOM, and license manifest.` |
| 18 | §Support and Known Limitations | `See the [S1 candidate support matrix](S1_SUPPORT_MATRIX_EN.md). The S1 installation model is a [manually extracted Windows portable ZIP](INSTALLATION_EN.md), invoked by explicit path without modifying PATH. No formal release package exists and the current interface is line-oriented CLI. Formal AgentRunner/CLI Agent are still undelivered; the supervised visible black-box candidate's product side is implemented by T028 (AT-24 conditionally closed) but its real CLI end-to-end depends on T027-B and it is not yet a runnable release. Even when implemented, the latter warrants post-return artifact acceptance only, not tools, network, cost or external effects, and requires explicit user risk acknowledgment and supervision. S1 also excludes signing/attestation, production Linux support, automatic updates, a Web console, the full TUI, and arbitrary executable hooks.` |
| 20 | §Remaining Pre-Release Gates (heading) | `## Remaining Pre-Release Gates` |
| 22 | Remaining Pre-Release Gates | `1. Rerun verification, first run, upgrade, rollback, and uninstall against the final ZIP.` |
| 23 | Remaining Pre-Release Gates | `2. Bind version, artifact filenames, download location, EOL notification channel, and support date to these notes and the support matrix.` |
| 24 | Remaining Pre-Release Gates | `3. Obtain the independent-review verdict and explicit product-owner S1 exit/release approval.` |
| 27 | Remaining Pre-Release Gates (slice-sync quote) | ``> Slice sync (P55 `EXIT-SEED-BUILD-EVIDENCE`): not applicable (no release artifact or feature change, evidence supplement only); reason: this slice produces no tag, artifact or feature item and only archives the AT-13 supplement under `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/`, so the candidate capabilities and remaining gates stay unchanged.`` |
| 29 | Remaining Pre-Release Gates (slice-sync quote) | ``> Slice sync (P57 `P5-P6-CLOSEOUT-ASSESSMENT`): not applicable (no tag, artifact or feature item); reason: P57 only consolidates the S1 Exit remaining-gate closeout inventory (see `docs/validation/p57-closeout-assessment_EN.md` §4) and leaves this file's "gates remaining before release" wording untouched.`` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。

## 文件：docs/INSTALLATION.md（共 206 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | 文档标题 | `# ProofRail S1 Windows 便携 ZIP 安装指南` |
| 5 | 标题下导语 | ``日期：2026-09-10。S1 安装模型已冻结为：Windows 11 amd64 便携 ZIP、手工解压到用户选择的独立目录、始终用 `prfrail.exe` 的显式路径运行，不修改 `PATH`、注册表或系统目录。ZIP 解压后只有一个 `prfrail/` 顶层目录，四个发行文件均在其中。当前尚无正式 GitHub Release；新目录结构必须通过新的可信候选 artifact 复验后，才能视为与最终发行物结构相同。`` |
| 11 | §1 环境要求与能力层次 | `\| 操作系统/架构 \| Windows 11 amd64（x86-64）；S1 唯一正式支持候选 \| 同核心要求；VS Code 与扩展也必须在该用户会话中运行 \|` |
| 15 | §1 环境要求与能力层次 | `\| 主机软件 \| Windows PowerShell 5.1+ 仅用于本指南的解压/核验命令；浏览器或可选 GitHub CLI 用于下载。仅启用正式 AgentRunner 时，另需 GitHub Copilot CLI 1.0.83（Windows x64 当前固定候选） \| VS Code 1.82+、已登录且可用的 GitHub Copilot Chat、SessionBridge 0.1.1 扩展；不要求安装 GitHub Copilot CLI \|` |
| 17 | §1 环境要求与能力层次 | ``\| 网络 \| 下载时访问 GitHub；下载和核验完成后，当前 `version`/`validate`/`preview` 可离线运行 \| Copilot Chat 模型调用需要宿主可用网络、账号/订阅和模型权限 \|`` |
| 19 | §1 环境要求与能力层次 | `……该模式必须由用户监督，并在启动前确认 ProofRail 无法证明工具/命令、越界读取、网络、实际费用、全部后台进程、会话连续性及 commit/push/发布等外部副作用。ProofRail 只保证隔离和归还后的产物扫描、独立门禁、证据与评审。` |
| 54 | §1.1 正式 AgentRunner 的 GitHub Copilot CLI 条件依赖 | ``当前 ProofRail 尚未实现 AgentRunner 与正式 CLI Agent adapter；黑箱模式所需的 visible 投递、TUI 风险确认与显式归还流程已由 T028 实现（产品侧），但真实 CLI 端到端仍依赖 T027-B 的 visible 通道接线，尚不可运行发布。……`` |
| 56 | §1.1 正式 AgentRunner 的 GitHub Copilot CLI 条件依赖 | ``当前发行候选没有图形窗口或完整统一 TUI。……S1 优先用“一个终端、短命令、可复制路径”降低学习成本，而不把原型界面冒充已交付功能。`` |
| 58 | §2 下载与离线核验（标题） | `## 2. 下载与离线核验` |
| 60 | §2 下载与离线核验 | `从 [S1 自托管证据](validation/s1-selfhost.md) 选择已通过可信矩阵的 GitHub Actions run。run 页的 `head_sha` 是 verifier commit，不是 candidate commit；必须同时核对证据报告中的 run ID、verifier commit、candidate commit 和全部 job 结论。` |
| 62 | §2 下载与离线核验 | `浏览器方式会保留 ZIP：登录 GitHub，打开 `https://github.com/larsonzh/prfrail/actions/runs/<run-id>`，在页面底部 **Artifacts** 区点击 `prfrail-Windows-<full-candidate-commit>`。下载文件通常名为 `prfrail-Windows-<full-candidate-commit>.zip`。` |
| 64 | §2 下载与离线核验 | `命令行下载会先检测已安装且已认证的 GitHub CLI：可用时，`gh run download` 会把 artifact 自动解压到 `--dir`，不会保留外层 ZIP；下载后四个文件位于 `$DownloadDir\prfrail`，因此跳过下方的 `Expand-Archive`。……API 回退要求把具有该仓库 Actions 读取权限的 token 放入进程环境变量 `GH_TOKEN` 或 `GITHUB_TOKEN`；……SSH 密钥不能认证 GitHub REST API。` |
| 69 | §2 下载与离线核验（脚本） | `$CandidateCommit = '<full-candidate-commit>'` |
| 70 | §2 下载与离线核验（脚本） | `$Artifact = "prfrail-Windows-$CandidateCommit"` |
| 71 | §2 下载与离线核验（脚本） | `$DownloadDir = Join-Path '<user-selected-independent-directory>' $CandidateCommit` |
| 82 | §2 下载与离线核验（脚本） | `gh run download $RunId --repo $Repository --name $Artifact --dir $DownloadDir` |
| 84 | §2 下载与离线核验（脚本） | `throw "Artifact download failed: $Artifact"` |
| 101 | §2 下载与离线核验（脚本） | `$ListUri = "https://api.github.com/repos/$Repository/actions/runs/$RunId/artifacts?name=$Artifact"` |
| 103 | §2 下载与离线核验（脚本） | `$Matches = @($Response.artifacts \| Where-Object { $_.name -eq $Artifact })` |
| 105 | §2 下载与离线核验（脚本） | `throw "Expected exactly one artifact named $Artifact in run $RunId; found $($Matches.Count)"` |
| 108 | §2 下载与离线核验（脚本） | `throw "Artifact has expired: $Artifact"` |
| 112 | §2 下载与离线核验（脚本） | `$Archive = Join-Path $DownloadDir "$Artifact.zip"` |
| 113 | §2 下载与离线核验（脚本） | `$DownloadUri = "https://api.github.com/repos/$Repository/actions/artifacts/$($Matches[0].id)/zip"` |
| 114 | §2 下载与离线核验（脚本） | `Invoke-WebRequest -UseBasicParsing -Uri $DownloadUri -Headers $Headers -OutFile $Archive` |
| 118 | §2 下载与离线核验 | `以下是浏览器或 REST API 下载 ZIP 后的解压路径；`gh` 分支已经自动解压，应跳过本段。将三个值替换为实际下载文件、完整 candidate commit 和当前用户可写的独立目录；若刚执行 REST 分支，则保留已经设置的 `$Archive`，不再执行第一行：` |
| 121 | §2 下载与离线核验（脚本） | `$Archive = (Resolve-Path '.\prfrail-Windows-<full-commit>.zip').Path` |
| 122 | §2 下载与离线核验（脚本） | `$CandidateCommit = '<full-commit>'` |
| 132 | §2 下载与离线核验（脚本） | `foreach ($line in Get-Content -LiteralPath '.\SHA256SUMS') {` |
| 144 | §2 下载与离线核验（脚本） | `if ($version.data.version -ne "0.1.0+$CandidateCommit") {` |
| 145 | §2 下载与离线核验（脚本） | `throw "Candidate version mismatch: $($version.data.version)"` |
| 152 | §2 下载与离线核验 | `预期解压根目录只含 `prfrail/`，该目录只含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json` 和 `licenses.json`。校验和只证明文件与清单一致，不认证发布者身份；run/verifier/candidate 必须从受信发布记录独立核对，包内 `version --json` 必须绑定同一 candidate commit。不要执行平台或 commit 不匹配的二进制。Actions artifact 可能过期，它不是长期正式下载入口。` |
| 158 | §3 五步快速上手 | `1. 按上一节下载、解压并核验四个发行文件。` |
| 195 | §5 升级与回滚 | `2. 按第 1 节把新 ZIP 解压到新的 commit 目录；不得覆盖旧目录。` |
| 206 | §6 卸载 | `正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定最终版本、commit、run 和文件名。` |

零命中主题：无（七类主题均有命中）。
注（关键词命中但未列入）：L9/L50「黑箱候选」、L21/L23/L28/L34/L44/L46/L48（GitHub Copilot CLI 安装/认证/模型选择）、L86/L90/L92/L97/L98/L111/L124/L134/L136（token/环境变量/解压管线代码行）。

## 文件：docs/INSTALLATION_EN.md（共 206 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | Title | `# ProofRail S1 Windows Portable ZIP Installation Guide` |
| 5 | Intro under title | ``Date: 2026-09-10. The frozen S1 installation model is a Windows 11 amd64 portable ZIP, manually extracted to an independent user-selected directory and always invoked through an explicit `prfrail.exe` path. It does not modify `PATH`, the registry, or system directories. The extracted ZIP has one top-level `prfrail/` directory containing all four release files. No formal GitHub Release exists yet; the new layout must pass a fresh trusted-candidate replay before it can be treated as matching the final carrier.`` |
| 11 | §1 Environment Requirements and Capability Levels | `\| OS/architecture \| Windows 11 amd64 (x86-64), the only S1 production-support candidate \| Same core requirement; VS Code and extensions must run in the same user session \|` |
| 15 | §1 Environment Requirements and Capability Levels | `\| Host software \| Windows PowerShell 5.1+ is needed only for this guide's extraction/verification commands; use a browser or optional GitHub CLI to download. Formal AgentRunner execution additionally requires GitHub Copilot CLI 1.0.83, the current pinned Windows x64 candidate \| VS Code 1.82+, an authenticated and usable GitHub Copilot Chat installation, and the SessionBridge 0.1.1 extension; GitHub Copilot CLI is not required \|` |
| 17 | §1 Environment Requirements and Capability Levels | ``\| Network \| GitHub access for download; current `version`/`validate`/`preview` can run offline after download and verification \| Copilot Chat model calls require host networking, account/subscription access, and model entitlement \|`` |
| 19 | §1 Environment Requirements and Capability Levels | `……That mode requires human supervision and prior acknowledgment that ProofRail cannot prove tools/commands, out-of-scope reads, network, actual cost, every background process, session continuity, or external effects such as commit/push/publication. ProofRail warrants isolation and post-return artifact scanning, independent gates, evidence and review only.` |
| 54 | §1.1 Conditional GitHub Copilot CLI Dependency for Formal AgentRunner | `ProofRail has not implemented AgentRunner or a formal CLI Agent adapter; the visible delivery, TUI risk acknowledgment and explicit return flow required by black-box mode are implemented by T028 (product side), but the real CLI end-to-end still depends on T027-B's visible channel wiring and it is not yet a runnable release. ……` |
| 56 | §1.1 Conditional GitHub Copilot CLI Dependency for Formal AgentRunner | `The current release candidate has no graphical window or complete unified TUI. ……` |
| 58 | §2 Download and Offline Verification (heading) | `## 2. Download and Offline Verification` |
| 60 | §2 Download and Offline Verification | `Select a GitHub Actions run that passed the trust matrix in the [S1 self-host evidence](validation/s1-selfhost_EN.md). The run page `head_sha` is the verifier commit, not the candidate commit. Independently check the report's run ID, verifier commit, candidate commit, and all job conclusions.` |
| 62 | §2 Download and Offline Verification | `The browser path preserves the ZIP: sign in to GitHub, open `https://github.com/larsonzh/prfrail/actions/runs/<run-id>`, and click `prfrail-Windows-<full-candidate-commit>` in the **Artifacts** section at the bottom of the page. The downloaded file is normally named `prfrail-Windows-<full-candidate-commit>.zip`.` |
| 64 | §2 Download and Offline Verification | `The command-line path first detects an installed and authenticated GitHub CLI. When available, `gh run download` automatically extracts the artifact into `--dir` and does not retain the outer ZIP; the four files are placed under `$DownloadDir\prfrail`, so skip the `Expand-Archive` block below. …… SSH keys cannot authenticate the GitHub REST API.` |
| 69 | §2 (script) | `$CandidateCommit = '<full-candidate-commit>'` |
| 70 | §2 (script) | `$Artifact = "prfrail-Windows-$CandidateCommit"` |
| 71 | §2 (script) | `$DownloadDir = Join-Path '<user-selected-independent-directory>' $CandidateCommit` |
| 82 | §2 (script) | `gh run download $RunId --repo $Repository --name $Artifact --dir $DownloadDir` |
| 84 | §2 (script) | `throw "Artifact download failed: $Artifact"` |
| 101 | §2 (script) | `$ListUri = "https://api.github.com/repos/$Repository/actions/runs/$RunId/artifacts?name=$Artifact"` |
| 103 | §2 (script) | `$Matches = @($Response.artifacts \| Where-Object { $_.name -eq $Artifact })` |
| 105 | §2 (script) | `throw "Expected exactly one artifact named $Artifact in run $RunId; found $($Matches.Count)"` |
| 108 | §2 (script) | `throw "Artifact has expired: $Artifact"` |
| 112 | §2 (script) | `$Archive = Join-Path $DownloadDir "$Artifact.zip"` |
| 113 | §2 (script) | `$DownloadUri = "https://api.github.com/repos/$Repository/actions/artifacts/$($Matches[0].id)/zip"` |
| 114 | §2 (script) | `Invoke-WebRequest -UseBasicParsing -Uri $DownloadUri -Headers $Headers -OutFile $Archive` |
| 118 | §2 Download and Offline Verification | `The following path applies after downloading the ZIP through the browser or REST API; skip it after the `gh` branch, which already extracted the artifact. Replace the three values with the downloaded file, full candidate commit, and an independent directory writable by the current user. If the REST branch just set `$Archive`, retain that value and skip the first line:` |
| 121 | §2 (script) | `$Archive = (Resolve-Path '.\prfrail-Windows-<full-commit>.zip').Path` |
| 122 | §2 (script) | `$CandidateCommit = '<full-commit>'` |
| 132 | §2 (script) | `foreach ($line in Get-Content -LiteralPath '.\SHA256SUMS') {` |
| 144 | §2 (script) | `if ($version.data.version -ne "0.1.0+$CandidateCommit") {` |
| 145 | §2 (script) | `throw "Candidate version mismatch: $($version.data.version)"` |
| 152 | §2 Download and Offline Verification | `The extraction root must contain only `prfrail/`, which in turn contains only `prfrail.exe`, `SHA256SUMS`, `sbom.spdx.json`, and `licenses.json`. Checksums prove agreement with the manifest, not publisher identity. Independently verify the run/verifier/candidate from a trusted record, and require package `version --json` to bind the same candidate commit. Never execute a binary for a different platform or commit. Actions artifacts may expire and are not a durable formal download location.` |
| 158 | §3 Five-Step Quick Start | `1. Download, extract, and verify the four release files as described above.` |
| 195 | §5 Upgrade and Rollback | `2. Follow section 1 to extract the new ZIP into a new commit directory. Never overwrite the old directory.` |
| 206 | §6 Uninstall | `Before formal release, rerun these commands against the final ZIP and bind its version, commit, run, and filename in the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md).` |

零命中主题：无（七类主题均有命中）。
注（关键词命中但未列入）：L9/L50 黑箱候选、L21/L23/L28/L34/L44/L46/L48（Copilot CLI 安装/认证/模型）、L86/L90/L92/L97/L98/L111/L124/L134/L136（管线代码行）、L19 全文另含黑箱候选语句。

## 文件：docs/INSTALLATION_PLAN.md（共 63 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 3 | 标题下引语 | `> 2026-09-11 T026 已完成，候选最终状态：**10 verified、2 unsupported、0 unknown、blocked**。Job Object 取消 verified；工具别名和 shell 网络可绕过 CLI deny，故固定候选不兼容。不得自动升级或改变安装模型；T027/AT-23 仍须等待外部强制边界或兼容候选。详见 [T026 最终结论](validation/t026-agent-runner-contract.md)。` |
| 7 | 标题下导语 | `日期：2026-09-09。状态：S1 安装模型已决定，正式发行仍待批准。当前没有正式发行包或受支持的生产部署；本文件记录已冻结模型、待决事项和发布前门禁，实际步骤见 [便携 ZIP 安装指南](INSTALLATION.md)。` |
| 15 | §2 目标部署模型 | `……当前发行候选已交付逐行 CLI 与 T025 聚焦终端交互收件箱；AgentRunner 仍未交付；受监督黑箱候选的产品侧已由 T028 实现（AT-24 条件闭合），但真实 CLI 端到端依赖 T027-B、尚不可运行发布；完整统一 TUI 未实现。` |
| 19 | §3 待决决策（表头） | `\| 决策 \| 候选 \| 冻结条件/责任 \| 当前状态 \|` |
| 21 | §3 待决决策 | `\| S1 发行载体 \| Windows amd64 便携 ZIP，手工解压；后续可另行评估 winget \| T017 原生验收、固定 commit、SBOM、许可证清单和校验和完成 \| 2026-09-09 所有者批准 \|` |
| 22 | §3 待决决策 | `\| Linux 支持 \| S1 核心 CI/预览；正式支持候选在 S2 \| 原生兼容、路径/权限/升级测试通过 \| S1 范围已定 \|` |
| 26 | §3 待决决策 | `\| 发布信任 \| S1 使用固定 commit、GitHub CI、SHA256SUMS、SBOM 和许可证清单 \| 校验和只证明内容完整性；签名/attestation 后续增强 \| S1 基线已定 \|` |
| 27 | §3 待决决策 | `\| 更新方式 \| 手工下载核验，新版本解压到新目录并显式切换路径；未来包管理器另议 \| 支持矩阵、schema 兼容、备份恢复门禁 \| S1 模型已定；最终 ZIP 待复验 \|` |
| 28 | §3 待决决策 | `\| 卸载与数据保留 \| 手工删除版本目录；run/store 默认保留，数据删除单独授权 \| T024 生命周期记录和显式删除授权 \| S1 模型与库级规则已冻结 \|` |
| 30 | §3 待决决策 | `\| VS Code visible 黑箱候选 \| 可选降级保证；要求 VS Code 1.82+、可用 Copilot Chat、SessionBridge 0.1.1 \| …… \| T028 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布；AT-24 条件闭合，S1 仍因 T027/AT-23 阻断 \|` |
| 32 | §3 待决决策 | `S1 不指定默认安装目录，不发布一键安装脚本，不修改 PATH，不声明自动更新，也不要求管理员权限。state/store 默认位置与正式下载入口仍待冻结；EOL 通知方式已由 ADR-010 获产品所有者批准（2026-11-16）。` |
| 34 | §4 已完成的候选便携包演练（标题） | `## 4. 已完成的候选便携包演练` |
| 36 | §4 已完成的候选便携包演练 | `2026-09-09 已对两个受信 Windows amd64 CI artifact 完成离线 verifier 核验及显式临时目录内的首次运行、升级、回滚和保留数据卸载。旧/新完整 candidate SHA 分别为 `1e7af676e8e84028c7ecfef1ccde91213738c20c` 与 `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`；详细证据见 [S1 Windows 便携包安装演练](validation/s1-install-drill.md)。` |
| 38 | §4 已完成的候选便携包演练 | `该结果关闭“候选 artifact 是否能执行生命周期”的工程疑问，并支撑 2026-09-09 获批的手工便携 ZIP、不改 PATH 模型。正式下载入口和最终 ZIP 尚未发布，仍须用最终载体重跑。` |
| 42 | §5 规划的用户路径 | `1. 从受信发布页选择与 OS/架构匹配的发行物。` |
| 43 | §5 规划的用户路径 | `2. 独立核验固定 commit、GitHub CI 结果、SHA256SUMS、SBOM/许可证清单和支持状态；明确校验和不认证发布者身份。` |
| 50 | §5 规划的用户路径 | `可复制命令见 [便携 ZIP 安装指南](INSTALLATION.md)。T024 已提供库级完整备份、新 store 核验、默认保留和显式删除授权规则；正式发行文件名和下载地址仍须在发行批准时绑定。` |
| 52 | §6 发布前安装门禁（标题） | `## 6. 发布前安装门禁` |
| 55 | §6 发布前安装门禁 | `- 发布物可离线核验 SHA256SUMS、SBOM、许可证和精确版本；这些证据不认证发布者身份。后续若增加签名/attestation，未知或撤销信息过期不得显示 fully verified。` |
| 59 | §6 发布前安装门禁 | `- README 与操作文档中的每条安装命令均由对应发行物在受支持平台实测。` |
| 63 | §7 文档完成条件 | `T016 已冻结 CLI、配置和当前源码运行方式；……S1 安装模型已于 2026-09-09 获批；最终 ZIP 复验和 S1 exit 获批前，本文仍不构成生产发布声明。[操作手册](OPERATIONS.md) 继续区分当前可执行与未来设计；候选发布边界见 [支持矩阵](S1_SUPPORT_MATRIX.md) 与 [发行说明](S1_RELEASE_NOTES.md)。` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。
注（关键词命中但未列入）：L11 命中「入口」（`开发中入口`，CLI 入口语义）；L29 命中「候选」（CLI Agent/AgentRunner 决策行）。

## 文件：docs/INSTALLATION_PLAN_EN.md（共 63 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 3 | Quote under title | `> T026 completed on 2026-09-11; final candidate status is **ten verified, two unsupported, zero unknown, blocked**. …… Do not automatically upgrade or change the installation model. T027/AT-23 still requires an external enforcement boundary or compatible candidate. See the [final T026 conclusion](validation/t026-agent-runner-contract_EN.md).` |
| 7 | Intro under title | `Date: 2026-09-09. Status: the S1 installation model is decided, while formal release remains pending. No formal release package or supported production deployment exists yet. This document records the frozen model, open decisions, and pre-release gates; see the [portable ZIP installation guide](INSTALLATION_EN.md) for operational steps.` |
| 15 | §2 Target Deployment Model | `……The current candidate delivers the line-oriented CLI and T025 focused terminal interaction inbox; AgentRunner is still undelivered; the supervised black-box candidate's product side is implemented by T028 (AT-24 conditionally closed) but its real CLI end-to-end depends on T027-B and it is not yet a runnable release; the complete unified TUI is not implemented.` |
| 19 | §3 Open Decisions (table header) | `\| Decision \| Candidate \| Freeze condition/owner \| Status \|` |
| 21 | §3 Open Decisions | `\| S1 artifact \| Manually extracted Windows amd64 portable ZIP; winget may be evaluated later \| T017 native acceptance, fixed commit, SBOM, license manifest, and checksums \| Owner approved 2026-09-09 \|` |
| 22 | §3 Open Decisions | `\| Linux support \| S1 core CI/preview; production support candidate in S2 \| Native compatibility, path/permission/upgrade tests \| S1 scope decided \|` |
| 26 | §3 Open Decisions | `\| Release trust \| S1 uses fixed commit, GitHub CI, SHA256SUMS, SBOM and license manifest \| Checksums prove only content integrity; signing/attestation is deferred \| S1 baseline decided \|` |
| 27 | §3 Open Decisions | `\| Update mechanism \| Manually verify and extract each version into a new directory, then explicitly switch paths; reconsider package managers later \| Support matrix, schema compatibility, backup/recovery gates \| S1 model decided; final ZIP awaits rerun \|` |
| 28 | §3 Open Decisions | `\| Uninstall/data retention \| Manually remove the version directory; retain run/store by default and separately authorize data deletion \| T024 lifecycle record and explicit deletion authorization \| S1 model and library rules frozen \|` |
| 30 | §3 Open Decisions | `\| VS Code visible black-box candidate \| …… \| …… \| T028 product side implemented (reduced assurance), real CLI end-to-end depends on T027-B, not yet a runnable release; AT-24 conditionally closed, S1 still blocked by T027/AT-23 \|` |
| 32 | §3 Open Decisions | `S1 defines no default install directory, publishes no one-click installer, does not modify PATH, claims no automatic updates, and requires no administrator privileges. The default state/store location and formal download entry point remain open; the EOL notification channel was approved by the product owner via ADR-010 (2026-11-16).` |
| 34 | §4 Completed Candidate Portable-Package Drill (heading) | `## 4. Completed Candidate Portable-Package Drill` |
| 36 | §4 Completed Candidate Portable-Package Drill | `On 2026-09-09, two trusted Windows amd64 CI artifacts completed offline verifier checks plus first run, upgrade, rollback, and evidence-preserving uninstall under an explicit temporary directory. The old and new full candidate SHAs are `1e7af676e8e84028c7ecfef1ccde91213738c20c` and `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`. See the [S1 Windows portable-package installation drill](validation/s1-install-drill_EN.md) for evidence.` |
| 38 | §4 Completed Candidate Portable-Package Drill | `This closes the engineering question of whether candidate artifacts can execute the lifecycle and supports the manually extracted portable-ZIP, no-PATH model approved on 2026-09-09. The formal download entry point and final ZIP do not yet exist, so the final carrier must rerun the drill.` |
| 42 | §5 Planned User Path | `1. Select an OS/architecture-matching artifact from a trusted release location.` |
| 43 | §5 Planned User Path | `2. Independently verify the fixed commit, GitHub CI result, SHA256SUMS, SBOM/license manifest, and support status; state explicitly that checksums do not authenticate the publisher.` |
| 50 | §5 Planned User Path | `See the [portable ZIP installation guide](INSTALLATION_EN.md) for copyable commands. T024 supplies library-level complete backup, new-store verification, retention-by-default, and explicit deletion-authorization rules; formal release filenames and download locations remain bound to release approval.` |
| 52 | §6 Pre-Release Installation Gates (heading) | `## 6. Pre-Release Installation Gates` |
| 55 | §6 Pre-Release Installation Gates | `- Artifacts support offline SHA256SUMS, SBOM, license, and exact-version verification; this evidence does not authenticate the publisher. If signing/attestation is added later, unknown or stale revocation information never displays fully verified.` |
| 59 | §6 Pre-Release Installation Gates | `- Every installation command in README and operations docs has been exercised against the corresponding artifact on a supported platform.` |
| 63 | §7 Documentation Completion | `T016 has frozen CLI, configuration, and current source-run behavior; …… The S1 installation model was approved on 2026-09-09. This document still does not claim a production release until the final ZIP rerun and S1 exit are approved; [Operations](OPERATIONS_EN.md) continues to distinguish executable behavior from future design. See the candidate [support matrix](S1_SUPPORT_MATRIX_EN.md) and [release notes](S1_RELEASE_NOTES_EN.md) for release boundaries.` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。
注：CN L11 对应句在 EN 侧无主题词命中。

## 文件：docs/ADR_REGISTER.md（共 117 行）—— 按请求范围（ADR-009/ADR-010 两行 + SLA 句）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 19 | 顶层 ADR 总表（表头 L9 `\| ID \| 决策/方案 \| 替代与理由 \| 状态/完成条件 \|`）——ADR-009 行，状态/完成条件列 | `\| ADR-009 \| S1 Windows amd64 使用手工解压的便携 ZIP；用户选择独立目录，以完整 `prfrail.exe` 路径运行，不修改 PATH/注册表/系统目录 \| 不采用安装器、包管理器、自动更新或隐式版本切换；最小化安装副作用并支持目录级并存/回滚 \| 2026-09-09 产品所有者明确批准；最终 ZIP 仍须绑定发行记录并复验 \|` |
| 20 | 顶层 ADR 总表——ADR-010 行 | `\| ADR-010 \| 产品版本采用 `vMAJOR.MINOR.PATCH`（首版候选 `v0.1.0`）；原型期新 minor 发布时，在 GitHub Release notes 与仓库支持矩阵同时公告前一 minor 的 EOL，初始窗口为发布后 90 个自然日 \| patch 不启动新的 EOL 倒计时；每次 minor 发布前按稳定性、用户量和维护能力评审是否延长未来窗口。不得追溯缩短已公告窗口；例外须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处 \| `APPROVED`（裁决人 = 用户；日期 = 2026-11-16；理由 = v0.1.0 进入 Release 流程前冻结 EOL 通知方式；依据 = 用户 2026-11-16 P77 起跑消息同轮裁定原文）；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准（**已于 2026-11-16 落实；P77**） \|` |
| 62 | 顶层（支持策略正文段，位于 ADR-010 行后） | `支持策略已由所有者于 2026-09-06 批准：……ADR-010 提议采用三段式 `vMAJOR.MINOR.PATCH`，首版候选为 `v0.1.0`；patch 发布不重置 EOL，新 minor 发布时用 GitHub Release notes 与仓库支持矩阵同时公告前一 minor 的 EOL。……ADR-010 已于 2026-11-16 获产品所有者批准；当前仍不是已发布 SLA。` |

补充（请求范围之外、含发布关键词的命中行，逐字摘录）：

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 14 | 顶层 ADR 总表——ADR-004 行 | ``\| ADR-004 \| SessionBridge v0.1.1 silent + 文件队列，核心自管持久幂等 \| …… \|（命中 `v0.1`，为 SessionBridge 版本号，非 ProofRail 发行版本）`` |
| 57 | 顶层 ADR 总表——ADR-046 行 | `\| ADR-046 \| **S1 Exit 收口口径与批准时机（P76 `S1-EXIT-CLOSEOUT-ASSESSMENT`）**：Q2 = (a) 判据字面（……）；Q3 = A 加法式（`s1-exit.md:52` 保持 `BLOCKED`，……）；Q10 = (C) 分两步批准（第一步批 ADR-010 + v0.1.0 Release，第二步在 S2 后批 S1 Exit 整体）；收口顺序 = P76 报告 §14 序列 \| …… \|` |
| 64 | 顶层（发布信任正文段） | `发布信任策略由所有者于 2026-09-08 简化：S1 不建立离线 Ed25519 密钥体系，以固定 commit SHA、GitHub CI、SHA256SUMS、SBOM 和许可证清单作为基础发布证据。SHA256SUMS 只证明下载内容与清单一致，不证明发布者身份；签名或 keyless attestation 移至后续增强，不阻断 T017/T018 或 S1 exit。`signature-receipt` 能力保留供可选跨主机/后续签名场景使用。` |
| 100 | §4 T001 所有者决策单 | `\| 产品所有者 \| 当前仓库用户作为最终范围/发布批准者 \| S1 编码、Git、发布仍须各自明确授权 \| 2026-09-06，产品所有者确认 \|` |
| 103 | §4 T001 所有者决策单 | `\| 发布信任 \| S1 使用固定 commit SHA、GitHub CI、SHA256SUMS、SBOM 和许可证清单；不建立离线密钥体系 \| SHA256SUMS 不得表述为身份认证；签名/attestation 后续增强 \| 2026-09-08，产品所有者批准简化策略 \|` |
| 104 | §4 T001 所有者决策单 | `\| 安装模型 \| Windows amd64 便携 ZIP 手工解压到用户自选独立目录；显式路径运行，不修改 PATH/注册表/系统目录 \| 最终 ZIP 仍需绑定发行记录并复验；不授权发布 \| 2026-09-09，产品所有者批准 \|` |

零命中主题（就 ADR-009/010 两行 + SLA 句范围而言）：下载入口/下载位置；复跑/再执行（全文件亦无「复跑/再执行」字样；「复验」见于 L19、L104）。
注（全文件关键词命中但未列入）：「下载/入口」另命中 L41（ADR-030 行「非生产入口」，非发布语义）；「v0.1.0」另命中 L57/L62（已列入补充）。

## 文件：docs/ADR_REGISTER_EN.md（共 118 行）—— 按请求范围

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 19 | Top ADR table (header L9 `\| ID \| Decision/proposal \| Alternatives/rationale \| Status/completion \|`) — ADR-009 row | `\| ADR-009 \| S1 Windows amd64 uses a manually extracted portable ZIP in a user-selected independent directory, invoked by full `prfrail.exe` path without modifying PATH, the registry, or system directories \| Reject installers, package managers, automatic updates, and implicit version switching; minimize installation side effects and support directory-level coexistence/rollback \| Product owner explicitly approved 2026-09-09; bind and rerun the final ZIP before release \|` |
| 20 | Top ADR table — ADR-010 row | `\| ADR-010 \| Product versions use `vMAJOR.MINOR.PATCH` (first candidate `v0.1.0`); during the prototype phase, a new minor announces the previous minor's EOL in both GitHub Release notes and the repository support matrix, initially 90 calendar days after release \| …… \| `APPROVED` (ruling party = the user; date = 2026-11-16; reason = freezing the EOL notification approach before v0.1.0 enters the release flow; basis = the user's same-turn P77 kickoff ruling of 2026-11-16); 90 days fits the early prototype, but the notice channels and evolution rule require explicit product-owner approval before the first release (**fulfilled on 2026-11-16; P77**) \|` |
| 63 | Top (support-policy paragraph) | `The owner approved this support policy on 2026-09-06: …… ADR-010 proposes three-part `vMAJOR.MINOR.PATCH` versions, beginning with candidate `v0.1.0`; …… ADR-010 was approved by the product owner on 2026-11-16; this is not yet a published SLA.` |

补充（请求范围之外、含发布关键词的命中行）：

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 22 | Top ADR table — ADR-012 row（状态列） | `\| ADR-012 \| …… \| …… \| `APPROVED / PARTIALLY IMPLEMENTED`; T026 completed the four wire/hash/replay contracts, capability matrix, and digest-pinned GitHub Copilot CLI 1.0.83 probes. …… T028 has completed the product-side `supervised-black-box` implementation and AT-24 is conditionally closed; the real CLI end-to-end depends on T027-B and it is not yet a runnable release; S1 remains blocked by T027/AT-23 \|` |
| 58 | Top ADR table — ADR-046 row | `\| ADR-046 \| **S1 Exit closeout caliber and approval timing (P76 `S1-EXIT-CLOSEOUT-ASSESSMENT`)**: Q2 = (a) criterion-literal (……); Q3 = A additive (`s1-exit.md:52` keeps `BLOCKED`; ……); Q10 = (C) two-step approval (step one approves ADR-010 plus the v0.1.0 Release; step two approves S1 Exit as a whole after S2); closeout order = the P76 report section 14 sequence \| …… \|` |
| 101 | §4 T001 Owner Decision Sheet | `\| Product owner \| Current repository user is final scope/release approver \| S1 coding, Git and release each still require explicit authorization \| Owner confirmed 2026-09-06 \|` |
| 104 | §4 T001 Owner Decision Sheet | `\| Release trust \| S1 uses fixed commit SHA, GitHub CI, SHA256SUMS, SBOM and license manifest; no offline key system \| Never describe SHA256SUMS as identity authentication; signing/attestation is a later enhancement \| Owner approved simplified policy 2026-09-08 \|` |
| 105 | §4 T001 Owner Decision Sheet | `\| Installation model \| Manually extract the Windows amd64 portable ZIP into a user-selected independent directory; invoke by explicit path without changing PATH, registry, or system directories \| The final ZIP still requires release-record binding and rerun; no publication authorization \| Owner approved 2026-09-09 \|` |

零命中主题（就 ADR-009/010 两行 + SLA 句范围而言）：下载入口/下载位置（全文件「下载/download/入口」零命中）；复跑/再执行（全文件零命中；「rerun」见于 L19、L105）。
注：「v0.1.0」另命中 L14（SessionBridge v0.1.1，SessionBridge 版本号）、L58/L63（已列入补充）。

## 文件：docs/DEV_PLAN.md（共 263 行）—— 按请求范围（T017/T018 相关行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 46 | §P0.3 技术探针（REQ-006/007/009/015/017/021/023/024/027） | ``T004 状态：`COMPLETE`。技术可行性已验证，但整体 S0 仍为 `NOT_READY`；依赖供应链、bootstrap、CI 发布证据与 S1 生产故障注入仍按 T017/后续任务关闭，本结果不授权 T005。`` |
| 90 | §P1.6 可信发布与影子验收（REQ-001/002/006/011/014/015/019/024/025/028） | ``- [x] T017 [Plan:P1.6] 建立 `.github/workflows/ci.yml`、`testdata/selfhost/`、`docs/validation/s1-selfhost.md`；固定工具链与 Actions 完整 SHA、两代隔离、较早可信 verifier commit、外部 oracle、原生平台、固定 candidate commit、SHA256SUMS、SBOM/许可证检查；依赖：T024 和独立 bootstrap 批准；验收：AT-13/14；禁止：候选构建/覆盖 verifier 或 seed、把校验和冒充身份认证、未授权发布。签名/attestation 属后续增强，不阻断 S1。完成：2026-09-09。……``；同 行续：``……push CI run [`34268076912`](https://github.com/larsonzh/prfrail/actions/runs/34268076912) 的 Ubuntu/Windows job 均成功。……[run `34277671704`](https://github.com/larsonzh/prfrail/actions/runs/34277671704) 的 test、candidate-build、candidate-probe、selfhost Windows/Linux 8 个 job 全部成功，6 个 artifact 的名称、大小和 GitHub archive digest 已归档。……`` |
| 92 | §P1.6 可信发布与影子验收 | ``T017 进度（P55 `EXIT-SEED-BUILD-EVIDENCE`）：AT-13 证据链已补齐并归档于[证据包](validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md)，……本片未改 CI 结构、未改产品代码，AT-13/AT-14 的判定口径与 T017 勾选状态不变。`` |
| 94 | §P1.6 可信发布与影子验收 | `- [ ] T018 [Plan:P1.6] 建立 `examples/{whois-shadow,go-minimal}/` 与 `docs/validation/s1-exit.md`；只读固定输入比较结果/失败分类，汇总全部 AT 和剩余风险；依赖：T017、T025、T026、T027、T028 及 whois 输入授权；验收：AT-15/22/23/24、RFC §14 S1 全部 exit；禁止：切换 whois 生产流程或“部分通过”冒充 S1 完成。` |
| 96 | §P1.6 可信发布与影子验收 | ``T018 状态：`IN_PROGRESS / BLOCKED`。……AT-01–AT-24 与剩余门禁见 [S1 exit 报告](validation/s1-exit.md)；T026/T027/AT-23、T028/AT-24、最终 ZIP 绑定/复验、EOL 通知方式和产品所有者 S1 exit 批准仍缺失，因此不勾选 T018。（**口径注记 2026-09-29，切片 `LEDGER-BACKFILL`**：……任务状态字段不变。）`` |
| 98 | §P1.6 可信发布与影子验收 | `> **T018 待议项（2026-09-22，B4 ① 实测）**：T018 剩余门（AT-24 / 最终 ZIP 绑定与复验 / EOL 通知方式 / 产品所有者 S1 exit 批准）需重新评估"**是否也依赖候选 / AgentRunner 装配进展**"——……` |
| 100 | §P1.6 可信发布与影子验收 | `T018 收口清单（2026-11-05，P57）：四道剩余门已展为八门清单（含 S-G2 的 AT-24 口径滞后、S-G6 的门级／任务级不同源），逐门判据与解锁条件见 `docs/validation/p57-closeout-assessment.md` §4；该待评估项**不立片**，清单即其收口产物。` |
| 104 | §P1.7 产品闭环与生命周期（在 P1.6 发布验收前执行） | `编号保持稳定，实际依赖顺序为 T016 → T019–T024 → T017 → T025 → T026 → T027 → T028 → T018。以下在 S0 授权后执行；T019–T024 已完成。每个任务可再按一个失败反例拆小，不一次生成完整功能。` |
| 149 | §P1.8 AI 执行、黑箱候选与操作员闭环（REQ-006/007/008/009/016/021/023/024/025/026/027） | `……③ **T018 剩余门**（AT-24、最终 ZIP 绑定与复验、EOL 通知方式、产品所有者 S1 exit 批准）**未闭合**，其是否同源依赖候选/装配进展已登记为 T018 待议项；④ **S1 exit** 未做。……生产派发、终局发布与生产门禁端口对在 Windows 被 fail-closed 拒绝，……**已提交 `e8cbff6` 并推送 `origin/main`（未推 gitee）**；CI run **`35812311023`** 双腿全绿，两腿 `Gates` 步骤均实测执行。` |

补充（§3 需求映射中含 T017/T018 的行，逐字摘录）：

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 210 | §3 需求映射 | `\| REQ-001 \| P0.1,P0.2,P1.6 \| T001,T002,T003,T018 \| AT-01/15；通用 Schema/非 C 样例 \|` |
| 211 | §3 需求映射 | `\| REQ-002 \| P1.2,P1.6 \| T009,T018 \| AT-03/15；chain/报告 \|` |
| 215 | §3 需求映射 | `\| REQ-006 \| P0.3,P1.4,P1.6,P1.8 \| T004,T013,T017,T026,T027 \| AT-08/14/23；adapter/能力/CI \|` |
| 220 | §3 需求映射 | `\| REQ-011 \| P1.5,P1.6 \| T015,T018 \| AT-11/15；三 harness \|` |
| 223 | §3 需求映射 | `\| REQ-014 \| P0.1,P1.3,P1.6 \| T001,T012,T018 \| AT-07/15；阶段报告 \|` |
| 224 | §3 需求映射 | `\| REQ-015 \| P0.3,P1.5,P1.6 \| T004,T016,T017 \| AT-12/13/14；构建/TUI \|` |
| 228 | §3 需求映射 | `\| REQ-019 \| P1.6 \| T017 \| AT-13/14；两代自托管 \|` |
| 233 | §3 需求映射 | `\| REQ-024 \| P0.2,P0.3,P1.1,P1.2,P1.3,P1.4,P1.6,P1.8 \| T002,T004,T005,T006,T008,T009,T011,T014,T017,T025,T026,T027,T028 \| AT-02/03/04/09/13/22/23/24；恢复 \|` |
| 234 | §3 需求映射 | `\| REQ-025 \| P0.2,P1.1,P1.3,P1.4,P1.6,P1.8 \| T002,T005,T011,T013,T014,T017,T025,T026,T027,T028 \| AT-04/06/10/14/22/23/24；执行证据 \|` |
| 237 | §3 需求映射 | `\| REQ-028 \| P0.1,P0.2,P1.1,P1.6 \| T001,T002,T003,T005,T017 \| AT-01/04/14；兼容 \|` |

零命中主题：版本号 `v0.1.0`/`0.1.0`（全文件零命中）。
注（范围外关键词命中）：「下载」另命中 L117/L135/L139/L150/L186（均属 T027 相关行）；「复跑」另命中 L147（T027 前置更新行）；「artifact」命中 L90/L96/L98/L149（均已列入）。

## 文件：docs/DEV_PLAN_EN.md（共 264 行）—— 按请求范围

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 47 | §P0.3 Technical Spikes | ``T004 status: `COMPLETE`. Technical feasibility is established, but overall S0 remains `NOT_READY`; dependency supply chain, bootstrap, CI release evidence and production crash injection remain T017/later gates. This result does not authorize T005.`` |
| 91 | §P1.6 Trusted Release and Shadow Acceptance | ``- [x] T017 [Plan:P1.6] Create .github/workflows/ci.yml, testdata/selfhost/ and docs/validation/s1-selfhost.md; pin tooling and Actions by full SHA, …… verify native platforms, a fixed candidate commit, SHA256SUMS, SBOM and licenses. Depends: T024 and bootstrap approval. Accept: AT-13/14. Exclude: candidate-built/overwritten verifier or seed, checksum-as-identity claims, or unauthorized release. Signing/attestation is a later enhancement and does not block S1. Completed: 2026-09-09. ……``；同行续：``……Both Ubuntu and Windows jobs succeeded in push CI [run `34268076912`](https://github.com/larsonzh/prfrail/actions/runs/34268076912). …… all eight test, candidate-build, candidate-probe, and selfhost Windows/Linux jobs succeeded in [run `34277671704`](https://github.com/larsonzh/prfrail/actions/runs/34277671704), with names, sizes, and GitHub archive digests for all six artifacts archived. ……`` |
| 93 | §P1.6 Trusted Release and Shadow Acceptance | `T017 progress (P55 `EXIT-SEED-BUILD-EVIDENCE`): the AT-13 evidence chain is complete and archived in the [evidence package](validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md), …… This slice changed neither the CI structure nor product code, so the AT-13/AT-14 judgment calibre and the T017 checkbox state stay unchanged.` |
| 95 | §P1.6 Trusted Release and Shadow Acceptance | `- [ ] T018 [Plan:P1.6] Create examples/{whois-shadow,go-minimal}/ and docs/validation/s1-exit.md; compare fixed read-only input/result/failure classifications and summarize all AT/risks. Depends: T017, T025, T026, T027, T028, and whois input authorization. Accept: AT-15/22/23/24 and every RFC S1 exit gate. Exclude: production whois cutover or partial-pass completion.` |
| 97 | §P1.6 Trusted Release and Shadow Acceptance | ``T018 status: `IN_PROGRESS / BLOCKED`. …… AT-01–AT-24 and remaining gates are summarized in the [S1 exit report](validation/s1-exit_EN.md). T026/T027/AT-23, T028/AT-24, final-ZIP binding/rerun, the EOL notification channel, and product-owner S1 exit approval remain absent, so T018 stays unchecked. ……`` |
| 99 | §P1.6 Trusted Release and Shadow Acceptance | `> **T018 open question (2026-09-22, measured by the B4 ①)**: the remaining T018 gates (AT-24 / final ZIP binding and re-verification / EOL notification / the product-owner S1-exit approval) need re-evaluation for ……` |
| 101 | §P1.6 Trusted Release and Shadow Acceptance | `T018 closeout inventory (2026-11-05, P57): the four remaining gates are expanded into an eight-gate inventory (including the S-G2 AT-24 wording lag and the S-G6 gate-versus-task source split), with per-gate criteria and unblock conditions in `docs/validation/p57-closeout-assessment_EN.md` §4; the open item **mints no slice** - the inventory is its closeout product.` |
| 105 | §P1.7 Product Flow and Lifecycle (Before P1.6 Release Acceptance) | `Keep IDs stable; execution order is T016 → T019–T024 → T017 → T025 → T026 → T027 → T028 → T018. These run under S0 approval; T019–T024 are complete. Split tasks by individual negative tests rather than generating whole features at once.` |
| 150 | §P1.8 AI Execution, Black-Box Candidates, and Operator Loop | `……③ **the remaining T018 gates** (AT-24, final ZIP binding and re-verification, the EOL notification channel, the product owner's S1-exit approval) are **not closed**, …… production dispatch, terminal publication and the production gate port pair are refused fail-closed on Windows, ……` |

补充（§3 Requirement Mapping 中含 T017/T018 的行）：

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 211 | §3 Requirement Mapping | `\| REQ-001 \| P0.1,P0.2,P1.6 \| T001,T002,T003,T018 \| AT-01/15; generic schema/non-C example \|` |
| 212 | §3 Requirement Mapping | `\| REQ-002 \| P1.2,P1.6 \| T009,T018 \| AT-03/15; chain/report \|` |
| 216 | §3 Requirement Mapping | `\| REQ-006 \| P0.3,P1.4,P1.6,P1.8 \| T004,T013,T017,T026,T027 \| AT-08/14/23; adapters/capability/CI \|` |
| 221 | §3 Requirement Mapping | `\| REQ-011 \| P1.5,P1.6 \| T015,T018 \| AT-11/15; three harnesses \|` |
| 224 | §3 Requirement Mapping | `\| REQ-014 \| P0.1,P1.3,P1.6 \| T001,T012,T018 \| AT-07/15; stage report \|` |
| 225 | §3 Requirement Mapping | `\| REQ-015 \| P0.3,P1.5,P1.6 \| T004,T016,T017 \| AT-12/13/14; build/TUI \|` |
| 229 | §3 Requirement Mapping | `\| REQ-019 \| P1.6 \| T017 \| AT-13/14; two generations \|` |
| 234 | §3 Requirement Mapping | `\| REQ-024 \| P0.2,P0.3,P1.1,P1.2,P1.3,P1.4,P1.6,P1.8 \| T002,T004,T005,T006,T008,T009,T011,T014,T017,T025,T026,T027,T028 \| AT-02/03/04/09/13/22/23/24; recovery \|` |
| 235 | §3 Requirement Mapping | `\| REQ-025 \| P0.2,P1.1,P1.3,P1.4,P1.6,P1.8 \| T002,T005,T011,T013,T014,T017,T025,T026,T027,T028 \| AT-04/06/10/14/22/23/24; execution evidence \|` |
| 238 | §3 Requirement Mapping | `\| REQ-028 \| P0.1,P0.2,P1.1,P1.6 \| T001,T002,T003,T005,T017 \| AT-01/04/14; compatibility \|` |

零命中主题：版本号 `v0.1.0`/`0.1.0`（全文件零命中）。
注（范围外关键词命中）：「download」另命中 L140（T027 行）；「artifact」另命中 L61（DG-01 行）、L136/L142（T027 行）；CN L149 对应行的 push/CI 语句在 EN L150 同有（本行截取处已含）。

## 文件：docs/OPERATIONS.md（共 148 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | 文档标题 | `# ProofRail 用户、操作员与可信发布手册` |
| 5 | 标题下导语 | `日期：2026-09-11；S1 操作设计稿。当前没有 ProofRail 正式发行包；手工解压 Windows amd64 便携 ZIP、显式路径运行且不修改 PATH 的 S1 安装模型已获批，候选步骤见 [安装指南](INSTALLATION.md)。……待决发行绑定见 [安装部署规划](INSTALLATION_PLAN.md)，完整用户流程见 [业务流程](BUSINESS_WORKFLOWS.md)。` |
| 9 | §1 当前可执行 | `开发环境 Go 1.22+；发布版未来不要求用户安装 Go，但被编排项目的工具链仍必需。在仓库根可执行：` |
| 33 | §1 当前可执行 | `……Bubble Tea v1.3.4 只完成过隔离原型，尚未进入发行二进制。……ProofRail 发布二进制不直接访问 GitHub、Go/npm registry 或其他公网 URL，发布证据核验只读本地文件；依赖下载和 GitHub Actions 网络仅属于开发/CI。……` |
| 37 | §1 当前可执行 | `**visible 调试观察（便利、非交付）**：……(1) visible **不是 S1 交付承诺**——`docs/S1_SUPPORT_MATRIX.md` 明写核心离线 CLI 不要求它，正式支持范围不变；……` |
| 41 | §2 安装与首次使用 | `S1 使用手工解压的 Windows amd64 便携 ZIP：核对固定 commit/GitHub CI 并验证 SHA256SUMS、SBOM 和许可证清单，解压到用户选择的独立版本目录，再以 `prfrail.exe` 完整路径执行 version 和静态预览。不得修改 PATH、注册表或系统目录，不得覆盖运行中的 host。可复制命令和升级/回滚/卸载步骤见 [安装指南](INSTALLATION.md)。SHA256SUMS 只证明内容与清单一致，不认证发布者身份；平台不匹配不得执行。Linux 在 S1 仅核心 CI/预览，不宣传正式支持。` |
| 45 | §2 安装与首次使用 | `核心离线模式不要求 VS Code。……AgentRunner 路径仍未交付（T027/AT-23 阻断 S1）；黑箱候选的产品侧实现已完成（T028，AT-24 条件闭合），但真实 CLI 端到端依赖 T027-B、尚不可运行发布；扩展已安装或消息可达不等于 AI 自动编程验收。` |
| 111 | §6 whois 影子迁移方案 | `whois v3.4.0 已封板，本轮及 S0 不改其生产代码、任务文件或发布目录。用户批准后，以只读、脱敏、固定摘要导出夹具；不得在 A/B 运行期间回滚源文件。` |
| 123 | §7 自托管、升级和发布（标题） | `## 7. 自托管、升级和发布` |
| 127 | §7 自托管、升级和发布 | `两代验收：稳定 N 捕获隔离输入并构建 N+1；候选不覆盖 host；新进程+独立 run/store 重放固定验收链；独立 oracle 比较 Schema、receipt、恢复和 fail-close goldens，完成 clean-room 演练。失败只阻断候选，不修改 seed 或历史证据。` |
| 129 | §7 自托管、升级和发布 | `发布清单：RFC/Schema/版本一致；固定 commit 与 GitHub CI；build/vet/test、原生平台和全部强制 AT；依赖许可证/漏洞清单及 SBOM；二进制对应平台/CGO_ENABLED=0；SHA256SUMS；升级/回滚测试；双语说明与已知限制；外部发布门禁和用户当轮 commit/push/release 授权。校验和不得冒充发布者身份认证；签名/attestation 是后续增强。默认远端 origin，Gitee 镜像不自动推送。` |
| 131 | §7 自托管、升级和发布 | `回退：停止候选进程并保留证据，只在明确存储兼容时恢复旧二进制；未知格式保持只读，不能降级覆写。发布信任、支持周期和漏洞响应策略见 [ADR_REGISTER.md](ADR_REGISTER.md)。` |
| 138 | §8 完整产品操作流程（规划，RFC §19） | `\| PC-02 交付 \| 选已接受结果，指定与源/run/store 不重叠的新目录，核验包摘要/排除项/评审引用 \| 拒绝覆盖；中断包不算交付；导出不等于外部产品发布 \|` |
| 144 | §8 完整产品操作流程（规划，RFC §19） | `S1 的交付包由用户自行决定后续合并/发布，ProofRail 不自动向原树写文件或执行 Git。安装包支持状态、SBOM/许可证、签名身份和撤销新鲜度分别检查；离线未知撤销不应显示“最新可信”。备份不包含凭据，需另行通过安全方式配置恢复环境。` |
| 146 | §8 完整产品操作流程（规划，RFC §19） | `……T024 已实现库级完整备份、新 store 核验、生命周期处置与发布清单；T017 发行物与安装命令仍待实现，完整产品流程当前仍不能作为可执行命令教程。` |
| 148 | §8 完整产品操作流程（规划，RFC §19）（切片同步引语） | ``> 切片同步（P55 `EXIT-SEED-BUILD-EVIDENCE`）：本片不适用（仅证据补链，无运维、值守与发布流程变更）；理由：本片只在 `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/` 归档既有 `internal/release/**` 行为的两平台捕获，未改 `.github/workflows/**` 结构与 Go 产品代码，故本文命令、值守动作与发布检查逐字不变。`` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。
注（关键词命中但未列入）：L43/L49/L69/L106 命中「绑定」（配置/链路绑定语义）；L77/L82/L91/L93/L106/L119/L139 命中「候选」（评审/黑箱候选语义）。

## 文件：docs/OPERATIONS_EN.md（共 148 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | Title | `# ProofRail User, Operator and Trusted Release Guide` |
| 5 | Intro under title | `Date: 2026-09-11; S1 operational design. No formal ProofRail release package exists. The S1 model of manually extracting a Windows amd64 portable ZIP, invoking it by explicit path, and not modifying PATH is approved; see the candidate [installation guide](INSTALLATION_EN.md). …… See the [installation plan](INSTALLATION_PLAN_EN.md) for unresolved release binding and [business workflows](BUSINESS_WORKFLOWS_EN.md) for the complete user flow.` |
| 9 | §1 Available Today | `Development needs Go 1.22+. Future release users will not need Go, but their project's harness toolchain remains required. From repository root:` |
| 33 | §1 Available Today | `……Bubble Tea v1.3.4 was only an isolated prototype and is not part of the release binary. …… The released ProofRail binary does not directly access GitHub, Go/npm registries, or other public URLs; release-evidence verification reads local files only. Dependency downloads and GitHub Actions networking belong only to development/CI. ……` |
| 41 | §2 Installation and First Run | `S1 uses a manually extracted Windows amd64 portable ZIP: verify the fixed commit/GitHub CI plus SHA256SUMS, SBOM, and license manifest; extract into a user-selected independent version directory; then invoke `prfrail.exe` by full path for version and static preview. Never modify PATH, the registry, or system directories, and never overwrite a running host. See the [installation guide](INSTALLATION_EN.md) for copyable commands and upgrade/rollback/uninstall steps. SHA256SUMS proves agreement with the manifest, not publisher identity. S1 Linux is core CI/preview, not production support. Do not execute wrong-platform binaries.` |
| 45 | §2 Installation and First Run | `Core offline mode does not require VS Code. …… The AgentRunner path is still undelivered (T027/AT-23 block S1); the black-box candidate's product-side implementation is complete (T028, AT-24 conditionally closed) but its real CLI end-to-end depends on T027-B and it is not yet a runnable release; an installed extension or reachable transport is not AI-coding acceptance.` |
| 111 | §6 whois Shadow Migration | `whois v3.4.0 is frozen. This round/S0 does not change its production code, tasks or release directories. After authorization, export fixed-hash read-only redacted fixtures; never restore source while A/B processes run.` |
| 123 | §7 Self-Hosting, Upgrade and Release (heading) | `## 7. Self-Hosting, Upgrade and Release` |
| 127 | §7 Self-Hosting, Upgrade and Release | `Stable N captures isolated input/builds N+1. The candidate never overwrites host. Replay a fixed chain in a new process and separate run/store, with an external oracle comparing schema/receipt/recovery/fail-close goldens and a clean-room drill. Failure blocks the candidate without changing seed or historical evidence.` |
| 129 | §7 Self-Hosting, Upgrade and Release | `Release checklist: aligned RFC/schema/version; fixed commit and GitHub CI; build/vet/test, native platforms and required AT; license/vulnerability inventory and SBOM; correct binaries/CGO_ENABLED=0; SHA256SUMS; upgrade/rollback tests; bilingual notes/limitations; independent gate and same-turn commit/push/release authorization. Checksums must not be presented as publisher authentication; signing/attestation is a later enhancement. Default remote is origin, never automatic Gitee mirror pushes.` |
| 131 | §7 Self-Hosting, Upgrade and Release | `Rollback stops candidate processes and preserves evidence; restore an older binary only when storage compatibility is explicit. Unknown formats stay read-only. Release-trust/support/response policies are in [ADR_REGISTER_EN.md](ADR_REGISTER_EN.md).` |
| 138 | §8 Complete Product Flow (Planned, RFC Section 19) | `\| PC-02 Delivery \| Select accepted results, choose a new nonoverlapping destination, verify hashes/exclusions/review references \| Reject overwrite; incomplete output is not delivery; export is not external release \|` |
| 144 | §8 Complete Product Flow (Planned, RFC Section 19) | `Users decide whether to merge/release S1 delivery packages externally; ProofRail does not write source or execute Git automatically. Check support status, SBOM/licenses, signing identity and revocation freshness separately. Unknown offline revocation must not display current trust. Backups exclude credentials; restore environments require separate secure configuration.` |
| 146 | §8 Complete Product Flow (Planned, RFC Section 19) | `…… T024 adds library-level complete backup, new-store verification, lifecycle disposition, and release inventory. T017 artifacts and installation commands remain pending, so the full product flow is still not an executable command tutorial today.` |
| 148 | §8 (slice-sync quote) | ``> Slice sync (P55 `EXIT-SEED-BUILD-EVIDENCE`): not applicable (evidence supplement only, no operations, on-call or release-process change); reason: this slice only archives the two-platform capture of the existing `internal/release/**` behaviour under `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/`, without touching the `.github/workflows/**` structure or Go product code, so the commands, on-call actions and release checks in this manual stay unchanged.`` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。
注：CN L37 对应句在 EN 侧无主题词命中；L118 命中「artifact hashes」（whois golden 表，未列入）；「bound」命中 L35/L43/L49/L69 等（配置/链路绑定语义，未列入）。

## 文件：docs/S1_SUPPORT_MATRIX.md（共 33 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | 文档标题 | `# ProofRail S1 候选支持矩阵` |
| 5 | 标题下导语 | ``日期：2026-09-09。状态：`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。`` |
| 7 | 顶层状态表（表头） | `\| 范围 \| S1 候选状态 \| 约束 \|` |
| 9 | 顶层状态表 | ``\| Windows 11 amd64 \| 正式支持候选 \| 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期；基础终端（逐行 CLI / `interactions tui` / `--progress`）保留；完整一体化 TUI 移出 S1、由 S2 承接 \|`` |
| 10 | 顶层状态表 | `\| Linux amd64 \| 核心 CI \| 仅构建、测试、候选探针和自托管；不宣称 S1 生产支持 \|` |
| 12 | 顶层状态表 | `\| generic/C/Go harness \| 配置与 noop 编排 \| language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close；Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力 \|` |
| 14 | 顶层状态表 | `\| VS Code / Copilot Chat / SessionBridge visible 黑箱候选 \| 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布 \| ……\|` |
| 16 | §版本与安全维护候选政策（标题） | `## 版本与安全维护候选政策` |
| 18 | §版本与安全维护候选政策 | ``- 产品版本采用 `vMAJOR.MINOR.PATCH`，首版候选为 `v0.1.0`；patch 发布不启动新的 EOL 倒计时。`` |
| 19 | §版本与安全维护候选政策 | `- 当前稳定 minor 与前一个 minor 接收安全修复；前一个 minor 的候选维护窗口为 90 天。` |
| 21 | §版本与安全维护候选政策 | `- **ADR-010（已批准，2026-11-16）**：新 minor 发布时在 GitHub Release notes 与本支持矩阵同时公告前一 minor 的 EOL 日期，日期为发布后第 90 个自然日。` |
| 22 | §版本与安全维护候选政策 | `- 90 天是原型期初始窗口；每次 minor 发布前按稳定性、用户量和维护能力评审，仅可延长后续新公告的窗口。` |
| 23 | §版本与安全维护候选政策 | `- 任何偏离须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处，且不得追溯缩短已公告窗口。ADR-010 已于 2026-11-16 获产品所有者批准；本页不替代 Release notes 的正式公告。` |
| 25 | §发行信任边界（标题） | `## 发行信任边界` |
| 27 | §发行信任边界 | `候选发行要求完整 commit SHA、GitHub CI、`SHA256SUMS`、SPDX SBOM 和许可证清单。校验和证明内容与清单一致，不认证发布者身份；S1 不宣称签名或 attestation。安装模型见 [便携 ZIP 安装指南](INSTALLATION.md)；具体文件名、下载位置和支持起始日期须与最终发行说明一起批准。` |
| 33 | §发行信任边界（切片同步引语） | `> 切片同步（P71 `S1-SCOPE-ANNOTATE`）：TUI 移出与 visible 降级保证口径已同步；理由：……Windows 行状态仍为「正式支持候选」，visible 行仍为「降级保证」。` |

零命中主题：绑定/bind（全文件零命中）。

## 文件：docs/S1_SUPPORT_MATRIX_EN.md（共 33 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 1 | Title | `# ProofRail S1 Candidate Support Matrix` |
| 5 | Intro under title | ``Date: 2026-09-09. Status: `CANDIDATE`, effective only after the formal release carrier and S1 exit are approved. This is not a published SLA.`` |
| 7 | Top status table (header) | `\| Scope \| S1 candidate status \| Constraints \|` |
| 9 | Top status table | `\| Windows 11 amd64 \| Production-support candidate \| Pure-Go `CGO_ENABLED=0` portable ZIP; manual extraction, explicit-path invocation, no PATH modification; rerun the lifecycle against the final ZIP; the basic terminal (line-oriented CLI / `interactions tui` / `--progress`) stays; the full unified TUI moves out of S1 and is taken up by S2 \|` |
| 10 | Top status table | `\| Linux amd64 \| Core CI \| Build, test, candidate probe, and self-host only; no S1 production-support claim \|` |
| 12 | Top status table | `\| generic/C/Go harnesses \| Configuration and noop orchestration \| Language scopes, read-only preview, and noop-only loops are verified; executable steps still fail closed; Python/Node official harness + three-language product tasks + the release-support commitment frozen to S2 (REQ-020 S2 commitment narrowing); existing tool scripts / open language IDs / generic capability are not frozen \|` |
| 14 | Top status table | `\| VS Code / Copilot Chat / SessionBridge visible black-box candidate \| Product side implemented (reduced assurance), real CLI end-to-end depends on T027-B, not yet a runnable release \| ……\|` |
| 16 | §Candidate Version and Security-Maintenance Policy (heading) | `## Candidate Version and Security-Maintenance Policy` |
| 18 | §Candidate Version and Security-Maintenance Policy | ``- Product versions use `vMAJOR.MINOR.PATCH`, beginning with candidate `v0.1.0`; patch releases do not start a new EOL countdown.`` |
| 19 | §Candidate Version and Security-Maintenance Policy | `- The current stable minor and previous minor receive security fixes; the candidate maintenance window for the previous minor is 90 days.` |
| 21 | §Candidate Version and Security-Maintenance Policy | `- **ADR-010 (approved 2026-11-16)**: when a new minor is released, announce the previous minor's EOL date in both the GitHub Release notes and this support matrix, with EOL on the 90th calendar day after release.` |
| 22 | §Candidate Version and Security-Maintenance Policy | `- Ninety days is the initial prototype-phase window; before each minor release, review stability, adoption, and maintenance capacity, and only lengthen future newly announced windows.` |
| 23 | §Candidate Version and Security-Maintenance Policy | `- Deviations require a dated ADR with rationale and expiry, explicit product-owner approval, and synchronized publication in both channels, and must not retroactively shorten an announced window. ADR-010 was approved by the product owner on 2026-11-16; this page does not replace the formal notice in the Release notes.` |
| 25 | §Release Trust Boundary (heading) | `## Release Trust Boundary` |
| 27 | §Release Trust Boundary | `A candidate release requires a full commit SHA, GitHub CI, `SHA256SUMS`, an SPDX SBOM, and a license manifest. Checksums prove agreement with a manifest, not publisher identity; S1 does not claim signing or attestation. See the [portable ZIP installation guide](INSTALLATION_EN.md); filenames, download locations, and the support start date must be approved with the final release notes.` |
| 33 | §Release Trust Boundary (slice-sync quote) | `> Slice sync (P71 `S1-SCOPE-ANNOTATE`): the TUI carve-out and the visible reduced-assurance wording are now aligned; reason: …… the Windows row status stays "Production-support candidate" and the visible row stays "reduced assurance".` |

零命中主题：bind/binding（全文件零命中）。
注：CN L13（CLI Agent/AgentRunner 行）在 EN 侧无主题词命中（EN L13 零命中）。

## 文件：docs/validation/s1-exit.md（共 77 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 7 | 标题下导语（# S1 Exit 验收报告） | `日期：2026-09-09。结论：`BLOCKED`。……可信 Windows artifact 的候选安装生命周期演练也已通过；手工解压便携 ZIP、不修改 PATH 的安装模型已获产品所有者批准，并已有双语教程、支持矩阵与发行说明草案。但最终 ZIP 尚未绑定和复验，EOL 通知方式尚未冻结，产品所有者也未批准最终 S1 exit，因此不得把部分通过标记为 S1 完成。……` |
| 34 | §AT-01–AT-24 汇总 | `\| AT-14 \| PASS \| [T017 self-host](s1-selfhost.md)，Windows/Linux 8/8 jobs 与 6 artifacts \|` |
| 46 | §AT-01–AT-24 汇总 | `这些 PASS 表示对应 AT 的现有确定性验收，不自动满足阶段文档与发布门禁。` |
| 56 | §RFC §14 S1 Exit | `\| 正式安装包教程 \| PARTIAL \| [便携 ZIP 安装指南](../INSTALLATION.md)与模型已冻结；最终 ZIP 尚未绑定和复验 \|` |
| 57 | §RFC §14 S1 Exit | `\| 安装/升级/回滚/卸载实测 \| PARTIAL \| [Windows 候选便携包演练](s1-install-drill.md)通过；最终载体确定后仍须复验 \|` |
| 58 | §RFC §14 S1 Exit | `\| 支持矩阵与发行说明 \| PARTIAL \| [候选支持矩阵](../S1_SUPPORT_MATRIX.md)与[发行说明草案](../S1_RELEASE_NOTES.md)已建立，尚未绑定正式发行物 \|` |
| 60 | §RFC §14 S1 Exit | `\| 产品所有者最终 S1 批准 \| BLOCKED \| 本轮只授权 whois 夹具读取，不等于 S1 发布/阶段批准 \|` |
| 66 | §剩余风险与下一步（标题 L64 `## 剩余风险与下一步`） | `1. 手工解压便携 ZIP、用户自选独立目录、显式路径运行且不修改 PATH 的模型已冻结；正式 ZIP 文件名和下载入口仍待发行绑定。EOL 通知与例外已由 ADR-010 获产品所有者批准（2026-11-16）。` |
| 67 | §剩余风险与下一步 | `2. Windows 候选 artifact 已完成首次运行、升级、回滚和保留数据卸载；正式载体确定后必须复验。Linux 只声明核心 CI，不宣称 S1 正式支持。` |
| 69 | §剩余风险与下一步 | `4. T025/AT-22 已完成聚焦终端人工交互闭环；……T028 已完成、AT-24 在 (b)/(c) 裁定下条件闭合（真实 CLI 端到端仍依赖 T027-B），其后仍需完成完整统一 TUI；再用最终 ZIP 复跑教程并绑定支持矩阵/发行说明，最后由产品所有者明确批准 S1 exit。在此之前 T018 与 S1 状态保持未完成。` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。
注（关键词命中但未列入）：L43 命中「artifacts」（AT-23 行内证据产物 `f2-pos artifacts`）；L71 命中「```artifacts」（证据文件清单围栏）；L11 命中「release/」（whois 目录名）。

## 文件：docs/validation/s1-exit_EN.md（共 77 行）

| 行号 | 小节 | 原文摘录 |
|---|---|---|
| 7 | Intro under title (# S1 Exit Validation Report) | `Date: 2026-09-09. Verdict: `BLOCKED`. …… a candidate lifecycle drill for trusted Windows artifacts also passes. The product owner approved manual portable-ZIP extraction without PATH modification, with bilingual installation, support-matrix, and release-note drafts now present. The final ZIP is not yet bound and rerun, the EOL notification channel remains open, and the product owner has not approved final S1 exit, so partial success must not be reported as S1 completion. ……` |
| 34 | §AT-01–AT-24 Summary | `\| AT-14 \| PASS \| [T017 self-host](s1-selfhost_EN.md), eight Windows/Linux jobs and six artifacts \|` |
| 46 | §AT-01–AT-24 Summary | `These PASS verdicts cover the existing deterministic AT evidence; they do not automatically satisfy stage documentation and release gates.` |
| 56 | §RFC Section 14 S1 Exit | `\| Formal install-package tutorial \| PARTIAL \| The [portable ZIP installation guide](../INSTALLATION_EN.md) and model are frozen; the final ZIP is not yet bound and rerun \|` |
| 57 | §RFC Section 14 S1 Exit | `\| Install/upgrade/rollback/uninstall exercises \| PARTIAL \| The [Windows candidate portable-package drill](s1-install-drill_EN.md) passes; the final carrier still requires rerun \|` |
| 58 | §RFC Section 14 S1 Exit | `\| Support matrix and release notes \| PARTIAL \| The [candidate support matrix](../S1_SUPPORT_MATRIX_EN.md) and [release-note draft](../S1_RELEASE_NOTES_EN.md) exist but are not bound to a formal release \|` |
| 60 | §RFC Section 14 S1 Exit | `\| Final product-owner S1 approval \| BLOCKED \| This session's fixture-read authorization is not release or stage approval \|` |
| 66 | §Residual Risks and Next Steps (heading L64 `## Residual Risks and Next Steps`) | `1. The manually extracted portable-ZIP model, user-selected independent directory, explicit-path invocation, and no-PATH policy are frozen. The formal ZIP filename and download entry point still require release binding. The EOL notice and exception policy was approved by the product owner via ADR-010 (2026-11-16).` |
| 67 | §Residual Risks and Next Steps | `2. Windows candidate artifacts now pass first-run, upgrade, rollback, and evidence-preserving uninstall; the final carrier must rerun the drill. Linux remains core CI only and is not S1 production support.` |
| 69 | §Residual Risks and Next Steps | `4. T025/AT-22 completes the focused terminal operator loop. …… then rerun the tutorial against the final ZIP, bind the support matrix/release notes, and obtain explicit product-owner S1 exit approval. Until then, T018 and S1 remain incomplete.` |

零命中主题：版本号 `v0.1.0`/`0.1.0`。
注：L43/L71 同 CN 侧（证据产物/围栏，未列入）。

## 关键句定位（CN/EN 对照）

**1. 三处绑定句**
- CN：`docs/INSTALLATION.md:206`（§6 卸载）→ `正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定最终版本、commit、run 和文件名。`
- EN：`docs/INSTALLATION_EN.md:206`（§6 Uninstall）→ `Before formal release, rerun these commands against the final ZIP and bind its version, commit, run, and filename in the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md).`

**2.「当前尚无正式 GitHub Release」句**
- CN：`docs/INSTALLATION.md:5`（标题下导语句内）→ `当前尚无正式 GitHub Release；新目录结构必须通过新的可信候选 artifact 复验后，才能视为与最终发行物结构相同。`
- EN：`docs/INSTALLATION_EN.md:5` → `No formal GitHub Release exists yet; the new layout must pass a fresh trusted-candidate replay before it can be treated as matching the final carrier.`
- 同主题相关（授权声明）：CN `docs/S1_RELEASE_NOTES.md:5` → `本文件不授权 tag、GitHub Release、发布或安装。`；EN `docs/S1_RELEASE_NOTES_EN.md:5` → `This file does not authorize a tag, GitHub Release, publication, or installation.`

**3.「当前仍不是已发布 SLA」句**
- CN：`docs/ADR_REGISTER.md:62`（段尾句）→ `ADR-010 已于 2026-11-16 获产品所有者批准；当前仍不是已发布 SLA。`
- EN：`docs/ADR_REGISTER_EN.md:63`（段尾句）→ `ADR-010 was approved by the product owner on 2026-11-16; this is not yet a published SLA.`
- 另有两处近似表述（无「仍」）：CN `docs/S1_SUPPORT_MATRIX.md:5` → `……待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。`；EN `docs/S1_SUPPORT_MATRIX_EN.md:5` → `……This is not a published SLA.`

**4. ADR-009 状态列句（P78 briefing 记作 `:20`，磁盘实读位于 `:19`）**
- CN：`docs/ADR_REGISTER.md:19` 状态/完成条件列全文 → `2026-09-09 产品所有者明确批准；最终 ZIP 仍须绑定发行记录并复验`（该行内无闭括号句；`:20` 为 ADR-010 行，其状态列尾部闭括号句为 `（**已于 2026-11-16 落实；P77**）`）
- EN：`docs/ADR_REGISTER_EN.md:19` → `Product owner explicitly approved 2026-09-09; bind and rerun the final ZIP before release`（`:20` 为 ADR-010 行，尾部闭括号为 `(**fulfilled on 2026-11-16; P77**)`）
- 事实补充：同文件 §4 T001 决策单「安装模型」行含近似绑定句——CN `docs/ADR_REGISTER.md:104` → `最终 ZIP 仍需绑定发行记录并复验；不授权发布`；EN `docs/ADR_REGISTER_EN.md:105` → `The final ZIP still requires release-record binding and rerun; no publication authorization`

**5. s1-exit 剩余风险第 1 条**
- CN：`docs/validation/s1-exit.md:66`（§剩余风险与下一步，标题在 L64）→ `1. 手工解压便携 ZIP、用户自选独立目录、显式路径运行且不修改 PATH 的模型已冻结；正式 ZIP 文件名和下载入口仍待发行绑定。EOL 通知与例外已由 ADR-010 获产品所有者批准（2026-11-16）。`
- EN：`docs/validation/s1-exit_EN.md:66`（§Residual Risks and Next Steps）→ `1. The manually extracted portable-ZIP model, user-selected independent directory, explicit-path invocation, and no-PATH policy are frozen. The formal ZIP filename and download entry point still require release binding. The EOL notice and exception policy was approved by the product owner via ADR-010 (2026-11-16).`

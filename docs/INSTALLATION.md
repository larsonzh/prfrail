# ProofRail S1 Windows 便携 ZIP 安装指南

[English](INSTALLATION_EN.md)

日期：2026-09-10。S1 安装模型已冻结为：Windows 11 amd64 便携 ZIP、手工解压到用户选择的独立目录、始终用 `prfrail.exe` 的显式路径运行，不修改 `PATH`、注册表或系统目录。ZIP 解压后只有一个 `prfrail/` 顶层目录，四个发行文件均在其中。当前尚无正式 GitHub Release；新目录结构必须通过新的可信候选 artifact 复验后，才能视为与最终发行物结构相同。

## 1. 环境要求与能力层次

| 项目 | 核心离线 CLI | 连接 Copilot Chat 的 AI 会话 |
|---|---|---|
| 操作系统/架构 | Windows 11 amd64（x86-64）；S1 唯一正式支持候选 | 同核心要求；VS Code 与扩展也必须在该用户会话中运行 |
| CPU/内存 | 无独立 GPU 要求；尚无可据以承诺固定最低 CPU、内存的 S1 性能基准 | 另受 VS Code、Copilot Chat 和目标项目工具链需求约束 |
| 磁盘 | 可写安装目录，以及彼此分离且容量足够的源目录、run workspace、state/store 和备份目录 | 还需 SessionBridge 文件通道空间；证据与快照容量取决于项目规模 |
| ProofRail 运行依赖 | 便携 `prfrail.exe` 不要求 Go、Python、Node.js、数据库、Docker 或管理员权限 | 核心二进制要求不变 |
| 主机软件 | Windows PowerShell 5.1+ 仅用于本指南的解压/核验命令；浏览器或可选 GitHub CLI 用于下载 | VS Code 1.82+、已登录且可用的 GitHub Copilot Chat、SessionBridge 0.1.1 扩展 |
| 项目工具 | `validate`/`preview` 不要求项目编译器；实际 harness 必须预装其声明的编译器、解释器、构建和测试工具 | 同左；ProofRail 不静默安装这些工具 |
| 网络 | 下载时访问 GitHub；下载和核验完成后，当前 `version`/`validate`/`preview` 可离线运行 | Copilot Chat 模型调用需要宿主可用网络、账号/订阅和模型权限 |

SessionBridge 不是核心离线 CLI 的安装依赖，但它的扩展是 ProofRail 通过 VS Code 与 Copilot Chat 建立 AI 会话连接的必要前置组件。AI 模式准备顺序是：安装 VS Code 和 GitHub Copilot Chat 并登录账号；从 SessionBridge [GitHub Releases](https://github.com/larsonzh/sessbridge/releases) 下载并安装 `0.1.1` 扩展，或从 VS Code Marketplace 安装扩展；然后重载 VS Code 窗口。ProofRail 的目标产品集成会内置 SessionBridge v1 文件 IPC 客户端，用户不需要另行下载或启动 Python、PowerShell、sh 客户端；[`client/`](https://github.com/larsonzh/sessbridge/tree/main/client) 中的三种实现只作为独立诊断、协议参考和 SessionBridge 手工联调工具。

不要混淆客户端程序与运行时通道：ProofRail 内置客户端或可选诊断客户端都通过文件 IPC 通信；通道目录默认在 Windows 为 `%TEMP%\sessbridge`，SessionBridge 参考客户端可由 `SESSBRIDGE_CHANNEL_DIR` 或 `--channel-dir`/`-ChannelDir` 覆盖；目标 VS Code 实例由 PID 选择。调用方和扩展必须使用同一通道目录及目标实例。ProofRail 的正式契约固定 `mode=silent`、`legacy=false`；每个连续 AI 作业必须发送非空、稳定的 `conversationId`，从而自动启用 SessionBridge RFC §5.1 的 silent 多轮历史，不能使用无状态空值。

SessionBridge 独立联调时可先用 `visible` 或 `@sbr-review` 验证面板投递和人工交互，再用带 `conversationId` 的 `silent` 验证 ProofRail 所需的自动回执及多轮上下文。但这些只是可选诊断能力；`visible`/`auto`/`@sbr-review` 不得成为 ProofRail 产品默认值或 silent 失败后的自动回退。正式人工干预由 ProofRail 自有 CLI/TUI 完成，改变该边界必须先修订 ADR-004/011、产品需求、契约和验收。

当前 ProofRail 尚未实现该内置 IPC 客户端：`internal/adapters/sessbridge.go` 只定义由消费方实现的 Go `SilentClient` 接口，产品运行时也没有通道目录、目标 PID 或会话身份的配置入口。扩展安装成功仍不等于 AI 闭环已经可用；当前 S1 候选的 executable step/adapter 产品闭环仍未交付，不能把 noop-only `run` 扩大解释为 AI 已能自动改码。

当前发行候选没有图形窗口、完整 TUI 或 AI/操作员交互收件箱。`prfrail.exe` 使用普通终端中的逐行 CLI 文本，默认无 ANSI 颜色，并为自动化提供 `--json`；可直接在 Windows Terminal、PowerShell 或 VS Code 集成终端中运行。Bubble Tea v1.3.4 仅完成过隔离可行性原型，尚未链接进当前二进制。S1 优先用“一个终端、短命令、可复制路径”降低学习成本，而不把原型界面冒充已交付功能。

## 2. 下载与离线核验

从 [S1 自托管证据](validation/s1-selfhost.md) 选择已通过可信矩阵的 GitHub Actions run。run 页的 `head_sha` 是 verifier commit，不是 candidate commit；必须同时核对证据报告中的 run ID、verifier commit、candidate commit 和全部 job 结论。

浏览器方式会保留 ZIP：登录 GitHub，打开 `https://github.com/larsonzh/prfrail/actions/runs/<run-id>`，在页面底部 **Artifacts** 区点击 `prfrail-Windows-<full-candidate-commit>`。下载文件通常名为 `prfrail-Windows-<full-candidate-commit>.zip`。

命令行下载会先检测已安装且已认证的 GitHub CLI：可用时，`gh run download` 会把 artifact 自动解压到 `--dir`，不会保留外层 ZIP；下载后四个文件位于 `$DownloadDir\prfrail`，因此跳过下方的 `Expand-Archive`。没有可用 `gh` 时，脚本改用 GitHub Actions REST API 下载 ZIP，再执行下方的 `Expand-Archive`。API 回退要求把具有该仓库 Actions 读取权限的 token 放入进程环境变量 `GH_TOKEN` 或 `GITHUB_TOKEN`；不要把 token 写进脚本、命令行参数或日志。SSH 密钥不能认证 GitHub REST API。

```powershell
$Repository = 'larsonzh/prfrail'
$RunId = '<run-id>'
$CandidateCommit = '<full-candidate-commit>'
$Artifact = "prfrail-Windows-$CandidateCommit"
$DownloadDir = Join-Path '<user-selected-independent-directory>' $CandidateCommit

$ghCommand = Get-Command gh -ErrorAction SilentlyContinue
$ghReady = $false
if ($ghCommand) {
    gh auth status *> $null
    $ghReady = ($LASTEXITCODE -eq 0)
}

if ($ghReady) {
    gh run view $RunId --repo $Repository
    gh run download $RunId --repo $Repository --name $Artifact --dir $DownloadDir
    if ($LASTEXITCODE -ne 0) {
        throw "Artifact download failed: $Artifact"
    }
    $InstallDir = Join-Path $DownloadDir 'prfrail'
    $Archive = $null
} else {
    $Token = $env:GH_TOKEN
    if (-not $Token) { $Token = $env:GITHUB_TOKEN }
    if (-not $Token) {
        throw 'GitHub CLI is unavailable or unauthenticated; set GH_TOKEN or GITHUB_TOKEN for the REST fallback'
    }

    $Headers = @{
        Authorization = "Bearer $Token"
        Accept = 'application/vnd.github+json'
        'X-GitHub-Api-Version' = '2022-11-28'
        'User-Agent' = 'ProofRail-install-guide'
    }
    $ListUri = "https://api.github.com/repos/$Repository/actions/runs/$RunId/artifacts?name=$Artifact"
    $Response = Invoke-RestMethod -Method Get -Uri $ListUri -Headers $Headers
    $Matches = @($Response.artifacts | Where-Object { $_.name -eq $Artifact })
    if ($Matches.Count -ne 1) {
        throw "Expected exactly one artifact named $Artifact in run $RunId; found $($Matches.Count)"
    }
    if ($Matches[0].expired) {
        throw "Artifact has expired: $Artifact"
    }

    New-Item -ItemType Directory -Path $DownloadDir -Force | Out-Null
    $Archive = Join-Path $DownloadDir "$Artifact.zip"
    $DownloadUri = "https://api.github.com/repos/$Repository/actions/artifacts/$($Matches[0].id)/zip"
    Invoke-WebRequest -UseBasicParsing -Uri $DownloadUri -Headers $Headers -OutFile $Archive
}
```

以下是浏览器或 REST API 下载 ZIP 后的解压路径；`gh` 分支已经自动解压，应跳过本段。将三个值替换为实际下载文件、完整 candidate commit 和当前用户可写的独立目录；若刚执行 REST 分支，则保留已经设置的 `$Archive`，不再执行第一行：

```powershell
$Archive = (Resolve-Path '.\prfrail-Windows-<full-commit>.zip').Path
$CandidateCommit = '<full-commit>'
$InstallRoot = '<user-selected-independent-directory>'
$ExtractDir = Join-Path $InstallRoot $CandidateCommit
$InstallDir = Join-Path $ExtractDir 'prfrail'

New-Item -ItemType Directory -Path $ExtractDir -Force | Out-Null
Expand-Archive -LiteralPath $Archive -DestinationPath $ExtractDir

Push-Location $InstallDir
try {
    foreach ($line in Get-Content -LiteralPath '.\SHA256SUMS') {
        if ($line -notmatch '^([0-9a-f]{64})  ([^\\/]+)$') {
            throw "Invalid SHA256SUMS line: $line"
        }
        $actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $Matches[2]).Hash.ToLowerInvariant()
        if ($actual -ne $Matches[1]) {
            throw "SHA-256 mismatch: $($Matches[2])"
        }
    }
    Get-Content -Raw -LiteralPath '.\sbom.spdx.json' | ConvertFrom-Json | Out-Null
    Get-Content -Raw -LiteralPath '.\licenses.json' | ConvertFrom-Json | Out-Null
    $version = & '.\prfrail.exe' version --json | ConvertFrom-Json
    if ($version.data.version -ne "0.1.0+$CandidateCommit") {
        throw "Candidate version mismatch: $($version.data.version)"
    }
} finally {
    Pop-Location
}
```

预期解压根目录只含 `prfrail/`，该目录只含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json` 和 `licenses.json`。校验和只证明文件与清单一致，不认证发布者身份；run/verifier/candidate 必须从受信发布记录独立核对，包内 `version --json` 必须绑定同一 candidate commit。不要执行平台或 commit 不匹配的二进制。Actions artifact 可能过期，它不是长期正式下载入口。

## 3. 五步快速上手

以下路径先验证无副作用的核心 CLI，不要求 SessionBridge，也不会调用 AI：

1. 按上一节下载、解压并核验四个发行文件。
2. 在独立测试 workspace 执行 `init`，生成最小 `proofrail.chain.json`。
3. 执行 `validate`，确认配置有效及当前 CLI 是否可运行。
4. 执行只读 `preview`，核对权限边界、unknown 和零调用计数。
5. 仅对 noop-only 链执行 `run`，随后用 `report` 查看结果；不要从此推断 executable/AI 闭环已可用。

```powershell
$PrfRail = Join-Path $InstallDir 'prfrail.exe'
$Workspace = '<empty-test-workspace>'
$RunRoot = '<separate-run-root>'
$RunDir = Join-Path $RunRoot 'quickstart'

New-Item -ItemType Directory -Path $Workspace, $RunRoot -Force | Out-Null

& $PrfRail init --workspace $Workspace
& $PrfRail validate --chain (Join-Path $Workspace 'proofrail.chain.json')
& $PrfRail preview --chain (Join-Path $Workspace 'proofrail.chain.json')
& $PrfRail run --chain (Join-Path $Workspace 'proofrail.chain.json') --run-id quickstart --run-dir $RunDir
& $PrfRail report --run-dir $RunDir
```

## 4. 首次运行检查

不把安装目录加入 `PATH`。在独立测试 workspace 中使用完整路径：

```powershell
$PrfRail = Join-Path $InstallDir 'prfrail.exe'
& $PrfRail version --json
& $PrfRail validate --chain 'C:\path\to\workspace\proofrail.chain.json' --json
& $PrfRail preview --chain 'C:\path\to\workspace\proofrail.chain.json' --json
```

`preview` 是离线只读操作。当前 `run` 仅支持 noop-only 任务链；executable `code/build/verify` step 会 fail-close。运行目录、state/store 和源目录必须显式分离。

## 5. 升级与回滚

1. 停止所有 ProofRail 写者，完成并验证 state/store 备份。
2. 按第 1 节把新 ZIP 解压到新的 commit 目录；不得覆盖旧目录。
3. 核验新版并用新版完整路径执行只读 `version`、`validate` 和 `preview`。
4. 升级通过后，后续命令改用新版完整路径；旧目录保留为回滚点。
5. 回滚时停止写者，改回旧版完整路径，并以只读方式核验历史 run。不要用旧程序写入已由不兼容新版本迁移的数据。

没有自动更新、`current` 链接或隐式版本切换。每次切换都由操作员显式选择二进制路径。

## 6. 卸载

确认没有运行中的 ProofRail 进程后，只删除对应版本的二进制目录。由于安装从未修改 `PATH`、注册表或系统目录，无需恢复这些设置。证据、配置、run 和 store 默认保留；删除它们需要单独授权并先确认保留/审计要求。

正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定最终版本、commit、run 和文件名。
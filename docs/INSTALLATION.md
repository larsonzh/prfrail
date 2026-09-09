# ProofRail S1 Windows 便携 ZIP 安装指南

[English](INSTALLATION_EN.md)

日期：2026-09-09。S1 安装模型已冻结为：Windows 11 amd64 便携 ZIP、手工解压到用户选择的独立目录、始终用 `prfrail.exe` 的显式路径运行，不修改 `PATH`、注册表或系统目录。当前尚无正式 GitHub Release；以下命令用于与最终发行物结构相同的可信候选 artifact，不能据此声称产品已发布。

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

SessionBridge 不是核心离线 CLI 的安装依赖，但它是 ProofRail 通过 VS Code 与 Copilot Chat 建立 AI 会话连接的必要前置组件。AI 模式准备顺序是：安装 VS Code 和 GitHub Copilot Chat并登录账号，从 SessionBridge [GitHub Releases](https://github.com/larsonzh/sessbridge/releases) 安装 `0.1.1` 扩展，然后重载 VS Code 窗口。扩展安装成功仍不等于 AI 闭环已经可用：还须确认 Copilot Chat 模型可用、ProofRail 与 SessionBridge 使用同一通道目录和目标 VS Code 实例，并以 `silent`、非 legacy 模式通信。当前 S1 候选的 executable step/adapter 产品闭环仍未交付；现阶段安装 SessionBridge 只能满足宿主连接前提，不能把 noop-only `run` 扩大解释为 AI 已能自动改码。

当前发行候选没有图形窗口或完整 TUI。`prfrail.exe` 使用普通终端中的逐行 CLI 文本，默认无 ANSI 颜色，并为自动化提供 `--json`；可直接在 Windows Terminal、PowerShell 或 VS Code 集成终端中运行。Bubble Tea v1.3.4 仅完成过隔离可行性原型，尚未链接进当前二进制。S1 优先用“一个终端、短命令、可复制路径”降低学习成本，而不把原型界面冒充已交付功能。

## 2. 下载与离线核验

从 [S1 自托管证据](validation/s1-selfhost.md) 选择已通过可信矩阵的 GitHub Actions run。run 页的 `head_sha` 是 verifier commit，不是 candidate commit；必须同时核对证据报告中的 run ID、verifier commit、candidate commit 和全部 job 结论。

浏览器方式会保留 ZIP：登录 GitHub，打开 `https://github.com/larsonzh/prfrail/actions/runs/<run-id>`，在页面底部 **Artifacts** 区点击 `prfrail-Windows-<full-candidate-commit>`。下载文件通常名为 `prfrail-Windows-<full-candidate-commit>.zip`。

已安装并登录 GitHub CLI 时，也可直接下载。`gh run download` 会把 artifact 自动解压到 `--dir`，不会保留外层 ZIP，因此执行此命令后跳过下方的 `Expand-Archive`，直接把 `$InstallDir` 设为 `$DownloadDir`：

```powershell
gh auth status
$Repository = 'larsonzh/prfrail'
$RunId = '<run-id>'
$CandidateCommit = '<full-candidate-commit>'
$Artifact = "prfrail-Windows-$CandidateCommit"
$DownloadDir = Join-Path '<user-selected-independent-directory>' $CandidateCommit

gh run view $RunId --repo $Repository
gh run download $RunId --repo $Repository --name $Artifact --dir $DownloadDir
if ($LASTEXITCODE -ne 0) {
    throw "Artifact download failed: $Artifact"
}
$InstallDir = $DownloadDir
```

以下是浏览器下载 ZIP 后的解压路径。将三个值替换为实际下载文件、完整 candidate commit 和当前用户可写的独立目录：

```powershell
$Archive = (Resolve-Path '.\prfrail-Windows-<full-commit>.zip').Path
$CandidateCommit = '<full-commit>'
$InstallRoot = '<user-selected-independent-directory>'
$InstallDir = Join-Path $InstallRoot $CandidateCommit

New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Expand-Archive -LiteralPath $Archive -DestinationPath $InstallDir

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

预期目录只含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json` 和 `licenses.json`。校验和只证明文件与清单一致，不认证发布者身份；run/verifier/candidate 必须从受信发布记录独立核对，包内 `version --json` 必须绑定同一 candidate commit。不要执行平台或 commit 不匹配的二进制。Actions artifact 可能过期，它不是长期正式下载入口。

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
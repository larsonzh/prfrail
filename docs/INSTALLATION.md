# ProofRail S1 Windows 便携 ZIP 安装指南

[English](INSTALLATION_EN.md)

日期：2026-09-09。S1 安装模型已冻结为：Windows 11 amd64 便携 ZIP、手工解压到用户选择的独立目录、始终用 `prfrail.exe` 的显式路径运行，不修改 `PATH`、注册表或系统目录。当前尚无正式 GitHub Release；以下命令用于与最终发行物结构相同的可信候选 artifact，不能据此声称产品已发布。

## 1. 下载与离线核验

从获批的 GitHub Actions run 下载名为 `prfrail-Windows-<full-commit>` 的 artifact ZIP，并核对 run、完整 candidate commit 和受支持平台。将下面三个值替换为实际下载文件、完整 commit 和当前用户可写的独立目录：

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

预期目录只含 `prfrail.exe`、`SHA256SUMS`、`sbom.spdx.json` 和 `licenses.json`。校验和只证明文件与清单一致，不认证发布者身份；run/commit 必须从受信发布记录独立核对。不要执行平台或 commit 不匹配的二进制。

## 2. 首次运行

不把安装目录加入 `PATH`。在独立测试 workspace 中使用完整路径：

```powershell
$PrfRail = Join-Path $InstallDir 'prfrail.exe'
& $PrfRail version --json
& $PrfRail validate --chain 'C:\path\to\workspace\proofrail.chain.json' --json
& $PrfRail preview --chain 'C:\path\to\workspace\proofrail.chain.json' --json
```

`preview` 是离线只读操作。当前 `run` 仅支持 noop-only 任务链；executable `code/build/verify` step 会 fail-close。运行目录、state/store 和源目录必须显式分离。

## 3. 升级与回滚

1. 停止所有 ProofRail 写者，完成并验证 state/store 备份。
2. 按第 1 节把新 ZIP 解压到新的 commit 目录；不得覆盖旧目录。
3. 核验新版并用新版完整路径执行只读 `version`、`validate` 和 `preview`。
4. 升级通过后，后续命令改用新版完整路径；旧目录保留为回滚点。
5. 回滚时停止写者，改回旧版完整路径，并以只读方式核验历史 run。不要用旧程序写入已由不兼容新版本迁移的数据。

没有自动更新、`current` 链接或隐式版本切换。每次切换都由操作员显式选择二进制路径。

## 4. 卸载

确认没有运行中的 ProofRail 进程后，只删除对应版本的二进制目录。由于安装从未修改 `PATH`、注册表或系统目录，无需恢复这些设置。证据、配置、run 和 store 默认保留；删除它们需要单独授权并先确认保留/审计要求。

正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定最终版本、commit、run 和文件名。
# ProofRail S1 Windows Portable ZIP Installation Guide

[中文](INSTALLATION.md)

Date: 2026-09-10. The frozen S1 installation model is a Windows 11 amd64 portable ZIP, manually extracted to an independent user-selected directory and always invoked through an explicit `prfrail.exe` path. It does not modify `PATH`, the registry, or system directories. The extracted ZIP has one top-level `prfrail/` directory containing all four release files. No formal GitHub Release exists yet; the new layout must pass a fresh trusted-candidate replay before it can be treated as matching the final carrier.

## 1. Environment Requirements and Capability Levels

| Item | Core offline CLI / CLI Agent | SessionBridge visible black-box candidate |
|---|---|---|
| OS/architecture | Windows 11 amd64 (x86-64), the only S1 production-support candidate | Same core requirement; VS Code and extensions must run in the same user session |
| CPU/memory | No discrete GPU requirement; S1 has no performance baseline that supports a fixed minimum CPU or memory commitment yet | Also constrained by VS Code, Copilot Chat, and the target project's toolchain |
| Disk | A writable install directory plus sufficiently sized, non-overlapping source, run-workspace, state/store, and backup directories | Also requires SessionBridge channel space; evidence and snapshot capacity depends on project size |
| ProofRail runtime | The portable `prfrail.exe` requires no Go, Python, Node.js, database, Docker, administrator privileges, or system service | Core binary requirements remain unchanged |
| Host software | Windows PowerShell 5.1+ is needed only for this guide's extraction/verification commands; use a browser or optional GitHub CLI to download | VS Code 1.82+, an authenticated and usable GitHub Copilot Chat installation, and the SessionBridge 0.1.1 extension |
| Project tools | `validate`/`preview` need no project compiler; real harnesses require their declared compilers, interpreters, build tools, and test tools to be preinstalled | Same as the core CLI; ProofRail never installs them silently |
| Network | GitHub access for download; current `version`/`validate`/`preview` can run offline after download and verification | Copilot Chat model calls require host networking, account/subscription access, and model entitlement |

SessionBridge is not an installation dependency for the core offline CLI or AgentRunner. Formal AI execution separately requires a version/digest-pinned CLI Agent supported by the T026 capability matrix. VS Code, usable Copilot Chat and SessionBridge 0.1.1 are required only when the user selects `supervised-black-box`. That mode requires human supervision and prior acknowledgment that ProofRail cannot prove tools/commands, out-of-scope reads, network, actual cost, every background process, session continuity, or external effects such as commit/push/publication. ProofRail warrants isolation and post-return artifact scanning, independent gates, evidence and review only.

Do not conflate a client program with its runtime channel. AgentRunner invokes a pinned CLI Agent directly. The embedded SessionBridge client and optional diagnostics use file IPC; the Windows default channel is `%TEMP%\sessbridge` and PID selects the VS Code instance. Black-box candidates fix `mode=visible` and `legacy=false`; a visible receipt proves UI delivery only, and the user must explicitly return the isolated workspace. `silent` is auxiliary analysis only and `auto` is excluded from the formal flow.

Standalone SessionBridge tests may verify visible delivery, `@sbr-review`, and silent multi-turn behavior; they do not prove the ProofRail product loop. Black-box mode is explicit only and never an automatic fallback from AgentRunner or silent. Risk acknowledgment cannot waive isolation, secret scanning, independent gates, or review.

ProofRail has not implemented AgentRunner, a formal CLI Agent adapter, or the visible IPC/TUI risk-acknowledgment and return flow required by black-box mode. `internal/adapters/sessbridge.go` exposes only a consumer-provided `SilentClient`; installing the extension enables no AI product loop. The noop-only `run` command is not automated AI coding.

The current release candidate has no graphical window, complete TUI, or AI/operator interaction inbox. `prfrail.exe` provides line-oriented CLI text in a normal terminal, ANSI-free by default, plus `--json` for automation; it runs in Windows Terminal, PowerShell, or the VS Code integrated terminal. Bubble Tea v1.3.4 was exercised only as an isolated feasibility prototype and is not linked into the current binary. S1 therefore prioritizes one terminal, short commands, and copyable paths for quick learning instead of presenting a prototype interface as delivered functionality.

## 2. Download and Offline Verification

Select a GitHub Actions run that passed the trust matrix in the [S1 self-host evidence](validation/s1-selfhost_EN.md). The run page `head_sha` is the verifier commit, not the candidate commit. Independently check the report's run ID, verifier commit, candidate commit, and all job conclusions.

The browser path preserves the ZIP: sign in to GitHub, open `https://github.com/larsonzh/prfrail/actions/runs/<run-id>`, and click `prfrail-Windows-<full-candidate-commit>` in the **Artifacts** section at the bottom of the page. The downloaded file is normally named `prfrail-Windows-<full-candidate-commit>.zip`.

The command-line path first detects an installed and authenticated GitHub CLI. When available, `gh run download` automatically extracts the artifact into `--dir` and does not retain the outer ZIP; the four files are placed under `$DownloadDir\prfrail`, so skip the `Expand-Archive` block below. Without usable `gh`, the script downloads the ZIP through the GitHub Actions REST API and then continues to the `Expand-Archive` block. The API fallback requires a token with Actions read access to this repository in the process environment as `GH_TOKEN` or `GITHUB_TOKEN`. Never put that token in the script, command-line arguments, or logs. SSH keys cannot authenticate the GitHub REST API.

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

The following path applies after downloading the ZIP through the browser or REST API; skip it after the `gh` branch, which already extracted the artifact. Replace the three values with the downloaded file, full candidate commit, and an independent directory writable by the current user. If the REST branch just set `$Archive`, retain that value and skip the first line:

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

The extraction root must contain only `prfrail/`, which in turn contains only `prfrail.exe`, `SHA256SUMS`, `sbom.spdx.json`, and `licenses.json`. Checksums prove agreement with the manifest, not publisher identity. Independently verify the run/verifier/candidate from a trusted record, and require package `version --json` to bind the same candidate commit. Never execute a binary for a different platform or commit. Actions artifacts may expire and are not a durable formal download location.

## 3. Five-Step Quick Start

This path first verifies the side-effect-free core CLI. It requires no SessionBridge and makes no AI call:

1. Download, extract, and verify the four release files as described above.
2. Run `init` in an isolated test workspace to create a minimal `proofrail.chain.json`.
3. Run `validate` to confirm the configuration and current CLI runnability.
4. Run read-only `preview` and inspect permission boundaries, unknowns, and zero call counters.
5. Run only the noop-only chain, then inspect it with `report`; do not infer executable/AI-loop availability from this result.

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

## 4. First-Run Checks

Do not add the install directory to `PATH`. Use the full path in an isolated test workspace:

```powershell
$PrfRail = Join-Path $InstallDir 'prfrail.exe'
& $PrfRail version --json
& $PrfRail validate --chain 'C:\path\to\workspace\proofrail.chain.json' --json
& $PrfRail preview --chain 'C:\path\to\workspace\proofrail.chain.json' --json
```

`preview` is offline and read-only. Current `run` supports noop-only chains; executable `code/build/verify` steps fail closed. Keep run directories, state/store, and source explicitly separate.

## 5. Upgrade and Rollback

1. Stop every ProofRail writer, then create and verify a state/store backup.
2. Follow section 1 to extract the new ZIP into a new commit directory. Never overwrite the old directory.
3. Verify the new version and run read-only `version`, `validate`, and `preview` through its full path.
4. After acceptance, use the new full path for subsequent commands and retain the old directory as the rollback point.
5. To roll back, stop writers, invoke the old full path, and verify historical runs read-only. Do not let an old binary write data migrated incompatibly by a newer version.

There is no automatic update, `current` link, or implicit version switch. The operator explicitly selects the binary path for every switch.

## 6. Uninstall

After confirming no ProofRail process is running, remove only the selected version's binary directory. Because installation never changed `PATH`, the registry, or system directories, none of them needs restoration. Evidence, configuration, runs, and stores are retained by default; deleting them requires separate authorization and prior retention/audit review.

Before formal release, rerun these commands against the final ZIP and bind its version, commit, run, and filename in the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md).
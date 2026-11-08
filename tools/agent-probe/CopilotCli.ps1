<#
.SYNOPSIS
    Pinned GitHub Copilot CLI wrapper (R-08.1).

.DESCRIPTION
    Single entry point for starting the GitHub Copilot CLI from probe tooling.

    buildCopilotArgs    Pure function with no side effects. Returns the argv array.
                        The constant prefix --no-auto-update is present unless
                        -Unlock is passed explicitly.

    Invoke-CopilotCli   Fail-closed launcher. With the default (-Unlock absent) it
                        verifies the pinned CLI version before running: first probe
                        plus one sample every fifth probe. The check runs
                        "<exe> --no-auto-update --version" and requires the stdout
                        text to contain the pinned version. On mismatch the call is
                        refused, evidence is written and the function throws (a
                        terminating error, so a script context exits non-zero).
                        With -Unlock the pin check is skipped, the observed version
                        string is recorded, and execution stops when that version is
                        outside the accepted set (version not controllable).

    Version evidence is only the text printed by "<exe> --version". winget and
    ProductVersion are not used (R-08.1).

.PARAMETER Subcommand
    Positional CLI subcommand, for example "version". Placed right after the
    prefix.

.PARAMETER Prompt
    Value for -p. Omitted when empty.

.PARAMETER Effort
    Value for --effort=<value>. Omitted when empty.

.PARAMETER Context
    Value for -C (workspace directory). Omitted when empty.

.PARAMETER ExePath
    Copilot CLI executable. Defaults to the pinned npm layout under %APPDATA%.

.PARAMETER Unlock
    Drop --no-auto-update so the SEA loader picks the highest installed version.
    Expected observed version is 1.0.85. Default is $false (fail-closed).

.NOTES
    PowerShell 5.1 compatible. Comments and output are English only.

    Known boundary (B1 DR-2): the version string is not a stable identity of the
    binary. This wrapper still uses the --version text as the evidence because
    R-08.1 requires that literal check; the SHA-256 pin stays authoritative
    wherever it is already enforced.

    The BYOK credential injection is deliberately NOT part of this wrapper; keys
    are injected by the tmp companion and cleared in a finally block.
#>

$script:CopilotPinnedVersion = '1.0.83'
$script:CopilotUnlockedAcceptedVersions = @('1.0.83', '1.0.85')
$script:CopilotProbeSampleEvery = 5
$script:CopilotDefaultExeRelative = 'npm\node_modules\@github\copilot\node_modules\@github\copilot-win32-x64\copilot.exe'
$script:CopilotCliLastExitCode = $null
$script:CopilotCliLastVersionText = ''

$script:CopilotProbeRoot = ''
if (-not [string]::IsNullOrEmpty($PSScriptRoot)) {
    $script:CopilotProbeRoot = Join-Path (Split-Path -Parent (Split-Path -Parent $PSScriptRoot)) 'tmp\P61-AT23-RESTRICTED-CANDIDATE'
}
if ($script:CopilotProbeRoot -eq '') {
    $script:CopilotProbeRoot = Join-Path (Get-Location).Path 'tmp\P61-AT23-RESTRICTED-CANDIDATE'
}
$script:CopilotProbeCountPath = Join-Path $script:CopilotProbeRoot 'wrapper-probe-count.txt'
$script:CopilotVersionEvidencePath = Join-Path $script:CopilotProbeRoot 'wrapper-version-evidence.txt'

function Get-CopilotUtf8NoBom {
    return (New-Object System.Text.UTF8Encoding($false))
}

function Update-CopilotProbeCount {
    # Append-only counter: it is only ever increased, never reset by this file.
    $directory = Split-Path -Parent $script:CopilotProbeCountPath
    if (-not (Test-Path -LiteralPath $directory)) {
        New-Item -ItemType Directory -Force -Path $directory | Out-Null
    }
    $count = 0
    if (Test-Path -LiteralPath $script:CopilotProbeCountPath) {
        $raw = ([IO.File]::ReadAllText($script:CopilotProbeCountPath, [Text.Encoding]::UTF8)).Trim()
        if ($raw -match '^\d+$') { $count = [int]$raw }
    }
    $count = $count + 1
    [IO.File]::WriteAllText($script:CopilotProbeCountPath, [string]$count, (Get-CopilotUtf8NoBom))
    return $count
}

function Write-CopilotVersionEvidence {
    param(
        [string]$ExePath,
        [string]$Mode,
        [bool]$Passed,
        [string]$VersionText
    )
    $directory = Split-Path -Parent $script:CopilotVersionEvidencePath
    if (-not (Test-Path -LiteralPath $directory)) {
        New-Item -ItemType Directory -Force -Path $directory | Out-Null
    }
    $firstLine = ''
    if (-not [string]::IsNullOrWhiteSpace($VersionText)) {
        foreach ($line in ($VersionText -split "`r?`n")) {
            if (-not [string]::IsNullOrWhiteSpace($line)) { $firstLine = $line.Trim(); break }
        }
    }
    if ($firstLine.Length -gt 120) { $firstLine = $firstLine.Substring(0, 120) }
    $verdict = 'mismatch'
    if ($Passed) { $verdict = 'ok' }
    $timestamp = (Get-Date).ToUniversalTime().ToString('o')
    $record = '{0} mode={1} verdict={2} exe={3} version={4}' -f $timestamp, $Mode, $verdict, $ExePath, $firstLine
    [IO.File]::AppendAllText($script:CopilotVersionEvidencePath, $record + "`n", (Get-CopilotUtf8NoBom))
}

function Read-CopilotVersionText {
    param(
        [string]$ExePath,
        [switch]$Pinned
    )
    $arguments = @('--version')
    if ($Pinned.IsPresent) { $arguments = @('--no-auto-update', '--version') }
    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    $code = -1
    try {
        $text = (& $ExePath @arguments 2>&1 | Out-String)
        if ($null -ne $LASTEXITCODE) { $code = [int]$LASTEXITCODE }
    } finally {
        $ErrorActionPreference = $previousPreference
    }
    return [pscustomobject]@{ Text = $text; ExitCode = $code }
}

function Get-CopilotVersionToken {
    param([string]$Text)
    $match = [regex]::Match($Text, '\d+\.\d+\.\d+')
    if ($match.Success) { return $match.Value }
    return ''
}

function buildCopilotArgs {
    [CmdletBinding()]
    param(
        [string]$Subcommand = '',
        [string]$Prompt = '',
        [string]$Effort = '',
        [string]$Context = '',
        [switch]$Unlock
    )

    $argv = New-Object System.Collections.Generic.List[string]

    # User ruling (design-r1 A.1): the SEA loader selects the highest installed
    # version when --no-auto-update is absent. The prefix is therefore dropped
    # only for an explicit -Unlock and kept for every other call.
    if (-not $Unlock.IsPresent) { $argv.Add('--no-auto-update') }

    if ($Subcommand -ne '') { $argv.Add($Subcommand) }
    if ($Prompt -ne '') { $argv.Add('-p'); $argv.Add($Prompt) }
    if ($Effort -ne '') { $argv.Add('--effort=' + $Effort) }
    if ($Context -ne '') { $argv.Add('-C'); $argv.Add($Context) }

    # Plain return so that an empty result pipes nothing; callers normalize with @(...).
    return $argv.ToArray()
}

function Invoke-CopilotCli {
    [CmdletBinding()]
    param(
        [string]$Subcommand = '',
        [string]$Prompt = '',
        [string]$Effort = '',
        [string]$Context = '',
        [string]$ExePath = '',
        [switch]$Unlock
    )

    $resolvedExe = $ExePath
    if ($resolvedExe -eq '') {
        if ([string]::IsNullOrEmpty($env:APPDATA)) {
            throw 'APPDATA is not defined; pass -ExePath explicitly.'
        }
        $resolvedExe = Join-Path $env:APPDATA $script:CopilotDefaultExeRelative
    }
    if (-not (Test-Path -LiteralPath $resolvedExe)) {
        throw ('Copilot CLI executable not found: {0}' -f $resolvedExe)
    }

    $probeCount = Update-CopilotProbeCount
    $mode = 'pinned'
    if ($Unlock.IsPresent) { $mode = 'unlock' }

    if (-not $Unlock.IsPresent) {
        # Stage 1 sampling: first probe plus one sample every fifth probe.
        $shouldVerify = $false
        if ($probeCount -eq 1) { $shouldVerify = $true }
        if (($probeCount % $script:CopilotProbeSampleEvery) -eq 0) { $shouldVerify = $true }
        if ($shouldVerify) {
            $sample = Read-CopilotVersionText -ExePath $resolvedExe -Pinned
            $matched = $sample.Text.Contains($script:CopilotPinnedVersion)
            if ($sample.ExitCode -ne 0) { $matched = $false }
            Write-CopilotVersionEvidence -ExePath $resolvedExe -Mode $mode -Passed $matched -VersionText $sample.Text
            if (-not $matched) {
                throw ('Pinned Copilot CLI version check failed: expected {0}, exit {1}. Refusing to run.' -f $script:CopilotPinnedVersion, $sample.ExitCode)
            }
            $script:CopilotCliLastVersionText = $sample.Text.Trim()
        }
    } else {
        $sample = Read-CopilotVersionText -ExePath $resolvedExe
        $token = Get-CopilotVersionToken -Text $sample.Text
        $accepted = $false
        foreach ($candidate in $script:CopilotUnlockedAcceptedVersions) {
            if ($token -eq $candidate) { $accepted = $true }
        }
        Write-CopilotVersionEvidence -ExePath $resolvedExe -Mode $mode -Passed $accepted -VersionText $sample.Text
        if (-not $accepted) {
            throw ('Copilot CLI version is not controllable: reported "{0}" is outside {1}. Refusing to run.' -f $token, ($script:CopilotUnlockedAcceptedVersions -join ', '))
        }
        $script:CopilotCliLastVersionText = $sample.Text.Trim()
    }

    $argv = @(buildCopilotArgs -Subcommand $Subcommand -Prompt $Prompt -Effort $Effort -Context $Context -Unlock:$Unlock.IsPresent)

    # The caller owns stdout and stderr redirection; nothing secret is printed here.
    & $resolvedExe @argv
    $exitCode = $LASTEXITCODE
    if ($null -eq $exitCode) { $exitCode = -1 }
    $script:CopilotCliLastExitCode = [int]$exitCode
    return [int]$exitCode
}

#requires -Version 5.1
<#
.SYNOPSIS
    Self-test for tools/agent-probe/Invoke-CandidateProbe.ps1.

.DESCRIPTION
    Exercises the runner end to end against a hermetic stub candidate, so the pin check, the dry-run
    gate, the artifact set, the OS argv capture and the version-text policy can all be verified
    without any billable model call.

    The test never touches the real candidate binary and never calls a model.

.PARAMETER KeepArtifacts
    Keep the temporary test directory instead of removing it (for inspection).
#>
[CmdletBinding()]
param(
    [switch]$KeepArtifacts
)

$ErrorActionPreference = 'Stop'

$RepoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$RunnerPath = Join-Path $PSScriptRoot 'Invoke-CandidateProbe.ps1'
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$root = Join-Path $RepoRoot "tmp\b1-probe\selftest-$stamp"
New-Item -ItemType Directory -Force -Path $root | Out-Null

$failures = New-Object System.Collections.Generic.List[string]
$checks = 0

function Assert-That {
    param([bool]$Condition, [string]$Message)
    $script:checks = $script:checks + 1
    if ($Condition) {
        Write-Output "ok   - $Message"
    } else {
        Write-Output "FAIL - $Message"
        $script:failures.Add($Message)
    }
}

function Invoke-Runner {
    param([string[]]$Arguments)
    $previous = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $output = & powershell -NoProfile -ExecutionPolicy Bypass -File $RunnerPath @Arguments 2>&1 | Out-String
        $code = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $previous
    }
    return [pscustomobject]@{ Output = $output; ExitCode = $code }
}

# Hermetic stub candidate: records that it ran, writes a usage file, emits a small JSON transcript,
# then stays alive briefly so the runner's Win32_Process polling can observe its command line.
$stubScript = Join-Path $root 'stub-candidate.ps1'
$stubMarker = Join-Path $root 'stub-executed.txt'
$stubBody = @'
$marker = Join-Path $PSScriptRoot 'stub-executed.txt'
[IO.File]::WriteAllText($marker, (Get-Date).ToUniversalTime().ToString('o'), (New-Object System.Text.UTF8Encoding($false)))
$usagePath = ''
for ($i = 0; $i -lt $args.Count; $i++) {
    if ($args[$i] -eq '--usage-output-file' -and $i + 1 -lt $args.Count) { $usagePath = $args[$i + 1] }
}
if ($usagePath -ne '') {
    [IO.File]::WriteAllText($usagePath, '{"totalPremiumRequestCost":1,"totalUserRequests":1}', (New-Object System.Text.UTF8Encoding($false)))
}
Write-Output '{"type":"assistant.message","data":{"content":"PONG"}}'
Write-Output '{"type":"result","data":{"exitCode":0}}'
Start-Sleep -Seconds 2
exit 0
'@
[IO.File]::WriteAllText($stubScript, $stubBody, (New-Object System.Text.UTF8Encoding($true)))

$stubCmd = Join-Path $root 'stub-candidate.cmd'
$stubCmdBody = "@echo off`r`npowershell -NoProfile -ExecutionPolicy Bypass -File `"%~dp0stub-candidate.ps1`" %*`r`n"
[IO.File]::WriteAllText($stubCmd, $stubCmdBody, (New-Object System.Text.UTF8Encoding($false)))
$stubSha256 = (Get-FileHash -LiteralPath $stubCmd -Algorithm SHA256).Hash

# 1. pin mismatch must fail closed and must not run anything.
$mismatchRoot = Join-Path $root 'mismatch'
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd,
    '-ExpectedSha256', '0000000000000000000000000000000000000000000000000000000000000000',
    '-Run', '-ProbeRoot', $mismatchRoot
)
Assert-That ($result.ExitCode -ne 0) 'pin mismatch exits non-zero'
Assert-That ($result.Output -match 'pin mismatch') 'pin mismatch reports the pin'
Assert-That (-not (Test-Path -LiteralPath $mismatchRoot)) 'pin mismatch creates no artifact root'
Assert-That (-not (Test-Path -LiteralPath $stubMarker)) 'pin mismatch does not execute the candidate'

# 2. dry run must validate and print, and must not execute the candidate or create artifacts.
$dryRoot = Join-Path $root 'dryrun'
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd,
    '-ExpectedSha256', $stubSha256, '-ProbeRoot', $dryRoot
)
Assert-That ($result.ExitCode -eq 0) 'dry run exits zero'
Assert-That ($result.Output -match 'dry run: scenario=provider-smoke') 'dry run reports the scenario'
Assert-That ($result.Output -match 'the candidate was NOT executed and no model call was made') 'dry run states that the candidate was not executed'
Assert-That ($result.Output -match 'versionText=<not sampled') 'dry run does not sample the version text'
Assert-That (-not (Test-Path -LiteralPath $stubMarker)) 'dry run does not execute the candidate'
Assert-That (-not (Test-Path -LiteralPath $dryRoot)) 'dry run creates no artifact root'

# 3. a real run must produce the full artifact set and summarise it.
$runRoot = Join-Path $root 'run'
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd,
    '-ExpectedSha256', $stubSha256, '-Run', '-ProbeRoot', $runRoot
)
Assert-That ($result.ExitCode -eq 0) 'stub run exits zero'
Assert-That (Test-Path -LiteralPath $stubMarker) 'a real run does execute the candidate'
$artifacts = @(Get-ChildItem -LiteralPath $runRoot -File -ErrorAction SilentlyContinue)
foreach ($suffix in @('.argv.json', '.osargv.txt', '.stdout.json', '.stderr.txt', '.usage.json', '.result.json')) {
    Assert-That (($artifacts | Where-Object { $_.Name.EndsWith($suffix) }).Count -ge 1) "run writes a $suffix artifact"
}
$osArgv = Get-ChildItem -LiteralPath $runRoot -Filter '*.osargv.txt' | Select-Object -First 1
$osArgvText = if ($osArgv) { Get-Content -LiteralPath $osArgv.FullName -Raw } else { '' }
Assert-That ($osArgvText -match 'stub-candidate') 'OS argv capture sees the spawned stub process'
$argvPath = Get-ChildItem -LiteralPath $runRoot -Filter '*.argv.json' | Select-Object -First 1
$argvRecord = if ($argvPath) { Get-Content -LiteralPath $argvPath.FullName -Raw | ConvertFrom-Json } else { $null }
$cValue = if ($argvRecord) {
    $index = [array]::IndexOf(@($argvRecord.args), '-C')
    if ($index -ge 0 -and $index + 1 -lt @($argvRecord.args).Count) { ([string]@($argvRecord.args)[$index + 1]).Trim('"') } else { '' }
} else { '' }
Assert-That ($cValue -ne '') 'the recorded -C value is not empty'
Assert-That ((('{0}' -f $cValue).TrimEnd('\')) -eq (((Join-Path $runRoot 'workspace')).TrimEnd('\'))) 'the recorded -C value points at the probe workspace'
$summaryPath = Get-ChildItem -LiteralPath $runRoot -Filter '*.result.json' | Select-Object -First 1
$summary = if ($summaryPath) { Get-Content -LiteralPath $summaryPath.FullName -Raw | ConvertFrom-Json } else { $null }
Assert-That ($summary -and $summary.exitCode -eq 0) 'summary records the exit code'
Assert-That ($summary -and $summary.premiumRequests -eq 1) 'summary records the premium requests'
Assert-That ($summary -and $summary.expectedReplyHits -ge 1) 'summary counts the expected reply marker'

# 4. the version text is informational by default and fatal only on request.
$versionRoot = Join-Path $root 'version'
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd, '-ExpectedSha256', $stubSha256,
    '-ExpectedVersion', 'definitely-not-present', '-Run', '-ProbeRoot', $versionRoot
)
Assert-That ($result.ExitCode -eq 0) 'version text mismatch is informational by default'
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd, '-ExpectedSha256', $stubSha256,
    '-ExpectedVersion', 'definitely-not-present', '-RequireVersionText', '-Run', '-ProbeRoot', (Join-Path $root 'version-required')
)
Assert-That ($result.ExitCode -ne 0) '-RequireVersionText turns a mismatch into a failure'
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd, '-ExpectedSha256', $stubSha256,
    '-RequireVersionText', '-ProbeRoot', (Join-Path $root 'version-required-dry')
)
Assert-That ($result.ExitCode -ne 0) '-RequireVersionText without -Run is refused'
Assert-That ($result.Output -match 'requires -Run') 'the refusal names the missing -Run gate'

# 5. the artifact root stays inside <repo>\tmp: no escape, no existing root, no reparse redirection.
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd, '-ExpectedSha256', $stubSha256,
    '-Run', '-ProbeRoot', (Join-Path $RepoRoot 'outside-probe-root')
)
Assert-That ($result.ExitCode -ne 0) 'a ProbeRoot outside tmp is refused'
Assert-That ($result.Output -match 'must stay inside') 'the refusal names the tmp subtree rule'
Assert-That (-not (Test-Path -LiteralPath (Join-Path $RepoRoot 'outside-probe-root'))) 'the refused root is not created'

$existingRoot = Join-Path $root 'existing'
New-Item -ItemType Directory -Force -Path $existingRoot | Out-Null
$result = Invoke-Runner -Arguments @(
    '-Scenario', 'provider-smoke', '-Executable', $stubCmd, '-ExpectedSha256', $stubSha256,
    '-Run', '-ProbeRoot', $existingRoot
)
Assert-That ($result.ExitCode -ne 0) 'an existing ProbeRoot is refused'
Assert-That ($result.Output -match 'already exists') 'the refusal names the overwrite rule'

$junctionTarget = Join-Path $root 'junction-target'
New-Item -ItemType Directory -Force -Path $junctionTarget | Out-Null
$junction = Join-Path $root 'junction'
$junctionCreated = $true
try { New-Item -ItemType Junction -Path $junction -Target $junctionTarget | Out-Null } catch { $junctionCreated = $false }
if ($junctionCreated) {
    $result = Invoke-Runner -Arguments @(
        '-Scenario', 'provider-smoke', '-Executable', $stubCmd, '-ExpectedSha256', $stubSha256,
        '-Run', '-ProbeRoot', (Join-Path $junction 'child')
    )
    Assert-That ($result.ExitCode -ne 0) 'a ProbeRoot behind a reparse point is refused'
    Assert-That ($result.Output -match 'reparse point') 'the refusal names the reparse rule'
} else {
    Assert-That $false 'the reparse-point fixture could not be created (junction support missing)'
}

if (-not $KeepArtifacts) {
    Remove-Item -LiteralPath $root -Recurse -Force -ErrorAction SilentlyContinue
    # Leave no empty scratch parent behind: remove tmp\b1-probe only when this run emptied it.
    $parent = Split-Path -Parent $root
    if ((Test-Path -LiteralPath $parent) -and -not (Get-ChildItem -LiteralPath $parent -Force)) {
        Remove-Item -LiteralPath $parent -Force -ErrorAction SilentlyContinue
    }
}

Write-Output ("selftest: checks=$checks failures=$($failures.Count)")
if ($failures.Count -gt 0) {
    $failures | ForEach-Object { Write-Output "failure: $_" }
    exit 1
}
exit 0

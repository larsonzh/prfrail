# B2: apply -> verify -> revert -> verify(zero difference) -> apply rehearsal.
# Proves restorability BEFORE any paid probe runs. Outputs every step into tmp/b2/rehearsal/<stamp>/.
# Run elevated.
[CmdletBinding()]
param(
    [string]$ContainerName = 'prfrail-b2-battery',
    [string]$BaseDir = '',
    [switch]$RefreshBaseline
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
Assert-B2Admin -Action 'Invoke-B2LoopbackRehearsal.ps1'
$ScratchRoot = Join-Path (Resolve-Path (Join-Path $PSScriptRoot '..\..\..')).Path 'tmp\b2'
if ($BaseDir -eq '') { $BaseDir = Join-Path $ScratchRoot 'state' }
New-Item -ItemType Directory -Force -Path $BaseDir | Out-Null

$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$outDir = Join-Path (Join-Path $ScratchRoot 'rehearsal') $stamp
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

function Invoke-Step([string]$Name, [string]$Script, [string[]]$Arguments) {
    $path = Join-Path $PSScriptRoot $Script
    $logPath = Join-Path $outDir ($Name + '.log')
    $previous = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $output = & powershell -NoProfile -ExecutionPolicy Bypass -File $path @Arguments 2>&1 | Out-String
        $stepCode = $LASTEXITCODE
        if ($null -eq $stepCode) { $stepCode = 0 }
    } finally {
        $ErrorActionPreference = $previous
    }
    [IO.File]::WriteAllText($logPath, $output, (New-Object System.Text.UTF8Encoding($false)))
    Write-Host ("step {0}: exit={1}" -f $Name, $stepCode)
    if ($stepCode -ne 0) {
        # A bare "exit=1" is useless for the operator: surface why the step failed and where the full log is.
        Write-Host ("step {0} failed; last output lines:" -f $Name)
        $tail = @($output -split "`r?`n" | Where-Object { $_.Trim() -ne '' } | Select-Object -Last 12)
        foreach ($line in $tail) { Write-Host ('  ' + $line.Trim()) }
        Write-Host ('  full log: ' + $logPath)
    }
    return [int]$stepCode
}

$common = @('-ContainerName', $ContainerName, '-BaseDir', $BaseDir)

# A baseline that predates schemaVersion 2 cannot answer the profile-existence question, so the first apply
# would refuse it. The rehearsal is itself an explicit, elevated, deliberate action, so it detects that case
# and forwards -RefreshBaseline to the first apply (apply still refuses unless the machine is at restoral
# state, so this cannot bake our own change into a fresh baseline).
$refreshFirstApply = [bool]$RefreshBaseline
$baselinePath = Join-Path $BaseDir 'baseline.json'
if (-not $refreshFirstApply -and (Test-Path -LiteralPath $baselinePath)) {
    $existingBaseline = Read-B2JsonFile -Path $baselinePath
    if (-not (Test-B2SnapshotSchemaCurrent $existingBaseline)) {
        $refreshFirstApply = $true
        Write-Host 'baseline-refresh: the recorded baseline predates schemaVersion 2; forwarding -RefreshBaseline to step 1-apply'
    }
}
$firstApplyArgs = if ($refreshFirstApply) { $common + @('-RefreshBaseline') } else { $common }
$results = New-Object System.Collections.Generic.List[object]
$failed = $false
foreach ($step in @(
        @{ name = '1-apply'; script = 'apply-b2-loopback.ps1'; args = $firstApplyArgs },
        @{ name = '2-verify-present'; script = 'verify-b2-loopback.ps1'; args = ($common + @('-Expect', 'present')) },
        @{ name = '3-revert'; script = 'revert-b2-loopback.ps1'; args = $common },
        @{ name = '4-verify-absent-zero-diff'; script = 'verify-b2-loopback.ps1'; args = ($common + @('-Expect', 'absent')) },
        @{ name = '5-apply-final'; script = 'apply-b2-loopback.ps1'; args = $common },
        @{ name = '6-verify-final-present'; script = 'verify-b2-loopback.ps1'; args = ($common + @('-Expect', 'present')) }
    )) {
    $code = [int](Invoke-Step -Name $step.name -Script $step.script -Arguments $step.args)
    $results.Add([pscustomobject]@{ step = $step.name; exitCode = $code })
    if ($code -ne 0) { $failed = $true; Write-Host ("rehearsal ABORTED at " + $step.name); break }
}

$summary = [pscustomobject]@{
    stamp        = $stamp
    container    = $ContainerName
    ok           = (-not $failed)
    steps        = $results
    outputDir    = $outDir
    revertedDiff = (Read-B2TextFile -Path (Join-Path $BaseDir 'revert-diff.txt'))
}
[IO.File]::WriteAllText((Join-Path $outDir 'summary.json'), ($summary | ConvertTo-Json -Depth 6), (New-Object System.Text.UTF8Encoding($false)))
Write-Output ('rehearsal-ok=' + (-not $failed) + ' output=' + $outDir)
if ($failed) { exit 1 }

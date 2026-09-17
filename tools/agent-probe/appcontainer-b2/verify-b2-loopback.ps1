# B2: verify the loopback-exemption state and compare the current snapshot with the baseline.
# Read-only apart from writing the snapshot file. Run elevated (Playback of the exemption list
# needs it on some builds; snapshot writing works either way).
[CmdletBinding()]
param(
    [string]$ContainerName = 'prfrail-b2-battery',
    [string]$BaseDir = '',
    [ValidateSet('present', 'absent', 'any')]
    [string]$Expect = 'any'
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
$ScratchRoot = Join-Path (Resolve-Path (Join-Path $PSScriptRoot '..\..\..')).Path 'tmp\b2'
if ($BaseDir -eq '') { $BaseDir = Join-Path $ScratchRoot 'state' }
New-Item -ItemType Directory -Force -Path $BaseDir | Out-Null

# Read-only: name-only derivation works whether or not the profile still exists, and nothing here creates it.
$sid = Get-B2ContainerSidByName -ContainerName $ContainerName
$profileExists = Get-B2ContainerProfileExists -ContainerName $ContainerName
$exempt = Get-B2ExemptionList
$present = ($exempt.sids -contains $sid)
Write-Output ('containerSid: ' + $sid)
Write-Output ('profileExists: ' + $profileExists)
Write-Output ('exemptionPresent: ' + $present)
Write-Output ('exemptionCount: ' + $exempt.sids.Count)

$baselinePath = Join-Path $BaseDir 'baseline.json'
if (Test-Path $baselinePath) {
    $currentPath = Join-Path $BaseDir 'current.json'
    $current = New-B2Snapshot -OutFile $currentPath -ContainerName $ContainerName
    $baseline = Read-B2JsonFile -Path $baselinePath
    $strict = @(Get-B2SnapshotDiff -Baseline $baseline -Current $current)
    $failures = @(Get-B2RestoralFailures -DiffLines $strict)
    $listeners = @(Get-B2ListenerDiff -Baseline $baseline -Current $current)
    $strictText = if ($strict.Count -eq 0) { 'zero-difference' } else { ($strict -join ' | ') }
    $listenerText = if ($listeners.Count -eq 0) { 'listener-zero-churn' } else { ($listeners -join ' | ') }
    Write-Output ('restoralDiff: ' + $strictText)
    Write-Output ('restoralFailures: ' + $failures.Count)
    Write-Output ('listenerDiff: ' + $listenerText)
} else {
    $failures = @()
    $strictText = 'no-baseline'
    $listenerText = 'no-baseline'
    Write-Output 'restoralDiff: no-baseline'
}

$ok = $true
switch ($Expect) {
    'present' { $ok = $present }
    'absent' { $ok = ($present -eq $false) -and ($failures.Count -eq 0 -or $strictText -eq 'no-baseline') }
    'any' { $ok = $true }
}
Write-Output ('verify: expect=' + $Expect + ' ok=' + $ok)
if (-not $ok) { exit 1 }

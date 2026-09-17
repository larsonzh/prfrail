# B2: remove the AppContainer loopback exemption and prove zero differences against the baseline.
# Run elevated. Pair with apply-b2-loopback.ps1.
[CmdletBinding()]
param(
    [string]$ContainerName = 'prfrail-b2-battery',
    [string]$BaseDir = '',
    [switch]$KeepProfile
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
$ScratchRoot = Join-Path (Resolve-Path (Join-Path $PSScriptRoot '..\..\..')).Path 'tmp\b2'
if ($BaseDir -eq '') { $BaseDir = Join-Path $ScratchRoot 'state' }

$baselinePath = Join-Path $BaseDir 'baseline.json'
$resultPath = Join-Path $BaseDir 'revert-result.json'
$diffPath = Join-Path $BaseDir 'revert-diff.txt'
$currentPath = Join-Path $BaseDir 'after-revert.json'
$steps = New-Object System.Collections.Generic.List[object]

function Add-Step($name, $detail) {
    $script:steps.Add([pscustomobject]@{ name = $name; detail = $detail; at = (Get-Date).ToUniversalTime().ToString('o') })
}

Assert-B2Admin -Action 'revert-b2-loopback.ps1'
# Name-only derivation: revert must still work after the profile has been deleted (the exemption is keyed by SID).
$sid = Get-B2ContainerSidByName -ContainerName $ContainerName
Add-Step 'container-sid' $sid

$before = Get-B2ExemptionList
Add-Step 'exemption-list-before' (($before.sids -join ','))

if ($before.sids -notcontains $sid) {
    Add-Step 'revert' 'exemption absent already; nothing to remove (idempotent)'
} else {
    $output = (& CheckNetIsolation.exe LoopbackExempt -d -p=$sid 2>&1 | Out-String).Trim()
    Add-Step 'revert-checknetisolation' $output
    $after = Get-B2ExemptionList
    if ($after.sids -contains $sid) {
        $target = @($after.sids | Where-Object { $_ -ne $sid })
        Add-Step 'revert-fallback' ('NetworkIsolationSetAppContainerConfig with ' + $target.Count + ' entries')
        $hr = [B2LoopbackNative]::SetExemptionList([string[]]$target)
        Add-Step 'revert-fallback-hresult' ('0x{0:X8}' -f $hr)
        $after = Get-B2ExemptionList
    }
    if ($after.sids -contains $sid) {
        Add-Step 'revert' 'FAILED (exemption still present after both mechanisms)'
        $failed = [pscustomobject]@{ ok = $false; sid = $sid; steps = $steps; exemptions = $after.sids }
        [IO.File]::WriteAllText($resultPath, ($failed | ConvertTo-Json -Depth 6), (New-Object System.Text.UTF8Encoding($false)))
        throw "loopback exemption could not be removed for $sid"
    }
    Add-Step 'revert' 'removed'
}

# Restoral means "back to the recorded baseline", and the baseline records profile absence. Deleting the
# profile here (unless -KeepProfile) is what makes the strict diff reach zero difference; -KeepProfile keeps
# it for a B4 reuse and then the profile-existence line is reported as intentional instead of a failure.
$profileDeleted = $false
$profileKept = $false
if (Get-B2ContainerProfileExists -ContainerName $ContainerName) {
    if ($KeepProfile) {
        $profileKept = $true
        Add-Step 'profile' 'kept by request (-KeepProfile); restoral covers the exemption only'
    } else {
        $deleteResult = [B2LoopbackNative]::DeleteAppContainerProfile($ContainerName)
        Add-Step 'profile-delete-hresult' ('0x{0:X8}' -f $deleteResult)
        if ($deleteResult -ne 0) { throw ('DeleteAppContainerProfile failed 0x{0:X8}' -f $deleteResult) }
        if (Get-B2ContainerProfileExists -ContainerName $ContainerName) { throw 'the AppContainer profile still exists after deletion' }
        $profileDeleted = $true
        Add-Step 'profile' 'deleted (restores the baseline profile-absence state)'
    }
} else {
    Add-Step 'profile' 'absent already'
}

$current = New-B2Snapshot -OutFile $currentPath -ContainerName $ContainerName
if (-not (Test-Path $baselinePath)) {
    Add-Step 'diff' 'no baseline recorded; cannot prove zero difference (record a baseline before applying)'
    [IO.File]::WriteAllText($diffPath, "no-baseline`n", (New-Object System.Text.UTF8Encoding($false)))
    throw 'no baseline.json found; revert cannot prove restoral'
}
$baseline = Read-B2JsonFile -Path $baselinePath
$diff = @(Get-B2SnapshotDiff -Baseline $baseline -Current $current)
$failures = @(Get-B2RestoralFailures -DiffLines $diff -ProfileKept $profileKept)
$listenerDiff = @(Get-B2ListenerDiff -Baseline $baseline -Current $current)
$diffText = if ($diff.Count -eq 0) { 'zero-difference' } else { ($diff -join "`n") }
$listenerText = if ($listenerDiff.Count -eq 0) { 'listener-zero-churn' } else { ($listenerDiff -join "`n") }
[IO.File]::WriteAllText($diffPath, ('restoral: ' + $diffText + "`n" + 'listeners (informational): ' + $listenerText + "`n"), (New-Object System.Text.UTF8Encoding($false)))
Add-Step 'diff' $diffText
Add-Step 'listener-diff' $listenerText

$result = [pscustomobject]@{
    ok            = ($failures.Count -eq 0)
    sid           = $sid
    profileDeleted = $profileDeleted
    profileKept   = $profileKept
    diffLines     = $diff
    failureLines  = $failures
    listenerLines = $listenerDiff
    baselinePath  = $baselinePath
    currentPath   = $currentPath
    steps         = $steps
}
[IO.File]::WriteAllText($resultPath, ($result | ConvertTo-Json -Depth 6), (New-Object System.Text.UTF8Encoding($false)))
Write-Output ('revert-diff: ' + $diffText)
Write-Output ('revert-result: ' + $resultPath)
if ($failures.Count -ne 0) { exit 1 }

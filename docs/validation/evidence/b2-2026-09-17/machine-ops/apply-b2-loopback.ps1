# B2: grant the AppContainer loopback exemption needed by the enforcement boundary.
# Single system change; idempotent; fails closed with rollback; records a baseline first.
# Run elevated. Pair with revert-b2-loopback.ps1.
[CmdletBinding()]
param(
    [string]$ContainerName = 'prfrail-b2-battery',
    [string]$BaseDir = '',
    [switch]$RefreshBaseline
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
# Scratch lives in the git-ignored repo tree, never next to the (versioned) script.
$ScratchRoot = Join-Path (Resolve-Path (Join-Path $PSScriptRoot '..\..\..')).Path 'tmp\b2'
if ($BaseDir -eq '') { $BaseDir = Join-Path $ScratchRoot 'state' }
New-Item -ItemType Directory -Force -Path $BaseDir | Out-Null

$baselinePath = Join-Path $BaseDir 'baseline.json'
$resultPath = Join-Path $BaseDir 'apply-result.json'
$steps = New-Object System.Collections.Generic.List[object]

function Add-Step($name, $detail) {
    $script:steps.Add([pscustomobject]@{ name = $name; detail = $detail; at = (Get-Date).ToUniversalTime().ToString('o') })
}

Assert-B2Admin -Action 'apply-b2-loopback.ps1'

# Baseline first: taking the snapshot before the (possibly profile-creating) SID resolution keeps profile
# creation visible to Get-B2SnapshotDiff instead of being absorbed into the baseline it just wrote.
if (Test-Path $baselinePath) {
    $existingBaseline = Read-B2JsonFile -Path $baselinePath
    # Fail closed on an untrustworthy baseline identity: a missing/empty containerName cannot be compared,
    # so it is refused instead of being treated as "matches anything".
    if ([string]::IsNullOrWhiteSpace($existingBaseline.containerName)) {
        throw 'the recorded baseline has no containerName and cannot be trusted; after a full revert run apply-b2-loopback.ps1 -RefreshBaseline to record a valid one'
    }
    if ($existingBaseline.containerName -ne $ContainerName) {
        throw ('the recorded baseline belongs to container "' + $existingBaseline.containerName + '", not "' + $ContainerName + '"; use a separate -BaseDir or refresh the baseline deliberately')
    }
    if (-not (Test-B2SnapshotSchemaCurrent $existingBaseline)) {
        if (-not $RefreshBaseline) {
            throw ('the recorded baseline predates schemaVersion 2 and cannot answer the profile-existence question; after a full revert run apply-b2-loopback.ps1 -RefreshBaseline to record a valid baseline')
        }
        # A refresh is only sound at restoral state: the container SID must not hold the exemption and the
        # profile must not exist, otherwise the refreshed "baseline" would bake in our own change.
        $preExempt = Get-B2ExemptionList
        $preSid = Get-B2ContainerSidByName -ContainerName $ContainerName
        if (@($preExempt.sids) -contains $preSid) {
            throw 'refusing to refresh the baseline while the container SID still holds the loopback exemption; run revert-b2-loopback.ps1 first'
        }
        if (Get-B2ContainerProfileExists -ContainerName $ContainerName) {
            throw 'refusing to refresh the baseline while the AppContainer profile still exists; run remove-b2-container-profile.ps1 first'
        }
        Remove-Item -LiteralPath $baselinePath -Force
        Add-Step 'baseline' 'legacy baseline discarded (restoral state verified: no exemption, no profile)'
    }
}

if (-not (Test-Path $baselinePath)) {
    $null = New-B2Snapshot -OutFile $baselinePath -ContainerName $ContainerName
    Add-Step 'baseline' ('written to ' + $baselinePath)
} else {
    Add-Step 'baseline' 'reused existing baseline (schemaVersion 2; -RefreshBaseline after a full revert replaces it)'
}

$sid = New-B2ContainerSidIfMissing -ContainerName $ContainerName
Add-Step 'container-sid' $sid

$before = Get-B2ExemptionList
Add-Step 'exemption-list-before' (($before.sids -join ','))

if ($before.sids -contains $sid) {
    Add-Step 'apply' 'already exempted; nothing to do (idempotent)'
} else {
    $output = (& CheckNetIsolation.exe LoopbackExempt -a -p=$sid 2>&1 | Out-String).Trim()
    Add-Step 'apply-checknetisolation' $output
    $after = Get-B2ExemptionList
    if ($after.sids -notcontains $sid) {
        $target = @($after.sids) + $sid
        Add-Step 'apply-fallback' ('NetworkIsolationSetAppContainerConfig with ' + $target.Count + ' entries')
        $hr = [B2LoopbackNative]::SetExemptionList([string[]]$target)
        Add-Step 'apply-fallback-hresult' ('0x{0:X8}' -f $hr)
        $after = Get-B2ExemptionList
    }
    if ($after.sids -notcontains $sid) {
        Add-Step 'apply' 'FAILED (exemption not present after both mechanisms)'
        $failed = [pscustomobject]@{ ok = $false; sid = $sid; steps = $steps; exemptions = $after.sids }
        [IO.File]::WriteAllText($resultPath, ($failed | ConvertTo-Json -Depth 6), (New-Object System.Text.UTF8Encoding($false)))
        throw "loopback exemption could not be applied for $sid; nothing else was changed"
    }
    Add-Step 'apply' 'applied'
}

$final = Get-B2ExemptionList
$result = [pscustomobject]@{
    ok           = $true
    sid          = $sid
    exemptions   = $final.sids
    baselinePath = $baselinePath
    steps        = $steps
}
[IO.File]::WriteAllText($resultPath, ($result | ConvertTo-Json -Depth 6), (New-Object System.Text.UTF8Encoding($false)))
Write-Output ('apply-ok sid=' + $sid + ' exemptionCount=' + $final.sids.Count)
Write-Output ('apply-result: ' + $resultPath)

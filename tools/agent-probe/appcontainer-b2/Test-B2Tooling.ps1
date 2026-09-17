# Hermetic self-test for the B2 machine-discipline tooling: no elevation, no machine-state change, no
# network. Guards the invariants that a field run already broke once:
#   1. reading a BOM-less UTF-8 JSON artifact with non-ASCII content (PowerShell 5.1 Get-Content mis-decodes it)
#   2. the baseline schema-currency predicate
#   3. the strict-diff identity gate and the legacy-baseline informational line
#   4. profile existence must never be inferred from name-only SID derivation
#   5. normalize-encoding.ps1 must skip run-artifact trees and only rewrite in-scope files
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')

$checks = 0
$failures = 0
function Check([string]$Name, [scriptblock]$Body) {
    $script:checks++
    try {
        $result = & $Body
        if ($result -eq $true) { Write-Output ("ok   - " + $Name) }
        else { $script:failures++; Write-Output ("FAIL - " + $Name + " (returned " + $result + ")") }
    } catch {
        $script:failures++
        Write-Output ("FAIL - " + $Name + " (threw " + $_.Exception.Message + ")")
    }
}

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..\..')).Path
$scratch = Join-Path $repoRoot 'tmp\b2\selftest'
New-Item -ItemType Directory -Force -Path $scratch | Out-Null

Check 'BOM-less UTF-8 JSON with non-ASCII content round-trips' {
    $path = Join-Path $scratch 'unicode.json'
    $payload = [pscustomobject]@{ containerName = 'prfrail-b2-battery'; note = ('中文路径 ' + $env:USERNAME + ' \防火墙'); count = 3 }
    [IO.File]::WriteAllText($path, ($payload | ConvertTo-Json -Depth 4), (New-Object System.Text.UTF8Encoding($false)))
    $roundTrip = Read-B2JsonFile -Path $path
    ($roundTrip.note -eq $payload.note) -and ($roundTrip.count -eq 3)
}

Check 'missing file raises instead of returning garbage' {
    try {
        $null = Read-B2TextFile -Path (Join-Path $scratch 'does-not-exist.txt')
        $false
    } catch {
        $true
    }
}

Check 'schema-currency predicate rejects legacy and partial baselines' {
    $legacy = [pscustomobject]@{ containerName = 'x'; exemptions = @() }
    $partial = [pscustomobject]@{ schemaVersion = 2; containerName = 'x'; exemptions = @() }
    $full = [pscustomobject]@{ schemaVersion = 2; containerName = 'x'; profileExists = $false; exemptions = @() }
    (-not (Test-B2SnapshotSchemaCurrent $legacy)) -and (-not (Test-B2SnapshotSchemaCurrent $partial)) -and (Test-B2SnapshotSchemaCurrent $full)
}

Check 'strict diff refuses to compare a baseline recorded for another container' {
    $baseline = [pscustomobject]@{ schemaVersion = 2; containerName = 'container-a'; profileExists = $false; exemptions = @(); firewallRules = @() }
    $current = [pscustomobject]@{ schemaVersion = 2; containerName = 'container-b'; profileExists = $false; exemptions = @(); firewallRules = @() }
    $diff = @(Get-B2SnapshotDiff -Baseline $baseline -Current $current)
    ($diff.Count -eq 1) -and ($diff[0] -like 'baseline-container-mismatch:*')
}

Check 'legacy baseline produces only the informational line' {
    $baseline = [pscustomobject]@{ containerName = 'container-a'; exemptions = @(); firewallRules = @() }
    $current = [pscustomobject]@{ schemaVersion = 2; containerName = 'container-a'; profileExists = $false; exemptions = @(); firewallRules = @() }
    $diff = @(Get-B2SnapshotDiff -Baseline $baseline -Current $current)
    ($diff.Count -eq 1) -and ($diff[0] -like 'profile-existence-unverifiable:*')
}

Check 'profile existence is not inferred from name-only derivation' {
    (-not (Get-B2ContainerProfileExists -ContainerName 'prfrail-b2-selftest-nonexistent')) -and
        ((Get-B2ContainerSidByName -ContainerName 'prfrail-b2-selftest-nonexistent') -like 'S-1-15-2-*')
}

Check 'restoral verdict ignores informational lines only' {
    $informational = @('profile-existence-unverifiable: baseline predates schemaVersion 2')
    $real = @('exemption-added: S-1-15-2-1')
    $mixed = @('profile-existence-changed: False -> True (prfrail-b2-battery)')
    $informationalFailures = @(Get-B2RestoralFailures -DiffLines $informational)
    $realFailures = @(Get-B2RestoralFailures -DiffLines $real)
    $keptFailures = @(Get-B2RestoralFailures -DiffLines $mixed -ProfileKept $true)
    $notKeptFailures = @(Get-B2RestoralFailures -DiffLines $mixed -ProfileKept $false)
    ($informationalFailures.Count -eq 0) -and ($realFailures.Count -eq 1) -and ($keptFailures.Count -eq 0) -and ($notKeptFailures.Count -eq 1)
}

Check 'normalize-encoding skips run-artifact trees and only rewrites in-scope files' {
    # The scratch root lives beside the selftest scratch, not under it: a 'selftest' path segment is itself
    # excluded, which would skip the whole synthetic tree and make this check vacuous.
    $root = Join-Path $repoRoot 'tmp\b2\enc-selftest'
    $excludedDir = Join-Path $root 'excluded\candidate-runs'
    $scopedDir = Join-Path $root 'in-scope'
    New-Item -ItemType Directory -Force -Path $excludedDir, $scopedDir | Out-Null
    $payload = "line one`r`nline two`r`n"
    $noBom = New-Object System.Text.UTF8Encoding($false)
    $excludedFile = Join-Path $excludedDir 'payload.ps1'
    $scopedFile = Join-Path $scopedDir 'authored.ps1'
    [IO.File]::WriteAllText($excludedFile, $payload, $noBom)
    [IO.File]::WriteAllText($scopedFile, $payload, $noBom)
    $excludedBefore = [Convert]::ToBase64String([IO.File]::ReadAllBytes($excludedFile))
    $output = & powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot 'normalize-encoding.ps1') -Root 'tmp\b2\enc-selftest' 2>&1
    if ($LASTEXITCODE -ne 0) { throw ('normalize-encoding.ps1 exited with ' + $LASTEXITCODE) }
    $text = ($output | Out-String)
    $excludedUntouched = ($excludedBefore -eq [Convert]::ToBase64String([IO.File]::ReadAllBytes($excludedFile)))
    $bytes = [IO.File]::ReadAllBytes($scopedFile)
    $scopedNormalized = ($bytes.Length -ge 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF) -and
        (-not [Text.Encoding]::UTF8.GetString($bytes, 3, $bytes.Length - 3).Contains("`r`n"))
    $reported = ($text -match 'candidates=1 rewritten=1') -and ($text -match 'violations=0')
    $excludedUntouched -and $scopedNormalized -and $reported
}

Write-Output ("selftest: checks={0} failures={1} scratch={2}" -f $checks, $failures, $scratch)
if ($failures -gt 0) { exit 1 }

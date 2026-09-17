# Assemble the B2 evidence bundle: re-run the two decisive free batteries with captured output,
# copy the decisive artifacts (sanitized, repo-encoding normalized) and hash everything.
[CmdletBinding()]
param(
    [string]$DateStamp = '2026-09-17',
    [switch]$SkipBatteries
)

$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..')
$ScratchRoot = Join-Path (Get-Location).Path 'tmp\b2'

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

$evidenceRoot = Join-Path (Get-Location) ('docs\validation\evidence\b2-' + $DateStamp)
if (Test-Path -LiteralPath $evidenceRoot) { Remove-Item -LiteralPath $evidenceRoot -Recurse -Force }
New-Item -ItemType Directory -Force -Path $evidenceRoot | Out-Null
foreach ($sub in @('machine-ops', 'restorability', 'boundary', 'container-exec', 'proxy-decisions', 'token-exchange', 'measuring')) {
    New-Item -ItemType Directory -Force -Path (Join-Path $evidenceRoot $sub) | Out-Null
}
$userName = $env:USERNAME
function Copy-Sanitized([string]$source, [string]$destination, [switch]$NormalizeText) {
    if (-not (Test-Path -LiteralPath $source)) { Show 'missingArtifact' $source; return }
    $text = [IO.File]::ReadAllText($source)
    $text = $text.Replace($userName, '<user>').Replace($userName.ToLowerInvariant(), '<user>')
    if ($NormalizeText) {
        $text = $text.Replace("`r`n", "`n")
        $bytes = [Text.Encoding]::UTF8.GetBytes($text)
        $preamble = [Text.Encoding]::UTF8.GetPreamble()
        $all = New-Object byte[] ($preamble.Length + $bytes.Length)
        [Array]::Copy($preamble, 0, $all, 0, $preamble.Length)
        [Array]::Copy($bytes, 0, $all, $preamble.Length, $bytes.Length)
        [IO.File]::WriteAllBytes($destination, $all)
    } else {
        [IO.File]::WriteAllText($destination, $text, (New-Object System.Text.UTF8Encoding($false)))
    }
}

# --- machine-operation scripts (verbatim copies, normalized to repo encoding) -------------------
foreach ($script in @('ac-lib.ps1', 'b2-loopback-lib.ps1', 'apply-b2-loopback.ps1', 'revert-b2-loopback.ps1', 'verify-b2-loopback.ps1', 'remove-b2-container-profile.ps1', 'Invoke-B2LoopbackRehearsal.ps1', 'boundary-battery.ps1', 'container-candidate.ps1', 'container-modelcall.ps1', 'container-nodecall.ps1', 'path-probe.ps1', 'Test-B2Tooling.ps1', 'assemble-evidence.ps1', 'normalize-encoding.ps1', 'node-client.js')) {
    $target = Join-Path $evidenceRoot ('machine-ops\' + $script)
    if ($script -eq 'node-client.js') { Copy-Sanitized (Join-Path $PSScriptRoot $script) $target } else { Copy-Sanitized (Join-Path $PSScriptRoot $script) $target -NormalizeText }
    Show 'machineOp' $script
}

# --- restorability: rehearsal transcript (free, decisive) ---------------------------------------
if (-not $SkipBatteries) {
    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    Show 'runningRehearsal' 'apply -> verify -> revert -> verify -> apply'
    $rehearsalOut = Join-Path $evidenceRoot 'restorability\rehearsal.txt'
    $rehearsalText = (& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot 'Invoke-B2LoopbackRehearsal.ps1') 2>&1 | Out-String)
    $rehearsalFailed = $rehearsalText -match 'elevated PowerShell session'
    [IO.File]::WriteAllBytes($rehearsalOut, ([Text.Encoding]::UTF8.GetPreamble() + [Text.Encoding]::UTF8.GetBytes((($rehearsalText.Trim()) + "`n"))))
    if ($rehearsalFailed) {
        Show 'rehearsalStatus' 'requires an elevated session - run tmp\b2\Invoke-B2LoopbackRehearsal.ps1 as administrator to capture the full transcript'
    } else {
        Show 'rehearsalStatus' 'completed'
        foreach ($line in ($rehearsalText -split "`r?`n")) { if ($line.Trim() -ne '') { Write-Output ('  ' + $line.Trim()) } }
    }
    $verifyOut = Join-Path $evidenceRoot 'restorability\verify-current-state.txt'
    $verifyText = (& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot 'verify-b2-loopback.ps1') 2>&1 | Out-String)
    [IO.File]::WriteAllBytes($verifyOut, ([Text.Encoding]::UTF8.GetPreamble() + [Text.Encoding]::UTF8.GetBytes((($verifyText.Trim()) + "`n"))))
    Show 'verifyStatus' 'captured'

    Show 'runningBoundaryBattery' 'free boundary checks'
    $batteryOut = Join-Path $evidenceRoot 'boundary\boundary-battery.txt'
    & powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $PSScriptRoot 'boundary-battery.ps1') 2>&1 | Tee-Object -FilePath $batteryOut | ForEach-Object { Write-Output ('  ' + $_) }
    Show 'batteryTranscript' $batteryOut
    $ErrorActionPreference = $previousPreference
} else {
    Show 'batteries' 'skipped by request'
}

# --- container execution (pinned candidate starts inside the boundary) ---------------------------
$execRun = Join-Path $ScratchRoot 'candidate-runs\20260917-163249-exec-check'
Copy-Sanitized (Join-Path $execRun 'result.json') (Join-Path $evidenceRoot 'container-exec\result.json')
Copy-Sanitized (Join-Path $execRun 'stdout.txt') (Join-Path $evidenceRoot 'container-exec\stdout.txt')
Copy-Sanitized (Join-Path $execRun 'argv.json') (Join-Path $evidenceRoot 'container-exec\argv.json')
# pre-subst run: the candidate's ESM loader failed with EPERM lstat 'D:\' inside the container
$preSubstRun = Join-Path $ScratchRoot 'candidate-runs\20260917-162646-exec-check'
Copy-Sanitized (Join-Path $preSubstRun 'stderr.txt') (Join-Path $evidenceRoot 'container-exec\exec-check-pre-subst-stderr.txt')

# --- proxy decisions on real client traffic (deny-by-default, then allow-listed) -----------------
$denyRun = Join-Path $ScratchRoot 'candidate-runs\20260917-163359-smoke'
Copy-Sanitized (Join-Path $denyRun 'proxy.jsonl') (Join-Path $evidenceRoot 'proxy-decisions\deny-run.jsonl')
Copy-Sanitized (Join-Path $denyRun 'result.json') (Join-Path $evidenceRoot 'proxy-decisions\deny-run-result.json')
$allowRun = Join-Path $ScratchRoot 'candidate-runs\20260917-163515-smoke'
Copy-Sanitized (Join-Path $allowRun 'proxy.jsonl') (Join-Path $evidenceRoot 'proxy-decisions\allow-run.jsonl')
Copy-Sanitized (Join-Path $allowRun 'stderr.txt') (Join-Path $evidenceRoot 'proxy-decisions\allow-run-stderr.txt')
Copy-Sanitized (Join-Path $allowRun 'usage.json') (Join-Path $evidenceRoot 'proxy-decisions\allow-run-usage.json')
Copy-Sanitized (Join-Path $allowRun 'result.json') (Join-Path $evidenceRoot 'proxy-decisions\allow-run-result.json')
$candidateLog = Get-ChildItem (Join-Path $allowRun 'logs') -File -Filter '*.log' -ErrorAction SilentlyContinue | Select-Object -First 1
if ($candidateLog) { Copy-Sanitized $candidateLog.FullName (Join-Path $evidenceRoot 'proxy-decisions\candidate-session-log.txt') }

# --- token-exchange finding (WAF block for non-sanctioned clients) -------------------------------
$tokenRun = Get-ChildItem (Join-Path $ScratchRoot 'nodecall-runs') -Directory -ErrorAction SilentlyContinue | Sort-Object Name | Select-Object -Last 1
if ($tokenRun) {
    # token.json is deliberately NOT archived: on success it would contain a live Copilot session token.
    # The harness deletes it right after parsing (see container-nodecall.ps1).
    foreach ($name in @('node-result.txt', 'node-out.txt', 'token-err.txt')) {
        Copy-Sanitized (Join-Path $tokenRun.FullName $name) (Join-Path $evidenceRoot ('token-exchange\' + $name))
    }
    Copy-Sanitized (Join-Path $tokenRun.FullName 'proxy.jsonl') (Join-Path $evidenceRoot 'token-exchange\proxy.jsonl')
    Show 'tokenRun' $tokenRun.FullName
}

# --- measuring: container filesystem/path behaviour ----------------------------------------------
# probe13 matrix comes from path-probe.ps1; the directory-listing capture is copied as raw bytes
# because cmd writes it in the OEM code page (saving it as text would corrupt the bytes).
$probeOut = Join-Path $ScratchRoot 'path-probe'
if (Test-Path -LiteralPath (Join-Path $probeOut 'probe13-path-spelling.txt')) {
    Copy-Sanitized (Join-Path $probeOut 'probe13-path-spelling.txt') (Join-Path $evidenceRoot 'measuring\probe13-path-spelling.txt')
} else { Show 'missingArtifact' 'path-probe\probe13-path-spelling.txt' }
if (Test-Path -LiteralPath (Join-Path $probeOut 'probe10-dir-listing.txt')) {
    Copy-Sanitized (Join-Path $probeOut 'probe10-dir-listing.txt') (Join-Path $evidenceRoot 'measuring\probe10-dir-listing.txt')
} else { Show 'missingArtifact' 'path-probe\probe10-dir-listing.txt' }
if (Test-Path -LiteralPath (Join-Path $probeOut 'probe10-dir-listing.raw.txt')) {
    Copy-Item -LiteralPath (Join-Path $probeOut 'probe10-dir-listing.raw.txt') -Destination (Join-Path $evidenceRoot 'measuring\probe10-dir-listing.raw.txt') -Force
} else { Show 'missingArtifact' 'path-probe\probe10-dir-listing.raw.txt' }

# Profile-existence detection evidence: the derivation API is name-only, so existence is decided by the
# per-user registry storage key. The probe creates and deletes a throwaway profile and never touches the B2 one.
$detectProbe = Join-Path $ScratchRoot 'detect-probe.ps1'
if (Test-Path -LiteralPath $detectProbe) {
    Copy-Sanitized $detectProbe (Join-Path $evidenceRoot 'measuring\profile-detection-probe.ps1') -NormalizeText
} else { Show 'missingArtifact' 'detect-probe.ps1' }
$detectOutput = Join-Path $ScratchRoot 'state\profile-detection.txt'
if (-not (Test-Path -LiteralPath $detectOutput)) {
    $detectOutput = Join-Path $ScratchRoot 'profile-detection.txt'
}
if (Test-Path -LiteralPath $detectOutput) {
    Copy-Sanitized $detectOutput (Join-Path $evidenceRoot 'measuring\profile-detection.txt')
} else { Show 'missingArtifact' 'profile-detection transcript (run tmp\b2\detect-probe.ps1 with Tee-Object first)' }
$batteryRuns = Get-ChildItem (Join-Path $ScratchRoot 'boundary-battery') -Directory -ErrorAction SilentlyContinue | Sort-Object Name -Descending
if ($batteryRuns) {
    # newest run that still carries the container artifacts (older scratch runs may be partially cleaned)
    $latest = $batteryRuns | Where-Object { Test-Path -LiteralPath (Join-Path $_.FullName 'ws\result.txt') } | Select-Object -First 1
    if (-not $latest) { $latest = $batteryRuns | Select-Object -First 1 }
    foreach ($name in @('ws\result.txt', 'ws\child-console.txt')) {
        $sourcePath = Join-Path $latest.FullName $name
        $targetName = ($name -replace '\\', '-')
        Copy-Sanitized $sourcePath (Join-Path $evidenceRoot ('measuring\' + $targetName))
    }
}

# --- secret scan (fail loudly if anything resembling a token landed in the bundle) --------------
$suspect = @()
foreach ($file in Get-ChildItem $evidenceRoot -Recurse -File) {
    $text = [IO.File]::ReadAllText($file.FullName)
    if ($text -match 'gho_[A-Za-z0-9]{20,}' -or $text -match 'ghp_[A-Za-z0-9]{20,}' -or $text -match 'github_pat_[A-Za-z0-9_]{20,}' -or $text -match '\b[0-9a-f]{40}\b') {
        $suspect += $file.FullName.Replace($evidenceRoot + '\', '')
    }
}
Show 'secretScanSuspects' (($suspect -join ', ') + $(if ($suspect.Count -eq 0) { 'none' } else { '' }))
if ($suspect.Count -gt 0) { throw 'potential secret material found in the evidence bundle' }

# --- manifest + hashes ---------------------------------------------------------------------------
$manifest = New-Object System.Collections.Generic.List[string]
$manifest.Add('# B2 evidence bundle manifest (' + $DateStamp + ')')
$manifest.Add('')
$manifest.Add('All artifacts are sanitized copies: the local user name is replaced by <user>; no token value')
$manifest.Add('appears anywhere (verified by an automated scan for gho_/ghp_/github_pat_/40-hex patterns).')
$manifest.Add('Script copies are normalized to the repository encoding convention (UTF-8 with BOM, LF);')
$manifest.Add('log and result copies stay UTF-8 without BOM. Hashes below cover the stored bytes.')
$manifest.Add('')
foreach ($file in (Get-ChildItem $evidenceRoot -Recurse -File | Sort-Object FullName)) {
    $relative = $file.FullName.Replace($evidenceRoot + '\', '')
    $hash = (Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
    $manifest.Add(('- {0}  {1} bytes  sha256 {2}' -f $relative, $file.Length, $hash))
}
$manifestPath = Join-Path $evidenceRoot 'MANIFEST.md'
$manifestText = ($manifest -join "`n") + "`n"
[IO.File]::WriteAllBytes($manifestPath, ([Text.Encoding]::UTF8.GetPreamble() + [Text.Encoding]::UTF8.GetBytes($manifestText)))

$sums = New-Object System.Collections.Generic.List[string]
foreach ($file in (Get-ChildItem $evidenceRoot -Recurse -File | Sort-Object FullName)) {
    $relative = $file.FullName.Replace($evidenceRoot + '\', '')
    $sums.Add(((Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant() + '  ' + $relative))
}
[IO.File]::WriteAllText((Join-Path $evidenceRoot 'SHA256SUMS.txt'), (($sums -join "`n") + "`n"), (New-Object System.Text.UTF8Encoding($false)))

Show 'evidenceRoot' $evidenceRoot
Show 'fileCount' (Get-ChildItem $evidenceRoot -Recurse -File).Count

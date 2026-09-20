# New-B3aEvidence.ps1 - assemble the B3a evidence bundle (design S7) and verify it.
#
# Bundle layout (B2 format):
#   <OutputRoot>\b3a-<date>\   (default: <repo>\docs\validation\evidence\b3a-<date>)
#     README.md               reproduction commands + premise + coverage (BOM+LF)
#     MANIFEST.md             '- <path>  N bytes  sha256 <hex>' (BOM+LF)
#     SHA256SUMS.txt          '<sha256 hex>  path\relative' (no BOM+LF, backslash)
#     calibration-verdict.json (if present)
#     selftest.txt            archived Test-B3aRig.ps1 output (checks=N failures=0)
#     rig\                    frozen byte copies of all rig scripts (incl. the
#                             reused 11-shutdown-audit.ps1)
#     journal\journal.jsonl + journal\chain-verify.txt
#     rounds\rNN\             plan.json, cut.json, write-probe.txt, audit.txt,
#                             inventory-1/2.json, dirty.txt, reconcile.json
#                             (missing files are reported, not fabricated)
#     session\session.log
#
# Assembly refuses to run when the journal chain fails or the hermetic self-test
# fails. After assembly the whole bundle is re-hashed against SHA256SUMS.txt;
# exit 0 = 0 mismatches.
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

param(
    [string]$SessionDir = '',
    [string]$RigRoot = 'D:\VirtualBox VMs\Win11B3\b3a-runs',
    [string]$OutputRoot = '',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep'
)

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'b3a-journal-lib.ps1')

if ($MyInvocation.InvocationName -eq '.') { return }

# Fail closed: an incomplete bundle must never look complete.
$ErrorActionPreference = 'Stop'

$DateTag = Get-Date -Format 'yyyy-MM-dd'
if (-not $SessionDir) { $SessionDir = Join-Path $RigRoot $DateTag }
if (-not (Test-Path -LiteralPath $SessionDir)) { 'session directory missing: ' + $SessionDir; exit 1 }
if ($OutputRoot) {
    $BundleDir = Join-Path $OutputRoot ('b3a-' + $DateTag)
} else {
    $BundleDir = Join-Path (Join-Path $PSScriptRoot '..\..') ('docs\validation\evidence\b3a-' + $DateTag)
}
# Normalize: Join-Path keeps '..' segments literally, and the MANIFEST/SHA256SUMS
# loop derives relative paths from $BundleDir.Length - a literal '..\..' path is
# longer than the real file paths, which made Substring throw (2026-09-20, found on
# the script's first real run).
$BundleDir = [System.IO.Path]::GetFullPath($BundleDir)
if (Test-Path -LiteralPath $BundleDir) { 'bundle already exists (refusing to overwrite): ' + $BundleDir; exit 1 }

$JournalPath = Join-Path $SessionDir 'journal.jsonl'
$chain = Test-B3aJournalChain -JournalPath $JournalPath
if (-not $chain.ok) {
    'journal chain broken at line {0}: {1} - evidence bundle refused' -f $chain.firstBadLine, $chain.firstError
    exit 1
}

$copied = New-Object System.Collections.Generic.List[string]
$missing = New-Object System.Collections.Generic.List[string]
# Required artifacts: a bundle that silently lacks them would still look complete.
$missingRequired = New-Object System.Collections.Generic.List[string]

function Copy-B3aBundleFile {
    param([Parameter(Mandatory = $true)][string]$Source, [Parameter(Mandatory = $true)][string]$Dest,
          [switch]$Required)
    if (-not (Test-Path -LiteralPath $Source)) {
        if ($Required) { $missingRequired.Add($Dest) } else { $missing.Add($Dest) }
        return
    }
    $dir = Split-Path -Parent $Dest
    if (-not (Test-Path -LiteralPath $dir)) { [void](New-Item -ItemType Directory -Path $dir -Force) }
    [System.IO.File]::Copy($Source, $Dest, $true)
    $copied.Add($Dest)
}

# rig scripts (frozen byte copies)
foreach ($f in @(
    'b3a-rig-lib.ps1', 'b3a-journal-lib.ps1', 'New-B3aSession.ps1', 'Invoke-B3aRound.ps1',
    'Invoke-B3aCalibration.ps1', 'Get-B3aReconcile.ps1', 'Invoke-B3aReconcile.ps1',
    'New-B3aEvidence.ps1', 'Test-B3aRig.ps1'
)) {
    Copy-B3aBundleFile -Source (Join-Path $PSScriptRoot $f) -Dest (Join-Path $BundleDir ('rig\' + $f)) -Required
}
foreach ($f in @('12-round-inventory.ps1', '13-round-write.ps1')) {
    Copy-B3aBundleFile -Source (Join-Path $PSScriptRoot ('guest\' + $f)) -Dest (Join-Path $BundleDir ('rig\' + $f)) -Required
}
Copy-B3aBundleFile -Source (Join-Path $GuestPrepDir '11-shutdown-audit.ps1') -Dest (Join-Path $BundleDir 'rig\11-shutdown-audit.ps1') -Required

# journal
Copy-B3aBundleFile -Source $JournalPath -Dest (Join-Path $BundleDir 'journal\journal.jsonl') -Required
$chainText = 'journal chain verification: ok=' + $chain.ok + ' records=' + $chain.count
if (-not $chain.ok) { $chainText += ' firstBadLine=' + $chain.firstBadLine + ' error=' + $chain.firstError }
Write-B3aTextFile -Path (Join-Path $BundleDir 'journal\chain-verify.txt') -Text $chainText
$copied.Add((Join-Path $BundleDir 'journal\chain-verify.txt'))

# rounds: every round that started (valid or void) must have produced reconcile.json;
# the remaining artifacts exist only for the steps a void round actually reached.
if (Test-Path -LiteralPath (Join-Path $SessionDir 'rounds')) {
    $roundDirs = @(Get-ChildItem -LiteralPath (Join-Path $SessionDir 'rounds') -Directory | Sort-Object Name)
    if ($roundDirs.Count -eq 0) { $missingRequired.Add('rounds\<no round directory>') }
    foreach ($rd in $roundDirs) {
        foreach ($art in @('plan.json', 'cut.json', 'write-probe.txt', 'write-stderr.txt', 'audit.txt', 'audit-pre.txt',
            'inventory-1.json', 'inventory-2.json', 'dirty.txt', 'reconcile.json')) {
            $isRequired = ($art -eq 'reconcile.json')
            Copy-B3aBundleFile -Source (Join-Path $rd.FullName $art) `
                -Dest (Join-Path $BundleDir ('rounds\' + $rd.Name + '\' + $art)) -Required:$isRequired
        }
    }
} else {
    $missingRequired.Add('rounds\<no round directory>')
}

# session log + calibration verdict
Copy-B3aBundleFile -Source (Join-Path $SessionDir 'session.log') -Dest (Join-Path $BundleDir 'session\session.log') -Required
Copy-B3aBundleFile -Source (Join-Path $SessionDir 'calibration-verdict.json') -Dest (Join-Path $BundleDir 'calibration-verdict.json') -Required

# hermetic self-test output (design S7: the checks=N failures=0 output is archived)
$psExe = Join-Path $env:WINDIR 'System32\WindowsPowerShell\v1.0\powershell.exe'
$self = Invoke-B3aProcess -FileName $psExe -Arguments ('-NoProfile -ExecutionPolicy Bypass -File "{0}"' -f (Join-Path $PSScriptRoot 'Test-B3aRig.ps1')) -TimeoutS 300
Write-B3aRawTextFile -Path (Join-Path $BundleDir 'selftest.txt') -Bytes $self.stdoutBytes
if ($self.exitCode -ne 0) { 'self-test FAILED (exit ' + $self.exitCode + ') - output archived as selftest.txt'; exit 1 }
$copied.Add((Join-Path $BundleDir 'selftest.txt'))

# verification/ - artifacts that make the report's traceability claims checkable
# (independent review, 2026-09-20: "frozen bytes == bytes that ran" and
# "analyzer = 0" were claims without a bundled artifact).
[void](New-Item -ItemType Directory -Path (Join-Path $BundleDir 'verification') -Force)

# 1. rig copies vs the repository, plus each script's mtime against the start of the
#    calibration as recorded in session.log.
$sessionLogPath = Join-Path $SessionDir 'session.log'
$calibStartText = ''
$calibStartDt = $null
if (Test-Path -LiteralPath $sessionLogPath) {
    foreach ($line in [System.IO.File]::ReadAllLines($sessionLogPath)) {
        $m = [regex]::Match($line, '^\[(\d\d:\d\d:\d\d)\] === B3a calibration begin')
        if ($m.Success) { $calibStartText = $m.Groups[1].Value }
    }
    if ($calibStartText) {
        $calibStartDt = (Get-Item -LiteralPath $sessionLogPath).LastWriteTime.Date.Add([TimeSpan]::Parse($calibStartText))
    }
}
$rigNames = @(
    'b3a-rig-lib.ps1', 'b3a-journal-lib.ps1', 'Get-B3aReconcile.ps1', 'Invoke-B3aReconcile.ps1',
    'Invoke-B3aRound.ps1', 'Invoke-B3aCalibration.ps1', 'New-B3aSession.ps1', 'Test-B3aRig.ps1',
    'New-B3aEvidence.ps1', '12-round-inventory.ps1', '13-round-write.ps1'
)
$verifyLines = New-Object System.Collections.Generic.List[string]
if ($calibStartDt) {
    $verifyLines.Add('calibration began at ' + $calibStartText + ' (last "=== B3a calibration begin" line in session.log)')
} else {
    $verifyLines.Add('WARNING: no "B3a calibration begin" line found in session.log; the mtime column is unjudged')
}
foreach ($n in $rigNames) {
    $repoPath = Join-Path $PSScriptRoot $n
    if (-not (Test-Path -LiteralPath $repoPath)) { $repoPath = Join-Path (Join-Path $PSScriptRoot 'guest') $n }
    $bundlePath = Join-Path $BundleDir ('rig\' + $n)
    $repoHash = Get-B3aSha256Hex -Path $repoPath
    $bundleHash = Get-B3aSha256Hex -Path $bundlePath
    $repoMtime = (Get-Item -LiteralPath $repoPath).LastWriteTime
    $verdict = 'unjudged'
    if ($calibStartDt) { if ($repoMtime -le $calibStartDt) { $verdict = 'before-calibration-start' } else { $verdict = 'AFTER-CALIBRATION-START' } }
    $verifyLines.Add(('{0,-28} identical={1,-6} repoMtime={2} {3}' -f `
        $n, ($repoHash -eq $bundleHash), $repoMtime.ToString('yyyy-MM-dd HH:mm:ss'), $verdict))
    $verifyLines.Add(('{0,-28}   repo={1}' -f '', $repoHash))
    $verifyLines.Add(('{0,-28}   bundle={1}' -f '', $bundleHash))
}
Write-B3aTextFile -Path (Join-Path $BundleDir 'verification\rig-vs-repo.txt') -Text ($verifyLines -join "`n")
$copied.Add((Join-Path $BundleDir 'verification\rig-vs-repo.txt'))

# 2. raw static-analysis output over the frozen rig directory
$analyzerProbe = Join-Path $env:TEMP ('b3a-analyzer-' + [guid]::NewGuid().ToString('N') + '.ps1')
$analyzerText = @"
`$r = @(Invoke-ScriptAnalyzer -Path '$PSScriptRoot' -Recurse -Severity Warning,Error)
'findings=' + `$r.Count
`$r | ForEach-Object { 'FINDING ' + `$_.RuleName + ' ' + (Split-Path `$_.ScriptName -Leaf) + ':' + `$_.Line + ' ' + `$_.Message }
"@
Write-B3aTextFile -Path $analyzerProbe -Text $analyzerText -WithBom
$analyzerRun = Invoke-B3aProcess -FileName $psExe -Arguments ('-NoProfile -ExecutionPolicy Bypass -File "{0}"' -f $analyzerProbe) -TimeoutS 600
Write-B3aRawTextFile -Path (Join-Path $BundleDir 'verification\analyzer.txt') -Bytes $analyzerRun.stdoutBytes
[void](Remove-Item -LiteralPath $analyzerProbe -Force -ErrorAction SilentlyContinue)
$copied.Add((Join-Path $BundleDir 'verification\analyzer.txt'))
foreach ($l in @([System.IO.File]::ReadAllLines((Join-Path $BundleDir 'verification\analyzer.txt')))) { $l }

# README.md (BOM+LF)
$readme = @'
# B3a evidence bundle <DATE>

Reproduction commands (host PowerShell 5.1; see docs/t027/B3A_RIG_DESIGN.md S10):

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aSession.ps1 `
  -VmName Win11B3 -CredFile "D:\VirtualBox VMs\Win11B3\host-only\cred.txt"

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aCalibration.ps1 `
  -VmName Win11B3 -CredFile "D:\VirtualBox VMs\Win11B3\host-only\cred.txt" -N 5

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aRound.ps1 `
  -Round 3 -Candidate negative-control

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Test-B3aRig.ps1

powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aEvidence.ps1
```

Premise: the full premise and coverage statement is docs/t027/B3A_RIG_DESIGN.md S1
(hypervisor VirtualBox 7.2.18 r175117, guest Windows 11 Pro 25H2 build 26200.8037
zh-CN retail unactivated, VBoxSVGA graphics adapter, no guest NIC, SATA/AHCI
useHostIOCache=false, VDI backend, host NVMe, NTFS, frozen snapshot 'pristine').
Conclusions cover only this premise (design S1.8).
Guest clocks are secondary evidence only; the journal hostTs is the authority.
Calibration results prove only the rig's discriminating power, never a product
candidate (C2/C3 conclusions belong to B3b).

Encoding: .ps1/.md files carry UTF-8 BOM + LF; .json/.txt are UTF-8 without BOM
+ LF (CODING_CONVENTIONS.md).
'@
$readme = $readme.Replace('<DATE>', $DateTag)
Write-B3aTextFile -Path (Join-Path $BundleDir 'README.md') -Text $readme -WithBom
$copied.Add((Join-Path $BundleDir 'README.md'))

# MANIFEST.md + SHA256SUMS.txt over every stored byte
$files = @(Get-ChildItem -LiteralPath $BundleDir -Recurse -File | Sort-Object FullName)
$manifestLines = New-Object System.Collections.Generic.List[string]
$sumsLines = New-Object System.Collections.Generic.List[string]
foreach ($fi in $files) {
    # $BundleDir is normalized above, so every file under it is strictly longer;
    # fail loudly rather than letting Substring throw an opaque range error.
    if (-not $fi.FullName.StartsWith($BundleDir, [System.StringComparison]::OrdinalIgnoreCase)) {
        'BLOCKER: file outside the bundle root: ' + $fi.FullName
        exit 1
    }
    $rel = $fi.FullName.Substring($BundleDir.Length).TrimStart('\', '/')
    $hash = Get-B3aSha256Hex -Path $fi.FullName
    $manifestLines.Add(('- {0}  {1} bytes  sha256 {2}' -f $rel, $fi.Length, $hash))
    $sumsLines.Add(('{0}  {1}' -f $hash, $rel.Replace('/', '\')))
}
Write-B3aTextFile -Path (Join-Path $BundleDir 'MANIFEST.md') -Text ($manifestLines -join "`n") -WithBom
Write-B3aTextFile -Path (Join-Path $BundleDir 'SHA256SUMS.txt') -Text ($sumsLines -join "`n")
$copied.Add((Join-Path $BundleDir 'MANIFEST.md'))
$copied.Add((Join-Path $BundleDir 'SHA256SUMS.txt'))

# whole-bundle verification: recompute every hash against SHA256SUMS.txt
$mismatches = 0
$checked = 0
foreach ($sumLine in @((Get-B3aFileTextUtf8 -Path (Join-Path $BundleDir 'SHA256SUMS.txt')) -split "`n")) {
    if (-not $sumLine.Trim()) { continue }
    $parts = $sumLine -split '\s+', 2
    if ($parts.Count -lt 2) { continue }
    $expected = $parts[0]
    $relPath = $parts[1].Trim().Replace('\', [System.IO.Path]::DirectorySeparatorChar)
    $actual = Get-B3aSha256Hex -Path (Join-Path $BundleDir $relPath)
    $checked++
    if ($actual -ne $expected) { $mismatches++; 'MISMATCH: ' + $relPath }
}
''
'bundle: {0}' -f $BundleDir
'files copied: {0}; files in bundle: {1}' -f $copied.Count, $files.Count
'missing optional artifacts (reported, not fabricated): {0}' -f $missing.Count
foreach ($m in $missing) { '  MISSING-OPTIONAL: ' + $m }
'missing required artifacts: {0}' -f $missingRequired.Count
foreach ($m in $missingRequired) { '  MISSING-REQUIRED: ' + $m }
'sha256sum check: {0} mismatches ({1} files checked)' -f $mismatches, $checked
if ($missingRequired.Count -gt 0) {
    'evidence bundle REFUSED: required artifacts are missing'
    exit 1
}
if ($mismatches -ne 0) { exit 1 }
exit 0

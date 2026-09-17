# Normalize text artifacts to the repository encoding convention and regenerate the evidence manifest.
[CmdletBinding()]
param(
    [string]$Root = 'docs\validation',
    # Directory names (wildcards allowed) that hold run artifacts / extracted payloads. They are never
    # normalized: they either belong to an external program (a candidate's extracted package) or are
    # scratch output copied into the evidence bundle, which is normalized separately. Skipping them keeps a
    # scratch-root scan fast - rewriting thousands of payload files looks like a hang to the operator.
    [string[]]$ExcludeDirNames = @('candidate-runs', 'nodecall-runs', 'modelcall-runs', 'boundary-battery', 'path-probe', 'selftest', 'exec-probe*', 'exec-isolation', 'ac-workspace*')
)

$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..')
$rootPath = (Resolve-Path -LiteralPath $Root).Path
$bomExtensions = @('.md', '.ps1')
$changed = 0
function Test-ExcludedPath([string]$Path) {
    foreach ($dirName in $Path.Split('\')) {
        foreach ($pattern in $ExcludeDirNames) {
            if ($dirName -like $pattern) { return $true }
        }
    }
    return $false
}
Write-Output ("normalizing {0} (skipping run-artifact trees: {1})" -f $rootPath, ($ExcludeDirNames -join ', '))
$textExtensions = @('.md', '.ps1', '.txt', '.json', '.jsonl', '.js', '.go')
$all = @(Get-ChildItem $rootPath -Recurse -File)
$candidates = @($all | Where-Object {
        $_.Extension -in $textExtensions -and
        $_.Name -ne 'probe10-dir-listing.raw.txt' -and
        -not (Test-ExcludedPath $_.FullName)
    })
$skipped = $all.Count - $candidates.Count
Write-Output ("scanned files: {0}; content candidates: {1}; skipped (excluded trees / non-text / raw): {2}" -f $all.Count, $candidates.Count, $skipped)
$index = 0
foreach ($file in $candidates) {
    $index++
    if ($index % 100 -eq 0) { Write-Output ("progress: {0}/{1}" -f $index, $candidates.Count) }
    $text = [IO.File]::ReadAllText($file.FullName)
    $normalized = $text.Replace("`r`n", "`n")
    $wantBom = $bomExtensions -contains $file.Extension.ToLowerInvariant()
    $bytes = [Text.Encoding]::UTF8.GetBytes($normalized)
    if ($wantBom) {
        $preamble = [Text.Encoding]::UTF8.GetPreamble()
        $allBytes = New-Object byte[] ($preamble.Length + $bytes.Length)
        [Array]::Copy($preamble, 0, $allBytes, 0, $preamble.Length)
        [Array]::Copy($bytes, 0, $allBytes, $preamble.Length, $bytes.Length)
    } else {
        $allBytes = $bytes
    }
    $existing = [IO.File]::ReadAllBytes($file.FullName)
    # Compare with native hashing, not a script-level per-byte loop: for large verbatim payload files the
    # loop dominated the runtime and made an otherwise healthy scan look stuck.
    $same = $false
    if ($existing.Length -eq $allBytes.Length) {
        $sha = [Security.Cryptography.SHA256]::Create()
        try {
            $same = ([Convert]::ToBase64String($sha.ComputeHash($existing)) -eq [Convert]::ToBase64String($sha.ComputeHash($allBytes)))
        } finally {
            $sha.Dispose()
        }
    }
    if (-not $same) { [IO.File]::WriteAllBytes($file.FullName, $allBytes); $changed++ }
}
Write-Output ("normalization done: candidates={0} rewritten={1} unchanged={2}" -f $candidates.Count, $changed, ($candidates.Count - $changed))

$evidenceRoot = Join-Path $rootPath 'evidence\b2-2026-09-17'
if (Test-Path -LiteralPath $evidenceRoot) {
    $artifacts = Get-ChildItem $evidenceRoot -Recurse -File | Where-Object { $_.Name -notin @('MANIFEST.md', 'SHA256SUMS.txt') } | Sort-Object FullName
    $manifest = New-Object System.Collections.Generic.List[string]
    $manifest.Add('# B2 evidence bundle manifest (2026-09-17)')
    $manifest.Add('')
    $manifest.Add('All artifacts are sanitized copies: the local user name is replaced by <user>; no token value')
    $manifest.Add('appears anywhere (verified by an automated scan for gho_/ghp_/github_pat_/40-hex patterns).')
    $manifest.Add('The machine-ops/ copies are the archived snapshot of the maintained scripts now living in')
    $manifest.Add('tools/agent-probe/appcontainer-b2/ (same bytes at archive time); scratch output of those scripts')
    $manifest.Add('goes to the git-ignored tmp/b2/ tree, never next to the versioned scripts.')
    $manifest.Add('Every text artifact is stored as UTF-8 with LF: repository encoding convention applies, so the')
    $manifest.Add('scripts (*.ps1) carry a BOM while logs, results and JSON do not. Exception:')
    $manifest.Add('measuring/probe10-dir-listing.raw.txt keeps the raw OEM code-page bytes captured from cmd; its')
    $manifest.Add('decoded counterpart is measuring/probe10-dir-listing.txt. The live scratch copies under tmp/')
    $manifest.Add('are normalized identically (normalize-encoding.ps1 -Root tmp\b2), so a live copy and its archived')
    $manifest.Add('counterpart differ only by the user-name sanitization. Hashes below cover the stored bytes.')
    $manifest.Add('')
    foreach ($file in $artifacts) {
        $relative = $file.FullName.Replace($evidenceRoot + '\', '')
        $manifest.Add(('- {0}  {1} bytes  sha256 {2}' -f $relative, $file.Length, (Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant()))
    }
    $manifestText = ($manifest -join "`n") + "`n"
    [IO.File]::WriteAllBytes((Join-Path $evidenceRoot 'MANIFEST.md'), ([Text.Encoding]::UTF8.GetPreamble() + [Text.Encoding]::UTF8.GetBytes($manifestText)))
    $sums = New-Object System.Collections.Generic.List[string]
    foreach ($file in $artifacts) {
        $sums.Add(((Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant() + '  ' + $file.FullName.Replace($evidenceRoot + '\', '')))
    }
    [IO.File]::WriteAllText((Join-Path $evidenceRoot 'SHA256SUMS.txt'), (($sums -join "`n") + "`n"), (New-Object System.Text.UTF8Encoding($false)))
    Write-Output ('artifacts in bundle: ' + $artifacts.Count)
}
Write-Output ('files rewritten: ' + $changed)
# Closing audit: report only violations, and never walk the excluded artifact trees. The previous version
# dumped one line per .md/.ps1 (thousands of lines for a scratch root), which hid the real result and made
# the command look like it had hung.
$audited = @($all | Where-Object { $_.Extension -in $bomExtensions -and $_.Name -ne 'probe10-dir-listing.raw.txt' -and -not (Test-ExcludedPath $_.FullName) } | Sort-Object FullName)
$offenders = New-Object System.Collections.Generic.List[string]
foreach ($file in $audited) {
    $bytes = [IO.File]::ReadAllBytes($file.FullName)
    $bom = ($bytes.Length -ge 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF)
    if (-not $bom) { $offenders.Add($file.FullName.Replace($rootPath + '\', '') + ' missing BOM') }
    elseif ([Text.Encoding]::UTF8.GetString($bytes, 3, $bytes.Length - 3).Contains("`r`n")) { $offenders.Add($file.FullName.Replace($rootPath + '\', '') + ' contains CRLF') }
}
Write-Output ("bom/crlf audit: checked={0} violations={1} (per-file BOM required for {2})" -f $audited.Count, $offenders.Count, ($bomExtensions -join ', '))
foreach ($offender in $offenders) { Write-Output ('  ' + $offender) }

# Path-spelling probe (free): reproduces the two measurements cited in the B2 report.
#   output 1: D: vs R: spelling matrix for one directory (exists/list/create)
#   output 2: directory enumeration via the subst spelling vs the host spelling
# The probe never creates machine state: it resolves the container SID derive-only and requires that the
# loopback exemption is already applied by apply-b2-loopback.ps1.
[CmdletBinding()]
param(
    [string]$ContainerName = 'prfrail-b2-battery',
    [string]$OutDir = ''
)

$ErrorActionPreference = 'Continue'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..')
. (Join-Path $PSScriptRoot 'ac-lib.ps1')
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
$ScratchRoot = Join-Path (Get-Location).Path 'tmp\b2'

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

if ($OutDir -eq '') { $OutDir = Join-Path $ScratchRoot 'path-probe' }
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$sid = Get-B2ContainerSidExisting -ContainerName $ContainerName
$exemptRaw = (& CheckNetIsolation.exe LoopbackExempt -s 2>&1 | Out-String)
if (-not ($exemptRaw -match [regex]::Escape($sid))) { throw 'the container SID does not hold the loopback exemption; run apply-b2-loopback.ps1 first' }

$runDir = Join-Path $OutDir (Get-Date -Format 'yyyyMMdd-HHmmss')
$ws = Join-Path $runDir 'ws'
New-Item -ItemType Directory -Force -Path $ws | Out-Null
Set-Content -LiteralPath (Join-Path $ws 'marker.txt') -Value 'marker-b2' -NoNewline
foreach ($dir in @($runDir, $ws)) { $null = (& icacls $dir /grant ("*" + $sid + ':(OI)(CI)M') 2>&1) }

$drive = 'R:'
$null = (& subst $drive $runDir 2>&1)
Show 'substDrive' ($drive + ' -> ' + $runDir)
trap { if ($drive) { $null = (& subst $drive /D 2>&1) }; Write-Error ('path probe aborted: ' + ($_ | Out-String)); exit 1 }

try {
    $envPairs = @()
    foreach ($key in [Environment]::GetEnvironmentVariables('Process').Keys) { $envPairs += ($key + '=' + [Environment]::GetEnvironmentVariable($key, 'Process')) }

    # --- output 1: exists/list/create matrix for the same directory under both spellings -------------
    $matrix = (
        "`$lines = @(); " +
        "`$lines += 'runDirD=$runDir'; " +
        "`$lines += 'runDirR=" + ($drive + '\') + "'; " +
        "`$lines += 'existsD=' + [IO.Directory]::Exists('$runDir'); " +
        "`$lines += 'existsR=' + [IO.Directory]::Exists('" + ($drive + '\') + "'); " +
        "try { `$c = [IO.Directory]::GetFileSystemEntries('$runDir'); `$lines += 'listD-count=' + `$c.Count } catch { `$lines += 'listD=DENIED' }; " +
        "try { `$c2 = [IO.Directory]::GetFileSystemEntries('" + ($drive + '\') + "'); `$lines += 'listR-count=' + `$c2.Count } catch { `$lines += 'listR=DENIED' }; " +
        "try { [IO.File]::WriteAllText('$runDir\create-d.txt','x'); `$lines += 'createD=ok' } catch { `$lines += 'createD=DENIED' }; " +
        "try { [IO.File]::WriteAllText('" + ($drive + '\create-r.txt') + "','x'); `$lines += 'createR=ok' } catch { `$lines += 'createR=DENIED' }; " +
        "[IO.File]::WriteAllText('" + ($drive + '\probe13.txt') + "', (`$lines -join [Environment]::NewLine))"
    )
    $encoded = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($matrix))
    $cmd = 'C:\Windows\System32\cmd.exe /c powershell.exe -NoProfile -NonInteractive -Command ' + $encoded
    # -Command with a base64 payload is not valid; use the short form instead by writing the script to disk.
    $scriptPath = Join-Path $runDir 'matrix.ps1'
    [IO.File]::WriteAllText($scriptPath, ($matrix -replace '; ', "`n"), (New-Object System.Text.UTF8Encoding($false)))
    $cmd = 'C:\Windows\System32\cmd.exe /c powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -File ' + ($drive + '\matrix.ps1')
    $lastError = 0; $childPid = 0
    $code = [AcLauncher2]::LaunchWithEnv($sid, @(), $cmd, ($drive + '\'), [string[]]$envPairs, [ref]$lastError, [ref]$childPid)
    Show 'matrixExitCode' $code
    $matrixPath = Join-Path $runDir 'probe13.txt'
    if (Test-Path -LiteralPath $matrixPath) {
        $text = ([IO.File]::ReadAllText($matrixPath) -split "`r?`n") -join "`n"
        [IO.File]::WriteAllText((Join-Path $OutDir 'probe13-path-spelling.txt'), ($text + "`n"), (New-Object System.Text.UTF8Encoding($false)))
        Show 'probe13' (Join-Path $OutDir 'probe13-path-spelling.txt')
    } else { Show 'probe13' 'missing' }

    # --- output 2: directory enumeration through both spellings -------------------------------------
    $listing = @(
        'echo === dir R:\ws ===', 'dir R:\ws',
        'echo === dir D:\ws ===', ('dir ' + $ws),
        'echo === marker via R ===', ('type ' + ($drive + '\ws\marker.txt')),
        'echo === marker via D ===', ('type ' + (Join-Path $ws 'marker.txt'))
    ) -join ' & '
    $listingPath = Join-Path $OutDir 'probe10-dir-listing.txt'
    $cmd2 = 'C:\Windows\System32\cmd.exe /c (' + $listing + ') > ' + ($drive + '\listing.txt') + ' 2>&1'
    $lastError = 0; $childPid = 0
    $code2 = [AcLauncher2]::LaunchWithEnv($sid, @(), $cmd2, ($drive + '\'), [string[]]$envPairs, [ref]$lastError, [ref]$childPid)
    Show 'listingExitCode' $code2
    $listingSource = Join-Path $runDir 'listing.txt'
    if (Test-Path -LiteralPath $listingSource) {
        # cmd writes in the OEM code page; keep the raw bytes and add a decoded, human-readable copy
        $rawBytes = [IO.File]::ReadAllBytes($listingSource)
        [IO.File]::WriteAllBytes((Join-Path $OutDir 'probe10-dir-listing.raw.txt'), $rawBytes)
        $consoleEncoding = [Console]::OutputEncoding
        $decoded = $consoleEncoding.GetString($rawBytes)
        if ($decoded -notmatch 'ws') { $decoded = [Text.Encoding]::GetEncoding(936).GetString($rawBytes) }
        [IO.File]::WriteAllText($listingPath, ($decoded -replace "`r`n", "`n"), (New-Object System.Text.UTF8Encoding($false)))
        Show 'probe10' $listingPath
        Show 'probe10raw' (Join-Path $OutDir 'probe10-dir-listing.raw.txt')
    } else { Show 'probe10' 'missing' }
} finally {
    $null = (& subst $drive /D 2>&1)
    Show 'substRestored' (-not (Test-Path -LiteralPath ($drive + '\')))
}
Show 'runDir' $runDir

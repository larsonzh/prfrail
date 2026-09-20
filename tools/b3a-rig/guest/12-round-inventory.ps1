# 12-round-inventory.ps1 - deterministic restart inventory for the B3a rig (design S2.3).
#
# Emits ONE JSON object to stdout (UTF-8, no BOM, fixed key order, files sorted
# by relative path with the ordinal comparer). The inventory is read-only: it
# writes nothing on the guest. All guest timestamps are labelled secondary
# ('guest-clock') and are never evidence.
#
# Repeatability contract (design S2.3): the rig runs this twice after every
# reboot and the two stdout byte sequences must be identical. To make that
# possible the stdout JSON contains only boot-stable values; the live guest
# clock and uptime (which change between the two runs) go to STDERR as
# secondary diagnostics. Implementation note: the design's boot section lists
# lastBoot + uptime; uptime cannot be byte-identical across two runs, so it is
# reported on stderr instead - to be folded back into the design deliberately.
#
# fsutil dirty query <drive>: is recorded in the 'dirty' field; on failure the
# value is 'unavailable'.
#
# Exit 0 always unless the process itself dies. PowerShell 5.1 compatible.
# Encoding: UTF-8 with BOM + LF. All output English.

param(
    [Parameter(Mandatory = $true)][string]$RoundRoot,
    [Parameter(Mandatory = $true)][string]$ProbeRoot
)

$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = New-Object System.Text.UTF8Encoding($false)

function Get-Sha256Hex {
    param([Parameter(Mandatory = $true)][byte[]]$Bytes)
    $sha = [System.Security.Cryptography.SHA256]::Create()
    return ([System.BitConverter]::ToString($sha.ComputeHash($Bytes))).Replace('-', '').ToLowerInvariant()
}

function Get-B3aSortedEntry {
    # ordinal sort on .path so the output is byte-stable
    param([Parameter(Mandatory = $true)]$List)
    $cmp = [System.Comparison[object]]{
        param($a, $b)
        return [System.StringComparer]::Ordinal.Compare([string]$a.path, [string]$b.path)
    }
    $List.Sort($cmp)
    return $List
}

function Get-FileEntry {
    param([Parameter(Mandatory = $true)][string]$Root)
    $list = New-Object System.Collections.Generic.List[object]
    if (Test-Path -LiteralPath $Root) {
        foreach ($fi in @(Get-ChildItem -LiteralPath $Root -Recurse -File -Force)) {
            $rel = $fi.FullName.Substring($Root.Length).TrimStart('\', '/')
            $list.Add([ordered]@{
                path       = $rel
                size       = [long]$fi.Length
                sha256     = Get-Sha256Hex -Bytes ([System.IO.File]::ReadAllBytes($fi.FullName))
                guestClock = [ordered]@{
                    creationTime  = $fi.CreationTime.ToString('yyyy-MM-dd HH:mm:ss')
                    lastWriteTime = $fi.LastWriteTime.ToString('yyyy-MM-dd HH:mm:ss')
                }
            })
        }
    }
    return Get-B3aSortedEntry -List $list
}

function Get-ProbeState {
    # absent | empty | entries + product-layout enumeration when entries
    # (design S2.3; R/C pairing is by base file name).
    param([Parameter(Mandatory = $true)][string]$Root)
    if (-not (Test-Path -LiteralPath $Root)) {
        return [ordered]@{ state = 'absent'; files = @(); summary = $null }
    }
    $items = @(Get-ChildItem -LiteralPath $Root -Force)
    if ($items.Count -eq 0) {
        return [ordered]@{ state = 'empty'; files = @(); summary = $null }
    }
    $files = Get-FileEntry -Root $Root
    $tmpFiles = New-Object System.Collections.Generic.List[string]
    $requestIds = New-Object System.Collections.Generic.List[string]
    $completionIds = New-Object System.Collections.Generic.List[string]
    $runDirs = New-Object System.Collections.Generic.List[string]
    foreach ($e in $files) {
        if ($e.path -like '*.tmp') { $tmpFiles.Add([string]$e.path) }
    }
    foreach ($e in $files) {
        if ($e.path -like 'requests\*.jsonl') {
            $requestIds.Add([System.IO.Path]::GetFileNameWithoutExtension([string]$e.path))
        }
        if ($e.path -like 'completions\*.jsonl') {
            $completionIds.Add([System.IO.Path]::GetFileNameWithoutExtension([string]$e.path))
        }
    }
    if (Test-Path -LiteralPath (Join-Path $Root 'runs')) {
        foreach ($d in @(Get-ChildItem -LiteralPath (Join-Path $Root 'runs') -Directory -Force)) {
            $runDirs.Add([string]$d.Name)
        }
    }
    $reqSorted = @($requestIds | Sort-Object)
    $compSorted = @($completionIds | Sort-Object)
    $runSorted = @($runDirs | Sort-Object)
    $tmpSorted = @($tmpFiles | Sort-Object)
    $missing = @(); $orphans = @()
    foreach ($r in $reqSorted) { if ($compSorted -notcontains $r) { $missing += $r } }
    foreach ($c in $compSorted) { if ($reqSorted -notcontains $c) { $orphans += $c } }
    $summary = [ordered]@{
        fileCount                 = $files.Count
        tmpFiles                  = $tmpSorted
        requestCount              = $reqSorted.Count
        completionCount           = $compSorted.Count
        requestsWithoutCompletion = @($missing)
        orphanCompletions         = @($orphans)
        runDirs                   = @($runSorted)
    }
    return [ordered]@{ state = 'entries'; files = @($files); summary = $summary }
}

$probe = Get-ProbeState -Root $ProbeRoot

$dirty = 'unavailable'
$drive = ([System.IO.Path]::GetPathRoot($RoundRoot)).TrimEnd('\')
if ($drive) {
    try {
        $fsutil = fsutil dirty query $drive 2>&1
        $joined = ((@($fsutil) | ForEach-Object { [string]$_ }) -join ' ').Trim()
        if ($joined) { $dirty = $joined }
    } catch {
        $dirty = 'unavailable'
    }
}

$lastBoot = ''
try {
    $os = Get-CimInstance Win32_OperatingSystem
    $lastBoot = $os.LastBootUpTime.ToString('yyyy-MM-dd HH:mm:ss')
    $now = (Get-Date)
    $uptimeMinutes = [math]::Round((($now - $os.LastBootUpTime).TotalMinutes), 1)
} catch {
    $now = Get-Date
    $uptimeMinutes = -1
}

$inventory = [ordered]@{
    schema            = 'b3a-inventory'
    v                 = 1
    roundRoot         = $RoundRoot
    roundRootExists   = (Test-Path -LiteralPath $RoundRoot)
    files             = @(Get-FileEntry -Root $RoundRoot)
    probeRoot         = $ProbeRoot
    probeState        = $probe.state
    probeEntries      = [ordered]@{ files = @($probe.files); summary = $probe.summary }
    dirty             = $dirty
    boot              = [ordered]@{ lastBoot = $lastBoot }
    secondaryEvidence = 'guest-clock'
}

[Console]::Out.WriteLine(($inventory | ConvertTo-Json -Compress -Depth 20))
$clkLine = 'guest-clock: now={0} uptimeMinutes={1}' -f $now.ToString('yyyy-MM-dd HH:mm:ss'), $uptimeMinutes
[Console]::Error.WriteLine($clkLine)
exit 0

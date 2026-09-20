# New-B3aSession.ps1 - session preflight + guest tool deployment (design S0 / S3.1).
#
# Checks (all green before any round may start):
#   1. VBoxManage --version equals the frozen expectation (default 7.2.18r175117)
#   2. VM exists and the machine-readable facts match the frozen premise:
#      SATA/AHCI controller, useHostIOCache=false (read from the .vbox XML, the
#      machinereadable dump does not expose it), both SATA disks attached
#      (slot 0 non-rotational), NIC 1-8 all none, snapshot 'pristine' present by
#      NAME (Q2: never reference snapshot UUIDs), firmware EFI, chipset piix3,
#      2 vCPU / 4096 MB, boot1=disk
#   3. host D: free space >= 5 GB (default)
#   4. host power plan read-only assertion: STANDBYIDLE and HIBERNATEIDLE are 0
#      on AC and DC (never changes host power settings, design S6 item 5)
#   5. journal directory created + full chain verification (design S6 item 10)
#   6. guest tool deployment + guest-side hash comparison - only when the VM is
#      already running; otherwise the deployment is deferred to each round
#      (see Invoke-B3aRound.ps1 S3.5 note: restore pristine reverts C:\prfrail-prep,
#      so a session-once deployment cannot survive the first restore)
#   7. credential file + 11-shutdown-audit.ps1 host copy present (design S6 item 7)
#
# All checks are logged to <RigRoot>\<date>\session.log. Exit 0 = all green;
# exit 1 = blocked (nothing was changed).
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
param(
    [string]$VmName = 'Win11B3',
    [string]$CredFile = 'D:\VirtualBox VMs\Win11B3\host-only\cred.txt',
    [string]$GuestUser = 'prfrail',
    [string]$RigRoot = 'D:\VirtualBox VMs\Win11B3\b3a-runs',
    [string]$VBoxManage = 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe',
    [string]$ExpectedVersion = '7.2.18r175117',
    [string]$GuestToolDir = 'C:\prfrail-prep',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep',
    [long]$MinFreeBytes = 5368709120
)

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'b3a-journal-lib.ps1')

if ($MyInvocation.InvocationName -eq '.') { return }

# Fail closed: a gate that cannot run (e.g. a parameter-binding error inside a
# check) must not leave the session looking green - 2026-09-20.
$ErrorActionPreference = 'Stop'

$DateTag = Get-Date -Format 'yyyy-MM-dd'
$SessionDir = Join-Path $RigRoot $DateTag
$JournalPath = Join-Path $SessionDir 'journal.jsonl'
$SessionLog = Join-Path $SessionDir 'session.log'

function Write-B3aSessionLog {
    param([Parameter(Mandatory = $true)][string]$Text)
    Add-B3aLogLine -Path $SessionLog -Text ('[{0}] {1}' -f (Get-Date -Format 'HH:mm:ss'), $Text)
    $Text
}

$script:failed = $false
function Write-B3aCheck {
    # Detail is optional on purpose: an empty-string binding error would skip the
    # check while the surrounding session still looked green (2026-09-20).
    param([Parameter(Mandatory = $true)][string]$Name, [Parameter(Mandatory = $true)][bool]$Ok,
          [string]$Detail = 'ok')
    [void](Write-B3aSessionLog ('CHECK {0} {1} :: {2}' -f $(if ($Ok) { 'PASS' } else { 'FAIL' }), $Name, $Detail))
    '{0}  {1}  {2}' -f $(if ($Ok) { 'PASS' } else { 'FAIL' }), $Name, $Detail
    if (-not $Ok) { $script:failed = $true }
}

[void](Write-B3aSessionLog ('=== B3a session preflight {0} begin ===' -f $DateTag))

# 1. version
$ver = Invoke-B3aProcess -FileName $VBoxManage -Arguments '--version'
if ($ver.exitCode -ne 0) {
    Write-B3aCheck -Name '1-version' -Ok $false -Detail ('cannot run VBoxManage: ' + $ver.stderrText)
} else {
    $verText = $ver.stdoutText.Trim()
    Write-B3aCheck -Name '1-version' -Ok ($verText -eq $ExpectedVersion) `
        -Detail ('VBoxManage --version = {0} (expected {1})' -f $verText, $ExpectedVersion)
}

# 2. VM facts
$info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
if ($info.exitCode -ne 0) {
    Write-B3aCheck -Name '2-vm-facts' -Ok $false -Detail 'showvminfo failed'
} else {
    $v = $info.values
    $factsOk = $true
    $factsDetail = New-Object System.Collections.Generic.List[string]
    if ($v['storagecontrollername0'] -ne 'SATA' -or $v['storagecontrollertype0'] -ne 'IntelAhci') {
        $factsOk = $false; $factsDetail.Add('storage controller is not SATA/AHCI')
    }
    $vboxCfg = ''
    if ($v.Contains('CfgFile')) { $vboxCfg = $v['CfgFile'] }
    if (-not $vboxCfg -or -not (Test-Path -LiteralPath $vboxCfg)) {
        $factsOk = $false; $factsDetail.Add('cannot locate the .vbox file')
    } else {
        $xml = [System.IO.File]::ReadAllText($vboxCfg)
        if ($xml -notmatch 'useHostIOCache="false"') {
            $factsOk = $false; $factsDetail.Add('useHostIOCache is not false')
        }
    }
    foreach ($slot in @('SATA-0-0', 'SATA-1-0')) {
        $attached = $v[$slot]
        if (-not $attached -or $attached -eq 'none') { $factsOk = $false; $factsDetail.Add($slot + ' not attached') }
    }
    foreach ($port in @('0', '1')) {
        # the frozen premise makes both disks non-rotational (single host NVMe); the
        # injection volume (port 1) is the one whose media semantics matter most
        if ($v['SATA-nonrotational-' + $port + '-0'] -ne 'on') {
            $factsOk = $false; $factsDetail.Add(('SATA-{0}-0 lost its non-rotational (SSD) flag' -f $port))
        }
    }
    for ($n = 1; $n -le 8; $n++) {
        $nic = ''
        if ($v.Contains('nic' + $n)) { $nic = $v['nic' + $n] }
        if ($nic -and $nic -ne 'none' -and $nic -ne 'null') {
            $factsOk = $false
            $factsDetail.Add(('NIC {0} = {1} (frozen premise: none; see design Q1)' -f $n, $nic))
        }
    }
    if ($v['firmware'] -ne 'EFI') { $factsOk = $false; $factsDetail.Add('firmware != EFI') }
    if ($v['chipset'] -ne 'piix3') { $factsOk = $false; $factsDetail.Add('chipset != piix3') }
    if ($v['cpus'] -ne '2') { $factsOk = $false; $factsDetail.Add('cpus != 2') }
    if ($v['memory'] -ne '4096') { $factsOk = $false; $factsDetail.Add('memory != 4096 MB') }
    if ($v['boot1'] -ne 'disk') { $factsOk = $false; $factsDetail.Add('boot1 != disk') }
    $snaps = Get-B3aSnapshotList -VBoxManage $VBoxManage -VmName $VmName
    $hasPristine = (@($snaps.names) -contains 'pristine')
    if (-not $hasPristine) { $factsOk = $false; $factsDetail.Add("snapshot 'pristine' not found (referenced by NAME, Q2)") }
    # The detail is the premise attestation recorded in the session log, so it names
    # the facts the checks actually verified.
    $factsSummary = ('VMState={0}; SATA nonrotational={1}/{2}; pristine snapshot present={3}' -f `
        $v['VMState'], $v['SATA-nonrotational-0-0'], $v['SATA-nonrotational-1-0'], $hasPristine)
    if ($factsDetail.Count -gt 0) { $factsSummary += '; ' + ($factsDetail -join '; ') }
    Write-B3aCheck -Name '2-vm-facts' -Ok $factsOk -Detail $factsSummary
}

# 3. free space
try {
    $driveInfo = New-Object System.IO.DriveInfo('D:')
    Write-B3aCheck -Name '3-free-space' -Ok ($driveInfo.AvailableFreeSpace -ge $MinFreeBytes) `
        -Detail ('D: free = {0:N1} GB (min {1:N1})' -f ($driveInfo.AvailableFreeSpace / 1GB), ($MinFreeBytes / 1GB))
} catch {
    Write-B3aCheck -Name '3-free-space' -Ok $false -Detail ('D: unavailable: ' + $_.Exception.Message)
}

# 4. host power plan (read-only assertion)
$powerOk = $true
$powerDetail = New-Object System.Collections.Generic.List[string]
$scheme = Invoke-B3aProcess -FileName 'powercfg' -Arguments '/getactivescheme'
if ($scheme.exitCode -ne 0) {
    $powerOk = $false; $powerDetail.Add('powercfg unavailable')
} else {
    foreach ($setting in @('SUB_SLEEP STANDBYIDLE', 'SUB_SLEEP HIBERNATEIDLE')) {
        $q = Invoke-B3aProcess -FileName 'powercfg' -Arguments ('/query SCHEME_CURRENT ' + $setting)
        if ($q.exitCode -ne 0) {
            $powerOk = $false; $powerDetail.Add($setting + ' query failed')
            continue
        }
        # The labels are localized (this host runs zh-CN), so anchor on structure:
        # value lines of the setting block end with 0x<hex>, and the last two are
        # the AC and DC indices (the earlier ones are min/max/increment). The raw
        # output is logged so the parse stays auditable.
        [void](Write-B3aSessionLog ('powercfg ' + $setting + ' raw: ' + ($q.stdoutText -replace "`r?`n", ' | ')))
        $vals = @([regex]::Matches($q.stdoutText, ':\s+0x([0-9a-fA-F]{1,8})\s*$', 'Multiline') | ForEach-Object { $_.Groups[1].Value })
        if ($vals.Count -lt 2) {
            $powerOk = $false; $powerDetail.Add($setting + ' indices not parseable (value lines found: ' + $vals.Count + ')')
            continue
        }
        $ac = $vals[$vals.Count - 2].PadLeft(8, '0')
        $dc = $vals[$vals.Count - 1].PadLeft(8, '0')
        if ($ac -ne '00000000' -or $dc -ne '00000000') {
            $powerOk = $false; $powerDetail.Add(($setting + ' not 0 (AC=' + $ac + ' DC=' + $dc + ')'))
        }
    }
}
$powerSummary = 'sleep and hibernate indices are 0 for AC and DC'
if ($powerDetail.Count -gt 0) { $powerSummary = ($powerDetail -join '; ') }
Write-B3aCheck -Name '4-host-power-plan' -Ok $powerOk -Detail $powerSummary

# 5. journal directory + full chain verification
[void](New-Item -ItemType Directory -Path $SessionDir -Force)
$chain = Test-B3aJournalChain -JournalPath $JournalPath
Write-B3aCheck -Name '5-journal-chain' -Ok $chain.ok `
    -Detail ('journal={0} records={1} badLine={2}' -f $JournalPath, $chain.count, $chain.firstBadLine)

# 6. guest tool deployment (only possible while the VM is running)
$state = ''
if ($info.values) { $state = $info.values['VMState'] }
if ($state -eq 'running') {
    $deploy = Invoke-B3aToolDeploy -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
        -GuestToolDir $GuestToolDir -GuestPrepDir $GuestPrepDir -ScriptsDir $PSScriptRoot
    $detail = (@($deploy.files) | ForEach-Object { $_.name + ':' + $(if ($_.hashMatch) { 'ok' } else { 'BAD' }) }) -join ' '
    Write-B3aCheck -Name '6-tool-deploy' -Ok $deploy.ok -Detail $detail
} else {
    Write-B3aCheck -Name '6-tool-deploy' -Ok $true `
        -Detail ('VM is {0}; deployment deferred to every round after restore pristine (see Invoke-B3aRound.ps1 S3.5)' -f $state)
}

# 7. credential file + host copy of the audit script
$credOk = (Test-Path -LiteralPath $CredFile)
$auditHost = Join-Path $GuestPrepDir '11-shutdown-audit.ps1'
$auditOk = (Test-Path -LiteralPath $auditHost)
Write-B3aCheck -Name '7-inputs' -Ok ($credOk -and $auditOk) `
    -Detail ('credFile={0} auditScript={1}' -f $(if ($credOk) { 'present' } else { 'MISSING' }), $(if ($auditOk) { 'present' } else { 'MISSING' }))

[void](Write-B3aSessionLog ('=== preflight done: ' + $(if ($script:failed) { 'BLOCKED' } else { 'all green' }) + ' ==='))
if ($script:failed) { exit 1 }
exit 0

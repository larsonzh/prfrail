# New-B3bSession.ps1 - B3b session preflight (design section 5, step S0).
#
# S0 reuses the frozen B3a preflight instead of reimplementing it: the environment
# assertions (VBoxManage version, VM facts, free space, host power plan, journal
# chain, tool deployment, credential and audit inputs) are identical for both
# slices, and a second implementation would be a second source of truth for the
# premise. This script runs that preflight as a child process, then adds the three
# checks that are specific to B3b:
#
#   8.  reuse pin: the bytes of the B3a files B3b depends on must equal the hashes
#       recorded when the session directory was created (fail closed on drift)
#   9.  module hygiene: no B3b module may shadow a B3a helper, and no B3b module
#       may carry a top-level param block (B3a defect D1)
#   10. writer binary: the guest writer must exist and its hash is recorded, because
#       the verdict is about the bytes that actually ran
#
# Exit 0 = all green; exit 1 = blocked (nothing was changed).
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
param(
    [string]$VmName = 'Win11B3',
    [string]$CredFile = 'D:\VirtualBox VMs\Win11B3\host-only\cred.txt',
    [string]$GuestUser = 'prfrail',
    [string]$RigRoot = 'D:\VirtualBox VMs\Win11B3\b3b-runs',
    [string]$VBoxManage = 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe',
    [string]$ExpectedVersion = '7.2.18r175117',
    [string]$GuestToolDir = 'C:\prfrail-prep',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep',
    [string]$B3aRigDir = '',
    [string]$WriterExe = '',
    # The session id is explicit so a run that crosses midnight stays in one session. An
    # existing pin.json wins: re-running the preflight inside a started session must not
    # move it.
    [string]$SessionId = '',
    [string]$StoreRoot = 'R:\agent-runner-replay',
    [string]$ScratchRoot = 'R:\calib',
    [switch]$RecordPin
)

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')

if ($MyInvocation.InvocationName -eq '.') { return }

# Fail closed: a gate that cannot run must not leave the session looking green.
$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrEmpty($B3aRigDir)) {
    $B3aRigDir = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\b3a-rig'))
}
if ([string]::IsNullOrEmpty($WriterExe)) {
    $repoRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\..'))
    $WriterExe = Join-Path $repoRoot 'tmp\b3b\writer.exe'
}

# Whether the caller asked for a specific session matters: an explicit id that disagrees with
# the pin is a real conflict (refused), while an auto-resolved date simply yields to the pin of
# the session that is already here.
$explicitSessionId = -not [string]::IsNullOrEmpty($SessionId)
if (-not $explicitSessionId) { $SessionId = Get-Date -Format 'yyyy-MM-dd' }
$DateTag = $SessionId
$SessionDir = Join-Path $RigRoot $DateTag
[void](New-Item -ItemType Directory -Path $SessionDir -Force)
$SessionLog = Join-Path $SessionDir 'session.log'
$PinPath = Join-Path $SessionDir 'pin.json'
if (Test-Path -LiteralPath $PinPath) {
    $existingPin = (Get-B3aFileTextUtf8 -Path $PinPath) | ConvertFrom-Json
    $pinnedId = [string]$existingPin.sessionId
    if ((-not [string]::IsNullOrEmpty($pinnedId)) -and ($pinnedId -ne $DateTag)) {
        # An explicitly requested session id that disagrees with the pin is a real conflict:
        # silently adopting either side would leave the preflight and the matrix pointed at
        # different directories. It is refused, and the operator has to say which session is
        # meant.
        if ($explicitSessionId) {
            Write-Output ('REFUSED: the session directory already pins session {0}, but -SessionId requested {1}' -f $pinnedId, $SessionId)
            exit 1
        }
        $DateTag = $pinnedId
        $SessionDir = Join-Path $RigRoot $DateTag
        $SessionLog = Join-Path $SessionDir 'session.log'
        $PinPath = Join-Path $SessionDir 'pin.json'
        Write-Output ('session id taken from the existing pin: ' + $DateTag)
    }
}
# One rig process per session: the preflight probes and may start the same VM the matrix drives,
# so it takes the same lock. It releases it on the way out because the preflight is short-lived
# and must not block the run that follows it.
$preflightLock = Enter-B3bSessionLock -SessionDir $SessionDir -Owner 'New-B3bSession'
if (-not $preflightLock.ok) {
    Write-Output ('REFUSED: ' + $preflightLock.reason)
    Write-Output 'Never probe or drive the experiment VM while another rig process holds the session.'
    exit 1
}

function Write-B3bSessionLog {
    param([Parameter(Mandatory = $true)][string]$Text)
    Add-B3aLogLine -Path $SessionLog -Text ('[{0}] {1}' -f (Get-Date -Format 'HH:mm:ss'), $Text)
    $Text
}

$script:failed = $false
function Write-B3bCheck {
    param([Parameter(Mandatory = $true)][string]$Name, [Parameter(Mandatory = $true)][bool]$Ok, [string]$Detail = 'ok')
    [void](Write-B3bSessionLog ('CHECK {0} {1} :: {2}' -f $(if ($Ok) { 'PASS' } else { 'FAIL' }), $Name, $Detail))
    '{0}  {1}  {2}' -f $(if ($Ok) { 'PASS' } else { 'FAIL' }), $Name, $Detail
    if (-not $Ok) { $script:failed = $true }
}

# A terminating error anywhere below must not leave a live-pid lock behind: the next honest run
# would then be refused until someone deleted the file by hand. A trap covers every terminating
# path (including one thrown from inside a helper) with a single handler, which is why it is used
# instead of re-indenting the whole preflight into a try/finally.
trap {
    if (-not [string]::IsNullOrEmpty($SessionDir)) { [void](Exit-B3bSessionLock -SessionDir $SessionDir) }
    Write-Output ('preflight aborted: ' + $_.Exception.Message)
    exit 1
}

[void](Write-B3bSessionLog ('=== B3b session preflight {0} begin ===' -f $DateTag))

# 1-7. frozen B3a preflight, run as a child process so its guard against dot-sourcing
# (and its parameter binding) stays intact.
$frozenSession = Join-Path $B3aRigDir 'New-B3aSession.ps1'
if (-not (Test-Path -LiteralPath $frozenSession)) {
    Write-B3bCheck -Name '1-frozen-preflight' -Ok $false -Detail ('missing: ' + $frozenSession)
} else {
    $arguments = '-NoProfile -ExecutionPolicy Bypass -File {0} -VmName {1} -CredFile {2} -GuestUser {3} -RigRoot {4} -VBoxManage {5} -ExpectedVersion {6} -GuestToolDir {7} -GuestPrepDir {8}' -f `
        (Get-B3aQuoted $frozenSession), (Get-B3aQuoted $VmName), (Get-B3aQuoted $CredFile), (Get-B3aQuoted $GuestUser), `
        (Get-B3aQuoted $RigRoot), (Get-B3aQuoted $VBoxManage), $ExpectedVersion, (Get-B3aQuoted $GuestToolDir), (Get-B3aQuoted $GuestPrepDir)
    $run = Invoke-B3aProcess -FileName 'powershell' -Arguments $arguments -TimeoutS 600
    foreach ($lineRaw in ($run.stdoutText -split "`n")) {
        $line = $lineRaw.Trim()
        if ($line) { [void](Write-B3bSessionLog ('  frozen: ' + $line)) }
    }
    if ($run.stderrText.Trim()) { [void](Write-B3bSessionLog ('  frozen stderr: ' + (Get-B3aBoundedText -Text $run.stderrText))) }
    Write-B3bCheck -Name '1-frozen-preflight' -Ok ($run.exitCode -eq 0) `
        -Detail ('New-B3aSession.ps1 exit={0} (1-7 inherited from the frozen preflight)' -f $run.exitCode)
}

# 8. reuse pin
$currentHashes = Get-B3bReusePinMap -B3aRigDir $B3aRigDir
if ($RecordPin -or -not (Test-Path -LiteralPath $PinPath)) {
    $pinDoc = [ordered]@{
        schema    = 'b3b-pin'
        v         = 1
        sessionId = $DateTag
        recordedAt = (Get-B3aHostTimestamp)
        b3aRigDir = $B3aRigDir
        hashes    = $currentHashes
    }
    Write-B3aTextFile -Path $PinPath -Text ($pinDoc | ConvertTo-Json -Depth 10) -WithBom
    Write-B3bCheck -Name '8-reuse-pin' -Ok $true -Detail ('recorded pin: ' + $PinPath)
} else {
    $pinDoc = (Get-B3aFileTextUtf8 -Path $PinPath) | ConvertFrom-Json
    $expected = @{}
    foreach ($property in $pinDoc.hashes.PSObject.Properties) { $expected[$property.Name] = [string]$property.Value }
    $pinResult = Test-B3bPinIntegrity -B3aRigDir $B3aRigDir -Expected $expected
    $detail = 'hashes match the session pin'
    if (-not $pinResult.ok) {
        $detail = ('drift: ' + (($pinResult.drift) -join '; ') + ' missing: ' + (($pinResult.missingPins) -join ','))
    }
    Write-B3bCheck -Name '8-reuse-pin' -Ok $pinResult.ok -Detail $detail
}

# 9. Guest store-root writability. The round driver needs these two paths on the injection
# volume, and a permission or mounting problem there would only surface around S4/S6, once
# the round has already cost a cut and a restart. Probing it during the preflight turns that
# into a free failure.
#
# The probe uses a unique name and sets ErrorActionPreference inside the guest command, and
# it deliberately does not end with a hard-coded exit 0: a probe that reports success no
# matter what would be worse than no probe at all.
$probeToken = [guid]::NewGuid().ToString('N').Substring(0, 8)
$probeCommand = ('$ErrorActionPreference = ''Stop''; ' +
    'New-Item -ItemType Directory -Force -Path ''{0}'',''{1}'' | Out-Null; ' +
    '$probe = ''b3b-preflight-{2}.txt''; ' +
    'Set-Content -LiteralPath (Join-Path ''{0}'' $probe) -Value probe; ' +
    'Copy-Item -LiteralPath (Join-Path ''{0}'' $probe) -Destination (Join-Path ''{1}'' $probe) -Force; ' +
    'Remove-Item -LiteralPath (Join-Path ''{0}'' $probe),(Join-Path ''{1}'' $probe) -Force') -f $StoreRoot, $ScratchRoot, $probeToken
$vmUp = Start-B3bVmIfOff -VBoxManage $VBoxManage -VmName $VmName
$guestReady = $false
if ($vmUp) {
    $guestReady = Wait-B3bGuestCommandReady -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile `
        -GuestUser $GuestUser -TimeoutS 300
}
if (-not $vmUp) {
    Write-B3bCheck -Name '9-guest-store-root-writable' -Ok $false -Detail ('the VM could not be started headless: ' + $VmName)
} elseif (-not $guestReady) {
    Write-B3bCheck -Name '9-guest-store-root-writable' -Ok $false -Detail 'the guest did not become reachable within 300 s'
} else {
    $probeRun = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
        -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
        -GuestArgs @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-Command', $probeCommand) -TimeoutMs 90000
    $probeDetail = '{0} store={1} scratch={2} exit={3}' -f $probeToken, $StoreRoot, $ScratchRoot, $probeRun.exitCode
    if ($probeRun.exitCode -ne 0) {
        $probeDetail = $probeDetail + ' :: ' + (Get-B3aBoundedText -Text ($probeRun.stdoutText + ' ' + $probeRun.stderrText))
    }
    Write-B3bCheck -Name '9-guest-store-root-writable' -Ok ($probeRun.exitCode -eq 0) -Detail $probeDetail
}

# 9. module hygiene
$hygiene = Test-B3bModuleHygiene -B3bRigDir $PSScriptRoot -B3aRigDir $B3aRigDir
$hygieneDetail = ('b3b functions={0}; collisions={1}; unsafe modules={2}' -f `
    (@($hygiene.b3bFunctions).Count), (($hygiene.collisions) -join ','), (($hygiene.topLevelParams) -join ','))
Write-B3bCheck -Name '9-module-hygiene' -Ok $hygiene.ok -Detail $hygieneDetail

# 10. writer binary
if (Test-Path -LiteralPath $WriterExe) {
    $writerSha = Get-B3aSha256Hex -Path $WriterExe
    Write-B3bCheck -Name '10-writer-binary' -Ok $true -Detail ('{0} sha256={1} bytes={2}' -f $WriterExe, $writerSha, (Get-Item -LiteralPath $WriterExe).Length)
} else {
    Write-B3bCheck -Name '10-writer-binary' -Ok $false `
        -Detail ('missing: ' + $WriterExe + ' (build it with: go test -c -tags b3bnative -o <path> .\internal\adapters\)')
}

[void](Write-B3bSessionLog ('=== preflight done: ' + $(if ($script:failed) { 'BLOCKED' } else { 'all green' }) + ' ==='))
[void](Exit-B3bSessionLock -SessionDir $SessionDir)
if ($script:failed) { exit 1 }
exit 0

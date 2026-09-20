# Invoke-B3aRound.ps1 - one B3a round, steps S1-S9 (design S3.2).
#
# Round protocol (gates in parentheses; a failed gate stops the round):
#   S1  assert VM powered off (hard poweroff first if needed, logged) + restore
#       the 'pristine' snapshot by NAME (restore exit 0)
#   S2  startvm --type headless (exit 0)
#   S3  poll guestcontrol until ready (<= MaxReadyWaitS)
#   S3.5 deploy the three guest scripts + guest-side hash compare (design S0 item 6
#       is repeated here on purpose: restore pristine reverts C:\prfrail-prep, so
#       a once-per-session deployment cannot survive the first restore - a
#       deliberate deviation to be folded back into the design)
#   S4  pre-injection gate: 11-shutdown-audit.ps1, uptime < 50 min and
#       VERDICT != LICENCE-SHUTDOWN (else void)
#   S4.5 round-root setup: cmd /c mkdir R:\calib\rNN (rig-side housekeeping, not
#       part of the measured write; restore pristine erases R: every round)
#   S5  host computes digest/size from the seed, writes the 'plan' journal record,
#       flush + read-back ACK (journal first: no ACK means no write and no cut)
#   S6  control write 13-round-write.ps1 (probe JSON on stdout only; guest
#       writes exactly one control file this round); the probe digest/size must
#       match the plan pre-commitment bi-directionally (else void, no cut)
#   S7  operator confirm (skip only with -NoConfirm, a session-level choice) ->
#       sleep PlanCutDelayS (2.0) -> controlvm poweroff -> 'cut' record + ACK
#   S8  restart, shutdown audit (VERDICT == HARD-POWER-LOSS else void),
#       inventory x2 (byte-identical stdout else void)
#   S9  Get-B3aReconcile -> reconcile.json + 'reconcile' journal record + ACK
#   post-round: poweroff + restore pristine (design S5.4: never leave the guest
#       running between steps)
#
# Void rules are design S3.3 items 1-7. A void round restores pristine and the
# calibration re-runs the same logical slot with a NEW round number.
#
# Exit codes: 0 valid round finished (outcome in rounds/rNN/reconcile.json)
#             3 void round (pristine restored; re-run with a NEW round number)
#             4 valid parity round but outcome=survived (operator decision, S4.4)
#             5 operational fail (re-run the SAME round number; budget 2)
#             6 session blocker (journal/config/ACK - stop everything)
#             7 operator declined the cut
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
param(
    [Parameter(Mandatory = $true)][ValidateRange(0, 999)][int]$Round,
    [Parameter(Mandatory = $true)][ValidateSet('parity', 'negative-control', 'positive-control')][string]$Candidate,
    [string]$VmName = 'Win11B3',
    [string]$CredFile = 'D:\VirtualBox VMs\Win11B3\host-only\cred.txt',
    [string]$GuestUser = 'prfrail',
    [string]$RigRoot = 'D:\VirtualBox VMs\Win11B3\b3a-runs',
    [string]$VBoxManage = 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe',
    [string]$GuestToolDir = 'C:\prfrail-prep',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep',
    [string]$ProbeRoot = 'R:\agent-runner-replay',
    [int]$Size = 262144,
    [double]$PlanCutDelayS = 2.0,
    [int]$MaxReadyWaitS = 120,
    [switch]$NoConfirm
)

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'b3a-journal-lib.ps1')
. (Join-Path $PSScriptRoot 'Get-B3aReconcile.ps1')

if ($MyInvocation.InvocationName -eq '.') { return }

# Fail closed: an unexpected error must abort the round loudly instead of skipping
# a gate silently.
$ErrorActionPreference = 'Stop'

# round derivation (design S5.2)
$DateTag = Get-Date -Format 'yyyy-MM-dd'
$SessionDir = Join-Path $RigRoot $DateTag
$JournalPath = Join-Path $SessionDir 'journal.jsonl'
$SessionLog = Join-Path $SessionDir 'session.log'
$RoundLabel = 'r{0:D2}' -f $Round
$RoundDir = Join-Path $SessionDir (Join-Path 'rounds' $RoundLabel)
$RoundRootGuest = 'R:\calib\' + $RoundLabel

$stage = 'parity'
$mode = 'negative'
$fileName = 'parity.bin'
$seed = 'b3a-parity-' + $Round
$literal = $null
$sizeEff = 14
if ($Candidate -eq 'negative-control') {
    $stage = 'temp-written'
    $fileName = 'neg.bin'
    $seed = 'b3a-neg-' + $Round
    $sizeEff = $Size
}
if ($Candidate -eq 'positive-control') {
    $stage = 'parent-synced'
    $mode = 'positive'
    $fileName = 'pos.bin'
    $seed = 'b3a-pos-' + $Round
    $sizeEff = $Size
}
if ($Candidate -eq 'parity') { $literal = 'pre-cut-marker' }
$targetPath = $RoundRootGuest + '\' + $fileName

function Write-B3aRoundLog {
    param([Parameter(Mandatory = $true)][string]$Text)
    Add-B3aLogLine -Path $SessionLog -Text ('[{0}] {1}' -f (Get-Date -Format 'HH:mm:ss'), $Text)
    $Text
}

function Invoke-B3aBaselineReset {
    # post-round / void recovery: poweroff (if running) + restore pristine (S5.4)
    $info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
    $state = ''
    if ($info.values) { $state = $info.values['VMState'] }
    if ($state -eq 'running' -or $state -eq 'paused' -or $state -eq 'saved') {
        $cut = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('controlvm "{0}" poweroff' -f $VmName)
        [void](Write-B3aRoundLog ('reset: controlvm poweroff (was {0}) exit={1}' -f $state, $cut.exitCode))
    }
    $restore = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('snapshot "{0}" restore pristine' -f $VmName)
    [void](Write-B3aRoundLog ('reset: restore pristine exit=' + $restore.exitCode))
    return ($restore.exitCode -eq 0)
}

function Invoke-B3aVoid {
    param([Parameter(Mandatory = $true)][string]$Reason)
    [void](Write-B3aRoundLog ('VOID: ' + $Reason))
    $rec = Invoke-B3aReconcile -OutputPath (Join-Path $RoundDir 'reconcile.json') `
        -JournalPath $JournalPath -Round $Round -Candidate $Candidate -Stage $stage `
        -TargetPath $targetPath -ContentDigest $contentDigestKnown -ContentSize $contentSizeKnown `
        -VoidReason $Reason
    if (-not $rec.ok) {
        [void](Write-B3aRoundLog ('BLOCKER: void reconcile failed: ' + $rec.error))
        exit 6
    }
    if (-not (Invoke-B3aBaselineReset)) { exit 6 }
    exit 3
}

function Wait-B3aReady {
    # The timeout is passed in explicitly: the ready-wait must not close over
    # script-scope state, and S3's budget has to stay visible at the call site.
    param([Parameter(Mandatory = $true)][int]$TimeoutS)
    $elapsed = New-Object System.Diagnostics.Stopwatch
    $elapsed.Start()
    while ($elapsed.Elapsed.TotalSeconds -lt $TimeoutS) {
        $ver = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
            -Exe 'C:\Windows\System32\cmd.exe' -GuestArgs @('/c', 'ver') -TimeoutMs 15000
        if ($ver.exitCode -eq 0) { $elapsed.Stop(); return $elapsed.Elapsed.TotalSeconds }
    }
    $elapsed.Stop()
    return -1
}

# --- round start ------------------------------------------------------------
[void](New-Item -ItemType Directory -Path $RoundDir -Force)
[void](Write-B3aRoundLog ('=== round {0} candidate={1} stage={2} begin ===' -f $Round, $Candidate, $stage))

$contentDigestKnown = ''
$contentSizeKnown = 0
if ($literal) {
    $contentDigestKnown = Get-B3aSeedPayloadDigest -Size $sizeEff -Literal $literal
    $contentSizeKnown = $sizeEff
} else {
    $contentDigestKnown = Get-B3aSeedPayloadDigest -Seed $seed -Size $sizeEff
    $contentSizeKnown = $sizeEff
}

# preflight (cheap local blockers)
if (-not (Test-Path -LiteralPath $VBoxManage)) { [void](Write-B3aRoundLog 'BLOCKER: VBoxManage missing'); exit 6 }
if (-not (Test-Path -LiteralPath $CredFile)) { [void](Write-B3aRoundLog 'BLOCKER: credential file missing'); exit 6 }
$chain = Test-B3aJournalChain -JournalPath $JournalPath
if (-not $chain.ok) {
    [void](Write-B3aRoundLog ('BLOCKER: journal chain broken at line {0}: {1}' -f $chain.firstBadLine, $chain.firstError))
    exit 6
}

# S1 restore baseline
$info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
if ($info.exitCode -ne 0) { [void](Write-B3aRoundLog 'BLOCKER: cannot query VM'); exit 6 }
$state = ''
if ($info.values) { $state = $info.values['VMState'] }
if ($state -ne 'poweroff') {
    $pre = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('controlvm "{0}" poweroff' -f $VmName)
    [void](Write-B3aRoundLog ('S1: pre-round poweroff (was {0}) exit={1}' -f $state, $pre.exitCode))
}
$restore = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('snapshot "{0}" restore pristine' -f $VmName)
[void](Write-B3aRoundLog ('S1: restore pristine exit=' + $restore.exitCode))
if ($restore.exitCode -ne 0) { [void](Write-B3aRoundLog 'BLOCKER: restore pristine failed'); exit 6 }

# S2 start
$start = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('startvm "{0}" --type headless' -f $VmName)
[void](Write-B3aRoundLog ('S2: startvm exit=' + $start.exitCode))
if ($start.exitCode -ne 0) { exit 5 }

# S3 ready wait
$readyS = Wait-B3aReady -TimeoutS $MaxReadyWaitS
if ($readyS -lt 0) { [void](Write-B3aRoundLog 'S3: ready wait timed out'); exit 5 }
[void](Write-B3aRoundLog ('S3: guestcontrol ready after {0:N1}s' -f $readyS))

# S3.5 tool deployment (per-round, see header note)
$deploy = Invoke-B3aToolDeploy -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -GuestToolDir $GuestToolDir -GuestPrepDir $GuestPrepDir -ScriptsDir $PSScriptRoot
foreach ($f in $deploy.files) {
    [void](Write-B3aRoundLog ('S3.5: deploy {0} copied={1} hashMatch={2} note={3}' -f $f.name, $f.copied, $f.hashMatch, $f.note))
}
if (-not $deploy.ok) { exit 5 }

# S4 pre-injection gate
$auditPre = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
    -GuestArgs @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\11-shutdown-audit.ps1')) `
    -TimeoutMs 60000
if ($auditPre.exitCode -ne 0) { [void](Write-B3aRoundLog 'S4: pre-injection audit failed'); exit 5 }
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'audit-pre.txt') -Bytes $auditPre.stdoutBytes
$preVerdict = Read-B3aAuditVerdict -Text $auditPre.stdoutText
[void](Write-B3aRoundLog ('S4: pre-audit verdict={0} uptime={1}' -f $preVerdict.verdict, $preVerdict.uptimeMinutes))
if ($preVerdict.verdict -eq 'LICENCE-SHUTDOWN') { Invoke-B3aVoid 'licence-shutdown' }
if ($null -eq $preVerdict.uptimeMinutes) { [void](Write-B3aRoundLog 'S4: uptime unparseable'); exit 5 }
if ($preVerdict.uptimeMinutes -ge 50) { Invoke-B3aVoid 'uptime-too-old' }

# S4.5 round-root setup (rig-side housekeeping, NOT part of the measured write).
# Every round runs from the pristine snapshot, so R:\calib\rNN never exists yet;
# without it the measured write dies with DirectoryNotFoundException (VBoxManage
# reported guest exit 33). It is created here, before the plan record, so the
# measured guest write stays exactly one file write (design S3.2 invariant).
$mkdir = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\cmd.exe' -GuestArgs @('/c', 'mkdir', (Get-B3aQuoted $RoundRootGuest)) -TimeoutMs 30000
[void](Write-B3aRoundLog ('S4.5: round root mkdir exit=' + $mkdir.exitCode))
if ($mkdir.exitCode -ne 0) { exit 5 }

# S5 plan record: committed and acknowledged BEFORE the guest writes anything
# (design S2.1 "journal first"). The record is a pre-commitment - digest and size
# are computed here on the host from the seed/literal, never copied from the probe.
$planRec = New-B3aJournalRecord -Round $Round -Candidate $Candidate -Stage $stage -TargetPath $targetPath `
    -ContentDigest $contentDigestKnown -ContentSize ([long]$contentSizeKnown) -PlanCutDelayS $PlanCutDelayS -Op 'plan'
$planAck = Add-B3aJournalRecord -JournalPath $JournalPath -Record $planRec
if (-not $planAck.ack) { Invoke-B3aVoid ('plan-no-ack: ' + $planAck.error) }
Write-B3aTextFile -Path (Join-Path $RoundDir 'plan.json') -Text $planAck.recordJson
[void](Write-B3aRoundLog 'S5: plan record ACK ok (committed before the guest write)')

# S6 control write (the guest writes exactly one control file this round)
$writeArgs = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\13-round-write.ps1'), '-Mode', $mode, '-Path', $targetPath, '-Size', ([string]$sizeEff))
if ($literal) { $writeArgs += @('-Literal', $literal) } else { $writeArgs += @('-Seed', $seed) }
$write = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' -GuestArgs $writeArgs -TimeoutMs 60000
if ($write.exitCode -ne 0) {
    # archive the probe output and stderr before leaving: the probe still prints its
    # JSON on failure, and without it a failed write is undiagnosable (2026-09-20)
    Write-B3aRawTextFile -Path (Join-Path $RoundDir 'write-probe.txt') -Bytes $write.stdoutBytes
    Write-B3aRawTextFile -Path (Join-Path $RoundDir 'write-stderr.txt') -Bytes $write.stderrBytes
    [void](Write-B3aRoundLog ('S6: control write exit=' + $write.exitCode + '; probe and stderr archived in the round directory'))
    exit 5
}
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'write-probe.txt') -Bytes $write.stdoutBytes
$probe = Get-B3aJsonFromText -Text $write.stdoutText
if ($null -eq $probe) { [void](Write-B3aRoundLog 'S6: probe JSON unparseable'); exit 5 }
if ([string]$probe.writeResult -ne 'ok') {
    [void](Write-B3aRoundLog ('S6: probe writeResult=' + $probe.writeResult + ' error=' + $probe.error))
    exit 5
}
[void](Write-B3aRoundLog ('S6: probe digest={0} size={1} apis={2}' -f $probe.contentDigest, $probe.size, (@($probe.apisCalled) -join ',')))
# bi-directional check against the S5 pre-commitment (mismatch => void, no cut)
if (([string]$probe.contentDigest) -ne $contentDigestKnown) {
    Invoke-B3aVoid ('digest-mismatch host={0} probe={1}' -f $contentDigestKnown, $probe.contentDigest)
}
if ([long]$probe.size -ne [long]$contentSizeKnown) {
    Invoke-B3aVoid ('size-mismatch host={0} probe={1}' -f $contentSizeKnown, $probe.size)
}

# S7 operator confirmation + hard cut
if (-not $NoConfirm) {
    '--- plan record (journal ACK ok) ---'
    $planAck.recordJson
    ''
    $answer = Read-Host ('Hard poweroff of {0} in {1} s? Type YES to cut, anything else aborts' -f $VmName, $PlanCutDelayS)
    if ($answer -ne 'YES') {
        [void](Write-B3aRoundLog 'S7: operator declined the cut; VM left running')
        'aborted by operator; VM left running (run Invoke-B3aRound again to continue)'
        exit 7
    }
}
Start-Sleep -Seconds $PlanCutDelayS
$cut = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('controlvm "{0}" poweroff' -f $VmName)
[void](Write-B3aRoundLog ('S7: controlvm poweroff exit=' + $cut.exitCode))
if ($cut.exitCode -ne 0) { Invoke-B3aVoid ('poweroff-nonzero exit=' + $cut.exitCode) }
$cutRec = New-B3aJournalRecord -Round $Round -Candidate $Candidate -Stage $stage -TargetPath $targetPath `
    -ContentDigest $contentDigestKnown -ContentSize ([long]$contentSizeKnown) -PlanCutDelayS $PlanCutDelayS -Op 'cut' `
    -Extra ([ordered]@{ cutExitCode = $cut.exitCode; cutHostTs = (Get-B3aHostTimestamp) })
$cutAck = Add-B3aJournalRecord -JournalPath $JournalPath -Record $cutRec
if (-not $cutAck.ack) { [void](Write-B3aRoundLog ('BLOCKER: cut record ACK failed: ' + $cutAck.error)); exit 6 }
Write-B3aTextFile -Path (Join-Path $RoundDir 'cut.json') -Text $cutAck.recordJson
[void](Write-B3aRoundLog 'S7: cut record ACK ok')

# S8 restart + audit + inventory x2
$start2 = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('startvm "{0}" --type headless' -f $VmName)
if ($start2.exitCode -ne 0) { [void](Write-B3aRoundLog ('S8: restart exit=' + $start2.exitCode)); exit 5 }
$readyS2 = Wait-B3aReady -TimeoutS $MaxReadyWaitS
if ($readyS2 -lt 0) { [void](Write-B3aRoundLog 'S8: ready wait after restart timed out'); exit 5 }
$audit = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
    -GuestArgs @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\11-shutdown-audit.ps1')) `
    -TimeoutMs 60000
if ($audit.exitCode -ne 0) { [void](Write-B3aRoundLog 'S8: post-cut audit failed'); exit 5 }
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'audit.txt') -Bytes $audit.stdoutBytes
$postVerdict = Read-B3aAuditVerdict -Text $audit.stdoutText
[void](Write-B3aRoundLog ('S8: post-audit verdict={0}' -f $postVerdict.verdict))
if ($postVerdict.verdict -ne 'HARD-POWER-LOSS') {
    Invoke-B3aVoid ('shutdown-verdict: ' + $postVerdict.verdict)
}
$invArgs = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\12-round-inventory.ps1'), '-RoundRoot', $RoundRootGuest, '-ProbeRoot', $ProbeRoot)
$inv1 = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' -GuestArgs $invArgs -TimeoutMs 120000
if ($inv1.exitCode -ne 0) { [void](Write-B3aRoundLog 'S8: inventory run 1 failed'); exit 5 }
$inv2 = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' -GuestArgs $invArgs -TimeoutMs 120000
if ($inv2.exitCode -ne 0) { [void](Write-B3aRoundLog 'S8: inventory run 2 failed'); exit 5 }
if (-not (Test-B3aBytesEqual -A $inv1.stdoutBytes -B $inv2.stdoutBytes)) {
    Invoke-B3aVoid 'inventory-not-repeatable'
}
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'inventory-1.json') -Bytes $inv1.stdoutBytes
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'inventory-2.json') -Bytes $inv2.stdoutBytes
$invParsed = Get-B3aJsonFromText -Text $inv1.stdoutText
if ($null -eq $invParsed) { [void](Write-B3aRoundLog 'S8: inventory JSON unparseable'); exit 5 }
$dirtyText = if ($invParsed.PSObject.Properties.Name -contains 'dirty') { [string]$invParsed.dirty } else { '' }
Write-B3aTextFile -Path (Join-Path $RoundDir 'dirty.txt') -Text $dirtyText
[void](Write-B3aRoundLog ('S8: inventory x2 byte-identical; dirty=' + $dirtyText))

# S9 reconcile
$recRes = Invoke-B3aReconcile -PlanJsonPath (Join-Path $RoundDir 'plan.json') `
    -InventoryJsonPath (Join-Path $RoundDir 'inventory-1.json') `
    -InventoryJson2Path (Join-Path $RoundDir 'inventory-2.json') `
    -AuditTxtPath (Join-Path $RoundDir 'audit.txt') `
    -OutputPath (Join-Path $RoundDir 'reconcile.json') `
    -JournalPath $JournalPath -Round $Round -Candidate $Candidate -Stage $stage `
    -TargetPath $targetPath -ContentDigest $contentDigestKnown -ContentSize $contentSizeKnown
if (-not $recRes.ok) { [void](Write-B3aRoundLog ('BLOCKER: reconcile failed: ' + $recRes.error)); exit 6 }
[void](Write-B3aRoundLog ('S9: outcome={0}' -f $recRes.outcome))

# post-round baseline reset (design S5.4)
if (-not (Invoke-B3aBaselineReset)) { exit 6 }
[void](Write-B3aRoundLog ('=== round {0} done: outcome={1} ===' -f $Round, $recRes.outcome))

if ($Candidate -eq 'parity' -and $recRes.outcome -eq 'survived') {
    'PAUSE: R0 parity round SURVIVED (design S4.4); operator decision required before calibration.'
    exit 4
}
exit 0

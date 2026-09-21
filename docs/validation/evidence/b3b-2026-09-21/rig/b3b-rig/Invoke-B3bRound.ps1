# Invoke-B3bRound.ps1 - one B3b round, steps S1-S9 (design section 5).
#
# The walk mirrors the B3a round step for step so the two slices stay comparable,
# with three deliberate differences:
#   - S3.5 deploys the Go writer in addition to the two reused guest scripts, and the
#     writer hash is compared on both sides.
#   - S5 asks the writer for the pre-commitment instead of recomputing it on the host:
#     the measured object is a real product record, and only the product encoder knows
#     its exact bytes. The commitment is still journaled and ACKed before any write.
#   - S6.5 polls the guest marker file until the target stage announces itself,
#     archives the transcript, cross-checks the pre-commitment, takes a survival probe
#     and only then allows the cut. A host-side timing kill cannot hit a chosen stage
#     (A7 measured CP4 and CP5 as indistinguishable), so the stage has to announce
#     itself from inside the guest.
#
# Void rules and the void vocabulary are design section 4. A void round restores
# pristine and the matrix re-runs the same logical slot with a NEW round number.
#
# Exit codes: 0 valid round finished (outcome in rounds/rNN/reconcile.json)
#             3 void round (pristine restored; re-run with a NEW round number)
#             5 operational fail (re-run the SAME round number; budget 2)
#             6 session blocker (journal/pin/protocol defect - stop everything)
#             7 operator declined the cut
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
param(
    [Parameter(Mandatory = $true)][ValidateRange(0, 999)][int]$Round,
    [Parameter(Mandatory = $true)][ValidateSet('c1', 'c2', 'c3', 'c2c3', 'positive-control')][string]$Candidate,
    [Parameter(Mandatory = $true)][ValidateSet('temp-written', 'temp-synced', 'temp-closed', 'record-linked', 'parent-synced')][string]$Stage,
    [string]$VmName = 'Win11B3',
    [string]$CredFile = 'D:\VirtualBox VMs\Win11B3\host-only\cred.txt',
    [string]$GuestUser = 'prfrail',
    [string]$RigRoot = 'D:\VirtualBox VMs\Win11B3\b3b-runs',
    [string]$VBoxManage = 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe',
    [string]$GuestToolDir = 'C:\prfrail-prep',
    # The session id is passed in by the matrix so that a run which crosses midnight keeps
    # writing into the session it started: a per-invocation date would silently split the
    # rounds across two directories and the finalize would only see half the evidence.
    [string]$SessionId = '',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep',
    [string]$StoreRoot = 'R:\agent-runner-replay',
    [string]$ScratchRoot = 'R:\calib',
    [string]$WriterExe = '',
    [double]$PlanCutDelayS = 2.0,
    [int]$StageDetectTimeoutS = 120,
    [int]$StageBlockS = 300,
    [int]$MaxReadyWaitS = 120,
    [switch]$NoConfirm,
    [switch]$Paired
)

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'Get-B3bReconcile.ps1')

if ($MyInvocation.InvocationName -eq '.') { return }

# Fail closed: an unexpected error must abort the round loudly instead of skipping a
# gate silently.
$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrEmpty($WriterExe)) {
    $repoRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\..'))
    $WriterExe = Join-Path $repoRoot 'tmp\b3b\writer.exe'
}

if ([string]::IsNullOrEmpty($SessionId)) { $SessionId = Get-Date -Format 'yyyy-MM-dd' }
$DateTag = $SessionId
$SessionDir = Join-Path $RigRoot $DateTag
$JournalPath = Join-Path $SessionDir 'journal.jsonl'
$SessionLog = Join-Path $SessionDir 'session.log'
$RoundLabel = 'r{0:D2}' -f $Round
$RoundDir = Join-Path $SessionDir (Join-Path 'rounds' $RoundLabel)
$ScratchGuest = $ScratchRoot.TrimEnd('\') + '\' + $RoundLabel
$SignalGuest = $GuestToolDir.TrimEnd('\') + '\runs\' + $RoundLabel
$MarkerGuest = $SignalGuest + '\marker.jsonl'
$PlanMarkerGuest = $SignalGuest + '\marker-plan.jsonl'
$GuestWriterExe = $GuestToolDir.TrimEnd('\') + '\writer.exe'
# The guest copy always carries this fixed name, while the host file may be named
# anything. Probing the guest by the host basename would void every round as soon as
# the operator renames the locally built binary.
$WriterName = [System.IO.Path]::GetFileName($GuestWriterExe)
$targetPath = $StoreRoot.TrimEnd('\') + '\requests\request.b3b-' + $RoundLabel + '.jsonl'
if ($Paired) {
    # Ordering round: the measured object is the completion record, and the request
    # record is published first by the same mechanism.
    $targetPath = $StoreRoot.TrimEnd('\') + '\completions\completion.b3b-' + $RoundLabel + '.jsonl'
}
$contentDigestKnown = ''
$contentSizeKnown = 0

function Write-B3bRoundLog {
    param([Parameter(Mandatory = $true)][string]$Text)
    Add-B3aLogLine -Path $SessionLog -Text ('[{0}] {1}' -f (Get-Date -Format 'HH:mm:ss'), $Text)
    $Text
}

function Invoke-B3bBaselineReset {
    # post-round / void recovery: poweroff (if running) + settle + restore pristine
    $info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
    $state = ''
    if ($info.values) { $state = $info.values['VMState'] }
    if ($state -eq 'running' -or $state -eq 'paused' -or $state -eq 'saved') {
        $cut = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('controlvm "{0}" poweroff' -f $VmName)
        [void](Write-B3bRoundLog ('reset: controlvm poweroff (was {0}) exit={1}' -f $state, $cut.exitCode))
        # The poweroff returns before the VM has finished shutting down; restoring inside that
        # window leaves a stale session that breaks the next startvm.
        [void](Wait-B3bVmPowerOff -VBoxManage $VBoxManage -VmName $VmName -TimeoutS 90)
    }
    $restore = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('snapshot "{0}" restore pristine' -f $VmName)
    [void](Write-B3bRoundLog ('reset: restore pristine exit=' + $restore.exitCode))
    return ($restore.exitCode -eq 0)
}

function Wait-B3bReady {
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

function Read-B3bGuestMarker {
    # cmd's type exits non-zero while the file does not exist yet, which is the normal
    # state before the writer reaches its first stage.
    param([Parameter(Mandatory = $true)][string]$GuestPath)
    return (Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
        -Exe 'C:\Windows\System32\cmd.exe' -GuestArgs @('/c', 'type', (Get-B3aQuoted $GuestPath)) -TimeoutMs 20000)
}

function Start-B3bWriter {
    # The writer is started without waiting: the host has to poll the marker while the
    # writer sits inside the target stage. stdout/stderr are drained asynchronously so
    # a full pipe can never stall the child.
    # ShouldProcess does not apply here: the function only launches a host process
    # handle and changes no configuration; the caller decides whether the round may
    # proceed and owns the cut. Same reasoning as the frozen New-B3aJournalRecord.
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '', Justification = 'Launching a host process handle changes no configuration; the caller owns the round decision and the cut.')]
    param([Parameter(Mandatory = $true)][hashtable]$Environment)
    $parts = @(
        'guestcontrol', (Get-B3aQuoted $VmName), 'run',
        '--exe', (Get-B3aQuoted $GuestWriterExe),
        '--username', $GuestUser,
        '--passwordfile', (Get-B3aQuoted $CredFile),
        '--wait-stdout', '--wait-stderr',
        '--timeout', ([string](($StageBlockS + 180) * 1000))
    )
    foreach ($name in $Environment.Keys) {
        $parts += ('--putenv=' + $name + '=' + [string]$Environment[$name])
    }
    $parts += '--'
    $parts += @('-test.run=TestB3bNativeGuestWriter', '-test.v')
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = $VBoxManage
    $psi.Arguments = ($parts -join ' ')
    $psi.UseShellExecute = $false
    $psi.CreateNoWindow = $true
    $psi.RedirectStandardOutput = $true
    $psi.RedirectStandardError = $true
    $process = New-Object System.Diagnostics.Process
    $process.StartInfo = $psi
    [void]$process.Start()
    return [pscustomobject]@{
        process    = $process
        stdoutTask = $process.StandardOutput.ReadToEndAsync()
        stderrTask = $process.StandardError.ReadToEndAsync()
    }
}

function Wait-B3bWriter {
    param([Parameter(Mandatory = $true)]$Handle, [int]$TimeoutMs = 60000)
    $exited = $Handle.process.WaitForExit($TimeoutMs)
    if (-not $exited) {
        try { $Handle.process.Kill() } catch {
            # the child may have exited between WaitForExit and Kill; the archive below
            # still records whatever the reader tasks produced
            Write-Verbose ('B3b: writer kill failed: ' + $_.Exception.Message)
        }
    }
    $stdout = ''
    $stderr = ''
    try { $stdout = $Handle.stdoutTask.Result } catch { Write-Verbose 'B3b: writer stdout task faulted' }
    try { $stderr = $Handle.stderrTask.Result } catch { Write-Verbose 'B3b: writer stderr task faulted' }
    $exitCode = -1
    if ($exited) { $exitCode = $Handle.process.ExitCode }
    return [pscustomobject]@{ exited = $exited; exitCode = $exitCode; stdoutText = $stdout; stderrText = $stderr }
}

function Invoke-B3bVoid {
    param([Parameter(Mandatory = $true)][string]$Reason)
    [void](Write-B3bRoundLog ('VOID: ' + $Reason))
    $rec = Invoke-B3bReconcile -OutputPath (Join-Path $RoundDir 'reconcile.json') `
        -JournalPath $JournalPath -Round $Round -Candidate $Candidate -Stage $Stage `
        -TargetPath $targetPath -ContentDigest $contentDigestKnown -ContentSize ([long]$contentSizeKnown) `
        -VoidReason $Reason -PlanCutDelayS $PlanCutDelayS
    if (-not $rec.ok) {
        [void](Write-B3bRoundLog ('BLOCKER: void reconcile failed: ' + $rec.error))
        exit 6
    }
    if (-not (Invoke-B3bBaselineReset)) { exit 6 }
    exit 3
}

# --- round start ------------------------------------------------------------
[void](New-Item -ItemType Directory -Path $RoundDir -Force)
[void](Write-B3bRoundLog ('=== round {0} arm={1} stage={2} begin ===' -f $Round, $Candidate, $Stage))

# preflight (cheap local blockers)
if (-not (Test-Path -LiteralPath $VBoxManage)) { [void](Write-B3bRoundLog 'BLOCKER: VBoxManage missing'); exit 6 }
if (-not (Test-Path -LiteralPath $CredFile)) { [void](Write-B3bRoundLog 'BLOCKER: credential file missing'); exit 6 }
if (-not (Test-Path -LiteralPath $WriterExe)) { [void](Write-B3bRoundLog ('BLOCKER: writer binary missing: ' + $WriterExe)); exit 6 }
$chain = Test-B3aJournalChain -JournalPath $JournalPath
if (-not $chain.ok) {
    [void](Write-B3bRoundLog ('BLOCKER: journal chain broken at line {0}: {1}' -f $chain.firstBadLine, $chain.firstError))
    exit 6
}

# S1 restore baseline
$info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
if ($info.exitCode -ne 0) { [void](Write-B3bRoundLog 'BLOCKER: cannot query VM'); exit 6 }
$state = ''
if ($info.values) { $state = $info.values['VMState'] }
if ($state -ne 'poweroff') {
    $pre = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('controlvm "{0}" poweroff' -f $VmName)
    [void](Write-B3bRoundLog ('S1: pre-round poweroff (was {0}) exit={1}' -f $state, $pre.exitCode))
    if (-not (Wait-B3bVmPowerOff -VBoxManage $VBoxManage -VmName $VmName -TimeoutS 90)) {
        [void](Write-B3bRoundLog 'BLOCKER: the VM did not reach poweroff; refusing to restore the snapshot on a moving target')
        exit 6
    }
}
# A stale headless session left by a hard poweroff has to go before the restore: restoring
# under it closes the session and the next startvm then fails with E_FAIL / SessionMachine.
$stale = @(Clear-B3bStaleVmSession -VBoxManage $VBoxManage -VmName $VmName)
if ($stale.Count -gt 0) {
    [void](Write-B3bRoundLog ('S1: cleared {0} stale headless session(s): {1}' -f $stale.Count, ($stale -join ',')))
}
$restore = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('snapshot "{0}" restore pristine' -f $VmName)
[void](Write-B3bRoundLog ('S1: restore pristine exit=' + $restore.exitCode))
if ($restore.exitCode -ne 0) { [void](Write-B3bRoundLog 'BLOCKER: restore pristine failed'); exit 6 }

# S2 start (bounded retry: a hard poweroff can leave a stale headless session that makes
# the first startvm fail with E_FAIL / SessionMachine)
$start = Start-B3bVm -VBoxManage $VBoxManage -VmName $VmName
[void](Write-B3bRoundLog ('S2: startvm ok={0} attempts={1}' -f $start.ok, $start.attempts))
if (-not $start.ok) { [void](Write-B3bRoundLog ('S2: startvm failed: ' + $start.lastError)); exit 5 }

# S3 ready wait
$readyS = Wait-B3bReady -TimeoutS $MaxReadyWaitS
if ($readyS -lt 0) { [void](Write-B3bRoundLog 'S3: ready wait timed out'); exit 5 }
[void](Write-B3bRoundLog ('S3: guestcontrol ready after {0:N1}s' -f $readyS))

# S3.5 deployment: the Go writer plus the reused guest scripts. restore pristine
# reverts C:\prfrail-prep, so this repeats every round.
$deployOk = $true
$hostWriterSha = Get-B3aSha256Hex -Path $WriterExe
[void](Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\cmd.exe' -GuestArgs @('/c', 'del', '/q', (Get-B3aQuoted $GuestWriterExe)) -TimeoutMs 30000)
$copyArgs = 'guestcontrol "{0}" copyto --username {1} --passwordfile "{2}" "{3}" "{4}"' -f $VmName, $GuestUser, $CredFile, $WriterExe, $GuestWriterExe
$copy = Invoke-B3aProcess -FileName $VBoxManage -Arguments $copyArgs
if ($copy.exitCode -ne 0) {
    [void](Write-B3bRoundLog ('S3.5: writer copyto exit=' + $copy.exitCode))
    $deployOk = $false
} else {
    $hashCommand = '(Get-FileHash -Algorithm SHA256 -LiteralPath ''' + $GuestWriterExe + ''').Hash.ToLower()'
    $guestHash = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
        -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
        -GuestArgs @('-NoProfile', '-Command', $hashCommand) -TimeoutMs 60000
    $guestWriterSha = ($guestHash.stdoutText -replace '\s', '').ToLowerInvariant()
    [void](Write-B3bRoundLog ('S3.5: writer host={0} guest={1}' -f $hostWriterSha, $guestWriterSha))
    if ($guestWriterSha -ne $hostWriterSha) { $deployOk = $false }
}
$deploy = Invoke-B3aToolDeploy -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -GuestToolDir $GuestToolDir -GuestPrepDir $GuestPrepDir -ScriptsDir (Join-Path $PSScriptRoot '..\b3a-rig')
foreach ($f in $deploy.files) {
    [void](Write-B3bRoundLog ('S3.5: deploy {0} copied={1} hashMatch={2} note={3}' -f $f.name, $f.copied, $f.hashMatch, $f.note))
}
if (-not $deploy.ok) { $deployOk = $false }
if (-not $deployOk) { [void](Write-B3bRoundLog 'S3.5: deployment failed'); exit 5 }

# S4 pre-injection gate
$auditPre = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
    -GuestArgs @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\11-shutdown-audit.ps1')) -TimeoutMs 60000
if ($auditPre.exitCode -ne 0) { [void](Write-B3bRoundLog 'S4: pre-injection audit failed'); exit 5 }
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'audit-pre.txt') -Bytes $auditPre.stdoutBytes
$preVerdict = Read-B3aAuditVerdict -Text $auditPre.stdoutText
[void](Write-B3bRoundLog ('S4: pre-audit verdict={0} uptime={1}' -f $preVerdict.verdict, $preVerdict.uptimeMinutes))
if ($preVerdict.verdict -eq 'LICENCE-SHUTDOWN') { Invoke-B3bVoid 'licence-shutdown' }
if ($null -eq $preVerdict.uptimeMinutes) { [void](Write-B3bRoundLog 'S4: uptime unparseable'); exit 5 }
if ($preVerdict.uptimeMinutes -ge 50) { Invoke-B3bVoid 'uptime-too-old' }

# S4.5 rig-side housekeeping: the scratch directory on the measured volume and the
# signaling directory on C:. Neither is part of the measured write, and both are
# created before the plan record so the measured publication stays exactly one record
# write.
foreach ($path in @($ScratchGuest, $SignalGuest)) {
    $mkdir = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
        -Exe 'C:\Windows\System32\cmd.exe' -GuestArgs @('/c', 'mkdir', (Get-B3aQuoted $path)) -TimeoutMs 30000
    [void](Write-B3bRoundLog ('S4.5: mkdir {0} exit={1}' -f $path, $mkdir.exitCode))
    if ($mkdir.exitCode -ne 0) { exit 5 }
}

# S5 pre-commitment: the writer builds the record, prints PLAN and stops without
# publishing. The host journals the digest before anything is written.
$planEnv = @{
    B3B_GUEST_WRITER = '1'; B3B_ROUND = $RoundLabel; B3B_CANDIDATE = $Candidate
    B3B_STAGE = 'temp-written'; B3B_ROOT = $StoreRoot; B3B_MARKER = $PlanMarkerGuest; B3B_PLAN_ONLY = '1'
}
if ($Paired) { $planEnv['B3B_PAIRED'] = '1' }
$planRun = Wait-B3bWriter -Handle (Start-B3bWriter -Environment $planEnv) -TimeoutMs 120000
Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-plan-stdout.txt') -Text $planRun.stdoutText
if ($planRun.exitCode -ne 0) {
    Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-plan-stderr.txt') -Text $planRun.stderrText
    [void](Write-B3bRoundLog ('S5: plan run exit=' + $planRun.exitCode))
    exit 5
}
$plan = Get-B3bPlanFromText -Text $planRun.stdoutText
if ($null -eq $plan) { [void](Write-B3bRoundLog 'S5: PLAN line not found in the writer output'); exit 5 }
if ([string]$plan.path -ne $targetPath) {
    [void](Write-B3bRoundLog ('BLOCKER: writer plan path {0} != expected {1}' -f $plan.path, $targetPath))
    exit 6
}
$contentDigestKnown = [string]$plan.digest
$contentSizeKnown = [long]$plan.size
Write-B3aTextFile -Path (Join-Path $RoundDir 'plan-source.json') -Text ($plan | ConvertTo-Json -Compress)
$planRec = New-B3aJournalRecord -Round $Round -Candidate $Candidate -Stage $Stage `
    -TargetPath $targetPath -ContentDigest $contentDigestKnown -ContentSize $contentSizeKnown `
    -PlanCutDelayS $PlanCutDelayS -Op 'plan'
$planAck = Add-B3aJournalRecord -JournalPath $JournalPath -Record $planRec
if (-not $planAck.ack) { Invoke-B3bVoid ('plan-no-ack: ' + $planAck.error) }
Write-B3aTextFile -Path (Join-Path $RoundDir 'plan.json') -Text $planAck.recordJson
[void](Write-B3bRoundLog ('S5: plan ACK ok digest={0} size={1}' -f $contentDigestKnown, $contentSizeKnown))

# S6 measured publication: the writer publishes the record and blocks inside the
# target stage. Its own PLAN line repeats the pre-commitment inside the measured
# transcript, which gives the host a second, independent digest check.
$writerEnv = @{
    B3B_GUEST_WRITER = '1'; B3B_ROUND = $RoundLabel; B3B_CANDIDATE = $Candidate
    B3B_STAGE = $Stage; B3B_ROOT = $StoreRoot; B3B_MARKER = $MarkerGuest; B3B_BLOCK_SECONDS = ([string]$StageBlockS)
}
if ($Paired) { $writerEnv['B3B_PAIRED'] = '1' }
$writerHandle = Start-B3bWriter -Environment $writerEnv
[void](Write-B3bRoundLog 'S6: writer started; polling the marker')
# S6.5 marker detection
$stageSeen = $false
$escape = $false
$transcriptText = ''
$detectStart = Get-Date
$deadline = $detectStart.AddSeconds($StageDetectTimeoutS)
while ((Get-Date) -lt $deadline) {
    $read = Read-B3bGuestMarker -GuestPath $MarkerGuest
    if ($read.exitCode -eq 0) {
        $transcriptText = $read.stdoutText
        $parsed = Read-B3bMarkerTranscript -Text $transcriptText
        $stages = Get-B3bMeasuredPhaseStageList -Entries $parsed.entries -TargetPath $targetPath
        if (Test-B3bStageEscape -Stages $stages -Target $Stage) { $escape = $true; break }
        if (Get-B3bTargetStageSeen -Stages $stages -Target $Stage) { $stageSeen = $true; break }
    }
    Start-Sleep -Milliseconds 2000
}
$detectS = ((Get-Date) - $detectStart).TotalSeconds
Write-B3bArtifactText -Path (Join-Path $RoundDir 'marker-transcript.txt') -Text $transcriptText
[void](Write-B3bRoundLog ('S6.5: detect after {0:N1}s stageSeen={1} escape={2}' -f $detectS, $stageSeen, $escape))
if ($escape) {
    # A writer that continued past the stage it was told to stop in is a protocol
    # defect: the whole session stops rather than turning it into a void round.
    $stray = Wait-B3bWriter -Handle $writerHandle -TimeoutMs 30000
    Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-stdout.txt') -Text $stray.stdoutText
    Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-stderr.txt') -Text $stray.stderrText
    [void](Write-B3bRoundLog 'BLOCKER: stage-escape')
    exit 6
}
if (-not $stageSeen) {
    $slow = Wait-B3bWriter -Handle $writerHandle -TimeoutMs 30000
    Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-stdout.txt') -Text $slow.stdoutText
    Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-stderr.txt') -Text $slow.stderrText
    # Tell a genuinely silent writer apart from one that already failed: a candidate
    # that errors out on every round would otherwise surface only as "void budget
    # exhausted", which hides the real cause and wastes the whole cell.
    $failureText = ''
    foreach ($line in @((($slow.stdoutText + "`n" + $slow.stderrText) -split "`n"))) {
        if ($line -match 'FAIL|panic:|unknown B3B') { $failureText = $line.Trim(); break }
    }
    if (-not [string]::IsNullOrEmpty($failureText)) {
        Invoke-B3bVoid ('writer-failed-before-stage: ' + $failureText)
    }
    Invoke-B3bVoid 'marker-missed'
}
$transcriptEntries = (Read-B3bMarkerTranscript -Text $transcriptText).entries
if (@($transcriptEntries).Count -lt 1) { Invoke-B3bVoid 'marker-transcript-empty' }
$measuredPlan = $null
foreach ($entry in @($transcriptEntries)) {
    if ($null -ne $entry.PSObject.Properties['digest']) { $measuredPlan = $entry; break }
}
if ($null -eq $measuredPlan) { Invoke-B3bVoid 'plan-line-missing-in-transcript' }
if ([string]$measuredPlan.digest -ne $contentDigestKnown) {
    Invoke-B3bVoid ('plan-digest-mismatch journal={0} transcript={1}' -f $contentDigestKnown, $measuredPlan.digest)
}

# S6.5 survival probe: the writer must still be alive, which is what makes "the cut
# landed inside the target stage" a claim the evidence supports.
$aliveProbe = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\cmd.exe' `
    -GuestArgs @('/c', 'tasklist', '/FI', (Get-B3aQuoted ('IMAGENAME eq ' + $WriterName)), '/NH') -TimeoutMs 30000
Write-B3bArtifactText -Path (Join-Path $RoundDir 'survival-probe.txt') -Text $aliveProbe.stdoutText
$alive = ($aliveProbe.exitCode -eq 0) -and ($aliveProbe.stdoutText -match [regex]::Escape($WriterName))
if (-not $alive) { Invoke-B3bVoid 'writer-exited-before-cut' }
[void](Write-B3bRoundLog 'S6.5: writer alive at the cut point')

# S7 operator confirmation + hard cut
if (-not $NoConfirm) {
    '--- plan record (journal ACK ok) ---'
    $planAck.recordJson
    ''
    '--- marker transcript ---'
    $transcriptText
    ''
    $answer = Read-Host ('Hard poweroff of {0} in {1} s? Type YES to cut, anything else aborts' -f $VmName, $PlanCutDelayS)
    if ($answer -ne 'YES') {
        [void](Write-B3bRoundLog 'S7: operator declined the cut; VM left running')
        'aborted by operator; VM left running (run the same round again to continue)'
        exit 7
    }
}
# The freeze delay is applied first, so the probe-to-cut gap is only the duration of the probe
# itself: the re-check is meaningful exactly because nothing else runs between it and the cut.
Start-Sleep -Seconds $PlanCutDelayS
$aliveProbe2 = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\cmd.exe' `
    -GuestArgs @('/c', 'tasklist', '/FI', (Get-B3aQuoted ('IMAGENAME eq ' + $WriterName)), '/NH') -TimeoutMs 30000
Write-B3bArtifactText -Path (Join-Path $RoundDir 'survival-probe-2.txt') -Text $aliveProbe2.stdoutText
$alive2 = ($aliveProbe2.exitCode -eq 0) -and ($aliveProbe2.stdoutText -match [regex]::Escape($WriterName))
if (-not $alive2) { Invoke-B3bVoid 'writer-exited-before-cut (second probe, immediately before the cut)' }
$cut = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('controlvm "{0}" poweroff' -f $VmName)
[void](Write-B3bRoundLog ('S7: controlvm poweroff exit=' + $cut.exitCode))
if ($cut.exitCode -ne 0) { Invoke-B3bVoid ('poweroff-nonzero exit=' + $cut.exitCode) }
$cutRec = New-B3aJournalRecord -Round $Round -Candidate $Candidate -Stage $Stage `
    -TargetPath $targetPath -ContentDigest $contentDigestKnown -ContentSize $contentSizeKnown `
    -PlanCutDelayS $PlanCutDelayS -Op 'cut' `
    -Extra ([ordered]@{ cutExitCode = $cut.exitCode; cutHostTs = (Get-B3aHostTimestamp) })
$cutAck = Add-B3aJournalRecord -JournalPath $JournalPath -Record $cutRec
if (-not $cutAck.ack) { [void](Write-B3bRoundLog ('BLOCKER: cut record ACK failed: ' + $cutAck.error)); exit 6 }
Write-B3aTextFile -Path (Join-Path $RoundDir 'cut.json') -Text $cutAck.recordJson
$writerOut = Wait-B3bWriter -Handle $writerHandle -TimeoutMs 60000
Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-stdout.txt') -Text $writerOut.stdoutText
Write-B3bArtifactText -Path (Join-Path $RoundDir 'writer-stderr.txt') -Text $writerOut.stderrText
# A writer that reported TIMEOUT never stayed in the stage it was told to block in, so
# the cut cannot have landed inside the target stage. The reason is recovered from the
# writer's own output and turned into a void round instead of a valid one.
$timeoutMatch = [regex]::Match(($writerOut.stdoutText + "`n" + $writerOut.stderrText), 'TIMEOUT\s+([A-Za-z-]+)')
$writerVoidReason = $null
if ($timeoutMatch.Success) { $writerVoidReason = ('writer-exited-before-cut: TIMEOUT ' + $timeoutMatch.Groups[1].Value) }
[void](Write-B3bRoundLog ('S7: cut record ACK ok; writer exit={0}' -f $writerOut.exitCode))

# S8 restart + evidence copy-back + audit + inventory x2
$start2 = Start-B3bVm -VBoxManage $VBoxManage -VmName $VmName
if (-not $start2.ok) { [void](Write-B3bRoundLog ('S8: restart failed: ' + $start2.lastError)); exit 5 }
$readyS2 = Wait-B3bReady -TimeoutS $MaxReadyWaitS
if ($readyS2 -lt 0) { [void](Write-B3bRoundLog 'S8: ready wait after restart timed out'); exit 5 }
$pullPairs = @(
    @{ Guest = $MarkerGuest; Host = (Join-Path $RoundDir 'marker.txt') },
    @{ Guest = $PlanMarkerGuest; Host = (Join-Path $RoundDir 'marker-plan.txt') }
)
foreach ($pair in $pullPairs) {
    $pullArgs = 'guestcontrol "{0}" copyfrom --username {1} --passwordfile "{2}" "{3}" "{4}"' -f $VmName, $GuestUser, $CredFile, $pair.Guest, $pair.Host
    $pull = Invoke-B3aProcess -FileName $VBoxManage -Arguments $pullArgs
    [void](Write-B3bRoundLog ('S8: copyfrom {0} exit={1}' -f $pair.Guest, $pull.exitCode))
    if ($pull.exitCode -ne 0) { [void](Write-B3bRoundLog 'S8: marker copy-back failed'); exit 5 }
}
$audit = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
    -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
    -GuestArgs @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\11-shutdown-audit.ps1')) -TimeoutMs 60000
if ($audit.exitCode -ne 0) { [void](Write-B3bRoundLog 'S8: post-cut audit failed'); exit 5 }
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'audit.txt') -Bytes $audit.stdoutBytes
$postVerdict = Read-B3aAuditVerdict -Text $audit.stdoutText
[void](Write-B3bRoundLog ('S8: post-audit verdict={0}' -f $postVerdict.verdict))
if ($postVerdict.verdict -ne 'HARD-POWER-LOSS') { Invoke-B3bVoid ('shutdown-verdict: ' + $postVerdict.verdict) }
$invArgs = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', ($GuestToolDir + '\12-round-inventory.ps1'), '-RoundRoot', $ScratchGuest, '-ProbeRoot', $StoreRoot)
# Read until two consecutive reads agree, then record that pair. What has to be identical is
# the pair the round records, not the first two reads (see Get-B3bSettledInventoryPair).
$invReads = New-Object System.Collections.Generic.List[object]
$invSettled = $null
$invFailures = New-Object System.Collections.Generic.List[string]
for ($invAttempt = 1; $invAttempt -le $B3bInventoryMaxReads; $invAttempt++) {
    $invRun = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
        -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' -GuestArgs $invArgs -TimeoutMs 120000
    if ($invRun.exitCode -ne 0) {
        # A failed read attempt is retried inside the round instead of failing the round: the
        # guest management channel is occasionally not ready right after the reboot, and
        # redoing the whole round for that costs minutes and a fresh cut. The exit code and a
        # bounded excerpt are logged, because the previous message said only 'failed' and left
        # the cause unknowable (same defect class as the dropped round stderr).
        $detail = ('attempt {0}: exit={1} stderr={2} stdout={3}' -f $invAttempt, $invRun.exitCode, `
                (Get-B3aBoundedText -Text $invRun.stderrText), (Get-B3aBoundedText -Text $invRun.stdoutText))
        $invFailures.Add($detail)
        [void](Write-B3bRoundLog ('S8: inventory read failed; ' + $detail))
        continue
    }
    $invReads.Add([pscustomobject]@{ bytes = $invRun.stdoutBytes; text = $invRun.stdoutText })
    Write-B3aRawTextFile -Path (Join-Path $RoundDir ('inventory-attempt-{0}.bin' -f $invAttempt)) -Bytes $invRun.stdoutBytes
    $invSettled = Get-B3bSettledInventoryPair -Reads $invReads.ToArray()
    if ($invSettled.stable) {
        [void](Write-B3bRoundLog ('S8: inventory settled after {0} read(s)' -f $invSettled.tries))
        break
    }
}
if ($invReads.Count -eq 0) {
    Write-B3bArtifactText -Path (Join-Path $RoundDir 'inventory-failures.txt') -Text ($invFailures -join "`n")
    [void](Write-B3bRoundLog 'S8: every inventory read failed; the round cannot be interpreted')
    exit 5
}
if ($null -eq $invSettled) { $invSettled = Get-B3bSettledInventoryPair -Reads $invReads.ToArray() }
Write-B3bArtifactText -Path (Join-Path $RoundDir 'inventory-reads.txt') -Text ([string]$invSettled.tries)
if (-not $invSettled.stable) {
    # Keep both raw streams when no pair settles, under names the matrix does not treat as
    # round evidence. Honest c1 rounds voided on this reason with nothing left to
    # inspect, so 'the store is still settling' could not be told apart from 'a diagnostic
    # field is unstable' - two situations that call for opposite fixes.
    Write-B3aRawTextFile -Path (Join-Path $RoundDir 'inventory-mismatch-1.bin') -Bytes $invSettled.first.bytes
    Write-B3aRawTextFile -Path (Join-Path $RoundDir 'inventory-mismatch-2.bin') -Bytes $invSettled.second.bytes
    Invoke-B3bVoid 'inventory-not-repeatable'
}
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'inventory-1.json') -Bytes $invSettled.first.bytes
Write-B3aRawTextFile -Path (Join-Path $RoundDir 'inventory-2.json') -Bytes $invSettled.second.bytes
$invParsed = Get-B3aJsonFromText -Text $invSettled.first.text
if ($null -eq $invParsed) { [void](Write-B3bRoundLog 'S8: inventory JSON unparseable'); exit 5 }
$dirtyText = ''
if ($invParsed.PSObject.Properties.Name -contains 'dirty') { $dirtyText = [string]$invParsed.dirty }
Write-B3bArtifactText -Path (Join-Path $RoundDir 'dirty.txt') -Text $dirtyText
[void](Write-B3bRoundLog ('S8: inventory pair settled in {0} read(s); dirty={1}' -f $invSettled.tries, $dirtyText))

# S9 reconcile
$targetKind = 'request'
if ($targetPath -match '[\\/]completions[\\/]') { $targetKind = 'completion' }
$recParams = @{
    PlanJsonPath = (Join-Path $RoundDir 'plan.json')
    InventoryJsonPath = (Join-Path $RoundDir 'inventory-1.json')
    InventoryJson2Path = (Join-Path $RoundDir 'inventory-2.json')
    AuditTxtPath = (Join-Path $RoundDir 'audit.txt')
    # Attribution reads the transcript archived before the cut; the copy pulled back
    # afterwards may have lost its tail to the power loss under test.
    PreCutTranscriptPath = (Join-Path $RoundDir 'marker-transcript.txt')
    MarkerTranscriptPath = (Join-Path $RoundDir 'marker.txt')
    OutputPath = (Join-Path $RoundDir 'reconcile.json')
    JournalPath = $JournalPath
    Round = $Round
    Candidate = $Candidate
    Stage = $Stage
    TargetPath = $targetPath
    ContentDigest = $contentDigestKnown
    ContentSize = $contentSizeKnown
    TargetKind = $targetKind
    PlanCutDelayS = $PlanCutDelayS
}
if ($Paired) { $recParams['Paired'] = $true }
if (-not [string]::IsNullOrEmpty($writerVoidReason)) { $recParams['VoidReason'] = $writerVoidReason }
$recRes = Invoke-B3bReconcile @recParams
if (-not $recRes.ok) { [void](Write-B3bRoundLog ('BLOCKER: reconcile failed: ' + $recRes.error)); exit 6 }
[void](Write-B3bRoundLog ('S9: outcome={0} stageHit={1} mechanismOk={2} orphan={3}' -f `
        $recRes.outcome, $recRes.stageHit, $recRes.mechanismOk, $recRes.orphanCompletion))

# post-round baseline reset
if (-not (Invoke-B3bBaselineReset)) { exit 6 }
[void](Write-B3bRoundLog ('=== round {0} done: outcome={1} ===' -f $Round, $recRes.outcome))
exit 0

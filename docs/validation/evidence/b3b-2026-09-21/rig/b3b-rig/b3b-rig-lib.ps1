# B3b rig library: the helpers that are specific to the durability-candidate
# matrix. It dot-sources the frozen B3a journal library on purpose: the journal
# convention (canonical field order, chained recordHash, write-then-read-back)
# must have exactly one implementation, and the B3a evidence bundle pins the bytes
# of that implementation.
#
# Reuse boundaries (see docs/t027/B3B_DESIGN.md section 8):
#   - The frozen writer drops any field that is not in Get-B3aRecordFieldOrder, and
#     New-B3aJournalRecord only accepts -Op plan/cut/reconcile. B3b therefore maps
#     its semantics onto the existing keys instead of extending the record shape.
#   - A dot-sourced module must not declare a top-level param() (B3a lesson D1:
#     it silently rebinds same-named caller variables).
#   - B3b helpers are all prefixed B3b so a dot-source can never shadow a B3a
#     helper.
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

. (Join-Path $PSScriptRoot '..\b3a-rig\b3a-journal-lib.ps1')

# ---------------------------------------------------------------------------
# frozen-dependency pinning

function Get-B3bPinnedFileList {
    # Files B3b reuses from B3a, relative to the B3a rig directory. Their bytes must
    # not drift: they are part of the B3a evidence bundle, and they carry the only
    # journal implementation.
    return @(
        'b3a-rig-lib.ps1',
        'b3a-journal-lib.ps1',
        'guest\12-round-inventory.ps1',
        'guest\13-round-write.ps1'
    )
}

function Get-B3bReusePinMap {
    param([Parameter(Mandatory = $true)][string]$B3aRigDir)
    $hashes = @{}
    foreach ($relative in (Get-B3bPinnedFileList)) {
        $path = Join-Path $B3aRigDir $relative
        if (-not (Test-Path -LiteralPath $path)) {
            $hashes[$relative] = ''
            continue
        }
        $hashes[$relative] = Get-B3aSha256Hex -Path $path
    }
    return $hashes
}

function Write-B3bArtifactText {
    # Write a round artifact that may legitimately be empty. The frozen Write-B3aTextFile
    # rejects an empty -Text through its parameter binding, and several artifacts here are
    # empty by design: an empty marker transcript (nothing was read yet), an empty writer
    # stdout, an empty dirty list. The artifact has to exist and be empty in that case, not
    # turn a normal empty result into an uncaught exception that kills the round after it has
    # already spent a cut.
    param(
        [Parameter(Mandatory = $true)][string]$Path,
        [AllowEmptyString()][string]$Text = '',
        [switch]$WithBom
    )
    $parent = Split-Path -Parent $Path
    if ($parent -and (-not (Test-Path -LiteralPath $parent))) { [void](New-Item -ItemType Directory -Path $parent -Force) }
    $normalized = ([string]$Text).Replace("`r`n", "`n")
    [IO.File]::WriteAllText($Path, $normalized, [Text.UTF8Encoding]::new([bool]$WithBom))
}

function Enter-B3bSessionLock {
    # One session, one rig process. Two concurrent invocations restore and start the same VM at
    # the same time, which produces 'cannot query VM' blockers and half-written rounds - that is
    # not hypothetical, it is what happened the first time the handshake smoke was retried
    # while the previous run was still alive.
    #
    # The lock is taken and deliberately not released: a later run decides whether the recorded
    # pid is still alive, which is also what makes a crashed run recoverable without a stale
    # file having to be deleted by hand.
    param(
        [Parameter(Mandatory = $true)][string]$SessionDir,
        [Parameter(Mandatory = $true)][string]$Owner
    )
    if (-not (Test-Path -LiteralPath $SessionDir)) { [void](New-Item -ItemType Directory -Path $SessionDir -Force) }
    $lockPath = Join-Path $SessionDir 'session.lock'
    if (Test-Path -LiteralPath $lockPath) {
        $existing = $null
        try { $existing = (Get-B3aFileTextUtf8 -Path $lockPath) | ConvertFrom-Json } catch { $existing = $null }
        $holderPid = 0
        if (($null -ne $existing) -and ($null -ne $existing.PSObject.Properties['pid'])) { $holderPid = [int]$existing.pid }
        # A lock that was copied together with the session tree belongs to another directory: it
        # must not be able to block the session it now sits in. The fingerprint makes that
        # decidable instead of leaving it to whoever remembers to delete the file.
        $holderDir = ''
        if (($null -ne $existing) -and ($null -ne $existing.PSObject.Properties['sessionDir'])) { $holderDir = [string]$existing.sessionDir }
        $sameDirectory = ([string]::IsNullOrEmpty($holderDir)) -or ($holderDir -eq $SessionDir)
        if ($holderPid -gt 0 -and $sameDirectory) {
            $holder = Get-Process -Id $holderPid -ErrorAction SilentlyContinue
            if ($null -ne $holder) {
                $holderOwner = ''
                if (($null -ne $existing) -and ($null -ne $existing.PSObject.Properties['owner'])) { $holderOwner = [string]$existing.owner }
                return [ordered]@{ ok = $false; reason = ('the session lock is held by a live process: pid {0} ({1})' -f $holderPid, $holderOwner) }
            }
        }
    }
    $doc = [ordered]@{
        schema = 'b3b-session-lock'; v = 1; pid = $PID; owner = $Owner
        sessionDir = $SessionDir
        host = $env:COMPUTERNAME; acquiredAt = (Get-B3aHostTimestamp)
    }
    Write-B3aTextFile -Path $lockPath -Text ($doc | ConvertTo-Json -Compress) -WithBom
    return [ordered]@{ ok = $true; reason = ''; path = $lockPath }
}

function Exit-B3bSessionLock {
    # Only the holder may release it, so a run that lost the race cannot free the winner's lock.
    param([Parameter(Mandatory = $true)][string]$SessionDir)
    $lockPath = Join-Path $SessionDir 'session.lock'
    if (-not (Test-Path -LiteralPath $lockPath)) { return }
    $doc = $null
    try { $doc = (Get-B3aFileTextUtf8 -Path $lockPath) | ConvertFrom-Json } catch { return }
    if (($null -ne $doc) -and ($null -ne $doc.PSObject.Properties['pid']) -and ([int]$doc.pid -eq $PID)) {
        Remove-Item -LiteralPath $lockPath -Force
    }
}

function Start-B3bVmIfOff {
    # Bring the experiment VM up headless if it is not already running. The preflight needs
    # the guest reachable to probe the injection volume, and a probe that runs while the VM
    # is powered off would report a false failure.
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '', Justification = 'Starting the dedicated experiment VM is the action the operator ran the preflight to perform; no other VM is touched.')]
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName
    )
    $info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
    $state = ''
    if ($info.values) { $state = $info.values['VMState'] }
    if ($state -eq 'running') { return $true }
    return (Start-B3bVm -VBoxManage $VBoxManage -VmName $VmName).ok
}

function Wait-B3bVmPowerOff {
    # controlvm poweroff returns immediately, while the VM is still shutting down. A snapshot
    # restore that runs in that window races the shutdown and leaves a stale session behind,
    # and the next startvm then fails with E_FAIL / SessionMachine ('The VM session was closed
    # before any attempt to power it on'). Waiting for the state to settle avoids it.
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName,
        [int]$TimeoutS = 90
    )
    $deadline = (Get-Date).AddSeconds($TimeoutS)
    while ((Get-Date) -lt $deadline) {
        $info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
        $state = ''
        if ($info.values) { $state = $info.values['VMState'] }
        if ($state -eq 'poweroff') { return $true }
        Start-Sleep -Seconds 2
    }
    return $false
}

function Test-B3bVmCommandLineMatch {
    # Does a headless process's command line name this VM? The launcher is given the uuid or the
    # name, and a bare substring test would also match a *different* VM whose name merely starts
    # with ours (Win11B3 vs Win11B30), so both forms are matched with word boundaries. This is a
    # separate function so the self-test can pin the boundaries instead of trusting the regex.
    param(
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$CommandLine,
        [Parameter(Mandatory = $true)][string]$VmName,
        [AllowEmptyString()][string]$VmId = ''
    )
    if ([string]::IsNullOrEmpty($CommandLine)) { return $false }
    if ($CommandLine -match ('\b' + [regex]::Escape($VmName) + '\b')) { return $true }
    if (-not [string]::IsNullOrEmpty($VmId)) {
        if ($CommandLine -match ('\b' + [regex]::Escape($VmId) + '\b')) { return $true }
    }
    return $false
}

function Clear-B3bStaleVmSession {
    # A hard poweroff can leave the VM's headless process alive. While it exists, a snapshot
    # restore closes the session underneath it and the next startvm fails with E_FAIL /
    # SessionMachine ('The VM session was closed before any attempt to power it on').
    #
    # The VM name alone is not enough to recognise the holder: the launcher is given the uuid.
    # Both are matched, and when no VM at all is running every headless process is stale by
    # definition - which is also the only way to clear a holder whose command line could not be
    # read. Another VM that is running is never touched unless its command line matches ours.
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName
    )
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSShouldProcess', '', Justification = 'Analyser hint only: this function only terminates a stale headless session of the experiment VM.')]
    $uuid = ''
    $info = Get-B3aVmInfo -VBoxManage $VBoxManage -VmName $VmName
    if ($info.values) { $uuid = [string]$info.values['UUID'] }
    $runningProbe = Invoke-B3aProcess -FileName $VBoxManage -Arguments 'list runningvms'
    if ($runningProbe.exitCode -ne 0) {
        # Without a trustworthy answer there is no safe way to tell a stale session of this VM
        # from a live one belonging to another VM, so nothing is cleared.
        Write-Verbose 'B3b: list runningvms failed; refusing to clear any headless session'
        return @()
    }
    $running = @($runningProbe.stdoutText -split "`n" | Where-Object { $_.Trim() -ne '' })
    $cleared = New-Object System.Collections.Generic.List[int]
    foreach ($holder in @(Get-Process -Name VBoxHeadless -ErrorAction SilentlyContinue)) {
        $holderCommand = ''
        try {
            $holderCommand = [string](Get-CimInstance -ClassName Win32_Process -Filter ('ProcessId = ' + $holder.Id)).CommandLine
        } catch {
            $holderCommand = ''
        }
        $matchesThisVm = Test-B3bVmCommandLineMatch -CommandLine $holderCommand -VmName $VmName -VmId $uuid
        $unidentifiable = [string]::IsNullOrEmpty($holderCommand)
        $noVmRunning = ($running.Count -eq 0)
        # A readable command line that names another VM is left alone even when the running list
        # is empty: another VM can be mid-startup and not listed yet. Only a matching process, or
        # an unidentifiable one while nothing at all is running, is the stale session this
        # function exists to remove.
        $safeToClear = $matchesThisVm -or ($unidentifiable -and $noVmRunning)
        if (-not $safeToClear) { continue }
        try {
            Stop-Process -Id $holder.Id -Force -ErrorAction Stop
            $cleared.Add($holder.Id)
        } catch {
            Write-Verbose ('B3b: stale headless session could not be stopped (pid ' + $holder.Id + ')')
        }
    }
    return $cleared.ToArray()
}

function Start-B3bVm {
    # Start the VM headless with the retry the hard-poweroff path needs. Only a process whose
    # command line names *this* VM is cleared, so another VM on the same host is never touched.
    # Returns @{ok; attempts; lastError}.
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '', Justification = 'Starting the dedicated experiment VM is the action the caller asked for; the cleanup it may perform is limited to a headless session whose command line names this VM.')]
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName,
        [int]$MaxAttempts = 3
    )
    $lastError = ''
    for ($attempt = 1; $attempt -le $MaxAttempts; $attempt++) {
        $start = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('startvm "{0}" --type headless' -f $VmName)
        if ($start.exitCode -eq 0) { return [ordered]@{ ok = $true; attempts = $attempt; lastError = '' } }
        $lastError = Get-B3aBoundedText -Text ($start.stdoutText + ' ' + $start.stderrText)
        [void](Clear-B3bStaleVmSession -VBoxManage $VBoxManage -VmName $VmName)
        Start-Sleep -Seconds 5
    }
    return [ordered]@{ ok = $false; attempts = $MaxAttempts; lastError = $lastError }
}

function Wait-B3bGuestCommandReady {
    # Wait until guestcontrol can actually run a trivial command in the guest. The guest
    # process is not available the instant startvm returns, and the round driver waits for
    # exactly this condition before it does anything real.
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName,
        [Parameter(Mandatory = $true)][string]$CredFile,
        [Parameter(Mandatory = $true)][string]$GuestUser,
        [int]$TimeoutS = 240
    )
    $deadline = (Get-Date).AddSeconds($TimeoutS)
    while ((Get-Date) -lt $deadline) {
        $probe = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
            -Exe 'C:\Windows\System32\cmd.exe' -GuestArgs @('/c', 'ver') -TimeoutMs 20000
        if ($probe.exitCode -eq 0) { return $true }
        Start-Sleep -Seconds 5
    }
    return $false
}

function Test-B3bPinIntegrity {
    # Compare the current bytes against the session-frozen expectations. Any drift
    # blocks the session (fail closed): a changed dependency invalidates the claim
    # that the rig that ran is the rig that was reviewed.
    param(
        [Parameter(Mandatory = $true)][string]$B3aRigDir,
        [Parameter(Mandatory = $true)][hashtable]$Expected
    )
    $drift = New-Object System.Collections.Generic.List[string]
    $missing = New-Object System.Collections.Generic.List[string]
    $current = Get-B3bReusePinMap -B3aRigDir $B3aRigDir
    foreach ($relative in (Get-B3bPinnedFileList)) {
        if (-not $Expected.ContainsKey($relative) -or [string]::IsNullOrEmpty($Expected[$relative])) {
            $missing.Add($relative)
            continue
        }
        if ($current[$relative] -ne $Expected[$relative]) {
            $drift.Add(('{0}: expected {1} got {2}' -f $relative, $Expected[$relative], $current[$relative]))
        }
    }
    return [ordered]@{
        ok           = (($drift.Count -eq 0) -and ($missing.Count -eq 0))
        drift        = $drift.ToArray()
        missingPins  = $missing.ToArray()
        current      = $current
    }
}

# ---------------------------------------------------------------------------
# matrix shape

$script:B3bStageTempWritten = 'temp-written'
$script:B3bStageTempSynced = 'temp-synced'
$script:B3bStageTempClosed = 'temp-closed'
$script:B3bStageRecordLinked = 'record-linked'
$script:B3bStageParentSynced = 'parent-synced'
# Upper bound on inventory reads per round when looking for a settled pair. Three reads are
# the minimum that can find a pair the first two reads missed; four leaves one spare attempt.
# Each read costs roughly half a minute, so the bound also caps the per-round time cost.
$script:B3bInventoryMaxReads = 4

function Get-B3bPublishStageList {
    return @(
        $script:B3bStageTempWritten,
        $script:B3bStageTempSynced,
        $script:B3bStageTempClosed,
        $script:B3bStageRecordLinked,
        $script:B3bStageParentSynced
    )
}

function Get-B3bCandidateList {
    # The arms that carry a durability claim.
    return @('c1', 'c2', 'c3', 'c2c3')
}

function Get-B3bControlList {
    # Calibration arms. They claim nothing about durability: the device-handle flush
    # arm exists only to show that the rig can observe a *clean* result, which is the
    # positive half of the discrimination gate (design section 6). It has to be a
    # first-class arm on both sides of the rig - the writer must be able to run it and
    # the matrix must be able to schedule it - otherwise the gate is structurally
    # false and every arm is capped at 'unresolved' for the wrong reason.
    return @('positive-control')
}

function Get-B3bArmList {
    # Every arm the rig can schedule.
    return @((Get-B3bCandidateList) + (Get-B3bControlList))
}

function Resolve-B3bCellSpec {
    # Single source of truth for parsing '<arm>:<stage>'. The arm pattern must allow hyphens:
    # the control arm is literally named 'positive-control', and when the matrix inlined its own
    # pattern the arm became unschedulable from the command line while the self-test still
    # reported it as schedulable (it only grepped the source text). Found during the live run.
    param(
        [Parameter(Mandatory = $true)][string]$Cell,
        [string[]]$ArmList = @(),
        [string[]]$StageList = @()
    )
    if ($Cell -notmatch '^([a-z0-9-]+):([a-z-]+)$') {
        throw ('cell must look like arm:stage, got ' + $Cell)
    }
    $arm = $Matches[1]
    $stage = $Matches[2]
    if (($ArmList.Count -gt 0) -and ($ArmList -notcontains $arm)) { throw ('unknown arm ' + $arm) }
    if (($StageList.Count -gt 0) -and ($StageList -notcontains $stage)) { throw ('unknown stage ' + $stage) }
    return [pscustomobject]@{ arm = $arm; stage = $stage }
}

function Get-B3bDecisionStageList {
    # The stages where a candidate claims durability. C2 replaces the publish
    # primitive, so its claim starts once the primitive returned (record-linked);
    # C3 replaces the parent-directory step, so its claim starts after that step
    # (parent-synced). Stages 1-3 run the same code in every arm and are
    # informational.
    #
    # c2c3 claims only from parent-synced. Its primitive set is "move + directory
    # flush", and the directory flush runs in the parent-sync step, i.e. *after*
    # record-linked. Claiming record-linked for c2c3 would require the flush token in
    # a stage where the flush cannot have run yet, so that cell could never satisfy
    # the mechanism evidence and the arm would be unprovable by construction. At
    # record-linked c2c3 is byte-for-byte c2 anyway, so nothing is lost by narrowing
    # the claim.
    param([Parameter(Mandatory = $true)][string]$Candidate)
    switch ($Candidate) {
        # The comma keeps the empty array an array: without it the pipeline
        # unrolls it to $null, and every caller that reads .Count would depend on
        # that unrolling instead of on an explicit empty set.
        'c1' { return , @() }
        'c2' { return , @($script:B3bStageRecordLinked, $script:B3bStageParentSynced) }
        'c3' { return , @($script:B3bStageParentSynced) }
        'c2c3' { return , @($script:B3bStageParentSynced) }
        'positive-control' { return , @() }
        default { throw ('Get-B3bDecisionStageList: unknown candidate ' + $Candidate) }
    }
}

function Get-B3bAllowedOutcomeList {
    # Allowed sets (design section 7). Before the record is linked nothing may be
    # visible, and a visible record has to be complete; from record-linked on, the
    # record has to survive, because that is the durability the candidate claims.
    param([Parameter(Mandatory = $true)][string]$Stage)
    $early = @($script:B3bStageTempWritten, $script:B3bStageTempSynced, $script:B3bStageTempClosed)
    if ($early -contains $Stage) { return @('survived', 'lost-absent') }
    return @('survived')
}

function Test-B3bOutcomeAllowed {
    param(
        [Parameter(Mandatory = $true)][string]$Stage,
        [Parameter(Mandatory = $true)][string]$Outcome
    )
    return ((Get-B3bAllowedOutcomeList -Stage $Stage) -contains $Outcome)
}

function Get-B3bVoidReasonList {
    # The void vocabulary the round driver can actually produce. A void round is not a
    # measurement, so relabelling a measured round as void must not be possible with an
    # invented reason: the reason has to be one the driver can write. This list must cover
    # every literal the driver passes to Invoke-B3bVoid - the rig self-test cross-checks
    # the two statically, because a missing entry would turn honest void rounds into
    # integrity failures and refuse the whole finalize.
    return @(
        'licence-shutdown'
        'uptime-too-old'
        'plan-no-ack'
        'marker-missed'
        'marker-transcript-empty'
        'plan-line-missing-in-transcript'
        'plan-digest-mismatch'
        'writer-exited-before-cut'
        'writer-failed-before-stage'
        'poweroff-nonzero'
        'shutdown-verdict'
        'inventory-not-repeatable'
    )
}

function Test-B3bVoidReasonKnown {
    param([AllowEmptyString()][string]$VoidReason = '')
    if ([string]::IsNullOrWhiteSpace($VoidReason)) { return $false }
    foreach ($known in (Get-B3bVoidReasonList)) {
        if ($VoidReason.StartsWith($known, [System.StringComparison]::Ordinal)) { return $true }
    }
    return $false
}

function Get-B3bRecordOutcome {
    # Classify one observation of the published record. The zero digest is derived
    # from the expected size, so "the directory entry is there but the data never
    # reached the disk" is a visible violation rather than an unexplained digest
    # mismatch (B3a classifier semantics, reused verbatim).
    param(
        [Parameter(Mandatory = $true)][bool]$Present,
        [long]$ActualSize = 0,
        [string]$ActualDigest = '',
        [long]$ExpectedSize = 0,
        [string]$ExpectedDigest = '',
        [string]$ZeroDigest = ''
    )
    $zero = $ZeroDigest
    if ([string]::IsNullOrEmpty($zero) -and $ExpectedSize -gt 0) {
        $zero = Get-B3aZeroDigest -Size $ExpectedSize
    }
    return Get-B3aOutcome -Present $Present -ActualSize $ActualSize -ActualDigest $ActualDigest `
        -ExpectedSize $ExpectedSize -ExpectedDigest $ExpectedDigest -ZeroDigest $zero
}

function Get-B3bSettledInventoryPair {
    # Picks the pair of consecutive inventory reads a round records. The contract being
    # honoured is the B3a one (design section 2.3): the two *recorded* reads must be
    # byte-identical, which is what proves the enumeration state is settled rather than a
    # momentary read. Requiring the first two reads to agree turned out to be stronger than
    # that, and wrong for honest rounds: right after a cut at the weakest publish stage the
    # guest's first enumeration can still report a pre-settle view. Observed live at
    # c1:temp-written - both reads agreed on requestCount=0, completionCount=0, the tmp-file
    # list, dirty and lastBoot, and differed in exactly one field: the size of the unlinked
    # temp file (0 then 1082). Scanning forward to the first consecutive identical pair keeps
    # the settled-state proof, keeps the void when no pair settles, and records how many
    # reads it took instead of discarding the rounds that needed more than two.
    param(
        [Parameter(Mandatory = $true)][object[]]$Reads,
        [int]$MaxReads = 0
    )
    $limit = $Reads.Count
    if ($MaxReads -gt 0 -and $MaxReads -lt $limit) { $limit = $MaxReads }
    for ($i = 1; $i -lt $limit; $i++) {
        if (Test-B3aBytesEqual -A $Reads[$i - 1].bytes -B $Reads[$i].bytes) {
            return [pscustomobject]@{ stable = $true; tries = ($i + 1); first = $Reads[$i - 1]; second = $Reads[$i] }
        }
    }
    $first = $null
    $second = $null
    if ($limit -ge 1) { $first = $Reads[0] }
    if ($limit -ge 2) { $second = $Reads[$limit - 1] }
    return [pscustomobject]@{ stable = $false; tries = $limit; first = $first; second = $second }
}

function Get-B3bMatrixCellList {
    # The full matrix as cell descriptors. $N is the rounds per cell.
    param(
        [Parameter(Mandatory = $true)][int]$N,
        [switch]$WithCombo,
        [switch]$WithU4
    )
    if ($N -lt 5) { throw 'Get-B3bMatrixCellList: N must be >= 5 (design section 6)' }
    $cells = New-Object System.Collections.Generic.List[object]
    foreach ($candidate in @('c1', 'c2', 'c3')) {
        foreach ($stage in (Get-B3bPublishStageList)) {
            $cells.Add([ordered]@{ arm = $candidate; stage = $stage; rounds = $N; kind = 'matrix' })
        }
    }
    if ($WithCombo) {
        # c2c3 claims from parent-synced only: at record-linked the directory flush has
        # not run yet, so that cell would be c2 relabelled (see Get-B3bDecisionStageList).
        foreach ($stage in @($script:B3bStageParentSynced)) {
            $cells.Add([ordered]@{ arm = 'c2c3'; stage = $stage; rounds = $N; kind = 'matrix' })
        }
    }
    if ($WithU4) {
        # Paired publication: a complete request record followed by a completion
        # record stopped at the target stage. Only the completion stages matter.
        foreach ($candidate in @('c1', 'c2', 'c3')) {
            foreach ($stage in @($script:B3bStageRecordLinked, $script:B3bStageParentSynced)) {
                $cells.Add([ordered]@{ arm = $candidate; stage = $stage; rounds = 3; kind = 'u4' })
            }
        }
    }
    $cells.Add([ordered]@{ arm = 'positive-control'; stage = $script:B3bStageParentSynced; rounds = $N; kind = 'calibration' })
    return $cells.ToArray()
}

# ---------------------------------------------------------------------------
# writer handshake

function Get-B3bPlanFromText {
    # Parse the writer's pre-commitment line: PLAN path=<p> digest=<hex> size=<n>.
    param([Parameter(Mandatory = $true)][string]$Text)
    foreach ($lineRaw in ($Text -split "`n")) {
        $line = $lineRaw.Trim()
        if (-not $line.StartsWith('PLAN ')) { continue }
        $path = ''
        $digest = ''
        $size = 0
        $pathMatch = [regex]::Match($line, 'path=([^\s]+)')
        if ($pathMatch.Success) { $path = $pathMatch.Groups[1].Value }
        $digestMatch = [regex]::Match($line, 'digest=([0-9a-f]{64})')
        if ($digestMatch.Success) { $digest = $digestMatch.Groups[1].Value }
        $sizeMatch = [regex]::Match($line, 'size=([0-9]+)')
        if ($sizeMatch.Success) { $size = [long]$sizeMatch.Groups[1].Value }
        if ([string]::IsNullOrEmpty($path) -or [string]::IsNullOrEmpty($digest) -or $size -le 0) {
            throw ('Get-B3bPlanFromText: malformed pre-commitment line: ' + $line)
        }
        return [ordered]@{ path = $path; digest = $digest; size = $size }
    }
    return $null
}

function Read-B3bMarkerTranscript {
    # Parse the marker transcript. Unparseable lines are reported, never skipped
    # silently: the marker is the evidence that power was cut inside a stage.
    param([Parameter(Mandatory = $true)][string]$Text)
    $entries = New-Object System.Collections.Generic.List[object]
    $bad = New-Object System.Collections.Generic.List[string]
    foreach ($lineRaw in ($Text -split "`n")) {
        $line = $lineRaw.Trim()
        if ([string]::IsNullOrEmpty($line)) { continue }
        try {
            $entries.Add(($line | ConvertFrom-Json))
        } catch {
            $bad.Add($line)
        }
    }
    return [ordered]@{ entries = $entries.ToArray(); unparseable = $bad.ToArray() }
}

function Get-B3bMarkerStageList {
    # Stage names in the order the writer reported them.
    param([Parameter(Mandatory = $true)]$Entries)
    $stages = New-Object System.Collections.Generic.List[string]
    foreach ($entry in @($Entries)) {
        if ($null -ne $entry.PSObject.Properties['stage']) { $stages.Add([string]$entry.stage) }
    }
    return $stages.ToArray()
}

function Get-B3bMeasuredPhaseStageList {
    # Stage names belonging to the measured publication only.
    #
    # The paired mode publishes the request record first and the completion record
    # second, and both phases emit the same stage names. Feeding every stage from the
    # transcript into the hit test would let the request phase satisfy "the target stage
    # was reached", and the host would then cut power during the wrong publication.
    # Each phase's entries carry the path that was being published, so the phase is
    # resolved by path.
    param(
        [Parameter(Mandatory = $true)]$Entries,
        [Parameter(Mandatory = $true)][string]$TargetPath
    )
    $leaf = [System.IO.Path]::GetFileName($TargetPath)
    $stages = New-Object System.Collections.Generic.List[string]
    foreach ($entry in @($Entries)) {
        if ($null -eq $entry.PSObject.Properties['stage']) { continue }
        if ($null -eq $entry.PSObject.Properties['path']) { continue }
        if ([System.IO.Path]::GetFileName([string]$entry.path) -ne $leaf) { continue }
        $stages.Add([string]$entry.stage)
    }
    return $stages.ToArray()
}

function Test-B3bStageEscape {
    # True when a stage later than the target appears in the transcript. That is a
    # protocol defect, not a void round: the writer continued past the stage the
    # host was told to cut in.
    param(
        [Parameter(Mandatory = $true)][string[]]$Stages,
        [Parameter(Mandatory = $true)][string]$Target
    )
    $order = Get-B3bPublishStageList
    $targetIndex = [array]::IndexOf($order, $Target)
    if ($targetIndex -lt 0) { throw ('Test-B3bStageEscape: unknown target stage ' + $Target) }
    foreach ($stage in $Stages) {
        $index = [array]::IndexOf($order, $stage)
        if ($index -gt $targetIndex) { return $true }
    }
    return $false
}

function Get-B3bTargetStageSeen {
    param(
        [Parameter(Mandatory = $true)][string[]]$Stages,
        [Parameter(Mandatory = $true)][string]$Target
    )
    return ($Stages -contains $Target)
}

# ---------------------------------------------------------------------------
# round accounting

function New-B3bRoundResult {
    # One round's bookkeeping. valid=false means the round may not enter any count;
    # voidReason explains why (design section 4).
    # ShouldProcess does not apply here: this function only builds an in-memory
    # ordered hashtable (the same reason the frozen New-B3aJournalRecord carries the
    # same suppression). Persisting a round result is the caller's job.
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '', Justification = 'Pure in-memory record construction (returns an ordered hashtable); no system state is created or changed.')]
    param(
        [Parameter(Mandatory = $true)][int]$Round,
        # Arm and Stage are labels, not gates, so an empty value is allowed: the matrix
        # records a round whose artifacts are missing as invalid, and it must be able to
        # say so instead of failing parameter binding.
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$Arm,
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$Stage,
        [bool]$Valid = $false,
        [string]$VoidReason = '',
        [string]$Outcome = '',
        [bool]$MechanismOk = $false,
        [bool]$AuditOk = $false,
        [bool]$InventoryRepeatableOk = $false,
        [bool]$JournalChainOk = $false,
        # Did the writer demonstrably stop inside the target stage? A round that never
        # reached the stage cannot testify about it, so this gates the counterexample
        # rule in Get-B3bCellVerdict together with MechanismOk.
        [bool]$StageHit = $false,
        [string]$ExpectedDigest = '',
        [string]$PlanDigest = '',
        [string]$MarkerStages = ''
    )
    return [ordered]@{
        round                = $Round
        arm                  = $Arm
        stage                = $Stage
        valid                = $Valid
        voidReason           = $VoidReason
        outcome              = $Outcome
        mechanismOk          = $MechanismOk
        stageHit             = $StageHit
        auditOk              = $AuditOk
        inventoryRepeatableOk = $InventoryRepeatableOk
        journalChainOk       = $JournalChainOk
        expectedDigest       = $ExpectedDigest
        planDigest           = $PlanDigest
        markerStages         = $MarkerStages
    }
}

function Get-B3bCellVerdict {
    # Cell verdict for one (arm, stage) cell.
    #   disproven : a valid round left the allowed set **and** the round carries the
    #               evidence that makes it a counterexample (the writer reached the
    #               target stage and the candidate primitive demonstrably ran).
    #   eligible  : N valid rounds, all inside the allowed set, the mechanism ran,
    #               and the device proved it can discriminate.
    #   unresolved: everything else, including a device without discrimination
    #               power, which may never be read as a pass.
    #
    # Why a violation needs mechanism evidence: a round whose record is missing but
    # whose marker trace does not show the target stage (and the primitive token) is
    # most likely a round that died before publishing, not a durability failure. Such
    # a round must not be able to disprove a cell - it caps the cell at 'unresolved'
    # and says so.
    param(
        [Parameter(Mandatory = $true)][string]$Arm,
        [Parameter(Mandatory = $true)][string]$Stage,
        [Parameter(Mandatory = $true)]$Rounds,
        [Parameter(Mandatory = $true)][int]$RequiredN,
        [bool]$DiscriminationGate = $false
    )
    if ($RequiredN -lt 5) { throw 'Get-B3bCellVerdict: RequiredN must be >= 5' }
    $valid = 0
    $violations = New-Object System.Collections.Generic.List[object]
    $unattributable = New-Object System.Collections.Generic.List[object]
    $losses = 0
    $mechanismMissing = New-Object System.Collections.Generic.List[int]
    $incomplete = New-Object System.Collections.Generic.List[int]
    foreach ($round in @($Rounds)) {
        if (-not $round.valid) { continue }
        if (-not $round.auditOk -or -not $round.inventoryRepeatableOk -or -not $round.journalChainOk) {
            $incomplete.Add([int]$round.round)
            continue
        }
        $valid++
        if (-not (Test-B3bOutcomeAllowed -Stage $Stage -Outcome ([string]$round.outcome))) {
            if ($round.stageHit -and $round.mechanismOk) {
                $violations.Add($round)
            } else {
                $unattributable.Add($round)
            }
        }
        if ([string]$round.outcome -ne 'survived') { $losses++ }
        if (-not $round.mechanismOk -and ($Arm -ne 'c1') -and ((Get-B3bDecisionStageList -Candidate $Arm) -contains $Stage)) {
            $mechanismMissing.Add([int]$round.round)
        }
    }
    $verdict = 'unresolved'
    $reasons = New-Object System.Collections.Generic.List[string]
    if ($violations.Count -gt 0) {
        $verdict = 'disproven'
        foreach ($violation in $violations) {
            $reasons.Add(('round {0}: {1} is outside the allowed set for {2}' -f $violation.round, $violation.outcome, $Stage))
        }
    } elseif ($unattributable.Count -gt 0) {
        foreach ($round in $unattributable) {
            $reasons.Add(('round {0}: {1} is outside the allowed set but the round carries no stage/mechanism evidence (stageHit={2} mechanismOk={3}), so it cannot be attributed to the candidate' -f `
                    $round.round, $round.outcome, $round.stageHit, $round.mechanismOk))
        }
    } elseif ($valid -lt $RequiredN) {
        $reasons.Add(('valid rounds {0} < required {1}' -f $valid, $RequiredN))
    } elseif (-not $DiscriminationGate) {
        $reasons.Add('the device did not demonstrate discrimination power, so a clean cell carries no information')
    } elseif ($mechanismMissing.Count -gt 0) {
        $reasons.Add(('mechanism evidence missing in rounds {0}' -f (($mechanismMissing.ToArray()) -join ',')))
    } else {
        $verdict = 'eligible'
        $reasons.Add('all valid rounds inside the allowed set with mechanism evidence')
    }
    if ($incomplete.Count -gt 0) {
        $reasons.Add(('rounds excluded for incomplete gates: {0}' -f (($incomplete.ToArray()) -join ',')))
    }
    return [ordered]@{
        arm             = $Arm
        stage           = $Stage
        verdict         = $verdict
        validRounds     = $valid
        losses          = $losses
        violations      = $violations.Count
        unattributable  = $unattributable.Count
        discrimination  = $DiscriminationGate
        reasons         = $reasons.ToArray()
    }
}

function Get-B3bDiscriminationGate {
    # Positive control must lose nothing and the baseline negative control must lose
    # at least once, measured with the same writer and the same record-level payload.
    # Without that pair, "nothing was lost" cannot be told apart from "this rig
    # cannot see losses at this payload".
    param(
        [Parameter(Mandatory = $true)][int]$NegativeValidRounds,
        [Parameter(Mandatory = $true)][int]$NegativeLosses,
        [Parameter(Mandatory = $true)][int]$PositiveValidRounds,
        [Parameter(Mandatory = $true)][int]$PositiveLosses,
        [int]$RequiredN = 5
    )
    if ($RequiredN -lt 5) { throw 'Get-B3bDiscriminationGate: RequiredN must be >= 5' }
    $ok = ($NegativeValidRounds -ge $RequiredN) -and ($NegativeLosses -ge 1) `
        -and ($PositiveValidRounds -ge $RequiredN) -and ($PositiveLosses -eq 0)
    return [ordered]@{
        ok                   = $ok
        negativeValidRounds  = $NegativeValidRounds
        negativeLosses       = $NegativeLosses
        positiveValidRounds  = $PositiveValidRounds
        positiveLosses       = $PositiveLosses
        requiredN            = $RequiredN
    }
}

function Get-B3bArmVerdict {
    # Arm verdict from its cells.
    #   disproven : a **decision** cell was disproven. Only the decision cells can
    #               refute an arm: stages 1-3 run identical code in every arm (the
    #               candidate primitive has not been reached yet), so a violation there
    #               is a property of the device, not of the candidate, and attributing
    #               it to one arm would be a misattribution. Informational violations
    #               are reported in infoViolations and cap the arm at 'unresolved'.
    #   proven    : every decision cell is eligible and no decision cell was disproven.
    param(
        [Parameter(Mandatory = $true)][string]$Arm,
        [Parameter(Mandatory = $true)]$Cells
    )
    $decision = Get-B3bDecisionStageList -Candidate $Arm
    $disproven = New-Object System.Collections.Generic.List[string]
    $infoViolations = New-Object System.Collections.Generic.List[string]
    $eligible = New-Object System.Collections.Generic.List[string]
    $notEligible = New-Object System.Collections.Generic.List[string]
    foreach ($cell in @($Cells)) {
        $cellIsDecision = ($decision -contains [string]$cell.stage)
        if ([string]$cell.verdict -eq 'disproven') {
            if ($cellIsDecision) { $disproven.Add([string]$cell.stage) } else { $infoViolations.Add([string]$cell.stage) }
        }
        if ($cellIsDecision) {
            if ([string]$cell.verdict -eq 'eligible') { $eligible.Add([string]$cell.stage) }
            elseif ([string]$cell.verdict -ne 'disproven') { $notEligible.Add([string]$cell.stage) }
        }
    }
    $verdict = 'unresolved'
    if ($disproven.Count -gt 0) {
        $verdict = 'disproven'
    } elseif (($decision.Count -gt 0) -and ($eligible.Count -eq $decision.Count) -and ($infoViolations.Count -eq 0)) {
        # An informational violation means the device misbehaved while this arm ran (a torn record
        # in stages 1-3), which is not something an arm may be declared proven through: the cell
        # violation cannot be attributed to the candidate, so it must cap the arm instead.
        $verdict = 'proven'
    }
    return [ordered]@{
        arm            = $Arm
        verdict        = $verdict
        decisionSets   = $decision
        eligible       = $eligible.ToArray()
        notEligible    = $notEligible.ToArray()
        disproven      = $disproven.ToArray()
        infoViolations = $infoViolations.ToArray()
    }
}

function Get-B3bClaimScope {
    # The only wording a B3b verdict may be read with. Single-sourced here so the
    # verdict document, the evidence bundle README and the ADR cannot drift into a
    # stronger claim than the experiment supports: the device demonstrated
    # discrimination at stage 1 only, the payload is record-level (~1 KiB) and the
    # observation window is 2 s, and 0 losses in N=5 has a one-sided 95% upper bound
    # near p<0.45. Nothing here is a long-term guarantee.
    return ('within premise P (VirtualBox Win11B3 / NTFS on SATA-AHCI, useHostIOCache=false, ' +
        'record-level ~1 KiB payload, 2 s cut window, N=5 per cell): no counterexample observed, ' +
        'candidate mechanism demonstrably executed, device discrimination demonstrated at ' +
        'stage temp-written only; not a long-term guarantee and not a claim about other ' +
        'hypervisors, filesystems, payload sizes or the host losing power')
}

function Get-B3bVerdictDocument {
    param(
        [Parameter(Mandatory = $true)]$Cells,
        [Parameter(Mandatory = $true)][int]$RequiredN,
        [Parameter(Mandatory = $true)]$Discrimination,
        $U4 = $null,
        [string]$WriterSha256 = '',
        [string]$Premise = '',
        [string]$ClaimScope = ''
    )
    if ([string]::IsNullOrWhiteSpace($ClaimScope)) { $ClaimScope = Get-B3bClaimScope }
    $arms = New-Object System.Collections.Generic.List[object]
    foreach ($arm in @('c1', 'c2', 'c3', 'c2c3')) {
        $armCells = @($Cells | Where-Object { $_.arm -eq $arm })
        if ($armCells.Count -eq 0) { continue }
        $arms.Add((Get-B3bArmVerdict -Arm $arm -Cells $armCells))
    }
    return [ordered]@{
        schema         = 'b3b-durability-verdict'
        v              = 1
        premise        = $Premise
        claimScope     = $ClaimScope
        requiredN      = $RequiredN
        writerSha256   = $WriterSha256
        discrimination = $Discrimination
        cells          = @($Cells)
        arms           = $arms.ToArray()
        u4             = $U4
        generatedAt    = (Get-B3aHostTimestamp)
    }
}

# ---------------------------------------------------------------------------
# static module hygiene
#
# Both checks exist because B3a learned them the hard way: a dot-sourced module
# that declares a top-level param() silently rebinds same-named caller variables
# (B3a defect D1), and a dot-source that redefines a helper silently replaces it.
# Reading the AST is cheaper and safer than probing behaviour after the fact.

function Get-B3bScriptFunctionList {
    # Function names defined by a script file, plus whether it declares a
    # top-level param block (which makes it unsafe to dot-source).
    param([Parameter(Mandatory = $true)][string]$Path)
    $tokens = $null
    $errors = $null
    $ast = [System.Management.Automation.Language.Parser]::ParseFile($Path, [ref]$tokens, [ref]$errors)
    $names = New-Object System.Collections.Generic.List[string]
    foreach ($function in $ast.FindAll({ param($node) $node -is [System.Management.Automation.Language.FunctionDefinitionAst] }, $true)) {
        $names.Add($function.Name)
    }
    return [ordered]@{
        path           = $Path
        functions      = $names.ToArray()
        parameterCount = $errors.Count
        hasTopLevelParam = ($null -ne $ast.ParamBlock)
    }
}

function Test-B3bModuleHygiene {
    # Returns @{ok; collisions; topLevelParams; b3bFunctions; b3aFunctions}.
    param(
        [Parameter(Mandatory = $true)][string]$B3bRigDir,
        [Parameter(Mandatory = $true)][string]$B3aRigDir
    )
    $b3bFiles = @(
        'b3b-rig-lib.ps1',
        'Get-B3bVerdict.ps1',
        'Get-B3bReconcile.ps1'
    )
    $b3aFiles = @('b3a-rig-lib.ps1', 'b3a-journal-lib.ps1')
    $b3bFunctions = New-Object System.Collections.Generic.List[string]
    $b3aFunctions = New-Object System.Collections.Generic.List[string]
    $topLevelParams = New-Object System.Collections.Generic.List[string]
    foreach ($relative in $b3aFiles) {
        $path = Join-Path $B3aRigDir $relative
        if (-not (Test-Path -LiteralPath $path)) { continue }
        $info = Get-B3bScriptFunctionList -Path $path
        foreach ($name in $info.functions) { $b3aFunctions.Add($name) }
    }
    foreach ($relative in $b3bFiles) {
        $path = Join-Path $B3bRigDir $relative
        if (-not (Test-Path -LiteralPath $path)) { continue }
        $info = Get-B3bScriptFunctionList -Path $path
        if ($info.parameterCount -gt 0) { $topLevelParams.Add(($relative + ' (parse errors)')) }
        if ($info.hasTopLevelParam) { $topLevelParams.Add(($relative + ' (top-level param block)')) }
        foreach ($name in $info.functions) { $b3bFunctions.Add($name) }
    }
    $collisions = New-Object System.Collections.Generic.List[string]
    foreach ($name in $b3bFunctions) {
        if ($b3aFunctions -contains $name) { $collisions.Add($name) }
    }
    return [ordered]@{
        ok             = (($collisions.Count -eq 0) -and ($topLevelParams.Count -eq 0))
        collisions     = $collisions.ToArray()
        topLevelParams = $topLevelParams.ToArray()
        b3bFunctions   = $b3bFunctions.ToArray()
        b3aFunctions   = $b3aFunctions.ToArray()
    }
}

function Get-B3bReplayPairing {
    # Pair request records with completion records using the product file naming.
    #
    # The frozen B3a inventory pairs by base file name, but product records are named
    # requests\request.<id>.jsonl and completions\completion.<id>.jsonl
    # (pathForRequestID adds the prefix), so a base-name comparison can never pair
    # them and its requestsWithoutCompletion/orphanCompletions fields are always both
    # populated. That summary is therefore reported for transparency only; the verdict
    # uses this pairing, which strips the documented prefixes.
    param([Parameter(Mandatory = $true)]$Files)
    $requestIds = New-Object System.Collections.Generic.List[string]
    $completionIds = New-Object System.Collections.Generic.List[string]
    foreach ($file in @($Files)) {
        $path = [string]$file.path
        $leaf = [System.IO.Path]::GetFileNameWithoutExtension($path)
        $normalized = $path.Replace('/', '\')
        if ($normalized.StartsWith('requests\') -and $leaf.StartsWith('request.')) {
            $requestIds.Add($leaf.Substring(8))
        }
        if ($normalized.StartsWith('completions\') -and $leaf.StartsWith('completion.')) {
            $completionIds.Add($leaf.Substring(11))
        }
    }
    $missing = New-Object System.Collections.Generic.List[string]
    $orphans = New-Object System.Collections.Generic.List[string]
    foreach ($id in $requestIds) { if (-not ($completionIds -contains $id)) { $missing.Add($id) } }
    foreach ($id in $completionIds) { if (-not ($requestIds -contains $id)) { $orphans.Add($id) } }
    return [ordered]@{
        requestIds                = $requestIds.ToArray()
        completionIds             = $completionIds.ToArray()
        requestsWithoutCompletion = $missing.ToArray()
        orphanCompletions         = $orphans.ToArray()
        pairingRule               = 'product naming: request.<id>.jsonl / completion.<id>.jsonl, prefixes stripped'
    }
}

# Invoke-B3bMatrix.ps1 - orchestrates the B3b durability matrix.
#
# Two modes:
#   -Cell '<arm>:<stage>'  run N rounds of one cell, each round through
#                          Invoke-B3bRound.ps1, with the void and operational retry
#                          budgets. Round numbers continue after the highest round
#                          already present, so a session can be resumed.
#   -Finalize              rebuild every per-round result from the evidence on disk
#                          (rounds/rNN/reconcile.json), compute the discrimination
#                          gate, aggregate the three-tier verdict and write
#                          matrix.json plus verdict.json.
#
# The matrix has no separate state file on purpose: every number the verdict uses is
# recomputable from the round artifacts and the journal, so a lost or edited
# accumulator cannot change a conclusion.
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
param(
    [string]$Cell = '',
    [switch]$Finalize,
    [Parameter(Mandatory = $true)][string]$RigRoot,
    # Resolved once per matrix invocation and passed down to every round, so a long run
    # that crosses midnight cannot split its evidence across two session directories.
    [string]$SessionId = '',
    [int]$Rounds = 5,
    [int]$MaxVoidPerCell = 3,
    [int]$MaxOperationalRetry = 2,
    [string]$WriterExe = '',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep',
    [string]$CredFile = 'D:\VirtualBox VMs\Win11B3\host-only\cred.txt',
    [string]$VmName = 'Win11B3',
    [string]$VBoxManage = 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe',
    [double]$PlanCutDelayS = 2.0,
    [int]$StageDetectTimeoutS = 120,
    [int]$StageBlockS = 300,
    [string]$Premise = '',
    [int]$StartRoundFloor = 0,
    [switch]$Paired,
    [switch]$Checkpoint,
    [switch]$NoConfirm
)

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'Get-B3bVerdict.ps1')
. (Join-Path $PSScriptRoot 'Get-B3bReconcile.ps1')

$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrEmpty($WriterExe)) {
    $repoRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\..'))
    $WriterExe = Join-Path $repoRoot 'tmp\b3b\writer.exe'
}
$DateTag = $SessionId
if ([string]::IsNullOrEmpty($DateTag)) { $DateTag = Get-Date -Format 'yyyy-MM-dd' }
$SessionDir = Join-Path $RigRoot $DateTag
# One rig process per session. The lock is intentionally never released here: the next run
# checks whether the recorded pid is still alive, which keeps a crashed run recoverable while
# making a second concurrent run impossible.
$lock = Enter-B3bSessionLock -SessionDir $SessionDir -Owner 'Invoke-B3bMatrix'
if (-not $lock.ok) {
    Write-Output ('REFUSED: ' + $lock.reason)
    Write-Output 'Never run two rig processes against the same VM: they restore and start it concurrently.'
    exit 6
}
$RoundsDir = Join-Path $SessionDir 'rounds'
$JournalPath = Join-Path $SessionDir 'journal.jsonl'
$SessionLog = Join-Path $SessionDir 'session.log'

function Write-B3bMatrixLog {
    param([Parameter(Mandatory = $true)][string]$Text)
    Add-B3aLogLine -Path $SessionLog -Text ('[{0}] matrix: {1}' -f (Get-Date -Format 'HH:mm:ss'), $Text)
    $Text
}

function Get-B3bRoundNumberList {
    # Round numbers already present as round directories.
    param([Parameter(Mandatory = $true)][string]$Path)
    $numbers = New-Object System.Collections.Generic.List[int]
    if (Test-Path -LiteralPath $Path) {
        foreach ($dir in @(Get-ChildItem -LiteralPath $Path -Directory)) {
            if ($dir.Name -match '^r([0-9]{2,})$') { $numbers.Add([int]$Matches[1]) }
        }
    }
    return $numbers.ToArray()
}

function Get-B3bRoundArtifactVerdict {
    # Verify that reconcile.json still agrees with the artifacts it summarises.
    # Without this the verdict pipeline would trust its own summary file: editing
    # reconcile.json (and the inventory it points at) could flip a cell from
    # 'disproven' to 'proven' while the journal chain stays green.
    param([Parameter(Mandatory = $true)][string]$RoundDir)
    $reasons = New-Object System.Collections.Generic.List[string]
    $reconcilePath = Join-Path $RoundDir 'reconcile.json'
    if (-not (Test-Path -LiteralPath $reconcilePath)) {
        return [ordered]@{ ok = $false; reasons = @('reconcile.json missing') }
    }
    $reconcile = (Get-B3aFileTextUtf8 -Path $reconcilePath) | ConvertFrom-Json
    if ([string]$reconcile.verdict -ne 'ok') {
        # A void round makes no claim about durability, so there is no outcome to verify.
        return [ordered]@{ ok = $true; reasons = @() }
    }
    $planPath = Join-Path $RoundDir 'plan.json'
    $invPath = Join-Path $RoundDir 'inventory-1.json'
    if (-not (Test-Path -LiteralPath $planPath)) { $reasons.Add('plan.json missing') }
    if (-not (Test-Path -LiteralPath $invPath)) { $reasons.Add('inventory-1.json missing') }
    if ($reasons.Count -eq 0) {
        try {
            $derived = Get-B3bDerivedOutcome -PlanJsonPath $planPath -InventoryJsonPath $invPath
            if ([string]$derived.outcome -ne [string]$reconcile.outcome) {
                $reasons.Add(('outcome recomputed as {0} from plan.json+inventory-1.json but reconcile.json claims {1}' -f `
                        $derived.outcome, $reconcile.outcome))
            }
            if ([string]$derived.expectedDigest -ne [string]$reconcile.contentDigest) {
                $reasons.Add('contentDigest disagrees with plan.json')
            }
        } catch {
            $reasons.Add('outcome could not be recomputed: ' + $_.Exception.Message)
        }
    }
    # Every artifact hash the reconcile record claims must still be the hash of *that*
    # file. The comparison is positional on purpose: the two inventory runs are
    # byte-identical by construction, so a set-membership test would let a tampered
    # inventory-1.json pass on the strength of inventory-2.json still holding the old
    # hash.
    $invDigests = @($reconcile.inventoryDigests)
    $invPairs = @()
    if ($invDigests.Count -ge 1) { $invPairs += @{ Recorded = [string]$invDigests[0]; Name = 'inventory-1.json' } }
    if ($invDigests.Count -ge 2) { $invPairs += @{ Recorded = [string]$invDigests[1]; Name = 'inventory-2.json' } }
    foreach ($pair in $invPairs) {
        $path = Join-Path $RoundDir $pair.Name
        if (-not (Test-Path -LiteralPath $path)) { $reasons.Add(($pair.Name + ' missing')); continue }
        $live = (Get-B3aSha256Hex -Path $path).ToLowerInvariant()
        if ($live -ne $pair.Recorded.ToLowerInvariant()) {
            $reasons.Add(($pair.Name + ' was modified after the round (recorded ' + $pair.Recorded.Substring(0, 12) + ', live ' + $live.Substring(0, 12) + ')'))
        }
    }
    $digestPairs = @(
        @{ Field = 'auditDigest'; Name = 'audit.txt' },
        @{ Field = 'markerPreCutDigest'; Name = 'marker-transcript.txt' }
    )
    # marker.txt (the copy pulled back after the cut) is deliberately NOT part of this
    # gate: a power cut can legitimately truncate it, so treating that as tampering would
    # refuse honest rounds. The pre-cut archive is the authoritative stage proof, and the
    # difference between the two is already recorded per round as markerPostCutDiffers.
    foreach ($pair in $digestPairs) {
        if ($null -eq $reconcile.PSObject.Properties[$pair.Field]) { continue }
        $recorded = [string]$reconcile.($pair.Field)
        if ([string]::IsNullOrEmpty($recorded)) { continue }
        $path = Join-Path $RoundDir $pair.Name
        if (-not (Test-Path -LiteralPath $path)) { $reasons.Add(($pair.Name + ' missing')); continue }
        $live = (Get-B3aSha256Hex -Path $path).ToLowerInvariant()
        if ($live -ne $recorded.ToLowerInvariant()) {
            $reasons.Add(($pair.Name + ' was modified after the round'))
        }
    }
    return [ordered]@{ ok = ($reasons.Count -eq 0); reasons = $reasons.ToArray() }
}

function Get-B3bJournalReconcileMap {
    # round number -> the reconcile journal record for that round. Used to cross-check
    # the round artifacts against the tamper-evident journal: the journal line is
    # hash-chained, so agreement between the two is evidence the artifacts were not
    # rewritten after the round finished.
    param([Parameter(Mandatory = $true)][string]$JournalPath)
    $map = @{}
    if (-not (Test-Path -LiteralPath $JournalPath)) { return $map }
    foreach ($record in @(Get-B3aJournalRecord -JournalPath $JournalPath)) {
        if ([string]$record.op -ne 'reconcile') { continue }
        $map[[int]$record.round] = $record
    }
    return $map
}

function Get-B3bJournalCrossCheck {
    param(
        [Parameter(Mandatory = $true)][string]$RoundDir,
        $JournalRecord
    )
    $reasons = New-Object System.Collections.Generic.List[string]
    $reconcilePath = Join-Path $RoundDir 'reconcile.json'
    if (-not (Test-Path -LiteralPath $reconcilePath)) { return @('reconcile.json missing') }
    $reconcile = (Get-B3aFileTextUtf8 -Path $reconcilePath) | ConvertFrom-Json
    if ($null -eq $JournalRecord) {
        $reasons.Add('no reconcile journal record for this round')
        return $reasons.ToArray()
    }
    if ([string]$JournalRecord.contentDigest -ne [string]$reconcile.contentDigest) {
        $reasons.Add('journal contentDigest disagrees with reconcile.json')
    }
    if ([string]$JournalRecord.candidate -ne [string]$reconcile.candidate) {
        $reasons.Add('journal candidate disagrees with reconcile.json')
    }
    if ([string]$JournalRecord.stage -ne [string]$reconcile.stage) {
        $reasons.Add('journal stage disagrees with reconcile.json')
    }
    $invNames = @('inventory-1.json', 'inventory-2.json')
    if ([string]$reconcile.verdict -eq 'ok') {
        # Positional comparison for the same reason as in the artifact verification:
        # two identical inventory runs must not be able to vouch for each other.
        $invDigests = @($reconcile.inventoryDigests)
        for ($i = 0; $i -lt $invDigests.Count; $i++) {
            if ([string]::IsNullOrEmpty([string]$invDigests[$i])) { continue }
            $path = Join-Path $RoundDir $invNames[$i]
            if (-not (Test-Path -LiteralPath $path)) { $reasons.Add(($invNames[$i] + ' missing')); continue }
            if ((Get-B3aSha256Hex -Path $path).ToLowerInvariant() -ne ([string]$invDigests[$i]).ToLowerInvariant()) {
                $reasons.Add(($invNames[$i] + ' no longer matches the digest recorded in reconcile.json'))
            }
        }
    }
    # The recorded verdict and void reason must agree with reconcile.json for **both**
    # kinds of round: otherwise a measured round edited on the reconcile side into a
    # 'void' would drop out of the valid set without tripping this gate at all.
    $extra = $null
    if ($null -ne $JournalRecord.PSObject.Properties['extra']) { $extra = $JournalRecord.extra }
    if ($null -ne $extra) {
        if ($null -ne $extra.PSObject.Properties['verdict']) {
            if ([string]$extra.verdict -ne [string]$reconcile.verdict) {
                $reasons.Add(('journal verdict ({0}) disagrees with reconcile.json ({1})' -f $extra.verdict, $reconcile.verdict))
            }
        }
        if ($null -ne $extra.PSObject.Properties['voidReason']) {
            if ([string]$extra.voidReason -ne [string]$reconcile.voidReason) {
                $reasons.Add('journal voidReason disagrees with reconcile.json')
            }
        }
        if (($null -ne $extra.PSObject.Properties['outcome']) -and ([string]$extra.outcome -ne [string]$reconcile.outcome)) {
            $reasons.Add(('journal outcome ({0}) disagrees with reconcile.json ({1})' -f $extra.outcome, $reconcile.outcome))
        }
        if ($null -ne $extra.PSObject.Properties['inventoryDigests']) {
            $extraDigests = @($extra.inventoryDigests)
            for ($i = 0; $i -lt $extraDigests.Count; $i++) {
                if ([string]::IsNullOrEmpty([string]$extraDigests[$i])) { continue }
                $path = Join-Path $RoundDir $invNames[$i]
                if (-not (Test-Path -LiteralPath $path)) { continue }
                if ((Get-B3aSha256Hex -Path $path).ToLowerInvariant() -ne ([string]$extraDigests[$i]).ToLowerInvariant()) {
                    $reasons.Add(($invNames[$i] + ' no longer matches the digest the journal recorded for this round'))
                }
            }
        }
    }
    return $reasons.ToArray()
}

function Get-B3bRoundResult {
    # One round's result rebuilt from its reconcile.json **and** verified against the
    # artifacts that produced it and against the journal record for that round. Any
    # missing or disagreeing artifact makes the round invalid rather than silently
    # passing a gate.
    param(
        [Parameter(Mandatory = $true)][string]$RoundDir,
        $JournalRecords = $null,
        [bool]$ChainOk = $true
    )
    $reconcilePath = Join-Path $RoundDir 'reconcile.json'
    # A round without its artifacts is invalid, and it carries unknown labels so it can
    # never join a real cell.
    $result = New-B3bRoundResult -Round 0 -Arm 'unknown' -Stage 'unknown'
    if (-not (Test-Path -LiteralPath $reconcilePath)) { return $result }
    $reconcile = (Get-B3aFileTextUtf8 -Path $reconcilePath) | ConvertFrom-Json
    $round = 0
    if ($null -ne $reconcile.PSObject.Properties['round']) { $round = [int]$reconcile.round }
    $arm = ''
    if ($null -ne $reconcile.PSObject.Properties['candidate']) { $arm = [string]$reconcile.candidate }
    $stage = ''
    if ($null -ne $reconcile.PSObject.Properties['stage']) { $stage = [string]$reconcile.stage }
    $paired = $false
    if ($null -ne $reconcile.PSObject.Properties['paired']) { $paired = [bool]$reconcile.paired }
    $targetKind = ''
    if ($null -ne $reconcile.PSObject.Properties['targetKind']) { $targetKind = [string]$reconcile.targetKind }
    $stageHit = $false
    if ($null -ne $reconcile.PSObject.Properties['stageHit']) { $stageHit = [bool]$reconcile.stageHit }
    $record = $null
    if ($null -ne $JournalRecords) { $record = $JournalRecords[$round] }
    if ([string]$reconcile.verdict -ne 'ok') {
        # A void round is not a measurement, but it must still be a *recorded* void: it
        # has to agree with the journal record for that round and its reason has to be one
        # the driver can produce. Without this, editing a measured round's reconcile (and
        # its journal line) into a void would drop a real counterexample from the valid set
        # without tripping the integrity gate.
        $voidReasons = New-Object System.Collections.Generic.List[string]
        foreach ($reason in @(Get-B3bJournalCrossCheck -RoundDir $RoundDir -JournalRecord $record)) {
            $voidReasons.Add([string]$reason)
        }
        if (-not (Test-B3bVoidReasonKnown -VoidReason ([string]$reconcile.voidReason))) {
            $voidReasons.Add('the void reason is not one the round driver can produce: ' + [string]$reconcile.voidReason)
        }
        $voidReason = [string]$reconcile.voidReason
        $voidIntegrityOk = ($voidReasons.Count -eq 0) -and $ChainOk
        if (-not $voidIntegrityOk) {
            $voidReason = 'evidence-mismatch: ' + (($voidReasons.ToArray()) -join '; ')
        }
        $voidResult = New-B3bRoundResult -Round $round -Arm $arm -Stage $stage -Valid $false -VoidReason $voidReason
        $voidResult['paired'] = $paired
        $voidResult['targetKind'] = $targetKind
        $voidResult['sessionBlock'] = [string]$reconcile.sessionBlock
        $voidResult['integrityOk'] = $voidIntegrityOk
        return $voidResult
    }
    $auditOk = ($null -ne $reconcile.PSObject.Properties['auditVerdict']) -and ([string]$reconcile.auditVerdict -eq 'HARD-POWER-LOSS')
    $inventoryOk = ($null -ne $reconcile.PSObject.Properties['inventoryIdentical']) -and ([bool]$reconcile.inventoryIdentical)
    # Integrity gate: the round may only enter a count when its own artifacts agree
    # with reconcile.json, with the journal record and with the journal chain.
    $verification = Get-B3bRoundArtifactVerdict -RoundDir $RoundDir
    $crossReasons = @()
    if ($verification.ok) { $crossReasons = @(Get-B3bJournalCrossCheck -RoundDir $RoundDir -JournalRecord $record) }
    $integrityReasons = New-Object System.Collections.Generic.List[string]
    foreach ($reason in @($verification.reasons)) { $integrityReasons.Add([string]$reason) }
    foreach ($reason in @($crossReasons)) { $integrityReasons.Add([string]$reason) }
    $integrityOk = ($integrityReasons.Count -eq 0) -and $ChainOk
    if (-not $integrityOk) {
        $integrityResult = New-B3bRoundResult -Round $round -Arm $arm -Stage $stage -Valid $false `
            -VoidReason ('evidence-mismatch: ' + (($integrityReasons.ToArray()) -join '; '))
        $integrityResult['paired'] = $paired
        $integrityResult['targetKind'] = $targetKind
        $integrityResult['integrityOk'] = $false
        return $integrityResult
    }
    $final = New-B3bRoundResult -Round $round -Arm $arm -Stage $stage -Valid $true `
            -Outcome ([string]$reconcile.outcome) `
            -MechanismOk ([bool]$reconcile.mechanismOk) `
            -StageHit $stageHit `
            -AuditOk $auditOk `
            -InventoryRepeatableOk $inventoryOk `
            -JournalChainOk $true `
            -ExpectedDigest ([string]$reconcile.contentDigest) `
            -PlanDigest ([string]$reconcile.contentDigest) `
            -MarkerStages ([string]$reconcile.markerStages)
    $final['paired'] = $paired
    $final['targetKind'] = $targetKind
    $final['sessionBlock'] = [string]$reconcile.sessionBlock
    $final['integrityOk'] = $true
    return $final
}

function Get-B3bAllRoundResult {
    param(
        [Parameter(Mandatory = $true)][string]$Path,
        [string]$JournalPath = ''
    )
    $results = New-Object System.Collections.Generic.List[object]
    if (-not (Test-Path -LiteralPath $Path)) { return $results.ToArray() }
    $records = $null
    $chainOk = $true
    if (-not [string]::IsNullOrEmpty($JournalPath)) {
        $chain = Test-B3aJournalChain -JournalPath $JournalPath
        $chainOk = [bool]$chain.ok
        $records = Get-B3bJournalReconcileMap -JournalPath $JournalPath
    }
    foreach ($dir in (@(Get-ChildItem -LiteralPath $Path -Directory) | Sort-Object Name)) {
        $results.Add((Get-B3bRoundResult -RoundDir $dir.FullName -JournalRecords $records -ChainOk $chainOk))
    }
    return $results.ToArray()
}

if ($Cell) {
    # Parsed by the shared resolver, so the command line and the self-test cannot disagree
    # about what a legal cell looks like (the control arm name contains a hyphen).
    $cellSpec = $null
    try {
        $cellSpec = Resolve-B3bCellSpec -Cell $Cell -ArmList (Get-B3bArmList) -StageList (Get-B3bPublishStageList)
    } catch {
        throw ('Invoke-B3bMatrix: ' + $_.Exception.Message)
    }
    $arm = $cellSpec.arm
    $stage = $cellSpec.stage
    $paired = [bool]$Paired
    # The two sub-matrices have reserved round ranges. Enforcing them here means a cell
    # cannot be run into the wrong range and be silently reclassified at finalize.
    if ($paired -and ($StartRoundFloor -lt 500)) { $StartRoundFloor = 500 }
    if ((-not $paired) -and ($StartRoundFloor -ge 500)) {
        throw 'Invoke-B3bMatrix: the main matrix must not write into the ordering range (round >= 500)'
    }
    $all = Get-B3bAllRoundResult -Path $RoundsDir -JournalPath $JournalPath
    $existing = @($all | Where-Object { $_.arm -eq $arm -and $_.stage -eq $stage -and ([bool]$_.paired -eq $paired) })
    $validCount = @($existing | Where-Object { $_.valid }).Count
    $voidCount = @($existing | Where-Object { -not $_.valid }).Count
    [void](Write-B3bMatrixLog ('cell {0} paired={1} before: valid={2} void={3} (target {4})' -f $Cell, $paired, $validCount, $voidCount, $Rounds))
    $numbers = Get-B3bRoundNumberList -Path $RoundsDir
    $next = 0
    if (@($numbers).Count -gt 0) { $next = ([int[]]$numbers | Measure-Object -Maximum).Maximum + 1 }
    if ($next -lt $StartRoundFloor) { $next = $StartRoundFloor }
    if ((-not $paired) -and ($next -ge 500)) {
        throw ('Invoke-B3bMatrix: the main matrix has reached round {0}; run the ordering sub-matrix with -Paired -StartRoundFloor 500 instead' -f $next)
    }
    $operational = 0
    while ($validCount -lt $Rounds) {
        if ($voidCount -ge $MaxVoidPerCell) {
            [void](Write-B3bMatrixLog ('cell {0} blocked: void budget exhausted ({1})' -f $Cell, $voidCount))
            exit 3
        }
        $arguments = @(
            '-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', (Get-B3aQuoted (Join-Path $PSScriptRoot 'Invoke-B3bRound.ps1')),
            '-Round', ([string]$next), '-Candidate', $arm, '-Stage', $stage,
            '-SessionId', $DateTag,
            '-RigRoot', (Get-B3aQuoted $RigRoot), '-WriterExe', (Get-B3aQuoted $WriterExe),
            '-GuestPrepDir', (Get-B3aQuoted $GuestPrepDir), '-CredFile', (Get-B3aQuoted $CredFile),
            '-VmName', (Get-B3aQuoted $VmName), '-VBoxManage', (Get-B3aQuoted $VBoxManage),
            '-PlanCutDelayS', ([string]$PlanCutDelayS),
            '-StageDetectTimeoutS', ([string]$StageDetectTimeoutS), '-StageBlockS', ([string]$StageBlockS)
        )
        if ($paired) { $arguments += '-Paired' }
        if ($NoConfirm) { $arguments += '-NoConfirm' }
        $run = Invoke-B3aProcess -FileName 'powershell' -Arguments ($arguments -join ' ') -TimeoutS 3600
        $code = $run.exitCode
        [void](Write-B3bMatrixLog ('round {0} cell={1} exit={2}' -f $next, $Cell, $code))
        if (-not (@(0, 3, 5, 6, 7) -contains $code)) {
            # An unexpected exit means the round driver threw. Its stderr carries the message and
            # would otherwise be lost, leaving a bare exit code in the session log - which is
            # exactly how the first post-cut failure of this rig presented itself.
            [void](Write-B3bMatrixLog ('round {0} unexpected exit {1}; stderr: {2}' -f $next, $code, (Get-B3aBoundedText -Text $run.stderrText)))
        }
        if ($code -eq 0) {
            $validCount++
            $operational = 0
            $next++
            continue
        }
        if ($code -eq 3) {
            $voidCount++
            $next++
            continue
        }
        if ($code -eq 7) {
            [void](Write-B3bMatrixLog 'operator declined the cut; stopping the cell')
            exit 7
        }
        if ($code -eq 6) {
            [void](Write-B3bMatrixLog 'session blocker reported by the round; stopping')
            exit 6
        }
        $operational++
        [void](Write-B3bMatrixLog ('operational failure (exit {0}), retry {1}/{2}' -f $code, $operational, $MaxOperationalRetry))
        if ($operational -gt $MaxOperationalRetry) {
            [void](Write-B3bMatrixLog 'operational retry budget exhausted; stopping the cell')
            exit 5
        }
    }
    [void](Write-B3bMatrixLog ('cell {0} complete: valid={1}' -f $Cell, $validCount))
    exit 0
}

if ($Checkpoint) {
    # Mid-run decision point. The discrimination gate is the one measurement that decides
    # whether the remaining cells can say anything at all, so it is evaluated as soon as
    # the baseline-early and control cells are complete instead of after the whole budget
    # has been spent: a gate that cannot pass would otherwise burn hours and end in
    # 'unresolved' for a reason nobody chose.
    $all = Get-B3bAllRoundResult -Path $RoundsDir -JournalPath $JournalPath
    $valid = @($all | Where-Object { $_.valid -and (-not [bool]$_.paired) })
    $c1Early = @($valid | Where-Object { $_.arm -eq 'c1' -and $_.stage -eq (Get-B3bPublishStageList)[0] })
    $positive = @($valid | Where-Object { $_.arm -eq 'positive-control' })
    $gate = Get-B3bDiscriminationGate -NegativeValidRounds $c1Early.Count `
        -NegativeLosses (@($c1Early | Where-Object { $_.outcome -ne 'survived' }).Count) `
        -PositiveValidRounds $positive.Count `
        -PositiveLosses (@($positive | Where-Object { $_.outcome -ne 'survived' }).Count) -RequiredN 5
    [void](Write-B3bMatrixLog ('checkpoint: gate ok={0} baseline-early valid={1} losses={2} control valid={3} losses={4}' -f `
            $gate.ok, $gate.negativeValidRounds, $gate.negativeLosses, $gate.positiveValidRounds, $gate.positiveLosses))
    if (-not $gate.ok) {
        Write-Output 'CHECKPOINT FAILED: the device did not demonstrate discrimination power.'
        Write-Output 'No cell of this matrix may be reported as proven. Operator decision required:'
        Write-Output '  - raise N or widen the cut window for the baseline-early cell and re-run the checkpoint, or'
        Write-Output '  - record the explicit decision to continue (the run can then only end in unresolved), or'
        Write-Output '  - abort and fix the device.'
        exit 4
    }
    Write-Output 'CHECKPOINT PASSED: discrimination power demonstrated; the matrix may proceed.'
    exit 0
}

if ($Finalize) {
    if ([string]::IsNullOrWhiteSpace($Premise)) {
        throw 'Invoke-B3bMatrix -Finalize: -Premise is required; a verdict without its premise cannot be read correctly'
    }
    $all = Get-B3bAllRoundResult -Path $RoundsDir -JournalPath $JournalPath
    # A finalize that finds no rounds at all must refuse instead of writing an empty verdict.
    # The session directory is derived from -SessionId (falling back to today's date), and the
    # pin that fixes the session id lives *inside* that directory, so a missing or stale
    # -SessionId resolves to a different, empty directory - which used to produce a verdict
    # built from zero rounds and exit 0. That silently discards a whole session's evidence
    # and is indistinguishable from an honest 'unresolved' at read time.
    if ($all.Count -eq 0) {
        Write-Output ('REFUSED: no round evidence under the resolved session directory: ' + $SessionDir)
        Write-Output 'Pass the same -SessionId that produced the rounds; a session id is never inferred from the evidence.'
        exit 4
    }
    # Ordering rounds are identified by the pairing flag each round recorded, not by a
    # round-number convention: a number-only split would let a mis-numbered round (or a
    # resume that walked the main matrix into the reserved range) move a real
    # counterexample out of its arm's verdict.
    $u4All = @($all | Where-Object { [bool]$_.paired })
    $matrixAll = @($all | Where-Object { -not [bool]$_.paired })
    $numberedWrong = @($all | Where-Object { ([int]$_.round -ge 500) -ne [bool]$_.paired })
    if ($numberedWrong.Count -gt 0) {
        Write-Output ('REFUSED: {0} round(s) disagree with the numbering convention (paired rounds must be >= 500, matrix rounds < 500): {1}' -f `
                $numberedWrong.Count, (($numberedWrong | ForEach-Object { 'r' + [int]$_.round }) -join ','))
        exit 4
    }
    $kindWrong = @($u4All | Where-Object { [string]$_.targetKind -ne 'completion' })
    if ($kindWrong.Count -gt 0) {
        Write-Output ('REFUSED: {0} ordering round(s) were not measured on a completion record' -f $kindWrong.Count)
        exit 4
    }
    $integrity = @($all | Where-Object { [string]$_.voidReason -like 'evidence-mismatch*' })
    if ($integrity.Count -gt 0) {
        Write-Output ('REFUSED: the evidence is not self-consistent for {0} round(s); no verdict may be derived from it.' -f $integrity.Count)
        foreach ($round in $integrity) { Write-Output ('  r{0}: {1}' -f [int]$round.round, $round.voidReason) }
        exit 4
    }
    $blocked = @($all | Where-Object { -not [string]::IsNullOrEmpty([string]$_.sessionBlock) })
    if ($blocked.Count -gt 0) {
        Write-Output ('REFUSED: a session blocker was recorded: {0}' -f (($blocked | ForEach-Object { 'r' + [int]$_.round + ':' + $_.sessionBlock }) -join ','))
        exit 6
    }
    $valid = @($matrixAll | Where-Object { $_.valid })
    $void = @($matrixAll | Where-Object { -not $_.valid })
    $c1Early = @($valid | Where-Object { $_.arm -eq 'c1' -and $_.stage -eq (Get-B3bPublishStageList)[0] })
    $positive = @($valid | Where-Object { $_.arm -eq 'positive-control' })
    $discrimination = Get-B3bDiscriminationGate -NegativeValidRounds $c1Early.Count `
        -NegativeLosses (@($c1Early | Where-Object { $_.outcome -ne 'survived' }).Count) `
        -PositiveValidRounds $positive.Count `
        -PositiveLosses (@($positive | Where-Object { $_.outcome -ne 'survived' }).Count) `
        -RequiredN 5
    $u4 = $null
    if ($u4All.Count -gt 0) {
        # The ordering question is answered from the inventory the round pinned, not from the
        # summary field: reconcile.orphanCompletions is a claim, and a single-field edit of it
        # would otherwise turn a real ordering falsification into 'unresolved' with every gate
        # still green - the same threat family the non-void path already closes.
        $u4Rounds = New-Object System.Collections.Generic.List[object]
        $orphanMismatch = New-Object System.Collections.Generic.List[string]
        foreach ($round in $u4All) {
            $roundDirName = 'r{0:D2}' -f [int]$round.round
            $roundDir = Join-Path $RoundsDir $roundDirName
            $reconcilePath = Join-Path $roundDir 'reconcile.json'
            $orphan = $false
            if (Test-Path -LiteralPath $reconcilePath) {
                $reconcile = (Get-B3aFileTextUtf8 -Path $reconcilePath) | ConvertFrom-Json
                $claimedOrphans = @()
                if ($null -ne $reconcile.PSObject.Properties['orphanCompletions']) { $claimedOrphans = @($reconcile.orphanCompletions) }
                if ([string]$reconcile.verdict -ne 'ok') {
                    # A void round writes no inventory - every void path exits before the inventory is
                    # taken - and its orphan claim is empty by construction. Its tamper surface is the
                    # void vocabulary plus the journal cross-check in Get-B3bRoundResult, so there is
                    # nothing to re-derive here. Demanding an inventory anyway would refuse an honest
                    # session whenever a U4 round was voided inside its budget, which is a normal
                    # outcome (up to three voids per cell are allowed).
                    $orphan = $false
                } else {
                    $derivedOrphans = $null
                    $inventoryPath = Join-Path $roundDir 'inventory-1.json'
                    if (Test-Path -LiteralPath $inventoryPath) {
                        $inventory = (Get-B3aFileTextUtf8 -Path $inventoryPath) | ConvertFrom-Json
                        $probeFiles = @()
                        if ($inventory.PSObject.Properties.Name -contains 'probeEntries') {
                            if ($null -ne $inventory.probeEntries.PSObject.Properties['files']) { $probeFiles = @($inventory.probeEntries.files) }
                        }
                        $derived = Get-B3bReplayPairing -Files $probeFiles
                        $derivedOrphans = @($derived.orphanCompletions)
                    }
                    if ($null -eq $derivedOrphans) {
                        $orphanMismatch.Add(('r{0}: inventory-1.json missing, the orphan claim cannot be re-derived' -f [int]$round.round))
                    } elseif ($derivedOrphans.Count -ne $claimedOrphans.Count) {
                        $orphanMismatch.Add(('r{0}: the inventory re-derives {1} orphan completion(s) but reconcile.json claims {2}' -f `
                                [int]$round.round, $derivedOrphans.Count, $claimedOrphans.Count))
                    } else {
                        $orphan = ($claimedOrphans.Count -gt 0)
                    }
                }
            }
            $entry = New-B3bRoundResult -Round ([int]$round.round) -Arm ([string]$round.arm) -Stage ([string]$round.stage) `
                -Valid ([bool]$round.valid) -VoidReason ([string]$round.voidReason) -Outcome ([string]$round.outcome) `
                -MechanismOk ([bool]$round.mechanismOk) -AuditOk ([bool]$round.auditOk) `
                -InventoryRepeatableOk ([bool]$round.inventoryRepeatableOk) -JournalChainOk ([bool]$round.journalChainOk)
            $entry['orphanCompletion'] = $orphan
            $u4Rounds.Add($entry)
        }
        if ($orphanMismatch.Count -gt 0) {
            Write-Output ('REFUSED: the ordering evidence is not self-consistent for {0} round(s).' -f $orphanMismatch.Count)
            foreach ($reason in $orphanMismatch) { Write-Output ('  ' + $reason) }
            exit 4
        }
        $u4 = Get-B3bU4Verdict -Rounds $u4Rounds.ToArray() -RequiredN 3
    }
    $writerSha = ''
    if (Test-Path -LiteralPath $WriterExe) { $writerSha = Get-B3aSha256Hex -Path $WriterExe }
    $chain = Test-B3aJournalChain -JournalPath $JournalPath
    $verdict = Invoke-B3bVerdict -Rounds $valid -RequiredN 5 -Discrimination $discrimination -U4 $u4 `
        -WriterSha256 $writerSha -Premise $Premise
    $matrix = [ordered]@{
        schema                  = 'b3b-matrix'
        v                       = 1
        sessionDir              = $SessionDir
        premise                 = $Premise
        requiredN               = 5
        writerSha256            = $writerSha
        journalChainOk          = $chain.ok
        journalRecordCount      = $chain.count
        validRoundCount         = $valid.Count
        voidRoundCount          = $void.Count
        voidReasons             = @($void | ForEach-Object { $_.voidReason } | Sort-Object -Unique)
        discrimination          = $discrimination
        rounds                  = @($all)
        u4                      = $u4
    }
    $matrixPath = Join-Path $SessionDir 'matrix.json'
    Write-B3aTextFile -Path $matrixPath -Text ($matrix | ConvertTo-Json -Depth 20)
    $verdictPath = Join-Path $SessionDir 'verdict.json'
    Write-B3aTextFile -Path $verdictPath -Text ($verdict | ConvertTo-Json -Depth 20)
    [void](Write-B3bMatrixLog ('finalize: valid={0} void={1} discrimination={2} journalChain={3}' -f `
            $valid.Count, $void.Count, $discrimination.ok, $chain.ok))
    foreach ($armResult in @($verdict.arms)) {
        '{0} => {1}' -f $armResult.arm, $armResult.verdict
    }
    'matrix: ' + $matrixPath
    'verdict: ' + $verdictPath
    exit 0
}

'usage: Invoke-B3bMatrix.ps1 -Cell <arm:stage> -RigRoot <dir> [-Rounds N] [-NoConfirm]'
'       Invoke-B3bMatrix.ps1 -Finalize -RigRoot <dir> [-Premise <text>]'
''
'cells: <arm> is c1|c2|c3|c2c3|positive-control; <stage> is one of the five publish stages.'
'The ordering sub-matrix uses the same cell names with -Paired and -StartRoundFloor 500;'
'its rounds are kept in the 500+ range so the finalize step can tell the two questions apart.'
'exit 1'
exit 1

# Get-B3bReconcile.ps1 - outcome classification MODULE for one B3b round.
#
# Why this exists instead of reusing the B3a classifier: the B3a classifier looks the
# target up in the round root, because its measured file lived there. B3b measures a
# real product record, which lives in the replay store layout, so the target has to
# be found in the inventory's probe section. Everything else is shared with B3a -
# the same Get-B3aOutcome classifier, the same zero-digest rule, the same journal
# record shape and the same ACK discipline.
#
# No top-level param() block: this file is dot-sourced by the round driver and by
# the rig self-test, and a top-level param() would silently rebind same-named
# caller variables (B3a defect D1).
#
# Powershell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')

function Get-B3bMechanismToken {
    # The syscall token a candidate is required to have executed, as reported in the
    # writer's marker trace. Without this the "nothing was lost" of a cell could
    # come from a candidate that silently degraded to the baseline.
    param([Parameter(Mandatory = $true)][string]$Arm)
    switch ($Arm) {
        'c1' { return , @('os.Link') }
        'c2' { return , @('MoveFileExW(WRITE_THROUGH,no-replace)') }
        'c3' { return , @('FlushFileBuffers(dir)') }
        'c2c3' { return , @('MoveFileExW(WRITE_THROUGH,no-replace)', 'FlushFileBuffers(dir)') }
        # The control claims nothing, so it requires no token: whichever flush it
        # managed (volume handle, or the record itself when the volume handle is
        # unavailable at that privilege level) is recorded in the trace and read by
        # the report, but it must not be able to hold an arm back.
        'positive-control' { return , @() }
        default { throw ('Get-B3bMechanismToken: unknown arm ' + $Arm) }
    }
}

function Get-B3bMechanismForbiddenToken {
    # Tokens that must NOT appear for the arm: the baseline must not secretly run a
    # candidate primitive, or the cell would be measuring the wrong code path.
    param([Parameter(Mandatory = $true)][string]$Arm)
    switch ($Arm) {
        'c1' { return , @('MoveFileExW(WRITE_THROUGH,no-replace)', 'FlushFileBuffers(dir)') }
        'c2' { return , @('FlushFileBuffers(dir)') }
        default { return , @() }
    }
}

function Get-B3bMechanismVerdict {
    # @{ok; missing; forbidden; trace}
    param(
        [Parameter(Mandatory = $true)][string]$Arm,
        [Parameter(Mandatory = $true)]$Trace
    )
    $traceList = @($Trace | ForEach-Object { [string]$_ })
    $missing = New-Object System.Collections.Generic.List[string]
    foreach ($token in (Get-B3bMechanismToken -Arm $Arm)) {
        if (-not ($traceList -contains $token)) { $missing.Add($token) }
    }
    $forbidden = New-Object System.Collections.Generic.List[string]
    foreach ($token in (Get-B3bMechanismForbiddenToken -Arm $Arm)) {
        if ($traceList -contains $token) { $forbidden.Add($token) }
    }
    return [ordered]@{
        ok        = (($missing.Count -eq 0) -and ($forbidden.Count -eq 0))
        missing   = $missing.ToArray()
        forbidden = $forbidden.ToArray()
        trace     = $traceList
    }
}

function Get-B3bStageHit {
    # Was the target stage reached, and did the writer stop there? A stage later than
    # the target means the writer escaped the stage the host was told to cut in: that
    # is a protocol defect, not a void round.
    param(
        [Parameter(Mandatory = $true)][string]$Target,
        [Parameter(Mandatory = $true)]$Transcript,
        [string]$TargetPath = ''
    )
    $entries = @($Transcript.entries)
    if ([string]::IsNullOrEmpty($TargetPath)) {
        $stages = Get-B3bMarkerStageList -Entries $entries
    } else {
        # Only the measured phase counts: the paired mode publishes the request record
        # first and the completion record second, and both phases emit the same stage
        # names, so without the path filter the request phase would satisfy "the target
        # stage was reached" and the host would cut during the wrong publication.
        $stages = Get-B3bMeasuredPhaseStageList -Entries $entries -TargetPath $TargetPath
    }
    $hit = Get-B3bTargetStageSeen -Stages $stages -Target $Target
    $escaped = $false
    if (@($stages).Count -gt 0) { $escaped = Test-B3bStageEscape -Stages $stages -Target $Target }
    # The last matching entry wins: it belongs to the measured phase.
    $targetTrace = @()
    foreach ($entry in $entries) {
        if ([string]$entry.stage -ne $Target) { continue }
        if (-not [string]::IsNullOrEmpty($TargetPath)) {
            if ($null -eq $entry.PSObject.Properties['path']) { continue }
            if ([System.IO.Path]::GetFileName([string]$entry.path) -ne [System.IO.Path]::GetFileName($TargetPath)) { continue }
        }
        if ($null -ne $entry.PSObject.Properties['trace']) { $targetTrace = @($entry.trace) }
    }
    return [ordered]@{
        hit          = $hit
        escaped      = $escaped
        stages       = $stages
        unparseable  = @($Transcript.unparseable).Count
        targetTrace  = $targetTrace
    }
}

function Get-B3bDerivedOutcome {
    # Recompute one round's outcome from the two artifacts that recorded the world
    # state: plan.json (the writer's pre-commitment) and inventory-1.json (the
    # post-restart view of the replay store).
    #
    # Why this is a function: the round reconciliation and the matrix's later
    # verification must derive the outcome with the *same* code. A verdict pipeline
    # that trusts reconcile.json as its own source of truth can be flipped from
    # 'disproven' to 'proven' by editing two JSON files while the journal chain stays
    # green, which would make the evidence chain worthless.
    param(
        [Parameter(Mandatory = $true)][string]$PlanJsonPath,
        [Parameter(Mandatory = $true)][string]$InventoryJsonPath
    )
    if (-not (Test-Path -LiteralPath $PlanJsonPath)) { throw 'plan.json missing' }
    if (-not (Test-Path -LiteralPath $InventoryJsonPath)) { throw 'inventory-1.json missing' }
    $plan = (Get-B3aFileTextUtf8 -Path $PlanJsonPath) | ConvertFrom-Json
    $inv1 = (Get-B3aFileTextUtf8 -Path $InventoryJsonPath) | ConvertFrom-Json

    # The measured record lives in the replay store layout, so the target is looked up
    # in the probe section: the round root holds no product record.
    $targetLeaf = [System.IO.Path]::GetFileName([string]$plan.targetPath)
    $probeState = ''
    if ($inv1.PSObject.Properties.Name -contains 'probeState') { $probeState = [string]$inv1.probeState }
    $probeFiles = @()
    if ($inv1.PSObject.Properties.Name -contains 'probeEntries') {
        if ($null -ne $inv1.probeEntries.PSObject.Properties['files']) { $probeFiles = @($inv1.probeEntries.files) }
    }
    $entry = $null
    foreach ($f in $probeFiles) {
        if ([System.IO.Path]::GetFileName([string]$f.path) -eq $targetLeaf) { $entry = $f; break }
    }
    $present = ($null -ne $entry)
    $actualSize = 0
    $actualDigest = ''
    if ($present) {
        $actualSize = [long]$entry.size
        $actualDigest = [string]$entry.sha256
    }
    $expectedSize = [long]$plan.contentSize
    $expectedDigest = [string]$plan.contentDigest
    $outcome = Get-B3aOutcome -Present $present -ActualSize $actualSize -ActualDigest $actualDigest `
        -ExpectedSize $expectedSize -ExpectedDigest $expectedDigest -ZeroDigest (Get-B3aZeroDigest -Size $expectedSize)
    $frozenSummary = $null
    if ($inv1.PSObject.Properties.Name -contains 'probeEntries') {
        if ($null -ne $inv1.probeEntries.PSObject.Properties['summary']) {
            $frozenSummary = $inv1.probeEntries.summary
        }
    }
    $dirty = $null
    if ($inv1.PSObject.Properties.Name -contains 'dirty') { $dirty = $inv1.dirty }
    return [ordered]@{
        outcome           = $outcome
        present           = $present
        actualSize        = $actualSize
        actualDigest      = $actualDigest
        expectedSize      = $expectedSize
        expectedDigest    = $expectedDigest
        planStage         = [string]$plan.stage
        planTargetPath    = [string]$plan.targetPath
        probeState        = $probeState
        probeFiles        = $probeFiles
        frozenProbeSummary = $frozenSummary
        dirty             = $dirty
    }
}

function Invoke-B3bReconcile {
    # Classify one round and write reconcile.json + the reconcile journal record.
    # Returns @{ok; outcome; verdict; voidReason; stageHit; mechanismOk; sessionBlock; error}.
    param(
        [string]$PlanJsonPath = '',
        [string]$InventoryJsonPath = '',
        [string]$InventoryJson2Path = '',
        [string]$AuditTxtPath = '',
        [string]$MarkerTranscriptPath = '',
        # The transcript the host archived *before* the cut. Preferred over the copy
        # that comes back from the guest afterwards: the post-cut file can lose its
        # tail to the very power loss under test, and reading that as "the stage was
        # never reached" would misattribute the loss.
        [string]$PreCutTranscriptPath = '',
        [Parameter(Mandatory = $true)][string]$OutputPath,
        [string]$JournalPath = '',
        [int]$Round = 0,
        [string]$Candidate = 'c1',
        [string]$Stage = '',
        [string]$TargetPath = '',
        [string]$ContentDigest = '',
        [long]$ContentSize = 0,
        [switch]$Paired,
        [string]$TargetKind = 'request',
        [string]$VoidReason = $null,
        [double]$PlanCutDelayS = 2.0
    )
    $result = [pscustomobject]@{
        ok = $false; outcome = $null; verdict = ''; voidReason = $null
        stageHit = $false; stageEscape = $false; mechanismOk = $false; orphanCompletion = $false
        sessionBlock = $null; error = $null
    }
    try {
        $stageHit = $false
        $stageEscape = $false
        $mechanismOk = $false
        $orphanCompletion = $false
        $sessionBlock = $null
        $markerSource = 'none'
        $markerPostCutDiffers = $false
        $preCutDigest = $null
        $markerDigest = $null
        $stagesText = ''
        if ($MarkerTranscriptPath -and (Test-Path -LiteralPath $MarkerTranscriptPath)) {
            $markerDigest = Get-B3aSha256Hex -Path $MarkerTranscriptPath
        }

        if (-not [string]::IsNullOrEmpty($VoidReason)) {
            $targetPath = $TargetPath; $digest = $ContentDigest; $size = $ContentSize; $stageEff = $Stage
            if ($PlanJsonPath -and (Test-Path -LiteralPath $PlanJsonPath)) {
                $plan = (Get-B3aFileTextUtf8 -Path $PlanJsonPath) | ConvertFrom-Json
                if (-not $targetPath) { $targetPath = $plan.targetPath }
                if (-not $digest) { $digest = $plan.contentDigest }
                if ($size -eq 0) { $size = [long]$plan.contentSize }
                if (-not $stageEff) { $stageEff = [string]$plan.stage }
            }
            $reconcile = [ordered]@{
                schema = 'b3b-reconcile'; v = 1
                round = $Round; candidate = $Candidate; stage = $stageEff
                targetPath = $targetPath; contentDigest = $digest; contentSize = $size
                paired = [bool]$Paired; targetKind = $TargetKind
                outcome = $null
                rule = 'void'
                inventoryDigests = @(); inventoryIdentical = $false
                probeState = ''; orphanCompletions = @()
                auditDigest = $null; auditVerdict = $null; dirty = $null
                markerDigest = $markerDigest; markerPreCutDigest = $preCutDigest
                markerSource = $markerSource; markerPostCutDiffers = $false
                markerStages = @(); stageHit = $false; stageEscape = $false; sessionBlock = $null
                mechanismOk = $false; mechanismMissing = @(); mechanismForbidden = @()
                verdict = 'void'; voidReason = $VoidReason
                hostTs = (Get-B3aHostTimestamp)
            }
            Write-B3aTextFile -Path $OutputPath -Text ($reconcile | ConvertTo-Json -Compress -Depth 20)
            $result.verdict = 'void'
            $result.voidReason = $VoidReason
            $result.ok = $true
        } else {
            if (-not $PlanJsonPath -or -not (Test-Path -LiteralPath $PlanJsonPath)) { throw 'plan.json missing' }
            if (-not $InventoryJsonPath -or -not (Test-Path -LiteralPath $InventoryJsonPath)) { throw 'inventory-1.json missing' }
            $plan = (Get-B3aFileTextUtf8 -Path $PlanJsonPath) | ConvertFrom-Json
            $inv1 = (Get-B3aFileTextUtf8 -Path $InventoryJsonPath) | ConvertFrom-Json
            # One derivation, shared with the matrix's later verification.
            $derived = Get-B3bDerivedOutcome -PlanJsonPath $PlanJsonPath -InventoryJsonPath $InventoryJsonPath
            $probeState = [string]$derived.probeState
            $probeFiles = @($derived.probeFiles)
            $outcome = $derived.outcome

            # Stage attribution prefers the transcript the host archived *before* the
            # cut; the post-cut copy is recorded for transparency and a disagreement
            # between the two is reported rather than silently resolved.
            $markerSource = 'none'
            $markerPostCutDiffers = $false
            $preCutDigest = $null
            $transcript = $null
            if ($PreCutTranscriptPath -and (Test-Path -LiteralPath $PreCutTranscriptPath)) {
                $preCutDigest = Get-B3aSha256Hex -Path $PreCutTranscriptPath
                $transcript = Read-B3bMarkerTranscript -Text (Get-B3aFileTextUtf8 -Path $PreCutTranscriptPath)
                $markerSource = 'pre-cut'
            } elseif ($MarkerTranscriptPath -and (Test-Path -LiteralPath $MarkerTranscriptPath)) {
                $transcript = Read-B3bMarkerTranscript -Text (Get-B3aFileTextUtf8 -Path $MarkerTranscriptPath)
                $markerSource = 'post-cut-only'
            }
            if ($null -ne $transcript) {
                $hitInfo = Get-B3bStageHit -Target ([string]$plan.stage) -Transcript $transcript -TargetPath ([string]$plan.targetPath)
                $stageHit = [bool]$hitInfo.hit
                $stageEscape = [bool]$hitInfo.escaped
                $stagesText = (@($hitInfo.stages) -join ',')
                $mechanism = Get-B3bMechanismVerdict -Arm $Candidate -Trace $hitInfo.targetTrace
                $mechanismOk = [bool]$mechanism.ok
                if (($markerSource -eq 'pre-cut') -and $MarkerTranscriptPath -and (Test-Path -LiteralPath $MarkerTranscriptPath)) {
                    $postTranscript = Read-B3bMarkerTranscript -Text (Get-B3aFileTextUtf8 -Path $MarkerTranscriptPath)
                    $postHit = Get-B3bStageHit -Target ([string]$plan.stage) -Transcript $postTranscript -TargetPath ([string]$plan.targetPath)
                    $markerPostCutDiffers = ((@($hitInfo.stages) -join ',') -ne (@($postHit.stages) -join ','))
                }
            } else {
                $mechanism = [ordered]@{ ok = $false; missing = @('marker transcript missing'); forbidden = @(); trace = @() }
            }
            # The writer was told to block inside the target stage. A later stage that
            # only the post-cut evidence reveals is still a protocol defect: it is
            # escalated as a session blocker instead of being folded into a round
            # verdict.
            if ($stageEscape) { $sessionBlock = 'stage-escape' }

            $orphans = @()
            $frozenSummary = $null
            if ($inv1.PSObject.Properties.Name -contains 'probeEntries') {
                if ($null -ne $inv1.probeEntries.PSObject.Properties['summary']) {
                    $frozenSummary = $inv1.probeEntries.summary
                }
            }
            # Pair R and C with the product naming; the frozen summary pairs by base
            # file name, which cannot match request.<id> against completion.<id>.
            $pairing = Get-B3bReplayPairing -Files $probeFiles
            $orphans = @($pairing.orphanCompletions)
            $orphanCompletion = ($orphans.Count -gt 0)

            $inventoryDigests = @((Get-B3aSha256Hex -Path $InventoryJsonPath))
            $inventoryIdentical = $false
            if ($InventoryJson2Path -and (Test-Path -LiteralPath $InventoryJson2Path)) {
                $inventoryDigests += (Get-B3aSha256Hex -Path $InventoryJson2Path)
                $inventoryIdentical = Test-B3aBytesEqual `
                    -A ([System.IO.File]::ReadAllBytes($InventoryJsonPath)) `
                    -B ([System.IO.File]::ReadAllBytes($InventoryJson2Path))
            }
            $auditDigest = $null; $auditVerdict = $null
            if ($AuditTxtPath -and (Test-Path -LiteralPath $AuditTxtPath)) {
                $auditDigest = Get-B3aSha256Hex -Path $AuditTxtPath
                $auditVerdict = (Read-B3aAuditVerdict -Text ([System.IO.File]::ReadAllText($AuditTxtPath))).verdict
            }
            $dirty = $null
            if ($inv1.PSObject.Properties.Name -contains 'dirty') { $dirty = $inv1.dirty }

            $reconcile = [ordered]@{
                schema = 'b3b-reconcile'; v = 1
                round = $Round; candidate = $Candidate; stage = [string]$plan.stage
                targetPath = $plan.targetPath
                contentDigest = $plan.contentDigest
                contentSize = $plan.contentSize
                paired = [bool]$Paired; targetKind = $TargetKind
                outcome = $outcome
                rule = 'outcome per B3b design section 7 (content+size), target resolved in the probe section'
                inventoryDigests = $inventoryDigests
                inventoryIdentical = $inventoryIdentical
                probeState = $probeState
                orphanCompletions = @($orphans)
                replayPairing = $pairing
                frozenProbeSummary = $frozenSummary
                auditDigest = $auditDigest
                auditVerdict = $auditVerdict
                dirty = $dirty
                markerDigest = $markerDigest
                markerPreCutDigest = $preCutDigest
                markerSource = $markerSource
                markerPostCutDiffers = $markerPostCutDiffers
                markerStages = $stagesText
                stageHit = $stageHit
                stageEscape = $stageEscape
                sessionBlock = $sessionBlock
                mechanismOk = $mechanismOk
                mechanismMissing = @($mechanism.missing)
                mechanismForbidden = @($mechanism.forbidden)
                verdict = 'ok'
                voidReason = $null
                hostTs = (Get-B3aHostTimestamp)
            }
            Write-B3aTextFile -Path $OutputPath -Text ($reconcile | ConvertTo-Json -Compress -Depth 20)
            $result.outcome = $outcome
            $result.verdict = 'ok'
        }

        if ($JournalPath) {
            $extra = [ordered]@{
                outcome = $result.outcome
                inventoryDigests = @($reconcile.inventoryDigests)
                auditDigest = $reconcile.auditDigest
                verdict = $reconcile.verdict
                voidReason = $reconcile.voidReason
            }
            $rec = New-B3aJournalRecord -Round $Round -Candidate $Candidate -Stage ([string]$reconcile.stage) `
                -TargetPath ([string]$reconcile.targetPath) `
                -ContentDigest ([string]$reconcile.contentDigest) `
                -ContentSize ([long]$reconcile.contentSize) `
                -PlanCutDelayS $PlanCutDelayS -Op 'reconcile' -Extra $extra
            $ack = Add-B3aJournalRecord -JournalPath $JournalPath -Record $rec
            if (-not $ack.ack) { throw ('reconcile journal ACK failed: ' + $ack.error) }
        }
        $result.stageHit = $stageHit
        $result.stageEscape = $stageEscape
        $result.mechanismOk = $mechanismOk
        $result.orphanCompletion = $orphanCompletion
        $result.sessionBlock = $sessionBlock
        $result.ok = $true
        return $result
    } catch {
        $result.error = $_.Exception.Message
        return $result
    }
}

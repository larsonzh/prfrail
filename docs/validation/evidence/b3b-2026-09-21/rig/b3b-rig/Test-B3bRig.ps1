# Test-B3bRig.ps1 - hermetic self-test for the B3b rig.
#
# Everything here runs on the host with no VM: static hygiene, the reuse pin, the
# verdict vocabulary, the replay pairing, the measured-phase filter and a full matrix
# rebuild from synthetic round artifacts. Each check is written so that deleting the
# guard it covers turns it red - the point of the self-test is to make the rig's rules
# falsifiable before any power is cut.
#
# Exit 0 when every check passes; exit 1 otherwise. PowerShell 5.1 compatible.
# Encoding: UTF-8 with BOM + LF. All output English.

param(
    [string]$B3bRigDir = '',
    [string]$B3aRigDir = '',
    [string]$PssaModulePath = 'C:\Users\妙妙呜\.vscode\extensions\ms-vscode.powershell-2025.4.0\modules\PSScriptAnalyzer'
)

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'Get-B3bVerdict.ps1')
. (Join-Path $PSScriptRoot 'Get-B3bReconcile.ps1')

$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrEmpty($B3bRigDir)) { $B3bRigDir = $PSScriptRoot }
if ([string]::IsNullOrEmpty($B3aRigDir)) {
    $B3aRigDir = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\b3a-rig'))
}

$script:checks = 0
$script:failures = 0

function Assert-B3bCheck {
    param(
        [Parameter(Mandatory = $true)][string]$Name,
        [Parameter(Mandatory = $true)][bool]$Ok,
        [string]$Detail = ''
    )
    $script:checks++
    if ($Ok) {
        Write-Output ('PASS  ' + $Name)
        return
    }
    $script:failures++
    Write-Output ('FAIL  ' + $Name + ' :: ' + $Detail)
}

$work = Join-Path $env:TEMP ('b3b-selftest-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $work -Force)

try {
    # --- A. static hygiene of every rig script ---------------------------------
    $scripts = @(Get-ChildItem -LiteralPath $B3bRigDir -Filter *.ps1 -File | Sort-Object Name)
    Assert-B3bCheck -Name 'scripts-present' -Ok ($scripts.Count -ge 7) -Detail ('found ' + $scripts.Count)
    $encodingBad = New-Object System.Collections.Generic.List[string]
    $parseBad = New-Object System.Collections.Generic.List[string]
    foreach ($script in $scripts) {
        if (-not (Test-B3aPs1Encoding -Path $script.FullName)) { $encodingBad.Add($script.Name) }
        $errors = $null
        [void][System.Management.Automation.Language.Parser]::ParseFile($script.FullName, [ref]$null, [ref]$errors)
        if ($errors.Count -gt 0) { $parseBad.Add(($script.Name + ':' + $errors.Count)) }
    }
    Assert-B3bCheck -Name 'ps1-encoding-bom-lf' -Ok ($encodingBad.Count -eq 0) -Detail (($encodingBad) -join ',')
    Assert-B3bCheck -Name 'parse-all-scripts' -Ok ($parseBad.Count -eq 0) -Detail (($parseBad) -join ',')

    if (Test-Path -LiteralPath $PssaModulePath) {
        Import-Module $PssaModulePath -ErrorAction Stop
        $analysis = @(Invoke-ScriptAnalyzer -Path $B3bRigDir -Recurse)
        Assert-B3bCheck -Name 'analyzer-clean' -Ok ($analysis.Count -eq 0) `
            -Detail (($analysis | ForEach-Object { $_.RuleName + '@' + $_.Line }) -join ',')
    } else {
        Assert-B3bCheck -Name 'analyzer-clean' -Ok $true -Detail 'skipped: PSScriptAnalyzer not found'
    }

    # --- B. module hygiene ------------------------------------------------------
    $hygiene = Test-B3bModuleHygiene -B3bRigDir $B3bRigDir -B3aRigDir $B3aRigDir
    Assert-B3bCheck -Name 'module-no-b3a-collision' -Ok ($hygiene.collisions.Count -eq 0) -Detail (($hygiene.collisions) -join ',')
    Assert-B3bCheck -Name 'module-no-toplevel-param' -Ok ($hygiene.topLevelParams.Count -eq 0) -Detail (($hygiene.topLevelParams) -join ',')
    Assert-B3bCheck -Name 'module-load-preserves-vars' -Ok ($hygiene.b3bFunctions.Count -ge 30) -Detail ('functions=' + $hygiene.b3bFunctions.Count)

    # --- C. reuse pin -----------------------------------------------------------
    $current = Get-B3bReusePinMap -B3aRigDir $B3aRigDir
    $missingPin = @((Get-B3bPinnedFileList) | Where-Object { [string]::IsNullOrEmpty($current[$_]) })
    Assert-B3bCheck -Name 'reuse-pin-complete' -Ok ($missingPin.Count -eq 0) -Detail (($missingPin) -join ',')
    $drifted = @{}
    foreach ($key in $current.Keys) { $drifted[$key] = $current[$key] }
    $drifted[(Get-B3bPinnedFileList)[0]] = 'deadbeef'
    $pinResult = Test-B3bPinIntegrity -B3aRigDir $B3aRigDir -Expected $drifted
    Assert-B3bCheck -Name 'reuse-pin-drift-red' -Ok (-not $pinResult.ok) -Detail 'a changed hash must block the session'

    # --- D. verdict vocabulary --------------------------------------------------
    $stageList = Get-B3bPublishStageList
    Assert-B3bCheck -Name 'stage-count' -Ok ($stageList.Count -eq 5) -Detail ($stageList -join ',')
    Assert-B3bCheck -Name 'allowed-early' -Ok ((Get-B3bAllowedOutcomeList -Stage 'temp-written') -join ',' -eq 'survived,lost-absent') -Detail 'early stages may be absent'
    Assert-B3bCheck -Name 'allowed-late' -Ok ((Get-B3bAllowedOutcomeList -Stage 'record-linked') -join ',' -eq 'survived') -Detail 'late stages may not be absent'
    Assert-B3bCheck -Name 'decision-c2' -Ok ((Get-B3bDecisionStageList -Candidate 'c2').Count -eq 2) -Detail 'c2 claims from record-linked'
    Assert-B3bCheck -Name 'decision-c1' -Ok ((Get-B3bDecisionStageList -Candidate 'c1').Count -eq 0) -Detail 'the baseline claims nothing'

    function Get-B3bSyntheticRound {
        param(
            [Parameter(Mandatory = $true)][int]$Round,
            [Parameter(Mandatory = $true)][string]$Arm,
            [Parameter(Mandatory = $true)][string]$Stage,
            [Parameter(Mandatory = $true)][string]$Outcome,
            [Parameter(Mandatory = $true)][bool]$Mechanism,
            [bool]$StageReached = $true,
            [bool]$Gates = $true
        )
        return (New-B3bRoundResult -Round $Round -Arm $Arm -Stage $Stage -Valid $true -Outcome $Outcome `
                -MechanismOk $Mechanism -StageHit $StageReached `
                -AuditOk $Gates -InventoryRepeatableOk $Gates -JournalChainOk $Gates)
    }
    $clean = @(1..5 | ForEach-Object { Get-B3bSyntheticRound -Round $_ -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true })
    Assert-B3bCheck -Name 'verdict-proven-full' -Ok ((Get-B3bCellVerdict -Arm 'c2' -Stage 'record-linked' -Rounds $clean -RequiredN 5 -DiscriminationGate $true).verdict -eq 'eligible') -Detail 'all green with mechanism and discrimination'
    Assert-B3bCheck -Name 'verdict-no-discrimination-unresolved' -Ok ((Get-B3bCellVerdict -Arm 'c2' -Stage 'record-linked' -Rounds $clean -RequiredN 5 -DiscriminationGate $false).verdict -eq 'unresolved') -Detail 'a device without discrimination power cannot pass'
    $lost = @(1..4 | ForEach-Object { Get-B3bSyntheticRound -Round $_ -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true })
    $lost += (Get-B3bSyntheticRound -Round 5 -Arm 'c2' -Stage 'record-linked' -Outcome 'lost-absent' -Mechanism $true)
    Assert-B3bCheck -Name 'verdict-lost-stage5-red' -Ok ((Get-B3bCellVerdict -Arm 'c2' -Stage 'record-linked' -Rounds $lost -RequiredN 5 -DiscriminationGate $true).verdict -eq 'disproven') -Detail 'one counterexample disproves the cell'
    $torn = @(1..4 | ForEach-Object { Get-B3bSyntheticRound -Round $_ -Arm 'c2' -Stage 'temp-written' -Outcome 'lost-absent' -Mechanism $true })
    $torn += (Get-B3bSyntheticRound -Round 5 -Arm 'c2' -Stage 'temp-written' -Outcome 'lost-torn-content' -Mechanism $true)
    Assert-B3bCheck -Name 'verdict-torn-red' -Ok ((Get-B3bCellVerdict -Arm 'c2' -Stage 'temp-written' -Rounds $torn -RequiredN 5 -DiscriminationGate $true).verdict -eq 'disproven') -Detail 'a torn record is never allowed'
    $noMechanism = @(1..5 | ForEach-Object { Get-B3bSyntheticRound -Round $_ -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $false })
    Assert-B3bCheck -Name 'verdict-mechanism-missing-unresolved' -Ok ((Get-B3bCellVerdict -Arm 'c2' -Stage 'record-linked' -Rounds $noMechanism -RequiredN 5 -DiscriminationGate $true).verdict -eq 'unresolved') -Detail 'no mechanism evidence means no claim'
    $earlySurvived = @(1..5 | ForEach-Object { Get-B3bSyntheticRound -Round $_ -Arm 'c3' -Stage 'temp-written' -Outcome 'survived' -Mechanism $true })
    Assert-B3bCheck -Name 'verdict-early-survived-legal' -Ok ((Get-B3bCellVerdict -Arm 'c3' -Stage 'temp-written' -Rounds $earlySurvived -RequiredN 5 -DiscriminationGate $true).verdict -eq 'eligible') -Detail 'surviving before the claim is not a violation'
    $floorRefused = $false
    try { [void](Get-B3bCellVerdict -Arm 'c2' -Stage 'record-linked' -Rounds $clean -RequiredN 4 -DiscriminationGate $true) } catch { $floorRefused = $true }
    Assert-B3bCheck -Name 'verdict-requiredN-floor' -Ok $floorRefused -Detail 'RequiredN below 5 must be refused'

    # A violation that carries no stage/mechanism evidence is a measurement artifact,
    # not a refutation: the writer may simply have died before publishing. It must cap
    # the cell at 'unresolved' instead of turning into 'disproven'.
    $unattributable = @(1..4 | ForEach-Object { Get-B3bSyntheticRound -Round $_ -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true })
    $unattributable += (Get-B3bSyntheticRound -Round 5 -Arm 'c2' -Stage 'record-linked' -Outcome 'lost-absent' -Mechanism $false -StageReached $false)
    $unattributableVerdict = Get-B3bCellVerdict -Arm 'c2' -Stage 'record-linked' -Rounds $unattributable -RequiredN 5 -DiscriminationGate $true
    Assert-B3bCheck -Name 'verdict-unattributable-violation-unresolved' -Ok ($unattributableVerdict.verdict -eq 'unresolved' -and $unattributableVerdict.unattributable -eq 1) `
        -Detail ('verdict=' + $unattributableVerdict.verdict + ' unattributable=' + $unattributableVerdict.unattributable)

    # Ruling that closes the c2c3 contradiction: its claim starts at parent-synced.
    $c2c3Decision = Get-B3bDecisionStageList -Candidate 'c2c3'
    Assert-B3bCheck -Name 'decision-c2c3-parent-synced-only' -Ok ($c2c3Decision.Count -eq 1 -and $c2c3Decision[0] -eq 'parent-synced') -Detail ($c2c3Decision -join ',')
    Assert-B3bCheck -Name 'decision-control-claims-nothing' -Ok ((Get-B3bDecisionStageList -Candidate 'positive-control').Count -eq 0) -Detail 'the control arm claims nothing'
    Assert-B3bCheck -Name 'arm-list-has-control' -Ok ((Get-B3bArmList) -contains 'positive-control') -Detail ((Get-B3bArmList) -join ',')
    # An informational violation cannot refute an arm (stages 1-3 run the same code in every
    # arm), but it must cap it: an arm may not be declared proven when the device misbehaved while
    # that arm ran. Both directions are pinned, so neither the gap nor an over-correction survives.
    $infoOnly = @(
        [ordered]@{ arm = 'c3'; stage = 'parent-synced'; verdict = 'eligible' },
        [ordered]@{ arm = 'c3'; stage = 'temp-written'; verdict = 'disproven' }
    )
    $infoVerdict = Get-B3bArmVerdict -Arm 'c3' -Cells $infoOnly
    Assert-B3bCheck -Name 'arm-info-cell-caps-the-arm-at-unresolved' -Ok (($infoVerdict.verdict -eq 'unresolved') -and ($infoVerdict.infoViolations.Count -eq 1)) `
        -Detail ([string]$infoVerdict.verdict)
    $cleanInfo = @(
        [ordered]@{ arm = 'c3'; stage = 'parent-synced'; verdict = 'eligible' },
        [ordered]@{ arm = 'c3'; stage = 'temp-written'; verdict = 'eligible' }
    )
    Assert-B3bCheck -Name 'arm-info-clean-can-still-be-proven' -Ok ((Get-B3bArmVerdict -Arm 'c3' -Cells $cleanInfo).verdict -eq 'proven') -Detail 'a clean informational cell must not block the claim'
    Assert-B3bCheck -Name 'verdict-document-carries-claim-scope' -Ok ((Get-B3bVerdictDocument -Cells @([ordered]@{ arm = 'c1'; stage = 'temp-written'; verdict = 'eligible' }) -RequiredN 5 -Discrimination @{ ok = $true } -Premise 'P').claimScope.Length -gt 80) -Detail 'claimScope must be carried by the verdict document'
    Assert-B3bCheck -Name 'arm-c1-never-proven' -Ok ((Get-B3bArmVerdict -Arm 'c1' -Cells @([ordered]@{ arm = 'c1'; stage = 'temp-written'; verdict = 'eligible' })).verdict -eq 'unresolved') -Detail 'the reference arm claims nothing'
    $c2Cells = @(
        [ordered]@{ arm = 'c2'; stage = 'record-linked'; verdict = 'eligible' },
        [ordered]@{ arm = 'c2'; stage = 'parent-synced'; verdict = 'eligible' }
    )
    Assert-B3bCheck -Name 'arm-c2-proven' -Ok ((Get-B3bArmVerdict -Arm 'c2' -Cells $c2Cells).verdict -eq 'proven') -Detail 'both decision cells eligible'
    $c2CellsPartial = @([ordered]@{ arm = 'c2'; stage = 'record-linked'; verdict = 'eligible' }, [ordered]@{ arm = 'c2'; stage = 'parent-synced'; verdict = 'unresolved' })
    Assert-B3bCheck -Name 'arm-c2-partial-unresolved' -Ok ((Get-B3bArmVerdict -Arm 'c2' -Cells $c2CellsPartial).verdict -eq 'unresolved') -Detail 'one unresolved decision cell holds the arm back'

    Assert-B3bCheck -Name 'discrimination-gate-ok' -Ok ((Get-B3bDiscriminationGate -NegativeValidRounds 5 -NegativeLosses 1 -PositiveValidRounds 5 -PositiveLosses 0).ok) -Detail 'one loss in the baseline plus a clean positive control'
    Assert-B3bCheck -Name 'discrimination-gate-no-loss-red' -Ok (-not (Get-B3bDiscriminationGate -NegativeValidRounds 5 -NegativeLosses 0 -PositiveValidRounds 5 -PositiveLosses 0).ok) -Detail 'a baseline that never loses proves nothing'

    # --- E. ordering verdict ----------------------------------------------------
    $u4Orphan = @(
        (Get-B3bSyntheticRound -Round 1 -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true) ,
        (Get-B3bSyntheticRound -Round 2 -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true) ,
        (Get-B3bSyntheticRound -Round 3 -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true)
    )
    $u4Orphan[0]['orphanCompletion'] = $true
    Assert-B3bCheck -Name 'verdict-orphan-C-disproven' -Ok ((Get-B3bU4Verdict -Rounds $u4Orphan -RequiredN 3).verdict -eq 'disproven') -Detail 'a completion without its request falsifies ordering'
    $u4Clean = @(1..3 | ForEach-Object {
        $r = Get-B3bSyntheticRound -Round $_ -Arm 'c2' -Stage 'record-linked' -Outcome 'survived' -Mechanism $true
        $r['orphanCompletion'] = $false
        $r
    })
    Assert-B3bCheck -Name 'verdict-orphan-clean-unresolved' -Ok ((Get-B3bU4Verdict -Rounds $u4Clean -RequiredN 3).verdict -eq 'unresolved') -Detail 'a clean sample cannot prove ordering'

    # --- F. replay pairing ------------------------------------------------------
    $paired = Get-B3bReplayPairing -Files @(
        [pscustomobject]@{ path = 'requests\request.b3b-r01.jsonl'; size = 10; sha256 = 'a' },
        [pscustomobject]@{ path = 'completions\completion.b3b-r01.jsonl'; size = 10; sha256 = 'b' }
    )
    Assert-B3bCheck -Name 'pairing-produces-no-false-pair-gap' -Ok ($paired.orphanCompletions.Count -eq 0 -and $paired.requestsWithoutCompletion.Count -eq 0) -Detail 'product naming must pair'
    $orphanOnly = Get-B3bReplayPairing -Files @([pscustomobject]@{ path = 'completions\completion.b3b-r01.jsonl'; size = 10; sha256 = 'b' })
    Assert-B3bCheck -Name 'pairing-detects-orphan' -Ok ($orphanOnly.orphanCompletions.Count -eq 1) -Detail 'a lone completion is an orphan'
    $missingCompletion = Get-B3bReplayPairing -Files @(
        [pscustomobject]@{ path = 'requests\request.b3b-r01.jsonl'; size = 10; sha256 = 'a' },
        [pscustomobject]@{ path = 'requests\request.b3b-r02.jsonl'; size = 10; sha256 = 'a' },
        [pscustomobject]@{ path = 'completions\completion.b3b-r01.jsonl'; size = 10; sha256 = 'b' }
    )
    Assert-B3bCheck -Name 'pairing-detects-missing-completion' -Ok ($missingCompletion.requestsWithoutCompletion.Count -eq 1) -Detail 'r02 has no completion'

    # --- G. measured-phase filter ----------------------------------------------
    $requestLeaf = 'requests\request.b3b-r09.jsonl'
    $completionLeaf = 'completions\completion.b3b-r09.jsonl'
    $entries = @()
    foreach ($stage in $stageList) { $entries += [pscustomobject]@{ stage = $stage; path = $requestLeaf; trace = @('stage:' + $stage) } }
    foreach ($stage in @('temp-written', 'temp-synced', 'temp-closed')) { $entries += [pscustomobject]@{ stage = $stage; path = $completionLeaf; trace = @('stage:' + $stage) } }
    $entries += [pscustomobject]@{ stage = 'record-linked'; path = $completionLeaf; trace = @('MoveFileExW(WRITE_THROUGH,no-replace)') }
    $transcript = [pscustomobject]@{ entries = $entries; unparseable = @() }
    $targetPath = 'R:\agent-runner-replay\completions\completion.b3b-r09.jsonl'
    $phaseStages = Get-B3bMeasuredPhaseStageList -Entries $entries -TargetPath $targetPath
    Assert-B3bCheck -Name 'phase-filter-excludes-request-phase' -Ok (($phaseStages -join ',') -eq 'temp-written,temp-synced,temp-closed,record-linked') -Detail ($phaseStages -join ',')
    $hit = Get-B3bStageHit -Target 'record-linked' -Transcript $transcript -TargetPath $targetPath
    Assert-B3bCheck -Name 'phase-filter-trace-from-measured-phase' -Ok ($hit.hit -and ($hit.targetTrace -join ',') -eq 'MoveFileExW(WRITE_THROUGH,no-replace)') -Detail ($hit.targetTrace -join ',')
    $escapedTranscript = [pscustomobject]@{ entries = ($entries + [pscustomobject]@{ stage = 'parent-synced'; path = $completionLeaf; trace = @('stage:parent-synced') }); unparseable = @() }
    $escapeHit = Get-B3bStageHit -Target 'record-linked' -Transcript $escapedTranscript -TargetPath $targetPath
    Assert-B3bCheck -Name 'stage-escape-detected' -Ok ($escapeHit.escaped) -Detail 'a later stage inside the measured phase is a protocol defect'

    # --- H. writer handshake parsing -------------------------------------------
    $plan = Get-B3bPlanFromText -Text "PLAN path=R:\s\requests\request.b3b-r01.jsonl digest=$('a' * 64) size=512"
    Assert-B3bCheck -Name 'plan-parse' -Ok ($null -ne $plan -and $plan.size -eq 512) -Detail 'pre-commitment line'
    Assert-B3bCheck -Name 'plan-absent' -Ok ($null -eq (Get-B3bPlanFromText -Text 'no plan here')) -Detail 'a missing pre-commitment must be $null'
    $markerParse = Read-B3bMarkerTranscript -Text ('{"stage":"a"}' + "`n" + 'garbage')
    Assert-B3bCheck -Name 'marker-unparseable-reported' -Ok ($markerParse.unparseable.Count -eq 1) -Detail 'unparseable marker lines are reported, never skipped silently'

    # --- I. journal round-trip with the frozen writer --------------------------
    $journal = Join-Path $work 'journal.jsonl'
    $record = New-B3aJournalRecord -Round 1 -Candidate 'c2' -Stage 'record-linked' `
        -TargetPath 'R:\agent-runner-replay\requests\request.b3b-r01.jsonl' `
        -ContentDigest ('b' * 64) -ContentSize 512 -PlanCutDelayS 2.0 -Op 'plan'
    $ack = Add-B3aJournalRecord -JournalPath $journal -Record $record
    Assert-B3bCheck -Name 'journal-ack' -Ok ($ack.ack) -Detail ([string]$ack.error)
    $chain = Test-B3aJournalChain -JournalPath $journal
    Assert-B3bCheck -Name 'journal-chain-ok' -Ok $chain.ok -Detail ([string]$chain.firstError)
    $text = [IO.File]::ReadAllText($journal)
    $digest = 'b' * 64
    $tamperedText = $text.Replace($digest, ('c' * 64))
    Assert-B3bCheck -Name 'journal-tamper-applied' -Ok ($tamperedText -ne $text) -Detail 'the synthetic tamper must really change bytes'
    [IO.File]::WriteAllText($journal, $tamperedText, [Text.UTF8Encoding]::new($false))
    $chainTampered = Test-B3aJournalChain -JournalPath $journal
    Assert-B3bCheck -Name 'journal-content-tamper-red' -Ok (-not $chainTampered.ok) -Detail 'an edited record body must fail the hash recomputation'
    [IO.File]::WriteAllText($journal, $text, [Text.UTF8Encoding]::new($false))
    $hashMatch = [regex]::Match($text, '"recordHash":"([0-9a-f]{64})"')
    Assert-B3bCheck -Name 'journal-recordHash-present' -Ok $hashMatch.Success -Detail 'the stored line must carry recordHash'
    if ($hashMatch.Success) {
        $storedHash = $hashMatch.Groups[1].Value
        $firstChar = 'a'
        if ($storedHash[0] -eq 'a') { $firstChar = 'b' }
        $editedHash = $firstChar + $storedHash.Substring(1)
        $hashTampered = $text.Replace('"recordHash":"' + $storedHash + '"', '"recordHash":"' + $editedHash + '"')
        Assert-B3bCheck -Name 'journal-hash-tamper-applied' -Ok ($hashTampered -ne $text) -Detail 'the edited recordHash must differ'
        [IO.File]::WriteAllText($journal, $hashTampered, [Text.UTF8Encoding]::new($false))
        $chainHash = Test-B3aJournalChain -JournalPath $journal
        Assert-B3bCheck -Name 'journal-hash-tamper-red' -Ok (-not $chainHash.ok) -Detail 'an edited recordHash must fail the recomputation'
        [IO.File]::WriteAllText($journal, $text, [Text.UTF8Encoding]::new($false))
    }

    # --- J. matrix rebuild from round artifacts --------------------------------
    # The fixtures are built through the *real* reconciliation so that the artifacts,
    # reconcile.json and the journal agree by construction - exactly as they do in a
    # live round. That is what makes the tamper checks below meaningful: the only
    # difference is the edit the check itself makes.
    function New-B3bFixtureRound {
        # Builds a fixture round directory (returns an in-memory result).
        [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '', Justification = 'Builds a throwaway fixture tree under the self-test scratch directory; the caller owns the directory and it is removed in finally.')]
        param(
            [Parameter(Mandatory = $true)][string]$SessionRoot,
            [Parameter(Mandatory = $true)][int]$Round,
            [Parameter(Mandatory = $true)][string]$Arm,
            [Parameter(Mandatory = $true)][string]$Stage,
            [Parameter(Mandatory = $true)][string]$Kind,
            [bool]$Mechanism = $true,
            [bool]$Paired = $false,
            [switch]$RequestPresent,
            [switch]$Void
        )
        $roundDir = Join-Path (Join-Path $SessionRoot 'rounds') ('r{0:D2}' -f $Round)
        [void](New-Item -ItemType Directory -Path $roundDir -Force)
        $digest = 'c' * 64
        $size = 512
        $sub = 'requests'
        $leaf = 'request.b3b-r{0:D2}.jsonl' -f $Round
        if ($Paired) {
            $sub = 'completions'
            $leaf = 'completion.b3b-r{0:D2}.jsonl' -f $Round
        }
        $targetPath = 'R:\agent-runner-replay\' + $sub + '\' + $leaf
        $plan = [ordered]@{
            schema = 'b3b-plan'; v = 1; round = $Round; candidate = $Arm; stage = $Stage
            targetPath = $targetPath; contentDigest = $digest; contentSize = $size
            hostTs = '2026-09-20T00:00:00Z'
        }
        Write-B3aTextFile -Path (Join-Path $roundDir 'plan.json') -Text ($plan | ConvertTo-Json -Compress -Depth 5)
        $files = @()
        if ($Kind -eq 'survived') { $files = @([ordered]@{ path = $targetPath; size = $size; sha256 = $digest }) }
        elseif ($Kind -eq 'torn') { $files = @([ordered]@{ path = $targetPath; size = $size; sha256 = ('e' * 64) }) }
        if ($RequestPresent) {
            # An ordering round whose request record is still there: the honest pairing result is
            # "no orphan", which is what makes the tamper fixture below a real falsification.
            $requestPath = 'R:\agent-runner-replay\requests\request.b3b-r{0:D2}.jsonl' -f $Round
            $files = @([ordered]@{ path = $requestPath; size = $size; sha256 = $digest }) + $files
        }
        $inv = [ordered]@{
            schema = 'b3a-round-inventory'; v = 1
            roundRoot = ('<scratch>\r' + $Round); probeRoot = 'R:\agent-runner-replay'
            probeState = 'ok'; dirty = 'false'
            probeEntries = [ordered]@{ files = $files; summary = [ordered]@{ present = @($files).Count; absent = 0 } }
            roundRootEntries = [ordered]@{ files = @(); summary = [ordered]@{ present = 0; absent = 0 } }
        }
        $invText = ($inv | ConvertTo-Json -Compress -Depth 8)
        Write-B3aTextFile -Path (Join-Path $roundDir 'inventory-1.json') -Text $invText
        Write-B3aTextFile -Path (Join-Path $roundDir 'inventory-2.json') -Text $invText
        Write-B3aTextFile -Path (Join-Path $roundDir 'audit.txt') -Text 'VERDICT: HARD-POWER-LOSS (uptime = 0.3 minutes)'
        $stageOrder = Get-B3bPublishStageList
        $targetIndex = [array]::IndexOf($stageOrder, $Stage)
        $entries = @()
        for ($i = 0; $i -le $targetIndex; $i++) {
            $lineStage = $stageOrder[$i]
            $trace = @('stage:' + $lineStage)
            if ($Mechanism) {
                # Cumulative trace, the way the writer emits it: a token appears from the
                # stage where that primitive runs onwards, so the later decision stages
                # still carry the evidence of the earlier primitive.
                $linkIndex = [array]::IndexOf($stageOrder, 'record-linked')
                $parentIndex = [array]::IndexOf($stageOrder, 'parent-synced')
                foreach ($token in (Get-B3bMechanismToken -Arm $Arm)) {
                    $tokenIndex = $linkIndex
                    if ($token -like 'FlushFileBuffers*') { $tokenIndex = $parentIndex }
                    if ($i -ge $tokenIndex) { $trace += $token }
                }
            }
            $entries += [ordered]@{
                round = $Round; candidate = $Arm; stage = $lineStage; path = $targetPath
                trace = $trace; timestamp = ('2026-09-20T00:00:{0:D2}Z' -f $i)
            }
        }
        $transcriptText = (($entries | ForEach-Object { $_ | ConvertTo-Json -Compress -Depth 4 }) -join "`n")
        Write-B3aTextFile -Path (Join-Path $roundDir 'marker-transcript.txt') -Text $transcriptText
        Write-B3aTextFile -Path (Join-Path $roundDir 'marker.txt') -Text $transcriptText
        $targetKind = 'request'
        if ($Paired) { $targetKind = 'completion' }
        $recParams = @{
            PlanJsonPath = (Join-Path $roundDir 'plan.json')
            InventoryJsonPath = (Join-Path $roundDir 'inventory-1.json')
            InventoryJson2Path = (Join-Path $roundDir 'inventory-2.json')
            AuditTxtPath = (Join-Path $roundDir 'audit.txt')
            PreCutTranscriptPath = (Join-Path $roundDir 'marker-transcript.txt')
            MarkerTranscriptPath = (Join-Path $roundDir 'marker.txt')
            OutputPath = (Join-Path $roundDir 'reconcile.json')
            JournalPath = (Join-Path $SessionRoot 'journal.jsonl')
            Round = $Round
            Candidate = $Arm
            Stage = $Stage
            TargetPath = $targetPath
            ContentDigest = $digest
            ContentSize = $size
            TargetKind = $targetKind
            PlanCutDelayS = 2.0
        }
        if ($Paired) { $recParams['Paired'] = $true }
        if ($Void) {
            # A real void round never reaches the inventory: every void path in the driver exits
            # before it is taken. The fixture removes the artifacts such a round would not have and
            # reconciles as void, so the finalize path that must tolerate a voided U4 round is
            # really exercised instead of being assumed.
            foreach ($name in @('plan.json', 'inventory-1.json', 'inventory-2.json', 'audit.txt', 'marker-transcript.txt', 'marker.txt')) {
                $stale = Join-Path $roundDir $name
                if (Test-Path -LiteralPath $stale) { Remove-Item -LiteralPath $stale -Force }
            }
            $voidParams = @{
                OutputPath = (Join-Path $roundDir 'reconcile.json')
                JournalPath = (Join-Path $SessionRoot 'journal.jsonl')
                Round = $Round
                Candidate = $Arm
                Stage = $Stage
                TargetPath = $targetPath
                ContentDigest = $digest
                ContentSize = $size
                TargetKind = $targetKind
                VoidReason = 'marker-missed'
            }
            if ($Paired) { $voidParams['Paired'] = $true }
            $processed = Invoke-B3bReconcile @voidParams
        } else {
            $processed = Invoke-B3bReconcile @recParams
        }
        if (-not $processed.ok) { throw ('fixture reconcile failed for round ' + $Round + ': ' + $processed.error) }
        return $processed
    }

    $matrixRoot = Join-Path $work 'matrix'
    # The fixture session id is pinned and threaded into every matrix call. Relying on the
    # matrix defaulting to today's date made the whole suite time-dependent: it passed while
    # the wall clock agreed with this literal and turned 7 assertions red after midnight,
    # because finalize then resolved a different (empty) session directory.
    $fixtureSessionId = '2026-09-20'
    $sessionDir = Join-Path $matrixRoot $fixtureSessionId
    [void](New-Item -ItemType Directory -Path (Join-Path $sessionDir 'rounds') -Force)
    $index = 0
    foreach ($arm in @('c1', 'c2', 'c3')) {
        foreach ($stage in $stageList) {
            for ($i = 1; $i -le 5; $i++) {
                $kind = 'survived'
                # One genuine loss in the baseline's earliest cell is what gives the
                # device any discrimination power at all.
                if ($arm -eq 'c1' -and $stage -eq $stageList[0] -and $i -eq 1) { $kind = 'lost-absent' }
                [void](New-B3bFixtureRound -SessionRoot $sessionDir -Round $index -Arm $arm -Stage $stage -Kind $kind)
                $index++
            }
        }
    }
    for ($i = 1; $i -le 5; $i++) {
        [void](New-B3bFixtureRound -SessionRoot $sessionDir -Round $index -Arm 'positive-control' -Stage 'parent-synced' -Kind 'survived')
        $index++
    }
    # One paired (ordering) round, so the U4 re-derivation path is exercised end to end. Both
    # records are present, so the honest reading is "no orphan observed"; the tamper fixture below
    # then claims a phantom orphan the pinned inventory does not support.
    [void](New-B3bFixtureRound -SessionRoot $sessionDir -Round 500 -Arm 'c2' -Stage 'record-linked' -Kind 'survived' -Paired $true -RequestPresent)
    # A voided U4 round: allowed inside the per-cell void budget, and it carries no inventory. The
    # finalize has to tolerate it rather than refuse an honest session.
    [void](New-B3bFixtureRound -SessionRoot $sessionDir -Round 501 -Arm 'c2' -Stage 'record-linked' -Kind 'survived' -Paired $true -RequestPresent -Void)
    # M2 tolerance fixture: the marker copy pulled back after the cut may lose its tail
    # to the very power loss under test, so a truncated post-cut copy must not be able to
    # invalidate an honest round. Round 14 is c3/parent-synced - c3's decision cell - so
    # the main assertions below also prove the tolerance.
    $truncatedDir = Join-Path (Join-Path $sessionDir 'rounds') 'r14'
    $preCutLines = @((([IO.File]::ReadAllText((Join-Path $truncatedDir 'marker-transcript.txt'))) -split "`n") | Where-Object { $_ -ne '' })
    [IO.File]::WriteAllText((Join-Path $truncatedDir 'marker.txt'), $preCutLines[0], [Text.UTF8Encoding]::new($false))
    Assert-B3bCheck -Name 'post-cut-marker-truncated-fixture' -Ok (([IO.File]::ReadAllText((Join-Path $truncatedDir 'marker.txt'))).Length -lt ([IO.File]::ReadAllText((Join-Path $truncatedDir 'marker-transcript.txt'))).Length) -Detail 'the tolerance fixture must really shorten the post-cut copy'
    $finalizeOutput = & (Join-Path $B3bRigDir 'Invoke-B3bMatrix.ps1') -Finalize -RigRoot $matrixRoot -SessionId $fixtureSessionId `
        -Premise 'self-test' -WriterExe (Join-Path $work 'writer.exe') 2>&1 | Out-String
    $finalizeExit = $LASTEXITCODE
    $verdictPath = Join-Path $sessionDir 'verdict.json'
    Assert-B3bCheck -Name 'matrix-finalize-writes-verdict' -Ok ((Test-Path -LiteralPath $verdictPath) -and $finalizeExit -eq 0) -Detail ('exit=' + $finalizeExit + ' ' + $finalizeOutput)
    if (Test-Path -LiteralPath $verdictPath) {
        $verdict = ([IO.File]::ReadAllText($verdictPath)) | ConvertFrom-Json
        $c1Verdict = (@($verdict.arms | Where-Object { $_.arm -eq 'c1' })).verdict
        $c2Verdict = (@($verdict.arms | Where-Object { $_.arm -eq 'c2' })).verdict
        Assert-B3bCheck -Name 'matrix-c1-unresolved' -Ok ($c1Verdict -eq 'unresolved') -Detail ([string]$c1Verdict)
        Assert-B3bCheck -Name 'matrix-c2-proven' -Ok ($c2Verdict -eq 'proven') -Detail ([string]$c2Verdict)
        $c3Verdict = (@($verdict.arms | Where-Object { $_.arm -eq 'c3' })).verdict
        Assert-B3bCheck -Name 'matrix-c3-proven-despite-truncated-post-cut-marker' -Ok ($c3Verdict -eq 'proven') -Detail ([string]$c3Verdict)
        Assert-B3bCheck -Name 'void-u4-round-does-not-block-finalize' -Ok ($null -ne $verdict.u4) `
            -Detail 'a voided paired round is normal and must not refuse the session'
        Assert-B3bCheck -Name 'matrix-verdict-carries-premise-and-scope' -Ok (([string]$verdict.premise -eq 'self-test') -and ([string]$verdict.claimScope).Length -gt 80) -Detail ('premise=' + $verdict.premise)
    }

    # A finalize pointed at the wrong (or missing) -SessionId must refuse instead of writing a
    # verdict built from zero rounds. That failure mode is invisible at read time: an empty
    # verdict is shaped exactly like an honest 'unresolved', so a whole session's evidence
    # could be discarded without anything looking wrong.
    $emptySessionRoot = Join-Path $work 'matrix-empty'
    [void](New-Item -ItemType Directory -Path $emptySessionRoot -Force)
    $emptyOut = & (Join-Path $B3bRigDir 'Invoke-B3bMatrix.ps1') -Finalize -RigRoot $emptySessionRoot -SessionId $fixtureSessionId `
        -Premise 'self-test' -WriterExe (Join-Path $work 'writer.exe') 2>&1 | Out-String
    Assert-B3bCheck -Name 'finalize-refuses-a-session-with-no-rounds' -Ok (($LASTEXITCODE -eq 4) -and ($emptyOut -match 'REFUSED')) `
        -Detail ('exit=' + $LASTEXITCODE + ' ' + $emptyOut.Trim())
    # ...and the verdict has to land in the requested session, not in a directory named after
    # the current date: that is what silently voided this suite after midnight.
    $sessionDirs = @(Get-ChildItem -LiteralPath $matrixRoot -Directory)
    Assert-B3bCheck -Name 'finalize-honours-the-requested-session-id' -Ok (($sessionDirs.Count -eq 1) -and ($sessionDirs[0].Name -eq $fixtureSessionId)) `
        -Detail (($sessionDirs | ForEach-Object { $_.Name }) -join ',')

    # The inventory pair a round records must be two *consecutive identical* reads, not simply
    # the first two reads. The live run voided three honest c1 rounds because the first read
    # after the cut could still report a pre-settle view (size 0 then 1082 for the same
    # unlinked temp file, every other field identical).
    $invSynthetic = @(
        @{ name = 'settles-on-the-first-pair'; data = @('a', 'a'); stable = $true; tries = 2; pair = @('a', 'a') }
        @{ name = 'settles-on-the-second-pair'; data = @('a', 'b', 'b'); stable = $true; tries = 3; pair = @('b', 'b') }
        @{ name = 'never-settles'; data = @('a', 'b', 'a', 'b'); stable = $false; tries = 4 }
        @{ name = 'single-read-never-settles'; data = @('a'); stable = $false; tries = 1 }
    )
    foreach ($case in $invSynthetic) {
        $reads = @($case.data | ForEach-Object { [pscustomobject]@{ bytes = [Text.Encoding]::UTF8.GetBytes($_); text = [string]$_ } })
        $got = Get-B3bSettledInventoryPair -Reads $reads
        $pairOk = $true
        if ($case.stable) {
            $pairOk = ([Text.Encoding]::UTF8.GetString($got.first.bytes) -eq $case.pair[0]) -and ([Text.Encoding]::UTF8.GetString($got.second.bytes) -eq $case.pair[1])
        }
        Assert-B3bCheck -Name ('inventory-' + $case.name) -Ok (($got.stable -eq $case.stable) -and ($got.tries -eq $case.tries) -and $pairOk) `
            -Detail ('stable=' + $got.stable + ' tries=' + $got.tries + ' pairOk=' + $pairOk)
    }
    Assert-B3bCheck -Name 'inventory-retry-bound-leaves-room-to-settle' -Ok ($B3bInventoryMaxReads -ge 3) `
        -Detail ('bound=' + $B3bInventoryMaxReads)
    # The driver must keep every read it took and record the settled pair (not the first two).
    $roundInvSource = [IO.File]::ReadAllText((Join-Path $B3bRigDir 'Invoke-B3bRound.ps1'))
    Assert-B3bCheck -Name 'round-preserves-every-inventory-read' -Ok ($roundInvSource -match 'inventory-attempt-') `
        -Detail 'an unstable read is the evidence that explains the retry'
    Assert-B3bCheck -Name 'round-records-the-settled-pair' -Ok (($roundInvSource -match 'inventory-1\.json.\) -Bytes \$invSettled\.first\.bytes') -and ($roundInvSource -match 'inventory-2\.json.\) -Bytes \$invSettled\.second\.bytes')) `
        -Detail 'the recorded pair must be the settled one, not the first two reads'
    Assert-B3bCheck -Name 'round-records-how-many-reads-were-needed' -Ok ($roundInvSource -match 'inventory-reads\.txt') `
        -Detail 'the retry count is part of the round evidence'
    Assert-B3bCheck -Name 'round-preserves-the-unsettled-pair' -Ok (($roundInvSource -match 'inventory-mismatch-1\.bin') -and ($roundInvSource -match 'inventory-mismatch-2\.bin')) `
        -Detail 'a void on this reason must leave both streams behind'
    # A failed read attempt is retried inside the round, not by redoing the whole round, and a
    # round that never got a single readable inventory must abort rather than guess.
    Assert-B3bCheck -Name 'round-retries-a-failed-inventory-read-in-place' -Ok (($roundInvSource -match 'inventory read failed; ') -and ($roundInvSource -match 'continue')) `
        -Detail 'redoing the whole round for a guest-channel hiccup costs a fresh cut'
    Assert-B3bCheck -Name 'inventory-failure-log-carries-the-exit-code' -Ok ($roundInvSource -match 'exit=\{1\}') `
        -Detail 'logging only that it failed leaves the cause unknowable'
    Assert-B3bCheck -Name 'round-aborts-when-every-inventory-read-fails' -Ok (($roundInvSource -match 'every inventory read failed') -and ($roundInvSource -match 'inventory-failures\.txt') -and ($roundInvSource -match 'exit 5')) `
        -Detail 'no readable inventory means the round cannot be interpreted'

    function Test-B3bFinalizeRefusal {
        param(
            [Parameter(Mandatory = $true)][string]$Name,
            [Parameter(Mandatory = $true)][scriptblock]$Tamper,
            [Parameter(Mandatory = $true)][int]$ExpectedExit
        )
        $clone = Join-Path $work ('clone-' + $Name)
        if (Test-Path -LiteralPath $clone) { Remove-Item -LiteralPath $clone -Recurse -Force }
        Copy-Item -LiteralPath $matrixRoot -Destination $clone -Recurse -Force
        # A copied session tree is a *new* session, so it must not inherit the lock of the run
        # that produced it: otherwise the tamper checks would be refused for the wrong reason.
        $inheritedLock = Join-Path (Join-Path $clone $fixtureSessionId) 'session.lock'
        if (Test-Path -LiteralPath $inheritedLock) { Remove-Item -LiteralPath $inheritedLock -Force }
        & $Tamper $clone
        $out = & (Join-Path $B3bRigDir 'Invoke-B3bMatrix.ps1') -Finalize -RigRoot $clone -SessionId $fixtureSessionId `
            -Premise 'self-test' -WriterExe (Join-Path $work 'writer.exe') 2>&1 | Out-String
        $code = $LASTEXITCODE
        Assert-B3bCheck -Name $Name -Ok (($code -eq $ExpectedExit) -and ($out -match 'REFUSED')) `
            -Detail ('exit=' + $code + ' expected=' + $ExpectedExit + ' output=' + $out.Trim())
    }

    # FT1: editing reconcile.json must not be able to change a verdict. The tamper
    # removes the baseline's single loss, i.e. it tries to erase the evidence that the
    # device can discriminate.
    Test-B3bFinalizeRefusal -Name 'refuses-tampered-reconcile-outcome' -ExpectedExit 4 -Tamper {
        param($cloneRoot)
        $target = Join-Path $cloneRoot ($fixtureSessionId + '\rounds\r00\reconcile.json')
        $text = [IO.File]::ReadAllText($target)
        [IO.File]::WriteAllText($target, $text.Replace('"outcome":"lost-absent"', '"outcome":"survived"'), [Text.UTF8Encoding]::new($false))
    }
    # FT1b: editing only the inventory is caught by the recorded-artifact hashes.
    Test-B3bFinalizeRefusal -Name 'refuses-tampered-inventory' -ExpectedExit 4 -Tamper {
        param($cloneRoot)
        $target = Join-Path $cloneRoot ($fixtureSessionId + '\rounds\r01\inventory-1.json')
        $before = (Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash.ToLowerInvariant()
        $text = [IO.File]::ReadAllText($target)
        [IO.File]::WriteAllText($target, ($text + ' '), [Text.UTF8Encoding]::new($false))
        $after = (Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash.ToLowerInvariant()
        # A tamper that silently changes nothing would make this check pass for the wrong
        # reason, so the tamper proves itself first.
        Assert-B3bCheck -Name 'tamper-changed-the-artifact' -Ok ($before -ne $after) -Detail 'the tamper must really change the file'
    }
    # FT4: a round that records its pairing flag must also sit in the reserved range;
    # a mismatch must be refused instead of silently re-bucketed.
    Test-B3bFinalizeRefusal -Name 'refuses-paired-numbering-mismatch' -ExpectedExit 4 -Tamper {
        param($cloneRoot)
        $target = Join-Path $cloneRoot ($fixtureSessionId + '\rounds\r04\reconcile.json')
        $text = [IO.File]::ReadAllText($target)
        [IO.File]::WriteAllText($target, $text.Replace('"paired":false', '"paired":true'), [Text.UTF8Encoding]::new($false))
    }
    # FT3: a stage-escape discovered after the cut is a session blocker, not a footnote.
    Test-B3bFinalizeRefusal -Name 'refuses-session-block' -ExpectedExit 6 -Tamper {
        param($cloneRoot)
        $target = Join-Path $cloneRoot ($fixtureSessionId + '\rounds\r06\reconcile.json')
        $text = [IO.File]::ReadAllText($target)
        [IO.File]::WriteAllText($target, $text.Replace('"sessionBlock":null', '"sessionBlock":"stage-escape"'), [Text.UTF8Encoding]::new($false))
    }
    # V4 Pro closure finding: the U4 verdict reads a claim from reconcile.json, so a single-field
    # edit must not be able to swallow a real ordering falsification.
    Test-B3bFinalizeRefusal -Name 'refuses-tampered-u4-orphan' -ExpectedExit 4 -Tamper {
        param($cloneRoot)
        $target = Join-Path $cloneRoot ($fixtureSessionId + '\rounds\r500\reconcile.json')
        $text = [IO.File]::ReadAllText($target)
        $edited = $text.Replace('"orphanCompletions":[]', '"orphanCompletions":["completion.phantom.jsonl"]')
        Assert-B3bCheck -Name 'u4-tamper-applied' -Ok ($edited -ne $text) -Detail 'the tamper must really change the orphan claim'
        [IO.File]::WriteAllText($target, $edited, [Text.UTF8Encoding]::new($false))
    }
    # M1: a measured round may not be relabelled as a void round (that would drop a real
    # counterexample out of the valid set). The rebuttal is the void vocabulary: an
    # invented reason is not something the round driver can produce.
    Test-B3bFinalizeRefusal -Name 'refuses-tampered-void-relabel' -ExpectedExit 4 -Tamper {
        param($cloneRoot)
        $target = Join-Path $cloneRoot ($fixtureSessionId + '\rounds\r10\reconcile.json')
        $text = [IO.File]::ReadAllText($target)
        [IO.File]::WriteAllText($target, $text.Replace('"verdict":"ok","voidReason":null', '"verdict":"void","voidReason":"invented-by-a-tamperer"'), [Text.UTF8Encoding]::new($false))
    }
    # FT8: the control arm must be schedulable, otherwise the discrimination gate is
    # structurally false and every arm is capped at 'unresolved' for the wrong reason.
    # This parses the real control cell through the same resolver the command line uses.
    # The previous version only grepped the matrix source for the arm-list comparison and
    # therefore stayed green while the arm pattern above it rejected 'positive-control' -
    # the control calibration cell was unschedulable until the live run hit it.
    $controlArm = @(Get-B3bControlList)[0]
    $controlSpec = $null
    $controlParseError = ''
    try {
        $controlSpec = Resolve-B3bCellSpec -Cell ($controlArm + ':' + $B3bStageParentSynced) `
            -ArmList (Get-B3bArmList) -StageList (Get-B3bPublishStageList)
    } catch {
        $controlParseError = $_.Exception.Message
    }
    Assert-B3bCheck -Name 'control-arm-cell-is-schedulable' -Ok (($null -ne $controlSpec) -and ($controlSpec.arm -eq $controlArm) -and ($controlSpec.stage -eq $B3bStageParentSynced)) `
        -Detail ('cell=' + $controlArm + ':' + $B3bStageParentSynced + ' arm=' + $(if ($null -ne $controlSpec) { $controlSpec.arm } else { '<threw>' }) + ' error=' + $controlParseError)
    # ...and the resolver still has to refuse a cell it cannot honour, so the assertion above
    # cannot be satisfied by a parser that accepts anything.
    $bogusRejected = $false
    try {
        [void](Resolve-B3bCellSpec -Cell ('not-an-arm:' + $B3bStageParentSynced) -ArmList (Get-B3bArmList))
    } catch {
        $bogusRejected = $true
    }
    Assert-B3bCheck -Name 'cell-resolver-still-refuses-unknown-arms' -Ok $bogusRejected `
        -Detail 'the shared cell resolver must validate the arm, not only its shape'
    # The matrix must go through that resolver instead of growing its own pattern again.
    # Without this, a future edit could reintroduce an arm pattern that cannot match the
    # control arm while both assertions above stay green - the exact shape of the defect
    # the live run hit (matrix had an inline pattern, the self-test grepped only the
    # comparison line below it).
    $matrixCellSource = [IO.File]::ReadAllText((Join-Path $B3bRigDir 'Invoke-B3bMatrix.ps1'))
    Assert-B3bCheck -Name 'matrix-uses-the-shared-cell-resolver' -Ok (($matrixCellSource -match 'Resolve-B3bCellSpec') -and ($matrixCellSource -notmatch "\^\\\(\[a-z0-9")) `
        -Detail 'the matrix must not carry its own arm pattern'
    $matrixSource = [IO.File]::ReadAllText((Join-Path $B3bRigDir 'Invoke-B3bMatrix.ps1'))
    Assert-B3bCheck -Name 'checkpoint-exists-before-the-budget' -Ok ($matrixSource -match 'if \(\$Checkpoint\)') `
        -Detail 'the discrimination gate must have a mid-run decision point'

    # The void vocabulary and the driver's call sites have to stay in step: a reason the
    # driver can write but the vocabulary does not know would make every honest void round
    # look like tampering and refuse the whole finalize.
    $driverSource = [IO.File]::ReadAllText((Join-Path $B3bRigDir 'Invoke-B3bRound.ps1'))
    $driverReasons = New-Object System.Collections.Generic.List[string]
    # Both call shapes count: a bare literal and a parenthesised concatenation, which is how
    # the driver writes the interpolated reasons.
    foreach ($match in [regex]::Matches($driverSource, "Invoke-B3bVoid \(?'([^']+)'")) {
        $literal = $match.Groups[1].Value
        # Only the leading token is the vocabulary entry ('plan-no-ack: <error>' -> 'plan-no-ack').
        $token = (([string]$literal -split ' ')[0]).TrimEnd(':')
        if (-not $driverReasons.Contains($token)) { $driverReasons.Add($token) }
    }
    $unknownDriverReasons = @($driverReasons | Where-Object { -not (Test-B3bVoidReasonKnown -VoidReason $_) })
    Assert-B3bCheck -Name 'void-vocabulary-covers-driver' -Ok (($unknownDriverReasons.Count -eq 0) -and ($driverReasons.Count -ge 8)) `
        -Detail ('driver reasons not in the vocabulary: ' + (($unknownDriverReasons) -join ',') + ' (found ' + $driverReasons.Count + ')')
    Assert-B3bCheck -Name 'void-vocabulary-is-not-empty' -Ok ((Get-B3bVoidReasonList).Count -ge 8) -Detail ((Get-B3bVoidReasonList) -join ',')

    # A session must not be split by midnight: the id is resolved once and passed down, so
    # the round driver must never re-derive it from the clock on its own.
    $sessionSource = [IO.File]::ReadAllText((Join-Path $B3bRigDir 'New-B3bSession.ps1'))
    Assert-B3bCheck -Name 'session-id-resolved-once-and-threaded' `
        -Ok (($matrixSource -match "'-SessionId', \`$DateTag") -and ($driverSource -match 'IsNullOrEmpty\(\$SessionId\)') -and ($sessionSource -match 'IsNullOrEmpty\(\$SessionId\)')) `
        -Detail 'the matrix must pass the resolved session id to every round and the session preflight must honour an existing pin'
    Assert-B3bCheck -Name 'session-pin-pins-the-session-id' -Ok ($sessionSource -match 'existingPin\.sessionId') `
        -Detail 'an existing pin.json must win, so a re-run cannot move a started session'
    Assert-B3bCheck -Name 'session-id-conflict-is-refused' -Ok ($sessionSource -match 'REFUSED: the session directory already pins session') `
        -Detail 'an explicit -SessionId that disagrees with the pin must be refused, not silently adopted'
    # A missing '#' turns a comment into a command that only fails at run time: the parser
    # accepts 'after the round ...' as a command named 'after', so a parse check cannot see
    # it. This walks the AST instead - any bare command name (no hyphen) that is neither a
    # function the rig defines nor a resolvable command is a defect.
    $definedFunctions = New-Object System.Collections.Generic.List[string]
    foreach ($script in $scripts) {
        $scriptAst = [System.Management.Automation.Language.Parser]::ParseFile($script.FullName, [ref]$null, [ref]$null)
        $functionNodes = $scriptAst.FindAll({ param($node) $node -is [System.Management.Automation.Language.FunctionDefinitionAst] }, $true)
        foreach ($functionNode in $functionNodes) {
            if (-not $definedFunctions.Contains($functionNode.Name)) { $definedFunctions.Add($functionNode.Name) }
        }
    }
    $scriptKeywords = @('if', 'else', 'elseif', 'foreach', 'while', 'switch', 'return', 'throw', 'try', 'catch',
        'finally', 'function', 'param', 'exit', 'break', 'continue', 'do', 'for', 'filter', 'using', 'class',
        'enum', 'hidden', 'static', 'default', 'begin', 'process', 'end', 'workflow', 'dynamicparam', 'data', 'in')
    $bareCommands = New-Object System.Collections.Generic.List[string]
    foreach ($script in $scripts) {
        $ast = [System.Management.Automation.Language.Parser]::ParseFile($script.FullName, [ref]$null, [ref]$null)
        $commandNodes = $ast.FindAll({ param($node) $node -is [System.Management.Automation.Language.CommandAst] }, $true)
        foreach ($commandNode in $commandNodes) {
            $name = [string]$commandNode.GetCommandName()
            if ([string]::IsNullOrEmpty($name)) { continue }
            if ($name.Contains('-')) { continue }
            if ($scriptKeywords -contains $name.ToLowerInvariant()) { continue }
            if ($definedFunctions.Contains($name)) { continue }
            if ($null -ne (Get-Command -Name $name -ErrorAction SilentlyContinue)) { continue }
            $bareCommands.Add(($script.Name + ':' + $commandNode.Extent.StartLineNumber + ' ' + $name))
        }
    }
    Assert-B3bCheck -Name 'no-unresolvable-bare-command' -Ok ($bareCommands.Count -eq 0) -Detail (($bareCommands) -join ' | ')

    # --- K. evidence bundle, in-repo -------------------------------------------------
    # The bundle checks used to live in a one-off scratch script, which meant the F5/F6 closure
    # rested on a manual run nobody could repeat. They are part of the suite now: a minimal
    # session is built, the real builder runs, and both the layout and the premise propagation
    # are asserted (including the case where the verdict carries no premise at all).
    $bundleSession = Join-Path $work 'bundle-session'
    [void](New-Item -ItemType Directory -Path (Join-Path $bundleSession 'rounds\r00') -Force)
    Write-B3aTextFile -Path (Join-Path $bundleSession 'rounds\r00\reconcile.json') -Text '{"schema":"b3b-reconcile","v":1,"round":0}'
    Write-B3aTextFile -Path (Join-Path $bundleSession 'session.log') -Text 'bundle session log'
    Write-B3aTextFile -Path (Join-Path $bundleSession 'matrix.json') -Text '{"schema":"b3b-matrix","v":1}'
    Write-B3aTextFile -Path (Join-Path $bundleSession 'verdict.json') -Text '{"schema":"b3b-durability-verdict","arms":[{"arm":"c2","verdict":"proven"}],"premise":"unit-test premise","claimScope":"not a long-term guarantee"}'
    $bundleJournal = Join-Path $bundleSession 'journal.jsonl'
    $bundleRecord = New-B3aJournalRecord -Round 0 -Candidate 'c2' -Stage 'record-linked' `
        -TargetPath 'R:\x\y.jsonl' -ContentDigest ('a' * 64) -ContentSize 512 -PlanCutDelayS 2.0 -Op 'plan'
    $bundleAck = Add-B3aJournalRecord -JournalPath $bundleJournal -Record $bundleRecord
    Assert-B3bCheck -Name 'bundle-fixture-journal-ack' -Ok $bundleAck.ack -Detail ([string]$bundleAck.error)
    $pinDoc = [ordered]@{
        schema = 'b3b-pin'; v = 1; sessionId = 'unit-test'; recordedAt = '2026-09-20T00:00:00Z'
        b3aRigDir = $B3aRigDir; hashes = (Get-B3bReusePinMap -B3aRigDir $B3aRigDir)
    }
    Write-B3aTextFile -Path (Join-Path $bundleSession 'pin.json') -Text ($pinDoc | ConvertTo-Json -Depth 10) -WithBom
    # A builder refusal is a terminating error, not an exit code, so every call goes through this
    # wrapper: a red check must be reported, never allowed to take the whole suite down with it.
    function Invoke-B3bEvidenceBuilder {
        param(
            [Parameter(Mandatory = $true)][string]$RigRoot,
            [Parameter(Mandatory = $true)][string]$BundleDir,
            [Parameter(Mandatory = $true)][string]$RepoRoot
        )
        $script:LASTEXITCODE = 0
        $text = ''
        $failure = ''
        try {
            $text = & (Join-Path $B3bRigDir 'New-B3bEvidence.ps1') -RigRoot $RigRoot -BundleDir $BundleDir -RepoRoot $RepoRoot 2>&1 | Out-String
        } catch {
            $failure = $_.Exception.Message
        }
        $code = 0
        if ($failure -ne '') { $code = 1 } elseif ($null -ne $LASTEXITCODE) { $code = [int]$LASTEXITCODE }
        return [pscustomobject]@{ Exit = $code; Output = $text; Failure = $failure }
    }

    $bundleDir = Join-Path $work 'bundle-out'
    $bundleRun = Invoke-B3bEvidenceBuilder -RigRoot $bundleSession -BundleDir $bundleDir -RepoRoot $work
    $bundleOutput = ($bundleRun.Output + ' ' + $bundleRun.Failure).Trim()
    Assert-B3bCheck -Name 'bundle-builds-from-the-repository-suite' -Ok ($bundleRun.Exit -eq 0) -Detail $bundleOutput
    $bundleLayout = @('journal\journal.jsonl', 'journal\chain-verify.txt', 'session\pin.json', 'matrix\verdict.json', 'README.md', 'MANIFEST.md', 'SHA256SUMS.txt')
    $missingLayout = @($bundleLayout | Where-Object { -not (Test-Path -LiteralPath (Join-Path $bundleDir $_)) })
    Assert-B3bCheck -Name 'bundle-layout-complete' -Ok ($missingLayout.Count -eq 0) -Detail (($missingLayout) -join ',')
    if (Test-Path -LiteralPath (Join-Path $bundleDir 'README.md')) {
        $bundleReadme = [IO.File]::ReadAllText((Join-Path $bundleDir 'README.md'))
        Assert-B3bCheck -Name 'bundle-readme-carries-the-premise' -Ok ($bundleReadme -match 'unit-test premise') -Detail 'the premise must travel from verdict.json into the bundle'
    }

    # The evidence builder fails closed on a verdict without a premise, and it must say so
    # instead of writing a silently empty line.
    $forensicSession = Join-Path $work 'bundle-session-nopremise'
    if (Test-Path -LiteralPath $forensicSession) { Remove-Item -LiteralPath $forensicSession -Recurse -Force }
    Copy-Item -LiteralPath $bundleSession -Destination $forensicSession -Recurse -Force
    Write-B3aTextFile -Path (Join-Path $forensicSession 'verdict.json') -Text '{"schema":"b3b-durability-verdict","arms":[],"claimScope":"scope"}'
    $noPremiseDir = Join-Path $work 'bundle-out-nopremise'
    $noPremiseRun = Invoke-B3bEvidenceBuilder -RigRoot $forensicSession -BundleDir $noPremiseDir -RepoRoot $work
    $noPremiseOut = ($noPremiseRun.Output + ' ' + $noPremiseRun.Failure).Trim()
    $noPremiseReadme = ''
    if (Test-Path -LiteralPath (Join-Path $noPremiseDir 'README.md')) { $noPremiseReadme = [IO.File]::ReadAllText((Join-Path $noPremiseDir 'README.md')) }
    Assert-B3bCheck -Name 'bundle-marks-a-missing-premise' -Ok ($noPremiseReadme -match 'UNKNOWN \(verdict.json carried no premise\)') `
        -Detail ($noPremiseOut.Trim() + ' :: premise section=' + (@(($noPremiseReadme -split "`n") | Where-Object { $_ -match 'premise' }) -join ' '))

    $escapeDir = Join-Path ([System.IO.Path]::GetTempPath()) ('b3b-escape-' + [guid]::NewGuid().ToString('N'))
    $escapeRun = Invoke-B3bEvidenceBuilder -RigRoot $bundleSession -BundleDir $escapeDir -RepoRoot $work
    $escapeRefused = ($escapeRun.Exit -ne 0) -or ($escapeRun.Failure -match 'outside the repository root')
    Assert-B3bCheck -Name 'bundle-refuses-to-write-outside-the-repo-root' -Ok ($escapeRefused -and (-not (Test-Path -LiteralPath $escapeDir))) `
        -Detail ('exit=' + $escapeRun.Exit + ' failure=' + $escapeRun.Failure + ' output=' + $escapeRun.Output.Trim())
    if (Test-Path -LiteralPath $escapeDir) { Remove-Item -LiteralPath $escapeDir -Recurse -Force }

    # The four invariants the closure re-review found untested.
    # 1) The pre-cut liveness re-probe has to run on every path, including -NoConfirm: the wait
    #    that makes it necessary is the freeze delay, not the operator turn.
    $roundAst = [System.Management.Automation.Language.Parser]::ParseFile((Join-Path $B3bRigDir 'Invoke-B3bRound.ps1'), [ref]$null, [ref]$null)
    $probeWrites = @($roundAst.FindAll({ param($node)
                ($node -is [System.Management.Automation.Language.CommandAst]) -and
                ($node.GetCommandName() -eq 'Write-B3bArtifactText') -and
                ($node.Extent.Text -match 'survival-probe-2')
            }, $true))
    $guardedProbes = 0
    foreach ($write in $probeWrites) {
        $ancestor = $write
        while ($null -ne $ancestor.Parent) {
            $ancestor = $ancestor.Parent
            if (($ancestor -is [System.Management.Automation.Language.IfStatementAst]) -and ($ancestor.Extent.Text -match 'NoConfirm')) {
                $guardedProbes++
                break
            }
        }
    }
    Assert-B3bCheck -Name 'second-cut-probe-runs-on-every-path' -Ok (($probeWrites.Count -eq 1) -and ($guardedProbes -eq 0)) `
        -Detail ('writes=' + $probeWrites.Count + ' inside-NoConfirm=' + $guardedProbes)
    # The order matters as much as the presence: the freeze delay has to come *before* the probe, so
    # that nothing but the probe itself runs between the re-check and the cut. Pinned statically,
    # because a future edit could move the sleep back between the two without any other check going
    # red (that is exactly the regression the closure review found once already).
    $sleepCommands = @($roundAst.FindAll({ param($node)
                ($node -is [System.Management.Automation.Language.CommandAst]) -and
                ($node.GetCommandName() -eq 'Start-Sleep') -and
                ($node.Extent.Text -match 'PlanCutDelayS')
            }, $true))
    $orderOk = $false
    if (($sleepCommands.Count -eq 1) -and ($probeWrites.Count -eq 1)) {
        $orderOk = ($sleepCommands[0].Extent.StartLineNumber -lt $probeWrites[0].Extent.StartLineNumber)
    }
    Assert-B3bCheck -Name 'freeze-delay-precedes-the-cut-probe' -Ok $orderOk `
        -Detail ('sleep line=' + (@($sleepCommands | ForEach-Object { $_.Extent.StartLineNumber }) -join ',') +
            ' probe line=' + (@($probeWrites | ForEach-Object { $_.Extent.StartLineNumber }) -join ','))
    # 2) The preflight and the matrix drive the same VM, so they have to exclude each other.
    Assert-B3bCheck -Name 'preflight-takes-the-session-lock' -Ok ($sessionSource -match 'Enter-B3bSessionLock') `
        -Detail 'the preflight must take the same session lock as the matrix'
    Assert-B3bCheck -Name 'preflight-releases-the-lock-on-a-terminating-error' -Ok ($sessionSource -match 'trap \{' -and $sessionSource -match 'Exit-B3bSessionLock -SessionDir \$SessionDir') `
        -Detail 'an aborted preflight must not leave a live-pid lock behind'
    # 3) Clearing a stale session must not be able to kill another VM whose name starts with ours.
    $matchCases = @(
        @{ Line = 'VBoxHeadless --startvm Win11B3'; Want = $true },
        @{ Line = 'VBoxHeadless --startvm 02205fb0-1933-48ba-bbdb-a2c06b4eb9b2'; Want = $true },
        @{ Line = 'VBoxHeadless --startvm Win11B30'; Want = $false },
        @{ Line = 'VBoxHeadless --startvm OtherVm'; Want = $false },
        @{ Line = ''; Want = $false }
    )
    $matchFailures = New-Object System.Collections.Generic.List[string]
    foreach ($case in $matchCases) {
        $got = Test-B3bVmCommandLineMatch -CommandLine $case.Line -VmName 'Win11B3' -VmId '02205fb0-1933-48ba-bbdb-a2c06b4eb9b2'
        if ($got -ne $case.Want) { $matchFailures.Add($case.Line + ' => ' + $got) }
    }
    Assert-B3bCheck -Name 'stale-session-match-is-boundary-anchored' -Ok ($matchFailures.Count -eq 0) -Detail (($matchFailures) -join ' | ')
    # 4) The lock blocks a live holder of this directory but not a lock that was copied here
    #    from another session directory.
    $lockDir = Join-Path $work 'lock-fingerprint'
    [void](New-Item -ItemType Directory -Path $lockDir -Force)
    $firstLock = Enter-B3bSessionLock -SessionDir $lockDir -Owner 'self-test'
    $secondLock = Enter-B3bSessionLock -SessionDir $lockDir -Owner 'self-test-second'
    Assert-B3bCheck -Name 'lock-blocks-a-second-live-holder' -Ok ($firstLock.ok -and (-not $secondLock.ok)) `
        -Detail ('first=' + $firstLock.ok + ' second=' + $secondLock.ok + ' reason=' + $secondLock.reason)
    $foreignDir = Join-Path $work 'lock-foreign'
    [void](New-Item -ItemType Directory -Path $foreignDir -Force)
    $foreignDoc = [ordered]@{
        schema = 'b3b-session-lock'; v = 1; pid = $PID; owner = 'copied-from-elsewhere'
        sessionDir = 'D:\somewhere\else'; host = $env:COMPUTERNAME; acquiredAt = 'copied'
    }
    Write-B3aTextFile -Path (Join-Path $foreignDir 'session.lock') -Text ($foreignDoc | ConvertTo-Json -Compress) -WithBom
    $afterCopy = Enter-B3bSessionLock -SessionDir $foreignDir -Owner 'self-test-third'
    Assert-B3bCheck -Name 'inherited-lock-does-not-block-a-copied-session' -Ok ($afterCopy.ok) -Detail $afterCopy.reason
} finally {
    if (Test-Path -LiteralPath $work) { Remove-Item -LiteralPath $work -Recurse -Force }
}

Write-Output ('checks={0} failures={1}' -f $script:checks, $script:failures)
if ($script:failures -gt 0) { exit 1 }
exit 0

# B3b verdict module: the aggregation that turns per-round results into the
# three-tier durability verdict. Kept free of a top-level param block on purpose:
# it is dot-sourced by the CLI wrapper and by the rig self-test, and a top-level
# param() would silently rebind same-named caller variables (B3a defect D1).
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')

function Get-B3bU4Verdict {
    # U4 asks whether directory entries can be lost out of order under a power cut:
    # a complete request record followed by a completion record stopped inside its
    # target stage. Only a counterexample decides this question - a clean sample can
    # never prove ordering, so the best outcome here is unresolved.
    param(
        [Parameter(Mandatory = $true)]$Rounds,
        # Mandatory with no default on purpose: PSScriptAnalyzer forbids a default
        # for a mandatory parameter, and the caller has to state the sample size it
        # is judging rather than inherit one silently.
        [Parameter(Mandatory = $true)][int]$RequiredN
    )
    $valid = 0
    $orphans = New-Object System.Collections.Generic.List[object]
    foreach ($round in @($Rounds)) {
        if (-not $round.valid) { continue }
        if (-not $round.auditOk -or -not $round.inventoryRepeatableOk -or -not $round.journalChainOk) { continue }
        $valid++
        if ($round.orphanCompletion) { $orphans.Add($round) }
    }
    $verdict = 'unresolved'
    if ($orphans.Count -gt 0) { $verdict = 'disproven' }
    return [ordered]@{
        verdict      = $verdict
        validRounds  = $valid
        orphans      = $orphans.Count
        requiredN    = $RequiredN
        note         = 'ordering can only be falsified by a counterexample, never proven by a clean sample'
    }
}

function Invoke-B3bVerdict {
    # Group the rounds by (arm, stage), judge every cell, then aggregate per arm.
    param(
        [Parameter(Mandatory = $true)]$Rounds,
        [Parameter(Mandatory = $true)][int]$RequiredN,
        [Parameter(Mandatory = $true)]$Discrimination,
        $U4 = $null,
        [string]$WriterSha256 = '',
        [string]$Premise = ''
    )
    if ($RequiredN -lt 5) { throw 'Invoke-B3bVerdict: RequiredN must be >= 5 (design section 6)' }
    if ($null -eq $Discrimination) { throw 'Invoke-B3bVerdict: the discrimination gate result is required' }
    $groups = @{}
    foreach ($round in @($Rounds)) {
        $key = ([string]$round.arm + '|' + [string]$round.stage)
        if (-not $groups.ContainsKey($key)) { $groups[$key] = New-Object System.Collections.Generic.List[object] }
        $groups[$key].Add($round)
    }
    $gate = [bool]$Discrimination.ok
    $cells = New-Object System.Collections.Generic.List[object]
    foreach ($key in (@($groups.Keys) | Sort-Object)) {
        $parts = $key -split '\|'
        $cells.Add((Get-B3bCellVerdict -Arm $parts[0] -Stage $parts[1] -Rounds $groups[$key].ToArray() `
                    -RequiredN $RequiredN -DiscriminationGate $gate))
    }
    return (Get-B3bVerdictDocument -Cells $cells.ToArray() -RequiredN $RequiredN -Discrimination $Discrimination `
            -U4 $U4 -WriterSha256 $WriterSha256 -Premise $Premise)
}

# Get-B3aReconcile.ps1 - outcome classification MODULE for one B3a round (design S2.4).
#
# This file defines Invoke-B3aReconcile and is dot-sourced by Invoke-B3aRound.ps1.
# It deliberately declares NO top-level parameters: dot-sourcing a script whose top
# level has a param() block rebinds the caller's variables of the same name, which
# silently replaced the round script's command-line -Round/-Candidate with this
# file's defaults (found 2026-09-20 - every round ran as round 0/negative-control).
# The standalone CLI wrapper is Invoke-B3aReconcile.ps1.
#
# Reads plan.json + inventory JSON(s) + audit.txt, classifies the target file with
# the machine-checkable rules of design S2.4, writes reconcile.json, and appends the
# 'reconcile' journal record (ACK-gated). Comparison is on content and size only -
# guest clocks are never compared.
#
# Void mode: pass -VoidReason to record a void round (outcome=null, verdict=void,
# voidReason set); plan/inventory/audit paths are then optional.
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'b3a-journal-lib.ps1')

function Invoke-B3aReconcile {
    param(
        [string]$PlanJsonPath = '',
        [string]$InventoryJsonPath = '',
        [string]$InventoryJson2Path = '',
        [string]$AuditTxtPath = '',
        [Parameter(Mandatory = $true)][string]$OutputPath,
        [string]$JournalPath = '',
        [int]$Round = 0,
        [string]$Candidate = 'negative-control',
        [string]$Stage = '',
        [string]$TargetPath = '',
        [string]$ContentDigest = '',
        [long]$ContentSize = 0,
        [string]$VoidReason = $null
    )
    $result = [pscustomobject]@{ ok = $false; outcome = $null; verdict = ''; voidReason = $null; error = $null }
    try {
        if (-not [string]::IsNullOrEmpty($VoidReason)) {
            # round context from explicit parameters, falling back to plan.json
            $targetPath = $TargetPath; $digest = $ContentDigest; $size = $ContentSize
            $stage = $Stage
            if ($PlanJsonPath -and (Test-Path -LiteralPath $PlanJsonPath)) {
                $plan = (Get-B3aFileTextUtf8 -Path $PlanJsonPath) | ConvertFrom-Json
                if (-not $targetPath) { $targetPath = $plan.targetPath }
                if (-not $digest) { $digest = $plan.contentDigest }
                if ($size -eq 0) { $size = [long]$plan.contentSize }
                if (-not $stage) { $stage = $plan.stage }
            }
            $reconcile = [ordered]@{
                schema = 'b3a-reconcile'; v = 1
                round = $Round; candidate = $Candidate
                targetPath = $targetPath; contentDigest = $digest; contentSize = $size
                outcome = $null
                rule = 'void'
                inventoryDigests = @(); inventoryIdentical = $false
                auditDigest = $null; auditVerdict = $null; dirty = $null
                verdict = 'void'; voidReason = $VoidReason
                hostTs = (Get-B3aHostTimestamp)
            }
            Write-B3aTextFile -Path $OutputPath -Text ($reconcile | ConvertTo-Json -Compress -Depth 20)
            $result.outcome = $null
            $result.verdict = 'void'
            $result.voidReason = $VoidReason
        } else {
            if (-not $PlanJsonPath -or -not (Test-Path -LiteralPath $PlanJsonPath)) { throw 'plan.json missing' }
            if (-not $InventoryJsonPath -or -not (Test-Path -LiteralPath $InventoryJsonPath)) { throw 'inventory-1.json missing' }
            $plan = (Get-B3aFileTextUtf8 -Path $PlanJsonPath) | ConvertFrom-Json
            $inv1 = (Get-B3aFileTextUtf8 -Path $InventoryJsonPath) | ConvertFrom-Json
            $targetLeaf = [System.IO.Path]::GetFileName([string]$plan.targetPath)
            $entry = $null
            if ($inv1.roundRootExists) {
                foreach ($f in @($inv1.files)) {
                    if ([string]$f.path -eq $targetLeaf) { $entry = $f; break }
                }
            }
            $present = ($null -ne $entry)
            $actualSize = if ($present) { [long]$entry.size } else { 0 }
            $actualDigest = if ($present) { [string]$entry.sha256 } else { '' }
            $expectedSize = [long]$plan.contentSize
            $expectedDigest = [string]$plan.contentDigest
            $zeroDigest = Get-B3aZeroDigest -Size $expectedSize
            $outcome = Get-B3aOutcome -Present $present -ActualSize $actualSize -ActualDigest $actualDigest `
                -ExpectedSize $expectedSize -ExpectedDigest $expectedDigest -ZeroDigest $zeroDigest
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
                $auditText = [System.IO.File]::ReadAllText($AuditTxtPath)
                $auditDigest = Get-B3aSha256Hex -Path $AuditTxtPath
                $auditVerdict = (Read-B3aAuditVerdict -Text $auditText).verdict
            }
            $dirty = $null
            if ($inv1.PSObject.Properties.Name -contains 'dirty') { $dirty = $inv1.dirty }
            $rule = 'outcome per design S2.4 (content+size only)'
            $reconcile = [ordered]@{
                schema = 'b3a-reconcile'; v = 1
                round = if ($plan.PSObject.Properties.Name -contains 'round') { $plan.round } else { $Round }
                candidate = if ($plan.PSObject.Properties.Name -contains 'candidate') { $plan.candidate } else { $Candidate }
                targetPath = $plan.targetPath
                contentDigest = $plan.contentDigest
                contentSize = $plan.contentSize
                outcome = $outcome
                rule = $rule
                inventoryDigests = $inventoryDigests
                inventoryIdentical = $inventoryIdentical
                auditDigest = $auditDigest
                auditVerdict = $auditVerdict
                dirty = $dirty
                verdict = 'ok'
                voidReason = $null
                hostTs = (Get-B3aHostTimestamp)
            }
            Write-B3aTextFile -Path $OutputPath -Text ($reconcile | ConvertTo-Json -Compress -Depth 20)
            $result.outcome = $outcome
            $result.verdict = 'ok'
        }
        # journal record (design S2.4: "reconcile.json + reconcile journal record")
        if ($JournalPath) {
            $recStage = $Stage
            if (-not $recStage -and $PlanJsonPath -and (Test-Path -LiteralPath $PlanJsonPath)) {
                $p = (Get-B3aFileTextUtf8 -Path $PlanJsonPath) | ConvertFrom-Json
                if ($p.PSObject.Properties.Name -contains 'stage') { $recStage = [string]$p.stage }
            }
            if (-not $recStage) { $recStage = 'reconcile' }
            $extra = [ordered]@{
                outcome = $result.outcome
                inventoryDigests = @($reconcile.inventoryDigests)
                auditDigest = $reconcile.auditDigest
                verdict = $reconcile.verdict
                voidReason = $reconcile.voidReason
            }
            $rec = New-B3aJournalRecord -Round $Round -Candidate $Candidate -Stage $recStage `
                -TargetPath ([string]$reconcile.targetPath) `
                -ContentDigest ([string]$reconcile.contentDigest) `
                -ContentSize ([long]$reconcile.contentSize) `
                -PlanCutDelayS 2.0 -Op 'reconcile' -Extra $extra
            $ack = Add-B3aJournalRecord -JournalPath $JournalPath -Record $rec
            if (-not $ack.ack) { throw ('reconcile journal ACK failed: ' + $ack.error) }
        }
        $result.ok = $true
        return $result
    } catch {
        $result.error = $_.Exception.Message
        return $result
    }
}

# Module only: the standalone CLI entry point lives in Invoke-B3aReconcile.ps1 so
# that this file stays free of top-level parameters (see the header note).

# Invoke-B3aReconcile.ps1 - standalone CLI wrapper for the reconcile module.
#
# Usage (all switches are optional except -OutputPath):
#   powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aReconcile.ps1 `
#     -OutputPath <path> -JournalPath <path> -Round <n> -Candidate <name> [-Stage <s>]
#     [-TargetPath <p>] [-ContentDigest <hex>] [-ContentSize <n>] [-VoidReason <text>]
#     [-PlanJsonPath <p>] [-InventoryJsonPath <p>] [-InventoryJson2Path <p>] [-AuditTxtPath <p>]
#
# Exit 0 = reconcile recorded; exit 1 = required input missing, invalid, or the
# journal record was not acknowledged.
#
# The implementation module is Get-B3aReconcile.ps1, which holds no top-level
# parameters so it can be dot-sourced safely. Do NOT dot-source this wrapper: its
# parameter names would rebind the caller's variables (2026-09-20).
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

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

$ErrorActionPreference = 'Stop'

. (Join-Path $PSScriptRoot 'Get-B3aReconcile.ps1')

$r = Invoke-B3aReconcile -PlanJsonPath $PlanJsonPath -InventoryJsonPath $InventoryJsonPath `
    -InventoryJson2Path $InventoryJson2Path -AuditTxtPath $AuditTxtPath -OutputPath $OutputPath `
    -JournalPath $JournalPath -Round $Round -Candidate $Candidate -Stage $Stage `
    -TargetPath $TargetPath -ContentDigest $ContentDigest -ContentSize $ContentSize `
    -VoidReason $VoidReason
if ($r.ok) {
    'reconcile: verdict={0} outcome={1}' -f $r.verdict, $(if ($r.outcome) { $r.outcome } else { 'n/a' })
    exit 0
}
'reconcile: FAILED - ' + $r.error
exit 1
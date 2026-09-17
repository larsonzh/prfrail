# B2 close-out helper: delete the AppContainer profile that apply-b2-loopback.ps1 created.
# Deletion only: this script never creates machine state, and it refuses to run while the profile's loopback
# exemption is still present (revert first, then delete).
# Safe to run unelevated; it is idempotent and reports whether the profile still exists afterwards.
[CmdletBinding()]
param(
    [string]$ContainerName = 'prfrail-b2-battery'
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

$exempt = Get-B2ExemptionList
$profileExists = Get-B2ContainerProfileExists -ContainerName $ContainerName

if (-not $profileExists) {
    Show 'profileExists' $false
    Show 'action' 'nothing to do'
    exit 0
}

$sid = Get-B2ContainerSidByName -ContainerName $ContainerName
Show 'containerName' $ContainerName
Show 'containerSid' $sid

if (@($exempt.sids) -contains $sid) {
    throw ('refusing to delete the profile while its loopback exemption is still present; run revert-b2-loopback.ps1 first (SID {0})' -f $sid)
}

$hr = [B2LoopbackNative]::DeleteAppContainerProfile($ContainerName)
if ($hr -ne 0) { throw ('DeleteAppContainerProfile failed 0x{0:X8}' -f $hr) }

$stillThere = Get-B2ContainerProfileExists -ContainerName $ContainerName
Show 'profileExistsAfterDelete' $stillThere
if ($stillThere) { throw 'the profile still exists after deletion' }
Show 'result' 'deleted'
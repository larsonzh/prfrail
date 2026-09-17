# One-off validation probe (scratch, disclosed): proves the registry-backed profile-existence check on a
# throwaway profile name, then removes every trace. It never touches the B2 container or its exemption.
$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot '..\..\tools\agent-probe\appcontainer-b2\b2-loopback-lib.ps1')

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

$name = 'prfrail-b2-detect-probe'
Show 'beforeCreate' (Get-B2ContainerProfileExists -ContainerName $name)
$sid = New-B2ContainerSidIfMissing -ContainerName $name
Show 'createdSid' $sid
Show 'afterCreate' (Get-B2ContainerProfileExists -ContainerName $name)
Show 'deriveOnlySidMatches' ((Get-B2ContainerSidByName -ContainerName $name) -eq $sid)
$deleteResult = [B2LoopbackNative]::DeleteAppContainerProfile($name)
Show 'deleteHresult' ('0x{0:X8}' -f $deleteResult)
Show 'afterDelete' (Get-B2ContainerProfileExists -ContainerName $name)
try {
    $null = Get-B2ContainerSidExisting -ContainerName $name
    Show 'existingHelper' 'UNEXPECTED-SUCCESS'
} catch {
    Show 'existingHelper' 'throws-as-designed'
}
Show 'b2ProfileExists' (Get-B2ContainerProfileExists -ContainerName 'prfrail-b2-battery')
Show 'storageKeysRemaining' (@(Get-ChildItem 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\AppContainer\Storage' -ErrorAction SilentlyContinue | Where-Object { $_.PSChildName -like 'prfrail-b2*' } | Select-Object -ExpandProperty PSChildName) -join ',')

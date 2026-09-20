# Classify the guest's most recent power event.
# The rig calls this after every restart so each round carries machine-checkable proof
# that the last shutdown was the intended hard cut and not a licence-service shutdown.
$ErrorActionPreference = 'Continue'
$hours = 6
$since = (Get-Date).AddHours(-$hours)

$ev = Get-WinEvent -FilterHashtable @{LogName = 'System'; Id = 1074, 1076, 41, 6008; StartTime = $since } -ErrorAction SilentlyContinue |
    Sort-Object TimeCreated

'--- power events in the last {0} h (oldest first) ---' -f $hours
foreach ($e in $ev) {
    $first = ($e.Message -split "`r?`n")[0]
    '{0}  Id={1}  {2}' -f $e.TimeCreated.ToString('yyyy-MM-dd HH:mm:ss'), $e.Id, $first
}

'--- verdict ---'
$last = $ev | Select-Object -Last 1
if (-not $last) {
    'VERDICT: no shutdown/power event recorded in the window'
} elseif ($last.Message -match 'wlms') {
    'VERDICT: LICENCE-SHUTDOWN  (wlms.exe; void the round)'
} elseif ($last.Id -eq 41 -or $last.Id -eq 6008) {
    'VERDICT: HARD-POWER-LOSS  (expected for an injection round)'
} elseif ($last.Id -eq 1074) {
    'VERDICT: GRACEFUL-SHUTDOWN requested by another process (inspect the line above)'
} else {
    "VERDICT: other (Id=$($last.Id))"
}

'--- context ---'
"now      = $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
"lastBoot = $((Get-CimInstance Win32_OperatingSystem).LastBootUpTime)"
$up = (Get-Date) - (Get-CimInstance Win32_OperatingSystem).LastBootUpTime
"uptime   = $([math]::Round($up.TotalMinutes,1)) minutes"
$lic = Get-CimInstance -ClassName SoftwareLicensingProduct -Filter 'PartialProductKey is not null'
"licence  = status $($lic.LicenseStatus), grace $($lic.GracePeriodRemaining)"

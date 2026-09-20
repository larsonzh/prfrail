# Invoke-B3aCalibration.ps1 - R0 + R1..RN negative + R(N+1)..R(2N) positive and
# the calibration-verdict.json (design S4 / S4.5 / S5.2).
#
# Slots (N defaults to 5, design Q3): R0 parity (informative, not counted in N),
# then N negative-control rounds, then N positive-control rounds. Each slot is
# run by invoking Invoke-B3aRound.ps1 as a child process:
#   exit 0  -> valid round, outcome recorded
#   exit 3  -> void round (new round number, void budget 3 per control candidate)
#   exit 4  -> R0 survived: pause for the operator (design S4.4)
#   exit 5  -> operational fail (same round number, retry budget 2)
#   exit 6/7 -> session blocker / operator abort: stop, propagate
#
# Final verdict is computed only from the rounds produced by THIS invocation
# plus the whole journal chain check (design S4.5). Writes
# <RigRoot>\<date>\calibration-verdict.json.
#
# Exit codes: 0 CALIBRATED; 1 UNUSABLE; 3 void budget exceeded;
#             4 R0-survived pause; 5 fail budget exhausted; 6 blocker.
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
param(
    [string]$VmName = 'Win11B3',
    [string]$CredFile = 'D:\VirtualBox VMs\Win11B3\host-only\cred.txt',
    [string]$GuestUser = 'prfrail',
    [string]$RigRoot = 'D:\VirtualBox VMs\Win11B3\b3a-runs',
    [string]$VBoxManage = 'C:\Program Files\Oracle\VirtualBox\VBoxManage.exe',
    [string]$GuestToolDir = 'C:\prfrail-prep',
    [string]$GuestPrepDir = 'D:\VirtualBox VMs\Win11B3\guest-prep',
    [string]$ProbeRoot = 'R:\agent-runner-replay',
    # N is the acceptance denominator; the design rules (S4.2/S4.3) are stated at
    # N=5 and a verdict over fewer rounds could not support "at least one loss in N".
    [ValidateRange(5, [int]::MaxValue)][int]$N = 5,
    [int]$Size = 262144,
    [double]$PlanCutDelayS = 2.0,
    [int]$MaxReadyWaitS = 120,
    [int]$VoidBudget = 3,
    [int]$FailBudget = 2,
    [switch]$NoConfirm
)

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'b3a-journal-lib.ps1')

if ($MyInvocation.InvocationName -eq '.') { return }

# Fail closed: an unexpected error stops the session instead of skipping a gate.
$ErrorActionPreference = 'Stop'

$DateTag = Get-Date -Format 'yyyy-MM-dd'
$SessionDir = Join-Path $RigRoot $DateTag
$JournalPath = Join-Path $SessionDir 'journal.jsonl'
$SessionLog = Join-Path $SessionDir 'session.log'

function Write-B3aCalLog {
    param([Parameter(Mandatory = $true)][string]$Text)
    Add-B3aLogLine -Path $SessionLog -Text ('[{0}] {1}' -f (Get-Date -Format 'HH:mm:ss'), $Text)
    $Text
}

[void](New-Item -ItemType Directory -Path $SessionDir -Force)
[void](Write-B3aCalLog ('=== B3a calibration begin: N={0} (R0 + {0}x negative + {0}x positive) ===' -f $N))

$chain = Test-B3aJournalChain -JournalPath $JournalPath
if (-not $chain.ok) {
    [void](Write-B3aCalLog ('BLOCKER: journal chain broken at line {0}: {1}' -f $chain.firstBadLine, $chain.firstError))
    exit 6
}

# absolute round numbering: monotonic, voids also consume numbers (design S3.3)
$startRound = 0
if ($chain.count -gt 0) {
    $existing = Get-B3aJournalRecord -JournalPath $JournalPath
    foreach ($r in $existing) { if ([int]$r.round -gt $startRound) { $startRound = [int]$r.round } }
}
$abs = $startRound + 1
[void](Write-B3aCalLog ('journal chain ok ({0} records); next absolute round = {1}' -f $chain.count, $abs))

# slot plan: parity first, then negative 1..N, then positive 1..N (design S5.2)
$slots = New-Object System.Collections.Generic.List[object]
$slots.Add([ordered]@{ logical = 'R0'; candidate = 'parity' })
for ($i = 1; $i -le $N; $i++) { $slots.Add([ordered]@{ logical = ('R{0}' -f $i); candidate = 'negative-control' }) }
for ($i = 1; $i -le $N; $i++) { $slots.Add([ordered]@{ logical = ('R{0}' -f ($N + $i)); candidate = 'positive-control' }) }

$voidCounts = @{ 'parity' = 0; 'negative-control' = 0; 'positive-control' = 0 }
# design S3.3: "same round number, retry budget 2" - the budget belongs to a round
# number, so a flaky round cannot consume the budget of the rounds that follow it.
$failCounts = @{}
$runResults = New-Object System.Collections.Generic.List[object]

foreach ($slot in $slots) {
    $candidate = [string]$slot.candidate
    while ($true) {
        $childArgs = @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', (Join-Path $PSScriptRoot 'Invoke-B3aRound.ps1'),
            '-Round', ([string]$abs), '-Candidate', $candidate,
            '-VmName', $VmName, '-CredFile', (Get-B3aQuoted $CredFile),
            '-RigRoot', (Get-B3aQuoted $RigRoot), '-VBoxManage', (Get-B3aQuoted $VBoxManage),
            '-GuestToolDir', $GuestToolDir, '-GuestPrepDir', (Get-B3aQuoted $GuestPrepDir),
            '-ProbeRoot', $ProbeRoot, '-Size', ([string]$Size), '-GuestUser', $GuestUser,
            '-PlanCutDelayS', ([string]$PlanCutDelayS), '-MaxReadyWaitS', ([string]$MaxReadyWaitS))
        if ($NoConfirm) { $childArgs += '-NoConfirm' }
        $psExe = Join-Path $env:WINDIR 'System32\WindowsPowerShell\v1.0\powershell.exe'
        [void](Write-B3aCalLog ('launch: round {0} candidate={1} ({2})' -f $abs, $candidate, $slot.logical))
        $child = Invoke-B3aProcess -FileName $psExe -Arguments ($childArgs -join ' ')
        $code = $child.exitCode
        if ($code -eq 0) {
            $reconcilePath = Join-Path $SessionDir ('rounds\r{0:D2}\reconcile.json' -f $abs)
            if (-not (Test-Path -LiteralPath $reconcilePath)) {
                [void](Write-B3aCalLog ('BLOCKER: round {0} exited 0 but reconcile.json is missing' -f $abs))
                exit 6
            }
            $rec = (Get-B3aFileTextUtf8 -Path $reconcilePath) | ConvertFrom-Json
            $runResults.Add([ordered]@{
                round = $abs; candidate = $candidate; logical = $slot.logical
                outcome = $rec.outcome; verdict = $rec.verdict; voidReason = $rec.voidReason
            })
            [void](Write-B3aCalLog ('round {0} ({1}) valid: outcome={2}' -f $abs, $slot.logical, $rec.outcome))
            if ($candidate -eq 'parity' -and $rec.outcome -eq 'survived') {
                'PAUSE: R0 parity round SURVIVED (design S4.4). Operator decision required before starting the calibration proper.'
                [void](Write-B3aCalLog 'PAUSE: R0 survived')
                exit 4
            }
            $abs++
            break
        } elseif ($code -eq 3) {
            $voidCounts[$candidate]++
            [void](Write-B3aCalLog ('round {0} void ({1}/{2}): re-running with a new number' -f $abs, $voidCounts[$candidate], $VoidBudget))
            if ($voidCounts[$candidate] -gt $VoidBudget) {
                [void](Write-B3aCalLog ('CALIBRATION FAILED: void budget exceeded for {0}' -f $candidate))
                exit 3
            }
            $abs++
        } elseif ($code -eq 5) {
            $key = [string]$abs
            if (-not $failCounts.ContainsKey($key)) { $failCounts[$key] = 0 }
            $failCounts[$key]++
            # log the child's own output BEFORE the budget check: on the exhausted
            # path this is the only record of the final failing attempt (2026-09-20)
            [void](Write-B3aCalLog ('round {0} child stderr: {1}' -f $abs, (Get-B3aBoundedText -Text $child.stderrText)))
            [void](Write-B3aCalLog ('round {0} child stdout: {1}' -f $abs, (Get-B3aBoundedText -Text $child.stdoutText)))
            [void](Write-B3aCalLog ('round {0} operational fail ({1}/{2} retries used): retrying same number' -f $abs, $failCounts[$key], $FailBudget))
            if ($failCounts[$key] -gt $FailBudget) {
                [void](Write-B3aCalLog ('CALIBRATION STOPPED: retry budget exhausted for round {0} (escalate to operator)' -f $abs))
                exit 5
            }
        } elseif ($code -eq 4) {
            exit 4
        } elseif ($code -eq 7) {
            [void](Write-B3aCalLog 'operator declined a cut; calibration stopped')
            exit 7
        } else {
            [void](Write-B3aCalLog ('BLOCKER: unexpected round exit code ' + $code))
            [void](Write-B3aCalLog ('round {0} child stderr: {1}' -f $abs, (Get-B3aBoundedText -Text $child.stderrText)))
            [void](Write-B3aCalLog ('round {0} child stdout: {1}' -f $abs, (Get-B3aBoundedText -Text $child.stdoutText)))
            exit 6
        }
    }
}

# ---- verdict from the rounds of THIS invocation + the whole journal chain ----
$negValid = 0; $negLosses = 0; $posValid = 0; $posLosses = 0
$inventoryOk = $true; $auditOk = $true
foreach ($r in $runResults) {
    $outcome = [string]$r.outcome
    if ($r.candidate -eq 'negative-control' -and $r.verdict -eq 'ok') {
        $negValid++
        if ($outcome -like 'lost-*') { $negLosses++ }
    }
    if ($r.candidate -eq 'positive-control' -and $r.verdict -eq 'ok') {
        $posValid++
        if ($outcome -like 'lost-*') { $posLosses++ }
    }
    if ($r.verdict -eq 'ok') {
        $recPath = Join-Path $SessionDir ('rounds\r{0:D2}\reconcile.json' -f $r.round)
        $rec = (Get-B3aFileTextUtf8 -Path $recPath) | ConvertFrom-Json
        if (-not $rec.inventoryIdentical) { $inventoryOk = $false }
        if ([string]$rec.auditVerdict -ne 'HARD-POWER-LOSS') { $auditOk = $false }
    }
}
$chain2 = Test-B3aJournalChain -JournalPath $JournalPath
$verdict = Get-B3aCalibrationVerdict -JournalChainOk $chain2.ok `
    -NegativeValidRounds $negValid -NegativeLosses $negLosses `
    -PositiveValidRounds $posValid -PositiveLosses $posLosses `
    -InventoryRepeatabilityOk $inventoryOk -AuditVerdictsOk $auditOk -RequiredN $N

$verdictJson = [ordered]@{}
foreach ($k in $verdict.Keys) { $verdictJson[$k] = $verdict[$k] }
$verdictJson['hostTs'] = Get-B3aHostTimestamp
$verdictJson['journalRecordCount'] = $chain2.count
# .ToArray() rather than @($runResults): see the binder note in Invoke-B3aToolDeploy
# (PS 5.1 fails to convert a List[object] through @() in this position).
$verdictJson['rounds'] = $runResults.ToArray()
$verdictPath = Join-Path $SessionDir 'calibration-verdict.json'
Write-B3aTextFile -Path $verdictPath -Text ($verdictJson | ConvertTo-Json -Compress -Depth 20)

[void](Write-B3aCalLog ('=== calibration verdict: {0} (neg {1}/{2} losses, pos {3}/{4} losses) ===' -f `
    $verdict.gate, $negLosses, $negValid, $posLosses, $posValid))
'verdict: {0}' -f $verdict.gate
'verdict file: {0}' -f $verdictPath
'rounds:'
foreach ($r in $runResults) {
    '  round {0} {1} ({2}): outcome={3} verdict={4} voidReason={5}' -f `
        $r.round, $r.candidate, $r.logical, $r.outcome, $r.verdict, $r.voidReason
}
if ($verdict.gate -eq 'CALIBRATED') { exit 0 } else { exit 1 }

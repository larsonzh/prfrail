# Test-B3aRig.ps1 - hermetic self-tests for the B3a rig (design S2.5 / S7).
#
# Hermetic: uses temporary fixture directories only and never touches the VM.
# Guest scripts are executed through Windows PowerShell 5.1 (powershell.exe)
# exactly as the rig does. Final line: "checks=N failures=0"; exit 0 only when
# failures = 0.
#
# Covered:
#   (a) tampered journal line detected, chain verify goes red
#   (b) failed/missing ACK is a hard failure (never proceeds; Assert-B3aAck throws)
#   (c) inventory script deterministic on a fixture tree (two runs byte-identical)
#   (d) outcome classifier returns the right one of survived / lost-absent /
#       lost-zero-filled / lost-torn-size / lost-torn-content (fixture cases)
#   (e) calibration verdict UNUSABLE when a control misbehaves, CALIBRATED when
#       both controls pass (fixture-driven, no VM)
#   plus: every rig script parses clean and carries BOM+LF; the negative write
#   probe reports its digest and the host recomputation matches.
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')
. (Join-Path $PSScriptRoot 'b3a-journal-lib.ps1')

# Fail closed: an unexpected error aborts the run instead of dropping checks.
$ErrorActionPreference = 'Stop'

$script:checks = 0
$script:failures = 0

function Assert-B3aCheck {
    param([Parameter(Mandatory = $true)][string]$Name,
          [Parameter(Mandatory = $true)][bool]$Ok,
          [string]$Detail = '')
    $script:checks++
    if ($Ok) { 'PASS  {0}' -f $Name }
    else {
        $script:failures++
        'FAIL  {0}  {1}' -f $Name, $Detail
    }
}

function Invoke-B3aTestGuest {
    param([Parameter(Mandatory = $true)][string]$ScriptName,
          [string[]]$argList = @())
    $psExe = Join-Path $env:WINDIR 'System32\WindowsPowerShell\v1.0\powershell.exe'
    $quoted = @()
    foreach ($a in $argList) { $quoted += (Get-B3aQuoted $a) }
    $argStr = '-NoProfile -ExecutionPolicy Bypass -File "' + (Join-Path $PSScriptRoot ('guest\' + $ScriptName)) + '"' + $(if ($quoted.Count) { ' ' + ($quoted -join ' ') } else { '' })
    return Invoke-B3aProcess -FileName $psExe -Arguments $argStr -TimeoutS 120
}

$FixtureRoot = Join-Path $env:TEMP ('b3a-rig-test-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $FixtureRoot -Force)

try {
    # ---- script parse + encoding ------------------------------------------
    $rigScripts = @(
        'b3a-rig-lib.ps1', 'b3a-journal-lib.ps1', 'New-B3aSession.ps1',
        'Invoke-B3aRound.ps1', 'Invoke-B3aCalibration.ps1', 'Get-B3aReconcile.ps1',
        'Invoke-B3aReconcile.ps1', 'New-B3aEvidence.ps1', 'Test-B3aRig.ps1',
        'guest\12-round-inventory.ps1', 'guest\13-round-write.ps1'
    )
    $parseOk = $true
    $parseDetail = ''
    foreach ($s in $rigScripts) {
        $path = Join-Path $PSScriptRoot $s
        $tokens = $null; $errors = $null
        [void][System.Management.Automation.Language.Parser]::ParseFile($path, [ref]$tokens, [ref]$errors)
        if ($null -ne $errors -and $errors.Count -gt 0) {
            $parseOk = $false
            $parseDetail += ($s + ': ' + $errors[0].Message + '; ')
        }
    }
    Assert-B3aCheck -Name 'parse-all-scripts' -Ok $parseOk -Detail $parseDetail

    $encOk = $true
    $encDetail = ''
    foreach ($s in $rigScripts) {
        if (-not (Test-B3aPs1Encoding -Path (Join-Path $PSScriptRoot $s))) {
            $encOk = $false; $encDetail += ($s + '; ')
        }
    }
    Assert-B3aCheck -Name 'ps1-encoding-bom-lf' -Ok $encOk -Detail $encDetail

    # ---- (a) journal chain + tamper ----------------------------------------
    $journal = Join-Path $FixtureRoot 'journal.jsonl'
    $r0 = New-B3aJournalRecord -Round 0 -Candidate 'parity' -Stage 'parity' `
        -TargetPath 'R:\calib\r00\parity.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](1, 2, 3))) `
        -ContentSize 3 -PlanCutDelayS 2.0 -Op 'plan'
    [void](Add-B3aJournalRecord -JournalPath $journal -Record $r0)
    $r1 = New-B3aJournalRecord -Round 1 -Candidate 'negative-control' -Stage 'temp-written' `
        -TargetPath 'R:\calib\r01\neg.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](4, 5, 6))) `
        -ContentSize 3 -PlanCutDelayS 2.0 -Op 'plan'
    [void](Add-B3aJournalRecord -JournalPath $journal -Record $r1)
    $r1c = New-B3aJournalRecord -Round 1 -Candidate 'negative-control' -Stage 'temp-written' `
        -TargetPath 'R:\calib\r01\neg.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](4, 5, 6))) `
        -ContentSize 3 -PlanCutDelayS 2.0 -Op 'cut' -Extra ([ordered]@{ cutExitCode = 0; cutHostTs = (Get-B3aHostTimestamp) })
    [void](Add-B3aJournalRecord -JournalPath $journal -Record $r1c)

    $chain1 = Test-B3aJournalChain -JournalPath $journal
    $firstPrevNull = ($null -eq $chain1.records[0].prevHash)
    Assert-B3aCheck -Name 'journal-chain-ok' -Ok ($chain1.ok -and $chain1.count -eq 3 -and $firstPrevNull) `
        -Detail ('count=' + $chain1.count + ' err=' + $chain1.firstError)

    # tamper: change a digest char in line 2 -> recordHash mismatch must go red
    $bytes = [System.IO.File]::ReadAllBytes($journal)
    $text = [System.Text.Encoding]::UTF8.GetString($bytes)
    $lines = $text -split "`n"
    $lines[1] = $lines[1].Replace('0', '1')
    [System.IO.File]::WriteAllBytes($journal, [System.Text.Encoding]::UTF8.GetBytes(($lines -join "`n")))
    $chainTampered = Test-B3aJournalChain -JournalPath $journal
    Assert-B3aCheck -Name 'journal-tamper-red' -Ok ((-not $chainTampered.ok) -and $chainTampered.firstBadLine -eq 2) `
        -Detail ('badLine=' + $chainTampered.firstBadLine + ' err=' + $chainTampered.firstError)

    # ---- (b) ACK ok / ACK failure blocks -----------------------------------
    $journal2 = Join-Path $FixtureRoot 'journal-ack.jsonl'
    $recA = New-B3aJournalRecord -Round 0 -Candidate 'parity' -Stage 'parity' `
        -TargetPath 'R:\calib\r00\parity.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](7))) `
        -ContentSize 1 -PlanCutDelayS 2.0 -Op 'plan'
    $ack1 = Add-B3aJournalRecord -JournalPath $journal2 -Record $recA
    $chain2 = Test-B3aJournalChain -JournalPath $journal2
    Assert-B3aCheck -Name 'journal-ack-ok' -Ok ($ack1.ack -and $chain2.ok -and $chain2.count -eq 1) -Detail $ack1.error

    # append failure (read-only file) -> ack=false -> Assert-B3aAck throws -> BLOCKED
    [System.IO.File]::SetAttributes($journal2, [System.IO.FileAttributes]::ReadOnly)
    $ackFail = Add-B3aJournalRecord -JournalPath $journal2 -Record `
        (New-B3aJournalRecord -Round 1 -Candidate 'negative-control' -Stage 'temp-written' `
            -TargetPath 'R:\calib\r01\neg.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](8))) `
            -ContentSize 1 -PlanCutDelayS 2.0 -Op 'plan')
    [System.IO.File]::SetAttributes($journal2, [System.IO.FileAttributes]::Normal)
    $threw = $false
    try { [void](Assert-B3aAck -AckResult $ackFail -Context 'test') } catch { $threw = $true }
    $gate = if ($ackFail.ack) { 'PROCEED' } else { 'BLOCKED' }
    Assert-B3aCheck -Name 'journal-ack-fail-blocks' -Ok ((-not $ackFail.ack) -and $threw -and $gate -eq 'BLOCKED') `
        -Detail $ackFail.error

    # ---- (b2) timestamp round-trip fixpoint (2026-09-20 intermittent root cause) --
    # The verifier must hash the stored line text, not a re-serialization: an ISO
    # timestamp whose 7th fractional digit is a trailing zero comes back from
    # ConvertFrom-Json as [datetime] and re-serializes one digit shorter. Both the
    # append path (tail verify) and the chain verify must survive that.
    $journalTs = Join-Path $FixtureRoot 'journal-ts.jsonl'
    $ackTs1 = Add-B3aJournalRecord -JournalPath $journalTs -Record `
        (New-B3aJournalRecord -Round 0 -Candidate 'parity' -Stage 'parity' `
            -TargetPath 'R:\calib\r00\parity.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](9))) `
            -ContentSize 1 -PlanCutDelayS 2.0 -Op 'plan' -HostTs '2026-09-20T03:35:28.8416210+08:00')
    $ackTs2 = Add-B3aJournalRecord -JournalPath $journalTs -Record `
        (New-B3aJournalRecord -Round 1 -Candidate 'negative-control' -Stage 'temp-written' `
            -TargetPath 'R:\calib\r01\neg.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](10))) `
            -ContentSize 1 -PlanCutDelayS 2.0 -Op 'cut' -HostTs '2026-09-20T03:35:28.8000000+08:00' `
            -Extra ([ordered]@{ cutExitCode = 0 }))
    $chainTs = Test-B3aJournalChain -JournalPath $journalTs
    Assert-B3aCheck -Name 'journal-timestamp-roundtrip' `
        -Ok ($ackTs1.ack -and $ackTs2.ack -and $chainTs.ok -and $chainTs.count -eq 2) `
        -Detail ('ack1=' + $ackTs1.ack + ' ack2=' + $ackTs2.ack + ' chain=' + $chainTs.ok + ' err=' + $chainTs.firstError)

    # appended bytes after the recordHash field sit outside the hashed text, so the
    # verifier must also enforce the canonical line shape; without that check this
    # tamper leaves the chain green (independent-review counterexample, 2026-09-20)
    $journalTail = Join-Path $FixtureRoot 'journal-tail.jsonl'
    [void](Add-B3aJournalRecord -JournalPath $journalTail -Record `
        (New-B3aJournalRecord -Round 0 -Candidate 'parity' -Stage 'parity' `
            -TargetPath 'R:\calib\r00\parity.bin' -ContentDigest (Get-B3aSha256Hex -Bytes ([byte[]](11))) `
            -ContentSize 1 -PlanCutDelayS 2.0 -Op 'plan'))
    $tailLines = @(Get-B3aRawJournalLine -JournalPath $journalTail)
    $tailLine = $tailLines[0]
    $tamperedLine = $tailLine.TrimEnd('}') + ',"tail":"tampered"}'
    [System.IO.File]::WriteAllBytes($journalTail, [System.Text.Encoding]::UTF8.GetBytes(($tamperedLine + "`n")))
    $chainTail = Test-B3aJournalChain -JournalPath $journalTail
    Assert-B3aCheck -Name 'journal-tail-append-red' `
        -Ok ((-not $chainTail.ok) -and $chainTail.firstError -eq 'stored-line-not-canonical') `
        -Detail ('ok=' + $chainTail.ok + ' err=' + $chainTail.firstError)

    # ---- (c) inventory determinism -----------------------------------------
    $roundRoot = Join-Path $FixtureRoot 'calib\r07'
    [void](New-Item -ItemType Directory -Path $roundRoot -Force)
    Write-B3aTextFile -Path (Join-Path $roundRoot 'neg.bin') -Text 'aaaa'
    [void](New-Item -ItemType Directory -Path (Join-Path $roundRoot 'sub') -Force)
    Write-B3aTextFile -Path (Join-Path $roundRoot 'sub\b.txt') -Text 'bbbbbbbb'
    $probeRoot = Join-Path $FixtureRoot 'agent-runner-replay'
    [void](New-Item -ItemType Directory -Path (Join-Path $probeRoot 'requests') -Force)
    [void](New-Item -ItemType Directory -Path (Join-Path $probeRoot 'completions') -Force)
    [void](New-Item -ItemType Directory -Path (Join-Path $probeRoot 'runs\req1\events') -Force)
    Write-B3aTextFile -Path (Join-Path $probeRoot 'requests\req1.jsonl') -Text '{}'
    Write-B3aTextFile -Path (Join-Path $probeRoot 'requests\req2.jsonl') -Text '{}'
    Write-B3aTextFile -Path (Join-Path $probeRoot 'completions\req2.jsonl') -Text '{}'
    Write-B3aTextFile -Path (Join-Path $probeRoot 'completions\req9.jsonl') -Text '{}'
    Write-B3aTextFile -Path (Join-Path $probeRoot 'ownership.json') -Text '{}'
    Write-B3aTextFile -Path (Join-Path $probeRoot 'stray.tmp') -Text 'x'
    Write-B3aTextFile -Path (Join-Path $probeRoot 'runs\req1\events\state-events.jsonl') -Text '[]'
    $invRun1 = Invoke-B3aTestGuest -ScriptName '12-round-inventory.ps1' -argList @('-RoundRoot', $roundRoot, '-ProbeRoot', $probeRoot)
    $invRun2 = Invoke-B3aTestGuest -ScriptName '12-round-inventory.ps1' -argList @('-RoundRoot', $roundRoot, '-ProbeRoot', $probeRoot)
    $identical = Test-B3aBytesEqual -A $invRun1.stdoutBytes -B $invRun2.stdoutBytes
    $invJson = Get-B3aJsonFromText -Text $invRun1.stdoutText
    $invGood = ($null -ne $invJson) -and (@($invJson.files).Count -eq 2) -and ($invJson.probeState -eq 'entries') `
        -and (@($invJson.probeEntries.summary.tmpFiles).Count -eq 1) `
        -and (@($invJson.probeEntries.summary.requestsWithoutCompletion) -contains 'req1') `
        -and (@($invJson.probeEntries.summary.orphanCompletions) -contains 'req9') `
        -and (@($invJson.probeEntries.summary.runDirs) -contains 'req1')
    Assert-B3aCheck -Name 'inventory-deterministic' -Ok ($identical -and $invGood -and $invRun1.exitCode -eq 0) `
        -Detail ('identical=' + $identical + ' json=' + ($null -ne $invJson))

    # ---- write probe (negative control) ------------------------------------
    $negPath = Join-Path $FixtureRoot 'calib\r03\neg.bin'
    [void](New-Item -ItemType Directory -Path (Split-Path -Parent $negPath) -Force)
    $wr = Invoke-B3aTestGuest -ScriptName '13-round-write.ps1' -argList @('-Mode', 'negative', '-Path', $negPath, '-Size', '4096', '-Seed', 'b3a-test-seed')
    $probeJson = Get-B3aJsonFromText -Text $wr.stdoutText
    $hostDigest = Get-B3aSeedPayloadDigest -Seed 'b3a-test-seed' -Size 4096
    $fileExists = (Test-Path -LiteralPath $negPath) -and ((Get-Item -LiteralPath $negPath).Length -eq 4096)
    Assert-B3aCheck -Name 'write-probe-negative' -Ok ($wr.exitCode -eq 0 -and $null -ne $probeJson `
        -and $probeJson.writeResult -eq 'ok' -and ([string]$probeJson.apisCalled -eq 'File.WriteAllBytes') `
        -and ([string]$probeJson.contentDigest -eq $hostDigest) -and $fileExists) `
        -Detail ('exit=' + $wr.exitCode + ' digestMatch=' + ([string]$probeJson.contentDigest -eq $hostDigest))

    # ---- (d) outcome classifier ---------------------------------------------
    $size = 262144
    $content = Get-B3aSeedPayload -Seed 'b3a-classifier' -Size $size
    $expectedDigest = Get-B3aSha256Hex -Bytes $content
    $zeroDigest = Get-B3aZeroDigest -Size $size
    $otherDigest = Get-B3aSha256Hex -Bytes (Get-B3aSeedPayload -Seed 'b3a-other' -Size $size)
    Assert-B3aCheck -Name 'outcome-survived' -Ok ((Get-B3aOutcome -Present $true -ActualSize $size -ActualDigest $expectedDigest -ExpectedSize $size -ExpectedDigest $expectedDigest -ZeroDigest $zeroDigest) -eq 'survived')
    Assert-B3aCheck -Name 'outcome-lost-absent' -Ok ((Get-B3aOutcome -Present $false -ExpectedSize $size -ExpectedDigest $expectedDigest -ZeroDigest $zeroDigest) -eq 'lost-absent')
    Assert-B3aCheck -Name 'outcome-lost-zero-filled' -Ok ((Get-B3aOutcome -Present $true -ActualSize $size -ActualDigest $zeroDigest -ExpectedSize $size -ExpectedDigest $expectedDigest -ZeroDigest $zeroDigest) -eq 'lost-zero-filled')
    Assert-B3aCheck -Name 'outcome-lost-torn-size' -Ok ((Get-B3aOutcome -Present $true -ActualSize ($size - 1) -ActualDigest $otherDigest -ExpectedSize $size -ExpectedDigest $expectedDigest -ZeroDigest $zeroDigest) -eq 'lost-torn-size')
    Assert-B3aCheck -Name 'outcome-lost-torn-content' -Ok ((Get-B3aOutcome -Present $true -ActualSize $size -ActualDigest $otherDigest -ExpectedSize $size -ExpectedDigest $expectedDigest -ZeroDigest $zeroDigest) -eq 'lost-torn-content')

    # ---- (e) calibration verdict (fixture-driven) ---------------------------
    $vUnusable = Get-B3aCalibrationVerdict -JournalChainOk $true `
        -NegativeValidRounds 5 -NegativeLosses 0 -PositiveValidRounds 5 -PositiveLosses 0 `
        -InventoryRepeatabilityOk $true -AuditVerdictsOk $true -RequiredN 5
    Assert-B3aCheck -Name 'verdict-unusable-neg-no-loss' -Ok ($vUnusable.gate -eq 'UNUSABLE')

    $vCalibrated = Get-B3aCalibrationVerdict -JournalChainOk $true `
        -NegativeValidRounds 5 -NegativeLosses 2 -PositiveValidRounds 5 -PositiveLosses 0 `
        -InventoryRepeatabilityOk $true -AuditVerdictsOk $true -RequiredN 5
    Assert-B3aCheck -Name 'verdict-calibrated' -Ok ($vCalibrated.gate -eq 'CALIBRATED')

    $vPosLoss = Get-B3aCalibrationVerdict -JournalChainOk $true `
        -NegativeValidRounds 5 -NegativeLosses 2 -PositiveValidRounds 5 -PositiveLosses 1 `
        -InventoryRepeatabilityOk $true -AuditVerdictsOk $true -RequiredN 5
    Assert-B3aCheck -Name 'verdict-unusable-pos-loss' -Ok ($vPosLoss.gate -eq 'UNUSABLE')

    $vChainBad = Get-B3aCalibrationVerdict -JournalChainOk $false `
        -NegativeValidRounds 5 -NegativeLosses 2 -PositiveValidRounds 5 -PositiveLosses 0 `
        -InventoryRepeatabilityOk $true -AuditVerdictsOk $true -RequiredN 5
    Assert-B3aCheck -Name 'verdict-unusable-chain-broken' -Ok ($vChainBad.gate -eq 'UNUSABLE')

    # design S4.2/S4.3: the acceptance rule is stated at N=5, so a verdict over a
    # smaller sample must be impossible rather than merely discouraged.
    $floorThrew = $false
    try {
        [void](Get-B3aCalibrationVerdict -JournalChainOk $true `
            -NegativeValidRounds 1 -NegativeLosses 1 -PositiveValidRounds 1 -PositiveLosses 0 `
            -InventoryRepeatabilityOk $true -AuditVerdictsOk $true -RequiredN 1)
    } catch { $floorThrew = $true }
    Assert-B3aCheck -Name 'verdict-requiredN-floor' -Ok $floorThrew -Detail 'RequiredN=1 must be rejected'

    # ---- (f) loading the round script's modules must not rebind its variables ---
    # Dot-sourcing a script whose top level declares param() rebinds same-named
    # variables in the caller's scope. That silently replaced the round script's
    # command-line -Round/-Candidate with the reconcile module's defaults, so every
    # round ran as round 0/negative-control (found 2026-09-20). The module list is
    # read out of the round script, so a newly dot-sourced file is covered here too.
    $roundText = [System.IO.File]::ReadAllText((Join-Path $PSScriptRoot 'Invoke-B3aRound.ps1'), [System.Text.Encoding]::UTF8)
    $modules = @([regex]::Matches($roundText, "\. \(Join-Path \`$PSScriptRoot '([^']+)'\)") | ForEach-Object { $_.Groups[1].Value })
    # The guard must not pass vacuously: the list is derived from the caller, so an
    # empty list or a missing known module means the extraction itself broke.
    $modulesOk = ($modules.Count -ge 3) `
        -and ($modules -contains 'b3a-rig-lib.ps1') `
        -and ($modules -contains 'b3a-journal-lib.ps1') `
        -and ($modules -contains 'Get-B3aReconcile.ps1')
    $probe = New-Object System.Collections.Generic.List[string]
    $probe.Add('$Round = 7')
    $probe.Add("`$Candidate = 'parity'")
    $probe.Add("`$JournalPath = 'J'")
    $probe.Add("`$stage = 'parity'")
    $probe.Add("`$RoundRootGuest = 'R'")
    foreach ($m in $modules) { $probe.Add(". (Join-Path '" + $PSScriptRoot + "' '" + $m + "')") }
    $probe.Add('if ($Round -ne 7) { "clobbered Round=$Round"; exit 9 }')
    $probe.Add('if ($Candidate -ne "parity") { "clobbered Candidate=$Candidate"; exit 10 }')
    $probe.Add('if ($JournalPath -ne "J") { "clobbered JournalPath=$JournalPath"; exit 11 }')
    $probe.Add('if ($stage -ne "parity") { "clobbered stage=$stage"; exit 12 }')
    $probe.Add('if ($RoundRootGuest -ne "R") { "clobbered RoundRootGuest=$RoundRootGuest"; exit 13 }')
    $probe.Add('exit 0')
    $probePath = Join-Path $FixtureRoot 'probe-module-load.ps1'
    Write-B3aTextFile -Path $probePath -Text ($probe -join "`n") -WithBom
    $psExe = Join-Path $env:WINDIR 'System32\WindowsPowerShell\v1.0\powershell.exe'
    $pr = Invoke-B3aProcess -FileName $psExe -Arguments ('-NoProfile -ExecutionPolicy Bypass -File "{0}"' -f $probePath) -TimeoutS 60
    Assert-B3aCheck -Name 'module-load-preserves-vars' -Ok ($modulesOk -and $pr.exitCode -eq 0) `
        -Detail ('modules={0} extracted={1} exit={2} {3}' -f ($modules -join ','), $modulesOk, $pr.exitCode, $pr.stdoutText.Trim())

    # ---- (g) deploy return shape (PS 5.1 @() on List[object] binder bug) -------
    # [pscustomobject]@{ files = @($results) } with $results a List[object] raises
    # System.ArgumentException inside PowerShell's PSToObjectArrayBinder, so
    # Invoke-B3aToolDeploy could never return - which broke every round (found
    # 2026-09-20 while S0 deployed to a running guest). The call below exercises the
    # return path hermetically: every guest operation fails fast against a bogus
    # VBoxManage path, which is all the shape check needs.
    $prepDir = Join-Path $FixtureRoot 'prep'
    [void](New-Item -ItemType Directory -Path $prepDir -Force)
    Write-B3aTextFile -Path (Join-Path $prepDir '11-shutdown-audit.ps1') -Text '# fixture' -WithBom
    $deployOk = $true
    $deployDetail = ''
    try {
        $dd = Invoke-B3aToolDeploy -VBoxManage (Join-Path $FixtureRoot 'no-such-vboxmanage.exe') `
            -VmName 'fixture' -CredFile (Join-Path $FixtureRoot 'no-cred.txt') -GuestUser 'nobody' `
            -GuestToolDir 'C:\nope' -GuestPrepDir $prepDir -ScriptsDir $PSScriptRoot
        # the shape matters, not just "did not throw": files must be an array of the
        # three expected entries (that is what the PS 5.1 binder bug destroyed)
        $deployOk = ($dd.files -is [array]) -and (@($dd.files).Count -eq 3)
        $deployDetail = ('ok=' + $dd.ok + ' files=' + (@($dd.files).Count) + ' array=' + ($dd.files -is [array]))
    } catch {
        $deployOk = $false
        $deployDetail = 'threw ' + $_.Exception.GetType().FullName + ': ' + $_.Exception.Message
    }
    Assert-B3aCheck -Name 'deploy-return-shape' -Ok $deployOk -Detail $deployDetail
} finally {
    Remove-Item -LiteralPath $FixtureRoot -Recurse -Force -ErrorAction SilentlyContinue
}

''
'checks={0} failures={1}' -f $script:checks, $script:failures
if ($script:failures -eq 0) { exit 0 } else { exit 1 }

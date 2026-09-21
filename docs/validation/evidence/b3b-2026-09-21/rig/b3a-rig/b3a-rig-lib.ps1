# b3a-rig-lib.ps1 - shared host-side helpers for the B3a power-loss injection rig.
#
# Frozen design: docs/t027/B3A_RIG_DESIGN.md (S2 rig architecture, S2.4 reconcile
# classifier, S4.5 calibration verdict). Do not change semantics here without
# revising that document first.
#
# This file holds the small dependency-free helpers shared by the host scripts:
# raw-byte process capture, guestcontrol wrapper, hashing, the deterministic seed
# payload derivation (kept in sync with guest/13-round-write.ps1), the outcome
# classifier, the calibration verdict predicate, and file-writing helpers that
# enforce the repository encoding rules.
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.
#
# NOTE (deviation from the S2.5 script list): the design lists nine scripts and no
# shared library besides b3a-journal-lib.ps1. The shared helpers below are kept in
# this extra file so that b3a-journal-lib.ps1 stays a pure journal module. This is
# an implementation decision to be folded back into the design deliberately.

# ---------------------------------------------------------------------------
# time + hashing

function Get-B3aHostTimestamp {
    # ISO-8601 local time with timezone offset; the authoritative cut time.
    return [DateTimeOffset]::Now.ToString('o')
}

function Get-B3aSha256Hex {
    param(
        [byte[]]$Bytes = $null,
        [string]$Path = ''
    )
    $sha = [System.Security.Cryptography.SHA256]::Create()
    $digest = $null
    if ($Path) {
        $fs = [System.IO.File]::OpenRead($Path)
        try { $digest = $sha.ComputeHash($fs) } finally { $fs.Close() }
    } else {
        $digest = $sha.ComputeHash($Bytes)
    }
    return ([System.BitConverter]::ToString($digest)).Replace('-', '').ToLowerInvariant()
}

function Get-B3aZeroDigest {
    param([Parameter(Mandatory = $true)][long]$Size)
    # SHA-256 of $Size zero bytes; precomputed once per session (design S2.4).
    if ($Size -lt 0) { throw 'Get-B3aZeroDigest: negative size' }
    $zero = New-Object byte[] $Size
    return Get-B3aSha256Hex -Bytes $zero
}

# ---------------------------------------------------------------------------
# deterministic payload derivation (host side)
#
# The guest probe derives its payload from the seed the same way; the host
# recomputes the digest independently to bi-directionally check the probe
# (design S3.2 S5/S6). Derivation:
#   payload[i*32 .. i*32+31] = SHA256( UTF8(seed) || LE32(i) ),  i = 0 .. blocks-1
#   blocks = ceil(size / 32); the result is truncated to size bytes.
# Keep this in sync with the Get-SeedPayload function inside
# guest/13-round-write.ps1; the self-test cross-checks both sides.

function Get-B3aSeedPayload {
    param([Parameter(Mandatory = $true)][string]$Seed, [Parameter(Mandatory = $true)][int]$Size)
    if ($Size -lt 1) { throw 'Get-B3aSeedPayload: size must be positive' }
    $sha = [System.Security.Cryptography.SHA256]::Create()
    $seedBytes = [System.Text.Encoding]::UTF8.GetBytes($Seed)
    $blocks = [int][math]::Ceiling($Size / 32.0)
    $out = New-Object byte[] ($blocks * 32)
    for ($i = 0; $i -lt $blocks; $i++) {
        $counter = [System.BitConverter]::GetBytes([uint32]$i)
        $block = New-Object byte[] ($seedBytes.Length + 4)
        [Array]::Copy($seedBytes, 0, $block, 0, $seedBytes.Length)
        [Array]::Copy($counter, 0, $block, $seedBytes.Length, 4)
        $digest = $sha.ComputeHash($block)
        [Array]::Copy($digest, 0, $out, $i * 32, 32)
    }
    if ($out.Length -eq $Size) { return $out }
    $trimmed = New-Object byte[] $Size
    [Array]::Copy($out, 0, $trimmed, 0, $Size)
    return $trimmed
}

function Get-B3aSeedPayloadDigest {
    param(
        [string]$Seed = '',
        [Parameter(Mandatory = $true)][int]$Size,
        [string]$Literal = $null
    )
    if (-not [string]::IsNullOrEmpty($Literal)) {
        return Get-B3aSha256Hex -Bytes ([System.Text.Encoding]::UTF8.GetBytes($Literal))
    }
    return Get-B3aSha256Hex -Bytes (Get-B3aSeedPayload -Seed $Seed -Size $Size)
}

# ---------------------------------------------------------------------------
# reconcile classifier (design S2.4)

function Get-B3aOutcome {
    param(
        [Parameter(Mandatory = $true)][bool]$Present,
        [long]$ActualSize = 0,
        [string]$ActualDigest = '',
        [long]$ExpectedSize = 0,
        [string]$ExpectedDigest = '',
        [string]$ZeroDigest = ''
    )
    if (-not $Present) { return 'lost-absent' }
    if ($ActualSize -ne $ExpectedSize) { return 'lost-torn-size' }
    if ($ActualDigest -eq $ExpectedDigest) { return 'survived' }
    if ($ActualDigest -eq $ZeroDigest) { return 'lost-zero-filled' }
    return 'lost-torn-content'
}

# ---------------------------------------------------------------------------
# calibration verdict predicate (design S4.5)

function Get-B3aCalibrationVerdict {
    param(
        [Parameter(Mandatory = $true)][bool]$JournalChainOk,
        [int]$NegativeValidRounds = 0,
        [int]$NegativeLosses = 0,
        [int]$PositiveValidRounds = 0,
        [int]$PositiveLosses = 0,
        [bool]$InventoryRepeatabilityOk = $false,
        [bool]$AuditVerdictsOk = $false,
        [int]$RequiredN = 5
    )
    # Defence in depth for the design S4.2/S4.3 acceptance rule (N=5): a verdict
    # computed from fewer rounds could not support "at least one loss in N".
    if ($RequiredN -lt 5) { throw 'Get-B3aCalibrationVerdict: RequiredN must be >= 5 (design S4.2/S4.3)' }
    $calibrated = $JournalChainOk `
        -and ($NegativeValidRounds -ge $RequiredN) `
        -and ($NegativeLosses -ge 1) `
        -and ($PositiveValidRounds -ge $RequiredN) `
        -and ($PositiveLosses -eq 0) `
        -and $InventoryRepeatabilityOk `
        -and $AuditVerdictsOk
    return [ordered]@{
        schema                   = 'b3a-calibration-verdict'
        v                        = 1
        journalChainOk           = $JournalChainOk
        negativeValidRounds      = $NegativeValidRounds
        negativeLosses           = $NegativeLosses
        positiveValidRounds      = $PositiveValidRounds
        positiveLosses           = $PositiveLosses
        inventoryRepeatabilityOk = $InventoryRepeatabilityOk
        auditVerdictsOk          = $AuditVerdictsOk
        requiredN                = $RequiredN
        gate                     = $(if ($calibrated) { 'CALIBRATED' } else { 'UNUSABLE' })
    }
}

# ---------------------------------------------------------------------------
# process capture (raw bytes, no console encoding mangling)

function Invoke-B3aProcess {
    param(
        [Parameter(Mandatory = $true)][string]$FileName,
        [string]$Arguments = '',
        [int]$TimeoutS = 0
    )
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = $FileName
    $psi.Arguments = $Arguments
    $psi.UseShellExecute = $false
    $psi.CreateNoWindow = $true
    $psi.RedirectStandardOutput = $true
    $psi.RedirectStandardError = $true
    $p = New-Object System.Diagnostics.Process
    $p.StartInfo = $psi
    try {
        [void]$p.Start()
    } catch {
        return [pscustomobject]@{
            exitCode    = -1
            stdoutBytes = ([byte[]]@())
            stderrBytes = ([byte[]]@())
            stdoutText  = ''
            stderrText  = ''
            timedOut    = $false
            error       = $_.Exception.Message
        }
    }
    $outMs = New-Object System.IO.MemoryStream
    $errMs = New-Object System.IO.MemoryStream
    $outCopy = $p.StandardOutput.BaseStream.CopyToAsync($outMs)
    $errCopy = $p.StandardError.BaseStream.CopyToAsync($errMs)
    $timedOut = $false
    if ($TimeoutS -gt 0) {
        if (-not $p.WaitForExit($TimeoutS * 1000)) {
            $timedOut = $true
            try { $p.Kill() } catch {
                # The child may have exited between WaitForExit and Kill; the
                # exit code read below is authoritative either way.
                Write-Verbose ('B3a: kill after timeout failed: ' + $_.Exception.Message)
            }
        }
    }
    $p.WaitForExit()
    try { $outCopy.Wait() } catch {
        # A faulted copy task means the pipe closed early; the bytes already
        # buffered in $outMs are still what the caller receives.
        Write-Verbose ('B3a: stdout copy faulted: ' + $_.Exception.Message)
    }
    try { $errCopy.Wait() } catch {
        Write-Verbose ('B3a: stderr copy faulted: ' + $_.Exception.Message)
    }
    $outBytes = $outMs.ToArray()
    $errBytes = $errMs.ToArray()
    return [pscustomobject]@{
        exitCode    = $p.ExitCode
        stdoutBytes = $outBytes
        stderrBytes = $errBytes
        stdoutText  = [System.Text.Encoding]::UTF8.GetString($outBytes)
        stderrText  = [System.Text.Encoding]::UTF8.GetString($errBytes)
        timedOut    = $timedOut
        error       = $null
    }
}

function Get-B3aQuoted {
    # Quote only when the value contains whitespace. A value that already carries a
    # double quote cannot be embedded in a VBoxManage command line safely, so fail
    # closed instead of emitting a line that could break out of the quotes.
    param([Parameter(Mandatory = $true)][string]$Text)
    if ($Text.Contains('"')) { throw 'Get-B3aQuoted: value contains a double quote and cannot be quoted safely' }
    if ($Text -match '\s') { return '"' + $Text + '"' }
    return $Text
}

function Invoke-B3aGuestRun {
    # guestcontrol run with raw stdout/stderr capture (design S2.3 command shape).
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName,
        [Parameter(Mandatory = $true)][string]$CredFile,
        [string]$GuestUser = 'prfrail',
        [Parameter(Mandatory = $true)][string]$Exe,
        [string[]]$GuestArgs = @(),
        [int]$TimeoutMs = 60000
    )
    $parts = @(
        'guestcontrol', (Get-B3aQuoted $VmName), 'run',
        '--exe', (Get-B3aQuoted $Exe),
        '--username', $GuestUser,
        '--passwordfile', (Get-B3aQuoted $CredFile),
        '--wait-stdout', '--wait-stderr',
        '--timeout', ([string]$TimeoutMs),
        '--'
    )
    foreach ($a in $GuestArgs) { $parts += $a }
    return Invoke-B3aProcess -FileName $VBoxManage -Arguments ($parts -join ' ')
}

# ---------------------------------------------------------------------------
# VM facts (read-only queries)

function Get-B3aVmInfo {
    param([Parameter(Mandatory = $true)][string]$VBoxManage,
          [Parameter(Mandatory = $true)][string]$VmName)
    $r = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('showvminfo "{0}" --machinereadable' -f $VmName)
    $map = @{}
    if ($r.exitCode -eq 0) {
        foreach ($lineRaw in ($r.stdoutText -split "`n")) {
            $line = $lineRaw.TrimEnd("`r").Trim()
            $eq = $line.IndexOf('=')
            if ($eq -lt 1) { continue }
            $key = $line.Substring(0, $eq).Trim('"').Trim()
            $val = $line.Substring($eq + 1).Trim()
            if ($val.Length -ge 2 -and $val.StartsWith('"') -and $val.EndsWith('"')) {
                $val = $val.Substring(1, $val.Length - 2)
            }
            # machinereadable escapes backslashes and quotes inside values
            $val = $val.Replace('\\', '\').Replace('\"', '"')
            $map[$key] = $val
        }
    }
    return [pscustomobject]@{ exitCode = $r.exitCode; values = $map; raw = $r.stdoutText }
}

function Get-B3aSnapshotList {
    param([Parameter(Mandatory = $true)][string]$VBoxManage,
          [Parameter(Mandatory = $true)][string]$VmName)
    $r = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('snapshot "{0}" list --machinereadable' -f $VmName)
    $names = New-Object System.Collections.Generic.List[string]
    if ($r.exitCode -eq 0) {
        foreach ($lineRaw in ($r.stdoutText -split "`n")) {
            $line = $lineRaw.TrimEnd("`r").Trim()
            if ($line -match '^SnapshotName="(.*)"$') { $names.Add($Matches[1]) }
        }
    }
    return [pscustomobject]@{ exitCode = $r.exitCode; names = @($names) }
}

# ---------------------------------------------------------------------------
# text helpers

function Get-B3aJsonFromText {
    # First line of $Text that parses as JSON, else $null.
    param([Parameter(Mandatory = $true)][string]$Text)
    foreach ($lineRaw in ($Text -split "`n")) {
        $line = $lineRaw.Trim()
        if (-not $line) { continue }
        try { return ($line | ConvertFrom-Json) } catch {
            # Expected: guest probes also print human-readable lines, so keep
            # scanning the remaining lines instead of failing the whole parse.
            Write-Verbose ('B3a: line is not JSON, skipping: ' + $line)
        }
    }
    return $null
}

function Get-B3aFileTextUtf8 {
    # Strict UTF-8 read for JSON/artifacts (Get-Content defaults to the ANSI
    # code page in PS 5.1 and would corrupt UTF-8 files without BOM).
    param([Parameter(Mandatory = $true)][string]$Path)
    return [System.IO.File]::ReadAllText($Path, [System.Text.Encoding]::UTF8)
}

function Read-B3aAuditVerdict {
    # Parse the 11-shutdown-audit.ps1 output. Returns @{verdict; uptimeMinutes; raw}.
    param([Parameter(Mandatory = $true)][string]$Text)
    $raw = ''
    $m = [regex]::Match($Text, 'VERDICT:\s*([^\r\n]*)')
    if ($m.Success) { $raw = $m.Groups[1].Value.Trim() }
    $verdict = 'OTHER'
    if ($raw -match 'LICENCE-SHUTDOWN') { $verdict = 'LICENCE-SHUTDOWN' }
    elseif ($raw -match 'HARD-POWER-LOSS') { $verdict = 'HARD-POWER-LOSS' }
    elseif ($raw -match 'GRACEFUL-SHUTDOWN') { $verdict = 'GRACEFUL-SHUTDOWN' }
    elseif ($raw -match 'no shutdown') { $verdict = 'NONE' }
    $uptime = $null
    $u = [regex]::Match($Text, 'uptime\s*=\s*([0-9]+(?:\.[0-9]+)?)\s*minutes')
    if ($u.Success) { $uptime = [double]$u.Groups[1].Value }
    return [ordered]@{ verdict = $verdict; uptimeMinutes = $uptime; raw = $raw }
}

function Test-B3aPs1Encoding {
    # UTF-8 with BOM + LF only (CODING_CONVENTIONS.md hard rule for .ps1).
    param([Parameter(Mandatory = $true)][string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) { return $false }
    $bytes = [System.IO.File]::ReadAllBytes($Path)
    if ($bytes.Length -lt 3 -or $bytes[0] -ne 0xEF -or $bytes[1] -ne 0xBB -or $bytes[2] -ne 0xBF) { return $false }
    for ($i = 0; $i -lt $bytes.Length - 1; $i++) {
        if ($bytes[$i] -eq 13 -and $bytes[$i + 1] -eq 10) { return $false }
    }
    return $true
}

function Write-B3aTextFile {
    # UTF-8 (optional BOM) + LF; creates parent directories.
    param([Parameter(Mandatory = $true)][string]$Path,
          [Parameter(Mandatory = $true)][string]$Text,
          [switch]$WithBom)
    $dir = Split-Path -Parent $Path
    if ($dir -and -not (Test-Path -LiteralPath $dir)) { [void](New-Item -ItemType Directory -Path $dir -Force) }
    $enc = New-Object System.Text.UTF8Encoding([bool]$WithBom)
    $preamble = $enc.GetPreamble()
    $body = $enc.GetBytes(($Text -replace "`r`n", "`n"))
    $all = New-Object byte[] ($preamble.Length + $body.Length)
    [Array]::Copy($preamble, 0, $all, 0, $preamble.Length)
    [Array]::Copy($body, 0, $all, $preamble.Length, $body.Length)
    [System.IO.File]::WriteAllBytes($Path, $all)
}

function Write-B3aRawTextFile {
    # Persist captured raw bytes as UTF-8 without BOM + LF.
    param([Parameter(Mandatory = $true)][string]$Path, [byte[]]$Bytes)
    $text = [System.Text.Encoding]::UTF8.GetString($Bytes)
    Write-B3aTextFile -Path $Path -Text $text
}

function Get-B3aBoundedText {
    # Single-line, length-bounded text for log lines: child process output can be
    # long and the log must stay readable and one line per entry.
    param([string]$Text = '', [int]$Max = 2000)
    if ([string]::IsNullOrEmpty($Text)) { return '' }
    $flat = ($Text -replace "`r?`n", ' | ').Trim()
    if ($flat.Length -gt $Max) { return ($flat.Substring(0, $Max) + '...') }
    return $flat
}

function Add-B3aLogLine {
    # Append one LF-terminated UTF-8 (no BOM) line; used for session.log.
    # The round script runs as a child process while the calibration keeps logging
    # to the same file, so the append must tolerate a short sharing conflict:
    # 2026-09-20 a sharing violation killed a round and its parent (both run with
    # $ErrorActionPreference = 'Stop'), which also lost the log line that would have
    # explained it. Open with FileShare.ReadWrite and retry a few times.
    param([Parameter(Mandatory = $true)][string]$Path,
          [Parameter(Mandatory = $true)][string]$Text)
    $dir = Split-Path -Parent $Path
    if ($dir -and -not (Test-Path -LiteralPath $dir)) { [void](New-Item -ItemType Directory -Path $dir -Force) }
    $enc = New-Object System.Text.UTF8Encoding($false)
    $bytes = $enc.GetBytes($Text + "`n")
    $attempt = 0
    while ($true) {
        $attempt++
        try {
            $fs = New-Object System.IO.FileStream($Path, [System.IO.FileMode]::Append, [System.IO.FileAccess]::Write, [System.IO.FileShare]::ReadWrite)
            try {
                $fs.Write($bytes, 0, $bytes.Length)
                $fs.Flush($true)
            } finally {
                $fs.Close()
            }
            return
        } catch {
            if ($attempt -ge 5) { throw }
            Start-Sleep -Milliseconds (50 * $attempt)
        }
    }
}

function Test-B3aBytesEqual {
    param([byte[]]$A, [byte[]]$B)
    if ($null -eq $A -or $null -eq $B) { return ($null -eq $A -and $null -eq $B) }
    if ($A.Length -ne $B.Length) { return $false }
    for ($i = 0; $i -lt $A.Length; $i++) { if ($A[$i] -ne $B[$i]) { return $false } }
    return $true
}

# ---------------------------------------------------------------------------
# guest tool deployment (design S0 item 6; also run per round, see the note in
# Invoke-B3aRound.ps1: restore pristine reverts C:\prfrail-prep, so deployment
# must be repeated after every restore).

function Invoke-B3aToolDeploy {
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingPlainTextForPassword', '', Justification = 'CredFile is a path to a host-side credential file, not a credential value; the file itself is kept outside the repository.')]
    param(
        [Parameter(Mandatory = $true)][string]$VBoxManage,
        [Parameter(Mandatory = $true)][string]$VmName,
        [Parameter(Mandatory = $true)][string]$CredFile,
        [string]$GuestUser = 'prfrail',
        [Parameter(Mandatory = $true)][string]$GuestToolDir,
        [Parameter(Mandatory = $true)][string]$GuestPrepDir,
        [Parameter(Mandatory = $true)][string]$ScriptsDir
    )
    $files = @(
        [ordered]@{ Name = '11-shutdown-audit.ps1';  Source = (Join-Path $GuestPrepDir '11-shutdown-audit.ps1') },
        [ordered]@{ Name = '12-round-inventory.ps1'; Source = (Join-Path $ScriptsDir 'guest\12-round-inventory.ps1') },
        [ordered]@{ Name = '13-round-write.ps1';     Source = (Join-Path $ScriptsDir 'guest\13-round-write.ps1') }
    )
    $results = New-Object System.Collections.Generic.List[object]
    $okAll = $true
    foreach ($f in $files) {
        $entry = [ordered]@{ name = $f.Name; copied = $false; hashMatch = $false; hostSha256 = ''; guestSha256 = ''; note = '' }
        $tmp = $null
        if (-not (Test-Path -LiteralPath $f.Source)) {
            $entry.note = 'source file missing'
            $okAll = $false
            $results.Add($entry)
            continue
        }
        $source = $f.Source
        if (-not (Test-B3aPs1Encoding -Path $f.Source)) {
            # Q7: normalize into a throwaway copy instead of touching the source file.
            $tmp = Join-Path $env:TEMP ('b3a-deploy-' + [guid]::NewGuid().ToString('N') + '.ps1')
            Write-B3aTextFile -Path $tmp -Text ([System.IO.File]::ReadAllText($f.Source)) -WithBom
            $source = $tmp
            $entry.note = 'deployed from normalized temp copy (source encoding violated BOM+LF)'
        }
        $entry.hostSha256 = Get-B3aSha256Hex -Path $source
        # drop a stale copy on the guest first; a missing file is not an error
        [void](Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
            -Exe 'C:\Windows\System32\cmd.exe' `
            -GuestArgs @('/c', 'del', '/q', (Get-B3aQuoted ($GuestToolDir + '\' + $f.Name))) -TimeoutMs 30000)
        $cp = Invoke-B3aProcess -FileName $VBoxManage -Arguments ('guestcontrol "{0}" copyto --username {1} --passwordfile "{2}" "{3}" "{4}\{5}"' -f `
            $VmName, $GuestUser, $CredFile, $source, $GuestToolDir, $f.Name)
        if ($cp.exitCode -ne 0) {
            $entry.note = 'copyto failed (exit ' + $cp.exitCode + ')'
            $okAll = $false
            $results.Add($entry)
            if ($tmp) { [void](Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue) }
            continue
        }
        $entry.copied = $true
        # guest-side hash comparison (design S0 item 6)
        $h = Invoke-B3aGuestRun -VBoxManage $VBoxManage -VmName $VmName -CredFile $CredFile -GuestUser $GuestUser `
            -Exe 'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
            -GuestArgs @('-NoProfile', '-Command', ('"(Get-FileHash -Algorithm SHA256 -Path ''{0}\{1}'').Hash"' -f $GuestToolDir, $f.Name)) `
            -TimeoutMs 60000
        $m = [regex]::Match($h.stdoutText, '[0-9A-Fa-f]{64}')
        $entry.guestSha256 = if ($m.Success) { $m.Value.ToLowerInvariant() } else { '' }
        $entry.hashMatch = ($entry.guestSha256 -eq $entry.hostSha256)
        if (-not $entry.hashMatch) { $entry.note = 'guest hash differs from host source'; $okAll = $false }
        # the Q7 normalization copy never outlives the deployment of its own file
        # (both exit paths are covered here and at the copyto failure above)
        if ($tmp) { [void](Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue) }
        $results.Add($entry)
    }
    # .ToArray() rather than @($results): PS 5.1 raises ArgumentException
    # ("parameter type mismatch") inside PSToObjectArrayBinder for
    # [pscustomobject]@{ files = @(<List[object]>) }, so the function could never
    # return (found 2026-09-20 while S0 deployed to a running guest).
    return [pscustomobject]@{ ok = $okAll; files = $results.ToArray() }
}

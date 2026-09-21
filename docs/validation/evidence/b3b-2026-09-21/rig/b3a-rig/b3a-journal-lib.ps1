# b3a-journal-lib.ps1 - host-side crash-point journal for the B3a rig (design S2.1).
#
# JSONL, UTF-8 without BOM + LF, one record per line, append-only. Write path:
#   FileStream(Append, FileShare.Read, 4096, FileOptions.WriteThrough)
#   -> Write -> Flush($true) -> Close -> read-back verify -> ACK.
# The ACK is observable (returned object + Assert-B3aAck) so callers cannot
# accidentally proceed on a failed append (design S2.1 confirmation protocol:
# "journal first; inject only after the record is acknowledged").
#
# Record schema (fixed key order):
#   schema, v, round, candidate, stage, targetPath, contentDigest, contentSize,
#   planCutDelayS, op, hostTs,
#   [cut: cutExitCode, cutHostTs | reconcile: outcome, inventoryDigests,
#        auditDigest, verdict, voidReason]
#   prevHash, recordHash
#
# Hash convention (the design leaves the exact bytes open; frozen here):
#   recordHash = SHA256( UTF8(prevHashHex) + UTF8(hashedText) ) hex lower
#   hashedText = the record serialized WITHOUT the recordHash field, in the fixed
#   key order above (prevHash already present; first record prevHash=null, the
#   UTF8 prefix is then the literal string "null").
# The stored line is hashedText with "recordHash" appended as the last key, so the
# hashed text is recoverable from the stored line alone. Verifiers additionally
# require the stored line to equal that canonical reconstruction byte for byte:
# bytes after the appended field are outside the hashed text, so without this check
# a forged/duplicated recordHash field or an appended tail would leave the chain
# green (independent review, 2026-09-20).
#
# Verification hashes that STORED LINE TEXT, never a re-serialization of the parsed
# object. ConvertFrom-Json turns ISO-8601 strings into [datetime], whose
# re-serialization drops a trailing zero in the 7th fractional digit
# ('...8416210+08:00' -> '...841621+08:00', measured 2026-09-20), which made the
# old re-serialize-and-compare verifier fail intermittently (~1 in 10 records).
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

. (Join-Path $PSScriptRoot 'b3a-rig-lib.ps1')

function Get-B3aRecordFieldOrder {
    # Fixed canonical key order shared by writer and verifier.
    return @(
        'schema', 'v', 'round', 'candidate', 'stage', 'targetPath', 'contentDigest',
        'contentSize', 'planCutDelayS', 'op', 'hostTs',
        'cutExitCode', 'cutHostTs',
        'outcome', 'inventoryDigests', 'auditDigest', 'verdict', 'voidReason',
        'prevHash', 'recordHash'
    )
}

function Get-B3aCanonicalRecordJson {
    # The record serialized in fixed key order without the recordHash field.
    param([Parameter(Mandatory = $true)][hashtable]$Record)
    $ordered = [ordered]@{}
    foreach ($k in (Get-B3aRecordFieldOrder)) {
        if ($k -eq 'recordHash') { continue }
        if ($Record.Contains($k)) { $ordered[$k] = $Record[$k] }
    }
    return ($ordered | ConvertTo-Json -Compress -Depth 20)
}

function Get-B3aStoredRecordPrefix {
    # The exact text that was hashed for a stored line. The write path always
    # appends ',"recordHash":"<hex>"}' last, so the hashed text is everything
    # before that marker, re-closed with '}'. Returns $null when the marker is
    # absent (i.e. the line was not produced by Add-B3aJournalRecord).
    param([Parameter(Mandatory = $true)][string]$Line)
    $marker = ',"recordHash":"'
    $idx = $Line.LastIndexOf($marker)
    if ($idx -lt 0) { return $null }
    return ($Line.Substring(0, $idx) + '}')
}

function Get-B3aStoredRecordLine {
    # The canonical stored line. HashedText is the canonical JSON object (closing
    # brace included); the stored line replaces that brace with the appended
    # recordHash field:
    #   hashedText = '<... fields ...>'
    #   line       = '<... fields ...>,"recordHash":"<64 hex>"}'
    # Both the writer and the verifiers go through this, so the line shape has one
    # definition and verification can compare the stored line byte for byte.
    param([Parameter(Mandatory = $true)][string]$HashedText, [Parameter(Mandatory = $true)][string]$RecordHash)
    if (-not $HashedText.EndsWith('}')) { throw 'Get-B3aStoredRecordLine: hashed text must be a JSON object' }
    return ($HashedText.Substring(0, $HashedText.Length - 1) + ',"recordHash":"' + $RecordHash + '"}')
}

function Get-B3aLineHash {
    # recordHash = SHA256( UTF8(prevHashHex | 'null') + UTF8(hashedText) ), hex lower.
    # The single hashing primitive: the write path hashes the canonical text it is
    # about to store, the verifier hashes the text that is stored.
    param([Parameter(Mandatory = $true)][string]$HashedText, [string]$PrevHashHex = '')
    $prefixText = 'null'
    if (-not [string]::IsNullOrEmpty($PrevHashHex)) { $prefixText = $PrevHashHex }
    $prefix = [System.Text.Encoding]::UTF8.GetBytes($prefixText)
    $body = [System.Text.Encoding]::UTF8.GetBytes($HashedText)
    $all = New-Object byte[] ($prefix.Length + $body.Length)
    [Array]::Copy($prefix, 0, $all, 0, $prefix.Length)
    [Array]::Copy($body, 0, $all, $prefix.Length, $body.Length)
    $sha = [System.Security.Cryptography.SHA256]::Create()
    $digest = $sha.ComputeHash($all)
    return ([System.BitConverter]::ToString($digest)).Replace('-', '').ToLowerInvariant()
}

function Get-B3aRecordHash {
    # Write path only: hash the canonical serialization that is about to be stored.
    param([Parameter(Mandatory = $true)][hashtable]$Record)
    $prev = ''
    if ($Record.Contains('prevHash') -and $null -ne $Record['prevHash']) { $prev = [string]$Record['prevHash'] }
    return Get-B3aLineHash -HashedText (Get-B3aCanonicalRecordJson -Record $Record) -PrevHashHex $prev
}

function New-B3aJournalRecord {
    # Base record for op plan/cut/reconcile; prevHash/recordHash are added by
    # Add-B3aJournalRecord. Extra carries the op-specific fields (design S2.1).
    # ShouldProcess does not apply here: this function only builds an in-memory
    # ordered hashtable. Persistence (and the confirm/WhatIf surface) lives in
    # Add-B3aJournalRecord.
    [Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '', Justification = 'Pure in-memory record construction (returns an ordered hashtable); no system state is created or changed. Add-B3aJournalRecord performs the actual persistence.')]
    param(
        [Parameter(Mandatory = $true)][int]$Round,
        [Parameter(Mandatory = $true)][string]$Candidate,
        [Parameter(Mandatory = $true)][string]$Stage,
        [Parameter(Mandatory = $true)][string]$TargetPath,
        [Parameter(Mandatory = $true)][string]$ContentDigest,
        [Parameter(Mandatory = $true)][long]$ContentSize,
        [Parameter(Mandatory = $true)][double]$PlanCutDelayS,
        [Parameter(Mandatory = $true)][ValidateSet('plan', 'cut', 'reconcile')][string]$Op,
        [hashtable]$Extra = @{},
        # Production callers omit HostTs (the host clock is read here). The
        # parameter exists so the timestamp round-trip regression test can pin a
        # value whose 7th fractional digit is a trailing zero.
        [string]$HostTs = $null
    )
    $ts = $HostTs
    if ([string]::IsNullOrEmpty($ts)) { $ts = Get-B3aHostTimestamp }
    $rec = [ordered]@{
        schema        = 'b3a-journal'
        v             = 1
        round         = $Round
        candidate     = $Candidate
        stage         = $Stage
        targetPath    = $TargetPath
        contentDigest = $ContentDigest
        contentSize   = $ContentSize
        planCutDelayS = $PlanCutDelayS
        op            = $Op
        hostTs        = $ts
    }
    foreach ($k in $Extra.Keys) { $rec[$k] = $Extra[$k] }
    return $rec
}

function Get-B3aRawJournalLine {
    param([Parameter(Mandatory = $true)][string]$JournalPath)
    if (-not (Test-Path -LiteralPath $JournalPath)) { return @() }
    $bytes = [System.IO.File]::ReadAllBytes($JournalPath)
    $text = [System.Text.Encoding]::UTF8.GetString($bytes)
    $list = New-Object System.Collections.Generic.List[string]
    foreach ($l in ($text -split "`n")) { $list.Add($l) }
    while ($list.Count -gt 0 -and $list[$list.Count - 1].Trim() -eq '') {
        $list.RemoveAt($list.Count - 1)
    }
    # callers always wrap with @(...); the flat array unrolls on return
    return @($list)
}

function Add-B3aJournalRecord {
    # Append + flush + close + read-back verify; returns
    # @{ack; recordHash; recordJson; error}. Never throws for expected failures.
    param(
        [Parameter(Mandatory = $true)][string]$JournalPath,
        [Parameter(Mandatory = $true)][hashtable]$Record
    )
    $result = [pscustomobject]@{ ack = $false; recordHash = $null; recordJson = ''; error = $null }
    $dir = Split-Path -Parent $JournalPath
    if ($dir -and -not (Test-Path -LiteralPath $dir)) {
        try { [void](New-Item -ItemType Directory -Path $dir -Force) } catch {
            $result.error = 'journal-directory-create-failed: ' + $_.Exception.Message
            return $result
        }
    }
    try {
        # extend the chain from the current tail, after re-verifying that tail
        $prevHash = $null
        if (Test-Path -LiteralPath $JournalPath) {
            $lines = @(Get-B3aRawJournalLine -JournalPath $JournalPath)
            if ($lines.Count -gt 0) {
                $lastJson = $lines[$lines.Count - 1]
                try { $last = $lastJson | ConvertFrom-Json } catch {
                    $result.error = 'tail-line-unparseable'
                    return $result
                }
                $tailPrefix = Get-B3aStoredRecordPrefix -Line $lastJson
                if ($null -eq $tailPrefix) {
                    $result.error = 'tail-recordHash-field-missing'
                    return $result
                }
                $tailPrev = ''
                if ($null -ne $last.prevHash) { $tailPrev = [string]$last.prevHash }
                $tailHash = Get-B3aLineHash -HashedText $tailPrefix -PrevHashHex $tailPrev
                if ($tailHash -ne ([string]$last.recordHash)) {
                    $result.error = 'tail-hash-mismatch'
                    return $result
                }
                if ((Get-B3aStoredRecordLine -HashedText $tailPrefix -RecordHash $tailHash) -ne $lastJson) {
                    $result.error = 'tail-line-not-canonical'
                    return $result
                }
                $prevHash = [string]$last.recordHash
            }
        }
        $rec = @{}
        foreach ($k in $Record.Keys) { $rec[$k] = $Record[$k] }
        $rec['prevHash'] = $prevHash
        $rec['recordHash'] = Get-B3aRecordHash -Record $rec
        $canonical = Get-B3aCanonicalRecordJson -Record $rec
        $json = Get-B3aStoredRecordLine -HashedText $canonical -RecordHash ([string]$rec['recordHash'])
        # design write path: WriteThrough + Flush($true) + close
        $bytes = [System.Text.Encoding]::UTF8.GetBytes($json + "`n")
        $fs = New-Object System.IO.FileStream($JournalPath, [System.IO.FileMode]::Append, [System.IO.FileAccess]::Write, [System.IO.FileShare]::Read, 4096, [System.IO.FileOptions]::WriteThrough)
        try {
            $fs.Write($bytes, 0, $bytes.Length)
            $fs.Flush($true)
        } finally {
            $fs.Close()
        }
        # read-back verification: reopen read-only, read the last stored line,
        # parse it and compare the stored recordHash (design S2.1 self-proof)
        $lines = @(Get-B3aRawJournalLine -JournalPath $JournalPath)
        if ($lines.Count -lt 1) { $result.error = 'read-back-empty'; return $result }
        $parsed = $null
        try { $parsed = $lines[$lines.Count - 1] | ConvertFrom-Json } catch {
            $result.error = 'read-back-unparseable'
            return $result
        }
        if (([string]$parsed.recordHash) -ne $rec['recordHash']) {
            $result.error = 'read-back-hash-mismatch'
            return $result
        }
        $result.ack = $true
        $result.recordHash = $rec['recordHash']
        $result.recordJson = $json
        return $result
    } catch {
        $result.error = 'append-failed: ' + $_.Exception.Message
        return $result
    }
}

function Test-B3aJournalChain {
    # Full-chain verify from the first record (design S0 item 5 / S6 item 10).
    # Returns @{ok; count; firstBadLine; firstError; records}.
    param([Parameter(Mandatory = $true)][string]$JournalPath)
    $result = [pscustomobject]@{ ok = $true; count = 0; firstBadLine = 0; firstError = $null; records = @() }
    if (-not (Test-Path -LiteralPath $JournalPath)) { return $result }
    $lines = @(Get-B3aRawJournalLine -JournalPath $JournalPath)
    $prevHash = $null
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $lineNo = $i + 1
        $parsed = $null
        try { $parsed = $lines[$i] | ConvertFrom-Json } catch {
            $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'unparseable-line'
            return $result
        }
        if ($null -eq $parsed.recordHash) {
            $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'missing-recordHash'
            return $result
        }
        $prefix = Get-B3aStoredRecordPrefix -Line $lines[$i]
        if ($null -eq $prefix) {
            $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'recordHash-field-missing'
            return $result
        }
        $ownPrev = ''
        if ($null -ne $parsed.prevHash) { $ownPrev = [string]$parsed.prevHash }
        $expect = Get-B3aLineHash -HashedText $prefix -PrevHashHex $ownPrev
        if ($expect -ne ([string]$parsed.recordHash)) {
            $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'recordHash-mismatch'
            return $result
        }
        if ((Get-B3aStoredRecordLine -HashedText $prefix -RecordHash $expect) -ne $lines[$i]) {
            $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'stored-line-not-canonical'
            return $result
        }
        if ($i -eq 0) {
            if ($null -ne $parsed.prevHash) {
                $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'first-prevHash-not-null'
                return $result
            }
        } else {
            if (([string]$parsed.prevHash) -ne $prevHash) {
                $result.ok = $false; $result.firstBadLine = $lineNo; $result.firstError = 'prevHash-chain-break'
                return $result
            }
        }
        $prevHash = [string]$parsed.recordHash
        $result.records += $parsed
    }
    $result.count = $lines.Count
    return $result
}

function Get-B3aJournalRecord {
    # Parsed records in order (used by the calibration statistics).
    param([Parameter(Mandatory = $true)][string]$JournalPath)
    $out = @()
    foreach ($line in @(Get-B3aRawJournalLine -JournalPath $JournalPath)) {
        try { $out += ($line | ConvertFrom-Json) } catch { throw ('journal record unparseable: ' + $_) }
    }
    return $out
}

function Assert-B3aAck {
    # Callers must gate on this: no ACK means never proceed (design S2.1).
    param([Parameter(Mandatory = $true)]$AckResult, [string]$Context = 'journal record')
    if ($null -eq $AckResult -or -not $AckResult.ack) {
        $detail = 'missing ack'
        if ($null -ne $AckResult -and $AckResult.error) { $detail = $AckResult.error }
        throw ('B3a journal ACK failed (' + $Context + '): ' + $detail)
    }
    return $true
}

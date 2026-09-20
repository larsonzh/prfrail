# 13-round-write.ps1 - guest-side control write for the B3a rig (design S4.2 / S4.3).
#
# negative: exactly one buffered [System.IO.File]::WriteAllBytes and NO flush of
#           any kind (no Flush(true), no FileOptions.WriteThrough, no directory/
#           volume-handle flush). This is the control that MUST be lost.
# positive: FileStream(Create, Write, FileShare.None, 4096, FileOptions.WriteThrough)
#           -> Write -> Flush($true) -> Close, then a volume-handle
#           FlushFileBuffers on \\.\<drive>: via P/Invoke. Every step is a hard
#           assertion; any failure exits non-zero. This control MUST survive.
#           The volume flush is a rig-only calibration means (design S1.8); it
#           proves nothing about product candidates C2/C3.
#
# Payload (deterministic, recomputable by the host; keep in sync with
# Get-B3aSeedPayload in b3a-rig-lib.ps1):
#   payload[i*32 .. i*32+31] = SHA256( UTF8(seed) || LE32(i) ), i = 0 .. blocks-1
#   blocks = ceil(size / 32); truncated to size bytes.
# -Literal <s> overrides the derivation: payload = UTF8 bytes of <s> and must be
# exactly -Size bytes (used by the R0 parity round to replicate the dryrun's
# 14-byte 'pre-cut-marker').
#
# stdout = one line of JSON (UTF-8, no BOM) with the digest of the bytes actually
# written and the list of APIs actually called. The guest clock in the probe is
# labelled secondary ('guest-clock') and is never evidence.
#
# Exit 0 = written as requested; exit 1 = write/commit error (probe still on stdout).
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

param(
    [Parameter(Mandatory = $true)][ValidateSet('negative', 'positive')][string]$Mode,
    [Parameter(Mandatory = $true)][string]$Path,
    [int]$Size = 262144,
    [string]$Seed = 'b3a-neg-0',
    [string]$Literal = $null
)

$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = New-Object System.Text.UTF8Encoding($false)

function Get-SeedPayload {
    param([Parameter(Mandatory = $true)][string]$Seed, [Parameter(Mandatory = $true)][int]$Size)
    if ($Size -lt 1) { throw 'size must be positive' }
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

function Get-Sha256Hex {
    param([Parameter(Mandatory = $true)][byte[]]$Bytes)
    $sha = [System.Security.Cryptography.SHA256]::Create()
    return ([System.BitConverter]::ToString($sha.ComputeHash($Bytes))).Replace('-', '').ToLowerInvariant()
}

$apis = New-Object System.Collections.Generic.List[string]
$derivation = 'sha256-blocks'
$contentDigest = $null
$commit = 'none'
$guestClock = (Get-Date -Format 'yyyy-MM-dd HH:mm:ss')

try {
    if (-not [string]::IsNullOrEmpty($Literal)) {
        $payload = [System.Text.Encoding]::UTF8.GetBytes($Literal)
        if ($payload.Length -ne $Size) {
            throw ('literal payload is {0} bytes but -Size is {1}' -f $payload.Length, $Size)
        }
        $derivation = 'literal'
    } else {
        $payload = Get-SeedPayload -Seed $Seed -Size $Size
    }
    $contentDigest = Get-Sha256Hex -Bytes $payload
    $writeResult = 'ok'

    if ($Mode -eq 'negative') {
        # frozen negative control: single buffered write, deliberately no flush
        [System.IO.File]::WriteAllBytes($Path, $payload)
        $apis.Add('File.WriteAllBytes')
    } else {
        # frozen positive control commit chain (design S4.3)
        $fs = New-Object System.IO.FileStream($Path, [System.IO.FileMode]::Create, [System.IO.FileAccess]::Write, [System.IO.FileShare]::None, 4096, [System.IO.FileOptions]::WriteThrough)
        try {
            $fs.Write($payload, 0, $payload.Length)
            $apis.Add('FileStream.Write')
            $fs.Flush($true)
            $apis.Add('FileStream.Flush(true)')
        } finally {
            $fs.Close()
            $apis.Add('FileStream.Close')
        }
        $drive = ([System.IO.Path]::GetPathRoot($Path)).TrimEnd('\')
        $volume = '\\.\' + $drive
        if (-not ('B3aNative' -as [type])) {
            Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class B3aNative {
    [DllImport("kernel32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern IntPtr CreateFile(string lpFileName, uint dwDesiredAccess, uint dwShareMode, IntPtr lpSecurityAttributes, uint dwCreationDisposition, uint dwFlagsAndAttributes, IntPtr hTemplateFile);
    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern bool FlushFileBuffers(IntPtr hFile);
    [DllImport("kernel32.dll", SetLastError = true)]
    public static extern bool CloseHandle(IntPtr hObject);
}
'@
        }
        # GENERIC_WRITE 0x40000000, FILE_SHARE_READ|FILE_SHARE_WRITE = 3, OPEN_EXISTING = 3
        $h = [B3aNative]::CreateFile($volume, 0x40000000, 3, [IntPtr]::Zero, 3, 0, [IntPtr]::Zero)
        $apis.Add('CreateFile(' + $volume + ')')
        if ($h -eq [IntPtr]::new(-1)) {
            throw ('CreateFile failed on {0} (error {1})' -f $volume, [System.Runtime.InteropServices.Marshal]::GetLastWin32Error())
        }
        try {
            if (-not [B3aNative]::FlushFileBuffers($h)) {
                throw ('FlushFileBuffers failed on {0} (error {1})' -f $volume, [System.Runtime.InteropServices.Marshal]::GetLastWin32Error())
            }
            $apis.Add('FlushFileBuffers')
        } finally {
            [void][B3aNative]::CloseHandle($h)
            $apis.Add('CloseHandle')
        }
        $commit = 'positive-commit-ok'
    }

    $probe = [ordered]@{
        schema            = 'b3a-write-probe'
        v                 = 1
        mode              = $Mode
        path              = $Path
        size              = $payload.Length
        seed              = $Seed
        derivation        = $derivation
        contentDigest     = $contentDigest
        writeResult       = $writeResult
        commit            = $commit
        apisCalled        = @($apis)
        guestClock        = $guestClock
        secondaryEvidence = 'guest-clock'
    }
    [Console]::Out.WriteLine(($probe | ConvertTo-Json -Compress -Depth 10))
    exit 0
} catch {
    $probe = [ordered]@{
        schema            = 'b3a-write-probe'
        v                 = 1
        mode              = $Mode
        path              = $Path
        size              = $Size
        seed              = $Seed
        derivation        = $derivation
        contentDigest     = $contentDigest
        writeResult       = 'error'
        commit            = $commit
        error             = $_.Exception.Message
        apisCalled        = @($apis)
        guestClock        = $guestClock
        secondaryEvidence = 'guest-clock'
    }
    [Console]::Out.WriteLine(($probe | ConvertTo-Json -Compress -Depth 10))
    exit 1
}

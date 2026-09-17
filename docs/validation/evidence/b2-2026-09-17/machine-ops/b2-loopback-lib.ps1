# B2 scratch library: AppContainer SID, loopback exemption list, baseline snapshots and diffs.
# Shared by apply / revert / verify / rehearsal scripts. Delete at slice closeout (after promotion
# into the B2 evidence bundle).
$ErrorActionPreference = 'Stop'

if (-not ('B2LoopbackNative' -as [type])) {
    $csharp = @'
using System;
using System.Runtime.InteropServices;

public static class B2LoopbackNative
{
    [StructLayout(LayoutKind.Sequential)]
    public struct SID_AND_ATTRIBUTES { public IntPtr Sid; public uint Attributes; }

    [DllImport("userenv.dll", CharSet = CharSet.Unicode, SetLastError = true)]
    public static extern int CreateAppContainerProfile(string name, string displayName, string description, IntPtr capabilities, uint capabilityCount, out IntPtr sid);

    [DllImport("userenv.dll", CharSet = CharSet.Unicode, SetLastError = true)]
    public static extern int DeriveAppContainerSidFromAppContainerName(string name, out IntPtr sid);

    [DllImport("userenv.dll", CharSet = CharSet.Unicode, SetLastError = true)]
    public static extern int DeleteAppContainerProfile(string name);

    [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern bool ConvertSidToStringSid(IntPtr sid, out IntPtr stringSid);

    [DllImport("advapi32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
    public static extern bool ConvertStringSidToSid(string sddlSid, out IntPtr sid);

    [DllImport("kernel32.dll")]
    public static extern IntPtr LocalFree(IntPtr handle);

    [DllImport("FirewallAPI.dll")]
    public static extern int NetworkIsolationSetAppContainerConfig(uint dwNumPublicAppCs, IntPtr pAppContainerSIDs);

    public static string SidToString(IntPtr sid)
    {
        IntPtr s;
        if (!ConvertSidToStringSid(sid, out s)) { return ""; }
        string r = Marshal.PtrToStringUni(s);
        LocalFree(s);
        return r == null ? "" : r;
    }

    public static int SetExemptionList(string[] sids)
    {
        int elem = Marshal.SizeOf(typeof(SID_AND_ATTRIBUTES));
        IntPtr array = IntPtr.Zero;
        if (sids != null && sids.Length > 0)
        {
            array = Marshal.AllocHGlobal(elem * sids.Length);
            for (int i = 0; i < sids.Length; i++)
            {
                IntPtr sid;
                if (!ConvertStringSidToSid(sids[i], out sid)) { return -1000 - i; }
                SID_AND_ATTRIBUTES entry = new SID_AND_ATTRIBUTES();
                entry.Sid = sid;
                entry.Attributes = 0;
                Marshal.StructureToPtr(entry, (IntPtr)(array.ToInt64() + i * elem), false);
            }
        }
        int hr = NetworkIsolationSetAppContainerConfig((uint)(sids == null ? 0 : sids.Length), array);
        if (array != IntPtr.Zero) { Marshal.FreeHGlobal(array); }
        return hr;
    }
}
'@
    Add-Type -TypeDefinition $csharp -Language CSharp
}

function Assert-B2Admin([string]$Action) {
    $isAdmin = ([Security.Principal.WindowsPrincipal]::new([Security.Principal.WindowsIdentity]::GetCurrent())).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if (-not $isAdmin) { throw ("$Action requires an elevated PowerShell session (run it yourself as administrator).") }
}

# Encoding discipline: every artifact this tooling writes is UTF-8 *without* BOM, and PowerShell 5.1's
# Get-Content decodes BOM-less UTF-8 as ANSI - which corrupts any non-ASCII path or display name (e.g. a
# Chinese user profile or firewall rule name) and makes ConvertFrom-Json fail on valid JSON. Always read
# through these helpers instead of Get-Content.
function Read-B2TextFile([string]$Path) {
    if (-not (Test-Path -LiteralPath $Path)) { throw ("file not found: {0}" -f $Path) }
    return [IO.File]::ReadAllText($Path)
}

function Read-B2JsonFile([string]$Path) {
    return (Read-B2TextFile -Path $Path | ConvertFrom-Json)
}

function Read-B2Lines([string]$Path) {
    return @(Read-B2TextFile -Path $Path) -split "`r?`n"
}

function New-B2ContainerSidIfMissing([string]$ContainerName) {
    $sidPtr = [IntPtr]::Zero
    $hr = [B2LoopbackNative]::CreateAppContainerProfile($ContainerName, $ContainerName, 'ProofRail B2 enforcement container', [IntPtr]::Zero, 0, [ref]$sidPtr)
    $alreadyExists = -2147024713
    if ($hr -eq $alreadyExists) {
        $sidPtr = [IntPtr]::Zero
        $null = [B2LoopbackNative]::DeriveAppContainerSidFromAppContainerName($ContainerName, [ref]$sidPtr)
    } elseif ($hr -ne 0) {
        throw ("CreateAppContainerProfile failed 0x{0:X8}" -f $hr)
    }
    return [B2LoopbackNative]::SidToString($sidPtr)
}

# Name-only derivation: the SID is a pure function of the profile name, so this works even when the profile
# has been deleted. It never creates anything and never claims the profile exists - use
# Get-B2ContainerProfileExists for that question.
function Get-B2ContainerSidByName([string]$ContainerName) {
    $sidPtr = [IntPtr]::Zero
    $hr = [B2LoopbackNative]::DeriveAppContainerSidFromAppContainerName($ContainerName, [ref]$sidPtr)
    if ($hr -ne 0 -or $sidPtr -eq [IntPtr]::Zero) {
        throw ("could not derive the AppContainer SID for '{0}' (0x{1:X8})" -f $ContainerName, $hr)
    }
    return [B2LoopbackNative]::SidToString($sidPtr)
}

# Derive-only resolver: used by every read-only or reverting step so that a "verify" never creates an
# AppContainer profile as a side effect.
function Get-B2ContainerSidExisting([string]$ContainerName) {
    if (-not (Get-B2ContainerProfileExists -ContainerName $ContainerName)) {
        throw ("AppContainer profile '{0}' does not exist; read-only steps never create it (run apply-b2-loopback.ps1 first)" -f $ContainerName)
    }
    return Get-B2ContainerSidByName -ContainerName $ContainerName
}

# Profile existence must NOT be inferred from DeriveAppContainerSidFromAppContainerName: that API derives a
# SID from the name alone and succeeds even for a name that has no profile (measured 2026-09-17). The
# profile's storage key under the current user's hive is the authoritative per-user marker.
function Get-B2ContainerProfileExists([string]$ContainerName) {
    $storageRoot = 'HKCU:\Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\AppContainer\Storage'
    if (-not (Test-Path -LiteralPath $storageRoot)) { return $false }
    return (Test-Path -LiteralPath (Join-Path $storageRoot $ContainerName))
}

# Parses the loopback exemption list. CheckNetIsolation prints "AppContainer SID:" lines (localised
# builds keep the SID value itself in ASCII), so the SID regex is the reliable part.
function Get-B2ExemptionList {
    $raw = (& CheckNetIsolation.exe LoopbackExempt -s 2>&1 | Out-String)
    $sids = @()
    foreach ($match in [regex]::Matches($raw, 'S-1-15-2(?:-\d+)+')) {
        if ($sids -notcontains $match.Value) { $sids += $match.Value }
    }
    return [pscustomobject]@{ raw = $raw.Trim(); sids = $sids }
}

function Get-B2FirewallSnapshot {
    $rules = @()
    try {
        $rules = Get-NetFirewallRule -ErrorAction Stop | ForEach-Object {
            [pscustomobject]@{
                name      = $_.Name
                display   = $_.DisplayName
                direction = [string]$_.Direction
                action    = [string]$_.Action
                enabled   = [string]$_.Enabled
                profile   = [string]$_.Profile
            }
        } | Sort-Object name
    } catch {
        $rules = @([pscustomobject]@{ name = '<query-failed>'; display = $_.Exception.Message; direction = ''; action = ''; enabled = ''; profile = '' })
    }
    return $rules
}

function Get-B2LoopbackListeners {
    $rows = @()
    $tcp = @()
    try { $tcp = Get-NetTCPConnection -State Listen -ErrorAction Stop | Where-Object { $_.LocalAddress -in @('127.0.0.1', '::1') } } catch { }
    foreach ($c in $tcp) {
        $proc = Get-Process -Id $c.OwningProcess -ErrorAction SilentlyContinue
        $name = ''
        $path = ''
        if ($proc) { $name = $proc.ProcessName; $path = [string]$proc.Path }
        $rows += [pscustomobject]@{ proto = 'tcp'; address = $c.LocalAddress; port = $c.LocalPort; pid = $c.OwningProcess; process = $name; path = $path }
    }
    $udp = @()
    try { $udp = Get-NetUDPEndpoint -ErrorAction SilentlyContinue | Where-Object { $_.LocalAddress -in @('127.0.0.1', '::1') } } catch { }
    foreach ($c in $udp) {
        $proc = Get-Process -Id $c.OwningProcess -ErrorAction SilentlyContinue
        $name = ''
        if ($proc) { $name = $proc.ProcessName }
        $rows += [pscustomobject]@{ proto = 'udp'; address = $c.LocalAddress; port = $c.LocalPort; pid = $c.OwningProcess; process = $name; path = '' }
    }
    return ($rows | Sort-Object proto, port)
}

function New-B2Snapshot([string]$OutFile, [string]$ContainerName) {
    $exempt = Get-B2ExemptionList
    $profileExists = Get-B2ContainerProfileExists -ContainerName $ContainerName
    $sid = ''
    if ($profileExists) { $sid = Get-B2ContainerSidByName -ContainerName $ContainerName }
    $snapshot = [pscustomobject]@{
        schemaVersion = 2
        takenAtUtc    = (Get-Date).ToUniversalTime().ToString('o')
        containerName = $ContainerName
        containerSid  = $sid
        profileExists = $profileExists
        exemptions    = $exempt.sids
        firewallRules = (Get-B2FirewallSnapshot)
        listeners     = (Get-B2LoopbackListeners)
    }
    $json = $snapshot | ConvertTo-Json -Depth 6
    [IO.File]::WriteAllText($OutFile, $json, (New-Object System.Text.UTF8Encoding($false)))
    return $snapshot
}

# A baseline is only usable when it carries the current schema *and* the profileExists field, otherwise the
# strict diff reports an informational line instead of silently assuming absence.
function Test-B2SnapshotSchemaCurrent([object]$Snapshot) {
    if ($null -eq $Snapshot) { return $false }
    $versionProperty = $Snapshot.PSObject.Properties['schemaVersion']
    # Not $profile: that name would shadow PowerShell's automatic $PROFILE variable.
    $profileProperty = $Snapshot.PSObject.Properties['profileExists']
    if ($null -eq $versionProperty -or $null -eq $profileProperty) { return $false }
    return ([int]$versionProperty.Value -ge 2)
}

# Single place that decides which strict-diff lines actually fail a restoral verdict. Informational lines
# (a legacy baseline that cannot answer the profile question) never fail; a profile kept on request is reported
# but is not a restoral failure either. revert, verify and the self-test all call this, so the rule cannot drift.
function Get-B2RestoralFailures([object[]]$DiffLines, [bool]$ProfileKept = $false) {
    return @($DiffLines | Where-Object {
            $_ -notlike 'profile-existence-unverifiable*' -and
            -not ($ProfileKept -and $_ -like 'profile-existence-changed:*')
        })
}

function Get-B2SnapshotDiff([object]$Baseline, [object]$Current) {
    $lines = @()
    # Identity gate first: comparing a baseline recorded for a different container could otherwise report a
    # clean restoral for the wrong object. A mismatch is a failure line, never informational.
    $baselineName = [string]$Baseline.containerName
    $currentName = [string]$Current.containerName
    if ($baselineName -ne $currentName) {
        $lines += ('baseline-container-mismatch: baseline records "' + $baselineName + '" but the current snapshot is "' + $currentName + '"; refusing to compare')
        return $lines
    }
    $baseSids = @($Baseline.exemptions)
    $nowSids = @($Current.exemptions)
    foreach ($sid in $nowSids) { if ($baseSids -notcontains $sid) { $lines += ('exemption-added: ' + $sid) } }
    foreach ($sid in $baseSids) { if ($nowSids -notcontains $sid) { $lines += ('exemption-removed: ' + $sid) } }

    $baseRules = @($Baseline.firewallRules | ForEach-Object { $_.name + '|' + $_.direction + '|' + $_.action + '|' + $_.enabled + '|' + $_.profile })
    $nowRules = @($Current.firewallRules | ForEach-Object { $_.name + '|' + $_.direction + '|' + $_.action + '|' + $_.enabled + '|' + $_.profile })
    foreach ($rule in $nowRules) { if ($baseRules -notcontains $rule) { $lines += ('firewall-added: ' + $rule) } }
    foreach ($rule in $baseRules) { if ($nowRules -notcontains $rule) { $lines += ('firewall-removed: ' + $rule) } }

    # Profile existence is part of the strict diff: creating the AppContainer profile is a machine change
    # even when the exemption table is restored, so it must never hide behind a "zero-difference" result.
    # A legacy baseline that predates the field is reported as unverifiable rather than as a change.
    if (Test-B2SnapshotSchemaCurrent $Baseline) {
        $baseProfile = [bool]$Baseline.profileExists
        $nowProfile = [bool]$Current.profileExists
        if ($baseProfile -ne $nowProfile) {
            $lines += ('profile-existence-changed: ' + $baseProfile + ' -> ' + $nowProfile + ' (' + $Current.containerName + ')')
        }
    } else {
        $lines += ('profile-existence-unverifiable: baseline predates schemaVersion 2; refresh it with apply-b2-loopback.ps1 -RefreshBaseline after a full revert')
    }

    return $lines
}

# Listener churn is environmental (apps start and stop): reported, never used to fail restoral.
function Get-B2ListenerDiff([object]$Baseline, [object]$Current) {
    $lines = @()
    $baseListeners = @($Baseline.listeners | ForEach-Object { $_.proto + ':' + $_.address + ':' + $_.port + ':' + $_.process })
    $nowListeners = @($Current.listeners | ForEach-Object { $_.proto + ':' + $_.address + ':' + $_.port + ':' + $_.process })
    foreach ($item in $nowListeners) { if ($baseListeners -notcontains $item) { $lines += ('listener-added: ' + $item) } }
    foreach ($item in $baseListeners) { if ($nowListeners -notcontains $item) { $lines += ('listener-removed: ' + $item) } }
    return $lines
}

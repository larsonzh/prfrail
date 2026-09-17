#requires -Version 5.1
<#
.SYNOPSIS
    Hermetic self-test for tools/agent-probe/enforcement-proxy.

.DESCRIPTION
    Builds the proxy, then exercises the real binary end to end against local targets only:
    allowlisted and refused forward-proxy requests, observe-mode recording, the fail-fast rule for
    an occupied fixed port, and the JSONL audit log. No external network access, no model calls.

.PARAMETER KeepArtifacts
    Keep the temporary self-test directory instead of removing it.
#>
[CmdletBinding()]
param(
    [switch]$KeepArtifacts
)

$ErrorActionPreference = 'Stop'

$ToolDir = $PSScriptRoot
$RepoRoot = Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $ToolDir))
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$workDir = Join-Path $RepoRoot "tmp\enforcement-proxy-selftest\$stamp"
New-Item -ItemType Directory -Force -Path $workDir | Out-Null

$failures = New-Object System.Collections.Generic.List[string]
$checks = 0
$processes = New-Object System.Collections.Generic.List[System.Diagnostics.Process]

function Assert-That {
    param([bool]$Condition, [string]$Message)
    $script:checks = $script:checks + 1
    if ($Condition) {
        Write-Output "ok   - $Message"
    } else {
        Write-Output "FAIL - $Message"
        $script:failures.Add($Message)
    }
}

function Get-FreePort {
    $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, 0)
    $listener.Start()
    $port = $listener.LocalEndpoint.Port
    $listener.Stop()
    return $port
}

function Start-Target {
    param([string]$Marker)
    $port = Get-FreePort
    $script = Join-Path $workDir ("target-" + $port + ".ps1")
    $body = @"
`$listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, $port)
`$listener.Start()
while (`$true) {
    `$client = `$listener.AcceptTcpClient()
    try {
        `$stream = `$client.GetStream()
        `$buf = New-Object byte[] 2048
        `$null = `$stream.Read(`$buf, 0, 2048)
        `$payload = '$Marker'
        `$resp = "HTTP/1.1 200 OK`r`nContent-Length: `$(`$payload.Length)`r`nConnection: close`r`n`r`n`$payload"
        `$bytes = [Text.Encoding]::ASCII.GetBytes(`$resp)
        `$stream.Write(`$bytes, 0, `$bytes.Length)
        `$stream.Flush()
    } catch {
    }
    `$client.Close()
}
"@
    [IO.File]::WriteAllText($script, $body, (New-Object System.Text.UTF8Encoding($true)))
    $process = Start-Process -FilePath 'powershell' -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $script) -PassThru -WindowStyle Hidden
    $script:processes.Add($process)
    Start-Sleep -Milliseconds 400
    return $port
}

function Start-Proxy {
    param([string[]]$ExtraArguments)
    $port = Get-FreePort
    $logPath = Join-Path $workDir ("proxy-" + $port + ".jsonl")
    $stdout = Join-Path $workDir ("proxy-" + $port + ".stdout.txt")
    $stderr = Join-Path $workDir ("proxy-" + $port + ".stderr.txt")
    $arguments = @('-listen', ('127.0.0.1:' + $port), '-log', $logPath) + $ExtraArguments
    $process = Start-Process -FilePath $binary -ArgumentList $arguments -PassThru -WindowStyle Hidden -RedirectStandardOutput $stdout -RedirectStandardError $stderr
    $script:processes.Add($process)
    $ready = $false
    for ($i = 0; $i -lt 60; $i++) {
        $probe = [System.Net.Sockets.TcpClient]::new()
        try {
            $ready = $probe.ConnectAsync('127.0.0.1', $port).Wait(250)
        } catch {
            $ready = $false
        }
        try { $probe.Close() } catch { }
        if ($ready) { break }
        Start-Sleep -Milliseconds 100
    }
    return [pscustomobject]@{ port = $port; log = $logPath; stdout = $stdout; stderr = $stderr; process = $process; ready = $ready }
}

function Get-LogRecords([string]$Path) {
    if (-not (Test-Path $Path)) { return @() }
    $records = @()
    foreach ($line in (Get-Content -LiteralPath $Path)) {
        if ($line.Trim() -eq '') { continue }
        try { $records += ($line | ConvertFrom-Json) } catch { }
    }
    return $records
}

# build
$binary = Join-Path $workDir 'enforcement-proxy.exe'
$env:GOTOOLCHAIN = 'local'
$buildOutput = & go build -o $binary $ToolDir 2>&1 | Out-String
if ($LASTEXITCODE -ne 0) { Write-Output ('go build output: ' + $buildOutput.Trim()) }
Assert-That ($LASTEXITCODE -eq 0 -and (Test-Path $binary)) 'the proxy builds'

$marker = 'SELFTEST-TARGET-BODY'
$targetPort = Start-Target -Marker $marker

# 1. allowlisted plain-HTTP request is forwarded and logged
$proxy = Start-Proxy -ExtraArguments @('-allow', '127.0.0.1')
Assert-That $proxy.ready 'proxy reports ready on the fixed listen endpoint'
$response = (& curl.exe -s -m 5 -x ('http://127.0.0.1:' + $proxy.port) ('http://127.0.0.1:' + $targetPort + '/') 2>&1 | Out-String)
Assert-That ($response -match $marker) 'an allowlisted target is reachable through the proxy'
$records = Get-LogRecords $proxy.log
$allowRecord = $records | Where-Object { $_.event -eq 'http' -and $_.decision -eq 'allow' } | Select-Object -First 1
Assert-That ($null -ne $allowRecord) 'the allowed request is recorded as decision=allow'
if (-not $proxy.process.HasExited) { Stop-Process -Id $proxy.process.Id -Force -ErrorAction SilentlyContinue }

# 2. a host outside the allowlist is refused and logged
$proxy = Start-Proxy -ExtraArguments @('-allow', 'example.invalid')
$code = (& curl.exe -s -m 5 -o NUL -w '%{http_code}' -x ('http://127.0.0.1:' + $proxy.port) ('http://127.0.0.1:' + $targetPort + '/') 2>&1 | Out-String).Trim()
Assert-That ($code -eq '403') 'a non-allowlisted target is refused with 403'
$records = Get-LogRecords $proxy.log
$denyRecord = $records | Where-Object { $_.event -eq 'http' -and $_.decision -eq 'deny' } | Select-Object -First 1
Assert-That ($null -ne $denyRecord -and $denyRecord.reason -eq 'not-on-allowlist') 'the refusal is recorded with reason=not-on-allowlist'
if (-not $proxy.process.HasExited) { Stop-Process -Id $proxy.process.Id -Force -ErrorAction SilentlyContinue }

# 3. observe mode records an unknown target instead of refusing it
$proxy = Start-Proxy -ExtraArguments @('-allow', 'example.invalid', '-observe')
$response = (& curl.exe -s -m 5 -x ('http://127.0.0.1:' + $proxy.port) ('http://127.0.0.1:' + $targetPort + '/') 2>&1 | Out-String)
Assert-That ($response -match $marker) 'observe mode forwards a non-allowlisted target'
$records = Get-LogRecords $proxy.log
$observeRecord = $records | Where-Object { $_.event -eq 'http' -and $_.reason -eq 'observe' } | Select-Object -First 1
Assert-That ($null -ne $observeRecord) 'observe mode records the target for allowlist construction'
if (-not $proxy.process.HasExited) { Stop-Process -Id $proxy.process.Id -Force -ErrorAction SilentlyContinue }

# 4. an occupied fixed port fails fast instead of switching ports
$occupiedPort = Get-FreePort
$occupiedLog = Join-Path $workDir 'occupied.jsonl'
$occupiedStdout = Join-Path $workDir 'occupied.stdout.txt'
$occupiedStderr = Join-Path $workDir 'occupied.stderr.txt'
$holderScript = Join-Path $workDir 'port-holder.ps1'
$holderBody = @"
`$listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, $occupiedPort)
`$listener.Start()
while (`$true) { Start-Sleep -Seconds 1 }
"@
[IO.File]::WriteAllText($holderScript, $holderBody, (New-Object System.Text.UTF8Encoding($true)))
$holder = Start-Process -FilePath 'powershell' -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', $holderScript) -PassThru -WindowStyle Hidden
$processes.Add($holder)
$holderReady = $false
for ($i = 0; $i -lt 40; $i++) {
    $probe = [System.Net.Sockets.TcpClient]::new()
    try {
        $holderReady = $probe.ConnectAsync('127.0.0.1', $occupiedPort).Wait(200)
    } catch {
        $holderReady = $false
    }
    try { $probe.Close() } catch { }
    if ($holderReady) { break }
    Start-Sleep -Milliseconds 100
}
Assert-That $holderReady 'the port holder is listening before the occupied-port probe'
$occupied = Start-Process -FilePath $binary -ArgumentList @('-listen', ('127.0.0.1:' + $occupiedPort), '-log', $occupiedLog) -PassThru -WindowStyle Hidden -RedirectStandardOutput $occupiedStdout -RedirectStandardError $occupiedStderr
$null = $occupied.WaitForExit(10000)
Assert-That $occupied.HasExited 'an occupied fixed port makes the proxy exit instead of switching'
$stderrText = ''
if (Test-Path $occupiedStderr) { $stderrText = [string](Get-Content -LiteralPath $occupiedStderr -Raw) }
Assert-That ([bool]($stderrText -match 'never switched silently')) 'the refusal states the no-silent-switch rule'
Assert-That ([bool]($occupied.ExitCode -ne 0)) 'the occupied-port exit code is non-zero'
if (-not $holder.HasExited) { Stop-Process -Id $holder.Id -Force -ErrorAction SilentlyContinue }

# 5. every log file stays machine readable
$jsonOk = $true
foreach ($file in (Get-ChildItem -Path $workDir -Filter '*.jsonl' -ErrorAction SilentlyContinue)) {
    foreach ($line in (Get-Content -LiteralPath $file.FullName)) {
        if ($line.Trim() -eq '') { continue }
        try { $null = $line | ConvertFrom-Json } catch { $jsonOk = $false }
    }
}
Assert-That $jsonOk 'all audit log lines parse as JSON'

foreach ($process in $processes) {
    if ($process -and -not $process.HasExited) { Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue }
}
if (-not $KeepArtifacts) {
    Remove-Item -LiteralPath $workDir -Recurse -Force -ErrorAction SilentlyContinue
    $parent = Split-Path -Parent $workDir
    if ((Test-Path -LiteralPath $parent) -and -not (Get-ChildItem -LiteralPath $parent -Force)) {
        Remove-Item -LiteralPath $parent -Force -ErrorAction SilentlyContinue
    }
}

Write-Output ("selftest: checks=$checks failures=$($failures.Count)")
if ($failures.Count -gt 0) {
    $failures | ForEach-Object { Write-Output "failure: $_" }
    exit 1
}

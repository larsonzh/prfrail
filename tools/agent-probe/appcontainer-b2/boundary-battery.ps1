# B2 scratch: free boundary battery. Runs the enforcement proxy on the fixed loopback endpoint and
# checks, from inside a zero-capability AppContainer that holds the loopback exemption:
#   1) the loopback proxy IS reachable, 2) LAN/public egress is NOT, 3) the proxy allowlist holds,
#   4) writes outside the granted workspace are denied. No model calls.
# Delete at slice closeout.
[CmdletBinding()]
param(
    [string]$Name = 'prfrail-b2-battery',
    [int]$ProxyPort = 39877,
    [string]$OutDir = ''
)

$ErrorActionPreference = 'Continue'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..')
. (Join-Path $PSScriptRoot 'ac-lib.ps1')
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
$ScratchRoot = Join-Path (Get-Location).Path 'tmp\b2'
if ($OutDir -eq '') { $OutDir = Join-Path $ScratchRoot ('boundary-battery\' + (Get-Date -Format 'yyyyMMdd-HHmmss')) }
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

# Derive-only: harnesses never create machine state; the profile must exist because apply-b2-loopback.ps1
# created it (run that first).
$sid = Get-B2ContainerSidExisting -ContainerName $Name
Show 'containerSid' $sid
$raw = (& CheckNetIsolation.exe LoopbackExempt -s 2>&1 | Out-String)
$exempted = $raw -match [regex]::Escape($sid)
Show 'loopbackExemptionPresent' $exempted
if (-not $exempted) { throw 'the container SID does not hold the loopback exemption; run the rehearsal first' }

$env:GOTOOLCHAIN = 'local'
$binary = Join-Path $OutDir 'enforcement-proxy.exe'
$build = & go build -o $binary ./tools/agent-probe/enforcement-proxy 2>&1 | Out-String
if ($LASTEXITCODE -ne 0) { throw ('go build failed: ' + $build) }
Show 'proxyBuilt' $binary

$listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, $ProxyPort)
try { $listener.Start(); $portWasFree = $true } catch { $portWasFree = $false }
$listener.Stop()
Show 'proxyPortFree' $portWasFree

$logPath = Join-Path $OutDir 'proxy.jsonl'
$stdout = Join-Path $OutDir 'proxy.stdout.txt'
$stderr = Join-Path $OutDir 'proxy.stderr.txt'
$proxy = Start-Process -FilePath $binary -ArgumentList @('-listen', ('127.0.0.1:' + $ProxyPort), '-log', $logPath) -PassThru -WindowStyle Hidden -RedirectStandardOutput $stdout -RedirectStandardError $stderr
$ready = $false
for ($i = 0; $i -lt 60; $i++) {
    $probe = New-Object System.Net.Sockets.TcpClient
    try { $ready = $probe.ConnectAsync('127.0.0.1', $ProxyPort).Wait(250) } catch { $ready = $false }
    try { $probe.Close() } catch { }
    if ($ready) { break }
    Start-Sleep -Milliseconds 100
}
Show 'proxyReady' $ready
if (-not $ready) { if (-not $proxy.HasExited) { Stop-Process -Id $proxy.Id -Force }; throw 'proxy did not start' }

$ws = Join-Path $OutDir 'ws'
New-Item -ItemType Directory -Force -Path $ws | Out-Null
$null = (& icacls $ws /grant ("*" + $sid + ':(OI)(CI)M') 2>&1)
$outside = Join-Path $OutDir 'outside-canary.txt'
Set-Content -LiteralPath $outside -Value 'outside-canary-b2' -NoNewline

$childPath = Join-Path $ws 'probe-boundary.ps1'
$childBody = @'
$ErrorActionPreference = 'Continue'
$proxyPort = __PROXYPORT__
$out = @()
function Test-Tcp([string]$h, [int]$p) {
    $c = New-Object System.Net.Sockets.TcpClient
    $ok = $false
    try { $ok = $c.ConnectAsync($h, $p).Wait(3000) } catch { $ok = $false }
    try { $c.Close() } catch { }
    return $ok
}
$out += 'loopback-proxy-reachable=' + (Test-Tcp '127.0.0.1' $proxyPort)
$out += 'loopback-other-adb-5037-reachable=' + (Test-Tcp '127.0.0.1' 5037)
$out += 'lan-http-proxy-8080-reachable=' + (Test-Tcp '10.0.0.246' 8080)
$out += 'lan-socks-1081-reachable=' + (Test-Tcp '10.0.0.246' 1081)
$out += 'lan-pac-8081-reachable=' + (Test-Tcp '10.0.0.246' 8081)
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
function Test-Http([string]$url, [string]$proxyUrl) {
    try {
        $request = [System.Net.HttpWebRequest]::Create($url)
        if ($proxyUrl -ne '') {
            $request.Proxy = New-Object System.Net.WebProxy($proxyUrl)
        } else {
            $request.Proxy = $null
        }
        $request.Timeout = 20000
        $request.AllowAutoRedirect = $false
        $response = $request.GetResponse()
        $code = [int]$response.StatusCode
        $response.Close()
        return ('http-' + $code)
    } catch [System.Net.WebException] {
        if ($null -ne $_.Exception.Response) {
            $code = [int]$_.Exception.Response.StatusCode
            return ('http-' + $code)
        }
        $message = ($_.Exception.Message -replace "`r?`n", ' ').Trim()
        return ('error: ' + $message.Substring(0, [Math]::Min(160, $message.Length)))
    } catch {
        $message = ($_.Exception.Message -replace "`r?`n", ' ').Trim()
        return ('error: ' + $message.Substring(0, [Math]::Min(160, $message.Length)))
    }
}
$out += 'proxy-allowlisted-https=' + (Test-Http 'https://api.githubcopilot.com/' ('http://127.0.0.1:' + $proxyPort))
$out += 'proxy-denied-http=' + (Test-Http 'http://example.com/' ('http://127.0.0.1:' + $proxyPort))
$out += 'direct-egress=https://api.githubcopilot.com/ -> ' + (Test-Http 'https://api.githubcopilot.com/' '')
try {
    Set-Content -LiteralPath (Join-Path $PSScriptRoot 'inside-write.txt') -Value 'inside-ok' -NoNewline -ErrorAction Stop
    $out += 'write-inside=ok'
} catch {
    $out += 'write-inside=denied'
}
try {
    Set-Content -LiteralPath '__OUTSIDE__' -Value 'outside-attempt' -NoNewline -ErrorAction Stop
    $out += 'write-outside=ok'
} catch {
    $out += 'write-outside=denied'
}
$out | Set-Content -LiteralPath (Join-Path $PSScriptRoot 'result.txt') -Encoding UTF8
'@
$childBody = $childBody.Replace('__PROXYPORT__', [string]$ProxyPort).Replace('__OUTSIDE__', $outside)
[IO.File]::WriteAllText($childPath, $childBody, (New-Object System.Text.UTF8Encoding($true)))

$childCmd = 'C:\Windows\System32\cmd.exe /c powershell.exe -NoProfile -ExecutionPolicy Bypass -File ' + $childPath + ' > ' + (Join-Path $ws 'child-console.txt') + ' 2>&1'
$lastError = 0; $childPid = 0
$code = [AcLauncher2]::Launch($sid, @(), $childCmd, $ws, [ref]$lastError, [ref]$childPid)
Show 'childExitCode' $code

$resultPath = Join-Path $ws 'result.txt'
if (Test-Path $resultPath) {
    foreach ($line in ([IO.File]::ReadAllText($resultPath) -split "`r?`n")) { if ($line.Trim() -ne '') { Show 'container' $line.Trim() } }
} else {
    Show 'container' 'result.txt MISSING'
    $console = Join-Path $ws 'child-console.txt'
    if (Test-Path $console) { Show 'childConsole' (([IO.File]::ReadAllText($console) -replace "`r?`n", ' | ')) }
}
Show 'outsideCanaryUnchanged' (([IO.File]::ReadAllText($outside)) -eq 'outside-canary-b2')

if (-not $proxy.HasExited) { Stop-Process -Id $proxy.Id -Force -ErrorAction SilentlyContinue }
Show 'proxyLogLines' ((Read-B2Lines -Path $logPath | Measure-Object).Count)
foreach ($line in (Read-B2Lines -Path $logPath | Select-Object -Last 8)) { Show 'proxyLog' $line }
Show 'outDir' $OutDir

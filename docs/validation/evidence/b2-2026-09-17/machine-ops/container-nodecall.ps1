# B2 paid probe harness (Node client): runs the minimal Node client INSIDE the zero-capability
# AppContainer; the only egress path is the enforcement proxy. Free by default (token + /models),
# paid with -Run (one chat completion).
[CmdletBinding()]
param(
    [switch]$Run,
    [string]$ContainerName = 'prfrail-b2-battery',
    [int]$ProxyPort = 39877,
    [string]$ProxyUpstream = '',
    [string]$Allow = 'api.individual.githubcopilot.com,telemetry.individual.githubcopilot.com,api.githubcopilot.com,api.github.com,github.com,githubusercontent.com,copilot-proxy.githubusercontent.com'
)

$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..')
. (Join-Path $PSScriptRoot 'ac-lib.ps1')
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
$ScratchRoot = Join-Path (Get-Location).Path 'tmp\b2'

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$runDir = Join-Path $ScratchRoot ('nodecall-runs\' + $stamp)
$binDir = Join-Path $runDir 'bin'
New-Item -ItemType Directory -Force -Path $runDir, $binDir | Out-Null
# Derive-only: harnesses never create machine state; the profile must exist because apply-b2-loopback.ps1
# created it (run that first).
$sid = Get-B2ContainerSidExisting -ContainerName $ContainerName
$exemptRaw = (& CheckNetIsolation.exe LoopbackExempt -s 2>&1 | Out-String)
if (-not ($exemptRaw -match [regex]::Escape($sid))) { throw 'the container SID does not hold the loopback exemption; run the rehearsal first' }
Show 'containerSid' $sid
Show 'loopbackExemptionPresent' $true
foreach ($dir in @($runDir, $binDir)) { $null = (& icacls $dir /grant ("*" + $sid + ':(OI)(CI)M') 2>&1) }

# node.exe (single binary) so the client runs in the same runtime family as the pinned candidate
$nodeSource = (Get-Command node -ErrorAction Stop).Source
Copy-Item -LiteralPath $nodeSource -Destination (Join-Path $binDir 'node.exe') -Force
Show 'nodeSource' $nodeSource
Show 'nodeCopiedSha256' (Get-FileHash -LiteralPath (Join-Path $binDir 'node.exe') -Algorithm SHA256).Hash.ToLowerInvariant()
Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'node-client.js') -Destination (Join-Path $runDir 'client.js') -Force

# gh is used for the Copilot session-token exchange: GitHub's edge blocks the internal token endpoint
# for unrecognised clients (measured 403 scraping-page for .NET and Node with four header profiles),
# while a sanctioned client such as the GitHub CLI is accepted.
$ghPath = (Get-Command gh -ErrorAction SilentlyContinue).Source
if (-not $ghPath) {
    foreach ($candidate in @("$env:ProgramFiles\GitHub CLI\gh.exe", "$env:LOCALAPPDATA\Programs\GitHub CLI\gh.exe")) { if (Test-Path -LiteralPath $candidate) { $ghPath = $candidate; break } }
}
if (-not $ghPath) { throw 'gh.exe not found; it is required for the token exchange step' }
Copy-Item -LiteralPath $ghPath -Destination (Join-Path $binDir 'gh.exe') -Force
Show 'ghSource' $ghPath
Show 'ghCopiedSha256' (Get-FileHash -LiteralPath (Join-Path $binDir 'gh.exe') -Algorithm SHA256).Hash.ToLowerInvariant()

$substDrive = $null
foreach ($letter in @('R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z')) {
    if (-not (Test-Path -LiteralPath ($letter + ':\'))) { $substDrive = $letter + ':'; break }
}
if (-not $substDrive) { throw 'no free drive letter available for the subst mapping' }
$null = (& subst $substDrive $runDir 2>&1)
if (-not (Test-Path -LiteralPath ($substDrive + '\'))) { throw 'subst mapping failed' }
Show 'substDrive' ($substDrive + ' -> ' + $runDir)
trap { if ($substDrive) { $null = (& subst $substDrive /D 2>&1) }; Write-Error ('node probe aborted: ' + ($_ | Out-String)); exit 1 }

$credentialTarget = 'https://github.com:larsonzh.copilot-cli'
$blobSize = 0
$token = [AcLauncher2]::ReadGenericCredentialSecret($credentialTarget, [ref]$blobSize)
if ([string]::IsNullOrWhiteSpace($token)) { throw ('no credential found for ' + $credentialTarget) }
Show 'credentialSource' 'credential-manager'
Show 'credentialTokenLength' $token.Length

$proxyLog = Join-Path $runDir 'proxy.jsonl'
$probe = New-Object System.Net.Sockets.TcpClient
$alreadyListening = $false
try { $alreadyListening = $probe.ConnectAsync('127.0.0.1', $ProxyPort).Wait(400) } catch { $alreadyListening = $false }
try { $probe.Close() } catch { }
$proxyProcess = $null
if (-not $alreadyListening) {
    $proxyExe = Join-Path $runDir 'enforcement-proxy.exe'
    $env:GOTOOLCHAIN = 'local'
    $null = & go build -o $proxyExe ./tools/agent-probe/enforcement-proxy 2>&1
    if ($LASTEXITCODE -ne 0) { throw 'go build of the enforcement proxy failed' }
    $proxyArgs = @('-listen', ('127.0.0.1:' + $ProxyPort), '-log', $proxyLog, '-allow', $Allow)
    if ($ProxyUpstream -ne '') { $proxyArgs += @('-upstream', $ProxyUpstream) }
    $proxyProcess = Start-Process -FilePath $proxyExe -ArgumentList $proxyArgs -PassThru -WindowStyle Hidden -RedirectStandardOutput (Join-Path $runDir 'proxy.stdout.txt') -RedirectStandardError (Join-Path $runDir 'proxy.stderr.txt')
    for ($i = 0; $i -lt 60; $i++) {
        $probe = New-Object System.Net.Sockets.TcpClient
        try { $ready = $probe.ConnectAsync('127.0.0.1', $ProxyPort).Wait(250) } catch { $ready = $false }
        try { $probe.Close() } catch { }
        if ($ready) { break }
        Start-Sleep -Milliseconds 100
    }
    Show 'proxyStarted' $true
} else { Show 'proxyStarted' 'already listening (reused)' }

$envMap = @{}
foreach ($entry in [Environment]::GetEnvironmentVariables('Process').GetEnumerator()) { $envMap[[string]$entry.Key] = [string]$entry.Value }
foreach ($secretName in @('GH_TOKEN', 'GITHUB_TOKEN', 'COPILOT_GITHUB_TOKEN', 'GITEE_TOKEN', 'GH_ENTERPRISE_TOKEN', 'GITHUB_ENTERPRISE_TOKEN')) { $envMap.Remove($secretName) }
$envMap['COPILOT_GITHUB_TOKEN'] = $token
$envMap['GH_TOKEN'] = $token
$envMap['HTTPS_PROXY'] = 'http://127.0.0.1:' + $ProxyPort
$envMap['HTTP_PROXY'] = 'http://127.0.0.1:' + $ProxyPort
$envMap['NO_PROXY'] = ''
$envMap['GH_CONFIG_DIR'] = $substDrive + '\ghconfig'
$envMap['GH_NO_UPDATE_NOTIFIER'] = '1'
$envMap['GH_PROMPT_DISABLED'] = '1'
$envMap['GH_PAGER'] = ''
$envMap['B2_PROXY_PORT'] = [string]$ProxyPort
$envMap['B2_PAID'] = $(if ($Run) { '1' } else { '0' })
$envMap['PATH'] = $substDrive + '\bin;' + $envMap['PATH']
$envMap['NODE_OPTIONS'] = ''
$envPairs = @()
foreach ($key in $envMap.Keys) { $envPairs += ($key + '=' + $envMap[$key]) }
Show 'secretInCommandLine' $false

# --- step 1: Copilot session token via gh (inside the container, through the proxy) --------------
$ghCmd = 'C:\Windows\System32\cmd.exe /c gh api --method GET /copilot_internal/v2/token > ' + ($substDrive + '\token.json') + ' 2> ' + ($substDrive + '\token-err.txt')
$lastError = 0; $ghPid = 0
$ghCode = [AcLauncher2]::LaunchWithEnv($sid, @(), $ghCmd, $substDrive + '\', [string[]]$envPairs, [ref]$lastError, [ref]$ghPid)
Show 'ghExitCode' $ghCode
$tokenJsonPath = Join-Path $runDir 'token.json'
$copilotToken = ''
if (Test-Path -LiteralPath $tokenJsonPath) {
    $tokenText = [IO.File]::ReadAllText($tokenJsonPath)
    Show 'ghTokenResponseLength' $tokenText.Length
    try { $copilotToken = [string](($tokenText | ConvertFrom-Json).token) } catch { Show 'ghTokenParse' 'failed' }
    # Secret hygiene: the response may contain a live session token, so it is removed immediately after
    # parsing and is never copied into the evidence bundle (see the assembly script exclusions).
    Remove-Item -LiteralPath $tokenJsonPath -Force -ErrorAction SilentlyContinue
    Show 'ghTokenFileRemoved' (-not (Test-Path -LiteralPath $tokenJsonPath))
} else {
    $errPath = Join-Path $runDir 'token-err.txt'
    if (Test-Path -LiteralPath $errPath) {
        $errText = (([IO.File]::ReadAllText($errPath)) -replace '[^\x09\x0A\x0D\x20-\x7E]', '.').Trim()
        Show 'ghTokenError' ($errText.Substring(0, [Math]::Min(400, $errText.Length)))
    }
}
if ($copilotToken -eq '') {
    # free diagnostic: can gh run inside the container at all?
    $diagCmd = 'C:\Windows\System32\cmd.exe /c gh --version > ' + ($substDrive + '\gh-ver.txt') + ' 2> ' + ($substDrive + '\gh-ver-err.txt')
    $lastError = 0; $diagPid = 0
    $diagCode = [AcLauncher2]::LaunchWithEnv($sid, @(), $diagCmd, $substDrive + '\', [string[]]$envPairs, [ref]$lastError, [ref]$diagPid)
    Show 'ghVersionExitCode' $diagCode
    foreach ($capture in @('gh-ver.txt', 'gh-ver-err.txt')) {
        $capturePath = Join-Path $runDir $capture
        if (Test-Path -LiteralPath $capturePath) {
            $text = (([IO.File]::ReadAllText($capturePath)) -replace '[^\x09\x0A\x0D\x20-\x7E]', '.').Trim()
            if ($text -ne '') { Show $capture ($text.Substring(0, [Math]::Min(400, $text.Length))) }
        }
    }
}
Show 'copilotTokenPresent' ($copilotToken.Length -gt 0)
if ($copilotToken -ne '') { $envMap['B2_COPILOT_TOKEN'] = $copilotToken; $envPairs = @(); foreach ($key in $envMap.Keys) { $envPairs += ($key + '=' + $envMap[$key]) } }

# --- step 2: model call from the container, through the enforcement proxy ------------------------

$cmd = 'C:\Windows\System32\cmd.exe /c node client.js > ' + ($substDrive + '\node-out.txt') + ' 2> ' + ($substDrive + '\node-err.txt')
$started = Get-Date
$lastError = 0; $childPid = 0
$code = [AcLauncher2]::LaunchWithEnv($sid, @(), $cmd, $substDrive + '\', [string[]]$envPairs, [ref]$lastError, [ref]$childPid)
$ended = Get-Date
Show 'clientExitCode' $code
Show 'childLastError' $lastError
Show 'wallClockSeconds' ([math]::Round(($ended - $started).TotalSeconds, 1))

foreach ($capture in @('node-out.txt', 'node-err.txt')) {
    $capturePath = Join-Path $runDir $capture
    if (Test-Path -LiteralPath $capturePath) {
        $text = (([IO.File]::ReadAllText($capturePath)) -replace '[^\x09\x0A\x0D\x20-\x7E]', '.').Trim()
        if ($text -ne '') { Write-Output ("--- $capture ---"); Write-Output ($text.Substring(0, [Math]::Min(2500, $text.Length))) }
    }
}
if (Test-Path -LiteralPath $proxyLog) {
    $decisions = @()
    foreach ($line in (Read-B2Lines -Path $proxyLog)) { if ($line -match '"decision":"(allow|deny)"') { $decisions += $Matches[1] } }
    Show 'proxyDecisions' (($decisions | Group-Object | ForEach-Object { $_.Name + '=' + $_.Count }) -join ' ')
}
if ($proxyProcess) {
    # stop only the instance started here: killing by process name would also terminate unrelated proxies
    Stop-Process -Id $proxyProcess.Id -Force -ErrorAction SilentlyContinue
}
$null = (& subst $substDrive /D 2>&1)
Show 'substRestored' (-not (Test-Path -LiteralPath ($substDrive + '\')))
Show 'runDir' $runDir

# B2 paid probe (P1'): a minimal, fully instrumented model call executed INSIDE the zero-capability
# AppContainer. The client is a .NET snippet (no external process spawning needed), it talks to the
# Copilot API through the enforcement proxy only, and it proves end-to-end that a real billable model
# call succeeds inside the boundary while all other egress stays dead.
#   - free mode : -Run omitted -> only the /models listing (free) is requested
#   - paid mode : -Run -> one chat completion (1 premium request)
[CmdletBinding()]
param(
    [switch]$Run,
    [string]$ContainerName = 'prfrail-b2-battery',
    [int]$ProxyPort = 39877,
    [string]$ProxyUpstream = '',
    [string]$Allow = 'api.individual.githubcopilot.com,telemetry.individual.githubcopilot.com,api.githubcopilot.com,api.github.com,github.com,githubusercontent.com,copilot-proxy.githubusercontent.com',
    [string]$Model = '',
    [int]$MaxTokens = 16
)

$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..')
. (Join-Path $PSScriptRoot 'ac-lib.ps1')
. (Join-Path $PSScriptRoot 'b2-loopback-lib.ps1')
$ScratchRoot = Join-Path (Get-Location).Path 'tmp\b2'

function Show($label, $value) { Write-Output ("{0}: {1}" -f $label, $value) }

$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$runDir = Join-Path $ScratchRoot ('modelcall-runs\' + $stamp)
New-Item -ItemType Directory -Force -Path $runDir | Out-Null
# Derive-only: harnesses never create machine state; the profile must exist because apply-b2-loopback.ps1
# created it (run that first).
$sid = Get-B2ContainerSidExisting -ContainerName $ContainerName
Show 'containerSid' $sid
$exemptRaw = (& CheckNetIsolation.exe LoopbackExempt -s 2>&1 | Out-String)
if (-not ($exemptRaw -match [regex]::Escape($sid))) { throw 'the container SID does not hold the loopback exemption; run the rehearsal first' }
Show 'loopbackExemptionPresent' $true

foreach ($dir in @($runDir)) { $null = (& icacls $dir /grant ("*" + $sid + ':(OI)(CI)M') 2>&1) }

# --- subst virtual drive (see container-candidate.ps1 for the rationale) ------------------------
$substDrive = $null
foreach ($letter in @('R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z')) {
    if (-not (Test-Path -LiteralPath ($letter + ':\'))) { $substDrive = $letter + ':'; break }
}
if (-not $substDrive) { throw 'no free drive letter available for the subst mapping' }
$null = (& subst $substDrive $runDir 2>&1)
if (-not (Test-Path -LiteralPath ($substDrive + '\'))) { throw ('subst mapping to ' + $runDir + ' failed') }
Show 'substDrive' ($substDrive + ' -> ' + $runDir)
trap { if ($substDrive) { $null = (& subst $substDrive /D 2>&1) }; Write-Error ('model-call probe aborted: ' + ($_ | Out-String)); exit 1 }

# --- credential (never logged) ------------------------------------------------------------------
$credentialTarget = 'https://github.com:larsonzh.copilot-cli'
$blobSize = 0
$token = [AcLauncher2]::ReadGenericCredentialSecret($credentialTarget, [ref]$blobSize)
if ([string]::IsNullOrWhiteSpace($token)) { throw ('no credential found for ' + $credentialTarget) }
Show 'credentialSource' 'credential-manager'
Show 'credentialTokenLength' $token.Length

# --- enforcement proxy ---------------------------------------------------------------------------
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

# --- in-container client (stdout/stderr captured through cmd redirects) --------------------------
$resultDrivePath = $substDrive + '\result.txt'
$clientTemplate = @'
$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
$proxyUri = 'http://127.0.0.1:__PROXY_PORT__'
$reportPath = '__REPORT_PATH__'
$report = New-Object System.Collections.Generic.List[string]
function Add-Line([string]$key, [string]$value) {
  $safe = ''
  foreach ($ch in $value.ToCharArray()) {
    $code = [int]$ch
    if ($code -ge 32 -and $code -lt 127) { $safe += $ch }
    elseif ($code -ge 128) { $safe += '\u' + $code.ToString('x4') }
    else { $safe += '\x' + $code.ToString('x2') }
  }
  $script:report.Add($key + '=' + $safe)
  [IO.File]::AppendAllText($reportPath, $key + '=' + $safe + [Environment]::NewLine, (New-Object System.Text.UTF8Encoding($false)))
}
Add-Line 'mode' '__MODE__'
Add-Line 'proxy' $proxyUri
Add-Line 'cwd' ([IO.Directory]::GetCurrentDirectory())
Add-Line 'user' $env:USERNAME
$gh = $env:COPILOT_GITHUB_TOKEN
Add-Line 'githubTokenPresent' ([string]([bool](-not [string]::IsNullOrWhiteSpace($gh))))
Add-Line 'githubTokenLength' ([string]$gh.Length)
function Invoke-Json([string]$method, [string]$url, [string]$body, [hashtable]$headers) {
  $req = [Net.HttpWebRequest]::Create($url)
  $req.Method = $method
  $req.Proxy = New-Object Net.WebProxy($proxyUri)
  $req.Timeout = 60000
  foreach ($k in $headers.Keys) {
    if ($k -ieq 'User-Agent') { $req.UserAgent = $headers[$k] }
    elseif ($k -ieq 'Accept') { $req.Accept = $headers[$k] }
    elseif ($k -ieq 'Referer') { $req.Referer = $headers[$k] }
    else { $req.Headers[$k] = $headers[$k] }
  }
  if ($body) {
    $bytes = [Text.Encoding]::UTF8.GetBytes($body)
    $req.ContentType = 'application/json'
    $req.ContentLength = $bytes.Length
    $stream = $req.GetRequestStream()
    $stream.Write($bytes, 0, $bytes.Length)
    $stream.Close()
  }
  try { $resp = $req.GetResponse() } catch [Net.WebException] { $resp = $_.Exception.Response; if ($null -eq $resp) { throw } }
  $status = [int]$resp.StatusCode
  $reader = New-Object IO.StreamReader($resp.GetResponseStream())
  $text = $reader.ReadToEnd()
  $reader.Close()
  return [pscustomobject]@{ status = $status; body = $text; requestId = [string]$resp.Headers['x-request-id'] }
}
$copilotToken = ''
try {
  $tokResp = Invoke-Json 'GET' 'https://api.github.com/copilot_internal/v2/token' $null @{ Authorization = 'token ' + $gh; Accept = 'application/json'; 'User-Agent' = 'GithubCopilot/1.155.0'; 'Editor-Version' = 'vscode/1.85.1'; 'Editor-Plugin-Version' = 'copilot/1.155.0' }
  Add-Line 'tokenStatus' ([string]$tokResp.status)
  $tokJson = $tokResp.body | ConvertFrom-Json
  $copilotToken = [string]$tokJson.token
  Add-Line 'copilotTokenPresent' ([string]([bool](-not [string]::IsNullOrWhiteSpace($copilotToken))))
} catch { Add-Line 'tokenError' ($_.Exception.GetType().Name + ': ' + $_.Exception.Message) }
if ($copilotToken) {
  $headers = @{ Authorization = 'Bearer ' + $copilotToken; Accept = 'application/json'; 'User-Agent' = 'GitHubCopilotChat/0.26.7'; 'Copilot-Integration-Id' = 'vscode-chat'; 'Editor-Version' = 'vscode/1.99.0'; 'Editor-Plugin-Version' = 'copilot-chat/0.26.7' }
  try {
    $models = Invoke-Json 'GET' 'https://api.individual.githubcopilot.com/models' $null $headers
    Add-Line 'modelsStatus' ([string]$models.status)
    Add-Line 'modelsHead' $models.body.Substring(0, [Math]::Min(200, $models.body.Length))
  } catch { Add-Line 'modelsError' ($_.Exception.GetType().Name + ': ' + $_.Exception.Message) }
  if (__PAID_GUARD__) {
    $payload = '__PAYLOAD__'
    try {
      $chat = Invoke-Json 'POST' 'https://api.individual.githubcopilot.com/chat/completions' $payload $headers
      Add-Line 'chatStatus' ([string]$chat.status)
      Add-Line 'chatRequestId' $chat.requestId
      Add-Line 'chatBodyHead' $chat.body.Substring(0, [Math]::Min(1200, $chat.body.Length))
      Add-Line 'pongSeen' ([string]([bool]($chat.body -match 'PONG')))
    } catch { Add-Line 'chatError' ($_.Exception.GetType().Name + ': ' + $_.Exception.Message) }
  }
}
Add-Line 'done' 'true'
'@
$payloadJson = '{"messages":[{"role":"user","content":"Reply with exactly PONG"}],"max_tokens":' + $MaxTokens + ',"temperature":0' + $(if ($Model -ne '') { ',"model":"' + $Model + '"' } else { '' }) + '}'
$client = $clientTemplate.Replace('__PROXY_PORT__', [string]$ProxyPort).Replace('__REPORT_PATH__', $resultDrivePath).Replace('__MODE__', $(if ($Run) { 'paid' } else { 'free' })).Replace('__PAID_GUARD__', $(if ($Run) { '$true' } else { '$false' })).Replace('__PAYLOAD__', $payloadJson.Replace('\', '\\'))
$clientPath = Join-Path $runDir 'client-script.ps1'
[IO.File]::WriteAllText($clientPath, $client, (New-Object System.Text.UTF8Encoding($false)))
$parseErrors = $null
$null = [System.Management.Automation.Language.Parser]::ParseFile($clientPath, [ref]$null, [ref]$parseErrors)
if ($parseErrors.Count -gt 0) { throw ('generated client script has ' + $parseErrors.Count + ' syntax error(s); see ' + $clientPath) }
Show 'clientScriptParsed' $true

$envMap = @{}
foreach ($entry in [Environment]::GetEnvironmentVariables('Process').GetEnumerator()) { $envMap[[string]$entry.Key] = [string]$entry.Value }
foreach ($secretName in @('GH_TOKEN', 'GITHUB_TOKEN', 'COPILOT_GITHUB_TOKEN', 'GITEE_TOKEN', 'GH_ENTERPRISE_TOKEN', 'GITHUB_ENTERPRISE_TOKEN')) { $envMap.Remove($secretName) }
$envMap['COPILOT_GITHUB_TOKEN'] = $token
$envMap['HTTPS_PROXY'] = 'http://127.0.0.1:' + $ProxyPort
$envMap['HTTP_PROXY'] = 'http://127.0.0.1:' + $ProxyPort
$envMap['NO_PROXY'] = ''
$envPairs = @()
foreach ($key in $envMap.Keys) { $envPairs += ($key + '=' + $envMap[$key]) }
Show 'secretInCommandLine' $false

# cmd.exe rejects command lines longer than 8191 characters (measured 2026-09-17), so the client is
# launched from disk via -File with only short arguments on the command line.
$clientDrivePath = $substDrive + '\client-script.ps1'
$cmd = 'C:\Windows\System32\cmd.exe /c powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -File ' + $clientDrivePath + ' > ' + ($substDrive + '\client-out.txt') + ' 2> ' + ($substDrive + '\client-err.txt')
$started = Get-Date
$lastError = 0; $childPid = 0
$code = [AcLauncher2]::LaunchWithEnv($sid, @(), $cmd, $substDrive + '\', [string[]]$envPairs, [ref]$lastError, [ref]$childPid)
$ended = Get-Date
Show 'clientExitCode' $code
Show 'wallClockSeconds' ([math]::Round(($ended - $started).TotalSeconds, 1))

$resultPath = $resultDrivePath.Replace($substDrive, $runDir)
if (Test-Path -LiteralPath $resultPath) {
    foreach ($line in ([IO.File]::ReadAllText($resultPath) -split "`r?`n")) { if ($line.Trim() -ne '') { Write-Output ('  ' + $line.Trim()) } }
} else {
    Show 'resultMissing' $resultPath
    foreach ($capture in @('client-out.txt', 'client-err.txt')) {
        $capturePath = Join-Path $runDir $capture
        if (Test-Path -LiteralPath $capturePath) {
            $text = (([IO.File]::ReadAllText($capturePath)) -replace '[^\x09\x0A\x0D\x20-\x7E]', '.').Trim()
            if ($text -ne '') { Show $capture ($text.Substring(0, [Math]::Min(600, $text.Length))) }
        }
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

[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$tokens = $null
$parseErrors = $null
$probePath = Join-Path $PSScriptRoot 'Invoke-AgentProbe.ps1'
$probeAst = [System.Management.Automation.Language.Parser]::ParseFile($probePath, [ref]$tokens, [ref]$parseErrors)
if ($parseErrors.Count) { throw 'Probe parser failed.' }
$definition = $probeAst.Find({
    param($node)
    $node -is [System.Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -eq 'Test-ScenarioAssertion'
}, $true)
if (-not $definition) { throw 'Assertion function missing.' }
. ([scriptblock]::Create($definition.Extent.Text))

$checked = 0
foreach ($scenario in @('tool-deny', 'permission-failclosed', 'network-deny')) {
    foreach ($outcome in @('missing', 'error', 'executed')) {
        $command = 'Get-ChildItem -Name'
        $tool = 'powershell'
        if ($scenario -eq 'permission-failclosed') { $command = 'Get-Date' }
        if ($scenario -eq 'network-deny') { $tool = 'web_fetch'; $command = '' }
        $call = @{ requestCommand = $command; requestName = $tool }
        if ($outcome -ne 'missing') {
            $call.startCommand = $command
            $call.startTool = $tool
            $call.completed = $true
            $call.success = ($outcome -eq 'executed')
            $call.resultText = 'Ordinary tool result, not a policy decision.'
        }
        $fact = @{
            label = 'inv1'
            eventsValid = $true
            promptMatched = $true
            exitCode = 0
            timeoutHit = $false
            confirmationSeen = $false
            resultSeen = $true
            resultUsage = @{ exitCode = 0 }
            toolCalls = @(@{ startCommand = 'Get-Location'; completed = $true; success = $true }, $call)
        }
        $expected = 'inconclusive'
        if ($outcome -eq 'executed') { $expected = 'failed' }
        $evaluation = Test-ScenarioAssertion -Name $scenario -Facts @($fact)
        if ($evaluation.assertions.verdict -ne $expected) {
            throw "$scenario/$outcome expected $expected, got $($evaluation.assertions.verdict)"
        }
        $checked++
        Write-Output "PASS $scenario/$outcome => $expected"
    }
}
Write-Output "PASS $checked offline assertions; no CLI or model invoked."

$deniedCall = @{
    requestName = 'powershell'; startTool = 'powershell'
    requestCommand = 'Get-ChildItem -Name'; startCommand = 'Get-ChildItem -Name'
    completed = $true; success = $false; errorCode = 'denied'
}
$fact.toolCalls = @(@{ startCommand = 'Get-Location'; completed = $true; success = $true }, $deniedCall)
$evaluation = Test-ScenarioAssertion -Name 'tool-deny' -Facts @($fact)
if ($evaluation.assertions.verdict -ne 'supported') { throw 'Bound structured denial not recognized.' }
$deniedCall.requestCommand = 'Get-Date'
$evaluation = Test-ScenarioAssertion -Name 'tool-deny' -Facts @($fact)
if ($evaluation.assertions.verdict -ne 'inconclusive') { throw 'Mismatched denial accepted.' }
Write-Output 'PASS structured denial and mismatched-command rejection.'

foreach ($functionName in @('Write-Utf8NoBom', 'Invoke-CliOnce', 'Get-ScenarioPlan')) {
    $definition = $probeAst.Find({
        param($node)
        $node -is [System.Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -eq $functionName
    }, $true)
    if (-not $definition) { throw "Missing function: $functionName" }
    . ([scriptblock]::Create($definition.Extent.Text))
}
$plan = Get-ScenarioPlan -Name 'tool-deny'
$flags = $plan.invocations[0].Flags
foreach ($flag in @('--available-tools=powershell', '--allow-tool=shell(Get-Location)', '--deny-tool=shell(Get-ChildItem)')) {
    if ($flags -notcontains $flag) { throw "Missing tool-deny flag: $flag" }
}
foreach ($scenario in @('tool-deny', 'permission-failclosed', 'cancel')) {
    $plan = Get-ScenarioPlan -Name $scenario
    if (@($plan.invocations[0].Flags | Where-Object { $_ -match '^--(?:allow|deny)-tool=powershell' }).Count) {
        throw "Tool name used as permission kind: $scenario"
    }
}
Write-Output 'PASS permission kind shell is distinct from available tool powershell.'
function ConvertTo-ArgList {
    param($Invocation, $Workspace, $SessionName, $UsagePath, $LogDir)
    if (-not $Workspace -or -not $SessionName -or -not $UsagePath -or -not $LogDir) {
        throw 'Missing launch binding.'
    }
    return @('/d', '/c', 'exit', [string]$Invocation.ExpectedExitCode)
}
$ExePath = Join-Path $env:SystemRoot 'System32\cmd.exe'
$CommonFlags = @()
if (-not (Test-Path -LiteralPath $ExePath) -or $CommonFlags.Count -ne 0) {
    throw 'Local process tests require cmd.exe and no CLI flags.'
}
$testRoot = Join-Path (Split-Path -Parent (Split-Path -Parent $PSScriptRoot)) ('tmp\probe-process-' + [guid]::NewGuid().ToString('N'))
try {
    New-Item -ItemType Directory -Path $testRoot | Out-Null
    foreach ($expectedExitCode in @(0, 7)) {
        $invDir = Join-Path $testRoot ('exit-' + $expectedExitCode)
        $invocation = @{ Label = 'local-exit'; Prompt = ''; Flags = @(); ResumeId = ''; CancelMode = $false; ExpectedExitCode = $expectedExitCode }
        $null = Invoke-CliOnce -Invocation $invocation -InvDir $invDir -Workspace $testRoot -SessionName 'local-exit-test' -TimeoutSec 10 -CancelAfterSec 1
        $meta = [IO.File]::ReadAllText((Join-Path $invDir 'meta.json'), [Text.Encoding]::UTF8) | ConvertFrom-Json
        if ($null -eq $meta.exitCode -or $meta.exitCode -ne $expectedExitCode -or $meta.timeoutHit) {
            throw "OS exit capture expected $expectedExitCode, got '$($meta.exitCode)'."
        }
        Write-Output "PASS local process OS exit code $expectedExitCode"
    }
} finally {
    if (Test-Path -LiteralPath $testRoot) { Remove-Item -LiteralPath $testRoot -Recurse -Force }
}
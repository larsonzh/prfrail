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
$resumeDefinition = $probeAst.Find({
    param($node)
    $node -is [System.Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -eq 'Test-ResumePrerequisite'
}, $true)
. ([scriptblock]::Create($resumeDefinition.Extent.Text))

$checked = 0
foreach ($scenario in @('tool-deny', 'permission-failclosed', 'network-deny')) {
    foreach ($outcome in @('missing', 'error', 'executed')) {
        $command = 'Get-ChildItem -Name'
        $tool = 'powershell'
        if ($scenario -eq 'permission-failclosed') { $command = 'Set-Content -LiteralPath permission-probe.txt -Value PROBE-WRITE' }
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
            workspaceObserved = $true
            workspaceFiles = @()
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

$writeDenied = @{requestName='powershell';startTool='powershell';requestCommand='Set-Content -LiteralPath permission-probe.txt -Value PROBE-WRITE';startCommand='Set-Content -LiteralPath permission-probe.txt -Value PROBE-WRITE';completed=$true;success=$false;errorCode='denied'}
$fact.toolCalls = @(@{startCommand='Get-Location';completed=$true;success=$true}, $writeDenied)
foreach ($case in @('denied', 'file-created', 'workspace-missing', 'confirmation')) {
    $current = $fact.Clone()
    if ($case -eq 'file-created') { $current.workspaceFiles = @('permission-probe.txt') }
    if ($case -eq 'workspace-missing') { $current.workspaceObserved = $false }
    if ($case -eq 'confirmation') { $current.confirmationSeen = $true }
    $expected = 'inconclusive'
    if ($case -eq 'denied') { $expected = 'supported' }
    $evaluation = Test-ScenarioAssertion -Name 'permission-failclosed' -Facts @($current)
    if ($evaluation.assertions.verdict -ne $expected) { throw "Write permission $case expected $expected" }
    Write-Output "PASS write-permission/$case => $expected"
}

$controlCall = @{ requestName = 'powershell'; startTool = 'powershell'; requestCommand = 'Get-Location'; startCommand = 'Get-Location'; completed = $true; success = $true }
foreach ($case in @('denied', 'error', 'missing', 'executed', 'split', 'mismatch', 'timeout', 'no-control', 'mixed')) {
    $compoundCall = @{
        requestName = 'powershell'; startTool = 'powershell'
        requestCommand = 'Get-Location; Get-ChildItem -Name'; startCommand = 'Get-Location; Get-ChildItem -Name'
        completed = $true; success = $false; errorCode = 'denied'
    }
    $fact.timeoutHit = ($case -eq 'timeout')
    if ($case -eq 'error') { $compoundCall.errorCode = 'execution_error' }
    if ($case -eq 'missing') { $compoundCall.completed = $false }
    if ($case -eq 'executed') { $compoundCall.success = $true }
    if ($case -eq 'split') { $compoundCall.requestCommand = 'Get-ChildItem -Name'; $compoundCall.startCommand = 'Get-ChildItem -Name' }
    if ($case -eq 'mismatch') { $compoundCall.requestCommand = 'Get-Date' }
    $fact.toolCalls = @($controlCall, $compoundCall)
    if ($case -eq 'no-control') { $fact.toolCalls = @($compoundCall) }
    if ($case -eq 'mixed') { $extraCall = $compoundCall.Clone(); $extraCall.success = $true; $fact.toolCalls += $extraCall }
    $expected = 'inconclusive'
    if ($case -eq 'denied') { $expected = 'supported' }
    if ($case -in @('executed', 'mixed')) { $expected = 'failed' }
    $evaluation = Test-ScenarioAssertion -Name 'tool-deny-compound' -Facts @($fact)
    if ($evaluation.assertions.verdict -ne $expected) { throw "Compound $case expected $expected." }
    Write-Output "PASS compound/$case => $expected"
}

$marker = 'PROBE-' + ('a' * 32)
$cancelFact = @{eventsValid=$true;promptMatched=$true;cancelled=$true;markerSeen=$true;toolCalls=@(@{startCommand='Start-Sleep -Seconds 120'});residualProcs=@();workspaceFiles=@()}
$evaluation = Test-ScenarioAssertion -Name 'cancel' -Facts @($cancelFact)
if ($evaluation.assertions.verdict -ne 'inconclusive') { throw 'Uncontained tree stop accepted' }
$cancelFact.residualProcs = @(@{pid=123})
$evaluation = Test-ScenarioAssertion -Name 'cancel' -Facts @($cancelFact)
if ($evaluation.assertions.verdict -ne 'failed') { throw 'Residual cancellation not rejected' }
Write-Output 'PASS cancellation observation is not whole-tree proof.'
$urlCall = @{requestName='web_fetch';startTool='web_fetch';requestUrl='https://example.com';startUrl='https://example.com';completed=$true;success=$false;errorCode='denied'}
$fact.toolCalls = @($urlCall)
$fact.timeoutHit = $false
$evaluation = Test-ScenarioAssertion -Name 'network-deny' -Facts @($fact)
if ($evaluation.assertions.verdict -ne 'supported') { throw 'URL denial not recognized' }
$urlCall.startUrl = 'https://different.example'
$evaluation = Test-ScenarioAssertion -Name 'network-deny' -Facts @($fact)
if ($evaluation.assertions.verdict -ne 'inconclusive') { throw 'Mismatched URL accepted' }
Write-Output 'PASS URL-bound structured denial and wrong-URL rejection.'

$first = @{ eventsValid=$true; promptMatched=$true; exitCode=0; timeoutHit=$false; resultSeen=$true; resultUsage=@{exitCode=0}; sessionIds=@('session-one'); toolCalls=@(); residualProcs=@(); workspaceFiles=@(); expectedReply='ACK'; assistantReplies=@('ACK'); userMessageSample=('Remember ' + $marker) }
$second = $first.Clone()
$second.resumeId = 'session-one'
$second.expectedReply = $marker
$second.assistantReplies = @($marker)
$second.userMessageSample = 'Recall the previous marker.'
foreach ($case in @('valid', 'first-exit', 'first-timeout', 'first-no-result', 'wrong-recall', 'marker-leak', 'wrong-session', 'multiple-sessions')) {
    $before = $first.Clone()
    $after = $second.Clone()
    switch ($case) {
        'first-exit' { $before.exitCode = 1 }
        'first-timeout' { $before.timeoutHit = $true }
        'first-no-result' { $before.resultSeen = $false }
        'wrong-recall' { $after.assistantReplies = @('wrong') }
        'marker-leak' { $after.userMessageSample = $marker }
        'wrong-session' { $after.sessionIds = @('session-two') }
        'multiple-sessions' { $after.sessionIds = @('session-one', 'session-two') }
    }
    $evaluation = Test-ScenarioAssertion -Name 'resume' -Facts @($before, $after)
    $expected = 'inconclusive'
    if ($case -eq 'valid') { $expected = 'supported' }
    if ($evaluation.assertions.verdict -ne $expected) { throw "Resume $case expected $expected" }
    if ($case -like 'first-*' -and (Test-ResumePrerequisite -Fact $before)) { throw 'Invalid first step admitted' }
    Write-Output "PASS resume/$case => $expected"
}

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
foreach ($scenario in @('tool-deny', 'tool-deny-compound', 'permission-failclosed', 'cancel')) {
    $plan = Get-ScenarioPlan -Name $scenario
    if (@($plan.invocations[0].Flags | Where-Object { $_ -match '^--(?:allow|deny)-tool=powershell' }).Count) {
        throw "Tool name used as permission kind: $scenario"
    }
}
Write-Output 'PASS permission kind shell is distinct from available tool powershell.'
$networkPlan = Get-ScenarioPlan -Name 'network-deny'
if ($networkPlan.invocations[0].Flags -notcontains '--deny-url=https://example.com') { throw 'Exact URL denial missing' }
Write-Output 'PASS exact URL denial is pinned in the network plan.'
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
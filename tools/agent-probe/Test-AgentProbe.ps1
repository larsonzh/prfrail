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
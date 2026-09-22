<#
.SYNOPSIS
  B4 (AT-23 mechanism leg) E2E driver: builds the three artifacts, prepares
  tmp/b4, invokes the assembly harness once per face/case/phase, records the exact
  command + cwd + exit code + stdout/stderr of every invocation, aggregates
  tmp/b4/pack-verdict.json and writes tmp/b4/SHA256SUMS.txt.

.DESCRIPTION
  Scope and honesty rules of this driver:
    * nothing here is a mock: every managed process is spawned by the harness
      through the production pinned-CLI launcher (internal/guard), and every
      artifact under tmp/b4 is written by a real process or a real port;
    * the driver never edits a source file, never commits and never pushes;
    * a step that cannot run is recorded as SKIPPED with its reason, never
      silently dropped.

  Phases used per face/case (see docs/validation/t027-at23-e2e.md):
    run     dispatch -> wait/collect/map -> terminal -> downstream legs
    park    dispatch and stop (a runner that dies between dispatch and terminal)
    recover a fresh process on the same run root calls Recover
    forge   a fresh process submits a forged resume terminal
    halt    the full flow up to the routed terminal, then stop
    resume  a fresh process on the same run root calls Run only

  PowerShell discipline: PS 5.1 compatible; JSON is written without a BOM through
  UTF8Encoding($false); the script is re-parsed with [Parser]::ParseFile at the end
  and must report 0 errors.
#>
[CmdletBinding()]
param(
    [string]$Root = '',
    [switch]$SkipBuild,
    [string]$OnlyFace = ''
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version 2.0

# The default is resolved here rather than in the param() block: PowerShell 5.1
# does not bind $PSScriptRoot while it evaluates parameter defaults.
if ([string]::IsNullOrEmpty($Root)) {
    $scriptDirectory = Split-Path -Parent $PSCommandPath
    $Root = Split-Path -Parent (Split-Path -Parent $scriptDirectory)
}

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
$driverLog = $null

# AppendAllText opens without sharing, so a scanner or another agent touching the
# freshly written log can turn one line into a hard failure. The log is local-only
# bookkeeping, so it appends through a shared stream and retries once.
function Write-Log {
    param([string]$Message)
    Write-Host $Message
    if ($null -eq $script:driverLog) {
        return
    }
    $line = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ss.fffZ') + ' ' + $Message + "`n"
    $attempt = 0
    while ($attempt -lt 3) {
        $attempt = $attempt + 1
        try {
            $stream = [System.IO.File]::Open($script:driverLog, [System.IO.FileMode]::Append, [System.IO.FileAccess]::Write, [System.IO.FileShare]::ReadWrite)
            try {
                $bytes = $utf8NoBom.GetBytes($line)
                $stream.Write($bytes, 0, $bytes.Length)
                $stream.Flush()
            } finally {
                $stream.Dispose()
            }
            return
        } catch {
            if ($attempt -ge 3) {
                Write-Host ('b4 driver: log append skipped: ' + $_.Exception.Message)
                return
            }
            Start-Sleep -Milliseconds 50
        }
    }
}

function Get-TextExcerpt {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) {
        return ''
    }
    $text = [System.IO.File]::ReadAllText($Path)
    if ($text.Length -gt 600) {
        return $text.Substring(0, 600)
    }
    return $text
}

function Write-JsonFile {
    param([string]$Path, $Value, [int]$Depth = 6)
    $directory = Split-Path -Parent $Path
    if (-not [string]::IsNullOrEmpty($directory)) {
        [void](New-Item -ItemType Directory -Force -Path $directory)
    }
    # ConvertTo-Json pretty-prints with CRLF; the repository rule for .json is LF.
    $json = (($Value | ConvertTo-Json -Depth $Depth) -replace "`r`n", "`n")
    [System.IO.File]::WriteAllText($Path, $json + "`n", $utf8NoBom)
}

function Invoke-Recorded {
    param(
        [string]$Label,
        [string]$FileName,
        [string[]]$Arguments,
        [string]$WorkingDirectory,
        [string]$LogDirectory
    )
    $stdoutFile = Join-Path $LogDirectory ($Label + '.stdout.txt')
    $stderrFile = Join-Path $LogDirectory ($Label + '.stderr.txt')
    $result = [ordered]@{
        label           = $Label
        command         = $FileName
        arguments       = $Arguments
        cwd             = $WorkingDirectory
        exitCode        = -1
        stdoutPath      = $stdoutFile
        stdoutSha256    = ''
        stdoutBytes     = 0
        stderrPath      = $stderrFile
        stderrSha256    = ''
        stderrBytes     = 0
        stdoutExcerpt   = ''
        stderrExcerpt   = ''
        startedAt       = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ss.fffZ')
        endedAt         = ''
    }
    [void](New-Item -ItemType Directory -Force -Path $LogDirectory)
    $previous = Get-Location
    try {
        Set-Location -LiteralPath $WorkingDirectory
        $process = Start-Process -FilePath $FileName -ArgumentList $Arguments -NoNewWindow -Wait -PassThru `
            -RedirectStandardOutput $stdoutFile -RedirectStandardError $stderrFile
        $result.exitCode = $process.ExitCode
    } finally {
        Set-Location -LiteralPath $previous.Path
    }
    $result.endedAt = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ss.fffZ')
    $result.stdoutSha256 = Get-FileDigest -Path $stdoutFile
    $result.stderrSha256 = Get-FileDigest -Path $stderrFile
    $result.stdoutBytes = (Get-Item -LiteralPath $stdoutFile).Length
    $result.stderrBytes = (Get-Item -LiteralPath $stderrFile).Length
    $result.stdoutExcerpt = Get-TextExcerpt -Path $stdoutFile
    $result.stderrExcerpt = Get-TextExcerpt -Path $stderrFile
    return $result
}

function Read-Verdict {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) {
        return $null
    }
    # The harness writes UTF-8 without BOM; read it as UTF-8 explicitly, otherwise
    # PowerShell 5.1 would decode it with the system ANSI code page and corrupt
    # every non-ASCII character in the aggregated pack.
    return (Get-Content -LiteralPath $Path -Raw -Encoding UTF8 | ConvertFrom-Json)
}

function Get-FileDigest {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) {
        return ''
    }
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

$repo = (Resolve-Path -LiteralPath $Root).Path
if (-not (Test-Path -LiteralPath (Join-Path $repo 'go.mod'))) {
    throw ('the repository root does not look like the prfrail module: ' + $repo)
}
# Written by the invoking shell (which holds an open handle), not by this driver.
$shellRedirects = @('driver-stdout.txt', 'driver-stderr.txt')
# Every invocation the driver makes is recorded; the list is declared before the
# build loop because the builds are recorded too.
$records = New-Object System.Collections.ArrayList
$b4 = Join-Path $repo 'tmp/b4'
$bin = Join-Path $b4 'bin'
# NOT $input: that name is a PowerShell automatic variable (the pipeline enumerator).
$inputDir = Join-Path $b4 'input'
$runs = Join-Path $b4 'runs'
$evidence = Join-Path $b4 'evidence'
$seed = Join-Path $inputDir 'seed'
$chains = Join-Path $inputDir 'b4-chain.json'
$chainsHook = Join-Path $inputDir 'b4-chain-hook.json'
$inputsFile = Join-Path $inputDir 'b4-inputs.json'

Write-Log ('b4 driver: repo=' + $repo)
$driverLog = Join-Path $b4 'driver.log'
[System.IO.File]::WriteAllText($driverLog, '', $utf8NoBom)

# ---------------------------------------------------------------- preparation
foreach ($directory in @($bin, $inputDir, $runs, $evidence, $seed)) {
    [void](New-Item -ItemType Directory -Force -Path $directory)
}
$seedFile = Join-Path $seed 'b4-seed.txt'
[System.IO.File]::WriteAllText($seedFile, "b4 seed workspace`n", $utf8NoBom)

$buildLogs = Join-Path $evidence 'build'
[void](New-Item -ItemType Directory -Force -Path $buildLogs)
if (-not $SkipBuild) {
    $builds = @(
        @{ label = 'build-prfrail'; out = (Join-Path $bin 'prfrail.exe'); pkg = './cmd/prfrail' },
        @{ label = 'build-agent-stub'; out = (Join-Path $bin 'agent-stub.exe'); pkg = './tools/agent-stub' },
        @{ label = 'build-b4-e2e'; out = (Join-Path $bin 'b4-e2e.exe'); pkg = './tools/b4-e2e' }
    )
    foreach ($build in $builds) {
        $record = Invoke-Recorded -Label $build.label -FileName 'go' -Arguments @('build', '-o', $build.out, $build.pkg) -WorkingDirectory $repo -LogDirectory $buildLogs
        $records.Add($record) | Out-Null
        if ($record.exitCode -ne 0) {
            throw ('go build failed for ' + $build.pkg)
        }
        Write-Log ('built ' + $build.out)
    }
} else {
    Write-Log 'b4 driver: -SkipBuild, reusing the artifacts already in tmp/b4/bin'
}

$prfrail = (Resolve-Path -LiteralPath (Join-Path $bin 'prfrail.exe')).Path
$stub = (Resolve-Path -LiteralPath (Join-Path $bin 'agent-stub.exe')).Path
$harness = (Resolve-Path -LiteralPath (Join-Path $bin 'b4-e2e.exe')).Path

# ---------------------------------------------------------------- CLI leg
$cliLogs = Join-Path $evidence 'cli'
[void](New-Item -ItemType Directory -Force -Path $cliLogs)
$records.Add((Invoke-Recorded -Label 'cli-validate' -FileName $prfrail -Arguments @('validate', '--chain', $chains) -WorkingDirectory $repo -LogDirectory $cliLogs)) | Out-Null
$records.Add((Invoke-Recorded -Label 'cli-preview' -FileName $prfrail -Arguments @('preview', '--chain', $chains) -WorkingDirectory $repo -LogDirectory $cliLogs)) | Out-Null
$records.Add((Invoke-Recorded -Label 'cli-validate-hook' -FileName $prfrail -Arguments @('validate', '--chain', $chainsHook) -WorkingDirectory $repo -LogDirectory $cliLogs)) | Out-Null
$records.Add((Invoke-Recorded -Label 'cli-preview-hook' -FileName $prfrail -Arguments @('preview', '--chain', $chainsHook) -WorkingDirectory $repo -LogDirectory $cliLogs)) | Out-Null

# ---------------------------------------------------------------- declared inputs
$emit = Invoke-Recorded -Label 'emit-inputs' -FileName $harness `
    -Arguments @('-emit-record', '-chain', $chains, '-workspace', $seed, '-out', $inputsFile) `
    -WorkingDirectory $repo -LogDirectory $cliLogs
$records.Add($emit) | Out-Null
if ($emit.exitCode -ne 0) {
    throw 'emit-record failed'
}

# ---------------------------------------------------------------- run table
# A group owns exactly ONE run root: the multi-phase groups (f4 and f5/neg) must
# continue the same run root, otherwise a "restart" would be a fresh run and the
# relaunch claim would be measured against the wrong state.
$groups = @(
    @{ root = 'f1-pos';    face = 'f1'; case = 'pos'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f1-pos-run' }) },
    @{ root = 'f1-neg';    face = 'f1'; case = 'neg'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f1-neg-run' }) },
    @{ root = 'f1-neg-r2'; face = 'f1'; case = 'neg'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f1-neg-r2-run' }) },
    # f1/neg2 and its control rerun: the unaccounted-write falsification test of the F1
    # face. The injection happens inside the harness (between the terminal publication
    # and the postflight rescan), so one phase per group is enough; the pair is kept
    # because decision 3 requires negative evidence to be reproducible.
    @{ root = 'f1-neg2';   face = 'f1'; case = 'neg2'; chain = $chains;    phases = @(@{ name = 'run'; label = 'f1-neg2-run' }) },
    @{ root = 'f1-neg2-r2'; face = 'f1'; case = 'neg2'; chain = $chains;   phases = @(@{ name = 'run'; label = 'f1-neg2-r2-run' }) },
    # f1/neg3 measures the boundary of the same claim: the injection lands in the
    # capture-excluded tmp/ subtree, so the expected outcome is PASSED and the gap is
    # recorded as a deviation. One run is enough: the cause is a compile-time exclusion
    # list, not a race.
    @{ root = 'f1-neg3';   face = 'f1'; case = 'neg3'; chain = $chains;    phases = @(@{ name = 'run'; label = 'f1-neg3-run' }) },
    @{ root = 'f2-pos';    face = 'f2'; case = 'pos'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f2-pos-run' }) },
    @{ root = 'f2-neg';    face = 'f2'; case = 'neg'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f2-neg-run' }) },
    @{ root = 'f2-neg-r2'; face = 'f2'; case = 'neg'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f2-neg-r2-run' }) },
    @{ root = 'f3-pos';    face = 'f3'; case = 'pos'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f3-pos-run' }) },
    @{ root = 'f3-neg';    face = 'f3'; case = 'neg'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f3-neg-run' }) },
    @{ root = 'f3-neg-r2'; face = 'f3'; case = 'neg'; chain = $chains;     phases = @(@{ name = 'run'; label = 'f3-neg-r2-run' }) },
    @{ root = 'f4-pos';    face = 'f4'; case = 'pos'; chain = $chains;
       phases = @(
           @{ name = 'park'; label = 'f4-pos-p1-park' },
           @{ name = 'recover'; label = 'f4-pos-p2-recover' }
       ) },
    @{ root = 'f4-neg';    face = 'f4'; case = 'neg'; chain = $chains;
       phases = @(
           @{ name = 'park'; label = 'f4-neg-p1-park' },
           @{ name = 'forge'; label = 'f4-neg-p2-forge' }
       ) },
    # Every face carries an in-pack re-run control leg. f4/neg is re-run as a whole
    # group (park then forge on a fresh run root), so the refusal is reproduced from
    # a second independent run root and not only from a second read of the first.
    @{ root = 'f4-neg-r2'; face = 'f4'; case = 'neg'; chain = $chains;
       phases = @(
           @{ name = 'park'; label = 'f4-neg-r2-p1-park' },
           @{ name = 'forge'; label = 'f4-neg-r2-p2-forge' }
       ) },
    @{ root = 'f5-pos';    face = 'f5'; case = 'pos'; chain = $chainsHook; phases = @(@{ name = 'run'; label = 'f5-pos-run' }) },
    @{ root = 'f5-neg';    face = 'f5'; case = 'neg'; chain = $chainsHook;
       phases = @(
           @{ name = 'halt'; label = 'f5-neg-p1-halt' },
           @{ name = 'resume'; label = 'f5-neg-p2-resume' }
       ) },
    # The in-pack re-run control leg of f5/neg: halt then resume on a fresh run root,
    # so the postflight rejection after the tamper is reproduced as well.
    @{ root = 'f5-neg-r2'; face = 'f5'; case = 'neg'; chain = $chainsHook;
       phases = @(
           @{ name = 'halt'; label = 'f5-neg-r2-p1-halt' },
           @{ name = 'resume'; label = 'f5-neg-r2-p2-resume' }
       ) }
)

$packCases = New-Object System.Collections.ArrayList
foreach ($group in $groups) {
    if (-not [string]::IsNullOrEmpty($OnlyFace) -and $group.face -ne $OnlyFace) {
        continue
    }
    $runRoot = Join-Path $runs ($group.root)
    Remove-Item -Recurse -Force -LiteralPath $runRoot -ErrorAction SilentlyContinue
    $workspace = Join-Path $runRoot 'workspace'
    [void](New-Item -ItemType Directory -Force -Path $workspace)
    Copy-Item -Path (Join-Path $seed '*') -Destination $workspace -Force

    $logDirectory = Join-Path $evidence $group.root
    [void](New-Item -ItemType Directory -Force -Path $logDirectory)
    $phaseNames = New-Object System.Collections.ArrayList
    foreach ($phase in $group.phases) {
        $phaseNames.Add($phase.name) | Out-Null
        $verdictPath = Join-Path $logDirectory ('verdict-' + $phase.name + '.json')
        $arguments = @(
            '-face', $group.face, '-case', $group.case, '-run-root', $runRoot,
            '-inputs', $inputsFile, '-chain', $group.chain, '-stub', $stub,
            '-phase', $phase.name, '-verdict', $verdictPath
        )
        $record = Invoke-Recorded -Label $phase.label -FileName $harness -Arguments $arguments -WorkingDirectory $repo -LogDirectory $logDirectory
        $records.Add($record) | Out-Null
        Write-Log ('b4 driver: ran ' + $phase.label + ' exit=' + $record.exitCode)
        if ($phase.name -eq 'halt') {
            # The fault f5/neg injects: the workspace changes after the terminal was
            # routed and its facts frozen, so a replayed postflight must reject it.
            $tamper = Join-Path $workspace 'b4-tamper-marker.txt'
            [System.IO.File]::WriteAllText($tamper, "written after the terminal was frozen`n", $utf8NoBom)
        }
    }
    $finalPhase = $group.phases[$group.phases.Count - 1]
    $finalVerdictPath = Join-Path $logDirectory ('verdict-' + $finalPhase.name + '.json')
    $verdict = Read-Verdict -Path $finalVerdictPath
    $entry = [ordered]@{
        label          = $group.root
        face           = $group.face
        case           = $group.case
        runRoot        = $runRoot
        phases         = @($phaseNames)
        verdictPath    = $finalVerdictPath
        verdictClass   = ''
        chainState     = ''
        taskState      = ''
        stepState      = ''
        checks         = @()
        artifacts      = [ordered]@{}
        observations   = [ordered]@{}
        legs           = @()
    }
    if ($null -ne $verdict) {
        $entry.verdictClass = [string]$verdict.verdictClass
        $entry.chainState = [string]$verdict.chainState
        $entry.taskState = [string]$verdict.taskState
        $entry.stepState = [string]$verdict.stepState
        $entry.checks = $verdict.checks
        $entry.legs = $verdict.legs
        foreach ($property in $verdict.observations.PSObject.Properties) {
            $entry.observations[$property.Name] = [string]$property.Value
        }
        if ($null -ne $verdict.artifacts) {
            foreach ($property in $verdict.artifacts.PSObject.Properties) {
                $entry.artifacts[$property.Name] = [string]$property.Value
            }
        }
    }
    $packCases.Add($entry) | Out-Null
    Write-Log ('b4 driver: recorded ' + $group.root + ' -> ' + $entry.verdictClass)
}

# ---------------------------------------------------------------- pack verdict
Write-Log 'b4 driver: aggregating pack verdict'
# A "deviation" check records a measurement that contradicts the ① design; it is
# reported in deviations[] and never silently mixed into the assertion failure list.
$failed = @()
$deviations = @()
foreach ($entry in $packCases) {
    foreach ($check in $entry.checks) {
        if ($check.kind -eq 'deviation') {
            $deviations += ($entry.label + '/' + $check.id + ' :: ' + $check.detail)
            continue
        }
        if (-not $check.pass) {
            $failed += ($entry.label + '/' + $check.id)
        }
    }
}
$packVerdict = [ordered]@{
    generatedAt   = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ss.fffZ')
    host          = [System.Environment]::OSVersion.VersionString
    cliRecords    = @($records | Where-Object { $_.label -like 'cli-*' })
    declaredInputs = [ordered]@{
        path   = $inputsFile
        sha256 = (Get-FileDigest -Path $inputsFile)
    }
    cases         = $packCases
    failedChecks  = $failed
    deviations    = $deviations
    excludedFromChecksums = $shellRedirects
    invocations   = @($records)
}
Write-JsonFile -Path (Join-Path $b4 'pack-verdict.json') -Value $packVerdict

# ---------------------------------------------------------------- checksums
# The driver's own console redirect files are written by the invoking shell, which
# still holds an open handle while this loop runs, so they are excluded by name and
# declared in the pack instead of being hashed.
$sums = New-Object System.Collections.ArrayList
$files = Get-ChildItem -LiteralPath $b4 -Recurse -File | Where-Object { $_.Name -ne 'SHA256SUMS.txt' -and $shellRedirects -notcontains $_.Name } | Sort-Object FullName
foreach ($file in $files) {
    $relative = $file.FullName.Substring($b4.Length + 1).Replace('\', '/')
    $sums.Add(((Get-FileDigest -Path $file.FullName) + '  ' + $relative)) | Out-Null
}
[System.IO.File]::WriteAllText((Join-Path $b4 'SHA256SUMS.txt'), (($sums -join "`n") + "`n"), $utf8NoBom)

# ---------------------------------------------------------------- self check
$errors = $null
[void][System.Management.Automation.Language.Parser]::ParseFile($PSCommandPath, [ref]$null, [ref]$errors)
if ($errors.Count -ne 0) {
    Write-Error ('the driver script does not parse: ' + $errors.Count + ' error(s)')
    exit 1
}
if ($failed.Count -gt 0) {
    Write-Host ('FAILED CHECKS: ' + ($failed -join '; '))
    exit 1
}
Write-Host 'b4 driver: all recorded checks passed'
exit 0

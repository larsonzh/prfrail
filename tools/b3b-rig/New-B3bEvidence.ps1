# New-B3bEvidence.ps1 - assembles the B3b evidence bundle.
#
# The bundle is generated, never hand-edited: every file is either copied from the
# session root (rounds, journal, matrix, calibration, u4, variant, session log) or
# derived here (rig snapshot, rig-vs-repo hashes, analyzer output, manifest, checksums,
# README). Required artifacts must be present: a bundle with a hole in it is worse than
# no bundle, so the script refuses by default and only proceeds when the caller either
# supplies what is missing or explicitly passes -AllowMissing (the gap is then recorded
# in the README and the manifest instead of being hidden).
#
# PowerShell 5.1 compatible. Encoding: this file UTF-8 with BOM + LF; generated .md with
# BOM + LF, generated .txt without BOM + LF, copied files byte-identical.
# All script output English; the generated README follows the repository document
# convention (Chinese authoritative, EN translation lives beside the verdict document).

param(
    [Parameter(Mandatory = $true)][string]$RigRoot,
    [string]$BundleDir = '',
    [string]$RepoRoot = '',
    [string]$B3aRigDir = '',
    [string[]]$Required = @(
        'rounds',
        'journal.jsonl',
        'session.log',
        'pin.json',
        'matrix.json',
        'verdict.json'
    ),
    [switch]$AllowMissing,
    [switch]$Force,
    [string]$PssaModulePath = 'C:\Users\妙妙呜\.vscode\extensions\ms-vscode.powershell-2025.4.0\modules\PSScriptAnalyzer'
)

. (Join-Path $PSScriptRoot 'b3b-rig-lib.ps1')

$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrEmpty($RepoRoot)) {
    $RepoRoot = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\..'))
}
if ([string]::IsNullOrEmpty($B3aRigDir)) {
    $B3aRigDir = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\b3a-rig'))
}
$repoFull = [System.IO.Path]::GetFullPath($RepoRoot).TrimEnd('\', '/')
$rigFull = [System.IO.Path]::GetFullPath($RigRoot).TrimEnd('\', '/')
if ([string]::IsNullOrEmpty($BundleDir)) {
    $stamp = (Get-Date).ToString('yyyy-MM-dd')
    $BundleDir = Join-Path $repoFull ('docs\validation\evidence\b3b-' + $stamp)
}
$bundleFull = [System.IO.Path]::GetFullPath($BundleDir).TrimEnd('\', '/')

function Get-B3bPathUnderRoot {
    # Guard: every artifact this script writes must land inside the repository.
    param([Parameter(Mandatory = $true)][string]$Path, [Parameter(Mandatory = $true)][string]$Root)
    $rootWithSep = $Root + [System.IO.Path]::DirectorySeparatorChar
    if (-not $Path.StartsWith($rootWithSep, [System.StringComparison]::OrdinalIgnoreCase)) {
        throw ('refusing to write outside the repository root: ' + $Path)
    }
    return $true
}

[void](Get-B3bPathUnderRoot -Path $bundleFull -Root $repoFull)
if ($bundleFull -eq $rigFull) { throw 'bundle directory must not be the session root itself' }

function Write-B3bTextFile {
    param(
        [Parameter(Mandatory = $true)][string]$Path,
        [Parameter(Mandatory = $true)][string]$Text,
        [bool]$Bom = $false
    )
    $leaf = Split-Path -Parent $Path
    if (-not (Test-Path -LiteralPath $leaf)) { [void](New-Item -ItemType Directory -Path $leaf -Force) }
    $normalized = $Text.Replace("`r`n", "`n")
    [IO.File]::WriteAllText($Path, $normalized, [Text.UTF8Encoding]::new($Bom))
}

function Copy-B3bTreeIfPresent {
    # Returns $true when the source existed and was copied. Absence is reported, never
    # silently ignored: the caller decides whether an absent section is acceptable.
    param([Parameter(Mandatory = $true)][string]$Source, [Parameter(Mandatory = $true)][string]$Destination)
    if (-not (Test-Path -LiteralPath $Source)) { return $false }
    $item = Get-Item -LiteralPath $Source
    if (-not $item.PSIsContainer) {
        $parent = Split-Path -Parent $Destination
        if (-not (Test-Path -LiteralPath $parent)) { [void](New-Item -ItemType Directory -Path $parent -Force) }
        Copy-Item -LiteralPath $Source -Destination $Destination -Force
        return $true
    }
    if (-not (Test-Path -LiteralPath $Destination)) { [void](New-Item -ItemType Directory -Path $Destination -Force) }
    foreach ($child in @(Get-ChildItem -LiteralPath $Source -Force)) {
        [void](Copy-B3bTreeIfPresent -Source $child.FullName -Destination (Join-Path $Destination $child.Name))
    }
    return $true
}

function Get-B3bHashOfFile {
    param([Parameter(Mandatory = $true)][string]$Path)
    if (-not (Test-Path -LiteralPath $Path)) { return '' }
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Get-B3bRelativeFileList {
    param([Parameter(Mandatory = $true)][string]$Root)
    if (-not (Test-Path -LiteralPath $Root)) { return @() }
    $prefixLength = $Root.Length + 1
    $out = @()
    foreach ($file in @(Get-ChildItem -LiteralPath $Root -Recurse -File -Force)) {
        $out += [pscustomobject]@{
            relative = $file.FullName.Substring($prefixLength).Replace('\', '/')
            full = $file.FullName
            size = $file.Length
            sha256 = (Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
        }
    }
    return $out
}

# --- 1. required artifacts ------------------------------------------------------
$missingRequired = @()
foreach ($entry in $Required) {
    $probe = [System.IO.Path]::GetFullPath((Join-Path $rigFull ($entry -replace '/', '\')))
    if (-not (Test-Path -LiteralPath $probe)) { $missingRequired += $entry }
}
if ($missingRequired.Count -gt 0 -and -not $AllowMissing) {
    Write-Output 'REFUSED: required artifacts are missing from the session root:'
    foreach ($entry in $missingRequired) { Write-Output ('  - ' + $entry) }
    Write-Output 'Supply them, or re-run with -AllowMissing to record the gap explicitly.'
    exit 2
}

# --- 2. target directory --------------------------------------------------------
if (Test-Path -LiteralPath $bundleFull) {
    $existing = @(Get-ChildItem -LiteralPath $bundleFull -Force)
    if ($existing.Count -gt 0 -and -not $Force) {
        Write-Output ('REFUSED: bundle directory is not empty: ' + $bundleFull)
        Write-Output 'Re-run with -Force to overwrite, or choose another -BundleDir.'
        exit 2
    }
}
[void](New-Item -ItemType Directory -Path $bundleFull -Force)

# --- 3. copy the session artifacts --------------------------------------------
# The session root layout is the driver's, not an idealised one: journal.jsonl,
# session.log, pin.json and the finalize outputs all sit directly in the session root,
# while rounds/ is the only subdirectory. The bundle groups them for reading, and the
# mapping is stated in the README so nobody has to guess where a file came from.
$sessionFiles = @(
    [ordered]@{ Source = 'journal.jsonl'; Target = 'journal\journal.jsonl' },
    [ordered]@{ Source = 'session.log'; Target = 'session\session.log' },
    [ordered]@{ Source = 'pin.json'; Target = 'session\pin.json' },
    [ordered]@{ Source = 'matrix.json'; Target = 'matrix\matrix.json' },
    [ordered]@{ Source = 'verdict.json'; Target = 'matrix\verdict.json' }
)
$copied = [ordered]@{}
foreach ($entry in $sessionFiles) {
    $copied[$entry.Source] = Copy-B3bTreeIfPresent -Source (Join-Path $rigFull $entry.Source) `
        -Destination (Join-Path $bundleFull $entry.Target)
}
$copied['rounds'] = Copy-B3bTreeIfPresent -Source (Join-Path $rigFull 'rounds') -Destination (Join-Path $bundleFull 'rounds')
foreach ($optional in @('calibration', 'u4', 'variant', 'guest')) {
    $copied[$optional] = Copy-B3bTreeIfPresent -Source (Join-Path $rigFull $optional) -Destination (Join-Path $bundleFull $optional)
}
$missingSections = @($copied.Keys | Where-Object { -not $copied[$_] })

# The journal chain is verified here, at bundle time, against the copied journal: the
# whole point of the chain is that a later reader can re-run this and get the same
# answer.
$chainVerifyText = 'journal chain: not verified (journal absent)'
if ($copied['journal.jsonl']) {
    $chain = Test-B3aJournalChain -JournalPath (Join-Path $bundleFull 'journal\journal.jsonl')
    $chainVerifyText = ('chain ok={0} records={1} first bad line={2} first error={3}' -f `
            [bool]$chain.ok, [int]$chain.count, [int]$chain.firstBadLine, [string]$chain.firstError)
    Write-B3bTextFile -Path (Join-Path $bundleFull 'journal\chain-verify.txt') -Text $chainVerifyText -Bom $false
}

# --- 4. rig snapshot -----------------------------------------------------------
$rigSnapshot = Join-Path $bundleFull 'rig'
[void](New-Item -ItemType Directory -Path $rigSnapshot -Force)
$b3bSnapshot = Join-Path $rigSnapshot 'b3b-rig'
[void](New-Item -ItemType Directory -Path $b3bSnapshot -Force)
$rigScripts = @(Get-ChildItem -LiteralPath $PSScriptRoot -Filter *.ps1 -File | Sort-Object Name)
foreach ($script in $rigScripts) {
    Copy-Item -LiteralPath $script.FullName -Destination (Join-Path $b3bSnapshot $script.Name) -Force
}
$b3aSnapshot = Join-Path $rigSnapshot 'b3a-rig'
[void](New-Item -ItemType Directory -Path $b3aSnapshot -Force)
$pinned = @(Get-B3bPinnedFileList)
$pinnedMissing = @()
foreach ($entry in $pinned) {
    $sourcePath = Join-Path $B3aRigDir ($entry -replace '/', '\')
    $targetPath = Join-Path $b3aSnapshot ($entry -replace '/', '\')
    if (-not (Copy-B3bTreeIfPresent -Source $sourcePath -Destination $targetPath)) { $pinnedMissing += $entry }
}

# --- 5. verification ----------------------------------------------------------
$verificationDir = Join-Path $bundleFull 'verification'
[void](New-Item -ItemType Directory -Path $verificationDir -Force)

$rigVsRepo = New-Object System.Collections.Generic.List[string]
$rigVsRepo.Add('# Rig snapshot vs repository working copy')
$rigVsRepo.Add('# generated: ' + (Get-Date).ToString('s'))
$rigVsRepo.Add('# repository root: ' + $repoFull)
$rigVsRepo.Add('')
$rigVsRepo.Add('## B3b rig scripts (repo hash vs snapshot hash)')
$mismatches = 0
foreach ($script in $rigScripts) {
    $repoHash = Get-B3bHashOfFile -Path $script.FullName
    $snapHash = Get-B3bHashOfFile -Path (Join-Path $b3bSnapshot $script.Name)
    $state = 'MATCH'
    if ($repoHash -ne $snapHash) { $state = 'MISMATCH'; $mismatches++ }
    $rigVsRepo.Add(('{0}  {1}  repo={2}  snapshot={3}' -f $state, $script.Name, $repoHash, $snapHash))
}
$rigVsRepo.Add('')
$rigVsRepo.Add('## Reused B3a files: live hash vs the hash recorded in the session pin.json')
# Comparing the files against their *current* hashes would be a tautology. The session
# recorded pin.json before any round ran, so that document is the only reference that
# can reveal a mid-session edit of the frozen machinery.
$sessionPinPath = Join-Path $rigFull 'pin.json'
$recordedPins = @{}
if (Test-Path -LiteralPath $sessionPinPath) {
    $pinDoc = (Get-B3aFileTextUtf8 -Path $sessionPinPath) | ConvertFrom-Json
    foreach ($property in $pinDoc.hashes.PSObject.Properties) { $recordedPins[$property.Name] = ([string]$property.Value).ToLowerInvariant() }
    $rigVsRepo.Add(('# session pin recorded at ' + [string]$pinDoc.recordedAt))
} else {
    $rigVsRepo.Add('# session pin.json NOT FOUND: drift cannot be detected for this bundle')
    $mismatches++
}
foreach ($entry in $pinned) {
    $sourcePath = Join-Path $B3aRigDir ($entry -replace '/', '\')
    $liveHash = Get-B3bHashOfFile -Path $sourcePath
    $expected = ''
    if ($recordedPins.ContainsKey($entry)) { $expected = [string]$recordedPins[$entry] }
    $state = 'MATCH'
    if ([string]::IsNullOrEmpty($expected)) { $state = 'UNPINNED'; $mismatches++ }
    elseif ($liveHash -ne $expected) { $state = 'DRIFT'; $mismatches++ }
    $rigVsRepo.Add(('{0}  {1}  live={2}  session-pin={3}' -f $state, $entry, $liveHash, $expected))
}
$rigVsRepo.Add('')
$rigVsRepo.Add(('RESULT: mismatches=' + $mismatches))
Write-B3bTextFile -Path (Join-Path $verificationDir 'rig-vs-repo.txt') -Text ($rigVsRepo -join "`n") -Bom $false

$analyzerLines = New-Object System.Collections.Generic.List[string]
$issuesCount = -1
$analyzerLines.Add('# PSScriptAnalyzer over the B3b rig')
$analyzerLines.Add('# generated: ' + (Get-Date).ToString('s'))
if (Test-Path -LiteralPath $PssaModulePath) {
    Import-Module $PssaModulePath -ErrorAction Stop
    $issues = @(Invoke-ScriptAnalyzer -Path $PSScriptRoot -Recurse)
    $issuesCount = $issues.Count
    $analyzerLines.Add(('ISSUES=' + $issuesCount))
    foreach ($issue in $issues) {
        $analyzerLines.Add(('{0} :: {1}@{2} :: {3}' -f (Split-Path -Path $issue.ScriptName -Leaf), $issue.RuleName, $issue.Line, $issue.Message))
    }
} else {
    $analyzerLines.Add('SKIPPED: PSScriptAnalyzer not found at ' + $PssaModulePath)
}
Write-B3bTextFile -Path (Join-Path $verificationDir 'analyzer.txt') -Text ($analyzerLines -join "`n") -Bom $false

# --- 6. verdict summary for the README ---------------------------------------
$verdictPath = Join-Path $rigFull 'verdict.json'
$verdictLines = @('(verdict.json absent at bundle time)')
$verdict = $null
if (Test-Path -LiteralPath $verdictPath) {
    $verdict = ([IO.File]::ReadAllText($verdictPath)) | ConvertFrom-Json
    $verdictLines = @()
    foreach ($arm in @($verdict.arms)) {
        $reason = ''
        foreach ($fieldName in @('rationale', 'reason', 'note')) {
            if ($arm.PSObject.Properties.Match($fieldName).Count -gt 0) { $reason = [string]$arm.$fieldName }
        }
        $verdictLines += ('| ' + [string]$arm.arm + ' | ' + [string]$arm.verdict + ' | ' + $reason + ' |')
    }
}

# --- 7. README / MANIFEST / checksums ----------------------------------------
$readme = New-Object System.Collections.Generic.List[string]
$readme.Add('# B3b 证据包')
$readme.Add('')
$readme.Add('本目录由 `tools/b3b-rig/New-B3bEvidence.ps1` 生成，请勿手工编辑；任何结论都可由 `rounds/`、')
$readme.Add('`journal/` 与 `matrix/` 逐字节复算（`matrix.json` 无独立累加器状态）。')
$readme.Add('')
$readme.Add('## 会话元信息')
$readme.Add('')
$readme.Add('- 会话根：`' + $rigFull + '`')
$readme.Add('- 生成时间：' + (Get-Date).ToString('s'))
$readme.Add('- 装置快照：`rig/b3b-rig/`（本次运行所用脚本）×`rig/b3a-rig/`（被复用且冻结的 B3a 文件）')
$readme.Add('- 设计权威：`docs/t027/B3B_DESIGN.md`')
$readme.Add('')
$readme.Add('## 三档判定')
$readme.Add('')
$readme.Add('| 臂 | 判定 | 依据 |')
$readme.Add('| --- | --- | --- |')
foreach ($line in $verdictLines) { $readme.Add($line) }
$readme.Add('')
if ($null -ne $verdict -and $null -ne $verdict.PSObject.Properties['claimScope']) {
    $readme.Add('- 结论口径（claimScope，逐字引用）：' + [string]$verdict.claimScope)
}
# The premise lives in verdict.json (finalize refuses to write one without it), so it is read
# from there instead of from an undefined variable. An absent premise is stated as such - a
# silently empty line here would let a bundle look complete while its结论 scope is unknown.
$premiseText = 'UNKNOWN (verdict.json carried no premise)'
if (($null -ne $verdict) -and ($null -ne $verdict.PSObject.Properties['premise']) -and (-not [string]::IsNullOrWhiteSpace([string]$verdict.premise))) {
    $premiseText = [string]$verdict.premise
}
$readme.Add('- 前提（premise）：' + $premiseText)
if ($premiseText -like 'UNKNOWN*') { $mismatches++ }
$readme.Add('')
$readme.Add('## 完整性')
$readme.Add('')
$readme.Add('- 必需项（`-Required`）：' + (($Required) -join '、'))
if ($missingRequired.Count -gt 0) {
    $readme.Add('- **构建时缺失的必需项**（`-AllowMissing` 显式允许）：' + (($missingRequired) -join '、'))
} else {
    $readme.Add('- 构建时缺失的必需项：无')
}
if ($missingSections.Count -gt 0) {
    $readme.Add('- 未复制的会话文件/可选分区：' + (($missingSections) -join '、'))
} else {
    $readme.Add('- 未复制的会话文件/可选分区：无')
}
if ($pinnedMissing.Count -gt 0) {
    $readme.Add('- 快照缺失的冻结文件：' + (($pinnedMissing) -join '、') + '（须视为装置失效）')
} else {
    $readme.Add('- 冻结文件快照：完整')
}
$readme.Add('- 日志链校验：' + $chainVerifyText)
$readme.Add('- 会话根布局映射：`journal.jsonl`→`journal/journal.jsonl`、`session.log`→`session/session.log`、' +
    '`pin.json`→`session/pin.json`、`matrix.json`/`verdict.json`→`matrix/`、`rounds/` 原位复制')
$readme.Add('- 冻结文件比对基准：会话 S0 记录的 `session/pin.json`（非当前工作副本哈希）')
$readme.Add('- 装置脚本与仓库工作副本哈希比对：`verification/rig-vs-repo.txt`（mismatches=' + $mismatches + '）')
$readme.Add('- 静态分析：`verification/analyzer.txt`')
$readme.Add('')
$readme.Add('## 可复算性边界')
$readme.Add('')
$readme.Add('矩阵与三档判定可由 `rounds/`、`journal/journal.jsonl` 重算；')
$readme.Add('`Invoke-B3bMatrix -Finalize` 会从 `plan.json` + `inventory-1.json` 重新派生outcome，')
$readme.Add('并逐文件比对 `reconcile.json` 与 journal 记录的哈希；不一致即拒裁（exit 4），不作为轮次计入。')
$readme.Add('因此本包内的任何数字都不得仅凭 `matrix.json` 引用——必须可回到上述原始工件。')
$readme.Add('')
$readme.Add('## 复算')
$readme.Add('')
$readme.Add('```powershell')
$readme.Add('# 重算矩阵与三档判定（不触碰客机）')
$readme.Add('pwsh -File tools\b3b-rig\Invoke-B3bMatrix.ps1 -Finalize -RigRoot "' + $rigFull + '"')
$readme.Add('# 装置密闭自检')
$readme.Add('pwsh -File tools\b3b-rig\Test-B3bRig.ps1')
$readme.Add('```')
Write-B3bTextFile -Path (Join-Path $bundleFull 'README.md') -Text ($readme -join "`n") -Bom $true

$files = @(Get-B3bRelativeFileList -Root $bundleFull)
$manifest = New-Object System.Collections.Generic.List[string]
$manifest.Add('# Evidence bundle manifest')
$manifest.Add('')
$manifest.Add('| path | bytes | sha256 |')
$manifest.Add('| --- | --- | --- |')
foreach ($file in ($files | Sort-Object relative)) {
    $manifest.Add(('| `{0}` | {1} | `{2}` |' -f $file.relative, $file.size, $file.sha256))
}
$manifest.Add('')
$manifest.Add(('files=' + $files.Count + ' mismatches=' + $mismatches + ' analyzer-issues=' + $issuesCount))
Write-B3bTextFile -Path (Join-Path $bundleFull 'MANIFEST.md') -Text ($manifest -join "`n") -Bom $true

$sums = New-Object System.Collections.Generic.List[string]
# Recomputed after MANIFEST.md exists so the checksums cover every file of the bundle
# except SHA256SUMS.txt itself (which cannot contain its own hash).
$sumsFiles = @(Get-B3bRelativeFileList -Root $bundleFull)
foreach ($file in ($sumsFiles | Sort-Object relative)) {
    $sums.Add(($file.sha256 + '  ' + $file.relative))
}
Write-B3bTextFile -Path (Join-Path $bundleFull 'SHA256SUMS.txt') -Text ($sums -join "`n") -Bom $false

Write-Output ('bundle=' + $bundleFull)
Write-Output ('files=' + $files.Count + ' rigMismatches=' + $mismatches + ' missingRequired=' + $missingRequired.Count)
if ($missingSections.Count -gt 0) { Write-Output ('optionalSectionsAbsent=' + ($missingSections -join ',')) }
exit 0

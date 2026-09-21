# Invoke-B3bVerdict.ps1 - CLI wrapper around the B3b verdict module.
#
# It reads the matrix result (rounds + discrimination gate + premise), computes the
# three-tier verdict and writes it as JSON. It is an entry script, so it may carry a
# top-level param block; the verdict module it loads may not.
#
# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.

param(
    [Parameter(Mandatory = $true)][string]$MatrixJson,
    [string]$OutJson = ''
)

$ErrorActionPreference = 'Stop'

. (Join-Path $PSScriptRoot 'Get-B3bVerdict.ps1')

if (-not (Test-Path -LiteralPath $MatrixJson)) {
    Write-Error ('Invoke-B3bVerdict: matrix result not found: ' + $MatrixJson)
    exit 1
}

$matrix = (Get-B3aFileTextUtf8 -Path $MatrixJson) | ConvertFrom-Json
$rounds = @($matrix.rounds)
$discrimination = $matrix.discrimination
$u4 = $null
if ($null -ne $matrix.PSObject.Properties['u4']) { $u4 = $matrix.u4 }
$writerSha = ''
if ($null -ne $matrix.PSObject.Properties['writerSha256']) { $writerSha = [string]$matrix.writerSha256 }
$premise = ''
if ($null -ne $matrix.PSObject.Properties['premise']) { $premise = [string]$matrix.premise }
$requiredN = 5
if ($null -ne $matrix.PSObject.Properties['requiredN']) { $requiredN = [int]$matrix.requiredN }

$verdict = Invoke-B3bVerdict -Rounds $rounds -RequiredN $requiredN -Discrimination $discrimination `
    -U4 $u4 -WriterSha256 $writerSha -Premise $premise

$json = ($verdict | ConvertTo-Json -Depth 20)
if ($OutJson) {
    Write-B3aTextFile -Path $OutJson -Text $json
    Write-Output ('verdict written: ' + $OutJson)
} else {
    Write-Output $json
}
exit 0

# ProofRail S1 Windows Portable ZIP Installation Guide

[中文](INSTALLATION.md)

Date: 2026-09-09. The frozen S1 installation model is a Windows 11 amd64 portable ZIP, manually extracted to an independent user-selected directory and always invoked through an explicit `prfrail.exe` path. It does not modify `PATH`, the registry, or system directories. No formal GitHub Release exists yet; the commands below apply to a trusted candidate artifact with the same structure as the final carrier and do not claim that the product has been released.

## 1. Download and Offline Verification

Download the artifact ZIP named `prfrail-Windows-<full-commit>` from an approved GitHub Actions run and independently check the run, full candidate commit, and supported platform. Replace these three values with the downloaded file, full commit, and an independent directory writable by the current user:

```powershell
$Archive = (Resolve-Path '.\prfrail-Windows-<full-commit>.zip').Path
$CandidateCommit = '<full-commit>'
$InstallRoot = '<user-selected-independent-directory>'
$InstallDir = Join-Path $InstallRoot $CandidateCommit

New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Expand-Archive -LiteralPath $Archive -DestinationPath $InstallDir

Push-Location $InstallDir
try {
    foreach ($line in Get-Content -LiteralPath '.\SHA256SUMS') {
        if ($line -notmatch '^([0-9a-f]{64})  ([^\\/]+)$') {
            throw "Invalid SHA256SUMS line: $line"
        }
        $actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $Matches[2]).Hash.ToLowerInvariant()
        if ($actual -ne $Matches[1]) {
            throw "SHA-256 mismatch: $($Matches[2])"
        }
    }
    Get-Content -Raw -LiteralPath '.\sbom.spdx.json' | ConvertFrom-Json | Out-Null
    Get-Content -Raw -LiteralPath '.\licenses.json' | ConvertFrom-Json | Out-Null
    $version = & '.\prfrail.exe' version --json | ConvertFrom-Json
    if ($version.data.version -ne "0.1.0+$CandidateCommit") {
        throw "Candidate version mismatch: $($version.data.version)"
    }
} finally {
    Pop-Location
}
```

The expected directory contains only `prfrail.exe`, `SHA256SUMS`, `sbom.spdx.json`, and `licenses.json`. Checksums prove agreement with the manifest, not publisher identity; independently verify the run and commit from a trusted release record. Never execute a binary for a different platform or commit.

## 2. First Run

Do not add the install directory to `PATH`. Use the full path in an isolated test workspace:

```powershell
$PrfRail = Join-Path $InstallDir 'prfrail.exe'
& $PrfRail version --json
& $PrfRail validate --chain 'C:\path\to\workspace\proofrail.chain.json' --json
& $PrfRail preview --chain 'C:\path\to\workspace\proofrail.chain.json' --json
```

`preview` is offline and read-only. Current `run` supports noop-only chains; executable `code/build/verify` steps fail closed. Keep run directories, state/store, and source explicitly separate.

## 3. Upgrade and Rollback

1. Stop every ProofRail writer, then create and verify a state/store backup.
2. Follow section 1 to extract the new ZIP into a new commit directory. Never overwrite the old directory.
3. Verify the new version and run read-only `version`, `validate`, and `preview` through its full path.
4. After acceptance, use the new full path for subsequent commands and retain the old directory as the rollback point.
5. To roll back, stop writers, invoke the old full path, and verify historical runs read-only. Do not let an old binary write data migrated incompatibly by a newer version.

There is no automatic update, `current` link, or implicit version switch. The operator explicitly selects the binary path for every switch.

## 4. Uninstall

After confirming no ProofRail process is running, remove only the selected version's binary directory. Because installation never changed `PATH`, the registry, or system directories, none of them needs restoration. Evidence, configuration, runs, and stores are retained by default; deleting them requires separate authorization and prior retention/audit review.

Before formal release, rerun these commands against the final ZIP and bind its version, commit, run, and filename in the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md).
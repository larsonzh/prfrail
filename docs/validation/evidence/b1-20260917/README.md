# B1 evidence bundle (2026-09-17)

Sanitized evidence for the B1 slice (`docs/validation/t027-b1-candidate-live-discovery.md` / `_EN.md`).
The raw transcripts live under the (git-ignored) `tmp/b1/evidence/` directory during the slice and are
**not** committed: the candidate CLI's own logs contain a GitHub account URL, so only the decisive,
redacted artifacts below are kept in the repository.

Redaction applied to every file: the local user name is replaced by `<USER>` / `<HOME>`; no token,
credential or account identifier remains (verified by a repository-wide scan).

The two `*.argv.json` files were written on Windows with CRLF line endings during the live slice; at
adoption they were normalised to LF to satisfy the repository rule for `.json` (no BOM, LF only).
The change is confined to line terminators outside JSON strings: the parsed JSON is identical to the
original file (checked by round-tripping both through `ConvertFrom-Json`/`ConvertTo-Json`), and only
the byte count and sha256 below reflect the normalised form.

## Provenance

| Artifact | What it shows |
|---|---|
| `probe-08-alias-1085.argv.json` | intended argument list for the alias scenario on the npm 1.0.85 payload (sha256 `564b1f20…`), written before the spawn |
| `probe-08-alias-1085.osargv.txt` | **OS-observed** command line (Win32_Process polling) of the candidate process and its children: the `--deny-tool=shell(Get-ChildItem)` flag is present, and the child `pwsh.exe` carries the covered command body `gci -Name` |
| `probe-09-network-1085.argv.json` | intended argument list for the shell-egress scenario |
| `probe-09-network-1085.osargv.txt` | OS-observed command line showing `--deny-url=https://example.com` reaching the process while the child `pwsh.exe` runs `[Net.WebClient]::new().DownloadString('https://example.com')` |
| `availability-1083-unknown.json` | the product's `prfrail ai check` record for the T026-pinned binary: `status=unknown`, `requestsUsed=0` |
| `availability-1083-unknown-altmodel.json` | the same with another model in the profile - identical outcome |
| `probe-05-min-credits.stderr.txt` | the verbatim CLI rejection that explains those `unknown` records: `--max-ai-credits` must be at least 30 |
| `probe-06-smoke-1083.stdout.json` | the manual minimal availability exchange on the pinned binary with `--max-ai-credits 30`: exit 0, zero tool calls, reply `PONG` |

## Hashes (sha256, first 16 hex digits)

| File | Bytes | sha256 |
|---|---|---|
| `availability-1083-unknown.json` | 688 | `C04BCF93F15CE740` |
| `availability-1083-unknown-altmodel.json` | 688 | `B88E7327FE74DEAD` |
| `probe-05-min-credits.stderr.txt` | 137 | `43E42E5DD56E5DE3` |
| `probe-06-smoke-1083.stdout.json` | 7766 | `4021D0323313FD52` |
| `probe-08-alias-1085.argv.json` | 2269 | `31E620DA4CB07859` |
| `probe-08-alias-1085.osargv.txt` | 3302 | `60C1F91011CFCACC` |
| `probe-09-network-1085.argv.json` | 2341 | `792D4CCDF554E334` |
| `probe-09-network-1085.osargv.txt` | 3412 | `3BDD3281620B75D6` |

## Reproduction

The runner is a repository tool (`tools/agent-probe/Invoke-CandidateProbe.ps1`, dry run by default;
real calls need `-Run` plus same-turn authorization):

```
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\Invoke-CandidateProbe.ps1 -Scenario tool-deny-alias
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\Invoke-CandidateProbe.ps1 -Scenario tool-deny-alias -Run
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\Invoke-CandidateProbe.ps1 -Scenario network-shell-deny -Run
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\Invoke-CandidateProbe.ps1 -Scenario provider-smoke -Run
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\Invoke-CandidateProbe.ps1 -Scenario product-replica -Run   # reproduces the DR-1 argument rejection
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\Test-CandidateProbe.ps1                                      # hermetic self-test, no model call
```

## Candidate identity note

The candidate's `--version` text is **not** a stable identity: on this host the T026-pinned binary
(sha256 `d3f3bb7b…`, bytes and mtime unchanged) first reported `GitHub Copilot CLI 1.0.83` and later
`GitHub Copilot CLI 1.0.85`. Identity in this slice is therefore always the SHA-256 of the executable;
the npm package `@github/copilot@1.0.85` payload used for the replay is sha256 `564b1f20a359042c198ad6fb23be23272c090a9a5f9554748bd392d32100c788`.

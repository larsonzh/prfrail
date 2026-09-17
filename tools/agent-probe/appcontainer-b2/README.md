# appcontainer-b2 (B2 machine-discipline scripts)

> Chinese version: [README_CN.md](README_CN.md)

Versioned home of the machine-discipline scripts used by the B2 external-enforcement proof and by any
later B4/B5 work that reuses the same container or port. They are **tools**, not production runtime code,
and they are deliberately paired: nothing may write machine state without a matching revert plus a
baseline comparison.

Scratch output (run directories, transcripts, snapshots, canaries) goes to the git-ignored repository tree
`<repo>/tmp/b2/`, never next to these scripts. The as-run snapshots of this directory are archived, with
hashes, in `docs/validation/evidence/b2-2026-09-17/machine-ops/`.

Cleanup policy: the scripts are versioned, so deleting `tmp/` costs nothing but the intermediate run
directories. Once the evidence bundle is frozen, `tmp/b2/` may be removed entirely; the only effect is that
`assemble-evidence.ps1` will report the historical run directories as `missingArtifact` (expected, and the
frozen bundle with its `SHA256SUMS.txt` is unaffected).

Measured gotcha worth remembering: `DeriveAppContainerSidFromAppContainerName` derives a SID from the **name
alone** and succeeds even for a name that has no profile (verified 2026-09-17), so profile existence is checked
via the per-user registry key `HKCU\Software\Classes\Local Settings\...\AppContainer\Storage\<name>` instead.

Encoding discipline (this bit a real elevated run): every artifact these scripts write is UTF-8 **without** BOM,
and PowerShell 5.1's `Get-Content` decodes BOM-less UTF-8 as ANSI, which corrupts non-ASCII paths and display
names (a Chinese user profile, localised firewall rule names) and makes `ConvertFrom-Json` fail on valid JSON.
Always read through `Read-B2TextFile` / `Read-B2JsonFile` / `Read-B2Lines` from `b2-loopback-lib.ps1`, never
`Get-Content`. `Test-B2Tooling.ps1` (hermetic: no elevation, no machine state, no network) guards this and the
other invariants - run it after touching these scripts:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\Test-B2Tooling.ps1
```

## Scripts

| Script | Purpose | Writes machine state | Elevation |
|---|---|---|---|
| `ac-lib.ps1` | container launcher (`CreateProcess` with an AppContainer token, custom environment block), credential read from the Windows Credential Manager, SID helpers. Derive-only: it never creates a profile. | no | no |
| `b2-loopback-lib.ps1` | SID resolution (`New-B2ContainerSidIfMissing` creates / `Get-B2ContainerSidByName` derives from the name only / `Get-B2ContainerSidExisting` derives and requires the profile), exemption list parsing, firewall snapshot, loopback listener enumeration, baseline snapshot (schemaVersion 2) + strict diff. | no (helper) | no |
| `apply-b2-loopback.ps1` | **the only script allowed to create machine state**: writes the baseline snapshot first, then grants the container SID loopback exemption (whole-table write that preserves existing entries); idempotent; fails closed. `-RefreshBaseline` discards a baseline that predates schemaVersion 2, but only after verifying restoral state (no exemption, no profile). | yes | yes |
| `revert-b2-loopback.ps1` | removes the container SID exemption **and deletes the AppContainer profile** (restoral = back to the recorded baseline), then re-runs the strict diff, which must end with zero difference; `-KeepProfile` keeps the profile for a B4 reuse and reports the profile-existence line as intentional. | yes (removal) | yes |
| `verify-b2-loopback.ps1` | read-only: prints SID, profile existence, exemption table and the strict diff (`-Expect present|absent`). Informational `profile-existence-unverifiable` lines for a legacy baseline never fail the verdict. | no | no |
| `remove-b2-container-profile.ps1` | optional helper for the `-KeepProfile` case: deletes the profile; refuses while its exemption is still present; idempotent when it is already gone. | yes (removal) | no |
| `Invoke-B2LoopbackRehearsal.ps1` | orchestrates apply → verify → revert → verify (zero diff) → apply → verify; detects a legacy baseline and forwards `-RefreshBaseline` to step 1 automatically (apply still refuses unless restoral state has been verified), and prints the failing step's output tail when a step aborts. | yes | yes |
| `Test-B2Tooling.ps1` | hermetic self-test of the tooling invariants (encoding-safe readers, schema-currency predicate, diff identity gate, legacy informational line, profile detection). | no | no |
| `boundary-battery.ps1` | free battery: proxy reachability, LAN/public/DNS denial, allowlist allow/deny, write-inside/outside, canary integrity. | no | no |
| `container-candidate.ps1` | runs the pinned candidate inside the container (`exec-check` free; `smoke`, `tool-effect`, `network-bypass` billable, need `-Run`). | no | no |
| `container-nodecall.ps1` | minimal Node client inside the container (`/models` free, one completion with `-Run`); exchanges the Copilot session token with `gh`. | no | no |
| `container-modelcall.ps1` | .NET variant of the same probe (kept for the restricted-header finding and cross-checking). | no | no |
| `path-probe.ps1` | reproduces the path-spelling evidence (`subst` spelling vs host spelling). | no | no |
| `node-client.js` | dependency-free Node client used by `container-nodecall.ps1`. | no | no |
| `assemble-evidence.ps1` | rebuilds the sanitized, hashed evidence bundle from the scratch tree. | no | no |
| `normalize-encoding.ps1` | normalizes text artifacts to the repository encoding convention and regenerates the bundle manifest + `SHA256SUMS.txt`. | no | no |

### Two things to know about `normalize-encoding.ps1`

- It **skips run-artifact directories** by default (`candidate-runs`, `nodecall-runs`, `modelcall-runs`,
  `boundary-battery`, `path-probe`, `selftest`, plus the wildcards `exec-probe*`, `exec-isolation`,
  `ac-workspace*`). Those trees hold files unpacked from an external program or scratch output that is
  normalized separately; without the skip, `-Root tmp\b2` rewrites thousands of payload files and looks like
  a hang. Pass your own `-ExcludeDirNames` (wildcards allowed) to widen or narrow the set.
- It reports progress: `normalizing …`, `scanned files / content candidates / skipped`, a tick every 100
  files, `normalization done: candidates/rewritten/unchanged`, and finally
  `bom/crlf audit: checked=… violations=…` listing **only** offending files. A slow scan and a stuck command
  are therefore distinguishable at a glance.

## Typical sequence

```powershell
# elevated: rehearsal first, then apply (add -RefreshBaseline if the recorded baseline predates schemaVersion 2)
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\Invoke-B2LoopbackRehearsal.ps1 -RefreshBaseline
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1

# unelevated: free evidence
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\boundary-battery.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\container-candidate.ps1 -Scenario exec-check
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\path-probe.ps1

# elevated close-out: verify -> revert -> verify absent -> optional profile removal
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\verify-b2-loopback.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\revert-b2-loopback.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\verify-b2-loopback.ps1 -Expect absent
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\appcontainer-b2\remove-b2-container-profile.ps1
```

`revert` already deletes the profile, so the last command is only needed when revert was run with
`-KeepProfile` (B4 reuse) and the profile should now go away.

See `../../agent-probe/enforcement-proxy/README.md` (中文版 `README_CN.md`) for the enforcement proxy itself
and for the boundary description these scripts realise.

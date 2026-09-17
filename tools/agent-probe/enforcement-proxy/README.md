# enforcement-proxy (B2 external-enforcement boundary)

> 中文版：[README_CN.md](README_CN.md)

Minimal, auditable forward proxy used by the B2 proof. It is a **tool** (`tools/agent-probe/`), not
production runtime code. Standard library only; no TLS interception (HTTPS travels as a `CONNECT`
tunnel), no third-party dependencies, nothing to install or uninstall.

## Boundary architecture it belongs to

1. The candidate CLI runs inside a **zero-network-capability AppContainer**; the OS then denies all
   public and LAN egress (including DNS) at kernel level.
2. The container SID holds a **loopback exemption**, so the only reachable endpoints are local
   loopback services.
3. This proxy listens on **one fixed loopback endpoint** and is the container's only egress path:
   `container -> 127.0.0.1:<port> (this proxy) -> host egress (direct, or through -upstream)`.
4. The **host-side egress leg is outside the container boundary**: the container cannot dial the
   internet or the LAN directly, and cannot bypass the proxy. The `networkControl` claim is therefore
   "no direct egress; all egress is proxied and allowlisted", not "the host itself has no network".

## Port policy (fixed, never silently changed)

- Default listen endpoint: `127.0.0.1:39877`; the listen address is configurable with `-listen`.
- The listen address must be a **loopback** address (`127.0.0.1`, `::1`, `localhost`). Anything else is
  refused at startup: exposing the allowlist (or, in observe mode, the whole upstream) on a routable
  address would turn it into an open relay for the network.
- At startup the proxy binds the exact endpoint; if the port is occupied it **exits non-zero** and
  logs `event=start decision=error`. It never falls back to another port.
- During B2 the port is frozen at 39877 and all evidence points at it. Changing it afterwards is
  allowed but requires re-running the reachability probe (P1) and the smoke probe (P4) and recording
  the new port plus results in the evidence pack. B4 reuse must use the same port as B2 evidence; if
  it cannot, B4 must re-verify the boundary and state the difference in its report.
- System rules (the AppContainer loopback exemption) do **not** reference the port, so a port change
  adds no restoral work.

## Allowlist semantics

- `-allow` takes a comma-separated list; a rule matches the host itself **and any subdomain**
  (dot-boundary), case-insensitively, ports ignored.
- A rule must be a bare host suffix. Schemas, paths, ports, wildcards and malformed labels are rejected
  at startup instead of being accepted as a rule that can never match.
- Default: `api.githubcopilot.com,api.github.com,github.com,githubusercontent.com,copilot-proxy.githubusercontent.com`.
- Everything else is refused with `403` and logged as `decision=deny reason=not-on-allowlist`.
- `-observe` mode forwards unknown hosts **and records them**, purely for discovery. Observe runs are
  never proof evidence; freeze the allowlist and re-run in enforce mode for the proof.

## Upstream

- `-upstream http://host:port` chains the proxy through an upstream HTTP proxy for the host-side leg.
  The container never sees the upstream; changing it does not touch any OS rule.
- Default is direct egress.

## Audit log (evidence)

JSONL, one object per line:

| Field | Meaning |
|---|---|
| `ts` | UTC timestamp (RFC 3339, nanoseconds) |
| `event` | `start`, `stop`, `connect`, `http` |
| `phase` | `authorized` (pre-flight allow record), `established` then `closed` for `connect`; `authorized` then `allow` for `http` |
| `decision` | `allow`, `deny`, or `error` |
| `reason` | `allowlisted`, `observe`, `not-on-allowlist`, `malformed-request: …`, `dial-failed: …`, `roundtrip-failed: …` |
| `host`, `port` | requested destination (plain-HTTP events record the URL port, defaulting to 80/443) |
| `upstream` | upstream proxy in use, if any |
| `client` | caller address (the container) |
| `bytesToTarget`, `bytesToClient` | tunnel byte counts, present on `connect closed` events only |
| `durationMs` | handle duration for `http` events |

Audit integrity: the log is the evidence for `networkControl`, so a write failure is **fail-closed**. The
allow paths write a pre-flight `authorized` record before any traffic is forwarded; if that write fails the
request is answered with `403` and the tunnel is never established (an already-established tunnel is closed
immediately). The first failure is also reported on stderr.

## Build, run, self-test

```
go build -o enforcement-proxy.exe ./tools/agent-probe/enforcement-proxy
enforcement-proxy.exe -listen 127.0.0.1:39877 -log evidence\proxy.jsonl
enforcement-proxy.exe -listen 127.0.0.1:39877 -allow "api.githubcopilot.com,github.com" -upstream http://10.0.0.246:8080
powershell -NoProfile -ExecutionPolicy Bypass -File tools\agent-probe\enforcement-proxy\Test-EnforcementProxy.ps1
go test ./tools/agent-probe/enforcement-proxy/
```

The self-test runs entirely on loopback (local targets only, no external traffic, no model calls) and
covers: allowlisted forward request, refused request, observe-mode recording, occupied-port fail-fast,
and log parseability. Scratch lives in `<repo>/tmp/enforcement-proxy-selftest/<stamp>` (the repository
`tmp/` is git-ignored) and is removed on success; `-KeepArtifacts` keeps it for inspection.

## Process requirements carried by this tool

1. **Restoral rehearsal before paid probes.** `tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1` →
   `verify-b2-loopback.ps1` → `revert-b2-loopback.ps1` → `verify-b2-loopback.ps1` →
   `apply-b2-loopback.ps1` (orchestrated by `Invoke-B2LoopbackRehearsal.ps1` in the same directory), with the
   revert diff (exemptions + firewall rules: must be zero difference) written into the evidence pack.
2. **Loopback baseline invalidation is executable.** Take a listener snapshot with
   `tools\agent-probe\appcontainer-b2\verify-b2-loopback.ps1` (it snapshots and diffs; `New-B2Snapshot` in its
   `b2-loopback-lib.ps1` is the underlying helper)
   **before P1** and **before any B4 reuse**, and compare against the B2 baseline. Any newly added
   loopback listener is reported to the human owner, who decides whether to continue.
3. **Proxy egress path is stated in the evidence.** The container reaches only
   `127.0.0.1:<port>`; the proxy performs host-side egress (direct or via `-upstream`). The
   `networkControl` claim covers the container, not the host.

## Restoring the machine after an elevated run (close-out)

The elevated run leaves exactly two possible machine-state changes, and this section is the procedure that
returns the machine to its pre-B2 state. The scripts live in `tools/agent-probe/appcontainer-b2/` (see its
`README.md`); their scratch output goes to the git-ignored `tmp/b2/`. Run every command in an **elevated**
PowerShell (the scripts refuse to run otherwise), from the repository root.

| Machine state | Created by | Removed by |
|---|---|---|
| Loopback exemption for the container SID (whole-table write, pre-existing entries preserved) | `tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1` | `tools\agent-probe\appcontainer-b2\revert-b2-loopback.ps1` |
| AppContainer profile `prfrail-b2-battery` | `tools\agent-probe\appcontainer-b2\apply-b2-loopback.ps1` (first run only) | `tools\agent-probe\appcontainer-b2\revert-b2-loopback.ps1` (same run, unless `-KeepProfile`); `remove-b2-container-profile.ps1` for the kept case |

Procedure:

1. `verify-b2-loopback.ps1` — read-only: prints the SID, the profile existence, the current exemption table and
   the strict diff against the recorded baseline. Expect the container SID to be present at this point.
2. `revert-b2-loopback.ps1` — removes the container SID exemption, deletes the AppContainer profile (restoral
   means back to the recorded baseline; `-KeepProfile` keeps it for a B4 reuse) and re-runs the strict diff
   (exemptions + firewall + profile existence). The run must end with **zero difference**; anything else means the
   machine did not return to its baseline and must be investigated before proceeding.
3. `verify-b2-loopback.ps1 -Expect absent` — read-only confirmation that the exemption is gone.
4. Optional: `remove-b2-container-profile.ps1` — deletes the profile only when it was kept (`-KeepProfile`).
5. Archive the transcripts of steps 1–3 into the evidence bundle and recompute `SHA256SUMS.txt`
   (`assemble-evidence.ps1` + `normalize-encoding.ps1`, same directory).

A baseline written before schemaVersion 2 cannot answer the profile-existence question; `revert`/`verify` report
that as an informational `profile-existence-unverifiable` line instead of a difference, and
`apply-b2-loopback.ps1 -RefreshBaseline` (or `Invoke-B2LoopbackRehearsal.ps1 -RefreshBaseline`) records a valid
one — only when the machine is already at restoral state (no exemption, no profile).

The rehearsal (`Invoke-B2LoopbackRehearsal.ps1`) runs apply → verify → revert → verify → apply and is the proof
that this close-out is repeatable; run it in an elevated session and archive its transcript.

What close-out does **not** cover, by design:

- the per-run directory ACL grants under `tmp/b2/` (they disappear with the directories themselves);
- `subst` drive mappings (per-session, released by the harness);
- any host-side effect of the proxy's own egress (outside the boundary; see the known limitations).

## Known limitations

- No TLS interception: HTTPS is opaque to the proxy (allowlist is enforced on the CONNECT target).
- Plain-HTTP upgrade flows (WebSockets) are not tunneled; only standard requests are forwarded.
- `-observe` mode is for discovery only and must never be cited as enforcement evidence.
- The tool does not manage the AppContainer or its exemption; that pairing is a versioned companion tool in
  `tools/agent-probe/appcontainer-b2/` (`apply-b2-loopback.ps1` / `revert-b2-loopback.ps1` /
  `verify-b2-loopback.ps1`), whose as-run snapshot is archived under
  `docs/validation/evidence/b2-2026-09-17/machine-ops/`.

# B3a power-loss injection rig and two-way calibration — verification report

- Slice: **B3a** (power-loss injection rig + two-way calibration), depends on A7
- Date: 2026-09-20
- Verdict: **calibration PASSED (`CALIBRATED`)** — the rig discriminates loss from survival in both directions;
  **B3b may start**
- Evidence bundle: `docs/validation/evidence/b3a-2026-09-20/` (128 files, `SHA256SUMS.txt` with 126 entries, **0 mismatches**)
- Live data: `D:\VirtualBox VMs\Win11B3\b3a-runs\2026-09-20\` (outside the repository)

## 0. Verdict summary

| Item | Result |
|---|---|
| Negative control (deliberately unflushed) | **4 losses in 5 valid rounds** (`lost-zero-filled` ×2, `lost-absent` ×2) + 1 survival → satisfies "at least one loss" |
| Positive control (explicit flush + committed directory entry) | **all 5 valid rounds `survived`** (`positiveLosses = 0`) |
| Journal chain | `journalChainOk = true`; 36 records (12 `plan`, 12 `cut`, 12 `reconcile`), complete after reboot and reconcilable round by round |
| Inventory repeatability | two consecutive inventories after every reboot are **byte-identical on stdout** (`inventoryIdentical = true` ×12) |
| Injection authenticity | **all 12 rounds** report `auditVerdict = HARD-POWER-LOSS` (no round was faked by a graceful shutdown) |
| Volume state | `fsutil dirty query R:` reports `Volume - R: is NOT Dirty` in all 12 rounds |
| Machine verdict | `docs/validation/evidence/b3a-2026-09-20/calibration-verdict.json`, `gate = CALIBRATED` |

## 1. Slice acceptance mapping (item by item)

| Slice requirement (`REMAINING_SLICES.md` B3a) | Evidence in this run |
|---|---|
| The negative control must lose data at least once within N rounds | 4/5 losses (r04/r06 zero-filled, r05/r07 absent), see the table in §4 |
| The positive control must lose nothing | 5/5 `survived` (r08–r12) |
| The journal survives reboot and reconciles round by round | chain verification passes; every round carries `plan`→`cut`→`reconcile` with `outcome` on the last record; `journalRecordCount = 36` |
| The inventory is repeatable (>= 2 identical enumerations) | two byte-identical inventories per round (r01–r12) |
| Freeze the rig scripts and the process, produce an evidence bundle | bundle + this report; frozen script copies in `evidence/b3a-2026-09-20/rig/` (11 files, byte-identical to the repository) |

> Boundary respected: **this slice draws no candidate (C2/C3) conclusion** and does not clear `unproven`; it only decides
> whether the rig can discriminate loss.

## 2. Frozen premise (the conclusions cover only this premise)

| Item | Value |
|---|---|
| Hypervisor | VirtualBox **7.2.18 r175117** |
| VM | `Win11B3`, UUID `02205fb0-1933-48ba-bbdb-a2c06b4eb9b2` (outside the repository: `D:\VirtualBox VMs\Win11B3`) |
| Guest | **retail Windows 11 Pro 25H2 build 26200.8037 (zh-CN, unactivated)**; EFI + vTPM 2.0; 2 vCPU / 4096 MB |
| Storage | SATA/AHCI, `useHostIOCache=false`; `os.vdi` 64 GB dynamic (port 0) + `data.vdi` 8 GB fixed (port 1); **both disks `nonrotational=on`** |
| Injection volume | `R:` = INJECT, NTFS 7.98 GB (on `data.vdi`) |
| Graphics / network | VBoxSVGA (GA WDDM driver bound); `nic1=none` (no real NIC device inside the guest) |
| Baseline | snapshot **`pristine`** (referenced by NAME) = `065a25e1-d95c-4f33-8081-4f06216ecce8` |
| Install media | `Win11_25H2_Chinese_Simplified_x64_v2.iso`, 8 543 608 832 B, sha256 `7408581e67bc455ebaafb9230e531abf45b1c8864a22114a1b03893f897102e4` (Image #4 = Pro) |
| Difference vs A7 | the A7 host is 24H2 build 26100; this guest is **25H2 build 26200** (same GA line, different build) |
| Difference vs design §1.6 | the evaluation-edition `TIMEBASED_EVAL` hourly self-shutdown defect is gone after the retail rebuild (70-minute idle run, no self-shutdown); the window/assertion/audit rules are now a **cheap safety net**, not licence evasion |

**Premise correction made by this run**: after the rebuild both disks had lost `--nonrotational on`; the S0 gate caught it
(`2-vm-facts` FAIL). Both flags were re-applied and the **baseline snapshot was re-taken** (`7e0287b6…` → `065a25e1…`).
A snapshot also stores machine configuration, so any configuration change requires re-taking it.

## 3. Rig structure and protocol highlights

- **Scripts (11)**: `b3a-rig-lib.ps1` (shared helpers), `b3a-journal-lib.ps1` (journal module), `Get-B3aReconcile.ps1`
  (reconcile **module**, no top-level parameters), `Invoke-B3aReconcile.ps1` (its CLI wrapper), `New-B3aSession.ps1` (S0),
  `Invoke-B3aRound.ps1` (one round, S1–S9), `Invoke-B3aCalibration.ps1` (calibration orchestration), `New-B3aEvidence.ps1`
  (bundle), `Test-B3aRig.ps1` (hermetic self-test), `guest/12-round-inventory.ps1`, `guest/13-round-write.ps1`
  (plus the reused guest-side `11-shutdown-audit.ps1`).
- **Per-round steps (all observed)**: S1 restore `pristine` → S2 start → S3 ready wait → **S3.5 tool deployment (repeated
  every round)** → S4 pre-injection audit → **S4.5 round-root setup** → **S5 `plan` record + ACK (journal first)** →
  S6 guest control write with a bi-directional probe/plan check → S7 hard poweroff (`controlvm poweroff`) + `cut` record →
  S8 restart + shutdown audit + inventory ×2 → S9 reconcile + `reconcile` record.
- **Frozen journal hash convention**: `recordHash = SHA256(UTF8(prevHashHex | "null") + UTF8(stored line text))`; the stored
  line is the canonical JSON with its closing `}` replaced by `,"recordHash":"<hex>"}`; verification **requires the stored
  line to equal the canonical reconstruction byte for byte**.
- **"Journal first"**: the `plan` record's digest/size are computed on the **host** from the seed (a pre-commitment, never
  copied from the probe); the guest write is allowed only after the ACK, and the probe's digest/size must match the plan
  or the round is voided without an injection.

## 4. Calibration result (12 rounds, all valid)

Round 1 was a **single-round smoke test** run before the calibration to validate the repaired S1–S9 chain; the 11 official
rounds then continued in the same session (the journal-based round numbering resumed correctly:
`journal chain ok (3 records); next absolute round = 2`).

| Round | Slot | Candidate | Stage | Outcome | Reading |
|---|---|---|---|---|---|
| r01 | — | parity | parity | `lost-zero-filled` | smoke round (reproduces the known mode) |
| r02 | R0 | parity | parity | `lost-zero-filled` | parity round (design §4.4 expected mode) |
| r03 | R1 | negative-control | temp-written | `survived` | no loss (allowed for the negative control) |
| r04 | R2 | negative-control | temp-written | `lost-zero-filled` | loss |
| r05 | R3 | negative-control | temp-written | `lost-absent` | loss |
| r06 | R4 | negative-control | temp-written | `lost-zero-filled` | loss |
| r07 | R5 | negative-control | temp-written | `lost-absent` | loss |
| r08 | R6 | positive-control | parent-synced | `survived` | survived |
| r09 | R7 | positive-control | parent-synced | `survived` | survived |
| r10 | R8 | positive-control | parent-synced | `survived` | survived |
| r11 | R9 | positive-control | parent-synced | `survived` | survived |
| r12 | R10 | positive-control | parent-synced | `survived` | survived |

- Payload 256 KiB (262144 B), injection delay fixed at **2.0 s**, `restore pristine` every round, hard cut via
  `VBoxManage controlvm poweroff`.
- **Discriminating power**: negative-control loss rate **4/5 = 0.8**; positive-control loss rate **0/5 = 0.0**. Under this
  premise and payload the rig reliably produces observable loss, and the explicitly flushed chain loses nothing — which is
  exactly the two-way discrimination this slice must prove.
- **Timing**: the official calibration ran 04:34:38 → 05:13:20 = **38 min 42 s for 11 rounds** (≈3 min 31 s per round); the
  smoke round took 3 min 44 s.

## 5. Per-round evidence and consistency

Every `rounds/rNN/` directory holds **9/9 artifacts** (nothing missing):

`plan.json`, `cut.json`, `write-probe.txt`, `audit-pre.txt`, `audit.txt`, `inventory-1.json`, `inventory-2.json`,
`dirty.txt`, `reconcile.json`

- `auditVerdict = HARD-POWER-LOSS` ×12 (`GRACEFUL-SHUTDOWN`/`LICENCE-SHUTDOWN` never appeared, so no round faked survival);
- `inventoryIdentical = true` ×12 (repeatable inventory; reconciliation compares content and size only, never time);
- `dirty = "Volume - R: is NOT Dirty"` ×12 (no dirty volume, no `chkdsk` recovery needed);
- the probe proves the write API: the negative control reports `apisCalled = ["File.WriteAllBytes"]` (no flush API at all),
  while the positive control reports `FileStream + Flush(true) + CreateFile(\\.\R:) + FlushFileBuffers + CloseHandle`
  (`commit = positive-commit-ok`).

## 6. Journal self-proof

- 36 records = 12 `plan` + 12 `cut` + 12 `reconcile`; every write path is
  `FileStream(Append, FileShare.Read, WriteThrough)` → `Flush($true)` → close → **read-back hash recomputation plus a
  byte-exact line-shape check** → ACK.
- The session starts with a **full chain verification** (`journalChainOk = true`); the calibration counts only the rounds
  produced by its own invocation and additionally requires each round's `inventoryIdentical` and `auditVerdict` to pass
  before it can report `CALIBRATED`.
- Physical domain isolation: the journal lives on the **host** NVMe (`D:\VirtualBox VMs\...`), outside the guest's power
  domain; the injection only terminates the VM process. Honest boundary: a host power loss is outside this slice's
  injection domain (consistent with design §2.1).

## 7. ④ Independent review record (Codex, read-only)

| Pass | Verdict | Highlights |
|---|---|---|
| 1 | **FAIL** | 3 High (journal order violated "journal first"; `Wait-B3aReady` missing `-TimeoutS`; `-N` had no lower bound) + 3 Medium + 2 Low |
| 2 | **FAIL** | all 8 fixes confirmed `fixed`; but one new High: the hash covered only the prefix before `recordHash`, so an appended tail could keep the chain green |
| 3 | **PASS** | that High confirmed `closed`; no Medium+ remaining |
| 4 (focused) | **FAIL** | the three late-defect fixes themselves hold, but 3 Medium findings targeted **guard strength** (a guard could pass vacuously, `files` was not asserted to be an array, and the exhausted-budget branch missed the child output) → all strengthened |

> The call budget (1 + 1 re-review) was exceeded. Reason: B3a is a hard-gate slice and a new High plus new fix batches
> appeared mid-slice, so the closed-loop rule ("④ until no Medium+; never self-certify") was followed; the cost is recorded honestly.

## 8. Defects found by step ⑤ (native verification), missed by both the self-test and the reviews

| # | Defect | Severity | Consequence | Fix + regression test |
|---|---|---|---|---|
| D1 | `Get-B3aReconcile.ps1` had a top-level `param()` yet was dot-sourced | **Critical (silent)** | dot-sourcing rebinds same-named caller variables → every round's `-Round`/`-Candidate` became `0`/`negative-control`, invalidating the whole calibration (the first attempt stalled at `r00`) | the module **must have no top-level parameters**; the CLI moved to `Invoke-B3aReconcile.ps1`; self-check `module-load-preserves-vars` replays the dot-source list extracted from the caller and requires a non-empty list containing the known modules |
| D2 | parent and child appending to `session.log` concurrently (`StreamWriter` defaults to `FileShare.Read`) | High | sharing violation plus fail-closed killed both processes with exit 1 and lost the explaining log line | `FileShare.ReadWrite` with up to 5 backoff retries; child stdout/stderr is logged on non-zero exit, and before the retry-budget check |
| D3 | `[pscustomobject]@{ files = @($List[object]) }` | High | the PS 5.1 binder throws `ArgumentException` ("parameter type mismatch") → `Invoke-B3aToolDeploy` could **never return**, so S0 check 6 and every round's S3.5 failed | use `.ToArray()` (two sites); self-check `deploy-return-shape` (drives the return path with a bogus `VBoxManage`, then requires `files` to be an array of 3) |
| D4 | the round root directory was never created | High | the measured write died with `DirectoryNotFoundException`; the guest script only declares `exit 0/1`, VBoxManage mapped it to 33, so the host saw an unexplained 33 | new **S4.5** (guest `mkdir R:\calib\rNN`, rig-side housekeeping that is *not* counted as the measured write); on a failed S6 the probe stdout/stderr are archived (`write-probe.txt`/`write-stderr.txt`) |
| D5 | `New-B3aEvidence.ps1` used a non-normalized `$BundleDir` (`Join-Path '..\..'`) | Medium | the literal path is longer than the real file paths → `Substring` went out of range, failing the very first real packaging run | normalize with `[System.IO.Path]::GetFullPath()` plus an explicit out-of-root block |

Same-class sweep (that review pass): no other script dot-sources a file with top-level parameters; the rig now consistently
uses `.ToArray()` instead of `@(<List[object]>)`.

## 9. Counterexample / mutation evidence (the tests are not vacuous)

| Test | Counterexample / mutation | Result |
|---|---|---|
| `journal-timestamp-roundtrip` | revert the verifier to "parse then re-serialize" | **fails deterministically** (the 7th fractional digit's trailing zero is dropped: `…8416210+08:00` → `…841621+08:00`; the old code misjudged ~1 in 10 records, so ~69% of 11-round sessions would have aborted with `journal chain broken`) |
| `journal-tail-append-red` | remove the byte-exact line-shape check | **fails deterministically** (the chain stays green after a tail field is appended) |
| `verdict-requiredN-floor` | `RequiredN = 1` | must throw (refuses a `CALIBRATED` verdict from a single-round sample) |
| `module-load-preserves-vars` | dot-source a module with a top-level `param()` | goes red as soon as a variable is rebound (class-level guard for D1) |
| `deploy-return-shape` | the `@(List[object])` binder bug | goes red on an exception or when `files` is not an array (class-level guard for D3) |

Hermetic self-test size: **22 checks, 0 failures** (10 consecutive fresh processes plus the bundle-time run, all
`checks=22 failures=0`). Static analysis: `Invoke-ScriptAnalyzer -Severity Warning,Error` = **0** (including the 14 style
warnings cleared in this slice) - the raw output is archived as `verification/analyzer.txt` in the bundle (`findings=0`). Encoding: all 11 rig scripts and the bundle's `.ps1`/`.md` files are **UTF-8 BOM + LF**;
the bundle's `.json`/`.txt` are BOM-free + LF.

## 10. Known limitations and honest boundaries

1. **Single-device calibration**: the conclusions cover this premise only (VirtualBox + VDI backend + host NVMe + NTFS +
   the builds above). Changing the hypervisor, the disk backend, or the injection delay/payload size ⇒ **recalibrate**.
2. **The loss rate is not a candidate conclusion**: 4/5 for the negative control only shows the rig can produce observable
   loss; the 5-stage × N × candidate durability verdict belongs to **B3b**. No single-round or single-sample extrapolation here.
3. **Journal threat model**: a keyless hash chain gives **tamper visibility**, not resistance to an actor that can rewrite
   the whole chain and recompute it; that boundary is covered by freezing the bundle and verifying the whole-bundle hashes.
   A host power loss is outside the injection domain.
4. **Bundle byte accounting**: the bundle's `rig/` copies are byte-identical to the repository, and all 10 measurement-path
   scripts were last modified **before** the calibration started (04:34:38) — i.e. "frozen bytes == bytes that ran"; only the
   packaging script `New-B3aEvidence.ps1` (05:14:28) was corrected afterwards for D5 and is not on the measurement path.   Both statements are independently checkable from `verification/rig-vs-repo.txt` (per script `identical`, `repoMtime` and
   the judgement against the calibration start: all 10 measurement-path scripts read `before-calibration-start`, only the
   packaging script reads `AFTER-CALIBRATION-START`).5. **The session directory contains a smoke round**: `r01` is the pre-calibration single-round smoke test (also a
   parity/negative-direction write), flagged in §4; the calibration verdict counts only the 11 rounds of its own invocation
   (see `negativeValidRounds`/`positiveValidRounds` in `calibration-verdict.json`).
6. **Credential file: retention decision pending**: `D:\VirtualBox VMs\Win11B3\host-only\cred.txt` is required for B3b to
   reuse this environment; deleting it now would make the environment inoperable. Recommendation: **keep it until B3b
   finishes**, or set a new one-time password before B3b. The file has always been outside the repository and was never printed.

## 11. B3b hand-off

B3b **must** reuse the frozen: injection delay **2.0 s**, per-round `restore pristine`, the S1–S9 round protocol (including
S3.5/S4.5), the journal convention and the reconcile/inventory tooling, and the premise statement. Changing any of these ⇒
**recalibrate**. B3b decides its own 5-stage × N × candidate matrix and N (this slice publishes the per-round negative
distribution as input for choosing N).

## 12. Reproduction commands (host PowerShell 5.1)

```powershell
# S0 session preflight (everything must be green before any round starts)
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aSession.ps1

# one round (waits for operator confirmation before each cut unless -NoConfirm is given)
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aRound.ps1 `
  -Round 1 -Candidate parity -NoConfirm

# two-way calibration (the run behind this report)
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Invoke-B3aCalibration.ps1 -NoConfirm

# hermetic self-test / evidence bundle
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\Test-B3aRig.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools\b3a-rig\New-B3aEvidence.ps1
```

## 13. Deliverables and open items

**Deliverables**: (1) environment and premise statement (§2); (2) the 11 rig scripts (`tools/b3a-rig/`); (3) the two-way
calibration result (§4); (4) per-round journal, inventory and reconcile artifacts (bundle `rounds/`, `journal/`); (5) the
evidence bundle (`MANIFEST.md` + `SHA256SUMS.txt`, 0 mismatches); (6) this report (bilingual).

**Open items**:
1. Credential file retention/rotation decision (§10 item 6).
2. The retired evaluation-edition environment once echoed its password in clear text — rotate that (retired) password if
   desired; it does not affect this slice.
3. Before starting B3b, confirm it reuses this slice's frozen parameters (§11).

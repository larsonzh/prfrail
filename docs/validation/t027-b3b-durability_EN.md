# T027 · B3b · Candidate injection matrix and the three-tier durability verdict — validation report

Date: 2026-09-21 | Status: **COMPLETE (report written; not committed, not pushed)** | Executing owner: DeepSeek V4.1 Flash (the ⑥ stage ran with the owner tier set to Deep/Max)
Authorities: `docs/t027/B3B_DESIGN.md` (single authority for design and remediation) | `docs/t027/B3A_RIG_DESIGN.md` (frozen rig) | Evidence pack: `docs/validation/evidence/b3b-2026-09-21/`

> **The tier and the premise must be read together**: this report uses only the three tiers `proven` / `unresolved` / `disproven`.
> Any `proven` holds **only under premise P** and **must not** be read as a long-term guarantee, universal durability, a cross-platform result, or AT-23 evidence.

---

## 0. Honesty rules (the premise of every conclusion here)

1. **No inference survives outside premise P** (§1); the conclusions do not cover physical hardware, other hypervisors/storage stacks, other filesystems or payload sizes, and not the host losing power.
2. **Three tiers**: `proven` / `unresolved` / `disproven`. `unresolved` is **not** a failure and **not** "not durable"; it means the claim could not be established under premise P and the device's resolution.
3. **The five conditions for `proven`** (design §7): (i) all 5 cells × N=5 valid rounds inside the allowed set; (ii) **mechanism evidence complete in the decision cells**; (iii) the discrimination gate holds; (iv) the global invariants hold; (v) the journal chain and plan/cut/reconcile are complete.
4. **Tests cannot certify themselves**: the self-test and its assertions only demonstrate falsifiability (a deleted guard reddens a named assertion, §6/§9); they are not the conclusion. The only source of the conclusion is ⑤'s native evidence.
5. **Forbidden wording**: no long-term-guarantee phrasing (`durable` / `guarantee` / `survives power loss`); only the verbatim `claimScope` sentence of §4.

## 1. Scope and boundaries

**Premise P (frozen; verbatim in `matrix/verdict.json` as `premise`)**

- hypervisor **VirtualBox 7.2.18 r175117**; guest **Windows 11 25H2 build 26200.8037** (retail, unactivated); filesystem **NTFS**.
- Storage: SATA/AHCI, **`useHostIOCache=false`**, both disks `--nonrotational on`; injection volume **`R:` = INJECT** (`data.vdi`, 7.98 GB, 4096 B clusters).
- Injection: **guest hard power-off** (`VBoxManage controlvm poweroff`); baseline snapshot `pristine` = `065a25e1-d95c-4f33-8081-4f06216ecce8` (not retaken in this slice).
- `nic1=none` (no network stack in the guest). **Signalling-domain declaration**: `C:\prfrail-prep\` is a second write domain used for stage markers and **does not take part in the measurement**; the only measured domain is the replay store root on `R:`, with **one** measured publication per round.
- Payload: record level, **~1 KiB** (1082 B); pre-cut freeze delay **2.0 s** (frozen B3a calibration, not retaken); **N=5** rounds per cell.

**Not covered (do not extrapolate)**

- Physical hardware/firmware write-back and true mains loss (this is a hypervisor-level hard power-off); other hypervisors, other filesystems, SMB/shared volumes/cloud disks (U5); the host itself losing power (U3).
- `.tmp` residual cleanup, concurrent multi-writer publication, resume/settlement and other paths.
- **AT-23**: this slice is **not** AT-23 evidence, does not lift the Windows `first-dispatch` refusal, and does not clear other platforms' `unproven`.
- Statistical resolution: a clean cell gives a one-sided 95% upper bound of p<0.45, so it can only read as "no counterexample observed under this premise and this device's resolution".

**Device fidelity (not proven here; kept verbatim)**: the slice observed one late-visibility event (an unlinked temp file's size read as 0 and then 1082, see R33).
The evidence favours the reading that the first enumeration saw a not-yet-settled metadata view, but that evidence **cannot** rule out host/virtual-disk layer visibility arriving after the cut.
This report therefore **must not** be read as "this device has been shown equivalent to a physical power loss".

## 2. Candidates and the rig

| Candidate | Semantics | Landed at |
| --- | --- | --- |
| **C1** (reference) | current baseline: `os.Link` no-replace plus the platform parent-directory step (a no-op on Windows) | the default path itself; **never `proven`** |
| **C2** | a no-replace move via `MoveFileExW` + `MOVEFILE_WRITE_THROUGH (0x8)` replacing `os.Link` | publish primitive hook `replayStorePublishPrimitiveHook` |
| **C3** | a directory-handle `FlushFileBuffers` replacing the empty parent-directory step on Windows | parent-directory hook `replayStoreParentSyncHook` |
| **C2C3** | C2 + C3 combined | both hooks |
| **positive-control** | baseline primitive plus directory/record handle flush (falls back and labels itself when the volume handle is unavailable) | both hooks; claims **no** decision cell, only calibrates the device |

- Both seams are **nil by default and zero-semantic**: the production path is byte-for-byte unchanged (checked at ③ and ④).
- **D3 (gate bypass, declared)**: the experimental writer builds in `publicationDurability=proven` and calls `writeReplayRecordNoReplace` directly, only to reach the measured publication point; production gates are untouched.
- **Variant arm (device self-falsification, a manual arm, not part of the 103 rounds)**: the compile-time constant `b3bVariantDegrade` (`-tags "b3bnative b3bvariant"`) degrades every candidate to the baseline primitive, while the trace keeps naming what actually ran. See §3.4.

**Frozen protocol (design §2/§4, reused from B3a)**

- Five publication stages: `temp-written` → `temp-synced` → `temp-closed` → `record-linked` → `parent-synced`.
- **Decision cells**: C2 = {record-linked, parent-synced}; C3 = {parent-synced}; C2C3 = {parent-synced}; C1 and the control = ∅.
- **Stage attribution uses the pre-cut `marker-transcript.txt`**; an escape becomes a `sessionBlock` (finalize exit 6).
- **Fail-closed rule (the crux of this slice's verdicts)**: a valid round in a decision stage with `mechanismOk=false` must not make its cell `eligible`, and is reported as `mechanism evidence missing in rounds N`.
  Meaning: **a "survived" whose mechanism was never shown to have run cannot be cashed into a durability claim** — this is where the rig's self-falsification ability formally lives.

## 3. Results

### 3.1 Gates (Windows host, final run before hand-off)

| Item | Result |
| --- | --- |
| Rig self-test `tools/b3b-rig/Test-B3bRig.ps1` | **`checks=103 failures=0`** |
| `Invoke-ScriptAnalyzer -Path tools/b3b-rig -Recurse` | **0** |
| `gofmt -l` (touched files + `internal`) | empty |
| `go build ./...` / `go vet ./...` / `go vet -tags b3bnative ./internal/adapters/` | all ok |
| `go test -count=1 ./...` | **14 packages ok, 0 FAIL** |
| native suite `go test -tags b3bnative -count=1 ./internal/adapters/` | **exit 0** (25 s) |
| Encoding and line endings | every touched `.md`/`.ps1` = UTF-8 **with BOM** + LF; `.go`/`.json` = **without BOM** + LF |

### 3.2 Main matrix (session `2026-09-21`, writer sha256 `06de8a72…`)

- **17 cells × 5 = 85 rounds, all valid: `ok=85 void=0`**; zero operational failures (the single operational failure, round r07, happened on the pre-fix revision and was cleared by R34; see §9.4).
- Calibration / negative control (the discrimination gate): the positive control `positive-control:parent-synced` is **5/5 survived**; the negative control `c1:temp-written` is **5/5 lost** (plus `c1:temp-synced` and `c1:temp-closed`, 5/5 lost each).
- Per-cell results (survived / lost):

| Cell | survived | lost | mechanism token present |
| --- | --- | --- | --- |
| `positive-control:parent-synced` | 5 | 0 | 5/5 |
| `c1:temp-written` / `temp-synced` / `temp-closed` | 0 / 0 / 0 | 5 / 5 / 5 | — |
| `c1:record-linked` | 4 | **1** | 5/5 |
| `c1:parent-synced` | 5 | 0 | 5/5 |
| `c2:temp-written` / `temp-synced` / `temp-closed` | 0 / 0 / 0 | 5 / 5 / 5 | — |
| `c2:record-linked` | **5** | 0 | **5/5** |
| `c2:parent-synced` | **5** | 0 | **5/5** |
| `c3:temp-written` / `temp-synced` / `temp-closed` | 0 / 0 / 0 | 5 / 5 / 5 | — |
| `c3:record-linked` | 5 | 0 | 0/5 (not its decision stage) |
| `c3:parent-synced` | 5 | 0 | **4/5** (see §4.2) |
| `c2c3:parent-synced` | 5 | 0 | **4/5** (see §4.2) |

- **Measured inventory reads**: `reads=2` or `reads=3`, **never 4** (R33's extra cost averages one inventory read, ~27 s/round).

### 3.3 U4 ordering sub-matrix (`-Paired`, round numbers ≥500)

- **6 cells × 3 = 18 rounds, all valid: `ok=18 void=0`**; all 18 rounds have `targetKind=completion` (enforced by finalize).
- **True orphan C count = 0** across all six cells.

### 3.4 Variant arm (D1, manual arm, separate rig root `b3b-runs-variant`)

- 2 rounds of `c2:parent-synced`: `survived` **and** `stageHit=true`, but `mechanismOk=false` (missing `MoveFileExW(WRITE_THROUGH,no-replace)`); the trace is verbatim `… os.Link … platform-parent-sync …`, so the rig **refuses** to cash those two survivals into a claim.
- Offline corroboration: under the variant build `TestB3bNativeCandidateTraceNamesTheRealPrimitive` **FAILS** (the diagnostic names the degraded primitives verbatim); under the default build the same test **PASSES** and the native suite is exit 0.
- Artifacts were merged into the evidence pack under `variant/` (with an authored `README.md`).

### 3.5 Evidence pack

`docs/validation/evidence/b3b-2026-09-21`: `files=2157 rigMismatches=0 missingRequired=0`;
`optionalSectionsAbsent=calibration,u4,guest` (**deviation D4**: every cell — including the positive control and the U4 paired rounds — is written into `rounds/`, with paired rounds flagged `paired=true`; the absence of those optional sections **must not** be read as "calibration/U4 did not run").

## 4. Verdict ledger

### 4.1 Three-tier verdict (`-Finalize` exit 0, 2026-09-21 12:08:21; source `matrix/verdict.json`)

| Arm | Tier | Basis (machine output) |
| --- | --- | --- |
| **C2** | **`proven`** | `decisionSets=[record-linked,parent-synced]`, both cells **`eligible`** (5/5 inside the allowed set with mechanism evidence) |
| C1 | `unresolved` | the reference arm has no decision cell (`decisionSets=[]`); `infoViolations=["record-linked"]` (that cell is `disproven`, see 4.2) |
| C3 | `unresolved` | `decisionSets=[parent-synced]` → `notEligible`, reason verbatim **`mechanism evidence missing in rounds 77`** |
| C2C3 | `unresolved` | same, reason verbatim **`mechanism evidence missing in rounds 81`** |
| positive-control | (calibration) | 5/5 survived, `eligible`; claims no durability |

**Discrimination gate (anti-false-`proven`)**: `ok=true`, negative control (c1 stage 1) **5 rounds, 5 losses**, positive control **5 rounds, 0 losses** (`requiredN=5`).

**U4 (ordering)**: `unresolved`, `validRounds=18`, `orphans=0`, note verbatim:
`ordering can only be falsified by a counterexample, never proven by a clean sample`.

**Claim scope (`claimScope`, quoted verbatim)**:

> within premise P (VirtualBox Win11B3 / NTFS on SATA-AHCI, useHostIOCache=false, record-level ~1 KiB payload, 2 s cut window, N=5 per cell): no counterexample observed, candidate mechanism demonstrably executed, device discrimination demonstrated at stage temp-written only; not a long-term guarantee and not a claim about other hypervisors, filesystems, payload sizes or the host losing power

**Qualifiers that must be stated sentence by sentence**:

- (a) the positive control actually performed a **volume/record handle flush** (the trace records which, and labels a fallback); the gate may only be read as "the device observed a clean result under that flush", **not** as "the volume-level device cache has been proven drained".
- (b) `c2c3` shares `c3`'s decision point at `parent-synced` and its combined primitive is equivalent to C3 there, so `c2c3` may only be reported as a **diagnostic/repeat arm**, never as an independent durability claim.
- (c) stage attribution comes from the **pre-cut** transcript, not the post-cut copy-back; C1's 1/5 loss at stage 4 is evidence of the device's discrimination only and is not extrapolated to production behaviour.

### 4.2 Verbatim per-cell refusals and falsifications

| Cell | Tier | Reason (machine output, verbatim) |
| --- | --- | --- |
| `c1/record-linked` | `disproven` | `round 21: lost-absent is outside the allowed set for record-linked` (that round had `stageHit && mechanismOk`, so it is attributable as a counterexample) |
| `c3/parent-synced` | `unresolved` | `mechanism evidence missing in rounds 77` |
| `c2c3/parent-synced` | `unresolved` | `mechanism evidence missing in rounds 81` |

## 5. Key decisions (anti-misreading)

### 5.1 Why C2 can be `proven` while C3/C2C3 cannot

C2's two decision cells total 10 rounds, all inside the allowed set, with complete mechanism tokens, a discriminating device and a complete journal chain.
C3 and C2C3 each have **1/5 round** in their decision cell where the cut landed before the claimed mechanism ran (no `FlushFileBuffers(dir)` token in the trace; rounds r77 and r81).
The rig fails closed and **refuses** to count a survival whose mechanism never ran as evidence for that arm — counting it would let a candidate claim durability it did not produce.

### 5.2 Conclusions that must **not** be read out of this report

- **Do not** write "c3/c2c3 are not durable" or "c3/c2c3 were disproven": they merely **failed to establish a claim** (`unresolved`).
- **Do not** write "U4 proved ordering": zero orphans only means "not falsified".
- **Do not** claim causality beyond "the candidate is better than the baseline at the decision cells" (except C1's 1/5 loss at stage 4, which is device-discrimination evidence).
- **Do not** claim "the journal chain detects round-artifact tampering": the chain only verifies journal lines; artifact integrity rests on re-derivation from plan+inventory and per-file hash comparison; a wholesale rewrite of artifacts and journal can only be caught by the pack checksums and the S0 pin, the two external anchors.
- **Do not** claim "the evidence pack recomputes itself automatically": only that it can be recomputed from `rounds/` + `journal.jsonl`.

## 6. Changes and test points (falsifiability)

- **Rig**: 10 scripts under `tools/b3b-rig/` (new) plus the evidence builder; four `tools/b3a-rig/` files are **reused as-is** (their hashes are recorded in the session pin).
- **Go side**: `internal/adapters/agent_runner_replay_store.go` (two nil-default seams plus the dispatcher),
  `agent_runner_replay_store_b3b_test.go`, `agent_runner_b3b_candidates_windows_test.go`,
  `agent_runner_native_b3b_windows_test.go`, and the two variant-arm files (`agent_runner_b3b_variant_default_test.go`, `agent_runner_b3b_variant_degraded_test.go`).
- **103 assertions** (`checks=103 failures=0`), including 6 tamper-refusal fixtures, 1 truncated-marker tolerance fixture, 6 in-repo evidence-pack assertions and 3 variant-arm wiring assertions.
- **Falsification experiments (delete the guard / revert the logic ⇒ a named assertion must redden)**, all actually executed:

| # | Mutation | Expected red | Observed |
| --- | --- | --- | --- |
| 1 | `-Cell` pattern reverted so it rejects hyphenated arms | `control-arm-cell-is-schedulable` | **red** ✓ |
| 2 | settled-pair scan reverted to "first pair only" | `inventory-settles-on-the-second-pair` | **red** ✓ |
| 3 | finalize empty-session guard disabled (`-eq -1`) | `finalize-refuses-a-session-with-no-rounds` | **red** ✓ |
| 4 | preflight lock acquisition removed (**A11 sample**) | `preflight-takes-the-session-lock` | **red** ✓ |
| 5 | real-primitive assertion run under the variant build | `TestB3bNativeCandidateTraceNamesTheRealPrimitive` | **red** ✓ (green under the default build) |

## 7. Known boundaries and follow-ups

- **Device discrimination risk** (design §10.1): at ~1 KiB record payload only 1/5 rounds falsified stage 4; if a larger payload or narrower window shows no loss, the ceiling on any claim is still `unresolved`.
- **Stages 4/5 are observationally equivalent** (A7's CP4≡CP5), so a failed marker channel leaves only `unresolved` (r77/r81 are instances of exactly that class).
- **Fidelity is not proven** (end of §1): do not claim the device equals a physical power loss; that claim would need an independent check (not done here).
- **Timing basis**: measured cost is ≈5.3–7.3 min/round, not the earlier 3.8–4.6 min estimate; recompute from each round's `inventory-reads.txt` distribution rather than the old "two inventories" basis.
- **D4**: the evidence pack's `calibration/` and `u4/` sections are absent (see §3.5).
- Follow-up: using `proven` in a product contract needs its own protocol-first contract slice; the Windows `unproven` state, the `first-dispatch` refusal and the AT-23 conclusions are **unchanged**.

## 8. Cost and accounting

- **Probe budget: 0/10 used** (this slice is fully offline: no paid probes, no real model calls).
- **⑤ wall clock**: `01:25:20` (main-session S0) → `12:08:21` (finalize) = **10 h 43 min**; evidence pack at `12:09:13`.
  Breakdown: calibration 5 rounds 26 min; c1 (including one operational failure retried twice as whole rounds) 61 min; the other 14 main-matrix cells ≈7 h 22 min; U4 18 rounds ≈1 h 35 min; variant arm (including the variant preflight) ≈12 min.
- **③ V4 Pro**: pre-analysis ×1; pre-review ×1 (`changes-required`, F1/F2/F3 blocking → R1–R12); closing re-review **PASS** (round counts per the design note §13 table);
  plus **1 call that failed on a provider connection** (no budget consumed, logged under the §3.9 availability pre-check).
- **③.5 MAI-Code-1.1-Flash**: first round ×1 (M1/M2) → remediation (R13/R14) → closing **PASS** (2 calls total, within the ≤3 of §3.10).
- **④ GPT-5.3 Codex: 6 calls** (**above the §2.3 1+1 cap; recorded truthfully per ruling D2**). Purpose and findings per call:

| # | Purpose | Verdict | Findings |
| --- | --- | --- | --- |
| 1 | first final review (no checklists attached) | `FINDINGS` | 4 items (C1–C4) |
| 2 | re-review after the C1–C4 remediation | `FINDINGS` | 5 items (MAJOR#1–#5) |
| 3 | closing #1 (after P1–P4 and the 5 items) | `FINDINGS` | 5/5 historical items closed + **1 new MAJOR (R29)** + 1 coverage gap (R30) |
| 4 | closing #2 (after R29/R30) | **`PASS` (explicit GO)** | 0 |
| 5 | targeted delta review (P5: the `-Cell` resolver) | **`PASS`** | 0 (plus 1 falsifiability strengthening → assertion added) |
| 6 | targeted delta review (P7: settled inventory pair) | **`PASS`** | 0 (plus 1 evidence-retention gap → `inventory-mismatch-*` added) |

  **Why the budget was exceeded (not an executor violation)**: the user requires that "after remediation the stage itself must be re-run" (the closure rule), and ⑤ surfaced two new review points (P5: the control arm was unschedulable; P7: the inventory repeatability contract). Re-running ④ after each remediation is an **obligation** under that rule, not a choice.
- **Tokens**: to be filled in by the user from the usage UI (the owner does not estimate them).
- **Gate wall clock**: self-test ≈20 s per run (≈10 runs this slice); native suite 25 s; `go test ./...` ≈40 s.
- **Bilingual consistency check (§四.10)**: the B3b status row, checkboxes, completion row and header symbol are structurally identical across `REMAINING_SLICES{,_EN}.md` (13 `B3b` markers each, `✅` on both sides); the ADR-013 row in `ADR_REGISTER{,_EN}.md` and the B3b paragraph in `DEV_PLAN{,_EN}.md` are mirrored as well; the numbers, commands and conclusions in this report and in `_EN` agree item by item. **Inconsistencies: 0**.

## 9. Review record (③ / ③.5 / ④) and closing loop

### 9.1 ③ V4 Pro

- Pre-review (after implementation): **`changes-required`** (F1/F2/F3 blocking) → remediation R1–R12 → closing re-review **`PASS`**.
- Focus: decision cells contradicting the mechanism-token rule (R1); the positive control having no driver, which made the discrimination gate structurally and permanently false (R2);
  an over-strong positive-control flush claim (R3); stage attribution reading the post-cut copy-back (R4); all closed and pinned by assertions.

### 9.2 ③.5 MAI (single-layer independent scan)

- First round **`INDEPENDENT SCAN: FINDINGS`**: 2 items (M1/M2, design note §13) → remediation R13/R14 → closing **`PASS`**.
- **A12 independent corroboration**: since MAI reported PASS, ⑤ must contain at least one dedicated experiment that independently proves the key invariants ⇒ the **variant arm** (device self-falsification, §3.4) plays that role: it gives a judgement about "is the mechanism check really guarding anything" that is independent of MAI and reproducible offline.

### 9.3 ④ Codex (6 calls, itemised in §8)

- Coverage: security (bypass paths, injection, zero-value boundaries), architectural consistency (slice goals and contracts), completeness (positive/counterexample/boundary cases),
  **test-counterexample review** (weak tests, happy-path-only, uncovered races), and the falsifiability audit (Section D).
- **A11 sampling (symmetric, actually executed)**:
  - ④ claimed "removing the preflight lock acquisition ⇒ `preflight-takes-the-session-lock` reddens" → **observed red** ✓ (this report's §6 row 4).
  - ④ claimed "reverting the settled-pair scan to first-pair-only ⇒ `inventory-settles-on-the-second-pair` reddens" → **observed red** ✓ (§6 row 2).
  - **Gap declared**: ③.5 MAI's closing round was a `PASS` and emitted no executable "delete guard → reddens" entry, so the sample could not be drawn from MAI.
    Recorded as a known limitation (sampling provides **existence** corroboration only, not a coverage proof).
- **Closing-loop record** (performed after the user's compliance challenge):

| Stage | Closing result | Note |
| --- | --- | --- |
| ③ V4 Pro | **PASS** | no new substantive finding |
| ③.5 MAI | **PASS** | — |
| ④ closing #1 | `FINDINGS` | 5/5 historical items closed + R29 + R30 |
| ④ closing #2 | **PASS (explicit GO)** | — |
| ④ delta (P5) | **PASS** | with a falsifiability strengthening, adopted |
| ④ delta (P7) | **PASS** | with an evidence-retention suggestion, adopted |

- **Landing check on review conclusions (§3.11)**: every `file:line` from the three review stages was located and re-checked before landing; none had to be downgraded to a lead (no such items this slice).

### 9.4 ⑤ fixes made in place (found only on the live rig; each is in design note R31–R35 / P5–P8)

| # | Symptom | Fix | Falsification |
| --- | --- | --- | --- |
| P1 | a stale `VBoxHeadless` session after a hard power-off made the next `startvm` fail with `E_FAIL` | `Wait-B3bVmPowerOff` + `Clear-B3bStaleVmSession` (word-boundary match, fail-closed) + bounded retry | live `startvm ok=True attempts=1` |
| P2 | the rig had no concurrency protection (two instances could drive the same VM) | session lock (pid + owner + `sessionDir` fingerprint; blocks only a live pid in the same directory) | two assertions (blocks / copied session not blocked) |
| P3 | the frozen writer rejected empty strings ⇒ an uncaught exception after the cut | `Write-B3bArtifactText` (`[AllowEmptyString]`) at every write site that can legitimately be empty | live round exit 0 |
| P4 | the matrix dropped per-round stderr for unexpected exit codes | record round stderr (this patch is what located P3) | — |
| P5 | the control arm's name contains a hyphen while `-Cell` validation rejected it ⇒ the calibration cell was **unschedulable from the CLI** | shared resolver `Resolve-B3bCellSpec` + behavioural assertions | reverting the pattern reddens the assertion |
| P6 | the session id defaulted to "today" ⇒ a `-Finalize` missing `-SessionId` wrote an **empty verdict with exit 0** | finalize fails closed on zero rounds; the self-test threads an explicit fixture session id | disabling the guard reddens the assertion |
| P7 | the inventory contract was implemented more strictly than intended ("first two reads" instead of "the enumeration is settled") | `Get-B3bSettledInventoryPair` (first *consecutive* identical pair, bound 4) + every read archived + the settled pair recorded | reverting to first-pair-only reddens the assertion |
| P8 | a single failed inventory read discarded the whole round and re-ran it | retry inside the round + exit code and stderr logged | three static assertions |
| P9 | the variant-arm filename ended in `_arm`, which Go read as an implicit `GOARCH=arm`, so the file **never built** | renamed to `…_degraded_test.go`; ownership verified via `IgnoredGoFiles` | `go list` in both tag views; the trap is recorded in repository memory |

### 9.5 Residual sampling (§3.11)

**The mandatory-sampling trigger applies**: this slice touches **stopping semantics** (hard power-off injection, reboot and inventory) and device-ownership semantics (the session lock).
Under §3.11 the rule would be one Codex blind audit **before** ④. **How it was actually satisfied**: ④'s **first call already ran with no ③/③.5 checklists attached** (the §3.10 anchoring isolation), which is a blind audit in fact, so no second blind audit was added.
The merge decision is recorded here as §3.11 requires: the two calls' inputs (slice definition + diff + contract sentences + read-only constraints) are structurally identical, so repeating the call would add cost without adding discriminating power.
**Findings from the blind audit**: 4 items (C1–C4); **did it overturn the ③.5 PASS**: no (the two finding sets are disjoint and were both closed by the subsequent closing loop).

## 10. Explicitly not executed

- **Not committed, not pushed** (needs explicit authorisation in the same turn; `origin` only, never gitee).
- **No real model calls / paid probes** (this slice is fully offline).
- **No independent check of device fidelity** (no external byte-level comparison of the virtual disk before and after the cut), so the §1 qualifier stands.
- **No long-form stress runs** (e.g. 200-round simultaneity races): this slice used a fixed N=5 per cell plus 3 rounds per U4 cell, as designed.
- The variant arm is a **manual arm**: it is not part of the 103-round budget and does not enter any arm's three-tier verdict.

## 11. Next-slice recommendation (§7.4)

> Per the dependency graph the next slice should be **B4 · AT-23 E2E evidence** (or, first, a contract-level slice deciding whether B3b's `proven` is adopted).
> Recommended tiers — owner **Deep (Max)** (for the design and conclusion-writing steps; implementation may return to Standard),
> V4 Pro **①③ mandatory** (B4 touches stopping/identity/ownership semantics), ③.5 **enabled** (behaviour-changing slice), Codex **Extra High (④ mandatory, plus the §3.11 blind audit)**.
> Tell me when to start and I will run the 7-step protocol with the closure discipline (V4 Pro pre-analysis first, protocol first).

- **Directive item to be discussed (recorded per ruling D2; the directive is not amended in this round)**: should hard-gate slices (A7/B2/B3/B4) have a **separate** ④ call cap?
  This slice used 6 calls (1 final + 1 re-review + 2 closing + 2 delta), four of which came from the "re-run the stage after remediation" closure rule.
  Suggested evaluation before the next slice starts: either raise the cap explicitly for hard-gate slices (e.g. ≤4), or state in the directive that closure re-runs do **not** consume the cap.

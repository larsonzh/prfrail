# T027 · A7 · Windows publication-durability candidates and falsification — validation report

Date: 2026-09-16. Status: `COMPLETE` (③ pre-review / ③.5 scan / ④ final review are backfilled in §9 as the slice proceeds).
Verdict (**the level and its premise must be read together**): under premise **P = process-level crash, no power loss, single host Windows 11 / NTFS on the local volume** — the **visibility and arbitration layer is `proven`** (P1–P5, constructive argument plus measurement); **power-loss durability is `unresolved`** (U1–U5, including the durability of candidates C2/C3 which still awaits B3 injection - **not the same as "the candidate is disproven"**); the **falsified claim set D1–D6 is `disproven` throughout**; and **Windows publication durability stays `unproven` with first-dispatch still refused**. Every summary sentence, the ADR and the ledgers must carry the same premise and the same three levels; **`proven` must never be read as a power-loss guarantee**. A7 also carries one measurement that **overturns the assumption this slice started from**: a directory-handle flush is *not* unavailable on Windows (§4 E4), so the ADR-013 candidate ranking had to be revised.

## 0. Honesty rules (the premise of every conclusion below)

1. **A process kill leaves the complete OS page cache in place, not a post-power-loss disk state** — a process ended by TerminateProcess loses none of its finished writes;
2. **the `record-linked` and `parent-synced` crash points are observationally equivalent under a kill** (§4 E2 measured the same distribution), so no green experiment in A7 can ever be evidence for directory-entry durability;
3. A7 therefore does not change Windows `unproven` and does not unlock `first-dispatch`; every power-loss claim is recorded `unresolved` and handed to B3 (real crash/power-loss injection).

## 1. Scope and boundaries

- **Done**: ADR candidate enumeration, the four-row falsification matrix, a minimal falsification prototype (in-suite model experiments + `a7native` real-process experiments), a claim ledger.
- **Not done** (hard gates): Windows `unproven` unchanged; no real dispatch; **the production publish primitive is untouched** (`writeReplayRecordNoReplace`'s visibility/durability steps keep their semantics); no CONTRACTS change; no schema/fixture/CI-workflow change; no real candidate, no network.
- The only production-file change is one **nil-default test seam** `replayStorePublishStageHook` (§6): a timing-based kill cannot reliably hit a chosen stage, and the crash-point experiments must.

## 2. Candidate protocols and their A7 outcome

| Candidate | Mechanism | Visibility claim | Durability claim | A7 outcome |
|---|---|---|---|---|
| C1 current baseline | `temp → file.Sync → close → os.Link → parent-durability step (a no-op on Windows) → remove temp` | atomic no-replace (already implemented) | **none** (the Windows step is a no-op; the gate refuses the write) | visibility `proven` (§4 E2/E3); durability `unresolved` |
| C2 `MoveFileEx(MOVEFILE_WRITE_THROUGH)` | replace `os.Link` with a WRITE_THROUGH no-replace move | the full content appears atomically | one documented sentence says the call returns only after the file really moved to disk, another limits the guarantee to copy-and-delete moves | **availability: measured (an evidence annotation)** (succeeds unprivileged, EEXIST maps to `fs.ErrExist`, refuses while a reader holds the target, §4 E6); the flush guarantee stays `unresolved` |
| C3 directory-handle `FlushFileBuffers` | `CreateFile(dir, GENERIC_READ\|GENERIC_WRITE, FILE_FLAG_BACKUP_SEMANTICS)` then flush | — (durability step) | claimed to be the equivalent of a dir fsync | **A7 overturns "unavailable"**: a read-only handle gives `Access is denied`, while **read-write and write-only handles return success without privileges** (§4 E4) ⇒ C3 joins C2 as a B3 injection candidate; whether NTFS really makes the entry durable stays `unresolved` |
| C4 volume-handle flush | open `\\.\X:` and flush | — | semantically covers directory entries | **Not adopted**: needs administrator/backup privilege, which conflicts with the portable product position (ADR-009); E5 was not executed (an unprivileged run would not be representative, registered as such) |
| C5 two-phase marker | publish one more directory entry beside R | one more name | claims marker visibility identifies the commit point | **`disproven` as a durability mechanism** (§4 E7): a marker is a new name, i.e. a new directory entry, so it inherits the same durability class; it adds no state and removes no cross-file consistency surface |
| C6 pre-created slot file | bootstrap a fixed slot; publish by writing inside it plus a file-level flush | length/flag inside the slot | converts directory-entry durability into **file-data** durability (which does have a standard guarantee) | recorded in the ADR as a later candidate; **not pluggable under the frozen contract** (it conflicts with one-file-per-requestId plus no-replace, and arbitration would need a different lock primitive) |
| Excluded | `ReplaceFile` (violates no-replace monotonic publication), `FILE_FLAG_WRITE_THROUGH` on the temp handle (covers file data only), TxF (deprecated), USN-journal reading (its own flush semantics apply; B3-only optional aid) | | | see row D in §3 |

**Symmetry statement**: the Unix-side `proven` is likewise a claim about a *local Unix filesystem class* (tmpfs and network filesystems do not guarantee `directory.Sync` either). Every durability conclusion must be stated over the `(OS, filesystem, storage stack)` triple; A7's host evidence covers **Windows 11 + NTFS on the local volume + this host's storage stack** only.

## 3. Falsification matrix (the four items the slice defines)

| Row | Claim | Falsifiable in A7? | Experiment | Criterion | Result |
|---|---|---|---|---|---|
| A | torn / stale writes: a visible record is either absent or complete and never changes | **partly** (the visibility layer yes; power-loss tearing no) | E2 kill matrix (5 stages × N rounds, a reader polls the published path) | red = partial content, decode failure or valid→different; green = every round inside the allowed set | **green**: the first three stages are always `absent`, the last two always `dispatched-unknown` (§4 E2) |
| B | R and marker cross-file consistency: after a crash the state ∈ {absent, R-only, R+C}, no orphan C, and a marker introduces no state needing new rules | **yes** (process level); directory entries **losing order** under power loss is not → B3 | **E9 completion-publication crash matrix** (same primitive, same injection) + E7 marker model | red = an orphan or an illegal state; green = always in the allowed set and marker intermediate states isomorphic to no-marker | **green** (E9): a pre-link kill leaves R alive with state `dispatched-unknown` (C invisible, no orphan); a post-link kill gives `terminal-receipt-present`; E7: with a visible marker and no receipt the state is still `dispatched-unknown`; **power-loss orphans stay with B3** |
| C | multi-process fencing / reuse: exactly one winner; killing the winner neither wedges the slot nor unlocks a second dispatch; identical hashes replay idempotently and a different hash conflicts | **yes** | E3 (8 in-process publishers plus killing the winner) + **E10 two processes publishing the same requestId in the same round** (the barrier only synchronises the start, so **simultaneity is not enforced** — see §4 E10) | red = two winners, an unreadable slot or a non-idempotent replay; green = invariants hold every round | **green**: E3 gives `first=1 / unknown-block=7` and a killed winner still `dispatched-unknown` with the next call `unknown-block`; **E10 gives exactly one first-dispatch and one unknown-block per round** with a readable store (§4) |
| D | provability per crash point (meta row) | **CP4 and CP5 are observationally equivalent under a kill** ⇒ whether the durability step took effect **cannot be falsified**; this is why A7 cannot promote any candidate to `proven` | the per-stage distribution comparison in E2 | red = the two distributions differ; green = they match | **green (equivalence holds)**: both stages are `dispatched-unknown` × N ⇒ the durability claim is `unresolved`, left to B3 |

## 4. Experiment log (host Windows 11 / NTFS; `a7native`, real processes)

> The `×N` values at the end of the table are the **sampling counts of this report**; the gate run uses `PROOFRAIL_A7_ROUNDS=20` (the request matrix, the completion matrix and two-process fencing are all green). §8 reconciles the two conventions.

| Experiment | Claim | Observation | Result |
|---|---|---|---|
| E2 kill matrix | after a kill at any publish stage the replay root holds only allowed outcomes | `temp-written`/`temp-synced`/`temp-closed` → `absent`×5; `record-linked`/`parent-synced` → `dispatched-unknown`×5 | green; CP4≡CP5 |
| E3 arbitration and killing the winner | concurrent single winner; a killed winner does not unlock dispatch | 8 concurrent → `first=1, unknown-block=7`; after killing the winner the state is `dispatched-unknown` and the next call returns `unknown-block` | green |
| E4 directory-handle flush probe | (initial assumption) a directory flush is unavailable on Windows | read-only handle → `Access is denied`; **read-write handle → `<nil>`; write-only handle → `<nil>`** | **assumption overturned**: callable unprivileged; whether it really reaches stable storage is `unresolved` (→ B3) |
| E5 volume-handle flush | whether it works with privileges | **not executed** (needs an administrator; conflicts with the portable position and would not be representative) | registered as not executed; adoption criterion in §2 C4 |
| E6 `MoveFileEx` probe | C2 availability, EEXIST semantics, open-reader behaviour | publishes unprivileged; an existing target gives `fs.ErrExist` (convergence semantics survive); **an existing target that nobody holds open is refused as well**; a target held open gives `ERROR_ALREADY_EXISTS`; mapping to `fs.ErrExist` is **hard-asserted** (mismatch fails the test) | green (availability: measured, an evidence annotation); **the open variable is not isolated**: the refusal is caused by the target existing, not by it being open, and this report claims no cause it did not measure |
| E7 marker model (in-suite) | a marker cannot lift the durability class and adds no state | the marker goes through the same durability seam (`syncCalls>0`); `PublicationDurability()` is unchanged; the state is still `dispatched-unknown`; the Windows declaration stays `unproven` and the write gate still refuses | `disproven` (marker as a durability mechanism) |
| E8 temp residue | whether residue from a pre-link kill is inert | 1 leftover `.tmp`; not treated as a record and does not block a later publication; **no startup cleanup exists, so it accumulates** | green + finding (§7) |
| E9 completion-publication crash matrix (real kills, **matched to the `completions/` path only**) | a kill at any stage of the completion publication must not produce an orphan and must not lose the request | `temp-written`/`temp-synced`/`temp-closed` → `dispatched-unknown`×5 (C invisible, R alive); `record-linked`/`parent-synced` → `terminal-receipt-present`×5 | green (no orphan) |
| E10 two-process fencing (real processes, same round) | two processes racing the same requestId | exactly one `first-dispatch` plus one `unknown-block` per round; the store stays readable with state `dispatched-unknown` | green; **two-way rendezvous fencing: each child writes its own ready file, and the parent releases them once after seeing both**, giving exactly one first-dispatch in each of 5 rounds (the concurrency semantics are unchanged: `unknown-block`, not `unknown-unlock`) |

## 5. Claim ledger

**Level semantics (identical to the slice's acceptance vocabulary)**: every conclusion level is **exactly one of `proven` / `disproven` / `unresolved`**; `observed`, `falsified` and `disproven-by-policy` are **not levels** and are used only as **evidence annotations** (e.g. "availability: measured", "overturned by measurement", "excluded by policy"). `proven` = supported by a constructive argument plus reproducible measurement **under an explicitly stated premise**. Every P1–P5 `proven` below carries premise P (process-level crash, no power loss, single host Windows 11 / NTFS on the local volume) and says nothing about power-loss durability.

**Constructive argument behind P1–P5**:

- **P1**: `os.Link` (CreateHardLink on Windows) atomically creates a new name in the same directory and **fails when the target exists** (`fs.ErrExist` passed through), so "visible" and "exclusive" are two faces of one atomic step.
- **P2**: the record is written to a temp file in the same directory, synced and closed before the link, so **visible implies complete**; a failed write or sync discards the temp, leaving no path on which half-written content can become visible (E2: zero tears over the rounds; sampling in §4/§8 - 5 rounds in this report, 20 in the gate run).
- **P3**: R is published before C (`RecordCompletion` calls `ensureRequestLocked` first) and C uses the same no-replace primitive, so a crash cannot produce the one ordering gap, "C without R" (E9: zero orphans over the rounds; sampling in §4/§8 - 5 rounds in this report, 20 in the gate run).
- **P4**: a no-replace collision enters bounded reread convergence (≤16 attempts) and compares digests, so identical content replays idempotently and different content yields `ErrAgentRunnerRequestConflict` with no write (E3/E10 plus the existing suite).
- **P5**: `PublicationDurability() != Proven` is refused at **all 8 write-path entries** (request, completion, launch receipt, launch identity, terminal intent, terminal closure, and the publisher's publication and recovery), and no written content can change the platform declaration (E7).

| Claim | Level | Evidence |
|---|---|---|
| P1 atomic no-replace visibility (**8-way true concurrency inside one process** plus **same-round concurrent attempts across processes**; cross-process strong synchronisation is not proved) | `proven` (constructive argument + measurement, under premise P) | argument above + E2/E3/E10 |
| P2 no torn visibility on a published path | `proven` (constructive argument + measurement) | argument above + E2 (zero tears over N rounds) |
| P3 after a process-level crash the state ∈ {absent, R-only, R+C} with no orphan | `proven` (constructive argument + measurement) | argument above + E2/E7/E9 |
| P4 single-winner arbitration, idempotent replay, different-hash conflict | `proven` (constructive argument + measurement) | argument above + E3/E10 + the existing suite |
| P5 the unproven gate is faithful (it refuses and deleting the guard reddens a test) | `proven` (constructive argument + measurement) | argument above + the two existing platform/dispatcher tests + the A7 re-check |
| D1 "a directory-handle flush is unavailable on Windows" | `disproven` (evidence annotation: overturned by measurement) | E4 (read-write and write-only handles return success, **every writable mode's success being a hard assertion**) |
| D2 "a volume flush is usable for the portable, unprivileged product" | `disproven` (evidence annotation: **excluded by policy, not an experimental refutation**) | needs administrator privilege (established knowledge) plus ADR-009; **E5 not executed** |
| D3 "a marker lifts the durability class" | `disproven` | E7 |
| D4 "a pre-created slot is pluggable under the frozen contract" | `disproven` | contract conflict (§2 C6) |
| D5 "`FILE_FLAG_WRITE_THROUGH` covers directory entries" | `disproven` | semantics (file data only) + the E6 observation |
| D6 "`ReplaceFile` is compatible with no-replace" | `disproven` | contract (a published file must not be replaced) |
| U1 power-loss durability of candidates C2/C3 | `unresolved` | **not provable in A7**; B3 real crash/power-loss injection |
| U2 C2's ambiguous documentation ("moved to disk" vs "copy-and-delete moves are flushed") | `unresolved` | the vendor text contradicts itself; A7 has no authority to adjudicate |
| U3 local storage-stack (device cache/firmware) flush semantics | `unresolved` | not observable in A7 |
| U4 directory entries **losing order** under power loss (a real orphan) | `unresolved` | entirely out of A7's reach; B3's minimum experiment is in §7 |
| U5 shared-volume (SMB/cluster) fencing | `unresolved` | outside the product scope |

**Falsifiability (what would overturn this)**: any kill round showing torn visibility, valid→different, or an orphan C overturns P1/P2/P3; a guard whose deletion does not redden a test overturns P5; a primitive that holds for a marker entry but not for the record would overturn D3 (theoretically impossible); only B3's real power-loss evidence may overturn U1 — no A7 experiment may stand in for it.

## 6. Changes and test points (with mutation bindings)

**Changes**:
1. `internal/adapters/agent_runner_replay_store.go`: a nil-default package-level seam `replayStorePublishStageHook` plus five stage constants and `replayStorePublishStage()` call sites (`temp-written`/`temp-synced`/`temp-closed`/`record-linked`/`parent-synced`). **Zero semantic change while nil.**
2. New `internal/adapters/agent_runner_replay_store_a7_test.go` (in-suite, runs in CI).
3. New `internal/adapters/agent_runner_native_a7_test.go` (`//go:build a7native`).
3.5 New `internal/adapters/agent_runner_native_a7_fencing_test.go` (`//go:build a7native`; the E9 completion matrix and the E10 two-process fencing).
4. New `internal/adapters/agent_runner_native_a7_windows_test.go` (`//go:build a7native && windows`).
5. `docs/ADR_REGISTER.md` / `_EN.md`: the ADR-013 row plus a §5 summary.
6. This report and its `_EN` counterpart.
7. Wrap-up: `docs/DEV_PLAN{,_EN}.md` and `docs/t027/REMAINING_SLICES{,_EN}.md` status write-back.

**Test points and mutation bindings**:

| Test | Guard it protects | Mutation | Expectation |
|---|---|---|---|
| `TestA7PublishStageHookIsNilByDefaultAndSeesEveryStage` | the seam must be nil in production and all five stage call sites must exist | delete any stage call site | red |
| `TestA7CrashAtEveryStageLeavesOnlyLegalReplayStates` | the visibility boundary (before/after the link) and a closed state set | make the record visible before the link / widen the state set | red |
| `TestA7MarkerPublicationAddsNoDurabilityClassAndNoNewState` | only the platform declaration may lift a durability class; a marker creates no state | let the marker bypass `replayStoreSyncParentDirectory` / flip the Windows declaration to proven | red (Windows leg) |
| `TestA7LeftoverTempArtifactsAreInvisibleAndDoNotBlockPublication` | the loader reads exact paths and never scans the directory | make the loader pick up `.tmp` files | red |
| `TestA7ReplayIdempotenceAndConflictGuardsStillBind` | the digest comparison in `ensureRequestLocked` | delete the digest comparison | red |
| `TestA7NativePublishCrashMatrix` | stage semantics and the allowed set under a real kill | decouple the stage decision from visibility | red |
| `TestA7NativeSingleWinnerAndKilledWinnerDoNotUnlockDispatch` | single winner + a killed winner not unlocking | treat the killed winner's slot as re-dispatchable | red |
| `TestA7NativeCompletionPublicationCrashMatrix` | completion (C) publication is separate from request (R) publication: a kill inside `completions/` must create no orphan and lose no R | drop the `completions/` path match (falling back to the R re-publication point) / make C visible before the link | red |
| `TestA7NativeConcurrentProcessFencing` | same-round two-process concurrency with exactly one first-dispatch (two-way rendezvous) | turn the two-way rendezvous back into a one-way start gate (a serial start would then pass, i.e. a false green) | red |
| `TestA7NativeTempResidueIsInertButAccumulates` | `.tmp` residue is completely inert for visibility and arbitration | make the loader scan the directory and pick up `.tmp` files | red |
| `TestA7NativeDirectoryHandleFlushProbe` | the access-mode dependency: **every writable mode must succeed** and read-only must be denied | widen the assertion to "at least one writable mode succeeded" / delete the read-only denial assertion | red |
| `TestA7NativeMoveFileExWriteThroughProbe` | no-replace refuses to overwrite | allow the overwrite | red |

## 7. Known boundaries and the hand-off to B3

- **U1/U2/U3/U4 all move to B3**; B3's minimum experiment: real power loss (or an equivalent device-level injection) + a per-round crash-point journal (written to a **different physical device** that can itself prove durability) + a restart inventory of the R/C pairs and directory entries; inject per candidate (C2, C3); a candidate may only be written into the contract as `proven` after B3 is fully green.
- **`a7native` does not run in CI**: the `a7native` experiments (including E4/E6/E9/E10) run **neither in the CI default path nor in the `-race` subset** (the same treatment as `a6native`, to avoid flakiness); the conclusions here are **single-host** evidence and must be re-run manually on a Windows host (`$env:PROOFRAIL_A7_ROUNDS=20; go test -tags a7native -count=1 -run TestA7Native ./internal/adapters/`).
- **Temp residue accumulates** (E8 finding): `.tmp` files left by a pre-link crash are never cleaned. This slice **does not fix it** (that is a production behaviour change needing its own slice decision); the recommendation is for B3 or an operations slice to evaluate whether startup cleanup introduces new read-side semantics.
- **E5 not executed**: a volume flush needs privilege and conflicts with the portable product position (ADR-009); a change of that position would require a new decision.
- **Scope of every claim**: durability conclusions cover `Windows 11 + NTFS on the local volume + this host's storage stack` only; the Unix-side `proven` holds over local filesystem classes only.
- **A7 changes nothing about**: Windows `unproven`, the `first-dispatch` refusal, the frozen CONTRACTS sentences, T027's `BLOCKED / NOT IMPLEMENTED` state, or the AT-23 conclusion.

## 8. Cost and metering

- Probe budget: **0 of 10 used** (A7 is fully offline; no paid probes).
- Model calls: ① V4 Pro pre-analysis ×1; ③ V4 Pro pre-review ×1 (`PASS WITH FIXES`) plus re-review ×3 (round 1 `RE-REVIEW: FINDINGS`: 1 Medium + 4 Low; round 2: 3 Low; round 3: `RE-REVIEW: PASS`); ③.5 MAI-Code-1.1-Flash ×2 (`INDEPENDENT SCAN: FINDINGS` → `PASS`); ④ Codex ×6 rounds (rounds 1–5 `FINDINGS`, per-round counts in the tables below, **1 High + 10 Medium + 5 Low** in total, all fixed) then **round 6 `RE-REVIEW: PASS`**; **§9 is the single authority for each round's findings and their integration**, so `×N` always equals the number of recorded rounds); one further V4 Pro call failed with a provider connection timeout (no probe budget consumed; registered under the §3.9 availability pre-check).
- Experiment rounds: 5 in-suite A7 tests; **7** `a7native` tests; each matrix runs `5 stages × rounds` (5 sampled in this report, 20 in the gate run); two-process fencing runs `rounds × 2 child processes` on the same convention.

## 9. Review record (③ / ③.5 / ④)

### ③ V4 Pro implementation pre-review (2026-09-16)

Verdict: **`PASS WITH FIXES`** (no High finding; no production-semantics regression; the frozen contract sentences are not violated).

| Level | Finding | Fix |
|---|---|---|
| Medium | M1 row B cited evidence that does not exist (E2 never publishes C) and P3's "no orphan" had no C-crash evidence | **Added E9**, the completion-publication crash matrix (real kills, 5 stages × N rounds, zero orphans), and pointed row B at E9 |
| Medium | M2 the header/register said `proven` while the ledger said `observed` | The ledger's P1–P5 were raised to `proven` with a **constructive-argument section** and the premise "process-level crash, no power loss" (§5) |
| Medium | M3 E6 did not isolate the open-reader variable (an existing target is refused regardless) | The probe **gained a control** (an existing target that nobody holds open is refused too) and the report now says the cause is the target existing and that the probe does not isolate openness (§4 E6) |
| Medium | M4 multi-process fencing had no experiment (row C was an in-process 8-goroutine race plus a post-hoc kill) | **Added E10**, the two-process fencing experiment (two child processes publishing the same requestId, exactly one first-dispatch per round) |
| Low | L1 `a7ReleaseLeakedTempHandles` failed silently (Windows cleanup could flake and the mutation stayed green) | It now calls `t.Fatal` when the poll budget is exhausted |
| Low | L2 the directory-flush probe was green-on-inversion (`flushAccepted==0` only logged) | It now fails directly (the file does not run in CI, so a Windows host must redden) |
| Low | L3 the seam contract was undeclared (runs under the caller's lock, must not call back into the store, covers every no-replace publication, failure paths fire no stage) | The seam comment now states that contract |
| Low | L4 `_ = errors.Is` was a placeholder | Removed, together with the import |
| Low | L5 D2's `disproven` label overreached (E5 was not executed) | Relabelled `disproven-by-policy` (not an experimental refutation) with E5 recorded as not executed (**converged again in ④ round 1 to `disproven` plus the annotation "excluded by policy"**) |

Re-review (③ re-run after the fixes): **`RE-REVIEW: FINDINGS`** (no High; 1 Medium + 4 Low).

| Level | Re-review finding | Fix |
|---|---|---|
| Medium | **the first three stages of E9 actually killed inside the request re-publication** (`ensureRequestLocked` still writes a temp for an existing request and only collides at the link), so "a crash before C" was never injected | The seam became `func(stage, path string)` (still nil-default, still semantics-free) and the completion helper **matches the `completions/` path only**; after the change the first three stages are genuine completion crash points |
| Low | the `_ = errors.Is` placeholder remained | Removed, together with its import |
| Low | the §5 P5 argument counted "four call sites" (there are 8) | Changed to "all 8 write-path entries" and named them |
| Low | §4's counters (×5/×3) could not be reconciled with the gate's `ROUNDS=20` | §4 gained a sampling note and §8 the unified convention |
| Low | E10's "simultaneous publication" was not enforced by the experiment (the barrier is a starting gate) | Wording narrowed to "same round", with §3 row C and §4 E10 recording that simultaneity is not enforced |

Re-review (③ re-run, second attempt): **`RE-REVIEW: FINDINGS`** (no High, no Medium; 3 Low, all document-wording or ledger consistency).

| Level | Finding | Fix |
|---|---|---|
| Low | §3 row C still said "at the same instant", contradicting §4 E10 and this record | Row C now says two processes publishing **in the same round**, with the note that the barrier only synchronises the start so simultaneity is not enforced |
| Low | the Chinese §3 carried a **duplicate row C** (the older E3-only row was never deleted), so the two languages disagreed structurally | The old row was deleted, leaving only the E10 row (and row B's typo "跨文件" was restored) |
| Low | `docs/t027/REMAINING_SLICES{,_EN}.md` still said A7 was "not started", contradicting the actual progress | The status row and the section heading now read in progress (the completion write-back still happens at ⑥) |

Re-review (third attempt): **`RE-REVIEW: PASS`** (no High/Medium/Low; L-A/L-B/L-C all closed and no new problem found). ③ ran as pre-review `PASS WITH FIXES` → re-review ×2 `FINDINGS` (1 Medium + 4 Low → 3 Low) → re-review 3 `PASS`.

Inherent limits of ③ (recorded under the Section E convention): the re-reviews were **read-only static reviews** (no terminal access in that environment; full-directory static diagnostics were zero), so the gate and `a7native` numbers rest on the host logs and were not independently re-run there; E4/E6/E9/E10 are Windows single-host evidence (§7).

### ③.5 MAI low-cost independent scan (2026-09-16)

Round 1: **`INDEPENDENT SCAN: FINDINGS`** (2 High + 1 Medium, no Critical; all about evidence level and wording, none about code semantics).

| Level | Finding | Fix |
|---|---|---|
| High | a top-level summary sentence omitted the premise and could be read as "A7 overall = proven" | The header now says **the level and its premise must be read together**, states premise P (process-level crash, no power loss, single host Windows 11 / NTFS on the local volume) and the three levels; the ADR-013 row and its §5 were aligned (never read `proven` as a power-loss guarantee) |
| High | P1's "including concurrency" claimed more than E10 measured (the barrier is only a starting gate) | P1 now reads **8-way true concurrency inside one process plus same-round concurrent attempts across processes; cross-process strong synchronisation is not proved** (same wording in §3 row C and §4 E10) |
| Medium | the level vocabulary could drift between the ADR, the ledger and the English version | The ADR-013 row gained a "levels and premise" sentence (`proven` under premise P; candidate-level results all `disproven`; `unresolved`) aligned with the report header and §5 |

Round 2: **`INDEPENDENT SCAN: PASS`** (all three closed, nothing new; a PASS about **document wording and evidence level**, not about product readiness).

### ④ Codex independent final review (2026-09-16, first round with no list attached)

Round 1: **`RE-REVIEW: FINDINGS`** (1 High + 2 Medium + 1 Low; no Critical; no production-semantics regression and the frozen contract sentences held).

| Level | Finding | Fix |
|---|---|---|
| High | the slice's acceptance vocabulary allows only `proven`/`disproven`/`unresolved`, yet the ledger and ADR introduced `observed`, `falsified` and `disproven-by-policy` as **peer levels** | **Levels converged to the three states**: the §5 header now says every level is one of the three and that `observed`/`falsified`/`disproven-by-policy` are evidence annotations only; D1/D2 rows became `disproven` plus annotations; §4 E6's "availability `observed`" became "availability: measured (annotation)"; the ADR-013 row and §5 were aligned in both languages |
| Medium | E10's children were released by a **one-way start gate**, so a serial start could still pass | **Two-way rendezvous**: each child writes its own ready file and parks before publishing; the parent waits for **both** ready files, writes the single `go` once and only then lets them race; re-ran 5 rounds with exactly one first-dispatch each |
| Medium | `MoveFileEx`'s "existing target → `fs.ErrExist`" was logged, never asserted | Both that branch and the open-reader branch now **assert** (`errors.Is(err, fs.ErrExist)`, `t.Fatal` otherwise) |
| Low | the request matrix compared CP4/CP5 with `fmt.Sprint(map)`, which is order-dependent | Replaced with a **per-key comparison** of the two states (`absent`, `dispatched-unknown`) |

Re-review (④ round 2): **`RE-REVIEW: FINDINGS`** (2 Medium + 1 Low; no High/Critical; all four round-1 fixes were verified in place, and the new findings are about **document-wording consistency** and **audit-count closure**).

| Level | Finding | Fix |
|---|---|---|
| Medium | in `DEV_PLAN{,_EN}.md` the sentence about Windows publication durability carried neither premise P nor the three-state vocabulary, so the vocabulary could drift across documents | That sentence now states **the level and its premise must be read together** (every level is one of `proven`/`disproven`/`unresolved`, under premise P = …), and the dev plan gained the **A7 paragraph** (premise P, three levels, boundaries and the retained `unproven`; the two languages are structurally identical) |
| Medium | §8's counts and §9's record disagreed (§8 said "one re-review" while §9 recorded three) and ④ round 2 was still "to be backfilled" | §8 now states the actual call counts, and each ④ round is backfilled **immediately after it returns** (this round included) |
| Low | E6's control case (an existing target that nobody holds open) only asserted "non-nil", not the error class | The control now adds a **hard assertion** `errors.Is(closedErr, fs.ErrExist)`, so "the open reader changes nothing about the error class" rests on an assertion rather than an inference |

Re-review (④ round 3): **`RE-REVIEW: FINDINGS`** (2 Medium + 2 Low; no High/Critical; all three round-2 fixes were confirmed in place, and the new findings centre on **round-count closure** and **candidate disposition across documents**).

| Level | Finding | Fix |
|---|---|---|
| Medium | §8 said Codex had run three rounds while §9 still said "round 3 to be backfilled", and the DEV_PLAN A7 paragraph said "two rounds", so three places disagreed | **Single source of truth**: the per-round detail lives in report §9 (§8 records only the total and points there) and the DEV_PLAN no longer enumerates rounds itself; each round is backfilled into §9 as soon as it returns |
| Medium | the DEV_PLAN A7 paragraph's "candidates C1–C6 are `disproven` throughout" contradicted the report and ADR, where C2/C3 stay `unresolved` for power-loss durability and are B3 injection candidates | Rewritten to match the report: **D1–D6 are the set of falsified *claims***; C2/C3 remain `unresolved` on power-loss durability and go to B3 (C1 is the shipped baseline; the C5 marker is `disproven` as a durability mechanism; C6 is pending a contract revision) |
| Low | §6's traceability was incomplete: it declared 7 `a7native` tests but bound only 4, and the change list omitted the fencing file | §6 now lists `agent_runner_native_a7_fencing_test.go` and the mutation bindings for the three missing experiments (completion matrix, two-process fencing, temp residue); the E4 binding now reads "widen to 'at least one writable mode succeeded'" |
| Low | E4's wording ("read-write and write-only both succeed") claimed more than the test asserted ("at least one succeeded" plus "read-only denied") | **The test was strengthened**: every writable access mode must succeed (`writableAccepted == writableModes`) while read-only stays denied |

Re-review (④ round 4): **`RE-REVIEW: FINDINGS`** (3 Medium + 1 Low; no High/Critical; all four round-3 fixes were confirmed in place, and the new findings are about **candidate disposition**, **the evidence denominator** and **count consistency**).

| Level | Finding | Fix |
|---|---|---|
| Medium | the ADR-013 row and ADR §5 still said "candidate-level results are `disproven`", contradicting both the C2/C3-are-B3-candidates sentence in the same section and the new DEV_PLAN wording, so a reader could take C2/C3 as already-disproven candidates | Both ADR places (both languages) now use the same vocabulary as the report and dev plan: **D1–D6 are the falsified claim set**, and **C2/C3 remain `unresolved` on power-loss durability as B3 injection candidates, which is not the same as "the candidate is disproven"** |
| Medium | the E4 probe used "writable modes that happened to open" as its denominator, so a writable mode whose `CreateFile` failed was excluded from the denominator and the test could still pass - the report's "every writable mode's success is a hard assertion" claimed more than the test could do | The denominator is now **fixed at 2** (`const a7WritableModes = 2`): a `CreateFile` failure on any writable mode calls `t.Fatal`, every writable mode must flush successfully (`writableFlushed != 2` fails), and the read-only mode stays denied |
| Medium | §8 declared four rounds while §9 still carried a round-4 placeholder, so closure did not hold at that instant | This section now holds the complete round-4 record (verdict + findings + fixes), and §8 no longer announces a round before it is backfilled (each round is backfilled on return, so `×N` always equals the number of recorded rounds) |
| Low | §5's constructive argument said "E2: three rounds, zero tearing / E9: three rounds, zero orphans" while §4 sampled ×5 and §8 says 5 rounds in this report | Now reads "sampling in §4/§8: 5 rounds in this report, 20 in the gate run", identically in both languages |

Re-review (④ round 5): **`RE-REVIEW: FINDINGS`** (1 Medium; no High/Critical; all four round-4 fixes were confirmed in place and Q1/Q3-Q11 are satisfied - only the **candidate-disposition wording in the header** had not converged).

| Level | Finding | Fix |
|---|---|---|
| Medium | the header still said "several candidates are `disproven` (D1–D6)", phrasing D1–D6 as candidate results; since the header is the highest-priority entry point a reader could take C2/C3 as already-disproven candidates | Now reads "the **falsified claim set D1–D6 is `disproven` throughout**" and keeps the anti-misreading anchor in the same sentence: **the durability of candidates C2/C3 still awaits B3 injection, which is not the same as "the candidate is disproven"**; mirrored in both languages |

Re-review (④ round 6, final): **`RE-REVIEW: PASS`** (zero Medium or higher; all of Q1–Q11 pass with a file location cited for each).

The re-review explicitly verified (with locations): the three-state vocabulary (§5 and the header); the same premise P and the same candidate-disposition vocabulary in the header, ledger, ADR row, ADR §5 and dev plan (D1–D6 the falsified claim set; C2/C3 still `unresolved` and handed to B3); the frozen contract intact with a semantics-free seam (`replay_store.go`'s seam plus the unproven gates at every write entry, and the frozen CONTRACTS section), evidence reconcilable with conclusions (§3 A–D / §4 E2–E10 / §5 P-D-U; the E4 denominator fixed at 2), E10's two-way rendezvous with per-round unique paths and stdout decisions, E9's `completions/`-only matching, E6's three `errors.Is(..., fs.ErrExist)` hard assertions claiming no unmeasured causality, the per-key CP4≡CP5 comparison, §6's changes and mutation bindings matching the seven experiments, complete honesty boundaries, and count closure (§8 `×N` = the rounds recorded in §9; the dev plan only points at §9).

Residual Low observations (the reviewer's own words; known boundaries, not new problems): the review was static and did not independently re-run `a7native` or the full gate set (no terminal); `a7native` evidence remains Windows single-host and out of CI, so power-loss injection in B3 and cross-environment validation are still required.

MAI's **minimal anti-misreading boundary sentence** (adopted): **"the same requestId on one host, one volume and one round is proved to have a single winner that does not unlock a second dispatch; cross-process strong synchronisation, cross-host strong synchronisation and power-loss durability remain unproved."**

**MAI round-2 Section D (invariants still without test coverage)**: power-loss durability (real power loss, device caches, firmware write-back, directory-entry ordering); cross-process **strong** synchronisation (E10 proves only a single winner per round on one host and volume); non-local storage stacks (SMB, shared volumes, cloud drives, clusters); `.tmp` residue cleanup. **Section E (not verifiable)**: real power-loss injection and device-layer semantics (B3's territory). **What would overturn it**: another premise-free `proven` summary, P1/E10 being widened back to "cross-process strong synchronisation is proved", or a demand to lift `unproven` before CONTRACTS is revised. MAI also confirmed the wording did **not** over-narrow (`proven` stays under premise P, and P1 keeps 8-way true in-process concurrency plus same-round cross-process attempts).

### Residual sampling (§3.11 criterion)

**Not triggered**: A7's only production change is a nil-default test seam (`replayStorePublishStageHook`, zero semantic change) and it alters neither ownership, stopping nor identity semantics - the seam covers those publications without changing their behaviour - so this slice adds no §3.11 mandatory Codex blind audit. ④ still runs as usual and carries "the seam's coverage and its zero-semantics claim" as a completeness requirement.

# ProofRail P37 report: T027-B production-grade architecture design (directive v1.58)

Slice `P37-T027B-BREAKTHROUGH` (assessment-type `[SLICE]`); executed 2026-10-30; directive v1.58; both rounds ran on the read-only carrier `independent-reviewer` and were transcribed verbatim by the control seat; R = 0, D = 0.

## 1 Executive summary

- **Call facts**: two rounds ran continuously, not in batches. Round 1 used `GPT-6.1 Sol (copilot)`, round 2 used `GPT-6 Astra (copilot)` (the off-list authorisation is OB-130 in §12.4). Two calls in total, within the cap of 3.
- **Transcription and checks**: both `review.md` files were transcribed verbatim and their `sha256` matches the normalised original message (only BOM and line-ending normalisation may differ); round 1 is 111 lines and 10723 bytes, round 2 is 68 lines and 9277 bytes.
- **Batch assistance**: two dispatches, both `deep-reasoner` + `GLM-5.3-Flash (glm)`; category 1/2 (input-pack statistics and number checks) and category 2 (two-round transcription number checks), each mechanically re-verified by the control seat.
- **Core judgement**: round 1 delivers a complete design (G0-G4 gate order, module split, state machine, key interfaces, resource constraints); round 2 exposes three fatal hand-over gaps (start ownership, evidence durability, gate-order prerequisites) and gives a minimal rebuild plus the November minimal boundary. **Together they are the workable baseline.**

## 2 Two rounds in review

### 2.1 Round 1 - architecture design (Sol)

- **Gate order**: G0 project approval (OB-55 ownership and frozen items) then G1 durability (per platform-times-premise) then G2 boundary (restricted candidate runs) then G3 assembly (zero-model three tasks) then G4 real loop. Gates bind evidence; a date never stands in for a condition.
- **Modules**: `chain` alone owns scheduling and acceptance; `guard` owns the managed process and stopping; adapters only translate external facts; SessionBridge stays off the main execution chain.
- **State machine**: completed is allowed only for a PASSED step; anything unknown lands in uncertain / paused, and completed never expresses the unknown.
- **Interfaces**: reuse the admission / dispatch / postflight signatures and add `Prepare` / `WaitDispatched` / `Decode` as internal interfaces, each requiring idempotence on the same binding and preserving the underlying failure class.
- **Resources**: pinned candidate version, one watchdog, persistent budget with unknown hold, minimal credential injection, network denied by default.

### 2.2 Round 2 - critique and rebuild (Astra)

- **Three fatal gaps**: (1) start-ownership hand-over (an identity-publication failure after spawn can be reported as "never spawned", while a failed reclaim releases the slot); (2) durable completion records do not make the referenced evidence originals durable; (3) G1 hides production prerequisites outside this boundary and G2 already needs real calls, so the original gate order cannot be a November commitment.
- **Six minimal rebuilds**: separate ownership from short mutexes, make the Prepare result a non-reusable managed resource, bind safety decisions to the opened object, add an evidence commit barrier, select an immutable context by request digest, and wire persistent reservation and settlement confirmation to the call and publish boundaries.
- **Three downgrades**: split "development order" from "production admission order"; narrow generic Wait re-entry to terminal-transaction recovery; demote full session resume to a later capability.
- **November minimal boundary**: one Windows premise P, one candidate, one task with one Agent step plus declared hooks; verify assembly and serial acceptance with a zero-model stand-in first; the real three-task loop needs its own conditional acceptance, and stand-in evidence never impersonates real-candidate evidence.

## 3 Input pack and measurement

- **Round 1 input**: 62 items (6 document items, 56 Go files), 14923 lines, 712662 bytes; `internal/adapters` is 38 files and 395919 bytes, `internal/chain` is 18 files and 184724 bytes. The per-item list is `01-sol-design/input-list.md`.
- **Round 2 input**: the full round-1 design book plus the same original pack, re-delivered item by item to avoid context drift.
- **Model calls**: 2 (cap 3). **Batch assistance**: 2, which the §2.2 note excludes from that cap.
- **Conventions**: whole files count raw bytes including BOM; line ranges are sliced, re-joined with LF and counted without the original BOM.

## 4 Awaiting a user ruling

1. Whether OB-55 belongs to a standalone T027-B implementation slice, and that slice boundary.
2. Which candidate and which Windows premise P are frozen; whether same-boundary durability and real probe evidence are authorised.
3. The AT-23 three-unit mapping and who signs it; whether the first acceptance covers session resume.
4. Resource, cost and human-intervention caps, plus the stop or postpone condition if the pilot is blocked.
5. Whether to accept "ship the bounded mechanism and assembly slice first, with the real three-task loop on separate conditional acceptance".

## 5 Pipeline trimming and stage declaration

- ⑥ skipped: this slice is an assessment slice with no implementation diff.
- ⑦ not applicable: this slice has no implementation review object; round 2 Astra only critiques architecture across models and produces no ⑦ verdict.
- ⑧a transcription: both rounds were transcribed verbatim by the control seat and the `sha256` matches the normalised original.
- ⑨: the encoding gate plus three scopes plus `SELFTEST` are recorded in the ⑩ receipt.
- ⑧c: the temp directory `tmp/P37-T027B-BREAKTHROUGH/` is cleaned with the three zero-residue read-backs.
- ⑩: stop point, awaiting commit / push authorisation.

## 6 Artefact registration

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/t027b-breakthrough.md	repo
docs/validation/t027b-breakthrough_EN.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/00-doc-inventory.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/freedom-list.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/MANIFEST.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/SHA256SUMS.txt	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/01-sol-design/input-list.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/01-sol-design/review.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/input-list.md	repo
docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/review.md	repo
```

## 7 Exception handling / fact pack

1. **Tier**: the user set 1M / Max on the client side; the control seat cannot observe it and records the fact as-is. The transcript offset is evidence that a call happened: the round-1 dispatch is around L98894-L98898 and the round-2 dispatch around L99179-L99183, with the transcript at 99502 lines at record time; that offset does **not** claim to prove the tier.
2. **Batch assistance**: 2 dispatches (categories 1/2 and 2), each mechanically re-verified; no ranking, Top-N selection or conclusive work was dispatched.
3. **At ⑨**: any red gate was reworded only inside the authorised boundary and re-run; the process and result are in the ⑩ receipt, and no criterion was bypassed.
4. **Nothing else**: no out-of-scope call, no unreported temp artifact, no contract or code change.

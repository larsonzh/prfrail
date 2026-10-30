# ProofRail P38 report: T027-B code-level rebuild blueprint (directive v1.59)

Slice `P38-T027B-TACTICAL-BLUEPRINT` (assessment type `[SLICE]`); run on 2026-10-30 through 2026-10-31; directive v1.59; both rounds were executed by read-only carriers and transcribed verbatim by the control seat; R = 0, D = 0.

## 1 Executive summary

- **Call facts**: round 1 (Sol) was split into two batches by user ruling (batch A = C1-C4, batch B = C5-C7) and round 2 (Astra) is the red-team acceptance; 3 usable bodies, plus 1 voided first dispatch.
- **Transcription and checks**: all three bodies were transcribed verbatim and `sha256` matches between source and landed copy; batch A is 10,052 characters, batch B 9,808 and round 2 2,362.
- **Batch assistance**: 2 dispatches (input-pack statistics, evidence-pack consistency recomputation), both `deep-reasoner` + `GLM-5.3-Flash (glm)`, each mechanically re-verified.
- **Core judgement**: the blueprint covers C1-C7 with file- and function-level changes, new tests and assertion points for each; Astra returns `VERDICT: REVISE` - 3 BLOCKED (C1 / C4 / C6) and 4 PARTIAL (C2 / C3 / C5 / C7), plus new holes N1-N4.

## 2 Two rounds in review

### 2.1 Round 1 - blueprint (Sol, batched)

- **Batch A (C1-C4)**: C1 replaces the unbounded reclaim wait with one reaper per process plus bounded termination; C2 adds `ProcessSpec.OnSpawn` and a pre-registered starting owner so a spawned process can no longer be retried as never-started; C3 binds the bounded read and the PID read to one open object; C4 reorders the terminal commit so artifact sealing precedes intent.
- **Batch B (C5-C7)**: C5 selects an immutable context by the digest `K=(run,task,step,attempt,requestId,requestHash)` and shares `outcomeOnce`; C6 adds the durable wrapper `DurableCostLedger` with `Reserve`, `RequireOutstandingReservation` and `Settle`; C7 uses `EvaluateT027BReadiness` to separate the probe, the real acceptance and the production release, removing the AT-23 self-dependency.
- **Cross-item constraints**: each batch states one invariant set (resident ownership, single-writer core, no nested locks, sealing then settlement then closure ordering).

### 2.2 Round 2 - red-team acceptance (Astra)

- **Verdicts**: C1 / C4 / C6 are BLOCKED; C2 / C3 / C5 / C7 are PARTIAL, retaining respectively the host-crash window, the directory-replacement window, the hard-coded attempt and the reusable probe quota.
- **New holes**: N1 is the host-crash window before Job ownership exists; N2 is ancestor replacement after the directory helper returns; N3 is the hard-coded attempt in `Reopen`; N4 is the reusable `ProbePermit` quota.
- **Over-design**: no necessary over-design was found; the advice is to bound `ArtifactCommitter` and `DurableCostLedger` to this cycle fixed flows and proven backends, adding no generic transaction framework.
- **Contract blocker**: `CONTRACT-BLOCK` - upgrading Windows publication durability requires revising `docs/CONTRACTS.md` section 7 with an independent review first.

## 3 Input pack and measurement

- **Input pack**: 59 items (3 document segments + 56 `.go`), line count 14,714, 648,707 bytes (about 634 KB); per-item detail in `01-sol-blueprint/input-list.md`.
- **Round 2 narrowing**: only the two blueprint bodies, no code pack (user instruction, to save inference time); declared as narrowed in this report per the v20 rule.
- **Output measurement**: batch A source 13,442 bytes, batch B source 14,583 bytes, round 2 source 4,242 bytes, all measured as BOM-free text.

## 4 Open items for the user

| # | Item | State | Blocks |
|---|---|---|---|
| 1 | Whether N1-N4 from the REVISE verdict get one more in-slice round | Open | November implementation |
| 2 | Whether the `CONTRACT-BLOCK` becomes a contract revision slice | Open | C4 and C6 landing |
| 3 | Whether PARTIAL for C2 / C3 / C5 / C7 is acceptable as the November start point | Open | T027-B implementation |

## 5 Pipeline trimming and process declarations

- **Trim basis**: this slice is an assessment type `[SLICE]` with no implementation diff and no implementation review object.
- **Batch assistance**: 2 dispatches, covering step-0 input-pack statistics and evidence-pack consistency recomputation; no ranking, selection or conclusive work was dispatched.
- **Declaration basis**: the two lines below are the G7-f skip channel; this report carries no `sixstage` fence.

- ⑥ skipped: this slice is an assessment slice with no implementation diff.
- ⑦ not applicable: this slice has no implementation review object; round 2 Astra only performs a red-team review and produces no ⑦ verdict.

## 6 Artifact registration

The change set holds 15 paths: the directive pair (header bumped to v1.59, the section 0.3 log row, OB-131 in section 12.4), the ledger pair (the P38 row and the P38 block), the report pair and 9 evidence-pack paths.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/t027b-tactical-blueprint.md	repo
docs/validation/t027b-tactical-blueprint_EN.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/00-doc-inventory.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/freedom-list.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/MANIFEST.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/SHA256SUMS.txt	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/input-list.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/input-list.md	repo
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/review.md	repo
```

## 7 Exception handling / fact pack

1. **Tier**: the user set 1M / Max on the client side; the control seat cannot observe it and records the fact as-is. The transcript offset is evidence that a call happened: the batch-A dispatch is around L101607, batch B around L101855 and round 2 around L102004, with the transcript at 102,123 lines at record time; that offset does **not** claim to prove the tier.
2. **Output-length resistance**: the first round-1 dispatch triggered self-restart and truncation twice; the 42,509-byte and 15,428-byte originals are kept as evidence but were not used as bodies; after batching, batch A still self-restarted once and only the first complete pass was taken, the second being voided.
3. **Stop-loss**: per the user ruling, retries stopped as soon as batch A first returned complete, consuming no further quota; the round-2 input was narrowed to the two blueprint bodies by the same ruling.
4. **Nothing else**: no out-of-scope call, no unreported temp artifact, no contract or code change.

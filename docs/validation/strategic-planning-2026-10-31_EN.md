# ProofRail P36 forward-planning report: three rounds and the synthesis ruling (directive v1.57)

Slice `STRATEGIC-PLANNING-2026-10-31` (assessment-type `[SLICE]`); **marker date 2026-10-31, actually executed 2026-10-30**; directive v1.57; all three rounds ran on the read-only carrier `independent-reviewer` and were transcribed verbatim by the control seat; R = 1 (the ⑨ gate flagged a new CJK character, the wording was changed and the gate re-run, counted as one rework), D = 0.

## 1 Executive summary

- **Call facts**: three rounds ran continuously, not in batches. Rounds 1 and 2 used `GPT-6.1 Sol (copilot)` (Max), round 3 used `GPT-6 Astra (copilot)` (Max; the off-list authorisation is OB-129 in §12.4, effective 2026-10-30). Three calls in total, within the cap of 10.
- **Transcription and checks**: all three `review.md` files were transcribed verbatim and their `sha256` matches the normalised original message (only BOM and line-ending normalisation may differ); logic lines are 60, 54 and 56, bytes including BOM are 5581, 5715 and 7154.
- **CI evidence (cross-slice support statement)**: the "both CI runs green" in the P35 ledger completion row is supported by two records — run `40832429061` (commit `397dc30`) in §1 of the P35 report, and run `40833256125` (headSha `70cb193`, both legs `success`) recorded in this slice pack at `00-doc-inventory.md`. **This report cites no run id that the P35 report does not carry.**
- **Core judgement**: the right order for November is "first open the legal execution entry and the real no-IDE loop, then make delivery recoverable"; on the governance side, subtraction frees capacity but **never lowers gate strength**.

## 2 Findings by round

### 2.1 Round 1 - product and architecture planning

- The S1 MVA should be defined under "single Windows user + pinned CLI + no IDE"; the items proposed for removal from S1 (the full TUI, two generations of self-hosting, the black-box sub-items) need approval with the requirements and acceptance updated, and **until approved S1 is still judged incomplete**.
- Platform: Windows first with Linux kept on core CI; the truly irreversible binding is not "Windows first" but writing drive letters, path folding, proprietary credentials or process semantics into persistent protocol, hashes and core state.
- Language: freeze the Python / Node official extensions and the formal support promise and **register it explicitly as a narrowing of the S2 commitment**; do not freeze the existing tool scripts or the generic capability.
- De-IDE in three steps: legal entry, then the real loop, then substitutability, each with a failure criterion; `TUI` / `Web` stay views only.

### 2.2 Round 2 - governance and process planning

- Three governance subtractions: stop auto-promoting every observation into a governance slice; cut duplicated manual rewriting (call counts and R / D summarised in one place); cut re-reads that carry no new information. **Not reducible**: contract-first, independent acceptance, the ⑤⑥⑦ enabled by the matrix, gates with raw evidence, and the permission and cost constraints.
- T027-B: conditionally projectable, but a project is not a production unlock; first rule on OB-55 ownership and freeze the candidate, target environment, AT-23 mapping, resources and probe authorisation; run November as four "evidence by deadline" boxes, and **never let a date stand in for a condition**.
- Advisory AI review node: a one-off consultation on a whitelisted model and a read-only carrier; mechanically it gets no state, approval or release write access, and the formal control flow never depends on its output.

### 2.3 Round 3 - synthesis ruling

- The first two rounds are not mutually exclusive; they differ in tier and prerequisites. Round 1 silence on ⑥⑦, Astra and the checklist **must not be read as support for an exemption**.
- Rulings: this round carries no ⑥⑦; a consultation never substitutes for a formal review; T027-B may be prepared for project approval at once but must not start coding on that basis; Windows durability is recorded per platform-times-premise and never upgraded platform-wide.
- On the eight checklist questions from round 2: three hold and five partly hold; the minimal disposition for each holding item stays inside the checklist, the dispatch and the evidence, and **touches no directive body**.

## 3 Synthesis ruling (November Top 3 and the first step)

1. **Legal execution entry**: pin the CLI version and the Windows deployment premise P; isolation and durability evidence must be reproducible, discriminating and mechanically executed; `unsupported` / `unproven` refuses the first dispatch.
2. **Real no-IDE loop**: once a real probe is authorised, close VS Code and verify the three real task success chains, T2 failure not starting T3, no re-execution after restart, budget and handoff re-checks, independent acceptance and a zero-write source tree.
3. **Recoverable delivery**: fixed ZIP, offline verification, backup and restore, retention evidence on decommission, and the cost plus human intervention recorded for every accepted task.

The three run in series on one resource pool and a later item never pre-empts an earlier one; **completing all three by month end is not promised**.

**First step (within 7 days)**: submit and rule on the pre-project admission pack for T027-B (OB-55 ownership, implementation boundary, pinned candidate, Windows / P, AT-23 and refusal-path mapping, evidence gaps, resource and probe cost caps). The ruling may only be "limited implementation allowed" or "blocked with release conditions", and **it grants no production admission**.

## 4 Awaiting a user ruling

1. OB-55 ownership, and whether T027-B is projected on the frozen boundary.
2. Whether the round-1 MVA, the S1 removals and the S2 narrowing are approved.
3. Candidate version, Windows / P, resource and real-probe cost authorisation, plus the stop or postpone plan if the pilot is blocked.
4. The schedule and budget of the independent formal stages of P36 as determined by the current matrix.
5. The acceptance boundary between AT-24 `unknown` and "unknown side effects stay paused".

## 5 Pipeline trimming and stage declaration

- ⑥ skipped: this slice is an assessment slice with no implementation diff.
- ⑦ not applicable: this slice has no implementation review object; round 3 Astra only rules on planning across models and produces no ⑦ verdict.
- ⑧a transcription: all three rounds were transcribed verbatim by the control seat and the `sha256` matches the normalised original.
- ⑨: the encoding gate plus three scopes plus `SELFTEST` are recorded in the ⑩ receipt.
- ⑧c: the temp directory `tmp/STRATEGIC-PLANNING-2026-10-31/` is cleaned with the three zero-residue read-backs.
- ⑩: stop point, awaiting commit / push authorisation.

## 6 Cost and measurement

- **Model calls**: two Sol rounds plus one Astra round, three calls in total, all at Max. The user confirms this round consumed 705.7 credits (68% to 71%) with no overrun risk.
- **Batch assistance**: two dispatches were actually made, both `deep-reasoner` + `GLM-5.3-Flash (glm)`, and every scenario stayed inside the four §2.2 categories (enumeration and statistics; excerpt plus number checks; bilingual mirror language check; manifest draft). **No ranking or Top-N selection was dispatched, and no conclusive work was dispatched.**
- **Why rounds 2 and 3 switched to main-agent script self-proof (kept truthfully)**: the transcription check of either round is deterministic inside one control-seat command (the script prints both `sha256` values), so an extra dispatch would add noise without information; under the reverse rule "stop what does not meet the bar, and never lower the bar to exhaust a resource pack" the verification moved to self-proof, cross-checked independently by the ⑨ gate and `PACK-CONSISTENT`.
- **Batch reading versus the control-seat recount**: the second dispatch reported the report pair as logic lines 79, heading lines 9, bold markers 38, bars 0 and backticks 48, identical item by item to the control-seat recount; the only difference is the blank-line convention (the control seat counts 21 under the `split` convention including the trailing empty segment, the batch assistant 20 excluding it), and the equality between the two sides is unaffected.

## 7 Exception handling

1. **Spelling fix outside the original boundary**: `GPT-6.1-Sol (copilot)` corrected to `GPT-6.1 Sol (copilot)` (CN and EN in step, the `modelId` column `gpt-6.1-sol` untouched); **this exceeds the original boundary, is a §2.2 spelling fix and is a direct user ruling in the same round on 2026-10-30**; the basis is that all three rounds of this slice accepted the space form and that it matches OB-62.
2. **One rework at ⑨ (R = 1)**: the first run flagged one new CJK character by `G4b`, the wording was changed and the gate re-ran green; the report records `R` as 1 truthfully.
3. **Temp-file line-ending difference**: the landed copy of the v25 checklist was written with CRLF by the editor tool and normalised to LF at transcription per §11.3, an allowed difference; the `sha256` comparison is equal under normalisation.
4. **Nothing else**: no unreported temp artifact, no out-of-scope call, no gate bypass.

## 8 Artefact registration

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/strategic-planning-2026-10-31.md	repo
docs/validation/strategic-planning-2026-10-31_EN.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/00-doc-inventory.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/freedom-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/MANIFEST.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/SHA256SUMS.txt	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/01-product-planning/input-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/01-product-planning/review.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/02-governance-planning/input-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/02-governance-planning/review.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/03-synthesis-astra/input-list.md	repo
docs/validation/evidence/STRATEGIC-PLANNING-2026-10-31/03-synthesis-astra/review.md	repo
```

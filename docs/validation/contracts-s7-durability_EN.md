# P54 `CONTRACTS-S7-DURABILITY` validation report (step 0 skeleton)

Slice: Windows durability contract upgrade and C6 landing (size L · contract & implementation `[SLICE]`). Anchor HEAD = `2b8ac3e`, directive v1.60. The R-02 ruling approved the slice on 2026-11-04.

> Status: step 0 skeleton; stage 1 architecture (double play) → 3 red test → 2 implementation → 4 gates → 5/6/7 reviews → 8 docs → 9 verification → 8c → 10 stop point are appended per round.

## 1. Exception registration

- ① Start-message goal contradiction (structural defect) — location: the MoveFileEx example in start-message section six item 1; the contradiction: the title and the goal sentence ("land C6") and the example mechanism (the B3b C2 MoveFileEx) belong to two different candidate numberings; ruling basis: the user ruling of 2026-11-04 "the MoveFileEx example in the start-message body is treated as an example of how to write it and does not change the goal of landing C6; the title and the goal sentence take precedence"; disposition: the goal mechanism of this slice = the ADR-013 C6 (the pre-created fixed slot), and the C6/C2 split of this table is kept.
- ② G7-f ⑥ skip declaration (the former next line was the G7-f fixed-form line) — status: already replaced by the 8a `sixstage` fence (the ⑥ text is transcribed verbatim in §6, ending with a closed-set conclusion line), and the transitional placeholder is deleted; allocated stage-6 quota = 3; this slice revises CONTRACTS §7 durability-grade semantics, so §7.2.1 item 3 does not hold and stage 6 is mandatory.
- ③ This temporary stop point (stage-1 completion → design review → release the red team) — location: after stage 1 completes, before the Sol red team; ruling holder: the user; basis: the user ruling of 2026-11-04 "this stop point is explicitly set by the user for this slice and takes effect only once; it is not repeated after the red team is released; it does not replace stage 10 and does not change the regular rhythm of start-message section seven"; effective only once, not generalised; the ruling text is archived verbatim in `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/ruling-2026-11-04.md`.
- ④ Temporary-authorisation registration gap (disclosed truthfully per the OB-134 precedent) — instances: the P53 stage-1 double-play Sol and this slice's first Sol red-team pass; registered: no independent §12.4/ADR registration (the independent trace does not hold); nature of the gap: §1.7.4 item 2 requires an OB or ADR independent registration; disposition: disclosed truthfully, no retroactive backfill. The user ruled on 2026-11-04 to count by the truth source (only OB-132 in this class), making this slice's second red team the 2nd and authorisable.
- ⑤ Governance-side input (kept per the user suggestion of 2026-11-04): the truth-source count of this temporary-authorisation class reaches 3 in P54, triggering §1.7.4 ⑤ (iii); later slices must not use a temporary authorisation again and must open a dedicated slice to bring the outside-whitelist red-team model into §2.2 or a standalone clause.
- ⑥ Landing-form deviation registration (2026-11-04): the red-team closure self-check table landed as report §5 (a standalone section), deviating from the earlier "into report §3" wording; reason = a standalone section is clearer and does not affect compliance.
- ⑦ tmp-root misplaced-file self-heal (§5.1 self-heal tier, 2026-11-04): the controller found `tmp\s7a-contrasts.js` at the ⑩ final-state reading (misplaced at the tmp root, most likely a leftover comparison script of the ②e execution role; gitignored, hence invisible to `git status`); order: self-report → record the evidence-② reading (253) → exact delete with `Remove-Item -LiteralPath` → zero-residual evidence (path `Test-Path`=False; `??` list empty; recursive file count 252 ≤ 253; TEMP hits 0); D (§5.1 temp-file-discipline calibre) = 1; because the violator attribution is uncertain (controller or ②e), the self-heal judgment is submitted, per the stricter "controller-involved" calibre, to the same-round independent review layer for confirmation — the independent review layer (MAI-Code-1.1-Flash) returned `CONFIRMED` (original in the evidence package `07-final/review.md` item three).

## 2. artifacts fence

The full change-set entries of this slice (evidence payloads, report pair, ledgers, spec docs and code/test changes; path + disposition, TAB-separated, path-sorted). Step 0 declares only the payloads currently on disk; later stages append per round.

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/CHECKLIST_v43.md	repo
docs/t027/CHECKLIST_v43_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/contracts-s7-durability.md	repo
docs/validation/contracts-s7-durability_EN.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/00-briefing.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2b.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2c.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review-r2.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/ruling-2026-11-04.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/ruling-2026-11-04-q1.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/05-prereview/review-r1.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/06-scan/review.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/07-final/review.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/cross-platform-linux.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/gate-readings.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/green-c6slot-tests.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/08-gates/red-baseline-c6slot-compile.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/MANIFEST.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/SHA256SUMS.txt	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/freedom-list.md	repo
docs/validation/evidence/CONTRACTS-S7-DURABILITY/input-list.md	repo
internal/adapters/agent_runner_replay_store.go	repo
internal/adapters/agent_runner_replay_store_c6_enable.go	repo
internal/adapters/agent_runner_replay_store_c6_noop_test.go	repo
internal/adapters/agent_runner_replay_store_c6_windows.go	repo
internal/adapters/agent_runner_replay_store_c6_windows_create.go	repo
internal/adapters/agent_runner_replay_store_c6_windows_test.go	repo
internal/adapters/agent_runner_replay_store_launches.go	repo
internal/adapters/agent_runner_replay_store_terminals.go	repo
internal/console/runtime_agent.go	repo
```

## 3. Step 0 record

- Anchor check: HEAD = `2b8ac3e`, origin/main = `2b8ac3e`, directive v1.60, `SELFTEST 697/697`, three scopes `TOTAL_FAIL=0`; all five items match (the generational-gap disclaimer is verified: a backfill commit advanced HEAD).
- The v43 section-10 opener "v42" → "v43" is fixed (CN+EN mirrored).
- The ledger is registered (status table `🔄 进行中` + section, R/D=0 baseline).
- Pre-check: the CONTRACTS §7 three-grade durability wording, the B3b verdict ledger and the A7 C6 definition with the D4 contract conflict were read; the key open question (the ADR-013 C6 numbering vs the B3b C2 `MoveFileEx` relation) is registered in `input-list.md` and handed to stage 1.

## 4. C6 upgrade judgment table (step 0 stop point → ruled)

Judged 2026-11-04; every basis row carries file:line, no inference. **Overall initial judgment = user ruling needed; user ruling (2026-11-04) = Q1 A (the goal mechanism of this slice = the ADR-013 C6 pre-created fixed slot) / Q2 B (the Windows publication durability grade stays `unproven` and first-dispatch keeps being rejected)**; the Windows platform grade must not be written `proven` in the P54 contract revision; the C2 candidate-level conclusion already lives in the B3b verdict ledger and P54 does not rewrite it; the ruling text is in §1 item ③ and the archived file.

| Row | Field | Judgment |
|---|---|---|
| C6 (the name used by this slice, ADR-013 numbering) | Mechanism description | bootstrap pre-creates a fixed slot; publish = write the in-slot data + a file-level flush, turning directory-entry durability into file-data durability (which has standard guarantees); not pluggable under the current contract (conflicts with "an independent file per requestId + no-replace", arbitration needs a different lock primitive). Sources: `docs/validation/t027-windows-publication-durability.md:28`, `docs/ADR_REGISTER.md:23` |
| C6 | B3b injection-matrix coverage | Not covered. Basis: `docs/validation/t027-b3b-durability.md:44-48` (the candidate table lists only C1/C2/C3/C2C3/positive-control); `:79` (the main matrix 17 cells x 5 = 85 rounds has no c6 cell); `:116-128` (the §4.1 verdict ledger has no C6 arm); `docs/t027/REMAINING_SLICES.md:443` (the B3b ledger entry states the same) |
| C6 | (a) reproducible injection or control evidence | Missing (no C6 injection round exists) |
| C6 | (b) discrimination gate | Missing (no C6 negative/positive control) |
| C6 | (c) mechanism evidence | Missing (no C6 mechanism-token record) |
| C6 | Initial judgment | User ruling needed (all three criteria missing; `proven` must not be written) |
| C2 (the example mechanism in the start message, not the name used by this slice) | Mechanism description | a `MoveFileExW` + `MOVEFILE_WRITE_THROUGH (0x8)` no-replace move replacing `os.Link`. Source: `docs/validation/t027-b3b-durability.md:45` |
| C2 | B3b injection-matrix coverage | Covered. Basis: `docs/validation/t027-b3b-durability.md:122` (§4.1 C2 candidate-level `proven` (only under premise P′; whether P′ covers the physical-machine target deployment environment is not argued), decisionSets=[record-linked,parent-synced] both cells fully eligible) |
| C2 | (a) reproducible injection or control evidence | Complete (the main matrix 85 rounds all valid; the two C2 decision cells 5/5 survived and 0 lost; `:79`, `:90-91`) |
| C2 | (b) discrimination gate | Complete (discrimination gate ok=true: negative control 5/5 lost, positive control 5/5 survived; `:128`) |
| C2 | (c) mechanism evidence | Complete (both decision cells 5/5 mechanism tokens; the variant arm proves the rig refuses survivals whose mechanism never ran; `:90-91`, `:104-107`) |
| C2 | Initial judgment | Upgradeable to proven (only under premise P′; whether P′ covers the physical-machine target deployment environment is not argued) |

Note: the proven of the C2 row of this table is candidate-level and inside premise P′, not platform-level; whether P′ (VirtualBox guest / 25H2 / hostiocache=false / ~1KiB / N=5) covers the physical-machine target deployment environment is not argued. The parenthetical must not be truncated in citations.

Sub-note (appended per the user ruling of 2026-11-04, the original text unchanged): (full qualification: only under premise P′; whether P′ covers the physical-machine target deployment environment is not argued).

## 5. Red-team closure self-check table (before stage 2, read-only)

This table is the controller's self-check, not an independent review; ⑤/⑥ must independently confirm against both verbatim red-team rounds.

The controller checked the r2b closure quality item by item against both verbatim red-team rounds (2026-11-04; read-only, no file changed).

| Red-team finding | Grade | r2b closure location | Self-check |
|---|---|---|---|
| Round-1 High 1 abnormal paths (empty/corrupt slot confusion, lost-slot rebuild, visibility vs flush conflation) | High | r2 A2 uncommitted-state freeze and check order, three time points, in-lock truncate retry; the round-2 dependency on prebuild-set existence is closed by the user-ruled degradation | Closed (degradation user-ruled) |
| Round-1 High 2 seam scope vs reader coverage mismatch | High | r2 A2 seam scope (R/C-only `publishReplayRCNoReplace`, ownership/launch/terminal keep the original path byte-for-byte); round 2: no remaining finding, closed at design level | Closed |
| Round-1 High 3 non-equivalent prebuild-lifecycle reading | High | r2 A2 rewritten per Q1=B plus "this ruling overrules r1's High 6 judgment"; round 2: no remaining finding, closed | Closed |
| Round-1 M4 core-claim argument gap | Medium | r2 D conditional claim + E/F three separate measurements; not re-opened in round 2 | Closed |
| Round-1 M5 discrimination measurement object undefined | Medium | r2 F three measurements (physical slot survival / uncommitted payload survival / logical record visibility); not re-opened | Closed |
| Round-1 M6 LockFileEx arbitration protocol insufficient | Medium | r2 A2 arbitration clauses (scope/same handle/bounded wait/lock order/unlock/mapped-view exclusion); not re-opened | Closed |
| Round-1 M3 citation discipline not fully propagated | Medium | full parentheticals added everywhere + the round-2 note Low closed by the user-ruled sub-note (report §4) | Closed |
| Round-1 Sol items ①②③ (all need revision) | — | ① Q1=B overrules the r1 reading; ② narrowed seam; ③ arbitration clauses | Closed |
| Round-2 High historical full-set binding gap | High | r2b change-point ① verbatim degradation sentence + "this degradation is accepted per the user ruling of 2026-11-04" + the ADR-029 supplement | Closed |
| Round-2 Medium allowed-set inconsistency | Medium | r2b annex C allowed-set table (process-kill vs power-loss recovery, self-consistent per cut point) + the table note | Closed |
| Round-2 Medium BA3 replay-probe error | Medium | r2b supplement B·BA3 two-arm comparison rewrite | Closed |
| Round-2 Low ordering argument over-broad | Low | r2b supplement A-4 limited to "first prebuild + current single-stage record format" | Closed |
| Round-2 Low own-hash wording conflict | Low | r2b supplement A·A-1 wording revision | Closed |
| Round-2 Low API-shape vs execution-path conflation | Low | r2b supplement B three-layer split, instruction-level guarantees removed | Closed |
| Round-2 Low judgment-table note missing the parenthetical | Low | report §4 sub-note (original text untouched, user-ruled) | Closed |

Conclusion: every finding = Closed, none partially closed or unclosed; stage 2 may enter.

Mirror status: the EN mirror of this table is already complete in this file (15 rows CN/EN equal), synced as part of ⑧a.

## 6. Review chain

- ① double play: r1 → first Sol red team `FINDINGS` → r2 → second Sol red team `FINDINGS` → r2b closure → r2c activation-semantics micro-fix; both red-team transcript sets are archived as `docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md` and `red-team-sol-review-r2.md`, with the four design rounds in the same directory.
- ②a–d implementation: contract first (§7 C6 Decision Path + Annex C) → narrowed R/C seam → store-level activation.
- ③a–c: red tests first, the BA3 two-arm rewrite, and three gap tests (`TestC6PrebuildNoopWithoutActivation`, `TestC6PrebuildPlanDriftIsCorruption`, `TestC6PrebuildRunIDMismatchIsRootConflict`).
- ④ all gates green: default, c6slot and linux build/vet, plus the tree and index scopes.
- ⑤ pre-review: `PRE-REVIEW: PASS WITH FIXES` (0H/0M/4L, L1–L4 disposition in §7).
- ⑥ independent scan: done (verbatim transcription below).
- ⑦ final review: Codex `FINAL REVIEW: PASS` (0 Critical / 0 High / 0 Medium / 0 Low; the red-team equivalence argument passed at §7.2.1 calibre; the impact-surface judgment for the ⑧a-before-⑦ sequencing deviation holds); original archived verbatim in the evidence package `07-final/review.md`.
- ⑦ re-review (#1/3, narrowed input: the ②e dead-code relocation): Codex `FINAL REVIEW: PASS` (0 findings; pure-relocation, dead-code elimination, c6slot completeness, build-tag and import hygiene all pass); original appended to `07-final/review.md` round two.
- ⑧/⑨/⑩: closed — stage ⑨ all three scopes green + `SELFTEST: PASS 697/697`, stage ⑧c cleanup done, and after user authorisation at ⑩: commit `7c0763a` pushed to origin (never gitee), run `41840557379` green on the first attempt on both legs (windows 15 packages ok / ubuntu 15 packages ok + Race 6 packages ok).

⑥ independent-scan verbatim transcription (headings written as Section A–E to satisfy the gate, body verbatim):
```sixstage
## Section A Scope
- Read-only independent scan of the §7 durability contract, the C6 prebuild activation path, the Windows slot implementation, and the no-op fallback path.
- Scope excludes code edits or any dependency on the earlier ⑤ pre-review conclusion.

## Section B Findings
- No findings.
- The contract wording and the implementation remain aligned: the C6 path is build-gated, store-activated only after a successful prebuild, and the non-activated build remains a no-op.

## Section C Evidence
- The production assembly order is explicit in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L129-L195): `NewAgentRunnerReplayStoreForRun` is followed immediately by `PrebuildReplaySlots(...)`, before any dispatch.
- The activation semantics are defined in [internal/adapters/agent_runner_replay_store_c6_enable.go](internal/adapters/agent_runner_replay_store_c6_enable.go#L13-L115): prebuild success sets the store active, failure leaves the run fail-closed, and reopen verification is idempotent.
- The Windows slot behavior is enforced in [internal/adapters/agent_runner_replay_store_c6_windows.go](internal/adapters/agent_runner_replay_store_c6_windows.go#L15-L239): header format, payloadLen validation order, empty-slot freeze semantics, and no-replace publication rules.
- The unactivated fallback is pinned in [internal/adapters/agent_runner_replay_store_c6_noop_test.go](internal/adapters/agent_runner_replay_store_c6_noop_test.go#L12-L87): a non-`c6slot` build cannot activate the slot path or touch the replay root.
- Fresh validation command: `cd /d d:\LZProjects\prfrail && go test ./... && go test -tags c6slot ./internal/adapters`
  - Result: all packages reported `ok`; the adapters package completed successfully with no failing assertions.

## Section D Falsifiability audit
- The Windows C6 behavior is explicitly testable and tested in [internal/adapters/agent_runner_replay_store_c6_windows_test.go](internal/adapters/agent_runner_replay_store_c6_windows_test.go): empty-header creation, committed-slot rewrite refusal, corruption ordering, uncommitted intermediate state handling, prebuild activation, and fail-closed prebuild conditions.
- Those tests are structured as negative/positive checks and therefore provide a valid falsification surface for the contract.

## Section E Unverified / overturn conditions
- None open for the current implementation.
- Any future platform or premise change would require re-running the same validation set and re-checking the C6 activation enablement and durability claims under that new premise.

INDEPENDENT SCAN: PASS
```

## 7. Deviations and Open Items

- L1 (dangling Annex C reference): closed by embedding Annex C into CONTRACTS{,_EN} §7 and re-pointing the A2 reference to "Annex C of this section".
- L2/L3 (BA1/BA2 unexecuted assertions, no pin on the drift and runId branches): closed by the three ③c tests (names in the §6 ③a–c row).
- L4 (the A2 tail sentence lacks the "limited observation surface" qualifier): left unchanged; rewriting needs user authorisation, so it is registered truthfully as an open item of this slice. The user's 2026-11-04 ruling maintains the original wording (reason: ⑦ already judged `FINAL REVIEW: PASS`; L4 is Low and ⑤ registered it truthfully as "unclosed within the authorisation boundary"; leaving it unchanged does not constitute a false pass), so this item is closed.
- r2c activation-semantics micro-fix: the "build-level vs store-level" conflict exposed by ③; design-r2c chose store-level activation and A2 gained the activation-semantics sentence.
- SKIP=1: `TestSpawnRecoveryHelper` is a subprocess design behaviour (also skipped in the default build), not a failure.
- gopls build-tag warning: appears only in the tool-configuration state of the `c6_enable` file and does not affect build/vet/test.
- Sequencing deviation (deviation registration 1st, not merged into the §5.1 D count; self-discovered by the controller): the ⑧a documentation sync ran before the ⑦ final review (the directive pipeline order is ⑦ then ⑧); cause = a controller orchestration slip; impact surface = the ⑧a products (six documents + ADR-030 + report §6/§7) are included in the ⑦ input package for review, so the ⑦ conclusion covers the final state; disposition = registered truthfully, no retroactive re-ordering.
- ②e dead-code relocation (triggered by the user's 2026-11-04 clarification): the `createReplaySlotWindows` creation chain had no call site in the default Windows build (gopls `unusedfunc`) → moved into the new file `agent_runner_replay_store_c6_windows_create.go` (windows && c6slot); byte-identical semantics, stage ⑨ all three states green; R counts 1 (fix re-run after a finding), R=1 < 2 does not trigger the high-risk row; the narrowed ⑦ re-review is in §6.

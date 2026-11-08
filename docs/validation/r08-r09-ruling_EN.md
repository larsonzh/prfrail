# R-08 / R-09 Ruling Record (AT-23 / AT-24 preconditions)

> Nature: `[ROUTINE]` docs-only landing (no slice number, no directive version bump); ruler = the user; ruling date = 2026-11-06; this document = finalized on 2026-11-07. **CN is authoritative**, `_EN` is the mirror.

## 1. Background and scope

- R-08 / R-09 are the pending-ruling items of §6 in v49; this document is the complete record of the user's 2026-11-06 ruling.
- Ruling baseline: aiming at "getting it runnable in the near term"; BYOK main path + local runs + conditional closure; loose at first, tightened step by step.

## 2. R-08 ruling (6 items)

### R-08.1 Candidate version = (a) GitHub Copilot CLI 1.0.83

- Wrapper functions `Invoke-CopilotCli` / `buildCopilotArgs`; `-Unlock` defaults to `$false` (fail-closed).
- Stage 1: verification before the first run + sampling 1 in every 5 probes; stage 2: verification on every run.
- Three-way evidence comparison (default wrapper / explicit parameters / `-Unlock` bare run).
- Version evidence uses the `copilot --version` output, not `winget list` or `ProductVersion`.
- New copilot call sites must go through the wrapper; otherwise flagged as a deviation in ⑤ pre-review.
- Known boundary: the lock relies on the SEA loader's `Pr()` behavior; auto-update paths outside the loader are not audited.
- Lock failure fallback: switch to option (b) "version not controllable".

### R-08.2 Probe budget

- Stage 1: declared per occurrence (user approves in the same turn), no preset cap on the first occurrence.
- Stage 2: Copilot 5 / BYOK 100 CNY / Copilot currency 0 / retries 5 / additions require same-turn user approval.
- "Declared per occurrence" = same-turn user approval; the controller must not decide on its own.

### R-08.3 Windows premise P = (a) P only (process-level crash, no power loss, local Win11/NTFS)

- The AT-23 report must truthfully state "does not cover power-loss scenarios; power-loss durability cites B3b (P′ covering the physical-machine target deployment environment is not argued)".

### R-08.4 VM / proxy = (a) neither allowed

- BYOK connectivity pre-check failure path: stop and escalate to the user — options (a) temporarily enable the proxy (same-turn approval required) / (b) switch to a locally reachable endpoint / (c) pause AT-23; the controller must not enable the proxy or switch endpoints on its own.

### R-08.5 Real-candidate-leg failure disposition = (b) conditional closure (reduced)

- Mechanism legs all green + the real candidate leg actually ran (≥3 times, covering start/execute/stop) with unsupported/unknown recorded.
- Every unsupported/unknown must carry concrete error information or evidence (raw stderr, exit code, failing stage), never just "cannot run".
- The report must state clearly "AT-23 is conditionally closed and does not constitute proof of real-candidate compatibility".

- **Applicability conditions of the inherited-drift exemption (five elements, all must hold simultaneously)**:
- ① **Same baseline**: the comparison baseline = the B4 baseline archive (`docs/validation/t027-at23-e2e.md` + `docs/validation/evidence/directive-history/b4/archive/`); only this archive may serve as the comparison target — no other baseline may be chosen or mixed in.
- ② **Same assertion**: the assertion id, the assertion meaning and the assertion's harness (`tools/b4-e2e/`) are all unchanged versus the B4 baseline; any harness change or any redefinition of the assertion fails.
- ③ **Drift direction stricter**: the assertion value changes in the more conservative / stricter direction (e.g. `process-tree-taken-away` from `true` to `false` = no process-tree takeover detected = a more conservative verdict); the reverse direction (e.g. `false` to `true` = more lenient) fails.
- ④ **Four states unchanged**: `verdictClass` / `chainState` / `taskState` / `stepState` are verbatim-identical to the B4 baseline; any drift in any of these fields fails.
- ⑤ **Cause traceable**: the drift cause is traceable to a specific product change after B4 (e.g. P61's `02d5797`, the A5 pid-reports binding as the sixth frozen fact); a cause that cannot be attributed to a specific change fails.
- **Consequence of the exemption**: all five elements hold simultaneously ⇒ the "premise ① unmet" verdict is exempted and premise ① is treated as satisfied (the inherited-drift exemption is established); any element fails ⇒ no exemption, and the existing criteria of this subsection apply.
- **Exemption boundary**: the exemption does not equal AT-23 conditional closure — premise ② of the real candidate leg must still hold; the exemption must be adjudicated by **⑦ (independent final review)**, and the controller must not adjudicate it on its own.

### R-08.6 BYOK parallel path

- (a) main path / (a) unconditional / (b) independently billed (does not consume Copilot premium requests, subject to the BYOK cost cap).
- Copilot-hosted is used only for 1 GPT-6 Luna connectivity test and does not enter the acceptance conclusion.
- BYOK configuration: `COPILOT_PROVIDER_TYPE=anthropic` + `BASE_URL=https://api.deepseek.com/anthropic` + `API_KEY` + `MODEL=deepseek-flash` + two token caps (`COPILOT_PROVIDER_MAX_PROMPT_TOKENS=840000` + `COPILOT_PROVIDER_MAX_OUTPUT_TOKENS=128000`).
- **Step 0 mandatory BYOK feasibility minimal verification**: `copilot --no-auto-update -p "say hello"`; on failure stop and escalate to the user — options (a) switch to the openai type / (b) adjust the R-08.2 cap / (c) pause AT-23.
- Register in the report: the BYOK conclusion cannot replace the B1 Copilot-hosted baseline.

## 3. R-09 ruling (9 items)

### R-09.1 Tool-use unknown = (c) case by case

Accepted if the post-return artifact scan proves "artifacts are inside the workspace with no external side effects"; blocked if traces outside the workspace are found.

### R-09.2 Resource-read unknown = (c) case by case

Accepted if the post-return scan proves "no out-of-bound reads" (manifest/diff comparison shows only workspace files were read); blocked if outside-workspace read traces are found.

### R-09.3 Network unknown (bottom line) = (c) case by case

Loopback acceptable, public network must be blocked; confirmed once before the start, does not block the start.

### R-09.4 Cost unknown = (a) explicitly marked unknown

CONTRACTS §7.1 already requires unknownDimensions to contain cost.

### R-09.5 Background-process unknown = (a) acceptable

The user is responsible for cleanup; the report states it clearly.

### R-09.6 Session-continuity unknown = (a) acceptable

A design boundary of the black-box mode.

### R-09.7 External-side-effect unknown (bottom line) = (c) case by case

Accepted if the user explicitly confirms "unknown external side effects accepted"; blocked otherwise. The confirmation record is signed before the start and does not block the start.

### R-09.8 Mandatory fields = (a) all seven required

assuranceLevel=reduced and unknownDimensions non-empty and containing cost.

### R-09.9 Judgement calibre = (b) per-item independent, no accumulation

One item's acceptable unknown does not pass the whole; one item's block does not fail the whole.

### AT-24 report three sentences (mandatory)

1. This report's assurance level is reduced, not full.
2. The acceptance boundaries of the seven unknown dimensions are: tool use / resource reads = case by case (accepted if the scan proves no overreach); background processes / session continuity = acceptable; network = loopback acceptable, public network blocked; cost = explicitly marked unknown; external side effects = the user has explicitly confirmed acceptance.
3. This report does not constitute a guarantee of the black-box Agent's process correctness; ProofRail only guarantees the post-return artifact scan, the independent gates, the evidence and the reviews.

## 4. Stage Progression

### Stage vs layer, the conceptual distinction

- Stage: the strictness of the judgement standards over time. Stage 1 loose → stage 2 strict. Scope: across multiple AT-23 executions.
- Layer: the test layers within one execution. Layer 0 → 1 → 2 → 3. Scope: a single AT-23 execution.
- The two are orthogonal: stage governs "when the standards tighten", layer governs "how the tests are organised within one run".

### Stage overview table

| Item | Stage 1 (loose) | Stage 2 (strict) | Upgrade trigger |
|---|---|---|---|
| R-08.1 version-check frequency | first + 1-in-5 sampling | every time | all three met |
| R-08.2 probe budget | per-occurrence declaration (same-turn user approval) | preset Copilot 5 / BYOK 100 / retries 5 | same |
| R-08.3 ~ R-08.6 | unchanged | unchanged | — |
| R-09.1 ~ R-09.9 | unchanged | unchanged | — |

Stage progression affects only R-08.1 and R-08.2.

### Stage 1 → stage 2 upgrade trigger (all three simultaneously)

1. The first AT-23 execution (stage 1) has layer 3 all green + conditional closure;
2. the user explicitly rules "upgrade from stage 1 to stage 2" (the controller must not decide on its own);
3. after the upgrade the key AT-23 scenarios must be re-run (layer 1 + layer 2 core scenarios + all layer 3 scenarios) to confirm the strict standards still pass.

### Stage × layer interaction

- The first AT-23 execution (layers 0→1→2→3) uses stage-1 standards.
- After the stage 1→2 upgrade the key scenarios are re-run; passing stage 2 counts as the final sign-off.
- A stage-2 failure rolls back to stage 1, but the report states clearly "not passed under the strict stage-2 standards", and the failure cause is registered as an observation item.
- A "stage 2 not passed" must not negate the stage-1 conclusion.

## 5. Ten revisions summary (v1 → v3)

| # | ID | Revision | Nature |
|---|---|---|---|
| 1 | R-08.6 | add the step-0 minimal verification + failure path | mandatory (highest risk) |
| 2 | R-09.2 | (a) → (c) case by case | mandatory (conflicts with SEC-15) |
| 3 | R-09.1 | (a) → (c) case by case | recommended |
| 4 | R-08.4 | add the BYOK pre-check failure path | recommended |
| 5 | R-08.5 | spell out "at least 3 times" | recommended |
| 6 | R-08.1 | add the "new call sites must use the wrapper" discipline | recommended |
| 7 | layer discipline | distinguish core / non-core scenarios | recommended |
| 8 | R-08.2 | spell out per-occurrence declaration = same-turn user approval | truthfully stated |
| 9 | P60 J3 | truthfully state "historical rows are a human duty, not mechanical coverage" | truthfully stated |
| 10 | AT-23 report | truthfully state "does not cover power-loss scenarios" | truthfully stated |

v3 addition: the stage-progression strategy became its own section; the to-do list gained the stage upgrade.

## 6. Write-back destinations

- §6 R-table: the R-08/R-09 rows of `docs/t027/CHECKLIST_v49{,_EN}.md` §6 flip to 「已裁」 (summary + a pointer to this document); v50 inherits this status when generated.
- One reference line is appended at the end of the P60 block in REMAINING_SLICES.
- The complete ruling is archived under `docs/validation/evidence/R08-R09-RULING/` (MANIFEST.md + SHA256SUMS.txt).

## 7. Relation to ADRs

This ruling does not register an ADR separately. Reasons:
1. The wrapper-function plan of R-08.1 and the BYOK main path of R-08.6 have their trade-offs fully recorded in §2; this document is authoritative.
2. If a later finding needs ADR-level traceability for some item (e.g. the wrapper plan is overturned, or the BYOK plan needs restructuring), it will be back-registered per §1.7.4 starting from ADR-035.
3. R-08.2~R-08.5 and R-09.1~R-09.9 are pure acceptance-judgement calibres without architecture decisions; no ADR is needed.

## 8. Change-cost grading

This table lists the change cost of each ruling, for reference in future revisions.

| Ruling item | Cost of change | Note |
|---|---|---|
| R-08.2 probe budget | Very low | Pure numbers |
| R-08.4 VM/proxy toggles | Very low | Pure environment |
| R-09.1~R-09.9 boundaries | Low (loosen) / Medium (tighten) | Pure docs; tightening costs more |
| AT-23 layering / staging | Low | Execution strategy, does not affect verdicts |
| R-08.1 wrapper function | Medium | **newly created by P61 ②** (Invoke-CopilotCli / buildCopilotArgs); must pass the review chain |
| R-08.6 BYOK config | Medium | Config layer, but AT-23 depends on it |
| P60 J2 criterion | Medium-high | Criterion body + version bump + re-review |
| R-08.3 P → P′ | High | Open VM + run injection matrix |
| R-08.5 conditional-closure → overall-fail | High | Overturns a signed conclusion |
| R-09 loosen → strict (after AT-24 signed) | High | Re-run AT-24 |

> **Point-in-time note (P62, 2026-11-08)**: this row originally read "landed code"; in reality the function was only newly created by P61 ② (before commit 6bb23a1 the repository had no such function). Per §11.1 historical rows are not rewritten, so this row gains a note. The factual error is independent of the design decision (the R-08.1 ruling itself is unchanged).

**Principle**: "things that will change" are designed to be changeable (budget, environment, layers, stages); "things that must not change" are explicitly frozen (verdict criteria, safety bottom lines, landed code).

- ⑥ skipped: this is a `[ROUTINE]` docs-only ruling-record landing with no code / criterion / contract change; §7.2.1 allows a docs-only slice to skip ⑥ with a recorded reason.

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/at23-restricted-candidate.md	repo
docs/validation/at23-restricted-candidate_EN.md	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/00-briefing.md	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/01-arch/design-r1.md	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/04-reviews/05-prereview-r1.txt	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/04-reviews/05-prereview-r2.txt	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/04-reviews/06-scan-r1.txt	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/04-reviews/07-final-r1.txt	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/MANIFEST.md	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/SHA256SUMS.txt	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/freedom-list.md	repo
docs/validation/evidence/P62-R08-R09-RULING-REVISION/input-list.md	repo
docs/validation/evidence/R08-R09-RULING/MANIFEST.md	repo
docs/validation/evidence/R08-R09-RULING/SHA256SUMS.txt	repo
docs/validation/evidence/R08-R09-RULING/r08-r09-ruling.md	repo
docs/validation/evidence/R08-R09-RULING/r08-r09-ruling_EN.md	repo
docs/validation/r08-r09-ruling.md	repo
docs/validation/r08-r09-ruling_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
```

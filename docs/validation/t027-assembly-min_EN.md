# T027-ASSEMBLY-MIN — minimal runnable assembly (segment A: the reject path) validation report

> Slice `T027-ASSEMBLY-MIN` (product implementation, segment A) | tier `[SLICE]` | size M | date 2026-10-03.
> **Status: ⑨ native validation, ⑧b finalisation and ⑧c wrap-up are executed; the ⑩ stop report is section 13 and the slice parks before commit authorisation** (commit / push require same-turn authorisation).

[English](t027-assembly-min_EN.md)

## 1. Change summary (including the execution freedom list)

- **What this slice did**: deliver the segment A minimal runnable assembly — wire the already frozen AgentRunner building blocks into a runnable reject-path chain (routed to the assembly entry only when all eight input flags are present: request record, capability record, enforcement record, availability record, authorisation ledger, cost ledger, pinned CLI executable and version). It adds three production files (`internal/console/runtime_agent.go`, `internal/console/runtime_noop.go`, `internal/adapters/agent_runner_postflight.go`) and two test files, and edits the existing dispatch in `internal/console/runtime.go` and `internal/console/cli.go`.
- **What this slice did not do**: no real enforcement landing (that is segment B); no change to `internal/chain/**`, to existing adapters and console files, to `schemas/**`, to CONTRACTS, to fixtures, to workflows or to `go.mod`; no new sentinel; no positive full-chain evidence on Windows.
- **Artifacts**: this report, the freedom list, the five Go files listed above, `tmp/assembly-min/DESIGN.md` (design plus three remediation records, workspace-local), the mutation script and the evidence readings (`tmp/assembly-min/`, workspace-local).
- **Execution freedom list**: `docs/validation/evidence/T027-ASSEMBLY-MIN-freedom-list.md` (F-1 to F-6).
- **One-line conclusion**: the assembly runs on the reject path and its fail-closed surface was tightened after the third review round (the stop check now demands both pid reports, and the two parsers are one); one hard gate (the F-02 root cause) lies outside this slice and closes as a **recorded state** under the user pre-authorisation.

## 2. Deliverables and their landing points

| Deliverable | Landing point | Note |
| --- | --- | --- |
| Assembly entry (reject path) | `internal/console/runtime_agent.go` (`ExecuteAgentRunnerRun`) | Construction order is fixed: marker, replay store, reject overlapping roots, snapshot, launcher, dispatcher, admission, intent, postflight, publisher; any step failing refuses the run |
| Single failure entry | Same file (the `fail` closure) | All seven post-spawn failure returns pass through the closure for the best-effort reclaim; the only naked failure return sits inside the closure |
| Stop-check criterion | Same file (`managedProcessesGone`) | Both pid reports must exist; a missing report, an unparsable entry or a failed liveness answer is fail-closed |
| Postflight port | `internal/adapters/agent_runner_postflight.go` | Six legs: manifest, diff, log, usage, process-alive, secret scan; an artifact digest must first bind to the frozen fact |
| Non-AgentRunner path unchanged | `internal/console/runtime_noop.go` | The noop runtime and its fail-closed `RunPostflight` stub moved as a unit, semantics unchanged |
| CLI surface | `internal/console/cli.go` (`executeRun`) | Routes to the assembly entry when all eight flags are present and an executable step exists; otherwise the pre-existing error stands |
| Guard tests | `internal/console/runtime_agent_e2e_test.go`, `internal/adapters/agent_runner_postflight_test.go` | Assembly legs in both directions, seven tamper classes for postflight, the parser matrix, the reclaim closure behaviour and the source-level invariant |

## 3. Incidents and fallbacks

1. **Transport failures at ① (four attempts)**: the architecture subagent returned transport errors four times in a row; per the user ruling the retry option was taken — no proxy, spacing of more than two minutes, reduced payload; the fifth attempt succeeded. The retry allowance was not exhausted and the record is kept.
2. **DR-8 candidate (Windows temporary-directory rename jitter)**: `TestRestoreWithHardlinksAndPermissions` in `internal/snapshot` went red during the first parallel full test run (rename in a Windows temporary directory reported `Access is denied`), then green when the package was run alone and green again in the next full run; the package is not in this slice change set, so it is classified as environment jitter rather than slice-introduced, and the first red run is preserved as evidence.
3. **Master error (recorded faithfully)**: to prepare a before-change snapshot for ⑤ the master copied nine files into `tmp/assembly-min/before/`; the `.go` copies were treated as new packages by `go build ./...` and the baseline gates went red. Disposition: the directory was renamed to `_before` (the underscore prefix is ignored by the Go toolchain), the scripts were updated and the baseline returned to green.
4. **Master error (recorded faithfully)**: on the first run of the extended mutation audit the restore write for M16 failed under a Windows file lock (`EBUSY`), leaving the M16 mutation applied in `agent_runner_postflight.go`. Disposition: the original text was reconstructed from the fact that the mutation is a single-occurrence deterministic literal substitution (no mutation marker survived and `gofmt` printed nothing), the audit script gained an atomic write, a `finally` restore and an `EBUSY` retry, and the full re-run reported byte-exact restoration true.
5. **Master error (recorded faithfully)**: after the first ④ the master mis-stated a count as thirteen ok and three failing; three of those lines were markers and only one test actually failed. The report corrects it.
6. **Fallbacks = 2** (items 3 and 4); no other fallback was needed, and no data loss or unauthorised write occurred in this slice.
7. **Gate findings fixed rather than whitelisted (two)**: on the first gate run of the new report `G4b` flagged an unregistered new CJK character and `G7-a` flagged the missing artifacts block; the disposition was to **reword** and to add the block, with no whitelist entry and no criterion change.

## 4. Execution pipeline

| Stage | Role | Status | Artifact |
| --- | --- | --- | --- |
| Step 0 preconditioning | Master | Executed | `DELIVERY_DIRECTIVE{,_EN}` raised to v1.17, ledger backfill and OB entries |
| ① Architecture | V4 Pro | Executed (v1 and v2) | Decision scheme and assembly design (segment A scope statement, reject-path enumeration) |
| ②b Implementation | Flash (master selection) | Executed (with r3 and r4 remediation) | Five Go files, tests and mutation anchors |
| ③ Tests | V4 Pro | Executed | Legs in both directions, tamper matrix, parser matrix, source-level invariant |
| ④ Master integration and gates | Master | Executed (re-run several times) | See section 9 |
| ⑤ Pre-review | V4 Pro | Executed (three rounds, the third a quota exception) | See section 6 |
| ⑥ Independent scan | MAI Flash | Executed | See section 6 |
| ⑦ Final review and re-review | Codex (not downgraded) | Executed (all three rounds used) | See section 6 |
| ⑧a Draft documents | Master | Executed | This report and the freedom list |
| ⑨ Native validation | Master | Executed | Read-only re-runs of the three gate scopes and the encoding check (see section 9) |
| ⑧b Finalisation | Master | Executed | Ledger write-back (the slice block in the bilingual ledger) and the status-line update |
| ⑧c Wrap-up | Master | Executed | Temporary-artifact retention check and work-tree state check |
| ⑩ Stop point | Master | Parked before commit authorisation | Section 13 |

## 5. Environment and tier

- Platform: the local Windows host (Ubuntu CI is the only platform for the positive full chain); tier `[SLICE]`; master tier standard.
- **⑦ used all three rounds without downgrade**: blind review, final review and re-review, model fixed at `gpt-5.3-codex`; the user red line forbids a fourth round.
- **The third ⑤ round is a quota exception**: recorded independently as ADR-018 under the OB-37 meta-rule (narrowed input, not a precedent).
- Probes = 0 (no real external call was made); the trial-deferral four elements are recorded with ⑩.

## 6. Review outcome summary (⑤⑥⑦)

- **⑤ round one** (V4 Pro): `PASS WITH FIXES` — F1 was a safety-relevant Medium (artifacts not bound to the frozen facts), with another Medium and several Low; all dispositions landed in-slice with matching mutations.
- **⑤ round two** (V4 Pro): `PASS WITH FIXES` — nothing above Medium; the remaining Low and two observations were disposed in-slice.
- **⑤ round three** (V4 Pro, quota exception ADR-018): `PASS WITH FIXES` — nothing above Medium; two Low were disposed by a wording fix and a new mutation entry.
- **⑥ independent scan** (MAI Flash, blind): `PASS`, zero findings; because that round material was not written to a standalone file, this report states it as-is.
- **⑦ round one blind review** (Codex, without the ⑤⑥ lists): one High, one Medium (hard gate) and one Low — the stop check was too weak, a failed identity record after spawn did not reclaim, and a comment contradicted the implementation.
- **⑦ round two final review** (Codex, with the ⑤⑥ lists and that remediation): `FINAL REVIEW: FINDINGS` — M-1 (the two parsers had drifted, measured as a console fail-open direction) and M-2 (the reclaim covered only the first non-awaiting branch, so the hard gate stayed open).
- **⑦ round three re-review** (Codex): `RE-REVIEW: PASS WITH FIXES` — M-1 closed; F-02 partially closed and still judged Medium (recorded state under the pre-authorisation); L-3 acceptable as-is; one new Low (the robustness of the source-level invariant) registered as the OB-61 candidate.

## 7. Review input packs (including blind-review isolation)

- **Round three re-review input pack**: `tmp/assembly-min/review-pack-r3.md` (index), `change-set.diff` (whole slice), `r3-remediation.diff`, `r4-remediation.diff`, `mutation.txt`, `pre-review-r1.md`, `pre-review-r2.md`, the gate and test readings, `v3.txt` + `v4.txt` + `snapshot-flake.txt` (the three DR-8 evidences) and `contracts-excerpt.txt`.
- **Earlier packs**: round one was blind (only the change set and the design, no ⑤⑥ list); round two carried the ⑤ round one and two lists, the ⑥ outcome and the evidence readings.
- **Origin of each remediation diff**: the r3 diff is the difference between a byte-exact reconstruction of the pre-r3 state and the work tree; the r4 diff is the difference between a pre-audit snapshot and the work tree, with an assertion that the reference side (the adapters source) is unchanged; both generators assert structure and file their readings.
- **Sanitisation and isolation**: the material sent for review had user names, host names, IPs, credentials and tokens filtered out; round one fed no list so that the blind review stays independent.
- This slice touches state-machine, authorisation-contract and evidence-model semantics, so the mandatory blind review condition held and was executed as round one.

## 8. Falsifiability (mechanical checks and mutation testing)

- Mechanical checks: see section 9 (gates and the full test run).
- **Mutation testing**: for each entry the script applies a unique replacement, runs the focused tests, diffs every per-test verdict against the unmutated baseline and restores byte-exactly (SHA-256 compared). Final state: **17 red, three truthfully marked superseded (M9, M11, M15) and byte-exact restoration true**; the baseline is 22 tests, 21 passing and one skipped.
- **The three superseded entries are recorded rather than deleted**: the old weak guards of M9 and M11 were replaced by the both-reports criterion; the M15 anchor no longer exists after the parsers were unified, and removing the error check is a no-op under a strict parser.
- **What no mechanical check covers (stated plainly)**: the source-level invariant test counts textual shapes and matches strings, so an equivalent rewrite could in theory evade it; this is registered as the OB-61 candidate.

## 9. Gate results

- `node tools/gates/gate.js --all --scope=tree` gives `TOTAL_FAIL=0` (`changed=16`), including a passing baseline `gofmt` / `go build` / `go vet` / `go test` and a passing encoding check.
- `G3(a)` symmetry: this report pair is a new file pair with an equal number of added lines on each side and none deleted (the on-the-spot gate reading governs); `G4b`: every new CJK character is registered and the unknown set is empty; `G7-a` to `G7-d`: the artifacts block is present and every artifact row resolves.
- `node tools/gates/gate.js --all --scope=index` gives `TOTAL_FAIL=0` (`changed=0`; the slice artifact paths are registered in the index per the G7 contract but **no content is staged**, so the index scope stays an empty set and the criteria report an empty set by design rather than a failure).
- `node tools/gates/gate.js --selftest` gives `SELFTEST: PASS 248/248`.
- `go test -count=1 ./...` gives 14 packages passing and 0 failing (the same package set as the baseline).
- Encoding: every changed file is UTF-8 without BOM (Go and JS) or with BOM (Markdown), with LF line endings.

## 10. Known divergences and coverage gaps

1. **The F-01(3) substitute criterion (verbatim)**: an `InspectProcess` error means **the process is treated as gone**. Reason: the literal criterion (an error is fail-closed) is refuted by counter-examples on both platforms — a pid that exited normally reports `ERROR_INVALID_PARAMETER` on Windows (measured line: `AgentRunner evidence invalid: cannot observe pid 49656 from stub.pid: The parameter is incorrect.`) and `ENOENT` from `/proc/<pid>` on Linux, so the literal reading would make every completed run unverifiable; the substitute matches the reading already used by the run side. The fail-closed surface is explicitly narrowed to a failed liveness answer and an unparsable entry.
2. **The F-02 recorded state (user pre-authorisation, verbatim)**: if the re-review still judges F-02 as a Medium (namely mitigated, not closed), the user accepts that state in advance as the final record of this slice: the slice closes with M-1 fixed, the F-02 mitigation maximised and the root cause tracked as OB-57, with no fourth round for F-02 and no scope expansion; this does not break the §5.3 discipline but is the user explicitly accepting an out-of-slice Medium as a recorded state. The re-review judged it partially closed and still Medium, so the recorded state stands.
3. **The three evidences for the DR-8 first-red run (verbatim)**: `tmp/assembly-min/v3.txt` (the first red run: rename in a Windows temporary directory reported `Access is denied`), `tmp/assembly-min/v4.txt` (the isolated package re-run, green) and `tmp/assembly-min/snapshot-flake.txt` (the jitter conclusion excerpt); the package is not in this slice change set.
4. **The master EBUSY error and its recovery (verbatim)**: on the first run of the extended mutation audit the M16 restore write failed under a Windows file lock and the mutation stayed applied in `agent_runner_postflight.go`; the master reconstructed the original from the fact that the mutation is a single-occurrence deterministic literal substitution, then gave the script an atomic write, a forced restore and a retry, after which the full re-run reported byte-exact restoration true.
5. **The eight CLI flags have no CLI-level test (verbatim)**: the eight-flag routing in `executeRun` is backed only by the static production path and the assembly legs, with no CLI-level test; the real caller is backed by the static production path together with the assembly legs in both directions.
6. **Three parity fields that are not byte-equal (verbatim)**: in the parity test `OutputManifestHash`, `Facts.ManifestHash` and `Facts.ProcessStopEvidenceHash` are not byte-equal, so substitute assertions are used instead (the field is present, non-empty and digest-shaped); if `snapshot.CapturedAt` ever becomes injectable these three should be tightened to byte equality.
7. **The F-01 misjudgement risk for the real CLI (verbatim)**: the both-reports criterion depends on the candidate CLI always writing both the parent and the child report; the current stub writes the child report only when it spawns a child, so the real CLI integration in segment B must assess whether the criterion would misjudge a completed run as unverifiable.

## 11. Candidate OBs (not landed in §12.4 by this slice)

| Candidate | Content | Origin |
| --- | --- | --- |
| OB-56 candidate | `AgentRunnerPinnedCLIRun.Run` has no production caller and the non-spawn wait entry was missing; this slice closed the gap with a non-spawn wait entry | this slice A6 legacy note |
| OB-57 candidate | the unclosed root cause of the reclaim: a failed identity record after a successful spawn returns without reclaiming; the root cause is in `internal/adapters/agent_runner_replay_dispatcher.go`, outside `T027-ASSEMBLY-MIN`, to be fixed in a separate slice | F-02 (consistent across all three ⑦ rounds) |
| OB-58 candidate | bring the pid report digest into the frozen facts and bind it to the postflight request (touches the frozen-fact and contract domain) | F-01(2) (user ruled a separate slice) |
| OB-59 candidate | a CLI contract gap: the stub and the real CLI write the child report only when they spawn a child; until it is closed the productive positive chain is unverifiable for want of the child report, and the fixture seeding proves only the landing point, not that production has a producer | the production consequence of F-01(1) (user ruled a separate slice) |
| OB-60 candidate | add two sections to the design-document template — (1) cross-package shared contracts: list every item shared across files (constants, parse functions, decision order, error classification), each with an authoritative side, a verification method and a mechanical check; (2) failure-path enumeration: for every operation that can fail, state who owns the failure, the consequence and how it is proven, exhaustively and without an others-are-similar escape. Applicability trimming: not required at tier S; required at tier M or above and for any cross-package change. | reserved by the user, sourced from the M-1 and F-02 design lessons of this slice |
| OB-61 candidate | the source-level invariant test counts textual shapes and matches strings, so an equivalent rewrite could in theory evade it; upgrading it to an AST-level structural assertion is recommended | the new Low from the ⑦ round-three re-review |

## 12. Open items and next steps

- ⑨ Native validation: executed (read-only re-runs of the three gate scopes and the encoding check, readings in section 9).
- ⑧b Finalisation: executed (the bilingual ledger slice block and completion-log row, plus this report status-line update).
- ⑧c Wrap-up: executed (encoding and residue checks plus a check of the temporary-artifact retention rule).
- ⑩ Stop point: parked before commit authorisation (see section 13).
- Commit / push require same-turn authorisation; no gitee push; the step 0 v1.17 changes are handled with ⑩ per the user ruling.

## 13. ⑩ Stop report

- **⑩ outcome**: stages ① to ⑧a, ⑨, ⑧b and ⑧c are all executed; the slice **parks before commit authorisation** — no commit, no push, no gitee, awaiting same-turn authorisation.
- **Artifact list with paths** (item by item):
  - Implementation: `internal/console/runtime_agent.go` (assembly entry, single-entry reclaim closure, stop-check criterion), `internal/console/runtime_noop.go` (the noop runtime and its fail-closed stub), `internal/adapters/agent_runner_postflight.go` (the six-leg postflight port), `internal/console/cli.go` (the eight-flag routing), `internal/console/runtime.go` (the pre-existing dispatch kept as-is)
  - Tests: `internal/console/runtime_agent_e2e_test.go`, `internal/adapters/agent_runner_postflight_test.go` (parser matrix, closure behaviour and the source-level invariant)
  - Documents: this report pair and the freedom list; the step 0 `DELIVERY_DIRECTIVE{,_EN}` (v1.17) and `ADR_REGISTER{,_EN}` (ADR-018); `tmp/assembly-min/` (design, the three remediation records, the mutation script and the evidence readings, workspace-local)
- **Gates, item by item**: `tree` gives `TOTAL_FAIL=0` (`G3(a)` report pair symmetric, `G4b` unknown CJK empty, `G7-a` to `G7-d` all PASS with every artifact row resolving); `index` gives `TOTAL_FAIL=0`; `SELFTEST PASS 248/248`; the full test run gives 14 packages ok and 0 failing; the encoding check passes throughout (Markdown with BOM, Go and JS without BOM, LF line endings)
- **⑤⑥⑦ outcome lines and quota**: ⑤ returned `PASS WITH FIXES` in all three rounds (the third was an exception beyond the 2/2 quota, ADR-018); ⑥ returned `INDEPENDENT SCAN: PASS` (blind, zero findings); ⑦ returned `FINDINGS` in the blind review, `FINAL REVIEW: FINDINGS` in the final review and `RE-REVIEW: PASS WITH FIXES` in the re-review (the 3/3 quota is used up, with no fourth round)
- **Cost and accounting** (line by line):
  - ① architecture 2 versions; ②b implementation 3 rounds (r2, r3, r4); ③ tests 1 round; the ① transport failures were handled per the user ruling (no proxy, spacing above two minutes, reduced payload, retry allowance not exhausted)
  - ⑤ 3 calls (V4 Pro, including one quota exception); ⑥ 1 call (MAI, blind); ⑦ 3 calls (Codex, not downgraded); real external calls 0
  - 20 mutation entries (17 red, three truthfully marked superseded) with byte-exact restoration true; the focused-test baseline is 22 items (21 passing, one skipped); the mutation audit was repeated 3 times (the first was interrupted by a Windows file lock and the incident is recorded)
  - The authorisation and cost ledgers are consumed read-only by the assembly entry through the existing record files; this slice performs no new billable action
- **The eight intent-to-add paths** (path registration only, no content staged; they become real staging on the authorised commit):
  - `docs/validation/t027-assembly-min.md`
  - `docs/validation/t027-assembly-min_EN.md`
  - `docs/validation/evidence/T027-ASSEMBLY-MIN-freedom-list.md`
  - `internal/adapters/agent_runner_postflight.go`
  - `internal/adapters/agent_runner_postflight_test.go`
  - `internal/console/runtime_agent.go`
  - `internal/console/runtime_agent_e2e_test.go`
  - `internal/console/runtime_noop.go`
- **Trial deferral, sixth time (four elements)**: ruling party = the user; date = 2026-10-03; reason = this product-implementation slice continues execution without entering the trial; basis = the authorisation text of this round. **Not usable as a precedent**.
- **Disposition wording for step 0 (v1.17)**: the v1.17 line in `DELIVERY_DIRECTIVE{,_EN}`, the OB-54 and OB-55 rows in §12.4 and the orchestrator-invocation record block are all step 0 artifacts; per the user ruling they are **handled together with ⑩** (committed in one go with the slice, not separately).
- **Not executed**: real CI observation (the slice is not pushed, so there is no CI verdict); the segment B real enforcement landing; Windows evidence for the positive full chain; the OB-57 root-cause fix; the registration and disposition of OB-58 to OB-61 (they land in the ledger when the next slice starts); the `tmp/assembly-min/` process files are retained per the existing instruction (workspace-local).
- **Authorisation pending**: `git commit` (the gate reads 16 changed paths, 8 of them intent-registered new artifacts) and `git push` (origin only); **no gitee push**.

## 14. Artifact list

- Rows below are tab-separated: repo-relative path and disposition (`repo` means in the repo and resolvable; `local-only` means workspace-only and not committed).

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/validation/t027-assembly-min.md	repo
docs/validation/t027-assembly-min_EN.md	repo
docs/validation/evidence/T027-ASSEMBLY-MIN-freedom-list.md	repo
internal/adapters/agent_runner_postflight.go	repo
internal/adapters/agent_runner_postflight_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_noop.go	repo
internal/console/cli.go	repo
internal/console/runtime.go	repo
tmp/assembly-min/**	local-only
```

- Note: under tree scope a `repo` row requires the path to be **registered in the index**; before commit authorisation this slice registers the path only, without staging content (`local-only` rows are exempt), so nothing enters the repo and the commit still needs same-turn authorisation.

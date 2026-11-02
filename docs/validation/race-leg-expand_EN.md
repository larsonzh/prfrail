# RACE-LEG-EXPAND Verification Report

- Positioning: v38 layer 1 item 2 (P49); CI Race leg (Ubuntu-only) expansion.
- Directive: v1.60; start anchor `81dd49f`; date 2026-11-02.
- Implementation class: `[SLICE]` (infrastructure); size tier M; ledger row `P49`.

## 1. Goal and Scope

- **Goal**: expand the CI Race leg (Ubuntu-only, `CGO_ENABLED=1 go test -race -count=1`) from `internal/adapters` and `internal/chain` to six core internal packages (adding `console` / `gates` / `guard` / `tickets`), bringing concurrency-sensitive paths under the race detector.
- **Change**: one line in the Race step `run:` of `.github/workflows/ci.yml`, plus one line at `internal/release/workflow_test.go:245` syncing the workflow frozen-contract digest (old `ff4826f2…` → new `825bacc3…`). No production code change; no existing step deleted or reordered.
- **Risk rating**: Low. The change is a single CI configuration line plus a frozen-contract digest sync with a narrow semantic surface; residual risk sits in the first CI run's wall-clock time and in the flake observation for newly added packages entering the Race leg (see 7).
- **Untouched**: `.go` semantics for production and tests (the only `.go` change is the frozen-contract digest constant), CONTRACTS / Schema / fixtures, other CI steps; no gitee push; stop at 10 for authorisation.
- **R / D baseline**: R = 1, D = 2. R counts the ② r2 contract-digest remediation; D counts the ① retrieval miss and the main controller's environment pre-check blind spot (see 5).

## 2. Pipeline Execution

| Step | agentName | model | Output |
|---|---|---|---|
| 0 | Main controller | `DeepSeek V4.1 Flash (deepseek)` | Ledger entry, evidence-pack skeleton, freedom list |
| 1 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `01-arch/design.md` |
| GLM 1/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-1-ciyml-scan.md` |
| GLM 2/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-2-race-time-estimate.md` |
| 3 | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | Static red baseline, coverage matrix, VM evidence archive |
| 2 | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | Race-leg expansion (r1), contract-digest sync (r2), `04-impl/implementer-report.md` |
| 4 | `prfrail-implementer` / `prfrail-tester` | same as above | `go build` / `go vet` / `go test ./...` all green |
| 5 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS WITH FIXES` |
| 6 r1 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS WITH FINDINGS` (output shape deviated from the A-E contract) |
| 6 r2 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS` (A-E compliant re-issue, budget 2/3) |
| 7 r1 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS WITH FIXES` |
| 7 r2 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS` |
| 8 | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | This report pair and the documentation sync |

## 3. Red-to-Green Evidence

- Red baseline (3, `ci.yml` unmodified): the assertion script exits **1** with a red-state verdict line; the prefix / race / count / Linux-guard checks and both legacy package checks are all PASS, so the red lands exactly on the missing coverage of the added packages and the package-order failure ⇒ a falsifiable red, ruling out a false red caused by a parser defect. Log `03-tdd-red-run.txt`.
- Independent re-check (script-free): a text search reads `ci.yml` directly; the Race line names only the two legacy packages, and the added package patterns match zero times in the whole file ⇒ the script and the text search agree.
- Flip (re-run after step 2): exit **0** with a green-state verdict line; each of the six packages appears exactly once in the order adapters → chain → console → gates → guard → tickets, with the prefix and the Linux guard still PASS. Log `03-tdd-green-run.txt` (the red-baseline file is not overwritten).
- Local run infeasible: the Windows host lacks gcc and `CGO_ENABLED=1 go test -race` fails at build time ⇒ local results are not documentary race evidence. Log `03-race-local-infeasible.txt`.
- Real `-race` on the VM (Ubuntu VM `10.0.0.199`, isolated toolchain go1.22.12, `rc=0`, run window about 50 seconds): all added packages green, no historical data race observed. Log `vm-race-four-pkgs-green.txt`.

| Package | VM `-race` reading | Evidence |
|---|---|---|
| `internal/console` | `ok` 20.410s | `03-tests/logs/vm-race-four-pkgs-green.txt` |
| `internal/gates` | `ok` 3.384s | same |
| `internal/guard` | `ok` 3.043s | same |
| `internal/tickets` | `ok` 7.420s | same |

- Workflow frozen contract (step 2 r2): before the fix the `internal/release` workflow-contract test family was red with a single cause (the `test/Race` step-body digest did not match the expected value), log `04-red-before-contract-sync.txt`; after the fix the package is green and the per-test verbose run is all PASS, logs `04-release-green.txt` and `04-release-workflow-tests-verbose.txt`.

## 4. Review Chain Results

- 5: `PRE-REVIEW: PASS WITH FIXES` - F1 (Medium) is a step-1 retrieval miss: the frozen contract of `internal/release` pins the `test/Race` digest, closed by the r2 sync and confirmed by an independent digest recomputation. Evidence `05-prereview/review.md`.
- 6 r1: `INDEPENDENT SCAN: PASS WITH FINDINGS` - one Low (the evidence pack must travel with MANIFEST and the checksum file) and one Info (self-reported status is not independent evidence); the independently recomputed digest matches `workflow_test.go`, and no third live verification site remains. Evidence `06-scan/review.md`.
- 6 r2: r1 did not follow the Section A-E contract of section 7.2.2, so it was re-issued once under the "format deviation allows one re-issue" rule; r2 produced a compliant `INDEPENDENT SCAN: PASS` with no findings in Section B. Evidence `06-scan/review-r2.md`. Verbatim transcript below:

```sixstage
## Section A
I reviewed the only-read diff for the CI Race expansion and the matching workflow freeze contract: one run line in .github/workflows/ci.yml, one digest entry in internal/release/workflow_test.go, and the ledger/evidence package for P49. The scope remains infrastructure-only: no production Go code changes, no CI step deletions, and no package coverage outside the six-package Linux race target.

## Section B
(no findings)

## Section C
git diff shows exactly the intended changes: ci.yml changed only the Race run target from two packages to six packages; workflow_test.go changed only the digest for test/Race from ff4826f2103a2457cd306662f2542e54d809546deb1f804964b63f2148185dd5 to 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9; the ledger files and evidence directory were updated as expected. The SHA-256 recomputation over the exact run string "CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/..." produced 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9, matching the frozen contract. For runtime evidence, the VM log in docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt shows rc=0 and these four package results: console ok 20.410s, gates ok 3.384s, guard ok 3.043s, tickets ok 7.420s. The evidence bundle also includes 04-impl/logs/04-release-workflow-tests-verbose.txt and SHA256SUMS.txt, and the file hashes are internally consistent with MANIFEST.md; the package body is UTF-8 BOM+LF as required.

## Section D
If the Linux guard were removed from the Race step, the workflow contract and red/green TDD assertion would fail because the review enforces runner.os == 'Linux' and the package list must be exactly the six targets. If the old digest were restored in workflow_test.go, the workflow contract test would fail immediately with a SHA mismatch at test/Race. If the package list were truncated back to adapters and chain, the red-baseline assertion and the VM coverage log would no longer match the accepted six-package scope. No uncovered invariant remains: the Linux-only condition, the six-package target set, the frozen digest, and the evidence-package integrity are all directly checked.

## Section E
I could not independently re-run GitHub Actions in the remote repo; only the archived VM and bundle evidence are available.
What would overturn this PASS: a current HEAD run where the workflow contract test fails on the digest or Linux guard; a mismatched SHA256SUMS.txt or MANIFEST.md for the evidence bundle; or any real race failure in the four newly added packages on the same Ubuntu toolchain.

INDEPENDENT SCAN: PASS
```

Disclosure of merge differences: the fence is a verbatim transcript of the 6 r2 final message, with no merged differences (the conclusion line appears only once, at the end of the block).
- 7 r1: `FINAL REVIEW: PASS WITH FIXES` - F49-07-01 (Medium) the evidence list did not cover the full set and F49-07-02 (Low) two review documents were encoding-non-compliant; both fixed. Evidence `07-final/review.md`.
- 7 r2: `FINAL REVIEW: PASS` - full recomputation of payload digests with no mismatch (count 24, missing 0), the digest list / MANIFEST archive list / on-disk payloads agree as three equal sets, sampled encodings all compliant, `internal/release` re-run green. Evidence `07-final/review-r2.md`.
- 6 output shape: r1 used a structured shape (checks performed / item-by-item conclusions / findings / overall conclusion) that did not follow the `Section A-E` contract; it was re-issued as r2 per the section 7.2 failure path and produced an A-E compliant output (transcript in the fence above), with the r1 original kept verbatim in `06-scan/review.md`.

## 5. Deviations and Anomalies

1. **Deviation facts**: two deviations occurred (D1 retrieval miss, D2 environment pre-check blind spot); neither changes a code conclusion, but each was first exposed outside step 7.
2. **Root cause**: technical layer - D1 because `internal/release` pins CI script step bodies by digest, so editing the `run:` line of `ci.yml` necessarily trips that contract; D2 because a non-interactive SSH session does not load the login environment (`~/.bashrc` returns early for non-interactive shells) and the system-directory scan did not cover user-level directories. Process layer - D1 because the step-1 impact retrieval did not include workflow frozen contracts in its required scope; D2 because the environment pre-check started installing before enumerating existing toolchains.
3. **Impact**: D1 turned the whole contract family red after r1 and required one added line of remediation in step 2 r2 (R plus 1); D2 produced no wrong conclusion but caused one unnecessary installation (D plus 1) and put the pre-check's credibility in question.
4. **Corrective actions**: for D1, step 2 r2 synced the digest to `825bacc3…` and re-ran green, with 5 r2 recomputing the same value and 6 confirming no third live verification site; for D2, the main controller recorded the factual position per the same-turn user ruling (the VM already had go1.27.1 at `/home/larson/.local/bin/go`; this slice's isolated verification used go1.22.12 at `~/go1.22/bin/go`) with no re-run.
5. **Evidence anchors**: D1 in `04-impl/logs/04-red-before-contract-sync.txt`, `04-impl/logs/04-contract-diff.txt` and `04-impl/logs/04-release-workflow-tests-verbose.txt`; D2 in the "main controller addendum - environment fact correction" section of `05-prereview/review.md`.
6. **Recurrence prevention**: step-1 retrieval must cover workflow frozen contracts such as `internal/release` (non-recurrence statement); remote SSH environment pre-checks must `source ~/.bashrc` or `bash -lc` to load the login environment and must explicitly probe user-level directories, confirming no existing toolchain before any minimal installation (folded into the section 7 environment pre-check rule of v39).
7. **User ruling correspondence**: the D2 position was established by the same-turn user ruling and is quoted verbatim - 「主控未加载用户环境，漏查了已有的 go1.27.1，自行安装了 1.22.12。虽工具链版本与 CI 完美对齐（高保真），但存在预检盲区事实。」 The D1 remediation duty comes from pre-review F1 and this slice's authorisation boundary.

## 6. Documentation Sync

- `docs/ARCHITECTURE.md` and `docs/ARCHITECTURE_EN.md`: a test-infrastructure paragraph appended at the end of section 1 - the tests of the six core internal packages run under the CI Ubuntu Race leg for race detection; the Windows leg skips the Race step and is unaffected by cgo.
- `docs/DEV_PLAN.md` and `docs/DEV_PLAN_EN.md`: a full-text search for Race-leg / race-detection / CI-coverage wording found no outdated milestone such as "incomplete CI coverage" or "the Race leg only covers" (the leg's package scope never entered DEV_PLAN), so they are unchanged.
- Ledger backfill: `docs/t027/REMAINING_SLICES.md` and its `_EN` mirror flip the status table and the P49 section to completed, and append a completion record.

## 7. Known Limitations

- The host is Windows without gcc, so `-race` cannot run locally; race truth is carried by the VM isolated verification and the Ubuntu CI leg.
- The added packages enter the CI Race leg for the first time, with flake observation owed (this repository has a history of timing-jitter first-run reds); a single zero-violation run is not a strict proof.
- First CI run obligation: after the commit, capture the `ok` lines of the six packages and the wall-clock duration, and append them once under `## Commit and CI`.
- The VM isolated verification used go1.22.12 (aligned with CI), not the VM's existing go1.27.1; cross-version behaviour differences are out of scope here.
- G7-f (the 6 output contract): after the r1 format deviation, 6 was re-issued as r2 under the section 7.2 failure path; r2 emitted the A-E shape and was transcribed verbatim into the section 4 fence; this item is closed.
- After 7 r2 the evidence pack was supplemented at 8b with `07-final/review-r2.md` (payload count moves from 24 to 25, with MANIFEST and the checksum file kept in sync).
- G7-b in the `tree` / `index` scope requires every `repo` path declared in the fence to be staged in the index; at hand-off the added payloads are untracked, so the main controller must run a precise `git add` before the step 9 gate run.

## 8. Gate Readings (9)

- Appended after step 9: the three-scope `TOTAL_FAIL` and `SELFTEST` readings.

## 9. Commit and CI

- Appended once after step 9.

```artifacts
.github/workflows/ci.yml	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/t027/CHECKLIST_v39.md	repo
docs/t027/CHECKLIST_v39_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/00-briefing.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/01-arch/design.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-race-local-infeasible.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-green-run.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-red-run.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tests-norace-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/assert-race-leg-coverage.js	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/red-baseline.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/tester-report.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/implementer-report.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-ciyml.diff	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-contract-diff.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-gotest-all-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-red-before-contract-sync.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-release-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-release-workflow-tests-verbose.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/05-prereview/review.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/06-scan/review.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/06-scan/review-r2.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/07-final/review-r2.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/08-gates/gate-readings.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/07-final/review.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/09-glm/glm-1-ciyml-scan.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/09-glm/glm-2-race-time-estimate.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/freedom-list.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/input-list.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt	repo
docs/validation/race-leg-expand.md	repo
docs/validation/race-leg-expand_EN.md	repo
internal/release/workflow_test.go	repo
tmp/RACE-LEG-EXPAND/	local-only
```

This slice stops at ⑩ awaiting authorization.

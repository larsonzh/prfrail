# ACCEPTANCE-SECURITY-COUNTEREXAMPLES (P52) acceptance report

Slice: `ACCEPTANCE-SECURITY-COUNTEREXAMPLES` (security-validation `[SLICE]`, tier M). Package time 2026-11-03 (`HEAD` = `941f60e`, directive v1.60). Controller model = `DeepSeek V4 Pro`.

## 1. Goal and boundary

- Land the GAP-51/52/55/60/64/70 security counterexamples against the acceptance-chain gap of breakpoint A, proving the system fails closed on illegal input, path forgery, handle overreach and boundary bypass.
- Test cases and assertions only, never touch production code (real security holes must be reported, not fixed in this slice); no acceptance-semantics change, no `docs/CONTRACTS{,_EN}.md` change.
- A passing counterexample proves the attack was blocked; any failing counterexample (attack not blocked) is a real hole - stop immediately and report.

## 2. Step-0 record

- Anchor check: HEAD = `941f60e`, origin/main = `941f60e`, directive v1.60, SELFTEST 697/697, three scopes TOTAL_FAIL=0; the HEAD `2b0af2d` recorded by v41 differs only because of the backfill commit (briefing disclaimer).
- Controller model: `DeepSeek V4 Pro` (explicit in the briefing, security adversarial/validation upgrade).
- Ledger: `REMAINING_SLICES{,_EN}.md` P52 entry registered `🔄 In progress`, starting from R/D = 0.
- Briefing verbatim: `00-briefing.md` (BOM+LF); freedom list: `freedom-list.md`; input list: `input-list.md`.
- **Breach finding (2026-11-03, pre-S2)**: T51-1 middle-symlink escape breached on Linux (go1.22.12) - `OpenBoundDirFile` returned an out-of-bound handle and read 23 out-of-bound bytes. Per rule 3: stop, archive the payload (`05-attack/t51-1-linux-escape.md`), register GAP-77. No production code is fixed in this slice.

## 2b. Suspected hole report

| Id | Severity | Status | Content |
|---|---|---|---|
| GAP-77 (V-1) | HIGH | breached (Linux real run) | `openat(O_NOFOLLOW)` call at `bounddir_unix.go:64` guards only the final component (:66 is the error-return line); a middle symlink escapes the bound dir; fix in the P53 slice |
| GAP-78 (V-2) | LOW | candidate (code reading) | `seal.go` validate admits an empty ContextKey (both consumers fail closed; the validate contract should tighten) |
| GAP-79 (V-3) | LOW | candidate (code reading) | asymmetric normalisation: AT23Hash TrimSpace vs SealHash whitespace rejection |
| GAP-80 (V-4) | obs. | candidate (residual) | Windows `isUnderDirectory` spelling prefix rather than object identity |

## 3. Counterexample matrix and results

| GAP | Counterexample | Landing file | Sentinel / assertion | Result |
|---|---|---|---|---|
| GAP-51 | T51-1 | `bounddir_escape_unix_test.go` (`!windows && !darwin`) | attack-outcome assertion (no sentinel, P53 introduces); skip wording `Known vulnerability GAP-77: tracked for P53` | SKIP (breach proven) |
| GAP-52 | T52-1 | `bounddir_escape_windows_test.go` (`windows`) | wording "outside the bound directory" | PASS |
| GAP-52 | T52-2 | same | positive control (in-bound junction admitted) | PASS |
| GAP-55 | T55-1..5 | `boundread_counterexample_test.go` | `ErrBoundFileClosed` / `ErrBoundArtifactExceedsBound` / `ErrBoundFileNotRegular` | all PASS |
| GAP-55 | T55-6/7 | `bounddir_identity_test.go` (`!darwin`) | wording "is not" / four escape shapes refused | all PASS |
| GAP-60 | T60-1..3 | `state_seal_counterexample_test.go` | `ErrPaused` + `evidence-seal-mismatch` / `ErrSealIncomplete` / current-state pin | all PASS |
| GAP-64 | T64-1..6 | `probe_permit_counterexample_test.go` | `ErrProbePermitMissing` / `ErrProbePermitConsumed` / triple negation | all PASS |
| GAP-70 | T70-1..4 | `readiness_counterexample_test.go` | `ErrReadinessEvidenceIncomplete` / current-state pin | all PASS |

## 4. Review chain

- Stage 1 architecture: DESIGN DONE (`01-arch/design.md`).
- Stage 3 test first: PROBE: ALL-GREEN (5 probe faces + four-package baseline, `03-tests/red-baseline.md`).
- Stage 2 implementation: ②a–②d four serial batches all green (`04-impl/implementer-reports.md`).
- Stage 5 pre-review: PRE-REVIEW: PASS (0 Medium+; 2 Low disposed + 3 Info; `05-prereview/review.md`).
- Stage 6 independent scan: INDEPENDENT SCAN: PASS (verbatim below).
- Stage 7 final review: FINAL REVIEW: PASS (0 Medium+; 2 Info; `07-final/review.md`).

## 5. Gate readings

- `gofmt -l` clean on the seven new files; `go build ./...`, `go vet ./...` pass; `GOOS=linux go vet ./internal/adapters/` opposite-OS recheck passes.
- `go test -count=1 ./...`: fifteen packages all `ok`, no FAIL.
- Three scopes `TOTAL_FAIL=0` and `SELFTEST 697/697` in `08-gates/gate-readings.txt`.
- Commit / CI: not committed (stopped at 10 for same-turn authorisation).

## 6. Stage-6 independent scan verbatim

The following is the stage-6 independent scan conclusion verbatim; no evaluation outside the fence.

```sixstage
## Section A
Scope: staged review covers the 22-file package for this slice; under internal/ the staged diff is exactly 7 added counterexample test files, with no non-_test.go changes. The audit focused on whether the new tests truly exercise the claimed fail-closed paths, whether skip logic is honest, and whether the assertions are tied to the real contract rather than incidental errors.

## Section B
Finding: no material defect in the slice as a fail-closed acceptance-security counterexample package. The skip branches are not blanket passes; they are conditional on environment limits or a known vulnerability backlog. The Linux escape case is intentionally skipped before the vulnerable call and is tracked as a known gap, while the red-path evidence is kept in the separate attack record rather than hidden. The Windows junction tests skip only when the environment cannot create a junction; they are not silently passing when the platform supports the primitive. The child-process helper skips only when the environment variable is absent; that is harness gating, not an evasion. The non-sentinel assertions are anchored to the documented contract: either the exact boundary wording is required, or the error must not be one of the known sentinels. These are not loose text matches on arbitrary failures. The cases are written as falsification checks: remove the corresponding guard and the test becomes red.

## Section C
Evidence reviewed: 1. Staged diff confirms the internal scope: 7 new test files and no production changes under internal/. 2. The skip is placed before the vulnerable call in the Unix escape case, so it does not automatically pass the attack path; it exits before entering the known vulnerable branch. The separate attack record documents the actual exploit outcome: ATTACK-SUCCEEDED with SECRET-OUTSIDE-BOUNDARY. 3. The sentinel-heavy tests bind directly to concrete sentinels via errors.Is, matching the contracts for closed handles, over-bound growth, invalid regular files, missing permits, incomplete readiness, and incomplete seals. 4. The non-sentinel tests are not arbitrary; they assert the exact documented behavior for a missing sentinel or a boundary message, such as the directory-identity mismatch and the outside-the-boundary junction rejection. 5. The replay test is not mock-only: it persists a real seal to the file-backed store and then drives the engine replay path that compares the stored context key to the derived one, so the failure is tied to the real acceptance logic and not a fake object. 6. The environment gates are narrow: second-writer failure, junction creation unavailable, helper env absent, and known vulnerability backlog. None of them are broad CI exemptions.

## Section D
Falsifiability audit: The suite is shaped to reject the underlying defense: if any corresponding guard is removed, the relevant case should fail. The sentinel assertions are precise and one-to-one with the contract, which is the strongest anti-false-positive pattern. The documented asymmetries (empty context key and trimmed AT23) are explicitly noted as current behavior, not hidden defects. They are record-keeping at the contract boundary rather than silent acceptance. The only deliberate exception is the known Linux middle-symlink gap, which is acknowledged as a tracked fix item and is not presented as a passing check. This is a legitimate gap marker, not a fake pass.

## Section E
Unverified or non-fully-validated parts: The Linux middle-symlink exploit is not run in the normal unit suite because it is intentionally skipped. It is only proven in the separate attack record, so the CI path is not directly executing the exploit. The Windows junction branch remains platform-dependent: if mklink/junction creation is unavailable, the test is skipped and no boundary proof is available on that runner. The empty-context and trim-space behaviors are currently accepted as asymmetries; they are part of the audit but not part of the fail-closed guard set itself. Reverse conditions: if any skip becomes unconditional, or if the vulnerable path executes without the documented precondition, the scan would flip to findings. If a production file under internal/ changes in the staged diff, or if a sentinel assertion can be satisfied by an unrelated error, the scan would also flip to findings. If the exploit can be reproduced without the skip gate in a normal test run, the T51-1 case is not valid as a tracked-vulnerability record and the package would fail the audit.

INDEPENDENT SCAN: PASS
```

## 7. Unclosed items and boundary

- GAP-77 fix deferred to P53 (this slice only registers it and keeps the evidence chain via the skip).
- The two Windows junction tests skip where mklink is unavailable (coverage availability depends on a capable runner).
- T64-6 simplifies the stage-1 design (sequential form, see stage-5 L2 backfill).
- Deviation D1: the ②b cross-compile artefact was written to the system temp directory, violating the v41 §8 temp-location closed set; the controller verified no residue remains.

## 8. artifacts fence

The following is this slice's archive item list (path + disposition, TAB separated, ASCII order).

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/acceptance-security-counterexamples.md	repo
docs/validation/acceptance-security-counterexamples_EN.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/00-briefing.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/01-arch/design.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/03-tests/red-baseline.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/04-impl/implementer-reports.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/t51-1-linux-escape.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-prereview/review.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/06-scan/review.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/07-final/review.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/08-gates/gate-readings.txt	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/09-glm/glm-1-bounddir-boundread-inventory.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/09-glm/glm-2-gates-chain-inventory.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/MANIFEST.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/SHA256SUMS.txt	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/freedom-list.md	repo
docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/input-list.md	repo
internal/adapters/bounddir_escape_unix_test.go	repo
internal/adapters/bounddir_escape_windows_test.go	repo
internal/adapters/bounddir_identity_test.go	repo
internal/chain/state_seal_counterexample_test.go	repo
internal/evidence/boundread_counterexample_test.go	repo
internal/gates/probe_permit_counterexample_test.go	repo
internal/gates/readiness_counterexample_test.go	repo
```

## 9. Metering

- R = 0 (stage 5 first round had no Medium+, no re-run triggered; stage 7 passed on the first round, no re-review triggered).
- D = 1 (D1: the ②b cross-compile artefact was written to the system temp directory, violating the v41 §8 temp-location closed set; the controller verified no residue remains).
- GLM batch assists dispatched twice; paid review rounds: stage 5 once, stage 6 once, stage 7 once.

# DR13-FIX-D3 Verification Report

- Positioning: v37 layer 1 item 1; the D3 remainder slice of the P45 prerequisite `DR13-FIX`.
- Directive: v1.60; start anchor `e742224`; date 2026-11-02.
- Implementation class: `[SLICE]`; size tier S; ledger row `P48`.

## 1. Goal and Scope

- **Goal**: `release()` now runs `os.Remove` first, then the unlock, then the close, composed with the existing D1 same-process ownership gate and D2 Stat re-check, closing the millisecond-wide concurrency window W1.
- **Landing**: `internal/adapters/agent_runner_starting_marker.go:108` function start; `:112` remove first; `:113` unlock next; `:114` close last.
- **Regression scope**: full Windows remove-while-open regression; dual-platform build and full-package tests.
- **Untouched**: D1 / D2 implementations, the `winners==1` assertion, CONTRACTS / Schema / fixtures, `ci.yml`; no gitee push.

## 2. Pipeline Execution

| Step | agentName | model | Output |
|---|---|---|---|
| 0 | Main controller | `DeepSeek V4 Pro (deepseek)` | Ledger entry, evidence skeleton, freedom list |
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `01-arch/design.md` |
| GLM 1/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-1-case-scan.md` |
| GLM 2/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-2-impact-scan.md` |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | Red baseline, new test set |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | Production change, `04-impl/implementer-report.md` |
| ④ | `prfrail-implementer` / `prfrail-tester` | same as above | `go build` / `go vet` / `go test ./...` all green |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS` |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS` |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS` |
| ⑧ | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | Report pair and documentation sync |

## 3. Red-to-Green Evidence

- First red run (unfixed code): `TestStartingMarkerReleaseRemovesBeforeUnlock` reported FAIL, because unlock sat before `os.Remove`; log `docs/validation/evidence/DR13-FIX-D3/03-tests/logs/go-test-red.txt` (the raw process artefact was archived at 8c).
- Stress red proof (supplementary, probabilistic stance): `TestStartingMarkerReleaseNeverExposesReleasedContent` hit violations 4/35 on its first run, same log.
- Main-controller independent re-check: `docs/validation/evidence/DR13-FIX-D3/03-tests/logs/maincontroller-red-verify.txt`.
- Focused run after the fix: `docs/validation/evidence/DR13-FIX-D3/03-tests/logs/mc-focused-green.txt`, 6/6 green.
- Green baseline: remove-while-open direct test, idempotence test, Windows fail-closed name-occupancy counter-test, ShareFact probe.

## 4. Review Chain Results

- ⑤: `PRE-REVIEW: PASS` (2 Low + 3 Info). Low-1 is the theoretical false-green shape of the order gate, handed to ⑦ and left without a code change; Low-2 is the unverified unix runtime leg on this host, covered by the dual-platform CI; the Info items need no fix.
- ⑥: `INDEPENDENT SCAN: PASS`, findings 0; the independent re-run was 6/6 green.
- ⑦: `FINAL REVIEW: PASS`, findings 0; a non-blocking observation is registered - the order gate judges by the earliest call position, so a counter-example in an unreachable branch that the behavioural tests do not catch should be upgraded to a defect with reachability-level assertions added.

```sixstage
## Section A

The verdict line of the transcript is merged into the closing line of this fence.

## Section B

(no findings)

## Section C

From the code and the targeted tests, the D3 reorder holds and introduces no new concurrency window or caller-semantics drift. The release sequence is remove first, then unlock, then close; when removal fails the lock and handle are released unconditionally, so the failure direction is fail-closed. The boundary audit found only the release function changed; the return semantics are unchanged.

## Section D

The judgment is overturned by any evidence that unlock runs before the remove, or that the lock stays held after a failed removal; the structural gate and the fail-closed counter-test fail immediately in those cases.

## Section E

The unix runtime leg was not verified on this host; a dual-platform CI failure in a targeted path overturns the pass.

INDEPENDENT SCAN: PASS
```
Disclosure: the transcript Section A verdict line is merged into this fence closing line, exactly one difference.

## 5. Documentation Sync

- `docs/ARCHITECTURE.md` and `docs/ARCHITECTURE_EN.md` gain a lock-order entry under §4.2: remove the path before releasing the lock, so a scanner cannot observe an idle lock with a path still present; when removal fails the lock is released unconditionally, and the residual converges through the ownership gate and the next scan.
- `docs/DEV_PLAN.md` and `docs/DEV_PLAN_EN.md`: no DR13 / starting marker / P48 milestone description exists, so they are unchanged.

## 6. Known Limitations

- The order gate is a source-order gate, not a runtime proof; see the ⑦ observation for its theoretical false-green shape.
- The stress regression is probabilistic (R3 stance): 3s window, no `-race`; zero violations is not a strict proof, and the structural guarantee comes from the order gate plus the D1 / D2 composition argument.
- The unix runtime leg did not run on this host (Windows); `RefusesLiveHostUnix` and the orphan-adoption path take the dual-platform CI as the final criterion.
- Fact R-N1: on this toolchain an `os.OpenFile` handle blocks `os.Remove` (errno 32), matching the comment in `agent_runner_starting_marker_windows.go:97-101`; a `release` returning false under high concurrency is therefore a design-scoped R5 convergence path.

## 7. Gate Readings (9)

- Three scopes: index / tree / ci all `TOTAL_FAIL=0` (changed=27; G6-1 / G6-4 scanned 282 added lines).
- `SELFTEST: PASS 697/697`.
- Base legs: gofmt empty, `go build ./...`, `go vet ./...`, `go test ./...` all green (ok=15, FAIL=0).

## 8. Commit and CI

- Main commit: `82e3ac8` (`feat: harden starting marker release order against W1 window (P48-DR13-FIX-D3)`), pushed to `origin/main` (no gitee push).
- CI: run `41379998477` (head `82e3ac8`) **green on both legs on the first attempt** (no rerun).
- Both legs report in `Test` the line `ok github.com/larsonzh/prfrail/internal/adapters`; the ubuntu `Race` step also reports ok for `internal/adapters` and `internal/chain`.
- Both legs `Gates`: `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=27`, `TOTAL_FAIL=0`.
- The three dispatch-only jobs (candidate build / candidate probe / bootstrap release) are skipped as usual.

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/CHECKLIST_v38.md	repo
docs/t027/CHECKLIST_v38_EN.md	repo
docs/validation/dr13-fix-d3.md	repo
docs/validation/dr13-fix-d3_EN.md	repo
docs/validation/evidence/DR13-FIX-D3/00-briefing.md	repo
docs/validation/evidence/DR13-FIX-D3/01-arch/design.md	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/red-baseline.md	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/tester-report.md	repo
docs/validation/evidence/DR13-FIX-D3/04-impl/implementer-report.md	repo
docs/validation/evidence/DR13-FIX-D3/05-prereview/review.md	repo
docs/validation/evidence/DR13-FIX-D3/06-scan/review.md	repo
docs/validation/evidence/DR13-FIX-D3/07-final/review.md	repo
docs/validation/evidence/DR13-FIX-D3/09-glm/glm-1-case-scan.md	repo
docs/validation/evidence/DR13-FIX-D3/09-glm/glm-2-impact-scan.md	repo
docs/validation/evidence/DR13-FIX-D3/MANIFEST.md	repo
docs/validation/evidence/DR13-FIX-D3/SHA256SUMS.txt	repo
docs/validation/evidence/DR13-FIX-D3/freedom-list.md	repo
docs/validation/evidence/DR13-FIX-D3/input-list.md	repo
internal/adapters/agent_runner_starting_marker.go	repo
internal/adapters/agent_runner_starting_marker_order_test.go	repo
internal/adapters/agent_runner_starting_marker_order_windows_test.go	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/logs/go-test-red.txt	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/logs/maincontroller-red-verify.txt	repo
docs/validation/evidence/DR13-FIX-D3/03-tests/logs/mc-focused-green.txt	repo
docs/validation/evidence/DR13-FIX-D3/08-gates/gate-readings.txt	repo
tmp/DR13-FIX-D3/	local-only
```

This slice stops at ⑩ awaiting authorization.

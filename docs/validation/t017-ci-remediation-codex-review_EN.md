# T017 CI Remediation GPT-5.3 Codex Independent Review Record

[中文](t017-ci-remediation-codex-review.md)

Date: 2026-09-09. Reviewed baseline: `12fa05d1f56c5e9a03a0964fce9d9f842bb85397`. Reviewed object: the complete current uncommitted T017 CI-remediation diff relative to that baseline. Reviewer: GPT-5.3 Codex. Final verdict: `PASS FOR BASELINE REVIEW`.

This verdict allows the fixes to proceed as a new verifier-baseline candidate. It is not proof of T017 completion, Go 1.22.12 or GitHub Windows/Linux CI success, or release authorization.

## Findings

No reproducible `BLOCKER`, `HIGH`, or `MEDIUM` was identified.

1. `LOW`: the read from `waitDone` in `internal/guard/process_test.go` has no independent timeout. If a future platform implementation reports the identity gone while `Process.Wait` does not return, the test can wait until its global timeout. Production remains fail-closed and local repeated validation was stable, so this does not block the baseline. The minimum hardening is a `select` with an independent timeout.

## Static Review Conclusions

1. `StopProcessIdentity` uses bounded waiting for a context without a deadline and falls back to one second when `grace <= 0`; an already-cancelled parent or existing deadline is not relaxed.
2. Uncertain termination still returns `ErrTerminationUncertain`; live, PID-reused, or identity-uncertain processes are not incorrectly reported as `stopped`.
3. Child `Wait` results are classified with `errors.As(..., *exec.ExitError)` rather than platform-specific error text.
4. The probe path change only removes the Unix collision between an extensionless binary and its output directory; production behavior is unchanged.
5. The license test builds the real `cmd/prfrail` binary and retains the fail-closed missing-license assertion.
6. Export comparison resolves the deepest existing parent and appends the absent suffix. Non-`IsNotExist` errors remain fail-closed, and bidirectional overlap checks remain intact.
7. Documentation accurately records the dual-platform failure of GitHub Actions run `34253473081`; T017 remains `[ ]` and `IMPLEMENTED_AWAITING_CI`.

## Actual Review Scope

The reviewer confirmed exactly 11 files relative to `12fa05d1f56c5e9a03a0964fce9d9f842bb85397`: five Go fixes, four status documents, and two DeepSeek V4 Pro remediation-review records. There were no additional changes; `sessbridge` was unchanged and `prfrail/tmp/` contained only `.gitkeep`.

## Local Validation Evidence

1. `git diff --check`, `git diff --stat`, and complete target-file diff review: exit 0.
2. `gofmt -l`: exit 0 with no output.
3. `CGO_ENABLED=0 GOTOOLCHAIN=local go build ./...`: exit 0.
4. `CGO_ENABLED=0 GOTOOLCHAIN=local go vet ./...`: exit 0.
5. `CGO_ENABLED=0 GOTOOLCHAIN=local go test -count=1 -timeout=2m ./...`: exit 0; all packages passed.
6. The guard stop test repeated 10 times, the two release remediation tests repeated three times, and the snapshot overlap test repeated five times: all passed.
7. Encoding checks passed for all 11 reviewed files, and editor diagnostics reported no errors.

The local toolchain was Go 1.27.0 windows/amd64. Validation did not run under Go 1.22.12 or on GitHub Windows/Linux runners.

## Remaining Gates

1. The remediation commit has not landed, and the new verifier baseline's Windows/Linux push CI has not run.
2. A different candidate SHA must still be created, followed by the complete candidate-build/probe/selfhost matrix manually dispatched from `main`.
3. T017 may be marked `COMPLETE` only after every Windows/Linux trust job succeeds and the run URL, verifier/candidate SHAs, platform results, and artifact summary are archived.

This review was strictly read-only: it modified, formatted, staged, or deleted no files; performed no Git write; triggered no GitHub workflow, Release, signature, or attestation; and used no network or dependency installation.
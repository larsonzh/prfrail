# T017 CI Remediation DeepSeek Independent Review Record

[中文](t017-ci-remediation-deepseek-review.md)

Date: 2026-09-09. Reviewed baseline: `12fa05d1f56c5e9a03a0964fce9d9f842bb85397`. Reviewed object: the current uncommitted T017 CI-remediation diff. Reviewer: DeepSeek V4 Pro. Final verdict: `PASS FOR BASELINE REVIEW`.

This verdict allows the fixes to proceed as a new verifier-baseline candidate. It is not proof of T017 completion, Go 1.22.12 or GitHub Windows/Linux CI success, or release authorization.

## Findings

No reproducible `BLOCKER`, `HIGH`, or `MEDIUM` was identified.

1. `LOW`: the read from `waitDone` in `internal/guard/process_test.go` has no independent timeout. If a future platform implementation reports the identity gone while the process object has not exited, the test can wait until its global timeout. The current `StopProcessIdentity` path calls `waitIdentityGone` before success, so this does not block the baseline. A timeout-select may harden the test later.

## Previous Blocker Recheck

The previous `HIGH` is fixed: the child-process `Wait` result no longer matches platform-specific error text. It uses `errors.As(waitErr, &exitErr)` to accept `*exec.ExitError`, covering Unix signal termination and Windows non-zero exits without broadening production success semantics.

## Actual Review Scope

The reviewer confirmed exactly nine changed files relative to `12fa05d1f56c5e9a03a0964fce9d9f842bb85397`:

1. `internal/guard/process.go`
2. `internal/guard/process_test.go`
3. `internal/release/probe_test.go`
4. `internal/release/verify_test.go`
5. `internal/snapshot/export.go`
6. `docs/validation/s1-selfhost.md`
7. `docs/validation/s1-selfhost_EN.md`
8. `docs/DEV_PLAN.md`
9. `docs/DEV_PLAN_EN.md`

There were no additional changes at review time; `sessbridge` was unchanged and `tmp/` contained only `.gitkeep`.

## Local Validation Evidence

1. `git diff --check`, `git diff --stat`, and `git diff --name-status`: exit 0; nine files, 79 insertions / 54 deletions.
2. `gofmt -l`: exit 0 with no output.
3. Encoding gate: all nine files passed; Markdown used UTF-8 BOM + LF and Go used UTF-8 without BOM + LF.
4. `CGO_ENABLED=0 GOTOOLCHAIN=local go build ./...`: exit 0.
5. `CGO_ENABLED=0 GOTOOLCHAIN=local go vet ./...`: exit 0.
6. `CGO_ENABLED=0 GOTOOLCHAIN=local go test -count=1 -timeout=2m ./...`: exit 0; all packages passed.
7. The guard stop test repeated 10 times, the two release remediation tests repeated three times, and the snapshot overlap test repeated five times: all passed.
8. `GOOS=linux GOARCH=amd64 go test -c ./internal/guard`: cross-compilation passed and the temporary artifact was removed.
9. Editor diagnostics for all nine files: no errors.

## Unverified Items and Remaining Gates

1. The local toolchain was Go 1.27.0. Tests did not run under Go 1.22.12, and Linux received cross-compilation rather than native execution.
2. The remediation commit has not landed, and GitHub Windows/Linux push CI has not been rerun.
3. GPT-5.3 Codex must independently review the final candidate diff, including this record.
4. Only after both reviews pass may the replacement verifier baseline be committed. After its dual-platform push CI passes, a different candidate SHA must be created and the complete candidate-build/probe/selfhost matrix manually dispatched from `main`.
5. T017 may be marked `COMPLETE` only after archiving the run URL, verifier/candidate SHAs, platform results, and artifact digests.

This review was strictly read-only: it modified no files, performed no Git write, and triggered no GitHub workflow, release, signature, or attestation.
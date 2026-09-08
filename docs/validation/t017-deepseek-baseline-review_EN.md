# T017 Verifier Baseline DeepSeek Independent Review Record

[中文](t017-deepseek-baseline-review.md)

Date: 2026-09-09. Verdict: `PASS FOR BASELINE REVIEW`. No reproducible `BLOCKER` or `HIGH` was found. Landing an independent verifier baseline, using a different candidate SHA, and completing the first real Windows/Linux GitHub CI run remain pending.

## Review Source and Boundary

1. This page archives the DeepSeek V4 Pro report manually recovered by the user from a collapsed chat window. Its two tables were reconstructed from the copied rows and columns; the remaining content preserves the report's meaning.
2. This review applied the same strict read-only requirements as the [GPT-5.3 Codex verifier-baseline independent review](t017-baseline-review_EN.md) and returned the same `PASS FOR BASELINE REVIEW` verdict.
3. The matching verdicts corroborate that the reviewed baseline received independent model review. Both remain conversational review records, not signed receipts, attestations, or independently authenticated model-identity evidence.
4. Code locations in the report identify the reviewed snapshot; later file changes may shift line numbers.

## Findings

**No BLOCKER or HIGH finding was identified.**

1. `MEDIUM` (known and disclosed trust boundary): script-digest gates detect content drift; they do not prove PowerShell semantics. If an attacker changes both a workflow `run` body and its expected test digest, local tests can become green again.

	- Evidence: `workflow_test.go:234`, `workflow_test.go:241`, item 4 under “Review Remediation Boundaries” in `s1-selfhost_EN.md`, and finding 1 in `t017-baseline-review_EN.md`.
	- Why existing tests do not prevent it: this is an inherent structural property of the gate; same-source tests cannot authenticate themselves.
	- Mitigation already implemented and recorded: digests live in independent constants, automatic refresh from the workflow under test is prohibited, and the baseline requires independent review before landing. This remains non-blocking.

2. `LOW`: the candidate's `actual` result is self-reported by the candidate process. `candidate-probe` must execute the candidate binary; a malicious candidate can emit conforming JSON or even rewrite `actual.json` on the runner. Equality with the external oracle detects **drift**, not malicious-candidate behavior.

	- Evidence: `ci.yml#L136` (the probe executes the candidate), `probe.go:47`, and `verify.go:82` (equality comparison only).
	- Why existing tests do not prevent it: a black-box probe cannot eliminate self-attestation; the `fake-oracle-write` mutation covers the script layer, not behavior inside the candidate process.
	- Recommendation (later enhancement, non-blocking): state in `s1-selfhost*.md` that actual is candidate-reported and the external oracle prevents drift rather than proving a candidate non-malicious. Long-term controls depend on T018/S2 process isolation and observability.

3. `LOW`: local Go 1.27.0 differs from CI-pinned Go 1.22.12. Local success cannot replace results from the target runners.

	- Evidence: `GO_VERSION: "1.22.12"` near `ci.yml#L48`; local version `go1.27.0 windows/amd64`.

4. `LOW`: production offline enforcement is an import/string static scan, not an OS network sandbox; its command vocabulary is not exhaustive.

	- Evidence: `network_test.go:18` and `policy.go:193` (`requireCapabilities` returns `ErrPolicyUnavailable` for deny/loopback/allowlist when capability is absent). The documentation already states this accurately in implementation-scope item 8 of `s1-selfhost_EN.md`; this remains non-blocking.

## Offline Gate Exit Codes (All 0)

| Command | Exit code |
|---|---:|
| `go build ./...` | 0 |
| `go vet ./...` | 0 |
| `go test -count=1 ./...` | 0 (all packages ok) |
| Focused workflow contract and mutation tests | 0 |
| `git diff --check` | 0 |

## Mutation Subtests: 540 Passed, by Group

| Group | Passed |
|---|---:|
| `TestWorkflowActionMutationsAreRejected` | 208 |
| `TestWorkflowActionParametersAreLocked` | 78 |
| `TestWorkflowScriptMutationsAreRejected` | 144 |
| `TestWorkflowScriptsRejectCommentedCommands` | 12 |
| `TestWorkflowStructureMutationsAreRejected` | 82 |
| `TestWorkflowYAMLMutationsAreRejected` | 16 |

## Focused Review Cross-Check

- **A. Complete-script gate: passed.** Contract and mutation tests share `decodeWorkflow` (`workflow_test.go:79`) and `validateWorkflowContract` (`workflow_test.go:147`). Twelve digests are independent constants (`workflow_test.go:241`). Comments, here-strings, uninvoked script blocks, unreachable branches, `exit 0`, variable substitution, function shadowing, and forged oracle writes are rejected. Digests are also bound to job/step name, order, `shell`, `if`, and working directory. The documentation accurately limits the “change both and become green” boundary through independent review.
- **B. Actions and workflow structure: passed.** Eleven Actions are locked by job/role to repository plus full SHA (`workflow_test.go:303`), with every `with` parameter locked (`workflow_test.go:317`). All eight inputs are required strings without defaults. Permissions, toolchain, matrix, explicit `fail-fast: false`, needs, and step order are validated. YAML node checks reject unknown/duplicate fields, multiple documents, nulls, anchors, aliases, and merges while preserving field presence.
- **C. Three-runner trust chain: passed.** candidate-build only builds and seals the candidate. probe builds its harness from verifier-source and uploads only an untrusted oracle. selfhost uses a fresh runner to rebuild verifier/seed, downloads the same binary/oracle, and never executes the candidate; candidate-metadata checkout uses only `git show --format=%cI` to read commit time. Artifact names bind `runner.os` plus candidate SHA. All three trust jobs use `cache: false` plus `GOTOOLCHAIN=local`, preventing untrusted state from entering verifier builds through the Go cache.
- **D. Verifier and release evidence: passed.** `VerifySelfHost` binds seed/candidate/oracle ownership and content. realpath, symlink, directory overlap, path traversal, trailing JSON, and duplicate cases fail closed. SHA256SUMS completely covers artifacts and rejects extras, omissions, and tampering. The SBOM binds the source commit and matches the license manifest module by module. `PublisherIdentityAuthenticated: false` is constant output (`verify.go:202`). The bootstrap lacks the new verifier, and the process explicitly requires independent review and landing before dispatch against a different candidate.
- **E. Runtime offline behavior: passed.** An independent scan found zero `net` imports. All five `os/exec` uses are approved boundaries. All three network modes fail closed before execution when capability is absent (`network_test.go:77`), and the static scan is not represented as sandbox proof.

## Report Items

- Changed Go/YAML/JSON files used UTF-8 without BOM + LF, while Markdown used UTF-8 BOM + LF: **all 30 changed files passed**.
- `tmp/`: only `.gitkeep`; no residue.
- Local Go version: `go1.27.0 windows/amd64`, **different from** CI-pinned 1.22.12.
- First real GitHub CI: **not run**. Independent verifier baseline: **not landed** (HEAD remained `fc99f553...`; no commit/push/tag/dispatch). Real-runner attack exercise: **not run**.
- T017 status: `IMPLEMENTED_AWAITING_CI`; DEV_PLAN remained unchecked, accurately.

This was a strict read-only review: it modified no files, performed no Git write, triggered no GitHub workflow, and did not touch sessbridge.
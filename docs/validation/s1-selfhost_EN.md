# T017 Self-Host and Trusted Release Validation Report

[中文](s1-selfhost.md)

Date: 2026-09-08. Verdict: `IMPLEMENTED_AWAITING_CI`. Local Windows validation passed, and both the [GPT-5.3 Codex review](t017-baseline-review_EN.md) and the [DeepSeek V4 Pro review](t017-deepseek-baseline-review_EN.md) concluded `PASS FOR BASELINE REVIEW`. T017 requires the first successful native Windows/Linux GitHub CI matrix before it can be marked `COMPLETE`.

## Implementation Scope

1. `.github/workflows/ci.yml`: pins Go 1.22.12, `CGO_ENABLED=0`, `GOTOOLCHAIN=local`, and GitHub Actions by full commit SHA; it runs build/vet/test on Windows and Linux and runs the self-host matrix only through an explicit maintainer `workflow_dispatch` from `main`, building the fixed seed, candidate, and trusted verifier in independent roots. All three trust jobs explicitly set setup-go `cache: false`, preventing Go cache restore/save from carrying build state between runners.
2. `internal/release/` and `cmd/prfrail-release/`: verify a full 40-character commit SHA, unchanged seed, path isolation, external oracle, SHA256SUMS, SPDX 2.3 SBOM, and licenses for modules actually linked into the binary.
3. `testdata/selfhost/`: fixes bootstrap commit `fc99f553c3066b24230a093d14bc17af0c6a398e`, four external oracle cases, noop/executable chains, and the license policy.
4. `cmd/prfrail/main.go`: makes the version link-injectable so CI binds the candidate to its fixed commit.

## AT-13/14 Evidence

1. Seed and candidate use separate source, build, and run roots. Real-path overlap, reversed oracle ownership, or seed mutation fails closed.
2. The external oracle fixes `version`, `validate-noop`, `run-noop`, and `reject-executable`; the candidate produces only actual results and cannot generate expected outcomes. Actual results are self-reported by the candidate process, and equality with the external oracle detects behavioral drift rather than proving the candidate non-malicious. At dispatch, the maintainer independently supplies candidate, verifier, and bootstrap commits plus five fixture SHA-256 pins; the workflow has no candidate-controlled defaults.
3. The verifier commit must equal the `main` commit containing the reviewed workflow, while the candidate must be a different full commit SHA. The `candidate-build` job only builds and seals the candidate binary; a fresh `candidate-probe` runner downloads that immutable artifact, executes the probe, and uploads only an untrusted oracle. The final `selfhost` runner downloads the same binary and oracle, builds its tool from verifier-source, never executes the candidate, and uses oracle/policy files from the verifier commit for `hash-tree/selfhost/generate/verify`.
4. The local two-generation drill built the seed from the fixed commit and the candidate from the current tree. All four oracle cases matched, with seed, candidate, and oracle hashes independently recomputed.
5. Release metadata comes from Go binary build info. The SPDX document is bound to the source commit and must match the license manifest module by module on path, version, and license. Unapproved licenses, empty or inconsistent SPDX package lists, symbolic revisions, missing/tampered artifacts, symlinks, and incomplete SHA256SUMS all fail closed.
6. `PublisherIdentityAuthenticated=false` is invariant. SHA256SUMS proves content agreement with the manifest, not publisher identity; signing/attestation is a later enhancement.
7. The GitHub workflow uploads only a CI artifact; it creates no Release, signature, commit, or push. CI checkout, setup-go, npm, and artifact upload are GitHub-hosted build operations, not ProofRail runtime behavior.
8. Regression tests prevent production Go files under `cmd/` and `internal/` from directly importing `net` or its subpackages or hard-coding common downloaders/networking Git subcommands. The current process runner advertises no network-isolation capability, so every network mode fails closed before execution. The SPDX namespace is an offline identifier and causes no network request.

## Review Remediation Boundaries

1. Script-safety assertions based on `strings.Contains` / `strings.Index` have been removed. All 12 `run` bodies are bound by job/step to full SHA-256 values in independent constants, alongside `shell`, `if`, working-directory, and mutually exclusive Action/script fields. Commented commands, quoted text, uninvoked script blocks, early successful exits, variable substitutions, or appended fake-oracle writes cannot preserve the original digest.
2. The workflow contract test and in-memory mutation tests share `decodeWorkflow` and `validateWorkflowContract`. Structural checks cover every job's input mappings, native matrix, dependencies, and step order. Action checks cover exact repositories, SHAs, all `with` parameters, and execution conditions, including test checkout `persist-credentials: false`, Setup Go version/cache policy, and the final artifact upload.
3. YAML decoding rejects unknown/duplicate fields, multiple documents, explicit defaults, aliases/anchors, and unapproved nulls. Raw field presence is retained, so combining `uses` with an empty `run`/`shell` cannot hide a mixed step. Negative tests mutate only in-memory YAML/structures; they neither execute replacement scripts nor modify the real workflow.
4. Fixed digests detect content drift; they do not prove PowerShell semantics, publisher identity, or runtime sandboxing. Even comment-only script changes require reviewing the complete script and execution context before explicitly updating baseline constants. Expected digests must never refresh automatically from the workflow under test. Changes to the tests themselves still require independent review.
5. This remediation does not rerun the two-generation drill or establish success on GitHub CI or real-runner attack exercises. The remediation passed both independent reviews but still needs to land as a new verifier baseline and pass subsequent CI; T017 remains `IMPLEMENTED_AWAITING_CI`.

## Current Gate Results

1. `CGO_ENABLED=0 go build ./...`: passed.
2. `CGO_ENABLED=0 go vet ./...`: passed.
3. `CGO_ENABLED=0 go test -count=1 ./...`: passed; `internal/release` includes focused tests for external trust pins, direct-network prohibition, and release evidence.
4. `cd tools/contracts && npm test`: passed, 2/2 tests with all 82 independent fixtures matching.
5. Native Windows 11 complete two-generation drill: passed; fixed-seed hashing, external oracle, SPDX/license/SHA256SUMS verification succeeded, and temporary worktrees/artifacts were removed.
6. GitHub push CI run `34253473081`: the first run on 2026-09-09 failed both Windows and Linux `test` jobs and is not AT-14 evidence. It exposed Go 1.22/native-runner differences: the process-stop test did not promptly reap its child, the Unix probe test reused the candidate-binary path as an output directory, the license test used a test binary whose main-module path was empty, and Windows overlap comparison did not resolve the existing parent of an absent destination consistently.
7. The fixes pass focused and full tests with the current local toolchain. Both the [DeepSeek V4 Pro independent review](t017-ci-remediation-deepseek-review_EN.md) and the [GPT-5.3 Codex independent review](t017-ci-remediation-codex-review_EN.md) returned `PASS FOR BASELINE REVIEW`. The fixes have not landed as a replacement verifier baseline, and its push CI and first complete maintainer-dispatched Windows/Linux matrix remain pending. AT-14 and T017 remain incomplete.

## First CI Inputs

The maintainer must calculate these required inputs from the independently reviewed commit/fixtures and enter them manually; candidate-generated values are not trusted. The workflow must be dispatched from `main`; the self-host job is rejected for every other ref:

| Input | Current reviewed value |
|---|---|
| `candidate_commit` | Create after the corrected verifier baseline lands; must be a full SHA different from the verifier commit |
| `verifier_commit` | Pending renewed independent review and landing; `12fa05d1f56c5e9a03a0964fce9d9f842bb85397` cannot be the final green baseline because its push CI failed |
| `bootstrap_commit` | `fc99f553c3066b24230a093d14bc17af0c6a398e` |
| `bootstrap_manifest_sha256` | `ae0d5db85d6a45de19dc75d40460705e0a7a0dea20de2a11632c25bf8b4543fc` |
| `expected_oracle_sha256` | `388f88842710abc90676ed0a85f4bd935b594fc552146577c304b0bbbd850aaf` |
| `license_policy_sha256` | `ff9d243c67d516a20147cb5b8684741c3f5ddb77a761af7aa1aab1df4a398013` |
| `noop_chain_sha256` | `fbeee2b64ab174637c77d9cdd1b04bd7fd6c681912f0b73f7a4847bc19a7b31c` |
| `executable_chain_sha256` | `dbc5fc415917398c88a6ed669971c417ed8ee6884a2c8ebfc1271807d757935b` |

## Boundaries and Follow-Up

1. `12fa05d1f56c5e9a03a0964fce9d9f842bb85397` completed independent review and its first landing, but its push CI failed. Because the fixes change verifier-related code, they require renewed independent review and a replacement green verifier baseline; that earlier verifier commit then evaluates a different candidate commit. This one-time trust bootstrap does not substitute candidate execution for human review.
2. This implementation publishes no install package and authorizes no `git commit`, `git push`, tag, signature, or GitHub Release.
3. The fixed seed is read-only; failures discard the candidate and never overwrite or roll back the seed.
4. Both the replacement verifier baseline's push CI and the later complete maintainer-dispatched matrix must pass on both platforms. Only after recording verifier/candidate/run evidence may T017 be marked `COMPLETE`; any failure keeps the gate closed.
5. T018, installation/upgrade/uninstall exercises, and final S1 exit remain incomplete.
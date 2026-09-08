# T017 Self-Host and Trusted Release Validation Report

[中文](s1-selfhost.md)

Date: 2026-09-08; updated 2026-09-09. Verdict: `COMPLETE`. The remediation passed both independent reviews and the verifier baseline's Windows/Linux push CI. The complete Windows/Linux matrix manually dispatched from `main` against a different candidate SHA also passed, and the AT-13/14 evidence is archived below.

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
5. This remediation does not rerun the local two-generation drill or establish success for real-runner attack exercises; passing CI does not expand those boundaries. The remediation passed both independent reviews, landed, and passed the verifier baseline's dual-platform push CI, followed by the complete manually dispatched matrix against a different candidate SHA; T017 is `COMPLETE`.

## Current Gate Results

1. `CGO_ENABLED=0 go build ./...`: passed.
2. `CGO_ENABLED=0 go vet ./...`: passed.
3. `CGO_ENABLED=0 go test -count=1 ./...`: passed; `internal/release` includes focused tests for external trust pins, direct-network prohibition, and release evidence.
4. `cd tools/contracts && npm test`: passed, 2/2 tests with all 82 independent fixtures matching.
5. Native Windows 11 complete two-generation drill: passed; fixed-seed hashing, external oracle, SPDX/license/SHA256SUMS verification succeeded, and temporary worktrees/artifacts were removed.
6. GitHub push CI run `34253473081`: the first run on 2026-09-09 failed both Windows and Linux `test` jobs and is not AT-14 evidence. It exposed Go 1.22/native-runner differences: the process-stop test did not promptly reap its child, the Unix probe test reused the candidate-binary path as an output directory, the license test used a test binary whose main-module path was empty, and Windows overlap comparison did not resolve the existing parent of an absent destination consistently.
7. The fixes pass focused and full tests with the current local toolchain. Both the [DeepSeek V4 Pro independent review](t017-ci-remediation-deepseek-review_EN.md) and the [GPT-5.3 Codex independent review](t017-ci-remediation-codex-review_EN.md) returned `PASS FOR BASELINE REVIEW`. The fixes landed as verifier baseline `18c9340c60deb30049554491f8905df51d0a3180`; both the [Ubuntu job](https://github.com/larsonzh/prfrail/actions/runs/34268076912/job/102202411754) and [Windows job](https://github.com/larsonzh/prfrail/actions/runs/34268076912/job/102202412040) succeeded in GitHub push CI [run 34268076912](https://github.com/larsonzh/prfrail/actions/runs/34268076912).
8. Maintainer `larsonzh` manually triggered GitHub Actions [run 34277671704](https://github.com/larsonzh/prfrail/actions/runs/34277671704) (CI #3, `workflow_dispatch`) from `main`. The run was bound to verifier `18c9340c60deb30049554491f8905df51d0a3180` and supplied candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c`, which are distinct. The run concluded `success`, with all 8 jobs successful: [Go Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603547), [Go Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603526), [candidate build Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603542), [candidate build Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603184), [candidate probe Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234809504), [candidate probe Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234809369), [bootstrap/release Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102235086119), and [bootstrap/release Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102235086172).

## First Manual-CI Artifacts

The public API for run `34277671704` reports these six unexpired GitHub Actions artifacts. Each digest is GitHub's digest of the artifact archive; it does not replace the verifier-checked SHA256SUMS inside the archive:

| Artifact | Size (bytes) | GitHub archive digest |
|---|---:|---|
| `candidate-binary-Linux-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,247,948 | `sha256:7e4d569e6754809ce3da97cdebcfad1db058eef9b51578f3dfc9cb56601b8d0e` |
| `candidate-binary-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,397,137 | `sha256:ab8edb01ddb7dafc54ee78bf6942b755fffe5437ec185617416c90c6d4983e66` |
| `candidate-oracle-Linux-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 261 | `sha256:44465973567919c34f39ba869a1b8030fa9611a0107d98ece948e6bb37195caf` |
| `candidate-oracle-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 261 | `sha256:fdecda0b17ee81a782f5d348f3e697b28ff56f11b2bb24f31ad09f5ae57b58d1` |
| `prfrail-Linux-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,249,050 | `sha256:43a99ee59076e3eae96c9ed61016159511f9ac9309169f170befead286d1d41a` |
| `prfrail-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,398,244 | `sha256:7adc8f21c81df88db31d18ea5b83e892d5e7d3a9416088c51b0ec12b15df6ca7` |

## First CI Inputs

The maintainer must calculate these required inputs from the independently reviewed commit/fixtures and enter them manually; candidate-generated values are not trusted. The workflow must be dispatched from `main`; the self-host job is rejected for every other ref:

| Input | Current reviewed value |
|---|---|
| `candidate_commit` | `1e7af676e8e84028c7ecfef1ccde91213738c20c`; a docs-only candidate distinct from and directly based on the verifier baseline |
| `verifier_commit` | `18c9340c60deb30049554491f8905df51d0a3180`; both independent reviews passed and both Windows/Linux jobs succeeded in push CI run `34268076912` |
| `bootstrap_commit` | `fc99f553c3066b24230a093d14bc17af0c6a398e` |
| `bootstrap_manifest_sha256` | `ae0d5db85d6a45de19dc75d40460705e0a7a0dea20de2a11632c25bf8b4543fc` |
| `expected_oracle_sha256` | `388f88842710abc90676ed0a85f4bd935b594fc552146577c304b0bbbd850aaf` |
| `license_policy_sha256` | `ff9d243c67d516a20147cb5b8684741c3f5ddb77a761af7aa1aab1df4a398013` |
| `noop_chain_sha256` | `fbeee2b64ab174637c77d9cdd1b04bd7fd6c681912f0b73f7a4847bc19a7b31c` |
| `executable_chain_sha256` | `dbc5fc415917398c88a6ed669971c417ed8ee6884a2c8ebfc1271807d757935b` |

## Boundaries and Follow-Up

1. After push CI failed for `12fa05d1f56c5e9a03a0964fce9d9f842bb85397`, the remediation passed both independent reviews and established `18c9340c60deb30049554491f8905df51d0a3180` as the new green verifier baseline. That earlier verifier commit must then evaluate a different candidate commit. This one-time trust bootstrap does not substitute candidate execution for human review.
2. This implementation publishes no install package and authorizes no `git commit`, `git push`, tag, signature, or GitHub Release.
3. The fixed seed is read-only; failures discard the candidate and never overwrite or roll back the seed.
4. Both the replacement verifier baseline's push CI and the later complete maintainer-dispatched matrix passed on both platforms. Verifier/candidate/run and artifact evidence is recorded, so T017 is marked `COMPLETE`.
5. T018, installation/upgrade/uninstall exercises, and final S1 exit remain incomplete.
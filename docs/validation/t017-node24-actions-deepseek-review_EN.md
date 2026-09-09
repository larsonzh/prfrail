# T017 Node 24 Actions DeepSeek V4 Pro Independent Review Record

[中文](t017-node24-actions-deepseek-review.md)

Date: 2026-09-09. Historical verifier: `18c9340c60deb30049554491f8905df51d0a3180`. Reviewed new verifier: `6024526893fd7acc47de44ee862572ac9ec33303`. Reviewer: DeepSeek V4 Pro. Final verdict: `PASS FOR NEW VERIFIER BASELINE`.

This verdict approves the commit as the new verifier baseline after the Node 24 Actions upgrade. It does not establish completion of the new baseline's T017 runtime matrix, authorize a release, authenticate a publisher, or rewrite the historical T017 `COMPLETE` evidence.

## Findings

No blocking `Critical`, `High`, or `Medium` finding was identified, and there was no blocking `Low` finding.

Non-blocking observations:

1. The ordinary `test` job uses setup-go with `cache: true`, while all three trust jobs use `cache: false`. This is an existing isolation design locked by the contract tests.
2. upload-artifact v7 defaults to `overwrite: false`. Each of the three upload names contains the platform and candidate SHA, preventing collisions within a run.
3. Artifact downloads do not preserve executable permissions. The Linux probe explicitly runs `chmod +x` before executing the candidate, so the existing mitigation remains in place.

## Independent Verification Conclusions

1. The range `18c9340..6024526` contains only two commits: `971417e`, which modifies four T017 completion-evidence documents, and `6024526`, which modifies only `.github/workflows/ci.yml`, `internal/release/workflow_test.go`, and `internal/release/workflow_mutation_test.go`.
2. The substantive changes in `6024526` are limited to version/full-SHA upgrades for four Actions, `digest-mismatch: error` on three download steps, and corresponding contract-test synchronization. The trust topology, triggers, permissions, and script bodies are unchanged.
3. Independent checks of the official release/tag/commit/action.yml sources confirm that the following full SHAs belong to the corresponding version tags in official `actions/*` repositories, are identified as verified commits on the release pages, and declare `runs.using: node24` in `action.yml`:

| Action | Version | Pinned SHA |
|---|---|---|
| `actions/checkout` | `v7.0.1` | `3d3c42e5aac5ba805825da76410c181273ba90b1` |
| `actions/setup-go` | `v7.0.0` | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |
| `actions/upload-artifact` | `v7.0.1` | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` |
| `actions/download-artifact` | `v8.0.1` | `3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` |

4. The breaking changes do not adversely affect this workflow: checkout v7.0.1 fixes and dependency updates do not change the used inputs; setup-go v7 ESM/cache dependency upgrades do not change `go-version`/`cache` semantics; upload-artifact v7 adds direct uploads through `archive`, which this workflow does not enable, so uploads remain compressed by default; download-artifact v8 uses Content-Type to control decompression and changes digest mismatch handling to default to `error`, strengthening fail-closed behavior. GitHub-hosted runners are unaffected by the v7 minimum runner-version requirement for self-hosted runners.
5. Valid download-artifact v8 `digest-mismatch` values are `ignore`, `info`, `warn`, and `error`, with `error` as the default. All three trusted download points explicitly select `error`. Upload v7 records the server-side artifact digest and download v8 validates against that digest, so the pairing is compatible.
6. Verifier equality with `GITHUB_SHA`, candidate inequality with verifier, main-only dispatch, the candidate-build → candidate-probe → selfhost dependency chain, the dual-platform matrix, `contents: read`, `persist-credentials: false`, `cache: false` in trust jobs, the external oracle, and fixed fixture SHA-256 values all remain unchanged.
7. Contract tests continue to lock Action names, full SHAs, step roles, parameters, permissions, triggers, inputs, job dependencies, the dual-platform matrix, and script digests. They reject foreign Actions, changed SHAs, role swaps, parameter deletion/modification, skipped steps, early success, forged oracles, and YAML anchor/merge/null/duplicate/trailing-document mutations. The synchronized changes did not remove or weaken existing negative cases.

Official sources:

1. <https://github.com/actions/checkout/releases/tag/v7.0.1>
2. <https://github.com/actions/setup-go/releases/tag/v7.0.0>
3. <https://github.com/actions/upload-artifact/releases/tag/v7.0.1>
4. <https://github.com/actions/download-artifact/releases/tag/v8.0.1>
5. <https://github.com/actions/upload-artifact/releases/tag/v7.0.0>
6. <https://github.com/actions/download-artifact/releases/tag/v8.0.0>

## Existing CI Evidence Boundary

Ordinary push CI run `34314888161` succeeded. It proves only that the Node 24 Action combination can complete the ordinary build, vet, test, and contract-fixture checks on GitHub-hosted Windows/Linux runners. The three trust jobs were excluded by their dispatch conditions, so this run does not establish completion of the T017 trust matrix.

## Remaining Runtime Evidence

1. Create a docs-only candidate with a SHA different from the new verifier. The candidate must not modify the reviewed workflow, release verifier, or its contract/mutation tests.
2. Manually dispatch the complete workflow from `main` and confirm all eight jobs succeed: test, candidate-build, candidate-probe, and selfhost on both Ubuntu and Windows.
3. Archive the run URL, verifier/candidate SHAs, each job result, the names/sizes/digests of all six artifacts, and the dispatch-input digests.
4. Until the matrix succeeds and its evidence is archived, the new verifier baseline's T017 runtime validation must not be described as `COMPLETE`.

This review was strictly read-only: it modified or created no file, made no commit or push, and triggered no workflow, Release, signature, or attestation. The historical T017 `COMPLETE` evidence for verifier `18c9340c60deb30049554491f8905df51d0a3180`, candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c`, and run `34277671704` remains valid.
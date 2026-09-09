# T017 Node 24 Actions GPT-5.3 Codex Independent Review Record

[中文](t017-node24-actions-codex-review.md)

Date: 2026-09-09. Historical verifier: `18c9340c60deb30049554491f8905df51d0a3180`. Reviewed new verifier: `6024526893fd7acc47de44ee862572ac9ec33303`. Reviewer: GPT-5.3 Codex. Final verdict: `PASS FOR NEW VERIFIER BASELINE`.

This verdict approves `6024526893fd7acc47de44ee862572ac9ec33303` as the new verifier baseline after the Node 24 Actions upgrade. It does not invalidate the historical T017 `COMPLETE` result, establish completion of the new baseline's T017 runtime matrix, authorize a release, or authenticate a publisher.

## Findings

No reproducible blocking `Critical`, `High`, `Medium`, or `Low` finding was identified.

1. Verifier/candidate SHA isolation, explicit `workflow_dispatch` from `main`, native dual-platform runners, read-only repository permissions, and disabled checkout credential persistence remain unchanged.
2. The candidate-build, candidate-probe, and selfhost dependency chain is unchanged. The trusted verifier, fixed bootstrap seed, candidate metadata, and candidate binary remain in independent directories.
3. Failure propagation through the external oracle, fixed fixture SHA-256 values, seed hash, selfhost operation, release metadata generation, and final verification is not weakened.
4. All three trusted artifact download points explicitly set `digest-mismatch: error`, matching the fail-closed digest verification semantics of `actions/download-artifact@v8.0.1`.
5. `workflow_test.go` locks every Action SHA and parameter. `workflow_mutation_test.go` continues to reject Action, script, structure, parameter, and YAML bypass mutations. The test synchronization did not remove or weaken existing negative cases.

## Action Supply-Chain Verification

The following official GitHub Actions are pinned by full commit SHA. Each official release page identifies the SHA as the verified commit for that version tag, and each corresponding `action.yml` declares `runs.using: node24`:

| Action | Version | Pinned SHA |
|---|---|---|
| `actions/checkout` | `v7.0.1` | `3d3c42e5aac5ba805825da76410c181273ba90b1` |
| `actions/setup-go` | `v7.0.0` | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |
| `actions/upload-artifact` | `v7.0.1` | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` |
| `actions/download-artifact` | `v8.0.1` | `3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` |

Official sources:

1. <https://github.com/actions/checkout/releases/tag/v7.0.1>
2. <https://github.com/actions/setup-go/releases/tag/v7.0.0>
3. <https://github.com/actions/upload-artifact/releases/tag/v7.0.1>
4. <https://github.com/actions/download-artifact/releases/tag/v8.0.1>

## Review Scope and Evidence Boundary

1. Commit `6024526893fd7acc47de44ee862572ac9ec33303` modifies only `.github/workflows/ci.yml`, `internal/release/workflow_test.go`, and `internal/release/workflow_mutation_test.go`.
2. The workflow's substantive changes are limited to the full-SHA/version upgrades for four Actions and the addition of `digest-mismatch: error` to three download-artifact steps. The contract tests synchronously lock those values.
3. Ordinary push CI run `34314888161` succeeded, demonstrating that the new Action combination completes normal Windows/Linux CI. That run did not execute the complete T017 trust matrix and is not matrix-completion evidence.
4. This was a read-only review: it did not modify code, Git history, or the workflow, and did not trigger a release, signature, or attestation.

## Remaining Runtime Evidence

1. Create a docs-only candidate with a SHA different from the new verifier. The candidate must not modify the reviewed workflow, release verifier, or its contract/mutation tests.
2. Manually dispatch the complete workflow from `main` and confirm all eight jobs succeed: test (Ubuntu/Windows), candidate-build (Ubuntu/Windows), candidate-probe (Ubuntu/Windows), and selfhost (Ubuntu/Windows).
3. Archive the run URL, verifier/candidate SHAs, each job result, the names/sizes/digests of all six artifacts, and the dispatch-input digests.
4. Until the full matrix succeeds and its evidence is archived, the new verifier baseline's T017 runtime validation must not be described as `COMPLETE`.

This record captures only the review verdict. The historical T017 `COMPLETE` evidence for verifier `18c9340c60deb30049554491f8905df51d0a3180`, candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c`, and run `34277671704` remains valid and unchanged.
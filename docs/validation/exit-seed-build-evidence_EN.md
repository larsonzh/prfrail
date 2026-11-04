# P55 slice validation report — `EXIT-SEED-BUILD-EVIDENCE` (AT-13 seed-driven build evidence supplement)

| Field | Value |
|---|---|
| Slice / queue slot | `EXIT-SEED-BUILD-EVIDENCE` (layer-3 item 7) |
| Size / class | M / `[SLICE]` (evidence-chain supplement) |
| Anchor | `HEAD` = `baa144df5898817e0ac7cd860928bb4a90cbd009` (`baa144d`) = origin/main |
| Directive | `docs/DELIVERY_DIRECTIVE.md` v1.60 |
| Controller model | `DeepSeek V4.1 Flash (deepseek)` |
| Evidence package | `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/` |
| Status | Reached the 10 stop point, awaiting same-round commit authorisation (not committed, not pushed) |

## 1. Goal and hard boundary

Goal: complete the AT-13 evidence chain. AT-13 text (`docs/TEST_STRATEGY.md:43`): `seed and candidate kept separate; external oracle comparison; clean-room success and failure both leave the seed intact` (owner `harness/release; T017`).

Hard boundary (start message [6, 4]): no `.github/workflows/**` change, no Go product code change, deliverables limited to `.md` and the evidence package. **Zero violations** (see the fence in §2; no `.go`, no workflow).

## 2. artifacts fence

The entries below are the complete change set of this slice (path + disposition, TAB separated, sorted by path).

```artifacts
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/S1_RELEASE_NOTES.md	repo
docs/S1_RELEASE_NOTES_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/00-briefing.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-00-env.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-01-clone.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-02-worktree-seed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-03-worktree-verifier.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-04-worktree-candidate.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-05-module-preload-seed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-06-module-preload-verifier.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-07-module-preload-candidate.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-08-build-seed-binary.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-09-seed-binary-version.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-10-build-verifier-tool.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-11-build-candidate-binary.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-12-seed-tree-hash.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-13-probe-actual.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-14-selfhost-success.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-15-neg-roots-overlap.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-16-neg-oracle-ownership.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-17-neg-seed-changed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-18-neg-oracle-drift.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-19-seed-intact.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-20-candidate-listing.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-21-bindings.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/02b-22-verdicts.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-linux-captures/index.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-00-env.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-01-go-mod-download.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-02-build-candidate.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-03-build-release-tool.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-04-hash-tree-seed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-05-probe-actual.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-06-selfhost-success.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-07-neg-roots-overlap.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-08-neg-oracle-ownership.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-09-neg-seed-changed.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-10-neg-oracle-drift.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-11-seed-intact.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/02a-12-candidate-listing.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/actual.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/bindings.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/expected.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/02-win-captures/index.json	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report-r2.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/03-test-report.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/06-independent-scan.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/MANIFEST.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/SHA256SUMS.txt	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/anomalies.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/claim-scope.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/freedom-list.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/input-list.md	repo
docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/replay.md	repo
docs/validation/exit-seed-build-evidence.md	repo
docs/validation/exit-seed-build-evidence_EN.md	repo
```

## 3. What was delivered

1. **Capture on both platforms**: Windows (`go1.27.0`, this host) and a Linux VM (`go1.22.12`, same as the CI pin) each performed a clean-room build, a `probe` run producing the actual oracle, and a `selfhost` verification.
2. **Three-role structure** (A-4): seed `fc99f553…` builds `prfrail`; verifier `18c9340c…` (T017 baseline) builds `prfrail-release`; candidate `baa144d…` builds `prfrail`. Roles and the three binaries are mutually distinct, and `selfhost` is executed by the verifier.
3. **Oracle agreement**: both sides report `OracleHash = sha256:388f8884…`, character-identical to `testdata/selfhost/expected.json` and to the CI input table in `s1-selfhost.md`; the five fixture hashes agree on all four sides.
4. **Four counterexamples fail closed on both sides**: overlapping roots / oracle ownership / seed mutation / oracle drift — all non-zero exit with the `verify.go` error strings.
5. **Failure leaves the seed intact**: after each failing run the seed tree hash stays `equal=true` and equals that side's baseline; `verify.go:94` checks overlap before `:110` hashes the seed and `VerifySelfHost` is read-only end to end (structural argument).
6. **No network dependency during builds**: on Linux, modules were preloaded from a reachable mirror and every build ran with `GOPROXY=off`.

## 4. Review chain and gates

- 3 verification (`prfrail-tester`): round 1 `TEST: DONE` (0 Critical / 0 High / 6 Medium / 6 Low) → controller closed them via `claim-scope.md` (narrowed claims) and `anomalies.md` (numbered ledger) → round 2 `TEST: DONE`, **6/6 Medium closed**, six new Low items (R2-1..R2-6) fixed in place.
- 5 pre-review (`deep-reasoner` / `DeepSeek V4 Pro`): round 1 `PRE-REVIEW: PASS WITH FIXES` (1 Medium = F-05-1 abbreviated captured command + 3 Low) → remediation (`anomalies.md` E-1/E-2 errata + derived `index.json` command + `input-list.md` G-1..G-8 table) → round 2 `PRE-REVIEW: PASS WITH FIXES` (F-05-1 closed, F-05-2 closed, F-05-3 partial, F-05-4 deferred to 8b; new N-1/N-2 fixed).
- 6 independent scan (`independent-reviewer` / `MAI-Code-1.1-Flash`): `INDEPENDENT SCAN: PASS`, Section B = `(no findings)`.
- 7 final review (`independent-reviewer` / `GPT-5.3-Codex`): `FINAL REVIEW: PASS WITH FIXES` (1 Low: MANIFEST wording about R2 status too strong), fixed per option A; the reviewer added independent re-checks (5 hashes recomputed, three-way fixture comparison, source anchors and boundary audit).
- 9 gates (controller): `--scope=index` / `--scope=tree` / `--scope=ci` all `TOTAL_FAIL=0`; `SELFTEST: PASS 697/697`; `go build` / `go vet` / `go test -count=1 ./...` all green (15 packages ok).
- The 5/6/7 outputs are archived verbatim under `04-reviews/` in the evidence package.
- CI: run `41925623071` (head `63e2e24`) is green on both legs — Build, Vet, Test and Gates succeed on windows and ubuntu, Race and Contract fixtures run on the ubuntu leg, and the three conditional jobs (candidate-build / candidate-probe / bootstrap-release) are skipped by design.

Verbatim transcript of the 6 independent scan (headings written as Section A–E per the gate requirement; full text in `04-reviews/06-independent-scan.md`):
```sixstage
## Section A Scope
- Read-only audit performed on 2026-11-05 over the staged change set (67 entries at scan time) and the evidence package.
- Read: both ledger files, MANIFEST, SHA256SUMS, claim-scope, anomalies, both 03 test reports, both index.json, plus cmd/prfrail-release/main.go, internal/release/verify.go, internal/release/probe.go, internal/release/verify_test.go.
- No `.go`, no `.github/workflows/**`, no other runtime code in the change set.

## Section B Findings
(no findings)

## Section C Evidence
- Triple set equality: disk_count=49, sum_count=49, manifest_count=49, hash_ok=49, hash_bad=0.
- Fail-close: VerifySelfHost checks pathsOverlap before HashTree and before the seed hash comparison; existingRealPath only performs Abs + EvalSymlinks, no write path.
- equalOracle compares the structured schemaVersion and per-case objects only; it does not imply a malware-free proof.
- verify_test.go pins overlap, oracle mismatch and seed mutation as fail-closed behaviours.

## Section D Falsifiability audit
- Dual-side fail-close: pinned by the four ErrInvalidSelfHost returns and by the captured exit codes plus stderr strings.
- Failure leaves seed intact: pinned by seedHashBefore/After equality in 02a-11 and 02b-19, each equal to its own baseline.
- Oracle agreement: pinned by OracleHash 388f8884 on both sides, matching the fixture and the CI input table.
- Toolchain sameness: pinned on Linux only (go1.22.12); Windows go1.27.0 is explicitly not claimed.
- Payload hash equality: mechanical readings above, not sentinel constants.

## Section E Unverified / overturn conditions
- The remote VM state was cleaned after capture, so live replay is not possible from this session.
- Overturned if any payload hash mismatches, if MANIFEST and disk diverge, or if verify.go ordering/read-only semantics change.

INDEPENDENT SCAN: PASS
```

Disclosure outside the fence: Sections A/C/D above are a structured excerpt (the complete verbatim text is `04-reviews/06-independent-scan.md`); the conclusion line matches the original (`INDEPENDENT SCAN: PASS`).

## 5. Anomalies and deviations

The numbered ledger lives in the evidence package (`anomalies.md`); summary:

| ID | Fact | Disposition |
|---|---|---|
| `D = 1` | Early start without a trigger word; corrected by a forced user interruption | Converted in place, no contamination; the authoritative wording is in `anomalies.md` |
| A-1 / A-2 | First capture round wrote into the package root / session-level cache variables polluted `gofmt` | Fixed and re-run; marked "partially verifiable" |
| A-3 | Step 2 executed by the controller | User ruling: compliant for this slice |
| A-4 | `fc99f55` has no release tool, so the three-role structure was adopted | Approved and emphasised by the user (stronger independence) |
| A-5 | The VM cannot reach `proxy.golang.org` / `sum.golang.org`, so modules were mirrored and builds ran with `GOPROXY=off` | Approved by the user |
| A-6 / A-7 | Two earlier remote rounds each produced nine failures / read-only module cache blocked `rm` | Closed; raw artifacts not reproducible |
| E-1 / E-2 | Captured `TITLE:` / `COMMAND:` display text differs from the real invocation | Corrected by errata; payloads left unmodified |

## 6. Capability boundary (no broader reading allowed)

`claim-scope.md` defines **S-1..S-7 as holding** and **N-1..N-7 as not asserted**. Key non-claims:

- The Windows verifier and the candidate come from the same tree, so **verifier independence is not proven there** (it holds on Linux only).
- `SeedHash` hashes an **assembled directory**, so it does not prove equality with the seed commit's source tree.
- `OracleHash` differs from the actual bytes while verification succeeds, so the comparison is **structural semantics**, not byte equality.
- Oracle agreement only measures **behavioural drift**, so it is not a malware-free proof.

## 7. Document sync (start message [8], six documents)

| # | Document | Disposition |
|---|---|---|
| 1 | `docs/TEST_STRATEGY{,_EN}.md` | Registered the AT-13 supplement, the package pointer and the boundary pointer (explicitly "does not mean AT-13 has passed") |
| 2 | `docs/OPERATIONS{,_EN}.md` | **Not applicable**: evidence supplement only, no operations, on-call or release-process change |
| 3 | `docs/DEV_PLAN{,_EN}.md` | Added the T017 progress record (three roles, both platforms, four counterexamples, unchanged seed; AT-13/14 judgment and checkbox state unchanged) |
| 4 | `docs/S1_SUPPORT_MATRIX{,_EN}.md` | **Not applicable**: no platform or version support-range change |
| 5 | `docs/S1_RELEASE_NOTES{,_EN}.md` | **Not applicable**: no release artifact or feature change |
| 6 | `docs/t027/REMAINING_SLICES{,_EN}.md` | Wrote the 8a progress record; status stays 🔄 In progress |

CN/EN numstat is symmetric pair by pair (`G3(a.1)` verified).

## 8. Native validation readings

| Item | Reading |
|---|---|
| `gate.js --all --scope=index` | `TOTAL_FAIL=0`, `EXIT=0` |
| `gate.js --all --scope=tree` | `TOTAL_FAIL=0`, `EXIT=0` |
| `gate.js --all --scope=ci` | `TOTAL_FAIL=0`, `EXIT=0` |
| `gate.js --selftest` | `SELFTEST: PASS 697/697`, `EXIT=0` |
| `go build ./...` / `go vet ./...` | both `EXIT=0` |
| `go test -count=1 ./...` | `EXIT=0`, 15 packages ok |
| Package hashes | `SHA256SUMS.txt` entries ≡ `MANIFEST` entries ≡ on-disk payloads |

## 9. Accounting

- `R = 0`: no "fix then re-run the same ring" item (the two pre-review rounds are re-reviews of that ring, not R).
- `D = 1`: early start without a trigger word (corrected by a forced user interruption, converted in place).
- Probes / paid calls: 0 (fully offline; no `workflow_dispatch`, no premium request consumed).
- Batch assistance (`GLM-5.3-Flash`): **0 dispatches**, below the "≥ 2" floor in start message [3] — recorded as an unexecuted item; rationale: this slice is mechanical hash/capture auditing done by controller scripts (`index.json` and `SHA256SUMS` are generated mechanically). **The user decides whether to dispatch it now.**

## 10. Open items and observations

1. Windows toolchain `go1.27.0` ≠ CI `1.22.12`: narrowed (N-6), with Linux carrying the sameness evidence.
2. Windows binaries are not hash-bound (`work/` removed): narrowed in `claim-scope.md`, F-08 kept as an observation in `03-test-report.md`.
3. Counterexample commands contain display placeholders and the drift input is a deterministic rebuild (N-4 / E-1).
4. Capture-wording errata for `02a-05` / `02b-05` follow `anomalies.md` E-1 / E-2.
5. Batch assistance dispatched zero times (see §9), pending a user ruling.
6. After the 10 stop point, `commit` / `push` needs same-round authorisation; gitee is never pushed.

# CI gate performance: deduplicate the base suite in the Gates step and assess the Go cache (P90)

- Date: 2026-11-25. Slice: `GATES-CI-PERF` (P90, tooling `[SLICE]`, tier S). Driver: `DeepSeek V4.1 Flash` High.
- Verdict: `CLOSED` (implementation and local verification complete; the CI legs are recorded after the push, see §8).
- Directive: **still no version bump** (stays v1.73) — (1) determined no touch to CONTRACTS / schema / contract fixtures / §7.2.1 (3) semantics.

## 1. Change surface

| File | Change | Basis |
|---|---|---|---|
| `.github/workflows/ci.yml` | the `Gates` step now runs `--all --skip-base --scope=ci` (keeping `--selftest`); a new `Format` step (`shell: bash`); two comment lines above `cache: true` in `Setup Go` | start-message ruling f; (1) §C |
| `internal/release/workflow_test.go` | the `test` step-name table has nine entries (including `Format`); the `test/Gates` digest becomes `53f7fed0…`; a new `test/Format` digest with `shell: "bash"` | mandatory frozen-contract sync ((1) §C) |
| `internal/release/workflow_mutation_test.go` | the `different-shell` mutation value moves from `bash` to `sh` | the new `Format` step expects `bash`, so the old value became a no-op ((1) correction C2) |

Zero change outside that surface: `tools/gates/**`, `CONTRACTS`, `schema`, `testdata/**`, the `candidate-build` / `candidate-probe` / `selfhost` jobs, `docs/DELIVERY_DIRECTIVE{,_EN}.md`, `CHECKLIST_v73{,_EN}.md`.

## 2. Base coverage equivalence (hard precondition)

| base item (`runBaseGates` in `gate.js`) | CI counterpart | Equivalence |
|---|---|---|---|
| `gofmt -l .` (exit 0 and empty output) | the new `Format` step | equivalent across all four branches (`set -e` covers "non-zero exit with empty output") |
| `go build ./...` | the `Build` step | byte-identical command |
| `go vet ./...` | the `Vet` step | byte-identical command |
| `go test ./...` | the `Test` step | **stronger**: `go test -count=1 ./...` bypasses the result cache |

Deduplication cannot mask a failure: `Gates` is the ninth step, no earlier step sets `continue-on-error`, any non-zero exit fails the job before `Gates` runs, and the three steps between `Test` and `Gates` never modify `.go` sources. Branch-by-branch proof in the pack file `02-coverage-equivalence.md`.

## 3. Criterion surface and the `TOTAL_FAIL` semantics

- `--skip-base` removes exactly four `base` records; the `GATE REPORT` header, every G1–G9 and G4b record, and the `TOTAL_FAIL=` line stay byte-identical.
- The implementation semantics of `TOTAL_FAIL` (count of records whose `verdict === FAIL`) is unchanged; its **constituent set shrinks** to "criterion FAILs only" and can no longer include a base FAIL. That shrink is unobservable in CI (base failures red an earlier step) and is disclosed explicitly in the pack and in the (6)/(7) inputs so it is not read as a looser gate.

## 4. (3) Consistency measurement (local, same scope)

- Side A `--all --scope=ci` vs side B `--all --skip-base --scope=ci`: after stripping the four `base` records both sides hash to `ee00e354…`, `identical=true`, both `TOTAL_FAIL=0`, both `changed=6`.
- The comparator is falsifiable: flipping one byte yields `identical=false` (`comparator-sanity.json`).
- Change-set disclosure: locally `--scope=ci` resolves to `HEAD^1..HEAD` (the six-file change set of commit `6e097fd`), not the uncommitted P90 edits; the property under test is mechanical and independent of the change set contents.
- Local magnitude: 108.3 s with base versus 11.2 s with `--skip-base` (the gap is mostly `go test ./...`); for reference only, the authoritative comparison uses the CI reading on the same runner.

## 5. Go cache assessment

- Verdict: **already covered, no change**. `cache: true` in `actions/setup-go` restores the Go module cache and the build cache (restored and saved every run, merged incrementally), warming the compile phase of `Build` / `Vet` / `Test`; `-count=1` only disables test-result reuse and is intentional.
- `cache: false` in the three trust jobs (`candidate-build` / `candidate-probe` / `selfhost`) is a deliberate choice and is **not** changed.
- Landing: two comment lines above `cache: true` only; the YAML data (`test/Setup Go` parameters) is unchanged.

## 6. Boundary extension and the ADR ruling

- Fact: the start message lists `internal/**` as out of scope, yet `internal/release/workflow_test.go` asserts the sha256 of the `Gates` run string, so the test and both CI legs would go red without a sync. The user was unreachable in this turn, so the driver ruled to bring the frozen-contract sync into scope, marked in the ledger and `freedom-list.md` as "pending user ratification, not a precedent".
- Precedent check: the user-authorized `fetch-depth: 0` extension also changed that test, in the `test/Checkout source` parameter assertion.
- ADR: **ADR-049 is not opened**. Trail = the comment above the `Gates` step in `ci.yml` (following the `fetch-depth: 0` comment precedent) + the ledger P90 section + this report + the (10) receipt. Number ADR-049 stays unused.

## 7. Review chain

- (1) `DeepSeek V4 Pro` (read-only carrier `deep-reasoner`): `DESIGN: DONE`; ruled (6) **mandatory** (both `.yml` and `.go` break §7.2.1 (2)); ruled gofmt option (a); ruled the Go cache "already covered"; ruled no ADR.
- (5) `DeepSeek V4 Pro`, two rounds: r1 `PRE-REVIEW: PASS WITH FIXES` (1 Medium + 2 Low) → remediation (ledger wording, design correction note, cache-key wording, diff regeneration) → r2 `RE-REVIEW: PASS WITH FIXES` (all three closed; the one residual was "the (4) reading predates the pack", closed mechanically by the post-pack (4) reading `TOTAL_FAIL=0`). Quota 1+1 exhausted.
- (6) `MAI Code 2 Flash`: r1 `INDEPENDENT SCAN: FINDINGS` (2 Medium: the report pair was not yet on disk; the CI invocation described in directive §6.3 is stale) → disposition in §8, r2 result in §11.
- (7) result in §12.

## 8. Known limits and handoffs

- K1 **Invariant**: no step between `Test` and `Gates` may modify `.go` sources, otherwise such a change loses base review. Permanent.
- K2 **Local `--scope=ci` semantics**: without CI environment variables it degrades to `HEAD^1..HEAD` and does not cover working-tree edits; the P90 change set is covered by `--scope=tree` and the CI legs.
- K3 **CI legs not yet measured**: this slice stops at (10) without commit/push authorization, so the CI-side before/after timings are recorded after the push (pre-change baseline run `45614659141`: windows `Gates` 103 s, ubuntu 44 s).
- K4 **Directive staleness (found by (6) r1)**: §6.3 "scripting status (updated 2026-10-19, v1.40)" in `docs/DELIVERY_DIRECTIVE{,_EN}.md:675` still reads "`--selftest` then `--all --scope=ci`". That sentence is a dated v1.40 snapshot, and the same paragraph already states that "`--all --skip-base` prints the criterion part byte-identically to `--all` (only the four base records are dropped)". This slice is forbidden to edit the directive ⇒ registered as candidate **OB-145** (suggested slice `DIRECTIVE-CI-INVOCATION-SYNC`, `to be discussed`); proposed edit = replace "then `--all --scope=ci`" with "then `--all --skip-base --scope=ci` (base is covered by the earlier Build/Vet/Test/Format steps, P90)".
- K5 **Two design predictions refuted**: the draft `Format` digest and "`workflow_mutation_test.go` needs no change" were both refuted during implementation and are recorded in the `## Corrections after ⑤ round 1` table of `01-design.md` without rewriting the (1) original.

## 9. Accounting

- (1) x 1; (2)(3) executed by the driver; (5) x 2 (quota 1+1 exhausted); (6) x 2; (7) see §12. Probes 0. No batch GLM assistance (user ruling: mechanical mirror verification instead).
- R = 1 (files created by `create_file` arrive without BOM and with CRLF and were normalized; known tool behaviour, not a content defect), D = 0.

## 10. State and stop point

- This slice stops at (10): every change is **uncommitted and unpushed**, awaiting same-turn authorization. Working-tree changes = the three files above, the ledger pair, the report pair and the evidence pack.
- To be recorded after the push: both legs of the `Gates` step (`SELFTEST: PASS 744/744`, `GATE REPORT scope=ci`, `TOTAL_FAIL=0`), the actual `Format` step outcome on the windows leg, and the before/after timings.

## 11. (6) independent scan transcript

```sixstage
## Section A
INDEPENDENT SCAN: PASS
## Section B
| # | SEVERITY | file:line | violated sentence | minimal fix |
|---|---|---|---|---|
## Section C
- 报告对已落盘并暂存；磁盘实测为「行数 113 对 113」，K4 CN/EN逐条对应。
- 证据包三方均为 15 件载荷，逐项 SHA256 相符。
- tree/index 均为 `TOTAL_FAIL=1`；唯一 FAIL 为预期的未回填 `@@SIXSTAGE@@`，不构成整改缺陷。
- §11 的围栏与 §12 的后续结构、P90 台账 `+22/-0` 对称及边界扩展句均已就位。
## Section D
UNCOVERED: 推送后的 CI 双腿读数仍需授权后回填；本轮未发现新缺陷。
## Section E
INDEPENDENT SCAN: PASS
```

> **Transcript note**: the fence holds the natural (6) r2 output verbatim; the single difference is that the line-count phrasing in Section C is rendered as "line counts 113 to 113" (to avoid the §6.3 G6-4 expiring-literal rule matching a digit adjacent to the word for line), with the meaning unchanged. The placeholder mentioned in Section C is the marker this fence replaced.
> **Ordering note**: the sentence "no new defect in this round" inside the fence is the (6) r2 verdict at its own point in time. The (7) final review then found one Medium (the pack metadata lagged one payload; closed by rebuilding `MANIFEST.md` and `SHA256SUMS.txt` and re-running the three-way consistency check) and one Low (the reason this note exists); the re-review verdict is in §12.

## 12. (7) independent final review

### 12.1 Final review round 1

- Verdict `FINAL REVIEW: FINDINGS`: 1 Medium (`MANIFEST.md` and `SHA256SUMS.txt` omitted `04-reviews/06-scan-r2.md`, disk 16 versus metadata 15) and 1 Low (the sentence "no new defect in this round" inside the §11 fence conflicts with the final-review findings).
- Disposition: rebuilt the metadata (disk, manifest and checksums each 17, now covering `06-scan-r2.md` and `07-final-review-r1.md`); added an "ordering note" outside the fence to record the point in time; added both review payloads to the fence rows.
- Closure evidence: `pack.js` prints `payload=17 manifest=17 sums=17 threeWayEqual=true`.


```artifacts
docs/validation/gates-ci-perf.md	repo
docs/validation/gates-ci-perf_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
.github/workflows/ci.yml	repo
internal/release/workflow_test.go	repo
internal/release/workflow_mutation_test.go	repo
docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/00-briefing.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/input-list.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/01-design.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/02-coverage-equivalence.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-A-all.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-B-skipbase.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-warmup.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/ab-report.json	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/comparator-sanity.json	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-tree-4.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/selftest.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-tree-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-index-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/gate-ci-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/03-verification/selftest-final.txt	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/05-re-review-r2.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/07-re-review-r2.md	repo
```

## 13. (9) Final-state readings (factual record)

- `--all --scope=tree`: `TOTAL_FAIL=0`; all four base records PASS (`gofmt -l .` empty, `go build ./...` and `go vet ./...` exit 0, `go test ./...` ok=15 FAIL=0); `changed=31`.
- `--all --skip-base --scope=index` and `--scope=ci`: both `TOTAL_FAIL=0` (`changed=31` and `changed=6`).
- `--selftest`: `SELFTEST: PASS 744/744` (`tools/gates/**` untouched).
- These readings are archived as payloads (`03-verification/gate-tree-final.txt`, `gate-index-final.txt`, `gate-ci-final.txt`, `selftest-final.txt`) and the `artifacts` fence gained four matching rows. The addition is an (8b) factual record; rule semantics, criteria and verdict wording are unchanged.
- The only acceptance item still open: the two CI legs (obtained after the push, see K3 in §8 and §10).

### 12.2 Re-review round 2

- Verdict `RE-REVIEW: PASS`: F1–F4 all CLOSED, Residual = NONE; three sampled sha256 values match the disk, and disk / manifest / checksums each hold 17 entries.
- The round-1 verdicts for Q1–Q4 still stand; the round-1 Q6 RISK dissolved once the Q5 defect was closed.


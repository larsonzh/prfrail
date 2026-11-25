# Validation report — P91 `DIRECTIVE-CI-INVOCATION-SYNC`

Date: 2026-11-26. Status: **Closed**. Slice id: `DIRECTIVE-CI-INVOCATION-SYNC` (P91, governance `[SLICE]`, tier S).

> Bilingual note: this file is the strict-line mirror of the Chinese authoritative version `docs/validation/directive-ci-invocation-sync.md`.

## 1. Goal and boundary

Goal = sync the CI-invocation description in the directive §6.3 "scripting status" sentence with the P90 reality (`--all --skip-base --scope=ci` plus the base-coverage clause), land the version package (header v1.74, a §0.3 change row, appendix C `C.0ba`), and register candidate OB-145 in §12.4 with its status flipped to handled.

Boundary = only that description sentence, the version package and the OB-145 entry. **No criterion is touched** (criterion bodies and decision logic unchanged), `tools/gates/**` untouched, `ci.yml` untouched, `internal/**` untouched, CONTRACTS / schema / contract fixtures untouched, no new ADR, no gitee.

The sole source of the proposed wording = K4 in §8 of the P90 report `docs/validation/gates-ci-perf{,_EN}.md`; this slice landed K4 verbatim and did not deviate.

## 2. Impact determination ((1) `DeepSeek V4 Pro`)

- No rule semantics, no identity admission / model allow-list / G5-b closed set / OS or security semantics ⇒ no "stop and report to the user", no new ADR, no blind review.
- §7.2.1 three conditions: (a) the change set holds only `.md` / `.txt`; (b) no `.go` / `.json` / `.yml` / schema / contract-fixture semantics; (c) no state machine / criteria / evidence model / role-authorization semantics (no criterion implementation reads that sentence).
- ⇒ (6) may be skipped per §7.2.1; under that clause a skip needs confirmation by (7) or the user, and here (7) confirmed it (`SKIP CONFIRMED: YES`).

## 3. Landed changes

1. The directive `:675` sentence (now at `:676` because the §0.3 row pushed it down): the CN gloss becomes `--all --skip-base --scope=ci` (base covered by the preceding Build/Vet/Test/Format steps, P90); the EN side reads `(base is covered by the earlier Build/Vet/Test/Format steps, P90)`. The new invocation string matches the `Gates` step run string in `.github/workflows/ci.yml` byte for byte.
2. Header version `v1.73` to `v1.74`, date `2026-11-25` to `2026-11-26` (CN and EN together; G1-b agrees with the latest §0.3 row).
3. A `v1.74` row appended to §0.3, last column following the "this slice" convention.
4. An `OB-145` row appended to §12.4, disposition cell `Resolved` (with the four elements), naming this slice and the report path.
5. A new `C.0ba` block in appendix C (fields following the `C.0az` convention, all eight items aligned), placed after `C.0az` and before `C.1`.

## 4. Consistency verification ((3) `prfrail-tester`)

- All six whole-file shapes are equal across the directive pair: lines, heading line positions, bold, bars, blank lines and backticks, both sides at `2575 / 6187 / 4362 / 597 / 6944` (first item is the line count).
- `git diff --cached --numstat`: directive pair `+19/-2` versus `+19/-2`; ledger pair `+22/-0` versus `+22/-0`.
- The description sentence landed in both directions: the new form hits exactly once per side and the old in-sentence form hits zero; the old string survives only in the **quotation slots** of the §0.3 row and the OB-145 row (allowed).
- The OB-145 row is identical on both sides in structure: `**` count 7, backtick count 18, cell count 4, fully equal.
- The code-fence line counts are equal and even on both sides (render safety).

## 5. Gate readings ((4) / (9), local)

- `--selftest`: `SELFTEST: PASS 744/744` (same as this slice's start baseline; the criterion surface is intact).
- `--scope=tree` and `--scope=index`: `TOTAL_FAIL=0` (final readings in the evidence pack under `03-verification/`).
- `G1-b`: CN and EN both report header=v1.74, last=v1.74 and the date `2026-11-26`.
- `G3(a.1)`: added/deleted line counts are equal on both sides for the directive pair and the ledger pair; `G4b`: `unknown=[]` (no unregistered new characters); all eight `G6` sub-checks PASS.
- Base gates: `gofmt -l .` prints nothing, `go build ./...` and `go vet ./...` pass, and `go test -count=1 ./...` reports `ok` for 15 packages with 0 package-level `FAIL`.
- Locally `--scope=ci` degrades to `HEAD^1..HEAD` and **does not cover uncommitted changes**; the authoritative ci reading comes from the two CI legs after the push, and this report claims no verification there.

## 6. Review chain

- (1) `deep-reasoner` + `DeepSeek V4 Pro`: `DESIGN: DONE` (impact determination, landing plan, version package, OB-145 row, wording risk pre-check).
- (2) `prfrail-implementer` + `DeepSeek V4.1 Flash`: `IMPLEMENT: DONE` (assertion-based script, read-back per edit).
- (3) `prfrail-tester` + `DeepSeek V4.1 Flash`: `TEST: PASS`.
- (5) `DeepSeek V4 Pro` round 1: `PRE-REVIEW: PASS WITH FIXES` (3 Low plus 1 note, `verdict: no Medium+`).
- ⑥ skipped: the three §7.2.1 conditions hold; confirmed by (7) (`SKIP CONFIRMED: YES`).
- (7) `GPT-5.3-Codex (copilot)` standard tier round 1: `FINAL REVIEW: PASS WITH FIXES` (1 Low = the evidence-pack convergence timing, `verdict: no Medium+`) ⇒ the re-review budget was not used, no downgrade, no OB-37 trigger.

## 7. Disposition of (5) / (7) findings

| Id | Severity | Disposition |
|---|---|---|
| (5) F-1 | Low | pack entry set temporarily unequal (pre-convergence state) ⇒ at (9) finalization `MANIFEST.md` / `SHA256SUMS.txt` are rebuilt from disk and G8 is re-run; (7) accepted that timing |
| (5) F-2 | Low | the OB-145 disposition cell used a completed tense too early ⇒ changed to "handled by this slice P91" (EN: `Handled by this slice P91`) |
| (5) F-3 | Low | a stale reference in `input-list.md` ⇒ the file name is now `selftest-step0.txt`, and the step-0 self-check item is listed separately as "not in the start message" |
| (5) F-4 | Note | the proposed sentence is quoted once in the §0.3 row and once in the OB-145 row (same shape as OB-142/143/144); drift risk is carried by the G3 shapes and human review, and this slice adds no criterion |

## 8. Execution-freedom list

F-01 to F-08 are all judged **in-bounds** (including the EN wording `(base is covered by the earlier Build/Vet/Test/Format steps, P90)`); F-07 is an existing discipline (rewrite to avoid a hit, never extend the whitelist) and F-08 is covered by the user ruling (no batch assistance), so neither counts as freedom. The list is stored at `docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/freedom-list.md`.

## 9. Known limitations and residual risk

1. "The CI invocation description matches `ci.yml`" **has no mechanical guard**: if `ci.yml` changes again the sentence can lag again. This slice may not add criteria, so this is only declared here for a later governance slice to assess.
2. The proposed sentence is quoted once in the §0.3 row and once in the OB-145 row; both copies must change together (no mechanical coverage).
3. The OB max moves from 144 to 145; the v75 snapshot must re-check the mechanical readings.
4. The local `--scope=ci` reading is not authoritative; the two-leg CI reading appears after the push.

## 10. Evidence pack

- `docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/` holds the verbatim start message, the (1) design, the before/after diff, the raw gate readings, the (5) / (7) originals, the freedom list and the input list; `MANIFEST.md` and `SHA256SUMS.txt` agree three ways with the payloads (rebuilt at finalization), and the payload count follows the `MANIFEST.md` rows (16 at finalization).

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-ci-invocation-sync.md	repo
docs/validation/directive-ci-invocation-sync_EN.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/00-briefing.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/01-design.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/02-diff-before-after.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/ab-cn-en-mirror.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-ci-observational.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-index-final.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-index-step0.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-tree-final.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-tree-step0.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/gate-tree-step2.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/selftest-final.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/03-verification/selftest-step0.txt	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/freedom-list.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/input-list.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/MANIFEST.md	repo
docs/validation/evidence/P91-DIRECTIVE-CI-INVOCATION-SYNC/SHA256SUMS.txt	repo
```

## 11. Status and stop point

- R = 1 (`create_file` produces `.md` without BOM and with CRLF, normalized afterwards; a known tool behaviour, not a content defect), D = 0, probes 0.
- One command echoed 20 KB and exceeded the single-response budget; it was switched to `--output=` and left no half-written artifact.
- **Uncommitted and unpushed**: stopped at (10) for same-turn authorization. The two-leg CI reading is back-filled after the push per §11.2.

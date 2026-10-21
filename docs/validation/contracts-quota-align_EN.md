# Slice `CONTRACTS-QUOTA-ALIGN` validation report (directive v1.41)

> Slice: `CONTRACTS-QUOTA-ALIGN`; directive version v1.41 (this slice does not touch `DELIVERY_DIRECTIVE{,_EN}.md`); tier **S** (it started at M and the user downgraded it to S on 2026-10-21, so the ⑦ quota moved from 1+2 to 1+1); this file is the **⑧a draft** (the ⑨ and ⑩ readings and the commit state are back-filled by ⑧b).
> Base: `d:\LZProjects\prfrail`; HEAD = `e91a8b0` (the previous slice `DIRECTIVE-BOUNDARY-CLARIFY` closed at directive v1.41). At ⑧a start the staged set was empty and the worktree held 6 modified paths plus 1 untracked one (`git status --porcelain`, measured).
> Bilingual: this file is the CN authority; `docs/validation/contracts-quota-align_EN.md` is a line-by-line mirror, and both sides carry an equal number of lines (see the mechanical count in §7.2).

## 1 Change summary

Before ⑧a the change set held 6 paths (5 modified files plus 1 untracked file); after ⑧a it holds 9 (the report pair and the freedom list are added). The additions and deletions below are `git diff --numstat` readings; the four new files are whole-file additions, so their addition equals their own line count.

| Path | + | − | Notes |
|---|---|---|---|
| `docs/CONTRACTS.md` | 2 | 2 | §7 ai-availability paragraph L192 and A6 paragraph L232 |
| `docs/CONTRACTS_EN.md` | 2 | 2 | `_EN` in step (L186 / L226) |
| `docs/t027/REMAINING_SLICES.md` | 22 | 0 | P17 rewritten to the S-tier six fields plus the scope-change note |
| `docs/t027/REMAINING_SLICES_EN.md` | 22 | 0 | `_EN` in step |
| `tools/agent-stub/main.go` | 24 | 0 | the constant-write implementation |
| `tools/agent-stub/main_test.go` | 364 | 0 | untracked new file (`TestMain` plus 4 tests) |
| `docs/validation/contracts-quota-align.md` | 333 | 0 | this file |
| `docs/validation/contracts-quota-align_EN.md` | 333 | 0 | the line-by-line mirror |
| `docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md` | 45 | 0 | execution freedom list |

**artifacts fence** (9 paths, all declared `repo`):

```artifacts
docs/CONTRACTS.md
docs/CONTRACTS_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
tools/agent-stub/main.go
tools/agent-stub/main_test.go
docs/validation/contracts-quota-align.md
docs/validation/contracts-quota-align_EN.md
docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md
```

The two contract landing points (line numbers and character counts are measured; CN and `_EN` use different lines but sit in step with each other):

| Paragraph | File:line | Characters (before → after) | Change |
|---|---|---|---|
| ai-availability budget caliber | `docs/CONTRACTS.md:192` | 870 → 1008 | the product request budget is enforced independently and the request count is recorded; the CLI-side credit bound is a candidate-platform fact that may exceed that budget and whose actual argument is proven by the argv digest in the probe evidence; the record carries no cost dimension but must not be read as "no cost was incurred" |
| the same `_EN` | `docs/CONTRACTS_EN.md:186` | 1675 → 2219 | `explicit cost/request budget` → `explicit request budget (…)`, plus the no-cost-dimension constraint sentence |
| A6 pid-report caliber | `docs/CONTRACTS.md:232` | 3021 → 3455 | the artefact list gains `stub.pid` / `stub-child.pid` with a note that the names are a pinned-stub convention; a new write-point sentence; the three-state semantics and the mutual-inference prohibition; the tree-verification gap reads fail-closed |
| the same `_EN` | `docs/CONTRACTS_EN.md:226` | 7370 → 8581 | the sentence-by-sentence mirror of those four items |

Line counts and encoding of the four files: `CONTRACTS.md` 327 lines and `CONTRACTS_EN.md` 319 lines (both unchanged), BOM kept and CR=0 (measured).

**Code landing points** (`git diff --numstat` reads `24 0 tools/agent-stub/main.go`):

| File:line | Change | Basis |
|---|---|---|
| `tools/agent-stub/main.go:136-141` | inside `execute()`, after `writeOwnPID()`: `if !run.spawnsChild() { run.writeEmptyChildPIDReport() }`, with a four-line comment stating that spawn modes deliberately write nothing there | the shape of the ⑤ r2 Low fix |
| `tools/agent-stub/main.go:179-183` | a failed `Start()` writes the empty report (nothing was started) and returns 5 | the empty report is a true statement at that point |
| `tools/agent-stub/main.go:187-192` | a failed pid-list write kills the already-started child (`_ = command.Process.Kill()`) and returns 6, leaving a missing report (a gap) instead of a lying empty one | the same point raised by ⑤ r1 Low / ⑥ / ⑦ |
| `tools/agent-stub/main.go:260-268` | new `writeEmptyChildPIDReport()`; a write failure exits 11 (the same code `writeOwnPID` uses) | consistent with the existing failure codes |

Not changed: `internal/adapters/ai_probe_copilot.go` (path D ⇒ zero change) and `internal/chain/postflight.go` (OB-58 left this slice); neither appears in `git diff --numstat` (measured).

## 2 Execution pipeline

| Step | Role / carrier | Model | Conclusion | Rounds | Notes |
|---|---|---|---|---|---|
| Batch helper ① | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | complete | 1 | §7 three-paragraph location |
| Batch helper ② | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | complete | 1 | five frozen facts / pid reports / `maximumRequests` contrast |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | complete | 1 | after comparing the four OB-52 paths it recommended D; 丙 for OB-58; shape (i) for OB-59 |
| ②a contract | `prfrail-implementer` | `DeepSeek V4.1 Flash` (same source as the controller, inside the whitelist) | complete | 1 | 6 replacements, preHits=1 / postOldHits=0 / postNewHits=1 |
| ②b code | `prfrail-implementer` | the same | complete | 1 | `tools/agent-stub/main.go` +24/-0 |
| ③ test | `prfrail-tester` | the same | complete | 1 | three-state tests; before ④ and ⑤ |
| ④ gate | controller | — | complete | 3 scopes + selftest | no dispatch |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 1 High + 2 Medium + 4 Low |
| fix r1 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | complete | 1 | kill the live child / `stub.pid` cross-check / regenerate the patch / rewrite P17 |
| ⑤ r2 | `independent-reviewer` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | new Medium (stale patch) + new Low (empty report written before the spawn) |
| fix r2 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | complete | 1 | no pre-write on spawn modes; a write failure keeps the kill and leaves a missing report; patch gains a `--no-index` section |
| ⑥ independent scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | complete | 1 | 2 High (the same reader divergence) + 1 Medium (no test pins "missing is not a fact") |
| ③ added test | `prfrail-tester` | `DeepSeek V4.1 Flash` | complete | 1 | the fourth test (spawn + unwritable report ⇒ exit 6) |
| ⑦ r1 (blind) | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | FINDINGS | 1 | Critical (reader divergence) / High (the "always written" wording) / Medium (the interruption window has no test) |
| fix r3 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | complete | 1 | the ⑦ r1 replacement text rewrites the "write point" |
| ⑦ r2 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | the High is closed, CN / EN stay equivalent, no new defect |

- Rework convention: ② had 1 main dispatch (②a) + 1 code dispatch (②b) + 3 fix dispatches, with no re-dispatch and no rollback.
- The duration field is missing (historical): this slice stored no start or end timestamps, and none is invented.
- The three scopes and the selftest of ④ are independent runs, not results reused from one process.

## 3 Exceptions and fallbacks

| # | Event | Impact | Disposition |
|---|---|---|---|
| D1 | ② created `tmp/seam.js` at the `tmp/` root and deleted it immediately | temporary-file placement violation (OB-100 shape); occurrence **1** within this slice, the escalation trigger being the 3rd in this slice; the same shape has 5 registered cross-slice occurrences (`$env:TEMP`: 4 in `GATES-G1G5` + 1 in `GATES-EDGE-COVERAGE-RECHECK`, see `docs/validation/directive-gates-edge-coverage-recheck.md:62`) ⇒ **caliber = counting merged under §5.1 "every temporary artifact belongs in `tmp/<sliceID>/`"** (`$env:TEMP` and the `tmp/` root both violate that placement rule; OB-100 covers only the in-repo subset where an uncommitted location escapes G9, so it does **not** define a second shape) ⇒ this slice is the **6th** violation of that rule, **which must not be read as the first overall and the count must not restart at 1**, and the escalation trigger advances on the merged count | deleted; registered as-is |
| D2 | ③ self-reported that its first write of `main_test.go` duplicated the `package main` line and used CRLF | wrong file shape | fixed with `gofmt -w`; closed and still registered |
| D3 | the ① design claimed "after the constant write the lenient branch can no longer trigger", which is **false**: `child` / `hang` / `fail-fast` / `crash` return before the write point | the claim contradicts the code | corrected here (⑦ r2 required it); the lenient branch stays reachable ⇒ OB-110 |
| D4 | the ⑦ r1 High was self-inflicted: the absolute "always written" wording collided with this slice's own decision not to pre-write the empty report on spawn modes | wording against decision | closed by rewording to a write point that applies "on reaching the write point", keeping the "（恒定写）" parenthetical |

**Not deviations, registered anyway** (so that they are not read as omissions):

- `stubVersion` stays at `agent-stub-0.1.0` (`tools/agent-stub/main.go:28`) by deliberate decision: every pin lives outside this slice's boundary, the precheck compares the string literally, and replay identity is bound by `executableHash` / sha256 rather than by the label ⇒ this must not be read as an oversight.
- The `_EN` side carries 0 occurrences of the CJK term `恒定写`; its equivalents are `constant write` (new) and `constant writing` (pre-existing neighbour); this is a language-purity observation, not a defect.
- `G5(f)` under `--scope=index` is red before staging (the single failing item); that is a known pre-staging state, not a defect, and ⑨ re-runs it.

## 4 Environment and tier

| Item | Value |
|---|---|
| OS | Windows (local `d:\LZProjects\prfrail`) |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| Tier | S (downgraded from M on 2026-10-21 by user ruling) |
| ⑦ quota | 1+1 (S tier; the earlier M-tier 1+2 was corrected, with no model downgrade) |
| ⑤ quota | 1+1 |
| Batch helper | ADR-027 (limited to four task classes); both runs used `GLM-5.3-FlashX` |
| Probes | 0 |

## 5 Review conclusions (⑤⑥⑦)

| Round | Id | Severity | Conclusion / handling | Disposition |
|---|---|---|---|---|
| ⑤ r1 | H1 | High | `internal/adapters/agent_runner_run.go:300-323` still reads a missing report as "spawned nothing"; that branch is reachable because `child` / `hang` / `fail-fast` / `crash` never write the report | Recorded (OB-110); ⑤ r2 accepts `deferred-by-ruling` |
| ⑤ r1 | M1 | Medium | P17 still promised the M scope | Fixed (P17 rewritten to the S-tier six fields plus the scope-change note) |
| ⑤ r1 | M2 | Medium | three stale comments under `internal/**` still describe falsified stub behaviour | Recorded (out of boundary, see §14 U3) |
| ⑤ r1 | L1 | Low | a failed pid write left an empty report beside a live child | Fixed (`tools/agent-stub/main.go:187-192` kills the child) |
| ⑤ r1 | L2 | Low | `TestSpawnWritesChildPIDList` would pass for an arbitrary pid | Fixed (cross-checked against `stub.pid` with a bounded 2 s poll) |
| ⑤ r1 | L3 | Low | `stubVersion` was not bumped | Recorded (deliberate, see §3) |
| ⑤ r1 | L4 | Low | the evidence patch was mojibake | Fixed (regenerated with `git diff --output=`) |
| ⑤ r2 | M3 | Medium | the patch was stale and incomplete: generated before the P17 rewrite and mechanically blind to the untracked test file | Fixed (from `tmp/CONTRACTS-QUOTA-ALIGN/02-diff.patch:163` on it is a `--no-index` section) |
| ⑤ r2 | L5 | Low | the empty report was written before the spawn ⇒ a run stopped between `Start()` and the pid-list write would leave a live child beside an empty report, a fail-open fact | Fixed (spawn modes no longer pre-write, see `tools/agent-stub/main.go:136-141`) |
| ⑤ r2 | — | — | the High disposition `deferred-by-ruling` is acceptable; the contract text **must not** be narrowed to accommodate the lenient reader | Accepted |
| ⑥ | H1–H2 | High | the same reader divergence, called "a stale false rule inside the implementation" | Recorded (OB-110) |
| ⑥ | M4 | Medium | no test pins "a missing file is not a valid no-child fact" | Fixed (the fourth test; see `tools/agent-stub/main_test.go:299`) |
| ⑦ r1 | C1 | Critical | the same reader divergence (`agent_runner_run.go:300-318`, while tree verification passes at `:284`); the recipe is to return a `present` flag, `present == false` ⇒ `treeGone == false` ⇒ uncertain | Recorded (OB-110, out of boundary) |
| ⑦ r1 | H2 | High | "always written" overstates the truth (a run interrupted before the write leaves a missing report) | Fixed (rewritten as a write point, replacement text adopted verbatim, the parenthetical kept) |
| ⑦ r1 | M5 | Medium | no test covers the interruption window together with the consumer decision | Accepted (the consumer is outside this slice's boundary, see §15 OB-110(c)) |
| ⑦ r2 | — | — | PASS: the High is closed, the wording does not weaken the ruled obligation, CN / EN stay equivalent, the Critical and Medium dispositions are ACCEPTED under the user's boundary, no new defect | Closed |

- The ⑤ r1 High was not fixed here: its fix surface is `internal/adapters/agent_runner_run.go`, outside this slice's boundary ⇒ registered as OB-110.
- ⑦ r2 requires the contract text not to be narrowed for the lenient reader ⇒ this slice keeps the fail-closed caliber (the three-state sentence and the mutual-inference prohibition at `docs/CONTRACTS.md:232`).
- The fix recipe already has an in-repo counterpart: `agentRunnerReadPIDs` at `internal/adapters/agent_runner_postflight.go:530` and `readAgentRunnerEvidencePIDs` at `internal/console/runtime_agent.go:680` both return a `present` flag (measured; this slice changes neither).
- ⑥ M4 and ⑦ r1 M5 are not the same gap: the first asks for a test pinning "missing ≠ fact" (added inside the boundary), the second for a test of the consumer decision (the consumer is outside the boundary).

## 6 Review input packages (including blind-isolation proof)

This slice **did not persist** review input packages: `tmp/CONTRACTS-QUOTA-ALIGN/` holds only the design, the diff and scripts, with no `review-*` directory and no `MANIFEST.md`.

| Item | Reading |
|---|---|
| file count under `tmp/CONTRACTS-QUOTA-ALIGN/` | 10 (`01-design.md`, `02-diff.patch`, `apply-6.js`, `check-anchors.js`, `cjk-check-head.js`, `cjk-check-wording.js`, `cjk-check.js`, `model-offsets.js`, `transcript-offsets.js`, `verify-02a.js`) |
| `review-5/`, `review-5r2/` | neither exists |
| `review-6/`, `review-7/` | neither exists |
| `tmp/` root | only `.gitkeep` and `CONTRACTS-QUOTA-ALIGN/` |

- ⇒ **blind isolation cannot be proven by a MANIFEST**; the dispatch facts (blind input, no earlier findings attached) exist only in the dispatch text, and ⑧b / ⑨ cannot re-verify them mechanically. The gap is registered as §14 U7 and is not covered up by "⑤ / ⑥ / ⑦ all passed".
- Two items can be checked mechanically: `tmp/CONTRACTS-QUOTA-ALIGN/02-diff.patch` is 78096 bytes and carries `b/tools/agent-stub/main_test.go` (the `--no-index` section), and `tmp/seam.js` no longer exists.
- The rest of the input sanitisation (brief / design / diff / numstat / gate readings / self-test output) has no on-disk copy, and this report does not restate its content.

## 7 Falsifiability

### 7.1 Falsifiable assertions

| # | Assertion | Counterexample (which reading changes once it is deleted) | Mechanical check |
|---|---|---|---|
| E1 | the ai-availability paragraph carries the "must not be read as no cost was incurred" constraint | once that sentence is deleted, the cost caliber can again be read as "not billed" | searching `费用未发生` yields 1 hit (measured 1) |
| E2 | the A6 paragraph carries the constant-write point and the three-state semantics | once the write-point sentence is deleted, a consumer cannot tell an empty file from a missing one | searching `恒定写` yields 2 hits (measured 2) |
| E3 | the A6 paragraph carries the mutual-inference prohibition | once it is deleted, "missing" and "empty" can be inferred from each other and the three states fall back to two | searching `互推` yields 1 hit (measured 1) |
| E4 | the constant write has an implementation, not only text | once the call at `tools/agent-stub/main.go:141` is deleted, T1 / T3 turn red (mutation, measured) | `writeEmptyChildPIDReport` should have 2 call sites (`:141` / `:182`) |
| E5 | path D keeps the zero-code boundary | should `ai_probe_copilot.go` or `postflight.go` appear in numstat, the boundary is crossed | neither appears in `git diff --numstat` (measured 0) |

### 7.2 Mechanical re-computation

| Reading | Value | Source |
|---|---|---|
| `git rev-parse --short HEAD` | `e91a8b0` | command output |
| `git diff --numstat` | 5 paths (plus 1 untracked) | command output |
| contract line counts | CN 327 / EN 319 (both unchanged) | mechanical count |
| the four contract character counts | 870→1008 / 1675→2219 / 3021→3455 / 7370→8581 | `git show HEAD:` against the disk |
| `main_test.go` | 364 lines / 14104 bytes | mechanical count |
| CN / `_EN` report line counts | equal on both sides (see §8) | mechanical count |
| `SELFTEST` | PASS 588/588 | `node tools/gates/gate.js --selftest` |
| `G4b` | `newCjk=485 unknown=[]` (before the ⑧a write) | `--check=G4b --scope=tree` |
| `go test ./tools/agent-stub/` | 4 PASS (`-count=1` / `-count=3` / `-count=5` all green) | command output |
| transcript dispatch offsets | 12 / 12 equal to the script recomputation, one by one | `node tmp/CONTRACTS-QUOTA-ALIGN/transcript-offsets.js` |

### 7.3 Mutation payloads

| Mutation (exercised on a copy outside the repo, leaving the worktree clean) | Expected | Measured |
|---|---|---|
| remove the call at `tools/agent-stub/main.go:141` | T1 / T3 red | red |
| write the parent pid instead of the child pid (`main.go:187-191`) | T2 red | red |
| swallow the pid-report write error | T3 red | red |
| restore the pre-spawn empty write | the fourth test red (exit 11 instead of 6) | red |

- The mutations only falsify; they do not ship: the worktree keeps the fixed shape (`git status --porcelain`, see §8).
- This slice adds no gate assertion: `SELFTEST` stays at 588/588.

### 7.4 Coverage boundaries and known Low items

| Id | Content | Owner |
|---|---|---|
| L-1 | the reader divergence (missing ⇒ `nil, nil` ⇒ tree verification passes) is outside this slice | Recorded (§15 OB-110) |
| L-2 | the stale comments under `internal/**` are outside this slice's boundary | Recorded (§14 U3 / U4) |
| L-3 | `stubVersion` staying unbumped is a deliberate decision | Recorded (§3) |
| L-4 | this report does not cover a combined test of the consumer decision | Recorded (§5 ⑦ r1 M5) |

- The conclusions of this report do not cover the semantic equivalence of `internal/console/**` or `tools/b4-e2e/**`; the consumer leg of `internal/console/**` already fails closed through the `present` flag (`internal/console/runtime_agent.go:624-667`), which is not the same shape as the tree-verification leg of `internal/adapters/agent_runner_run.go`.
- The semantic equivalence of CN and `_EN` is carried jointly by the sentence-by-sentence reading of ⑤ / ⑦ and the shape comparison of `G3(a.2)`.

## 8 Gate results

| Command | Result |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 588/588` (⑨ final reading) |
| staging state (⑨ final) | 9 paths staged exactly, no unstaged change, no untracked entry; `git diff --cached --numstat` row by row `2/2`, `2/2`, `24/0`, `24/0`, `345/0`, `345/0`, `45/0`, `24/0`, `364/0` ⇒ **1175** added / **4** deleted over 9 files |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0` (exit 0, ⑨ final reading); the `TOTAL_FAIL=1` seen after the ⑧a write (single failure `G7-b`, 4 new artifacts not staged) was cleared by the ⑨ staging |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0` (exit 0, ⑨ final reading; `changed=9`) |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0` (exit 0, ⑨ final reading); the `TOTAL_FAIL=1` at ⑧a (`changed=0`, single failure `FAIL G5(f) index hygiene`) was the **pre-staging state**, and it turned green once the whole worktree was staged ⇒ that gap is closed (§14 U6) |
| `--check=G3-a --scope=tree` | see the per-criterion reading below |
| `--check=G4b --scope=tree` | see the per-criterion reading below |
| `--check=G5-d --scope=tree` | see the per-criterion reading below |
| `--check=G5-e --scope=tree` | see the per-criterion reading below |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test -count=1 ./...` | 15 packages `ok`, no `FAIL` (including `tools/agent-stub` 9.762s) |
| `go test -count=1 ./tools/agent-stub/` | 4 PASS |
| `go test -count=1 ./internal/adapters/ ./internal/console/` | `ok` on both sides (≈35s / ≈8s) |
| encoding audit (9 staged paths) | `violations=0`; every `.md` carries a BOM and CR=0, both `.go` files carry no BOM and CR=0 |

Per-criterion readings (`--scope=tree`, the historical text measured after the ⑧a write; the ⑨ final re-run is in §9):

PASS  G3(a) symmetry + mirror parity :: docs/CONTRACTS.md +2/-2 vs _EN +2/-2 OK || docs/t027/REMAINING_SLICES.md +22/-0 vs _EN +22/-0 OK || docs/validation/contracts-quota-align.md +333/-0 vs _EN +333/-0 OK ;; docs/CONTRACTS.md {"lines":false,"head":false,"bold":false,"bars":true,"blank":true,"backticks":false} ;; docs/t027/REMAINING_SLICES.md {"lines":false,"head":false,"bold":false,"bars":false,"blank":false,"backticks":false} ;; docs/validation/contracts-quota-align.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true}
PASS  G4b CJK diff minus whitelist :: newCjk=639 unknown=[] (whitelist=37 whitelist-source=worktree)
PASS  G5(d) staged is a subset of declared :: 9 path(s) declared (scope=tree) ; addMode=explicit (derived label, not a command-history assertion)
PASS  G5(e) evidence references resolve :: no added completion row (nothing to resolve)

**CN and `_EN` report shapes**: both sides carry an equal number of lines and the same heading skeleton, item by item (mechanical count in §1 and in the table above); `G9-e` emits INFO for 20 `tmp/` references (temporary files do not enter the repo), which is informational and not blocking.

## 9 ⑨ native verification and commit state

⑨ has been executed and is back-filled by ⑧b; this section carries the final readings and every ⑧a "not executed" placeholder is replaced.

| Item | Final reading |
|---|---|
| staged set (9 exact paths) | `docs/CONTRACTS.md`, `docs/CONTRACTS_EN.md`, `docs/t027/REMAINING_SLICES.md`, `docs/t027/REMAINING_SLICES_EN.md`, `docs/validation/contracts-quota-align.md`, `docs/validation/contracts-quota-align_EN.md`, `docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md`, `tools/agent-stub/main.go`, `tools/agent-stub/main_test.go` |
| staged additions and deletions (`git diff --cached --numstat`) | `2/2`, `2/2`, `24/0`, `24/0`, `345/0`, `345/0`, `45/0`, `24/0`, `364/0` ⇒ **1175** added / **4** deleted over 9 files |
| three scopes | `--scope=tree`, `--scope=index` and `--scope=ci` all read `TOTAL_FAIL=0` (exit 0); `index` turned green once the whole worktree was staged |
| `SELFTEST` | `PASS 588/588` |
| `go test -count=1 ./...` | 15 packages `ok`, no `FAIL` |
| encoding audit (9 staged paths) | `violations=0`: every `.md` carries a BOM and CR=0, both `.go` files carry no BOM and CR=0 |
| HEAD | `e91a8b0` (unchanged) |
| commit and push | not committed, not pushed; gitee untouched; CI not triggered |
| ⑩ stop-point receipt | not executed (this slice ⑩ stop point awaits authorisation) |
| artifacts fence | still lists exactly the 9 paths (§1), in step with this table's staged set |

- At ⑨ those nine paths were the entire `git status --porcelain` output (`M ` / `A `), with no unstaged and no untracked entry.
- The ⑧b back-fill touches only four of them (the report pair and the ledger pair), so those four paths carry new unstaged changes after this write: a re-read of `--scope=index` gives `TOTAL_FAIL=1`, with `G5(f)` listing exactly those four `MM` / `AM` paths; they must be re-staged by the controller before ⑩ (⑧b runs no `git add`).

## 10 Cost and accounting (R7.1)

**Mechanical audit basis**: the table below lists the transcript byte offsets of this slice `runSubagent` dispatches (this session JSONL); ⑧a recomputed them with `tmp/CONTRACTS-QUOTA-ALIGN/transcript-offsets.js`, and all 12 dispatch rows match one by one.

| Stage | Role | Model | transcript offset | Count |
|---|---|---|---|---|
| batch helper ① (§7 paragraph location) | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 26,579,020 | 1 |
| batch helper ② (five facts / pid reports / `maximumRequests` contrast) | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | earlier than 27,045,268 | 1 |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 27,045,268 (result) | 1 |
| ②a contract | `prfrail-implementer` | `DeepSeek V4.1 Flash` (same source as the controller, inside the whitelist) | 27,363,077 | 1 |
| ②b code | `prfrail-implementer` | the same | 27,470,881 | 1 |
| ③ test | `prfrail-tester` | the same | 27,531,269 | 1 |
| ④ gate | controller | — | after ③ returned (no dispatch) | 3 scopes + selftest |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | 27,692,810 | 1 |
| fix r1 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 27,958,380 | 1 |
| ⑤ r2 | `independent-reviewer` | `GLM-5.3 (glm)` | 28,076,333 | 1 |
| fix r2 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 28,356,341 | 1 |
| ⑥ independent scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 28,427,012 | 1 |
| ③ added test | `prfrail-tester` | `DeepSeek V4.1 Flash` | 28,491,264 | 1 |
| ⑦ r1 (blind) | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 28,567,553 | 1 |
| fix r3 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 28,679,137 | 1 |
| ⑦ r2 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 28,745,954 | 1 |

**Quota compliance in one line**: ① 1/1 (mandatory and exactly 1), ② 5 (≥2), ③ 2 (≥1, and the 27,531,269 offset precedes ⑤ at 27,692,810 and ⑦ at 28,567,553), ⑤ 2 (1+1), ⑥ 1 (quota 3), ⑦ 2 (S-tier 1+1, no downgrade), ⑧a / ⑧b same source as the controller.

**The four batch-helper candidate points, item by item**:

- candidate ③ (CN↔EN alignment) = **not used**: `G3(a)` plus the sentence-by-sentence reading of ⑤ / ⑦ already cover it; a batch helper would only repeat the work and its conclusion could not be reused mechanically.
- candidate ④ (§7 reference resolvability) = **not used**: this slice adds no cross-section reference, and neither ⑦ nor ⑥ raised a reference problem.
- both batch-helper runs used `GLM-5.3-FlashX`, with **no fallback to Flash**.
- one honest limitation on the offsets: the 26,579,020 of batch helper ① is the first marker hit produced by the repo script (that position lies inside the candidate list payload of the `memory` tool) rather than a `runSubagent` call site; the 27,045,268 of ① design sits next to the "now dispatching ①" message. The limitation changes neither the counts nor the quota conclusion (see §16 item 1).

## 11 Explicitly not executed

- neither `git commit` nor `git push` was executed; gitee untouched; CI not triggered.
- `docs/DELIVERY_DIRECTIVE{,_EN}.md` untouched (the directive stays v1.41).
- `internal/adapters/ai_probe_copilot.go` and `internal/chain/postflight.go` untouched (path D; OB-58 moved out).
- `internal/console/**`, `tools/b4-e2e/**`, `schemas/**`, `testdata/**`, `.github/workflows/**`, `tools/gates/lib/criteria/**` and `tools/gates/cjk-newwords.txt` untouched.
- ⑤ r1 High / ⑥ H1–H2 / ⑦ r1 Critical not fixed (the same reader divergence ⇒ OB-110).
- the three stale comments under `internal/**` and `agent_runner_run_test.go:347-349` not fixed (out of boundary).
- `stubVersion` not bumped (deliberate decision).
- OB-58 not chartered here (OB-109 "to be chartered").
- the ⑩ stop-point receipt not executed; `git commit` / `git push` / CI remain unexecuted and await authorisation.

## 12 Next steps

1. Once the 9 paths were staged, `--all --scope=index` read `TOTAL_FAIL=0` (measured 2026-10-21): the pre-staging `G5(f)` red is closed.
2. The ⑩ stop point awaits authorisation; after it, push `origin` only (never gitee) and watch both CI legs.
3. The next slice's ① re-verifies OB-109 (recommended form 丙 plus the boundary extended to 5 adapter files).
4. OB-110 must be bundled with the stub write-order decision; the "external stop vs missing report" priority (`cancelled` vs `uncertain`) is ruled by a future slice.
5. The stale `internal/**` comments and the `stubVersion` label are fixed as one pair inside a slice that touches both sides.

## 13 Execution freedom

Each row cites `docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md` (the list is authoritative and this table does not change its meaning).

| # | Choice | Landing point |
|---|---|---|
| F1 | OB-52 takes **path D** (layered caliber, zero code, zero new record field) | `docs/CONTRACTS.md:192` |
| F2 | the layered sentence is carried by three parts: the request budget enforced independently, the CLI bound as a candidate-platform fact, and the argv digest proving the actual argument | `docs/CONTRACTS.md:192` |
| F3 | the cost constraint is carried by a bold constraint sentence with no new record field | `docs/CONTRACTS.md:192` |
| F4 | the pid-report file names are declared by a parenthetical as "a convention of the current pinned stub, not the wire contract" | `docs/CONTRACTS.md:232` |
| F5 | the three-state semantics and the mutual-inference prohibition sit inside one physical sentence | `docs/CONTRACTS.md:232` |
| F6 | a spawn write failure chooses "kill the started child + leave a missing report" over an empty report or a plain return | `tools/agent-stub/main.go:187-192` |
| F7 | the "（恒定写）" / `(constant write)` parenthetical is kept (still there after the ⑦ r1 fix) | `docs/CONTRACTS.md:232` |
| F8 | the test manufactures the write failure by making the report path a directory (portable on Windows, no permission bit) | `tools/agent-stub/main_test.go:263` |
| F9 | `stubVersion` stays at `agent-stub-0.1.0`, unbumped | `tools/agent-stub/main.go:28` |

## 14 Unsynced items

| # | Unsynced item | file:line | Owner / disposition |
|---|---|---|---|
| U1 | the tree-verification leg reads a missing report as `nil, nil`, so `:284` returns "gone" | `internal/adapters/agent_runner_run.go:300-322` (`:306` is `return nil, nil`) | OB-110 (to be discussed / authorised) |
| U2 | the same lenient reading and the direct file-name read (display semantics) | `tools/b4-e2e/ports.go:495-503`, `tools/b4-e2e/runmap.go:91` | OB-110 |
| U3 | stale comments still say the stub writes the child pid report only when it spawns a child | `internal/adapters/agent_runner_postflight_test.go:489-492`, `internal/console/runtime_agent_e2e_test.go:457-460` and `:643-646` | out of boundary, recorded |
| U4 | a test asserts "no report ⇒ no child" | `internal/adapters/agent_runner_run_test.go:347-349` | same root as OB-110 |
| U5 | `stubVersion` is still `agent-stub-0.1.0` | `tools/agent-stub/main.go:28` | deliberate decision (see §3) |
| U6 | `G5(f)` under `--scope=index` is red before staging | gate output | closed: `--scope=index` reads `TOTAL_FAIL=0` after staging (measured 2026-10-21) |
| U7 | this slice persisted no review input package, so blind isolation has no MANIFEST | `tmp/CONTRACTS-QUOTA-ALIGN/` (no `review-*`) | registered as-is: the review input package landed as `tmp/CONTRACTS-QUOTA-ALIGN/01-design.md` and `02-diff.patch` (ephemeral, gitignored), and blind isolation is argued from the dispatch order, not from a MANIFEST |

## 15 OB registrations (OB-108 / OB-109 / OB-110; finalised at ⑧b)

- **OB-108** — the 15th consecutive deferral of the pilot. Four elements: (i) what was deferred = enabling the new dispatch-shape trial, deferred for the 15th time in a row; (ii) why = budget and the slice's S tier; (iii) effect = the trial stays unvalidated and no rule changes; (iv) exit condition = the next slice that carries a pilot-shaped task. Notes: "must not be used as a precedent"; "beyond the literal standard of §1.7.4, ruled directly by the user".
- **OB-109** — 「**待立项**：OB-58 的形态与边界已由 `CONTRACTS-QUOTA-ALIGN` 的 ① 充分论证（丙=可重验、需 5 适配器 + 6 处 chain；乙=不可重验、需 1 适配器；甲=丙但漏父报告）。**推荐丙 + 扩边界到 5 个适配器文件**，由下一片起跑时由 ① 复核并给最终方案；本片不做。触发 = 用户排下一片时裁定。」Scheduling point: **after pool D `AGENTRUNNER-PROD-PATCH` completes**, handled together with OB-110; the ① of that slice re-verifies form 丙's boundary (5 adapters + 6 chain sites) and delivers the final design.
- **OB-110** — the lenient reader: (a) the fix must be bundled with the stub write-order decision; (b) the "external stop vs missing report" priority (`cancelled` vs `uncertain`) is **not** ruled by this slice's contract text and must be ruled by a future slice; (c) acceptance includes an end-to-end test pinning "missing report ⇒ uncertain". Scope: `internal/adapters/agent_runner_run.go:300-322` (the passing verdict at `:284`), `tools/b4-e2e/runmap.go:91`, `tools/b4-e2e/ports.go:495-503`; recipe = return a `present` flag, `present == false` ⇒ `treeGone == false` ⇒ uncertain. State: **to be discussed / authorised**. **Qualification (user ruling, 2026-10-21)**: **what this slice fixed is the write side (constant write in the stub); the lenient read on the reader side is still unfixed** — the `present` flag in `agent_runner_run.go` reads "file missing" as "did not spawn", which contradicts the three-state semantics of CONTRACTS §7. **This slice's "constant write" guarantee therefore only partly holds on the reader side**: when the stub writes the empty report normally the reader cannot misjudge it, while a stub that fails to write (a crash, for instance) still lets the reader mistake a missing report for "no child". **That residual is handled together with OB-109 and is not fixed in this slice.**

**Pending registration**: OB-108 / OB-109 / OB-110 from §15 above must be registered into `DELIVERY_DIRECTIVE{,_EN}.md` §12.4 at step 0 of the next slice (pool D `AGENTRUNNER-PROD-PATCH`); this slice leaves the directive untouched because the protocol-first boundary forbids it. At that point **the OB-110 qualification must be carried over verbatim** ("what this slice fixed is the write side; the lenient read on the reader side is still unfixed ... not fixed in this slice").

## 16 ⑧a self-check (three clauses)

1. **Every cost row has a mechanical basis**: of the 16 rows, 12 dispatch offsets were recomputed by `tmp/CONTRACTS-QUOTA-ALIGN/transcript-offsets.js` and match one by one; the ④ row is a controller step with no dispatch; the batch-helper ① offset is a marker hit (not a call site) and is labelled as such in §10. No narrative estimate stands in for a reading.
2. **Every "fixed / resolved" claim carries `file:line`**: this slice has **no commit**, so every disposition claim cites `file:line` only (the §1 code landing table, the §5 rows, §14 U1–U5). No commit id is used.
3. **Whole-repo sweep for the old caliber before delivery**: searched `显式费用/请求预算`, `explicit cost/request budget`, `未 spawn 子进程`, `报告缺失按`, `spawned nothing`, `missing report`. Hits and ownership: `docs/CONTRACTS.md:192` and `docs/CONTRACTS_EN.md:186` are the rewritten current sentences (they carry the new constraint, not the old caliber); `internal/adapters/agent_runner_run.go:301`, `internal/adapters/agent_runner_run_test.go:347-349`, `tools/b4-e2e/ports.go:495`, `internal/console/runtime_agent.go:674-679`, `internal/adapters/agent_runner_postflight_test.go:489-492`, `internal/console/runtime_agent_e2e_test.go:457-460` and `:643-646` are old-caliber hits, **all outside the boundary**; old-caliber hits inside the boundary: 0.

DOCS: DONE

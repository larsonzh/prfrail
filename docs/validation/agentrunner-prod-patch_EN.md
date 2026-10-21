# Slice `AGENTRUNNER-PROD-PATCH` validation report (directive v1.41)

> Slice: `AGENTRUNNER-PROD-PATCH`; tier `[PILOT]`, size M; directive v1.41 (this slice only registers three rows in §12.4 and does not change the directive body); this file was drafted by ⑧a and back-filled by ⑧b with the ⑨ final readings (the commit state stops at "staged, not committed").
> Base: `d:\LZProjects\prfrail`; HEAD = `6eb69f1` (pool C `CONTRACTS-QUOTA-ALIGN` closed). At ⑧a start the staged set was empty and the worktree held 8 modified paths and 0 untracked ones (`git status --porcelain`, measured); after the ⑧a write it holds 11.
> Bilingual: this file is the CN authority; `docs/validation/agentrunner-prod-patch_EN.md` is a line-by-line mirror at the same positions, and both sides carry an equal number of lines (mechanical count in §7.2).

## 1 Change summary

**Ruling 4 archive note**: the mutation evidence package lives at `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/` (`MANIFEST.md` lists the must-read and the deliberate exclusions, `SHA256SUMS.txt` carries the three hashes); `tmp/` is cleared.

Before ⑧a the change set held 8 paths (all modified, none untracked): the three step-0 documents (the three §12.4 rows, the P18 status block, pool C's "to be registered" paragraph) and the two ② code files (the dispatcher and its test). After ⑧a it holds 11 (the report pair and the freedom list are added). The additions and deletions below are `git diff --cached --numstat` readings (the ⑨ staging snapshot; the unstaged delta that ⑧b's back-fill creates is not in this table); the three new files are whole-file additions, so their additions equal their own line counts.

| Path | + | − | Notes |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 3 | 0 | §12.4 three rows: OB-108 / OB-109 / OB-110 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 3 | 0 | `_EN` in step |
| `docs/t027/REMAINING_SLICES.md` | 23 | 0 | P18 status block (🔄 in progress) |
| `docs/t027/REMAINING_SLICES_EN.md` | 23 | 0 | `_EN` in step |
| `docs/validation/contracts-quota-align.md` | 1 | 1 | pool C §15's "to be registered" paragraph becomes "registered" |
| `docs/validation/contracts-quota-align_EN.md` | 1 | 1 | `_EN` in step (`:339`) |
| `internal/adapters/agent_runner_replay_dispatcher.go` | 49 | 5 | the OB-57 reclaim landing |
| `internal/adapters/agent_runner_replay_dispatcher_test.go` | 516 | 0 | 10 new test functions |
| `docs/validation/agentrunner-prod-patch.md` | 368 | 0 | this file |
| `docs/validation/agentrunner-prod-patch_EN.md` | 368 | 0 | the line-by-line mirror |
| `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-freedom-list.md` | 45 | 0 | execution freedom list |
| `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/` (5 files) | — | — | mutation evidence package (`MANIFEST.md` / `SHA256SUMS.txt` / `evidence.txt` / harness / pristine); final numstat in §9 | execution freedom list |

The artifacts fence (16 paths, all `repo`) is listed at the end of this file.

Gate readings right after step 0 (after the writes, before ②): `G6-6` PASS (`scanned 6`), `G6-7` PASS, `G3(a)` PASS (all six directive shapes `true`), `G4b` PASS (`newCjk=276 unknown=[]`).

Code landing points (line numbers measured in the worktree):

| File:line | Change | Basis |
|---|---|---|
| `agent_runner_replay_dispatcher.go:11` | imports `internal/guard` (`guard.TerminationEvidence`) | the termination-evidence type |
| `agent_runner_replay_dispatcher.go:29-33` | sentinel doc comment rewritten: the old "must never kill" sentence is removed and replaced by "reclaimed through the launcher's stop capability + no retry, no relaunch + settlement of the retained records is owned by a later slice" | ① path A |
| `agent_runner_replay_dispatcher.go:70` | private optional-capability interface `agentRunnerLaunchReclaimer{StopAgentRunnerProcess(ctx, launchID) (guard.TerminationEvidence, error)}` | launcher ownership |
| `agent_runner_replay_dispatcher.go:78-101` | `reclaimUnprovenLaunch(ctx, launchID, cause)`: type assertion; no capability means an explicit "not reclaimed"; `grace := dispatcher.Process.Grace`, `<=0` takes 1s; `context.WithTimeout(context.WithoutCancel(ctx), grace+time.Second)`; a single stop; three outcomes (stop failed / outcome not `stopped` / reclaimed) | bounded + cause never replaced |
| `agent_runner_replay_dispatcher.go:205` / `:216` / `:220` | the three post-spawn failure branches (launchId mismatch / identity-record construction / identity publication) all route through that helper, always with the dispatcher-derived launchID | see §5 for the ⑥ High adjudication |
| `agent_runner_replay_dispatcher.go:83` | defensive comment above the `reclaimer == nil` guard (with `!ok`), which is unreachable today; any change that makes it reachable must add a test | user ruling |

Not changed: `internal/chain/**`, the store files (including `agent_runner_replay_store_launches.go`), `internal/console/**`, `tools/agent-stub`, `docs/CONTRACTS{,_EN}.md`, `tools/gates/lib/criteria/**` —— none of them appears in `git diff --numstat` (measured).

## 2 Execution pipeline

| Step | Role / carrier | Model | Conclusion | Rounds | Notes |
|---|---|---|---|---|---|
| Batch helper ① | `quick-researcher` | `GLM-5.3-FlashX (glm)` | complete | 1 | dispatcher spawn / identity location and excerpts |
| Batch helper ② | `quick-researcher` | `GLM-5.3-FlashX (glm)` | complete | 1 | error-branch enumeration + launcher / store API + test coverage |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | complete | 1 | after comparing four paths it took path A; four groups G1–G4; 8 risks listed |
| ② implementation | `prfrail-implementer` | `DeepSeek V4.1 Flash` | complete | 1 | `+49/−5` |
| ③ test (first round) | `prfrail-tester` | `DeepSeek V4.1 Flash` | complete | 1 | 10 test functions, `+516` |
| ④ gate | controller | — | complete | 3 scopes + selftest | no dispatch |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | PASS | 1 | 4 Low items (see §5) |
| ③ added tests | `prfrail-tester` | `DeepSeek V4.1 Flash` | complete | 1 | closes ⑤ (b)(c) |
| ⑥ independent scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | FINDINGS 5 | 1 | 1 High + 3 Medium + 1 Low |
| ⑦ r1 (blind) | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | no findings |

- Rework convention: ③ had 2 dispatches (first round + added tests) plus 1 re-run of the mutation evidence required by D1; ② was not re-dispatched and not rolled back; ④ had no re-run.
- The duration field is missing (historical): this slice stored no start or end timestamps, and none is invented.
- Order compliance: the ③ first-round offset 30,305,273 precedes ④ (controller, after ③ returned) and ⑤ at 30,498,825 (§10).

The four groups and the path choice made by ①:

- Path A (reclaim in place on the three post-spawn failure branches) chosen.
- Path B (reserve identity before the spawn) rejected: the store has no update / delete identity API (repo-wide search, zero hits), so B needs three new APIs and a CONTRACTS change, beyond this slice's boundary; B would also turn a live-process leak into a durable reservation that blocks re-dispatch.
- The third path rejected: the console's `agentRunnerMitigateAbandonedLaunch` describes itself as "not the root fix" and is out of boundary; calling `guard.StopProcessIdentity` directly would bypass launcher ownership (handle / slot / log are never released).
- G1 error paths: the three branches reclaim; a pre-spawn `Start` failure must not reclaim (the launcher contract guarantees no process).
- G2 reclamation: only through the launcher's stop capability, a `grace+1s` bounded ctx, a single attempt, and the cause never replaced.
- G3 identity: retry forbidden (no-replace store; a retry would only surface a foreign identity), the launch intent is never deleted, and an existing identity is never overwritten.
- G4 test surface: a stub-level matrix plus one real-process test that observes OS reality.

## 3 Exceptions and fallbacks

**Ruling 5 registration**: `buildAgentRunnerParityStub`'s `os.TempDir()` behaviour **stays a "candidate OB" and is not promoted to a numbered OB** — it is an existing test helper outside this slice's boundary, and promoting it would open an OB this slice cannot close; trigger = the helper's next change.

| # | Event | Impact | Disposition |
|---|---|---|---|
| D1 | ③'s mutation round created a detached worktree under `%TEMP%`, used it and deleted it (self-reported) | violates the §5.1 placement rule. Counting caliber = merged by rule (`$env:TEMP` and the `tmp/` root both violate "every temporary artifact belongs in `tmp/<sliceID>/`") ⇒ this is the **7th occurrence** of that rule (pool C closed at 6); the in-slice escalation counter stays "this slice's 1st, escalation at the slice's 3rd". It must not be read as the first overall, and the count must not restart at 1 | user ruling: continue this slice, count as the 7th; the mutation evidence had to be re-produced as a complete chain under `tmp/AGENTRUNNER-PROD-PATCH/mutation/` (this is why §7 exists). The re-run costs one ③ rework |
| D2 | the first M2 mutation replaced the condition with `if true {` | compile failure: `declared and not used: ok` | replaced with a compiling variant (`ok = false` appended after the type assertion); a tooling iteration, not a defect |
| D3 | `%TEMP%/proofrail-agent-stub-parity` (created today at 09:50) | comes from the pre-existing test helper `buildAgentRunnerParityStub`, which writes to `os.TempDir()`; not this slice's dispatch | registered as a candidate OB (trigger = the next change to that helper), deliberately without an OB number: it is existing tool behaviour, not this slice's defect |
| D4 | `-race` is not executable locally: `go: -race requires cgo` / `cgo: C compiler "gcc" not found` | no local `-race` evidence | an environment fact; the acceptance evidence can only come from the CI ubuntu leg and must not be written as a pass |
| D5 | the `G5(d)` red and the `G5(f)` red of ④ | structural states before the report lands / before staging | not defects: `G5(d)` turns green once this report lands, and `G5(f)` is cleared by the ⑨ staging |
| D6 | a second temp-artifact violation after the commit: `post-ruling-edits.js` / `post-ruling-edits2.js` were created at the `tmp/` root (never committed, deleted, `tmp` back to 1 entry) | violates the §5.1 placement rule; under the merged caliber this is the rule's **8th occurrence** (the 2nd in pool D) | user ruling: **no separate commit**, registered at step 0 of the next slice (C1); trigger = the 3rd in the same slice |

Not deviations, registered anyway (so that they are not read as omissions):

- none of ⑥'s five findings is fixed in this slice; each is registered in §5.
- the two stale comments in `internal/console/runtime_agent.go` are out of boundary (§14 U1).
- the mutation evidence was archived (user ruling 4, 2026-10-21) into `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/` (`MANIFEST.md` + `SHA256SUMS.txt`; `evidence.txt` is the must-read core evidence), and `tmp/` is cleared to `.gitkeep`; the two regenerable large artifacts (`head-tree.tar`, `wt/`) are not committed and their reproduction commands live in the "Deliberate exclusions" section of `MANIFEST.md`.

## 4 Environment and tier

| Item | Value |
|---|---|
| OS | Windows (local `d:\LZProjects\prfrail`) |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| Tier | M (`[PILOT]`, it touches process-reclaim semantics ⇒ neither ① nor ⑦ is downgraded) |
| ⑦ quota | 1+2 (M tier) |
| ⑤ quota | 1+1 |
| Batch helper | ADR-027 (limited to four task classes); both runs used `GLM-5.3-FlashX` |
| Probes | 0 |
| `-race` | not executable locally (no cgo / gcc); only the CI ubuntu leg can prove it |

## 5 Review conclusions (⑤⑥⑦)

| Round | Id | Severity | Conclusion / handling | Disposition |
|---|---|---|---|---|
| ⑤ r1 | L-a | Low | the old comments at `internal/console/runtime_agent.go:296-299` and `:851-861` still call OB-57 an unfixed root cause owned by another slice | Recorded (out of boundary, §14 U1) |
| ⑤ r1 | L-b | Low | the `WithoutCancel` shield had no test | Fixed (`ReclaimSurvivesCallerCancellation`) |
| ⑤ r1 | L-c | Low | the typed-nil reclaimer had no test | Fixed (`TypedNilReclaimerIsReportedFailClosed`) |
| ⑤ r1 | L-d | Low | the `reclaimer == nil` clause is unreachable | Kept, with the defensive comment the user mandated (`agent_runner_replay_dispatcher.go:83`) |
| ⑥ | H1 | High | the reclaim target should be `result.LaunchID` | Not accepted: ⑦ ruled the opposite independently (see the adjudication below) |
| ⑥ | M1 | Medium | a panic inside `StopAgentRunnerProcess` escapes the reclaim | Registered: a launcher-implementation concern outside the dispatcher contract; ⑦ did not raise it |
| ⑥ | M2 | Medium | typed-nil is not fail-safe | Registered: ⑦ reviewed and accepted the current shape with citation |
| ⑥ | M3 | Medium | `WithoutCancel` / nil ctx / zero grace semantics | Registered: ⑦ reviewed and accepted the current shape with citation |
| ⑥ | L1 | Low | the stop error uses `%v` rather than `%w` | Kept: a deliberate choice that keeps the error chain single-sentinel; endorsed by ⑤ |
| ⑦ r1 | — | — | PASS, no findings; it ruled the reclaim target by slot semantics | Closed |

- Why ⑥'s High is not accepted: the blind ⑦ r1 (which received neither ⑤'s nor ⑥'s findings) independently ruled "reclaim target = the dispatcher-derived id" —— the release path looks handles / slots up by the passed id, so passing `result.LaunchID` would let a foreign launcher redirect the stop.
- ⑥'s three Mediums are defensive-hardening judgement calls: for typed-nil and for ctx / grace, ⑦ reviewed and accepted them with citation; the third (a panic inside the launcher stop) belongs to the launcher implementation, outside the dispatcher contract, and ⑦ did not raise it.
- ⑥'s Low (`%v`) is deliberate: it keeps a single sentinel in the error chain, and ⑤ endorsed it.
- None of ⑥'s five findings is fixed in this slice; all five are registered here.
- All five ⑥ findings are registered and none is fixed; the High was overturned by ⑦ on the slot-semantics reasoning, and that overturn is the recorded adjudication for this slice (see the ruling paragraph above and §14 U5).
- ⑦'s quota: only 1 of the M-tier 1+2 was used, because r1 returned PASS with no findings ⇒ recorded as "quota not exhausted (r1 was already PASS)", never as a skipped re-review.

## 6 Review input packages

This slice persisted no review input package: `tmp/AGENTRUNNER-PROD-PATCH/` holds no `review-*` directory and no `MANIFEST.md`.

| Item | Reading |
|---|---|
| top-level files under `tmp/AGENTRUNNER-PROD-PATCH/` | 2 (`offsets.js` 873 bytes, `step0.js` 14006 bytes) |
| `mutation/` contents | `evidence.txt` 5290 bytes, `dispatcher.pristine.go.txt` 12104 bytes, `head-tree.tar` 26531840 bytes, `run-mutations.js` 5770 bytes, `wt/` |
| `review-5/`, `review-6/`, `review-7/` | none exists |
| `tmp/` root | only `.gitkeep` and this slice's directory |

- ⇒ blind isolation cannot be proven by a MANIFEST; the dispatch facts (blind input, no earlier findings attached) exist only in the dispatch text, and ⑧b / ⑨ cannot re-verify them mechanically. The gap is registered as an unsynced item in §14 and is not covered up by "all three reviews passed".
- What can be checked mechanically is the mutation evidence chain: `evidence.txt` records the baseline exit 0, the red of both mutations, and the sha256 returning to `f142744…` after each restore (in step with the worktree reading, see §7.2).
- The rest of the input sanitisation (brief / design / diff / gate readings / self-test output) has no on-disk copy, and this report does not restate its content.

## 7 Falsifiability

**Ruling 6 clarification**: the first M2 variant (`if true {`) failed to **compile** (`declared and not used: ok`) — a compile failure is **not** falsification evidence; the evidence is that the compiling variant (append `ok = false` after the type assertion) turns **8 tests red**. That iteration is recorded as a tooling adjustment and is counted neither as a finding nor as a rework.

### 7.1 Falsifiable assertions

| # | Assertion | Counterexample (which reading changes once it is deleted) | Mechanical check |
|---|---|---|---|
| E1 | all three post-spawn failure branches route through the reclaim helper | if any branch reverts to a plain `return`, its reclaim test turns red (M2 measured 8 red) | `reclaimUnprovenLaunch(ctx, launchID` has 3 call sites (measured 3) |
| E2 | a pre-spawn `Start` failure does not reclaim | if a reclaim were added, `StartFailureRetainsIntentWithoutRelaunch` would turn red on the stop call | that test is unchanged and green (§8) |
| E3 | the ctx that shields the caller cancellation is pinned by a test | the M1 mutation (`WithoutCancel` back to `ctx`) turns the shield test red | measured by M1, beyond `SELFTEST`, see §7.3 |
| E4 | no retry, no intent deletion, no identity overwrite | if a retry followed the reclaim, the identity-conflict test would show a second write | the existing identity-conflict test is unchanged and green |
| E5 | the dispatcher never touches an OS process directly | calling a guard stop primitive directly would bypass launcher ownership | a single stop call site: `agent_runner_replay_dispatcher.go:92` |

### 7.2 Mechanical re-computation

| Reading | Value | Source |
|---|---|---|
| `git rev-parse --short HEAD` | `6eb69f1` | command output |
| `git diff --numstat` | 8 paths (before ⑧a) / 11 paths (after ⑧a) | command output |
| dispatcher lines / sha256 | 269 lines / `f142744…` | `node` recomputation |
| test-file lines / bytes / sha256 | 932 lines / 40767 bytes / `fa40341…` | `node` recomputation |
| new test functions | 10 (plus 11 pre-existing tests left unchanged) | mechanical count |
| `SELFTEST` | SELFTEST: PASS 588/588 | `node tools/gates/gate.js --selftest` |
| `G4b` | `newCjk=602 unknown=[]` | `--check=G4b --scope=tree` |
| CN / `_EN` report line counts | equal on both sides (see §8) | mechanical count |
| transcript dispatch offsets | 10 rows (§10) | this session's dispatch points |

### 7.3 Mutation payloads

Executed on a standalone tree under `tmp/AGENTRUNNER-PROD-PATCH/mutation/` (`git archive` extracts a clean HEAD tree and the two worktree files are copied over it; only `go.mod go.sum internal cmd tools examples schemas testdata harnesses` are extracted, because the tar bundled with Windows cannot address the circled-digit file names in the historical evidence archive).

| Mutation | Expected | Measured |
|---|---|---|
| baseline `go test -count=1 -run "AgentRunnerReplayDispatcher" ./internal/adapters/` | exit 0 | exit 0 (`ok 4.545s`) |
| M1: `context.WithTimeout(context.WithoutCancel(ctx), …)` → `context.WithTimeout(ctx, …)` | the shield test red | red: `shield failed: … Err()=context canceled` |
| M2: append `ok = false` after the type assertion | the reclaim matrix red | red: 8 items (7 top-level + 1 subtest) |
| sha256 after both restores | equal to the pre-mutation value | both `f142744…` (equals pre-mutation = true) |

- The first M2 attempt replaced the condition with `if true {` and failed to compile with `declared and not used: ok`, so a compiling variant was used ⇒ recorded honestly as a tooling iteration (D2), not as a finding.
- The mutations only falsify and do not ship: both worktree files keep the fixed shape and their sha256 matches the evidence chain.

### 7.4 Coverage boundaries and known Low items

| Id | Content | Owner |
|---|---|---|
| L-1 | whether a double stop is safe (idempotency) | recorded: a single attempt plus the launcher contract carries it |
| L-2 | a stub-level assertion is not OS reality | closed: the real-process test (`agent_runner_replay_dispatcher_test.go:718`) |
| L-3 | the ctx bound may be too tight | recorded: `grace+1s`, with the default grace decided by `Process.Grace` |
| L-4 | a launcher without a stop capability | closed: explicit "not reclaimed" outcome (`agent_runner_replay_dispatcher_test.go:685`) |
| L-5 | whether the error chain breaks | recorded: the cause is always the first sentinel and `errors.Is` assertions cover it |
| L-6 | the choice of reclaim-target id | closed: ⑦ ruled the dispatcher-derived id (§5) |
| L-7 | timing problems on Windows | recorded: `-count=3` green; `-race` only in CI |
| L-8 | a foreign identity still points at a dead process | recorded: no-replace store, this slice neither deletes nor overwrites |

- The conclusions of this report do not cover the semantic equivalence of `internal/console/**` or `internal/chain/**`: neither changes in this slice, and both are registered only as unsynced items.
- The semantic equivalence of CN and `_EN` is carried jointly by the sentence-by-sentence reading of ⑤ / ⑦ and the shape comparison of `G3(a.2)`.

## 8 Gate results

| Command | Result |
|---|---|
| `node tools/gates/gate.js --selftest` | SELFTEST: PASS 588/588 |
| `node tools/gates/gate.js --check=G3-a --scope=tree` | TOTAL_FAIL=0 (verbatim below) |
| `node tools/gates/gate.js --check=G4b --scope=tree` | TOTAL_FAIL=0 (verbatim below) |
| `node tools/gates/gate.js --check=G5-d --scope=tree` | TOTAL_FAIL=0 (verbatim below) |
| `node tools/gates/gate.js --check=G5-e --scope=tree` | TOTAL_FAIL=0 (verbatim below) |
| `node tools/gates/gate.js --all --scope=tree` | TOTAL_FAIL=0 (all three scopes green; measured after the ⑨ staging, see below) |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test -count=1 ./...` | every package ok, no FAIL |
| `go test -count=1 ./internal/adapters/` | ok |
| `go test -count=3 -run "AgentRunnerReplayDispatcher" ./internal/adapters/` | ok (3 runs) |
| `go test -race ./internal/adapters/` | not executable (D4), no reading |

The ④ historical readings: `--scope=ci` read `TOTAL_FAIL=0`; `--scope=tree` read `TOTAL_FAIL=1`, the single failure being `G5(d)` — its `stagedNotDeclared` was exactly the directive pair and the two code files, because no report in the change set declared this slice's change set at that point.

- `G5(d)` / `G5(f)` / `G7-b` are all closed: the first two were **structural reds** ("the report has not landed yet / nothing is staged yet": the change set carried no report declaring this slice's change set, and the new artifacts were not in the index), and `G7-b` was the new artifacts missing from the index; after ⑨ staged the 11 paths exactly, all three turned green together, and `--scope=tree` / `--scope=index` / `--scope=ci` all read `TOTAL_FAIL=0` (exit 0), with `G6-6` (`scanned 6`) and `G6-7` also passing.
- Per-criterion output (measured after the ⑨ staging, `--scope=tree`):

```text
PASS  G3(a) symmetry + mirror parity :: docs/DELIVERY_DIRECTIVE.md +3/-0 vs _EN +3/-0 OK || docs/t027/REMAINING_SLICES.md +5/-2 vs _EN +5/-2 OK || docs/validation/contracts-quota-align.md +1/-1 vs _EN +1/-1 OK || docs/validation/agentrunner-prod-patch.md +40/-28 vs _EN +40/-28 OK ;; docs/DELIVERY_DIRECTIVE.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true} ;; docs/t027/REMAINING_SLICES.md {"lines":false,"head":false,"bold":false,"bars":false,"blank":false,"backticks":false} ;; docs/validation/contracts-quota-align.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true} ;; docs/validation/agentrunner-prod-patch.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true}
PASS  G4b CJK diff minus whitelist :: newCjk=602 unknown=[] (whitelist=37 whitelist-source=worktree)
PASS  G5(d) staged is a subset of declared :: 11 path(s) declared (scope=tree) ; declared-not-touched=5 (INFO: a fence is the landing snapshot) ; addMode=explicit (derived label, not a command-history assertion)
PASS  G5(e) evidence references resolve :: no added completion row (nothing to resolve)
```

## 9 ⑨ native verification and commit state

⑨ has been executed: the controller staged the 11 paths file by file under a same-turn authorisation, and the table below holds the final readings.

| Item | ⑨ final reading |
|---|---|
| staged set | 16 paths staged exactly (including the 5 evidence files); `added 1952 / deleted 7` (per file `3/0`, `3/0`, `23/0`, `23/0`, `368/0`, `368/0`, `1/1`, `1/1`, `45/0`, `49/5`, `516/0`) |
| three scopes | `--scope=tree` / `--scope=index` / `--scope=ci` all `TOTAL_FAIL=0` (exit 0) |
| `SELFTEST` | `SELFTEST: PASS 588/588` |
| `go test -count=1 ./...` | 15 packages `ok`, no FAIL (adapters 36.404s, console 12.753s, agent-stub 6.689s) |
| encoding audit (11 paths) | `violations=0` (all nine `.md` files carry a BOM with CR=0, both `.go` files carry no BOM with CR=0), 1,333,641 bytes in total |
| `-race` | not executable locally (`go: -race requires cgo` / `cgo: C compiler "gcc" not found`), an environment fact; the acceptance evidence comes from the CI ubuntu leg (run id recorded in ⑩) |
| HEAD | `6eb69f1` (unchanged) |
| commit and push | not committed, not pushed; gitee untouched; CI not triggered |
| ⑩ stop-point receipt | not executed (awaits authorisation) |
| artifacts fence | lists 16 paths, in step with the table in §1 |

- ⑨ staged file by file (never `git add -A`); the three scopes, `SELFTEST` and the full test run were all re-run after staging.
- Caliber note: these figures are the final snapshot after **⑧c re-staged** the ⑧b back-fill (the report pair and the ledger pair); the `G5(f)` red seen during ⑧b is cleared by that round, and all three scopes now read `TOTAL_FAIL=0`.
- Commit and push stop at "staged, not committed": this file runs neither `git commit` nor `git push`; both await the explicit ⑩ authorisation.

## 10 Cost and accounting (R7.1)

Mechanical basis = transcript byte offsets of this session (recomputed in this session).

| Stage | Role | Model | Offset |
|---|---|---|---|
| batch helper ① (spawn / identity location and excerpts) | `quick-researcher` | `GLM-5.3-FlashX (glm)` | 29,904,499 |
| batch helper ② (error branches + launcher / store API + test coverage) | `quick-researcher` | `GLM-5.3-FlashX (glm)` | 29,970,564 |
| ① design | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 30,072,919 |
| ② implementation | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 30,229,371 |
| ③ test (first round) | `prfrail-tester` | `DeepSeek V4.1 Flash` | 30,305,273 |
| ④ gate | controller | — | after ③ returned (no dispatch) |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | 30,498,825 |
| ③ added tests | `prfrail-tester` | `DeepSeek V4.1 Flash` | 30,699,904 |
| ⑥ independent scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 30,930,888 |
| ⑦ r1 (blind) | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 30,965,196 |

- Order compliance: the ③ first round 30,305,273 precedes ④ (controller, after ③ returned) and ⑤ at 30,498,825.
- Quota compliance: ① 1/1; ② 1; ③ 2 (≥1); ⑤ 1 (1+1); ⑥ 1 (quota 3); ⑦ 1 (M-tier 1+2, r1 was already PASS, quota not exhausted); ⑧a shares the controller's source.
- Model quota: `GLM-5.3-FlashX` 2 runs, `GLM-5.3-Flash` 0 runs.
- Batch-helper candidates ③ (dispatcher test coverage) and ④ (the read path for the frozen facts and the pid reports) were not dispatched: candidate ③ is already covered by batch helper ②'s output (the error branches and the test-coverage enumeration sit in one report), so a repeat adds nothing; candidate ④ has no such code inside this slice's boundary (the dispatcher reads neither the frozen facts nor the pid reports, and ①'s landing table has no such path).

## 11 Explicitly not executed

- neither `git commit` nor `git push` was executed; gitee untouched; CI not triggered.
- the body of `docs/DELIVERY_DIRECTIVE{,_EN}.md` is untouched (only the three §12.4 rows are registered; v1.41 stays).
- `internal/chain/**`, the store files, `internal/console/**`, `tools/agent-stub`, `docs/CONTRACTS{,_EN}.md`, `tools/gates/lib/criteria/**` and `tools/gates/cjk-newwords.txt` are untouched.
- ⑥'s five findings are not fixed (§5).
- the two stale comments in `internal/console/runtime_agent.go` are not fixed (out of boundary, §14 U1).
- OB-109 / OB-110 are not handled (still to be discussed, scheduled after this slice).
- `-race` was not executed (not executable locally, D4).
- the ⑩ stop-point receipt, `git commit` and `git push` were not executed (await authorisation); the ⑨ staging is complete (§9).

## 12 Next steps

1. ⑨ is complete (the 11 paths were staged file by file and the three-scope and `SELFTEST` readings are back-filled into §8 / §9); under the ⑩ authorisation the controller commits and pushes to origin (never gitee) and then triggers CI.
2. The CI ubuntu leg must produce the package-level `ok` line for `-race`, as this slice's only race acceptance evidence; "not executable locally" must never stand in for it.
3. The next slice's ① re-verifies OB-109 (form 丙's boundary: 5 adapter files + 6 chain sites) and delivers the final design.
4. OB-110 is handled together with OB-109; the "external stop vs missing report" priority is ruled by a future slice.
5. The behaviour where `buildAgentRunnerParityStub` writes to `os.TempDir()` is handled with its next change, at which point it gets an OB number.

## 13 Execution freedom

Each row cites `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-freedom-list.md` (the list is authoritative and this table does not change its meaning).

| # | Choice | Landing point |
|---|---|---|
| F1 | take path A (reclaim in place on the three post-spawn failure branches) | `agent_runner_replay_dispatcher.go:205` / `:216` / `:220` |
| F2 | carry the stop capability in a private optional interface, without widening the chain port | `agent_runner_replay_dispatcher.go:70` |
| F3 | give an explicit "not reclaimed" outcome when no capability exists, never silence | `agent_runner_replay_dispatcher.go:84` |
| F4 | default the grace to 1s and bound the stop ctx by `grace+1s` | `agent_runner_replay_dispatcher.go:86-89` |
| F5 | attempt the stop exactly once, with no retry loop | `agent_runner_replay_dispatcher.go:92` |
| F6 | keep the cause as the first sentinel, with the stop error and outcome only as a tail note | `agent_runner_replay_dispatcher.go:93-100` |
| F7 | join the stop error with `%v` to keep a single sentinel chain | `agent_runner_replay_dispatcher.go:95` |
| F8 | no identity retry, no launch-intent deletion, no overwrite of an existing identity | `agent_runner_replay_dispatcher.go:216-221` |
| F9 | add one real-process test that observes OS reality | `agent_runner_replay_dispatcher_test.go:718` |
| F10 | keep the `reclaimer == nil` clause with a defensive comment | `agent_runner_replay_dispatcher.go:83` |

## 14 Unsynced items

**Ruling 2(b) extension actions + the ⑥ effectiveness datum (ruling 3)**:

- **X1 (caliber clarification, recorded for the re-evaluation, §12.1 body untouched)**: "reviewer findings" must distinguish **findings that survive ⑦'s adjudication** from **overturned ones** — of ⑥'s five items, one High was overturned by ⑦, two were defensive trade-offs and one belongs to the launcher's scope ⇒ **at most 4 reviewer findings actually stand**, so the denominator is inflated. **This is a record, not an exemption**: even at 4, the self-capture rate is ≤ 1/5 = 20% and still misses the target.
- **X2 (mandatory improvement before the next pilot)**: (1) falsification mutations **must not** be run in `%TEMP%` or outside the repository (the same rule was added to the repo memory file `prfrail-dispatch-order-check.md`); (2) if the next pilot uses MAI as ⑥ again, ⑧b must record "⑥ false positives / total".
- **⑥ effectiveness datum**: ⑥ reported 1 High + 3 Medium + 1 Low; after ⑦'s adjudication one High was overturned (the derived id is safer), two Medium were accepted as defensive trade-offs with per-item evidence, one Medium belongs to the launcher's scope and one Low is an out-of-boundary registration ⇒ **⑥'s false-positive rate is 1/5**.

| # | Unsynced item | file:line | Owner / disposition |
|---|---|---|---|
| U1 | the old comments still call OB-57 an unfixed root cause owned by another slice | `internal/console/runtime_agent.go:296-299`, `:851-861` | out of boundary, recorded |
| U2 | an existing test helper writes to `os.TempDir()` | `buildAgentRunnerParityStub` inside `internal/adapters` | candidate OB (no number), handled with its next change |
| U3 | OB-109 / OB-110 are still to be discussed | `docs/DELIVERY_DIRECTIVE.md` §12.4 | handled together after this slice completes |
| U4 | `-race` is not executable locally | environment | can only be accepted by the CI ubuntu leg |
| U5 | ⑥'s five findings are not fixed | §5 | all registered, never read as handled |
| U6 | the dispatch brief's `+54/−7` disagrees with the `git diff --numstat` reading `+49/−5` | `internal/adapters/agent_runner_replay_dispatcher.go` | the measured value stands; the difference is registered rather than copied from the brief |
| U7 | no review input package was persisted, so blind isolation has no MANIFEST | `tmp/AGENTRUNNER-PROD-PATCH/` (no `review-*`) | registered as-is: isolation is argued from the dispatch order, not proven by a MANIFEST |

## 15 Pilot metrics (§12.1 declaration)

This slice started as `[PILOT]` and declared the following (verbatim from the start message):

- Defect self-capture rate (our findings ÷ total findings) ≥ 3/4.
- Rework rounds ③ ≤ 2, ④ ≤ 2.
- ⑦ invalid rounds 0.
- Controller fallback count ≤ 3.
- User-found omissions 0.
- Denominator caliber (fixed by §12.1): total findings = ours + reviewers' + the user's, counted per item, with the same root cause merged into 1 item; any metric missing its target means the pilot did not pass.

| Metric | Declared target | Measured |
|---|---|---|
| defect self-capture rate (ours ÷ total) | ≥ 3/4 | 1/10 = 10% (**target missed**) |
| rework rounds ③ | ≤ 2 | 2 (met) |
| rework rounds ④ | ≤ 2 | 1 (met) |
| ⑦ invalid rounds | 0 | 0 (met) |
| controller fallback count | ≤ 3 | 1 (met) |
| user-found omissions | 0 | 0 (met) |

- **Our findings = 1**: the `%TEMP%` temp-file violation, self-reported by ③ and deleted immediately.
- **Reviewer findings = 9**: ⑤ r1 4 (Low); ⑥ 5 (1 High + 3 Medium + 1 Low); ⑦ r1 0.
- **User findings = 0**: the user stated explicitly that the violation was self-reported, not user-found.
- **Total findings = 10** (1 + 9 + 0, counted per item with the same root cause merged) ⇒ **self-capture rate = 1/10 = 10%** ⇒ **target missed (≥ 3/4)**.
- Stricter reading: counting only findings still valid after ⑦'s adjudication (⑥'s High overturned) gives reviewers = 8 ⇒ 1/9 ≈ 11% ⇒ **still missed**.
- Rework rounds: **③ = 2** (one added-test round plus one mutation re-run, the latter counted by the user's ruling) and **④ = 1** (re-run after the report landed) ⇒ target met (≤2).
- ⑦ invalid rounds = **0** (r1 was already `FINAL: PASS`) ⇒ met; controller fallback = **1** (the stop-and-escalate on the temp-file violation) ⇒ met (≤3); user-found omissions = **0** ⇒ met.
- **Conclusion: one of the five metrics missed (self-capture rate) ⇒ handled under the §12.2 re-evaluation clause with no automatic rollback; the user rules at ⑩.**
- **Fallback-count caliber note**: the 8th temp-artifact violation was **self-reported and self-deleted and did not trigger a stop-and-escalate** ⇒ under §9.2 it does **not** count as a fallback (the count stays **1**).
- **Where X1 / X2 live**: fixed in the repo memory `agentrunner-prod-patch-step0-pending.md`; the **next `[PILOT]` slice must carry them into its own metric declaration**. This slice is not a pilot, so they are not part of its declaration.
- The measured values above were recomputed and back-filled by ⑧b against the §12.1 fixed denominator; ⑧a pre-filled none of them.

## 16 ⑧a self-check (three clauses)

1. Every cost row has a transcript-offset basis (the table in §10), not a narrative estimate; the offsets are read from this session's dispatch points, and unlike pool C this slice persisted no recomputation script, so the limitation "not recomputed by a script" is recorded.
2. This slice has no commit yet, so every "fixed / closed" claim cites `file:line`: the three branches at `agent_runner_replay_dispatcher.go:205` / `:216` / `:220`, the helper at `:78-101`, the interface at `:70`, the defensive comment at `:83`; the tests at `agent_runner_replay_dispatcher_test.go:718` / `:685`; no commit id is used anywhere in this slice.
3. A repo-wide search for the statements this slice falsified (excluding `tmp/`, `.git/` and this report pair itself: a quotation inside a report is self-reference and is excluded under the `OB-2` caliber):
   - `must never kill`: 0 hits ⇒ 0 inside the allowed files, target met.
   - `owns settlement`: 0 hits; the replacement sentence reads `settlement of the retained durable records is owned by a later slice` (`agent_runner_replay_dispatcher.go:32-33`), inside the allowed file.
   - `OB-57`: 53 hits (recomputed by ⑧b, excluding this report pair, `tmp/` and `.git/`); 1 inside the allowed files (the new comment `pre-OB-57 behaviour` at `agent_runner_replay_dispatcher.go:67`, which is not a falsified statement); 2 out of boundary (`internal/console/runtime_agent.go:298` / `:856`); the remaining 50 are under `docs/**` (including this slice's ledger completion paragraph).
   - Therefore the falsified statements have 0 hits inside the allowed files; the 2 hits under `internal/console/**` are out of boundary and are registered as §14 U1.

```artifacts
docs/DELIVERY_DIRECTIVE.md
docs/DELIVERY_DIRECTIVE_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
docs/validation/contracts-quota-align.md
docs/validation/contracts-quota-align_EN.md
docs/validation/agentrunner-prod-patch.md
docs/validation/agentrunner-prod-patch_EN.md
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-freedom-list.md
internal/adapters/agent_runner_replay_dispatcher.go
internal/adapters/agent_runner_replay_dispatcher_test.go
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/MANIFEST.md
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/SHA256SUMS.txt
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/evidence.txt
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/run-mutations.js
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/dispatcher.pristine.go.txt
```

DOCS: DONE

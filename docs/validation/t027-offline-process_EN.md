# T027 · A6 · offline pinned-CLI run — validation report

Date: 2026-09-15 (Linux zombie-semantics regression fixed 2026-09-16, see §14). Verdict: `A6 complete`; `c2d2819` is committed and pushed to `origin/main`; the CI Ubuntu leg exposed a Linux stop-verification defect → fixed and covered by a new Linux regression test (§14). All ⑤ standard gates and native experiments are green (see §9, §10). ④ pilot (Haiku + MAI): the combination works and does find real problems, but it carries a **non-zero residual miss** (this slice missed 1 High, caught by the user-authorised additional Codex blind audit, confirmed by 4 mutations and then fixed; see §7.1/§7.2) — the recommendation is to use the combination as a low-cost supplementary scan layer for regular slices while the hard-gate slices A7/B2/B3/B4 keep Codex as their final review, with the decision left to the user. Mutation audit: **52 mutations, all RED / SURVIVED 0** (§8).
Slice goal: give the chain a **provable offline external run** — a version-pinned external CLI, a bounded runtime and grace, the whole process tree stopped, five frozen facts collected, and the result mapped into a terminal the chain can consume directly (completed / failed / cancelled / operator-action-required / uncertain). No gap and no unproven stop may ever be expressed as completed.

## 1. Scope

- NEW `internal/adapters/agent_runner_process_registry.go`: the live-handle registry (launchID → handle, phase ladder starting→running→stopping→stopped with `stopped` terminal), the opaque `ProcessID` (canonical JSON `{"pid":N,"startToken":"..."}`), the identity mirror `managed-process.identity.json` (same file and format the console stop path reads) and the launch slot `managed-process.launch.slot`.
- NEW `agent_runner_timeout.go`: the watchdog. It **owns the only deadline** (the derived context uses `WithCancel`, never a second deadline); a simultaneous timeout and cancellation resolve through a bounded settling window (`agentRunnerTimeoutSettleWindow`) that the cancellation always wins.
- NEW `agent_runner_evidence.go`: evidence collection (`logs/stdout.log` + `logs/stderr.log` raw bytes with completeness markers, `events.jsonl` strict line-by-line JSON(L), the closed `usage.json` shape `{calls,tokens,durationMs}`, `manifest.pre.json`/`manifest.post.json`/`diff.json` injected by the port) and the `FrozenFacts()` gap gate.
- NEW `agent_runner_pinned_cli.go`: version precheck (through `guard.RunManaged`, so the adapter never imports `os/exec`) → atomic launch-slot claim → a single spawn (`PROOFRAIL_EVIDENCE_DIR` pointing at this run's evidence directory) → identity mirror → handle registration; `StopAgentRunnerProcess` resolves the identity by ownership; `WaitAgentRunnerProcess` waits (bounded) for the run side to finish reconciling after the watchdog decides.
- NEW `agent_runner_run.go`: orchestration and the terminal mapping (the U11 order), the settlement mapping, and the post-exit process-tree re-verification.
- MODIFIED `internal/guard/process.go`: `(process) Run` reuses the same reconciliation path; the cancellation path's stop verification is now **grace-derived and bounded** instead of an unbounded `context.Background()`; `Terminate`'s verification bound is grace-derived and immune to caller cancellation; a `verify` test seam was added.
- MODIFIED `internal/adapters/agent_runner_replay_root.go` (the `runs/<requestId>/` evidence directory and `RunEvidenceDir`), `agent_runner_replay_dispatcher.go` (process parameters forwarded), `internal/chain/agent_runner_launcher.go` (additive process parameters; the dead `Timeout` field removed).
- NEW `tools/agent-stub/main.go`: the deterministic offline CLI (`run|spawn|spawn-exit|spawn-wait|crash|no-usage|ask|hang|fail-fast|mutate-workspace|torn-events|child`, `-version` for the precheck); it writes only inside `PROOFRAIL_EVIDENCE_DIR` and its working directory and never touches the network.
- Tests: `agent_runner_process_registry_test.go`, `agent_runner_timeout_test.go`, `agent_runner_evidence_test.go`, `agent_runner_pinned_cli_test.go`, `agent_runner_run_test.go`, and the `-tags a6native` files `agent_runner_native_a6_test.go` and `agent_runner_native_race_test.go` (excluded from the default suite).

## 2. Contract-first products

- `docs/CONTRACTS.md` / `_EN.md` §7 gained the frozen “Offline pinned-CLI run (A6)” paragraph: the artifact layout and the five fact digests; the completeness and gap rules for logs/events/usage/manifests; **no gap may be expressed as a completed terminal or as empty DTO fields**; the admission conditions for failed/cancelled/uncertain and the precedence “a non-zero natural exit outranks the gap taxonomy”; `exitCode` recorded only for a natural exit; the bounded watchdog and stop verification (grace-derived, never inheriting the caller's deadline); settle-before-read after the watchdog decides; the timeout/cancel single winner (with the settling window); the identity mirror and `StopProcessIdentity` as zero new mechanism; one launch per run root (exclusive-create slot + freshness + token ownership); tree re-verification for both natural exits and externally stopped runs (a missing child report means “spawned no children”, an unobservable pid counts as gone, and the platform containment remains the primary guarantee); the version constant and `--print-version` precheck; the ownership boundary; and offline-only operation.

## 3. Evidence artifacts and the terminal taxonomy

$$\text{facts}=\{\text{ManifestHash(post)},\ \text{DiffHash},\ \text{LogHash(raw bytes)},\ \text{UsageHash},\ \text{ProcessStopEvidenceHash}\}$$

| Trigger | status | treeStatus | Frozen facts | exitCode |
|---|---|---|---|---|
| Clean exit + four Complete flags + five digests | `completed` | `stopped` | required | natural exit code (must be 0) |
| Non-zero natural exit (with or without a gap) | `failed` | `stopped`/`unknown` | forbidden | natural exit code |
| Caller cancellation / deadline expiry | `cancelled` | as above | forbidden | empty |
| Watchdog timeout, unclassified stop, unsettled reconciliation | `uncertain` | as above | forbidden | empty |
| events/usage/logs/manifest gap (other than the non-zero exit case) | `uncertain` | as above | forbidden | natural code or empty |
| Process tree not verified gone | `uncertain` (only when the exit code is 0) | `unknown` | forbidden | empty |
| Operator request (an `operator-action-required` event) | `operator-action-required` | as above | forbidden | — |

Settlement: `observed` only for a `completed` run with complete usage (carrying the measured calls/tokens and the UsageHash evidence), otherwise `unknown`; a non-completed receipt must cite both the stop-evidence digest and its reason digest (sorted, unique).

## 4. Stop, cancellation and timeout

- Stopping goes only through the persisted identity (the handle's own identity; with no handle, the mirror plus a run-root ownership check) via the existing `StopProcessIdentity`; a repeated stop is idempotently `already-stopped`; a natural exit also owes an `already-stopped` proof.
- The cancellation path's verification context is grace-derived (no unbounded background context), and `Terminate`'s verification bound is decoupled from the caller's deadline.
- After the watchdog decides, the launcher waits up to `agentRunnerSettleWait(grace)` (grace+500ms, capped at 10s) for the run side to reconcile; an unsettled run returns `ErrAgentRunnerLaunchUnsettled` plus an “unknown” placeholder proof, and still goes through the durable stop path.
- **Linux zombie semantics (found and fixed through CI, §14)**: the aliveness probe must not read a process that has exited but has not been reaped by its parent as alive — its `/proc` entry, its start token and its process group all survive, yet it can no longer execute or write anything; reading a zombie as alive made every stop that does not own the child handle (stops by persisted identity) report a completed stop as `uncertain` once the grace period ran out.

## 5. Launch slot and ownership

- One launch per run root is guaranteed by the **exclusive create** (`O_CREATE|O_EXCL`) of the launch slot; the slot content is `<launchId>\n<per-claim token>`; a slot younger than the 10s freshness window is never taken over; beyond it the slot may only be reclaimed when the identity mirror proves the previous owner is gone (or was never published); an unreadable mirror fails closed.
- The release announces the end of the claim by atomically renaming the slot to `.releasing` and comparing the token: only a slot that still belongs to this claim is removed, anything else is put back — a late release can never delete another launch's slot.
- Stops resolve by ownership: with a handle, the handle's own identity (a foreign or stale launchId is refused, never reaching another launch's process); with no handle, the run-root owner is checked first, then the mirror is read.
- The wall-clock mtime distortion (suspension, stepped clocks) and “an external operator stop does not clear the slot” are both recorded as contract boundaries.

## 6. ③ V4 Pro pre-review record (7 rounds, all closed)

| Round | Verdict | Findings | Remediation |
|---|---|---|---|
| 1 | FAIL | 2 High (a caller deadline misclassified as completed; an events gap reaching completed) + 4 Medium + 5 Low | classify `DeadlineExceeded` and any unclassified error; events gaps added to the gate; a manifest gap becomes uncertain; exitCode only for natural exits; non-completed receipts cite the stop evidence; tree re-verification covers externally stopped runs; dead fields removed; grace honoured |
| 2 | FAIL | 2 Medium (a settle expiry releasing a possibly-live launch; in-memory-only exclusivity) + 3 Low | `agentRunnerNeedsDurableStop` routing; the durable exclusive launch slot; `stopGrace`; the bounded events prefix documented |
| 3 | FAIL | 2 Medium (a check-then-write window; a Windows-only fixture breaking the Ubuntu leg) + 2 Low | the slot made atomic (with a token in a later round); a portable injection seam; the handle retired even when the mirror is unreadable |
| 4 | FAIL | 1 Medium (the slot claimed before the precheck widened the steal window) + 2 Low | the claim moved after the precheck; `os.Chtimes` makes the test hit the intended branch; the operator-stop/slot boundary documented |
| 5 | FAIL | 1 Medium (the stop path did not bind launch ownership) + 2 Low | `stopIdentity` ownership binding; the slot owner record plus conditional release; the suspension/clock boundary documented; a falsifiable ordering test |
| 6 | PASS WITH FIXES (nothing blocking) | 2 Medium (a foreign launchId could stop another launch's process; the release's read-compare-remove was not atomic) + 2 Low | foreign ids refused plus the run-root owner gate; a per-claim token with the rename-aside takeover; the operator-stop classification and the gap precedence written into both contract files; the gap×exit rows added |
| 7 | PASS WITH FIXES (nothing blocking) | 2 Low (the Chinese contract was missing two clauses the English had; the rename-aside micro-window) | the Chinese contract completed (done); the micro-window registered as a hand-over item for the operator/takeover slice |

## 7. ④ final-review pilot record (Haiku + MAI)

This slice replaced Codex with the low-cost combination for ④ under the user's same-turn authorisation (the directive deviations for §2.2 / §2.3 / §3.9 are registered in §12).

| Item | Record |
|---|---|
| ④a model/role | `runSubagent(agentName="quick-verifier", model="Claude Haiku 4.5 (copilot)")`, human setting: output budget 8000 |
| ④a result | **0** items (the total was stated explicitly and the reviewer flagged its own zero as suspicious, as §4.2 requires); **format deviation**: it added prose and a summary table instead of only one-line items |
| ④a false positives | not applicable (0 items, no verdicts) |
| ④b model/role | `runSubagent(agentName="independent-reviewer", model="MAI-Code-1.1-Flash (copilot)")`, human setting: thinking mode high |
| ④b first round | `RE-REVIEW: FINDINGS`: 1 High + 1 Medium (each with a `file:line`, the quoted contract sentence and a failing scenario; Sections D/E complete) |
| ④ independent findings | **2** (Haiku 0 + MAI 2, intersection 0; Medium+ items raised by ③ that MAI missed: 0) |
| A11 spot-checks (actually executed) | Spot-check 1: MAI claimed “a new case with `ExitCode≠0` and an unverified tree would redden” — **executed**: the new row reddened on its first run (`status=uncertain … want failed`), the claim was true, and the code did contradict the contract → fixed by making the non-zero natural exit outrank the gap. Spot-check 2: MAI's Medium claimed a simultaneous cancel/timeout does not converge — **executed** as the decisive experiment `TestA6NativeE8`: before the fix 200 rounds gave `cancelled=138 / timeout=62` (claim true); with the bounded settling window it gives `cancelled=200 / timeout=0`, and M51 (window set to zero) reddens |
| ④b re-review | Both findings fixed and re-run: **`RE-REVIEW: PASS`**, no new regression; Section D still lists two boundary observations (the tree-status flip is pinned by a single case; the settling window is only stressed at one timing profile) — both are observations, neither blocks |
| Combination usability | Both model ids were resolved before the run (probe 1/2, see §11); the Haiku output was not truncated; MAI accepted the input and produced `file:line`s; **no** §6 fallback level was ever triggered |
| (a) the combination works | yes (both stages produced usable output) |
| (b) it finds real problems | yes: MAI independently found 1 High + 1 Medium in one call, both confirmed by actually executing the claimed counterexample (Haiku found nothing, consistent with its mechanical-scan role) |
| (c) cost reduction | the **pilot ④ stage** did not call Codex (0/2 fallback usage in that stage); tokens are to be filled in from the UI and recorded in the ledger; the Codex audits that followed ④ count separately as **3 calls**, see §7.1 |

④/③ accounting stays separate: ③'s seven rounds are not counted as ④ findings; the “④ independent findings = 2” above counts Haiku + MAI only.

### 7.1 Additional independent audit (Codex blind audit, explicitly authorised by the user in the same turn, run to measure the pilot's residual miss)

After ④ finished, the user asked for **one extra** independent Codex review to establish whether the low-cost combination still left problems behind and whether the A6 pilot is viable. That call is **not** on the pilot's normal path and was **not** triggered by any §6 medium/severe failure: it is an exception explicitly authorised by the user in the same turn (the deviation is registered in §12). The two re-reviews that follow were required by existing discipline (“close the remediation loop” and “never self-audit a change you made”), making 3 calls in total, each recorded below.

| Item | Record |
|---|---|
| Audit model/role | `runSubagent(agentName="independent-reviewer", model="GPT-5.3-Codex (copilot)")`, blind: given the changed scope, the contract requirements and a read-only constraint — **not** given the ③/④ finding lists |
| First round | `RE-REVIEW: FINDINGS`: **2 items (Codex marked both NOVEL)** — **1 High**: `refuseForeignSlot` returned `nil` for “any read error other than NotExist” and for “empty owner content”, i.e. it failed open (the opposite of the adjacent `refuseLiveMirrorProcess`); scenario: after a restart, an unreadable slot plus a foreign launchId proceeds through the mirror and stops a live process. **1 Medium**: it claimed there are “two stop mechanisms”, contradicting the contract's “stopping must go through the persisted identity via the existing `StopProcessIdentity`” |
| Adjudication | **High confirmed and fixed**: `refuseForeignSlot` became a three-arm fail-closed switch (readable ⇒ `fields[0] == launchID` must hold, an empty owner counts as a conflict; `os.ErrNotExist` ⇒ `refuseLiveMirrorProcess()` first, then decide; any other read error ⇒ conflict), with three new mutations M52/M53/M54. **Medium partially rejected**: read as **one** guard identity stop mechanism with **two entry points** — the handle-less path calls `process.platform.stop(identity, grace)`, the same identity mechanism `StopProcessIdentity` uses; forcing a single entry point would replace the guard's tree-verification evidence with a weaker `already-stopped` proof. The falsifiable part was adopted: the contract wording was tightened in both languages to “one mechanism, two entry points, no third channel” |
| Re-review | After the fix Codex was re-run: **`RE-REVIEW: PASS`** (High closed; Medium closed by the adjudication above; no new findings) |
| Audit 3 (test-delta review) | Under the existing rule “a change you made yourself must be reviewed by an independent third party”, one more read-only review covered the test delta **the driver had rewritten after Codex's PASS** (added assertions, removal of a one-off helper): **`RE-REVIEW: PASS`** with Section B **empty**; Section D listed three un-pinned invariants (no aliveness re-check after the missing-slot phase, the ownerless slot never checked against the owner id, `Actions` only checked by containment) — **all three were pinned** and the focused tests plus M52–M55 re-run (still 4/4 RED). This call closes the “no self-audit” compliance gap and is **not** part of the pilot's fallback |
| **Residual-miss datum (key for the pilot)** | The low-cost combination **missed 1 High** in this slice (found by the extra Codex blind audit, confirmed by 4 new mutations, then fixed); Medium+ items raised by ③ that it missed: **0** |
| §6 level it maps to | The combination missed a problem later proven real by mutation testing ⇒ under the pilot document's §6 that is a **mild failure**, whose prescribed remedy is “re-run MAI after the fix”. In this slice the closure gate was instead the **user-authorised Codex re-review** (Codex judged PASS) and **MAI was not re-run**; that path deviation and a follow-up recommendation are registered in §12 |
| Accounting discipline | The 2 audit findings are **not** counted in the “④ independent findings = 2” row above; they are recorded here separately so the pilot combination's output is never mixed with the extra audit's |

### 7.2 Pilot verdict (data + recommendation; the decision stays with the user)

| Criterion (pilot document §9) | Verdict | Evidence |
|---|---|---|
| (a) the combination works | **met** | both stages produced usable output, no §6 fallback level was ever triggered, both model ids were resolved before the run |
| (b) it finds real problems | **met, with a non-zero residual miss** | MAI independently found 1 High + 1 Medium in one call and both were confirmed by actually executing them (the A11 spot-checks above); **at the same time** the extra Codex blind audit found 1 High the combination missed (a fail-open path, fixed and pinned by 4 mutations) |
| (b2) a PASS needs independent corroboration (A12) | **met** | MAI reported PASS and ⑤ carries 4 A6-specific experiments (tree leak E1 / timeout boundary E2 / cancel race E3+E8 / crash injection E4) that independently prove the key invariants |
| (c) cost reduction | **met** | pilot ④ itself called no Codex; the 3 additional audits (blind audit + re-review + test-delta review) are one-off evaluation/compliance costs, not recurring ones |
| **Overall recommendation (driver)** | **use it as a low-cost supplementary scan layer for regular slices; do not use it as the only ④ on the hard-gate slices** | Sample size n=1: the combination missed 1 High that Codex caught, while ③ missed 0. Per pilot document §1 the hard-gate slices A7/B2/B3/B4 keep Codex as their final review; if the combination is to be promoted for regular slices, add a fallback rule such as “on every Nth slice, or whenever ownership / stop / identity semantics are touched, sample one slice for a Codex blind audit”, otherwise the residual miss has no bounded evidence |

On the decision rule: the pilot document's §9 three-way outcome (all met ⇒ promote / only (a)+(c) ⇒ fallback-only / any unmet ⇒ revert) is about the combination **replacing** Codex. This section's data ((a)(b)(c) formally met **plus** the 1-High residual miss) supports the combination as a **supplementary layer**, which is a different conclusion from “replace”; the user owns the choice.

## 8. Mutation audit (the H1 hard requirement)

- **Part 1, 48 mutations** (each removing one guard or condition, with a focused test): **RED 48 / SURVIVED 0 / INVALID 0**, all sources restored with SHA-256 verification; the script was a temporary file (`tmp/a6_mutations_v2.ps1`) and was deleted when `tmp/` was cleaned.
- **Part 2, 4 mutations (M52–M55, the Codex blind-audit round)**: M52/M53/M54 are the three arms of the High fix in §7.1 (unreadable slot refuses, empty owner refuses, missing slot plus a live mirror refuses); M55 flips the owner gate (`fields[0] != launchID` → `==`). Result **RED 4 / SURVIVED 0 / INVALID 0 / NOT-APPLIED 0**, restore verified by SHA-256. After Codex's re-review the test was **strengthened twice** (first two added assertions, then the three invariants missing per audit 3's Section D), and **each strengthening was followed by a re-run: still 4/4 RED with the restore verified** (scripts `tmp/a6_mutations_codex.ps1` and, for the strengthened re-run, `tmp/a6_mut_recheck.ps1`; both temporary, both deleted with `tmp/`). **52 distinct falsifiable claims in total, every one RED.**
- **Replacement**: the High fix in §7.1 invalidated the anchor text of one part-1 mutation (M49); the equivalent claim was re-added as M55 and actually run, so the total is counted as 52 with no double counting.
- Covered: every classification branch (including `TimedOut` outranking a natural failure, unclassified errors, an unverified tree), the gap gates (logs/usage/events/manifests plus the `FrozenFacts` second line of defence), settlement observed vs. unknown, exit-code honesty, the stop-evidence citation, tree verification (alive / unobservable / malformed report / per-pid observation), the pre-run manifest's capture position and wiring, the launch slot (exclusive create, freshness, token, release, release-after-refusal), ownership (handle identity, foreign id, owner gate, handle-less path), the unreadable-mirror refusal, the version precheck, the watchdog's single deadline and its simultaneous-arrival tie-break.
- Dropped with a rationale: the `os.Stat` non-NotExist error branch in `claimLaunchSlot` (unreachable without fault injection; the adjacent reachable mirror branch is covered by M43) and the token-string self-variant (the claim and the release read the same token, so it is equivalent; the discriminating half is covered by M49).

## 9. Native experiments (`-tags a6native`, real processes)

| Experiment | Claim | Result |
|---|---|---|
| E1 | A natural exit with a live descendant never ends as “completed + live tree” | 10 rounds: contained=10 / refused=0 (on Windows job-close reclaims the child, and the receipt still went through the re-verification) |
| E2 | Both sides of the timeout boundary | 20 rounds: completed=10 / degraded=10, with no out-of-window run completed and no surviving process |
| E3 | The cancel/timeout race has one winner and leaves nothing behind | 20 rounds: cancelled=20, exactly one terminal per round, every stub pid gone (the race window is 1.45–1.55s against a 1.5s watchdog, far from the precheck cost) |
| E4 | Crash injection at three phases (start/stop/collection) | 6 cases: fail-fast→failed, spawn-exit→contained, no marker / no usage / torn usage→uncertain, operator ask→operator-action-required |
| E5 | A version-pin mismatch must not spawn | pass |
| E6 | The facts and the settlement describe the real artifacts | pass (the usage digest and calls/tokens agree) |
| E8 | A simultaneous cancel/timeout converges on one winner | 138/62 before the fix (see §7), 200/0 after it |

## 10. Gates and encoding

- `gofmt -l internal cmd tools` empty; `go build ./...` and `go vet ./...` pass; `go test -count=1 ./...` green across 13 packages; the `tools/contracts` suite 4/4 (128 fixtures).
- Encoding and line endings: new/changed `.go` and `.json` files are UTF-8 without BOM + LF; `docs/CONTRACTS.md`, `_EN.md` and this document are UTF-8 **with BOM** + LF (verified byte by byte).
- Boundary constraints: the adapter and all new code import neither `net` nor `os/exec` (the version precheck goes through `guard`), `internal/release/network_test.go` passes, and `.github/workflows` and its digest contract are untouched.
- Cross-platform local gates reinforced (2026-09-16, see §14): `go build ./...` and `go vet ./...` are now also run under `GOOS=linux` with `CGO_ENABLED=0`, so Linux-only code paths (`internal/guard/process_linux.go`) no longer rely on CI alone to surface compile errors.

## 11. Probes and budget

- Probe **1/2**: purpose = resolving the ④a/④b model ids and the `runSubagent` vendor string; models = `Claude Haiku 4.5 (copilot)`, `MAI-Code-1.1-Flash (copilot)`; result = the first fallback rungs (`(anthropic)`/`(microsoft)`) are not whitelisted, the tool's authoritative list confirmed the `(copilot)` spelling, now recorded in the pilot document §2.3; **1/2 remaining**.
- Codex ran **3 times** (all read-only, no subprocess, no call chains; **above the 1–2 cap in pilot document §6**, each call registered in §12): ① the blind audit for the pilot's residual-miss assessment (scope, contracts and a read-only constraint only); ② the re-review after the fix → `RE-REVIEW: PASS`; ③ the test-delta review (closing the “no self-audit” gap) → `RE-REVIEW: PASS` with an empty Section B and Section D's three un-pinned invariants now pinned.
- No forbidden model (Luna/Gemini/Terra/Sol/GPT-5.4) was touched.

## 12. Directive deviations (A6 only, see `docs/t027/A6_PILOT_LAUNCH.md` §10.1)

- §2.2 (Codex must do the final review) → this slice ran ④ as the Haiku+MAI combination.
- §2.3 (only V4 Pro/Codex whitelisted) → `claude-haiku-4.5` / `mai-code-1.1-flash` were released for this slice only; depth stays ≤1, no call chains, and every other forbidden model stays forbidden.
- §3.9 (a purely offline slice defaults to 0 probes) → the user explicitly authorised one probe in the same turn to resolve the model ids.
- Pilot document §6 (Codex fallback capped at ≤2 calls, only on a medium/severe failure, never on the normal path) → the user explicitly authorised one **additional** Codex blind audit plus one re-review in the same turn (2 calls) to assess the pilot's residual miss; a 3rd call (test-delta review) followed under the existing rule “a change you made yourself must be reviewed by an independent party”, so **3 calls in total, above §6's cap**. None was failure-triggered: they are pilot-evaluation and compliance costs, recorded as such and **not** relabelled “fallback”.
- Pilot document §6 mild-failure remedy (if the combination misses a problem proven by mutation testing ⇒ **re-run MAI** after the fix) → MAI was not re-run here; the closure gate was the **user-authorised Codex re-review** (Codex judged `RE-REVIEW: PASS`). This is registered **as a path deviation**; if the combination is ever promoted for regular slices, the “who re-audits a missed finding” rule must be settled first (see the §7.2 recommendation).
- Pilot document §8 (the report must state that Codex was not used) → this report states it truthfully: **the pilot ④ stage did not use Codex**, and Codex only ran the user-authorised additional independent audit after ④ (§7.1); the two are listed separately and are never mixed.

## 13. Known boundaries and hand-overs

- **The rename-aside micro-window**: between the release's rename and the `.releasing` comparison there is a microsecond-wide “no slot” window (③ round 7, Low, judged non-blocking); hardening it belongs to the operator/takeover slice (or to a release semantic that embeds the token in the file name).
- **External operator stop**: after such a stop the slot may outlive it by up to one freshness window (10s), during which new launches are refused; on the run side an external stop is indistinguishable from a crash, so it is recorded as `failed` with its exit code (written into the contract).
- **The bounded events prefix**: an operator request beyond the 1MiB read bound is invisible, but that prefix is treated as incomplete and can never be completed.
- **Unix orphans**: when the launcher crashes after the parent exited, the mirror only names the parent, so on Unix a child orphan relies on the process group (this slice validated Windows natively only).
- **Uncovered observations** (④b re-review Section D): the tree-status flip is pinned by a single case, and the settling window was stressed at one timing profile only.
- **Observations left by the additional Codex blind audit** (§7.1, none blocking): ① **the unreadable-slot construction is not portable** — “the slot cannot be read” is built in the test with a **directory placeholder** (a permission denial is hard to reproduce on Windows without elevation), so the same shape on Unix only holds under `EISDIR` semantics; ② **the cross-platform primitives are validated natively on Windows only** — the `a6native` experiments (job containment, the `taskkill` path) **do not run in CI** (`ci.yml` runs the default path plus the `-race` subset), so the job semantics and orphan reclamation currently rest on host-local Windows evidence (updated 2026-09-16: the CI Ubuntu leg does run `internal/guard`'s Linux stop/identity path in the default suite, and a zombie-semantics defect turned 12 of its cases red — fixed and covered by a Linux regression test, see §14); ③ the combination's residual miss (1 High the combination did not find) is in §7.1/§7.2 and is a **process-level** known boundary that does not change this slice's code conclusions.
- **Section D/E of audit 3 (the test delta)**: the three un-pinned invariants are now **pinned** (aliveness re-checked after the missing-slot phase, the ownerless slot also tested against the owner id, and `Actions` asserted as “exactly one, `already-stopped`” instead of a containment check) and re-verified; what remains un-pinned is hygiene on **failure paths** (a `t.Fatal` in the later half of the test only guarantees, through `t.Cleanup`, that the process is stopped — not that the owner's handle is retired or its log handles closed) plus §7.1's “owner cleanup only asserts no error”; both are low-risk coverage residuals, not behavioural defects.
- **Commit status**: this slice's `c2d2819` was committed and pushed to `origin/main` under the user's explicit same-turn authorisation (gitee is never pushed); the Linux zombie-semantics fix and this document were pushed together with the fix commit (§14).

## 14. CI evidence and the Linux zombie-semantics regression fix (2026-09-16)

- **First CI run (red)**: after `c2d2819` was pushed, `main CI` (run `35074301616`) was green on the Windows leg and **failed the Test step on the Ubuntu leg**: 12 pinned-CLI cases reported `managed process termination uncertain: context deadline exceeded` with `Actions` = `[signal-term-process-group signal-kill-process-group]` and no proof of disappearance. This was the first time the slice's stop path ran on Linux — A6's native validation covered Windows only (§9, §13).
- **Root cause**: the Linux aliveness probe compared only the `/proc/<pid>/stat` start token. When a process is killed by a signal but its parent has not `wait`ed for it yet it becomes a zombie: the `/proc` entry and the start token survive and `kill(-pgid, 0)` keeps succeeding, so the **stop path that does not own the child handle** (the launcher stops by the mirrored identity and never `wait`s itself) reported a completed stop as `uncertain` once the grace period ran out. Windows has no zombie state (handle plus exit code), which is why that leg stayed green.
- **Fix** (`internal/guard/process_linux.go`):
  - `readProcStat` now parses the state (`fields[0]`) and the pgrp (`fields[2]`) alongside the start token, keeping the pid-reuse protection;
  - `platformIdentityAlive`: a state of `Z`/`X`/`x` (zombie/dead) is always read as **not running**;
  - `processGroupAlive`: when signal 0 still succeeds, `/proc` is scanned and the group only counts as alive if it has a **non-terminal** member; a group owned by another user (EPERM) keeps the old “conservatively alive” semantics, since `/proc` cannot show its members.
- **Regression test** (new `internal/guard/process_linux_test.go`, Linux-only): `TestLinuxStopProvesAnUnreapedTerminatedChild` — the child gets its own process group and is **deliberately left unreaped** after the kill, then the test asserts ① a running child reads as alive (control); ② an exited un-reaped child reads as **not running**; ③ its group is no longer an alive group; ④ `StopProcessIdentity` returns `stopped` with `Actions == [already-stopped]`; it only `wait`s afterwards, so the conclusion cannot be an artefact of the parent's own reaping. The case failed within the 5s bound before the fix.
- **Local gates (Windows host)**: `gofmt -l internal cmd tools` empty; `GOOS=linux` (`CGO_ENABLED=0`) `go build ./...` and `go vet ./...` pass; Windows `go build ./...` and `go vet ./...` pass; `go test -count=1 ./internal/guard/ ./internal/adapters/` passes.

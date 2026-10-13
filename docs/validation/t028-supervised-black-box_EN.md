# Slice validation report · `T028` (supervised-black-box explicit mode)

> Slice id: `T028`; directive v1.27; date 2026-10-11. Tier `[PILOT]`; size M-L; probes 0; `[ROUTINE]` not applicable.
> **Status**: **⑧b final + commit and CI writeback** (⑨ local validation and the §7.7 A12 dedicated independent experiment are in §7; ⑧c wrap-up records are in §11; the commit and CI facts are in §9). **Authorization boundary**: the commit and push ran under the user's **same-turn bundled authorization** (origin-only / SSH); gitee no; probes 0 (no additions).

## 1. Objective and scope

- Slice nature: feature implementation slice, tier `[PILOT]`, size M-L, **probes 0** (zero probe budget: the visible channel and the candidate are both stub-driven; a real-channel requirement escalates under §1.7.2 and was not triggered here).
- Deliverable objective: implement supervised-black-box explicit mode - versioned risk/liability notice + binding confirmation + isolated workspace + SessionBridge visible delivery + human supervision and explicit return + post-return full scan + independent gates/review + reduced/unknown report; **AT-24 fully closed** (conditional reading in §7.2 and §8).
- Four prohibitions (each must be mechanically verifiable): ① no automatic downgrade from AgentRunner; ② a visible receipt never impersonates completion (it produces no task PASS); ③ risk confirmation never waives a safety gate; ④ unknown cost/process/side effects are never displayed as safe.
- Out of scope (boundary stated): real SessionBridge visible end-to-end (depends on T027-B wiring), live observation of a real external agent (network/cost/process/side effects), `target.schema.json` registry-level glob comparison, the S1 release surface, `.github/workflows/**`, `go.mod`.
- Start-up rulings: PILOT confirmed; ⑦ budget `1 final + 1 re-review + 1 blind`; touching `postflight.go` / `models.go` was ruled additive and thus §1.7.1 execution freedom.

## 2. Change summary

- Change set = **commit `b095cd2`**: **28 items at the ⑨ point** (`+6998/−22`, 12 added / 16 modified); **34 items at the final state** (`+7544/−24`, also covering the report pair, the ledger note and the `DEV_PLAN` writeback, see the last three rows below), 0 untracked; one further **docs-only CI evidence writeback** follows (§9). Every artifact was **precisely staged** with `git add <path>`; `git add -A` was never used.
- Caliber definition (verified with `git show --stat` / `git diff --shortstat`): the **functional change set** = code / contract / schema / tests plus this slice's evidence list, **28 items / `+6998/−22`**; the **main commit `b095cd2`** = functional change set + closeout documents (two reports, two ledger rows, two `DEV_PLAN` entries), **34 items / `+7544/−24`** (matching CI `changed=34`); the **docs-only CI evidence writeback `049b71f`** = 4 items / `+18/−16`.

| Area | Artifact | numstat |
|---|---|---|
| Contract | `docs/CONTRACTS{,_EN}.md` | +16/−1 (each) |
| Contract | `schemas/chain.schema.json` | +33/−1 |
| Contract | `schemas/supervised-risk-confirmation.schema.json` | +166/−0 (added) |
| Contract | `schemas/supervised-report.schema.json` | +160/−0 (added) |
| Contract | `testdata/contracts/**` (4 items) | +605/−0 |
| Contract | `tools/contracts/{contracts.test.js,README.md}` | +3/−3 |
| Chain | `internal/chain/supervised_{blackbox,risk_confirmation,report}.go` | +819/−0 (3 added) |
| Chain | `internal/chain/{models,router,engine,postflight}.go` | +197/−12 |
| Chain | `internal/chain/supervised_blackbox_test.go` | +896/−0 (added) |
| Console | `internal/console/{runtime_supervised.go,supervised_cmd.go}` | +2102/−0 (2 added) |
| Console | `internal/console/{cli,config,preview}.go` | +115/−4 |
| Console | `internal/console/runtime_supervised_test.go` | +1549/−0 (added) |
| Adapters | `internal/adapters/sessbridge_visible{,_test}.go` | +290/−0 (2 added) |
| Evidence | `docs/validation/evidence/T028-freedom-list.md` | +31/−0 (added) |
| Report | `docs/validation/t028-supervised-black-box{,_EN}.md` | +266/−0 (added) |
| Report | `docs/t027/REMAINING_SLICES{,_EN}.md` | +6/−0 (completion row + classification boundary note) |
| Plan | `docs/DEV_PLAN{,_EN}.md` | +1/−1 (T028 entry checked off + report path) |

- Core semantic point: the return is the handoff `complete` receipt (no second return record is introduced); two transitions are added, `supervised-return-accepted` and `supervised-return-completed`; opening the window persists both the binding confirmation and a **policy snapshot**.
- Contract point: `docs/CONTRACTS{,_EN}.md` gains `### 7.1 supervised-black-box` explicit mode; the two new schemas force `assuranceLevel = reduced` and an `unknownDimensions` list that contains `cost`.

## 3. Pipeline and environment

- Pipeline (§5.1): ① architecture → ②a contract → ②b implementation → ③ tests → ④ gates → ⑤ first pre-review → ⑥ independent scan → ⑦ blind review → ②b repair rounds (②b-3 … ②b-8) → ③ re-run → ④ re-run → ⑤ re-review → ⑥ final scan → ⑦ final review → ⑧a this report → ⑨ native validation → ⑧b finalize → ⑧c wrap-up → ⑩ stop point.
- Roles and models: ① `deep-reasoner` / `DeepSeek V4 Pro (deepseek)`; ②a / ②b / ③ / ⑧ are carried by `DeepSeek V4.1 Flash (deepseek)`; ⑤ `GLM-5.3 (glm)`; ⑥ `MAI-Code-1.1-Flash (copilot)`; ⑦ `GPT-5.3-Codex (copilot)`, **not downgraded**.
- Environment: Windows host; repository `d:\LZProjects\prfrail` (the `prfrail` folder of a multi-root workspace); gate entry `node tools/gates/gate.js`; Go 1.22 toolchain.

## 4. Anomalies and fallbacks

- **D-14 (positive data point)**: ②b self-reported all green, but ③ independently ruled AT-24 clause 4 **not closed** (the scope / side-effect scans were missing) ⇒ ②b-3 filled them in and ③ was re-run. This event counts into the **our-side findings** numerator of the defect self-capture rate.
- **D-16 carrier outage and audit-first recovery**: the first ②b-5 dispatch died on `UND_ERR_SOCKET` (7 items already written, `+421/−40`) ⇒ the driver first took evidence into `partial-state.txt`, then re-dispatched as an "audit first, then repair" round; the audit found F-3 half-done and red on entry plus a real console-side gap ⇒ repaired, with mutation restoration proven by SHA256 equality.
- **D-17 environment incident**: the system proxy was ON (`http://10.0.0.246:8080`), which made the DeepSeek / GLM carriers unusable (direct connectivity was fine, the proxy path was the root cause) ⇒ recovered after the user turned the proxy off in the same turn; ⑥ / ⑦ ride the Copilot channel and were unaffected.
- **A-1 the ⑤ third call was not narrowed (declared honestly)**: that dispatch predated the user's "check only ②b-7" narrowing instruction ⇒ the full final-state input package was in fact delivered; the consequence was only restatements of registered/settled items plus two Lows, with no harm beyond a repeated probe.
- **A-2 the driver's first ④ re-run was misread as interrupted by the terminal**: one long command chain was misread as `^C` ⇒ it was split into short commands and **every reading was taken again**; the artifacts were re-checked and are complete.
- **A-3 evidence form for read-only carriers**: the ⑥ / ⑦ carriers have read-only tools only and cannot execute anything ⇒ every reading is **materialized** into the input package by the driver (`head-baseline.txt` / `staged-diff.txt` / `gate-*.txt`); this report therefore claims only a **text-level** review by ⑥ / ⑦.
- **A-4 first-round self-inflicted A12 criteria (declared honestly)**: **all** 5 first-round failures were defects of the script's own criteria (4 wrong JSON paths - the record body nests under `$defs`, while the schema itself is correct; 1 wrong assumption about the decoding function); after the criteria were corrected, `A12_EXPERIMENT_FAILURES=0`; that round **changed no artifact** to make the experiment pass.
- **A-5 report authoring order**: this report was written in one pass **after** the ⑨ / ⑧c readings existed (rather than writing empty values first and back-filling) ⇒ this deviates from the literal "⑧a → ⑨ → ⑧b" order of §5.1; **it must not be used as a precedent**.
- **A-6 the reports themselves join the change set**: the two report files, the ledger writeback and the `DEV_PLAN` writeback raise the final-state staged count from 28 to **34**; the ④ reading (28) and the ⑧c reading (34) are reported separately in §7.1 and §11.
- **A-7 the ⑧b writeback missed `DEV_PLAN` (found by the user's audit)**: the first report version judged `DEV_PLAN` "not applicable", but §11.2 requires a slice to write back `DEV_PLAN` at closeout and the OB-12 exception covers directive-governance slices only ⇒ it does not apply here (a code slice). **Corrected**: the T028 entries in `DEV_PLAN{,_EN}` are checked off and carry the report path, and §15 of this report is corrected accordingly.
- **A-8 the `local-only` wording of the §17 artifact fence (found by the user's audit; partial gap)**: per §1.5.1, process artifacts cited by a version-controlled document must not be treated as category ③; archiving is owned by the "tmp/ lifecycle execution" slice ⇒ this slice only registers candidate **C-15** (§13) and performs no archiving.

## 5. Review conclusions (⑤⑥⑦)

- **⑤ first read**: `PRE-REVIEW: PASS WITH FIXES` (1 High + 4 Medium + 3 Low; boundary compliant, contract consistent across three layers, P1–P4 hold).
- **⑤ re-review (2nd)**: F-1…F-8 came out as 7 CLOSED and F-5 PARTIALLY CLOSED; the code / schema / contract layers produced **zero new findings**; NEW-1 / NEW-2 / NEW-3 are registration and evidence-completeness items.
- **⑤ re-review (3rd, same-turn user-authorized exception)**: `PRE-REVIEW: PASS WITH FIXES` - the drift guard is real and correctly anchored, the `abort` exemption is acceptable, both ②b-6 repairs pass, the semantic P1 test is load-bearing, and all seven AT-24 clauses close under the (b)/(c) readings.
- **⑥ round 1**: `INDEPENDENT SCAN: FINDINGS`, but with a **format deviation** (no severities, no `file:line`, no cannot-verify list) ⇒ per §7.2 **one resend was permitted**; its three findings, on inspection, were **all restatements of registered / settled items**.
- **⑥ round 2**: `INDEPENDENT SCAN: PASS` (Sections A–E complete, 8 cannot-verify items); the format deviation was corrected.
- **⑥ round 3 (final scan, final state)**: `INDEPENDENT SCAN: PASS` - the findings table states **none**; the verdicts are `PASS` for boundary / contract / P1 / P2 / P3 and **conditional** `PASS` for P4 and AT-24; 3 cannot-verify items; the driver mechanically spot-checked its `file:line` citations and **every one of them resolved** (§6).
- **⑦ blind review (no list attached)**: `BLIND REVIEW: FINDINGS` - 1 Critical (= **restatement of ruling (b)**) + **1 High (a genuinely new finding: the return phase never verified the frozen policy ⇒ silent policy drift)** + 1 Medium (= **restatement of D-18**) + 1 Low (the P1 lexical evidence can be bypassed by an equivalent rewrite).
- **⑦ final review (freedom list attached)**: `FINAL REVIEW: PASS` - **0 new executable defects**; the High is judged **genuinely closed**; F-1…F-18 are all "in-bound / implementation choice" with **no overrun**; L-1 is closed and taking the (b) registration path for L-2 is acceptable; ⑥'s `PASS` is judged **credible but limited** (it never raised the drift dimension on its own).
- Disposition chain: the blind High was fixed by **②b-7** (policy snapshot + three-anchor guard + zero-persist path + `abort` exemption regression); the blind Low's convergence advice was satisfied by **semantic invariant tests**; the two ⑤ Lows were handled by **②b-8** (L-1 fixed, L-2 registered).

## 6. Review input packages (including blind-review isolation proof)

- **⑦ blind input package**: **only** "the change set (28 items) + the acceptance criterion (AT-24 + the four prohibitions + boundaries)"; the design book, the orchestration record, the freedom list and every earlier review conclusion were **withheld** ⇒ isolation holds (all 4 of its findings were independently derived).
- **⑦ final input package**: the final-state change set and the materialized readings + `PLAN-v1.md` + the orchestration record (§9.1–§9.11) + the **freedom list** (a closure round must attach it per §1.7.3 / OB-48).
- **Materialized readings list**: `head-baseline.txt` (36 items, per-path HEAD vs worktree sha256 / bytes / lines), `staged-diff.txt` (331325 bytes), `gate-{tree,index,ci,selftest}.txt`, `go-test.txt`, `npm-test.txt`, `numstat.txt`.
- **Driver mechanical spot-check (this does not replace a review)**: the ⑦ citations `TEST_STRATEGY.md:54` (the AT-24 row), `CONTRACTS.md:272` (§7.1 heading), `supervised-report.schema.json:133` (`reduced`) and `:150` (`cost`), `runtime_supervised_test.go:644` (drift negative), `:746` (`abort` exemption), `:1376` (lexical smoke), `:1417` (semantic invariant) **all resolved**.
- **Sanitization**: the review inputs contain repository paths and text only; no credentials and no external connection material.

## 7. Falsifiability and gate results

### §7.1 Mechanical readings

| Item | Reading |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `changed=28`, `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=index` | `changed=28`, `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=ci` | `changed=2`, `TOTAL_FAIL=0` (**not this slice's change set**, see §9) |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 335/335` |
| `go test -count=1 ./...` | `ok=14 FAIL=0` |
| `npm --prefix tools/contracts test` | `tests 4 / pass 4 / fail 0` |
| IDE diagnostics (three packages) | `internal/chain`, `internal/console`, `internal/adapters` ⇒ **0 errors** |
| `gofmt -l` / `go build` / `go vet` | empty / exit 0 / exit 0 |

- The table above was measured by ④ **after ②b-8** (final state); the ⑧c re-run taken after the reports, the ledger and the `DEV_PLAN` writeback joined the change set is reported in §11.

### §7.2 AT-24 clause verdicts

| Clause | Verdict | Implementation point and evidence |
|---|---|---|
| Display and persist the binding risk confirmation before start-up | closed | the `supervised-risk-confirmation` record + `SupervisedRiskNoticeVersion = 2`; `confirmAndOpenWindow` runs before any durable write; two negative cases |
| Isolated workspace only | closed | the isolated workspace root + `ChangeSetInScope` (reusing the same rule as `gates.ErrWorkspaceEscape`) |
| Delivery is not completion | closed | `DeliverVisible` produces no receipt record; `SubmitSupervisedReturn` accepts only the handoff `complete` receipt plus the five facts |
| Full manifest/diff after the explicit return | closed | the return path captures the before/after manifest and diff and persists them; `TestSupervisedAssemblyFullChain*` |
| Scope/secret/file/side-effect scans plus independent gates/review | closed ((c) condition) | `runScopeScan` / `runSecretAndTypeScan` / `runSideEffectScan` / `runAfterReturnGates`; the scope scan is a **bounded approximation** (registered) |
| Unconfirmed/offline/unknown process-cost-external-side-effect blocking | **conditionally closed ((b))** | anything detectable is blocked; anything unobservable in principle is labelled `unknown` and lands in the report. **Under the literal reading this clause is Critical** (see §8) |
| The report is marked reduced/unknown and no automatic downgrade is possible | closed | the schema forces `assuranceLevel = reduced` and a `cost` entry in `unknownDimensions`; the renderer emits `unknown` only |

### §7.3 Four-prohibition (P1–P4) evidence matrix

| Prohibition | Mechanical implementation point | Regression test |
|---|---|---|
| P1 no automatic downgrade | `models.go` intent guard + the `Definition.validate` case; `cli.go` mutual exclusion; `router.go` port branch | lexical check `:1376`; semantic invariant `:1417` (sixfold in total) |
| P2 delivery never impersonates completion | `sessbridge_visible.go` produces no receipt; `supervised_blackbox.go` accepts only the `complete` receipt | the delivery-only rejection negative in `supervised_blackbox_test.go` |
| P3 confirmation never waives the gates | scope / secret / type / side-effect / process / gates all run before the return completes | every one of the five side-effect blocking branches has a negative; `runAfterReturnGates` negative |
| P4 an unknown dimension is never displayed as safe | `ValidateSupervisedReport` rejects a non-`reduced` level, an empty `unknownDimensions` list and a missing `cost` | negatives at both the schema and the renderer ends |

### §7.4 Mutation bindings

| Mutation | Expectation | Observed |
|---|---|---|
| remove the drift check from the return entry | red | `TestSupervisedReturnRefusesDriftedPolicy` red |
| add the drift check back to `abort` | red | `TestSupervisedReturnAllowsAbortUnderDriftedPolicy` red, SHA256 restored byte-identically |
| break the `ErrSupervisedScopeUnderivable` branch | red | the matching negative red |
| break the ancestor-containment clause of `allows` | red | the matching negative red |
| put `independent-review` back into the write-time set (F-3) | red | red on entry, measured by the recovery round |

### §7.7 A12 dedicated independent experiment

- Trigger: ⑥ reported `PASS` ⇒ ⑨ must prove the load-bearing invariants with a criterion **independent of ⑥**; the experiment script is `tmp/t028/verify-a12.js` (a local process artifact) and contains static criteria plus one child-process re-run.
- Readings: `A12_CHECKS=58 A12_EXPERIMENT_FAILURES=0` (exit 0); the focused re-run of `go test -count=1 -run Supervised -v` over the chain / console / adapters packages ⇒ `81` `--- PASS` entries and `0` `--- FAIL` entries.
- Eight parts: ① the frozen-policy guard (three anchors + zero durable write + the `abort` exemption + the call from the return entry); ② contract invariants (`reduced`, a `cost` entry in `unknownDimensions`, `selection = confirm`, the third variant, the coverage constants); ③ P4 "an unknown must not render as safe"; ④ P2 "delivery does not impersonate completion"; ⑤ P3 "confirmation waives no gate" (all six pre-return checks are present); ⑥ the mechanical anchors and test registrations of the four prohibitions; ⑦ the independent re-run (tests / `go vet` / `gofmt`); ⑧ the D-18 fail-closed boundaries really exist.
- The 5 first-round failures were **self-inflicted script criteria** (see A-4) and all passed after correction; **no repository artifact** was changed to satisfy a criterion.

## 8. Known boundaries and residual risk

- **D-18 known boundary**: the **real CLI end-to-end of `prfrail supervised` is unreachable in this slice**, limited by two independent fail-closed boundaries: ① the visible channel is the `unavailableVisibleClient` placeholder; ② the assembly provides no hook-definition source while the schema forces a non-empty `hooksAfterReturn` ⇒ any `supervised return --complete` is refused with `ErrSupervisedGateUnavailable`. **Honest conclusion: this slice is not a runnable release**; the real end-to-end depends on **T027-B** wiring.
- **(b) conditional closure (wording conflict)**: the AT-24 clause "unknown process / cost / external side effects **block**" conflicts with §8 of the design book ("detectable ⇒ block; unobservable in principle ⇒ label `unknown`"); the user ruled §8 authoritative and moved the wording clarification to **a later documentation slice** (whose priority was raised to **the next slice**). **Under the literal reading this clause is Critical**, so the report must never claim an unconditional "AT-24 fully closed".
- **(AT24-WORDING correction note)**: this line is a correction note and does not alter the (b) line - that line is kept per §11.1. The wording clarification of the AT-24 clause has been completed by the `AT24-WORDING` slice: the AT-24 row of `docs/TEST_STRATEGY{,_EN}.md` now reads, per §8 of the design book, "anything detectable is blocked; anything unobservable in principle is labelled `unknown`", the acceptance verdict remains **conditionally closed ((b))**, and no other conclusion line of this report changes.
- **(c) conditional closure (the scope scan is an approximation)**: `AllowedTargets` are target-registry ids whose path globs live in `target.schema.json`, and neither the chain definition nor the console configuration **carries that table** ⇒ a glob-level comparison exceeds the `internal/console/**` boundary; the slice judged on "workspace containment + fail-closed on an empty allow-list" and registered it as a known boundary (candidate **C-10**).
- **L-2 registered boundary (no boundary markers in the digest)**: `supervisedHandoffPolicyHash` joins the parts with `\n` and digests them with **no field / list boundary markers** ⇒ moving an element across an adjacent list boundary can collide with the original value; there is **no enforcement bypass** (once the guard holds, the runtime uses the snapshot policy) ⇒ detection is merely weakened. Disposition = registration (b) + candidate **C-12**.
- **Reachability boundaries (D-12 / D-13)**: the file-type guard is reachable only through an in-process constructed manifest; through the production entry `ErrSupervisedScopeUnderivable` is refused **earlier** by `chain.New` as `ErrInvalidDefinition` (stricter, never looser) ⇒ double-level assertions.
- **(F) residual guards (registered as known coverage boundaries)**: `ErrSupervisedAbortOnly`, `ErrSupervisedNotOpen`, the three definition checks of `supervisedAssemblyStep`, the "managed identity has exited" proof of `processStopEvidence`, the nil runner / unregistered hook paths of `runAfterReturnGates`, and the non-regular-file skip of `runSecretAndTypeScan` ⇒ candidate **C-9**.
- **The link leg has no local red evidence**: the link-resolution escape branch of `allows` always SKIPs on this host (symbolic-link privilege is required) ⇒ its evidence depends on a host or CI leg that has link privileges.
- **C-13 / C-14 registered**: C-13 = `complete` can still overwrite store-local artifacts after an `abort` (**pre-existing**; chain state is unaffected); C-14 = the guard freezes the hook ids inside the policy while the hook **definitions** are supplied by the caller at return time (unreachable under the current fail-closed wiring) ⇒ it must be an input to T027-B.
- **C-11 known skew**: the confirmation fixture `riskNoticeVersion: 1` produced by this slice in `testdata/contracts/valid/runtime-records.json` while the production constant is already `2` ⇒ schema-legal but semantically behind; per ⑤'s advice it is folded into the AT-24 wording documentation slice.
- **(AT24-WORDING correction note)**: the C-11 note has been landed by the `AT24-WORDING` slice - the version skew between the confirmation fixture `riskNoticeVersion: 1` in `testdata/contracts/valid/runtime-records.json` and the production constant `2` (`SupervisedRiskNoticeVersion` in `internal/chain/supervised_blackbox.go`) is **kept as is**: this slice's boundary forbids touching fixtures, so it is only registered honestly here and the file stays unchanged; aligning the fixture later needs its own slice with contract-first work.
- **⑥'s credibility boundary**: ⑥'s final `PASS` does **not** cover the "frozen policy consistency" dimension (⑦'s blind review found it independently and ②b-7 repaired it) ⇒ it must not be read as ⑥ independently confirming that repair.

## 9. Explicitly not executed

- **Commit and push (done)**: commit `b095cd2` (`feat(supervised): add the supervised-black-box explicit mode (T028)`) is pushed to **origin** (SSH, `0135f85..b095cd2`, change size `+7544/−24`); **gitee not pushed**.
- **CI legs (done, green on the first run)**: run `38137422277` (push event, head `b095cd2`) concluded `completed / success`. Both legs' `Gates` step readings - Ubuntu (`repo=/home/runner/work/prfrail/prfrail`) and Windows (`repo=D:\a\prfrail\prfrail`) - show `SELFTEST: PASS 335/335`, `GATE REPORT scope=ci changed=34` and `TOTAL_FAIL=0`; package-level `ok` lines are **16** for Ubuntu (including the Race leg's `internal/adapters` and `internal/chain`) and **14** for Windows, and both legs cover the target packages `internal/chain` / `internal/console` / `internal/adapters` / `internal/gates`; `FAIL` lines are **0** on both. **No red on the first run** (so no §6.4 handling was needed).
- **CI evidence writeback commit (docs-only)**: `049b71f` (`docs(t028): write back the commit id and both CI leg readings (no rule change)`, 4 items / `+18/−16`) is pushed to origin (`b095cd2..049b71f`); its own run `38137996587` = `completed / success` (**confirmed green, and per the authorization no further record commit was created for it**).
- **Timing note on the local `--scope=ci` reading**: at the ④ point that scope read `changed=2` (the previous commit's path set), **not this slice's change set**; the real CI legs measured `changed=34` as recorded above.
- **Probes 0**: no paid / real-model probe was issued; the behavior of the real SessionBridge visible channel and of a real external agent was **not observed** (see D-18 in §8).
- **Not executed**: `go test -race` on this Windows host (the race leg is carried by the CI Ubuntu leg), a real hook-definition source assembly, and the real visible end-to-end (depends on T027-B).
- **Call counts**: ① 1; ②a 2; ②b 8 (including repair rounds); ③ 3; ④ 3; ⑤ 3 (capped, see §12); ⑥ 3 (capped); ⑦ 2 (blind + final); probes 0.

## 10. Deviations and freedom

- Same source as the freedom list, item by item in `docs/validation/evidence/T028-freedom-list.md` (F-1 to F-18); the verdict is **F-1…F-18 all in-bound / implementation choice, no overrun** (driver self-report plus item-by-item confirmation by the ⑦ final review).
- F-1 / F-2: the implementation shape and port naming of the P1 decoupling (in-bound).
- F-3: the write-time set and its containment relation with the chain closed set (slice pre-authorized).
- F-4: raising the risk-notice version constant to `2` (in-bound).
- F-5: contract over design book (the P2 reading, see D-1) (in-bound).
- F-6 / F-7 / F-8: parking through the `Failpoint` seam, documenting the failpoint as a production API, and deleting the dangling comment (in-bound / implementation choice).
- F-9…F-18: scan boundary lines, renderer labels, the fail-closed assembly placeholder, candidate registration shapes and so on (in-bound / implementation choice).
- **Deviations (declared honestly)**: A-1 (the ⑤ input was not narrowed), A-2 (terminal interruption forced re-reading), A-3 (evidence form for read-only carriers), A-4 (first-round self-inflicted A12 criteria), A-5 (report authoring order), A-6 (the reports join the change set); **A-1 and A-5 must not be used as precedents**.

## 11. ⑧c wrap-up cleanliness record

- Encoding check (**whole tree**, §5.1 ⑧c reading, `git ls-files -z` over every tracked file type): **tracked 3431 / text 3342 / binary 89**; **0 violations in the non-frozen main tree**; **the frozen evidence tree is counted separately at 9** (registered readings: B2's CRLF, B3b's BOM ×2, and 6 under `glm-eval/archive/**`).
- Encoding check (**change-set** reading): BOM / EOL classification over this slice's **34 changed artifacts** gives **0 violations** (UTF-8 without BOM + LF for `.go` / `.js` / `.json`, UTF-8 with BOM + LF for `.md`).
- Residue: untracked residue **0**; staged state **34 items** (including the two reports, the ledger writeback and the `DEV_PLAN` writeback), with no `??` line in `git status --short`.
- IDE diagnostics: **0 errors** on the changed artifacts.
- `node tools/gates/gate.js --all --scope=tree` re-run at the ⑧c point ⇒ `changed=34`, `TOTAL_FAIL=0`; `--scope=index` gives the same reading; `SELFTEST: PASS 335/335`; `npm test` 4/4. `G7-a` … `G7-d` flip from `INFO` to `PASS` (both reports are in scope and every artifact-fence row resolves).
- `tmp/t028/` holds one-shot process artifacts (§1.5.1 category ③); every ⑥ / ⑦ input and reading is re-verifiable inside that directory.
- **Ordering**: see A-5 (this report was written in one pass after the ⑨ / ⑧c readings were taken).

## 12. Cost and accounting

- R7.1 line by line (field shape = `{role, model, stage, startedAt, reason, retryOf}`): ① `DeepSeek V4 Pro (deepseek)` 1 call; ②a / ②b / ③ / ⑧ `DeepSeek V4.1 Flash (deepseek)` 15 calls; ⑤ `GLM-5.3 (glm)` 3 calls (capped: first read + post-②b-5 re-review + one same-turn user-authorized exception); ⑥ `MAI-Code-1.1-Flash (copilot)` 3 calls (capped: the format-deviating first round + one permitted resend + one final-state scan); ⑦ `GPT-5.3-Codex (copilot)` 2 calls (blind 1 + final 1), with the 1 re-review allowance **unused**; probes 0.
- **Budget-exception registration (§7.8.3 / §7.5)**: ⑤'s third call was a same-turn user-authorized exception that is **valid for this slice only and does not generalize**; its input **was not narrowed** as instructed (see A-1).
- Carriers (all `DeepSeek V4.1 Flash (deepseek)`): ②a 2 / ②b 8 / ③ 3 / ⑧a 1 / ⑧b 1 ⇒ **15** calls in total.
- **The cost vector lives in this section only and never enters the directive** (G6-2).

## 13. Candidate observations (not landed in this slice)

- **C-9**: add negatives for the supervised CLI-argument and truncated-run-root guards (next slice or a later hardening slice).
- **C-10**: once a target registry is reachable, add a **glob-level** authorized-scope comparison for the return, replacing today's approximation.
- **C-11**: the fixture's `riskNoticeVersion: 1` version skew, folded into the AT-24 wording documentation slice.
- **C-12**: harden the digest (label / count prefix per list), assigned to a later hardening slice.
- **C-13**: `complete` overwrites store-local artifacts after an `abort` (**pre-existing**) ⇒ a later slice should check the projection is non-terminal before `complete` entry.
- **C-14**: T027-B wiring warning - the hook **definitions** must be frozen with the window (unreachable under the current fail-closed wiring).
- **C-8**: sweep the repository for other stale count wording (carried over from an earlier slice).
- **C-15**: this slice's `tmp/t028/` process artifacts (`ORCH-NOTES.md` / `verify-a12.js` / `partial-state.txt` / `head-baseline.txt` / `staged-diff.txt` / `gate-*.txt` / `go-test.txt` / `npm-test.txt` / `numstat.txt`) are **cited by a version-controlled report yet still labelled `local-only`** ⇒ this conflicts with the §1.5.1 wording; archiving belongs to the "tmp/ lifecycle execution" slice (not to this slice's ⑧c), so it is registered as a candidate for that batch.
- **OB-82 sampling anchor**: this slice used a self-defined criterion ("the 4th slice after the last sampled slice"); the next slice should formally write it into §7.7 or §12.4.

## 14. Follow-up recommendations (including tier advice)

- Follow-up recommendations (advice only, not registered): C-8 … C-14 wait for registration as a batch when the next slice starts; the **AT-24 wording documentation slice has been promoted to the next slice**.
- Tier advice: for later slices touching code / concurrency / ownership semantics, ⑦ stays `GPT-5.3-Codex (copilot)` and is **not downgraded**; a documentation-only slice may drop to `DeepSeek V4 Pro (deepseek)` under §7.1.
- ⑥'s coverage gap (it never raised the drift dimension on its own) suggests the next slice should evaluate **layering the must-verify list** (which items are inherently owned by ⑦ / ⑨).
- The pilot metric fell short (see §16) ⇒ the next slice should evaluate the "self-capture rate" criterion and how deeply ⑤ should intervene.

## 15. ⑩ stop-point report (including the writeback checklist)

- Executed through **⑧b finalize + ⑧c wrap-up + user-audit closeout + commit and CI**; ⑨'s local part and the §7.7 A12 experiment have both been executed; the **⑩ stop point has been issued** and was lifted by the user's **same-turn bundled authorization** (commit + push + CI observation + one docs-only CI evidence writeback).
- **Writeback checklist** (§11.2 reading): validation report = **produced** (this file pair); `DEV_PLAN` = **written back** (the T028 entry is **checked off** and carries the report path); ledger = **written back** (completion row plus classification boundary note); ADR = **not applicable** (no model downgrade and no budget exception requiring one).
- Authorization boundary (as of the ⑩ stop point): **commit no / push no / gitee no**; afterwards an **origin-only** (SSH) commit and push were executed under the user's same-turn bundled authorization, **gitee no**, probes 0; the commit and CI facts are in the "Commit and push" and "CI legs" entries of §9.

## 16. Pilot notes (metrics and honest declarations)

| Metric | Target | Observed |
|---|---|---|
| Defect self-capture rate | ≥ 3/4 | **2/15 (short of target)** |
| ③ rework count | ≤ 2 | 2 (at the cap) |
| ④ rework count | ≤ 2 | 1 |
| ⑦ invalid rounds | 0 | 0 |
| Driver fallback count | ≤ 3 | 0 |
| Misses found by the user | 0 | 2 |

- **Finding distribution (honest reading)**: our side found 2 (the GAP-1 hard defect stopped by ③ and the F-3 half-done state caught by the recovery audit); the reviewers found 11 (⑤'s 8 first-round items + 3 re-review items); the user found 2 (CL-1, a fixture-count typo, and CL-2, stale README count wording) ⇒ 15 in total.
- **Honest conclusion: the self-capture rate of 2/15 is below the pilot target of ≥ 3/4** and must not be papered over by swapping criteria. Attribution: this slice introduced a brand-new mode at the M-L tier for the first time, and ⑤'s contribution concentrated on three classes of semantic / soundness problems - "a missing enforcement point for a contract invariant", "a recovery-matrix defect" and "evidence honesty" - which are exactly the blind spots of mechanical gates and implementer self-checks.
- **Driver self-check misses (not counted as user findings)**: the ③ task brief reused the self-reported counts of the first ③ report without verifying them ⇒ corrected in the orchestration record.

## 17. Artifact list

```artifacts
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/T028-freedom-list.md	repo
docs/validation/t028-supervised-black-box.md	repo
docs/validation/t028-supervised-black-box_EN.md	repo
internal/adapters/sessbridge_visible.go	repo
internal/adapters/sessbridge_visible_test.go	repo
internal/chain/engine.go	repo
internal/chain/models.go	repo
internal/chain/postflight.go	repo
internal/chain/router.go	repo
internal/chain/supervised_blackbox.go	repo
internal/chain/supervised_blackbox_test.go	repo
internal/chain/supervised_report.go	repo
internal/chain/supervised_risk_confirmation.go	repo
internal/console/cli.go	repo
internal/console/config.go	repo
internal/console/preview.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_test.go	repo
internal/console/supervised_cmd.go	repo
schemas/chain.schema.json	repo
schemas/supervised-report.schema.json	repo
schemas/supervised-risk-confirmation.schema.json	repo
testdata/contracts/invalid/core-config.json	repo
testdata/contracts/invalid/runtime-records.json	repo
testdata/contracts/valid/core-config.json	repo
testdata/contracts/valid/runtime-records.json	repo
tools/contracts/README.md	repo
tools/contracts/contracts.test.js	repo
tmp/t028/ORCH-NOTES.md	local-only
tmp/t028/verify-a12.js	local-only
tmp/t028/staged-diff.txt	local-only
tmp/t028/head-baseline.txt	local-only
```
> **Archive note (2026-10-13, A1-a)**: the `tmp/` items cited by this report's fence have been archived per §1.5.1 into `docs/validation/evidence/directive-history/t028/` (with `MANIFEST.md` and `SHA256SUMS.txt`); the original paths are not rewritten, and the correspondence is recorded in that unit's `MANIFEST.md`. This is a follow-on to the v1.23 archiving execution (run at the `[ROUTINE]` tier, without bumping the directive version).


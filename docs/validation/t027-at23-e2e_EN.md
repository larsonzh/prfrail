# T027 · AT-23 mechanism-leg end-to-end evidence pack (B4)

> Status: **complete — the B4 deliverables are in place, the ④ gate is green in all three scopes and the ⑤/⑥/⑦ reviews are closed; AT-23 is not passed**. ③ finished the end-to-end re-run and the independent recheck from the pack's own bytes, the ⑤ pre-review findings F1 to F5 are remediated and re-verified by ⑤, and the ⑦ final-review findings F1 to F5 plus R-1 to R-4 were disposed of under the user's 2026-09-23 rulings ⒜/⒝/⒞; the five §8 entries are of mixed nature — deviation 1 is repaired output (integrated into this slice's code and report), deviations 2, 4 and 5 are recorded measurements and deviation 3 is a boundary measurement. The review summary is in §9.1, the OB-1 observation record in §9.2 and the gate conclusions in §9.4.

## 1. Conclusion (four questions, kept apart)

| Question | Conclusion |
|---|---|
| Is B4 complete | Yes (the ② deliverables are in place, ③ finished the end-to-end re-run and the independent recheck, the ④ gate is green in all three scopes and the ⑤/⑥/⑦ reviews are closed; see the review summary in §9.1) |
| Is AT-23 fully passed | **No** |
| Which faces were skipped | See the five-face coverage matrix in §7 (per face: ran, skipped, reason, alternative evidence) |
| Is the alternative evidence for the skipped faces sufficient | **No** - ⑦ judged all five insufficient under the equivalent-coverage-of-production-semantics reading (⑤ judged them sufficient under the honest-fail-closed reading; the disagreement is registered as OB-23); **this slice claims no equivalent coverage** |

**AT-23 status statement (must not be dressed up)**: AT-23 is **not passed**. What has been verified is the part of the mechanism leg that can be executed honestly on this platform; starting a real AI candidate and making billed calls is **measurably unreachable** here (the B1/B2 frozen conclusion) and is **not "not attempted"**. **Do not** write "not passed" as "partially passed".

## 2. Input summary

| Input | Purpose | Digest |
|---|---|---|
| `tmp/b4/input/b4-chain.json` | F1 to F4 chain definition (`code` plus `isolated` step) | `bb8c100d8d36387c40d7966044c50c994ccd45833647b0fa144a8efb7369361a` |
| `tmp/b4/input/b4-chain-hook.json` | F5 chain definition (`code` step plus `verify` hook step) | `c7ad22c5741dbc02d5b35697898b5cf8a85a89b2fb7b5a04e1398ccf9c795237` |
| `tmp/b4/input/b4-inputs.json` | Declared inputs: identity binding plus parent snapshot digest | recorded as `declaredInputs.sha256` in the pack |
| `tmp/b4/bin/prfrail.exe` | Real CLI (the `validate` and `preview` legs) | `204aa5bea0b2fbd493f4522b2a9ef08b025db44a61deeacf0fd1428e8a54d824` |
| `tmp/b4/bin/agent-stub.exe` | Managed workload (pinned to `agent-stub-0.1.0`) | `01ac99b6ef5684785e90187fa90f78c09f91b9805df63d10a07e18e1ed47314d`, and `-version` echoes `agent-stub-0.1.0` |
| `tmp/b4/bin/b4-e2e.exe` | Assembly harness (the ⑦ review object) | `2f376099b437e664da362533efe7a3e607ea043ea2b4d74ababae30c704e2cab` |

The digests bind the three binaries that `SHA256SUMS.txt` pins inside this pack's `tmp/b4/bin/**`; re-running the driver rebuilds them and yields different digests (Go builds are not byte-reproducible here), so the digests are valid for this run only.

## 3. Environment summary

| Item | Value |
|---|---|
| Operating system | `Microsoft Windows [版本 10.0.26100.9457]` (raw `cmd /c ver`) |
| Host | Windows 11 Home Chinese edition, build 26100, 64-bit; 11th Gen Intel Core i7-1165G7, 8 logical cores; 16801927168 bytes of visible memory |
| Go toolchain | `go version go1.27.0 windows/amd64` |
| git HEAD | `d76ca7e74ed15499699e4115c51bb0f8f4aa9b95` (measured at run time; ③ staged and committed nothing) |
| Platform scope | Windows hosts only; **native Linux validation belongs to C1 and is neither collected nor inferred here** |

## 4. Commands and cwd

The **exact command text plus cwd plus exit code plus stdout/stderr** of every invocation is recorded by the driver: the `invocations` array of `tmp/b4/pack-verdict.json` carries the command, the cwd, the exit code and each stream's file path, sha256, byte count and first 600 characters; the raw text is on disk as `tmp/b4/evidence/<case>/<label>.stdout.txt` and `.stderr.txt`. The invocation shapes are as follows.

| Step | Command shape | cwd |
|---|---|---|
| Build | `go build -o tmp/b4/bin/<artifact>.exe <pkg>` | repository root |
| CLI leg | `tmp/b4/bin/prfrail.exe validate` or `preview --chain <chain.json>` | repository root |
| Declared inputs | `tmp/b4/bin/b4-e2e.exe -emit-record -chain <chain.json> -workspace <seed> -out <inputs.json>` | repository root |
| Per face and phase | `tmp/b4/bin/b4-e2e.exe -face <fN> -case <pos or neg> -run-root <dir> -inputs <json> -chain <json> -stub <exe> -phase <run or park or recover or forge or halt or resume> -verdict <json>` | repository root |

## 5. Case statistics

The ③ re-run results follow (verdict classes come from the `verdict-*.json` files; the chain, task and step states come from the same files and are cross-checked against `runs/<case>/events/state-events.jsonl` line by line).

| Face and leg | Forward verdict class | Reverse verdict class | Control rerun agrees | Note |
|---|---|---|---|---|
| F1 forward | `completed` (chain `COMPLETED`, task `PASSED`) | — | — | five frozen facts present and mutually distinct; the workspace diff is empty |
| F1 reverse · unaccounted write (`f1/neg2`) | — | `uncertain-or-paused` (task `REPAIR_PENDING`, chain `PAUSED`) | agrees (both runs share the verdict class and the `REPAIR_PENDING` / `postflight-rejected` event-log transition) | the true falsification leg added by ③, see §6 and deviation 1 in §8 |
| F1 control leg · accounted write (`f1/neg`) | `completed` (task `PASSED`) | — | agrees (both runs `completed`) | the same bytes written inside the managed run, so the frozen facts already contain them; see deviation 2 in §8 |
| F1 boundary leg · capture-excluded path (`f1/neg3`) | `completed` (task `PASSED`) | — | single run (the cause is a compile-time exclusion set, not a race) | the unaccounted write lands under `tmp/**`, so the rescan cannot see it; see deviation 3 in §8 |
| F2 | `uncertain-or-paused` (chain `PAUSED`) | `uncertain-or-paused` | agrees | the timeout fires near the 5000ms bound (a range assertion); both legs assert the reported pids are gone |
| F3 | `uncertain-or-paused` (`logsComplete=false`) | `completed` (control leg, `logsComplete=true`) | agrees | the gap criterion is calibrated in both directions, not an all-reject |
| F4 | `uncertain-or-paused` (zero relaunch after the restart) | `refused` (zero writes) | agrees (both runs share the verdict class and the refusal reason) | the forged resume terminal was refused |
| F5 | `completed` (the whole chain ran) | `uncertain-or-paused` (task `REPAIR_PENDING`) | agrees (both runs share the verdict class and `REPAIR_PENDING`) | postflight rejected the tampered workspace |

③ independent recheck (the ② summary is not taken on trust): `tmp/b4/SHA256SUMS.txt` is recomputed file by file, so the declared set and the on-disk set are equal with zero mismatches (`driver-stdout.txt` and `driver-stderr.txt` are written by the shell that invokes the driver, whose handle is still open while the driver hashes, so both are excluded by name and the exclusion is declared in `pack-verdict.json`'s `excludedFromChecksums`); each face is read from the artifacts under `runs/<case>/**` - event log, pid reports, identity mirror - instead of from the `verdictClass` field. The in-pack re-run evidence holds face by face: the six re-run groups `f1-neg-r2`, `f1-neg2-r2`, `f2-neg-r2`, `f3-neg-r2`, `f4-neg-r2` and `f5-neg-r2` each reproduce the verdict class of their own first-run leg, so all five faces carry in-pack re-run evidence. The only single-run leg, `f1/neg3`, is the boundary measurement (the cause is a compile-time exclusion set, not a race) and face F1 carries two further re-run groups, `f1-neg-r2` and `f1-neg2-r2`; no leg was therefore demoted to a "**suspected reverse**".

Byte determinism of a re-run (⑤ remediation F1): `fillProjection` used to write `taskState` and `stepState` through a single slot inside a `for … range map` loop, so Go's randomised iteration order both made the verdict bytes irreproducible and silently dropped states (measured on f5-pos, where only the verify step survived). It now writes every entry in sorted key order (`key=state; key=state`). Comparing the two legs of each identical-shape pair (each on its own run root) field by field, the differences fall only on genuinely time-varying fields: `startedAt` and `endedAt`, `journal[*].at`, `mapping.elapsedMs`, the process identity (`pid`/`startToken` in `journal[].detail.processId`), the digests of the timestamp-carrying manifests and of the pid-carrying `stub.pid`/`stub-child.pid`, the `requestRecordHash` bound to measured inputs and to the run instant (`notes[1]`), the `eventLogHash` over clock-generated event ids, and the `checks[].detail` and `candidateHash` values that quote those digests; `verdictClass`, `chainState`, `taskState`, `stepState`, `observations`, `diff` and every assertion's `pass` value are byte-identical.

A reverse case that cannot be reproduced **may only** be recorded as "suspected reverse" and never as "reverse".

## 6. Failure evidence

| Case | Injection | Observable failure point | Artifact and digest |
|---|---|---|---|
| F1 reverse (`f1/neg2`) | after the managed process had exited and the terminal had been published and routed, the experimenter writes the byte-identical `mutate-workspace` file into the workspace | the frozen diff stays empty; postflight returns `proofrail:b4-postflight-manifest-mismatch:1` (detail is that path); the event log records `REPAIR_PENDING` / `postflight-rejected`; the task never reaches `PASSED` | `tmp/b4/evidence/f1-neg2/verdict-run.json` (the journal carries `injection` and `postflight:failed`) and `tmp/b4/runs/f1-neg2/events/state-events.jsonl` |
| F1 boundary leg (`f1/neg3`) | the same injection placed under the workspace's `tmp/` subtree (the snapshot capturer excludes `tmp/**` by default) | no failure point: the postflight rescan cannot see it and the task reaches `PASSED` - that is the boundary of the unchanged-workspace claim, registered as a deviation | `tmp/b4/evidence/f1-neg3/verdict-run.json`; the cause is the default exclusion set in `internal/snapshot/capture.go` |
| F2 forward | Load `-mode hang`, watchdog bound 5s | mapped to `uncertain`, no frozen facts, five error-evidence digests | `tmp/b4/evidence/f2-pos/verdict-run.json` |
| F3 forward | Load `-mode run -no-marker` | the log completeness gap degrades the outcome to `uncertain` | `tmp/b4/evidence/f3-pos/verdict-run.json` |
| F4 reverse | Forged resume terminal with a made-up priorCompletionHash | `ErrAgentRunnerResumeContinuityUnproven`; the event log digest is equal before and after | `tmp/b4/evidence/f4-neg/verdict-forge.json`, plus the re-run leg `tmp/b4/evidence/f4-neg-r2/verdict-forge.json` with the same verdict class and refusal reason |
| F5 reverse | A marker file written into the workspace after the terminal was frozen | postflight returns `failed`, the task goes to `REPAIR_PENDING` and the chain to `PAUSED` | `tmp/b4/evidence/f5-neg/verdict-resume.json`, plus the re-run leg `tmp/b4/evidence/f5-neg-r2/verdict-resume.json` with the same verdict class and `REPAIR_PENDING` |
| Whole slice | The production dispatch path is refused by a platform constant on Windows | the fail-closed text is recorded, see appendix A | the first `dispatch` entry of each verdict `journal` |

## 7. Five-face coverage matrix (ran, skipped, reason, alternative evidence)

| Face | Leg | Status | Reason | Alternative evidence |
|---|---|---|---|---|
| F1 isolated workspace | Assembly harness starts the load through the production launcher, facts complete, workspace reconciled; plus the unaccounted-write falsification leg and the capture-exclusion boundary leg | ran | — | — |
| F2 timeout | Watchdog timeout decision and whole-process-tree stop | ran | — | — |
| F3 log gap | Completeness-marker gap degrades (with the control leg) | ran | — | — |
| F4 unknown recovery | Zero relaunch after the restart; the forged resume terminal refused with zero writes | ran | — | — |
| F5 rescan after exit 0 | postflight and gate hook and freeze and review and promotion | ran | — | — |
| Whole slice | Production dispatch path (`AgentRunnerReplayDispatcher` publishes R) | skipped | platform constant: replay publication durability is `unproven` on Windows | **non-equivalent substitute**: after recording the fail-closed text the spawn goes through the same production launcher; see appendix A |
| Whole slice | Production terminal publication (`AgentRunnerTerminalPublisher`) | skipped | same root cause, the request record is missing | **non-equivalent substitute**: the terminal chain is built with the public record constructors plus `ToChainAgentRunnerTerminal`; the chain-side routing is unchanged |
| Whole slice | Production gate port pair (`gates.ChainPort` plus `GuardExecutor`) | skipped | capability declaration mismatch: the `memory limit` gap makes any hook unexecutable | **non-equivalent substitute** (**policy port uncovered**): the harness runs a real gate command over the same `guard.RunManaged` boundary and keeps the evidence |
| Whole slice | Starting a real AI candidate and billed calls | skipped | frozen by B1/B2: no candidate can be started legally and the billing evidence is blocked by §8.3 and §8.4 | **non-equivalent substitute**: a transparently declared deterministic CLI is the managed workload; **no claim of candidate compatibility** |
| Whole slice | Native Linux validation | skipped | this slice's platform scope is Windows; native Linux belongs to C1 | **non-equivalent substitute**: not collected, not inferred (given in place, see §9.6) |
| Whole slice | SessionBridge `visible` human interaction | skipped | belongs to T028 or AT-24 | **non-equivalent substitute**: zero cover in this slice (given in place, see §9.6) |

**This slice claims no equivalent coverage**: the alternative evidence in the rows marked skipped above is **not** equivalent coverage of the production semantics, only the evidence this slice can honestly collect; the disagreement about how sufficiency is judged is registered as OB-23 (directive §12.4).

## 8. Known defects

| Id | Status | Note |
|---|---|---|
| DR-6 | observed (does not block B4) | `internal/gates` failed once on a first run and passed on every rerun; registered only |
| DR-2 to DR-5 | fixed (DR-FIX slice) | see `docs/validation/dr-fix-enforcement-proxy.md`, not restated here |
| Deviation 1 (registered by ②, fixed by ③) | repaired output | The original F1 reverse case was not a falsification test: the write happened inside the managed run, so the frozen facts described the mutated tree, postflight passed and the task `PASSED` (both runs). ③ added `f1/neg2`: the same bytes are injected by the experimenter after the terminal was published, so the frozen diff is empty, postflight returns `manifest-mismatch`, the task goes to `REPAIR_PENDING` and the chain to `PAUSED` (both runs share the verdict class and the event sequence) |
| Deviation 2 (③ recheck) | recorded as measured | `f1/neg` is kept as the control leg: the same bytes written inside the managed run are part of the diff and pass. With `f1/neg2` it forms a controlled pair proving postflight is a real discriminator and not an all-reject |
| Deviation 3 (③ boundary measurement) | recorded as measured | `f1/neg3`: the same write placed inside the capturer's default exclusion set (`internal/snapshot/capture.go` excludes `.git/**`, `.prfrail/**` and `tmp/**` by default) is invisible to both the diff and the rescan, so the task `PASSED`. **Claim boundary (must be declared)**: the "workspace unchanged" claim means "**unchanged outside the capture exclusion set**", and `f1/neg3` is exactly the measurement showing that a write inside the exclusion set is not detected. This slice does **not** change that exclusion set; whether the boundary is acceptable is **left to ⑦** |
| Deviation 4 (③ recheck, rewritten as three parts by ⑤) | recorded as measured | (a) ①'s original assumption: zero relaunch is judged by a non-empty `launches/` receipt directory; (b) the measurement: `launches/` is empty - normal under the A2 single-winner mechanism, not a failure; (c) the corrected criterion: zero relaunch is carried by the **dispatch-event count** (always one) **and the identity mirror** (the managed-process identity mirror is byte-equal to the `pid`/`startToken` published by the park-phase spawn). Whether that criterion is acceptable is again **left to ⑦** |
| Deviation 5 (③ recheck) | recorded as measured | Under the `spawn` load the descendant overwrites `stub.pid`, so `stub.pid` and `stub-child.pid` degenerate to the same pid in the tree evidence. The tree claim loses no force (the verified pid is exactly the long-lived descendant; see the unmanaged calibration in `tmp/b4/probe/f2-tree-calibration.txt`), but the parent pid cannot be read from that artifact and is carried by the stop proof and the job object instead |

## 9. Execution pipeline, review, gates, cost and unrun items

### 9.1 Review summary (⑤ pre-review, ⑥ independent scan, ⑦ final review)

- **⑤ pre-review (1 round)**: reported 5 findings F1 to F5, all re-verified as **CONFIRMED FIXED** (code and artifacts); it also reported N1 Medium (report §2 digest versus pack manifest) plus three Low findings N2 to N4, all fixed; N1 was closed by a deterministic proof (three binary digests byte-identical across CN, `_EN` and `SHA256SUMS.txt`, 0 mismatches) instead of a paid round.
- **⑥ independent scan (1 round, MAI-Code-1.1-Flash)**: `INDEPENDENT SCAN: PASS`, **0 findings**; it carried the conclusion line plus sections A to E, with substantive D (per-guard falsifiability audit and uncovered invariants) and E (what it could not verify and what would overturn it).
- **⑦ final review (round 1, GPT-5.3-Codex)**: `FINDINGS` — F1 Critical (contract internal ordering conflict), F2 High (publication legs on a non-equivalent substitute), F3 Medium (gate port pair not equivalently covered), F4 Medium (B-1 boundary needs contract-layer wording), F5 Low; R-1 ruled that the contract wording needs revision, which triggered the pre-declared hard stop and an escalation; remediated under the user's 2026-09-23 rulings ⒜/⒝/⒞.
- **⑦ re-review (round 2)**: `FINDINGS` — 2 Medium (the same issue in CN and `_EN`: the scope of the "no code" wording in §9.7); four claims accepted (the ordering is now consistent with no residual contradiction anywhere in the contract, the skip-row labelling does not overclaim, the B-1 and R-3 boundaries are not dressed up as contract consensus, and AT-23 is still recorded as not passed). The sentence was re-scoped per the minimal fix and the finding was closed deterministically under the OB-11 exception (**no round 3 of ⑦ was run**, registered as OB-24).
The ③ independent recheck conclusions are in §5.

### 9.2 OB-1 observation record (fixed subsection; omitting any row leaves the write-back incomplete)

| Observation | Record |
|---|---|
| Did the ⑥ first output carry the conclusion line and sections A to E | Yes — `INDEPENDENT SCAN: PASS` plus all five sections: A a one-line conclusion, B `(none)`, C the evidence used, D per-face falsifiability conditions, E what it could not verify |
| Did it produce executable findings (`file:line` plus contract sentence plus minimal fix) | No — **0 findings**; within the same slice ⑤ and ⑦ both produced substantive findings |
| Evidence for closing or escalating OB-1 | ⑥'s output was **fully format-compliant** (weakening the documentation-only-task attribution) yet carried no finding; the sensitivity gap against ⑤/⑦ is on record ⇒ OB-1 stays observed, neither closed nor escalated |

### 9.3 Falsifiability

| Face | Condition that would refute the claim |
|---|---|
| F1 | a workspace change outside the set the load wrote, with postflight not rejecting it |
| F2 | a hanging load judged `completed`; or a descendant tree still alive past the verification window |
| F3 | a load without the marker still judged `completed` |
| F4 | a second dispatch after the restart; or a forged resume terminal accepted |
| F5 | exit 0 producing a task `PASSED` directly; or `REVIEW_PENDING` or `PASSED` reached after the tamper |

### 9.4 Gate results

`②` (measured after the ⑤ pre-review remediation re-run): `gofmt -l .` empty; `go build ./...` and `go vet ./...` exit 0; `go test -count=1 ./...` with no `FAIL` package; `node tools/gates/gate.js --selftest` reports `SELFTEST: PASS`; `node tools/gates/gate.js --all --scope=tree` and `--all --scope=index` both report `TOTAL_FAIL=0` (`G7-b` is green too: this slice's artifacts are staged); `G3-a` reports `PASS` for this bilingual pair (line count, heading positions, bold markers, pipes and blank lines all equal, and both sides' added/deleted counts equal); `G6-1` to `G6-4` have **zero hits** on the added lines; the `G4b` difference set is empty (the characters `庭`, `捉` and `悖` were registered through the R7 mechanism in `tools/gates/cjk-newwords.txt` with their reason and context excerpt); `run-faces.ps1` re-parses with `[Parser]::ParseFile` reporting zero errors.
`③` (kept for comparison): the four base gates and `--selftest` gave the same results; `--all --scope=tree` reported `TOTAL_FAIL=1` before this remediation round, the single failure being `G7-b` (the declared repo artifacts were not staged yet, which is ④'s job); `G6-3` and `G6-4` had zero hits on the draft and the `G4b` difference set was empty. Also measured: `git cat-file -e 'HEAD:tools/b4-e2e/**'` reports `fatal: path does not exist`, so the artifact block uses one concrete path per file (G7-b resolves repo rows with `git cat-file -e <ref>:<path>` in ci and range scope, where a glob form never resolves; ③ recheck: the block's repo rows correspond one to one with this change set and nothing is missing); `tmp/b4/**` stays `local-only` and G7-b never resolves it against the repository.
`④` (measured on 2026-09-23 after re-staging, including the contract-wording revision and the report-annotation fixes): `gofmt -l .` empty; `go build ./...` and `go vet ./...` exit 0; `go test -count=1 ./...` with no `FAIL` package; `node tools/gates/gate.js --selftest` reports `SELFTEST: PASS`; `--all --scope=tree`, `--scope=index` and `--scope=ci` **all report `TOTAL_FAIL=0`** (tree scope changed=14, `G7-b` resolves all 30 artifact entries, `G6-1` to `G6-4` have zero hits on the added lines and the `G4b` difference set is empty); `G3(a)` reports symmetric added/deleted counts for all three bilingual pairs; the independent positional mirror check (master-side script) reports 1070/1070 for the directive pair and 228/228 for this report pair, with all five shape properties equal and 0 per-line mismatches (`CONTRACTS` is historically not a positional mirror and `G3(a)` only requires equal added/deleted counts for it, which +4/-4 satisfies).

### 9.5 Cost accounting

Budget: ① times 1 + ⑤ times 1 + ⑥ times 1 + ⑦ times 2; realized: ① times 1 + ⑤ times 1 + ⑥ times 1 + ⑦ times 2 (one final review plus one re-review for this slice; the contract-wording revision is counted against its own separate quota and does not consume this slice's).
**The contract revision's derivative cost (a data point for future reference)**: the 2 Medium findings of the ⑦ re-review (round 2) and this `CONTRACTS` §7 wording revision **share one cause** - the document wording the revision introduced itself consumed one further review round. **Suggested**: besides counting a contract revision against its own quota, budget one extra review round for it (measured here: one revision, one derivative review round).

### 9.6 Unrun items and next steps

| Item | Owner | Blocking condition |
|---|---|---|
| Production dispatch registration (R, launch receipt, terminal publication) | later slice | needs platform publication durability or a contract-level substitute; this slice only records the fail-closed text |
| Making the production gate port executable | later slice | the gap between `GuardExecutor`'s capability declaration and `requireCapabilities` |
| Native Linux validation | C1 | — |
| Real candidate compatibility and billing evidence | later slice | DR-1 and B2 §8.3 and §8.4 |
| SessionBridge `visible` human interaction | T028 or AT-24 | — |
| Driving an AgentRunner step from the product layer (console assembly) | later slice | this slice may not change product code |

### 9.7 Order accounting (the R5 item handed over by A5)

Measured after the F5 forward run (from `runs/f5-pos/events/state-events.jsonl` and the same run's verdict `journal`): the gate-hook `verify` step passes during the step phase (reason=`step-passed`), the task then writes `REVIEW_PENDING` (reason=`postflight-passed`) and finally `PASSED` (reason=`task-accepted`); the measured event sequences are the three adjacent numbers of this run (the hook step sits one place before `REVIEW_PENDING`), and in the `journal` `gate-hook:ran` precedes `acceptance:frozen`. ⑤ remediation F5: the reconciliation now matches `record.Event.Entity.StepID` against the declared gate step ids (measured `gateHookStepIds=b4-verify-step`), so a build step that passed earlier can no longer mispoint it.

Accounting conclusion (**no hedging**): the measured gate order is **contrary to** the old `docs/CONTRACTS.md` §7 wording ("freeze the candidate, then run the declared build/test/verify gates independently") - the engine's gate hook step is a member of the task's step list and therefore runs before postflight and before the candidate freeze. ⑦'s final review ruled **R-1: the contract wording needs revision**, on the ground that the contract contradicts itself in two places (L205/L227 versus L229) while the implementation agrees with L229. Under the user's 2026-09-23 ruling ⒜, **for the R-1 sub-change only** this slice changed the contract wording and left the implementation code untouched (deterministic criterion: `git diff --cached --numstat -- internal/` is empty; the slice's other artifacts are listed in §10): the three ordering statements in `CONTRACTS` §7 are unified to the measured order "in-step gate → postflight → freeze / review / promotion" (CN and `_EN` in sync), registered as **OB-22** in the directive's §12.4. Note the order is decided by the engine's step scheduling and is independent of the gate execution channel (the production channel is unexecutable here; see §7 and appendix A).

### 9.8 Native validation by ⑨ (measured on this platform)

⑨ ran on this platform (Windows, Go 1.27.0 windows/amd64) as **deterministic checks plus a same-command re-run**.

| Check | Measured |
|---|---|
| Pack integrity | `tmp/b4/SHA256SUMS.txt` recomputed in full: 347 entries, 0 missing, 0 mismatched |
| Binary identity | the three binary digests are byte-identical to report §2 (`prfrail.exe`, `agent-stub.exe`, `b4-e2e.exe`) |
| Encoding and EOL | all 14 changed files meet the convention (`.md` with BOM, the rest without; all LF): 0 violations |
| Gate stability | `--all --scope=tree` and `--scope=index` report `TOTAL_FAIL=0` on both rounds; `--selftest` reports `PASS 203/203` on both rounds |
| Test stability | `go test -count=1 ./...` has no `FAIL` package on either round |
| PowerShell syntax | `run-faces.ps1` re-parsed with `[Parser]::ParseFile`: 0 errors |
| Same-command re-run | the f5 forward leg was re-run on an independent run root with the same pinned binaries (`exitCode=0`) and compared leaf by leaf against the pack verdict: all 27 differences **fall on the declared volatile fields** (`runRoot`/`workspace` paths, `startedAt`/`endedAt`/`mapping.elapsedMs`/`journal[*].at`, `journal[].detail.processId`, the timestamp- or pid-carrying manifest and `stub.pid` digests, `eventLogHash`, `notes[1]`'s `requestRecordHash` and the `candidateHash` derived from it); `verdictClass`, `chainState`, `taskState`, `stepState`, `observations` (including `gateHookSequence=13` and `reviewPendingSequence=14`), `diff` and `checks[].pass` are **byte-identical** |

**A fail-closed observation from the first ⑨ attempt (recorded as measured)**: the first re-run did not pre-seed the run root's `workspace` the way the driver does (the driver copies `seed/b4-seed.txt` into it), so the harness refused immediately with `blocked`, reporting that the declared parent snapshot digest did not match the workspace digest - **that is the declared-input binding doing its job**, not a defect; re-running after seeding per the driver's recipe passed.

**⑨ not covered (declared explicitly)**: the driver was not re-run over the whole pack (re-running rebuilds the binaries and would invalidate the §2 digests, which are therefore still declared for this run); native Linux validation was not attempted (it belongs to C1); the production dispatch and terminal-publication legs were not re-run (they are refused by the platform constant, see appendix A).

## 10. Artifact list (G7-a)

```artifacts
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/t027-at23-e2e.md	repo
docs/validation/t027-at23-e2e_EN.md	repo
tools/b4-e2e/main.go	repo
tools/b4-e2e/fixtures.go	repo
tools/b4-e2e/ports.go	repo
tools/b4-e2e/postflight.go	repo
tools/b4-e2e/runmap.go	repo
tools/b4-e2e/verdict.go	repo
tools/b4-e2e/run-faces.ps1	repo
tools/gates/cjk-newwords.txt	repo
tmp/b4/**	local-only
```

## 11. Meta-verification statement

This slice is the **first slice judged by its own G6 and G7**: the report and the ledgers both sit inside the judging scope, so the report carries its own `artifacts` block and the ledgers carry no expiring literals. A `G6` or `G7` red on this slice's output means the output is fixed, **never the criteria and never a whitelist bypass**; a genuine exemption must be registered through the criteria's own mechanism with its reason.

## Appendix A · Platform boundaries measured by ②

| Item | Measured result |
|---|---|
| Platform publication durability | `replayStorePublicationDurability()` returns `unproven` on Windows (`internal/adapters/agent_runner_replay_store_windows.go:5`) and no public API can change that constant |
| Production dispatch | `Store.RecordRequest` fails closed with `AgentRunner replay store publication durability unproven: request publication`, spawning no process and writing no record |
| Production terminal publication | `PublishTerminal` fails closed with `AgentRunner request conflict: terminal publication requires a persisted request record`, the same root cause |
| Production gate port | `gates.ChainPort` plus `GuardExecutor` refuse every hook with `hook policy enforcement unavailable: memory limit` |
| Load determinism | two runs of the same `-mode run -sleep 2s` produce byte-identical `usage.json`, `events.jsonl`, both log streams and `diff.json`; only the harness-side manifest carries a capture timestamp and `stub.pid` carries a pid |

## Appendix B · Nature of the declared inputs (must not be read as a candidate probe)

| Item | Declaration |
|---|---|
| Nature | The capability, enforcement, availability, authorization and budget records admission needs are **declared** by the harness and constructed for a fixed workload |
| Measured constraint | the capability matrix must be fully `verified` to obtain a legal disposition, so the matrix is a declaration and not a probe result |
| No claim | **no masquerading as an AI candidate and no claim of candidate compatibility**; this slice claims the ProofRail mechanism leg only |
| Budget record | the declared availability record carries `requestsUsed=1` (the schema rejects an available record with zero); this slice makes zero model calls and zero paid requests |

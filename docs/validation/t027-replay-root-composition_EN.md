# T027 A0 replay-root composition Prerequisite Validation Report

[中文](t027-replay-root-composition.md)

Date: 2026-09-14. Conclusion: `PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`.

## Validation scope

This report covers only the T027 offline slice A0 replay-root composition and does not claim a production AgentRunner implementation. A0 delivered:

- A single production constructor `NewAgentRunnerReplayStoreForRun(runRoot, runID)`: the replay root is deterministically derived from the durable run root as `<runRoot>/agent-runner-replay` and no longer accepts an arbitrary caller root; the arbitrary-absolute-path constructor became the package-internal `newAgentRunnerReplayStoreAt`, retained only for tests and internal composition.
- Durable run-root marker: the run root must already exist and contain the regular file `events/state-events.jsonl`; cross-restart location uses only `(runRoot, runID)`, runID must be a valid evidence ID, and scanning or guessing is forbidden.
- runID binding: loading an R or C whose `runId` differs fails closed as corruption; writing an R or C with a different runId is rejected as a root conflict on both the request and completion sides.
- Overlap and alias rejection: `RejectWriterRoots` keeps the replay root disjoint in both directions from writer roots (including 8.3 short names and file-identity aliases); `RejectProtectedRoots` rejects equality or a protected root inside the replay root, allows a protected root containing the replay root only when the containment chain covers the run root, and fails closed whenever containment or identity cannot be decided.
- symlink/reparse: component-level rejection of symlinks (Unix) and reparse points (Windows, including junctions and mounted folders) for the root, `requests/`, `completions/`, and the ownership path at construction and again after the lock in every read/write entry point (`State`, `RecordRequest`, `RecordCompletion`); failures fail closed and must not land outside the root.
- ownership: `ownership.json` (kind, schemaVersion, runId, runRootPath, runRootHash, createdAt; domain `proofrail:agent-runner-replay-ownership:1\n`) is published monotonically without replacement and is not gated on publication durability `proven`; a missing marker with published R/C is corruption; a missing marker over an empty root is rebuilt; a runId or hash mismatch, or a moved or copied directory, is a root conflict; if the parent sync fails after the marker link, the call fails closed and leaves the now-visible marker in place; a failed construction (including bootstrap and the low-level constructor) leaves no ownership metadata behind.
- Boundaries: no real dispatch was introduced; no arbitrary caller root is reopened; `unproven` publication semantics are unchanged; T027 remains `BLOCKED / NOT IMPLEMENTED`.

## Execution results

### Native Windows execution

| Command | Result |
|---|---|
| `gofmt -l internal/adapters` | Clean (no output) |
| `go build ./...` | Pass |
| `go vet ./...` | Pass |
| `go test -count=1 ./...` | All pass |
| `go test -count=1 -run 'TestAgentRunnerReplay' ./internal/adapters` | Pass |
| `node tools/contracts/contracts.test.js` | 4/4 pass (no schema/fixture changes; regression confirmation) |

### Windows platform counterexample tests (executed, not skipped)

| Test | Result |
|---|---|
| Junction replacing the replay root → unsafe path rejection | PASS |
| Reopening the same root through a case-alias spelling | PASS |
| 8.3 short-name alias overlap rejection and short-name reopen | PASS |
| Protected short-name alias → conflict | PASS |
| `\\?\` extended-length prefix reopening the same root | PASS |
| Drive-relative `C:relative` → conflict | PASS |
| Concurrent distinct spellings all succeed | PASS |
| Runtime `requests` swapped to a junction → `RecordRequest` unsafe-path with the external target left empty | PASS |

### Linux native regression (GitHub Actions Ubuntu)

After commit `7ad1ac0` was pushed, CI run [34799747443](https://github.com/larsonzh/prfrail/actions/runs/34799747443) succeeded:

| Job | Steps | Result |
|---|---|---|
| Go ubuntu-latest | Build / Vet / Test / Contract fixtures | all success |
| Go windows-latest | all steps | success |

The Ubuntu job actually exercised the A0 Unix-side tests (symlink rejection, sync ordering, sync-failure fail-closed, no ownership on failed construction, convergence exhaustion, runtime link swap), closing the earlier "awaiting Linux native confirmation" boundary; Windows-only tests (junction, 8.3 short names, `\\?\` prefix) run only in the Windows job by build tag.

### Review conclusions

- V4 Pro implementation pre-review: PASS (conditional) → fixed the Medium (caller-root chain reparse check) and every Low (completion-side write classification, contract wording, marker regular-file check, createdAt validation, marker IO error sentinel wrapping, zero-value root guard on the read path) plus 12 added counterexample tests.
- Codex independent final review: changes required → Medium (runtime TOCTOU re-check) implemented with counterexample tests; Low ① (ownership publication ordering) tightened so the low-level constructor completes before ownership is published; Low ② (counterexample gaps) closed with tests.
- Codex closing re-review: PASS.

## Known boundaries

- Unix-side tests (symlink rejection, sync ordering, sync-failure fail-closed, no ownership on failed construction, convergence exhaustion, runtime link swap) passed the GitHub Actions Ubuntu native regression (run 34799747443, commit `7ad1ac0`); broader Unix crash injection or non-Ubuntu distribution differences remain outside this slice.
- A residual TOCTOU window between user-mode component checks and file operations cannot be fully eliminated; it is recorded in CONTRACTS as an explicitly accepted boundary.
- Full-chain reparse/symlink rejection means run directories inside OneDrive-redirected folders, junction trees, or mounted-folder layouts are rejected fail-closed; this is an intentional safety trade-off and any relaxation requires a separate ADR.
- The 8.3 short-name and `\\?\` cases skip explicitly when the environment lacks the capability; they executed on this machine without skipping.

## Explicitly not performed

- No real dispatch was introduced; `AgentRunnerPort`/Engine wiring is untouched;
- No real model CLI was invoked, no network access, no SecretStore reads;
- No commit, push, or publish happened while this report was written; both followed under same-turn user authorization (`7ad1ac0`), with CI evidence in the section above;
- T027 remains `BLOCKED / NOT IMPLEMENTED` and AT-23 has not passed.

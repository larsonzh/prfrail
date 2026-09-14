# T027 A2 dispatch and launch ambiguity prerequisite validation report

[中文](t027-dispatch-launch-ambiguity.md)

Date: 2026-09-14. Verdict: `PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`.

## Validation scope

This report covers only the T027 offline slice A2 dispatch and launch ambiguity and makes no claim about the AgentRunner production implementation. A2 delivered:

- Contract freeze (before implementation): `docs/CONTRACTS.md` / `CONTRACTS_EN.md` gained four paragraphs — ① dispatch publishes R first, then a narrow reconfirmation after publication and before launch; failures are fail-closed (R retained, no receipt, no spawn, no auto-retry, errors identifiable via `errors.Is` with the underlying cause preserved); ② the two-file store-local receipt protocol (pre-spawn `launches/intent.<requestId>.jsonl` carrying launchId, requestId+requestHash, run/task/step/attempt, adapterId, reconfirmation time, authorization and budget digests, with no-replace plus bounded reread convergence and pre-write rejection on unproven platforms; post-spawn `launches/identity.<launchId>.jsonl` binding the intent digest and the opaque process identity; invariants: no intent means never spawned, a missing identity means unproven, and neither permits a blind relaunch); ③ dispatch replay handling (R-only never launches; an existing intent is never republished or relaunched; intent+identity never relaunches; an R+C terminal state never launches; the only permitted in-process retry is a port-contractual unspawned failure, and an intent is never republished; takeover and settlement remain with later slices); ④ `launches/` shares the path-safety, ownership, and sync discipline of requests/completions.
- `AgentRunnerReplayDispatcher` (adapters, implementing `chain.AgentRunnerPort`): binding validation → `RecordRequest` (terminal/unknown blocked with distinct sentinels) → authorization ledger snapshot plus the `AgentRunnerLaunchReconfirmer` narrow reconfirmation → deterministic launchId → `RecordLaunchReceipt` (single winner) → `Launcher.StartAgentRunnerProcess` → `RecordLaunchIdentity`; success returns a three-part evidence chain (R, receipt, identity recordHash).
- `AgentRunnerLaunchReconfirmer` (value semantics): reuses the U1-extracted package functions `validateAgentRunnerAuthorization` / `validateAgentRunnerBudget` for the narrow recheck of "authorization active, unrevoked, unexpired + reservation outstanding", preserving underlying causes in the error chain.
- Store extensions: `RecordLaunchReceipt` / `LaunchReceipt` / `RecordLaunchIdentity` / `LaunchIdentity`; write paths follow the R/C discipline (read-only replay resolution first → durability gate → binding checks → no-replace plus bounded reread); receipts and identities are store-local records (separate domain digests, never on the wire, never part of R/C recordHash); `launches/` is now part of the production constructor, `bootstrapReplayStoreRoot`, the locked runtime path-safety recheck, and the ownership-publish preflight enumeration.
- U1 extraction: the A1 private validators were promoted to package functions and the admission call sites rewritten, with A1 test guards keeping behavior unchanged.
- Focused tests: 24-way cross-store single winner (exactly 1 first-launch / 23 already-launched), both no-replace collision branches (same payload → replay, foreign payload → conflict), identity without intent = corruption, unproven-platform write rejection with exact replay allowed, runId and full-field binding (write conflict + read corruption), zero-value store fail-closed; the full dispatcher branch set (R-only / intent-only / intent+identity / terminal blocks; revocation including "R retained, no receipt, no spawn, no auto-retry"; start failure including "intent retained, never relaunched"; launchId mismatch and identity conflict → `ErrAgentRunnerLaunchIdentityUnproven` with no kill and no retry; 8 concurrent dispatchers with exactly one spawn; binding mismatches writing no records); Unix specifics (construction-time symlinked `launches/` rejection, runtime swap rejection with the external directory proven empty, bounded reread convergence on a dangling slot, sync failure with no rollback and a reopen resolving as replay).

## Execution results

### Windows native run (native evidence)

| Command | Result |
|---|---|
| `gofmt -l internal` | clean (no output) |
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `go test -count=1 ./...` | all pass |
| `go test -count=1 -run 'AgentRunnerLaunch\|AgentRunnerReplayDispatcher\|AgentRunnerReplayStoreLaunch\|AgentRunnerReplayRoot' ./internal/adapters` | pass |

Notes: the A2 concurrency counterexamples (8 concurrent dispatchers, 8 concurrent store calls, 24 cross-store calls) ran natively on Windows; Windows has no gcc so `-race` was not run locally.

### GitHub Actions CI (post-push Linux native regression)

- Committed as `aeb504e` and pushed to `origin/main`; run [34844216048](https://github.com/larsonzh/prfrail/actions/runs/34844216048) concluded `success`: `Go ubuntu-latest` (Build/Vet/Test/**Race**/Contract fixtures) and `Go windows-latest` (Build/Vet/Test; Race skipped by condition) both passed; `Bootstrap and release evidence`, `Candidate build`, and `Candidate probe` were skipped by the push condition (same as A0/A1).
- The Unix-specific tests (launches symlink swap, dangling-slot convergence, sync failure without rollback) and the `-race` concurrency regression were actually executed by that run on native Linux.

### Review outcomes

- V4 Pro pre-review: conditional pass (no High). M1 (no real counterexample for the cross-store launch-receipt single winner and the collision classification) and M2 (no test for the identity-without-intent corruption branch) were remediated with new tests; L1–L5 and L7 were remediated (full explicit field cross-checks, launches sync-failure counterexample, launchId mismatch counterexample, Windows launches junction swap counterexample, missing-marker-with-non-empty-launches corruption counterexample, `launches/` added to `bootstrapReplayStoreRoot`); L6 (unproven ordering semantics) and L8 (zero `StartedAt` guaranteed by the A6 port implementation) are recorded as known boundaries.
- Codex independent final review (including a per-test counterexample audit): conditional pass (no High/Medium); the single Low tail item (explicit binding to also compare RunID/AuthorizationHash/BudgetHash plus a persisted-runId check in the replay branch) was remediated; the final review confirmed that no double-spawn path exists.

## Known boundaries

- W1/W2 do not auto-take-over: an already-owned launch slot (including the concurrent first-dispatch window) and a failed revocation reconfirmation both return identifiable blocks only, with no automatic takeover or compensation; takeover and settlement belong to A3.
- The recheck→spawn residual window (decided in U4): an irreducible window between reconfirmation and the real `Start` remains; this slice only guarantees the state within that window is auditable and that failures are fail-closed.
- Unproven-platform ordering semantics: receipts and identities allow "write-free exact replay" and reject every new write (the same rule as `RecordCompletion`); `RecordRequest` rejects unproven first because it carries first-dispatch eligibility.
- Zero `StartedAt`: the port contract states that a result is returned only when the process was spawned and the identity is populated, so implementations (A6) must guarantee it; this slice adds no extra rejection.
- Auto-retry stance: the contract only *permits* an in-process retry for a port-contractual unspawned failure; this slice chooses the more conservative never-auto-retry, leaving post-failure handling to humans or upper layers.

## Explicitly not executed

- No completion published, no Engine wiring, no real candidate run, no takeover/settlement (those belong to A3 and later);
- no real model CLI invoked, no network access, no SecretStore reads;
- at the time of writing this report no commit, push, or publish was performed; the commit and push (`aeb504e`) were later completed under same-turn user authorization, with the CI evidence above;
- T027 remains `BLOCKED / NOT IMPLEMENTED` and AT-23 has not passed.

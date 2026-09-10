# T025 Operator Interaction Validation Report

[中文](t025-operator-interaction.md)

Date: 2026-09-10. Verdict: `T025 COMPLETE / AT-22 PASS`.

## Implementation Scope

1. `schemas/operator-interaction.schema.json` and `testdata/contracts/{valid,invalid}/operator-interaction.json`: freeze request/response wire format, actors, binding fields, and unknown-field rejection; included in the independent gate covering 30 Schemas and 86 cases carried by 12 fixture files.
2. `internal/chain/operator_interaction.go`: RFC 8785 canonical digests, strict decoding, append-only ledger, pending replay, duplicate/tamper/binding rejection, and a same-conversation/new-request resume command.
3. Engine single-writer control API: `OpenOperatorInteraction` persists chain/task/step pause state; `ResumeOperatorInteraction` revalidates request/response, run/task/step/attempt/workspace/conversation/context, current state, and latest event evidence before restoring step/task state. The chain stays paused until an explicit run entry continues it.
4. `internal/console/interactions.go`: `interactions list/respond/tui`, with JSON automation output and a focused ANSI-free terminal inbox.

## AT-22 Coverage

1. Only requests passing strict structure, actor, digest, and evidence validation can open an interaction. Free text cannot trigger a state transition.
2. A response must match the operator, attempt, workspace, conversation, context, and request hash, and its selection must be allowed. Stale, duplicate, tampered, and foreign responses fail closed.
3. The ledger is append-only JSONL. Restart reconstructs pending items from request/response records without trusting chat history, TUI memory, or SessionBridge `@sbr-review`.
4. The TUI presents the question, reason, risk, required action, and allowed responses. Invalid input or persistence failure does not remove the pending item.
5. Engine state events use request/response record hashes as evidence. Open pauses the chain before writing task/step waiting state; resume restores the step before the task and returns no executable resume command until both writes complete.
6. Without a cross-event EventStore transaction, a partial write keeps the chain paused. Retry with the same evidence hashes converges; another response cannot take over the partial state.
7. General clarification returns only a structured selection. Authorization, review, and manual writes still use their respective authorization, review, or handoff paths and cannot be bypassed through operator interaction.

## Gate Results

1. `go build ./...`: passed.
2. `go vet ./...`: passed.
3. `go test -count=1 ./...`: passed, with 13 tested packages and two command packages without tests.
4. `node tools/contracts/contracts.test.js`: passed, 2/2; all 30 Schemas, 86 cases carried by 12 fixture files, and two canonical vectors pass.
5. `git diff --check`: passed.
6. The repository has no standalone encoding checker. This change keeps `.go/.json` as UTF-8 without BOM + LF and `.md` as UTF-8 with BOM + LF, confirmed by read-only byte checks.

## Focused Tests

- Nine chain cases: single bound response, broken binding, transitions, strict decode/tamper, current-binding revalidation, Engine persistence/replay, write failure, partial-write convergence, and foreign-run/stale-state rejection.
- Seven console cases: restart inbox rebuild, response persistence, stale/unstructured input, missing-ledger no-op, valid TUI selection, invalid TUI input, and persistence failure preserving pending state.
- One contract group: positive/negative operator-interaction fixtures plus the full Schema/canonical validator.

## Boundaries

1. `interactions tui` is a focused ANSI-free terminal interaction, not the complete unified TUI. The Bubble Tea prototype is not in the release binary.
2. T025 returns a validated resume command and persists chain state, but does not start or resume a real CLI Agent. T026/T027 own AgentRunner wire formats, capability probes, and session execution.
3. Timeout/disconnection has no default response. Without an appended valid response, restart leaves the request pending.
4. No real model/network call, `git commit`, `git push`, signing, or release action was performed.

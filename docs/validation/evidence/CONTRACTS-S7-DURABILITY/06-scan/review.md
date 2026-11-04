# ⑥ 独立扫描原文（P54 CONTRACTS-S7-DURABILITY）

载体：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（只读，2026-11-04；A–E 形状）。以下为返回文本逐字归档。

## A. Scope
- Read-only independent scan of the §7 durability contract, the C6 prebuild activation path, the Windows slot implementation, and the no-op fallback path.
- Scope excludes code edits or any dependency on the earlier ⑤ pre-review conclusion.

## B. Findings
- No findings.
- The contract wording and the implementation remain aligned: the C6 path is build-gated, store-activated only after a successful prebuild, and the non-activated build remains a no-op.

## C. Evidence
- The production assembly order is explicit in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L129-L195): `NewAgentRunnerReplayStoreForRun` is followed immediately by `PrebuildReplaySlots(...)`, before any dispatch.
- The activation semantics are defined in [internal/adapters/agent_runner_replay_store_c6_enable.go](internal/adapters/agent_runner_replay_store_c6_enable.go#L13-L115): prebuild success sets the store active, failure leaves the run fail-closed, and reopen verification is idempotent.
- The Windows slot behavior is enforced in [internal/adapters/agent_runner_replay_store_c6_windows.go](internal/adapters/agent_runner_replay_store_c6_windows.go#L15-L239): header format, payloadLen validation order, empty-slot freeze semantics, and no-replace publication rules.
- The unactivated fallback is pinned in [internal/adapters/agent_runner_replay_store_c6_noop_test.go](internal/adapters/agent_runner_replay_store_c6_noop_test.go#L12-L87): a non-`c6slot` build cannot activate the slot path or touch the replay root.
- Fresh validation command: `cd /d d:\LZProjects\prfrail && go test ./... && go test -tags c6slot ./internal/adapters`
  - Result: all packages reported `ok`; the adapters package completed successfully with no failing assertions.

## D. Falsifiability audit
- The Windows C6 behavior is explicitly testable and tested in [internal/adapters/agent_runner_replay_store_c6_windows_test.go](internal/adapters/agent_runner_replay_store_c6_windows_test.go): empty-header creation, committed-slot rewrite refusal, corruption ordering, uncommitted intermediate state handling, prebuild activation, and fail-closed prebuild conditions.
- Those tests are structured as negative/positive checks and therefore provide a valid falsification surface for the contract.

## E. Unverified / overturn conditions
- None open for the current implementation.
- Any future platform or premise change would require re-running the same validation set and re-checking the C6 activation enablement and durability claims under that new premise.

INDEPENDENT SCAN: PASS

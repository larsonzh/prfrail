# T026 AgentRunner Contract and Capability Validation Report

[中文](t026-agent-runner-contract.md)

Date: 2026-09-10. Verdict: `CONTRACT COMPLETE / RUNTIME PROBE BLOCKED BY AUTHENTICATION / AT-23 NOT PASS`.

## Completed Scope

1. Four Schemas freeze the request, capability, event, and completion wire formats. The contract gate now covers 34 Schemas, 96 cases carried by 14 fixture files, and two canonical vectors.
2. `internal/adapters/agent_runner*.go` implements RFC 8785 domain-separated hashes, strict decoding, sorted sets, create/resume constraints, request/event/completion binding, and idempotency/conflict indexes reconstructable from durable records.
3. Capability uses a fixed 12-entry matrix and is compatible only when every entry has runtime evidence and is verified. Events lock eventId, session/sequence, and request/session so one request cannot drift across sessions. Completions lock both completionId and requestId; completed proves only external execution end, process-tree stop, complete logs/usage, and captured output manifest, never task PASS. The pinned capability test also reads the referenced evidence file and verifies its byte SHA-256 so the record cannot cite missing or stale evidence.

## Candidate Probe

- Candidate: GitHub Copilot CLI 1.0.83 on Windows x64, build commit `e6a98f1`. The native `copilot.exe` is 144,796,448 bytes with SHA-256 `d3f3bb7b8bbf68357ad29f514a179d09f76135483d8bfb643131b8600f671ee2`; the platform package metadata digest is `d6d07e0943462974bba6380170ca6c6d147e6e609b2cc5ceff165f8f16faa5ed`.
- Direct native execution bypassing both the VS Code bootstrap and npm shim reports version 1.0.83. `--help` was used only to design probes and is not accepted as capability-pass evidence.
- The operator explicitly authorized one minimal, potentially billable model call. The first launch script failed before Copilot started because Windows PowerShell 5.1 does not support `Get-Date -AsUTC`, so it did not consume the authorization. The corrected call used an isolated workspace, JSONL streaming, a 30 AI-credit soft limit, only the `shell` tool, exact `shell(pwd)` approval, HTTP/HTTPS denial, and disabled built-in MCP, system-temp access, custom instructions, remote control/export, and automatic updates.
- The corrected CLI process produced no JSONL, stderr, or `usage.json`. A 468-byte dedicated process log (SHA-256 `7d1b6cd0a3f696a9233b51071fff75c94542e33b2a0c1c697ef2008841ac5129`) shows that authentication orchestration rejected the current GitHub credential because classic PATs are unsupported; the log contains no credential value. The process was interrupted after that error, and no matching probe process remained. No model response or usage evidence exists, and the artifacts cannot prove model-endpoint contact or billing, so the one-call authorization was not used for a retry.
- After fresh authorization, a second preflight removed `GH_TOKEN` only from the Copilot child process, without changing the user environment or configuration, to try falling back to stored OAuth credentials. The CLI explicitly reported that no authentication information was found and exited 1. Its 381-byte `usage.json` (SHA-256 `f9f90cadcbcb1d9ba7176d520df7585abc79ff3ef751937204c366c6765b0c98`) records zero user requests, zero input/output tokens, zero API duration, zero premium request cost, and no model metrics. The second authorization therefore also caused no model call or billing.
- The second preflight explicitly used `--model auto`. Copilot CLI does not inherit the model selected in the VS Code chat input. Use `--model <name>` to select a CLI-supported model for one run, or `--model auto` to delegate routing to the CLI. Without the flag, the CLI uses its own `model` configuration or default.
- All 12 capabilities remain unknown: noninteractive use, cwd, events/complete logs, session create/resume, cancellation, process-tree stop, tool/network/permission controls, usage, and unattended confirmations. Disposition remains blocked. Machine evidence and the hashed blocked capability record are under `testdata/agent-runner/capability-probes/`.

## Gate Results

1. `go build ./...`: passed.
2. `go vet ./...`: passed.
3. `go test -count=1 ./...`: passed, with 13 tested packages and two command packages without tests.
4. `node tools/contracts/contracts.test.js`: passed, 2/2; 34 Schemas, 96 cases, and two canonical vectors passed.
5. `git diff --check`: passed.

## Unblock Conditions

A retry first requires completing a CLI 1.0.83-supported OAuth login in the user's own terminal or configuring a fine-grained PAT with `Copilot Requests` permission, followed by fresh explicit authorization for one minimal, potentially billable model call. The next run should first expose foreground output proving that the CLI entered non-interactive mode, then execute the 12 probes in an isolated workspace and update the capability record. T026 can close and T027 can start only after every required entry is verified; AT-23 must not be reported as passed now.
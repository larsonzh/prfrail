# ProofRail User, Operator and Trusted Release Guide

[简体中文](OPERATIONS.md)

Date: 2026-09-09; S1 operational design. No formal ProofRail release package exists. The S1 model of manually extracting a Windows amd64 portable ZIP, invoking it by explicit path, and not modifying PATH is approved; see the candidate [installation guide](INSTALLATION_EN.md). T016/T019/T020/T023 provide the CLI baseline and delivery foundation: `version/init/validate/config explain/preview/cost report/run/report` plus internal export modules (`internal/snapshot/export.go`, `internal/evidence/delivery.go`). `preview` is an offline read-only static report and does not execute commands, network probes, model calls, or credential reads. `run` currently supports noop-only chains and fail-closes with a non-zero exit for executable `code/build/verify` steps. Remaining product capabilities are still planned RFC interfaces. See the [installation plan](INSTALLATION_PLAN_EN.md) for unresolved release binding and [business workflows](BUSINESS_WORKFLOWS_EN.md) for the complete user flow.

## 1. Available Today

Development needs Go 1.22+. Future release users will not need Go, but their project's harness toolchain remains required. From repository root:

```powershell
go version
go build ./...
go vet ./...
go test ./...
go run ./cmd/prfrail version
go run ./cmd/prfrail init --workspace .
go run ./cmd/prfrail validate --chain .\proofrail.chain.json
go run ./cmd/prfrail config explain --chain .\proofrail.chain.json
go run ./cmd/prfrail preview --chain .\proofrail.chain.json
go run ./cmd/prfrail cost report --ledger .\cost-ledger.json
go run ./cmd/prfrail run --chain .\proofrail.chain.json --run-id run-demo
go run ./cmd/prfrail report --run-dir .\tmp\prfrail-runs\run-demo
```

Core packages now have automated tests, but full product E2E/TUI is still in later slices. The current visual surface is a line-oriented CLI in Windows Terminal, PowerShell, or the VS Code integrated terminal, ANSI-free by default and with `--json`; there is no separate graphical window. Bubble Tea v1.3.4 was only an isolated prototype and is not part of the release binary. gopls is a development tool, not a core runtime dependency. The released ProofRail binary does not directly access GitHub, Go/npm registries, or other public URLs; release-evidence verification reads local files only. Dependency downloads and GitHub Actions networking belong only to development/CI. System/Git proxies do not automatically configure Go downloads; use process HTTP_PROXY/HTTPS_PROXY when needed and restore them afterward, without arbitrary global GOPROXY changes or disabled verification.

The target TUI will show pending requests and allowed responses after AI returns structured `operator-action-required`; the task remains `WAITING_FOR_OPERATOR` until the chain control API validates and persists the response. Before CLI Agent resume, revalidate attempt/workspace/session/process and authority; uncertainty creates a new attempt. This loop belongs to T025 and is not delivered. SessionBridge `@sbr-review` remains standalone diagnostics only.

## 2. Installation and First Run

S1 uses a manually extracted Windows amd64 portable ZIP: verify the fixed commit/GitHub CI plus SHA256SUMS, SBOM, and license manifest; extract into a user-selected independent version directory; then invoke `prfrail.exe` by full path for version and static preview. Never modify PATH, the registry, or system directories, and never overwrite a running host. See the [installation guide](INSTALLATION_EN.md) for copyable commands and upgrade/rollback/uninstall steps. SHA256SUMS proves agreement with the manifest, not publisher identity. S1 Linux is core CI/preview, not production support. Do not execute wrong-platform binaries.

Current copyable flow is `prfrail init --workspace <path>`, `prfrail validate --chain <chain-file>`, `prfrail config explain --chain <chain-file>`, `prfrail preview --chain <chain-file> [--json]`, `prfrail cost report [--ledger <path>] [--json]`, `prfrail run --chain <chain-file> [--run-id] [--run-dir]`, and `prfrail report --run-dir <run-dir>`. New users should first complete the side-effect-free [five-step quick start](INSTALLATION_EN.md#3-five-step-quick-start). `preview` reports `previewHash`, permission boundaries, unknowns, and zero call counters in read-only mode; `cost report` reads the local cost ledger and reports reserved/settled/unknown holds without default telemetry export; `run` is currently noop-only. Export support is implemented at library level, but a dedicated CLI `export` command remains in later slices; full executable-step gate/adapter integration is also still in later slices.

Core offline mode does not require VS Code. Target formal AI execution uses a pinned CLI Agent through AgentRunner. Only supervised black-box candidates require VS Code, usable Copilot Chat and SessionBridge 0.1.1 visible mode. Before launch, acknowledge unobservable process/network/cost/external effects, supervise the Agent, explicitly return the workspace, then let ProofRail validate independently. Neither path is delivered; an installed extension or reachable transport is not AI-coding acceptance.

The current CLI baseline uses a single chain config file (`--chain` first, then `./proofrail.chain.json`, then `./proofrail.json`). The `proofrail.toml/workspace.toml` layered model remains planned. Runtime state belongs to the engine. Never manually edit state/journals/receipts/accepted manifests. config explain exposes origins, tools/network permissions and budgets; model prose is not configuration.

Before running, verify distinct source/run/store roots, secret exclusions, disk, tools, model/cost limits, review actors and enforceable isolation. Uncommitted source files do not require reset. External source changes are not automatically absorbed after baseline capture. Never give agents the original workspace to edit.

## 3. Operator Decisions

| Observation | Allowed action | Prohibition |
|---|---|---|
| REVIEW_PENDING | Verify parent/candidate/evidence hashes and code/docs; approve/reject/authorized waive | Technical success as acceptance, receipt editing |
| Blocking failure | Preserve evidence, prove stop, bounded repair/cancel | Skipping failure into next task |
| Uncertain process liveness | Read-only identity/termination/lease inspection, stay paused | Arbitrary PID killing, early restore |
| Low disk | Pause, GC dry-run for expired unreferenced data or increase quota | Deleting baseline/accepted/failure evidence |
| Adapter timeout | Check durable dispatch/requestId and uncertainty | New-ID infinite retries, visible/auto fallback |
| User selects black-box candidate | Read and acknowledge bound risk, supervise the visible Agent, explicitly return only after work ends | Treat as AgentRunner, expose source/store, or show unknown cost/process/external effects as safe |
| AI requests operator input | Use the ProofRail TUI to inspect the structured request and submit an allowed response; open handoff for file edits | Authorize through free-form chat, edit state files, or depend on `@sbr-review` to advance |
| Budget exhaustion | Human fingerprint/cost/new-authorization review | Resetting counters to bypass hard_block |
| Corrupt store | Quarantine and verify references offline | Rehashing corrupt content as valid history |

## 4. Handoff and Recovery

Planned actions are prfrail handoff open/complete/abort/request-agent. open requires stopped managed writers, flushed journals and a frozen manifest before the exclusive operator lease. Read goals, scope, reproduction commands, budget, expiry and return conditions. Edit only the specified run-workspace.

complete returns control, not approval: rescan, check boundaries/secrets/unknown processes, rerun all prescribed gates and obtain independent review. abort preserves evidence and discards this candidate. request-agent creates a new attempt from confirmed conclusions. Disconnection/expiry pauses without reassignment. Enter passwords/tokens/MFA directly into a safe terminal, never chat/issues.

Recovery: preserve scene, prove managed-process state/old-writer termination, verify store/events, replay journal, rebuild state, verify complete acceptance evidence, then explicitly continue. A PASSED state file is insufficient. pause/cancel use the control API, not manual file edits.

## 5. Minimal Diagnostic Pack

Include version/OS, run/task/step/attempt, first error category, requestId when relevant, input/evidence hashes, minimal reproduction, commands/exits, expected/actual and redacted summaries. Exclude tokens, keys, full conversations and entire repositories.

| Symptom | First evidence/owner | Smallest check |
|---|---|---|
| Config rejection | taskdef/value origin | AT-01 fixture |
| Capture/restore failure | snapshot/path manifest | AT-02 native paths/permissions |
| Uncertain application | applier/journal | AT-04 single fault-point replay |
| Stuck command | guard/gates/process identity | AT-06, not an immediate timeout increase |
| No response | adapters/durable request/instance | AT-08 offline before host |
| No progress after review | chain/candidate binding/publication receipt | AT-05 expiry/hash changes |
| Code passes, docs fail | impact rule/generation/doc hook | AT-11, do not delete the rule |

## 6. whois Shadow Migration

whois v3.4.0 is frozen. This round/S0 does not change its production code, tasks or release directories. After authorization, export fixed-hash read-only redacted fixtures; never restore source while A/B processes run.

| whois concept | ProofRail mapping | Evidence |
|---|---|---|
| A/B | Ordered T1/T2 tasks | Input/output snapshot hashes |
| D/V rounds | Labels plus code/build/verify/noop | Explicit kinds, reasons, results |
| start-file | Chain/environment/machine-state separation | Mapping differences, no interface-compatibility claim |
| Step47/C goldens | C harness hooks | Fixed input, exit categories, artifact hashes |
| Repair candidate | Prepare/Inspect/Validate/Promote | Old/new failure classification and promotion evidence |

S0 uses a small non-C Go fixture and dual/triple-language configurations to challenge core coupling. S1 replays only isolated copies, comparing inputs, decisions and failure classes; explain/test differences rather than demanding identical nondeterministic logs. Discuss cutover only after S1 exit, rollback drills and explicit approval. Never automatically replace frozen workflows.

## 7. Self-Hosting, Upgrade and Release

Bootstrap the first seed from a fixed commit through normal Go build, independent tests and human audit. Record commit/source hash, toolchain, oracle and approval actor. The current stub is not a trusted seed.

Stable N captures isolated input/builds N+1. The candidate never overwrites host. Replay a fixed chain in a new process and separate run/store, with an external oracle comparing schema/receipt/recovery/fail-close goldens and a clean-room drill. Failure blocks the candidate without changing seed or historical evidence.

Release checklist: aligned RFC/schema/version; fixed commit and GitHub CI; build/vet/test, native platforms and required AT; license/vulnerability inventory and SBOM; correct binaries/CGO_ENABLED=0; SHA256SUMS; upgrade/rollback tests; bilingual notes/limitations; independent gate and same-turn commit/push/release authorization. Checksums must not be presented as publisher authentication; signing/attestation is a later enhancement. Default remote is origin, never automatic Gitee mirror pushes.

Rollback stops candidate processes and preserves evidence; restore an older binary only when storage compatibility is explicit. Unknown formats stay read-only. Release-trust/support/response policies are in [ADR_REGISTER_EN.md](ADR_REGISTER_EN.md).

## 8. Complete Product Flow (Planned, RFC Section 19)

| Stage | Action/success evidence | On failure |
|---|---|---|
| PC-01 Trial | Preview supplied no-AI fixtures, inspect permissions/budgets/unknowns, separately authorize executable probes before running | No config hook execution or automatic tool installation |
| PC-02 Delivery | Select accepted results, choose a new nonoverlapping destination, verify hashes/exclusions/review references | Reject overwrite; incomplete output is not delivery; export is not external release |
| PC-03 Approval/revocation | Inspect inbox/candidate hash; verify blocked new dispatch and in-flight stop evidence after revocation | Remain paused if stopping fails; no hard-gate waiver |
| PC-04 Recovery/diagnostics | Read-only redacted plan identifies last acceptance, journal, processes and unknown effects | No manual lock deletion/blind redispatch; file restore cannot undo uploads/deployments |
| PC-05 Costs | Run `prfrail cost report --ledger <path>` to inspect reserved/settled/unknown holds and measurement sources; authorize increases explicitly | Timeout is not free; subscription mode guarantees authorized call/token caps, not monetary bills |
| PC-06 Retirement/uninstall | Stop managed writers/revoke authorization, verify complete backup and new-store restore drill; retain run/store by default | Preserve unknown formats; do not delete shared SessionBridge/toolchains; confirm data deletion separately |

Users decide whether to merge/release S1 delivery packages externally; ProofRail does not write source or execute Git automatically. Check support status, SBOM/licenses, signing identity and revocation freshness separately. Unknown offline revocation must not display current trust. Backups exclude credentials; restore environments require separate secure configuration.

For historical secrets, quarantine/block export/revoke credentials, then obtain a human retention/deletion decision; never rewrite old hashes as intact history. Deletion does not guarantee SSD/external-copy erasure. Reports stay local without default telemetry; support submissions contain only minimal redacted reproduction data. T020 now provides library-level export capability; T021 adds authorization/revocation with an approvals inbox, T022 adds side-effect classification/observation and recovery diagnosis (library level), T023 adds the cost ledger with local cost reporting (`cost report`), and T024 adds library-level complete backup, new-store verification, lifecycle disposition, and release inventory. T017 artifacts and installation commands remain pending, so the full product flow is still not an executable command tutorial today.
# ProofRail User, Operator and Trusted Release Guide

[简体中文](OPERATIONS.md)

Date: 2026-09-08; S1 operational design. No production ProofRail package exists, but T016/T019/T020 now provide the CLI baseline and delivery foundation: `version/init/validate/config explain/preview/run/report` plus internal export modules (`internal/snapshot/export.go`, `internal/evidence/delivery.go`). `preview` is an offline read-only static report and does not execute commands, network probes, model calls, or credential reads. `run` currently supports noop-only chains and fail-closes with a non-zero exit for executable `code/build/verify` steps. Except for the next section, remaining product capabilities are still planned RFC interfaces. See the [installation plan](INSTALLATION_PLAN_EN.md) for unresolved distribution choices and [business workflows](BUSINESS_WORKFLOWS_EN.md) for the complete user flow.

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
go run ./cmd/prfrail run --chain .\proofrail.chain.json --run-id run-demo
go run ./cmd/prfrail report --run-dir .\tmp\prfrail-runs\run-demo
```

Core packages now have automated tests, but full product E2E/TUI is still in later slices. Current CLI output is ANSI-free by default and supports `--json`. gopls is a development tool, not a core runtime dependency. System/Git proxies do not automatically configure Go downloads; use process HTTP_PROXY/HTTPS_PROXY when needed and restore them afterward, without arbitrary global GOPROXY changes or disabled verification.

## 2. Installation and First Run Design

Future installation: select the Windows amd64 package, verify signature/checksum, extract separately, check version, preview statically, then authorize capability probes separately. S1 Linux is core CI/preview, not production support. Do not execute unknown signatures or wrong-platform binaries. Back up the complete object/event/state-reference closure and verify historical runs read-only before upgrades; never overwrite a running host.

Current copyable flow is `prfrail init --workspace <path>`, `prfrail validate --chain <chain-file>`, `prfrail config explain --chain <chain-file>`, `prfrail preview --chain <chain-file> [--json]`, `prfrail run --chain <chain-file> [--run-id] [--run-dir]`, and `prfrail report --run-dir <run-dir>`. `preview` reports `previewHash`, permission boundaries, unknowns, and zero call counters in read-only mode; `run` is currently noop-only. Export support is implemented at library level, but a dedicated CLI `export` command remains in later slices; full executable-step gate/adapter integration is also still in later slices. Start with a file-queue fixture consumer without AI; VS Code is optional.

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

Bootstrap the first seed through normal Go build, independent tests and human audit. Record source hash, toolchain, oracle and signing/approval actors. The current stub is not a trusted seed.

Stable N captures isolated input/builds N+1. The candidate never overwrites host. Replay a fixed chain in a new process and separate run/store, with an external oracle comparing schema/receipt/recovery/fail-close goldens and a clean-room drill. Failure blocks the candidate without changing seed or historical evidence.

Release checklist: aligned RFC/schema/version; build/vet/test, native platforms and required AT; license/vulnerability inventory; correct binaries/CGO_ENABLED=0; checksums plus approved signature; upgrade/rollback tests; bilingual notes/limitations; independent gate and same-turn commit/push/release authorization. Default remote is origin, never automatic Gitee mirror pushes.

Rollback stops candidate processes and preserves evidence; restore an older binary only when storage compatibility is explicit. Unknown formats stay read-only. Signing/support/response proposals are in [ADR_REGISTER_EN.md](ADR_REGISTER_EN.md).

## 8. Complete Product Flow (Planned, RFC Section 19)

| Stage | Action/success evidence | On failure |
|---|---|---|
| PC-01 Trial | Preview supplied no-AI fixtures, inspect permissions/budgets/unknowns, separately authorize executable probes before running | No config hook execution or automatic tool installation |
| PC-02 Delivery | Select accepted results, choose a new nonoverlapping destination, verify hashes/exclusions/review references | Reject overwrite; incomplete output is not delivery; export is not external release |
| PC-03 Approval/revocation | Inspect inbox/candidate hash; verify blocked new dispatch and in-flight stop evidence after revocation | Remain paused if stopping fails; no hard-gate waiver |
| PC-04 Recovery/diagnostics | Read-only redacted plan identifies last acceptance, journal, processes and unknown effects | No manual lock deletion/blind redispatch; file restore cannot undo uploads/deployments |
| PC-05 Costs | Inspect reserved/settled/unknown holds and measurement sources; authorize increases explicitly | Timeout is not free; subscription mode guarantees authorized call/token caps, not monetary bills |
| PC-06 Retirement/uninstall | Stop managed writers/revoke authorization, verify complete backup and new-store restore drill; retain run/store by default | Preserve unknown formats; do not delete shared SessionBridge/toolchains; confirm data deletion separately |

Users decide whether to merge/release S1 delivery packages externally; ProofRail does not write source or execute Git automatically. Check support status, SBOM/licenses, signing identity and revocation freshness separately. Unknown offline revocation must not display current trust. Backups exclude credentials; restore environments require separate secure configuration.

For historical secrets, quarantine/block export/revoke credentials, then obtain a human retention/deletion decision; never rewrite old hashes as intact history. Deletion does not guarantee SSD/external-copy erasure. Reports stay local without default telemetry; support submissions contain only minimal redacted reproduction data. T020 now provides library-level export capability; T021–T024 remain pending, and the full product flow is still not an executable command tutorial today.
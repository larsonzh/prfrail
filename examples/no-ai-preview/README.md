# No-AI Preview Example

This example demonstrates PC-01 (`AT-16`) static preview without running commands, models, network checks, or credential input.

## Usage

From repository root:

```powershell
go run ./cmd/prfrail preview --chain .\proofrail.chain.json --json
```

Expected properties:

1. `preview.preview.mode` is `offline-read-only`.
2. `callCounters.commandCalls/networkCalls/modelCalls/credentialReads/versionProbes` are all `0`.
3. `unknowns` lists next actions for capabilities that need explicit authorization.
4. The command does not create run directories and does not mutate source/store state.

For noop-only chains, `preview.preview.outcome=ready`; if executable steps are present, outcome is `blocked` with deterministic `blockingEvidence`.

# Go Minimal Example

This deterministic `AT-15` fixture exercises a non-C Go language scope without running external commands or modifying a product repository.

From the ProofRail repository root:

```powershell
go run ./cmd/prfrail validate --chain .\examples\go-minimal\proofrail.chain.json --json
go run ./cmd/prfrail preview --chain .\examples\go-minimal\proofrail.chain.json --json
go run ./cmd/prfrail run --chain .\examples\go-minimal\proofrail.chain.json --run-id go-minimal --run-dir .\tmp\go-minimal-run --json
go run ./cmd/prfrail report --run-dir .\tmp\go-minimal-run --json
```

Expected facts:

1. Validation reports one task, one step, and `runnableInCli=true`.
2. Preview reports `mode=offline-read-only`, `outcome=ready`, and zero command, network, model, credential, and version-probe calls.
3. Run and report both return `chainState=COMPLETED`.
4. The chain fixture remains byte-for-byte unchanged. Runtime evidence is isolated under the explicit temporary run directory.
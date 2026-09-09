# Generic Minimal Example

This deterministic fixture exercises a generic language scope without running external commands or modifying a product repository.

From the ProofRail repository root:

```powershell
go run ./cmd/prfrail validate --chain .\examples\generic-minimal\proofrail.chain.json --json
go run ./cmd/prfrail preview --chain .\examples\generic-minimal\proofrail.chain.json --json
go run ./cmd/prfrail run --chain .\examples\generic-minimal\proofrail.chain.json --run-id generic-minimal --run-dir .\tmp\generic-minimal-run --json
go run ./cmd/prfrail report --run-dir .\tmp\generic-minimal-run --json
```

Validation reports one generic scope and a runnable noop-only chain. Preview remains offline and read-only, and runtime evidence is isolated under the explicit temporary run directory.
# C Minimal Example

This deterministic fixture exercises a C language scope without invoking a compiler or modifying a product repository. It proves configuration and noop orchestration independently of local compiler availability.

From the ProofRail repository root:

```powershell
go run ./cmd/prfrail validate --chain .\examples\c-minimal\proofrail.chain.json --json
go run ./cmd/prfrail preview --chain .\examples\c-minimal\proofrail.chain.json --json
go run ./cmd/prfrail run --chain .\examples\c-minimal\proofrail.chain.json --run-id c-minimal --run-dir .\tmp\c-minimal-run --json
go run ./cmd/prfrail report --run-dir .\tmp\c-minimal-run --json
```

Validation reports one C scope and a runnable noop-only chain. Preview remains offline and read-only, and runtime evidence is isolated under the explicit temporary run directory. This example does not claim that ProofRail executed `gcc`.
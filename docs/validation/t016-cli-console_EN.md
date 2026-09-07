# T016 CLI/Console Validation Report

[简体中文](t016-cli-console.md)

Date: 2026-09-08. Verdict: `T016 COMPLETE`. This slice delivers a no-IDE CLI baseline with `init`, `validate`, `config explain`, `run`, and `report`, and replaces placeholder entry behavior with a testable executor.

## Scope

1. `cmd/prfrail/main.go`
   - Switches the entrypoint to `internal/console` command execution and removes placeholder `init/run` handling.
2. `internal/console/config.go`
   - Adds chain-config models, strict JSON decoding (unknown-field rejection), and semantic validation.
   - Adds config search order: `--chain` → `./proofrail.chain.json` → `./proofrail.json`.
   - Adds `config explain` output with value-source attribution (builtin-default/chain/task).
3. `internal/console/cli.go`
   - Adds command routing and argument parsing for `version/init/validate/run/report/config explain`.
   - Adds stable `--json` output and explicit exit codes (0/1/2).
4. `internal/console/runtime.go`
   - Adds noop-only local execution and report summary loading on top of `internal/chain`.
   - Fail-closes with a non-zero exit for executable `code/build/verify` steps.
5. `internal/chain/models.go`
   - Exports `ValidateDefinition` so CLI validation reuses core chain semantics.

## AT-12 Coverage

1. No-IDE flow: `init/validate/config explain/run/report` can be executed directly from the terminal.
2. Parseable output: all core commands support `--json`, including failure responses.
3. No-color output: text mode emits no ANSI escape sequences.
4. No fake success for unsupported runtime behavior: executable steps are rejected by `run` with non-zero exit.
5. Explainability: `config explain` reports effective policy and source pointers.

## Tests

- `internal/console/cli_test.go`
  - `TestInitValidateAndExplainJSON`
  - `TestRunAndReportNoopChain`
  - `TestRunRejectsExecutableSteps`
  - `TestValidateRejectsUnknownField`
  - `TestTextOutputHasNoANSI`
  - `TestUnknownCommandUsageExitCode`

## Gates

1. `go test -count=1 ./...`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `2026-09-08` remote high-count race regression (Ubuntu VM `10.0.0.199`, Go 1.27.1 + gcc 13.3.0): `go test -race -count=20 ./internal/applier ./internal/console ./internal/taskdef` and `go test -race -count=1 ./...` passed on first pass (no proxy fallback needed).

## Boundary

1. `run` currently supports noop-only chains; full executable-step gate/adapter integration belongs to later slices.
2. Layered `proofrail.toml/workspace.toml` merge remains planned; this slice freezes single-file chain parsing and explain behavior.
3. The local `run` uses synthetic in-memory review/promotion decisions without persisted receipts; it demonstrates the noop loop only and is not acceptance evidence.
4. No commit and no push in this run.

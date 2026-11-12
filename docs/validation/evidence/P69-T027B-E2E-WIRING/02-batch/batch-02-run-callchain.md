# P69 批量辅助 ② 归档 — prfrail run 调用链枚举

载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 批量辅助首选；显式传 model）。实得 5288 字符（CHARS 自报）。
以下为批量辅助最终消息逐字原文。

# prfrail run call-chain enumeration

Paths relative to `d:\LZProjects\prfrail`. Categories: [P]=CLI parsing, [A]=assembly entry, [E]=task execution, [R]=routing.
Read-only enumeration; no code modified; no commands executed.

## 1. CLI entry
- `Execute` decl `cli.go:70`; `case "run":` `cli.go:97`; `cli.executeRun(...)` call `cli.go:98`. [P]+[R]

## 2. executeRun - decl `internal/console/cli.go:274` (all refs below are cli.go)
CLI parsing [P]:
- flag set `cli.go:275`; `--chain` `:276`; `--run-dir` `:277`; `--run-id` `:278`; `--json` `:279`
- `--request-record` parsed `:280`
- `--request-records` parsed: slice init `:281`; `set.Func` `:282`
- `--capability-record` `:289`; `--enforcement-record` `:290`; `--availability-record` `:291`; `--authorization-ledger` `:292`; `--cost-ledger` `:293`; `--agent-cli` `:294`; `--agent-version` `:295`
- `set.Parse(args)` `:296`; `cli.getwd()` `:303`; `ResolveChainConfigPath` `:307`; `LoadChainConfig` `:311`
- `FindExecutableSteps` `:315`; `FindSupervisedSteps` `:316`; `agentRunnerRunInputs` literal `:317-327`
Routing [R] - branch conditions:
- `--request-record`/`--request-records` mutual exclusion `:328`
- request-records with no executable/supervised step refused `:334`
- supervised vs AgentRunner exclusion: `len(supervised) > 0 && inputs.anyProvided()` `:339`
- executable steps with incomplete inputs refused `:343`
- three-way split: supervised `if len(supervised) > 0` `:375` -> `cli.executeSupervisedRun(...)` `:376`; AgentRunner `if len(executable) > 0` `:378` -> `cli.executeAgentRunnerRun(...)` `:379`; noop fall-through -> `ExecuteNoopRun(...)` `:381`; summary write `:385`
Definitions: `agentRunnerRunInputs` struct `cli.go:396`; `complete()` decl `cli.go:408`; `anyProvided()` decl `cli.go:430`; `runUsage` const `cli.go:388`

## 3. Supervised branch
- `executeSupervisedRun` decl `internal/console/supervised_cmd.go:34` [A]; `cli.supervisedAssemblyOptions(...)` call `:35`; `ExecuteSupervisedRun(...)` call `:39`; summary write `:43`
- `supervisedAssemblyOptions` decl `supervised_cmd.go:50` [A]; `cli.supervisedGateRunner(runDir)` call `:51`; `SupervisedAssemblyOptions{...}` literal `:55`; `supervisedGateRunner` decl `supervised_cmd.go:66`
- `ExecuteSupervisedRun` decl `internal/console/runtime_supervised.go:267` [A]; `supervisedClockAndIDs` `:271`; `supervisedAssemblyPlan` `:272`; `ensureAgentRunnerEventLogMarker` `:283`; `newSupervisedEngine(...)` call `:286`; parked replay branch `:291-292`; `engine.Run(ctx)` call `:294` [E]
- `newSupervisedEngine` decl `runtime_supervised.go:587` [E]; `newSupervisedRuntime` `:588`; `chain.NewFileEventStore` `:593`; `chain.New(ctx, chain.Options{...})` `:597`

## 4. AgentRunner branch
- `executeAgentRunnerRun` decl `internal/console/cli.go:454` [P] record decode:
  - `readRecord` closure `:459`
  - repeatable form: `readAgentRunnerRequestRecords(...)` call `:478` (helper decl `cli.go:580`)
  - single form: file read `:483`; `adapters.DecodeAgentRunnerRequestRecord` `:487`
  - capability: read `:492`; `adapters.DecodeAgentRunnerCapabilityRecord` `:496`
  - enforcement: read `:500`; `adapters.DecodeAgentRunnerEnforcementRecord` `:504`
  - availability: read `:508`; `adapters.DecodeAIAvailabilityRecord` `:512`
  - authorization ledger: path resolve `:516`; `loadAuthorizationLedger` `:520`
  - cost ledger: path resolve `:524`; file read `:528`; `tickets.DecodeCostLedgerRecord` `:532`; `tickets.LoadCostLedger` `:536`
  - agent CLI path resolve `:540`; `AgentRunnerAssemblyOptions` literal `:544-564`
  - `ExecuteAgentRunnerRun(ctx, options)` call `:565`
- `ExecuteAgentRunnerRun` decl `internal/console/runtime_agent.go:126` [A]:
  - Step 1 preconditions (comment `:139`): `agentRunnerAssemblyPlan` call `:144` (decl `:1063`); `agentRunnerAssemblyRecords` call `:148` (decl `:1087`); capability validate `:152`; enforcement validate `:156`; availability validate `:160`; cost-ledger check `:163`; run-root check `:166`; `agentRunnerValidateCLI` `:169` (decl `:1304`)
  - Step 2 run-root marker (comment `:173`): `ensureAgentRunnerEventLogMarker` call `:176` (decl `:1482`)
  - Step 3 replay store (comment `:180`): `adapters.NewAgentRunnerReplayStoreForRun` `:181`
  - Step 3b prebuild (comment `:186`): `store.PrebuildReplaySlots(...)` `:193`
  - Step 4 snapshot store (comment `:197`): `snapshot.NewStore` `:199`
  - Step 5 pinned-CLI launcher (comment `:204`): `adapters.NewAgentRunnerPinnedCLILauncher` `:207`
  - Step 6 per-task ports (comment `:220`): `agentRunnerBindTasks` call `:226` (decl `:1144`)
  - Step 7 runtime ports + wait carrier (comment `:231`): port maps `:235-241`; bind loop `:242`; pre-manifest port wrap `:279`; wait carrier `agentRunnerDispatchWaiter` `:285`; admission mux `:289`; `chain.StepRouter` literal `:293`
  - Step 8 engine (comment `:298`) [E]: `chain.NewFileEventStore` `:318`; `chain.New(ctx, chain.Options{...})` `:322`; `fail` closure `:368`; serial driver `for {` `:392`; `engine.Run(ctx)` `:393`

## 5. Noop branch
- `ExecuteNoopRun` decl `internal/console/runtime.go:33` [E]; `chain.NewFileEventStore` `runtime.go:34`; `localRuntime` literal `:39`; `chain.Options` literal `:40-55`; `chain.New(ctx, options)` `:56`; `engine.Run(ctx)` `:60`; summary return `:63`

CHARS: 5288

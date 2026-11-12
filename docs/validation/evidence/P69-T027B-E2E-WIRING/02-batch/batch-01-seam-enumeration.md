# P69 批量辅助 ① 归档 — P44 判定出口与生产装配接缝枚举

载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 批量辅助首选；显式传 model）。实得 5546 字符（CHARS 自报）。
以下为批量辅助最终消息逐字原文。

## A. P44 exported surface (internal/gates)

readiness.go:
- ErrReadinessProbeNotAcceptance — readiness.go:53 — sentinel: acceptance evaluation only holds level-1 probe evidence.
- ErrReadinessEvidenceIncomplete — readiness.go:56 — sentinel: missing or malformed readiness material (seal fields, hash syntax, unknown goal).
- ErrReadinessAcceptanceMissing — readiness.go:59 — sentinel: production evaluation without level-2 seal material.
- ErrReadinessProductionAT23Required — readiness.go:62 — sentinel: production evaluation without the AT23 acceptance hash.
- ReadinessGoal (type) — readiness.go:70 — selects the readiness level to evaluate.
- GoalProbe (const) — readiness.go:74 — level 1: a probe may run.
- GoalAcceptance (const) — readiness.go:76 — level 2: work counts as accepted.
- GoalProduction (const) — readiness.go:78 — level 3: production release admitted.
- Readiness (interface) — readiness.go:84 — sealed evaluation result; unexported tag method closes the level set.
- SealBinding (type) — readiness.go:89 — seal attempt carrying level-2 material (Attempt/TaskKey/SealHash).
- ReleaseBinding (type) — readiness.go:96 — level-3 admission material (Seal/AT23Hash/EvidenceRootHash).
- ReadinessInput (type) — readiness.go:104 — raw evaluation input (Permit/Seal/AT23Hash/ProbeEvidence/EvidenceRootHash); read-only.
- ProbeReadiness (type) — readiness.go:113 — level-1 readiness value holding a bound permit.
- AcceptanceReadiness (type) — readiness.go:118 — level-2 readiness value holding a seal binding.
- ProductionReadiness (type) — readiness.go:124 — level-3 readiness value holding a release binding.
- (*ProbeReadiness).Permit — readiness.go:134 — accessor returning the admitting probe permit.
- (*AcceptanceReadiness).Acceptance — readiness.go:142 — accessor returning the seal binding.
- (*ProductionReadiness).Release — readiness.go:150 — accessor returning the release binding.
- EvaluateT027BReadiness — readiness.go:159 — pure decision function; the only constructor of a Readiness level.

probe_permit.go:
- ErrProbePermitMissing — probe_permit.go:66 — sentinel: probe conclusion/permit operation with no bound permit.
- ErrProbePermitInUse — probe_permit.go:69 — sentinel: the single claim is held by an unfinished attempt.
- ErrProbePermitConsumed — probe_permit.go:71 — sentinel: the quota is spent for good.
- ProbePermit (type) — probe_permit.go:107 — single-use probe admission token linearised by a durable on-disk marker.
- NewProbePermit — probe_permit.go:126 — constructor: fresh unused permit bound to (id, path); maxRequests normalised to >=1.
- LoadProbePermit — probe_permit.go:140 — rebuilds a permit from its marker file; missing marker = unused, existing marker = spent.
- (*ProbePermit).Require — probe_permit.go:210 — claims the quota; admission persists first via O_CREATE|O_EXCL marker + Sync.
- (*ProbePermit).Complete — probe_permit.go:253 — marks a claimed permit spent; idempotent when spent, refuses a never-required permit.

## B. Assembly seams (internal/console)

- ExecuteAgentRunnerRun (runtime_agent.go): func :126; precondition step starts :139 (Step 1 comment; first check call agentRunnerAssemblyPlan :144); event-log marker creation step starts :173 (Step 2 comment; runEventLogPath :175, ensureAgentRunnerEventLogMarker :176).
- ExecuteSupervisedRun (runtime_supervised.go): func :267; precondition step starts :268 (ctx.Err) / :272 (supervisedAssemblyPlan); event-log marker creation :282-283 (runEventLogPath :282, ensureAgentRunnerEventLogMarker :283).
- executeRun (cli.go): func :274; routing block to the two assembly entries :375-380 (if len(supervised) > 0 :375 → executeSupervisedRun :376; if len(executable) > 0 :378 → executeAgentRunnerRun :379); immediately preceding decision blocks: event-log-already-exists check :365-369, cfg.ToDefinition() error check :372-374.
- supervised_cmd.go: executeSupervisedRun :34 (calls supervisedAssemblyOptions :35, ExecuteSupervisedRun :39); supervisedAssemblyOptions :50 (builds SupervisedAssemblyOptions :55).

## C. Existing gates references in console

- runtime_supervised.go — import :18; usages :234 (GateRunner *gates.Runner), :235 (GateHooks map[string]gates.Hook), :466 (gateRunner *gates.Runner), :467 (gateHooks map[string]gates.Hook), :1371 (gates.DefaultSecretScanner()), :1434 (gates.DefinitionHash), :1438 (gates.Request{} / gateRunner.Run).
- supervised_cmd.go — import :12; usages :66 (*gates.Runner return), :71 (gates.NewFileResultStore), :75-77 (gates.Runner{, gates.GuardExecutor{}, gates.DefaultSecretScanner()).
- runtime_supervised_test.go (test) — import :16; usages :39-49, :58-74, :125-132 (fake executor/results/hook; gates.Runner + gates.Hook assembly).

## D. No-consumer fact

- Repo-wide grep (internal/**/*.go) for the P44 exported symbols (EvaluateT027BReadiness, ReadinessGoal, GoalProbe/GoalAcceptance/GoalProduction, Readiness, SealBinding, ReleaseBinding, ReadinessInput, ProbeReadiness, AcceptanceReadiness, ProductionReadiness, ErrReadiness*, ProbePermit, NewProbePermit, LoadProbePermit, ErrProbePermit*) matched 6 files, all under internal/gates.
- No production consumer in internal/console, internal/chain, internal/adapters: 无生产消费方（仅 gates 包内测试引用）.
- Test files referencing them: internal/gates/readiness_test.go, internal/gates/readiness_counterexample_test.go, internal/gates/probe_permit_test.go, internal/gates/probe_permit_counterexample_test.go.
- Note: console consumes other gates symbols (Runner/Hook/etc., section C) but none of the P44 readiness/probe-permit exports.

CHARS: 5546

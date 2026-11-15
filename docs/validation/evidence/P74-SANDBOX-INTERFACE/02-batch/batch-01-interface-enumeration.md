# BATCH-01 现有 guard / chain / adapters / permit 事实枚举

- 载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单批量辅助）
- 任务限定：只列事实、不判断、不排序、不做 Top-N 择优
- 实派次数：第 1 次（§二 强制 ≥2）

## A. `internal/guard/` 包

**文件清单**（14 个 .go）：`process.go`、`lease.go`、`process_linux.go`、`process_linux_io.go`、`process_windows.go`、`process_unsupported.go`、及测试文件 7 个（`process_test.go`、`process_stop_windows_test.go`、`process_spawn_ownership_windows_test.go`、`process_spawn_ownership_test.go`、`process_reclaim_test.go`、`process_reaper_test.go`、`process_linux_test.go`）。

- `internal/guard/process.go:46` — `ProcessSpec`：

```go
type ProcessSpec struct {
	Command string
	Args    []string
	Dir     string
	Env     []string
	Stdout  io.Writer
	Stderr  io.Writer
	// OnSpawn, when set, runs after the process exists with a complete identity
	// and its ownership anchor verified, and before it is allowed to run.
	OnSpawn func(ctx context.Context, identity ProcessIdentity) error
}
```

- `process.go:122` — `ProcessIdentity`：

```go
type ProcessIdentity struct {
	PID        int    `json:"pid"`
	StartToken string `json:"startToken"`
}
```

- `process.go:136` — `ManagedProcess`（结构字段含 `cmd *managedCommand`、`identity`、`platform *platformProcess`、`verify`、`waitOnce sync.Once`、`waitDone chan struct{}`、`waitErr error`）：

```go
type ManagedProcess struct {
	cmd      *managedCommand
	identity ProcessIdentity
	platform *platformProcess
	verify   func(context.Context) error
	waitOnce sync.Once
	waitDone chan struct{}
	waitErr  error
}
```

- `ProcessSupervisor`：`internal/guard/` 全目录 grep 无匹配（不存在）。
- 方法/函数签名（file:line，签名行原文）：
  - `process.go:190`：`func StartManaged(ctx context.Context, spec ProcessSpec) (*ManagedProcess, error) {`
  - `process.go:266`：`func RunManaged(ctx context.Context, spec ProcessSpec, grace time.Duration) (ProcessResult, error) {`
  - `process.go:280`：`func (process *ManagedProcess) Run(ctx context.Context, grace time.Duration, startedAt time.Time) (ProcessResult, error) {`
  - `process.go:290`：`func (process *ManagedProcess) runManaged(ctx context.Context, grace time.Duration, startedAt time.Time) (ProcessResult, error) {`
  - `process.go:352`：`func (process *ManagedProcess) Identity() ProcessIdentity { return process.identity }`
  - `process.go:371`：`func StopProcessIdentity(ctx context.Context, identity ProcessIdentity, grace time.Duration) (TerminationEvidence, error) {`
  - `process.go:426`：`func (process *ManagedProcess) Terminate(ctx context.Context, grace time.Duration) (TerminationEvidence, error) {`
  - `process.go:465`：`func (process *ManagedProcess) Reclaim(ctx context.Context) error {`
  - 另有 `process.go:159`：`func (process *ManagedProcess) waitResult() <-chan struct{}`；`process.go:175`：`func (process *ManagedProcess) verifyGone(ctx context.Context) error`；`process.go:243`：`func (process *ManagedProcess) abortStartup(cause error) error`；`process.go:355`：`func InspectProcess(pid int) (ProcessIdentity, error)`；`process.go:361`：`func ProcessAlive(identity ProcessIdentity) (bool, error)`。

## B. `internal/chain/` 启动端口

- `internal/chain/agent_runner_launcher.go:14`：

```go
type AgentRunnerLauncher interface {
	StartAgentRunnerProcess(ctx context.Context, launch AgentRunnerLaunchRequest) (AgentRunnerLaunchResult, error)
}
```

- `agent_runner_launcher.go:26`：

```go
type AgentRunnerLaunchRequest struct {
	RequestID string
	RunID     string
	AdapterID string
	// Process parameters (additive, A6).
	Command       string
	Args          []string
	Dir           string
	Env           []string
	Grace         time.Duration
	WorkspaceRoot string
}
```

## C. `internal/adapters/agent_runner_pinned_cli.go`

- spawn 路径，`:218`（OnSpawn 钩子定义于 `:189`）：

```go
	process, err := guard.StartManaged(ctx, guard.ProcessSpec{
		Command: launcher.config.Executable,
		Args:    append([]string(nil), args...),
		Dir:     launch.Dir,
		Env:     append(launcher.childEnv(launch.Env), agentRunnerEvidenceDirEnvVar+"="+evidenceDir),
		Stdout:  stdout,
		Stderr:  stderr,
		OnSpawn: onSpawn,
	})
```

- 停机路径，`:300`（`StopAgentRunnerProcess` 起于 `:284`）：

```go
	// The launch may have been started with a tighter grace than the configuration,
	// and the operator path must not be more patient than the run it is stopping.
	proof, stopErr := guard.StopProcessIdentity(ctx, identity, launcher.stopGrace(launchID))
	launcher.releaseLaunch(launchID)
	return proof, stopErr
```

## D. P69 permit 检查点

- `internal/console/runtime_agent.go:197`（前 4 行 + 后 10 行）：

```go
	// The probe permit quota is claimed only after every precondition passed: a
	// configuration refusal must not burn the one-shot quota. The release defers
	// Complete so the quota is spent on every path past this point, success or
	// failure (N4: a claimed permit is never refunded).
	releaseProbePermit, err := requireProbePermit(ctx, options.ProbePermit)
	if err != nil {
		return RunSummary{}, err
	}
	defer releaseProbePermit()

	// Step 2: the run-root marker comes first. The replay store refuses to be
	// constructed without it, so the order is a correctness rule, not a script.
	eventLog := runEventLogPath(options.RunRoot)
	if err := ensureAgentRunnerEventLogMarker(eventLog); err != nil {
```

- `internal/console/runtime_supervised.go:297`（`freshStart` 判定）与 `:305`：

```go
	freshStart := false
	if _, statErr := os.Stat(SupervisedInteractionsLedgerPath(options.RunRoot)); statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			return RunSummary{}, fmt.Errorf("supervised probe permit ledger check: %w", statErr)
		}
		freshStart = true
	}
	if freshStart {
		releaseProbePermit, permitErr := requireProbePermit(ctx, options.ProbePermit)
```

- `requireProbePermit` 定义：`internal/console/readiness_gate.go:14`。
- 出现顺序（runtime_agent.go，仅行号序）：permit `:197` → Step 2 run-root marker `:205-210`（`ensureAgentRunnerEventLogMarker :209`）→ Step 3 replay store `:217` → Step 3b `PrebuildReplaySlots :227` → Step 5 launcher `:238` → Step 7 绑定 → 首次派发 `engine.Run(ctx)` `:429`。runtime_supervised.go：`freshStart :297-303` → permit `:305` → marker `:312` → `engine.Run(ctx)` `:328`。

## E. `internal/gates/probe_permit.go`

- `:210` `Require` 关键行：

```go
func (p *ProbePermit) Require(ctx context.Context) error {
	if p == nil {
		return fmt.Errorf("%w: nil probe permit", ErrProbePermitMissing)
	}
```

```go
	// The marker file, not the cache, is the linearisation point: persisting the
	// claim is what makes an admission final.
	if err := p.persistClaim(); err != nil {
```

```go
	if err := ctx.Err(); err != nil {
		// The attempt was abandoned after the claim became durable; the quota stays
		// consumed (an abandoned attempt keeps the claim, it never refunds).
		return fmt.Errorf("probe permit %q: require: %w", p.id, err)
	}
```

- `:253` `Complete`（CAS）：

```go
func (p *ProbePermit) Complete() error {
	if p == nil {
		return fmt.Errorf("%w: nil probe permit", ErrProbePermitMissing)
	}
	if p.state.CompareAndSwap(uint32(probePermitClaimed), uint32(probePermitSpent)) {
		return nil
	}
```

- 配套：`:271` `func (p *ProbePermit) refusal() error`；`:284` `func (p *ProbePermit) persistClaim() error`（注释：`O_EXCL makes the file itself the cross-process lock: exactly one writer can create it.`）。

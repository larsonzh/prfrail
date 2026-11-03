# GLM 批量辅助 #2 —— internal/console CLI 端到端测试范式枚举（原文归档）

载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`；只读检索；实得输出字符数 7983（自报 `OUTPUT_CHARS=7983`）。

## 1. CLI 调用范式

`NewCLI(` 测试出现:cli_test.go:404(helper `newTestCLI`)、runtime_supervised_test.go:409、:871。`.Execute(context` 直呼:cli_test.go:418(`runCLI` 内)、runtime_supervised_test.go:411、:873。经 `runCLI` 驱动:

- `TestInitValidateAndExplainJSON` cli_test.go:22(init/validate/explain 回路)
- `TestRunAndReportNoopChain` cli_test.go:90(run noop+report)
- `TestLanguageExamplesCompleteReadOnlyCLILoop` cli_test.go:142
- `TestRunRejectsExecutableSteps` cli_test.go:227(输入不全→code=1)
- `TestValidateRejectsUnknownField` cli_test.go:256
- `TestTextOutputHasNoANSI` cli_test.go:281
- `TestPreviewJSONAndTextParityWithoutSideEffects` cli_test.go:297
- `TestPreviewReadyForNoopChain` cli_test.go:361
- `TestUnknownCommandUsageExitCode` cli_test.go:391
- `TestSupervisedReportRenderGolden` runtime_supervised_test.go:401(report 子命令)
- `TestRunCommandRefusesMixedSupervisedAgentRunnerInputs` runtime_supervised_test.go:869

无 CLI 层传满八件套的既有测试(grep `--capability-record`/`--agent-cli`/`--agent-version` 于 *_test.go 0 命中)。最贴近 run+记录旗标样例(runtime_supervised_test.go:869-880):

```go
func TestRunCommandRefusesMixedSupervisedAgentRunnerInputs(t *testing.T) {
	chainPath := writeSupervisedChainConfig(t)
	cli := NewCLI("test")
	var stdout, stderr strings.Builder
	code := cli.Execute(context.Background(), []string{"run", "--chain", chainPath, "--request-record", "declared.json"}, &stdout, &stderr)
	if code == exitSuccess {
		t.Fatal("P1: a supervised definition with AgentRunner inputs must be refused")
	}
	if !strings.Contains(stderr.String(), "mutually exclusive") {
		t.Fatalf("P1: the refusal must name the mutual exclusion: %s", stderr.String())
	}
}
```

另一式(cli_test.go:237-253):`runCLI` 返回 (code, stdout, stderr),`code != 1` Fatal,`decodeResponse` 断 `result.OK==false`;退出码 exitSuccess=0/exitFailure=1/exitUsage=2(cli.go:23-25)。

## 2. 链配置构造范式

- 生成:`cfg := DefaultChainConfig("chain-one")`(cli_test.go:231、260、301、365;runtime_supervised_test.go:1519 等),默认 1 task/1 step。
- 改造:`cfg.Tasks[0].Steps = []StepConfig{{ID:"build-one", Kind:"build", Hooks:[]string{"go-build"}}}`(cli_test.go:232-236)。
- 落盘:`EncodeChainConfig(cfg)` → `os.WriteFile(path, content, 0644)`(cli_test.go:237-241);supervised 式先 `ValidateChainConfig` 再写 t.TempDir() 下 `DefaultChainConfigName`(runtime_supervised_test.go:1543-1555)。
- TaskConfig 最小字段集(唯一多字段样例 runtime_supervised_test.go:1529-1541):`TaskConfig{ID, Review:"policy", Steps:[]StepConfig{{ID, Kind:"code", Execution:<mode>, HandoffPolicy:…}}}`。IsolatedWorkspace 仅以代码层 `chain.Step{ID, Kind:"code", Mode: chain.IsolatedWorkspace}` 表达;TaskConfig 侧无 isolated 样例。
- 多 task 样例:存在,但为 `chain.Definition` 形态非 TaskConfig:`multiTaskIsolatedDefinition()`(runtime_multitask_red_test.go:57-71,2 task×1 isolated code step)。TaskConfig 形态多 task 样例:无。

## 3. 记录文件构造范式

- 入口恒为 `adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{...})`,返回 (record, error);RecordHash 构造器内算,测试无需手工,自写 digest 仅用于 Evidence 域(runtime.go:125 `digest`)。
- helper:`multiTaskRequestRecord(t, taskID, stepID, requestID, parentHash)` runtime_multitask_agent_test.go:53-80;`driverRequestRecord(t, runID, taskID, stepID, requestID, parentHash)` runtime_multitask_driver_test.go:154-183;单任务全量内联 runtime_agent_e2e_test.go:329-356。
- 必填字段集(runtime_agent_e2e_test.go:329-348):RequestID、RunID、TaskID、StepID、Attempt、CreatedAt(timestampLayout 字符串)、AdapterID、Mode:"create"、WorkspaceHash、ContextHash、ParentSnapshotHash、AuthorizationHash、BudgetHash、AllowedTargets、AllowedEffects、EffectMappingVersion("1")、EffectMappingHash(契约常量 assemblyEffectMappingHash,runtime_agent_e2e_test.go:35-37)、Evidence。WorkspaceHash/ParentSnapshotHash:e2e 用真实 baseline digest(:191-196),offline harness 用假 digest。
- 落盘:`json.Marshal(record)`(internal/adapters/agent_runner_test.go:41;approvals_test.go:54 同式)→ os.WriteFile;严格解码 `adapters.DecodeAgentRunnerRequestRecord`(internal/adapters/agent_runner.go:73);注册键 `agentRunnerRequestRecordKey(taskID, stepID)`=taskID+"/"+stepID(runtime_agent.go:750-753)。

摘录(runtime_multitask_agent_test.go:53-80,逐字):

```go
func multiTaskRequestRecord(t *testing.T, taskID, stepID, requestID, parentHash string) adapters.AgentRunnerRequestRecord {
	t.Helper()
	record, err := adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{
		RequestID:            requestID,
		RunID:                multiTaskRunID,
		TaskID:               taskID,
		StepID:               stepID,
		Attempt:              1,
		CreatedAt:            multiTaskClock,
		AdapterID:            "multitask-stub",
		Mode:                 "create",
		WorkspaceHash:        multiTaskTestDigest("workspace", taskID),
		ContextHash:          multiTaskTestDigest("context", taskID),
		ParentSnapshotHash:   parentHash,
		AuthorizationHash:    multiTaskTestDigest("authorization", taskID),
		BudgetHash:           multiTaskTestDigest("budget", taskID),
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"workspace-write", "local-process"},
		EffectMappingVersion: assemblyEffectMappingVersion,
		EffectMappingHash:    assemblyEffectMappingHash,
		Evidence:             []string{multiTaskTestDigest("request-evidence", taskID)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return record
}
```

## 4. 断言与辅助函数清单

均为同包(package console)可复用 helper;`digest` 为生产函数:

- `newTestCLI` cli_test.go:403 `(cwd string) CLI`(注入 cli.now/cli.getwd)
- `runCLI` cli_test.go:414 `(t *testing.T, cli CLI, args ...string) (int, string, string)`
- `decodeResponse` cli_test.go:424 `(t *testing.T, text string) testResponse`
- `digest` runtime.go:125;`multiTaskTestDigest` runtime_multitask_agent_test.go:34 `(parts ...string) string`
- `multiTaskRequestRecord` runtime_multitask_agent_test.go:53(见 §3);`driverRequestRecord` runtime_multitask_driver_test.go:154
- `multiTaskFacts` runtime_multitask_agent_test.go:82 `(taskID string) chain.AgentRunnerFrozenFacts`
- `multiTaskIsolatedDefinition` runtime_multitask_red_test.go:57 `() chain.Definition`
- `requireDefinitionShapeAccepted` runtime_multitask_red_test.go:107 `(t *testing.T, definition chain.Definition)`
- `requireNoRetiredRefusal` runtime_multitask_red_test.go:115 `(t *testing.T, err error, refusals ...string)`
- `writeSupervisedChainConfig` runtime_supervised_test.go:1517 `(t *testing.T) string`
- `newAssemblyFixture` runtime_agent_e2e_test.go:157 `(t *testing.T, knobs assemblyKnobs) *assemblyFixture`(装配全量 options)
- `requireAgentRunnerZeroSpawn` runtime_agent_e2e_test.go:424 `(t *testing.T, runRoot string)`
- `requireAgentRunnerEventReason` runtime_agent_e2e_test.go:440 `(t *testing.T, runRoot, reason string)`
- `agentRunnerDirEntries` runtime_agent_e2e_test.go:408;`zeroModelRunProjection` runtime_zero_model_e2e_test.go:41;`writeAIAvailabilityFixture` ai_test.go:14

## 5. 命名与放置建议(事实)

- 命名:多任务 `runtime_multitask_{red,agent,driver,supervised}_test.go`;e2e `runtime_agent_e2e_test.go`、`runtime_agent_concurrent_e2e_test.go`、`runtime_zero_model_e2e_test.go`;纯 CLI 命令集中 cli_test.go。
- `t.TempDir()` 普遍使用(cli_test.go:23/91/392 及各 fixture);`t.Setenv` 包内 0 处。
- 全局状态:无 TestMain;包级 var 仅 assemblyStubOnce/assemblyStubPath/assemblyStubErr(runtime_agent_e2e_test.go:40-44),余为包级 const(assembly*/multiTask*/supervised* 各族);无 os.Chdir,cwd 经 `newTestCLI` 注入;`TestZeroModelAgentCLIHelper`(runtime_agent_concurrent_e2e_test.go:24)占用测试二进制按 `--` 再入的约定。

## 6. 端到端可跑性事实

- 外部二进制依赖仅 `buildAssemblyStub` 的 `go build ./tools/agent-stub`(runtime_agent_e2e_test.go:60-68,失败 t.Skipf);`TestZeroModelAgentCLIHelper` 以测试二进制自身作假 agent CLI;无网络;noop CLI 回路零外部依赖。
- 临时目录假文件可达 request-record 严格解码,机械依据(internal/console/cli.go):
  1. `--request-record` 旗标:cli.go:280;
  2. 路由需非 supervised 且 `len(executable)>0`(cli.go:359-362),且 `inputs.complete()` 八值全非空(cli.go:378-400);
  3. 前置仅:LoadChainConfig(cli.go:303)、`evidence.ValidID`(cli.go:334)、event log 不存在(cli.go:346-352);
  4. `executeAgentRunnerRun`(cli.go:428)spawn 前仅 os.ReadFile+`adapters.DecodeAgentRunnerRequestRecord`(cli.go:448),失败经 writeCommandError 返回非零。
- ⇒ 八件套旗标传临时目录内 `json.Marshal(record)` 产物路径即达解码点,无需 agent-stub 或网络。

OUTPUT_CHARS=7983

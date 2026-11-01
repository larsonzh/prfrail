# ProofRail Validation Report — P45-T027B-E2E-ZERO-MODEL

| Field | Content |
|---|---|
| Slice | `P45-T027B-E2E-ZERO-MODEL` |
| Type | Assembly-validation `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `eb36992` |
| Date | 2026-11-01 |

## Goal and Scope

- **Goal**: give the zero-model three-task chain end-to-end assembly evidence — T1 executes and gets accepted (`PASSED` + parent-snapshot binding), T2 fails naturally (`FAILED`) and T3 never starts, with no model call anywhere.
- **This slice is assembly validation with zero production-code change**: `git diff --name-only -- internal/` is empty; the deliverables are three new test files, one example fixture and README command-path fix, plus the doc backfill.
- **E2E placement**: the main E2E lives in `internal/chain` (engine level, read-only reuse of `threeTaskDefinition` / `testOptions` / `fakeSteps.failTask`); the two console legs cover the `ExecuteNoopRun` multi-task noop chain and the natural `localRuntime` refusal of a non-noop step, plus `examples/no-ai-preview` staying runnable with zero AI; the concurrency leg lives in console (same launcher / same runRoot, 8-way Start with exactly one winner).
- **The console single-task limit stays** (`runtime_agent.go:774`, `runtime_supervised.go:335`): multi-task assembly would collide with A6's one identity mirror per runRoot, the exclusion contract DR13-FIX pinned down.
- **Mechanical zero-model criteria**: the engine leg asserts `StepIntents==nil`, `ExecutionTarget==DefaultExecution`, `AgentRunnerFacts==nil` and no `agentRunnerDispatchedReason` event; the preview leg asserts `RunnableInCLI==true`, `Outcome==ready`, an empty `BlockingEvidence` and every step `kind==noop`.
- **Counter semantics**: the five preview counters are static constants of `BuildPreviewReport`; they are a text lock on the README promise, not a model-channel sentinel.
- **Scope**: `internal/chain/**`, `internal/console/**`, `examples/no-ai-preview/**` and the doc backfill.
- **Boundary**: `internal/adapters/**` untouched (read-only calls of exported API only); `docs/CONTRACTS*` and `docs/t027/CHECKLIST_v34.md` untouched; no gitee push; `commit` / `push` wait for authorisation at step 10.

## Verdict Summary

- The task-state order holds independently on both the engine and console sides: T1 `PASSED`, T2 `FAILED`, chain `FAILED`; T3 never appears in the workspace calls, the event stream, the projection keys or the step request set.
- The parent-snapshot binding holds point by point: T1's parent is `baselineHash`, T2's is `acceptedOne`, T3's is `acceptedTwo`, and the final `AcceptedParent` is `acceptedThree`; the event chain passes `evidence.VerifyEventChain` and carries 3 seal records.
- Both zero-model sides hold, and the "counters are not a sentinel" reading is written into the test comments and into this report.
- Concurrent Start pins the DR13-FIX strong input with out-of-package behavioural equivalence: exactly one winner, every loser hits one of two sentinels, the winner is `ProcessAlive`, a late start is refused while the winner is alive, and the slot is released after stop so a new start succeeds.
- The first run exposed a real defect (the missing example fixture); it was handled as a test-fixture fix and no production code was touched.

## Review Chain Verdicts

- **① Architecture**: `deep-reasoner` + DeepSeek V4 Pro; the zero-model three-task E2E assembly plan, `DESIGN: DONE`; placement trade-off, red-test list, concurrency assertion design and 4 unconfirmed items in `01-arch-design/review.md`.
- **③ First red run**: 7 green and 1 red out of 8 cases; the red one is the preview leg, rooted in the missing `examples/no-ai-preview/proofrail.chain.json` (the README example is not runnable), a real defect.
- **② Fix round**: after adding that fixture and the minimal README command-path fix the suite is fully green.
- **③ Second run**: per the ⑤ MINOR the tautological assertions were dropped and a mutation was self-proved — flipping one fixture step from `kind=noop` to `kind=code` turns the preview leg red, and after restoring it the file SHA256 equals the pre-change value.
- **④ Gates**: `TOTAL_FAIL=0`, readings in the gate-readings section.
- **⑤ Pre-review**: `PRE-REVIEW: PASS WITH FIXES` (2 MINOR: tautological counter assertions; a ledger evidence row exceeding the actual assertions).
- **⑥ Scan**: `INDEPENDENT SCAN: PASS`, Section B empty, Section D lists 4 uncovered areas; slice quota 1/3.
- **⑦ Final review**: `FINAL REVIEW: PASS`, plus 1 MINOR and 6 release conditions, each discharged in the deviations section; slice quota 1/3.
- **R / D**: R = 0, D = 0 (not inherited).

## Six-Stage Scan Transcript

Source `docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/06-scan/review.md`; the fence transcribes the ⑥ body in A→B→C→D→E order, and three normalisations are disclosed outside it.

```sixstage
## Section A 扫描范围与输入
范围：只读核查 internal/chain/engine_zero_model_e2e_test.go、internal/console/runtime_zero_model_e2e_test.go、internal/console/runtime_agent_concurrent_e2e_test.go、examples/no-ai-preview/proofrail.chain.json 与 examples/no-ai-preview/README.md，聚焦三任务状态序、T3 未启动、零模型判据与并发 Start 互斥。

## Section B 发现清单
| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
| --- | --- | --- | --- | --- |
|  |  |  |  |  |

## Section C 与 ①/⑤ 的一致性核对
1. T3 不启动的多路互证成立，且没有假绿空间。
   - 在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 中，T2 失败后会检查：`workspaces.calls` 只有 T1/T2，`projection.TaskStates` 不含 T3，`store.events` 中不出现 T3 任务实体，`steps.requests` 中不出现 T3，并且 `len(steps.requests)==3` 这个断言强制卡住“额外发起”。
   - 在 [internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 中，控制台对 T2 用 `kind="code"`、`Mode=ManagedChangeSet` 的真实拒绝路径断言了 `localRuntime` 自然失败，且 `task-three` 不在 `taskIDs` 中。
   - 这不是单一字段判读，而是“任务启动序/事件流/投影键/步骤请求集”四重交叉，真实覆盖了“T3 never started”。

2. 零模型证明机械、且没有把计数器误当哨兵。
   - 引擎腿在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 里直接断言：`options.StepIntents == nil`；并在执行时检查 `request.ExecutionTarget == DefaultExecution`，`request.AgentRunnerFacts == nil`，并且事件中不存在 `agentRunnerDispatchedReason`。
   - 预览腿在 [internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 中，计数器依然保留为 0，但测试明确说明它只是“README 运营承诺的静态锁定”，真正机械判断是 `RunnableInCLI`、`Outcome == ready`、`BlockingEvidence` 为空、且所有 `preview.Steps` 都是 `noop`。这符合“可证伪 + 计数器非哨兵”的要求。

3. 父快照绑定与事件链语义一致，且无越界。
   - [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 断言：T1 的父是 `baselineHash`，T2 的父是 `acceptedOne`，T3 的父是 `acceptedTwo`，终态 `AcceptedParent == acceptedThree`；这与真实链语义一致。
   - 事件链断言也直接调用了 [internal/evidence/event.go](internal/evidence/event.go) 的 `VerifyEventChain`，校验 `Sequence` 单调、`LastEventHash` 与最终事件 hash 一致，并确认 `seal` 数量 3。该证明与真实 engine state machine 一致，不是“测试造假”。

4. 并发腿是同 launcher/runRoot 的真正排他场景，且无墙钟判据。
   - [internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go) 复现的是同一 `launcher` 与同一 `runRoot` 下 8 个并发 start；赢家唯一，失败者只可看到 `ErrAgentRunnerLaunchConflict` / `ErrAgentRunnerProcessDuplicate`，并验证 `ProcessAlive`、slot owner 与 mirror 一致。
   - 这与 DR13-FIX 强输入一致，且没有使用 10s freshness 作为排他依据；它使用“状态/存在性/完成序”而非墙钟，是合理的。

5. 夹具和 README 与严格解码一致，且最小必要变更。
   - [examples/no-ai-preview/proofrail.chain.json](examples/no-ai-preview/proofrail.chain.json) 结构符合 [internal/console/config.go](internal/console/config.go) 的 `LoadChainConfig` 与 [schemas/chain.schema.json](schemas/chain.schema.json) 的 schema：`schemaVersion`、`chain`、`tasks`、`workspace` 全部合法，且文件大小 1044 bytes，起始 `0x7B 0x0A 0x20` 表明无 BOM，符合要求。
   - [examples/no-ai-preview/README.md](examples/no-ai-preview/README.md) 的修正是命令路径修正，最小必要，且与预览示例的行为一致。

6. 边界/卫生与证据包自洽性。
   - 实际只读边界被尊重：新测试文件是 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go)、[internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go)、[internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go)，生产代码未被改动；`git diff --name-only -- internal/` 在本轮检查中没有输出 tracked diff，只有这 3 个新测试文件分布在 internal tree，符合“零改动生产代码”的要求。
   - 证据包的 [docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md](docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md) 明确写着“步骤 0 骨架落盘；后续载荷按轮追加”，所以它当前是合法的 start-of-slice skeleton，而不是最终闭环包；该状态在文件中自说明，自洽。
   - Fresh verification evidence: `go test -count=1 ./...` returned all package lines as `ok` with no `FAIL` lines, including the internal/chain and internal/console packages.

## Section D 未覆盖面与残余风险
- 这次只读复核未对真实 VS Code UI/浏览器路径做端到端点击验证；它依赖代码和事件日志的静态断言。
- 未对 Windows 之外的 OS/locale 运行环境做差异验证；并发 snapshot 与 ProcessAlive 的行为可能受平台差异影响。
- 证据包未进入最终 06/07/08 载荷阶段，因此未覆盖终审与 CI 归档闭环。
- 例子中的 CLI 命令路径已被核验为合法 JSON 例子，但未运行 `go run ./cmd/prfrail preview --chain ...` 的真实命令行交互输出。

## Section E 结论要旨
- 只要把 [examples/no-ai-preview/proofrail.chain.json](examples/no-ai-preview/proofrail.chain.json) 中任一 step 的 `kind` 改成 `code`，[internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 会立即让 `preview.Steps` 或 `BlockingEvidence` 断言失败。
- 只要让 T2 在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 里继续推进到 T3 事件，`len(workspaces.calls)` 或 `task-three` 事件断言会立刻变红。
- 只要把并发测试里的“唯一赢家”逻辑去掉，或让 loser 观察到非 `Conflict`/`Duplicate` 错误，[internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go) 会直接失败。

INDEPENDENT SCAN: PASS
```

- Normalisation 1: the ⑥ section titles `## A 结论` / `## B 缺陷` / `## C 核验结论` / `## D 未覆盖面` / `## E 可证伪性` were renamed inside the fence to this template's section names; the five section bodies are verbatim.
- Normalisation 2: the Section A scope line is taken from the file set ⑥ verified in Section C item 6; the archived original carries no separate scope line.
- Normalisation 3: the trailing `§A: INDEPENDENT SCAN: PASS` line of ⑥ is normalised into the fence-final `INDEPENDENT SCAN: PASS`; the empty Section B header and empty row are kept verbatim.

## Finding Closure

Severity follows the review wording and is not bolded; the ⑦ conclusion closed set is `FINAL REVIEW: PASS` (0 new executable defects).

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| ⑤ r1 1 | The five preview counter assertions are tautological; no implementation change can turn them red | MINOR | `internal/console/runtime_zero_model_e2e_test.go` | Closed: the test comment declares them a static text lock on the README promise, not a mechanical criterion; the mechanical criteria are `RunnableInCLI` / `Outcome` / `BlockingEvidence` / step `kind` |
| ⑤ r1 2 | The ledger evidence row named `registry.len()==1` and slot owner/mirror agreement, beyond the actual assertions | MINOR | `docs/t027/REMAINING_SLICES.md`, `_EN` | Closed: the acceptance row now states the out-of-package equivalence, and D7 below registers the missing direct assertion |
| ⑦ B1 | The README line number and command snippet inside the failure message were stale | MINOR | `internal/console/runtime_zero_model_e2e_test.go` | Closed: the message now names `README.md:10` and `--chain .\examples\no-ai-preview\proofrail.chain.json` |

## Deliverables

- `internal/chain/engine_zero_model_e2e_test.go`: the engine-level main E2E (T1 / T2 / T3 three states + parent-snapshot binding + sealed monotonic event chain + zero AgentRunner dispatch).
- `internal/console/runtime_zero_model_e2e_test.go`: the console three-noop-task chain leg, the stop-on-natural-failure leg and the runnable-example preview leg.
- `internal/console/runtime_agent_concurrent_e2e_test.go`: the 8-way Start concurrency leg on one launcher / one runRoot, with its self-exec helper.
- 8 test cases (`internal/chain` 4 + `internal/console` 4) plus one self-exec helper `TestZeroModelAgentCLIHelper`.
- `examples/no-ai-preview/proofrail.chain.json`: the two-task pure-noop fixture the README promised and that was missing — a test-fixture fix outside the default run path.
- `examples/no-ai-preview/README.md`: minimal command-path fix (1 line).
- Doc backfill: `docs/ARCHITECTURE{,_EN}.md` (three-task E2E assembly path and concurrency assertion points), `docs/DEV_PLAN{,_EN}.md` (T027 milestone E2E status), `docs/t027/REMAINING_SLICES{,_EN}.md` (P45 summary row and detail block).

## Deviations and Registration

| ID | Content | Source | Ruling |
|---|---|---|---|
| D1 | ⑦ MINOR: the README line number and command snippet inside the failure message drifted | ⑦ release condition 1 | Closed: the message names `README.md:10` and the real command path (Finding Closure row 3) |
| D2 | The 2 MINOR of ⑤ and their disposition | ⑦ release condition 2 | Closed: the counter assertions became a declared text lock and the ledger evidence row became behavioural equivalence (Finding Closure rows 1 and 2) |
| D3 | The 4 uncovered areas of ⑥: real UI end-to-end, cross-OS and locale difference, final evidence-pack closure, real preview command execution evidence | ⑦ release condition 3 | Accepted and registered: the first three go to Known Limitations and Follow-Up Slice Candidates; the fourth is closed by the manual run in Gate Readings |
| D4 | The 4 unconfirmed items of ①: line numbers of `engine.go` and the `launcher` must be re-checked with `grep -n`; the `localRuntime` refusal wording would need the red test updated; `step-failed` is a deliberate contract; the preview leg depends on the example config staying put; plus the boundary "the console single-task limit and the A6 mirror collision are left alone" | ⑦ release condition 4 | Accepted and registered: the line numbers were re-checked during implementation; the other three are standing contracts and needed no change because nothing moved |
| D5 | This slice is assembly validation with zero production-code change | ⑦ release condition 5 | Declared: `git diff --name-only -- internal/` is empty and the only new files are 3 `_test.go` (Goal and Scope, item 2) |
| D6 | `examples/no-ai-preview` is a test-fixture fix, not part of the default run path | ⑦ release condition 6 | Declared: the fixture is pure noop with no execution side effect and is referenced only by the test and the README |
| D7 | `registry.len()==1` and slot-owner / mirror agreement are `internal/adapters` package-private state and are not asserted directly | ⑤ r1 2 / ⑥ Section C item 4 | Accepted: replaced by out-of-package behavioural equivalence (one winner, two loser sentinels, a live winner refusing a late start, release and re-start after stop); the cross-restart mirror is covered by the adapters' own tests |
| D8 | The five preview counters are static constants of `BuildPreviewReport` | ⑤ r1 1 / ⑥ Section C item 2 | Accepted: they are only a text lock on the README promise and must not be read as a model-channel sentinel; zero model use is carried by the four engine-leg assertions and the preview-leg derived values |
| D9 | The console single-task limit is kept; `runtime_agent.go:774` and `runtime_supervised.go:335` stay closed | ① section A / ⑦ release condition 4 | Accepted: multi-task assembly collides with A6's one identity mirror per runRoot, which DR13-FIX pinned; opening it would red the existing refusal cases and disturb the DR-9 sensitive surface |
| D10 | The ⑥ fence introduces 2 Chinese characters absent from the HEAD corpus; they are registered in the repository CJK new-word registry | ⑧ gates (G4b) | Accepted and registered: the transcript cannot be rewritten, so the entries carry the reason and context verbatim; this is the only write outside the doc-backfill list and it is declared as such |
| D11 | The package count of `go test ./...` disagrees with the kick-off brief | ⑧ measurement | Measured value wins: 15 test packages all `ok`, 0 `FAIL`, plus 4 packages without test files; this report states no unmeasured package count |

## Models and Metering

| Stage | Carrier | Model |
|---|---|---|
| ① Architecture | `deep-reasoner` | DeepSeek V4 Pro |
| ② Implementation | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ Tests | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ Pre-review | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ Scan | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ Final review | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ Docs | `prfrail-documenter` | DeepSeek V4.1 Flash |
| Batch assistance | `deep-reasoner` | GLM-5.3-Flash (2 dispatches: assembly-entry census / interface and port contract census) |

- **Main controller**: `DeepSeek V4.1 Flash` (High tier); ① and ⑤ stay on V4 Pro per the allow-list, which is a separate matter from the controller choice.

## Gate Readings

- `node tools/gates/gate.js --all --skip-base` (tree scope): `TOTAL_FAIL=0`; G3(a) / G3(b) PASS, G4b `unknown=[]`, G6-4 no hit, G7-a to G7-f all PASS, G8-a to G8-d all PASS.
- `gofmt -l internal/chain internal/console`: empty output.
- `go vet ./...`: no output.
- `go test -count=1 ./...`: 15 test packages all `ok`, 0 `FAIL`; 4 further packages carry no test files (`cmd/prfrail`, `cmd/prfrail-release`, `tools/agent-probe/supervisor`, `tools/b4-e2e`).
- `go test -count=1 -run 'TestZeroModel' ./internal/chain/ ./internal/console/ -v`: all 8 cases `PASS`, and the self-exec helper also `PASS`.
- Manual preview: `go run ./cmd/prfrail preview --chain examples/no-ai-preview/proofrail.chain.json --json` → `runnableInCli=true`, `preview.outcome=ready`, empty `blockingEvidence`, both steps `kind=noop`, all five counters 0, 3 `unknowns`.

## Commit and CI

The commit and CI readings are appended to this section once, after step ⑨, and the CN/EN pair is updated in the same edit.

## Follow-Up Slice Candidates

- Console multi-task assembly (if ever needed): lifting D9 first requires solving the one-identity-mirror-per-runRoot exclusion contract.
- Cross-OS and locale concurrency difference: a native measurement of `ProcessAlive` and snapshot semantics outside Windows.
- Extending the `-race` leg to `internal/console`: the local host lacks cgo, so the concurrency leg has no race-detector run.
- Real UI end-to-end: the in-editor preview and run paths are not yet aligned with this slice's tests.

## Known Limitations

- **⑥ uncovered 1**: no end-to-end click-through of the real UI or browser path; the verdict rests on static assertions over code and event logs.
- **⑥ uncovered 2**: no difference measurement across OS and locale outside Windows; concurrent snapshots and `ProcessAlive` may be platform sensitive, so the ubuntu CI leg is the final judge.
- **⑥ uncovered 3**: at the ⑧ point the evidence pack had not reached the final 06/07/08 payload stage; the final-review and CI closure is completed at ⑨.
- **⑥ uncovered 4**: the real preview command execution evidence is supplied by the manual run in Gate Readings, not by an automated test.
- **⑤ residual**: `registry.len()==1` and slot-owner / mirror agreement have no direct assertion (D7); the cross-restart mirror is covered by the adapters' own tests.
- **⑦ MINOR closed**: the failure-message drift is handled in Finding Closure row 3 and does not affect run semantics or gate strength.
- **D9 summary**: the console single-task limit stays closed; multi-task assembly and the A6 exclusion contract are mutually exclusive and sit outside this slice's boundary.

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/p45-e2e-zero-model.md	repo
docs/validation/p45-e2e-zero-model_EN.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/SHA256SUMS.txt	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/00-briefing.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/freedom-list.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/01-arch-design/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/03-tests/red-baseline.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/05-prereview/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/06-scan/review.md	repo
docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/07-final/review.md	repo
examples/no-ai-preview/README.md	repo
examples/no-ai-preview/proofrail.chain.json	repo
internal/chain/engine_zero_model_e2e_test.go	repo
internal/console/runtime_zero_model_e2e_test.go	repo
internal/console/runtime_agent_concurrent_e2e_test.go	repo
tools/gates/cjk-newwords.txt	repo
```

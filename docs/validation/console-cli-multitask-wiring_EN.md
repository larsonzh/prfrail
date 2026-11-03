# CONSOLE-CLI-MULTITASK-WIRING (P51) Acceptance Report

Slice: `CONSOLE-CLI-MULTITASK-WIRING` (implementation slice `[SLICE]`, tier M). Closing time 2026-11-03 (`HEAD` = `3aa53d5`, directive v1.60).

## 1. Goal and boundary

- Back-fill the P50 leftover: the command line could not declare per-task request records, so multi-task serial assembly could only be driven by the offline harness.
- This slice only parses CLI arguments and injects them into the assembly layer; it changes no assembly semantics and does not touch `docs/CONTRACTS{,_EN}.md`.
- The single-value `--request-record` single-task compatibility path is kept verbatim: `complete()` degrades to the old equivalence and `anyProvided()` only prepends a branch.

## 2. Design and interface

- A repeatable flag `--request-records <path>` is added and is strictly mutually exclusive with the single-value `--request-record`, rejected with exit 2 when both are given.
- Records are decoded by `(taskID, stepID)` with the same key function as the assembly layer: `agentRunnerRequestRecordKey(taskID, stepID)`.
- The registry is injected through `AgentRunnerAssemblyOptions.RequestRecords`; the single-value path keeps the original single-record branch and leaves that field nil.
- Seven fail-closed cases: both flags, an empty path, a duplicate key, a key not bound by the chain definition, a definition coverage gap, a record `runId` differing from `--run-id`, and multiple records declared on a noop-only chain.
- Every refusal precedes run-root creation and event-log writes: the earlier IO is read-only and `writeCommandError` writes only to stdout/stderr, so the refusal is zero-write.

## 3. Implementation changes

- `internal/console/cli.go`: +132 / -13 (flag registration, mutual exclusion and the noop-only guard, record reading with deterministic key order, injection into the assembly options).
- New `internal/console/cli_multitask_wiring_test.go`: a positive read, six refusals and one single-value compatibility case (the case names are listed in the review chain and the evidence pack).
- `internal/console/runtime_agent.go`, `internal/console/runtime_supervised.go` and `cmd/prfrail/main.go` are unchanged.
- `docs/CONTRACTS{,_EN}.md` is unchanged: the protocol is untouched, so no revision is required.

## 4. Red-green reversal evidence

- Stage 3 baseline: four red and one green among the five skeleton tests, all four reds failing with `unknown flag` because the flag was not registered (`03-tests/logs/03-tdd-red-run.txt`).
- After stages 2a/2b the five were all green; stage 2c added the noop-only chain leg, which first went red (old code exited 0 silently) and then green.
- The baseline list is in `03-tests/red-baseline.md` and the step-by-step implementation record is in `04-impl/implementer-reports.md`.
- The stage 5 F1 fix went red first: `TestRunRefusesRequestRecordsWithoutExecutableSteps` asserts `FindExecutableSteps==0` upfront, so the old code was bound to take the noop path and exit 0.

## 5. Review chain

- 1 Architecture: DESIGN DONE (`01-arch/design.md`).
- 5 Pre-review r1: PASS WITH FIXES (1 High / 2 Low / 3 Info, F1 being the noop-only chain silent degradation); r2: PASS (`05-prereview/review.md`, `05-prereview/review-r2.md`).
- 6 Independent scan: INDEPENDENT SCAN: PASS (`06-scan/review.md`).
- 7 Final review: FINAL REVIEW: PASS (1 Info: no positive wiring leg; `07-final/review.md`).

## 6. Gate readings

- `gofmt -l internal/console` is silent; `go build ./...` and `go vet ./...` pass.
- `go test -count=1 ./...`: all fifteen packages are `ok` with no FAIL (`04-impl/logs/04-full-suite-green.txt`).
- The three-scope `TOTAL_FAIL=0` and `SELFTEST 697/697` are in `08-gates/gate-readings.txt`.
- Commit `2b0af2d` was pushed to origin (never gitee) and run `41569353460` is green on the first attempt on both legs: the `Gates` step reports `SELFTEST 697/697` and `TOTAL_FAIL=0` on each leg.

## 7. Open items and boundary

- 7 Info "no positive wiring leg": this slice adds no cross-process positive assertion; wiring correctness is covered by the code-path verification and the refusal legs.
- The supervised assembly has no record input, so `--request-records` is always refused for a supervised chain; its record injection is out of scope here.
- The real R/C write paths on Windows are carried by the CI legs; locally this slice covers only the CLI layer and the unit-test regression.
- The `...` ambiguity in the `runUsage` text `(--request-record <path> | --request-records <path>)...` is given verbatim by the design and is not an implementation defect.
- Deviation D2: while counting receipt characters a scratch statistics file was written to the user temp directory (`%TEMP%`) outside the repository, violating the v40 section 8 temp-artefact location closed set; the file is deleted and the ledger D count is corrected to 2.

## 8. artifacts fence

The archived entries for this slice are listed below (path plus owner, TAB separated, ASCII dictionary order).

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/CHECKLIST_v41.md	repo
docs/t027/CHECKLIST_v41_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/console-cli-multitask-wiring.md	repo
docs/validation/console-cli-multitask-wiring_EN.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/00-briefing.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/01-arch/design.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/03-tests/logs/03-tdd-red-run.txt	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/03-tests/red-baseline.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/04-impl/implementer-reports.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/04-impl/logs/04-full-suite-green.txt	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/05-prereview/review-r2.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/05-prereview/review.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/06-scan/review.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/07-final/review.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/08-gates/gate-readings.txt	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/09-glm/glm-1-cli-flag-inventory.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/09-glm/glm-2-cli-test-enumeration.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/MANIFEST.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/SHA256SUMS.txt	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/freedom-list.md	repo
docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/input-list.md	repo
internal/console/cli.go	repo
internal/console/cli_multitask_wiring_test.go	repo
```

## 9. Verbatim stage 6 independent scan

The verbatim text of the stage 6 independent scan conclusion follows, with no commentary outside the fence.

```sixstage
## Section A
范围：只读审查 diff 与 CLI/测试；未修改文件/未做 git 写操作。
证据：git diff HEAD -- internal/console/cli.go internal/console/cli_multitask_wiring_test.go + build/vet/test。

## Section B
(no findings)

## Section C
internal/console/cli.go:319-346: --request-records 解析、与 --request-record 互斥、noop-only 失败闭环。
internal/console/cli.go:398-439: complete()/anyProvided() 扩展 repeatable 形式且保留单值兼容。
internal/console/cli.go:467-547: executeAgentRunnerRun() 把 RequestRecords 传入 AgentRunnerAssemblyOptions。
internal/console/runtime_agent.go:81-88,1078-1093: 运行时按 <taskID>/<stepID> 取 map，并拒绝缺映射。
go test -count=1 ./internal/console/: ok github.com/larsonzh/prfrail/internal/console 10.148s

## Section D
静默降级：executeRun()/readAgentRunnerRequestRecords() 在 no-op chain 上 fail closed，且 assembly 读取 map by key。
7个拒绝面：双flag/重复键/错 runId/键未命中/覆盖缺失/空路径/noop-only 均可命中，且错误文案带关键字。
complete()/anyProvided()：保留旧单值语义，新增 repeatable 分支为加法，不改变单值行为。
误绿风险低：测试断言 stderr 关键字，不只看 exitCode；装配消费点已核实。

## Section E
无法验证：未真实启动 AgentRunner 进程/多任务串行执行，仅做 CLI 层和编译/单测回归。
推翻条件：若 --request-records 在 noop chain 或错键场景 exit 0 且继续；或 single-value path 失去 request record: 文案；或 map key 读错。

INDEPENDENT SCAN: PASS
```

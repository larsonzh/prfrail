# CONSOLE-CLI-MULTITASK-WIRING（P51）验收报告

切片：`CONSOLE-CLI-MULTITASK-WIRING`（实现类 `[SLICE]`，M 档）。封片时点 2026-11-03（`HEAD` = `3aa53d5`，准则 v1.60）。

## 1. 目标与边界

- 补齐 P50 遗留：命令行此前无法逐任务声明请求记录，多任务串行装配只能由离线 harness 驱动。
- 本片只做 CLI 参数解析与向装配层的注入，不改装配语义，不改 `docs/CONTRACTS{,_EN}.md`。
- 单值 `--request-record` 的单任务兼容路径逐字保留：`complete()` 退化等价，`anyProvided()` 只做前置追加。

## 2. 设计与接口

- 新增可重复 flag `--request-records <path>`，与单值 `--request-record` 严格互斥，同时给出以 exit 2 拒绝。
- 记录按 `(taskID, stepID)` 解码，键函数与装配层同源：`agentRunnerRequestRecordKey(taskID, stepID)`。
- 注册表经 `AgentRunnerAssemblyOptions.RequestRecords` 注入；单值路径仍走原单记录分支，该字段保持 nil。
- 七条 fail-closed：双 flag、空路径、重复键、键未命中链定义、定义缺覆盖、记录 `runId` 与 `--run-id` 不一致、noop-only 链上声明多记录。
- 全部拒绝早于 run-root 创建与事件日志写入：此前的 IO 只读，`writeCommandError` 只写 stdout/stderr，故为零写入。

## 3. 实现改动

- `internal/console/cli.go`：+132 / −13（flag 注册、互斥与 noop-only 守卫、记录读取与键序确定化、注入装配选项）。
- 新增 `internal/console/cli_multitask_wiring_test.go`：含正向读取、六条拒绝与一条单值兼容用例（逐条名称见第三节审查链与证据包）。
- `internal/console/runtime_agent.go`、`internal/console/runtime_supervised.go`、`cmd/prfrail/main.go` 零改动。
- `docs/CONTRACTS{,_EN}.md` 零改动：协议未变，无需修订。

## 4. 红绿反转证据

- ③ 基线：五条骨架测试四点红一绿，四红均因 flag 未注册而报 `unknown flag`（`03-tests/logs/03-tdd-red-run.txt`）。
- ②a/②b 后五条全绿；②c 补 noop-only 链腿，先红（旧代码 exit 0 静默执行）后绿。
- 基线清单见 `03-tests/red-baseline.md`，分步实施记录见 `04-impl/implementer-reports.md`。
- ⑤ F1 的修复先行落红：`TestRunRefusesRequestRecordsWithoutExecutableSteps` 前置断言 `FindExecutableSteps==0`，旧代码必走 noop 路径并 exit 0。

## 5. 审查链

- ① 架构：DESIGN DONE（`01-arch/design.md`）。
- ⑤ 预审 r1：PASS WITH FIXES（1 High / 2 Low / 3 Info，F1 为 noop-only 链静默降级）；r2：PASS（`05-prereview/review.md`、`05-prereview/review-r2.md`）。
- ⑥ 独立扫描：INDEPENDENT SCAN: PASS（`06-scan/review.md`）。
- ⑦ 终审：FINAL REVIEW: PASS（1 Info：缺正向接线腿；`07-final/review.md`）。

## 6. 门禁读数

- `gofmt -l internal/console` 无输出；`go build ./...`、`go vet ./...` 通过。
- `go test -count=1 ./...`：十五个包全 `ok`，无 FAIL（`04-impl/logs/04-full-suite-green.txt`）。
- 三作用域 `TOTAL_FAIL=0` 与 `SELFTEST 697/697` 见 `08-gates/gate-readings.txt`。
- 提交 `2b0af2d` 已推 origin（未推 gitee），run `41569353460` 首跑双腿全绿：两腿 `Gates` 均为 `SELFTEST 697/697`、`TOTAL_FAIL=0`。

## 7. 未闭环项与边界

- ⑦ Info「无正向接线腿」：本片未新增跨进程正向断言，接线正确性由代码路径核验与拒绝腿共同覆盖。
- supervised 装配无记录输入端，故 `--request-records` 对 supervised 链一律拒绝；其记录注入不在本片范围。
- Windows 上真实 R/C 写路径由 CI 腿承担，本片本地只做 CLI 层与单测回归。
- `runUsage` 文案 `(--request-record <path> | --request-records <path>)...` 的 `...` 歧义为设计逐字给定，非实施缺陷。
- 偏离登记 D2：回执统计时曾把中间统计文件写到仓库外用户临时目录（`%TEMP%`），违反 v40 §八「临时件位置」闭集；该文件已删除，台账 D 计数修正为 2。

## 8. artifacts 围栏

以下为本片归档条目清单（路径 + 归属，TAB 分隔，ASCII 字典序）。

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

## 9. ⑥ 独立扫描原文

以下为 ⑥ 独立扫描结论的逐字原文，围栏外不追加评价。

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

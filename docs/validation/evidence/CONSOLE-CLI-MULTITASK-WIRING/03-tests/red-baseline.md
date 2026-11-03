# ③ 红测报告原文 —— P51 `CONSOLE-CLI-MULTITASK-WIRING`

载体：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（TDD 先行）。实得输出字符数 3890（自报 `OUTPUT_CHARS=3890`）。

主控机械复核：新增文件 `internal/console/cli_multitask_wiring_test.go` 存在、`go vet ./internal/console/` 无输出、红基线原始输出已由主控复跑并落盘 `03-tests/logs/03-tdd-red-run.txt`（结论见报告 §红基线实测 与主控复跑一致：4 红 1 绿）。

---

## 1. 新增文件与测试清单（文件:行号 + 每个测试一句断言摘要）

新增（唯一写入）：`internal/console/cli_multitask_wiring_test.go`（`package console`，gofmt 已执行）。生产代码零改动。

| 测试 | 行号 | 断言摘要 |
|---|---|---|
| `TestRunRequestRecordsPathsAreRead` | 157 | 两 task 链 + 七件套 + `--request-records <不存在路径>`：exit≠0、stderr 含 `request record`、且不含 `flag provided but not defined` |
| `TestRunRefusesRequestRecordWithRequestRecords` | 178 | 单 task 链 + `--request-record a.json` 与 `--request-records b.json` 同时出现：exit≠0、stderr 含 `mutually exclusive` |
| `TestRunRefusesDuplicateRequestRecordKey` | 199 | 两 task 链 + 两条键同为 `task-one/code-one`（requestId 不同）的有效记录：exit≠0、stderr 含 `duplicate key` |
| `TestRunRefusesRequestRecordRunIDMismatch` | 219 | 两 task 链 + 两条各命中一 task 的有效记录，其一 `RunID=run-p51-other` ≠ `--run-id`：exit≠0、stderr 含 `does not match run id` |
| `TestRunSingleRequestRecordStaysCompatible` | 240 | 单 task 链 + 完整八件套 + 仅单值 `--request-record` 指向不存在文件：exit≠0、stderr 含 `request record:`（带冒号） |

同文件内局部 helper（全部 `p51` 前缀，无重名）：`p51Digest`(41)、`p51TwoTaskChain`(53)、`p51SingleTaskChain`(61)、`p51WriteChainConfig`(67，先 `ValidateChainConfig` 再落盘)、`p51WriteRequestRecord`(95，`adapters.NewAgentRunnerRequestRecord` + `json.Marshal`)、`p51SevenAgentRunnerArgs`(134)、`p51RunArgs`(149)。

## 2. 红基线实测

命令：`go test -count=1 -run '<5 个测试>' -v ./internal/console/ 2>&1 | Tee-Object -FilePath tmp\CONSOLE-CLI-MULTITASK-WIRING\red-baseline.txt`
落盘：`tmp/CONSOLE-CLI-MULTITASK-WIRING/red-baseline.txt`（5530 字节，Tee-Object 可用）。

| 测试 | 状态 | 关键输出原文（首行） |
|---|---|---|
| `TestRunRequestRecordsPathsAreRead` | FAIL | `cli_multitask_wiring_test.go:172: --request-records must be a defined flag, stderr="run: flag provided but not defined: -request-records` |
| `TestRunRefusesRequestRecordWithRequestRecords` | FAIL | `cli_multitask_wiring_test.go:192: the refusal must state the flags are mutually exclusive, stderr="run: flag provided but not defined: -request-records` |
| `TestRunRefusesDuplicateRequestRecordKey` | FAIL | `cli_multitask_wiring_test.go:213: the refusal must name the duplicate key, stderr="run: flag provided but not defined: -request-records` |
| `TestRunRefusesRequestRecordRunIDMismatch` | FAIL | `cli_multitask_wiring_test.go:233: the refusal must state the run id mismatch, stderr="run: flag provided but not defined: -request-records` |
| `TestRunSingleRequestRecordStaysCompatible` | PASS | `--- PASS: TestRunSingleRequestRecordStaysCompatible (0.01s)` |

尾部汇总原文：`FAIL` / `FAIL	github.com/larsonzh/prfrail/internal/console	1.092s` / `FAIL`（exit code 1）。
说明：4 条红腿失败原因一致——HEAD `3aa53d5` 未定义 `--request-records`，flag 解析层即失败（exit 2，stderr `run: flag provided but not defined: -request-records` + `Usage of run:` 全量 flag 清单 + `usage: prfrail run ...`），未进入任何装配/记录读取逻辑；这正是本片要消除的行为。

## 3. 绿腿证据

```
=== RUN   TestRunSingleRequestRecordStaysCompatible
--- PASS: TestRunSingleRequestRecordStaysCompatible (0.01s)
```

该腿证明单值 `--request-record` 现有路径仍走到 `executeAgentRunnerRun` 的记录读取分支（stderr 标签含冒号 `request record:`），接线改造不得破坏它。

## 4. 编译与格式检查结果

- `gofmt -l` 首跑列出该文件 → 执行 `gofmt -w internal\console\cli_multitask_wiring_test.go` → 复跑 `gofmt -l` 无输出（已格式化）。
- `go vet ./internal/console/`：无任何输出（通过）。
- `go test -run` 编译阶段通过，未出现任何编译错误或未定义符号。

## 5. 偏差与风险

1. **第 3/4 条测试的未来行为依赖 ② 的实现细节**：`duplicate key`、`does not match run id` 两段文案由本片规格逐字固定；若 ② 采用不同措辞，这两条会保持红色。属规格约束，非测试缺陷，需 ② 按此文案实现。
2. **第 1 条测试的“七件套”未做落盘**：仅提供 7 个不存在的路径，用于保证互斥/读取判据先于其他记录读取触发；若 ② 实现顺序改为先读其它记录，stderr 首行会变成 `capability record:`，但因断言仅要求含 `request record`，仍为绿——即该断言对“读顺序”不敏感（有意放宽以免过度锁定实现）。
3. **记录文件使用固定时钟常量** `2026-11-03T00:00:00.000Z`，无真实时钟断言；无网络、无外部二进制、无 TestMain 依赖。
4. **未覆盖**：`--request-records` 与 supervised 定义互斥、以及“records 注册表非空时忽略单记录”的 CLI 侧语义不在本片 5 条清单内（规格未要求），归 ③ 后续或 ② 单测范围。
5. 未执行 `git add/commit/push`；工作区改动仅新增该测试文件（另有步骤 0 既有台账/证据包文件），`tmp/CONSOLE-CLI-MULTITASK-WIRING/` 已 gitignore。

TEST: DONE
OUTPUT_CHARS=3890

# T016 CLI/Console 验证报告

[English](t016-cli-console_EN.md)

日期：2026-09-08。结论：`T016 COMPLETE`。本任务交付无 IDE 可用的 CLI 基线：`init`、`validate`、`config explain`、`run`、`report`，并将占位入口替换为可测试执行器。

## 实现范围

1. `cmd/prfrail/main.go`
   - 入口改为调用 `internal/console` 命令执行器，不再使用占位 `init/run` 逻辑。
2. `internal/console/config.go`
   - 新增链配置模型、严格 JSON 解析（未知字段拒绝）、语义校验。
   - 新增配置搜索顺序：`--chain` → `./proofrail.chain.json` → `./proofrail.json`。
   - 新增 `config explain` 解析结果与来源说明（含 builtin-default/chain/task 来源）。
3. `internal/console/cli.go`
   - 新增命令路由与参数解析：`version/init/validate/run/report/config explain`。
   - 新增稳定 `--json` 输出与错误退出码（0/1/2）。
4. `internal/console/runtime.go`
   - 新增基于 `internal/chain` 的 noop-only 本地执行与 `report` 汇总读取。
   - 对可执行 step（`code/build/verify`）明确 fail-close 并返回非零退出。
5. `internal/chain/models.go`
   - 导出 `ValidateDefinition`，供 CLI 复用核心语义校验。

## AT-12 覆盖

1. 无 IDE CLI 闭环：可直接通过命令行执行 `init/validate/config explain/run/report`。
2. 稳定可解析输出：所有核心命令支持 `--json`，失败场景同样返回结构化 JSON。
3. 无色输出：文本模式不输出 ANSI 转义序列。
4. 非实现能力不冒充成功：`run` 对 `code/build/verify` step 明确失败并非零退出。
5. 配置解释可追溯：`config explain` 输出有效策略与来源指针。

## 测试用例

- `internal/console/cli_test.go`
  - `TestInitValidateAndExplainJSON`
  - `TestRunAndReportNoopChain`
  - `TestRunRejectsExecutableSteps`
  - `TestValidateRejectsUnknownField`
  - `TestTextOutputHasNoANSI`
  - `TestUnknownCommandUsageExitCode`

## 门禁

1. `go test -count=1 ./...`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `2026-09-08` 远程 race 高次数回归（Ubuntu VM `10.0.0.199`，Go 1.27.1 + gcc 13.3.0）：`go test -race -count=20 ./internal/applier ./internal/console ./internal/taskdef` 与 `go test -race -count=1 ./...` 通过（首轮通过，未触发代理兜底）。

## 边界

1. 当前 `run` 仅支持 noop-only 任务链；执行型 step 的真实 gate/adapter 闭环属于后续切片。
2. `proofrail.toml/workspace.toml` 分层合并仍是规划项，本轮固定单文件链配置解析与 explain。
3. 本地 `run` 的 review/promotion 为内存合成决定，回执不落盘，仅用于演示 noop 闭环，不构成验收证据。
4. 本轮默认未提交、未推送。

# RACE-LEG-EXPAND · ② 实现角色报告

- 切片：P49 `RACE-LEG-EXPAND`
- 角色：② 实现者（prfrail-implementer）
- 日期：2026-11-02
- 授权边界：只改 `.github/workflows/ci.yml` 的 Race 步骤 `run:` 行；不改任何 `.go` 源码；不改任何既有步骤；产出仅落 `tmp/RACE-LEG-EXPAND/` 与本目录

## 1. 改动点清单（方案书 → 落地）

| 编号 | 改动点 | 文件 | 类型 | 落地状态 |
|---|---|---|---|---|
| CP-1 | Race 腿目标包由 2 包扩容至 6 包 | `.github/workflows/ci.yml`（步骤 `Race` 的 `run:` 行，第 84 行） | 修改（单行，1 insertion / 1 deletion） | **已完成** |

未覆盖改动点：无。

## 2. 修改前后逐字对照

修改前（`ci.yml` 第 84 行，逐字）：

```
        run: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...
```

修改后（`ci.yml` 第 84 行，逐字）：

```
        run: CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...
```

形状保持：`CGO_ENABLED=1 ` 前缀保留；`-race`、`-count=1` 保留；六包顺序为 adapters → chain → console → gates → guard → tickets；每包 `./internal/<pkg>/...` 形态各出现 1 次（由 ③ 的断言脚本逐包计数确认）。

## 3. 四项自验结论

| # | 自验项 | 命令 | 结果 | 证据（磁盘路径） |
|---|---|---|---|---|
| 1 | diff 只有一行改动 | `git diff -- .github/workflows/ci.yml` | **PASS**：`--numstat` = `1  1`；`--stat` = `1 file changed, 1 insertion(+), 1 deletion(-)` | `tmp/RACE-LEG-EXPAND/04-ciyml.diff` |
| 2 | 红→绿翻转复跑 | `node tmp\RACE-LEG-EXPAND\assert-race-leg-coverage.js` | **GREEN**，`EXIT=0`（末行 `RACE-COVERAGE: GREEN - all six packages covered by the Race leg`） | `tmp/RACE-LEG-EXPAND/03-tdd-green-run.txt`（红基线 `03-tdd-red-run.txt` 未被覆盖） |
| 3 | YAML 语法检查 | `go run .vscode\scripts\yamlcheck.go .github\workflows\ci.yml` | **PASS**：`PASS .github\workflows\ci.yml (yaml parse ok, 16432 bytes)`，exit 0 | 本报告第 3 节（命令输出已记录，见下） |
| 4 | 编码 / 行尾 | `node` 读原始字节 | **PASS**：`BOM=false`、`CRLF_PAIR=false`、`LF_TOTAL=408`、`LF_ONLY=408`、`BYTES=16432`、`SHA256=627d71cc893f9414ab33d010b3af05d105db9a2772cbad05d43608f6208dc264` | 本报告第 3 节（同上） |

第 2 项断言脚本逐项结论（摘自 `03-tdd-green-run.txt`）：

```
[CHECK] CGO_ENABLED=1 prefix: PASS - "CGO_ENABLED=1 " present at command head
[CHECK] go test -race: PASS
[CHECK] -count=1: PASS
[CHECK] Linux-only guard preserved: PASS
[CHECK] package covered: ./internal/adapters/...: PASS - occurrences=1
[CHECK] package covered: ./internal/chain/...: PASS - occurrences=1
[CHECK] package covered: ./internal/console/...: PASS - occurrences=1
[CHECK] package covered: ./internal/gates/...: PASS - occurrences=1
[CHECK] package covered: ./internal/guard/...: PASS - occurrences=1
[CHECK] package covered: ./internal/tickets/...: PASS - occurrences=1
[CHECK] package order matches documented order: PASS
```

第 3、4 项原始输出：

```
PASS .github\workflows\ci.yml (yaml parse ok, 16432 bytes)
YAMLEXIT=0
BOM=false
CRLF_PAIR=false
LF_TOTAL=408 LF_ONLY=408
BYTES=16432
SHA256=627d71cc893f9414ab33d010b3af05d105db9a2772cbad05d43608f6208dc264
```

## 4. 未做事项声明

- 未修改任何 `.go` 生产代码或测试源码。
- 未删除、未重排、未改写着任何既有 CI 步骤：`- name: Race` 与 `if: runner.os == 'Linux'` 原样保留，`Contract fixtures`、`Gates`、`Build`、`Vet`、`Test`、`Setup Go`、`checkout` 等均未触碰。
- 未修改 workflow 顶层 `env:` 或任何其它行（由 `--numstat = 1  1` 佐证）。
- 未改任何 `.md` 正文（本报告为本角色唯一的 `.md` 产出）。
- 未新增/升级任何依赖。
- 未执行 `git add` / `git commit` / `git push`。

## 5. 已知限制

1. 本地无法实跑 `-race`：本机缺 cgo/gcc 工具链，`CGO_ENABLED=1 go test -race` 构造期即 `build failed`（取证：`tmp/RACE-LEG-EXPAND/03-race-local-infeasible.txt`）。因此本片**不**以本机 race 结果作为证据。
2. `-race` 真实可运行性由两处承担：① VM 上六包 `-race` 已绿（`docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt`）；② 未来 Linux CI 腿（`runner.os == 'Linux'`，`CGO_ENABLED=1`）的实跑结论。本片只保证「CI 腿的目标包清单已扩容」这一静态事实。
3. CI 墙钟时间会增加：Race 腿从 2 包扩到 6 包，新增 console/gates/guard/tickets 四包的 race 编译与执行开销（③ 与 ⑨ 的分析有相应估计）。本角色不做时间预算判断。
4. 编辑器脏缓冲风险已排除取证：第 4 项用 Node 直接读磁盘字节（`BYTES=16432` 与 `yamlcheck` 报告一致，SHA256 已登记），且 `git diff` 可见该改动，故磁盘为真相、内存缓冲未覆盖磁盘。

IMPLEMENT: DONE

---

## r2：工作流冻结契约同步

- 轮次：r2（主控发现 r1 改动触发工作流冻结契约）
- 本 r2 修改文件（唯一）：`internal/release/workflow_test.go`（1 行）

### r2-1 契约更新依据

机制：`validateWorkflowScripts` 的 `expectedScripts` 映射逐字锁定每个脚本步骤 `step.Run` 的 SHA-256；
`Race` 步骤体在 r1 扩容为六包后摘要变更，期望值仍为旧值，故整个 `internal/release` 工作流契约测试族红。

改前实测（`go test -count=1 ./internal/release/`，exit 1，原文摘录，完整见
`tmp/RACE-LEG-EXPAND/04-red-before-contract-sync.txt`）：

```
--- FAIL: TestWorkflowScriptMutationsAreRejected (0.01s)
    workflow_mutation_test.go:43: baseline workflow: script step "test/Race" differs from the reviewed body: SHA-256=825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9, want ff4826f2103a2457cd306662f2542e54d809546deb1f804964b63f2148185dd5
--- FAIL: TestWorkflowActionParametersAreLocked (0.01s)
--- FAIL: TestWorkflowStructureMutationsAreRejected (0.00s)
--- FAIL: TestWorkflowYAMLMutationsAreRejected (0.00s)
--- FAIL: TestSelfhostWorkflowUsesIndependentTrustedVerifier (0.01s)
--- FAIL: TestWorkflowActionMutationsAreRejected (0.01s)
--- FAIL: TestWorkflowScriptsRejectCommentedCommands (0.01s)
FAIL    github.com/larsonzh/prfrail/internal/release    6.528s
```

共 **7** 个测试红，错误原因完全一致（同一 `test/Race` 摘要不符）。失败测试自行算出的实际摘要
`825bacc3...` 即为新 run 行的 SHA-256。

改动（`internal/release/workflow_test.go:245`，唯一一处）：

```diff
-			digest: "ff4826f2103a2457cd306662f2542e54d809546deb1f804964b63f2148185dd5", condition: "runner.os == 'Linux'",
+			digest: "825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9", condition: "runner.os == 'Linux'",
```

### r2-2 新摘要核对过程

口径：`step.Run` 的内容字节（即 `run:` 之后的纯脚本体，不含行首缩进、不含行尾换行）。

取数方式：Node 读取 `d:\LZProjects\prfrail\.github\workflows\ci.yml` 原文，定位第 84 行
（正则 `^\s*run:\s*CGO_ENABLED=1 go test -race`），按 `^\s*run:\s?` 剥除前缀后对字符串做 SHA-256：

```
LINE_NO=84
BODY_LEN=169
SHA256=825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9
BODY=CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...
```

核对结论：计算值 `825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9` 与待写入期望值**逐字符一致**，
且与失败测试输出的 `SHA-256=825bacc3...` 一致 ⇒ 契约摘要更新正确，非猜测值。

### r2-3 自验结论（三项 + 逐条工作流测试）

| # | 自验项 | 命令 | 结果 | 证据路径 |
|---|---|---|---|---|
| 1 | 契约 diff 仅 1 行 | `git diff -- internal/release/workflow_test.go` | **PASS**：`--numstat` = `1  1`；diff 正文仅 1 行 `-/+`（digest 值） | `tmp/RACE-LEG-EXPAND/04-contract-diff.txt` |
| 2 | `internal/release` 全绿 | `go test -count=1 ./internal/release/` | **PASS**：`ok github.com/larsonzh/prfrail/internal/release 7.036s`，EXIT=0 | `tmp/RACE-LEG-EXPAND/04-release-green.txt` |
| 3 | 全仓全绿 | `go test -count=1 ./...` | **PASS**：`OK_COUNT=15`、`FAIL_COUNT=0`，EXIT=0（与改前口径一致） | `tmp/RACE-LEG-EXPAND/04-gotest-all-green.txt` |
| 4 | 新摘要核对（Node/SHA-256） | 见 r2-2 | **PASS**：`825bacc3...` 一致 | 本报告 r2-2 代码块 |
| 5 | 7 个工作流测试逐条 | `go test -count=1 -run 'TestWorkflow\|TestSelfhost' -v ./internal/release/` | **PASS**：`WF_PASS=7`、`WF_FAIL=0`，EXIT=0 | `tmp/RACE-LEG-EXPAND/04-release-workflow-tests-verbose.txt` |

第 5 项逐条结果：`TestWorkflowScriptMutationsAreRejected`、`TestWorkflowActionParametersAreLocked`、
`TestWorkflowStructureMutationsAreRejected`、`TestWorkflowYAMLMutationsAreRejected`、
`TestSelfhostWorkflowUsesIndependentTrustedVerifier`、`TestWorkflowActionMutationsAreRejected`、
`TestWorkflowScriptsRejectCommentedCommands` —— 全部 PASS。

### r2-4 未做事项声明

- 未改 `workflow_test.go` 其它任何行（`--numstat = 1  1` 佐证）。
- 未改 `.github/workflows/ci.yml`（r1 成果保持；本 r2 未触碰）。
- 未改任何其它 `.go` / `.md` / Schema / 依赖（本报告为本轮唯一 `.md` 产出）。
- 未动生产代码；未执行 `git add` / `git commit` / `git push`。
- 临时件全部落 `tmp/RACE-LEG-EXPAND/`（只增不删）。

IMPLEMENT: DONE

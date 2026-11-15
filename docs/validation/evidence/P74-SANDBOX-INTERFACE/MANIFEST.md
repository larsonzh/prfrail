# MANIFEST — P74-SANDBOX-INTERFACE 证据包

切片：`SANDBOX-INTERFACE`（P74，层 1 主线 B 段引入独立片，实现类 `[SLICE]`，M 档）。建包时点 2026-11-15（`HEAD` = `f834ce5`，准则 v1.69）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `03-tests/` / `04-reviews/` / `06-gates/` 载荷；⑧ 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10；含 2026-11-14 同轮裁定授权） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `01-arch/design-r1.md` | ① 设计原文（deep-reasoner + DeepSeek V4 Pro；含 r2 补正，`DESIGN: DONE` + `DESIGN-R2: DONE`） |
| `02-batch/batch-01-interface-enumeration.md` | 批量辅助 ①（GLM-5.3-Flash）guard / chain / adapters / permit 事实枚举原文归档（实得 7178 字符） |
| `02-batch/batch-02-contracts-structure.md` | 批量辅助 ②（GLM-5.3-Flash）CONTRACTS §7 结构与台账 P74 段事实枚举原文归档（实得 3016 字符） |
| `03-tests/red-guard-contract.txt` | ③ 红测（`TestSandboxInterfaceContract`，惰性骨架期 failure） |
| `03-tests/red-agent-runner.txt` | ③ 红测（`TestAgentRunnerRunSandboxRefusesBeforePermit` 三子例断言红） |
| `03-tests/red-supervised.txt` | ③ 红测（`TestSupervisedRunSandboxRefusesBeforePermit` 三子例断言红） |
| `03-tests/green-guard-contract.txt` | ②b 后绿测（guard 契约 + 门函数四分支） |
| `03-tests/green-agent-runner.txt` | ②b 后绿测（AgentRunner 调用点级） |
| `03-tests/green-supervised.txt` | ②b 后绿测（supervised 调用点级） |
| `03-tests/mutation-agent-runner.txt` | 变异检验（删 `runtime_agent.go` 检查块后变红） |
| `03-tests/mutation-supervised.txt` | 变异检验（删 `runtime_supervised.go` 检查块后变红） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（deep-reasoner + DeepSeek V4 Pro；PASS，3 Low + 3 Info） |
| `04-reviews/prereview-r2.md` | ⑤ 预审 r2 复跑原文（三条 Low 整改后；PASS，1 Low + 2 Info） |
| `04-reviews/scan-r1.md` | ⑥ 独立扫描原文（independent-reviewer + MAI-Code-1.1-Flash；PASS，Section A–E） |
| `04-reviews/final-r1.md` | ⑦ 终审原文（independent-reviewer + GPT-5.3-Codex；PASS，盲审不抽） |
| `06-gates/gate-index.txt` | ④ `gate.js --all --scope=index` 原始输出 |
| `06-gates/gate-tree.txt` | ④ `gate.js --all --scope=tree` 原始输出 |
| `06-gates/gate-ci.txt` | ④ `gate.js --all --scope=ci` 原始输出 |
| `06-gates/gate-selftest.txt` | ④ `gate.js --selftest` 原始输出（SELFTEST 737/737） |
| `06-gates/linux-build.txt` | ④ `GOOS=linux GOARCH=amd64` 的 build / vet 读数（含显式退出码） |
| `06-gates/full-test.txt` | ④ `go test ./... -count=1` 全量回归原始输出（ok 15 / FAIL 0） |

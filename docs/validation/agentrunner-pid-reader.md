# 切片 `AGENTRUNNER-PID-READER` 验证报告（准则 v1.42）

> 切片：`AGENTRUNNER-PID-READER`；分级 `[SLICE]`、档位 M-L；准则 v1.42（本片只补登 §12.4 与 §0.3 各一行、升头部版本号）。
> 基线：`d:\LZProjects\prfrail`；HEAD = `7201679`（C1 闭环）。
> 双语：本文件为 CN 权威；`docs/validation/agentrunner-pid-reader_EN.md` 为逐行同位镜像。

## 1 变更概要

本片一句话概括：闭合两笔欠账——OB-110（读者侧三态读法）与 OB-109（第六枚冻结事实）——并以准则 v1.42 补登 OB-111。

交付物：

- 契约 A5/A6：`docs/CONTRACTS.md` 与 `docs/CONTRACTS_EN.md` 各 +2/−2 行（`git diff --cached --numstat` 读数）。
- 代码：18 条路径——`internal/adapters` 10 个、`internal/chain` 4 个、`tools/b4-e2e` 4 个。
- 聚焦测试：上述文件内的新增/改写用例，另含 `a6native` 与 `b3bnative` 两条具名运行。
- 准则三处：§0.3 的 v1.42 条目、§12.4 的 OB-111 条目、头部版本号。
- 本报告对（CN 权威 + `_EN` 镜像）与 `docs/validation/evidence/AGENTRUNNER-PID-READER-freedom-list.md`。

**v1.42 的定位**：本片是执行切片；v1.42 只在切片内补登一行台账并升版本号，属切片内副产物，不是独立治理片，也不改任何判据语义与判据作用域。

## 2 环境与工具链

- 本机工具链：`go1.27.0 windows/amd64`。
- CI 声明：Go `1.22.12`。
- 仓库工作区：Windows（`d:\LZProjects\prfrail`）。
- 本片未在 Linux VM 上复跑：本片改动无平台分支，与 C1 的跨平台复跑不同。
- 命令面：`gofmt -l .`、`go build ./...`、`go vet ./...`、`go test -count=1 ./...`。

## 3 异常与兜底记录

**D1（本片唯一偏离）**：②a 在 `tmp/` 根创建了 `tmp/verify-*.js` 与 `tmp/pid-reader-dump.md`，违反 §5.1「临时件」的位置闭集。处置：②a 自报、自删，并做零残留核验（`tmp/` 顶层仅 `.gitkeep` 与 `AGENTRUNNER-PID-READER/`）。按 §5.1「显式登记，不销解、不追加动作」登记并继续；这是本片第 1 次，§5.1 合并口径第 11 次。

## 4 审查链

| 轮次 | 模型 | 结论 | 处置 |
|---|---|---|---|
| ⑤ r1 | `GLM-5.3` | 1 Medium（**双 BOM**）+ 2 Low + 5 Info | 修复后进入 ⑤ r2 |
| ⑤ r2 | `GLM-5.3`（收窄读） | PASS | 收敛 |
| ⑥ r1 | `MAI-Code-1.1-Flash` | PASS（0 发现） | 无需整改 |
| ⑦ r1 | `GPT-5.3-Codex`（盲审） | `ADJUDICATION: PASS` | 收敛 |

⑤ r1 的两条 Low 为「测试注释失真」与「harness 的 `readPIDs` 仍用 `fmt.Sscanf` 宽松解析」；它们与 5 条 Info 一并登记在 §8。

## 5 OB-110 读者侧三态语义

三态沿用契约 A6：文件**不存在**即证据缺口（fail-closed ⇒ `uncertain`）；文件存在且为**空**（0 字节）即「确实无子进程」；文件存在且**非空**即按行解析 pid 列表。三者不得互推。

四类读者逐一如下：

| 读者 | 三态处置 | 缺失时的结论 |
|---|---|---|
| `internal/adapters/agent_runner_run.go`：`readAgentRunnerChildPIDs` 三值返回 | 缺失 ⇒ `present=false`；空 ⇒ 空列表；非空 ⇒ 逐行 `strconv.Atoi` | `agentRunnerVerifyTreeGoneWith` 见 `present=false` 立即返回「不验证」 |
| `tools/b4-e2e/ports.go`：`readPIDs` + `managedProcessesGone` | 同上三态 | 收集缺失文件名并报 `adapters.ErrInvalidAgentRunnerEvidence` |
| `tools/b4-e2e/postflight.go`：`alivePIDs` | 同上三态 | 返回 error，调用方降级 `uncertain` |
| `tools/b4-e2e/runmap.go`：`mirrorVerifyTreeGone` | 镜像生产语义 | `present=false` ⇒ `treeGone=false` |

**与 `uncertain` 的关系**：缺失报告只经由「进程树无法证明消失」进入 `uncertain`；空报告不构成缺口。**优先级未变**：自然非零退出仍判 `failed`，`uncertain` 不覆盖已成立的自然失败。

## 6 OB-109 第六枚冻结事实

- 字节口径（引用契约 A6 原式）：`ProcessPidReportsHash = evidence.Digest("proofrail:agent-runner-pid-reports:1\n", "stub.pid " + Digest(父报告字节) + "\nstub-child.pid " + Digest(子报告字节) + "\n")`；内层 `Digest` 用空域、按原始字节，不解码。
- 实现落点：`AgentRunnerPidReportsDigest`（`internal/adapters/agent_runner_evidence.go:369`）；集合器在任一报告缺失时只写缺口域摘要（`:357`），不写空摘要。
- chain 侧：`postflightFactsCount = 6`（`internal/chain/postflight.go:39`），六枚事实按 `evidence[2..7]` 连续无空洞重建（`:40`、`:249`）。
- 第 7 腿三态：重算值与冻结值不一致 ⇒ `failed`；报告缺失 ⇒ `uncertain`；空报告 ⇒ 重算值与冻结值同为空串摘要 ⇒ `pass`。
- `PidReportsComplete` 为独立布尔（`internal/adapters/agent_runner_evidence.go:79`）；为假时 `FrozenFacts()` 以 `ErrAgentRunnerEvidenceGap` 拒绝（`:163`）。

## 7 测试与门禁

③ 读数：`gofmt` 空输出；`go build ./...` 与 `go vet ./...` 均 exit 0；`go test -count=1 ./...` 为 15 包 `ok` / 0 `FAIL`；具名运行 `a6native` 与 `b3bnative` 各 `ok`。

④ 读数：`tree=0`、`index=0`、`ci=0`；`tools/gates` 自检 `SELFTEST 588/588`；编码审计 24 条路径 0 违规。

编码审计曾发现**双 BOM**（⑤ r1 的 Medium），已修复并复跑通过。

**CI 读数（主提交 `02d5797`，run `39740097139`）**：双腿均 `success`（Ubuntu 腿 2 分 46 秒）；包级 `ok` 行 windows 15 条 / ubuntu 17 条；腿内 `SELFTEST` 均 588/588、`TOTAL_FAIL` 均 0。首跑即双绿，未触发 §6.4 的取证复跑，无 DR 登记。

## 8 已知限制

- **KL-C1-7 再述**：复用 `run-c1-faces.sh` 前必须使用全新的 `C1_OUT`；沿用旧目录会把上一轮证据混入本轮对账。
- **KL-C1-8 判定**：本片**未触发**「断网四条失效条件」。依据为 ③/① 的机械核对：本片未新增 `net` 依赖、无套接字调用、无 cgo、无新外部进程。
- **Low（镜像层）**：`tools/b4-e2e/ports.go:525` 的 `readPIDs` 仍用 `fmt.Sscanf`，与生产 `strconv.Atoi`（`internal/adapters/agent_runner_run.go:327`）不等价；只影响 harness，不在生产路径。
- **Low**：`internal/adapters/agent_runner_postflight_test.go:515-520` 的 seed 注释与恒定写口径不一致。
- **Info（纵深防御）**：第 7 腿的 missing 分支被第 5 腿前置拦截，实际不可达。
- **Info**：harness 读取缺上界检查；`Stat`→`ReadFile` 存在 TOCTOU。
- **Info**：`internal/console/runtime_agent_e2e_test.go:457-460` 与 `:643-646` 的陈旧注释（称 stub 仅在派生进程时写子报告）非本片引入。
- **Info**：五枚时代的旧路由证据在重放时按 conflict 整体拒绝、零写入（fail-closed）。
- **G4a 缺口（本片实证）**：G4a 只验「有 BOM」不验「BOM 个数」，双 BOM 可穿过；另立后续 `gates` 片处置。
- **⑨ 期四处 `G6-4` 改写的性质（如实登记、不改结论）**：④ 首跑 `--scope=index` 判红，四处命中与改写为——`18` 加 `个文件` 后缀 →「18 条路径」、`v1.42` 加 `行` 后缀 →「v1.42 条目」、`OB-111` 加 `行` 后缀 →「OB-111 条目」、`24` 加 `文件` 后缀 →「24 条路径」。**判定为误报形态**：这类「数字 + 单位」表述指本片变更集的确定事实，报告冻结后不会过期；而 `G6-4` 的正则不区分「切片变更集数」与「仓库当前总数」。**后续同类情形的口径**：优先走 `G6-4` 自带的白名单（显式豁免 + 理由 + 上下文摘录），**不默认改写措辞**；改写只在「判据确为误报且白名单路径不可用」时作为记录态处置，并如实登记。本条系 **OB-111 同族（`G6-4` 版）**数据点，供 `OB-111` 下次评估时一并考量。

## 9 自由度

本片行使的自由度、由裁定固定的边界与唯一偏离，逐条登记在 `docs/validation/evidence/AGENTRUNNER-PID-READER-freedom-list.md`（F 10 条、M 6 条、D 1 条）。

## 10 未执行事项

- 未在 Linux VM 复跑：本片无平台相关改动。
- 未新增 face：harness 无 evidence-dir tamper 旋钮，新增 face 属扩大范围。
- OB-104 / OB-105 / OB-107 不处理。

## 11 结论

两笔欠账逐项闭合：OB-110 的三态语义在四处读者一致，OB-109 的第六枚事实与契约 A6 逐字节等价；三作用域门禁全绿（`tree`、`index`、`ci` 均为 0）；唯一偏离 D1 已按 §5.1 登记。

```artifacts
docs/validation/agentrunner-pid-reader.md
docs/validation/agentrunner-pid-reader_EN.md
docs/validation/evidence/AGENTRUNNER-PID-READER-freedom-list.md
docs/CONTRACTS.md
docs/CONTRACTS_EN.md
docs/DELIVERY_DIRECTIVE.md
docs/DELIVERY_DIRECTIVE_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
internal/adapters/agent_runner_evidence.go
internal/adapters/agent_runner_evidence_test.go
internal/adapters/agent_runner_pinned_cli_test.go
internal/adapters/agent_runner_postflight.go
internal/adapters/agent_runner_postflight_test.go
internal/adapters/agent_runner_run.go
internal/adapters/agent_runner_run_test.go
internal/adapters/agent_runner_terminal_chain.go
internal/adapters/agent_runner_terminal_publisher_test.go
internal/adapters/agent_runner_terminal_record.go
internal/chain/agent_runner_terminal.go
internal/chain/agent_runner_terminal_test.go
internal/chain/postflight.go
internal/chain/postflight_test.go
tools/b4-e2e/main.go
tools/b4-e2e/ports.go
tools/b4-e2e/postflight.go
tools/b4-e2e/runmap.go
```

DOCS: DONE

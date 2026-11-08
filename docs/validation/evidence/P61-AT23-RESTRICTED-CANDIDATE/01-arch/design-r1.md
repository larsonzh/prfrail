# ① 执行计划确认（design-r1）— P61 AT23-RESTRICTED-CANDIDATE

落盘：2026-11-08（① 收尾）。来源：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)` 的 ① 回执，经用户 2026-11-08 四项裁定修订后定稿。

## A. 用户裁定修订（四要素：裁决人 = 用户；日期 = 2026-11-08；理由 = ① 澄清轮；依据 = 本轮裁定消息）

1. **`-Unlock` 语义（修正）**：默认（`Unlock=$false`）argv 含 `--no-auto-update` ⇒ 预期 `1.0.83`；`-Unlock` 显式传时 argv **不含** `--no-auto-update`（SEA loader 按数字前缀选最高版本，B1 已实测）⇒ 预期 `1.0.85`。裸跑版本串 ∉ {1.0.83, 1.0.85} ⇒ 停下，按 R-08.1 降级分支登记「版本不可控」。前置只读探查已做：`%LOCALAPPDATA%\copilot\pkg\win32-x64\` 下 1.0.83 与 1.0.85 两版目录均在位，无下载顾虑。
2. **Go 侧偏离登记（措辞修正）**：Go 侧 `copilotCLIProbeArgs`（`internal/adapters/ai_probe_copilot.go`）含 `--no-auto-update`，机制上安全；但其与 R-08.1 字面的 `buildCopilotArgs(base, unlock)` **命名与签名不同** ⇒ 如实登记「**与 R-08.1 字面的命名/签名偏离**」，进 ⑤ 预审必核面，由 ⑤ 判断是否可接受。本片不改 Go 侧（改 Go = 产品代码语义改动，超出 ② 范围）。
3. **第 2 笔（② 内 Copilot-hosted GPT-6 Luna 1 次）**：已批准。执行串 `copilot --no-auto-update --model gpt-6-luna -p "say hello"`（隔离 workspace、无 BYOK env）；`--model` 形态经 `--help` 实测（`--model <model>`，示例 `gpt-5.4`）；报 unknown/usage error ⇒ 记 unknown、不重试不追加。计 1 个 Copilot premium（独立于 BYOK 计费），不进验收结论。
4. **层 2 口径（写入定稿）**：机制腿 = **B4 harness 重跑**（`tools/b4-e2e/run-faces.ps1`，复用生产包、无 `os/exec`），**「非等价替代，与 B4 同口径」**；生产派发/终局发布/门禁端口对在 Windows `unproven` 下被 fail-closed 拒绝（B4 附录 A 冻结事实，引用为基线、不重跑）。层 2 核心场景「必通」= 本片证据包内 harness verdict 类与 B4 基线**逐项一致**。
5. **费用预批**：层 3 真实候选腿 3~5 次 BYOK，预估 **¥0.03~0.05 CNY**（按层 0 实测 ¥0.01/次）；全片 BYOK 预估 ¥0.05~0.07 + 1 个 Copilot premium；层 3 前仍逐次声明。

## B. 层 0（已绿，一行结论）

2026-11-08 全绿：钉版 1.0.83 `--no-auto-update -p "say hello"` exit 0、真实计费 ¥0.01 CNY（DeepSeek 侧核实）、`BYOK env cleanup: verified (all six cleared)`；证据 `01-arch/layer0-result-2026-11-08.md`。

## C. 层 1 骨架 5 项（含修正语义）

| # | 动作 | 命令形态 | 机械断言 | 证据落点 |
|---|---|---|---|---|
| 1 | 默认封装版本 | `Invoke-CopilotCli -Subcommand version`（argv 含 `--no-auto-update`） | stdout 含 `GitHub Copilot CLI 1.0.83.` 且 exit 0 | `03-test/layer1/l1-01-wrapper-version.txt` |
| 2 | 裸跑对照 | `copilot --no-auto-update --version` | 同项 1 版本串逐字 | `l1-02-bare-version.txt` |
| 3 | `-Unlock` 裸跑 | `Invoke-CopilotCli -Subcommand version -Unlock`（argv **不含** `--no-auto-update`） | 报 `1.0.85`；∉ {1.0.83,1.0.85} ⇒ 停走 R-08.1 降级分支 | `l1-03-unlock-version.txt` |
| 4 | BYOK 1 次 | tmp 凭据注入件包封装：`Invoke-CopilotCli -Prompt "say hello"` | exit 0 + 正常回复；finally 六项清理读数全空 | `l1-04-byok-call.txt` |
| 5 | `--effort=high` 三态 | BYOK：`copilot --no-auto-update --effort=high -p "say hello"` | 支持 ⇒ 追加测 low/medium/xhigh/max 记集合，层 3 固定 high；不支持/静默忽略 ⇒ 记已知边界不阻断 | `l1-05-effort-probe.txt` |

闸门：任一红 ⇒ 停修（单点突破），不进层 2。

## D. 层 2（核心必通）+ 层 3（全场景 + 条件闭合）

- 层 2 核心：F1 隔离 workspace、F4 未知恢复 —— harness 正反例断言（F1：`completed`/五枚事实互异/工作树 diff 空；反例 `uncertain-or-paused` + `REPAIR_PENDING`。F4：重启零 relaunch、派发计数恒 1、身份镜像逐字相等；反例伪 resume `refused` 零写）。非核心：F2 timeout、F3 日志缺失（挂了记录后继续）。
- 层 3：F2/F3/F4 三场景重跑 + postflight 六枚事实重建逐字节比对 + 真实候选腿 ≥3 次（启动/执行/停机，每次 stderr 原文/退出码/失败阶段；预期复现 T026 两项 unsupported）。
- 条件闭合前提核对（④ 前勾选）：① 机制腿全绿（verdict 类与 B4 基线逐项一致）② 真实候选腿 ≥3 次且证据具体 ③ 报告三句话逐字出现 ④ BYOK 结论未替代 B1 基线 ⑤ `--effort`/`--context` 已知边界已登记。

## E. ② 实现清单（P61 首个生产文件改动）

| 文件 | 动作 | 内容 |
|---|---|---|
| `tools/agent-probe/CopilotCli.ps1` | **新建**（UTF-8 BOM+LF；英文注释/输出） | `buildCopilotArgs`（纯函数；默认恒定前缀 `--no-auto-update`；`-Unlock` 时移除该前缀）+ `Invoke-CopilotCli`（`-Subcommand/-Prompt/-Effort/-Context/-ExePath/-Unlock`；默认 fail-closed；版本钉住 1.0.83，阶段 1 首次 + 每 5 次抽样、证据仅 `--version` 输出；失配拒绝；`-Unlock` 记录实际版本串；三条证据对照） |
| `tools/agent-probe/Invoke-AgentProbe.ps1`、`Invoke-CandidateProbe.ps1`、`appcontainer-b2/container-candidate.ps1` | 修改 | 既有 6 处直接调用点收敛为走封装（行为不变） |
| `internal/adapters/ai_probe_copilot*.go` | **不改** | 偏离登记见 §A.2；⑤ 必核面 |

- BYOK 密钥注入**不进封装**：继续走 tmp 凭据读取器配套件（CredReadW → 注入 → finally 清理），密钥不落任何文件。
- ② 内 Luna 1 次（第 2 笔，已批准）：证据 `02-impl/luna-connectivity.txt`。

## F. 付费清单（修订后编号，R-08.2 逐次声明）

| # | 时点 | 内容 |
|---|---|---|
| 1 | 层 0（已付） | BYOK say hello，¥0.01 CNY |
| 2 | ② 内（已批准） | Copilot-hosted GPT-6 Luna 1 次 = 1 个 Copilot premium |
| 3 | 层 1 项 4 | BYOK 封装路径 say hello，约 ¥0.01 |
| 4 | 层 1 项 5 | `--effort=high` 三态探测（BYOK），约 ¥0.01 |
| 5~7+ | 层 3 | 真实候选腿 3~5 次 BYOK，预估 ¥0.03~0.05 |

## G. 证据包结构与机械断言

- 结构：`01-arch/`（本件 + layer0 系列）、`02-impl/`（CopilotCli.ps1 快照 + sha256、收敛 diff、impl-notes、luna-connectivity）、`03-test/`（layer1/、layer2/、layer3/、real-candidate/run-1..n、postflight-rebuild/）、`04-reviews/`（05/06/07 各轮原文）、`06-gates/`（go build/vet/test、三作用域、selftest）、根 `MANIFEST.md`/`SHA256SUMS.txt`（每轮复算）。
- 收尾机械断言：`go build/vet/test ./...` 全绿；`gate.js --all --scope=index/tree/ci` 全绿；`--selftest` 722/722。

## H. 风险与边界（不简化）

R-08.1 fail-closed 与「新增调用点走封装」（⑤ 标偏离）；R-08.5 条件闭合口径与报告三句话；R-09.3 网络底线（回环可接受、公网阻断；BYOK → api.deepseek.com 为裁定授权例外）；R-09.7 外部副作用须用户显式确认；CONTRACTS §7 冻结句不改（⑧ 标「不适用 + 理由」）；阶段 1→2 三条同时满足、主控不自行升级；BYOK 预检失败/密钥缺失 ⇒ 停下升级；⑦ 额度 1+3+1（含盲审），用尽 ⇒ 升级用户裁决。

DESIGN: DONE（修订版，2026-11-08）

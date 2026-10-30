# ProofRail（证轨）

> *Provable rails for unattended AI engineering* —— 为 AI 无人值守编程铺设可证明的轨道。

[简体中文](#简体中文) · [English](#english)

---

## 简体中文

ProofRail 让 AI 在无人值守下安全地改代码、跑验证、出证据：每一步都可验证、可回滚、可审计，
人的审批始终保留在关键节点。

- 首错即停 + 哈希绑定产物 + 原子提交回滚的“可证明写入”模型
- 相同失败指纹预算与有效修复证据判定的“防 AI 无限循环”机制
- 编辑边界矩阵与停机门禁的“AI 权限最小化”模型
- 隔离候选事务（candidate 不污染正式定义）的自愈安全模型

### 状态

**防御性重构期（P39–P45）**（核对时点 2026-10-31；C6 预算已落盘，C1–C7 逐项攻克中；各项以 [docs/DEV_PLAN.md](docs/DEV_PLAN.md) 的任务状态为准）

| 能力 | 状态 |
|---|---|
| 离线 no-AI 链路（`init` / `validate` / `config explain` / `preview` / `approvals` / `interactions` / `cost report` / `secret` / `report`） | 可用 |
| `ai check` 单请求 live availability probe | 可用（只证明可达性，不证明 AgentRunner enforcement 或 AT-23） |
| `export`、`effects` / `diagnostics` | 部分可用（库级实现已具备，CLI 入口待后续切片） |
| 受监督黑箱模式（supervised black-box） | 部分可用（显式风险确认 + 隔离 workspace + 人工监督 + 归还后全量扫描；只保证归还后产物，不保证工具、网络、费用或外部副作用） |
| `run` 执行 `code` / `build` / `verify` 任务链 | 未交付（当前只执行 noop-only 链，遇到可执行步骤 fail-close 并返回非零；T027 装配链已有可跑拒绝路径，持久预算已落盘） |
| 正式 AI 执行（AgentRunner 真实候选腿） | 未交付（AT-23 未通过；C1–C7 与 N1–N4 已归档，P39 起逐项修复） |
| 完整统一 TUI / `serve` | 未交付 |
| 正式发行包与安装 | 未交付 |

本项目仍不是正式发行包。分域权威与文档导航见 [docs/DOCUMENTATION_PLAN.md](docs/DOCUMENTATION_PLAN.md)；
执行顺序与任务状态见 [docs/DEV_PLAN.md](docs/DEV_PLAN.md)，就绪缺口见 [docs/ADR_REGISTER.md](docs/ADR_REGISTER.md)；
独立产品端到端叙事见 [docs/BUSINESS_WORKFLOWS.md](docs/BUSINESS_WORKFLOWS.md)；尚未冻结的安装发行方案见 [docs/INSTALLATION_PLAN.md](docs/INSTALLATION_PLAN.md)。
产品建议书与历史设计来源见 [docs/RFC-proofrail-unattended-ai-engineering-product.md](docs/RFC-proofrail-unattended-ai-engineering-product.md)。
产品完整性评审补充了无副作用预览、已接受结果导出、授权撤销、外部副作用边界、成本预留结算和备份/停用流程；见 [docs/PRODUCT_REQUIREMENTS.md](docs/PRODUCT_REQUIREMENTS.md) §8。文件回滚不保证撤销外部操作。

### 快速开始（当前 CLI 基线）

当前界面是 Windows Terminal、PowerShell 或 VS Code 集成终端中的逐行 CLI，不是图形窗口或完整 TUI。发布候选的软硬件要求、下载核验和五步上手见 [Windows 便携 ZIP 安装指南](docs/INSTALLATION.md)。核心离线 CLI 不要求 VS Code。目标正式 AI 执行使用 AgentRunner/CLI Agent，尚未交付（AT-23 未通过）。基于 SessionBridge visible 的受监督黑箱模式已实现并提供显式风险确认与人工监督，但只保证归还后的产物检查，不保证工具、网络、费用或外部副作用。

```text
prfrail init --workspace .
prfrail validate --chain ./proofrail.chain.json
prfrail config explain --chain ./proofrail.chain.json
prfrail preview --chain ./proofrail.chain.json
# 仅在明确授权一次真实模型请求后执行：
prfrail ai check --chain ./proofrail.chain.json --channel agent-runner-cli --copilot C:/path/to/copilot.exe --workspace . --max-requests 1 --out ./tmp/ai-availability.json --json
prfrail ai verify --chain ./proofrail.chain.json --record ./tmp/ai-availability.json --channel agent-runner-cli --max-age 10m --max-requests 1 --json
prfrail approvals list --ledger ./authorization-ledger.jsonl
prfrail interactions list --ledger ./operator-interactions.jsonl
prfrail interactions tui --ledger ./operator-interactions.jsonl --actor-id operator-one
prfrail run --chain ./proofrail.chain.json --run-id run-demo
prfrail report --run-dir ./tmp/prfrail-runs/run-demo
```

说明：`run` 目前仅执行 noop-only 任务链；`interactions tui` 是 T025 的聚焦终端收件箱，不是完整统一 TUI，`serve`/完整 TUI 仍在后续切片。`ai verify` 可重复指定 `--prior-record <path>`，离线拒绝同一 probe ID 的历史分叉；相同记录重放保持幂等。黑箱模式须先展示降级保证和责任边界。SessionBridge visible 投递不等于完成，正式流程也不依赖 `@sbr-review`。

### 构建（需 Go 工具链）

```text
go build ./...
go vet ./...
go test ./...
```

### 路线图（S0–S3）

| 阶段 | 目标 | 关键交付 |
|---|---|---|
| S0 规格化 | 领域模型与协议定稿 | 文档包、JSON Schema、威胁模型、ADR |
| S1 核心任务链 | 可运行的 MVP | Chain Engine、checker、gate runner、adapter、快照/评审/恢复、TUI |
| S2 平台与语言扩展 | 走向通用 | 多 harness、Linux 支持、Web 控制台、模型策略 |
| S3 高级能力 | 产品化完整闭环 | 生成钩子场景 B、多编辑器、多语言文档 |

### 命名

| 场景 | 写法 |
|---|---|
| 正式品牌、标题、对外文档 | **ProofRail** |
| 紧凑视觉标识 | **PrfRail** |
| 仓库、CLI、包名、模块前缀 | `prfrail` |
| 中文文档 | **证轨** |

### License

MIT（见 [LICENSE](LICENSE)）。

---

## English

ProofRail enables AI to safely modify code, run validations, and produce evidence in unattended mode: every step is verifiable, reversible, and auditable, while human approval always remains at critical checkpoints.

- "Provable write" model: fail-fast on first error + hash-bound artifacts + atomic commit/rollback
- Anti-infinite-loop mechanism: identical failure-fingerprint budgets + valid-fix-evidence determination
- Least-privilege AI model: edit boundary matrix + halt gates
- Self-healing security model: isolated candidate transactions (candidate never contaminates formal definitions)

### Status

**Defensive refactor period (P39–P45)** (checked 2026-10-31; C6 budget landed and C1–C7 under attack item by item; each row follows the task state in [docs/DEV_PLAN_EN.md](docs/DEV_PLAN_EN.md))

| Capability | State |
|---|---|
| Offline no-AI path (`init` / `validate` / `config explain` / `preview` / `approvals` / `interactions` / `cost report` / `secret` / `report`) | available |
| `ai check` single-request live availability probe | available (proves reachability only, not AgentRunner enforcement or AT-23) |
| `export`, `effects` / `diagnostics` | partial (library level shipped; CLI entrypoints in later slices) |
| Supervised black-box mode | partial (explicit risk acknowledgement + isolated workspace + human supervision + full post-return scan; warrants post-return artifacts only, not tools, network, cost or external effects) |
| `run` executing `code` / `build` / `verify` chains | not delivered (noop-only today; executable steps fail closed with a non-zero exit; the T027 assembly chain already has a runnable refusal path and the durable budget has landed) |
| Formal AI execution (real AgentRunner candidate leg) | not delivered (AT-23 not passed; C1–C7 and N1–N4 are archived and fixed item by item from P39) |
| Complete unified TUI / `serve` | not delivered |
| Production release package and installation | not delivered |

This is still not a production release package. See [docs/DOCUMENTATION_PLAN_EN.md](docs/DOCUMENTATION_PLAN_EN.md) for domain authorities and navigation;
[docs/DEV_PLAN_EN.md](docs/DEV_PLAN_EN.md) for task state and [docs/ADR_REGISTER_EN.md](docs/ADR_REGISTER_EN.md) for readiness gaps;
[docs/BUSINESS_WORKFLOWS_EN.md](docs/BUSINESS_WORKFLOWS_EN.md) for the independent end-to-end narrative and [docs/INSTALLATION_PLAN_EN.md](docs/INSTALLATION_PLAN_EN.md) for unresolved installation and distribution decisions.
Project proposal and historical design source: [docs/RFC-proofrail-unattended-ai-engineering-product.md](docs/RFC-proofrail-unattended-ai-engineering-product.md).
Product review adds side-effect-free preview, accepted-result export, revocation, external-effect boundaries, cost reservation/settlement and backup/retirement flows; see [docs/PRODUCT_REQUIREMENTS_EN.md](docs/PRODUCT_REQUIREMENTS_EN.md) section 8. File rollback cannot guarantee undoing external effects.

### Quick Start (current CLI baseline)

The current interface is a line-oriented CLI in Windows Terminal, PowerShell, or the VS Code integrated terminal, not a graphical window or complete TUI. See the [Windows portable ZIP installation guide](docs/INSTALLATION_EN.md) for candidate hardware/software requirements, download verification, and the five-step quick start. Core offline mode needs no VS Code. Target formal AI execution uses AgentRunner/CLI Agent and is not delivered yet (AT-23 not passed). A supervised black-box mode built on SessionBridge visible delivery is implemented with explicit risk acknowledgement and human supervision, but it warrants post-return artifacts only, not tools, network, cost or external effects.

```text
prfrail init --workspace .
prfrail validate --chain ./proofrail.chain.json
prfrail config explain --chain ./proofrail.chain.json
prfrail preview --chain ./proofrail.chain.json
# Run only after explicit authorization for one live model request:
prfrail ai check --chain ./proofrail.chain.json --channel agent-runner-cli --copilot C:/path/to/copilot.exe --workspace . --max-requests 1 --out ./tmp/ai-availability.json --json
prfrail ai verify --chain ./proofrail.chain.json --record ./tmp/ai-availability.json --channel agent-runner-cli --max-age 10m --max-requests 1 --json
prfrail approvals list --ledger ./authorization-ledger.jsonl
prfrail interactions list --ledger ./operator-interactions.jsonl
prfrail interactions tui --ledger ./operator-interactions.jsonl --actor-id operator-one
prfrail run --chain ./proofrail.chain.json --run-id run-demo
prfrail report --run-dir ./tmp/prfrail-runs/run-demo
```

Notes: `run` currently executes noop-only chains. `interactions tui` is T025's focused terminal inbox, not the complete unified TUI; `serve` and the complete TUI remain later slices. `ai verify` accepts repeatable `--prior-record <path>` inputs, rejects historical forks for one probe ID offline, and treats identical record replay as idempotent. The black-box mode must show its reduced assurance and responsibility boundary before selection. SessionBridge visible delivery is not completion, and the formal path does not depend on `@sbr-review`.

### Build (requires Go toolchain)

```text
go build ./...
go vet ./...
go test ./...
```

### Roadmap (S0–S3)

| Phase | Goal | Key deliverables |
|---|---|---|
| S0 Spec | Finalize domain model and protocols | Docs package, JSON Schema, threat model, ADR |
| S1 Core chain | Runnable MVP | Chain Engine, checker, gate runner, adapter, snapshot/review/recover, TUI |
| S2 Platform & language extension | Go general | Multiple harnesses, Linux support, web console, model policies |
| S3 Advanced capabilities | Full productized loop | Generated hook scenario B, multiple editors, multilingual docs |

### Naming

| Context | Spelling |
|---|---|
| Official brand, titles, external docs | **ProofRail** |
| Compact visual identity | **PrfRail** |
| Repo, CLI, package, module prefix | `prfrail` |
| Chinese docs | **证轨** |

### License

MIT (see [LICENSE](LICENSE)).

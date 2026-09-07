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

**S1 实施中（2026-09-08）**：T005–T016 已完成证据、快照、进程/租约守卫、托管变更集与事务应用，以及无 IDE CLI 基线（`init/validate/config explain/run/report`）。
当前仍不是正式发行包：`run` 只支持 noop-only 任务链；遇到 `code/build/verify` 步骤会 fail-close 并返回非零退出码。
项目建议书与历史设计来源见 [docs/RFC-proofrail-unattended-ai-engineering-product.md](docs/RFC-proofrail-unattended-ai-engineering-product.md)；分域权威见 [docs/DOCUMENTATION_PLAN.md](docs/DOCUMENTATION_PLAN.md)。

从 [docs/DOCUMENTATION_PLAN.md](docs/DOCUMENTATION_PLAN.md) 阅读文档导航与低成本模型流程；
执行顺序见 [docs/DEV_PLAN.md](docs/DEV_PLAN.md)，就绪缺口见 [docs/ADR_REGISTER.md](docs/ADR_REGISTER.md)。
独立产品端到端叙事见 [docs/BUSINESS_WORKFLOWS.md](docs/BUSINESS_WORKFLOWS.md)；尚未冻结的安装发行方案见 [docs/INSTALLATION_PLAN.md](docs/INSTALLATION_PLAN.md)。

产品完整性评审已补充无副作用预览、已接受结果导出、授权撤销、外部副作用边界、成本预留结算和备份/停用流程；见 [docs/PRODUCT_REQUIREMENTS.md](docs/PRODUCT_REQUIREMENTS.md) §8。均为待批准设计，不是现有功能；文件回滚不保证撤销外部操作。

### 快速开始（当前 CLI 基线）

```text
prfrail init --workspace .
prfrail validate --chain ./proofrail.chain.json
prfrail config explain --chain ./proofrail.chain.json
prfrail run --chain ./proofrail.chain.json --run-id run-demo
prfrail report --run-dir ./tmp/prfrail-runs/run-demo
```

说明：`run` 目前仅执行 noop-only 任务链；`serve`/完整 TUI 仍在后续切片。

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

**S1 implementation in progress (2026-09-08)**: T005–T016 implement evidence, snapshots, process/lease guards, managed change sets, transactional apply, and a no-IDE CLI baseline (`init/validate/config explain/run/report`).
This is still not a production release package: `run` currently supports noop-only chains and fail-closes with a non-zero exit for executable `code/build/verify` steps.
Project proposal and historical design source: [docs/RFC-proofrail-unattended-ai-engineering-product.md](docs/RFC-proofrail-unattended-ai-engineering-product.md). Domain authorities: [docs/DOCUMENTATION_PLAN_EN.md](docs/DOCUMENTATION_PLAN_EN.md).

Start with [docs/DOCUMENTATION_PLAN_EN.md](docs/DOCUMENTATION_PLAN_EN.md) for navigation and the low-cost model workflow;
see [docs/DEV_PLAN_EN.md](docs/DEV_PLAN_EN.md) for tasks and [docs/ADR_REGISTER_EN.md](docs/ADR_REGISTER_EN.md) for readiness gaps.
See [docs/BUSINESS_WORKFLOWS_EN.md](docs/BUSINESS_WORKFLOWS_EN.md) for the independent end-to-end product narrative and [docs/INSTALLATION_PLAN_EN.md](docs/INSTALLATION_PLAN_EN.md) for unresolved installation and distribution decisions.

Product review adds side-effect-free preview, accepted-result export, revocation, external-effect boundaries, cost reservation/settlement and backup/retirement flows. See [docs/PRODUCT_REQUIREMENTS_EN.md](docs/PRODUCT_REQUIREMENTS_EN.md) section 8. These are unapproved designs, not available features; file rollback cannot guarantee undoing external effects.

### Quick Start (current CLI baseline)

```text
prfrail init --workspace .
prfrail validate --chain ./proofrail.chain.json
prfrail config explain --chain ./proofrail.chain.json
prfrail run --chain ./proofrail.chain.json --run-id run-demo
prfrail report --run-dir ./tmp/prfrail-runs/run-demo
```

Notes: `run` currently executes noop-only chains; `serve` and the full TUI remain in later slices.

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

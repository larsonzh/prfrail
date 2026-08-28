# ProofRail（证轨）

> *Provable rails for unattended AI engineering* —— 为 AI 无人值守编程铺设可证明的轨道。

ProofRail 让 AI 在无人值守下安全地改代码、跑验证、出证据：每一步都可验证、可回滚、可审计，
人的审批始终保留在关键节点。

- 首错即停 + 哈希绑定产物 + 原子提交回滚的“可证明写入”模型
- 相同失败指纹预算与有效修复证据判定的“防 AI 无限循环”机制
- 编辑边界矩阵与停机门禁的“AI 权限最小化”模型
- 隔离候选事务（candidate 不污染正式定义）的自愈安全模型

## 状态

**提案（Proposal）**：本文档对应产品设计基线 v0.1（2026-08-28），实现尚未开始。
设计规范见 [docs/RFC-proofrail-unattended-ai-engineering-product.md](docs/RFC-proofrail-unattended-ai-engineering-product.md)。

## 快速开始（规划中）

```text
prfrail init      # 向导：语言/项目类型 → 任务链模板 → 环境 → 生成 chain-file
prfrail validate  # schema + 预检（进程/锁/工具链/远程）
prfrail run       # 单命令执行，TUI 实时进度
prfrail report    # 链报告
prfrail serve     # 本地 Web 控制台（P1）
```

## 构建（需 Go 工具链）

```text
go build ./...
go test ./...
```

## 路线图（S0–S3）

| 阶段 | 目标 | 关键交付 |
|---|---|---|
| S0 规格化 | 领域模型与协议定稿 | 文档包、JSON Schema、威胁模型、ADR |
| S1 核心任务链 | 可运行的 MVP | Chain Engine、checker、gate runner、adapter、快照/评审/恢复、TUI |
| S2 平台与语言扩展 | 走向通用 | 多 harness、Linux 支持、Web 控制台、模型策略 |
| S3 高级能力 | 产品化完整闭环 | 生成钩子场景 B、多编辑器、多语言文档 |

## 命名

| 场景 | 写法 |
|---|---|
| 正式品牌、标题、对外文档 | **ProofRail** |
| 紧凑视觉标识 | **PrfRail** |
| 仓库、CLI、包名、模块前缀 | `prfrail` |
| 中文文档 | **证轨** |

## License

MIT（见 [LICENSE](LICENSE)）。

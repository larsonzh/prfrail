# Copilot Instructions for prfrail（证轨）

面向：在 VS Code 中协助维护本仓库的 AI 代理。请先掌握结构与规范，避免破坏既有契约。

## 定位与权威
- 产品：ProofRail（证轨，`prfrail`）——为 AI 无人值守编程铺设**可证明轨道**：
  每个步骤可验证、可回滚、可审计，人工审批保留在关键节点。
- **项目建议书与历史设计来源**：`docs/RFC-proofrail-unattended-ai-engineering-product.md`；文件名为兼容保留，
  不再是全域唯一权威。分域权威见 `docs/DOCUMENTATION_PLAN.md`：需求、流程、架构、契约、安全、验收和
  实施分别由对应专项文档负责；协议修订先改 `docs/CONTRACTS.md` 及 Schema/样例，再改代码。
- 工程规范：`docs/CODING_CONVENTIONS.md`（编码/行尾/Git 纪律）。
- 姊妹仓库 `sessbridge`（同多根工作区）：ProofRail 只消费其 `silent` 消息层 + 文件队列；
  `visible` 人工交互不属于 ProofRail 正式协议。

## 结构速览
- 模块：`module github.com/larsonzh/prfrail`，Go 1.22+；CLI 入口 `cmd/prfrail/main.go`。
- `internal/`：`adapters`（适配器）、`applier`（写入器）、`chain`（任务链）、`console`、
  `evidence`（证据）、`gates`（门禁）、`guard`、`repair`、`snapshot`（快照）、
  `taskdef`（任务定义）、`tickets`（票据）。

## 编码 + 行尾（硬规则，见 docs/CODING_CONVENTIONS.md）
- `.md` / `.ps1`：UTF-8 **with BOM** + **LF**。
- 所有 `.json`：UTF-8 **without BOM** + **LF**；仅明确用于 BOM 兼容性测试的 fixture 可例外。
- `.github/hooks/context-mode.json` 还必须保留 `"$schema": "../context-mode-schema.json"`；
  `context-mode upgrade` 后需检查无 BOM + 保留 `$schema`。
- 其它（`.go`/`.py`/`.js`/`.txt` 等）：UTF-8 **without BOM** + **LF**。

## 工程约定
- Go 代码 `gofmt`；`go build ./...`、`go vet ./...`、`go test ./...` 必须通过。
- 修改任务链/门禁/证据模型等核心语义前，先读 `docs/CONTRACTS.md` 对应章节；仅追溯设计理由或处理
  尚未迁移内容时再读项目建议书。跨文档冲突必须先修订，不得自行选择。
- `.vscode/scripts/` 在 `.gitignore` 中（不入库），改动仅限本机。
- 临时/一次性文件放仓库根 `tmp/`（已 gitignore，仅 `tmp/.gitkeep` 保留目录）；
  **用完立即清除**，不得长期遗留。

## Git 纪律（硬规则）
- 默认只推 `origin`（GitHub）；**gitee 仅镜像备份**——未经用户同一轮显式要求，
  **禁止任何 gitee push**（即使说“push all changes”也只推 origin）。
- **未经用户同一轮显式授权，禁止 `git commit` / `git push`**；授权后使用
  `git add <具体文件>` 精确暂存，避免 `git add -A`。
- 提交信息风格：`<type>: <summary>`（chore/feat/fix/docs/test/build）。

## 测试与验证
- 核心验证：`go build ./...`、`go test ./...`；符合 Go 1.22 工具链。
- 涉及 context-mode/hook 配置：编辑器内 MCP doctor 为准；
  CLI 独立 doctor 的 hooks 目录误报可忽略（无工作区上下文时）。

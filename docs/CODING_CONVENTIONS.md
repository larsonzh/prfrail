# ProofRail（证轨）编码与工程规范（Coding Conventions）

> 状态：2026-09-04 定稿。适用：本仓库所有提交文件的创建与维护。

## 1. 编码格式 + 行尾序列（硬规则）

为保证跨平台使用（Windows / Linux / macOS、Go / Python / PowerShell 5.1），
所有**提交进仓库**的文件必须满足：

| 文件类型 | 编码 | 行尾 | 说明 |
|----------|------|------|------|
| Markdown（`.md`） | UTF-8 **with BOM** | **LF** | 含 `docs/*.md`、`copilot-instructions.md` |
| PowerShell（`.ps1`） | UTF-8 **with BOM** | **LF** | PS 5.1 读写兼容、中文注释稳定 |
| JSON（`.json`） | UTF-8 **with BOM** | **LF** | 纯数据/样例/文档引用 |
| JSON（例外） | UTF-8 **without BOM** | **LF** | 被严格 JSON 解析器消费的文件，见 §1.1 |
| 其它（`.go`/`.py`/`.js`/`.txt`/`.gitignore`/`LICENSE`…） | UTF-8 **without BOM** | **LF** | 代码与明文文本 |

### 1.1 例外（必须无 BOM）

- **`.github/hooks/context-mode.json`**：context-mode CLI 用 `JSON.parse` 读取，
  带 UTF-8 BOM 会解析失败（已验证）。必须为 **无 BOM + LF**，并保留
  `"$schema": "../context-mode-schema.json"` 行（VS Code 校验需要关联 schema）。
  - 升级 `context-mode` 后检查该文件：仍无 BOM + 保留 `$schema`；如被重写需恢复。
- 其它被 Node / Go `json.Decoder` 等严格消费的 JSON 同理，不得带 BOM。

### 1.2 机械门禁

- 提交前自查：`.md/.ps1/.json` 带 BOM+LF；`.go` 等无 BOM+LF。
- 需要时在 `.vscode/scripts/`（不入库）增加批量规范化脚本，但不得改变内容语义。
- `.vscode/` 目录位于 `.gitignore`（如 `.vscode/scripts/proxy-toggle.ps1`），不入库。

## 2. Go 工程约定

- 模块：`github.com/larsonzh/prfrail`，Go 1.22+。
- 代码格式：`gofmt`；命名遵循 Go 惯例（导出符号大写、注释以符号名开头）。
- 结构：`cmd/prfrail/`（CLI 入口）、`internal/`（`adapters` `applier` `chain` `console`
  `evidence` `gates` `guard` `repair` `snapshot` `taskdef` `tickets`）。
- 验证：`go build ./...`、`go vet ./...`、`go test ./...`。
- 配置文件/产物路径约定遵循 RFC 与现有代码，避免魔法路径散落。

## 3. 文档与协议

- 权威设计规范：`docs/RFC-proofrail-unattended-ai-engineering-product.md`；
  协议修订走 RFC 流程，先改 RFC 再改代码。
- 与 `sessbridge`（姊妹仓库，同工作区）边界：ProofRail 正式协议只消费
  SessionBridge `silent` 消息层 + 文件队列；`visible` 人工交互属于 SessionBridge 产品能力，
  不进入 ProofRail 正式协议。

## 4. Git 纪律

- 默认只推 `origin`（GitHub）；**gitee 仅为镜像备份**，未经用户同一轮显式要求
  **禁止**任何 gitee push。
- 未经用户同一轮显式授权，**禁止** `git commit` / `git push`；
  仅允许在获得明确授权后执行，且新增文件用 `git add <file>` 精确暂存。
- 提交信息风格：`<type>: <summary>`（chore/feat/fix/docs/test/build）。

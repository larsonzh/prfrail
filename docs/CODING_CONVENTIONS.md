# ProofRail（证轨）编码与工程规范（Coding Conventions）

> 状态：2026-09-04 编码/Git 规则定稿；2026-09-06 补充原型前工程流程。适用：本仓库所有提交文件的创建与维护。

[English](CODING_CONVENTIONS_EN.md)

## 1. 编码格式 + 行尾序列（硬规则）

为保证跨平台使用（Windows / Linux / macOS、Go / Python / PowerShell 5.1），
所有**提交进仓库**的文件必须满足：

| 文件类型 | 编码 | 行尾 | 说明 |
|----------|------|------|------|
| Markdown（`.md`） | UTF-8 **with BOM** | **LF** | 含 `docs/*.md`、`copilot-instructions.md` |
| PowerShell（`.ps1`） | UTF-8 **with BOM** | **LF** | PS 5.1 读写兼容、中文注释稳定 |
| JSON（`.json`） | UTF-8 **without BOM** | **LF** | 配置、Schema、样例及运行时 JSON 统一格式 |
| 其它（`.go`/`.py`/`.js`/`.txt`/`.gitignore`/`LICENSE`…） | UTF-8 **without BOM** | **LF** | 代码与明文文本 |

### 1.1 JSON 例外与兼容性 fixture

- `.github/hooks/context-mode.json` 与 `schemas/*.schema.json` 同样遵循 JSON 默认规则：
  **UTF-8 无 BOM + LF**。前者还必须保留 `"$schema": "../context-mode-schema.json"`。
- 仅用于验证 BOM 输入兼容性的 golden fixture 可带 BOM；文件名或 fixture 元数据必须明确测试目的，
  且不得复用为配置、Schema、canonical/wire 数据或普通黄金样例。
- UTF-8 无 BOM 完整支持中文等 Unicode 文本；BOM 不是多语言兼容所必需，并可能使严格 JSON
  解析器失败。

### 1.2 机械门禁

- 提交前自查：`.md/.ps1` 带 BOM+LF；`.json`、`.go` 等无 BOM+LF。
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

- 项目建议书与历史设计来源：`docs/RFC-proofrail-unattended-ai-engineering-product.md`；
  分域权威见 `docs/DOCUMENTATION_PLAN.md`。协议修订先改 `docs/CONTRACTS.md` 及 Schema/样例，再改代码。
- 与 `sessbridge`（姊妹仓库，同工作区）边界：ProofRail 正式协议只消费
  SessionBridge `silent` 消息层 + 文件队列；`visible` 人工交互属于 SessionBridge 产品能力，
  不进入 ProofRail 正式协议。

## 4. Git 纪律

- 默认只推 `origin`（GitHub）；**gitee 仅为镜像备份**，未经用户同一轮显式要求
  **禁止**任何 gitee push。
- 未经用户同一轮显式授权，**禁止** `git commit` / `git push`；
  仅允许在获得明确授权后执行，且新增文件用 `git add <file>` 精确暂存。
- 提交信息风格：`<type>: <summary>`（chore/feat/fix/docs/test/build）。

## 5. 临时文件（tmp/）

- 仓库根目录 `tmp/` 用于放置**临时/一次性**文件（构建/测试产物、调试快照、中间脚本等）；
  已被 `.gitignore` 忽略（保留 `tmp/.gitkeep` 占位以便克隆后目录存在）。
- **临时文件使用完后应及时清除**：验证结束或调试结束即删除，不得长期遗留。
- `.vscode/` 与 `tmp/` 均不入库；需要共享的工程脚本应放入入库目录（如 `tools/`）。

## 6. Go 实现与审查规则

- 先依据分域权威文档固定行为，再写 table-driven 反例测试与最小实现。不得实现待批准的默认值、权限或 wire 字段。
- 接口由消费者定义；CLI 负责装配，核心不得反向依赖具体 adapter/console。仅为真实测试/替换边界抽象，避免一次生成全套空接口。
- 使用 `context.Context` 传播取消/期限；子进程与 goroutine 有明确停止和回收路径。时钟、进程、文件故障注入可替换，不能用 sleep 猜测完成。
- 错误保留原因和对象身份，使用包装及 errors.Is/As 判定，不以字符串匹配控制状态；禁止吞掉写入/关闭/flush/rollback 错误。
- 不从包级 init 启动进程或写文件；禁止库内 os.Exit/log.Fatal。外部命令 executable/args 分离，不隐式 shell。
- 配置/协议使用结构化解析器；未知字段拒绝；哈希用冻结 canonical 算法与原始内容字节，不能用 map 遍历顺序或字符串拼接。
- 文件路径验证覆盖穿越、别名、symlink/reparse 和检查-使用竞态；写状态/对象先遵循 journal/atomic 协议，不把 rename 当整组事务。
- 测试使用 t.TempDir 和非秘密夹具，故障注入只杀 helper；不修改真实工作树、已发布 seed 或用户进程。任何相关进程仍运行时禁止回滚文件。
- 依赖限标准库/经审查纯 Go 包；锁定版本、许可证、最低 Go、来源。不得顺手升级 Go 或增添运行时服务。

## 7. 低成本实施与完成定义

从 [文档导航](DOCUMENTATION_PLAN.md) 找当前 T/REQ/AT，只读相关契约小节与附近测试；一条可证伪假设、一次最小编辑、立即窄验证，再做全 build/vet/test。格式化用 gofmt，不能为了绿灯删断言、弱化门禁、扩大路径或重置预算。

每次交付说明改动范围、实际命令/退出码/用例数、未验证项、下一步。连续两次无新证据的同类失败交维护者；没有授权不增加收费模型调用。无测试不算功能验收，cross-compile 不算原生运行，模型自述不算证据。

## 8. 双语与贡献流程

文档先改对应分域规范，再同步测试/示例和用法。中文保持原名，英文 `_EN`，互链；REQ/AT/T/ADR 等编号相同、同次评审。冲突按 `DOCUMENTATION_PLAN.md` 的权威矩阵定位并阻断，不因翻译改变状态机或字段。

提交评审材料包括需求/任务 ID、最小行为差异、测试证据、协议兼容和文档影响。安全/协议变更需独立评审，所有者批准阶段与资源预算。不得把文档创建或计划覆盖率当 S0 readiness 通过。

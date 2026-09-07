# ProofRail 用户、操作员与可信发布手册

[English](OPERATIONS_EN.md)

日期：2026-09-08；S1 操作设计稿。当前没有 ProofRail 正式安装包；但 T016 已交付无 IDE CLI 基线：`version/init/validate/config explain/run/report`。`run` 目前只支持 noop-only 任务链，遇到 `code/build/verify` 会 fail-close 并返回非零退出码。除以下“当前可执行”外，其余产品能力仍是 RFC 规划接口。待决发行方案见 [安装部署规划](INSTALLATION_PLAN.md)，完整用户流程见 [业务流程](BUSINESS_WORKFLOWS.md)。

## 1. 当前可执行

开发环境 Go 1.22+；发布版未来不要求用户安装 Go，但被编排项目的工具链仍必需。在仓库根可执行：

```powershell
go version
go build ./...
go vet ./...
go test ./...
go run ./cmd/prfrail version
go run ./cmd/prfrail init --workspace .
go run ./cmd/prfrail validate --chain .\proofrail.chain.json
go run ./cmd/prfrail config explain --chain .\proofrail.chain.json
go run ./cmd/prfrail run --chain .\proofrail.chain.json --run-id run-demo
go run ./cmd/prfrail report --run-dir .\tmp\prfrail-runs\run-demo
```

核心包已有自动化测试，但完整产品 E2E/TUI 仍在后续切片。当前 CLI 输出默认无 ANSI 颜色并支持 `--json`；开发环境 gopls 与核心运行无关。系统/Git 代理不保证 Go 使用代理；安装 Go 工具需要进程 HTTP_PROXY/HTTPS_PROXY，使用后恢复原环境，不随意全局改 GOPROXY 或关闭校验。

## 2. 安装与首次使用设计

未来安装顺序：选择 Windows amd64 对应发行包→验证签名和校验和→解压到独立目录→验证 version→静态预览→单独授权能力探测。Linux 在 S1 仅核心 CI/预览，不宣传正式支持。未知签名或平台不匹配不得执行。升级前备份对象、事件和状态引用的完整闭包并只读核验历史 run；不覆盖运行中的 host。

当前可复制流程为 `prfrail init --workspace <path>`、`prfrail validate --chain <chain-file>`、`prfrail config explain --chain <chain-file>`、`prfrail run --chain <chain-file> [--run-id] [--run-dir]`、`prfrail report --run-dir <run-dir>`。`run` 当前仅支持 noop-only 链；执行型 step 的真实 gate/adapter 闭环仍在后续切片。没有 AI 时先使用文件队列 fixture consumer 完成闭环；VS Code 不是必需依赖。

当前 CLI 基线使用单文件链配置（优先 `--chain`；否则按 `./proofrail.chain.json` → `./proofrail.json` 搜索）。`proofrail.toml/workspace.toml` 的分层配置仍属后续规划。运行状态由引擎维护，用户不得修改 runtime-state、journal、receipt 或已接受 manifest。使用 config explain 检查最终来源、工具/网络权限和预算，不从模型自然语言推断配置。

首次运行前确认源树/运行区/store 不重叠，secret 排除、磁盘、工具链、模型/费用上限、review 主体和可实施隔离。源树含未提交文件不要求 reset；baseline 捕获后不自动吸收外部变化。不要把 original workspace 交给 agent 编辑。

## 3. 操作员决策表

| 观察 | 允许动作 | 禁止动作 |
|---|---|---|
| REVIEW_PENDING | 核对父/候选/证据摘要、代码和文档 diff，approve/reject 或有授权的 waive | 把技术通过当接受、改 receipt |
| 阻断门禁失败 | 保存证据，确认停机，有限修复或取消 | 跳过失败继续下一任务 |
| 进程存活不明 | 只读检查受管身份、终态和租约，保持暂停 | 按 PID 名称随意 kill、提前恢复文件 |
| 磁盘不足 | 暂停、GC dry-run 检查无引用过期对象或增加配额 | 删除 baseline/已接受/失败证据取空间 |
| adapter 超时 | 核对持久投递记录/requestId/结果不确定性 | 改 ID 无限重发、切到 visible/auto |
| 预算耗尽 | 人工复核失败指纹、费用与新授权 | 重置计数规避 hard_block |
| store 损坏 | 隔离对象并离线核验引用，保持依赖链阻断 | 用新哈希覆盖旧摘要冒充恢复 |

## 4. 人工交接和恢复

规划操作为 `prfrail handoff open/complete/abort/request-agent`。open 只在受管写者停机、journal 刷新、manifest 固化后创建唯一人工写租约。操作员读取 handoff pack，确认目标、允许范围、复现命令、预算、期限和归还条件；只编辑指定 run-workspace。

complete 归还而非批准：重新扫描、检查越界/秘密/未知进程、执行全部规定 gates，然后独立 review。abort 保留证据并丢弃本次候选；request-agent 只交回经确认结论并开启新 attempt。断线/到期保持暂停，不擅自把权限交给代理。密码/token/MFA 直接输入安全终端，禁止贴入聊天或 issue。

崩溃恢复顺序：保存现场→确认受管进程状态和旧写者失效→验证 store/事件→重放 journal→重建状态→检查候选是否具有完整接受依据→显式继续。不能因 state 文件说 PASSED 就越过核验；pause/cancel 应通过控制 API，禁止手改文件模拟。

## 5. 排错最小证据包

提交问题时附版本/OS、run/task/step/attempt、首个错误类别、requestId（如涉及代理）、输入与证据哈希、最小复现、已跑命令/退出码、预期/实际和脱敏日志摘要。不要附 token、私钥、完整聊天历史或整仓源码。

| 症状 | 首查模块/证据 | 最小检查 |
|---|---|---|
| 配置拒绝 | taskdef/字段来源 | 对应 AT-01 fixture |
| 无法捕获/恢复 | snapshot/路径清单 | AT-02，同平台路径/权限 |
| 原子应用不明 | applier/journal | AT-04，单崩溃点重放 |
| 命令挂起 | guard/gates/进程身份 | AT-06，不先改 timeout |
| 投递后无回复 | adapters/持久请求、目标实例 | AT-08，先离线后实机 |
| 评审后未推进 | chain/候选绑定和 publication receipt | AT-05，检查批准是否过期/摘要改变 |
| 代码通过文档阻断 | impact rule/生成关系/文档 hook | AT-11，不删规则 |

## 6. whois 影子迁移方案

whois v3.4.0 已封板，本轮及 S0 不改其生产代码、任务文件或发布目录。用户批准后，以只读、脱敏、固定摘要导出夹具；不得在 A/B 运行期间回滚源文件。

| whois 示例概念 | ProofRail 映射 | 保留证据 |
|---|---|---|
| A/B | T1/T2 有序 task | 输入/输出快照摘要 |
| D/V 轮次 | step label + code/build/verify/noop | 显式 kind、原因与结果 |
| start-file | chain-file + 环境配置 + 机器 state 分离 | 映射差异表，不声称接口兼容 |
| Step47/C 编译黄金 | C harness hooks | 固定输入、退出分类、产物摘要 |
| 自愈候选 | Prepare/Inspect/Validate/Promote | 旧新失败分类和提升证据 |

S0 用小型 Go 非 C 夹具和 C+Go/三语言配置反证核心耦合。S1 仅在隔离副本重放，比较同输入、同判定、同失败分类；差异需要解释和测试，不要求非确定日志逐字相同。只有 S1 exit、回滚演练和用户明确授权后才讨论切换；ProofRail 不能自动替换已封板流程。

## 7. 自托管、升级和发布

bootstrap：维护者用常规 Go 构建与独立测试、人工审计建立首个 seed；记录源码摘要、工具链、oracle、签名/批准主体。不得把当前占位版本直接当可信 seed。

两代验收：稳定 N 捕获隔离输入并构建 N+1；候选不覆盖 host；新进程+独立 run/store 重放固定验收链；独立 oracle 比较 Schema、receipt、恢复和 fail-close goldens，完成 clean-room 演练。失败只阻断候选，不修改 seed 或历史证据。

发布清单：RFC/Schema/版本一致；build/vet/test、原生平台和全部强制 AT；依赖许可证/漏洞清单；二进制对应平台/CGO_ENABLED=0；校验和+已批准签名；升级/回滚测试；双语说明与已知限制；外部发布门禁和用户当轮 commit/push/release 授权。默认远端 origin，Gitee 镜像不自动推送。

回退：停止候选进程并保留证据，只在明确存储兼容时恢复旧二进制；未知格式保持只读，不能降级覆写。签名、支持周期和漏洞响应提案见 [ADR_REGISTER.md](ADR_REGISTER.md)。

## 8. 完整产品操作流程（规划，RFC §19）

| 环节 | 操作与成功依据 | 失败时 |
|---|---|---|
| PC-01 试用 | 随附无 AI 夹具做静态预览，检查权限/预算/unknown；单独授权命令探测后才运行 | 不执行配置中的 hook；缺能力按提示准备，不自动装工具 |
| PC-02 交付 | 选已接受结果，指定与源/run/store 不重叠的新目录，核验包摘要/排除项/评审引用 | 拒绝覆盖；中断包不算交付；导出不等于外部产品发布 |
| PC-03 审批/撤销 | 查看待审批队列与候选摘要；撤销后确认无新投递及在途停机证据 | 不能停机则维持暂停；不通过 waive 放宽硬门禁 |
| PC-04 恢复/诊断 | 先只读获取脱敏恢复计划，确认最后接受点、journal、进程和未知副作用 | 不手删锁、不盲目重投；上传/部署不能靠文件恢复撤销 |
| PC-05 费用 | 查看预留/已结算/未知占用及计量来源；超限申请新授权 | 超时未必免费；订阅模式只保证授权调用/token 上限，不保证货币账单 |
| PC-06 停用/卸载 | 停受管写者、撤销授权，核验完整备份并在新 store 演练；默认保留 run/store | 未知格式保留原件；不删除共享 SessionBridge/工具链；数据删除另确认 |

S1 的交付包由用户自行决定后续合并/发布，ProofRail 不自动向原树写文件或执行 Git。安装包支持状态、SBOM/许可证、签名身份和撤销新鲜度分别检查；离线未知撤销不应显示“最新可信”。备份不包含凭据，需另行通过安全方式配置恢复环境。

历史证据含秘密：隔离、停止外发、撤销凭据，再由人决定审计保留与删除冲突；不改旧哈希伪造完整链。删除不能保证 SSD/外部副本安全擦除。报告默认本地无遥测；提交支持材料只提供复现所需脱敏最小包。以上能力待 T019–T024 实现，当前不能作为可执行命令教程。
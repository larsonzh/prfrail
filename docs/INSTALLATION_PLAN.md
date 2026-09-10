# ProofRail 安装与部署规划

[English](INSTALLATION_PLAN_EN.md)

日期：2026-09-09。状态：S1 安装模型已决定，正式发行仍待批准。当前没有正式发行包或受支持的生产部署；本文件记录已冻结模型、待决事项和发布前门禁，实际步骤见 [便携 ZIP 安装指南](INSTALLATION.md)。

## 1. 当前可用方式

当前仅支持开发者从源码构建和运行测试：Go 1.22+、仓库源码及目标项目自身工具链。`cmd/prfrail` 仍是开发中入口，不能据此宣称已安装产品。准确的当前命令见 [操作手册](OPERATIONS.md)。

## 2. 目标部署模型

S1 目标是本地单用户、单写者、无常驻云服务的 CLI/TUI：ProofRail 二进制、状态/store、隔离 run-workspace 和只读源目录边界明确。正式执行采用固定 CLI Agent/AgentRunner；SessionBridge 只在辅助消息和显式 visible 黑箱候选中使用。黑箱路径以快速复用现有 Agent 换取较低过程保证，须展示风险/责任并由用户监督。当前发行候选已交付逐行 CLI 与 T025 聚焦终端交互收件箱；两条 AI 执行路径和完整统一 TUI 均未实现。

## 3. 待决决策

| 决策 | 候选 | 冻结条件/责任 | 当前状态 |
|---|---|---|---|
| S1 发行载体 | Windows amd64 便携 ZIP，手工解压；后续可另行评估 winget | T017 原生验收、固定 commit、SBOM、许可证清单和校验和完成 | 2026-09-09 所有者批准 |
| Linux 支持 | S1 核心 CI/预览；正式支持候选在 S2 | 原生兼容、路径/权限/升级测试通过 | S1 范围已定 |
| 安装目录与 PATH | 用户选择独立目录；完整路径调用，不修改 PATH/注册表/系统目录 | 不需管理员权限、可并存升级和回滚 | 2026-09-09 所有者批准 |
| state/store 默认位置 | 用户数据目录或显式 workspace 配置 | 不与源/run 嵌套，权限与备份策略可验证 | 待定 |
| 配置搜索顺序 | 显式参数、workspace、用户 profile | T016 冻结 CLI 与 config explain 行为 | 待定 |
| 发布信任 | S1 使用固定 commit、GitHub CI、SHA256SUMS、SBOM 和许可证清单 | 校验和只证明内容完整性；签名/attestation 后续增强 | S1 基线已定 |
| 更新方式 | 手工下载核验，新版本解压到新目录并显式切换路径；未来包管理器另议 | 支持矩阵、schema 兼容、备份恢复门禁 | S1 模型已定；最终 ZIP 待复验 |
| 卸载与数据保留 | 手工删除版本目录；run/store 默认保留，数据删除单独授权 | T024 生命周期记录和显式删除授权 | S1 模型与库级规则已冻结 |
| CLI Agent/AgentRunner | 核心离线模式可选；正式 AI 执行必需 | 固定版本/摘要、能力矩阵、隔离、停止/恢复、证据与独立验收通过 AT-23 | T026/T027 待实现 |
| VS Code visible 黑箱候选 | 可选降级保证；要求 VS Code 1.82+、可用 Copilot Chat、SessionBridge 0.1.1 | 显式风险确认、`visible`/非 legacy、人工监督/归还、后置全扫与 AT-24；禁止自动 fallback | T028 待实现 |

S1 不指定默认安装目录，不发布一键安装脚本，不修改 PATH，不声明自动更新，也不要求管理员权限。state/store 默认位置、正式下载入口与 EOL 通知仍待冻结。

## 4. 已完成的候选便携包演练

2026-09-09 已对两个受信 Windows amd64 CI artifact 完成离线 verifier 核验及显式临时目录内的首次运行、升级、回滚和保留数据卸载。旧/新完整 candidate SHA 分别为 `1e7af676e8e84028c7ecfef1ccde91213738c20c` 与 `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`；详细证据见 [S1 Windows 便携包安装演练](validation/s1-install-drill.md)。

该结果关闭“候选 artifact 是否能执行生命周期”的工程疑问，并支撑 2026-09-09 获批的手工便携 ZIP、不改 PATH 模型。正式下载入口和最终 ZIP 尚未发布，仍须用最终载体重跑。

## 5. 规划的用户路径

1. 从受信发布页选择与 OS/架构匹配的发行物。
2. 独立核验固定 commit、GitHub CI 结果、SHA256SUMS、SBOM/许可证清单和支持状态；明确校验和不认证发布者身份。
3. 解压/安装到独立目录，不覆盖运行中的 host；执行版本与只读自检。
4. 初始化显式 state/store 与 workspace 配置，确认它们不和源目录或 run-workspace 重叠。
5. 先运行无副作用静态预览，再单独授权能力探测和实际运行。
6. 升级前停止写者并建立完整、可恢复的新 store 备份；失败时保留旧二进制和原 store。
7. 卸载时默认保留证据和配置；数据删除是单独授权动作。

可复制命令见 [便携 ZIP 安装指南](INSTALLATION.md)。T024 已提供库级完整备份、新 store 核验、默认保留和显式删除授权规则；正式发行文件名和下载地址仍须在发行批准时绑定。

## 6. 发布前安装门禁

- Windows 11 amd64 原生安装、首次运行、升级、回滚和卸载通过；Linux 仅按当期支持矩阵声明。
- 发布物可离线核验 SHA256SUMS、SBOM、许可证和精确版本；这些证据不认证发布者身份。后续若增加签名/attestation，未知或撤销信息过期不得显示 fully verified。
- 全新用户环境无需 Go 即可运行 ProofRail；项目任务需要的外部工具明确列出且不自动安装。
- 安装不修改目标源码或 Git，不删除 SessionBridge/用户工具链，不把 secret 写入配置、日志或备份。
- state/store 备份覆盖完整引用闭包，并已在新目录完成恢复演练。
- README 与操作文档中的每条安装命令均由对应发行物在受支持平台实测。

## 7. 文档完成条件

T016 已冻结 CLI、配置和当前源码运行方式；T024 已冻结库级备份、升级前阻断、停用和数据保留规则；T017 的候选 CI artifact、固定 commit、SHA256SUMS、SBOM/许可证门禁与双平台可信矩阵已完成。签名/attestation 是后续增强。S1 安装模型已于 2026-09-09 获批；最终 ZIP 复验和 S1 exit 获批前，本文仍不构成生产发布声明。[操作手册](OPERATIONS.md) 继续区分当前可执行与未来设计；候选发布边界见 [支持矩阵](S1_SUPPORT_MATRIX.md) 与 [发行说明](S1_RELEASE_NOTES.md)。
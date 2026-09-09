# ProofRail S1 候选发行说明

[English](S1_RELEASE_NOTES_EN.md)

日期：2026-09-09。状态：`DRAFT / NOT RELEASED`。版本、下载链接、正式 artifact 文件名、支持起始日期和批准人均待 S1 exit 冻结；本文件不授权 tag、GitHub Release、发布或安装。

## 候选能力

- 本地单用户、单写者 CLI 基线：`init`、`validate`、`config explain`、`preview`、`approvals`、`run`、`report`。
- 默认拒绝未授权能力；静态 preview 不运行命令、网络、模型或凭据探测。
- noop-only 任务链可完成，executable `code/build/verify` 步骤 fail-close。
- 证据、快照、进程/租约守卫、托管变更集、事务应用、授权撤销、备份/恢复和默认保留数据的卸载规则。
- generic/C/Go 最小样例与脱敏 whois 冻结夹具只读影子验证。
- 固定 verifier/candidate commit、双平台 candidate probe/self-host、校验和、SPDX SBOM 与许可证清单的候选发布链。

## 支持与已知限制

支持范围以 [S1 候选支持矩阵](S1_SUPPORT_MATRIX.md) 为准。S1 安装模型是[手工解压 Windows 便携 ZIP](INSTALLATION.md)、显式路径运行且不修改 PATH；当前仍无正式发行包，候选包仅在隔离临时目录完成生命周期演练。S1 不提供签名/attestation、Linux 正式支持、自动更新、Web 控制台、完整 TUI 或任意 executable hook。

## 发布前剩余门禁

1. 对最终 ZIP 复跑核验、首次运行、升级、回滚和卸载。
2. 将版本、artifact 文件名、下载位置、EOL 通知方式和支持日期绑定到本说明及支持矩阵。
3. 获得独立评审结论和产品所有者明确 S1 exit/发行批准。
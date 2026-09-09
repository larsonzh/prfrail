# ProofRail S1 候选支持矩阵

[English](S1_SUPPORT_MATRIX_EN.md)

日期：2026-09-09。状态：`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。

| 范围 | S1 候选状态 | 约束 |
|---|---|---|
| Windows 11 amd64 | 正式支持候选 | 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期 |
| Linux amd64 | 核心 CI | 仅构建、测试、候选探针和自托管；不宣称 S1 生产支持 |
| macOS / Windows on Arm / Linux arm64 | 不支持 | 未进入 S1 验收矩阵 |
| generic/C/Go harness | 配置与 noop 编排 | language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close |
| VS Code / SessionBridge | 可选适配器 | 不是核心安装依赖，不扩大 ProofRail 权限 |

## 版本与安全维护候选政策

- 当前稳定 minor 与前一个 minor 接收安全修复；前一个 minor 的候选维护窗口为 90 天。
- 历史未知 run 只读保留，不自动迁移。
- 首版正式发布前仍须冻结 EOL 通知方式和例外；本页不能替代最终公告。

## 发行信任边界

候选发行要求完整 commit SHA、GitHub CI、`SHA256SUMS`、SPDX SBOM 和许可证清单。校验和证明内容与清单一致，不认证发布者身份；S1 不宣称签名或 attestation。安装模型见 [便携 ZIP 安装指南](INSTALLATION.md)；具体文件名、下载位置和支持起始日期须与最终发行说明一起批准。
# ProofRail S1 候选支持矩阵

[English](S1_SUPPORT_MATRIX_EN.md)

日期：2026-09-09。状态：`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。

| 范围 | S1 候选状态 | 约束 |
|---|---|---|
| Windows 11 amd64 | 正式支持候选 | 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期 |
| Linux amd64 | 核心 CI | 仅构建、测试、候选探针和自托管；不宣称 S1 生产支持 |
| macOS / Windows on Arm / Linux arm64 | 不支持 | 未进入 S1 验收矩阵 |
| generic/C/Go harness | 配置与 noop 编排 | language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close |
| CLI Agent / AgentRunner | 规划、未交付 | 目标正式 AI 执行通道；须固定版本/配置、能力探针、隔离 workspace、进程/会话证据与独立验收；T026/T027/AT-23 阻断 S1 |
| VS Code / Copilot Chat / SessionBridge visible 黑箱候选 | 规划、降级保证、未交付 | 核心离线 CLI 不要求；仅在用户显式确认风险并监督时桥接隔离候选，归还后独立验收；不保证工具/网络/费用/外部副作用，T028/AT-24 阻断 S1 |

## 版本与安全维护候选政策

- 产品版本采用 `vMAJOR.MINOR.PATCH`，首版候选为 `v0.1.0`；patch 发布不启动新的 EOL 倒计时。
- 当前稳定 minor 与前一个 minor 接收安全修复；前一个 minor 的候选维护窗口为 90 天。
- 历史未知 run 只读保留，不自动迁移。
- 待审批的 ADR-010 提议：新 minor 发布时在 GitHub Release notes 与本支持矩阵同时公告前一 minor 的 EOL 日期，日期为发布后第 90 个自然日。
- 90 天是原型期初始窗口；每次 minor 发布前按稳定性、用户量和维护能力评审，仅可延长后续新公告的窗口。
- 任何偏离须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处，且不得追溯缩短已公告窗口。该提议尚未获批，本页不能替代最终公告。

## 发行信任边界

候选发行要求完整 commit SHA、GitHub CI、`SHA256SUMS`、SPDX SBOM 和许可证清单。校验和证明内容与清单一致，不认证发布者身份；S1 不宣称签名或 attestation。安装模型见 [便携 ZIP 安装指南](INSTALLATION.md)；具体文件名、下载位置和支持起始日期须与最终发行说明一起批准。
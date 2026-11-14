# ProofRail S1 候选支持矩阵

[English](S1_SUPPORT_MATRIX_EN.md)

日期：2026-09-09。状态：`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效；当前不是已发布 SLA。

| 范围 | S1 候选状态 | 约束 |
|---|---|---|
| Windows 11 amd64 | 正式支持候选 | 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期；基础终端（逐行 CLI / `interactions tui` / `--progress`）保留；完整一体化 TUI 移出 S1、由 S2 承接 |
| Linux amd64 | 核心 CI | 仅构建、测试、候选探针和自托管；不宣称 S1 生产支持 |
| macOS / Windows on Arm / Linux arm64 | 不支持 | 未进入 S1 验收矩阵 |
| generic/C/Go harness | 配置与 noop 编排 | language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close；Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力 |
| CLI Agent / AgentRunner | 规划、未交付 | 目标正式 AI 执行通道；须固定版本/配置、能力探针、隔离 workspace、进程/会话证据与独立验收；T026/T027/AT-23 阻断 S1 |
| VS Code / Copilot Chat / SessionBridge visible 黑箱候选 | 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布 | 核心离线 CLI 不要求；仅在用户显式确认风险并监督时桥接隔离候选，归还后独立验收；不保证工具/网络/费用/外部副作用，AT-24 条件闭合，S1 仍因 T027/AT-23 阻断 |

## 版本与安全维护候选政策

- 产品版本采用 `vMAJOR.MINOR.PATCH`，首版候选为 `v0.1.0`；patch 发布不启动新的 EOL 倒计时。
- 当前稳定 minor 与前一个 minor 接收安全修复；前一个 minor 的候选维护窗口为 90 天。
- 历史未知 run 只读保留，不自动迁移。
- 待审批的 ADR-010 提议：新 minor 发布时在 GitHub Release notes 与本支持矩阵同时公告前一 minor 的 EOL 日期，日期为发布后第 90 个自然日。
- 90 天是原型期初始窗口；每次 minor 发布前按稳定性、用户量和维护能力评审，仅可延长后续新公告的窗口。
- 任何偏离须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处，且不得追溯缩短已公告窗口。该提议尚未获批，本页不能替代最终公告。

## 发行信任边界

候选发行要求完整 commit SHA、GitHub CI、`SHA256SUMS`、SPDX SBOM 和许可证清单。校验和证明内容与清单一致，不认证发布者身份；S1 不宣称签名或 attestation。安装模型见 [便携 ZIP 安装指南](INSTALLATION.md)；具体文件名、下载位置和支持起始日期须与最终发行说明一起批准。

> 切片同步（P55 `EXIT-SEED-BUILD-EVIDENCE`）：本片不适用（无平台与版本支持范围增减）；理由：补证只归档既有 seed 驱动构建的两平台捕获，平台清单、支持层级与 EOL 政策不变。

> 切片同步（P57 `P5-P6-CLOSEOUT-ASSESSMENT`）：本片不适用（无平台清单、支持层级、EOL 政策或安装模型的增减）；理由：P57 只把 T018 剩余门整理为收口清单，未改变任何门的判定，本文件措辞逐字不动。

> 切片同步（P71 `S1-SCOPE-ANNOTATE`）：TUI 移出与 visible 降级保证口径已同步；理由：R-07 b1 附条件 2 落地只改 REQ 阶段归属与调试指引——平台清单与支持层级不变，TUI 移出系 P68 已裁定的**功能阶段归属同步**，Windows 行状态仍为「正式支持候选」，visible 行仍为「降级保证」。
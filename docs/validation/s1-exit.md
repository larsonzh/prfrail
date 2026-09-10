# S1 Exit 验收报告

[English](s1-exit_EN.md)

日期：2026-09-09。结论：`BLOCKED`。T018 已完成 AT-15 的只读 whois 冻结夹具对比和 generic/C/Go 最小闭环，可信 Windows artifact 的候选安装生命周期演练也已通过；手工解压便携 ZIP、不修改 PATH 的安装模型已获产品所有者批准，并已有双语教程、支持矩阵与发行说明草案。但最终 ZIP 尚未绑定和复验，EOL 通知方式尚未冻结，产品所有者也未批准最终 S1 exit，因此不得把部分通过标记为 S1 完成。

## AT-15 证据

1. 本轮得到仅限 `D:\LZProjects\whois\testdata\`、仅非秘密已脱敏夹具、仅本轮有效的只读检查与导出授权；未读取活动 start-file、日志、`out/` 或 `release/`，未修改或执行 whois。
2. `examples/whois-shadow/` 从 `testdata/cidr_matrix_cases_draft.tsv` 导出 9 个固定用例。源 SHA-256 为 `8a6abecb0178a970677b84b3e5403ce8830f281e1a197d792ae54089b3e9efb0`，1,981 字节、10 行，秘密标记扫描 0 命中。
3. stdlib-only 比较器按 case ID 比较 input、result 与 failure classification；9/9 正例通过，输入漂移、结果漂移和失败分类漂移三个负例均 fail-close。
4. `examples/{generic-minimal,c-minimal,go-minimal}/` 分别声明 `generic-standard`、`c-standard`、`go-standard` language scope；三组 CLI `validate → preview → run → report` 均在临时 run 目录闭环，状态为 `COMPLETED`，零命令/网络/模型/凭据/版本探测调用，配置文件字节不变。
5. 边界：whois 结果是获授权冻结夹具的离线影子投影，不是在线查询、生产执行或切换证据；三语言样例使用当前 noop-only CLI，不冒充 executable hook 或编译器已执行。

## AT-01–AT-24 汇总

| AT | 状态 | 主要证据 |
|---|---|---|
| AT-01 | PASS | `DEV_PLAN.md` T002/T003：29 Schema、82 fixture、JCS 向量 |
| AT-02 | PASS | [T006 snapshot](t006-snapshot.md) |
| AT-03 | PASS | [T007 process/lease](t007-process-lease.md)、[T009 chain](t009-chain.md) |
| AT-04 | PASS | [T005 evidence](t005-evidence.md)、[T008 applier](t008-applier.md) |
| AT-05 | PASS | [T011 review/publish](t011-chain.md) |
| AT-06 | PASS | [T010 gates](t010-gates.md) |
| AT-07 | PASS | [T012 tickets/repair](t012-tickets-repair.md) |
| AT-08 | PASS | [T013 adapters](t013-adapters.md) |
| AT-09 | PASS | [T014 handoff](t014-handoff.md) |
| AT-10 | PASS | [T014 handoff](t014-handoff.md) |
| AT-11 | PASS | [T015 harness/docs](t015-taskdef-documentation.md) |
| AT-12 | PASS | [T016 CLI](t016-cli-console.md) |
| AT-13 | PASS | [T017 self-host](s1-selfhost.md)，含 Node 24 run `34325758548` |
| AT-14 | PASS | [T017 self-host](s1-selfhost.md)，Windows/Linux 8/8 jobs 与 6 artifacts |
| AT-15 | PASS | `examples/whois-shadow/` 与三个 `*-minimal` 语言样例；本报告聚焦验证 |
| AT-16 | PASS | [T019 preview](t019-preview.md) |
| AT-17 | PASS | [T020 export](t020-export.md) |
| AT-18 | PASS | [T021 approvals](t021-approvals.md) |
| AT-19 | PASS | [T022 effects/diagnostics](t022-effects-diagnostics.md) |
| AT-20 | PASS | [T023 cost ledger](t023-cost-ledger.md) |
| AT-21 | PASS | [T024 lifecycle](t024-lifecycle.md) |
| AT-22 | BLOCKED | T025 尚未实现：缺少 `operator-interaction` Schema/golden、TUI 交互收件箱、结构化响应及恢复测试 |
| AT-23 | BLOCKED | T026/T027 尚未实现：缺少 AgentRunner 契约、真实 CLI Agent 能力探针、adapter、停止/恢复和独立验收闭环 |
| AT-24 | BLOCKED | T028 尚未实现：缺少 visible 黑箱候选的风险确认、隔离投递、显式归还、全量扫描及 reduced/unknown 报告 |

这些 PASS 表示对应 AT 的现有确定性验收，不自动满足阶段文档与发布门禁。

## RFC §14 S1 Exit

| 门禁 | 状态 | 证据/缺口 |
|---|---|---|
| 核心任务链与 fail-close 场景 | BLOCKED | 既有 AT-01–AT-14、AT-16–AT-21 通过；T025–T028/AT-22–AT-24 尚未实现 |
| seed 自托管与真实平台 | PASS | T017 Windows/Linux 可信矩阵及 Windows 本机两代演练 |
| whois 只读影子 | PASS | 9/9 固定用例与 3 类漂移负例；不含生产切换 |
| generic/C/Go harness | PASS | 三套 harness 定义及三个最小用户样例均通过只读 CLI 闭环 |
| 正式安装包教程 | PARTIAL | [便携 ZIP 安装指南](../INSTALLATION.md)与模型已冻结；最终 ZIP 尚未绑定和复验 |
| 安装/升级/回滚/卸载实测 | PARTIAL | [Windows 候选便携包演练](s1-install-drill.md)通过；最终载体确定后仍须复验 |
| 支持矩阵与发行说明 | PARTIAL | [候选支持矩阵](../S1_SUPPORT_MATRIX.md)与[发行说明草案](../S1_RELEASE_NOTES.md)已建立，尚未绑定正式发行物 |
| 独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审 |
| 产品所有者最终 S1 批准 | BLOCKED | 本轮只授权 whois 夹具读取，不等于 S1 发布/阶段批准 |

## 剩余风险与下一步

1. 手工解压便携 ZIP、用户自选独立目录、显式路径运行且不修改 PATH 的模型已冻结；正式 ZIP 文件名和下载入口仍待发行绑定。EOL 通知与例外已有 ADR-010 待审批方案，但尚未获产品所有者批准。
2. Windows 候选 artifact 已完成首次运行、升级、回滚和保留数据卸载；正式载体确定后必须复验。Linux 只声明核心 CI，不宣称 S1 正式支持。
3. generic/C/Go 样例只证明语言作用域和 noop 编排；当前 CLI 对 executable step 仍 fail-close。
4. 完成 T025/AT-22 的 TUI 人工交互、T026/T027/AT-23 的 AgentRunner 完整执行，以及 T028/AT-24 的 visible 黑箱候选降级保证；再用最终 ZIP 复跑教程并绑定支持矩阵/发行说明，最后由产品所有者明确批准 S1 exit。在此之前 T018 与 S1 状态保持未完成。
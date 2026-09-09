# S1 Exit 验收报告

[English](s1-exit_EN.md)

日期：2026-09-09。结论：`BLOCKED`。T018 已启动，AT-15 的只读 whois 冻结夹具对比与非 C Go 最小闭环均通过；但正式安装包教程、可运行 generic/C 用户样例、安装/升级/回滚/卸载实测和产品所有者最终 S1 批准尚未完成，因此不得把部分通过标记为 S1 完成。

## AT-15 证据

1. 本轮得到仅限 `D:\LZProjects\whois\testdata\`、仅非秘密已脱敏夹具、仅本轮有效的只读检查与导出授权；未读取活动 start-file、日志、`out/` 或 `release/`，未修改或执行 whois。
2. `examples/whois-shadow/` 从 `testdata/cidr_matrix_cases_draft.tsv` 导出 9 个固定用例。源 SHA-256 为 `8a6abecb0178a970677b84b3e5403ce8830f281e1a197d792ae54089b3e9efb0`，1,981 字节、10 行，秘密标记扫描 0 命中。
3. stdlib-only 比较器按 case ID 比较 input、result 与 failure classification；9/9 正例通过，输入漂移、结果漂移和失败分类漂移三个负例均 fail-close。
4. `examples/go-minimal/` 声明独立 `go-standard` language scope；CLI 的 `validate → preview → run → report` 在临时 run 目录闭环，状态为 `COMPLETED`，零命令/网络/模型/凭据/版本探测调用，配置文件字节不变。
5. 边界：whois 结果是获授权冻结夹具的离线影子投影，不是在线查询、生产执行或切换证据；Go 样例是当前 noop-only CLI 能力，不冒充 executable hook 已实现。

## AT-01–AT-21 汇总

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
| AT-15 | PASS | `examples/whois-shadow/` 与 `examples/go-minimal/`；本报告聚焦验证 |
| AT-16 | PASS | [T019 preview](t019-preview.md) |
| AT-17 | PASS | [T020 export](t020-export.md) |
| AT-18 | PASS | [T021 approvals](t021-approvals.md) |
| AT-19 | PASS | [T022 effects/diagnostics](t022-effects-diagnostics.md) |
| AT-20 | PASS | [T023 cost ledger](t023-cost-ledger.md) |
| AT-21 | PASS | [T024 lifecycle](t024-lifecycle.md) |

这些 PASS 表示对应 AT 的现有确定性验收，不自动满足阶段文档与发布门禁。

## RFC §14 S1 Exit

| 门禁 | 状态 | 证据/缺口 |
|---|---|---|
| 核心任务链与 fail-close 场景 | PASS | AT-01–AT-14、AT-16–AT-21 报告 |
| seed 自托管与真实平台 | PASS | T017 Windows/Linux 可信矩阵及 Windows 本机两代演练 |
| whois 只读影子 | PASS | 9/9 固定用例与 3 类漂移负例；不含生产切换 |
| generic/C/Go harness | PARTIAL | harness 定义与 Go 用户样例存在；generic/C 可运行用户样例缺失 |
| 正式安装包教程 | BLOCKED | `INSTALLATION_PLAN.md` 仍明确为规划，当前无正式发行包 |
| 安装/升级/回滚/卸载实测 | BLOCKED | T024 只有库级生命周期能力，无安装入口实测 |
| 支持矩阵与发行说明 | PARTIAL | 平台/版本政策已冻结，尚无与正式发行物绑定的最终文档 |
| 独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审 |
| 产品所有者最终 S1 批准 | BLOCKED | 本轮只授权 whois 夹具读取，不等于 S1 发布/阶段批准 |

## 剩余风险与下一步

1. 发行载体、安装目录和 PATH 尚未冻结；不得制作虚假安装命令或发布未授权 artifact。
2. Windows 安装、首次运行、升级、回滚和卸载必须对最终发行物实测；Linux 只声明核心 CI，不宣称 S1 正式支持。
3. generic/C 用户样例需在隔离临时目录运行并证明源树不变；当前 CLI 对 executable step 仍 fail-close。
4. 完成上述证据并更新支持矩阵/发行说明后，须由产品所有者明确批准 S1 exit；在此之前 T018 与 S1 状态保持未完成。
# 批量辅助 ① 转录 — REQ 现状枚举与关联引用（GLM-5.3-Flash via deep-reasoner）

派单：2026-11-14；实得 5179 字符；任务限定 = 只列事实，不做判断。主控已抽验关键行（REQ-020 CN `:82` / EN `:80`、REQ-011 CN `:73` / EN `:71`、REQ 计数 28）。

---

## 面 1

### 1. REQ-020 所在行号与原文

`docs/PRODUCT_REQUIREMENTS.md:82`（单行条目，不跨行）：

```text
| REQ-020 | R20 | S0-S2 | S0 验证双/三语言配置；S1 C+Go 契约；S2 三语言实跑，失败整任务阻断 |
```

`docs/PRODUCT_REQUIREMENTS_EN.md:80`（单行条目，不跨行）：

```text
| REQ-020 | R20 | S0-S2 | S0 dual/triple-language configs; S1 C+Go contracts; S2 three-language execution with whole-task failure |
```

### 2. REQ-020「阶段」列当前值

- `docs/PRODUCT_REQUIREMENTS.md:82`：`S0-S2`
- `docs/PRODUCT_REQUIREMENTS_EN.md:80`：`S0-S2`

### 3. REQ-020「可验收要求」列当前值

- CN `:82`（第 4 列）：`S0 验证双/三语言配置；S1 C+Go 契约；S2 三语言实跑，失败整任务阻断`
- EN `:80`（第 4 列）：`S0 dual/triple-language configs; S1 C+Go contracts; S2 three-language execution with whole-task failure`

### 4. REQ-011 / REQ-022 / PC-01 / PC-05 / PC-06 行号与原文

REQ-011：

```text
docs/PRODUCT_REQUIREMENTS.md:73    | REQ-011 | R11 | S1/S2 | S1 generic/C/Go harness 外置，核心无语言判断分支；S2 扩展语言 |
docs/PRODUCT_REQUIREMENTS_EN.md:71 | REQ-011 | R11 | S1/S2 | External generic/C/Go harnesses in S1, no language branching in core; extend languages in S2 |
```

REQ-022：

```text
docs/PRODUCT_REQUIREMENTS.md:84    | REQ-022 | R22 | S1/S2 | 文档一等 target、确定性 impact rule、生成方向、waiver 和整体验收；S2 智能提示 |
docs/PRODUCT_REQUIREMENTS_EN.md:82 | REQ-022 | R22 | S1/S2 | First-class docs, deterministic impact rules, generation direction, waivers and aggregate acceptance; S2 intelligent hints |
```

PC-01：

```text
docs/PRODUCT_REQUIREMENTS.md:119    | PC-01 | 首次使用难判断权限、费用和准备程度 | REQ-009/010/016/027 | 无 AI 示例、无副作用预览、unknown 原因与下一步；S2 图形向导/界面本地化 | AT-16 |
docs/PRODUCT_REQUIREMENTS_EN.md:117 | PC-01 | New users cannot assess permissions, cost or readiness | REQ-009/010/016/027 | No-AI example, side-effect-free preview, unknown reasons/next action; S2 graphical onboarding/localization | AT-16 |
```

PC-05：

```text
docs/PRODUCT_REQUIREMENTS.md:123    | PC-05 | 超时/重复请求可能仍计费，难解释成本和收益 | REQ-006/026 | 持久预留/结算、未知费用占用、共享上限、无默认遥测的本地报告；S2 趋势报表 | AT-20 |
docs/PRODUCT_REQUIREMENTS_EN.md:121 | PC-05 | Timeouts/retries may be charged; costs/value are unclear | REQ-006/026 | Durable reservation/settlement, unknown holds, shared caps, local reports without default telemetry; S2 trends | AT-20 |
```

PC-06：

```text
docs/PRODUCT_REQUIREMENTS.md:124    | PC-06 | 升级、备份、停用、卸载和保留冲突未形成闭环 | REQ-019/025/028 | SBOM/支持矩阵、完整备份闭包/恢复演练、默认保留证据的卸载；S2 显式迁移 | AT-21 |
docs/PRODUCT_REQUIREMENTS_EN.md:122 | PC-06 | Upgrade, backup, retirement and retention conflicts lack a complete flow | REQ-019/025/028 | SBOM/support matrix, complete backup closure/restore drill, evidence-preserving uninstall; S2 explicit migration | AT-21 |
```

表头：

```text
docs/PRODUCT_REQUIREMENTS.md:61    | ID | RFC 来源 | 阶段 | 可验收要求 |
docs/PRODUCT_REQUIREMENTS.md:117   | 切片 | 业务缺口与用户结果 | REQ | S1 最小能力 / 后续边界 | 验收 |
docs/PRODUCT_REQUIREMENTS_EN.md:59 | ID | RFC | Stage | Testable requirement |
docs/PRODUCT_REQUIREMENTS_EN.md:115| Slice | Business gap and user outcome | REQ | S1 minimum / later boundary | Acceptance |
```

### 5. 五条关键词命中统计

CN：

| 条目 | Python | Node | 三语言 | C+Go | 语言支持 | harness |
|---|---|---|---|---|---|---|
| REQ-011:73 | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ `generic/C/Go harness` |
| REQ-022:84 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| PC-01:119 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| PC-05:123 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| PC-06:124 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |

补充事实：REQ-011:73 含片段 `无语言判断分支`、`扩展语言`；PC-06:124 含片段 `支持矩阵`——均非关键词字面。

EN：

| 条目 | Python | Node | three-language | C+Go | language support | harness |
|---|---|---|---|---|---|---|
| REQ-011:71 | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ `harnesses` |
| REQ-022:82 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| PC-01:117 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| PC-05:121 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| PC-06:122 | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |

补充事实：REQ-011:71 含片段 `language branching`、`extend languages`——均非关键词字面。

### 6. REQ 条目总数与编号连续性

```text
§4 表：REQ-001–REQ-022（行 63–84），共 22 条，连续
§5 表：REQ-023–REQ-028（行 90–95），共 6 条，连续
合计：REQ-001–REQ-028，共 28 条，连续
```

## 面 2

### 7. 「阶段」列取值闭集及行号

CN：

```text
S0: 行 79
S0-S2: 行 82
S0-S3: 行 76
S0/S1: 行 63
S1: 行 64,65,66,67,70,72,74,77,80,81
S1-S3: 行 68
S1/S2: 行 73,78,83,84
S1/S2（黑箱候选子项推迟 S2）: 行 69
S1/S3: 行 75
S1 不要求，推迟至 S2（基础终端保留）: 行 71
```

EN：

```text
S0: 行 77
S0-S2: 行 80
S0-S3: 行 74
S0/S1: 行 61
S1: 行 62,63,64,65,68,70,72,75,78,79
S1-S3: 行 66
S1/S2: 行 71,76,81,82
S1/S2 (the black-box candidate sub-item is deferred to S2): 行 67
S1/S3: 行 73
Not required by S1, deferred to S2 (the basic terminal stays): 行 69
```

输出字符数：5179
BATCH-01 DONE

---

## 主控抽验记录

- `docs/PRODUCT_REQUIREMENTS.md:82` 与 `:73` 原文一致（2026-11-14 主控机械复核）。
- REQ 计数 CN=28 / EN=28（机械复核）。

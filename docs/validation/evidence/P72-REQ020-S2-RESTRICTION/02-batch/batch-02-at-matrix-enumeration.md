# 批量辅助 ② 转录 — AT 与支持矩阵的语言相关枚举（GLM-5.3-Flash via deep-reasoner）

派单：2026-11-14；实得 7866 字符；任务限定 = 只列事实，不做判断。主控已抽验关键行（AT-01 CN/EN `:31`、AT-11 `:41`、S1_SUPPORT_MATRIX 表头 `:7` 与 harness 行 `:12`）。

---

约定声明（方法，非结论）：以下关键词扫描为**子串匹配**；语言名 `C`、`Go` 按**大小写敏感**匹配；其余关键词**不区分大小写**；命中片段照录原文。

## 面 A：TEST_STRATEGY.md / TEST_STRATEGY_EN.md

### A1. AT-01 / AT-11 / AT-20 行号与逐字全文

表头列名（定位用）：

```text
TEST_STRATEGY.md:29    | ID | 需求 | 场景与断言 | 所有者/计划 |
TEST_STRATEGY_EN.md:29 | ID | Requirements | Scenario/assertion | Owner/plan |
```

```text
TEST_STRATEGY.md:31 | AT-01 | REQ-001/010/017/018/020/022/028 | Schema 合法配置；未知字段/版本、重复 ID、空 steps、非法 noop、重叠 target、依赖/生成环拒绝；C+Go/三语言/交接/文档样例 | taskdef；T002/T003 |
TEST_STRATEGY.md:41 | AT-11 | REQ-011/013/020/022 | generic/C/Go、锁模板生成 A；自批/手改生成物/文档缺失阻断；C+Go 同 task 整体恢复 | taskdef/gates；T015 |
TEST_STRATEGY.md:50 | AT-20 | REQ-006/026 | PC-05：预留前后/投递后/结算点崩溃，超时保留余额占用，重复结算去重，两 run 竞争共享额度不超配；订阅模式不报精确账单 | tickets/adapters；T023 |
```

```text
TEST_STRATEGY_EN.md:31 | AT-01 | REQ-001/010/017/018/020/022/028 | Valid schemas; reject unknown fields/versions, duplicate IDs, empty steps, invalid noop, overlapping targets and dependency/generation cycles; dual/triple-language, handoff and doc fixtures | taskdef; T002/T003 |
TEST_STRATEGY_EN.md:41 | AT-11 | REQ-011/013/020/022 | generic/C/Go, locked template A; self-approval/manual generated edits/missing docs block; C+Go whole-task recovery | taskdef/gates; T015 |
TEST_STRATEGY_EN.md:50 | AT-20 | REQ-006/026 | PC-05: crashes before/after reservation, dispatch and settlement; timeout holds, deduplicated settlement, two runs cannot overallocate shared cap; no precise subscription-bill claim | tickets/adapters; T023 |
```

### A2. 关键词逐条标注（11 词 × 6 行）

| 关键词 | md:31 AT-01 | md:41 AT-11 | md:50 AT-20 | EN:31 AT-01 | EN:41 AT-11 | EN:50 AT-20 |
|---|---|---|---|---|---|---|
| `C+Go` | 命中「C+Go/三语言」 | 命中「C+Go 同 task」 | 未命中 | 未命中 | 命中「C+Go whole-task recovery」 | 未命中 |
| `Go` | 命中（C+Go 内） | 命中「generic/C/Go」 | 未命中 | 未命中 | 命中「generic/C/Go」 | 未命中 |
| `C` | 命中（C+Go 内） | 命中「generic/C/Go」 | 命中「PC-05」 | 未命中 | 命中「generic/C/Go」 | 命中「PC-05」 |
| `三语言` | 命中「三语言」 | 未命中 | 未命中 | 未命中 | 未命中 | 未命中 |
| `three-language` | 未命中 | 未命中 | 未命中 | 未命中（原文「dual/triple-language」，非连续子串） | 未命中 | 未命中 |
| `Python` | 未命中 | 未命中 | 未命中 | 未命中 | 未命中 | 未命中 |
| `Node` | 未命中 | 未命中 | 未命中 | 未命中 | 未命中 | 未命中 |
| `双语言` | 未命中 | 未命中 | 未命中 | 未命中（原文「dual/triple-language」） | 未命中 | 未命中 |
| `dual-language` | 未命中 | 未命中 | 未命中 | 未命中（同上） | 未命中 | 未命中 |
| `语言` | 命中（「三语言」内） | 未命中 | 未命中 | 未命中 | 未命中 | 未命中 |
| `language` | 未命中 | 未命中 | 未命中 | 命中「triple-language」 | 未命中 | 未命中 |

### A3. AT 条目总数与编号范围

```text
TEST_STRATEGY.md 与 TEST_STRATEGY_EN.md：以「| AT-」起始的表行各 24 条（line 31–54）；编号范围 AT-01–AT-24，连续无缺号。共 24 条。
```

## 面 B：S1_SUPPORT_MATRIX.md / S1_SUPPORT_MATRIX_EN.md

### B4. 表格结构

```text
S1_SUPPORT_MATRIX.md:7    | 范围 | S1 候选状态 | 约束 |
S1_SUPPORT_MATRIX.md:8    |---|---|---|
S1_SUPPORT_MATRIX_EN.md:7 | Scope | S1 candidate status | Constraints |
S1_SUPPORT_MATRIX_EN.md:8 |---|---|---|
```

### B5. 全部表格数据行（每文件 6 行，全给）

```text
S1_SUPPORT_MATRIX.md:9  | Windows 11 amd64 | 正式支持候选 | 纯 Go `CGO_ENABLED=0` 便携 ZIP；手工解压、显式路径运行、不修改 PATH；最终 ZIP 须重跑安装生命周期；基础终端（逐行 CLI / `interactions tui` / `--progress`）保留；完整一体化 TUI 移出 S1、由 S2 承接 |
S1_SUPPORT_MATRIX.md:10 | Linux amd64 | 核心 CI | 仅构建、测试、候选探针和自托管；不宣称 S1 生产支持 |
S1_SUPPORT_MATRIX.md:11 | macOS / Windows on Arm / Linux arm64 | 不支持 | 未进入 S1 验收矩阵 |
S1_SUPPORT_MATRIX.md:12 | generic/C/Go harness | 配置与 noop 编排 | language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close |
S1_SUPPORT_MATRIX.md:13 | CLI Agent / AgentRunner | 规划、未交付 | 目标正式 AI 执行通道；须固定版本/配置、能力探针、隔离 workspace、进程/会话证据与独立验收；T026/T027/AT-23 阻断 S1 |
S1_SUPPORT_MATRIX.md:14 | VS Code / Copilot Chat / SessionBridge visible 黑箱候选 | 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布 | 核心离线 CLI 不要求；仅在用户显式确认风险并监督时桥接隔离候选，归还后独立验收；不保证工具/网络/费用/外部副作用，AT-24 条件闭合，S1 仍因 T027/AT-23 阻断 |
```

```text
S1_SUPPORT_MATRIX_EN.md:9  | Windows 11 amd64 | Production-support candidate | Pure-Go `CGO_ENABLED=0` portable ZIP; manual extraction, explicit-path invocation, no PATH modification; rerun the lifecycle against the final ZIP; the basic terminal (line-oriented CLI / `interactions tui` / `--progress`) stays; the full unified TUI moves out of S1 and is taken up by S2 |
S1_SUPPORT_MATRIX_EN.md:10 | Linux amd64 | Core CI | Build, test, candidate probe, and self-host only; no S1 production-support claim |
S1_SUPPORT_MATRIX_EN.md:11 | macOS / Windows on Arm / Linux arm64 | Unsupported | Not included in the S1 acceptance matrix |
S1_SUPPORT_MATRIX_EN.md:12 | generic/C/Go harnesses | Configuration and noop orchestration | Language scopes, read-only preview, and noop-only loops are verified; executable steps still fail closed |
S1_SUPPORT_MATRIX_EN.md:13 | CLI Agent / AgentRunner | Planned, not delivered | Target formal AI execution channel; requires pinned version/config, capability probe, isolated workspace, process/session evidence and independent acceptance; T026/T027/AT-23 block S1 |
S1_SUPPORT_MATRIX_EN.md:14 | VS Code / Copilot Chat / SessionBridge visible black-box candidate | Product side implemented (reduced assurance), real CLI end-to-end depends on T027-B, not yet a runnable release | Not needed for core offline CLI; bridges an isolated candidate only with explicit risk acknowledgment and human supervision, followed by independent acceptance; does not warrant tools/network/cost/external effects; AT-24 conditionally closed, S1 still blocked by T027/AT-23 |
```

harness 相关行事实：两文件的 harness 行均为 line 12 单行，同一行同时覆盖 generic / C / Go；两文件内均无单独的「C harness」行或「Go harness」行。

### B6. 数据行逐行关键词标注

- `S1_SUPPORT_MATRIX.md:9`：命中 `支持`「正式支持候选」；其余未命中。
- `S1_SUPPORT_MATRIX.md:10`：命中 `支持`「不宣称 S1 生产支持」；其余未命中。
- `S1_SUPPORT_MATRIX.md:11`：命中 `支持`「不支持」（子串）；其余未命中。
- `S1_SUPPORT_MATRIX.md:12`：命中 `harness`「generic/C/Go harness」；其余未命中。
- `S1_SUPPORT_MATRIX.md:13` / `:14`：全部未命中。
- `S1_SUPPORT_MATRIX_EN.md:9` / `:10`：全部未命中（`Production-support` 非 `supported` 连续子串）。
- `S1_SUPPORT_MATRIX_EN.md:11`：命中 `supported`「Unsupported」（子串）；其余未命中。
- `S1_SUPPORT_MATRIX_EN.md:12`：命中 `harness`「generic/C/Go harnesses」；其余未命中。
- `S1_SUPPORT_MATRIX_EN.md:13` / `:14`：全部未命中。

### B7. 数据行总数与语言名去重列表

- 数据行总数：每文件 6 行（md line 9–14；EN line 9–14）。
- 表格数据行中出现的语言名（去重）：`C`、`Go`（出自 md:9「纯 Go」、md:12「generic/C/Go harness」；EN:9「Pure-Go」、EN:12「generic/C/Go harnesses」）。
- `Python`、`Node` 在两文件表格数据行中均未出现。

### B8. 表格外政策段/尾注段

语言名与 `三语言`/`three-language`/`harness`/`冻结`/`frozen` 在政策段、信任边界段、尾注段均未命中；`支持` 命中行逐字如下：

```text
S1_SUPPORT_MATRIX.md:21 - 待审批的 ADR-010 提议：新 minor 发布时在 GitHub Release notes 与本支持矩阵同时公告前一 minor 的 EOL 日期，日期为发布后第 90 个自然日。
S1_SUPPORT_MATRIX.md:29 > 切片同步（P55 `EXIT-SEED-BUILD-EVIDENCE`）：本片不适用（无平台与版本支持范围增减）；理由：补证只归档既有 seed 驱动构建的两平台捕获，平台清单、支持层级与 EOL 政策不变。
S1_SUPPORT_MATRIX.md:31 > 切片同步（P57 `P5-P6-CLOSEOUT-ASSESSMENT`）：本片不适用（无平台清单、支持层级、EOL 政策或安装模型的增减）；理由：P57 只把 T018 剩余门整理为收口清单，未改变任何门的判定，本文件措辞逐字不动。
S1_SUPPORT_MATRIX.md:33 > 切片同步（P71 `S1-SCOPE-ANNOTATE`）：TUI 移出与 visible 降级保证口径已同步；理由：R-07 b1 附条件 2 落地只改 REQ 阶段归属与调试指引——平台清单与支持层级不变，TUI 移出系 P68 已裁定的**功能阶段归属同步**，Windows 行状态仍为「正式支持候选」，visible 行仍为「降级保证」。
```

段行号区间与标题行：

```text
S1_SUPPORT_MATRIX.md:16-23 / _EN:16-23   ## 版本与安全维护候选政策 / ## Candidate Version and Security-Maintenance Policy
S1_SUPPORT_MATRIX.md:25-27 / _EN:25-27   ## 发行信任边界 / ## Release Trust Boundary
S1_SUPPORT_MATRIX.md:29,31,33 / _EN:29,31,33   切片同步尾注三条
```

输出字符数：7866
BATCH-02 DONE

---

## 主控抽验记录

- AT-01 CN/EN `:31`、AT-11 `:41` 原文一致（2026-11-14 主控机械复核）。
- S1_SUPPORT_MATRIX 表头 `:7` 与 harness 行 `:12` 原文一致；数据行 6 行（表内 `|` 起始行共 7 行 = 表头 + 6）。

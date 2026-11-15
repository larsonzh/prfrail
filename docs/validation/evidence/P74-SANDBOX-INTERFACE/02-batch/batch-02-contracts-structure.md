# BATCH-02 CONTRACTS §7 结构与台账 P74 段事实枚举

- 载体：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单批量辅助）
- 任务限定：只列事实、不判断、不排序、不做 Top-N 择优
- 实派次数：第 2 次（§二 强制 ≥2）

## A. `docs/CONTRACTS.md` 章节结构（总行数 403）

```
docs/CONTRACTS.md:7:## 1. 契约成熟度与兼容
docs/CONTRACTS.md:27:## 2. Schema 与静态检查责任
docs/CONTRACTS.md:56:### 2.1 Chain 状态、事件与恢复（T009 权威）
docs/CONTRACTS.md:96:## 3. 配置解析与冻结
docs/CONTRACTS.md:124:## 4. 快照、证据与接受事务
docs/CONTRACTS.md:168:### 4.1 证据封存屏障与不可变上下文摘要（P43 权威）
docs/CONTRACTS.md:178:### 4.2 就绪分级与探针许可一次性语义（P44 权威）
docs/CONTRACTS.md:190:## 5. Gate Hook 协议
docs/CONTRACTS.md:202:## 6. 票据、修复与交接
docs/CONTRACTS.md:210:## 7. AgentRunner 与辅助消息通道边界
docs/CONTRACTS.md:310:### 7.1 supervised-black-box 显式模式
docs/CONTRACTS.md:321:### 7.2 Sandbox 检查层接口
docs/CONTRACTS.md:360:## 8. 错误与验收边界
docs/CONTRACTS.md:382:## 9. 产品与生命周期契约（RFC §19）
```

## B. `docs/CONTRACTS_EN.md` 章节结构（总行数 395）

```
docs/CONTRACTS_EN.md:7:## 1. Maturity and Compatibility
docs/CONTRACTS_EN.md:22:## 2. Schema and Checker Responsibilities
docs/CONTRACTS_EN.md:54:### 2.1 Chain State, Events, and Recovery (T009 Authority)
docs/CONTRACTS_EN.md:84:## 3. Configuration Resolution
docs/CONTRACTS_EN.md:109:## 4. Snapshots, Evidence and Acceptance
docs/CONTRACTS_EN.md:159:### 4.1 Evidence Seal Barrier and Immutable Context Digest (P43 authority)
docs/CONTRACTS_EN.md:169:### 4.2 Readiness Levels and One-Shot Probe Permit Semantics (P44 authority)
docs/CONTRACTS_EN.md:181:## 5. Gate Hooks
docs/CONTRACTS_EN.md:194:## 6. Tickets, Repair and Handoff
docs/CONTRACTS_EN.md:202:## 7. AgentRunner and Auxiliary Messaging Boundaries
docs/CONTRACTS_EN.md:298:### 7.1 supervised-black-box explicit mode
docs/CONTRACTS_EN.md:309:### 7.2 Sandbox check-layer interface
docs/CONTRACTS_EN.md:348:## 8. Errors and Acceptance
docs/CONTRACTS_EN.md:373:## 9. Product and Lifecycle Contracts (RFC Section 19)
```

## C. `docs/CONTRACTS.md` `## 7.` 段结构

- 段起点：`docs/CONTRACTS.md:210: ## 7. AgentRunner 与辅助消息通道边界`
- 段终点：行 210–359（行 360 为下一标题 `## 8. 错误与验收边界`）
- 段内 `### ` 子段（共 2 个；行 211–309 为段正文，无其他 `### ` 标题）：
  - 行 310–320：`### 7.1 supervised-black-box 显式模式`
  - 行 321–359：`### 7.2 Sandbox 检查层接口`

## D. `### 7.2` 子段原文范围

`docs/CONTRACTS.md` 行 321–359，八条 (a)–(h)：检查层定位 / 裁决三态 / 错误分类 / 检查点位置（方案 i）/ 触发面（S1）/ 平台实现 / 边界（不做）/ B 段未解锁声明。英文侧 `docs/CONTRACTS_EN.md` 行 309–347 对应同名子段。

## E. `docs/t027/REMAINING_SLICES.md` P74 段（行 2042–2061）

段标题：`### P74 — \`SANDBOX-INTERFACE\`（层 1 主线 B 段引入 · 实现类）· 🔄 进行中`；字段表含标识 / 状态 / 目标 / 依赖 / 规模 / 步骤 / 现状缺口 / 交付物 / 验收证据 / 边界 / 分级 / 探针额度 / 执行自由度清单落盘路径 / R/D 基线 / 临时件落点声明。

P73 状态标记行原文：`docs/t027/REMAINING_SLICES.md:2026: | 状态 | ✅ 已完成（2026-11-14 收口） |`

## F. `docs/t027/REMAINING_SLICES_EN.md` P74 段（行 2029–2048）

段标题：`### P74 — \`SANDBOX-INTERFACE\` (layer-1 part-B introduction · implementation) · 🔄 In progress`；字段表与中文侧逐项对应（Identifier / Status / Goal / Dependency / Size / Steps / Current gap / Deliverables / Acceptance evidence / Boundary / Tier / Probe budget / Freedom-list landing path / R/D baseline / Temp-file location declaration）。

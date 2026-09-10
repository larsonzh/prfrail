﻿# T025 操作员交互验证报告

[English](t025-operator-interaction_EN.md)

日期：2026-09-10。结论：`T025 COMPLETE / AT-22 PASS`。

## 实现范围

1. `schemas/operator-interaction.schema.json` 与 `testdata/contracts/{valid,invalid}/operator-interaction.json`：冻结 request/response wire、角色、绑定字段和额外字段拒绝；纳入 30 份 Schema、12 个 fixture 文件所承载 86 个用例的独立契约门禁。
2. `internal/chain/operator_interaction.go`：RFC 8785 规范摘要、严格解码、追加式 ledger、pending 重放、重复/篡改/错误绑定拒绝，以及同 conversationId、新 requestId 的恢复命令。
3. Engine 单写者控制 API：`OpenOperatorInteraction` 持久化 chain/task/step 暂停，`ResumeOperatorInteraction` 重新验证 request/response、run/task/step/attempt/workspace/conversation/context、当前状态和最新事件证据后恢复 step/task。chain 保持暂停，直到显式运行入口继续。
4. `internal/console/interactions.go`：`interactions list/respond/tui`；支持 JSON 自动化输出和默认无 ANSI 的聚焦终端收件箱。

## AT-22 覆盖

1. 只有通过严格结构、角色、摘要和证据校验的 request 才能开启交互；自由文本不能触发状态转换。
2. 响应必须属于指定 operator、attempt、workspace、conversation、context 和 request hash，且 selection 必须来自允许集合；陈旧、重复、篡改和异源响应 fail-close。
3. ledger 采用追加式 JSONL；重启后从 request/response 重建 pending，不依赖聊天历史、TUI 内存或 SessionBridge `@sbr-review`。
4. TUI 展示问题、原因、风险、所需动作与允许答复；无效输入或持久化失败不移除 pending。
5. Engine 以 request/response record hash 作为状态事件证据。open 先暂停 chain，再写 task/step 等待；resume 先恢复 step、再恢复 task，完成前不返回可执行恢复命令。
6. EventStore 无跨事件事务时，部分写入保持 chain 暂停；同一证据哈希的重试可收敛，不同响应不能借半完成状态接管。
7. 普通澄清只回传结构化 selection；需要授权、评审或人工写入时仍分别进入 authorization、review 或 handoff，不能借 operator interaction 绕过。

## 门禁结果

1. `go build ./...`：通过。
2. `go vet ./...`：通过。
3. `go test -count=1 ./...`：通过，13 个含测试包及 2 个无测试 command 包。
4. `node tools/contracts/contracts.test.js`：通过，2/2；30 份 Schema、12 个 fixture 文件所承载 86 个用例和两条 canonical 向量均通过。
5. `git diff --check`：通过。
6. 仓库无独立 encoding checker；本次 `.go/.json` 保持 UTF-8 无 BOM + LF，`.md` 保持 UTF-8 BOM + LF，并用只读字节检查复核。

## 聚焦测试

- chain 9 项：ledger 单响应、错误绑定、状态转换、严格解码/篡改、当前绑定复验、Engine 持久化/重放、写失败、部分写入收敛、异源 run/陈旧状态拒绝。
- console 7 项：重启重建 inbox、响应持久化、陈旧/非结构化输入、缺 ledger 无副作用、TUI 正常选择、TUI 无效输入、响应写失败保持 pending。
- contract 1 组：operator-interaction 正反 fixture 与全量 Schema/canonical validator。

## 边界

1. `interactions tui` 是聚焦、无 ANSI 的终端交互，不是完整统一 TUI；Bubble Tea 原型未进入发行二进制。
2. T025 返回经过验证的恢复命令并持久化 chain 状态，但没有启动或恢复真实 CLI Agent；AgentRunner wire、能力探针和 session 执行由 T026/T027 完成。
3. timeout/断线没有默认答复；没有追加有效 response 时，重启后的 request 仍为 pending。
4. 本轮未执行真实模型/网络调用、`git commit`、`git push`、签名或发布。

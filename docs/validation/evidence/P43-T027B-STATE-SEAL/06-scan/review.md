# P43-T027B-STATE-SEAL · ⑥ 独立扫描（第二轮原样）

载体：`independent-reviewer`；模型：`MAI-Code-1.1-Flash (copilot)`；只读。首轮判定 `INDEPENDENT SCAN: BLOCKED`（C5 的 K 未接线），修复后本轮重扫。

## Section A 扫描范围与输入
范围：只读核查 internal/chain/attempt.go、seal.go、engine.go、state_seal_test.go，聚焦 ContextKey/TaskContextKey、SealRecord.validate、sameClosure、acceptFromExistingSeal、runStep，以及相关回归测试。

## Section B 发现清单
1 · High · internal/chain/engine.go:306-318 · 运行时已在 seal 生成前写入 ContextKey，并用 TaskContextKey 绑定 accepted closure；K 已真正接入最终 acceptance 路径 · 保持该路径作为唯一 seal 入口
2 · High · internal/chain/engine.go:349-353 · acceptFromExistingSeal 重新派生 taskKeyNow 并与 record.ContextKey 逐项比较；跨 attempt / stale parent / wrong workspace 会被标记 REPAIR · 当前判定已覆盖主 stale-accept 问题
3 · Medium · internal/chain/seal.go:74-82,199-228 · validate 要求 ContextKey 非空时必须是合法 hash，sameClosure 比较 ContextKey；stale seal 的 K 差异会被拒绝 · 补直接 mismatch 用例
4 · Low · internal/chain/attempt.go:70-93 · TaskContextKey 与 ContextKey 共用同一根函数，task 版是 step/request 维度的投影；未来新增字段可能漂移 · 约束 task 版恒为 ContextKey(..., stepID="", requestID="", requestHash="", executionTarget="") 的投影
5 · Low · internal/chain/supervised_blackbox.go、internal/chain/router.go · supervised / handoff / cancel 路径当时无 ContextKey/TaskContextKey · 保持「最终 accepted seal 只来自 engine 的审核/推广路径」

## Section C 与 ①/⑤ 的一致性核对
- ① 运行时接线已闭合：seal 写入 K（engine.go）、复用前重派生比较（engine.go:349-353）、`runStep` 为每个 request 计算 K（engine.go:423）、seal 载入校验（seal.go:74-82）、`sameClosure` 比较 K（seal.go:211-212）。
- ⑤ 派生一致性：两个派生函数共享同一根函数，参数序列固定；当前无现实漂移，残余风险是未来维护忘记同步新增字段。
- 结论：C4 与 N3 对齐，K 已进入运行时判定链。

## Section D 未覆盖面与残余风险
- 已覆盖：TestContextKeyStableAcrossReentry、TestSealRecordBindsAcceptedClosure、TestSealBindsTaskContextKey、TestPassedTaskSealMustVerifyOnRecovery、TestTerminalRouteRefusesForeignAttempt。
- 缺失：stale 父快照/工作区差异的直接用例；K 仅因 attempt/parentHash 变化的表格用例；supervised / handoff / cancel 路径的 K 绑定（非最终 acceptance 路径）；seal 目录级耐久性无 crash 注入证明。

## Section E 结论要旨
本轮变更已把 K 从「函数 API」修正为「真正进入运行时判定链」：seal 写入、seal 再用、step request 派生、跨 attempt stale 拒绝均落地。残余为非最终路径未绑 K、stale mismatch 直接矩阵测试不完整，属维护提醒级，不构成阻断。

INDEPENDENT SCAN: PASS WITH FINDINGS

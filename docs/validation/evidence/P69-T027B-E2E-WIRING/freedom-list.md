# P69 T027B-E2E-WIRING 执行自由度清单（freedom-list）

切片：`T027B-E2E-WIRING`（P69 · 层 2 治理侧 · 实现类 [SLICE] · M 档）
步骤 0 落盘：2026-11-12；**⑥ 之后、⑦ 之前定稿**（§1.7.3 时序说明，OB-48）。本清单由主控落盘，⑦ 复核其审查输入包中的同一清单。

| # | 行使的执行自由度 | 依据（§1.7.1） | 备注 |
|---|---|---|---|
| F-1 | permit 注入形态取「装配选项持 `ProbePermit *gates.ProbePermit` 指针 + 判定点现场构造 `ReadinessInput`」，而非 P44 报告 :156 处方字面（注入 `ReadinessInput`） | 契约已明确范围内的实现方式 | design-r1 Q2 结论（① 定稿）；D3 触发条件「注入形态定稿」就此满足；与 P44 处方的措辞偏差在本片报告如实登记 |
| F-2 | CLI 判定目标级别取 `GoalProbe`（装配起跑 = 一次 probe 起跑；L2/L3 在装配入口结构上不可达） | 单文件实现选择 | design-r1 Q1 结论 |
| F-3 | supervised 可重入判别取「`SupervisedInteractionsLedgerPath` 不存在 = 首次起跑」（fresh 门），`ExecuteSupervisedReturn` 不受门控 | 不影响契约语义的重构 | design-r1 Q4 结论；ledger 写入点在 Require 之后，首启已消费额度 |
| F-4 | CLI 只做 presence 判定不消费额度；AgentRunner 判定两次（Step 0 + helper 内）为纯函数冗余 | 契约已明确范围内的实现方式 | design-r1 Q4；Require 在 Step 1 后、Step 2 marker 前 |
| F-5 | `--probe-permit` 用 `LoadProbePermit`（生产不得每跑重签额度）；缺旗标且路由到装配 ⇒ fail-closed | 单文件实现选择 | design-r1 Q3 |
| F-6 | noop 路径与 `ai check` 路径本片不接线（② 步骤清单未含；`ai check` 的 permit 生命周期需独立设计，属后续切片） | 不影响契约语义的重构 | design-r1 Q3 结论 |
| F-7 | 测试侧 permit 形态统一 `NewProbePermit(id, t.TempDir()/marker, 1)`；`LoadProbePermit` 仅 CLI 测试经旗标天然使用 | 已授权包内的实现方式 | design-r1 Q5 |
| F-8 | ③ 反例由 design-r1 的 5 条拆为 7 条（反例 1/2 各补 supervised 同型） | 单文件实现选择 | ⑤ r1 以 Info 登记（计数口径），报告如实登记 |
| F-9 | ⑤ r1 唯一 Low（`_ = permit.Complete()` 吞错）的最小修复取注释级论证失败不可达（不改行为、不加日志机制） | 不影响契约语义的重构 | ⑤ 无 Medium+ ⇒ 不触发 ⑤ 重跑；修复后 ④ 全量重跑全绿 |
| F-10 | 变异检验表达：M3 保编译补 `_ = *probePermit`、M5 以 `_ = releaseProbePermit` 表达删 defer（删行本身是编译错误） | 已授权包内的实现方式 | ③ 变异证据口径说明第 3/4 条 |
| F-11 | 正例测试平台感知双分支（durability proven 断言成功 / unproven 断言耐久拒绝），两端都验证额度恰消费一次与二跑 Consumed | 单文件实现选择 | ③ 口径说明第 1 条 |

R/D 基线：**R = 0、D = 0**（§13.1.1 步骤 0 登记；至本定稿时点仍为 0）。

# P85 · 主控对 ① 方案书的评审注记（2026-11-22）

方案书 = `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design.md`（① `deep-reasoner` + DeepSeek V4 Pro Max，462 行）。

## 评审结论：采纳（无退回项）

- 起跑消息六的逐项覆盖核：file:line 改动点 ✅；契约影响判定（触/不触 + 理由）✅；
  回归测试点（反例 + 变异）✅；每项可证伪闭合断言 ✅；跨包共享契约 ✅；失败路径枚举 ✅。
- 契约引用行独立核对（主控实读 CONTRACTS.md）：:144-148（review-receipt / policyHash
  引用真实存在）✅；:214（fail-closed 预检 + 失败原因闭集）✅；:216（windows-credential 段）✅；
  :220（正式通道启动前探测）✅。判定依据与磁盘实况一致。
- 三项判定 = T2 不触 / T3 触 / T8a 触，理由充分：T2 制品 store-local 同构
  handoff-policy.json，CONTRACTS:147 的「引用是否真实存在」即其语义落点；
  T3 / T8a 均触及 :214 的预检判据与失败原因闭集 ⇒ 协议先行。
- 设计边界核：无新 schema / 新 CLI 命令 / 新依赖（Go 1.22 stdlib）✅；
  T8a 的 S1 恒拒发 BYOK 与起跑消息「不能证明 ⇒ 零启动」一致 ✅；
  T3 的 Linux fd-exec + Windows ACL 方向与 P83 ⑦ r3 High 整改方向一致 ✅。
- 主控补充要求（已并入 ② 派单）：
  a. ②a 契约修订的 EN 侧落点须按语义同位段落 grep 定位（CONTRACTS 双语对历史非同位镜像，
     只保证增删计数对称 + 语义同构）；
  b. T2 制品目录 `<runRoot>/agent-runner-policy/` 须确认与既有 runRoot 目录族命名一致，
     并经「受管 agent 写入根不相交」断言（② 落地时实测）；
  c. 契约修订段落其余句逐字节不动（协议先行的最小 diff）。

## 下一步

②a（契约先行，prfrail-implementer + DeepSeek V4.1 Flash）→ ②b 代码 → ③ 测试。

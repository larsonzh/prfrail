=========================================
切片：P44-T027B-GATE-READINESS 门控就绪分离（M 档 · 重构片 · [SLICE]）
=========================================

【〇·零 起跑前自检】
核对仓库实况（HEAD=8a46b01、准则 v1.60、SELFTEST 697/697），若 ≥2 项不符，立即停止并报终端读数。

【〇 本片起跑前的裁定】
1. 切片标识：`P44-T027B-GATE-READINESS`
2. 目标：解决 C7 缺陷与 N4 漏洞。具体实现：
   - **C7 门控分离**：实现 `EvaluateT027BReadiness`，严格分离“探针就绪”、“真实验收就绪”与“生产发布就绪”三个级别，严禁混用。
   - **N4 一次性许可**：实现 `ProbePermit` 的一次性消费机制（调用前原子预留，消费后立即失效），防止并发或重复取得额度。
3. 依赖：P43 闭环。
4. 规模：M。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）。
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（可写）。**必须派单，严禁主控代写！**
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（可写，测试先行）。**必须派单，严禁主控代写！**
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，额度 1+1）。
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（只读，额度 3）。
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（只读，1+2）。
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`。
- ⑨ / ⑧c / ⑩：由主控执行。
- 边界：仅修改 `internal/console/**` 与 `internal/gates/**`；不动 `internal/chain/**` 等；不推 gitee；停 ⑩ 等授权。

【二、派单硬性要求】
- **严禁模型擅自升级**：每次 `runSubagent` 必须显式传 `agentName` 与 `model`。
- 顺序约束：③ 偏移必须早于 ④ 与 ⑤；⑥ 先于 ⑦。

【三、批量辅助】
- 本片必须实派 `GLM-5.3-Flash (glm)` ≥ 2 次。场景：枚举现有的 readiness 判定点、分析 `ProbePermit` 的并发竞争日志。

【四、临时脚本与中间产物】
- 一律落盘 `tmp/P44-T027B-GATE-READINESS/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录后，首先执行「起跑消息落盘」：将本段消息（从切片名到回执要求结束）落盘为 `docs/validation/evidence/P44-T027B-GATE-READINESS/00-briefing.md`（BOM+LF）。
2. 在 `REMAINING_SLICES{,_EN}.md` 登记 P44 台账（16 字段块、`R/D=0`、临时件落点）。
3. 落盘 `freedom-list.md`，并生成 `MANIFEST.md` 与 `SHA256SUMS.txt`（须包含 00-briefing.md 与 freedom-list.md）。
4. 落盘 v33 清单副本。

【六、① 架构要求】
1. 派发 V4 Pro 输出《C7+N4 门控就绪与一次性许可方案》。
2. 必须明确回答：三个级别的 readiness 如何确定义与隔离？`ProbePermit` 的原子抢占与过期机制如何保证严格一次性？并发拒绝时返回什么错误码？
3. 单次输出限 10,000 字符内，禁止自重启。

【七、执行节奏（常规不停）】
1. 步骤 0 → ① → ②③ → ④ → ⑤ → ⑥ → ⑦ → ⑧ → ⑨ → ⑧c → ⑩。
2. ④ 门禁全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 则原位重跑；⑩ 是唯一常规停点，无需逐环请示。

【八、文档同步更新与补录要求】
1. `docs/CONTRACTS.md`：补充 readiness 三级定义与 `ProbePermit` 一次性消费语义。
2. `docs/ARCHITECTURE.md`：补充门控分离机制与许可发放模块。
3. `docs/DEV_PLAN.md`：更新 T027 重构里程碑关于 C7/N4 的完成状态。
4. 台账回写：最终事实包中体现回写状态。

【九、⑩ 回执要求】
包含派单表（含真实时序、agentName、model、transcript 偏移）、③ 顺序自查、⑤⑥⑦结论链、三作用域门禁（含 scanned N > 0）、最终事实包（12 字段，含文档回写状态、00-briefing.md 的哈希）。

请先落盘 v33（含起跑消息存档的新纪律），跑 CI 双绿后，立刻进入 P44 步骤 0。完成后无需请示，直接连轴推进至 ⑩ 停点。

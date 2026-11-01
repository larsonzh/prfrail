# P45 起跑消息（原文落盘）

> 落盘纪律：v33 §七 第 10 条「起跑消息即工件」。以下为用户起跑消息原文，未改写（仅剥离附件清单，本消息无附件）。

---

DR13-FIX 停点回执已收。判定：定性精准、修复彻底、反例先红后绿、主控零代写。GLM 辅助火力全开，堪称并发排查的教科书。
DR13-FIX 正式闭环。

【关于 v35 快照】
由于 P45 是本阶段最后一个 M 档大装配片，且 DR13-FIX 属于前置插入片，v35 快照将在 P45 闭环后统一生成。当前直接以 `eb36992` 为锚点起跑 P45。

【P45 正式起跑】
=========================================
切片：P45-T027B-E2E-ZERO-MODEL 零模型三任务 E2E（M 档 · 装配验证片 · [SLICE]）
=========================================

【〇·零 起跑前自检】
核对仓库实况（HEAD=eb36992、准则 v1.60、SELFTEST 697/697），若 ≥2 项不符，立即停止并报终端读数。

【〇 本片起跑前的裁定】
1. 切片标识：`P45-T027B-E2E-ZERO-MODEL`
2. 目标：验证 P39–P44 的所有零件（预算、锁序、归属、绑定、封存、门控）能在 `console` 装配 + `chain` 状态机中贯通成一条零模型三任务链：
   - T1 正常执行并被接受（进入 `PASSED`，父快照绑定）；
   - T2 触发自然失败（`FAILED`），且 **T3 不启动**（链暂停或终止）；
   - 全程零模型调用，所有步骤用声明式 stub 代替真实 AI。
3. 依赖：P44 闭环；DR13-FIX 闭环（并发隐患已排除）。
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
- 边界：仅 `internal/chain/**`、`internal/console/**` 与 `examples/**`；不动 `internal/adapters/**` 等；不推 gitee；停 ⑩ 等授权。

【二、派单硬性要求】
- **严禁主控代写 ②③**：每次 `runSubagent` 必须显式传 `agentName` 与 `model`。
- 顺序约束：③ 偏移必须早于 ④ 与 ⑤；⑥ 先于 ⑦。

【三、批量辅助】
- 本片必须实派 `GLM-5.3-Flash (glm)` ≥ 2 次（建议 4 次）。场景：枚举三任务装配路径、梳理 console 与 chain 的接口契约、比对测试日志。

【四、临时脚本与中间产物】
- 一律落盘 `tmp/P45-T027B-E2E-ZERO-MODEL/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录后，首先执行「起跑消息落盘」：将本段消息（从切片名到回执要求结束）落盘为 `docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/00-briefing.md`（BOM+LF）。
2. 在 `REMAINING_SLICES{,_EN}.md` 登记 P45 台账（16 字段块、`R/D=0`、临时件落点）。
3. 落盘 `freedom-list.md`，并生成 `MANIFEST.md` 与 `SHA256SUMS.txt`（须包含 00-briefing.md 与 freedom-list.md）。
4. 落盘 v34 清单副本。

【六、① 架构要求（含 DR13-FIX 强输入）】
1. 派发 V4 Pro 输出《零模型三任务 E2E 装配方案》。
2. 必须明确回答：
   - `console` 的哪个入口负责启动三任务链？如何注入 stub 让 T1 成功、T2 自然失败？
   - T2 失败后如何保证 T3 不启动（链状态机在哪一步拦截）？
   - 父快照绑定与事件链连续性如何断言？
   - 零模型调用如何机械证明（例如 stub 的调用计数 = 0）？
3. **必须吸收 DR13-FIX 的 F 节输入**：E2E 测试必须包含「同一 launcher/runRoot 并发 Start」场景断言：
   - 恰一个赢家，败者全命中 `Conflict`/`Duplicate` 哨兵；closeout 时 `registry.len()==1`。
   - 必测「赢家进程在装配结束仍存活」（`ProcessAlive`）。
   - 核对 slot owner token == 赢家 launchID、mirror 与 slot 一致。
   - **禁止使用 10s freshness 墙钟作排他断言依据**，改用完成序/状态断言。
4. 单次输出限 10,000 字符内，禁止自重启。

【七、执行节奏（常规不停）】
1. 步骤 0 → ① → ②③ → ④ → ⑤ → ⑥ → ⑦ → ⑧ → ⑨ → ⑧c → ⑩。
2. ④ 门禁全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 则原位重跑；⑩ 是唯一常规停点，无需逐环请示。

【八、文档同步更新与补录要求】
1. `docs/ARCHITECTURE.md`：补充三任务 E2E 装配路径、失败隔离机制，以及并发 Start 场景断言。
2. `docs/DEV_PLAN.md`：更新 T027 重构里程碑关于 E2E 的完成状态。
3. 台账回写：最终事实包中体现回写状态。

【九、⑩ 回执要求】
包含派单表、③ 顺序自查、⑤⑥⑦结论链、三作用域门禁（含 scanned N > 0）、最终事实包（12 字段，含文档回写状态、00-briefing.md 的哈希）。

现在开始执行 P45 步骤 0 与 ① 架构。完成后直接连轴推进至 ⑩ 停点。

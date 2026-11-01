# DR13-FIX 起跑消息（原文落盘）

> 落盘纪律：v33 §七 第 10 条「起跑消息即工件」。以下为用户起跑消息原文，未改写（仅剥离附件清单，本消息无附件）。

---

v34 已落盘（`8e834b2`）并通过 CI 双绿，队列 3/4 已补全，审核通过。

【关于 DR-13 的裁定：选 A】
1. DR-13 二次复现，升格为阻断级。先起 DR13-FIX 小片定性并修复，再起跑 P45。
2. **纠正一处预算判断**：你看到的「19,786 / 20,000」是 10 月终值，11 月 1 日 08:00 Copilot Max 额度已重置为满额。请按 11 月实际可用额度规划，不要用 10 月终值做约束。
3. DR13-FIX 的定性结论将直接输入 P45 的 ① 设计（P45 涉及并发装配），所以这不是额外的开销，而是 P45 的前置输入。
4. **定位声明（写进 DR13-FIX 报告 §1 与台账）**：本片是 P45 的前置修复片，**不占 v34 队列 1–4 的编号**；v34 保持为「P44 闭环 → P45 起跑前」的基线，本片闭环后出 v35 统一反映。**不需要先修订 v34**。

【DR13-FIX 起跑消息】
=========================================
切片：DR13-FIX 并发启动 flaky 定性修复（S 档 · 修复片 · [SLICE]）
=========================================

【〇·零 起跑前自检】
核对仓库实况（HEAD=8e834b2、准则 v1.60、SELFTEST 697/697），若 ≥2 项不符，立即停止并报终端读数。

【〇 本片起跑前的裁定】
1. 切片标识：`DR13-FIX`
2. **定位**：P45 的前置修复片，不占 v34 队列 1–4 的编号；本片闭环后出 v35 统一反映。
3. 目标：定性并修复 `TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner` 的 flaky 根因。两种可能：
   - **A：测试装置过严**——8 个 goroutine 抢同一 launch slot，断言「全程恰 1 个赢家」，但 stub 进程可能在竞争者全部返回前退出并触发 reclaim，导致「先后两个赢家」。若定性为此，修断言（改为「同一时刻至多 1 个存活」，而非「全程恰 1 个成功」）。
   - **B：真实排他窗口漏洞**——launch slot 的 O_EXCL 或 reclaim 路径存在真实竞态。若定性为此，修实现。
4. 依赖：P44 闭环；DR-13 二次复现（run `41253798647` + `41257990185`）。
5. 规模：S（单测试文件 + 可能的实现微调）。
6. **本片定性结论将输入 P45 的 ① 设计**。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，定性根因）。
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（可写）。**必须派单，严禁主控代写！**
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（可写，若需补强反例）。**必须派单，严禁主控代写！**
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，额度 1+1）。
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（只读，额度 3）。
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（只读，1+1，S 档）。
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`。
- ⑨ / ⑧c / ⑩：由主控执行。
- 边界：仅 `internal/adapters/agent_runner_pinned_cli*.go` 与其测试文件；不动其他包；不推 gitee；停 ⑩ 等授权。

【二、派单硬性要求】
- 严禁主控代写 ②③；每次 `runSubagent` 必须显式传 `agentName` 与 `model`。
- 顺序约束：③ 偏移必须早于 ④ 与 ⑤；⑥ 先于 ⑦。

【三、批量辅助】
- 本片必须实派 `GLM-5.3-Flash (glm)` ≥ 2 次。场景：分析并发失败日志、枚举 launch slot 的 O_EXCL 与 reclaim 路径。

【四、临时脚本与中间产物】
- 一律落盘 `tmp/DR13-FIX/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录后，首先执行「起跑消息落盘」：将本段消息落盘为 `docs/validation/evidence/DR13-FIX/00-briefing.md`（BOM+LF）。
2. 在 `REMAINING_SLICES{,_EN}.md` 登记 DR13-FIX 台账（16 字段块、`R/D=0`、临时件落点、定位声明）。
3. 落盘 `freedom-list.md`，生成 `MANIFEST.md` 与 `SHA256SUMS.txt`。

【六、① 定性要求】
1. 派发 V4 Pro 做只读定性：分析 `internal/adapters/agent_runner_pinned_cli.go` 的 launch slot 独占创建 + reclaim 路径，与 `agent_runner_pinned_cli_test.go:902` 的断言逻辑，判断是装置过严还是真实漏洞。
2. 必须给出：根因归属（A 或 B）、证据链（file:line）、修法建议。
3. 单次输出限 10,000 字符内，禁止自重启。

【七、执行节奏（常规不停）】
1. 步骤 0 → ① 定性 → ②③ 修复 → ④ → ⑤ → ⑥ → ⑦ → ⑧ → ⑨ → ⑧c → ⑩。
2. ④ 门禁全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 则原位重跑；⑩ 是唯一常规停点。

【八、⑩ 回执要求】
1. ① 的定性结论（A 或 B）与证据链。
2. 修复 diff 与断言变更说明。
3. 补强反例（若定性为 B，须有能复现原竞态的测试）。
4. 三作用域门禁（含 scanned N > 0）。
5. 最终事实包（12 字段），含 DR-13 状态更新（观测级 → 已修复）。

现在开始执行 DR13-FIX 步骤 0。完成后直接连轴推进至 ⑩ 停点，无需逐环请示。DR13-FIX 闭环后，再起跑 P45。

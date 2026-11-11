# freedom-list — P64 AT24-SUPERVISED-BLACK-BOX（步骤 0 落盘）

> 本片执行自由度清单（§1.7.3 固定路径；G7 扫描域外 ⇒ 无机械判据保护，由 ⑦ 人工复核）。
> 基线：R = 0、D = 0（2026-11-11 起跑登记）。

| 编号 | 自由度项 | 边界与出处 |
|---|---|---|
| F-1 | 不自行 commit / push | ⑩ 停点等同一轮授权；提交信息英文 `<type>: <summary>`（起跑消息【六】11 与仓库指令） |
| F-2 | 不启动任何新切片；不替用户裁定；不自行升准则版本 | §1.7.2 三条 bullet（起跑消息【七】2） |
| F-3 | 环境：不开 Ubuntu VM、不开网络代理、不启用 GPT-6.1 Sol 双打 | 起跑消息【〇】1 |
| F-4 | R-09.1~R-09.9 逐项判定 + AT-24 验收结论的唯一主体 = ⑦；主控只整理事实、不判定 | 起跑消息 3.2、【七】4 |
| F-5 | 不改产品代码（除 ② 修 T028 缺陷且范围限于 `internal/console/**` 与 `internal/chain/**` 的 supervised 相关文件）；不改 CONTRACTS / schema / fixtures / workflows / 准则 / R-09 / T028 报告原文 | 起跑消息 3.3、【四】8–17、【七】8 |
| F-6 | 模型绑定按起跑消息：①⑤ V4 Pro、②③⑧ Flash、⑥ MAI、⑦ Codex、批量辅助 GLM-5.3-Flash；每次派单同时传 agentName 与 model；严格串行 | 起跑消息【一】【二】 |
| F-7 | 额度：⑤ 1+1、⑥ 3、⑦ 1+3（高风险档）+ 盲审另计 1；批量辅助 ≥2 次；用尽 ⇒ 停下升级用户裁决 | 起跑消息【一】【二】 |
| F-8 | 单次输出 ≤ 10,000 字符；批量辅助 max_tokens ≥ 2048、回执报实得字符数 | 起跑消息【〇·零】4、【二】 |
| F-9 | 验收腿 = 机制腿（stub visible 投递 + 模拟人工归还 + 后置扫描）；真实端到端登记为 unknown/边界 | 起跑消息【〇】3、【七】6 |
| F-10 | ② 范围扩展（用户裁决①，2026-11-11 同轮授权）：`schemas/supervised-report.schema.json` + `testdata/contracts/{valid,invalid}/runtime-records.json` 中 supervised 相关 fixture（如确需同步）+ 同步 `tools/contracts/contracts.test.js` 断言。定性 = schema 实现对齐 CONTRACTS §7.1（非契约语义变更，不走 §1.7.2 升级边界） | 用户裁决①原文（逐字）：「授权扩展：本轮授权 ② 修复范围从 `internal/console/**` 与 `internal/chain/**` 扩展为 `+ schemas/supervised-report.schema.json + testdata/contracts/{valid,invalid}/runtime-records.json` 中 supervised 相关 fixture（如确需同步）。」 |
| F-11 | ② 范围扩展（用户裁决②分支 A，2026-11-11 同轮授权）：新增 junction 等价越界反例正式测试 `internal/console/runtime_supervised_junction_probe_test.go`（supervised 相关）；前置探针等价成立（见 `03-test/L09-01-junction-probe.txt`） | 用户裁决②（探针 PASS ⇒ 分支 A：补 `TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot`） |

## 读数校正注记（起跑前已知事实，v54 生成时的读数修正候选）

1. **实测读数**：`docs/t027/CHECKLIST_v52{,_EN}.md` 按 LF（0x0A 字节）计数，实况均为 **277**（CN/EN 对称）；双实现复核一致（PowerShell `ReadAllBytes` 逐字节计数 + Node `Buffer` 逐字节计数）；两文件末尾字节均为 `0x0A`；`git status --short -- docs/t027/CHECKLIST_v52.md docs/t027/CHECKLIST_v52_EN.md` 干净无改动。
2. **与 v53 §一 记录的差异**：v53 §一 记 v52=278，实况 277，差 1。定性 = **v53 §一 生成时的读数偏差**，与 v51 §一 记 v50=233（v52 起修正为 232）的处理口径**完全同构**（详见 v53 §一 的 v50 口径注记）。
3. **对 P64 的影响**：无。按 P64 起跑消息【〇·零】第 2 条自身规则——「除 HEAD 预期前进外，若仍有 ≥2 项不符 ⇒ 停止并报实况」——8 个锚点中仅此 1 项不符（**不是** ≥2 项），**不触发停止**。其余 7 项（HEAD / origin / 准则 v1.64 / SELFTEST 722/722 / 三作用域 TOTAL_FAIL=0 / v50=232 / v51=270）全部符合预期。
4. **处置**：v53 §一 的 278 读数**不追改**（历史快照冻结）；⑩ 回执「给下版清单的事实包」增一条「v52 全文件行数 v54 起修正为 277（CN/EN 对称），v53 §一 记 278 系读数偏差」。
5. **性质**：本注记不构成 P64 的越界登记，也不改变 P64 的任何验收语义。

## R-09.3 / R-09.7 补录说明

- 时间线：步骤 0（2026-11-11 起跑）→ 索取确认时用户不在线（自动应答）→ 用户本人于 ③ 执行前（2026-11-11 本轮）逐字补确认两句。
- 原因：主控按纪律不把自动应答当作用户确认（§1.7.2 授权可回溯性；不得伪造授权）。
- 处置：保留两件 `01-arch/` 原「未确认」记录段（逐字不删不改）+ 追加「补录确认」段 + 同轮重算 MANIFEST / SHA256SUMS 两件哈希（其余条目不动）并回读校验三方一致。

# 00-briefing — P56 `G4B-CORPUS-DECOUPLE` 起跑消息逐字原文

性质：用户 2026-11-05 起跑消息的逐字转录（围栏内为原文，围栏外为本包的性质说明）。
用途：按准则 §七 第 10 条「起跑消息即工件」入包，并由 `MANIFEST.md` / `SHA256SUMS.txt` 与报告 `artifacts` 围栏共同声明。

```text
切片：P56 G4B-CORPUS-DECOUPLE G4b 语料枚举与索引解耦（S/M 档 · 判据优化类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读刚才发送的 `docs/t027/CHECKLIST_v45.md`（跨会话快照）。
2. 核对第一节的 5 个锚点（HEAD 应为 `c76c6ce`，若 v45 提交后 HEAD 前进，则以 v45 提交号为锚点）。
3. **主控模型声明**：本片为门禁工具纯代码重构，无底层并发/状态机变更，默认使用 `DeepSeek V4.1 Flash`。请在回执首行标 `[自执行]`。
4. **记忆唤醒**：
   - 「实现类切片 ⑨ 全绿后必停 ⑩」：未获授权绝不自行 commit/push。
   - 「v45 §九 门禁坑」：尤其是写 `.md` 必须补 BOM；修改 `tools/gates` 后必须重跑 `--selftest` 并关注判定逻辑漂移。
   - 「纪律 V · 大任务拆分」：遇超时先查 `git status` 半截产物。

【〇 本片起跑前置说明与授权】
1. 根据 v45 快照，队列第 8 项 `AT23-RESTRICTED-CANDIDATE` 因依赖 R-08/R-09 裁定而处于 ⛔ 阻断状态。
2. **用户同轮裁定**：暂缓第 8 项，激活并起跑队列中的备选治理项 `G4B-CORPUS-DECOUPLE`。此决定已由主控在步骤 0 的 `00-briefing.md` 中如实登记。
3. **架构双打（Sol）**：本片**不启用**。直接进入执行，无需 ① 架构与红队。

【一、授权边界（流水线角色与模型）】
- ① 架构：跳过。
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+2，M 档）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控执行。

【二、派单硬性要求与纪律 V 预警】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。严格串行。
- **强制执行纪律 V**：修改 `g4b.js` 和 `selftest.js` 时，拆分小步派发。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：检索 `ctx.listAllMarkdown` 的所有调用点并列出影响面、交叉核对新的变异断言），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/G4B-CORPUS-DECOUPLE/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录 `tmp/G4B-CORPUS-DECOUPLE/` 与 `docs/validation/evidence/G4B-CORPUS-DECOUPLE/`。
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 BOM + LF）。
3. 在 `docs/t027/REMAINING_SLICES{,_EN}.md` 登记台账（状态 `🔄 进行中`，`R/D=0`）。
4. 落盘 `freedom-list.md`、`MANIFEST.md` 与 `SHA256SUMS.txt`、`artifacts` 围栏。
5. 前置核对：查阅 `docs/DELIVERY_DIRECTIVE.md` 的 G4b 条款及 `tools/gates/lib/criteria/g4b.js`、`tools/gates/lib/ctx.js`。

【六、执行流水线与核心要求】
本片目标是解决 OB-120 指出的“G4b 语料枚举与索引耦合”导致的判定漂移和假阳性问题：
1. **定位与修改（②）**：
   - 阅读 `tools/gates/lib/ctx.js` 中的 `listAllMarkdown()` 及 `g4b.js`。
   - 将 `listAllMarkdown()` 的 tree/index 作用域分支从依赖索引（`git ls-files`）改为基于 HEAD 的对象库（如 `git ls-tree -r --name-only HEAD`），使语料读取与工作区暂存状态彻底解耦。
   - 确保缓存键或判定逻辑不再受“暂存删除 .md 文件”这一状态影响。
2. **测试（③）**：
   - 在 `tools/gates/selftest.js` 中新增测试。
   - 必须补一条针对“暂存删除 .md 文件”场景的变异断言：模拟暂存删除一个 `.md`，证明冷跑（清缓存）和热跑（走缓存）给出的判定完全一致，且语料不发生错误收缩。
3. **性能验证**：
   - 修改后需对比修改前后 G4b 的耗时，确保 `--selftest` 和三作用域门禁不出现显著退化。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v45 §八 防呆硬规则与 §九 门禁坑。

【八、文档同步更新（v45 第 18 条与 §11.2）】
在 ⑧ 阶段，必须同步更新或显式声明不适用以下文档（不得遗漏）：
1. `docs/validation/g4b-corpus-decouple{,_EN}.md`：新建本片验证报告。
2. `docs/DELIVERY_DIRECTIVE{,_EN}.md`：若 `listAllMarkdown` 的枚举语义变更导致 §6.3 的 G4b 条款描述需调整，则同步修订判据文本；并在附录 C 追加本治理片的 C.0xx 条目。若判据文本无需改变，如实声明。
3. `docs/t027/REMAINING_SLICES{,_EN}.md`：台账回写（状态翻 ✅，写入完成记录）。
4. `docs/DEV_PLAN{,_EN}.md`：治理片按 §11.2 规定不回写，必须在此声明“不适用 + 理由（治理片不回写）”。
5. `docs/ARCHITECTURE{,_EN}.md`：若架构文档中提及了 G4b 的语料读取机制，需同步更新；若未提及，声明“不适用”。

【九、⑩ 回执要求】
1. 派单表（逐字机械读取）。
2. 修改说明：如何解耦、如何确保判定 `hit ≡ cold`。
3. 测试证据：新增的变异断言及运行结果。
4. 性能对比：修改前 vs 修改后的 `--selftest` 耗时及三作用域耗时。
5. 三作用域门禁与 `SELFTEST`。
6. 文档同步更新的 `git diff --stat`。

现在开始执行 P56 步骤 0。
```

## 说明

- 本文件为**逐字转录**（围栏内不改写、不补全），围栏外的本段为本包的性质说明。
- 起跑前自检结论（锚点核对）记录在 `input-list.md`；自由度清单见 `freedom-list.md`。

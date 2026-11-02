# P50 CONSOLE-MULTITASK-ASSEMBLY 起跑消息（原文落盘）

> 落盘纪律：v39 §七 第 10 条「起跑消息即工件」。以下为用户起跑消息原文，未改写。

---

切片：P50 CONSOLE-MULTITASK-ASSEMBLY console 多任务生产装配扩展（L 档 · 实现类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v39.md`（跨会话快照）。
2. 核对第一节的 5 个锚点（HEAD 应为 `8bd8224`）。
3. **主控模型声明**：本片为涉及状态机与核心装配的 L 档实现切片，按 v39「纪律 P」，本会话主控模型已提升为 **`DeepSeek V4 Pro`**。请在回执中明确记录。
4. **记忆唤醒（红线复述）**：请在回执中复述以下规则，证明你已掌握：
   - 「评估/实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获同轮授权绝不自行 commit/push。
   - 「纪律 S（会话长度与压缩纪律）」：上下文达 100K 或出现压缩提示，立即准备结片。
   - 「环境预检与复用铁律」：本片为本地代码装配，若涉及环境检查，先探查已有环境再安装。
   - 「v39 §九 门禁坑」：特别注意契约同步、BOM 补齐、tmp 污染等 13 条经验。
5. 若锚点 ≥2 项不符，立即停止并报仓库实况。

【〇 本片起跑前的裁定】
1. **目标**：解除 console 层的单任务限制。目前 `runtime_agent.go:774` 与 `runtime_supervised.go:335` 强制单任务运行。本片须将其重构为支持多任务并发装配，且不破坏既有的 A6 排他契约。
2. **前置依赖（关键！）**：P45 的 D9 曾明确记载“多任务装配与 A6『每 runRoot 唯一 identity mirror』冲突”。因此，**本片必须协议先行**。如果架构推演发现必须修改 `CONTRACTS` 或底层锁机制才能支持多任务，必须在 ① 阶段给出明确的协议修订方案。
3. **风险定级**：L 档核心架构重构。严禁把任务降级为 `[ROUTINE]`。必须走完整流水线，严禁主控代写。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1）
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（必须派单，严禁代写）
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（必须派单，TDD 先行）
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+3，L 档/高风险配置）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：由主控（V4 Pro）执行。

【二、派单硬性要求】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。
- **严格串行**：绝对禁止并行调用多个子代理（R1.1 / R5.1）。
- 主控不得代写任何 ②③ 代码或测试。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：旧单任务限制的引用面清查、相关测试用例枚举），严格串行，遵守输出上限纪律 Q。

【四、临时脚本与中间产物】
一律落盘 `tmp/CONSOLE-MULTITASK-ASSEMBLY/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录 `tmp/CONSOLE-MULTITASK-ASSEMBLY/` 与 `docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/`。
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF）。
3. 在 `REMAINING_SLICES{,_EN}.md` 登记台账（状态设为 `🔄 进行中`，`R/D=0`，临时件落点）。
4. 落盘 `freedom-list.md`、`MANIFEST.md` 与 `SHA256SUMS.txt`。
5. `artifacts` 围栏声明 `00-briefing.md`。
6. 前置核对：查阅 `docs/validation/p45-e2e-zero-model.md` 中的 D9 结论，以及 `internal/console/runtime_agent.go` 和 `runtime_supervised.go` 的单任务限制逻辑。

【六、执行流水线与核心要求】
1. **① 架构（关键）**：
   - 必须回答：如何在解除单任务限制的同时，不破坏 A6（每 runRoot 唯一 identity mirror）的并发契约？
   - 若需修改 `CONTRACTS.md`，必须在 ① 中给出精确的“协议先行”修改方案，并作为单独的第一步实施。
   - 必须评估现有 `ExecuteNoopRun` 与 `localRuntime` 在并发下的安全边界。
2. **③ 测试先行（TDD 强制）**：在修改生产代码前，必须先写一个红测，证明当前 console 在尝试装配并发多任务时会被拒绝或发生冲突。
3. **② 实现**：按架构方案和契约修改实现多任务支持。任何失败必须原位修复。
4. **④ 门禁**：本地 `go build/vet/test ./...` 全绿。若涉及 race 敏感，注意利用本地/远程环境。
5. **⑤⑥⑦ 审查**：重点审查“并发多任务”下的上下文覆盖、锁竞争以及 A6 镜像是否被破坏。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v39 §八 纪律 T/U（防字面量误杀）以及 §九 的 13 条门禁坑。

【八、文档同步更新（v39 第 18 条铁律 + 历史欠账清理）】
在 ⑧ 阶段，必须同步更新以下权威文档：
1. `docs/ARCHITECTURE.md`（及 `_EN`）：更新 console 多任务装配、A6 镜像共享、调度机制的相关章节。
2. **`docs/DEV_PLAN.md`（及 `_EN`）：这是必须做的历史欠账清理。** 请在 T027-B（防御性重构期）的段落中，补齐 P39–P49 的核心里程碑事实（如：C1-C7 修复、持久预算、锁序、归属、证据、封存、门控就绪、零模型 E2E、Race 腿扩包等），并标记本片 P50 为“console 多任务装配完成”。
3. 台账回写：完成记录中体现文档更新状态（尤其要报告 DEV_PLAN 的补齐情况）。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. 红绿反转证据（测试路径与日志）。
3. 三作用域门禁（含 `scanned N > 0`）、`SELFTEST 697/697`。
4. **文档同步更新**的 `git diff --stat`，必须明确包含 `DEV_PLAN.md` 的补齐行数。
5. 最终事实包，含未闭环项与 v39 快照的差异点。

现在开始执行 P50 步骤 0。

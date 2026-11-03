# 起跑消息原文（P51 CONSOLE-CLI-MULTITASK-WIRING）

以下为用户起跑消息逐字原文，未作任何增删。

---

切片：P51 CONSOLE-CLI-MULTITASK-WIRING console 多任务 CLI 接线（M 档 · 实现类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v40.md`（跨会话快照）。
2. 核对第一节的 5 个锚点。
   **⚠️ 锚点代差免责声明**：v40 快照中写的 HEAD 是 `81f9dc9`（P50 主提交）。但实际上，生成 v40 后紧接着产生了一个 docs-only 回填提交 `3aa53d5`。因此，**当前实际的 HEAD 是 `3aa53d5`（或更新）**。请直接以 `git rev-parse --short HEAD` 的实况为准，不要花时间去 diff 这两个 commit，这是正常的文档回填代差。
3. **主控模型声明**：本会话主控模型为默认的 **`DeepSeek V4.1 Flash`**（依据 v40 纪律 P 第 8 条）。若执行中遇到 Flash 连续 2 次行为漂移或报错，主控应主动停下并提示切换 Pro。
4. **记忆唤醒（红线复述）**：请在回执中复述以下规则：
   - 「评估/实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获授权绝不自行 commit/push。
   - **「纪律 V · 大任务拆分与半截防污染」**：如果派发实施类任务遇到网络超时，先查 `git status` 是否有半截产物；若有则精确回退，严禁盲目重试；大任务必须拆分成小步串行派发。
   - 「v40 §九 门禁坑」：BOM 补齐、tmp 污染、契约同步等 13 条。
5. 若锚点 ≥2 项不符（除上述代差外），立即停止并报仓库实况。

【〇 本片起跑前的裁定】
1. **目标**：补齐 P50 遗留的 CLI 接线工作。让用户可以通过命令行参数（如 `prfrail run` 或 `prfrail preview`）直接驱动 P50 已经实现的多任务串行装配。
2. **前置依赖**：P50 已实现 `runtime_agent.go` 和 `runtime_supervised.go` 的多任务支持，本片只做 CLI 层（通常是 `cmd/prfrail` 和 `internal/console/cli.go`）的参数解析、路由和装配注入。
3. **风险定级**：M 档。不触及底层引擎状态机，主要防范参数解析错误和装配漏项。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1）
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（必须派单，严禁代写）
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（必须派单，TDD 先行）
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+2）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：由主控（Flash）执行。

【二、派单硬性要求与纪律 V 预警】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。**严格串行，禁止并行。**
- **强制执行纪律 V**：本片虽然是 M 档，但为了防范网络超时和通道假死，**② 实施阶段的派单必须拆分为小步**（如：第一步加 CLI 参数定义，第二步在 runtime 里挂接，第三步改测试），**严禁一次性把整个 CLI 改造全部塞给一个子代理**。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：现有 CLI flag 定义检索、CLI 端到端测试用例枚举），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/CONSOLE-CLI-MULTITASK-WIRING/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录 `tmp/CONSOLE-CLI-MULTITASK-WIRING/` 与 `docs/validation/evidence/CONSOLE-CLI-MULTITASK-WIRING/`。
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF）。
3. 在 `REMAINING_SLICES{,_EN}.md` 登记台账（状态设为 `🔄 进行中`，`R/D=0`）。
4. 落盘 `freedom-list.md`、`MANIFEST.md` 与 `SHA256SUMS.txt`。
5. `artifacts` 围栏声明 `00-briefing.md`。
6. 前置核对：查阅 P50 报告及 `internal/console/cli.go` / `cmd/prfrail/main.go` 中当前的参数接收情况，以及 `runtime_agent.go` 多任务入口的函数签名。

【六、执行流水线与核心要求】
1. **① 架构**：排查从 CLI 入口到 `AgentRunnerAssembly` 的调用链。设计 CLI 参数接收形态（如是否支持传入多 task 的 JSON/YAML 定义），以及如何将解析后的对象传递给 P50 实现的串行装配入口。明确需要修改的文件清单。
2. **③ 测试先行（TDD 强制）**：写红测，从 CLI 层面调用（或模拟调用），断言当前 CLI 尚未把多任务参数正确传递给装配层（或直接被拒绝）。红测必须稳、准、狠。
3. **② 实现（按纪律 V 拆分）**：
   - 必须先加 CLI 参数定义与解析逻辑。
   - 然后接驳多任务装配。
   - 如果涉及 `cmd/prfrail`，确保不破坏现有的单任务兼容性（如果有）。
4. **④ 门禁**：`go build ./...`、`go vet ./...` 全绿。
5. **⑤⑥⑦ 审查**：重点审查 CLI 路由是否有遗漏、错误处理是否完备、是否意外改变了原有的单任务行为。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格防范 v40 §九 门禁坑。

【八、文档同步更新（v40 第 18 条）】
在 ⑧ 阶段，必须同步更新以下权威文档：
1. `docs/ARCHITECTURE.md`（及 `_EN`）：在 CLI/TUI 交互或 Console 章节，补充多任务 CLI 的用法说明。
2. `docs/DEV_PLAN.md`（及 `_EN`）：更新 T027-B 里程碑中关于 CLI 装配的完成状态。
3. 台账回写：完成记录中体现文档更新状态。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. 红绿反转证据（路径与日志）。
3. 三作用域门禁与 `SELFTEST 697/697`。
4. 文档同步更新的 `git diff --stat`。
5. 最终事实包，含未闭环项与 v40 快照的差异点。

现在开始执行 P51 步骤 0。

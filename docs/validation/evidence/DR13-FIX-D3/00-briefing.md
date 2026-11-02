# P48 DR13-FIX-D3 起跑消息（原文落盘）

> 落盘纪律：v37 §七 第 9 条「起跑消息即工件」。以下为用户起跑消息原文，未改写（附件 `CHECKLIST_v37.md` 为仓库文件，不复制入包）。

---

切片：P48 DR13-FIX-D3 跨进程并发加固（S 档 · 实现类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v37.md`（跨会话快照）。
2. 核对第一节的 5 个锚点：
   - **HEAD 应为 `e742224`**（注意：v37 正文里写的 `50a70ac` 是自指代差，以 `e742224` 为准）。
   - 准则 v1.60，SELFTEST 697/697，三作用域全绿。
3. **主控模型声明**：本片为底层并发/句柄类实现切片，按 v37「纪律 P」，本会话主控模型已提升为 **`DeepSeek V4 Pro`**。请在回执中明确记录此事实（作为防漂移的基础依据）。
4. **记忆唤醒（红线复述）**：请在回执中复述以下四条规则，证明你已掌握：
   - 「评估/实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获同轮授权绝不自行 `commit`/`push`。
   - 「纪律 S（会话长度与压缩纪律）」：上下文达 100K 或出现压缩提示，立即准备结片。
   - 「纪律 T（禁写量词形态）」与「纪律 U（禁逐字引用触发串）」：防止文档被 `G6-4` 误杀。
5. 若锚点 ≥2 项不符，立即停止并报仓库实况。

【〇 本片起跑前的裁定】
1. **标识**：`DR13-FIX-D3`（v37 §四 层 1 第 1 项）。
2. **目标**：修复 DR13-FIX 遗留的 D3 缺口。`release` 逻辑目前是先 `unlock` 锁，再 `os.Remove` 文件，这两步之间存在毫秒级的并发窗口（W1）。本片须将其改为**先 `os.Remove` 再 `unlock`**，配合 D2 完全闭合该并发窗口，并完成 Windows remove-while-open 全量回归。
3. **风险定级**：底层并发/句柄操作。**绝对不可因为“只改两行”就当做 [ROUTINE] 处理**。必须走完整流水线，严禁主控代写。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（必须派单，严禁代写）
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（必须派单，严禁代写）
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+2）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：由主控（V4 Pro）执行。

【二、派单硬性要求】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。
- **严格串行**：`runSubagent` 串行阻塞，绝对禁止并行调用多个子代理（R1.1 / R5.1）。
- 主控不得代写任何 ②③ 代码或测试。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：用例检索、回归影响面枚举），严格串行，单次输出 ≤ 10000 字符。

【四、临时脚本与中间产物】
一律落盘 `tmp/DR13-FIX-D3/`；只增不删；由 ⑧c 统一清理；清理后须确认 `tmp` 顶层仅剩 `.gitkeep` 和 `gates-cache`。

【五、步骤 0 待办（按序执行）】
1. 建目录 `tmp/DR13-FIX-D3/` 与 `docs/validation/evidence/DR13-FIX-D3/`。
2. **起跑消息落档**：将本段消息逐字落盘为 `docs/validation/evidence/DR13-FIX-D3/00-briefing.md`（UTF-8 with BOM + LF，不改写不摘要）。
3. 在 `REMAINING_SLICES{,_EN}.md` 登记 DR13-FIX-D3 台账（状态设为 `🔄 进行中`，`R/D=0` 基线，临时件落点 `tmp/DR13-FIX-D3/`）。
4. 落盘 `freedom-list.md`，生成 `MANIFEST.md` 与 `SHA256SUMS.txt`（必须包含 00-briefing 和 freedom-list）。
5. 在最终报告的 `artifacts` 围栏声明 00-briefing.md。
6. 核对前置阅读：`docs/validation/dr13-fix.md` 的「D3」与 `internal/adapters/agent_runner_pinned_cli.go` 的 `release()` 逻辑。

【六、执行流水线与核心要求】
1. **① 架构**：输出《D3 顺序反转与回归方案》。必须回答：`os.Remove` 失败时，锁是否仍能释放？Windows 的 `remove-while-open` 语义在不同 OS 版本上的差异是什么？
2. **③ 测试先行（TDD 强制）**：**必须先在未修复的代码上，写出一个确定性失败的红测**（能稳定卡住 unlock 与 remove 之间的空隙并触发竞态）。验证红测跑红后，才能进入 ②。
3. **② 实现**：实现“先 remove 再 unlock”，在 Windows 上跑通全量回归。任何失败必须原位修复。
4. **④ 门禁**：`go build` / `go vet` / `go test ./...` 必须全绿。
5. **⑤⑥⑦ 审查**：审查要重点关注“删除失败”与“锁未释放”的组合状态是否安全（fail-closed），不得遗留竞态隐患。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**（D-02 教训）。若 CI 已触发，只作过程事实记录。
3. 严格遵守「纪律 T（禁写数字+量词形态）」与「纪律 U（禁逐字引用触发串）」，防止被 G6-4 误杀。
4. 若上下文达到 100K 或出现压缩提示，立即结束当前片（纪律 S）。

【八、文档同步更新（v37 第 15 条铁律）】
本片改变了底层锁与文件释放的执行顺序，涉及 `internal/adapters/agent_runner_pinned_cli.go`。在 ⑧ 阶段，**必须同步更新**：
1. `docs/ARCHITECTURE.md`（及 `_EN`）：在 §4.2 并发锁序原则中，补充“文件锁释放与物理删除的原子顺序”说明。
2. `docs/DEV_PLAN.md`（及 `_EN`）：若有相关里程碑描述，需同步更新。
3. 台账回写：完成记录中体现文档更新状态。漏写此项将按 v37 规则视为结构缺陷。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. 红绿反转证据（红测首跑红的日志路径、转绿后的日志路径）。
3. 三作用域门禁（含 `scanned N > 0`）、`SELFTEST 697/697` 读数。
4. **文档同步更新**的 `git diff --stat` 记录。
5. 最终事实包（12 字段），含未闭环项与 v37 快照的差异点。

现在开始执行 P48 步骤 0。

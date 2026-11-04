切片：P54 CONTRACTS-S7-DURABILITY Windows 耐久契约升级与 C6 落地（L 档 · 契约与实现类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v43.md`（跨会话快照）。
2. 核对第一节的 5 个锚点。**⚠️ 锚点代差免责声明**：v43 中记录的 HEAD 是 `7802819`，实际 HEAD 因回填提交 `2b8ac3e` 已前进。请以实际 `git rev-parse --short HEAD` 为准。同时请顺手将 v43 §十 第一句里的“v42”改为“v43”。
3. **主控模型声明**：本片为底层契约与持久化 API 修改，按 v43「纪律 P」提升为 **`DeepSeek V4 Pro`**。请在回执中记录。
4. **记忆唤醒（红线复述）**：
   - 「实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获授权绝不自行 commit/push。
   - 「纪律 V · 大任务拆分」：大任务拆分小步串行，遇超时先查 `git status` 半截产物。
   - 「v43 §九 门禁坑 15 条」。
5. 若锚点 ≥2 项不符，立即停止并报仓库实况。

【〇 本片起跑前的裁定与双打授权】
1. **R-02 裁定**：正式批准 `CONTRACTS-S7-DURABILITY` 立项。
2. **目标**：将 Windows 平台的发布耐久等级从 `unproven` 升级为 `proven`，需落实 C6，并全面修订 `CONTRACTS.md` §7。
3. **① 架构双打授权（§1.7.4 临时授权）**：
   本片涉及不可测的“断电/耐久语义”与契约修订，极易出现“自认为防住了”的假象。主控在完成 ① 架构后，必须临时调用 `independent-reviewer` + `GPT-6.1 Sol (copilot)` 进行红队抗压审查（限本片、一次性）。红队需重点审查：契约边界是否夸大、C6 的同步机制在异常路径下是否退化为普通写、非 Windows 平台的抽象是否纯净。
4. **环境说明**：本片主要改动 Windows 平台 API。**无需开启 Ubuntu VM**。Unix 侧的“零影响”只需在本地使用 `GOOS=linux go build/vet` 验证，并由 CI 的 ubuntu 腿兜底即可。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（主推） + `independent-reviewer` + `GPT-6.1 Sol (copilot)`（红队）
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+3，L 档高配）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（V4 Pro）执行。

【二、派单硬性要求与纪律 V 预警】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。严格串行。
- **强制执行纪律 V**：② 实施必须拆分小步（如：先改 CONTRACTS 文本，再实现 C6 封装，再改 replay_store 调用，再改测试）。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：契约修订点枚举、C6 相关的 Windows syscall 用法索引），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/CONTRACTS-S7-DURABILITY/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录 `tmp/CONTRACTS-S7-DURABILITY/` 与 `docs/validation/evidence/CONTRACTS-S7-DURABILITY/`。
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF）。
3. 在 `REMAINING_SLICES{,_EN}.md` 登记台账（状态 `🔄 进行中`，`R/D=0`）。
4. 落盘 `freedom-list.md`、`MANIFEST.md` 与 `SHA256SUMS.txt`。
5. `artifacts` 围栏声明。
6. 前置核对：查阅 `docs/CONTRACTS.md` §7 的耐久定义，以及 A7/B3b 相关验证报告。

【六、执行流水线与核心要求】
1. **① 架构**：设计 C6 的落地方式（例如：在 `agent_runner_replay_store.go` 的 `publish` 原语中引入 Windows 平台特有的 `MoveFileEx` 调用）。必须确保 Unix 平台零影响。同时明确 CONTRACTS 的精确修订文本。
2. **③ 测试先行（TDD 强制）**：写红测证明当前平台在 Windows 下发布失败/处于 `unproven`，无法进行首派发。由于是 Windows 专属 API，测试需按平台构建标签隔离。
3. **② 实现**：先改 CONTRACTS（协议先行），再写 Windows 实现代码，确保接口签名不变。
4. **④ 门禁**：全量测试绿，注意跨平台构建（GOOS=windows 与 linux）。
5. **⑤⑥⑦ 审查**：Codex 会极其严苛地核对 CONTRACTS 修订与实际代码行为的一致性。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v43 §八 防呆硬规则与 v43 §九 门禁坑。

【八、文档同步更新（v43 第 18 条）】
在 ⑧ 阶段，必须同步更新以下权威文档（不得遗漏）：
1. `docs/CONTRACTS{,_EN}.md`：本次的核心交付，§7 耐久条件与定义修订。
2. `docs/ARCHITECTURE{,_EN}.md`：在存储/数据流相关章节补充 Windows 耐久落盘机制。
3. `docs/SECURITY{,_EN}.md`：更新与持久化数据完整性相关的安全边界说明。
4. `docs/DEV_PLAN{,_EN}.md`：将 T027-B 的发布耐久项更新为已完成。
5. `docs/ADR_REGISTER{,_EN}.md`：新增一条 ADR（记录本片使用 C6 的决策依据与平台前提）。
6. `docs/t027/REMAINING_SLICES{,_EN}.md`：台账回写。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. ① 双打红队审查结论。
3. 红绿反转证据（含 Windows 环境下的跑红、跑绿记录）。
4. 跨平台编译证据（GOOS=linux 与 windows）。
5. 三作用域门禁与 `SELFTEST 697/697`。
6. 文档同步更新的 `git diff --stat`（必须覆盖上述六项文档）。
7. 最终事实包。

现在开始执行 P54 步骤 0。

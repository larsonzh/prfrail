切片：P52 ACCEPTANCE-SECURITY-COUNTEREXAMPLES 接受链安全反例（M 档 · 安全验证类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v41.md`（跨会话快照）。
2. 核对第一节的 5 个锚点。
   **⚠️ 锚点代差免责声明**：v41 中记录的 HEAD 是 `2b0af2d`（P51 主提交）。实际 HEAD 因回填提交 `941f60e` 已前进。请以实际 `git rev-parse --short HEAD` 为准，不用花时间排查。
3. **主控模型声明**：本片为安全对抗/验证类切片，按 v41「纪律 P」提升为 **`DeepSeek V4 Pro`**。请在回执中明确记录。
4. **记忆唤醒（红线复述）**：
   - 「实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获授权绝不自行 commit/push。
   - 「纪律 V · 大任务拆分」：大任务必须拆分小步串行，遇超时先查 `git status` 半截产物。
   - 「v41 防呆硬规则」：测字符数用内联脚本，严禁在 `$env:TEMP` 等磁盘生成临时文件。
   - 「v41 §九 门禁坑 13 条」。
5. 若锚点 ≥2 项不符，立即停止并报仓库实况。

【〇 本片起跑前的裁定】
1. **目标**：针对断点 A 的接受链缺口，落地 GAP-51/52/55/60/64/70 所代表的安全反例验证。证明当前系统在面对"非法输入、路径伪造、句柄越权、边界绕过"时，会正确 `fail-closed`。
2. **前置依赖**：P51 已闭环，多任务装配与 CLI 链路已打通。
3. **边界**：只做测试用例与断言，**绝不修改生产代码（除非发现真实安全漏洞需上报，但本片内仍不修）**。不改接受语义。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，评估威胁模型）
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`（只写测试文件）
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`（运行并验证）
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+2）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：由主控（V4 Pro）执行。

【二、派单硬性要求与纪律 V 预警】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。**严格串行，禁止并行。**
- **强制执行纪律 V**：② 实施阶段的派单必须拆分（例如：先写 GAP-51/52 的反例，验证后；再写 GAP-55/60 的反例）。**严禁一次性把几十个反例全塞给一个子代理**，防止网络超时与半截污染。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：GAP-51 等具体条目的代码引用面清查、已有的拒绝用例枚举），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录 `tmp/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/` 与 `docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/`。
2. 起跑消息逐字落盘为 `00-briefing.md`（UTF-8 with BOM + LF）。
3. 在 `REMAINING_SLICES{,_EN}.md` 登记台账（状态 `🔄 进行中`，`R/D=0`）。
4. 落盘 `freedom-list.md`、`MANIFEST.md` 与 `SHA256SUMS.txt`。
5. `artifacts` 围栏声明 `00-briefing.md`。
6. 前置核对：查阅 P46 轮 2 的缺口台账 `GAP-51 / 52 / 55 / 60 / 64 / 70` 的具体定义，以及相关代码（如 `internal/adapters/bounddir_*.go`、`internal/evidence/boundread.go` 等）。

【六、执行流水线与核心要求】
1. **① 架构**：对 GAP-51/52/55/60/64/70 提取"攻击面"。例如：攻击者如何伪造路径替换？如何让句柄读取超界？设计出一套 `TestXxx` 的命名清单与断言矩阵（必须断言 `errors.Is` 特定的拒绝哨兵）。
2. **③ 测试先行（TDD 强制）**：**本片本质上就是写反例，所以红测必须表现为"攻击被成功阻断"（测试 PASS，但验证的是 fail-closed 逻辑）**。若发现某个攻击手段未阻断（即反例测试 FAIL），立刻停止，这属于真实漏洞，必须上报，绝不允许在本片内顺手修生产代码。
3. **② 实现**：针对架构师的清单，在 `*_test.go` 中编写安全反例。
4. **④ 门禁**：`go test ./...` 全绿。
5. **⑤⑥⑦ 审查**：Codex 极其擅长安全审查。它会检查：反例是否真的覆盖了攻击面？测试有没有造假（比如用错包名导致没跑到）？句柄断言是否够死？

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v41 §八 防呆硬规则与 v41 §九 门禁坑。

【八、文档同步更新（v41 第 18 条）】
在 ⑧ 阶段，必须同步更新以下权威文档：
1. `docs/ARCHITECTURE.md`（及 `_EN`）：在安全边界或验证章节，补充接受链安全反例矩阵的说明。
2. `docs/DEV_PLAN.md`（及 `_EN`）：在 T027-B 里程碑中，标记接受链反例已覆盖。
3. 台账回写。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. 反例矩阵与结果证据（每条 GAP 对应哪个测试用例，如何证明其 fail-closed）。
3. 三作用域门禁与 `SELFTEST 697/697`。
4. 文档同步更新的 `git diff --stat`。
5. 最终事实包，含未闭环项与 v41 快照的差异点。

现在开始执行 P52 步骤 0。

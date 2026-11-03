切片：P53 FIX-GAP-77-SYMLINK-ESCAPE 修复中间符号链接逃逸（M 档 · 底层安全修复类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v42.md`（跨会话快照）。
2. 核对第一节的 5 个锚点。
   **⚠️ 锚点代差免责声明**：v42 中记录的 HEAD 是 `c636dd8`（P52 主提交）。实际 HEAD 因回填提交 `dc785ed` 已前进，请以实际 `git rev-parse --short HEAD` 为准。
3. **主控模型声明**：本片为底层并发/句柄类安全修复切片，按 v42「纪律 P」提升为 **`DeepSeek V4 Pro`**。请在回执中记录。
4. **记忆唤醒（红线复述）**：
   - 「实现类切片 ⑨/⑧c 全绿后必停 ⑩」：未获授权绝不自行 commit/push。
   - 「纪律 V · 大任务拆分」：大任务拆分小步串行，遇超时先查 `git status` 半截产物。
   - 「v42 防呆硬规则」：测字符数用内联脚本，严禁在 `$env:TEMP` 等磁盘生成临时文件。
   - 「v42 §九 门禁坑 15 条」。
5. 若锚点 ≥2 项不符，立即停止并报仓库实况。

【〇 本片起跑前的环境确认】
1. **VM 复用**：远程 Ubuntu VM `10.0.0.199` 已开启，请使用 `ssh -i "c:\Users\妙妙呜\.ssh\id_rsa" larson@10.0.0.199` 连通。按 v42 §七 第 19 条，请使用 `bash -lc` 加载环境。
2. **网络代理关闭**：当前环境直连，无需配置代理。

【〇 本片起跑前的裁定】
1. **目标**：修复 **GAP-77**（Linux 平台 `bounddir_unix.go:64` 的 `openat(O_NOFOLLOW)` 仅护末组件，导致中间符号链接逃逸）。
2. **核心交付物**：
   - `internal/adapters/bounddir_unix.go` 的生产代码修复。
   - 移除 `internal/adapters/bounddir_escape_unix_test.go` 中 `TestBoundDirMiddleSymlinkEscapesBound` 的 `t.Skip`，验证红转绿。
   - **规范同步**：更新 `docs/SECURITY.md`（SEC-02 威胁矩阵补充该攻击载荷与验证方法）与 `docs/TEST_STRATEGY.md`（AT-02 或相关基线列为强制反例）。

【一、授权边界（流水线角色与模型）】
- ① 架构：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 实现：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（V4 Pro）执行。

【二、派单硬性要求与纪律 V 预警】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。严格串行。
- **强制执行纪律 V**：本片涉及底层句柄与路径打开改造。派单必须严格拆分：先设计逐级解析方案，再修改一个函数，再验证，严禁打包大改。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：现有 openat 用法枚举、Linux 兼容性评估），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/FIX-GAP-77-SYMLINK-ESCAPE/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办】
1. 建目录。
2. 落盘 `00-briefing.md`。
3. 登记台账。
4. 落盘 `freedom-list.md` 等。
5. `artifacts` 围栏声明。

【六、执行流水线与核心要求】
1. **① 架构**：对比两种修复方案：(A) 逐级循环 `openat` + `O_NOFOLLOW`；(B) 使用 `O_PATH` 打开中间目录，再通过 `/proc/self/fd` 相对打开。必须给出 Linux 兼容性、性能与代码复杂度评估。
2. **③ 红测先行（TDD 强制）**：必须在修改生产代码之前，**移除 T51-1 的 `t.Skip`**，在 Ubuntu VM 上跑出真实的红（读出界外字节），保存日志。
3. **② 实现**：按架构方案修改 `bounddir_unix.go`。**禁止改变接口签名**。修复后，必须保障界内多层目录读取依旧正常（不可误杀），且界外中间软链接被 `fail-closed` 阻断。跑通 T51-1 转绿。
4. **④ 门禁**：本地 `go test ./...` 与 VM Linux 环境 `go test -race` 全绿。
5. **⑤⑥⑦ 审查**：重点盯防 `openat` 的越界与循环漏洞。Codex 将深度审查句柄与路径安全性。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v42 §八 防呆硬规则与 v42 §九 门禁坑。

【八、文档同步更新（v42 第 18 条）】
在 ⑧ 阶段，必须同步更新：
1. `docs/ARCHITECTURE.md`（及 `_EN`）：更新 §5.1 安全反例矩阵，标记 GAP-77 已修复。
2. `docs/SECURITY.md`：在 SEC-02 威胁矩阵补入“中间符号链接逃逸”的真实攻击载荷与防范方法。
3. `docs/TEST_STRATEGY.md`：在 AT-02 或相关测试基线中，将其列为强制反例。
4. `docs/DEV_PLAN.md`（及 `_EN`）：更新 P53 任务状态。
5. 台账回写。

【九、⑩ 回执要求】
1. 派单表（`agentName` 与 `model` 逐字机械读取）。
2. 红绿反转证据（包含 VM 上的真实跑红日志、修复后跑绿日志）。
3. 三作用域门禁与 `SELFTEST 697/697`。
4. 文档同步更新的 `git diff --stat`。
5. 最终事实包，含 GAP-77 状态。

现在开始执行 P53 步骤 0。

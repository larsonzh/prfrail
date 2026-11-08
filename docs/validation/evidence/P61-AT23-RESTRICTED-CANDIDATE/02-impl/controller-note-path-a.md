# ② 收敛裁定注记（主控，2026-11-08）

- ② 回执 `IMPLEMENT: FINDINGS`：任务一（生产封装）与任务三（Luna 连通性）完成；任务二（3 脚本 6 处既有调用点收敛）未执行——实现方判定为契约级冲突，给出 A/B 两条待裁定路径。
- **主控裁定：采纳路径 A（既有调用点维持原样）**。依据：
  1. R-08.1 字面义务 =「**新增** copilot 调用点必须走封装」；既有调用点先于封装存在，不属「新增」。（保留）
  2. **（用户 2026-11-08 纠正）**起跑消息步骤 0 第 7 项的「散落 5 处以上」增强条款**按调用点计数**（PowerShell 侧 6 处 + Go 侧 1 处 = **7 处 > 5**）是**触发**的；但 R-08.1 字面只约束「新增」调用点，既有 7 处不在 R-08.1 的收敛义务内；起跑消息的增强条款与 R-08.1 字面冲突时，**以 R-08.1 为准**（§1.7.2 授权依据须可回溯到 R-08 裁定原文）。故路径 A 如实登记为「与起跑消息增强条款的偏离」，进 ⑤ 预审判断技术理由是否充分。（修正：原「阈值未触发」表述作废）
  3. 实现方技术理由成立：`Invoke-AgentProbe.ps1` / `Invoke-CandidateProbe.ps1` 的 `Start-Process` + 进程树测量是 T026 探针的**测量对象**，收敛会丢测量语义；`container-candidate.ps1` 为 AppContainer 跨容器执行（任务书预授权例外）；CandidateProbe 有意容忍版本文本漂移（记录态），封装的文本钉版会把记录态转拒绝态。（保留）
- 登记口径（如实）：既有 7 处直接调用点维持原样；**本片所有新增 copilot 调用一律走 `Invoke-CopilotCli`**。⑤ 预审必核面四项：① 收敛后无遗漏**新增**直接调用点；② Go 侧 `copilotCLIProbeArgs` 与 R-08.1 `buildCopilotArgs` 的命名/签名偏离；③ **路径 A 的三个技术理由是否独立成立**（`Start-Process` 进程树测量语义 / AppContainer 跨容器预授权 / CandidateProbe 版本漂移记录态）——⑤ 需逐项核对「收敛到封装后行为是否真会变」，不得只采信实现方说法；④（用户 2026-11-08 追加）层 2 的 `five-facts-distinct` → `six-facts-distinct` 断言 id 变更（成因 `02d5797`，A5 pid reports 绑定）——⑤ 需核对断言逻辑**真覆盖原 5 枚 + 新增第 6 枚**，而非仅改名。
- 路径 B（扩封装契约 `-LauncherMode/-ExtraArgumentList`）本片不采纳；若未来需要，另行立项。
- 计费：第 2 笔（② 内 Copilot-hosted GPT-6 Luna 1 次，1 个 Copilot premium）已消耗，结果 `connected`（exit 0，回复首行 Hello!）。

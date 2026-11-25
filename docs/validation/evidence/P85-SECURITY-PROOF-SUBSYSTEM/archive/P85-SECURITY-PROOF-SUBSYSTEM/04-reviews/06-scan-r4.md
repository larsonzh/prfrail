# ⑥ 独立扫描 r4 原文（P85 · MAI Code 2 Flash · 用户 2026-11-23 裁定追加 · OB-140）

INDEPENDENT SCAN: PASS

**A 类**
- F02 — CLOSED：落盘于 agent_runner_policy.go#L635，绑定句柄后复核目录身份，再打开末组件。
- F05 — CLOSED：ExecutableFileID 已落盘于 process.go#L73，digest 后比对于 executable_binding.go#L89，执行器传入于 ai_probe_copilot_executor.go#L670。空值保留兼容分支，生产载体拒绝空身份。
- F20 — CLOSED：落盘于 ai_install_proof.go#L225，读写载体均复核派生目录句柄；兼容路径仍按失败关闭处理。

**B 类移交文案**
- F01 — 如实：正文与 SECURITY.md#L11 一致；仅引用行号应为 :11，不是文案的 :10。
- F04 — 如实：承认内核最终名称解析残余，未夸大闭合。
- F08 — 如实：明确 named ACL、已打开写句柄属于威胁模型外残余，不作已闭合陈述。
- F10 — 如实：准确说明 Windows 不可目录 sync，Linux 已 sync；未声称 Windows 提供额外 crash-durable 证明。

**C 类计划**
- 完整：逐条覆盖 F06/F07、F16、F17、F18、F22，以及旧 VM 基线同步；08-bc-handoff.md#L26 的 VM 命令、豁免登记、T8a rerun 与双平台 CI 包级证据均可执行。

无新根因发现。

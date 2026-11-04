# 00-briefing — 起跑消息原文与主控指令（P55 `EXIT-SEED-BUILD-EVIDENCE`）

归档时点：2026-11-05；锚点 HEAD = `baa144d`（= origin/main）；准则 v1.60。

**性质**：本文件为用户本轮消息的**逐字落盘**（不改写、不摘要、不删节）。其中「一、纪律登记」与「二、关键覆盖警告」为主控指令（含一条纪律登记要求），「三、P55 正式起跑」为切片起跑消息本体。

**覆盖声明（如实登记）**：本文件的前一版为步骤 0 抢跑期由主控自撰的草稿（其转录的是更早一轮「核锚点并看队列」的消息，非正式起跑原文）。按用户本轮裁定，以本次转录的原文**整份覆盖**，并同步重算 `SHA256SUMS.txt`；草稿内容不得被引用为本片的事实来源。

## 用户消息原文（逐字）

```text
收到实况核对。你主动停手并如实报告越界，处置得当。不回退，我们将这次抢跑作为“步骤 0 的草稿”就地转正，但必须留下纪律痕记。

一、纪律登记（硬性）
在即将创建的 P55 验证报告中，登记一条异常：D = 1（无触发词抢跑，违反 §1.6 / §5.0，主控已主动发现并停手，处置方式 = 就地转正，未造成污染）。

二、关键覆盖警告
你抢跑时生成的 `00-briefing.md` 大概率是你自己脑补的切片描述，**并非用户的起跑原文**。现在，请立即用我接下来提供的这则正式起跑消息全文，**覆盖**该文件，并重新计算 `SHA256SUMS.txt`。同理，检查 `input-list.md` 与 `freedom-list.md` 是否符合下文要求，不符合就重写。

三、P55 正式起跑（以此覆盖你的草稿）

切片：P55 EXIT-SEED-BUILD-EVIDENCE AT-13 seed 驱动构建补证（M 档 · 证据补全类 [SLICE]）
=====================================================================

【〇·零 起跑前自检与上下文锚定】
1. 请先完整阅读 `docs/t027/CHECKLIST_v44.md`（跨会话快照）。
2. 核对第一节的 5 个锚点（HEAD 应为 `baa144d`）。
3. **主控模型声明**：本片为补证类工程任务，无底层并发/状态机变更，默认使用 `DeepSeek V4.1 Flash`。**本片你的回执首行统一标注 `[自执行]`**。仅在连续 2 次行为漂移或遇复杂日志根因时切 Pro。
4. **记忆唤醒**：
   - 「实现/验证类切片 ⑨ 全绿后必停 ⑩」：未获授权绝不自行 commit/push。
   - 「纪律 V · 大任务拆分」：大任务拆分小步串行，遇超时先查 `git status` 半截产物。
   - 「v44 §九 门禁坑 15 条」。

【〇 本片起跑前的环境声明与双打授权】
1. **环境预检**：本片需补全 AT-13 证据。
   - 若需 Linux 侧构建/验证，请复用远程 VM `10.0.0.199`（go1.22.12 @ `~/go1.22` 已就位）。访问前须 `bash -lc` 加载环境。
   - 若需拉取外部依赖，请确保代理环境已配置。
2. **架构双打（Sol）**：本片**不启用**。直接进入执行，无需 ① 架构与红队。
3. **目标**：补全 AT-13（seed 驱动构建与外部 oracle 比对）的证据链，**不改 CI 结构，不改产品代码**。

【一、授权边界（流水线角色与模型）】
- ① 架构：跳过。
- ② 实现/补证：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`
- ③ 测试/验证：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（额度 1+2，M 档）
- ⑧ 文档：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（Flash）执行，回执标 `[自执行]`。

【二、派单硬性要求与纪律 V 预警】
- 每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`。严格串行。
- **强制执行纪律 V**：涉及多个证据点的收集，必须拆分小步派发。若某步涉及远程 VM 执行，须单独编写脚本并传输，严禁长 SSH 嵌套命令。

【三、批量辅助】
实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` ≥ 2 次（场景：证据文件 SHA 清单核对、双语段落一致性检查），严格串行。

【四、临时脚本与中间产物】
一律落盘 `tmp/EXIT-SEED-BUILD-EVIDENCE/`；只增不删；⑧c 统一清理。

【五、步骤 0 待办（你已完成部分，请据此修正）】
1. 建目录 `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/`（已有）。
2. **重写** `00-briefing.md`，将本条起跑消息全文逐字落盘，更新 MANIFEST 与 SHA256SUMS。
3. `REMAINING_SLICES{,_EN}.md` 台账（已有）确认无误。
4. 检查并修正 `freedom-list.md`、`input-list.md`。
5. 前置核对：查阅 `docs/validation/s1-selfhost.md` 的历史数据与 `docs/DELIVERY_DIRECTIVE.md` 的测试策略。

【六、执行流水线与核心要求】
本片本质是**补证**，主控应亲自执行或指派 Flash 执行命令并收集证据：
1. **定位缺口**：查阅 T017 的 `s1-selfhost.md`，对比 AT-13 的要求（seed 分离、外部 oracle、clean-room 成功及失败不损 seed），列出缺失的证据条目。
2. **执行与捕获**：
   - 在本机（Windows）与 VM（Ubuntu）上分别模拟或触发 seed 驱动的构建流程。
   - 确保捕获完整的原始 stdout/stderr、退出码、环境变量版本（Go version 等）、以及比对结果。
3. **生成报告**：编写补充证据说明，落盘至 `docs/validation/exit-seed-build-evidence{,_EN}.md`，将核心日志脱敏后存入 `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/`。
4. **禁止事项**：绝对不可修改 `.github/workflows/*.yml`，绝对不可改动 Go 核心代码。本片产出仅限 `.md` 与证据包。

【七、执行节奏与红线】
1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环。
2. **⑩ 是唯一常规停点**。⑨/⑧c 全绿后，**必须停在 ⑩ 等授权，严禁自行 commit/push**。
3. 严格遵守 v44 §八 防呆硬规则与 §九 门禁坑。

【八、文档同步更新（v44 第 18 条）】
在 ⑧ 阶段，必须同步更新以下权威文档（不得遗漏）：
1. `docs/TEST_STRATEGY{,_EN}.md`：更新 AT-13 的状态与证据链接。
2. `docs/OPERATIONS{,_EN}.md`：补充“自托管、升级和发布”章节中的实测证据。
3. `docs/DEV_PLAN{,_EN}.md`：更新 T017/T018 中对应的构建/自托管进度。
4. `docs/S1_SUPPORT_MATRIX{,_EN}.md`：若涉及跨平台构建的最终结论，同步更新矩阵状态。
5. `docs/S1_RELEASE_NOTES{,_EN}.md`：若证据链补全达成发行要求，同步更新状态。
6. `docs/t027/REMAINING_SLICES{,_EN}.md`：台账回写。

【九、⑩ 回执要求】
1. 派单表（逐字机械读取）。
2. 红绿/对比证据（原始输出的摘要或文件路径）。
3. 环境与版本（Windows 与 Linux VM 的验证结果）。
4. 跨平台证据（如有）。
5. 三作用域门禁与 `SELFTEST 697/697`。
6. 文档同步更新的 `git diff --stat`（必须覆盖上述六项文档）。
7. 最终事实包。

现在开始执行 P55 步骤 0 的正式确认与修正。
```

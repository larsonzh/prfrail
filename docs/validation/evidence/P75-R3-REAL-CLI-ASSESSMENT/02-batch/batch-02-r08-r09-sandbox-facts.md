# batch-02 — R-08 / R-09 与 Sandbox 边界事实枚举（GLM-5.3-Flash 原文归档）

切片：P75-R3-REAL-CLI-ASSESSMENT。派单：`deep-reasoner` + `GLM-5.3-Flash (glm)`。实得字符数：5808。原文如下（逐字）。

---

【④ R-08 / R-09 与 R-3 相关条款】源：`docs/validation/r08-r09-ruling.md`

- `r08-r09-ruling.md:1` — 标题原文：「# R-08 / R-09 裁定记录（AT-23 / AT-24 前置裁定）」
- `:3` — 头部原文：「性质：`[ROUTINE]` docs-only 落地（不占切片号、不升准则版本）；裁定人 = 用户；裁定日期 = 2026-11-06；本文档 = 2026-11-07 落地定稿。**CN 为权威版**，`_EN` 为同位镜像」
- `:7` — 「R-08 / R-09 是 v49 §六 的待裁定项，本文件是 2026-11-06 用户裁定的完整记录」
- `:12` — R-08.1 标题：「候选版本 = (a) GitHub Copilot CLI 1.0.83」；`:14`–`:20` 要点：封装 `Invoke-CopilotCli`/`buildCopilotArgs`、`-Unlock` 默认 `$false`（fail-closed）；阶段 1 首次校验+每 5 次探针抽样 1 次、阶段 2 每次校验；版本证据用 `copilot --version`；「新增 copilot 调用点必须走封装」；锁定失败降级「转 (b) 案『版本不可控』」
- `:22` — R-08.2 标题「探针预算」；`:24`–`:26`：「阶段 1：逐次声明（用户同轮批准），首次不预设上限」「阶段 2：Copilot 5 次 / BYOK 100 元 / Copilot 货币 0 / 重试 5 次 / 追加须用户同轮批准」「『逐次声明』= 用户同轮批准，主控不得自行决定」
- `:28` — R-08.3 标题：「Windows 前提 P = (a) 只用 P（进程级崩溃、无断电、本机 Win11/NTFS）」；`:30` — 「AT-23 报告必须如实标注『不覆盖断电场景；断电耐久引用 B3b（P′ 覆盖物理机目标部署环境未论证）』」
- `:32` — R-08.4 标题：「VM / 代理 = (a) 都不允许」；`:34` — 「BYOK 连通性预检失败路径：停下升级用户裁决——可选 (a) 临时启用代理（须同轮批准）/ (b) 改用本机可达端点 / (c) 暂停 AT-23；主控不得自行启用代理或切端点」
- `:36` — R-08.5 标题：「真实候选腿失败处置 = (b) 条件闭合（reduced）」；`:38`–`:40`：「机制腿全绿 + 真实候选腿实际跑过（≥3 次，覆盖启动/执行/停机）并记录 unsupported/unknown」「每次 unsupported/unknown 必须有具体错误信息或证据（stderr 原文、退出码、失败阶段）」「报告必须明写『AT-23 条件闭合，不构成真实候选兼容性证明』」
- `:42`–`:49` — 继承漂移豁免五要素（须同时满足）：①同基线（对照基线 = B4 基线存档）②同断言（harness `tools/b4-e2e/` 与断言不变）③漂移方向更严格（反向不满足）④四态不变（`verdictClass`/`chainState`/`taskState`/`stepState` 与 B4 基线逐字一致）⑤成因可追溯（如 P61 的 `02d5797`）；`:48` — 五要素齐 ⇒ 前提①视为满足；`:49` — 「豁免成立不等于 AT-23 条件闭合——真实候选腿的前提②仍须满足；豁免须由 **⑦（独立终审）** 判定，主控不得自行判定」
- `:51` — R-08.6 标题「BYOK 平行路径」；`:53`–`:55` — 「(a) 主路径 / (b) 独立计费（不占 Copilot premium request）」「Copilot-hosted 仅用于 1 次 GPT-6 Luna 连通性测试，不计入验收结论」；BYOK 配置 `COPILOT_PROVIDER_TYPE=anthropic`+`BASE_URL=https://api.deepseek.com/anthropic`+`MODEL=deepseek-flash`+两 token 上限；`:56` — 「步骤 0 必做 BYOK 可行性最小验证：`copilot --no-auto-update -p "say hello"`；失败则停下升级用户裁决」；`:57` — 「BYOK 结论不能替代 B1 的 Copilot-hosted 基线」
- `:61`,`:63` — R-09.1「工具使用 unknown = (c) 视情况」；「归还后产物扫描能证明『产物在 workspace 内、无外部副作用』则接受；发现 workspace 外痕迹则阻断」
- `:65`,`:67` — R-09.2「资源读取 unknown = (c) 视情况」；「归还后扫描能证明『未越界读取』（manifest/diff 比对只读了 workspace 内文件）则接受；发现 workspace 外读取痕迹则阻断」
- `:69`,`:71` — R-09.3「网络 unknown（底线）= (c) 视情况」；「回环可接受、公网必须阻断；起跑前一次性确认，不阻断起跑」
- `:73`,`:75` — R-09.4「费用 unknown = (a) 显式标 unknown」；「CONTRACTS §7.1 已要求 unknownDimensions 必含 cost」
- `:77`,`:79` — R-09.5「后台进程 unknown = (a) 可接受」；「用户负责清理；报告明写」
- `:81`,`:83` — R-09.6「会话连续性 unknown = (a) 可接受」；「黑箱模式设计边界」
- `:85`,`:87` — R-09.7「外部副作用 unknown（底线）= (c) 视情况」；「用户显式确认『已接受未知外部副作用』则接受；没有确认则阻断。起跑前签署确认记录，不阻断起跑」
- `:89`,`:91` — R-09.8「强制字段 = (a) 七项全必填」；「assuranceLevel=reduced 且 unknownDimensions 非空且含 cost」
- `:93`,`:95` — R-09.9「判定口径 = (b) 每项独立判定，不累积」；「某项 unknown 可接受不整体放行；某项阻断不整体失败」
- `:97`–`:101` — 「AT-24 报告三句话（硬性）」：①「本报告的保证级别是 reduced，不是 full」②七项 unknown 接受边界逐项列出 ③「本报告不构成对黑箱 Agent 过程正确性的保证；ProofRail 只保证归还后的产物扫描、独立门禁、证据和评审」
- 与 R-3 关系：本文件全文无「R-3」字面；头部仅以 `:1` 标题「AT-23 / AT-24 前置裁定」承载与 AT-23 的前置关系——未找到

【⑤ P74 `Sandbox` 抽象层与 S2 真实 enforcement】源：`docs/validation/sandbox-interface.md`、`docs/validation/t027b-admission.md`

- `sandbox-interface.md:1`,`:3` — 标题「ProofRail SANDBOX-INTERFACE 报告（P74 · 层 1 主线 B 段引入 · 实现类 [SLICE] · M 档）」；「日期：2026-11-15。状态：⑧ 文档装配定稿（①–⑦ 已闭环）」
- `:8` — 定义与作用范围：「`internal/guard/sandbox.go`（`Sandbox` 检查层接口 + `Verdict` 三态 + 双哨兵 + `CheckSandbox` 门函数）、`sandbox_windows.go`（形态 α）、`sandbox_unsupported.go`（`!windows` fail-closed）、`sandbox_preflight.go`（注入点+`preflightSandbox`+`sandboxEnforcementFromEnv`）、两条派发路径的装配检查点（均在 permit 之前）、CLI 触发源注入、`CONTRACTS{,_EN}.md` §7.2 契约段」
- `:9` — 「检查层落在 permit 认领之前，`Denied` / `Unknown` / 检查错误全部 fail-closed；S1 默认 `off` ⇒ 默认控制流与引入前逐句等价；`PRFRAIL_SANDBOX_CHECK=enforce` 时两条派发路径均被拒绝（S1 平台实现恒 `Unknown`）。**B 段不解锁 S-G1 / S-G8**」
- `:40` — 「检查点：方案 (i) —— permit 认领之前、Step 2 run-root marker 之前。理由：permit 配额一次性不可退」
- `:43` — 停升级边界 7 条「逐条裁决全部不触（不改 §7 之外契约段 / 不改 chain·adapters 判定本体 / 不改 schema·fixtures·CI·候选版本·R-08·R-09 / 不引入 P′ / 不扩域 / 不改 S1 Exit 八门 / 不改变既有拒绝语义）」
- `:49` — 契约落点：「§7 末尾新增 `### 7.2 Sandbox 检查层接口`…八条 (a)–(h)：检查层定位 / 裁决三态 / 错误分类 / 检查点位置（方案 i）/ 触发面（S1）/ 平台实现 / 边界（不做）/ B 段未解锁声明」
- `:105` — 「**B 段不解锁 S-G1 / S-G8**：本接口只落检查层抽象，`S1 Exit` 八门判定不变；真实 enforcement 边界留 S2」
- `:106` — P73 附条件三条落实：「(a) 接口与调用点同片落地；(b) 死代码禁令处置——由 ADR-044 + B2 指针 + v1.69 处置，本片以环境变量提供可达入口，`SandboxEnforced` 非死代码；(c) 报告明写『B 段不解锁 S-G1 / S-G8』」
- `:107` — 「`freshStart` 分支在 `enforce` 下**恒被拒绝**（平台实现恒 `Unknown`）…『enforcement 未接入时显式要求即拒绝』；真实接入由 S2 负责。**再入路径（非 `freshStart`）不检查、不消费配额**」
- `:133` — 「层 1 B 段主线状态：抽象层已落，真实 enforcement 边界仍在 S2；S-G1 / S-G8 未解锁」
- `:215`,`:216` — ⑥ 扫描转录：「边界不碰 chain/adapters 判定本体，也不落真实 enforcement，符合 (g)」「B 段不解锁 S-G1/S-G8，符合 (h)」
- `t027b-admission.md:7` — 「**OB-55 裁定 = 方案 b**（提前落平台无关 `Sandbox` 接口），与 R-06（2026-11-12）同向，属确认与细化；三条附条件全部接受」
- `:19` — 「B2 §8 第 8 条明确『S1 不实做』`Sandbox` 接口，理由 = 提前落地会产生无人调用的死代码」
- `:24` — 「R-06 已把『方向』定在方案 b，其节奏条件（先补 R-1）已由 P69 收口」
- `:31`,`:32` — OB-55 选项②「提前落 `Sandbox` 接口：契约先行 + 平台无关接口 + Windows 包装 + fail-closed 空实现 + 同片调用点」；选项③「维持 S2：冲突留到 S2，届时须从零补接口」
- `:43` — 裁定 5：「『B 段不解锁 S-G1 / S-G8』写报告与 `DEV_PLAN`，`s1-exit` 不改」
- `:53` — 口径：「**B 段不解锁 S-G1 / S-G8**——门级判定不因接口落地改变（S-G1 依据 = AT-23 阶段 1 条件闭合，S-G8 与 S-G1 同源）」
- `:57`,`:58` — 「候选版本：维持 GitHub Copilot CLI **1.0.83**（R-08.1 冻结）；`Sandbox` 抽象层不依赖也不改变候选版本」「Windows / P：维持 R-08.3 的『只 P』…**不引入 P′**」
- `:64` — 「拒绝路径：语义等价，拒绝理由措辞由『无 enforcement 记录』变为『Sandbox unsupported』，fail-closed 行为不变」
- `:69` — 证据缺口：「S-G1 / S-G6 / S-G7 不因 B 段获得任何解锁证据」
- `:74`,`:78` — 「P74 实施片：M 档；探针需求 = 0（接口与空实现不调用真实候选）；R-08.2 阶段 1『逐次声明』探针额度仅用于 R-3 重评，另案处理」「AT-23 重评（R-3）与 B 段**无强制先后**：R-3 不依赖 `Sandbox` 接口，B 段落不落都不影响 R-3 可达性」
- S2 真实 enforcement 的接入条件/前提：两文件均无显式条款，仅归属表述（sandbox-interface `:105`「真实 enforcement 边界留 S2」、`:107`「真实接入由 S2 负责」、`:133`「真实 enforcement 边界仍在 S2」；t027b `:32`「冲突留到 S2，届时须从零补接口」）——未找到

实得字符数：5808

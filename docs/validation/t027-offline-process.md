# T027 · A6 · 离线 pinned CLI 运行 — 验证报告

日期：2026-09-15（Linux 僵尸语义回归修复 2026-09-16，见 §14）。结论：`A6 完成`；`c2d2819` 已提交并推送 `origin/main`，修复提交 `6db5a10`，CI 全绿（run `35076045445`）。⑤ 标准门禁与原生实验全绿（详见 §9、§10）。④ 试点（Haiku + MAI）：组合可用、能发现真问题，但存在**非零残差**（本片漏 1 项 High，由用户同轮授权的追加 Codex 盲审捕获、经 4 条变异证实并整改；详见 §7.1/§7.2）——建议组合作常规切片的低成本补充扫描层、硬门切片 A7/B2/B3/B4 保留 Codex 终审，最终由用户决定。变异审计 **52 项全 RED / SURVIVED 0**（§8）。
切片目标：为 chain 提供一次**可证明的离线外部运行**——用代码常量固定版本的外部 CLI，绑定运行时长与宽限，停止整棵进程树，采集五枚冻结事实，并把结果映射为可被 chain 直接消费的终局（completed / failed / cancelled / operator-action-required / uncertain），任何缺口与未证明的停机都不得以 completed 表达。

## 1. 交付范围

- 新增 `internal/adapters/agent_runner_process_registry.go`：live 句柄注册表（launchID → 句柄，阶段阶梯 starting→running→stopping→stopped，`stopped` 为终态）、不透明 `ProcessID`（规范 JSON `{"pid":N,"startToken":"..."}`）、身份镜像 `managed-process.identity.json`（与 console 停机路径同一文件与格式）、启动槽 `managed-process.launch.slot`。
- 新增 `agent_runner_timeout.go`：`AgentRunnerTimeoutManager` 看门狗。**看门狗独占唯一截止时间**（派生 context 只用 `WithCancel`，不再派生第二个 deadline）；超时与取消同刻到达时有界落定窗口（`agentRunnerTimeoutSettleWindow`）保证取消获胜。
- 新增 `agent_runner_evidence.go`：证据采集（`logs/stdout.log`+`logs/stderr.log` 原始字节与完整性标记、`events.jsonl` 严格 JSON(L) 逐行解码、`usage.json` 封闭形状 `{calls,tokens,durationMs}`、`manifest.pre.json`/`manifest.post.json`/`diff.json` 由端口注入）、`FrozenFacts()` 缺口门槛。
- 新增 `agent_runner_pinned_cli.go`：版本预检（经 `guard.RunManaged`，adapter 不 import `os/exec`）→ 原子取得启动槽 → 单次 spawn（`PROOFRAIL_EVIDENCE_DIR` 指向本 run 的证据目录）→ 身份镜像 → 注册句柄；`StopAgentRunnerProcess` 按所有权解析身份；`WaitAgentRunnerProcess` 在判负后等待运行侧完成自身对账（有界）。
- 新增 `agent_runner_run.go`：编排与终局映射（U11 顺序）、结算映射、自然退出后的进程树复验。
- 改动 `internal/guard/process.go`：拆出 `(process) Run` 复用同一对账路径；取消路径的停机验证由无界 `context.Background()` 改为**由 grace 派生**的有界 context；`Terminate` 的验证界同样由 grace 派生且不受调用方取消影响；新增 `verify` 测试接缝。
- 改动 `internal/adapters/agent_runner_replay_root.go`（`runs/<requestId>/` 证据目录与 `RunEvidenceDir`）、`agent_runner_replay_dispatcher.go`（转发进程参数）、`internal/chain/agent_runner_launcher.go`（加法式进程参数，去除死字段 `Timeout`）。
- 新增 `tools/agent-stub/main.go`：确定性离线 CLI（`run|spawn|spawn-exit|spawn-wait|crash|no-usage|ask|hang|fail-fast|mutate-workspace|torn-events|child`，`-version` 预检），只写 `PROOFRAIL_EVIDENCE_DIR` 与工作目录，不触网。
- 测试：`agent_runner_process_registry_test.go`、`agent_runner_timeout_test.go`、`agent_runner_evidence_test.go`、`agent_runner_pinned_cli_test.go`、`agent_runner_run_test.go`、`agent_runner_native_a6_test.go` 与 `agent_runner_native_race_test.go`（均 `-tags a6native`，不进默认套件）。

## 2. 协议先行产物

- `docs/CONTRACTS.md` / `_EN.md` §7 新增“离线 pinned CLI 运行（A6）”段（冻结）：工件布局与五枚事实口径；日志/events/usage/manifest 的完整性与缺口口径；**任一缺口不得以 completed 或 DTO 空字段表达**；failed/cancelled/uncertain 的准入条件与“非零自然退出优先于缺口归类”；`exitCode` 只在自然退出时记录；有界看门狗与停机验证（由 grace 派生、不继承调用方 deadline）；判负后必须等运行侧完成对账；超时与取消的单胜者（含落定窗口）；身份镜像与 `StopProcessIdentity` 零新机制；run root 单 launch（启动槽独占创建 + 新鲜度 + 令牌所有权）；自然退出与被外部停机均需在上界内复验子进程集合（报告缺失按无子进程、不可观察 pid 按已消失，平台遏制为主保证）；版本常量与 `--print-version` 预检；所有权边界；全程离线。

## 3. 证据工件与终局分类

$$\text{事实}=\{\text{ManifestHash(post)},\ \text{DiffHash},\ \text{LogHash(原始字节)},\ \text{UsageHash},\ \text{ProcessStopEvidenceHash}\}$$

| 触发 | status | treeStatus | 冻结事实 | exitCode |
|---|---|---|---|---|
| 干净退出 + 四个 Complete + 五枚摘要 | `completed` | `stopped` | 必带 | 自然退出码（须为 0） |
| 非零自然退出（含缺口并存） | `failed` | `stopped`/`unknown` | 禁止 | 自然退出码 |
| 调用方取消 / deadline 到期 | `cancelled` | 同上 | 禁止 | 留空 |
| 看门狗超时、未归类停机、未对账 | `uncertain` | 同上 | 禁止 | 留空 |
| events/usage/logs/manifest 缺口（非零退出除外） | `uncertain` | 同上 | 禁止 | 自然退出码或留空 |
| 进程树未验证消失 | `uncertain`（仅当退出码为 0） | `unknown` | 禁止 | 留空 |
| operator 请求（events 中的 `operator-action-required`） | `operator-action-required` | 同上 | 禁止 | — |

结算：仅 `completed` 且 usage 完整时为 `observed`（携带实测 calls/tokens 与 UsageHash 证据），否则 `unknown`；非 completed 收据必须同时引用停机证据摘要与原因摘要（排序去重）。

## 4. 停机、取消与超时

- 停机只经已持久身份（句柄自身身份；无句柄时经镜像 + run root 所有者校验）走既有 `StopProcessIdentity`，重复停机幂等为 `already-stopped`；自然退出同样补一份 `already-stopped` 证明。
- 取消路径的验证 context 由 grace 派生（不再是无界后台 context）；`Terminate` 的验证界与调用方 deadline 解耦。
- 看门狗判负后，launcher 在 `agentRunnerSettleWait(grace)`（grace+500ms，封顶 10s）内等待运行侧完成对账；未对账按 `ErrAgentRunnerLaunchUnsettled` + “unknown”占位证据返回，且该 launch 仍走持久停机路径。
- **Linux 僵尸语义（CI 发现并修复，§14）**：存活判定不得把“已退出但尚未被父进程回收”的进程当作存活——它的 `/proc` 条目、start token 与进程组都还在，但已不可能执行或写入任何东西；把僵尸读为存活会让**不持有子句柄的停机路径**（按持久身份停机）在 grace 用尽后把已完成的停机误报为 `uncertain`。

## 5. 启动槽与所有权

- run root 单 launch 由**独占创建**（`O_CREATE|O_EXCL`）的启动槽保证；槽内容为 `<launchId>\n<per-claim token>`；新鲜度窗口 10s 内一律拒绝接管；超出窗口仅在身份镜像证明前任已消失（或从未发布）时可回收；镜像不可读时 fail-closed。
- 释放采用 `rename → .releasing` 的原子先占，再按令牌比对：仅当仍属于本次 claim 才删除，否则原样放回——迟到释放不可能删除他人槽位。
- 停机按所有权解析：持有句柄 → 该句柄自身身份（外来/陈旧 launchId 一律拒绝，不越界停他人进程）；无句柄 → 校验 run root 所有者后再读镜像。
- 墙钟 mtime 的可失真性（挂起/时钟回拨）与“operator 外部停机不清槽位”均已写入契约边界。

## 6. ③ V4 Pro 预审记录（7 轮，全部闭环）

| 轮次 | 结论 | 发现 | 整改 |
|---|---|---|---|
| 1 | FAIL | High 2（调用方 deadline 被误判为 completed；events 缺口可达 completed）+ Medium 4 + Low 5 | 分类补 `DeadlineExceeded`/任意未归类错误；events 缺口纳入门槛；manifest 缺口改判 uncertain；exitCode 仅自然退出；非 completed 引用停机证据；树复验覆盖被外部停机；清死字段；grace 覆盖 |
| 2 | FAIL | Medium 2（结算超时过早释放句柄；排他性仅内存）+ Low 3 | `agentRunnerNeedsDurableStop` 路由；持久启动槽（独占创建）；`stopGrace`；events 有界前缀文档化 |
| 3 | FAIL | Medium 2（检查-写入窗口；Windows-only 夹具破坏 Ubuntu 腿）+ Low 2 | 启动槽原子化（后续轮再加令牌）；改用可移植注入接缝；句柄在镜像不可读时仍退役 |
| 4 | FAIL | Medium 1（槽位先于预检取得，慢预检放大窃槽窗口）+ Low 2 | 槽位取得后移至预检之后；测试用 `os.Chtimes` 精确命中分支；operator 不清槽写入契约 |
| 5 | FAIL | Medium 1（停机路径不绑定启动所有权）+ Low 2 | `stopIdentity` 所有权绑定；槽位 owner 记录 + 条件释放；挂起/时钟边界文档化；claim 顺序补可证伪测试 |
| 6 | PASS WITH FIXES（无阻塞） | Medium 2（外来 launchId 可停他人进程；释放读-比-删非原子）+ Low 2 | 外来 ID 拒绝 + 槽位 owner 门；令牌 + rename 先占；operator 停机分类与缺口优先级写入契约双语；补齐 gap×exit 行 |
| 7 | PASS WITH FIXES（无阻塞） | Low 2（中文契约缺 EN 已有两句；rename-aside 微窗口） | 中文契约补齐（已改）；微窗口登记为后续 operator/接管切片移交项 |

## 7. ④ 终审试点记录（Haiku + MAI）

本片按用户同轮授权以低成本组合替代 Codex 执行 ④（准则 §2.2 / §2.3 / §3.9 的偏离登记见 §12）。

| 项 | 记录 |
|---|---|
| ④a 模型/角色 | `runSubagent(agentName="quick-verifier", model="Claude Haiku 4.5 (copilot)")`，人工档位：输出预算 8000 |
| ④a 结果 | 条目数 **0**（显式给出总数并自标“可疑”，符合 §4.2 空扫描规定）；**格式偏离**：附加了散文与统计表，未严格只输出单行条目 |
| ④a 误报 | 不适用（0 条，无判定） |
| ④b 模型/角色 | `runSubagent(agentName="independent-reviewer", model="MAI-Code-1.1-Flash (copilot)")`，人工档位：思考模式 high |
| ④b 首轮结果 | `RE-REVIEW: FINDINGS`：High 1 + Medium 1（均含 `file:line`、契约原句与失败场景；Section D/E 齐备） |
| ④ 独立发现数 | **2**（Haiku 0 + MAI 2，两者交集 0；③ 已发现项被 MAI 漏掉数：0） |
| A11 抽查（实际执行） | 抽查 1：MAI 声明“新增 `ExitCode≠0` + 树未验证的用例会红”——**已实际执行**：新行首跑即红（`status=uncertain … want failed`），声明属实，并因此发现代码确与契约冲突 → 已按“非零自然退出优先”修正。抽查 2：MAI 的 Medium 声明“同刻取消/超时不收敛”——**已实际执行**决定性实验 `TestA6NativeE8`：修复前 200 轮 `cancelled=138 / timeout=62`（声明属实），加入有界落定窗口后 `cancelled=200 / timeout=0`，且 M51（窗口置零）红 |
| ④b 复审 | 首轮两项已整改后重跑：**`RE-REVIEW: PASS`**，无新回归；Section D 仍列出两处边界提示（树状态翻转仅一条用例钉住；落定窗口只在单一时序剖面下压测），均为观察项、不阻塞 |
| 组合可用性 | 两模型标识均在启动前解析（探针 1/2，见 §11）；Haiku 输出未被截断；MAI 可接收输入并给出 `file:line`；**未发生**任何触发 §6 回退级别的失败 |
| (a) 组合可用 | 满足（两段均产出、MAI 输出可执行） |
| (b) 可发现真问题 | 满足：MAI 在一次调用内独立发现 1 High + 1 Medium，且两项经实际执行确认为真（Haiku 未发现，属机械扫描定位） |
| (c) 成本下降 | **试点 ④ 阶段**未调用 Codex（该阶段兜底调用数 0/2），token 由用户在 UI 回填后记入台账；④ 之后追加的 Codex 审查另计 **3 次**，见 §7.1 |

④ 与 ③ 记账分离：③ 累计 7 轮发现不计入 ④ 发现数；上表“④ 独立发现数=2”仅计 Haiku+MAI。

### 7.1 追加独立审查（Codex 盲审，用户同轮显式授权，用于试点残差评估）

④ 完成后，用户要求**额外**用 Codex 做一次独立审查，以判断低成本组合是否仍遗留问题、A6 试点是否可行。该调用**不在**试点正常路径上，也**不**由试点文档 §6 的中度/重度失败触发，属用户同轮显式授权的例外（偏离登记见 §12）；紧随其后的两次复审为“整改闭环”与“禁止自审”两项既有纪律所需，合计 3 次，逐次记录于下表。

| 项 | 记录 |
|---|---|
| 追加审查模型/角色 | `runSubagent(agentName="independent-reviewer", model="GPT-5.3-Codex (copilot)")`，盲审：只给改动范围、契约要求与只读约束，**不给** ③/④ 的发现清单 |
| 首轮结果 | `RE-REVIEW: FINDINGS`：**2 项（Codex 自标 NOVEL）**——**High 1**：`refuseForeignSlot` 对“除 NotExist 外的任何读错误”与“空 owner 内容”一律返回 `nil`，即 fail-open（与相邻 `refuseLiveMirrorProcess` 的 fail-closed 相反）；场景：重启后槽位不可读 + 外来 launchId → 继续走镜像 → 停机存活进程。**Medium 1**：声称停机存在“两套机制”，与契约“停机只能经已持久身份、走既有 `StopProcessIdentity`”冲突 |
| 主控裁定 | **High 确认并整改**：`refuseForeignSlot` 改为三分支 fail-closed（可读 ⇒ 必须 `fields[0] == launchID`，owner 为空按冲突；`os.ErrNotExist` ⇒ 先 `refuseLiveMirrorProcess()` 再判；其他读错误 ⇒ 冲突），并新增 M52/M53/M54 三条变异。**Medium 部分驳回**：读作“**一套** guard 身份停机机制、**两个入口**”——无句柄路径走 `process.platform.stop(identity, grace)`，与 `StopProcessIdentity` 是同一身份机制；若强行统一入口，会把 guard 的树复验证据换成更弱的 `already-stopped` 证明。采纳其中可证伪的部分：契约措辞双语收紧为“一套机制、两个入口、无第三通道” |
| 复审 | 整改后重跑 Codex：**`RE-REVIEW: PASS`**（High 已关闭；Medium 按上述裁定关闭；无新发现） |
| 追加审查 3（测试增量复审） | 主控按既有纪律“自身改写代码须由独立第三方复核”，对**主控在 Codex 复审之后自行改写**的测试增量（新增断言、移除一次性助手函数）再跑 1 次只读复审：**`RE-REVIEW: PASS`**，Section B **无发现**；Section D 指出 3 处未钉住不变量（缺槽阶段后未复检存活、ownerless 槽未核对 owner id、`Actions` 仅做包含判断），**已全部补钉**并重跑聚焦测试与 M52–M55（仍 4/4 RED）。该次调用用于补上“禁止自审”的合规缺口，**不是**试点兜底 |
| **残差数据（试点关键）** | 低成本组合在本片**漏掉 1 项 High**（由追加 Codex 盲审发现，并经 4 条新变异实测确认后整改）；漏掉 ③ 已发现的 Medium+：**0** |
| §6 分级回照 | 组合漏掉已被变异检验证实的问题 ⇒ 按试点文档 §6 属**轻度失败**，规定动作为“修复后重跑 MAI”。本片实际以**用户同轮授权的 Codex 复审**作为收口门（Codex 判 PASS），**未重跑 MAI**；该路径偏离与建议一并登记于 §12 |
| 计数纪律 | 追加审查的 2 项发现**不计入**上表“④ 独立发现数=2”，单列于本节，避免把试点组合产出与追加审查混算 |

### 7.2 试点判定（数据 + 建议，最终由用户决定）

| 判据（试点文档 §9） | 判定 | 依据 |
|---|---|---|
| (a) 组合可用 | **满足** | 两段均产出可用输出，未触发 §6 任何回退级别，模型标识启动前解析完成 |
| (b) 能发现真问题 | **满足，但存在非零残差** | MAI 一次调用内独立发现 1 High + 1 Medium，两项均经主控实际执行确认（见上表 A11 抽查）；**同时**追加 Codex 盲审发现组合漏掉的 1 项 High（fail-open，已整改并由 4 条变异钉住） |
| (b2) PASS 需独立佐证（A12） | **满足** | MAI 报 PASS，⑤ 有 4 项 A6 专项实验（进程树泄漏 E1 / 超时边界 E2 / 取消竞态 E3+E8 / 崩溃注入 E4）独立证明关键不变量 |
| (c) 成本下降 | **满足** | 试点 ④ 本身未调用 Codex；追加审查共 3 次（盲审 + 复审 + 测试增量复审）均属一次性评估/合规成本，非常规成本项 |
| **综合建议（主控）** | **可作常规切片的低成本补充扫描层，不建议作为硬门切片的唯 ④** | 样本量 n=1：本片组合漏 1 High（Codex 捕获）、③ 预审漏项 0。硬门切片 A7/B2/B3/B4 按试点文档 §1 继续保留 Codex 终审；若要在常规切片推广，建议补一条“组合复核后按比例抽片用 Codex 盲审”的兜底规则（例如每 N 片或每次涉及所有权/停机/身份语义时抽审 1 片），否则残差无上界证据 |

判定口径说明：试点文档 §9 的三档结论（全满足⇒可推广 / 仅 (a)(c)⇒降级方案 / 任一不满足⇒回退）针对的是**试点组合取代 Codex**；本节数据（(a)(b)(c) 形式满足 + 1 项 High 残差）支持的是**组合作补充层**，与“取代”是两个不同结论，取舍权在用户。

## 8. 变异审计（H1 硬要求）

- **第一部分 48 项变异**（每条 = 去掉一个守卫/条件下的聚焦测试），结果 **RED 48 / SURVIVED 0 / INVALID 0**，源码按 SHA-256 校验还原；脚本为临时文件（`tmp/a6_mutations_v2.ps1`），已在 `tmp/` 清理时删除。
- **第二部分 4 项变异（M52–M55，Codex 盲审轮）**：M52/M53/M54 = §7.1 的 High 整改三条分支（不可读槽位拒绝、owner 为空拒绝、缺槽 + 存活镜像拒绝）；M55 = owner 门翻转变体（`fields[0] != launchID` → `==`）。结果 **RED 4 / SURVIVED 0 / INVALID 0 / NOT-APPLIED 0**，还原经 SHA-256 校验。测试在 Codex 复审后**两次**被加强（先补两处断言，再按第 3 次复审的 Section D 补钉三处不变量），**每次加强后均重跑：仍 4/4 RED / 还原校验通过**（脚本 `tmp/a6_mutations_codex.ps1` 与加强后的复跑脚本 `tmp/a6_mut_recheck.ps1`，均为临时文件，验证后随 `tmp/` 清除）。合计 **52 项不同的可证伪主张，全部 RED**。
- **替代关系**：§7.1 的 High 整改使第一部分一条旧变异的锚文本失效（M49），已用等价新条目 M55 补上并实跑，故总数按 52 计而不重复计数。
- 覆盖：分类顺序各支（含 `TimedOut` 优先于自然失败、未归类错误、未验证树）、缺口门槛（logs/usage/events/manifest 与 `FrozenFacts` 双重拒绝）、结算 observed/unknown、exitCode 诚实性、停机证据引用、树复验（存活/不可观察/畸形报告/逐一观测）、前置 manifest 的捕获时点与接线、启动槽（独占、新鲜度、令牌、释放、被拒启动释放槽）、所有权（持有句柄身份、外来 ID、所有者门、无句柄路径）、镜像不可读拒绝、版本预检、看门狗唯一 deadline 与同刻落定窗口。
- 弃用项（附理由）：`claimLaunchSlot` 中 `os.Stat` 非 NotExist 错误分支（无故障注入不可达，另一可达镜像分支由 M43 覆盖）；令牌字符串自变体（claim 与 release 读同一令牌 ⇒ 等价，判别半面由 M49 覆盖）。

## 9. 原生实验（`-tags a6native`，真实进程）

| 实验 | 主张 | 结果 |
|---|---|---|
| E1 | 自然退出而子进程仍存活时，绝不以“completed + 活树”收场 | 10 轮：contained=10 / refused=0（Windows 由 job-close 收走子进程，收据仍经复验） |
| E2 | 超时边界两侧 | 20 轮：completed=10 / degraded=10，且无一轮让越界运行 completed、无残留进程 |
| E3 | 取消/超时竞态单胜者且无残留 | 20 轮：cancelled=20，每轮恰一个终局，stub 自身 pid 全部消失（竞态窗口 1.45–1.55s 对 1.5s 看门狗，远离预检耗时） |
| E4 | 三阶段崩溃注入（start/stop/collection） | 6 例：fail-fast→failed、spawn-exit→contained、无 marker/无 usage/撕裂 usage→uncertain、operator ask→operator-action-required |
| E5 | 版本 pin 不匹配不得 spawn | 通过 |
| E6 | 事实与结算描述真实工件 | 通过（usage 摘要与 calls/tokens 一致） |
| E8 | 同刻取消/超时收敛到单胜者 | 修复前 138/62（见 §7），修复后 200/0 |

## 10. 门禁与编码

- `gofmt -l internal cmd tools` 空；`go build ./...`、`go vet ./...` 通过；`go test -count=1 ./...` 13 个包全绿；`tools/contracts` 契约套件 4/4（含 128 夹具）。
- 编码与行尾：新增/改动 `.go`、`.json` 为 UTF-8 无 BOM + LF；`docs/CONTRACTS.md`、`_EN.md` 与本文档为 UTF-8 **带 BOM** + LF（逐字节核验）。
- 边界约束：adapter/新代码不 import `net`/`os/exec`（版本预检经 `guard`），`internal/release/network_test.go` 通过；未触碰 `.github/workflows` 与其摘要契约。
- 跨平台本地门禁补强（2026-09-16，见 §14）：新增 `GOOS=linux`、`CGO_ENABLED=0` 下的 `go build ./...` 与 `go vet ./...` 交叉编译校验；Linux-only 代码路径（`internal/guard/process_linux.go`）不再只依赖 CI 暴露编译错误。

## 11. 探针与预算

- 探针 **1/2**：目的=解析 ④a/④b 的模型标识与 `runSubagent` vendor 串；模型=`Claude Haiku 4.5 (copilot)`、`MAI-Code-1.1-Flash (copilot)`；结果=首级回落梯的 `(anthropic)`/`(microsoft)` 不在白名单，工具回显权威清单确认 `(copilot)` 写法并已写入试点文档 §2.3；**剩余 1/2**。
- Codex 调用 **3 次**（均为读-only、无子进程、无调用链；**超出试点文档 §6 的 1–2 次上限**，逐次登记于 §12）：① 试点残差评估盲审（仅给范围/契约/只读约束）；② 整改后复审 → `RE-REVIEW: PASS`；③ 测试增量复审（补“禁止自审”合规缺口）→ `RE-REVIEW: PASS`，Section B 无发现、Section D 三处未钉不变量已补钉。
- 未动用任何禁止模型（Luna/Gemini/Terra/Sol/GPT-5.4）。

## 12. 准则偏离登记（仅 A6 有效，见 `docs/t027/A6_PILOT_LAUNCH.md` §10.1）

- §2.2（Codex 必须终审）→ 本片 ④ 由 Haiku+MAI 组合执行。
- §2.3（白名单仅 V4 Pro/Codex）→ 限时放行 `claude-haiku-4.5` / `mai-code-1.1-flash`；深度仍 ≤1、无调用链；其余禁止模型仍禁。
- §3.9（纯离线切片默认 0 探针）→ 用户同轮显式预授权使用 1 次探针解析模型标识。
- 试点文档 §6（Codex 兜底上限 ≤2 次，仅中度/重度失败触发，不在正常路径调用）→ 用户同轮显式授权**追加**一次 Codex 盲审 + 一次复审（2 次）用于试点残差评估；后又按“主控自行改写代码须由独立方复核”的既有纪律追加第 3 次（测试增量复审），**合计 3 次、超出 §6 上限**。三次均非失败触发，属试点评估与合规成本，不计入兜底预算口径（如实登记，不改写为“兜底”）。
- 试点文档 §6 轻度失败处置（组合漏掉已被变异证实的问题 ⇒ 修复后**重跑 MAI**）→ 本片未重跑 MAI，改以用户同轮授权的 **Codex 复审**作为收口门（Codex 判 `RE-REVIEW: PASS`）。**如实登记为路径偏离**；若后续要在常规切片推广组合，须先补齐“漏项后由谁复核”的规则（见 §7.2 建议）。
- 试点文档 §8（报告须写明“未用 Codex”）→ 本片报告如实写明：**试点 ④ 阶段未用 Codex**，Codex 仅在 ④ 之后按用户授权做追加独立审查（§7.1），两者在报告中分列、不混算。

## 13. 已知边界与后续移交

- **rename-aside 微窗口**：释放的 rename 与 `.releasing` 比对之间存在微秒级“无槽”窗口（③ 第 7 轮 Low，判定不阻塞）；建议在 operator/接管切片一并加固（或改为文件名内嵌令牌的释放语义）。
- **operator 外部停机**：停机后槽位最多多存活一个新鲜度窗口（10s），期间新启动被拒；外部停机在运行侧与崩溃不可区分，按 `failed` 记录并保留退出码（已写入契约）。
- **events 有界前缀**：超出 1MiB 读取上界的事件请求不可见，但该前缀按不完整处置且绝不 completed。
- **Unix 孤儿**：launcher 崩溃且父进程已退出时，镜像只命名父进程，Unix 侧子进程孤儿需依赖 process group（本片仅 Windows 原生验证）。
- **未覆盖的观测项**（④b 复审 Section D）：树状态翻转仅一条用例钉住；落定窗口仅单一时序剖面压测。
- **追加 Codex 盲审遗留观测项**（§7.1，均不阻塞）：① **不可读槽位的构造不可移植**——“槽位不可读”在测试中用**目录占位**（Windows 上难以非提权复现权限拒绝），Unix 侧同形构造只在 `EISDIR` 语义下成立；② **跨平台原语仅 Windows 原生验证**——`a6native` 实验（进程 job 遏制、`taskkill` 路径）**不进 CI**（`ci.yml` 只跑默认路径与 `-race` 子集），故 job 语义、孤儿回收等结论目前只有 Windows 本机证据（2026-09-16 更新：CI 的 Ubuntu 腿已在默认套件中真实跑过 `internal/guard` 的 Linux 停机/身份路径，并因僵尸语义缺陷暴露 12 个红灯——已修复并加 Linux 回归测试，见 §14）；③ 组合残差（组合漏掉 1 项 High）见 §7.1/§7.2，属**方案层面**的已知边界，不改变本片代码结论。
- **第 3 次复审（测试增量）的 Section D/E**：三处“未钉住不变量”**已补钉**（缺槽阶段后复检存活、ownerless 槽兼测 owner id、`Actions` 断言由包含改为“恰一条且为 `already-stopped`”）并重跑验证；仍未钉住的是**失败路径**上的卫生性断言（若在测试后半段 `t.Fatal`，`t.Cleanup` 只保证进程被停，不保证 owner 句柄退役与日志句柄关闭），以及 §7.1 测试的“owner 清理只断言不报错”，二者均为低风险覆盖残差、非行为缺陷。
- **提交状态**：本片 `c2d2819` 已按用户同轮显式授权提交并推送 `origin/main`（gitee 永不推送）；Linux 僵尸语义修复与本文档同步随修复提交推送（§14）。

## 14. CI 证据与 Linux 僵尸语义回归修复（2026-09-16）

- **CI 首跑（红）**：`c2d2819` 推送后 `main CI`（run `35074301616`）Windows 腿全绿、Ubuntu 腿 **Test 失败**：12 个 pinned-CLI 用例报 `managed process termination uncertain: context deadline exceeded`，证据里 `Actions` 为 `[signal-term-process-group signal-kill-process-group]` 而无消失判定。这是本片首次在 Linux 上跑停机路径——A6 的原生验证只覆盖 Windows（§9、§13）。
- **根因**：`process_linux.go` 的存活判定只比对 `/proc/<pid>/stat` 的 start token。进程被信号杀死但父进程尚未 `wait` 回收时为僵尸：`/proc` 条目与 start token 仍在，`kill(-pgid, 0)` 也仍成功，于是**不持有子句柄的停机路径**（launcher 用镜像身份停机，自身从不 `wait`）在 grace 用尽后把已完成的停机判成 `uncertain`。Windows 无僵尸概念（句柄 + 退出码），故该腿全绿。
- **修复**（`internal/guard/process_linux.go`）：
  - `readProcStat` 同时解析 state（`fields[0]`）与 pgrp（`fields[2]`），保留 start token 的 pid 复用防护；
  - `platformIdentityAlive`：state 为 `Z`/`X`/`x`（僵尸/已死）一律视为**不在运行**；
  - `processGroupAlive`：信号 0 仍成功时再扫描 `/proc`，仅当组内存在**非终态**成员才判存活；组属于其他用户（EPERM）时保持“保守视为存活”的旧语义（`/proc` 看不到其成员）。
- **回归测试**（新增 `internal/guard/process_linux_test.go`，Linux-only）：`TestLinuxStopProvesAnUnreapedTerminatedChild` —— 子进程独占进程组、杀掉后**故意不回收**，断言 ① 运行中读为存活（对照组）；② 退出未回收时读为**不在运行**；③ 该组不再是存活组；④ `StopProcessIdentity` 给出 `stopped` 且 `Actions == [already-stopped]`；最后才 `wait` 回收，保证结论不是由父进程的回收造成的。修复前该用例在 5s 上界内失败。
- **本地门禁（Windows 主机）**：`gofmt -l internal cmd tools` 空；`GOOS=linux`（`CGO_ENABLED=0`）`go build ./...` 与 `go vet ./...` 通过；Windows `go build ./...`、`go vet ./...` 通过；`go test -count=1 ./internal/guard/ ./internal/adapters/` 通过。
- **远端 CI 证据（绿）**：GitHub Actions run `35076045445`（提交 `6db5a10`）全绿——Ubuntu 腿在默认套件中真实执行了 `Race` 与 `Contract fixtures` 两个步骤（Windows 腿依 `runner.os == 'Linux'` 条件跳过该两步），其余三个候选/引导 job（Candidate build / Candidate probe / Bootstrap and release evidence）仍按设计跳过，与上一次绿跑（A5 run `34951615184`）同形。

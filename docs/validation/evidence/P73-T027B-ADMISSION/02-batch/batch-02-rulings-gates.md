## 面 A — `docs/validation/r08-r09-ruling.md`

**3）裁定日期行（行 3）**：

```text
> 性质：`[ROUTINE]` docs-only 落地（不占切片号、不升准则版本）；裁定人 = 用户；裁定日期 = 2026-11-06；本文档 = 2026-11-07 落地定稿。**CN 为权威版**，`_EN` 为同位镜像。
```

「状态」行：不存在（全文无独立状态行，日期承载于行 3）。

**1）R-08.1 – R-08.4（行 12–34，逐字连续）**：

```text
### R-08.1 候选版本 = (a) GitHub Copilot CLI 1.0.83

- 封装函数 `Invoke-CopilotCli` / `buildCopilotArgs`；`-Unlock` 默认 `$false`（fail-closed）。
- 阶段 1：首次运行前校验 + 每 5 次探针抽样 1 次；阶段 2：每次校验。
- 三条证据对照（默认封装 / 显式参数 / `-Unlock` 裸跑）。
- 版本证据用 `copilot --version` 输出，不用 `winget list` 或 `ProductVersion`。
- 新增 copilot 调用点必须走封装；否则在 ⑤ 预审标记为偏离。
- 已知边界：锁定依赖 SEA loader 的 `Pr()` 行为；loader 之外的自动更新路径未审计。
- 锁定失败降级：转 (b) 案「版本不可控」。

### R-08.2 探针预算

- 阶段 1：逐次声明（用户同轮批准），首次不预设上限。
- 阶段 2：Copilot 5 次 / BYOK 100 元 / Copilot 货币 0 / 重试 5 次 / 追加须用户同轮批准。
- 「逐次声明」= 用户同轮批准，主控不得自行决定。

### R-08.3 Windows 前提 P = (a) 只用 P（进程级崩溃、无断电、本机 Win11/NTFS）

- AT-23 报告必须如实标注「不覆盖断电场景；断电耐久引用 B3b（P′ 覆盖物理机目标部署环境未论证）」。

### R-08.4 VM / 代理 = (a) 都不允许

- BYOK 连通性预检失败路径：停下升级用户裁决——可选 (a) 临时启用代理（须同轮批准）/ (b) 改用本机可达端点 / (c) 暂停 AT-23；主控不得自行启用代理或切端点。
```

**R-08.5（行 36–49）**：

```text
### R-08.5 真实候选腿失败处置 = (b) 条件闭合（reduced）

- 机制腿全绿 + 真实候选腿实际跑过（≥3 次，覆盖启动/执行/停机）并记录 unsupported/unknown。
- 每次 unsupported/unknown 必须有具体错误信息或证据（stderr 原文、退出码、失败阶段），不能只有「跑不通」。
- 报告必须明写「AT-23 条件闭合，不构成真实候选兼容性证明」。

- **继承漂移豁免的适用条件（五要素，须同时满足）**：
- ① **同基线**：对照基线 = B4 基线存档（`docs/validation/t027-at23-e2e.md` + `docs/validation/evidence/directive-history/b4/archive/`）；只允许以该存档为对照，不得另选或混用其它基线。
- ② **同断言**：断言 id、断言含义、断言所在 harness（`tools/b4-e2e/`）与 B4 基线均不变；harness 有改动或断言被重定义即不满足。
- ③ **漂移方向更严格**：断言值变化方向更保守/更严格（如 `process-tree-taken-away` 从 `true` 变 `false` = 未检测到进程树被取走 = 判定更保守）；反向变化（如 `false` 变 `true` = 更宽松）不满足。
- ④ **四态不变**：`verdictClass` / `chainState` / `taskState` / `stepState` 与 B4 基线逐字一致；任一字段漂移即不满足。
- ⑤ **成因可追溯**：漂移成因可追溯到 B4 之后的特定产品改动（如 P61 的 `02d5797`，A5 pid reports 绑定为第六枚冻结事实）；无法归因到特定改动即不满足。
- **豁免后果**：五要素同时满足 ⇒ 豁免「前提①不满足」的判定，前提①视为满足（继承漂移豁免成立）；任一不满足 ⇒ 不豁免，按本小节既有判据处理。
- **豁免边界**：豁免成立不等于 AT-23 条件闭合——真实候选腿的前提②仍须满足；豁免须由 **⑦（独立终审）** 判定，主控不得自行判定。
```

**R-08.6（行 51–57）**：

```text
### R-08.6 BYOK 平行路径

- (a) 主路径 / (a) 无条件 / (b) 独立计费（不占 Copilot premium request，受 BYOK 费用上限约束）。
- Copilot-hosted 仅用于 1 次 GPT-6 Luna 连通性测试，不计入验收结论。
- BYOK 配置：`COPILOT_PROVIDER_TYPE=anthropic` + `BASE_URL=https://api.deepseek.com/anthropic` + `API_KEY` + `MODEL=deepseek-flash` + 两个 token 上限（`COPILOT_PROVIDER_MAX_PROMPT_TOKENS=840000` + `COPILOT_PROVIDER_MAX_OUTPUT_TOKENS=128000`）。
- **步骤 0 必做 BYOK 可行性最小验证**：`copilot --no-auto-update -p "say hello"`；失败则停下升级用户裁决——可选 (a) 改用 openai 类型 / (b) 调整 R-08.2 上限 / (c) 暂停 AT-23。
- 报告中登记：BYOK 结论不能替代 B1 的 Copilot-hosted 基线。
```

**2）R-09.1 – R-09.9（行 61–95，逐字连续）**：

```text
### R-09.1 工具使用 unknown = (c) 视情况

归还后产物扫描能证明「产物在 workspace 内、无外部副作用」则接受；发现 workspace 外痕迹则阻断。

### R-09.2 资源读取 unknown = (c) 视情况

归还后扫描能证明「未越界读取」（manifest/diff 比对只读了 workspace 内文件）则接受；发现 workspace 外读取痕迹则阻断。

### R-09.3 网络 unknown（底线）= (c) 视情况

回环可接受、公网必须阻断；起跑前一次性确认，不阻断起跑。

### R-09.4 费用 unknown = (a) 显式标 unknown

CONTRACTS §7.1 已要求 unknownDimensions 必含 cost。

### R-09.5 后台进程 unknown = (a) 可接受

用户负责清理；报告明写。

### R-09.6 会话连续性 unknown = (a) 可接受

黑箱模式设计边界。

### R-09.7 外部副作用 unknown（底线）= (c) 视情况

用户显式确认「已接受未知外部副作用」则接受；没有确认则阻断。起跑前签署确认记录，不阻断起跑。

### R-09.8 强制字段 = (a) 七项全必填

assuranceLevel=reduced 且 unknownDimensions 非空且含 cost。

### R-09.9 判定口径 = (b) 每项独立判定，不累积

某项 unknown 可接受不整体放行；某项阻断不整体失败。
```

## 面 B — `docs/validation/p57-closeout-assessment.md`

**4）§3.3 R-1 / R-2 / R-3（行 68–70）**：

```text
| R-1 | P44 的端到端接线 | F4：判定实现与测试已完成，接线未做 | **无队列位**（v46 §四 未列） | 需用户在 ⑩ 裁定是否入队；若入队须按 §1.8.2 登 `REMAINING_SLICES` |
| R-2 | enforcement 生产落地 B 段 | F6：归属三选一未裁 | OB-55 待议 + R-06 待裁 | 用户裁定 OB-55 归属 |
| R-3 | 真实 CLI 端到端（AT-23 真实候选腿） | F7 + `docs/validation/s1-exit.md:43` | 层 1 第 8 项 `AT23-RESTRICTED-CANDIDATE` | R-08 / R-09 裁定 + 兼容候选 |
```

**5）§4.2 八门表 S-G1..S-G8（行 84–91）**：

```text
| S-G1 AT-23 真实候选腿与整体 | CLI 可真实启动经准入的 AI 候选并计费 | 不可达 / 未通过：`docs/validation/s1-exit.md:43` | R-08、R-09、T027 装配、兼容候选 | 三者齐备后重跑真实 E2E |
| S-G2 AT-24 与 T028 | `supervised-black-box` 实现完成且 AT-24 判定成立 | **口径滞后**：`docs/DEV_PLAN.md:119` 已 `COMPLETED`（AT-24 在 (b)/(c) 裁定下条件闭合），而 `docs/validation/s1-exit.md:44` 仍写 `BLOCKED` 且理由为「T028 尚未实现」 | 裁定 (b)/(c) 的落地表述 | 在 S1 exit 报告中按 (b)/(c) 重述为「条件闭合」，并保留「真实 CLI 端到端依赖 T027-B」的边界 |
| S-G3 最终 ZIP 绑定与复验 | 最终 ZIP 绑定发行记录并复验 | 未通过：`s1-exit.md:56`（`PARTIAL`）+ `docs/ADR_REGISTER.md:19`（ADR-009） | S-G2、S-G4、发行记录 | 其余门闭合后执行绑定与复验 |
| S-G4 EOL 通知方式 | ADR-010 获产品所有者批准 | 未通过：`docs/ADR_REGISTER.md:20`（`PROPOSED / NOT APPROVED`） | 产品所有者 | 批准 ADR-010 |
| S-G5 产品所有者 S1 exit 批准 | 产品所有者作出批准 | 未知（人工门）：`s1-exit.md:60` | S-G1 至 S-G4 | 由产品所有者作出 |
| S-G6 T018 整体 | 依赖字段满足且四道剩余门闭合 | 不可达：`docs/DEV_PLAN.md:94`（依赖含 T027、T028）、`:96`（`IN_PROGRESS / BLOCKED`） | S-G1 与 S-G2 与任务级依赖 | 门级与任务级**不同源**：四门清空后 T018 仍须 T027、T028 完成 |
| S-G7 S1 exit 总体（RFC §14） | 全 AT 通过且 T018 完成且 §14 全 exit | 不可达：`s1-exit.md:52`（核心任务链 `BLOCKED`）、`:56` 至 `:58`（三项 `PARTIAL`）、`:60` | S-G1 至 S-G6 | 上列全绿 |
| S-G8 真实 CLI 端到端 | 受监督路径可在真实 CLI 与真实候选下端到端运行 | 未通过：`docs/DEV_PLAN.md:119`（D-18 明写依赖 T027-B） | P5 的 R-3 | 与 S-G1 同源解锁 |
```

**6）§1 / §2 中关于「S1 Exit 是否可达」**：§1 行 14 为唯一涉及句：

```text
目标（两项，均属 §1.8.1 的**评估动作**）：① 判定 `T027-B` 是否还需拆分为正式切片；② 产出「S1 Exit 剩余门收口清单」。
```

§2（行 18–37）仅 artifacts 围栏，无该结论句；「不可达」字样见 §4.2 行 84/89/90（已引于 5）。

## 面 C — `docs/validation/s1-exit.md`

**7）门表（RFC §14）全部行（数据行 52–60，共 9 行；表头行 50）**：

```text
| 核心任务链与 fail-close 场景 | BLOCKED | 既有 AT-01–AT-22 通过；T026 完成但固定候选为 10 verified、2 unsupported、0 unknown，整体不兼容；T027 未完成、AT-23 未通过；T028 已完成、AT-24 已验收 = 条件闭合（P64，2026-11-11；[P64 报告](at24-supervised-black-box.md)）；**依赖链**：AT-23 阶段 1 条件闭合达成（前提① F2 修复片满足、前提② P61 已满足；阶段 2 不适用，用户 2026-11-09 裁定 A1）⇒ **S-G1 依据 = AT-23 阶段 1 条件闭合**（阶段 2 不适用，用户 2026-11-09 裁定 A1；依据本片 ⑦ 判定，而非 S1 Exit 整体通过）；S-G7 依赖链同步：依赖链中 AT-23 部分已阶段 1 条件闭合（阶段 2 不适用），其余门（S-G2~S-G6、S-G8）仍待，**不得声称 S1 Exit 整体通过**；**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过 |
| seed 自托管与真实平台 | PASS | T017 Windows/Linux 可信矩阵及 Windows 本机两代演练 |
| whois 只读影子 | PASS | 9/9 固定用例与 3 类漂移负例；不含生产切换 |
| generic/C/Go harness | PASS | 三套 harness 定义及三个最小用户样例均通过只读 CLI 闭环 |
| 正式安装包教程 | PARTIAL | [便携 ZIP 安装指南](../INSTALLATION.md)与模型已冻结；最终 ZIP 尚未绑定和复验 |
| 安装/升级/回滚/卸载实测 | PARTIAL | [Windows 候选便携包演练](s1-install-drill.md)通过；最终载体确定后仍须复验 |
| 支持矩阵与发行说明 | PARTIAL | [候选支持矩阵](../S1_SUPPORT_MATRIX.md)与[发行说明草案](../S1_RELEASE_NOTES.md)已建立，尚未绑定正式发行物 |
| 独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审 |
| 产品所有者最终 S1 批准 | BLOCKED | 本轮只授权 whois 夹具读取，不等于 S1 发布/阶段批准 |
```

状态分布（事实）：BLOCKED×2（行 52、60）；PASS×4（行 53、54、55、59）；PARTIAL×3（行 56、57、58）。

**8）S-G1 / S-G8 / AT-23 / AT-24**：AT-23 = 行 43、AT-24 = 行 44（AT-01–AT-24 汇总表）：

```text
| AT-23 | 阶段 1 条件闭合达成（2026-11-09，F2 修复片收口；阶段 2 不适用，用户 2026-11-09 裁定 A1）：⑦ 终审 `FINAL: PASS WITH FIXES` + 盲审四问均成立 ⇒ 前提①满足 ⇒ **AT-23 阶段 1 条件闭合达成**（前提② P61 已满足）；P63 豁免判定与 P61 PARTIAL 为历史时点结论、原文保留 | [F2 修复片结论](at23-restricted-candidate.md) §11（⑦ 判定原文逐字 + F2 三腿与 B4 基线逐字段对照：verdict 四态一致、`process-tree-taken-away` 回 `true`、`failedChecks` 空集、字段漂移仅 1 处 = f2-pos artifacts 新增 stub-child.pid 键（修复预期产物））；[P61 报告](at23-restricted-candidate.md)：机制腿与 B4 基线 0 漂移（仅 F2 一处断言继承漂移）、真实候选腿实际跑过启动 / 执行 / 停机并记录 unsupported，不构成真实候选兼容性证明；[T026 报告](t026-agent-runner-contract.md) 的历史结论不变 |
| AT-24 | 条件闭合（2026-11-11，P64 验收；⑦ 复审 RE-REVIEW: PASS） | [T028 黑箱模式](t028-supervised-black-box.md)：T028 已完成（`DEV_PLAN.md:119`），visible 黑箱候选的风险确认、隔离投递、显式归还、全量扇描与 reduced/unknown 报告已实现；闭合依据为 (b)/(c) 裁定，P64 按该口径验收并判条件闭合（[P64 报告](at24-supervised-black-box.md)）；**真实 CLI 端到端仍依赖 T027-B**（D-18）；黑箱子项阶段归属 S2（P68 b1），条件闭合不撤销 |
```

S-G1 / S-G8 另见行 62 中句：

```text
`S-G1`…`S-G8` 编号体系定义于 `docs/validation/p57-closeout-assessment.md` §4.2，本文件不引入该编号，亦不建立其与门表行名的映射。
```

## 面 D — `docs/validation/remaining-rulings.md` §4（R-06）

**9）选项 b（行 72）与用户裁定（行 78）**：

```text
- **方案 b（提前落平台无关 `Sandbox` 接口）**——预算 S/M 档；时间 1–2 片；风险 高（与 B2 §8 第 8 条死代码禁令冲突，须先修该条）。
```

```text
- **用户裁定（2026-11-12，同轮）：方案 b**（先补 P44 端到端接线 R-1 收口，再起 T027-B）；**R-1 补全本身由后续切片承载，P68 不推进**。
```

**10）OB-55 表述（行 65–66）**：

```text
- `docs/DELIVERY_DIRECTIVE.md:1091`：OB-55 处置列「**待议**：B 段立项时由 ① 阶段评估（独立切片 / 提前落 `Sandbox` 接口 / 维持 S2）」。
- `docs/DEV_PLAN.md:151` 立项登记；`:153` R-2 = OB-55 归属待 R-06；`:161`/`:163` P37/P38 基线；`:177`–`:185` P39–P44（`:185` P44 端到端接线为后续切片）。
```

## 面 E — `docs/S1_SUPPORT_MATRIX.md`

**11）行 13**：

```text
| CLI Agent / AgentRunner | 规划、未交付 | 目标正式 AI 执行通道；须固定版本/配置、能力探针、隔离 workspace、进程/会话证据与独立验收；T026/T027/AT-23 阻断 S1 |
```

检索（大小写不敏感）：`enforcement` = 0 命中；`Sandbox` = 0 命中；`T027-B` = 1 命中（行 14）：

```text
| VS Code / Copilot Chat / SessionBridge visible 黑箱候选 | 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布 | 核心离线 CLI 不要求；仅在用户显式确认风险并监督时桥接隔离候选，归还后独立验收；不保证工具/网络/费用/外部副作用，AT-24 条件闭合，S1 仍因 T027/AT-23 阻断 |
```

## 面 F — `docs/validation/B2-EXTERNAL-ENFORCEMENT.md`

**12）§8 第 8 条（行 184–186）**：

```text
8. **S2 待做（只登记，不实做）**：定义平台无关的 `Sandbox` 接口（`internal/guard/sandbox.go`）+ Windows 实现包装
   （`sandbox_windows.go`）+ fail-closed 空实现（`sandbox_unsupported.go`）。**S1 不实做**，仅在 S2 启动时按此接口扩展；
   在 S1 提前落地该接口会产生无人调用的死代码，反而增加维护面。
```

输出字符数: 约 8000
BATCH-02 DONE

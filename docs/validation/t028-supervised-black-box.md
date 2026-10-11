# 切片验证报告 · `T028`（supervised-black-box 显式模式）

> 切片标识：`T028`；准则 v1.27；日期 2026-10-11。分级 `[PILOT]`；规模 M-L；探针 0；`[ROUTINE]` 不适用。
> **状态**：**⑧b 定稿 + 提交与 CI 回写**（⑨ 的本地验证与 §7.7 的 A12 专项独立实验见 §7；⑧c 收尾记录见 §11；提交与 CI 事实见 §9）。**授权边界**：commit 与 push 经用户**同轮打包授权**执行（origin-only / SSH）；gitee 否；探针 0（不追加）。

## 1. 目标与范围

- 切片性质：功能实现片，分级 `[PILOT]`，规模 M-L，**探针 0**（零探针预算：可见通道与候选均以 stub 驱动；真实通道需求按 §1.7.2 升级，本片未触发）。
- 落地目标：实现 supervised-black-box 显式模式——版本化风险/责任提示 + 绑定确认 + 隔离 workspace + SessionBridge visible 投递 + 人工监督与显式归还 + 后置全量扫描 + 独立 gates/review + reduced/unknown 报告；**AT-24 完整闭合**（条件读法见 §7.2 与 §8）。
- 四条禁令（每条须可机械验证）：① 不自动从 AgentRunner 降级；② visible 回执不冒充完成（不产生 task PASS）；③ 风险确认不豁免安全门禁；④ 未知费用/进程/副作用不显示为安全。
- 范围外（边界明写）：真实 SessionBridge visible 端到端（依赖 T027-B 接线）、真实外部 agent 的现场观测（网络/费用/进程/副作用）、`target.schema.json` 注册表级 glob 比对、S1 发布面、`.github/workflows/**`、`go.mod`。
- 起跑裁定：PILOT 确认；⑦ 预算 `1 终审 + 1 复审 + 1 盲审`；`postflight.go` / `models.go` 触碰判为加法式、属 §1.7.1 执行自由度。

## 2. 变更概要

- 变更集 = **提交 `b095cd2`**：**⑨ 时点 28 项**（`+6998/−22`、12 新增 / 16 修改）；**定稿态 34 项**（`+7544/−24`，另含报告对、台账注记与 `DEV_PLAN` 回写，见下表末三行），0 未跟踪；其后另有**一次 docs-only CI 证据回写**（§9）。全部工件以 `git add <路径>` **精确暂存**，未使用 `git add -A`。
- 口径定义（以 `git show --stat` / `git diff --shortstat` 核对）：**功能变更集** = 代码 / 契约 / schema / 测试 + 本片证据清单，**28 项 / `+6998/−22`**；**主提交 `b095cd2`** = 功能变更集 + 收口文档（两份报告、两行台账、两处 `DEV_PLAN`），**34 项 / `+7544/−24`**（CI `changed=34` 与此一致）；**docs-only CI 证据回写 `049b71f`** = 4 项 / `+18/−16`。

| 区域 | 工件 | numstat |
|---|---|---|
| 契约 | `docs/CONTRACTS{,_EN}.md` | +16/−1（各） |
| 契约 | `schemas/chain.schema.json` | +33/−1 |
| 契约 | `schemas/supervised-risk-confirmation.schema.json` | +166/−0（新增） |
| 契约 | `schemas/supervised-report.schema.json` | +160/−0（新增） |
| 契约 | `testdata/contracts/**`（4 项） | +605/−0 |
| 契约 | `tools/contracts/{contracts.test.js,README.md}` | +3/−3 |
| 链 | `internal/chain/supervised_{blackbox,risk_confirmation,report}.go` | +819/−0（3 新增） |
| 链 | `internal/chain/{models,router,engine,postflight}.go` | +197/−12 |
| 链 | `internal/chain/supervised_blackbox_test.go` | +896/−0（新增） |
| 控制台 | `internal/console/{runtime_supervised.go,supervised_cmd.go}` | +2102/−0（2 新增） |
| 控制台 | `internal/console/{cli,config,preview}.go` | +115/−4 |
| 控制台 | `internal/console/runtime_supervised_test.go` | +1549/−0（新增） |
| 适配器 | `internal/adapters/sessbridge_visible{,_test}.go` | +290/−0（2 新增） |
| 证据 | `docs/validation/evidence/T028-freedom-list.md` | +31/−0（新增） |
| 报告 | `docs/validation/t028-supervised-black-box{,_EN}.md` | +266/−0（新增） |
| 报告 | `docs/t027/REMAINING_SLICES{,_EN}.md` | +6/−0（完成行 + 分类边界注记） |
| 计划 | `docs/DEV_PLAN{,_EN}.md` | +1/−1（T028 条目勾选 + 报告路径） |

- 核心语义要点：归还即 handoff `complete` 回执（不新增第二枚归还记录）；新增转移 `supervised-return-accepted` 与 `supervised-return-completed`；开窗即持久化绑定确认与**策略快照**。
- 契约要点：`docs/CONTRACTS{,_EN}.md` 新增 `### 7.1 supervised-black-box` 显式模式；两份新 schema 强制 `assuranceLevel = reduced` 与 `unknownDimensions` 含 `cost`。

## 3. 执行流水线与环境

- 流水线（§5.1）：① 架构 → ②a 契约 → ②b 实现 → ③ 测试 → ④ 门禁 → ⑤ 首轮预审 → ⑥ 独立扫描 → ⑦ 盲审 → ②b 修复轮（②b-3 … ②b-8）→ ③ 重跑 → ④ 重跑 → ⑤ 复审 → ⑥ 终扫 → ⑦ 终审 → ⑧a 本报告 → ⑨ 原生验证 → ⑧b 定稿 → ⑧c 收尾 → ⑩ 停点。
- 角色与模型：① `deep-reasoner` / `DeepSeek V4 Pro (deepseek)`；②a / ②b / ③ / ⑧ 由 `DeepSeek V4.1 Flash (deepseek)` 承担；⑤ `GLM-5.3 (glm)`；⑥ `MAI-Code-1.1-Flash (copilot)`；⑦ `GPT-5.3-Codex (copilot)`、**不降级**。
- 环境：Windows 本机；仓库 `d:\LZProjects\prfrail`（多根工作区中的 `prfrail`）；门禁入口 `node tools/gates/gate.js`；Go 1.22 工具链。

## 4. 异常与兜底记录

- **D-14（正向数据点）**：②b 自报全绿，③ 独立判定 AT-24 子句 4 **未闭合**（缺 scope / 副作用扫描）⇒ ②b-3 补齐后重跑 ③。该事件计入缺陷自捕获率的**我方发现**分子。
- **D-16 载体中断与审计式恢复**：②b-5 首轮因 `UND_ERR_SOCKET` 中断（已写 7 项、`+421/−40`）⇒ 主控先取证 `partial-state.txt`，再以「先审计后修」重派；审计发现 F-3 半成品且入口即 RED、控制台侧缺口真实 ⇒ 修复并以 SHA256 等值证明变异还原。
- **D-17 环境事故**：系统代理 ON（`http://10.0.0.246:8080`）致 DeepSeek / GLM 载体不可用（直连可达、代理路径为根因）⇒ 用户同轮关闭代理后恢复；⑥ / ⑦ 走 Copilot 通道、未受影响。
- **A-1 ⑤ 第 3 次调用输入未收窄（如实声明）**：该轮派单早于用户「只核 ②b-7」的收窄指示 ⇒ 实际递送完整最终态输入包；后果仅为在册/已裁定项的重述与两条 Low，未造成重复探针之外的损害。
- **A-2 主控 ④ 首次重跑被终端误判中断**：一次长命令链被误判为 `^C` ⇒ 拆分为短命令**全部重取**读数；产物复核齐全。
- **A-3 只读载体的证据形态**：⑥ / ⑦ 载体只有只读工具、无执行能力 ⇒ 全部读数由主控**材料化**进输入包（`head-baseline.txt` / `staged-diff.txt` / `gate-*.txt`）；本片据实只主张 ⑥ / ⑦ 的**文本级**复核。
- **A-4 A12 脚本首轮自伤判据（如实登记）**：首轮 5 项失败**全部**系脚本自身判据缺陷（4 项误设 JSON 路径——记录体嵌在 `$defs`，schema 本体正确；1 项误判解码函数），修正判据后 `A12_EXPERIMENT_FAILURES=0`；该轮**未**因判据修改任何工件。
- **A-5 报告撰写时序**：本报告在 ⑨ / ⑧c 读数**取得之后**一次成文（而非先写空值再回填）⇒ 与 §5.1 的「⑧a → ⑨ → ⑧b」字面次序有偏差；**不得作为先例**。
- **A-6 报告自身计入变更集**：两份报告、台账回写与 `DEV_PLAN` 回写使最终态暂存项由 28 增至 **34**；④ 读数（28）与 ⑧c 读数（34）分列于 §7.1 与 §11。
- **A-7 ⑧b 遗漏 `DEV_PLAN` 回写（用户审核发现）**：报告初版把 `DEV_PLAN` 判为「不适用」，但 §11.2 要求切片收尾必须回写 `DEV_PLAN`，OB-12 的例外仅适用于准则治理切片 ⇒ 本片（代码切片）不适用该例外。**已更正**：`DEV_PLAN{,_EN}` 的 T028 条目勾选完成并写入报告路径，本报告 §15 同步更正。
- **A-8 §17 产物栅栏的 `local-only` 口径（用户审核发现，属部分缺口）**：`tmp/t028/` 中被受版本控制文档引用的过程件按 §1.5.1 不得归入 ③；归档职责属「tmp/ 生命周期执行」片 ⇒ 本片只登记候选 **C-15**（§13），本片不改归档。

## 5. 审查结论摘要（⑤⑥⑦）

- **⑤ 第 1 读**：`PRE-REVIEW: PASS WITH FIXES`（1 High + 4 Medium + 3 Low；边界合规、契约三层一致、P1–P4 成立）。
- **⑤ 复审（第 2 次）**：F-1…F-8 判定为 7 项 CLOSED、F-5 PARTIALLY CLOSED；代码 / Schema / 契约层**零新发现**；NEW-1 / NEW-2 / NEW-3 属登记与证据完整性。
- **⑤ 复审（第 3 次，用户同轮授权例外）**：`PRE-REVIEW: PASS WITH FIXES`——漂移守卫真实且锚定正确、`abort` 豁免可接受、②b-6 两处修复通过、语义 P1 测试承重、AT-24 七子句在 (b)/(c) 下全部闭合。
- **⑥ 第 1 轮**：`INDEPENDENT SCAN: FINDINGS`，但**格式偏离**（无严重度、无 `file:line`、无「无法验证清单」）⇒ 按 §7.2 **允许重发 1 次**；其三条 findings 经核对**全为在册 / 已裁定项的重述**。
- **⑥ 第 2 轮**：`INDEPENDENT SCAN: PASS`（Section A–E 齐备，8 项 cannot-verify）；格式偏离已纠正。
- **⑥ 第 3 轮（终扫，最终态）**：`INDEPENDENT SCAN: PASS`——发现表明写**无**；判定 = 边界 / 契约 / P1 / P2 / P3 为 `PASS`，P4 与 AT-24 为**条件** `PASS`；cannot-verify 3 项；主控机械抽查其 `file:line` 引用**逐条命中**（§6）。
- **⑦ 盲审（不附任何清单）**：`BLIND REVIEW: FINDINGS`——1 Critical（= **(b) 裁定重述**）+ **1 High（真新发现：归还期未校验冻结策略 ⇒ 策略漂移）** + 1 Medium（= **D-18 重述**）+ 1 Low（P1 词法检查可被等价改写绕过）。
- **⑦ 终审（含自由度清单）**：`FINAL REVIEW: PASS`——**新增可执行缺陷 0**；High 判为**真闭合**；F-1…F-18 全部「授权内 / 实现选择」、**无越界**；L-1 已闭、L-2 选登记路径 (b) 可接受；⑥ 的 `PASS` 判为**有限采信**（其从未自主提出漂移维度）。
- 处置链：盲审 High 由 **②b-7** 修复（策略快照 + 三锚守卫 + 零持久写 + `abort` 豁免回归）；盲审 Low 的收敛建议由**语义级不变量测试**补齐；⑤ 的两条 Low 由 **②b-8** 处置（L-1 修、L-2 登记）。

## 6. 审查输入包摘要（含盲审隔离证明）

- **⑦ 盲审输入包**：**仅**「变更集（28 项）+ 验收口径（AT-24 + 四条禁令 + 边界）」；**未**递设计书、编排记录、自由度清单与任何先前审查结论 ⇒ 隔离成立（其 4 条发现均为独立导出）。
- **⑦ 终审输入包**：最终态变更集与材料化读数 + `PLAN-v1.md` + 编排记录（§9.1–§9.11）+ **自由度清单**（闭环轮按 §1.7.3 / OB-48 须附）。
- **材料化读数清单**：`head-baseline.txt`（36 项，逐路径 HEAD 与工作区的 sha256 / bytes / lines）、`staged-diff.txt`（331325 字节）、`gate-{tree,index,ci,selftest}.txt`、`go-test.txt`、`npm-test.txt`、`numstat.txt`。
- **主控机械抽查（不构成审查替代）**：⑦ 引用的 `TEST_STRATEGY.md:54`（AT-24 验收行）、`CONTRACTS.md:272`（§7.1 标题）、`supervised-report.schema.json:133`（`reduced`）与 `:150`（`cost`）、`runtime_supervised_test.go:644`（漂移负例）、`:746`（`abort` 豁免）、`:1376`（词法检查）、`:1417`（语义不变量）**逐条命中**。
- **消毒**：审查输入只含仓库内路径与文本；无凭据、无外部连接材料。

## 7. 可证伪性与门禁结果

### §7.1 机械读数

| 项 | 读数 |
|---|---|
| `node tools/gates/gate.js --all --scope=tree` | `changed=28`、`TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=index` | `changed=28`、`TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=ci` | `changed=2`、`TOTAL_FAIL=0`（**非本片变更集**，见 §9） |
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 335/335` |
| `go test -count=1 ./...` | `ok=14 FAIL=0` |
| `npm --prefix tools/contracts test` | `tests 4 / pass 4 / fail 0` |
| IDE 诊断（三包） | `internal/chain`、`internal/console`、`internal/adapters` ⇒ **0 错误** |
| `gofmt -l` / `go build` / `go vet` | 空 / exit 0 / exit 0 |

- 上表为 ④ 于 **②b-8 后**（最终态）实测；⑧c 在加入报告、台账与 `DEV_PLAN` 回写后的复跑读数见 §11。

### §7.2 AT-24 七子句判定

| 子句 | 判定 | 实现点与证据 |
|---|---|---|
| 启动前显示并持久化绑定风险确认 | 闭合 | `supervised-risk-confirmation` 记录 + `SupervisedRiskNoticeVersion = 2`；`confirmAndOpenWindow` 先于任何持久写；两条负例 |
| 只在隔离 workspace | 闭合 | 隔离工作区根 + `ChangeSetInScope`（复用 `gates.ErrWorkspaceEscape` 同规则） |
| 投递不等于完成 | 闭合 | `DeliverVisible` 不产生回执记录；`SubmitSupervisedReturn` 只接受 handoff `complete` 回执 + 五事实 |
| 显式归还后全量 manifest/diff | 闭合 | 归还路径捕获前后 manifest / diff 并落盘；`TestSupervisedAssemblyFullChain*` |
| 范围/秘密/文件/副作用扫描与独立 gates/review | 闭合（(c) 条件） | `runScopeScan` / `runSecretAndTypeScan` / `runSideEffectScan` / `runAfterReturnGates`；范围扫描为**有界近似**（已登记） |
| 未确认/断线/未知进程-费用-外部副作用阻断 | **条件闭合（(b)）** | 可检测者阻断；原理不可观测者标 `unknown` 并进报告。**字面读法下该条为 Critical**（见 §8） |
| 报告标记 reduced/unknown 且禁止自动降级 | 闭合 | schema 强制 `assuranceLevel = reduced`、`unknownDimensions` 含 `cost`；渲染只输出 `unknown` |

### §7.3 四条禁令（P1–P4）证据矩阵

| 禁令 | 机械实现点 | 回归测试 |
|---|---|---|
| P1 不自动降级 | `models.go` 意图守卫 + `Definition.validate` 分支；`cli.go` 互斥；`router.go` 端口分支 | 词法检查 `:1376`；语义不变量 `:1417`（合计六重） |
| P2 投递不冒充完成 | `sessbridge_visible.go` 零回执；`supervised_blackbox.go` 仅接受 `complete` 回执 | `supervised_blackbox_test.go` 的 delivery-only 拒绝负例 |
| P3 确认不豁免门禁 | 归还前 scope / 秘密 / 类型 / 副作用 / 进程 / gates 全部执行 | 副作用扫描五条阻断分支各有负例；`runAfterReturnGates` 负例 |
| P4 未知不显示为安全 | `ValidateSupervisedReport` 拒绝非 `reduced`、空 `unknownDimensions`、缺 `cost` | schema 与渲染两端负例 |

### §7.4 变异绑定

| 变异 | 期望 | 实测 |
|---|---|---|
| 归还入口移除漂移校验 | 红 | `TestSupervisedReturnRefusesDriftedPolicy` 红 |
| 给 `abort` 加回漂移检查 | 红 | `TestSupervisedReturnAllowsAbortUnderDriftedPolicy` 红，SHA256 逐字还原 |
| `ErrSupervisedScopeUnderivable` 分支失效 | 红 | 对应负例红 |
| `allows` 上溯包含子句失效 | 红 | 对应负例红 |
| F-3 写时集合重新放入 `independent-review` | 红 | 恢复轮实测入口即红 |

### §7.7 A12 专项独立实验

- 触发理由：⑥ 报 `PASS` ⇒ ⑨ 须有一项**与 ⑥ 无关**的判据独立证明关键不变量；实验脚本 `tmp/t028/verify-a12.js`（本地过程件），含静态判据与一次子进程复跑。
- 读数：`A12_CHECKS=58 A12_EXPERIMENT_FAILURES=0`（exit 0）；聚焦复跑 `go test -count=1 -run Supervised -v` 于链 / 控制台 / 适配器三包 ⇒ `81` 条 `--- PASS`、`0` 条 `--- FAIL`。
- 八个部分：① 冻结策略守卫（三锚 + 零持久写 + `abort` 豁免 + 返回入口调用）；② 契约不变量（`reduced`、`unknownDimensions` 含 `cost`、`selection = confirm`、第三变体、覆盖常量）；③ P4「未知不得渲染为安全」；④ P2「投递不冒充完成」；⑤ P3「确认不豁免门禁」（归还前六类检查全在）；⑥ 四条禁令的机械锚与测试登记；⑦ 独立复跑（测试 / `go vet` / `gofmt`）；⑧ D-18 的 fail-closed 边界确实存在。
- 首轮 5 项失败为**脚本自伤判据**（见 A-4），修正后全绿；**未**因判据修改任何仓库工件。

## 8. 已知边界与残余风险

- **D-18 已知边界**：`prfrail supervised` 的**真实 CLI 端到端在本片不可达**，受两道独立 fail-closed 边界所限：① 可见通道为 `unavailableVisibleClient` 占位；② 装配不提供 hook 定义来源，而 schema 强制 `hooksAfterReturn` 非空 ⇒ 任何 `supervised return --complete` 均以 `ErrSupervisedGateUnavailable` 拒绝。**如实结论：本片不是可运行发布**，真实端到端依赖 **T027-B** 接线。
- **(b) 条件闭合（措辞冲突）**：AT-24 子句「未知进程 / 费用 / 外部副作用**阻断**」与设计书 §8 的「可检测 ⇒ 阻断；原理不可观测 ⇒ 标 `unknown`」冲突；用户已裁定以 §8 读法为权威、措辞澄清**改由下一篇文档片**承担（优先级已提升为**下一片**）。**在字面读法下该条为 Critical**，报告不得写无条件的「AT-24 全闭」。
- **(c) 条件闭合（范围扫描为近似）**：`AllowedTargets` 是 target 注册表 ID，其路径 glob 存于 `target.schema.json`，而链定义与控制台配置**均不携带该表** ⇒ glob 级比对超出 `internal/console/**` 边界；本片以「工作区包含 + 空 allow-list fail-closed」为判据，登记为已知边界（候选 **C-10**）。
- **L-2 登记边界（摘要无边界标记）**：`supervisedHandoffPolicyHash` 以 `\n` 连接各部件后取摘要、**无字段 / 列表边界标记** ⇒ 跨相邻列表边界搬移元素可与原值碰撞；**无执法旁路**（守卫生效后运行时采用快照策略）⇒ 仅弱化检测。处置 = 登记 (b) + 候选 **C-12**。
- **可达性边界（D-12 / D-13）**：file-type 守卫只能以进程内构造的 manifest 触达；`ErrSupervisedScopeUnderivable` 经生产入口会被 `chain.New` **更早**以 `ErrInvalidDefinition` 拒（更严非更松）⇒ 双层断言。
- **(F) 残余守卫（登记为已知覆盖边界）**：`ErrSupervisedAbortOnly`、`ErrSupervisedNotOpen`、`supervisedAssemblyStep` 的三条定义校验、`processStopEvidence` 的「已退出受管身份」证明、`runAfterReturnGates` 的 nil runner / 未注册 hook、`runSecretAndTypeScan` 跳过非普通文件 ⇒ 候选 **C-9**。
- **链接腿无本地红证**：`allows` 的链接解析逃逸分支在本机恒 SKIP（需符号链接权限）⇒ 其证据依赖具备链接权限的宿主或 CI 腿。
- **C-13 / C-14 已登记**：C-13 = `complete` 在 `abort` 后仍可覆写 store-local 工件（**预存在**，链状态不受影响）；C-14 = 守卫冻结 policy 内的 hook ID，而 hook **定义**由调用方在归还时提供（当前 fail-closed 下不可达）⇒ 必须作为 T027-B 的输入。
- **C-11 已知时差**：`testdata/contracts/valid/runtime-records.json` 中本片自产的确认夹具 `riskNoticeVersion: 1`，而生产常量已为 `2` ⇒ schema 合法但与当前语义有时差，按 ⑤ 建议并入 AT-24 措辞文档片处理。
- **⑥ 的采信边界**：⑥ 的最终 `PASS` **不覆盖**「冻结策略一致性」维度（该维度由 ⑦ 盲审独立发现、②b-7 修复）⇒ 不得读作 ⑥ 独立证实了该修复。

## 9. 明确未执行事项

- **提交与推送（已完成）**：提交 `b095cd2`（`feat(supervised): add the supervised-black-box explicit mode (T028)`）已推送 **origin**（SSH，`0135f85..b095cd2`，变更规模 `+7544/−24`）；**gitee 未推**。
- **CI 双腿（已完成，首跑即绿）**：run `38137422277`（push 事件、head `b095cd2`）结论 `completed / success`。两腿 `Gates` 步骤读数 = Ubuntu 腿（`repo=/home/runner/work/prfrail/prfrail`）与 Windows 腿（`repo=D:\a\prfrail\prfrail`）均 `SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=34`、`TOTAL_FAIL=0`；包级 `ok` 行 = Ubuntu 腿 **16** 条（含 Race 腿 `internal/adapters` 与 `internal/chain`）、Windows 腿 **14** 条，两腿均含目标包 `internal/chain` / `internal/console` / `internal/adapters` / `internal/gates`；两腿 `FAIL` 行 **0**。**首跑无红**（无 §6.4 处置需求）。
- **CI 证据回写提交（docs-only）**：`049b71f`（`docs(t028): write back the commit id and both CI leg readings (no rule change)`，4 项 / `+18/−16`）已推送 origin（`b095cd2..049b71f`）；其自身 run `38137996587` = `completed / success`（**确认绿，但按授权口径不为其再起记录提交**）。
- **本地 `--scope=ci` 的时点说明**：④ 时点该作用域读到 `changed=2`（上一提交的路径集），**不是本片变更集**；真实 CI 腿已按上条实测 `changed=34`。
- **探针 0**：未发起任何付费 / 真实模型探针；真实 SessionBridge visible 通道与真实外部 agent 的行为**未观测**（见 §8 的 D-18）。
- **未执行**：`go test -race` 本机（Windows）未跑（race 腿由 CI 的 Ubuntu 腿承担）、真实 hook 定义来源装配、真实 visible 端到端（依赖 T027-B）。
- **调用计数**：① 1；②a 2；②b 8（含修复轮）；③ 3；④ 3；⑤ 3（触顶，见 §12）；⑥ 3（触顶）；⑦ 2（盲审 + 终审）；探针 0。

## 10. 偏差与自由度

- 与自由度清单同源，逐条见 `docs/validation/evidence/T028-freedom-list.md`（F-1 至 F-18）；判定 = **F-1…F-18 全部授权内 / 实现选择，无越界**（主控自报 + ⑦ 终审逐条确认）。
- F-1 / F-2：P1 解耦的实现形态与端口命名（授权内）。
- F-3：写时集合口径与链闭集的包含关系（切片预授权）。
- F-4：风险提示版本常量升为 `2`（授权内）。
- F-5：契约优先于方案书（P2 口径，见 D-1）（授权内）。
- F-6 / F-7 / F-8：`Failpoint` 接缝承载驻停、失败点文档声明为生产 API、悬空注释删除（授权内 / 实现选择）。
- F-9…F-18：扫描分界、渲染标签、装配 fail-closed 占位、候选登记形态等（授权内 / 实现选择）。
- **偏差（如实登记）**：A-1（⑤ 输入未收窄）、A-2（终端中断重取）、A-3（只读载体证据形态）、A-4（A12 首轮自伤判据）、A-5（报告撰写时序）、A-6（报告自身计入变更集）；**A-1 与 A-5 不得作为先例**。

## 11. ⑧c 收尾清洁记录

- 编码核验（**全量**，§5.1 ⑧c 口径，`git ls-files -z` 全类型逐文件）：**tracked 3431 / text 3342 / binary 89**；**非冻结主树违规 0**；**冻结证据树另计 9 项**（已登记口径：B2 的 CRLF、B3b 的 BOM ×2、`glm-eval/archive/**` 6 项）。
- 编码核验（**变更集**口径）：本片 **34 项**变更工件的逐文件 BOM / EOL 判定 **违规 0**（`.go` / `.js` / `.json` 为 UTF-8 without BOM + LF，`.md` 为 UTF-8 with BOM + LF）。
- 残留：未跟踪残留 **0**；暂存态 **34 项**（含两份报告、台账回写与 `DEV_PLAN` 回写），`git status --short` 无 `??` 行。
- IDE 诊断：变更工件 **0 错误**。
- `node tools/gates/gate.js --all --scope=tree` 于 ⑧c 时点复跑 ⇒ `changed=34`、`TOTAL_FAIL=0`；`--scope=index` 同读数；`SELFTEST: PASS 335/335`；`npm test` 4/4。`G7-a` … `G7-d` 由 `INFO` 转为 `PASS`（两份报告在作用域内，产物栅栏行全部可解析）。
- `tmp/t028/` 为一次性过程件（§1.5.1 第 ③ 类）；⑥ / ⑦ 的全部输入与读数均在该目录内可复核。
- **时序**：见 A-5（本报告一次成文，⑨ / ⑧c 读数在其之前取得）。

## 12. 成本与计量

- R7.1 逐行（字段口径 = `{role, model, stage, startedAt, reason, retryOf}`）：① `DeepSeek V4 Pro (deepseek)` 1 次；②a / ②b / ③ / ⑧ `DeepSeek V4.1 Flash (deepseek)` 15 次；⑤ `GLM-5.3 (glm)` 3 次（触顶：首审 + ②b-5 后复审 + 用户同轮授权例外 1 次）；⑥ `MAI-Code-1.1-Flash (copilot)` 3 次（触顶：首轮格式偏离 + 允许重发 1 次 + 最终态补扫 1 次）；⑦ `GPT-5.3-Codex (copilot)` 2 次（盲审 1 + 终审 1），复审额度 1 **未用**；探针 0。
- **额度例外登记（§7.8.3 / §7.5）**：⑤ 的第 3 次调用系用户同轮明示例外，**仅本片有效、不推广**；其输入**未按指示收窄**（见 A-1）。
- 承担者（均为 `DeepSeek V4.1 Flash (deepseek)`）：②a 2 / ②b 8 / ③ 3 / ⑧a 1 / ⑧b 1 ⇒ 合计 **15** 次调用。
- **成本向量只进本节，不进准则**（G6-2）。

## 13. 候选观察项评估（本片不落）

- **C-9**：为 supervised 模式的 CLI 参数与残缺 run root 守卫补负例（下一片或后续加固片）。
- **C-10**：目标注册表可达时，为归还补 **glob 级**授权范围比对，替换当前近似。
- **C-11**：夹具 `riskNoticeVersion: 1` 的版本漂移，并入 AT-24 措辞文档片。
- **C-12**：摘要加固（每列表带标签 / 计数前缀），归后续加固片。
- **C-13**：`complete` 在 `abort` 后覆写 store-local 工件（**预存在**）⇒ 后续片在 `complete` 入口先查投影非终态。
- **C-14**：T027-B 接线预警——hook **定义**须随窗口冻结（当前 fail-closed 下不可达）。
- **C-8**：全仓同类过期计数措辞清扫（前片遗留）。
- **C-15**：本片的 `tmp/t028/` 过程件（`ORCH-NOTES.md` / `verify-a12.js` / `partial-state.txt` / `head-baseline.txt` / `staged-diff.txt` / `gate-*.txt` / `go-test.txt` / `npm-test.txt` / `numstat.txt`）**已被受版本控制的报告引用，却仍标为 `local-only`** ⇒ 与 §1.5.1 的口径句冲突；归档属「tmp/ 生命周期执行」片的职责（不属本片 ⑧c），故登记为候选、与该批同批处理。
- **OB-82 抽样锚点**：本片自定义口径（「上次抽中片之后第 4 片」），建议下一片起跑正式写入 §7.7 或 §12.4。

## 14. 后续建议（含档位建议）

- 后续建议（仅建议、不登记）：C-8 … C-14 留至下一片起跑同批登记；**AT-24 措辞文档片的优先级已提升为下一片**。
- 档位建议：后续涉及代码 / 并发 / 所有权语义的切片，⑦ 仍用 `GPT-5.3-Codex (copilot)`、**不降级**；纯文档片可依 §7.1 分档降为 `DeepSeek V4 Pro (deepseek)`。
- ⑥ 的覆盖缺口（其未自主提出漂移维度）建议下一片评估**必核面清单的分层**（哪些项固有由 ⑦ / ⑨ 承担）。
- 试点指标未达标（见 §16）⇒ 建议下一片评估「自捕获率」口径与 ⑤ 的介入深度。

## 15. ⑩ 停点报告（含回写核对行）

- 执行至 **⑧b 定稿 + ⑧c 收尾 + 用户审核收口 + 提交与 CI**；⑨ 的本地部分与 §7.7 的 A12 实验均已执行；**⑩ 停点已出具**，并由用户**同轮打包授权**（commit + push + CI 观察 + 一次 docs-only CI 证据回写）解除。
- **回写核对行**（§11.2 口径）：验证报告 = **已产出**（本文件对）；`DEV_PLAN` = **已回写**（T028 条目**勾选完成** + 报告路径）；台账 = **已回写**（完成行 + 分类边界注记）；ADR = **不适用**（无模型降级、无额度例外需 ADR）。
- 授权边界（⑩ 停点当时）：**commit 否 / push 否 / gitee 否**；其后按用户同轮打包授权执行 **origin-only**（SSH）提交与推送，**gitee 否**、探针 0；提交与 CI 事实见 §9 的「提交与推送」与「CI 双腿」。

## 16. 试点说明（指标与如实声明）

| 指标 | 目标 | 实测 |
|---|---|---|
| 缺陷自捕获率 | ≥ 3/4 | **2/15（未达标）** |
| ③ 返工次数 | ≤ 2 | 2（达标上限） |
| ④ 返工次数 | ≤ 2 | 1 |
| ⑦ 无效轮次 | 0 | 0 |
| 主控兜底次数 | ≤ 3 | 0 |
| 用户发现的遗漏数 | 0 | 2 |

- **缺陷发现分布（诚实口径）**：我方发现 2（③ 拦截的 GAP-1 硬缺陷 + 恢复轮审计捕获的 F-3 半成品）；审查方发现 11（⑤ 首轮 8 + ⑤ 复审 3）；用户发现 2（CL-1 夹具计数笔误、CL-2 README 过期计数口径）⇒ 合计 15。
- **如实结论：自捕获率 2/15 低于试点目标 ≥ 3/4**，不得以调换口径的方式改写。归因：本片首次在 M-L 档引入全新模式，⑤ 的贡献集中在「契约不变式的执法点缺失」「恢复矩阵缺陷」「证据诚实性」三类语义 / 健全性问题，而这三类恰是机械门禁与实现者自检的盲区。
- **主控自检失误（不计入用户发现）**：③ 任务书沿用了首次 ③ 报告的自述计数且未核实 ⇒ 已在编排记录更正。

## 17. 产物清单

```artifacts
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/T028-freedom-list.md	repo
docs/validation/t028-supervised-black-box.md	repo
docs/validation/t028-supervised-black-box_EN.md	repo
internal/adapters/sessbridge_visible.go	repo
internal/adapters/sessbridge_visible_test.go	repo
internal/chain/engine.go	repo
internal/chain/models.go	repo
internal/chain/postflight.go	repo
internal/chain/router.go	repo
internal/chain/supervised_blackbox.go	repo
internal/chain/supervised_blackbox_test.go	repo
internal/chain/supervised_report.go	repo
internal/chain/supervised_risk_confirmation.go	repo
internal/console/cli.go	repo
internal/console/config.go	repo
internal/console/preview.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_supervised_test.go	repo
internal/console/supervised_cmd.go	repo
schemas/chain.schema.json	repo
schemas/supervised-report.schema.json	repo
schemas/supervised-risk-confirmation.schema.json	repo
testdata/contracts/invalid/core-config.json	repo
testdata/contracts/invalid/runtime-records.json	repo
testdata/contracts/valid/core-config.json	repo
testdata/contracts/valid/runtime-records.json	repo
tools/contracts/README.md	repo
tools/contracts/contracts.test.js	repo
tmp/t028/ORCH-NOTES.md	local-only
tmp/t028/verify-a12.js	local-only
tmp/t028/staged-diff.txt	local-only
tmp/t028/head-baseline.txt	local-only
```

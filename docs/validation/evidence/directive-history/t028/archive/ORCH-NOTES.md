# T028 编排记录（主控；〔PILOT〕）

> `tmp/` 一次性过程件（§1.5.1 第 ③ 类）。用途：累计本片的编排决策、实证偏差与待复核旗标，供 ⑤/⑥/⑦ 输入包与 ⑧a/⑧b 报告直接引用。

## 1. 阶段与产物

| 阶段 | 载体 / 模型 | 产物 | 机械读数 |
|---|---|---|---|
| ① 架构 | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` max | `tmp/t028/PLAN-v1.md`（234 行） | 关键符号主控实证存在（`HandoffSession`/`BuildHandoffReceiptRecord`/`OpenHandoffTransition`/`ReturnHandoffTransition`/`ErrUnqualifiedReview`/`requirePostflightQualifiedReview`/`mapDispatchOutcome`/`managedProcessesGone`/`DefaultExecution`/`StepRouter`）|
| ②a-1 契约 | `prfrail-implementer` / `DeepSeek V4.1 Flash (deepseek)` | `docs/CONTRACTS{,_EN}.md` 各 +15/−0（`### 7.1` + §2 行 + §7 行为表两行） | tree/index `changed=2`/`TOTAL_FAIL=0`；`npm test` 4/4 |
| ②a-2 契约 | 同上 | `schemas/chain.schema.json`（第三变体）、`schemas/supervised-risk-confirmation.schema.json`（新）、`schemas/supervised-report.schema.json`（新）、夹具 +15 条目（正 3 / 反 12）、`tools/contracts/contracts.test.js` 覆盖常量 36→38 / 128→143、CONTRACTS §2 句子去数字 | tree/index `changed=10`/`TOTAL_FAIL=0`；`npm test` 4/4；`G4b unknown=[]` |
| ②b-1 实现 | 同上 | `internal/chain` 8 文件 **+1619/−12**（3 新文件 + `postflight.go` 加法分支） | `go test ./...` 14 包全 `ok`；tree/index `changed=18`/`TOTAL_FAIL=0`；selftest 335/335；变异 M1–M4 红、M5 如实标注为冗余守卫 |
| ②b-2 实现 | 同上 | `internal/console` + `internal/adapters` 9 文件 **+2654/−3**（4 新文件）+ `tools/contracts/README.md` 1 行 | `go test ./...` 14 包全 `ok`；tree/index `changed=27`/`TOTAL_FAIL=0`；selftest 335/335；`npm test` 4/4；变异 M1–M4 红且 byte-identical 还原 |
| ③ 测试（首轮） | `prfrail-tester` / `DeepSeek V4.1 Flash (deepseek)` | `TEST: DONE` + 缺口清单（GAP-1 硬缺陷、GAP-1a/3/4 覆盖缺口、歧义 F） | 全量套件绿；AT-24 子句 4 判定 **未闭合** ⇒ 拦截；2 项自做变异红 |
| ②b-3 修补（②b 的修复轮，**非新环**） | `prfrail-implementer` / 同上 | `internal/console/runtime_supervised.go` **+190/−12**、`runtime_supervised_test.go` **+298/−29**（仅此二文件，未触其他包、未重构既有代码） | tree/index `changed=27`/`TOTAL_FAIL=0`；selftest 335/335；`go test ./...` 14 包全 `ok`；`npm test` 4/4；变异 M1–M4 红且 `runtime_supervised.go` SHA256 恒为 `726E4CA4…B220F42` |
| ③ 测试（重跑，§5.3） | `prfrail-tester` / 同上 | `TEST: DONE` + 残余缺口 G-NEW-1/2/3 + 三项待裁 | 套件全绿；基线 13/4/22 逐名存活；两项独立变异红 + 逐字还原（SHA256 与 ②b-3 记录一致）；AT-24 子句 ④ 判定**已闭合**（范围子扫描为**近似**）、子句 ⑤ **部分**（措辞冲突待用户裁） |
| ②b-4 修补（tests-only） | `prfrail-implementer` / 同上 | `internal/console/runtime_supervised_test.go` **+116/−0**（生产文件零改动：索引 blob ≡ 工作区 blob） | 副作用扫描五条阻断分支全部有负例；3 项新变异红 + SHA 回 `726E4CA4…`；tree/index `changed=27`/`TOTAL_FAIL=0`；selftest 335/335 |
| ③ 聚焦重跑（§5.3） | `prfrail-tester` / 同上 | `TEST: DONE`；**G-NEW-1 已闭合**、G-NEW-2 部分（链接腿本机 SKIP）；(F) 残余守卫清单 | 控制台 23 / 适配器 4 / 链 22 逐名 `PASS`（1 条 SKIP 已声明）；2 项新变异（权限位、相对链接逃逸）红 + 逐字还原；套件全绿 |
| ④ 主控集成 + 门禁 | 主控 | 输入包材料化（`head-baseline.txt`/`staged-diff.txt`/`gate-*.txt`） | 暂存 **27** 项（+6046/−20，0 未跟踪）；tree/index `changed=27`/`TOTAL_FAIL=0`；base `build`/`vet`/`test ok=14 FAIL=0`；selftest 335/335 |

## 2. 实证偏差（如实登记，供 ⑤/⑥/⑦）

- **D-1 契约优先于方案书（P2 口径）**：方案书 §4-C3 写「五事实与归还回执哈希互异」，与 CONTRACTS §7.1 的「归还回执本身就是五事实之一」冲突 ⇒ ②b-1 依契约判定：`facts.HandoffReceiptHash == receipt.ReceiptHash`（绑定），互异施加于**窗口开启所记录的确认/投递类摘要**。**请 ⑤/⑦ 重点复核该解释**。
- **D-2 `Definition.validate` 需接纳新模式**（方案书 C1 未列）：纯加法 case 标签，既有模式行为不变。
- **D-3 归还步进路径**：表内无 `WAITING_FOR_OPERATOR→PASSED` ⇒ 归还走 `WAITING_FOR_OPERATOR→RUNNING`（`supervised-return-accepted`）`→PASSED`（`supervised-return-completed`，`InputEvidence = 五事实`）；零新增状态与转移。
- **D-4 par景机制（**高关注旗标**）**：②b-1 判明「经 `engine.Run` 派发 supervised 步骤会把端口的 `ErrAwaitingSupervisedConfirmation` 当作步骤失败」⇒ ②b-2 改由引擎**已导出的 `Failpoint` seam** 在 `runStep` 的 `step→RUNNING` 之后、`Steps.Execute` 之前驻停（`errSupervisedParked`）。**该 seam 原为测试接缝** ⇒ 属「用测试接缝承载生产流程」的实现选择，已 ②b-2 自述；**请 ⑤/⑥/⑦ 评估**是否须改为显式端口/入口（如需改则属扩范围，须用户裁决）。
- **D-5 `abort` 终态 = `CANCELLED`**（非 `FAILED`）：控制台不得直写链状态 ⇒ 走既有 `RequestCancel` + `Run`；`TestSupervisedReturnRefusesTimeout` 即断言 `CANCELLED`。
- **D-6 新增 store-local 文件**：除方案书列举者，另有 `delivery.json`（跨次调用冻结确认所需）、`operator-interactions.jsonl`、`snapshot-store/`、`workspaces/`。
- **D-7 渲染标签**：已证维度渲染为 `proved: <id>`（非 `verified`），以使 P4 的「无 `safe`/`verified` 字样」断言有意义。
- **D-8 端口接线**：`chain.New` 仍要求 agent `Postflight` 端口 ⇒ supervised-only 运行提供 fail-closed 实现（`ErrPostflightUnavailable`），链在 supervised 分支下从不调用它（②b-2 已实证）。
- **D-9 方案书 B1 误标**：`schemas/chain.schema.json` 实为**既有文件**（方案书标「新」）⇒ 按「加法编辑」处理。
- **D-10 范围扫描判据（**待 ⑦ 复核**）**：`allowedTargets` 在仓内是 **target 注册表 ID**，其路径 glob 存于独立 `target.schema.json` 注册表，而 `chain.Definition` 与控制台配置**均不携带该表** ⇒ 控制台层可观测的许可物化只有隔离工作区根，故 scope scan 以「工作区包含」（复用 `gates.ErrWorkspaceEscape` 同规则）+ 空 allow-list **fail-closed** 为判据；**glob 级比对在本片边界内不可实现**（须引注册表来源，超出 `internal/console/**`）。
- **D-11 副作用分界（**待 ⑦ 复核**）**：工作区**内新建**条目 = 申报工作（不阻断）；**删除 / 类型变更 / 权限变迁 / 逃逸链接** = 副作用（阻断）。
- **D-12 file-type 守卫可达性**：`snapshot.Capture` 结构上只产出已接受类型 ⇒ 该守卫只能以进程内构造的 manifest 触达 ⇒ 负例直接测**生产扫描函数本身**（非端到端 return）。
- **D-13 `ErrSupervisedScopeUnderivable` 可达性**：经生产入口时，空 allow-list 定义会被 `chain.New` **更早**以 `ErrInvalidDefinition` 拒（更严非更松）⇒ 双层断言（纯函数 + 装配级）。
- **D-14 GAP 拦截事件（正向数据点）**：②b（②b-1/②b-2）自报全绿，但 ③ 独立判定 **AT-24 子句 4 未闭合**（缺 scope / side-effect 扫描）⇒ ②b-3 补齐后重跑 ③。本事件计入「缺陷自捕获率」的**我方发现**分子，报告 §4 如实登记。
- **D-15 歧义 F（待用户裁定）**：AT-24 原文将「未知费用 / 外部副作用」与阻断项并列，而 PLAN §8（已获用户采纳）划线为「**可检测 ⇒ 阻断；原理不可观测 ⇒ 标 `unknown`**」。实现按 §8，且在两种读法下均不输出 `safe`/`verified`。若 AT-24 字面权威 ⇒ 须回改 `docs/TEST_STRATEGY.md` 该行措辞；若 §8 权威 ⇒ 无需改文（报告 §4 登记）。**本片不自行裁决**。
- **D-16 ②b-5 首轮中断与恢复（主控编排 + 环境事故，如实登记）**：① ②b-5 首轮派单遭 **`UND_ERR_SOCKET`（DeepSeek 载体连接被中断）**，但该轮**已写入部分改动**（7 文件未暂存，+421/−40）；主控按纪律先取证（`tmp/t028/partial-state.txt`：逐文件 SHA256 + numstat + 变异残留探针，探针为空），再以「**先审计后修**」方式重派。② 恢复轮的审计发现：**F-3 处于半成品状态且两条新测试在入口即为 RED**（注释声称已移除 `independent-review`，集合仍包含）⇒ 恢复轮修正；F-1 的**链侧守卫已在**、**控制台侧缺口**才被补上；F-2/F-4/F-6/F-7/F-8 已完整。**审计式恢复避免了将半成品当作完成态**。③ 恢复轮以 SHA256 等值证明变异还原（因这些文件本就带本片未提交改动，`git diff --stat` 不能为空）⇒ 证据形态替代已如实申报，交 ⑤ 复核。④ 环境：`internal/console/runtime_supervised.go` 在编辑器内有未保存缓冲，子代理已把缓冲回置为与磁盘一致；**磁盘为权威**（主控已核 `git diff --stat` 为空 ⇒ 工作区 ≡ 索引），若编辑器仍显示脏标记应**回退缓冲而非保存**。
- **D-17 载体连通性事故（环境，阻断中）**：③…二者相继失败：DeepSeek 载体 `UND_ERR_SOCKET`、GLM 载体 `UND_ERR_CONNECT_TIMEOUT`。主控探针结论：**系统代理 = ON（`http://10.0.0.246:8080`）**；对该代理的 TCP 探测可达，但**直连**（绕开代理）对 `api.deepseek.com:443` 与 `open.bigmodel.cn:443` 均可达（直连 HTTPS 返回 401 = 端点正常响应）⇒ **代理路径为根因**。两个第三方模型载体（②① 与 ⑤）均循 WinINet 系统代理 ⇒ 在代理为 ON 时不可用。**处置待用户裁定**（代理开关属用户授权范围）；⑥/⑦ 的 MAI/Codex 载体为 Copilot 原生通道，未受影响。

## 3. 澄清项（用户 2026-10-11 要求登记）

- **CL-1 夹具计数**：实测夹具条目 **143（HEAD 128，Δ+15）**，逐文件 Δ = `valid/core-config +1`、`valid/runtime-records +2`、`invalid/core-config +2`、`invalid/runtime-records +10` ⇒ 覆盖常量 **143 正确**；主控在 ②a 小结中的「+14」系**主控计数笔误**（合并了两项 notice 反例）。产物清单与覆盖常量自洽（`schemaFiles.length=38` / `fixtureCount=143` / 语义套件 4）。
- **CL-2 README「29 schemas」**：`tools/contracts/README.md` 与 `schemas/*.schema.json` **同口径**；HEAD 实为 **36 份**（README 的 29 在本片前即**过期 7 份**），本片后 **38 份** ⇒ 未写新数字，改为「compiles every Draft 2020-12 schema in `schemas/`」（不绑数字），已在 ②b-2 落地（`bom=true crlf=false`）。报告 §4 登记该结论；**C-8 候选观察项** = 全仓同类过期计数措辞清扫（本片不裁）。

## 4. 用户裁定（2026-10-11）

- **PILOT 确认**：改以 `试点：` 起跑；试点指标五行见 `BRIEFING.md` §7；DR-FIX 复评按 §12.2 双轨转为记录项。
- **⑦ 预算 = 1 终审 + 1 复审 + 1 盲审 = 3 次**（盲审不与复审合并）；按比例抽查抽中（上次抽中片之后第 4 片）+ 必抽保守计入。
- **OB-82 抽样锚点**：本片自定义口径，须在报告如实登记，并建议下一片起跑正式写入 §7.7 或 §12.4（**本片不裁**）。
- **核心文件触碰（`postflight.go` / `models.go`）**：用户判定**不触发边界**（加法式分支 + 既有路径零改动，属 §1.7.1 执行自由度）；PLAN §9 已显式申报。
- **PLAN 开放问题**：两枚新记录、preview 映射近似、CLI 面 `prfrail supervised return|report`、超时/断线 PILOT 语义、自由度清单落盘点与 ⑦ 输入包材料化 —— 全部**采纳**；歧义 B/D/E 留后续。

## 5. 试点指标追踪（口径见 §12.1）

| 指标 | 目标 | 当前 |
|---|---|---|
| 缺陷自捕获率 | ≥ 3/4 | **我方发现已入账 1 条硬缺陷**（D-14：②b 遗漏 scope/side-effect 扫描，由 ③ 拦截）+ ②a 澄清 CL-1/CL-2（主控自查）⇒ 分子至少 3；审查方 / 用户发现待收尾统计 |
| 返工次数 ③ ≤ 2、④ ≤ 2 | — | ③ 未起跑；④ 未起跑 |
| ⑦ 无效轮次 | 0 | 未起跑 |
| 主控兜底次数 | ≤ 3 | 0（暂无兜底） |
| 用户发现的遗漏数 | 0 | 2（CL-1 计数笔误、CL-2 README 口径 —— 均由用户指出，**计入「用户发现」**） |

## 6. §5.3 闭环与重跑记录

- **重跑 ③ 的依据**：`②b-3` 属对 ②b 的修补（§5.3 的「整改后复审」形态），故**必须重跑 ③**（测试工程师），不得直接进 ④（用户 2026-10-11 提醒；HAIKU-BUMP 的 A-3 与 B4 先例同）。
- **重跑基线**：②b-2 已通过的测试集（链 19 顶层 + 8 子测试、控制台 13 测试含 2 子测试、适配器 4）必须**逐名存活且仍绿**；②b-3 只允许**新增**测试，不得改名 / 删除既有测试。
- **待 ④ 前的门禁口径**：tree / index / selftest 三读数 + 契约套件需在**最终暂存态**重取（本片当前 `changed=27`）。
- **重跑结果**：`TEST: DONE`。基线**逐名存活**（控制台 13、适配器 4、链 22）；修复未泄漏（新符号仅于 `internal/console/**` 命中）；两项独立变异（`ErrSupervisedScopeUnderivable` 分支、`allows` 上溯包含子句）均红且逐字还原（SHA256 `726E4CA4…` 与 ②b-3 记录一致）。
- **计数口径更正（主控失误）**：任务书里链侧基线「19 顶层 + 8 子测试」不可复现（实为 `-run Supervised` = 22+11）——该数字系主控沿用**首次 ③ 报告的自述计数**而未核实；两轮均无改名/删除，现集合为任何 19 项基线的超集。如实登记，计入主控自检失误（不计入「用户发现」）。

## 7. ③ 报出的残余缺口与待裁项

| 编号 | 内容 | 主控处置 |
|---|---|---|
| **G-NEW-1（中）** | `runSideEffectScan` 五条阻断分支中四条（类型变更 / 权限与可执行位 / 绝对链接 / 逃逸链接）无负例 | **本轮补测**（②b-4，仅补测试不改生产码） |
| **G-NEW-2（低）** | `allows` 的符号链接解析逃逸分支与 `newSupervisedScope` 的 `EvalSymlinks` 失败分支无测试 | **本轮补测**（同上） |
| **G-NEW-3（观察）** | scope scan 在生产路径**不可达**（`snapshot.Capture` 只产出工作区相对路径）⇒ 实为对不可信输入的纵深防御；`ErrSupervisedScopeUnderivable` 经生产入口会被 `chain.New` 更早拒 | 报告 §4/§9 显式登记可达性边界；不扩范围 |
| 待裁 (a) 链侧基线计数口径 | 已由主控更正（见 §6） | 已处理 |
| **待裁 (b) 子句 5 措辞冲突** | AT-24 字面（未知费用/外部副作用 ⇒ 阻断）vs PLAN §8（可检测⇒阻断；不可观测⇒标 unknown） | **升级用户裁定** |
| **待裁 (c) 是否接受范围扫描的「近似」性质** | 根因：`AllowedTargets` 是 target 注册表 ID，其路径 glob 存于 `target.schema.json`注册表，而 `chain.Definition` 与控制台配置**均不携带该表** ⇒ glob 级比对超出 `internal/console/**` 边界 | **升级用户裁定**（接受近似 ⇒ 登记为已知边界；不接受 ⇒ 需扩片） |
## 8. ③ 循环止损（主控裁定）与残余守卫登记

- **止损理由**：焦点 ③ 重跑已把 AT-24 七子句与 P1–P4 所需的正反例**全部覆盖**（子句 ⑤ 的措辞待你裁定）；(F) 清单剩余项均为 **CLI 参数校验**与「仅经被篡改 / 残缺 run root 可达」的纵深防御支，**不属**本片验收证据的口径（「聚焦测试含正反例」）。同时试点指标「返工次数 ③ ≤ 2」**已用满**（③ 共 3 次调用 = 首轮 + 2 次重跑）⇒ **不再追测**，改以登记方式收口。
- **(F) 残余守卫（登记为已知覆盖边界，报告 §8/§9 引用）**：生产可达但未钉住者 = `ErrSupervisedAbortOnly`（`runtime_supervised.go:247`，非法 `--outcome`）、`ErrSupervisedNotOpen`（`:270`，未开窗即归还）、`supervisedAssemblyStep` 的三条定义校验（`:283/:287/:294`）、`processStopEvidence` 的「已退出受管身份」证明（`:1007`）、`runAfterReturnGates` 的 nil runner / 未注册 hook（`:1014/:1020`）、`runSecretAndTypeScan` 跳过非普通文件（`:968`）；其余仅经**被篖改工件**可达（含 `chain/supervised_*.go` 的解码/校验分支）。
- **候选观察项 C-9**：为 supervised 模式的 CLI 参数与残缺 run root 守卫补负例（下一片或后续加固片）。
- **已知边界（链接腿）**：`allows` 的链接解析逃逸分支（`:887`）在本机恒 SKIP（需符号链接权限）⇒ **无本地红证**，其证据依赖具备链接权限的宿主/CI 腿；报告 §9 必须如实登记。

## 9. 增补（⑤ 复审 NEW-1/NEW-2 要求；本节的阶段行、F 台账与指标**取代 §1/§5/§8 中的过期条目**）

### 9.1 阶段行补录（②b-5 起）

| 阶段 | 载体 / 模型 | 产物 | 机械读数 |
|---|---|---|---|
| ⑤ 首轮预审 | `deep-reasoner` / `GLM-5.3 (glm)` | `PRE-REVIEW: PASS WITH FIXES`（1 High + 4 Medium + 3 Low；边界合规、契约三层一致、P1–P4 成立） | 只读载体（未执行命令）；基于主控材料化的输入包 |
| ②b-5 第 1 次 | `prfrail-implementer` / `DeepSeek V4.1 Flash (deepseek)` | **中断**（`UND_ERR_SOCKET`）但**已写部分改动**：7 文件 +421/−40 | 主控取证 `partial-state.txt`（逐文件 SHA256 + 变异残留探针空） |
| ②b-5 恢复轮 | 同上 | 审计→补齐：**F-3 半成品被救回**（入口即 RED）、F-1 控制台侧缺口补齐；F-2/F-4/F-6/F-7/F-8 判定已完整 | 5 项新变异红 + SHA256 逐字还原；tree/index `changed=27`/`TOTAL_FAIL=0`；selftest 335/335 |
| ②b-6 IDE 告警 | 同上 | `config.go` QF1003 ⇒ tagged switch（分支逐字保留 + 8 场景 A/B 等价探针）；`runtime_supervised.go` unusedparams ⇒ 删冗余参数（非接口方法、全仓 2 处引用） | 主控 `get_errors` 六文件 **No errors found**；聚焦测试与门禁复跑绿 |
| ⑤ 复审 | `deep-reasoner` / `GLM-5.3 (glm)` | **F-1…F-8：7 项 CLOSED、F-5 PARTIALLY CLOSED**；代码/Schema/契约层**零新发现**；NEW-1/2/3 属登记与证据完整性 | AT-24 七子句按 (b)/(c) 裁定**全部闭合**；P1–P4 维持成立 |

### 9.2 F-1…F-8 台账（⑧ 报告 §4/§5 直接引用）

| # | 级别 | 复审判定 | 落点 |
|---|---|---|---|
| F-1 | High | **CLOSED** | 链侧 `chain/supervised_blackbox.go`（`Delivery.DeliveryStatus != accepted` ⇒ `ErrInvalidSupervisedConfirmation`，零写）；控制台侧 `runtime_supervised.go` const + `confirmAndOpenWindow` 守卫（先于任何持久写）+ 两条负例 |
| F-2 | Medium | **CLOSED** | 相位 D 状态感知重放 + 持久开窗标记 + `ErrSupervisedWindowUnproven` + 三条崩溃注入测试（`ResumeOperatorInteraction` 幂等性经实读确认） |
| F-3 | Medium | **CLOSED** | 写时集合 10 项、不含 `independent-review`；链闭集 11 ⊇ 写时集合；渲染路径唯一且不做合成 |
| F-4 | Medium | **CLOSED** | 通知文本如实（明写非 OS 沙箱）+ `SupervisedRiskNoticeVersion = 2`；七项不可证明维度与 unknown 闭集一一对应 |
| F-5 | Medium | **PARTIALLY CLOSED**（代码侧正确；**登记缺失已在 §9.4 的 D-18 补齐**） | 两道 fail-closed 边界：`unavailableVisibleClient` 占位 + 无 hook 定义来源 ⇒ 真实 CLI 端到端依赖 T027-B 接线 |
| F-6/F-7/F-8 | Low | **CLOSED** | `Failpoint` doc 声明为生产 API；悬空注释删除；`ScopeClean` → `ChangeSetInScope` |

### 9.3 试点指标重算（诚实口径，⑧ 报告 §5 直接引用）

- **我方发现**：③ 拦截的 GAP-1 硬缺陷（D-14）1 条 + 主控自查的 F-3 半成品（D-16，属 ②b-5 中断所致、由恢复轮审计捕获）1 条 ⇒ **2**。
- **审查方发现**：⑤ 首轮 8 条（1H+4M+3L）+ ⑤ 复审 3 条（NEW-1/2/3，均 Medium 及以下、属登记/证据完整性）⇒ **11**。
- **用户发现**：CL-1（夹具计数笔误）、CL-2（README 过期计数口径）⇒ **2**（另有一次流程提醒不计入发现）。
- **总发现 = 15 ⇒ 自捕获率 = 2/15**，**低于试点目标 ≥ 3/4** ⇒ **必须如实登记为未达标**，不得以旧表或口径调换蒙混。归因：本片首次在 M-L 档引入全新模式，⑤ 在「契约不变式的执法点缺失（F-1）」「恢复矩阵缺陷（F-2）」「证据诚实性（F-3/F-4）」三类**语义/健全性**问题上贡献了主要发现面，而这三类恰是机械门禁与实现者自检的盲区。
- 其余指标：③ 返工 2（达标上限）、④ 返工 1（②b-5 后复跑）≤2、⑦ 无效轮次待 ⑦、主控兜底 = 0（无兜底触发）、用户发现遗漏数 = 2（>0，未达标）。

### 9.4 新登记项

- **D-18 已知边界（F-5 / NEW-2，进报告 §9）**：`prfrail supervised` 的真实 CLI **端到端在本片不可达**，受**两道独立 fail-closed 边界**所限：① 可见通道为 `unavailableVisibleClient` 占位（投递步即以 `visible_delivery_unavailable` 拒绝）；② 装配不提供 hook 定义来源，而 schema 强制 `hooksAfterReturn` 非空 ⇒ 任何合法定义的 `supervised return --complete` 都会以 `ErrSupervisedGateUnavailable` 拒绝。**交付事实**：装配级全链已由 stub 驱动证明（`TestSupervisedAssemblyFullChain*`），真实端到端完成依赖 **T027-B** 的可见通道接线与 hook 定义来源。
- **C-10（来自 (c) 裁定）**：目标注册表可达时，为 supervised 归还补 **glob 级**授权范围比对，替换当前的「工作区包含 + 空 allow-list fail-closed」近似。
- **C-11（来自 NEW-3）**：`testdata/contracts/valid/runtime-records.json` 中本片自产的确认夹具 `riskNoticeVersion: 1`，而生产常量已升至 2 ⇒ schema 合法但与当前语义有时差（今日的 v1 记录会被引擎判为版本漂移拒绝）。**按 ⑤ 建议在 [PILOT] 预算下延后**，并入后续 AT-24 措辞文档片顺带处理。
- **行号口径声明**：§8 所列行号取自**修复前**修订；修复后行号整体漂移（如 `:247`/`:1007` 现约 `:170`/`:1158` 等）⇒ **以符号名定位为准**，报告引用行号时须重新导出（不得沿用旧行号）。
- **非缺陷观察（⑤ (d) 末，记入报告 §13 候选）**：① `runSupervisedTaskPostflight` 的 uncertain 分支复用 `agentRunnerRecoveryUncertainReason`（线缆值为中性 `"recovery-uncertain"`，仅 Go 标识符名带前缀）；② `DeliveredAt` 取发送前时刻（时钟语义微瑕）；③ Go 侧 `validateHandoffPolicy` 似未强制 `HooksAfterReturn` 非空，而 schema 强制（T014 既有接缝，非本片引入）。

### 9.5 ⑥ 第 1 次（MAI）结果与处置

- 结论行 = `INDEPENDENT SCAN: FINDINGS`；但**未给严重度、无 file:line 级新发现、Section E 未给出「无法验证清单」**（不合其模板的必核面覆盖要求）。
- 其三条「findings」经主控核对**全部为在册 / 已裁定项的重述**：① 真实可见通道未端到端 = **F-5 / D-18 已知边界**（且本片零探针、以 stub 驱动是你的起跑裁定）；② 范围 / 副作用为有界近似 = **(c) 已裁定接受**；③ 措辞冲突未决 = **(b) 已裁定采纳 PLAN §8 读法**（⑥ 似未读到该裁定）。⇒ **无可执行的待修项**。
- **处置**：按 §7.2「MAI 输出格式偏离 ⇒ 允许重发 1 次」**重发 1 次**（第 2 次有效调用，额度 ≤3），并在派单中校准：三条重述须标注为「在册 / 已裁定项的重述」而非新缺陷（如坚持为缺陷须给理由与 file:line）；要求给出 file:line 级**新**发现或**明写「无」**；补齐严重度、必核面覆盖表与完整的「无法验证清单」。
- **预设兜底**（若重发仍无可执行新发现）：按 §7.2 的「⑥ 未产出可执行发现」路径——① 如实留痕；② **进入 ⑦ 独立终审并按 ⑦ 的结论闭环**；③ 由 ⑨ 追加专项独立佐证实验承担 A12 角色；**该层结论不构成审查充分性证明**。

### 9.6 ⑥ 重发结果

- 结论行 = **`INDEPENDENT SCAN: PASS`**；三条旧项均归类为 `info` 级的「在册 / 已裁定项的重述」并**明写无新发现**；Section D 必核面覆盖表与 Section E「无法验证清单」（8 项）齐备 ⇒ 格式偏离已纠正，§7.2 的兜底路径**未触发**。

### 9.7 ⑦ 盲审轮（§7.7 必抽 + 按比例抽查，**不附任何清单**）

- 载体 `independent-reviewer` / `GPT-5.3-Codex (copilot)`；输入仅「变更集（28 路径）+ 验收口径（AT-24 + 四条禁令 + 边界）」，**未递**设计书 / 编排记录 / 自由度清单 / 任何先前审查结论。
- 结论行 = **`BLIND REVIEW: FINDINGS`**（4 条）：
  - **Critical**：按 AT-24 **字面**读法，「未知进程/费用/外部副作用阻断」未被实现（实现将其标 `unknown` 并允许完成）⇒ 这是 **(b) 已裁定项的重述**（用户已裁 §8 读法为权威、AT-24 措辞澄清延后文档片）。**处置**：不改代码；报告 §4/§9 必须明写「本片 AT-24 完整闭合依赖 (b) 裁定；在字面读法下该条为 Critical」并引用户裁定。⑦ 自己给出的收敛路径 (2) 与此一致（先修条款再宣称验收）。
  - **High（**真新发现**）**：归还阶段以「返回时重新加载的 definition/policy」执法（超时 / 范围 / gates），但**未校验其与开窗时冻结的 `HandoffPolicyHash` / `AllowedReturnActions` 一致**（会话已存但未被执法）⇒ 中途换配置可致**策略漂移**。**处置：②b-7 修**（归还入口强校验 policy 哈希与冻结值相等，不等则 fail-closed；最好持久化策略快照并以快照执法）。
  - **Medium**：默认 CLI 装配为 fail-closed 占位（visible sender 不可用 + 无 hook 源）= **D-18 的重述**，附要求「若本片目标含可运行交付则需补接线；否则验收结论须明写『非可运行发布』」⇒ 报告 §9 照此写明。
  - **Low**：P1 解耦的词法证据可被等价改写绕过 ⇒ 建议保留为烟雾项并**补语义级不变量测试**（状态迁移 / 端口路由层）。**处置：②b-7 一并补**。
- **§7.7 抽样记录（进报告 §6）**：片号 = `T028`（口径 = 「上次抽中片（`DIRECTIVE-EVIDENCE-SCOPE`）之后第 4 片」，主控自定义锚点，待下一片正式写入 §7.7/§12.4）；**盲审发现数 = 4**（1 Critical + 1 High + 1 Medium + 1 Low；其中 Critical/Medium 为在册 / 已裁定项的重述，**High 为新可执行发现**）；**是否推翻 ⑥ 的 PASS = 部分推翻**（⑥ 的 `PASS` 未覆盖「冻结 policy 一致性」这一维度；但未被整体推翻）。
- **处置**：②b-7 修 High + 补 Low 建议的语义不变量测试 ⇒ 重跑 ④ 聚焦门禁 ⇒ ⑦ **终审**（含自由度清单）复核修复并做全片闭环；若终审仍报 Medium+ ⇒ 用第 3 次（复审）额度。

### 9.8 ⑤ 复审 #2（第三次调用）与 ②b-8

- **⑤ #2 = `PRE-REVIEW: PASS WITH FIXES`**（载体 `GLM-5.3 (glm)`）。判定：②b-7 漂移守卫**真实且锚定正确**（冻结时机、单一哈希约定、三检查不冗余、**快照统治后续超时/范围/门禁执法**、fail-closed 零写、真实 CLI 路径可达）；`abort` 豁免**可接受**（理由成立）但当时**无测试钉住**；②b-6 两修复**通过**（八场景 A/B 逐字节等价、删参数无信息损失）；语义 P1 测试**真语义、承重**；盲审四项台账 **Critical=(b) 裁定 / High CLOSED / Medium CLOSED（待报告陈述）/ Low CLOSED**；AT-24 七子句**在 (b)/(c) 下全部闭合**（报告不得写无条件的「全闭」）。
- **额度留痕（§7.8.3 用户授权例外）**：⑤ 调用 = **3 次**（首审 1 + ②b-5 后复审 1 + 本轮用户同轮明示例外 1）⇒ **已触顶**；授权以《用户同轮明示》登记，**仅本片有效、不推广**；④ 之后**不再跑 ⑤ 第 4 次**（与用户 2026-10-11 裁定一致）。
- **输入形态如实声明（不得读作已收窄）**：本轮 ⑤ 的**派单早于**用户「只核 ②b-7 策略漂移修复、不重审 ②b-1…②b-6」的收窄指示，故**实际递送的是完整最终态输入包**（非收窄）。后果如实登记：多出的覆盖面仅产生**在册/已裁定项的重述**与两条 Low，**未产生重复探针之外的损害**；用户的**预期收窄口径**（若后续同类片沿用）记录于此：收窄到 `runtime_supervised.go` 归还入口 + `chain` 的 `SubmitSupervisedReturn` + 相应新增测试；不重审前轮已通过改动；不要求全面偏离项扫描。
- **②b-8（最终机械清理，`LOWFIX: DONE`）**：① L-1 **已修**（`deliverSupervisedNotice` 签名去掉死返回值 `chain.SupervisedDelivery`，2 调用点 + 函数体 6 条 return 同步；改动前全仓 4 命中、两调用点均以 `_` 丢弃、`_test.go` 零直接调用 ⇒ **无消费者丢失**；`chain.SupervisedDelivery{}` 在包内出现次数 5→0）；② **abort 豁免回归守卫已补**（`TestSupervisedReturnAllowsAbortUnderDriftedPolicy`：漂移下 abort 必须被接受、终态 CANCELLED、无 PASS/review 输入；变异 = 给 abort 加回 drift 检查 ⇒ 该测试**变红**，还原 SHA256 逐字相等）；③ **可选 Nit-1 已采纳**（漂移负例的 hooks 腿注册 `hook-two` ⇒ 该腿现证明「无守卫时会跑错门禁并收敛」的强命题；仅作用于该子测试实例，不触共享夹具）；④ 读数：gofmt/build/vet 清、`go test ./...` 14 包全 `ok`、tree/index `changed=28`/`TOTAL_FAIL=0`、selftest 335/335、`npm test` 4/4。

### 9.9 L-2 处置路径（**明确裁定：选 (b) 登记，不改码**）

- **L-2（Low）内容**：`supervisedHandoffPolicyHash` 将各部件以 `\n` 连接后摘要、**无字段/列表边界标记** ⇒ 「跨相邻列表边界搬移元素」的漂移可与原值碰撞（检查 ③ 不报）；**无执法旁路**（守卫生效后运行时采用**快照** policy，超时/范围/门禁均读冻结值；`ReturnActions` 另由会话冻结集在回执构造时执法）⇒ 该局限**只弱化检测**。
- **选 (b) 的理由**：① ⑤ 自己判定「非验收必需」且「实际执法不受影响」；② 本片零在途窗口（技术上可现在改），但改摘要约定属安全相关函数的最后改动，而 ⑦ 余量为 **0**（无复审余地）；③ 按用户 2026-10-11 的分支裁定，(b) = 登记为已知边界 + 候选观察项，写报告 §8/§9 与 §13。
- **落点**：报告 §9 已知边界（含 **adoption 缓解说明**）+ **候选 C-12**（摘要加固：每列表带标签/计数前缀）⇒ 归后续加固片。
- **新增候选（同一批登记形态）**：**C-13** = 「complete 在 abort/完成后覆写 store-local 工件」（**预存在**，非本批引入：`completeSupervisedRun` 在 `SubmitSupervisedReturn` 之前写回执/报告 ⇒ 对已终态 run 再发起 complete 会在链拒绝前覆写工件；链状态不受影响，仅事后审计对账变弱）⇒ 后续片在 complete 入口先查投影非终态；**C-14** = T027-B 接线预警：守卫冻结的是 policy 内的 **hook ID**，hook **定义**由调用方在归还时提供（`GateHooks`），**定义内容未随窗口冻结**（当前 D-18 fail-closed 下不可达）⇒ 必须写入报告 §9 供 T027-B 输入。Nit-1/Nit-2 已分别采纳/落地，无需另列候选。
- **⑥ / ⑦ 余量台账**：⑥ 已用 **2/3**（首轮格式偏离 + 允许的 1 次重发 = PASS），**余 1 次** ⇒ ②b-8 后对最终态补扫；⑦ 已用 **1/3**（盲审），**余 1 终审 + 1 复审**、**无余量**。
### 9.10 ②b-8 后 ④ 重跑（驱动侧最终态）

- **三作用域门禁**：`--all --scope=tree` ⇒ `changed=28` / `TOTAL_FAIL=0`；`--scope=index` ⇒ `changed=28` / `TOTAL_FAIL=0`；`--scope=ci` ⇒ `changed=2` / `TOTAL_FAIL=0`。**如实标注**：本地 `ci` 作用域读的是**上一提交**（HAIKU-BUMP 文档回写）的 2 路径，**不是本片变更集**；本片无提交 ⇒ 该作用域在本片**无切片相关读数**、CI 双腿不运行（授权边界：commit/push = 否）⇒ ⑨ 必须由本地替代验证承担该缺口，并如实登记「CI 未跑」。
- **门禁自检**：`--selftest` ⇒ `SELFTEST: PASS 335/335`。
- **全包测试**：`go test -count=1 ./...` ⇒ **14 包全 `ok`、0 FAIL**（`examples/whois-shadow` + `tools/agent-probe/enforcement-proxy` + 12 个 `internal/*`）。
- **契约套件**：`npm --prefix tools/contracts test` ⇒ `tests 4 / pass 4 / fail 0`（schema 编译与目录期望一致性、冻结向量复现、AgentRunner 常量对齐、fail-closed 映射）。
- **IDE 诊断**：`internal/chain`、`internal/console`、`internal/adapters` ⇒ **No errors found**。
- **⑦ 输入包重材料化（最终态）**：`head-baseline.txt` `lines=36`（28 路径 + 表头/空行）、`staged-diff.txt` `bytes=331325`（②b-8 前为 327372，增量与 `+14/−9`、`+79/−3` 相称）；`gate-tree.txt` / `gate-index.txt` / `gate-ci.txt` / `gate-selftest.txt` 同步刷新为 ②b-8 后状态。
- **环境事件（登记）**：一次长命令链（`go test` + `cd` + `npm` + 基线脚本串联）被终端误判为 `^C` 中断 ⇒ 已拆分为短命令**全部重取**读数；复核产物齐全（`go-test.txt` 首轮未生成，重跑后生成并读取）。
- **下一步**：⑥ 最终态扫描（第 3 次 = 最后一次额度）⇒ ⑦ 终审（含自由度清单）。

### 9.11 ⑥ 最终态扫描（第 3 次 = 本片最后一次 ⑥ 额度）

- **载体** `independent-reviewer` / `MAI-Code-1.1-Flash (copilot)`；输入 = 最终态材料 10 项（`head-baseline.txt` / `staged-diff.txt` / `gate-{tree,index,ci,selftest}.txt` / `go-test.txt` / `npm-test.txt` / `PLAN-v1.md` / `ORCH-NOTES.md`）+ 8 个实现文件直读；未执行任何命令（只读载体，如实声明）。
- **结论行 = `INDEPENDENT SCAN: PASS`**。结构齐备：A 范围与输入（含无法读取/无法执行清单）；B 发现表明写「**无**」且**未把已知在册项重复列为新缺陷**；C 逐项判定 = 边界 **PASS（有条件：范围近似已登记为边界）**、契约 **PASS**、P1 **PASS**、P2 **PASS**、P3 **PASS**、P4 **PASS（有条件：实现层强制 ≠ 运行时可观测）**、AT-24 七子句 **PASS（有条件，依赖 §8 裁定）**；D 八项必核面逐条「已验证」；E cannot-verify **3 项**（真实 visible 端到端 / 真实外部 agent 的现场观测 / 生产 hook 定义与 gate 装配）。
- **驱动侧引用抽查（机械核对，不构成审查替代）**：其 `file:line` 引用全部落在真实范围内且与符号一致——`internal/console/runtime_supervised.go` 实测 **1879 行**（其引 `L1518-L1622` 在界内）、`newSupervisedScope` 实测 **L1055**（其引 `L1039-L1119`）、`ErrSupervisedPolicyDrift` 实测 **L165** 且漂移执法链实测 **L1615–1653**、`ErrSupervisedScopeUnderivable` 实测 **L183** 且被 **L1057/L1061/L1065** 使用 ⇒ 该 `PASS` **不是空转判断**，引用可信度经机械证实。
- **相对第 2 轮的差异**：第 2 轮的 8 项 cannot-verify 在最终态收窄为 3 项（其余经 ②b-7/②b-8 修复 + 材料化读数变得可核）；**未发现新增缺口**，亦未发现 ②b-7/②b-8 削弱既有守卫。
- **与 ⑦ 盲审的关系（如实登记，防止误读）**：⑥ 的最终 `PASS` **不覆盖** ⑦ 盲审独立报出的 High（归还期冻结 policy 一致性）——该维度在 ⑥ 第 2 轮与最终轮**均未被 ⑥ 自主提出**；本轮 ⑥ 是在**已知该修复**的输入下复核 ⇒ 其「无新发现」**不得**读作「⑥ 独立证实了该修复正确性」，该证实责任在 **⑦ 终审** 与本片 **⑨** 的专项独立实验。
- **下一步**：⑦ 终审（含自由度清单；余 `1 终审 + 1 复审`，无余量）。

### 9.12 ⑦ 终审（`1 终审` 用毕；`1 复审` 未用）

- **载体** `independent-reviewer` / `GPT-5.3-Codex (copilot)`；输入 = 完整最终态材料 + `PLAN-v1.md` + `ORCH-NOTES.md`（§9.1–§9.11）+ **自由度清单** `docs/validation/evidence/T028-freedom-list.md`（盲审轮按设计不附清单，本轮为闭环轮故全量递送）。只读载体，未执行命令（如实声明）。
- **结论行 = `FINAL REVIEW: PASS`**，且**新增可执行缺陷 = 无**；其发现表仅两条 `[在册重述]`（Critical = (b) 裁定重述、Medium = D-18 重述）⇒ **不触发复审**；⑦ 额度终态 = 终审 `1/1` 用毕、复审 `1` **未用**（如实登记「未用」，不结转）。
- **闭环判定**：① 盲审 High（冻结 policy 一致性）= **真闭合**（三锚守卫 + 先守卫后 complete 路径 + 零写负例 + abort 豁免回归）；② 盲审 Low 收敛建议 = **已落为语义承重测试**（词法烟雾项与语义不变量项并存）；③ (b)/D-18 处置与用户裁定一致，并给出报告**条件化表述要点**；④ L-1 已闭、**L-2 选登记 (b) 可接受**（无需改码）；⑤ 四条禁令各有机械实现点 + 回归测试；⑥ 边界/契约/P1–P4 = `PASS`，AT-24 = `PASS（有条件，依赖 (b)/(c)）`；⑦ **F-1…F-18 全部判为「授权内 / 实现选择」，无越界**；⑧ ⑥ 的最终 `PASS` = **有限采信**（其未自主提出冻结 policy 维度一事已被正确登记）；⑨ ②b-6 两个 IDE 修复 = 行为等价（条件：依据台账的 8 场景 A/B 等价探针）；⑩ 报告必含项清单（10 条）已给出。
- **驱动侧引用抽查（机械核对，不构成审查替代）**：`docs/TEST_STRATEGY.md:54` = AT-24 验收行 ✓；`docs/CONTRACTS.md:272` = `### 7.1 supervised-black-box 显式模式` ✓；`schemas/supervised-report.schema.json:133` = `"const": "reduced"` ✓、`:150` = `"const": "cost"` ✓；`internal/console/runtime_supervised_test.go:644` = `TestSupervisedReturnRefusesDriftedPolicy` ✓、`:746` = `TestSupervisedReturnAllowsAbortUnderDriftedPolicy` ✓、`:1376` = `TestNoSupervisedReferencesInAgentRunnerPaths`（词法烟雾）✓、`:1417` = `TestSupervisedPathNeverRoutesThroughAgentRunner`（语义承重）✓ ⇒ **引用逐条命中，非空转判断**。
- **⑦ 要求报告必含（⑧a 直接引用）**：1 AT-24 为**条件闭合**（不得写无条件「全闭」）；2 **非可运行发布**（真实 visible 端到端依赖 T027-B）；3 ⑥ 最终 `PASS` 的能力边界与不覆盖项；4 `scope=ci` 的 `changed=2` **非本片变更集**、CI 双腿未跑；5 L-2 走登记路径 (b) + 后续候选；6 C-10…C-14 在册状态与去向；7 试点指标**未达标**（自捕获 2/15）；8 cannot-verify 清单；9 High 真闭环证据链；10「新增可执行缺陷 = 无；在册重述不计新增」。
- **下一步**：⑧a 报告对（`docs/validation/t028-supervised-black-box{,_EN}.md`，位置镜像、BOM+LF、含试点说明节）。

### 9.13 ⑧a / ⑨ / ⑧b / ⑧c 与 ⑩ 停点

- **⑧a 产出**：`docs/validation/t028-supervised-black-box{,_EN}.md`（各 257 行；位置镜像；BOM + LF）。文件级五形状核验 = `lines/headings/bold/pipes/blanks` = `257/23/307/214/48`，CN 与 `_EN` **逐项相等**。
- **⑨ 原生验证（最终冻结态，`changed=32`）**：`--scope=tree` ⇒ `TOTAL_FAIL=0`；`--scope=index` ⇒ `TOTAL_FAIL=0`；`--selftest` ⇒ `SELFTEST: PASS 335/335`；`npm --prefix tools/contracts test` ⇒ `tests 4 / pass 4 / fail 0`；`go test -count=1 ./...` ⇒ 14 包全 `ok`；IDE 诊断（链 / 控制台 / 适配器 / schemas / contracts / 夹具）⇒ 0 错误；**A12 专项独立实验** ⇒ `A12_CHECKS=58 A12_EXPERIMENT_FAILURES=0`（聚焦复跑 `81 PASS / 0 FAIL`）。
- **⑧b 定稿**：报告 §15 的回写核对行；台账完成行 `docs/t027/REMAINING_SLICES{,_EN}.md`（各 +1 行，**只写报告路径**）。
- **⑧c 收尾**：32 项变更工件 BOM / EOL 逐文件核验 ⇒ `VIOLATIONS=0`；未跟踪残留 **0**；`G7-a` = 2 份报告 OK、`G7-b` = **72 行**产物栅栏全部可解析、`G7-c` / `G7-d` = PASS、`G8-a` / `G8-d` = PASS；`G4b` **由 FAIL 转 PASS**（3 个新 CJK 字 `遭` / `雾` / `蒙` ⇒ **改文、不扩白名单**）。
- **⑩ 停点**：**commit 否 / push 否 / gitee 否**（须同轮授权）。最终调用计数：① 1 / ②a 2 / ②b 8 / ③ 3 / ④ 3 / ⑤ 3（触顶，含 §7.8.3 用户授权例外 1 次）/ ⑥ 3（触顶）/ ⑦ 2（盲审 + 终审；复审额度未用）/ 探针 0。

### 9.14 用户审核收口（PASS WITH FIXES，三处处理）

- **① 台账分类边界注记（结构调整，不改清单结构）**：在 `docs/t027/REMAINING_SLICES{,_EN}.md` 的 T028 完成行**上方**加入注记块（标题为「T027 剩余切片」但自 v1.14 起为 T027 域统一临时台账；T028 属 P1.8 独立任务、不依赖 T027，登记系惯例延续）。为**保持表格结构完整**，注记下方**重列该表的表头与分隔行**，使 T028 行仍处于合法表格内。读数：两文件各 **+6/−0**（1 完成行 + 5 注记/表头行）。
- **② `DEV_PLAN{,_EN}` 的 T028 条目回写（真实缺口，已修）**：`- [ ]` ⇒ `- [x]`；`状态：\`PLANNED / NOT IMPLEMENTED\`` ⇒ `完成：2026-10-11；状态：\`COMPLETED\``（括号内明写「AT-24 在 (b)/(c) 裁定下条件闭合；本片非可运行发布，真实 CLI 端到端依赖 T027-B」）+ 报告路径；REQ 表**未动**（按 PLAN §4-E4）。读数：两文件各 **+1/−1**。**报告初版把 `DEV_PLAN` 判为「不适用」系判定错误**（§11.2 要求切片收尾回写 `DEV_PLAN`；OB-12 例外仅适用于准则治理切片，本片为代码切片）⇒ 已登记 **A-7** 并同步更正报告 §15。
- **③ `tmp/t028/` 归档（部分缺口，路径 B）**：§17 产物栅栏把被受版本控制文档引用的过程件标为 `local-only`，与 §1.5.1 的口径句冲突；归档职责属「tmp/ 生命周期执行」片、不属本片 ⑧c ⇒ 登记候选 **C-15**（与 C-8…C-14 同批），**不另立片**。已登记 **A-8**。
- **报告对同步（除上述外）**：§2（拆为「⑨ 时点 28 项」与「定稿态 34 项」，表末新增三行：报告对 / 台账 / `DEV_PLAN`）；§7.1 说明句；§11（暂存态 34 项、`changed=34`）；§13（新增 C-15）；§17（栅栏新增 `docs/DEV_PLAN{,_EN}.md` 两行）。
- **复核读数（收口后最终态）**：报告对五形状 `lines/headings/bold/pipes/blanks` = `265/22/321/226/48`，**CN 与 `_EN` 逐项相等**；tree / index `changed=34`、`TOTAL_FAIL=0`；`G4b unknown=[]`；`G7-a` 2 份报告 OK、`G7-b` **76** 行产物栅栏全部可解析；`SELFTEST: PASS 335/335`；34 项工件 BOM/EOL 核验 `VIOLATIONS=0`、未跟踪 **0**；VS Code 诊断 0 错误；`git diff --cached --shortstat` = `34 files changed, 7542 insertions(+), 24 deletions(-)`。
- **停点**：**commit 否 / push 否 / gitee 否** —— ②③ 的落盘改动均已**精确暂存**，提交与推送须同轮授权。

### 9.15 提交、CI 双腿与 CI 证据回写（⑩ 终态）

- **主提交** `b095cd2`（`feat(supervised): add the supervised-black-box explicit mode (T028)`）⇒ **34 项 / +7544 −24**；已推送 **origin**（SSH，`0135f85..b095cd2`）；**gitee 未推**；探针 0。
- **CI run `38137422277`**（push 事件、head `b095cd2`）= **`completed / success`**，**首跑即绿**（无 §6.4 处置需求）：
  - 两腿 `Gates` 步骤读数：Ubuntu 腿（`repo=/home/runner/work/prfrail/prfrail`）与 Windows 腿（`repo=D:\a\prfrail\prfrail`）均 `SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=34`、`TOTAL_FAIL=0`。
  - 包级 `ok` 行：Ubuntu 腿 **16**（含 Race 腿 `internal/adapters` 83.702s 与 `internal/chain` 1.703s）、Windows 腿 **14**；两腿 `FAIL` 行 **0**；目标包 `internal/chain` / `internal/console` / `internal/adapters` / `internal/gates` 两腿均现。
- **一次 docs-only CI 证据回写** `049b71f`（`docs(t028): write back the commit id and both CI leg readings (no rule change)`；4 项 / +18 −16）⇒ 已推送 origin（`b095cd2..049b71f`）；其自身 run `38137996587` = **`completed / success`**（**确认绿但不另起记录提交**，按用户授权口径）。
- **⑧c 全量编码口径补录（用户提醒）**：`git ls-files -z` 全类型逐文件 ⇒ `tracked 3431 / text 3342 / binary 89`，**非冻结主树违规 0**，**冻结证据树另计 9 项**（已登记口径）；报告 §11 现**并列**全量与变更集两个口径。
- **⑩ 终态**：工作区洁净（`git status -sb` = `## main...origin/main`，无未跟踪 / 无未暂存）；`git log` 顶两条 = `049b71f` → `b095cd2`。

### 9.16 用户审核（整体通过，一处需修）与 S1_RELEASE_NOTES 第 4 项更正

- **认可项**：台账去号（四号零命中）、报告口径定义（三口径实测）、`049b71f`/`38137996587` 只进报告、六组文档状态同步、门禁读数；**ADR-012 保留 `APPROVED / PARTIALLY IMPLEMENTED` 可接受**（尾句已说明黑箱侧已实现、AgentRunner 未交付）。
- **需修（已改）**：`docs/S1_RELEASE_NOTES{,_EN}.md` 的「发布前剩余门禁」第 4 项，由「完成 T025–T028，并分别通过 AT-22/AT-23/AT-24 验收」改为当前事实：**「完成 T027 并通过 AT-23；T028 已完成、AT-24 条件闭合；S1 仍因 T027/AT-23 阻断」**（EN 同步）。只改这一项。
- **修后读数（16 项暂存，未提交）**：`--scope=tree` = `changed=16` / `TOTAL_FAIL=0`；`--scope=index` = `changed=16` / `TOTAL_FAIL=0`；`--selftest` = `SELFTEST: PASS 335/335`；`G4a` = 16 file(s) OK；`G4b` = `newCjk=319 unknown=[]`；`G3(a)` = 8 对全部 `+N/-N vs _EN +N/-N OK`（`S1_RELEASE_NOTES` 为 `+2/-2` 对 `+2/-2`）；`G7-b` = 76 行栅栏 OK；`go test ./...` = `ok=14 FAIL=0`。
- **编码两口径**：全量 `git ls-files -z` = `tracked 3431 / text 3342 / binary 89`、主树违规 0、冻结证据树 9（已登记）、`MAIN_TREE_CLEAN`；变更集 = 16 项 `VIOLATIONS=0`、未跟踪 0、**全部 `.md`**（`NON_MD=0`）。
- **对称性**：`ASYMMETRY_INTRODUCED=0`（`S1_RELEASE_NOTES` 对形状增量 `[0,0,0,0,0]`）。
- **规模**：`git diff --cached --shortstat` = `16 files changed, 20 insertions(+), 16 deletions(-)`。
- **停点**：**commit 否 / push 否 / gitee 否**，等用户同轮提交授权。

### 9.17 用户审核（PASS WITH FIXES）与 INSTALLATION_PLAN 两处补修

- **认可项**：台账去号、`S1_RELEASE_NOTES` 第 4 项、6 组状态文档、报告 §2 口径定义、`049b71f` 与其 run 只进报告；ADR-012 状态词保留原样可接受。
- **补修（本轮遗漏，已改）**：`docs/INSTALLATION_PLAN{,_EN}.md` 两处，CN/EN 同步，**只改这两处**：
  - **a. §2 目标部署模型末句**：由「两条 AI 执行路径和完整统一 TUI 均未实现」改为「AgentRunner 仍未交付；受监督黑箱候选的产品侧已由 T028 实现（AT-24 条件闭合），但真实 CLI 端到端依赖 T027-B、尚不可运行发布；完整统一 TUI 未实现」（EN 逐句对齐）。
  - **b. §3 待决决策表「VS Code visible 黑箱候选」行末列**：由 `T028 待实现` 改为「T028 产品侧已实现（降级保证）、真实 CLI 端到端依赖 T027-B、尚不可运行发布；AT-24 条件闭合，S1 仍因 T027/AT-23 阻断」（EN 同步）。
- **边界遵守**：未动该文件其它内容、未动 AgentRunner 行的「T027/AT-23 仍阻断」、未动 §1 / §6 / §7 的其它表述。
- **§1.7.2 扫描（未顺手改）**：对该文件全量扫描后，除上述两处外未发现其它 T028 / supervised-black-box 的过时现状表述（其余命中为 T026 候选状态、ZIP 与下载入口 / EOL 待冻结、T016 / T024 / T017 事实、`## 5. Planned User Path` 标题）⇒ **无需停下询问**。
- **修后读数（18 项暂存，未提交）**：`--scope=tree` = `changed=18` / `TOTAL_FAIL=0`；`--scope=index` = `changed=18` / `TOTAL_FAIL=0`；`--selftest` = `SELFTEST: PASS 335/335`；`G4a` = 18 file(s) OK；`G4b` = `newCjk=342 unknown=[]`；`G3(a)` = **9 对全部 OK**（`INSTALLATION_PLAN` `+2/-2` 对 `+2/-2`；报告对 `+2/-0` 对 `+2/-0`）；`G7-b` = 76 行栅栏 OK。
- **编码两口径 / 变更集构成**：全量 = `tracked 3431 / text 3342 / binary 89`、主树违规 0、冻结证据树 9（已登记）、`MAIN_TREE_CLEAN`；变更集 = 18 项 `VIOLATIONS=0`、未跟踪 0、**全部 `.md`**（`NON_MD=0`）；`git diff --cached --shortstat` = `18 files changed, 24 insertions(+), 20 deletions(-)`；`ASYMMETRY_INTRODUCED=0`。
- **停点**：**commit 否 / push 否 / gitee 否**，等用户同轮提交授权。

### 9.18 更正轮提交与 CI 双腿（⑩ 新停点）

- **本轮授权（同轮）**：对已暂存的 18 项**一次提交** + **push origin**（仅 origin、不推 gitee）+ 观察 CI。
- **提交 `ee83b77`**（`docs(t028): sync T028 status to INSTALLATION_PLAN and close the docs-only correction round`）⇒ **18 项 / +24 −20**；已推送 origin（SSH，`049b71f..ee83b77`）；**gitee 未推**；探针 0。
- **CI run `38147035122`**（push 事件、head `ee83b77`）= **`completed / success`**，**首跑即绿**（无 §6.4 处置）：
  - **Windows 腿**（`repo=D:\a\prfrail\prfrail`）：`SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=18`、`TOTAL_FAIL=0`；包级 `ok` = **14**、`FAIL` = 0。
  - **Ubuntu 腿**（`repo=/home/runner/work/prfrail/prfrail`）：`SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=18`、`TOTAL_FAIL=0`；包级 `ok` = **16**（含 Race 腿 `internal/adapters` 84.329s、`internal/chain` 1.915s）、`FAIL` = 0。
  - 两腿 `Gates` 步骤**均实测执行**（原始行取自 `gh run view 38147035122 --log` → `tmp/t028/ci-log3.txt` → `ci-summary.js`）。
- **回写决策（按授权口径）**：本次 docs-only 更正**即本片唯一提交**，**不再为其加记录提交**；报告 §9 的提交/CI 段落仍以主提交 `b095cd2` 与其 run `38137422277` 为准，`049b71f` / `38137996587` 亦在册；`ee83b77` / `38147035122` **仅记于本编排记录与本轮回贴**，如后续需写入报告，须**同轮授权**再开一次 docs-only 补全提交。
- **⑩ 停点**：工作区洁净（`git status -sb` = `## main...origin/main`），HEAD = `ee83b77` = origin/main；**未对 gitee 执行任何操作**。

### 9.19 用户裁定：C.0p 状态行修正（docs-only，待授权提交）

- **范围**：只改 `docs/DELIVERY_DIRECTIVE.md` 附录 C.0p 的状态行（+ `_EN` 镜像同步，系 G3(a) 同构要求）。未扩范围：未触 §0.3 / §7.2.3 / §12.4 / §2.2 / §2.4 / §7.5 / §6.3 / CONTRACTS / schema / 代码 / workflows / `go.mod` / 角色集。
- **改动**（各 `+1/−1`）：
  - CN（`docs/DELIVERY_DIRECTIVE.md:1459`）：`| 状态 | ⬜ **进行中**（本片 ⑩ 前的如实状态：待 ⑤⑥⑦ 与 ⑨ 完成后由 ⑧b 更新） |` → `| 状态 | ✅ **已完成**（报告见 docs/validation/directive-haiku-bump{,_EN}.md；⑩ 前状态曾为进行中，已由 ⑧b 更新） |`
  - EN（`docs/DELIVERY_DIRECTIVE_EN.md:1459`）：`| Status | ⬜ **In progress** (the truthful state before ⑩ ...) |` → `| Status | ✅ **Complete** (report at docs/validation/directive-haiku-bump{,_EN}.md; the pre-⑩ state was in progress and ⑧b updated it) |`
- **与报告一致性**：`docs/validation/directive-haiku-bump{,_EN}.md` 的结论为已定稿（⑦ 终审 + 用户授权收窄读 PASS、⑨ 本地与 CI 全绿、提交 `ccacd26` 与 docs 回写 `0135f85` 均已推送）⇒ 与「已完成」一致，**无分歧、无需行内另行注明**。
- **口径不变（未改）**：T027 仍 `BLOCKED / NOT IMPLEMENTED`；AT-23 未通过；S1 仍因 T027/AT-23 阻断；T028 非可运行发布；AT-24 条件闭合。
- **不新增记录提交**：`ee83b77` / `38147035122` **不另开**「记录其自身 run」的 docs-only 提交；报告 §9 保持主提交 `b095cd2`/`38137422277` 与 `049b71f`/`38137996587`。若后续需把 `ee83b77` 写入报告，则并入下一片起跑前的 docs-only 补全，且不为那次补全再新增提交。
- **修后读数（2 项暂存，未提交）**：`--scope=tree` = `changed=2` / `TOTAL_FAIL=0`（`G1-b` = `header=v1.27 last=v1.27 2026-10-11/2026-10-11 rows=28 OK`；`G3(a)` = `+1/-1 vs _EN +1/-1 OK`；`G4a` = 2 file(s) OK；`G4b` = `newCjk=17 unknown=[]`；`G6-1` / `G6-4` PASS）；`--scope=index` = `changed=2` / `TOTAL_FAIL=0`；`--selftest` = `SELFTEST: PASS 335/335`。
- **编码两口径 / 构成**：全量 = `tracked 3431 / text 3342 / binary 89`、主树违规 0、冻结证据树 9（已登记）、`MAIN_TREE_CLEAN`；变更集 = 2 项 `VIOLATIONS=0`、未跟踪 0、全部 `.md`；`shortstat` = `2 files changed, 2 insertions(+), 2 deletions(-)`；`ASYMMETRY_INTRODUCED=0`（C.0p 对增量 `[0,0,0,0,0]`）。
- **停点**：⑩ —— 待用户同轮授权后提交与推送（仅 origin、观察 CI 双腿、不推 gitee）。

### 9.20 C.0p 定稿提交与 CI 双腿（⑩ 停点）

- **授权（同轮）**：`git commit` + `git push origin main`（仅 origin、不推 gitee）。
- **提交 `6cb7dce`**（`docs: finalize C.0p status for DIRECTIVE-HAIKU-BUMP`）⇒ **2 项 / +2 −2**（`docs/DELIVERY_DIRECTIVE.md`、`docs/DELIVERY_DIRECTIVE_EN.md`）；推送 **origin（SSH）**：`ee83b77..6cb7dce`；**gitee 零操作**；探针 0。
- **CI run `38149827513`**（push 事件、head `6cb7dce`）= **`completed / success`**，**首跑即绿**（无 §6.4 处置，亦无以复跑掩盖首跑红的情形）：
  - **Ubuntu 腿**（`repo=/home/runner/work/prfrail/prfrail`）：`SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=2`、`TOTAL_FAIL=0`；包级 `ok` = **16**（含 Race 腿 `internal/adapters` 84.219s、`internal/chain` 1.951s）；`FAIL` = 0。
  - **Windows 腿**（`repo=D:\a\prfrail\prfrail`）：`SELFTEST: PASS 335/335`、`GATE REPORT scope=ci changed=2`、`TOTAL_FAIL=0`；包级 `ok` = **14**；`FAIL` = 0。
  - 两腿 `Gates` 步骤**均实测执行**（原始行：`gh run view 38149827513 --log` → `tmp/t028/ci-log4.txt` → `ci-summary.js`）。
- **不新增记录提交**：按 §11.2 与用户口径，**不为该 run 再开 docs-only 记录提交**；若日后需把本次 run 写入报告，并入下一片起跑前的 docs-only 补全，且不为那次补全新增提交。
- **⑩ 停点**：工作区洁净（`git status -sb` = `## main...origin/main`）；HEAD = `6cb7dce` = origin/main；文档集自本次起可作为 **T028 文档基线**。

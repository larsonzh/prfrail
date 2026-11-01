# 任务执行清单（跨会话快照 v35 · 装配验证收官版）

生成时点：2026-11-01，HEAD = origin/main = `3a7541d`，准则 v1.60。
替代关系：完全取代 v34 及任何早于本日的快照。
本版修订（相对 v34）：① 锚点更新至 P45 全链闭环（`3a7541d`）；② §二 同时收录 **DR13-FIX 前置修复片**（不改编号）与 P45 事实包两段；③ 队列重排为「耐久升级（1）/ 触发式（2）/ 低优先级（3）/ P45 后续片候选（4）」；④ 新增三条装配验证类教训（零模型机械判据集合、并发断言的存活与释放面、console 单任务限制的读数边界）；⑤ 继承 v34 的 24 项非终态观察项与 P39–P45 全部后续片候选。

---

## 一、当前锚点（开场先核这 5 个数）

| 项 | 值 | 核对命令 |
|---|---|---|
| HEAD | `3a7541d` | `git rev-parse --short HEAD` |
| origin/main | `3a7541d` | `git rev-parse --short origin/main` |
| 准则版本 | **v1.60** | `docs/DELIVERY_DIRECTIVE.md` 头部 |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest`（见下方「命令坑」） |
| 三作用域 | **全 TOTAL_FAIL=0** | `gate.js --all --scope=index\|tree\|ci` |

**若 ≥2 项不符 ⇒ 以仓库实况为准**，先让主控报终态读数再继续。

**`--selftest` 命令坑**：全量约需数分钟，`| Select-Object -Last 3` 会长时间无输出、易误判挂起。建议用 `| Select-String -Pattern 'SELFTEST'` 并容忍长时间静默；不要中途终止。

**【三作用域全绿口径】** 干净树时 `index/tree` 均 `scanned 0`（vacuous pass），只有 `ci` 扫到 N 行——须按纪律 E 报 `scanned N` 实际值。v35 生成时干净树实测 `ci` 的 `G6-1 scanned 4`（`--scope=ci` 以 `HEAD` 端点比对，非零）。

**【tmp/ 固定断言】** `tmp` 顶层 = `.gitkeep,gates-cache`；`gates-cache` 子目录数 = 1（`g4b`）；**`tmp/gates-cache/` 下 `g4b/` 之外的文件数 = 0**；`Test-Path tmp/P40-T027B-LOCK-REAPER`、`tmp/P41-T027B-SPAWN-OWNERSHIP`、`tmp/P42-T027B-EVIDENCE-BIND`、`tmp/P43-T027B-STATE-SEAL`、`tmp/P44-T027B-GATE-READINESS`、`tmp/P45-T027B-E2E-ZERO-MODEL`、`tmp/DR13-FIX` 均为 `False`。

> **口径区分**：「`gates-cache` 下 `g4b` 之外 = 0」与「`tmp/` 排除 `gates-cache` 后 = 1（`.gitkeep`）」是两个不同读数，不可互相替代。

---

## 二、P45 与前置修复片事实包

### 2.1 DR13-FIX（P45 前置修复片，不改编号）

**定位**：P45 的前置修复片，不占 v34/v35 队列 1–4 的编号；其价值是把 P45 需要的并发隐患先清掉。

| 项 | 值 |
|---|---|
| 提交链 | 主 `b123db1`（15 files，+1057/−2）＋ 回填 `eb36992`（5 files，+33/−4） |
| CI | `41269372878` **首跑双腿全绿**（ubuntu `ok adapters 12.343s`＋Race `ok adapters 89.396s`；windows `ok adapters 34.937s`）；`41270106374`（docs 回填）双绿 |
| 定性 | **`DIAGNOSIS: B`（真实排他窗口）**：`RecoverStarting` 把「starting marker 的 flock 可获取」当作「前宿主已死」，同进程竞争者锁住已 unlink 的孤儿 inode 后，unix 分支 `reclaimStartingProcess(alive=true)` 杀掉活进程并夺走 slot ⇒ 两个 Start 先后都成功 |
| 修法 | **D1** 同进程所有权门（`registry.lookup(marker.LaunchID)` 命中即拒绝回收）；**D2** 锁后 `os.Stat(path)` 复核孤儿 inode；**D3 不做**（先把路径 remove 再 unlock，留后续片） |
| 反例 | `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns`：未修复代码确定性红（Windows 失败点 `spawn_recovery_test.go:281`），修复后绿 |
| 审查链 | ⑤ r1 `PASS WITH FIXES`（2 Medium）→ r2 `PASS`；⑥ `INDEPENDENT SCAN: PASS`（§B 零缺陷）；⑦ `FINAL REVIEW: PASS`（6 条放行条件全部落地） |
| 台账 | `DR-13` 处置列已翻「已修复」 |
| 遗留 | D3；跨进程 W1 残留；`StopAgentRunnerProcess` 失败分支无条件 `releaseLaunch`；W2 内核 flock 语义未实证；unix `StopProcessIdentity` 未本机复现；跨 launcher/跨进程镜像失效窗口 |

### 2.2 P45-T027B-E2E-ZERO-MODEL（零模型三任务 E2E）

| 项 | 值 |
|---|---|
| 提交链 | 主 **`4273d6e`**（23 files，+1599/−3）＋ 回填 **`3a7541d`**（5 files，+32/−4） |
| CI | `41315652405` **首跑双腿全绿**（ubuntu `ok chain 0.564s` / `ok console 2.856s` / `ok adapters 12.291s` / `ok gates 0.159s`＋Race `ok adapters 89.514s` / `ok chain 2.521s`；windows `ok chain 5.011s` / `ok console 20.796s` / `ok adapters 34.297s` / `ok gates 0.786s`；两腿 `SELFTEST 697/697`、`changed=23`、`TOTAL_FAIL=0`）；`41316975346`（docs 回填）双绿 |
| 审查链 | ① `DESIGN: DONE` → ③ 首轮 **7 绿 1 红**（红点= `examples/no-ai-preview/proofrail.chain.json` 缺失，**真实缺陷**）→ ② 补夹具 + README 命令路径最小修正 → 8 绿 → ④ 门禁 `TOTAL_FAIL=0` → ⑤ `PASS WITH FIXES`（2 MINOR）→ ⑥ `INDEPENDENT SCAN: PASS`（§B 零缺陷）→ ⑦ `FINAL REVIEW: PASS`（1 MINOR 已修正，6 条放行条件落地） |
| 交付 | 新增测试三文件（chain 引擎腿 4 例 / console 旁证腿 3 例 / console 并发腿 1 例）；`examples/no-ai-preview/proofrail.chain.json`；`examples/no-ai-preview/README.md` 单行修正；报告对 + 证据包 + 台账 + `ARCHITECTURE` / `DEV_PLAN` 双语回写 |
| **生产代码** | **零改动**（`git diff --name-only -- internal/` 为空；三个新文件全为 `_test.go`） |
| 四路互证（T3 未启动） | workspace 调用序 / 事件流无 `task-three` / 投影键无该任务 / 步骤请求集无该任务 |
| 零模型机械判据 | `StepIntents==nil` + `ExecutionTarget==DefaultExecution` + `AgentRunnerFacts==nil` + 无 `agentRunnerDispatchedReason` 事件；预览腿 `RunnableInCLI==true` + `Outcome==ready` + `BlockingEvidence` 空 + 全部 step `kind==noop` |
| 并发断言 | 同 launcher/同 runRoot、8 路 barrier；恰一赢家；败者全命中 `ErrAgentRunnerLaunchConflict`/`ErrAgentRunnerProcessDuplicate`；赢家 `ProcessAlive`；活体赢家拒绝迟到启动；`Stop` 后槽位释放且新 Start 成功 |
| 边界声明 | **console 装配层仍强制单任务**（`runtime_agent.go:774`、`runtime_supervised.go:335`），三任务 E2E 落在 **chain 引擎层** + console 两条旁证腿；**不得读作 console 已支持三任务** |
| 偏差 | **D10** `tools/gates/cjk-newwords.txt` +2 条（`营`/`尊`，仅因 ⑥ 正文逐字转录不可改写；白名单 56 → 58）；**D11** 起跑消息写「16 包」，实测 15 个测试包 + 4 个无测试文件包 |
| `00-briefing.md` | SHA256 `562ce6e885f45342258fadc7cec508bd975678559184dcdddcf182b3132797fa`（入 MANIFEST + SHA256SUMS + artifacts 围栏） |

### 2.3 P45 新增教训（摘要，全文并入 §八）

① 零模型证明必须是**机械判据集合**，不是「没看到调用」；② 并发断言必须含**赢家存活**与**槽位释放可重发**，不能只测「失败者是否报错」；③ **console 装配层单任务限制是已知边界**，E2E 的三任务贯通只在 chain 层，不得扩大读法。

---

## 三、已完成切片（近期基线）

| 版本 | 切片 | 主提交 + 回填 | 结论 |
|---|---|---|---|
| v1.56 | P34 | `acf676e` + `cc703ac` + `ebec7d1` | 闭环（R=3、D=2） |
| v1.56 | P35 | `397dc30` + `70cb193` | 闭环（R=1、D=0） |
| v1.57 | P36 | `9dca0e7` + `e691b37` + `00ddc53` | 闭环（R=1、D=0） |
| v1.58 | P37 | `1bf57a4` + `00a4627` + `fa0fa3a` + `344c206` | 闭环（R=4、D=0） |
| v1.59 | P38 | `c59246b` + `5de5657` + `9ea54dd` + `8b16e21` | 闭环（R=8、D=0） |
| v1.60 | P39 | `4d5ebf3` + `35bede7` + `4994fa8` | 闭环（R=4、D=0） |
| v1.60 | P40 | `3fab6ba` + `3af9c9a` + `a142eaa` | 闭环（R=0、D=0） |
| v1.60 | P41 | `8799015` + `caaf64d` | 闭环（R=0、D=0） |
| v1.60 | P42 | `cab3e08` + `be612fe` + `d3a511b` | 闭环（R=0、D=0） |
| v1.60 | P43 | `b6e0f34` + `609bb19` + `7004282` + `8a46b01` | 闭环（R=0、D=0；D7/D7b/D8） |
| v1.60 | P44 | `9061b67` + `16abe2b` | 闭环（R=0、D=0；D1–D6，新增 DR-13） |
| v1.60 | **DR13-FIX** | `b123db1` + `eb36992` | 闭环（P45 前置修复片；DR-13 已修复） |
| **v1.60** | **P45** | **`4273d6e` + `3a7541d`** | **闭环（R=0、D=0；D10/D11）** |
| v1.60 | v34 快照 | `8e834b2` | 基线（「P44 闭环 → P45 起跑前」） |

---

## 四、执行队列

**总原则：预算 → 锁序 → 归属 → 证据 → 封存 → 门控就绪 → 末装配（均已完）→ 耐久升级。**

### 队列 0 · P46 战略回看（当前插片，草案待起草）

| 字段 | 内容 |
|---|---|
| 标识 | `STRATEGIC-REVIEW-2026-11-02`（日期可按实际起跑日调整） |
| 内含 | 回看 P39–P45 闭环后的项目状态，回答 P35 双断点是否闭合、下一步最优路径、11 月优先级是否需调整 |
| 依赖 | P45 闭环（`3a7541d`） |
| 规模 | L（评估类） |
| ⑦ 额度 | 不适用（流水线裁剪） |
| 边界 | 不改代码、不改准则正文；报告对 + 证据包；切片草案由主控起草、用户审核后起跑 |

### 队列 1 · `CONTRACTS-S7-DURABILITY`（顺延）

| 字段 | 内容 |
|---|---|
| 标识 | `CONTRACTS-S7-DURABILITY` |
| 内含 | Windows 发布耐久升级：`CONTRACTS.md` §7 修订 + 独立评审 |
| 依赖 | 依赖 P39（DurableCostLedger）与 P44（GateReadiness §4.2 就绪分级）已闭环；**且依赖 P46 战略回看的裁定结果**（若 P46 认为应调整优先级，则本队列随之调整） |
| 规模 | M |
| ⑦ 额度 | 1+1 |
| 边界 | 仅 `docs/CONTRACTS{,_EN}.md` 与 schema 校验；**先协议后代码** |

### 队列 2 · `G4B-CORPUS-DECOUPLE`（触发式，当前挂起）

| 字段 | 内容 |
|---|---|
| 标识 | `G4B-CORPUS-DECOUPLE`（OB-120） |
| 内含 | `listAllMarkdown` 的 tree/index 用 `ls-files`，暂存删除 `.md` 时语料收缩 ⇒ G4b 冷跑从严多报 |
| 依赖 | 触发条件成立：出现「暂存删除 `.md` 后门禁假阳性」的实际误判 |
| 规模 | M |
| ⑦ 额度 | 不适用（触发式） |
| 边界 | 候选修法见 OB-120 现象列（A / B / C） |

### 队列 3 · 低优先级（长期挂起）

| 字段 | 内容 |
|---|---|
| 标识 | `OB-50` / `spawnGit` / `ci 冻结成本` / `%TEMP% 历史遗留 22 件` / `GATES-EXEMPT` |
| 内含 | 逐项见 §12.4 对应 OB 行；`GATES-EXEMPT` = P33 的 ①-3(iv) 拆片，未立项，触 4 个判据文件 |
| 依赖 | 长期挂起，无明确触发 |
| 规模 | 不等（含 S 到 M-L） |
| ⑦ 额度 | 无 |
| 边界 | 逐项见 §12.4 |

### 队列 4 · P45 后续片候选（不并入队列 0–3）

| 候选 | 内含 | 触发 |
|---|---|---|
| **console 多任务装配扩展** | 放开 `runtime_agent.go:774` / `runtime_supervised.go:335` 的单任务限制，需先解决「每 runRoot 唯一 identity mirror」与多任务并发的冲突；**P46 战略回看将评估其优先级** | 产品确需 console 层三任务 |
| **Race 腿扩包** | 把 `./internal/console/...`（及 `internal/gates`）纳入 CI Race 腿，属 `ci.yml` 变更；P45 并发腿目前无竞态检测证据 | CI 预算允许 |
| **`registry.len()` 与 mirror 直接断言** | P45 以包外行为等价替代；若要直接断言，需在 adapters 侧新增导出只读查询面 | 另立片 + 协议先行 |
| **⑥ 的 4 项未覆盖面** | UI 端到端、跨 OS/locale 差异、证据包终态核对、真实 `preview` 命令的 CI 化 | 后续验证片 |
| **D3 / 跨进程 W1 / Stop 失败分支** | DR13-FIX 遗留的并发加固面 | 并发加固片 |

> 队列 1–2 为已排期/触发式项；队列 3–4 不占当期排期，登记即可。

### 后续片候选（P39–P45，不并入队列 1–4）

**P39 ⑤ 四项非阻塞发现**：持锁进程崩溃后的锁释放语义无测试腿；`rename` 成功后报错的可见分支无测试腿；`newDurableAck` 的 `confirmed` 字段命名与语义不符；矩阵「确认前取消」腿未落地。

**P40 报告登记的后续片候选**：F-07a 静态锁序断言入 CI 判据；F-11 `lockPath` canonical 化；Reclaim 生产接线（uncertain 路径）；CI Race leg 扩包覆盖 `internal/guard` 与 `internal/tickets`。

**P41 报告登记的后续片候选**：**I1** kill-on-close 半证；**I2** 四缝注入（Start 后 / attach 前 / OnSpawn 内 / resume 前）；⑦ F-01 拒绝态三态化；⑦ F-02 rename-away 对象绑定删除；⑦ F-03 mirror 删除断言；Pdeathsig 独立切片；⑤ Low-1 槽位新鲜度兜底评估。

**P42 报告登记的后续片候选**：**N2 读路径 sweep**（`hashLogs` / `readPIDs` / `fileHash` / `scanWorkspace`）；Unix 两段式 dirfd；Windows reparse 对象身份级校验；`CloseBoundDirs` 生产接线；darwin 旧文档 supersede；原语层测试补齐（关闭后拒读 / 超界拒摘要 / 身份同源对抗）。

**P43 报告登记的后续片候选**：生产装配接线 `Attempts`（console 侧 `TaskAttemptResolver` 来源）；supervised attempt>1 端到端；seal 目录级耐久 crash 注入；stale/mismatch 判定矩阵补测；`router.go` 与 handoff 的 K 绑定；`WorkspaceReopener` 落地；`StepRequest.ContextKey` 的读取者落地。

**P44 报告登记的后续片候选**：console 接线（R-3）；Race 腿扩展（R-4）；类型层绑定（R-2，`ProductionReadiness` 构造签名接收 L2 值）；⑥ 的 4 处未覆盖面（`LoadProbePermit` 非 `NotExist` 读失败 / `Require` 的 `ctx == nil` / `Complete` 未知态 / 真跨进程双写）。

**P45 报告登记的后续片候选**：见「队列 4」四行，另加「console 单任务限制解除需协议先行」。

**低优先级（长期挂起）**：`OB-50` / `spawnGit` / `ci 冻结成本` / `%TEMP% 历史遗留` / `GATES-EXEMPT`。

---

## 五、观察项

### 观察中（12）
`OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113`

### 待议（11）
`OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`

### 待办（1）
`OB-69`

**非终态合计 24 项**（观察中 12 + 待议 11 + 待办 1）；本轮未改 §12.4，读数与 v34 一致；其余按 §12.4 检索。

### 已裁决（28，含既有窗口新增）
OB-117–OB-124（顺延第 16–22 次）、OB-125（ADR-028）、OB-126/127/128（第 23/24/25 次）、OB-129–OB-131（Astra 三次授权）、OB-132（① 位角色-模型绑定外授权，P39）。P40–P45 窗口的用户裁定（方案 A、F 项全盘、DR-12、各片起跑授权、主控模型策略、P43 疑点裁定、P44 起跑与闭环授权、DR13-FIX 选 A、P45 起跑与闭环授权）属切片级裁定，未入 OB。

### 机械读数（可核对口径）
- §12.4 最大 OB = **OB-132**；最大顺延 = 第 25 次；ADR max = **28**；DR max = **13**（已修复，台账实测）。
- §12.4 条目数 132；`已处置 75 / 已裁决 28 / 已登记 5 / 观察中 12 / 待议 11 / 待办 1`。

---

## 六、待用户裁定

| # | 项 | 现状 | 阻塞谁 |
|---|---|---|---|
| 0 | **P46 战略回看是否立项**（草案待审核） | 待裁（v35 落盘后由主控起草 P46 草案） | 队列 0 |
| 1 | **`CONTRACTS-S7-DURABILITY` 立项节奏**（队列 1） | 待裁 | C6 的 Windows 落地 |
| 2 | P45 后续片的立项归属（console 多任务装配 / Race 腿扩包 / `registry.len()` 直接断言） | 待裁 | 后续装配与验证片 |
| 3 | DR-13 是否继续保持「已修复 + 观察」口径（若同一断言串再现即重开） | 待裁 | CI 稳定度 |
| 4 | DR13-FIX 的 D3（先 remove 再 unlock）是否升格为独立加固片 | 待裁 | 跨进程并发加固 |
| 5 | OB-55 归属与 T027-B 立项节奏 | 待裁 | 后续实施 |
| 6 | MVA、S1 移出项与 S2 承诺收窄批准 | 待裁 | 后续切片 |
| 7 | 候选与 Windows 前提 P 及资源/探针授权 | 待裁 | G2 门 |
| 8 | AT-24 `unknown` 与「未知副作用保持暂停」接受边界 | 待裁 | AT-24 收口 |

---

## 七、起跑与步骤 0 通则（代码切片）

1. 核锚点 → 步骤 0：登记本片台账（`R/D = 0` 基线、临时件落点 `tmp/<切片ID>/`）、预取输入包并落 `input-list.md` 与 `freedom-list.md`。
2. 每个修复切片必须给出**回归测试**（对 P38 报告列出的测试名逐条落地或说明替代），并把「未执行测试/lint」显式登记。
3. ⑨ 三作用域全绿 + `SELFTEST 697/697` + `go build/vet/test ./...` 全绿，方可提交。
4. 模型绑定（P39 教训，硬约束）：① `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`；⑤ 同；②③⑧ `DeepSeek V4.1 Flash (deepseek)`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`。**每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`，严禁擅自升级或替换**。
5. **执行节奏纪律（用户 2026-10-31 裁定，硬规则）**：**常规不停**——④ 门禁全绿 ⇒ 直接派 ⑤；⑤/⑥/⑦ 报 Medium+ ⇒ 修复后**原位重跑该环**；**⑩ 才是唯一常规停点**。仅遇架构冲突、超预算或无解时停下升级。
6. **跨版本 / 跨平台排查纪律（DR-12 + P42）**：CI 与本机行为分歧时优先核对工具链版本与**构建标签覆盖面**；交叉构建后必须 `go env GOOS GOARCH` 复核；删除「本机未使用」符号前必须做至少一次对向 `GOOS` 的 `go vet`。
7. 批量辅助：每片实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` **≥ 2 次**（策略见 §八「GLM 批量辅助策略」）。
8. 主控模型策略：默认 `DeepSeek V4.1 Flash`；遇底层并发 / 状态机 / 句柄类复杂片、单份超 10,000 字符设计文档、大量 CI 日志深度根因分析、或 Flash 连续 2 次行为漂移时提示切换 Pro；Max 不作默认。**任何降级须按 OB-37 元规则独立登记**。
9. **角色载体不可替代（P43 D7 起，硬规则）**：② 必须由 `prfrail-implementer`、③ 必须由 `prfrail-tester` 派单执行；主控只做编排、落盘、门禁与回执，**不得代写产品代码或测试**。违反即登记为最高级独立性降级（D 项 + 独立影响评估）。
10. **起跑消息即工件（用户 2026-11-01 裁定，硬规则）**：步骤 0 建目录时，必须将用户的起跑消息原文（剥离附件清单后）落盘为 `docs/validation/evidence/<切片ID>/00-briefing.md`，并加入证据包的 `MANIFEST.md` 与 `SHA256SUMS.txt`，同时在报告的 `artifacts` 围栏中声明。缺此件不得进入 ①。
11. **线性化点纪律（P44 ⑦ FR-01 起，硬规则）**：凡涉及持久化状态机，**「持久化成功」必须是线性化点**；内存 CAS 只能作同进程快速拒绝缓存。注释与契约须写明「取额成功 ⟺ 标记已在盘」并配正向红测。
12. **证据包与围栏同步纪律（P44 ⑦ r2 NR-01 + G5(d)）**：新增任何证据载荷必须**同轮**完成三件事——补 `MANIFEST.md` 归档条目；重算 `SHA256SUMS.txt`；在报告 `artifacts` 围栏补行。任缺其一都会被 `G8-b` 或 `G5(d)` 拦下。
13. **审查方归档纪律（P44 起）**：⑤⑥⑦ 的产出归档必须**逐字原文**；主控不得按要点重建；若确需摘要，摘要另立文件并在围栏外披露。
14. **装配验证片纪律（P45 起，硬规则）**：断言必须落在**可证伪的机械判据**上；不得把「常量恒等」写成哨兵（须显式声明其非哨兵属性）；报告必须声明「生产代码是否零改动」，且「某层未支持的能力」不得被读作「已支持」。「层边界」只适用于**已明确记载的已知限制**；新增未记载的能力缺口不得据此免除登记义务，必须登记为后续片候选。

---

## 八、跨片协作纪律（起跑消息必须落实）

| 要素 | 内容 |
|---|---|
| 载体表 | ① ⑤ ⑥ ⑦ 用 `deep-reasoner` / `independent-reviewer`（只读）；② ③ ⑧ 用 `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`（可写）；禁止 `prfrail-tester` 跑任何审查层环节；批量辅助载体必须是 `deep-reasoner`，禁止 `quick-researcher` |
| 派单表机械读取 | ⑩ 派单表每行 `agentName` 必须从 `transcript.arguments.agentName` 逐字机械读取，附脚本 stdout 原文 |
| R2.1 补 agentName | 每次 `runSubagent` 的 `arguments` 必须同时含 `agentName` 与 `model`，两者都不得缺 |
| 提交信息含非 ASCII 须用 -F 文件 | 提交信息含任何非 ASCII 字符时须用 `git commit -F <UTF-8 文件>`；纯 ASCII 可用 `-m` |
| 「可独立核对 vs 不可独立核对」判别标准 | 可独立核对的产出（枚举、事实比对、引用点清单），模型错误不构成作废理由；不可独立核对的产出（判断、结论、审查意见），模型错误必须重做 |
| ⑥ 先于 ⑦ | ⑥ 必须在 ⑦ 之前完成；重派 ⑥ 时若文档有变更，⑦ 结论须重新评估 |
| 临时件位置 | 一律 `tmp/<切片ID>/`；逐件精确清理 + 回读证明 |
| 临时目录只增不删 | 执行角色在 `tmp/<切片ID>/` 内只增不删；只能由 ⑧c 清理；清理时**只删本片目录**，严禁递归删父目录（P44 D5 实证） |
| 长命令重定向 | 不要用 PowerShell `>`；用 `\| Tee-Object -FilePath <path>` 或 `cmd /c "..."` |
| 门禁执行 | 中间环节用 `--all --skip-base`；提交前必须全绿（DR-11 红线） |
| 纪律 E · index 作用域的正确顺序 | 先精确 `git add` → 再跑 `--scope=index`；回执必须同时报 `G6-1 scanned N` 且 N > 0 |
| 纪律 F · 规范优先的适用边界 | 起跑消息的具体示例与措辞规范冲突时按规范执行，但须回执登记冲突 |
| 批量辅助 | 只用 `GLM-5.3-Flash (glm)`；**≥ 2 次**，按 §八「GLM 批量辅助策略」尽量多派 |
| 审查方只读纪律 | ⑥ ⑦ 不得创建、修改、删除任何文件（含 tmp/）；数据来源只能写「读取主控落盘的 <路径>」 |
| 纪律 A · 推定授权 | ⑩ 停点回执的 CI 取证需要 push，视为该片授权范围内的必要动作；用户同轮明确「不 push」则以用户措辞为准，不得作先例 |
| 纪律 B · 证据缺口表述 | 证据缺口必须如实登记「丢失了什么、什么未丢」，并明写「未主张证据缺口已消除」 |
| 纪律 C · 枚举优先于推导 | 机械读数优先于算术反推；遇「总量 − 已知 = 差值」时多跑一条枚举命令 |
| 纪律 D · 白名单政策优先于授权数量 | 授权数量与政策冲突时以政策为准，但执行时须预置「数量可能超出预估」的说明 |
| 纪律 G · 围栏即落地快照 | `artifacts` 围栏按路径排序、TAB 分隔、CN/EN 逐字节相同；每次触及新文件必须同步补行，否则 G5(d) 报 staged-not-declared |
| 纪律 H · 提交前对向 GOOS 复核 | 删除或改动常量/函数前至少跑一次对向 `GOOS` 的 `go vet`；改动平台文件后跑三平台 `build + vet` |
| 纪律 I · 判定词闭集 | ⑥ ⑦ 的判定行必须落在 §4.1 闭集内；围栏归一须在围栏外披露原文与理由 |
| 纪律 J · 审查语义 | ⑥ 的 Section B 只允许真实缺陷；「已接线 / 已覆盖」归 Section C；严重度不得用于「已实现」陈述 |
| 纪律 K · 持久化优先 | 见 §七 第 11 条：持久化成功即线性化点，内存 CAS 不作准入门 |
| 纪律 L · 三方自洽与围栏同步 | 见 §七 第 12 条：载荷集 = MANIFEST = SHA256SUMS，且同轮补围栏 |
| 纪律 M · 归档即原文 | 见 §七 第 13 条：审查方产出逐字归档，不得按要点重建 |
| 纪律 N · 机械判据与非哨兵声明 | 见 §七 第 14 条：断言必须可证伪；常量恒等须声明非哨兵；须声明生产代码是否零改动与「未支持能力」的读数边界 |

**GLM 批量辅助策略（用户 2026-11-01 裁定，能力近乎无限）**：`GLM-5.3-Flash (glm)` 资源充足，**「≥ 2 次」仅是下限**；应尽可能把机械核对、日志分析、枚举对比外包给它。允许范围：枚举/统计、长文档摘录与数字核对、双语镜像检查、MANIFEST 草案、日志机械提取、路径与文件操作普查、测试断言点枚举。禁止：复杂逻辑判断、排序与 Top-N 择优、审查结论、裁定决策。每次实派次数必须计入 ⑩ 的成本与计量。

**P38 新增教训（v28 起）**：单个子代理单次输出有硬上限；长文档必须分批（每批 ≤ 10,000 字符、单主题），派单写明「禁止自重启」与 `SELF-RESTART-ABORT`。

**P40 新增教训（v30 起）**：① 交叉构建后 `go env GOOS GOARCH` 复核；② CI 与本机分歧先查工具链版本（DR-12）；③ ⑥ Section B 的 `INFO` 行不是 severity 行。

**P41 新增教训（v31 起）**：① Section B 写 `Low | …` 管道形态会触发 G7-f 的 severity 行闭集；② Windows 持锁删除需 `FILE_SHARE_DELETE`；③ Job 属性缺 `EXTENDED_STARTUPINFO_PRESENT` 被静默忽略；④ 自写 `CreateProcessW` 必须配 `PROC_THREAD_ATTRIBUTE_HANDLE_LIST`。

**P42 新增教训（v32 起）**：① 「本机未使用」≠「全仓未使用」，删符号前须对向 `GOOS` 复核；② `//go:build !windows` 覆盖 darwin，不可用平台配 fail-closed 桩；③ gopls 的 `unusedfunc` / `unusedparams` 会针对构建标签变体报；④ 单腿红灯先怀疑构建标签排除；⑤ 围栏转录必须逐字；⑥ 批量辅助是强制项。

**P43 新增教训（v33 起）**：① 主控不得代写 ②③；② ⑥ 的 Section B 语义必须收紧（纪律 J）；③ 额度按轮次簿记；④ 偏差须附独立影响评估；⑤ 「声称已改」必须用机械判据复核，否则会重演 ⑦ r2 的 BLOCKED。

**P44 新增教训（v34 起）**：① 持久化优先的线性化点（纪律 K）；② ⑦ 的 FR 项必须真修，不得以注释或申辩代替改造；③ 证据包三方自洽 + 围栏同步（纪律 L）；④ 归档即原文（纪律 M）；⑤ 新增证据载荷会击穿 `G5(d)`，必须同轮补围栏；⑥ 首跑单腿红必须先 `gh run rerun --failed` 定性，再决定登记或修复。

**P45 新增教训（v35 起）**：① **零模型证明必须是机械判据集合**（`StepIntents==nil` + 默认执行目标 + 无 facts + 无派发事件），不能是「没看到调用」，更不能把常量恒等当哨兵；② **并发断言必须双向**——既要「败者全命中哨兵」，也要「赢家存活 / 槽位释放可重发」，只测失败面会漏掉「双赢家后互杀」这一整类缺陷（DR13-FIX 即此类）；③ **层边界必须写进报告口径**：console 装配层仍强制单任务，三任务 E2E 只在 chain 层贯通，**不得读作 console 已支持三任务**；④ 装配验证片可以「生产代码零改动」，此时证据全部来自新增 `_test.go`，报告须显式声明零改动并给出 `git diff --name-only -- internal/` 为空的事实。

---

## 九、清单元规则

**引用时用首句引文，不用「第 N 条」。**

| 规则（首句引文） | 内容 |
|---|---|
| 「每 2 个切片 / 待裁定变化 / 准则升版 / 队列变化 / 用户要求 / 锚点 ≥2 不符 ⇒ 重新生成」 | 触发重新生成的时机 |
| 「待裁定项必须逐项列名」 | §六 必须逐项列名 |
| 「队列项必须给出六要素」 | 标识 / 内含 / 规模分级 / ⑦ 额度 / 依赖 / 边界 |
| 「规则义务 ≠ 门禁义务」 | 不把规则文本当机械判据 |
| 「禁物理时钟断言」 | 观察项等条款不得以「N 个月后」为触发，应以事件为触发 |
| 「提交前门禁绿」 | 硬门（DR-11） |
| 「派单表 agentName 机械读取」 | ⑩ 派单表每行 agentName 必须从 transcript.arguments.agentName 逐字机械读取 |
| 「治理片不回写 DEV_PLAN」 | 准则治理切片的权威台账 = 附录 C.0–C.0x |
| 「片名一经进台账即不改」 | v16 |
| 「授权数量与政策执行冲突时以政策为准」 | 但不得以此为由判定越权或拒绝执行 |
| 「准则行号随增删动态变化」 | 引用时以节名 + 行内标识定位为准 |
| 「主控登记文本可措辞补全；审查方原文必须逐字」 | v19 |
| 「事实错误独立于设计决策」 | 台账与报告互斥属事实错误，不应等设计裁定 |
| 「同一事实只在一处描述」 | 正文或指针择一 |
| 「提交信息含非 ASCII 字符时须用 git commit -F <UTF-8 文件>」 | v20 |
| 「派单必须显式传 agentName + model 两者」 | v20 |
| 「审查输入包范围须在 ⑩ 回执声明（完整 / 收窄）」 | v20 |
| 「⑩ 回执须附『给下版清单的事实包』」 | v20 修订 4。字段：HEAD / origin/main；主提交 + 回填；CI run；台账新增行；新增观察项；新增纪律候选；§12.4 最大 OB / 最大顺延；ADR 最大号；tmp 固定断言实测；本片暴露的清单可能写错的事实；与上版清单的差异点；未闭环项 |
| 「报告 §1 异常登记区加 `<!-- gate-exempt -->` 标记」 | v21-OBS-3，未升格，待裁 |
| 「[ROUTINE] 小修不占新切片号但不减纪律」 | v32：告警/编码/临时件类小修走同一 ⑨ 门禁 + CI + 围栏补行，不得顺带改语义 |
| 「起跑消息即工件」 | v33：`00-briefing.md` 落盘 + 入 MANIFEST 与 SHA256SUMS + 围栏声明，缺此件不得进入 ① |
| 「执行载体偏离＝最高级独立性降级」 | v33：须登记 D 项并附独立影响评估 |
| 「持久化成功即线性化点」 | v34：内存 CAS 只作快速拒绝缓存；须写明「取额成功 ⟺ 标记已在盘」并配正向红测 |
| 「证据包三方自洽」 | v34：载荷集 = MANIFEST 归档条目 = SHA256SUMS 条目；新增载荷同轮补围栏，否则 G8-b / G5(d) 红 |
| 「审查方产出归档逐字」 | v34：⑤⑥⑦ 产出必须原文落盘，禁止按要点重建；摘要另立文件并披露 |
| 「机械判据优于见证性陈述」 | **v35 新增**：零模型/零副作用类断言必须是可证伪的判据集合；常量恒等要标注为非哨兵 |
| 「并发断言须覆盖存活与释放」 | **v35 新增**：除失败面外，必须断言赢家存活、槽位释放后可重发 |
| 「层边界不得扩大读法」 | **v35 新增**：某层未支持的能力必须在报告中显式声明，禁止把「另一层已贯通」读作「本层已支持」 |

---

## 十、给新会话的一句话开场

```
以下是我的项目状态快照（v35），请先核第一节 5 个锚点，相符则按第四节队列继续，不符以仓库实况为准。

【粘贴本清单全文 v35】

【补充文档】
- REMAINING_SLICES.md（台账，P45 已 ✅；DR13-FIX 已 ✅；DR max = 13 且 DR-13 已修复）
- DELIVERY_DIRECTIVE.md（v1.60）
- docs/validation/p45-e2e-zero-model{,_EN}.md（P45 报告，含 D10/D11 与 6 条放行条件）
- docs/validation/dr13-fix{,_EN}.md（DR13-FIX 报告，含 DIAGNOSIS: B 与 D1/D2）
- docs/CONTRACTS.md §4.2（就绪分级与探针许可一次性语义）与 §4.1（封存屏障与 K）
- docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/review.md（C1–C7 + N1–N4）

当前片：队列 0 = P46 战略回看（草案待起草）。
Flash 状态：[待起跑 / 执行中 / 待回执]
```

---

## 十一、当前状态一句话

P35–P45 全部闭环（`3a7541d`，v1.60）；DR13-FIX 作为 P45 前置修复片清除了并发排他窗口，P45 以零模型三任务 E2E 证明 P39–P44 的零件（预算、锁序、归属、绑定、封存、门控）在 chain 引擎层贯通且装配层有旁证腿；**当前插片 = P46 战略回看（待起草）**；`CONTRACTS-S7-DURABILITY` 顺延为队列 1，且其起跑依赖 P46 的裁定结果；触发式与低优先级队列挂起待条件。

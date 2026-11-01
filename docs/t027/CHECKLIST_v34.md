# 任务执行清单（跨会话快照 v34 · P45 起跑交接版）

生成时点：2026-11-01，HEAD = origin/main = `16abe2b`，准则 v1.60。
替代关系：完全取代 v33 及任何早于本日的快照。
本版修订（相对 v33）：① 锚点更新至 P44 全链闭环（`16abe2b`）；② 核心队列推进到 P45（零模型三任务 E2E）；③ P44 事实包（两段提交链与 2 次 CI，含首跑 flake 与 DR-13）纳入；④ P44 新增教训（持久化优先的线性化点、⑦ FR 项必须真修、证据包三方自洽、归档须原文不得重建）纳入 §八 与 §九；⑤ §七 新增两条起跑纪律；⑥ 继承 v33 的 24 项非终态观察项与 P40–P44 全部后续片候选。

---

## 一、当前锚点（开场先核这 5 个数）

| 项 | 值 | 核对命令 |
|---|---|---|
| HEAD | `16abe2b` | `git rev-parse --short HEAD` |
| origin/main | `16abe2b` | `git rev-parse --short origin/main` |
| 准则版本 | **v1.60** | `docs/DELIVERY_DIRECTIVE.md` 头部 |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest`（见下方「命令坑」） |
| 三作用域 | **全 TOTAL_FAIL=0** | `gate.js --all --scope=index\|tree\|ci` |

**若 ≥2 项不符 ⇒ 以仓库实况为准**，先让主控报终态读数再继续。

**`--selftest` 命令坑**：全量约需数分钟，`| Select-Object -Last 3` 会长时间无输出、易误判挂起。建议用 `| Select-String -Pattern 'SELFTEST'` 并容忍长时间静默；不要中途终止。

**【三作用域全绿口径】** 干净树时 `index/tree` 均 `scanned 0`（vacuous pass），只有 `ci` 扫到 N 行——须按纪律 E 报 `scanned N` 实际值。v34 生成时 `ci` 的 `G6-1 scanned 14`。

**【tmp/ 固定断言】** `tmp` 顶层 = `.gitkeep,gates-cache`；`gates-cache` 子目录数 = 1（`g4b`）；**`tmp/gates-cache/` 下 `g4b/` 之外的文件数 = 0**；`Test-Path tmp/P40-T027B-LOCK-REAPER`、`tmp/P41-T027B-SPAWN-OWNERSHIP`、`tmp/P42-T027B-EVIDENCE-BIND`、`tmp/P43-T027B-STATE-SEAL`、`tmp/P44-T027B-GATE-READINESS` 均为 `False`。

> **口径区分**：「`gates-cache` 下 `g4b` 之外 = 0」与「`tmp/` 排除 `gates-cache` 后 = 1（`.gitkeep`）」是两个不同读数，不可互相替代。

---

## 二、P44 闭环事实包（C7+N4 门控就绪与一次性许可，已入库）

### 提交链（全部已推 `origin/main`，SSH，未推 gitee）

| # | 提交 | 内容 | 规模 |
|---|---|---|---|
| 1 | **`9061b67`** | 主提交：三级 readiness 分离 + `ProbePermit` 一次性语义 + P44 用例（含变异证据）+ 报告对 + 证据包 + 台账 + 三文档回写 | 文件数 23，+2432 / −2 |
| 2 | **`16abe2b`** | CI 证据回填 + 台账闭环（完成记录仅记报告路径）+ 证据包 `08-ci/ci-evidence.txt` + DR-13 登记 | 文件数 7，+43 / −11 |

### CI（每次推送均实测；首跑红已定性并登记）

| run | head | 结论 | 备注 |
|---|---|---|---|
| `41253798647` | `9061b67` | **attempt 1 ubuntu 腿红 → attempt 2 双绿** | 红点 = `internal/adapters` 的 `TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner`（`exactly one contender may win the launch slot, got 2 (losers 6)`），P44 未触碰该包 ⇒ 登记 **DR-13**；attempt 2 两腿 `ok gates`、Race 腿 `ok adapters`/`chain` 均取到 |
| `41255254907` | `16abe2b` | 双绿 | 文档回填，观察运行 |

### 审查链

① `DESIGN: DONE`（含 §八 四条修订注记 R-1～R-4）→ ③ `RED-CONFIRMED`（编译红 `undefined: EvaluateT027BReadiness` 等）→ ② 转绿 → ④ 门禁 `TOTAL_FAIL=0`（`scanned 44`）→ ⑤ r1 `PRE-REVIEW: PASS WITH FIXES`（2 Medium：规格漂移未记录 / 无生产接线）→ r2 `PASS WITH FIXES`（3 Low）→ ⑥ `INDEPENDENT SCAN: PASS`（Section B 零缺陷，额度 1/3）→ ⑦ r1 `FINAL REVIEW: BLOCKED`（FR-01 High 崩溃窗口 / FR-02 Medium 缺红测 / FR-03 Low 证据未落盘）→ r2 `BLOCKED`（NR-01 Medium 证据包三方不自洽 / NR-02 Low 偏差未登记）→ **r3 `FINAL REVIEW: PASS`（额度 3/3）**。**R = 0、D = 0**（偏差 D1–D6 已登记）。

### 偏差登记（P44）

**D1** `GoalProduction` 缺 L2 判据取 `Seal == nil`（不看 `ProbeEvidence`），经 ⑤ 裁定为正确；**D2** `ReadinessInput` 改用哈希串替代 chain 类型 ⇒ ① 的「L3 构造签名强制先有 L2 值」类型层绑定未落地，改运行时判据；**D3** `internal/console/**` 零改动 ⇒ C7/N4 未端到端接线，接线后置；**D4** 本机无 cgo ⇒ `-race` 未执行（替代证据 `-count=50 -cpu 1,2,4` 与 `-count=20`），CI Race 腿不含 `internal/gates`；**D5** ② 误删 `tmp/` 中间件；**D6** ③ 首轮变异证据未落盘，已复算补救。

### 交付

`internal/gates/readiness.go`（`ReadinessGoal` 三值 / sealed `Readiness` / `SealBinding` / `ReleaseBinding` / `ReadinessInput` / 三具体类型与互斥访问器 / `EvaluateT027BReadiness` 纯函数 / 4 枚哨兵）；`internal/gates/probe_permit.go`（`ProbePermit` 三态 / `NewProbePermit` / `LoadProbePermit` / `Require` 持久化优先 / `Complete` / 3 枚哨兵）；测试两文件（P44 用例 20 例 + 9 子例）；变异证据 11 项（M1–M11 全 RED）。

### P44 新增教训

① **持久化优先的线性化点**（⑦ FR-01）：状态机的准入判定必须以「持久化成功」为线性化点，内存 CAS 只能作快速拒绝缓存，否则崩溃窗口会让额度复活；② **⑦ 的 FR 项必须真修而非申辩**：FR-01 的修法是改造 `Require` 顺序，不是补注释；③ **证据包三方自洽**（⑦ r2 NR-01）：文件系统载荷集、`MANIFEST.md` 归档条目、`SHA256SUMS.txt` 条目必须三方相等，缺一即 G8-b 红；④ **归档必须是原文**：审查方产出归档逐字落盘，不得按要点重建；⑤ **门禁 G5(d) 会被新载荷击穿**：新增证据文件必须同轮补入报告 `artifacts` 围栏。

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
| **v1.60** | **P44** | **`9061b67` + `16abe2b`** | **闭环（R=0、D=0；D1–D6 已登记，新增 DR-13）** |

---

## 四、执行队列（11 月代码重构）

**总原则：预算 → 锁序 → 归属 → 证据 → 封存 → 门控就绪（已完）→ 末装配。**

### 队列 1 · P45 `P45-T027B-E2E-ZERO-MODEL` 🔄 待起跑

| 字段 | 内容 |
|---|---|
| 标识 | `P45-T027B-E2E-ZERO-MODEL` |
| 目标 | 零模型三任务装配 E2E：T1 接受 → T2 自然失败 → T3 待执行；端到端贯通 `console` 装配 + `chain` 状态机，全程零模型调用 |
| 依赖 | P44 闭环 |
| 规模 | M |
| ⑦ 额度 | 1+1 |
| 边界 | 仅 `internal/chain/**`、`internal/console/**` 与 `examples/**`；不推 gitee；停 ⑩ 等授权 |

### 队列 2 · `CONTRACTS-S7-DURABILITY`（11 月立项）

| 字段 | 内容 |
|---|---|
| 标识 | `CONTRACTS-S7-DURABILITY` |
| 内含 | Windows 发布耐久升级：`CONTRACTS.md` §7 修订 + 独立评审 |
| 依赖 | 与 P45 并行 |
| 规模 | M |
| ⑦ 额度 | 1+1 |
| 边界 | 仅 `docs/CONTRACTS{,_EN}.md` 与 schema 校验；先协议后代码 |

### 队列 3 · `G4B-CORPUS-DECOUPLE`（触发式，当前挂起）

| 字段 | 内容 |
|---|---|
| 标识 | `G4B-CORPUS-DECOUPLE`（OB-120） |
| 内含 | `listAllMarkdown` 的 tree/index 用 `ls-files`，暂存删除 `.md` 时语料收缩 ⇒ G4b 冷跑从严多报 |
| 依赖 | 触发条件成立：出现「暂存删除 `.md` 后门禁假阳性」的实际误判 |
| 规模 | M |
| ⑦ 额度 | 不适用（触发式） |
| 边界 | 候选修法见 OB-120 现象列（A / B / C） |

### 队列 4 · 低优先级（长期挂起）

| 字段 | 内容 |
|---|---|
| 标识 | `OB-50` / `spawnGit` / `ci 冻结成本` / `%TEMP% 历史遗留 22 件` / `GATES-EXEMPT` |
| 内含 | 逐项见 §12.4 对应 OB 行；`GATES-EXEMPT` = P33 的 ①-3(iv) 拆片，未立项，触 4 个判据文件 |
| 依赖 | 长期挂起，无明确触发 |
| 规模 | 不等（含 S 到 M-L） |
| ⑦ 额度 | 无 |
| 边界 | 逐项见 §12.4 |

> 队列 1–2 为已排期项，步骤 0 登记台账时以该片 16 字段块的边界为准；队列 3–4 不占当期排期。

### 后续片候选（P39–P44，不并入队列 1–4）

**P39 ⑤ 四项非阻塞发现**：持锁进程崩溃后的锁释放语义无测试腿；`rename` 成功后报错的可见分支无测试腿；`newDurableAck` 的 `confirmed` 字段命名与语义不符；矩阵「确认前取消」腿未落地。

**P40 报告登记的后续片候选**：F-07a 静态锁序断言入 CI 判据；F-11 `lockPath` canonical 化；Reclaim 生产接线（uncertain 路径）；CI Race leg 扩包覆盖 `internal/guard` 与 `internal/tickets`。

**P41 报告登记的后续片候选**：**I1** kill-on-close 半证；**I2** 四缝注入（Start 后 / attach 前 / OnSpawn 内 / resume 前）；⑦ F-01 拒绝态三态化；⑦ F-02 rename-away 对象绑定删除；⑦ F-03 mirror 删除断言；Pdeathsig 独立切片；⑤ Low-1 槽位新鲜度兜底评估。

**P42 报告登记的后续片候选**：**N2 读路径 sweep**（`hashLogs` / `readPIDs` / `fileHash` / `scanWorkspace`）；Unix 两段式 dirfd；Windows reparse 对象身份级校验；`CloseBoundDirs` 生产接线；darwin 旧文档 supersede；原语层测试补齐（关闭后拒读 / 超界拒摘要 / 身份同源对抗）。

**P43 报告登记的后续片候选**：生产装配接线 `Attempts`（console 侧 `TaskAttemptResolver` 来源）；supervised attempt>1 端到端；seal 目录级耐久 crash 注入；stale/mismatch 判定矩阵补测；`router.go` 与 handoff 的 K 绑定；`WorkspaceReopener` 落地；`StepRequest.ContextKey` 的读取者落地。

**P44 报告登记的后续片候选**：**console 接线（R-3）**——`runtime_agent.go:103`、`runtime_supervised.go:248`、`runtime.go:33` 与 `ai.go:39-90` 注入 `ReadinessInput`，解除端到端口径限制；**Race 腿扩展（R-4）**——把 `./internal/gates/...` 纳入 CI Race 腿（属 `ci.yml` 变更）；**类型层绑定（R-2）**——`ProductionReadiness` 构造签名接收 L2 值，恢复 ① 的类型层约束；⑥ 未覆盖面（`LoadProbePermit` 非 `NotExist` 读失败 / `Require` 的 `ctx == nil` / `Complete` 未知态 / 真跨进程双写）。

**低优先级（长期挂起）**：`OB-50` / `spawnGit` / `ci 冻结成本` / `%TEMP% 历史遗留` / `GATES-EXEMPT`。

---

## 五、观察项

### 观察中（12）
`OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113`

### 待议（11）
`OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`

### 待办（1）
`OB-69`

**非终态合计 24 项**（观察中 12 + 待议 11 + 待办 1）；本轮未改 §12.4，读数与 v33 一致；其余按 §12.4 检索。

### 已裁决（28，含既有窗口新增）
OB-117–OB-124（顺延第 16–22 次）、OB-125（ADR-028）、OB-126/127/128（第 23/24/25 次）、OB-129–OB-131（Astra 三次授权）、OB-132（① 位角色-模型绑定外授权，P39）。P40–P44 窗口的用户裁定（方案 A、F 项全盘、DR-12、各片起跑授权、主控模型策略、P43 疑点裁定、P44 起跑与闭环授权）属切片级裁定，未入 OB。

### 机械读数（可核对口径）
- §12.4 最大 OB = **OB-132**；最大顺延 = 第 25 次；ADR max = **28**；DR max = **13**（台账实测）。
- §12.4 条目数 132；`已处置 75 / 已裁决 28 / 已登记 5 / 观察中 12 / 待议 11 / 待办 1`。

---

## 六、待用户裁定

| # | 项 | 现状 | 阻塞谁 |
|---|---|---|---|
| 1 | **P45 起跑授权**（含 ⑦ 额度、边界、批量辅助次数） | 待裁（v34 审核通过后起跑） | P45 |
| 2 | P44 后续片（R-3 console 接线 / R-4 Race 腿 / R-2 类型层绑定）的立项归属 | 待裁 | P45 之后的装配片 |
| 3 | CONTRACTS-S7-DURABILITY 立项节奏 | 待裁 | C6 的 Windows 落地 |
| 4 | DR-13 是否升格为独立修复切片 | 待裁（当前观测级） | CI 稳定度 |
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
7. 批量辅助：每片实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` **≥ 2 次**（建议 3–4 次），并在报告 §模型与计量 回填实际次数。
8. 主控模型策略：默认 `DeepSeek V4.1 Flash`；遇底层并发 / 状态机 / 句柄类复杂片、单份超 10,000 字符设计文档、大量 CI 日志深度根因分析、或 Flash 连续 2 次行为漂移时提示切换 Pro；Max 不作默认。**任何降级须按 OB-37 元规则独立登记**。
9. **角色载体不可替代（P43 D7 起，硬规则）**：② 必须由 `prfrail-implementer`、③ 必须由 `prfrail-tester` 派单执行；主控只做编排、落盘、门禁与回执，**不得代写产品代码或测试**。违反即登记为最高级独立性降级（D 项 + 独立影响评估）。
10. **起跑消息即工件（用户 2026-11-01 裁定，硬规则）**：步骤 0 建目录时，必须将用户的起跑消息原文（剥离附件清单后）落盘为 `docs/validation/evidence/<切片ID>/00-briefing.md`，并加入证据包的 `MANIFEST.md` 与 `SHA256SUMS.txt`，同时在报告的 `artifacts` 围栏中声明。缺此件不得进入 ①。
11. **线性化点纪律（P44 ⑦ FR-01 起，硬规则）**：凡涉及持久化状态机，**「持久化成功」必须是线性化点**；内存 CAS 只能作同进程快速拒绝缓存，不得作为准入闸门。注释与契约须写明「取额成功 ⟺ 标记已在盘」的不变式，并配一条「成功即落盘」的正向红测。
12. **证据包与围栏同步纪律（P44 ⑦ r2 NR-01 + G5(d)）**：新增任何证据载荷必须**同轮**完成三件事——① 补 `MANIFEST.md` 归档条目；② 重算 `SHA256SUMS.txt`；③ 在报告 `artifacts` 围栏补行。任缺一项都会被 `G8-b` 或 `G5(d)` 拦下。
13. **审查方归档纪律（P44 起）**：⑤⑥⑦ 的产出归档必须**逐字原文**；主控不得按要点重建，若确需摘要，摘要另立文件并在围栏外披露。

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
| 批量辅助 | 只用 `GLM-5.3-Flash (glm)`；**≥ 2 次，建议 3–4 次**分别派单 |
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

**P38 新增教训（v28 起）**：单个子代理单次输出有硬上限；长文档必须分批（每批 ≤ 10,000 字符、单主题），派单写明「禁止自重启」与 `SELF-RESTART-ABORT`。

**P40 新增教训（v30 起）**：① 交叉构建后 `go env GOOS GOARCH` 复核；② CI 与本机分歧先查工具链版本（DR-12）；③ ⑥ Section B 的 `INFO` 行不是 severity 行。

**P41 新增教训（v31 起）**：① Section B 写 `Low | …` 管道形态会触发 G7-f 的 severity 行闭集；② Windows 持锁删除需 `FILE_SHARE_DELETE`；③ Job 属性缺 `EXTENDED_STARTUPINFO_PRESENT` 被静默忽略；④ 自写 `CreateProcessW` 必须配 `PROC_THREAD_ATTRIBUTE_HANDLE_LIST`。

**P42 新增教训（v32 起）**：① 「本机未使用」≠「全仓未使用」，删符号前须对向 `GOOS` 复核；② `//go:build !windows` 覆盖 darwin，不可用平台配 fail-closed 桩；③ gopls 的 `unusedfunc` / `unusedparams` 会针对构建标签变体报；④ 单腿红灯先怀疑构建标签排除；⑤ 围栏转录必须逐字；⑥ 批量辅助是强制项。

**P43 新增教训（v33 起）**：① 主控不得代写 ②③；② ⑥ 的 Section B 语义必须收紧（纪律 J）；③ 额度按轮次簿记；④ 偏差须附独立影响评估；⑤ 「声称已改」必须用机械判据复核（精确符号检索 + 编译/测试），否则会重演 ⑦ r2 的 BLOCKED。

**P44 新增教训（v34 起）**：① 持久化优先的线性化点（纪律 K）；② ⑦ 的 FR 项必须真修，不得以注释或申辩代替改造；③ 证据包三方自洽 + 围栏同步（纪律 L）；④ 归档即原文（纪律 M）；⑤ 新增证据载荷会击穿 `G5(d)`，必须同轮补围栏；⑥ 首跑单腿红必须先 `gh run rerun --failed` 定性，再决定登记或修复。

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
| 「持久化成功即线性化点」 | **v34 新增**：内存 CAS 只作快速拒绝缓存；须写明「取额成功 ⟺ 标记已在盘」并配正向红测 |
| 「证据包三方自洽」 | **v34 新增**：载荷集 = MANIFEST 归档条目 = SHA256SUMS 条目；新增载荷同轮补围栏，否则 G8-b / G5(d) 红 |
| 「审查方产出归档逐字」 | **v34 新增**：⑤⑥⑦ 产出必须原文落盘，禁止按要点重建；摘要另立文件并披露 |

---

## 十、给新会话的一句话开场

```
以下是我的项目状态快照（v34），请先核第一节 5 个锚点，相符则按第四节队列继续，不符以仓库实况为准。

【粘贴本清单全文 v34】

【补充文档】
- REMAINING_SLICES.md（台账，P44 已 ✅，DR max = 13）
- DELIVERY_DIRECTIVE.md（v1.60）
- docs/validation/p44-gate-readiness{,_EN}.md（P44 报告，含 D1–D6）
- docs/CONTRACTS.md §4.2（就绪分级与探针许可一次性语义的权威条款）
- docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/review.md（C1–C7 + N1–N4）

当前片：队列 1 = P45-T027B-E2E-ZERO-MODEL（零模型三任务 E2E）。
Flash 状态：[待起跑 / 执行中 / 待回执]
```

---

## 十一、当前状态一句话

P35–P44 全部闭环（`16abe2b`，v1.60，P44 含两段提交链与 2 次 CI，其中首跑 ubuntu 腿红由 `internal/adapters` 既有 flake 引起并登记 DR-13，⑥/⑦ 额度用满且终判 PASS）；11 月重构队列推进，**下一个起跑片 = P45（零模型三任务 E2E）**：T1 接受 → T2 自然失败 → T3 待执行，全程零模型调用，且 ②③ 一律派子代理执行。

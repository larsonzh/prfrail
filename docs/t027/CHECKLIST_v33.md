# 任务执行清单（跨会话快照 v33 · P44 起跑交接版）

生成时点：2026-11-01，HEAD = origin/main = `8a46b01`，准则 v1.60。
替代关系：完全取代 v32 及任何早于本日的快照。
本版修订（相对 v32）：① 锚点更新至 P43 全链闭环（`8a46b01`）；② 核心队列推进到 P44（门控就绪，C7+N4）；③ P43 事实包（四段提交链与 5 次 CI）纳入；④ P43 新增教训（主控不得代写 ②③、⑥ Section B 语义、额度簿记、偏差须附影响评估）纳入 §八；⑤ 新增硬纪律「起跑消息即工件」；⑥ 继承 v32 的 24 项非终态观察项与 P40–P43 全部后续片候选。

---

## 一、当前锚点（开场先核这 5 个数）

| 项 | 值 | 核对命令 |
|---|---|---|
| HEAD | `8a46b01` | `git rev-parse --short HEAD` |
| origin/main | `8a46b01` | `git rev-parse --short origin/main` |
| 准则版本 | **v1.60** | `docs/DELIVERY_DIRECTIVE.md` 头部 |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest`（见下方「命令坑」） |
| 三作用域 | **全 TOTAL_FAIL=0** | `gate.js --all --scope=index\|tree\|ci` |

**若 ≥2 项不符 ⇒ 以仓库实况为准**，先让主控报终态读数再继续。

**`--selftest` 命令坑**：全量约需数分钟，`| Select-Object -Last 3` 会长时间无输出、易误判挂起。建议用 `| Select-String -Pattern 'SELFTEST'` 并容忍长时间静默；不要中途终止。

**【三作用域全绿口径】** 干净树时 `index/tree` 均 `scanned 0`（vacuous pass），只有 `ci` 扫到 N 行——须按纪律 E 报 `scanned N` 实际值。

**【tmp/ 固定断言】** `tmp` 顶层 = `.gitkeep,gates-cache`；`gates-cache` 子目录数 = 1（`g4b`）；**`tmp/gates-cache/` 下 `g4b/` 之外的文件数 = 0**；`Test-Path tmp/P40-T027B-LOCK-REAPER`、`tmp/P41-T027B-SPAWN-OWNERSHIP`、`tmp/P42-T027B-EVIDENCE-BIND`、`tmp/P43-T027B-STATE-SEAL` 均为 `False`。

> **口径区分**：「`gates-cache` 下 `g4b` 之外 = 0」与「`tmp/` 排除 `gates-cache` 后 = 1（`.gitkeep`）」是两个不同读数，不可互相替代。

---

## 二、P43 闭环事实包（C4+C5+N3 证据封存与幂等重放，已入库）

### 提交链（全部已推 `origin/main`，SSH，未推 gitee）

| # | 提交 | 内容 | 规模 |
|---|---|---|---|
| 1 | **`b6e0f34`** | 主提交：seal 屏障 + 真实 attempt + K + 10 例测试 + 报告对 + 证据包 + 台账 + 三文档回写 | 文件数 32，+1796 / −110 |
| 2 | `609bb19` | 报告 §提交与 CI 回填 | 文件数 2 |
| 3 | `7004282` | D7 执行载体偏离登记 | 文件数 2 |
| 4 | **`8a46b01`** | ⑥ 三轮全貌 + D8 严重度误用 + D7b 影响评估 + ⑥ r3 未覆盖面 | 文件数 2 |

### CI（每次推送均实测）

| run | head | 结论 | 备注 |
|---|---|---|---|
| `41228033279` | `b6e0f34` | **双绿** | 两腿 `changed=32`、`SELFTEST 697/697`、`TOTAL_FAIL=0`；ubuntu `ok chain 0.348s` / `console 1.904s` / `tickets 1.999s` + Race `chain 2.012s`；windows `ok chain 5.607s` / `console 22.474s` / `tickets 5.034s` |
| `41228611379` | `609bb19` | 双绿 | 文档回填，观察运行 |
| `41229538385` | `7004282` | 绿 | D7 登记，观察运行 |
| `41240794403` | `8a46b01` | 绿 | ⑥ 三轮与 D8 登记，观察运行 |

### 审查链

① `DESIGN: DONE` → ③ `RED-CONFIRMED`（编译红 `undefined: SealRecord`）→ ② 两轮转绿 → ④ 门禁 → ⑤ r1 `PRE-REVIEW: PASS WITH FIXES`（2 High + 4 Medium + 2 Low + 4 Info）→ r2 `PRE-REVIEW: PASS` → ⑥ r1 `INDEPENDENT SCAN: BLOCKED`（K 未接线）→ 修复 → r2 `PASS WITH FINDINGS`（Section B 严重度误用：条目实为正面核验）→ **r3 `INDEPENDENT SCAN: PASS`（Section B 零缺陷，额度 3/3）** → ⑦ r1 `FINAL REVIEW: BLOCKED`（F1–F6）→ r2 `BLOCKED`（F1 未落地）→ **r3 `FINAL REVIEW: PASS`（额度 3/3）**。**R = 0、D = 0**。

### 交付

`internal/chain/seal.go`（`SealRecord` 含 `ContextKey` / `SealSink` / `SealTaskEvidence` O_EXCL + fsync + 失败删半写 / `LoadTaskSeal` / `sameClosure`）；`seal_verify.go`（`New` 与 `Recover` 双点校验）；`attempt.go`（`TaskAttemptResolver` / `resolveTaskAttempt` / `projectionAttempt` / `taskAttempts` / `ContextKey` / `TaskContextKey`）；`engine.go`（`acceptFromExistingSeal` 含 K 与 decision 校验、`sealAcceptedClosure`、step 闭包校验、`StepRequest.ContextKey`）；attempt 贯通 dispatch / routing / operator / facts / cancel / supervised；测试 10 例。

### P43 新增教训

① **主控不得代写 ②③**（D7：⑦ r2 因主控补丁未落盘而 BLOCKED）；② ⑥ 的 Section B 只允许真实缺陷，核验结论归 Section C，严重度不得用于「已实现」陈述（D8）；③ 额度必须按轮次精确簿记（⑥ 3/3、⑦ 3/3）；④ 偏差须同时给独立影响评估，不能只「声明接受」（D7b）。

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
| **v1.60** | **P43** | **`b6e0f34` + `609bb19` + `7004282` + `8a46b01`** | **闭环（R=0、D=0；D7/D7b/D8 已登记）** |

---

## 四、执行队列（11 月代码重构）

**总原则：预算 → 锁序 → 归属 → 证据 → 封存（均已完）→ 门控就绪 → 末装配。**

| # | 标识 | 内含 | 规模 | ⑦ 额度 | 依赖 | 边界 |
|---|---|---|---|---|---|---|
| 1 | **P44-T027B-GATE-READINESS** | **C7+N4**：`EvaluateT027BReadiness` 分离探针 / 真实验收 / 生产发布三级；`ProbePermit` 一次性消费 | M | 1+1 | P43 | 仅 `internal/console/**` 与 `internal/gates/**`；不推 gitee；停 ⑩ 等授权 |
| 2 | P45-T027B-E2E-ZERO-MODEL | 零模型三任务装配 E2E：T1 接受 → T2 自然失败 → T3 待执行 | M | 1+1 | P44 | 仅 `internal/chain/**`、`internal/console/**` 与 examples；停 ⑩ 等授权 |
| 3 | CONTRACTS-S7-DURABILITY（11 月立项） | Windows 发布耐久升级：`CONTRACTS.md` §7 修订 + 独立评审 | M | 1+1 | 与 P44 并行 | 仅 `docs/CONTRACTS{,_EN}.md` 与 schema 校验；先协议后代码 |

> 上表边界为规划级；每片步骤 0 登记台账时以该片 16 字段块的边界为准。

**P39 ⑤ 四项非阻塞发现**：持锁进程崩溃后的锁释放语义无测试腿；`rename` 成功后报错的可见分支无测试腿；`newDurableAck` 的 `confirmed` 字段命名与语义不符；矩阵「确认前取消」腿未落地。

**P40 报告登记的后续片候选**：F-07a 静态锁序断言入 CI 判据；F-11 `lockPath` canonical 化；Reclaim 生产接线（uncertain 路径）；CI Race leg 扩包覆盖 `internal/guard` 与 `internal/tickets`。

**P41 报告登记的后续片候选**：**I1** kill-on-close 半证；**I2** 四缝注入（Start 后 / attach 前 / OnSpawn 内 / resume 前）；⑦ F-01 拒绝态三态化；⑦ F-02 rename-away 对象绑定删除；⑦ F-03 mirror 删除断言；Pdeathsig 独立切片；⑤ Low-1 槽位新鲜度兜底评估。

**P42 报告登记的后续片候选**：**N2 读路径 sweep**（`hashLogs` / `readPIDs` / `fileHash` / `scanWorkspace`）；Unix 两段式 dirfd；Windows reparse 对象身份级校验；`CloseBoundDirs` 生产接线；darwin 旧文档 supersede；原语层测试补齐（关闭后拒读 / 超界拒摘要 / 身份同源对抗）。

**P43 报告登记的后续片候选**：生产装配接线 `Attempts`（console 侧 `TaskAttemptResolver` 来源）；supervised attempt>1 端到端；seal 目录级耐久 crash 注入；stale/mismatch 判定矩阵补测；`router.go` 与 handoff 的 K 绑定；`WorkspaceReopener` 落地；`StepRequest.ContextKey` 的读取者落地。

**低优先级（长期挂起）**：`OB-50` / `spawnGit` / `ci 冻结成本` / `%TEMP% 历史遗留` / `GATES-EXEMPT`。

---

## 五、观察项

### 观察中（12）
`OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113`

### 待议（11）
`OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`

### 待办（1）
`OB-69`

**非终态合计 24 项**（观察中 12 + 待议 11 + 待办 1）；其余按 §12.4 检索。

### 已裁决（28，含既有窗口新增）
OB-117–OB-124（顺延第 16–22 次）、OB-125（ADR-028）、OB-126/127/128（第 23/24/25 次）、OB-129–OB-131（Astra 三次授权）、OB-132（① 位角色-模型绑定外授权，P39）。P40–P43 窗口的用户裁定（方案 A、F 项全盘、DR-12、各片起跑授权、主控模型策略、P43 疑点裁定）属切片级裁定，未入 OB。

### 机械读数（可核对口径）
- §12.4 最大 OB = **OB-132**；最大顺延 = 第 25 次；ADR max = **28**；DR max = **12**（台账实测）。
- §12.4 条目数 132；`已处置 75 / 已裁决 28 / 已登记 5 / 观察中 12 / 待议 11 / 待办 1`。

---

## 六、待用户裁定

| # | 项 | 现状 | 阻塞谁 |
|---|---|---|---|
| 1 | **P44 起跑授权**（含 ⑦ 额度、边界、批量辅助次数） | 已授权（2026-11-01 起跑消息） | P44 |
| 2 | CONTRACTS-S7-DURABILITY 立项节奏 | 待裁 | C6 的 Windows 落地 |
| 3 | C7/N4 的 PARTIAL 是否可作起点 | 随 P44 起跑消息裁定为本片范围 | P44 验收标准 |
| 4 | OB-55 归属与 T027-B 立项节奏 | 待裁 | 后续实施 |
| 5 | MVA、S1 移出项与 S2 承诺收窄批准 | 待裁 | 后续切片 |
| 6 | 候选与 Windows 前提 P 及资源/探针授权 | 待裁 | G2 门 |
| 7 | AT-24 `unknown` 与「未知副作用保持暂停」接受边界 | 待裁 | AT-24 收口 |
| 8 | **P41/P42/P43 后续片候选的立项归属**（I1/I2、N2 sweep、attempt 装配接线等） | 待裁 | 相关片范围 |

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
| 临时目录只增不删 | 执行角色在 `tmp/<切片ID>/` 内只增不删；只能由 ⑧c 清理 |
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

**P38 新增教训（v28 起）**：单个子代理单次输出有硬上限；长文档必须分批（每批 ≤ 10,000 字符、单主题），派单写明「禁止自重启」与 `SELF-RESTART-ABORT`。

**P40 新增教训（v30 起）**：① 交叉构建后 `go env GOOS GOARCH` 复核；② CI 与本机分歧先查工具链版本（DR-12）；③ ⑥ Section B 的 `INFO` 行不是 severity 行。

**P41 新增教训（v31 起）**：① Section B 写 `Low | …` 管道形态会触发 G7-f 的 severity 行闭集；② Windows 持锁删除需 `FILE_SHARE_DELETE`；③ Job 属性缺 `EXTENDED_STARTUPINFO_PRESENT` 被静默忽略；④ 自写 `CreateProcessW` 必须配 `PROC_THREAD_ATTRIBUTE_HANDLE_LIST`。

**P42 新增教训（v32 起）**：① 「本机未使用」≠「全仓未使用」，删符号前须对向 `GOOS` 复核；② `//go:build !windows` 覆盖 darwin，不可用平台配 fail-closed 桩；③ gopls 的 `unusedfunc` / `unusedparams` 会针对构建标签变体报；④ 单腿红灯先怀疑构建标签排除；⑤ 围栏转录必须逐字；⑥ 批量辅助是强制项。

**P43 新增教训（v33 起）**：① 主控不得代写 ②③；② ⑥ 的 Section B 语义必须收紧（纪律 J）；③ 额度按轮次簿记；④ 偏差须附独立影响评估；⑤ 「声称已改」必须用机械判据复核（精确符号检索 + 编译/测试），否则会重演 ⑦ r2 的 BLOCKED。

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
| 「起跑消息即工件」 | **v33 新增**：`00-briefing.md` 落盘 + 入 MANIFEST 与 SHA256SUMS + 围栏声明，缺此件不得进入 ① |
| 「执行载体偏离＝最高级独立性降级」 | v33 新增：须登记 D 项并附独立影响评估 |

---

## 十、给新会话的一句话开场

```
以下是我的项目状态快照（v33），请先核第一节 5 个锚点，相符则按第四节队列继续，不符以仓库实况为准。

【粘贴本清单全文 v33】

【补充文档】
- REMAINING_SLICES.md（台账，P43 已 ✅）
- DELIVERY_DIRECTIVE.md（v1.60）
- docs/validation/p43-state-seal{,_EN}.md（P43 报告，含 D1–D8）
- docs/CONTRACTS.md §4.1（封存屏障与不可变上下文摘要 K 的权威条款）
- docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/review.md（C1–C7 + N1–N4）

当前片：队列 1 = P44-T027B-GATE-READINESS（门控就绪，C7+N4）。
Flash 状态：[待起跑 / 执行中 / 待回执]
```

---

## 十一、当前状态一句话

P35–P43 全部闭环（`8a46b01`，v1.60，P43 含四段提交链与 5 次 CI，⑥/⑦ 额度用满且终判 PASS）；11 月重构队列推进，**下一个起跑片 = P44（门控就绪，C7+N4）**：`EvaluateT027BReadiness` 三级分离 + `ProbePermit` 一次性消费，且 ②③ 一律派子代理执行。

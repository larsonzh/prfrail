# 任务执行清单（跨会话快照 v32 · P43 起跑交接版）

生成时点：2026-10-31，HEAD = origin/main = `d3a511b`，准则 v1.60。
替代关系：完全取代 v31 及任何早于本日的快照。
本版修订（相对 v31）：① 锚点更新至 P42 全链闭环（`d3a511b`）；② 核心队列推进到 P43（状态封存，C4+C5+N3）；③ P42 三段提交链与 6 次 CI 读数入事实包；④ P42 新增教训（跨 GOOS 删除校验、平台约束按支持面收口、gopls 变体死码、单腿红灯的版本分裂根因）纳入 §八；⑤ 继承 v31 的 24 项非终态观察项与 P40/P41/P42 全部后续片候选。

---

## 一、当前锚点（开场先核这 5 个数）

| 项 | 值 | 核对命令 |
|---|---|---|
| HEAD | `d3a511b` | `git rev-parse --short HEAD` |
| origin/main | `d3a511b` | `git rev-parse --short origin/main` |
| 准则版本 | **v1.60** | `docs/DELIVERY_DIRECTIVE.md` 头部 |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest`（见下方「命令坑」） |
| 三作用域 | **全 TOTAL_FAIL=0** | `gate.js --all --scope=index\|tree\|ci` |

**若 ≥2 项不符 ⇒ 以仓库实况为准**，先让主控报终态读数再继续。

**`--selftest` 命令坑**：全量约需数分钟，`| Select-Object -Last 3` 会长时间无输出、易误判挂起。建议用 `| Select-String -Pattern 'SELFTEST'` 并容忍长时间静默；不要中途终止。

**【三作用域全绿口径】** 干净树时 `index/tree` 均 `scanned 0`（vacuous pass），只有 `ci` 扫到 N 行——须按纪律 E 报 `scanned N` 实际值。

**【tmp/ 固定断言】** `tmp` 顶层 = `.gitkeep,gates-cache`；`gates-cache` 子目录数 = 1（`g4b`）；**`tmp/gates-cache/` 下 `g4b/` 之外的文件数 = 0**；`Test-Path tmp/P40-T027B-LOCK-REAPER` = `False`；`Test-Path tmp/P41-T027B-SPAWN-OWNERSHIP` = `False`；`Test-Path tmp/P42-T027B-EVIDENCE-BIND` = `False`。

> **口径区分**：「`gates-cache` 下 `g4b` 之外 = 0」与「`tmp/` 排除 `gates-cache` 后 = 1（`.gitkeep`）」是两个不同读数，不可互相替代。

---

## 二、P42 闭环事实包（C3+N2 句柄与目录绑定 + 静态告警清理，已入库）

### P42-T027B-EVIDENCE-BIND（C3+N2，v1.60）

- **提交链**（全部已推 `origin/main`，SSH，未推 gitee）：

| # | 提交 | 内容 | 规模 |
|---|---|---|---|
| 1 | **`cab3e08`** | 主提交：绑定读原语 + 目录绑定原语 + 读取点替换 + 报告对 + 证据包 + 台账 | 文件数 34，+2350 / −50 |
| 2 | **`be612fe`** | 报告 §提交与 CI 回填 | 文件数 2，+10 / −2 |
| 3 | `a3998dc` | [ROUTINE] 三处静态告警修复 + D7 补登记 + tmp 探针清理 | 文件数 5，+31 / −26 |
| 4 | `cd57142` | 修复 ubuntu 腿红灯：`spawnRecoverySettleBound` 移入 unix 测试文件 | 文件数 1，+5 |
| 5 | `cf90c53` | darwin fail-closed 桩（`ErrBoundDirUnsupported`）+ 约束收口 | 文件数 5，+75 / −9 |
| 6 | **`d3a511b`** | 绑定名守卫移出 darwin 变体（`bounddir_clean.go`）+ 桩参数显式化 | 文件数 5，+38 / −25 |

- **CI（每次推送均实测）**：

| run | head | 结论 | 备注 |
|---|---|---|---|
| `41194972405` | `cab3e08` | **双绿** | 主体；两腿 `changed=34`、`SELFTEST 697/697`、`TOTAL_FAIL=0` |
| `41195571714` | `be612fe` | 双绿 | 文档回填，观察运行不计轮次 |
| `41207400308` | `a3998dc` | **ubuntu 红 / windows 绿** | `vet: spawn_recovery_live_host_unix_test.go:92: undefined: spawnRecoverySettleBound` |
| `41208049995` | `cd57142` | 双绿 | 修复后；ubuntu `ok adapters 12.253s` / windows `ok adapters 37.761s` |
| `41209821231` | `cf90c53` | 双绿 | darwin 桩；三平台 build/vet 静默 |
| **`41211542105`** | **`d3a511b`** | **双绿** | 最终锚点；两腿 `changed=5`、`SELFTEST 697/697`、`TOTAL_FAIL=0` |

- **审查链**：① 架构 `DESIGN: DONE`（含修订注记 D1/D4/F1）→ ③ 红测 `RED-CONFIRMED`（I1–I5）→ ② 两轮转绿（主体 + F 修复轮，含 `confirmOpenedArtifact` 逐读复核）→ ④ 门禁 → ⑤ r1 `PRE-REVIEW: PASS WITH FIXES`（1 Medium + 3 Low + 2 Info）→ ⑤ r2 `PRE-REVIEW: PASS`（追加 N1–N6 登记要求）→ ⑥ `INDEPENDENT SCAN: PASS`（首轮无 A–E 段不合规，重发后合规，额度 2/3）→ ⑦ `FINAL REVIEW: PASS WITH FIXES`（3 Low + 4 Info，非阻断）。**R = 0、D = 0**。
- **交付**：`internal/evidence` 新增 `BoundFile`（`OpenBoundFile` / `OpenBoundFileTruncating` / `NewBoundFile` / `ReadAll` / `Hash` / `Identity` / `Close`；bound+1 探测截断、仅全量读后给摘要、`ErrNotExist` 透传）；`internal/adapters` 新增 `BoundDir`（Windows 目录句柄 + `GetFinalPathNameByHandle` 真名 + 文件句柄最终路径复核；Unix `openat` + `O_NOFOLLOW` + `O_CLOEXEC`，`dev/ino` 校验）；读取点替换与 `os.Stat` 探针删除；`bounddir_clean.go` 守卫（`!darwin`）；`bounddir_darwin.go` fail-closed 桩。报告对 `docs/validation/p42-evidence-bind{,_EN}.md`（行数 197，六形状 197/16/66/237/32/4 全等）；证据包 5 件载荷 + MANIFEST + SHA256SUMS 自洽；artifacts 围栏条目数 38；台账 P42 所在行 ✅ 已闭环。
- **P42 新增教训**：
  ① **跨 GOOS 凭证**：`internal/adapters` 的常量/函数「本机未使用」不等于全仓未使用——unix 专属测试文件在 Windows 不参与编译，删符号前必须 `GOOS=linux go vet` 复核（`spawnRecoverySettleBound` 即此型）；
  ② **平台约束按真实支持面收口**：`//go:build !windows` 会覆盖 darwin，而标准库 darwin 无 `syscall.Openat`；不可用平台应配 fail-closed 桩而非退化为按路径打开；
  ③ **gopls 变体诊断**：`unusedfunc` / `unusedparams` 会针对构建标签变体报；平台专属守卫应放在平台专属文件（`bounddir_clean.go` 用 `!darwin`）；
  ④ **单腿红灯先查版本分裂**：ubuntu `vet` 红、windows 绿 ⇒ 先看是否文件被构建标签排除（DR-12 的构建标签分支）；
  ⑤ 交叉构建后必须 `go env GOOS GOARCH` 复核（沿用 P40）。
- **G4b 白名单新增**：`兄` / `弟`（① 架构原文 `D1` 实测依据逐字归档；依 2026-10-25 裁定「既定原文登记、不改写」）。

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
| v1.60 | P40 | `3fab6ba` + `3af9c9a` + `a142eaa` | 闭环（R=0、D=0；DR-12 已修） |
| v1.60 | P41 | `8799015` + `caaf64d` | 闭环（R=0、D=0） |
| **v1.60** | **P42** | **`cab3e08` + `be612fe` + `d3a511b`**（含 `a3998dc`/`cd57142`/`cf90c53` 告警清理尾链） | **闭环（R=0、D=0）** |

---

## 四、执行队列（11 月代码重构）

**总原则：先预算（已完）、再锁序（已完）、再归属（已完）、续证据（P42 已完）、末装配。**

| # | 标识 | 内含 | 规模 | ⑦ 额度 | 依赖 | 边界 |
|---|---|---|---|---|---|---|
| 1 | **P43-T027B-STATE-SEAL** | **C4+C5+N3**：证据封存屏障先于 intent；上下文摘要 K；`Reopen` 传真实 attempt；证据封存与状态机漏洞收口 | L | 1+1 | P39、P42 | 仅 `internal/chain/**` 与 `internal/tickets/**`；不推 gitee；停 ⑩ 等授权 |
| 2 | P44-T027B-GATE-READINESS | C7+N4：`EvaluateT027BReadiness` 分离探针 / 真实验收 / 生产发布 + `ProbePermit` 一次性 | M | 1+1 | P43 | 仅 `internal/console/**` 与 `internal/gates/**`；不推 gitee；停 ⑩ 等授权 |
| 3 | P45-T027B-E2E-ZERO-MODEL | 零模型三任务装配 E2E：T1 接受 → T2 自然失败 → T3 待执行，单 Engine / 单事件链 | M | 1+1 | P44 | 仅 `internal/chain/**`、`internal/console/**` 与 examples；不推 gitee；停 ⑩ 等授权 |
| 4 | CONTRACTS-S7-DURABILITY（11 月立项） | Windows 发布耐久升级：`CONTRACTS.md` §7 修订 + 独立评审 | M | 1+1 | 与 P43 并行 | 仅 `docs/CONTRACTS{,_EN}.md` 与 schema 校验；先协议后代码；独立评审 |

> 上表边界为规划级；每片步骤 0 登记台账时以该片 16 字段块的边界为准。

**P39 ⑤ 四项非阻塞发现（随 #4 或后续重构片处理）**：持锁进程崩溃后的锁释放语义无测试腿；`rename` 成功后报错的可见分支无测试腿；`newDurableAck` 的 `confirmed` 字段命名与语义不符；矩阵「确认前取消」腿未落地。

**P40 报告登记的后续片候选（待用户裁立项归属）**：F-07a 静态锁序断言入 CI 判据；F-11 `lockPath` canonical 化；Reclaim 生产接线（`agent_runner_pinned_cli.go` uncertain 路径）；CI Race leg 扩包覆盖 `internal/guard` 与 `internal/tickets`。

**P41 报告登记的后续片候选（待用户裁立项归属）**：**I1** kill-on-close 半证（宿主崩溃 ⇒ 挂起树被内核击杀的测试半）；**I2** 四缝注入（Start 后 / attach 前 / OnSpawn 内 / resume 前各注入一次宿主崩溃）；⑦ F-01 拒绝态三态化（拒绝与无 marker 可区分）；⑦ F-02 rename-away 对象绑定删除；⑦ F-03 mirror 删除断言；Pdeathsig 独立切片（缩短 Linux 最坏回收延迟）；⑤ Low-1 槽位新鲜度兜底评估。

**P42 报告登记的后续片候选（待用户裁立项归属）**：**N2 读路径 sweep**（`hashLogs` / `readPIDs` / `fileHash` / `scanWorkspace` 的绑定读迁移）；Unix 两段式 dirfd（中间组件符号链接向量）；Windows reparse 对象身份级校验（替代前缀比较）；`CloseBoundDirs` 生产接线 + 注释口径统一；darwin 旧文档 supersede（`docs/validation/t027-replay-convergence.md` 交叉编译条目）；原语层测试补齐（`BoundFile` 关闭后拒读 / `Hash` 超界拒摘要 / 身份同源对抗）。

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
OB-117–OB-124（顺延第 16–22 次）、OB-125（ADR-028）、OB-126/127/128（第 23/24/25 次）、OB-129（Astra 首次授权）、OB-130（Astra 二次授权）、OB-131（Astra 三次授权）、OB-132（① 位角色-模型绑定外授权，P39）。P40/P41/P42 窗口的用户裁定（方案 A、F 项全盘、DR-12、各片起跑授权、主控模型 Pro→Flash 回归与长期策略）属切片级裁定，未入 OB。

### 机械读数（可核对口径）
- §12.4 最大 OB = **OB-132**；最大顺延 = 第 25 次；ADR max = **28**；DR max = **12**（台账实测）。
- §12.4 条目数 132；`已处置 75 / 已裁决 28 / 已登记 5 / 观察中 12 / 待议 11 / 待办 1`。

---

## 六、待用户裁定

| # | 项 | 现状 | 阻塞谁 |
|---|---|---|---|
| 1 | **P43 起跑授权**（含 ⑦ 额度、边界、批量辅助次数） | 待裁 | P43 |
| 2 | CONTRACTS-S7-DURABILITY 立项节奏 | 待裁 | C4/C6 的 Windows 落地 |
| 3 | C4/C5/C7 的 PARTIAL 是否可作起点 | 待裁 | P43/P44 验收标准 |
| 4 | OB-55 归属与 T027-B 立项节奏 | 待裁 | 后续实施 |
| 5 | MVA、S1 移出项与 S2 承诺收窄批准 | 待裁 | 后续切片 |
| 6 | 候选与 Windows 前提 P 及资源/探针授权 | 待裁 | G2 门 |
| 7 | AT-24 `unknown` 与「未知副作用保持暂停」接受边界 | 待裁 | AT-24 收口 |
| 8 | **P41/P42 后续片候选的立项归属**（I1/I2、N2 sweep、原语测试补齐等，并入 P43–P45 或另立片） | 待裁 | 相关片范围 |

---

## 七、起跑与步骤 0 通则（代码切片）

1. 核锚点 → 步骤 0：登记本片台账（`R/D = 0` 基线、临时件落点 `tmp/<切片ID>/`）、预取输入包并落 `input-list.md` 与 `freedom-list.md`。
2. 每个修复切片必须给出**回归测试**（对 P38 报告列出的测试名逐条落地或说明替代），并把「未执行测试/lint」显式登记。
3. ⑨ 三作用域全绿 + `SELFTEST 697/697` + `go build/vet/test ./...` 全绿，方可提交。
4. 模型绑定（P39 教训，硬约束）：① `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`；⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`；②③⑧ `DeepSeek V4.1 Flash (deepseek)`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`（常规档）。**每次 `runSubagent` 必须显式同时传 `agentName` 与 `model`，严禁擅自升级或替换**。
5. **执行节奏纪律（用户 2026-10-31 裁定，硬规则）**：**常规不停**——④ 门禁全绿 ⇒ 直接派 ⑤；⑤ 报 Medium+ ⇒ 修复后**原位重跑该环**；**⑩ 才是唯一常规停点**。仅在遇到架构冲突、超预算或无解时停下升级。
6. **跨版本 / 跨平台排查纪律（DR-12 + P42）**：CI 与本机行为分歧时，优先核对工具链版本与**构建标签覆盖面**；交叉构建后必须 `go env GOOS GOARCH` 复核，防 GOOS 残留；删除「本机未使用」符号前必须做至少一次对向 `GOOS` 的 `go vet`。
7. 批量辅助：每片实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` **≥ 2 次**（建议 3–4 次），并在报告 §模型与计量 回填实际次数。
8. 主控模型策略（用户 2026-10-31 裁定）：默认 `DeepSeek V4.1 Flash`（High 档）；遇底层并发 / 状态机 / 句柄类复杂片、单份超 10,000 字符设计文档、大量 CI 日志深度根因分析、或 Flash 连续 2 次行为漂移时，提示切换 Pro 或获授权后切换；Max 不作默认。**任何降级须按 OB-37 元规则独立登记，并在报告如实标注独立性降低**。

---

## 八、跨片协作纪律（起跑消息必须落实）

| 要素 | 内容 |
|---|---|
| 载体表 | ① ⑤ ⑥ ⑦ 用 `deep-reasoner` / `independent-reviewer`（只读）；② ③ ⑧ 用 `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`（可写）；禁止 `prfrail-tester` 跑任何审查层环节；批量辅助载体必须是 `deep-reasoner`，禁止 `quick-researcher` |
| 派单表机械读取 | ⑩ 派单表每行 `agentName` 必须从 `transcript.arguments.agentName` 逐字机械读取，附脚本 stdout 原文 |
| R2.1 补 agentName | 每次 `runSubagent` 的 `arguments` 必须同时含 `agentName` 与 `model`，两者都不得缺。任务包与回执都写全 |
| 提交信息含非 ASCII 须用 -F 文件 | 提交信息含任何非 ASCII 字符（中文 / 全角符号 / § 等）时，须用 `git commit -F <UTF-8 文件>`，不得用 `-m` 内联。纯 ASCII 提交建议用 `-m` 避免临时文件违规 |
| 「可独立核对 vs 不可独立核对」判别标准 | 可独立核对的产出（枚举、事实比对、引用点清单），模型错误不构成作废理由；不可独立核对的产出（判断、结论、审查意见），模型错误必须重做 |
| ⑥ 先于 ⑦ | ⑥ 必须在 ⑦ 之前完成；重派 ⑥ 时若文档有变更，⑦ 结论须重新评估 |
| 临时件位置 | 一律 `tmp/<切片ID>/`；逐件精确清理 + 回读证明；整改轮同样落此目录 |
| 临时目录只增不删 | 执行角色在 `tmp/<切片ID>/` 内只增不删；该目录只能由 ⑧c 清理。违反即触发 §9.3「② 提前清理」同型偏离计数 |
| 长命令重定向 | 不要用 PowerShell `>`；用 `\| Tee-Object -FilePath <path>` 或 `cmd /c "..."` |
| 门禁执行 | 中间环节用 `--all --skip-base`；提交前必须全绿（DR-11 红线） |
| 纪律 E · index 作用域的正确顺序 | 先精确 `git add` → 再跑 `--scope=index`。回执中若报告「门禁全绿」必须同时报 `G6-1 scanned N` 且 N > 0 |
| 纪律 F · 规范优先的适用边界 | 当起跑消息的「具体示例文本」与「措辞规范段落」冲突时，执行角色按规范执行，但须回执登记冲突 |
| 批量辅助 | 只用 `GLM-5.3-Flash (glm)`；**≥ 2 次，建议 3–4 次**分别派单 |
| 审查方只读纪律 | ⑦ 不得创建、修改、删除任何文件（含 tmp/）；数据来源只能写「读取主控落盘的 <路径>」 |
| 纪律 A · 推定授权 | ⑩ 停点回执的 CI 取证需要 push，视为该片授权范围内的必要动作；若用户在同轮明确写「不 push」或「push 另行授权」，则以用户措辞为准。推定授权可接受，但不得作先例 |
| 纪律 B · 证据缺口表述 | 证据缺口必须如实登记「丢失了什么、什么未丢」，并明写「未主张证据缺口已消除」 |
| 纪律 C · 枚举优先于推导 | 机械读数优先于算术反推；遇到「总量 − 已知 = 差值」的推导时，多跑一条枚举命令 |
| 纪律 D · 白名单政策优先于授权数量 | 授权数量与政策执行冲突时，以政策为准，但主控在执行授权时须在派单/任务包文本中预置「数量可能超出预估」的说明 |
| 纪律 G · 围栏即落地快照 | 报告 `artifacts` 围栏按路径排序、TAB 分隔、CN/EN 逐字节相同；每次触及新文件（含 `[ROUTINE]` 小修、工具白名单）必须同步补行，否则 G5(d) 报 staged-not-declared |
| 纪律 H · 提交前对向 GOOS 复核 | 删除或改动任何常量/函数前，至少跑一次对向 `GOOS` 的 `go vet`（P42 `spawnRecoverySettleBound` 事件）；改动平台文件后跑三平台 `build + vet` |

**P38 新增教训（v28 起适用）**：单个子代理单次输出有硬上限；要求长文档时**必须先分批**（每批 ≤ 10,000 字符、每批单主题），并在派单文本写明「禁止自重启」与 `SELF-RESTART-ABORT` 逃生标记。

**P40 新增教训（v30 起适用）**：① 交叉构建后同一终端 `go env GOOS GOARCH` 复核（GOOS 残留 ⇒ 测试二进制报「不是有效的 Win32 应用程序」）；② CI 与本机分歧先查工具链版本（DR-12）；③ ⑥ 六段围栏 Section B 的 `INFO` 行不是 severity 行，`INDEPENDENT SCAN: PASS` + 零 severity 行自洽。

**P41 新增教训（v31 起适用）**：① ⑥ 的 Section B 写 `Low | …` 管道形态会触发 G7-f 的 severity 行闭集（severity 行 ⇔ FINDINGS）⇒ 报告围栏内分隔符归一为 `·` + 围栏外逐字披露，verbatim 原文保留在证据包；② Windows 持锁删除需 `FILE_SHARE_DELETE` 句柄；③ `CreateProcessW` 的 Job 属性缺 `EXTENDED_STARTUPINFO_PRESENT` 会被静默忽略；④ 自写 `CreateProcessW` 必须配 `PROC_THREAD_ATTRIBUTE_HANDLE_LIST` 限制继承。

**P42 新增教训（v32 起适用）**：① 「本机未使用」≠「全仓未使用」——unix 专属测试文件在 Windows 不参与编译，删符号前必须 `GOOS=linux go vet`；② `//go:build !windows` 覆盖 darwin，不可用平台配 fail-closed 桩而非退化为按路径打开；③ gopls `unusedfunc` / `unusedparams` 会对构建标签变体报 ⇒ 平台专属守卫放平台专属文件；④ 单腿红灯（ubuntu 红 / windows 绿）先怀疑构建标签排除，再怀疑代码；⑤ 报告与围栏的 ⑥ 转录必须逐字，围栏外才做措辞归一；⑥ 批量辅助是强制项（本片 2 次，发现原语层三处测试缺口）。

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
| 「⑩ 回执须附『给下版清单的事实包』」 | v20 修订 4。字段：HEAD / origin/main；本片主提交 + 回填（hash + 变更面）；本片 CI run；台账新增行；新增观察项；新增纪律候选；§12.4 最大 OB / 最大顺延次数；ADR 最大号；tmp 固定断言实测；本片暴露的、清单可能写错的事实；与上版清单的差异点；本片发现的未闭环项 |
| 「报告 §1 异常登记区加 `<!-- gate-exempt -->` 标记」 | v21-OBS-3：报告 §1 异常登记文字本身会触发 G2/G6-4/G4b ⇒ **未升格**，待用户裁定 |
| 「[ROUTINE] 小修不占新切片号但不减纪律」 | v32 新增：告警/编码/临时件类小修须走同一 ⑨ 门禁 + CI + 围栏补行，且不得顺带改语义 |

---

## 十、给新会话的一句话开场

```
以下是我的项目状态快照（v32），请先核第一节 5 个锚点，相符则按第四节队列继续，不符以仓库实况为准。

【粘贴本清单全文 v32】

【补充文档】
- REMAINING_SLICES.md（台账，P42 已 ✅）
- DELIVERY_DIRECTIVE.md（v1.60）
- docs/validation/p42-evidence-bind{,_EN}.md（P42 报告）
- docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/02-astra-redteam/review.md（C1–C7 判定 + N1–N4）
- docs/DEV_PLAN.md（已含 11 月重构里程碑映射）

当前片：队列 1 = P43-T027B-STATE-SEAL（状态封存，C4+C5+N3）。
Flash 状态：[待起跑 / 执行中 / 待回执]
```

---

## 十一、当前状态一句话

P35–P42 全部闭环（`d3a511b`，v1.60，P42 含三段提交链与 6 次 CI 实测，末次 `41211542105` 双绿）；11 月重构队列推进，**下一个起跑片 = P43（状态封存，C4+C5+N3）**：证据封存屏障先于 intent、上下文摘要 K、`Reopen` 传真实 attempt。

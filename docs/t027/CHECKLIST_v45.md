# 任务执行清单（跨会话快照 v45 · P55 闭环版）

生成时点：2026-11-05，HEAD = origin/main = `c76c6ce`（P55 回写提交），准则 v1.60。P55 已闭环并推送 origin（`63e2e24` 主提交 + `c76c6ce` 回写，均未推 gitee）。
替代关系：完全取代 v44 及任何更早快照。回填提交（docs-only）只在 HEAD 前进，不改本节锚点语义。

本版修订（相对 v44）：① P55 `EXIT-SEED-BUILD-EVIDENCE` 闭环（`63e2e24` 主提交 + `c76c6ce` 回写；run `41925623071` 与 `41928327517` 双绿；跨平台证据包 53 载荷 + 三角色两代构建）；② §二 事实包换为 P55，P54 归入 §三；③ §三 追加 P55 该切片条目；④ §四 第 7 项出队，**第 8 项 `AT23-RESTRICTED-CANDIDATE` 标为 ⛔ 阻断（依赖 R-08 / R-09）**，并把挂起项与待评估项提上台面作备选；⑤ §五 机械读数按实况复核（OB max 136、顺延 25、ADR 30、DR 13）；⑥ §六 R-02 / R-03 / R-05 / R-11 记为已完成，R-08 / R-09 / R-10 保持待裁；⑦ §九 门禁坑新增 P55 四条（隔离模块缓存污染 `gofmt`、停止态措辞触发 `G7-d`、捕获 `COMMAND` 行保真、跨平台三角色结构）。

## 一、当前锚点（开场先核这 5 个数）

| 项 | 值 | 核对命令 |
|---|---|---|
| HEAD | `c76c6ce` | `git rev-parse --short HEAD` |
| origin/main | `c76c6ce` | `git rev-parse --short origin/main` |
| 准则版本 | **v1.60** | `docs/DELIVERY_DIRECTIVE.md` 头部 |
| SELFTEST | **697/697** | `node tools/gates/gate.js --selftest` |
| 三作用域 | **全 TOTAL_FAIL=0** | `gate.js --all --scope=index`、`tree`、`ci` |

若 ≥2 项不符 ⇒ 以仓库实况为准，先报终态读数再继续。

`tmp` 固定断言：顶层 = `.gitkeep` + `gates-cache`；`gates-cache` 下仅 `g4b`；P55 的 `tmp/EXIT-SEED-BUILD-EVIDENCE/` 已于 ⑧c 清理。

## 二、P55 事实包

| 字段 | 值 |
|---|---|
| 切片 / 规模 | `EXIT-SEED-BUILD-EVIDENCE`（证据链补证类 `[SLICE]`）/ M |
| 目标 | 补齐 AT-13 证据链（seed 与 candidate 分离；外部 oracle 比较；clean-room 成功及失败不损 seed），**不改 CI 结构、不改产品代码** |
| 改动面 | 仅 `.md` 与证据包：证据包 `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/`（53 载荷 + 2 元数据）+ 报告对 `docs/validation/exit-seed-build-evidence{,_EN}.md` + 6 对权威文档同步 + 台账 P55 条目；**零 `.go`、零 workflow** |
| 核心产出 | ① 双侧捕获（Windows `go1.27.0` / Linux VM `go1.22.12`）；② **三角色两代结构**：seed `fc99f55…` 构建 `prfrail`、verifier `18c9340…` 构建 `prfrail-release`、candidate `baa144d…` 构建 `prfrail`；③ `OracleHash=sha256:388f8884…` 四面一致（双侧 + 夹具 + `s1-selfhost.md` CI 输入表）；④ 四条反例双侧 fail-close；⑤ 失败后 seed 树哈希 `equal=true`；⑥ Linux 构建期 `GOPROXY=off`（先镜像预取后断网） |
| 审查链 | ③ 两轮（`prfrail-tester`）：`TEST: DONE`（0C/0H/6M/6L）→ 6 项 Medium 以 `claim-scope.md` + `anomalies.md` 处置 → r2 `TEST: DONE`（6/6 闭合，R2-1..R2-6 已修）→ ⑤ 两轮（`deep-reasoner` + `DeepSeek V4 Pro`）：`PRE-REVIEW: PASS WITH FIXES`（1 Medium = 捕获 `COMMAND` 显示缩写 → E-1 勘误 + 派生 index 命令修正）→ ⑥（`independent-reviewer` + `MAI-Code-1.1-Flash`）：`INDEPENDENT SCAN: PASS`（Section B `(no findings)`）→ ⑦（`independent-reviewer` + `GPT-5.3-Codex`）：`FINAL REVIEW: PASS WITH FIXES`（1 Low 措辞，已修正）→ ⑧a 文档同步（`prfrail-documenter`）→ ⑨ 三作用域 + SELFTEST → ⑧c 清理 → ⑩ 停点 → 用户授权后提交推送 |
| 提交 / CI | `63e2e24`（`test: add AT-13 seed-driven build evidence across Windows and Linux (P55)`，69 项）+ `c76c6ce`（`docs: record P55 slice completion and CI evidence`，4 项）；run `41925623071`（63e2e24）双绿、run `41928327517`（c76c6ce）双绿；ubuntu 腿含 Race 与 Contract fixtures；三个条件作业按设计 skipped |
| 计量 | `R = 0`、`D = 1`（无触发词抢跑，经用户强制打断纠正，就地转正）；`runSubagent` 共 9 次（③×2、⑤×2、⑥、⑦、⑧a、批量辅助×2）；探针/付费调用 0 |
| 裁定 | 用户批准 A-4（三角色替代「seed 跑 selfhost」，因 `fc99f55` 无发布工具）并认为独立性更强；批准 A-5（`goproxy.cn` 预取 + 构建期 `GOPROXY=off`）；裁定 A-3（② 主控自执行）合规；要求批量辅助必须补派（实得 2/2，`CONSISTENT` + `PARITY: OK`） |
| 文档同步 | `TEST_STRATEGY`（AT-13 补证登记 + 边界指针）、`DEV_PLAN`（T017 进度）、`OPERATIONS` / `S1_SUPPORT_MATRIX` / `S1_RELEASE_NOTES`（三项「不适用 + 理由」）、`REMAINING_SLICES`（台账回写）；CN/EN numstat 逐对对称 |
| 能力边界 | `claim-scope.md`：**S-1..S-7 成立 / N-1..N-7 不 assert**（Windows 核验者与候选同源 ⇒ 不证独立性；`SeedHash` 为组装目录哈希；oracle 为结构化比较；行为一致 ≠ 无恶意） |
| 未闭环 | 无硬性未闭环；观察项：Windows 工具链 `go1.27.0` ≠ CI `1.22.12`（N-6）、Windows 二进制未哈希（F-08）、反例命令含占位符（N-4/E-1）、E-2 文案勘误 |

## 三、已完成切片（近期基线）

| 版本 | 切片 | 主提交 + 回填 | 结论 |
|---|---|---|---|
| v1.60 | P49 | `75990a8` + 回填提交 | 闭环（R = 1、D = 2；run `41445657595` 首跑双绿） |
| v1.60 | P50 | `81f9dc9` + 回填提交 | 闭环（R = 1、D = 1；run `41507722672`） |
| v1.60 | P51 | `2b0af2d` + 回填提交 | 闭环（R = 1、D = 2；run `41569353460`） |
| v1.60 | P52 | `c636dd8` + 回填提交 | 闭环（R = 0、D = 1；run `41642139932`；GAP-77 击穿实证） |
| v1.60 | P53 | `7802819` + 回填提交 | 闭环（R = 0、D = 0；run `41707160563`；GAP-77 已修复） |
| v1.60 | P54 | `7c0763a` + `f38e0ce` | 闭环（R = 1、D = 1；run `41840557379`；C6 判定路径 + 降级裁定） |
| v1.60 | P55 | `63e2e24` + `c76c6ce` | 闭环（R = 0、D = 1；run `41925623071` + `41928327517` 双绿；AT-13 双侧证据链 + 三角色两代构建） |

## 四、执行队列

总原则：层 1 清记录链 → 层 2 战略主线 → 层 3 发布前门 → 层 4 出口门；层内序固定，层间不跳跃。

### 层 1（已清空）

无剩余项。

### 层 2（战略主线）

无剩余项（P53 已出队，层 2 彻底清空）。

### 层 3 / 层 4（当前主线，已被阻断）

第 7 项 P55 已出队（`63e2e24` 闭环）。**下一片 = 第 8 项 `AT23-RESTRICTED-CANDIDATE`，当前 ⛔ 阻断（依赖 R-08 / R-09 裁定）。**

| 序 | 切片 | 标识与内含 | 规模 | ⑦ 额度 | 依赖 | 边界 |
|---|---|---|---|---|---|---|
| 8 | AT23-RESTRICTED-CANDIDATE | AT-23 受限真实候选与签收 | M | Codex 1 次终审 | **R-08 / R-09 裁定（未裁 ⇒ ⛔ 阻断）** | 无真实候选授权不得签收 |

### 备选（阻断解除前唯一可动的两类，均需用户点名）

| 类别 | 项 | 说明 | 起跑条件 |
|---|---|---|---|
| 挂起项（触发式） | `G4B-CORPUS-DECOUPLE` | 语料与判定解耦；触发式挂起 | 用户点名或触发条件成立 |
| 待评估项（治理侧） | P7 治理侧拆出 OB 汇总 | 「后续批次切片」节待评估项；含 OB-77 等拆出项 | 按 §1.8.1 先起评估动作 |
| 待评估项（产品侧） | P5 `T027-B` / P6 `S1-EXIT-DOC` 收口清单 | 产品侧收口；字段只有三项，须先补 §1.2 全表 | 同为 §1.8.1 评估动作 |
| 低优先级治理 | 五项低优先级任务 | 无触发不启动；须用户显式点名 | 用户点名 |

> 口径：以上四类**都不是**层 3 主线，不得读作「层 3 恢复」；`AT23-RESTRICTED-CANDIDATE` 的阻断不因它们起跑而解除。

## 五、观察项

- 观察中（13）：`OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135`。
- 待议（11）：`OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75 / 120`。
- 待办（1）：`OB-69`。非终态合计 **25**。
- GAP 登记：**GAP-77 已修复（P53 闭环）**；GAP-78/79/80 候选（登记于 P52 报告 §2b 与台账）。
- P55 新增观察（不占 OB 号，属切片级登记）：① Windows 工具链与 CI 不同形（N-6）；② 捕获 `COMMAND` 行保真类问题（E-1/E-2，已勘误）；③ `index.json` 曾出现派生态滞后（已重建为 v2 口径）。
- 机械读数：`§12.4` 最大 OB = **OB-136**；顺延最大 = 第 25 次；ADR max = **30**；DR max = **13**。

## 六、待用户裁定（R-01 至 R-11）

| 编号 | 项 | 现状 |
|---|---|---|
| R-01 | P46 战略回看立项 | 已完成 |
| R-02 | `CONTRACTS-S7-DURABILITY` 节奏与归属 | 已完成（P54 闭环） |
| R-03 | P45 后续片归属 | 已完成（`RACE-LEG-EXPAND` 闭环） |
| R-04 | DR-13「已修复 + 观察」口径 | 待裁 |
| R-05 | DR13-FIX-D3 是否独立立项 | 已完成（独立立项并闭环） |
| R-06 | OB-55 归属与 T027-B 节奏 | 待裁 |
| R-07 | MVA、S1 移出项与 S2 承诺收窄 | 待裁 |
| R-08 | 候选、Windows 前提 P 与探针授权 | **待裁（阻塞第 8 项 `AT23-RESTRICTED-CANDIDATE`）** |
| R-09 | AT-24 `unknown` 接受边界 | **待裁（阻塞第 8 项签收口径）** |
| R-10 | 总预算与各片模型上限 | 待裁（阻塞全部排片） |
| R-11 | 历史纠偏范围 | 已完成（B 档 + 逐项差异化） |

## 七、起跑与步骤 0 通则

1. 核锚点 → 步骤 0：登记本片台账（`R/D = 0` 基线、临时件落点 `tmp/<切片ID>/`）、落 `input-list.md` 与 `freedom-list.md`。
2. 每个修复切片必须给出**回归测试**，并把「未执行测试/lint」显式登记。
3. ⑨ 三作用域全绿 + `SELFTEST 697/697` + `go build/vet/test ./...` 全绿，方可提交。
4. 模型绑定（硬约束）：①⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`；②③⑧ `DeepSeek V4.1 Flash (deepseek)`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`；战略回看轮 1/3 `independent-reviewer` + `GPT-6.1 Sol (copilot)`；批量辅助 `deep-reasoner` + `GLM-5.3-Flash (glm)`。**每次 `runSubagent` 必须同时传 `agentName` 与 `model`**。
5. **执行节奏纪律**：常规不停——④ 全绿 ⇒ 直接派 ⑤；⑤/⑥/⑦ 报 Medium+ ⇒ 修复后原位重跑该环；**⑩ 才是唯一常规停点**。
6. 跨版本/跨平台排查：CI 与本机分歧先查工具链版本与构建标签覆盖；交叉构建后 `go env GOOS GOARCH` 复核；删「本机未使用」符号前至少一次对向 `GOOS` 的 `go vet`。
7. 批量辅助：每片实派 `deep-reasoner` + `GLM-5.3-Flash (glm)` **≥ 2 次**（策略与输出上限见 §八）；**严格串行**，禁止并行子代理。
8. 主控模型策略：默认 `DeepSeek V4.1 Flash`；遇底层并发/状态机/句柄类复杂片、单份超 10,000 字符设计文档、大量 CI 日志深度根因分析、或 Flash 连续 2 次行为漂移时切换 Pro；**任何降级须按 OB-37 元规则独立登记**。
9. **角色载体不可替代**：② 必须 `prfrail-implementer`、③ 必须 `prfrail-tester` 派单；主控只做编排、落盘、门禁与回执。**例外**：纯证据补证/抓日志类切片可由主控自执行，但须用户同轮明示（P55 先例：起跑消息【六】授权）。
10. **起跑消息即工件**：步骤 0 必须把用户起跑消息原文落盘为 `docs/validation/evidence/<切片ID>/00-briefing.md`，并入 `MANIFEST.md` 与 `SHA256SUMS.txt`，同时在 `artifacts` 围栏声明；缺此件不得进入 ①。
11. **线性化点纪律**：涉及持久化状态机时「持久化成功」必须是线性化点，内存 CAS 只作同进程快速拒绝缓存。
12. **证据包与围栏同步纪律**：新增证据载荷必须同轮补 `MANIFEST.md` 条目、重算 `SHA256SUMS.txt`、在 `artifacts` 围栏补行。
13. **审查方归档纪律**：⑤⑥⑦ 产出必须逐字原文归档（P55 先例：`04-reviews/`）；主控不得按要点重建。
14. **装配验证片纪律**：断言必须落在可证伪的机械判据上；常量恒等须显式声明非哨兵；报告必须声明生产代码是否零改动。
15. **评估类切片 ⑨ 后必停 ⑩**：评估类无实现对象，⑨ 全绿后停在 ⑩ 等同一轮授权。
16. **单次输出不得超过 10,000 字符**：超限前分批；禁止自重启；按需读取。
17. **压缩即结片**：上下文到 100K 或出现压缩提示 ⇒ 立即结束当前片，一会话只跑一个 L 档切片。
18. **起跑消息必须含「文档同步更新」章节**：若切片改变了项目阶段状态或契约能力，⑧ 阶段须同步更新权威文档；漏写即视为结构缺陷。
19. **环境预检与复用铁律**：需外部环境时须在步骤 0 显式声明；缺工具须先探查（PATH、配置、`find`）再最小化安装；**远程 SSH 必须 `bash -lc` 加载登录环境并探查用户级目录**；安装须固定与 CI 对齐的版本并核验校验和。P55 复用经验：VM `10.0.0.199` 的 `~/go1.22/bin/go` = go1.22.12（登录 shell 默认是更新的 1.27.1，须显式钉住）。
20. **证据捕获的可重放性与保真（v45 新增）**：捕获件的 `COMMAND:` 行必须与真实 argv 一致；若只能写显示缩写，须在 `anomalies.md` 立勘误条目并在派生 `index.json` 写入忠实命令 + `commandErrata` 指针；漂移/构造类输入件须在 `replay.md` 给出可重放重建步骤。

## 八、跨片协作纪律（起跑消息必须落实）

| 要素 | 内容 |
|---|---|
| 载体表 | ①⑤⑥⑦ 用 `deep-reasoner` / `independent-reviewer`（只读）；②③⑧ 用 `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`（可写）；禁止 `prfrail-tester` 跑审查层环节；批量辅助必须是 `deep-reasoner`，禁止 `quick-researcher` |
| 派单表机械读取 | ⑩ 派单表每行 `agentName` 必须从 `transcript.arguments.agentName` 逐字机械读取，附脚本 stdout 原文 |
| 派单串行 | `runSubagent` 串行阻塞，禁止并行调用多个子代理；每次传 `agentName` 与 `model` |
| 提交信息 | 一律英文，格式 `<type>: <summary>`；含任何非 ASCII 字符时须用 `git commit -F <UTF-8 文件>` |
| 可独立核对 vs 不可独立核对 | 可独立核对的产出（枚举、事实比对、引用点清单）模型错误不构成作废；不可独立核对者（判断、结论、审查意见）必须重做 |
| ⑥ 先于 ⑦ | ⑥ 必须在 ⑦ 之前完成；重派 ⑥ 且文档有变更时，⑦ 结论须重新评估 |
| 临时件位置与清理 | 一律 `tmp/<切片ID>/`；只能由 ⑧c 清理，且只删本片目录，严禁递归删父目录；**测字符数或记录中间过程一律使用内联脚本，严禁在任何磁盘路径（尤其 `$env:TEMP`）生成临时文件。** |
| 长命令重定向 | 不用 PowerShell `>`；用 `\| Tee-Object -FilePath <path>` 或 `cmd /c "..."` |
| 门禁执行 | 中间环节用 `--all --skip-base`；提交前必须全绿 |
| 纪律 E · index 作用域顺序 | 先精确 `git add` → 再跑 `--scope=index`；回执须报 `G6-1 scanned N` 且 N 大于 0 |
| 纪律 F · 规范优先边界 | 起跑示例与规范冲突时按规范执行，但须回执登记冲突 |
| 纪律 G · 围栏即落地快照 | `artifacts` 围栏按路径排序、TAB 分隔、CN/EN 逐字节相同；触及新文件必须同步补行 |
| 纪律 H · 对向 GOOS 复核 | 删除或改动常量/函数前至少跑一次对向 `GOOS` 的 `go vet` |
| **纪律 O · 输出长度防御** | 单次输出不得超过 10,000 字符；超限前分批；禁止自重启 |
| **纪律 P · 主控模型切换触发** | 见 §七 第 8 条；降级按 OB-37 登记 |
| **纪律 Q · 批量辅助输出上限** | 派单须写明字符上限；回执须报实得字符 |
| **纪律 R · 时态滞后批量清账** | 台账与事实不一致走一次性清账片，以代码与提交为事实源 |
| **纪律 S · 会话长度与压缩纪律** | 上下文达到 100K 或出现压缩提示 ⇒ 立即结束当前片 |
| **纪律 T · 量词形态** | 台账与报告禁写「标识后紧跟量词」形态：会被 `G6-4` 读作行数字面量，应写「该条目」或「该记录」 |
| **纪律 U · 禁逐字引用触发串** | 报告不得逐字引用触发判据的字符串 |
| **纪律 V · 大任务拆分与半截防污染** | 大任务拆分为小步串行派发；超时后必须 `git status`/`git diff` 检查半截产物，必要时精确 `git checkout` 回退 |
| **纪律 W · 双语报告结构对齐（v45 新增）** | 报告对必须同序同型：`##` 标题数量一致、`artifacts` 围栏逐字节相同、`sixstage` 围栏五段齐全且末行结论一致；建议每片以机械脚本自检（P55 先例：批量辅助 `PARITY: OK`） |

**GLM 批量辅助策略**：`GLM-5.3-Flash (glm)` 资源充足，「≥ 2 次」仅是下限，机械核对、长文档摘录、双语镜像检查、派单草案、日志提取应尽量外派；禁止复杂判断与择优。**P55 教训：主控用脚本自查不能替代补派**——欠派须在同轮补齐并留回执。

## 九、清单元规则

**引用时用首句引文，不用「第 N 条」。**

| 规则（首句引文） | 内容 |
|---|---|
| 「每 2 个切片 / 待裁定变化 / 准则升版 / 队列变化 / 用户要求 / 锚点 ≥2 不符 ⇒ 重新生成」 | 触发重新生成的时机 |
| 「待裁定项必须逐项列名」 | §六 必须逐项列名 |
| 「队列项必须给出六要素」 | 标识 / 内含 / 规模分级 / ⑦ 额度 / 依赖 / 边界 |
| 「规则义务 ≠ 门禁义务」 | 不把规则文本当机械判据 |
| 「禁物理时钟断言」 | 观察项等条款以事件为触发 |
| 「提交前门禁绿」 | 硬门 |
| 「派单表 agentName 机械读取」 | 从 `transcript.arguments.agentName` 逐字读取 |
| 「治理片不回写 DEV_PLAN」 | 权威台账 = 附录 C.0–C.0x |
| 「片名一经进台账即不改」 | 片名稳定 |
| 「授权数量与政策冲突时以政策为准」 | 不得以此判定越权或拒绝执行 |
| 「同一事实只在一处描述」 | 正文或指针择一 |
| 「主控登记文本可措辞补全；审查方原文必须逐字」 | 两侧标准不同 |
| 「事实错误独立于设计决策」 | 台账与报告互斥属事实错误 |
| 「⑩ 回执须附『给下版清单的事实包』」 | 字段含 HEAD、CI run、台账新增、新增观察项、新增纪律候选、最大 OB / 顺延、ADR/DR 最大号、tmp 断言、与上版差异、未闭环项 |
| 「[ROUTINE] 小修不占新切片号但不减纪律」 | 同一 ⑨ 门禁 + CI + 围栏补行 |
| 「起跑消息即工件」 | `00-briefing.md` 落盘 + 入 MANIFEST 与 SHA256SUMS + 围栏声明 |
| 「执行载体偏离＝最高级独立性降级」 | 须登记 D 项并附独立影响评估 |
| 「机械判据优于见证性陈述」 | 零模型/零副作用类断言须可证伪 |
| 「并发断言须覆盖存活与释放」 | 失败面之外须断言赢家存活、槽位释放后可重发 |
| 「层边界不得扩大读法」 | 某层未支持的能力必须显式声明 |
| 「评估类切片 ⑨ 后必停 ⑩」 | 禁止自行 commit/push |
| 「单次输出不得超过 10000 字符」 | 超限前分批；禁止自重启 |
| 「压缩即结片」 | 上下文 100K 或出现压缩提示 ⇒ 立即结束当前片 |
| 「批量辅助下限不可由主控自查替代（v45）」 | 实派次数不足 ⇒ 同轮补派，不得以「脚本已核」豁免 |

### 门禁坑（P48 六条 + P49 七条 + P52 两条 + P55 四条）

### P48 起（沿用）

1. **`tmp/` 下的 `.go` 副本会被 `go build ./...` 编译**：审查输入包中的源码副本一律用 `*.go.txt`。
2. **`G7-b` 围栏处置词闭集**：未入库临时件写 `local-only`，写 `local` 会报 unknown disposition。
3. **报告出现 `INDEPENDENT SCAN` 标记必须配 `sixstage` 围栏**：围栏内 `## Section A` 至 `## Section E` 各一次，结论行只保留块内最后一个非空行，并在围栏外披露合并差异。
4. **写 `.md` 必须显式补 BOM**：用 node 回写时漏传 BOM 参数会立刻触发 `G4a`。
5. **`G8-b` 要求 MANIFEST 逐件列出载荷**：列目录行会被判「payload not listed」。
6. **读数必须最后回填**：收档会改变 `changed=` 计数与 `scanned=` 读数。

### P49 起（沿用）

7. **改 `.github/workflows/*.yml` 必须同步 `internal/release` 工作流冻结契约**。
8. **远程 SSH 预检**：非交互会话不加载用户 PATH，必须 `bash -lc` 并探查 `~/.local/bin` 等用户级目录。
9. **PowerShell → ssh 嵌套引号**：跨壳复杂命令一律写远端 shell 脚本上传执行。
10. **⑥ 输出必须 A–E 形状**：`Section A`（范围）→ `B`（发现或 `(no findings)`）→ `C`（证据）→ `D`（可证伪性审计）→ `E`（无法验证部分 + 推翻条件）→ 块末结论行。
11. **围栏与磁盘枚举必须一致**：新增载荷后须同轮补围栏行 + MANIFEST 条目 + 重算 SHA256SUMS。
12. **`gofmt -l .` 会扫进隐藏目录与 `tmp/`**：见 P55 第 16 条。
13. **重算哈希用只读脚本**：编辑器占用载荷时写式归一化会 `EBUSY`。

### P52 起（沿用）

14. **小节标题的状态符号必须与状态行同步**：`G1(d)` 会比对两处闭集一致性。
15. **行号统一后必须重算 SHA256SUMS**：内容类勘正会改变载荷哈希。

### P55 起（新增）

16. **`go build`/`go test` 在 `tmp/<切片ID>/` 内建隔离 `GOCACHE`/`GOMODCACHE` 会污染 `gofmt -l .`**：模块缓存含第三方 `.go` 源码（如 `gopkg.in/yaml.v3`），`base gofmt` 立即判红。处置：捕获完成后整目录删除 **且清除会话级 `$env:GOCACHE`/`$env:GOMODCACHE`**，否则后续 `go test` 会重建该目录。
17. **报告写「已停止 / Stopped」会触发 `G7-d` 停止态判据**：⑩ 停点措辞应写「⑩ 停点已到达 / Reached the 10 stop point」，只有真正停止的切片才写「已停止——未完成」。
18. **捕获件的 `COMMAND:` 行必须与真实 argv 同形**：显示缩写（漏 `--noop-chain` 等必要旗标）会使该捕获「按记录不可重放」，属保真缺陷；⑤ 会以 Medium 判出。修法：`anomalies.md` 立 E-N 勘误 + 派生 `index.json` 写忠实命令与 `commandErrata`。
19. **远端 Go 模块获取**：VM 到 `proxy.golang.org` / `sum.golang.org` 可能是 IPv6 黑洞且 `GOPROXY=direct` 亦因虚名解析失败 ⇒ 用可达镜像预取后**构建期设 `GOPROXY=off`**（既解决下载又反证离线可构建）；模块缓存只读，删除前须 `chmod -R u+w`。

## 十、给新会话的一句话开场

以下是我的项目状态快照（v45），请先核第一节 5 个锚点，相符则按第四节队列继续，不符以仓库实况为准。

【粘贴本清单全文 v45】

【补充文档】`REMAINING_SLICES.md`（P55 已完成）、`DELIVERY_DIRECTIVE.md`（v1.60）、`docs/validation/exit-seed-build-evidence{,_EN}.md`（P55 报告，`63e2e24`，run 见 §二）。

当前片：P55 已闭环；下一片 = 层 3 第 8 项 `AT23-RESTRICTED-CANDIDATE`，**⛔ 阻断（等 R-08 / R-09 裁定）**。

## 十一、当前状态一句话

P55 `EXIT-SEED-BUILD-EVIDENCE` 闭环（`63e2e24` + `c76c6ce` 已推 origin，run `41925623071` 与 `41928327517` 双绿；**AT-13 双侧证据链 + 三角色两代构建落地**——seed `fc99f55` / verifier `18c9340` / candidate `baa144d`，`OracleHash` 四面一致，四条反例双侧 fail-close，失败不损 seed，R=0、D=1）；下一步 = 层 3 第 8 项 `AT23-RESTRICTED-CANDIDATE`，**当前 ⛔ 阻断，等 R-08 / R-09 裁定**。

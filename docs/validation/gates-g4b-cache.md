# G4b-CACHE 切片验证报告（业务片：G4b 跨进程语料缓存）

> 切片：`P27 / G4b-CACHE`（准则 v1.49，非治理片，M 档）。日期：2026-10-25。基线 `HEAD` = `7b84528`（P26 交付后）。变更落点见 §2，审查链见 §3，门禁读数见 §4。

## 1 切片概览

| 字段 | 内容 |
|---|---|
| 标识 | `P27 / G4b-CACHE` |
| 日期 | 2026-10-25 |
| 规模与分级 | M（业务片，性能优化，非治理片） |
| 基线与准则版本 | `HEAD` = `7b84528`（P26 交付后，准则 v1.48），本片落地后为 v1.49 |
| 目标 | 给 G4b 的「基线语料 CJK 字符集」加**跨进程缓存**：`--scope=tree` 整片从冷跑中位 **68539 ms** 降至热跑中位 **7226 ms** |
| 边界 | 只缓存语料字符集；判定结论、新增行字符与白名单一律不缓存；不改任何判据的判定逻辑 |

**设计五原则**（① 裁定方案，② 实现，③ 与 ⑥ 各自独立核验）：

1. 只存语料字符集：条目值恒为 `{v:1,corpus:string[]}`，不含任何判定信息。
2. 绝不缓存判定结论：`out.rec('G4b', …)` 每次无条件执行，命中也必须重算。
3. 键六段：版本、解析 sha、`ctx.kind`、白名单源、白名单摘要、判据源码摘要。
4. 读前重算键：键不单独持久化，命中判断前必须重新计算。
5. 任何异常 fail-open：建目录、读取、解析、写入与 `rev-parse` 失败一律静默降级为全量扫描。

## 2 变更落点

| 文件 | 规模 | 作用 |
|---|---|---|
| `tools/gates/lib/criteria/g4b.js` | **280 物理行**，numstat `241/34` | 缓存实现：键、存储、fail-open 与注释 |
| `tools/gates/selftest.js` | 新增 T108–T113，合计 **23** 条（674 → **697**），numstat `354/0`（删除行 = 0） | 缓存键、fail-open、索引脏与承重性变异 |
| `docs/DELIVERY_DIRECTIVE.md` | **2034 物理行**，numstat `18/2` | 头部 v1.49、§0.3 行、§6.3 的 G4 缓存条款与**升版义务**、§12.4 的 **OB-120**、附录 **C.0aj** |
| `docs/DELIVERY_DIRECTIVE_EN.md` | **2034 物理行**，numstat `18/2`（与 CN 对称） | 镜像侧，逐句对齐 |
| `docs/t027/REMAINING_SLICES.md` | numstat `21/0` | P27 那一行与详情块，含「观察项 / 下片建议」与「已知限制」两行 |
| `docs/t027/REMAINING_SLICES_EN.md` | numstat `21/0`（与 CN 对称） | 镜像侧 |
| `tools/gates/cjk-newwords.txt` | 新增 **6** 条登记（**37** → **43**） | 逐字转录与撤回措辞所必需汉字的白名单条目（每条带理由与上下文摘录） |

补充说明（落盘形态与顺序自查）：

- **① (a) 落盘形态**：缓存跨进程落盘，路径 = `repoRoot/tmp/gates-cache/g4b/`；不放在 `tmp/G4b-CACHE/` 的理由 = 后者是切片临时件、⑧c 会清零，清零即毁缓存，与 C.0aj「缓存目录为运行时状态」的自洽要求冲突。
- **③ 顺序自查**：③ = `L68115`、⑤ = `L69717` ⇒ ③ < ⑤ 成立；④ 是主控环节（不在派单表），按 §5.3 在 ③ 之后执行。
- **⑥ 先于 ⑦**：⑥ r3 = `L71059`、⑦ r1 = `L71123` ⇒ ⑥ < ⑦ 成立。
- **R7.1 逐条计量**：逐条见 §3.1 派单表（26 条）；`retryOf` 语义 = ②-2…②-7 分别对前次 ② 的整改重跑；⑥ r2/r3 分别对 r1 的格式重发与整改重跑；⑦ r2 对 r1 的整改复审。

## 3 审查链

### 3.1 派单表（机械读取）

| 阶段 | 角色（`agentName`） | 载体模型 | 额度 | 结论 |
|---|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 3 轮 | 初版、增量复核与 ⑦-M 裁决，均 `DESIGN: DONE` |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 6 轮 | 实现、条件③移除、注释精确化、四处回写、§12.4 的 OB-120、恢复措辞又撤回 |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 3 轮 | 三轮均 `TEST: DONE`（Medium+ 归零） |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 1+1 用尽 | r1 `PRE-REVIEW: PASS WITH FIXES`，r2 同 |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 3 用尽 | r1 输出不合规作废，r2 `FINDINGS`，r3 `PASS` |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 1+2 | r1 `FINDINGS`（1 High），r2 `RE-REVIEW: PASS` |
| 批量协助 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 3 次 | 枚举、关键事实与测试草稿 |

- 批量协助的三次调用只做只读整理，不产生判定，也不改变任何额度计数；⑥ 的 r1 因输出契约不合规而作废（见 §3.5），既不计审查价值也照扣额度。

### 3.2 ① 四项关键判定

| 判定项 | 结论 |
|---|---|
| 缓存值纯度 | 只存 `{v:1,corpus:string[]}`，verdict 每次重算（`g4b.js` 的 `out.rec('G4b', …)` 无条件执行） |
| 键完备性 | 三反例分别由 sha 段、白名单摘要段与 `sourceDigest()` 段承担；`ctx.kind` 为**承重段**（range 用 `ls-tree`，tree 与 index 用 `ls-files`） |
| fail-open | `mkdirSync`、`get`、`parse`、`put` 与 `rev-parse` 全部静默降级；唯一 `UsageError` 是缓存之前的 `whitelist.error` |
| **条件 ③ 移除** | 初版设计要求「tree 与 index 索引洁净才启用缓存」，① 在增量复核中认定其论证有误（语料只依赖 HEAD 对象库）；移除后 `ctx.kind` 升格为承重段，并由 **T113 与变异**钉住不得复活 |

### 3.3 ③ 三轮独立核验

| 轮次 | 结论 | 要点 |
|---|---|---|
| r1 | `TEST: DONE`（Medium+ = 0） | 断言计数 674 → 691；三反例与 fail-open 由独立探针复现；每条键要素经变异证明承重 |
| r2 | `TEST: DONE`（Medium+ = 0） | 691 → 697 的前序回写核验；生产文件多轮变异后 `sha256` 逐位复原；tree 三组耗时与 tree↔ci 归因 |
| r3 | `TEST: DONE`（Medium+ = 0） | 补 `sourceDigest` 前后同为 `10d7ec758d494d7570d7e98e227802c5732ef6db0a59b6f604f02dbf46b5e21f`（`run.toString()` 与 `cjkChars.toString()` 拆分 sha256 逐字节相等）；真实 git 端到端六步；tree 组中位复核 |

三次核验的共同要点：断言计数 674 → 691 → **697**（净增 **23** 条）；三反例（HEAD 改动、白名单改动、判据源码改动）与 fail-open 由**独立探针**复现；每条键要素经变异证明**承重**；生产文件经多轮变异后 `sha256` 逐位复原；真实 git 端到端的「干净写入 → 变脏读取仍命中」六步落在临时仓 `repo-clean`（`git write-tree` 前后相等 `c9bfc762…`）。

### 3.4 ⑤ 预审（额度 1+1，两轮用尽）

- **r1 `PRE-REVIEW: PASS WITH FIXES`**：1 Medium（`G4B-CORPUS-DECOUPLE` 观察项与下片建议未在**任何**已暂存文档登记）+ 3 Low（T108–T112 未含 T113；§6.3 缺升版义务；`keyString` 二次 `rev-parse` 的 `'none'` 键 TOCTOU）+ 3 Info。
- **r2 `PRE-REVIEW: PASS WITH FIXES`**：r1 的 1M 与 3L 全部闭环；新增 1 Low（P27「状态」行仍写「① 待派」）+ 3 Info（§6.3 升版义务无机械强制；#4 与 #5 待报告落实；编辑器缓冲卫生）；⑤ 明确裁定主控「恢复『彻底』并登记白名单」的处置**正确**。

### 3.5 ⑥ 独立扫描（三轮，含逐字转录）

- r1 结论为 `PASS`，但**输出契约不合规**（`## Section A — 结论` 形态、Section D 非可证伪性审计、末行非结论行）⇒ 作废。
- r2 `FINDINGS`：1 Medium（缓存键未机械绑定 `ctx.listAllMarkdown` 与 `ctx.showAt` 的语义）。
- r3 `PASS`：该 Medium 经 ① 裁决与 ⑥ 独立复核后**撤回**，降为已知限制。

以下为 ⑥ r3 的完整输出**逐字转录**（未作任何改写）；围栏内两处结论行（Section A 下第一行与块末）取值相等、位置合规，符合 `§6.3` 的 `G7-f` 口径。

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
无

## Section C
- 命令 1: cd /d/LZProjects/prfrail; git rev-parse --short HEAD; git status --porcelain -uall
  关键输出摘要: HEAD=7b84528；状态为 6 个 M、0 个 ??，没有未跟踪残留。
- 命令 2: git log --oneline -- tools/gates/lib/ctx.js
  关键输出摘要: 仅 2 条提交；即 ctx.js 的语义扩展历史极短。
- 命令 3: git show --stat --oneline 8a2841f -- tools/gates/lib/ctx.js tools/gates/lib/criteria/g4b.js tools/gates/selftest.js
  关键输出摘要: 8a2841f 同时改了 3 个文件，包含 ctx.js、g4b.js、selftest.js。
- 命令 4: git show --stat --oneline 025158f -- tools/gates/lib/ctx.js tools/gates/lib/criteria/g4b.js tools/gates/selftest.js
  关键输出摘要: 025158f 同时改了 ctx.js 和 selftest.js；两次历史都未显示“只改 ctx.js 且不联动 g4b.js/selftest.js”。
- 命令 5: git grep -nE 'actions/cache|gates-cache|tmp/' -- .github/workflows/ci.yml
  关键输出摘要: 无命中；CI 工作流未配置 actions/cache，也未使用 gates-cache 或 tmp/ artifact 复用。
- 命令 6: node tools/gates/gate.js --selftest
  关键输出摘要: SELFTEST: PASS 697/697；TOTAL_FAIL=0；G4b unknown=[]；whitelist=37；退出码 0。
- 命令 7: 读取 .github/workflows/ci.yml 的 Gates 步骤和 checkout 片段
  关键输出摘要: `Checkout source` 配置 fetch-depth: 0；Gates 步骤位于全新 checkout 之后；与用户给出的主控事实一致，CI 是冷跑，不存在本地 tmp/gates-cache 复用路径。

## Section D
(i) 失效方向：部分成立。它的“危险方向”确实是“旧缓存命中导致少报（false negative）”，但它并不能越过 CI；这是 ① 的核心判断。这里的关键不是“可能有少报”，而是“少报是否能通过冻结提交门禁”。在真实工作流里，CI 每次都是全新 checkout，且 `.github/workflows/ci.yml` 中没有 actions/cache、没有 gates-cache、没有 tmp/ artifact 下载；因此本地旧条目无法成为 CI 证据源。结论：该论点只证明“本地预检可能瞬时不一致”，不能证明“提交门禁会被绕过”。
(ii) 概率：成立。三条件连乘确实极低：ctx 语义变更 × 未升 CACHE_VERSION × 本地存在旧条目；而且 git 历史非常短：ctx.js 只有 2 条提交，且两个关键变更都同步改了 g4b.js/selftest.js。`keyString` 也把 resolved sha、whitelist 摘要、sourceDigest、CACHE_VERSION 一起绑定，且 `selftest.js` 的 T109–T111 已明确验证 HEAD 改动、白名单改动、判据源码改动都必须 miss。于是“语义漂移而未更新缓存键”的可发生性被强力压低。
(iii) 可检测性：部分成立。这里确实没有一个“在语义漂移时自动变红”的断言，T108–T113 主要覆盖的是缓存键、fail-open、索引脏、不变更/变异手段，而不是“真实 ctx.js 改了但没升级缓存键”的现场监测。这个事实是真，但它不等于“能绕过 CI”，因为结构性冷跑已经把风险边界压在本地、而不是提交门禁。也就是说，检测性缺口存在，但它只影响本地快速反馈，不改变 CI 的有效性。
(iv) 现行义务：成立。`docs/DELIVERY_DIRECTIVE.md` 中已经写明“升版义务”：ctx.listAllMarkdown / ctx.showAt / git 枚举或读取语义变化时，必须同步升级缓存格式版本前缀，否则旧条目一律视为陈旧。并且 `OB-120` 作为候选修法已登记，说明本问题已被纳入治理项，非无视。这里真正的问题不是“没有规则”，而是“该风险在当前工作流中不构成 Medium 级别的 CI 间接威胁”。
结论：本题不成立为 Medium；至少应归为 Low/已知限制，并且更稳妥的裁决是撤回 Medium。没有发现一个“能同时骗过本地命中与 CI 冷跑”的可发生反例序列；这类风险在当前仓库结构下只剩本地预检瞬时不一致，不足以支撑 Medium 级别结论。

## Section E
- r2 的格式问题：已修复。本轮严格按照模板输出，且最终行与 Section A 完全一致。
- 无法核验项：没有对“未来某次 ctx.js 语义改动 + 未升 CACHE_VERSION + 本地旧缓存被复用”的特定历史分支做真实演练；若存在该类实测记录，结论会被推翻。
- 可能推翻 PASS 的观测：若未来实测显示 CI 中确实存在 persistent cache 复用，或出现“ctx.js 改动未同步升级 g4b.js / selftest.js 但能通过新 checkout 结论”，则 Medium 可能重新成立。
INDEPENDENT SCAN: PASS
```

### 3.6 ⑦ 终审（盲审 r1 + 复审 r2）

- r1 的结论行为 `FINAL-REVIEW: FINDINGS`：1 **High**（`SELFTEST: FAIL 695/697`，T87 的两条断言 `unresolvable`）。
- 主控按 **§9.1** 做三步根因仲裁：清场前 `CAND_COUNT=0`（唯一 node 为 PID 4820，即 MCP 宿主）；**单跑 7 次全 `PASS 697/697`**（`EXIT=0`、`FAIL_LINES=0`）；**并发 A/B 均 `PASS 697/697`**；`resolveRepoArtifact` 两分支均 `ok:true`；`git ls-files --error-unmatch -- tmp/.gitkeep` 退出 0。
- 判定：关于「本片代码缺陷」与「selftest 不耐并发」的两个假设**均被证伪**。
- 根因候选 **A（一次性并发，最可能为 ⑥ r2 自报的残留 selftest 进程；PID 未记录、存活窗口不可回溯，`Kernel-Process/Analytic` 本机无记录）** 与 **B（`spawnGit` 一次性 spawn 失败）** 并列，均未 100% 机械坐实。
- **关键限定**：`spawnGit` 把 stderr 收进结果对象而不打印，因此 `codex-st.txt` 的「无 git 错误噪声」**不能**区分「git 报错」与「spawn 失败」。
- r2 的结论行为 `RE-REVIEW: PASS`：确认该 High 不成立为代码缺陷，7 次单跑与并发 A/B 的同位断言均为 PASS，`res-probe` 两分支 `ok:true`；并指出本片改动集中于 G4b 缓存与 T108–T113，未触及 T87 与 G9 的逻辑。

## 4 门禁读数

### 4.1 定稿读数（六个文件）

| 作用域 | EXIT | 末行 | G1-b | G3(a) | G4a | G4b | G6-4/5/6/7 | G9 |
|---|---|---|---|---|---|---|---|---|
| `--all --scope=index` | 0 | `TOTAL_FAIL=0` | OK | OK | OK | OK | OK | OK |
| `--all --scope=tree` | 0 | `TOTAL_FAIL=0` | OK | OK | OK | OK | OK | OK |
| `--all --scope=ci` | 0 | `TOTAL_FAIL=0` | OK | OK | OK | OK | OK | OK |
| `--selftest` | 0 | `SELFTEST: PASS 697/697` | — | — | — | — | — | — |

判据读数：`G1-b` 为 `header=v1.49 last=v1.49 rows=50 OK`；`G3(a)` 准则对 `18/2` 两侧对称、台账对 `21/0` 两侧对称；`G4a` 为 `6 file(s) OK`；`G4b` 为 `newCjk=401 unknown=[]`（whitelist=37，即本报告落盘前的白名单状态；index 与 tree 腿），ci 腿为 `newCjk=148 unknown=[]`（`endpoint:HEAD`，ci 腿只覆盖 `HEAD..base`）；`G9` 为 `6 path(s) OK`。另：`SELFTEST: PASS 697/697`（7 次单跑与 2 次并发全绿）。

### 4.2 含本报告后的定稿读数

本报告与自由度清单落地后，变更集为 **10 路径**（7 个 `M` 与 3 个 `A`，未跟踪 `??` = 0）；⑨ 定稿读数：

| 作用域 | EXIT | 末行 | changed | G4a | G4b | G9 |
|---|---|---|---|---|---|---|
| `index` | 0 | `TOTAL_FAIL=0` | 10 | `10 file(s) OK`（1 excluded: evidence tree） | `newCjk=658 unknown=[]`（`whitelist=43`） | `10 path(s) OK`（1 excluded） |
| `tree` | 0 | `TOTAL_FAIL=0` | 10 | 同上 | 同上 | 同上 |
| `ci` | 0 | `TOTAL_FAIL=0` | 4 | `4 file(s) OK` | `newCjk=148 unknown=[]`（`whitelist=37`，`endpoint:HEAD`） | `4 path(s) OK` |

- `SELFTEST: PASS 697/697`（EXIT 0，`^FAIL` = 0）。
- 说明：本表为**收尾时刻**读数；`tmp/gates-cache/g4b/` 的缓存条目数随运行增减（收尾时刻 9），不影响任何判据读数。
- `G7-f`：`2 changed sixstage fence(s) complete (A→B→C→D→E, closed-set conclusion, SEVERITY consistency)`。
- `G9-e` 对本报告引用**未暂存**的运行日志登记 INFO（informational、永不判红）；该行为属设计——临时件不入库，引用它们只被「提示」而不被「判红」。

### 4.3 `SELFTEST` 计数变化

`674` → `697`（净增 **23** 条）：T108 3 条、T109 3 条、T110 3 条、T111 3 条、T112 5 条、T113 6 条。附 `git diff --numstat HEAD -- tools/gates/selftest.js` 为 `354 0`（删除行 0）。

## 5 低级别发现处置清单

| # | 来源 | 级别 | 位置 | 处置 |
|---|---|---|---|---|
| 1 | ① | — | `g4b.js:33` 为 **83 列**（① 主张逐字文本优先） | 接受列宽超限：替换前该处原单行注释为 87 列（净改善），`g4b.js:137` 的 81 列与文件既有最大值持平，因此**未抬高**文件既有列宽上限 |
| 2 | ② | — | 行数口径 | 已处置：② 曾把行数报成 **281 物理行**，实为 **280 物理行**（`split('\n').length` 差一） |
| 3 | ③ | Low | `({)` 零字节未跟踪件 | 已处置（临时件纪律，自愈档）：③ 发现，经授权后 `Remove-Item -LiteralPath '({)'` 清除，`git status --porcelain -uall` 的 `??` 归零，`--scope=index` 由 `TOTAL_FAIL=1` 转 `TOTAL_FAIL=0`。来源：mtime 05:22:34（本片窗口内），全量 transcript 扫描该窗口内没有任何记录的 `> ({` 重定向，邻近匹配均为 JS 的 `=> ({` 误命中，因此无法定位确切命令（如实登记，不编造）；机制族为 PowerShell 把表达式片段当作重定向目标。同时构成 G5(f) 的正面例证：该件使 `--scope=index` 报 `TOTAL_FAIL=1`，证明该判据真能抓到 stray 未跟踪件 |
| 4 | ⑤ r1 | Low | `keyString` 二次 `rev-parse` 的 `'none'` 键 TOCTOU，以及 T112 五面未直接覆盖 `mkdirSync` 失败、`put` 抛错与 JSON 语法非法 | 记录在案：本片不修，已在台账「已知限制」行登记，理由为方向 fail-safe 且与已测路径同构 |
| 5 | 本片新增 | Low | `selftest.js` 的 `spawnGit` 静默 stderr | 观察中：候选修法为 spawn 失败时把 stderr 尾部附进 `detail`（不改变判定，只增诊断）；触发条件为下一次任何 git 相关断言失败且原因不明时由 ① 评估；本片不改 |

## 6 边界与不追溯

- `git diff HEAD -- lib/ctx.js tools/gates/gate.js .gitignore internal .github docs/DEV_PLAN.md` 的输出为空。
- 未改 §2.2、§2.4、§5.2、§7.5、其他判据本体、CONTRACTS、schema、fixtures 与 workflows。
- `tools/gates/cjk-newwords.txt` 相对 `HEAD` 新增 6 条登记（**37** → **43**，每条带理由与上下文摘录）；编码门禁：`.md` 为 BOM 加 LF，`.js` 为无 BOM 加 LF。

## 7 遗留与后续

### 7.1 已知限制（逐字采用）

> ⑦ r1 独立报 High（`SELFTEST FAIL 695/697`，T87 两条 `unresolvable`）。经主控三步根因仲裁：清场前无 gate/selftest 进程；单跑 7 次均 `PASS 697/697`；并发 A/B 均 `PASS 697/697`；`resolveRepoArtifact` 两分支与 `git ls-files --error-unmatch -- tmp/.gitkeep` 均正常。⇒ 「代码缺陷」与「selftest 本身不耐并发」两个假设均被证伪。根因候选至少两个，均未 100% 机械坐实：(A) ⑦ 执行时刻的一次性并发，最可能为 ⑥ r2 自报的残留 selftest 进程——该进程 PID 未记录、存活窗口不可回溯（`Kernel-Process/Analytic` 无记录）；(B) `spawnGit` 在该次运行中一次性 spawn 失败——`codex-st.txt` 无 git 错误噪声，但 `spawnGit` 把 stderr 收进结果对象而不打印，故「无噪声」无法区分「git 报错」与「spawn 失败」。⑦ 的读数真实，归因修正为环境瞬时（候选 A 或 B，未定）。防护措施 = 本片不加锁（并发 A/B 不复现，且加锁仅覆盖候选 A；若真因是 B 则加锁无效），登记观察项：若再现且清场确认无残留进程 ⇒ 另立片给 selftest 入口加进程级互斥（候选 A）或给 `spawnGit` 加诊断与重试（候选 B）。⑦ 同时存在两处流程偏离（§7.3② 执行命令、§7.3③ 创建 `tmp/G4b-CACHE/codex-st.txt`），已登记，不改变读数有效性。

### 7.2 设计偏离（四要素式）

- **D-1（① 自认）**：初版设计所列启用条件 ③（tree 与 index 仅当 `git diff --cached --quiet HEAD` 退出 0 时启用缓存）论证不完备；正确的不变量是「缓存语料恒为 HEAD 真值语料的子集、命中方向天然 fail-safe」，与索引洁净无关。该条已移除（`ctx.kind` 升格为承重键段），T113 与变异钉住其不得复活。
- **D-2（① 自认）**：① 交付的 `cacheUsable` docblock 曾写「the full-scan path computes exactly the same index-dependent value」，该等式只在索引状态自写入起未变时成立；暂存删除 `.md` 后冷跑语料收缩，与缓存值不再相等（③ 实测 `HIT_VS_COLD_SAME_VERDICT=false`）。已改写为「缓存值是写入时索引状态下的 HEAD 语料、恒为真值语料子集、后续收缩只会让冷跑更严格、命中方向 fail-safe」。
- **D-3（执行自由度扩张，已撤回）**：回写时主控自行把 §12.4 的 OB-120 措辞「完全解耦」改为「彻底解耦」，并新增白名单 `彻` 条目；机械查询显示「完」与「全」均在 HEAD 语料中（含「完」的文件为 165 个），故「完全」本就合法，该改动属**反向操作**。已撤回措辞改动与白名单条目，白名单回到 **37** 条（本报告另行登记的 6 条见 §2 与 §6，与 D-3 的撤回无关）。教训：执行自由度扩张（不在审查方发现范围内的措辞改动）必须先做机械判定，不得凭直觉优化；「自主优化」在治理严谨的项目里是风险源。
- **D-4（临时件纪律，自愈档）**：见 §5 第 3 条（`({)`）。
- **D-5（单方声称标注）**：② 曾报 `--selftest` 的 `EXIT=0` 系「由代码路径推得」，非同一轮捕获；③ 随后两次同一轮捕获 `EXIT=0`，该疑点归零。
- **D-6（审查方流程偏离，⑦ r1）**：⑦ r1 执行了 `node … --selftest`（§7.3②）并创建 `tmp/G4b-CACHE/codex-st.txt`（§7.3③）；该文件时间戳 08:28:11–08:34:49（约 6.5 分钟，等于一次完整 selftest）可作佐证。**不改变读数有效性**；若 ⑦ r2 再犯则按 §7.3 升级。r2 未再犯。
- **D-7（同轮临时授权例外）**：§5.1 的 ⑩ 前终态回读 ② ③ 在本片不通过（`tmp/` 递归文件数 > 1、顶层含 `gates-cache`）。处置 = 用户 2026-10-25 同轮授权例外（四要素见 §7.5），**非规则违反**；`tmp/gates-cache/g4b/` 的运行时缓存条目保留（数量随运行增减，收尾时刻为 9）、⑧c 不清零；本条不得作先例。

### 7.3 观察项与下片建议

- **观察项**：`listAllMarkdown` 对 tree 与 index 用 `git ls-files`（读索引），因此暂存删除 `.md` 时冷跑语料相对 HEAD 真值收缩，冷跑从严多报；在「写入后暂存删除 `.md`」的窗口内，缓存命中与冷跑可能给出不同 verdict（实测 warm=PASS、cold=FAIL），方向 fail-safe（缓存值恒为真值语料子集、从不洗白违规）。**候选 OB `G4B-CORPUS-DECOUPLE`**：影响面已核实——`listAllMarkdown` **仅 G4b 消费**（`g4b.js:223`），G4a 用 `changedPaths` 加 `readBytes`、G7 用 `changedPaths` 与 `exists`、G8 用 `listTreeFiles`（已改用 `ls-tree <headRef>`）、G9 与 G1 与 G5 用 `changedMarkdown` 等。本片不执行，已在 §12.4 登记 OB-120（首状态词为 `待议`）。
- **候选 OB（性能）**：ci 作用域约 **24–37 s** 的冻结成本来源未定，非 G4b 缓存可解（`--check=G4b` 单判据 ci 冷 59435 ms、热 3090 ms，降幅 94.80%；整片 ci 冷 91579 ms、热 27428–39727 ms；`gate.js:147` 为 `await Promise.all(jobs)`，因此整片是并发关键路径）。
- **候选 OB（诊断）**：`spawnGit` 静默 stderr（见 §5 第 5 条）。
- **候选 OB：C.0aj ↔ §5.1 规则冲突**——C.0aj 允许缓存目录放 `tmp/` 下保留，而 §5.1 的 ⑩ 前终态回读要求 `tmp/` 递归文件数 = 1 且顶层名恒为 `.gitkeep`，两者无法同时成立。候选修法：(i) 缓存目录统一移出 `tmp/`（如仓库根 `.gates-cache/`），§5.1 加指针说明；(ii) §5.1 的 ② ③ 增加「运行时缓存目录例外」条款。本片不修，交下一治理片。

### 7.4 性能三组（中位，`--scope=tree`）

| 组 | 3 次原始毫秒 | 中位 | 日志行 |
|---|---|---|---|
| A 未命中（每次清缓存） | 68539 / 66608 / 75815 | **68539** | `miss 88364fcfdfda corpus=1333` |
| B 命中（不清缓存） | 7232 / 7226 / 6808 | **7226** | `hit 88364fcfdfda corpus=1333` |
| C 失效后重建（清空后跑） | 81089 / 72017 / 63941 | **72017** | `miss 88364fcfdfda corpus=1333` |

- 命中与未命中之比为 **10.54%**，即压降 **89.46%**（门槛 ≥30%）；C 组接近 A 组（差 5%，落在冷跑自身变化范围内）；口径＝冷跑前清空 `tmp/gates-cache`、至少 3 次取中位、命中仍需重算键与 verdict（非零成本）、数据取自本机 Windows、随 `.md` 语料规模缩放。
- 固定开销参照（`--check=G1-a --scope=tree`）为 1846 ms；扣除后 `(7226−1846)/(68539−1846)` = **8.07%**，即压降 **91.9%**。

### 7.5 同轮临时授权例外（§1.7.4 四要素）

| 要素 | 内容 |
|---|---|
| 裁决人 | 用户 |
| 日期 | 2026-10-25 |
| 理由 | C.0aj（缓存目录为运行时状态，按运行期状态保留）与 §5.1 的 ⑩ 前终态回读（`tmp/` 递归文件数 = 1、顶层名恒为 `.gitkeep`）存在规则内部冲突；本片不修规则、不改代码，选临时授权例外 |
| 依据 | 本轮用户回复原文「选 D」 |

- 性质：超出 §5.1 的 ⑩ 前终态回读 ② ③ 字面标准，系用户直接裁定。
- 范围：仅 `tmp/gates-cache/g4b/` 的运行时缓存条目保留（数量随运行增减，收尾时刻为 9）；§5.1 ② ③ 对本片不适用。
- 不得作先例；规则冲突另立候选 OB（见 §7.3）交后续治理片。

## 8 提交与 CI

本报告与自由度清单为 ⑧a 的落盘物；⑨ 重跑三作用域与 `--selftest` 后把定稿读数补记于 §4.2。提交与 CI 证据如下（回写提交之后的 run 只作观察记录，按 §11.2 不再为其新增提交）。

| 项 | 值 |
|---|---|
| 提交 | `64d1d58f82446bd9c1675df23523411785f5f29d` |
| 推送 | `origin` 的 `7b84528..64d1d58`（不推 gitee） |
| CI run | `40183876643`（`https://github.com/larsonzh/prfrail/actions/runs/40183876643`） |
| 结论 | `completed / success` |

- 双腿：`Go windows-latest` = `success`、`Go ubuntu-latest` = `success`；其余 3 个 job `skipped`（候选与发布链未触发，属预期）。
- `Gates` 步骤命令（两腿相同，仓内相对路径）：`node gate.js --selftest && node gate.js --all --scope=ci`。
- Windows 腿：`SELFTEST: PASS 697/697`、`TOTAL_FAIL=0`、`ok github.com/larsonzh/prfrail/internal/gates`（`0.244s`）。
- Ubuntu 腿：`SELFTEST: PASS 697/697`、`TOTAL_FAIL=0`、`ok github.com/larsonzh/prfrail/internal/gates`（`0.140s`）。
- 独立佐证：全新 checkout、无本地缓存的 CI 环境下，含 T87 与 T108–T113 的 697 条断言在 Windows 与 Ubuntu 双平台均全绿 ⇒ 为 ⑦ r1 的 T87 疑点提供独立反证。

```artifacts
tools/gates/lib/criteria/g4b.js	repo
tools/gates/selftest.js	repo
tools/gates/cjk-newwords.txt	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/gates-g4b-cache.md	repo
docs/validation/gates-g4b-cache_EN.md	repo
docs/validation/evidence/G4B-CACHE-freedom-list.md	repo
```

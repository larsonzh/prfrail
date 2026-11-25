# P88 ① 架构方案书（逐字转录；载体 deep-reasoner + DeepSeek V4 Pro）

> 说明：本文件是 ① 的原始输出逐字转录（仅补文件名标题行），未做内容改写。

# P88 ① 影响判定与设计方案（结论先行）

## 结论速览

| # | 裁决 | 一句话 |
|---|---|---|
| ① 影响判定 | 不触 CONTRACTS / schema / fixtures；落在「判据实现文件」**成立**；**必须升 MINOR v1.72→v1.73**（须用户同轮授权，§1.7.2） | 依据 = §0.1 权威表 + v1.61 先例 |
| O-1 | **做**：ctx 层 `cat-file --batch` 批量 blob 预取（g4b.js 零改动、CACHE_VERSION 保持 v2） | 冷扫 490s → ~100s 级，CI 每次 run 直接受益 ~300–400s |
| O-2 | **不做**：base 四件套并行 | 冷缓存下有变慢风险、CI 收益≈0、抬高 flake 面 |
| O-3 | O-3c ctx 读数 memo **做**（~10–20s）；O-3a g4b.js 内部去重 **不做**（key12 刚性 + 收益≈0）；「loadWhitelist 重复读」**前提不成立**（检索证伪） | — |
| A/B | 冷-冷 + 热-热配对；输出**不含耗时数字 ⇒ 零归一化**；热-热配对直接证明 cache 行逐字一致 | 旧实现快照 `tmp/P88-GATES-PERF-ASSESS/orig-gates/` 已存在 |
| selftest | 既有 744 全部保持通过；新增断言致总数变化 ⇒ **须用户授权口径**（推荐按 v1.40/v1.61 先例改基数） | — |
| 升版授权 | v1.73 的落版是 §1.7.2 升级边界，**必须用户同轮授权后才能执行文档回写** | 本报告即授权请求 |

---

## 1. 影响判定

### 1.1 是否触 CONTRACTS / schema / fixtures —— **否**（检索证据）

- **§0.1 权威关系表**（`docs/DELIVERY_DIRECTIVE.md` 实测读取）：「门禁判据实现与真值源 | `tools/gates/**`（判据文本见 §6.3）」——本片全部代码落点（`tools/gates/lib/ctx.js`、`tools/gates/selftest.js`，gate.js 不动）都在 `tools/gates/**` 内，不落在「协议与语义 = CONTRACTS / schemas / 契约夹具」权威行。
- **§1.4 协议先行不触发**：该规则适用于「语义变化」；本片验收硬条件①（三作用域输出逐字一致）本身就是「无语义变化」的机器证明，故 CONTRACTS 双语 → schema → 代码的顺序链条无需启动。
- **fixtures**：`tools/gates/testdata/` 与 `schemas/`、`testdata/` 均不改；`selftest.js` 是自测代码而非契约夹具。
- **CI workflow**：`.github/workflows/ci.yml` 不改（边界明确）。

### 1.2 「判据实现文件」是否成立 —— **成立**

- §0.1 把 `tools/gates/**` 定义为判据实现与真值源；`ctx.js` 是全部判据共享的读取/枚举层，改动它 = 改动判据执行路径（正是起跑消息用户裁定 b 的措辞：「触门禁判据执行路径（tools/gates 实现，不触判据本体）」）。
- 后果链：§7.2.1 逐条——① 否（脚本逻辑变更）、② 否（脚本逻辑变更）、③ 判据实现层变更 ⇒ **⑥ MAI 必做**（与台账步骤、起跑消息一致）；⑦ 档位 = Codex 常规档（用户裁定 b，与 §5.2 判据一致，无需再判定）。

### 1.3 是否必须升 MINOR v1.72→v1.73 —— **是（裁决 + 依据）**

1. **先例（决定性）**：§0.3 v1.61 `G4B-CORPUS-DECOUPLE`——纯 `tools/gates` 实现变更（ctx.js + g4b.js + selftest 697→703），明确写「**不新增判据、不改 §6.3 判据文本本体**」，**仍升 v1.61**。本片与之同型（判据优化片、实现层改动、判据文本零改动）。
2. **权威定位**：§0.1 定位 `tools/gates/**` 为真值源 ⇒ 真值源变更必须在 §0.3 变更日志留痕；新增 §0.3 行受 **G6-5**（新增行节号引用解析）与 **G6-7**（声明行 ↔ 台账 ID 对账）机械判据约束。
3. **版本元数据一致性（G1-b）**：头部版本必须等于 §0.3 最新行版本 ⇒ 加行必同步头部 v1.73；历史无「两行同版本」先例，不加行则属隐瞒。
4. **MAJOR 排除**：§0.2「MAJOR = 角色集或流水线结构变更」不成立 ⇒ 不升 MAJOR（台账边界亦禁）。

**升版动作本身属 §1.7.2 升级边界**（「未经用户同轮明确授权的准则升版」）⇒ ① 报告提交此裁决，**用户同轮授权后 ② 才执行**头部/§0.3/附录 C.0az 回写；不得在 ② 静默升版。§0.3 变更行与 C.0az 块（形状按 C.0aj）须载明：「不新增判据、不改 §6.3 判据文本本体、不改 G5-b 闭集与白名单、不触 CONTRACTS/schema/fixtures/workflows、不动 internal/** 与 CI workflow、不升 MAJOR」。

---

## 2. 优化落点方案

### O-1 ctx 层批量 blob 预取（核心，做）

**文件**：`tools/gates/lib/ctx.js`（g4b.js **零改动**）。

**改动点**：
- `showAt(ref, p)`：先查 per-context `Map`（键 = `ref + '\u0000' + p`，避免 `:` 歧义），命中即返。
- `listAllMarkdown()`：枚举后、返回前，若 `list.length ≥ BATCH_MIN`（建议 8；实测语料 1131 远大于阈值；selftest 既有夹具 ≤2 个文件 ⇒ 不触发，744 不受扰动）则同步调用内部 `prefetchBlobs(ref, list)`。
- `spawnGit` 微调一处（必要）：`encoding` 改为 `binary && input !== undefined ? undefined : (binary ? 'buffer' : 'utf8')`。**实测依据**：本机 Node v24.17.0 下 `encoding:'buffer'` + 字符串 input 抛 `ERR_UNKNOWN_ENCODING`，而 + Buffer input 通过；`encoding:undefined` + 字符串 input 通过（stdout 为 Buffer）。现有调用不受影响（`isIgnoredBatch` 是 binary=false；G4a 的 binary 读无 input）。

**请求/响应格式**（均已本机实测）：
- 请求：每行 `<ref>:<path>`，`\n` 结尾，整串 `Buffer.from(specs,'utf8')`；无 shell、无 core.quotePath 转义 ⇒ CJK/空格路径逐字直通（本机实测当前语料 0 个此类路径，合成夹具钉住）。
- 响应：字节级循环——`<oid> SP <type> SP <size> LF` 头行；`<spec> missing LF` ⇒ 该路径记 null（实测 `HEAD~3:nope/nope.md missing`、exit 0）；按 `size` 精确切片 blob、断言其后恰一个 `LF`。
- **二进制**：blob 原始字节 `buf.toString('utf8')` 与 `spawnSync encoding utf8` 的 `git show` 返回逐字节一致（实测 4 个含 BOM 文件 OK）。
- **decodeBuffer 语义**：`showAt` 现语义 = 原样 utf8 解码、**不**剥 BOM、**不**做 CRLF→LF ⇒ 预取缓存同样零归一化（与 showAt 对齐，而非 readText）。
- 路径含 `\n`：不参与预取（回退逐文件 showAt）——确定性且安全。

**fail-open 判定**（任一命中 ⇒ 整批丢弃、不部分采用、静默回退逐文件、无 throw、无任何新输出行）：spawnError / `!ok` / out 为空或不可解析 / 头行不匹配 / size 越界 / 尾部缺 LF / 应答条数 ≠ 请求条数 / 超 maxBuffer（200MB；实测语料 15.85MB，余量充足）。

**缓存生命周期**：per-context Map，随进程消亡；跨进程缓存仍由 g4b.js 既有 corpus 缓存承担，两层互不干扰。

**预计收益**：1131 次 spawn（实测单 spawn ≈137ms，另加每文件对象读取）→ 1 次 spawn + 15.85MB 解析（数百 ms 级）。`--all --skip-base --scope=range:HEAD~3..HEAD` 490s → ~100s 级；CI Gates 步每次 run 的 G4b 冷扫（HEAD^1 语料、tmp 不入库）直接受益。tree/index 热命中路径不触发任何 showAt（hit 分支不调 listAllMarkdown），零开销。

**关键论证——CACHE_VERSION 保持 v2（不升 v3）**：§6.3 G4b 缓存「升版义务」触发条件 =「`ctx.listAllMarkdown` / `ctx.showAt` 的枚举或读取**语义**变化」。本实现是**传输层重构**：`(ref, path)→文本` 的数学函数逐字节不变（A/B 字节 diff + S-1 等价断言 + T115/T116 冷热等价兜底）。若升版则 key12 必变 ⇒ INFO 行 `miss <v3…>` ≠ `miss <v2…>` ⇒ 与验收①「diff 为空」**不可兼得**。⇒ 本论证必须写入 ⑦ 输入包；若 ⑦ 采严格字面读法，O-1 放弃（仅留 O-3c），本片仍可闭环。**回退条件**：A/B 任一作用域 diff 非空 / 任一既有 744 断言红 ⇒ 整块回退（含 spawnGit 微调）。

### O-2 base 四件套并行 —— **裁决：本片不做，进大改建议**

依据（用户前提已复核，结论成立）：
1. **冷缓存下有变慢风险**：实测构成 gofmt 1.1 + build 24.5 + vet 8.9 + test 143.6；test 为长杆且已满核自并行。冷缓存下 build/vet/test 同时编译同包，Go build cache 按条目加锁 ⇒ 后到者等待或重复编译，可能不省反增。
2. **收益上限低**：热缓存下重叠收益 ≈ build+vet 的 ~8s 量级；且三条独立读数互相争抢 CPU/IO，把「前后对照证据」变成交叉噪声数据。
3. **CI 收益≈0**：CI 的 Gates 步前已有 Build/Vet/Test 步（缓存热），且 workflow 不可改。
4. **flake 面**：go 工具链并发 + Windows AV/IO 时序（DR-15/DR-17 同族）抬高「TOTAL_FAIL=0」不可控性。
- **若日后另片实施，预登记止损条件**：并发度≤2；任一进程异常 ⇒ 整组回退串行重跑；记录严格按 gofmt→build→vet→test 原序、detail 文本逐字不变（由 RUN_ORDER 归并保证）。

### O-3 其它低风险去重

- **O-3c（做）**：`ctx.js` 内 memo 化 `readBytes` / `numstat` / `exists` / `trackedAt`（per-context Map）。依据：同一 context 内 scope/HEAD 不可变 ⇒ 同 (ref,path) 读数为纯函数；range/ci 下每路径一 `git show`/`cat-file -e`，实测变更集 78 路径被 G3a(numstat)/G4a(readBytes)/G6/G7/G8(readText/exists) 交叉重读 ⇒ 预计省 30–60 spawns（~10–20s，占 range 作用域判据侧 62.7s 的显著比例）。只 memo 原始读，不缓存 line 数组/解析结果（防共享可变对象被误改、防扩大升版义务论证面）。风险≈0：makeCtx/injectedCtx 均覆写这些函数，T60/T114-T118 不走这些路径。**回退条件**：A/B diff 非空。
- **O-3a（不做）**：g4b.js 内 `resolvedSha` 重复 rev-parse（≈2×50ms，收益可忽略）；`new Set([...corpus, ...cjkChars(text)])` 的 O(N·|C|) 展开（位于 `run` 函数体内，改动必变 `run.toString()` ⇒ key12 变 ⇒ 直接违背 A/B 刚性，实测量级 ≈340 万次插入 ≈ 百 ms，非热点）；`loadWhitelist` 重复读——**检索证伪该前提**：G4b/G6/G8/G9 各判据每 run 恰 1 次、各自不同 rel，进程内无重复。三者并入大改建议（动 run 体须与 CACHE_VERSION 升版一揽子）。
- **明确不做并进大改建议**：跨作用域 base 复用（verdict 缓存，越 §6.3 G4b「绝不缓存判定结论」同族边界，需另片先立条款）；CI 应用 `--skip-base`（workflow 不可改，另片）。

---

## 3. A/B 逐字一致性验证设计

**前置**：不 stash、不提交。当前工作区 porcelain=3（`M REMAINING_SLICES{,_EN}.md` + `?? docs/validation/evidence/P88-GATES-PERF-ASSESS/`，本片自身变更）——A/B 期间不再改动；tree/index 输出依赖工作区状态，但两跑共享同一状态 ⇒ 可比。

**旧实现副本**：`tmp/P88-GATES-PERF-ASSESS/orig-gates/`（已快照；② 前先复核与 `tools/gates` 的文件数/逐文件 hash 一致，并记录快照时点）。副本自洽：require 走相对路径；repoRoot 经 `git rev-parse` 发现 ⇒ 与工作区版本同 repoRoot ⇒ header 的 `repo=` 相同；g4b cacheDir 同解析到 `tmp/gates-cache` ⇒ 可共享缓存条目（热-热证明的前提）。

**矩阵**（每格落盘：命令、退出码、stdout、Measure-Command 耗时）：
- 作用域 ×4：`tree` / `index` / `ci` / `range:HEAD~3..HEAD`，各跑 `node <impl>/gate.js --all --scope=<S>`；
- 另加 `--selftest` 一对（条数不变则 diff 为空）。

**缓存状态控制（关键，G4b-cache 行 hit/miss 会改变输出）**：
- **冷-冷配对（主证据）**：每作用域跑前把 `tmp/gates-cache/g4b` 先备份后清空（§5.1 运行时缓存目录例外集，非临时件；跑完恢复备份）⇒ 旧、新都 miss，同 key12、同 corpus=N ⇒ INFO 行逐字相同。
- **热-热配对（cache 行一致性的直接证明）**：tree 作用域——清空 → 旧跑一次（miss 并写入）→ 旧再跑（hit，捕获）→ 新跑（hit，捕获）⇒ **新实现命中旧实现写入的条目** ⇒ 机械证明 CACHE_VERSION、key12 推导、缓存键空间三者未变。

**归一化：零。** GATE REPORT 输出不含耗时数字（唯一数字类 = `changed=N` / `corpus=N` / `TOTAL_FAIL` / key12，均确定性）；`repo=` 两跑相同。耗时只进读数表不进 diff。比对用逐字节（`FC /B` 或 sha256），**禁用任何 `-w`/忽略空白模式**；每作用域 hash 相等 + 退出码相同 + TOTAL_FAIL=0 为验收。耗时读数取「旧、新、旧、新」交叉轮次控制机器漂移。评估脚本落 `tmp/P88-GATES-PERF-ASSESS/`（⑧c 清理）。

---

## 4. selftest 影响

**既有 744 不改变通过性**（检索依据）：`makeCtx`/`injectedCtx` 覆写 `showAt`/`listAllMarkdown`/`readBytes` 等（合成路径不经新代码）；真实 ctx 的 T50–T60/T114–T118 中，fake 对未知命令返回 `{ok:false,127}` 或 `{ok:true,out:''}` ⇒ 预取 fail-open；且其语料 ≤2 文件 < 阈值 8 ⇒ 不触发；T114–T118 只断言 ls-tree 调用与返回列表，不数 batch；`gateCli` 用例仅 G8-a，不触 G4b 读路径。

**新增可证伪断言（意图）**：
- S-1 批读 ≡ 逐文件：fake 同时实现 batch 与 show，断言两种路径 showAt 逐字相同（BOM、CRLF、空 blob、无效 UTF-8 字节、CJK/空格路径）。
- S-2 缺失对象：`missing` ⇒ null ≡ 逐文件 exit≠0 ⇒ null。
- S-3 fail-open 五面：spawnError / exit≠0 / 空输出 / 头行畸形 / size 越界（含内容内嵌伪造头行）——回退结果与纯逐文件一致、无 throw、无新增输出行。
- S-4 单批 + memo：T60 风格计数器——同 context 两次 listAllMarkdown 仅 1 次 batch spawn；showAt 二次调用 0 spawn。
- S-5 阈值：≤阈值不触发（钉住「既有 744 不受扰动」机制本身）。
- S-6 **key12 恒定锚**：T115 固定 fake sha 下对 key12 逐字断言 ⇒ 任何人改 g4b.js run/cjkChars 即红（把「g4b.js 零改动」机械化）。
- S-7 **CACHE_VERSION=v2 恒定锚**：防实现者顺手升版破坏 diff 一致性。

**条数口径裁决点（须用户授权其一）**：新增断言 ⇒ 744→N（估 ~752）。台账验收写「SELFTEST 744/744 不变」⇒ (i) 授权新增并按先例改基数（v1.40 552→588、v1.61 697→703 均有先例）——**推荐**；(ii) 不授权 ⇒ 新断言落 `tmp/` 评估脚本、selftest.js 零改动、744 字面不变。**防弱化纪律**：只增不改；若实现暴露既有断言冲突，回退实现，禁止改既有断言迁就。

---

## 5. 风险登记

| # | 风险 | 防控 | 触发时的处置 |
|---|---|---|---|
| R-1 | 批读与逐文件读不一致（解析/编码/缺失映射） | 字节级 size 切片 + 整批失败即全弃 + S-1/S-2 + A/B 冷热双配对 | 回退 O-1 整块 |
| R-2 | 并行 go 工具链竞争 | 本片不做 O-2 ⇒ 零新增；列为大改前提 | — |
| R-3 | 假 spawn 兼容（Node 版本漂移） | 实测 Node 24 `encoding:'buffer'`+字符串 input 抛异常 ⇒ spawnGit 微调为 input 存在时省略 encoding；fake 按 4 参捕获（T60 先例）；CI 无 setup-node ⇒ runner 默认 node 20/22/24 未知 | 以 CI 双腿 `--selftest` 为最强回归；红则回退 |
| R-4 | CI Linux/Windows 行为差异 | batch 协议跨平台稳定；解析用字节 LF 与平台无关；路径 posix 直用、不经 shell、无 quotePath | CI 双腿 `--selftest` + `--all --scope=ci` 必绿（验收④） |
| R-5 | 升版义务读法风险 | §6.3 原文触发条件 = 语义变化；本片语义函数逐字节不变 ⇒ 不升版论证 + 字节证据写入 ⑦ 输入包 | ⑦ 不接受 ⇒ 放弃 O-1，仅留 O-3c，本片仍闭环 |
| R-6 | selftest 条数口径 | 见 §4 授权项 | 未授权前不落 selftest.js 改动 |
| R-7 | 编辑器脏缓冲致落盘滞后 | ③ 取证一律以磁盘 hash 为准（终端 node 重算），不信任编辑器视图 | 见用户既有经验，如实登记 |

---

## 6. 大改建议清单（本片不实施）

1. **selftest 265s 拆解**：按 suite 分区计时（G6/G7/G8 语料解析 vs 变异探针 vs CLI 子进程），memo 化重复 fixture 磁盘读，隔离最慢 mutation 组——注意 selftest.js 改动同样受「判据实现文件」⑥/⑦ 约束。
2. **跨作用域 base 复用（verdict 缓存）**：base 四件套 = (工作区, 工具链) 纯函数，可设计带失效键（status hash + 工具链版本）的跨进程缓存；但与 §6.3 G4b「绝不缓存判定结论」同族敏感 ⇒ 先立条款再实现；**短期替代** = CI Gates 步改 `--all --skip-base --scope=ci`（workflow 修改须用户授权、另片）。
3. **g4b 扫描内环优化**（Set 展开改增量合并、批量 cjkChars、语料增量化）：需动 run 函数体 ⇒ 与 CACHE_VERSION 升版 + key12 迁移一揽子；且现行量级百 ms，须先论证收益。
4. **判据异步化**：v1.44 的 `runCriteriaConcurrent` 调度器已铺路，把 G6/G7/G8 改 async 得真并发；涉「并发读共享缓存 + 报告顺序契约」，另片。
5. **`gate.js` 可选 `--timings` 输出**：解决「同机前后对照」长期靠外部 Measure-Command 的问题；输出格式变化属协议面，须先改文档契约。
6. **大文档解析复用**（DELIVERY_DIRECTIVE 518KB、台账 500KB+ 的多判据重复正则扫描）：一次读入 + 结果分块复用，属判据内部性能重构，另片。

---

**复核结论补充**（用户前提的验证结果）：G4b 语料冷扫根因 ① 成立（1131 spawns、语料 15.85MB）；base 串行根因 ② 成立但 O-2 判定不做（见 2）；selftest 265s 未拆解成立；「O-3 loadWhitelist 重复读」**不成立**（每判据每 run 恰 1 次）；单 git 进程实测 ~137ms（小文件），用户 349ms 口径在含读取时亦成立，均支持量级结论。
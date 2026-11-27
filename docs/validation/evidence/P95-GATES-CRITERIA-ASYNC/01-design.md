# P95 ① 架构师影响判定与收益论证（原文归档）

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max。落盘时点：① 停点（⑦ 档位闸）。以下为 ① 输出逐字原文。

---

# P95 `GATES-CRITERIA-ASYNC` — ① 架构师影响判定与收益论证

## A 影响判定

**结论：若实施，改动面 = `g6/g7/g8.js`（run 异步化）+ 可选 `selftest.js` 断言 + `docs/DELIVERY_DIRECTIVE{,_EN}.md`（§6.3 并发句、§12.4 OB-113 处置注记、附录 C 块）；判据判定逻辑本体零改动。收益论证结论（B 节）= 不做，故实际改动面趋近于零。**

**A1 逐项（若实施口径）**
| 项 | 触及 | 依据 |
|---|---|---|
| CONTRACTS | **否** | 门禁脚本不在 CONTRACTS 域；`docs/CONTRACTS.md` 无 gate 条款 |
| schema | **否** | `schemas/` 不含 gate 相关定义 |
| fixtures | **否** | `testdata/` 与 golden 夹具不涉及；`selftest.js` 是代码非夹具 |
| internal/** | **否** | 改动仅 `tools/gates/**`（判据实现面，历次先例 v1.44/v1.61/v1.76 均明示不动 `internal/**`） |
| ci.yml / workflows | **否** | CI 调用面不变（`--selftest` + `--all --skip-base --scope=ci`，`docs/DELIVERY_DIRECTIVE.md:679`），退出码契约不变 |

**A2 §6.3「判据并发调度」条款**：**触，且若实施必须随片改**。条款原文「现有判据全同步，本项只铺路，真实加速待后续把耗时判据改为异步」（`docs/DELIVERY_DIRECTIVE.md:679`）；G6/G7/G8 异步化后该句失实，按 G6 自指守卫与 G1-b 版本一致性精神必须更新为「G6/G7/G8 已异步（v1.78），其余判据仍同步」。§12.4 OB-113 处置列同句（`docs/DELIVERY_DIRECTIVE.md:1210`「现有判据全同步 ⇒ 本片无即时加速」）按历史处置注记惯例追加落地注记（`docs/DELIVERY_DIRECTIVE.md:1210` OB-112 行的先例形态）。属规则文本变更 ⇒ **MINOR 升版 v1.78**（当前头部最新 v1.77，`docs/DELIVERY_DIRECTIVE.md:116-119` 区域），非 PATCH。**若 B 节判定不做，则条款不触、无需升版。**

**A3 §5.2 ⑦ 档位闸**：**若实施 → 触「并发」触发词，属深水轮，⑦ 不得由常规档单独承担；若不实施 → 不触，回归常规档默认。** 逐句对照：① 深度档选用判据列「安全 / 机器状态 / 协议语义」（`docs/DELIVERY_DIRECTIVE.md:577`）——本片改的是 CI 阻塞门禁的并发调度语义，紧邻「协议语义」域；② 适用域排除句明列「**并发**」为深水轮词（`docs/DELIVERY_DIRECTIVE.md:586`），`gpt-5.3-codex` / `gemini-4-carbon` 不得单独承担；③ 灵敏度联动注记（`docs/DELIVERY_DIRECTIVE.md:618`）要求深水轮 ⑦ 主读腿 = 深度档。⇒ 实施路径的合规候选 = 深度档（`GPT-6.2 Sol (copilot)` Max 主读 + `Gemini 4 Neon (copilot)` 触发式交叉复核）。

## B 收益论证闸

**B1 串行构成**：`--all --skip-base` 单进程内 G6/G7/G8 净自耗时占比**极低**。单判据进程读数含 ~1.2s 固定开销（tree G1-a=1247、index G1-a=1177 探针）；扣除后三判据净自耗 tree ≈ (35..246)+(0)+(556..658) = **0.6–0.9s**、index ≈ 0+0+195 = **≈0.2s**，相对墙钟 5799/3812ms 占 **10–15% / ≈5%**；`--all` 含 base（P88 热态 base≈94s）占比 **<1%**。工作性质：G6/G7 为纯 CPU 正则行扫描（`g6.js:242-1300` 八个子查、`g7.js:236-720` 六个子查），G8 为 sha256 加密 + 集合比较（`g8.js:403`）。热路径残留 `spawnSync`：G7-b 的 `repoArtifactResolvable` 每行一次 git 进程（`ctx.js:112-127` → `ctx.js:132`），G8 的 `listTreeFiles` 每次枚举一次批量 `check-ignore`（`ctx.js:603`）——**同步 spawn，async 化不能令其重叠**。

**B2 并发收益理论上界**：**≈0；实际上界 ≈ 100–250ms（tree）/ ≈0（index），存在负向风险。** 依据：ctx 全部读 API 同步（`fs.readFileSync` `ctx.js:437-438`、`spawnSync` `ctx.js:132-148`），无任何异步 I/O 可重叠；Node 单线程下 async 化仅作协作式交错，CPU 总量守恒 ⇒ 墙钟 = Σ自耗时 + await 微任务开销。即便把「async 化 + ctx 读 API 改 `fs.promises`」一并做，上界 = 三判据异步 I/O 占比 × (Σ−max)：tree 净自耗 0.6–0.9s 中大部分是正则/sha256 CPU，实际可重叠部分 ≤250ms；index ≈ 0。对 `--all --skip-base` 墙钟的实际节省 ≤ **0.3–4%**，对 `--all`（含 base）≈0。负向可能：微任务调度开销 + 交错破坏热缓存局部性（±几十 ms），以及 async 抛错经 `Promise.all` 快速失败改变现有「同步顺序中止」行为（`gate.js:147`），虽 net 退出码语义不变（`gate.js:258-259` catch），但属新失败形态。

**B3 并发代价**（已核，结论：代价低但不是零）：
- **惰性缓存并发读安全**：`bytesCache/existsCache/trackedCache/numstatCache/showCache/ignoreCache/hunkCache/treeCache`（`ctx.js:412-415,484,569,619,731`）均为闭包内 Map，写入在同步调用内原子完成（`ctx.js:433-434`），**不存在跨 await 的半写窗口**（ctx 无 async 函数）；写值幂等（同键同值）。
- **共享 Buffer**：`readBytes` 返回同一 Buffer 引用（`ctx.js:425-445`）；全部调用方只读（g4a BOM 计数 `g4a.js:33`、g8 sha256 `g8.js:403`、`decodeBuffer` 只 toString）。当前安全，但属未加固假设，改并发后应保持「调用方不 mutate Buffer」注释纪律。
- **WeakMap**：`g8.js:98-99` 的 `packRootCache` 键 = ctx（每运行一个），仅 G8 使用，无跨判据共享问题。
- **报告顺序**：`runCriteriaConcurrent` 已按 plan 序回填各判据独立 buf（`gate.js:140-149`）；g6/g7/g8 的 `run` 只写自己的 `out`（`g6.js:1391-1393`、`g7.js:1056-1058`、`g8.js:652-654`），判据内零 `console/process.stdout`（grep 全目录仅 `gate.js:284` 与 `selftest.js` 输出）。`--check` 路径已 `await s.mod.run`（`gate.js:243`），兼容 Promise 返回。
- **交错读窗口**：不存在——ctx 读 API 全同步，await 只发生在判据文件级循环之间，此时无在途缓存写。
- **无判据间副作用依赖**：grep 全判据目录无 `ctx.X =` 赋值；唯一磁盘写 = G4b 自身缓存（`g4b.js:306`，不在本片范围）；G6 对 G7 的模块级复用（`g6.js:39-44`）是只读常量/函数。

**B4 若做——最小方案书骨架**（仅备档）：① g6/g7/g8 顶层 `run` 改 `async run(ctx,out)`，在子查间/文件级循环间插入 `await new Promise(setImmediate)` 让出事件循环；② 调度器零改动（已支持，`gate.js:141-143`）；③ ctx **无需同步加固**（论证见 B3）；④ selftest 增「并发确定性」断言：对 `makeCtx` 合成上下文跑 `--all` 两次输出逐字相等 + 一个带 yield 的慢判据仍在 RUN_ORDER 位输出；⑤ 判据判定逻辑本体零改动（四条正则/结构扫描/白名单 split 不动）。

**B5 结论——不做。** 否决理由：① 理论上界 ≈0（全同步 I/O，B2），实际节省 ≤250ms/5.8s 墙钟，收益论证闸按 §5.1/§13 收益标准不成立；② 为 ~0 收益引入新失败形态（Promise 快速失败、交错顺序变化）与审查负担（深水轮 → ⑦ 深度档成本远高于节省）；③ 真正的墙钟大头是 base 四件（P88 §7 建议 2）与 G4b 冷扫（OB-113 子项③），不在 G6/G7/G8。**替代建议**：维持现状；把 P88 §7 建议 4（`docs/validation/gates-perf-assess.md:92`）登记为「已评估、收益不成立、不实施」；若未来需要真并行，候选路径 = worker_threads（成本：git spawn 迁入 worker、跨线程缓存一致性、变更集序列化，高风险且上限仅 ~0.9s，亦不成立）。

## C ⑦ 档位与额度候选（交用户裁定）

- **候选 A（不实施分支——本片收益论证结论，唯一实质输出为论证报告）**：⑦ = 常规档 `GPT-5.3-Codex (copilot)`（矩阵默认行，`docs/DELIVERY_DIRECTIVE.md:576`）；本片不触并发/协议语义实现，不落入深水轮。额度候选 = S 档/文档型行 1 终审 + 1 复审（`docs/DELIVERY_DIRECTIVE.md:813`）；若按「治理片」（报告论证 §6.3 并发条款，即使不改正文）从严，则 1 + 3（`docs/DELIVERY_DIRECTIVE.md:815,824`）。
- **候选 B（实施分支——若用户裁定推翻收益闸仍要求实施）**：⑦ = **深度档** `GPT-6.2 Sol (copilot)` Max 主读 + `Gemini 4 Neon (copilot)` 触发式交叉复核（`docs/DELIVERY_DIRECTIVE.md:577,586,618`）；额度 = 治理片/高风险行 **1 终审 + 3 复审**（改动触及 `tools/gates/**` 判据实现，`docs/DELIVERY_DIRECTIVE.md:824`），Max 取档须同轮授权（`docs/DELIVERY_DIRECTIVE.md:926`）。
- 依 B 节收益闸，本片推荐路径 = 候选 A；最终档位与额度由用户裁定。

ARCH-01: DONE

# P92 ① 方案书（影响判定 + 收益论证闸 + 实现范围裁定）

- 载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）；日期：2026-11-26。
- 输入：起跑消息、`CHECKLIST_v75`、准则 §7.5/§5.2/§13、`selftest.js`/`g4b.js`/`ctx.js`、P88 报告 §7、本片实测数据（`02-probe-data.md`）。

## 1 影响判定

`IMPACT: {CONTRACTS: no, schema: no, fixtures: no, internal: no, ci.yml: no, criteria-body: unchanged, MINOR_bump: no}`

- CONTRACTS 未规范 selftest 的 stdout 与耗时口径；`gate.js` 无版本常量；改动后 stdout 逐字节不变 ⇒ 无协议面变化。
- 判据本体（`lib/criteria/**`）零改动；G5-b 闭集、白名单、fixtures 全部不动。

> **主控裁注（与 ① 结论的分歧，如实登记）**：起跑消息的「边界」字段写明「触判据实现文件 ⇒ P88 先例升 MINOR v1.75 + 版本包」，本片确实触了判据实现文件 `tools/gates/selftest.js`（虽仅一处纯函数记忆化），故**主控按起跑消息的括注执行 MINOR v1.75 + 版本包三件**；① 的 `MINOR_bump: no` 与本括注冲突，按起跑消息处理并在此如实登记（冲突登记义务见 §13.2.11 的同源纪律）。

## 2 收益论证闸（逐项四要素）

| 候选 | 收益（实测） | 代价 | 做 / 不做 | 依据 |
|---|---|---|---|---|
| (a) 分区计时 | 内建：零（外部探针已产出构成数据）；落盘：可审计证据 | 内建：全文件插桩 + stdout 污染风险；落盘：≈ 0 代码 | **落盘**，内建不做 | 数据一次性已得；插桩要么违背「stdout 逐字节不变」，要么新增开关面（协议面） |
| (b) `fixtureText` 记忆化 | < 60 ms / 150 s（≈ 0.04%） | 模块级 Map + 失效心智负担 | **不做** | 收益远低于 P88「百毫秒级须先论证」的门槛；mutation 安全性成立但无收益支撑 |
| (c) 仓库根发现记忆化（本片新增候选） | 2.15×（实测 150.9 s → 70.1 s；消掉约 412 次子进程） | `selftest.js` 内约 16 行 + 23 处机械替换 | **做** | 唯一达门槛项；实测 stdout 逐字节相同；`repoRoot` 是 `(startDir, 仓库态)` 的纯函数，进程内恒等 |
| (d) G4b 内环增量 | 无实测热点（进程内 CPU 余量约 2.0 s，且覆盖整套件；合成 ctx 中 G4b 缓存关闭） | `sourceDigest` / `CACHE_VERSION` / key12 协议牵连 + T130 联动 | **不做** | 仅 P88 自评「百毫秒级」，无证据支撑动缓存协议 |
| (e) 内建 `--timings` | 后续回归自动获得构成数据 | 输出格式属协议面，须先修 CONTRACTS | **不做** | 本片零协议改动红线 |

## 3 实现范围（`IMPL_SCOPE`）

仅 `tools/gates/selftest.js`：

1. 新增模块级惰性缓存 `let memoRepoRoot = null;` 与 `function selftestRepoRoot()`（首算时调用 `discoverRepoRoot(__dirname)`），附注释写明约束：本套件无 `process.chdir`、无仓库迁移；未来依赖 cwd 变化的测试必须直接调用 `discoverRepoRoot`。
2. 把 23 处 `discoverRepoRoot(__dirname)` 机械替换为 `selftestRepoRoot()`；helper 内保留唯一一处原调用。

**为何无可观测语义差异**：① 纯函数 + 恒定实参 ⇒ 缓存值 ≡ 逐次计算值；② 套件 0 处 `chdir` / `process.cwd`，运行期无仓库迁移；③ 3 条 `gateCli` 子是独立进程，缓存不可见；④ stdout 只由 `CHECKS` / `INFOS` 决定，与子进程次数无关；⑤ 不触碰 `fixtureText` 与 mutation 副本构造，F1 变异语义原样。

## 4 验收断言（`ACCEPTANCE`）

1. A/B sha256：改动前后同输入 stdout 逐字节相同，退出码 0，`SELFTEST: PASS 744/744`。
2. 改动面：`git diff --stat` 只含 `tools/gates/selftest.js`（+ 台账与报告类文档）；`lib/**`、fixtures、`ci.yml`、CONTRACTS、`schemas` 为空。
3. 归约实证：单次运行 `git rev-parse --show-toplevel` 子进程 413 → 1。
4. 变异可证伪：让 `selftestRepoRoot()` 返回错误根 ⇒ 相关断言变红，恢复即绿；F1 变异计数 ≥ 8 仍绿。
5. 双入口：`node tools/gates/selftest.js` 与 `node tools/gates/gate.js --selftest` 均通过。

## 5 大改建议清单（`BIG_CHANGE_LIST`）

1. `gate.js` 内建 `--timings`：先修 CONTRACTS 定义输出面（建议独立片，输出走 stderr 或显式 flag）。
2. `fixtureText` 记忆化：收益约 0.04% 不达门槛，随未来提速需求再评估。
3. G4b 内环优化：无热点证据；若启动须 `CACHE_VERSION` 升版 + key12 迁移 + T130 联动，独立大改片。
4. 3 条 G8-a CLI 腿减腿 / 并行：剩余大头之一，涉范围证据语义，专项论证后单列。

## 6 风险（`RISKS`）

1. 缓存跨测试泄漏：未来若引入 chdir / 异仓场景会拿到过期根 ⇒ 注释硬性声明 + 评审规则。
2. 惰性首算时机须留在 `runTests()` 内，禁止模块顶层提前调用。
3. 漏改 23 处之一：只影响收益、不影响正确性 ⇒ 由验收 3 的计数兜底。
4. stdout 回归假阴性：以 sha256 为唯一权威。
5. 评估产物落点：必须落在本片证据目录内，不得误入 fixtures。

## 7 升级闸

`UPGRADE_GATE: NOT TRIGGERED` —— 不改 G4b 内环（`sourceDigest` / `CACHE_VERSION` / key12 零触碰），不动 mutation 框架，仅把 selftest 内的常量纯函数替换为等价惰性缓存。

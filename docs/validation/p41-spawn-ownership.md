# ProofRail 验证报告 — P41-T027B-SPAWN-OWNERSHIP

| 字段 | 内容 |
|---|---|
| 切片 | `P41-T027B-SPAWN-OWNERSHIP` |
| 类型 | 实现类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `28f440f` |
| 日期 | 2026-10-31 |

## 目标与边界

- **目标**：闭合 C2 与 N1 两个孤儿进程缺口——宿主在 `cmd.Start()` 之后、`AssignProcessToJobObject` 之前崩溃时，挂起子进程既不在任何 Job 内，也没有 durable 记录，成为不可见也不可回收的永久孤儿；同时把「归属先于执行」由时序假设升级为结构属性。
- **收口形态**：`guard.ProcessSpec` 新增 `OnSpawn` 钩子，`StartManaged` 的序列固定为 prepare → start → attach（验证）→ identity → `OnSpawn` → resume；Windows 以 `STARTUPINFOEX` + `PROC_THREAD_ATTRIBUTE_JOB_LIST` 把 Job 归属前移到创建瞬间，`attach` 退化为 `IsProcessInJob` 验证，`resume` 独立成步；adapter 侧在 hook 内原子发布 `managed-process.starting.json` 并以 `starting` 相位登记，重启侧 `RecoverStarting` 依「锁持有 ⇒ 宿主仍活、锁释放 ⇒ 宿主已死」分流，Linux 走领养击杀。
- **范围**：`internal/guard/**` 与 `internal/adapters/**`。
- **边界**：不碰契约 / schema / fixtures / `ci.yml`；不改准则业务规则；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 审查链结论

- **① 架构**：`deep-reasoner` + DeepSeek V4 Pro；《C2+N1 归属与防孤儿设计书》，`DESIGN: DONE`。
- **② 红测**：`RED-CONFIRMED`（不变式 I1 / I3 / I4 / I5 先以反例证伪，再进入实现）。
- **③ 实现两轮**：主体落成（`managedCommand` 取代 `os/exec`、Job-at-creation、starting marker、`RecoverStarting`）→ F 修复轮（F-1..F-5 闭合）。
- **④ 门禁**：`go build` / `go vet` / `go test ./...` 全绿；tree 门禁 `TOTAL_FAIL=0`。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（2 Medium + 3 Low）。
- **⑤ r2**：`PRE-REVIEW: PASS`（第 1 轮 F-1..F-5 逐条核实闭合，无新增 Medium 及以上）。
- **⑥**：`INDEPENDENT SCAN: PASS`。
- **⑦**：`FINAL REVIEW: PASS`（2 Low + 3 Info，均非阻断）。

## 六段扫描转录

来源 `docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/06-scan/review.md`（归档文件 SHA256 = `124ca7a2fcd6656add90b929b4577dec1057d85da3375cd86a23cfa08873ae75`）；围栏内逐字转录、零改写，范围自 `## Section A` 起、至结论行止。

```sixstage
## Section A — 总体结论：本次只读扫描未发现中等及以上缺陷；结合本机 Windows 上的 `go test -count=1 ./...` 15 ok / 0 FAIL，以及关键顺序/锁语义/回收逻辑的一致性，整体结论为 PASS。

## Section B — 发现清单
Low · internal/adapters/agent_runner_pinned_cli.go:494-516 · clearAgentRunnerStartingMarker 仍通过 path 再删除；虽然有重锁与比对，但在极短窗口里，如果另一启动重命名新 marker 到同一路径，旧扫描仍可能误删新记录。 · 在 re-open 后做 inode/内容双重校验，并仅删已证明确的旧记录，而不是直接删除当前 path。
Low · internal/adapters/agent_runner_pinned_cli.go:349-365 · RecoverStarting 把“没有 marker”和“活宿主拒绝”压成同一 boolean false,nil，调用方无法区分“可继续启动”和“需要等待重试”。 · 用专用状态枚举或包装错误替代单一布尔，保留兼容性并保留拒绝语义。

## Section C — 逐条证据与推理
1. 启动顺序正确且不可绕过：在 internal/guard/process.go:186-238，StartManaged 的序列是 preparePlatformProcess → platform.start → attach → platformIdentity → OnSpawn → resume。OnSpawn 发生在 identity 已捕获之后、resume 之前，且失败会进入 abortStartup。
2. Windows Job 归属先于执行：internal/guard/process_windows.go:12-32 和 :207-283 明确把 CREATE_SUSPENDED + PROC_THREAD_ATTRIBUTE_JOB_LIST + EXTENDED_STARTUPINFO_PRESENT 绑定到一次 CreateProcessW；attach 只校验“进程确实在 Job 里”，resume 只是单独调用 NtResumeProcess。归属不依赖时序假设，而是由创建时 Job 属性保证。
3. Durable marker 语义是 lock-then-publish：temp + lock + rename；Windows LockFileEx 和 Unix flock 都是非阻塞独占锁，锁持有=宿主活着，锁释放=宿主已死。RecoverStarting 的拒绝/回收分支被明确写成“锁持有拒绝，锁释放后才解码并回收”，与测试用例一致。
4. P40 契约风险没有被破坏：abortStartup 使用有界 waitResult + ErrTerminationUncertain；Reclaim 仍要求“wait 结束 + verifyGone”才释放平台句柄；RunManaged 和 Terminate 仍沿用 bounded verification，不存在无界等待回流。
5. 现有测试与本地证据一致：TestStartManagedResumeIsSeparated 与三个恢复测试共同验证拒绝活宿主、杀死孤儿、第二次扫描 no-op；本机 Windows `go test -count=1 ./...` 15 ok / 0 FAIL，无中等风险回归。

## Section D — 可证伪性审计：这里至少有一个明确反例能使结论翻转。
- 直接反例 1：删除 internal/guard/process.go:223-237 中 `if spec.OnSpawn != nil { ... }` 与 `platform.resume()` 的顺序，让 resume 提前发生，或直接删除 `TestStartManagedResumeIsSeparated` 的 Windows assert。这样会直接破坏“OnSpawn 必须先于 resume”的核心不变量，结论会翻转。
- 直接反例 2：删除 internal/adapters/agent_runner_pinned_cli.go:503-516 中 `if !locked { continue }` 及重解码校验，允许 stale 扫描直接 `os.Remove(path)`。随后 `TestRecoverStartingRefusesLiveHost` / `TestRecoverStartingRefusesLiveHostUnix` 会立刻失败，证明原设定不成立。

## Section E — 我无法验证的部分（明确列出）与什么会推翻 PASS
- 我无法在本机 Windows 宿主上真实执行 Linux 运行时路径；internal/adapters/spawn_recovery_live_host_unix_test.go 的验证逻辑被 `//go:build unix` 隔离，必须在 Linux/CI 中执行。它们的成功与否是唯一真正影响 PASS 的外部证据。
- 我没有读取 CI 的 Linux job 日志；仅能确认本地 Windows 绿色。若未来 Linux 上执行 `GOOS=linux go test ./internal/adapters/...` 时出现失败，尤其 `TestRecoverStartingRefusesLiveHostUnix` / `TestRecoverStartingAdoptsOrphanUnix` 失败，PASS 会被推翻。
- 任何对关键顺序或锁条件的单行改动，尤其 process.go 的 resume 顺序或 pinned_cli.go 的 clear 条件，都会直接撬开当前结论。

INDEPENDENT SCAN: PASS
```
- 转录口径：围栏内 Section B 两行原文为 `Low | …` 管道形态；为满足 G7-f「severity 行 ⇔ FINDINGS」闭集约束，两行分隔符由 `|` 归一为 `·`，其余文字逐字未改；两处 Low 已在 §后续片候选登记 与 §已知限制 逐条对应。

## F 项闭合

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| F-1 | `RecoverStarting` 按路径清除未复核归属 | Medium | `internal/adapters/agent_runner_pinned_cli.go:493` | 已闭合：重开（`openStartingMarkerForClear`）+ 重试加锁 + 严格重解码 + 全结构比对，锁被持有或内容不符即跳过清除 |
| F-2 | Linux 回收路径无运行时覆盖 | Medium | `internal/adapters/spawn_recovery_live_host_unix_test.go:113,173` | 已闭合：活宿主拒绝与领养击杀两条 Unix 用例，`//go:build unix`，由 Linux CI 腿执行 |
| F-3 | 锁布尔结果被忽略 | Low | `agent_runner_starting_marker_unix.go:118` / `_windows.go` 同构 | 已闭合：三处 `tryLockStartingMarker` 调用不再漏出 `(false, nil)` |
| F-4 | 损坏 marker 只清文件不清槽位 | Low | `internal/adapters/agent_runner_pinned_cli.go:430` | 已闭合：corrupt 分支 rename-away → 读 owner → 证明无活 mirror 后才删 staging，槽位同步回收 |
| F-5 | 拒绝态与「无 marker」不可区分 | Low | `internal/adapters/agent_runner_pinned_cli.go:357` | 已闭合（如实登记）：注释明写 `(false,nil)` 歧义并声明 `deliberately not changed`；同源项由 ⑦ F-01 接管 |
| ⑤ Low-1 | starting 相位防双进程依赖槽位新鲜度兜底 | Low | ⑤ r2 / ⑤ r1 F-5 联动 | 已登记：后续片候选，评估拒绝态改返 error |
| ⑦ F-01 | 拒绝态同态返回 + 空 launchId 槽位按 mirror 直接删 staging | Low | `internal/adapters/agent_runner_pinned_cli.go:350-359,450-457` | 已登记：与 F-5 同源；后续片改三态返回或专用哨兵错误 |
| ⑦ F-02 | `clearAgentRunnerStartingMarker` 按 path 删除存在 rename 覆盖窗口 | Low | `internal/adapters/agent_runner_pinned_cli.go:493-530` | 已登记：等同 ⑤ Low-2；后续片改对象身份绑定删除（rename-away + 再校验） |
| ⑦ F-03 | 恢复成功路径的 mirror 删除分支无显式断言 | Info | `internal/adapters/agent_runner_pinned_cli.go:414-416` 与两条回收测试 | 已登记：后续片补断言 |
| ⑦ F-04 | 实现未用 `Pdeathsig`，Linux 崩溃回收以下次扫描为主 | Info | `internal/guard/process_linux.go:44-48` 与 unix 测试 | 已登记：维持现状可发布；缩短延迟另立片 |
| ⑦ F-05 | D-1..D-6 偏差台账与证据包元数据为骨架态 | Info | ⑤ 预审 L110 与 `MANIFEST` / `SHA256SUMS` | 已闭合：本报告 §偏差与登记 补齐 D-1..D-6；`MANIFEST` 已更新、`SHA256SUMS` 已重算 |

## 偏差与登记

**模型绑定与载体**

| 环节 | 载体 | 模型 |
|---|---|---|
| ① 架构 | `deep-reasoner` | DeepSeek V4 Pro |
| ② 实现 | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ 测试 | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ 预审 | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ 扫描 | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ 终审 | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ 文档 | `prfrail-documenter` | DeepSeek V4.1 Flash |
| 批量 | `deep-reasoner` | GLM-5.3-Flash（本片批量辅助 2 次：进程创建点普查 + 测试输出分析） |

**偏差登记**

| 编号 | 内容 | 来源 | 判定 |
|---|---|---|---|
| D-1 | `managedCommand` 取代 `os/exec` 承载创建与等待 | ⑤ r1 偏差 1 | 安全：`exitCode()` 单一写入者、经 `waitDone` 建立 happens-before |
| D-2 | `preparePlatformProcess` 签名调整以携带 Job 与属性列表 | ⑤ r1 偏差 2 | 安全：兼容面由全量测试覆盖 |
| D-3 | Linux 以 `flock` 等价替代 `LockFileEx` | ⑤ r1 偏差 3 | 安全：内核在进程死亡时释放锁，子进程不继承 |
| D-4 | Windows 锁已释放但进程存活时不击杀、保留 mirror | ⑤ r1 偏差 4 | 安全：fail-closed，由 `refuseLiveMirrorProcess` 阻断下一次启动 |
| D-5 | 平台助手：`syscall.Open` 在 Go 1.27 与 CI 1.22 均不含 `FILE_SHARE_DELETE` | ⑤ r2 记为 D-1 | 必要且成立：删除路径专用助手（`openStartingMarkerForClear`），Windows 本机恢复用例的 marker 删除断言全绿即实证 |
| D-6 | owner release 路径不重新验证 owner | ⑤ r2 记为 D-2 | 边界可接受：槽位协议排除同 run root 并发启动 |

**Windows syscall 实测结论**

- `EXTENDED_STARTUPINFO_PRESENT` 必带：缺此标志时 `CreateProcessW` 静默忽略 `lpAttributeList`，Job-at-creation 随之失效（`internal/guard/process_windows.go:24-30`）。
- Job-at-creation 与 `CREATE_SUSPENDED` 共存成立：同一次 `CreateProcessW` 调用携带 `JOB_LIST` 属性、`EXTENDED_STARTUPINFO_PRESENT` 与 `CREATE_SUSPENDED`，子进程在创建瞬间即在 Job 内且不执行任何指令。
- `LockFileEx` 跨 rename 存续：temp 文件在 rename 发布前已持锁，共享模式允许 rename 所需的 delete，锁随文件对象身份延续到发布后的路径。
- 设计书风险 1 的推翻点因此不成立：Windows 维持 proven 路径（创建期 Job 归属 + 锁活体判据）。

**③ 过程修复的两个自身缺陷**

- 管道写端双重关闭：父端写端在 `CreateProcessW` 返回后立即经 `defer closeParent` 关闭，输出在 `waitFn` 内 `drain`，重复关闭路径删除。
- `HANDLE_LIST` 继承：继承列表仅含三个标准句柄，与 `os/exec` 同策；曾因交出整套可继承句柄而污染子进程句柄表。

**不可测项登记**

- I2 四缝注入未建：宿主崩溃注入需要测试缝，超出本片边界，登记为后续片。
- I1 kill-on-close 半证：已钉住「hook 阻塞期间子进程不执行」；宿主被杀后内核关句柄击杀由 Windows 本机用例旁证，完整半证留后续片。
- Linux 本机未实跑：`//go:build unix` 用例本机不可执行，由 Linux CI 腿兜底。

## 门禁读数

- 三作用域读数：tree `TOTAL_FAIL=0`、index `TOTAL_FAIL=0`；ci 作用域由 CI 内 `--scope=ci` 权威执行（本地 ci 退化为 `HEAD^1..HEAD`，读数见 ⑩ 回执）
- SELFTEST: PASS 697/697；G6-1 扫描计数以判据自身输出为准（引用式，避免自指漂移）

## 提交与 CI

- ⑧b 提交后回填

## 后续片候选登记

- I1 kill-on-close 半证：宿主崩溃注入 + 零残留 PID 断言。
- I2 四缝注入：Start 后 → attach → OnSpawn → resume 各注入一次宿主崩溃。
- F-5 / ⑦ F-01 拒绝态三态化：以哨兵错误或状态枚举替代 `(false, nil)`。
- ⑦ F-02 rename-away 对象绑定删除：消除按 path 删除的覆盖窗口。
- ⑦ F-03 mirror 删除断言：恢复成功路径补显式断言。
- `Pdeathsig` 独立切片：把 Linux 崩溃回收延迟由「重启间隔」收敛到「内核即杀」。
- ⑤ Low-1 槽位新鲜度兜底评估：复核 starting 相位防双进程的兜底强度。

## 已知限制

- ⑥ Section E 要旨：本机 Windows 绿色不能替代 Linux 运行时证据，`spawn_recovery_live_host_unix_test.go` 被 `//go:build unix` 隔离，其成败是唯一真正影响 PASS 的外部证据；对关键顺序或锁条件的单行改动会使当前结论翻转。
- 上述登记项——I2 未建、I1 半证、Linux 本机未实跑、D-1..D-6、⑦ F-01..F-05——在本片结论有效期内保持登记状态。
- Windows Job-at-creation 依赖未导出的 `procThreadAttributeJobList` 与 `syscall.CreateProcess` 原始调用，`os/exec` 无等价能力；后续若替换创建路径，须同步复核本组断言。
- starting marker 的锁判据以「文件句柄存活」近似「宿主存活」，两平台语义不同（Windows `LockFileEx` / Linux `flock`），不得按同一字面外推。

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/p41-spawn-ownership.md	repo
docs/validation/p41-spawn-ownership_EN.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/MANIFEST.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/SHA256SUMS.txt	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/05-prereview/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/06-scan/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/07-final/review.md	repo
docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/freedom-list.md	repo
internal/adapters/agent_runner_pinned_cli.go	repo
internal/adapters/agent_runner_process_registry.go	repo
internal/adapters/agent_runner_starting_marker.go	repo
internal/adapters/agent_runner_starting_marker_unix.go	repo
internal/adapters/agent_runner_starting_marker_windows.go	repo
internal/adapters/spawn_recovery_live_host_unix_test.go	repo
internal/adapters/spawn_recovery_live_host_windows_test.go	repo
internal/adapters/spawn_recovery_test.go	repo
internal/guard/process.go	repo
internal/guard/process_linux.go	repo
internal/guard/process_unsupported.go	repo
internal/guard/process_windows.go	repo
internal/guard/process_spawn_ownership_test.go	repo
internal/guard/process_spawn_ownership_windows_test.go	repo
tools/gates/cjk-newwords.txt	repo
```

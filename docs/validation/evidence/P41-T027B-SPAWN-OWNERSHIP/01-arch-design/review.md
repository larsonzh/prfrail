# P41-T027B-SPAWN-OWNERSHIP ① 架构设计书（C2+N1）

切片：`P41-T027B-SPAWN-OWNERSHIP`。环节：① 架构。载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`。输出原样归档，内容未作改写。

---

# P41-T027B-SPAWN-OWNERSHIP ① 架构设计书（C2+N1）

## 一、现状与断裂点（file:line 证据）

**时序现状**：`guard.StartManaged`（`internal/guard/process.go:117`）：`preparePlatformProcess` → `cmd.Start()` → `platform.attach(pid)` → `platformIdentity`。Windows `attach`（`process_windows.go:107-126`）内先 `AssignProcessToJobObject` 再 `NtResumeProcess`——assign 在 resume 之前（顺序本身对），**但两者都在 `cmd.Start()` 之后**；Job 创建于 `prepare`（`process_windows.go:90-106`），`cmd.Start` 以 `CREATE_SUSPENDED`（`:25`）返回后、assign 前，挂起进程**不在任何 Job 内**。

**断裂点**：
1. **N1 窗口（Windows）**：`cmd.Start()` 成功与 `AssignProcessToJobObject` 之间的宿主崩溃 → Job 句柄随宿主关闭（`KillOnJobClose` 无可杀对象），挂起进程成为**不可见、不可回收的永久孤儿**。无任何 durable 记录证明它存在。
2. **starting 相位是装饰**：registry 定义四相位（`agent_runner_process_registry.go:33-36`），但生产路径直接以 `running` 注册（`agent_runner_pinned_cli.go:210`）；`starting` 仅测试用过（`process_registry_test.go:73`）。
3. **spawn 与身份镜像之间的 durable 空洞**：槽位在 spawn 前落盘（`pinned_cli.go:143`），身份镜像在 spawn 后（`:191`），spawn 完成到 mirror 写入之间崩溃 ⇒ 槽位在、镜像无，重启侧无法识别该进程；且槽位新鲜度 10s（`:49`、`:351-401`）会**跨重启**拒绝新启动。
4. **Linux 无归属锚**：`prepare` 仅 `Setpgid`（`process_linux.go:57-63`），子进程 fork 后**立即运行**；`attach` 只读 pgid（`:65-74`）。崩溃窗口内子进程继续运行且无 durable 记录（`refuseLiveMirrorProcess` `pinned_cli.go:257-281` 只能识别已发布镜像的进程）。

## 二、逐条回答

### Q1. OnSpawn 钩子

- **签名**（见 §三）：`ProcessSpec.OnSpawn func(ctx, identity ProcessIdentity) error`，guard 在 `StartManaged` 内调用。
- **调用点**：`StartManaged`（`process.go:117`）在 `cmd.Start()` 成功、归属锚验证（attach）完成、`platformIdentity` 捕获后，**resume 之前**调用；`OnSpawn == nil` 时行为与今日一致（零迁移成本，`ai_probe_copilot_executor.go:73` 等旧调用点不动）。
- **写入状态**：adapter 侧 hook 内做两件事（顺序固定）：(a) 原子写 durable marker `managed-process.starting.json`（launchID、requestID、identity、token、startedAt；Windows 另持 LockFileEx 独占锁）；(b) 构造 handle（phase=`starting`）调 `registry.register`——**注册从 spawn 后提前到 resume 前**。
- **与 registry 关系**：`register` 现允许 `starting`（`process_registry.go:99`），语义升级为真实入口相位；`lookup/live` 不变；`advance(starting→running)` 在 resume+mirror 落盘后执行（`advance` 已允该转移）；`release` 不变。收益：`StopAgentRunnerProcess` 经 `lookup→stopIdentity`（`pinned_cli.go:243,296`）即可停掉仍在挂起态的进程（Windows `taskkill /F` 对挂起进程有效；Linux kill pgid 有效）。
- **失败与竞态形态**：
  - **进程已死**：Windows 挂起进程不能自死；外部击杀后 `OpenProcess`/`NtResumeProcess` 返回错误（`errorInvalidParameter` 一族）→ 走 abort。Linux 子进程 exec 即死 ⇒ `cmd.Start` 本身报错（hook 不可达）或 hook 内 `ProcessAlive` 为 false ⇒ hook 返回错误 → abort。
  - **attach 失败**：Windows 用 Job-at-creation 后 attach 退化为**验证**；验证失败=创建即错 → `platform.terminate()` + `abortStartup`。Linux pgid 不符（防御性）同样 terminate+abort。
  - **OnSpawn 自身失败**：进程仍未 resume ⇒ 语义「未启动」；guard 执行 terminate（杀挂起树）+ `abortStartup`（`process.go:155`，有界 reaper + 不确定保留）。**hook 契约**：失败必须自清理已写 marker（写后验、失败先删后返回），guard 不代删（guard 不拥有 adapter 的 durable 命名空间）。
  - **并发 stop**：starting 已注册 ⇒ 可被 stop；stop 与 resume 的竞态由「resume 前 identity 已确定」消解——两个动作操作同一 Job/pgid，先到先赢，失败侧得 `already-stopped`。

### Q2. 归属先于 resume 的结构保证

**断裂根因**：Windows 的归属动作（assign）与 resume 共居 `attach` 且都在 Start 后；Linux 无 suspend、无归属锚。**结构修复 = 把「归属」从 attach 期前移到创建期**：

- Windows：`preparePlatformProcess` 改用 `STARTUPINFOEX + PROC_THREAD_ATTRIBUTE_JOB_LIST`，进程**创建瞬间即在 Job 内**（`KillOnJobClose` 同步生效）；`attach` 退化为验证 assign 结果；新增独立 `resume()`（`NtResumeProcess`）。⇒ 任意时刻宿主崩溃，OS 关句柄即杀挂起树，**窗口按构造消失**。
- Linux：`prepare` 保留 `Setpgid`，新增 `Pdeathsig: SIGKILL`（Go 1.22 支持），父线程死亡即内核击杀；`resume()` 为 no-op。
- **状态机顺序**（创建→绑定→持久登记→resume）：`StartManaged` 内硬编码顺序，resume 是唯一不可达（OnSpawn 成功后）的下一步；hook 阻塞时进程保持挂起（Windows），即「marker 先于首条用户指令」。
- **测试钉住点**：红测用 hook 阻塞 + 子进程首动作写 marker 文件——hook 未放行前断言子进程 marker 不存在（从未执行）；放行后出现。Windows 另用 kill-on-close 实证「hook 期间宿主被杀 ⇒ 挂起子进程被内核击杀」。
- **与 P40 衔接**：`abortStartup`（`process.go:155`）原样复用为 OnSpawn/resume/attach 三类失败的统一出口（有界 wait + `ErrTerminationUncertain` 保留句柄 + 后台 reaper 收尾）；`Reclaim`（`:377`）、`RunManaged`、`Run` 零改动；新增错误包一层可 `errors.Is` 的哨兵（见 §三）。

### Q3. Job 归属前宿主崩溃的回收

**识别**：重启侧（`claimLaunchSlot`/`refuseLiveMirrorProcess` 前，新增 `RecoverStarting`）扫描 run root 的 `managed-process.starting.json`。判定规则：
- marker 存在 + LockFileEx 锁**被持有**（Windows）⇒ 宿主仍活 → 与今日新鲜度拒绝同口径。
- marker 存在 + 锁已释放 ⇒ 宿主已死：
  - **Windows**：Job-at-creation 已保证子进程随句柄关闭被杀 ⇒ 进程必已死（`ProcessAlive` 复核）→ **直接清理** marker/slot/镜像，立即放行新启动（不再等 10s 新鲜度）。
  - **Linux**：`ProcessAlive` 为真（孤儿仍在跑）→ **领养-击杀**：以 marker 身份经 `StopProcessIdentity`（pgid 杀树，`process_linux.go:70`），`StartToken` 比对防 PID 复用误杀；击杀验证后清理。
  - 两平台 marker 身份 `StartToken` 不匹配（PID 已复用）⇒ 不杀、不清镜像，只清 marker+slot（fail-closed 留痕）。
- **清理槽位**：复用 `releaseLaunchSlot` 的 rename-away 协议（`pinned_cli.go:409`），token 校验防抢新槽。
- **N1 最坏界**：Windows **零窗口**——崩溃即由内核杀，重启首次扫描即知「已无进程」；Linux 若不引入 `Pdeathsig`，孤儿持续运行至重启后首次扫描才判死，**最坏延迟 = 重启间隔**（无自愈能力）；引入 `Pdeathsig` 后同为崩溃即杀。口径：Windows=Job 归属（创建期）+ 锁活体；Linux=pgid 归属 + 启动扫描领养击杀（`Pdeathsig` 为可选加固，落地需实机验证）。

## 三、关键签名与状态转移表

**签名（≤3）**：
```go
// 1. guard：钩子挂载（其余字段不动）
type ProcessSpec struct { /* 既有字段 */; OnSpawn func(ctx context.Context, identity ProcessIdentity) error }
// 2. guard 平台内部：分离归属与运行
func (p *platformProcess) attach(pid int) error // Windows: 验证 job-list 归属；Linux: 记录 pgid
func (p *platformProcess) resume() error        // Windows: NtResumeProcess；Linux: no-op
// 3. adapters：重启回收（claimLaunchSlot 前调用）
func (l *AgentRunnerPinnedCLILauncher) RecoverStarting(ctx context.Context) (reclaimed bool, err error)
```
新增哨兵：`guard.ErrSpawnOwnership`（`errors.Is` 可识别，包装 cause）。

**状态转移表（≤6 行）**：

| # | 现态 | 事件/条件 | 次态 | 平台锚 |
|---|---|---|---|---|
| 1 | 无（slot 已 claim，pre-spawn） | `cmd.Start` 成功，Job/pgid 创建期绑定完成 | `starting`（未持久，挂起/运行中未 resume） | Job 已含挂起树 / pgid 已设 |
| 2 | `starting`（未持久） | OnSpawn 成功（marker 落盘 + registry 注册 starting） | `starting`（durable）→ 允许 resume | marker = 启动树身份证明 |
| 3 | `starting`（durable） | resume 成功 + identity mirror 写 + `advance(running)` | `running` | 归属与身份双重可证明 |
| 4 | `starting` | OnSpawn/attach/resume 失败或进程已死 | abortStartup → terminate → 清 marker/slot → 无 | 无残留 |
| 5 | `starting`（durable，宿主已崩） | 重启扫描：Windows Job 已杀 / Linux 孤儿存活 | 回收：验身份 → Linux 领养击杀 → 清 marker+slot+mirror | Windows 免杀、Linux 杀 pgid |
| 6 | `running` | 停机路径 | `stopping→stopped`（release） | P40 语义不变 |

## 四、红测不变式清单（③ 必须钉住，可证伪）

1. **I1 归属先于 resume（Windows）**：OnSpawn 阻塞期间，子进程首条用户指令不可执行（子 marker 缺失），且此时 Job 归属已成立——用 kill-on-close 反证：宿主测试进程死于 hook 期间 ⇒ 挂起子进程被内核击杀、零残留 PID。
2. **I2 崩溃无孤儿（全平台）**：在「Start 后→attach→OnSpawn→resume」四缝各注入一次宿主崩溃（测试子进程扮演宿主）⇒ 重启侧扫描后目标 PID 一律不存在（Windows Job 杀 / Linux Pdeathsig 或领养击杀）。
3. **I3 starting 的 durable 证据**：宿主在 OnSpawn 返回前崩溃 ⇒ 磁盘 marker 完整（原子写）且 identity mirror **不存在**（未达 running）⇒ 重启侧可区分「starting 崩溃」与「从未启动」。
4. **I4 失败回滚**：OnSpawn 返回错误 ⇒ 无残留进程、槽位释放、marker/mirror 均无、错误 `errors.Is` 可识别，且后续 launch 可立即复用槽位。
5. **I5 回收幂等 + 防误杀**：对同一 marker 二次 `RecoverStarting` 为 no-op；PID 已被复用（StartToken 不符）时绝不击杀新进程、mirror 保留 fail-closed。

## 五、风险与反例

1. **Job-at-creation 的代价**：STARTUPINFOEX 属性列表需 `syscall` 级扩展（`procThreadAttributeJobList`），Go `os/exec` 无直接支持——须在 `preparePlatformProcess` 手写 `STARTUPINFOEXW`+`UpdateProcThreadAttribute`。**推翻点**：若目标 Windows 的 Job 归属与 `CREATE_SUSPENDED` 组合实测不可共存（属性不生效时进程照常运行），方案退化为「Start 后立即 assign + 崩溃测试证残余窗口」，Windows 维持 unproven。
2. **Linux `Pdeathsig` 语义**：只对「创建它的线程」死亡触发；宿主多线程崩溃可能不触发 ⇒ 只能作加固，领养-击杀仍是主回收路径。反例：fork 后宿主线程迁移（Go runtime 无此情形，但需注释登记）。
3. **OnSpawn 钩子契约漂移**：hook 内做重活（阻塞 I/O）会拉长挂起窗口；约定 hook 只写 marker+注册，超时由调用方 ctx 约束，否则挂起树成为新型滞留。任何「把 resume 提前到 hook 前」的改动直接推翻本方案。
4. **回收扫描的误判**：marker 存在但宿主**未死**（锁被持有）却因锁实现缺陷被误读为可回收 → 双进程同 run root。LockFileEx 语义必须红测覆盖「宿主存活时扫描必须拒绝」。
5. **兼容回归**：`agent_runner_pinned_cli_test.go` 大量依赖今日「spawn 后注册 running」与 10s 新鲜度（`:326,692,912` 等）；时序前移会触发既有测试重排，属预期改动，须全量 `go test -count=1 ./...` 红绿双验。

## 六、与 P42（EVIDENCE-BIND）衔接点

- **不重复建设**：P41 的 marker 写入/回收读取**只用既有原语**——identity 经 `platformIdentity`（单次句柄绑定，P40 已审）、liveness 经 `ProcessAlive`（StartToken 比对）；不新写「按路径读 PID→重开路径」模式，不新增目录绑定机制。
- **登记新文件**：`managed-process.starting.json` 作为 run root 新工件**登记入 P42 的目录绑定/ACL 覆盖清单**——P41 只声明路径与原子写协议（temp+rename+LockFileEx），其「目录绑定延续到实际 I/O」（N2）的收紧由 P42 完成。
- **边界**：P41 不碰 `internal/evidence/**`、不碰 `RunEvidenceDir` 语义；若 P42 结论要求 marker 读写改经绑定句柄，P41 只改调用点不改协议。

`DESIGN: DONE`

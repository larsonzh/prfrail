# D3 顺序反转与回归方案（P48 DR13-FIX-D3 · ① 架构 · 只读载体）

## 1. 结论摘要

**最终推荐顺序：`os.Remove(claim.path)` → `unlockStartingMarker(claim.file)` → `claim.file.Close()`（unlock 与 Close 无条件执行）**。
一句话理由：把「路径消失」严格前移到锁释放之前，使扫描者不再可能观察到「锁空闲 + 路径仍在」的 W1 态；「删除失败仍解锁关闭」是唯一可实现且 fail-closed 的失败方向，残留 marker 由 D1/D2 与下一次扫描按既有证据收敛。

## 2. 逐条结论

### A. 失败组合安全性（fail-closed）

现序证据：`internal/adapters/agent_runner_starting_marker.go:92-101`（unlock :96 → Close :97 → os.Remove :99）。

**(1) 常态（remove 成功）**：路径在锁仍被 claim 句柄持有时消失。claim 句柄由 `createLockedStartingMarkerFile` 以 `startingMarkerShareRead|startingMarkerShareWrite|startingMarkerShareDelete` 创建（`agent_runner_starting_marker_windows.go:146`，常量定义 :31），Windows 上本进程句柄允许删除；unix 上 unlink 已打开文件恒允许。remove 后 unlock 针对的是已脱离目录的同一文件对象（Windows）/ 同一 inode（flock 随 inode 保留，`agent_runner_starting_marker_unix.go` 头部注释），语义正确；最后 Close 兜底。

**(2) remove 失败**（并发 scan 已清 / 共享冲突 / 权限）：**锁仍必须且无条件释放**。理由：① 「为保锁而不 Close」必然泄漏文件句柄，而两平台的锁都由内核在 Close 或进程死亡时释放（`agent_runner_starting_marker_windows.go:68-70` 注释「Closing the file releases it as well」；unix 头部 flock 注释「released by the kernel on any kind of host death」），「remove 失败后仍持锁」不可实现且无意义；② 残留的「锁空闲 + 路径仍在」marker 走既有收敛路径：同进程被 D1 门拦截（`agent_runner_pinned_cli.go:390-397`），跨进程由下一次扫描的 D2 复核与 liveness 判定处置——这正是 DR13-FIX 已登记的 R5 残留口径（`docs/validation/dr13-fix.md` §偏差与登记 R5），本片收窄、不消除。**推荐最终代码顺序：`remove → unlock → Close`**，unlock 与 Close 放在同一无条件路径（无 defer 分支）；`release()` 返回值语义不变（仍为「是否实际删除」bool），现有调用方全部忽略返回值（`agent_runner_process_registry.go:82-87`、`agent_runner_pinned_cli.go:689/702`），零契约变化。

**(3) unlock 失败**：当前实现忽略返回值（`agent_runner_starting_marker_windows.go:71-78` 三下划线）。**不需改**：unlock 仅是显式形式，权威保证是 Close/进程死亡的内核释放；新顺序中 Close 紧随 unlock 无条件执行，unlock 报错也不会留下「幽灵锁」（锁不依附于路径、只依附于句柄/inode）。

**(4) 安全方向**：「删除失败但锁已释放」是安全且唯一可实现的方向——失败态降级为「可被下次扫描按既有证据收敛的残留」，不泄漏句柄、不阻塞 run root；「删除成功但锁未释放」在锁随句柄/进程存续的语义下不可持久存在，是伪命题。真正要防的是「锁释放早于路径消失」，这正是本片要消灭的顺序。

### B. Windows remove-while-open 语义

- Windows 删除语义自 NT 起统一：DeleteFile 在存在任何不含 FILE_SHARE_DELETE 的打开句柄时返回共享冲突，Win7/10/11/Server 无版本差异。claim 句柄共享模式含 `startingMarkerShareDelete`（`agent_runner_starting_marker_windows.go:146`），与 `writeAgentRunnerStartingMarker` 的「锁内 rename」共用同一前提（`:125-133` 注释「share mode that stays renameable」）。在仓先例已实证「持锁删除」：`clearAgentRunnerStartingMarker` 就是「持锁 os.Remove → unlock → Close」（`agent_runner_pinned_cli.go:570-574`），其句柄经 `openStartingMarkerForClear` 显式补 delete 共享（`agent_runner_starting_marker_windows.go:102-122`，share 参数 :110）。
- **一个须回归实证的事实分歧点**：仓库注释称 Go 的 `os.OpenFile` 不含 delete 共享（windows.go:97-101 注释），而 `RecoverStarting` 用 `os.OpenFile` 打开路径（`agent_runner_pinned_cli.go:341`）。若属实：扫描者持柄期间 release 的 remove 会失败 → 落入 A(2) 收敛；若注释过时（Go 某版本起 `os.OpenFile` 已含 FILE_SHARE_DELETE）：remove 成功、扫描者可能锁到孤儿 inode——恰由 D2 的 Stat 复核兜住（`:366-377`）。两条路径都闭合，但回归必须显式钉住本工具链实际行为（见 §3）。
- **核心不变量成立**：remove-before-unlock 结构性保证「路径消失时锁仍被本进程持有」；Close 在 unlock 之后只影响句柄，不影响路径存在性。
- 跨版本差异仅在环境层（AV/过滤器/WSL 可瞬态制造共享冲突），非版本语义差异；覆盖方式 = Windows CI 腿 + 本机 + 显式共享冲突反测。

### C. 与 D2 / D1 的协同

新顺序不使 D2 失效、也不放宽 D1：
- D2（`agent_runner_pinned_cli.go:366-377`）针对「扫描者锁到已 unlink 孤儿 inode」。D3 后孤儿 inode 的产生点收窄为「扫描者 open 先于 release remove、tryLock 后于 remove」这一窗口，Stat 复核仍在其后无条件执行，拦截力不变；不存在 D3 使 Stat 复核被跳过的路径。
- D1（`:390-397`）同进程所有权门与 marker 顺序无关，仍拦截本进程在册 launch。
- 组合不变式：claim 自身释放不再产生「锁空闲 + 路径仍在」；扫描者侧两个残余窗口（孤儿 inode、remove 失败残留）分别由 D2、D1+下次扫描覆盖。**未发现 D3 引入的新绕过路径。**

### D. 回归方案

**(1) 红测设计**：纯测试层**无法**确定性卡住 unlock↔remove 间隙——`release()` 无注入边界，`unlockStartingMarker` 是普通包级函数（测试文件不能重定义），tryLock 双平台非阻塞（与 DR13-FIX 已登记的 B2/R3 同类结论，`docs/validation/dr13-fix.md` R3）。可证伪替代（按推荐序）：
- **结构性红测（首选，确定性、零生产改动）**：`TestStartingMarkerReleaseRemovesBeforeUnlock` 用 `go/parser` 解析 `agent_runner_starting_marker.go`，断言 release 函数体内 `os.Remove(claim.path)` 语句先于 `unlockStartingMarker(claim.file)`。旧代码确定性红、新代码绿；诚实声明为「源码顺序门」而非运行时证明。
- **fail-closed 语义红测（确定性，Windows-only）**：测试用 raw `CreateFile` 以「read|write 共享、无 delete 共享」打开已发布路径（第三方占名句柄），再调 `release()`：断言返回 false、路径仍在，但随后 `openStartingMarkerForClear` + `tryLockStartingMarker` 必须成功——证明「remove 失败时锁确已释放」。
- **压力回归（概率型，登记口径同 R3）**：并发循环「release + 重写 claim」×「模拟扫描（open+tryLock+Stat+decode）」N 次迭代，断言扫描者从不观察到「tryLock 成功 + Stat 成功 + 内容属于已被 release 的 claim」。旧代码该态概率非零、新代码结构性为零；仅作回归补充。

**(2) Windows 全量回归范围**（`go test -count=1 ./internal/adapters`，另跑 ubuntu 腿）：
- 必跑既有：`TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner`（`agent_runner_pinned_cli_test.go:857`）、`StartsAndStops`（:224）、`RefusesASecondLiveLaunch`（:278）、`RefusesASecondInstanceWhileTheMirrorIsLive`（:331）、`StopGraceFollowsTheLaunch`（:380）、`StopNeverReachesAnotherLaunch`（:431）；`TestRecoverStartingIdempotent`（`spawn_recovery_test.go:138`）、`TokenMismatchFailsClosed`、`NeverReclaimsALaunchThisProcessStillOwns`（:221）；`TestRecoverStartingCleansCrashedHostWindows`（`spawn_recovery_live_host_windows_test.go:75`）、`RefusesLiveHost`（:124）；unix 腿 `RefusesLiveHostUnix`（`spawn_recovery_live_host_unix_test.go:118`）、`AdoptsOrphanUnix`（:160）。
- 新增断言：① 结构性顺序断言；② release 后「路径 NotExist ∧ 返回 true」的 remove-while-open 直测（钉住 share 模式前提，防回归改坏共享参数）；③ double-release 幂等（返回 false、无 panic）；④ 共享冲突反测（上述 fail-closed 测试）；⑤ 可选：用 `os.OpenFile` 复刻 `RecoverStarting:341` 的打开方式后验证 `os.Remove` 是否返回共享冲突，钉住 §B 的事实分歧点。

**(3) 结构性红测若视为不足的备选**：最小生产 seam——把 `unlockStartingMarker` 改为包级 var（一行生产改动，超出 S 档边界，需 Sol 显式授权并登记偏差），测试替换桩在 remove 与 unlock 之间阻塞，构造确定性运行时红测；其可证伪性在于旧代码不存在该插入点、桩无法制造「路径已删 + 锁仍持」态。默认不采用。

### E. 边界与风险

**不该动**：DR-13 定性结论与 `winners==1` 断言（`agent_runner_pinned_cli_test.go:901-902`，其语义即协议不变式，DR-13 已裁定不可放松）；`docs/CONTRACTS.md` 与 Schema；D1 门与 D2 复核；`clearAgentRunnerStartingMarker` 既有「锁内 remove」顺序（已是正确顺序，`:570-574`）；`writeAgentRunnerStartingMarker` 的 discard 路径（`agent_runner_starting_marker.go:74-78` 操作的是 temp 名，扫描不可见，不属 W1）；`unlockStartingMarker` 签名与返回值忽略；`closeStartingMarker`/`releaseLaunch`/`discardStartingLaunch` 调用结构（`agent_runner_pinned_cli.go:679-690`/`:697-704`）。
**新引入风险**：① remove 失败残留的跨进程 W1 口径不变（R5），本片只收窄不消除，须如实标注；② 若 §B 事实分歧点与注释相反，孤儿 inode 窗口略增、由 D2 吸收——回归必须显式钉住；③ `release()` 返回语义不变，调用方忽略，无行为漂移；④ 无新竞态——新顺序中间态「路径已删 + 锁仍持」对扫描者等价于「无 marker」（open 得 ErrNotExist 走 `:346-350` no-op），不会诱发误回收。

## 3. 回归测试清单（汇总）

| 类型 | 测试 | 平台 | 可证伪性 |
|---|---|---|---|
| 结构性红测 | release 体内 Remove 先于 unlock（AST） | 双平台 | 确定性 |
| 语义直测 | release → 路径 NotExist ∧ 返回 true（remove-while-open） | 双平台 | 确定性 |
| 失败方向 | 无 delete 共享第三方句柄 → release false 且锁已释放 | Windows | 确定性 |
| 幂等 | 二次 release false、无 panic | 双平台 | 确定性 |
| 共享事实 | `os.OpenFile` 打开后 `os.Remove` 行为钉桩 | Windows | 确定性 |
| 压力 | release × 扫描并发 0 违例 | 双平台 | 概率型（R3 口径） |
| 既有回归 | §2D(2) 全部既有测试 + `-race` | 双平台 | 既有 |

## 4. 风险与边界清单

- 不做：DR-13 结论/断言、CONTRACTS/Schema、D1/D2、clear 路径既有顺序、discard 路径、unlock 签名与返回值处理。
- 残留（如实登记）：remove 失败时跨进程 W1 残留（R5 口径）；W2 flock 语义不受本片影响亦不闭合；AV/过滤器瞬态共享冲突无法在 CI 确定性覆盖。
- 需实证：Go `os.OpenFile` 在本工具链的 share 模式真伪（`agent_runner_starting_marker_windows.go:97-101` 注释）。

## 5. 起跑两问

- 「os.Remove 失败时锁是否仍能释放？」——**能且必须**：Close 在两平台由内核保证释放锁，release 必须无条件 unlock+close；任何「为保锁而不 Close」都泄漏句柄且不可实现，残留 marker 由 D1/D2 与下一次扫描收敛。
- 「Windows remove-while-open 在不同 OS 版本上的差异？」——**共享模式判定自 NT 起在 Win7/10/11/Server 语义一致**（任一缺 FILE_SHARE_DELETE 的句柄即拒绝删除）；差异仅在环境层（AV/过滤器/WSL 瞬态冲突），且 claim 句柄带 delete 共享的「持锁删除」先例已由 `clearAgentRunnerStartingMarker` 在仓落地验证。

# ⑥ 独立扫描（原样归档）

> 载体：runSubagent(agentName="independent-reviewer", model="MAI-Code-1.1-Flash (copilot)")。偏移 120306。

## A 结论

INDEPENDENT SCAN: PASS

## B 缺陷

| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
|---|---|---|---|---|
|  |  |  |  |  |

## C 核验结论

1. D1 的正确性与充分性：成立。`registry.lookup` 是进程内唯一来源的“同 launcher 仍持有该 launchID”判定，且持有 `registry.mu` 读出句柄后不做任何写操作；同进程所有权判断位于 `RecoverStarting` 的 decode 成功后、liveness 判定前，正好覆盖“锁空闲但同进程仍持有启动句柄”的关键窗口 [internal/adapters/agent_runner_process_registry.go](internal/adapters/agent_runner_process_registry.go#L141-L149) 和 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L326-L399)。成功路径上，`StartAgentRunnerProcess` 只在 `advance(..., running)` 后调用 `handle.closeStartingMarker()`，并不会在成功返回前 `release` 该 handle；因此，注册表中的活句柄确实仍在，并在同进程竞态下拦截误回收 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L254-L293)。  
   结论：D1 不仅正确，而且对当前修复边界是充分的；它锁定了“同进程 handle 仍在”的死角，而不是放大到跨进程 orphan 的正常回收。

2. D1 是否误伤孤儿收养：成立。`registry` 纯内存，仅存在于当前进程；崩溃后它必然消失，所以跨进程 orphan 不会命中 D1。相反，正路径仍保留：`TestRecoverStartingIdempotent`、`TestRecoverStartingTokenMismatchFailsClosed`、`TestRecoverStartingRefusesLiveHostUnix`、`TestRecoverStartingAdoptsOrphanUnix` 继续覆盖“已死主机/TokenMismatch/活主机/孤儿收养”的正确分支；D1 只在 same-process ownership 上做 fail-closed—且这正是其目标 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L135-L210) 和 [internal/adapters/spawn_recovery_live_host_unix_test.go](internal/adapters/spawn_recovery_live_host_unix_test.go#L113-L198)。  
   结论：D1 不会误伤 orphan 收养；它只拦截“本进程仍持有该 launch 的句柄”。

3. D2 的正确性：成立。`os.Stat(path)` 在 `tryLock` 成功后立刻做；若 `ErrNotExist`，说明锁已被取到的 inode 已脱离当前 published path，属于`orphan inode` 的典型场景，必须返回 `(false, nil)` 并清理 lock/FD，而不是继续 decode 或 kill [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L357-L377)。这正是把“锁空闲但 marker 不再存在”识别为“无可回收记录”的 fail-closed 语义，不构成误判；它并不掩盖真实错误，因为只有 `NotExist` 才是此路径的合法闭环，而其它 I/O 错误仍会返回 `ErrAgentRunnerLaunchConflict`。  
   结论：D2 对“orphan inode”窗口是正确的，且对调用方不构成误判。

4. E1 的可证伪性：成立。新增测试 `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` 使用了真实 helper 进程、真实 `claimLaunchSlot`、真实 marker 文件和真实 registry 注册，断言四件事全部落在“同进程仍 owns launch”这个对错分界上：`reclaimed == false`、helper 继续活着、marker 仍在、slot owner 仍为 launchID [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
   这不是假绿通道：若去掉 D1 门，旧实现会落入“live process + free lock”的回收分支，触发 `StopProcessIdentity` 杀活进程或删 slot；这正与 `reclaimStartingProcess`/`reclaimLaunchSlot` 分支的语义一致 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L406-L417)。  
   结论：该测试是可证伪的，并且确实覆盖了真正的回归窗口。

5. 并发回归：成立。`TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner` 在同一 launcher / 同一 run root 下保留原逻辑；其设计目标是确保恰一名赢家，且其余竞争者返回 `ErrAgentRunnerLaunchConflict` 或 `ErrAgentRunnerProcessDuplicate` [internal/adapters/agent_runner_pinned_cli_test.go](internal/adapters/agent_runner_pinned_cli_test.go#L857-L920)。在当前实现下，适配器包测试已通过，说明这是结构性约束，而不是偶然性；并发 winner 语义依赖的仍是“slot O_EXCL”与 registry ownership gate，不依赖测试抖动。  
   结论：`winners == 1` 的结论结构性成立。

6. 边界与卫生：成立。允许路径仅限 `internal/adapters/agent_runner_pinned_cli*.go` 与其测试文件；证据包则落于 `docs/validation/evidence/DR13-FIX/`，并无新增/删改或临时脚本落盘；代码中无 TODO，且新增断言均为英文说明 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
   只读验证证据：`go test ./internal/adapters` 输出为 `ok github.com/larsonzh/prfrail/internal/adapters 35.127s`；`go vet ./internal/adapters` 输出为空（无 vet issue）；`gofmt -l` 也未输出任何文件，说明当前代码已格式化。  
   证据包自洽：`MANIFEST.md` 列出的三个载荷与 `SHA256SUMS.txt` 一致，且对应文件均存在 [docs/validation/evidence/DR13-FIX/MANIFEST.md](docs/validation/evidence/DR13-FIX/MANIFEST.md) 与 [docs/validation/evidence/DR13-FIX/SHA256SUMS.txt](docs/validation/evidence/DR13-FIX/SHA256SUMS.txt)。  
   结论：卫生与证据三方一致，扫描通过。

## D 未覆盖面

1. `StopAgentRunnerProcess` 的失败分支仍能在“输入解析/停止失败”时无条件 `releaseLaunch`; 这不等于 D1 失效，但它是下一个最值得加固的边界，尤其在停止中途状态不一致时可能清掉 `registry` 句柄 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L271-L301)。  
   重要性：它是 stop/start 竞态的镜像窗口，若下一轮再复现 stop-while-starting，可能出现“已杀/未杀状态模糊”。

2. D2/TOCTOU 仍依赖平台层对 `flock`/inode 语义的一致性；本地仅能读代码而不能在 unix CI 上直接探针 `flock` 的行为差异，因此 W2 的内核级细节未被本机独立实证 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L343-L377)。  
   重要性：它关乎“跨进程锁释放是否会让另一个扫描持到 orphan inode”的真实发生率。

3. 证据包的 `01-diagnosis` 是已归档的定性结论，但并未在本机 Unx 宿主上复现 `StopProcessIdentity` 对活进程的具体行为；当前的证据仍依赖 CI/设计推演 [docs/validation/evidence/DR13-FIX/01-diagnosis/review.md](docs/validation/evidence/DR13-FIX/01-diagnosis/review.md)。  
   重要性：这里是最终发布前最需要的真实 runtime 证据，尤其是 ubuntu 腿。

4. D1 仅覆盖“同一 launcher 仍持有同一 launchID”的漏洞；它不覆盖“跨 launcher、跨进程、同一 runRoot 但镜像已失效”的不同窗口，后者仍需要 slot/镜像一致性检查和更严格的 stop 语义 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L620-L700)。  
   重要性：这类窗口不会触发 D1，但仍可能导致不完全的“owner 检查”。

## E 可证伪性

可证伪性成立：当前修复的核心是两层“先证实 marker 仍指向真实 published path，再确认 launcher.registry 没有同 launchID 活句柄”。这两个门都可以被单步反证：

- 若删掉 D1 门，`TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` 会直接失败，且 helper 进程会被误杀或 slot 被删除 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
- 若删掉 D2 门，`tryLock` 成功后到 `os.Stat` 之间的 orphan inode window 就能导致误回收；这在旧实现下会在“marker path 已 unlink”时制造 false positive [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L343-L377)。  
- 若没有真实 helper 进程、真实 slot claim、真实 marker 写入，测试就会失去证伪性；当前测试恰好使用了实际的 live process 和 slot 创建，不存在假绿通道。  
- 这一点也与 `go test` 结果一致：当前实现通过 package 级回归测试，说明反例是可复现且被封住的。

§A INDEPENDENT SCAN: PASS
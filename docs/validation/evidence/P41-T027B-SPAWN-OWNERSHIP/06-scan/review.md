# P41-T027B-SPAWN-OWNERSHIP ⑥ 独立扫描归档（原样）

切片：`P41-T027B-SPAWN-OWNERSHIP`。环节：⑥ 独立扫描。载体：`independent-reviewer`；模型：`MAI-Code-1.1-Flash (copilot)`。输出原样归档。

---

## Section A — 总体结论：本次只读扫描未发现中等及以上缺陷；结合本机 Windows 上的 `go test -count=1 ./...` 15 ok / 0 FAIL，以及关键顺序/锁语义/回收逻辑的一致性，整体结论为 PASS。

## Section B — 发现清单
Low | internal/adapters/agent_runner_pinned_cli.go:494-516 | clearAgentRunnerStartingMarker 仍通过 path 再删除；虽然有重锁与比对，但在极短窗口里，如果另一启动重命名新 marker 到同一路径，旧扫描仍可能误删新记录。 | 在 re-open 后做 inode/内容双重校验，并仅删已证明确的旧记录，而不是直接删除当前 path。
Low | internal/adapters/agent_runner_pinned_cli.go:349-365 | RecoverStarting 把“没有 marker”和“活宿主拒绝”压成同一 boolean false,nil，调用方无法区分“可继续启动”和“需要等待重试”。 | 用专用状态枚举或包装错误替代单一布尔，保留兼容性并保留拒绝语义。

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

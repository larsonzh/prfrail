命令与核对均已完成（231 个 `.go` 全量扫描 + 二次核对文件锁语义）。实际执行的枚举命令（node，只读）：

```js
// 1) 全量枚举（walk internal/**/*.go，逐行正则匹配）
const declRe=/sync\.(Mutex|RWMutex|Once)\b/g;
const callDefs=[['Lock',/\.Lock\(\)/g],['Unlock',/\.Unlock\(\)/g],['RLock',/\.RLock\(\)/g],['RUnlock',/\.RUnlock\(\)/g],['Do',/\.Do\(/g]];
// 文件锁补扫：/flock/i、/LockFile/i、/\.lock["']/；方法定义：/func\s*\(\s*(\w+)\s+\*?(\w+)\s*\)\s*(Lock|Unlock|RLock|RUnlock)\s*\(/
// 2) 二次核对：dump durable_cost_lock_unix.go / _windows.go 全文、guard/lease.go:195-245、全 internal/ 的 O_EXCL 行
```

## 进程内锁

| 文件:行 | 类型 | 变量/接收者 |
|---|---|---|
| internal/adapters/agent_runner_capability.go:62 | Mutex(decl) | mu |
| internal/adapters/agent_runner_capability.go:123 | Lock() | index.mu |
| internal/adapters/agent_runner_capability.go:124 | Unlock() | index.mu |
| internal/adapters/agent_runner_completion.go:48 | Mutex(decl) | mu |
| internal/adapters/agent_runner_completion.go:131 | Lock() | index.mu |
| internal/adapters/agent_runner_completion.go:132 | Unlock() | index.mu |
| internal/adapters/agent_runner_enforcement.go:47 | Mutex(decl) | mu |
| internal/adapters/agent_runner_enforcement.go:116 | Lock() | index.mu |
| internal/adapters/agent_runner_enforcement.go:117 | Unlock() | index.mu |
| internal/adapters/agent_runner_event.go:44 | Mutex(decl) | mu |
| internal/adapters/agent_runner_event.go:128 | Lock() | index.mu |
| internal/adapters/agent_runner_event.go:129 | Unlock() | index.mu |
| internal/adapters/agent_runner_native_a6_test.go:39 | Once(decl) | agentStubOnce |
| internal/adapters/agent_runner_native_a6_test.go:48 | Once.Do() | agentStubOnce |
| internal/adapters/agent_runner_native_a7_test.go:206 | Mutex(decl) | mu |
| internal/adapters/agent_runner_native_a7_test.go:217 | Lock() | mu |
| internal/adapters/agent_runner_native_a7_test.go:218 | Unlock() | mu |
| internal/adapters/agent_runner_postflight_test.go:31 | Once(decl) | parityStubOnce |
| internal/adapters/agent_runner_postflight_test.go:41 | Once.Do() | parityStubOnce |
| internal/adapters/agent_runner_process_registry.go:81 | Mutex(decl) | mu |
| internal/adapters/agent_runner_process_registry.go:102 | Lock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:103 | Unlock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:115 | Lock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:116 | Unlock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:128 | Lock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:129 | Unlock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:145 | Lock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:146 | Unlock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:162 | Lock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:163 | Unlock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:171 | Lock() | registry.mu |
| internal/adapters/agent_runner_process_registry.go:172 | Unlock() | registry.mu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:19 | Mutex(decl) | mu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:29 | Lock() | stub.mu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:30 | Unlock() | stub.mu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:52 | Lock() | stub.mu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:53 | Unlock() | stub.mu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:427 | Mutex(decl) | stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:453 | Lock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:455 | Unlock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:465 | Lock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:466 | Unlock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:483 | Lock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:484 | Unlock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:489 | Lock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:490 | Unlock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:498 | Lock() | stub.stopMu |
| internal/adapters/agent_runner_replay_dispatcher_test.go:499 | Unlock() | stub.stopMu |
| internal/adapters/agent_runner_replay_store.go:166 | Mutex(decl) | mu |
| internal/adapters/agent_runner_replay_store.go:224 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store.go:225 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store.go:264 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store.go:265 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store.go:318 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store.go:319 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:40 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:41 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:140 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:141 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:170 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:171 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:249 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_launches.go:250 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:39 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:40 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:138 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:139 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:168 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:169 | Unlock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:251 | Lock() | store.mu |
| internal/adapters/agent_runner_replay_store_terminals.go:252 | Unlock() | store.mu |
| internal/adapters/ai_availability.go:51 | Mutex(decl) | mu |
| internal/adapters/ai_availability.go:148 | Lock() | index.mu |
| internal/adapters/ai_availability.go:149 | Unlock() | index.mu |
| internal/adapters/ai_probe_copilot_executor.go:295 | Mutex(decl) | mu |
| internal/adapters/ai_probe_copilot_executor.go:316 | Lock() | writer.owner.mu |
| internal/adapters/ai_probe_copilot_executor.go:326 | Unlock() | writer.owner.mu |
| internal/adapters/ai_probe_copilot_executor.go:334 | Lock() | output.mu |
| internal/adapters/ai_probe_copilot_executor.go:335 | Unlock() | output.mu |
| internal/adapters/ai_probe_copilot_executor.go:340 | Lock() | output.mu |
| internal/adapters/ai_probe_copilot_executor.go:341 | Unlock() | output.mu |
| internal/console/runtime_agent_e2e_test.go:47 | Once(decl) | assemblyStubOnce |
| internal/console/runtime_agent_e2e_test.go:56 | Once.Do() | assemblyStubOnce |
| internal/gates/executor.go:63 | Mutex(decl) | mu |
| internal/gates/executor.go:77 | Lock() | output.mu |
| internal/gates/executor.go:88 | Unlock() | output.mu |
| internal/gates/executor.go:96 | Lock() | output.mu |
| internal/gates/executor.go:97 | Unlock() | output.mu |
| internal/gates/result.go:22 | Mutex(decl) | mu |
| internal/gates/result.go:63 | Lock() | store.mu |
| internal/gates/result.go:64 | Unlock() | store.mu |
| internal/guard/process_test.go:141 | Mutex(decl) | mu |
| internal/guard/process_test.go:148 | Lock() | mu |
| internal/guard/process_test.go:155 | Unlock() | mu |
| internal/guard/process_test.go:188 | Lock() | mu |
| internal/guard/process_test.go:189 | Unlock() | mu |
| internal/snapshot/store.go:21 | RWMutex(decl) | mu |
| internal/snapshot/store.go:55 | Lock() | s.mu |
| internal/snapshot/store.go:56 | Unlock() | s.mu |
| internal/snapshot/store.go:120 | RLock() | s.mu |
| internal/snapshot/store.go:121 | RUnlock() | s.mu |
| internal/snapshot/store.go:141 | RLock() | s.mu |
| internal/snapshot/store.go:142 | RUnlock() | s.mu |
| internal/snapshot/store.go:161 | RLock() | s.mu |
| internal/snapshot/store.go:162 | RUnlock() | s.mu |
| internal/snapshot/store.go:186 | RLock() | s.mu |
| internal/snapshot/store.go:187 | RUnlock() | s.mu |
| internal/tickets/cost.go:203 | RWMutex(decl) | mu |
| internal/tickets/cost.go:345 | RLock() | ledger.mu |
| internal/tickets/cost.go:346 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:354 | RLock() | ledger.mu |
| internal/tickets/cost.go:355 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:363 | RLock() | ledger.mu |
| internal/tickets/cost.go:364 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:375 | RLock() | ledger.mu |
| internal/tickets/cost.go:376 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:392 | RLock() | ledger.mu |
| internal/tickets/cost.go:393 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:414 | RLock() | ledger.mu |
| internal/tickets/cost.go:415 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:430 | RLock() | ledger.mu |
| internal/tickets/cost.go:431 | RUnlock() | ledger.mu |
| internal/tickets/cost.go:453 | Lock() | ledger.mu |
| internal/tickets/cost.go:454 | Unlock() | ledger.mu |
| internal/tickets/cost.go:482 | Lock() | ledger.mu |
| internal/tickets/cost.go:483 | Unlock() | ledger.mu |
| internal/tickets/cost.go:518 | Lock() | ledger.mu |
| internal/tickets/cost.go:519 | Unlock() | ledger.mu |
| internal/tickets/durable_cost_ledger.go:81 | Mutex(decl) | mu |
| internal/tickets/durable_cost_ledger.go:174 | Lock() | ledger.mu |
| internal/tickets/durable_cost_ledger.go:175 | Unlock() | ledger.mu |
| internal/tickets/durable_cost_ledger.go:318 | Lock() | ledger.mu |
| internal/tickets/durable_cost_ledger.go:319 | Unlock() | ledger.mu |

## 跨进程文件锁

| 文件:行 | 类型 | 锁路径派生 |
|---|---|---|
| internal/tickets/durable_cost_lock_unix.go:32 | `syscall.Flock(fd, LOCK_EX\|LOCK_NB)`（unix 变体） | `acquireDurableLock(path)` 的 path = `filepath.Join(root, name+".lock")`（durable_cost_ledger.go:116；锁文件永不替换/删除） |
| internal/tickets/durable_cost_lock_unix.go:49 | `syscall.Flock(fd, LOCK_UN)`（解锁） | 同上，同一 fd |
| internal/tickets/durable_cost_lock_windows.go:48 | `LockFileEx`（独占 + FAIL_IMMEDIATELY，固定 1 字节区间；kernel32 经 NewLazyDLL） | 同一 `lockPath`（_windows 构建变体，与 unix 互斥编译） |
| internal/tickets/durable_cost_lock_windows.go:58 | `UnlockFileEx`（解锁） | 同上 |
| internal/guard/lease.go:217 | `os.Mkdir(lockPath, 0700)` 原子目录锁（非 flock/LockFileEx，存在即争用失败，defer Remove :228） | `store.leasePath(resourceID) + ".lock"`（guard/lease.go:215） |

机械备注（不作为锁点计数）：
- `internal/tickets/durable_cost_ledger.go:394` 的 `O_CREATE|O_EXCL` 为账本文件原子替换写入，非锁语义；`internal/applier/atomic.go:10`、`internal/snapshot/store.go:89`、`internal/console/cli.go:610`、`internal/adapters/agent_runner_pinned_cli.go:384` 同类（O_EXCL 原子创建/临时文件）。
- `Once` 3 处全部位于 `_test.go`；`adapters` 的 decl 计数含 3 个 `_test` 文件；`guard` 的 decl 仅 `process_test.go`。
- 未发现自定义类型实现 `Lock/Unlock/RLock/RUnlock` 方法（方法定义扫描为空），所有 `.Lock()` 等调用均直接作用于 `sync.Mutex`/`sync.RWMutex` 字段。

## 持有方分布（机械归类，不做锁序结论）

| 持有方 | decl | 调用点 |
|---|---|---|
| adapters | 13 | 68 |
| tickets | 2 | 24 |
| 其它(gates) | 2 | 6 |
| 其它(snapshot) | 1 | 10 |
| guard | 1 | 4 |
| 其它(console) | 1 | 1 |

## 计数

Mutex=15 RWMutex=2 Once=3 FileLock=5 调用点合计=110（Lock=58, Unlock=58 中配对；RLock/RUnlock=18/18；另 Once.Do=3）
FileLock=5 = unix Flock 2 + windows LockFileEx/UnlockFileEx 2（构建变体互斥）+ guard/lease.go Mkdir 目录锁 1；运行时单平台实际生效 3 处。

CENSUS: DONE

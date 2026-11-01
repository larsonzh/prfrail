# ProofRail 验证报告 — DR13-FIX

| 字段 | 内容 |
|---|---|
| 切片 | `DR13-FIX` |
| 类型 | 修复类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `8e834b2` |
| 日期 | 2026-11-01 |

## 目标与边界

- **定位**：本片是 P45 的前置修复片，不占 v34 队列 1–4 的编号；v34 保持「P44 闭环 → P45 起跑前」的基线，本片闭环后由 v35 统一反映（用户 2026-11-01 裁定）。
- **目标**：定性 `TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner` 的 flaky 根因，区分「测试装置过严（A）」与「真实排他窗口漏洞（B）」，并交付可证伪的修复与反例。
- **范围**：`internal/adapters/agent_runner_pinned_cli.go`（排他窗口）与 `internal/adapters/spawn_recovery_test.go`（反例）。
- **边界**：不动其它包（含 `internal/chain/**`、`internal/gates/**`、`internal/console/**`）、schema、fixtures、`ci.yml`；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 定性结论

- **DIAGNOSIS: B**：真实排他窗口漏洞。同进程、同 launcher 并发 Start 时，`RecoverStarting` 的「marker 锁空闲 ⇒ 前宿主已死」证明失效；unix 分支对该「孤儿」执行停进程，杀掉首位赢家的活子进程并夺走其 slot，使两个 Start 先后成功。
- **排除 A**：两次成功之间没有任何合法 slot 释放点被触发——无 Stop 调用，stub 的 `run` 模式睡 60 秒且 launcher 从不 reap，首胜者也不曾失败；slot 在首胜者存活期间消失只能经由 reclaim 误判，而 reclaim 把活进程当崩溃孤儿杀掉，属实现缺陷，不是「断言把跨时间成功数当同一时刻存活数」的装置过严。
- **断言不可放松**：`agent_runner_pinned_cli_test.go:901-902` 的断言语义正是协议不变式本身（slot 原子独占 ⇒ 第二个 Start 在首胜者存活期间必须失败），因此本片修实现、不动断言。
- **证据链摘要**：
  - `agent_runner_pinned_cli.go:124` 的入口门在 claim 之前，8 个 contender 可同时通过（register 要到 OnSpawn 才发生）⇒ 排他只能靠 slot 兜底。
  - `agent_runner_pinned_cli.go:345` 的 `tryLock` 成功后未复核路径，锁可能落在已被 unlink 的孤儿 inode 上。
  - `agent_runner_starting_marker_unix.go` 的 `reclaimStartingProcess` 在 `alive=true` 时执行 `StopProcessIdentity` 并返回 `(true,nil)`——ubuntu 腿杀活根因；`agent_runner_starting_marker_windows.go` 同名函数在 `alive=true` 时返回 `(false,nil)`（fail closed，不杀），这解释了 windows 腿为何不红。
  - `agent_runner_starting_marker.go:89-99` 的 `release()` = unlock → close → `os.Remove`，锁释放严格早于路径消失，扫描者可观察到「锁空闲 + 路径仍在」（W1 窗口）。
  - fail-closed 与完整回收分支都会 `reclaimLaunchSlot` 并清 marker，构成夺槽路径；`agent_runner_pinned_cli_test.go:894-896` 的败者白名单只有 Conflict / Duplicate，与观测到的 `got 2 (losers 6)` 逐字吻合。
- **W1 时序摘要**（不依赖任何特殊 flock 内核语义）：A 抢到 slot 并 spawn → A 的 OnSpawn 发布 marker 并 register → B 在「marker 已发布、A 尚未 `closeStartingMarker`」的毫秒窗口内 open 成功后被抢占 → A 释放锁、`os.Remove` marker、返回成功（stub 仍活）→ B 恢复运行并锁住已被 unlink 的孤儿 inode → decode 成功、`ProcessAlive(A)` 为真 → unix 分支杀掉 A 的活 stub 并夺走其 slot → B 的 claim 见 slot 缺失、镜像已死而放行 ⇒ 赢家 #2；其余 6 人全部命中 Conflict / Duplicate 哨兵。
- **W2 时序摘要**（放大器，内核相关）：B 的 `tryLock` 失败后关闭自己的 FD，按 man flock(2)「进程关闭任一指向该文件的 fd 即释放该进程对该文件全部 flock」的语义，A 的锁被这一关闭动作提前释放，使「锁空闲 + 进程存活」贯穿 A 的整个收尾窗口并显著提高命中率；W1 独立成立，两次 CI 红究竟经 W1 还是 W2 无法从日志区分。
- **修法取舍**：D1 为主修——`RecoverStarting` 在 decode 成功之后、liveness 判定之前加同进程所有权门（本 launcher 在册持有 `marker.LaunchID` 的活句柄 ⇒ 永不回收该 launch），确定性、双平台、不涉协议、不改 CONTRACTS / Schema；D2 为配套——`tryLock` 成功后 `os.Stat(path)` 复核，`os.ErrNotExist` 即锁住的是已 unlink 的孤儿 inode，返回 `(false,nil)`，关掉跨进程 W1 残留；D3（`release` 改为先 `os.Remove` 再 unlock）触及成功 / 失败 / 退休热路径且依赖 Windows remove-while-open 语义，本片不做，登记后续片。

## 审查链结论

- **① 定性**：`deep-reasoner` + DeepSeek V4 Pro；闭集判定 `DIAGNOSIS: B`，附 15 条 file:line 证据链与对 P45 的 5 条输入。
- **②③ 修复与反例**：D1 / D2 与反例 `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` 落地；反例在**未修复代码**上确定性红——Windows 失败点 `spawn_recovery_test.go:281`（marker 断言），unix 预期在「子进程存活」断言红（本机未执行）。
- **④ 门禁**：`TOTAL_FAIL=0`；`go build` / `go vet` / `go test ./...` 全绿，读数见 §门禁读数。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（0 Critical / High；2 Medium + 1 Low + 2 Info）。
- **⑤ r2**：`PRE-REVIEW: PASS`（B3 已闭合；B1 / B2 / B4 / B5 处置批准；附 N1–N4 报告级义务）。
- **⑥**：`INDEPENDENT SCAN: PASS`，Section B 判 `0 · 无`。
- **⑦**：`FINAL REVIEW: PASS`，0 新增发现，附 6 条放行条件（逐条落地见 §偏差与登记）。
- **R / D**：R = 0、D = 0（不继承上片）。

## 六段扫描转录

来源 `docs/validation/evidence/DR13-FIX/06-scan/review.md`；围栏内按 A→B→C→D→E 顺序转录 ⑥ 正文，围栏外披露两处归一。

```sixstage
## Section A 扫描范围与输入
INDEPENDENT SCAN: PASS
范围：只读核查 internal/adapters/agent_runner_pinned_cli.go、agent_runner_process_registry.go、spawn_recovery_test.go 与 spawn_recovery_live_host_unix_test.go，聚焦 D1 / D2 落点、反例构造与并发回归。

## Section B 发现清单
| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
|---|---|---|---|---|
|  |  |  |  |  |

## Section C 与 ①/⑤ 的一致性核对
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

## Section D 未覆盖面与残余风险
1. `StopAgentRunnerProcess` 的失败分支仍能在“输入解析/停止失败”时无条件 `releaseLaunch`; 这不等于 D1 失效，但它是下一个最值得加固的边界，尤其在停止中途状态不一致时可能清掉 `registry` 句柄 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L271-L301)。  
   重要性：它是 stop/start 竞态的镜像窗口，若下一轮再复现 stop-while-starting，可能出现“已杀/未杀状态模糊”。
2. D2/TOCTOU 仍依赖平台层对 `flock`/inode 语义的一致性；本地仅能读代码而不能在 unix CI 上直接探针 `flock` 的行为差异，因此 W2 的内核级细节未被本机独立实证 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L343-L377)。  
   重要性：它关乎“跨进程锁释放是否会让另一个扫描持到 orphan inode”的真实发生率。
3. 证据包的 `01-diagnosis` 是已归档的定性结论，但并未在本机 Unx 宿主上复现 `StopProcessIdentity` 对活进程的具体行为；当前的证据仍依赖 CI/设计推演 [docs/validation/evidence/DR13-FIX/01-diagnosis/review.md](docs/validation/evidence/DR13-FIX/01-diagnosis/review.md)。  
   重要性：这里是最终发布前最需要的真实 runtime 证据，尤其是 ubuntu 腿。
4. D1 仅覆盖“同一 launcher 仍持有同一 launchID”的漏洞；它不覆盖“跨 launcher、跨进程、同一 runRoot 但镜像已失效”的不同窗口，后者仍需要 slot/镜像一致性检查和更严格的 stop 语义 [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L620-L700)。  
   重要性：这类窗口不会触发 D1，但仍可能导致不完全的“owner 检查”。

## Section E 结论要旨
可证伪性成立：当前修复的核心是两层“先证实 marker 仍指向真实 published path，再确认 launcher.registry 没有同 launchID 活句柄”。这两个门都可以被单步反证：

- 若删掉 D1 门，`TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns` 会直接失败，且 helper 进程会被误杀或 slot 被删除 [internal/adapters/spawn_recovery_test.go](internal/adapters/spawn_recovery_test.go#L209-L284)。  
- 若删掉 D2 门，`tryLock` 成功后到 `os.Stat` 之间的 orphan inode window 就能导致误回收；这在旧实现下会在“marker path 已 unlink”时制造 false positive [internal/adapters/agent_runner_pinned_cli.go](internal/adapters/agent_runner_pinned_cli.go#L343-L377)。  
- 若没有真实 helper 进程、真实 slot claim、真实 marker 写入，测试就会失去证伪性；当前测试恰好使用了实际的 live process 和 slot 创建，不存在假绿通道。  
- 这一点也与 `go test` 结果一致：当前实现通过 package 级回归测试，说明反例是可复现且被封住的。

INDEPENDENT SCAN: PASS
```

- 归一披露一：⑥ 的段标题 `## A 结论` / `## B 缺陷` / `## C 核验结论` / `## D 未覆盖面` / `## E 可证伪性` 在围栏内改写为本模板的段名，五段正文逐字未改。
- 归一披露二：Section A 由两行组成——⑥ 的总体结论行（逐字）与范围行（取自 ⑥ Section C 所核查的文件集合）；⑥ Section E 的末行 `§A INDEPENDENT SCAN: PASS` 归一为围栏末行 `INDEPENDENT SCAN: PASS`。

## 项闭合

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 B1 | 本轮读数全部来自 Windows 本机，`//go:build unix` 用例本机不编译，而 DR-13 的原始复现腿是 ubuntu | Medium | `spawn_recovery_live_host_unix_test.go` | 已处置：⑨ 以 Ubuntu 腿为最终判据，登记见 §偏差与登记 R2 |
| ⑤ r1 B2 | D2 无专属确定性测试，唯一触发窗口依赖调度时序 | Medium | `agent_runner_pinned_cli.go:367-376` | 已处置：登记为「仅概率覆盖」，见 §偏差与登记 R3 |
| ⑤ r1 B3 | `io/fs` 仅用于 `fs.ErrNotExist`，与文件其余处的 `os.ErrNotExist` 风格不一致 | Low | `agent_runner_pinned_cli.go:370` | 已闭合：改用 `os.ErrNotExist` 并删除 `io/fs` import，全文件共 6 处 |
| ⑤ r1 B4 | 切片登记行超出声明边界 | Info | `docs/t027/REMAINING_SLICES.md` 与 `_EN`（各 `+23`） | 已处置：属开工簿记，按边界豁免登记 |
| ⑤ r1 B5 | 反例的「锁空闲 + 活进程」为模拟态，防御价值成立但非可达时序复现 | Info | `spawn_recovery_test.go:221-290` | 已处置：测试注释已如实声明，不改代码 |
| ⑤ r2 N1–N4 | 报告级义务：ubuntu 包级 `ok` 行、B2 / D3 登记、措辞修正并入边界簿记 | — | 本报告 | 已闭合：见 §偏差与登记 R1–R6 与 §已知限制 |
| ⑥ | Section B 判 `0 · 无`，Section C 六条核验全部成立 | — | `06-scan/review.md` | 已闭合：`INDEPENDENT SCAN: PASS` |
| ⑦ | `FINAL REVIEW: PASS`，0 新增发现，附 6 条放行条件 | — | `07-final/review.md` | 已处置：逐条落地见 §偏差与登记 R1–R6 |

## 交付

- `internal/adapters/agent_runner_pinned_cli.go`（`+37/−0`）：D2 路径复核——`tryLock` 成功后 `os.Stat(path)`，`os.ErrNotExist` 即锁住孤儿 inode 而返回 `(false,nil)`，其余 I/O 错误包 `ErrAgentRunnerLaunchConflict`（`:367-376`）；D1 同进程所有权门——decode 成功后 `registry.lookup(marker.LaunchID)` 命中即 unlock、close、`return false, nil`（`:390-397`）。两段均带设计注释，行为零协议变更。
- `internal/adapters/spawn_recovery_test.go`（`+84/−0`）：反例 `TestRecoverStartingNeverReclaimsALaunchThisProcessStillOwns`（`:221-290`），以真实 helper 进程、真实 `claimLaunchSlot`、真实 `registry.register` 与无锁 marker 构造「锁空闲 + 进程存活 + 本进程在册」三重态；四条断言为 `reclaimed==false`、helper 仍活、marker 仍在、slot owner 仍为该 launchID；`t.Cleanup` 收掉子进程与持久记录。
- 证据包 `docs/validation/evidence/DR13-FIX/`：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余 7 件为载荷（`00-briefing.md`、`freedom-list.md`、`01-diagnosis/review.md`、`03-tests/red-baseline.md`、`05-prereview/review.md`、`06-scan/review.md`、`07-final/review.md`）。
- 台账：`docs/t027/REMAINING_SLICES{,_EN}.md` 的 DR13-FIX 汇总行与详情块翻为完成态，DR-13 条目的处置计划改为已修复并保留重开口径。

## 偏差与登记

| 编号 | 内容 | 来源 | 判定 |
|---|---|---|---|
| R1 | Ubuntu 腿证据：并发单赢家与 unix recovery 两个用例的包级 `ok` 行须由 ⑨ 提取 | ⑦ 放行条件 1 | 已处置：移交 ⑨，读数在 §提交与 CI 追加 |
| R2 | 本机无法执行 `//go:build unix` 用例，ubuntu 腿为最终判据 | ⑦ 放行条件 2 与 ⑤ r1 B1 | 已处置：登记为 B1；⑨ 任一腿红即回 ⑤ 重开 |
| R3 | D2 无专属「孤儿 inode」确定性测试，仅由 D1 的确定性反例与并发回归概率覆盖 | ⑦ 放行条件 3 与 ⑤ r1 B2 | 已处置：登记接受；`tryLock` 双平台均非阻塞，零成本的可注入边界不存在 |
| R4 | D3（`release` 改为先 `os.Remove` 再 unlock）本片不做 | ⑦ 放行条件 4 | 已处置：触及成功 / 失败 / 退休热路径且依赖 Windows remove-while-open，登记为后续片 |
| R5 | 跨进程 W1 残留：`os.Remove` 失败叠加跨进程扫描时 D1 不覆盖 | ⑦ 放行条件 5 | 已处置：登记适用边界为同 runRoot、跨 launcher |
| R6 | ⑥ Section D 四项未覆盖面 | ⑦ 放行条件 6 | 已处置：四项全部登记于 §后续片候选登记 与 §已知限制，本片不按缺陷计 |
| D-a | ⑤ r1 B4：台账登记行在 ②③ 之前落盘，超出代码边界声明 | ⑤ r1 B4 | 接受：属开工簿记（切片条目与证据包骨架），无代码影响 |
| D-b | ⑤ r2 措辞修正：`os.ErrNotExist` 在文件中合共 6 处（原 5 处未动 + D2 分支 1 处由 `fs.ErrNotExist` 统一而来），「无中文注释」限新增行口径 | ⑤ r2 B-r2-1 与 B-r2-2 | 接受并整改：本报告两处均按修正后口径表述 |

## 模型与计量

| 环节 | 模型 |
|---|---|
| ① 定性 | DeepSeek V4 Pro (deepseek) |
| ② 实现 | DeepSeek V4.1 Flash (deepseek) |
| ③ 测试 | DeepSeek V4.1 Flash (deepseek) |
| ⑤ 预审（r1 / r2） | DeepSeek V4 Pro (deepseek) |
| ⑥ 独立扫描 | MAI-Code-1.1-Flash (copilot) |
| ⑦ 终审 | GPT-5.3-Codex (copilot) |
| ⑧ 文档 | DeepSeek V4.1 Flash (deepseek) |

- **批量辅助**：`GLM-5.3-Flash (glm)` 实派 4 次——批次 1 排他路径普查、批次 2 测试装置、批次 3 日志比对、批次 4 文件操作普查；另有批次 5 的 diff 核对，作为 ⑤ r2 逐行核对的输入。
- **主控模型**：`DeepSeek V4.1 Flash`（High 档）；①⑤ 按白名单仍由 V4 Pro 承担，与主控选型是两件事。

## 门禁读数

- 命令：`node tools/gates/gate.js --all --skip-base`；tree 作用域 `TOTAL_FAIL=0`（归档工件按 G7-b 要求精确 `git add` 后实跑）。
- `G3(a)`：台账镜像两侧的新增 / 删除行数相等，报告镜像两侧亦相等；报告对的六形状（行数 / 标题行号序列 / `**` 总数 / `|` 总数 / 空行数 / 围栏与反引号数）全等。
- `G4b`：`unknown=[]`；`G6-4`：`no expiring literal`；`G6-8`：DR 集合含 13，无撞号与悬挂引用。
- `go build ./...` 与 `go vet ./...` 退出码 0；`gofmt -l` 无输出。
- `go test -count=1 ./...`：`ok` 行 15 条、`FAIL` 0 条（19 个包条目，其中 4 个 `[no test files]`）；起跑消息记载的「16 包」与本次实测差 1，按实测登记。

## 提交与 CI

- 提交与 CI 读数在 ⑨ 之后由本节追加一次，追加时同步更新 CN/EN。

## 后续片候选登记

- D3 先删后解锁：`release` 改为 `os.Remove` 先于 unlock，配合 D2 全闭 W1，需 Windows remove-while-open 全量回归。
- 跨进程 W1 残留：`os.Remove` 失败叠加跨进程扫描时，建 marker 残留的可观测告警或补偿路径。
- `StopAgentRunnerProcess` 失败分支：解析 / 停止失败时不再无条件 `releaseLaunch`，改为按已证状态分流。
- W2 flock 内核语义：在 unix 侧实证「关闭任一 fd 即释放该进程全部 flock」的现行行为。
- unix `StopProcessIdentity` 对活进程的行为：本机（Windows）无法复现，移交 CI / unix 宿主取证。
- 跨 launcher / 跨进程镜像失效窗口：补 slot 与镜像一致性检查，覆盖 D1 触达不到的窗口。

## 已知限制

- **B1**：本机无法执行 `//go:build unix` 用例，unix 分支（含杀活路径与孤儿收养）的修复后判定以 Ubuntu 腿为最终判据。
- **B2**：D2 无专属确定性测试，仅由 D1 反例与并发回归概率覆盖；`tryLock` 双平台均非阻塞，外部无法确定性占据 stat 前窗口。
- **D3**：`release` 仍是 unlock 后 remove，路径消失晚于锁释放的毫秒窗口在跨进程侧残留，本片以 D1（同进程）确定性拦截。
- **W1**：跨进程残留需 `os.Remove` 失败为前提，概率极低且本片未闭合，适用边界为同 runRoot、跨 launcher。

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/dr13-fix.md	repo
docs/validation/dr13-fix_EN.md	repo
docs/validation/evidence/DR13-FIX/00-briefing.md	repo
docs/validation/evidence/DR13-FIX/01-diagnosis/review.md	repo
docs/validation/evidence/DR13-FIX/03-tests/red-baseline.md	repo
docs/validation/evidence/DR13-FIX/05-prereview/review.md	repo
docs/validation/evidence/DR13-FIX/06-scan/review.md	repo
docs/validation/evidence/DR13-FIX/07-final/review.md	repo
docs/validation/evidence/DR13-FIX/MANIFEST.md	repo
docs/validation/evidence/DR13-FIX/SHA256SUMS.txt	repo
docs/validation/evidence/DR13-FIX/freedom-list.md	repo
internal/adapters/agent_runner_pinned_cli.go	repo
internal/adapters/spawn_recovery_test.go	repo
```

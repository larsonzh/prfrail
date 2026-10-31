# P40-T027B-LOCK-REAPER ⑤ 预审归档（两轮原样）

切片：`P40-T027B-LOCK-REAPER`。环节：⑤ 预审。载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`。第 1 轮与第 2 轮输出原样归档，内容未作改写。

---

## 第 1 轮输出（原样）

预审核验完成。以下为 ⑤ 预审报告。

### 一、独立核验摘要（本人实测，非转述）

| 项 | 命令 | 结果 |
|---|---|---|
| 构建 | `go build ./...` | exit=0 ✓ |
| 静态检查 | `go vet ./...` | exit=0 ✓ |
| 全量测试 | `go test ./... -count=1` | exit=0，16 包全 `ok`、0 FAIL ✓ |
| `-race` 本机 | `go test -race ./internal/guard/` | exit=2（`-race requires cgo`）→ 偏差 G 前提成立 |

变更范围与切片声明一致（`git diff --stat`：`internal/guard/` 4 生产文件 + 1 新测试、`internal/tickets/` 4 生产文件 + 1 新测试、docs 2 文件）。

### 二、审阅重点逐条结论

1. 状态机一致性——基本对齐，一处承接缺口（StartManaged 失败清理到界未关时句柄无回收路径 → F-01）。
2. 锁序守卫正确性——主干对齐，两个子路径未钉住（相等路径无直测、panic-pop 无直测）。
3. 偏差 A–G 安全评审——A 平台指针安全；B 10s 宽窗+自证安全、需文档回写；C 脱离式 taskkill 安全、零测试；D 自然退出 close 条件方向正确、零测试；E stop 不加 settle 安全；F goroutine id 退化安全；G -race 靠 CI 但 Race leg 不含 guard/tickets。
4. 覆盖完整性映射——R1/R2/R3、L1-L4 与既有测试钉住清单如报告正文；未覆盖窗口：Reclaim 零直测、taskkill 脱离/超时零测试、自然退出树残留零测试、锁序相等路径与 panic-pop 零测试、同进程同 scope 非嵌套并发零测试、-race 不含 guard/tickets。
5. 可证伪判据（§4.2.1 第 4 条）：删除 `waitOnce.Do` → R3 红；锁序 check 移到 mu 后 → F-06 死锁红；删除 backoff ctx 检查 → L3 红。

### 三、发现清单

| # | 严重度 | 位置 | 问题 | 建议修法 |
|---|---|---|---|---|
| F-01 | Medium | `internal/guard/process.go:139-166` | `abortStartup` 到界返回 uncertain 且不 close，但 `process` 不逃逸出 `StartManaged` → Windows Job 句柄永久泄漏 | 到界支路追加 `go func(){ <-process.waitDone; process.platform.close() }()`；回写设计行 4 注记 |
| F-02 | Medium | `internal/guard/process.go:366-385`；`process_reaper_test.go:96` | `Reclaim` 无任何直接测试、无生产调用方 | 新增 `TestReclaimReleasesHandleOnce`：幂等/拒绝/句柄恰释放一次 |
| F-03 | Medium | `.github/workflows/ci.yml:84` | Race leg 仅覆盖 adapters/chain，guard/tickets 不跑 `-race` | 扩展 Race leg 或如实登记原覆盖范围 |
| F-04 | Medium | `internal/guard/process_windows.go:219-253`、`process.go:288-309` | 10s 宽窗、脱离式回收、到界自证、Actions 追加零测试 | 极小 bound 断言 `errStopCommandExpired`+`taskkill-detached`；settled 动作断言 |
| F-05 | Medium | `internal/guard/process.go:191-216` | `ErrProcessTreeRemained` 支路无 guard 级回归（仅 Linux 可达） | `process_linux_test.go` 增父退子残留用例 |
| F-06 | Medium | `internal/tickets/durable_cost_lock_shared.go:57-91` | 同栈相等路径拒绝与 panic-pop 无测试 | `nested(ledgerA, ledgerA)` 断言 Reentrant + panic 后可重新取锁 |
| F-07a | Low | 设计 §2 断言 1 | 静态 CI 嵌套断言未实现 | 登记为后续片候选，本片不强制 |
| F-07b | Low | 设计 §1 行 5/行 8 vs `process.go:288-309` | 实现语义比设计字面更严，文档未回写 | 设计文档行 5 加注记 |
| F-08 | Low | `process_windows.go:241` | `ctx.Done` 支路与 timer 等价冗余 | 可保留（无害） |
| F-09 | Low | `process.go:272,327` | nil-ctx 防御风格不一致 | 可延后 |
| F-10 | Low | `durable_cost_ledger.go:321-350` | 同 scope 非嵌套并发无测试 | 可并入 F-06 |
| F-11 | Low | `durable_cost_lock_shared.go:67` | lockPath 字符串比较非 canonical | 登记即可 |
| F-12 | Low | `console/stopper.go:15-16` | 10s+grace 与 stopper 外层 12s 关系 | 注释登记 |

FINAL REVIEW: PASS WITH FIXES

（0 Critical / 0 High；6 Medium 全部为覆盖缺口/潜在缺陷类且修法均已给出、可在 ⑥ 前落地；核心正确性——唯一 reaper、锁序单调、有界等待——经本人独立复验全绿。）

---

## 第 2 轮输出（原样）

# P40-T027B-LOCK-REAPER · ⑤ 预审第 2 轮（F 项清零确认，只读）

审查依据：`tmp/P40-T027B-LOCK-REAPER/fix.diff`（45,720 字节全量）、`fix.numstat`/`status.txt`、工作区直读、独立只读验证（`git status`、`gofmt -l`、BOM/CRLF、SHA256、`GOOS=linux go vet`、焦点测试三组重跑）。本环未改任何文件、未做任何 git 写操作。

### 逐条核对

- F-01 — CLOSED（附 2 条 NOTE）：后台 reaper `process.go:155-176`；快路径同步 close（:160-161）、到界支路先起 goroutine 再返回 uncertain（:164-175）；无句柄无人回收窗口；platform.close 幂等（Windows `job != 0` 判 + mutex、Linux `released` 标志、unsupported 空实现），abortStartup 的 process 永不出 StartManaged，Reclaim 不可能触达。NOTE-1：后台 close 不做 verifyGone（Windows close 即 kill-on-job-close，正确终局化）。NOTE-2：生产侧 Terminate 调用方（`agent_runner_pinned_cli.go:577`）无人接 Reclaim，uncertain 停树句柄保留至进程退出（有界），接线属后续片。
- F-02 — CLOSED：`process_reclaim_test.go` `TestReclaimReleasesHandleOnce`（:39-108）钉住四断言（uncertain 后 closeCount=0 → Kill 后 Reclaim nil + closeCount=1 → 二次 Reclaim 幂等 → 未停进程返 ErrTerminationUncertain）；helper `TestReclaimHelperStayAlive`（:17-23）用 `time.Sleep` 循环不触死锁检测。本机重跑 PASS（0.07s）。
- F-03 — CLOSED：`.github/**`、`internal/release/**` 零改动；「登记而非改代码」可接受；登记物须落在 ⑨ 报告/`REMAINING_SLICES.md` 后续批次。
- F-04 — CLOSED WITH NOTE：纳秒 bound → `errStopCommandExpired` + `taskkill-detached` + `taskkill-process-tree`；live-tree → stopped、already-exited → already-stopped。未覆盖的仅 settled 标记追加与 detached 错误 join 两条平凡组合行；无遗留可疑分支。NOTE：`review.md` 行 13 的 `exec.CommandContext` 措辞与实现漂移，建议 ⑨ 同步措辞（已由主控在 ⑥ 前收口）。
- F-05 — CLOSED WITH NOTE：`process_linux_test.go:110-172` 父 helper 裸 spawn 子进程构成「父退出后进程组仍活」；断言 ErrProcessTreeRemained + 证据 + `signal-term-process-group` + closeCount==1；`GOOS=linux go vet` 零输出；CI ubuntu 腿风险极低（无 hidepid、SIGKILL 升级、30s 裕量）。
- F-06 / F-10 — CLOSED：同栈重入拒绝、panic 后 pop、非嵌套并发串行化三用例钉住；check 在 `mu.Lock()` 与文件锁获取之前执行，拒绝路径不会先卡在 flock。
- F-07b — CLOSED：行 12/行 13 补注与原行语义无冲突；BOM+LF；SHA256SUMS 3 件载荷与磁盘实算逐一相等。
- 第 1 轮 6 项 Low 处置口径 — 可接受（条件：F-07a 与 F-11 的登记必须出现在 ⑨ 报告或台账后续批次）。
- ② 观测性改动（`releases`/`closeCount`）— 无语义风险（三形态齐全、mutex 保护、per-instance 计数）。

### 新发现

| 级别 | 位置 | 内容与修法 |
|---|---|---|
| Low/Note | `01-arch-design/review.md:13` | 行 5 补注与实现的 `exec.CommandContext` 措辞漂移；修法：⑨ 报告装配时同步措辞（已由主控收口） |
| Note | `internal/adapters/agent_runner_pinned_cli.go:577` | 生产 Terminate 无 Reclaim 接线；与 F-07a 同批登记为后续片候选 |
| Note | `process_reaper_test.go:48,115`、`process_test.go:82` | cleanup 中重复 Process.Wait 依赖 os 层互斥/缓存语义；非缺陷，登记脆弱性 |

第 1 轮结论（0 Critical / 0 High / 6 Medium + 6 Low）在本轮全部核实闭合；无新增 Critical/High/Medium。

PRE-REVIEW: PASS

# P41-T027B-SPAWN-OWNERSHIP ⑤ 预审归档（两轮原样）

切片：`P41-T027B-SPAWN-OWNERSHIP`。环节：⑤ 预审。载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`。两轮输出原样归档。

---

## 第 1 轮输出（原样）

抽查验证完成（`go test -count=1 ./internal/guard/... ./internal/adapters/...` 复跑绿；`-v` 计数 348 PASS / 0 FAIL / 4 SKIP——4 个 SKIP 均为 helper 子进程；`go vet` 干净；`full1.txt` 15 ok；`gate-tree4.txt` TOTAL_FAIL=0）。另注：③ 偏差清单原文未在证据包内（仅 01-arch-design 已落盘），6 项偏差以派单文本描述为准审校，建议 ⑧c 前补齐 ③ 归档。

# P41-T027B-SPAWN-OWNERSHIP ⑤ 预审结论

## 逐条结论

**1. 归属先于 resume — CLOSED**
- 顺序硬编码正确且不可绕过：`preparePlatformProcess`（建 Job，`KillOnJobClose`，process_windows.go:137-150）→ `start`（CreateProcessW 一次调用即含归属）→ `attach`（仅验证）→ `platformIdentity` → `OnSpawn` → `resume`（process.go:200-237）。
- 创建期归属成立：`CreateProcessW` 携带 `JOB_LIST` 属性 + `EXTENDED_STARTUPINFO_PRESENT` + `CREATE_SUSPENDED`，`bInheritHandles=1`（process_windows.go:163-210、259-267）；`attach` 退化为 `IsProcessInJob` 验证（process_windows.go:485-505）。kill-on-close 语义随创建期归属成立。
- 句柄继承与 os/exec 同策：`HANDLE_LIST` 仅含 3 个标准句柄（process_windows.go:237），父端写端在 CreateProcessW 返回后立即关闭（defer `closeParent`，process_windows.go:237/401-418），输出在 `waitFn` 内 `drain`——两个管道缺陷均闭合。

**2. OnSpawn 契约 — CLOSED**
- guard 失败处理与设计一致：hook 失败 → `platform.terminate()` → `abortStartup` → `ErrSpawnOwnership` 包装（process.go:224-230），attach/identity/resume 失败同出口（214-235）。
- hook 自清理边界正确：register 失败时 `claim.release()` 自删 marker（pinned_cli.go:206-209）；resume 失败由 `discardStartingLaunch` 回收。
- nil hook 零迁移成立：`OnSpawn == nil` 直接跳过，旧调用点不动。
- 小注：并发 stop 先赢时，输侧返回的是 `ErrSpawnOwnership` 包 resume NTSTATUS 错误，而非设计 Q1 措辞的「already-stopped」——语义安全，措辞漂移。

**3. starting durable marker — CLOSED WITH NOTE**
- 原子发布 temp+锁+rename（marker.go:52-87）；LockFileEx 字节 0 起 1 字节、EX|FAIL_IMMEDIATELY（marker_windows.go:47-67）；创建即带 share-delete；跨 rename 锁存续成立。
- 严格解码：`DecodeStrictJSON`（Canonicalize + DisallowUnknownFields + 尾部值拒绝 + 64KB 上限）。
- Linux flock 等价（偏差 3）：LOCK_EX|LOCK_NB 于 rename 前加锁，内核死亡即释放，子进程不继承——机制合理；新竞态见 F-1。
- `advance(starting→running)` 时点在 mirror 写后、marker 释放前；`stopIdentity` 经 `hasSpawnIdentity` 支持 starting 态。
- NOTE：`createLockedStartingMarkerFile` 忽略 `tryLock` 的 `locked` 布尔，见 F-3。

**4. RecoverStarting — CLOSED WITH NOTE**
- 调用点正确：`RecoverStarting` 在 `claimLaunchSlot` 之前。
- 锁持有 ⇒ 拒绝；锁释放 ⇒ 解码 → `ProcessAlive` → PID 复用复核 → 平台回收。
- 偏差 4（Windows 锁释但进程存活 ⇒ 不杀保留）安全：fail-closed 保留 mirror 后 `refuseLiveMirrorProcess` 阻断下一次启动。
- NOTE：F-1、F-4、F-5。

**5. 偏差 1/2（managedCommand / prepare 签名）— CLOSED WITH NOTE**
- 兼容面充分：348 PASS；`exitCode()` 取值正确（单一写入者、经 `waitDone` 建立 happens-before）。
- NOTE：Windows `managedProcessHandle.Wait()` 返回 nil `ProcessState`；argv[0] 用原始 `spec.Command` 而非解析后路径，与 os/exec 有外观差异。

**6. 红测不变式与不可测项登记 — CLOSED WITH NOTE**
- I1 顺序已钉；I4、I3/I5 均有红测且转绿。
- I1 kill-on-close 半证：登记为后续片候选（仅测试注释，须随报告落盘）。
- I2 四缝注入未建：按本片范围可接受，报告须写清理由。
- Linux 运行时未实测：flock 拒绝路径与领养击杀路径均无平台运行时覆盖（F-2），报告须如实补登记。

**7. 全绿真实性 — 抽查通过**
- 复跑双包 ok；348 PASS / 0 FAIL / 4 SKIP；7 条新红测逐一 PASS；vet 干净；全量 15 ok；tree 门禁 TOTAL_FAIL=0。

## 新发现（按严重度）

- **Medium F-1** — RecoverStarting 按路径清除未复核归属：扫描期间另一宿主 rename 发布新 marker，旧扫描会误删活 marker ⇒ Linux 孤儿失去 durable 证据、槽位新鲜度过后可被再启动（重开 C2 窗口）。修复：清除前重开路径 + 重试加锁 + 重解码比对身份，锁被持有或内容不符则跳过清除。
- **Medium F-2** — Linux 回收路径无运行时覆盖：flock 活宿主拒绝测试仅 Windows 有；领养击杀在任何平台都无用例。建议补两条 Unix 用例或登记。
- **Low F-3** — 锁布尔结果被忽略：`(false, nil)` 会签发未持锁的 claim。
- **Low F-4** — 损坏 marker 只清文件不清槽位：槽位新鲜度阻塞 10s。
- **Low F-5** — 拒绝态与「无 marker」不可区分：病理场景可绕过拒绝直达槽位抢夺，由 hook 速度契约缓解。

## 总体评价

设计→实现落地一致性高；6 项偏差中 1/2/3/4 均安全，5/6 待 ③ 归档核对。7 个审阅点全部可闭合，但 F-1 与 F-2 属本片协议语义缺口，且修复成本低，应在 ⑥ 前闭合。

PRE-REVIEW: PASS WITH FIXES

---

## 第 2 轮输出（原样）

# P41-T027B-SPAWN-OWNERSHIP · ⑤ 预审第 2 轮（F 项清零确认，只读）

输入：`tmp/P41-T027B-SPAWN-OWNERSHIP/fix.diff`（60,687 B）、`status2.txt` + 工作区直读。本环未改动任何文件、无 git 写操作。

## 回归复核（先记事实）

- `go test ./...`：**15 ok / 0 FAIL**（adapters 58.8s、guard 12.3s 本机实跑，非缓存）✓
- `go test ./internal/guard/ ./internal/adapters/ -count=1 -v`：**348 PASS / 4 SKIP / 0 FAIL** ✓；4 SKIP 均为 helper 子进程自检跳过，无意外跳过 ✓
- ④ 门禁 `node tools/gates/gate.js --all --skip-base`：**TOTAL_FAIL=0**（G9-e 两条 INFO 为台账 tmp 引用未 staged 提示，非 FAIL）✓
- 交叉验证：`GOOS=linux go vet ./internal/adapters/ ./internal/guard/` exit 0，随后复核 `go env GOOS` = windows（P40 教训）✓

## F 项逐条

**F-1 — CLOSED WITH NOTE**（`agent_runner_pinned_cli.go:493` `clearAgentRunnerStartingMarker`）
- 重证流程完整：重开（`openStartingMarkerForClear`）→ 重试加锁（3 次 × 5ms，`!locked` 仅 continue 不删）→ 严格重解码 → 全结构比对 → 锁持有时 `os.Remove(path)`；open 失败 / 不匹配 / 超次一律不删。
- D-1 平台助手实证：Go 1.27 与 CI 1.22 的 `syscall.Open` 均不含 FILE_SHARE_DELETE ⇒ 助手必要；Windows 本机恢复用例 marker 删除断言全绿 = 实证。
- NOTE：① 理论残余窗口——验证 open 时 inode、删除 remove 时 path，另一宿主恰在 lock→remove 间 rename 覆盖会删到新记录；微秒级且 slot 协议排除同 run root 并发启动，实际不可达，留待后续片评估改 rename-away 清除；② D-2（owner release 不重证）边界可接受。

**F-2 — CLOSED**（`spawn_recovery_live_host_unix_test.go:113/:173`）
- `TestRecoverStartingRefusesLiveHostUnix`：独立 fd flock（同进程双 open 冲突）模拟活宿主 ⇒ 扫描拒绝、marker 与活进程不受扰；释放锁+进程死亡后回收。构造正确。
- `TestRecoverStartingAdoptsOrphanUnix`：`Setsid` 独立 session/进程组；`InspectProcess` 取真实身份；marker 无锁+live 孤儿 ⇒ 领养击杀断言（orphan 死、marker/slot 清理、二次扫描 no-op）。构造正确。
- Linux CI 可运行性：`unix` 内建 tag；ubuntu-latest 腿 Test 步实跑 adapters 包；Go 1.22 支持 `Setsid`/`Flock`；helper 有界 sleep 且 cleanup 必 reap。风险低。
- 本机未实跑已如实登记：build tag 平台分裂 + 两文件注释声明；正式「未执行测试」登记待报告落盘时补全。

**F-3 — CLOSED**（`agent_runner_starting_marker_unix.go:118` / `_windows.go` 同构）
- 三处 `tryLockStartingMarker` 调用无 `(false,nil)` 漏出；错误文案与红测无冲突。

**F-4 — CLOSED**（`agent_runner_pinned_cli.go:430` `reclaimLaunchSlot`）
- 损坏 marker 走 corrupt 分支：rename-away → 读 owner → 空字段放回 → `refuseLiveMirrorProcess` 证明无活 mirror 之后才删 staging；拒绝失败放回 slot。删的是自己 rename 走的 staging 旧文件；新 launch 槽位 O_EXCL 新建独立文件，不会误删新槽。

**F-5 — CLOSED**（`agent_runner_pinned_cli.go:357`）
- 登记如实：注释明写 `(false,nil)` 歧义与「deliberately not changed」。
- 附加观察（Info）：红测写法为 `reclaimed, _ :=` 且注释自认「error 也是可接受拒绝」，故拒绝态改返 error 并不破坏红测；上轮已裁「不修」，本轮不推翻。

## 新发现（非阻断）

- **Low-1**：starting 相位防双进程完全依赖 slot 10s 新鲜度兜底——设计书风险 3 已预判（hook 只做本地写+注册，毫秒级），实际风险极低；与 F-5 联动建议后续片评估拒绝态改返 error。
- **Low-2**：F-1 残余的 remove-by-path vs rename 替换 inode 窗口（微秒级、实际不可达）。
- **Info-1**：恢复成功路径的「mirror 存在且匹配 ⇒ 删除」分支无测试断言。
- **Info-2**：Linux `Pdeathsig` 未引入——设计书列为「可选加固」，实现选择领养-击杀为主路径，属设计书允许路径；Linux 最坏回收延迟=重启间隔。
- **Info-3**：D-1..D-6 正式偏差登记未见落盘文件——代码注释为当前唯一就地登记，⑧ 回写时须逐项补全。

PRE-REVIEW: PASS

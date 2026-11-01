# ⑤ 预审记录（两轮原样归档）

> 载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。r1 偏移 120049，r2 偏移 120219。

## r1

## A 结论

`PRE-REVIEW: PASS WITH FIXES`

## B 发现

| 编号 | 严重度 | file:line | 问题 | 修法 |
|---|---|---|---|---|
| 1 | Medium | `spawn_recovery_live_host_unix_test.go` L118 / L155（`//go:build !windows` 文件） | 本轮验证读数（`^TestRecoverStarting` 全绿、`go test ./...` 16 包、`-count=20`）均系 Windows 本机执行；`TestRecoverStartingRefusesLiveHostUnix` 与 `TestRecoverStartingAdoptsOrphanUnix` 因 build tag 在本机不编译。而 DR-13 原始复现腿正是 ubuntu CI（run `41253798647` / `41257990185`），D1 的「不误伤」与孤儿正路径恰由这两个测试证明 | 门禁前补一条 ubuntu CI 腿重跑证据（含 `ok github.com/larsonzh/prfrail/internal/adapters` 包级行），否则 unix 分支本轮无执行证据 |
| 2 | Medium | `agent_runner_pinned_cli.go` L367–377 | D2 无专属确定性测试：唯一能触发「open 后、tryLock 前被 unlink」窗口的是并发胜者测试，其命中依赖调度时序（loser 的 open 须落在 winner `closeStartingMarker` 的 Remove 之前）。`-count=20` 绿不能证明 D2 被覆盖；若 D2 被改坏，并发测试可能仍绿 | 补 D2 专属回归（持锁删名后断言 `(false,nil)` 且不 reclaim），或登记「D2 仅由并发回归概率性覆盖」并接受 |
| 3 | Low | `agent_runner_pinned_cli.go` L7、L371 | `io/fs` 仅用于 `fs.ErrNotExist`；Go 1.16+ 中 `os.ErrNotExist` 即 `fs.ErrNotExist` 的别名（os/error.go），二者 `errors.Is` 完全等价。本文件其余 4 处（L339、L599、claimLaunchSlot、refuseLiveMirrorProcess）一律用 `os.ErrNotExist`，风格不一致 | 改用 `os.ErrNotExist` 并删除 `io/fs` import（非必需改动） |
| 4 | Info | `docs/t027/REMAINING_SLICES.md` / `_EN`（+23×2） | 工作树存在切片表 DR13-FIX 登记，超出声明边界「仅 `internal/adapters/agent_runner_pinned_cli*.go` 与其测试文件」；属切片开工簿记，无代码影响 | 交付时在切片 closeout 中注明该簿记项，或视为边界豁免 |
| 5 | Info | `spawn_recovery_test.go` L221–284 | E1 的三重态（marker 在路径 + 锁空闲）是模拟态：真实 `closeStartingMarker`（L273→`release()`）会删除 marker 文件，真实中该窗口极窄（仅 `os.Remove` 失败或 release 与 scan 交错）。防御价值成立，但属「防御性场景」而非可达时序复现 | 无需改动；在测试注释已如实声明（L214–219），保持 |

未发现 Critical/High 缺陷。

## C 交叉核对

1. **独立实测（Windows 本机）**：`go test ./internal/adapters/ -run "^TestRecoverStarting|…ConcurrentStartsHaveOneWinner" -count=1 -v` —— 8 个测试全 PASS（含 E1、Idempotent、TokenMismatchFailsClosed、RefusesLiveHost、CleansCrashedHostWindows、并发胜者），包级 `ok … 2.864s`，与主控读数一致。
2. **行号与规模独立复核**：实现文件 969 行、测试文件 291 行（+84 从 L221 起），`git diff --stat` = pinned_cli +38 / spawn_recovery_test +84，与报告一致；D2 实为 L367–377、D1 实为 L389–397（lookup L392），位置正确：D2 在 tryLock 后、decode 前；D1 在 decode 成功后、liveness 判定前。
3. **必查 1 独立走查**：注册点 L210（onSpawn，marker 写入之后）；Start 成功路径 `advance(running)` L267 → `closeStartingMarker` L273 → 返回，此间无任何 `registry.release` 调用；`release` 仅存在于 `releaseLaunch` L691（Stop L293/L302 路径）与 `discardStartingLaunch` L704（Start 失败路径，且 guard 保证无进程在跑，L229 注释）。⇒「Start 成功返回且进程活 ⇒ handle 必在册」成立。唯一反例是 Stop 失败路径 L300–302 无条件 release，但该路径 `closeStartingMarker`（L690）已删 marker 文件 ⇒ `RecoverStarting` 无扫描目标 ⇒ D1 失效后果不可达。
4. **必查 2/3 独立复核**：registry 为纯内存（`agent_runner_process_registry.go` L110–117），`NewAgentRunnerPinnedCLILauncher` 每次新建（L90）⇒ 崩溃后必空，跨进程孤儿回收不误伤；D2 的 `fs.ErrNotExist ⇒ (false,nil)` 与「无 marker」同义（孤儿 inode 无可回收记录），其它错误包 `ErrAgentRunnerLaunchConflict` 与既有 open 失败路径（L341）分类一致，fail-closed 恰当；`TestRecoverStartingAdoptsOrphanUnix` 的 orphan 为同机进程但未注册，实测绿，证明 D1 判定粒度是「launcher 在册」而非「进程属本机」，粒度正确。
5. **必查 5 独立论证**：胜者唯一性由 `createLaunchSlot` 的 `O_CREATE|O_EXCL`（L759 区域）结构性保证；修复后 loser 的 preclaim 扫描在任何时序下都不可能 reclaim winner（同进程 D1 挡 / running 前锁挡 / running 后 marker 已删）⇒ `winners==1` 不依赖调度。

## D 残余风险

1. unix 腿无本轮执行证据（见 B1）；若 ubuntu CI 腿未重跑，unix `reclaimStartingProcess(alive=true)` 的杀进程路径未获修复后实证。
2. `closeStartingMarker` 的 `os.Remove` 失败 + 跨进程 scan 并发：marker 残留未锁时，另一进程 scan 会 reclaim 并（unix）杀掉仍活的 winner 子进程。D1 无法跨进程覆盖；概率极低且需 Remove 失败前提。
3. D2 的 `os.Stat(path)` 只证「路径有文件」、不绑定锁定的 inode（极端 TOCTOU：stat 成功但路径已被 rename 覆盖为他人 marker）。由 `clearAgentRunnerStartingMarker` 的 F-1 内容比较与 `reclaimLaunchSlot` 的 owner 校验兜底，后果有限。
4. D3 未做（用户裁定留后续）：已死但未 Stop 的 handle 仍占 registry，`live()` 会阻塞同 run root 新 Start —— 既有行为，本片未扩大。
5. E1 断言在 Windows 上经「marker/slot 被清」间接证明 D1（Windows `reclaimStartingProcess` 从不 kill），非直接 kill 断言 —— 可接受。

## E 可证伪性

- E1 四条断言均可红：去掉 D1 ⇒ unix 上 `reclaimStartingProcess(alive=true)` 杀 helper ⇒ 断言 2（alive）红；Windows 上 fail-closed 分支清 marker、夺 slot ⇒ 断言 3（marker 仍在，L281）与断言 4（slot owner）红。已在未修复代码上确定性红（与用户记录一致）。
- E1 未放宽 helper：`writeStartingMarker` 确为无锁原子写（CreateTemp+Sync+Rename，无 flock）；`registry.register` 用真实现且 phase=`starting`、identity 为真实活进程；slot 经真 `claimLaunchSlot` 写入（fields[0]==launchID，reclaim 可夺）。未发现放宽。
- 缺口：D2 无专属可证伪测试（B2）；E1 不区分 D1 的 phase 过滤变体（若实现改为「仅 starting 相位拒绝」，E1 仍绿）——该变体不可达缺陷（running 相位 marker 已删），仅登记。

**PRE-REVIEW: PASS WITH FIXES**

---

## r2

## A 结论

`PRE-REVIEW: PASS`（B3 已闭合；B1/B2/B4/B5 处置批准；附 4 条报告级必办义务 N1–N4，见 §C/§D）

## B 发现

- **B-r2-1（Low）**：处置口径「全文件 5 处统一为 `os.ErrNotExist`」与事实不符。实测：HEAD 版 5 处 `os.ErrNotExist`、0 处 `fs.` 前缀；工作树 6 处 `os.ErrNotExist`、0 处 `fs.`。准确口径应为「原 5 处保持不变 + D2 分支 1 处由 `fs.ErrNotExist` 转换为 `os.ErrNotExist`，全文件共 6 处统一」。报告措辞须改，否则 ⑥ 机械核对可能抓计数不一致。
- **B-r2-2（Info）**：「无中文注释」仅对新增行成立。`spawn_recovery_test.go:18/24` 存在 HEAD 既有中文章节号引用（`§二 Q3, §三, §四 I3/I5`、`§二 Q1, §六`，P41 遗留、指向设计文档中文章节编号），本片未引入、未改动。

无代码级缺陷，无新增 High/Medium/Low。

## C r1 各项闭合核验

**B3 — 已闭合**。逐行核对：当前 `git diff` 的 `agent_runner_pinned_cli.go` 无 import hunk ⇒ import 列表与 HEAD 完全一致，`"io/fs"` 确认删除；batch5 记录的 `errors.Is(statErr, fs.ErrNotExist)` 现为 L370 `errors.Is(statErr, os.ErrNotExist)`，其余 5 处（L339/599/666/716/734）未动；Go 标准库中 `os.ErrNotExist` 即 `fs.ErrNotExist` 的别名（同一错误值），`errors.Is` 语义完全等价 ⇒ **零行为变更**。实测：`gofmt -l internal/adapters` 空、`go vet ./...` 静默（exit 0）、`go test -count=1 ./internal/adapters/` 绿（36.178s）、并发回归 `-count=5` 绿（3.629s）。

**B1 — 批准处置，部分闭合（证据义务移交 ⑨）**。核实 `ci.yml`：test job 双 OS；Test 步 `go test -count=1 ./...` 在 ubuntu 编译并执行 `//go:build unix` 的 `TestRecoverStartingRefusesLiveHostUnix` / `TestRecoverStartingAdoptsOrphanUnix`（文件未改动、build tag 为 `unix`，ubuntu 满足）；Race 步（仅 Linux）`CGO_ENABLED=1 go test -race ./internal/adapters/...` 再次执行。本机无 unix/WSL/docker ⇒ CI 是唯一原生执行渠道，处置充分、不构成以登记代替修复（无代码可修，属证据完备性）。**N1（硬）**：⑨ 后报告须逐字摘录 ubuntu 腿 Test 与 Race 两步的包级 `ok github.com/larsonzh/prfrail/internal/adapters` 行；若任一腿红，须回 ⑤ 重开。

**B2 — 批准登记，诚实**。实测 `tryLockStartingMarker` 双平台均非阻塞（unix flock `LOCK_EX|LOCK_NB`；windows LockFileEx `FAIL_IMMEDIATELY`），open→tryLock→stat 序列无可注入边界，外部无法确定性占据 stat 前窗口；不新增生产钩子的确定性构造不存在（可挂接的 pause 点为零，钩子即给生产代码加缝，风险大于收益）。D2 分支 3 行逻辑清晰，由 D1 确定性证伪测试（E1）+ 并发回归概率性覆盖，登记理由成立。**N2**：报告「偏差与登记」须载明 D2 无专属确定性测试、仅概率覆盖，及 D3（先删后解锁）与跨进程 W1 残留移交后续片。

**B4 / B5 — 批准 Info 登记，不改代码**，符合本仓 Low/Info 由 ⑧b 登记闭合的既有实践。**N3**：B-r2-1/B-r2-2 两项措辞修正并入 B4 边界簿记；「无中文」声明按「新增行口径」表述。

**新问题核验 — 无**。`git status` 仍为 9 项本片范围（4 M + 5 ??）；新增 84 行测试与 37 行生产代码全英文、无 TODO/FIXME/XXX；E1 测试断言完备（`reclaimed==false` + 活进程 + marker + slot owner 全钉，且唯一能产出 (false,nil) 的前置分支即 D1 门，open 必成功 ⇒ 无假绿通道）。

## D 残余风险

1. B1 证据未出：⑨ ubuntu 腿结果为本片放行的**后验前置**，红即回 ⑤（N1）。
2. D2 跨进程 W1 残留（诊断 §G 已认）与 D3 未做：移交后续片/P45 ① 输入。
3. N1–N4 未落地前，⑥/⑦ 机械复核可能抓「5 处/6 处」与「无中文」措辞不一致。
4. E1 仅覆盖同进程防御场景，跨进程孤儿收养正确性依赖 ⑨ CI unix 腿。

闭集判定：`PRE-REVIEW: PASS`
# P53 FIX-GAP-77-SYMLINK-ESCAPE · ⑤ 预审报告 r1（deep-reasoner，只读，逐字归档）

## 审查范围

| 文件 | 版本锚点 | 状态 |
|---|---|---|
| `internal/adapters/bounddir_unix.go` | HEAD=dc785ed 基础上 +170/−12，diff 已逐行核对 | 生产改动 |
| `internal/adapters/bounddir_escape_unix_test.go`（改） | ③b 落盘，tag `!windows && !darwin && !p53_red_baseline` | 测试 |
| `internal/adapters/bounddir_escape_payload_unix_test.go`（新） | 共享载荷，tag 仅平台限定 | 测试 |
| `internal/adapters/bounddir_escape_redbaseline_test.go`（新） | tag `p53_red_baseline && !windows && !darwin` | 测试 |
| `internal/adapters/bounddir_symlink_storm_test.go`（新） | tag `linux`，F1/F2 | 测试 |

权威输入：`01-arch/design-r2.md`、`01-arch/rulings.md`（裁定 (a)–(e)）、`01-arch/red-team-sol-review.md`、`03-tests/red-baseline.md`、`03-tests/green-suite-report.md`、`freedom-list.md` 均通读。审查为只读：仅执行 `git status/diff`、grep、`go vet`（本机与 `GOOS=linux GOARCH=arm64/freebsd/darwin` 交叉）、`gofmt -l`、`go build ./...`、两种 tag 的 `go test -c`，未修改任何文件。

## 发现

### M1（Medium）— F1/F2 风暴测试的平台 tag 与运行时门禁不一致，linux/arm64 必然假红

- 位置：`internal/adapters/bounddir_symlink_storm_test.go:1`（`//go:build linux`）、`:174`（`if successes == 0` 断言）、`:300`（`must be EXDEV` 断言）。
- 机理：生产门禁是 `runtime.GOOS=="linux" && GOARCH=="amd64"` 才走 openat2 主路径，其余组合一律 fail-closed ENOSYS（`bounddir_unix.go:163-165`）。但 F1/F2 的 build tag 只有 `linux`，本机已实测 `GOOS=linux GOARCH=arm64 go test -c ./internal/adapters/` 编译通过、两测试进入套件。在 arm64 上：F1 所有打开被 ENOSYS 拒绝 ⇒ `successes==0` ⇒ `:174` 必然 FAIL；F2 子进程若 mount 成功，拒绝 errno 为 ENOSYS 而非 EXDEV ⇒ `:300` 必然 FAIL（若 mount 失败则恰巧 skip）。这与设计 r2 §5"闭集仅在 linux/amd64 断言，其它平台显式 t.Skipf"的口径矛盾——T51-1 有 `gap77SkipUnlessOpenat2`，F1/F2 没有对应门。
- 修复指令：将 tag 收窄为 `//go:build linux && amd64`；或（更强）在 F1/F2 测试开头复用/复制 `gap77SkipUnlessOpenat2`（含 GOARCH 判断），确保非主路径环境按"环境限制 ⇒ Skipf + 登记"处理，不留必然假红。

### M2（Medium）— fail-closed（ENOSYS）契约在测试层零钉住，且 T51-1 在 <5.6 内核上假红

- 位置：`internal/adapters/bounddir_unix.go:162-176`（fail-closed 分支，无任何测试引用；grep 证实 `openat2Disabled` 在全部测试文件中零出现）、`internal/adapters/bounddir_escape_unix_test.go:60-64`（skip 门只查 GOOS/GOARCH，不探内核）、`:98-99`（闭集断言）。
- 机理：(1) 裁定 (a) 将"不支持 openat2 ⇒ fail-closed 返回 ENOSYS 包装错误"定为本片安全语义支柱之一，但该契约没有测试：非 linux/amd64 unix 无断言、ENOSYS 置位路径无断言、`openat2Disabled` 置位后行为无断言。设计 r2 §5 原"降级模拟：包内测试置位 openat2Disabled 复跑红绿断言"未落盘（③b 登记"留 ⑧"，但 ⑧ 是文档/复核角色，不补测试）。(2) 在 linux/amd64 但内核 <5.6 的主机上，T51-1 不 skip（平台判断通过），openat2 返回 ENOSYS ⇒ 闭集 {ELOOP, ENOTDIR, EXDEV} 断言 FAIL——正确 fail-closed 行为被判为失败（假红），且未登记环境。
- 修复指令：新增包内测试（如 `bounddir_failclosed_unix_test.go`）：置位 `openat2Disabled`（包内可见）后断言 `OpenBoundDirFile` 返回 `*os.PathError` 且 `errors.Is(pathErr.Err, syscall.ENOSYS)`、`Op=="openat"`；测试末尾 `openat2Disabled.Store(false)` 复位或按名字排序置于套件末并注释包内串行前提。同时在 `gap77SkipUnlessOpenat2` 增加环境探测：以一次真实打开（或直接调 `openat2Resolve`）探测 ENOSYS，命中则 `t.Skipf` 登记"fail-closed host"，避免 <5.6 内核假红。

### L1（Low）— store.mu → dir.mu 单向嵌套的阻塞窗口，与设计陈述不完全一致

- 位置：`internal/adapters/agent_runner_replay_root.go:368-371`（`bindEvidenceDir` 持 `store.mu` 时调用 `_ = bound.Close()`）。
- 机理：`Close` 现取 `dir.mu`，构成 store.mu→dir.mu 单向嵌套；dir.mu 持有者（openat2/boundDirVerify）不请求 store.mu，无死锁。但设计 r2 §4 写"store.mu 与 dir.mu 无嵌套"，与实现有出入；若 `boundDirVerify` 的 `os.Stat(expected)`（`bounddir_unix.go:245-250`，锁内慢 I/O）偶发长阻塞，持 store.mu 的 Close 会排队。
- 修复指令：无需改代码（单向无环）；在 ⑧ 回写 SEC-02/ARCHITECTURE 时按实现如实描述锁序（store.mu→dir.mu 单向），并注明 dir.mu 持有路径不回调 store.mu。

### L2（Low）— F4 并发 Close×读取测试未落盘（设计 §5 承诺，③b 已登记）

- 位置：`03-tests/red-baseline.md` §7（"F4 并发 Close 测试：不在本轮派单，未编写"）。
- 机理：设计 r2 §5 明确"F4 并发测试：并发 Close×读取循环"。③b 全包 `-race`（100s）PASS 可作部分补偿，但无定向 Close 竞态测试。
- 修复指令：后续片（或本片 ⑧ 授权下 ③ 补单）补充：读循环 × 并发 `dir.Close()`，断言结果仅含"nil bound directory"或精确界内字节、无 panic；本片不阻塞。

### L3（Low）— 红基线变体无平台门，变异检验须限定 linux/amd64 执行

- 位置：`internal/adapters/bounddir_escape_redbaseline_test.go:1`（tag 无 amd64）。
- 机理：⑧ 若在非 linux/amd64 上跑变异检验，会得 `red baseline not reproduced`（ENOSYS）而非期望的变红，判读混淆。③b 报告的命令已写明在 VM（linux/amd64）执行。
- 修复指令：⑧ 执行文档中固化"变异检验仅 linux/amd64 内核 ≥5.6"；代码可不动。

## 逐条核对

1. **openat2 ABI**：✅ `sysOpenat2Nr=437`（linux/amd64 正确）；`openHow` 3×uint64、无填充、24B（`bounddir_unix.go:43-51,199`）；RESOLVE 常量 0x01/0x02/0x04/0x08 与 `linux/openat2.h` 一致；`Syscall6(437, dirfd, name, &how, 24, 0, 0)` 参数序与 size 传参正确；`mode=0`（无 O_CREAT/O_TMPFILE，合法）。零新依赖 ✓。
2. **运行时门禁**：✅ `bounddir_unix.go:163-165` 非 linux/amd64 直接返回 ENOSYS，绝不执行 437（FreeBSD 等平台 437 语义不同已规避）；`openat2Disabled` 仅 ENOSYS 置位（`:171-173`），其余 errno 原样上抛，语义精确。
3. **错误契约**：✅ `&os.PathError{Op:"openat", Path: clean, Err: err}` 与旧版逐字一致；ENOENT 经包装保持 `errors.Is(err, fs.ErrNotExist)`（`evidence_bind_dir_unix_test.go:94` 依赖成立）；`NewBoundFile` 三参调用与 `allowTruncation` 传递不变（`agent_runner_replay_store.go:261` 传 true 路径核对无误）。
4. **F4 生命周期**：✅ grep 证实 `dir.fd`/`dir.path` 全部读写点（`bounddir_unix.go:129,136,218,221,222,237,241,253`）均在 `dir.mu` 内；构造点（`:111`）发布前无竞争；`BoundDir` 全程指针传递（`boundDirs map[string]*BoundDir`），跨平台 vet 均无 copylocks；`Close` 幂等（fd<0 直接返回）；无反向锁嵌套死锁（见 L1 单向嵌套说明）。
5. **裁定 (a)–(c)**：✅ (a) fail-closed：ENOSYS 上抛、walk 未实现、无 `allowUnsafeFallback` 开关；(b) `boundResolveMask=0x0F` 恒含 NO_XDEV；(c) `cleanBoundRelativePath` 无任何组件计数，无深度帽。
6. **测试有效性**：F5 闭集断言正确（ENOSYS/ENOENT/EACCES 等均在闭集外防假绿）；正向对照（3 层 + 单段）真实；红基线 tag 互斥成立（两种 tag `go test -c` 均通过，③b 报告确认单 RUN 行）；F1 断言有效（每轮唯一 canary、violations 全量比对、`successes==0` 防空转、风暴后载荷完整）；F2 skip 口径符合"环境限制 ⇒ Skipf + 登记"（marker + 原因）。**缺口见 M1/M2/L2**。
7. **回归面**：✅ `evidence_bind_dir_unix_test.go` 三例（祖先替换读成功、injected.txt 走 ENOENT⇒fs.ErrNotExist、叶子软链接 ELOOP 拒绝、CloseBoundDirs 可删根）语义兼容；`bounddir_identity_test.go`（boundDirVerify 互换拒绝）兼容；`evidence_seam_signatures_test.go` 签名钉子编译通过；`bounddir_darwin.go`/`bounddir_windows.go` 零改动、各自平台 vet 通过。

## 证据

- 代码锚点：`bounddir_unix.go:43-51,66,80,111,127-136,162-176,188-204,211-253`；`bounddir_escape_unix_test.go:50-64,86-99,139-197`；`bounddir_escape_payload_unix_test.go:24-46,53-97`；`bounddir_escape_redbaseline_test.go:51-88`；`bounddir_symlink_storm_test.go:1,65-188,220-300`。
- 运行取证（本机只读）：`go vet ./internal/adapters/`=0；`gofmt -l`=空；`GOOS=linux GOARCH=arm64|freebsd|darwin go vet`=0×3；`go build ./...`=0；`GOOS=linux GOARCH=amd64 go test -c`（默认 tag 与 `-tags p53_red_baseline`）=0；`GOOS=linux GOARCH=arm64 go test -c`=0（M1 前提实证）。
- 文档锚点：`rulings.md` (a)–(e)；`design-r2.md` §1/§2/§4/§5/§6；`red-baseline.md` §2–§7（红文案逐字命中、F1 PASS、F2 SKIP 登记、F4 未派单）；`green-suite-report.md` §2–§6（全包 PASS、race PASS、红变体预期 FAIL 且 errno=ELOOP 佐证闭集）；`freedom-list.md`（"改变 fail-closed 语义属必须升级事项"——生产代码未触犯）。

发现汇总：2 Medium（M1、M2）+ 3 Low（L1、L2、L3）。生产代码（openat2 ABI、门禁、错误契约、锁生命周期、裁定落实）**零缺陷**；全部发现集中在测试口径层。

PRE-REVIEW: FINDINGS

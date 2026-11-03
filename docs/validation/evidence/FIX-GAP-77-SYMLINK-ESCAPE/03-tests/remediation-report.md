## ③c 测试角色报告 · P53 FIX-GAP-77-SYMLINK-ESCAPE

### 一、改动清单（全部测试域，4 项 5 文件）

**T1 消除 gopls 告警（红基线改为平台 tag + 运行时门控）**
- `internal/adapters/bounddir_escape_redbaseline_test.go`：`//go:build p53_red_baseline && !windows && !darwin` → `//go:build !windows && !darwin`（L1）；函数改名 `TestBoundDirMiddleSymlinkEscapesBound` → `TestBoundDirRedBaselineMiddleSymlinkEscape`（L62）；函数首部新增 `os.Getenv(redBaselineEnv) != "1"` ⇒ `t.Skipf`（L63-65，`redBaselineEnv` 常量 L51），并加 `linux/amd64` 环境跳过（非变异目标不产出假红，L71-73）；文件头改写互斥机制说明与新启用命令（L5-40）。
- `internal/adapters/bounddir_escape_unix_test.go`：tag 去掉 `&& !p53_red_baseline`（L1）；头部 L30-37、L120-124 由「同名 build tag 互斥」改为「按名区分 + 环境变量门控」。
- `internal/adapters/bounddir_escape_payload_unix_test.go`：L14-19 载荷注释同步。
- grep 复核：代码内已无 `-tags p53_red_baseline` 生效引用，仅保留 1 处「该 tag 已退役、为 no-op」的提示（L29）。

**T2（M1）风暴测试平台门与生产门禁对齐**
- `bounddir_symlink_storm_test.go`：tag `linux` → `linux && amd64`（L1）；`TestBoundDirRenameStormNeverYieldsOutsideBytes`（L87）与 `TestBoundDirBindMountImportRefused`（L238）首部调用 `gap77SkipUnlessOpenat2(t)`；bind-mount 的跳过放在 helper 分支**之后**，避免子进程单独 SKIP 被父进程当作通过（L239-244）。

**T3（M2）fail-closed 契约 + 内核探测**
- `gap77SkipUnlessOpenat2`（`bounddir_escape_unix_test.go` L66-88）：平台判断后追加真实内核探测；`gap77ProbeOpenat2`（L90-117）以 `openat2Resolve` 解析 `t.TempDir()` 下不存在名——ENOENT=主路径可用，ENOSYS=失败关闭主机 SKIP，其它 errno（seccomp 过滤等）亦 SKIP；`sync.Once` 缓存、**不改 `openat2Disabled`**、不创建任何文件（无副作用）。
- 新建 `bounddir_failclosed_unix_test.go`（tag `!windows && !darwin`）：`TestBoundDirFailClosedWithoutOpenat2` — 锁存 `openat2Disabled=true`（cleanup 复位），以内界合法名 `original` 调用，断言 `*os.PathError`、`Op=="openat"`、`Path` 为清洗名、`errors.Is(..., syscall.ENOSYS)`、无 panic、无句柄；非 linux/amd64 时清锁后断言 `openBoundArtifactLocked(-1,…)` 由**常量门**直接返回 ENOSYS（dirfd=-1 不可能进内核）；再断言锁存后二次读取仍同形拒绝。

**T4（L2）F4 并发 Close×读取**
- 同文件 `TestBoundDirConcurrentCloseAndRead`：4 读者循环 `OpenBoundDirFile/ReadAll` × 主协程 16 次 `dir.Close()`；结果只允许「界内 27 字节精确 canary」或 `adapters: nil bound directory`，EBADF/截断/界外字节/读失败/Close 失败一律登记；`ready` 通道保证至少一次成功读（10s 上限防真空通过）；终止态再断言目录已关闭且拒绝文案一致。

### 二、覆盖矩阵

| 测试点 | 测试名 | 类型 | 可证伪性 |
|---|---|---|---|
| F5 中间符号链接逃逸被拒 | `TestBoundDirMiddleSymlinkEscapesBound` | 反例+正控 | 删守卫 ⇒ 逃逸成功 ⇒ 红 |
| 错误契约/拼写矩阵/确定性 | `TestBoundDirMiddleSymlinkRefusedWithClosedErrno` | 反例+边界 | 放宽闭集 ⇒ 红 |
| M2 fail-closed ENOSYS | `TestBoundDirFailClosedWithoutOpenat2` | 反例（契约） | 删门 ⇒ open 成功 ⇒ 红 |
| F4 并发 Close×读 | `TestBoundDirConcurrentCloseAndRead` | 并发 | 去掉 `dir.mu` ⇒ race/EBADF ⇒ 红 |
| F1 rename 风暴 | `TestBoundDirRenameStormNeverYieldsOutsideBytes` | 并发 | 回归多步解析 ⇒ 界外字节 ⇒ 红 |
| F2 bind mount 导入 | `TestBoundDirBindMountImportRefused` | 环境依赖 | VM 上 SKIP（已登记原因） |
| T51-1 变异对照 | `TestBoundDirRedBaselineMiddleSymlinkEscape` | 变异检验 | 默认 SKIP |

**不可测项**：① F2 在 VM 因 `unshare(CLONE_NEWUSER)=EINVAL`、`unshare(CLONE_NEWNS)=EPERM`、无免密 sudo 而 SKIP（环境限制，非通过）；② 真实「无 openat2 内核」主机的绿集判据不可测——以 `openat2Disabled` 锁存契约测试替代（VM 内核 6.8 具备 openat2）。

### 三、验证证据

**本机（Windows，磁盘真值）**：`gofmt -l internal/adapters/` → 空；`go vet ./internal/adapters/` exit 0；`GOOS=linux GOARCH=amd64 go vet` exit 0（另 linux/arm64、darwin/amd64、freebsd/amd64 亦 exit 0）；`GOOS=linux go test -c ./internal/adapters/` exit 0（产物已删）；`go build ./...` exit 0；`go test -count=1 ./internal/adapters/` → `ok … 38.448s`；`gopls check bounddir_escape_redbaseline_test.go` 在默认与 `GOOS=linux` 两视图均 exit 0；编辑器诊断无 error。5 个测试文件 CR=0（LF）、无 BOM。

**VM 10.0.0.199（go1.22.12，`bash -lc`）**：
- 上传后 5 文件 SHA256 与本机逐一相同；生产 `bounddir_unix.go` 前后均 `88fab435…c5c44`。
- `go test -count=1 ./internal/adapters/` → `ok … 16.720s`，`EXIT=0`。
- 定向 `-v`：`PASS: MiddleSymlinkEscapesBound / MiddleSymlinkRefusedWithClosedErrno / FailClosedWithoutOpenat2 / ConcurrentCloseAndRead / RenameStormNeverYieldsOutsideBytes`；`SKIP: RedBaselineMiddleSymlinkEscape（mutation check only (⑧)…）`、`SKIP: BindMountImportRefused（userns EINVAL/NEWNS EPERM）`；F4 实证 `27 byte-exact reads, 12 refusals with "adapters: nil bound directory", 16 Close calls`；F1 `200 iterations, 6336 successful reads, 1495 refusals`。
- `CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/` → `ok … 99.138s`，无 DATA RACE，`EXIT=0`。
- `PRFRAIL_P53_RED_BASELINE=1 go test -run TestBoundDirRedBaselineMiddleSymlinkEscape -v -count=1` → `--- FAIL`，文案 `red baseline not reproduced: … (openat link/secret: too many levels of symbolic links)`，`EXIT=1`（固定树上预期对照成立）。

### 四、⑧ 变异检验新命令与判读

```
# 1) 默认套件（恒 SKIP，非 FAIL）
go test -run TestBoundDirRedBaselineMiddleSymlinkEscape -v -count=1 ./internal/adapters/
# 2) 变异后复跑（新对照命令，替代已退役的 -tags p53_red_baseline）
PRFRAIL_P53_RED_BASELINE=1 go test -run TestBoundDirRedBaselineMiddleSymlinkEscape -v -count=1 ./internal/adapters/
```

判读：出现 `GAP-77 regression: middle symlink escape succeeded and read outside the bound directory` ⇒ **变异生效且缺陷复现**；`red baseline not reproduced` ⇒ 守卫仍有效（变异未落地）；`canary mismatch` ⇒ 载荷/环境问题，不计判据；`SKIP` 且含 `the GAP-77 mutation target is linux/amd64` ⇒ 环境非变异目标（无效轮次，需换 linux/amd64 重跑）。

### 五、声明

- **生产代码零改动**：`internal/adapters/bounddir_unix.go` 本机与 VM SHA256 前后恒为 `88fab4356a113d426dd396858c9a08cf4190e3d9cd2df1243706e4170d0c5c44`；`git status` 中该文件的 M 态来自 ②a 实现轮，本轮未触。
- **未 commit、未 push**（含未推 gitee）；临时产物仅 `tmp/FIX-GAP-77-SYMLINK-ESCAPE/` 内一个已删的 `adapters.test`。
- **已声明的偏离**：`docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/green-suite-report.md` 与 `MANIFEST.md` 中 ③b 期 `-tags p53_red_baseline` 记载属历史归档，未改写（改写会伪造 ③b 的实际执行记录）；替代命令已写入代码注释与本报告，请 ⑧ 以本报告为准。

TEST: DONE

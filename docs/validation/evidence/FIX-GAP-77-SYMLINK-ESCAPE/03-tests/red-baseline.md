## P53 FIX-GAP-77-SYMLINK-ESCAPE · ③ 红测基线报告

### 1. 环境事实

| 项 | 值 |
|---|---|
| 红测主机 | 远程 VM `10.0.0.199`（用户 `larson`，`bash -lc`），内核 `6.8.0-85-generic` |
| 工具链 | `go version go1.22.12 linux/amd64`（`~/go1.22/bin/go`）；`go env GOOS GOARCH` = `linux` / `amd64` |
| 仓库 | `~/LZProjects/prfrail`，`HEAD = dc785ed607989b61cf1e75c970ca6f1b2dcfdc0e`（原克隆在 `941f60e`） |
| 基线同步方式 | **VM 无外网**（`github.com:443` 超时，`git fetch origin` 失败，`dc785ed` 对象不存在）⇒ 本机 `git bundle create ... main`（8.0 MB，`941f60e` 是 `dc785ed` 祖先）→ scp → VM `git fetch <bundle> main && git checkout dc785ed`；核对 `git status --porcelain` 为空后开始改测试。bundle 用完即删 |
| 测试文件同步 | scp 覆盖两个 `_test.go`；两侧 `sha256sum` 逐字节一致（`6b925a6f…` / `4b3d91e1…`） |
| 生产代码 | **未改**：VM `git status` 仅 ` M bounddir_escape_unix_test.go` + `?? bounddir_symlink_storm_test.go`，`git diff --stat` 只含测试文件 |

### 2. T51-1 红证据（固定文案逐字命中）

命令：`~/go1.22/bin/go test -run TestBoundDirMiddleSymlinkEscapesBound -v -count=1 ./internal/adapters/`

```
=== RUN   TestBoundDirMiddleSymlinkEscapesBound
    bounddir_escape_unix_test.go:134: GAP-77 regression: middle symlink escape succeeded and read outside the bound directory
--- FAIL: TestBoundDirMiddleSymlinkEscapesBound (0.00s)
FAIL
FAIL    github.com/larsonzh/prfrail/internal/adapters   0.013s
EXIT_CODE=1
```

`RUN` 行存在、无 `SKIP`、退出码 1；失败前已满足「打开成功 → 经句柄读出 → 逐字节等于 23 字节 `SECRET-OUTSIDE-BOUNDARY`」，故非假红。三条互斥失败文案保证判别：`red baseline not reproduced`（打开被拒）/ `canary mismatch`（读取失败、截断或字节不符）。

### 3. F1 / F2 / F5 的 HEAD 行为

**F1 rename 风暴（`TestBoundDirRenameStormNeverYieldsOutsideBytes`）→ PASS（符合预期，作为实现期护栏）**

```
--- PASS: TestBoundDirRenameStormNeverYieldsOutsideBytes (0.16s)
    F1 storm: 200 iterations, 8886 successful reads, 9116 refusals
EXIT_CODE=0
```

`CGO_ENABLED=1 CC=gcc go test -race` 同二例：`PASS`，无 race 报告（`3420 successful reads / 1246 refusals`）。

**F2 bind mount（`TestBoundDirBindMountImportRefused`）→ SKIP（环境限制，已登记）**

```
--- SKIP: TestBoundDirBindMountImportRefused (0.01s)
    skipped: host refused the setup (unshare(CLONE_NEWUSER|CLONE_NEWNS), then
    unshare(CLONE_NEWNS) as a privileged fallback: user namespace: invalid
    argument; mount namespace: operation not permitted)
EXIT_CODE=0
```

探针实测（tmp 内一次性 Go 程序，已删）：`unshare(CLONE_NEWUSER)=EINVAL`、`unshare(CLONE_NEWNS)=EPERM`、`CapEff=0`、`apparmor_restrict_unprivileged_userns=1`、`sudo -n` 需密码 ⇒ **本机无 CAP_SYS_ADMIN 且 userns 被禁**，无法构造 bind mount。测试已实现 user-ns 优先 + root 挂载 ns 回退两条路径，换到有权限环境即可真实执行（HEAD 预期红：读到界外 23 字节；修复后要求 `EXDEV`）。

**F5 绿侧（HEAD 预期红，已验证）**

```
--- FAIL: TestBoundDirMiddleSymlinkRefusedWithClosedErrno (0.00s)
    the middle symlink "link/secret" must be refused, but the open succeeded
    and read "SECRET-OUTSIDE-BOUNDARY" (read error <nil>)
EXIT_CODE=1
```

**全包回归**：`go test -count=1 ./internal/adapters/` → `FAIL` 且失败集合**恰为**上述两个预期红测，其余测试全绿（18.6 s）。

### 4. 绿侧断言设计（F5 防假绿）

三重联合断言，任何单一廉价改动都无法满足：

1. `OpenBoundDirFile(dir,"link/secret")` 必须失败；
2. 必须 `errors.As` 得 `*os.PathError`，且 `Op=="openat"`、`Path=="link/secret"`、`Err ∈ {ELOOP, ENOTDIR, EXDEV}`；`ENOENT/EACCES/EPERM/EBADF/EMFILE/EINTR/EIO/ENOSYS` 一律 FAIL（防"任意错误假绿"）；
3. 同目录同调用形状的正向对照：`sub/nested/inside`（3 层）必须打开成功、未截断、逐字节等于 `IN-SIDE-BOUNDARY`（防"一刀切拒绝假绿"）。

配套说明：闭集仅在 `linux/amd64`（openat2 主路径）断言，其它平台显式 `t.Skipf` 并说明归入 fail-closed 判据，不静默放行；变异检验（把 `openBoundArtifactLocked` 退化为单次 `openat` 应回到红文案）与降级模拟（置位 `openat2Disabled`）均需生产符号，本环禁止改生产代码且 HEAD 引用即编译失败，**留 ⑧ 复核执行**。

### 5. 测试文件改动清单（仅 `_test.go`）

| 文件 | 变更 |
|---|---|
| `internal/adapters/bounddir_escape_unix_test.go` | 删 `t.Skip`；T51-1 重写为真实红基线（句柄读取 + 23 字节逐字节比对 + 三条互斥文案）；新增 `gap77Tree/buildGAP77Tree` 共用载荷；新增 `TestBoundDirMiddleSymlinkRefusedWithClosedErrno`（F5 绿侧） |
| `internal/adapters/bounddir_symlink_storm_test.go`（新建，`//go:build linux`） | F1 `TestBoundDirRenameStormNeverYieldsOutsideBytes`（2 reader × 200 轮确定性 rename 风暴，O 侧每轮唯一 canary，成功读一律须等于界内 canary，并断言成功数 > 0 防空转、风暴后载荷完整）；F2 `TestBoundDirBindMountImportRefused`（re-exec 子进程 + user/mount namespace + `mount --bind` + `EXDEV` 断言 + 不支持时打印 `BIND-SETUP-UNSUPPORTED:` 供父进程 skip） |

**生产代码零改动**：`internal/adapters/bounddir_unix.go`、`bounddir.go`、`bounddir_clean.go`、`internal/evidence/*` 及 Windows/darwin 分支均未触碰；无 commit / push。

### 6. 覆盖矩阵

| 测试点 | 测试名 | 类型 | 源 |
|---|---|---|---|
| T51-1 红基线 | `TestBoundDirMiddleSymlinkEscapesBound` | 反例（红基线） | r2 §5 红 |
| F5 绿闭集 + 正向对照 | `TestBoundDirMiddleSymlinkRefusedWithClosedErrno` | 正例 + 反例 | r2 §5 绿 |
| F1 并发 rename | `TestBoundDirRenameStormNeverYieldsOutsideBytes` | 并发 + 边界 | r2 §1/§5、裁定 (a) |
| F2 挂载导入 | `TestBoundDirBindMountImportRefused` | 反例（环境受限） | r2 §1/§5、裁定 (b) |
| r2 §6 回归清单 | 既有 7 项测试（未改动） | 回归 | r2 §6 |

本机（Windows）冒烟：`gofmt -l` 空、`go vet ./internal/adapters/` exit 0、`GOOS=linux GOARCH=amd64 go vet` exit 0、`GOOS=linux go test -c` 构建成功（产物已删）；`go test -run TestBoundDirMiddleSymlinkEscapesBound` 在 Windows 返回 `no tests to run`（构建标签排除，符合预期）。

### 7. 未执行项与环境限制登记

- **F2 未在 VM 实跑**：无 `CAP_SYS_ADMIN`、user namespace `EINVAL`、`apparmor_restrict_unprivileged_userns=1`、无免密 sudo（详见 §3）。HEAD 红行为未观测，属**未验证项**，不得视为已覆盖。
- **F4 并发 Close 测试**：不在本轮派单（需修复后的 `dir.mu` 才有可证伪对象），未编写。
- **降级模拟 / 变异检验**：需生产符号或生产代码退化，本环越界，留 ⑧。
- **darwin / 非 linux 平台**：F1/F2 不参与构建（`//go:build linux`）；F5 在非 linux/amd64 上 skip 并打印原因。
- **VM 工作树**：保留上述两个未提交测试文件（供 ②a 实现轮直接复用），生产代码仍为 `dc785ed`；VM 上的 bundle、探针与运行脚本已全部清除，本机 `tmp/FIX-GAP-77-SYMLINK-ESCAPE/` 内临时产物（测试二进制、bundle、探针、运行脚本）亦已清除。

TEST: DONE

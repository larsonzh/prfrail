# P53 FIX-GAP-77-SYMLINK-ESCAPE — ① 架构设计 r2（红队 F1–F6 修订版）

切片：P53（M 档）。目标文件：`internal/adapters/bounddir_unix.go` 及其 unix 测试。红线：公共签名冻结、红测先行、无 `/proc` 降级、Windows/darwin 零改动。

## 0. 对 r1 的处置总纲

r1 全文作废，不再沿用。r1 的 A/B 混合、`/proc/self/fd` 探测、`O_PATH` 伪分支、深度帽 64 全部删除。r2 定稿为：

> **主路径 = `openat2(2)` 从 `dir.fd` 单次受限解析**（`RESOLVE_BENEATH|RESOLVE_NO_SYMLINKS|RESOLVE_NO_MAGICLINKS|RESOLVE_NO_XDEV`）；**降级 = 逐级 `openat(O_NOFOLLOW|O_DIRECTORY)` 游走**（即自由清单方案 A）。全程无 `/proc`、无 `O_PATH`、无深度帽。

**F1 核心裁定：包含关系的判定时点 = 解析完成的瞬间。** openat2 在**单个系统调用内**完成全部组件解析与打开，游走期间不存在可被 rename 命中的窗口；降级路径（多次 openat）只钉住"第 i 级已打开目录对象"，**不钉住祖先关系**——该差异在下文降级矩阵中如实写明，不再宣称等价。

## 1. 威胁模型（重写）

| 载荷 | 构造 | 主路径（openat2） | 降级（walk） |
|---|---|---|---|
| 中间软链接（GAP-77 本体） | `D/link→O`，读 `link/secret` | `NO_SYMLINKS` 在任一组件遇链接即 ELOOP，未返回句柄 | 中间组件 `openat(O_NOFOLLOW)` ⇒ ELOOP |
| 叶子软链接 | `D/g→O/f` | ELOOP（`NO_SYMLINKS` 先行，`O_NOFOLLOW` 兜底） | ELOOP |
| 并发 `rename(D/sub→O/sub)` 后注入叶子（F1） | 两次系统调用间移出中间目录再放叶子 | **闭合**：解析+打开同一 syscall，无中间窗口；成功后 fd 钉住对象 | **不闭合**：残余竞态窗口（前提=攻击者对 D 有写权限+精确命中窗口） |
| `RENAME_EXCHANGE` 原子换目录树 | 换入 D 内的树按"解析瞬间位于 D 内"计，非越界；换出+注入同左行 | 闭合（同上） | 不闭合（同上） |
| bind mount 导入（F2） | `mount --bind O D/m`，读 `m/secret`；同文件系统 bind 不可用 `st_dev` 排除 | `NO_XDEV` 跨越任何挂载点（含 bind）即 EXDEV；`BENEATH` **单独不挡**子挂载（man 页语义，红队 F2 成立） | **不闭合**：挂载点不是 symlink，`O_NOFOLLOW`/`O_DIRECTORY` 均放行 |
| 目录硬链接 / 叶子硬链接 | 目录硬链接内核拒 EPERM；叶子硬链接的名字在 D 内 | 界内名字 ⇒ 界内（来源隔离不承诺） | 同左 |
| `/proc/self/fd/N` 魔链输入 | 仅存在于 procfs；若被 bind 进 D，`NO_XDEV` 已先拒 | `NO_MAGICLINKS` 显式 ELOOP | 中间组件若为链接 ⇒ ELOOP |
| `EAGAIN`（BENEATH 的 `..` 竞态） | clean 守卫已保证净化后无 `..` ⇒ 不可达；若出现 | 拒绝，不重试（fail-closed） | n/a |
| `EINTR`/seccomp `EPERM` | 信号/过滤 | 拒绝（fail-closed）；仅 `ENOSYS` 触发降级 | 拒绝 |

**已明确的边界（写进 SEC-02 口径）**：绑定读取拒绝的对象越界判定在"解析完成瞬间、命名空间内 dirfd 后代"；叶子 fd 一旦返回，其对象身份固定（后续 rename 不影响已开 fd，P42 既有语义保留）。硬链接别名、换入 D 内的新树属"界内名字"，不在本片威胁模型。

## 2. 方案定稿与降级矩阵

**主路径（Linux ≥5.6, amd64）**：一次 `syscall.Syscall6(437, dirfd, path, &how, 24, 0, 0)`，`how = {flags: O_RDONLY|O_CLOEXEC|O_NOFOLLOW, mode: 0, resolve: 0x0F}`（`RESOLVE_NO_XDEV 0x01 | NO_MAGICLINKS 0x02 | NO_SYMLINKS 0x04 | BENEATH 0x08`）。成功返回普通 fd 交 `NewBoundFile`；失败返回 errno。不重试 EAGAIN。

**降级**：`openat2Disabled atomic.Bool` 初始 false；首次 openat2 返回 `ENOSYS` ⇒ 置位并立即以 walk 重试本次打开，此后永久走 walk。walk：逐组件 `openat(cur, comp, O_RDONLY|O_DIRECTORY|O_NOFOLLOW|O_CLOEXEC)`，叶子 `openat(last, O_RDONLY|O_CLOEXEC|O_NOFOLLOW)`，失败即关闭链上全部中间 fd 后返回。

**降级矩阵（保证如实陈述）**：

| 环境 | GAP-77 静态载荷 | F1 并发 rename | F2 挂载导入 | 合法界内读取 |
|---|---|---|---|---|
| Linux ≥5.6 amd64（主路径） | 阻断（ELOOP） | 阻断（单次解析） | 阻断（EXDEV） | 放行 |
| 降级 walk（<5.6、非 amd64、其它 unix） | 阻断（ELOOP/ENOTDIR） | **不承诺**（有窗口） | **不承诺** | 放行 |

若 §8 裁定 (a) 与 (b) 均为"入威胁模型"，降级路径降为**直接拒绝绑定读取**（fail-closed，返回 `ENOSYS` 形状错误），仅主路径提供服务；由裁定决定。

**红线"无 /proc 降级"落地**：主路径与降级均不依赖 `/proc`；r1 的探测逻辑整体删除，`OpenBoundDir` 零改动。**F3 矛盾消除**：无 A/B 双分支、无魔链重开、不用 `O_PATH`；叶子拒绝由 `NO_SYMLINKS`/`O_NOFOLLOW` 产生 ELOOP，不存在"返回链接自身 fd 须 fstat 拒绝"的分支。

## 3. 工具链事实（go1.22.12 linux/amd64，已联网核对源码原文）

| 符号 | 标准库 `syscall` | 处置 |
|---|---|---|
| `O_DIRECTORY(0x10000)`、`O_NOFOLLOW(0x20000)`、`O_CLOEXEC`、`O_RDONLY` | **已导出**（zerrors_linux_amd64.go） | 直接用 |
| `ELOOP/ENOTDIR/EXDEV/ENOSYS/E2BIG/EINVAL/EAGAIN` | **已导出** | 直接用 |
| `O_PATH(0x200000)` | **未导出**（zerrors 无此项） | **不使用**（本设计无需 O_PATH） |
| `SYS_OPENAT2(437)` | **未导出**（zsysnum_linux_amd64.go 止于 `SYS_PRLIMIT64=302`） | 文件内私有 `const sysOpenat2Nr = 437` |
| `RESOLVE_*` 四个值 | 未导出 | 文件内私有 const（值见 §2，取自 `linux/openat2.h`） |
| `open_how` 结构（u64×3，24 字节） | 无 | 私有 `type openHow struct{ flags, mode, resolve uint64 }` |

- **零依赖**：`go.mod` 仅有 `golang.org/x/text` 与 `gopkg.in/yaml.v3`，**无 `golang.org/x/sys`**；不新增依赖、不改 go.mod。调用经 `syscall.Syscall6`，**glibc 无 openat2 包装（man 页明示）、musl 亦无**，纯内核 ABI，与 libc 无关。
- **运行时门禁**：仅 `runtime.GOOS=="linux" && runtime.GOARCH=="amd64"` 时尝试 openat2（FreeBSD 等系统调用号 437 含义不同，绝不可裸调）；其余组合直接 walk。保持单文件改动（不新增 `bounddir_openat2_linux.go`，守住自由清单的文件边界）。
- **内核矩阵**：≥5.6 主路径；<5.6 首次调用 ENOSYS 自动降级；S1 支持矩阵仅 linux/amd64 与 windows/amd64，其它 unix 为编译级支持。

## 4. 函数级设计（仅 `bounddir_unix.go`；签名全部不变）

```go
type BoundDir struct {
    fd   int
    path string
    mu   sync.Mutex // F4：私有互斥，覆盖游走/verify/Close
}

// 私有常量：sysOpenat2Nr=437；resolveNoXdev=0x01；resolveNoMagic=0x02；
// resolveNoSymlinks=0x04；resolveBeneath=0x08；boundResolveMask=0x0F
type openHow struct{ flags, mode, resolve uint64 } // 24 字节，linux/openat2.h 布局
var openat2Disabled atomic.Bool

// openBoundDirFile：签名不变；新增加锁与分派
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error) {
    if dir == nil { return nil, fmt.Errorf("adapters: nil bound directory") }
    dir.mu.Lock(); defer dir.mu.Unlock()
    if dir.fd < 0 { return nil, fmt.Errorf("adapters: nil bound directory") }
    clean, err := cleanBoundRelativePath(rel)
    if err != nil { return nil, err }
    fd, err := openBoundArtifactLocked(dir.fd, clean)
    if err != nil { return nil, &os.PathError{Op: "openat", Path: clean, Err: err} } // 包装不变
    return evidence.NewBoundFile(os.NewFile(uintptr(fd), clean), bound, allowTruncation)
}

func openBoundArtifactLocked(dirfd int, clean string) (int, error) {
    if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" && !openat2Disabled.Load() {
        fd, err := openat2Resolve(dirfd, clean)
        if err == syscall.ENOSYS { openat2Disabled.Store(true); return walkBoundComponents(dirfd, clean) }
        return fd, err // ELOOP/EXDEV/ENOTDIR/ENOENT/EACCES/EAGAIN/EINTR... 一律原样上抛
    }
    return walkBoundComponents(dirfd, clean)
}

func openat2Resolve(dirfd int, clean string) (int, error) {
    p, err := syscall.BytePtrFromString(clean)
    if err != nil { return -1, err }
    how := openHow{flags: syscall.O_RDONLY | syscall.O_CLOEXEC | syscall.O_NOFOLLOW, resolve: boundResolveMask}
    fd, _, errno := syscall.Syscall6(sysOpenat2Nr,
        uintptr(dirfd), uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(&how)), uintptr(unsafe.Sizeof(how)), 0, 0)
    if errno != 0 { return -1, errno }
    return int(fd), nil
}

func walkBoundComponents(dirfd int, clean string) (int, error) {
    comps := strings.Split(clean, string(os.PathSeparator)) // Clean 后无空段、无 "."/".."
    cur, opened := dirfd, []int{}
    for i := 0; i < len(comps)-1; i++ {
        next, err := syscall.Openat(cur, comps[i],
            syscall.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
        if err != nil { closeChain(opened); return -1, err }
        opened = append(opened, next); cur = next
    }
    fd, err := syscall.Openat(cur, comps[len(comps)-1], syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0)
    closeChain(opened)
    return fd, err
}

// Close：加锁幂等；boundDirVerify：加锁读 fd（agent_runner_replay_root.go:360 调用点兼容）
```

**要点**：① 新增 import 仅 stdlib：`runtime`、`strings`、`sync`、`sync/atomic`、`unsafe`。② 错误包装与哨兵闭集不变（`evidence_bind_dir_unix_test.go:94` 的 `fs.ErrNotExist` 依赖保持）。③ `OpenBoundDir` 逐字不动。④ walk 中间目录要求 `+r`（`O_RDONLY`）：证据目录为进程自建（可读），仅降级环境存在此收窄，属可用性差异非安全放行。⑤ 无深度帽：openat2 无组件数限制；walk 由输入长度界定；内核字符串上限为天然上限。

**F4 生命周期闭合**：`dir.mu` 覆盖 `openBoundDirFile` 全游走、`boundDirVerify`、`Close`；持锁期间 `fd` 不会被关闭，fd 数字复用不可发生；`Close` 幂等置 `-1`。调用方兼容：`readBoundEvidence`（`agent_runner_replay_store.go:258`）先短暂取 `store.mu` 再释放后调用，`store.mu` 与 `dir.mu` 无嵌套；`BoundDir` 全程指针传递，`go vet` copylocks 无告警；`evidence_seam_signatures_test.go` 的签名钉子继续通过。

## 5. 红绿判据（F5 加固）

- **红（修复前 VM 实跑，自由清单时序强制）**：仅重写 T51-1 并删 `t.Skip`（先于任何生产改动）；HEAD 上 `OpenBoundDirFile(dir,"link/secret")` 成功 ⇒ 测试**必须经句柄逐字节读出** `SECRET-OUTSIDE-BOUNDARY`（23 字节比对），再以唯一文案 `GAP-77 regression: middle symlink escape succeeded and read outside the bound directory` 失败。若打开失败或 canary 不符 ⇒ 以**不同文案**（`red baseline not reproduced` / `canary mismatch`）失败，杜绝假红。归档：go1.22.12、`bash -lc`、`-run TestBoundDirMiddleSymlinkEscapesBound -v -count=1 ./internal/adapters/`、退出码、RUN 行。
- **绿（同份测试）**：打开必须失败且 `errors.As` 得 `*os.PathError`，其 `Err` 属**安全拒绝闭集 {ELOOP, ENOTDIR, EXDEV}**（不含 ENOENT/EACCES/EPERM/EBADF/EMFILE/EINTR/EIO——命中这些即 FAIL，防"任意错误假绿"）；随后**同条件正向读取** `sub/inside` 成功且逐字节等于 `IN-SIDE-BOUNDARY`（证明非一刀切拒绝）。预期主/降级路径对 GAP-77 载荷均给 ELOOP，断言唯一稳定。
- **变异检验（⑧ 复核）**：临时把 `openBoundArtifactLocked` 退化为单次 `openat(O_NOFOLLOW)` ⇒ 测试必须回到红文案；恢复后转绿。
- **降级模拟**：包内测试置位 `openat2Disabled` 复跑红绿断言（无生产钩子）。
- **F1/F2 实测**：F1 用确定性 rename 风暴（O 侧每次迭代唯一 canary，断言任何成功读取均非 O 值）；F2 测试 `unshare(CLONE_NEWNS=0x20000 本地常量, SYS_UNSHARE=272 已导出)+mount --bind`，`EPERM` 时 `t.Skipf` 并登记环境，成功时断言 EXDEV；两者证明机制已执行，核心保证以 §2 单系统调用/NO_XDEV 语义为判据。
- **F4 并发测试**：并发 `Close`×读取循环，断言结果只含"nil bound directory"错误或精确界内字节，无 panic；VM `-race` 门禁。

## 6. 回归清单（必须全绿）

`TestBoundDirReadSurvivesAncestorReplacement`、`TestBoundDirRejectsSymlinkLeaf`、`TestBoundDirVerifyRefusesSwappedDirectory`、`TestCleanBoundRelativePathRefusesEscapeShapes`、`evidence_bind_dir_unix_test.go` 三例（f.txt/fs.ErrNotExist/CloseBoundDirs 可删根）、`agent_runner_replay_store` 相关读取测试、`evidence_seam_signatures_test.go`。门禁：`gofmt -l` 空、`go build/vet/test ./...` 全绿、`GOOS=linux go vet ./internal/adapters/`、VM `go1.22.12` 全量 + `-race ./internal/adapters/`。

## 7. F1–F6 逐条处置表

| 发现 | 处置 | 落点 |
|---|---|---|
| F1 High（fd 不钉从属关系） | 主路径改 openat2 单次解析（无中间窗口）；判定时点明写为解析完成瞬间；降级如实声明不承诺 | §0/§1/§2 |
| F2 High（挂载导入） | `RESOLVE_NO_XDEV` 主路径闭合；`BENEATH` 单独不挡已写明；降级不闭合如实登记；是否入威胁模型交裁定 (b) | §1/§2/§8 |
| F3 Medium（A/B 矛盾、O_PATH 不编译） | 单一算法定稿；O_PATH//proc 探测/伪分支/魔链重开全部删除；"ELOOP 非必然"已修正；go1.22.12 事实列 §3 | §0/§2/§3 |
| F4 Medium（fd 生命周期） | 私有 `sync.Mutex` 覆盖游走/verify/Close，签名不变；调用方兼容性论证 | §4 |
| F5 Medium（假绿） | 红=canary 逐字节+专属文案；绿=闭集 {ELOOP,ENOTDIR,EXDEV}+正向读取+变异检验+降级模拟 | §5 |
| F6 Medium（深度帽） | 删除 64 帽；自然上限=输入长度+内核字符串限制；是否引入交裁定 (c)（建议不引入） | §4/§8 |

## 8. 用户裁定请求（① 验收时拍板）

1. **(a) 威胁边界是否含"绑定根目录树可被外部写入者 rename/替换"？** 含 ⇒ openat2 单次解析（或等价原子保证）为硬要求，降级环境应 fail-closed 拒绝而非 walk；不含 ⇒ walk 降级保留，其"只抗静态载荷"的前提写入 SEC-02。
2. **(b) 挂载导入是否属于信任前提？** 属攻击面 ⇒ `RESOLVE_NO_XDEV` 必选（已内置），并接受"界内合法子挂载也被拒"的可用性代价；属信任前提 ⇒ 可去掉 NO_XDEV（设计默认保留）。
3. **(c) 深度帽**：设计建议**不引入**（删除 r1 的 64），确认接受。
4. **(d) 文档同步归属**：CHECKLIST v42 与自由清单要求本片 ⑧ 同步 `SECURITY{,_EN}` SEC-02、`TEST_STRATEGY{,_EN}` AT-02、`ARCHITECTURE{,_EN}` §5.1、`DEV_PLAN{,_EN}`，确认纳入派单（代码改动仍仅 `bounddir_unix.go`+测试）。
5. **(e) 方案载体签认**：自由清单原列二选一（A/B），openat2 为红队指令评估的第三候选且为 r2 主路径，需用户签认超出清单。

## 9. 结论（一句话派单）

在 `bounddir_unix.go` 内：`openBoundDirFile` 加私有 `sync.Mutex`，解析改为**openat2 单次受限解析**（`RESOLVE_BENEATH|NO_SYMLINKS|NO_MAGICLINKS|NO_XDEV`，仅 linux/amd64 运行时启用，常量 437/0x0F 与 24 字节 `open_how` 文件内私有定义、零新依赖）；首次 ENOSYS 永久降级为**逐级 `openat(O_NOFOLLOW|O_DIRECTORY)` 游走**（不承诺并发 rename 与挂载导入，待 §8(a)(b) 裁定）；无 `/proc`、无 `O_PATH`、无深度帽；公共签名与哨兵闭集不变，T51-1 按 §5 红绿判据重写并移除 skip，Windows/darwin 零改动。
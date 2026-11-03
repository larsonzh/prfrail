# P53 GLM 批量辅助 1/2 · openat 用法与调用面枚举（只读事实清单，逐字归档）

枚举完成（只读，未做任何修改）。`tools/` 下所有 Go 文件均无匹配，全部命中集中在 `internal/adapters/`。

## 1. openat / openat2 / openat2Resolve / openBoundArtifactLocked 出现点

**生产代码（4 文件）**
- `bounddir_unix.go`：`:9,11,19-25` 注释（openat2 主路径取代旧单次 openat+O_NOFOLLOW；ENOSYS fail-closed）；`:42-47` `sysOpenat2Nr=437` 私有常量（go1.22 未导出）；`:69` openHow 镜像 open_how ABI；`:80` `openat2Disabled atomic.Bool` 闩锁；`:136` 调用 `openBoundArtifactLocked`；`:141` 拒绝时 `PathError{Op:"openat"}`；`:150,162` `openBoundArtifactLocked` 定义（持 `dir.mu` 下按 dirfd 解析）；`:152-158` 设计注释（裁定 a：无 walk 降级）；`:166-172` 闩锁检查+`openat2Resolve` 调用+ENOSYS 闩存；`:183-195` `openat2Resolve` 定义（`RESOLVE_BENEATH|NO_SYMLINKS|NO_MAGICLINKS|NO_XDEV`，flags 另保 `O_NOFOLLOW|O_CLOEXEC`）
- `bounddir.go:15`：共享表面文档（Unix=openat+O_NOFOLLOW+O_CLOEXEC）
- `bounddir_windows.go:8,101`：注释（Windows 无 openat，靠句柄派生名）
- `bounddir_darwin.go:7,32`：注释（stdlib 无 syscall.Openat ⇒ 拒绝绑定）

**测试代码（5 文件）**
- `bounddir_escape_unix_test.go:61-114` `gap77SkipUnlessOpenat2`/`gap77ProbeOpenat2`（内核一次性探测）；`:112` 调 `openat2Resolve` 探测；`:181-223` op="openat" 契约断言
- `bounddir_failclosed_unix_test.go:44-129` ENOSYS 契约（`:55,99,105` 闩锁 `openat2Disabled`）；`:100` 直接调 `openBoundArtifactLocked(-1, rel)` 验 nil-fd 拒绝
- `bounddir_symlink_storm_test.go:27,67,217`（F1/F2 前 skip 门）
- `bounddir_escape_redbaseline_test.go:13,68` 注释（红基线=回退单 openat）
- `evidence_bind_dir_unix_test.go:10,16,99,133` I3 叶符号链接 O_NOFOLLOW 拒绝钉扎

## 2. openBoundDirFile / OpenBoundDirFile 调用方

**生产（仅 2 处）**
- `agent_runner_replay_store.go:261` → `openBoundDirFile(dir, relative, bound, true)`（replay 消费读，allowTruncation=true 日志类）
- `agent_runner_replay_root.go:356` → `OpenBoundDir(base)`（replay 根证据目录绑定）

**测试**：`bounddir_escape_unix_test.go:137,157,210,241`；`bounddir_failclosed_unix_test.go:79,110,194,257`；`bounddir_symlink_storm_test.go:124,189,299`；`bounddir_escape_redbaseline_test.go:84`；`bounddir_escape_windows_test.go:69,110`；`evidence_bind_dir_unix_test.go:70,94,128`；`evidence_bind_dir_windows_test.go:64,88`；`evidence_seam_signatures_test.go:40`（签名钉扎）

**实现/转发**：`bounddir_unix.go:117-123`、`bounddir_windows.go:72-78`、`bounddir_darwin.go:38-43`（三平台同名转发到 `openBoundDirFile`）

## 3. bounddir_unix.go 之外的目录绑定解析

- `bounddir_windows.go` — Windows 侧同表面实现（share-delete 句柄 + 句柄派生真实名，非 openat）
- `bounddir_darwin.go` — darwin fail-closed 占位（`ErrBoundDirUnsupported`，永不解析）
- `bounddir.go` — 共享表面/错误定义
- 其余：`internal/evidence/boundread_test.go:17` 仅注释引用 CloseBoundDirs；`tools/` 全部 Go 源（agent-stub、agent-probe、b4-e2e）零匹配。除上述平台变体外，无其它文件做绑定解析。

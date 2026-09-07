# T006 快照与存储底座验证报告

[English](t006-snapshot_EN.md)

日期：2026-09-07。结论：`T006 COMPLETE`。本报告覆盖快照捕获、内容寻址存储（CAS）、无 Git 恢复及 AT-02 验收；任务链执行与事务 journal 属于 T008/T009。

## 实现范围

- `types.go`：`snapshot-manifest` 数据模型与严格校验，支持 baseline 和 candidate 模式；Unicode 码点字典序排序、目录闭包保证、Windows 保留名（CON/PRN/AUX/NUL/COM/LPT 等及带后缀变体）拒绝、大小写冲突检测、软链接逃逸防御及域分离摘要 `proofrail:snapshot-manifest:1\n`。
- `store.go`：内容寻址存储（CAS），支持临时文件刷盘与原子重命名落位、内容去重、读取时哈希完整性校验与损坏隔离、存储字节配额强制检查，以及确定性 GC dry-run（严格保留可达和审计锁定对象，不执行物理删除）。
- `capture.go`：只读遍历源工作区并捕获未提交/未跟踪文件，源树元数据和内容保持严格不变；支持默认（`.git/**`、`tmp/**`、`.prfrail/**`）及用户自定义排除、文件角色识别、硬链接分组检测、环境探测（OS/Arch/大小写敏感/环境变量名白名单）与捕获中文件变动检测（mtime/size 二次校验）。
- `restore.go`：独立工作区物化，先执行完整引用闭包预检（缺对象即 fail-close），在不调用任何 Git 命令、不依赖 Git 历史的前提下从 CAS 纯内容重建目录层级、普通文件权限（可执行/只读）、硬链接与符号链接。
- `reparse_windows.go` / `reparse_other.go`：平台特定的重解析点与目录联接（junction）检测，非标链接和不可移植重解析点一律 fail-close。

## 验收矩阵（AT-02）

- 含未提交文件树捕获：已验证包含修改、未跟踪文件和锁文件，生成合规 manifest 与内容哈希。
- 源树不变：捕获前后源工作区文件大小、mtime、内容与权限比对完全一致。
- 保留名：CON、PRN、AUX、NUL、COM1-9、LPT1-9 等及其带扩展名文件被严格拒绝（`ErrReservedPath`）。
- 大小写冲突：大小写不同但同名路径碰撞被检测并拒绝（`ErrCaseCollision`）。
- 符号链接：工作区内相对符号链接正确捕获与物化，指向父级以上的逃逸目标被阻断；Windows 无特权提示跳过，Ubuntu 实机完全通过。
- 重解析点/junction：非标重解析点与目录联接被识别并阻断（`ErrUnsupportedLink`）。
- 硬链接：同 inode/文件索引的多个硬链接被识别为同一 `hardlinkGroup`，恢复时优先硬链接物化。
- 长路径：1024 字节以内合法相对路径支持，超长路径被拒绝（`ErrInvalidPath`）。
- 捕获中变动：文件在捕获读取期间发生大小或修改时间变动被检测并终止（`ErrConcurrentChange`）。
- 配额不足：写入或捕获超过 `maxBytesQuota` 时立即 fail-close 并报告 `ErrQuotaExceeded`。
- GC dry-run：可达与锁定对象受严格保护，仅标记超期孤立对象，物理文件零删除。
- 恢复闭包验证：目标非空阻断、对象缺失阻断、哈希篡改阻断，恢复产物逐字节一致。

## 双平台验证记录

```powershell
go test -v ./internal/snapshot
go vet ./internal/snapshot
go test -cover ./internal/snapshot
```

Windows 11 (NTFS)：19 项测试中 18 项 PASS，1 项符号链接因无管理员特权依规范 SKIP；语句覆盖率 76.3%。
Ubuntu 24.04 (ext4)：经 SSH 在远程测试环境中运行，19 项测试全部 100% PASS（包括原生符号链接创建与解析）。

## 边界

T006 仅建立快照与存储基础，不执行任何 Git commit/push，不从 Git 历史恢复数据，不对存储执行实际自动清理。工作区写事务与崩溃回滚点由 T008 承接。

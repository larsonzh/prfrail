# P53 FIX-GAP-77-SYMLINK-ESCAPE · ③b 测试报告（消除 T51-1 恒红）

## 1. 改动清单（仅 `_test.go`）

| 文件 | 状态 | 内容 |
|---|---|---|
| `internal/adapters/bounddir_escape_unix_test.go` | 改 | 构建标签改为 `!windows && !darwin && !p53_red_baseline`；删除红基线三分支逻辑；`TestBoundDirMiddleSymlinkEscapesBound` 改为**修复后 PASS 形态**；`TestBoundDirMiddleSymlinkRefusedWithClosedErrno` 改为**互补断言**（不再重复闭集） |
| `internal/adapters/bounddir_escape_payload_unix_test.go` | 新建 | 共享载荷（`gap77OutsideCanary` 23 字节 / `gap77InsideCanary` / `gap77RegressionMessage` / `gap77Tree` / `buildGAP77Tree`）。标签只做平台限定（`!windows && !darwin`），**默认构建与 `-tags p53_red_baseline` 两种模式都编译**，保证变异运行比对的是同一份字节 |
| `internal/adapters/bounddir_escape_redbaseline_test.go` | 新建 | `//go:build p53_red_baseline && !windows && !darwin`；原样搬运三条互斥逻辑（`GAP-77 regression…` / `red baseline not reproduced` / `canary mismatch`），函数名与默认文件**同名**，靠标签互斥，同一时刻只编译一份 |

**命名与合并关系**（已在两处文件头注释写明）：`TestBoundDirMiddleSymlinkEscapesBound` 保留原名以维持 T51-1 证据链可追溯，但语义从「红基线」翻转为「常驻回归」；`TestBoundDirMiddleSymlinkRefusedWithClosedErrno` 保留为**互补轴**，不重复闭集。

**断言分工（避免冗余）**

| 断言 | 归属 |
|---|---|
| 打开必须失败 + `errors.As` → `*os.PathError` + `Err ∈ {ELOOP, ENOTDIR, EXDEV}` | `…EscapesBound`（唯一一处） |
| 正向对照 `sub/nested/inside` 逐字节 == `IN-SIDE-BOUNDARY`、未截断 | `…EscapesBound`（按要求保留） |
| `Op=="openat"`、`Path==清理后名`、四种拼写（`link/secret`、`./link/secret`、`link/../link/secret`、`sub/../link/secret`）**errno 漂移检测**、重复 2 次仍拒（防一次性拒绝） | `…RefusedWithClosedErrno` |
| 单段界内名 `original` 逐字节读取（互补三段对照，防"只放行子目录"） | `…RefusedWithClosedErrno` |

两个绿测在非 `linux/amd64` 上 `t.Skipf`（openat2 主路径判据，与 ③ 基线理由一致），不静默放行。

## 2. VM 验证证据（`10.0.0.199`，`bash -lc`，go1.22.12 linux/amd64，`~/LZProjects/prfrail`）

工作树：`M bounddir_escape_unix_test.go` + `M bounddir_unix.go`（②a 生产修复，**非本轮改动**）+ 三个未跟踪测试文件；`git diff --stat` 与本机逐行一致。三文件 scp 后两侧 SHA256 全等（`f90a84e7…` / `85d59bca…` / `78cf55af…`）。

**① 全包（确认无恒红）**
```
ok      github.com/larsonzh/prfrail/internal/adapters   17.255s
GROUP1_EXIT=0
```

**② 聚焦三测**
```
=== RUN   TestBoundDirMiddleSymlinkEscapesBound
--- PASS: TestBoundDirMiddleSymlinkEscapesBound (0.00s)
=== RUN   TestBoundDirMiddleSymlinkRefusedWithClosedErrno
--- PASS: TestBoundDirMiddleSymlinkRefusedWithClosedErrno (0.00s)
=== RUN   TestBoundDirRenameStormNeverYieldsOutsideBytes
    bounddir_symlink_storm_test.go:169: F1 storm: 200 iterations, 4992 successful reads, 1395 refusals
--- PASS: TestBoundDirRenameStormNeverYieldsOutsideBytes (0.09s)
EXIT_CODE=0
```
（均为 `PASS` 而非 `SKIP`，闭集断言真实执行。）

**③ race**
```
ok      github.com/larsonzh/prfrail/internal/adapters   100.072s
GROUP3_EXIT=0
```
无 `WARNING: DATA RACE`。

**④ 红基线变体（预期 FAIL，⑧ 对照基准）**
```
=== RUN   TestBoundDirMiddleSymlinkEscapesBound
    bounddir_escape_redbaseline_test.go:61: red baseline not reproduced: OpenBoundDirFile("link/secret")
    was refused on the vulnerable tree (openat link/secret: too many levels of symbolic links); ...
--- FAIL: TestBoundDirMiddleSymlinkEscapesBound (0.00s)
EXIT_CODE=1
```
文案命中 `red baseline not reproduced`，且括号内 errno 恰为 **ELOOP** —— 同时独立印证安全闭集判据。该次运行只出现一个 `RUN` 行，证明标签互斥生效（未出现重复定义编译错误）。

## 3. 红基线变体启用方式（交 ⑧）

```bash
# 变异后复跑（当前修复树＝预期 FAIL: red baseline not reproduced）
~/go1.22/bin/go test -tags p53_red_baseline \
  -run TestBoundDirMiddleSymlinkEscapesBound -v -count=1 ./internal/adapters/
```
判读规则（写在文件头）：打开被拒 ⇒ `red baseline not reproduced`；读取失败/截断/字节不符 ⇒ `canary mismatch`；**只有**打开成功且经句柄逐字节读到 23 字节 `SECRET-OUTSIDE-BOUNDARY` 才落 `GAP-77 regression: …`。变异（如把 `openBoundArtifactLocked` 退化回单次 `openat(O_NOFOLLOW)`）后该运行应出现 `GAP-77 regression…`。默认套件不含此文件，`go test ./internal/adapters/` 不受影响。

## 4. 本机（Windows）冒烟

| 检查 | 结果 |
|---|---|
| `gofmt -l internal/adapters/` | 空（退出 0） |
| `go vet ./internal/adapters/` | 退出 0 |
| `GOOS=linux GOARCH=amd64 go vet ./internal/adapters/` | 退出 0 |
| `GOOS=linux go test -c ./internal/adapters/` | 退出 0（产物已删除） |
| `GOOS=linux go vet` / `test -c` + `-tags p53_red_baseline` | 退出 0（红变体可编译） |
| `go build ./...` | 退出 0 |
| 编码/行尾 | 三个文件 UTF-8 **无 BOM** + **LF**（新建文件首次落盘为 CRLF，已机械归一，CR=0；无值语义变更） |

## 5. 覆盖矩阵

| 测试点 | 测试名 | 类型 | 套件 |
|---|---|---|---|
| T51-1 逃逸必须被拒 | `TestBoundDirMiddleSymlinkEscapesBound` | 反例 | 默认（常驻绿） |
| T51-1 红基线（变异检验） | 同名（redbaseline 文件） | 反例 | `-tags p53_red_baseline` |
| F5 错误契约 + 拼写矩阵 + 确定性 + 单段正向 | `TestBoundDirMiddleSymlinkRefusedWithClosedErrno` | 正例 + 反例 + 边界 | 默认 |
| F1 并发 rename 风暴 | `TestBoundDirRenameStormNeverYieldsOutsideBytes` | 并发 | 默认（未改动） |
| F2 挂载导入 | `TestBoundDirBindMountImportRefused` | 反例（环境受限） | 默认（未改动） |

## 6. 不可测 / 未覆盖项（显式登记）

- **变异检验本身**不可由本环执行：需改生产代码（本环禁止），已交 ⑧ 用上述 tag 命令复跑。
- **崩溃类**（进程被 kill 时句柄/描述符状态）本片无载荷入口，未覆盖；无对应测试点。
- **F2 bind mount** 在 VM 上仍 `SKIP`（无 `CAP_SYS_ADMIN`、userns `EINVAL`），属环境受限项，未验证。

## 7. 声明

- 生产代码**零改动**：`bounddir_unix.go` 前后 SHA256 = `88fab435…`，本机与 VM 逐字节相同；`bounddir.go`/`bounddir_clean.go`/`bounddir_windows.go`/`bounddir_darwin.go`/`internal/evidence/*` 未触碰。
- 未执行 `git commit` / `git push`；未使用 `git add`。
- 本机新增临时件仅 `tmp/FIX-GAP-77-SYMLINK-ESCAPE/` 下测试二进制，**已删除**；VM 侧无新增临时文件。

TEST: DONE

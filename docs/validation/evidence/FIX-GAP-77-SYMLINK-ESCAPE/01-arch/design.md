# P53 FIX-GAP-77-SYMLINK-ESCAPE — ① 架构设计（bounddir 中间符号链接逃逸修复）

切片：`FIX-GAP-77-SYMLINK-ESCAPE`（P53，M 档）。目标文件：`internal/adapters/bounddir_unix.go`（+ 其测试）。范围红线：只动 Unix 文件；`bounddir.go` 公共面冻结；红测先行。

---

## 1. 威胁模型

### 1.1 精确机理（GAP-77）

当前实现 `openBoundDirFile`（`bounddir_unix.go:64`）对净化后的相对名 `clean` 做**单次** `syscall.Openat(dir.fd, clean, O_RDONLY|O_CLOEXEC|O_NOFOLLOW, 0)`。Linux 语义：`O_NOFOLLOW` 仅约束解析中**最后一个组件**；`"link/secret"` 的中间组件 `link` 是符号链接时内核照常跟随，最终对象为界外 `O/secret`。`cleanBoundRelativePath` 只做拼写级守卫（空名/绝对/`..`），`NewBoundFile` 只拒目录/非正则/超界，均不复核对象身份 ⇒ 界外字节完整进入 `ReadAll`/`Hash`（P52 实跑读出 23 字节，证据 `05-attack/t51-1-linux-escape.md`）。

**不变式（本片要重建的）**：绑定读取返回的句柄，其对象必须落在绑定目录对象（由 `dir.fd` 钉住的 dentry 树）之内；任何一级解析越界 ⇒ 拒绝并关闭，绝不降级为路径重解析。

### 1.2 对抗载荷推演（方案 A=逐级 openat，方案 B=O_PATH + /proc/self/fd）

| 载荷 | 逐字形态 | 方案 A 行为 | 方案 B 行为 |
|---|---|---|---|
| 多重软链接 | `a→O1, a/b` 或 `a→b→…` | 中间组件 `a` 被 `O_NOFOLLOW|O_DIRECTORY` 打开 ⇒ ELOOP 拒绝 | `a` 被 `O_PATH|O_NOFOLLOW` 打开 ⇒ ELOOP 拒绝 |
| `..` 与软链接混合 | `sub/../link/secret`、`a/../../b` | `filepath.Clean` 已折叠 `..`（`sub/../link`→`link`，`a/../../b`→`../b` 被拼写守卫拒绝）；净化后不含 `..`，无法越界 | 同左（净化共用，两方案一致） |
| 硬链接 | `ln` 目录硬链接（需 CAP，Linux 拒 EPERM）；叶子硬链接 | 叶子硬链接仍指向界内文件，不越界 | 同左 |
| 并发 `rename` 替换中间目录（TOCTOU） | 每级 open 之间把 `sub` 换为指向 O 的软链接 | 免疫：下一级 `openat` 从**上一级已打开的目录 fd** 出发，解析对象钉死在原目录 inode；`sub` 名字被替换不影响已钉 fd | 免疫：O_PATH fd 钉住目录对象；`/proc/self/fd/N` 魔法链接解析到 fd 指向的 dentry，`rename` 不改变 fd↔对象绑定 |
| 符号链接环 | `a→a`、`a→b, b→a` | `O_NOFOLLOW` 对组件 `a` 直接返回 ELOOP，不跟随、不循环 | 同左；另有深度上限（见 §4）双保险 |
| `/proc` 自身被挂载点遮蔽 | 低权限者无权挂载；极简容器无 `/proc` 挂载 | 完全不依赖 `/proc` | 必须在 `OpenBoundDir` 时探测 `/proc/self/fd` 可用性；探测失败 ⇒ 降级 A（见 §5）。**探测用绝对路径 `/proc/self/fd`（O_PATH|O_DIRECTORY）**，内核从进程挂载命名空间根解析，不受 cwd/绑定目录影响；成功后才使用 B |

结论：两方案在全部载荷下均 **fail-closed**（拒绝 + 关闭句柄），核心差异只在实现载体与 `/proc` 依赖。A 无需 TOCTOU 补偿：fd 链（`dir.fd`→各级目录 fd→叶子）使解析全程锚定目录对象而非路径名，这是与 Windows 半部不同之处（Windows 因无 openat 才需要事后 `confirmOpenedArtifact` 复核；Unix 修复后复核仍不必要）。

---

## 2. 逐维对比

| 维度 | A：逐级 openat | B：O_PATH + /proc/self/fd |
|---|---|---|
| Linux 兼容性 | 仅需 `openat`+`O_NOFOLLOW`（内核 2.6.16+），**无 `/proc` 依赖**，chroot/极简容器可用；所有文件系统（ext4/xfs/btrfs/tmpfs/overlayfs/NFS）一致 | 需 `O_PATH`（内核 2.6.39+，现代发行版满足）**且 `/proc` 已挂载**；无 `/proc` 的 chroot/硬化容器不可用；NFS 有 ESTALE/权限边缘（失败即拒绝，不越界） |
| 性能 | 每级 1×openat + 1×close；典型深度 2 路径 = 3 syscall（相对现值 +2），可忽略 | 每级 1×openat(O_PATH) + 1×close，叶子经魔法链接 openat；同数量级 |
| 代码复杂度 | 低：新增 1 个游走辅助函数，`openBoundDirFile` 改为调用之；无探测、无分支 | 中：多探测函数、双分支、fd 数字符串拼接、fd 泄漏风险 |
| 安全性 | 每级 `O_NOFOLLOW` 天然拒绝非目录组件；O_DIRECTORY 让内核拒非常规组件；fail-closed | 同安全度；O_PATH 不校验类型，需自查 `S_IFDIR`（见 §4）；`/proc/self/fd` 探测失败降级 A 后安全度=A |
| 边界条件 | 深路径时游走期间每级仅 1 fd 存活，无 fd 压力 | 若不及时关闭中间 fd 会线性占用；设计为逐级关闭（O_PATH fd 钉住对象，关闭后已打开的下级 fd 独立存在，不影响解析） |

---

## 3. 推荐方案（定稿）

**B 为主 + A 为无 `/proc` 降级（混合）**。理由：

1. B 的 O_PATH 游走语义最贴合"逐组件解析、全程锚定对象"的目标，且叶子打开经 `/proc/self/fd/N` 与中间 O_PATH fd 天然对齐 P42 F1 的"同一次解析完成"诉求。
2. 红线要求 B 必须带降级路径，本设计以 **A 作为天然降级实现**：A 本就是本片的候选最小修复，代码量小、无环境依赖，探测失败时启用它不增加任何安全假设（A 自身即安全）。
3. 若评测结论倾向"最小改动"，A 单独也足以修复 GAP-77；但混合方案在无 `/proc` 环境与主路径上都有明确可论证的安全边界，故定为 A/B 混合。**实现以 B 为默认，A 为降级，二者共用"逐级 fstat 校验 + 叶子 O_NOFOLLOW 打开"骨架，分支点仅一处。**

---

## 4. 函数级设计（`bounddir_unix.go` 内，签名全部不变）

### 4.1 改动清单

| 函数/方法 | 变更 |
|---|---|
| `BoundDir` 结构体 | 新增私有字段 `procFdOK bool`（Unix 文件内私有，不影响公共契约；`bounddir.go` 无结构体定义，Windows 版独立） |
| `OpenBoundDir` | 末尾增加 `/proc/self/fd` 探测，写入 `procFdOK`；探测失败不报错，只置 false |
| `openBoundDirFile` | 单次 `openat` 替换为游走辅助函数调用；错误包装保持 `&os.PathError{Op:"openat", Path: clean, Err: err}` |
| 新增 `walkBoundComponents`（包内私有） | 逐级打开 + 叶子打开；深度上限 64 |
| 新增 `probeProcSelfFd`（包内私有） | 探测函数 |
| 常量 | `boundDirMaxDepth = 64` |

错误哨兵保持现有闭集：任何逃逸迹象（ELOOP/ENOTDIR/EMLINK 等 errno 或非目录 fstat）⇒ 关闭全部已开 fd 并返回错误，**不降级为宽松路径**。不新增哨兵（T51-1 已明确断言"无哨兵、看攻击结果"，不引入新错误类型）。

### 4.2 伪代码级流程

```go
// OpenBoundDir 末尾追加：
if d, err := syscall.Open("/proc/self/fd",
    syscall.O_PATH|syscall.O_DIRECTORY|syscall.O_CLOEXEC, 0); err == nil {
    _ = syscall.Close(d)
    dir.procFdOK = true
}
// 失败：dir.procFdOK 保持 false，绝不因探测失败拒绝绑定（非逃逸信号）。

// openBoundDirFile 主体替换：
clean, err := cleanBoundRelativePath(rel)
if err != nil { return nil, err }

// walkBoundComponents(dir, clean) → *os.File 或 error
//   1. comps := strings.Split(clean, "/")；若 len(comps) > 64 ⇒ 返回 ENAMETOOLONG 包装错误（fail-closed 深度帽）
//   2. cur := dir.fd（借用，不 close）
//   3. 逐级循环 i := 0; i < len(comps)-1; i++:
//        levelFd, err := syscall.Openat(cur, comps[i],
//            syscall.O_PATH|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
//        // 错误（ELOOP=软链接, ENOENT, ENOTDIR 等）⇒ 关闭链上已开 fd，返回 &os.PathError{Op:"openat", ...}
//        fstat(levelFd) 校验 S_IFMT == S_IFDIR，否则 ENOTDIR 拒绝（O_PATH 不校验类型，必须自证）
//        if cur != dir.fd { _ = syscall.Close(cur) }   // 逐级释放，防止深路径 fd 泄漏
//        cur = levelFd
//   4. 叶子打开：
//        if dir.procFdOK {
//            fd, err := syscall.Openat(cur, comps[last],
//                syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0)
//            // 注意：叶子直接对 cur（O_PATH fd）openat，无需 /proc/self/fd 魔法链接；
//            // 此步与旧实现唯一差异是 base fd 从 dir.fd 变为最后一级目录 fd
//        } else {
//            fd, err := syscall.Openat(cur, comps[last],
//                syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0)
//        }
//        // 实际上 A/B 叶子打开相同（都对 cur 用 openat）；B 的 /proc 价值只体现在
//        // "中间组件用 O_PATH 打开" 这一语义上——此处无分支差异，分支只在 cur 构造方式
//        若错误 ⇒ 关闭链上已开 fd，返回 PathError
//   5. if cur != dir.fd { _ = syscall.Close(cur) }
//   6. return os.NewFile(uintptr(fd), clean)
```

**实现要点**：

- **`/proc/self/fd` 的价值**：O_PATH fd 不触发文件系统实际打开（无权限检查、无引用计数开销），对 NFS/overlayfs 上不可读但可 `O_PATH` 的目录更宽容；A 路径对无读权限目录会因 `O_RDONLY` 被拒（fail-closed，可接受）。两者行为差异仅此，均不会放行越界。
- **叶子的 O_NOFOLLOW 语义不变**：叶子组件若是软链接仍被拒（覆盖 P52 T51-2 正向对照"界内 junction/leaf 放行"之外的叶子拒绝测试）。
- **深度帽 64**：正常证据路径（`requestId/state-events.jsonl` 等）深度 ≤3，64 为保守上限；超限返回 ENAMETOOLONG 包装错误（拒绝，非放行）。
- **错误包装**：游走与叶子错误统一 `&os.PathError{Op:"openat", Path: clean, Err: err}`，保留 `errors.Is(err, fs.ErrNotExist)`（`evidence_bind_dir_unix_test.go:94` 依赖此包装）；不改变任何既有哨兵。
- **`close` 语义**：`BoundDir.Close` 不变；游走失败时链上已开 fd 必须在返回前全部 close（防止泄漏）。

### 4.3 与 darwin/windows 半部的关系

darwin 半部（`bounddir_darwin.go`）不解析任何名字，不受影响；windows 半部（`bounddir_windows.go`）由构建标签隔离，本片零改动。`bounddir_clean.go`（`!darwin`）不改——其拼写守卫仍是第一道防线，逐级解析是第二道。

---

## 5. 降级策略（红线）

- **检测点**：`OpenBoundDir` 成功返回前，以**绝对路径** `syscall.Open("/proc/self/fd", O_PATH|O_DIRECTORY|O_CLOEXEC, 0)` 探测；成功 ⇒ `dir.procFdOK = true`（B 可用），失败 ⇒ `false`（A 降级）。
- **可靠性论证**：探测失败 = 内核进程挂载命名空间内不存在可用 `/proc/self/fd`，此时 B 的魔法链接必然失败 ⇒ 降级 A 是唯一可行路径；A 不依赖 `/proc`，任何 Linux 内核 ≥2.6.16 均可用。探测**不得**基于路径拼接（如 cwd 相对路径），避免被绑定目录内容干扰。
- **无运行时再探测**：绑定期内 `/proc` 挂载状态不会在无特权下改变；不做二次探测，消除探测与使用之间的新 TOCTOU 面。若探测成功但后续魔法链接失败（如容器异常），该失败作为普通 `openat` 错误返回 ⇒ 拒绝，仍 fail-closed。
- **A 路径自身安全**：无 `/proc` 环境下 A 的逐级 `openat(O_NOFOLLOW)` 同样阻断所有逃逸载荷（§1.2），降级不降低安全等级。

---

## 6. 红测先行设计（红线）

### 6.1 预期红形态

T51-1（`TestBoundDirMiddleSymlinkEscapesBound`）移除 `t.Skip` 后，在修复前的 HEAD 上运行：

- `OpenBoundDirFile(dir, "link/secret", 1024)` 返回 nil 错误 ⇒ 测试以 `t.Fatalf("GAP-77 regression: middle symlink escape succeeded and read outside the bound directory")` 失败；等价证据：手工探针读出 `SECRET-OUTSIDE-BOUNDARY`（23 字节，P52 已实证 `ATTACK-SUCCEEDED`）。
- 红的**判据** = 测试 FAIL 且失败文案为 GAP-77 regression（攻击得手），不是其它环境错误（如 EACCES）。

### 6.2 执行环境与步骤（红测）

- 环境：远程 Ubuntu VM `10.0.0.199`，`go1.22.12 linux/amd64`（已就位 `~/go1.22`），仓库克隆 `~/LZProjects/prfrail`；与 CI `GO_VERSION=1.22.12` 对齐。SSH 需 `bash -lc` 加载登录环境（CHECKLIST v42 铁律 19）。
- 步骤：
  1. 同步仓库至本片基线 HEAD，仅删除 `bounddir_escape_unix_test.go:39` 的 `t.Skip` 行（**不碰生产代码**）。
  2. `bash -lc '~/go1.22/bin/go test -run TestBoundDirMiddleSymlinkEscapesBound -v -count=1 ./internal/adapters/'` ⇒ 预期 **FAIL**，文案含 `GAP-77 regression`。
  3. 归档红测输出（原文 + 退出码）作为红证据。
  4. 落地 §4 生产代码后重跑同命令 ⇒ 预期 **PASS**；并复跑 P52 探针（`05-attack` 载荷）确认 `ATTACK-BLOCKED` 且无界外字节。
- 绿判定：T51-1 PASS；`go test -count=1 ./...` 全绿（十五包全 `ok`）；`gofmt -l` 无输出；`go build ./...`、`go vet ./...` 通过；`GOOS=linux go vet ./internal/adapters/` 对向复核通过（CHECKLIST 门禁口径）。

### 6.3 回归用例清单（不可误杀）

| 用例 | 落点/来源 | 必须仍通过 |
|---|---|---|
| 界内多层目录正常读取（深度 ≥2，如 `sub/x.txt`） | 本片新增到 `bounddir_escape_unix_test.go` 或独立 unix 测试（正向对照） | 读取内容与直读一致 |
| 单层文件读取 `f.txt` | `evidence_bind_dir_unix_test.go:70` | PASS 不变 |
| 不存在文件 ⇒ `fs.ErrNotExist` | `evidence_bind_dir_unix_test.go:94` | `errors.Is` 匹配不变 |
| 叶子软链接拒绝（`O_NOFOLLOW` 语义） | `evidence_bind_dir_unix_test.go:99` | 拒绝不变 |
| 目录替换后 `boundDirVerify` 拒绝/原名接受 | `bounddir_identity_test.go` | PASS 不变 |
| 拼写守卫四形态拒绝 | `bounddir_identity_test.go` | PASS 不变 |
| replay store 经 `openBoundDirFile`（深度 2，`requestId/state-events.jsonl`）读取 | `agent_runner_replay_store.go:261` 相关测试 | 无回归 |
| 降级 A 分支（`procFdOK=false`）行为 | 本片新增 unix 测试：直接调用包内辅助函数或注入探测结果 | 界内读取成功、逃逸拒绝 |

---

## 7. Windows 对向非影响声明

- 本片只改 `//go:build !windows && !darwin` 的 `bounddir_unix.go` 及其 unix 测试；`bounddir_windows.go`、`bounddir_escape_windows_test.go`、`bounddir_clean.go`（`!darwin`）零改动。
- 无需对向 GOOS 复核（Windows 半部文件未动）；但按 CHECKLIST 门禁口径保留 `GOOS=linux go vet ./internal/adapters/` 与全量 `go test ./...` 作为兜底。darwin 半部为 fail-closed stub，不受本片影响，无需 GOOS=darwin 复测（可顺手 `GOOS=darwin go build ./internal/adapters/` 确认构建不破）。

---

## 8. 开放问题清单

1. **交付物范围**：CHECKLIST v42 §四 P53 交付物要求本片同步更新 `docs/SECURITY.md`（SEC-02 威胁矩阵补充中间软链接逃逸载荷）与 `docs/TEST_STRATEGY.md`（AT-02 强制反例）。用户本片硬约束为"只改 `bounddir_unix.go`（及其测试）"，与切片交付物存在交叉。**建议**：代码修复归本片；两处规范文档更新作为同一片的伴随交付物，由 ② 实施完成，属"契约修订先改文档"纪律。请 Sol 裁定是否把文档更新纳入本片派单（本设计不阻止，但需明确归属）。
2. **降级分支测试注入**：A 降级路径需要可测；建议用包内测试直接以 `procFdOK=false` 的 `BoundDir` 实例驱动辅助函数，避免测试钩子污染生产代码。若 Sol 认为需公共可测注入，需另行裁定（默认按包内测试处理，不扩公共面）。
3. **叶子打开是否复用 O_PATH 魔法链接**：本设计对叶子仍用 `openat(O_NOFOLLOW)`（语义等价旧实现、无 `/proc` 依赖），仅中间组件用 O_PATH。若未来要求叶子也走 `/proc/self/fd`（如对抗 exotic 文件系统），属后续加固，不阻塞本片。
4. **深度帽数值**：64 为建议值，② 实施可结合最大合法证据路径深度（≤3）复核；原则是"远超合法深度 + 深于任何现实路径"，避免误杀。
5. **红测证据归档路径**：红测输出建议归档至 `docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/red-baseline.md`（按 P52 目录范式），请 ③ 测试先行确认归档命名。

---

## 9. 结论（一句话派单）

在 `bounddir_unix.go` 内将单次 `openat(O_NOFOLLOW)` 替换为**逐级组件解析**：中间组件以 `O_PATH|O_NOFOLLOW` 逐级打开并 fstat 校验目录（主路径 B），`OpenBoundDir` 时探测 `/proc/self/fd` 不可用则降级为**逐级 `openat(O_NOFOLLOW|O_DIRECTORY)`**（降级 A），叶子统一 `openat(O_NOFOLLOW)`，全程错误即关闭、拒绝、fail-closed；公共签名与哨兵闭集不变，T51-1 移除 skip 红→绿，Windows/darwin 半部零改动。

---

**硬约束遵守声明**：本设计不改 `bounddir.go` 公共面；红测先行（§6）；B 含无 `/proc` 降级（§5）；改动范围仅 `bounddir_unix.go` 及其 unix 测试（§4.1/§7，文档更新归属见 §8-1 待 Sol 裁定）。
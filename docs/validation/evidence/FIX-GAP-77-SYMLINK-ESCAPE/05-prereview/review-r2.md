# ⑤ 预审 r2 复核报告 · P53 FIX-GAP-77-SYMLINK-ESCAPE（原位重跑，只读，逐字归档）

## 整改闭环核对

**M1（风暴测试平台门 + SKIP 位置）— 成立 ✓**
- `bounddir_symlink_storm_test.go` 标签实测为 `linux && amd64`；F1 首行即 `gap77SkipUnlessOpenat2(t)`。
- F2 关键点验证通过：子进程分支（`bindMountHelperEnv=="1"`）**不**调用 skip 门，直接跑 helper；`gap77SkipUnlessOpenat2(t)` 只放在父进程分支且在子分支判定**之后**——子进程 SKIP 退出码 0 被父进程读成通过的风险（M1 原指）已消除。父进程通过探测 ⇒ 同主机子进程必然具备 openat2，逻辑自洽。
- 可证伪性：arm64 上生产恒 ENOSYS，"至少一次成功读/EXDEV" 必然假红，收窄 + 探测双保险均就位。

**M2（fail-closed 契约 + 内核探测）— 成立 ✓，可证伪性成立（含一处固有残余）**
- `TestBoundDirFailClosedWithoutOpenat2`：锁存 `openat2Disabled=true`（cleanup 复位）、用合法界内名 `original` 调用（拒绝必须来自门而非名字）、断言 `*os.PathError`/`Op=="openat"`/Path 为清洗名/`errors.Is(ENOSYS)`/无 panic/无句柄/二次同形拒绝。非 linux/amd64 分支清锁后以 `dirfd=-1` 调 `openBoundArtifactLocked`，断言常量平台门直接返回 ENOSYS——与生产代码首检 `runtime.GOOS/GOARCH` 完全对应。
- **可证伪性实测推理**：删除 `openat2Disabled.Load()` 门 → linux/amd64 上 `openat2Resolve` 直通 → 界内 open 成功 → `gotHandle` → Fatal。变异可拦。
- 残余（Low，无需改代码）：linux/amd64 上单独删除**平台门**不可被本测试证伪（该门对 linux/amd64 行为本身无影响），仅由本测试在非 linux/amd64 主机运行时覆盖。属平台门控固有属性。
- `gap77ProbeOpenat2`：直连 `openat2Resolve`，完全绕过 `openBoundArtifactLocked`，与锁存互不干扰（源码注释与 grep 双确认）；`sync.Once` 在首个调用者测试内同步执行，临时目录生命周期安全；无文件创建。

**L1（锁序口径）— 口径维持，留 ⑧ ✓**
- 生产 `bounddir_unix.go` SHA256 与 r1 记录逐位一致（见下）⇒ 锁序事实自 r1 起零变化。抽查确认生产路径均非嵌套持有：`readBoundEvidence`/`CloseBoundDirs` 先释 store.mu 再触 dir.mu，`bindEvidenceDir` 先 dir.mu 后 store.mu（顺序执行）。⑤ r1 的「设计陈述与实现出入」属文档写回事宜，⑧ 按实况描述锁序的处置口径继续成立。

**L2（F4 并发 Close×读取）— 成立 ✓**
- `TestBoundDirConcurrentCloseAndRead`：4 读者 × 16 次 Close；只允许「界内全串精确 canary」或 `adapters: nil bound directory`，EBADF/截断/其它字节/读失败/Close 失败全部登记为 problem；`ready` 通道 + 10s 超时防真空；终态断言目录已失效且文案一致。
- 说明一处表述差异：canary 实际为 16 字节（`IN-SIDE-BOUNDARY`，经核验）；③c 报告中的 "27 byte-exact reads" 指 VM 实测**成功读次数**而非 canary 长度，任务包中的「27 字节」系转录误读，不影响断言正确性（全串比较即字节精确）。

**L3（红基线变体改造）— 成立 ✓**
- 标签 `!windows && !darwin`、独立名 `TestBoundDirRedBaselineMiddleSymlinkEscape`（与常驻回归不重名不遮蔽）、env 门 `PRFRAIL_P53_RED_BASELINE=1` 在前、linux/amd64 SKIP 在后并给出显式判读理由。默认套件 SKIP 不 FAIL。
- 代码内 `p53_red_baseline` 仅存于注释（退役说明）与 ③b 历史归档（green-suite-report.md/MANIFEST.md，③c 已声明保留历史不改写），无任何生效 build-tag 引用。
- 边界推理通过：<5.6 内核 fail-closed 主机上，env 武装后未变异 → "red baseline not reproduced"（守卫有效，判读正确）；变异落地 → open 成功 → 回归文案。非 linux/amd64 SKIP 正确（生产常量门 ENOSYS，变异目标定义在 linux/amd64）。

## 新增问题排查

1. **锁存残留**：全包 grep `t.Parallel()` 为空 ⇒ 测试串行，锁存窗口不泄漏；cleanup 兜底复位；触碰 `openat2Disabled` 的文件仅生产代码 + failclosed 测试（`escape_unix_test.go` 中仅注释）。无残留。
2. **探测缓存**：进程内一次缓存，与锁存正交，无副作用。无风险。
3. **红基线环境门与默认套件**：默认 SKIP 不影响绿集；残余（Low）为环境级误设 `PRFRAIL_P53_RED_BASELINE=1` 会让默认套件红——env 名高度专属，可接受。
4. **M2 残余**：见上（平台门删除在 linux/amd64 上不可单独证伪，其它平台分支覆盖）。Low，固有属性。
5. 未发现整改引入的编译/竞态/语义新问题。

## 门禁读数复核

- **生产零改动（硬项）**：`Get-FileHash` = `88FAB4356A113D426DD396858C9A08CF4190E3D9CD2DF1243706E4170D0C5C44`，与任务给定值逐位一致（大小写不敏感比较相同）。git 状态 6 项中 `bounddir_unix.go` 的 M 态属 ②a 实现轮，与 ③c 声明一致。
- **本机抽查**：`go build ./...` EXIT=0；`go vet ./...` EXIT=0；`go test -count=1 ./internal/adapters/` ok 36.1s；`GOOS=linux GOARCH=amd64/arm64 go vet` 均 EXIT=0（红基线/风暴/failclosed 文件在两视图可编译）。
- **VM 侧**（③c 报告 + 主控回执，与本机代码态互证）：go1.22.12 全包 ok、定向 -v PASS×5/SKIP×2（默认红基线 SKIP、F2 userns 受限 SKIP）、-race 无 DATA RACE、`PRFRAIL_P53_RED_BASELINE=1` → FAIL/EXIT=1 文案 `red baseline not reproduced`。默认 SKIP_EXIT=0 与 env 版 ENV_EXIT=1 均与代码逻辑吻合。

## 结论

M1/M2/L2/L3 四项整改全部闭合且可证伪性成立（M2 平台门证伪能力属平台固有残余）；L1 锁序口径维持、⑧ 写回处置有效；未发现新增 Medium+ 问题，仅两处 Low 残余（记录如上，不需改代码）。

PRE-REVIEW: PASS

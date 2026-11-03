# ② 实施角色报告合集（逐字归档）

载体：runSubagent(agentName="prfrail-implementer", model="DeepSeek V4.1 Flash (deepseek)")，四个批次（②a–②d）严格串行。逐字归档。

---

IMPL-S1: DONE

## 交付文件 A：`internal/evidence/boundread_counterexample_test.go`（package evidence，无 build tag）

新增测试函数：
- `TestBoundFileReadAfterCloseRefused` — `if data, truncated, err := bound.ReadAll(); !errors.Is(err, ErrBoundFileClosed)` / `if digest, err := bound.Hash("probe"); !errors.Is(err, ErrBoundFileClosed)`
- `TestBoundFileHashRefusesGrowthPastBound` — 断言 `truncated==true && len(data)==4`，随后 `if digest, err := bound.Hash("probe"); !errors.Is(err, ErrBoundArtifactExceedsBound)`
- `TestBoundFileNewRefusesNilAndNegativeBound` — `!errors.Is(err, ErrBoundFileNotRegular)`（nil 句柄）/ `!errors.Is(err, ErrBoundArtifactExceedsBound)`（bound=-1）
- `TestBoundFileOpenRefusesDirectoryWithSentinel` — 用 `errors.Is` 双形态哨兵（见偏离 D1）
- `TestBoundFileOpenRefusesOverBoundWithSentinel` — `!errors.Is(err, ErrBoundArtifactExceedsBound)`（9 字节 @ bound 8）

## 交付文件 B：`internal/adapters/bounddir_identity_test.go`（package adapters，`//go:build !darwin`）

新增测试函数：
- `TestBoundDirVerifyRefusesSwappedDirectory` — `swapped := boundDirVerify(dir, boundPath)`；`swapped == nil` ⇒ Fatal；`!strings.Contains(swapped.Error(), "is not")` ⇒ Fatal；控制组 `boundDirVerify(dir, boundPath+".bak") == nil`（见偏离 D4）
- `TestCleanBoundRelativePathRefusesEscapeShapes` — 四形态 `err == nil` ⇒ Fatal；对照 `clean != filepath.FromSlash("sub/x")` ⇒ Fatal（见偏离 D2）

签名已核实：`boundDirVerify(dir *BoundDir, expected string) error`（unix L90 / windows L168）、`cleanBoundRelativePath(rel string) (string, error)`（`bounddir_clean.go` L20，同为 `!darwin`）。

## 自跑门禁

```
--- gofmt -l ---            (无输出)
--- go vet ./internal/evidence/ ./internal/adapters/ ---   (无输出)
--- go test -count=1 ./internal/evidence/ ./internal/adapters/ ---
ok      github.com/larsonzh/prfrail/internal/evidence   0.883s
ok      github.com/larsonzh/prfrail/internal/adapters   37.853s
```
七条新测试逐条 `--- PASS`（verbose 实跑，无 SKIP）。

## git status --short（本切片相关行）

```
?? internal/adapters/bounddir_identity_test.go
?? internal/evidence/boundread_counterexample_test.go
```
其余 `M/A` 条目（`docs/t027/*`、`docs/validation/**`）为其他角色既有改动，本环未触碰；无生产代码与既有测试文件被修改。

## 偏离（逐条）

- **D1（断言形态，必须登记）**：测试 4 在 Windows 上无法拿到 `ErrBoundFileNotRegular`。实测 `OpenBoundFile(<目录>, 64)` 在 Windows 由 opener 直接失败：`*fs.PathError ... Access is denied`（`errors.Is(err, fs.ErrPermission)==true`，`ErrBoundFileNotRegular==false`），因为 `boundopen_windows.go` 的 `CreateFile` 带 `FILE_ATTRIBUTE_NORMAL` 且无 `BACKUP_SEMANTICS`，句柄根本没打开，`NewBoundFile` 的类型检查无机会执行。断言仍为 `errors.Is`（非 `err != nil`），采用显式两态：`ErrBoundFileNotRegular`（unix 形态）或 `fs.ErrPermission`（Windows 形态），其它错误一律 Fatal；理由与实测已写入测试注释。
- **D2（路径期望值）**：对照组期望值用 `filepath.FromSlash("sub/x")`。`filepath.Clean(filepath.FromSlash(...))` 在 Windows 返回 `sub\x`，字面 `"sub/x"` 会误报；语义与「折叠 `.`、归一化分隔符」一致。
- **D3（编码/行尾）**：两个新文件首次落盘为 CRLF，`gofmt -l` 报了这两条。已按仓库约定改写为 **UTF-8 无 BOM + LF**（实测 `crlf=0 bom=false`），重跑 `gofmt -l` 为空。
- **D4（超范围增强）**：测试 6 增加一条控制组断言（同一 `boundDirVerify` 对 `P+".bak"` 必须返回 nil），用于证明拒绝来自「对象被替换」而非「校验本身报错」。不改变被测语义，本机通过。
- 测试 2 未 skip：第二写者打开成功，实测 `readall len=4 trunc=true`、`hash IsExceeds=true`。

IMPLEMENT: DONE
IMPL-S1: END

---

IMPL-S2: DONE

## 交付文件 C：`internal/adapters/bounddir_escape_unix_test.go`（`//go:build !windows && !darwin`）
测试函数：`TestBoundDirMiddleSymlinkEscapesBound`
- L39 skip 文案（逐字）：`t.Skip("Known vulnerability GAP-77: tracked for P53")`（位于被测调用之前，因当前 HEAD 该调用**成功返回界外句柄**）
- 载荷：`D`/`O`、`D/original`="inside"、`O/secret`="SECRET-OUTSIDE-BOUNDARY"、`D/sub`、`os.Symlink(O, D/link)`
- 关键断言（L61-66）：`_, err = OpenBoundDirFile(dir, filepath.Join("link","secret"), 1024)`；`err == nil` ⇒ `t.Fatalf("GAP-77 regression: middle symlink escape succeeded and read outside the bound directory")`。无哨兵（P53 才引入 openat 形态错误），断言落在**攻击结果**上，理由已写入文件内注释。
- 文件头注释含 GAP-77/V-1、证据路径 `05-attack/t51-1-linux-escape.md`、P53 修复方式（逐级 openat）与"删除 skip 行后应 PASS"。

## 交付文件 D：`internal/adapters/bounddir_escape_windows_test.go`（`//go:build windows`）
签名已核对：`OpenBoundDir(string) (*BoundDir, error)`、`OpenBoundDirFile(*BoundDir, string, int64) (*evidence.BoundFile, error)`、`(*BoundDir).Close() error`、`(*BoundFile).ReadAll/Close`。

1. `TestBoundDirRefusesJunctionEscape`
   - `cmd /c mklink /J D\j O`（绝对路径）；失败 ⇒ `t.Skipf("junction creation unavailable: %v (%s)", ...)`
   - L69-76：`err == nil` ⇒ 关闭句柄后 `t.Fatalf("a junction leaving the bound directory must be refused, got a handle")`；非空则断言 `strings.Contains(err.Error(), "outside the bound directory")`（`bounddir_windows.go` L128 `confirmOpenedArtifact` 文案）。
2. `TestBoundDirJunctionInsideBoundAdmitted`
   - junction → `D/sub`；`OpenBoundDirFile(dir, "j/f", 1024)` 必须成功；`bytes.Equal(data, []byte("inside"))` 且 `truncated == false`（正向对照，证明拒绝不是"凡 junction 皆拒"）。
- 两条均 `defer os.Remove(link)` + `defer/explicit dir.Close()`，文件句柄亦显式 `Close`。

## 自跑门禁（本机 Windows，cwd=`d:\LZProjects\prfrail`）
```
gofmt -l internal\adapters\   → 无输出
go vet ./internal/adapters/   → VET_EXIT=0
go test -count=1 ./internal/adapters/
ok  github.com/larsonzh/prfrail/internal/adapters  38.265s   (PKG_TEST_EXIT=0)
本批逐条：
--- PASS: TestBoundDirRefusesJunctionEscape (0.16s)
--- PASS: TestBoundDirJunctionInsideBoundAdmitted (0.06s)
PASS / ok ... 1.355s
TestBoundDirMiddleSymlinkEscapesBound 在 Windows 被 build tag 排除（仅 linux 生效，其 skip 保证 CI 不红）
```

## 交叉编译自检
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0
go vet ./internal/adapters/            → LINUX_VET_EXIT=0
go test -c -o %TEMP%\p52s2_linux.test  → LINUX_BUILDTEST_EXIT=0
```
（未在 Linux 实跑，本机无 WSL/容器通道；C 文件的 import 与 API 已在 linux 下完成类型检查与链接级编译。）

## 偏离清单
1. 文件创建工具写出的是 **CRLF**，不符合 `.go` 的 LF 约定；已用 `gofmt -w`（仅作用于这两个新文件）规范化为 LF + 标准格式，常量块对齐随之调整。
2. D-1 在 `err == nil` 分支先 `file.Close()` 再 `Fatal`（防句柄泄漏），原单未要求。
3. D-1 增加 `defer os.Remove(link)` 显式摘除 reparse point，避免清理阶段任何跟随 junction 的风险。
4. C 文件的 `err != nil` 属用户指定的"攻击结果"断言形态；已在注释中写明无哨兵原因，未追加 `errors.As(*os.PathError)`（避免在 P53 实现形态未定前写死错误类型）。
5. 未 commit / 未 push；`git status` 确认本次仅新增上述两个未跟踪文件（`git diff --stat` 为空，无任何生产代码或既有测试被修改）。

IMPL-S2: END

---

IMPL-S3: DONE

## 文件 E：`internal/gates/probe_permit_counterexample_test.go`（新建，package gates，无 build tag）

测试函数清单：
1. `TestProbePermitRequireNilPermitRejected`（L29）
2. `TestProbePermitRequireNilContextRejected`（L41）
3. `TestProbePermitCompleteUnknownStateIsConsumed`（L65）
4. `TestProbePermitLoadNonNotExistReadFailure`（L90）
5. `TestProbePermitLoadEmptyPathRejected`（L113）
6. `TestProbePermitCrossProcessSingleClaim`（L130）+ helper `TestProbePermitCrossProcessHelper`（L154）

关键断言摘录：
- L33：`if !errors.Is(err, ErrProbePermitMissing)`（nil receiver）
- L46：`if !errors.Is(err, ErrProbePermitMissing)`（`Require(nil)`，L215 分支）；L51 追加 `os.Stat` 不得落 marker（前置检查无副作用）
- L75/L77：`permit.state.Store(uint32(3))`（合法常量 0/1/2 之外）→ `!errors.Is(err, ErrProbePermitConsumed)`；注释说明越界态必须 fail-closed
- L102：三连 `!` 全枚举——`errors.Is(err, ErrProbePermitMissing) || errors.Is(err, ErrProbePermitInUse) || errors.Is(err, ErrProbePermitConsumed)` 必须为假；L105 追加 `strings.Contains(err.Error(), "load probe permit")`；注释写明契约依据（读失败原样返回、无哨兵）
- L115：`LoadProbePermit("")` → `!errors.Is(err, ErrProbePermitMissing)`
- L146：父进程 `LoadProbePermit(marker)` 后 `Require` → `!errors.Is(err, ErrProbePermitConsumed)`；子进程 `exec.Command(os.Args[0], "-test.run=TestProbePermitCrossProcessHelper")`，env `P52_MARKER_PATH`，成功 `os.Exit(0)`；非零退出时 `t.Fatalf` 附 `CombinedOutput`

## 文件 F：`internal/gates/readiness_counterexample_test.go`（新建，package gates，无 build tag）

测试函数清单：
1. `TestReadinessProductionRejectsMalformedNonEmptyAT23`（L28，4 形态子测试）
2. `TestReadinessProductionValidatesSealBeforeAT23`（L84）
3. `TestReadinessAcceptanceRejectsPaddedSealHash`（L109）
4. `TestReadinessAcceptanceAdmitsTrimmedAT23`（L134）

关键断言摘录：
- L64：四形态（`sha256:`+63hex / `sha256:`+65hex / `sha256:`+63hex+`g` / `SHA256:`+64hex）均 `!errors.Is(err, ErrReadinessEvidenceIncomplete)` 且 `ready != nil` 必须为假
- L90：`SealBinding{Attempt:0}`（非法）+ 合法 AT23 → `ErrReadinessEvidenceIncomplete` 且 `ready == nil`；L98 追加 `strings.Contains(err.Error(), "seal attempt")`，机械钉住「seal 先于 AT23」
- L113：`SealHash: 64hex+" "` → `ErrReadinessEvidenceIncomplete`（SealHash 不 trim）
- F4：`AT23Hash: "  "+64hex+"  "` → `err == nil`；`Release().AT23Hash == bare`（TrimSpace 行为钉，注释标注 GAP-79 候选）

## 自跑门禁摘录

```
gofmt -l <两文件>   →  无输出（GOFMT-CLEAN）
go vet ./internal/gates/  →  VET-OK（无输出）
go test -count=1 ./internal/gates/  →  ok  github.com/larsonzh/prfrail/internal/gates  1.625s
```

本批 10 条逐条 `--- PASS`（含子测试 4 条）：
`TestProbePermitRequireNilPermitRejected` / `...RequireNilContextRejected` / `...CompleteUnknownStateIsConsumed` / `...LoadNonNotExistReadFailure` / `...LoadEmptyPathRejected` / `...CrossProcessSingleClaim`(0.17s) / `TestReadinessProductionRejectsMalformedNonEmptyAT23`(+4 子测试) / `...ValidatesSealBeforeAT23` / `...RejectsPaddedSealHash` / `...AdmitsTrimmedAT23`，均 PASS；`PASS` 后 `ok`。

全量包运行唯一 SKIP 为 `TestProbePermitCrossProcessHelper`（设计如此，无 env 时 skip）。

## 偏离说明（逐条）

1. **文件编码**：初次写入为 CRLF，`gofmt -l` 报两文件；已用 `gofmt -w` 规范化为 LF（`hasCR=false`、`bom=false`），符合仓库 `.go` = UTF-8 无 BOM + LF。
2. **子进程 stderr**：要求「附子进程 stderr」，实现用 `cmd.CombinedOutput()`（`go test` 的 `t.Fatalf` 走 stdout，仅取 stderr 会丢诊断）。注释已标明。
3. **F2 追加断言**：除哨兵断言外增加 `strings.Contains(err.Error(), "seal attempt")`，用于真正钉住「seal 先于 AT23」的时序（仅哨兵断言不足）。不改变契约。
4. **未修改**：生产代码、既有测试文件均零改动（`git status` 中新文件为 `??` 未跟踪；其余 `A/M` 条目属本片前序批次，非本环产出）。
5. 未执行任何 `git add` / `git commit` / `git push`。
6. **（⑤ 预审 L2 补登记）**：T64-6 实现相对 ① 设计简化——设计规格含 claimed/release 信号握手与「父新实例再 Require」两段，实现改为顺序式（子退出→父 Load+Require）；覆盖目标（跨进程持久声明）仍达成，属设计-实现偏差登记缺口，已按 ⑤ 预审要求补登记。

IMPL-S3: END

---

IMPL-S4: DONE

## 测试函数清单与关键断言

新增唯一文件：`internal/chain/state_seal_counterexample_test.go`（`package chain`，无 build tag，约 300 行，含自足最小装配）。

| 测试 | 断言 |
|---|---|
| `TestSealValidateRefusesNonHashContextKey` | `ContextKey="not-a-hash"` → `errors.Is(err, ErrSealIncomplete)`（`seal.go:81-83`） |
| `TestSealValidateAdmitsEmptyContextKey` | `ContextKey=""` → `validate() == nil`；注释钉明这是 **GAP-78 候选的现状记录**（validate 契约允许空 K，两消费端各自 fail-closed） |
| `TestAcceptanceReplayRefusesForeignContextKey` | 两个子测试（见下） |

第 3 条两个子测试：
- `foreign K is refused`：先用真实 `FileEventStore.SealTaskEvidence` 落盘一份 `ContextKey=cxForeignContextKey`（合法 `sha256:` digest）的密封记录，再 `Engine.Run`。断言：`err != nil`、`errors.Is(err, ErrPaused)`、`TaskStates[taskKey("task-one",1)] == "REPAIR_PENDING"`、事件日志中该 task 的 `REPAIR_PENDING` 事件 `Reason.Code == "evidence-seal-mismatch"`。
- `bound K replays`（对照）：同构记录但 `ContextKey = TaskContextKey("run-one","task-one",1,cxParentHash,"task-one")` → `Run` 返回 nil、task 状态 `PASSED`、chain `COMPLETED`。证明拒绝来自 K 差异而非构造错误。

## 等效性依据

- 走的不是 mock 或手搓 state：`Engine.Run` → `runTask`（`engine.go:196-229`，steps/postflight/REVIEW_PENDING 全真实）→ `acceptFromExistingSeal`（**调用点 `engine.go:230`**，K 比较 **`engine.go:352-354`**）。
- `acceptFromExistingSeal` 虽非导出，但测试文件就在 `package chain` 内，实际触碰的是生产入口 `Run`，非直接调用非导出方法，故无需 skip。
- 文件自足：未复用 `state_seal_test.go` / `engine_test.go` 的辅助，自写 9 个最小端口实现与 `cxOptions`；`FileEventStore` 为生产实现（密封写读均真实 fsync 路径）。

## 自跑门禁摘录

```
gofmt -l internal/chain/state_seal_counterexample_test.go   # 无输出
go vet ./internal/chain/                                    # 无输出
=== RUN   TestSealValidateRefusesNonHashContextKey
--- PASS: TestSealValidateRefusesNonHashContextKey (0.00s)
=== RUN   TestSealValidateAdmitsEmptyContextKey
--- PASS: TestSealValidateAdmitsEmptyContextKey (0.00s)
=== RUN   TestAcceptanceReplayRefusesForeignContextKey
    --- PASS: .../foreign_K_is_refused (0.03s)
    --- PASS: .../bound_K_replays (0.02s)
--- PASS: TestAcceptanceReplayRefusesForeignContextKey (0.06s)
ok  github.com/larsonzh/prfrail/internal/chain  1.380s   # 全包测试
```

## 偏离清单

1. **无 skip**：第 3 条已用 `Run` 真实路径落地，未触发预案中的 `t.Skip`。
2. **一处措辞级说明（非偏离）**：第 2 条按任务要求断言 `nil`，这是"现状钉"，与 GAP-78 修复方向相悖是预期内的；注释已写明修复 GAP-78 时需同步改此断言。
3. **文件规模**：为满足"独立自足"硬约束，装配代码比复用既有辅助多约 100 行；未引入任何外部依赖。
4. **范围校验**：`git status` 显示本环只新增该 1 个文件，生产代码与既有测试文件零改动（`internal/adapters/*`、`internal/gates/*` 的未跟踪文件为前序批次遗留，本环未触碰）。
5. 未执行任何 `git add/commit/push`（无同轮授权）。

IMPL-S4: END

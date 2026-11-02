# DR13-FIX-D3 · ③ 测试红基线（red baseline）

- 切片：P48 `DR13-FIX-D3`（starting marker release 顺序反转）
- 输入方案：`docs/validation/evidence/DR13-FIX-D3/01-arch/design.md`；辅助扫描 `09-glm/glm-1-case-scan.md`、`09-glm/glm-2-impact-scan.md`
- 被测现状（生产代码未修复）：`internal/adapters/agent_runner_starting_marker.go:92-101`，`release()` 现序 = `unlockStartingMarker`(:96) → `file.Close`(:97) → `os.Remove`(:99)
- 结论：**确定性顺序红测 FAIL**（主证据）；压力回归亦 FAIL（补充证据，R3 概率口径）

## 1. 新增测试文件清单

| file:line | 测试 | 类型 | 对应方案条目 |
|---|---|---|---|
| `internal/adapters/agent_runner_starting_marker_order_test.go:154` | `TestStartingMarkerReleaseRemovesBeforeUnlock` | 结构性（go/ast 源码顺序） | §3 结构性红测（首选） |
| `internal/adapters/agent_runner_starting_marker_order_test.go:235` | `TestStartingMarkerReleaseRemovesThePublishedPath` | 语义直测（remove-while-open） | §3 语义直测 |
| `internal/adapters/agent_runner_starting_marker_order_test.go:261` | `TestStartingMarkerReleaseIsIdempotent` | 幂等直测 | §3 幂等 |
| `internal/adapters/agent_runner_starting_marker_order_test.go:296` | `TestStartingMarkerReleaseNeverExposesReleasedContent` | 压力回归（概率型） | §3 压力（R3 口径） |
| `internal/adapters/agent_runner_starting_marker_order_windows_test.go:55` | `TestStartingMarkerReleaseFailClosedWhenTheNameIsOccupied` | Windows fail-closed 反测 | §3 失败方向 |
| `internal/adapters/agent_runner_starting_marker_order_windows_test.go:118` | `TestStartingMarkerReleaseShareFactProbe` | 共享事实探针（只记录） | §3 共享事实（GLM-2 R-N1） |

测试辅助（同文件内，均为测试态私有，未触碰生产符号）：`dr13SourcePath`:36、`dr13ReleaseFuncBody`:50、
`dr13IsClaimPathExpr`:79、`dr13ClaimPathAliases`:91、`dr13IsRemoveOfClaimPath`:111、`dr13StatementListing`:133、
`dr13OccupyStartingMarkerName`（Windows）:29。

## 2. 红/绿结论（生产代码未修复状态下）

| 测试 | 结果 | 说明 |
|---|---|---|
| `TestStartingMarkerReleaseRemovesBeforeUnlock` | **FAIL（红）** | 确定性红测：`unlockStartingMarker(claim.file)`(line 96) 先于 `os.Remove(claim.path)`(line 99) |
| `TestStartingMarkerReleaseNeverExposesReleasedContent` | **FAIL（红）** | 补充证据：4/35 次「扫描者持锁观察」命中的记录属于「随后删名成功」的 claim（144 次发布、62 次残余删除） |
| `TestStartingMarkerReleaseRemovesThePublishedPath` | PASS | 现序先 Close 后 Remove，故 remove-while-open 语义测试在旧代码上为绿；修复后仍须绿（钉住 claim 句柄 FILE_SHARE_DELETE 前提） |
| `TestStartingMarkerReleaseIsIdempotent` | PASS | 旧代码二次 release 因 `claim.file = nil` 返回 false；修复后仍须绿 |
| `TestStartingMarkerReleaseFailClosedWhenTheNameIsOccupied` | PASS | 失败方向语义在新旧序下同结论；修复后仍须绿（并证明 remove 失败时锁确已释放） |
| `TestStartingMarkerReleaseShareFactProbe` | PASS（仅日志） | 本工具链实测：`os.OpenFile` 句柄阻止 `os.Remove`，errno=32(0x20) sharing violation |

## 3. 红测失败原文（关键摘录）

```text
=== RUN   TestStartingMarkerReleaseRemovesBeforeUnlock
    agent_runner_starting_marker_order_test.go:225: release() drops the host lock before the published name is gone: "unlockStartingMarker(claim.file)" (line 96) runs before "return os.Remove(claim.path) == nil" (line 99), so a scan can lock a marker that is still published (W1). Required order: remove -> unlock -> close.
        statements:
          line 93: if claim == nil || claim.file == nil {
          line 96: unlockStartingMarker(claim.file)
          line 97: _ = claim.file.Close()
          line 98: claim.file = nil
          line 99: return os.Remove(claim.path) == nil
--- FAIL: TestStartingMarkerReleaseRemovesBeforeUnlock (0.00s)
```

```text
    agent_runner_starting_marker_order_test.go:403: 4 of 35 release/scan interleavings let a scan lock a published marker whose own release then deleted the name (W1 observed between 144 iterations, 62 residual deletes); launch ids: launch-release-stress-153, launch-release-stress-168, launch-release-stress-236, launch-release-stress-57
--- FAIL: TestStartingMarkerReleaseNeverExposesReleasedContent (3.01s)
```

```text
    agent_runner_starting_marker_order_windows_test.go:134: os.OpenFile handle blocks os.Remove: err=remove ...share-probe.txt: The process cannot access the file because it is being used by another process. errno=32 (0x20), matching the windows.go comment
--- PASS: TestStartingMarkerReleaseShareFactProbe (0.01s)
```

## 4. 跑测命令与耗时

```powershell
New-Item -ItemType Directory -Force tmp\DR13-FIX-D3\03-red | Out-Null
cmd /c "go test -count=1 -run TestStartingMarkerRelease -v ./internal/adapters/ > tmp\DR13-FIX-D3\03-red\go-test-red.txt 2>&1"
```

- 原始输出：`tmp/DR13-FIX-D3/03-red/go-test-red.txt`（cmd 重定向，字节级原文）
- 退出码：1（FAIL）
- 包内耗时：`4.265s`（其中压力回归窗口 3.01s）；整条命令含编译/启动 9.9s
- 未使用 `-race`（本机 Windows 无 gcc，与 R3 登记一致）

## 5. 可证伪性与覆盖说明

- 结构性红测的证伪路径：删除/移位 `os.Remove(claim.path)` 或把它改成 defer ⇒ 红；恢复为「remove → unlock → close」⇒ 绿。
  诚实边界：它是**源码顺序门**，不是运行时证明；运行时后果由语义直测、fail-closed 反测与压力回归分别钉住。
- 压力回归的计数口径（防假红）：仅当「观察时该 claim 处于 releasing」且「该 claim 的 release 最终返回 true（删名成功）」才计违例；
  两类设计已登记的残余（R5 删名被拒的残余 marker；D2 覆盖的孤儿 inode + 新 marker 占名，用 `os.SameFile` 区分）不计违例。
- 不可测项：unlock↔remove 的确定性运行时窗口（`release()` 无注入 seam，`unlockStartingMarker` 不可在测试中重定义）——
  按 design §2D(3) 判定为**不可测**，本片不引入生产 seam。

## 6. 生产代码零改动声明

本轮 ③ 仅新增上述两个 `_test.go` 文件；`internal/adapters` 下**零生产代码改动**（`git status --porcelain` 中
`internal/adapters` 只出现两个 `??` 新测试文件，`git diff --stat` 无 `internal/adapters` 条目）。
`docs/t027/REMAINING_SLICES.md` / `_EN.md` 的修改为进入本轮前的既有工作区状态，非本角色产生。

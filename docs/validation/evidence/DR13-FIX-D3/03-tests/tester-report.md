## ① 新增/修改的测试文件清单

| file:line | 内容 |
|---|---|
| `internal/adapters/agent_runner_starting_marker_order_test.go:154` | `TestStartingMarkerReleaseRemovesBeforeUnlock`（go/ast 源码顺序门，跨平台） |
| `internal/adapters/agent_runner_starting_marker_order_test.go:235` | `TestStartingMarkerReleaseRemovesThePublishedPath`（remove-while-open 直测） |
| `internal/adapters/agent_runner_starting_marker_order_test.go:261` | `TestStartingMarkerReleaseIsIdempotent` |
| `internal/adapters/agent_runner_starting_marker_order_test.go:296` | `TestStartingMarkerReleaseNeverExposesReleasedContent`（压力回归，3s 窗口） |
| `internal/adapters/agent_runner_starting_marker_order_windows_test.go:55` | `TestStartingMarkerReleaseFailClosedWhenTheNameIsOccupied`（无 FILE_SHARE_DELETE 占名反测） |
| `internal/adapters/agent_runner_starting_marker_order_windows_test.go:118` | `TestStartingMarkerReleaseShareFactProbe`（只 `t.Logf`，不硬断言） |
| `docs/validation/evidence/DR13-FIX-D3/03-tests/red-baseline.md` | 红基线（UTF-8+BOM+LF，已用 node 归一） |

测试内部辅助：`dr13SourcePath`:36 / `dr13ReleaseFuncBody`:50 / `dr13IsRemoveOfClaimPath`:111 / `dr13StatementListing`:133 / `dr13OccupyStartingMarkerName`(win):29。

## ② 红/绿结论（生产代码未修复）

| 测试 | 结果 | 依据 |
|---|---|---|
| `TestStartingMarkerReleaseRemovesBeforeUnlock` | **FAIL（确定性红，主证据）** | unlock(:96) 先于 os.Remove(:99) |
| `TestStartingMarkerReleaseNeverExposesReleasedContent` | **FAIL（补充红，R3 概率口径）** | 4/35 次观察命中「随后删名成功」的 claim（144 次发布 / 62 次残余删名） |
| `TestStartingMarkerReleaseRemovesThePublishedPath` | PASS | 旧序先 Close 后 Remove，天然满足；修复后须仍绿 |
| `TestStartingMarkerReleaseIsIdempotent` | PASS | 旧代码依 `claim.file=nil` 返回 false |
| `TestStartingMarkerReleaseFailClosedWhenTheNameIsOccupied` | PASS | 失败方向语义新旧同结论；修复后须仍绿 |
| `TestStartingMarkerReleaseShareFactProbe` | PASS（仅日志） | 本工具链 os.OpenFile 阻止 os.Remove，errno=32(0x20) |

## ③ 红证日志与关键摘录

- 日志：`tmp/DR13-FIX-D3/03-red/go-test-red.txt`（cmd 字节级重定向）
- 命令：`cmd /c "go test -count=1 -run TestStartingMarkerRelease -v ./internal/adapters/ > tmp\DR13-FIX-D3\03-red\go-test-red.txt 2>&1"`，退出码 1；包内 4.265s（压力窗口 3.01s），整条命令 9.9s；未用 `-race`

```text
release() drops the host lock before the published name is gone: "unlockStartingMarker(claim.file)" (line 96) runs before "return os.Remove(claim.path) == nil" (line 99), so a scan can lock a marker that is still published (W1). Required order: remove -> unlock -> close.
    line 96: unlockStartingMarker(claim.file)
    line 99: return os.Remove(claim.path) == nil
--- FAIL: TestStartingMarkerReleaseRemovesBeforeUnlock (0.00s)

4 of 35 release/scan interleavings let a scan lock a published marker whose own release then deleted the name (W1 observed between 144 iterations, 62 residual deletes)
--- FAIL: TestStartingMarkerReleaseNeverExposesReleasedContent (3.01s)

os.OpenFile handle blocks os.Remove: err=remove ...: The process cannot access the file because it is being used by another process. errno=32 (0x20), matching the windows.go comment
--- PASS: TestStartingMarkerReleaseShareFactProbe (0.01s)
```

## ④ 遇到的问题 / 对 ② 的注意事项

1. **AST 门的具体形状**：测试要求 `os.Remove(<claim.path 或由 claim.path 初始化的局部别名>)` 存在且**非 defer/go**，且其源码位置早于 `unlockStartingMarker(...)`。只约束这两者顺序，不限制 `Close` 写法（`defer close` 不会被判红，但 design 建议无 defer 分支）。
2. **返回值语义不得变**：`Test...RemovesThePublishedPath` 要求首次 `release()` 返回 **true**；`Test...IsIdempotent` 要求二次返回 **false** 且不 panic。建议 `removed := os.Remove(claim.path) == nil; unlock; close; claim.file = nil; return removed`。
3. **remove 失败仍须无条件 unlock+Close**：`Test...FailClosedWhenTheNameIsOccupied` 会重开（`openStartingMarkerForClear`）并 `tryLockStartingMarker` 断言锁已释放；「为保锁而不 Close」必红。
4. **不要动 claim 句柄共享模式**：`Test...RemovesThePublishedPath` 直接钉住 `createLockedStartingMarkerFile` 的 `FILE_SHARE_DELETE` 前提。
5. **R-N1 事实已实测**：本工具链 `os.OpenFile` 句柄**确实**阻止删除（errno 32），与 `_windows.go:97-101` 注释一致；因此在 Windows 高并发下 `release()` 会频繁返回 false（本轮 62/144 次残余），属 design A(2)/R5 收敛径，**不是** D3 回归，请勿在 release 里加 errno 映射或重试逻辑。
6. **压力测试防假红口径**（供 ② 复审）：仅当「观察时 claim 处于 releasing」且「该 claim 的 release 最终返回 true」才计违例；两类残余（删名被拒残留、D2 覆盖的孤儿 inode+新 marker 占名，用 `os.SameFile` 区分）不计。
7. **跨平台编译已验证**：`GOOS=linux go vet ./internal/adapters/` = 0；`go build ./...` = 0（Windows 腿同样 0）。
8. 工作区里 `docs/t027/REMAINING_SLICES.md` / `_EN.md` 的 `M` 状态为进入本轮前既存，非本角色产生。

## ⑤ 生产代码零改动确认

`internal/adapters` 零生产改动（`git diff --stat -- internal/adapters` 为空，仅两个 `??` 新测试文件）；除 gofmt 两个新测试文件与写入上述 `.md` 证据外未改动任何其它文件。

TEST: DONE

## 1. `release()` 改动后完整原文（`internal/adapters/agent_runner_starting_marker.go:89-120`）

```go
// release removes the published marker, drops the host lock, and closes the
// claim handle, in that order. It reports whether the marker file was actually
// removed; a claim that was already released, or a marker another scan already
// cleared, reports false.
//
// The path is removed first, while the claim handle still holds the lock. That
// makes the deletion and the lock release one continuous step from the point of
// view of a recovering scan: the name is already gone by the time the lock can
// look free, so a scan can never pair an idle lock with a live marker still
// published under that name.
//
// The removal can still fail, for example when a concurrent scan holds the name
// open without allowing deletion. The lock and the handle are then released
// unconditionally anyway: the kernel drops the lock when the handle closes or the
// host dies, so withholding the close to keep the lock would only leak a handle
// and cannot keep the marker protected. Such a residual marker fails closed into
// the existing convergence path - the same-process ownership gate refuses it
// here, and the next scan re-checks it - rather than being retried inside
// release.
func (claim *startingMarkerClaim) release() bool {
	if claim == nil || claim.file == nil {
		return false
	}
	removed := os.Remove(claim.path) == nil
	unlockStartingMarker(claim.file)
	_ = claim.file.Close()
	claim.file = nil
	return removed
}
```

顺序：`os.Remove` → `unlockStartingMarker` → `Close` → `claim.file = nil` → `return removed`。无 defer、无重试、无 errno 映射；签名与返回语义不变。

## 2. 六项自验结果表（证据均落 `tmp/DR13-FIX-D3/04-impl/`）

| # | 命令 | 退出码 | 关键行 |
|---|---|---|---|
| a | `gofmt -l internal/adapters/agent_runner_starting_marker.go` → `gofmt.txt` | 0 | 输出为空（0 行，即无待格式化文件） |
| b | `go test -count=1 -run TestStartingMarkerRelease -v ./internal/adapters/` → `go-test-focused.txt` | 0 | 6/6 PASS：`RemovesBeforeUnlock` 0.00s、`RemovesThePublishedPath` 0.01s、`IsIdempotent` 0.01s、`NeverExposesReleasedContent` 3.05s、`FailClosedWhenTheNameIsOccupied` 0.02s、`ShareFactProbe` 0.01s；`ok ... 4.359s`，无 FAIL |
| c1 | `go test -count=1 ./internal/adapters/` → `go-test-adapters.txt` | 0 | `ok github.com/larsonzh/prfrail/internal/adapters 48.625s`（FAIL=0） |
| c2 | `go test -count=5 -run "…NeverExposesReleasedContent\|…ConcurrentStartsHaveOneWinner" ./internal/adapters/` → `go-test-stress.txt` | 0 | `ok ... 20.808s`（5×压力 3s 窗口 + 5×并发单赢家，0 违例） |
| d1 | `go build ./...` → `build.txt` | 0 | 空输出 |
| d2 | `go vet ./...` → `vet.txt` | 0 | 空输出 |
| e | `GOOS=linux go vet ./internal/adapters/` → `vet-linux.txt` | 0 | 空输出（跨平台编译检查通过） |

## 3. 改动范围确认

```
git diff --stat -- internal/
 internal/adapters/agent_runner_starting_marker.go | 25 +++++++++++++++++++----
 1 file changed, 21 insertions(+), 4 deletions(-)
```

工作区其余 `M`：`docs/t027/REMAINING_SLICES.md` / `_EN.md`（进入本环前既存，同 ③ 报告第 8 条）；`??` 为 ③ 的两个测试文件与证据目录。**未触碰**：`writeAgentRunnerStartingMarker` 的 discard 闭包、`clearAgentRunnerStartingMarker`、windows/unix 锁实现、`closeStartingMarker`/`releaseLaunch`/`discardStartingLaunch`、任何测试文件。

## 4. 问题与对 ⑤ 审查的说明

1. **命令引号坑（已修复，非代码问题）**：首跑 c2 用 `cmd /c "… -run A|B …"` 时外层引号被 cmd 剥除，`|` 被当成管道，报「不是内部或外部命令」且产出无意义输出（退出码 255）。改用 PowerShell 单引号包裹、内层保留 `"A|B"` 后重跑，已覆盖同一证据文件。
2. **`vet-linux.txt` 的落地过程需如实说明**：我先用默认 GOOS 写了一次该文件，再用 `$env:GOOS='linux'` 覆盖重跑；两次均 exit 0、输出为空，文件中记录的是 `GOOS=linux` 那一次结果，环境变量已 `Remove-Item`。
3. **fail-closed 语义已由测试实证**：`ShareFactProbe` 本轮日志仍显示 `os.OpenFile` 句柄阻止 `os.Remove`（errno 32），即 Windows 高并发下 `release()` 返回 false 属 design A(2)/R5 收敛径，不是 D3 回归；因此我**没有**加 errno 映射或重试，请 ⑤ 重点核对这一点是否与「删除失败仍无条件 unlock+Close」一致。
4. **请 ⑤ 复核的落点**：`:113` 的 `removed :=` 在 `unlock`/`Close` 之前求值，确保「删除是否成功」与「锁/句柄是否释放」解耦——返回值不得因后续 unlock 失败而变（unlock 返回值仍被忽略，与 design A(3) 一致）；注释未引用 ③ 判据字符串形态，仅陈述顺序与 fail-closed 理由。
5. 测试文件零改动：6 个红测全部转绿是唯一改动效果，无断言调整、无 fixture 变更。

IMPLEMENT: DONE

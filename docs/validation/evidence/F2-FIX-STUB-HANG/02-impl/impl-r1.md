# ② 实现回执（impl-r1）

- 切片：F2-FIX-STUB-HANG（层 1 第 10 项·修复类）；角色：prfrail-implementer（DeepSeek V4.1 Flash）；落盘 2026-11-09。
- 仓库 HEAD：`de37c58`；工作目录 `D:\LZProjects\prfrail`。纪律：未 commit / 未 push；仅触碰 `tools/agent-stub/main.go` 一个文件。

## 变更（唯一变更点，与 design-r1 §B 精确 diff 一致）

```go
	case "hang":
		// No artifacts beyond its own pid: the run is stopped from the outside.
		run.writeOwnPID()
		// A hang never spawns a child, so the empty report is the positive fact
		// the constant-write contract requires: consumers can tell "no child"
		// apart from "report never written" (the latter must stay fail-closed).
		run.writeEmptyChildPIDReport()
		time.Sleep(time.Hour)
		return 0
```

- 插入点：`case "hang"` 分支内 `run.writeOwnPID()` 之后、`time.Sleep(time.Hour)` 之前（main.go L112–120）。
- `writeEmptyChildPIDReport()` 为既有方法（L264–268，写 0 字节 `stub-child.pid`，失败 `os.Exit(11)`），未重复实现。
- hang 分支不落到公共路径（L121 起的 `case "run", …` 与 L135+ 公共写路径原结构未动）。
- 注释英文（国际化硬规则）；缩进 tab；`.go` = UTF-8 无 BOM + LF。

## 磁盘实况核验（主控，防脏缓冲）

- 终端 `Get-Content` 直读磁盘：hang 分支含 `run.writeEmptyChildPIDReport()`，与派单 diff 一致（本回执落盘时点）。
- ③ 绿测与 layer3 重跑均以磁盘编译产物为准。

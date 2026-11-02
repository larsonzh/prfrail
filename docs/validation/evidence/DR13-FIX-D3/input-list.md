# input-list — `DR13-FIX-D3`（P48）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-02，`HEAD` = `e742224`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/validation/dr13-fix.md` | 24008 | DR13-FIX 报告：D3 缺口定义（R4 放行条件与边界簿记） |
| `internal/adapters/agent_runner_starting_marker.go` | 4776 | D3 落点：`startingMarkerClaim.release()`（现序 unlock → close → remove） |
| `internal/adapters/agent_runner_starting_marker_windows.go` | 6140 | Windows remove-while-open 语义（LockFileEx / 共享删除模式） |
| `internal/adapters/agent_runner_starting_marker_unix.go` | 3515 | Unix 侧锁语义对照 |
| `internal/adapters/agent_runner_pinned_cli.go` | 43978 | release 链入口（`releaseLaunch` / `closeStartingMarker`） |
| `docs/t027/CHECKLIST_v37.md` | 10765 | 跨会话快照（锚点、四层队列、纪律 T / U） |

## 未复制入包

Go 源码与测试不复制；如需引用仅登记路径与行号。

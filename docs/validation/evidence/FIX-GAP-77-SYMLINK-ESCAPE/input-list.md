# input-list — `FIX-GAP-77-SYMLINK-ESCAPE`（P53）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-03，`HEAD` = `dc785ed`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/t027/CHECKLIST_v42.md` | 23383 | 跨会话快照（锚点、四层队列、§八纪律 A–V、§九清单元规则与门禁坑 15 条） |
| `docs/t027/REMAINING_SLICES.md` | 287127 | 台账登记落点（状态表条目 + 疑似漏洞上报 GAP-77 行） |
| `docs/t027/REMAINING_SLICES_EN.md` | 321603 | 台账英文镜像落点 |
| `internal/adapters/bounddir_unix.go` | 4108 | 生产修复对象（:64 `openat(O_NOFOLLOW)` 仅护末组件） |
| `internal/adapters/bounddir_escape_unix_test.go` | 3356 | T51-1 `TestBoundDirMiddleSymlinkEscapesBound`（`t.Skip` 待移除，红转绿） |
| `internal/adapters/bounddir.go` | 1253 | 目录屏障公共契约面（close 语义与错误哨兵，接口签名冻结参考） |
| `internal/adapters/bounddir_windows.go` | 10584 | Windows 对向实现（本片不改，只作对向 GOOS 复核参考） |
| `docs/validation/acceptance-security-counterexamples.md` | 12316 | P52 报告（GAP-77 击穿实证的登记来源） |
| `docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/t51-1-linux-escape.md` | 3654 | T51-1 攻击载荷与击穿证据原文（读出界外 23 字节） |
| `docs/SECURITY.md` | 8570 | ⑧ 规范同步落点（SEC-02 威胁矩阵） |
| `docs/SECURITY_EN.md` | 9387 | SECURITY 英文镜像落点 |
| `docs/TEST_STRATEGY.md` | 11214 | ⑧ 规范同步落点（AT-02 或相关基线，强制反例） |
| `docs/TEST_STRATEGY_EN.md` | 12072 | TEST_STRATEGY 英文镜像落点 |
| `docs/ARCHITECTURE.md` | 18344 | ⑧ 规范同步落点（§5.1 安全反例矩阵） |
| `docs/ARCHITECTURE_EN.md` | 21067 | ARCHITECTURE 英文镜像落点 |
| `docs/DEV_PLAN.md` | 80029 | ⑧ P53 任务状态更新落点 |
| `docs/DEV_PLAN_EN.md` | 91596 | DEV_PLAN 英文镜像落点 |

## 未复制入包

Go 源码与测试不复制；如需引用仅登记路径与行号。P52 证据包其余载荷（①/③/④/⑤/⑥/⑦ 归档原文）按需引用，不复制入本包。

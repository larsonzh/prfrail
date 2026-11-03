# input-list — `ACCEPTANCE-SECURITY-COUNTEREXAMPLES`（P52）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-03，`HEAD` = `941f60e`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/t027/CHECKLIST_v41.md` | 21931 | 跨会话快照（锚点、四层队列、§八纪律 A–V、§九清单元规则与门禁坑） |
| `docs/t027/REMAINING_SLICES.md` | 285808 | 台账登记落点（状态表条目 + 完成记录） |
| `docs/t027/REMAINING_SLICES_EN.md` | 320215 | 台账英文镜像落点 |
| `docs/validation/evidence/STRATEGIC-REVIEW-2026-11-02/02-round2/review.md` | 14230 | P46 轮 2 缺口台账（GAP-51/52/55/60/64/70 定义原文） |
| `docs/validation/strategic-review-2026-11-02.md` | 13312 | P46 战略回看：断点 A 结论、第 10 位（接受前）切片定义 |
| `internal/adapters/bounddir.go` | 1253 | 目录屏障公共契约面（close 语义与错误哨兵） |
| `internal/adapters/bounddir_unix.go` | 4108 | GAP-51 相关：Unix 目录绑定实现（两段式 dirfd 缺口候选） |
| `internal/adapters/bounddir_windows.go` | 10584 | GAP-52 相关：Windows reparse 前缀比较实现 |
| `internal/adapters/bounddir_clean.go` | 1214 | 目录屏障清理实现 |
| `internal/evidence/boundread.go` | 7912 | GAP-55 相关：有界读原语（关闭拒读/超界拒摘要/同源对抗） |
| `internal/evidence/boundopen_unix.go` | 356 | GAP-51 相关：Unix 有界打开 |
| `internal/evidence/boundopen_windows.go` | 1177 | GAP-52 相关：Windows 有界打开 |
| `internal/gates/probe_permit.go` | 12926 | GAP-64 相关：LoadProbePermit 等门禁 fail 分支 |
| `internal/gates/readiness.go` | 9321 | GAP-70 相关：GoalProduction / GoalAcceptance 判定与 AT23Hash / SealBinding 校验 |
| `internal/chain/seal.go` | 8760 | GAP-60 相关：seal 记录与 K 绑定校验（ContextKey validate） |
| `internal/chain/seal_verify.go` | 2281 | GAP-60 相关：恢复期封条核验（verifyAcceptedSeals） |

## 未复制入包

Go 源码与测试不复制；如需引用仅登记路径与行号。

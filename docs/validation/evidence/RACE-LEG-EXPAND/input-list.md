# input-list — `RACE-LEG-EXPAND`（P49）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-02，`HEAD` = `81dd49f`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `.github/workflows/ci.yml` | — | 修改落点：Race 腿目标包清单（第 78–80 行区域） |
| `docs/validation/p45-e2e-zero-model.md` | — | 前置阅读：候选登记「`-race` 腿扩包至 `internal/console`：本机缺 cgo，并发腿未过竞态检测」 |
| `docs/t027/CHECKLIST_v38.md` | — | 跨会话快照（锚点、四层队列、§九门禁坑、纪律 T / U） |
| `docs/t027/REMAINING_SLICES.md` | — | 台账登记落点（状态表 + P49 小节） |
| `docs/t027/REMAINING_SLICES_EN.md` | — | 台账登记落点（英文镜像） |
| `docs/DELIVERY_DIRECTIVE.md` | — | 准则 v1.60（流水线角色、门禁、审查纪律） |

## 未复制入包

Go 源码与测试不复制；如需引用仅登记路径与行号。`ci.yml` 全文不复制入包，仅在审查环节以 git diff 呈现。

## 锚点核对（步骤 0 实况）

| 项 | 值 |
|---|---|
| HEAD / origin/main | `81dd49f`（工作区干净） |
| 准则版本 | v1.60 |
| SELFTEST | 697/697 |
| 三作用域 | index / tree / ci 全 TOTAL_FAIL=0 |
| tmp 断言 | 顶层 = `.gitkeep` + `gates-cache`；`gates-cache` 下仅 `g4b` |

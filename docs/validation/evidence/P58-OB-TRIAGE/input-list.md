# P58 审查输入清单（input list）

| 序 | 输入 | 路径 | 用途 |
|---|---|---|---|
| 1 | 准则观察项台账 | `docs/DELIVERY_DIRECTIVE.md` §12.4（自「### 12.4 观察项台账」起至「### 12.5」前） | 唯一事实源：五条 OB 的当前状态与处置列 |
| 2 | 准则缺口清单 | `docs/DELIVERY_DIRECTIVE.md` §12.5（自「### 12.5 缺口与未来扩展」起至本节末） | 判断是否存在与本片相关的缺口行 |
| 3 | 台账 P7 行 | `docs/t027/REMAINING_SLICES.md` 的 `### P7 — 治理侧拆出 OB 汇总` 小节 | 承载集出处 |
| 4 | 台账 P4 行 | `docs/t027/REMAINING_SLICES.md` 的 `### P4 — DIRECTIVE-OB-SWEEP` 小节的「建议拆出」表 | 交叉核对 |
| 5 | 执行队列 | `docs/t027/CHECKLIST_v47.md` §四 层 2 的 `G-1` 行 | 本片在队列中的位置与起跑条件 |
| 6 | 本报告草稿 | `docs/validation/p58-ob-triage-assessment.md` 与 `_EN.md` | 被审对象 |
| 7 | 自由度清单 | `docs/validation/evidence/P58-OB-TRIAGE/freedom-list.md` | 越界判定的依据 |

## 输入一致性断言

- 上列七项在派单前逐件存在，缺件即停止派发（§5.1 的 OB-93 时序说明）。
- 输入未附任何既有审查清单（⑤ 的独立性要求）。
- 输入已消毒：无 IP、用户名、主机名、SSH 路径、凭据或 token。

## 旧口径检索记录

- 在 `docs/t027/REMAINING_SLICES.md` 的 P7 行内检索到五个旧行号锚点：`:860`（OB-13）、`:903`（OB-48）、`:923`（OB-68）、`:930`（OB-75）、`:932`（OB-77）。
- 上述锚点与 `docs/DELIVERY_DIRECTIVE.md` 当前行号不一致，属登记时快照；按 §11.1 不改写，本片按 OB 编号解析。
- 在 `docs/DELIVERY_DIRECTIVE.md` §12.4 内检索「拆出」与「P7」两串，台账行内无命中；两串只出现在 `docs/t027/REMAINING_SLICES.md` 的 P4 与 P7 行。

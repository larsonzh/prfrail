# 00 文档清单与继承事实 — P38-T027B-TACTICAL-BLUEPRINT

切片：`P38-T027B-TACTICAL-BLUEPRINT`（P38，评估类 `[SLICE]`）。准则 v1.59。建包时点 2026-10-30（`HEAD` = `344c206`）。

## 一、继承事实（P37，只作输入，不重复举证）

- 报告对：`docs/validation/t027b-breakthrough{,_EN}.md`（各 78 行）。
- 证据包：`docs/validation/evidence/P37-T027B-BREAKTHROUGH/`（载荷 6 件，`PACK-CONSISTENT`）。
- 提交链：`1bf57a4`（主）→ `00a4627`（CI 回填）→ `fa0fa3a`（台账闭环）→ `344c206`（R 计数修正）。
- CI：`41034345827` / `41036123038` / `41038767155` / `41041824162`，每次 Ubuntu + Windows 双腿 `success`。
- P37 两轮正文是本片第 1 轮的固定输入：Sol 111 行 / 10,723 B，Astra 68 行 / 9,277 B。

## 二、本片文档清单

| 文档 | 版本 / 位置 | 本片用途 |
|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` / `_EN.md` | v1.59；头部 `:3`；§0.3 `:105`；§12.4 `:1154`（OB-131 行） | 步骤 0 升版与 OB-131 登记 |
| `docs/t027/REMAINING_SLICES.md` / `_EN.md` | P38 行与 P38 块（含 R/D 基线与临时件落点） | 台账登记 |
| `docs/CONTRACTS.md` §7 | `:188–:283`（96 行 / 48,064 B） | 第 1 轮输入：AgentRunner 与辅助消息通道边界 |
| `docs/validation/evidence/P37-T027B-BREAKTHROUGH/01-sol-design/review.md` | 111 行 / 10,723 B | 第 1 轮输入：架构设计书 |
| `docs/validation/evidence/P37-T027B-BREAKTHROUGH/02-astra-critique/review.md` | 68 行 / 9,277 B | 第 1 轮输入：C1–C7 缺陷源 |
| `internal/adapters/**` 非测试 `.go` | 38 件 / 395,919 B | 第 1 轮输入：代码事实 |
| `internal/chain/**` 非测试 `.go` | 18 件 / 184,724 B | 第 1 轮输入：代码事实 |

## 三、输入包合计（实测）

- 条目 **59**（3 件文档段 + 56 个 `.go`），行数 **14,714**，字节 **648,707**（≈ 633.5 KB）。
- 逐件明细见 `01-sol-blueprint/input-list.md`。

## 四、本片不触碰的域

代码、契约文本、schema、fixtures、`tools/` 一律不改；本片对准则的改动仅限 OB-131 登记、头部升 v1.59、§0.3 日志行与 P38 台账登记。

## 五、执行阻力与分批记录

- 第 1 轮首次派单连续两次因单次输出长度上限触发**自重启 + 截断**：原始件保留为 `tmp/P38-T027B-TACTICAL-BLUEPRINT/r1-attempt1-selfrestart.md`（字节 42,509）与 `r1-attempt2-restart-truncated.md`（字节 15,428），二者**未**作为转录底稿。
- 用户裁定改为分批：批次 A（C1–C4）→ `01-sol-blueprint/review.md`；批次 B（C5–C7）→ `01-sol-blueprint/review-batch-b.md`。两批输出字符数 分别为 10,052 与 9,808。
- 批次 A 返回件仍含一次自重启（第二遍以 `## C1` 重述并在代码块中截断）：主控**只取第一遍完整段**（C1–C4 与跨项约束；行数 47），第二遍作废留证。
- 第 2 轮 Astra 红队输入仅两份施工图（未附 0.6 MB 代码包）；输出行数 22，字节 4,245，结论 `VERDICT: REVISE`。

# P94 审查与设计输入清单（input-list）

> 本清单记录本片读取的权威输入与其用途，供 ⑤ / ⑥ / ⑦ 复核输入完整性。

## 1. 锚点与快照

| 输入 | 用途 | 实测读数 |
|---|---|---|
| `docs/t027/CHECKLIST_v77.md` / `_EN.md` §一 | 步骤 0 锚点七项 | HEAD = origin = `9217915`；准则 v1.76；SELFTEST 757/757；三作用域 `TOTAL_FAIL=0`；v75 / v76 = 215 / 216（LF） |
| `docs/t027/CHECKLIST_v77.md` §四 / §五 / §六 | 下一项定义、P93 移交 a / c 登记、待裁定项 | 下一项 = `GATES-EXEMPT` 合并片（M 档，⑥⑦ 必做，须用户点名） |
| 机械读数 | ADR / OB / DR / 试点顺延 | ADR max 48；OB max 145；DR max 18；试点顺延第 25 次 |
| `tmp` 顶层断言 | 临时件位置闭集 | `.gitkeep` / `gates-cache`（其下仅 `g4b`）/ `find-quotes.js` |

## 2. 规则与判据

| 输入 | 用途 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` §6.3 | G2 / G4b / G6（含 G6-4）判据条款全集与作用域口径 |
| `docs/DELIVERY_DIRECTIVE.md` §12.4 | OB 台账现状（`GATES-EXEMPT` 无条目） |
| `docs/DELIVERY_DIRECTIVE.md` §0.3 / 头部 / 附录 C | 版本包落点与 C.0bd 位置 |
| `tools/gates/lib/criteria/g2.js` | 占位符判据实现（豁免通道候选落点） |
| `tools/gates/lib/criteria/g4b.js` | 新汉字语料判据实现与缓存口径 |
| `tools/gates/lib/criteria/g6.js` | G6 族（含 G6-4 计数型字面量与 P93 的 `sixstage` 作用域子规则） |
| `tools/gates/selftest.js` | T23 零预载三断言与新增断言落点 |
| `tools/gates/g6-whitelist.txt` | 白名单零条目现状与「作用域豁免不是白名单职责」注释头 |

## 3. 出处与移交

| 输入 | 用途 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` 附录 C.0ap（P33 段） | `GATES-EXEMPT` 拆片出处（①-3(iv) 的字面措辞） |
| `docs/validation/directive-v21-obs-sweep.md` 尾部建议 | ①-3(iv) 的原始意图（报告 §1 异常登记区的 `gate-exempt` 区间） |
| `docs/validation/g6-4-transcript-exempt.md` §15 与 §1 | P93 的三选一裁定与未选项去向 |
| `docs/validation/evidence/P93-DIRECTIVE-G6-4-TRANSCRIPT-EXEMPT/01-design.md` §3 | P93 移交 (a) / (c) 的逐条理由原文 |

## 4. 工程规范

| 输入 | 用途 |
|---|---|
| `docs/CODING_CONVENTIONS.md` | `.js` / `.txt` = 无 BOM + LF；`.md` = BOM + LF |

## 5. 未纳入的输入（如实披露）

- 未读取 `internal/**`：本片边界明确不改 Go 实现。
- 未读取 CONTRACTS / schema / fixtures：本片不触协议、不触夹具。
- 未启用 Ubuntu VM 与代理；CI 双腿按 §11.2 观察记录。

# P95 审查与设计输入清单（input-list）

> 本清单记录本片读取的权威输入与其用途，供 ⑤ / ⑥ / ⑦ 复核输入完整性。

## 1. 锚点与快照

| 输入 | 用途 | 实测读数 |
|---|---|---|
| `docs/t027/CHECKLIST_v78.md` / `_EN.md` §一 | 步骤 0 锚点七项 | HEAD = origin = `22370fe`；准则 v1.77；SELFTEST 770/770；三作用域 `TOTAL_FAIL=0`；v76 / v77 = 216 / 217（LF），v78 = 211（LF，自复） |
| `docs/t027/CHECKLIST_v78.md` §四 / §五 / §六 | 层 3 性能候选（`GATES-CRITERIA-ASYNC` 记录态）与队列状态 | 下一项 = 无（队列清空，暂停排队）；`GATES-CRITERIA-ASYNC` = 层 3 记录态 |
| 机械读数 | ADR / OB / DR / 试点顺延 | ADR max 48；OB max 145；DR max 18；试点顺延第 25 次 |
| `tmp` 顶层断言 | 临时件位置闭集 | `.gitkeep` / `gates-cache`（其下仅 `g4b`）/ `find-quotes.js` |

## 2. 规则与判据

| 输入 | 用途 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` §2.2 | 角色—载体—模型三元组（① 架构师限定名、批量辅助首选与 ⑤ 备用档型号） |
| `docs/DELIVERY_DIRECTIVE.md` §5.2 | ⑦ 选型矩阵（深度档判据：并发 / 协议语义；适用域排除句） |
| `docs/DELIVERY_DIRECTIVE.md` §7.5 | ⑦ 额度分层（M 档 = 1+2）与计价真值源 |
| `docs/DELIVERY_DIRECTIVE.md` §13.1 / §13.2 | 起跑与步骤 0 通则、跨片协作纪律（派单表机械读取 / 串行 / 输出上限） |
| `docs/CODING_CONVENTIONS.md` | `.js` = UTF-8 无 BOM + LF；`.md` = UTF-8 with BOM + LF |

## 3. 代码与出处

| 输入 | 用途 |
|---|---|
| `tools/gates/gate.js` | `runCriteriaConcurrent` 调度契约与 `RUN_ORDER` 回填现状（:140/:231 起） |
| `tools/gates/lib/criteria/g6.js` / `g7.js` / `g8.js` | 三个最重判据的同步串行实现现状 |
| `docs/validation/gates-perf-assess{,_EN}.md` §7 建议 4 | 出处 + G6/G7/G8 耗时基线 |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | 台账 P95 段落点与格式 |

## 4. 本片工件

| 输入 | 用途 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10） |
| `freedom-list.md` | 执行自由度清单（R/D = 0 基线） |

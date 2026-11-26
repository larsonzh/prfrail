# P93 执行自由度清单（⑥ 后、⑦ 前落盘）

本清单列出本片在「怎么做可自主」范围内做出的实质性选择（§1.7 执行自由度）。所有选择均在 ① 裁定的通道 b 内，未改变契约语义。

1. **单源复用方式**：`g6.js` 以 `require('./g7')` 取 `internals`（`targets` / `sixstageBlocks` / `SIXSTAGE_SECTION_RE` / `SIXSTAGE_SECTION_SEQ`），而不是在 `g6.js` 复制围栏正则。理由：避免「同一口径写两处必然漂移」；`g7.js` 只依赖 `../ctx`，无循环依赖。可证伪：T140 断言 `g6.js` 不含本地围栏字面量。
2. **豁免区间取 `[start, end)`**：含开栏行、不含闭栏行。理由：闭栏行是裸 ``` 不可能命中四条正则；与 G7-f 的「触及」区间口径一致。
3. **结构面作为 ⑥ r1 Medium 的整改方式**：要求围栏体按序出现 `Section A`–`Section E`。理由：不引入新机制（哈希锚定 / 登记位）即可消除「外观围栏」这个最廉价的逃避面，且与 G7-f 判定同源。
4. **`.trim()` 对齐**：`sixstageStructurallyValid` 对每行先 `trim()`，与 `g7.js` 的 `SIXSTAGE_SECTION_RE.exec(String(body[k]).trim())` 逐字同口径。
5. **`fence-exempt` 的落点**：只在无 finding 时进入 detail，并通过把准则措辞精确化为「无 finding 时」承接 ⑤ 的 Low；不改共享 `emit`（改它会波及 G6-1 / G6-2 / G6-3 / G6-5 / G6-6 / G6-7 / G6-8 的 detail 语义）。
6. **T133 的定位方式**：以子串 `fence-exempt 13` 钉住「开栏 1 行 + 体 12 行」的区间算式；体行数变化时该断言会随之失配（预期耦合，已在测试注释说明）。
7. **报告对生成方式**：用数组脚本一次性写两份，脚本内断言 `cn.length === en.length`，从机制上保证 G3(a.1) 行数对称。
8. **新断言编号**：自 T133 起（实测既有最大为 T132），避免与 P81 / P87 的 T131 系列撞号。

# ③ 测试回执（prfrail-tester · DeepSeek V4.1 Flash）—— 逐字转录

**TEST: FINDINGS** — 夹具层已补齐且全绿（717/717），但探针发现 `G5(e2)` 存在一处实现级盲区（不是夹具问题），交 ⑤ 裁定；未改 `g1.js`/`g5.js`。

## ① T119 / T120 逐断言核对（全部真实，非空判）

用临时探针把每条输入的 `G1(e2)`/`G5(e2)` **record detail** 打出来（不只是 verdict），读数如下：

| 断言 | 实际 record | 真实性 |
|---|---|---|
| T119-1 变更行 + `✅` ⇒ PASS | `PASS :: 1 judged, 0 skipped, 0 historical` | 真命中 PASS 分支（计数证明不是"没收集到所以 PASS"）|
| T119-2 变更行 + `🔄` ⇒ FAIL | `FAIL :: ["docs/t027/REMAINING_SLICES.md:11 in-block completion row but heading A0=🔄"]` | 红基线真实 |
| T119-3 历史行（added 为空）⇒ PASS 不追溯 | `PASS :: 0 judged, 0 skipped, 1 historical` | 真走 historical 分支，非解析失败 |
| T119-4 变更行 + 标题无状态符号 ⇒ PASS 跳过 | `PASS :: 0 judged, 1 skipped, 0 historical` | 真走 skip 分支 |
| T120-1 路径可解析 ⇒ PASS | `PASS :: 1 reference(s) resolved` | 真解析到 1 条引用 |
| T120-2 路径不可解析 ⇒ FAIL | `FAIL :: ["docs/validation/ghost.md"]` | 红基线真实 |
| T120-3 无变更块内行 ⇒ PASS | `PASS :: no added in-block completion row (nothing to resolve)` | 真走空集合分支 |

结论：七条断言全部真实；原夹具的弱点只在**只断言 verdict**，故 PASS 分支的"可达性"隐含依赖反例。

## ② 补的夹具（T121 / T122，仅 `tools/gates/selftest.js`）

- **T121**（G1(e2)，5 条）：① `✅` PASS 必须带 `1 judged, 0 skipped, 0 historical` 计数；② 历史行必须带 `1 historical`；③ **无任何标题**时走 `!c.headKey` 半边跳过（T119 只覆盖 `!c.headSymbol` 半边）；④ **就近绑定**：`✅` 标题上面再压一个 `🔄` 标题，必须 FAIL 且点名 `heading B1=🔄`；⑤ EN 字段名 `| Completion record |` 同样被收集并判定。
- **T122**（G5(e2)，2 条）：① 正例必须带 `1 reference(s) resolved`（排除"零解析仍 PASS"）；② EN 字段名行同样被提取并解析。
- 附带：`g1Probe`/`g5Probe` 各增一个 `e2detail` 返回字段（夹具层改动，不影响既有调用）；文件头覆盖清单更新为 `T119-T122`。
- 增量：**710 → 717（+7）**。注：`02-impl/impl-r1.md` 内 `SELFTEST: PASS 710/710` 与报告中的该数字将变旧，需 ⑤/② 回写（我不改文档）。

## ③ SELFTEST verbatim

```
SELFTEST: PASS 717/717
```

（`node --check tools\gates\selftest.js` 通过；文件 UTF-8 **无 BOM**（首 3 字节 `39,117,115`）+ **LF**（CRLF=0 / LF=6963））

## ④ 工作区探针读数（真实台账）

```
PASS G1(e2) :: 0 judged, 0 skipped, 40 historical (v1.64)
docs/t027/REMAINING_SLICES.md     -> inblock rows = 20
docs/t027/REMAINING_SLICES_EN.md  -> inblock rows = 20
```

与 ② 报告的「40 historical 0 judged」**一致**（合计 40 行、0 判定）。

## ⑤ FINDINGS（交 ⑤，未改判据）

**F-1（Medium）`G5(e2)` 对现行台账约定完全空判** — `tools/gates/lib/criteria/g5.js:81` 与 `:86`

- 现象：`if (raw.indexOf('{') >= 0) continue;` 会丢弃花括号展开路径。台账现行写法是 `报告对：`docs/validation/<slice>{,_EN}.md``，实测（当前工作区文本）：
  - CN：20 条块内完成行，**14 条**含 `{,_EN}`，其中 **13 条解析出 0 条引用**（另 1 条只靠附带的正规路径出 1 条）；
  - EN：同样 14 / 13。
  - 合成验证：新增 `| 完成记录 | 报告对：`docs/validation/x{,_EN}.md` |` ⇒ `PASS :: no added in-block completion row (nothing to resolve)`，即**新切片沿用现行写法时 `G5(e2)` 永远不判**，本片"收口 OB-138 盲区"的目标对该写法未达成。
- 说明：准则 v1.64 文本写的是"同 (e) 口径"，代码与字面一致；但 (e) 的这个跳过没有任何夹具保护（selftest 中无相关断言），也无文档登记。故属"实现覆盖盲区"，请 ⑤ 二选一：
  1. 修 `evidencePathsFromRow`：把 `{a,b}` 展开为候选路径（至少 `{,_EN}` → 两个文件），要求可解析；
  2. 或在 §6.3 + OB-138 登记该边界，并规定"新块内完成行的证据列必须写出两条完整路径"（但那样需 ⑤ 确认是否算规则语义变更）。
- 影响面：`G5(e)`（追加表日期行）同源同病；`G1(e2)` 判"是否已完成"不看路径，**不受影响**。

**F-2（Low）`g5.js:231-239` detail 文案把"有新增行但无路径"说成"没有新增行"**

- 现象：`refs2.length === 0` 时统一输出 `no added in-block completion row (nothing to resolve)`，即使确实存在新增块内行（F-1 即触发此文案）。建议区分：`no added in-block row` / `N added row(s) but no parsable evidence path`，并给出两个计数（行数、引用数），使静默跳过可见。

**G1(e2) 实现**：未发现缺陷（计数语义、就近绑定、无标题跳过、历史排除、EN 字段名均符合 v1.64 文本）。

## ⑥ 越界声明

- 仅修改 `tools/gates/selftest.js`（+128 / −3）；未改 `g1.js` / `g5.js` / 任何文档 / 白名单。
- 未 `git add`、未 commit、未 push；未派单。
- 临时探针脚本置于 `tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/` 并**已整体删除**；`tmp/` 下残留 `gates-cache`、`gate-v49-tree.txt` 为**他人既有**产物，未触碰。

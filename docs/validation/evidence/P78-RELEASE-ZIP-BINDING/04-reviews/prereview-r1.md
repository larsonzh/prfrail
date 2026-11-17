# ⑤ 预审 r1（DeepSeek V4 Pro，只读）

切片：P78 `RELEASE-ZIP-BINDING`。审查对象：方案书 r1 + 落库 diff（28 文件）+ 复验证据 + 门禁原始输出。

```
PRE-REVIEW: PASS WITH FIXES
```

## 发现清单

| # | SEVERITY | file:line | 被违反的句子或判据 | 最小修复 |
|---|---|---|---|---|
| 1 | Medium | `docs/validation/release-zip-binding.md:79-85`（`_EN` 同位） | §9 门禁读数表宣称 index / tree / ci 三作用域 `TOTAL_FAIL=0`、§8 宣称「④ 门禁全绿」；但现存变更集捕获件为 tree=`TOTAL_FAIL=2`（`04-gates-tree-pre.txt`）、index=`TOTAL_FAIL=1`（`pre-index.txt`，G5(d) 暂存瞬态）。§9 数值仅与步 0 空变更集捕获（`00-*`，vacuous pass）一致，是对未发生读数的预先陈述 | §9 标注「待 ⑥ 围栏注入 + ⑨ 暂存后回填」，或如实登记当前读数与两处红因处置；⑧a 前必须用变更集范围重跑三作用域并归档 `06-gates/` |
| 2 | Medium | `docs/validation/s1-exit.md:7`、`:56-58`（`_EN:7`、`:56-58` 同位） | 头部结论段仍写「最终 ZIP 尚未绑定和复验，EOL 通知方式尚未冻结」；门表 3 行证据列仍写「尚未绑定 / 仍须复验 / 尚未绑定正式发行物」——与 P78 绑定事实及 ADR-010（2026-11-16 已冻结 EOL）矛盾；门表 9 行零改动系设计裁定，但该残留未列入报告 §10 | 仿 `S1-EXIT-SG2-SYNC`（`0f80c1e`）先例登记后续同步片；至少先在报告 §10 增列该残留与处置 |
| 3 | Low | `docs/validation/release-zip-binding.md:133-140`（`_EN` 同位） | artifacts 块声明的 8 个证据件（`04-reviews/`×3、`06-gates/`×5）在磁盘与变更集中均不存在；「⑨ 前 git add 即解」只覆盖另 10 个已存在未暂存件 | ⑥/⑦/⑨ 按 `MANIFEST.md`「后续追加计划」产出后再声明；G7-b 的披露口径应区分「未暂存」与「未产出」 |
| 4 | Low | `docs/validation/release-zip-binding.md:42` | §4 声称「第一层六份按目标文本落盘」，但 `S1_RELEASE_NOTES:5` 末句、`OPERATIONS:5`、`INSTALLATION_PLAN:32` 与方案书 r1 目标文本存在增补偏离（语义更优且 CN/EN 对称，但非逐字） | 报告中注明三处增补偏离及理由，或将方案书目标文本同步为实际文本 |
| 5 | Low | `docs/INSTALLATION_PLAN.md:7`、`docs/S1_RELEASE_NOTES.md:18`（`_EN:7`、`:18` 同位） | 「当前没有正式发行包 / 当前仍无正式发行包 / No formal release package exists」在绑定完成后语义陈旧；三个 grep 的关键词未覆盖该表述，亦未列入报告 §10 | 后续同步片统一改口径（如「正式发行仍为候选状态」）或列入 §10 残留清单 |

## 逐条证据与推理（摘录）

- 绑定四元组（实读值）：版本 `v0.1.0`；commit `c93571352f456a275d6e71ce3e084c9f96360e86`；run `44043705691`；文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`。三处落点 + `INSTALLATION` + `OPERATIONS` + 报告对逐字一致。
- `INSTALLATION{,_EN}` 的绑定句因 §2 插入两行由 `:206` 后移至 `:208`，属预期漂移。
- 门级：`p57-closeout-assessment.md` §4.2 S-G3 依赖 S-G2 / S-G4 / 发行记录，三者均满足；报告 §7 明确「不构成 S1 Exit 整体通过」。
- CN/EN：numstat 十对文档增删逐对相等；报告对 145 行 ↔ 145 行、`|` 123 ↔ 123、反引号 234 ↔ 234。
- 边界：28 个文件全部位于 `docs/**`；`S1_SUPPORT_MATRIX` 平台表（`:7-15`）与 `:5` 零改动；`ADR_REGISTER:62` 与 `:104` 未改写；`DEV_PLAN` 仅追加注记。
- 复验：`20-reverify.log` 25 条断言 `TOTAL_FAIL=0`，四项均对最终 ZIP。
- 证据包：`MANIFEST.md` 6 条载荷与 `SHA256SUMS.txt` 6 行一一对应。

## 我无法验证的部分

1. GitHub 外部实体当前状态（Release 在线状态、资产完整性、Release notes 正文、CI run 的 8 个 job 结论）——仅以事实基线为准，未做在线复核。
2. Release notes 草稿正文未在输入包内，未逐字核验。
3. `REMAINING_SLICES` P78 段当前为「进行中」，⑧b 完成记录与 R/D 终值尚不存在，属正常时序。
4. `04-reviews/`、`06-gates/` 声明件尚未产出，内容无从核验。
5. `pre-index.txt` 的 G5(d) 红因消除路径需以 ⑨ 实跑为准。

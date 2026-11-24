# P87 ⑤ 预审 r2（收窄复审；`DeepSeek V4 Pro (deepseek)`；2026-11-24；逐字原文）

## Section A
r1 六条 Medium 逐条核：
1. 大上下文文档档阈值 — **闭合**：`DELIVERY_DIRECTIVE.md:573`「审查输入包超 400K token 且为纯文档 / S 档」，EN:573 同位逐字对应。
2. Luna 判据悬空 — **闭合**：`DELIVERY_DIRECTIVE.md:575` 补回平行备选档行（`gpt-6-luna` / Max / 不可用或 §7.5 预算受限平行替换），EN:575 同位。
3. 历史行追溯改名 — **闭合**：CN:90（v1.43 §0.3 行）/ :1158（OB-68 处置列）/ :2265（附录 C P73 规模行）三处还原「§5.2 分档表选定」；EN 三处同位还原「§5.2 tier table」。
4. 备选额度口径 — **闭合**：CN:832「次数上限由切片档决定、模型档只决定模型与单价；Sol 备选沿用所在切片档、轻量档同常规档」，与 §7.5 表体（S/M/L 三行键 = 切片档，v1.35 分档判定规则）一致；EN:832 同位。
5. 三作用域读数 — **闭合**：gate-tree/index/ci 均 `TOTAL_FAIL=0`；selftest `744/744`（740→744 = +4，对应 4 枚新 g5b 夹具 3 合法 + 1 非法）；MANIFEST/SHA256SUMS 已重建，G8-b pack self-consistent PASS。
6. 活文残留 — **闭合**：CN:313 / :1589、EN:313 / :1589 均改「§5.2 的 ⑦ 选型矩阵」/「the §5.2 ⑦ selection matrix」；全库 grep 其余「分档表」仅存于 §0.3/§12.4/附录 C 历史行（应保留）。

## Section B
- B1 Luna 行 vs 四档：**Not a problem**。平行备选档系替换机制而非第五档；常规档行「其余三档判据均不满足」计数自洽；台账「四档」指四个模型档，与权威输入目标①一致。
- B2 生效机制句 vs §8.2：**Not a problem**（无矛盾）。CN:583 的声明义务与 §8.2:919 的 Sol-Max 取档授权义务互补不冲突，§7.5:808「不得中途升档」同向，仅轻微职责重叠。**Low**：权威输入目标⑤尾句「沿用档位启动前设定、不得中途升档」未随句落地（语义由 §7.5:808 / §8.2:919 覆盖）。
- B3 备选额度口径 vs §7.5 表体：**Not a problem**。口径锚定切片档、表体三行键即切片档，「轻量档同常规档」不与任何行矛盾。

## Section C
**成立**。`shape-six.txt` 六形状全 OK（lines 2540 / bold 6160 / pipes 4287 / blank 589 / ticks 6838 / headings 152）。`numstat.txt` 三对逐字对称：准则对 +84/-22（双侧）、ADR 对 +1/-0、台账对 +21/-0；gate G3(a) 同值 PASS。**Low**：fix-summary 称「准则对 +9/-8」与 numstat 实得 +84/-22 措辞不一致（疑指 r2 增量差；对称结论本身成立）。

## Section D
**达标**。gate-tree.txt（scope=tree changed=27）、gate-index.txt（changed=27）、gate-ci.txt（scope=ci changed=2）均 `TOTAL_FAIL=0`；gate-selftest.txt `SELFTEST: PASS 744/744`；G1-b header=v1.72 rows=73 双语一致。build/vet/test-all.txt 为空与③环境声明（纯文档+JS、CI 双腿承担）相符，Not a problem。

## Section E
**Medium ×1**（未闭合，亦非本轮修复引入）：
- 深度档 Neon 触发条件未落：权威输入目标①「（Sol 报 Medium+ 或结构疑点 ⇒ 追加 Neon 深挖）」在 CN:572 深度档行缺失，全库「结构疑点」0 命中、EN 同缺失；适用域排除句 CN:581 的「触发式交叉复核」因此成为**无定义触发**，条款不可判（验收②「字面可判」受损）。
**Low ×4**（不阻断）：① 生效机制句缺目标⑤尾句（见 B2）；② §7.5 分档判定规则（v1.35）「⑦ 三模型分档」计数陈旧（现四档+平行备选行，CN/EN 同态）；③ 轻量档行新增「定稿复核与定向小读」限定（权威输入仅「文档型/治理片」，解读性收窄非冲突）；④ fix-summary「+9/-8」与 numstat 措辞不一致（见 C）。

PRE-REVIEW: FINDINGS

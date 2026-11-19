# P82 架构方案书 r1

## 0 结论摘要

- **通道选择：自愈档（新增子款），非降级档**。理由：① 与「同命令建删特例」同构——该特例已对仓库内同命令建删豁免「次序强制 + 读数②基线」，仓库外版只需平移判据要素，无结构性理由另立第三处置桶；② 自愈档本已把 `$env:TEMP` 列入自删对象域，仓库外误放本在其覆盖意图内，缺口只是同命令建删令次序强制不可满足；③ 从严顾虑由「残留断言落在逐路径 `Test-Path` + `$env:TEMP` 命中数 + 主控亲自执行 + 主控违规须同轮独立确认」补偿，D 计数与 3 次上界照常。**无需改判据代码，不触发 ESCALATE**。
- 修订面：仅 §5.1（新增两处 + 修订一句）、§0.3 一行、附录 C 一块、§12.4 一行；CN/EN 逐句同位。

## Q1 处置通道（CN/EN 成对）

**修订句（CN）——「仓库外位置从严」句由 v1.71 修订，前后对照：**

- 修订前：`**仓库外位置从严**：『同命令建删 + 仓库外位置』组合**无任何自愈路径**，属有意从严。`
- 修订后：`**仓库外位置从严（本句由 v1.71 修订；修订前原文：『同命令建删 + 仓库外位置』组合无任何自愈路径，属有意从严）**：『同命令建删 + 仓库外位置』组合**不适用本特例**、仍属从严；其处置通道为紧随本段之后的『仓库外同命令建删自愈子款』，该子款不减 D 计数、不豁免累计上界，**不得读作本特例的放宽**。`

**修订句（EN）**：
- Before: `**Out-of-repo locations are strict**: the combination "same-command create-and-delete + out-of-repo location" has **no self-heal path at all** and is deliberately strict.`
- After: `**Out-of-repo locations are strict (this sentence revised by v1.71; the pre-revision wording read: the combination "same-command create-and-delete + out-of-repo location" has no self-heal path at all, and is deliberately strict)**: the combination "same-command create-and-delete + out-of-repo location" **does not take this exception** and remains strict; its disposition channel is the "out-of-repo same-command create-and-delete self-heal sub-clause" immediately following this paragraph, which reduces no D count and waives no per-slice cap, and **must not be read as a relaxation of this exception**.`

**新增子款（CN，插在「同命令建删特例」段之后、「非自愈档」段之前）**：

`**仓库外同命令建删自愈子款（v1.71）**：仅当「仓库外位置（`$env:TEMP` / `%TEMP%` / 用户目录 / 其他盘 / 网络共享）+ 同命令建删 + 事后零残留」三者同时成立时适用；仓库外误放若**非同命令建删**，走自愈档常规路径（自报 → 主控记录 → 自删 → 三项读数），本子款不适用。**五项门槛（五项同时成立方可按自愈档取证自愈、不停片）**：① 命令原文（含全部管道、重定向、外部调用、变量展开）粘贴于回执；② 命令原始 stdout 与退出码随附；③ 事后零残留机械证据——自报路径**每一条**的 `Test-Path` 均为 `False`，且固定命令 `Get-ChildItem -LiteralPath $env:TEMP -Filter "*<切片ID>*" -Force` 的命中数 = 0（无论自报路径是否以切片 ID 命名，该读数均须取得并 = 0，作为 `$env:TEMP` 面无残留的独立断言）；`git status --porcelain -uall` **不作为**本子款判定要素（其对仓库外路径恒不显示、示证力为零）；④ 主控核对命令原文：删除动作**之后**不得再连接任何命令 / 语句（`&` / `;` / 管道续段），**「创建 → 删除」本身的连接允许**；且不含 `Remove-Item` 以外的删除工具、`-Force` 之外的隐蔽删除；⑤ 自报须载明违规实际位置与自删动作。**不可自证约束**：③ 的逐路径回读与 `$env:TEMP` 读数**均由主控亲自执行**（命令 + 原始 stdout + 退出码留原文），违规者只提交自报，不得以违规者提供的截图、转述或自制清单代替。**豁免范围**：本子款豁免自愈档「自报 → 主控记录读数 → 自删」的次序强制、**读数②基线**与**预快照完备性**检查（同命令建删下基线物理不可得、命令原文即路径全集来源）；豁免**仅限于本子款成立时**。**完备性替代**：以 ①④ 替代预快照完备性；③ 的路径集合 = 命令原文可推导的创建路径全集；自报路径集合须 ⊇ 该全集，不足即按非自愈档处置。**审查层确认**：违规者为主控本人时，沿自愈档既有机制，本子款判定同样须经**同轮独立审查层确认**后方生效。**计数**：本子款成立的事件**仍计入本片 D 次数**；「不停片」**不免除**「D 累计达 3 次一律停片升级」的义务。**任一要素缺失 ⇒ 按非自愈档处置。**`

**新增子款（EN）**：

`**Out-of-repo same-command create-and-delete self-heal sub-clause (v1.71)**: it applies only when "an out-of-repo location (`$env:TEMP` / `%TEMP%` / the user profile / other drives / network shares) + same-command create-and-delete + zero residue afterwards" all hold together; an out-of-repo misplacement that is **not** same-command create-and-delete takes the self-heal tier's regular path (self-report → master records → self-delete → the three readings), and this sub-clause does not apply. **Five gates (all five must hold for evidence-based self-heal under the self-heal tier without stopping the slice)**: ① the verbatim command (including every pipe, redirection, external call and variable expansion) is pasted into the receipt; ② the command's raw stdout and exit code are attached; ③ mechanical evidence of zero residue afterwards - **every** self-reported path is `Test-Path` `False`, and the fixed command `Get-ChildItem -LiteralPath $env:TEMP -Filter "*<sliceID>*" -Force` yields a hit count = 0 (the reading is taken and must equal 0 whether or not any self-reported path is named by the slice ID, as an independent assertion that the `$env:TEMP` face is residue-free); `git status --porcelain -uall` is **not** a judging element of this sub-clause (it never shows out-of-repo paths, zero probative value); ④ the master checks the verbatim command: **after** the deletion action no further command or statement may be chained (`&` / `;` / a pipe continuation), **chaining of the "create → delete" itself is allowed**; and it contains no deletion tool other than `Remove-Item` and no hidden deletion beyond `-Force`; ⑤ the self-report must state the actual violated location and the self-deletion action. **No-self-attestation constraint**: the per-path readback and the `$env:TEMP` reading of ③ are **all executed by the master personally** (command + raw stdout + exit code kept verbatim); the violator submits only the self-report, and screenshots, paraphrases or self-made lists provided by the violator must never stand in for it. **Scope of the waiver**: this sub-clause waives the self-heal tier's ordering mandate "self-report → master records the reading → self-delete", **reading ②'s baseline** and the **pre-snapshot completeness** check (under same-command create-and-delete the baseline is physically unobtainable, and the verbatim command is the sole source of the path set); the waiver is **confined to "when this sub-clause holds"**. **Completeness substitute**: ①④ take the place of the pre-snapshot completeness check; the path set of ③ = the full set of created paths derivable from the verbatim command; the self-reported path set must be a superset of that full set, any shortfall ⇒ the non-self-heal tier. **Review-layer confirmation**: when the violator is the master itself, following the existing self-heal-tier mechanism, this sub-clause's ruling likewise takes effect only after the **same-turn independent review layer confirms it**. **Counting**: an event where this sub-clause holds **still counts toward this slice's D count**; "no stopping" **does not waive** the duty "D cumulative reaching 3 always halts and escalates". **Any missing element ⇒ the non-self-heal tier.**`

## Q2 `tmp/` 根自愈口径明确化

**CN**：修订前 = 自愈档「预快照与完备性」句原文（不动）；修订后 = 原句不动，其下**新增一句**：
`**`tmp/` 根自愈的逐要素口径（v1.71 明确化；本档三项读数文本零改动）**：五要素固定次序——**自报**（载明逐路径实际位置与自删动作）→ **主控取改前读数**（固定命令 `(Get-ChildItem -LiteralPath tmp -Recurse -File -Force | Measure-Object).Count` 取 `tmp/` 递归文件数，与 `??` 清单、`$env:TEMP` 命中清单同批留原文）→ **自删**（违规者逐件精确自删）→ **主控取改后读数**（同一条固定命令复取）→ **零残留判定**（改后读数**不大于**改前读数，且自报路径逐条 `Test-Path` 均为 `False`，且 `??` 项不含自报路径；判定依据仍是本档三项读数，本句只固化命令与时序，不新增读数）。`

**EN**：Before = the existing "Pre-snapshot and completeness" sentence (untouched); After = same sentence plus, appended below it:
`**Per-element caliber for `tmp/`-root self-heal (v1.71 clarification; this tier's three readings are textually untouched)**: five elements in fixed order - **self-report** (states each path's actual location and the self-deletion action) → **master takes the pre-deletion reading** (fixed command `(Get-ChildItem -LiteralPath tmp -Recurse -File -Force | Measure-Object).Count` for the `tmp/` recursive file count, recorded verbatim together with the `??` list and the `$env:TEMP` hit list) → **self-delete** (the violator deletes its misplaced artifacts item by item) → **master takes the post-deletion reading** (re-runs the same fixed command) → **zero-residue ruling** (the post-deletion reading is **no greater than** the pre-deletion reading, every self-reported path is `Test-Path` `False`, and the `??` list contains no self-reported path; the ruling basis remains this tier's three readings, this sentence only fixes the commands and the order).`

## Q3 既有条款零扰动证明（可机械核验的断言）

1. **三项读数**：实施后 `git diff -U0 docs/DELIVERY_DIRECTIVE.md` 在「零残留机械证据三项」所在段无任何 +/− 行；三个序数锚点（`① 本次误放件` / `② tmp/ 递归文件数` / `③ $env:TEMP`）各自在文件中仍**恰好出现一次**且所在行未变；改前整段 sha256 与改后逐字节相同。
2. **D 上界 3 次**：`累计上界（v1.45）` 段在 diff 中无改动；锚点串「D 次数累计达 3 次」与「含同命令建删特例的取证自愈」改前后逐字一致。
3. **防再犯句**：在起跑消息侧（v64 快照 §五 内），**不在准则正文**——本片不动其文本；v64 §五 只改计数行与 OB-139 条目状态（见 Q5）；机械断言 = v64 的 diff 中该句所在行无变化。
4. **§6.3 与判据代码**：变更集不含 §6.3 区间与 `tools/gates/**`——`git status --porcelain` 清单无该目录任何路径；`git diff -U0` 在 §6.3 区间无 hunk。判定结论：**零代码改动**。

## Q4 v1.71 版本与日志

- MINOR 依据（§0.2 原文）：`MAJOR` = 角色集或流水线结构的变更；`MINOR` = 规则增补。本片为规则增补（新增子款 + 口径句），不触角色集与流水线结构 ⇒ **MINOR**。头部版本行：`版本：v1.71。日期：2026-11-19。`
- §0.3 新行（完整单元格）：`| v1.71 | 2026-11-19 | **切片 \`P82 TMP-DISCIPLINE-OUT-OF-REPO\` 落地（治理片 \`[SLICE]\`；M 档；用户 2026-11-19 裁定立项）**：① §5.1 新增「仓库外同命令建删自愈子款」——为「仓库外位置 + 同命令建删」补自愈通道，五要素含主控亲自取证的逐路径 \`Test-Path\` 与 \`$env:TEMP\` 切片 ID 命中数断言，次序强制与读数②基线豁免仅限本子款；「仓库外位置从严」句修订为指向本子款（修订前原文就地留档）；② §5.1 自愈档补「\`tmp/\` 根自愈逐要素口径」句（固定命令与时序固化，三项读数文本零改动）；③ §12.4 的 OB-139 处置列翻 \`已处置\`；④ 附录 C 增 C.0ax。**不改自愈档三项读数文本、D 上界、同命令建删特例五项门槛、§6.3 判据与 \`tools/gates/**\`、CONTRACTS / schema / fixtures / CI。 | \`P82 TMP-DISCIPLINE-OUT-OF-REPO\`（本片） |`（EN 镜像逐句同位，表头列取值照搬。）
- 附录 C 新块 **C.0ax**（照 C.0aw 体例，CN/EN 对称）：

| 项 | 内容 |
|---|---|
| 切片 / 规模 | `P82 TMP-DISCIPLINE-OUT-OF-REPO`（治理类 `[SLICE]`；M 档；探针 0；Windows 本机） |
| 触发 | 用户 2026-11-19 裁定：OB-139 第 6 次触发 ⇒ 立项 P82（修 §5.1 临时件纪律，升 v1.71，起跑须用户点名） |
| 目标 | §5.1 补「仓库外同命令建删自愈子款」与「`tmp/` 根自愈逐要素口径」；准则升 v1.71；OB-139 翻 `已处置` |
| 交付物 | 准则对 `docs/DELIVERY_DIRECTIVE{,_EN}.md`；台账 P82 段；报告对 `docs/validation/tmp-discipline-out-of-repo{,_EN}.md`；证据包 `docs/validation/evidence/P82-TMP-DISCIPLINE-OUT-OF-REPO/` |
| 范围 | 不改自愈档三项读数文本、D 上界、同命令建删特例五项门槛、§6.3 判据本体与 `tools/gates/**`、CONTRACTS / schema / fixtures / 产品代码 / CI / p57 / `s1-exit` / v63 快照；不新增 ADR / OB |
| 审查 | ⑤ 预审、⑥ 独立扫描、⑦ 终审（治理片额度）；读数以本片报告为准 |
| 版本 | v1.71（头部版本行与 §0.3 最新行一致，G1-b） |

## Q5 OB-139 收口

- §12.4 该行处置列目标文本（现象列**逐字不动**，登记日不动）：`**已处置**（P82 \`TMP-DISCIPLINE-OUT-OF-REPO\` 落地，v1.71）：§5.1 新增「仓库外同命令建删自愈子款」，为「\`%TEMP%\` 仓库外位置 + 同命令建删」补处置通道，并固化 \`tmp/\` 根自愈逐要素口径；D 计数口径与累计上界不变；本处置不宽免（P80 一次性宽免不作先例，本行不引用）。`
- EN 行：`**Resolved** (landed by P82 \`TMP-DISCIPLINE-OUT-OF-REPO\`, v1.71): §5.1 adds the "out-of-repo same-command create-and-delete self-heal sub-clause", giving the "\`%TEMP%\` out-of-repo location + same-command create-and-delete" combination a disposition channel, and fixes the per-element caliber of \`tmp/\`-root self-heal; the D-count caliber and the per-slice cap are unchanged; this disposition grants no waiver (the P80 one-off waiver is no precedent and is not cited here).`
- v64 快照 §五 同步点：待议 10 → **9**（OB-139 移出待议清单）；观察中保持 **13**；非终态 24 → **23**；OB-139 状态条目由「已转待议 + 立项 P82」改写为「**已处置**（P82 落地，v1.71）」；防再犯句文本不动。台账 `REMAINING_SLICES{,_EN}.md` 的 P82 段同步为已完成（状态行 ✅ + 完成记录行 = 报告路径，按 §11.2 完成行格式）。

## Q6 文档同步三层归属

- **第一层（本片必改）**：准则对 §5.1 / §0.3 / 附录 C.0ax / §12.4 的 OB-139 行（CN+EN 逐句同位）；台账 `REMAINING_SLICES{,_EN}.md` P82 段；v64 快照 §五（计数行 + OB-139 条目状态）；报告对 `docs/validation/tmp-discipline-out-of-repo{,_EN}.md`；证据包。
- **第二层（需实读确认后决定）**：v64 快照的 `_EN` 镜像同步；v64 §二 附段 A「义务传递」段对 P82 的时态措辞（若为将来式须翻完成式）；v64 §六「下一裁定项（预判）」行是否删除（P82 已起跑落实）；`DEV_PLAN` 是否已有 P82 引用行（按 §11.2 治理片不回写 DEV_PLAN，仅核对既有引用）。
- **第三层（明确不动）**：v50–v63 全部快照；P81 报告正文（§3 的 D=2 记录不追改，见 Q7）；CONTRACTS / schema / fixtures；`tools/gates/**` 与 §6.3；产品代码与 CI；p57 / `s1-exit`；自愈档三项读数、D 上界、同命令建删特例五项门槛；§12.4 其他行与附录 C 既有块；不新增 ADR / OB。

## Q7 历史行判定

- §0.3 历史行（v1.70 及之前）、§12.4 其他行、附录 C 既有块：按 §11.1 **不追溯**——v1.70 行内「升 v1.71，待排期，起跑须用户点名」属历史时点记录，不改写。
- P81 报告 §3 的「临时件 D = 2」两条事件记录**不追改**；本片报告仅以**落点溯源指针**引用其所在小节与事件编号（`docs/validation/directive-model-upgrade.md` §3），不复制其正文、不修正其定性。
- 本片新增内容均落在新增行/新增块，符合「只判变更行、历史行不追溯」的判据作用域口径。

DESIGN: DONE
# 05 — ⑤ 预审第 1 轮原文（`DeepSeek V4 Pro`，2026-11-26）

# PRE-REVIEW: PASS WITH FIXES

**审查对象**：暂存区 5 文件变更（`DELIVERY_DIRECTIVE{,_EN}.md` +19/−2、`REMAINING_SLICES{,_EN}.md` +22/−0）+ P91 证据包 11 件。已实读 `ci.yml`、K4、§0.2/§0.3/§6.3/§7.2.1/§4.3、附录 C.0az、P90 台账段、证据包全部载荷与两轮门禁读数，并独立重算证据包全文件 SHA256。

## 必查清单逐条判定

1. **判据语义影响**：**无**。grep 实证 `tools/gates/lib/criteria/*` 对「脚本化状态」句零读取路径——判据对准则的读取限于 G1-b（头部版本行）、G3(a.2)（全文件形状）、G4b（语料字）、G5-c（附录 B 调用串）、G6 系列（§0.3/§12.4/附录 C 区新增行元数据）。本片改动仅为新增行落入既有判据作用域，判据本体与判定逻辑零改动。
2. **事实正确性**：**实**。实读 `ci.yml`：`Format` 步在 95 行、`Gates` 步在 107 行，`Gates` run 串为 `--selftest && …--all --skip-base --scope=ci`，与准则新句逐字一致；`Build`/`Vet`/`Test`/`Format` 均在 `Gates` 之前（其间 `Race`/`Contract fixtures` 非 base 四条，不影响表述）。K4 拟改文本与实落逐字一致，且未偏离「K4 为唯一出处」纪律。
3. **CN/EN 对位**：**未破**。步骤 2 tree 读数实证：两指令对 numstat 全等（+19/−2、+22/−0），G3(a.2) 六形状 `lines/head/bold/bars/blank/backticks` 全 true；括注位置与反引号范围双侧对应。
4. **版本包三件**：**齐全且合规**。G1-b 实测 PASS（CN/EN 均 header=v1.74、last=v1.74、日期一致）；§0.3 新行末列「`P91 DIRECTIVE-CI-INVOCATION-SYNC`（本片）」体例正确；附录 C 新块 `C.0ba` 位于 `C.0az` 之后、`C.1` 之前，字段集与 C.0az 同构（8 项全对齐），编号进位合理。
5. **OB-145 行**：现象列已说清「v1.40 带日期快照句」性质；处置列四要素齐全；首状态词「已处置」∈ G6-6 闭集（step2 读数 scanned 2 PASS）。**但有超前措辞**（见 F-2）。
6. **台账 P91 段**：17 字段与起跑消息 §1.2 字段表逐项对应，与 P90 段体例一致（含「现状缺口」保留起跑时态的先例写法）；无重复实体新增（见 F-4 注记）。
7. **⑥ 可跳过（独立判定）**：**可跳过**。§7.2.1 三条同时满足：① 变更集仅 `.md`/`.txt`；② 无 `.go`/`.json`/`.yml`/schema/fixtures/脚本逻辑；③ 不涉状态机/判据/证据模型/角色授权语义——判据实现无该句读取路径（第 1 条实证）。注意 SW-1/SW-2 反面样例不适用：本片未改任何 G 子规则条文。**前提**：跳过须 ⑦ 或用户确认，不得主控单方认定。
8. **门禁风险**：step2 tree 读数实证：G6-4（scanned 82 无命中）、G6-3（无非法 run 号形态）、G4b（newCjk=313、unknown=[]）、G1-a、G2 全 PASS；G7-e INFO（变更集无报告文件，无目标）；G6-5 refs 18 全解析、G6-7 declarations 1 对账通过、G6-6 scanned 2 PASS。**实际命中仅 G8-b**（见 F-1，已知中间态）。
9. **证据包三方相等**：已登记 6 项哈希与磁盘**全部一致**；但条目集不等（见 F-1）。

## 发现

**F-1 / Low / 证据包条目集不等（G8-b 当前 FAIL）**：磁盘载荷实为 8 件，MANIFEST 与 `SHA256SUMS.txt` 各只列 6 件；未登记 3 件 = `01-design.md`、`02-diff-before-after.txt`、`03-verification/gate-tree-step2.txt`（step2 读数原文实录 `G8-b pack self-consistency :: FAIL`）。已登记条目 SHA256 全部相符，属收敛阶段重建项。建议：⑨ 前按磁盘派生重建 MANIFEST/SHA256SUMS 并重跑 `--check=G8`。

**F-2 / Low / OB-145 与 §0.3 措辞超前**：OB-145 处置列「**本片 P91 落地（v1.74）**：准则 :675 句改为…」与 §0.3 v1.74 行「切片 `P91 …` **落地**…」均为完成式，而台账 P91 段状态为「🔄 执行中」（未过 ⑤⑥⑦⑧a⑨）。G6-1 不判此形态（step2 PASS），但属描述性超前。建议：改未完成式「由本片 P91 处置（v1.74）」，⑧a/⑨ 闭环后再回写「落地」；EN 侧同步。

**F-3 / Low / input-list 引用失准**：第 9 行「见 `03-verification/selftest.txt`」实际文件名为 `selftest-step0.txt`（缺 `-step0`）；且「期望（起跑消息）」列含起跑消息原文未出现的「准并行数 2557 / 2557」锚点，期望列来源不实。建议：修正引用名并标注该锚点来源（或移出「起跑消息」列）。

**F-4 / 注记 / 同句两处引述**：拟改句「再 `--all --skip-base --scope=ci`（base 由前序 Build/Vet/Test/Format 步覆盖，P90）」在 OB-145 处置列与 §0.3 v1.74 行各引述一遍，当前逐字一致（与 OB-142/143/144 处置列复述体例同款）；漂移风险由 G3 形状 + ⑤⑥⑦ 人工承担，无机械判据覆盖，后续改句须两处同步。

## 结论

- 边界合规（无 `ci.yml`/`tools/gates/**`/`internal/**`/CONTRACTS/schema/fixtures/新 ADR 越界）；事实正确；G3/G4b/G6 全绿；仅 3 项 Low（F-1 为已知收敛项、F-2/F-3 为描述性修正）+ 1 注记。
- ⑥ 可跳过（须 ⑦/用户确认）；建议 ⑦ 复核 F-2 措辞收敛与 F-1 重建时点。

verdict: no Medium+

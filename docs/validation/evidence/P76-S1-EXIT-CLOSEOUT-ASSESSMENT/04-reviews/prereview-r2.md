# P76 ⑤ 预审 r2（retryOf = r1）

审查者：deep-reasoner + DeepSeek V4 Pro（独立于实现方）
日期：2026-11-16。审查范围：仅 r1 的 6 项发现闭合判定 + 整改引入的新问题 + H-1 处置可接受性；不重审全片。

## 输入（实际读取的文件）

1. `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/05-prereview-r1.md`（r1 原文，全文实读）
2. `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/b-diff-r2.patch`（整改后完整暂存 diff，全文实读 + 分文件 hunk 提取）
3. `docs/validation/s1-exit-closeout-assessment.md` / `_EN.md`（本片报告，磁盘全文实读）
4. `docs/validation/s1-exit.md` / `_EN.md`（`:52` 及被引用行，磁盘实读）
5. `docs/validation/p57-closeout-assessment.md` / `_EN.md`（`:101` 指针，磁盘实读）
6. `docs/ADR_REGISTER.md` / `_EN.md`（ADR-046 / ADR-047 两行，磁盘实读）
7. `docs/SECURITY.md` / `_EN.md`、`docs/PRODUCT_REQUIREMENTS.md` / `_EN.md`（注记行，磁盘实读）
8. `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/` 下 `freedom-list.md`、`SHA256SUMS.txt`、`01-arch/design-r1.md` §3、`02-ruling/ruling-2026-11-16.md`（全文）
9. `docs/DEV_PLAN.md`、`docs/INSTALLATION.md`、`docs/INSTALLATION_PLAN.md`、`docs/validation/s1-reach-reachability.md`、`docs/t027/CHECKLIST_v59.md`（被引用行定点复核）
10. 终端实跑：`git diff --cached --numstat`、BOM/LF 检查、SHA256 复算、报告对六形状统计

## 一、r1 发现逐项闭合判定

| 编号 | 判定（CLOSED / NOT CLOSED） | 依据 |
|---|---|---|
| H-1 | CLOSED（处置方式可接受，详见下节 H-1 判定） | `:52` 的 S-G2 句按边界保留；残余在报告 §4 / §17、p57 指针 `:101`、ADR-046 影响面、freedom-list §四共五处透明登记；ADR-046 硬边界保持安全不变量 |
| L-1 | CLOSED | 报告 §1 改为「落地六项裁定（Q2 / Q3 / Q7 / Q9 / Q10 与收口顺序）…」；EN 同位同义 |
| L-2 | CLOSED | 报告 §9 改为「`docs/validation/` 全域无独立安全评审报告」；EN「the absence of any independent security review report」同位 |
| L-3 | CLOSED | SECURITY{,_EN}:7 与 PRODUCT_REQUIREMENTS{,_EN}（CN `:111` / EN `:109`）注记均补「且不晚于 2027-03-31」 |
| L-4 | CLOSED | freedom-list §四已定稿：F-1..F-5 + 新增 F-6/F-7/F-8 + 未行使项；「待回填 / 待定」清零；SHA256SUMS.txt 重算且 8/8 与磁盘哈希一致（本环复算） |
| L-5 | CLOSED（附新发现 L-6） | 报告 §5 新增方案 A/B 文本对照；ADR-046 依据补「+ 附加项」（EN 同位）。方案 B 文本可回溯性见 L-6 |

### H-1 判定（处置方式是否可接受）

**磁盘实证**：`s1-exit{,_EN}.md:52` 的改动仅限 AT-23 片段（`git diff --cached --numstat` = `1 1`；「AT-23 未通过」→「AT-23 整体未通过」/「AT-23 as a whole is not passed」，恰为方案 A 的「整体」2 字）。S-G2 句「S-G2 门级判定待 v54 快照，不得据此声称该门通过」在磁盘上逐字保留，零改动。

**登记点核实（五处，均实读确认）**：

1. 报告 §4「已登记残余」：明写「属本片边界外的历史表述，本片未改该句；其同步修正留后续切片」；
2. 报告 §17 已知边界：「S-G2 句未同步（本片仅改 AT-23 片段），已登记为残余」；
3. p57 `:101` 指针括注：「（`s1-exit.md:52` 的 S-G2 句尚未同步，属已登记残余，留后续切片。）」；
4. ADR-046 影响面：「`s1-exit.md:52` 的 S-G2 句留后续切片同步（本片不改该片段）」（EN 同位）；
5. freedom-list §四未行使项：「不改 `s1-exit.md:52` 的 S-G2 句（本片边界；已登记为报告 §4 / §17 的残余）」。

**可接受性判断**：

- 交接消息【三】改动 1 的边界只授权 `:52` 行内 AT-23 片段、其余行零改动；本片执行与边界逐字相符（numstat 1/1）。若本轮越界改 S-G2 句，将违反用户授权——r1 建议的两条路径中「补齐」被本轮边界禁止。
- 「经用户裁定保留」路径由边界本身表达：授权范围只覆盖 AT-23 片段，S-G2 句同步即被本轮边界排除在外。故「透明登记残余 + ADR-046 影响面声明留后续切片」是本轮约束下最优处置。
- 安全不变量保持：报告未宣告 S1 Exit 通过（S-G7 不可达、S-G1 不通过、S-G8 blocked）；ADR-046 硬边界「`s1-exit.md:52` 未修正前不得宣告」成立。
- r1 所报的静默冲突（报告声称与磁盘不符）已消除：§4/§17 与磁盘现状一致，§14 第 2 条「本片已成文」现指门级判定已在报告与 ADR-046 成文，且被 §4/§17 限定。

**遗留条件（不阻塞本片，但须追账）**：后续切片必须执行 `:52` 的 S-G2 句同步；该残余的承载切片尚未具名（见 L-7）。

## 二、最低核验动作（五条逐条结果）

### ① 报告声明的 `file:line` 与磁盘一致 —— 通过

| 报告引用 | 磁盘实读结果 |
|---|---|
| `s1-exit.md:43` 不构成真实候选兼容性证明 | ✓ 逐字一致 |
| `s1-exit.md:44` 条件闭合 | ✓ |
| `s1-exit.md:52`（修文对象） | ✓ :52 即核心任务链行，AT-23 片段已改 |
| `s1-exit.md:62` P71 注记 / `:66` EOL | ✓ |
| `DEV_PLAN.md:119` T028 `COMPLETED` / `:94` 依赖字段 / `:96` 状态 / `:27` 独立安全评审 | ✓ 逐字一致 |
| `PRODUCT_REQUIREMENTS.md:69/:71/:82/:107` | ✓ 逐字一致（注记插入于 :111，不漂移以上行） |
| `INSTALLATION.md:206` 三处绑定 | ✓ |
| `INSTALLATION_PLAN.md:32` EOL | ✓ |
| `SECURITY.md:5`「未完成独立安全签署」 | ✓ |
| `p57:88`–`:95` 八门行 | ✓ 八行齐全 |
| `s1-reach:22`–`:31` | ✓ G-A 起十门表在列 |
| `ADR_REGISTER` ADR-046 = `:57`、ADR-047 = `:58`（EN `:58/:59`） | ✓ |
| p57 P76 指针 = `:101`（CN/EN 同位） | ✓ |
| SECURITY 注记 = `:7`（CN/EN）、PR 注记 = CN `:111` / EN `:109` | ✓ |

### ② diff 与报告描述一致 —— 通过

- `s1-exit{,_EN}.md`：numstat `1 1`，唯一替换 = AT-23 片段加「整体」/「as a whole」2 字（方案 A），与报告 §5 描述逐字一致；其余行零改动 ✓。
- p57{,_EN}：`+2` 行加法式指针，门表八行零改动 ✓。
- ADR_REGISTER{,_EN}：`+2` 行（ADR-046 / ADR-047）✓。
- SECURITY / PRODUCT_REQUIREMENTS{,_EN}：`+2` 行加法式注记，含 2027-03-31 ✓。
- S1_SUPPORT_MATRIX{,_EN}：`2 2`（仅 `:13`/`:14` 句尾）✓。
- 新报告对：209 行新建 ✓；证据包 8 件新建，SHA256SUMS 条目集与文件集一致 ✓。
- 报告 §1/§5/§9/§17 的措辞与 diff 全部相符（CN/EN 双份）✓。

### ③ CN/EN 逐字同位、语义等价 + 报告对六形状 —— 通过

| 形状 | CN | EN | 结论 |
|---|---|---|---|
| 行数 | 210 | 210 | ✓ 等长 |
| 标题数 / 位置 | 20 处 @ 9,16,28,44,52,64,70,81,87,94,100,109,115,126,138,153,157,163,169,182 | 同 CN 全同位 | ✓ |
| `**` 数 | 2 | 2 | ✓ |
| `|` 数 | 202 | 202 | ✓ |
| 空行数 | 50 | 50 | ✓ |
| 反引号数 | 220 | 220 | ✓ |

各修文 CN/EN 语义等价逐条核对：§1 六项裁定表述、§4 已登记残余、§5 方案 A/B 对照、§9 全域无独立安全评审报告、§17 残余边界条、ADR-046/047 两行、p57 指针、SECURITY/PR 注记、s1-exit `:52` 片段（`整体` = `as a whole`）、S1_SUPPORT_MATRIX `:13/:14` ✓。

### ④ 边界声明与 diff 相符 —— 通过

- diff 触及 24 个文件，全部落在起跑消息【八】第一层（s1-exit 修 `:52`、报告对新建、证据包追加）与第二层（p57 / S1_SUPPORT_MATRIX / PRODUCT_REQUIREMENTS / ADR_REGISTER / SECURITY）；第三层（准则 / CONTRACTS / CHECKLIST_v59 / v50–v58 / s1-reach / 代码与 CI）零触及 ✓。
- `s1-exit{,_EN}.md` 其余行零改动（numstat 1/1）✓；状态列 `BLOCKED`、门表行名与行数不变 ✓。
- 交接消息【三】改动 1 的边界（仅 AT-23 片段、方案 A 加 2 字）逐字相符 ✓。
- `REMAINING_SLICES` 未触（⑧b 待做，符合流程）✓。

### ⑤ 机械读数可复算 —— 通过

| 项 | 报告 §2 值 | 本环复算 |
|---|---|---|
| ADR 最大号 | 47 | 磁盘 CN/EN 均 = 47（ADR-047 为最大号）✓ |
| OB 最大号 = 139 / DR 最大号 = 14 / 试点顺延第 25 次 | 与 r1 相同 | `CHECKLIST_v59.md:47` / `:145` 逐字核对 ✓（该文件本片未触） |
| SELFTEST | 737/737 | `tmp/.../gate-selftest-c.txt` 尾部 = `SELFTEST: PASS 737/737` ✓ |
| SHA256SUMS | 8 项 | 对 8 个证据包文件逐一 sha256 复算，全部 MATCH（含已定稿 freedom-list）✓ |
| 编码/EOL | （未声明） | 本片改动全部 `.md` 为 UTF-8 with BOM + LF；`SHA256SUMS.txt` 无 BOM + LF ✓ |

## 三、新发现（按严重度）

- **L-6（Low）**：报告 §5 方案 B 文本可回溯性不足。报告写「方案 B = `AT-23 整体未通过（阶段 1 条件闭合已达成，详见本行依赖链说明）`」，但与磁盘可查的两个权威源均不一致：design-r1 §3 建议文本 = `T027 未完成、AT-23 未通过（阶段 1 条件闭合已达成、阶段 2 不适用，2026-11-09 裁定 A1）；`，ruling §八只提方案 A/B 之名未给文本；且「详见本行依赖链说明」在证据包任何文件中均不存在。若该文本出自交接消息【三】改动 1，建议报告注明出处行；否则改为逐字引用权威源文本，避免「供回溯」的对照本身不可回溯。
- **L-7（Low）**：H-1 残余的承载切片未具名。报告 §18「下一步」三条（ADR-010 批准片 / Release 与 ZIP 绑定复验片 / 独立安全评审补评片）均不覆盖「`s1-exit.md:52` 的 S-G2 句同步」；p57 指针与 ADR-046 只写「留后续切片」。建议在 ⑧b 的 `REMAINING_SLICES` 回写中为该残余具名承载切片（或并入既有切片），确保可追踪闭合、不被遗忘。

（r2 未发现 Critical / High / Medium。以上两条均为 Low，不阻塞本片收口。）

## 四、结论

- r1 六项发现：**全部 CLOSED**。H-1 的「登记残余」处置在本轮边界（仅授权 AT-23 片段）下**可接受**：零越界（numstat 1/1）、安全不变量保持（未宣告 S1 Exit 通过、ADR-046 硬边界成立）、冲突五处透明登记不再静默；同步义务明确留后续切片（附 L-7 追账建议）。
- 新发现：Critical 0 / High 0 / Medium 0 / Low 2（L-6、L-7）。
- 建议：L-6 由 ②c 附带修正（出处注明或逐字引用）；L-7 由 ⑧b 的 `REMAINING_SLICES` 回写一并处置。两者均为 Low，不需再跑 ⑤（r2 已达硬上限 2 次）。

PRE-REVIEW: PASS WITH FIXES

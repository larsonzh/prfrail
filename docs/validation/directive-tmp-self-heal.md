# 证轨 P34 治理片验证报告：临时件自愈特例（准则 v1.56）

切片 `DIRECTIVE-TMP-SELF-HEAL`；基线 `HEAD` = `53ecf7f`；准则版本 v1.56；规模 M（治理片）；报告日期 2026-10-28。

## 1 变更概要（含异常与兜底记录）

**问题**：§5.1 现行「自删先于主控读数 ⇒ 非自愈档」使「同命令建删」在物理上无法自愈；P32（D-1 / D-2 / D-3）与 P33（D-1）已连续三次踩到该缺口。

**裁定**：用户 2026-10-28 分三组逐组裁定，选「在 §5.1 自愈档下新增同命令建删特例」；① 只出表、不表态，出表后停等裁定。

**结果**：准则升 v1.56（头部 / §0.3 / §5.1 新增段与两处最小闭合 / §9.3 处置格扩界 / 附录 C.0aq），台账 P34 块与边界行同步；三作用域门禁 `TOTAL_FAIL=0`，自测 `SELFTEST: PASS 697/697`；未改代码与判据本体。

| 文件 | 落点 | 内容 |
| --- | --- | --- |
| 两准则 | 头部 `:3` | `版本：v1.56。日期：2026-10-28。` |
| 两准则 | §0.3 `:102` | v1.56 变更行 |
| 两准则 | §5.1 `:516` | 新增「同命令建删特例（v1.56）」段 |
| 两准则 | §5.1 `:518` 与 `:520` | 非自愈档、累计上界两处最小闭合 |
| 两准则 | §9.3 `:902` | 处置格双版本标记（扩界） |
| 两准则 | 附录 C `:2080` | `C.0aq` 切片块 |
| 两台账 | 状态行与 P34 块 | 状态行同步、边界行追加扩界 |

**第 1 项 R = 2**：R1 = ② 首轮把任务包的改写指令 `W3`–`W9` 照抄进准则正文（CN 7 处、EN 7 处）；发现方式 = 主控侧机械核验（不采信自报）；修复 = 派 ② 返工，转写为自包含标签；根因 = 任务包未明写「`W1`–`W13` 是改写指令、不是正文内容」，执行侧亦未判断；`G6-5` / `G6-7` 抓不到（只判 §X.Y 与 OB / SW / DR 类 ID）。R2 = ⑤ r1 两条 Low（台账边界行 `:899` 行号已漂移到 `:902`、`OB-128` 触控面漏列 §9.3）⇒ 定点整改后 ④ 复跑。

**第 2 项 D = 1**：② 返工轮在 `tmp/DIRECTIVE-TMP-SELF-HEAL/` 内建 3 支一次性脚本并自行删除；定性 = **非自愈档**（自删先于主控读数，读数不可得）；处置 = 用户以 `OB-37` 元规则裁定继续（不停片），**不得作先例**。

**第 3 项 元层面发现**：本片正在落的 v1.56 特例明确排除「位置合法而仅次序违规」的情形，**首例由本片自身触发，恰落新规则排除侧**——不是执行侧省力，是 v1.56 排除条件过严的首次实证。

**第 4 项 证据缺口（纪律 B）**：丢失 = 3 支被删脚本（`p34-rework-relabel.js` / `p34-rework-postcheck.js` / `p34-rework-scope.js`）的正文与 stdout；未丢 = 两 `.bak` 前像、`apply.js` / `verify.js` 全文、改动后六形状与门禁读数、`git diff` 全量。**未主张证据缺口已消除。**

**第 5 项 §12.1 引用精度**：起跑消息把「治理片不满足试点选择标准」记为 §12.1 本体表述，实为 §12.4 各 OB 行的既定写法（§12.1 本体无该字面句）；本片不修 §12.1 本体（超边界），仅登记事实，供 v23 清单改用「§12.4 OB 行既定写法」或明写引用链。

**第 6 项 ③ 的 `TEST: FAIL (1)` 留痕**：FAIL 项 = 任务包自定的「CN / EN 句子数相等」判据；误报依据 = CN 用「，」连写、EN 用「;」断句，语义逐项对应；**未重跑 ③ 直接进 ④ 的合规性** = §5.3 对 ③ 的重跑触发条件是「修代码 / 测试」，本片既未修代码也未修测试，触发条件不成立，③ 的实质检查项已全部通过。

**第 7 项 批量辅助**：本片曾遗漏一次安排（① 之后、裁定之前），性质 = **执行安排错误**（非「认为冗余」）；已补派 3 次（枚举 / 比对 / 影响面）；登记为 **v22-OBS-C 的第二实例**。

**第 8 项 `W1`–`W13` 采纳与 X1 扩界**：逐条采纳清单见 1.1；扩界 = §9.3「② 提前清理」处置格与台账边界行同步，见 1.1 末段。

**第 9 项 边界守护正面例证**：本片发现「完成表缺 P34 条目」时，因边界严令（扩界前规定“除三处外不得改其他行”）而拒绝擅自补加、如实上报——与同片刚落地的「同命令建删特例」在精神上一致（守边界、不夹带）。

### 1.1 `W1`–`W13` 采纳清单与 X1 扩界

| 编号 | 采纳内容 |
| --- | --- |
| `W1` | ④ 精确化为「删除动作之后不得再连接」 |
| `W2` | ⑤ 理由句改为「`git status` 恒不显示 + §5.1 零容忍」 |
| `W3` | 显式豁免次序强制与读数②基线，范围限定为「特例成立时」 |
| `W4` | 以 ①④ 替代预快照完备性检查，③ 路径集合 = 创建路径全集，③ 沿既有机制由主控执行 |
| `W5` | ①② 属「不可重放命令」证据例外 |
| `W6` | 违规者为主控本人时须经同轮独立审查层确认 |
| `W7` | 不改变三路径次序、不免除 ⑩ 前终态回读 |
| `W8` | ⑧c 后事件亦不豁免「不得为整改动作重建」禁令 |
| `W9` | 「同命令建删 + 仓库外」无自愈路径，有意从严 |
| `W10` | §5.1 两处最小闭合：触发枚举补「特例五项任一缺失」、上界补「含特例取证自愈」 |
| `W11` | `C.0aq` 首笔口径核准为「两次裁定、三次事件」 |
| `W12` | `C.0aq` 措辞取「事实真值源仍以两片报告为准，附录仅作汇总留痕」 |
| `W13` | 原四项追加第 ⑤ 项「建删位置为仓库内路径」 |

**X1 扩界**：§9.3「② 提前清理」处置格内「自愈档」一处扩为「自愈档（v1.45）或同命令建删特例（v1.56）」，保留两个版本标记、不动症状列与其余格；台账 P34 边界行同步追加「+ §9.3 处置格（扩界）」。

**第 9 项 v23 观察项候选 3 条**：

1. 任务包自定判据的设计——本片 ③ 因「CN / EN 句子数相等」判据不严谨产生 1 条误报。
2. v1.56 的排除条件是否过严——位置合法 + 同命令建删应否纳入特例（触发 = 本片 D = 1 首次实证）。
3. 批量辅助安排的时点——本片遗漏一次。

**本报告快照说明**：本报告为 ⑩ 时点快照；清单层事实偏差 4 项（§12.1 引用精度、gates-cache 例外集范围、⑥ 输出契约 `##` 前缀、检查器判据形态）与未闭环项 4 项（v1.56 排除条件、⑥ 格式偏离系统性、⑧ 取证落盘路径、任务包自定判据）见 ⑩ 回执的事实包第 10/12 字段，将于 docs-only 回填时下沉到本报告 §1 与 §11。

**清单层事实偏差（4 项，docs-only 回填下沉）**：① `§12.1` 本体无「治理片不满足试点选择标准」字面句（该表述是 `§12.4` 各 OB 行的既定写法）；② `tmp/gates-cache/` 的运行时例外集仅 `g4b/` 子目录，不含 `gates-cache` 根层；③ ⑥ 输出契约要求五段标题带 `##` 前缀，裸行 `Section A` 必致 `G7-f` FAIL；④ 检查器判据形态不严谨（③ 的「CN/EN 句子数相等」与 ⑨ 的「标题式 § 引用」各产生 1 条误报）。

**F-1 补句的反引号对齐**：F-1 补句原文中 CN/EN 反引号数不对称（CN 无 `##`、EN 有），由主控补齐使两侧均为 354；属格式对齐，非语义改动；已如实登记。

## 2 执行流水线

| 步骤 | 角色 / 载体 | 模型 | 状态 | 返工次数 |
| --- | --- | --- | --- | --- |
| 〇·零与步骤 0 | 主控 | 主控模型 | 自检 0 发现 | 0 |
| ① 方案 | `deep-reasoner` | `DeepSeek V4 Pro` | 三组分叉对比表，停等裁定 | 1 |
| 批量辅助 | 只读子代理 ×3 | `GLM-5.3-Flash` | 枚举 / 比对 / 影响面 | 3 |
| ② 实现 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 首落笔 + 返工 + 定点整改 | 3 |
| ③ 测试 | `prfrail-tester` | `DeepSeek V4.1 Flash` | `TEST: FAIL (1)`（判据误报，见 §1 第 6 项） | 1 |
| ④ 门禁 | 主控 | 门禁脚本 | 三作用域 `TOTAL_FAIL=0` | 2 |
| ⑤ 预审 | 预审员 | `DeepSeek V4 Pro` | `PRE-REVIEW: PASS` | 1 |
| ⑥ 独立扫描 | 独立扫描方 | 见 ⑩ 派单表 | 结论为 `PASS`（原文见 §5） | 1 |
| ⑦ 终审 | 独立终审方 | 见 ⑩ 派单表 | `FINAL REVIEW: PASS` | 1 |
| ⑧a 报告 | `prfrail-documenter` | `DeepSeek V4.1 Flash` | 本报告 | 1 |

耗时未逐轮计时，统一标 `未计时`（登记缺口见 §11 第 6 项）。

## 3 异常与兜底记录

见 §1：九项全部落在 §1 的第 1–9 项（本章与 §1 合并），本章只给计数口径。

| 计数 | 值 | 口径 |
| --- | --- | --- |
| R（兜底） | 2 | §9.2；R1 = ② 首轮照抄 `W3`–`W9`，R2 = ⑤ r1 两条 Low |
| 临时件 D | 1 | §5.1；② 返工轮 3 支脚本自建自删，非自愈档 |
| 偏离表编号 | 0 | 本片无临时件之外的偏离 |
| 探针 | 0 | 额度 0 |

## 4 环境与档位

| 项 | 值 |
| --- | --- |
| 工作目录 | `d:\LZProjects\prfrail` |
| 平台 | Windows（PowerShell 终端） |
| 基线 | `HEAD` = `53ecf7f` |
| 准则版本 | v1.56 |
| 档位 | M（治理片；⑦ 额度 1+3） |
| ⑥ 与 ⑦ 载体 | 独立扫描方、独立终审方（模型见 ⑩ 派单表） |
| 探针额度 | 0 |

## 5 审查结论摘要（⑤⑥⑦）

| 环节 | 结论行 | 判定 |
| --- | --- | --- |
| ⑤ 预审 r1 | `PRE-REVIEW: PASS` | 2 Low + 若干 Info；两条 Low 已在 ② 定点整改 |
| ⑥ 独立扫描 | `INDEPENDENT SCAN: PASS` | Section A–E 齐备；Section D 无 High/Medium 发现；原文逐字转录见下方围栏 |
| ⑦ 终审 | `FINAL REVIEW: PASS` | 八维全通过、发现全 0 |

**三件套**：报告对 `docs/validation/directive-tmp-self-heal{,_EN}.md`；自由度清单 `docs/validation/evidence/DIRECTIVE-TMP-SELF-HEAL-freedom-list.md`；门禁读数见 §8。

**⑥ 独立扫描回执原文（逐字转录，未改写）**：

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

No findings.

## Section C

Q1: Is the v1.55 → v1.56 metadata update complete and internally consistent?  
PASS. The header and version bump are aligned in the CN and EN directive, and the changelog entry is present at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L102) and [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md#L102). The §0.3 version line is also updated in the same change surface.

Q2: Does the new same-command create/delete exception remain narrowly scoped, and are the non-self-heal and cumulative-cap closures present?  
PASS. The new clause is limited to the §5.1 self-heal tier and explicitly states the three required conditions: position valid, same-command create+delete, and zero residual evidence after the fact; it also rejects out-of-scope placements such as tmp root, repo root, or project subdirectories. The minimal closures are present in the same section at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L516-L520).

Q3: Is the §9.3 "提前清理" disposition cell extended without widening the policy beyond the allowed boundary?  
PASS. The cell now references both the self-heal exception and the new same-command exception, and the boundary expansion is explicitly recorded in the slice scope. See [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L902).

Q4: Are the appendix and OB tracking consistent with the slice's actual scope and deferral status?  
PASS. Appendix C.0aq is present and framed as a summary trace rather than a new source of law; OB-128 is recorded as the 25th deferral and tied to the governance-slice nature of this change at [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1148) and [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L2080-L2084).

Q5: Do the P34 ledger rows and boundary lines in both CN and EN correctly reflect the slice status and boundaries?  
PASS. The P34 rows show the slice as in progress and the boundary notes are aligned in both language mirrors: [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L66), [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L1139-L1157), [docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L66), and [docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L1126-L1144).

The recheck across these exact edit points finds no semantic drift, no missing closure, and no boundary leak beyond the declared slice.

## Section D

UNCOVERED: none.

The guardrails that matter here are already spelled out in the text itself: the same-command exception is explicitly bounded, the non-self-heal path remains enforced, and the cumulative D cap still applies. If those protections were removed, the likely failure would be human misinterpretation rather than a mechanically visible contradiction in the directive text; this is exactly why the new section and the closure lines were added in the same patch.

## Section E

I cannot verify downstream adoption beyond the repository text itself, but the current evidence supports the PASS conclusion. What would overturn this PASS would be:
- a broader same-command exception than the tmp/<slice-id> bounded path,
- a bypass of the non-self-heal fail-stop rule,
- or a silent override of the cumulative D limit.

The current wording explicitly forbids all three, and the corresponding ledger and appendix entries remain consistent with that rule set.

INDEPENDENT SCAN: PASS
```

⑥ 独立扫描：结论 `INDEPENDENT SCAN: PASS`，Section A–E 齐备，Section D 无 High/Medium 发现；原文见上方围栏。本围栏为 ⑥ 重发版（首轮版五段标题缺 `##` 前缀致 G7-f FAIL，按 §7.2.2 ⑥ 内部额度重发；不计入 §9.2 R）。

### ⑧a 自查

- ① 成本表机械审计：有 —— §10 九项合计 `12` 与 §2 表的调用行一一对应。
- ② 「已处置」`file:line`：有 —— `OB-128` 的「已裁决」在 `docs/DELIVERY_DIRECTIVE.md:1148`，镜像同位。
- ③ 旧口径 grep：有 —— 读数见 §8 末表（两条检索键、CN 与 EN 各一组）。
- ④ ⑥ 与 ⑦ 的 Info / Low 汇总：有 —— ⑤ 两条 Low 见 §1 第 1 项、⑦ 发现 0 见 §5、⑥ 原文逐字转录见 §5 围栏。

### A-NN 编号完整性

本报告不使用 `A-NN` 编号；例外与缺口以 §1 的编号项与文件路径定位，故无须补实体定义。

## 6 审查输入包摘要（含盲审隔离证明）

| 环节 | 输入范围（按 §7.2.2 口径） | 隔离要求 |
| --- | --- | --- |
| ⑤ | 起跑消息、① 方案书、实现 diff、门禁读数 | 可读裁定记录 |
| ⑥ | 切片字段、实现 diff、契约原句、只读约束 | 盲包，剔除裁定记录与 ⑤ 结论 |
| ⑦ | ①–⑥ 全套、偏离与豁免登记、门禁读数 | 全量输入 |

**盲审隔离证明**：⑥ 的输入包按 §7.2.2 不含 ⑤ 的结论与裁定记录；本 ⑧a 输入未附输入包的逐件清单，故此处只登记口径、不登记包内路径（缺口见 §11 第 5 项）。

## 7 可证伪性（机械检查 + 变异检验）

| 层 | 机制 | 覆盖 | 不覆盖 |
| --- | --- | --- | --- |
| L1 | ⑩ 前终态回读与 §5.1 三项读数 | 特例要素缺项 ⇒ 判红 | 特例是否被误用，本身无机械侦测 |
| L2 | `G3(a)` 六形状与 `G3(b)` 关键字段 | 单边改字 ⇒ 判红 | 两侧语义是否等义 |
| L3 | `G6-5` / `G6-7` 段引用与台账对账 | 引用悬空 ⇒ 判红 | v1.56 排除条件是否过严 |

变异检验中的 `M1`、`M2` 为推演（改写准则正文超本片边界），`M3` 为内存实测（只读文件、不写盘）：

| 编号 | 变异体 | 会红的判据 | 不会红的判据 |
| --- | --- | --- | --- |
| `M1` | 删「位置合法而仅次序违规的，不适用本特例」句 | 无 ⇒ 只能由 ⑤ / ⑦ 人工发现 | 全部门禁 |
| `M2` | 删「审查层确认」句 | 无 ⇒ `G6-5` 只判 §X.Y 引用 | 全部门禁 |
| `M3` | 单边删 CN 的一个 `**` | `G3(a)` 六形状报 `bold` 不等 | `G3(b)` 关键字段 |

## 8 门禁结果

| 作用域 | 读数 |
| --- | --- |
| `--scope=index` | `TOTAL_FAIL=0`；`G6-1 scanned 86`；`G3(a)` directive `+20/-4` 对称、六形状全 `true`；`G3(b)` hash `1/1`；`G4b newCjk=453 unknown=[]` |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `--selftest` | `SELFTEST: PASS 697/697` |

上述读数为 ⑥ 扫描时点的快照；本 ⑧a 新增报告后 `G4b` 的 `newCjk` 随新增行增长（属时点差异，口径见 §3），`TOTAL_FAIL=0` 的结论不变。

**旧口径 grep 读数**（⑧a 自查 ③ 的取证）：

| 侧 | 检索键 | 命中 |
| --- | --- | --- |
| CN | `自愈档（v1.45）或同命令建删特例（v1.56）` | 3 |
| EN | `self-heal tier (v1.45) or same-command create-and-delete exception (v1.56)` | 3 |
| CN | `自愈档（v1.45 同步）` | 0 |
| EN | `v1.45 sync` | 0 |

## 9 ⑧c 收尾清洁记录

**⑧c 清理登记（记录先于清理，依 §5.1「清理与落痕的时序」）**：

- 待清理清单（逐件，5 件，均在本片临时目录内）：`tmp/DIRECTIVE-TMP-SELF-HEAL/apply.js`、`tmp/DIRECTIVE-TMP-SELF-HEAL/verify.js`、`tmp/DIRECTIVE-TMP-SELF-HEAL/step0-insert.js`、`tmp/DIRECTIVE-TMP-SELF-HEAL/DELIVERY_DIRECTIVE.md.bak`、`tmp/DIRECTIVE-TMP-SELF-HEAL/DELIVERY_DIRECTIVE_EN.md.bak`。两件 `.bak` 为编辑前快照，**亦属临时件**，一并列入。
- 清理方式：`Remove-Item -LiteralPath` 逐件执行，不用通配符；清单外文件一律不动（含 `tmp/gates-cache/g4b/` 运行时缓存与 `tmp/.gitkeep`）。
- 清理后读数：见本节末尾「清理后回读」。

**清理后回读（五条，均于 ⑧c 实测）**：① `git status --porcelain -uall` 的 `??` 项 = 0、第二列非空 = 0（`EXIT=0`）；② `tmp` 递归文件数（排除 `gates-cache/g4b`）= 1（`EXIT=0`）；③ `tmp` 顶层排除 `gates-cache` 后 = `.gitkeep`（1 项，`EXIT=0`）；④ `Test-Path tmp/DIRECTIVE-TMP-SELF-HEAL` = `False`（`EXIT=0`）；⑤ `gates-cache` 递归目录数 = 1（`g4b`，`EXIT=0`）。清理登记 5 件 → 逐件 `deleted` → `files_left=0` → 目录已删。

现状（⑧a 时点）：`tmp/DIRECTIVE-TMP-SELF-HEAL/` 共 5 件（`apply.js`、`DELIVERY_DIRECTIVE.md.bak`、`DELIVERY_DIRECTIVE_EN.md.bak`、`step0-insert.js`、`verify.js`）；`tmp/gates-cache/` 为运行时缓存（例外集成员）；`tmp/` 顶层另有 `.gitkeep`。

## 10 成本与计量

按 R7.1 逐环节列调用次数：

| 环节 | 调用次数 |
| --- | --- |
| ① 方案 | 1 |
| 批量辅助 | 3 |
| ② 实现（首落笔 + 返工 + 定点整改） | 3 |
| ③ 测试 | 1 |
| ⑤ 预审 | 1 |
| ⑥ 独立扫描 | 1 |
| ⑦ 终审 | 1 |
| ⑧a 报告 | 1 |
| 探针 | 0 |

## 11 明确未执行事项

1. commit / push / CI：均未执行，须 ⑩ 停点后同轮授权。**⑩ 停点后已完成**：commit `acf676e`；CI run `40695233778`，ubuntu 与 windows 两作业均 success（两腿 `SELFTEST: PASS 697/697`、`G6-1 scanned 652`、`TOTAL_FAIL=0`）。
2. 真实探针：未执行（额度 0）。
3. `GATES-EXEMPT` 实现：未执行（转办项，留后续片）。
4. v1.56 排除条件的放宽：未执行（留 v23 观察项候选）。
5. 审查输入包逐件清单：未登记（输入未附清单，只登记口径，见 §6）。
6. 逐轮耗时：未计时（§2 统一标 `未计时`）。

**未闭环项（4 项，供 v23 队列参考）**：① v1.56 排除条件是否过严（位置合法 + 同命令建删应否纳入特例；触发 = 本片 D = 1 首次实证）；② ⑥ 格式偏离跨两片系统性（P33 与 P34 各一次，均属契约细节类偏离）；③ ⑧ 取证落盘路径未限定（本片 D = 2 即由 `gates-cache` 根写日志引起）；④ 任务包自定判据缺前置评审（本片 ③ 与 ⑨ 各产生 1 条误报）。

## 12 先例禁用核对

- **本片例外两项**：`R = 2` 与 `D = 1`，二者均为超常规上限的事件登记。
- `D = 1` 的 `OB-37` 元规则授权**不得作先例**（§1.7.4 ④ 与 §7.8.4 的元规则 ③ 同标记）。
- **扩界**（§9.3 处置格）属用户同轮授权，登记为**事实**而非先例。
- 其余引用（`C.0aq` 的「不得作先例」句）只作跨片汇总留痕，本片不据此主张豁免。

## 13 下一步建议

1. ⑨ 原生验证 → ⑧b 回写 → ⑧c 收尾清洁 → ⑩ 停点回执；⑧c 后整改不得重建 `tmp/DIRECTIVE-TMP-SELF-HEAL/`。
2. ⑥ 原文已逐字转录于 §5 的 `sixstage` 围栏（重发版，五段标题带 `##` 前缀，满足 G7-f 契约）；§11 已删除原第 5 项。
3. v23 观察项候选 3 条（§1 第 9 项）与 `GATES-EXEMPT` 转办项一并评估。
4. 档位建议：本类治理片维持 M 档；触及判据语义的轮次由 ⑦ 取高强度载体，纯文档轮次可降档。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-tmp-self-heal.md	repo
docs/validation/directive-tmp-self-heal_EN.md	repo
docs/validation/evidence/DIRECTIVE-TMP-SELF-HEAL-freedom-list.md	repo
```

# P72 REQ020-S2-RESTRICTION — ① 现状核对与差异清单（design-r1）

## §0 输入与依据

- **切片卡**：`docs/t027/REMAINING_SLICES.md:2000–2018`（CN）/ `REMAINING_SLICES_EN.md:1987–2005`（EN）。目标原文（CN `:2006`）：
  ```text
  | 目标 | REQ-020 的阶段归属列改为「S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力）」；同步 `TEST_STRATEGY` 的 AT-01 与 `S1_SUPPORT_MATRIX` 的语言行 |
  ```
  边界（CN `:2013`）：只改 REQ-020 阶段列 + 关联文档；不改 MVA 复合标准 / 判据本体 / 产品代码；不推 gitee。
- **冻结依据链**：P68 R-07 b1 ③（Q3 语言冻结）→ 用户裁定（`remaining-rulings{,_EN}.md` §5.4 `:105`）→ Q3 冻结范围原文（`01-arch/r07-refined.md:49`）→ P71 ⑤ Info 建议登记（`s1-scope-annotate{,_EN}.md:90`）。
- **核验方式**：2026-11-14 工作树实测（`Select-String` 行号 + node 字节级 JSON 转义）；终端显示的长行中偶见幻影空格（渲染伪影），正文一律以磁盘字节为准。

## §1 现状核对

### 1.1 REQ 表（`docs/PRODUCT_REQUIREMENTS{,_EN}.md`）

实测：CN 表头 `:61`；REQ 条目 28 条、编号 001–028 连续（机械复核，见 §4）；EN 同。

- CN `:73` REQ-011：
  ```text
  | REQ-011 | R11 | S1/S2 | S1 generic/C/Go harness 外置，核心无语言判断分支；S2 扩展语言 |
  ```
- EN `:71` REQ-011：
  ```text
  | REQ-011 | R11 | S1/S2 | External generic/C/Go harnesses in S1, no language branching in core; extend languages in S2 |
  ```
- CN `:82` REQ-020：
  ```text
  | REQ-020 | R20 | S0-S2 | S0 验证双/三语言配置；S1 C+Go 契约；S2 三语言实跑，失败整任务阻断 |
  ```
- EN `:80` REQ-020：
  ```text
  | REQ-020 | R20 | S0-S2 | S0 dual/triple-language configs; S1 C+Go contracts; S2 three-language execution with whole-task failure |
  ```
- CN `:84` REQ-022 / EN `:82`：
  ```text
  | REQ-022 | R22 | S1/S2 | 文档一等 target、确定性 impact rule、生成方向、waiver 和整体验收；S2 智能提示 |
  | REQ-022 | R22 | S1/S2 | First-class docs, deterministic impact rules, generation direction, waivers and aggregate acceptance; S2 intelligent hints |
  ```
- PC 行：CN PC-01 `:119` / PC-05 `:123` / PC-06 `:124`；EN PC-01 `:117` / PC-05 `:121` / PC-06 `:122`。实测该 6 行（CN+EN）对 `Python|Node|三语言|triple|dual` **零命中**（PC-06 的「支持矩阵」指文档名，见 CN `:124`）。

### 1.2 AT 表（`docs/TEST_STRATEGY{,_EN}.md`）

- CN `:31` AT-01：
  ```text
  | AT-01 | REQ-001/010/017/018/020/022/028 | Schema 合法配置；未知字段/版本、重复 ID、空 steps、非法 noop、重叠 target、依赖/生成环拒绝；C+Go/三语言/交接/文档样例 | taskdef；T002/T003 |
  ```
- EN `:31` AT-01：
  ```text
  | AT-01 | REQ-001/010/017/018/020/022/028 | Valid schemas; reject unknown fields/versions, duplicate IDs, empty steps, invalid noop, overlapping targets and dependency/generation cycles; dual/triple-language, handoff and doc fixtures | taskdef; T002/T003 |
  ```
- CN `:41` AT-11（`C+Go` 为「同 task 整体恢复」恢复语义）：
  ```text
  | AT-11 | REQ-011/013/020/022 | generic/C/Go、锁模板生成 A；自批/手改生成物/文档缺失阻断；C+Go 同 task 整体恢复 | taskdef/gates；T015 |
  ```
- EN `:41` AT-11：
  ```text
  | AT-11 | REQ-011/013/020/022 | generic/C/Go, locked template A; self-approval/manual generated edits/missing docs block; C+Go whole-task recovery | taskdef/gates; T015 |
  ```
- CN `:50` AT-20 / EN `:50`：PC-05 成本类，无任何语言表述（已核）。

### 1.3 `docs/S1_SUPPORT_MATRIX{,_EN}.md`

CN 表头 `:7` = `| 范围 | S1 候选状态 | 约束 |`；数据行 6 行（`:9`–`:14`）：Windows 11 amd64（正式支持候选，`:9`）、Linux amd64（核心 CI，`:10`）、macOS/WoA/Linux arm64（不支持，`:11`）、generic/C/Go harness（`:12`）、CLI Agent/AgentRunner（`:13`）、VS Code visible 黑箱候选（`:14`）。**无独立 C 行 / Go 行**；`^\| (Python|Node)` 零命中（实测）。harness 行：

- CN `:12`：
  ```text
  | generic/C/Go harness | 配置与 noop 编排 | language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close |
  ```
- EN `:12`：
  ```text
  | generic/C/Go harnesses | Configuration and noop orchestration | Language scopes, read-only preview, and noop-only loops are verified; executable steps still fail closed |
  ```

### 1.4 冻结原文（P68 R-07 b1 Q3）

- `docs/validation/evidence/P68-REMAINING-RULINGS/01-arch/r07-refined.md:49`（Q3 冻结范围原文，**逐句对应目标文本**）：
  ```text
  **S2 收窄（REQ-020 三语言）**：`docs/PRODUCT_REQUIREMENTS.md:82` 定 S0 验证双/三语言配置、**S1 只要求 C+Go 契约**、S2 三语言实跑。Q3 建议的冻结范围为「Python/Node **官方 harness**、三语言产品任务与正式支持承诺」，**不冻结**既有工具脚本、开放语言 ID、generic 能力或黄金配置，也不新增 Python/Node 核心运行时依赖 ⇒ **冻结不改变 S1 范围与既有可用性**，只收窄 S2 承诺面。`docs/S1_SUPPORT_MATRIX.md` 现状亦为 generic/C/Go = 配置与 noop 编排已验证 ⇒ generic 覆盖不受影响。
  ```
- `remaining-rulings.md:87`（EN `:87`）：
  ```text
  - 附加：AT-13 未通过（…）；S2 语言冻结不影响 S1（REQ-020 的 S1 只要求 C+Go；generic/C/Go 已验证）。
  - Additional: AT-13 has not passed (…); freezing S2 languages does not affect S1 (REQ-020's S1 only requires C+Go; generic/C/Go already verified).
  ```
- `remaining-rulings.md:92`（EN `:92`）方案 b1 定义；`:105` 用户裁定「方案 b1 + 附条件」（CN/EN 同结构）。

### 1.5 P71 ⑤ Info 建议原文

- CN `docs/validation/s1-scope-annotate.md:90`（§7 已知边界 ③）：
  ```text
  - ③ **REQ-020 的「S2 承诺收窄登记」不在本片范围**（⑤ Info 建议登记，留待后续切片）。
  ```
- EN `:90`：
  ```text
  - (3) **The "S2 commitment narrowing registration" for REQ-020 is outside this slice** (a step-5 Info item suggests registering it; it is left to a later slice).
  ```
- 切片卡来源字段：CN `REMAINING_SLICES.md:2004` = 「来源 = P71 ⑤ Info 建议登记 + P68 b1 ③ 的 S2 收窄（Q3 语言冻结）」。

### 1.6 P57 §4.2 八门清单

`docs/validation/p57-closeout-assessment.md:80`（EN `:80`）门表 S-G1..S-G8 = AT-23 真实候选腿、AT-24 与 T028、最终 ZIP 绑定、EOL 通知、产品所有者批准、T018 整体、S1 exit 总体、真实 CLI 端到端（CN `:84`–`:91` / EN `:84`–`:91`）。**八门无语言支持门**。`:93` 已有「P71 指针（2026-11-13）」行。

### 1.7 `docs/validation/s1-exit{,_EN}.md` 的 AT-01 行

**有 AT-01 行**（状态汇总行，非阶段归属行）：

- CN `:21`：
  ```text
  | AT-01 | PASS | `DEV_PLAN.md` T002/T003：29 Schema、82 fixture、JCS 向量 |
  ```
- EN `:21`：
  ```text
  | AT-01 | PASS | `DEV_PLAN_EN.md` T002/T003: 29 schemas, 82 fixtures, JCS vectors |
  ```
- CN `:15` 边界句与冻结方向一致：「三语言样例使用当前 noop-only CLI，不冒充 executable hook 或编译器已执行」。

### 1.8 `docs/S1_RELEASE_NOTES{,_EN}.md` 支持段

CN `:13` 能力清单 = 「generic/C/Go 最小样例与脱敏 whois 冻结夹具只读影子验证」；`:16`–`:18` 支持与已知限制段**无** Python/Node 官方支持、**无**三语言实跑承诺（EN `:13`/`:16` 对应）。全文件 `Python|Node|三语言|triple` 零命中（实测）。

## §2 差异清单

| 文档 | 现状 vs Q3 冻结要求 | 性质 | 理由 |
|---|---|---|---|
| `PRODUCT_REQUIREMENTS{,_EN}.md` | REQ-020 阶段列 `S0-S2` 无收窄注记 vs 冻结要求 S1 仅 C+Go 契约、S2 承诺收窄 | **修改**（仅第 3 列单元格） | 切片核心目标；第 4 列「可验收要求」零改动 |
| `TEST_STRATEGY{,_EN}.md` AT-01 | 场景列含 `C+Go/三语言`（EN `dual/triple-language`）未区分阶段归属 vs 冻结要求三语言实跑归 S2 | **修改**（仅该片段） | 阶段归属同步表述；其余单元格逐字不变 |
| `TEST_STRATEGY{,_EN}.md` AT-11 | `C+Go 同 task 整体恢复` = S1 恢复语义，非支持承诺阶段归属 | **不动** | 见 §3.3 判定 |
| `S1_SUPPORT_MATRIX{,_EN}.md` | harness 行约束列无三语言承诺阶段注记 | **加法式加注**（约束列末尾追加一句） | 冻结属承诺收窄登记，不动状态列与行数 |
| `REMAINING_SLICES{,_EN}.md` P72 段 | 状态 🔄 进行中、无完成记录行 | **修改**（⑧ 阶段） | ⑧ 阶段 🔄→✅ + 完成记录行 |
| `remaining-rulings{,_EN}.md` §5.4 | 无 Q3 冻结落地指针 | **加法式**（1 行指针） | 仿 P71 指针先例，闭合裁定→落地追溯链 |
| `s1-scope-annotate{,_EN}.md` §7 ③ | 「留待后续切片」待办未消解 | **加法式**（1 行指针） | P72 即该后续切片，待办须闭环 |
| `S1_RELEASE_NOTES` / `s1-exit` / `p57` / `CHECKLIST_v57` / `DEV_PLAN` / `OPERATIONS` / 准则 / `CONTRACTS` / `SECURITY` | 无语言阶段归属或可验收要求引用 | **不动** | 见 §6 判定 |

补充观察：阶段列新值不显式写 `S0`，但「S0 验证双/三语言配置」事实保留于第 4 列（零改动不变量）；与 P71 对 REQ-009/007/021 的 X 标注同构（阶段列只表述未收口的承诺）。

## §3 回写方案

### 3.1 REQ-020「阶段」列（修改）

CN `docs/PRODUCT_REQUIREMENTS.md:82` 目标整行：

```text
| REQ-020 | R20 | S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力） | S0 验证双/三语言配置；S1 C+Go 契约；S2 三语言实跑，失败整任务阻断 |
```

EN `docs/PRODUCT_REQUIREMENTS_EN.md:80` 目标整行（逐句对应 CN；术语直接采用仓库既有译本 `REMAINING_SLICES_EN.md:1993` 目标字段的既定用词 `official harness` / `release-support commitment` / `language IDs`，不再新造）：

```text
| REQ-020 | R20 | Not required by S1 (C+Go contract only); S2 commitment narrowed (Python/Node official harness + three-language product tasks + the release-support commitment frozen; existing tool scripts / open language IDs / generic capability are not frozen) | S0 dual/triple-language configs; S1 C+Go contracts; S2 three-language execution with whole-task failure |
```

逐句对称核对：`S1 不要求（仅 C+Go 契约）`↔`Not required by S1 (C+Go contract only)`；`S2 承诺收窄`↔`S2 commitment narrowed`；`Python/Node 官方 harness`↔`Python/Node official harness`；`三语言产品任务`↔`three-language product tasks`；`正式支持承诺冻结`↔`the release-support commitment frozen`；`不冻既有工具脚本 / 开放语言 ID / generic 能力`↔`existing tool scripts / open language IDs / generic capability are not frozen`。

### 3.2 TEST_STRATEGY AT-01（修改，仅场景列片段）

CN `:31` 目标整行（其余单元格逐字不变）：

```text
| AT-01 | REQ-001/010/017/018/020/022/028 | Schema 合法配置；未知字段/版本、重复 ID、空 steps、非法 noop、重叠 target、依赖/生成环拒绝；C+Go 契约（三语言实跑冻结至 S2，见 REQ-020）/交接/文档样例 | taskdef；T002/T003 |
```

EN `:31` 目标整行（EN 列表以逗号分隔，替换 `dual/triple-language` 片段）：

```text
| AT-01 | REQ-001/010/017/018/020/022/028 | Valid schemas; reject unknown fields/versions, duplicate IDs, empty steps, invalid noop, overlapping targets and dependency/generation cycles; C+Go contract (three-language execution frozen to S2, see REQ-020), handoff and doc fixtures | taskdef; T002/T003 |
```

### 3.3 AT-11 判定：**不动**

- 判定：不改。依据：AT-11 场景列的 `C+Go 同 task 整体恢复`（EN `C+Go whole-task recovery`）是 **S1 内的恢复语义断言**（任务级原子恢复），不是三语言支持承诺的阶段归属；Q3 冻结范围（`r07-refined.md:49`）为「Python/Node 官方 harness、三语言产品任务与正式支持承诺」，不触及 C+Go 恢复语义；收窄文本明示「不冻…generic 能力」（generic/C/Go 是 S1 已承诺能力）。
- 附注：AT-15（`小 Go 非 C 夹具`，CN `:46`）与 AT-20 同属 S1 能力面 / 成本类，均不动。

### 3.4 S1_SUPPORT_MATRIX generic/C/Go 行（加法式加注）

「S1 候选状态」列确认保持 `配置与 noop 编排`（EN `Configuration and noop orchestration`）不变；行数 6 不变；不新增 Python/Node 独立行。

追加段（CN）：

```text
Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力
```

追加段（EN）：

```text
Python/Node official harness + three-language product tasks + the release-support commitment frozen to S2 (REQ-020 S2 commitment narrowing); existing tool scripts / open language IDs / generic capability are not frozen
```

CN `:12` 目标整行（`【追加段起】`/`【追加段止】`为标注，不写入文件）：

```text
| generic/C/Go harness | 配置与 noop 编排 | language scope、只读 preview 和 noop-only 闭环已验证；executable step 仍 fail-close；【追加段起】Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力【追加段止】 |
```

EN `:12` 目标整行：

```text
| generic/C/Go harnesses | Configuration and noop orchestration | Language scopes, read-only preview, and noop-only loops are verified; executable steps still fail closed; 【append starts】Python/Node official harness + three-language product tasks + the release-support commitment frozen to S2 (REQ-020 S2 commitment narrowing); existing tool scripts / open language IDs / generic capability are not frozen【append ends】 |
```

### 3.5 REMAINING_SLICES P72 段（⑧ 阶段目标文本）

CN（`REMAINING_SLICES.md`）：

```text
:2000 行 → ### P72 — `REQ020-S2-RESTRICTION`（REQ-020 的 S2 承诺收窄登记） · ✅ 已完成（2026-11-14 收口）
:2005 行 → | 状态 | ✅ 已完成（2026-11-14 收口） |
:2018 行后新增 → | 完成记录 | 报告对 `docs/validation/req020-s2-restriction{,_EN}.md`；审查链 ① `DESIGN: DONE` → ② 分步全绿（<⑧ 实填项数> 项） → ③ 一致性 → ⑤ `PRE-REVIEW: PASS` → ⑥ 跳过（⑦ 确认） → ⑦ 盲审 `PASS` + 终审 `PASS`；GLM 批量 <⑧ 实填> 次；R=<⑧ 实填>、D=<⑧ 实填>；停于 ⑩ 等授权（不写 hash/run，§11.2） |
```

EN（`REMAINING_SLICES_EN.md`）：

```text
:1987 行 → ### P72 — `REQ020-S2-RESTRICTION` (registering the S2 commitment narrowing for REQ-020) · ✅ Completed (closed 2026-11-14)
:1992 行 → | Status | ✅ Completed (closed 2026-11-14) |
:2005 行后新增 → | Completion record | Report pair `docs/validation/req020-s2-restriction{,_EN}.md`; review chain ① `DESIGN: DONE` → ② all green step by step (<filled at ⑧> items) → ③ consistency → ⑤ `PRE-REVIEW: PASS` → ⑥ skipped (confirmed by ⑦) → ⑦ blind `PASS` + final `PASS`; GLM batch <filled at ⑧> times; R=<...>, D=<...>; stopped at ⑩ pending authorisation (no hash/run recorded, §11.2) |
```

注：`<⑧ 实填>` 为 ②–⑦ 实际结果占位，⑧ 阶段以真实值替换；收口日期以实际 ⑧ 日为准；报告对文件名按切片名小写化惯例推定（P71 先例 `s1-scope-annotate{,_EN}.md`），⑧ 以实际定稿路径为准。

### 3.6 指针行（加法式，适用项，见 §6 判定 1/5）

`remaining-rulings.md:107` 后新增一行：

```text
- **P72 落地指针（2026-11-14）**：R-07 b1 ③ 的 S2 收窄（Q3 语言冻结）已由切片 `REQ020-S2-RESTRICTION` 落地，报告见 `docs/validation/req020-s2-restriction{,_EN}.md`；本文件原行不改（§11.1）。
```

`remaining-rulings_EN.md:107` 后：

```text
- **P72 landing pointer (2026-11-14)**: the S2 narrowing in R-07 b1 item 3 (the Q3 language freeze) has been landed by the `REQ020-S2-RESTRICTION` slice; the report is `docs/validation/req020-s2-restriction{,_EN}.md`. The original lines above stay untouched (§11.1).
```

`s1-scope-annotate.md:90` 后新增一行：

```text
- ③a **P72 落地指针（2026-11-14）**：③ 所指「REQ-020 的 S2 承诺收窄登记」已由切片 `REQ020-S2-RESTRICTION` 落地，报告见 `docs/validation/req020-s2-restriction{,_EN}.md`；③ 原行不改。
```

`s1-scope-annotate_EN.md:90` 后：

```text
- (3a) **P72 landing pointer (2026-11-14)**: the "S2 commitment narrowing registration" for REQ-020 referred to in item 3 has been landed by the `REQ020-S2-RESTRICTION` slice; the report is `docs/validation/req020-s2-restriction{,_EN}.md`. The original line above stays untouched.
```

## §4 不变量

1. **REQ 条目数 28 不变、编号连续**：`(Select-String -Path docs\PRODUCT_REQUIREMENTS.md -Pattern '^\| REQ-(\d{3}) ').Count` = 28 且提取的编号序列 = `001..028`（node 逐项比对；EN 同）。基线实测已通过。
2. **REQ-020「可验收要求」列零改动**：git diff 该行仅第 3 个 `|` 字段变化；node 按 `|` 拆列比对第 4 字段 before/after 逐字节 identical。
3. **MVA 复合标准不动**：`git diff --name-only` 结果 ⊆ 本文 §3 所列 6 对文件；不含 `strategic-planning*` / `DELIVERY_DIRECTIVE*` / `CONTRACTS*` / `SECURITY*`。
4. **S1 Exit 八门判定本身不动**：diff 不含 `s1-exit*`、`p57-closeout-assessment*`；八门表逐行 unchanged。
5. **S1_SUPPORT_MATRIX 行数 6 不变**：`(Select-String -Path docs\S1_SUPPORT_MATRIX.md -Pattern '^\| ').Count` = 7（表头 + 6 数据行）；`'^\| (Python|Node)'` 命中 = 0；harness 行第 2 字段 = `配置与 noop 编排`（EN `Configuration and noop orchestration`）。基线实测已通过。
6. **AT-01 其余单元格逐字不变**：diff 该行仅 `C+Go/三语言`（EN `dual/triple-language`）片段变化；改后行与 §3.2 目标整行逐字节相等。

## §5 反例清单

1. **禁止**改 REQ-020 第 4 列「可验收要求」。检出：§4-2 拆列比对 identical；或 `git diff --word-diff` 该行仅第 3 字段。
2. **禁止**增删 REQ 行或改编号。检出：§4-1 计数 = 28 且编号序列 = `001..028`。
3. **禁止**在 S1_SUPPORT_MATRIX 新增 Python/Node 独立行（行数 6→7）。检出：§4-5 计数 = 7 且 `^\| (Python|Node)` = 0。
4. **禁止**改 harness 行「S1 候选状态」列（须仍为 `配置与 noop 编排`）。检出：§4-5 第 2 字段比对。
5. **禁止**把冻结扩大为「不冻」保留句缺失（如删除「不冻既有工具脚本 / 开放语言 ID / generic 能力」）。检出：REQ-020 行与 §3.1 目标整行逐字节比对。
6. **禁止**顺带修改准则正文 / `CONTRACTS` / `SECURITY` / p57 八门 / s1-exit 门表 / CHECKLIST_v57。检出：§4-3/§4-4 的 diff 文件集合限定。

## §6 需评估 5 项判定

1. **`remaining-rulings{,_EN}.md` §5**：**适用**。依据：P68 §5.4 `:107` 已有 P71 指针行先例（「已由切片落地…本文件原行不改」）；P72 落地后 R-07 b1 ③ 缺落地指针将使「建议→裁定→落地」链断裂。加法式 1 行（§3.6），不删改原行。
2. **`S1_RELEASE_NOTES{,_EN}.md`**：**不适用**。依据：CN `:13`/`:18`、EN 对应行未提及 Python/Node 官方支持或三语言实跑（全文件关键词零命中，实测）；冻结是承诺收窄、不产生既有文本冲突；且不在 P72 交付物清单。
3. **`s1-exit{,_EN}.md` 的 AT-01 行**：**不适用（行存在，不改）**。依据：CN `:21`/EN `:21` 为状态汇总行（PASS + T002/T003 证据 = schema/fixture 层面，非 S2 三语言实跑）；CN `:15` 边界句已与冻结方向一致；本片不动门判定。
4. **`p57-closeout-assessment{,_EN}.md` §4.2**：**不适用（不改）**。依据：八门 S-G1..S-G8（CN `:84`–`:91` / EN `:84`–`:91`）无语言支持门；`:93` 已有 P71 指针；本片不动门判定，无需追加。
5. **`s1-scope-annotate{,_EN}.md` §5 待办项**：**适用**。依据：`:90` 明写「REQ-020 的「S2 承诺收窄登记」…留待后续切片」，P72 即该后续切片，待办须以加法式指针行消解（§3.6），原行不改（③ 作为 P71 历史记录保留）。

## §7 停升级边界逐条裁决

- (a) 改 REQ 表结构（增删条目/改编号/改可验收要求列）：**不触发**。仅改 REQ-020 第 3 列单元格；§4-1/§4-2 锁定。
- (b) 改 MVA 复合标准：**不触发**。P68 裁定方案 b1 明确「不动 MVA 作为验收复合标准」（`remaining-rulings.md:92`）；本片不改任何 MVA 文本（§4-3）。
- (c) 改 S1 Exit 八门判定本身：**不触发**。不改 `s1-exit` / `p57`（§4-4）。
- (d) 改 S1_SUPPORT_MATRIX 正式支持范围（非加法式加注）：**不触发**。仅约束列末尾追加一句；状态列、行数、平台行均不变（§4-5）。
- (e) 改准则正文 / CONTRACTS 语义 / SECURITY 边界：**不触发**。§4-3 diff 文件集合限定。
- (f) REQ-020 当前阶段列已不是 `S0-S2`：**不触发**。实测 CN `:82` / EN `:80` 现值为 `S0-S2`，停升级前提成立（即本片正当性成立）。
- (g) 新增 Python/Node 独立行：**不触发**。§3.4 明确不加行；§4-5 检出 `^\| (Python|Node)` = 0。
- (h) REQ-011/REQ-022/PC-01/PC-05/PC-06 是否引用「三语言/语言支持的阶段归属或可验收要求」：**裁决 = 采纳窄读，不触发**。论证：
  - 窄读依据（冻结原文 `r07-refined.md:49`）：被冻结对象有明确界定——「Python/Node 官方 harness、三语言产品任务与正式支持承诺」；同一句明示「**不冻结**既有工具脚本、开放语言 ID、generic 能力或黄金配置」。
  - REQ-011（CN `:73`）的「S2 扩展语言」属「开放语言 ID / generic 能力」范畴：该行全文为「S1 generic/C/Go harness 外置，核心无语言判断分支；S2 扩展语言」——即语言机制经外置 harness + 开放语言 ID 实现，正是收窄文本明示保留的对象 ⇒ 不改、不触发。
  - REQ-022 / PC-01 / PC-05 / PC-06：实测零命中（§1.1）⇒ 不存在引用，不触发。
  - 宽读后果：若采纳「任何语言阶段字样都算引用」⇒ REQ-011 的「S2 扩展语言」触发 ⇒ 本片必须 `DESIGN: NEEDS_CLARIFICATION`。本设计**不采纳宽读**，理由：停升级门槛的目的（防止承诺语义变更逃逸审查）在窄读下已被满足——Q3 冻结对象（三语言正式支持承诺）与 REQ-011 措辞（开放语言扩展）语义不重叠，且收窄文本显式保留后者；宽读会把「不冻 open language IDs」的保留句变成自我矛盾。

## §8 结论

DESIGN: DONE

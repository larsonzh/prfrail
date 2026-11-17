# batch-01 — ADR-010 原文与 ADR_REGISTER 格式枚举（`GLM-5.3-Flash`）

> 性质：批量辅助事实枚举，不含判定。落盘者 = 主控（逐字转录，未改写）。

以下为纯事实枚举结果（所有摘录均经脚本逐字节核验；摘录中的空格/标点与磁盘文件逐字一致）。

════════════════════════════════════════
一、两张主表的列结构
════════════════════════════════════════

【CN `docs/ADR_REGISTER.md`】
- 表头行原文（`docs/ADR_REGISTER.md:9`）：

| ID | 决策/方案 | 替代与理由 | 状态/完成条件 |

- 分隔行（`docs/ADR_REGISTER.md:10`）：

|---|---|---|---|

- 列名逐个：① `ID`；② `决策/方案`；③ `替代与理由`；④ `状态/完成条件`。

【EN `docs/ADR_REGISTER_EN.md`】
- 表头行原文（`docs/ADR_REGISTER_EN.md:9`）：

| ID | Decision/proposal | Alternatives/rationale | Status/completion |

- 分隔行（`docs/ADR_REGISTER_EN.md:10`）：

|---|---|---|---|

- 列名逐个：① `ID`；② `Decision/proposal`；③ `Alternatives/rationale`；④ `Status/completion`。

【ADR 行单元格分隔形态（事实数据）】
- 两文件中 ADR-001 ~ ADR-047 共 47 行，全部按 `|` 切分后均为 **4 列**（distinct cell counts = [4]，无 3 列或 5 列行）。
- 两文件全文 **转义竖线 `\|` 出现次数 = 0**（即无任何行使用转义竖线）。
- 附带事实：两文件内还存在其它表格的分隔行——CN 在第 75、97 行，EN 在第 76、98 行（不属于 ADR 主表）。
- 文件总行数：CN 117 行，EN 118 行。

════════════════════════════════════════
二、ADR-010 完整行原文（CN 与 EN）
════════════════════════════════════════

【CN】`docs/ADR_REGISTER.md:20`（原始行共 311 字符），4 列单元格逐字内容：
- 第 1 列：`ADR-010`
- 第 2 列：`产品版本采用 \`vMAJOR.MINOR.PATCH\`（首版候选 \`v0.1.0\`）；原型期新 minor 发布时，在 GitHub Release notes 与仓库支持矩阵同时公告前一 minor 的 EOL，初始窗口为发布后 90 个自然日`
- 第 3 列：`patch 不启动新的 EOL 倒计时；每次 minor 发布前按稳定性、用户量和维护能力评审是否延长未来窗口。不得追溯缩短已公告窗口；例外须有带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处`
- 第 4 列（完整，含分号后段落）：`` `PROPOSED / NOT APPROVED`；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准 ``

【EN】`docs/ADR_REGISTER_EN.md:20`（原始行共 797 字符），4 列单元格逐字内容：
- 第 1 列：`ADR-010`
- 第 2 列：`` Product versions use `vMAJOR.MINOR.PATCH` (first candidate `v0.1.0`); during the prototype phase, a new minor announces the previous minor's EOL in both GitHub Release notes and the repository support matrix, initially 90 calendar days after release ``
- 第 3 列：`Patch releases do not start a new EOL countdown; before each minor release, review stability, adoption, and maintenance capacity to decide whether future windows should grow. Never retroactively shorten an announced window; exceptions require a dated ADR with rationale and expiry, explicit product-owner approval, and synchronized publication in both channels`
- 第 4 列（完整，含分号后段落）：`` `PROPOSED / NOT APPROVED`; 90 days fits the early prototype, but the notice channels and evolution rule require explicit product-owner approval before the first release ``

════════════════════════════════════════
三、ADR-006 / 007 / 008 / 009 的「批准记录」形态
════════════════════════════════════════

【CN 状态列逐字原文】
- ADR-006（`docs/ADR_REGISTER.md:16`）：`2026-09-08 所有者批准简化；T017 落实基础发布证据`
- ADR-007（`docs/ADR_REGISTER.md:17`）：`2026-09-06 所有者批准；S0 新增付费调用为 0，超出须另授权`
- ADR-008（`docs/ADR_REGISTER.md:18`）：`2026-09-06 所有者批准范围；T002 已冻结记录，T003 已补样例，T019–T024 运行验证；未授权 S1 编码`
- ADR-009（`docs/ADR_REGISTER.md:19`）：`2026-09-09 产品所有者明确批准；最终 ZIP 仍须绑定发行记录并复验`

【EN 状态列逐字原文】
- ADR-006（`docs/ADR_REGISTER_EN.md:16`）：`Owner simplified policy 2026-09-08; T017 supplies baseline release evidence`
- ADR-007（`docs/ADR_REGISTER_EN.md:17`）：`Owner approved 2026-09-06; zero new paid S0 calls, separate authorization for increases`
- ADR-008（`docs/ADR_REGISTER_EN.md:18`）：`Owner approved scope 2026-09-06; T002 froze records, T003 added goldens and T019–T024 run checks; no S1 coding authorization`
- ADR-009（`docs/ADR_REGISTER_EN.md:19`）：`Product owner explicitly approved 2026-09-09; bind and rerun the final ZIP before release`

【语法形态（事实描述）】
- 四行均**不使用** `APPROVED` 关键字，也**没有**「APPROVED + 括注四要素」形态。
- 实际形态 = `日期 + 主体名词 + 批准动词`（CN：`批准简化` / `批准` / `批准范围` / `明确批准`；EN：`simplified policy` / `approved` / `approved scope` / `explicitly approved`）+ 分号（CN 全角 `；`，EN 半角 `; `）+ 完成条件/义务子句。
- 均无加粗、无 `=` 号、无括注。

════════════════════════════════════════
四、ADR-014 ~ ADR-045 状态列以 `APPROVED` 开头的行
════════════════════════════════════════

【事实结果：0 条】
- 对 CN 与 EN 两文件 ADR-014 ~ ADR-045（各 32 行）的状态列做精确前缀扫描：以裸 `APPROVED` 开头 = 0 条；以 `**APPROVED` / `*APPROVED` 开头 = 0 条；大小写不敏感前缀 = 0 条。
- 按任务规则「不足 6 条则全列」：枚举清单为空（0 条）。因此「共同语法骨架归纳」对该形态不存在可归纳对象。
- 补充事实：对**全表** 47 行扫描，状态列含子串 `APPROVED` 的行仅 2 行，均在 ADR-014~045 区间之外：
  - ADR-010（两文件均 line 20）：状态列以反引号代码段 `` `PROPOSED / NOT APPROVED` `` 开头（首字符是反引号，非字母 A）。
  - ADR-012（两文件均 line 22）：状态列以反引号代码段开头，CN 逐字：`` `APPROVED / PARTIALLY IMPLEMENTED`；T026 已完成四类 wire/hash/replay、能力矩阵及 GitHub Copilot CLI 1.0.83 固定摘要探测。…``；EN 逐字：`` `APPROVED / PARTIALLY IMPLEMENTED`; T026 completed the four wire/hash/replay contracts, … ``（首字符同样是反引号）。

【补充事实：区间内实际存在的状态列开头形态（3 类，共 32 行）】
- 形态 A `**已裁（…）**。四要素：…`：21 行（ADR-014~029、031、033、036、044、045）
- 形态 B `**已登记（…）**。四要素：…`：10 行（ADR-032、034、035、037~043）
- 形态 C `` `ACCEPTED / IMPLEMENTED`（本片）。四要素：… ``：1 行（ADR-030，CN line 41 / EN line 42）
- 样本逐字摘录（开头片段）：
  - CN ADR-014（line 24）：`**已裁（用户直接裁定，仅本片有效，不得作为先例）**。四要素：裁决人 = 用户； 日期 = 2026-09-24；理由 = …`
  - CN ADR-024（line 34）：`**已裁（用户直接裁定；仅本片及以同款起跑消息起跑的片有效；用户手动切换时以用户裁定为准）**。四要素：裁决人 = 用户…`
  - CN ADR-027（line 37）：`**已裁（用户直接裁定；不得作为先例）**。四要素：裁决人 = 用户；日期 = 2026-10-19；理由 = …`（此行分号后**无空格**）
  - CN ADR-032（line 43）：`**已登记（步骤 0，2026-11-06）**。四要素：裁决人 = 用户；日期 = 2026-11-06；理由 = …`
  - EN ADR-014（line 25）：`**Ruled (direct user ruling, valid for this slice only, not a precedent)**. Four elements: ruler = the user; date = …`
  - EN ADR-032（line 44）：`**Registered (step 0, 2026-11-06)**. Four elements: ruler = the user; date = …`
- 骨架事实（对上述现存形态）：加粗 `**状态词（括注）**` + 全角句号 `。` + `四要素：` + `裁决人 = 用户；日期 = YYYY-MM-DD；理由 = …`；用半角 `=`（两侧带空格）；CN 用全角括号 `（）`/句号 `。`/分号 `；`，EN 用半角 `()`/`.`/`; `；CN 行内存在分号后有无空格两种变体（如 ADR-014 为 `； 日期`，ADR-027 为 `；日期`）。

════════════════════════════════════════
五、CN 与 EN 的 ADR-010 行号/同位核对
════════════════════════════════════════

- ADR-010：CN = `docs/ADR_REGISTER.md:20`，EN = `docs/ADR_REGISTER_EN.md:20` → **行号相同、结构同位**。
- 补充事实：ADR-001 ~ ADR-012 在两文件中占行号完全一致（lines 11–22）。分歧自 ADR-013 起：CN = line 23，EN = line 24；直接原因是 EN 第 23 行为**空行**（位于 ADR-012 行 line 22 与 ADR-013 行 line 24 之间），CN 第 23 行则直接是 ADR-013 行。此后所有行 EN 相对 CN 整体 +1（至 ADR-047：CN line 58 / EN line 59）。

BATCH-01: DONE

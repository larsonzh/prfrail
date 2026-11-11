# P66 CHECKLIST-RULES-HOIST · ⑤ 预审 r2（复审轮）

## 修复核验

- **F-1 | CLOSED**：`migration-mapping.md:118` 新增「节号重指记录」段，声明「目标位置」列保留规划态、实施按起跑消息 §3.2 按来源块 1:1 落点（§七→§13.1…§13.5），并逐字指向 `step-A-check.md:28` `### R2-1 结构分歧（已按规范收敛）`（该段 :41 反向引用「migration-mapping.md 的节号重指记录」，双向闭合）与 `step-B:26` / `step-C:19` / `step-D:19` / `step-E:25` 四处「偏离 ① 映射的登记」节——四处均真实存在；抽核 r1 所举 R04：实落 `13.1.4`（DELIVERY_DIRECTIVE.md:1221「模型与载体绑定（指针）」= R04 内容）与 1:1 方案一致。
- **F-2 | CLOSED**：`dedup-analysis.md:128` 新增附注，:130 给出总则句未计入 98/82 的理由并指向 §13.3.1 落点；`migration-mapping.md:112` 新增 `## 1.1 补登` G-0 行（§九 表前总则句 L206 → §13.3.1，原样，迁入）；`docs/DELIVERY_DIRECTIVE.md:1274` 实存 `**13.3.1 引用时用首句引文，不用「第 N 条」**`，与 dedup/step-D 描述的总则句逐字一致。
- 备注（非阻断）：补登段位置靠后致 `migration-mapping.md` 出现两个 `## 2.` 标题且 `## 1.1` 排在 `## 3` 之后，纯编号观感问题；门禁 `TOTAL_FAIL=0`，不构成新发现。**（主控已在 r2 后把该段标题改为「附：节号重指记录」，消除重号。）**

## 回归读数

- numstat：`docs/DELIVERY_DIRECTIVE.md` **184/3**、`docs/DELIVERY_DIRECTIVE_EN.md` **184/3** —— 与 r1 报告实测（CN 184/3 = EN 184/3）**逐字节一致**，即修复轮对两份规范文档零改动。
- porcelain：`M  docs/DELIVERY_DIRECTIVE{,_EN}.md`（仅首列，无第二列脏缓冲标记）；`A  …/01-arch/migration-mapping.md`（122/0）、`A  …/01-arch/dedup-analysis.md`（130/0）；`?? …/04-reviews/05-prereview-r1.md`（r1 报告，未暂存）。
- mtime 佐证：DELIVERY_DIRECTIVE{,_EN}.md = **16:41:50**（步骤 F 后未再动）；仅 migration-mapping.md / dedup-analysis.md（17:11:25）+ MANIFEST.md / SHA256SUMS.txt（resync 机械簿记，符合 G8 口径）落入 17:00 后窗口；step-B 未改（其原悬空引用因目标已存在而自动闭合，无需回改）。

## F-3/F-4/F-5 取舍

- **F-3：ACCEPT**。三列目标完好（§13.4 三列 26 行全非空），来源分组内容未丢失、落为表前块引用注记（DELIVERY_DIRECTIVE.md:1305），且 step-E 机械核验 #6 已如实记载实际落点——微差双向可见、无悬空引用；建议归档时在 step-E 偏离登记补一句「decision-list #5 括注方案改为表前注记」，非阻断。
- **F-4：ACCEPT**。坑 20「（候选纪律，P56 实证）」属**登记册口径**下的逐字事故原文（症状列），纪律 X 已另有规范源（13.2.23 指针 + §5.1）；与第 26 条的不对称源于后者升格发生在本片内并已按「已登记改写」声明；后续片可选在处置列加「现行见 §5.1」括注，本轮不动规范文档是正确取舍。
- **F-5：ACCEPT**。13.3.5 指针含 §6.3 有效支（「记录态」实体在 :688 等），§9.3（:901 起）确无字面但弱支不构成重复定义；13.3.16「环次见 §1.6」——全文档「环次」仅出现 1 次（即指针自身），§1.6 只有概念无字面；二者均为指针精度瑕疵、无规范语义风险，按 r1 建议移交后续片修（改 §1.6 字面或改指针措辞）。

## 未闭合清单

- 无。Critical/High/Medium 均已闭合；F-3/F-4/F-5 三项 Low 按上述判定 ACCEPT 留存，需在切片结案登记中列为「已知低风险留存」随档。

RE-REVIEW: PASS

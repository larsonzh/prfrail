BLIND REVIEW: FINDINGS

## 1) 审查范围与输入（file:line）
- 裁定原文：`docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/01-arch/ruling-2026-11-06-p60.md:12`、`:13`、`:14`、`:24`、`:41`。
- 起跑消息：`docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/00-briefing.md:7`、`:8`、`:136`、`:178`。
- 变更集：`git diff --cached`（按文件核验 `g1.js` / `g5.js` / `selftest.js`、`DELIVERY_DIRECTIVE{,_EN}.md`、`ADR_REGISTER{,_EN}.md`、`REMAINING_SLICES{,_EN}.md`、`p60-directive-judge-blindspots{,_EN}.md`）。
- 门禁复核：`node tools/gates/gate.js --check=G7 --all --scope=index`（本轮独立执行）。

## 2) 逐项结论

### 2.1 G7-f 闭环确认
- 通过（子项）：两份 P60 报告均含 1 个 `sixstage` 围栏，且 `Section A/B/C/D/E` 各 1 次，末行为 `INDEPENDENT SCAN: PASS`（CN: `docs/validation/p60-directive-judge-blindspots.md:96-123`；EN: `docs/validation/p60-directive-judge-blindspots_EN.md:96-123`）。
- 通过（门禁子项）：`G7-f sixstage transcript contract` 在本轮 `--scope=index` 读数为 PASS。
- 未通过（整条要求）：同一命令整体并非“全 PASS”，返回 `TOTAL_FAIL=2`（`G5(f)` 与 `G8-b` 失败，见第 4 节）。

### 2.2 J2 判据独立复核
- 一致：`G1(e2)` 已按“新增行判定、最近标题绑定、无标题/无状态跳过、历史不追溯”落地（`tools/gates/lib/criteria/g1.js:103`、`:131`、`:217-226`、`:238-242`）。
- 一致：`G5(e2)` 已按“块内完成记录字段行证据解析”落地，并补 `{,_EN}` 单组花括号展开与“可见跳过”明细（`tools/gates/lib/criteria/g5.js:99-111`、`:184-196`、`:277-290`）。
- 一致：`selftest` 已覆盖 T119-T124 的正反断言与计数细化（`tools/gates/selftest.js:2899`、`:2959`、`:3271`、`:3319`、`:3354`、`:3404`）。
- 核验通过：`G1(c)/(d)/(e)/(f)` 与 `G5(d)/(e)/(f)` 在 `git diff --cached --unified=0` 下未见改写；仅新增 e2 相关代码块。
- 核验通过：`evidencePathsFromRow` 字节级未改（old/new `sha256` 同为 `482c3a19e451f1eafaf54ea35da8ebb5555611db971293fc9861f83ebfc2dd18`，equal=true）。

### 2.3 J3 范围核验
- 通过：§12.4 的 OB-138 仅该行翻转至 `已处置` 并补注记（`docs/DELIVERY_DIRECTIVE.md:1171`；EN 镜像 `docs/DELIVERY_DIRECTIVE_EN.md:1171`）。
- 通过：`tools/gates/g6-whitelist.txt` 暂存 diff 为 0。
- 通过：`tools/gates/cjk-newwords.txt` 仅新增 1 行「稽」（附理由与上下文）。

### 2.4 授权可回溯
- 通过：头部版本与 §0.3 已到 v1.64（`docs/DELIVERY_DIRECTIVE.md:3`、`:111`；EN `docs/DELIVERY_DIRECTIVE_EN.md:3`、`:111`）。
- 通过：§0.3 v1.64 行明确“由用户 2026-11-06 P60 裁定（J2）授权”。
- 通过：ADR-034 已登记四要素（`docs/ADR_REGISTER.md:45`；EN `docs/ADR_REGISTER_EN.md:46`）。
- 通过：OB-138 处置注记与 J2/J1 口径一致（`docs/DELIVERY_DIRECTIVE.md:1171`；EN `docs/DELIVERY_DIRECTIVE_EN.md:1171`）。

### 2.5 D 项登记合规
- 通过：报告 §1 已登记 D-1 / D-2（CN `docs/validation/p60-directive-judge-blindspots.md:21-22`；EN `docs/validation/p60-directive-judge-blindspots_EN.md:21-22`）。
- 通过：台账 R/D 基线已同步 `D=2`（CN `docs/t027/REMAINING_SLICES.md:1721`；EN `docs/t027/REMAINING_SLICES_EN.md:1708`）。
- 发现：同一 P60 行的“完成记录”叙述末尾仍写 `R=0、D=0` / `R=0, D=0`（CN `docs/t027/REMAINING_SLICES.md:1723`；EN `docs/t027/REMAINING_SLICES_EN.md:1710`），与基线字段不一致。

### 2.6 编码与对称
- 通过：指定 `.md` 文件均为 UTF-8 with BOM + LF；指定 `.js` 均为 UTF-8 without BOM + LF（逐文件脚本核验）。
- 通过：`DELIVERY_DIRECTIVE` CN/EN 六形状门禁为 PASS（`G3(a.2)`，本轮 gate 输出）。

### 2.7 独立新增问题结论
- 本轮新增 Medium 级问题 2 项（见第 4 节）。
- 本轮新增 Low 观察 1 项（R/D 叙述与基线字段不一致）。

## 3) 发现清单

1. [Medium] `--check=G7 --all --scope=index` 非全 PASS，`G5(f)` 失败（index hygiene）。
   - 证据：本轮门禁输出 `FAIL  G5(f) index hygiene`，并列出 `MM/AM` 路径（`docs/t027/REMAINING_SLICES{,_EN}.md`、`docs/validation/p60-directive-judge-blindspots{,_EN}.md`、`.../freedom-list.md`）。
   - 影响：不满足“工作区 == 暂存区”的闭环前提，且不满足任务 1 的“应全 PASS”。
   - 修法：先消除 unstaged 与 staged 偏差（确保 index/worktree 一致），再重跑同一门禁命令并记录新读数。

2. [Medium] 证据包自洽失败，`G8-b` 报哈希不一致。
   - 证据：本轮门禁输出 `FAIL  G8-b pack self-consistency`，指向 `docs/validation/evidence/P60-DIRECTIVE-JUDGE-BLINDSPOTS/06-gates/gate-readings.txt` 与 `.../freedom-list.md` 的 recorded/actual hash mismatch。
   - 影响：冻结包不可验，⑨/⑧c 证据链不闭合。
   - 修法：以最终文件内容重算并回写 `SHA256SUMS.txt`（必要时连同 `MANIFEST.md` 同步），确保包内载荷集合与哈希一致后重跑门禁。

3. [Low] P60 台账行内口径冲突：`R / D 基线` 已为 D=2，但“完成记录”叙述仍写 D=0。
   - 证据：CN `docs/t027/REMAINING_SLICES.md:1721` vs `:1723`；EN `docs/t027/REMAINING_SLICES_EN.md:1708` vs `:1710`。
   - 建议：在同一行叙述口径统一为 D=2，避免读者误判。

## 4) 声明
- 本轮为独立盲审；除本回执文件外，未改任何文件；未 commit/push；未派单。
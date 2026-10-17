# DIRECTIVE-TEMP-AUTH-LIMIT 切片验证报告（临时授权阈值与额度上限修订：§1.7.4 ⑤ + §7.5 + OB-92 + OB-93）

> 切片：`DIRECTIVE-TEMP-AUTH-LIMIT`（v1.35）。状态：② 实现完成、门禁全绿；⑦ r1 `PASS`；停在 ⑩，提交须同轮授权。

## 1. 变更概要

- **§1.7.4 ⑤ 增四项判定要素**：(i) 「同类」= 同一被突破条款 + 同一例外类别（五类闭集）；(ii) 计数以 §12.4 台账与 ADR 登记表为唯一真值源（同一实例只计一次），窗口自该条款最近一次正式修订起算且修订后归零；(iii) 第三次不得授权、步骤 0 停机；(iv) 同款起跑消息不延续前片授权、须同轮新裁定并计入该族累计计数（OB-92 落地句）。
- **§7.5 ⑦ 额度按档分层**：S 1+1、M 1+2、L·治理·高风险 1+3，盲审另计 +1；新增「分档判定规则」（本表「档」= §1.2 规模档、`M-L` 视同 M、多命中取复审最多者、治理片与文档型定义、在飞切片过渡）与「修订依据」（ADR-025 的 7 次实例；最新实例 `GATES-TIGHTEN` 的 ⑦ r1–r3）。
- **§5.1 两条时序规范**：审查输入包生成器须逐件断言申报输入存在并附输入清单（OB-93）；⑧c 清理须在记录落盘之后、逐件精确、可机械回读，临时件一律 `tmp/<切片ID>/`。
- **台账与附录**：§12.4 的 OB-92 / OB-93 翻 `已处置` 并新增 OB-94（⑤ 额度用尽后的残差登记，⑦ 复核后翻 `已处置`）；§12.5 三行缺口启动时登记、本片落地后翻 `已落地`；附录 C.0w；ADR-024 处置列追加 OB-92 修订注记（**不改写**历史裁定原文）。
- **边界**：未改 §2.2 / §2.4 闭集、§5.2 矩阵 / §6.3 判据本体、CONTRACTS / schema / fixtures / workflows 与 `internal/**`；**未新增 ADR**（上限仍 25，步骤 0 判定修订走切片路径即可）。

## 2. 单元实施与收口读数

| 单元 | 内容 | 收口读数 |
|---|---|---|
| 1 | §1.7.4 ⑤ 四要素 + OB-92 落地句 | 门禁全绿 |
| 2 | §7.5 分层额度 + 分档判定规则 + 修订依据 | 门禁全绿 |
| 3 | §5.1 两条时序规范 + OB-93 | 门禁全绿 |
| 4 | 台账（OB-92/93/94、§12.5 三行、C.0w、ADR 注记、P10） | 门禁全绿 |

- **末次全量运行（`--scope=index`，本报告对定稿快照）**：`SELFTEST PASS 417/417`；`TOTAL_FAIL=0`；`G1-b` 两文件 `header=v1.35 last=v1.35`；`G3(a)` 三对 numstat 全对称且六形状全 `true`；`G4b unknown=[]`。
- **新增自测**：无（纯文档片，自测分母不变，仍 417/417）。
- **⑨ 定稿读数**：全树编码 `ENC-TREE tracked=3617 text=3528 binary=89 violations=0`（`docs/validation/evidence/` 的 17 件冻结样例按例豁免）；变更行 `§X.Y` 引用 `total=277 unresolved=0`；`go build ./...` 与 `go vet ./...` 通过、`go test ./...` `ok=14 FAIL=0`（由 base 门禁执行）。

## 3. 片内自查自纠（不计偏离）

门禁当场红 ⇒ 就地修 ⇒ 不计偏离，逐条如下：

- `G1-b`：v1.35 日志行初版插在 §0.3 表首 ⇒ header 与末行版本不一致 ⇒ 移到表尾（末行 = 当前版本）。
- `G4b`：三次命中基线语料之外的字 ⇒ 分别改写为「让……只能落到」「三条约束」「ADR 登记表」，**未扩白名单**。
- `G6-4` 误报：OB-94 的记录把版本号后接「行」字的写法当成过期的行数表述 ⇒ 改为「条目」。
- `G3(a.1)`：`P10` 块尾空行使 CN/EN 的 `numstat` 变成 17/18 ⇒ 对齐为 18/18。
- 自检脚本 `checks.js` 两处断言把通过误判为失败（A4 用固定长度切片、A5 计数口径含不相关行）⇒ 已修正，现均正确通过。
- ADR-024 变更行的 EN 侧加粗对数多 2 属**既有漂移**（HEAD 即 CN 20 / EN 22；本片双侧对称各 +2，未扩大漂移）⇒ **不改**，仅登记。

## 4. 未决项与偏离台账

- 同型偏离：**2/3**（本片无新增）。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（4 Medium + 5 Low）⇒ 逐条原位整改。
- **⑤ r2**：`PRE-REVIEW: PASS WITH FIXES`（2 Medium 残差 + 4 Low）⇒ 全部就地整改；⑤ 额度 1+1 用尽 ⇒ 按 §7.8.3 并入 ⑦ 输入并登记 OB-94，**未开第 3 轮**。
- **⑥**：`INDEPENDENT SCAN: PASS`（盲输入，Section B 空）。
- **⑦ r1**：`FINAL REVIEW: PASS`（0 发现；F-1 / F-2 闭合；⑦ 不降级，`GPT-5.3-Codex` Extra High）。
- **自指情形登记（本片特性，必须显式）**：本片在修订 §1.7.4 ⑤ 阈值的过程中**自身撞上该阈值**；按 §7.8.3「额度用尽」行接受「fix + 并入 ⑦ 输入」路径，F-1 / F-2 的闭合路径交 ⑦ 复核，**OB-94 为登记载体**；**不销解、不追加 ⑤ 轮次**。
- **§12.1 试点顺延（第 8 次）**：本片为准则治理 / 规则修订片，不满足 §12.1 的试点选择标准（命中配套②）；按 §12.4 唯一真值源记为第 **8** 次（v1.30–v1.34 五片同例），本片**未新增顺延 OB**；如需正式登记由用户裁定。
- ⑩ 停点：提交与推送须同轮授权；本片尚无任何提交。
- **⑧c 两次执行（本片自指情形的正向实证）**：首轮清理后仍须把清理记录落盘 ⇒ 按本片新写入 §5.1 的时序规范「先落痕再清理」重开临时目录落痕、末尾再跑一次 ⑧c（`Test-Path` 为 `False`、`tmp/` 仅 `.gitkeep`）——**新规则在本片收尾中被实际使用并生效**。
- **同型第 3 次：`$env:TEMP` 复发（形态偏离，不销解、不加新 OB）**：CI 观察阶段用 `$env:TEMP\w.txt` 承接 `gh run watch` 输出；按起跑消息第 7 条**停下升级**；已立即 `Remove-Item -LiteralPath` 删除并回读 `Test-Path` 为 `False`；未入仓库、不影响门禁与 CI；**升格为下一片起的前置约束与机械判据候选**。
- **⑧c 清理（已执行）**：`tmp/DIRECTIVE-TEMP-AUTH-LIMIT/` 的全部辅助件（脚本、审查输入包、门禁原始读数）逐项删除；回读 `Test-Path` 为 `False`、`tmp/` 仅余 `.gitkeep`。

## 5. 已知缺口与登记（裁定所称「§8」）

- **OB-92 / OB-93 / OB-94 均翻 `已处置`**：OB-92 以「⑤(iv) 落地句 + ADR 修订注记 + 不改写历史裁定」三件套闭合；OB-93 以 §5.1 两条时序规范闭合；OB-94 由 ⑦ r1 复核 F-1 / F-2 后闭合。
- **F-3..F-5（Low）的整改登记**：F-3「本节末句」→「本节表后句」；F-4 §12.5 的 ⑦ 预算行补「v1.35 起按该片档位取值」；F-5 步骤 0 锚点补「§1.2『步骤』字段首项；未写该项的切片须在起跑回执中补做」。
- **本片范围内仍未闭合项**：无。
- **候选（只在报告登记，不占 OB 号）**：把「临时件纪律」从授权边界升格为**机械判据**（如 `G9-c` 或并入 `G6` / `G8`）：**⑧c 回读证据 + 临时件目录形态检查**；本片不落地，作为后续治理片的输入。
- **候选（只在报告登记，不进 §12.4）**：本片新增的 §7.5 分层额度尚无真实切片验证（n=0）⇒ 建议下一片按新表执行并记录实际用量；若 M 档无需用满即 `PASS`，可在后续治理片把 M 档收回 1+1。

## 6. 下一步建议（待用户裁定）

- **建议下一片用新额度表实跑一次**（M 档或 L 档 / 治理片），验证「1+2 / 1+3 + 收窄输入」是否够用，并把实际用量写入该片报告的「成本与计量」节（R7.1 口径）。
- **§12.1 试点**仍待用户裁定（本片为治理片，不满足试点标准）。

## 7. 提交与 CI

- **主提交**：`9d5bf9e` — `feat: revise the temporary-authorization threshold and the review quota cap (v1.35)`（9 文件，`+286/−16`）；已推 `origin/main`（SSH，范围 `aa49e0c..9d5bf9e`）；**未推 gitee**。
- **CI 观察（该提交自身的 run）**：run `39045058091`（workflow `CI`，event push，attempt 1）**首跑双腿全绿**——`Go ubuntu-latest` 2m39s、`Go windows-latest` 3m00s；**未触发首跑红处置**。
- **Gates 步硬证据（两腿实测）**：`SELFTEST: PASS 417/417`；`GATE REPORT scope=ci changed=9`；`TOTAL_FAIL=0`；`G3(a)` 三对 numstat 全对称、指令对六形状全 `true`。
- **包级证据**：两腿合计覆盖 14 个包（`internal/**` 的 12 个、`examples/whois-shadow`、`tools/agent-probe/enforcement-proxy`），`FAIL` 行 0。
- **回填口径**：本节由 docs-only 回填写入；该回填提交自身的 CI 只作观察，不再为其新增提交。

## 8. 工件

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-temp-auth-limit.md	repo
docs/validation/directive-temp-auth-limit_EN.md	repo
docs/validation/evidence/DIRECTIVE-TEMP-AUTH-LIMIT-freedom-list.md	repo
```

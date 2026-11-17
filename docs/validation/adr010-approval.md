# P77 验证报告：ADR-010 批准与 EOL 通知方式同步

日期：2026-11-17。状态：`[SLICE]` 文档类，S 档。切片标识：`ADR010-APPROVAL`（P77）。

## 1. 切片标识与授权

- 目标：产品所有者批准 ADR-010（EOL 通知方式）后执行文档落地——ADR-010 状态翻 `APPROVED`，并把生效的 EOL 通知方式同步到所有相关文档。
- 批准原文四要素：裁决人 = 用户（产品所有者）；日期 = 2026-11-16；理由 = v0.1.0 候选版进入 Release 流程前须冻结 EOL 通知方式以支持 S-G4 收口；依据 = 用户 2026-11-16 P77 起跑消息同轮裁定。
- 基线：`HEAD` = `origin/main` = `fb90a20`；准则 v1.69；`SELFTEST` 737/737；三作用域 `TOTAL_FAIL=0`。
- R / D 基线：R = 0、D = 0（步骤 0 登记，§13.1.1）。

## 2. 边界与不改清单确认

- 未改 S1 Exit 八门判据；`p57-closeout-assessment{,_EN}` 门表与 `s1-exit{,_EN}` 门表零改动。
- `s1-exit{,_EN}` 的剩余风险第 1 条属叙事段（非门表），按 P76 报告 §8 的「批准后改动点」清单同步；同文件日期为 2026-09-09 的结论段按 §11.1 历史不改写保留（见 §8）。
- 未改 EOL 通知方式本身：ADR-010 的第 2 列与第 3 列逐字未动。
- 未改 `S1_SUPPORT_MATRIX` 平台清单与支持层级；未改 INSTALLATION 安装模型。
- 无 GitHub Release notes 外部动作（属 P78 `RELEASE-ZIP-BINDING`）。
- 未改其它 ADR；未新增 ADR，未新增 OB。
- 未改准则版本（维持 v1.69）；未触 `CHECKLIST_v59/v60{,_EN}` 与 v50–v58 快照。
- 未改 `.go` / `.json` / `.yml` / schema / fixtures / CI / `go.mod`；未推 gitee。

## 3. 文档改动逐项证据

| 文件 | 位置 | 改动性质 | CN 增删 | EN 增删 |
|---|---|---|---|---|
| `docs/ADR_REGISTER.md` / `_EN` | ADR-010 状态列 | 替换 | `+3/-3` | `+3/-3` |
| `docs/ADR_REGISTER.md` / `_EN` | 支持策略段尾句 | 替换 | 计入上行 | 计入上行 |
| `docs/ADR_REGISTER.md` / `_EN` | 决策单「平台/支持」状态列 | 追加闭括号 | 计入上行 | 计入上行 |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | 尾部政策段首句与尾句 | 替换 | `+2/-2` | `+2/-2` |
| `docs/INSTALLATION_PLAN.md` / `_EN` | EOL 未冻结并列句 | 替换 | `+1/-1` | `+1/-1` |
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | §7 待落实句 | 追加加法式指针 | `+1/-1` | `+1/-1` |
| `docs/validation/s1-exit.md` / `_EN` | 剩余风险第 1 条 | 替换 | `+1/-1` | `+1/-1` |
| `docs/t027/REMAINING_SLICES.md` / `_EN` | P77 台账段 | 追加 | `+20/-0` | `+20/-0` |

### 3.1 ADR-010 状态栏改前与改后

改前（CN）：

> `` `PROPOSED / NOT APPROVED`；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准 ``

改后（CN）：

> `` `APPROVED`（裁决人 = 用户；日期 = 2026-11-16；理由 = v0.1.0 进入 Release 流程前冻结 EOL 通知方式；依据 = 用户 2026-11-16 P77 起跑消息同轮裁定原文）；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准（**已于 2026-11-16 落实；P77**） ``

EN 侧改动同构：状态词取 `APPROVED`，四要素改为英文同一句型，历史完成条件句逐字保留，尾部闭括号与 CN 同构。

### 3.2 六份候选文档与两处同源表述的判定

| 文档 | 判定 | 依据 |
|---|---|---|
| `S1_SUPPORT_MATRIX{,_EN}` | 第一层，必改（两处） | 含「待审批」「尚未获批」的未批准态表述 |
| `INSTALLATION_PLAN{,_EN}` | 第一层，必改 | 该句把 EOL 通知与另两项并列称为「仍待冻结」，属未冻结态 |
| `PRODUCT_REQUIREMENTS{,_EN}` | 第二层，触发，加法式指针 | 待落实清单含 EOL 通知项 |
| `S1_RELEASE_NOTES{,_EN}` | 不触（由第二层降级） | 该行是绑定义务，不声称未批准；批准不等于绑定完成 |
| `INSTALLATION{,_EN}` | 不触 | 全文无 EOL 表述 |
| `OPERATIONS{,_EN}` | 不触 | 全文无 EOL 表述 |
| `ADR_REGISTER` 支持策略段尾句 | 第一层，必改 | 与批准事实直接矛盾的存量断言 |
| `ADR_REGISTER` 决策单平台/支持格 | 第一层，追加闭括号 | 该格是「当前约束待落实」状态列 |
| `DEV_PLAN{,_EN}` | 不改（边界排除，如实披露） | 三处均为带日期的历史记录，且该文件在边界中明令不改 |

### 3.3 组合检索排查

- 五类模式各自独立检索，并集作用域为证据树以外（`docs/validation/evidence/**` 内的命中属本片自身的转录与证据载荷，不计入波及面）；并集共涉及二十三个文件。
- 归属四类：本片改动、本片登记（台账）、不改（历史冻结）、不改（边界排除并披露）。逐文件归属见 ① 方案书 §Q4。
- 冻结类文件（`CHECKLIST_v59/v60{,_EN}`、`p57-closeout-assessment{,_EN}`、`s1-exit-closeout-assessment{,_EN}`、`s1-reach-reachability{,_EN}`）中的时点表述保留，靠 ADR-010 状态列溯源。

### 3.4 ② 落盘形态与 ① 目标文本的差异披露

| 位置 | ① 方案书目标文本 | ② 实际落盘 | 取舍依据 |
|---|---|---|---|
| ADR-010 状态列 | 加粗状态词 + 「四要素：」前缀 + 性质 / 可引用 / 生效日 / 回退尾注 | 起跑消息 §四 ②a 参照格式：反引号状态词 + 括注四要素 + 历史句 + 加粗闭括号 | 起跑消息 §四 ②a 是 ② 的权威格式；性质等尾注不入状态列，S-G4 收口结论由本报告 §5 承载、四要素与日期见 §1，避免与 ADR-045 / ADR-046 / ADR-047 重复 |
| `S1_SUPPORT_MATRIX` 首句与尾句 | 首句为「已批准的 ADR-010（2026-11-16）：」；尾句保留「本页不能替代最终公告」 | 首句为加粗的「ADR-010（已批准，2026-11-16）」；尾句改为「本页不替代 Release notes 的正式公告」 | 起跑消息 §四 ②a 的句 1 与句 2 明示该形态；通用限制保留，只消除「尚未获批」层 |
| `PRODUCT_REQUIREMENTS` §7 待落实句 | 句内最小切分 | 句尾加法式指针 | 起跑消息 §四 ②b 第 4 项明示「加法式指针」，且加法式不改写既有句子（§11.1） |
| `s1-exit` 日期为 2026-09-09 的结论段 | 列为改动点 | 未改 | 该段是带日期的历史结论，§11.1 不改写；P76 报告 §8 的最小清单只列剩余风险第 1 条 |

- 说明：① 方案书 §Q3 的「状态列已按 Q1 目标文本落盘」是写作时点的前瞻表述；② 实际按上表落盘，该表述由本节取代。

## 4. 验收证据与门禁

| 判据 | 结果 |
|---|---|
| ADR-010 状态列含首状态词 `APPROVED` | 是（CN 与 EN 同为真） |
| 四要素逐项存在 | 是（裁决人 / 日期 / 理由 / 依据） |
| 历史完成条件句逐字保留并追加闭括号 | 是（CN 与 EN 各一组断言） |
| `S1_SUPPORT_MATRIX` 的「该提议尚未获批」零命中 | 是（CN 与 EN 均零命中） |
| 改动行数 CN 与 EN 对称 | 是（每对文件的两侧增删一致） |
| 三作用域与 `SELFTEST` | 见 §7.2 |

## 5. S-G4 收口

判据原文（不改）：`p57-closeout-assessment{,_EN}` §4.2 的 S-G4 行为「EOL 通知方式 | ADR-010 获产品所有者批准」，本片对该行及其所在门表零改动。

- 证据：`docs/ADR_REGISTER.md:20` 的状态列已按 §3.1 的改后文本落盘（EN 侧为 `docs/ADR_REGISTER_EN.md:20`）；本报告与本报告对承载结论；批准四要素见 §1。
- 结论：**S-G4 由「未通过」转为「通过」**；该判定由本报告承载，`p57-closeout-assessment{,_EN}` 与 `s1-exit{,_EN}` 的门表据 P76 裁定保持零改动。
- 下游：P78 `RELEASE-ZIP-BINDING` 解锁（S-G4 为其依赖项之一）；S-G3 待 P78 完成后判定。
- 不变项：本次批准不得被解读为 S1 Exit 通过；S-G5 仍为人工门；`s1-exit` 的核心任务链口径与 S-G2 句留后续切片同步。

## 6. 下游解锁与残留

- 与 ADR-046 的 Q10 = (C) 分两步保持一致：本片完成第一步的 ADR-010 批准部分；第一步的 v0.1.0 Release 部分由 P78 承接；第二步（S1 Exit 整体批准）在 S2 之后。
- `S1_SUPPORT_MATRIX` 尾部政策段的其余三条要点（版本三段式、九十天窗口与评审、偏离须有带日期 ADR）逐字未动，与 ADR-010 第 2 列与第 3 列一致。

## 7. 计量与审查链

| 环节 | 次数 |
|---|---|
| ① 架构 | 1 次（`deep-reasoner` + `DeepSeek V4 Pro`） |
| 批量辅助 | 2 次（`GLM-5.3-Flash`） |
| 探针 | 0 |

### 7.1 审查链结果

- ① 架构：`DESIGN: DONE`，方案书见证据包 `01-arch/design-r1.md`。
- ⑥ 跳过：按 §7.2.1 三条件成立（改动只落文档与台账类文件；不含代码、schema 与 fixtures 变更；不涉状态机、门禁判据与授权契约语义），由 ⑦ 确认。
- ⑤ 预审：2 轮——r1 `PRE-REVIEW: PASS WITH FIXES`（1 中 3 低）；r2 `PRE-REVIEW: PASS WITH FIXES`（0 中 3 低）；额度 2/2 用尽，r2 的 3 条低风险整改交 ⑦ 代为确认。
- ⑦ 终审：1 轮——`FINAL REVIEW: PASS WITH FIXES`（0 中及以上，1 低；含 `SKIP CONFIRMED`）；额度 1/2 使用。
- ⑥ 跳过确认（由 ⑦ 作出）：覆盖面 = 变更面边界与文件类型、双语对称与报告形状、机械门禁收敛、S-G4 收口承载链；遗漏面 = 推送后真实 CI 运行态、用户原始起跑消息的字节级对照、运行时行为验证（本片不涉代码执行语义）。

### 7.2 门禁读数

| 判据 | 读数 |
|---|---|
| `--scope=index` | `TOTAL_FAIL=0` |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0`（本地读数取自上一提交的变更集；推送后的真实读数由 ⑩ 回执载明） |
| `--selftest` | `SELFTEST: PASS 737/737` |

## 8. 已知限制与残留差异

- `p57-closeout-assessment{,_EN}` §4.2 的 S-G4 「现状与证据」格仍引用旧状态文本；按 P76 裁定门表零改动，读者须以 `ADR_REGISTER` 的状态列为准。
- `s1-exit{,_EN}` 的日期为 2026-09-09 的结论段仍写「EOL 通知方式尚未冻结」；该段是带日期的历史结论，按 §11.1 不改写。
- `DEV_PLAN{,_EN}` 三处历史记录仍把 EOL 通知方式列为 T018 的剩余门；该文件在边界中明令不改，且三处均为带日期的进度记录。
- 冻结快照类文件（`CHECKLIST_v59/v60{,_EN}`、`s1-reach-reachability{,_EN}`）中的「未通过 / 未批准」时点表述按 §3.3 口径保留，不追溯改写。

## 9. ⑩ 回执要素

- 终态读数、提交号与 CI 运行号见 ⑩ 回执（本报告按 §11.2 不写入提交号与运行号）。
- 给下一片的锚点：P78 `RELEASE-ZIP-BINDING` 以本片闭环后的 `HEAD` 与 CI 运行号为步骤 0 锚点。

## 10. CI 事实链

本节由 CI 双绿后的 docs-only 回填提交补入（§11.2）；该回填提交自身的 run 只作观察记录。

- 主提交 `33272be`（`docs: approve ADR-010 and sync the EOL notification policy (P77)`）；推送 `origin/main`，范围 `fb90a20..33272be`（未推 gitee）。
- CI run `43986942087` 首跑双腿 success：ubuntu `3m22s`、windows `4m32s`；两腿 `SELFTEST: PASS 737/737`、`GATE REPORT scope=ci changed=29`、`TOTAL_FAIL=0`；包级 `ok` ubuntu 21 / windows 15；base `go test ok=15 FAIL=0`。
## artifacts

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/S1_SUPPORT_MATRIX.md	repo
docs/S1_SUPPORT_MATRIX_EN.md	repo
docs/INSTALLATION_PLAN.md	repo
docs/INSTALLATION_PLAN_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/adr010-approval.md	repo
docs/validation/adr010-approval_EN.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/00-briefing.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/01-arch/design-r1.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/02-batch/batch-01-adr010-register-format.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/02-batch/batch-02-eol-survey.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/04-reviews/final-r1.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-ci.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-index.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/06-gates/gate-tree.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/MANIFEST.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/SHA256SUMS.txt	repo
docs/validation/evidence/P77-ADR010-APPROVAL/freedom-list.md	repo
docs/validation/evidence/P77-ADR010-APPROVAL/input-list.md	repo
```

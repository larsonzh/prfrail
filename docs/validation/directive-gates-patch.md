# 切片 `GATES-PATCH` 验证报告（准则 v1.40）

> 切片：`GATES-PATCH`；准则版本 v1.40；档位 L；本文件为 **⑧a 初稿**（⑨ 与 ⑩ 的读数由 ⑧b 定义稿补入）。
> 基线 HEAD = `2482971`（⑧a 启动时读取）；索引与工作区的一致性由主控在 ⑧b 复核。
> 双语：本文件为 CN 权威，`docs/validation/directive-gates-patch_EN.md` 为逐节同位镜像；两侧行数相等。

## 1 变更概要

⑧a 之前的变更集为 13 项路径，⑧a 之后为 16 项路径；其中 `tools/**` 下 5 个判据文件为代码改动，其余为文档与台账。

| 路径 | + | − | 说明 |
|---|---|---|---|
| `docs/ADR_REGISTER.md` | 1 | 0 | ADR-027 登记行 |
| `docs/ADR_REGISTER_EN.md` | 1 | 0 | `_EN` 同位 |
| `docs/DELIVERY_DIRECTIVE.md` | 43 | 21 | v1.40 准则正文 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 43 | 21 | `_EN` 严格同位 |
| `docs/DEV_PLAN.md` | 1 | 1 | 过期字面量的措辞调整 |
| `docs/DEV_PLAN_EN.md` | 1 | 1 | `_EN` 同位 |
| `docs/t027/REMAINING_SLICES.md` | 2 | 2 | P2 状态行翻完成与残差转出说明 |
| `docs/t027/REMAINING_SLICES_EN.md` | 2 | 2 | `_EN` 同位 |
| `tools/gates/gate.js` | 2 | 0 | 注册 G5-d |
| `tools/gates/lib/criteria/g5d.js` | 143 | 0 | 新增 G5-d |
| `tools/gates/lib/criteria/g6.js` | 203 | 13 | 续延集补 `:` 与 G6-7 |
| `tools/gates/lib/criteria/g9.js` | 87 | 6 | G9-e |
| `tools/gates/selftest.js` | 494 | 10 | 断言 552 → 588 |

**artifacts 围栏**（16 项路径，处置均为 `repo`）：

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-patch.md	repo
docs/validation/directive-gates-patch_EN.md	repo
docs/validation/evidence/GATES-PATCH-freedom-list.md	repo
tools/gates/gate.js	repo
tools/gates/lib/criteria/g5d.js	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/lib/criteria/g9.js	repo
tools/gates/selftest.js	repo
```

## 2 执行流水线

本轮由主控串行派单；各环节的实际读数如下。

| 环节 | 角色 / 载体 | 模型 | 结论 | 轮次 | 备注 |
|---|---|---|---|---|---|
| ① 设计 | 主控与 `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 设计完成 | 1 | 四组编号与 SELFTEST 增量表 |
| ② 实现 | `prfrail-implementer` | 主控模型 | 完成 | 3 | 判据代码、文档与两轮整改 |
| ③ 测试 | `prfrail-tester` | 主控模型 | 完成 | 1 | 三作用域读数与 5 条 Low |
| ⑤ 预审 r1 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 2 Medium + 4 Low |
| ② 整改 | `prfrail-implementer` | 主控模型 | 完成 | 1 | 两条 Medium 闭合 + 2 条断言 + 6 处数字同步 |
| ⑤ 预审 r2 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS | 1（收窄） | 两条 Medium 均 CLOSED，零新发现 |
| ⑥ 独立盲审 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | PASS | 1 | 盲输入，零 Medium+ |
| ⑦ 终审 r1 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | 五组逐组结论均通过，零 Medium+ |
| ⑦ 复审 r2 – r3 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | r2 FINDINGS（1 Medium）→ r3 PASS | 2（收窄） | 自由度清单漏登 §5.1 的 `*.raw.txt` 优先级次序；整改 = 清单追加 F12 |
| ⑧a 初稿 | `prfrail-documenter` | 主控指定 | 本文件 | 1 | 三份文件落盘 |

- 返工口径：② 共 5 次派单（步骤 0 落地、实施 1、实施 2、整改 M-1/M-2、断言数同步），本表按环节拆为「实现 3」与「整改 1」两行，另加步骤 0 落地 1 次；⑤ 共 2 次（配额 1+1 用尽）；⑦ 共 3 次（r1 终审 + r2 收窄 + r3 收窄）。
- 耗时字段历史缺失：本片未落盘起止时间戳，不补造。
- ⑦ r2 / r3 的处置口径见 §5（⑧b 已回填）。

## 3 异常与兜底记录

| # | 事件 | 影响 | 处置 |
|---|---|---|---|
| D1 | 步骤 0 的两处机械归一（§12.5 的一行改写、§5.1 的空行归一） | 文本形态变化，语义不变 | 已处置 |
| D2 | ① 预计新增断言 40 条（→592），实际 34 条（→586），r2 补 2 条（→588） | 差额 6 条来自 T88 前置门由 9 条降到 3 条；⑦ r1 只核对了「静态新增 36 处 = 运行时增益 36」的一致性，**未逐条核对降级是否合理** ⇒ **登记为残余** | 记录在案 |
| D3 | ③ 的 5 条 Low（D-1 至 D-5） | 无阻断项 | 记录在案 |
| D4 | 两份候选形状（第七形状「CJK 字符间空格」、第八形状「变更行形状多重集」）**不做判据实现**，仅登记为记录态 | 实测 **13/50** 触及指令对的历史提交存在 CN≠EN 漂移（26%，属普遍形态而非边缘形态）；CN 侧第七形状实测 16 处命中、逐处判定为有意形态。原「抽样 4 个 / 2 个不可满足」为**错误估计**（把普遍形态误判为边缘形态），已作废、不得再引用 | 记录在案（**OB-104**） |
| D5 | 主控曾把 ⑤ 包构建脚本写到 `tmp/` 根一次 | 随即移入 `tmp/GATES-PATCH/`；属 OB-100 第二处实证 | 已处置 |
| D6 | ② 整改后误删整个 `tmp/GATES-PATCH/`（同型偏离第 2 次；起跑消息第 7 条：第 3 次停下升级） | 被删内容 = ⑤ r1 输入包（11 项 + `MANIFEST.md`）+ 历史重放脚本 + 日志（合计 28 个条目）。**⑤ r1 的原始输入包不可逐字节重建**——其 `02-changes.diff` / `03-numstat.txt` / `05-selftest-tail.txt` 取自整改前的暂存树，该树从未提交、已不存在；⑤ r1 的结论仅存于 transcript（字节偏移 `23604898`、行 `25213`）。**⑦ r1 用的是整改后新建的 `review-6/` 包**（含 `MANIFEST.md`），该包在 ⑦ r1 **之前**生成，不是 ⑤ r1 包的副本 | 已整改（不可追溯部分如实登记） |
| D7 | ② 首轮整改用渲染文本作 `oldString` 未命中 | 改用基于暂存索引的确定性脚本重写整行 | 已处置 |
| D8 | ⑤ 配额 1+1 用尽 | r2 为收窄原位重跑，未开第 3 次 | 记录在案 |
| D9 | 主控首次构建 r2 与 ⑥ 包时把三次 `--all` 内嵌进脚本 | 后台挂起并被终止；改为先落盘日志再轻量装配，读数未受影响 | 已处置 |
| D10 | 自由度清单产出时点晚于 ⑦ r1（§1.7 要求 ⑥ 之后、⑦ 之前） | ⑦ r1 输入包内无清单 | 记录在案 |
| D11 | `DECL_ID_RE` 的边界断言只加了前侧（`OB-97x` 仍截出 `OB-97`） | 两侧共用同一常量，对称且不产生假绿 | 记录在案 |

- D6 的偏离型态已是第 2 次，前例见 `docs/validation/evidence/G9-C-freedom-list.md` 第 4 节的第 2 条。
- D4 的降级依据已按实测口径改写：漂移是普遍形态（13/50），故两项形状在本片只登记、不实现；「保守」不再是理由。

## 4 环境与档位

| 项 | 值 |
|---|---|
| OS | Windows（本机 `d:\LZProjects\prfrail`） |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| 基线 HEAD | `2482971` |
| 档位 | L |
| ⑥ | 必派（判据本体变更，§7.2.1 第 ③ 条不成立） |
| ⑦ 额度 | 1+3 |
| ⑤ 额度 | 1+1（用尽） |
| 批量辅助 | ADR-027（`GLM-5.3-FlashX (glm)`，限四类任务） |
| 探针 | 0 |

## 5 审查结论摘要（⑤⑥⑦）

| 轮次 | 编号 | 级别 | 结论 / 处置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 | M1 | Medium | `DECL_ID_RE` 缺前置词边界，`ADR-027` 被截为 `DR-027` | 已整改 |
| ⑤ r1 | M2 | Medium | 「脚本化状态」段漏列 G5-d / G6-7 / G9-e | 已整改 |
| ⑤ r1 | L1–L4 | Low | EN 动词闭集、`翻` 词边界、`TMP_REF_RE` 边界、两条补登行未入枚举 | 记录在案 |
| ⑤ r2 | — | — | 两条 Medium 均 CLOSED，零新发现 | 已闭合 |
| ⑥ r1 | — | — | 盲审零 Medium+ | 记录在案 |
| ⑦ r1 | — | — | 五组逐组结论均通过，零 Medium+ | 记录在案 |
| ⑦ r2 | — | Medium | 收窄复审自由度清单：漏登 v1.40 新增的 §5.1「冻结证据包内 `*.raw.txt` 的优先级」次序 | 已整改（清单追加 F12） |
| ⑦ r3 | — | — | 收窄复核该 Medium 已闭合、无新问题 | 已闭合 |

- **⑤ r1**：两条 Medium 的具体内容见上表；整改后 r2 以收窄输入原位重跑，判 CLOSED。
- **⑤ r1 保留的 4 条 Low**：EN 动词闭集全为过去分词，v1.40 的该行用主动态 `lands` 时 EN 侧 0 声明；`翻` 无词边界（`翻倍` 命中）；`TMP_REF_RE` 无前置边界（`attmp/x` 命中）；`OB-101` 与 `OB-102` 两条历史补登台账行未在交付物枚举中出现。
- **③ 保留的 5 条 Low**：见 §7.4 的 D-1 至 D-5。
- **⑥ 与 ⑦**：两者均为盲输入；⑦ r1 对五组逐组给出通过结论、未报 Medium+。
- **⑦ 收窄链**：r2 只核自由度清单，出 1 条 Medium（清单漏登 §5.1 的 `*.raw.txt` 优先级次序）⇒ 清单追加 F12 整改 ⇒ r3 判该 Medium 已闭合、无新问题。
- **⑦ 额度使用 = 3 / 3**：r1 终审 + r2 收窄 + r3 收窄。

## 6 审查输入包摘要（含盲审隔离证明）

| 包 | 位置 | 内容 | 清单 |
|---|---|---|---|
| ⑤ r1 | 历史缺失（随 D6 的越权清理丢失） | brief / 切片定义 / 设计 / 全量 diff / 门禁读数 | 历史缺失 |
| ⑤ r2 | `tmp/GATES-PATCH/review-5r2/` | brief / 两条 Medium 的 hunk / 计数同步 / 自测 hunk / 读数 / 三作用域门禁输出 | `MANIFEST.md` |
| ⑥ r1 | `tmp/GATES-PATCH/review-6/` | brief / 全量 diff / numstat / 状态 / 三作用域门禁 / 自测 / 读数 | `MANIFEST.md` |
| ⑦ r1 | 历史缺失（⑧a 未读取到该包） | 与 ⑥ 同形的输入，brief 为对抗式 | 历史缺失 |

- ⑥ 与 ⑦ 的 brief 均按盲输入派单：未附 ⑤ 结论、未附任何前轮发现。
- 包内不含前轮结论文件；审查者的暂存目录与输入包分离。
- ⑥ 与 ⑦ r1 未附自由度清单（清单在 ⑦ r1 之后产出，见 D10）。
- 输入消毒：包内只有切片定义、设计、测试报告、diff、numstat、门禁读数、自测输出与判据源码。

## 7 可证伪性

### 7.1 边界探针（9 条）

| # | 探针 | 期望 | 实际来源 |
|---|---|---|---|
| E1 | `lang:zh:en` 整 token 消费 | `UsageError`（exit 2） | T84 正例 |
| E2 | `lang:zh:` 整 token 消费 | `UsageError`（exit 2） | T84 正例 |
| E3 | 合法 `lang:zh` 仍判定为 zh | 通过 | T84 反例 |
| E4 | `ADR-027` 不产出 `DR-027` | 无 `DR-027` 声明 | T85 新增断言 |
| E5 | `登记 ADR-028 与 OB-103` | 只产出 `OB-103` | T85 新增断言 |
| E6 | `OB-97x` 只加前侧边界 | 截出 `OB-97`（两侧同构，D11） | ③ 手测 |
| E7 | G5-d 缺键 | FAIL | T86 断言 |
| E8 | G5-d 重复键、空值或缺少 frontmatter | FAIL | T86 断言 |
| E9 | G9-e 指向不存在的路径 | 只出 INFO | G9-e 实现与门禁读数 |

### 7.2 历史重放

- G6-7 历史全量重放：提交数 50、文件路径次 96、发现数 0、错误数 0。
- 重放覆盖全部历史提交，不是抽样；重放脚本随 D6 的清理丢失，读数为主控既有记录。

### 7.3 变异载荷

| 断言块 | 条数 | 载荷 |
|---|---|---|
| T84 | 7 | 2 条正例、3 条反例、2 条变异 |
| T85 | 2（本片新增） | 前置边界两向 |
| T86 | 8 | G5-d 的五键形态：缺键、重复、空值、无 frontmatter，额外键不判 |
| T88 | 3 | 前置门降级后只保留记录态的载荷断言 |

- 变异动作与被断言结论不变；载荷细节见自由度清单的 M1 与 M3。

### 7.4 覆盖边界与已知 Low

| 编号 | 内容 | 归属 |
|---|---|---|
| D-1 | `DECL_ID_RE` 仅认 ASCII 连字符，U+2011 可绕过 | 记录在案 |
| D-2 | `TMP_REF_RE` 大小写敏感且不认反斜杠；OB-77 自身台账行产生 2 条 INFO | 记录在案 |
| D-3 | `--check=G5-d` 连带执行 G5 模块并多出一条 INFO | 记录在案 |
| D-4 | `--check=G6-7` 在缺少 `lang:` 的夹具上 exit 2 | 记录在案 |
| D-5 | 嵌套的 `.github/agents/sub/*.agent.md` 不判 | 记录在案 |

- G9-e 只出 INFO、不参与结论（F5）。
- G6-7 只机械化 ID 子类；无 ID 的措辞级对账由 ⑤⑥⑦ 与 G6-5 承担。

## 8 门禁结果

| 命令 | 结果 |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 588/588` |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=1`（仅 G5(f)：三条 ⑧a 新增路径未跟踪） |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=1`（仅 G7-b：三条 ⑧a 新增路径未暂存） |
| `G6-7` 真实变更集 | `scanned 2, declarations 3` |
| `G9-e` index 与 tree | ⑧a 前 `scanned 92 added line(s), 4 reference(s), 2 unresolvable`；⑧a 后 tree `scanned 640, 16 reference(s), 14 unresolvable`，均为 INFO |
| `G5-d` 本轮 | `no changed .github/agents/*.agent.md in scope` |
| 角色文件普查 | 14 项角色文件中仅 `sol-orchestrator.agent.md` 缺 `user-invocable`，属既有偏离 |
| 基础 `go test ./...` | 由 `--all` 内嵌执行 |

⑧a 的逐条自检读数：

| 判据 | 读数 |
|---|---|
| G7-a | PASS |
| G1-a | PASS |
| G2 | PASS |
| G3-a | PASS |
| G6-4 | PASS |
| G8-a | PASS |

- index 与 tree 两个作用域的 `TOTAL_FAIL=1` 分别只来自 G5(f)（三条新增路径未跟踪）与 G7-b（三条声明为 `repo` 的新增路径未暂存）；⑧a 不执行 `git add`，两者由主控暂存后同时消解。首次 `--all` 的 G4b 报出 1 个未登记字（码位 U+53F3）；未改白名单，改以同义表述消除后复跑为 `unknown=[]`。
- 三作用域复跑口径为独立运行，不复用同一进程内的结果。

**CN 与 `_EN` 六形状（`docs/DELIVERY_DIRECTIVE.md` 对 `_EN`）**：

| 形状 | CN | EN |
|---|---|---|
| 行数 | 1845 | 1845 |
| 标题 | 122 | 122 |
| `**` | 4438 | 4438 |
| `|` | 2947 | 2947 |
| 空行 | 447 | 447 |
| 反引号 | 4528 | 4528 |

- 两侧均为 UTF-8 with BOM 与 LF。

### 8.1 已知缺口（⑨ 追加事实登记）

`gate.js --list` 把 `base` 列为第 19 个 id，但 `gate.js --check=base` 返回 exit 2 与 `unknown criterion id: base` ⇒ `--list` 与 `--check` 的 id 集不一致；`base` 仅在 `--all` 下执行。此为既有缺陷、非本片引入，本片不修，登记备查。 已登记为 **OB-105**（见 §12.4 台账）。

### 8.2 断言预算残余

断言预算低于 ① 预计 6 条，降级理由为 T88 前置门由 9 条降到 3 条；**⑦ 未逐条核对降级合理性**，登记为残余（见 §3 的 D2）。本片不补跑 ⑦（额度已 3/3），亦不以「数字自洽」替代「降级合理」。

## 9 ⑨ 原生验证与提交状态

**⑨ 已执行**：2026-10-20，主控原生验证；本节读数均为实测。

| 项 | 读数 |
|---|---|
| 逐判据 `--check=<ID> --scope=index` | 18 项全部 exit 0：`G1-a` / `G1-b` / `G1` / `G2` / `G3-a` / `G3-b` / `G4a` / `G4b` / `G5-a` / `G5-b` / `G5-c` / `G5-d` / `G5` / `R2.5` / `G6` / `G7` / `G8` / `G9`；其中 `G5-d` 为 INFO（作用域内无 `.github/agents/*.agent.md` 变更），其余 17 项为 PASS |
| 三作用域 `--all` | `index` / `tree` / `ci` 均 `TOTAL_FAIL=0`、exit 0 |
| `selftest.js` | `SELFTEST: PASS 588/588`（不低于本片门槛 586） |
| 干净 worktree 重放 | `git worktree add --detach tmp/GATES-PATCH/pristine HEAD`（HEAD = `2482971`）→ `git diff --cached --binary` 导出暂存补丁 → `git apply --index` → 在 pristine 内 `gate.js --all --scope=index` 得 `TOTAL_FAIL=0`、`selftest.js` 得 `588/588` ⇒ 变更集自洽，不依赖工作区或未跟踪文件；随后 `git worktree remove --force`、补丁文件删除，`git worktree list` 仅剩主工作区 |
| 退出码矩阵 | 合法判据 id → 0；`--check=<不存在的 ID>` → **2**（`unknown criterion id`）；`--scope=<非法值>` → **2**（`invalid --scope value`） |
| `G5(d)` 围栏闭合 | `G5(d) staged is a subset of declared :: 16 path(s) declared` ⇒ artifacts 围栏与暂存集双向闭合 |

**⑩ 已执行（2026-10-20）**：主控出 ⑩ 停点回执，用户总体接受并提出 5 项补正（见 §3 的 D2 / D4 / D6 与 §8.1 / §10），全部在本轮回填提交内落地。

**提交与 CI**：主提交 `cd6e974`（`feat: land the G6-7 / G5-d / G9-e criteria and the v1.40 directive package`，覆盖 16 项路径）已推送 `origin/main`；CI run **`39384174218`** 双绿——`Go ubuntu-latest` 与 `Go windows-latest` 均 `success`，两腿 `Build` / `Test` / `Gates` 步骤全 `success`，两腿 `SELFTEST: PASS 588/588`、`scope=ci changed=16`、`TOTAL_FAIL=0`。回填提交见本报告落盘后的本次 docs-only 提交。

## 10 成本与计量

| 位 | 次数 | 模型 / 说明 |
|---|---|---|
| ① | 1 | `DeepSeek V4 Pro (deepseek)` |
| ② | 5 | 主控模型（步骤 0 落地、实施 1、实施 2、整改 M-1/M-2、断言数同步） |
| ③ | 1 | 主控模型 |
| ⑤ | 2 | `GLM-5.3 (glm)`，配额 1+1 用尽 |
| ⑥ | 1 | `MAI-Code-1.1-Flash (copilot)`，配额 3 |
| ⑦ | 3 | `GPT-5.3-Codex (copilot)`，配额 1+3（r1 终审 + r2 / r3 收窄） |
| 批量辅助 | 0 | 未动用 ADR-027 的授权：本片②的实现与③的自测断言各为单批派单，§2.2 的四类批量任务（枚举 / 统计、长文档摘录与数字核对、双语镜像语言侧检查、MANIFEST 草案）在本片无对应的独立重复工作量；两处真正耗时的取证（G6-7 收窄的 63 例复算、50 提交历史重放）是决策依赖的串行工作，不属该四类 |
| 探针 | 0 | 未使用探针额度 |
| ⑧ | 3 | ⑧a 报告与清单、清单补登 F12、⑧b 定稿 |

- 按档位 1+3 计，⑦ 实耗 **3** 轮（r1 终审 + r2 / r3 收窄）；本轮不给 credits 数值，只给可核的模型与轮次事实。

## 11 明确未执行事项

- 主提交与回填提交均已执行并推送 `origin/main`；未触 gitee。
- ⑨ 原生验证、⑧b 定稿、⑧c 收尾清洁与 ⑩ 停点均已执行（见 §9）。
- ⑦ 三读均已执行（r1 终审、r2 收窄、r3 收窄），额度使用 3 / 3。
- 未改 CONTRACTS、schema、fixtures、workflows 与 `internal/**`。

## 12 下一步建议

1. 主提交 `cd6e974` 已按 artifacts 围栏逐件提交并推送（覆盖 16 项路径），CI run `39384174218` 双绿。
2. ⑨ 已在纯净 worktree 复现变更集，读数见 §9。
3. ⑧c 已在记录落盘后逐件清理 `tmp/`：`Test-Path tmp/GATES-PATCH` = False、`tmp` 递归文件计数 = 1。
4. 下一片由用户指定（池 A `DIRECTIVE-BOUNDARY-CLARIFY`），本片不自行起片。

## 13 执行自由度

逐条引用 `docs/validation/evidence/GATES-PATCH-freedom-list.md`（清单为权威，本表不改变其含义）。

| # | 选择 | 落点 |
|---|---|---|
| F1 | `DECL_ID_RE` 前置边界字符类取 `[A-Za-z0-9_-]`，并补 2 条 T85 断言 | `g6.js`、`selftest.js` |
| F2 | `G6-7` 动词闭集取 CN 9 词与 EN 7 词，全为过去式或完成体 | `g6.js`、§6.3 |
| F3 | `G6-7` 证据范围取整个 §12.4 与附录 C 区，判定行自身不作证据 | `g6.js` |
| F4 | 文档既无 §12.4 又无附录 C 时 `G6-7` 不判 | `g6.js`、§6.3 |
| F5 | `G9-e` 只出 INFO、不设白名单路径、不参与结论 | `g9.js`、§6.3 |
| F6 | `G5-d` 只判 `.github/agents/*.agent.md`，额外键不判 | `g5d.js` |
| F7 | 自测块 T88 由设计的 9 条缩为 3 条 | `selftest.js` |
| F8 | 两个候选形状降级为记录态 | §6.3、§12.5 |
| F9 | 降级以「记录态裁决」段落承载于 §6.3 的 G3 行，不新增节号 | §6.3 |
| F10 | 版本包以行内补注承载，不新增段落 | §6.3 |
| F11 | 「组 1 至组 4」是报告与派单的组织单位，不改变判据注册结构 | 本报告、派单 |

- 已知偏离 D1 至 D11 的逐条登记见 §3 与清单第 4 节；D12（⑦ r2 的 Medium，清单漏登项）登记在清单第 4 节。

DOCS: DONE

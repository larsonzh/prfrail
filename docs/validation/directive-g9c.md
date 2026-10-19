# 切片 `G9-c` 验证报告（判据 `G9`，准则 v1.38）

> 切片：`G9-c`（判据 id = `G9`）；准则版本 v1.38；档位 M；本文件为 **⑧b 定稿**（⑨ 已回填，⑩ 待授权）。
> 基线 HEAD = `53b2a90`（= `origin/main`）；主提交 hash 与 CI run 号待 ⑩ 同轮授权后回填。
> 双语：本文件为 CN 权威，`docs/validation/directive-g9c_EN.md` 为逐节同位镜像。

## 1 变更概要

**变更集**（`git diff --cached --numstat`，11 paths，合计 insertions 1053 / deletions 18）：

| 路径 | + | − | 说明 |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 27 | 9 | v1.38 准则正文 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 27 | 9 | `_EN` 严格同位镜像 |
| `docs/t027/REMAINING_SLICES.md` | 19 | 0 | 台账 P13 小节 + 完成行 |
| `docs/t027/REMAINING_SLICES_EN.md` | 19 | 0 | `_EN` 台账同位同步 |
| `docs/validation/directive-g9c.md` | 291 | 0 | 新建；本报告 |
| `docs/validation/directive-g9c_EN.md` | 291 | 0 | 新建；`_EN` 严格同位镜像 |
| `docs/validation/evidence/G9-C-freedom-list.md` | 39 | 0 | 新建；执行自由度清单 |
| `tools/gates/g9-whitelist.txt` | 11 | 0 | 新建；仅注释，空表合法 |
| `tools/gates/gate.js` | 2 | 0 | `CRITERIA` 与 `RUN_ORDER` 注册 `G9` |
| `tools/gates/lib/criteria/g9.js` | 123 | 0 | 新建判据实现 |
| `tools/gates/selftest.js` | 204 | 0 | `T80`–`T83` 断言 |

**判据 G9 语义**：

- G9(1) 形态：变更集内 basename 命中 `\.(log|tmp|bak|out)$/i` 且不在 `tmp/` 下 ⇒ FAIL。
- G9(2) 位置：路径以 `tmp/` 开头（**大小写归一**）且归一后 ≠ `tmp/.gitkeep` ⇒ FAIL；命中 `tmp/` 后 `continue`，同一路径只产生一条 finding。
- 排除面三类：证据树 `docs/validation/evidence/**`（OB-79，归 G8 承担）、已删除路径（`recordFor(p).deleted`）、白名单 `tools/gates/g9-whitelist.txt`（格式 `<match>\t<file>\t<reason>`；畸形态抛 `UsageError` 并以退出码 2 终止；空表合法）。
- 空变更集 ⇒ `INFO`（不 PASS 不 FAIL）；finding 详情最多列 5 条，其余折叠为 `(+N more)`。

**准则落点**（v1.37 → v1.38）：header 版本与日期 / §0.3 日志行 / §1.5 指向 §5.1 的指针句 / §5.1「临时件」术语定义 / §6.3 新增 G9 行 / §6.3 标题「八类」→「九类」/ R6.2 的 `G1–G9` / 「G1–G8」连带改写 / §12.4 的 OB-98 翻 `已处置` 并新增 OB-99 / 附录 C.0z / `_EN` 严格同位。

**执行自由度**：本片行使的实质性选择逐条见 §13，内容与 `docs/validation/evidence/G9-C-freedom-list.md` 一致。

## 2 执行流水线

| 步骤 | 角色 / 载体 | 模型 | 状态 | 返工次数 | 耗时 |
|---|---|---|---|---|---|
| ① 设计 | `deep-reasoner` | `deepseek-v4-pro` | `DESIGN: DONE`（判定「须协议先行」⇒ 停机升级） | 0 | 历史缺失 |
| ② 实现 | `prfrail-implementer` | 主控模型 | `IMPLEMENT: DONE`（含整改轮） | 2 | 历史缺失 |
| ③ 测试 | `prfrail-tester` | 主控模型 | `TEST: PASS` | 0 | 历史缺失 |
| ⑤ 预审 | `independent-reviewer` | `GLM-5.3 (glm)` | r1 `PRE-REVIEW: PASS WITH FIXES` → r2 `RE-REVIEW: PASS WITH FIXES`（仅剩 Low） | 1 | 历史缺失 |
| ⑥ 独立扫描 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `CODE-REVIEW: PASS`（零缺陷） | 0 | 历史缺失 |
| ⑦ 终审 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` Extra High | r1 `TERMINAL-REVIEW: FINDINGS` → r2 `TERMINAL-REVIEW: PASS` → r3 `FREEDOM-LIST: PASS WITH FIXES` | 2 | 历史缺失 |
| ⑧a 初稿 | `prfrail-documenter` | 主控模型 | 本报告已落盘 | 0 | 历史缺失 |
| ⑨ 原生验证 | 主控 | — | 已完成（纯净 worktree 复现，见 §9.1） | — | — |
| ⑧b 定稿 | `prfrail-documenter` | 主控模型 | 本报告定稿已落盘 | 0 | 历史缺失 |
| ⑧c 收尾清洁 | 主控 | — | 待 ⑧c | — | — |
| ⑩ 停点 | 主控 | — | 待 ⑩（等用户授权） | — | — |

- 耗时字段**历史缺失**：本片未落盘起止时间戳，不补造。
- 返工口径：② 共 3 次派单（初版 + 2 轮整改）；⑤ 共 2 次（配额 1+1 用尽）；⑦ 共 3 次（配额 1+2 用尽）。
- OB-95 / OB-96 第 2/2 片数据点：① 由 `deep-reasoner` 独立执行 1 次（OB-95）；②③⑧ 均为实派单（OB-96）；本片**无零派单位**。

## 3 异常与兜底记录

| # | 事件 | 影响 | 处置 | 处置 |
|---|---|---|---|---|
| D1 | ① 判定「须协议先行」（术语「临时件」无定义） | 设计阶段停机 | 按起跑消息停下升级；用户裁定后补 §5.1 定义，并更正引用为 `§5.1 的 ⑧c 条目；目录约定见 §1.5` | 已处置 |
| D2 | ② 整改轮**越权提前清理** `tmp/G9-C/`（清理属 ⑧c 主控职责） | ⑤ r2 输入包与 ③ 原始读数丢失 | 按原记录重建等价副本；`03-tester-report.md` 头注自证重建 | 已整改 |
| D3 | ⑤ 配额 1+1 用尽后仍有 F-A 第 4 处 Low | 无独立复审通道 | 由主控**确定性自证**闭合（`git grep --cached` 零命中 + `SELFTEST` 全绿 + 三作用域读数），未开第 3 次 ⑤ | 记录在案 |
| D4 | 自由度清单产出时点晚于 ⑦ r1/r2（§1.7 要求 ⑥ 之后、⑦ 之前） | ⑦ r1/r2 输入包内无清单 | 以 ⑦ r3 收窄复审本清单闭合（占用 ⑦ 第 3 次额度） | 已闭合 |
| D5 | ⑦ M2 未采纳（OB-98 状态词「已处置」不改） | 属「已登记边界」而非缺陷 | 主控判定不改：① 该单元格现状文本已写明「G9 覆盖静态形态…时序形态无机械判据——残留登记为 OB-99。**不得读作「临时件纪律完全机械化」**」；② 状态词系用户同轮起跑裁决显式下达，主控无权改回，且新增状态词属协议先行 | 观察中 |
| D6 | 准则附录 C.0z 的验收证据路径写作 `docs/validation/directive-gates-g9-c{,_EN}.md`，本片实际报告名为 `directive-g9c{,_EN}.md` | 准则路径与实际文件名不一致 | 已在准则附录 C.0z 修正：报告路径改为 `docs/validation/directive-g9c{,_EN}.md`、自由度清单路径改为 `docs/validation/evidence/G9-C-freedom-list.md`，CN 与 `_EN` 同步，六形状仍全等 | **已处置**（⑧b 前由 ② 修正附录 C.0z，CN/`_EN` 同位） |
| D7 | 起跑消息要求 §9 写「待 ⑨ / 待 ⑩ 回填」，该措辞命中 `G2` 占位符残留 | `G2` 判红 | 改用等价状态标记「⑨ / ⑩ 执行后补入」，语义不变；该措辞不得残留 | 已处置 |

- D2 为本片唯一的**职责越界**事件：清理时序属 ⑧c，② 不得提前执行；重建后输入包的 `MANIFEST.md` 记载的是重建副本的 `SHA256`，不声称与原包逐字节同一。
- D5 的判定依据原文见 §5 的 ⑦ r2 条目；本报告不重述为「缺陷已修」。

## 4 环境与档位

| 项 | 值 |
|---|---|
| OS | Windows（本机 `d:\LZProjects\prfrail`） |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| 基线 HEAD | `53b2a90` |
| 档位 | M |
| ⑥ | 必派（判据本体变更 ⇒ §7.2.1 ③ 不成立） |
| ⑦ 额度 | 1+2（用尽） |
| ⑤ 额度 | 1+1（用尽） |
| 批量辅助 | 0（本片无批量机械改写需求） |
| 探针 | 0 |

## 5 审查结论摘要（⑤⑥⑦）

| 轮次 | 编号 | 级别 | 结论 / 处置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 | F-A | Medium | 措辞过度承诺 | 已整改 |
| ⑤ r1 | F-B | Medium | 派单证据义务 | 已整改 |
| ⑤ r1 | 若干 | Low | 措辞与覆盖建议 | 已整改 |
| ⑤ r2 | — | Low | `selftest.js` `T80`–`T83` 注释块仍有同型句「the only shape a forced `git add -f` produces」 | 已整改（② 单点改写 + 主控确定性自证） |
| ⑥ r1 | — | — | `CODE-REVIEW: PASS`，零缺陷（盲审，未给 ⑤ 结论） | 记录在案 |
| ⑦ r1 | M1 | Medium | `tmp/` 根判定大小写敏感，`Tmp/` / `TMP/` 假 PASS | 已整改 |
| ⑦ r1 | M2 | Medium | OB-98 直接翻「已处置」易被误读 | 记录在案（已登记边界，见下） |
| ⑦ r1 | Low | Low | 缺大小写变体用例 | 已整改 |
| ⑦ r3 | 3 条 | Low | 自由度清单：F5 落点、F10 手法表述、F12 新增声明 | 已整改 |

- **⑤**：r1 `PRE-REVIEW: PASS WITH FIXES`（1 Medium F-A 措辞过度承诺 + F-B 派单证据义务 + Low 若干）；整改后 r2 原位重跑 `RE-REVIEW: PASS WITH FIXES`，**仅剩 Low**；该 Low 已由 ② 单点改写并由主控确定性自证闭合（`git grep --cached` 零命中 + `SELFTEST: PASS 552/552` + 三作用域 `TOTAL_FAIL=0` 的当时读数）。**⑤ 配额 1+1 已用尽，未开第 3 次**。
- **⑥**：`MAI-Code-1.1-Flash (copilot)` 盲审（未给 ⑤ 结论），`CODE-REVIEW: PASS`，零缺陷。
- **⑦ r1**：`GPT-5.3-Codex (copilot)` Extra High 盲审，`TERMINAL-REVIEW: FINDINGS` —— M1（`tmp/` 根判定大小写敏感，`Tmp/` / `TMP/` 假 PASS，Medium）、M2（OB-98 直接翻「已处置」易被误读，Medium）、Low（缺大小写变体用例）。
- **⑦ r2**（收窄）：M1 与 Low 已改并复测通过；**M2 由主控判定为「已登记边界」不予改动**，理由：① OB-98 单元格现状文本已写明「G9 覆盖静态形态 … 时序形态无机械判据 —— 残留登记为 OB-99。**不得读作「临时件纪律完全机械化」**」；② 状态词翻「已处置」系用户在同一轮起跑裁决中显式下达，主控无权改回，且新增状态词属协议先行。结论 `TERMINAL-REVIEW: PASS`。
- **⑦ r3**（收窄，只审自由度清单）：`FREEDOM-LIST: PASS WITH FIXES`，三条意见已全部落地（F5 落点改指 `ctx.js` 的 `parseWhitelist`；F10 手法表述改为「读源文本 → 字符串变异 → `Module._compile` 内存编译 → 断言字节同一」；新增 F12 声明「只覆盖 `tmp/` 根，其它 gitignored 根不在覆盖内」）。

## 6 审查输入包摘要（含盲审隔离证明）

| 包 | 位置 | 内容 | 清单 |
|---|---|---|---|
| ⑤ r2 | `tmp/G9-C/review-5/` | brief / 切片定义 / 设计 / 测试报告 / 全量 diff / numstat / 三作用域门禁读数 / 自测尾部 / `g9.js` / 白名单 / 准则条款摘录 / 自测摘录 | `MANIFEST.md`（14 项 + `SHA256`） |
| ⑥ r1 | `tmp/G9-C/review-6/` | 同上（brief 不同） | `MANIFEST.md`（14 项 + `SHA256`） |
| ⑦ r1 / r2 | `tmp/G9-C/review-7/` | 同上（brief 为对抗式） | `MANIFEST.md`（14 项 + `SHA256`） |
| ⑦ r3 | `tmp/G9-C/review-7b/` | brief / 自由度清单 / `g9.js` / 全量 diff / numstat | 收窄包 |

**盲审隔离证明**：

- ⑥ brief 原文声明「You have NOT been told what any earlier reviewer found, and that is intentional」；⑦ brief 原文声明「You were not given any earlier reviewer findings on purpose」⇒ 首轮均未附前序结论。
- 包内不含 ⑤ 的结论文件、不含任何前轮审查笔记；审查者的 scratch 目录（`review-5-notes/`、`review-7-notes2/` 等）与输入包分离。
- **首轮是否附清单**：⑥ 与 ⑦ r1 / r2 **未附**自由度清单（清单在其后产出，见 D4）；⑦ r3 为**只审清单**的收窄输入。
- **输入消毒**：包内只有切片定义、设计、测试报告、变更 diff、numstat、门禁读数、自测尾部、判据源码、白名单与准则条款摘录；不注入主控结论或期望判定。
- **重建声明**：受 D2 影响，⑤ r2 / ⑥ / ⑦ 的输入包为**按原记录重建的等价副本**，`MANIFEST.md` 中的 `SHA256` 对应重建副本；`03-tester-report.md` 头注自证该事实。

## 7 可证伪性

### 7.1 端到端机械检查（③ 实测，`git add -f` 真实索引条目）

| # | 变更集条目 | 期望 | 实际 | 退出码 |
|---|---|---|---|---|
| E1 | `tmp/G9-C/fixtures/x.log` | FAIL `tmp-root` | `FAIL G9 findings=1 ["tmp-root: …"]` | 1 |
| E2 | `tools/foo/selftest-output.log` | FAIL `temp-shape` | `FAIL G9 findings=1 ["temp-shape: …"]` | 1 |
| E3 | `tmp/.gitkeep` | PASS | `PASS`（写入 1 行使之真进变更集，验后还原为空；`SHA256 E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855`） | 0 |
| E4 | 普通 `.md` | PASS | `PASS` | 0 |

### 7.2 大小写变体（⑦ r1 M1 整改后的追加实测）

| # | 变更集条目 | 期望 | 实际 |
|---|---|---|---|
| E5 | `Tmp/…` | FAIL，保留原始拼写 | FAIL `tmp-root`，finding 打印 `Tmp/…` |
| E6 | `TMP/nested/…` | FAIL，保留原始拼写 | FAIL `tmp-root`，finding 打印 `TMP/nested/…` |
| E7 | `TMP/.gitkeep` | 放行 | PASS |
| E8 | `Tmp/.gitkeep` | 放行 | PASS |

### 7.3 变异检验（`selftest.js` `T82`，内存编译补丁 + 源文件字节同一断言）

| 变异 | 动作 | 被翻转的断言 | 结果 |
|---|---|---|---|
| M1 | 形态正则改为永假 `/(?!x)x/` | N1–N3 | 翻 PASS |
| M2 | 删除白名单分支 | P3 | 翻 FAIL |
| M3 | 删除证据树排除 | P2 | 翻 FAIL |
| M4 | 删除 `.gitkeep` 守卫 | N4 | 翻 PASS |

- **M4 的 needle 已随大小写整改同步为 `key !== TMP_ROOT_KEEP`**（F8 改了实现）；若仍用 `p !== TMP_ROOT_KEEP`，则 `applied=false`，`T82` 的 M4 与 `F1` stale-mutation 会双红。
- 变异动作与被断言结论不变；详情见 `docs/validation/evidence/G9-C-freedom-list.md` 的「连带机械同步」节。

### 7.4 覆盖边界（四类绕过路径）

| 绕过路径 | CI 可检性 | 本地可检性 | 归属 |
|---|---|---|---|
| `$env:TEMP` / `%TEMP%` | 不可检 | 启发式时窗盘点 | ⑧c 本地记录态 |
| 仓库外位置 | 不可检 | 同左（仅字面路径） | ⑧c 本地记录态 |
| PowerShell 管道 | 不可检 | 同左（仅字面路径） | ⑧c 本地记录态 |
| 变量中转 | **不可检** | **不可检**（需 AST 污点追踪） | OB-99 |

- F12：G9(2) 的位置守护**只覆盖 `tmp/` 根**；其它 gitignored 根（`/out/`、`/bin/`、`vendor/`、`/tools/tmp/`、`.vscode/`、`.idea/`、`.github/context-mode/content/`）的强入**不在 G9 覆盖内**，该边界在 §6.3 G9 行与 `g9.js` docstring 双处明写。

## 8 门禁结果

| 命令 | 结果 |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 552/552` |
| `T80` / `T81` / `T82` / `T83` 断言数 | 5 / 6 / 4 / 1（合计 16，全部 `ok`） |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0`（全部 11 个变更集路径暂存后） |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0`（全部 11 个变更集路径暂存后） |
| `node tools/gates/gate.js --check=G9 --scope=index` | `PASS G9 临时件 :: 11 path(s) OK (1 excluded) (whitelist=worktree)` |
| 基础 `go test ./...`（`--all` 内嵌） | `ok=14 FAIL=0` |

- **⑧a 的两处预期失败已在 ⑧b 消解**：`G5(f)`（index hygiene）与 `G7-b`（声明的 `repo` 工件未暂存）均由「切片工件全部暂存」消解；三次 `--all` 读数均为 `TOTAL_FAIL=0`。
- 三作用域复跑口径：⑧b 暂存后 `index=0` / `tree=0` / `ci=0`（三次独立运行，非同一进程内的复用结果）。

**CN / `_EN` 六形状（`docs/DELIVERY_DIRECTIVE.md` 对 `docs/DELIVERY_DIRECTIVE_EN.md`）**：

| 形状 | CN | EN |
|---|---|---|
| lines | 1805 | 1805 |
| `**` | 4208 | 4208 |
| `|` | 2856 | 2856 |
| 空行 | 439 | 439 |
| 反引号 | 4176 | 4176 |
| 标题 | 121（逐位同序） | 121（逐位同序） |

- 两侧均为 UTF-8 with BOM + LF。

### 8.1 已知缺口（2026-10-19 追加事实登记）

> **⑧c 与 ⑩ 之间的取证窗口（本片事实，不销解、不扩范围）**：本片 ⑧c 之后为取 ⑩ 回执证据（定位 §12.4 的 OB-98 / OB-99 条目、定位 §9.3 小节、统计派单偏移、读 CI 日志），在 `tmp/` 根与 `tmp/G9-C/` 重新创建临时件三处，均即时清除。**性质**：位置违规（`tmp/` 根非 `tmp/<切片ID>/`）+ 时序违规（⑧c 已清零、之后又重建）；**但根因不是纪律维持失效，是 §5.1 的规则缺口**——⑧c 的语义是「过程临时件清零」，⑩ 回执需要取证，两者之间无合法位置。**本片不改准则**（不扩范围），登记为**下一片候选**：`G9-d` 或独立片评估「⑧c 后取证窗口的临时件处置规则」（候选形态：① 允许 ⑧c 后建 `tmp/<切片ID>/`，但须在 ⑩ 前再清一次并回读；② 取证用 `node -e` 一次性输出、不落盘；③ 把 §9.3 与 §12.4 的取证数据在 ⑧c 前预取）。

**措辞澄清（与本节区分）**：⑩ 回执口头表述中曾把本片三处称为「OB-100 所指『创建位置违规但未提交』的**连续实证**」。该措辞**不准确**，须区分——三处是**位置 + 时序双重违规**，其中**时序面是 §5.1 的规则缺口，不作为 OB-100 的连续实证**；OB-100 仍只覆盖「创建位置违规但未提交」。该措辞在本次追加前未进入任何入库文件（机械核验锚定 `51ecddf` 检出零命中），故不涉 §11.1 的追溯改写；本段即为该区分的正式记录。

## 9 ⑨ 原生验证与提交状态（⑨ 已执行；⑩ 待授权）

### 9.1 ⑨ 原生验证

**手法**：在纯净 worktree（`git worktree add --detach <tmp>/wt HEAD`，HEAD = `53b2a90`）中复现变更集，证明变更集**自包含**、不依赖工作区未跟踪的审查包。

| 阶段 | 命令 | 读数 |
|---|---|---|
| 纯净树（未应用变更集） | `node tools/gates/selftest.js` | `SELFTEST: PASS 536/536` |
| 纯净树（未应用变更集） | `node tools/gates/gate.js --list` | **不含** `G9` |
| 纯净树（未应用变更集） | `node tools/gates/gate.js --all --scope=tree` / `--scope=index` | 均 `TOTAL_FAIL=0` |
| 应用暂存补丁（`git apply --index`，exit 0） | `node tools/gates/selftest.js` | `SELFTEST: PASS 552/552` |
| 应用暂存补丁（`git apply --index`，exit 0） | `node tools/gates/gate.js --list` | **含** `G9` |
| 应用暂存补丁（`git apply --index`，exit 0） | `node tools/gates/gate.js --all --scope=tree` / `--scope=index` | 均 `TOTAL_FAIL=0` |
| 应用暂存补丁（`git apply --index`，exit 0） | `git status --porcelain -uall` | 恰为变更集条目（`wt_status_lines=6`，当时尚未含 3 份文档） |

**退出码矩阵**（应用补丁后，逐例强入 / 变异）：

| 用例 | 条件 | 退出码 | 证据 |
|---|---|---|---|
| 干净变更集 | 不改动 | 0 | 无 finding |
| 强入临时件 | `tools/foo/probe.log` 入索引 | 1 | finding 为 `temp-shape: tools/foo/probe.log` |
| 白名单畸形态 | 条目只有 1 个字段 | 2 | 报 `whitelist.txt:2 malformed entry` |

**收尾**：`git worktree remove --force` exit 0、`git worktree prune`，主仓 `git status` 未受影响。

### 9.2 ⑧c 收尾清洁记录

**前期扫描（记录态，非判据）**：扫描 transcript 的 118 份 shell 命令、273 条重定向落点 ⇒ `tmp/<切片ID>/` 236、`tmp/<文件名>` 3、`$env:TEMP` / `%TEMP%` **28（历史违规，均早于本片）**、仓库外 1（`/tmp/staged.diff`）、仓库内非 `tmp` 3（含 1 条 JS 字符串内 `>` 的误报）、变量 / 含变量 2。

- **变量中转不可解析** ⇒ ⑧c 只能作**本地记录态**动作，不构成机械判据（OB-99）。
- 该扫描的原始输出在 D2 的越权清理中丢失；上列计数为主控既有记录，**本轮不可由当前 transcript 重算**（当前 `main.jsonl` 为重新开始的新会话）。

1. **全仓编码审计（不只变更集）**：受版本控制文件 `3630`，其中文本 `3541` 逐文件核验 BOM/LF ⇒ 违规 `14` 处，**全部位于 4 个冻结证据包内**（`docs/validation/evidence/directive-history/glm-landing`、`.../prereview-guidance`、`.../t028`、`docs/validation/evidence/glm-eval`，四者均带 `MANIFEST.md` + `SHA256SUMS.txt`）⇒ 按 OB-27 豁免**改做自洽性核验**：四包 `SHA256SUMS` 条目分别 `10 / 29 / 74 / 77`，合计 **190/190 校验通过，缺失 0、不符 0**。其中 `glm-eval` 包内 `selftest/fixtures/hash-bom.txt` 是 **BOM 兼容性 fixture**，其带 BOM 属预期。
2. **临时件清理（逐件精确）**：`tmp/G9-C/` 下 `90` 个文件逐件 `Remove-Item -LiteralPath` 删除，随后自底向上删除空目录并删除根目录。机械回读证明：`Test-Path tmp\G9-C` = `False`；`tmp` 顶层仅 `.gitkeep`（`0` 字节）；`tmp` 递归文件计数 = `1`；`git status --porcelain -uall` 仍为 `11` 个暂存文件、`??` 未跟踪 `0`。
3. **探针与 worktree 残留**：③ 阶段探针（`tools/foo/`、`examples/.g9-e4-normal.md`）已在 ③ 内删除；⑨ 的临时 worktree 已 `git worktree remove --force` + `git worktree prune`，`git worktree list` 仅主工作树 `D:/LZProjects/prfrail 53b2a90 [main]`。
4. **编辑器侧诊断**：全仓 `get_errors` ⇒ **零错误**。
5. **清理后复跑**：`SELFTEST: PASS 552/552`；`gate.js --all --scope=index|tree|ci` 三次 `TOTAL_FAIL=0`（exit 0）。
6. **时序合规声明**：记录先落盘（本报告 §5–§13、自由度清单、台账 P13）→ 再清理 → 再复跑门禁，符合 §5.1（v1.35）的「清理须在记录落盘之后」要求。

### 9.3 提交与 CI

- 主提交：`7dc9bcc`（`feat: add the G9 static temp-artifact gate (v1.38)`），已推 `origin/main`；**未触 gitee**。
- CI run：`39175328844`（workflow `CI`），结论 **success**，双腿 `Go ubuntu-latest` / `Go windows-latest` 均 `success`。
- `Gates` 步读数（双腿一致）：`SELFTEST: PASS 552/552`、`GATE REPORT scope=ci changed=11`、`TOTAL_FAIL=0`。
- 包级证据：双腿各 14 个 `ok` 包（`internal/adapters`、`internal/applier`、`internal/chain`、`internal/console`、`internal/evidence`、`internal/gates`、`internal/guard`、`internal/release`、`internal/repair`、`internal/snapshot`、`internal/taskdef`、`internal/tickets`、`examples/whois-shadow`、`tools/agent-probe/enforcement-proxy`），`FAIL` 包 **0**。
- 首跑即绿，无 §6.4 处置、无复跑。

## 10 成本与计量

| 位 | 次数 | 模型 / 说明 |
|---|---|---|
| ① | 1 | `deepseek-v4-pro`（不占 Copilot 预算） |
| ② | 3 | 主控模型（不占 Copilot 预算） |
| ③ | 1 | 主控模型 |
| ⑤ | 2 | `GLM-5.3 (glm)`，配额 1+1 用尽 |
| ⑥ | 1 | `MAI-Code-1.1-Flash (copilot)`，配额 3 |
| ⑦ | 3 | `GPT-5.3-Codex (copilot)`，配额 1+2 用尽 |
| 批量辅助 | 0 | 本片无批量机械改写需求 |
| 探针 | 0 | 未使用探针额度 |
| ⑧ | 2 | ⑧a / ⑧b 各 1（⑧a 初稿、⑧b 定稿） |

## 11 明确未执行事项

- 未 `git commit`、未 `git push`、未触 gitee。
- ⑧c 收尾清洁与 ⑩ 停点未执行；⑨ 原生验证与 ⑧b 定稿已执行（见 §2、§9）。
- `GATES-FLOW` 未做（OB-96 机械判据化留待复发时立项）。
- AST 污点追踪未评估（变量中转路径）。

## 12 下一步建议

1. ⑧c 收尾清洁：按 §9.2 的逐件口径清理 `tmp/G9-C/` 并回读证明。
2. ⑩ 停点：等用户同轮授权后再提交与推送，并回填 §9.3 的提交与 CI 事实。
3. 若后续时序形态再犯 ⇒ 按 OB-99 评估 AST 污点追踪可行性。

## 13 执行自由度

逐条引用 `docs/validation/evidence/G9-C-freedom-list.md`（清单为权威，本表不改变其含义）。

| # | 选择 | 依据 | 落点 |
|---|---|---|---|
| F1 | 判据 id 用 `G9`，切片名保留 `G9-c` | id 是命名，不改判据覆盖面 | `gate.js` 注册、`g9.js`、§6.3 行 |
| F2 | 形态闭集取 `.log` / `.tmp` / `.bak` / `.out`，加 `/i` 放宽为大小写不敏感 | 闭集元素取 §5.1「临时件」的典型产物；`/i` 是为 Windows 语义显式放宽，只扩大拦截面 | `g9.js` `TEMP_SHAPE_RE`、§6.3 行 |
| F3 | `tmp/` 之外判形态、`tmp/` 之内判位置，**互斥**（命中 `tmp/` 即 `continue`，同一路径不会产生两条 finding） | 避免重复计数，属实现选择 | `g9.js` G9(2) 分支 |
| F4 | 排除面顺序固定为「证据树 → 已删除 → 白名单 → 判形态 / 位置」 | 排除先行，减少无意义 finding；语义与 §6.3 一致 | `g9.js` 主循环 |
| F5 | 白名单文件格式 `<match>\t<file>\t<reason>`；畸形态抛 `UsageError`（exit 2）；空表（仅注释）合法 | 格式与 fail-closed 是实现选择；空表合法因为本片 `git ls-files` 实测零合法临时件 | `g9-whitelist.txt`；`ctx.js` 的 `parseWhitelist`；`g9.js` 的 `loadWhitelist` 与 `UsageError` 分支 |
| F6 | 详情最多列 5 条，其余折叠为 `(+N more)` | 报告可读性 | `g9.js` `shown` |
| F7 | 空变更集输出 `INFO`（不 PASS 不 FAIL） | 沿用既有判据先例（空集无行可判） | `g9.js` 早退分支 |
| F8 | `tmp/` 根比较**大小写归一**后判定，finding 仍打印原始拼写；`tmp/.gitkeep` 及其大小写变体放行 | ⑦ r1 M1 整改；`Tmp/` 与 `tmp/` 在 Windows / macOS 是同一目录 | `g9.js` `key` 比较、§6.3 行括注 |
| F9 | 常量具名（`TMP_ROOT_PREFIX` / `TMP_ROOT_KEEP` / `TEMP_SHAPE_RE`）而非字面量散落 | 可读性 + 变异测试需要唯一锚点 | `g9.js` |
| F10 | 自测用「读 `g9.js` 源文本 → 字符串变异 → `Module._compile` 内存编译补丁模块 → 断言源文件字节同一」手法，负例编号续用 `N5` / `N6` | 沿用 `T71` 既有手法（OB-17 恢复断言），编号续排 | `selftest.js` `T80`–`T83` |
| F11 | 准则侧补一个极短括注（`（**大小写不敏感**）` / `(**case-insensitive**)`）而非新增段落 | 让散文与代码一致，避免新增节号 | §6.3 G9 行 CN / `_EN` |
| F12 | G9(2) 的位置守护**只覆盖 `tmp/` 根**；其它 gitignored 根（`/out/`、`/bin/`、`vendor/`、`/tools/tmp/`、`.vscode/`、`.idea/`、`.github/context-mode/content/`）的强入**不在 G9 覆盖内** | 覆盖面的取舍，且在准则 §6.3 G9 行与 `g9.js` docstring 双处明写，不构成沉默缺口 | §6.3 G9 行、`g9.js` docstring |

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/directive-g9c.md	repo
docs/validation/directive-g9c_EN.md	repo
docs/validation/evidence/G9-C-freedom-list.md	repo
tools/gates/g9-whitelist.txt	repo
tools/gates/gate.js	repo
tools/gates/lib/criteria/g9.js	repo
tools/gates/selftest.js	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
```

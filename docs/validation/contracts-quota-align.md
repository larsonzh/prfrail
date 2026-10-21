# 切片 `CONTRACTS-QUOTA-ALIGN` 验证报告（准则 v1.41）

> 切片：`CONTRACTS-QUOTA-ALIGN`；准则版本 v1.41（本片不改 `DELIVERY_DIRECTIVE{,_EN}.md`）；档位 **S**（起跑为 M，2026-10-21 经用户裁定降为 S，⑦ 额度随之由 1+2 更正为 1+1）；本文件为 **⑧a 初稿**（⑨ / ⑩ 的读数与提交状态由 ⑧b 回填）。
> 基线：`d:\LZProjects\prfrail`；HEAD = `e91a8b0`（上一片 `DIRECTIVE-BOUNDARY-CLARIFY` 闭环，准则 v1.41）。⑧a 启动时暂存集为空，工作区为 6 项已修改 + 1 项未跟踪（`git status --porcelain` 实测）。
> 双语：本文件为 CN 权威；`docs/validation/contracts-quota-align_EN.md` 为逐行同位镜像，两侧行数相等（机械计数见 §7.2）。

## 1 变更概要

⑧a 之前 6 项路径（5 项已修改 + 1 项未跟踪），⑧a 之后 9 项（新增报告对与自由度清单）。下表增删数为 `git diff --numstat` 机械读数；四个新文件为整文件新增，其增删数取自身行数。

| 路径 | + | − | 说明 |
|---|---|---|---|
| `docs/CONTRACTS.md` | 2 | 2 | §7 ai availability 段 L192 与 A6 段 L232 |
| `docs/CONTRACTS_EN.md` | 2 | 2 | `_EN` 同位（L186 / L226） |
| `docs/t027/REMAINING_SLICES.md` | 22 | 0 | P17 重写为 S 档六字段 + 范围变更说明 |
| `docs/t027/REMAINING_SLICES_EN.md` | 22 | 0 | `_EN` 同位 |
| `tools/agent-stub/main.go` | 24 | 0 | 恒定写实现 |
| `tools/agent-stub/main_test.go` | 364 | 0 | 未跟踪新文件（`TestMain` + 4 个测试） |
| `docs/validation/contracts-quota-align.md` | 333 | 0 | 本文件 |
| `docs/validation/contracts-quota-align_EN.md` | 333 | 0 | 逐行镜像 |
| `docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md` | 45 | 0 | 执行自由度清单 |

**artifacts 围栏**（9 项路径，全部 `repo`）：

```artifacts
docs/CONTRACTS.md
docs/CONTRACTS_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
tools/agent-stub/main.go
tools/agent-stub/main_test.go
docs/validation/contracts-quota-align.md
docs/validation/contracts-quota-align_EN.md
docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md
```

两处契约落点（行号与字符数为机械读数；CN 与 `_EN` 行号不同，各自同位）：

| 段落 | 文件:行 | 字符数（前 → 后） | 变更要点 |
|---|---|---|---|
| ai availability 额度口径 | `docs/CONTRACTS.md:192` | 870 → 1008 | 产品请求预算独立强制执行并记录请求用量；CLI 侧额度为候选平台事实、可高于该预算，实际传参以探测证据的 argv 摘要为准；本记录不含费用维度但不得读作「费用未发生」 |
| 同上 `_EN` | `docs/CONTRACTS_EN.md:186` | 1675 → 2219 | `explicit cost/request budget` → `explicit request budget (…)`，并加入 no-cost-dimension 约束句 |
| A6 pid 报告口径 | `docs/CONTRACTS.md:232` | 3021 → 3455 | 工件清单增列 `stub.pid` / `stub-child.pid` 并注明为 pinned stub 约定；新增恒定写入口径句；新增三态语义与互推禁令；树验证的缺失口径改为 fail-closed |
| 同上 `_EN` | `docs/CONTRACTS_EN.md:226` | 7370 → 8581 | 上述四项的逐句镜像 |

四份文件行数与编码：`CONTRACTS.md` 为 CN 327、`CONTRACTS_EN.md` 为 EN 319（两侧均未变），BOM 保留、CR=0（机械读数）。

**代码落点**（`git diff --numstat` 为 `24 0 tools/agent-stub/main.go`）：

| 文件:行 | 变更 | 依据 |
|---|---|---|
| `tools/agent-stub/main.go:136-141` | `execute()` 内在 `writeOwnPID()` 之后加 `if !run.spawnsChild() { run.writeEmptyChildPIDReport() }`，并以四行注释说明 spawn 模式故意不在此写 | ⑤ r2 Low 的整改形态 |
| `tools/agent-stub/main.go:179-183` | `Start()` 失败时写空报告（确实什么都没起）并返回 5 | 空报告在此时是真实陈述 |
| `tools/agent-stub/main.go:187-192` | pid 列表写失败时杀掉已起的子进程（`_ = command.Process.Kill()`）并返回 6，留下缺失报告（缺口）而非虚假空报告 | ⑤ r1 Low / ⑥ / ⑦ 的同一点 |
| `tools/agent-stub/main.go:260-268` | 新增 `writeEmptyChildPIDReport()`，写失败 exit 11（与 `writeOwnPID` 同码） | 与既有失败码一致 |

未改动：`internal/adapters/ai_probe_copilot.go`（路径 D ⇒ 零改动）与 `internal/chain/postflight.go`（OB-58 已移出本片），二者在 `git diff --numstat` 中均零出现（机械读数）。

## 2 执行流水线

| 环节 | 角色 / 载体 | 模型 | 结论 | 轮次 | 备注 |
|---|---|---|---|---|---|
| 批量辅助 ① | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 完成 | 1 | §7 三处段落定位 |
| 批量辅助 ② | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 完成 | 1 | 五枚事实 / pid 报告 / `maximumRequests` 对照 |
| ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 完成 | 1 | OB-52 四路径比较后推荐 D；OB-58 推荐丙；OB-59 形态 (i) |
| ②a 契约 | `prfrail-implementer` | `DeepSeek V4.1 Flash`（主控同源，白名单内） | 完成 | 1 | 6 处替换，preHits=1 / postOldHits=0 / postNewHits=1 |
| ②b 代码 | `prfrail-implementer` | 同上 | 完成 | 1 | `tools/agent-stub/main.go` +24/-0 |
| ③ 测试 | `prfrail-tester` | 同上 | 完成 | 1 | 三态测试；早于 ④ 与 ⑤ |
| ④ 门禁 | 主控 | — | 完成 | 3 作用域 + selftest | 无派单 |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 1 High + 2 Medium + 4 Low |
| 整改 r1 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 完成 | 1 | 杀掉活子进程 / `stub.pid` 交叉核对 / 重出 patch / 重写 P17 |
| ⑤ r2 | `independent-reviewer` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 新 Medium（patch 陈旧）+ 新 Low（spawn 前预写空报告） |
| 整改 r2 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 完成 | 1 | 取消 spawn 前预写；写失败保持 kill + 缺失报告；patch 补 `--no-index` 段 |
| ⑥ 独立扫描 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 完成 | 1 | 2 High（同一读者背离）+ 1 Medium（缺「缺失不是事实」的钉住测试） |
| ③ 追加测试 | `prfrail-tester` | `DeepSeek V4.1 Flash` | 完成 | 1 | 第四测（spawn + 写不可用 ⇒ exit 6） |
| ⑦ r1（盲审） | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | FINDINGS | 1 | Critical（读者背离）/ High（「恒定写」措辞）/ Medium（中断窗口无测试） |
| 整改 r3 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 完成 | 1 | 按 ⑦ r1 替换文本改写「写入口径」 |
| ⑦ r2 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | High 闭合，CN / EN 仍等价，无新缺陷 |

- 返工口径：② 共 1 次主派单（②a）+ 1 次代码派单（②b）+ 3 次整改派单，未重派、未回滚。
- 耗时字段历史缺失：本片未落盘起止时间戳，不补造。
- ④ 的三作用域与 selftest 为独立运行，不复用同一进程内的结果。

## 3 异常与兜底记录

| # | 事件 | 影响 | 处置 |
|---|---|---|---|
| D1 | ② 在 `tmp/` 根创建 `tmp/seam.js` 并立即删除 | 临时件落点违规（OB-100 形状）；本片内第 **1** 次，升级触发为本片第 3 次；同形跨片已登记 5 次（`$env:TEMP`：`GATES-G1G5` 4 次 + `GATES-EDGE-COVERAGE-RECHECK` 1 次，见 `docs/validation/directive-gates-edge-coverage-recheck.md:62`）⇒ **口径 = 按 §5.1「临时件一律 `tmp/<切片ID>/`」合并计数**（`$env:TEMP` 与 `tmp/` 根同属对该放置规则的违反；OB-100 只是其中「仓内未提交落点不被 G9 检出」的机械盲区子集，**不另立形态**）⇒ 本片为该规则的**第 6 次**，**不得读作全局首次、计数不得从 1 重新起算**，升级触发按合并计数推进 | 已删除；如实登记 |
| D2 | ③ 自报首次写入 `main_test.go` 时重复 `package main` 行且用 CRLF | 文件形态错误 | 已用 `gofmt -w` 修正；闭合，仍登记 |
| D3 | ① 的设计断言「恒定写之后宽松分支不会再触发」为**假**：`child` / `hang` / `fail-fast` / `crash` 在写点前返回 | 断言与代码不符 | 本报告更正该断言（⑦ r2 要求）；宽松分支仍可达 ⇒ OB-110 |
| D4 | ⑦ r1 的 High 属自相矛盾：绝对措辞「恒定写」与本片「spawn 模式不预写空报告」的决策冲突 | 措辞与决策互相冲突 | 以「写入口径 + 当执行到写入点」改写闭合，保留「（恒定写）」括注 |

**非偏离但一并登记**（避免读作遗漏）：

- `stubVersion` 维持 `agent-stub-0.1.0`（`tools/agent-stub/main.go:28`）系有意决定：所有 pin 在本片边界外，precheck 按字符串字面比较，重放身份由 `executableHash` / sha256 绑定而非标签 ⇒ 不得读作遗漏。
- `_EN` 侧 CJK 词 `恒定写` 出现 0 次，对应表达为 `constant write`（新增）与 `constant writing`（既有邻句）；属语言纯度观察，非缺陷。
- `--scope=index` 的 `G5(f)` 在暂存前为红（唯一失败项），属已知暂存前状态、非缺陷，⑨ 复跑。

## 4 环境与档位

| 项 | 值 |
|---|---|
| OS | Windows（本机 `d:\LZProjects\prfrail`） |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| 档位 | S（2026-10-21 由 M 降档，用户裁定） |
| ⑦ 额度 | 1+1（S 档；原 M 档 1+2 已更正，不降级模型） |
| ⑤ 额度 | 1+1 |
| 批量辅助 | ADR-027（限四类任务）；两次均用 `GLM-5.3-FlashX` |
| 探针 | 0 |

## 5 审查结论摘要（⑤⑥⑦）

| 轮次 | 编号 | 级别 | 结论 / 处置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 | H1 | High | `internal/adapters/agent_runner_run.go:300-323` 仍把缺失报告读作「未 spawn」；该分支可达，因为 `child` / `hang` / `fail-fast` / `crash` 从不写报告 | 记录在案（OB-110）；⑤ r2 裁定 `deferred-by-ruling` 可接受 |
| ⑤ r1 | M1 | Medium | P17 仍承诺 M 档范围 | 已整改（P17 重写为 S 档六字段 + 范围变更说明） |
| ⑤ r1 | M2 | Medium | `internal/**` 三处旧注释仍在描述已被证伪的 stub 行为 | 记录在案（越界，见 §14 U3） |
| ⑤ r1 | L1 | Low | pid 写失败时留下空报告且子进程仍活 | 已整改（`tools/agent-stub/main.go:187-192` 杀掉子进程） |
| ⑤ r1 | L2 | Low | `TestSpawnWritesChildPIDList` 对任意 pid 都会通过 | 已整改（与 `stub.pid` 交叉核对，2 秒有界轮询） |
| ⑤ r1 | L3 | Low | `stubVersion` 未 bump | 记录在案（有意决定，见 §3） |
| ⑤ r1 | L4 | Low | 证据 patch 出现乱码 | 已整改（`git diff --output=` 重出） |
| ⑤ r2 | M3 | Medium | patch 陈旧且不完整：生成于 P17 重写之前，且对未跟踪测试文件机械盲 | 已整改（`tmp/CONTRACTS-QUOTA-ALIGN/02-diff.patch:163` 起为 `--no-index` 段） |
| ⑤ r2 | L5 | Low | 空报告在 spawn 之前写 ⇒ `Start()` 与 pid 列表写之间被停会留下「活子进程 + 空报告」的 fail-open 事实 | 已整改（spawn 模式不再预写，见 `tools/agent-stub/main.go:136-141`） |
| ⑤ r2 | — | — | High 处置 `deferred-by-ruling` 可接受；契约文本**不得**为适配宽松读者而收窄 | 已接受 |
| ⑥ | H1–H2 | High | 同一读者背离，被表述为「实现内的陈旧假规则」 | 记录在案（OB-110） |
| ⑥ | M4 | Medium | 无测试钉住「缺失文件不是有效的无子进程事实」 | 已整改（第四测；见 `tools/agent-stub/main_test.go:299`） |
| ⑦ r1 | C1 | Critical | 同一读者背离（`agent_runner_run.go:300-318`，树验证在 `:284` 通过）；配方为返回 `present` 标志，`present == false` ⇒ `treeGone == false` ⇒ uncertain | 记录在案（OB-110，越界） |
| ⑦ r1 | H2 | High | 「始终写」措辞过度断言（写前被中断则报告缺失） | 已整改（改写为写入口径，替换文本逐字采纳，保留「（恒定写）」） |
| ⑦ r1 | M5 | Medium | 无测试同时覆盖中断窗口与消费方决策 | 已接受（消费方不在本片边界内，见 §15 OB-110(c)） |
| ⑦ r2 | — | — | PASS：High 闭合、措辞未削弱裁定义务、CN / EN 仍等价、Critical 与 Medium 处置在用户边界下 ACCEPTED、无新缺陷 | 已闭合 |

- ⑤ r1 的 High 未在本片修复：其修复面在 `internal/adapters/agent_runner_run.go`，属本片边界外 ⇒ 已登记 OB-110。
- ⑦ r2 明确契约文本不得收窄以适配宽松读者 ⇒ 本片保留 fail-closed 口径（`docs/CONTRACTS.md:232` 三态语义句与互推禁令）。
- 修复配方在仓内已有同型实现：`internal/adapters/agent_runner_postflight.go:530` 的 `agentRunnerReadPIDs` 与 `internal/console/runtime_agent.go:680` 的 `readAgentRunnerEvidencePIDs` 都返回 `present` 标志（机械读数，本片未改）。
- ⑥ 的 M4 与 ⑦ r1 的 M5 不是同一缺口：前者要「缺失 ≠ 事实」的钉住测试（已在边界内补），后者要消费方决策的组合测试（消费方在边界外）。

## 6 审查输入包摘要（含盲审隔离证明）

本片**未落盘**审查输入包：`tmp/CONTRACTS-QUOTA-ALIGN/` 下只有设计、diff 与脚本，没有 `review-*` 目录，也没有 `MANIFEST.md`。

| 项 | 读数 |
|---|---|
| `tmp/CONTRACTS-QUOTA-ALIGN/` 文件数 | 10（`01-design.md`、`02-diff.patch`、`apply-6.js`、`check-anchors.js`、`cjk-check-head.js`、`cjk-check-wording.js`、`cjk-check.js`、`model-offsets.js`、`transcript-offsets.js`、`verify-02a.js`） |
| `review-5/`、`review-5r2/` | 均不存在 |
| `review-6/`、`review-7/` | 均不存在 |
| `tmp/` 根 | 仅 `.gitkeep` 与 `CONTRACTS-QUOTA-ALIGN/` |

- ⇒ **盲审隔离无法用 MANIFEST 证明**；派单事实（盲输入、未附前轮结论）只存在于派单文本，⑧b / ⑨ 无法机械复核。该缺口如实登记为 §14 U7，不以「⑤ / ⑥ / ⑦ 均通过」掩盖。
- 可机械核对的两项：`tmp/CONTRACTS-QUOTA-ALIGN/02-diff.patch` 为 78096 字节、含 `b/tools/agent-stub/main_test.go`（`--no-index` 段），且 `tmp/seam.js` 已不存在。
- 输入消毒的其余部分（brief / 设计 / diff / numstat / 门禁读数 / 自测输出）无落盘副本，本报告不复述其内容。

## 7 可证伪性

### 7.1 可证伪断言

| # | 断言 | 反例（删去后应有哪条读数变化） | 机械检查 |
|---|---|---|---|
| E1 | ai availability 段含「不得读作费用未发生」约束 | 删去该约束句后，费用口径再次可被读作「未计费」 | 检索 `费用未发生` 应得 1 处（实测 1） |
| E2 | A6 段含恒定写入口径与三态语义 | 删去写入口径句后，消费方无从区分空文件与缺失 | 检索 `恒定写` 应得 2 处（实测 2） |
| E3 | A6 段含互推禁令 | 删去后「缺失」与「空」可互推，三态退化为二态 | 检索 `互推` 应得 1 处（实测 1） |
| E4 | 恒定写有实现，而非仅文本 | 删去 `tools/agent-stub/main.go:141` 调用后 T1 / T3 转红（变异实测） | `writeEmptyChildPIDReport` 调用点应得 2 处（`:141` / `:182`） |
| E5 | 路径 D 的零代码边界成立 | 若 `ai_probe_copilot.go` 或 `postflight.go` 出现在 numstat，即越界 | `git diff --numstat` 中二者零出现（实测 0） |

### 7.2 机械读数复算

| 读数 | 值 | 来源 |
|---|---|---|
| `git rev-parse --short HEAD` | `e91a8b0` | 命令输出 |
| `git diff --numstat` | 5 项（另 1 项未跟踪） | 命令输出 |
| 契约行数 | CN 327 / EN 319（均未变） | 机械计数 |
| 契约四处字符数 | 870→1008 / 1675→2219 / 3021→3455 / 7370→8581 | `git show HEAD:` 与磁盘实测 |
| `main_test.go` | 行数 364 / 14104 字节 | 机械计数 |
| CN / `_EN` 报告行数 | 两侧相等（见 §8） | 机械计数 |
| `SELFTEST` | PASS 588/588 | `node tools/gates/gate.js --selftest` |
| `G4b` | `newCjk=485 unknown=[]`（⑧a 写盘前） | `--check=G4b --scope=tree` |
| `go test ./tools/agent-stub/` | 4 PASS（`-count=1` / `-count=3` / `-count=5` 全绿） | 命令输出 |
| transcript 派单偏移 | 12 / 12 与脚本复算逐条等值 | `node tmp/CONTRACTS-QUOTA-ALIGN/transcript-offsets.js` |

### 7.3 变异载荷

| 变异（在仓库外副本上执行，不污染工作区） | 预期 | 实测 |
|---|---|---|
| 删去 `tools/agent-stub/main.go:141` 的调用 | T1 / T3 红 | 红 |
| 写父 pid 而非子 pid（`main.go:187-191`） | T2 红 | 红 |
| 吞掉 pid 报告写错误 | T3 红 | 红 |
| 恢复 spawn 前的空报告预写 | 第四测红（exit 11 而非 6） | 红 |

- 变异只做证伪，不进入交付：工作区文件保持整改后形态（`git status --porcelain` 见 §8）。
- 本片未新增门禁断言：`SELFTEST` 维持 588/588。

### 7.4 覆盖边界与已知 Low

| 编号 | 内容 | 归属 |
|---|---|---|
| L-1 | 读者背离（缺失 ⇒ `nil, nil` ⇒ 树验证通过）不在本片 | 记录在案（§15 OB-110） |
| L-2 | `internal/**` 旧注释不在本片边界 | 记录在案（§14 U3 / U4） |
| L-3 | `stubVersion` 未 bump 为有意决定 | 记录在案（§3） |
| L-4 | 本报告不覆盖消费方决策的组合测试 | 记录在案（§5 ⑦ r1 M5） |

- 本报告的结论不覆盖 `internal/console/**` 与 `tools/b4-e2e/**` 的语义等价性；`internal/console/**` 的消费腿已用 `present` 标志 fail-closed（`internal/console/runtime_agent.go:624-667`），与 `internal/adapters/agent_runner_run.go` 的树验证腿不同型。
- CN 与 `_EN` 的语义等价性由 ⑤ / ⑦ 逐句对读与 `G3(a.2)` 形状比较共同承担。

## 8 门禁结果

| 命令 | 结果 |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 588/588`（⑨ 终局读数） |
| 暂存状态（⑨ 终局） | 9 项路径精确暂存、无未暂存改动、无未跟踪项；`git diff --cached --numstat` 逐行 `2/2`、`2/2`、`24/0`、`24/0`、`345/0`、`345/0`、`45/0`、`24/0`、`364/0` ⇒ 新增 **1175** / 删除 **4**，9 个文件 |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0`（exit 0，⑨ 终局读数）；⑧a 落盘后的 `TOTAL_FAIL=1`（唯一失败 `G7-b`，4 项新工件未暂存）已由 ⑨ 暂存消解 |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0`（exit 0，⑨ 终局读数；`changed=9`） |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0`（exit 0，⑨ 终局读数）；⑧a 时点的 `TOTAL_FAIL=1`（`changed=0`，唯一失败 `FAIL G5(f) index hygiene`）属**暂存前状态**，工作树全部暂存后转绿 ⇒ 该缺口关闭（§14 U6） |
| `--check=G3-a --scope=tree` | 见下方逐项读数 |
| `--check=G4b --scope=tree` | 见下方逐项读数 |
| `--check=G5-d --scope=tree` | 见下方逐项读数 |
| `--check=G5-e --scope=tree` | 见下方逐项读数 |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test -count=1 ./...` | 15 个包 `ok`、无 `FAIL`（含 `tools/agent-stub` 9.762s） |
| `go test -count=1 ./tools/agent-stub/` | 4 PASS |
| `go test -count=1 ./internal/adapters/ ./internal/console/` | 两侧均 `ok`（≈35s / ≈8s） |
| 编码审计（9 项暂存路径） | `violations=0`；`.md` 全部 BOM + CR=0，两个 `.go` 无 BOM + CR=0 |

逐判据读数（`--scope=tree`，⑧a 写盘后实测的历史原文；⑨ 终局复跑见 §9）：

PASS  G3(a) symmetry + mirror parity :: docs/CONTRACTS.md +2/-2 vs _EN +2/-2 OK || docs/t027/REMAINING_SLICES.md +22/-0 vs _EN +22/-0 OK || docs/validation/contracts-quota-align.md +333/-0 vs _EN +333/-0 OK ;; docs/CONTRACTS.md {"lines":false,"head":false,"bold":false,"bars":true,"blank":true,"backticks":false} ;; docs/t027/REMAINING_SLICES.md {"lines":false,"head":false,"bold":false,"bars":false,"blank":false,"backticks":false} ;; docs/validation/contracts-quota-align.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true}
PASS  G4b CJK diff minus whitelist :: newCjk=639 unknown=[] (whitelist=37 whitelist-source=worktree)
PASS  G5(d) staged is a subset of declared :: 9 path(s) declared (scope=tree) ; addMode=explicit (derived label, not a command-history assertion)
PASS  G5(e) evidence references resolve :: no added completion row (nothing to resolve)

**CN 与 `_EN` 报告形状**：两侧行数相等、标题骨架逐条相同（机械计数见 §1 与本节上表）；`G9-e` 对 20 处 `tmp/` 引用出 INFO（临时件不入库），属信息级、不阻断。

## 9 ⑨ 原生验证与提交状态

⑨ 已执行并由 ⑧b 回填；本节为终局读数，⑧a 的「未执行」占位已全部替换。

| 项 | 终局读数 |
|---|---|
| 暂存集（9 项精确路径） | `docs/CONTRACTS.md`、`docs/CONTRACTS_EN.md`、`docs/t027/REMAINING_SLICES.md`、`docs/t027/REMAINING_SLICES_EN.md`、`docs/validation/contracts-quota-align.md`、`docs/validation/contracts-quota-align_EN.md`、`docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md`、`tools/agent-stub/main.go`、`tools/agent-stub/main_test.go` |
| 暂存增删（`git diff --cached --numstat`） | `2/2`、`2/2`、`24/0`、`24/0`、`345/0`、`345/0`、`45/0`、`24/0`、`364/0` ⇒ 新增 **1175** / 删除 **4**，9 个文件 |
| 三作用域 | `--scope=tree`、`--scope=index`、`--scope=ci` 均 `TOTAL_FAIL=0`（exit 0）；`index` 在工作树全部暂存后转绿 |
| `SELFTEST` | `PASS 588/588` |
| `go test -count=1 ./...` | 15 个包 `ok`、无 `FAIL` |
| 编码审计（9 项暂存路径） | `violations=0`：`.md` 全部 BOM + CR=0，两个 `.go` 无 BOM + CR=0 |
| HEAD | `e91a8b0`（未变） |
| 提交与推送 | 未提交、未推送；未触 gitee；CI 未触发 |
| ⑩ 停点回执 | 未执行（本片 ⑩ 停点等授权） |
| artifacts 围栏 | 仍精确列 9 项路径（§1），与本表暂存集逐项一致 |

- ⑨ 时点上表 9 项即 `git status --porcelain` 的全部条目（`M ` / `A `），无未暂存、无未跟踪项。
- ⑧b 的回写只触及其中四个文档路径（报告对与台账对），故落盘后这四个路径出现新的未暂存改动：`--scope=index` 复读为 `TOTAL_FAIL=1`，`G5(f)` 列出恰为这四个 `MM` / `AM` 路径；需在 ⑩ 前由主控重新暂存（⑧b 不做 `git add`）。

## 10 成本与计量（R7.1）

**机械审计依据**：下表为本片 `runSubagent` 派单的 transcript 字节偏移（本会话 JSONL），⑧a 用 `tmp/CONTRACTS-QUOTA-ALIGN/transcript-offsets.js` 复算，12 条派单逐条等值。

| 环节 | 角色 | 模型 | transcript 偏移 | 次数 |
|---|---|---|---|---|
| 批量辅助 ①（§7 段落定位） | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 26,579,020 | 1 |
| 批量辅助 ②（五枚事实 / pid 报告 / `maximumRequests` 对照） | `deep-reasoner` | `GLM-5.3-FlashX (glm)` | 早于 27,045,268 | 1 |
| ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 27,045,268（结果） | 1 |
| ②a 契约 | `prfrail-implementer` | `DeepSeek V4.1 Flash`（主控同源，白名单内） | 27,363,077 | 1 |
| ②b 代码 | `prfrail-implementer` | 同上 | 27,470,881 | 1 |
| ③ 测试 | `prfrail-tester` | 同上 | 27,531,269 | 1 |
| ④ 门禁 | 主控 | — | ③ 返回后（无派单） | 3 作用域 + selftest |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | 27,692,810 | 1 |
| 整改 r1 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 27,958,380 | 1 |
| ⑤ r2 | `independent-reviewer` | `GLM-5.3 (glm)` | 28,076,333 | 1 |
| 整改 r2 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 28,356,341 | 1 |
| ⑥ 独立扫描 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 28,427,012 | 1 |
| ③ 追加测试 | `prfrail-tester` | `DeepSeek V4.1 Flash` | 28,491,264 | 1 |
| ⑦ r1（盲审） | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 28,567,553 | 1 |
| 整改 r3 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 28,679,137 | 1 |
| ⑦ r2 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 28,745,954 | 1 |

**额度合规一行**：① 1/1（必须且仅 1）、② 5（≥2）、③ 2（≥1，且偏移 27,531,269 早于 ⑤ 的 27,692,810 与 ⑦ 的 28,567,553）、⑤ 2（1+1）、⑥ 1（额度 3）、⑦ 2（S 档 1+1，不降级）、⑧a / ⑧b 主控同源。

**批量辅助四个候选点逐项登记**：

- 候选 ③（CN↔EN 对位）= **未用**：`G3(a)` 与 ⑤ / ⑦ 两轮逐句对读已覆盖，批量辅助只会重复，且结论不可机械复用。
- 候选 ④（§7 引用可解析性）= **未用**：本片未新增跨节引用，⑦ / ⑥ 均未提出引用问题。
- 两次批量辅助均用 `GLM-5.3-FlashX`，**未回退 Flash**。
- 偏移读数的一处如实限定：批量辅助 ① 的 26,579,020 是仓内脚本命中的 marker 首次出现位置（该位置落在 `memory` 工具的候选清单载荷内），不是 `runSubagent` 调用点；① 设计的 27,045,268 紧邻「现派 ①」派单消息。该限定不改变次数与额度结论（见 §16 第 1 条）。

## 11 明确未执行事项

- 未执行 `git commit` / `git push`；未触 gitee；CI 未触发。
- 未改 `docs/DELIVERY_DIRECTIVE{,_EN}.md`（准则仍 v1.41）。
- 未改 `internal/adapters/ai_probe_copilot.go` 与 `internal/chain/postflight.go`（路径 D；OB-58 移出）。
- 未改 `internal/console/**`、`tools/b4-e2e/**`、`schemas/**`、`testdata/**`、`.github/workflows/**`、`tools/gates/lib/criteria/**`、`tools/gates/cjk-newwords.txt`。
- 未修 ⑤ r1 High / ⑥ H1–H2 / ⑦ r1 Critical（同一读者背离 ⇒ OB-110）。
- 未整改 `internal/**` 三处旧注释与 `agent_runner_run_test.go:347-349`（越界）。
- 未 bump `stubVersion`（有意决定）。
- OB-58 未在本片立项（OB-109 待立项）。
- ⑩ 停点回执未执行；`git commit` / `git push` / CI 仍未执行，等授权。

## 12 下一步建议

1. 9 项路径暂存后 `--all --scope=index` 已转绿（`TOTAL_FAIL=0`，2026-10-21 实测），`G5(f)` 的暂存前红已闭合。
2. ⑩ 停点等授权；授权后只推 `origin`（不推 gitee），并观察 CI 双绿。
3. 下一片起跑时由 ① 复核 OB-109（推荐形态丙 + 扩边界到 5 个适配器文件）。
4. OB-110 须与 stub 写序决策捆绑立项；「外部停机 vs 报告缺失」的优先级（`cancelled` vs `uncertain`）由未来切片裁定。
5. `internal/**` 旧注释与 `stubVersion` 标签在同时触及两侧的切片内一次性配对处置。

## 13 执行自由度

逐条引用 `docs/validation/evidence/CONTRACTS-QUOTA-ALIGN-freedom-list.md`（清单为权威，本表不改变其含义）。

| # | 选择 | 落点 |
|---|---|---|
| F1 | OB-52 取**路径 D**（分层口径、零代码、零新记录字段） | `docs/CONTRACTS.md:192` |
| F2 | 分层句以「请求预算独立执行 + CLI 额度为候选平台事实 + argv 摘要证明实际传参」三段承载 | `docs/CONTRACTS.md:192` |
| F3 | 费用约束以加粗约束句承载，不新增记录字段 | `docs/CONTRACTS.md:192` |
| F4 | pid 报告文件名以括注声明为「当前 pinned stub 约定、非 wire 契约」 | `docs/CONTRACTS.md:232` |
| F5 | 三态语义与互推禁令合并在同一物理句内 | `docs/CONTRACTS.md:232` |
| F6 | spawn 写失败选「杀掉已起子进程 + 留缺失报告」，而非留空报告或直接返回 | `tools/agent-stub/main.go:187-192` |
| F7 | 保留「（恒定写）」/`(constant write)` 括注（⑦ r1 整改后仍保留） | `docs/CONTRACTS.md:232` |
| F8 | 测试以「报告路径为目录」制造写失败（Windows 可移植，不依赖权限位） | `tools/agent-stub/main_test.go:263` |
| F9 | `stubVersion` 维持 `agent-stub-0.1.0`，不 bump | `tools/agent-stub/main.go:28` |

## 14 未同步项清单

| # | 未同步项 | file:line | 归属 / 处置 |
|---|---|---|---|
| U1 | 树验证腿把缺失报告读作 `nil, nil`，`:284` 因此返回「已消失」 | `internal/adapters/agent_runner_run.go:300-322`（`:306` 为 `return nil, nil`） | OB-110（待议 / 待授权） |
| U2 | 同型宽松读法与文件名直读（显示语义） | `tools/b4-e2e/ports.go:495-503`、`tools/b4-e2e/runmap.go:91` | OB-110 |
| U3 | 旧注释仍称 stub「仅在派生子进程时写子 pid 报告」 | `internal/adapters/agent_runner_postflight_test.go:489-492`、`internal/console/runtime_agent_e2e_test.go:457-460` 与 `:643-646` | 越界，记录在案 |
| U4 | 测试以「无报告 ⇒ 无子进程」为断言 | `internal/adapters/agent_runner_run_test.go:347-349` | OB-110 同源 |
| U5 | `stubVersion` 仍为 `agent-stub-0.1.0` | `tools/agent-stub/main.go:28` | 有意决定（见 §3） |
| U6 | `--scope=index` 的 `G5(f)` 暂存前为红 | 门禁输出 | 已闭合：暂存后 `--scope=index` 读出 `TOTAL_FAIL=0`（2026-10-21 实测） |
| U7 | 本片未落盘审查输入包，盲审隔离无 MANIFEST 可证 | `tmp/CONTRACTS-QUOTA-ALIGN/`（无 `review-*`） | 如实登记：审查输入包以 `tmp/CONTRACTS-QUOTA-ALIGN/01-design.md` 与 `02-diff.patch` 落地（临时件，已 gitignore）；盲审隔离由派单顺序论证，不由 MANIFEST 证明 |

## 15 OB 登记（OB-108 / OB-109 / OB-110；已在 ⑧b 定稿）

- **OB-108** — 试点顺延第 15 次。四要素：(i) 被顺延的是什么 = 新派单形态试点的启用，第 15 次连续顺延；(ii) 为何 = 预算与切片 S 档；(iii) 影响 = 该试点仍未经检验，规则不变；(iv) 解除条件 = 下一个承载试点型任务的切片。附注：「不得作为先例」；「超出 §1.7.4 字面标准、系用户直接裁定」。
- **OB-109** — 「**待立项**：OB-58 的形态与边界已由 `CONTRACTS-QUOTA-ALIGN` 的 ① 充分论证（丙=可重验、需 5 适配器 + 6 处 chain；乙=不可重验、需 1 适配器；甲=丙但漏父报告）。**推荐丙 + 扩边界到 5 个适配器文件**，由下一片起跑时由 ① 复核并给最终方案；本片不做。触发 = 用户排下一片时裁定。」**排片时点 = 池 D `AGENTRUNNER-PROD-PATCH` 完成之后**，与 OB-110 合并处理；届时由 ① 复核形态丙的边界（5 适配器 + 6 处 chain）并给最终方案。
- **OB-110** — 宽松读者：(a) 修复必须与 stub 写序决策捆绑；(b) 「外部停机 vs 报告缺失」的优先级（`cancelled` vs `uncertain`）**不由**本片契约文本裁定，须由未来切片裁定；(c) 验收含「缺失报告 ⇒ uncertain」的端到端测试。范围：`internal/adapters/agent_runner_run.go:300-322`（判定通过点 `:284`）、`tools/b4-e2e/runmap.go:91`、`tools/b4-e2e/ports.go:495-503`；配方 = 返回 `present` 标志，`present == false` ⇒ `treeGone == false` ⇒ uncertain。状态：**待议 / 待授权**。**限定（用户 2026-10-21 裁定）**：**本片已修的是写侧（stub 恒定写），读者侧的宽松读法仍未修**——`agent_runner_run.go` 的 `present` 标志把「文件缺失」读作「未 spawn」，与 CONTRACTS §7 的三态语义不一致。**本片的「恒定写」保障在读者侧只部分生效**：当 stub 正常写空文件时，读者侧不会误判；当 stub 写入失败（进程崩溃等）时，读者侧仍可能把缺失误读为「无子进程」。**该残余与 OB-109 合并处理，本片不修。**

**待补登**：本报告 §15 的 OB-108 / OB-109 / OB-110 三项，须在下一片（池 D `AGENTRUNNER-PROD-PATCH`）的步骤 0 补登至 `DELIVERY_DIRECTIVE{,_EN}.md` §12.4；本片因协议先行边界不改准则。补登时 **OB-110 的限定句须原样带过去**（「本片已修的是写侧，读者侧的宽松读法仍未修……本片不修」）。

## 16 ⑧a 自查三条

1. **成本表每行有机械依据**：表中 16 条记录含 12 条派单偏移，均由 `tmp/CONTRACTS-QUOTA-ALIGN/transcript-offsets.js` 复算、逐条等值；④ 行为主控无派单；批量辅助 ① 的偏移为 marker 命中位置（非调用点），已在 §10 如实标注。不以叙事估算替代读数。
2. **每处「已处置 / 已修复」均带 `file:line`**：本片**无 commit**，故全部处置声明只引 `file:line`（§1 代码落点表、§5 各行、§14 U1–U5）。未使用 commit id。
3. **交付前旧口径全仓检索**：检索 `显式费用/请求预算`、`explicit cost/request budget`、`未 spawn 子进程`、`报告缺失按`、`spawned nothing`、`missing report`。命中与归属：`docs/CONTRACTS.md:192` 与 `docs/CONTRACTS_EN.md:186` 为改写后的现行句（含新约束，非旧口径）；`internal/adapters/agent_runner_run.go:301`、`internal/adapters/agent_runner_run_test.go:347-349`、`tools/b4-e2e/ports.go:495`、`internal/console/runtime_agent.go:674-679`、`internal/adapters/agent_runner_postflight_test.go:489-492`、`internal/console/runtime_agent_e2e_test.go:457-460` 与 `:643-646` 为旧口径命中，**全部在边界外**；边界内旧口径命中 0。

DOCS: DONE

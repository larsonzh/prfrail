# 切片 `AGENTRUNNER-PROD-PATCH` 验证报告（准则 v1.41）

> 切片：`AGENTRUNNER-PROD-PATCH`；分级 `[PILOT]`、档位 M；准则 v1.41（本片只在 §12.4 补登三行，未改准则正文）；本文件由 ⑧a 起草、⑧b 按 ⑨ 终局读数回填（提交状态停在「已暂存、未提交」）。
> 基线：`d:\LZProjects\prfrail`；HEAD = `6eb69f1`（池 C `CONTRACTS-QUOTA-ALIGN` 闭环）。⑧a 启动时暂存集为空，工作区 8 项已修改、0 项未跟踪（`git status --porcelain` 实测）；⑧a 落盘后为 11 项。
> 双语：本文件为 CN 权威；`docs/validation/agentrunner-prod-patch_EN.md` 为逐行同位镜像，两侧行数相等（机械计数见 §7.2）。

## 1 变更概要

**裁定 4 归档说明**：变异证据包见 `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/`（`MANIFEST.md` 列必读项与排除项、`SHA256SUMS.txt` 列三件哈希），`tmp/` 已清空。

⑧a 之前 8 项路径（全部已修改，无未跟踪项）：步骤 0 的三项文档（准则 §12.4 三行、P18 状态块、池 C 的「待补登」句）与 ② 的两项代码（dispatcher 及其测试）。⑧a 之后 11 项（新增本报告对与自由度清单）。下表增删为 `git diff --cached --numstat` 的机械读数（⑨ 暂存快照；⑧b 回填产生的未暂存增量不在本表）；三个新增文件为整文件新增，其增删数取自身行数。

| 路径 | + | − | 说明 |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 3 | 0 | §12.4 三行：OB-108 / OB-109 / OB-110 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 3 | 0 | `_EN` 同位 |
| `docs/t027/REMAINING_SLICES.md` | 23 | 0 | P18 状态块（🔄 进行中） |
| `docs/t027/REMAINING_SLICES_EN.md` | 23 | 0 | `_EN` 同位 |
| `docs/validation/contracts-quota-align.md` | 1 | 1 | 池 C §15 的「待补登」段改为「已补登」 |
| `docs/validation/contracts-quota-align_EN.md` | 1 | 1 | `_EN` 同位（`:339`） |
| `internal/adapters/agent_runner_replay_dispatcher.go` | 49 | 5 | OB-57 回收落点 |
| `internal/adapters/agent_runner_replay_dispatcher_test.go` | 516 | 0 | 10 个新测试函数 |
| `docs/validation/agentrunner-prod-patch.md` | 368 | 0 | 本文件 |
| `docs/validation/agentrunner-prod-patch_EN.md` | 368 | 0 | 逐行镜像 |
| `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-freedom-list.md` | 45 | 0 | 执行自由度清单 |
| `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/`（5 文件） | — | — | 变异证据包（`MANIFEST.md` / `SHA256SUMS.txt` / `evidence.txt` / harness / pristine），终局 numstat 见 §9 | 执行自由度清单 |

artifacts 围栏（16 项路径，全部 `repo`）列在文末。

步骤 0 的门禁读数（落盘后、② 之前）：`G6-6` PASS（`scanned 6`）、`G6-7` PASS、`G3(a)` PASS（准则六形态全 `true`）、`G4b` PASS（`newCjk=276 unknown=[]`）。

代码落点（行号为工作区实测）：

| 文件:行 | 变更 | 依据 |
|---|---|---|
| `agent_runner_replay_dispatcher.go:11` | 引入 `internal/guard`（`guard.TerminationEvidence`） | 停止证据类型 |
| `agent_runner_replay_dispatcher.go:29-33` | sentinel 注释改写：删去「must never kill」，改为「经启动器的停止能力回收 + 不重试不重启 + 保留记录的结清归后续切片」 | ① 路径 A |
| `agent_runner_replay_dispatcher.go:70` | 私有可选能力接口 `agentRunnerLaunchReclaimer{StopAgentRunnerProcess(ctx, launchID) (guard.TerminationEvidence, error)}` | 启动器所有权 |
| `agent_runner_replay_dispatcher.go:78-101` | `reclaimUnprovenLaunch(ctx, launchID, cause)`：类型断言；无能力即显式「未回收」；`grace := dispatcher.Process.Grace`，`<=0` 取 1s；`context.WithTimeout(context.WithoutCancel(ctx), grace+time.Second)`；单次停止；三个出口（停止失败 / 结论非 `stopped` / 已回收） | 有界 + cause 不被替换 |
| `agent_runner_replay_dispatcher.go:205` / `:216` / `:220` | spawn 后三条失败分支（launchId 不匹配 / 身份记录构造 / 身份发布）全部经该 helper，且一律传 dispatcher 自证 launchID | ⑥ High 的裁定见 §5 |
| `agent_runner_replay_dispatcher.go:83` | `reclaimer == nil` 守卫（连 `!ok`）上方加防御性注释（今日不可达；任何使其可达的改动必须补测试） | 用户裁定 |

未改动：`internal/chain/**`、store 文件（含 `agent_runner_replay_store_launches.go`）、`internal/console/**`、`tools/agent-stub`、`docs/CONTRACTS{,_EN}.md`、`tools/gates/lib/criteria/**` —— 在 `git diff --numstat` 中均零出现（机械读数）。

## 2 执行流水线

| 环节 | 角色 / 载体 | 模型 | 结论 | 轮次 | 备注 |
|---|---|---|---|---|---|
| 批量辅助 ① | `quick-researcher` | `GLM-5.3-FlashX (glm)` | 完成 | 1 | dispatcher 的 spawn / identity 定位与摘录 |
| 批量辅助 ② | `quick-researcher` | `GLM-5.3-FlashX (glm)` | 完成 | 1 | 错误分支枚举 + launcher / store API + 测试覆盖 |
| ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 完成 | 1 | 四路径比较后取路径 A；四组 G1–G4；列 8 项风险 |
| ② 实现 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 完成 | 1 | `+49/−5` |
| ③ 测试（首轮） | `prfrail-tester` | `DeepSeek V4.1 Flash` | 完成 | 1 | 10 个测试函数，`+516` |
| ④ 门禁 | 主控 | — | 完成 | 3 作用域 + selftest | 无派单 |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | PASS | 1 | 4 项 Low（见 §5） |
| ③ 追加测试 | `prfrail-tester` | `DeepSeek V4.1 Flash` | 完成 | 1 | 关闭 ⑤ 的 (b)(c) |
| ⑥ 独立扫描 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | FINDINGS 5 | 1 | 1 High + 3 Medium + 1 Low |
| ⑦ r1（盲审） | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | 无发现 |

- 返工口径：③ 共 2 次派单（首轮 + 追加测试）外加 D1 要求的变异证据重出 1 次；② 未重派、未回滚；④ 无重跑。
- 耗时字段历史缺失：本片未落盘起止时间戳，不补造。
- 顺序合规：③ 首轮偏移 30,305,273 早于 ④（主控、③ 返回后）与 ⑤ 的 30,498,825（§10）。

① 的四组裁定与路径取舍：

- 路径 A（三条 spawn 后失败分支就地回收）选定。
- 路径 B（spawn 前预留身份）否决：store 无 update / delete identity API（全仓检索零命中），B 需三个新 API 并改 CONTRACTS，超出本片边界；B 还会把活进程泄漏转换成阻塞重派的持久预留。
- 第三条路径否决：console 的 `agentRunnerMitigateAbandonedLaunch` 自述「非根因修复」且越界；直接调 `guard.StopProcessIdentity` 会绕过启动器所有权（handle / slot / log 均不释放）。
- G1 错误路径：三条分支回收；spawn 前的 `Start` 失败不得回收（启动器契约保证无进程）。
- G2 回收：只经启动器的停止能力，`grace+1s` 有界 ctx，单次尝试，cause 不被替换。
- G3 身份：禁止重试（no-replace store，重试只会暴露外来身份），不删启动意图，不覆盖既有身份。
- G4 测试面：stub 级矩阵 + 一个观察 OS 事实的真实进程测试。

## 3 异常与兜底记录

**裁定 5 登记**：`buildAgentRunnerParityStub` 的 `os.TempDir()` 行为**保持「候选 OB」形态、不升格为正式 OB 号**——它是既有测试助手、不在本片边界内，升格会开出一个本片无法闭合的 OB；触发条件 = 该 helper 下次被改动。

| # | 事件 | 影响 | 处置 |
|---|---|---|---|
| D1 | ③ 的变异轮在 `%TEMP%` 下建 detached worktree，用后删除（自报） | 违反 §5.1 的落点规则。计数口径 = 按规则合并（`$env:TEMP` 与 `tmp/` 根同属「临时件一律 `tmp/<切片ID>/`」的违反）⇒ 本项为该规则的**第 7 次**（池 C 已在 6 次关闭）；片内升级计数仍为「本片第 1 次，触发在本片第 3 次」。不得读作全局首次，计数不得从 1 重新起算 | 用户裁定：继续本片、计第 7 次；变异证据须在 `tmp/AGENTRUNNER-PROD-PATCH/mutation/` 重出完整链（§7 即因此存在）。重出计 ③ 1 次返工 |
| D2 | 首个 M2 变异把条件替换为 `if true {` | 编译失败：`declared and not used: ok` | 改用可编译变体（类型断言后追加 `ok = false`）；工具迭代，非缺陷 |
| D3 | `%TEMP%/proofrail-agent-stub-parity`（今日 09:50 创建） | 来自既有测试助手 `buildAgentRunnerParityStub` 写 `os.TempDir()`，非本片派单 | 登记为候选 OB（触发 = 该助手下次被改动），故意不给 OB 号：属既有工具行为，不是本片缺陷 |
| D4 | `-race` 本机不可执行：`go: -race requires cgo` / `cgo: C compiler "gcc" not found` | 本机无 `-race` 证据 | 环境事实；验收证据只能来自 CI 的 ubuntu 腿，不得写成通过 |
| D5 | ④ 的 `G5(d)` 红与 `G5(f)` 红 | 报告未落地 / 尚未暂存的结构状态 | 非缺陷：本报告落盘后 `G5(d)` 转绿，`G5(f)` 由 ⑨ 暂存消解 |

非偏离但一并登记（避免读作遗漏）：

- ⑥ 的五条发现在本片均未修复，逐条登记见 §5。
- `internal/console/runtime_agent.go` 两处陈旧注释越界（§14 U1）。
- 变异证据已按用户 2026-10-21 裁定 4 归档至 `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/`（`MANIFEST.md` + `SHA256SUMS.txt`，`evidence.txt` 为必读核心证据），`tmp/` 已清空至仅 `.gitkeep`；两个可重建的大件（`head-tree.tar`、`wt/`）未入库，重建命令写在 `MANIFEST.md` 的「Deliberate exclusions」节。

## 4 环境与档位

| 项 | 值 |
|---|---|
| OS | Windows（本机 `d:\LZProjects\prfrail`） |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| 档位 | M（`[PILOT]`，触进程回收语义 ⇒ ① 与 ⑦ 均不降级） |
| ⑦ 额度 | 1+2（M 档） |
| ⑤ 额度 | 1+1 |
| 批量辅助 | ADR-027（限四类任务）；两次均用 `GLM-5.3-FlashX` |
| 探针 | 0 |
| `-race` | 本机不可执行（无 cgo / gcc）；仅 CI ubuntu 腿可证 |

## 5 审查结论摘要（⑤⑥⑦）

| 轮次 | 编号 | 级别 | 结论 / 处置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 | L-a | Low | `internal/console/runtime_agent.go:296-299` 与 `:851-861` 的旧注释仍把 OB-57 写成由他片负责的未修根因 | 记录在案（越界，§14 U1） |
| ⑤ r1 | L-b | Low | `WithoutCancel` 屏蔽没有测试 | 已整改（`ReclaimSurvivesCallerCancellation`） |
| ⑤ r1 | L-c | Low | typed-nil 回收器无测试 | 已整改（`TypedNilReclaimerIsReportedFailClosed`） |
| ⑤ r1 | L-d | Low | `reclaimer == nil` 子句不可达 | 保留 + 用户要求的防御性注释（`agent_runner_replay_dispatcher.go:83`） |
| ⑥ | H1 | High | 回收目标应为 `result.LaunchID` | 不采纳：⑦ 独立裁定相反（见下方裁定段） |
| ⑥ | M1 | Medium | `StopAgentRunnerProcess` 内 panic 会逃逸回收 | 登记：启动器实现面，dispatcher 契约外，⑦ 未提出 |
| ⑥ | M2 | Medium | typed-nil 不是 fail-safe | 登记：⑦ 引用确认接受当前形态 |
| ⑥ | M3 | Medium | `WithoutCancel` / nil ctx / 零 grace 语义 | 登记：⑦ 引用确认接受当前形态 |
| ⑥ | L1 | Low | 停止错误用 `%v` 而非 `%w` | 保留：单 sentinel 错误链的刻意选择，⑤ 背书 |
| ⑦ r1 | — | — | PASS，无发现；明确以 slot 语义裁定回收目标 | 已闭合 |

- ⑥ 的 High 不采纳的依据：⑦ 盲审（未附 ⑤ / ⑥ 结论）独立裁为「回收目标 = dispatcher 自证 id」——释放路径按传入 id 查 handle / slot，改传 `result.LaunchID` 会让外来启动器改写停止目标。
- ⑥ 的三条 Medium 属防御性加固的判断项：其中 typed-nil 与 ctx / grace 两项 ⑦ 已审阅并带引用接受；第三条（启动器停止内的 panic）属启动器实现，在 dispatcher 契约之外，⑦ 未提出。
- ⑥ 的 Low（`%v`）是刻意选择：保持错误链只有单个 sentinel，⑤ 已背书。
- ⑥ 的五条发现在本片均未修复，全部在此登记。
- ⑥ 的五条发现全部登记、无一修复；其中 High 由 ⑦ 以 slot 语义裁定推翻，该推翻即为本片记录在案的裁决（见上方裁定段与 §14 U5）。
- ⑦ 的额度：M 档 1+2 中只用 1 次，因 r1 即 PASS 且无发现 ⇒ 记为「额度未用尽（r1 即 PASS）」，不得读作跳过复审。

## 6 审查输入包摘要

本片未落盘审查输入包：`tmp/AGENTRUNNER-PROD-PATCH/` 下无 `review-*` 目录、无 `MANIFEST.md`。

| 项 | 读数 |
|---|---|
| `tmp/AGENTRUNNER-PROD-PATCH/` 顶层文件 | 2（`offsets.js` 873 字节、`step0.js` 14006 字节） |
| `mutation/` 内容 | `evidence.txt` 5290 字节、`dispatcher.pristine.go.txt` 12104 字节、`head-tree.tar` 26531840 字节、`run-mutations.js` 5770 字节、`wt/` |
| `review-5/`、`review-6/`、`review-7/` | 均不存在 |
| `tmp/` 根 | 仅 `.gitkeep` 与本片目录 |

- ⇒ 盲审隔离无法用 MANIFEST 证明；派单事实（盲输入、未附前轮结论）只存在于派单文本，⑧b / ⑨ 不能机械复核。该缺口登记为 §14 的未同步项，不以「三审均通过」掩盖。
- 可机械核对的是变异证据链：`evidence.txt` 记录基线 exit 0、两次变异的红、以及两次还原后 sha256 回到 `f142744…`（与工作区实测一致，见 §7.2）。
- 输入消毒的其余部分（brief / 设计 / diff / 门禁读数 / 自测输出）无落盘副本，本报告不复述其内容。

## 7 可证伪性

**裁定 6 澄清**：M2 首次写法 `if true {` 触发 `declared and not used: ok` **编译失败**——编译失败本身**不构成证伪证据**；证伪证据是「可编译变体（类型断言后追加 `ok = false`）使 8 项测试变红」。该次迭代登记为工具调整，**不计入发现数、不计入返工数**。

### 7.1 可证伪断言

| # | 断言 | 反例（删去后应有哪条读数变化） | 机械检查 |
|---|---|---|---|
| E1 | 三条 spawn 后失败分支均经回收 helper | 任一分支退回直接 `return`，对应回收测试转红（M2 变异实测 8 项红） | `reclaimUnprovenLaunch(ctx, launchID` 出现 3 处调用（实测 3） |
| E2 | spawn 前 `Start` 失败不回收 | 若误加回收，`StartFailureRetainsIntentWithoutRelaunch` 会因停止调用而变红 | 该测试未改动且全绿（§8） |
| E3 | 屏蔽调用方取消的 ctx 有测试钉住 | M1 变异（`WithoutCancel` 改回 `ctx`）使 shield 测试红 | `SELFTEST` 之外由 M1 实测，见 §7.3 |
| E4 | 不重试、不删意图、不覆盖身份 | 若加回收后重试，身份冲突测试会出现二次写入 | 既有身份冲突测试未改动且全绿 |
| E5 | dispatcher 不直接操作 OS 进程 | 若直呼 guard 停止原语，启动器所有权旁路 | 停止调用点唯一：`agent_runner_replay_dispatcher.go:92` |

### 7.2 机械读数复算

| 读数 | 值 | 来源 |
|---|---|---|
| `git rev-parse --short HEAD` | `6eb69f1` | 命令输出 |
| `git diff --numstat` | 8 项（⑧a 前）/ 11 项（⑧a 后） | 命令输出 |
| dispatcher 行数 / sha256 | 行数 269 / `f142744…` | `node` 复算 |
| 测试文件行数 / 字节 / sha256 | 行数 932 / 字节 40767 / `fa40341…` | `node` 复算 |
| 新增测试函数 | 10（另有 11 个既有测试未改动） | 机械计数 |
| `SELFTEST` | SELFTEST: PASS 588/588 | `node tools/gates/gate.js --selftest` |
| `G4b` | `newCjk=602 unknown=[]` | `--check=G4b --scope=tree` |
| CN / `_EN` 报告行数 | 两侧相等（见 §8） | 机械计数 |
| transcript 派单偏移 | 10 条（§10） | 本会话派单点 |

### 7.3 变异载荷

在 `tmp/AGENTRUNNER-PROD-PATCH/mutation/` 的独立树上执行（`git archive` 取 HEAD 的干净树，再覆盖两个工作区文件；仅解开 `go.mod go.sum internal cmd tools examples schemas testdata harnesses`，因为随 Windows 自带的 tar 无法处理历史证据归档中的带圈数字文件名）。

| 变异 | 预期 | 实测 |
|---|---|---|
| 基线 `go test -count=1 -run "AgentRunnerReplayDispatcher" ./internal/adapters/` | exit 0 | exit 0（`ok 4.545s`） |
| M1：`context.WithTimeout(context.WithoutCancel(ctx), …)` → `context.WithTimeout(ctx, …)` | shield 测试红 | 红：`shield failed: … Err()=context canceled` |
| M2：类型断言后追加 `ok = false` | 回收矩阵红 | 红：8 项（7 顶层 + 1 子测） |
| 两次还原后的 sha256 | 等于变异前 | 均 `f142744…`（等于变异前 = true） |

- 首个 M2 尝试把条件替换为 `if true {`，因 `declared and not used: ok` 编译失败，改用可编译变体 ⇒ 如实记为工具迭代（D2），不是发现。
- 变异只做证伪，不进交付：工作区两文件保持整改后形态，sha256 与证据链一致。

### 7.4 覆盖边界与已知 Low

| 编号 | 内容 | 归属 |
|---|---|---|
| L-1 | 重复停止是否安全（幂等性） | 记录：单次尝试 + 启动器契约承担 |
| L-2 | stub 级断言不等于 OS 事实 | 关闭：真实进程测试（`agent_runner_replay_dispatcher_test.go:718`） |
| L-3 | ctx 上界可能偏紧 | 记录：`grace+1s`，默认 grace 由 `Process.Grace` 决定 |
| L-4 | 无停止能力的启动器 | 关闭：显式「未回收」出口（`agent_runner_replay_dispatcher_test.go:685`） |
| L-5 | 错误链是否断裂 | 记录：cause 始终为首 sentinel，`errors.Is` 断言覆盖 |
| L-6 | 回收目标 id 的选择 | 关闭：⑦ 裁定 dispatcher 自证 id（§5） |
| L-7 | Windows 上的时序问题 | 记录：`-count=3` 全绿；`-race` 只能 CI |
| L-8 | 外来身份仍指向已死进程 | 记录：no-replace store，本片不删不覆盖 |

- 本报告的结论不覆盖 `internal/console/**` 与 `internal/chain/**` 的语义等价性：二者在本片零改动，只作为未同步项登记。
- CN 与 `_EN` 的语义等价由 ⑤ / ⑦ 的逐句对读与 `G3(a.2)` 形状比较共同承担。

## 8 门禁结果

| 命令 | 结果 |
|---|---|
| `node tools/gates/gate.js --selftest` | SELFTEST: PASS 588/588 |
| `node tools/gates/gate.js --check=G3-a --scope=tree` | TOTAL_FAIL=0（原文见下） |
| `node tools/gates/gate.js --check=G4b --scope=tree` | TOTAL_FAIL=0（原文见下） |
| `node tools/gates/gate.js --check=G5-d --scope=tree` | TOTAL_FAIL=0（原文见下） |
| `node tools/gates/gate.js --check=G5-e --scope=tree` | TOTAL_FAIL=0（原文见下） |
| `node tools/gates/gate.js --all --scope=tree` | TOTAL_FAIL=0（三作用域全绿；⑨ 暂存后实测，见下） |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test -count=1 ./...` | 全部包 ok，无 FAIL |
| `go test -count=1 ./internal/adapters/` | ok |
| `go test -count=3 -run "AgentRunnerReplayDispatcher" ./internal/adapters/` | ok（3 次） |
| `go test -race ./internal/adapters/` | 不可执行（D4），无读数 |

④ 的历史读数：`--scope=ci` 为 `TOTAL_FAIL=0`；`--scope=tree` 为 `TOTAL_FAIL=1`，唯一失败是 `G5(d)`——其 `stagedNotDeclared` 恰为准则对与两个代码文件，因为当时变更集里还没有报告声明本片变更集。

- `G5(d)` / `G5(f)` / `G7-b` 三项红已全部关闭：前两者是「报告未落盘 / 尚未暂存」的**结构性红**（变更集里没有报告声明本片变更集；新工件不在索引），`G7-b` 是新工件未进索引；⑨ 精确暂存 11 路径后三项同时转绿，`--scope=tree` / `--scope=index` / `--scope=ci` 均为 `TOTAL_FAIL=0`（exit 0），`G6-6`（`scanned 6`）与 `G6-7` 亦 PASS。
- 逐判据原文（⑨ 暂存后、`--scope=tree` 实测）：

```text
PASS  G3(a) symmetry + mirror parity :: docs/DELIVERY_DIRECTIVE.md +3/-0 vs _EN +3/-0 OK || docs/t027/REMAINING_SLICES.md +5/-2 vs _EN +5/-2 OK || docs/validation/contracts-quota-align.md +1/-1 vs _EN +1/-1 OK || docs/validation/agentrunner-prod-patch.md +40/-28 vs _EN +40/-28 OK ;; docs/DELIVERY_DIRECTIVE.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true} ;; docs/t027/REMAINING_SLICES.md {"lines":false,"head":false,"bold":false,"bars":false,"blank":false,"backticks":false} ;; docs/validation/contracts-quota-align.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true} ;; docs/validation/agentrunner-prod-patch.md {"lines":true,"head":true,"bold":true,"bars":true,"blank":true,"backticks":true}
PASS  G4b CJK diff minus whitelist :: newCjk=602 unknown=[] (whitelist=37 whitelist-source=worktree)
PASS  G5(d) staged is a subset of declared :: 11 path(s) declared (scope=tree) ; declared-not-touched=5 (INFO: a fence is the landing snapshot) ; addMode=explicit (derived label, not a command-history assertion)
PASS  G5(e) evidence references resolve :: no added completion row (nothing to resolve)
```

## 9 ⑨ 原生验证与提交状态

⑨ 已执行：主控在同一轮授权下逐个文件暂存 11 项路径，下表为终局读数。

| 项 | ⑨ 终局读数 |
|---|---|
| 暂存集 | 16 路径精确暂存（含证据包 5 文件）；`added 1952 / deleted 7`（逐项 `3/0`、`3/0`、`23/0`、`23/0`、`368/0`、`368/0`、`1/1`、`1/1`、`45/0`、`49/5`、`516/0`） |
| 三作用域 | `--scope=tree` / `--scope=index` / `--scope=ci` 均 `TOTAL_FAIL=0`（exit 0） |
| `SELFTEST` | `SELFTEST: PASS 588/588` |
| `go test -count=1 ./...` | 15 个包 `ok`，无 FAIL（adapters 36.404s、console 12.753s、agent-stub 6.689s） |
| 编码审计（11 项路径） | `violations=0`（九个 `.md` 均 BOM + CR=0，两个 `.go` 均无 BOM + CR=0），合计 1,333,641 字节 |
| `-race` | 本机不可执行（`go: -race requires cgo` / `cgo: C compiler "gcc" not found`），属环境事实；验收证据由 CI 的 ubuntu 腿提供（run id 记录在 ⑩） |
| HEAD | `6eb69f1`（未变） |
| 提交与推送 | 未提交、未推送；未触 gitee；CI 未触发 |
| ⑩ 停点回执 | 未执行（等授权） |
| artifacts 围栏 | 列 16 项路径，与 §1 的表逐项一致 |

- ⑨ 逐个文件暂存（非 `git add -A`），三作用域、`SELFTEST` 与全量测试均复跑于暂存之后。
- 口径说明：本行数字是 **⑧c 重新暂存后**的终局快照（已含 ⑧b 对报告对与台账对的回填增量）；⑧b 期间 `--scope=index` 的 `G5(f)` 红已由该轮暂存消解，三作用域现均为 `TOTAL_FAIL=0`。
- 提交与推送停在「已暂存、未提交」：本文件不代行 `git commit` / `git push`，二者留给 ⑩ 的显式授权。

## 10 成本与计量（R7.1）

机械依据 = 本会话 transcript 的字节偏移（本会话重算）。

| 环节 | 角色 | 模型 | 偏移 |
|---|---|---|---|
| 批量辅助 ①（spawn / identity 定位与摘录） | `quick-researcher` | `GLM-5.3-FlashX (glm)` | 29,904,499 |
| 批量辅助 ②（错误分支 + launcher / store API + 测试覆盖） | `quick-researcher` | `GLM-5.3-FlashX (glm)` | 29,970,564 |
| ① 设计 | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 30,072,919 |
| ② 实现 | `prfrail-implementer` | `DeepSeek V4.1 Flash` | 30,229,371 |
| ③ 测试（首轮） | `prfrail-tester` | `DeepSeek V4.1 Flash` | 30,305,273 |
| ④ 门禁 | 主控 | — | ③ 返回后（无派单） |
| ⑤ r1 | `independent-reviewer` | `GLM-5.3 (glm)` | 30,498,825 |
| ③ 追加测试 | `prfrail-tester` | `DeepSeek V4.1 Flash` | 30,699,904 |
| ⑥ 独立扫描 | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | 30,930,888 |
| ⑦ r1（盲审） | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 30,965,196 |

- 顺序合规：③ 首轮 30,305,273 早于 ④（主控、③ 返回后）与 ⑤ 的 30,498,825。
- 额度合规：① 1/1；② 1；③ 2（≥1）；⑤ 1（1+1）；⑥ 1（额度 3）；⑦ 1（M 档 1+2，r1 即 PASS，额度未用尽）；⑧a 与主控同源。
- 模型配额：`GLM-5.3-FlashX` 2 次，`GLM-5.3-Flash` 0 次。
- 批量辅助候选 ③（dispatcher 测试覆盖）与候选 ④（冻结事实与 pid 报告的读取路径）均未派单：候选 ③ 已被批量辅助 ② 的产出覆盖（错误分支与测试覆盖在同一份报告内），重复派无增量；候选 ④ 在本片边界内不存在该类代码（dispatcher 不读冻结事实与 pid 报告，① 的落点表也无此路径）。

## 11 明确未执行事项

- 未执行 `git commit` / `git push`；未触 gitee；CI 未触发。
- 未改 `docs/DELIVERY_DIRECTIVE{,_EN}.md` 的准则正文（只补登 §12.4 三行，v1.41 不变）。
- 未改 `internal/chain/**`、store 文件、`internal/console/**`、`tools/agent-stub`、`docs/CONTRACTS{,_EN}.md`、`tools/gates/lib/criteria/**`、`tools/gates/cjk-newwords.txt`。
- 未修 ⑥ 的五条发现（§5）。
- 未整改 `internal/console/runtime_agent.go` 的两处陈旧注释（越界，§14 U1）。
- 未处理 OB-109 / OB-110（仍待议，排在本片之后）。
- 未执行 `-race`（本机不可执行，D4）。
- 未执行 ⑩ 停点回执、`git commit` 与 `git push`（等授权）；⑨ 的暂存已完成（§9）。

## 12 下一步建议

1. ⑨ 已完成（11 路径逐个暂存，三作用域与 `SELFTEST` 读数已回填 §8 / §9）；⑩ 授权后由主控提交并推送 origin（不触 gitee），随后触发 CI。
2. CI 的 ubuntu 腿必须给出 `-race` 的包级 `ok` 行，作为本片唯一的竞态验收证据；不得以「本机不可执行」代替。
3. 下一片起跑时由 ① 复核 OB-109（形态丙的边界：5 个适配器 + 6 处 chain）并给最终方案。
4. OB-110 与 OB-109 合并处理；「外部停机 vs 报告缺失」的优先级由未来切片裁定。
5. `buildAgentRunnerParityStub` 写 `os.TempDir()` 的行为随其下次改动处置，届时给 OB 号。

## 13 执行自由度

逐条引用 `docs/validation/evidence/AGENTRUNNER-PROD-PATCH-freedom-list.md`（清单为权威，本表不改变其含义）。

| # | 选择 | 落点 |
|---|---|---|
| F1 | 取路径 A（三条 spawn 后失败分支就地回收） | `agent_runner_replay_dispatcher.go:205` / `:216` / `:220` |
| F2 | 停止能力用私有可选接口承载，不扩 chain 端口 | `agent_runner_replay_dispatcher.go:70` |
| F3 | 无停止能力时给出显式「未回收」出口，不静默 | `agent_runner_replay_dispatcher.go:84` |
| F4 | grace 缺省取 1s，停止 ctx 用 `grace+1s` 有界 | `agent_runner_replay_dispatcher.go:86-89` |
| F5 | 停止只尝试一次，不做循环重试 | `agent_runner_replay_dispatcher.go:92` |
| F6 | cause 始终为首 sentinel，停止错误与结论只作尾注 | `agent_runner_replay_dispatcher.go:93-100` |
| F7 | 停止错误用 `%v` 拼接，保持单 sentinel 链 | `agent_runner_replay_dispatcher.go:95` |
| F8 | 不重试身份、不删启动意图、不覆盖既有身份 | `agent_runner_replay_dispatcher.go:216-221` |
| F9 | 补一个真实进程测试观察 OS 事实 | `agent_runner_replay_dispatcher_test.go:718` |
| F10 | `reclaimer == nil` 子句保留，配防御性注释 | `agent_runner_replay_dispatcher.go:83` |

## 14 未同步项清单

**裁定 2(b) 扩展动作 + 裁定 3 的 ⑥ 有效性数据点**：

- **X1（口径澄清，写入复评记录，不改 §12.1 正文）**：「审查方发现」须区分 **⑦ 裁决后成立者** 与 **被推翻者**——本片 ⑥ 的 5 项中 1 High 被 ⑦ 推翻、2 项为防御性取舍、1 项属 launcher 范畴 ⇒ **真正成立的审查方发现 ≤ 4**，分母虚高。**这是记录，不是豁免**：即使按 4 计，自捕获率 ≤ 1/5 = 20%，仍未达标。
- **X2（下一试点的前置改进）**：① 取证性变异**不得**在 `%TEMP%` / 仓库外进行（同一规则已补入仓库记忆 `prfrail-dispatch-order-check.md`）；② 若下一试点仍用 MAI 作 ⑥，须在 ⑧b 记录「⑥ 假阳性数 / 总数」。
- **⑥ 有效性数据点**：本片 ⑥ 报 1 High + 3 Medium + 1 Low，⑦ 裁决后 1 High 被推翻（推导 id 更安全）、2 Medium 经逐条给证被接受为防御性取舍、1 Medium 属 launcher 范畴、1 Low 越界登记 ⇒ **⑥ 的假阳性率为 1/5**。

| # | 未同步项 | file:line | 归属 / 处置 |
|---|---|---|---|
| U1 | 旧注释仍把 OB-57 写成由他片负责的未修根因 | `internal/console/runtime_agent.go:296-299`、`:851-861` | 越界，记录在案 |
| U2 | 既有测试助手写 `os.TempDir()` | `internal/adapters` 内的 `buildAgentRunnerParityStub` | 候选 OB（无号），随其下次改动处置 |
| U3 | OB-109 / OB-110 仍待议 | `docs/DELIVERY_DIRECTIVE.md` §12.4 | 排在本片完成后合并处理 |
| U4 | `-race` 本机不可执行 | 环境 | 只能由 CI ubuntu 腿验收 |
| U5 | ⑥ 的五条发现未修 | §5 | 全部登记，不读作已处理 |
| U6 | 派单简报的 `+54/−7` 与 `git diff --numstat` 实测 `+49/−5` 不一致 | `internal/adapters/agent_runner_replay_dispatcher.go` | 以实测为准；差异如实登记，未按简报数字回写 |
| U7 | 未落盘审查输入包，盲审隔离无 MANIFEST 可证 | `tmp/AGENTRUNNER-PROD-PATCH/`（无 `review-*`） | 如实登记：隔离由派单顺序论证，不由 MANIFEST 证明 |

## 15 试点指标（§12.1 声明）

本片以 `[PILOT]` 起跑，声明如下（逐字取自起跑消息）：

- 缺陷自捕获率（我方发现 ÷ 总发现）≥ 3/4。
- 返工次数 ③ ≤ 2、④ ≤ 2。
- ⑦ 无效轮次 0。
- 主控兜底次数 ≤ 3。
- 用户发现的遗漏数 0。
- 分母口径（§12.1 固定）：总发现 = 我方 + 审查方 + 用户，按条计，同根因合并为 1 条；任一指标未达标即判试点未达标。

| 指标 | 声明目标 | 实测 |
|---|---|---|
| 缺陷自捕获率（我方 ÷ 总发现） | ≥ 3/4 | 1/10 = 10%（**未达标**） |
| 返工次数 ③ | ≤ 2 | 2（达标） |
| 返工次数 ④ | ≤ 2 | 1（达标） |
| ⑦ 无效轮次 | 0 | 0（达标） |
| 主控兜底次数 | ≤ 3 | 1（达标） |
| 用户发现的遗漏数 | 0 | 0（达标） |

- **我方发现 = 1**：`%TEMP%` 临时件违规，由 ③ 自报并随即删除。
- **审查方发现 = 9**：⑤ r1 4 项（Low）；⑥ 5 项（1 High + 3 Medium + 1 Low）；⑦ r1 0 项。
- **用户发现 = 0**：用户明确说明该违规系我方自报，非用户发现。
- **总发现 = 10**（1 + 9 + 0，同根因合并后逐条计）⇒ **自捕获率 = 1/10 = 10%** ⇒ **未达标（目标 ≥ 3/4）**。
- 更严口径：若只计 ⑦ 裁定后仍成立者（⑥ 的 High 被推翻），审查方 = 8 ⇒ 1/9 ≈ 11% ⇒ **仍未达标**。
- 返工次数：**③ = 2**（① 追加测试 1 次 + ② 变异重跑 1 次，后者由用户裁定计入）、**④ = 1**（报告落地后复跑）⇒ 达标（≤2）。
- ⑦ 无效轮次 = **0**（r1 即 `FINAL: PASS`）⇒ 达标；主控兜底 = **1**（临时件违规的停下升级）⇒ 达标（≤3）；用户发现遗漏 = **0** ⇒ 达标。
- **结论：五项中一项未达标（自捕获率）⇒ 按 §12.2 的复评条款处理，不自动回退；须由用户在 ⑩ 裁决。**
- 上表实测值由 ⑧b 按 §12.1 的固定分母复算并回填；⑧a 未预填任何数值。

## 16 ⑧a 自查三条

1. 成本表每行都有 transcript 偏移依据（§10 表），不是叙述估算；偏移取自本会话的派单点，本片未像池 C 那样落盘复算脚本，故记录「未用脚本复算」这一限定。
2. 本片尚无提交，凡「已修 / 已关闭」均引 `file:line`：三条分支 `agent_runner_replay_dispatcher.go:205` / `:216` / `:220`，helper `:78-101`，接口 `:70`，防御注释 `:83`；测试 `agent_runner_replay_dispatcher_test.go:718` / `:685`；全片未使用提交号。
3. 全仓检索本片证伪的语句（排除 `tmp/`、`.git/` 与本报告对自指：报告内引用的被检索串按 `OB-2` 口径排除）：
   - `must never kill`：0 处 ⇒ 允许文件内为 0，达标。
   - `owns settlement`：0 处；替代句为 `settlement of the retained durable records is owned by a later slice`（`agent_runner_replay_dispatcher.go:32-33`），在允许文件内。
   - `OB-57`：53 处（⑧b 复算，排除本报告对与 `tmp/`、`.git/`）；允许文件内 1 处（`agent_runner_replay_dispatcher.go:67` 的新注释 `pre-OB-57 behaviour`，非被证伪句）；越界 2 处（`internal/console/runtime_agent.go:298` / `:856`）；其余 50 处在 `docs/**`（含本片台账完成段）。
   - 因此允许文件内被证伪语句的命中为 0；`internal/console/**` 的 2 处属越界，已登记 §14 U1。

```artifacts
docs/DELIVERY_DIRECTIVE.md
docs/DELIVERY_DIRECTIVE_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
docs/validation/contracts-quota-align.md
docs/validation/contracts-quota-align_EN.md
docs/validation/agentrunner-prod-patch.md
docs/validation/agentrunner-prod-patch_EN.md
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-freedom-list.md
internal/adapters/agent_runner_replay_dispatcher.go
internal/adapters/agent_runner_replay_dispatcher_test.go
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/MANIFEST.md
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/SHA256SUMS.txt
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/evidence.txt
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/run-mutations.js
docs/validation/evidence/AGENTRUNNER-PROD-PATCH-mutation/dispatcher.pristine.go.txt
```

DOCS: DONE

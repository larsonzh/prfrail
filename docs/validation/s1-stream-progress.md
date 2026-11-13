# ProofRail S1-STREAM-PROGRESS 报告（P70 · 层 2 治理侧 · 实现类 [SLICE] · S 档）

日期：2026-11-13。状态：⑧ 文档装配定稿（①–⑦ 已闭环）。切片标识：`S1-STREAM-PROGRESS`（P70 · 层 2 治理侧 · 实现类 `[SLICE]` · S 档）。

## 1. 执行摘要

- 切片：`S1-STREAM-PROGRESS`（P70，层 2 治理侧，实现类 `[SLICE]`，S 档），主控 `DeepSeek V4.1 Flash (deepseek)`。
- 交付：`run --progress` 轻量流式开关，默认关闭；新增 `internal/console/progress.go`。
- 结论：形态 A 落地；默认行为不变；`--json` 语义不变；门禁全绿。
- 计量：⑤ 一次（`PRE-REVIEW: PASS`）、⑥ 一次（`INDEPENDENT SCAN: PASS`）、⑦ 终审 1 + 复审 1（`FINAL REVIEW: PASS WITH FIXES` 之后接 `RE-REVIEW: FINDINGS`）；盲审 0；GLM 批量辅助 2 次。
- 计数：**R = 1、D = 2**（R = 1 = CI 首跑红后的修复重推，§9.2；临时件 D 累计 2 次；均为⑧c 后取证窗口内主控将捕获件误落 `tmp/` 根，已逐件精确清除、零残留，按 §5.1 非自愈档登记并升级用户裁决）。

### 1.1 异常与偏离登记

a) ③ 测试整改期，测试载体曾误用 `git checkout -- internal/console/progress_test.go`，而该文件在索引中为 intent-to-add 空条目 ⇒ 工作副本被清空；经 VS Code Local History 恢复（790 条内容），并以「行号一致 + 恢复后基线全绿 + 改动前 `git diff --stat` 为 790 插入」三证确认与探针前一致。
b) ⑦ 复审的 Low 整改后**未开第 3 轮 ⑦**：S 档额度 = 终审 1 + 复审 1 已用尽，按 §7.8.3「额度用尽」停表处置，整改经机械读数自证（见 §4）。
c) **临时件 D = 1**：⑧c 清空后、⑩ 回执前的取证窗口内，主控把一次门禁捕获写到了 `tmp/` 根（`tmp/g71.txt`，位置不在 `tmp/<切片ID>/` 内），违反 §5.1 的位置闭集；已当场自报、逐件 `Remove-Item -LiteralPath` 删除，并提交零残留读数：`Test-Path` 由 `True` 转 `False`，`tmp` 递归文件计数 102 → 101（排除 `gates-cache/g4b` 后为 1），顶层名回到 `gates-cache,.gitkeep`，未跟踪路径数为 0。按 §5.1「自愈档」要求，主控本人违规须经**同轮独立审查层确认**方才生效，而本片 ⑦ 额度（1 + 1）已用尽、本轮无可用审查层 ⇒ 按**非自愈档**登记并升级用户裁决；该违规不影响任何交付物字节（捕获件为纯 ASCII 门禁摘要、从未入库）。**第二起（同类）**：复查门禁读数时再次将两个捕获件写到 `tmp/` 根（`tmp/gi.txt`、`tmp/gt.txt`）；同样当场自报并逐件删除，零残留读数：两路径 `Test-Path` 均 `False`，`tmp` 递归文件计数 103 → 101（排除 `gates-cache/g4b` 后为 1），顶层名回到 `gates-cache,.gitkeep`，未跟踪路径数为 0。两起合计 D = 2，已接近 §5.1「D 累计达 3 次一律停片升级」的上界；本片末段已改为只用内联脚本打印、不再落任何仓库内文件。

## 2. 开关设计

- 数据源：引擎所有状态事件经 `Engine.transition` → `stateWriter.append` → `FileEventStore.Append`，后者**每事件 `file.Sync()`** ⇒ 事件按序落盘。
- 形态 A：`progressStore` 装饰 `*chain.FileEventStore`；`Append` 先转发落盘、成功后打印；`Load` 纯转发。
- **SealSink 保真**：装饰器必须实现并转发 `SealTaskEvidence` / `LoadTaskSeal`（`internal/chain/engine.go` 与 `seal_verify.go` 共 3 处对 `Options.Events` 的类型断言）；不转发会静默关闭封存重放或使接受闭包 fail-closed。
- 输出分流：`--progress` 只写 **stderr**；stdout 保持 `writeRunSummary` / `writeJSONResponse` 原路径 ⇒ 默认与 `--json` 均逐字节不变。
- 行格式（固定）：`progress: <occurredAt> <entity> <fromState>-><toState> <reasonCode>`；`entity` 为 `chain` / `task <taskId>` / `step <taskId>/<stepId>`；空字段以 `-` 占位。
- 三处注入：`ExecuteNoopRun`、`executeAgentRunnerRun`、`supervisedOpenWindowAt` 的只读重放点**不装饰**。

## 3. 数据源与形态选择

- `events/state-events.jsonl`：运行中可读（每事件 Sync、非独占追加、每次 Close）。
- AgentRunner 的 `events.jsonl`：由外部 pinned CLI 写入、无 Sync、prfrail 仅结束跳读取 ⇒ 不可作为运行中进度源。
- 否决 B（轮询）：需 goroutine 生命周期、ctx 取消、poll 成本与末行顺序竞态。
- 否决 C（仅 AgentRunner 事件）：运行中不可达且看不到 task/step 起止。

## 4. 反例与变异检验

测试清单（R1–R9，一句话各一条）：

- R1 默认关闭时 stdout 逐字节不变且 stderr 空；R2 `--json` 与 `--json --progress` 的 stdout 逐字节相等；R3 坏 writer 不影响 Append 结果与事件字节；R4 装饰器不改事件字节；R5 行格式快照；R6a/R6b SealSink 能力与写路径转发；R7 usage 与 flag 接受；R8 AgentRunner 路由的投影/事件形状/行数不变；R9 Supervised 路由的 stdout 与 `--json` 报告不变、progress 仅进 stderr。

变异（M1–M8，逐条一句话，全部实跑变红且逐字还原、恢复后 SHA256 一致）：

- M1 nil 守卫；M2 写 stdout；M3 返回打印错误；M4 改事件字节；M5 占位符；M6 `SealTaskEvidence` 不转发；M7 关闭 JSON 归一化；M8 关闭注入。

## 5. 审查链结论

- ⑤ `PRE-REVIEW: PASS`（3 Low + 3 Info）；F1/F5 已整改，F2 由后续 R8/R9 覆盖，F3/F4/F6 为 Info 不整改。
- ⑥ `INDEPENDENT SCAN: PASS`（0 发现）。
- ⑦ 终审 `FINAL REVIEW: PASS WITH FIXES`（1 Low：两路缺运行期覆盖）→ 加 R8/R9 → 复审 `RE-REVIEW: FINDINGS`（Low 未闭合：缺显式 JSON 不变性）→ 补 JSON 不变性断言 → 机械自证。
- **⑦ 盲审判定**：本片不触所有权 / 停机 / 身份语义 ⇒ 必抽不成立；按比例抽查亦**不抽**——P67 的盲审计入 §7.7 窗口的「已抽片」（§12.4 清账登记的「不计 §7.7 必抽 +1」只表示不占必抽名额，不表示不属于 §7.7 窗口；P69 事实包已按「窗口 P67=1 / P68=2 / P69=3，未满 4 片不抽」执行同一口径）⇒ 窗口内 P68 / P69 / P70 为 3 片未抽，未满 4 片。
- **P68 是否计入窗口**：计入——它在台账「后续批次切片」节以 `✅ 已完成` 登记，非 `⏸ 待评估`（后者才是非正式切片标记）⇒ P71 起窗口为 P68 / P69 / P70 / P71 四片未抽，按 §7.7 应抽。
- ⑦ 复审报 1 Low（AgentRunner/Supervised 两路缺显式 JSON 不变性断言）；主控已补 R8/R9 断言并经 `go test ./internal/console/ -run Progress` 全绿自证；⑦ 额度 1+1 已用尽，未开第 3 轮确认 ⇒ 按 §7.8.3「额度用尽」处理，待办：并入下次 ⑦ 输入（P71 起跑时须引用本行）。

## 6. 门禁读数

本片 ④ 批次原始读数（Windows 本机，2026-11-13；⑧ 文档装配后由 ⑨ 复跑给出定稿读数）：

- `gofmt -l .` 空输出；`go build ./...` exit 0；`go vet ./...` exit 0；`go test ./... -count=1` ok 包 15 / FAIL 0。
- `node tools/gates/gate.js --selftest` ⇒ `SELFTEST: PASS 737/737`。
- `node tools/gates/gate.js --all --scope=tree` ⇒ `changed=19`、`TOTAL_FAIL=0`；`--scope=ci` ⇒ `changed=2`、`TOTAL_FAIL=0`。
- `GOOS=linux GOARCH=amd64 go build ./...` 与 `go vet ./...` 均空输出、exit 0（Unix 零影响的本机证据；CI ubuntu 腿兜底）。
- 聚焦测试 `go test ./internal/console/... -count=1 -run Progress -v`：九条 `--- PASS` 全绿。
- ⑨ 定稿读数（⑧ 文档装配完成后复跑，先精确暂存再取 index）：`--scope=index` ⇒ `changed=33`、`TOTAL_FAIL=0`；`--scope=tree` ⇒ `changed=33`、`TOTAL_FAIL=0`；`--scope=ci` ⇒ `changed=2`、`TOTAL_FAIL=0`；`SELFTEST: PASS 737/737`；`go test ./... -count=1` exit 0；`GOOS=linux` 的 `go build` / `go vet` 均 exit 0。
- CI 链路（提交后）：主提交 `f77ab49` 首跑 **ubuntu 腿 `Test` 红**（run `43400311736`，exit 1；windows 腿绿）⇒ 同 run `gh run rerun --failed` 的 attempt 2 **ubuntu 仍红、windows 仍绿、同一测试同一行号同一断言串**（非 flake，登记 DR-14）⇒ 修复提交 `6a405bd` 的 run `43403891205` **双腿 success**（ubuntu 3m17s 含 Race / Contract fixtures / Gates、windows 4m28s；两腿包级 `ok github.com/larsonzh/prfrail/internal/console` 分别为 `3.477s` / `15.783s` 与 `24.133s`，`--- FAIL` 行数均为 0）。

## 7. 已知边界

- 门禁命令的进程退出码：`StateEvent` 字段闭集无 exitCode，`reason.code` 为语义码 ⇒ 不可达，补法需改事件模型（越界）。
- AI 工具调用摘要行：外部子 CLI 批写且无 Sync、prfrail 仅末跳读取 ⇒ 不可达。
- port 内命令型门禁（baseline / acceptance / review / scan）的起止：无对应事件 ⇒ 仅 hook/gate 形态的步骤可覆盖。
- AgentRunner 的 **CLI 成功腿**未在 CLI 层断言：Windows 平台 `replayStorePublicationDurability()` 恒 `Unproven` ⇒ 生产 dispatcher fail-closed；R8 改为覆盖 `run` 命令实际调用的同一生产入口。
- ⑦ 复审 Low 的整改未获第 3 轮独立复核（额度用尽，§7.8.3），并入下次 ⑦ 的输入。

## 8. 明确未执行事项

- 未改 `internal/chain/**`、`internal/adapters/**`、`internal/gates/**`、`internal/tickets/**`、`internal/guard/**`、`schemas/**`、`testdata/contracts/**`、`.github/workflows/**`、`go.mod`。
- 未改 `docs/DELIVERY_DIRECTIVE{,_EN}.md`（准则正文零改动）、未改 `docs/CONTRACTS{,_EN}.md`。
- 未改 v50–v56 历史快照；未 commit / 未 push；未推 gitee。
- D=2 的自愈确认者缺失：因 ⑦ 额度 1+1 已用尽，无法同轮独立审查层确认；按 §5.1 非自愈档登记并升级，用户 2026-11-13 裁定接受，限定条件 = 若第 3 次同类违规（D=3）发生，按 §5.1 累计上界强制停片。
- 防再犯：⑧c 后取证窗口内的门禁捕获件一律使用内联打印（`node -e` / `execFileSync`）或落 `tmp/<切片ID>/`；不得落 `tmp/` 根。本片 D=2 两起均为该违规的实例。
- P68 计入 §7.7 窗口（判定见 §5）⇒ **P71 起跑应抽盲审**；P71 起跑时须引用本行。
- DR-14（CI 首跑红，本片修复）：载体 = `internal/console/progress_test.go:172` 起的 `progressSummaryJSON` 守卫；台账登记见 `docs/t027/REMAINING_SLICES{,_EN}.md` 的「已知缺陷」表 DR-14 条目；证据 = run `43400311736`（head `f77ab49`）ubuntu 腿 `Test` 步红、attempt 2 仍红、windows 腿绿；修复提交 `6a405bd`（run `43403891205` 双腿 success）。
- **P71 落地指针（2026-11-13）**：本片 §8 记录的「P71 起跑应抽盲审」已执行——P71 `S1-SCOPE-ANNOTATE` 已抽 1 次盲审（先于终审），结论见 `docs/validation/s1-scope-annotate{,_EN}.md` §5；本片原行不改（§11.1）。

以下为本片归档条目清单（仓库相对路径 + 处置，TAB 分隔；CN 与 `_EN` 两侧围栏逐字节相同）。

```artifacts
docs/OPERATIONS.md	repo
docs/OPERATIONS_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/s1-stream-progress.md	repo
docs/validation/s1-stream-progress_EN.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/00-briefing.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/01-arch/design-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/02-batch/batch-01-state-events-enumeration.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/02-batch/batch-02-agentrunner-events-enumeration.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/scan-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/final-r1.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/final-r1b.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/06-gates/gate-readings.txt	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/06-gates/mutation-evidence.txt	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/MANIFEST.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/SHA256SUMS.txt	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/freedom-list.md	repo
docs/validation/evidence/P70-S1-STREAM-PROGRESS/input-list.md	repo
internal/console/cli.go	repo
internal/console/progress.go	repo
internal/console/progress_test.go	repo
internal/console/runtime.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_supervised.go	repo
internal/console/runtime_zero_model_e2e_test.go	repo
internal/console/stopper_test.go	repo
internal/console/supervised_cmd.go	repo
```

以下为 ⑥ 独立扫描 r1 的 Section A–E 逐字原文（转录自 `docs/validation/evidence/P70-S1-STREAM-PROGRESS/04-reviews/scan-r1.md`）。

```sixstage
## Section A

INDEPENDENT SCAN: PASS

## Section B

(no findings)

## Section C

- 透明性判断：在 [internal/console/progress.go](internal/console/progress.go) 里，装饰器先转发到内层事件存储，只有在 `inner.Append(...)` 成功后才写一行 progress；它没有重编码、没有改写 event，且不吞掉真实事件错误。真正会返回的错误仍然来自 `inner.Append`，这保留了原始失败语义。没有新 goroutine、没有后台 poll、没有额外状态对象，也没有泄漏路径；它是纯观察器，不改变事件日志的持久化字节。
- SealSink 保护：装饰器显式声明 `var _ chain.SealSink`，并且两个方法都转发到内层存储；因此 [internal/chain/engine.go](internal/chain/engine.go) 中的三处 `Events.(SealSink)` 断言仍然成立，且 [internal/chain/seal_verify.go](internal/chain/seal_verify.go) 的恢复校验也仍可读到 seal。若装饰器未转发，R6 相关断言会立刻红；当前代码没有这个缺口。
- `--progress` / `--json`：在 [internal/console/cli.go](internal/console/cli.go) 中，`progressWriter` 只在 `--progress` 真值时绑定到 `stderr`，默认关闭时保持 `nil`，因此无 progress 输出，stdout 字节流保持原样；`--json` 仍走 stdout，和原有输出路径一致，不存在"把 progress 丢进 stdout"的回归。
- 改动面界定：输入包中的原始 diff 只涉及 docs 台账与 [internal/console](internal/console) 目录；它没有越过到 [internal/chain](internal/chain) 等生产代码。也没有修改运行时协议、事件模型、CONTRACTS 语义或外部 adapters 层。
- 测试审计：新增的 [internal/console/progress_test.go](internal/console/progress_test.go) 设计了 R1–R7，对默认关闭、stderr-only、writer 失败、事件字节不变、占位格式、SealSink 转发、usage 接受等关键面都做了断言；它们不是恒真语句，而且用"删守卫/改写输出/改写 SealSink"来制造显式红灯，未发现过度掩码或假绿。

## Section D

- 删除 `if progress != nil` 这个守卫：R1 会红，且会出现 nil writer panic；默认关闭路径会被破坏。
- 把 progress 改写到 stdout：R2 会红；`--json` 与非 JSON 的 stdout 对比直接不一致。
- 去掉 "打印失败忽略" 的处理：R3 会红；broken writer 直接暴露错误。
- 在 `Append` 前 mutate 事件：R4 会红；事件文件字节发生变化。
- 把 placeholder `"-"` 改为空：R5 会红；固定快照断言失败。
- 删除 `SealTaskEvidence` 转发：R6b 会红；seal 文件不会落盘且 hash 不一致。
- 删除 `--progress` 用法接受：R7 会红；usage 断言直接失败。
- 反之，当前代码在这些守卫下都能保持对齐；没有看到"删掉某守卫后仍旧绿"的情况。

## Section E

- 我无法验证的部分：本次结论基于静态读取、设计基线和新增测试的断言意图，不是基于真实进程执行输出的现场复现；因此对特定 live 环境下的并发/终端交互仍保留运行时观察限制。
- 什么会推翻我的 PASS：如果未来提交把 progress 改回 stdout、移除 SealSink 转发、修改 event 记录前后字节，或者在 [internal/console](internal/console) 之外引入生产代码修改；若测试被改成掩码式"只看 stdout 不看磁盘"，或有真实并发下的 race/泄漏出现，也会推翻本结论。按现有读取材料，没有发现这些信号。

INDEPENDENT SCAN: PASS
```

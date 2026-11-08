# ③ 层 2 核心场景复跑回执（test-r2）

- 落盘：2026-11-08；仓库 HEAD = `05e8e35`；工作目录 `D:\LZProjects\prfrail`。
- 范围：**层 2 核心场景**（F1 隔离 workspace、F4 未知恢复）。**零计费**（纯 harness 重跑，无 BYOK、无 Copilot 调用）。
- 纪律：未 commit / 未 push；**未修改 harness**（`tools/b4-e2e/**` 的 `git status` 无改动）；未改任何生产/工具脚本；证据内无密钥。
- 口径（用户裁定 + design-r1 §D.4）：机制腿 = **B4 harness 重跑**（`tools/b4-e2e/run-faces.ps1`，复用生产包、无 `os/exec`），**非等价替代、与 B4 同口径**；「必通」= 本片 harness verdict 类与 B4 基线**逐项一致**（基线取归档包 `docs/validation/evidence/directive-history/b4/archive/pack-verdict.json`）。

## 1. 执行命令摘要

| 步 | 命令（原文） | 退出码 |
|---|---|---|
| 预播种 1 | `New-Item -ItemType Directory -Force -Path tmp\b4` | 0 |
| 预播种 2 | `Copy-Item docs\...\b4\archive\input\b4-chain.json tmp\b4\input\ -Force`；同 `b4-chain-hook.json` | 0 |
| F1 | `powershell -NoProfile -ExecutionPolicy Bypass -File tools\b4-e2e\run-faces.ps1 -OnlyFace f1 -SkipBuild` | **0** |
| F4 | `powershell -NoProfile -ExecutionPolicy Bypass -File tools\b4-e2e\run-faces.ps1 -OnlyFace f4 -SkipBuild` | **0** |

预播种输入与 B4 基线摘要逐字相等：`b4-chain.json` = `bb8c100d…39361a`，`b4-chain-hook.json` = `c7ad22c5…795237`（与 `t027-at23-e2e.md` §2 表及 `b4/MANIFEST.md` 一致）。run root 的 `workspace` 由驱动器按自身口径 `Remove-Item` + `Copy-Item` 播种（seed = `tmp/b4/input/seed/b4-seed.txt`），未手工干预。

环境坑（记为既有 B4 教训，未改脚本）：首次直跑 `exit=1`，`DirectoryNotFoundException: tmp\b4\driver.log`——驱动器在 `Write-Log` 初始化处先写 `tmp/b4/driver.log`，再建 `bin/input/runs/evidence/seed`，故 **`tmp/b4` 必须预先存在**；补建后第二次 `exit=1`，`emit-record failed: load chain: open tmp\b4\input\b4-chain.json`（脚本不自建链定义）。两者均属环境预播种缺失，非 harness 缺陷，按「照做勿改」补齐后即绿。

## 2. F1 隔离 workspace（核心）——断言对照

| 腿 | verdict 类 | 链 / 任务 / 步骤 | 机械断言（实测 detail） | 与基线 |
|---|---|---|---|---|
| `f1/pos` 正例 | `completed` | `COMPLETED` / `b4-task1=PASSED` / `b4-code-step1=PASSED` | 冻结事实齐备且互异（`distinct fact digests: 6`）；`workspace-unchanged`：`diff added=[] removed=[] changed=[]` | 一致 |
| `f1/neg2` 反例（终局后注入写入） | `uncertain-or-paused` | `PAUSED` / `b4-task1=REPAIR_PENDING` / 步骤 `PASSED` | `frozen-diff-excludes-injection`：冻结 diff 三项全空；`postflight-rejected-unaccounted-write`=REPAIR_PENDING；`never-passed`：`task reached PASSED: false`；`task-final-transition-is-postflight-rejection`：`sequence=11 toState=REPAIR_PENDING reason=postflight-rejected` | 一致 |
| `f1/neg2-r2` 包内重跑 | `uncertain-or-paused` | 同上 | 同上（同 `sequence=11`、同 reason） | 一致（与 `f1/neg2` 同 verdict 类、同末次转移） |
| `f1/neg` 控制腿 | `completed` | `COMPLETED` / `PASSED` / `PASSED` | `mutation-captured-in-diff`：`diff added=["agent-stub-mutation.txt"]`；`postflight-rejection(design expectation)`=`false` 且 kind=`deviation`（B4 偏差 2） | 一致 |
| `f1/neg3` 边界腿 | `completed` | `COMPLETED` / `PASSED` / `PASSED` | `injection-present-on-disk`=true；冻结 diff 空；`unaccounted-change-rejected(design expectation)`=`false` 且 kind=`deviation`（B4 偏差 3，捕捉排除集） | 一致 |

结论：正例断言（`completed` + 事实互异 + 工作树 diff 空）与反例断言（`uncertain-or-paused` + `REPAIR_PENDING` + `postflight-rejected` + 冻结 diff 空）**全部成立**，且反例在独立 run root 上二次复现。

## 3. F4 未知恢复（核心）——断言对照

| 腿 | verdict 类 | 机械断言（实测 detail） | 与基线 |
|---|---|---|---|
| `f4/pos` p1 park | `parked`（`RUNNING` / `STEPS_RUNNING` / `TERMINAL_PENDING`） | 派发腿：经生产 launcher 启动并停在 `TERMINAL_PENDING` | 一致 |
| `f4/pos` p2 recover | `uncertain-or-paused`（链 `PAUSED`） | `dispatchEvents`=`1`；`identityMirrorStable`=`true`；`recover-refused`（`engine.Run after Recover returned chain recovery uncertain`）；`single-dispatch`（`agent-dispatched TERMINAL_PENDING events for the step: 1`）；`no-new-process-identity`（`mirror unchanged across the restart: true`）；`chain-paused` | 一致 |
| `f4/neg` p2 forge | `refused` | `resume-continuity-refused`：`AgentRunner resume continuity unproven: prior completion sha256:032c24e8…07326e is not in step b4-code-step evidence history`；`zero-write`：`event log digest before=sha256:745165bf…a5f80 after=sha256:745165bf…a5f80`（前后逐字相等） | 一致（拒绝原因串与基线逐字相同） |
| `f4/neg-r2` p2 forge | `refused` | 同上（`progress` 独立 run root，digest 对为 `6ae45440…eb24b` 前后相等） | 一致 |

身份镜像逐字证据（`tmp/b4/runs/f4-pos/managed-process.identity.json`，单一文件、单一 launch slot）：`{"pid":5568,"startToken":"01dd6f1a7e4467ed"}`；`managed-process.launch.slot` = `launch-request-run-b4`；`events/state-events.jsonl` 共 10 行，其中派发/`TERMINAL_PENDING` 相关 1 行 ⇒ **重启后零 relaunch**、`pid`/`startToken` 仅一组。伪 resume 零写另有磁盘旁证：`f4-neg` 下 `terminal/`、`events/terminal.jsonl` 均 `exists=false`。

## 4. 与 B4 基线逐项对照（机械，`l2-baseline-compare.txt`）

- 覆盖 9 条腿（F1 六条 + F4 三条）：`verdictClass` / `chainState` / `taskState` / `stepState` / `observations` **逐项一致，0 处不符**；两次运行的 `failedChecks` 均为 `[]`；`deviations` 集合与基线同（`f1-neg`、`f1-neg-r2`、`f1-neg3` 各 1 项，均 kind=`deviation`）。
- **唯一标签漂移（如实登记，非失败）**：F1 面断言 id 由基线的 `five-facts-distinct`（detail：`distinct fact digests: 5`）变为当前 `six-facts-distinct`（`distinct fact digests: 6`）。成因可定位到 B4 之后的产品改动 `02d5797`（2026-10-22，`feat: bind pid reports into A5 frozen facts …`）：`factHashes` 现返回 6 枚（Manifest/Diff/Log/Usage/ProcessStopEvidence/ProcessPidReports）。语义为**严格超集**（对同一集合要求 6 枚互异），实测全绿；不构成 verdict 类不一致，交 ⑤/⑦ 知情。

## 5. 覆盖缺口登记（不补写、不改 harness）

| # | 缺口 | 性质 |
|---|---|---|
| G1 | `f4/pos` 无包内重跑腿（结构上只有 `park`→`recover` 一个 run root）；F4 面的可重复性仅由 `f4/neg-r2` 承担 | 既有覆盖缺口（B4 结构如此，照抄不补） |
| G2 | `f1/neg3` 单跑（成因为编译期排除集而非竞态）；F1 面另有 `f1-neg-r2`、`f1-neg2-r2` 两组重跑腿 | 既有（B4 偏差 3） |
| G3 | 生产派发 / 终局发布 / 门禁端口对**未重跑**（Windows `unproven` 下 fail-closed），引用 B4 附录 A 为基线 | 本片口径明示 |
| G4 | 机制腿为 **非等价替代**（harness 复用生产包），不主张等价覆盖 | 用户裁定 |
| G5 | 层 2 未覆盖 F2/F3/F5；层 3 未跑 | 本单范围 |

## 6. 证据清单（`03-test/layer2/`，8 个顶层文件 + 2 个证据目录，共 44 文件）

| 文件 | sha256 |
|---|---|
| `l2-f1-driver.txt` | `b0c47905e68348fdc6738754cb8935d5520702ddec085709192c99a83d52acf4` |
| `l2-f1-pack-verdict.json` | `c227a0dd08cdd4b3f66d0753f39fd4dfa9a9c55bf08c52873bd22968c4633150` |
| `l2-f4-driver.txt` | `4b78d330d6b041fa03050b6511dcf14027d518e40bb84a0ed669cd6d71428109` |
| `l2-f4-pack-verdict.json` | `c79828ab8111c93374ad74e4b749a0888d5cb8e7b7a27c79b1808dedef12442e` |
| `l2-f4-assertions.txt` | `19291b459eb3669153769e2f078b722d0ad04dc9da27d00086108f81f7061f7f` |
| `l2-baseline-compare.txt` | `1db4400a6fdaaf42786b57bdb81ba5382b87bcddb71afe973a9ed8b6317961c0` |

另含 `l2-f1-SHA256SUMS.txt`、`l2-f4-SHA256SUMS.txt`（驱动器按 run root 逐文件复算）与 `l2-f1-evidence/`、`l2-f4-evidence/` 下各腿的 `verdict-*.json` 与 `.stdout/.stderr` 原文。编码已核：全部 no-BOM + LF（0 违规）；密钥扫描：无（`sk-` 命中项均为 `task-final-…`/`postflight-…` 词内子串）。

## 7. 未做事项

未跑层 3；未 commit / 未 push；零计费（无 BYOK、无 Copilot premium）；未改 harness 与任何生产/工具脚本；未注册 MANIFEST / SHA256SUMS（由主控负责）；`tmp/b4`（run root 与 pack）与本件一并留存待清理。

TEST: DONE (layer2 core green)

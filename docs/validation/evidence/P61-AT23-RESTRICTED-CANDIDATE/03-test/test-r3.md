# ③ 层 3 全场景回执（test-r3）

- 落盘：2026-11-08；仓库 HEAD = `05e8e35`；工作目录 `D:\LZProjects\prfrail`。
- 范围：**层 3 全场景**（机制腿 F2/F3/F4 重跑 + postflight 六事实重建 + 真实候选腿 3 笔）。
- 纪律：未 commit / 未 push；**未修改 harness 与任何生产/工具脚本**（`git status -- tools internal cmd schemas` 仅有 ② 已暂存的
  `tools/agent-probe/CopilotCli.ps1`，实测 sha256 = `F1899250…826BCA`，与 ② 记录逐字相同，`git diff -- tools` 为空）；
  临时件全部落在 `tmp/P61-AT23-RESTRICTED-CANDIDATE/`；证据内无密钥（扫描 0 命中）。
- 口径：机制腿 = **B4 harness 重跑**（`tools/b4-e2e/run-faces.ps1`，`-SkipBuild` 复用已固定二进制），
  **非等价替代、与 B4 同口径**；「必通」= verdict 类与归档 B4 基线逐项一致。

## 1. 任务 A：机制腿三场景（零计费）

| 步 | 命令（原文） | 退出码 | 证据 |
|---|---|---|---|
| 预播种 | `tmp/b4`、`tmp/b4/input/{b4-chain.json,b4-chain-hook.json}` 沿用层 2 播种，摘要复核 `BB8C100D…` / `C7AD22C5…` | 0 | 与 `t027-at23-e2e.md` §2 一致 |
| F2 | `powershell -NoProfile -ExecutionPolicy Bypass -File tools\b4-e2e\run-faces.ps1 -OnlyFace f2 -SkipBuild` | **1** | `layer3/l3-f2-driver.txt`、`l3-f2-pack-verdict.json`、`l3-f2-evidence/` |
| F3 | 同上 `-OnlyFace f3` | **0** | `layer3/l3-f3-*` |
| F4 | 同上 `-OnlyFace f4` | **0** | `layer3/l3-f4-*` |

环境坑（新登记）：从 `cmd /c` 启动 Windows PowerShell 5.1 时 `PSModulePath` 继承 PowerShell 7 模块目录，`Get-FileHash` 解析失败
（`CommandNotFoundException`），驱动器首跑 `exit=1`。改为由当前 shell 直接调用 `powershell -File` 后正常；**未改脚本**。

| 腿 | verdict 类 | 链/任务/步骤 | 机械断言（实测） |
|---|---|---|---|
| `f2/pos`（hang，5s 看门狗） | `uncertain-or-paused` | `PAUSED` / `FAILED` / `FAILED` | 超时触发 `timedOut=true elapsedMs=5020 boundMs=5000`；无冻结事实；错误证据 **6** 枚；**`process-tree-taken-away` = FAIL**：「reported pids gone: false (… missing required pid report(s) stub-child.pid …)」 |
| `f2/neg`（spawn 1h 子进程，3s 超时） | `uncertain-or-paused` | `PAUSED` / `FAILED` / `FAILED` | 无冻结事实；`process-tree-taken-away` PASS（`gone: true`） |
| `f2/neg-r2`（包内重跑） | `uncertain-or-paused` | 同上 | 同上（与 `f2-neg` 同 verdict 类、同断言结果） |
| `f3/pos`（`-no-marker`） | `uncertain-or-paused` | `PAUSED` / `FAILED` / `FAILED` | `status=uncertain logsComplete=false`；`factCount=0`；错误证据 3 枚 |
| `f3/neg` | `completed` | `COMPLETED` / `PASSED` / `PASSED` | `logsComplete=true`、`usageComplete=true`、`manifestsComplete=true`、`factCount=6`、错误证据 0 |
| `f3/neg-r2` | `completed` | 同上 | 同上（重跑一致） |
| `f4/pos`（park→recover） | `uncertain-or-paused` | `PAUSED` / `STEPS_RUNNING` / `TERMINAL_PENDING` | `dispatchEvents=1`；`identityMirrorStable=true`；`recover-refused`（`chain recovery uncertain`）；`single-dispatch`；身份镜像单组 `{"pid":46768,"startToken":"01dd6f1d56208b91"}`；`state-events.jsonl` 共 10 行；重启后零 relaunch |
| `f4/neg`（伪 resume） | `refused` | `RUNNING` / `STEPS_RUNNING` / `TERMINAL_PENDING` | `resume-continuity-refused`：`prior completion sha256:032c24e8…07326e is not in step b4-code-step evidence history`；`zero-write`：事件日志摘要 `b3773920…06e2` 前后逐字相等 |
| `f4/neg-r2` | `refused` | 同上 | 同上（拒绝原因串逐字相同）；`zero-write` 摘要对 `d5f9ba95…0dbc` 相等 |

**与派单措辞的两处偏离（如实登记，以归档基线为准）**

1. 派单写「F2 反例（对照腿）`completed`」，但归档 B4 基线中 `f2/neg`、`f2/neg-r2` **同为 `uncertain-or-paused`**——该腿是
   `-mode spawn -child-life 1h -sleep 5s` 配 3s 超时，必然超时。本轮实测与基线一致；**若按派单字面判 `completed` 则会把红判成绿**，
   故按 design-r1 §D.4「与 B4 基线逐项一致」判定，并把措辞差异交主控。
2. 派单写 F2 正例「错误证据 5 枚」，本轮实测 **6 枚**（基线 5 枚）。多出的是第六枚冻结事实引入的 pid 报告缺口摘要
   （`agentRunnerPidReportsGapDomain`），成因同下条。

## 2. 与 B4 基线逐项对照（`layer3/l3-baseline-compare.txt`）

- 覆盖 9 条腿：`verdictClass` / `chainState` / `taskState` / `stepState` / `observations` **字段漂移 0 处**。
- 断言 `pass` 值差异 **1 处**：`f2-pos/process-tree-taken-away`，基线 `true` → 本轮 `false`。
- 根因（可定位）：`02d5797`（2026-10-22「bind pid reports into A5 frozen facts」）使 pid 报告成为 fail-closed 缺口；
  而 `tools/agent-stub` 的 `hang` 模式在 `switch` 内 `return`，**不经过** `writeEmptyChildPIDReport()`，因此 `stub-child.pid`
  不存在 ⇒ 收集器判缺口 ⇒ 选择臂 `recordedPIDsGone` 返回错误 ⇒ 断言红。基线（2026-09-22）尚无该事实，故为 `true`。
- 该红是**真实、可复现、且方向为 fail-closed** 的观察：它不推翻 F2 的超时判定（`uncertain-or-paused` 仍成立），但推翻了
  「所报 pid 已消失」这一断言的**可证性**（进程证明降级为「不可证明」）。本片**不改** stub / 收集器 / harness，交 ⑤/⑦ 裁决。

## 3. 任务 B：postflight 六事实重建（零计费）

- 口径：按 A5 位置化证据（终局路由事件 `evidence[2..7]`）取冻结值，按 `internal/adapters/agent_runner_evidence.go` 的
  摘要规则从磁盘工件逐字节重算，逐项比对。
- 结果（`postflight-rebuild/l3-six-fact-rebuild.txt`）：
  - **2 条有冻结事实的腿**（`f3-neg`、`f3-neg-r2`）：manifest / diff / logs / usage / **pid-reports** 五项**全部 MATCH**（逐字节）。
  - `process-stop` 事实：**不可重建**——`guard.TerminationEvidence` 不落盘，只有摘要进入链。已显式标注，并以
    「所报 pid 存活检查」作语义旁证：全部 `gone`（`f3-neg` 32624、`f3-neg-r2` 18292、`f2-pos` 44844 等）。
  - 7 条**无冻结事实**的腿（F2 三腿、F3 正例、F4 三阶段）：确认终局路由事件缺席（无 `agent-terminal-completed`），
    即「无冻结事实」与 verdict 类自洽；同时记录工件缺口原样：`f2-pos` 日志空且 `logsComplete=false`、`usage.json` 缺失、
    `stub-child.pid` 缺失；`f4-pos` manifest/diff/usage 缺失（`crash` 模式不产出）。

## 4. 任务 C：真实候选腿（BYOK，批准 3 笔，实际计费 2 笔）

| 腿 | 形态 | argv（生产封装 buildCopilotArgs，无 `-Unlock`） | exit | 失败阶段 | 结果 |
|---|---|---|---|---|---|
| run-1 启动 | `-InvokeWrapper -WrapperPrompt "List the files in the current directory" -WrapperEffort high` | `--no-auto-update -p "List the files in the current directory" --effort=high` | **0** | 无 | 正常回复并实际执行一次 shell（`Get-ChildItem -Force`），六项 BYOK env 全部清零 |
| run-2 执行 | 既有工具 `Invoke-AgentProbe.ps1 -Scenario tool-deny -Run`（BYOK env 由配套件注入） | 计划 argv：`-p "<tool-deny prompt>" -C <workspace> --name proofrail-t026-tool-deny-<stamp> --available-tools=powershell --allow-tool=shell(Get-Location) --deny-tool=shell(Get-ChildItem) …`（**未执行**） | **1** | **start** | 三次尝试均 `exit=1`、**零模型调用、零费用** |
| run-3 停机 | 派单兜底：后台起封装 BYOK 调用 → `taskkill /PID <host> /T /F` | `--no-auto-update -p "Run this PowerShell command and wait for it to finish: Start-Sleep -Seconds 120. …" --effort=high` | taskkill **0**；宿主退出码**不可得** | **stop** | 40s 后杀整棵树；`taskkill` 回显终止 16172（子）与 43976（宿主）；杀后进程 `gone` |

- run-2 三次尝试明细（`real-candidate/run-2/command.txt`）：① `powershell -File` 无法把 `string[]` 从命令行绑定，`Scenario` 收到
  `-Provider`；② 改为进程内数组 splat 仍是**按位置**绑定，同样失败；③ 给 tmp 配套件加**命名哈希 splat**（仅 tmp 件改动）后，探针进入
  自身预检，被**版本钉死**拒绝：裸 `--version` 经 SEA loader 解析为 `1.0.85`，而探针 pin 期望 `1.0.83`
  （封装以 `--no-auto-update` 维持 1.0.83）。**未达到 T026 的 toolControl 测量点**，属「启动阶段失败」，不产生候选执行语义。
- run-3 说明：模型侧多次尝试 `Start-Sleep` 等命令均被判 `Permission denied and could not request permission from user`
  （headless 下无 `--allow-tool`），故进程处于活跃会话时被杀；**未打印 token 块**（进程树被杀，CLI 未走到汇总），
  `stderr` 仅含 1.0 fallback 告警。

**CLI token 计量摘录（run-1，原文 `real-candidate/run-1/stderr.txt`）**

```
! No token count multiplier is known for Anthropic model "deepseek-flash". Token estimates will use
  the 1.0 fallback; configure a recognized canonical model ID to enable model-specific token counting.
Changes    +0 -0
Duration   9s
Tokens     〔↑〕36.5k（24.8k cached）〔箭头字形两组被控制台码页破坏〕189
Resume     copilot --resume=d5e88f22-34cd-4590-a5f7-60347ed7d7f9
```

- **计量性质**：CLI 自报的 token 为 **1.0 fallback 估算**，**不是计费依据**；与 DeepSeek 侧实测费用对账时以服务端账单为准。
- **编码瑕疵（如实登记）**：run-1 的 stderr 首次捕获时控制台码页（CP936）把 `↑`/`→`/`↓` 字形破坏为 `鈫?`/`鈥?`，末字节被 `?` 吞掉；
  原文另存 `stderr.raw-cp936.txt`，`stderr.txt` 仅在**首部加 7 行说明**、正文逐字不动。数字 36.5k / 24.8k / 189 精确可用。

## 5. 预算与计量

- 批准：BYOK **恰好 3 笔**（第 5~7 笔），实际 **2 笔产生模型调用**（run-1、run-3），run-2 **0 笔**（预检前失败）。
- 未使用的 1 笔**未追加**（遵守「不重试、不追加」）；是否补做由主控/用户裁定。
- 无 Copilot premium 消耗；无 `--effort` 换档（固定 `high`）。

## 6. 条件闭合前提自查（④ 前）

| # | 前提 | 自查 |
|---|---|---|
| ① | 机制腿全绿（verdict 类与 B4 基线逐项一致） | **未达成**：verdict 类 9/9 一致，但 `f2-pos/process-tree-taken-away` 1 条断言由绿转红（成因 `02d5797`，见 §2） |
| ② | 真实候选腿 ≥3 次且证据具体 | **未达成**：2 笔真实调用（证据具体）+ 1 笔启动阶段阻塞（零费用、有 stderr/exit/阶段） |
| ③ | 报告三句话逐字出现 | 属 ④/⑤ 阶段产物，③ 无法自查（本片报告尚未撰写） |
| ④ | BYOK 结论未替代 B1 基线 | **成立**：run-1/run-3 不主张候选兼容、不替代 B1/B2 冻结结论；run-2 无测量结果 |
| ⑤ | `--effort`/`--context` 边界已登记 | **部分**：`--effort=high` 层 1 已判支持并固定；`--context`（`-C`）在 run-1/run-3 未使用（探针计划使用 `-C` 但未执行），登记为未覆盖边界 |

## 7. 覆盖缺口与未做事项

| # | 缺口 | 性质 |
|---|---|---|
| G1 | `f4/pos` 无包内重跑腿（结构如此），可重复性由 `f4/neg-r2` 承担 | 既有（同层 2） |
| G2 | F2 正例 1 条断言红 | 本轮新增观察（§2） |
| G3 | run-2 未取得 toolControl 测量（工具自身 pin 拒绝） | 本轮新增阻塞 |
| G4 | run-3 无 token 块、宿主退出码不可得 | 停机形态固有 |
| G5 | 生产派发/终局发布/门禁端口对、Linux 原生验证、SessionBridge `visible` | 不在本单范围（同 B4 §7） |

## 8. 证据清单

| 路径（相对 `03-test/`） | 内容 |
|---|---|
| `layer3/l3-f{2,3,4}-driver.txt` | 驱动器原文输出（no-BOM+LF） |
| `layer3/l3-f{2,3,4}-pack-verdict.json` | 三面 pack（含 `failedChecks`/`invocations`） |
| `layer3/l3-f{2,3,4}-SHA256SUMS.txt` | 驱动器逐文件复算 |
| `layer3/l3-f*-evidence/<case>/` | 各腿 `verdict-*.json` 与 `.stdout/.stderr` 原文 |
| `layer3/l3-baseline-compare.txt` | 与归档基线的机械对照 |
| `postflight-rebuild/l3-six-fact-rebuild.txt` | 六事实重建与逐字节比对 |
| `real-candidate/run-1..3/` | argv/command、stdout、stderr（含原文/注解双份）、exit、失败阶段 |

- 编码核验：`layer3/`、`postflight-rebuild/`、`real-candidate/` 共 **60 文件**，BOM/CRLF 违规 **0**；密钥扫描命中 **0**。

## 9. 结论

机制腿方面：F3、F4 全绿且与 B4 基线逐字一致；F2 的 verdict 类一致，但正例出现 1 条 fail-closed 断言红（pid 报告缺口），
属真实可复现观察而非脚本噪声。真实候选腿方面：启动腿成功、停机腿取得具体证据，执行腿被既有工具的版本钉死阻塞（零费用），
**未取得 toolControl 测量**。因此**不满足** design-r1 §D 的两项闭合前提（①、②），本片不得进入「层 3 全绿」结论。明细如下：

- FINDING-1（阻断闭合前提①）：`f2-pos/process-tree-taken-away` 红——`hang` 模式不写 `stub-child.pid`，`02d5797` 将其升格为
  fail-closed 缺口；证据 `layer3/l3-f2-evidence/f2-pos/verdict-run.json`、`layer3/l3-baseline-compare.txt`。处置建议：改 stub 的
  `hang` 分支补写空子报告（需改 harness 侧 `tools/agent-stub`，本角色无授权）。
- FINDING-2（阻断闭合前提②）：`Invoke-AgentProbe.ps1` 在 1.0.83/1.0.85 共存的主机上**不可运行**（裸 `--version`→1.0.85 与 pin 冲突），
  证据 `real-candidate/run-2/attempt-3.*`。处置建议：把该工具收敛到生产封装 `CopilotCli.ps1`（属 ② 后续切片，需授权）或改由
  封装直接执行等价探针（会改变 T026 测量语义，须 ⑤ 判定）。
- FINDING-3（措辞级，非缺陷）：派单对 F2 反例的期望（`completed`）与归档基线（`uncertain-or-paused`）冲突；本轮以基线为准。
TEST: FINDINGS

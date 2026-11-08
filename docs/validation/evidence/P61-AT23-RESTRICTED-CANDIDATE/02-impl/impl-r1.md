# ② 实现回执（impl-r1）— P61 AT23-RESTRICTED-CANDIDATE

落盘：2026-11-08。角色：② 实现（`prfrail-implementer`，模型 `DeepSeek V4.1 Flash`）。
输入：`01-arch/design-r1.md`（① 定稿，§A 用户裁定修订、§E ② 清单）+ `docs/validation/r08-r09-ruling.md` §2 R-08.1。

边界声明（本片全程遵守）：**未改任何 Go 文件**、**未改契约/schema/fixtures**、
**未执行 `git commit` / `git push`**、**未执行任何 BYOK 调用**（BYOK 层 1 属 ③）、
**未修改本包内 `MANIFEST.md` / `SHA256SUMS.txt`**、未处理或打印任何密钥。

## 1. 文件清单

| 文件 | 动作 | 机械属性 | sha256 |
|---|---|---|---|
| `tools/agent-probe/CopilotCli.ps1` | **新建** | UTF-8 **BOM=1**、**CR=0**（LF）、9,655 字节 | `f189925094085414c0e4c42b5924139f6fbb7fd00462dd63cf5e1bd5bc826bca` |
| `tools/agent-probe/Invoke-AgentProbe.ps1` | **未修改** | — | 见 §3 阻塞说明 |
| `tools/agent-probe/Invoke-CandidateProbe.ps1` | **未修改** | — | 见 §3 阻塞说明 |
| `tools/agent-probe/appcontainer-b2/container-candidate.ps1` | **未修改** | — | 见 §3 已知例外 |

新文件当前为未跟踪状态（`?? tools/agent-probe/CopilotCli.ps1`），暂存/提交由主控决定。
证据包内 `02-impl/CopilotCli.ps1` 快照未复制（仓库路径即权威，避免重复 payload）；如需快照，由主控在登记时复制。

## 2. 函数契约摘要（`tools/agent-probe/CopilotCli.ps1`）

常量与证据落点：

- `CopilotPinnedVersion='1.0.83'`（:58）、`CopilotUnlockedAcceptedVersions=@('1.0.83','1.0.85')`（:59）、
  采样间隔 `CopilotProbeSampleEvery=5`（:60）、默认钉版路径（:61）。
- 证据目录 = `<repo>/tmp/P61-AT23-RESTRICTED-CANDIDATE/`（:65-73）：
  `wrapper-probe-count.txt`（只增不删）、`wrapper-version-evidence.txt`（追加式，单行含时间/模式/结论/版本首行）。

`buildCopilotArgs`（:146-170，纯函数、无副作用）：

- 入参 `-Subcommand/-Prompt/-Effort/-Context/-Unlock`；返回 `string[]`。
- 恒定前缀 `--no-auto-update` 仅在 `-Unlock` 缺省时加入（:157 附近）；`-Unlock` 显式传 ⇒ argv 不含该前缀（用户裁定 design-r1 §A.1）。
- 参数拼装：`Subcommand` 位次紧随前缀；`-p <prompt>`；`--effort=<v>`；`-C <ctx>`；空值一律省略。
- 返回 `$argv.ToArray()`（空结果不产出流水线对象，调用方用 `@(...)` 归一）。

`Invoke-CopilotCli`（:172-235，fail-closed 启动器）：

- 入参同上 + `-ExePath`（默认钉版 `%APPDATA%\npm\node_modules\@github\copilot\node_modules\@github\copilot-win32-x64\copilot.exe`）。
- 每次调用先推进探针计数（`Update-CopilotProbeCount` :79-93）。
- `-Unlock` 缺省：仅在 **首次（count==1）或每 5 次（count%5==0）** 抽样；抽样 = 实跑 `<exe> --no-auto-update --version`，
  要求退出码 0 且 stdout 文本含 `1.0.83`（`Read-CopilotVersionText` :120-137，不用 winget/ProductVersion）；
  失配 ⇒ 写证据（`Write-CopilotVersionEvidence` :95-118）后 `throw` 拒绝执行（终止性错误 ⇒ 脚本上下文非 0 退出）。
- `-Unlock` 显式：跳过钉住校验，改为记录实际版本串；提取首个 `\d+\.\d+\.\d+`（`Get-CopilotVersionToken` :139-144），
  必须 ∈ {1.0.83, 1.0.85}，否则写证据后 `throw`「版本不可控」。
- 执行：`& $resolvedExe @argv`（:230），stdout/stderr 由调用方重定向；返回 `[int]$LASTEXITCODE` 并写入
  `$script:CopilotCliLastExitCode`（:233）。**调用方注意**：子进程 stdout 与返回的退出码同属成功流，
  需要同时取文本与退出码时，用 `$script:CopilotCliLastExitCode` 取码、`Invoke-CopilotCli ... 2>&1 | Out-String` 取文本。
- 已知边界注释已写入文件头：版本字符串非稳定身份（B1 DR-2），本封装按 R-08.1 字面以 `--version` 文本为证据。

## 3. 任务二（收敛既有直接调用点）＝ **未完成，需澄清**

**结论**：6 处调用点全部**保持原样**，无法在不改变场景语义的前提下改走 `Invoke-CopilotCli`。
原因为**契约层面硬冲突**（不是实现偷懒），按本角色规矩「方案书模糊时不得猜测」停下并登记。

逐站点判定：

| # | 站点 | 现状形态 | 不能收敛的硬原因 |
|---|---|---|---|
| 1-4 | `appcontainer-b2/container-candidate.ps1` :94/:96/:231/:251 | 复制钉版到 `bin\copilot.exe`，在零能力 AppContainer 内以 `cmd /c copilot.exe …`（**按名**经 PATH 启动，绝对路径被拒） | 跨容器路径 + 按名解析 + argv 字符串内联重定向；封装在容器外执行，且 `-ExePath` 语义（容器内文件名）不成立 ⇒ 任务书已预授权为**已知例外** |
| 5 | `Invoke-AgentProbe.ps1` :540（另 :440/:443 为 `--version`/`--help` 身份采样） | `Start-Process -FilePath $ExePath -ArgumentList $argList -RedirectStandard{Output,Error,Input} -PassThru -NoNewWindow`，随后 `$proc.Handle`/`HasExited`/`ExitCode`/`taskkill /T /F` 超时清理 | 封装只有 `& exe @argv` 单一模式：无 `-PassThru` 句柄、无 stdin/stdout/stderr 文件重定向、无超时 kill；而上述字段本身是 T026 探针的**测量对象**（证据里的 pid/exitCode/超时判定），收敛即改场景逻辑 |
| 6 | `Invoke-CandidateProbe.ps1` :259（另 :157 为 `--version` 采样） | 同上 `Start-Process … -PassThru`，并轮询 `Win32_Process` 采集进程树证据（`osargv.txt`） | 同上；额外冲突：该脚本**允许自定义 `-Executable` + 自定义 `-ExpectedSha256`**，并明文记录「版本文本非稳定身份、1.0.83→1.0.85 同一字节」⇒ 封装的 fail-closed 文本钉版会把该脚本**记录态**的既定容忍条件变成拒绝态（语义回归） |

另一层原因（两处 Start-Process 站点共有）：封装入参不含 `--model` / `--output-format` / `--stream` /
`--max-ai-credits` / `--available-tools` / `--allow-tool` / `--deny-tool` / `--deny-url` / `--secret-env-vars` /
`--usage-output-file` / `--log-dir` / `--name` / `--resume` 等，收敛必然丢参数。设计 §E 已冻结封装参数集，
② 不得自行扩参（避免以「收敛」为名做架构决策）。

**需澄清（交 ⑤ 预审 + 主控裁定，② 不自行选择）**：

- 选项 A：维持 6 处原样，按 R-08.1「新增调用点必须走封装」的既有偏离流程，在 ⑤ 预审登记为偏离
  （理由：既有调用点非新增；且其执行语义与封装契约不兼容）。**改动量最小、无行为风险**。
- 选项 B：扩封装契约（例：新增 `-LauncherMode StartProcess` + `-ExtraArgumentList` 透传），再把两处探针改为走封装。
  该选项触及 §E 冻结参数集 ⇒ 需 ①/用户裁定，且仍需复核「进程树证据/超时 kill」的等价性，② 不预判。

## 4. 任务三（Copilot-hosted GPT-6 Luna 连通性，第 2 笔付费，1 个 Copilot premium）

- 隔离 workspace：`tmp/P61-AT23-RESTRICTED-CANDIDATE/luna-ws/`（新建，空目录）。
- argv：`<钉版 copilot.exe> --no-auto-update --model gpt-6-luna -p "say hello"`（经钉版 exe 直接调用，未走 `-Unlock`）。
- BYOK env：**未设置任何 BYOK 变量**；运行器仅做「移除」（六项 BYOK/模型变量 before=[] / after=[]，`DEEPSEEK_API_KEY` 亦未出现）。
- 结果：**exit 0，回复非空（首行 `Hello!`）⇒ `connected`**。未重试、未追加调用。
- 证据：`tmp/P61-AT23-RESTRICTED-CANDIDATE/luna-connectivity.txt`（含 argv/workingDirectory/exit/判定/首行 ≤120 字符；
  含 `formatRepairedAtUtc` 说明：首次落盘的换行被数组运算符优先级缺陷拆断，**仅重排格式、值未变、未发生第二次调用**）。
- 与设计差异（需主控登记）：design §E 写 `02-impl/luna-connectivity.txt`，本轮任务书指定 `tmp/…/luna-connectivity.txt`；
  本文件按任务书写入 tmp，如需入证据包由主控复制登记。Luna 结果不进验收结论（R-08.6）。

## 5. 机械校验（原始输出摘要）

- `Parser::ParseFile`（Windows PowerShell 5.1）：`PS51_PARSE_ERRORS=0`。
- 编码：`tools/agent-probe/CopilotCli.ps1` → `bom=1 cr=0`。
- `buildCopilotArgs` 行为：
  - `DEFAULT=--no-auto-update|version`
  - `UNLOCK=version`（不含前缀）
  - `FULL=--no-auto-update|-p|say hello|--effort=high|-C|C:\w`
  - `EMPTY_UNLOCK_COUNT=0`（空 argv 不产出流水线对象）
- `Invoke-CopilotCli` 运行时自检：
  - 钉版正常路径：`PINNED_LAST_EXIT=0`，`PINNED_LAST_VERSION=GitHub Copilot CLI 1.0.83.`
  - 版本失配（以 `where.exe` 冒充候选，第 5 次抽样点）：`THREW=Pinned Copilot CLI version check failed: expected 1.0.83, exit 1. Refusing to run.` + 证据行 `mode=pinned verdict=mismatch`
  - 可执行缺失：`MISSING_THREW=True`
  - `-Unlock` 不可控版本：`UNLOCK_THREW=Copilot CLI version is not controllable: reported "" is outside 1.0.83, 1.0.85. Refusing to run.`
- 全部 `--version` 自检为本地只读调用，**无模型调用、无 premium 消耗**；唯一付费调用 = §4 的 Luna 1 次。

## 6. 未做事项 / 遗留与风险

1. **未收敛 6 处既有调用点**（§3，需澄清）。
2. 未改 Go（`internal/adapters/ai_probe_copilot.go` 维持 R-08.1 字面的命名/签名偏离，见 design §A.2）。
3. 未 commit / push；未加入暂存区。
4. `tmp/.../wrapper-probe-count.txt` 因 ② 自检已推进至 **5**（只增不删规则下不回滚）。层 1 起按此基数继续抽样
   （下一次抽样点为 count=10）。如主控认为需从 1 起算，需在 ③ 开工前以书面裁定处理，② 不擅自删除。
5. 新增 `02-impl/impl-r1.md` 尚未登记进本包 `MANIFEST.md`/`SHA256SUMS.txt` ⇒ `gate.js` 的 G8-b（payload 未列入
   `## 归档条目清单`）预期变红，**由主控登记**后复算（与 `01-arch/design-r1.md` 当时同型）。
6. ② 自检产生的 `tmp/P61-AT23-RESTRICTED-CANDIDATE/*.ps1`、`*.js`、`*-out.txt` 属临时件，用完可清（`tmp/` 已 gitignore）。

IMPLEMENT: FINDINGS

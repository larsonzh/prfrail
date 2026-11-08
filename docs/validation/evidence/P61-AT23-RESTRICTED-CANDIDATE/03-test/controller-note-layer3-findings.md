# 层 3 FINDINGS 主控登记与升级注记（2026-11-08）

## FINDING-1（阻断前提①：机制腿 F2 一处断言漂移）
- 事实：F2 正例 `process-tree-taken-away` 由 B4 基线的 `true` 变 `false`；9 腿 `verdictClass/chainState/taskState/stepState/observations` 漂移 **0**，断言 pass 差异仅此 1 处；错误证据 5 → 6 枚。
- 成因：`02d5797`（2026-10-22，A5 pid reports 绑定为第六枚冻结事实）使 `tools/agent-stub` 的 `hang` 分支「不写 `stub-child.pid` 子报告」被新事实判为 fail-closed 缺口。
- 性质：**继承性漂移**（B4 基线若在 02d5797 之后重跑同样会红；未实跑，按代码语义推断）。不是 P61 引入的缺陷。
- **F2 verdictClass 确认（用户要求）**：`l3-baseline-compare.txt` 实测 `field drifts=0`；f2 三腿 `now=uncertain-or-paused/PAUSED` 与基线逐字一致，checks 数同（6/4/4）；**仅断言值** `process-tree-taken-away` true→false 漂移，不影响层 2/3 结论。
- **裁定（2026-11-08，用户）：选 (b)**——登记继承漂移、不修 harness（B4/C1 先例「验证片不得改被测 harness」），交 ⑦ 判条件闭合前提①的满足性。

## FINDING-2（阻断前提②：真实候选腿执行阶段零计费失败）
- 事实：run-2（执行阶段，探针 `-Scenario tool-deny`）三次尝试全部在 **start 阶段失败、零计费**：① `-File` 数组绑定失败；② 数组 splat 按位置绑定；③ 配套件改为命名 splat 后，`Invoke-AgentProbe.ps1` 自身 pin 校验拒绝（裸 `--version` = 1.0.85 ≠ 工具钉版 1.0.83；1.0.83/1.0.85 共存主机的既有行为）。每次均有 argv/stderr/exit 具体证据。
- 计费实况（**更正版，以 DeepSeek 实测为准**）：用户 2026-11-08 DeepSeek 网页实测层 3 消费 **¥0.04**（累计 ¥0.08 = 层 0 ¥0.01 + 层 1 ¥0.03 + 层 3 ¥0.04）。原「证实计费 1 笔」表述**撤回**——本地 token 块不足以支撑费用归因，如实更正如下：run-2 三次尝试零模型调用证据成立（前两次未启动 copilot；第三次仅 pin `--version` 后拒绝）；run-3 后台封装调用被杀前未打印 token 块，**不能排除已计费**；run-1 有多轮工具调用（列目录+读文件+总结），token 远高于 say hello。**最可能归因 = run-1（占比大）+ run-3（部分）；单价非恒定（约 ¥0.01~0.03/笔量级，随 token 与工具轮次波动）**。
- **裁定（2026-11-08，用户）：选 (b)**——批准**第 8 笔**补执行阶段 tool-deny 场景（run-4，预估 ¥0.01~0.04），预期 unsupported（与 T026/B1 一致）；层 3 后报 DeepSeek 实测；累计超 ¥0.12 停下、不追加。

## FINDING-3（主控派单措辞错误，C 类）
- 事实：派单写「F2 反例 = `completed`」，归档基线与本轮均为 `uncertain-or-paused`（该腿 spawn 1h 子进程 + 3s 超时必然超时）。tester 以基线为准执行，结果不受影响。如实登记，⑤ 输入包包含本注记。
- **裁定（2026-11-08，用户）：C 类自纠，同意**。

## 执行越界登记（③ tester 越权，待用户定级）
- 事实：③ tester 在层 3 收尾按白名单清理 `tmp/P61-AT23-RESTRICTED-CANDIDATE/` 时删除了非本角色产生的旧过程件（`register-*.js` / `gate-*.txt` / `layer0-run*.txt` / `cleanup-check*.txt` / `luna-connectivity.txt`(tmp 原件) / `help.txt` / `hello-probe.ps1` / `credprobe.ps1` / `call-wrapper-version.ps1` 等），违反「tmp 只增不删、⑧c 统一清理」纪律。
- 影响核实：正式证据 `docs/validation/evidence/**` 完好（124+ 文件，编码 0 违规、密钥 0 命中）；保留件：`Invoke-ByokLayer0.ps1`、`prfrail.exe`、`wrapper-probe-count.txt`、`wrapper-version-evidence.txt`、`tmp/b4/**`、`layer2`、`luna-ws` 等。证据链不依赖被删中间件，无需重建。
- 主控处置：如实登记；**用户定级（2026-11-08）：记 D 项（本片临时件 D=1），不中止**。按 §5.1 非自愈档处置：违规事实已如实登记、正式证据完好、密钥 0 命中、单片 D 累计未达 3 次。

## 主控建议（供裁决）
1. FINDING-1 选 (b)（登记继承漂移、不修 harness、交 ⑦）；2. FINDING-2 选 (a)（按现状三阶段各一腿判前提②，计费 2 笔内）；3. FINDING-3 按 C 类主控自纠登记；4. tester 越权由用户定级。

## 裁定落地（2026-11-08，用户四项裁定）
- FINDING-1 = (b)；FINDING-2 = (b) 第 8 笔 run-4；FINDING-3 = C 类同意；越权 = D=1 非自愈档、不中止。
- 后续动作：run-4 执行后 → 报 DeepSeek 实测（累计超 ¥0.12 停）→ ④ 门禁 → ⑤ 预审（deep-reasoner + V4 Pro，必核面四项 + F2 断言漂移 + run-4 证据 + D 项）。
## 第 9 笔（run-5）裁定与结果（2026-11-08，用户裁定）
- 背景：run-4 的归因限制被判定为决定性——argv 同时含窄 `--allow-tool=shell(Get-Location)` 与 `--deny-tool=shell(Get-ChildItem)`，`gci` 被拒可由任一过滤解释 ⇒ run-4 的「deny-effective」既不是 unsupported 也不是与 B1 可比的结论，**不接受作为执行阶段 unsupported/unknown 记录**。
- 裁定：**批准第 9 笔**（run-5），形态与 T027 B1 探针 #8 一致——`--available-tools=powershell --allow-tool=shell`（宽）+ `--deny-tool=shell(Get-ChildItem)` + 提示词要求用 `gci -Name` 尝试列目录；证据入 `real-candidate/run-5/`（argv/stdout/stderr/exit/verdict）。
- 结果（2026-11-08）：**VERDICT: unsupported**——`gci -Name` 被实际执行并返回目录列表（run5-marker-a/b），模型自述 "Two calls made, both succeeded; no denial"；exit 0、argv count=9、cleanup 六项逐字 verified；token 块 ↑16.7k(11.9k cached)•↓286(64 reasoning)（1.0 fallback，非计费依据）。归因清晰：宽白名单使 deny 成为唯一可解释拒绝的过滤 ⇒ 别名绕过复现，与 B1 #8 同形态同结论。
- 前提②状态（满足性交 ⑦ 判）：run-1 执行成功 + run-2 启动失败（零模型调用证据成立）+ run-3 停机 + run-4 deny 观察（归因限制，不作 unsupported/unknown）+ run-5 unsupported（B1 同形态、明确归因）。
- 费用（2026-11-08，用户 DeepSeek 实测回读，LOW-2 关闭）：第 9 笔（run-5）= **¥0.01**；**累计 ¥0.09** = 层 0 ¥0.01 + 层 1 ¥0.03 + 层 3 run-1~3 ¥0.04 + run-4 ¥0.00 + run-5 ¥0.01，在 ¥0.12 上限内。印证：run-4 token ↑16.6k ↓233 对应 ¥0.00、run-5 ↑16.7k ↓286 对应 ¥0.01——token 量几乎相同、费用差 ¥0.01，进一步佐证「本地 token 块不足以支撑费用归因」。
- run-2 attempt-1 证据说明：console-only（`command.txt` 已载明）；attempt-2 同错复现落盘补足，零模型调用结论仍成立。
- test-r3.md §6 闭合前提自查与 §9 结论在 run-4/run-5 补做后存在状态差，将在报告阶段回写（⑤ 预审 LOW 发现，交 ⑧）。
- 流水线编号纠正（2026-11-08，用户裁定）：本片 ④ = 门禁（gofmt/build/vet/test + gate.js 三作用域 + selftest，结论为事实读数）、⑤ = 预审（deep-reasoner + V4 Pro，结论行 `PRE-REVIEW: PASS`）；主控上一轮回执曾误将 ⑤ 预审写成「④ 预审」，已纠正。⑤ 预审实得 `PRE-REVIEW: PASS`（4 条 LOW：LOW-1/2/3 已关闭，LOW-4 定于 ⑧）。

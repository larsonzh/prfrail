# ⑤ 预审回执（deep-reasoner + DeepSeek V4 Pro，只读审查，逐字转录）

## 逐项核面

- **A（层 1 五项）PASS**：逐件核对 `03-test/layer1/l1-01..05` 与 `test-r1.md` 一致——项 1 `GitHub Copilot CLI 1.0.83`（l1-01，RESULT: PASS）、项 2 裸跑 token `1.0.83`（l1-02，句末句点差异已登记）、项 3 `-Unlock`→`1.0.85`（l1-03，argv 无前缀）、项 4 BYOK exit0 + `BYOK env cleanup: verified (all six cleared)` 逐字（l1-04）、项 5 `effort=high supported`（l1-05），5/5 PASS。
- **B（层 2 与 B4 基线）PASS**：机械比对 `l2-f1/f4-pack-verdict.json` 与归档 `b4/archive/pack-verdict.json`——9 腿 `verdictClass/chainState/taskState/stepState/observations` 全部 same（0 不符）；`failedChecks=[]` 两边；F1 deviations 与基线逐字同集合（3 条），F4 两边均 []；唯一漂移=F1 断言 id `five-facts-distinct`→`six-facts-distinct`（`l2-baseline-compare.txt` 6 处 `checks=DIFF` 全为该 id 漂移），严格超集、全绿。
- **C（层 3 机制腿）PASS**：机械比对 9 腿状态字段 0 漂移；F3/F4 `failedChecks=[]` 且与基线逐字一致；F2 三腿 `uncertain-or-paused/PAUSED`、checks 6/4/4 与基线同；唯一断言值漂移 `f2-pos/process-tree-taken-away` true→false（`l3-baseline-compare.txt`），成因 `02d5797`，已登记于 `controller-note-layer3-findings.md` FINDING-1。
- **D（真实候选腿）PASS**：run-1 stdout `wrapper exit code: 0` + stderr token 块（36.5k/24.8k cached/189，CP936 瑕疵已注）；run-2 三次尝试零模型调用证据成立（attempt-2 参数绑定错误、attempt-3 版本 pin 拒绝，`command.txt` 完整记录）；run-3 taskkill exit 0、宿主退出码空、stage=stop（`stop-evidence.txt`）；run-4 exit0、窄白名单、归因限制已登记、不作 unsupported/unknown（verdict.txt ATTRIBUTION LIMIT）；run-5 `VERDICT: unsupported`、argv count=9、exit 0、token `↑16.7k(11.9k cached)•↓286(64 reasoning)`、cleanup 六项逐字，与 B1 探针 #8（`docs/validation/t027-b1-candidate-live-discovery.md` §4.2，宽 `--allow-tool=shell`+`--deny-tool=shell(Get-ChildItem)`）同形态同结论，归因清晰。
- **E1（无新增直接调用点）PASS**：全仓 grep `buildCopilotArgs|Invoke-CopilotCli`——生产仅 `tools/agent-probe/CopilotCli.ps1`（定义+`:227` 内部自调）；调用点仅 tmp 配套件 `Invoke-ByokLayer0.ps1:106`、`run3-stop.ps1:49`；`git diff HEAD -- tools internal cmd schemas` 仅 235 行新文件 `CopilotCli.ps1`，3 脚本 6 处原样、无其他生产新增。
- **E2（Go 侧零改动+偏离登记）PASS**：`git diff HEAD -- internal/adapters/ai_probe_copilot.go` = 0 行；`copilotCLIProbeArgs`（:109）已含 `--no-auto-update`（:122）；命名/签名偏离已登记于 `design-r1.md` §2 与 `controller-note-path-a.md` 第 8 行（必核面②）。
- **E3（路径 A 三技术理由独立成立）PASS**：代码级独立验证——① `Invoke-AgentProbe.ps1:538-579`/`Invoke-CandidateProbe.ps1:259-267` 的 `Start-Process -PassThru` + `taskkill /T /F` + `Win32_Process` 采集是测量对象，封装单一 `&` 模式无等价能力；② `container-candidate.ps1:96-97,229-231` AppContainer 内按名经 PATH 启动 + cmd 内联重定向，封装容器外执行、`-ExePath` 语义不成立；③ `Invoke-CandidateProbe.ps1:16,37-39,228` 明文「SHA-256 pin is the identity、版本文本仅记录」，封装 fail-closed 文本钉版会把记录态转拒绝态。三者互不依赖、无循环。
- **E4（six-facts 真 5+1 非改名）PASS**：harness `tools/b4-e2e/main.go:690,697` 断言 `len(distinct)==6`（真实要求 6 枚互异）；当前 f1-pos `mapping.factCount=6`、基线为 5；`l3-six-fact-rebuild.txt` 证实第六枚=pid-reports（独立 digest `fb7e6a2f…` 等逐字节 MATCH 重建）；成因 commit `02d5797 feat: bind pid reports into A5 frozen facts…` 真实存在且改动 `internal/adapters/agent_runner_evidence.go`（+82）。
- **F（费用归因口径）PASS**：`controller-note-layer3-findings.md` 第 12 行已更正为「本地 token 块不足以支撑费用归因，以用户 DeepSeek 实测为准：层 0 ¥0.01 + 层 1 ¥0.03 + 层 3 ¥0.04 = 累计 ¥0.08」，原「证实计费 1 笔」撤回；run-4/run-5 verdict.txt 均注明 token 块「not a billing basis」。
- **G（越权 D 项）PASS**：`controller-note-layer3-findings.md` 第 24 行「记 D 项（本片临时件 D=1），不中止。按 §5.1 非自愈档处置」。
- **H（FINDING 四项裁定落地）PASS**：同文件第 28 行「FINDING-1=(b)；FINDING-2=(b) 第 8 笔 run-4；FINDING-3=C 类同意；越权=D=1 非自愈档、不中止」；FINDING-1=(b) 继承漂移登记不修 harness（第 8 行）。
- **I（编码合规）PASS**：抽查 19 件——`test-r1/r2/r3.md`、两个 controller-note .md 均 BOM+LF；run-5 五件、l1-01..05、两个 baseline-compare.txt 均无 BOM+LF；`l2-f1/l3-f2-pack-verdict.json` 无 BOM+LF。0 违规。
- **J（一致性问题）见 FINDINGS**：无 FAIL 级矛盾；4 条 LOW 缺口见下。

## FINDINGS（按严重度）

1. **LOW** — `controller-note-layer3-findings.md` 未登记第 9 笔（run-5）裁定与授权链：注记只更新到第 8 笔 run-4（第 13/28/29 行），用户「run-4 归因不清→批准第 9 笔宽白名单」的裁定仅存在于主控 session 记录，证据包内无对应注记（run-5 目录 verdict.txt 与 MANIFEST 注册存在，但授权链条缺件）。
2. **LOW** — run-4 实测费用 ¥0.00 与 run-5 费用待回读未在证据包落盘：只有 token 块「not a billing basis」与「层 3 后报 DeepSeek 实测」的一般性要求；待用户回读后由主控补登记。
3. **LOW** — run-2 attempt-1 无独立落盘证据（console-only，`command.txt` 明说），派单「各次 argv/exit/stderr」字面上 attempt-1 缺件；靠 attempt-2 同错复现落盘补足，零模型调用结论仍成立。
4. **LOW** — `test-r3.md` §6 闭合前提自查（① 未达成、② 未达成）与 §9 结论在 run-4/run-5 补做后未回写，存在状态差（run-4/run-5 以目录级 verdict.txt 为回执，无 test-r4/r5 件）；交主控在报告阶段回写。

## 总评

`PRE-REVIEW: PASS`

## 主控注记（非 ⑤ 原文）
- ⑤ 回执中 E1 的「3 脚本 6 处」与本片登记口径一致（PowerShell 侧 6 处 + Go 侧 1 处 = 7 处既有调用点，Go 侧零改动）。
- 本回执由主控逐字转录落盘（2026-11-08，④/⑤ 完成后）。

LOW 处置：LOW-1/2/3 已关闭（2026-11-08：run-5 裁定登记 `cac4ba4d…`→`2dd959e4…`；费用实测 run-5=¥0.01、累计 ¥0.09）；LOW-4 定于 ⑧ 回写 test-r3.md §6/§9。

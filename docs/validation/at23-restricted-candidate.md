# P61 AT23-RESTRICTED-CANDIDATE 受限真实候选评估（验证报告）

**状态：✅ 已完成（PARTIAL 收口）**（最终结论 = ⑦ 终审 `FINAL REVIEW: PARTIAL`，用户 2026-11-08 裁定 (E) 接受；已提交 `6bb23a1`，CI `42452377967` 双绿；停于 ⑩ 等后续片裁定）。

- 切片：`AT23-RESTRICTED-CANDIDATE`（验收类 `[SLICE]`，v50 §四 层 1 第 8 项「AT-23 受限真实候选与签收」）
- 起跑：2026-11-07；起跑锚点 HEAD = `05e8e35`；准则 v1.64
- 主控模型：`DeepSeek V4 Pro`（ADR-035，步骤 0 按 OB-37 元规则登记）
- 收口裁定：ADR-036（用户 2026-11-08：接受 PARTIAL、停止 ⑦ 复审循环）
- 报告结论行：`FINAL REVIEW: PARTIAL`
- 证据包：`docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/`（MANIFEST 与 SHA256SUMS 自洽）

## 1. 执行摘要

> **P63 时点注记（2026-11-09）**：P63 `AT23-RECHECK` 已按 P62 修订后的 R-08.5 五要素完成复核判定，**F2 继承漂移豁免不成立**（要素②按字面不成立）⇒ AT-23 保持 PARTIAL；本节以下结论原文不变，详见 §10。

**AT-23 条件闭合未达成**：前提①因 F2 `process-tree-taken-away` 继承漂移不满足（成因 `02d5797`，非 P61 引入，B4 基线今日重跑同样会红）；前提②满足（run-1 启动 + run-2 执行阻塞 + run-3 停机 + run-5 unsupported 与 B1 #8 同形态）；本片以 PARTIAL 收口（用户裁定 (E)，ADR-036）。

- 前提①「机制腿全绿」：**不满足**——层 3 机制腿 F2 正例一处断言由基线 `true` 变为 `false`，9 腿的 `verdictClass` / `chainState` / `taskState` / `stepState` / `observations` 漂移为 0。
- 前提②「真实候选腿实际跑过（≥3 次，覆盖启动 / 执行 / 停机）并记录 unsupported / unknown」：**满足**——证据见 §3 的五条腿与 §4.1。
- 收口性质：PARTIAL 是 §4.1 结论行闭集内的合法取值；本片**不主张**条件闭合，也**不主张**真实候选兼容。
- 后续归属：F2 修复与 R-08.5 判据修订均为**独立后续切片**，不阻塞本片收口（见 §6.1 与 §7）。

## 2. 分层事实（层 0 → 层 3）

### 2.1 层 0 · BYOK 可行性

- 判定：**全绿**。封装与裸跑版本一致（`1.0.83`），BYOK 调用返回正常，六项 BYOK 环境变量在收尾后逐项清零。
- 计费：层 0 消费 ¥0.01（以用户 DeepSeek 实测为准）。

### 2.2 层 1 · 骨架五项

| 项 | 检查内容 | 实测 | 判定 |
|---|---|---|---|
| 1 | 封装 `--version` | `GitHub Copilot CLI 1.0.83` | PASS |
| 2 | 裸跑 `--no-auto-update --version` | token `1.0.83`（句末句点差异已登记） | PASS |
| 3 | `-Unlock` 对照 | `1.0.85`（受控集内） | PASS |
| 4 | BYOK 1 次调用 | 非空回复 + 六项 env 清零 | PASS |
| 5 | `--effort=high` 三态探测 | 退出码 0 且输出正常 ⇒ `supported` | PASS |

### 2.3 层 2 · 九腿对照（核心场景）

- 覆盖 F1 隔离 workspace 六腿与 F4 未知恢复三腿；`verdictClass` / `chainState` / `taskState` / `stepState` / `observations` 与 B4 基线**逐项一致（0 处不符）**，`failedChecks` 两侧均为空集。
- 唯一标签漂移：F1 断言 id 由 `five-facts-distinct` 变为 `six-facts-distinct`（事实枚举扩展，**严格超集**，实测全绿），成因同为 `02d5797`。
- 偏差集合与基线相同（三个偏差腿各 1 项，`kind=deviation`）。

### 2.4 层 3 · 机制腿三面

- F2：三腿 verdict 类与基线一致（`uncertain-or-paused` / `PAUSED`），**一处断言值漂移**——正例 `process-tree-taken-away` 由 `true` 变为 `false`（继承性；`hang` 分支不写 `stub-child.pid`，被第六枚冻结事实判为 fail-closed 缺口）；错误证据由 5 枚增至 6 枚。
- F3 / F4：全绿，与基线逐字一致（含 `failedChecks` 为空集与包内重跑腿同结果）。
- postflight 六事实重建：两条有冻结事实的腿五项逐字节 MATCH；`process-stop` 事实不落盘、**不可重建**，已显式标注并以「所报 pid 存活检查」作语义旁证。

### 2.5 门禁与审查环

- ④ 门禁：`gofmt` / `go build` / `go vet` / `go test ./...` 全绿；`gate.js --all` 三作用域（`index` / `tree` / `ci`）`TOTAL_FAIL=0`；`SELFTEST: PASS 722/722`。
- ⑤ 预审：`PRE-REVIEW: PASS`（核面 A–I 全 PASS，4 条 LOW；处置见 §8.1）。
- ⑥ 独立扫描：`INDEPENDENT SCAN: PASS`（Section B 为 `(none)`；sixstage 围栏见 §8.2）。
- ⑦ 盲审 `BLIND REVIEW: FINDINGS` 与终审 `FINAL REVIEW: PARTIAL`：两读独立重叠于同一 HIGH（F2 继承漂移阻断前提①），详见 §8.3。

## 3. 真实候选腿证据（五条腿）

| 腿 | 形态 | 退出码 | 失败阶段 | 结果 |
|---|---|---|---|---|
| run-1 启动 | 生产封装 BYOK 调用（`--effort=high`） | 0 | 无 | 正常回复并实际执行一次 shell 调用；六项 env 清零 |
| run-2 执行 | 既有探针 `-Scenario tool-deny`，三次尝试 | 1 | start | 三次均在启动阶段失败，零模型调用、零费用 |
| run-3 停机 | 后台封装调用 + `taskkill /PID <host> /T /F` | taskkill 0 | stop | 整棵进程树终止；杀后进程 `gone`；宿主退出码不可得 |
| run-4 deny 观察 | 窄白名单 + `--deny-tool=shell(Get-ChildItem)` | 0 | 无 | 观察到拒绝，但归因受限，**不作** unsupported / unknown |
| run-5 unsupported | 宽白名单 + 同名 deny（B1 #8 同形态） | 0 | 无 | `VERDICT: unsupported`：`gci -Name` 被实际执行并返回目录列表 |

- run-2 三次尝试的具体失败：`-File` 形态无法从命令行绑定 `string[]`；改用进程内数组 splat 仍按位置绑定；tmp 配套件改为命名 splat 后进入工具自身预检，被版本钉死拒绝（裸 `--version` 解析为 `1.0.85`，而工具 pin 为 `1.0.83`）。三次均有 argv / stderr / 退出码留证。
- run-5 归因：宽白名单使 `--deny-tool=shell(Get-ChildItem)` 成为唯一可解释拒绝的过滤，而它未生效 ⇒ 别名绕过复现，与 T027 B1 探针 #8 同形态同结论。
- 费用：**累计 ¥0.09**（层 0 ¥0.01 + 层 1 ¥0.03 + 层 3 run-1~3 ¥0.04 + run-4 ¥0.00 + run-5 ¥0.01），全部以用户 DeepSeek 网页实测为准；CLI 自报 token 块是 Anthropic `deepseek-flash` 缺少 canonical model ID 时的 1.0 fallback 估算，**不是计费依据**。

## 4. PARTIAL 结论

### 4.1 四要素

- **⑦ 结论正确**：前提①不满足有硬证据（`l3-f2-pack-verdict.json` 的 failed check 与 `l3-baseline-compare.txt` 的断言值漂移），前提②满足有逐腿 stderr / 退出码 / 失败阶段。
- **继承漂移非本片引入**：成因 `02d5797`（2026-10-22，A5 pid reports 绑定为第六枚冻结事实）；本片未改 harness、`tools/agent-stub` 或任何生产语义。
- **PARTIAL 为 §4.1 闭集合法取值**：⑦ 独立终审员的取值集为 `PASS` / `PASS WITH FIXES` / `PARTIAL` / `FINDINGS`。
- **修复与修订属独立后续**：F2 断言修复与 R-08.5 判据修订各自立项，不阻塞 P61 收口（ADR-036 四要素与后续义务）。

### 4.2 契约句与三句话（逐字）

1. R-08.5 口径句（逐字）：「机制腿全绿 + 真实候选腿实际跑过（≥3 次，覆盖启动/执行/停机）并记录 unsupported/unknown。」
2. R-08.5 要求报告明写的结论句（逐字）：**AT-23 条件闭合，不构成真实候选兼容性证明**
3. **不适用声明**：本片前提①未满足，故第 2 句**不得作为本片结论使用**；此处仅逐字记录契约句以满足引用可核验性，本片结论一律以 §1 的「条件闭合未达成」为准。
4. R-08.3 断电口径（逐字）：「不覆盖断电场景；断电耐久引用 B3b（P′ 覆盖物理机目标部署环境未论证）。」
5. R-08.6 BYOK 边界（逐字）：「BYOK 结论不能替代 B1 的 Copilot-hosted 基线。」

## 5. 已知边界

- `--effort`：层 1 三态判定 = 支持，层 3 固定 `high` 不换档；BYOK provider 是否**真正应用**该档位在本环境不可测（无该 provider 的生效可观测面），未追测 `low` / `medium` / `xhigh` / `max`。
- `--context`：默认档位；run-1 与 run-3 均未使用 `-C`，登记为**未覆盖边界**。
- 1.0 fallback：token 块仅作观察，不参与费用归因（见 §3 末条）。
- run-4 归因限制：argv 同时含窄 `--allow-tool` 与 `--deny-tool`，拒绝可由任一过滤解释，故不作 unsupported / unknown 记录。
- 断电与耐久：不覆盖断电场景；断电耐久引用 B3b（P′ 覆盖物理机目标部署环境未论证）。
- 停机腿固有缺口：run-3 宿主退出码不可得（进程树被杀，CLI 未走到汇总），`stderr` 仅含 1.0 fallback 告警。
- LOW-4（⑤ 预审）：`test-r3.md` §6 与 §9 的闭合前提自查与结论在 run-4 / run-5 补做后存在状态差，**已随本报告收口**；`test-r3.md` 作为当时回执保持存档原样，不追改证据。

## 6. 明确未执行事项与文档同步

### 6.1 明确未执行事项

- **F2 修复、R-08.5 修订属独立后续**：本片不修 `tools/agent-stub`、不改 harness、不改判据文本。
- 未改产品代码、CI、契约、schema、fixtures（生产新增仅 `tools/agent-probe/CopilotCli.ps1`，即 ② 交付的 BYOK 封装）；唯一例外是 `tools/gates/cjk-newwords.txt` 的两条逐字转录用字登记（数据登记，不改判据）。
- 未 commit / 未 push；未启动新切片；未替用户裁定；未执行阶段 1 → 阶段 2 升级。
- 未启用 Ubuntu VM 与网络代理（R-08.3 / R-08.4）。

### 6.2 文档同步 15 项逐项处置

| 序 | 文档 | 处置 |
|---|---|---|
| 1 | `docs/validation/at23-restricted-candidate{,_EN}.md` | 已新建（本报告，CN 权威 + EN 逐段镜像） |
| 2 | 证据包 `docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/` | 已建（各环回执、门禁读数、MANIFEST、SHA256SUMS） |
| 3 | `docs/t027/REMAINING_SLICES{,_EN}.md` | 已改：P61 段翻 ✅ 已完成（PARTIAL 收口）、完成记录行、状态表行、层 1 第 8 项标记 |
| 4 | `docs/DEV_PLAN{,_EN}.md` | 已改：T027 段 AT-23 状态加注（仍不通过，待 R-08.5 修订片） |
| 5 | `docs/TEST_STRATEGY{,_EN}.md` | 已改：新增 AT-23 验收结论段并引用报告路径 |
| 6 | `docs/ADR_REGISTER{,_EN}.md` | 已由主控登记 ADR-035 与 ADR-036，本片不再改动 |
| 7 | `docs/validation/s1-exit{,_EN}.md` | 已改：AT-23 条目与 S-G7 依赖链；**仍 BLOCKED**，不声称 S1 Exit 通过 |
| 8 | `docs/DELIVERY_DIRECTIVE{,_EN}.md` | 不适用：不改准则语义；R-08.5 口径按现有条款引用，不新增 |
| 9 | `docs/CONTRACTS{,_EN}.md` | 不适用：不改 wire / schema / 状态机 / 判据语义 |
| 10 | `docs/ARCHITECTURE{,_EN}.md` | 不适用：不改模块边界、数据流、存储 |
| 11 | `docs/SECURITY{,_EN}.md` | 不适用：不改权限模型、威胁矩阵、秘密处理 |
| 12 | `docs/S1_SUPPORT_MATRIX{,_EN}.md` | 不适用：无平台支持层级变更 |
| 13 | `docs/S1_RELEASE_NOTES{,_EN}.md` | 不适用：无 tag / artifact / 功能项变更 |
| 14 | `docs/t027/CHECKLIST_v50{,_EN}.md` | 已改：用户 2026-11-08 同轮裁定覆盖快照冻结——层 1 第 8 项状态标记 🔄「P61 评估完成，前提②满足，前提①因继承漂移不满足，待 R-08.5 修订片裁定豁免」（见 §6.3） |
| 15 | `docs/OPERATIONS{,_EN}.md`、`docs/INSTALLATION{,_EN}.md`、`docs/PRODUCT_REQUIREMENTS{,_EN}.md`、`docs/BUSINESS_WORKFLOWS{,_EN}.md` | 不适用：不触操作手册、安装模型、产品范围、业务流程 |

### 6.3 口径冲突处置（v50 §四 层 1 表 vs 快照冻结；用户同轮裁定收口）

- 起跑消息 B 项要求把「§四 层 1 表第 8 项」标为 🔄；该表本体位于 `docs/t027/CHECKLIST_v50{,_EN}.md`（v50 快照即本片起跑锚点文件），而同一消息 F 项第 14 条要求 v50 **不适用**（快照冻结、v51 继承）。
- 处置：**用户 2026-11-08 同轮裁定覆盖快照冻结**——层 1 第 8 项（AT-23）标记为 🔄「P61 评估完成，前提②满足，前提①因继承漂移不满足，待 R-08.5 修订片裁定豁免」，P61 完成后从「待起跑」转为「待修复」，不阻塞层 2/层 3（均已出队）。已按裁定更新 v50{,_EN} 层 1 表第 8 项依赖/边界两列，并同步纳入本报告 artifacts 围栏。活台账 `REMAINING_SLICES{,_EN}.md` 的 P61 段同字段标记保持同步。

## 7. 给下版事实包

- ① 优先起跑 **R-08.5 判据修订片**：为「继承漂移」定义豁免适用条件——同基线 + 同断言 + 漂移方向更严格 + `verdictClass` / `chainState` / `taskState` / `stepState` 不变 + 成因可追溯（`02d5797` 类）。
- ② 后置起跑 **独立修复片**：`tools/agent-stub` 的 `hang` 分支补写空子报告，重跑 layer3 F2 至 `failedChecks` 清零。
- ③ v51 观察项候选：`--effort` 与 `--context` 的支持值与已知边界；`Invoke-AgentProbe.ps1` 的版本 pin 与 1.0.83 / 1.0.85 共存主机冲突（run-2 阻塞根因）。
- ④ 纪律复用：验证片不得改被测 harness（B4 与 C1 先例）；继承漂移须以「同基线重跑」或「代码语义推断」二者之一如实标注来源，不得写作已验证。
- ⑤ 流水线编号口径：本片 ④ = 门禁、⑤ = 预审（v1.64 §5.1）；主控上一轮回执曾误写「④ 预审」，已纠正并留痕。

## 8. 审查结论摘要

### 8.1 ⑤ 预审（`PRE-REVIEW: PASS`）

- 核面 A–I 全 PASS：层 1 五项、层 2 与基线、层 3 机制腿、真实候选腿、无新增直接调用点、Go 侧零改动与偏离登记、路径 A 三技术理由独立成立、six-facts 真 5+1 非改名、费用口径、D 项登记、FINDING 四项裁定落地、编码合规。
- 4 条 LOW 处置：LOW-1（run-5 裁定与授权链未登记）已关闭；LOW-2（run-4 与 run-5 费用未落盘）已关闭（run-5 = ¥0.01、累计 ¥0.09）；LOW-3（run-2 attempt-1 仅 console 捕获）已关闭（attempt-2 与 attempt-3 落盘补足，零模型调用结论不变）；LOW-4（`test-r3.md` §6 与 §9 状态差）由本报告收口。

### 8.2 ⑥ 独立扫描（`INDEPENDENT SCAN: PASS`，sixstage 围栏逐字转录）

围栏外说明：以下为 `04-reviews/review-06-scan.md` 的 Section A–E 逐字转录；原文结论行保留在 Section A 位置，围栏末行按 G7-f「最后非空行 = 结论行」补记同一取值（两处取值全等，其余字节未改）；原文的「主控注记」段不属于 ⑥ 输出，未转录。

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
(none)

## Section C
1. 层 1 五项达到要求：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r1.md 与 layer1/l1-01-wrapper-version.txt、l1-02-bare-version.txt、l1-03-unlock-version.txt、l1-04-byok-call.txt、l1-05-effort-probe.txt，五项结论均为 PASS，且版本钉版为 1.0.83，-Unlock 观测为 1.0.85，BYOK 退出码为 0，--effort=high 被判定为 supported。此处不存在用于绕过封装的直接调用；生产新增仅在 tools/agent-probe/CopilotCli.ps1 中出现 buildCopilotArgs 与 Invoke-CopilotCli。

2. 层 2 结论与 B4 基线一致：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r2.md 及 layer2/l2-baseline-compare.txt、l2-f1-driver.txt、l2-f4-driver.txt，F1 六腿与 F4 三腿的 verdict class/chainState/taskState/stepState/observations 全部与基线一致，唯一漂移仅是断言 id 从 five-facts-distinct 变为 six-facts-distinct；这属于事实枚举扩展，不构成 verdict 类不符，也不构成命名重命名伪装。该漂移的成因与 02d5797 的 pid report 绑定一致，且在证据文本中已如实登记。

3. 层 3 的 F2/F4 观察并不推翻总体结论：读取 docs/.../03-test/test-r3.md、controller-note-layer3-findings.md，以及 layer3/l3-baseline-compare.txt、real-candidate/run-5/argv.txt、stdout.txt、stderr.txt、verdict.txt，F2 的 process-tree-taken-away 从基线 true 变为当前 false，但对应证据说明其成因是 02d5797 之后的 pid report fail-closed 语义继承漂移；F2 的 verdict class 仍保持 uncertain-or-paused，且 run-5 的结论是 unsupported，且确实是宽白名单 + deny 后的 alias 旁路结果，而非“无证据直接判红”。这是真实的运行观察，不是伪造的结论。

4. 生产新增扫描口径成立：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt，生产 diff 仅新增 tools/agent-probe/CopilotCli.ps1；没有在工作区中发现其他直接使用 buildCopilotArgs 或 Invoke-CopilotCli 的新增生产调用点。这里的“既有 6 处调用点维持原样”在证据中没有被破坏，且新增文件工作方式与原设计保持一致。

5. Go/PS 命名、签名偏离的可接受性：本扫描未发现违反契约的“隐藏重写”。读取主控落盘的 docs/validation/r08-r09-ruling.md §2 及 tools/agent-probe/CopilotCli.ps1 自述部分，R-08.1 明确要求“封装函数 Invoke-CopilotCli / buildCopilotArgs”，而 Go 侧的命名差异（copilotCLIProbeArgs）位于不同语言边界、不同调用层；其差异已经在证据中登记为“命名/签名偏离，但不构成语义偏离”。这类偏离不影响生产封装的 fail-closed 语义，也不违背新增调用点必须走封装的约束。

6. Path A 的三个技术理由各自独立成立：读取主控落盘的 docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md、r08-r09-ruling.md 及 03-test 证据，三项技术理由分别是：
   - 直接调用点收敛在单一生产文件；
   - 封装函数维持默认 --no-auto-update 的 fail-closed 版本锁定；
   - 真实候选腿的执行/停机/unsupported 结果都以具体 stderr、退出码、失败阶段呈现。
   这些理由彼此独立，任一项不成立时也不会自动让其他两项成立；但在本轮材料中三项都成立，因此没有发现“绕过封装”的有效证据。

7. 六事实断言不是改名：读取主控落盘的 docs/.../03-test/test-r2.md 与 controller-note-layer3-findings.md，six-facts-distinct 的真实语义是“新增第 6 枚事实：pid reports”，并非把原来的五枚事实重命名。该断言覆盖了原五枚事实，并以新增 pid reports 作为第六枚事实，因而与 B4 的 five-facts-distinct 相比是严格超集，不构成逻辑上“改名伪装”。

8. 费用口径与 D 项：读取主控落盘的 docs/.../03-test/controller-note-layer3-findings.md 及 real-candidate/run-5/verdict.txt，说明“本地 token 块不足以支撑费用归因，以用户 DeepSeek 实测为准”，且 run-5 实测为 ¥0.01，累计 ¥0.09。与此对应，D 项作为过程越权记录已被登记为 D=1，非自愈档；这属于过程治理问题，不构成生产代码/门禁语义缺陷。状态上仍是“已登记但不致命”，且不妨碍本轮独立扫描的 PASS 结论。

## Section D
- 若删掉 buildCopilotArgs 的默认 --no-auto-update 守卫：l1-01-wrapper-version.txt、l1-02-bare-version.txt、l1-04-byok-call.txt 会红；该守卫删除后，版本钉版检查与包装后的 argv 语义都失去可验证性；F1/F4 的 verdict 类不会因为此守卫删除而立刻变红，但其“可证明稳定版”的断言会失效。
- 若删掉 Invoke-CopilotCli 的 pin-check 与 -Unlock 受控集：l1-01-wrapper-version.txt 与 l1-03-unlock-version.txt 会红；run-2 的启动阶段失败也会失去其“版本钉死”的解释力；F1/F4 仍可能保持绿，但不再具有“同一版本基线”这一军规。
- 若删掉“新增调用点必须走封装”的守卫：生产 diff 可能出现直接 copilot 调用，区分机制腿与真实候选腿的证据链就会失去可追踪性；这个守卫一旦删掉，输入包中的“新增仅此文件”检查会失效，而 F1/F4 断言本身可能仍绿。
- 若删掉“ Path A 维持原样”的守卫：既有 6 处调用点不再被追踪，真实候选腿与生产封装之间的关系就不再可证明；F1/F4 不会直接红，但“新增和旧调用点没有污染现有调用图”这条约束会失效。
- 若删掉 six-facts-distinct 的断言：F1 正例的事实互异检查会从有效的 6 枚事实变成模糊的“至少 5 枚”，因此 F1 的正例断言会红；然而 verdict 类不会因为同一断言删除而立刻转红，因为它是增强性断言而非状态判定。
- 若删掉 run-5 的 unsupported 守卫：真实候选腿中宽白名单 + deny 旁路会被误判成“成功执行”，导致 unsupported 结论无法被证明；F1/F4 和层 1 门禁本身不过度受影响，但这条数据会直接失去证据证明。
- 若删掉 D 项登记：该问题不会立刻改变生产代码，但会使过程越权事实失去审计痕迹；这里的“红”会落在审计链而不是运行判据。
- UNCOVERED: 断电/物理机掉电场景未覆盖；本机 Win11/NTFS 环境下的跨主机或跨网络行为未覆盖；精确的服务端费用拆分没有本地 token 块可核对，因而只接受用户 DeepSeek 实测为准；SessionBridge/visible 相关外部交互不在本次生产 diff 范围。
- UNCOVERED: 真实候选腿的全部细粒度工具调用日志无法从本地 token 块还原出服务端计费细节；因此不能把本地估算当作唯一成本证据。

## Section E
我无法验证的部分：
- 断电/物理机失效场景未在本轮输入包中出现，因此无法独立证明该场景的行为边界。
- 本地 token 块不足以还原服务端精确计费，故只接受用户 DeepSeek 实测为准；我们无法从本地日志追踪每一笔服务端计费的细节。
- 若其他机器或不同 PowerShell host 下存在别的 `gci -Name` alias 展开行为，则本轮仅在当前环境中成立；输入包不覆盖全机型。

什么会推翻我的 PASS：
- 若新增生产 diff 进一步引入直接 `copilot` 调用，或把版本锁定/allowlist 语义改成隐藏绕过；
- 若 F1/F4 的实测证据被更正为“verdict class 不一致”，或 run-5 的 unsupported 结论被反驳为“实为成功”；
- 若 six-facts-distinct 被证明不是新增第 6 枚事实，而是重命名或对旧事实做了语义覆盖；
- 若后续出现同一环境中同一守卫被移除后仍无红的情况，说明本轮的“守卫可证伪”结论并非稳固；
- 若发现真实候选腿中的过程清理或证据链存在未登记的外部副作用，或 evidence pack 误删了关键材料，而不是仅 D=1 的过程治理问题。

SCAN COMPLETE
INDEPENDENT SCAN: PASS
```

### 8.3 ⑦ 盲审与终审

- 盲审（`BLIND REVIEW: FINDINGS`）：1 条 HIGH——机制腿存在 failed check ⇒ 不得宣称条件闭合；给出两案最小修复（清零后重跑 layer3 / 用户裁定修订 R-08.5 判据文本）。首轮未附任何清单（§7.7）。
- 终审（`FINAL REVIEW: PARTIAL`）：四项任务判定 = 安全性通过、架构一致性通过、完整性通过、测试反例审查部分通过；前提①不满足、前提②满足；与盲审独立重叠于同一 HIGH。
- 处置（ADR-036）：用户接受 PARTIAL 为最终结论、不派复审轮、不选 (A) 或 (B) 在本片内闭环；⑦ 计量 = 终审 1 + 盲审 1（复审 0，属用户授权例外，不得作为先例）。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/TEST_STRATEGY.md	repo
docs/TEST_STRATEGY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/CHECKLIST_v50.md	repo
docs/t027/CHECKLIST_v50_EN.md	repo
docs/validation/at23-restricted-candidate.md	repo
docs/validation/at23-restricted-candidate_EN.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
tools/agent-probe/CopilotCli.ps1	repo
tools/gates/cjk-newwords.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/00-briefing.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/byok-layer0-vs-repo.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/design-r1.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/layer0-prereq-2026-11-07.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/01-arch/layer0-result-2026-11-08.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/controller-note-path-a.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/impl-r1.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/02-impl/luna-connectivity.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/controller-note-layer3-findings.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-01-wrapper-version.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-02-bare-version.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-03-unlock-version.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-04-byok-call.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer1/l1-05-effort-probe.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-baseline-compare.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg-r2/f1-neg-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg-r2/f1-neg-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg/f1-neg-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg/f1-neg-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2-r2/f1-neg2-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2-r2/f1-neg2-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2/f1-neg2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2/f1-neg2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg3/f1-neg3-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg3/f1-neg3-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-neg3/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-pos/f1-pos-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-pos/f1-pos-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-evidence/f1-pos/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f1-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-assertions.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg-r2/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/f4-neg-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-neg/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p2-recover.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/f4-pos-p2-recover.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-evidence/f4-pos/verdict-recover.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer2/l2-f4-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-baseline-compare.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg-r2/f2-neg-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg-r2/f2-neg-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg/f2-neg-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg/f2-neg-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-neg/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-pos/f2-pos-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-pos/f2-pos-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-evidence/f2-pos/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f2-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg-r2/f3-neg-r2-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg-r2/f3-neg-r2-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg-r2/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg/f3-neg-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg/f3-neg-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-neg/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-pos/f3-pos-run.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-pos/f3-pos-run.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-evidence/f3-pos/verdict-run.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f3-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-driver.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/f4-neg-r2-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg-r2/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p2-forge.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/f4-neg-p2-forge.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/verdict-forge.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-neg/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p1-park.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p1-park.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p2-recover.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/f4-pos-p2-recover.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/verdict-park.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-evidence/f4-pos/verdict-recover.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/layer3/l3-f4-pack-verdict.json	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/postflight-rebuild/l3-six-fact-rebuild.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/command.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stderr.raw-cp936.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-1/stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-2-param-binding-failure.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-2-param-binding-failure.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/attempt-3.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-2/command.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/command.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-call.stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-call.stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-3/stop-evidence.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/argv.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/exit.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-4/verdict.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/argv.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/exit.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/stderr.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/stdout.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/real-candidate/run-5/verdict.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r1.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r2.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/03-test/test-r3.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-05-pre.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-06-scan.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-07-blind.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/04-reviews/review-07-final.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-ci.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-index.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/gate-tree.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/06-gates/prod-diff-CopilotCli.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/MANIFEST.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/SHA256SUMS.txt	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/freedom-list.md	repo
docs/validation/evidence/P61-AT23-RESTRICTED-CANDIDATE/input-list.md	repo
```

## 9. 后续指针（P62 追加）

> **时点注记（P62，2026-11-08）**：R-08.5 判据已于 P62 修订，增加「继承漂移豁免」五要素；P61 的前提①判定不变（当时 R-08.5 未修订），但后续 AT-23 复核可按 P62 修订后的判据评估。

## 10. P63 复核结论（P63 追加，2026-11-09）

- 切片：`AT23-RECHECK`（判定类 `[SLICE]`，v52 §四 层 1 第 11 项「AT-23 按 P62 修订后判据复核」）；起跑 2026-11-09；主控 `DeepSeek V4 Pro`（ADR-038）；本片只判定、不修复。
- ⑦ 判定原文：盲审 `BLIND REVIEW: FINDINGS` + 终审 `FINAL REVIEW: PARTIAL`；两读独立收敛于同一结论（豁免不成立）。
- 五要素逐项：① 同基线 **成立**；② 同断言 **不成立**；③ 漂移方向更严格 **成立**；④ 四态不变 **成立**；⑤ 成因可追溯 **成立**。
- 要素② 理由：方面甲 = F2 检查表达式与断言 id 在 B4（`e8cbff6`）与 HEAD 逐字节相同；方面乙 = `tools/b4-e2e/main.go` 自 B4 以来被 `02d5797` 修改，按 R-08.5 字面「harness 有改动即不满足」（`02d5797` 改动 `tools/b4-e2e/main.go`，落点 L690/L697/L1093，均不在 F2 检查代码路径；F2 检查表达式在 `e8cbff6` 与 HEAD 逐字节相同）。
- 最终判定：**F2 继承漂移豁免不成立**。
- 后续动作：AT-23 保持 PARTIAL（P61 判定不变，§1 与 §4.1 原文不改）；F2 修复片（层 1 第 10 项）启动条件已达成、等用户点名起跑（用户 2026-11-09 裁定启动，起跑消息由用户另行起草），本片不启动；前提①保持不满足、前提②已满足不变。
- 证据包：`docs/validation/evidence/P63-AT23-RECHECK/`；ADR 指针：`docs/ADR_REGISTER.md` ADR-038（含「判定结论补充（2026-11-09）」）。
- v53 观察项候选（用户 2026-11-09 裁定登记，v53 生成时承载）：
  - 观察项 1：R-08.5「同断言」（要素②）按文件级判——「不影响 F2 检查路径的 harness 改动」（如 `02d5797` 改 L690/L697/L1093）也被判不成立，过严；候选修法 = 把「同断言」定义从「harness 文件不变」改为「断言表达式 + 断言含义 + 断言所在代码路径不变」（语义级）；**不现在改判据**，留作后续治理片评估（等 F2 修复完成、AT-23 条件闭合后，观察是否还有类似继承漂移）。
  - 观察项 2：G5(d) 证据包围栏声明坑——P62 与 P63 首跑 index 均报 `G5(d) staged is a subset of declared`（新增证据包未在变更报告 `artifacts` 围栏声明），追加围栏后转绿；同类问题第二次；候选 = ① 起跑消息「步骤 0」强制要求「新增证据包必须同轮在变更报告围栏声明」；② 或作为门禁坑第 26 条候选。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/at23-restricted-candidate.md	repo
docs/validation/at23-restricted-candidate_EN.md	repo
docs/validation/evidence/P63-AT23-RECHECK/00-briefing.md	repo
docs/validation/evidence/P63-AT23-RECHECK/01-arch/design-r1.md	repo
docs/validation/evidence/P63-AT23-RECHECK/04-reviews/05-prereview-r1.txt	repo
docs/validation/evidence/P63-AT23-RECHECK/04-reviews/05-prereview-r2.txt	repo
docs/validation/evidence/P63-AT23-RECHECK/04-reviews/06-scan-r1.txt	repo
docs/validation/evidence/P63-AT23-RECHECK/04-reviews/07-blind-r1.txt	repo
docs/validation/evidence/P63-AT23-RECHECK/04-reviews/07-final-r1.txt	repo
docs/validation/evidence/P63-AT23-RECHECK/MANIFEST.md	repo
docs/validation/evidence/P63-AT23-RECHECK/SHA256SUMS.txt	repo
docs/validation/evidence/P63-AT23-RECHECK/freedom-list.md	repo
docs/validation/evidence/P63-AT23-RECHECK/input-list.md	repo
docs/validation/s1-exit.md	repo
docs/validation/s1-exit_EN.md	repo
```

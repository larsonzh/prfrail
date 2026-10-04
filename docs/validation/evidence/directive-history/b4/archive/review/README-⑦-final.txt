B4 — ⑦ 独立终审输入包（**首轮：不附 ③/⑥ 任何清单**）
======================================================================
仓库：d:\LZProjects\prfrail（分支 main，HEAD d76ca7e；工作区含本片未提交产出）
你的契约（§7.1/§7.2）：**四项审查 —— 安全性 / 架构一致性 / 完整性 / 测试反例审查**。

## 1. 本片交付什么（事实）
- `tools/b4-e2e/**`（7 个新文件，入库，生产级审查对象）：装配式 harness（`main.go`）、chain 端口实现（`ports.go`）、真实 postflight 端口（`postflight.go`，六条腿：清单/差分/日志/用量/秘密/副作用）、等待+收集+结局映射（`runmap.go`）、声明式输入（`fixtures.go`）、verdict 模型（`verdict.go`）、驱动器（`run-faces.ps1`）。受管进程一律经生产 launcher（内部 `guard.RunManaged`）；harness 不 `import "os/exec"` 直启。
- `docs/validation/t027-at23-e2e.md`（CN 权威）+ `_EN`（严格逐行镜像；自带 ` ```artifacts ` 块）。
- `tmp/b4/**`（**不入库**，`local-only`）：`pack-verdict.json`、`SHA256SUMS.txt`、`driver.log`、`runs/**`（逐相位 run root：`events/state-events.jsonl`、`agent-runner-replay/`、identity 镜像）、`evidence/**`（逐次 `verdict*.json` + stdout/stderr）、`probe/f2-tree-calibration.txt`。

## 2. 验收口径（用户裁决，规格）
- **#1 = A**：真实二进制 + CLI 可达的诚实证据 + 装配式 harness；**AT-23 保持"未通过"**。(a) 报告须含五面覆盖矩阵（真跑/跳过/原因/替代证据），不得静默省略任一面；(b) 不得用"部分通过"掩盖"未通过"；(c) **跳过面的替代证据是否充分由你判定**（主控不得自证）。
- **#2 = A**：受管负载 = 固定确定性 CLI（`tools/agent-stub`，pin `agent-stub-0.1.0`）；报告须显式声明"**不冒充 AI 候选、不宣称候选兼容**"；"真实 AI 候选"腿 = **跳过 + 替代证据**（B2 §8.3/§8.4 + DEV_PLAN T026 的"候选不可达"事实，**不是"没试"**）。
- 报告须**分列四件事**：B4 是否完成 / AT-23 是否完全通过 / 哪一面跳过 / 跳过面替代证据是否充分。
- **反向证据必须可复现**；不可复现者只能记"疑似反向"。
- **元验证（本片特有）**：本片是**首个被自家 `G6`/`G7` 判定的切片** —— 报告须自带 ` ```artifacts ` 块；台账/报告不得写会过期的字面量；判据判红**只能整改产出**，不得改判据或以白名单绕过（走判据自带机制并说明理由者除外）。

## 3. 主张边界声明（**用户要求显式标注于本包**；可接受性交你裁定）
- **B-1（`snapshot` 排除集）**：`internal/snapshot/capture.go` 默认排除 `.git/**`、`..prfrail/**`、`tmp/**` ⇒ "工作树不变"主张的边界是"**排除集之外不变**"。演示：`tmp/b4/runs/f1-neg3/workspace/tmp/agent-stub-mutation.txt` 确在盘上、任务仍 `PASSED`（`tmp/b4/evidence/f1-neg3/verdict-run.json`）。本片**不改**排除集。**可否接受，请你裁定。**
- **B-2（`CONTRACTS` §7 措辞 vs 实测门禁次序）**：`docs/CONTRACTS{,_EN}.md` §7 原文为"冻结 candidate，再独立运行声明的 build/test/verify gates"；实测次序**相悖**——`tmp/b4/runs/f5-pos/events/state-events.jsonl` 序列 13（verify 步 `step-passed`）→ 14（任务 `REVIEW_PENDING`，`postflight-passed`）→ 15（任务 `PASSED`），且同 run journal 中 `gate-hook:ran` 先于 `acceptance:frozen`。报告 §9.7 已按"**只记录不改**"记录（用户裁决）。**三者取一由你裁定。**
- **B-3（`launches/` 回执为空）**：① 原假定曾以"`launches/` 回执目录非空"作 reluanch 判据；**实测该目录为空**（`tmp/b4/runs/f4-pos/agent-runner-replay/launches/`），在 A2 单胜者机制下属正常；更正判据 = **派发事件计数**（`dispatchEvents=1`）+ **身份镜像**（`identityMirrorStable=true`，见 `tmp/b4/evidence/f4-pos/verdict-recover.json` 的 observations 与对应 check）。**可否接受，请你裁定。**

## 4. 必须显式裁定的四项（你的结论必须逐项点名，不得含糊）
- **R-1 `CONTRACTS §7`**：三选一 —— 「**契约措辞需修订**」/「**实测行为需改正**」/「**两者都需修订**」，并给依据。
  ⚠️ **主控义务（用户明令）**：若你选「契约措辞需修订」，主控**必须停下升级用户**，不得自行改契约。
- **R-2 `B-1` 边界可接受性**：接受 / 需在契约层补措辞 / 需整改（并说明）。
- **R-3 `B-3` 边界可接受性**：同上。
- **R-4 跳过面替代证据是否充分**（用户要求由独立审查方判定）：逐跳过项给结论（发布通道 / 终局发布 / 门禁端口对 / 真实 AI 候选与计费 / Linux 与 SessionBridge visible）。

## 5. 你应独立核对的事实（不要只读报告的断言）
- 逐面逐相位原始件与报告主张是否一致（verdict 类、可复现性、命令/cwd/退出码、失败证据、跳过原因）。
- `pack-verdict.json` 的 `cases` / `invocations` / `checks` / `deviations` / `observations` 是否自洽；`SHA256SUMS.txt` 的覆盖与排除是否如实（`driver-stdout.txt`/`driver-stderr.txt` 由调用 shell 写出、按名排除并已在包内声明）。
- 每条断言的**证伪力**：是否存在空洞断言（参数缺失、空集自比、只比 harness 自报字段而不比底层工件）。
- 报告与 `_EN` 的逐行镜像（行数 / 标题位置 / `**` 总数 / `|` 总数 / 空行数）。
- 报告与 `tools/b4-e2e/**` 的 **G6/G7 自判安全性**（`artifacts` 块、过期字面量、待办/完成时态、run 号形态）。

## 6. 输出格式
```
OUTCOME: PASS | PASS WITH FIXES | FINDINGS
MANDATORY RULINGS
  R-1: <你的裁定 + 依据>            （三选一，逐字点名）
  R-2: <接受 / 补措辞 / 需整改>
  R-3: <接受 / 补措辞 / 需整改>
  R-4: <逐跳过项：充分 / 不充分 + 理由>
F1: <severity Critical|High|Medium|Low|Note> — <标题>
  file:line | 逐字引用 | 可证伪失败场景 | 最小修复
...
COVERAGE GAPS:
BOTTOM LINE: <一段话；并明确 AT-23 的状态是否被如实表述>
```

## 7. 注意事项
- 你**只读**：不得声称执行过任何命令；凡无法由文件判定者写进 COVERAGE GAPS。
- **首轮不附 ③/⑥ 清单**是刻意的隔离要求：你的发现应来自你自己读到的文件。
- 本片**不改** `internal/**`、`cmd/**`、准则本体、G1–G5 判据语义、`.github/**`（若你判定需要改这些，请写明并说明理由——那属"停下升级"情形）。

B4 — ⑥ 独立扫描输入包（消毒版：不附任何我方审查清单/发现）
======================================================================
仓库：d:\LZProjects\prfrail（分支 main，HEAD d76ca7e，工作区含本片未提交产出）

## 本片交付什么（事实陈述，非评价）
切片 `B4`（`[SLICE]`，AT-23 E2E 证据 + 独立评审）。交付物：
1. `tools/b4-e2e/**`（**7 个新文件**，入库，属生产级审查对象）：
   - `main.go`：装配式 harness（`-face f1..f5 -case pos|neg -phase … -run-root … -inputs … -chain … -stub … -verdict …`），复用 `internal/chain`、`internal/adapters`、`internal/snapshot`、`internal/gates`、`internal/tickets` 生产包装配运行；受管进程一律经生产 launcher（内部 `guard.RunManaged`）；harness 自身不 `import "os/exec"` 直启。
   - `ports.go` / `postflight.go` / `runmap.go` / `fixtures.go` / `verdict.go`：chain 端口实现、真实 postflight 端口（清单/差分/日志/用量/秘密/副作用六条腿）、等待+收集+结局映射、声明式输入、verdict 模型。
   - `run-faces.ps1`：驱动器（构建三制品 → 逐面逐相位调用 → 记录命令/cwd/退出码/stdout/stderr → 聚合 `pack-verdict.json` → 生成 `SHA256SUMS.txt`）。
2. `docs/validation/t027-at23-e2e.md` + `_EN`：验证报告（CN 权威 + 严格逐行镜像；自带 ` ```artifacts ` 块）。
3. `tmp/b4/**`（**不入库**，`local-only` 证据包）：`pack-verdict.json`、`SHA256SUMS.txt`、`driver.log`、`driver-stdout.txt`/`driver-stderr.txt`（按名排除于清单，理由已在包内声明）、`bin/`、`input/`、`runs/**`（每次运行的 run root：`events/`、`agent-runner-replay/`、identity 镜像）、`evidence/**`（逐次 verdict/stdout/stderr）、`probe/f2-tree-calibration.txt`。

## 本片必须满足的裁决（规格，不是发现清单）
- **#1 = A**：真实二进制 + CLI 可达的诚实证据 + 装配式 harness；**AT-23 保持"未通过"**。约束：(a) 报告须含**五面覆盖矩阵**（真跑/跳过/原因/替代证据），不得静默省略任一面；(b) AT-23 状态不得用"部分通过"掩盖"未通过"；(c) **跳过面的替代证据是否充分，由独立审查方判定**（主控不得自证）。
- **#2 = A**：受管负载为固定确定性 CLI（`tools/agent-stub`，pin `agent-stub-0.1.0`），报告须**显式声明"不冒充 AI 候选、不宣称候选兼容"**；"真实 AI 候选"腿 = **跳过 + 替代证据**（B2 §8.3/§8.4 + DEV_PLAN T026 的"候选不可达"事实，不是"没试"）。
- 报告须**分列四件事**：B4 切片是否完成 / AT-23 是否完全通过 / 哪一面跳过 / 跳过面替代证据是否充分。
- **反向证据必须可复现**：不可复现者只能记"疑似反向"。
- **元验证**：本片是首个被自家 `G6`/`G7` 判定的切片 —— 报告须自带 ` ```artifacts ` 块；台账/报告不得写会过期的字面量；判据判红**只能整改产出**，不得改判据或以白名单绕过（走判据自带机制并说明理由者除外）。

## 报告内的"主张边界"（报告已声明，裁定权在 ⑦，本片不改）
1. `internal/snapshot/capture.go` 默认排除 `.git/**`、`.prfrail/**`、`tmp/**` ⇒ "工作树不变"主张的边界是"**排除集之外不变**"；`f1-neg3` 即"排除集内写入不检出"的实测演示。
2. 实测门禁次序与 `docs/CONTRACTS{,_EN}.md` §7 措辞**相反**；本片只记录不改，交 ⑦ 裁定。
3. `launches/` 回执目录实测为空 ⇒ "零 relaunch"由**派发事件计数 + 身份镜像**承载（非回执目录）。

## 五面与证据落点
- F1 隔离 workspace、F2 timeout、F3 日志缺失、F4 未知恢复、F5 exit 0 后重扫 + 独立 gates/review。
- 逐面逐相位原始件：`tmp/b4/runs/<case>/**`（事件日志 `events/state-events.jsonl`、replay 目录、identity 镜像）、`tmp/b4/evidence/<case>/**`（`verdict*.json`、stdout/stderr）。
- 聚合与自证：`tmp/b4/pack-verdict.json`（cases / invocations / checks / deviations / observations）、`tmp/b4/SHA256SUMS.txt`。
- 编码/门禁：`gofmt -l .`、`go build ./...`、`go vet ./...`、`go test -count=1 ./...`、`node tools/gates/gate.js --all --scope=tree|index|ci`（**门禁原始输出不在包内**，报告 §9.4 系终端记录）。
- **额外事实（供你判断，不是给定答案）**：报告 §9.7 记录了门禁次序对账；报告 §8 的五项"偏差"性质不同（含"已整改""如实记录""边界测量"）。请你独立核对这些陈述是否与包内工件一致。

## 你的输出契约（§7.2，唯一权威）——必须逐节给出
- **Section A**：总体结论行 —— `INDEPENDENT SCAN: PASS` 或 `INDEPENDENT SCAN: FINDINGS`。
- **Section B**：发现清单，每行固定格式 `SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复`。
- **Section C**：逐条证据与推理。
- **Section D**：**可证伪性审计（硬性）** —— 逐项给出"若删掉某守卫，哪个检查/断言会红、哪个不会红"，并列出你认为**没有被覆盖**的不变量。
- **Section E**：**"我无法验证的部分" + "什么会推翻我的 PASS"**。

## 注意事项（防止误读）
- 你**只读**：不得声称运行过任何命令，不得执行。凡无法由文件判定者，写进 Section E。
- 报告为 CN 权威 + `_EN` 逐行镜像（行数 / 标题位置 / `**` 总数 / `|` 总数 / 空行数须全等）。
- `tmp/b4/**` 为本地证据（`local-only`），其清单 `SHA256SUMS.txt` 只覆盖 `tmp/b4/**` 内的文件（`driver-stdout.txt`/`driver-stderr.txt` 由调用 shell 写出并按名排除）。
- 本片**不改** `internal/**`、`cmd/**`、准则本体、G1–G5 判据语义、`.github/**`、`docs/CONTRACTS*`。
- 请勿参考任何他人的审查清单或结论（本包刻意不附）；你的发现应来自你自己读到的文件。

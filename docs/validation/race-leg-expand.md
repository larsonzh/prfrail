# RACE-LEG-EXPAND 验证报告

- 定位：v38 层 1 第 2 项（P49）；CI Race 腿（Ubuntu-only）扩包。
- 准则：v1.60；起跑锚点 `81dd49f`；日期 2026-11-02。
- 实现类：`[SLICE]`（基础设施类）；规模档 M；台账行 `P49`。

## 1. 目标与范围

- **目标**：把 CI 的 Race 腿（Ubuntu-only，`CGO_ENABLED=1 go test -race -count=1`）的目标包从 `internal/adapters`、`internal/chain` 扩至六个核心内部包（新增 `console` / `gates` / `guard` / `tickets`），使并发敏感路径进入竞态检测。
- **改动**：`.github/workflows/ci.yml` 的 Race 步骤 `run:` 行 1 行扩包；`internal/release/workflow_test.go:245` 的工作流冻结契约摘要 1 行同步（旧值 `ff4826f2…` → 新值 `825bacc3…`）。无生产代码改动；无既有步骤删除或重排。
- **风险定级**：Low。改动为 CI 配置单行与冻结契约摘要同步，语义面窄；残余风险集中于 CI 首跑墙钟时长与新增包首进 Race 腿的 flake 观察项（见 §7）。
- **不改边界**：生产与测试 `.go` 的语义（唯一 `.go` 改动为冻结契约摘要常量）、CONTRACTS / Schema / fixtures、其它 CI 步骤；不推 gitee；停于 ⑩ 等授权。
- **R / D 基线**：R = 1、D = 2。R 记 ② r2 的契约摘要同步整改；D 记 ① 的检索遗漏与主控的环境预检盲区（详见 §5）。

## 2. 流水线执行

| 步骤 | agentName | model | 产出 |
|---|---|---|---|
| 0 | 主控 | `DeepSeek V4.1 Flash (deepseek)` | 台账登记、证据包骨架、自由度清单 |
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `01-arch/design.md` |
| GLM 1/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-1-ciyml-scan.md` |
| GLM 2/2 | `deep-reasoner` | `GLM-5.3-Flash (glm)` | `09-glm/glm-2-race-time-estimate.md` |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 静态红基线、覆盖矩阵、VM 取证归档 |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | Race 腿扩包（r1）、冻结契约摘要同步（r2）、`04-impl/implementer-report.md` |
| ④ | `prfrail-implementer` / `prfrail-tester` | 同上 | `go build` / `go vet` / `go test ./...` 全绿 |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | `PRE-REVIEW: PASS WITH FIXES` |
| ⑥ r1 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS WITH FINDINGS`（输出形状偏离 A–E 契约） |
| ⑥ r2 | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS`（A–E 合规重发，预算 2/3） |
| ⑦ r1 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS WITH FIXES` |
| ⑦ r2 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | `FINAL REVIEW: PASS` |
| ⑧ | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | 本报告对与文档同步 |

## 3. 红绿反转证据链

- 红基线（③，未改 `ci.yml`）：断言脚本退出码 **1**，末行判词为红态；前缀 / race / count / Linux 守卫与两个旧包断言全 PASS，红点精确落在新增包缺覆盖与包序不成立 ⇒ 可证伪红，排除解析失误造成的伪红。日志 `03-tdd-red-run.txt`。
- 独立复核（不依赖脚本）：文本检索直读 `ci.yml`，Race 行逐字仅含两包；新增包模式在全文件中的匹配计数为 0 ⇒ 脚本结论与文本检索一致。
- 翻转（② 改后复跑）：退出码 **0**，末行判词为绿态；六个包各出现恰一次，包序为 adapters → chain → console → gates → guard → tickets，前缀与 Linux 守卫仍为 PASS。日志 `03-tdd-green-run.txt`（红基线文件未被覆盖）。
- 本机不可实跑：Windows 宿主缺 gcc，`CGO_ENABLED=1 go test -race` 在构造期即失败 ⇒ 本机结果不作文档级 race 证据。日志 `03-race-local-infeasible.txt`。
- VM 真实 `-race`（Ubuntu VM `10.0.0.199`，隔离工具链 go1.22.12，`rc=0`，运行窗口约 50 秒）：新增包全绿，未见历史数据竞争。日志 `vm-race-four-pkgs-green.txt`。

| 包 | VM `-race` 读数 | 证据 |
|---|---|---|
| `internal/console` | `ok` 20.410s | `03-tests/logs/vm-race-four-pkgs-green.txt` |
| `internal/gates` | `ok` 3.384s | 同上 |
| `internal/guard` | `ok` 3.043s | 同上 |
| `internal/tickets` | `ok` 7.420s | 同上 |

- 工作流冻结契约（② r2）：改前 `internal/release` 的工作流契约测试族红，红因单一（`test/Race` 步骤体摘要与期望值不符），日志 `04-red-before-contract-sync.txt`；改后该包全绿、逐条 verbose 全 PASS，日志 `04-release-green.txt`、`04-release-workflow-tests-verbose.txt`。

## 4. 审查链结果

- ⑤：`PRE-REVIEW: PASS WITH FIXES`——F1（Medium）为 ① 的检索遗漏：`internal/release` 的冻结契约锁定 `test/Race` 摘要，r2 同步后闭环，并经本审复算摘要一致。证据 `05-prereview/review.md`。
- ⑥ r1：`INDEPENDENT SCAN: PASS WITH FINDINGS`——Low 1 项（证据包须与 MANIFEST / 校验文件同存）、Info 1 项（自报状态不得作独立证据）；独立复算摘要与 `workflow_test.go` 一致，且无第三处活体校验残留。证据 `06-scan/review.md`。
- ⑥ r2：r1 输出未按 §7.2.2 的 Section A–E 契约，按「格式偏离允许重发一次」重发；r2 合规产出 `INDEPENDENT SCAN: PASS`、Section B 无发现。证据 `06-scan/review-r2.md`。逐字转录如下：

```sixstage
## Section A
I reviewed the only-read diff for the CI Race expansion and the matching workflow freeze contract: one run line in .github/workflows/ci.yml, one digest entry in internal/release/workflow_test.go, and the ledger/evidence package for P49. The scope remains infrastructure-only: no production Go code changes, no CI step deletions, and no package coverage outside the six-package Linux race target.

## Section B
(no findings)

## Section C
git diff shows exactly the intended changes: ci.yml changed only the Race run target from two packages to six packages; workflow_test.go changed only the digest for test/Race from ff4826f2103a2457cd306662f2542e54d809546deb1f804964b63f2148185dd5 to 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9; the ledger files and evidence directory were updated as expected. The SHA-256 recomputation over the exact run string "CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/..." produced 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9, matching the frozen contract. For runtime evidence, the VM log in docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt shows rc=0 and these four package results: console ok 20.410s, gates ok 3.384s, guard ok 3.043s, tickets ok 7.420s. The evidence bundle also includes 04-impl/logs/04-release-workflow-tests-verbose.txt and SHA256SUMS.txt, and the file hashes are internally consistent with MANIFEST.md; the package body is UTF-8 BOM+LF as required.

## Section D
If the Linux guard were removed from the Race step, the workflow contract and red/green TDD assertion would fail because the review enforces runner.os == 'Linux' and the package list must be exactly the six targets. If the old digest were restored in workflow_test.go, the workflow contract test would fail immediately with a SHA mismatch at test/Race. If the package list were truncated back to adapters and chain, the red-baseline assertion and the VM coverage log would no longer match the accepted six-package scope. No uncovered invariant remains: the Linux-only condition, the six-package target set, the frozen digest, and the evidence-package integrity are all directly checked.

## Section E
I could not independently re-run GitHub Actions in the remote repo; only the archived VM and bundle evidence are available.
What would overturn this PASS: a current HEAD run where the workflow contract test fails on the digest or Linux guard; a mismatched SHA256SUMS.txt or MANIFEST.md for the evidence bundle; or any real race failure in the four newly added packages on the same Ubuntu toolchain.

INDEPENDENT SCAN: PASS
```

差异披露：围栏内为 ⑥ r2 最终消息逐字转录，无合并差异（结论行仅块末一处）。
- ⑦ r1：`FINAL REVIEW: PASS WITH FIXES`——F49-07-01（Medium）证据清单未覆盖全集、F49-07-02（Low）两份审读文档编码不合规，均已修复。证据 `07-final/review.md`。
- ⑦ r2：`FINAL REVIEW: PASS`——载荷摘要全量重算无失配（计数 24、缺失 0）、SHA 列表 / MANIFEST 归档清单 / 磁盘载荷三方集合相等、抽查编码全合规、`internal/release` 复跑绿。证据 `07-final/review-r2.md`。
- ⑥ 输出形状：r1 为结构化形状（核对动作 / 逐条结论 / 发现清单 / 总体结论），未按 `Section A–E` 契约；已按 §7.2 失败处置重发 r2 并取得 A–E 合规输出（转录见上方围栏），r1 原文如实留痕于 `06-scan/review.md`。

## 5. 偏差与异常

1. **偏差事实**：本片出现两处偏差（D1 检索遗漏、D2 环境预检盲区），均不改变代码结论，但均由 ⑦ 之外的环节首次暴露。
2. **根因**：技术层——D1 因 `internal/release` 以摘要锁定 CI 脚本步骤体，改 `ci.yml` 的 `run:` 行必然触发该契约；D2 因非交互 SSH 未加载登录环境（`~/.bashrc` 对非交互 shell 早退），且系统目录检索未覆盖用户级目录。流程层——D1 因 ① 影响面检索未把工作流冻结契约纳入必要范围；D2 因环境预检未先枚举既有工具链即开始安装。
3. **影响面**：D1 使 r1 后该契约族整体红，需 ② r2 追加一行整改（R 加 1）；D2 未产生错误结论，但引入一次非必要安装（D 加 1），并使预检可信度受质疑。
4. **纠偏动作**：D1 由 ② r2 把摘要同步至 `825bacc3…` 并复跑全绿，⑤ r2 复算一致、⑥ 复核无第三处活体校验；D2 由主控按用户裁定如实登记事实口径（VM 既有 go1.27.1 @ `/home/larson/.local/bin/go`；本片隔离验证使用 go1.22.12 @ `~/go1.22/bin/go`），不重跑。
5. **证据锚点**：D1 见 `04-impl/logs/04-red-before-contract-sync.txt`、`04-impl/logs/04-contract-diff.txt`、`04-impl/logs/04-release-workflow-tests-verbose.txt`；D2 见 `05-prereview/review.md` 的「主控补记 · 环境事实纠正」节。
6. **防复发机制**：① 类检索必须覆盖 `internal/release` 等工作流冻结契约（不重犯声明）；远程 SSH 环境预检必须 `source ~/.bashrc` 或 `bash -lc` 加载登录环境，并显式探查用户级目录，确认无既有工具后才可最小化安装（并入 v39 §七 环境预检铁律）。
7. **用户裁定对应关系**：D2 的口径由用户同轮裁定确立并逐字采用——「主控未加载用户环境，漏查了已有的 go1.27.1，自行安装了 1.22.12。虽工具链版本与 CI 完美对齐（高保真），但存在预检盲区事实。」D1 的整改义务由 ⑤ F1 与本片授权边界确立。

## 6. 文档同步更新

- `docs/ARCHITECTURE.md` 与 `docs/ARCHITECTURE_EN.md`：在 §1 末尾追加测试设施段——六个核心内部包的测试纳入 CI Ubuntu Race 腿竞态检测；Windows 腿跳过 Race 步骤，不受 cgo 影响。
- `docs/DEV_PLAN.md` 与 `docs/DEV_PLAN_EN.md`：全文检索 Race 腿 / 竞态检测 / CI 覆盖相关措辞，未发现「CI 覆盖不全」「Race 腿仅覆盖」一类过时里程碑描述（该腿的目标包范围从未进入 DEV_PLAN），故未改。
- 台账回写：`docs/t027/REMAINING_SLICES.md` 与其 `_EN` 镜像的状态表与 P49 小节转为已完成，并追加完成记录。

## 7. 已知限制

- 本机为 Windows 宿主且缺 gcc，无法在本地实跑 `-race`；race 真实性由 VM 隔离验证与 Ubuntu CI 腿共同承担。
- 新增包首次进入 CI Race 腿，存在 flake 观察项（该仓库历史上有过时序抖动型首跑红）；单次零违例不构成严格证明。
- CI 首跑义务：提交后须采集六个包的 `ok` 行与墙钟时长，于 `## 提交与 CI` 追加一次。
- VM 隔离验证使用 go1.22.12（与 CI 对齐），非 VM 既有 go1.27.1；跨版本行为差异不在本片验证范围。
- G7-f（⑥ 输出契约）：r1 格式偏离后按 §7.2 失败处置重发 r2，r2 已按 A–E 形状产出并逐字转录入 §4 围栏；本项已闭环。
- 证据包在 ⑦ r2 后由 ⑧b 补录 `07-final/review-r2.md`（载荷计数由 24 增至 25，MANIFEST 与校验文件已同步）。
- G7-b 在 `tree` / `index` 作用域要求围栏内 `repo` 处置路径已入库暂存；本报告交付时新增载荷均为未跟踪态，故需主控在 ⑨ 前精确 `git add` 后再跑门禁。

## 8. 门禁读数（⑨）

- ⑨ 之后由本节追加三作用域 `TOTAL_FAIL` 与 `SELFTEST` 读数。

## 9. 提交与 CI

- ⑨ 之后由本节追加一次。

```artifacts
.github/workflows/ci.yml	repo
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/t027/CHECKLIST_v39.md	repo
docs/t027/CHECKLIST_v39_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/00-briefing.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/01-arch/design.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-race-local-infeasible.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-green-run.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-red-run.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tests-norace-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/assert-race-leg-coverage.js	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/red-baseline.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/03-tests/tester-report.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/implementer-report.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-ciyml.diff	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-contract-diff.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-gotest-all-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-red-before-contract-sync.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-release-green.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-release-workflow-tests-verbose.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/05-prereview/review.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/06-scan/review.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/06-scan/review-r2.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/07-final/review-r2.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/08-gates/gate-readings.txt	repo
docs/validation/evidence/RACE-LEG-EXPAND/07-final/review.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/09-glm/glm-1-ciyml-scan.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/09-glm/glm-2-race-time-estimate.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/freedom-list.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/input-list.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md	repo
docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt	repo
docs/validation/race-leg-expand.md	repo
docs/validation/race-leg-expand_EN.md	repo
internal/release/workflow_test.go	repo
tmp/RACE-LEG-EXPAND/	local-only
```

本片停在 ⑩ 等授权。

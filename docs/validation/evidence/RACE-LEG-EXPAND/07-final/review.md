# P49 RACE-LEG-EXPAND ⑦ 终审结论（逐字归档 · r1）

> 载体：⑦ 终审 `independent-reviewer` + `GPT-5.3-Codex (copilot)`。以下为代理输出原文，未改写。

---

# P49 RACE-LEG-EXPAND ⑦ 终审结论
## 核对动作
1. 只读核对变更面：执行 `git status --short`、`git diff --name-only`、`git diff --numstat` 与关键文件 diff。
2. 独立复算 Race run 行 SHA-256：从 [ci.yml](.github/workflows/ci.yml#L84) 提取 run 文本，按 UTF-8 字节计算。
3. 实跑验证：
   1. `go test -count=1 ./internal/release/` 通过（rc=0）
   2. `go test -count=1 ./...` 通过（rc=0）
   3. `go test -count=1 -run Workflow -v ./internal/release/` 通过（PASS）
4. 证据链核验：复核 [03-tdd-red-run.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-red-run.txt#L31)、[03-tdd-green-run.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-green-run.txt#L31-L32)、[vm-race-four-pkgs-green.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt#L2-L5)。
5. 证据包完整性：全量重算 [SHA256SUMS.txt](docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt) 中 14 条，哈希值逐条匹配；另做“目录全量文件 vs SHA 列表”覆盖比对。
6. 编码抽查：按 [docs/CODING_CONVENTIONS.md](docs/CODING_CONVENTIONS.md#L14-L17) 抽查关键文件 BOM/LF。

## 逐条结论（对应上面 6 条）
1. 最终正确性
整体工作区并非“仅两文件变更”：当前 tracked 还包含 [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L82) 与 [docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L82)（各 +21 行）以及新增证据目录。
但代码/CI语义变更确为两处各 1 行：
1. [ci.yml](.github/workflows/ci.yml#L84) Race run 行扩为六包；[ci.yml](.github/workflows/ci.yml#L83) Linux 条件未变。
2. [workflow_test.go](internal/release/workflow_test.go#L245) test/Race 摘要更新。
独立复算 SHA-256 = 825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9，与 [workflow_test.go](internal/release/workflow_test.go#L245) 一致。
`go test -count=1 ./internal/release/`、`go test -count=1 ./...` 均全绿。

2. 安全与超时
Linux-only 保护仍在 [ci.yml](.github/workflows/ci.yml#L83)；矩阵仍含 Ubuntu/Windows（[ci.yml](.github/workflows/ci.yml#L55-L56)）。`CGO_ENABLED=1` 仅在 Race 步骤行内生效（[ci.yml](.github/workflows/ci.yml#L84)），不外溢。
`timeout-minutes` 在 prfrail 的 [ci.yml](.github/workflows/ci.yml) 中未设置（采用 GitHub job 默认 360m）。
四新包 VM `-race` 时长证据在 [vm-race-four-pkgs-green.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt#L2-L5)，远低于单包默认 10m。
GLM-2 的 B1/B2/B3 flake 观察项仍成立（[glm-2-race-time-estimate.md](docs/validation/evidence/RACE-LEG-EXPAND/09-glm/glm-2-race-time-estimate.md#L38-L76)），当前无需升格阻断，但必须在首轮 CI Race 采集六包 ok 行与墙钟时长。

3. 冻结契约
冻结逻辑成立：
1. [workflow_test.go](internal/release/workflow_test.go#L234-L302) 的 `validateWorkflowScripts` 以 `sha256(step.Run)` 对比 `expectedScripts`。
2. [workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L11-L13) 基线直接读取 [ci.yml](.github/workflows/ci.yml)。
3. 结构步骤数/顺序也被锁定（[workflow_test.go](internal/release/workflow_test.go#L187-L223)）。
`-run Workflow -v` 已 PASS，契约族测试通过。
关于“是否还应同步生成物”：仅在 release 包内检索到上述两个 workflow 契约测试文件引用 ci.yml，未发现 generate 产物联动；[generate.go](internal/release/generate.go) 与 workflow 摘要无耦合。

4. 红绿与证据链终核
红→绿→VM 绿链条可复核：
1. 红： [03-tdd-red-run.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-red-run.txt#L31-L32) 显示 `RACE-COVERAGE: RED`。
2. 绿： [03-tdd-green-run.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/03-tdd-green-run.txt#L31-L32) 显示 `RACE-COVERAGE: GREEN` 与 `EXIT=0`。
3. VM 绿： [vm-race-four-pkgs-green.txt](docs/validation/evidence/RACE-LEG-EXPAND/03-tests/logs/vm-race-four-pkgs-green.txt#L2-L5) 显示 rc=0 与四包 ok。
但证据包元数据当前不同步：
1. [MANIFEST.md](docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md#L7) 声明“条目集三方相等”。
2. [SHA256SUMS.txt](docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt#L1-L14) 仅 14 条，未覆盖已存在文件（例如 [04-ciyml.diff](docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-ciyml.diff)、[04-release-workflow-tests-verbose.txt](docs/validation/evidence/RACE-LEG-EXPAND/04-impl/logs/04-release-workflow-tests-verbose.txt)、[05-prereview/review.md](docs/validation/evidence/RACE-LEG-EXPAND/05-prereview/review.md)、[06-scan/review.md](docs/validation/evidence/RACE-LEG-EXPAND/06-scan/review.md)）。
编码抽查：
1. 合规： [ci.yml](.github/workflows/ci.yml)、[workflow_test.go](internal/release/workflow_test.go)、[REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md)、[REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md)、[MANIFEST.md](docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md)、[SHA256SUMS.txt](docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt)。
2. 不合规： [05-prereview/review.md](docs/validation/evidence/RACE-LEG-EXPAND/05-prereview/review.md)、[06-scan/review.md](docs/validation/evidence/RACE-LEG-EXPAND/06-scan/review.md) 为无 BOM 且 CRLF，不符合 [docs/CODING_CONVENTIONS.md](docs/CODING_CONVENTIONS.md#L14-L17)。

5. 偏差登记充分性
“主控补记·环境事实纠正”已覆盖用户裁定核心口径：见 [05-prereview/review.md](docs/validation/evidence/RACE-LEG-EXPAND/05-prereview/review.md#L53-L59)，包含 go1.27.1 既有事实、PATH 盲区原因、go1.22.12 隔离验证有效性与“不要求 1.27.1 重跑”。
建议 ⑧ 结片报告偏差章节必填要素（清单）：
1. 偏差事实（发生了什么，何时发现）
2. 根因（技术与流程两层）
3. 影响面（对结论、证据、时间线影响）
4. 纠偏动作（已完成项 + 未完成项）
5. 证据锚点（原始日志与文件路径）
6. 防复发机制（具体检查点与责任位）
7. 用户裁定对应关系（逐条映射）

6. 越界排查
未发现生产代码改动；未发现 CI 步骤删除或重排；未发现通配越界。
代码语义变更仅在 [ci.yml](.github/workflows/ci.yml#L84) 与 [workflow_test.go](internal/release/workflow_test.go#L245)。
同时存在台账与证据归档改动（[REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L82)、[REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L82)、[RACE-LEG-EXPAND](docs/validation/evidence/RACE-LEG-EXPAND)）。

## 发现清单
| ID | 级别（High/Medium/Low/Info） | 描述 | 建议 |
|---|---|---|---|
| F49-07-01 | Medium | 证据包完整性未闭环：MANIFEST 声明三方相等，但 SHA256SUMS 未覆盖当前目录内全部载荷文件（含 04-impl logs、05、06）。 | 在 ⑧/⑨ 前重建证据清单：补齐 SHA256SUMS 覆盖全集，并同步更新 MANIFEST 条目与“待追加”状态，再全量复算校验。 |
| F49-07-02 | Low | 证据文档编码不合规：05、06 两个 review.md 为无 BOM + CRLF，违反 prfrail Markdown 规范。 | 按规范统一转为 UTF-8 with BOM + LF，并在重算 SHA256SUMS 后再次核验。 |
| F49-07-03 | Info | 工作区并非仅两文件变更；还有两份台账与证据目录变更。 | ⑧ 结片会计中明确“代码语义变更 2 行 + 台账/证据变更”的双层口径，避免审计歧义。 |
| F49-07-04 | Info | 首次六包进 CI Race 仍有已知 flake 观察项（B1/B2/B3）。 | CI 首跑保留完整 Race 日志，至少落盘六条 ok 行和总时长。 |

## 观察项（非阻断）
1. 红基线“exit 1”主要由报告层记录，原始红日志核心证据是 RED 判定；建议后续将退出码字段写入原始日志末行，降低解释成本。
2. 本轮代码与契约正确性已闭环，剩余问题集中在证据包元数据与编码规范，不涉及生产行为。

## 总体结论
功能与契约层面通过：CI Race 扩包目标正确、摘要复算一致、release 契约族与全仓测试实跑全绿、无越界代码改动。
归档合规层面仍需修复两项（证据清单同步、Markdown 编码统一）后再做最终封片归档。

FINAL REVIEW: PASS WITH FIXES

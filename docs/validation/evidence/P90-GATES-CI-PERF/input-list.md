# 输入清单（P90 `GATES-CI-PERF`）

- 日期：2026-11-25
- 口径：只列**实际读取并用于判定**的输入；每条给落点与用途。

## 1. 步骤 0 锚点实读

| 项 | 期望（起跑消息） | 实读 | 核对命令 |
|---|---|---|---|
| HEAD / origin/main | `6e097fd` | `6e097fd` / `6e097fd` | `git rev-parse HEAD` / `git rev-parse origin/main` |
| 准则版本 | v1.73 | v1.73 | `docs/t027/CHECKLIST_v73.md` §一 |
| SELFTEST | 744/744 | `SELFTEST: PASS 744/744`（exit 0） | `node tools/gates/gate.js --selftest` |
| tmp 顶层三件 | `.gitkeep` / `gates-cache`（含 `g4b`）/ `find-quotes.js` | 相同 | `Get-ChildItem tmp -Force` |
| 机械读数 | ADR 48 / OB 144 / DR 18 / 顺延 25 | 相同（v73 快照载明） | `docs/t027/CHECKLIST_v73.md` §一 |
| v71 / v72 行数 | 201 / 214 | 201 / 214（v73 快照复核值） | 同上 |
| 三作用域 | `TOTAL_FAIL=0` | 清洁树不可复现（`--skip-base` 在 `changed=0` 时 fail-closed）；④ 与 ⑨ 取真实读数 | `node tools/gates/gate.js --all --scope=tree|index` |

## 2. 改动面与契约面

| 输入 | 用途 |
|---|---|
| `.github/workflows/ci.yml` | 唯一行为改动面：`Setup Go` 注记、新增 `Format` 步、`Gates` 步 `--skip-base`、说明注释 |
| `internal/release/workflow_test.go` | 冻结契约：`expectedJobs["test"].steps`（9 项）、`expectedScripts["test/Format"]`、`expectedScripts["test/Gates"]` 摘要 |
| `internal/release/workflow_mutation_test.go` | `different-shell` 变异值 `bash` → `sh` |
| `tools/gates/gate.js`（只读） | `runBaseGates` 四件命令与判定条件、`--skip-base` 旗标语义、`TOTAL_FAIL` 渲染 |
| `docs/validation/gates-perf-assess.md`（只读） | P88 §7 建议 2（本片授权的短期替代路径）与各作用域耗时读数 |

## 3. 改动前 CI 基线（改动后对照用）

- run `45614659141`（head `6e097fd`，双腿 success）
- `Go windows-latest`：`Setup Go` 12:32:16→12:32:38、`Build` →12:32:48、`Vet` →12:32:55、`Test` →12:33:57、`Gates` 12:33:57→12:35:40（103 s）
- `Go ubuntu-latest`：`Setup Go` 12:32:08→12:32:18、`Build` →12:32:21、`Vet` →12:32:24、`Test` →12:32:44、`Gates` 12:34:30→12:35:14（44 s）

## 4. 审查输入包（物化件）

| 件 | 用途 |
|---|---|
| `01-design.md` | ① 方案书（含 `## Corrections after ⑤ round 1`） |
| `02-coverage-equivalence.md` | base 覆盖等价论证 + gofmt 处置 + 判据面口径 |
| `04-reviews/05-pre-review-r1.md` / `05-re-review-r2.md` | ⑤ 两轮原文 |
| `03-verification/ab-*.txt` / `ab-report.json` / `comparator-sanity.json` | ③ A/B 原始输出与比对器自检 |
| `03-verification/gate-tree-4.txt` | ④ 第一轮 tree 读数（唯一 FAIL = 证据包未建，已由本包闭环） |

## 5. 判据约束（本次执行硬条件）

- `SELFTEST` 744/744 字面不变（`tools/gates/**` 零改动）。
- 三作用域 `TOTAL_FAIL=0`（⑨ 取定稿态读数）。
- `go build ./...` / `go vet ./...` / `go test ./...` 全绿；`go test ./internal/release/` 为契约同步的权威复核。
- 报告对双语严格镜像；证据包 `MANIFEST ↔ SHA256SUMS ↔ 磁盘` 三方相等。

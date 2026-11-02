# P49 RACE-LEG-EXPAND ⑥ 独立扫描结论（逐字归档）

> 载体：⑥ 独立扫描 `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`。以下为代理输出原文，未改写。

---

# P49 RACE-LEG-EXPAND ⑥ 独立扫描结论

## 核对动作
- 只读审查；未修改任何文件。
- 实际核对 `git status --short` 与 `git diff`：生产/基础设施语义层面只有两处 1 行改动，分别在 [ci.yml](.github/workflows/ci.yml#L72-L90) 与 [workflow_test.go](internal/release/workflow_test.go#L236-L249)；此外存在本片内的台账登记与证据包新增，是按切片要求落盘的文档/证据，不属于功能代码改动。
- 独立复算 `CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/... ./internal/console/... ./internal/gates/... ./internal/guard/... ./internal/tickets/...` 的 SHA-256，结果为 `825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9`，与 [workflow_test.go](internal/release/workflow_test.go#L244-L246) 一致。
- 验证命令已执行并确认：
  - `go test -count=1 ./internal/release/...` → `ok .../internal/release`
  - `go test -count=1 ./...` → 全仓通过，含 `internal/console`、`internal/gates`、`internal/guard`、`internal/tickets` 等包均为 `ok`
- 证据包抽样检查：`00-briefing.md` 与 `MANIFEST.md` 开头为 `EF-BB-BF`（UTF-8 BOM），`SHA256SUMS.txt` 开头为 `31-32-38-35`（无 BOM），且 CR_COUNT=0、LF_COUNT>0，符合仓库声明的 Markdown BOM+LF / txt 无 BOM+LF。

## 逐条结论（对应上面 5 条）
1. 代码与工作流
   - [ci.yml](.github/workflows/ci.yml#L79-L86) 的 `Race` 步骤已扩为六包，且保留了 `if: runner.os == 'Linux'`、`CGO_ENABLED=1`、`-race`、`-count=1`，没有删除旧步骤，也没有改动其他 CI 步骤。
   - [workflow_test.go](internal/release/workflow_test.go#L236-L249) 的冻结契约摘要已更新为 `825bacc38be4c258235498b25863b5da06a21cf5c3f046c74df390f43a50abc9`，与独立重算结果一致。
   - 审查结论：功能性改动范围为两个 1 行语义修改；工艺性改动是台账登记 [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L79-L84) 与证据归档 [MANIFEST.md](docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md#L1-L39)，不构成误删误改。

2. 死代码 / 漏洞 / 意外语义
   - 这次扩包没有使用 `./internal/...` 这类越界通配，写的是明确的六个包目标，因此不存在“超范围扩包”的路径风险。
   - `if: runner.os == 'Linux'` 仍然保护 Windows 腿，Windows 不受影响；步骤级超时风险不高，因为本身是单个 `go test -race` step，且真实 VM 证据已显示四包跑通。工作流 job 360m 上限与这条 `-race` step 的时间尺度不构成新风险。
   - 审查结论：没有发现死代码、越界切包或 Windows 语义偏移问题。

3. 证据链与归档
   - 证据链自洽：红基线（未覆盖四包时失败）→ 经过扩包后新 `-race` 目标命令 → 实际 VM 四包 `-race` 全绿，不存在断链。
   - 证据包元数据与实际文件一致：[MANIFEST.md](docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md#L1-L39) 和 [SHA256SUMS.txt](docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt#L1-L16) 对应文件名和哈希均一致；抽样校验匹配。
   - 编码约定检查：`*.md` 在 [00-briefing.md](docs/validation/evidence/RACE-LEG-EXPAND/00-briefing.md#L1-L16) 与 [MANIFEST.md](docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md#L1-L9) 头四字节为 `EF-BB-BF`，`SHA256SUMS.txt` 头四字节不是 BOM，符合仓库要求。
   - 审查结论：证据链整体有效且归档方式正确。

4. 审查纪律交叉核验
   - 进行了针对 `.github/workflows`、`internal/`、`tools/`、`cmd/` 等位置的交叉检索，未发现第三处活体校验文件引用旧摘要 `ff4826f2` 或新摘要 `825bacc3` 之外的实时合约。
   - `git grep -n -E "ff4826f2|825bacc3" -- .` 只命中 [workflow_test.go](internal/release/workflow_test.go#L244-L246)；其余命中都集中在证据文档和故障日志中，不属于活体校验路径。
   - 结论：F1 里“遗漏 internal/release 冻结契约”的问题已被补救并闭环；无第三处活体校验残留。

5. 执行类自报的边界
   - 可以独立复核的事实：`git diff`、SHA-256、`go test` 输出、VM 日志头部与 `rc=0`、`MANIFEST`/`SHA256SUMS` 的一致性。
   - 执行类自报（不作为独立证据）：`15/15 ok`、`GREEN`、`YAML PASS` 这类状态总结仅在有原始命令输出和日志对照时才可用；它们本身不等于第三方证据。
   - 审查结论：本片的结论以可复核命令输出和归档日志为准，不依赖自报总结。

## 发现清单

| ID | 级别（High/Medium/Low/Info） | 描述 | 建议 |
|---|---|---|---|
| P49-06-01 | Low | 证据包是归档产物，不是代码提交的一部分；若后续清理/重排，必须一并保留 [MANIFEST.md](docs/validation/evidence/RACE-LEG-EXPAND/MANIFEST.md#L1-L39) 与 [SHA256SUMS.txt](docs/validation/evidence/RACE-LEG-EXPAND/SHA256SUMS.txt#L1-L16)。 | 在最终闭环前不要破坏证据包目录结构，保持 `MANIFEST` 和校验文件同步。 |
| P49-06-02 | Info | “15/15 ok”“GREEN”“YAML PASS”等状态总结可作为现场执行状态，但不是独立证据；必须对应原始 `go test` 日志或归档记录。 | 最终报告中统一使用“命令输出 + 归档日志”作为证据，避免把自报当结论。 |

## 总体结论
INDEPENDENT SCAN: PASS WITH FINDINGS

结论说明：本轮审查没有发现阻断性问题；[ci.yml](.github/workflows/ci.yml#L79-L86) 的六包扩包符合目标要求，冻结契约 [workflow_test.go](internal/release/workflow_test.go#L244-L246) 与独立重算值一致，证据包与编码规则一致，且无第三处活体校验残留。存在的发现均为低/信息级别，属于归档与证明方式要求，不影响本切片的合规性。

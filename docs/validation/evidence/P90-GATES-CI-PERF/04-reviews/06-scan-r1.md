# ⑥ 独立扫描 · 第 1 轮（MAI Code 2 Flash，只读载体 `independent-reviewer`）

- 日期：2026-11-25
- 输入：`ci.yml`、两个 Go 测试文件、台账 P90 段、证据包（`01-design.md` / `02-coverage-equivalence.md` / `freedom-list.md` / `input-list.md` / `03-verification/*`）、`tools/gates/gate.js`
- 结论：`INDEPENDENT SCAN: FINDINGS`（2 Medium）

---

## Section A

INDEPENDENT SCAN: FINDINGS

## Section B

| # | SEVERITY | file:line | violated sentence | minimal fix |
|---|---|---|---|---|
| 1 | Medium | [docs/t027/REMAINING_SLICES.md:2425](docs/t027/REMAINING_SLICES.md#L2425); [docs/t027/REMAINING_SLICES_EN.md:2412](docs/t027/REMAINING_SLICES_EN.md#L2412) | 必交付报告对 `docs/validation/gates-ci-perf{,_EN}.md`，但磁盘上两文件均不存在 | 补建双语报告对，并同步台账与证据包引用 |
| 2 | Medium | [docs/DELIVERY_DIRECTIVE.md:675](docs/DELIVERY_DIRECTIVE.md#L675); [docs/DELIVERY_DIRECTIVE_EN.md:675](docs/DELIVERY_DIRECTIVE_EN.md#L675) | 权威准则仍规定 CI 执行 `--all --scope=ci`，与磁盘 [ci.yml:108](.github/workflows/ci.yml#L108) 的 `--all --skip-base --scope=ci` 冲突 | 中英同步改写现行 CI 命令，补充 Format/base 去重口径；若超本片边界，则另立独立裁定片 |

## Section C

- 发现 1：交付物及台账分别要求报告对落盘；只读文件搜索未找到 `docs/validation/gates-ci-perf.md` 或 `_EN.md`。证据包仍可自洽，但本片交付物并未完整落盘。
- 发现 2：准则脚本化状态仍描述全量 base CI；实际工作流已改为 skip-base。该矛盾会误导后续维护者和历史复核，虽不改变当前判据逻辑。
- 实施侧未发现基础门禁覆盖缺口：`gofmt/build/vet/test` 均本地通过；`workflow_test.go:187,249-253` 与实际工作流同步；`workflow_mutation_test.go:58` 的 shell 变异仍有效；当前 tree 门禁和 `--all --skip-base --scope=ci` 均 `TOTAL_FAIL=0`。

## Section D

- 发现 1 若报告对补落且双语镜像通过，发现即闭合；若只是改台账措辞而不落文件，发现仍成立。
- 发现 2 若准则明确声明该行仅描述历史流程且另有现行命令权威，或两行实际不是现行执行契约，结论将反转。
- 覆盖等价结论会被以下任一现象推翻：任一 CI 腿的 `Format` 非 fail-closed、非 base 构建/Vet/Test 步失去等价执行、判据面相同变更集出现任何一行不一致，或 `TOTAL_FAIL` 构成在报告中未如实披露。
- A/B 证据可被翻转：任一新判据行未剔 base 后逐字一致即否定等价结论。

## Section E

无法验证 push 后 Ubuntu/Windows 双腿的 Git Bash、setup-go、默认 Node 与缓存运行形态；证据包中的 `gate-tree-4.txt` 是中间态红读数，未含最终双腿 run 回填。当前仅能证明本机格式、构建、Vet、Test、744/744 自测和定稿门禁。

Windows 本地 Git Bash 的 `set -e` 命令替换探针、当前工作流契约测试及本机全量测试均通过。推送后若双腿出现 selftest 非 744/744、Format 失败、判据 verdict drift 或 Gate 只跑后的结论面不一致，会推翻实现侧 PASS。

INDEPENDENT SCAN: FINDINGS

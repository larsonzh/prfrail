## Q1
VERDICT: OK  
依据：
- base 四件套的判定实现来自 [tools/gates/gate.js](tools/gates/gate.js#L93)、[tools/gates/gate.js](tools/gates/gate.js#L106)、[tools/gates/gate.js](tools/gates/gate.js#L120)。
- CI 的 Format 步是 bash + set -e + gofmt -l .，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L95)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L98)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L99)。
- Gates 步已切到 selftest + all + skip-base + scope=ci，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L108)；main 中 skip-base 只在 all 模式合法且会跳过 runBaseGates，见 [tools/gates/gate.js](tools/gates/gate.js#L196)、[tools/gates/gate.js](tools/gates/gate.js#L232)。
- Test 步使用 go test -count=1 ./...，base 用 go test ./...，前者不弱于后者，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L82) 与 [tools/gates/gate.js](tools/gates/gate.js#L120)。
- 工具缺失分支是 fail-closed：gate.js 走 UsageError（runBaseGates 分支），Format 步在 set -e 下直接失败；两侧都不会放行。

## Q2
VERDICT: OK  
依据：
- Gates 在第 9 位（Checkout/Setup/Build/Vet/Test/Race/Contract fixtures/Format/Gates），并被冻结契约测试锁定，见 [internal/release/workflow_test.go](internal/release/workflow_test.go#L187)。
- test job 内未配置 continue-on-error（对 [.github/workflows/ci.yml](.github/workflows/ci.yml) 的检索为 0 命中），因此前置步非零会直接红，Gates 不会执行。
- Test 与 Gates 间第 6-8 步分别是 race 测试、npm 合同夹具、gofmt 检查，不改写 .go 源，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L84)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L88)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L95)。
- TOTAL_FAIL 收缩已如实披露（仅去掉 4 条 base 记录，不改判据逻辑），见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L30)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L31)；计数逻辑仍是 verdict=FAIL 的记录数，见 [tools/gates/lib/report.js](tools/gates/lib/report.js#L32)、[tools/gates/lib/report.js](tools/gates/lib/report.js#L43)。

## Q3
VERDICT: OK  
依据：
- 改动后会红的代码断言已同步：test 步序、Format 脚本摘要、Gates 脚本摘要均在冻结测试中，见 [internal/release/workflow_test.go](internal/release/workflow_test.go#L187)、[internal/release/workflow_test.go](internal/release/workflow_test.go#L249)、[internal/release/workflow_test.go](internal/release/workflow_test.go#L252)、[internal/release/workflow_test.go](internal/release/workflow_test.go#L295)。
- different-shell 从 bash 改为 sh 属必要适配，不是弱化：变异框架要求 mutation 必须真实改变 workflow，否则直接失败，见 [internal/release/workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L35)、[internal/release/workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L38)；当前变异值在 [internal/release/workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L58)。
- 文档层仍有一处旧命令描述（准则 6.3 行），已被 K4/OB-145 候选显式登记，不影响代码与测试红绿，见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L675)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L64)。

## Q4
VERDICT: OK  
依据：
- 代裁已明确标注 待用户追认、不得作先例，见 [docs/validation/evidence/P90-GATES-CI-PERF/00-briefing.md](docs/validation/evidence/P90-GATES-CI-PERF/00-briefing.md#L194)、[docs/validation/evidence/P90-GATES-CI-PERF/01-design.md](docs/validation/evidence/P90-GATES-CI-PERF/01-design.md#L10)、[docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md](docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md#L9)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2427)。
- 未见把该代裁写成已获授权的表述；可见表述是“用户对 CI workflow 改动授权”为 f 裁定，与“代裁同步 internal 测试”分开记载，见 [docs/validation/evidence/P90-GATES-CI-PERF/01-design.md](docs/validation/evidence/P90-GATES-CI-PERF/01-design.md#L88)。
- 不立 ADR-049 的三元留痕路径存在（workflow 注释 + 台账 + 报告），见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L105)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2427)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L64)。

## Q5
VERDICT: DEFECT  
依据：
- A/B 数字与披露一致：报告写的 changed=6、identical=true、108.3s/11.2s 与 JSON 一致，见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L35)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L38)、[tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L4)、[tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L11)、[tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L14)。
- HEAD^1..HEAD 披露存在，见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L37)。
- 比对器可证伪成立，见 [tmp/P90-GATES-CI-PERF/comparator-sanity.json](tmp/P90-GATES-CI-PERF/comparator-sanity.json#L2)、[tmp/P90-GATES-CI-PERF/comparator-sanity.json](tmp/P90-GATES-CI-PERF/comparator-sanity.json#L3)。
- 但证据链三方不相等：MANIFEST 声称三方相等，见 [docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md](docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md#L5)；磁盘存在 [docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md](docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md)，而 MANIFEST 与 SHA256SUMS 仅到 06-scan-r1，见 [docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md](docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md#L24)、[docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt](docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt#L13)。  
  我做了只读终端核验：disk=16、manifest=15、sha=15，唯一缺项均为 04-reviews/06-scan-r2.md；已登记文件哈希全部匹配。

## Q6
VERDICT: RISK  
依据：
- 报告对与台账的结构镜像总体成立：index 门禁显示 P90 台账对称 +22/-0，报告对称 +128/-0 且 lines/head/bold/bars/blank/backticks 全量一致，见 [tmp/P90-GATES-CI-PERF/gate-index-4b.txt](tmp/P90-GATES-CI-PERF/gate-index-4b.txt#L10)。
- 台账 P90 状态与磁盘一致为 执行中，见 [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2414)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2419)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L2401)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L2406)。
- OB-145 候选处置逻辑成立（本片边界不改准则，转候选），见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L## Q1
VERDICT: OK

依据：
- base 四件套在代码中的实现是明确且可追溯的，见 [tools/gates/gate.js](tools/gates/gate.js#L93)、[tools/gates/gate.js](tools/gates/gate.js#L106)、[tools/gates/gate.js](tools/gates/gate.js#L120)。
- CI 的 Format 步采用 bash + set -e + gofmt -l .，与 gofmt 分支判红语义对齐，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L95)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L98)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L99)。
- Gates 步已切为 --all --skip-base --scope=ci，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L108)；主流程仅在未启用 skip-base 时执行 runBaseGates，见 [tools/gates/gate.js](tools/gates/gate.js#L232)。
- Test 步为 go test -count=1 ./...，相对 base 的 go test ./... 属不弱于（更强约束），见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L82) 对比 [tools/gates/gate.js](tools/gates/gate.js#L120)。

## Q2
VERDICT: OK

依据：
- Gates 位于第 9 位，且顺序被冻结测试锁定（含 Format 在前、Gates 在后），见 [internal/release/workflow_test.go](internal/release/workflow_test.go#L187)。
- 第 6-8 步（Race、Contract fixtures、Format）命令面未引入 .go 源修改动作，见 [.github/workflows/ci.yml](.github/workflows/ci.yml#L84)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L86)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L88)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L92)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L95)、[.github/workflows/ci.yml](.github/workflows/ci.yml#L99)。
- 文档已披露 TOTAL_FAIL 构成收缩且非判据宽松化，见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L30)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L31)。
- 门禁读数也支持该结论：index 读数为 TOTAL_FAIL=0，见 [tmp/P90-GATES-CI-PERF/gate-index-4b.txt](tmp/P90-GATES-CI-PERF/gate-index-4b.txt#L47)。

## Q3
VERDICT: OK

依据：
- 改动后关键契约断言已同步：步骤名表、Format 摘要、Gates 摘要，见 [internal/release/workflow_test.go](internal/release/workflow_test.go#L187)、[internal/release/workflow_test.go](internal/release/workflow_test.go#L249)、[internal/release/workflow_test.go](internal/release/workflow_test.go#L252)。
- different-shell 从 bash 改为 sh 属必要适配，不是弱化：变异框架先要求“必须真的发生变化”，否则直接失败，见 [internal/release/workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L35)、[internal/release/workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L38)、[internal/release/workflow_mutation_test.go](internal/release/workflow_mutation_test.go#L58)；且脚本契约会严格比对 shell 字段，见 [internal/release/workflow_test.go](internal/release/workflow_test.go#L295)。
- 文档层仍有准则描述滞后（非测试断言红线）：见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L675)，并已在报告 K4 与候选 OB-145 中披露，见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L64)。

## Q4
VERDICT: OK

依据：
- 代裁“待用户追认、不得作先例”有明确留痕：见 [docs/validation/evidence/P90-GATES-CI-PERF/00-briefing.md](docs/validation/evidence/P90-GATES-CI-PERF/00-briefing.md#L194)、[docs/validation/evidence/P90-GATES-CI-PERF/01-design.md](docs/validation/evidence/P90-GATES-CI-PERF/01-design.md#L10)、[docs/validation/evidence/P90-GATES-CI-PERF/01-design.md](docs/validation/evidence/P90-GATES-CI-PERF/01-design.md#L88)、[docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md](docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md#L9)。
- 台账 P90 段同样明确“主控代裁、待追认、不得先例”，见 [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2427)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L2414)。
- 未发现把“代裁”写成“已获授权”的直接表述；“用户裁定 f”表述用于 CI workflow 改动授权面，不是把代裁本身写成已追认，见 [docs/validation/evidence/P90-GATES-CI-PERF/01-design.md](docs/validation/evidence/P90-GATES-CI-PERF/01-design.md#L88)。

## Q5
VERDICT: DEFECT

依据：
- A/B 数字与报告 §4 一致：见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L35)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L38) 对照 [tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L4)、[tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L5)、[tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L11)、[tmp/P90-GATES-CI-PERF/ab-report.json](tmp/P90-GATES-CI-PERF/ab-report.json#L14)。
- HEAD^1..HEAD 披露到位：见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L37)。
- 比对器可证伪性成立：见 [tmp/P90-GATES-CI-PERF/comparator-sanity.json](tmp/P90-GATES-CI-PERF/comparator-sanity.json#L2)、[tmp/P90-GATES-CI-PERF/comparator-sanity.json](tmp/P90-GATES-CI-PERF/comparator-sanity.json#L3)。
- 但证据链三方不相等：MANIFEST 声称“三方相等”，见 [docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md](docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md#L5)；同时目录中存在 [docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md](docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md#L5)，而 MANIFEST 和 SHA256SUMS 仅登记到 06-scan-r1（见 [docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md](docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md#L24)、[docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt](docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt#L13)）。  
  只读脚本核验结果：disk=16、manifest=15、sha=15，缺口唯一为 04-reviews/06-scan-r2.md，且已登记文件哈希全匹配。

## Q6
VERDICT: RISK

依据：
- 报告对的镜像结构本身是成立的：index 读数给出 CN/EN +22/-0 和报告对 +128/-0，且 lines/head/bold/bars/blank/backticks 全部为 true，见 [tmp/P90-GATES-CI-PERF/gate-index-4b.txt](tmp/P90-GATES-CI-PERF/gate-index-4b.txt#L10)。
- 台账 P90 段状态与磁盘一致（执行中）：见 [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2414)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2419)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L2401)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md#L2406)。
- 候选 OB-145 的“本片不改准则”处置链完整：见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L64) 与 [docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md#L2427)。
- 风险点：K1-K5 未覆盖当前审出的问题（Q5 的证据链三方不等），因此“已知限制全覆盖”目前不充分。

## Q7
VERDICT: RISK

依据：
- CI 双腿读数未取被如实标注，且未用本机读数冒充 CI：见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L63)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L74)、[docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L75)、[docs/validation/gates-ci-perf_EN.md](docs/validation/gates-ci-perf_EN.md#L63)、[docs/validation/gates-ci-perf_EN.md](docs/validation/gates-ci-perf_EN.md#L74)、[docs/validation/gates-ci-perf_EN.md](docs/validation/gates-ci-perf_EN.md#L75)。
- 但“本轮未发现新缺陷”与当前终审发现不一致，见 [docs/validation/gates-ci-perf.md](docs/validation/gates-ci-perf.md#L91)、[docs/validation/gates-ci-perf_EN.md](docs/validation/gates-ci-perf_EN.md#L91)。因此“唯一未闭合项仅 CI 双腿读数”在当前状态下不成立。

| # | SEVERITY | file:line | 被违反的句子 | 最小修复 |
|---|---|---|---|---|
| 1 | Medium | [docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md:5](docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md#L5)；[docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md:24](docs/validation/evidence/P90-GATES-CI-PERF/MANIFEST.md#L24)；[docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt:13](docs/validation/evidence/P90-GATES-CI-PERF/SHA256SUMS.txt#L13)；[docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md:5](docs/validation/evidence/P90-GATES-CI-PERF/04-reviews/06-scan-r2.md#L5) | 条目集三方相等（磁盘 / 本清单 / SHA256SUMS.txt） | 将 04-reviews/06-scan-r2.md 同步加入 MANIFEST 与 SHA256SUMS，重算并回写计数口径（15 改 16），再重跑三方一致性核验并更新报告对应句子。 |
| 2 | Low | [docs/validation/gates-ci-perf.md:91](docs/validation/gates-ci-perf.md#L91)；[docs/validation/gates-ci-perf_EN.md:91](docs/validation/gates-ci-perf_EN.md#L91) | 本轮未发现新缺陷 | 在修复 #1 前，将该句改为“仍有未闭合缺陷（证据链三方不等）”；修复后再恢复“无新缺陷”表述。 |

FINAL REVIEW: FINDINGS

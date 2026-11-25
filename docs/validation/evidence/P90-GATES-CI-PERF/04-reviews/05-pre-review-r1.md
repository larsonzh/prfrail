# ⑤ 预审 · 第 1 轮（DeepSeek V4 Pro，只读载体 `deep-reasoner`）

- 日期：2026-11-25
- 输入：`01-design.md`、`impl.diff`、`numstat.txt`、`gate-tree-4.txt`、`ab-*` 证据、改动后 `ci.yml` / `workflow_test.go` / `workflow_mutation_test.go` / 台账 P90 段
- 结论：`PRE-REVIEW: PASS WITH FIXES`（1 Medium + 2 Low）

---

## Q1 覆盖等价

`VERDICT: OK`

`tools/gates/gate.js:93-125`（`runBaseGates`）四件与 CI 对等步逐字核对：

- gofmt：gate 条件 `gf.status === 0 && gfOut === ''`（gate.js:111）与 Format 步（ci.yml:92-101）`set -e` + `[ -n "$files" ] → exit 1` 逐分支等价：① 退出 0+无输出→双 PASS；② 退出 0+有输出→gate FAIL / 步 exit 1 红；③ 非零+无输出→gate FAIL（status≠0）/ `set -e` 在赋值语句上取命令替换的退出码（bash 语义）直接红；④ 非零+有输出→同③。gofmt 缺失时 gate 走 UsageError exit 2、步走 127，均为红，无「逃绿」分支。命令字面 `gofmt -l .` 与工作目录（仓库根）逐字相同。
- build/vet：CI `Build`/`Vet` 步（ci.yml:76-81）与 base 第 2/3 件命令逐字相同。
- test：CI `Test` 步 `go test -count=1 ./...`（ci.yml:83-84）**严格强于** base 第 4 件 `go test ./...`（gate.js:121，无 `-count=1`，可被结果缓存供绿）。

## Q2 不掩盖失败

`VERDICT: OK`

磁盘步序：Checkout source→Setup Go→Build→Vet→Test→Race(Linux)→Contract fixtures(Linux)→Format→Gates（ci.yml:55-110，Gates 确为第 9 位）。test job 内**无任何** `continue-on-error`（全文核对），任一步非零即 job 红且 Gates 不执行 ⇒ base 型失败必红、且红得更早。第 6–8 步均不写 `.go` 源（Race 只跑测试、Contract fixtures 只跑 npm、Format 只读），不存在「早步通过、base 才红」的确定性窗口；旧实现 Gates 本就在末尾、判据按 ctx.js 契约只读，无新增漂移源。不变量「Test 与 Gates 之间不得插入改源步骤」已在方案书 §I R1 显式登记为已知边界，可接受。

## Q3 判据面不变

`VERDICT: OK`（附一条 Low 观测）

- `ab-report.json`：`strippedBaseLines=4`、`baseLinesInB=0`、`aPrimeSha256 == bSha256 == ee00e354…`、`identical: true`、`aTotalFail == bTotalFail == "TOTAL_FAIL=0"`、A/B 头部 `changed=6` 相等。逐字读 `ab-A-all.txt`/`ab-B-skipbase.txt` 复核：除 4 条 `PASS base …` 行外全 46 行逐字一致。
- `G4b-cache` 两文件均为 `hit 80b2bbc3953c corpus=1472`（同一 ci 作用域键，已被预热，无 miss/hit 不对称）。
- `comparator-sanity.json`：`baselineIdentical=true`、翻转 1 字节后 `afterOneByteFlipIdentical=false`（sha `c0aa4903…`）⇒ 比对器可证伪。
- `out.fails()` 口径（`lib/report.js`：FAIL 记录计数）与「构成集合收缩为纯判据 FAIL」已在方案书 §E 显式披露，要求写入 ⑥/⑦ 输入 ✓。
- **Low 观测**：A/B 的 scope=ci 变更集 = `HEAD^1..HEAD` = v73 快照提交 `6e097fd` 的 6 文件 diff（`git log -1` 实证），并非 P90 未提交的 5 文件集（`git diff --numstat` 实证 15/1+22/0+22/0+1/1+5/2）。等价性是旗标机械性质、与变更集无关，故证明本身成立；P90 自身变更集的判据面由 ④ tree 读数（`gate-tree-4.txt`，`TOTAL_FAIL=1` 唯一 `G8-a`，已如实披露）与推送后 CI 复读覆盖。**建议报告写明 A/B 所用变更集**。

## Q4 冻结契约同步完备性

`VERDICT: OK`

全仓搜索实证：读取 `ci.yml` 的 Go 代码仅 `internal/release/workflow_test.go` 与 `workflow_mutation_test.go` 两处（`**/*.go` 全量 grep）；`tools/contracts/` 不解析 workflow 步骤（仅 README 声明运行期门禁在外部）；`TestWorkflowScriptsRejectCommentedCommands`（workflow_test.go:455）按 `Run != ""` 泛化遍历，自动覆盖新 Format 步；`validateWorkflowActions` 不需改（无新 Action、`test/Setup Go` 参数含 `cache: true` 冻结未变，注释不进 YAML 数据）。生活文档（DEV_PLAN/OPERATIONS/TEST_STRATEGY）无步骤名表/步数/摘要断言；历史 validation 报告与归档为冻结证据，不应同步。无遗漏必须同步点。

## Q5 `shell: bash` 可行性

`VERDICT: OK`

GitHub 文档化行为：非 Windows 默认 bash；Windows 显式 `shell: bash` 路由至镜像预装 Git for Windows bash（文档支持），`gofmt` 经 setup-go 在 PATH；唯一参数 `.` 无 MSYS 路径转换风险。失败形态（若镜像移除）为步骤启动失败（无法定位 bash/gofmt），fail-closed 红腿、Gates 不执行——不构成逃逸。已登记为推送后双腿验证点：台账验收证据 ⑥「CI 两腿读数于推送后回填」+ 方案书 §C 回退预案（证伪则本片回滚另议，不降级为近似形态）✓。

## Q6 Go 缓存结论

`VERDICT: OK`（附一条 Low 措辞警示）

「已覆盖不动」成立：setup-go `cache: true` 覆盖模块缓存与构建产物缓存两类；diff 只加注记、`cache: true` 行零改动（numstat 15/1 逐行核对）；`-count=1` 仅禁用测试结果复用（有意设计，非缺口）；`candidate-build`/`candidate-probe`/`selfhost` 三信任 job 的 `cache: false` 是**刻意选择**（`docs/validation/s1-selfhost.md:9` 明载「不通过 Go 缓存恢复或保存跨 runner 的构建状态」）⇒ **不应改**。**Low 警示**：ci.yml:73 注释与方案书 §F 的「go.mod-keyed」键派生措辞与 setup-go README 通行描述（模块缓存键通常由 `go.sum` 哈希派生，`cache-dependency-path` 默认 `**/go.sum`；构建缓存键由平台/arch/Go 版本派生）不符，结论不受影响（键过期只会补下载、内容寻址不会错用），但报告定稿前应复核该措辞。

## Q7 ADR 判定

`VERDICT: OK`

不立 ADR-049 成立：`docs/ADR_REGISTER.md` 现存 ADR-040–048 均为模型选型/裁定类，**无任何 CI workflow 扩展 ADR**——先例 `gates-ext` 的 `fetch-depth: 0` 扩展确以 ci.yml:61-63 注释单一留痕、无 ADR（归档证据 `directive-history/gates-ext/` 佐证）；ADR-049 编号未被消费。三元留痕已落其二：ci.yml:104-105 注释（含「user-authorized extension of the frozen CI contract」）+ 台账 P90 边界扩展句；报告对 §6/freedom-list/⑩ 回执为后续步骤（切片执行中，属正常时序）。

## Q8 边界与代裁

`VERDICT: RISK`　`SEVERITY: Medium`

- 代裁本身成立且标注到位：台账 P90 段含「**边界扩展（主控代裁，待用户追认，不得作先例）**」逐字在案（REMAINING_SLICES.md:2433 行域磁盘核实）；不改 `workflow_test.go` 则 CI 双腿必红的依据真实（Gates run 串被 sha256 冻结断言）。
- mutation `bash→sh`（workflow_mutation_test.go:58）**是必要适配、非弱化**：旧值 `bash` 作用于新 Format 步时为 no-op → `reflect.DeepEqual` 命中 → `mutation did not change the workflow` fatal（假红）；改 `sh` 后与全部既有 shell（`""`/`pwsh`/`bash`）均构成偏离，所有步骤的 shell 篡改仍被拒收，断言强度不减。
- **不一致点（规则②违规面）**：台账边界扩展句声明「仅摘要与步骤名表，`internal/**` **其余零改动**」，方案书 §C「`workflow_mutation_test.go` …无需改」、§J「`workflow_mutation_test.go`、…及全部测试**零改动**」——与 numstat 实证（`workflow_mutation_test.go 1/1`）直接矛盾。方案书在该点预测错误。
- **最小修复**：① 台账边界扩展句改为「仅摘要与步骤名表 + `workflow_mutation_test.go` 的 different-shell 变异值必要适配 1 行（bash→sh，否则对 Format 步为 no-op 假红）」；② 方案书 §C/§J 同步补正并撤销「无需改」表述；③ 该适配写入 freedom-list 与 ⑥/⑦ 输入（如实声明「方案书 §J 预测偏差 + 实际 1 行必要适配」）；④ 无代码回滚需求（`go test ./internal/release/` 已绿）。

## Q9 遗漏面

`VERDICT: OK`（两条 Low 注记）

- 报告对 `docs/validation/gates-ci-perf{,_EN}.md`、freedom-list、证据包 MANIFEST/SHA256SUMS、`evidence-forms.txt` 登记：均未落盘——⑤ 阶段属正常时序；`gate-tree-4.txt` 唯一 FAIL `G8-a`（证据树未登记）正是该状态的如实披露，闭环条件在台账验收证据 ④/⑤，须于 ⑩ 前完成。
- `cjk-newwords.txt`：无需登记（tree-4 与 A/B 两轮 `G4b unknown=[]`，tools/** 零改动与边界一致）✓。
- 描述 CI 步骤序的文档：无生活文档需同步（Q4 已证）✓。
- **Low 1**：`01-design.md:47` 的 Format 摘要为草稿值 `ff770e4b…`，与实际落盘 `5a9a2b57…`（workflow_test.go:250，`go test` 权威复核已绿）不符——证据包登记前须更正或加注「以落盘测试值为准」。
- **Low 2**：`tmp/P90-GATES-CI-PERF/review/impl.diff` 非 ASCII 段呈 GBK 乱码（PowerShell 管道产物，与方案书 §H「node 直写规避乱码」的自定口径相悖）——该件不入库（tmp/），但 ⑥ 输入 diff 应改用 UTF-8 安全通道重生成。

`PRE-REVIEW: PASS WITH FIXES`

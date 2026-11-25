## A 影响判定（逐条）

| 项 | 判定 | 依据 |
|---|---|---|
| CONTRACTS 正文 | **不触** | 全仓 grep：CONTRACTS.md 无 CI workflow/`ci.yml` 契约段；「frozen CI contract」的载体是 `internal/release/workflow_test.go`（Go 测试），不在 CONTRACTS 语义域。本片不改 `docs/CONTRACTS{,_EN}.md` 任何字节。 |
| schema | **不触** | `schemas/**` 零改动。 |
| 契约夹具 | **不触** | `testdata/**`（含 selfhost 五 JSON + 链夹具）与 `tools/contracts/**` 零改动；Contract fixtures 步本身及其 run 串不动。 |
| 准则正文 | **不触** | 不升版；`DELIVERY_DIRECTIVE{,_EN}.md` 零改动。`--skip-base` 语义由 OB-113（v1.36）早已确立，本片只改 CI 用法，不新增准则条款。 |
| §7.2.1 ③ 语义（状态机/门禁判据/证据模型/角色授权） | **不触** | 判据本体与判定逻辑零改动（`gate.js`/`selftest.js`/`ctx.js`/`lib/criteria/**` 零改动）；不涉状态机/证据模型/角色/授权。需如实声明的唯一面：§6.1 基础门禁在 CI 的**执行编排**变化（位置移动 + gofmt 步化），等价性由 §D 论证；`.github/workflows/**` 属 §1.6 增量边界，已获用户裁定 f 授权（起跑消息「f. CI workflow 改动授权」）。 |
| 冻结契约实体（`workflow_test.go`） | **触（执行纪律面）** | 非 CONTRACTS/schema/夹具；它是仓库级硬规则要求的强制伴随同步（§J）。已由主控代裁登记入台账 P90 段「边界扩展，待用户追认，不得作先例」。 |

## B §7.2 ⑥ 判定

按 §7.2.1 三条件逐条：

1. **① 不满足**：改动含 `.github/workflows/ci.yml`（机器结构 YAML），不属于「只落文档与配置文案类文件」的列举（`.md`/注释/文案/`.gitignore`/`.txt`/`.html`）。
2. **② 不满足**：改 `.yml`（步骤与脚本语义）+ 改 `.go`（`workflow_test.go` 契约断言），双项都在 ② 的排除清单内。
3. **③ 单独成立**：判据判定逻辑零改动、不涉状态机/证据模型/角色授权——但三条须**同时**满足。

**结论：⑥ 必做，不可跳过**（MAI Code 2 Flash，≤3，按起跑消息二节额度）。反面样例提醒（SW-1/SW-2）：不得按「看起来像配置/工具片」推断跳过。**修正项**：起跑消息二节「⑥ 预计跳过」以本 ① 判定为准改为「⑥ 必做」，写入 freedom-list 与报告，⑥ 输入包须含 `ci.yml` diff + `workflow_test.go` diff + 本方案书。

## C gofmt 处置

**裁定 a）新增独立格式步**。理由：(b) 仅靠本机纪律 + ADR 使 PR/push 路径失去 gofmt 自动防线，相对现状是倒退；(c) 其它等价方案（并入 Build 步 / `go fmt` / gofumpt）改变命令字面，失去与 gate.js `gofmt -l .` 的逐字对应，等价论证负担反而加重。裁定与台账 P90 段目标 ② 一致。

**精确 YAML**（插入位置 = test job 的 Gates 步正上方；对既有步字段零改动）：

```yaml
      - name: Format
        shell: bash
        run: |
          set -e
          files="$(gofmt -l .)"
          if [ -n "$files" ]; then
            printf '%s\n' "$files"
            exit 1
          fi

      - name: Gates
        run: node tools/gates/gate.js --selftest && node tools/gates/gate.js --all --skip-base --scope=ci
```

`set -e` 使 gofmt 自身非零退出也判红，与 gate.js 的 `gf.status !== 0 || gfOut !== ''` 双条件严格等价（非空输出 = 旧 FAIL 分支；仅非零无输出 = `set -e` 兜底）。

**冻结契约改动（按锚点，行号随插入偏移）**：
- `validateWorkflowStructure` 的 `expectedJobs["test"].steps`：插入 `"Format",` 于 `"Gates"` 之前 → 9 项：`{"Checkout source", "Setup Go", "Build", "Vet", "Test", "Race", "Contract fixtures", "Format", "Gates"}`（步数与步名两断言自动通过）。
- `expectedScripts` 表：新增一行 `"test/Format": {digest: "ff770e4b1c46d3ca1bcb29a3e85e930b9d69734f6e943fcbe485d2ad836bc478", shell: "bash"}`；`"test/Gates"` 的 digest 改为 `53f7fed093a67fae3c5af018a308509dd96120491f6362501abfed818a99fbd5`。
- 其余**零改动**：`validateWorkflowActions` 不变（无新 Action、`test/Setup Go` 参数 `{"go-version": ..., "cache": true}` 不动）；`TestWorkflowScriptsRejectCommentedCommands` 按 `Run != ""` 泛化遍历，自动覆盖新 Format 步；`workflow_mutation_test.go` 的 extra-step/missing-step/reordering 变异泛化，无需改。

**摘要实证**（本机已复算）：旧 Gates 串 = `ca20623b…c8d0`（与冻结值逐字相符）；新串 = `53f7fed0…fbd5`（与起跑消息给定值相符）；Format 串摘要按 YAML 字面块标量含**末尾换行**的口径计算（用冻结的 `test/Contract fixtures` 摘要 `94577fb7…` 反向验证了该口径，node 复算逐字吻合）。② 落盘后以 `go test -count=1 ./internal/release/` 为唯一权威复核（不匹配时该测试打印实算值）。

**两侧 runner 可行性**：
- ubuntu-latest：`shell: bash` = 原生 bash；gofmt 经 setup-go 已在 PATH。
- windows-latest：GitHub Actions 把 `shell: bash` 路由到镜像预装的 Git Bash（`C:\Program Files\Git\bin\bash.exe`），属文档化支持；gofmt 同在 PATH（windows 腿 Gates 步经 Node spawn 跑 gofmt 已有 P84–P89 多轮双绿实证）。MSYS 路径转换对本步无影响（唯一参数 `.`）；gofmt 输出的反斜杠路径仅影响日志观感，不影响 `-n` 判空。
- 本仓库无 `shell: bash` 先例 ⇒ 列为 ③ 推送后双腿验证点（台账验收证据 ⑥）。若 windows 腿证伪，回退 = 本片回滚另议，不回退为近似形态。

## D base 覆盖等价论证（逐件）

| base 件 | CI 对等步 | 论证 |
|---|---|---|
| gofmt | 新增 Format 步 | 同一命令 `gofmt -l .`、同一工作目录（仓库根）、同一工具链；判定条件 `非零退出 ∨ 非空输出` 与 gate.js 双条件严格等价（§C）。 |
| build | Build 步 `go build ./...` | 与 base 第 2 件**逐字相同**的命令与工作目录。 |
| vet | Vet 步 `go vet ./...` | 与 base 第 3 件逐字相同。 |
| test | Test 步 `go test -count=1 ./...` | **严格更强**：base 第 4 件用 `go test ./...`（可被 Go 测试缓存供结果），CI Test 的 `-count=1` 绕过结果缓存、每次真跑；包集同为 `./...`。 |

**去重不掩盖失败**：Gates 步排在第 8 位，Build/Vet/Test/Format 均先于它且均无 `continue-on-error`（默认 fail-closed），任一步非零退出即红并终止 job，Gates 根本不会执行。base 型失败只能由更早步骤表达——而早先旧实现中 Gates 的 base 重复跑（同树、同命令、同工具链）与这些步骤同输入同输出，属确定性冗余，不存在「早步通过、base 才红」的差异窗口（Test 与 Gates 之间仅 Race/Contract fixtures/Format 三步，均不改 `.go` 源）。故去重后任何 base 类失败仍必红、且红得更早。

## E 判据面不变

- **消失的行**：恰好 4 条 `base` 记录（`PASS/FAIL base gofmt -l . :: ...`、`go build ./... :: exit 0`、`go vet ./... :: exit 0`、`go test ./... :: ok=N FAIL=0`），无其它行受影响。
- **必须逐字不变的行**：头部 `GATE REPORT scope=ci repo=… changed=N`（`changed=` 只取决于 HEAD^1 diff，与旗标无关）；全部 18 条判据记录（G1-a…G9）；`TOTAL_FAIL=0`；以及 Gates 步前半段 `--selftest`（SELFTEST 744/744 不变）。
- **`out.fails()` 口径**（`lib/report.js`：`records.filter(r => r.verdict === 'FAIL')`）：TOTAL_FAIL 语义本身不变（FAIL 记录计数）；但 `--skip-base` 后不再追加 base 记录 ⇒ **TOTAL_FAIL 构成集合收缩为「纯判据 FAIL 数」**，理论上不再可能含 ≤4 条 base FAIL。在 CI 上该收缩不可观测（base 型失败已由更早步骤把 job 红掉，Gates 不产出报告）；判据逐条的 verdict 计算逻辑零改动。此口径收缩须在报告与 ⑥/⑦ 输入中显式声明，防止被误读为「判据宽松化」。

## F Go 缓存评估

**结论：三选一之「已覆盖不动」**。依据（setup-go 官方 README「Caching dependency files and build outputs」段，本次实读）：

- `cache: true` 覆盖两类：① **Go 模块缓存**（GOMODCACHE），键派生自 `go.mod` 哈希（v6 起默认；v5 用 go.sum）+ OS/arch/Go 版本；模块缓存内容寻址，键过期只会补下载、不会错用。② **Go 构建产物缓存**（GOCACHE，Linux `~/.cache/go-build` / Windows `%LocalAppData%\go-build`），独立键（平台/arch/Go 版本），每轮恢复+保存、增量合并。v7 中 `cache` 默认即为 true，显式写 `cache: true` 与冻结参数一致。
- 对 base 四件套的作用：build 全量受益（依赖 + 编译产物双缓存）；vet 受益于编译缓存（vet 分析本身不缓存但重编译省掉）；test 的**编译阶段**受益，**结果复用被 `-count=1` 有意禁用**（属刻意设计，不是缺口）；gofmt 无缓存可言（纯 I/O，P88 实测 1.1 s）。
- 佐证：`docs/validation/s1-selfhost.md` 明载三个信任 job 显式 `cache: false`「不通过 Go 缓存恢复或保存跨 runner 的构建状态」——即 test job 的 `cache: true` 是既有刻意选择，勿动。
- 缺口注记（不构成本片改动）：`Contract fixtures` 的 `npm ci` 无任何缓存（非 Go 面，与 base 四件套无关，列后续片建议）。
- 因结论非「仅注记」，**不加 YAML 注释**（零改动优先）。若 ② 坚持注记，落点 = Setup Go 步 `cache: true` 行上方，文案（不含会过期字面量）：`# setup-go restores the module cache (go.mod-keyed) and the Go build cache, warming Build/Vet/Test; keep cache enabled.`

## G ADR 判定

**代裁成立，ADR-049 不必要**。

- 代裁成立理由：① 目标（改 Gates run 串）与仓库硬规则（`workflow_test.go` 对 run 串做 sha256 断言）强制绑定，不改测试则 `go test ./internal/release/` 与 CI 双腿必红，不存在合规的「只改 ci.yml」路径；② 代裁范围最小化（仅摘要值 + 步骤名表 + 一个新脚本条目，`internal/**` 其余零改动），并标注「不得作先例 + 待用户追认」；③ §1.7.3 判定：该决定错了要改的是「代码（测试）」，属执行自由度的伴随动作而非契约/方向/授权变更——授权本身已由用户裁定 f 给出。
- 不需 ADR 的理由：本片无新决策需要 ADR 记录——去重方向 = P88 报告 §7 建议 2（用户点名），gofmt 处置 = 台账目标 ② 已定，缓存结论 = 现状评估。**留痕三元组**（与先例一致）：(1) ci.yml Gates 步上方注释（参照 :63 句式）：`# Skip the base gates here: Build/Vet/Test/Format above already cover section 6.1 (user-authorized extension of the frozen CI contract).`；(2) 台账 P90 段边界扩展句（已落）；(3) 报告对 §6 + freedom-list F 项 + ⑩ 回执显式点名待追认。
- 先例：`gates-ext`（新增 Gates 步 + `fetch-depth: 0` 冻结契约扩展）仅 ci.yml 注释留痕、无 ADR；`RACE-LEG-EXPAND`（ci.yml 1 行 + 契约摘要 1 行）无 ADR；T017 `6024526` 创建契约本体亦无 ADR。ADR-049 编号保留不消费；若用户 ⑩ 追认时要求强化留痕，再补立。

## H ③ 一致性验证方案（本机，同一变更集同一作用域）

1. **环境锁定**：跑 A 前记录 `git rev-parse HEAD` 与 `git status --porcelain`，A/B 之间与跑后复读，三方必须逐字相同（禁止任何 git 写、暂存、编辑）。
2. **A/B 执行**（node 驱动直写，规避 PowerShell 管道 GBK 乱码，沿 P88 先例）：
   ```
   node tools/gates/gate.js --all --scope=ci            → A.txt（含 base 段）
   node tools/gates/gate.js --all --skip-base --scope=ci → B.txt
   ```
3. **逐字比对**（node 单行脚本）：读 A，剔除匹配 `/^[A-Z]+\s+base\s/` 的行得 A′；断言 B 中 `base` 记录行数 = 0；`A′ === B` 逐字节相等；输出两枚 sha256 与 `IDENTICAL`/`DIFFERENT`。
4. **比对器可证伪性**：把 B 翻转 1 字节后重跑，必须报 `DIFFERENT`（同 P88 `comparator-sanity.txt`）。
5. **附加断言**：两文件 `TOTAL_FAIL=` 行相等且为 0；头部 `changed=N` 相等。
6. **CI 侧口径（区分声明）**：跨 commit 的 CI 输出 `changed=N` 不同，**不**做跨 run 逐字比对；CI 验收 = 双腿 Gates 步 `SELFTEST: PASS 744/744` + `TOTAL_FAIL=0` + 形态完整（台账验收证据 ①）。

## I 风险与反例（可证伪）

1. **R1 覆盖等价依赖「Test 之后无步骤改 .go」的不变量**——若未来有人在 Test 与 Gates 之间插入会改 `.go` 的步骤，该改动将无人复核。反例：插入改 `.go` 的探针步后，旧实现 Gates base 红而新实现绿（本片报告须把该不变量显式登记为已知边界）。
2. **R2 gofmt 等价性依赖 Format 步与 gate.js 逐字同构**——若 ② 把 `files="$(gofmt -l .)"` 换成近似形态（如管道 grep），「gofmt 非零退出且无输出」分支的行为会偏离。反例：用假 gofmt（退出 1、无输出）探针，两侧判定不一致；契约摘要断言（R7）与 ⑥ 输入 diff 可拦。
3. **R3 `shell: bash` 在 windows-latest 的可用性假设**——若 Windows 镜像移除 Git Bash 或 PATH 不含 gofmt，windows 腿 Format 步红。反例：推送后 windows 腿 Format 步 failure（本机不可模拟，故列入 ③ 双腿验证点）。
4. **R4 缓存结论依赖 setup-go v7 的 build-cache 行为**——若行为回归（只缓存模块不缓存构建产物），去重后的总耗时收益低于 P88 预估（141.8 s 基线）。反例：改动后 run 的 Test/Build 步耗时相对改动前 run 显著上升。
5. **R5 测试件「更强」论证依赖 CI Test 保持 `-count=1`**——若被改回 `go test ./...`，缓存通过 + Gates 不复核构成双失守。反例：`test/Test` 摘要变化（契约测试会红，风险已被冻结契约覆盖，属低但列明）。
6. **R6 TOTAL_FAIL 口径收缩被误读为判据宽松化**——若 ⑥/⑦ 或用户据「不再含 base」质疑判据面。反例：A/B 比对中任一条判据行不一致（③ 机械比对即证伪面）。
7. **R7 冻结契约同步遗漏**（步骤名表 / Gates 摘要 / Format 条目任一漏改或摘要多算/少算换行）。反例：`go test -count=1 ./internal/release/` 非零退出（③ 必跑项，机械可核）。

## J 边界与代裁核对

**必须不改（文件级清单）**：
- `tools/gates/**` 全部（`gate.js`/`selftest.js`/`ctx.js`/`lib/**` 零改动，SELFTEST 744/744 字面不变）。
- `internal/**` 除 `internal/release/workflow_test.go` 三处最小改动（test 步骤名表 1 行、`test/Gates` 摘要 1 行、新增 `test/Format` 条目 1 行）；`workflow_mutation_test.go`、`probe.go`、`verify.go`、`generate.go` 及全部测试零改动。
- `.github/workflows/ci.yml` 除：Gates run 行 1 行、新增 Format 步、可选注释。`candidate-build`/`candidate-probe`/`selfhost` 三 job 与 test job 的 Checkout/Setup Go/Build/Vet/Test/Race/Contract fixtures 各步字段零改动（含 `fetch-depth: 0` 与 `cache: true` 冻结参数）。
- `docs/CONTRACTS{,_EN}.md`、`schemas/**`、`testdata/**`、`tools/contracts/**`、`harnesses/**`、`examples/**`。
- `docs/DELIVERY_DIRECTIVE{,_EN}.md`、`docs/t027/CHECKLIST_v73{,_EN}.md`（v73 快照冻结）；台账改动仅 `docs/t027/REMAINING_SLICES{,_EN}.md` 的 P90 段（步骤 0 已落）。
- sessbridge 仓库（多根工作区另一仓）。

**代裁核对结论：成立且无遗漏、无多余**。全仓实证：冻结 Gates run 串的摘要 `ca20623b…` 仅存在于 `workflow_test.go:250` 一处（其余命中为归档证据与简报原文）；读取 `ci.yml` 的 Go 测试仅 `workflow_test.go` 与 `workflow_mutation_test.go` 两个，后者无步骤名表、变异泛化自动适配新步；`TestWorkflowScriptsRejectCommentedCommands` 自动覆盖 Format；`validateWorkflowActions` 无需改（无新 Action、参数冻结不变）。代裁列举的「仅摘要 + 步骤名表 + 必要脚本契约条目」恰好等于本片实际需要的全部契约同步面。

DESIGN: DONE

## Corrections after ⑤ round 1

The sections above are the architect output as delivered; they are not rewritten. Three predictions were corrected during implementation, each with the observable evidence that refuted it.

| # | Predicted in the design | Actual (disk truth) | Evidence |
|---|---|---|---|
| C1 | §C: the `Format` step digest `ff770e4b…` | `5a9a2b5725f2311225eaa69c3be661ca5caa8602563a1158d693057af3777471` | `go test -count=1 ./internal/release/` printed the computed value; `workflow_test.go` now pins it |
| C2 | §C/§J: `workflow_mutation_test.go` needs no change | one necessary adaptation: `different-shell` moves from `bash` to `sh`, because `bash` is a no-op for the new `Format` step (expected shell is `bash`) and trips `mutation did not change the workflow` | `go test` failure text before the edit; `git diff --numstat` shows 1/1 |
| C3 | §F: module cache key described as "go.mod-keyed" | setup-go keys the module cache from `go.sum` by default (build cache keyed by platform, arch and Go version); the conclusion "already covered, unchanged" is unaffected | the ci.yml comment was reworded to avoid the inaccurate key description |

Additional disclosure requested by ⑤ (Q3): the ③ A/B pair ran with `--scope=ci`, which locally resolves to `HEAD^1..HEAD`, i.e. the change set of commit `6e097fd` (six files), not the uncommitted P90 change set. The equivalence under test is a property of the flag, so the proof stands; the P90 change set itself is covered by the ④ `--scope=tree` reading and by the post-push CI legs.

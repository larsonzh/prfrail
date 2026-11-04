# 03-test-report — P55 `EXIT-SEED-BUILD-EVIDENCE` ③ 校验测试

| 字段 | 内容 |
|---|---|
| 切片 | P55 `EXIT-SEED-BUILD-EVIDENCE`（AT-13 seed 驱动构建补证；层 3 第 7 项） |
| 环节 | ③ 校验测试（证据链审计；**不重跑构建、不改产品代码、不改 CI**） |
| 锚点 | `HEAD` = `baa144df5898817e0ac7cd860928bb4a90cbd009`（`baa144d`，实测 `git rev-parse HEAD` 一致） |
| 复核时点 | 2026-11-05 |
| 包内文件 | 46 件（44 件载荷 + `MANIFEST.md` + `SHA256SUMS.txt`） |
| 审计性质 | 只读。未修改任何既有证据文件、`SHA256SUMS.txt`、`MANIFEST.md`、台账、产品代码或 `.github/workflows/**` |
| 本环节新增件 | `tmp/EXIT-SEED-BUILD-EVIDENCE/verify-sums.js`、`tmp/EXIT-SEED-BUILD-EVIDENCE/verify-cross.js`、本文件 |

> **范围声明**：本报告只判定「包内读数能否复算、断言能否被证伪、双侧能否交叉对齐」，不重新执行 seed 驱动构建。
> 因此所有「PASS」均指**证据链自洽且可复算**，不指「本次现场重跑通过」。

## 一、复核范围与方法

| 编号 | 任务 | 脚本 / 命令 | 退出码 | 结论 |
|---|---|---|---|---|
| A | 三方集合相等（`MANIFEST` ≡ `SHA256SUMS` ≡ 磁盘载荷） | `node tmp\EXIT-SEED-BUILD-EVIDENCE\verify-sums.js` | 见下 | **相等** |
| B | 全量 sha256 重算（逐条比对 `SHA256SUMS.txt`） | 同上，B 段 | 见下 | **44/44 命中，0 失配** |
| C | `index.json` 元数据自洽（`sha256` / `rawBytes` / `EXIT`） | 同上，C 段 | 见下 | **23 条不一致（均为元数据层，载荷层无影响）** |
| D | 编码 / 行尾（`.md` = BOM+LF；`.txt`/`.json` = 无 BOM+LF） | 同上，D 段 | 见下 | **0 违规** |
| E | 双侧交叉核对（夹具 / oracle / seed / candidate / 二进制 / 提交） | `node tmp\EXIT-SEED-BUILD-EVIDENCE\verify-cross.js` | **0** | 全部布尔断言 TRUE |
| F | 失败路径逐条断言（退出码 / stderr / stdout / seed 完好） | 同上，T3 段 | **0** | 8 反例 + 2 完好断言全通过 |
| G | 独立外部核对（不依赖包内自述） | `git rev-parse HEAD`、`git cat-file -t`、`git ls-tree -r fc99f55`、`Get-FileHash testdata\selfhost\*.json` | 0 | 见 §三 与 §六 |

`verify-sums.js` 的退出码为 **1**，其唯一来源是 C 段（`index.json` 元数据不一致）；A / B / D 三段内部判定均为 OK。
即：**载荷哈希链无失配，`index.json` 辅助元数据存在滞后与字段语义分歧**（详见 F-03）。

复算口径：
- `SHA256SUMS.txt` 按 `<64 hex><两空格><相对路径>` 解析，共 44 行，重复条目 0。
- 载荷集 = 包内全部文件 − `{MANIFEST.md, SHA256SUMS.txt}` − `03-test-report.md`（本环节新增，见 §七）。
- `MANIFEST.md` 归档条目按 `| \`path\` |` 表格行机械抽取，共 44 条。

## 二、哈希复核结果

| 项 | 读数 |
|---|---|
| 包内文件总数 | 46 |
| 元数据件 | 2（`MANIFEST.md`、`SHA256SUMS.txt`） |
| 磁盘载荷件 | 44 |
| `SHA256SUMS.txt` 条目 | 44（重复 0） |
| `MANIFEST.md` 归档条目 | 44 |
| 磁盘载荷集 ≡ `SHA256SUMS` 集 | **true**（未列于磁盘 0 条 / 列出但缺失 0 条） |
| `MANIFEST` 集 ≡ `SHA256SUMS` 集 | **true**（`manifest_only` 0 / `sums_not_in_manifest` 0） |
| `MANIFEST` 条目全部存在于磁盘 | **true** |
| **sha256 重算命中** | **44 / 44** |
| **失配清单** | **（空）** |
| 编码 / 行尾违规 | **0**（46 件快照：`.md` 4 件均 BOM+LF；`.txt`/`.json` 42 件均无 BOM+LF；`\r\n` 计数 0。本件落盘后复跑仍为 **0**，`.md` 计 5 件） |

**结论**：载荷完整性与 `MANIFEST` 声明（「条目集三方相等」）**实测成立**。`MANIFEST.md` 与 `SHA256SUMS.txt` 未被列入 `SHA256SUMS.txt`，与 `MANIFEST.md` 的元数据声明一致。

### `index.json` 元数据不一致明细（C 段）

| 类别 | 件数 | 明细 | 判定 |
|---|---|---|---|
| `index.sha256` 与磁盘 / `SHA256SUMS` 均不符 | 3 | `02b-02-worktree-seed`（index `8b19b1b3…` / 磁盘 `7dadf79b…`）、`02b-03-worktree-verifier`（`4e2f2f09…` / `f1d5185e…`）、`02b-04-worktree-candidate`（`2b84a039…` / `2fe8e78a…`） | **滞后元数据**（磁盘与 `SHA256SUMS` 一致，载荷无损） |
| `rawBytes` 与存储文件字节数不符 | 10（全部 Windows） | `02a-01` 0/168、`02a-02` 0/269、`02a-03` 0/220、`02a-04` 72/251、`02a-05` 0/324、`02a-06` 339/664、`02a-07` 61/285、`02a-08` 64/292、`02a-09` 41/225、`02a-10` 74/327 | **字段语义分歧**：Windows 侧 `rawBytes` = stdout 正文长度；Linux 侧 `rawBytes` = 存储文件长度（`02b-12` 296/296、`02b-15` 636/636） |
| 捕获体无 `EXIT:` 字段 | 6 | `02a-11`、`02a-12`、`02b-19`、`02b-20`、`02b-22`（衍生件）+ `02a-00`（改用 `EXIT=` 分隔符） | 格式差异；其中 5 件的退出码仅存在于 `index.json` 自述（见 F-11） |

## 三、双侧交叉核对（Windows ↔ Linux）

| 核对项 | Windows | Linux | 判定 |
|---|---|---|---|
| `HEAD` / `SourceCommit` | `baa144d…`（`bindings.head`、`02a-06.SourceCommit`） | `baa144d…`（`02b-21.candidateCommit`、`02b-14.SourceCommit`） | **一致**（均 = 本地 `git rev-parse HEAD`） |
| `bootstrap.json` | `ae0d5db8…` | `ae0d5db8…` | **四面一致**（含 CI 输入表与本地磁盘重算） |
| `expected.json` | `388f8884…` | `388f8884…` | **四面一致** |
| `noop-chain.json` | `fbeee2b6…` | `fbeee2b6…` | **四面一致** |
| `executable-chain.json` | `dbc5fc41…` | `dbc5fc41…` | **四面一致** |
| `licenses-policy.json` | `ff9d243c…` | `ff9d243c…` | **四面一致** |
| `OracleHash`（selfhost 输出） | `sha256:388f8884…` | `sha256:388f8884…` | **一致**，且 = `02-win-captures/expected.json` 文件哈希 = `testdata/selfhost/expected.json` 本地重算值 |
| `SeedHash` | `sha256:a42f547c…` | `sha256:67ea456b…` | **按构造本就不同**（非失配）：Windows 为 `work/seed` 树哈希，Linux 为含 `bin/prfrail` 的组装根哈希 |
| `CandidateHash` | `sha256:d0e7c97b…` | `sha256:f689e7f7…` | **按构造本就不同**（非失配）：Windows 候选根仅 `actual.json`；Linux 候选根含 `actual.json` + `bin/prfrail` + `candidate-commit.txt` |
| 同侧 seed ≠ candidate | true | true | **满足「分离」的路径级要求** |
| `SeedHash` 与基线捕获自洽 | = `02a-04` 输出 | = `02b-12` 输出 | **一致** |
| 三角色提交 | `bindings.json` **无** `verifierCommit` 字段（仅 `head`） | seed `fc99f55…` / verifier `18c9340…` / candidate `baa144d…` | Linux **三者互异且均可在本地仓库解析**；Windows **无独立核验者绑定**（F-02） |
| 二进制哈希 | **未记录**（`work/bin/*` 未哈希，工作区已删） | 3 枚：`1e7f78ef…` / `4a1c0bff…` / `9b7f72e4…`（互异） | Linux 完备；Windows 缺口（F-08） |
| 工具链 | `go1.27.0 windows/amd64` | `go1.22.12 linux/amd64` | **两侧不一致**（已登记）；CI 固定 `GO_VERSION: "1.22.12"`（`ci.yml:46` 实测）⇒ **仅 Linux 侧与 CI 同版** |
| 隔离构建缓存 | `tmp\…\work\gocache`、`…\gomodcache`（`02a-00`） | `/tmp/p55work.oIFQtI/{gocache,gomodcache}`（`02b-00`） | **两侧均有隔离读数**；Linux 另有 `GOPROXY_build=off` 构造期无网络读数 |
| expected 载体位置 | `testdata/selfhost/expected.json`（仓库只读输入）+ 包内副本 `02-win-captures/expected.json` | `seed/fixtures/expected.json`（seed 根内）；`02-linux-captures/` **无** `02b-oracle-expected.*` 文件 | **命名与位置不对称**（用户提示的命名差异经实测确认）；oracle 权威值仍由 `02b-21.fixtures.expected.json` 与 `02b-14` 命令路径给出 |
| oracle 归属独立性 | 候选根无 `expected.json`（`02a-12`，`expectedAbsentInCandidate=true`） | 候选根无 `expected.json`（`02b-20`） | 一致；但「absent=true」为**自撰断言行**，非工具输出（F-10） |

**关于「两侧 seed/candidate 哈希不同」的定性**：这不是失配。两侧的 seed 根与候选根**构成不同**（Windows 候选根仅 1 个文件，Linux 候选根 3 个条目），且 `CandidateHash` 是对**候选产物包**做树哈希而非对候选源码做树哈希。因此**跨侧 `SeedHash`/`CandidateHash` 不可比**，报告与 ⑧ 文档同步时不得把两侧数值当作同一对象比对。

## 四、失败路径逐条断言结果

### 4.1 反例捕获（退出码 / stderr / stdout）

| 捕获 | 退出码 | 非 0 | stderr 命中预期串 | stdout 为空 | 判定 |
|---|---|---|---|---|---|
| `02a-07-neg-roots-overlap` | 1 | ✅ | ✅ `seed and candidate roots overlap` | ✅ | PASS |
| `02a-08-neg-oracle-ownership` | 1 | ✅ | ✅ `oracle ownership is not independent` | ✅ | PASS |
| `02a-09-neg-seed-changed` | 1 | ✅ | ✅ `seed changed` | ✅ | PASS |
| `02a-10-neg-oracle-drift` | 1 | ✅ | ✅ `candidate result differs from external oracle` | ✅ | PASS |
| `02b-15-neg-roots-overlap` | 1 | ✅ | ✅ `seed and candidate roots overlap` | ✅ | PASS |
| `02b-16-neg-oracle-ownership` | 1 | ✅ | ✅ `oracle ownership is not independent` | ✅ | PASS |
| `02b-17-neg-seed-changed` | 1 | ✅ | ✅ `seed changed` | ✅ | PASS |
| `02b-18-neg-oracle-drift` | 1 | ✅ | ✅ `candidate result differs from external oracle` | ✅ | PASS |

八条均为 fail-close（非 0 退出 + 预期错误串 + 无 stdout），且 `02b-22-verdicts.txt` 的逐件 `id / exit / verdict` 与本表逐条吻合。

### 4.2 失败后 seed 完好

| 捕获 | `seedHashBefore` = `seedHashAfter` | `equal` | 等于本侧基线哈希 | 绑定的失败运行 | 判定 |
|---|---|---|---|---|---|
| `02a-11-seed-intact` | ✅（`sha256:a42f547c…`） | `true` | ✅ = `02a-04` 输出 | `02a-10-neg-oracle-drift` | PASS |
| `02b-19-seed-intact` | ✅（`sha256:67ea456b…`） | `true` | ✅ = `02b-12` 输出 | `02b-18-neg-oracle-drift` | PASS |

### 4.3 可证伪条件（这些断言「怎样才会变红」）

| 断言 | 可证伪读数 |
|---|---|
| 反例守卫存在 | 删掉守卫 ⇒ 对应捕获的 `EXIT` 变为 0，`--- STDERR ---` 段失去预期子句，`02b-22` 的 verdict 翻转 |
| 失败不损 seed | 若 seed 根被写入 ⇒ `seedHashAfter` ≠ `seedHashBefore`、`equal=false`；且 `02a-04`/`02b-12` 的基线哈希不再可复现 |
| 分离（同侧） | 若 seed 与 candidate 指同一根 ⇒ 生成 `02a-07`/`02b-15` 同型错误，成功捕获（`02a-06`/`02b-14`）不可得 |
| oracle 独立性 | 若 `expected` 落在候选根内 ⇒ `02a-08`/`02b-16` 同型错误 |
| oracle 一致性 | 若 actual 与 expected 结构化不等 ⇒ `02a-10`/`02b-18` 同型错误 |

**但须注意**：`02a-11`/`02b-19` 只绑定**漂移反例**，而漂移反例的注记明确写「drift applied to a separate candidate copy」——该失败路径**结构上不会写 seed 根**。真正有写入 seed 风险的是 `02a-07`/`02b-15`（candidate 指向 seed 子目录），但**该反例之后没有再次读取 seed 哈希**，也无时间戳或执行顺序读数。见 F-05。

## 五、发现清单

| 编号 | 严重度 | 位置 | 发现 |
|---|---|---|---|
| F-02 | **Medium** | `02-win-captures/bindings.json:1-14`、`02-win-captures/index.json:1-60`（`02a-03-build-release-tool`） | Windows 侧「核验者」`prfrail-release.exe` 由**同一 HEAD 工作树**构建（`go build … ./cmd/prfrail-release`），`bindings.json` 不含任何 `verifierCommit`；因此 Windows 成功运行**不能**证明「核验者 ≠ 候选」。独立核验者结构仅由 Linux 侧（`18c9340` vs `baa144d`）提供。 |
| F-03 | **Medium** | `02-linux-captures/index.json:24`（`02b-02`）、`:38`（`02b-03`）、`:52`（`02b-04`）；对照 `SHA256SUMS.txt:4-6` | 3 条 `index.sha256` 滞后：与磁盘实算及 `SHA256SUMS.txt` 均不符（载荷本身无损）。同包两个 `index.json` 的 `rawBytes` **语义不同**（Windows = stdout 正文；Linux = 存储文件），字段同名异义。 |
| F-04 | **Medium** | `02-win-captures/02a-07..02a-10`（`COMMAND:` 行）、`02-linux-captures/02b-18` | 反例**不可重放**：Windows 四条命令含 `<seed>` / `<hash>` / `<64>` / `<drifted actual.json>` 占位符；漂移反例所依赖的 modified `actual.json` **未归档**。故「反例确实由该输入触发」只能靠 stderr 串与退出码认定，无法从包内复原输入。 |
| F-05 | **Medium** | `02-win-captures/02a-11-seed-intact.txt:4`、`02-linux-captures/02b-19-seed-intact.txt:4` | 「失败不损 seed」的证据只绑定漂移反例（`02a-10` / `02b-18`），该类失败**不写 seed 根**；`roots-overlap` 反例（唯一指向 seed 子目录者）之后**无 post-run seed 哈希重读**，且全部捕获无时间戳/顺序读数 ⇒ 该断言对最危险路径**可证伪性弱**。 |
| F-06 | **Medium** | `00-briefing.md:15`（逐字原文）vs `freedom-list.md:34`、`MANIFEST.md:21`、`docs/t027/REMAINING_SLICES.md:1627` | `D = 1` 的**成因表述自相矛盾**：归档的逐字原文写「主控已主动发现并停手」，而登记写「**经用户强制打断纠正**（非主控自行发现）」并声明「主控不得以『主动发现』自助降级」。同一事件两种互斥描述并存于同一证据包。 |
| F-07 | **Medium** | 全仓库检索 | 任务书所述 `A-1`…`A-7` 的**编号化登记在磁盘上不存在**：`docs/**` 内检索「误落包根」「A-1」等（P55 上下文）无命中；仅 `MANIFEST.md:9` 与 `MANIFEST.md:11` 以**无编号**段落覆盖 A-4 / A-5 的实质内容。③ 时点无法按编号核验 A-1 / A-2 / A-3 / A-6 / A-7。 |
| F-08 | **Low** | `02-win-captures/`（无二进制捕获） | Windows 侧 `work/bin/prfrail.exe`、`prfrail-release.exe` **未做哈希绑定**，且 `work/` 已整目录删除 ⇒ 该侧构建产物不可复核（Linux 侧 3 枚二进制哈希齐备）。 |
| F-09 | **Low** | `02-win-captures/expected.json`（249 B）vs `actual.json`（313 B） | `OracleHash` 是**expected 文件的字节哈希**；而 actual 文件字节哈希为 `sha256:adc6e927…`（≠ `388f8884…`）却判定成功 ⇒ 比较是**结构化/规范化**语义，非字节相等。该口径未在包内任何文件中写明，易被误读为字节一致。 |
| F-10 | **Low** | `02-win-captures/02a-12-candidate-listing.txt:7`、`02-linux-captures/02b-20-candidate-listing.txt:11` | `expectedAbsentInCandidate=true` 是**自撰断言行**而非命令输出；虽可由同文件目录清单佐证，但不构成独立读数。 |
| F-11 | **Low** | `02a-11`、`02a-12`、`02b-19`、`02b-20`、`02b-22` | 5 件载荷体内无 `EXIT` 字段，其退出码仅来自 `index.json` 自述（`02b-19`/`02b-20`/`02b-22` 另有 `02b-22` 的逐件表兜底；`02a-11`/`02a-12` 无兜底）。 |
| F-12 | **Low** | `02-linux-captures/02b-00-env.txt:16-19`；`MANIFEST.md:9` | 预取阶段 `GOPROXY=https://goproxy.cn,direct` + `GOSUMDB=off` 的**副作用无读数**：是否改写 `go.sum`、是否影响随后被测的 seed 树哈希（`02b-12` 在预取之后计算）均未捕获；`GOPROXY=off` 只能证明**构建期**无网络，不能证明**被测树**未被预取步骤改动。 |
| F-13 | **Low** | `02a-04`（`work/seed`）、`02b-12`（`/tmp/…/seed`） | 被哈希的 seed 根**与固定 seed 提交无绑定读数**：Windows 侧无 seed 提交记录；Linux 侧 `02b-02` 的 worktree 在 `…/seed-src`，而 `02b-12` 哈希的是**另一个路径** `…/seed`，两者之间的组装/拷贝步骤未捕获（`02b-19` 仅注记「written only once at setup」）。seed 二进制本身的构造是提交绑定的（`02b-02` → `02b-08`），故「seed」前提**部分成立**，但 `SeedHash` 覆盖的对象无法回溯到 `fc99f55` 的树。 |

**无 Critical / High 发现。** 载荷哈希链、三方集合相等、双侧夹具/oracle 一致、8 条反例与 2 条 seed 完好断言均实测成立；上述发现全部为**证据链完备性 / 元数据一致性**问题，不涉及「断言落空」或「以叙述代替读数」的硬缺陷（`F-10`/`F-11` 为叙述型残留，已单列 Low）。

## 六、异常登记复审意见

| 登记 | 主张 | 与证据的一致性 | ③ 意见 |
|---|---|---|---|
| `D = 1` | 无触发词抢跑，经用户强制打断纠正，就地转正，未造成污染 | 载体齐备（`freedom-list.md:34`、`MANIFEST.md:21`、台账 `REMAINING_SLICES.md:1627`）；「未 commit / 未 push、HEAD 保持 `baa144d`」可复算（`git rev-parse HEAD` = `baa144d`） | **成立**，但成因表述与 `00-briefing.md` 逐字原文冲突（**F-06**），须在 ⑧b 统一为单一口径 |
| `A-1` | 首轮捕获误落包根 | 包内无残留文件；`SHA256SUMS` 与磁盘均无包根捕获件 | **不可复算**（无残留、无编号载体，**F-07**）；与「只增不删」口径存在张力，建议在 ⑧b 补一句「误落件已迁入 `02-win-captures/` 并在 `index.json` 以 `02a-oracle-*` 编号登记」的可核验说明 |
| `A-2` | 会话级 `GOCACHE`/`GOMODCACHE` 污染 | `02a-00` 显示隔离路径、`MANIFEST` 说明 `work/` 已删 | **部分成立**：仅能证明最终使用了隔离缓存，污染事件本身无读数 |
| `A-3` | ② 主控自执行（用户同轮授权） | `00-briefing.md:66`「主控应亲自执行或指派 Flash 执行命令并收集证据」 | **成立且可逐字核验** ✅ |
| `A-4` | seed 提交 `fc99f55` 无发布工具 ⇒ 改用 T017 三角色结构；用户已批准 | 独立复算：`git ls-tree -r fc99f55` **不含** `cmd/prfrail-release/**` 与 `internal/release/**`；`git cat-file -t fc99f55` = commit；`fc99f55` 主题 = `feat: simplify S1 release trust to baseline evidence` | **技术前提成立** ✅。「用户批准」在包内无载体（**F-07**）。**正面作用应加重记录**：该改法使 Linux 侧形成 seed `fc99f55` / verifier `18c9340` / candidate `baa144d` 三者互异、三枚二进制哈希互异的结构，**独立性显著强于** Windows 侧「同一工作树既出候选又出核验者」的形态（对照 **F-02**） |
| `A-5` | VM 无法访问 `proxy.golang.org`/`sum.golang.org` ⇒ `goproxy.cn` 预取 + 构建期 `GOPROXY=off` | `02b-00`：`GOPROXY_preload=https://goproxy.cn,direct`、`GOPROXY_build=off`、`GOSUMDB=off`、双 note 行；`MANIFEST.md:9` 同步登记 | **成立** ✅；预取副作用未观测（**F-12**） |
| `A-6` | 前两轮远端各 9 项 FAIL | 包内仅有最终轮 `02b-22`（21 项全 PASS），无前轮 FAIL 捕获 | **不可复算**（VM 已清理）；与「只增不删」口径的偏离须在 ⑧b 如实说明 |
| `A-7` | 远端模块缓存只读导致 `rm` 残留 | 无读数（远端工作区已清理） | **不可复算**；建议记录为「环境约束」而非「证据」 |

**复审结论**：`D = 1`、`A-3`、`A-4`、`A-5` 实质成立且可核验；`A-1`、`A-2` 部分成立；`A-6`、`A-7` 不可复算；**全部七条的编号化载体在 ③ 时点缺失（F-07）**，`D = 1` 的成因表述存在内部矛盾（F-06）。二者必须在 ⑧b 落盘为单一、可检索的编号清单。

## 七、未覆盖 / 无法验证事项

1. **Windows 侧工具链与 CI 不同版**：本机 `go1.27.0`，CI 固定 `1.22.12`（`ci.yml:46`）。Windows 侧不构成本片「真 seed 构建」证据，与 CI 同版证据只在 Linux 侧（`02b-00`）。
2. **Windows 侧无 seed 二进制构建**：`02-win-captures/` 无 seed 构建与 seed 二进制版本捕获，`SeedHash` 是对 `work/seed` 目录的树哈希。
3. **`SeedHash` 的对象不可回溯到提交**（F-13）：两侧均无「seed 根 ≡ 某提交树」的等式读数。
4. **`CandidateHash` 不是源码哈希**：Windows 候选根仅 `actual.json`；Linux 候选根 3 条目（含 `candidate-commit.txt`）。跨侧数值不可比。
5. **反例不可重放**（F-04）：占位符命令 + 漂移输入件未归档。
6. **seed 完好断言的路径覆盖不足**（F-05）：未覆盖「candidate 指向 seed 子目录」后的 seed 状态。
7. **oracle 比较不是字节相等**（F-09）：`OracleHash` 只覆盖 expected 文件；actual 以结构化比较通过。**该比较只证明「行为未漂移」，不构成「候选无恶意」的证明**——恶意但行为等价的候选可以通过本证据链。
8. **远端环境副作用不可复核**（F-12、A-6、A-7）：VM 工作区已清理，`go.sum` 是否被预取步骤改动、前两轮 FAIL 的原始日志均不可得。
9. **并发 / 崩溃类测试点不可测（N/A）**：AT-13 在本片无并发与崩溃断言，五类覆盖中的「并发」「崩溃」两类**不适用**，非遗漏。
10. **`03-test-report.md` 不在 `SHA256SUMS.txt` / `MANIFEST.md` 内**：③ 环节按授权只新增文件、不改元数据；**⑨ 复核时必须把本件补入 `MANIFEST.md` 与 `SHA256SUMS.txt`**，否则下一轮三方相等检查将出现 1 件未列载荷（本次审计已将其排除在 44 件载荷集之外）。

## 八、覆盖矩阵

| 编号 | 断言 / 测试名 | 类型 | 对应测试点 |
|---|---|---|---|
| A | 三方集合相等（`MANIFEST` ≡ `SHA256SUMS` ≡ 磁盘） | 正例（结构） | 证据链完整性 |
| B | 44 条 sha256 全量重算 | 正例（机械） | 证据链完整性 |
| C | `index.json` 元数据自洽 | 正例（机械，发现滞后） | 证据链完整性 / 元数据一致性 |
| D | 编码与行尾（BOM / LF） | 正例（机械） | 仓库编码契约 |
| E-1 | `HEAD` / `SourceCommit` 三方相等 | 正例（交叉） | seed 与 candidate 分离（锚点） |
| E-2 | 五枚夹具四面一致（win / linux / CI 输入表 / 本地磁盘重算） | 正例（交叉） | 外部 oracle 比较（基准） |
| E-3 | `OracleHash` 双侧相等且 = expected 文件哈希 | 正例（交叉） | 外部 oracle 比较 |
| E-4 | 同侧 seed ≠ candidate；跨侧差异定性为构造差异 | 正例 + 定性（边界） | seed 与 candidate 分离 |
| E-5 | 三角色提交互异、三枚二进制哈希互异（Linux） | 正例（边界） | 核验独立性 |
| E-6 | Windows 侧无 `verifierCommit`、无二进制哈希 | 反例（缺口识别） | 核验独立性 |
| T3-1 | roots-overlap 反例（win / linux） | 反例 | seed 与 candidate 分离 |
| T3-2 | oracle-ownership 反例（win / linux） | 反例 | 外部 oracle 比较（独立性） |
| T3-3 | seed-changed 反例（win / linux） | 反例 | seed 固定性 |
| T3-4 | oracle-drift 反例（win / linux） | 反例 | 外部 oracle 比较（漂移检测） |
| T3-5 | 失败后 seed 哈希前后对照（win / linux） | 反例后置断言 | clean-room 失败不损 seed |
| — | 并发类 | **N/A（不可测）** | 本片无并发断言 |
| — | 崩溃类 | **N/A（不可测）** | 本片无崩溃注入断言 |

**测试 diff**：无产品代码或测试代码改动（本片边界禁止改 Go 代码）；本环节产出为两个只读审计脚本（`tmp/EXIT-SEED-BUILD-EVIDENCE/`，⑧c 统一清理）与本报告。

## 九、结论

- 载荷哈希链：**44/44 命中，0 失配**；三方集合相等 **成立**；编码行尾 **0 违规**。
- 双侧交叉核对：`HEAD`、五枚夹具（四面）、`OracleHash` **全部一致**；seed/candidate 跨侧差异属**构造差异**，非失配。
- 失败路径：8 条反例（非 0 退出 + 预期 stderr + 空 stdout）与 2 条 seed 完好断言 **全部成立**。
- 发现：0 Critical / 0 High / 6 Medium / 6 Low（详见 §五）；主要缺口为 **Windows 核验者非独立（F-02）**、**反例不可重放（F-04）**、**seed 完好断言路径覆盖弱（F-05）**、**`index.json` 元数据滞后与字段同名异义（F-03）**、**D=1 成因表述矛盾（F-06）**、**A-1…A-7 编号载体缺失（F-07）**。
- 上述发现均不改变「证据链自洽」的判定，但须由 ⑤ / ⑦ 决定是否升至整改；F-07 与 `03-test-report.md` 的元数据补录属 **⑧b / ⑨ 必办项**。

TEST: DONE

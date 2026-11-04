# 03-test-report-r2 — P55 `EXIT-SEED-BUILD-EVIDENCE` ③ 校验测试 第 2 轮（复审）

| 字段 | 内容 |
|---|---|
| 切片 | P55 `EXIT-SEED-BUILD-EVIDENCE`（AT-13 seed 驱动构建补证；层 3 第 7 项） |
| 环节 | ③ 校验测试 · **第 2 轮（复审）**：只验证第 1 轮 6 项 Medium 的处置是否闭合，不重做整轮审计 |
| 锚点 | `HEAD` = `baa144df5898817e0ac7cd860928bb4a90cbd009`（实测 `git rev-parse HEAD` 一致；A-4 / D=1 的「HEAD 未变」主张成立） |
| 复核时点 | 2026-11-05 |
| 输入 | 主控处置 4 件：`anomalies.md`、`claim-scope.md`、`replay.md`、两份 `index.json` 重建 |
| 审计性质 | 只读。未修改任何既有证据文件、`03-test-report.md`、`SHA256SUMS.txt`、`MANIFEST.md`、产品代码或 `.github/workflows/**` |
| 本环节新增件 | `tmp/EXIT-SEED-BUILD-EVIDENCE/verify-sums-r2.js`、本文件 |

> **范围声明**：本轮**未重跑任何构建**。所有判定基于「重算哈希 / 重读 index 元数据 / 核对新增工件与磁盘事实是否一致 / 静态核对重放指令与实现签名」。

## 一、本轮范围

| 编号 | 任务 | 方法 | 结论 |
|---|---|---|---|
| V1 | 重跑第 1 轮哈希脚本，比较 44 → 48 的条目变化 | `node tmp\EXIT-SEED-BUILD-EVIDENCE\verify-sums.js` | 条目 48 已达成；脚本仍 exit 1，原因见 §二.3 |
| V2 | 用**正确载荷口径**重算三方相等与全量哈希 | `node tmp\EXIT-SEED-BUILD-EVIDENCE\verify-sums-r2.js` | **三方相等 true；48/48 命中；0 失配；exit 0** |
| V3 | `index.json` 派生口径核验（`bytes`/`sha256` 对磁盘、`captureMeta` 保留） | V2 的 C 段 + PowerShell 独立复算 | **39/39 条硬校验全通过** |
| V4 | 编码 / 行尾 | V1 的 D 段 + V2 的 D 段 | **0 违规** |
| V5 | 6 项 Medium 逐条闭合判定（含与磁盘事实的一致性核对） | 静态核对 + 独立复算 | 见 §三 |
| V6 | 重放指令可执行性（F-04）静态核对 | 与 `cmd/prfrail-release/main.go` 参数表比对；一次性最小包实测 `go build -o` / `Copy-Item` 的目录前置条件 | 指令签名一致；残留 2 项 Low（§五） |

## 二、重跑读数

### 2.1 条目数变化（44 → 48）

| 项 | 读数 |
|---|---|
| 包内文件总数 | **50**（= 48 载荷 + `MANIFEST.md` + `SHA256SUMS.txt`） |
| 磁盘载荷 | **48** |
| `SHA256SUMS.txt` 条目 | **48**（解析失败 0、重复 0） |
| `MANIFEST.md` 归档条目 | **48** |
| 组成式 | 第 1 轮 44 + `anomalies.md` + `claim-scope.md` + `replay.md` + `03-test-report.md` = **48** |

> 说明：第 1 轮 44 件不含 `03-test-report.md`（当时按授权未登记）。主控本轮已把它连同 3 件新工件一并登记，故增量为 **+4**。这与第 1 轮 §七 第 10 条「⑨ 必须把本件补入 `MANIFEST` 与 `SHA256SUMS`」的必办项一致，属**提前完成**，非偏差。

### 2.2 哈希与集合（V2，正确口径）

| 项 | 读数 |
|---|---|
| `set_disk_eq_sums` | **true**（`unlisted_on_disk` 0、`listed_but_absent` 0） |
| `set_manifest_eq_sums` | **true**（`manifest_only` 0、`sums_not_in_manifest` 0） |
| sha256 重算 | **48 / 48 命中** |
| 失配清单 | **（空）** |
| 编码 / 行尾违规 | **0**（`.md` 均 BOM+LF；`.txt`/`.json` 均无 BOM+LF；`\r\n` 计数 0） |
| 脚本退出码 | **0**（`RESULT=OK`） |

### 2.3 第 1 轮脚本仍 exit 1 的原因（任务书要求逐项指认）

`verify-sums.js` 本轮输出 `RESULT=FAIL / EXIT=1`，剩余原因**两条**，均**不是**载荷失配：

| # | 原因 | 读数 | 定性 |
|---|---|---|---|
| 1 | **脚本口径假阳性**：该脚本把 `03-test-report.md` 硬编码排除出载荷集（当时它确实未登记）。主控现已把它登记为载荷，于是 `payload(47) ≠ sums(48)`，报 `listed_but_absent=["03-test-report.md"]`、`manifest_absent_on_disk=["03-test-report.md"]` | `disk_payload_files=47` / `sums_entries=48` | **脚本问题，非包缺陷**。`Test-Path` 实测该文件存在（23973 B）；V2 用正确口径复算即 `true` |
| 2 | **`index.json` 段的结构性残项**：`index_issues=7`，**全部**为 `no_EXIT_line_in_capture`（7 件捕获体正文无 `EXIT:` 行，退出码只存在于 `index.json` 自述） | 7 条 | **低危结构性**，即第 1 轮的 **F-11**。第 1 轮 C 段的**滞后类**问题（`index_sha256_mismatch` / `rawBytes_mismatch` / `index_vs_sums_divergence`）**已全部消失**（V3 硬校验 0 条） |

第 1 轮「退 1 的唯一来源 = `index` 元数据段」这一判断**已不再成立**：元数据滞后已被消除，剩下的 7 条属 F-11 的原始形态，另外多出一条脚本口径假阳性。**任务书预期的 exit 0 在正确口径下达成（V2 exit 0）**。

### 2.4 `index.json` 派生口径核验（F-03 本体）

| 项 | 读数 |
|---|---|
| `index` 条目总数 | **39**（Windows 16 + Linux 23） |
| `captureMeta` 覆盖 | **39 / 39** |
| `bytes` 与磁盘字节数一致 | **39 / 39** |
| `sha256` 与磁盘实算一致 | **39 / 39** |
| `sha256` 与 `SHA256SUMS.txt` 一致 | **39 / 39** |
| 硬校验问题数 | **0** |
| 第 1 轮记录的 3 条滞后哈希 | 已修正：`02b-02` 现为 `7dadf79b…`、`02b-03` `f1d5185e…`、`02b-04` `2fe8e78a…`（与 `SHA256SUMS.txt:4-6` 逐字节一致） |
| 第 1 轮记录的 10 条 `rawBytes` 同名异义 | 已消除：`captureMeta.rawBytes` 保留捕获期读数（如 `02a-04` = 72，`bytes` = 251），`bytes` 改为派生口径 |

## 三、6 项 Medium 闭合判定

| 编号 | 处置工件 | 判据（文件 + 读数） | 判定 |
|---|---|---|---|
| **F-02**<br>Windows 核验者非独立 | `claim-scope.md` §二 **N-1** | ① `claim-scope.md:21` 明写「Windows 侧证明了核验者独立性」**不成立**，并限定独立性仅在 Linux 侧；② 与磁盘事实一致：`02-win-captures/bindings.json` 键集 = `head`/`seedHash`/`localGo`/`ciGo`/`fixtures`/`directiveVersion`，**无 `verifierCommit`**；③ `02-win-captures/index.json` 中 `02a-03-build-release-tool.command` = `go build -trimpath -o work/bin/prfrail-release.exe ./cmd/prfrail-release`（同一 HEAD 工作树）；④ `S-1` 把「三角色 / 三二进制互异」限定在 Linux | **闭合** |
| **F-03**<br>`index` 元数据滞后 + 字段同名异义 | 两份 `index.json` 重建 + `MANIFEST.md` 口径段 | ① §2.4 读数：39/39 `bytes`、`sha256` 硬校验通过，`captureMeta` 39/39 保留；② `MANIFEST.md` 末段「`index.json` 口径（v2）」明确 `bytes`/`sha256` 为派生元数据、捕获期读数移入 `captureMeta.rawBytes`/`captureMeta.crlfNormalized`；③ 第 1 轮 3 条滞后哈希与 10 条同名异义均已消失 | **闭合** |
| **F-04**<br>反例不可重放 | `replay.md` + `claim-scope.md` **N-4** | ① `replay.md` 给出 Windows（PowerShell）与 Linux（bash）两套可执行形式；② 命令签名与实现一致：`cmd/prfrail-release/main.go:26-36`（`hash-tree --root`）、`:38-49`（`probe --binary/--noop-chain/--executable-chain/--work-root/--output/--expected-version`）、`:50-62`（`selfhost --commit/--seed/--seed-hash/--candidate/--expected/--actual`）**逐个命名吻合**；③ 漂移输入件重建脚本与 `tmp/EXIT-SEED-BUILD-EVIDENCE/p55-capture-linux.sh:150-159` **逐语句等价**（读 `fixtures/expected.json` → `run-noop.outcome=failed` → 写目标），故 `02b-18` 的输入可确定性复原；④ `N-4` 明确捕获行是「显示占位符」、不可原地粘贴 | **闭合**（残留 2 项 Low，见 R2-1 / R2-6） |
| **F-05**<br>「失败不损 seed」覆盖路径弱 | `claim-scope.md` **N-5** | ① `N-5` 把断言从「文件系统级审计」收窄为「代码层顺序 + 只读」；② **行号顺序实测成立**：`internal/release/verify.go:94` `if pathsOverlap(seedRoot, candidateRoot)` → `:95` 返回 `seed and candidate roots overlap`；`:110` `seedHash, err := HashTree(seedRoot)` —— overlap 检查**先于** seed 哈希 16 行；③ `VerifySelfHost`（`:82-135`）体内**无任何写调用**：`existingRealPath`（`filepath.Abs` + `EvalSymlinks`）、`pathWithin`（`filepath.Rel`）、`HashTree`（读目录树）、`readOracle`（`os.ReadFile`）；④ 推论可验证：`02a-07`/`02b-15` 的 overlap 反例在 `:94` 即返回，`:110` 的 `HashTree` **根本不会执行** ⇒ 「反例运行写入 seed」在 verify 进程内结构性不可能；⑤ `p55-capture-linux.sh:127-134` 证实 overlap 反例只做 `cp` + `hash-tree`（只读），未调用 probe | **闭合** |
| **F-06**<br>`D = 1` 成因表述矛盾 | `anomalies.md` 末节 + `freedom-list.md` §三 | ① 冲突事实存在且未被抹除：`00-briefing.md` 围栏内保留原话「主控已主动发现并停手」（逐字归档）；② 权威表述落盘：`anomalies.md` 末节 = **「经用户强制打断纠正」（非主控自行发现）**，并给出「归档即原文、旧措辞属被作废表述、不得引用」的层级规则；③ 同口径扩散：`freedom-list.md` §三 `D-1` / `D-1 附带`（含「主控不得以『主动发现』自助降级」）；④ 台账一致：`docs/t027/REMAINING_SLICES.md:1627`（`R / D 基线` 行，行号实测命中）= 「无触发词抢跑，**经用户强制打断纠正**，就地转正，未造成污染」 | **闭合** |
| **F-07**<br>`A-1…A-7` 编号载体缺失 | `anomalies.md` 全表 | ① 编号化落盘：`D = 1` + `A-1`…`A-7`，每行含「事实 / 证据·判据 / 处置 / 复核状态」四列（第 1 轮 F-07 的直接要求）；② 独立复算逐条：**A-1**「`index.json` 16 条为修正后产物」→ 实测 Windows 索引 **16** 条 ✔；**A-3** 引文「主控应亲自执行或指派 Flash 执行命令并收集证据」→ `00-briefing.md` §六第 1 条**逐字存在** ✔；**A-4** `git cat-file -t fc99f55` = `commit`、`git ls-tree --name-only fc99f55 cmd/` = `cmd/prfrail`（仅 1 项）、`git ls-tree -r fc99f55 internal/release` **为空**、subject = `feat: simplify S1 release trust to baseline evidence` ✔；**A-5** `02b-00-env.txt:16-19` 含 `GOPROXY_preload=https://goproxy.cn,direct`、`GOPROXY_build=off`、`GOSUMDB=off` + 2 条 note ✔；**A-6** `02b-22-verdicts.txt` 21 行（20 PASS + 1 INFO，**0 FAIL**）✔；③ 不可复核项**已如实标注**：A-2（门禁读数不在包内）、A-6（原始件随 VM 清理）、A-7（远端副作用） | **闭合**（编号化 + 判据 + 不可复核声明齐备；A-1/A-2 的状态措辞建议见 R2-4） |

**小结**：6 项 Medium **全部闭合**，无「部分闭合」、无「未闭合」。

## 四、仍未闭合项

**Medium 及以上：无。**

以下为**本轮残留**（均 Low，不阻断，登记供 ⑧b / ⑨ 决定是否顺带修）：

| 编号 | 位置 | 内容 | 严重度 |
|---|---|---|---|
| R2-1 | `replay.md` 反例表 | 三条反例（根重叠 / oracle 归属 / 漂移）都要求候选根**已存在**。`VerifySelfHost` 在 `:85-93` 先做 `existingRealPath(candidateRoot)`，目录不存在时报 `candidate root: …`，**不会**产生表中期望的 stderr。实测佐证：`go build -o` 与 `Copy-Item` 会自动建父目录，故成功路径无碍；但 `replay.md` 未给出 `mkdir -p`/`New-Item`，照抄执行会得到非预期错误串（`p55-capture-linux.sh:127,136,149` 是有 `mkdir -p` 的） | Low |
| R2-2 | `claim-scope.md` N-2 | Linux seed 根实为 **4 项**：`bin/prfrail` + `fixtures/` + `seed-commit.txt` + `toolchain.txt`（`p55-capture-linux.sh:95,106,107,110`，且 `:110` 写在 `:112` 之前 ⇒ `toolchain.txt` 在 `SeedHash` 覆盖范围内）。N-2 只列了前三项。**结论方向不变**（仍非 `git archive` 树哈希），但「SeedHash 含工具链字符串」是可用信息，建议补一句 | Low |
| R2-3 | `03-test-report.md` §二 | 第 1 轮「捕获体无 `EXIT:` 字段」计 **6** 件，实测 **7** 件（漏计 `02b-00-env.txt`：该件同样无 `EXIT:` 行，退出码仅见 `index.json`）。计数更正，不影响任何断言 | Low |
| R2-4 | `anomalies.md` A-1 / A-2 | 两行「复核状态」写「已闭环」，但其判据在包内**不可复核**（A-1 的 13 个误落件仅能由「修正后索引 16 条存在」反推；A-2 引用的 `FAIL base gofmt -l . :: tmp\…` 门禁读数不在包内）。建议改为「**部分可核验**」，与 A-6 / A-7 的「不可复核」口径对齐 | Low |
| R2-5 | `claim-scope.md` S-6 | S-6 措辞为「失败运行后 seed 根哈希不变」，单引时会读作覆盖四类反例；实测只有漂移反例前后重读过哈希（`02a-11` / `02b-19`）。N-5 已限定，但 S-6 建议加注「（漂移反例）」，避免下游只引 S-6 | Low |
| R2-6 | `replay.md` §二 | Linux 段为「要点」级，依赖 `tmp/EXIT-SEED-BUILD-EVIDENCE/p55-capture-linux.sh`（文内已声明「⑧c 清理前可取」）。若该脚本被 ⑧c 删除，Linux 侧重放将缺少 `$WORK`/`$SEED`/`$VERIFIER_BIN`/`$CAND` 等变量定义与 `mkdir -p` 序列。建议把该脚本的关键片段固化进 `replay.md`，或把脚本转为证据包载荷 | Low |

## 五、12 项 Low 的处置归属

> **口径更正**：第 1 轮的**编号 Low 为 6 项**（F-08…F-13），与任务书所述「12 项 Low」不符。若把第 1 轮 §七「未覆盖 / 无法验证事项」（10 条）计入，则遗留项共 16 条。下表先按 **6 项编号 Low** 逐条归属，再补 §七 的归属。

### 5.1 编号 Low（F-08…F-13）

| 编号 | 内容摘要 | 本轮工件是否覆盖 | 归属 |
|---|---|---|---|
| F-08 | Windows 侧二进制（`prfrail.exe`/`prfrail-release.exe`）未哈希，`work/` 已删 | **部分覆盖**：`claim-scope.md` N-1 把 Windows 侧能力收窄为「工具在自身产物上可跑通 + oracle 语义 + 反例 fail-close」，S-1 把「三枚二进制互异」限定到 Linux ⇒ 不再声称二进制可复核 | 留作观察（本片边界内不可补：工作区已删除） |
| F-09 | `OracleHash` 是 expected **文件字节哈希**，actual 以**结构化**比较通过 | **已覆盖**：`claim-scope.md` N-3 明写「`OracleHash` 证明 actual 与 expected 字节相同」不成立，并给出 `adc6e927…` ≠ `388f8884…` 的读数；实测 `internal/release/verify.go:125` 走 `equalOracle(expected, actual)`，`:259`（`readOracle` 内）`evidence.Digest("", data)` 只对 expected 取字节哈希 | **闭合** |
| F-10 | `expectedAbsentInCandidate=true` 是自撰断言行，非命令输出 | **未覆盖**（`claim-scope.md` / `anomalies.md` / `replay.md` 均未提及） | 留作观察 |
| F-11 | 5 件载荷体内无 `EXIT` 字段，退出码仅来自 `index.json` 自述 | **未覆盖**；本轮重测实为 **7 件**（见 R2-3） | 留作观察（附计数更正） |
| F-12 | 预取阶段 `goproxy.cn` + `GOSUMDB=off` 的副作用无读数（是否改写 `go.sum`、是否影响 `02b-12` 树哈希） | **部分覆盖**：`claim-scope.md` S-7 与 `anomalies.md` A-5 均把结论限定在「**构建期** `GOPROXY=off` ⇒ 构建无网络依赖」，不再声称预取无副作用 | 留作观察（VM 已清理，不可补） |
| F-13 | 被哈希的 seed 根与固定 seed 提交无绑定读数 | **已覆盖**：`claim-scope.md` N-2 明写「`SeedHash` 等于 seed 提交源码树哈希」不成立，并说明 seed 根是**组装目录**；`p55-capture-linux.sh:95-110` 实测构成（见 R2-2） | **闭合** |

### 5.2 第 1 轮 §七「未覆盖 / 无法验证」10 条的归属

| §七 条目 | 归属 |
|---|---|
| 1 Windows 工具链与 CI 不同版 | 已由 `claim-scope.md` N-6 显式收窄 ⇒ **闭合（登记化）** |
| 2 Windows 无 seed 二进制构建 | 与 F-08 同源，N-1/N-2 覆盖 ⇒ **部分覆盖** |
| 3 `SeedHash` 对象不可回溯到提交（F-13） | N-2 ⇒ **闭合** |
| 4 `CandidateHash` 不是源码哈希 | 未新增载体（第 1 轮已在 §三 定性）⇒ 留作观察 |
| 5 反例不可重放（F-04） | `replay.md` + N-4 ⇒ **闭合** |
| 6 seed 完好断言路径覆盖不足（F-05） | N-5 ⇒ **闭合** |
| 7 oracle 比较非字节相等（F-09） | N-3 / N-7 ⇒ **闭合** |
| 8 远端环境副作用不可复核（F-12、A-6、A-7） | `anomalies.md` A-6/A-7 + S-7 ⇒ **已登记为不可复核**，实质仍不可补 |
| 9 并发 / 崩溃类 N/A | `claim-scope.md` §三「未覆盖（本片范围外）」⇒ **闭合（口径声明）** |
| 10 `03-test-report.md` 未登记入元数据 | 主控已登记（48 条含该件）⇒ **闭合**；本轮新件 `03-test-report-r2.md` 由 **⑨ 补登记** |

## 六、未验证事项

1. **未重跑构建**（本环节授权边界）：两侧 `selfhost` 成功路径与四条反例均未现场执行；本轮只做静态可执行性核对（命令签名 + 目录前置条件）。
2. **重放指令未实测执行**：`replay.md` 的命令**签名**与实现一致，但端到端跑通需落地工作区（Windows 侧 `work/` 与 Linux VM 均已清理），本轮未验证。
3. **远端不可复核项**：A-6（前两轮 9 项 FAIL 原始日志）、A-7（`rm` 残留）、F-12（预取是否改写 `go.sum`）—— VM 工作区已清理，均不可补。
4. **`03-test-report-r2.md` 未登记入 `MANIFEST.md` / `SHA256SUMS.txt`**：按本环节授权只能新增文件，**⑨ 必须补录**，否则下一轮三方相等检查将出现 1 件未列载荷（本轮脚本已显式排除它）。
5. **并发 / 崩溃类仍未覆盖（N/A）**：AT-13 在本片无并发与崩溃断言，非本片遗漏。
6. **`claim-scope.md` 与 `anomalies.md` 的「用户批准」类主张无包内载体**（A-4/A-5/A-3 的裁定来源），只能由会话记录佐证；本轮按「事实登记 + 不可复核声明」接受。

## 七、覆盖矩阵（本轮）

| 编号 | 测试 / 断言 | 类型 | 对应处置 |
|---|---|---|---|
| V1-A | 三方集合相等（正确载荷口径） | 正例（结构） | F-03 |
| V1-B | 48 条 sha256 全量重算 | 正例（机械） | 载荷完整性 |
| V1-C | `index.json` 派生元数据（`bytes`/`sha256`/`captureMeta`） | 正例（机械） | F-03 |
| V1-D | 编码 / 行尾 | 正例（机械） | 仓库编码契约 |
| V2-1 | `claim-scope.md` N-1 与 `bindings.json` 键集 / `02a-03` 命令一致性 | 正例（交叉） | F-02 |
| V2-2 | `replay.md` 命令签名 vs `cmd/prfrail-release/main.go` 参数表 | 正例（交叉） | F-04 |
| V2-3 | 漂移重建脚本 vs `p55-capture-linux.sh:150-159` | 正例（交叉） | F-04 |
| V2-4 | `verify.go:94` < `:110` 顺序 + `VerifySelfHost` 无写调用 | 正例（源码锚定） | F-05 |
| V2-5 | `anomalies.md` 末节 vs `00-briefing.md` 原文 vs `REMAINING_SLICES.md:1627` | 正例（三文件） | F-06 |
| V2-6 | `A-1`/`A-3`/`A-4`/`A-5`/`A-6` 逐条复算 | 正例 + 反例（不可复核识别） | F-07 |
| V2-7 | `go build -o` / `Copy-Item` 目录前置条件实测（一次性最小包） | 边界 | F-04（R2-1） |

## 八、结论

- 载荷完整性：**48/48 命中、0 失配**；三方集合相等 **true**；编码行尾 **0 违规**（正确口径下 **exit 0**）。
- 第 1 轮脚本 exit 1 的**元数据滞后来源已消失**（硬校验 0 条）；剩余为 F-11 的 7 条结构性缺 `EXIT:` 行 + 1 条脚本口径假阳性。
- **6 项 Medium 全部闭合**（F-02/F-03/F-04/F-05/F-06/F-07），且新增工件与磁盘事实逐条一致（含 `verify.go` 行号顺序、`SeedHash` 构成、`git ls-tree` 三项事实）。
- **无 Critical / High / Medium 残留**；新增 6 项 Low（R2-1…R2-6），均不阻断，其中 R2-6（`replay.md` 依赖 tmp 脚本）建议 ⑧b 顺带固化。
- ⑨ 必办：把 `03-test-report-r2.md` 补入 `MANIFEST.md` 与 `SHA256SUMS.txt`。

TEST: DONE

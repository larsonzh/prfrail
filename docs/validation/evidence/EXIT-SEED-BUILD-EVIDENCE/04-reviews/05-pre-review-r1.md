# 05-pre-review-r1 — ⑤ 预审第 1 轮（原文归档）

归档说明：本文件为该轮审查方最终消息的**逐字原文**（纪律：审查方产出原文归档，主控不得按要点重建）。审查方 = `deep-reasoner`；调用模型 = `DeepSeek V4 Pro (deepseek)`；时点 2026-11-05。

```text
# ⑤ 预审报告 — P55 `EXIT-SEED-BUILD-EVIDENCE`

预审员：deep-reasoner（DeepSeek V4 Pro）｜只读，未修改任何文件｜复核时点 2026-11-05

---

## 一、核验动作与结果速览（§4.2.1 四项全部实际执行）

| 动作 | 结果 |
|---|---|
| 实际读取 12+ 关键工件（claim-scope / anomalies / 两轮 ③ 报告 / 双侧 6 个关键捕获 / 02b-21 / MANIFEST / SHA256SUMS / input-list / freedom-list / replay / 00-briefing / verify.go / probe.go / disposition.go / main.go / p55-capture-linux.sh / 台账 P55 条目） | 全部完成 |
| 独立复算哈希：`SHA256SUMS.txt` 条目数 = **49** = 磁盘载荷数 = 49；抽查 9 条（win 3 / linux 3 / 包根 3）**9/9 OK**；SHA256SUMS 无 BOM；`.md` BOM+LF、`.json` 无 BOM（抽查 2 件） | 全部一致 |
| `testdata/selfhost/` 五枚夹具独立 `Get-FileHash`：`ae0d5db8…` / `dbc5fc41…` / `388f8884…` / `ff9d243c…` / `fbeee2b6…` = 双侧 bindings = `s1-selfhost.md:91-95` CI 输入表（**四面一致**） | 全部一致 |
| 代码层断言：`verify.go:94` overlap 检查**先于** `:110` seed 哈希 16 行；`VerifySelfHost`（`:82-:138`）全部被调用函数均为只读（`existingRealPath`=Abs+EvalSymlinks、`pathWithin`=Rel、`HashTree`→`BuildReleaseInventory`（Abs/Stat/WalkDir/ReadFile）、`readOracle`=ReadFile+Digest、`equalOracle`=纯比较） | 成立 |
| 越界检查：`git diff --cached --name-only` 无 `.github/workflows/**`、无任何 `.go`；53 项全部位于 `docs/t027/` 与 `docs/validation/evidence/`；CN/EN 台账 diff 各 19 行且仅含 P55 条目（G3(a.1) 对称） | 无越界 |
| 锚点：`git rev-parse HEAD` = `baa144d…` ✓；`fc99f55` / `18c9340` / `baa144d` 三提交本地均可解析（`cat-file -t` = commit）；`fc99f55` 的 `cmd/` 仅 `cmd/prfrail`、`internal/release` **为空**（A-4 前提成立）；其 subject = `feat: simplify S1 release trust to baseline evidence` | 全部命中 |
| tmp 纪律：`tmp/EXIT-SEED-BUILD-EVIDENCE/` 仍存 14 件（含 `p55-capture-linux.sh`、`verify-sums*.js`），「只增不删」成立；其 `:127-134`（overlap 反例仅 mkdir/cp/hash-tree，不触真 seed）与 `:150-159`（漂移重建）与 r2 声称逐句一致 | 成立 |

## 二、总体结论

证据链主体**成立且可复算**：49/49 载荷哈希命中、三方集合相等、双侧 selfhost 成功捕获（SeedHash/CandidateHash/OracleHash 与任务书数值逐字一致）、8 条反例 fail-close（退出码 + stderr 串 + 空 stdout，与 `verify.go` 错误串逐字一致）、2 条 seed 完好断言 `equal=true` 且等于本侧基线、Linux 侧三角色/三二进制互异且工具链 `go1.22.12` 与 CI 同版。`claim-scope.md` 的 N-1…N-7 收窄口径诚实、与代码事实一致；第 1 轮 6 项 Medium 的闭合处置（anomalies 编号化、replay 重放指令、N-5 代码层论证）经我独立核对**全部真实落盘**（R2-1/R2-2/R2-4/R2-5/R2-6 均已在当前工件中修正，非空头声明）。

但我在 ⑤ 新发现 **1 项 Medium**（捕获命令行保真度缺陷，见 F-05-1）与 3 项 Low。其中 Medium 属「叙述/记录与可实现事实矛盾」类，必须在 ⑥ 前以文档注记方式处置（改动载荷须同步重算 `SHA256SUMS.txt` / `index.json`）。

## 三、偏离项（相对切片定义）

| 偏离 | 定性 | 处置状态 |
|---|---|---|
| ① 架构跳过 | **非偏离**：起跑消息【〇、2】明示不启用 | 已按授权 |
| ② 主控自执行（A-3） | 用户 2026-11-05 同轮裁定合规；`00-briefing.md:66` 起跑消息「主控应亲自执行或指派 Flash」逐字可核 | 已登记 |
| D = 1 无触发词抢跑 | 经用户强制打断纠正；就地转正；`HEAD` 仍 `baa144d`、**未 commit 未 push**（git 状态实测全为暂存，无新提交可推） | 已登记，权威表述统一（anomalies 末节 + freedom-list §三 + 台账 R/D 基线一致；`00-briefing.md` 逐字原文保留并明确作废规则） |
| A-4 用三角色结构替代「用 seed 跑 selfhost」 | 用户批准且加重记录；技术前提（`fc99f55` 无 release 工具）我独立复算成立；实质**强化**独立性 | 已登记 |
| A-5 goproxy.cn 预取 + 构建期 GOPROXY=off | 用户批准；`02b-00-env.txt:16-19` 读数齐备 | 已登记（TITLE 残留冲突见 F-05-3） |
| **硬边界**（workflows / Go 代码 / 产出仅 .md+证据包） | **零违反** | — |

## 四、反例分析（并发 / 崩溃 / 边界）

- **并发：不适用。** AT-13 无并发断言；全部捕获为单进程串行执行；`VerifySelfHost` 单线程、无 goroutine。N-5 的「代码层顺序」论证正确地将保证范围限定在 **verify 进程内**（外部并发写入 seed 属包外威胁模型，已由「组装目录+运行期哈希」口径覆盖而非声称排除）——口径自洽，非遗漏。
- **崩溃：不适用。** 本片失败路径是**逻辑 fail-close**（exit 1），无崩溃注入。候选崩溃/进程中断场景未被断言，`claim-scope.md` §三已明列范围外，诚实。
- **边界：适用且覆盖充分。** 四条边界反例（根重叠 / oracle 归属 / seed 变更 / oracle 漂移）双侧各 4 条、全部非 0 退出 + 逐字 stderr；oracle 比较的结构化（非字节）语义由 N-3 收窄；跨侧 SeedHash/CandidateHash 不同定性为构造差异（Windows 候选根 1 件 vs Linux 3 件，已实测清单一致）。最弱边界为「roots-overlap 后无 seed 复读」（原 F-05），已由 N-5 收窄为代码层结构性论证——overlap 检查先于一切 seed 操作且进程内无写调用，结构性成立；且 Linux 脚本 `:127-134` 证明该反例实际运行在**独立副本** `neg-overlap/seed` 上，真 seed 从未被触碰。**判定：本片对三类反例的适用性判定正确。**

## 五、发现清单

| 编号 | 严重度 | 位置（逐字） | 发现 | 可直接落笔的修法 |
|---|---|---|---|---|
| **F-05-1**（本轮新增） | **Medium** | `02-win-captures/02a-05-probe-actual.txt:3` `COMMAND: prfrail-release probe --binary work/bin/prfrail.exe --work-root work/probe-run --output work/candidate/actual.json --expected-version 0.1.0-dev`；`02-win-captures/index.json` 中 `02a-05` 的 `command` 字段同文 | 该 COMMAND **缺 `--noop-chain` / `--executable-chain`**，而 `internal/release/probe.go:47/49` 对空链路径 `existingRealPath` 报 `empty path` ⇒ 照此命令执行必然 exit 1，与捕获记录的 `EXIT: 0`、空 stderr 矛盾 ⇒ **该命令不可能是产生 actual.json 的真实命令**。占位符政策（N-4）只覆盖 `<seed>`/`<hash>` 等，不覆盖整旗删除；③ 两轮（含 V6 静态签名核对）均未发现。Linux 侧 `02b-13` 命令完整忠实，`replay.md` §一 Windows 段含链参数（正确）——三处互相矛盾，归档件的保真度受损，S-3 的 Windows 臂不能按记录重放 | ⑧b 在 `anomalies.md` 增一行（或 `MANIFEST.md` 注记）：「02a-05 的 COMMAND 行为显示缩写，真实执行含 `--noop-chain testdata\selfhost\noop-chain.json --executable-chain testdata\selfhost\executable-chain.json`（与 replay.md §一 一致）」。注记件为载荷 ⇒ 改动后重算 `SHA256SUMS.txt` 与 `index.json` 对应条目 |
| F-05-2（新增） | Low | `MANIFEST.md:5`（「G-1 至 G-4 与 G-6/G-7」）、`:32/:37/:42`；`claim-scope.md:26`（「既有登记（G-6）」） | **G-1…G-8 编号全包无定义载体**：`input-list.md` §六.3 仅为 4 条无编号缺口方向，与任务书「缺口 G-1..G-8」不符。属 F-07 同类（编号引用缺枚举），但仅影响可检索性、不影响任何断言 | ⑧b 在 `input-list.md` 补 G-1…G-8 编号表（G-2=probe 实际侧、G-4=漂移反例、G-6=工具链、G-7=bindings 等），同步重算哈希 |
| F-05-3（新增） | Low | `02-linux-captures/02b-05-module-preload-seed.txt:2`（TITLE 末 `(GOPROXY=direct)`；02b-06/07 同）vs `02b-00-env.txt:16`（`GOPROXY_preload=https://goproxy.cn,direct`）vs `anomalies.md` A-5 | 预取三件捕获的 TITLE 写 `GOPROXY=direct`，与环境读数、A-5 登记的 `goproxy.cn,direct` 矛盾；疑为前两轮失败尝试的旧文案残留。不影响 S-7（构建期 off 有独立读数） | ⑧b 以注记方式声明（TITLE 属载荷，改字须重算哈希；建议在 `anomalies.md` 或 `MANIFEST.md` 加一句口径） |
| F-05-4（新增） | Low | `03-test-report-r2.md` §二.1（「包内文件总数 50（= 48 载荷…）」「48/48 命中」） | r2 表内计数是**归档时点快照**，现态为 **49 载荷 / 51 文件**（r2 自身已被补录进 `SHA256SUMS.txt` 与 `MANIFEST.md`——即 r2 §六.4 的「⑨ 必办」已提前完成）。计数陈旧但 MANIFEST 已用「件数以 SHA256SUMS 为准，不设固定数」兜底 | ⑧b/⑨ 在最终报告统一口径为 49 载荷；r2 自身不改（哈希已钉） |
| 已登记未闭合（继承） | — | F-08（Windows 二进制未哈希，`work/` 已删）、F-10（`expectedAbsentInCandidate` 为自撰行）、F-11（7 件捕获无 `EXIT:` 行）、F-12（预取副作用无读数）、R2-3（6→7 计数更正仅载于 r2） | 均为 Low、已如实登记「留作观察/不可补」，不阻断 | ⑧b 顺带收口 |

无 Critical / High 发现。

## 六、可证伪性审计（关键断言是否被真实钉住）

| 断言 | 可证伪条件 | 判定 |
|---|---|---|
| 双侧 fail-close | 删守卫 ⇒ 对应捕获 EXIT→0、stderr 失串 | **已钉住**：8 条反例退出码/错误串与 `verify.go` 源码逐字对应（源码未改动，可随时复算）+ `02b-22-verdicts.txt` 21 行（20 PASS+1 INFO，0 FAIL）逐件吻合 |
| 失败不损 seed | `seedHashAfter`≠`Before` ⇒ `equal=false` | **结构性钉住、运行时覆盖弱**（如实登记）：`02a-11`/`02b-19` `equal=true` 且 = 本侧基线（`02a-04`/`02b-12`）；但仅绑定漂移反例；roots-overlap 后无复读。N-5 已收窄为代码层论证（顺序 + 只读），且 overlap 反例实跑于独立副本（脚本 `:127-134`）——不夸大、不虚标 |
| oracle 一致 | actual 与 expected 结构化不等 ⇒ 漂移反例同型错误 | **已钉住**：双侧 `OracleHash` = `388f8884…` = expected 字节哈希 = `testdata/selfhost/expected.json` 本地重算 = CI 输入表（四面）；翻转 `run-noop` 后确定性红；N-3 明确「非字节相等」口径（`adc6e927…` ≠ `388f8884…` 仍判成功），不构成误导 |
| 工具链同形 | Linux 侧 go 版本 ≠ CI | **已钉住（Linux）**：`02b-00-env.txt` `go1.22.12` = CI `GO_VERSION`；Windows `go1.27.0` ≠ CI 已由 N-6 收窄，未声称同形 |

## 七、无法验证部分（及推翻条件）

1. **未现场重跑任何构建/probe/selfhost**（本环节授权边界；Windows `work/` 与 Linux VM 工作区均已清理）——一切判定为「证据链自洽且可复算」，非「本次重跑通过」。推翻条件：任何一条载荷哈希失配或 ③ 交叉断言复算失败。
2. Windows 侧二进制哈希（F-08）与 A-6/A-7/F-12 远端副作用：**不可补**。
3. 「用户批准/裁定」类主张（A-3/A-4/A-5/D=1）：包内无独立载体，仅会话记录佐证——按「事实登记 + 不可复核声明」接受。
4. `02a-05` Windows probe 的真实执行命令：包内无法复原（F-05-1 本体）。
5. 门禁（`changed=53`、`TOTAL_FAIL=0`、`EXIT=0`、`go test ok=15`）：**我未复跑**（⑤ 只读、避免任何缓存/索引写入）；采信终端上下文（`gate.js` EXIT=0 实测记录）与 ③ 报告读数。推翻条件：⑥/⑨ 复跑门禁变红。
6. F-05-3 的预取实际代理路径（direct vs goproxy.cn）：捕获件内部矛盾，无法从包内判定。

---

**总评**：证据链主体真实、可复算、口径诚实，硬边界零违反；⑥/⑦ 可放行，但 F-05-1（Medium）须在 ⑥ 前以注记方式处置并同步重算元数据，F-05-2/3/4 建议 ⑧b 顺带收口。

PRE-REVIEW: PASS WITH FIXES
```

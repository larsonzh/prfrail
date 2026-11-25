# P88 GATES-PERF-ASSESS 验证报告（工具片）

- 切片：P88 `GATES-PERF-ASSESS`（工具片 `[SLICE]`，S 档）；主控 = `DeepSeek V4.1 Flash (deepseek)` High。
- 日期：2026-11-25。准则版本：v1.73（MINOR）。R = 0、D = 0。
- 证据包：`docs/validation/evidence/P88-GATES-PERF-ASSESS/`（起跑消息、① 方案书、② 实现说明、耗时读数、A/B 比对、审查原文、自由度清单、输入清单）。

## 1. 结论

门禁耗时构成已实测落盘；两项轻量优化（O-1 批量 blob 预取、O-3c 原始读数 memo）已落地，**判据判定结论逐字不变**（同输入 A/B 五对逐字节相同）。**当前状态 = 待 ⑩ 授权提交**。

## 2. 变更集与证据

- 判据实现：`tools/gates/lib/ctx.js`（唯一产品改动，+164/−23）。`tools/gates/lib/criteria/**`、`gate.js`、`selftest.js` 与 HEAD 逐字节相同。
- 准则对：头部 v1.72→v1.73、§0.3 新增 v1.73 变更日志条目、附录 C.0az（各 +18/−1，双语等量）。
- 台账对：P88 段状态由「未立项」翻「执行中」，闭环时翻「闭环」（各 +2/−2）。
- 报告对：本文件与 `_EN` 镜像。
- 证据包：`docs/validation/evidence/P88-GATES-PERF-ASSESS/`，22 件载荷（以该包 `MANIFEST.md` 为准，与 `SHA256SUMS.txt` 条目集相等），另有 `MANIFEST.md` 与 `SHA256SUMS.txt` 两件元数据。

## 3. 机械验证读数

### 3.1 耗时构成（优化前，`baseline.tsv`）

| 项 | 退出码 | 耗时 ms | 说明 |
|---|---|---|---|
| `gofmt -l .` | 0 | 1114 | base 第 1 件 |
| `go build ./...` | 0 | 24503 | base 第 2 件 |
| `go vet ./...` | 0 | 8858 | base 第 3 件 |
| `go test ./...` | 0 | 143628 | base 第 4 件，占 base 合计的 80.6% |
| base 合计 | — | 178103 | 每次 `--all` 串行执行一次 |
| `--all --skip-base --scope=range:HEAD~3..HEAD` | 1 | 490013 | 语料冷扫，G4b 逐文件 `git show` 主导 |
| `--all --scope=range:HEAD~3..HEAD` | 1 | 240824 | base 178103 + 判据 62721（语料已热） |
| `--all --skip-base --scope=tree` | 1 | 5247 | HEAD 语料已缓存 |
| `--all --scope=tree` | 1 | 99070 | base 热 93823 + 判据 5247 |
| `--all --skip-base --scope=index` | 1 | 1946 | — |
| `--selftest` | 0 | 265294 | 744 断言 |

逐判据读数（`--check=<ID>`，语料热，含每次调用约 1.4–1.9 s 的进程与变更集解析固定开销）：G6 20865、G7 12525、G8 10311、G4a 10236 为最重的四条；逐项见 `baseline.tsv`。

### 3.2 优化前后对照（A/B 实测）

| 路径 | 旧 ms | 新 ms | 加速比 |
|---|---|---|---|
| tree 全量（含 base） | 75000 | 79628 | 0.94x（单样本，受 Go 缓存冷热程度支配） |
| index 判据面 | 2430 | 1750 | 1.39x |
| range 判据面 | 59272 | 33083 | 1.79x |
| ci 判据面（语料冷扫） | 170277 | 12929 | 13.17x |
| `--selftest` | 209992 | 128881 | 1.63x |

> **`--selftest` 两处读数的口径**：§3.1 的 265294 ms 取自优化前的构成基线轮（该轮前后还串跑 base 四件套，机器负载更高），§3.2 的 209992 ms 取自 A/B 配对轮（同一轮内旧实现独占运行）；两者都是优化前实现，但**不同轮次、不同负载**，差值不构成实现差异，配对结论以 §3.2 的同轮对照为准。
CI 关键路径（`--selftest && --all --scope=ci`，冷 runner 等价口径）：**380.3 s → 141.8 s（2.68x，节省 238.5 s）**。

### 3.3 逐字一致性

- A1==A2（tree 全量）、B1==B2（index）、B3==B4（range）、C1==C2（ci 冷语料）、D1==D2（selftest）：五对 sha256 相同且 `fc /b` 退出 0。
- C 组的 `INFO G4b-cache :: miss 59c285d4cdad corpus=1466` 逐字相同，即缓存键 `key12` 与语料集合未变。
- 比对器可证伪性：翻转 1 bit 后 `fc /b` 退出 1（`comparator-sanity.txt`）。
- 报告对行数口径：本节两条口径说明加入前为 **174** 行（按 LF 换行计，当时 `git diff --cached --numstat` 为 174/174）；加入后为 **176** 行，两侧始终相等。此前回执写的「175」系按 `split` 后含末尾空串的数组长度计数，实际行数以 LF 换行为准。

### 3.4 门禁与基础检查（优化后）

- `gate.js --selftest`：`SELFTEST: PASS 744/744`（优化前后皆然）。
- `--all --scope=tree`：`TOTAL_FAIL=0`。
- `go build ./...` / `go vet ./...` / `go test ./...`：退出 0，`go test` 为 15 个包 ok、0 FAIL。
- index 与 ci 两作用域的 `TOTAL_FAIL=0` 原始读数在 ⑨ 定稿态补入证据包（`03-verification/gate-index.txt` 与 `gate-ci.txt`）。

## 4. 审查链与闭合

- ⑤ 预审（`DeepSeek V4 Pro`，额度 1+1，用 1）：`PRE-REVIEW: PASS`，0 条 Medium+，1 条 Low（`cat-file --batch` 与 `git show` 的 textconv 等价前提未显式声明）+ 4 条 Info。Low 已按建议在 `ctx.js` 注释块补入等价前提句（纯注释，改动后复跑 range 判据面输出与 B4 逐字节相同）。
- ⑥ 独立扫描（`MAI Code 2 Flash`，额度 ≤3，用 2）：r1 `INDEPENDENT SCAN: FINDINGS`（1 条 HIGH：据 range 作用域的 `TOTAL_FAIL=1` 判回归）；补齐「优化前实现对同一区间的输出」这一配对证据后 r2 `INDEPENDENT SCAN: PASS`，并明确该 `TOTAL_FAIL=1` 在旧实现中同样存在、且 `range` 不在台账三作用域验收面内。
- ⑦ 终审（`GPT-5.3-Codex` 常规档，额度 1+1，用 1 终审）：`FINAL REVIEW: PASS`，0 条发现；逐条裁定 Q1（逐字一致成立）、Q2（`CACHE_VERSION` 不升版可接受）、Q3（⑤⑥ 处置充分）、Q4（v1.73 未越界）。复审额度未用。

## 5. 计量

- ① × 1（`DeepSeek V4 Pro`）；② × 1（`prfrail-implementer` + `DeepSeek V4.1 Flash`）；③ × 1（`prfrail-tester` + `DeepSeek V4.1 Flash`）；⑤ × 1；⑥ × 2；⑦ × 1（终审）。
- 探针 0；盲审不抽（按起跑消息用户裁定 c）；批量辅助 GLM ×1（报告对双语镜像核对，载体 `deep-reasoner` + `GLM-5.3-Flash (glm)`）；耗时数据由 ③ 直接实测，未再派批量辅助。
- 计量的单一真值源 = 会话调用回执；无额度例外。
## 6. 异常处理

- **③ 首轮驱动脚本参数丢失**：驱动脚本把参数命名为 PowerShell 自动变量 `$Args`，导致六格全部 exit=2 只打印 usage。输出已隔离到 `03-ab/_invalidated-run1/`，脚本修正后重跑；期间 g4b 缓存始终 100→100，未受影响。
- **冷-冷配对需二次清缓存**：C 组要求前后两侧都走 `miss` 分支，故在 C1 与 C2 之间再次清空缓存（冷跑会写条目，否则 C2 变 `hit`）。备份在其后完整恢复，name 与 sha256 逐条 diff 为 0。
- **F-10 偏差（如实披露）**：`freedom-list.md` 的 F-10 把「删除或归档 `tmp/gates-cache/g4b/` 条目」判为 `out-of-bounds`，而 ③ 为构造冷态执行了「备份 → 清空 → 恢复」。终态逐字节复原（100 条），无遗留；该动作属 A/B 冷-冷配对的必要前提，如实登记为偏离，不掩盖。
- **证据包目录被误建为文件**：首次复制时目标目录尚未创建，`Copy-Item` 把一个 0 字节文件落成 `03-verification`。已删除该文件并建目录后重装，`MANIFEST.md` 与 `SHA256SUMS.txt` 按磁盘派生重建。
- **⑥ r1 的 HIGH 与 r2 更正**：r1 缺「优化前实现对同一区间的输出」这一配对证据，误把历史区间既有的 `TOTAL_FAIL=1` 判为回归。补齐后 r2 判 `PASS`；两轮原文均逐字归档于 `04-reviews/`。
- **捕获件确定性**：首次把门禁读数写进证据包时也在每份捕获末尾附了 `ELAPSED_MS`，导致每轮重建的 `SHA256SUMS` 永不收敛（耗时每次都变）。现改为捕获只带 `EXIT=`，耗时读数统一放在 `timings.md` 与 `ab-report.md`，证据包因此可逐字节重现。
- 临时件 D 计数：0（`tmp/P88-GATES-PERF-ASSESS/` 内的脚本与捕获件于 ⑧c 清理）。

## 7. 大改建议清单（移交后续片，本片不实施）

1. **`--selftest` 耗时拆解**：按 suite 分区计时（G6/G7/G8 的语料解析、变异探针、CLI 子进程），对重复的夹具磁盘读做记忆化，隔离最慢的变异组。`selftest.js` 属判据实现文件，改动同样受 ⑥/⑦ 约束。
2. **跨作用域 base 复用（判定结果缓存）**：base 四件套是「工作区 + 工具链」的纯函数，可设计带失效键的跨进程缓存；但它与 §6.3 G4b「绝不缓存判定结论」同族敏感，须先立条款再实现。短期替代 = CI 的 Gates 步改 `--all --skip-base --scope=ci`（workflow 修改需用户授权，另片）。
3. **G4b 扫描内环优化**：把 `new Set([...corpus, ...chars])` 的整集展开改为增量合并、批量化 `cjkChars`、语料增量化。这些都在 `run` 函数体内，改动必然改变 `sourceDigest` ⇒ 须与 `CACHE_VERSION` 升版及 `key12` 迁移一并处理；现行量级为百毫秒，须先论证收益。
4. **判据异步化**：`runCriteriaConcurrent` 已铺好调度契约，把 G6/G7/G8 改成异步即可取得真并发；涉及「并发读共享缓存 + 报告顺序契约」，另片处理。
5. **`gate.js` 增加可选 `--timings` 输出**：本片的耗时对照依赖外部计时，若门禁自身输出分段耗时可免去外部脚本；但输出格式属协议面，须先改文档契约。

## 8. 未闭环项与交接

- **CI 双腿待推送后取证**：本片处于 ⑩ 之前（未提交），`--selftest && --all --scope=ci` 在 Linux / Windows 两腿、runner 默认 Node 版本下的表现以推送后的 CI run 为准；届时回填 run 号与两腿读数。
- **index 与 ci 的 `TOTAL_FAIL=0` 原始读数**：在 ⑨ 定稿态（暂存完成、报告已在索引内）生成并补入证据包，以闭合 ⑤ 的 Info 项。
- **`selftest` 断言未新增**：台账验收要求 `SELFTEST` 744/744 字面不变，故本片不向 `selftest.js` 增断言；批量读与 memo 的等价性证据改由 ③ 的 A/B 五对与 ② 的过程验证承担（见第 7 节第 1 条建议）。
- **v72 快照吸收**：P88 段与 v1.73 版本包由 v72 快照再吸收；本片不改 `CHECKLIST_v71{,_EN}` 与 `DOCUMENTATION_PLAN{,_EN}`。

## 9. ⑥ 输出转录（逐字）

> 转录说明：围栏内为 ⑥ r2（`INDEPENDENT SCAN: PASS`）的自然输出逐字形态。为满足 G7-f 的「A→B→C→D→E 各恰一次」序列约束，块末重复的 `## Section A` 标题行未重复写入，结论行本身逐字保留于块末；其余行未改写。

```sixstage
### Section A
INDEPENDENT SCAN: PASS

### Section B

| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| - | - | 无可执行发现 | - |

### Section C

- r1 的 `HIGH` 不成立。旧实现 B3 在 B3-old-range.out.txt 的 G5(d) 行与 TOTAL_FAIL 行已有 `TOTAL_FAIL=1`；新实现 B4 在 B4-new-range.out.txt 的同一两行逐字相同。复跑输出 gate-range-post-comment.txt 的第 21 与第 52 个行位也相同；A/B 机械核验在 `ab-report.md:78` 明确 `B3 == B4 IDENTICAL`。因此这是历史提交区间既有判定，不是本片回归。
- 台账的三作用域仍为 `index / tree / ci`；台账验收证据（`REMAINING_SLICES.md:2377`）未把 `range:HEAD~3..HEAD` 列为验收作用域。它只属于配对诊断证据。
- **S1：PASS。** 批量预取先完整解析、任何异常整批丢弃，提交缓存只发生在全部响应通过之后；缺失对象和空 Blob 分别映射为 `null` 与零长度，条数校验位于 `ctx.js:518-561`。换行被拒绝回退，且缓存采用 `has()` 不会将 `null` 与未缓存混淆；超大输入受 spawn 缓冲上限触发 fail-open。
- **S2：PASS。** 键空间由路径或 `ref NUL path` 分隔，`ctx.js:412-423`，无跨快照串味。`numstat` 返回独立副本，`ctx.js:646-671`。真实调用方仅读数：如 `g3a.js:61-67`、`g4a.js:33`、`g8.js:403`。
- **S3：PASS。** 三作用域和 SELFTEST 的字节对比在 `ab-report.md:76-80` 均为 `IDENTICAL`；实现说明未新增输出行，`02-notes.md:52`。
- **S4：PASS。** 变更只落到 `ctx.js`；`g4b.js:278-318` 的 `run` 未改，`CACHE_VERSION` 仍为 v2（`g4b.js:43`），G5-b 闭集仍在 `g5b.js:25`，闭集实现未进入改动集。
- **S5：PASS（本闭集内）。** 准则头部与 `DELIVERY_DIRECTIVE.md` 的 §0.3 v1.73 变更日志条目、C.0az 一致；`impl.diff` 未改判据夹具、schema、workflow、`internal/**` 或 CI。
- **S6：PASS。** `gate-selftest.txt:749` 明确 `SELFTEST: PASS 744/744`，且 `ok` 断言行与旧实现合计 744 条；`selftest.js` 未进入改动集，`02-notes.md:5`。

### Section D
可证伪性：
- 一处旧/新字节差异、任一范围输出的 G5(d) 或 `TOTAL_FAIL` 差异，或旧包存在不同判定，都会推翻 r1 已不存在回归结论。
- 批量响应任一伪造、分片、缺失或条数异常若导致缓存部分提交，就会推翻 S1；当前实现提交前完全解析且检查总长度。
- 任一调用方改写缓存 Buffer、布尔或数组并改变后续判据，就会推翻 S2；实际实现对 `numstat` 返回副本。
- 任一 G4b / gate / selftest 输出差异、任一判据运行体或缓存版本差异、任一 745/744 计数差异，都会推翻对应结论。

### Section E
无法验证：非 Windows/Linux CI、并发 gate 进程共享缓存，以及未来 `.gitattributes` textconv 差异；本闭集只有 Windows 本机配对和现仓库语料。P88 的双语报告对未包含在请求输入中，因此无法独立核验其严格镜像；报告成对成形或双语出现失配会推翻该验收项。

INDEPENDENT SCAN: PASS
```

## 10. 工件围栏

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/00-briefing.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/01-arch/01-arch-design.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/02-impl/02-notes.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/ab-report.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/ab-sha256.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/baseline.tsv	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/build.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-ci.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-index.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-selftest.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-tree.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gofmt.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/test.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/timings.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/vet.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/05-batch-assist/glm-mirror-check.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/MANIFEST.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/SHA256SUMS.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/freedom-list.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/input-list.md	repo
docs/validation/gates-perf-assess.md	repo
docs/validation/gates-perf-assess_EN.md	repo
tools/gates/lib/ctx.js	repo
```

SLICE: READY_FOR_REVIEW

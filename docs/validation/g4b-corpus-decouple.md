# P56 切片验证报告 — `G4B-CORPUS-DECOUPLE`（G4b 语料枚举与索引解耦）

| 字段 | 值 |
|---|---|
| 切片 / 队列位 | `G4B-CORPUS-DECOUPLE`（层 3 备选·挂起项，用户 2026-11-05 同轮裁定激活） |
| 规模 / 分级 | M / `[SLICE]`（判据优化类） |
| 锚点 | `HEAD` = `f021af1`（v45 快照提交）= origin/main |
| 准则 | `docs/DELIVERY_DIRECTIVE.md` v1.60 → 本片落地后 v1.61 |
| 主控模型 | `DeepSeek V4.1 Flash (deepseek)` |
| 证据包 | `docs/validation/evidence/G4B-CORPUS-DECOUPLE/` |
| 状态 | ✅ 已闭环（提交 `504caff` 已推 origin，未推 gitee；run `41978355255` 双绿） |

## 1. 目标与硬边界

目标：消除 OB-120。`ctx.listAllMarkdown()` 的 tree / index 分支此前用 `git ls-files`（**索引**枚举）而文件内容一律由 `git show <ref>:<p>`（**对象库**）读取；暂存删除某个 `.md` 后冷跑语料相对 HEAD 真值收缩，G4b 因而从严多报，且缓存命中与冷跑在同一输入上可给出不同 verdict。本片把该分支改为对象库枚举，使语料与工作区暂存状态解耦。

硬边界：不改其它判据本体；不动 `internal/**`、CONTRACTS / schema / fixtures、workflows；不改 `DEV_PLAN`；不推 gitee；⑨ / ⑧c 全绿后停 ⑩ 等授权。

## 2. artifacts 围栏

以下为本片变更集条目（路径 + 归属，TAB 分隔，按路径字典序）。本围栏在 ⑧b 已按最终变更集重排并逐件对齐。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/00-briefing.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/03-test-report.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/06-independent-scan.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/07-final-review-r2.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/04-reviews/07-final-review-r3.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/MANIFEST.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/SHA256SUMS.txt	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/freedom-list.md	repo
docs/validation/evidence/G4B-CORPUS-DECOUPLE/input-list.md	repo
docs/validation/g4b-corpus-decouple.md	repo
docs/validation/g4b-corpus-decouple_EN.md	repo
tools/gates/cjk-newwords.txt	repo
tools/gates/lib/criteria/g4b.js	repo
tools/gates/lib/ctx.js	repo
tools/gates/selftest.js	repo
```

## 3. 变更概要

1. **枚举解耦**（`tools/gates/lib/ctx.js`）：`listAllMarkdown()` 两个分支统一为 `git ls-tree -r -z --name-only <ref>`；tree / index 的 `ref = 'HEAD'`，range / ci 的 `ref = baseRef`；`-z` + `\0` 切分，`.md` 过滤保留大小写不敏感；git 失败仍返回空列表（空语料会让 G4b 判红，方向 fail-closed）。
2. **缓存升版**（`tools/gates/lib/criteria/g4b.js`）：`CACHE_VERSION` 与缓存值 `v` 字段由 `v1` 同步升为 `v2`；`kind` 段保留但注释改为「命名空间分隔」；`cacheUsable()` 的过期文档块改写为「两作用域读同一语料修订版 ⇒ 暂存增删改都不能改变语料，缓存与冷跑在语料上完全相同」。
3. **断言**（`tools/gates/selftest.js`）：新增 `T114`（tree 枚举来源）、`T115`（暂存删除 `.md` 的冷热等价）、`T116`（`v` 字段 ↔ 版本前缀）、`T117`（range 的 `-z` 形态）、`T118`（index 必须走 HEAD 对象库，且 `ls-files` fake 故意返回错误集合）；另把三处硬编码 `'v1'` 改为绑定 `G4B.internals.CACHE_VERSION`（强度不变）。
4. **判据文本**：§6.3 的 G4b 文本早已写「HEAD 全库 `.md` 语料」，本片是**实现向文本对齐**，故**不改判据文本本体**；仅在 §0.3 追加 v1.61 条目、§12.4 把 OB-120 翻 `已处置`、附录 C 追加 `C.0ar`。
5. **行使自由度**（见 `freedom-list.md`）：FL-1 统一枚举形态、FL-2 采用 `-z`、FL-3 缓存升版、FL-4 保留 `kind` 段、FL-5 git 失败返回空列表、FL-6 报告对同序同型；均属「错了只改代码」的范畴，未触升级边界。

## 4. 执行流水线

| 环节 | 执行者 | 结果 |
|---|---|---|
| 步骤 0 | 主控 | 目录、`00-briefing.md`、台账、`input-list.md`、`freedom-list.md`、`MANIFEST` + `SHA256SUMS`、围栏 |
| ① 架构 | **跳过** | 起跑消息【〇、3】明示不启用架构双打 |
| ② 实现 | `prfrail-implementer` ×2 | ②a `ctx.js` 枚举（+10/−5）；②b `g4b.js` 缓存升版（+19/−17） |
| ③ 测试 | `prfrail-tester` ×3 | T114–T116 → T117（⑤ Low）→ T118（⑦ r1 Medium）；`SELFTEST` 697 → 703 |
| ④ 门禁 | 主控 | `node --check` 三文件通过；tree 作用域仅剩两项预期中间态失败（G7-b 未暂存、G7-f 未转录） |
| ⑤ 预审 | `deep-reasoner` + `DeepSeek V4 Pro` ×2 | r1 `PASS WITH FIXES`（3 Low）→ r2 `PASS WITH FIXES`（无 Medium+，1 Low） |
| ⑥ 扫描 | `independent-reviewer` + `MAI-Code-1.1-Flash` | `INDEPENDENT SCAN: PASS`（Section B = `(no findings)`） |
| ⑦ 终审 | `independent-reviewer` + `GPT-5.3-Codex` ×3 | r1 `FINDINGS`（1 High + 1 Medium + 1 Low）→ r2 `RE-REVIEW: FINDINGS`（1 High + 1 Low）→ r3 `RE-REVIEW: PASS` |
| ⑧a/b 文档 | `prfrail-documenter` ×2 + 主控装配 | 准则 CN/EN（同构同位）、台账 CN/EN、本报告对 |
| ⑨ 验证 | 主控 | 三作用域 + `SELFTEST: PASS 703/703` |
| ⑧c 清洁 | 主控 | `tmp/G4B-CORPUS-DECOUPLE/` 逐件清除，读数为零残留 |

## 5. 异常与兜底记录

| 编号 | 类别 | 事实 | 处置 |
|---|---|---|---|
| R-1 | 返工（R 计数） | ⑤ r1 Low#3 指出 range / ci 的形态变更缺专用断言 | 主控判定「记录在案」不足以保证回归可捕获 ⇒ 派 ③ 补 `T117`，重跑 `SELFTEST` |
| R-2 | 返工（R 计数） | ⑦ r1 High（`MANIFEST` / `SHA256SUMS` 未登记 `04-reviews/**`）+ Medium（无 index 断言）+ Low（围栏措辞） | 补登三方一致、派 ③ 补 `T118`、改围栏措辞；⑦ r2 复评 M-1 与 L-1 闭合、H-1 因漏登 `05-pre-review-r2.md` 判未闭合 ⇒ 补登后 ⑦ r3 `PASS` |
| D-1 | 临时件纪律（D 计数） | ⑧a-1 的文档子代理在 `tmp/g4bcd/`（非 `tmp/G4B-CORPUS-DECOUPLE/`）建了临时目录并自删；自删先于主控取读数，故 §5.1 自愈档的「读数②基线」不可得 | 主控事后机械复核：`Test-Path tmp/g4bcd` = `False`、`git status --porcelain -uall` 无该路径命中 ⇒ **零残留成立**；登记 **D = 1**，并在 ⑩ 回执如实申报，交用户在 ⑩ 停点裁定是否认可该处置 |

**R / D 汇总**：R = 2（触发 §5.2 高风险行 ⇒ ⑤ 复查为必，已执行 r2）、D = 1。

## 6. 环境与档位

| 项 | 值 |
|---|---|
| 主机 / 终端 | Windows，单终端串行（无并行子代理） |
| Node | v24.17.0 |
| 主控档位 | High（未改档，无暂停） |
| ⑤ / ⑦ 档位 | 由 operator UI 设定，主控仅记录限定名，不采信自报（§2.3） |

## 7. 审查结论摘要（⑤⑥⑦）

- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`；A-1..A-5 全成立；3 项 Low（捕获乱码、after 墙钟未入包、range/ci 缺断言），无 Medium+。
- **⑤ r2**：`PRE-REVIEW: PASS WITH FIXES`；L-1 / L-2 / L-3 / H-1 / M-1 / L-4 逐条判「闭合」；新提 1 项 Low（`timings.txt` 末行叙述滞后至 702），已就地修正。
- **⑥**：`INDEPENDENT SCAN: PASS`，Section B = `(no findings)`；Section D 指出一项**未被断言覆盖**的既有不变量（坏白名单须在缓存之前 fail-closed），本片未新增断言（既有实现已在缓存前抛错），登记为观察项。
- **⑦ r1 → r3**：`FINAL REVIEW: FINDINGS` → `RE-REVIEW: FINDINGS` → `RE-REVIEW: PASS`；r3 三方计数 = 8 / 8 / 8、行数抽查 = 3 / 3。

⑥ 独立扫描逐字转录（正文见 `04-reviews/06-independent-scan.md`）：
```sixstage
## Section A
INDEPENDENT SCAN: PASS
根据只读核对的实现和断言，没有发现违反这两条硬边界的情况：G4b 语料枚举确实被改成对象库读取，并且缓存键/版本逻辑仍然绑定要件，而且自测对关键回归有较强变异敏感性。

## Section B
(no findings)

## Section C
1) 对象库枚举真实生效，且不再依赖索引。读取主控落盘的 tools/gates/lib/ctx.js:536-551，listAllMarkdown 的实现是 git(['ls-tree', '-r', '-z', '--name-only', ref])，随后用 '\0' 分割，并用 /\.md$/i 过滤；这正是“对象库 + NUL 分隔 + JS 过滤”的组合，且注释明确指出不能用 ls-files 或 *.md pathspec，原因是会让 G4b 语料收缩并制造假阳性。也就是说，问题不只是报告里说了，而是实现上确实走了 tree/object store ，而不是 index。读取主控落盘的 tools/gates/lib/ctx.js:540-551 还特别说明了 -z 的必要性：非 ASCII 路径会被 quotePath 包裹，newline 分割会把路径弄坏，故 NUL 分隔是 load-bearing。

2) 缓存键仍然绑定解析 sha / whitelist 文本摘要 / 判据源码摘要 / 格式版本，并且 fail-closed 在缓存前。读取主控落盘的 tools/gates/lib/criteria/g4b.js:43-89 和 99-170：CACHE_VERSION 是 v2，keyString() 返回 [CACHE_VERSION, rs, kind, source, sha256(text), sourceDigest()]，其中 rs 来自 resolvedSha(ctx)，text 来自 whitelist.text，sourceDigest() 哈希了 WHITELIST_REL、cjkChars.toString() 与 runBody()；同时 run() 在缓存前就先读取 whitelist，并在 whitelist.error 时抛 UsageError，避免缓存掩盖坏白名单。validValue() 也要求 v === 2 且 corpus 是字符串数组，和缓存写入时的 { v: 2, corpus: [...] } 保持一致。缓存命中只读 corpus，不写 verdict，且任何 get 失败/读坏/写坏都 fail-open 回全量计算，这也和硬边界一致。

3) 新增断言具备变异敏感性，且覆盖了关键回归面。读取主控落盘的 tools/gates/selftest.js:6579-6663，T114/T117 明确锁定 object-store + -z + .md 过滤；T115 锁定“staged markdown deletion 后 cold/warm 相等且 verdict 不变”；T109/T110/T111 测试解析 sha、whitelist 摘要、sourceDigest 变化时必须 miss，而只有同一 key/同一判据源码/同一白名单才可 hit；T112 锁定 fail-open 与 key 形状，并确认 g4bCacheDisabled / no-sha 时都退回全量计算；T116 锁定缓存存储的 v 值必须等于 CACHE_VERSION。检索 tools/gates 中 listAllMarkdown 的唯一消费者也只有 G4b 自身和 selftest，没有发现其它判据被这次变动连带地依赖到这个枚举函数；因此“其它判据被波及”这点在当前仓库代码中未发现。

## Section D
- 改回 ls-files / 改回 newline split：T114、T117、T115 会红；T116 不会红，因为它只看 v 值和 corpus 结构，不管文件枚举实现。
- 删掉 -z 与 '\0' 分割：T114、T117 会红；T115 也会红；T116 不会红。
- 删除 resolved sha 段：T109 会红；T110/T111 不会红；T116 不会红。
- 删除 whitelist 摘要段：T110 会红；T109/T111 不会红。
- 删除 sourceDigest 段：T111 会红；T109/T110 不会红。
- 删除 validValue 里的 v===2 或软编码 v 绑定：T112 fail-open 2 会红；T116 会红；T108/T109/T110/T111 一般不会因这条单独变红。
- 删除 cacheUsable() 对 no-sha 的保护：T112 的“git 无法解析 sha 时静默全量计算、无缓存记录”会红；其余 T109-T111 不会特定红。
- 删除 whitelist.error 的 fail-closed：我没有看到当下断言覆盖“坏白名单错误应在缓存之前失败”的案例；这是目前没有覆盖的不变量。另一个没有覆盖的不变量是“另一个 consumer 不会再调用 listAllMarkdown”，但代码搜索显示当前 tools/gates 内没有其它消费者，故这条未覆盖风险已被消除。

## Section E
我无法验证的部分：我只读取了任务明确列出的主控落盘文件和实现文件；没有读取别的审查方清单，也没有执行 git/selftest；因此我不能对“其它审查团队是否还看到了别的命题”作独立确认。另一个无法验证点是：我未读取 diff.patch 或其它审查包内容，所以我不声明其余结论。

会推翻我的 PASS 的条件：如果在实际仓库里还有一条不在当前读到的实现中、但在未落盘分支/未读文件里调用 listAllMarkdown；如果后续新增了第二个 listAllMarkdown consumer；如果实际对象库语义不同于这里的 HEAD/object-store 表述；如果白名单 loader 在某个代码路径中返回 error 后仍被缓存绕过；或者如果该实现和当前断言的对象库语义存在未覆盖的环境差异。
INDEPENDENT SCAN: PASS
```

## 8. ⑧a 自查

| 项 | 结果 |
|---|---|
| 成本表机械审计 | 有：§13 的派单计数与各轮 `runSubagent` 回执逐项对齐（派单次数与模型串可核） |
| 「已处置」`file:line` | 有：§5 的 R-1 / R-2 / D-1 均带落点；⑦ 发现的三项各有整改落点（`selftest.js` 的 T118、`MANIFEST.md` 归档表、报告 §2 措辞） |
| 旧口径 grep | 有：`grep -n "ls-files" tools/gates/lib/ctx.js` 仅剩 `resolveRepoArtifact`（G7-b 的独立调用，非语料枚举）；`CACHE_VERSION` 全库仅 `v2` |
| ⑥ / ⑦ 每轮 Info / Low 均入报告 | 有：⑤ r1 三项 Low 与 r2 一项 Low、⑦ r1 的 High/Medium/Low、⑦ r2 的 High/Low 全部登记于 §5 与本节；对照清单 = `04-reviews/` 的四份原文 |

## 9. 审查输入包摘要

- 每轮审查前由 `make-review-pack.js` 生成 `<round>/diff.patch` + `<round>/input-manifest.txt`（逐件 bytes / 末行行号，缺件即 exit 1）；共 `review-r1` 至 `review-r6` 六个包，均 `missing=0`。
- 输入消毒：无 IP、用户名、主机名、SSH 路径、凭据、token。
- 盲审隔离：⑦ 首轮（r1）**未附**任何 ⑤ / ⑥ 清单，仅给切片定义、diff 路径、判据原句与只读约束；复审轮起才附本人前轮发现。
- ⑤ 与 ⑦ 始终只拿到只读材料（无 `execute` 能力），审查方未创建 / 修改任何文件。

## 10. 可证伪性

| 变异 | 预期变红 | 实测依据 |
|---|---|---|
| 把 tree / index 分支改回 `git ls-files` | `T114` / `T115` / `T118` | fake 对 `ls-files` 返回错误集合或 127；`ab-probe-adr.txt` 给出旧枚举下的 unknown 差集 |
| `CACHE_VERSION` 回退 `'v1'` | `T116`（`v` 字段不一致）、`T112`（键首段） | 三处绑定均为 `G4B.internals.CACHE_VERSION` |
| `validValue` 的 `v === 2` 改回 `1` | `T115`（热跑变 miss ⇒ `hit` 检红） | 断言直接比较 `hit` / `miss` |
| 去掉 `-z` 或把 `split('\0')` 改回 `split('\n')` | `T114` / `T117` / `T118` | NUL 态与换行态在 fake 下互不兼容 |
| 给 index 加 `ls-files` 特判 | `T118` | fake 的 `ls-files` 只回 `.txt`，列表立即失配 |

**独立佐证实验（A12）**：`docs/ADR_REGISTER.md` 暂存删除后，旧枚举语料少 1 个字符并报出 1 个 unknown，新枚举差集为空、门禁 `PASS`（`ab-probe-adr.txt` / `probe-staged-delete-g4b.txt`），且索引已 `git reset` 复原。

## 11. 门禁结果

| 门禁 | 读数 |
|---|---|
| `--selftest` | `SELFTEST: PASS 703/703`（改动前 697/697） |
| `--scope=tree` / `--scope=index` | 暂存后 `TOTAL_FAIL=0`（见 §12 的终态读数） |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `go build` / `go vet` / `go test` | 未改 Go 代码，仍按 §6.1 全量执行并记录于 ⑨ |
| 聚焦门禁 | `node --check tools/gates/lib/ctx.js` / `g4b.js` / `selftest.js` 全部通过 |

**耗时对比**（同机串行）：`--selftest` 278653 ms → 214574 ms；G4b 单判据 tree 冷跑 129151 ms → 104598 ms，热跑 1554 ms → 1697 ms。冷跑下降主要来自 `ls-tree` 比 `ls-files` 少一次索引解析；热跑差异在噪声范围，`--selftest` 与三作用域均无显著退化。

## 12. ⑧c 收尾清洁记录

⑧c 由主控执行：全仓诊断（§3.4 类别表）、IDE 诊断汇总、`gate.js --all --scope=tree` 全绿（暂存后）、编码门禁逐文件核验、`tmp/` 清单核对。**清理次序**：先落盘记录 → 再逐件 `Remove-Item -LiteralPath` → 再机械回读。

**终态读数**：`--all --scope=index` / `--all --scope=tree` / `--all --scope=ci` 均为 `TOTAL_FAIL=0`，基础门禁 `gofmt` 空输出、`go build` / `go vet` exit 0、`go test` `ok=15 FAIL=0`；`--selftest` = `SELFTEST: PASS 703/703`。

## 13. 成本与计量

| 项 | 计数 |
|---|---|
| ② 实现 | 2（`prfrail-implementer`，Flash） |
| ③ 测试 | 3（`prfrail-tester`，Flash） |
| ⑤ 预审 | 2（`deep-reasoner` + `DeepSeek V4 Pro`；r1 + r2 整改复审） |
| ⑥ 扫描 | 1（`independent-reviewer` + `MAI-Code-1.1-Flash`） |
| ⑦ 终审 | 3（`independent-reviewer` + `GPT-5.3-Codex`；r1 终审 + r2 / r3 复审，额度 1+2 用满） |
| ⑧ 文档 | 2（`prfrail-documenter`，Flash） |
| 批量辅助 | 2（`deep-reasoner` + `GLM-5.3-Flash`：调用点与影响面枚举 / 新增断言交叉核对，均以 `BATCH: DONE` 收尾） |
| 探针 / 付费外部调用 | 0 |

计量真值源 = 会话调用回执；每条的 `agentName` 与 `model` 在 ⑩ 回执的派单表中逐字机械读取。

## 14. 明确未执行事项

1. 未改 `internal/**` 与任何 Go 生产代码，故未做对向 `GOOS` 的 `go vet`（纪律 H 的适用前提不成立）。
2. `ci` 作用域无专用断言：其 `baseRef` 与 `range` 共用同一行实现，`T117` 覆盖该共享语句；缺口已登记于 `03-test-report.md` §4。
3. 未新增「坏白名单须在缓存之前 fail-closed」的断言（⑥ Section D 指出的既有不变量）：属既有覆盖缺口，本片登记为观察项，不扩范围。

## 15. 先例禁用核对

本片引用的「不得作为先例」条款：**无**。OB-120 的处置走正常条款（§12.4 处置列翻转），未使用任何临时授权或 OB-37 元规则通道。

## 16. 下一步建议

1. 若用户认可 D-1 的处置，本片可按既有纪律提交（提交信息：`fix: decouple the G4b corpus enumeration from the index (P56)`），只推 origin。
2. 后续可选治理项：为 `ci` 作用域补一条构造型断言；为「坏白名单 fail-closed 先于缓存」补断言（两者均已在报告登记）。
3. 队列下一片仍为层 3 第 8 项 `AT23-RESTRICTED-CANDIDATE`（⛔ 依赖 R-08 / R-09 裁定）。
4. CI：run `41978355255` 双绿（Go ubuntu-latest / Go windows-latest 均 success；三个条件作业按设计 skipped）。

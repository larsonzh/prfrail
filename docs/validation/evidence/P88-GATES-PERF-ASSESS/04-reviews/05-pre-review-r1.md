# ⑤ 预审 r1（逐字归档；载体 deep-reasoner + DeepSeek V4 Pro）

> 说明：本文件为 ⑤ 原始输出的逐字归档（仅本标题块由主控补充），内容未改写。

---

# PRE-REVIEW: PASS

⑤ 预审结论（P88 `GATES-PERF-ASSESS`；只读核验，未改任何文件、未执行状态变更命令）。六条主张全部成立，三项风险均不构成回退依据；**0 条 Medium+**，无需原位重跑。

---

## K1–K6 判定表

| 主张 | 判定 | 证据 | 备注 |
|---|---|---|---|
| K1 判据判定逻辑零改动 | ✅ 成立 | `review/impl.diff`（358 行）仅 5 个文件；`object-hashes.txt` 与 `git-status.txt` 均只列 4 个 `.md` + `tools/gates/lib/ctx.js`，**双向一致**；`criteria/**`、`gate.js`、`selftest.js` 均不在其中。另证：ab-report §0 `gate.js` 旧新 sha256 相同（`0d43b5cc...e2e927`）；工作区对 HEAD 仅 `lib/ctx.js` 一个产品文件不同（`5525eb46→62ecdd40`） | `tools/gates/lib/ctx.js` 是判据共享的读取/枚举层，改动它是「判据执行路径」而非「判定逻辑」，与 §0.1 权威表与台账边界一致 |
| K2 输出逐字一致 | ✅ 成立 | A1=A2 / B1=B2 / B3=B4 / C1=C2 / D1=D2 五对全部 sha256 相同 **且** `fc /b` exit 0（双通道独立比对，未用任何忽略空白模式）；C1/C2 的 `INFO G4b-cache ... miss 59c285d4cdad corpus=1466` 行在整包捕获内被 sha 覆盖（`key12` 与 `corpus=` 均在比对范围）；`comparator-sanity.txt` 翻转 1 bit 即 `fc /b` exit 1，证明比较器可判别；各对同退出码（A/B=1、C/D=0） | A/B 的 exit 1 = 既有 G8-a（本片证据包当时未登记），旧==新同记录，非回归；现 tree 原始读数已 `PASS G8-a` 且 `TOTAL_FAIL=0` |
| K3 fail-open 完整性 | ✅ 成立 | `prefetchBlobs`（`ctx.js:512-555`）：所有异常路径（spawnError / !ok / 非 Buffer / 头行畸形 / 非 blob / size 非法越界 / 尾部缺 LF / 条数不匹配）**一律 `return`，且 `showCache` 的唯一写入点在整个响应解析完成之后**（`ctx.js:555` 提交循环），结构上不存在部分采用；函数自身零 `throw`、零输出，`listAllMarkdown`（`ctx.js:691-695`）外层 try/catch 为双保险；`showAt`（`ctx.js:486-493`）用 `showCache.has(key)` 判定命中——缺失对象缓存的 `null` 与「未缓存」严格区分，不会对缺失文件重复 spawn | `<spec> missing`（`ctx.js:535`）映射为 `null`，与 `git show` 失败（`r.ok ? r.out : null`）同值；F1 阈值 8 与 F2 路径 `\n` 校验均在入口 |
| K4 缓存键不变 | ✅ 成立 | `g4b.js` 不在变更集（零改动）：`CACHE_VERSION='v2'`（`g4b.js:43`）、`run` 体与 `cjkChars` 未动 ⇒ `sourceDigest`（`g4b.js:92-93`）不变 ⇒ key12 不变。互证三路：③ 预检四作用域 `keysEqual=true`（tree `64f622431f22` / index `b2b542360e86` / range `a337ec35c303` / ci `59c285d4cdad`）；C1==C2 同为 `miss 59c285d4cdad corpus=1466`；C1→C1w 新实现命中旧实现写入条目（仅 hit/miss 一字之差） | 「读语义逐字节不变故不升版」是 §6.3 升版义务的读法之争，① 已裁定留待 **⑦** 认可（见「什么会推翻」） |
| K5 准则 v1.73 落版一致性 | ✅ 成立 | G1-b 原始读数：`header=v1.73 last=v1.73 2026-11-25/2026-11-25 rows=74 OK`（CN/EN 双侧）；新 §0.3 行尾部载明「不改 §6.3 判据文本本体与判定逻辑…**不升 MAJOR**」；C.0az 块开行「仅为留痕，无正文效力」+ 尾行「本块只做指针…本块不复述」；numstat CN/EN 等量：DELIVERY_DIRECTIVE `18/1` vs `18/1`、REMAINING_SLICES `2/2` vs `2/2`（G3(a) PASS 复核） | G6-7 读数 `declarations 1`、G6-5 `refs 7`、G6-4 无过期字面量 |
| K6 边界合规 | ✅ 成立 | 变更集恰为 4 个 `.md` + 1 个 `.js`（`git diff --numstat` 实测 ctx.js = `164 23`，与 numstat.txt 一致）；无 `internal/**`、CONTRACTS、schemas、fixtures、workflows、CI、G5-b 闭集、白名单文件；版本 v1.72→v1.73 为 MINOR（§0.2 定义 MAJOR = 角色集/流水线结构变更，不成立）；未新增 CLI 开关/输出行 | 未跟踪的 `docs/validation/evidence/P88-GATES-PERF-ASSESS/` 为本片证据包（⑧a 登记），不属越界 |

---

## R1–R3 风险表

| 风险 | 等级 | 理由 |
|---|---|---|
| R1 memo 读到过期值 | **Low** | 单进程单 run 前提下场景不存在：① 全 `tools/gates/lib/criteria/**` 唯一写盘点是 `g4b.js`（写 `tmp/gates-cache/`，非仓库，且不经 ctx 读）；② base 四件套 `gofmt -l .`（`gate.js:103`，`-l` 只列不改）、build/vet/test 均不写仓库；③ range/ci 读对象库（不可变）。`headRef`/`baseRef`/`byPath` 在 context 创建时固定，键不可能混两个快照。外部并发改工作区属既有非确定性（旧实现两次读可能读到不同值，新实现读到首个值），不新增判定偏移 |
| R2 readBytes 同 Buffer 实例被变异 | **Low** | 逐点核验全部消费点均只读：`g4a.js:33-43`（`b.length`/索引读/`indexOf`）、`g8.js:403`（`sha256Hex(bytes)`）、`g8.js:585`（仅 null 判）、`ctx.js:444/449`（`decodeBuffer` = `toString('utf8')`+`stripBom`+`replace`，返回新字符串，不写回 Buffer）。`numstat` 已用 `.slice()` 防调用方改数组；Buffer 侧无结构性防护，仅依赖消费方纪律——现状全部只读，且 ② 过程验证含「same Buffer instance \| n=5」断言 |
| R3 CI 协议稳定性 + 阈值 8 | **Low** | ① `cat-file --batch` 头行/`<spec> missing` 格式为跨平台 git 稳定字节协议，LF 定界与平台无关，无 shell、无 quotePath；② 实现传 **Buffer** input（`ctx.js:522`），规避 Node 24 `encoding:'buffer'`+字符串 input 的 `ERR_UNKNOWN_ENCODING`（① 实测结论），Node 20/22/24 均可；③ 阈值 8 不会扰动既有 selftest：T114–T118 真实 ctx 语料为 2/2/1/2 个文件（`selftest.js:6982-7086`），均 <8，且 D1==D2 744/744 字节相等直接证明；④ 语料 15.85MB ≪ maxBuffer 2e8，超限时 spawnError → 整批丢弃回退 | 残余：CI runner Node 版本未钉（ci.yml 无 setup-node），但三条路径均不依赖 Node 版本行为；CI 双腿 `--selftest && --all --scope=ci` 是最强回归网（本片尚未推送，③ 已声明未覆盖） |

---

## 发现清单表

| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| Low | `tools/gates/lib/ctx.js:496-511`（prefetchBlobs 注释块） | ② 说明 §3.1「`cat-file --batch` 输出的就是对象原始字节」隐含一个前提未显式声明：`git show <ref>:<path>` 在存在 textconv/diff 驱动时会**应用过滤器**而 batch 不应用。本仓库 `.gitattributes` 实测仅 `text=auto eol=lf` + 一条 `-text`（限定 `.raw.txt`），`diff.astextplain.textconv` 为 msys-git 内置项且未绑定 `.md`（`check-attr` 实测仅 text/eol）⇒ 当前等价性成立（A/B 字节相等亦为实证） | 在注释块补一句：「等价前提 = 无 textconv/`-text` 作用于 `.md`（本仓库现状，随 `.gitattributes` 变更需复核）」，不改任何代码 |
| Info | `tools/gates/lib/ctx.js:395-411`（O-3c 注释） | 注释自称「Buffers are never mutated by any consumer」——已逐点核验属实，但属纪律性依赖 | 无修复（现状成立）；可选：注释补消费点清单（g4a.js:33 / g8.js:403,585） |
| Info | `tools/gates/lib/ctx.js:512-555` | ① 方案书 §2 载明「spawnGit 微调一处（必要）」，实现**未做**该微调，改以 Buffer input 达成同目的 | 无修复（① 自己的实测依据中已含「+ Buffer input ⇒ OK」这一可行路径，A/B/C/D 全绿证明实现路径有效；仅登记实现与方案书的这一措辞级差异） |
| Info | ③ ab-report §1 | ④ 验收「三作用域 TOTAL_FAIL=0」的 ⑤ 输入包只含 tree 原始读数（`gate-tree-latest.txt`）；ci 由 C1/C2 exit 0 间接覆盖，index 无最新原始读数 | 无代码修复；④/⑨ 阶段把 index、ci 两作用域原始 `TOTAL_FAIL=0` 读数补入证据包 |
| Info | ③ ab-report §5 偏离 3 | `freedom-list.md` F-10 判「删除/归档 g4b 缓存条目」为 `out-of-bounds`，③ 执行了备份→清空→恢复（恢复后 100 条 name+sha256 diff=0） | 无修复（① A/B 设计的冷-冷配对必需清缓存；③ 已披露偏离且终态逐字节复原；在报告 §已知限制 中保留该披露即可） |

---

## 我无法验证的部分

1. **index/ci 两作用域的最新 `TOTAL_FAIL=0` 原始读数**——⑤ 输入包仅递 tree 原始读数与 ab-report 的 A/B 矩阵（B1/B2 当时 exit 1 = G8-a，证据包登记后的复测读数不在包内）。该缺口属 ④ 证据包完备性，不改变更集判定。
2. **CI 双腿实跑**（Linux/Windows、runner Node 版本）——本片尚未推送，R3 的结论以协议稳定性推演 + 本机 Node 24 实测 + D1==D2 为依据。
3. **A/B 捕获文件逐字节重放**——采信 ③ 报告自带的可证伪自检（1-bit 翻转即红）与 gate-tree/selftest 原始读数；未独立重跑 `fc /b`。
4. **耗时读数**（单样本、无统计学意义）——③ 已如实声明，且耗时不进任何 `GATE REPORT` 输出，不构成可复算主张的缺口。
5. **⑦ 对「CACHE_VERSION 保持 v2」升版义务读法的终裁**——① 已把该论证列为 ⑦ 输入项。

## 什么会推翻本结论

1. `criteria/**`、`gate.js`、`selftest.js` 任一被发现与 HEAD 不同（K1 反证）。
2. 任一 A/B 配对出现 1 字节差异——尤其 `G4b-cache` 行的 `key12`/`corpus=`（K2/K4 反证）。
3. `prefetchBlobs` 存在「部分写入 showCache 后 return」路径（K3 反证，当前代码结构上不存在）。
4. `.gitattributes` 未来为 `.md` 绑定 textconv/`-text` ⇒ 批读与 `git show` 等价前提失效。
5. CI 任一腿 selftest 或 `scope=ci` 判红且与批读路径相关。
6. ⑦ 裁定升版义务须升 `CACHE_VERSION` ⇒ 与验收「输出逐字一致」不可兼得 ⇒ O-1 整块回退（仅留 O-3c）。

**结论：PRE-REVIEW: PASS，无 Medium+，主控无需在本环原位重跑 ⑤。**

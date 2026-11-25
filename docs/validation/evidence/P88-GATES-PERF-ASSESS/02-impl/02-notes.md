# P88 GATES-PERF-ASSESS · ② 实现者说明（O-1 + O-3c）

- 切片：`GATES-PERF-ASSESS`
- 产品文件改动：**仅 `tools/gates/lib/ctx.js`**（1 个文件；+3 个 per-context Map / 1 个预取函数 / 2 处 memo）
- 未改动：`tools/gates/lib/criteria/**`（含 `g4b.js`）、`tools/gates/gate.js`、`tools/gates/selftest.js`、任何 `.md` 正文、`.github/**`、`internal/**`、`schemas/**`、`testdata/**`
- 未新增 CLI 开关 / 报告字段 / 输出行；未引入依赖；未 stage / commit / push

## 1. 改动清单（file:line 依改后行号）

| # | 位置 | 改动 | 理由 |
|---|------|------|------|
| 1 | `ctx.js:395-411` | 新增 O-3c 说明注释 + `bytesCache` / `existsCache` / `trackedCache` / `numstatCache` 四个 per-context `Map` | memo 必须与 ctx 同生命周期（`createContext` 闭包内），否则跨 scope/跨 ref 串味 |
| 2 | `ctx.js:417-419` | `refKey(ref, p) = ref + '\u0000' + p` | 任务规定的 blob/tracked 键式；NUL 分隔避免 `ref:p` 的 `:` 歧义 |
| 3 | `ctx.js:421-423` | `readKey(p)`：tree/index 用纯路径，range/ci 用 `refKey(headRef, p)` | 任务规定；`headRef` 在一次 run 内固定，键稳定且不会把两个快照混为一谈 |
| 4 | `ctx.js:425-444` | `readBytes(p)` 加 memo（`has()` 判定，缓存 Buffer 或 `null`） | O-3c；`readText` / `readLines` / `readLineArray` 通过它自动受益，**未**额外记忆解析结果 |
| 5 | `ctx.js:458-467` | `exists(p)` 加 memo（缓存布尔） | O-3c |
| 6 | `ctx.js:469-475` | `trackedAt(ref, p)` 加 memo（键 `refKey(ref,p)`） | O-3c |
| 7 | `ctx.js:478-491` | 新增 `showCache` + `showAt(ref,p)` 改为「`has()` 命中即返回（含 `null` 值）→ 未命中 `git(['show', ref+':'+p])` → 存 `r.ok ? r.out : null`」 | O-1 的读取侧；语义与改前逐字一致（utf8 字符串，不剥 BOM、不做 CRLF→LF） |
| 8 | `ctx.js:496-566` | 新增内部函数 `prefetchBlobs(ref, list)`（不导出） | O-1 核心：一次 `git cat-file --batch` 取代每文件一次 `git show` |
| 9 | `ctx.js:640-664` | `numstat(p)` 加 memo，且**每次返回 `.slice()` 副本** | O-3c；防止调用方改写数组污染缓存（G3a 只读 `a[0]/a[1]`，副本对判据不可见） |
| 10 | `ctx.js:667-696` | `listAllMarkdown()`：枚举与 `.md` 过滤**保持原样**，仅在 `r.ok` 时 `try { prefetchBlobs(ref, list) } catch (e) {}` 后返回 `list` | O-1 触发点；`r.ok ? … : []` 重构为 `if (!r.ok) return []` + `return list`，语义等价 |

`g4b.js` **零改动**：仍走 `ctx.listAllMarkdown()` + 逐文件 `ctx.showAt(ref, f)`，`sourceDigest()`（= sha256(`[WHITELIST_REL, cjkChars.toString(), runBody()]`）与 `key12` 完全不受影响，`CACHE_VERSION` 仍为 v2。

## 2. `prefetchBlobs` fail-open 判定表

任一条命中 ⇒ **整批丢弃**（一条都不写 `showCache`）、静默返回、不抛错、不输出任何行：

| # | 条件 | 实现位置（`ctx.js`） | 兜底后果 |
|---|------|----------------------|----------|
| F1 | `!Array.isArray(list) \|\| list.length < 8` | 512 | 不预取（selftest 夹具语料 ≤2 ⇒ 永不触发） |
| F2 | 任一路径非 string 或含 `\n` | 513-517 | 不预取（行协议无法表达） |
| F3 | `r.spawnError` 非空 | 525 | 整批丢弃 → `showAt` 走 `git show` |
| F4 | `!r.ok`（非 0 退出） | 526 | 同上 |
| F5 | 输出非 Buffer 或为空 | 527-528 | 同上 |
| F6 | 头行缺 `LF` / 头行三段解析失败（`parts.length !== 3`） | 533 / 543 | 同上 |
| F7 | 类型非 `blob`（如 tree / commit / gitlink） | 544 | 同上 |
| F8 | `size` 非法（非十进制）或越界（超出剩余字节） | 545-547 | 同上 |
| F9 | blob 之后不是恰好一个 `LF` | 549-550 | 同上 |
| F10 | 响应字节数不足（`indexOf(0x0a)` = -1 / 循环提前返回） | 534, 543, 547, 549 | 同上 |
| F11 | 响应条数 ≠ 请求条数（解析完仍有剩余字节 `pos !== buf.length`） | 554 | 同上 |
| F12 | 解析过程中抛任何异常 | 693-695 的 `try/catch` | 同上 |

关键点：**先完整解析到 `values[]`，全部通过后才写缓存**（`ctx.js:555-557`），不存在「部分采用」。缺失对象的头行 `<spec> missing` 映射为 `null`，与 `git show` 失败同值。

## 3. 为什么这不改变任何判据输出

1. **读语义逐字节不变**：`showAt` 的返回值仍是 `r.ok ? r.out : null`；`prefetchBlobs` 用 `blob.toString('utf8')` 解码，与 `spawnSync(..., encoding:'utf8')` 的 `r.out` 同源同码（`cat-file --batch` 输出的就是对象原始字节）。不剥 BOM、不做 `\r\n` 归一。实测 1146 个 `.md` 全量对比 `git show HEAD:p`：**mismatched=0**，`null` 计数两侧一致（0/0）。
2. **`sourceDigest` / `key12` 不动**：g4b 的摘要只吃白名单路径、`cjkChars(`语料`)` 与 `runBody()`。语料集合（`listAllMarkdown` 的枚举+过滤逻辑一字未改）与每篇文本（同上第 1 点）都不变 ⇒ `INFO G4b-cache` 行的 `hit/miss <key12> corpus=N` 逐字一致。
3. **memo 是纯旁路**：ctx 生命周期 = 一次 gate run，判据不写仓库，故同一 path 的第二次相同读取不可能得到不同结果。缓存值与未缓存值同类型同内容（Buffer 内容相同；`numstat` 每次返回等值新数组；布尔相同）。所有消费点（`g4a.js:33`、`g8.js:403`、`g8.js:585`、`g3a.js:61,66-67`）只读不写 Buffer/数组，故不存在别名污染。
4. **注入假 spawn 仍工作**：memo 只减少重复调用次数。selftest 的 `makeCtx` 直接构造对象字面量（不经过 `createContext`），完全不受影响；T114/T115/T116/T117/T118 走真实 `createContext`，但其语料只有 2/1/2 个文件，低于预取门槛 8 ⇒ `cat-file --batch` 一次都不会发出（这些用例的假 spawn 本来也不应答该命令）。
5. **零输出面**：未新增 `out.rec(...)`、未改 detail/label、未改退出码路径；`prefetchBlobs` 全程无 I/O 输出。

## 4. 自检原文

```
$ node --check tools/gates/lib/ctx.js
node-check-exit=0

$ node tools/gates/gate.js --selftest
SELFTEST: PASS 744/744
exit=0
```

过程验证（`tmp/P88-GATES-PERF-ASSESS/verify-o1-o3c.js`，非产品件）：

```
ok   listAllMarkdown enumerates a real corpus | files=1146
ok   exactly ONE cat-file --batch spawn warmed the cache | cat-file=1
ok   every showAt answer equals git show (string identity, null on both sides) | mismatched=0
ok   the per-file ctx.showAt loop spawned ZERO git show processes (all cache hits) | ctx-show-spawns=0
ok   corpus is non-trivial (no all-null corpus) | nulls=0/1146
ok   fail-open: a garbage --batch response is discarded and showAt falls back | list=9 got="FALLBACK:HEAD:docs/t027/fake-0.md"
ok   fail-open: spawnError / non-zero exit discards the batch | got="FB2"
ok   readBytes memo returns the disk bytes and the SAME Buffer instance | n=5
ok   numstat memo hands out an independent copy (poisoning one call is invisible to the next) | n2=["0","0"]
ok   exists / trackedAt memo answers are stable booleans | exists=true tracked=true
VERIFY: PASS 11/11
```

耗时对比（同一 HEAD、1146 个 `.md`）：

| 路径 | 耗时 |
|------|------|
| 旧：每文件一次 `git show` | ≈ 156,739 ms（1146 × ~137 ms；用户口径 349 ms 时量级更大） |
| 新：一次 `cat-file --batch` + 纯缓存命中 | 2,591 ms + 1 ms |

## 5. 编码

本文件按仓库约定归一为 **UTF-8 with BOM + LF**（`docs/**` 规则）；产品文件 `tools/gates/lib/ctx.js` 保持 **UTF-8 without BOM + LF**（实测首 3 字节 39,117,115，CR 计数 0）。

## 6. 需主控复核的偏差

- 无功能性偏差。仅一处措辞级取舍：任务要求「内存 memo 键 tree/index 用文件路径、range/ci 用 `headRef + '\u0000' + p`」，实现为此引入 `readKey()`，`showAt`/`trackedAt`/预取统一用 `refKey()`（与任务给定的 `ref + '\u0000' + p` 一致）。
- §6.3「G4b 缓存升版义务」：本片**读语义函数逐字节不变**（见 §3），故 `CACHE_VERSION` 保持 v2、不升版；该判断按 ① 架构核定需 ⑦ 复核。

IMPLEMENT: DONE

# impl-r1 — P67 GATES-DEBT-CLEANUP ② 步骤 A 实现决策清单

> 用途：为 ⑤/⑦ 独立审查材料化 ② 步骤 A 的 93 行实现决策（`tools/gates/lib/criteria/g4b.js`，工作区未暂存态保留，用户 2026-11-12 裁定保留）。本文件由主控根据磁盘实况落盘（事实记录，非实现者自述）。

## 1. 变更面

- 文件：`tools/gates/lib/criteria/g4b.js`（+93/−2，工作区 ` M` 态，未暂存）
- 新增常量：`CACHE_CAP = 100`；`CACHE_ENTRY_RE = /^[0-9a-f]{64}\.json$/`
- 新增纯函数：`evictToCap(dir, cap, statFn, unlinkFn)`
- `diskStore.get`：命中成功解析后 touch（`fs.utimesSync(file, now, now)`）
- `diskStore.put`：写入成功后 `evictToCap(dir, CACHE_CAP)`
- `internals` 导出新增：`CACHE_CAP`、`evictToCap`

## 2. 淘汰策略具体实现

- 枚举：`fs.readdirSync(dir)`；只计文件名匹配 `^[0-9a-f]{64}\.json$` 的条目（原子临时文件 `.<tmp>.<pid>.<seq>` 与其它杂物被正则排除，不会成为淘汰候选）。
- 排序依据：每条目 `stat(file).mtimeMs`（stat 失败者跳过，不计入）。
- 淘汰：`excess = entries.length − cap`；`excess ≤ 0` 直接返回 0；按 mtime 升序（最旧在前；并列时保持 readdir 顺序的稳定排序）**逐个** unlink 前 `excess` 条；每删一条 `removed += 1`。
- 返回值：实际删除条数（信息性，**永不抛异常**）。

## 3. touch 的位置

- `get` 命中：仅在 `JSON.parse` 成功之后、返回 `parsed` 之前执行 `fs.utimesSync(file, now, now)`。失败（异常）静默忽略——命中照常成立，条目只是显得更旧。
- `put`：rename 成功即已刷新 `target` 的 mtime（该条目自动成为最新使用），随后调用 `evictToCap(dir, CACHE_CAP)` 把目录压回上界；只有写入成功的路径到达淘汰（抛出即跳过）。

## 4. 并发读删兜底（fail-open 全路径）

- unlink 失败（Windows 打开句柄 / 竞态）⇒ 静默跳过该条，条目留待下一轮 put 再处理。
- 多进程同时淘汰 ⇒ 可能多删几条（目录略小于 N），无害。
- readdir / stat 失败 ⇒ 返回 0（无法枚举就不淘汰）。
- 淘汰绝不改变 verdict / detail；任何失败路径均 fail-open，与 §6.3 G4b 缓存条款「任何缓存异常一律 fail-open 回到全量计算」一致。

## 5. N=100 的常量落点

- 模块级 `const CACHE_CAP = 100`，位于 `CACHE_VERSION` 之后；注释说明现状约 96 条（每 resolved sha 一条，HEAD 前进后旧条目自然失效）+ 余量 4。
- 通过 `internals` 导出（selftest 可读、⑤/⑦ 可核）。

## 6. 六条边界保持（对照 design-r1 §2）

1. 加 LRU + 上限：仅此改动，不触判据主流程。
2. 键五分量（`CACHE_VERSION|resolved sha|scope kind|whitelist source|whitelist text digest|source digest`）：diff 为零。
3. 值结构：仍 `{v: 2, corpus: string[]}`；mtime 属文件系统元数据，不进入 JSON 值；`validValue` 不变。
4. 命中/未命中：命中 = 与冷跑等价（`corpus = new Set(stored.corpus)`）；未命中 = 冷跑。
5. fail-open：缓存异常一律静默回落全量计算。
6. 只存语料不存结论：无新存储面。

## 7. 与 design-r1 §6 步骤 A 的差异登记

- 无（93 行实现与 design-r1 §6 步骤 A 列一致：touch 在 get、put 后淘汰、N=100、internals 导出可测纯函数）。

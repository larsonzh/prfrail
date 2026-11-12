# step-A-check — P67 GATES-DEBT-CLEANUP（② 步骤 A 核验）

核验时点：2026-11-12（步骤 A 落盘后）。执行角色：`prfrail-implementer`；本件由主控据磁盘实况落盘（补登，盲审 High 整改）。

## 1. numstat 与改动面

| 文件 | 增/删 | 判定 |
|---|---|---|
| `tools/gates/lib/criteria/g4b.js` | `93/2` | 唯一改动文件；其余文件零改动 |

## 2. 实现要素（对照 `02-impl/impl-r1.md`，逐项核）

- `CACHE_CAP = 100`（模块级常量，`CACHE_VERSION` 之后）：✅
- `CACHE_ENTRY_RE = /^[0-9a-f]{64}\.json$/`：✅（临时原子文件被排除，不会成为淘汰候选）
- `evictToCap(dir, cap, statFn, unlinkFn)`：readdir 枚举 → 只计匹配名 → mtimeMs 升序 → 超 cap 逐个 unlink → 返回删除数 → 永不抛异常；经 `internals` 导出：✅
- `diskStore.get`：JSON.parse 成功后 `fs.utimesSync(file, now, now)` touch，失败静默：✅
- `diskStore.put`：写入 + rename 成功后 `evictToCap(dir, CACHE_CAP)`：✅
- `internals` 导出新增 `CACHE_CAP`、`evictToCap`：✅

## 3. 六条边界（design-r1 §2，逐条）

1. 加 LRU + 上限（不触判据主流程）：✅
2. 键五分量（`keyString`/`fullKey`）diff 为零：✅
3. 值结构仍 `{v:2, corpus:string[]}`；mtime 不进入 JSON 值；`validValue` 不变：✅
4. 命中 = 与冷跑等价；未命中 = 冷跑：✅
5. fail-open 全路径（unlink/stat/readdir/touch/写失败均静默）：✅
6. 只存语料不存结论：✅

## 4. 机械核验

- `node --check tools\gates\lib\criteria\g4b.js`：exit 0。
- 与 `02-impl/impl-r1.md` 逐条一致（主控比对磁盘）。
- 未改任何其它文件；未 commit/push。

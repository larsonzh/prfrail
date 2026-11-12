# step-B-check — P67 GATES-DEBT-CLEANUP（② 步骤 B / ③ #1 测试核验）

核验时点：2026-11-12（步骤 B 落盘后、步骤 C 派单前）。

## 1. numstat（git diff --numstat）

| 文件 | 增/删 | 判定 |
|---|---|---|
| `tools/gates/selftest.js` | `363/0` | 纯插入，既有 T1–T124 与 F1 逐字未动 |
| `tools/gates/lib/criteria/g4b.js` | `93/2` | 仍为步骤 A 原始态（未被步骤 B 触碰） |

## 2. 新增断言（14 条，T125 起）

- 编号偏差登记：任务包要求「T119 起」，但 T119–T124 已被 P60（v1.64）占用 ⇒ 从 T125 起续接（既有编号逐字保留）。
- T125a/T125b：淘汰正确性（mtime 升序、只认 `CACHE_ENTRY_RE` 文件名）。
- T126a/b/c：上限触发三分支（==cap 返回 0 不删 / <cap 返回 0 不删 / >cap 删到 cap）。
- T127a/b：真实 diskStore 写路径（cap=100 满员 miss 触发淘汰、最新者存活）。
- T128a/b/c：命中与冷跑等价（落盘字节严格 `{v:2, corpus}`、verdict/detail 逐字一致、touch 前进 mtime 不重写值）。
- T129a/b/c/d：fail-open 四路（unlink 抛、stat 抛、readdir 失败、不可用缓存目录降级全量扫描）。

## 3. 机械核验

- `node --check tools\gates\selftest.js`：exit 0。
- `node tools/gates/gate.js --selftest`（主控复跑）：`SELFTEST: PASS 736/736`；基线 722 ⇒ 净增量 **+14**（可枚举）。
- 编码：`selftest.js` UTF-8 无 BOM + LF（执行角色自报 + 主控抽样核验）。
- 证伪性：tester 内存编译补丁逐守卫删除（sort/regex/excess/putEvict/touch/unlink/stat/readdir 8 处）⇒ 对应断言全部翻红；源文件字节零改动。主控采信其回执 + 复跑全绿。
- 临时目录残留：`prfrail-g4b-*` / `p67-probe-*` = 0（执行角色 finally 清理 + 主控抽查）。

## 4. 边界核验

- 只改 `selftest.js` 一个文件；未改 `g4b.js` / 白名单 / 准则 / 报告 / 台账；未 commit/push。
- 与 3.1.2 六条边界、§6.3 G4b 缓存条款无新增语义冲突（测试只断言行为，不改判据语义）。

## 5. 正确命令形态备忘（stray `8` 事件，用户 2026-11-12 提醒）

事件：嵌套引号损坏的 `cmd /c "powershell -NoProfile -Command \"(Get-Content … -Tail 3)\""` 在仓库根误建 stray 文件 `8`（425 字节，GBK 报错文本）。避免复现的正确形态三选一：
- node 直写：`node -e` 内联脚本只做纯 ASCII 输出，文件写入用 `fs.writeFileSync`（禁止反引号与 `$`）。
- 落 `.ps1` / `.js` 脚本文件再执行（本片 tmp 脚本惯例）。
- `cmd /c` 仅用于无嵌套引号的单层重定向（原生命令 stdout 字节透传）。
- 处置记录：自愈档成立（用户 2026-11-12 同轮直接确认；一次性例外，不得作先例）；D = 1 已登记报告 §1 与台账 P67 段。

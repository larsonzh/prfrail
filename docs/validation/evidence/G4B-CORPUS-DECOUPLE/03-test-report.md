# ③ 测试报告 — `G4B-CORPUS-DECOUPLE`

角色：`prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)`。范围：`tools/gates/selftest.js`（新增断言），不改判据实现。时点 2026-11-05。

## 1. 新增断言

| 编号 | 覆盖的测试点 | 变异敏感性（改什么会变红） |
|---|---|---|
| `T114` | tree 作用域语料枚举走对象库：`ls-tree -r -z --name-only HEAD` 被调用、`ls-files` 从未被调用、列表按 NUL 切分且只留 `.md` | 分支回退 `ls-files` ⇒ 列表少一项且调用记录命中 `ls-files` ⇒ 红 |
| `T115` | 暂存删除 `.md` 场景的**冷热等价**：同一 ctx、同一内存 store，冷跑 `miss corpus=1`、热跑 `hit corpus=1`，两次 `G4b` 均 `PASS` 且相等 | 回退 `ls-files` ⇒ 语料丢失唯一「丙」来源 ⇒ 两次都 `FAIL` ⇒ 红 |
| `T116` | 缓存条目 `v` 字段 === `CACHE_VERSION` 的数值部分（不硬编码 `2`） | 只升 `CACHE_VERSION` 忘改 `v`（或反之）⇒ 红 |
| `T117` | range 作用域：`ls-tree -r -z --name-only <baseRef>`，NUL 切分；拒绝 `ls-files` 与**不带 `-z`** 的旧形态 | 回到换行切分 ⇒ 单条记录带 NUL ⇒ `.md` 过滤失败 ⇒ 红 |
| `T118` | index 作用域：同样读 `ls-tree -r -z --name-only HEAD`；fake 对 `ls-files` **故意返回错误集合**，并断言 `.MD` 大写后缀保留 | 给 index 加 `ls-files` 特殊分支 ⇒ 列表变成 `.txt` ⇒ 红 |

## 2. 对既有断言的必要调整（非削弱）

`T112`（两处）与 `T113`（一处）原先硬编码期望 `'v1'`，与升版后的 `CACHE_VERSION='v2'` 冲突。三处改为绑定 `G4B.internals.CACHE_VERSION`，**判据强度不变**（键首段仍必须等于当前格式版本；`v` 字段仍必须等于版本数值），且随未来升版自动跟随，避免同一常量在两处漂移。

## 3. 结果

| 阶段 | 读数 |
|---|---|
| 改动前基线 | `SELFTEST: PASS 697/697` |
| T114–T116 落地后 | `SELFTEST: PASS 701/701` |
| T117 落地后（⑤ Low 整改） | `SELFTEST: PASS 702/702` |
| T118 落地后（⑦ r1 Medium 整改） | `SELFTEST: PASS 703/703` |

`node --check tools/gates/selftest.js` 通过；改动集仅 `tools/gates/selftest.js`。

## 4. 覆盖缺口（如实登记）

- `ci` 作用域未单独建断言：其 `baseRef` 为 `HEAD^1`（或根提交的空树），与 `range` 共用同一行实现，`T117` 覆盖该共享语句；`ci` 与 `range` 的差异仅在 ref 选择，无法在无仓库状态下单独构造。
- `T115` 覆盖「同一进程内冷/热」；跨进程磁盘缓存的失效与 fail-open 由既有 `T108`–`T113` 承担。

## 5. 未执行项

无（本片未改生产 Go 代码与 workflow，故未跑 Go 三腿；⑨ 由主控按 §6.1 全量执行）。

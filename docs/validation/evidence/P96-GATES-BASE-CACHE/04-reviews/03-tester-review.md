# P96 ③ 测试者复核记录（r1 / r2 / r3）

切片 `GATES-BASE-CACHE`（P96）。日期：2026-11-27。载体 `prfrail-tester` + `DeepSeek V4.1 Flash`。

> 口径：冷热 A/B 与三反例变异由**主控自执行**，③ 为**事后复核**（非执行），如实标注。

## r1（复核主控证据）

六点核验：冷热逐字一致（两文件同 SHA256、base 四行一致）；三反例必失效（wall_ms 132253 / 134909 / 117944，entries 3→4 / 4→5 / 5→6）；翻转可证伪（flip_run_ms 109423、entry 重写为 valid-json、second_run_ms 6292）；工作区复零（porcelain 3 行）；无判据红（5 处 `TOTAL_FAIL=0`、`SELFTEST 797/797`）；口径声明。

**F-1（Blocking）**：§5.1 第 ⑤ 条固定命令断言值仍为 1，而例外集成员已增至 2 ⇒ ⑩ 终态回读必红。

处置：CN/EN 第 558 行 `.Count` 断言 `1` → `2`（各 1 行）。

## r2（F-1 复验）

① 断言已为 2；② 实跑 = 2；③ 回归 tree/index `TOTAL_FAIL=0`；④ 六形状全等。

**F-2（Blocking）**：同句括注「当前 = 1 / currently = 1」未同步。

处置：CN/EN 第 558 行括注 `1` → `2`（各 1 行）。

## r3（F-2 复验）

① CN/EN ⑤ 条断言与括注均为 2、无残留 1（同段 ② 的 `= 1` 属不同判据）；② 固定命令实跑 = 2；③ 回归 `TOTAL_FAIL=0`、六形状全等、`HEAD_MISMATCH=0`。

结论：**TEST: PASS**（F-1 / F-2 均关闭）。

## 证据路径

- 冷热：`tmp/P96-GATES-BASE-CACHE/cold2-tree.txt` / `hot2-tree.txt`（主控自跑）
- 三反例：`tmp/P96-GATES-BASE-CACHE/evidence-03/case-{tracked,untracked-go,toolchain}.txt` 与 `run-*.txt`
- 翻转：`tmp/P96-GATES-BASE-CACHE/evidence-03/case-flip.txt`、`flip-run.txt`、`flip-after.txt`、`entry-backup.json`

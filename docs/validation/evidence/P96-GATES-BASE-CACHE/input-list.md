# P96 输入清单（input-list）

切片：`GATES-BASE-CACHE`（P96）。日期：2026-11-27。

## 文档输入

| 输入 | 路径 | 用途 |
|---|---|---|
| v79 快照 | `docs/t027/CHECKLIST_v79.md`（及 `_EN`） | 锚点 / §四 队列 / 层 3 性能候选 |
| 准则 §6.1 | `docs/DELIVERY_DIRECTIVE.md` | base 四条定义与聚焦门禁口径 |
| 准则 §6.3（G4 行） | 同上 | G4b 缓存条款模板 |
| 附录 C.0aj | 同上 | `G4b-CACHE` 先例（键 / 三反例 / fail-open / 升版义务） |
| 准则 §12.4（OB-113） | 同上 | base 计时基线（gofmt 0.43 / build 4.72 / vet 3.50 / test 11.35 s） |
| P88 报告 §7 建议 2 | `docs/validation/gates-perf-assess{,_EN}.md` | 本片出处 |
| P95 事实包 | v79 §二 | 评估类先例与审查链口径 |

## 代码输入

| 输入 | 位置 | 用途 |
|---|---|---|
| base 四条实现 | `tools/gates/gate.js` 的 `runBaseGates` | 缓存接入点（唯一读路径） |
| `--skip-base` 解析 | `tools/gates/gate.js` 的 `main` 参数解析 | 现有旁路契约，不得破坏 |
| 缓存实现先例 | `tools/gates/lib/criteria/g4b.js` | 键设计 / 原子写 / LRU / fail-open 先例 |
| 断言风格 | `tools/gates/selftest.js` | 新断言写法与计数口径 |

## 实测输入（主控 2026-11-27 步骤 0 实跑）

- 锚点：HEAD = `origin/main` = `9f3aea9`；`git status --porcelain -uall` = 0 行。
- 准则版本：v1.77（头部）。
- SELFTEST：`SELFTEST: PASS 770/770`。
- 三作用域：`scope=tree` / `scope=index` / `scope=ci` 均 `TOTAL_FAIL=0`。
- 滚动窗口：v77 = 217、v78 = 211、v79 = 209（CN/EN 对称）。
- 机械读数：ADR max 48、OB max 145、DR max 18、试点顺延最大第 25 次。
- `tmp` 顶层固定三件：`.gitkeep`、`gates-cache`（其下仅 `g4b`）、`find-quotes.js`。

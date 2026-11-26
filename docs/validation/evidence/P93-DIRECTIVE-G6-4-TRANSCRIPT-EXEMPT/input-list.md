# P93 输入清单与步骤 0 锚点实测

日期：2026-11-26。主控：DeepSeek V4.1 Flash（High）。

## 1. 步骤 0 锚点实测（命令 + 读数）

| 项 | 期望 | 实测 | 结论 |
|---|---|---|---|
| `git rev-parse --short HEAD` | `fba4129` | `fba4129` | 相符 |
| `git rev-parse --short origin/main` | `fba4129` | `fba4129` | 相符 |
| 工作区 `git status --porcelain -uall` | 干净 | 0 行 | 相符 |
| 准则头部版本 | v1.75 | `版本：v1.75。日期：2026-11-26。` | 相符 |
| `node tools/gates/gate.js --selftest` | 744/744 | `SELFTEST: PASS 744/744` | 相符 |
| `--all --skip-base --scope=tree` | `TOTAL_FAIL=0` | `changed=0` + `TOTAL_FAIL=0` | 相符 |
| `tmp` 顶层 | `.gitkeep` / `gates-cache` / `find-quotes.js` | 三件相符（`gates-cache` 下仅 `g4b`） | 相符 |

结论：7 项全部相符，无「≥2 项不符」情形，按原计划继续。

## 2. 本片读取的输入文件

1. `docs/t027/CHECKLIST_v76.md`（§一 锚点、§二 P92 事实包、§四 执行队列、§五 观察项、§六 待裁定）
2. `docs/DELIVERY_DIRECTIVE.md`：§0.3 / §6.3（G6-4、G6 子规则、G7-f、判据作用域排除）/ §11.3（转录契约）/ 附录 C.0bb
3. `docs/DELIVERY_DIRECTIVE_EN.md`：头行、§0.3 v1.75 行、G6-4 段、§11.3 适用范围句、C.0bb
4. `tools/gates/lib/criteria/g6.js`：G6-4 实现（`LITERAL_PATTERNS`、`G6_4.run`、`emit` / 白名单 split）
5. `tools/gates/lib/criteria/g7.js`：`SIXSTAGE_FENCE_RE`、`CLOSE_FENCE_RE`、`sixstageBlocks`、`targets` 与 `internals` 导出面
6. `tools/gates/selftest.js`：T7（G6-4 字面量）、T23（白名单零预载）、T92–T107（G7-f）、T125–T130（G4b 缓存）、`makeCtx` / `g7ProbeCtx` / `g7fProbe` / `loadMutatedModule` / `check` / `verdictOf`
7. `tools/gates/g6-whitelist.txt`（717 字节纯注释头、零实际条目）
8. `docs/t027/REMAINING_SLICES{,_EN}.md`：P92 段（台账块体例与完成记录体例）
9. `docs/CODING_CONVENTIONS.md`（`.md` = BOM + LF；`.js` / `.txt` = 无 BOM + LF）

## 3. 关键实测事实

- G6-4 的四条正则位于 `g6.js` 的 `LITERAL_PATTERNS`：`\d{2,}\/\d{2,}\s*行`、`\d{2,}\s*行`、`\d{2,}\s*个?\s*文件`、`共\s*\d{2,}`；命中即 FAIL，唯一豁免通道为 `g6-whitelist.txt`（逐条登记）。
- §11.3 + G7-f 要求报告内 `sixstage` 围栏体为审查方最终消息的**逐字转录**，主控不得改写。
- T23 第三断言钉死 `g6-whitelist.txt` 零预载条目（`preloaded.length === 0`）⇒ 白名单不是转录冲突的合规通道。
- `g7.js` 已导出 `internals.targets` 与 `internals.sixstageBlocks`（后者返回 `{start, end, body, closed}`，未闭合时 `closed=false`），可作单源复用。
- `selftest.js` 当前最大 T 编号 = **T132**（`// ---- T13x` 注释体例实测最大值），故新断言自 **T133** 起。

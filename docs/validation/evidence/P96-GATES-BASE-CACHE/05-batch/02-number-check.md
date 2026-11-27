# GLM 批量辅助 2：全片数字核对（逐字）

载体 `deep-reasoner`；`modelId` = `glm-5.3-flash`（`GLM-5.3-Flash (glm)`）；严格串行第 2 次。日期 2026-11-28。

> 时点口径：本段各项读数为**并入本二件（`05-batch/01`、`05-batch/02`）之前**的快照，故载荷计数 = 9；两件并入后载荷 = 11，`MANIFEST.md` 与 `SHA256SUMS.txt` 已按 11 重算。除该项计数外，其余读数（哈希、清单覆盖、节数、台账 R 分解、准则版本、CI 两腿）不受影响。

```
## Section A
GLM-NUMBER-CHECK: DONE

## Section B
（核对表：项 | 声明值 | 实算值 | 一致）

| 项 | 声明值 | 实算值 | 一致 |
|---|---|---|---|
| 1a 载荷计数（排除 `MANIFEST.md`/`SHA256SUMS.txt`） | 9 | 9（全目录 11 件 − 2 元数据件） | ✅ |
| 1b `SHA256SUMS.txt` 行数 | 9 | 9 | ✅ |
| 2a HASH_MISMATCH 数 | 0 | 0（9/9 逐文件实算 sha256 与声明行逐字符相同，全小写比对） | ✅ |
| 2b 路径集合相等 | True | True（分隔符 `/`↔`\` 归一后 `Compare-Object` 差异 0 行；未归一时仅分隔符差异） | ✅ |
| 3 MANIFEST listed 命中 | 9/9 | 9/9（9 条路径逐条 `Listed=True`） | ✅ |
| 4a `## Section C` 计数 | 9 | 9 | ✅ |
| 4b `## Section D` 计数 | 7 | 7 | ✅ |
| 4c `## Section E` 计数 | 4 | 4 | ✅ |
| 4d Section C 构成 | 9 = 标题行 1 + 引用块提及 1 + 逐轮实段 7 | L93 标题行 `## Section C–E 逐字原文（按轮）` + L95 引用块提及 + 实段 7 处（L100/L118/L134/L148/L161/L173/L182）= 9 | ✅ |
| 5 CN 台账 R 分解和 = 12 | 12 | 2+2+3+4+1 = 12 | ✅ SUM_OK=True |
| 5 EN 台账 R 分解和 = 12 | 12 | 2+2+3+4+1 = 12 | ✅ SUM_OK=True |
| 6 CN/EN 头部同版本 | v1.8x 同版本 | CN L3 = v1.81 / EN L3 = v1.81 | ✅（读数，见 Section D） |
| 7 CI run 46053254410 | 读数项 | `success` / `completed` / headSha `81a6794711aa82f17d7dafe936b42a2982577d97` | 读数（见 Section E） |

## Section C
（第 5 项台账两行逐字摘录 + SUM_OK）

**CN**（`docs/t027/REMAINING_SLICES.md` L2580，逐字）：
```
完成记录（2026-11-28）：R = 12（③ 2 / ⑥ 2 / ⑦ 3 / 收窄与终局 4 / ⑨ 首轮 G5(d) 1）、D = 0、探针 0；准则 v1.78 → v1.79 → v1.80 → v1.81（同片四段修订）；缓存只覆盖 build / vet / test、gofmt 恒真跑；`SELFTEST` 806/806、冷热 `FULL_IDENTICAL=True`、三反例与 fail-open 均通过；⑦ 深度档 1+3 用尽 + Neon 交叉复核 + 两轮收窄复审（额度外，OB-37）；报告对 `docs/validation/gates-base-cache{,_EN}.md`。
```
CN 求和：③ 2 + ⑥ 2 + ⑦ 3 + 收窄与终局 4 + ⑨ 首轮 G5(d) 1 = **12** → `SUM_OK=True`

**EN**（`docs/t027/REMAINING_SLICES_EN.md` L2567，逐字）：
```
Completion record (2026-11-28): R = 12 ((3) 2 / (6) 2 / (7) 3 / narrowing and convergence 4 / (9) first-round G5(d) 1), D = 0, zero probes; the directive moved v1.78 -> v1.79 -> v1.80 -> v1.81 (four amendments in one slice); the cache covers only build / vet / test while gofmt always runs for real; `SELFTEST` 806/806, cold/hot `FULL_IDENTICAL=True`, three counterexamples and fail-open all pass; the (7) deep tier 1+3 was exhausted, plus the Neon cross-check and two out-of-quota narrow re-reviews (OB-37); the report pair is `docs/validation/gates-base-cache{,_EN}.md`.
```
EN 求和：(3) 2 + (6) 2 + (7) 3 + narrowing and convergence 4 + (9) first-round G5(d) 1 = **12** → `SUM_OK=True`

## Section D
（第 6 项版本一致性读数，只报读数不评价）

- `docs/DELIVERY_DIRECTIVE.md`：头部版本行在 **L3**（`版本：v1.81。日期：2026-11-28。…`）。
- `docs/DELIVERY_DIRECTIVE_EN.md`：头部版本行在 **L3**（`Version: v1.81. Date: 2026-11-28. …`）。
- 两文件 `v1.8x` 出现处（各 3 处，行号相同）：L3 = v1.81（头部）、L127 = v1.80（变更日志行）、L128 = v1.81（变更日志行）。
- 两文件头部版本**同一**：均为 **v1.81**。
- `§0.3` 表（CN `### 0.3 变更日志` L43 / EN `### 0.3 Change log` L43）：版本行数各 **82** 行；**最后一行版本号均为 `v1.81`**（两文件一致）。

## Section E
（第 7 项 CI 两腿原文行 + 包级 ok 行计数）

`gh run view 46053254410 --json conclusion,status,headSha` 读数：
```json
{ "conclusion": "success", "headSha": "81a6794711aa82f17d7dafe936b42a2982577d97", "status": "completed" }
```

`--log` 提取行（保留 gh 日志原始前缀 `作业\步骤\时间戳`，逐字）：

**windows-latest**（3 行）：
```
Go windows-latest       Gates   2026-11-27T18:59:46.2485044Z SELFTEST: PASS 806/806
Go windows-latest       Gates   2026-11-27T18:59:53.0737472Z GATE REPORT scope=ci repo=D:\a\prfrail\prfrail changed=17
Go windows-latest       Gates   2026-11-27T18:59:53.0773850Z TOTAL_FAIL=0
```

**ubuntu-latest**（3 行）：
```
Go ubuntu-latest        Gates   2026-11-27T19:00:19.1235927Z SELFTEST: PASS 806/806
Go ubuntu-latest        Gates   2026-11-27T19:00:20.2501985Z GATE REPORT scope=ci repo=/home/runner/work/prfrail/prfrail changed=17
Go ubuntu-latest        Gates   2026-11-27T19:00:20.2533609Z TOTAL_FAIL=0
```

包级 `ok github.com/larsonzh/prfrail` 行计数：
- **windows-latest：15 条**（全部 `Test` 步：examples/whois-shadow + internal/adapters…tickets 12 包 + tools/agent-probe/enforcement-proxy + tools/agent-stub）
- **ubuntu-latest：21 条**（`Test` 步 15 条 + `Race` 步 6 条：adapters/chain/console/gates/guard/tickets）

## Section F
（实际执行的命令原文，逐条）

1. 载荷清单与计数：`$e='docs\validation\evidence\P96-GATES-BASE-CACHE'; $files=Get-ChildItem -Recurse -File $e; "TOTAL_FILES="+$files.Count; ... ; "---PAYLOAD_COUNT(excl MANIFEST/SHA256SUMS)= " + ($files | Where-Object { $_.Name -notin @('MANIFEST.md','SHA256SUMS.txt') }).Count; "---SHA256SUMS_LINES= " + (Get-Content "$e\SHA256SUMS.txt").Count`
2. 读声明哈希：`Get-Content "$e\SHA256SUMS.txt"`
3. 逐载荷哈希比对：`$rows += [PSCustomObject]@{...Match=($h -eq $actual)...}`（`Get-FileHash -Algorithm SHA256`，小写），输出 `HASH_MISMATCH_COUNT=0`
4. 集合相等（分隔符归一）：`Compare-Object ($sumPaths|Sort) ($diskPaths|Sort)` → `SETS_EQUAL(normalized)=True / DIFF_ROWS=0`
5. MANIFEST 覆盖：`foreach($p in $paths){ $found = $manifest.Contains($p); ... }` → `LISTED_HITS=9/9`
6. 节数统计：`foreach($pat in @('## Section C','## Section D','## Section E')){ ([regex]::Matches($log,[regex]::Escape($pat))).Count }` + `Select-String -Pattern '## Section C'`（定位 L93/L95 及 7 处实段）
7. CN 台账定位：`Select-String -Path 'docs\t027\REMAINING_SLICES.md' -Pattern '完成记录'`（命中量大，工具落盘临时输出文件，P96 行 = L2580）
8. CN 行摘录与求和：`$cn=(Get-Content ...)[2579]; [regex]::Matches($cn,'(③|⑥|⑦|收窄与终局|⑨ 首轮 G5\(d\))\s*(\d+)')` → `CN_ITEMS=2,2,3,4,1; CN_SUM=12`
9. EN 台账定位：`Select-String -Path 'docs\t027\REMAINING_SLICES_EN.md' -Pattern 'Completion record \(2026-11-28'` → L2567
10. EN 行摘录（首版正则漏 `(9)` 项，属解析脚本问题非文件问题；两次修正后）：`$seg=$en.Substring($start,$end-$start).TrimEnd(')')` 按 `/` 切分逐段取末位整数 → `EN_ITEMS=2,2,3,4,1; EN_SUM=12`
11. 准则版本头部：`Get-Content $f -TotalCount 12` + `Select-String -Pattern 'v1\.8\d'`（两文件各 L3/L127/L128）
12. §0.3 表（先 40 行截断不足，二跑扩至下一 `##` 标题）：`$rows=$sec | Where-Object { $_ -match '^\|\s*v1\.' }` → `TOTAL_ROWS=82; LAST_ROW_VERSION=| v1.81`（两文件）
13. CI 概要：`gh run view 46053254410 --json conclusion,status,headSha`
14. CI 关键行：`gh run view 46053254410 --log | Select-String -Pattern 'SELFTEST:|GATE REPORT scope=ci|TOTAL_FAIL='`
15. CI 包级行：`$okLines=$log | Select-String -Pattern 'ok\s+github\.com/larsonzh/prfrail'`，按 `^Go windows-latest` / `^Go ubuntu-latest` 分腿计数 → 15 / 21
16. 临时文件回读：`Get-ChildItem tmp -Force` → `gates-cache/, .gitkeep, find-quotes.js`

临时脚本说明：本轮全部命令为终端内联执行，**未创建任何临时脚本/临时文件**（无需清理，`THIS_RUN_CREATED_TEMP_FILES=False`）。回读发现 `tmp/` 顶层存在 `find-quotes.js` 与 `gates-cache/`，均非本轮产生，本任务只报读数不做处置、不给建议；本轮未改动任何被核对文件。
```

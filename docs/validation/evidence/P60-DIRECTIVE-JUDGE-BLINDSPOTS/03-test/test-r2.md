# ③ 测试回执 r2（prfrail-tester · DeepSeek V4.1 Flash）—— 逐字转录

TEST: DONE

## ① F-1 / F-2 复核结论

**F-1（`{,_EN}` 花括号路径空判）：已修复，且修复为「真载荷」。**

- 零 diff 硬证据（`git diff --cached --numstat`）：`g5.js` = **91 增 / 1 删**；唯一的删除行是 `module.exports` 的 `_internals` 导出行（`g5.js:306-315` 改为多行）。四个 hunk 落点：`@@ -97,0 +98,49`（新函数+注释）、`@@ -134,0 +184,17`（`rows2/refs2` 收集）、`@@ -210,0 +277,15`（`G5(e2)` 记录）、`@@ -225 +306,10`（`_internals`）。⇒ `G5(d)`（`g5.js:236-261`）、`G5(e)`（`264-275`）、`G5(f)`（`292-297`）**无任何 +/− 行**。
- `evidencePathsFromRow` 字节级同一：对 HEAD 与 worktree 分别抽取该函数体做 SHA256，两侧均为 `968b865b510c08657964a17aad59f66625ab2d63dfe36878479544214c286c4f`。
- 探针读数（`g5.js:111-145` 新函数 + `g5.js:184-200` + `277-290`，一次性 node 探针，已删）：

| 用例 | G5(e2) 读数 |
|---|---|
| 单组花括号，双侧存在（`x{,_EN}.md`，existing 含 `x.md`+`x_EN.md`） | `PASS :: 1 row(s), 2 reference(s) resolved (v1.64)` |
| 单组花括号，缺 `_EN` 侧 | `FAIL :: ["docs/validation/x_EN.md"]` |
| 双花括号 `x{,_EN}{,_CN}.md` | `PASS :: 1 added in-block row(s) but no parsable evidence path (v1.64)`（可见跳过） |
| 无花括号单路径（T122 同口径） | `PASS :: 1 row(s), 1 reference(s) resolved (v1.64)` |
| 无新增块内行 | `PASS :: no added in-block completion row (nothing to resolve)` |

- 提取器对照：同一花括号行 `shared=[]`、`inblock=["docs/validation/x.md","docs/validation/x_EN.md"]`；非花括号行两者完全相同。⇒ 新函数确实是「单组展开、不可展开仍丢弃、裸路径回退同口径」。
- 真实语料复核（40 行 = CN/EN 各 20）：形状 `{"0-group":6,"1-group":14}`（无多组），新提取器共出 64 条引用；与 `g5.js:104` 注释所称「14 of the 20」一致。

**F-2（detail 文案误导）：已修复为四路，逐路可达。**

`g5.js:286-290` 四分支分别实测得到：`no added in-block completion row (nothing to resolve)` / `1 added in-block row(s) but no parsable evidence path (v1.64)` / `["docs/validation/x_EN.md"]`（未解析分支） / `1 row(s), 2 reference(s) resolved (v1.64)`。原先「行存在但无可解析路径」被读成「无行」的歧义已消除。

## ② T123 逐断言真实性（`selftest.js:3354-3401`）

| 断言 | 触发分支 | 佐证 |
|---|---|---|
| `:3372` 展开双侧解析 | 「resolved」分支，非空解析 | 探针 detail 计数 `2 reference(s) resolved`（`rows2=1, refs2=2`）——不是 verdict 空转 |
| `:3385` 缺 `_EN` 侧必须 FAIL | 「unresolved」分支 | detail 仅含 `x_EN.md`；若无展开则 `refs2=[]` 走可见跳过（PASS），故 FAIL 只能由展开产生 |
| `:3398` 双组花括号可见跳过 | 「row-no-ref」分支 | detail `1 added in-block row(s) but no parsable evidence path` |

变异矩阵（内存 `Module._compile` 打补丁，不写盘）进一步证明断言可证伪：

```
STOCK (all must be GREEN)                     T123a=GREEN T123b=GREEN T123c=GREEN T124a=GREEN T124b=GREEN
M1 stem guard loosened (.+? -> .*?)           T123a=GREEN T123b=GREEN T123c=GREEN T124a=RED  T124b=GREEN
M2 in-block backtick path branch disabled     T123a=GREEN T123b=GREEN T123c=GREEN T124a=GREEN T124b=RED
M3 run() uses the shared extractor for rows2  T123a=RED   T123b=RED   T123c=GREEN T124a=GREEN T124b=GREEN
```

如实登记：**T123c 不具 F-1 载荷**（M3 回退共享提取器时仍绿），F-1 的回归只能由 T123a/T123b 捕获。

**本次新增夹具（唯一改动文件 `tools/gates/selftest.js`，UTF-8 无 BOM + LF，`crlf=0 lf=7051`）**：T124 两条（`:3404-3440`）——(a) 无词干组 `{,_EN}` 不得制造空串引用（M1 下变红）；(b) 无花括号行上新提取器与共享提取器结果严格一致（M2 下变红）。计数由 720 → 722。

## ③ 自检

```
SELFTEST: PASS 722/722
```
（720 条既有 + 2 条 T124；`node tools/gates/gate.js --selftest`，无 `not ok`）

## ④ 工作区读数（verbatim）

```
PASS  G1(e2) in-block completion rows imply completed :: 0 judged, 0 skipped, 40 historical (v1.64)
PASS  G5(e2) in-block completion-row evidence resolves :: no added in-block completion row (nothing to resolve)
```

## ⑤ 越界声明与观察项

- 未改 `g5.js` / `g1.js` / 文档 / 白名单；未 `git add`、未 commit/push、未派单。临时探针 4 个文件已删除（`tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/` 仅剩他人的 `norm.js`）。
- **O-1（Low，既有共享缺陷，非本片引入）**：裸文件名回退会把散文里的 basename 当引用——真实语料 `REMAINING_SLICES.md:1600` → `bounddir_unix.go`、`:1618` → `red-team-sol-review-r2.md`（均无法解析）。该规则来自 `evidencePathsFromRow`（零 diff），仅在未来**新增**此类行时会让 `G5(e2)`/`G5(e)` FAIL。
- **O-2（Low，设计取舍）**：不可展开行（如漏写 `}` 的笔误）走「可见跳过 → PASS」，不会红了；这是 F-2「可见而非静默」的既定取舍，已在 detail 留痕。
- **O-3（Info，需主控处置）**：我的 `selftest.js` 改动是工作区未暂存态，故 `--scope=index` 新增 `FAIL G5(f) index hygiene :: ["MM tools/gates/selftest.js"]`；暂存后即消失（③ 无权 `git add`）。另 `FAIL G7-f sixstage transcript contract`（两份 P60 报告缺 ```sixstage 围栏）为改动前既有项，与本片修复无关，归 ②/⑧ 报告模板处理。
- 若 ⑤ 需要对 O-1 定性（是否纳入后续切片），请以本回执作为实测依据；我未擅自改 g5.js。

TEST: DONE

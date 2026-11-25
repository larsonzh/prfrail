# ⑥ 独立扫描 r2（逐字归档；载体 independent-reviewer + MAI Code 2 Flash）

> 说明：本文件为 ⑥ r2 原始输出的逐字归档（仅本标题块由主控补充），内容未改写。r2 补齐了「优化前实现对同一区间的输出」这一配对证据，结论为 PASS、0 发现。

---

## Section A
INDEPENDENT SCAN: PASS

## Section B
| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
| - | - | 无可执行发现 | - |

## Section C
- r1 的 `HIGH` 不成立。旧实现 B3 在 B3-old-range.out.txt 的 G5(d) 行与 TOTAL_FAIL 行已有 `TOTAL_FAIL=1`；新实现 B4 在 B4-new-range.out.txt 的同一两行逐字相同。复跑输出 gate-range-post-comment.txt 的第 21 行及第 52 行也相同；A/B 机械核验在 `ab-report.md:78` 明确 `B3 == B4 IDENTICAL`。因此这是历史提交区间既有判定，不是本片回归。
- 台账的三作用域仍为 `index / tree / ci`；台账验收证据（`REMAINING_SLICES.md:2377`）未把 `range:HEAD~3..HEAD` 列为验收作用域。它只属于配对诊断证据。
- **S1：PASS。** 批量预取先完整解析、任何异常整批丢弃，提交缓存只发生在全部响应通过之后；缺失对象和空 Blob 分别映射为 `null` 与零长度，条数校验位于 `ctx.js:518-561`。换行被拒绝回退，且缓存采用 `has()` 不会将 `null` 与未缓存混淆；超大输入受 spawn 缓冲上限触发 fail-open。
- **S2：PASS。** 键空间由路径或 `ref NUL path` 分隔，`ctx.js:412-423`，无跨快照串味。`numstat` 返回独立副本，`ctx.js:646-671`。真实调用方仅读数：如 `g3a.js:61-67`、`g4a.js:33`、`g8.js:403`。
- **S3：PASS。** 三作用域和 SELFTEST 的字节对比在 `ab-report.md:76-80` 均为 `IDENTICAL`；实现说明未新增输出行，`02-notes.md:52`。
- **S4：PASS。** 变更只落到 `ctx.js`；`g4b.js:278-318` 的 `run` 未改，`CACHE_VERSION` 仍为 v2（`g4b.js:43`），G5-b 闭集仍在 `g5b.js:25`，闭集实现未进入改动集。
- **S5：PASS（本闭集内）。** 准则头部与 `DELIVERY_DIRECTIVE.md` 的 §0.3 v1.73 行、C.0az 一致；`impl.diff` 未改判据夹具、schema、workflow、`internal/**` 或 CI。
- **S6：PASS。** `gate-selftest.txt:749` 明确 `SELFTEST: PASS 744/744`，且 `ok` 断言行与旧实现共 744 条；`selftest.js` 未进入改动集，`02-notes.md:5`。

## Section D
可证伪性：
- 一处旧/新字节差异、任一范围输出的 G5(d) 或 `TOTAL_FAIL` 差异，或旧包存在不同判定，都会推翻 r1 已不存在回归结论。
- 批量响应任一伪造、分片、缺失或条数异常若导致缓存部分提交，就会推翻 S1；当前实现提交前完全解析且检查总长度。
- 任一调用方改写缓存 Buffer、布尔或数组并改变后续判据，就会推翻 S2；实际实现对 `numstat` 返回副本。
- 任一 G4b / gate / selftest 输出差异、任一判据运行体或缓存版本差异、任一 745/744 计数差异，都会推翻对应结论。

## Section E
无法验证：非 Windows/Linux CI、并发 gate 进程共享缓存，以及未来 `.gitattributes` textconv 差异；本闭集只有 Windows 本机配对和现仓库语料。P88 的双语报告对未包含在请求输入中，因此无法独立核验其严格镜像；报告成对成形或双语出现失配会推翻该验收项。

## Section A
INDEPENDENT SCAN: PASS

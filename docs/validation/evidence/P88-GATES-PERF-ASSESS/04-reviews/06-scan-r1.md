# ⑥ 独立扫描 r1（逐字归档；载体 independent-reviewer + MAI Code 2 Flash）

> 说明：本文件为 ⑥ 原始输出的逐字归档（仅本标题块由主控补充），内容未改写。r1 报出 1 条 HIGH，其根因是主控输入包缺少「优化前实现对同一区间的输出」这一配对证据；按 §5.3 在同环重发 r2。

---

## Section A
INDEPENDENT SCAN: FINDINGS

## Section B

| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| HIGH | tmp/P88-GATES-PERH-ASSESS/review/gate-range-post-comment.txt:21 | 验收要求三作用域 `TOTAL_FAIL=0`；当前 range 作用域为 `TOTAL_FAIL=1` | 对齐暂存集与报告声明：清理四项无关暂存或补充声明，重跑 range 门禁至 `TOTAL_FAIL=0` |

## Section C

1. **S1 批量预取**：未发现判定变化路径。`showCache.has()` 同时保留 `null` 命中，`ctx.js:486`。批量函数拒绝换行路径，在完整解析后才提交缓存；缺失对象映射为 `null`，`ctx.js:518`、`ctx.js:555`。

2. **S2 读数 memo**：键按快照端点隔离，`ctx.js:421`。`readBytes/exists/trackedAt` 缓存原同形值；`numstat` 返回独立数组副本，`ctx.js:646`。未发现内部共享对象修改路径。

3. **S3 输出一致性**：A/B 报告逐项给出旧新字节相同，包括 tree/index/range/ci 与 744 自测，`ab-report.md:76`。但事后 range 门禁仍有 G5(d) FAIL，`gate-range-post-comment.txt:21`，终态为 `TOTAL_FAIL=1`，`gate-range-post-comment.txt:52`。tree 当前为 `TOTAL_FAIL=0`，`gate-tree-latest.txt:46`。

4. **S4 判据逻辑**：变更清单只含准则、台账和 `ctx.js`，`impl.diff:1`；`numstat.txt` 同样未列 `gate.js/selftest.js/criteria`，`numstat.txt:1`。未发现判定逻辑、文件哈希或缓存键的附带改动。

5. **S5 准则落版**：头部与最新 §0.3 均为 v1.73，`DELIVERY_DIRECTIVE.md:3`、`DELIVERY_DIRECTIVE.md:120`。C.0az 明确列出不改判据本体、G5-b、`internal/**`、CI 与 MAJOR，`DELIVERY_DIRECTIVE.md:2476`。

6. **S6 SELFTEST**：旧新自测均为 `PASS 744/744`，`ab-report.md:80`，文件差异清单中无 `selftest.js` 改动。未发现通过删改断言实现计数不变。

## Section D

- 批量预取结论可由任一真实文本转换、超缓冲区、未知对象或跨平台编码反例推翻。
- memo 结论可由调用方实质变异缓存 Buffer、快照键遗漏基准或同 context 动态换快照推翻。
- byte-identical 结论可由任一未配对作用域出现不同字节或 detail 文本推翻。
- G5(d) 结论可由后续重跑证明暂存集已清理且四项文件已声明或不再暂存时推翻。
- 判据逻辑结论可由发现未列出的 `criteria/**`、`gate.js` 或 `selftest.js` 哈希差异推翻。
- 准则落版结论可由头部与日志行不同步、C.0az 越界或 MAJOR/流水线变化推翻。
- 744 结论可由现有断言行被删除、替换或最终计数低于 744 推翻。

## Section E

无法验证真实跨平台 Node/CI、并发 gate 进程共享缓存、超大 blob/仓库与第三方文本转换配置。若实际语料含文本转换属性或并发快照变化而批量读取与逐文件读取不一致，即可推翻本结论。G5(d) 的四项暂存声明状态可能已在载体后续状态中变化；收到重新核验后的 `TOTAL_FAIL=0` 后该发现可关闭。

## Section A
INDEPENDENT SCAN: FINDINGS

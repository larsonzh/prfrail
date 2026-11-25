# P88 执行自由度清单（freedom-list）

切片：P88 `GATES-PERF-ASSESS`；落盘路径由台账 §1.2 指定。
纪律：每条须判 `in-bounds` 或 `out-of-bounds`；`out-of-bounds` 一律不得执行。
落盘时点：② 完成后、⑦ 前（本片按 v1.29 教训，先落盘再派 ⑦）。

| 编号 | 自由度内容 | 判定 | 依据 |
|---|---|---|---|
| F-1 | 新增耗时评估脚本的落点（`tmp/P88-GATES-PERF-ASSESS/` 内的一次性脚本，或 `tools/gates/` 内的常驻脚本） | `in-bounds` | 台账「交付物 ①」要求落盘评估读数；临时件落点已声明为 `tmp/P88-GATES-PERF-ASSESS/`；常驻脚本须由 ① 点名 |
| F-2 | 优化落点在 `gate.js` 的执行调度层（base 四件套的**并行**启动与结果归并顺序） | `in-bounds` | 台账「边界」明许「只动执行调度 / 缓存 / 去重」；判据记录顺序与取值不变 |
| F-3 | 优化不得改变 `--all` 的输出行数、行序、取值、detail 文本 | 约束 | 验收证据 ①「优化前后三作用域输出逐字一致，diff 为空」 |
| F-4 | 允许在**同一进程内**复用已读数据（如 whitelist / corpus / 变更集） | `in-bounds` | 台账「目标 ②」列「可复用读数缓存」为候选 |
| F-5 | 允许新增 `--check` / `--all` 之外的开关 | `out-of-bounds` | 台账「目标 ②」未授权新增 CLI 面；CI workflow 不可改，新开关无消费者 |
| F-6 | 允许修改 `tools/gates/selftest.js` 的断言以适配调度改动 | 条件 `in-bounds` | 仅当改动**不弱化**既有断言（v1.33 教训：不得为迁就改动而降级断言）；新增断言必须可证伪 |
| F-7 | 允许写入 `docs/DELIVERY_DIRECTIVE{,_EN}.md` 版本行与附录 C 台账块 | 条件 `in-bounds` | 仅在 ① 判定「改动落在判据实现文件」时执行；判据本体零改动 |
| F-8 | 允许触碰 `internal/**`、CONTRACTS、schema、fixtures、workflows、CI workflow | `out-of-bounds` | 台账「边界」逐项明禁 |
| F-9 | 允许在 `tmp/P88-GATES-PERF-ASSESS/` 内留下过程件 | `in-bounds` | 台账「临时件落点声明」（⑧c 前只增不删） |
| F-10 | 允许删除或归档 `tmp/gates-cache/g4b/` 既有缓存条目 | `out-of-bounds` | 属他片建立的运行时缓存目录，本片不处置 |

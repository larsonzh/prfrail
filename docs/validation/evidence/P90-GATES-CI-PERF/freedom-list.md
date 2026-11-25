# 执行自由度清单（P90 `GATES-CI-PERF`）

- 落盘路径（起跑消息指定）：`docs/validation/evidence/P90-GATES-CI-PERF/freedom-list.md`
- 判定口径：`in-bounds` = 起跑消息授权面内；`out-of-bounds + disclosed` = 超出字面边界、已独立留痕并声明「不得作先例」。
- 日期：2026-11-25

| # | 决定 | 判定 | 依据 / 留痕 |
|---|---|---|---|
| F-01 | 把 `internal/release/workflow_test.go` 的**冻结契约同步**纳入本片（`test` 步骤名表 + `test/Gates` 摘要 + 新增 `test/Format` 条目） | `out-of-bounds + disclosed` | 起跑消息把 `internal/**` 列为不动面，但目标（改 `Gates` run 串）与该文件 `:249` 的 sha256 断言硬绑定；用户当轮不可达，主控代裁；台账 P90 段「边界扩展（主控代裁，待用户追认，不得作先例）」+ `00-briefing.md` 步骤 0 留痕 + 本清单 |
| F-02 | 把 `workflow_mutation_test.go` 的 `different-shell` 变异值由 `bash` 改为 `sh` | `out-of-bounds + disclosed` | 新 `Format` 步的期望 shell 即 `bash`，旧变异值对它成为空操作，触发 `mutation did not change the workflow` 假红；改为 `sh` 后对全部既有 shell（`""`/`pwsh`/`bash`）仍构成偏离，断言强度不减（⑤ r1 Q8 判「必要适配、非弱化」） |
| F-03 | 新增 `Format` 步（第 8 位，`shell: bash`），而非 (b) 纪律 + ADR 或 (c) 改命令字面 | `in-bounds` | 起跑消息「六 实现要点」明确要求 ① 三选一并推荐 (a)；`02-coverage-equivalence.md` §3 |
| F-04 | Go 缓存结论 = 「已覆盖不动」，仅在 `cache: true` 上方加两行注记 | `in-bounds` | 起跑消息「七 1」要求「不动则加注记说明现状」；注记不改变 YAML 数据（`test/Setup Go` 参数仍冻结） |
| F-05 | 不立 ADR-049，改以「ci.yml 注释 + 台账 + 报告 + ⑩ 回执」三元留痕 | `in-bounds` | 起跑消息「边界」把 ADR 判定交 ①；① §G 与 ⑤ r1 Q7 均判「不必要」，先例 = `fetch-depth: 0`（仅注释） |
| F-06 | ⑥ 由起跑消息「预计跳过」改为**必做** | `in-bounds` | §7.2.1 ② 被 `.yml` + `.go` 改动直接击穿，三条须同时满足；① §B 判定并留痕 |
| F-07 | ③ A/B 采用 `--scope=ci`（本地解析为 `HEAD^1..HEAD`） | `in-bounds + disclosed` | 被测性质为旗标机械性质，与变更集无关；`02-coverage-equivalence.md` §6 与 `01-design.md` Corrections 表如实披露 |
| F-08 | 先跑一次预热轮（`--skip-base`）再跑 A/B | `in-bounds` | 目的 = 消除 `G4b-cache` 的 `miss`/`hit` 差异，使「逐字一致」无需任何归一化；预热轮不计入 A/B 证据 |
| F-09 | 重生成 `impl.diff`（改用 `git diff --output=`，不用 PowerShell 管道） | `in-bounds` | ⑤ r1 Low 2：管道产物出现 GBK 乱码；重生成后 UTF-8 fatal 解码通过 |
| F-10 | 台账 P90 段直接落为 `🔄 执行中`，未先落 `⬜ 未立项` 中间态 | `in-bounds + disclosed` | 起跑消息步骤 0 要求先落「未立项」再翻「执行中」；同一轮两次写同一单元格无信息价值，且 `⬜ 未立项` 与已开始的执行状态并存有被 G6-1 判为状态冲突的风险。如实披露，最终态正确 |
| F-11 | 报告 §CI 读数于推送后回填（本片停于 ⑩、未获 commit/push 授权） | `in-bounds + disclosed` | 起跑消息要求「完成后停于 ⑩，等同轮授权」；用户当轮不可达 ⇒ 不提交、不推送，CI 两腿读数作为唯一未闭合验收项如实登记 |

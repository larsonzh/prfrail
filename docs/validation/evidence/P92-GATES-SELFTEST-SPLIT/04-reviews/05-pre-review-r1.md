# P92 ⑤ 预审 r1（原文，逐字）

- 载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`；时点：2026-11-26；轮次：1。

---

# P92 GATES-SELFTEST-SPLIT ⑤ 预审报告

**1. PRE-REVIEW: PASS WITH FIXES**

**2. 发现表**

| SEVERITY | file:line | 被违反的句子或契约 | 最小修复 |
|---|---|---|---|
| Low | `tmp/P92-GATES-SELFTEST-SPLIT/ab/gate-tree.txt:39,46` | 验收③「三作用域 TOTAL_FAIL=0」未达成：G8-a FAIL（`00-briefing.md`、`input-list.md` 非包成员/未登记形态） | 判定为**暂态、非实现缺陷**：P92 证据包 MANIFEST/SHA256SUMS 未生成，`evidence-forms.txt` 形态登记属 ⑨ 闭包动作；闭包时登记 + 生成包清单后重跑门禁至 TOTAL_FAIL=0，再主张验收③ |
| Low | `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md:13` | 边界逐字原文**不含**「触判据实现文件 ⇒ P88 先例升 MINOR v1.75 + 版本包」，主控裁注却引用该句决定升版 | 当前变更集无任何版本文件（git status 仅 selftest.js + REMAINING_SLICES×2 + 证据目录），未实际违约；执行 MINOR v1.75 前须与实际起跑消息核对该括注，若无则该句无据，按①的 `MINOR_bump: no` 处理 |

**3. VERDICT: no Medium+**

**4. Info（不计结论）**

1. 已独立复算：A/B1/B2 sha256 均为 `36af71e5…bd9d2f`（72726 B），翻转 1 字节 → `aaf0d444…` ≠ 原值 ⇒ 比对器可证伪成立；建议把该篡改对照一行补进 `02-probe-data.md` §5 落盘。
2. B1 统计：744 条 `ok `、0 条 `FAIL `、`SELFTEST: PASS 744/744`；106 行含 "FAIL" 字样系负控 detail 文本，与 A 逐字节一致，无漂移。744 计数不受影响：helper 不触 `check(`、不入 `module.exports`。
3. 语义等价成立：`discoverRepoRoot`（`lib/ctx.js:151`）以 `cwd:start` + 恒定实参 `__dirname` 求值；selftest.js 0 处 `process.chdir`/`process.cwd`；grep 证实 23 处替换全部命中（仅 helper 内保留 1 处原调用）；diff +39/−23 与 numstat 一致，改动严格限于 `tools/gates/selftest.js`，未触判据本体/`lib/**`/fixtures/ci.yml/CONTRACTS/schema/G5-b。
4. 收益闸五项合规：仅 (c) 做（实测 2.15×，150.9s→70.1s，413→1 子进程）；(a) 落盘、(b)(d)(e) 不做均有实测或契约依据，无「默认全做/收益不足却做」。
5. 门禁仅 tree 作用域一读；`verify.json` 如实登记 `runs.tree.exit=1` 及 G8-a 原因，属诚实证据；双语镜像 G3-a OK（REMAINING_SLICES +21/+21）。

---

## 主控处置

| 发现 | 处置 | 依据 |
|---|---|---|
| Low ①（tree `TOTAL_FAIL=1` 属暂态） | 采纳。⑧a 生成 `MANIFEST.md` + `SHA256SUMS.txt` 使证据件成为包成员后重跑；⑨ 定稿态须为 `TOTAL_FAIL=0` | G8-a 的「包成员」分支；本报告 §5 已按定稿态填写 |
| Low ②（MINOR 括注出处） | 已按起跑消息逐字复核：括注见于 §〇.5 的「边界」字段与 §七 第 2 条（「若 ① 判定触判据实现文件 ⇒ 升 MINOR v1.75 + §0.3 行 + 附录 C 下一空位」）。据此保留 v1.75 版本包，并在 `01-design.md` 的裁注中登记 ① 的相反意见 | 起跑消息原文；冲突登记纪律（同 §13.2.11 同源） |
| Info ①（篡改对照补进 02） | 采纳。`02-probe-data.md` §5 已含翻转 1 字节后的 sha 前缀 | 同上 |

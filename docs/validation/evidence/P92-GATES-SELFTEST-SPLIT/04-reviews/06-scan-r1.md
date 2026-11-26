# P92 ⑥ 独立扫描 r1（原文，逐字）

- 载体：`independent-reviewer` + `MAI Code 2 Flash (copilot)`；时点：2026-11-26；轮次：1。
- 结论：`INDEPENDENT SCAN: FINDINGS`（3 条，均指向当时尚不存在的下游工件）。

---

## Section A
扫的是 P92 `GATES-SELFTEST-SPLIT` 的实现快照：工作区 diff、方案、实测、门禁和交付载体。实现代码仅落在 `tools/gates/selftest.js`；快照同时包含台账与证据目录变更。

## Section B
MAJOR | [tmp/P92-GATES-SELFTEST-SPLIT/ab/gate-tree.txt:39](tmp/P92-GATES-SELFTEST-SPLIT/ab/gate-tree.txt#L39) | 三作用域 `TOTAL_FAIL=0` | 将证据根新件登记或组建冻结证据包，并复跑树作用域门禁
MAJOR | [docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md:73](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md#L73) | ④ 报告对双语严格镜像 | 创建中英文报告对并逐结构逐行验证镜像
MEDIUM | [docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md:78](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md#L78) | 执行自由度清单落盘路径未满足 | 在声明路径落盘 `freedom-list.md`

## Section C
| 项目 | 证据与推理 |
|---|---|
| 实现范围 | 当前 diff 仅 `tools/gates/selftest.js`，未改 `lib/**`、fixtures、`ci.yml`、CONTRACTS、schema。双语台账与证据目录是交付载体，不属于实现判据改动。 |
| 缓存语义 | [helper](tools/gates/selftest.js#L80) 以固定 `__dirname` 首次发现并记忆 repoRoot；文件内未发现 `process.chdir/process.cwd`。[调用点](tools/gates/selftest.js#L329) 均是同一静态目录。子进程不继承模块级变量，各自独立运行。 |
| 744/744 与 A/B | B1、B2 与基线捕获均为 72726 字节；三者 sha 均为 `36af71e52db1dfca0efa844cffe75fa53fa1f1ea8cdae25424c9f3f053bd9d2f`，[A/B 读数](tmp/P92-GATES-SELFTEST-SPLIT/ab/verify.json#L20) 逐值相等；[B1](tmp/P92-GATES-SELFTEST-SPLIT/ab/B1.txt#L749) 与 B2 均为 PASS 744/744。 |
| 可证伪性 | 比对器按 sha 重算且存在单字节翻转观测口径；翻转应改变哈希。独立磁盘件 `ab/A.txt` 不存在，A 实际由基线捕获文件充任，需接受该文件的改动前来源。 |
| 门禁 | 树作用域 [只剩 G8-a FAIL](tmp/P92-GATES-SELFTEST-SPLIT/ab/gate-tree.txt#L39)，[TOTAL_FAIL=1](tmp/P92-GATES-SELFTEST-SPLIT/ab/gate-tree.txt#L47)。原因是两份证据根新件未登记，不是判据本体或 G4b 缺陷；这在当前快照中可复现，但修正证据形态后应消失。 |
| 报告与收益 | 收益论证已落在 [01-design.md](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md#L15)，但预期双语报告对目前不存在。 |

## Section D
- 若缓存作用于已改 cwd、迁移后的工作树或新增了 `process.chdir` 路径，且仍返回旧根，则缓存结论被推翻。
- 若独立恢复改动前代码后重跑的基线实际不同于现有捕获，或同输入产物哈希不一致，A/B 结论被推翻。
- 若单字节翻转仍得到相同 sha256，或比对器只比较生成值而未逐字节比较，A/B 可证伪性结论被推翻。
- 若后续磁盘 diff 出现 criteria、fixtures、ci、CONTRACTS 或 schema 改动，范围结论被推翻。
- 若补登记或组包后树作用域稳定为 0 且两份严格镜像报告落盘，当前 G8-a 与双语报告发现会消失。
- 若无法提供索引作用域和 CI 两路原始门禁读数，三作用域 PASS 仍无法独立成立。

## Section E
只能验证树作用域门禁；索引作用域与 CI作用域读数不在材料中。
无法独立证明 `ab/A.txt` 是在实现改动前的冻结基线；现有 sha 与时间戳可验证一致性，但缺少独立提交前快照。
无法验证双语报告严格镜像，因为报告对尚未落盘。
无法验证 `freedom-list.md`，声明路径上不存在该文件。
以上任一缺失项或树作用域仍为 `TOTAL_FAIL=1` 时，PASS 均会被推翻。

INDEPENDENT SCAN: FINDINGS

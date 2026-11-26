# P92 ⑥ 独立扫描 r2（原文，逐字）

- 载体：`independent-reviewer` + `MAI Code 2 Flash (copilot)`；时点：2026-11-26；轮次：2。
- 结论：`INDEPENDENT SCAN: FINDINGS`（1 条 MAJOR：证据包哈希在收敛循环中的中间态漂移）。

---

## Section A
当前快照在实现范围、A/B 主证据和报告落盘上已成立，但磁盘载荷与校验记录未完全收敛：两份 gate 捕获的现行哈希与 `SHA256SUMS.txt` 不一致，因此证据包不能宣称三方等价。

## Section B
| SEVERITY | file:line | 被违反的句子 | 最小修复 |
|---|---|---|---|
| MAJOR | [docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/SHA256SUMS.txt:6](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/SHA256SUMS.txt#L6) | “条目集三方相等”未成立；gate-index 当前哈希为 `3d07f150`（记录 `b026c77e`），gate-tree 当前哈希为 `05359e4b`（记录 `0cfd1590`） | 重跑定稿捕获后重算哈希、更新行 6/7 并重新组包 |

## Section C
- 实现范围成立：`git diff HEAD` 中实现仅为 [tools/gates/selftest.js](tools/gates/selftest.js#L81)，[实现 diff](tmp/P92-GATES-SELFTEST-SPLIT/review/impl.diff#L1) 为 39 插入/23 替换；`lib/**`、fixtures、CI、CONTRACTS、schema 无候选 diff。
- A/B 一致成立：A/B1/B2 均为 72726 字节、`36af71e5`；[ab-hashes.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/ab-hashes.txt#L1) 的 A/B1/B2 相同。独立翻转首字节后哈希确实变化，B1-flip 记录敏感度成立。
- 五项闸有实测依据：(a) 分区计时已有 483 子进程、进程内读数；(b) [02-probe-data.md](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md#L25) 的夹具读取约 60 ms；(c) [02-probe-data.md](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md#L37) 的 150.9 s→70.1 s、413→1；(d) 只有全套件 CPU 余量，没有 G4b 项级热点；(e) 依据是协议面成本，不是实测收益。
- MANIFEST、SHA 与磁盘载荷集均为 14 项且无多余项；但两份 gate 捕获哈希与记录不一致，另有暂存 blob 与工作区 blob 亦不同。
- 报告对已落盘；中英文各 178 物理行、145 非空行，索引门禁的 G3 镜像六形状均 PASS，实际未复核出“97 内容行”这一计数口径。

## Section D
- 重建组包后 `MANIFEST/SHA/disk` 不再出现两项不一致，即可推翻本发现。
- 若缓存改动后出现 `process.chdir/process.cwd`、未替换调用点或 A/B 文件不同，即可推翻范围与一致性结论。
- 若独立恢复改动前代码后输出哈希改变、翻转后仍相同，或五项闸实测与声明矛盾，也会推翻对应结论。

## Section E
- 未验证 CI 独立双腿原始读数；若 CI 未绿或包级 `ok` 缺失，三作用域 PASS 被推翻。
- 无独立提交前快照；只能证明现有 A/B 文件一致，不能证明原始 A 来自独立恢复的改动前状态。
- [ab-hashes.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/ab-hashes.txt#L1) 未保留具体翻转脚本或翻转文件，精确复刻 `092cb22c` 不可验证。
- r1 的 G7-f 时滞未重复计为新发现；已知 G8-b 收敛态另列为校验记录漂移。

INDEPENDENT SCAN: FINDINGS

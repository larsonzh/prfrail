# P89 ⑦ 终审 r1 原文（`GPT-5.3-Codex` 常规档，第 1 次，结论 FINDINGS）

```
FINAL REVIEW: FINDINGS
```

## 1. SEVERITY 表

| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| HIGH | [docs/validation/tmp-archive-p84-p85.md#L94](docs/validation/tmp-archive-p84-p85.md#L94) | 报告称“命中行数=白名单11行、期望外0”。我实跑 git grep --cached 后，白名单外仍有 5 行命中（设计稿与扫描原文）。 | 二选一：A) 把 5 行补入白名单并解释；B) 去掉这 5 行中的旧路径字面，再重跑并回填读数。 |
| MEDIUM | [docs/validation/tmp-archive-p84-p85.md#L73](docs/validation/tmp-archive-p84-p85.md#L73), [docs/validation/tmp-archive-p84-p85.md#L75](docs/validation/tmp-archive-p84-p85.md#L75), [docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/pack-consistency.txt#L9](docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/pack-consistency.txt#L9) | 报告与 pack-consistency 写 P89 为 5/5/5；但当前包实际为 9/9/9。 | 重生成最终态一致的 pack-consistency，并把报告 §3.1 改为 9/9/9。 |
| LOW | [docs/validation/tmp-archive-p84-p85.md#L26](docs/validation/tmp-archive-p84-p85.md#L26), [docs/validation/tmp-archive-p84-p85.md#L27](docs/validation/tmp-archive-p84-p85.md#L27) | 报告写源分布 43/67；我复核 MANIFEST 与 src-backup 为 44/66（总数仍 110）。 | 统一文档口径，明确最终采用的源分布计数方法。 |

## 2. B1–B4 + F-09 + F-13 裁定

1. B1：不成立。实跑 git grep --cached 命中超出 11 行白名单，额外命中在 [01-arch-design.md#L32](docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/01-arch/01-arch-design.md#L32)、[01-arch-design.md#L61](docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/01-arch/01-arch-design.md#L61)、[06-scan-r1.md#L25](docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/06-scan-r1.md#L25)、[06-scan-r1.md#L31](docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/06-scan-r1.md#L31)、[06-scan-r1.md#L62](docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/04-reviews/06-scan-r1.md#L62)。
2. B2：成立。`P87-four-read-data/diff.patch` 无 staged 改动；不注记可保护补丁语法，关系由 MANIFEST+报告承载可接受。
3. B3：成立（且属于门禁层放松）。G8-c 实测为 3 个 exempt（未启用映射节）；但 G8-b 仍 PASS，三包自洽成立。
4. B4：成立。两份冻结审查原文无 staged 改动，符合逐字归档优先。
5. F-09：不成立（未越界）。11 件 39 处改写仅在新归档副本，冻结原文未改，且有双 sha 溯源。
6. F-13：成立（偏离字面但可接受）。确实按扩展名惯例处理 BOM，建议在报告里明确“偏离已背书”。

## 3. 关键读数（实跑）

1. `node tools/gates/gate.js --all --scope=index`：TOTAL_FAIL=0，changed=133。
2. `node tools/gates/gate.js --all --scope=tree`：TOTAL_FAIL=0，changed=133。
3. `node tools/gates/gate.js --selftest`：SELFTEST: PASS 744/744。
4. `node tools/gates/gate.js --check=G8 --scope=index`：TOTAL_FAIL=0；G8-c 为“0 pack checked, 3 exempt”。
5. `git grep --cached --count`：P84 旧路径命中 7（1+1+1+2+2）；P85 旧路径命中 16（1+1+1+1+4+2+1+2+3）；与报告白名单 11 不一致。
6. P89 包三方条目数（独立统计）：payload=9，manifest=9，sums=9。
7. `pack-consistency` 当前文本：P89 payload=5/manifest=5/sums=5（与第 6 条不一致）。
8. 归档总量：P84=44，P85=66，总计 110；两旧目录存在性：False/False。
9. 抽查 5 件哈希：SAMPLE_OK。
10. 活跃引用抽查 3 处目标：均 EXISTS。
11. 冻结面：`git diff --cached --numstat` 对 P85 `00-briefing.md` = 4 0（只增不删）；`CHECKLIST_v72{,_EN}.md` staged diff 为空。

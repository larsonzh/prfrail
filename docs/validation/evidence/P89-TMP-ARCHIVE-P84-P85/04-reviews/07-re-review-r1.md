# P89 ⑦ 终审复审 r1 原文（`GPT-5.3-Codex` 常规档，第 2 次 = 复审，结论 PASS WITH FIXES）

```
FINAL REVIEW: PASS WITH FIXES
```

## 1. SEVERITY 表

| SEVERITY | 发现 | 证据 | 最小修复 |
|---|---|---|---|
| MEDIUM | 证据快照口径不一致：P89 包在 pack-consistency 仍记为 hash_mismatch=1，而实时复算已为 0 | `docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/pack-consistency.txt` 显示 P89: payload=9/manifest=9/sums=9 且 hash_mismatch=1；同轮自跑 verify 显示 P89 hash_mismatch=0、drift=0、verdict=PASS | 仅重生成并回填该文件一次，使 pack-consistency 与 verify 终态一致（P89 hash_mismatch=0） |

## 2. 上轮 3 项逐条核验

1. **#1 已闭环**。两条旧路径在 index 命中集合为 16 行，且包含上轮新增背书的 5 行：`01-arch-design.md` 的 32/61 与 `06-scan-r1.md` 的 25/31/62；同时报告白名单段已写入“16 行”口径。
2. **#2 未闭环**（差 1 个文件快照一致性）。报告正文 §3.1 已是 9/9/9，verify 也为三包三方相等且漂移 0；但 `pack-consistency` 文件内 P89 仍是 hash_mismatch=1（指向 `input-list.md`）。
3. **#3 已闭环**。源分布口径为 44/66，且归档目录自数一致：P84 archive=44、P85 archive=66。

## 3. 终态读数（自跑）

1. `gate --all --scope=index` / `--scope=tree`：均 TOTAL_FAIL=0，changed=133。
2. `gate --selftest`：PASS 744/744。
3. `gate --check=G8 --scope=index`：TOTAL_FAIL=0，G8-c 为 exempt（3 exempt）。
4. 三包复算：P84 58、P85 95、P89 9，逐条 sha256 漂移 0；但 `pack-consistency` 现存快照仍有 P89 hash_mismatch=1。
5. P85 `00-briefing.md` 暂存差异为只增不删（numstat 4/0）；`CHECKLIST_v72{,_EN}.md` 无暂存改动。
6. 中英报告 LF 计数相等（309/309，CR=0/0）；`artifacts` 条目数 133，等于 staged 变更集 133。

## 4. 结论

当前不建议直接提交；先做且只做 1 项最小修复：重生成并回填 `pack-consistency.txt` 使 P89 `hash_mismatch=0`，与 verify 终态一致后即可按 PASS 提交。

---

> **主控事后记录**：该唯一修复已执行——`verify.js` 重跑 → 替换包内 `pack-consistency.txt`（P84/P85/P89 = 58/95/9，三包 `hash_mismatch=0`、`set_issues=0`）→ `meta.js` 重建 `SHA256SUMS.txt` → 复跑 `verify.js` 确认 `hit_lines=16 / whitelist=16 / unexpected=0 / verdict=PASS`。⑦ 额度 1+1 已用尽，未开第 3 轮。

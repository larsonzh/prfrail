# P86 ⑦ 盲审（§7.7 必抽，终审前）逐字转录 · 2026-11-22

- 载体：`independent-reviewer`（只读）；模型 `GPT-5.3-Codex (copilot)`，常规档。
- 额度：盲审必抽 1 次（另计 +1、终审前执行）；判定理由由 ① 在启动前给出（bounddir 触身份/所有权语义）并落盘于方案书 §7。
- 盲审纪律：输入包不含 ⑤/⑥ 审查发现、不含 ① 方案书裁定。
- 消毒声明（§7.3）：输入包不含凭据 / token / 私密主机名 / IP。

---

FINAL REVIEW: PASS WITH FIXES

发现表
|级别|位置|发现|最小修复|
|---|---|---|---|
|Low|`internal/evidence/p86_bounddir_unc_test.go:169`、`tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/test-post-restore.txt:140`|生产入口用例在当前证据里为 SKIP（权限不足），因此“经生产入口放行”尚缺非 SKIP 实证。|补 1 条管理员环境的 go test -v 日志，出现 TestP86UncShareRootProductionEntry 的 PASS（或 CI 等效证据）。|

1. 关于“ASCII 大小写折叠是否扩大放行对象集合”：结论是不扩大对象集合，只扩大同一 UNC 共享根的可接受拼写集合。依据是分支只在 share-root 体位触发（`bounddir_windows.go:997`、`:1046`、`:1071`），且仅做 ASCII 字母折叠；非 ASCII 与其余字节仍逐字节比较。按 UNC/SMB 语义，主机名与共享名大小写不区分，不会引入“仅大小写不同但对象不同”的新目标。
2. 界外放行逐类审计：未见可放行界外目标。不同主机/不同共享/尾空白共享名/非 ASCII 大小写对/大小写变体后代/cross share（如 `\\host\share2\sub`）由 share-root 比较失败后落入严格后代前缀守卫拒绝（`bounddir_windows.go:997`）；host-only（`\\host`, `\\host\`）不满足 share-root 分类并被后续绝对/前缀守卫拒绝；`\\?\` 与 `\\.\` 在 share-root 分类被显式排除并不会进入等值分支（`bounddir_windows.go:1046`）；`..`、重复分隔符、非根尾分隔符、正斜杠先被规范形守卫拒绝（`bounddir_windows.go:939`）。
3. 与 P84 Ordinal 纪律冲突性：未冲突，未削弱。W1 平面目录大小写孪生拒绝仍在独立平面，P86 分支未触及该比较面（`bounddir_windows.go:1650`）；且有“case-variant-descendant-stays-ordinal”反例钉住（`p86_bounddir_unc_test.go:97`）。
4. W10 注释与 r3 相容性：相容，且已把 M1/M2 分开标注。现注释明确“r3 记录的是单分隔前缀形态，双分隔是叠加变异”，与 r3 原文一致（`p84_bounddir_residual_test.go:918`、`07-re-review-r3.md:41`）。
5. 测试证明力：分支必要性已由变异实跑证明（删分支后 P86 正例红、恢复后绿，见 `tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/mutation-record.txt:4`、`:7`、`mutation-branch-removed.txt:1`）；边界反例覆盖充分（`p86_bounddir_unc_test.go:97`）。生产入口 SKIP 不是代码安全问题，但在“必须证明生产入口放行”的证据口径下属于可修复的 Low 证据缺口。
6. Medium 及以上未决：未发现。

可证伪性审计
1. 删除 UNC 等值分支后，三条拼写差异正例应立即失败；恢复后应回到全绿（已见实跑证据）。
2. 若把规范形守卫放宽，repeated-separator、`..`、non-root trailing separator 负例应转红。
3. 若把等值分支误扩到后代，case-variant-descendant-stays-ordinal 负例会转红。

无法验证部分
- 当前输入未提供“非 SKIP 的生产入口 PASS”日志；仅见当前主机因 net share 权限不足而 SKIP（`tmp/P86-BOUNDDIR-UNC-EQUIVALENCE/test-post-restore.txt:143`）。

END OF BLIND REVIEW

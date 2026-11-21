# P86 ⑤ 预审 r1 逐字转录 · 2026-11-22

- 载体：`independent-reviewer`（只读）；模型 `DeepSeek V4 Pro (deepseek)`；额度 1 + 1 的第 1 次。
- 输入包：① 方案书、生产改动、W10 注释、新测试、P84 报告 §5/§8、⑦ 收窄复审原文、⑦ 复审 r3。
- 消毒声明（§7.3）：输入包不含凭据 / token / 私密主机名 / IP。

---

PRE-REVIEW: PASS WITH FIXES

## 分级发现表

| 严重度 | 位置 | 发现 | 最小修复 |
|---|---|---|---|
| Medium | `01-arch-design.md` §5.4 | 生产入口「CI windows 腿必须真跑」的验收证据定为「日志出现 `ok internal/evidence`」——该行**无法区分 PASS 与 SKIP**，而夹具存在三条合法 skip 路径（net share、UNC 句柄、尾分隔符前置）。 | 验收证据改为 `go test -v` 输出含 `--- PASS: TestP86UncShareRootProductionEntry`；若 SKIP 则报告原文登记原因。测试本身 skip 口径可接受。 |
| Low | `bounddir_windows.go` :987-996 | `boundRelativeSuffix` 函数注释仍整体宣称「The comparison is Ordinal…a case-variant twin…is refused」，与新增分支矛盾（共享根 ASCII 大小写变体现被放行）。 | 注释补一句唯一例外：UNC 共享根卷名 ASCII 折叠由 P86 分支处理；W1 平面目录大小写孪生仍拒。 |
| Low | `bounddir_windows.go` :1038-1047 | `uncShareRootName` 注释声称「UNC volume name…and "" otherwise」，但对 `\\?\C:\`、`\\?\UNC\host\share\` 等扩展形态会返回非空卷名（`filepath.VolumeName` 返回该形态）。生产不可达（base 经 `handleFinalPath` DOS 形态 + `stripExtendedPathPrefix`；target 经 `normalizeReparseTarget` 剥离），故无实际放宽，但 A1「恰为」仅在生产层成立。 | 收紧守卫（拒绝 `vol` 以 `\\?\` 起头）或改注释为「volume name of the ordinary DOS form」；并在分类器测试补 `\\?\UNC\...`/`\\?\C:\` 负例。 |
| Low | `01-arch-design.md` §6 | 文档同步清单未含 `SECURITY.md:5` 状态注记与 `bounddir-residual{,_EN}.md` §4/§8 的 N13-UNC「残余/顺延」行——P86 闭合后两者将过时。 | ⑧b 收口时更新上述两处（属登记，非安全基线改动）。 |
| Low | `01-arch-design.md` §5.5 ① | 「必须**恰好**正例变红」措辞含糊（可误读为 5 子例全红）。测试文件注释已精确（恰 3 个拼写差异子例红，identical/strict-descendant 两世界同绿）。 | 方案书措辞对齐测试文件注释。 |

## 逐项结论（审查重点）

1. **正确性**：分支仅放行同长、ASCII 折叠相等的两个共享根卷名；不同主机/共享名（字节或长度不同）不可折叠通过；界外后代 `uncShareRootName` 返空、严格后代 Ordinal 前缀仍拒；中间 reparse 走链机制零改动；非规范拼写被先行 `nonCanonicalReparseTarget` 拒。**交叉**：`\\host\share` 对 `\\host\share\sub` 不变（后者根名返空走后代分支）。
2. **fail-closed**：**无放宽**。四个对照面（`nonCanonicalReparseTarget`、严格后代、`boundRelativeSuffixShortAlias`（UNC 明确拒）、`isUnderDirectory`）diff 均为 0；`D:` 相对形态被 `IsAbs` 先行拒；长度不等拒绝；非 ASCII 逐字节相等。唯一披露点 = 上述 Low 2（扩展形态在助手层可被识别，生产不可达）。
3. **§3 裁定可辩护**：是。同轮用户起跑消息（§六）优先于 P84 ⑦ 遗留建议（准则 §0.1），四要素留痕 +「不得作为先例」符合 OB-37；安全论证成立——SMB 主机/共享名大小写不敏感，同一主机不存在仅大小写不同的两个共享 ⇒ 折叠不新增放行集合、仅消误拒；范围限 ASCII + 卷名。**不需改 T4(b)/W1/SECURITY 规范性内容**（仅状态注记需更新，见 Low 4）。
4. **W10 注释**：两变异（M1 仅回退等值 / M2 叠加未修剪前缀）已分开且归因正确（M1 红 1 例、M2 复合红 3 例）；「`D:\\` = 轮 11 前原形」错误表述已删，改为「NOT the round-11-before arrangement」；与 r3「TrimRight 把基准变成 D:、目标刚好等于单分隔符前缀」逐字相容。
5. **测试充分性**：正例双向 + 折叠 + 恒等 + 后代回归齐备；反例覆盖主机/共享/后代大小写/重复分隔符/非根尾分隔符/斜杠/`..`/界外 + 助手层非 ASCII 与长度。`identical-trailing-spellings` 两世界同绿**不是缺口**——变异信号由 3 个拼写差异子例承载，测试注释已披露；建议补 `\\?\` 形态分类器负例（Low 2）。生产入口 skip 口径可接受，唯证据判据需升级（Medium 1）。
6. **A1–A4**：均可机械核验。A1 由正/反例表 + 代码封闭性（唯一新增路径 = 步骤 4）支持；A2 由 diff 范围核验；A3 由「删分支→3 子例红、恢复→哈希逐字节一致」核验；A4 由 W10 五例 + W1 全绿核验。变异两步足以证明**必要性**（无他路可放行 3 子例），与反例表（界放行集）合看即可证明充分性。

## 我无法验证的部分

- 未运行任何测试/变异/哈希（只读，且为 Windows build tag 用例）。
- 未独立确认工作树 `git status --porcelain` 是否仅含这三处改动（diff 文件为预先抽取；新测试文件为未跟踪文件，不在 diff 内）。
- 生产入口在 CI windows 腿的实跑（net share 权限、句柄名尾分隔符形态）为预测，未执行。
- `SELFTEST 740/740`、三作用域、全量 `go test` 日志未独立核验。
- `p84SetMountPointReparse` 夹具写 `\??\UNC\...` 并经 `DeviceIoControl` 原样存储为预期，未实机验证。

END OF PRE-REVIEW r1

# P86 ⑥ 独立扫描 r1 逐字转录 · 2026-11-22

- 载体：`independent-reviewer`（只读）；模型 `MAI Code 2 Flash (copilot)`；额度 ≤ 3 的第 1 次。
- 输入：生产改动、新测试、W10 注释、P84 报告 §5/§8、⑦ 收窄复审、① 方案书（§3 / §10）。
- 消毒声明（§7.3）：输入包不含凭据 / token / 私密主机名 / IP。

---

SCAN: FINDINGS

| 严重度 | file:line | 发现 | 最小修复 |
|---|---|---|---|
| Low | `internal/evidence/bounddir_windows.go:1046` | `uncShareRootName` 未验证卷名确为 `\\host\share`；实测 `\\host` / `\\host\` 也会返回非空卷名，故纯函数可构造 `\\HOST` 与 `\\host` 被等值放行。生产入口未实证可达，但分类器不满足“仅共享根”的声明。 | 校验恰有非空 host/share 双段，拒绝 host-only、重复分隔符及空段。 |
| Low | `internal/evidence/p86_bounddir_unc_test.go:97` | 负例未显式覆盖 host-only 与尾空白；变异检验目前只是注释，未保留可执行变异/实跑证据。 | 增加对应负例，并记录一次实际删除分支后的红测日志。 |

**已排除的攻击面**
- `\\host\share\` 与 `\\host\share`：本片原定等值，仅取消误拒，不扩大对象集合。
- 跨主机、跨共享、不同长度：比较长度及逐字节值均不匹配。
- 非 ASCII、UTF-8 多字节、控制字符：`lowerASCII` 不折叠，字节保持精确。
- `\\?\`、`\\.\`：分类器主动返回空，不进入等值分支。
- 后代、`sub` 交叉：后代保持 Ordinal 前缀和组件守卫。
- `..`、重复分隔符、正斜杠、非根尾分隔符：规范形守卫先拒绝。
- 尾空白：与无空白共享名长度不同，不等；与同名尾空白名称仅是同共享名的有/无尾分隔符形态。
- `\` 与 `.`：不由 `lowerASCII` 折叠，且 UNC 根名需精确匹配。
- 生产入口：句柄基准经 `stripExtendedPathPrefix`，存储目标经 `normalizeReparseTarget`；实时形态不会因 `\\?\` 提前误拒。

**无法验证部分**
- 本机生产入口因 `net share` 返回拒绝访问而明确 SKIP；纯函数测试与 P84 W10 均通过。
- 未实际执行分支删除变异，不能现场证明每项红测日志。
- 未在真实 UNC 环境运行大小写、扩展/设备及多主机生产路径。

定向及完整 `internal/evidence` 测试均通过；未触 F3/F6/F8、`go.mod`、schema、CLI、门禁或准则。

END OF SCAN r1
